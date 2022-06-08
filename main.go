package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kevinburke/ssh_config"
	"github.com/ngrok/ngrok-api-go/v4"
	"github.com/ngrok/ngrok-api-go/v4/tunnels"
	"ngrok_url/utils"
	"os"
	"path/filepath"
	"strings"
)

var (
	onlyPrint         bool
	sshConfigLocation string
)

func main() {

	flag.BoolVar(&onlyPrint, "p", false, "Only print the tunnel urls")
	flag.StringVar(&sshConfigLocation, "c", "", "Location of the ssh config file")
	flag.Parse()
	if sshConfigLocation != "" && onlyPrint {
		fmt.Println("You can't use both the -p and -c flags")
		os.Exit(1)
	} else {
		sshConfigLocation = filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	}

	var hosts = make([]*utils.Host, 0)
	var allTokens = readTokenFromFile(".token")
	for _, token := range allTokens {
		tmpHost := utils.NewHost(token[1], token[2], listTunnelUrls(context.Background(), token[0]))
		hosts = append(hosts, tmpHost)
		fmt.Println("Process name:  " + tmpHost.Name)
		fmt.Println(fmt.Sprintf("ssh -p %s %s@%s ", tmpHost.Port, tmpHost.Usr, tmpHost.Host))
	}
	if !onlyPrint {
		editConfig(hosts)
	}
}

func readTokenFromFile(filePath string) [][]string {
	token, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var res [][]string
	// Split on newline, then split on comma
	for _, line := range strings.Split(string(token), "\n") {
		if line == "" {
			continue
		}
		res = append(res, strings.Split(line, ","))
	}
	return res
}

func listTunnelUrls(ctx context.Context, token string) *utils.ApiRes {
	// construct the api client
	clientConfig := ngrok.NewClientConfig(token)
	var out string
	var port string
	// list all online tun
	tun := tunnels.NewClient(clientConfig)
	iter := tun.List(nil)
	for iter.Next(ctx) {
		url := iter.Item().PublicURL
		out = strings.TrimPrefix(url, "tcp://")
		totPort := strings.Split(out, ":")
		out = totPort[0]
		port = totPort[1]

	}
	if err := iter.Err(); err != nil {
		return &utils.ApiRes{Error: err}
	}
	return &utils.ApiRes{Host: out, Port: port}
}

func editConfig(allHosts []*utils.Host) {

	f, _ := os.Open(sshConfigLocation)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	cfg, _ := ssh_config.Decode(f)

	for hostIndex, hosts := range cfg.Hosts {
		// Check if hosts.Patterns[0].String() is in the list of allHosts for the Name
		for _, host := range allHosts {
			if host.Name == hosts.Patterns[0].String() {
				cfg.Hosts = remove(cfg.Hosts, hostIndex)
				fmt.Println("Removed: ", host.Name)
			}
		}
	}

	f, _ = os.OpenFile(sshConfigLocation, os.O_WRONLY, 0644)
	_, _ = f.WriteString(cfg.String())
	for _, host := range allHosts {
		_, _ = f.WriteString("\nHost " + host.Name + "\n     HostName " + host.Host + "\n     Port " + host.Port + "\n     User " + host.Usr + "\n     ServerAliveInterval 300\n     ServerAliveCountMax 3")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
}
func remove(slice []*ssh_config.Host, s int) []*ssh_config.Host {
	return append(slice[:s], slice[s+1:]...)
}

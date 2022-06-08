package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kevinburke/ssh_config"
	"github.com/ngrok/ngrok-api-go/v4"
	"github.com/ngrok/ngrok-api-go/v4/tunnels"
	"os"
	"path/filepath"
	"strings"
)

var (
	onlyPrint         bool
	sshConfigLocation string
)

type Hosts struct {
	ID     string `json:"ID,omitempty"`
	Name   string `json:"Name,omitempty"`
	ApiKey string `json:"ApiKey,omitempty"`
}

func main() {

	flag.BoolVar(&onlyPrint, "p", false, "Only print the tunnel urls")
	flag.StringVar(&sshConfigLocation, "c", filepath.Join(os.Getenv("HOME"), ".ssh", "config"), "Location of the ssh config file")
	flag.Parse()

	var allTokens = readTokenFromFile(".token")
	for _, token := range allTokens {
		fmt.Println("ProcessName: ", token[1])
		fmt.Print("URL & PORT: ")
		out := listTunnelUrls(context.Background(), token[0])
		fmt.Print(out)
		fmt.Println()
		kk := strings.Split(out, ":")
		if !onlyPrint {
			editConfig(token[1], kk[0], kk[1])
		}
	}

	fmt.Println()
}

func readTokenFromFile(file_path string) [][]string {
	token, err := os.ReadFile(file_path)
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

func listTunnelUrls(ctx context.Context, token string) string {
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
		tot_port := strings.Split(out, ":")
		out = tot_port[0]
		port = tot_port[1]

	}
	if err := iter.Err(); err != nil {
		return ""
	}
	return out + ":" + port
}

func editConfig(name string, host string, port string) {
	// Open for ssh config
	// Find the host, and remove it from the config if it exsits
	// Then add the host, to keep valid data
	if name == "" {
		return
	}
	f, _ := os.Open(sshConfigLocation)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	cfg, _ := ssh_config.Decode(f)
	for host_index, hosts := range cfg.Hosts {
		if hosts.Patterns[0].String() == name {
			cfg.Hosts = remove(cfg.Hosts, host_index)
			fmt.Println("Removed: ", name)
		}
	}

	f, _ = os.OpenFile(sshConfigLocation, os.O_WRONLY, 0644)
	_, _ = f.WriteString(cfg.String())
	_, _ = f.WriteString("\nHost " + name + "\n     HostName " + host + "\n     Port " + port + "\n     User maskine\n     ServerAliveInterval 300\n     ServerAliveCountMax 3")
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

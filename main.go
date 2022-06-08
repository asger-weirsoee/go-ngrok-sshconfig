package main

import (
	"context"
	"flag"
	"fmt"
	"ngrok_url/utils"
	"os"
	"path/filepath"
)

var (
	onlyPrint         bool
	sshConfigLocation string
	tokenPath         string
)

func main() {

	flag.BoolVar(&onlyPrint, "p", false, "Only print the tunnel urls")
	flag.StringVar(&sshConfigLocation, "c", "", "Location of the ssh config file")
	flag.StringVar(&tokenPath, "t", ".token", "Location of the tokens file")
	flag.Parse()
	if sshConfigLocation != "" && onlyPrint {
		fmt.Println("You can't use both the -p and -c flags")
		os.Exit(1)
	} else {
		sshConfigLocation = filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	}

	var hosts = make([]*utils.Host, 0)
	var allTokens = utils.ReadTokenFromFile(tokenPath)
	for _, token := range allTokens {
		tmpHost := utils.NewHost(token[1], token[2], utils.ListTunnelUrls(context.Background(), token[0]))
		hosts = append(hosts, tmpHost)
		fmt.Println("Process name:  " + tmpHost.Name)
		fmt.Println(fmt.Sprintf("ssh -p %s %s@%s ", tmpHost.Port, tmpHost.Usr, tmpHost.Host))
	}
	if !onlyPrint {
		fmt.Println("Editing ssh config: ", sshConfigLocation)
		utils.EditConfig(hosts, sshConfigLocation)
	}
}

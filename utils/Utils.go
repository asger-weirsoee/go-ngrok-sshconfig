package utils

import (
	"fmt"
	"github.com/kevinburke/ssh_config"
	"os"
)

func Contains(s []*Host, e *ssh_config.Host) bool {
	for _, a := range s {
		if a.Name == e.Patterns[0].String() {
			return true
		}
	}
	return false
}

func CreateFileIfNotExists(filepath string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		_, err := os.Create(filepath)
		if err != nil {
			panic(err)
		}
	}
}

func EditConfig(allHosts []*Host, sshConfigLocation string) {
	f, _ := os.Open(sshConfigLocation)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	cfg, _ := ssh_config.Decode(f)

	out := make([]*ssh_config.Host, 0)
	for _, hosts := range cfg.Hosts {
		if !Contains(allHosts, hosts) {
			out = append(out, hosts)
		}
	}
	cfg.Hosts = out

	f, _ = os.OpenFile(sshConfigLocation, os.O_WRONLY, 0644)
	_, _ = f.WriteString(cfg.String())
	for index, host := range allHosts {
		_, _ = f.WriteString("Host " + host.Name + "\n     HostName " + host.Host + "\n     Port " + host.Port + "\n     User " + host.Usr + "\n     ServerAliveInterval 300\n     ServerAliveCountMax 3")
		if index+1 != len(allHosts) {
			_, _ = f.WriteString("\n\n")
		}
	}
}

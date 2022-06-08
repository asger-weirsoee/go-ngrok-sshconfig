package utils

import (
	"context"
	"github.com/ngrok/ngrok-api-go/v4"
	"github.com/ngrok/ngrok-api-go/v4/tunnels"
	"os"
	"strings"
)

func ReadTokenFromFile(filePath string) [][]string {
	CreateFileIfNotExists(filePath)
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
	if len(res) == 0 {
		panic("No token found in file")
	}
	return res
}

func ListTunnelUrls(ctx context.Context, token string) *ApiRes {
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
		return &ApiRes{Error: err}
	}
	return &ApiRes{Host: out, Port: port}
}

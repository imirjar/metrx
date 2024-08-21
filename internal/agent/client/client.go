package client

import (
	"context"
	"log"
	"net"

	"github.com/imirjar/metrx/internal/agent/client/http"
)

type Client interface {
	POST(ctx context.Context, body []byte) error
}

func New(secret, crypto, addr string) Client {
	host, err := getMyIP()
	if err != nil {
		log.Print(err)
	}

	log.Print(host)

	client := http.New(secret, crypto, addr, host)
	return client
}

func getMyIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Print("client.go Dial ERROR", err)
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()
	return ip, nil
}

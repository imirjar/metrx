package http

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Client application part witch provide
type Client struct {
	headers map[string]string
	secret  string         // secret is using for hash function which place in HashSHA256 HTTP header
	target  string         // server ip where we will sent request
	host    string         // our ip for X-Real-IP HTTP header
	pk      *rsa.PublicKey // public key for security
	client  *http.Client
	sync.Mutex
}

func New(secret, crypto, target, host string) *Client {

	cli := Client{
		headers: map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
			"X-Real-IP":        "",
		},
		secret: secret,
		host:   host,
		target: target,
		client: &http.Client{Timeout: time.Millisecond * 200},
	}

	if crypto != "" {
		file, err := os.ReadFile(crypto)
		if err != nil {
			log.Print(err)
		}
		block, _ := pem.Decode(file)
		if block == nil || block.Type != "PUBLIC KEY" {
			log.Print(errors.New("failed to decode PEM block containing public key"))
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			log.Print(err)
		}
		switch pub := pub.(type) {
		case *rsa.PublicKey:
			// log.Print("PUB", pub)
			cli.pk = pub
		default:
			log.Print(fmt.Errorf("unexpected key type: %T", pub))
		}
	}
	return &cli
}

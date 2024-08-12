package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/imirjar/metrx/internal/models"
)

func NewClient(secret, crypto, host string) *Client {

	cli := Client{
		secret: secret,
		host:   host,
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

// Client application part witch provide
type Client struct {
	secret string
	host   string
	pk     *rsa.PublicKey
	client *http.Client
	sync.Mutex
}

// Using http.Client to sent our metrics to host+/updates
func (c *Client) POST(ctx context.Context, metrics []models.Metrics) error {

	// Convert metrics to slice of bytes
	ms, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	// fmt.Println("#####1", len(ms))

	// Rewrite compressed version of metrics list in buffer
	var cms bytes.Buffer
	w := gzip.NewWriter(&cms)
	w.Write(ms)
	w.Close()

	if c.pk != nil {

		rng := rand.Reader

		bms := cms.Bytes()
		// hash := encrypt.EncryptSHA256(hex.EncodeToString(ms), c.secret)
		ciphertext, err := rsa.EncryptPKCS1v15(rng, c.pk, bms)
		if err != nil {
			return err
		}

		cms.Reset()
		cms.Write(ciphertext)

	}

	// Create request with necessary params
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/updates", &cms)
	if err != nil {
		log.Print("client.go NewRequestWithContext ERROR", err)
		return err
	}
	// Adding necessary headers

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	ip, err := getMyIP()
	if err != nil {
		log.Println("#######", err)
	}
	log.Print(ip)
	req.Header.Add("X-Real-IP", ip)
	// If secret key exists in .env then we must add encoding header
	// if c.secret != "" {
	// hash := encrypt.EncryptSHA256(hex.EncodeToString(ms), c.secret)
	// 	req.Header.Add("HashSHA256", hex.EncodeToString(hash))
	// }

	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	// log.Print(resp.Status)

	return err
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

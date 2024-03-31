package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

type Client struct {
	Client http.Client
}

func (c *Client) POST(path, secret string, body []byte) error {
	log.Println("client.go SECRET", secret)
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(body)
	if err != nil {
		log.Print("client.go GZIP ERROR")
		return err
	}
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, path, &buf)

	if err != nil {
		log.Print("client.go REQUEST ERROR")
		return err
	}

	if secret != "" {
		hash, err := encrypt.EncryptSHA256(body, []byte(secret))
		log.Println("client.go hash", hex.EncodeToString(hash))
		if err != nil {
			log.Fatal(err)
			return err
		}
		req.Header.Add("HashSHA256", hex.EncodeToString(hash))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print("client.go CLIENT ERROR")
		log.Println(err)
		return err
	}

	return resp.Body.Close()
}

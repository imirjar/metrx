package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

type Client struct {
	Client http.Client
}

func (c *Client) POST(ctx context.Context, path, secret string, body []byte) error {
	log.Println("client.go SECRET", secret)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(body)
	gz.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, &buf)
	if err != nil {
		log.Print("REQUEST ERROR")
		return err
	}

	if secret != "" {
		hash, err := encrypt.EncryptSHA256(body, []byte(secret))
		log.Println("client.go hash", hex.EncodeToString(hash))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("HashSHA256", hex.EncodeToString(hash))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print("CLIENT ERROR")
		return err
	}

	resp.Body.Close()
	return err
}

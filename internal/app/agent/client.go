package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"errors"
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
	_, err := gz.Write(body)
	if err != nil {
		return errors.New("client.go GZIP ERROR")
	}
	gz.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, &buf)
	if err != nil || req == nil {
		return errors.New("client.go REQUEST ERROR")
	}

	if secret != "" {
		hash, err := encrypt.EncryptSHA256(body, []byte(secret))
		if err != nil {
			return errors.New("client.go HASH ERROR")
		}
		req.Header.Add("HashSHA256", hex.EncodeToString(hash))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print("REASON IS", err)
		// log.Println(err)
		return errors.New("client.go CLIENT DO ERROR")
	}
	defer resp.Body.Close()

	return nil
}

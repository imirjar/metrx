package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/imirjar/metrx/pkg/encrypt"
)

type Client struct {
	Client http.Client
	sync.Mutex
}

func (c *Client) POST(ctx context.Context, path, secret string, body []byte) error {
	ctxT, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*200))
	defer cancel()
	c.Lock()
	defer c.Unlock()
	log.Println("client.go SECRET", secret)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(body)
	gz.Close()

	req, err := http.NewRequestWithContext(ctxT, http.MethodPost, path, &buf)
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
		return err
	}
	defer resp.Body.Close()
	log.Print(resp.Status)

	return err
}

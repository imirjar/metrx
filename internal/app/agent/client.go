package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
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
	// log.Println("client.go SECRET", secret)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(body)
	gz.Close()

	req, err := http.NewRequestWithContext(ctxT, http.MethodPost, path, &buf)
	if err != nil {
		// log.Print("client.go NewRequestWithContext ERROR", err)
		return errNewRequestWithContextErr
	}

	if secret != "" {
		hash := encrypt.EncryptSHA256(hex.EncodeToString(body), secret)
		log.Println("HashSHA256", secret, hex.EncodeToString(hash))

		req.Header.Add("HashSHA256", hex.EncodeToString(hash))

		qwert := hmac.New(sha256.New, []byte(secret))
		qwert.Write(body)
		resQWE := qwert.Sum(nil)
		log.Printf(hex.EncodeToString(resQWE))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print(err)
		return errClientDoErr
	}
	defer resp.Body.Close()

	log.Print(resp.Status)

	return err
}

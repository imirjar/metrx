package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/encrypt"
)

func NewClient(secret, host string) *Client {
	return &Client{
		client: &http.Client{Timeout: time.Millisecond * 200},
		secret: secret,
		host:   host,
	}
}

// Client application part witch provide
type Client struct {
	secret string
	host   string
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

	// Rewrite compressed version of metrics list in buffer
	var cms bytes.Buffer
	w := gzip.NewWriter(&cms)
	w.Write(ms)
	w.Close()

	// Create request with necessary params
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/updates", &cms)
	if err != nil {
		log.Print("client.go NewRequestWithContext ERROR", err)
		return err
	}
	// Adding necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	// If secret key exists in .env then we must add encoding header
	if c.secret != "" {
		hash := encrypt.EncryptSHA256(hex.EncodeToString(ms), c.secret)
		req.Header.Add("HashSHA256", hex.EncodeToString(hash))
	}

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

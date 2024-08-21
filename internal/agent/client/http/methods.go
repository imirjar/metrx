package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/encrypt"
)

// Using http.Client to sent our metrics to host+/updates
func (c *Client) POST(ctx context.Context, body []byte) error {

	// Rewrite compressed version of metrics list in buffer
	var cms bytes.Buffer
	w := gzip.NewWriter(&cms)
	w.Write(body)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.target+"/updates", &cms)
	if err != nil {
		log.Print("client.go NewRequestWithContext ERROR", err)
		return err
	}

	// Adding necessary headers
	for h, i := range c.headers {
		req.Header.Add(h, i)
	}

	// If secret key exists in .env then we must add encoding header
	if c.secret != "" {
		hash := encrypt.EncryptSHA256(string(body), c.secret)
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

func (c *Client) MakeBatchRequest(ctx context.Context, batch []models.Metrics) (request, error) {

	var (
		r   request
		err error
	)

	body, err := json.Marshal(batch)
	if err != nil {
		return r, err
	}

	if c.secret != "" {
		r.hash = encrypt.EncryptSHA256(string(body), c.secret)
	}

	if c.pk != nil {
		rnd := rand.Reader
		encBody, err := rsa.EncryptPKCS1v15(rnd, c.pk, body)
		if err != nil {
			return r, nil
		}
		r.body = encBody
	} else {
		r.body = body
	}

	return r, nil
}

type request struct {
	hash []byte
	body []byte
}

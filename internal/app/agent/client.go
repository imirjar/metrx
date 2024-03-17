package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"log"
	"net/http"
)

type Client struct {
	Client http.Client
}

func (c *Client) POST(path string, body []byte, hash ...[]byte) error {
	// log.Print("POST")
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(body)
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, path, &buf)
	if err != nil {
		log.Print("REQUEST ERROR")
		return err
	}

	if len(hash) > 0 {
		log.Print("IS HASH")
		secret := hex.EncodeToString(hash[0])
		log.Print(secret)
		req.Header.Add("HashSHA256", secret)
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

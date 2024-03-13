package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"net/http"
)

type Client struct {
	Client http.Client
}

func (c *Client) POST(path string, body []byte, hash ...[]byte) error {

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(body)
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, path, &buf)
	if err != nil {
		return err
	}

	if len(hash) > 0 {
		req.Header.Add("HashSHA256", hex.EncodeToString(hash[0]))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := c.Client.Do(req)

	if err != nil {
		return err
	}

	resp.Body.Close()
	return err
}

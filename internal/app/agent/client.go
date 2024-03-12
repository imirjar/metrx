package agent

import (
	"bytes"
	"compress/gzip"
	"net/http"
)

type Client struct {
	Client http.Client
}

func (c *Client) POST(path string, metrics []byte) error {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(metrics)
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, path, &buf)
	if err != nil {
		return err
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

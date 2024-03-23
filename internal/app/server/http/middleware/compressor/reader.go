package compressor

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

func NewCompressReader(req *http.Request) (*http.Request, error) {
	log.Println("reader.go", req.Method)
	reqIsZipped := strings.Contains(req.Header.Get("Content-Encoding"), "gzip")
	if reqIsZipped {
		unzipR := req
		zipBody, err := gzip.NewReader(req.Body)
		if err != nil {
			return nil, err
		}

		unzipR.Body = zipBody

		return unzipR, nil
	}

	return req, nil

}

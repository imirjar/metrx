package compressor

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

func NewCompressReader(r *http.Request) (*http.Request, error) {
	log.Println("reader.go", r.Method)
	if r.Method == "POST" {
		reqIsZipped := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if reqIsZipped {
			log.Println("reader.go reqIsZipped", r.Method)
			zipBody, err := gzip.NewReader(r.Body)
			if err != nil {
				return nil, err
			}
			r.Body = zipBody
		}
	}
	return r, nil

}

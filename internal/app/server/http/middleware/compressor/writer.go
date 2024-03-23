package compressor

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

func NewCompressWriter(w http.ResponseWriter, r *http.Request) http.ResponseWriter {

	log.Println("writer.go True", r.Method)
	wMustBeZip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
	if wMustBeZip {
		cw := &compressWriter{
			w:    w,
			zipW: gzip.NewWriter(w),
		}
		defer cw.Close()
		return cw
	}
	return w
}

type compressWriter struct {
	w    http.ResponseWriter
	zipW *gzip.Writer
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zipW.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Add("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zipW.Close()
}

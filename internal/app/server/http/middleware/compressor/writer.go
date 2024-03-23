package compressor

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

func NewCompressWriter(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	if r.Method == "POST" {
		log.Println("writer.go True", r.Method)
		wMustBeZip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if wMustBeZip {
			cw := &compressWriter{
				w:    w,
				zipW: gzip.NewWriter(w),
			}
			defer cw.Close()
			return cw
		}
	} else {
		log.Println("writer.go False", r.Method)
	}

	return w

}

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
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
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zipW.Close()
}

// func (c compressReader) Read(p []byte) (n int, err error) {
// 	return c.gzReq.Read(p)
// }

// func (c *compressReader) Close() error {
// 	if err := c.req.Close(); err != nil {
// 		return err
// 	}
// 	return c.gzReq.Close()
// }

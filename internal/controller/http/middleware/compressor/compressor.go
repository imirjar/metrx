package compressor

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

func Compressing() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

			acceptEncoding := req.Header.Get("Accept-Encoding")
			contentEncoding := req.Header.Get("Content-Encoding")

			//client can read
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			sendsGzip := strings.Contains(contentEncoding, "gzip")

			if supportsGzip && sendsGzip {
				log.Println("supportsGzip")
				cResp := NewCompressWriter(resp)
				defer cResp.Close()
				resp = cResp
			}

			if sendsGzip {
				log.Println("sendsGzip")
				cr, err := NewCompressReader(req.Body)
				if err != nil {
					resp.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer cr.Close()
				req.Body = cr
			}

			next.ServeHTTP(resp, req)

		})
	}
}

type CompressReader struct {
	Req   io.ReadCloser
	GzReq *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		Req:   r,
		GzReq: zr,
	}, nil
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		Resp:   w,
		GzResp: gzip.NewWriter(w),
	}
}

// CompressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type CompressWriter struct {
	Resp   http.ResponseWriter
	GzResp *gzip.Writer
}

func (c *CompressWriter) Header() http.Header {
	return c.Resp.Header()
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.GzResp.Write(p)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.Resp.Header().Set("Content-Encoding", "gzip")
	}
	c.Resp.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *CompressWriter) Close() error {
	return c.GzResp.Close()
}

func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.GzReq.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.Req.Close(); err != nil {
		return err
	}
	return c.GzReq.Close()
}

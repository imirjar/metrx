package compressor

import (
	"compress/gzip"
	"net/http"
)

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		resp:   w,
		gzResp: gzip.NewWriter(w),
	}
}

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	resp   http.ResponseWriter
	gzResp *gzip.Writer
}

func (c *compressWriter) Header() http.Header {
	return c.resp.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.gzResp.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.resp.Header().Set("Content-Encoding", "gzip")
	}
	c.resp.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.gzResp.Close()
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.gzReq.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.req.Close(); err != nil {
		return err
	}
	return c.gzReq.Close()
}

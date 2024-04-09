package compressor

import (
	"compress/gzip"
	"net/http"
)

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

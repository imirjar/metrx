package middleware

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/imirjar/metrx/pkg"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	body, err := pkg.Compress(b)
	if err != nil {
		panic(err)
	}
	return w.Writer.Write(body)
}

func Encoder(next http.Handler) http.Handler {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Accept-Encoding") == "gzip" {

			cmpReader, err := gzip.NewReader(req.Body)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusBadRequest)
				return
			}
			defer cmpReader.Close()
			req.Body = cmpReader

			dcmpWriter := gzip.NewWriter(resp)
			defer dcmpWriter.Close()

			resp.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: resp, Writer: dcmpWriter}, req)
		} else {
			next.ServeHTTP(resp, req)
			return
		}

	})
}

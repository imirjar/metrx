package logger

import (
	"log"
	"net/http"
	"time"
)

type (
	MetaData struct {
		host        string
		Status      int
		method      string
		ContentType string
		Size        int
	}

	loggedW struct {
		http.ResponseWriter
		MetaData *MetaData
	}
)

func Logger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			method := r.Method
			contentType := r.Header.Get("Content-Type")

			var lw = loggedW{ResponseWriter: w}
			metaData := &MetaData{
				method:      method,
				host:        r.Host,
				ContentType: contentType,
				Size:        0,
			}

			lw.MetaData = metaData

			next.ServeHTTP(lw, r)
			duration := time.Since(start)

			log.Println(lw.MetaData.host, duration, lw.MetaData.method, lw.MetaData.ContentType, lw.MetaData.Size)
		})
	}
}

func (lw loggedW) Write(b []byte) (int, error) {
	size, err := lw.ResponseWriter.Write(b)
	lw.MetaData.Size += size
	return size, err
}

func (lw loggedW) WriteHeader(statusCode int) {
	lw.ResponseWriter.WriteHeader(statusCode)
}

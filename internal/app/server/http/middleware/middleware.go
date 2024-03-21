package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/imirjar/metrx/internal/app/server/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/logger"
	"github.com/imirjar/metrx/pkg/encrypt"
)

type Middleware struct {
	*http.Request
	http.ResponseWriter
}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Compressing() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			zipR, err := compressor.NewCompressReader(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			zipW := compressor.NewCompressWriter(w, r)

			next.ServeHTTP(zipW, zipR)
			// supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
			// sendsGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")

			// if supportsGzip {
			// 	cResp := compressor.NewCompressWriter(w)
			// 	defer cResp.Close()
			// 	w = cResp
			// }

			// if sendsGzip {
			// 	cr, err := compressor.NewCompressReader(r.Body)
			// 	if err != nil {
			// 		w.WriteHeader(http.StatusInternalServerError)
			// 		return
			// 	}
			// 	defer cr.Close()
			// 	m.Request.Body = cr
			// }

			// next.ServeHTTP(w, r)

		})
	}
}

func (m *Middleware) Encrypting(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// log.Print(r.URL.Path)

			// if strings.Contains(r.URL.Path, "/update") {
			if r.Method == "POST" {
				headerHash := r.Header.Get("HashSHA256")
				if headerHash == "" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// log.Print("router.go ", key)
				// log.Println("router.go headerHash", headerHash)

				body, err := io.ReadAll(r.Body)
				defer r.Body.Close()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				hash, err := encrypt.EncryptSHA256(body, []byte(key)) //h.cfg.SECRET
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				bodyHash := hex.EncodeToString(hash)
				w.Header().Set("HashSHA256", bodyHash)
				if bodyHash != headerHash {
					log.Printf("Заголовк%s", headerHash)
					log.Printf("Тело зап%s", bodyHash)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))

				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func (m *Middleware) Logging() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			method := r.Method

			//data for logging
			responseData := &logger.ResponseData{
				Status: 0,
				Size:   0,
			}

			m.ResponseWriter = &logger.LoggedResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				ResponseData:   responseData,
			}

			next.ServeHTTP(m.ResponseWriter, r)

			duration := time.Since(start)

			reqLog := log.WithFields(log.Fields{
				"URI":      r.RequestURI,
				"method":   method,
				"duration": duration,
			})
			reqLog.Info("request")

			respLog := log.WithFields(log.Fields{
				"status": responseData.Status,
				"size":   responseData.Size,
			})
			respLog.Info("response")
		})
	}
}

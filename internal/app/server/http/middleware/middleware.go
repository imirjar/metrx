package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
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
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

			acceptEncoding := req.Header.Get("Accept-Encoding")
			contentEncoding := req.Header.Get("Content-Encoding")

			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			sendsGzip := strings.Contains(contentEncoding, "gzip")

			if supportsGzip {
				cResp := compressor.NewCompressWriter(resp)
				defer cResp.Close()
				resp = cResp
			}

			if sendsGzip {
				cr, err := compressor.NewCompressReader(req.Body)
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

func (m *Middleware) Encrypting(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == "POST" && strings.Contains(r.URL.Path, "/updates") {
				req := r
				if key != "" {
					body, err := io.ReadAll(req.Body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					headerHash := req.Header.Get("HashSHA256")
					hashByte, err := encrypt.EncryptSHA256(body, []byte(key)) //h.cfg.SECRET

					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					computedHash := hex.EncodeToString(hashByte)

					if headerHash != computedHash {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					req.Body = io.NopCloser(bytes.NewReader(body))
				}
			}

			next.ServeHTTP(w, r)

			// headerHash := r.Header.Get("HashSHA256")
			// if r.Method == "POST" {
			// 	headerHash := r.Header.Get("HashSHA256")
			// 	if headerHash == "" {
			// 		w.WriteHeader(http.StatusInternalServerError)
			// 		return
			// 	}

			// 	body, err := io.ReadAll(r.Body)
			// 	defer r.Body.Close()
			// 	if err != nil {
			// 		w.WriteHeader(http.StatusInternalServerError)
			// 		return
			// 	}

			// 	hash, err := encrypt.EncryptSHA256(body, []byte(key)) //h.cfg.SECRET
			// 	if err != nil {
			// 		w.WriteHeader(http.StatusInternalServerError)
			// 		return
			// 	}

			// 	bodyHash := hex.EncodeToString(hash)
			// 	w.Header().Set("HashSHA256", bodyHash)
			// 	if bodyHash != headerHash {
			// 		log.Printf("Заголовк%s", headerHash)
			// 		log.Printf("Тело зап%s", bodyHash)
			// 		w.WriteHeader(http.StatusInternalServerError)
			// 		return
			// 	}
			// 	r.Body = io.NopCloser(bytes.NewBuffer(body))

			// 	next.ServeHTTP(w, r)
			// } else {
			// 	next.ServeHTTP(w, r)
			// }
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

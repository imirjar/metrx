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
	"github.com/imirjar/metrx/pkg/encrypt"
)

type Middleware struct {
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggedResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

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
			headerHash := r.Header.Get("HashSHA256")

			if headerHash != "" {
				log.Print("Безопасный запрос")
				body, err := io.ReadAll(r.Body)
				defer r.Body.Close()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				hash, err := encrypt.EncryptSHA256(body, []byte("SHA-256")) //h.cfg.SECRET
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				bodyHash := hex.EncodeToString(hash)
				if bodyHash != headerHash {
					log.Printf("Заголовк%s", headerHash)
					log.Printf("Тело зап%s", bodyHash)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				log.Printf("HASH IS EQUAL")
				r.Body = io.NopCloser(bytes.NewBuffer(body))

				log.Print("ХЭШ равен")
				next.ServeHTTP(w, r)

			} else {
				// log.Print("Небезопасный запрос")
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

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			loggedResp := loggedResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}
			next.ServeHTTP(&loggedResp, r)
			duration := time.Since(start)

			reqLog := log.WithFields(log.Fields{
				"URI":      r.RequestURI,
				"method":   method,
				"duration": duration,
			})
			reqLog.Info("request")

			respLog := log.WithFields(log.Fields{
				"status": responseData.status,
				"size":   responseData.size,
			})
			respLog.Info("response")
		})
	}
}

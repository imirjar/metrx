package middleware

import (
	"bytes"
	"encoding/hex"
	"errors"
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
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.Request = r

			supportsGzip := strings.Contains(m.Request.Header.Get("Accept-Encoding"), "gzip")
			sendsGzip := strings.Contains(m.Request.Header.Get("Content-Encoding"), "gzip")

			if supportsGzip {
				cResp := compressor.NewCompressWriter(w)
				defer cResp.Close()
				m.ResponseWriter = cResp
			}

			if sendsGzip {
				cr, err := compressor.NewCompressReader(m.Request.Body)
				if err != nil {
					m.ResponseWriter.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer cr.Close()
				m.Request.Body = cr
			}

			next.ServeHTTP(m.ResponseWriter, m.Request)
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

				// wr := w
				// hw := EncryptWriter{
				// 	ResponseWriter: w,
				// 	w:              w,
				// 	key:            key,
				// }

				// wr = hw
				// defer hw.Close()
				// next.ServeHTTP(wr, r)
				// log.Print("ХЭШ равен")

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

type EncryptWriter struct {
	http.ResponseWriter
	w   io.Writer
	key string
}

func (hw EncryptWriter) Write(b []byte) (int, error) {
	hash, err := encrypt.EncryptSHA256(b, []byte(hw.key))
	if err != nil {
		return 0, err
	}
	log.Print(hex.EncodeToString(hash))
	hw.Header().Add("HashSHA256", hex.EncodeToString(hash))
	return hw.w.Write(b)
}
func (hw *EncryptWriter) Close() error {
	if c, ok := hw.w.(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("middlewares: io.WriteCloser is unavailable on the writer")
}

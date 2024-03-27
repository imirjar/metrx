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

			// if r.Method == "POST" && strings.Contains(r.URL.Path, "/updates") {
			// 	// log.Print("POST")
			// 	if key != "" {
			// 		log.Println("SECRET", key)
			// 		// log.Print("Key")
			// 		body, err := io.ReadAll(r.Body)
			// 		// log.Print(body)
			// 		if err != nil {
			// 			w.WriteHeader(http.StatusInternalServerError)
			// 			return
			// 		}

			// 		headerHash := r.Header.Get("HashSHA256")
			// 		hashByte, err := encrypt.EncryptSHA256(body, []byte(key)) //h.cfg.SECRET

			// 		if err != nil {
			// 			w.WriteHeader(http.StatusInternalServerError)
			// 			return
			// 		}
			// 		computedHash := hex.EncodeToString(hashByte)

			// 		if headerHash != computedHash || headerHash == "" {
			// 			log.Print(computedHash)
			// 			w.WriteHeader(http.StatusInternalServerError)
			// 			return
			// 		}
			// 		r.Body = io.NopCloser(bytes.NewReader(body))
			// 		// log.Print(r.Body)
			// 	}
			// }

			// next.ServeHTTP(w, r)
			if key != "" {
				headerHash := r.Header.Get("HashSHA256")
				if headerHash != "" {
					log.Println("SECRET", key)
					// log.Print("Key")
					body, err := io.ReadAll(r.Body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					hashByte, err := encrypt.EncryptSHA256(body, []byte(key)) //h.cfg.SECRET
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					computedHash := hex.EncodeToString(hashByte)

					if headerHash != computedHash {
						log.Print(computedHash)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					r.Body = io.NopCloser(bytes.NewReader(body))
					// log.Print(r.Body)
				}
			}
			next.ServeHTTP(w, r)

			// hashSHA256 := r.Header.Get("HashSHA256")
			// log.Print(hashSHA256)
			// if hashSHA256 == "" {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }
			// h := hmac.New(sha256.New, []byte(key))
			// b, err := io.ReadAll(r.Body)
			// if err != nil {
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	return
			// }
			// if _, err := h.Write(b); err != nil {
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			// d := h.Sum(nil)
			// hh, err := hex.DecodeString(hashSHA256)
			// log.Println(hh, d)
			// if err != nil {
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			// if !hmac.Equal(d, hh) {
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	return
			// }
			// r.Body = io.NopCloser(bytes.NewReader(b))
			// next.ServeHTTP(w, r)
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

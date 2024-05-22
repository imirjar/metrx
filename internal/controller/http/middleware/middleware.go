package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/imirjar/metrx/internal/controller/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/controller/http/middleware/logger"
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

			//client can read
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			sendsGzip := strings.Contains(contentEncoding, "gzip")

			if supportsGzip && sendsGzip {
				log.Println("supportsGzip")
				cResp := compressor.NewCompressWriter(resp)
				defer cResp.Close()
				resp = cResp
			}

			if sendsGzip {
				log.Println("sendsGzip")
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
		fn := func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			hashHeader := r.Header.Get("HashSHA256")
			if key != "" && hashHeader != "" {
				hashByte := encrypt.EncryptSHA256(hex.EncodeToString(body), key) //h.cfg.SECRET

				qwert := hmac.New(sha256.New, []byte(key))
				qwert.Write(body)

				if hashHeader != hex.EncodeToString(hashByte) {
					// w.WriteHeader(http.StatusInternalServerError)
					http.Error(w, "", http.StatusInternalServerError)
					resQWE := qwert.Sum(nil)
					log.Printf("key %s hashHeader: %s hash: %s may be: %s", key, hashHeader, hex.EncodeToString(hashByte), hex.EncodeToString(resQWE))

					return
				} else {
					log.Printf("HASH IS EQUAL! ALL RIGHT!")
				}

			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)

		// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if key != "" {

		// 	headerHash := r.Header.Get("HashSHA256")
		// 	body, err := io.ReadAll(r.Body)
		// 	if err != nil {
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		return
		// 	}

		// 	hashByte, err := encrypt.EncryptSHA256(body, []byte(key)) //h.cfg.SECRET
		// 	if err != nil {
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		return
		// 	}

		// 	if headerHash != "" {
		// 		log.Println("SECRET", key)
		// 		// log.Print("Key")

		// 		computedHash := hex.EncodeToString(hashByte)

		// 		if headerHash != computedHash {
		// 			log.Print(computedHash)
		// 			w.WriteHeader(http.StatusInternalServerError)
		// 			return
		// 		}

		// 		// log.Print(r.Body)
		// 	}
		// 	r.Body = io.NopCloser(bytes.NewReader(body))
		// }
		// next.ServeHTTP(w, r)
		// })
	}
}

func (m *Middleware) EncWrite(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			wr := w
			if key != "" {
				hw := hashWriter{
					ResponseWriter: w,
					w:              w,
					key:            key,
				}

				wr = hw
				defer hw.Close()

			}
			next.ServeHTTP(wr, r)
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

type hashWriter struct {
	http.ResponseWriter
	w   io.Writer
	key string
}

func (hw hashWriter) Write(b []byte) (int, error) {
	hashByte := encrypt.EncryptSHA256(hex.EncodeToString(b), hw.key) //h.cfg.SECRET

	hw.Header().Add("HashSHA256", hex.EncodeToString(hashByte))
	return hw.w.Write(b)
}

func (hw *hashWriter) Close() error {
	if c, ok := hw.w.(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("middlewares: io.WriteCloser is unavailable on the writer")
}

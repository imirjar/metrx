package middleware

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

func (m *Middleware) CheckReqestHashHeader(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			hashHeader := r.Header.Get("HashSHA256")
			if key != "" && hashHeader != "" {
				h, err := encrypt.EncryptSHA256(body, []byte(key))
				if err != nil {
					log.Fatal(err)
				}
				hash := hex.EncodeToString(h)
				if hashHeader != hash {
					w.WriteHeader(http.StatusTeapot)
					return
				}

			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

type hashWriter struct {
	http.ResponseWriter
	w   io.Writer
	key string
}

func (hw hashWriter) Write(b []byte) (int, error) {
	hash, err := encrypt.EncryptSHA256(b, []byte(hw.key))
	if err != nil {
		return 0, err
	}
	hw.Header().Add("HashSHA256", hex.EncodeToString(hash))
	return hw.w.Write(b)
}
func (hw *hashWriter) Close() error {
	if c, ok := hw.w.(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("middlewares: io.WriteCloser is unavailable on the writer")
}

func (m *Middleware) ResposeHeaderWithHash(key string) func(http.Handler) http.Handler {
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

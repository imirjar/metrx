package encryptor

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

func Encrypting(key string) func(next http.Handler) http.Handler {
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
	}
}

func EncWrite(key string) func(next http.Handler) http.Handler {
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

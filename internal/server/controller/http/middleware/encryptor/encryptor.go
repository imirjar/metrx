package encryptor

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/imirjar/metrx/pkg/encrypt"
)

func DecryptR(path string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			byteKey, err := os.ReadFile(path)
			if err != nil {
				log.Print(err)
			}

			var hk encrypt.Housekeeper
			hk.LoadPrivateKey(byteKey)

			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Print(err)
			}

			decrBody, err := hk.DecryptByKey(body)
			if err != nil {
				log.Print(err)
			}
			r.Body = io.NopCloser(bytes.NewBuffer(decrBody))

			next.ServeHTTP(w, r)
		})
	}
}

func EcryptW(byteKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			hw := encWriter{
				w: w,
			}

			next.ServeHTTP(hw, r)
		})
	}
}

type encWriter struct {
	w      http.ResponseWriter
	pubKey []byte
}

func (ew encWriter) Header() http.Header {
	return ew.w.Header()
}

func (ew encWriter) Write(b []byte) (int, error) {
	var hk encrypt.Housekeeper
	hk.LoadPublicKey(ew.pubKey)
	newB, err := hk.EncryptByKey(b)
	if err != nil {
		log.Print(err)
	}
	return ew.w.Write(newB)
}

func (ew encWriter) WriteHeader(s int) {
	ew.w.WriteHeader(s)
}

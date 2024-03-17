package http

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

type Encryptor struct{}

func (h *HTTPGateway) Encrypting(next http.Handler) http.Handler {
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

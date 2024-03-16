package encryptor

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

func Encryptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerHash := r.Header.Get("HashSHA256")

		log.Println("Header Hash===>", headerHash)

		if headerHash != "" {
			log.Print("Безопасный запрос")
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			hash, err := encrypt.EncryptSHA256(body, []byte(headerHash))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			bodyHash := hex.EncodeToString(hash)
			if bodyHash != headerHash {
				log.Printf("%s", bodyHash)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			log.Print("ХЭШ равен")
			next.ServeHTTP(w, r)

		} else {
			// log.Print("Небезопасный запрос")
			next.ServeHTTP(w, r)
		}
	})
}

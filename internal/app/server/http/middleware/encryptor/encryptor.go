package encryptor

import (
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

func Encryptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		headerHash := req.Header.Get("HashSHA256")

		if headerHash != "" {
			body, err := io.ReadAll(req.Body)
			defer req.Body.Close()
			if err != nil {
				log.Print("####ERROR BODY")
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}

			hash, err := encrypt.EncryptSHA256(body, "HashSHA256")
			if err != nil {
				log.Print("####ERROR CRYPTO")
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}

			bodyHash := hex.EncodeToString(hash)
			if bodyHash != headerHash {
				log.Printf("%s", bodyHash)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Print("ХЭШ равен")
			next.ServeHTTP(resp, req)

		} else {
			log.Print("Небезопасный запрос")
			next.ServeHTTP(resp, req)
		}
	})
}

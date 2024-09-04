package encryptor

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"log"
	"net/http"
)

func DecryptR(pk *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// log.Println("#####", pk)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//For post metrhods only
			if pk != nil && r.Method == "POST" {
				encryptedData, err := io.ReadAll(r.Body)
				if err != nil {
					log.Print(err)
					http.Error(w, "Failed to read request body", http.StatusInternalServerError)
					return
				}
				defer r.Body.Close()

				decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, pk, encryptedData)

				if err != nil {
					log.Print(err)
					http.Error(w, "Failed to decrypt data", http.StatusInternalServerError)
					return
				}

				r.Body = io.NopCloser(bytes.NewReader(decryptedData))

			}
			next.ServeHTTP(w, r)
		})
	}

}

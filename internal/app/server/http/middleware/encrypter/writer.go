package encrypter

import (
	"encoding/hex"
	"net/http"

	"github.com/imirjar/metrx/pkg/encrypt"
)

type EncWriter struct {
	Key []byte
	W   http.ResponseWriter
}

func (e EncWriter) Header() http.Header {
	return e.W.Header()
}

func (e EncWriter) Write(b []byte) (int, error) {
	hashByte, err := encrypt.EncryptSHA256(b, e.Key) //h.cfg.SECRET
	if err != nil {
		e.WriteHeader(http.StatusInternalServerError)
		return 0, err
	}
	e.W.Header().Set("HashSHA256", hex.EncodeToString(hashByte))
	// log.Println("HASH HASH HASH", hex.EncodeToString(hashByte))
	return e.W.Write(b)
}

func (e EncWriter) WriteHeader(statusCode int) {
	e.W.WriteHeader(statusCode)
}

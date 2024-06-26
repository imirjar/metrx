package encrypt

import (
	"crypto/sha256"
)

func EncryptSHA256(value, key string) []byte {
	h := sha256.New()
	h.Write([]byte(value))
	return h.Sum(nil)
}

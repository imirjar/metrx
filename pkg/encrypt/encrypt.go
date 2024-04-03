package encrypt

import (
	"crypto/sha256"
)

func EncryptSHA256(value, key string) []byte {
	// h := hmac.New(sha256.New, []byte(key))
	// h.Write([]byte(value))
	// res := h.Sum(nil)
	// return res, nil
	h := sha256.New()
	h.Write([]byte(value))
	return h.Sum(nil)
}

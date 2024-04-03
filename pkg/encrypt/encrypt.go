package encrypt

import (
	"crypto/hmac"
	"crypto/sha256"
)

func EncryptSHA256(value, key string) ([]byte, error) {
	// h := sha256.New()
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	// dst := h.Sum(nil)
	// if _, err := h.Write(value); err != nil {
	// 	log.Print(err)
	// 	return value, err
	// }
	res := h.Sum(nil)
	return res, nil
}

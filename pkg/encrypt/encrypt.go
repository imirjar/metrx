package encrypt

import (
	"crypto/sha256"
	"log"
)

func EncryptSHA256(value []byte, key string) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(value); err != nil {
		log.Print(err)
		return value, err
	}
	res := h.Sum(nil)
	return res, nil
}

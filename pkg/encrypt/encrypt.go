package encrypt

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
)

func EncryptSHA256(value, key string) []byte {
	h := sha256.New()
	h.Write([]byte(value))
	return h.Sum(nil)
}

func GetRSA(path string) (*rsa.PublicKey, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("CRYPTO ERR:", err)
		return nil, err
	}
	block, _ := pem.Decode(file)
	if block == nil || block.Type != "PUBLIC KEY" {
		err := errors.New("failed to decode PEM block containing public key")
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println("CRYPTO ERR:", err)
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		// log.Print("PUB", pub)
		return pub, nil
	default:
		log.Println("CRYPTO ERR:", err)
		return nil, err
	}
}

package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

type Housekeeper struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

// LoadPrivateKey from path and convert it to RSA
func (hk *Housekeeper) LoadPrivateKey(keyBytes []byte) error {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		err := fmt.Errorf("failed to decode PEM block containing private key")
		log.Print(err)
		return err
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	switch priv := priv.(type) {
	case *rsa.PrivateKey:
		hk.private = priv
		return nil
	default:
		err := errors.New("not an RSA private key")
		log.Print(err)
		return err
	}
}

// LoadPublicKey from path and convert it to RSA
func (hk *Housekeeper) LoadPublicKey(keyBytes []byte) error {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	switch pub := pub.(type) {

	case *rsa.PublicKey:
		hk.public = pub
		return nil
	default:
		return errors.New("not an RSA public key")
	}
}

func (hk *Housekeeper) DecryptByKey(ciphertext []byte) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, hk.private, ciphertext, nil)
	if err != nil {
		log.Print(err)
		return plaintext, err
	}
	return plaintext, nil
}

func (hk *Housekeeper) EncryptByKey(msg []byte) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, hk.public, msg, nil)
	if err != nil {
		log.Print(err)
		return ciphertext, err
	}
	return ciphertext, nil
}

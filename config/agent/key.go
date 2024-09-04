package config

import (
	"crypto/rsa"
	"encoding/json"

	"github.com/imirjar/metrx/pkg/encrypt"
)

type PKey struct {
	Path string
	Pub  rsa.PublicKey
	Priv *rsa.PrivateKey
}

func (pk *PKey) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &pk.Path)
	if err != nil {
		return err
	}

	rsa, err := encrypt.GetRSA(pk.Path)
	if err != nil {
		return err
	}

	pk.Pub = *rsa

	// log.Println(pk.Pub)
	// log.Println(pk.Path)
	return err
}

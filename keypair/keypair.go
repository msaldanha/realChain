package keypair

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

func New() (*KeyPair, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return &KeyPair{PrivateKey:private.D.Bytes(), PublicKey:pubKey}, nil
}

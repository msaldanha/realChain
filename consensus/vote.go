package consensus

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"github.com/msaldanha/realChain/crypto"
	"strconv"
)

func (m *Vote) Hash() []byte {
	ok := []byte(strconv.FormatBool(m.Ok))
	hashableBytes := [][]byte{ok, []byte(m.Reason)}
	headers := bytes.Join(hashableBytes, []byte{})
	hash := sha256.Sum256(headers)
	return []byte(hex.EncodeToString(hash[:]))
}

func (m *Vote) Sign(privateKey *ecdsa.PrivateKey) error {
	s, err := crypto.Sign(m.Hash(), privateKey)
	if err != nil {
		return err
	}
	m.Signature = s
	return nil
}

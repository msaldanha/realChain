package address

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
	"errors"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Address struct {
}

func New() (*Address) {
	return &Address{}
}

func (addr *Address) GenerateForKey(pubKey []byte) (string, error) {
	pubKeyHash, err := hashPubKey(pubKey)
	if err != nil {
		return "", err
	}

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return string(address), nil
}

func (addr *Address) IsValid(address string) (bool, error) {
	pubKeyHash := Base58Decode([]byte(address))
	var chksum [4]byte
	copy(chksum[:], pubKeyHash[len(pubKeyHash) - addressChecksumLen:])
	if bytes.Compare(checksum(pubKeyHash[:len(pubKeyHash) - addressChecksumLen]), chksum[:]) != 0 {
		return false, errors.New("invalid checksum")
	}
	return true, nil
}

func hashPubKey(pubKey []byte) ([]byte, error) {
	sha256Hash := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(sha256Hash[:])
	if err != nil {
		return nil, err
	}
	return RIPEMD160Hasher.Sum(nil), nil
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
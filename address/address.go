package address

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
	"github.com/msaldanha/realChain/Error"
	"github.com/msaldanha/realChain/keypair"
	"github.com/davecgh/go-xdr/xdr2"
	"fmt"
)

const version = byte(0x00)
const addressChecksumLen = 4

const (
	ErrInvalidChecksum = Error.Error("invalid checksum")
)

type Address struct {
	Keys *keypair.KeyPair
	Address string
}

func New() (*Address) {
	return &Address{}
}

func NewAddressWithKeys() (*Address, error) {
	keys, err := keypair.New()
	if err != nil {
		return nil, err
	}

	return NewAddressForKeys(keys)
}

func NewAddressForKeys(keys *keypair.KeyPair) (*Address, error) {
	addr := &Address{Keys: keys}
	hash, err := generateAddressHash(addr.Keys.PublicKey)
	if err != nil {
		return nil, err
	}
	addr.Address = string(hash)
	return addr, nil
}

func NewAddressFromBytes(a []byte) *Address {
	var acc Address
	decoder := xdr.NewDecoder(bytes.NewReader(a))
	decoder.Decode(&acc)
	return &acc
}

func MatchesPubKey(addr []byte, pubKey []byte) bool {
	hash, err := generateAddressHash(pubKey)
	if err != nil {
		return false
	}
	return bytes.Equal(addr, hash)
}

func IsValid(addr string) (bool, error) {
	pubKeyHash := Base58Decode([]byte(addr))
	var chksum [4]byte
	copy(chksum[:], pubKeyHash[len(pubKeyHash) - addressChecksumLen:])
	chkCalc := checksum(pubKeyHash[:len(pubKeyHash) - addressChecksumLen])
	if bytes.Compare(chkCalc, chksum[:]) != 0 {
		return false, ErrInvalidChecksum
	}
	return true, nil
}

func (a *Address) ToBytes() []byte {
	var result bytes.Buffer
	encoder := xdr.NewEncoder(&result)
	count, err := encoder.Encode(a)
	if err != nil {
		fmt.Printf("Encoded %d, Error: %s", count, err.Error())
	}
	return result.Bytes()
}

func generateAddressHash(pubKey []byte) ([]byte, error) {
	pubKeyHash, err := hashPubKey(pubKey)
	if err != nil {
		return nil, err
	}

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	if ok, _ := IsValid(string(address)); !ok {
		fmt.Printf("WARNING: Not valid addr generated: %s\n", string(address))
	}

	return address, nil
}

func (a *Address) IsValid() (bool, error) {
	return IsValid(a.Address)
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
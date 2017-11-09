package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"encoding/binary"
	"io"
	"math/big"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func CreateKeyPair() (*ecdsa.PrivateKey, error) {
	return CreateKeyPairInt(rand.Reader)
}

func CreateKeyPairInt(reader io.Reader) (*ecdsa.PrivateKey, error) {
	pubkeyCurve := elliptic.P256()

	key := new(ecdsa.PrivateKey)
	key, err := ecdsa.GenerateKey(pubkeyCurve, reader)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GetPublicKeyHash(key ecdsa.PublicKey) string {
	var pk []byte

	pk = append(pk, 0x04)

	pk = append(pk, key.X.Bytes()...)
	pk = append(pk, key.Y.Bytes()...)

	s := sha256.Sum256(pk)

	h := ripemd160.New()
	h.Write(s[:])
	r := h.Sum(nil)

	pk = []byte{16}
	pk = append(pk, r...)

	s = sha256.Sum256(pk)
	s = sha256.Sum256(s[:])
	chk := s[0:4]

	pk = append(pk, chk...)

	return base58.Encode(pk)
}

// Key will be encoded like this:
// type + len(D) + D + len(x) + X + len(y) + y + hash
// type is 1 for private key, 2 for public key
// len are bytes
// hash is last 4 bytes of sha256(len(D) + ... y)
func PrivateKeyToBytes(key ecdsa.PrivateKey) []byte {
	var pk []byte

	pk = append(pk, BigIntToBytes(*key.D)...)
	pk = append(pk, BigIntToBytes(*key.X)...)
	pk = append(pk, BigIntToBytes(*key.Y)...)

	s := sha256.Sum256(pk)

	pk = append(pk, s[0:4]...)

	return append([]byte{0x01}, pk...)
}

// Returns key &  bytes read
func BytesToPrivateKey(b []byte) (ecdsa.PrivateKey, int, error) {
	var key ecdsa.PrivateKey
	idx := 0

	// Read Key.D
	key.D, idx = BytesToBigInt(b, idx)

	// Read Key.X
	key.X, idx = BytesToBigInt(b, idx)

	// Read Key.Y
	key.Y, idx = BytesToBigInt(b, idx)

	// Read hash
	hash := make([]byte, 4)
	copy(hash, b[idx:idx+4])

	// s := sha256.Sum256(b[0:idx])
	idx += 4

	return key, idx, nil
}

func BigIntToBytes(bigInt big.Int) []byte {
	bs := make([]byte, 4)
	size := len(bigInt.Bytes())

	if size > 2^32 {
		panic("Message too big")
	}

	binary.LittleEndian.PutUint32(bs, uint32(size))

	return append(bs, bigInt.Bytes()...)
}

func BytesToBigInt(b []byte, idx int) (*big.Int, int) {
	intsize := 4
	i := new(big.Int)

	size := int(binary.LittleEndian.Uint32(b[idx : idx+intsize]))
	i.SetBytes(b[idx+intsize : idx+intsize+size])

	return i, idx + intsize + size
}

func ReadBigIntFromFile(fd *os.File) (*big.Int, error) {
	bs := make([]byte, 4)
	bigInt := new(big.Int)

	_, err := fd.Read(bs)
	if err != nil {
		fd.Close()
		return nil, err
	}

	size := binary.LittleEndian.Uint32(bs)

	buffer := make([]byte, size)

	_, err = fd.Read(buffer)
	if err != nil {
		fd.Close()
		return nil, err
	}

	bigInt.SetBytes(buffer)

	return bigInt, err
}

func SignMessage(key ecdsa.PrivateKey, message []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, &key, message)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, 0)
	bytes = append(bytes, BigIntToBytes(*r)...)
	bytes = append(bytes, BigIntToBytes(*s)...)

	return bytes, nil
}

func SignVerify(key ecdsa.PublicKey, message []byte, signature []byte) bool {
	r, idx := BytesToBigInt(signature, 0)
	s, _ := BytesToBigInt(signature, idx)

	return ecdsa.Verify(&key, message, r, s)
}

func PublicKeyToBytes(key ecdsa.PublicKey) []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, BigIntToBytes(*key.X)...)
	bytes = append(bytes, BigIntToBytes(*key.Y)...)

	return bytes
}

func GetPublicKeyFromBytes(bytes []byte) ecdsa.PublicKey {
	var idx int
	var key ecdsa.PublicKey

	key.Curve = elliptic.P256()
	key.X, idx = BytesToBigInt(bytes, 0)
	key.Y, idx = BytesToBigInt(bytes, idx)

	return key
}

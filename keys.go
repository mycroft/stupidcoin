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

func WriteKeyToFile(key ecdsa.PrivateKey, filepath string) error {
	// X, Y, D.
	fd, err := os.Create(filepath)
	if err != nil {
		return err
	}

	err = WriteBigIntToFile(fd, *key.PublicKey.X)
	if err != nil {
		return err
	}

	err = WriteBigIntToFile(fd, *key.PublicKey.Y)
	if err != nil {
		return err
	}

	err = WriteBigIntToFile(fd, *key.D)
	if err != nil {
		return err
	}

	fd.Close()

	return nil
}

func WriteBigIntToFile(fd *os.File, bigInt big.Int) error {
	bs := make([]byte, 4)
	size := len(bigInt.Bytes())
	binary.LittleEndian.PutUint32(bs, uint32(size))

	_, err := fd.Write(bs)
	if err != nil {
		return err
	}

	_, err = fd.Write(bigInt.Bytes())
	return err
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

func LoadKeyFromFile(filepath string) (*ecdsa.PrivateKey, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	key := new(ecdsa.PrivateKey)
	key.PublicKey.X, err = ReadBigIntFromFile(fd)
	if err != nil {
		return nil, err
	}

	key.PublicKey.Y, err = ReadBigIntFromFile(fd)
	if err != nil {
		return nil, err
	}

	key.D, err = ReadBigIntFromFile(fd)
	if err != nil {
		return nil, err
	}

	key.Curve = elliptic.P256()

	return key, nil
}

func SignMessage(key ecdsa.PrivateKey, message []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, &key, message)
	if err != nil {
		return nil, err
	}

	var b []byte
	b = append(r.Bytes(), s.Bytes()...)

	return b, nil
}

func SignVerify(key ecdsa.PublicKey, message []byte, signature []byte) bool {
	var r, s big.Int

	r.SetBytes(signature[:32])
	s.SetBytes(signature[32:])

	return ecdsa.Verify(&key, message, &r, &s)
}

func GetPublicKeyHashBytes(key ecdsa.PublicKey) []byte {
	bytes := make([]byte, 64)

	copy(bytes[0:32], key.X.Bytes())
	copy(bytes[32:64], key.Y.Bytes())

	return bytes
}

func GetPublicKeyFromBytes(bytes []byte) ecdsa.PublicKey {
	var x, y big.Int
	var key ecdsa.PublicKey

	x.SetBytes(bytes[:32])
	key.X = &x
	y.SetBytes(bytes[32:])
	key.Y = &y

	key.Curve = elliptic.P256()

	return key
}

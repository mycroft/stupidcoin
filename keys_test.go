package main

import (
	"testing"
)

type FakeReader struct{}

func (r *FakeReader) Read(b []byte) (n int, err error) {
	for i := 0; i < cap(b); i++ {
		b[i] = 0x0
	}

	return cap(b), nil
}

func TestGetPublicKeyHash(t *testing.T) {
	fakeReader := new(FakeReader)
	key, err := CreateKeyPairInt(fakeReader)
	if err != nil {
		t.Errorf("CreateKeyPairInt: %s\n", err)
	}

	hash := GetPublicKeyHash(key.PublicKey)
	hashControl := "7ZRMGSqmCdVW3twY8HJy3vu5JCJDvWb1C8"

	if hashControl != hash {
		t.Errorf("GetPublicKeyHash: Invalid control hash: %s != %s", hashControl, hash)
	}
}

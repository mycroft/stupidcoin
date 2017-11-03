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

func TestReadWriteKey(t *testing.T) {
	// Create key
	key, err := CreateKeyPair()
	if err != nil {
		t.Errorf("Could not create key...")
		return
	}

	filepath := "testkey.dat"
	err = WriteKeyToFile(*key, filepath)
	if err != nil {
		t.Errorf("Could not write key: %s\n", err)
		return
	}

	key_copy, err := LoadKeyFromFile(filepath)
	if err != nil {
		t.Errorf("Could not read key: %s\n", err)
		return
	}

	message := []byte("testaroo")

	s1, err := SignMessage(*key, message)
	if err != nil {
		t.Errorf("Could not sign message: %s\n", err)
		return
	}

	s2, err := SignMessage(*key, message)
	if err != nil {
		t.Errorf("Could not sign message: %s\n", err)
		return
	}

	if !SignVerify(key.PublicKey, message, s1) {
		t.Errorf("Could not validate signature")
		return
	}

	if !SignVerify(key.PublicKey, message, s2) {
		t.Errorf("Could not validate signature")
		return
	}

	if !SignVerify(key_copy.PublicKey, message, s1) {
		t.Errorf("Could not validate signature")
		return
	}

	if !SignVerify(key_copy.PublicKey, message, s2) {
		t.Errorf("Could not validate signature")
		return
	}
}

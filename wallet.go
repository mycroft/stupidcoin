package main

import (
	"crypto/ecdsa"

	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Wallet struct {
	PublicKeys  []ecdsa.PublicKey
	PrivateKeys []ecdsa.PrivateKey
}

func LoadWallet(config Config) (*Wallet, error) {
	w := new(Wallet)

	if _, err := os.Stat(config.Wallet); os.IsNotExist(err) {
		return w, nil
	}

	// If wallet file exist, load it
	// fd, err := os.Open(config.Wallet)
	// if err != nil {
	// 	return w, err
	// }

	bytes, err := ioutil.ReadFile(config.Wallet)
	if err != nil {
		return w, err
	}

	idx := 0
	for idx < len(bytes) {
		// First byte: type
		typ := bytes[idx]
		idx++

		switch typ {
		case 0x01:
			// Read PrivateKey
			key, idxtmp, err := BytesToPrivateKey(bytes[idx:])
			idx += idxtmp
			if err != nil {
				return w, err
			}

			w.AddPrivateKey(key)
		default:
			fmt.Println("Invalid type", typ)
		}
	}

	return w, nil
}

func (w *Wallet) WriteWallet(config Config) error {
	fd, err := os.Create(config.Wallet)
	if err != nil {
		return err
	}
	defer fd.Close()

	for _, key := range w.PrivateKeys {
		pkbytes := PrivateKeyToBytes(key)
		fd.Write(pkbytes)
	}

	return nil
}

func (w *Wallet) AddPrivateKey(key ecdsa.PrivateKey) {
	w.PrivateKeys = append(w.PrivateKeys, key)
}

func (w *Wallet) AddPublicKey(key ecdsa.PublicKey) {
	w.PublicKeys = append(w.PublicKeys, key)
}

func (w *Wallet) List() {

	if len(w.PrivateKeys) > 0 {
		fmt.Println("Private keys")
		for _, key := range w.PrivateKeys {
			fmt.Println(GetPublicKeyHash(key.PublicKey))
		}
	}

	if len(w.PublicKeys) > 0 {
		fmt.Println("Public keys")
		for _, key := range w.PublicKeys {
			fmt.Println(key)
		}
	}
}

func (w *Wallet) GetPublicKeyByHash(hash string) (ecdsa.PublicKey, error) {
	for _, key := range w.PrivateKeys {
		current_hash := GetPublicKeyHash(key.PublicKey)
		if current_hash == hash {
			return key.PublicKey, nil
		}
	}

	for _, key := range w.PublicKeys {
		current_hash := GetPublicKeyHash(key)
		if current_hash == hash {
			return key, nil
		}
	}

	return ecdsa.PublicKey{}, errors.New("Could not find key")
}

func (w *Wallet) GetPrivateKeyByHash(hash string) (ecdsa.PrivateKey, error) {
	for _, key := range w.PrivateKeys {
		current_hash := GetPublicKeyHash(key.PublicKey)
		if current_hash == hash {
			return key, nil
		}
	}

	return ecdsa.PrivateKey{}, errors.New("Could not find key")
}

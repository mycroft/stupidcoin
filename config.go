package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Blockchain    string `json:"blockchain"`
	Wallet        string `json:"wallet"`
	key           ecdsa.PublicKey
	MiningAddr    string `json:"mining-addr"`
	WebListenAddr string `json:"listen-addr"`
}

func LoadConfiguration(path string) (Config, error) {
	var config Config

	// Default values...
	config.Blockchain = ".blocks.dat"
	config.Wallet = "wallet.key"
	config.WebListenAddr = ":8080"

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

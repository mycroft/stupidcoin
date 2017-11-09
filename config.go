package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Wallet     string `json:"wallet"`
	MiningAddr string `json:"mining-addr"`
}

func LoadConfiguration(path string) (Config, error) {
	var config Config

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

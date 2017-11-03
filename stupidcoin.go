package main

import (
	"flag"
	"fmt"
	"os"
)

var flagCreateKey bool
var flagMine bool

func init() {
	flag.BoolVar(&flagCreateKey, "create-key", false, "Create key pair")
	flag.BoolVar(&flagMine, "mine", false, "Mine block")
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if flagCreateKey {
		key, err := CreateKeyPair()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = WriteKeyToFile(*key, "wallet.key")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Your new key hash: %s\n", GetPublicKeyHash(key.PublicKey))

		return
	}

	// Load key
	key, err := LoadKeyFromFile("wallet.key")
	if err != nil {
		fmt.Printf("Could not read key: %s\n", err)
		os.Exit(1)
	}

	if flagMine {
		chain, err := LoadBlockchain()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = chain.MineBlock(key.PublicKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		chain.Dump()

		err = chain.SaveBlockchain()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return
	}

	// Nothing done. Showing options.
	Usage()
}

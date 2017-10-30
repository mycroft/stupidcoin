package main

import (
	"flag"
	"fmt"
	"os"
)

var flagCreateKey bool

func init() {
	flag.BoolVar(&flagCreateKey, "create-key", false, "Create key pair")
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

		os.Exit(0)
	}

	// Nothing done. Showing options.
	Usage()
}

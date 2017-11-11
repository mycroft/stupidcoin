package main

import (
	"flag"
	"fmt"
	"os"
)

var flagConfigFile string
var flagCreateKey bool
var flagListKeys bool
var flagMine, flagDumpChain bool
var flagWeb bool
var flagScan bool

func init() {
	flag.BoolVar(&flagCreateKey, "create-key", false, "Create key pair")
	flag.BoolVar(&flagListKeys, "list-keys", false, "List keys in wallet")
	flag.BoolVar(&flagMine, "mine", false, "Mine block")
	flag.StringVar(&flagConfigFile, "config", "config.json", "Configuration file to use")
	flag.BoolVar(&flagDumpChain, "dump", false, "Dump chain (debug)")
	flag.BoolVar(&flagWeb, "web", false, "Launch API server")
	flag.BoolVar(&flagScan, "scan", false, "Scan blockchain for our funds")
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	config, err := LoadConfiguration(flagConfigFile)
	if err != nil {
		panic(err)
	}

	wallet, err := LoadWallet(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flagCreateKey {
		key, err := CreateKeyPair()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		wallet.AddPrivateKey(*key)

		err = wallet.WriteWallet(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Your new key hash: %s\n", GetPublicKeyHash(key.PublicKey))

		return
	}

	if flagListKeys {
		wallet.List()

		return
	}

	key, err := wallet.GetPublicKeyByHash(config.MiningAddr)
	config.key = key
	if err != nil {
		panic(err)
	}

	chain, err := LoadBlockchain(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flagMine {
		err = chain.MineBlock(key)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		chain.Dump()

		err = chain.SaveBlockchain(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return
	}

	if flagDumpChain {
		chain.Dump()

		return
	}

	if flagScan {
		funds := chain.GetFunds(wallet)
		for _, fund := range funds {
			fmt.Printf("%x: %f\n", fund.txn.hash, fund.txn.outputs[fund.output_id].amount)
		}
		return
	}

	if flagWeb {
		err := WebRun(config, *wallet, chain)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Nothing done. Showing options.
	Usage()
}

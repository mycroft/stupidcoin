package main

import (
	"fmt"
	// "html"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type WebDaemon struct {
	Config     Config
	Mine       chan bool
	Txn        chan *TxnOrder
	Blockchain *Blockchain
	Wallet     Wallet
}

func (wd *WebDaemon) MineHandler(w http.ResponseWriter, r *http.Request) {
	// Mine a block.
	wd.Mine <- true
}

func (wd *WebDaemon) AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Add a transaction to transaction queue
	// Input
	// - A destination address (hash)
	// - An amount.

	r.ParseForm()
	fmt.Println(r.PostForm)

	txnOrder := new(TxnOrder)
	txnOrder.Addr = ""
	txnOrder.Amount = 0

	if _, ok := r.PostForm["amount"]; ok {
		f, err := strconv.ParseFloat(r.PostForm["amount"][0], 64)
		if err != nil {
			fmt.Fprintf(w, "NOT OK")
			return
		}

		txnOrder.Amount = f
	}

	if _, ok := r.PostForm["dest"]; ok {
		txnOrder.Addr = r.PostForm["dest"][0]
	}

	select {
	case wd.Txn <- txnOrder:
	default:
		fmt.Fprintf(w, "NOT QUEUED")
		return
	}

	fmt.Fprintf(w, "OK")
}

func WebRun(config Config, wallet Wallet, chain *Blockchain) error {
	daemon := new(WebDaemon)
	daemon.Blockchain = chain
	daemon.Config = config
	daemon.Wallet = wallet
	daemon.Mine = make(chan bool)
	daemon.Txn = make(chan *TxnOrder)

	// Start mining routine
	go func(wd *WebDaemon) {
		for {
			fmt.Println("Waiting for new request.")
			<-wd.Mine
			fmt.Println("got mining request...")

			wd.Blockchain.MineBlock(config.key)

			err := wd.Blockchain.SaveBlockchain(config)
			if err != nil {
				fmt.Println(err)
			}
		}
	}(daemon)

	// Start transaction creator
	go func(wd *WebDaemon) {
		for {
			txnOrder := <-wd.Txn

			txn, err := wd.Blockchain.CreateTransfertTransaction(wd.Wallet, txnOrder)
			if err != nil {
				fmt.Println(err)
				continue
			}
			wd.Blockchain.QueueTransaction(txn)
		}
	}(daemon)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/mine", daemon.MineHandler)
	router.HandleFunc("/txn/add", daemon.AddTransactionHandler)

	err := http.ListenAndServe(config.WebListenAddr, router)

	return err
}

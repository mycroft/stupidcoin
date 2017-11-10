package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/mux"
)

type WebDaemon struct {
	Mine       chan bool
	Blockchain *Blockchain
}

func (wd *WebDaemon) MineHandler(w http.ResponseWriter, r *http.Request) {
	// Mine a block.
	wd.Mine <- true
}

func (wd *WebDaemon) AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Add a transaction to transaction queue
}

func WebRun(config Config, chain *Blockchain) error {
	daemon := new(WebDaemon)
	daemon.Blockchain = chain
	daemon.Mine = make(chan bool)

	// Start mining routine
	go func(wd *WebDaemon) {
		for {
			fmt.Println("Waiting for new request.")
			<-wd.Mine
			fmt.Println("got mining request...")

			wd.Blockchain.MineBlock(config.key)
		}
	}(daemon)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/mine", daemon.MineHandler)

	err := http.ListenAndServe(":8080", router)

	return err
}

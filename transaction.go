package main

import ()

type TxInput struct {
	txhash []byte
	script []byte
}

type TxOutput struct {
	script []byte
	amount float64
}

type Transaction struct {
	hash    []byte
	inputs  []TxInput
	outputs []TxOutput
}

func CreateTransaction() *Transaction {
	tx := new(Transaction)

	return tx
}

func (tx *Transaction) AddInput(input []byte) {

}

func (tx *Transaction) AddOutput(input []byte) {

}

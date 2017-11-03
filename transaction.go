package main

import ()

type TxInput struct {
	txhash []byte
	script *Script
}

type TxOutput struct {
	script *Script
	amount float64
}

type Transaction struct {
	hash    []byte
	inputs  []*TxInput
	outputs []*TxOutput
}

func CreateTransaction() *Transaction {
	tx := new(Transaction)

	return tx
}

func CreateTxOutput(script *Script, amount float64) *TxOutput {
	output := new(TxOutput)
	output.script = script
	output.amount = amount

	return output
}

func (tx *Transaction) AddInput(input *TxInput) {
	tx.inputs = append(tx.inputs, input)
}

func (tx *Transaction) AddOutput(output *TxOutput) {
	tx.outputs = append(tx.outputs, output)
}

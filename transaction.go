package main

import (
	"crypto/sha256"

	"fmt"
	"os"
)

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

func (tx *Transaction) ComputeHash(update bool) []byte {
	h := sha256.New()

	for _, input := range tx.inputs {
		h.Write(input.txhash)
		h.Write(input.script.data)
	}

	for _, output := range tx.outputs {
		h.Write(output.script.data)
		// XXX Add amount
	}

	hash := h.Sum(nil)

	if update {
		tx.hash = hash
	}

	return hash
}

// XXX to rewrite using bytes...
func (tx *Transaction) SaveTransaction(fd *os.File) error {
	WriteBytesToFd(fd, tx.hash)

	WriteUint32ToFd(fd, uint32(len(tx.inputs)))
	for _, input := range tx.inputs {
		WriteBytesToFd(fd, input.txhash)
		WriteBytesToFd(fd, input.script.data)
	}

	WriteUint32ToFd(fd, uint32(len(tx.outputs)))
	for _, output := range tx.outputs {
		WriteBytesToFd(fd, output.script.data)
		WriteFloat64ToFd(fd, output.amount)
	}

	return nil
}

func CreateTransactionFromFd(fd *os.File) (*Transaction, error) {
	var err error
	var i uint32
	txn := new(Transaction)

	txn.hash, err = ReadBytesFromFd(fd)
	if err != nil {
		return nil, err
	}

	input_cnt, err := ReadUint32FromFd(fd)
	if err != nil {
		return nil, err
	}

	for i = 0; i < input_cnt; i++ {
		input := new(TxInput)
		fmt.Println(input_cnt, i)

		input.txhash, err = ReadBytesFromFd(fd)
		if err != nil {
			return nil, err
		}

		input.script = new(Script)
		input.script.data, err = ReadBytesFromFd(fd)
		if err != nil {
			return nil, err
		}

		txn.AddInput(input)
	}

	output_cnt, err := ReadUint32FromFd(fd)
	if err != nil {
		return nil, err
	}

	for i = 0; i < output_cnt; i++ {
		output := new(TxOutput)

		output.script = new(Script)
		output.script.data, err = ReadBytesFromFd(fd)
		if err != nil {
			return nil, err
		}

		output.amount, err = ReadFloat64FromFd(fd)
		if err != nil {
			return nil, err
		}

		txn.AddOutput(output)
	}

	return txn, nil
}

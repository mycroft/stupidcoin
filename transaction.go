package main

import (
	"crypto/sha256"

	"fmt"
	"math"
	"os"
	"strconv"
	"time"
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
	hash      []byte
	timestamp uint64
	inputs    []*TxInput
	outputs   []*TxOutput
}

func CreateTransaction() *Transaction {
	tx := new(Transaction)
	tx.timestamp = uint64(time.Now().Unix())

	return tx
}

func CreateTxOutput(script *Script, amount float64) *TxOutput {
	output := new(TxOutput)
	output.script = script
	output.amount = amount

	return output
}

func (tx *Transaction) String() string {
	var dump string

	dump += fmt.Sprintf("Txn: %x\n", tx.hash)

	for j := 0; j < len(tx.inputs); j++ {
		dump += fmt.Sprintf("- Input: %v\n",
			tx.inputs[j].script.String())
	}

	for j := 0; j < len(tx.outputs); j++ {
		dump += fmt.Sprintf("- Output: %f - %v\n",
			tx.outputs[j].amount,
			tx.outputs[j].script.String())
	}

	return dump
}

func (tx *Transaction) AddInput(input *TxInput) {
	tx.inputs = append(tx.inputs, input)
	tx.ComputeHash(true)
}

func (tx *Transaction) AddOutput(output *TxOutput) {
	tx.outputs = append(tx.outputs, output)
	tx.ComputeHash(true)
}

func (tx *Transaction) ComputeHash(update bool) []byte {
	h := sha256.New()
	h.Write([]byte(strconv.FormatUint(tx.timestamp, 10)))

	for _, input := range tx.inputs {
		h.Write(input.txhash)
		h.Write(input.script.data)
	}

	for _, output := range tx.outputs {
		h.Write(output.script.data)
		h.Write([]byte(strconv.FormatUint(math.Float64bits(output.amount), 10)))
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
	WriteUint64ToFd(fd, tx.timestamp)

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
	txn := CreateTransaction()

	txn.hash, err = ReadBytesFromFd(fd)
	if err != nil {
		return nil, err
	}

	txn.timestamp, err = ReadUint64FromFd(fd)
	if err != nil {
		return nil, err
	}

	input_cnt, err := ReadUint32FromFd(fd)
	if err != nil {
		return nil, err
	}

	for i = 0; i < input_cnt; i++ {
		input := new(TxInput)

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

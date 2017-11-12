package main

import (
	"crypto/ecdsa"

	"errors"
	"fmt"
	"os"
)

type TxnOrder struct {
	Addr   string
	Amount float64
}

type OutputFund struct {
	txn       *Transaction
	output_id int
	script    *Script
}

type Blockchain struct {
	last_index uint64
	blocks     []*Block

	// Do not store in blockchain
	txnQueue []*Transaction
}

func CreateBlockchain() *Blockchain {
	blockchain := new(Blockchain)
	blockchain.last_index = 0

	return blockchain
}

func LoadBlockchain(config Config) (*Blockchain, error) {
	blockchain := CreateBlockchain()

	if _, err := os.Stat(config.Blockchain); os.IsNotExist(err) {
		fmt.Printf("No existing block chain found...\n")

		return blockchain, nil
	}

	fd, err := os.Open(config.Blockchain)
	if err != nil {
		return nil, err
	}

	blockchain.last_index, err = ReadUint64FromFd(fd)
	if err != nil {
		return nil, err
	}

	for {
		// Read blocks
		block, err := CreateBlockFromFd(fd)
		if err != nil {
			return nil, err
		}

		blockchain.blocks = append(blockchain.blocks, block)

		if block.index == blockchain.last_index {
			break
		}
	}

	return blockchain, nil
}

func (bc *Blockchain) SaveBlockchain(config Config) error {
	fd, err := os.OpenFile(config.Blockchain, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	WriteUint64ToFd(fd, bc.last_index)

	// Write blocks
	for i := 0; i < len(bc.blocks); i++ {
		bc.blocks[i].SaveBlock(fd)
	}

	fd.Close()

	return nil
}

func (bc *Blockchain) MineBlock(key ecdsa.PublicKey) error {
	var last, b *Block

	if len(bc.blocks) == 0 {
		// Create genesis block
		b = CreateBlock(0, nil)
	} else {
		last = bc.blocks[bc.last_index]
		b = CreateBlock(bc.last_index+1, last.hash)
	}

	// Add a money creation output
	txn := CreateTransaction()
	scp := BuildP2PKScript(PublicKeyToBytes(key))
	txOutput := CreateTxOutput(scp, 100)
	txn.AddOutput(txOutput)
	txn.ComputeHash(true)

	b.AddTransaction(txn)

	// Add Txn from queue
	for _, txn = range bc.txnQueue {
		b.AddTransaction(txn)
	}

	bc.txnQueue = []*Transaction{}

	bc.blocks = append(bc.blocks, b)
	bc.last_index = b.index

	return nil
}

func (bc *Blockchain) Dump() {
	for i := 0; i < len(bc.blocks); i++ {
		fmt.Printf("### Block %d ###\n", i)
		fmt.Printf(bc.blocks[i].Dump())
		fmt.Println()
	}

	fmt.Printf("%d block(s).\n", len(bc.blocks))
}

func (bc *Blockchain) CreateTransfertTransaction(wallet Wallet, txnOrder *TxnOrder) (*Transaction, error) {
	required_amount := txnOrder.Amount
	used_funds := make([]*OutputFund, 0)
	funds := bc.GetFunds(&wallet)

	for _, fund := range funds {
		required_amount -= fund.txn.outputs[fund.output_id].amount
		used_funds = append(used_funds, fund)

		if required_amount <= 0 {
			break
		}
	}

	if required_amount > 0 {
		// could not create transaction.
		return nil, errors.New("Not enough funds.")
	}

	txn := new(Transaction)

	for _, used_fund := range used_funds {
		input := new(TxInput)
		input.script = used_fund.script
		input.txhash = used_fund.txn.hash

		txn.AddInput(input)
	}

	output := new(TxOutput)
	output.amount = txnOrder.Amount
	output.script = BuildP2PKHScript([]byte(txnOrder.Addr))
	txn.AddOutput(output)

	// Add remaining funds into a new output
	if required_amount < 0 {
		output = new(TxOutput)
		output.amount = 0 - required_amount
		output.script = BuildP2PKHScript([]byte(GetPublicKeyHash(wallet.PrivateKeys[0].PublicKey)))
		txn.AddOutput(output)
	}

	for _, used_fund := range used_funds {
		// Copy other outputs
		for i, output := range used_fund.txn.outputs {
			if used_fund.output_id != i {
				txn.AddOutput(output)
			}
		}
	}

	txn.ComputeHash(true)

	return txn, nil
}

func (bc *Blockchain) QueueTransaction(txn *Transaction) {
	bc.txnQueue = append(bc.txnQueue, txn)
}

func TryOutput(wallet *Wallet, outputScript *Script) (*Script, bool) {
	// Prepare a VM to execute output scripts.
	vm := new(VM)

	for _, pk := range wallet.PrivateKeys {
		input_scr := new(Script)

		sign, _ := SignMessage(pk, outputScript.data)

		input_scr.addPushBytes(sign)

		res, _ := vm.runInputOutput(*input_scr, *outputScript)

		if res {
			return input_scr, true
		}

		// Try P2PKH
		input_scr.addPushBytes(PublicKeyToBytes(pk.PublicKey))

		res, _ = vm.runInputOutput(*input_scr, *outputScript)

		if res {
			return input_scr, true
		}
	}

	return nil, false
}

// Scan blockchain for unspent output transactions matching out wallet private keys
func (bc *Blockchain) GetFunds(wallet *Wallet) []*OutputFund {
	used_inputs := make(map[string]bool)
	funds := make([]*OutputFund, 0)

	// Scan all blocks for unspent outputs matching our private keys.
	for j := len(bc.blocks) - 1; j >= 0; j-- {
		for _, tx := range bc.blocks[j].txns {
			for _, input := range tx.inputs {
				used_inputs[string(input.txhash)] = true
			}
			if _, ok := used_inputs[string(tx.hash)]; ok {
				// This transaction was already used. Skipping.
				continue
			}
			for k, output := range tx.outputs {
				// Try output.
				script, res := TryOutput(wallet, output.script)
				if res {
					// Adding this transaction in list
					of := new(OutputFund)
					of.output_id = k
					of.txn = tx
					of.script = script

					funds = append(funds, of)

				}
			}
		}
	}

	return funds
}

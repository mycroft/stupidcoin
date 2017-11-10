package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"
)

type TxnOrder struct {
	Addr   string
	Amount float64
}

type Blockchain struct {
	last_index uint64
	blocks     []*Block

	// Do not store in blockchain
	txnQueue []*Transaction
}

func LoadBlockchain(config Config) (*Blockchain, error) {
	blockchain := new(Blockchain)

	if _, err := os.Stat(config.Blockchain); os.IsNotExist(err) {
		fmt.Printf("No existing block chain found...\n")
		blockchain.last_index = 0

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

func (bc *Blockchain) CreateTransaction(wallet Wallet, txnOrder *TxnOrder) {
	// Go through block chain to find output transaction with available coins.
	vm := new(VM)

	pk := wallet.PrivateKeys[0]

	// for each block...
	for i := len(bc.blocks) - 1; i > 0; i-- {
		current_block := bc.blocks[i]
		// for each transaction...
		for j := 0; j < len(current_block.txns); j++ {
			current_txn := current_block.txns[j]
			// for each output...
			for k := 0; k < len(current_txn.outputs); k++ {
				input_scr := new(Script)

				sign, err := SignMessage(pk, current_txn.outputs[k].script.data)
				if err != nil {
					fmt.Println(err)
					return
				}

				input_scr.addPushBytes(sign)

				res, err := vm.runInputOutput(input_scr, current_txn.outputs[k].script)
				if err != nil {
					fmt.Println(err)
					return
				}

				if res && current_txn.outputs[k].amount >= txnOrder.Amount {
					// We found a transaction that suits !

					txn := new(Transaction)
					input := new(TxInput)
					input.script = input_scr
					input.txhash = current_txn.hash

					txn.AddInput(input)

					output := new(TxOutput)
					output.amount = txnOrder.Amount
					output.script = BuildP2PKHScript([]byte(txnOrder.Addr))

					txn.AddOutput(output)

					if current_txn.outputs[k].amount-txnOrder.Amount > 0 {
						output = new(TxOutput)
						output.amount = current_txn.outputs[k].amount - txnOrder.Amount
						output.script = BuildP2PKScript([]byte(GetPublicKeyHash(pk.PublicKey)))
						txn.AddOutput(output)
					}

					txn.ComputeHash(true)

					bc.txnQueue = append(bc.txnQueue, txn)

					return
				}
			}
		}
	}

	return
}

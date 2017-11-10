package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"
)

type Blockchain struct {
	last_index uint64
	blocks     []*Block
}

func LoadBlockchain() (*Blockchain, error) {
	blockchain := new(Blockchain)

	if _, err := os.Stat(".blocks.dat"); os.IsNotExist(err) {
		fmt.Printf("No existing block chain found...\n")
		blockchain.last_index = 0

		return blockchain, nil
	}

	fd, err := os.Open(".blocks.dat")
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

func (bc *Blockchain) SaveBlockchain() error {
	fd, err := os.OpenFile(".blocks.dat", os.O_RDWR|os.O_CREATE, 0644)
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
	b.AddTransaction(txn)

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

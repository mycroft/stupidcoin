package main

import (
	"crypto/sha256"

	"fmt"
	"os"
	"strconv"
	"time"
)

type Block struct {
	index     uint64
	last_hash []byte
	timestamp uint64
	hash      []byte
	txns      []*Transaction
}

func CreateBlock(index uint64, last_hash []byte) *Block {
	block := new(Block)

	block.index = index
	block.last_hash = last_hash
	block.timestamp = uint64(time.Now().Unix())

	block.ComputeHash()

	return block
}

func (b *Block) ComputeHash() {
	h := sha256.New()

	h.Write([]byte(strconv.FormatUint(b.index, 10)))
	h.Write([]byte(strconv.FormatUint(b.timestamp, 10)))

	b.hash = h.Sum(nil)

	return
}

func (b *Block) SaveBlock(fd *os.File) error {
	WriteUint64ToFd(fd, b.index)
	WriteBytesToFd(fd, b.last_hash)
	WriteUint64ToFd(fd, b.timestamp)
	WriteBytesToFd(fd, b.hash)

	return nil
}

func (b *Block) Dump() string {
	var dump string

	dump = fmt.Sprintf("Hash: %x\n", b.hash)
	dump += fmt.Sprintf("LastHash: %x\n", b.last_hash)

	for i := 0; i < len(b.txns); i++ {
		dump += fmt.Sprintf("Txn: %x\n", b.txns[i].hash)

		for j := 0; j < len(b.txns[i].inputs); j++ {
			dump += fmt.Sprintf("Input: %v\n", b.txns[i].inputs[j].script)
		}

		for j := 0; j < len(b.txns[i].outputs); j++ {
			dump += fmt.Sprintf("Output: %v\n", b.txns[i].outputs[j].script)
		}

	}

	return dump
}

func (b *Block) AddTransaction(txn *Transaction) {
	b.txns = append(b.txns, txn)
}

func CreateBlockFromFd(fd *os.File) (*Block, error) {
	var err error
	b := new(Block)

	b.index, err = ReadUint64FromFd(fd)
	if err != nil {
		return nil, err
	}

	b.last_hash, err = ReadBytesFromFd(fd)
	if err != nil {
		return nil, err
	}

	b.timestamp, err = ReadUint64FromFd(fd)
	if err != nil {
		return nil, err
	}

	b.hash, err = ReadBytesFromFd(fd)
	if err != nil {
		return nil, err
	}

	return b, nil
}

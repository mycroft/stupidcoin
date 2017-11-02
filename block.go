package main

import (
	"crypto/sha256"
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

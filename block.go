package main

import (
	"crypto/sha256"

	"bytes"
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

	block.ComputeHash(true)

	return block
}

func (b *Block) ComputeHash(update bool) []byte {
	h := sha256.New()

	h.Write([]byte(strconv.FormatUint(b.index, 10)))
	h.Write(b.last_hash)
	h.Write([]byte(strconv.FormatUint(b.timestamp, 10)))

	for _, txn := range b.txns {
		h.Write(txn.hash)
	}

	hash := h.Sum(nil)

	if update {
		b.hash = hash
	}

	return hash
}

// XXX to rewrite using bytes...
func (b *Block) SaveBlock(fd *os.File) error {
	WriteUint64ToFd(fd, b.index)
	WriteBytesToFd(fd, b.last_hash)
	WriteUint64ToFd(fd, b.timestamp)
	WriteBytesToFd(fd, b.hash)

	// Save transactions
	WriteUint32ToFd(fd, uint32(len(b.txns)))

	for _, txn := range b.txns {
		txn.SaveTransaction(fd)
	}

	return nil
}

func (b *Block) Dump() string {
	var dump string

	dump = fmt.Sprintf("Hash:\t\t%x\n", b.hash)
	dump += fmt.Sprintf("LastHash:\t%x\n", b.last_hash)
	dump += fmt.Sprintf("Txn count:\t%d\n", len(b.txns))

	for i := 0; i < len(b.txns); i++ {
		dump += b.txns[i].String()
	}

	return dump
}

func (b *Block) AddTransaction(txn *Transaction) {
	b.txns = append(b.txns, txn)

	// Recompute hash
	b.ComputeHash(true)
}

func CreateBlockFromFd(fd *os.File) (*Block, error) {
	var err error
	var i uint32
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

	txn_cnt, err := ReadUint32FromFd(fd)
	if err != nil {
		return nil, err
	}

	for i = 0; i < txn_cnt; i++ {
		txn, err := CreateTransactionFromFd(fd)
		if err != nil {
			return nil, err
		}

		b.AddTransaction(txn)
	}

	return b, nil
}

/* Verifying a block.
 * - Verify that hash is correct
 * - Verify that inputs are valid
 *   - First transaction can have a null input as this is block generation
 *
 */
func (b *Block) VerifyBlock() bool {
	hash := b.ComputeHash(false)
	if 0 != bytes.Compare(hash, b.hash) {
		return false
	}

	// XXX to do

	return true
}

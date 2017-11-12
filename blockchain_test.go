package main

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func CreateTestingWallet() *Wallet {
	key, err := CreateKeyPair()
	if err != nil {
		panic(fmt.Sprintf("Could not create key: %s", err))
	}

	wallet := new(Wallet)
	wallet.AddPrivateKey(*key)

	return wallet
}

func TestCreateBlockchain(t *testing.T) {
	bc := CreateBlockchain()
	if bc.last_index != 0 {
		t.Error("Blockchain last index is not 0")
	}
}

func TestMining(t *testing.T) {
	bc := CreateBlockchain()
	wallet := CreateTestingWallet()

	// Mine 1st block
	bc.MineBlock(wallet.PrivateKeys[0].PublicKey)

	if bc.last_index != 0 {
		t.Error("Blockchain last index is not 0")
	}

	if len(bc.blocks) != 1 {
		t.Error("Invalid block count")
	}

	if len(bc.blocks[0].txns) != 1 {
		t.Error("Invalid txn count")
	}

	// Mine 2nd block
	bc.MineBlock(wallet.PrivateKeys[0].PublicKey)

	if bc.last_index != 1 {
		t.Error("Blockchain last index is not 1")
	}

	if len(bc.blocks) != 2 {
		t.Error("Invalid block count")
	}

	if len(bc.blocks[1].txns) != 1 {
		t.Error("Invalid txn count")
	}

	// Scan chain
	funds := bc.GetFunds(wallet)

	if len(funds) != 2 {
		t.Error("Invalid available funds number")
	}

	amount := float64(0)
	for _, fund := range funds {
		amount += fund.txn.outputs[fund.output_id].amount
	}

	if amount != 200 {
		t.Errorf(fmt.Sprintf("Invalid fund value (%f != %f)", amount, float64(200)))
	}
}

func CheckFunds(bc *Blockchain, wallet *Wallet) float64 {
	funds := bc.GetFunds(wallet)

	amount := float64(0)
	for _, fund := range funds {
		amount += fund.txn.outputs[fund.output_id].amount
	}

	return amount
}

func ControlFunds(t *testing.T, wallet *Wallet, bc *Blockchain, amount float64) {
	fund := CheckFunds(bc, wallet)
	fundControl := amount

	tolerance := 0.000001

	if diff := math.Abs(fundControl - fund); diff > tolerance {
		t.Error(fmt.Sprintf("Invalid fund value (%f != %f)", fundControl, fund))
	}
}

func TestTransaction(t *testing.T) {
	// Create 2 wallet
	wallet1 := CreateTestingWallet()
	wallet2 := CreateTestingWallet()

	// Create a blockchain
	bc := CreateBlockchain()

	// Mine 2 initial blocks with 1st wallet

	// Mine 1st block
	bc.MineBlock(wallet1.PrivateKeys[0].PublicKey)
	bc.MineBlock(wallet1.PrivateKeys[0].PublicKey)

	// Check funds for wallets
	ControlFunds(t, wallet1, bc, 200)
	ControlFunds(t, wallet2, bc, 0)

	// Send 150 units from wallet1 to wallet2
	txnOrder := new(TxnOrder)
	txnOrder.Amount = 150.55
	txnOrder.Addr = GetPublicKeyHash(wallet2.PrivateKeys[0].PublicKey)

	txn, err := bc.CreateTransfertTransaction(*wallet1, txnOrder)
	if err != nil {
		t.Error(err)
	}
	bc.QueueTransaction(txn)

	// Not mined yet.
	// Check funds for wallets
	ControlFunds(t, wallet1, bc, 200)
	ControlFunds(t, wallet2, bc, 0)

	// Mine block
	bc.MineBlock(wallet1.PrivateKeys[0].PublicKey)

	// Mined: Control
	ControlFunds(t, wallet1, bc, 300-txnOrder.Amount)
	ControlFunds(t, wallet2, bc, 0+txnOrder.Amount)

	// Send some money back to its owner

	txnOrder.Amount = 130
	txnOrder.Addr = GetPublicKeyHash(wallet1.PrivateKeys[0].PublicKey)

	txn, err = bc.CreateTransfertTransaction(*wallet2, txnOrder)
	if err != nil {
		t.Error(err)
	}
	bc.QueueTransaction(txn)

	// Mine block
	bc.MineBlock(wallet1.PrivateKeys[0].PublicKey)

	// Mined: Control
	ControlFunds(t, wallet1, bc, 400-150.55+txnOrder.Amount)
	ControlFunds(t, wallet2, bc, 150.55-txnOrder.Amount)

	txnOrder.Amount = 20.55
	txnOrder.Addr = GetPublicKeyHash(wallet1.PrivateKeys[0].PublicKey)

	txn, err = bc.CreateTransfertTransaction(*wallet2, txnOrder)
	if err != nil {
		t.Error(err)
	}
	bc.QueueTransaction(txn)

	// Mine block
	bc.MineBlock(wallet1.PrivateKeys[0].PublicKey)

	// Mined: Control
	ControlFunds(t, wallet1, bc, 500)
	ControlFunds(t, wallet2, bc, 0)
}

func TransferFund(bc *Blockchain, w1 *Wallet, w2 *Wallet, amount float64) {
	txnOrder := new(TxnOrder)
	txnOrder.Amount = amount
	txnOrder.Addr = GetPublicKeyHash(w2.PrivateKeys[0].PublicKey)

	txn, err := bc.CreateTransfertTransaction(*w1, txnOrder)
	if err != nil {
		panic(err)
	}
	bc.QueueTransaction(txn)
}

func TestSaveLoadChain(t *testing.T) {
	bc := CreateBlockchain()
	w1 := CreateTestingWallet()
	w2 := CreateTestingWallet()

	// Mine 1st block
	bc.MineBlock(w1.PrivateKeys[0].PublicKey)
	TransferFund(bc, w1, w2, 50)
	bc.MineBlock(w1.PrivateKeys[0].PublicKey)

	f1 := CheckFunds(bc, w1)
	f2 := CheckFunds(bc, w2)

	if f1 != 150 {
		t.Errorf("Invalid sum")
	}

	if f2 != 50 {
		t.Errorf("Invalid sum")
	}

	c := Config{Blockchain: "chain.dat"}
	// Save blockchain.

	err := bc.SaveBlockchain(c)
	if err != nil {
		t.Error(err)
	}

	// Load blockchain.
	c2, err := LoadBlockchain(c)
	if err != nil {
		t.Error(err)
	}

	// Check c2
	if len(bc.blocks) != len(c2.blocks) {
		t.Error("Invalid block number")
	}

	for i := range bc.blocks {
		if bytes.Compare(bc.blocks[i].hash, c2.blocks[i].hash) != 0 {
			t.Errorf("Invalid hash %x / %x", bc.blocks[i].hash, c2.blocks[i].hash)
		}

		if len(bc.blocks[i].txns) != len(c2.blocks[i].txns) {
			t.Error("Invalid txn number")
		}

		for j := range bc.blocks[i].txns {
			txn1 := bc.blocks[i].txns[j]
			txn2 := c2.blocks[i].txns[j]

			if txn1.timestamp != txn2.timestamp {
				t.Errorf("Invalid txn timestamp (%d/%d)", txn1.timestamp, txn2.timestamp)
			}

			if bytes.Compare(bc.blocks[i].txns[j].hash, c2.blocks[i].txns[j].hash) != 0 {
				t.Error("Invalid txn hash")
			}
		}
	}
}

func TestVerifyblock(t *testing.T) {
	bc := CreateBlockchain()
	w1 := CreateTestingWallet()
	w2 := CreateTestingWallet()

	// Mine 1st block
	bc.MineBlock(w1.PrivateKeys[0].PublicKey)
	TransferFund(bc, w1, w2, 50)
	bc.MineBlock(w1.PrivateKeys[0].PublicKey)

	// f1 := CheckFunds(bc, w1)
	// f2 := CheckFunds(bc, w2)

	// Controling blocks

	for _, b := range bc.blocks {
		h := b.ComputeHash(false)

		if bytes.Compare(h, b.hash) != 0 {
			t.Error("Invalid hash")
		}

		for _, tx := range b.txns {
			h := tx.ComputeHash(false)

			if bytes.Compare(h, tx.hash) != 0 {
				t.Errorf("Invalid hash")
			}
		}

		v := b.VerifyBlock()
		if v != true {
			t.Errorf("Invalid block")
		}
	}
}

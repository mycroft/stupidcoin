package main

import (
	"testing"
)

func TestScript1(t *testing.T) {
	scp := new(Script)
	scpOutput := new(Script)

	scp.addInstruction(OP_NOP)

	str := []byte("Hello World")
	scp.addInstruction(OP_PUSH_BYTES)
	scp.addWord(uint16(len(str)))
	scp.addBytes(str)

	hash := []byte("b10a8db164e0754105b7a99be72e3fe5")
	scp.addInstruction(OP_PUSH_BYTES)
	scp.addWord(uint16(len(hash)))
	scp.addBytes(hash)

	scp.addInstruction(OP_SWAP)
	scp.addInstruction(OP_HASH_MD5)
	scp.addInstruction(OP_HASH_TOHEX)

	scp.addInstruction(OP_EQUAL)

	vm := new(VM)
	res, err := vm.runInputOutput(scp, scpOutput)
	if err != nil {
		t.Error(err)
	}

	if res != true {
		t.Error("Result is not true.")
	}

	return
}

func TestScript2(t *testing.T) {
	scp := new(Script)
	scpOutput := new(Script)

	// Input
	str := []byte("Hello World")
	scp.addInstruction(OP_PUSH_BYTES)
	scp.addWord(uint16(len(str)))
	scp.addBytes(str)

	// Script
	scp.addInstruction(OP_HASH_MD5)
	scp.addInstruction(OP_HASH_TOHEX)

	hash := []byte("b10a8db164e0754105b7a99be72e3fe5")
	scp.addInstruction(OP_PUSH_BYTES)
	scp.addWord(uint16(len(hash)))
	scp.addBytes(hash)

	scp.addInstruction(OP_EQUAL)

	vm := new(VM)
	res, err := vm.runInputOutput(scp, scpOutput)
	if err != nil {
		t.Error(err)
	}

	if res != true {
		t.Error("Result is not true.")
	}

	return
}

func TestScript3(t *testing.T) {
	// Create some key pair
	key, err := CreateKeyPair()
	if err != nil {
		t.Errorf("Could not create key...")
	}

	// Create an output script.
	output := BuildP2PKScript(GetPublicKeyHashBytes(key.PublicKey))

	// Create input script (signature)

	input := new(Script)

	sign, err := SignMessage(*key, output.data)
	if err != nil {
		t.Errorf("Could not sign output.")
	}

	input.addInstruction(OP_PUSH_BYTES)
	input.addWord(uint16(len(sign)))
	input.addBytes(sign)

	// Run vm over this input & output

	vm := new(VM)
	res, err := vm.runInputOutput(input, output)
	if err != nil {
		t.Error(err)
	}

	if res != true {
		t.Error("Result is not true.")
	}

	return
}

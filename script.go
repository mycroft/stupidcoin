package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Instruction byte

const (
	OP_NOP         Instruction = iota
	OP_PUSH_BYTE               = 0x10
	OP_PUSH_WORD               = 0x11
	OP_PUSH_DWORD              = 0x12
	OP_PUSH_BYTES              = 0x13
	OP_DUP                     = 0x14
	OP_SWAP                    = 0x15
	OP_EQUAL                   = 0x20
	OP_HASH_BASE58             = 0x30
	OP_HASH_BASE64             = 0x31
	OP_HASH_TOHEX              = 0x32
	OP_HASH_MD5                = 0x40
	OP_HASH_KEY                = 0x41
	OP_CHECKSIG                = 0x50
)

type Script struct {
	data []byte
}

func (script *Script) addInstruction(inst Instruction) {
	script.data = append(script.data, byte(inst))
}

func (script *Script) addByte(value byte) {
	script.data = append(script.data, value)
}

func (script *Script) addWord(value uint16) {
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, value)
	script.data = append(script.data, buffer...)
}

func (script *Script) addDword(value uint32) {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, value)
	script.data = append(script.data, buffer...)
}

func (script *Script) addBytes(bytes []byte) {
	script.data = append(script.data, bytes...)
}

func (script *Script) addPushBytes(bytes []byte) {
	script.addInstruction(OP_PUSH_BYTES)
	script.addWord(uint16(len(bytes)))
	script.data = append(script.data, bytes...)
}

func (script *Script) String() string {
	elem := make([]string, 0)

	for i := 0; i < len(script.data); {
		inst := script.data[i]
		i++

		switch Instruction(inst) {
		case OP_NOP:
			elem = append(elem, "OP_NOP")
		case OP_PUSH_BYTE:
			elem = append(elem, "OP_PUSH_BYTE")
			// xxx add byte
		case OP_PUSH_WORD:
			elem = append(elem, "OP_PUSH_WORD")
			// xxx
		case OP_PUSH_DWORD:
			elem = append(elem, "OP_PUSH_DWORD")
			// xxx
		case OP_PUSH_BYTES:
			elem = append(elem, "OP_PUSH_BYTES")
			size_bytes := int(binary.BigEndian.Uint16(script.data[i : i+2]))
			i += 2
			bytes := script.data[i : i+size_bytes]
			elem = append(elem, fmt.Sprintf("0x%x", bytes))
			i += size_bytes
		case OP_DUP:
			elem = append(elem, "OP_DUP")
		case OP_SWAP:
			elem = append(elem, "OP_SWAP")
		case OP_HASH_KEY:
			elem = append(elem, "OP_HASH_KEY")
		case OP_CHECKSIG:
			elem = append(elem, "OP_CHECKSIG")
		default:
			elem = append(elem, fmt.Sprintf("UNKNOWN:0x%x", inst))
		}
	}

	return strings.Join(elem, " ")
}

func (script *Script) dump() {
	fmt.Println(script.data)
}

func BuildP2PKScript(key []byte) *Script {
	s := new(Script)

	s.addPushBytes(key)
	s.addInstruction(OP_CHECKSIG)

	return s
}

func BuildP2PKHScript(hash []byte) *Script {
	s := new(Script)

	// <Sig> <PubKey> OP_DUP OP_HASH160 <PubkeyHash> OP_EQUALVERIFY OP_CHECKSIG

	s.addInstruction(OP_DUP)
	s.addInstruction(OP_HASH_KEY)

	s.addPushBytes(hash)

	s.addInstruction(OP_EQUAL)
	s.addInstruction(OP_CHECKSIG)

	return s
}

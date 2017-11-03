package main

import (
	"encoding/binary"
	"fmt"
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

func (script *Script) dump() {
	fmt.Println(script.data)
}

func BuildP2PKScript(key []byte) *Script {
	s := new(Script)

	s.addInstruction(OP_PUSH_BYTES)
	s.addWord(uint16(len(key)))
	s.addBytes(key)
	s.addInstruction(OP_CHECKSIG)

	return s
}

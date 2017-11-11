package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
)

type Stack struct {
	data [][]byte
}

type VM struct {
	script      *Script
	stack       *Stack
	current_idx int
}

func NewStack() *Stack {
	stack := new(Stack)

	return stack
}

func (stack *Stack) Push(data []byte) {
	stack.data = append(stack.data, data)
}

func (stack *Stack) Pop() ([]byte, error) {
	var elem []byte

	if len(stack.data) == 0 {
		return nil, errors.New("No element in stack")
	}

	elem = stack.data[len(stack.data)-1]
	stack.data = stack.data[0 : len(stack.data)-1]

	return elem, nil
}

func (stack *Stack) Empty() bool {
	return len(stack.data) == 0
}

func (vm *VM) hasEnough(bytes int) bool {
	return len(vm.script.data) >= (vm.current_idx + bytes)
}

func (vm *VM) runInputOutput(input Script, output Script) (bool, error) {
	vm.stack = NewStack()
	vm.script = &input
	vm.script.data = append(vm.script.data, output.data...)

	for vm.current_idx = 0; vm.current_idx < len(vm.script.data); {
		// pick instruction
		inst := Instruction(vm.script.data[vm.current_idx])
		vm.current_idx++

		switch inst {
		case OP_NOP:
			// Do nothing
			continue
		case OP_PUSH_BYTE:
			// Push one byte to stack
			if !vm.hasEnough(1) {
				return false, errors.New("Not enough bytes in script")
			}
			vm.stack.Push(vm.script.data[vm.current_idx : vm.current_idx+1])
			vm.current_idx++

		case OP_PUSH_WORD:
			// Push 2 bytes to stack
			if !vm.hasEnough(2) {
				return false, errors.New("Not enough bytes in script")
			}
			vm.stack.Push(vm.script.data[vm.current_idx : vm.current_idx+2])
			vm.current_idx += 2

		case OP_PUSH_DWORD:
			// Push 4 bytes to stack
			if !vm.hasEnough(4) {
				return false, errors.New("Not enough bytes in script")
			}
			vm.stack.Push(vm.script.data[vm.current_idx : vm.current_idx+4])
			vm.current_idx += 4

		case OP_PUSH_BYTES:
			// Get number of bytes to push
			if !vm.hasEnough(2) {
				return false, errors.New("Not enough bytes in script")
			}
			size_bytes := vm.script.data[vm.current_idx : vm.current_idx+2]
			vm.current_idx += 2
			size := binary.BigEndian.Uint16(size_bytes)

			if !vm.hasEnough(int(size)) {
				return false, errors.New("Not enough bytes in script")
			}

			vm.stack.Push(vm.script.data[vm.current_idx : vm.current_idx+int(size)])
			vm.current_idx += int(size)
		case OP_DUP:
			elem1, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			vm.stack.Push(elem1)
			vm.stack.Push(elem1)
		case OP_SWAP:
			elem1, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}
			elem2, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			vm.stack.Push(elem1)
			vm.stack.Push(elem2)

		case OP_EQUAL:
			// Picks 2 elements in stack, if not equal, fail.
			elem1, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}
			elem2, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			if len(elem1) != len(elem2) {
				return false, errors.New("OP_EQUAL: Can't compare elements: Invalid sizes")
			}

			for j := 0; j < len(elem1); j++ {
				if elem1[j] != elem2[j] {
					return false, errors.New(fmt.Sprintf("OP_EQUAL: Elements are not equal (%s / %s)", string(elem1), string(elem2)))
				}
			}

		case OP_HASH_BASE64:
			elem1, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			str := base64.StdEncoding.EncodeToString(elem1)
			vm.stack.Push([]byte(str))

		case OP_HASH_TOHEX:
			elem1, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			str := fmt.Sprintf("%x", elem1)
			vm.stack.Push([]byte(str))

		case OP_HASH_MD5:
			elem1, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			hash := md5.Sum(elem1)
			vm.stack.Push(hash[:])

		case OP_HASH_KEY:
			elem1, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			// Recreate key
			pk := GetPublicKeyFromBytes(elem1)

			// Get hash
			hash := GetPublicKeyHash(pk)

			vm.stack.Push([]byte(hash))

		case OP_CHECKSIG:
			// Pop public key
			key, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}
			// Pop signature
			sign, err := vm.stack.Pop()
			if err != nil {
				return false, errors.New("Not enough elements in stack")
			}

			// Rebuild key
			pbkey := GetPublicKeyFromBytes(key)

			// Check signature over script
			ret := SignVerify(pbkey, output.data, sign)
			if ret != true {
				return false, errors.New("Invalid signature")
			}

		default:
			return false, errors.New(fmt.Sprintf("Invalid instruction: 0x%x", inst))
		}
	}

	if !vm.stack.Empty() {
		return false, errors.New(fmt.Sprintf("Remaining elements in stack."))
	}

	return true, nil
}

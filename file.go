package main

import (
	"encoding/binary"
	"math"
	"os"
)

func WriteUint32ToFd(fd *os.File, i uint32) error {
	buffer := make([]byte, 4)

	// Write last_index
	binary.LittleEndian.PutUint32(buffer, i)
	_, err := fd.Write(buffer)

	return err
}

func WriteUint64ToFd(fd *os.File, i uint64) error {
	buffer := make([]byte, 8)

	// Write last_index
	binary.LittleEndian.PutUint64(buffer, i)
	_, err := fd.Write(buffer)

	return err
}

func WriteBytesToFd(fd *os.File, bytes []byte) error {
	// Write size of bytes
	WriteUint32ToFd(fd, uint32(len(bytes)))

	// Writes bytes
	_, err := fd.Write(bytes)

	return err
}

func WriteFloat64ToFd(fd *os.File, f float64) error {
	bits := math.Float64bits(f)
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, bits)

	_, err := fd.Write(buffer)

	return err
}

func ReadUint64FromFd(fd *os.File) (uint64, error) {
	buffer := make([]byte, 8)

	_, err := fd.Read(buffer)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(buffer), nil
}

func ReadUint32FromFd(fd *os.File) (uint32, error) {
	buffer := make([]byte, 4)

	_, err := fd.Read(buffer)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buffer), nil
}

func ReadBytesFromFd(fd *os.File) ([]byte, error) {
	i, err := ReadUint32FromFd(fd)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, i)

	_, err = fd.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func ReadFloat64FromFd(fd *os.File) (float64, error) {
	i, err := ReadUint64FromFd(fd)
	if err != nil {
		return 0, err
	}

	f := math.Float64frombits(i)

	return f, nil
}

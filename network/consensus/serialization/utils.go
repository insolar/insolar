package serialization

import (
	"encoding/binary"
	"io"
)

func setBit(n uint, pos uint) uint {
	return n | (1 << pos)
}

func hasBit(n uint, pos uint) bool {
	return (n & (1 << pos)) > 0
}

func clearBit(n uint, pos uint) uint {
	return n & ^(1 << pos)
}

func toggleBit(n uint, pos uint, val bool) uint {
	if val {
		return setBit(n, pos)
	}

	return clearBit(n, pos)
}

func uintFromBits(bits, start, end uint) uint {
	return (1<<(end-start+1) - 1) & (bits >> start)
}

func read(reader io.Reader, data interface{}) error {
	return binary.Read(reader, defaultByteOrder, data)
}

func write(writer io.Writer, data interface{}) error {
	return binary.Write(writer, defaultByteOrder, data)
}

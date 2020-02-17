// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bits

// ResetBits returns a new byte slice with all bits in 'value' reset,
// starting from 'start' number of bit.
//
// If 'start' is bigger than len(value), the original slice will be returned.
func ResetBits(value []byte, start uint8) []byte {
	if int(start) >= len(value)*8 {
		return value
	}

	startByte := start / 8
	startBit := start % 8

	result := make([]byte, len(value))
	copy(result, value[:startByte])

	// Reset bits in starting byte.
	mask := byte(0xFF)
	mask <<= 8 - startBit
	result[startByte] = value[startByte] & mask

	return result
}

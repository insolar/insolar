//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package longbits

import (
	"math/bits"
	"testing"

	"github.com/stretchr/testify/require"
)

var testSample = []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0xFD, 0xFF, 0xFF}

func copyBits(bytes []byte) []byte {
	return append([]byte(nil), bytes...)
}

func reverseBits(bytes []byte) []byte {
	r := make([]byte, len(bytes))
	for i, v := range bytes {
		r[i] = bits.Reverse8(v)
	}
	return r
}

func Test_BitReader_ReadByte(t *testing.T) {
	bytes := copyBits(testSample)
	br := NewBitArrayReader(FirstLow, bytes)
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		require.Equal(t, b, br.ReadByte())
	}
	require.True(t, br.IsArrayDepleted())

	br = NewBitArrayReader(FirstHigh, reverseBits(bytes))
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		require.Equal(t, b, br.ReadByte())
	}
	require.True(t, br.IsArrayDepleted())
}

func Test_BitStrReader_ReadByte(t *testing.T) {
	bytes := copyBits(testSample)
	br := NewBitStrReader(FirstLow, NewByteString(bytes))
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		require.Equal(t, b, br.ReadByte())
	}
	require.True(t, br.IsArrayDepleted())

	br = NewBitStrReader(FirstHigh, NewByteString(reverseBits(bytes)))
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		require.Equal(t, b, br.ReadByte())
	}
	require.True(t, br.IsArrayDepleted())
}

type testReader interface {
	AlignOffset() uint8
	ReadBit() int
	ReadByte() byte
	ReadSubByte(bitLen uint8) byte
	IsArrayDepleted() bool
}

func testBitReaderRead(t *testing.T, br testReader) {
	require.Equal(t, uint8(0), br.AlignOffset())
	require.Equal(t, 0, br.ReadBit())
	require.Equal(t, uint8(1), br.AlignOffset())
	require.Equal(t, 0, br.ReadBit())
	require.Equal(t, uint8(2), br.AlignOffset())
	require.Equal(t, 1, br.ReadBit())
	require.Equal(t, uint8(3), br.AlignOffset())
	require.Equal(t, 0, br.ReadBit())
	require.Equal(t, uint8(4), br.AlignOffset())
	require.Equal(t, 1, br.ReadBit())
	require.Equal(t, uint8(5), br.AlignOffset())

	require.Equal(t, byte(0xFF), br.ReadByte())
	require.Equal(t, uint8(5), br.AlignOffset())
	require.Equal(t, byte(0xAB), br.ReadByte())
	require.Equal(t, uint8(5), br.AlignOffset())

	require.Equal(t, 0, br.ReadBit())
	require.Equal(t, uint8(6), br.AlignOffset())
	require.Equal(t, 0, br.ReadBit())
	require.Equal(t, uint8(7), br.AlignOffset())
	require.Equal(t, 0, br.ReadBit())
	require.Equal(t, uint8(0), br.AlignOffset())

	require.Equal(t, byte(0x01), br.ReadByte())
	require.Equal(t, byte(0x02), br.ReadByte())

	require.Equal(t, byte(0xFD), br.ReadByte())
	require.Equal(t, byte(0xFF), br.ReadByte())
	require.Equal(t, byte(0xFF), br.ReadByte())

	require.True(t, br.IsArrayDepleted())
	require.Equal(t, uint8(0), br.AlignOffset())
}

func testBitReaderReadSubByte(t *testing.T, br testReader) {
	require.Equal(t, byte(0), br.ReadSubByte(0))
	require.Equal(t, uint8(0), br.AlignOffset())

	require.Equal(t, byte(0), br.ReadSubByte(1))
	require.Equal(t, uint8(1), br.AlignOffset())

	require.Equal(t, byte(2), br.ReadSubByte(2))
	require.Equal(t, uint8(3), br.AlignOffset())
	require.Equal(t, byte(0x7E), br.ReadSubByte(7))
	require.Equal(t, uint8(2), br.AlignOffset())

	require.Equal(t, byte(0x1F), br.ReadSubByte(6))
	require.Equal(t, uint8(0), br.AlignOffset())

	require.Equal(t, byte(0x15), br.ReadSubByte(7))

	require.Equal(t, byte(0x02), br.ReadSubByte(7))
	require.Equal(t, byte(0x08), br.ReadSubByte(7))
	require.Equal(t, byte(0x68), br.ReadSubByte(7))
	require.Equal(t, byte(0x7F), br.ReadSubByte(7))
	require.Equal(t, byte(0x7F), br.ReadSubByte(7))
	require.Equal(t, byte(0x3F), br.ReadSubByte(6))

	require.Equal(t, uint8(0), br.AlignOffset())
	require.True(t, br.IsArrayDepleted())
}

func Test_BitReader_FirstLow_Read(t *testing.T) {
	bytes := copyBits(testSample)
	testBitReaderRead(t, NewBitArrayReader(FirstLow, bytes))
	testBitReaderRead(t, NewBitStrReader(FirstLow, NewByteString(bytes)))
	testBitReaderReadSubByte(t, NewBitArrayReader(FirstLow, bytes))
	testBitReaderReadSubByte(t, NewBitStrReader(FirstLow, NewByteString(bytes)))
}

func testBitReaderReadSubByteCycle(t *testing.T, br testReader) {
	for i := byte(0); i < 8; i++ {
		require.Equal(t, i, br.ReadSubByte(3))
	}
}

func Test_BitReader_FirstHigh_Read(t *testing.T) {
	bytes := reverseBits([]byte{0x88, 0xC6, 0xFA})
	testBitReaderReadSubByteCycle(t, NewBitArrayReader(FirstHigh, bytes))
	testBitReaderReadSubByteCycle(t, NewBitStrReader(FirstHigh, NewByteString(bytes)))
}

func Test_BitReader_FirstLow_SubByte(t *testing.T) {
	bytes := []byte{0x88, 0xC6, 0xFA}
	testBitReaderReadSubByteCycle(t, NewBitArrayReader(FirstLow, bytes))
	testBitReaderReadSubByteCycle(t, NewBitStrReader(FirstLow, NewByteString(bytes)))
}

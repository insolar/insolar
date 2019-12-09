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
	"bytes"
	"io"
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
	br := NewBitArrayReader(LSB, bytes)
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		vb, err := br.ReadByte()
		require.NoError(t, err)
		require.Equal(t, b, vb)
	}
	require.True(t, br.IsArrayDepleted())

	br = NewBitArrayReader(MSB, reverseBits(bytes))
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		vb, err := br.ReadByte()
		require.NoError(t, err)
		require.Equal(t, b, vb)
	}
	require.True(t, br.IsArrayDepleted())
}

func Test_BitStrReader_ReadByte(t *testing.T) {
	bytes := copyBits(testSample)
	br := NewBitStrReader(LSB, CopyBytes(bytes))
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		vb, err := br.ReadByte()
		require.NoError(t, err)
		require.Equal(t, b, vb)
	}
	require.True(t, br.IsArrayDepleted())

	br = NewBitStrReader(MSB, CopyBytes(reverseBits(bytes)))
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		vb, err := br.ReadByte()
		require.NoError(t, err)
		require.Equal(t, b, vb)
	}
	require.True(t, br.IsArrayDepleted())
}

type testReader struct {
	internal interface {
		ReadBit() (int, error)
		ReadByte() (byte, error)
		ReadSubByte(bitLen uint8) (byte, error)
	}
}

func (v testReader) ReadBit() int {
	if r, err := v.internal.ReadBit(); err != nil {
		panic(err)
	} else {
		return r
	}
}

func (v testReader) ReadByte() byte {
	if r, err := v.internal.ReadByte(); err != nil {
		panic(err)
	} else {
		return r
	}
}

func (v testReader) ReadSubByte(bitLen uint8) byte {
	if r, err := v.internal.ReadSubByte(bitLen); err != nil {
		panic(err)
	} else {
		return r
	}
}

type arrayReader interface {
	AlignOffset() uint8
	IsArrayDepleted() bool
}

func (v testReader) IsArray() bool {
	_, ok := v.internal.(arrayReader)
	return ok
}

func (v testReader) testAlignOffset(t *testing.T, ofs uint8) {
	if ar, ok := v.internal.(arrayReader); ok {
		require.Equal(t, ofs, ar.AlignOffset())
	}
}

func (v testReader) testArrayDepleted(t *testing.T, isDepleted bool) {
	if ar, ok := v.internal.(arrayReader); ok {
		if isDepleted {
			require.True(t, ar.IsArrayDepleted())
		} else {
			require.False(t, ar.IsArrayDepleted())
		}
	}
}

func testBitReaderRead(t *testing.T, br testReader) {
	br.testAlignOffset(t, 0)
	require.Equal(t, 0, br.ReadBit())
	br.testAlignOffset(t, 1)
	require.Equal(t, 0, br.ReadBit())
	br.testAlignOffset(t, 2)
	require.Equal(t, 1, br.ReadBit())
	br.testAlignOffset(t, 3)
	require.Equal(t, 0, br.ReadBit())
	br.testAlignOffset(t, 4)
	require.Equal(t, 1, br.ReadBit())
	br.testAlignOffset(t, 5)

	require.Equal(t, byte(0xFF), br.ReadByte())
	br.testAlignOffset(t, 5)
	require.Equal(t, byte(0xAB), br.ReadByte())
	br.testAlignOffset(t, 5)

	require.Equal(t, 0, br.ReadBit())
	br.testAlignOffset(t, 6)
	require.Equal(t, 0, br.ReadBit())
	br.testAlignOffset(t, 7)
	require.Equal(t, 0, br.ReadBit())
	br.testAlignOffset(t, 0)

	require.Equal(t, byte(0x01), br.ReadByte())
	require.Equal(t, byte(0x02), br.ReadByte())

	require.Equal(t, byte(0xFD), br.ReadByte())
	require.Equal(t, byte(0xFF), br.ReadByte())
	require.Equal(t, byte(0xFF), br.ReadByte())

	br.testArrayDepleted(t, true)
	br.testAlignOffset(t, 0)
}

func testBitReaderReadSubByte(t *testing.T, br testReader) {
	require.Equal(t, byte(0), br.ReadSubByte(0))
	br.testAlignOffset(t, 0)

	require.Equal(t, byte(0), br.ReadSubByte(1))
	br.testAlignOffset(t, 1)

	require.Equal(t, byte(2), br.ReadSubByte(2))
	br.testAlignOffset(t, 3)

	require.Equal(t, byte(0x7E), br.ReadSubByte(7))
	br.testAlignOffset(t, 2)

	require.Equal(t, byte(0x1F), br.ReadSubByte(6))
	br.testAlignOffset(t, 0)

	require.Equal(t, byte(0x15), br.ReadSubByte(7))

	require.Equal(t, byte(0x02), br.ReadSubByte(7))
	require.Equal(t, byte(0x08), br.ReadSubByte(7))
	require.Equal(t, byte(0x68), br.ReadSubByte(7))
	require.Equal(t, byte(0x7F), br.ReadSubByte(7))
	require.Equal(t, byte(0x7F), br.ReadSubByte(7))
	require.Equal(t, byte(0x3F), br.ReadSubByte(6))

	br.testAlignOffset(t, 0)
	br.testArrayDepleted(t, true)
}

func newByteReader(b []byte) io.ByteReader {
	return bytes.NewReader(b)
}

func Test_BitReader_FirstLow_Read(t *testing.T) {
	bytes := copyBits(testSample)
	testBitReaderRead(t, testReader{NewBitIoReader(LSB, newByteReader(bytes))})
	testBitReaderRead(t, testReader{NewBitArrayReader(LSB, bytes)})
	testBitReaderRead(t, testReader{NewBitStrReader(LSB, CopyBytes(bytes))})
	testBitReaderReadSubByte(t, testReader{NewBitIoReader(LSB, newByteReader(bytes))})
	testBitReaderReadSubByte(t, testReader{NewBitArrayReader(LSB, bytes)})
	testBitReaderReadSubByte(t, testReader{NewBitStrReader(LSB, CopyBytes(bytes))})
}

func testBitReaderReadSubByteCycle(t *testing.T, br testReader) {
	for i := byte(0); i < 8; i++ {
		require.Equal(t, i, br.ReadSubByte(3))
	}
}

func Test_BitReader_FirstHigh_SubByte(t *testing.T) {
	bytes := reverseBits([]byte{0x88, 0xC6, 0xFA})
	testBitReaderReadSubByteCycle(t, testReader{NewBitIoReader(MSB, newByteReader(bytes))})
	testBitReaderReadSubByteCycle(t, testReader{NewBitArrayReader(MSB, bytes)})
	testBitReaderReadSubByteCycle(t, testReader{NewBitStrReader(MSB, CopyBytes(bytes))})
}

func Test_BitReader_FirstLow_SubByte(t *testing.T) {
	bytes := []byte{0x88, 0xC6, 0xFA}
	testBitReaderReadSubByteCycle(t, testReader{NewBitIoReader(LSB, newByteReader(bytes))})
	testBitReaderReadSubByteCycle(t, testReader{NewBitArrayReader(LSB, bytes)})
	testBitReaderReadSubByteCycle(t, testReader{NewBitStrReader(LSB, CopyBytes(bytes))})
}

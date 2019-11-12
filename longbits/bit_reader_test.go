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
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_BitReader_ReadByte(t *testing.T) {
	bytes := []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0xFF, 0xFF, 0xFF}
	br := NewBitArrayReader(FirstLow, bytes)
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		require.Equal(t, b, br.ReadByte())
	}
	require.True(t, br.IsArrayDepleted())

	br = NewBitArrayReader(FirstHigh, bytes)
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		require.Equal(t, b, br.ReadByte())
	}
	require.True(t, br.IsArrayDepleted())
}

func Test_BitStrReader_ReadByte(t *testing.T) {
	bytes := []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0xFF, 0xFF, 0xFF}
	br := NewBitStrReader(FirstLow, NewByteString(bytes))
	for _, b := range bytes {
		require.False(t, br.IsArrayDepleted())
		require.Equal(t, b, br.ReadByte())
	}
	require.True(t, br.IsArrayDepleted())

	br = NewBitStrReader(FirstHigh, NewByteString(bytes))
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
	IsArrayDepleted() bool
}

func test_BitReader_Read(t *testing.T, br testReader) {
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

	require.Equal(t, byte(0xFF), br.ReadByte())
	require.Equal(t, byte(0xFF), br.ReadByte())
	require.Equal(t, byte(0xFF), br.ReadByte())

	require.True(t, br.IsArrayDepleted())
	require.Equal(t, uint8(0), br.AlignOffset())
}

func Test_BitReader_ReadFirstLow(t *testing.T) {
	bytes := []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0xFF, 0xFF, 0xFF}
	test_BitReader_Read(t, NewBitArrayReader(FirstLow, bytes))
	test_BitReader_Read(t, NewBitStrReader(FirstLow, NewByteString(bytes)))
}

func Test_BitReader_ReadFirstHigh(t *testing.T) {
	bytes := []byte{0x2F, 0xFD, 0x58, 0x01, 0x02, 0xFF, 0xFF, 0xFF}
	test_BitReader_Read(t, NewBitArrayReader(FirstHigh, bytes))
	test_BitReader_Read(t, NewBitStrReader(FirstHigh, NewByteString(bytes)))
}

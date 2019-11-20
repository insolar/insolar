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

func TestBitBuilder_FirstLow_BitOrder(t *testing.T) {
	bb := BitBuilder{}
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(1)
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(1)
	bb.AppendBit(0)
	bb.AppendBit(1)

	bb.AppendByte(0xB6)
	bb.AppendSubByte(0x06, 4)
	bb.AppendSubByte(0x0B, 4)

	require.Equal(t, []byte{0xB6, 0xB6, 0xB6}, bb.dump())
	require.Equal(t, 24, bb.Len())
}

func TestBitBuilder_FirstHigh_BitOrder(t *testing.T) {
	bb := NewBitBuilder(MSB, 0)
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(1)
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(1)
	bb.AppendBit(0)
	bb.AppendBit(1)

	bb.AppendByte(0xB6)
	bb.AppendSubByte(0x06, 4)
	bb.AppendSubByte(0x0B, 4)

	require.Equal(t, []byte{0x6D, 0x6D, 0x6D}, bb.dump())
	require.Equal(t, 24, bb.Len())
}

func TestBitBuilder_FirstLow(t *testing.T) {
	bb := BitBuilder{}
	bb.AppendBit(0)
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(1)

	require.Equal(t, []byte{0x34}, bb.dump())
	require.Equal(t, 6, bb.Len())

	bb.AppendNBit(7, 1)
	require.Equal(t, []byte{0xF4, 0x1F}, bb.dump())
	require.Equal(t, 13, bb.Len())

	bb.AppendByte(0xAB)
	require.Equal(t, []byte{0xF4, 0x7F, 0x15}, bb.dump())
	require.Equal(t, 21, bb.Len())

	bb.PadWithBit(0)
	require.Equal(t, []byte{0xF4, 0x7F, 0x15}, bb.dump())
	require.Equal(t, 24, bb.Len())

	bb.PadWithBit(0)
	require.Equal(t, []byte{0xF4, 0x7F, 0x15}, bb.dump())
	require.Equal(t, 24, bb.Len())

	bb.AppendAlignedByte(0x01)
	bb.AppendByte(0x02)

	require.Equal(t, []byte{0xF4, 0x7F, 0x15, 0x01, 0x02}, bb.dump())
	require.Equal(t, 40, bb.Len())

	bb.AppendSubByte(0x5D, 4)
	require.Equal(t, []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0x0D}, bb.dump())
	require.Equal(t, 44, bb.Len())

	bb.AppendNBit(15, 1)
	require.Equal(t, []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0xFD, 0xFF, 0x07}, bb.dump())
	require.Equal(t, 59, bb.Len())

	bb.PadWithBit(1)
	require.Equal(t, []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0xFD, 0xFF, 0xFF}, bb.dump())
	require.Equal(t, 64, bb.Len())

	bb.PadWithBit(1)
	require.Equal(t, []byte{0xF4, 0x7F, 0x15, 0x01, 0x02, 0xFD, 0xFF, 0xFF}, bb.dump())
	require.Equal(t, 64, bb.Len())
}

func TestBitBuilder_FirstHigh(t *testing.T) {
	bb := NewBitBuilder(MSB, 0)
	bb.AppendBit(0)
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(0)
	bb.AppendBit(1)
	bb.AppendBit(1)

	require.Equal(t, []byte{0x2C}, bb.dump())
	require.Equal(t, 6, bb.Len())

	bb.AppendNBit(7, 1)
	require.Equal(t, []byte{0x2F, 0xF8}, bb.dump())
	require.Equal(t, 13, bb.Len())

	bb.AppendByte(0xAB)
	require.Equal(t, []byte{0x2F, 0xFE, 0xA8}, bb.dump())
	require.Equal(t, 21, bb.Len())

	bb.PadWithBit(0)
	require.Equal(t, []byte{0x2F, 0xFE, 0xA8}, bb.dump())
	require.Equal(t, 24, bb.Len())

	bb.PadWithBit(0)
	require.Equal(t, []byte{0x2F, 0xFE, 0xA8}, bb.dump())
	require.Equal(t, 24, bb.Len())

	bb.AppendAlignedByte(0x01)
	bb.AppendByte(0x02)

	require.Equal(t, []byte{0x2F, 0xFE, 0xA8, 0x80, 0x40}, bb.dump())
	require.Equal(t, 40, bb.Len())

	bb.AppendSubByte(0x5D, 4)
	require.Equal(t, []byte{0x2F, 0xFE, 0xA8, 0x80, 0x40, 0xB0}, bb.dump())
	require.Equal(t, 44, bb.Len())

	bb.AppendNBit(15, 1)
	require.Equal(t, []byte{0x2F, 0xFE, 0xA8, 0x80, 0x40, 0xBF, 0xFF, 0xE0}, bb.dump())
	require.Equal(t, 59, bb.Len())

	bb.PadWithBit(1)
	require.Equal(t, []byte{0x2F, 0xFE, 0xA8, 0x80, 0x40, 0xBF, 0xFF, 0xFF}, bb.dump())
	require.Equal(t, 64, bb.Len())

	bb.PadWithBit(1)
	require.Equal(t, []byte{0x2F, 0xFE, 0xA8, 0x80, 0x40, 0xBF, 0xFF, 0xFF}, bb.dump())
	require.Equal(t, 64, bb.Len())
}

func TestBitBuilder_FirstLow_SubByte(t *testing.T) {
	bb := BitBuilder{}
	for i := byte(0); i < 8; i++ {
		bb.AppendSubByte(i, 3)
	}
	require.Equal(t, 24, bb.Len())
	require.Equal(t, []byte{0x88, 0xC6, 0xFA}, bb.dump())
}

func TestBitBuilder_FirstHigh_SubByte(t *testing.T) {
	bb := NewBitBuilder(MSB, 0)
	for i := byte(0); i < 8; i++ {
		bb.AppendSubByte(i, 3)
	}
	require.Equal(t, 24, bb.Len())
	require.Equal(t, []byte{0x11, 0x63, 0x5F}, bb.dump())
}

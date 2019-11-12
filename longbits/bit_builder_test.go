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
	bb := NewBitBuilder(FirstHigh, 0)
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
	require.Equal(t, []byte{0x2F, 0xFD, 0x58}, bb.dump())
	require.Equal(t, 21, bb.Len())

	bb.PadWithBit(0)
	require.Equal(t, []byte{0x2F, 0xFD, 0x58}, bb.dump())
	require.Equal(t, 24, bb.Len())

	bb.PadWithBit(0)
	require.Equal(t, []byte{0x2F, 0xFD, 0x58}, bb.dump())
	require.Equal(t, 24, bb.Len())

	bb.AppendAlignedByte(0x01)
	bb.AppendByte(0x02)

	require.Equal(t, []byte{0x2F, 0xFD, 0x58, 0x01, 0x02}, bb.dump())
	require.Equal(t, 40, bb.Len())

	bb.AppendSubByte(0x5D, 4)
	require.Equal(t, []byte{0x2F, 0xFD, 0x58, 0x01, 0x02, 0x50}, bb.dump())
	require.Equal(t, 44, bb.Len())

	bb.AppendNBit(15, 1)
	require.Equal(t, []byte{0x2F, 0xFD, 0x58, 0x01, 0x02, 0x5F, 0xFF, 0xE0}, bb.dump())
	require.Equal(t, 59, bb.Len())

	bb.PadWithBit(1)
	require.Equal(t, []byte{0x2F, 0xFD, 0x58, 0x01, 0x02, 0x5F, 0xFF, 0xFF}, bb.dump())
	require.Equal(t, 64, bb.Len())

	bb.PadWithBit(1)
	require.Equal(t, []byte{0x2F, 0xFD, 0x58, 0x01, 0x02, 0x5F, 0xFF, 0xFF}, bb.dump())
	require.Equal(t, 64, bb.Len())
}

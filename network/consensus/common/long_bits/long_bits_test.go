///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
///

package long_bits

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBits64(t *testing.T) {
	bits := NewBits64(0x2211)
	require.Equal(t, uint8(0x11), bits.AsBytes()[0])

	require.Equal(t, uint8(0x22), bits.AsBytes()[1])

	require.Equal(t, uint8(0), bits.AsBytes()[7])
}

func TestBits64WriteTo(t *testing.T) {
	bits := NewBits64(1)
	n, err := bits.WriteTo(&writerToComparer{other: &bits})
	require.Equal(t, nil, err)

	require.Equal(t, int64(8), n)

	require.Equal(t, uint8(1), bits.AsBytes()[0])

	require.Panics(t, func() { bits.WriteTo(&writerToComparer{}) })
}

func TestBits64Read(t *testing.T) {
	bits := NewBits64(1)
	dest := make([]byte, 2)
	n, err := bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 2, n)

	require.Equal(t, uint8(1), dest[0])

	dest = make([]byte, 9)
	n, err = bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 8, n)

	require.Equal(t, uint8(1), dest[0])

	n, err = bits.Read(nil)

	require.Equal(t, nil, err)

	require.Equal(t, 0, n)
}

func TestBits64FoldToUint64(t *testing.T) {
	b := uint64(0x807060504030201)
	bits := NewBits64(b)
	require.Equal(t, b, bits.FoldToUint64())
}

func TestBits64FixedByteSize(t *testing.T) {
	bits := NewBits64(1)
	require.Equal(t, 8, bits.FixedByteSize())
}

func TestBits64AsByteString(t *testing.T) {
	bits := NewBits64(0x4142434445464748)
	require.Equal(t, "HGFEDCBA", bits.AsByteString())
}

func TestBits64String(t *testing.T) {
	require.True(t, NewBits64(1).String() != "")
}

func TestBits64AsBytes(t *testing.T) {
	bits := NewBits64(0x807060504030201)
	require.Equal(t, []uint8{1, 2, 3, 4, 5, 6, 7, 8}, bits.AsBytes())
}

func TestFoldToBits64(t *testing.T) {
	require.Equal(t, NewBits64(0x807060504030201), FoldToBits64([]byte{1, 2, 3, 4, 5, 6, 7, 8}))

	require.Equal(t, NewBits64(0x1808080808080808), FoldToBits64([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}))

	require.Panics(t, func() { FoldToBits64([]byte{1}) })
}

func TestNewBits128(t *testing.T) {
	bits := NewBits128(0x11, 0x22)
	require.Equal(t, uint8(0x11), bits.AsBytes()[0])

	require.Equal(t, uint8(0x22), bits.AsBytes()[8])
}

func TestBits128WriteTo(t *testing.T) {
	bits := NewBits128(1, 2)
	n, err := bits.WriteTo(&writerToComparer{other: &bits})
	require.Equal(t, nil, err)

	require.Equal(t, int64(16), n)

	require.Equal(t, uint8(1), bits.AsBytes()[0])

	require.Equal(t, uint8(2), bits.AsBytes()[8])

	require.Panics(t, func() { bits.WriteTo(&writerToComparer{}) })
}

func TestBits128Read(t *testing.T) {
	bits := NewBits128(1, 2)
	dest := make([]byte, 2)
	n, err := bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 2, n)

	require.Equal(t, uint8(1), dest[0])

	dest = make([]byte, 17)
	n, err = bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 16, n)

	require.Equal(t, uint8(1), dest[0])

	require.Equal(t, uint8(2), dest[8])

	n, err = bits.Read(nil)

	require.Equal(t, nil, err)

	require.Equal(t, 0, n)
}

func TestBits128FoldToUint64(t *testing.T) {
	l := uint64(0x807060504030201)
	h := uint64(0x10F0E0D0C0B0A09)
	bits := NewBits128(l, h)
	require.Equal(t, uint64(0x908080808080808), bits.FoldToUint64())
}

func TestBits128FixedByteSize(t *testing.T) {
	bits := NewBits128(1, 2)
	require.Equal(t, 16, bits.FixedByteSize())
}

func TestBits128String(t *testing.T) {
	require.True(t, NewBits128(1, 2).String() != "")
}

func TestBits128AsByteString(t *testing.T) {
	bits := NewBits128(0x4142434445464748, 0x494A4B4C4D4E4F50)
	require.Equal(t, "HGFEDCBAPONMLKJI", bits.AsByteString())
}

func TestBits128AsBytes(t *testing.T) {
	bits := NewBits128(0x807060504030201, 0x10F0E0D0C0B0A09)
	require.Equal(t, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1}, bits.AsBytes())
}

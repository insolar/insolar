//
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
//

package longbits

import (
	"encoding/binary"
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

	require.Panics(t, func() { _, _ = bits.WriteTo(&writerToComparer{}) })
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

	require.Equal(t, NewBits64(0x1808080808080808), FoldToBits64([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16}))

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

	require.Panics(t, func() { _, _ = bits.WriteTo(&writerToComparer{}) })
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

func TestBits224WriteTo(t *testing.T) {
	bits := Bits224{1}
	n, err := bits.WriteTo(&writerToComparer{other: &bits})
	require.Equal(t, nil, err)

	require.Equal(t, int64(24), n)

	require.Equal(t, uint8(1), bits.AsBytes()[0])

	require.Panics(t, func() { _, _ = bits.WriteTo(&writerToComparer{}) })
}

func TestBits224Read(t *testing.T) {
	bits := Bits224{1, 2, 3}
	dest := make([]byte, 2)
	n, err := bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 2, n)

	require.Equal(t, uint8(1), dest[0])

	dest = make([]byte, 25)
	n, err = bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 24, n)

	require.Equal(t, uint8(1), dest[0])

	n, err = bits.Read(nil)
	require.Equal(t, nil, err)

	require.Equal(t, 0, n)
}

func TestBits224FoldToUint64(t *testing.T) {
	bits := Bits224{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x0908070605040302))
	require.Equal(t, uint64(0xf0e0d0c0b0a), bits.FoldToUint64())
}

func TestBits224FixedByteSize(t *testing.T) {
	bits := Bits224{}
	require.Equal(t, 24, bits.FixedByteSize())
}

func TestBits224String(t *testing.T) {
	require.True(t, Bits224{}.String() != "")
}

func TestBits224AsByteString(t *testing.T) {
	bits := Bits224{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x4142434445464748))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x494A4B4C4D4E4F50))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x5152535455565758))
	require.Equal(t, "HGFEDCBAPONMLKJIXWVUTSRQ", bits.AsByteString())
}

func TestBits224AsBytes(t *testing.T) {
	bits := Bits224{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x908070605040302))
	require.Equal(t, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		1, 2, 3, 4, 5, 6, 7, 8, 9}, bits.AsBytes())
}

func TestBits256WriteTo(t *testing.T) {
	bits := Bits256{1}
	n, err := bits.WriteTo(&writerToComparer{other: &bits})
	require.Equal(t, nil, err)

	require.Equal(t, int64(32), n)

	require.Equal(t, uint8(1), bits.AsBytes()[0])

	require.Panics(t, func() { _, _ = bits.WriteTo(&writerToComparer{}) })
}

func TestBits256Read(t *testing.T) {
	bits := Bits256{1, 2, 3}
	dest := make([]byte, 2)
	n, err := bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 2, n)

	require.Equal(t, uint8(1), dest[0])

	dest = make([]byte, 33)
	n, err = bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 32, n)

	require.Equal(t, uint8(1), dest[0])

	n, err = bits.Read(nil)
	require.Equal(t, nil, err)

	require.Equal(t, 0, n)
}

func TestBits256FoldToUint64(t *testing.T) {
	bits := Bits256{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x0908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	require.Equal(t, uint64(0x201000000000000), bits.FoldToUint64())
}

func TestBits256FoldToBits128(t *testing.T) {
	bits := Bits256{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x0908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	require.Equal(t, Bits128{3, 1, 7, 1, 3, 1, 15, 1, 3, 1, 7, 1, 3, 1, 14, 3}, bits.FoldToBits128())
}

func TestBits256FoldToBits224(t *testing.T) {
	bits := Bits256{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x0908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	require.Equal(t, Bits224{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		1, 2, 3, 4, 5, 6, 7, 8, 9}, bits.FoldToBits224())
}

func TestBits256FixedByteSize(t *testing.T) {
	bits := Bits256{}
	require.Equal(t, 32, bits.FixedByteSize())
}

func TestBits256String(t *testing.T) {
	require.True(t, Bits256{}.String() != "")
}

func TestBits256AsByteString(t *testing.T) {
	bits := Bits256{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x4142434445464748))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x494A4B4C4D4E4F50))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x5152535455565758))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x595A5B5C5D5E5F60))
	require.Equal(t, "HGFEDCBAPONMLKJIXWVUTSRQ`_^]\\[ZY", bits.AsByteString())
}

func TestBits256AsBytes(t *testing.T) {
	bits := Bits256{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	require.Equal(t, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2}, bits.AsBytes())
}

func TestBits512WriteTo(t *testing.T) {
	bits := Bits512{1}
	n, err := bits.WriteTo(&writerToComparer{other: &bits})
	require.Equal(t, nil, err)

	require.Equal(t, int64(64), n)

	require.Equal(t, uint8(1), bits.AsBytes()[0])

	require.Panics(t, func() { _, _ = bits.WriteTo(&writerToComparer{}) })
}

func TestBits512Read(t *testing.T) {
	bits := Bits512{1, 2, 3}
	dest := make([]byte, 2)
	n, err := bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 2, n)

	require.Equal(t, uint8(1), dest[0])

	dest = make([]byte, 65)
	n, err = bits.Read(dest)
	require.Equal(t, nil, err)

	require.Equal(t, 64, n)

	require.Equal(t, uint8(1), dest[0])

	n, err = bits.Read(nil)
	require.Equal(t, nil, err)

	require.Equal(t, 0, n)
}

func TestBits512FoldToUint64(t *testing.T) {
	bits := Bits512{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x0908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	binary.LittleEndian.PutUint64(bits[32:40], uint64(0x0A09080706050403))
	binary.LittleEndian.PutUint64(bits[40:48], uint64(0x0302010F0E0D0C0B))
	binary.LittleEndian.PutUint64(bits[48:56], uint64(0x0B0A090807060504))
	binary.LittleEndian.PutUint64(bits[56:64], uint64(0x040302010F0E0D0C))
	require.Equal(t, uint64(0x403020100000000), bits.FoldToUint64())
}

func TestBits512FoldToBits256(t *testing.T) {
	bits := Bits512{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x0908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	binary.LittleEndian.PutUint64(bits[32:40], uint64(0x0A09080706050403))
	binary.LittleEndian.PutUint64(bits[40:48], uint64(0x0302010F0E0D0C0B))
	binary.LittleEndian.PutUint64(bits[48:56], uint64(0x0B0A090807060504))
	binary.LittleEndian.PutUint64(bits[56:64], uint64(0x040302010F0E0D0C))
	require.Equal(t, Bits256{2, 6, 6, 2, 2, 14, 14, 2, 2, 6, 6, 2, 2, 15, 13, 2, 6, 6, 2, 2, 14, 14,
		2, 2, 6, 6, 2, 2, 15, 13, 2, 6}, bits.FoldToBits256())
}

func TestBits512FoldToBits224(t *testing.T) {
	bits := Bits512{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x0908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	binary.LittleEndian.PutUint64(bits[32:40], uint64(0x0A09080706050403))
	binary.LittleEndian.PutUint64(bits[40:48], uint64(0x0302010F0E0D0C0B))
	binary.LittleEndian.PutUint64(bits[48:56], uint64(0x0B0A090807060504))
	binary.LittleEndian.PutUint64(bits[56:64], uint64(0x040302010F0E0D0C))
	require.Equal(t, Bits224{2, 6, 6, 2, 2, 14, 14, 2, 2, 6, 6, 2, 2, 15, 13, 2, 6, 6, 2, 2, 14, 14,
		2, 2}, bits.FoldToBits224())
}

func TestBits512FixedByteSize(t *testing.T) {
	bits := Bits512{}
	require.Equal(t, 64, bits.FixedByteSize())
}

func TestBits512String(t *testing.T) {
	require.True(t, Bits512{}.String() != "")
}

func TestBits512AsByteString(t *testing.T) {
	bits := Bits512{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x4142434445464748))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x494A4B4C4D4E4F50))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x5152535455565758))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x595A5B5C5D5E5F60))
	binary.LittleEndian.PutUint64(bits[32:40], uint64(0x6162636465666768))
	binary.LittleEndian.PutUint64(bits[40:48], uint64(0x696A6B6C6D6E6F70))
	binary.LittleEndian.PutUint64(bits[48:56], uint64(0x7172737475767778))
	binary.LittleEndian.PutUint64(bits[56:64], uint64(0x797A7B7C7D7E7F80))
	require.Equal(t, "HGFEDCBAPONMLKJIXWVUTSRQ`_^]\\[ZYhgfedcbaponmlkjixwvutsrq\x80\u007f~}|{zy",
		bits.AsByteString())
}

func TestBits512AsBytes(t *testing.T) {
	bits := Bits512{}
	binary.LittleEndian.PutUint64(bits[:8], uint64(0x807060504030201))
	binary.LittleEndian.PutUint64(bits[8:16], uint64(0x10F0E0D0C0B0A09))
	binary.LittleEndian.PutUint64(bits[16:24], uint64(0x908070605040302))
	binary.LittleEndian.PutUint64(bits[24:32], uint64(0x02010F0E0D0C0B0A))
	binary.LittleEndian.PutUint64(bits[32:40], uint64(0x0A09080706050403))
	binary.LittleEndian.PutUint64(bits[40:48], uint64(0x0302010F0E0D0C0B))
	binary.LittleEndian.PutUint64(bits[48:56], uint64(0x0B0A090807060504))
	binary.LittleEndian.PutUint64(bits[56:64], uint64(0x040302010F0E0D0C))
	require.Equal(t, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2, 3, 4}, bits.AsBytes())
}

func TestFoldedFoldToUint64(t *testing.T) {
	require.Equal(t, uint64(0x807060504030201), FoldToUint64([]byte{1, 2, 3, 4, 5, 6, 7, 8}))

	require.Equal(t, uint64(0x1808080808080808), FoldToUint64([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16}))

	require.Panics(t, func() { FoldToUint64([]byte{1}) })
}

func TestFillBitsWithStaticNoise(t *testing.T) {
	bytes := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	FillBitsWithStaticNoise(5, bytes)
	require.Equal(t, []byte{0xde, 0xaa, 0xdb, 0x79, 0x6d, 0xd5, 0xed, 0x3c}, bytes)

	bytes = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	FillBitsWithStaticNoise(31, bytes)
	require.Equal(t, []byte{0xc9, 0xaa, 0xcc, 0xf9, 0x65, 0x55, 0xe6, 0x3c, 0xc8, 0x2a, 0xcd,
		0xb9, 0x64, 0x95, 0xe6, 0x3c}, bytes)

	require.Panics(t, func() { FillBitsWithStaticNoise(1, []byte{1}) })
}

func TestReadFromArray(t *testing.T) {
	t.Skipped()
	// var d, s []byte
	// n, err := readFromArray(d, s)
	// require.Equal(t, 0, n)
	//
	// require.Equal(t, nil, err)
	//
	// d = make([]byte, 1)
	// n, err = readFromArray(d, s)
	// require.Equal(t, 0, n)
	//
	// require.Equal(t, nil, err)
	//
	// d = nil
	// s = []byte{1}
	// n, err = readFromArray(d, s)
	// require.Equal(t, 0, n)
	//
	// require.Equal(t, nil, err)
	//
	// d = make([]byte, 1)
	// n, err = readFromArray(d, s)
	// require.Equal(t, 1, n)
	//
	// require.Equal(t, nil, err)
	//
	// require.Equal(t, s[0], d[0])
	//
	// require.Equal(t, uint8(1), d[0])
	//
	// d = make([]byte, 2)
	// n, err = readFromArray(d, s)
	// require.Equal(t, 1, n)
	//
	// require.Equal(t, nil, err)
	//
	// d = make([]byte, 1)
	// s = make([]byte, 2)
	// n, err = readFromArray(d, s)
	// require.Equal(t, 1, n)
	//
	// require.Equal(t, nil, err)
}

func TestBitsToStringDefault(t *testing.T) {
	bits := Bits64{1}
	require.True(t, bitsToStringDefault(&bits) != "")
}

func TestBytesToDigestString(t *testing.T) {
	bits := Bits64{1}
	require.True(t, BytesToDigestString(&bits, "abc") != "")
}

func TestCopyToFixedBits(t *testing.T) {
	var d, s []byte
	copyToFixedBits(d, s, 0)
	require.Len(t, d, 0)

	require.Len(t, s, 0)

	d = make([]byte, 1)
	copyToFixedBits(d, s, 0)
	require.Len(t, d, 1)

	require.Len(t, s, 0)

	d = nil
	s = []byte{1}
	copyToFixedBits(d, s, 1)
	require.Len(t, d, 0)

	require.Len(t, s, 1)

	d = make([]byte, 1)
	copyToFixedBits(d, s, 1)
	require.Len(t, s, 1)

	require.Len(t, s, 1)

	require.Equal(t, s[0], d[0])

	require.Equal(t, uint8(1), d[0])

	d = make([]byte, 2)
	copyToFixedBits(d, s, 1)
	require.Len(t, d, 2)

	require.Len(t, s, 1)

	d = make([]byte, 1)
	s = []byte{1, 2}
	copyToFixedBits(d, s, 2)
	require.Len(t, d, 1)

	require.Len(t, s, 2)

	require.Equal(t, s[0], d[0])

	require.Equal(t, uint8(1), d[0])

	require.Panics(t, func() { copyToFixedBits(d, s, 3) })
}

func TestNewBits64FromBytes(t *testing.T) {
	bytes := []byte{}
	for i := 0; i < 8; i++ {
		bytes = append(bytes, byte(i%8))
	}
	bits := NewBits64FromBytes(bytes)
	require.Equal(t, bytes, bits.AsBytes())

	require.Panics(t, func() { NewBits64FromBytes([]byte{1}) })
}

func TestNewBits128FromBytes(t *testing.T) {
	bytes := []byte{}
	for i := 0; i < 16; i++ {
		bytes = append(bytes, byte(i%8))
	}
	bits := NewBits128FromBytes(bytes)
	require.Equal(t, bytes, bits.AsBytes())

	require.Panics(t, func() { NewBits128FromBytes([]byte{1}) })
}

func TestNewBits224FromBytes(t *testing.T) {
	bytes := []byte{}
	for i := 0; i < 24; i++ {
		bytes = append(bytes, byte(i%8))
	}
	bits := NewBits224FromBytes(bytes)
	require.Equal(t, bytes, bits.AsBytes())

	require.Panics(t, func() { NewBits224FromBytes([]byte{1}) })
}

func TestNewBits256FromBytes(t *testing.T) {
	bytes := []byte{}
	for i := 0; i < 32; i++ {
		bytes = append(bytes, byte(i%8))
	}
	bits := NewBits256FromBytes(bytes)
	require.Equal(t, bytes, bits.AsBytes())

	require.Panics(t, func() { NewBits256FromBytes([]byte{1}) })
}

func TestNewBits512FromBytes(t *testing.T) {
	bytes := []byte{}
	for i := 0; i < 64; i++ {
		bytes = append(bytes, byte(i%8))
	}
	bits := NewBits512FromBytes(bytes)
	require.Equal(t, bytes, bits.AsBytes())

	require.Panics(t, func() { NewBits512FromBytes([]byte{1}) })
}

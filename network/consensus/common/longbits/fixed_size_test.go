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
	"errors"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFoldUint64(t *testing.T) {
	require.Equal(t, uint32(0), FoldUint64(0))

	require.Equal(t, uint32(2), FoldUint64(2))

	require.Equal(t, uint32(math.MaxUint32), FoldUint64(math.MaxUint32))

	require.Equal(t, uint32(1), FoldUint64(math.MaxUint32+1))

	require.Equal(t, uint32(0), FoldUint64(math.MaxUint64))
}

func TestEqualFixedLenWriterTo(t *testing.T) {
	require.False(t, EqualFixedLenWriterTo(nil, nil))

	bits1 := NewBits64(0)
	require.False(t, EqualFixedLenWriterTo(&bits1, nil))

	require.False(t, EqualFixedLenWriterTo(nil, &bits1))

	bits2 := NewBits64(0)
	require.True(t, EqualFixedLenWriterTo(&bits1, &bits2))

	bits2 = NewBits64(1)
	require.False(t, EqualFixedLenWriterTo(&bits1, &bits2))
}

func TestCompare(t *testing.T) {
	require.False(t, (&writerToComparer{}).compare(nil, nil))

	bits1 := NewBits64(0)
	require.False(t, (&writerToComparer{}).compare(&bits1, nil))

	require.False(t, (&writerToComparer{}).compare(nil, &bits1))

	bits2 := NewBits64(0)
	require.True(t, (&writerToComparer{}).compare(&bits1, &bits2))

	bits3 := NewBits128(0, 0)
	require.False(t, (&writerToComparer{}).compare(&bits1, &bits3))

	bits1 = NewBits64(1)
	require.False(t, (&writerToComparer{}).compare(&bits1, &bits2))
}

func TestWrite(t *testing.T) {
	require.Panics(t, func() { _, _ = (&writerToComparer{}).Write(nil) })

	bits := NewBits64(0)
	fr := NewFixedReaderMock(t)
	fr.WriteToMock.Set(func(io.Writer) (int64, error) { return 0, errors.New("test") })
	n, err := (&writerToComparer{other: fr}).Write(bits.AsBytes())
	require.NotEqual(t, nil, err)

	require.Equal(t, 0, n)

	n, err = (&writerToComparer{other: &bits}).Write(bits.AsBytes())
	require.Equal(t, nil, err)

	require.Equal(t, 8, n)
}

func TestAsByteString(t *testing.T) {
	fs := &fixedSize{}
	require.Equal(t, "", fs.AsByteString())

	fs = &fixedSize{data: []byte{'a', 'b', 'c'}}
	require.Equal(t, "abc", fs.AsByteString())
}

/*func TestWriteTo(t *testing.T) {
	fs := &fixedSize{data: []byte{}}

	bits := NewBits64(0)

	n, err := fs.WriteTo(&writerToComparer{other: &bits})
	require.True(t, err == nil)

	require.Equal(t, 8, n)
}*/

func TestRead(t *testing.T) {
	item := byte(3)
	fs := &fixedSize{data: []byte{item}}
	buf := make([]byte, 2)
	n, err := fs.Read(buf)
	require.Equal(t, 1, n)

	require.Equal(t, nil, err)

	require.Equal(t, item, buf[0])
}

func TestFoldToUint64(t *testing.T) {
	fs := &fixedSize{data: []byte{1}}
	require.Panics(t, func() { fs.FoldToUint64() })

	fs.data = append(fs.data, 2, 3, 4, 5, 6, 7, 8)
	require.Equal(t, uint64(0x807060504030201), fs.FoldToUint64())
}

func TestFixedByteSize(t *testing.T) {
	fs := &fixedSize{data: []byte{1, 2}}
	require.Equal(t, len(fs.data), fs.FixedByteSize())
}

func TestAsBytes(t *testing.T) {
	fs := &fixedSize{data: []byte{1, 2}}
	require.Len(t, fs.AsBytes(), len(fs.data))

	require.Equal(t, fs.data, fs.AsBytes())
}

func TestNewFixedReader(t *testing.T) {
	data := []byte{1, 2, 3}
	fr := NewFixedReader(data)
	require.Len(t, fr.AsBytes(), len(data))

	require.Equal(t, data[1], fr.AsBytes()[1])
}

func TestCopyFixedSize(t *testing.T) {
	item := 0x7777
	bits := NewBits64(uint64(item))
	fr := CopyFixedSize(&bits)

	require.Len(t, fr.AsBytes(), len(bits))

	require.Equal(t, uint8(item), fr.AsBytes()[0])

	require.Equal(t, bits[0], fr.AsBytes()[0])
}

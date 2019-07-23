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
	"bytes"
	"io"

	"github.com/insolar/insolar/network/consensus/common/args"
)

type Foldable interface {
	FoldToUint64() uint64
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/longbits.FixedReader -o . -s _mock.go

type FixedReader interface {
	io.WriterTo
	io.Reader
	AsBytes() []byte
	AsByteString() string

	FixedByteSize() int
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/longbits.FoldableReader -o . -s _mock.go

type FoldableReader interface {
	FixedReader
	Foldable
}

func FoldUint64(v uint64) uint32 {
	return uint32(v) ^ uint32(v>>32)
}

// TODO ?NeedFix - current implementation can only work for limited cases
func EqualFixedLenWriterTo(t, o FixedReader) bool {
	if args.IsNil(t) || args.IsNil(o) {
		return false
	}
	return (&writerToComparer{}).compare(t, o)
}

type writerToComparer struct {
	thisValue *[]byte
	other     io.WriterTo
	result    bool
}

func (c *writerToComparer) compare(this, other FixedReader) bool {
	c.thisValue = nil
	if this == nil || other == nil || this.FixedByteSize() != other.FixedByteSize() {
		return false
	}
	c.other = other
	_, _ = this.WriteTo(c)
	return c.other == nil && c.result
}

func (c *writerToComparer) Write(otherValue []byte) (int, error) {
	if c.other == nil {
		panic("content of FixedReader must be read/written all at once")
	}
	if c.thisValue == nil {
		c.thisValue = &otherValue // result of &var is never nil
		_, err := c.other.WriteTo(c)
		if err != nil {
			return 0, err
		}
	} else {
		c.other = nil // mark "done"
		c.result = bytes.Equal(*c.thisValue, otherValue)
	}
	return len(otherValue), nil
}

type fixedSize struct {
	data []byte
}

func (c *fixedSize) AsByteString() string {
	return string(c.data)
}

func (c *fixedSize) WriteTo(w io.Writer) (n int64, err error) {
	n32, err := w.Write(c.data)
	return int64(n32), err
}

func (c *fixedSize) Read(p []byte) (n int, err error) {
	return copy(p, c.data), nil
}

func (c *fixedSize) FoldToUint64() uint64 {
	return FoldToUint64(c.data)
}

func (c *fixedSize) FixedByteSize() int {
	return len(c.data)
}

func (c *fixedSize) AsBytes() []byte {
	return c.data
}

func NewFixedReader(data []byte) FixedReader {
	return &fixedSize{data: data}
}

func CopyFixedSize(v FoldableReader) FoldableReader {
	r := fixedSize{}
	r.data = make([]byte, v.FixedByteSize())
	n, err := v.Read(r.data)
	if err != nil {
		panic(err)
	}
	if n != len(r.data) {
		panic("unexpected")
	}
	return &r
}

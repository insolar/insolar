// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package longbits

import (
	"bytes"
	"io"
)

type Foldable interface {
	FoldToUint64() uint64
}

//go:generate minimock -i github.com/insolar/insolar/longbits.FixedReader -o . -s _mock.go -g
type FixedReader interface {
	io.WriterTo
	io.Reader
	AsBytes() []byte
	AsByteString() ByteString

	FixedByteSize() int
}

//go:generate minimock -i github.com/insolar/insolar/longbits.FoldableReader -o . -s _mock.go -g
type FoldableReader interface {
	FixedReader
	Foldable
}

func FoldUint64(v uint64) uint32 {
	return uint32(v) ^ uint32(v>>32)
}

func EqualFixedLenWriterTo(t, o FixedReader) bool {
	if t == FixedReader(nil) || o == FixedReader(nil) {
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

func (c *fixedSize) AsByteString() ByteString {
	return ByteString(c.data)
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

func ReadFixedSize(v FoldableReader) []byte {
	data := make([]byte, v.FixedByteSize())
	n, err := v.Read(data)
	if err != nil {
		panic(err)
	}
	if n != len(data) {
		panic("unexpected")
	}
	return data
}

func NewFixedReader(data []byte) FixedReader {
	return &fixedSize{data: data}
}

func NewMutableFixedSize(data []byte) FixedReader {
	return &fixedSize{data}
}

func CopyToMutable(v FoldableReader) FoldableReader {
	return &fixedSize{ReadFixedSize(v)}
}

func NewImmutableFixedSize(data []byte) FixedReader {
	return NewByteString(data).AsReader()
}

func CopyToImmutable(v FoldableReader) FoldableReader {
	return NewByteString(ReadFixedSize(v)).AsReader()
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

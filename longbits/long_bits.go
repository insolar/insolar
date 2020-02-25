// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package longbits

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/bits"
	"strings"
)

var _ FoldableReader = &Bits128{}

const BitsStringPrefix = "0x"

type Bits64 [8]byte

func NewBits64(v uint64) Bits64 {
	r := Bits64{}
	binary.LittleEndian.PutUint64(r[:], v)
	return r
}

func (v *Bits64) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*v)[:])
	return int64(n), err
}

func (v *Bits64) Read(p []byte) (n int, err error) {
	return copy(p, (*v)[:]), nil
}

func (v *Bits64) FoldToUint64() uint64 {
	return binary.LittleEndian.Uint64(v[:])
}

func (v *Bits64) FixedByteSize() int {
	return len(*v)
}

func (v *Bits64) AsByteString() ByteString {
	return ByteString(v[:])
}

func (v Bits64) String() string {
	return bitsToStringDefault(&v)
}

func (v *Bits64) AsBytes() []byte {
	return v[:]
}

func (v Bits64) Compare(other Bits64) int {
	return bytes.Compare(v[:], other[:])
}

/* Array size doesnt need to be aligned */
func FoldToBits64(v []byte) Bits64 {
	var folded Bits64
	if len(v) == 0 {
		return folded
	}

	alignedLen := len(v) & (len(folded) - 1)
	copy(folded[alignedLen:], v)

	for i := 0; i < alignedLen; i += len(folded) {
		folded[0] ^= v[i+0]
		folded[1] ^= v[i+1]
		folded[2] ^= v[i+2]
		folded[3] ^= v[i+3]
		folded[4] ^= v[i+4]
		folded[5] ^= v[i+5]
		folded[6] ^= v[i+6]
		folded[7] ^= v[i+7]
	}
	return folded
}

func NewBits128(lo, hi uint64) Bits128 {
	r := Bits128{}
	binary.LittleEndian.PutUint64(r[:8], lo)
	binary.LittleEndian.PutUint64(r[8:], hi)
	return r
}

type Bits128 [16]byte

func (v *Bits128) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*v)[:])
	return int64(n), err
}

func (v *Bits128) Read(p []byte) (n int, err error) {
	return copy(p, (*v)[:]), nil
}

func (v *Bits128) FoldToUint64() uint64 {
	return FoldToUint64(v[:])
}

func (v *Bits128) FixedByteSize() int {
	return len(*v)
}

func (v Bits128) String() string {
	return bitsToStringDefault(&v)
}

func (v *Bits128) AsByteString() ByteString {
	return ByteString(v[:])
}

func (v *Bits128) AsBytes() []byte {
	return v[:]
}

func (v Bits128) Compare(other Bits128) int {
	return bytes.Compare(v[:], other[:])
}

type Bits224 [28]byte

func (v *Bits224) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*v)[:])
	return int64(n), err
}

func (v *Bits224) Read(p []byte) (n int, err error) {
	return copy(p, (*v)[:]), nil
}

func (v *Bits224) FoldToUint64() uint64 {
	return binary.LittleEndian.Uint64(v[:]) ^
		binary.LittleEndian.Uint64(v[8:]) ^
		binary.LittleEndian.Uint64(v[16:])
}

func (v *Bits224) FixedByteSize() int {
	return len(*v)
}

func (v Bits224) String() string {
	return bitsToStringDefault(&v)
}

func (v *Bits224) AsBytes() []byte {
	return v[:]
}

func (v *Bits224) AsByteString() ByteString {
	return ByteString(v[:])
}

func (v Bits224) Compare(other Bits224) int {
	return bytes.Compare(v[:], other[:])
}

type Bits256 [32]byte

func (v *Bits256) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*v)[:])
	return int64(n), err
}

func (v *Bits256) Read(p []byte) (n int, err error) {
	return copy(p, (*v)[:]), nil
}

func (v *Bits256) FoldToUint64() uint64 {
	return FoldToUint64(v[:])
}

func (v *Bits256) FoldToBits128() Bits128 {
	r := Bits128{}
	for i := range r {
		r[i] = v[i] ^ v[i+len(r)]
	}
	return r
}

func (v *Bits256) FoldToBits224() Bits224 {
	r := Bits224{}
	for i := range r {
		r[i] = v[i]
	}
	return r
}

func (v *Bits256) FixedByteSize() int {
	return len(*v)
}

func (v Bits256) String() string {
	return bitsToStringDefault(&v)
}

func (v *Bits256) AsBytes() []byte {
	return v[:]
}

func (v *Bits256) AsByteString() ByteString {
	return ByteString(v[:])
}

func (v Bits256) Compare(other Bits256) int {
	return bytes.Compare(v[:], other[:])
}

type Bits512 [64]byte

func (v *Bits512) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*v)[:])
	return int64(n), err
}

func (v *Bits512) Read(p []byte) (n int, err error) {
	return copy(p, (*v)[:]), nil
}

func (v *Bits512) FoldToUint64() uint64 {
	return FoldToUint64(v[:])
}

func (v *Bits512) FoldToBits256() Bits256 {
	r := Bits256{}
	for i := range r {
		r[i] = v[i] ^ v[i+len(r)]
	}
	return r
}

func (v *Bits512) FoldToBits224() Bits224 {
	r := Bits224{}
	for i := range r {
		r[i] = v[i] ^ v[i+32]
	}
	return r
}

func (v *Bits512) FixedByteSize() int {
	return len(*v)
}

func (v Bits512) String() string {
	return bitsToStringDefault(&v)
}

func (v *Bits512) AsBytes() []byte {
	return v[:]
}

func (v *Bits512) AsByteString() ByteString {
	return ByteString(v[:])
}

func (v Bits512) Compare(other Bits512) int {
	return bytes.Compare(v[:], other[:])
}

/* Array size must be aligned to 8 bytes */
func FoldToUint64(v []byte) uint64 {
	folded := FoldToBits64(v)
	return folded.FoldToUint64()
}

/*
This implementation DOES NOT provide secure random!
This function has a fixed implementation and MUST remain unchanged as some elements of Consensus rely on identical behavior of this functions.
Array size must be aligned to 8 bytes.
*/
func FillBitsWithStaticNoise(base uint32, v []byte) {

	if bits.OnesCount32(base) < 8 {
		base ^= 0x6206cc91 // add some noise
	}

	for i := uint32(0); i < uint32(len(v)); i += 8 {
		var n = base + i>>3
		u := uint64((^n) ^ (n << 16))
		u |= (u + 1) << 31
		u ^= u >> 1
		t := v[i:]
		binary.LittleEndian.PutUint64(t, u)
	}
}

func bitsToStringDefault(s FoldableReader) string {
	return BytesToDigestString(s, BitsStringPrefix)
	// return BytesToGroupedString(s.AsBytes(), BitsStringPrefix, "_", 8)
}

func BytesToDigestString(s FoldableReader, prefix string) string {
	return fmt.Sprintf("bits[%d]%s%08x", s.FixedByteSize()*8, prefix, s.FoldToUint64())
}

func BytesToGroupedString(s []byte, prefix string, separator string, everyN int) string {
	if everyN == 0 || len(separator) == 0 {
		return prefix + hex.EncodeToString(s)
	}

	buf := strings.Builder{}
	buf.WriteString(prefix)
	dst := make([]byte, hex.EncodedLen(len(s)))
	hex.Encode(dst, s)

	i := 0
	for i < len(s) {
		if i > 0 {
			buf.WriteString(separator)
		}
		n := i + everyN
		if i < len(s) {
			buf.Write(dst[i:n])
		} else {
			buf.Write(dst[i:])
			break
		}
		i = n
	}
	return buf.String()
}

func copyToFixedBits(dst, src []byte, expectedSize int) {
	size := len(src)
	if size != expectedSize {
		panic(fmt.Sprintf("Length missmatch, expected: %d, actual: %d", expectedSize, size))
	}

	copy(dst, src)
}

func NewBits64FromBytes(bytes []byte) *Bits64 {
	b := Bits64{}
	copyToFixedBits(b[:], bytes, b.FixedByteSize())
	return &b
}

func NewBits128FromBytes(bytes []byte) *Bits128 {
	b := Bits128{}
	copyToFixedBits(b[:], bytes, b.FixedByteSize())
	return &b
}

func NewBits224FromBytes(bytes []byte) *Bits224 {
	b := Bits224{}
	copyToFixedBits(b[:], bytes, b.FixedByteSize())
	return &b
}

func NewBits256FromBytes(bytes []byte) *Bits256 {
	b := Bits256{}
	copyToFixedBits(b[:], bytes, b.FixedByteSize())
	return &b
}

func NewBits512FromBytes(bytes []byte) *Bits512 {
	b := Bits512{}
	copyToFixedBits(b[:], bytes, b.FixedByteSize())
	return &b
}

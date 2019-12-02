//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package protokit

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const MaxVarintSize = binary.MaxVarintLen64
const MaxFieldId = math.MaxUint32 >> 3

type WireTag uint32

func SafeWireTag(v uint64) (WireTag, error) {
	if v > math.MaxUint32 {
		return 0, fmt.Errorf("invalid wire tag, overflow, %x", v)
	}
	wt := WireTag(v)
	if wt.IsValid() {
		return wt, nil
	}
	return 0, fmt.Errorf("invalid wire tag: %v", v)
}

func (v WireTag) IsZero() bool {
	return v != 0
}

func (v WireTag) IsValid() bool {
	return v.Type().IsValid() && v.FieldId() > 0
}

func (v WireTag) Type() WireType {
	return WireType(v & maskWireType)
}

func (v WireTag) FieldId() int {
	return int(v >> lenWireType)
}

func (v WireTag) TagSize() int {
	return SizeVarint32(uint32(v))
}

func (v WireTag) MaxFieldSize() (bool, int) {
	switch minSize, maxSize := v.Type().DataSize(); {
	case minSize == maxSize:
		return true, maxSize + v.TagSize()
	case maxSize < math.MaxInt32:
		maxSize += v.TagSize()
		fallthrough
	default:
		return false, maxSize
	}
}

func (v WireTag) FieldSize(u uint64) uint64 {
	return v.Type().FieldSize(v.TagSize(), u)
}

func (v WireTag) FixedFieldSize() int {
	if fixedSize, maxSize := v.MaxFieldSize(); !fixedSize {
		panic("illegal state - not fixed size")
	} else {
		return maxSize
	}
}

func (v WireTag) EnsureFixedFieldSize(sz int) WireTag {
	if v.FixedFieldSize() != sz {
		panic("illegal value - size mismatched")
	}
	return v
}

func (v WireTag) _checkTag(expected WireTag) error {
	if v == expected {
		return nil
	}
	return fmt.Errorf("tag mismatch: actual=%v, expected=%v", expected, v)
}

func (v WireTag) CheckType(t WireType) error {
	switch {
	case !t.IsValid():
		panic("illegal value")
	case t != v.Type():
		return fmt.Errorf("type mismatch: actual=%v, expectedType=%v", v, t)
	}
	return nil
}

func (v WireTag) CheckTag(expected WireTag) error {
	if !expected.IsValid() {
		panic("illegal value")
	}
	return v._checkTag(expected)
}

func (v WireTag) Check(expectedType WireType, expectedId int) error {
	return v._checkTag(expectedType.Tag(expectedId))
}

func (v WireTag) EnsureType(expectedType WireType) {
	if err := v.CheckType(expectedType); err != nil {
		panic(err)
	}
}

func (v WireTag) EnsureTag(expected WireTag) {
	if err := v.CheckTag(expected); err != nil {
		panic(err)
	}
}

func (v WireTag) Ensure(expectedType WireType, expectedId int) {
	if err := v.Check(expectedType, expectedId); err != nil {
		panic(err)
	}
}

func (v WireTag) ExpectDecoded(x uint64, err error) error {
	if err != nil {
		return err
	}
	if wt, err := SafeWireTag(x); err != nil {
		return err
	} else {
		return v.CheckTag(wt)
	}
}

func (v WireTag) ExpectFrom(r io.ByteReader) error {
	return v.ExpectDecoded(DecodeVarint(r))
}

func (v WireTag) DecodeFrom(r io.ByteReader) (uint64, error) {
	switch v.Type() {
	case WireVarint, WireBytes:
		if err := v.ExpectFrom(r); err != nil {
			return 0, err
		}
		return DecodeVarint(r)
	case WireFixed64:
		if err := v.ExpectFrom(r); err != nil {
			return 0, err
		}
		return DecodeFixed64(r)
	case WireFixed32:
		if err := v.ExpectFrom(r); err != nil {
			return 0, err
		}
		return DecodeFixed32(r)
	default:
		panic("illegal value")
	}
}

func (v WireTag) EncodeTo(w io.ByteWriter, u uint64) error {
	if err := EncodeVarint(w, uint64(v)); err != nil {
		return err
	}
	switch v.Type() {
	case WireVarint, WireBytes:
		return EncodeVarint(w, u)
	case WireFixed64:
		return EncodeFixed64(w, u)
	case WireFixed32:
		if u > math.MaxUint32 {
			panic("illegal value")
		}
		return EncodeFixed32(w, uint32(u))
	default:
		panic("illegal value")
	}
}

func (v WireTag) String() string {
	if v == 0 {
		return "zeroTag"
	}
	return fmt.Sprintf("%d:%v", v.FieldId(), v.Type())
}

type WireType uint8

const (
	WireVarint WireType = iota
	WireFixed64
	WireBytes
	WireStartGroup
	WireEndGroup
	WireFixed32

	MaxWireType = iota - 1
)
const lenWireType = 3
const maskWireType = 1<<lenWireType - 1

func (v WireType) IsValid() bool {
	return v <= MaxWireType
}

func (v WireType) isValidByte(firstByte byte) bool {
	return v.IsValid() && firstByte&maskWireType == byte(v) && firstByte>>lenWireType > 0
}

func (v WireType) Tag(fieldId int) WireTag {
	if fieldId <= 0 || fieldId > MaxFieldId {
		panic("illegal value")
	}
	return WireTag(fieldId<<lenWireType | int(v))
}

func (v WireType) DataSize() (minSize, maxSize int) {
	switch v {
	case WireVarint:
		return 1, MaxVarintSize
	case WireBytes:
		return 1, math.MaxInt32
	case WireFixed64:
		return 8, 8
	case WireFixed32:
		return 4, 4
	default:
		panic("illegal value")
	}
}

func (v WireType) FieldSize(tagSize int, u uint64) uint64 {
	switch v {
	case WireVarint:
		return uint64(tagSize + SizeVarint64(u))
	case WireBytes:
		return uint64(tagSize+SizeVarint64(u)) + u
	case WireFixed64:
		return uint64(tagSize) + 8
	case WireFixed32:
		return uint64(tagSize) + 4
	default:
		panic("illegal value")
	}
}

func (v WireType) ExpectDecoded(x uint64, err error) (WireTag, error) {
	if err != nil {
		return 0, err
	}
	if wt, err := SafeWireTag(x); err != nil {
		return 0, err
	} else {
		return wt, wt.CheckType(v)
	}
}

func (v WireType) ExpectFrom(r io.ByteReader) (WireTag, error) {
	return v.ExpectDecoded(DecodeVarint(r))
}

func (v WireTag) _wrapDecoded(x uint64, err error) (WireTag, uint64, error) {
	return v, x, err
}

func (v WireType) DecodeFrom(r io.ByteReader) (WireTag, uint64, error) {
	switch v {
	case WireVarint, WireBytes:
		if wt, err := v.ExpectFrom(r); err != nil {
			return 0, 0, err
		} else {
			return wt._wrapDecoded(DecodeVarint(r))
		}
	case WireFixed64:
		if wt, err := v.ExpectFrom(r); err != nil {
			return 0, 0, err
		} else {
			return wt._wrapDecoded(DecodeFixed64(r))
		}
	case WireFixed32:
		if wt, err := v.ExpectFrom(r); err != nil {
			return 0, 0, err
		} else {
			return wt._wrapDecoded(DecodeFixed32(r))
		}
	default:
		panic("illegal value")
	}
}

func (v WireType) String() string {
	switch v {
	case WireVarint:
		return "varint"
	case WireFixed64:
		return "fixed64"
	case WireBytes:
		return "bytes"
	case WireStartGroup:
		return "groupStart"
	case WireEndGroup:
		return "groupEnd"
	case WireFixed32:
		return "fixed32"
	default:
		return fmt.Sprintf("unknown%d", v)
	}
}

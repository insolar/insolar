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
	"errors"
	"fmt"
	"io"
)

var errOverflow = errors.New("proto: integer overflow")

func DecodeVarint(r io.ByteReader) (uint64, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return decodeVarint(b, r)
}

// Continues to read Varint that was stated with the given (b)
// So we can't use binary.ReadUvarint(r) here

func decodeVarint(b byte, r io.ByteReader) (n uint64, err error) {
	v := uint64(b & 0x7F)

	for i := uint8(7); i < 64; i += 7 {
		if b&0x80 == 0 {
			return v, nil
		}
		if b, err = r.ReadByte(); err != nil {
			return 0, err
		}
		v |= uint64(b&0x7F) << i
	}

	if b > 1 {
		return 0, errOverflow
	}
	return v, nil
}

func TryDecodeTag(expectedType WireType, r io.ByteScanner) (WireTag, error) {
	b, err := r.ReadByte()
	switch {
	case err != nil:
		return 0, err
	case !expectedType.isValidByte(b):
		return 0, r.UnreadByte()
	}
	if v, err := decodeVarint(b, r); err != nil {
		return 0, err
	} else if tag, err := SafeWireTag(v); err != nil {
		return 0, err
	} else {
		return tag, nil
	}
}

func DecodeExpectedType(expectedType WireType, r io.ByteReader) (WireTag, error) {
	if v, err := DecodeVarint(r); err != nil {
		return 0, err
	} else if tag, err := SafeWireTag(v); err != nil {
		return 0, err
	} else if err = tag.CheckType(expectedType); err != nil {
		return 0, err
	} else {
		return tag, nil
	}
}

func DecodeExpectedTag(expected WireTag, r io.ByteReader) error {
	if v, err := DecodeVarint(r); err != nil {
		return err
	} else if tag, err := SafeWireTag(v); err != nil {
		return err
	} else if err = tag.CheckTag(expected); err != nil {
		return err
	} else {
		return nil
	}
}

func DecodeField(expected WireTag, r io.ByteReader) error {
	if v, err := DecodeVarint(r); err != nil {
		return err
	} else if tag, err := SafeWireTag(v); err != nil {
		return err
	} else if err = tag.CheckTag(expected); err != nil {
		return err
	} else {
		return nil
	}
}

func DecodeFixed32(r io.ByteReader) (uint64, error) {
	return decodeFixed32(r)
}

func DecodeFixed64(r io.ByteReader) (uint64, error) {
	if v, err := decodeFixed32(r); err != nil {
		return 0, err
	} else if v2, err := decodeFixed32(r); err != nil {
		return 0, err
	} else {
		return v2<<32 | v, nil
	}
}

func decodeFixed32(r io.ByteReader) (v uint64, err error) {
	var b byte
	if b, err = r.ReadByte(); err != nil {
		return 0, err
	}
	v = uint64(b)
	if b, err = r.ReadByte(); err != nil {
		return 0, err
	}
	v |= uint64(b) << 8
	if b, err = r.ReadByte(); err != nil {
		return 0, err
	}
	v |= uint64(b) << 16
	if b, err = r.ReadByte(); err != nil {
		return 0, err
	}
	v |= uint64(b) << 24
	return v, nil
}

func MatchStringByteReader(s string, r io.ByteReader) (bool, error) {
	n := len(s)
	if n == 0 {
		panic("illegal value")
	}
	for i := 0; i < n; i++ {
		if cr, err := r.ReadByte(); err != nil {
			return false, err
		} else if cr != s[i] {
			return false, nil
		}
	}
	return true, nil
}

func MatchString(s string, r io.Reader, msg string) error {
	n := len(s)
	if n == 0 {
		panic("illegal value")
	}
	buf := make([]byte, n)
	if nr, err := r.Read(buf); err != nil {
		return err
	} else if nr != n {
		return fmt.Errorf("insufficient read%s actual=%d expected=%d", msg, nr, n)
	}

	if s != string(buf) {
		return fmt.Errorf("mismatched read%s actual=%s expected=%s", msg, s, buf)
	}
	return nil
}

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

import "io"

func SizeVarint32(x uint32) int {
	switch {
	case x < 1<<7:
		return 1
	case x < 1<<14:
		return 2
	case x < 1<<21:
		return 3
	case x < 1<<28:
		return 4
	}
	return 5
}

// SizeVarint returns the varint encoding size of an integer.
func SizeVarint64(x uint64) int {
	switch {
	case x < 1<<7:
		return 1
	case x < 1<<14:
		return 2
	case x < 1<<21:
		return 3
	case x < 1<<28:
		return 4
	case x < 1<<35:
		return 5
	case x < 1<<42:
		return 6
	case x < 1<<49:
		return 7
	case x < 1<<56:
		return 8
	case x < 1<<63:
		return 9
	}
	return 10
}

func EncodeVarint(w io.ByteWriter, u uint64) error {
	for u > 0x7F {
		if err := w.WriteByte(byte(u&0x7F | 0x80)); err != nil {
			return err
		}
		u >>= 7
	}
	return w.WriteByte(byte(u))
}

func EncodeFixed64(w io.ByteWriter, u uint64) error {
	if err := EncodeFixed32(w, uint32(u)); err != nil {
		return err
	}
	return EncodeFixed32(w, uint32(u>>32))
}

func EncodeFixed32(w io.ByteWriter, u uint32) error {
	if err := w.WriteByte(byte(u)); err != nil {
		return err
	}
	u >>= 8
	if err := w.WriteByte(byte(u)); err != nil {
		return err
	}
	u >>= 8
	if err := w.WriteByte(byte(u)); err != nil {
		return err
	}
	u >>= 8
	if err := w.WriteByte(byte(u)); err != nil {
		return err
	}
	return nil
}

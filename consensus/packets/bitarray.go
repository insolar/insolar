/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package packets

import (
	"encoding/binary"
	"math"

	"github.com/pkg/errors"
)

type bitArray struct {
	array    []uint8
	bitsSize uint
}

func newBitArray(size uint) *bitArray {
	totalSize := uint64(math.Round(float64(size/8) + 0.5))
	return &bitArray{
		array:    make([]uint8, totalSize),
		bitsSize: uint(size),
	}
}

func (arr *bitArray) Len() uint {
	return arr.bitsSize
}

func (arr *bitArray) put(bit, index int) error {
	if uint(index) >= arr.bitsSize {
		return errors.New("[ put ] failed to put a bit. out of range")
	}
	block := uint8(index / 8)
	step := uint8(8 - index%8 - 1)

	mask := uint8(1 << step)
	if bit == 0 {
		arr.array[block] &= ^(mask)
	} else if bit == 1 {
		arr.array[block] |= mask
	}
	return nil
}

func (arr *bitArray) serialize() ([]byte, error) {
	result := allocateBuffer(int(math.Round(float64(arr.bitsSize/8) + 0.5)))
	for _, byte := range arr.array {
		err := binary.Write(result, defaultByteOrder, byte)
		if err != nil {
			return nil, errors.Wrap(err, "[ serialize] failed to serialize a bitarray")
		}
	}

	return result.Bytes(), nil
}

func (arr *bitArray) get(index int) (uint8, error) {
	if uint(index) >= arr.bitsSize {
		return 0, errors.New("failed to get a bit - index out of range")
	}

	bitsN := uint8(index / 8)
	step := uint8(8 - index%8 - 1)
	res := arr.array[bitsN] >> step

	return res & 0x01, nil
}

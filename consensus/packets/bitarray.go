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

const sizeOfBlock = 8

type bitArray struct {
	array    []uint8
	bitsSize int
}

func round(dividend, divider int) float64 {
	return math.Round(float64(dividend/divider) + 0.5)
}

func newBitArray(size int) *bitArray {
	totalSize := uint64(round(size, sizeOfBlock))
	return &bitArray{
		array:    make([]uint8, totalSize),
		bitsSize: int(size),
	}
}

func (arr *bitArray) Len() int {
	return arr.bitsSize
}

func (arr *bitArray) put(bit, index int) error {
	if index >= arr.bitsSize {
		return errors.New("[ put ] failed to put a bit. out of range")
	}
	block := getBlockInBitArray(index)
	step := getStepToMove(index)

	mask := uint8(1 << step)
	if bit == 0 {
		arr.array[block] &= ^(mask) // change index bit to 0
	} else if bit == 1 {
		arr.array[block] |= mask // change index bit to 1
	} else {
		return errors.New("trying to set a wrong bit value")
	}
	return nil
}

func (arr *bitArray) serialize() ([]byte, error) {
	result := allocateBuffer(int(round(arr.bitsSize, sizeOfBlock)))
	for _, byte := range arr.array {
		err := binary.Write(result, defaultByteOrder, byte)
		if err != nil {
			return nil, errors.Wrap(err, "[ serialize] failed to serialize a bitarray")
		}
	}

	return result.Bytes(), nil
}

func (arr *bitArray) get(index int) (uint8, error) {
	if index >= arr.bitsSize {
		return 0, errors.New("failed to get a bit - index out of range")
	}

	block := getBlockInBitArray(index)
	step := getStepToMove(index)
	res := arr.array[block] >> step // get bit by index from block

	return res & lastBitMask, nil
}

func getStepToMove(index int) uint8 {
	return uint8(sizeOfBlock - index%sizeOfBlock - 1)
}

func getBlockInBitArray(index int) uint8 {
	return uint8(index / sizeOfBlock)
}

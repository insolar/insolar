/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package packets

import (
	"bytes"
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
	if (dividend % divider) == 0 {
		return float64(dividend / divider)
	}
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

func (arr *bitArray) set(bit, index int) error {
	if index >= arr.bitsSize {
		return errors.New("[ set ] failed to set a bit. out of range")
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

func (arr *bitArray) serialize(compressed bool) ([]byte, error) {
	if compressed {
		return arr.serializeCompressed()
	}
	var result bytes.Buffer
	for _, b := range arr.array {
		err := binary.Write(&result, defaultByteOrder, b)
		if err != nil {
			return nil, errors.Wrap(err, "[ serialize ] failed to serialize a bitarray")
		}
	}

	return result.Bytes(), nil
}

func (arr *bitArray) serializeCompressed() ([]byte, error) {
	var result bytes.Buffer
	last, err := arr.GetState(0)
	if err != nil {
		return nil, errors.Wrap(err, "[ serializeCompressed ] failed to get state from bitarray")
	}
	count := uint16(1)
	for i := 1; i < arr.bitsSize/2; i++ { // cuz 2 bits == 1 state
		current, err := arr.GetState(i)
		if err != nil {
			return nil, errors.Wrap(err, "[ serializeCompressed ] failed to get state from bitarray")
		}
		if (last != current) || (i+1 >= arr.bitsSize/2) {
			err := binary.Write(&result, binary.BigEndian, count)
			if err != nil {
				return nil, errors.Wrap(err, "[ serializeCompressed ] failed to write to buffer")
			}
			err = binary.Write(&result, binary.BigEndian, last)
			if err != nil {
				return nil, errors.Wrap(err, "[ serializeCompressed ] failed to write to buffer")
			}
			count = 1
			last = current
		} else {
			count++
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

func (arr *bitArray) GetState(index int) (TriState, error) {
	if index >= arr.bitsSize {
		return 0, errors.New("failed to get a bit - index out of range")
	}

	stateFirstBit, err := arr.get(2 * index)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetState ] failed to get a bit from bitarray")
	}
	stateSecondBit, err := arr.get(2*index + 1)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetState ] failed to get a bit from bitarray")
	}
	result := (stateFirstBit << 1) + stateSecondBit

	return TriState(result), nil
}

func (arr *bitArray) SetState(index int, state TriState) error {
	err := arr.putLastBit(state>>1, 2*index) // set first bit to array
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to set last bit")
	}
	err = arr.putLastBit(state, 2*index+1) // set second bit to array
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to set last bit")
	}
	return nil
}

func (arr *bitArray) putLastBit(state TriState, index int) error {
	bit := int(state & lastBitMask)
	return arr.set(bit, index)
}

func getStepToMove(index int) uint8 {
	return uint8(sizeOfBlock - index%sizeOfBlock - 1)
}

func getBlockInBitArray(index int) uint8 {
	return uint8(index / sizeOfBlock)
}

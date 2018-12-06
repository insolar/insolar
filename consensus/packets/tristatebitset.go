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
	"bytes"
	"encoding/binary"
	"io"
	"math/bits"

	"github.com/pkg/errors"
)

const lastBitMask = 0x01
const lowLengthSize = 6
const firstBitMask = 0x80
const lowBitLengthSize = 6

// TriStateBitSet bitset implementation.
type TriStateBitSet struct {
	CompressedSet bool

	array *bitArray
}

// NewTriStateBitSet creates and returns a tristatebitset.
func NewTriStateBitSet(size int) (*TriStateBitSet, error) {
	bitset := &TriStateBitSet{
		array: newBitArray(size * 2),
	}
	for i := 0; i < size; i++ {
		err := bitset.changeBitState(i, TimedOut)
		if err != nil {
			return nil, err
		}
	}
	return bitset, nil
}

func (dbs *TriStateBitSet) GetCells(mapper BitSetMapper) ([]BitSetCell, error) {
	return dbs.parseCells(mapper)
}

func (dbs *TriStateBitSet) ApplyChanges(changes []BitSetCell, mapper BitSetMapper) error {
	for _, cell := range changes {
		index, err := mapper.RefToIndex(cell.NodeID)
		if err != nil {
			return errors.Wrap(err, "[ ApplyChanges ] failed to get index from ref")
		}
		err = dbs.changeBucketState(index, cell.State)
		if err != nil {
			return errors.Wrap(err, "[ ApplyChanges ] failed to change bucket state")
		}
	}
	return nil
}

func (dbs *TriStateBitSet) Serialize() ([]byte, error) {
	var firstByte uint8 // compressed and hBitLength bits
	if dbs.CompressedSet {
		firstByte = 0x01
	} else {
		firstByte = 0x00
	}

	data, err := dbs.array.serialize(dbs.CompressedSet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to serialize a bitarray")
	}

	length := len(data)
	var result bytes.Buffer
	firstByte = firstByte << 1
	if bits.Len(uint(length)) > lowLengthSize {
		err = dbs.serializeWithHLength(firstByte, length, &result)
		if err != nil {
			return nil, errors.Wrap(err, "[ Serialize ] failed to serialize first bytes")
		}
	} else {
		err = dbs.serializeWithLLength(firstByte, length, &result)
		if err != nil {
			return nil, errors.Wrap(err, "[ Serialize ] failed to serialize first bytes")
		}
	}

	err = binary.Write(&result, defaultByteOrder, data)
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to write binary")
	}

	return result.Bytes(), nil
}

func (dbs *TriStateBitSet) serializeWithHLength(
	firstByte uint8,
	tmpLen int,
	result *bytes.Buffer,
) error {
	var secondByte uint8 // hBitLength
	firstByte++
	firstByte = firstByte << lowBitLengthSize // move compressed and hBitLength bits to right
	secondByte = uint8(tmpLen)
	err := binary.Write(result, defaultByteOrder, firstByte)
	if err != nil {
		return errors.Wrap(err, "[ serializeWithHLength ] failed to write binary")
	}
	err = binary.Write(result, defaultByteOrder, secondByte)
	if err != nil {
		return errors.Wrap(err, "[ serializeWithHLength ] failed to write binary")
	}
	return nil
}

func (dbs *TriStateBitSet) serializeWithLLength(
	firstByte uint8,
	tmpLen int,
	result *bytes.Buffer,
) error {
	firstByte = firstByte << lowLengthSize // move compressed and hbit flags to right
	firstByte += uint8(tmpLen)
	err := binary.Write(result, defaultByteOrder, firstByte)
	if err != nil {
		return errors.Wrap(err, "[ serializeWithLLength ] failed to write binary")
	}
	return nil
}

func DeserializeBitSet(data io.Reader) (BitSet, error) {
	firstbyte := uint8(0)
	err := binary.Read(data, defaultByteOrder, &firstbyte)
	var array *bitArray
	if err != nil {
		return nil, errors.Wrap(err, "[ Deserialize ] failed to read first byte")
	}
	compressed, hbitFlag, length := parseFirstByte(firstbyte)
	if hbitFlag {
		err = binary.Read(data, defaultByteOrder, &length)
		if err != nil {
			return nil, errors.Wrap(err, "[ Deserialize ] failed to read second byte")
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "[ Deserialize ] failed to create a bitset")
	}
	if compressed {
		array, err = deserializeCompressed(data, int(length))
		if err != nil {
			return nil, errors.Wrap(err, "[ DeserializeBitSet ] failed to deserialize a compressed bitarray")
		}
	} else {
		payload := make([]uint8, length)
		for i := uint8(0); i < length; i++ {
			err := binary.Read(data, defaultByteOrder, &payload[i])
			if err != nil {
				return nil, errors.Wrap(err, "[ Deserialize ] failed to read payload")
			}
		}
		array, err = parseBitArray(payload)
		if err != nil {
			return nil, errors.Wrap(err, "[ Deserialize ] failed to parse a bitarray")
		}
	}
	bitset := &TriStateBitSet{
		array: array,
	}
	return bitset, nil
}

func (dbs *TriStateBitSet) parseCells(mapper BitSetMapper) ([]BitSetCell, error) {
	cellSize := int(dbs.array.bitsSize / 2)
	cells := make([]BitSetCell, cellSize)
	for i := 0; i < cellSize; i++ {
		id, err := mapper.IndexToRef(i)
		if err != nil {
			return nil, err
		}
		cells[i].NodeID = id
		cells[i].State, err = parseState(dbs.array, i)
		if err != nil {
			return nil, errors.Wrap(err, "[ parseCells ] failed to parse TriState")
		}
	}
	return cells, nil
}

func deserializeCompressed(data io.Reader, size int) (*bitArray, error) {
	count := uint8(0)
	value := uint8(0)
	var payload []uint8
	blockSize := 0
	block := uint8(0)
	var err error
	for i := 0; i < size; i = i + 2 {
		err = binary.Read(data, binary.BigEndian, &count)
		if err != nil {
			return nil, errors.Wrap(err, "[ deserializeCompressed ] failed to read from data")
		}
		err = binary.Read(data, binary.BigEndian, &value)
		if err != nil {
			return nil, errors.Wrap(err, "[ deserializeCompressed ] failed to read from data")
		}
		for j := uint8(0); j < count; j++ {
			block += value
			blockSize += 2
			if (blockSize >= sizeOfBlock) || (j+1 >= count) {
				if j+1 >= count {
					block = block << uint(sizeOfBlock-blockSize)
				}
				payload = append(payload, block)
				blockSize = 0
			}
			block = block << 2
		}
	}
	return parseBitArray(payload)
}

func parseState(array *bitArray, index int) (TriState, error) {
	state, err := array.getState(index)
	return TriState(state), err
}

func parseBitArray(payload []uint8) (*bitArray, error) {
	len := len(payload)
	array := newBitArray(len*sizeOfBlock - 4) // bits count from bytes size
	for i := 0; i < len; i++ {
		array.array[i] = payload[i]
	}
	return array, nil
}

func (dbs *TriStateBitSet) changeBucketState(index int, newState TriState) error {
	return dbs.changeBitState(index, newState)
}

func (dbs *TriStateBitSet) putLastBit(state TriState, index int) error {
	bit := int(state & lastBitMask)
	return dbs.array.put(bit, index)
}

func (dbs *TriStateBitSet) changeBitState(i int, state TriState) error {
	err := dbs.putLastBit(state>>1, 2*i) // put first bit to array
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to put last bit")
	}
	err = dbs.putLastBit(state, 2*i+1) // put second bit to array
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to put last bit")
	}
	return nil
}

func parseFirstByte(b uint8) (compressed bool, hbitFlag bool, lbitLength uint8) {
	lbitLength = uint8(0)
	compressed = false
	hbitFlag = false
	if (b & firstBitMask) == firstBitMask { // check compressed flag bit
		compressed = true
	}
	check := (b << 1) & firstBitMask // check hBitLength flag bit
	if check == firstBitMask {
		hbitFlag = true
		return
	}
	lbitLength = (b << 2) >> 2 // remove 2 first bits
	return
}

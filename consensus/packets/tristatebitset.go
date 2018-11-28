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

	totalSize := int(round(dbs.array.Len()*2, sizeOfBlock)) + 1 // size of result bytes
	var result *bytes.Buffer
	var err error
	firstByte = firstByte << 1
	if bits.Len(uint(dbs.array.Len())) > lowLengthSize {
		result, err = dbs.serializeWithHLength(firstByte, dbs.array.Len(), totalSize)
		if err != nil {
			return nil, errors.Wrap(err, "[ Serialize ] failed to serialize first bytes")
		}
	} else {
		result, err = dbs.serializeWithLLength(firstByte, dbs.array.Len(), totalSize)
		if err != nil {
			return nil, errors.Wrap(err, "[ Serialize ] failed to serialize first bytes")
		}
	}

	data, err := dbs.array.serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to serialize a bitarray")
	}
	err = binary.Write(result, defaultByteOrder, data)
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to write binary")
	}

	return result.Bytes(), nil
}

func (dbs *TriStateBitSet) serializeWithHLength(
	firstByte uint8,
	tmpLen int,
	totalSize int,
) (res *bytes.Buffer, err error) {
	var result *bytes.Buffer
	var secondByte uint8 // hBitLength
	firstByte++
	firstByte = firstByte << lowBitLengthSize // move compressed and hBitLength bits to right
	secondByte = uint8(tmpLen)
	totalSize++ // secondbyte is optional
	result = allocateBuffer(totalSize)
	err = binary.Write(result, defaultByteOrder, firstByte)
	if err != nil {
		return result, errors.Wrap(err, "[ serializeWithHLength ] failed to write binary")
	}
	err = binary.Write(result, defaultByteOrder, secondByte)
	if err != nil {
		return result, errors.Wrap(err, "[ serializeWithHLength ] failed to write binary")
	}
	return result, nil
}

func (dbs *TriStateBitSet) serializeWithLLength(
	firstByte uint8,
	tmpLen int,
	totalSize int,
) (res *bytes.Buffer, err error) {
	result := allocateBuffer(totalSize)
	firstByte = firstByte << 1 // move compressed flag to right
	firstByte += uint8(tmpLen)
	err = binary.Write(result, defaultByteOrder, firstByte)
	if err != nil {
		return nil, errors.Wrap(err, "[ serializeWithLLength ] failed to write binary")
	}
	return result, nil
}

func DeserializeBitSet(data io.Reader) (BitSet, error) {
	firstbyte := uint8(0)
	err := binary.Read(data, defaultByteOrder, &firstbyte)
	if err != nil {
		return nil, errors.Wrap(err, "[ Deserialize ] failed to read first byte")
	}
	compressed, hbitFlag, length := parseFirstByte(firstbyte)
	if compressed {
		panic("[ DeserializeBitSet ] not implemented yet")
	}
	if hbitFlag {
		err = binary.Read(data, defaultByteOrder, &length)
		if err != nil {
			return nil, errors.Wrap(err, "[ Deserialize ] failed to read second byte")
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "[ Deserialize ] failed to create a bitset")
	}
	blockCount := uint64(round(int(length), sizeOfBlock))
	payload := make([]uint8, blockCount)
	for i := 0; uint64(i) < blockCount; i++ {
		err := binary.Read(data, defaultByteOrder, &payload[i])
		if err != nil {
			return nil, errors.Wrap(err, "[ Deserialize ] failed to read first byte")
		}
	}
	array, err := parseBitArray(payload, int(length))
	if err != nil {
		return nil, errors.Wrap(err, "[ Deserialize ] failed to parse a bitarray")
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

func parseState(array *bitArray, index int) (TriState, error) {
	stateFirstBit, err := array.get(2 * index)
	if err != nil {
		return 0, err
	}
	stateSecondBit, err := array.get(2*index + 1)
	if err != nil {
		return 0, err
	}
	return TriState((stateFirstBit << 1) + stateSecondBit), nil
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
	if (b & firstBitMask) == 1 { // check compressed flag bit
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

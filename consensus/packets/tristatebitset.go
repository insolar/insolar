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

	cells  []BitSetCell
	mapper BitSetMapper
}

// NewTriStateBitSet creates and returns a tristatebitset.
func NewTriStateBitSet(cells []BitSetCell, mapper BitSetMapper) (*TriStateBitSet, error) {
	if (mapper == nil) || (cells == nil) {
		return nil, errors.New("[ NewTriStateBitSet ] failed to create tristatebitset")
	}
	bitset := &TriStateBitSet{
		cells:  make([]BitSetCell, mapper.Length()),
		mapper: mapper,
	}
	bitset.ApplyChanges(cells)
	return bitset, nil
}

func (dbs *TriStateBitSet) GetCells() []BitSetCell {
	return dbs.cells
}

func (dbs *TriStateBitSet) ApplyChanges(changes []BitSetCell) {
	for _, cell := range changes {
		err := dbs.changeBucketState(&cell)
		if err != nil {
			panic(err)
		}
	}
}

func (dbs *TriStateBitSet) Serialize() ([]byte, error) {
	var firstByte uint8 // compressed and hBitLength bits
	if dbs.CompressedSet {
		firstByte = 0x01
	} else {
		firstByte = 0x00
	}

	array, err := dbs.cellsToBitArray()
	if err != nil {
		return nil, errors.Wrap(err, "[ Serialize ] failed to get bitarray from cells")
	}

	totalSize := int(round(array.Len()*2, sizeOfBlock)) + 1 // size of result bytes
	var result *bytes.Buffer
	firstByte = firstByte << 1
	if bits.Len(uint(array.Len())) > lowLengthSize {
		result, err = dbs.serializeWithHLength(firstByte, array.Len(), totalSize)
		if err != nil {
			return nil, errors.Wrap(err, "[ Serialize ] failed to serialize first bytes")
		}
	} else {
		result, err = dbs.serializeWithLLength(firstByte, array.Len(), totalSize)
		if err != nil {
			return nil, errors.Wrap(err, "[ Serialize ] failed to serialize first bytes")
		}
	}

	data, err := array.serialize()
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

func (dbs *TriStateBitSet) Deserialize(data io.Reader) error {
	var firstbyte uint8
	err := binary.Read(data, defaultByteOrder, &firstbyte)
	if err != nil {
		return errors.Wrap(err, "[ Deserialize ] failed to read first byte")
	}
	compressed, hbitFlag, length := parseFirstByte(firstbyte)
	if hbitFlag {
		err = binary.Read(data, defaultByteOrder, &length)
		if err != nil {
			return errors.Wrap(err, "[ Deserialize ] failed to read second byte")
		}
	}
	blockCount := uint64(round(int(length), sizeOfBlock))
	payload := make([]uint8, blockCount)
	for i := 0; uint64(i) < blockCount; i++ {
		err = binary.Read(data, defaultByteOrder, &payload[i])
		if err != nil {
			return errors.Wrap(err, "[ Deserialize ] failed to read first byte")
		}
	}
	var array *bitArray
	if compressed {
		panic("we have no implementation for this branch")
	} else {
		array, err = parseBitArray(payload, int(length))
		if err != nil {
			return err
		}
	}
	cells, err := dbs.parseCells(array)
	if err != nil {
		return err
	}

	dbs.cells = cells
	dbs.CompressedSet = compressed
	return nil
}

func (dbs *TriStateBitSet) parseCells(array *bitArray) ([]BitSetCell, error) {
	cellSize := int(array.bitsSize / 2)
	cells := make([]BitSetCell, cellSize)
	for i := 0; i < cellSize; i++ {
		id, err := dbs.mapper.IndexToRef(i)
		if err != nil {
			return nil, err
		}
		cells[i].NodeID = id
		cells[i].State, err = parseState(array, i)
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

func parseBitArray(payload []uint8, size int) (*bitArray, error) {
	array := newBitArray(size)
	for i := 0; i < size; i++ {
		block := getBlockInBitArray(i)
		step := getStepToMove(i)
		bit := (payload[block] >> step) & lastBitMask
		err := array.put(int(bit), i)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func parseFirstByte(byte uint8) (compressed bool, hbitFlag bool, lbitLength uint8) {
	lbitLength = uint8(0)
	compressed = false
	hbitFlag = false
	if (byte & firstBitMask) == 1 { // check compressed flag bit
		compressed = true
	}
	check := (byte << 1) & firstBitMask // check hBitLength flag bit
	if check == firstBitMask {
		hbitFlag = true
		return
	}
	lbitLength = (byte << 2) >> 2 // remove 2 first bits
	return
}

func (dbs *TriStateBitSet) changeBucketState(cell *BitSetCell) error {
	n, err := dbs.mapper.RefToIndex(cell.NodeID)
	if err != nil {
		return errors.Wrap(err, "[ changeBucketState ] failed to get index from ref")
	}
	dbs.cells[n] = *cell
	return nil
}

func putLastBit(array *bitArray, state TriState, index int) error {
	bit := int(state & lastBitMask)
	return array.put(bit, index)
}

func changeBitState(array *bitArray, i int, state TriState) error {
	err := putLastBit(array, state>>1, 2*i) // put first bit to array
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to put last bit")
	}
	err = putLastBit(array, state, 2*i+1) // put second bit to array
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to put last bit")
	}
	return nil
}

func (dbs *TriStateBitSet) cellsToBitArray() (*bitArray, error) {
	array := newBitArray(dbs.mapper.Length() * 2)
	for i := 0; i < len(dbs.cells); i++ {
		cell := dbs.cells[i]
		index, err := dbs.mapper.RefToIndex(cell.NodeID)
		if err != nil {
			return nil, errors.Wrap(err, "[ cellsToBitArray ] failed to get index from ref")
		}
		err = changeBitState(array, index, cell.State)
		if err != nil {
			return nil, errors.Wrap(err, "[ cellsToBitArray ] failed to change bit state")
		}
	}
	return array, nil
}

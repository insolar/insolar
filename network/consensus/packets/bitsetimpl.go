//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package packets

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/bits"

	"github.com/pkg/errors"
)

const lowLengthSize = 6
const firstBitMask = 0x80
const last6BitsMask = 0x3f
const lastTwoBitsMask = 0x3

// bitsetImpl bitset implementation
type bitsetImpl struct {
	compressed bool
	array      bitArray
}

func (dbs *bitsetImpl) Clone() BitSet {
	clone := bitsetImpl{compressed: dbs.compressed}
	clone.array = make(bitArray, len(dbs.array))
	copy(clone.array, dbs.array)
	return &clone
}

func (dbs *bitsetImpl) GetCells(mapper BitSetMapper) ([]BitSetCell, error) {
	cells := make([]BitSetCell, len(dbs.array))
	for i := 0; i < len(dbs.array); i++ {
		id, err := mapper.IndexToRef(i)
		if err != nil {
			return nil, err
		}
		cells[i] = BitSetCell{NodeID: id, State: dbs.array[i]}
	}
	return cells, nil
}

// NewBitSetImpl creates and returns a bitset implementation
func NewBitSetImpl(size int, compressed bool) (BitSet, error) {
	bitset := &bitsetImpl{
		compressed: compressed,
		array:      make(bitArray, size),
	}
	for i := 0; i < size; i++ {
		bitset.array[i] = TimedOut
	}
	return bitset, nil
}

func (dbs *bitsetImpl) GetTristateArray() ([]BitSetState, error) {
	result := make([]BitSetState, len(dbs.array))
	copy(result, dbs.array)
	return result, nil
}

func (dbs *bitsetImpl) ApplyChanges(changes []BitSetCell, mapper BitSetMapper) error {
	for _, cell := range changes {
		index, err := mapper.RefToIndex(cell.NodeID)
		if err != nil {
			return errors.Wrapf(err, "failed to map reference %s to bitset index", cell.NodeID)
		}
		dbs.array[index] = cell.State
	}
	return nil
}

func (dbs *bitsetImpl) Serialize() ([]byte, error) {
	var firstByte uint8 // compressed and hBitLength bits
	if dbs.compressed {
		firstByte = 0x01
	} else {
		firstByte = 0x00
	}

	data, err := dbs.array.Serialize(dbs.compressed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize bitarray")
	}

	length := len(dbs.array)
	var result bytes.Buffer
	firstByte <<= 1
	if bits.Len(uint(length)) > lowLengthSize {
		err = dbs.serializeWithHLength(firstByte, length, &result)
		if err != nil {
			return nil, errors.Wrap(err, "failed to write 2-byte bitset header")
		}
	} else {
		err = dbs.serializeWithLLength(firstByte, length, &result)
		if err != nil {
			return nil, errors.Wrap(err, "failed to write 1-byte bitset header")
		}
	}

	err = binary.Write(&result, defaultByteOrder, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write bitset binary representation")
	}

	return result.Bytes(), nil
}

func (dbs *bitsetImpl) serializeWithHLength(firstByte uint8, length int, result io.Writer) error {
	var secondByte uint8 // hBitLength
	firstByte++
	firstByte <<= lowLengthSize // move compressed and hBitLength bits to right
	secondByte = uint8(length & 0xff)
	lowByte := uint8(length >> 8)
	if lowByte != 0 {
		lowByte &= last6BitsMask
		firstByte |= lowByte
	}
	err := binary.Write(result, defaultByteOrder, firstByte)
	if err != nil {
		return errors.Wrap(err, "failed to write first byte of bitset header")
	}
	err = binary.Write(result, defaultByteOrder, secondByte)
	if err != nil {
		return errors.Wrap(err, "failed to write second byte of bitset header")
	}
	return nil
}

func (dbs *bitsetImpl) serializeWithLLength(firstByte uint8, length int, result io.Writer) error {
	firstByte <<= lowLengthSize // move compressed and hbit flags to right
	firstByte += uint8(length)
	err := binary.Write(result, defaultByteOrder, firstByte)
	if err != nil {
		return errors.Wrap(err, "failed to write first byte of bitset header")
	}
	return nil
}

func DeserializeBitSet(data io.Reader) (BitSet, error) {
	compressed, length, err := parseHeader(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse bitset header")
	}

	var array bitArray
	if compressed {
		array, err = deserializeCompressed(data, length)
		if err != nil {
			return nil, errors.Wrap(err, "failed to deserialize compressed bitarray")
		}
	} else {
		payload := make([]uint8, div(length, statesInByte))
		err := binary.Read(data, defaultByteOrder, &payload)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read bitarray payload")
		}
		array, err = deserialize(payload, length)
		if err != nil {
			return nil, errors.Wrap(err, "failed to deserialize bitarray")
		}
	}
	bitset := &bitsetImpl{
		compressed: compressed,
		array:      array,
	}
	return bitset, nil
}

func parseHeader(data io.Reader) (bool, int, error) {
	firstbyte := uint8(0)
	err := binary.Read(data, defaultByteOrder, &firstbyte)
	if err != nil {
		return false, 0, errors.Wrap(err, "failed to read first byte")
	}
	var length int
	compressed, hbitFlag, lowLength := parseFirstByte(firstbyte)
	if hbitFlag {
		var highLength uint8
		err = binary.Read(data, defaultByteOrder, &highLength)
		if err != nil {
			return false, 0, errors.Wrap(err, "failed to read second byte")
		}
		length = int(lowLength)<<8 | int(highLength)
	} else {
		length = int(lowLength)
	}
	if length < 0 || length > 1024 {
		return false, 0, errors.Errorf("got bitset with incorrect size: %d", length)
	}
	return compressed, length, nil
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
	}
	lbitLength = (b << 2) >> 2 // remove 2 first bits
	return
}

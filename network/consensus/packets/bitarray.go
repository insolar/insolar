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
	"math"

	"github.com/pkg/errors"
)

const bitsForState = 2
const statesInByte = 4

type bitArray []BitSetState

func div(dividend, divider int) int { // nolint: unparam
	if (dividend % divider) == 0 {
		return dividend / divider
	}
	return int(math.Round(float64(dividend/divider) + 0.5))
}

func parseStatesFromByte(b uint8) [statesInByte]BitSetState {
	var result [statesInByte]BitSetState
	for i := statesInByte; i > 0; i-- {
		result[i-1] = BitSetState(b & lastTwoBitsMask)
		b >>= bitsForState
	}
	return result
}

func deserialize(data []byte, length int) (bitArray, error) {
	if len(data) != div(length, statesInByte) {
		return nil, errors.Errorf("wrong size of data buffer, expected: %d, got: %d", div(length, statesInByte), len(data))
	}
	result := make(bitArray, length)
	statesLeft := length
	for i := 0; i < len(data); i++ {
		parsedStates := parseStatesFromByte(data[i])
		statesRange := statesInByte
		if statesLeft < statesInByte {
			statesRange = statesLeft
		}
		startIndex := i * statesInByte
		copy(result[startIndex:startIndex+statesRange], parsedStates[:statesRange])
		statesLeft -= statesInByte
	}
	return result, nil
}

func deserializeCompressed(data io.Reader, size int) (bitArray, error) {
	count := uint16(0)
	index := 0
	var value BitSetState
	var err error

	statesLeft := uint16(size)
	result := make(bitArray, size)
	for statesLeft > 0 {
		err = binary.Read(data, binary.BigEndian, &count)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read states count from buffer")
		}
		err = binary.Read(data, binary.BigEndian, &value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read state from buffer")
		}
		for i := 0; i < int(count); i++ {
			result[index] = value
			index++
		}
		statesLeft -= count
	}
	return result, nil
}

func (ba bitArray) serialize() []byte {
	length := div(len(ba), statesInByte)
	result := make([]byte, length)
	for i := 0; i < len(result); i++ {
		startIndex := i * statesInByte
		endIndex := i*statesInByte + statesInByte
		if endIndex > len(ba) {
			endIndex = len(ba)
		}
		result[i] = writeStatesToByte(ba[startIndex:endIndex])
	}
	return result
}

func writeStatesToByte(states []BitSetState) uint8 {
	var result uint8
	result |= uint8(states[0]) & lastTwoBitsMask
	for i := 1; i < len(states); i++ {
		result <<= bitsForState
		result |= uint8(states[i]) & lastTwoBitsMask
	}
	for i := len(states); i < statesInByte; i++ {
		result <<= bitsForState
	}
	return result
}

func (ba bitArray) Serialize(compressed bool) ([]byte, error) {
	if compressed {
		return ba.serializeCompressed()
	}
	return ba.serialize(), nil
}

func (ba bitArray) serializeCompressed() ([]byte, error) {
	var result bytes.Buffer
	last := ba[0]
	count := 1
	for i := 1; i < len(ba); i++ {
		current := ba[i]
		if last == current {
			count++
			continue
		}

		err := ba.writeSequence(&result, last, count)
		if err != nil {
			return nil, err
		}
		count = 1
		last = current
	}
	err := ba.writeSequence(&result, last, count)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func (ba bitArray) writeSequence(buf io.Writer, state BitSetState, count int) error {
	err := binary.Write(buf, binary.BigEndian, uint16(count))
	if err != nil {
		return errors.Wrap(err, "failed to write states count to buffer")
	}
	err = binary.Write(buf, binary.BigEndian, state)
	if err != nil {
		return errors.Wrap(err, "failed to write state to buffer")
	}
	return nil
}

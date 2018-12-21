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

package sign

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
)

const (
	expectedBigIntBytesLength = 32
	sizeBytesLength           = 2
	TwoBigIntBytesLength      = (expectedBigIntBytesLength * 2) + sizeBytesLength
)

func SerializeTwoBigInt(one, two *big.Int) []byte {
	oneBytes := one.Bytes()
	twoBytes := two.Bytes()

	oneBytesLen := len(oneBytes)
	twoBytesLen := len(twoBytes)

	if oneBytesLen > expectedBigIntBytesLength || twoBytesLen > expectedBigIntBytesLength {
		err := fmt.Sprintf(
			"[ serializeTwoBigInt ] wrong one, two length. one: %d; two: %d; needed: %d. One was: %s, Two was: %s",
			oneBytesLen,
			twoBytesLen,
			expectedBigIntBytesLength,
			one.String(),
			two.String(),
		)
		panic(err)
	}

	var serialized [TwoBigIntBytesLength]byte

	serialized[0] = uint8(oneBytesLen)
	serialized[1] = uint8(twoBytesLen)

	oneStartPos := sizeBytesLength
	oneEndPos := oneStartPos + oneBytesLen
	copy(serialized[oneStartPos:oneEndPos], oneBytes)

	twoStartPos := sizeBytesLength + expectedBigIntBytesLength
	twoEndPos := twoStartPos + twoBytesLen
	copy(serialized[twoStartPos:twoEndPos], twoBytes)

	return serialized[:]
}

func DeserializeTwoBigInt(data []byte) (*big.Int, *big.Int, error) {
	if len(data) != TwoBigIntBytesLength {
		return nil, nil, errors.Errorf("[ DeserializeTwoBigInt ] wrong data length: %d", len(data))
	}

	var one, two big.Int

	oneBytesLen := int(data[0])
	twoBytesLen := int(data[1])

	if oneBytesLen > expectedBigIntBytesLength || twoBytesLen > expectedBigIntBytesLength {
		return nil, nil, errors.Errorf("[ DeserializeTwoBigInt ] wrong data to parse one len: %d, two len: %d", oneBytesLen, twoBytesLen)
	}

	oneStartPos := sizeBytesLength
	oneEndPos := oneStartPos + oneBytesLen

	twoStartPos := sizeBytesLength + expectedBigIntBytesLength
	twoEndPos := twoStartPos + twoBytesLen

	one.SetBytes(data[oneStartPos:oneEndPos])
	two.SetBytes(data[twoStartPos:twoEndPos])
	return &one, &two, nil
}

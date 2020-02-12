// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sign

import (
	"math/big"

	"github.com/pkg/errors"
)

const (
	expectedBigIntBytesLength = 32
	TwoBigIntBytesLength      = expectedBigIntBytesLength * 2
)

func SerializeTwoBigInt(one, two *big.Int) []byte {
	oneBytes := mustCanonicalizeInt(one)
	twoBytes := mustCanonicalizeInt(two)

	var serialized [TwoBigIntBytesLength]byte

	copy(serialized[:expectedBigIntBytesLength], oneBytes)
	copy(serialized[expectedBigIntBytesLength:TwoBigIntBytesLength], twoBytes)

	return serialized[:]
}

func DeserializeTwoBigInt(data []byte) (*big.Int, *big.Int, error) {
	if len(data) != TwoBigIntBytesLength {
		return nil, nil, errors.Errorf("[ DeserializeTwoBigInt ] wrong data length: %d", len(data))
	}

	var one, two big.Int

	one.SetBytes(data[:expectedBigIntBytesLength])
	two.SetBytes(data[expectedBigIntBytesLength:TwoBigIntBytesLength])

	return &one, &two, nil
}

func canonicalizeInt(val *big.Int) ([]byte, error) {
	bytes := val.Bytes()
	size := len(bytes)

	if size > expectedBigIntBytesLength {
		return nil, errors.Errorf("Failed to canonicalize big.Int - wrong length: %d", size)
	}

	paddingSize := expectedBigIntBytesLength - size
	if paddingSize > 0 {
		paddedBytes := make([]byte, size+paddingSize)

		copy(paddedBytes[paddingSize:], bytes)
		return paddedBytes, nil
	}

	return bytes, nil
}

func mustCanonicalizeInt(val *big.Int) []byte {
	bytes, err := canonicalizeInt(val)

	if err != nil {
		panic(err)
	}

	return bytes
}

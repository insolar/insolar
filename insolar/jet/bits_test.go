//
// Copyright 2019 Insolar Technologies GmbH
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
//

package jet

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// bitsToString converts byte slice to plain string representation in bits.
func bitsToString(bitslice []byte) string {
	var b strings.Builder
	for _, bits := range bitslice {
		fmt.Fprintf(&b, "%08b", bits)
	}
	return b.String()
}

// setBitsPrefix copies count bits from []bits to []in.
func setBitsPrefix(in []byte, bits []byte, count int) []byte {
	out := make([]byte, len(in))
	// copy in to out
	_ = copy(out, in)
	if count > len(bits)*8 {
		panic(fmt.Sprintf("count=%v is greater whan bits slice size in bits %v*8=%v",
			count, len(bits), len(bits)*8))
	}

	// copy first bytes from bits, instead last byte if count%8 != 0
	bytesTailOffset := count / 8
	copy(out, bits[:bytesTailOffset])

	bitsTailOffset := count % 8
	if bitsTailOffset == 0 {
		return out
	}

	// preserve last bits in last modified byte in []in
	// overwrite other bits by first bits from []bits

	// switch first bits in last byte of []in to 0
	// [1011 0110], bitsTailOffset=5 -> 0000 0111 (mask: 00000111)
	inMask := byte(0xFF)
	inMask >>= byte(bitsTailOffset)
	inMask &= in[bytesTailOffset]
	// switch last bits in last byte of []bits to 0
	// [0110 1101], bitsTailOffset=5 -> 0110 1000 (mask: 11111000)
	bitsMask := byte(0xFF)
	bitsMask <<= 8 - byte(bitsTailOffset)
	bitsMask &= bits[bytesTailOffset]
	// bits[bytesTailOffset][:bitsTailOffset] + in[bytesTailOffset][bitsTailOffset:]
	out[bytesTailOffset] = inMask | bitsMask

	return out
}

func TestJetStorage_bitsToString(t *testing.T) {
	cases := map[string][]byte{
		"11110000" + "10100111": {0xF0, 0xA7},
		"01001001":              {0x49},
	}
	for str, b := range cases {
		assert.Equal(t, str, bitsToString(b))
	}
}

func TestJetStorage_setBitsPrefix(t *testing.T) {
	// inBits + overwriteBits -> expectedBits
	var (
		inBits        = "01000111" + "11000101"
		overwriteBits = "01100010" + "0110"
		expectedBits  = "01100010" + "01100101"
		// check is setBitsPrefix ignores tail of overwriteBits buffer
		overwriteData = overwriteBits + "0111"
	)

	var (
		in        = parsePrefix(inBits)
		overwrite = parsePrefix(overwriteData)
		expected  = parsePrefix(expectedBits)
	)

	got := setBitsPrefix(in, overwrite, len(overwriteBits))
	assert.Equal(t, expected, got)
}

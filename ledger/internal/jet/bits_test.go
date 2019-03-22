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
	"testing"

	"github.com/magiconair/properties/assert"
)

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

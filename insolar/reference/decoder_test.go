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

package reference

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/network/consensus/common/pulse"
)

func TestDecoder_Decode_legacy(t *testing.T) {
	legacyReference_ok := "1tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.11111111111111111111111111111111"
	{ // good old reference, ok to parse
		dec := NewDecoder(AllowLegacy)
		global, err := dec.Decode(legacyReference_ok)
		if assert.NoError(t, err) {
			assert.Equal(t, global.addressLocal, global.addressBase)
			assert.Equal(t, pulse.Number(0x1000000), global.addressLocal.GetPulseNumber())
			assert.Equal(t, uint8(0x0), global.addressBase.getScope())
		}
	}
	{ // good old reference, disallow parsing
		dec := NewDecoder(0)
		_, err := dec.Decode(legacyReference_ok)
		assert.Error(t, err)
	}

	legacyReference_bad := "1tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.1tJDSkaUggSvNZBDPpPgENR2j3QhzC1wbZS9uyxK2f"
	{ // bad legacy reference (domain isn't empty)
		dec := NewDecoder(AllowLegacy)
		_, err := dec.Decode(legacyReference_bad)
		assert.Error(t, err)
	}

	legacyReference_empty := "11111111111111111111111111111111.11111111111111111111111111111111"
	{ // empty legacy reference
		dec := NewDecoder(AllowLegacy)
		_, err := dec.Decode(legacyReference_empty)
		assert.NoError(t, err)
	}
}

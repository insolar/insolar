/*
 *    Copyright 2018 INS Ecosystem
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

package badgerdb

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/record"
)

func MustDecodeHexString(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func TestStore_prefixkey(t *testing.T) {
	passRecPulse0 := record.LockUnlockRequest{}
	raw, err := record.EncodeToRaw(&passRecPulse0)
	assert.Nil(t, err)
	ref := &record.Reference{
		Domain: record.ID{Pulse: 0, Hash: raw.Hash()},
		Record: record.ID{Pulse: 0, Hash: raw.Hash()},
	}
	key := ref.CoreRef()
	keyP := prefixkey(0, key[:])
	emptyHexStr := strings.Repeat("00", record.IDSize)
	emptyKey := MustDecodeHexString(emptyHexStr + emptyHexStr)
	emptyKeyPrefix := MustDecodeHexString("00" + emptyHexStr + emptyHexStr)

	assert.NotEqual(t, emptyKey, key)
	assert.NotEqual(t, emptyKeyPrefix, keyP)

	expectHexKey := "00000000416ad5cadc41ad8829bdc099b3b20f04dce93217219487fb64cbced600000000416ad5cadc41ad8829bdc099b3b20f04dce93217219487fb64cbced6"
	expectHexKeyP := "00" + expectHexKey
	assert.Equal(t, MustDecodeHexString(expectHexKey), key[:])
	assert.Equal(t, MustDecodeHexString(expectHexKeyP), keyP)
}

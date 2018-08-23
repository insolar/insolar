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

package leveldb

// tests for non public stuff
//
// TODO:
// refactor public API tests to check same things via public APIs and remove internals test.

import (
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

func TestMain(m *testing.M) {
	if err := DropDB(); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func MustDecodeHexString(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
func TestPrefixkey(t *testing.T) {
	passRecPulse0 := record.LockUnlockRequest{}
	raw, err := record.EncodeToRaw(&passRecPulse0)
	assert.Nil(t, err)
	ref := &record.Reference{
		Domain: record.ID{Pulse: 0, Hash: raw.Hash()},
		Record: record.ID{Pulse: 0, Hash: raw.Hash()},
	}
	key := ref.Key()
	keyP := prefixkey(0, key)
	emptyHexStr := strings.Repeat("00", record.IDSize)
	emptyKey := MustDecodeHexString(emptyHexStr + emptyHexStr)
	emptyKeyPrefix := MustDecodeHexString("00" + emptyHexStr + emptyHexStr)

	assert.NotEqual(t, emptyKey, key)
	assert.NotEqual(t, emptyKeyPrefix, keyP)
	// log.Printf("emptyKey:  %x\n", emptyKey)
	// log.Printf("k:         %x\n", k)
	// log.Printf("prefixk: %x\n", kPrefix)

	expectHexKey := "00000000416ad5cadc41ad8829bdc099b3b20f04dce93217219487fb64cbced600000000416ad5cadc41ad8829bdc099b3b20f04dce93217219487fb64cbced6"
	expectHexKeyP := "00" + expectHexKey
	assert.Equal(t, MustDecodeHexString(expectHexKey), key)
	assert.Equal(t, MustDecodeHexString(expectHexKeyP), keyP)
}

func setRawRecord(ll *LevelLedger, ref *record.Reference, raw *record.Raw) error {
	k := prefixkey(scopeIDRecord, ref.Key())
	return ll.ldb.Put(k, record.MustEncodeRaw(raw), nil)
}

func TestSetRawRecord(t *testing.T) {
	ledger, err := InitDB(dbDirPath, nil)
	assert.Nil(t, err)
	defer ledger.Close()

	// prepare record and it's raw representation
	passRecPulse0 := record.LockUnlockRequest{}
	raw, err := record.EncodeToRaw(&passRecPulse0)
	assert.Nil(t, err)
	ref := &record.Reference{
		Domain: record.ID{Pulse: 0, Hash: raw.Hash()},
		Record: record.ID{Pulse: 0, Hash: raw.Hash()},
	}

	// record should not exists
	rec, err := ledger.GetRecord(ref)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, rec)

	// put record in storage by key
	err = setRawRecord(ledger, ref, raw)
	assert.Nil(t, err)

	// get record from storage by key
	gotrec, err := ledger.GetRecord(ref)
	assert.Nil(t, err)
	assert.Equal(t, &passRecPulse0, gotrec)
}

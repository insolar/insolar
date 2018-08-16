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

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

func TestMain(m *testing.M) {
	if err := DropDB(); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func setRawRecord(ll *LevelLedger, ref *record.Reference, raw *record.Raw) error {
	k := prefixkey(scopeIDRecord, ref.Key())
	return ll.ldb.Put(k, record.MustEncodeRaw(raw), nil)
}

func TestGetRecordNotFound(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	ref := &record.Reference{}
	rec, err := ledger.GetRecord(ref)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, rec)
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

func TestSetRawRecord(t *testing.T) {
	ledger, err := InitDB()
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

func zerohash() []byte {
	b := make([]byte, record.HashSize)
	return b
}

func randhash() []byte {
	b := zerohash()
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func hexhash(hash string) []byte {
	b := zerohash()
	if len(hash)%2 == 1 {
		hash = "0" + hash
	}
	h, err := hex.DecodeString(hash)
	if err != nil {
		panic(err)
	}
	_ = copy(b, h)
	return b
}

func referenceWithHashes(domainhash, recordhash string) record.Reference {
	dh := hexhash(domainhash)
	rh := hexhash(recordhash)

	return record.Reference{
		Domain: record.ID{Hash: dh},
		Record: record.ID{Hash: rh},
	}
}

func TestSetRecord(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	// mock pulse source
	pulse1 := record.PulseNum(1)
	ledger.pulseFn = func() record.PulseNum { return pulse1 }
	defer ledger.Close()

	passRecPulse1 := &record.LockUnlockRequest{}
	idPulse1 := pulse1.ID(passRecPulse1)

	refPulse1 := &record.Reference{
		Domain: record.ID{},
		Record: idPulse1,
	}
	rec, err := ledger.GetRecord(refPulse1)
	assert.Nil(t, rec)
	assert.Equal(t, storage.ErrNotFound, err)

	gotRef, err := ledger.SetRecord(passRecPulse1)
	assert.Nil(t, err)
	assert.Equal(t, idPulse1, gotRef.Record)
	assert.Equal(t, refPulse1, gotRef)

	gotRec, err := ledger.GetRecord(gotRef)
	assert.Nil(t, err)
	assert.Equal(t, passRecPulse1, gotRec)

	// check is record IDs in different pulses are not the same
	pulse0 := record.PulseNum(0)
	idPulse0 := pulse0.ID(gotRec)

	idPulse0Hex := fmt.Sprintf("%x", idPulse0)
	idPulse1Hex := fmt.Sprintf("%x", idPulse1)
	assert.NotEqual(t, idPulse1Hex, idPulse0Hex, "got hash")
}

// TODO: uncomment when record storage is functional
// func TestCreatesRootRecord(t *testing.T) {
// 	ledger, err := InitDB()
// 	assert.Nil(t, err)
// 	defer ledger.Close()
//
// 	var zeroID record.ID
// 	copy([]byte(zeroRecordBinary)[:record.IDSize], zeroID[:])
// 	zeroRef, ok := ledger.GetRecord(record.ID2Key(zeroID))
// 	assert.True(t, ok)
// 	assert.Equal(t, ledger.zeroRef, zeroRef)
// }

func TestGetClassIndexOnEmptyDataReturnsNotFound(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	ref := &record.Reference{
		Record: record.ID{Pulse: 1},
	}

	idx, err := ledger.GetClassIndex(ref)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, idx)
}

func TestSetClassIndexStoresDataInDB(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	zerodomain := record.ID{Hash: zerohash()}
	refgen := func() record.Reference {
		recID := record.ID{
			Hash: randhash(),
		}
		return record.Reference{
			Domain: zerodomain,
			Record: recID,
		}
	}
	latestRef := refgen()
	idx := index.ClassLifeline{
		LatestStateRef: latestRef,
		AmendRefs:      []record.Reference{refgen(), refgen(), refgen()},
	}
	zeroRef := record.Reference{
		Domain: zerodomain,
		Record: record.ID{
			Hash: hexhash("122444"),
		},
	}
	err = ledger.SetClassIndex(&zeroRef, &idx)
	assert.Nil(t, err)

	storedIndex, err := ledger.GetClassIndex(&zeroRef)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

func TestGetObjectIndexOnEmptyDataReturnsNotFound(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	ref := referenceWithHashes("1000", "5000")
	idx, err := ledger.GetObjectIndex(&ref)
	assert.Equal(t, storage.ErrNotFound, err)
	assert.Nil(t, idx)
}

func TestSetObjectIndexStoresDataInDB(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	idx := index.ObjectLifeline{
		ClassRef:       referenceWithHashes("50", "60"),
		LatestStateRef: referenceWithHashes("10", "20"),
		AppendRefs: []record.Reference{
			referenceWithHashes("", "1"),
			referenceWithHashes("", "2"),
			referenceWithHashes("", "3"),
		},
	}
	zeroref := referenceWithHashes("", "")
	err = ledger.SetObjectIndex(&zeroref, &idx)
	assert.Nil(t, err)

	storedIndex, err := ledger.GetObjectIndex(&zeroref)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

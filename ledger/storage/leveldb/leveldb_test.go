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

package leveldb_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/leveldb/leveltestutils"
)

func TestLevelLedger_GetRecordNotFound(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	ref := &record.Reference{}
	rec, err := ledger.GetRecord(ref)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, rec)
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

func TestLevelLedger_SetRecord(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()
	// mock pulse source
	pulse1 := record.PulseNum(1)
	ledger.SetCurrentPulse(pulse1)

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

func TestLevelLedger_GetClassIndex_ReturnsNotFoundIfNoIndex(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	ref := &record.Reference{
		Record: record.ID{Pulse: 1},
	}

	idx, err := ledger.GetClassIndex(ref)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, idx)
}

func TestLevelLedger_SetClassIndex_StoresCorrectDataInStorage(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

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
	err := ledger.SetClassIndex(&zeroRef, &idx)
	assert.Nil(t, err)

	storedIndex, err := ledger.GetClassIndex(&zeroRef)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

func TestLevelLedger_SetObjectIndex_ReturnsNotFoundIfNoIndex(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	ref := referenceWithHashes("1000", "5000")
	idx, err := ledger.GetObjectIndex(&ref)
	assert.Equal(t, storage.ErrNotFound, err)
	assert.Nil(t, idx)
}

func TestLevelLedger_SetObjectIndex_StoresCorrectDataInStorage(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

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
	err := ledger.SetObjectIndex(&zeroref, &idx)
	assert.Nil(t, err)

	storedIndex, err := ledger.GetObjectIndex(&zeroref)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

func TestLevelLedger_GetDrop_ReturnsNotFoundIfNoDrop(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	drop, err := ledger.GetDrop(1)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, drop)
}

func TestLevelLedger_SetDrop_StoresCorrectDataInStorage(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	// it references on 'fake' zero
	fakeDrop := jetdrop.JetDrop{
		Hash: []byte{0xFF},
	}

	ledger.SetCurrentPulse(42)
	drop42, err := ledger.SetDrop(42, &fakeDrop)
	assert.NoError(t, err)
	got, err := ledger.GetDrop(42)
	assert.NoError(t, err)
	assert.Equal(t, got, drop42)
}

func TestLevelLedger_SetCurrentPulse(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	ledger.SetCurrentPulse(42)
	assert.Equal(t, record.PulseNum(42), ledger.GetCurrentPulse())
}

func TestLevelLedger_SetEntropy(t *testing.T) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")
	defer cleaner()

	ledger.SetEntropy(42, []byte{1, 2, 3})
	entropy, err := ledger.GetEntropy(42)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, entropy)
	entropy, err = ledger.GetEntropy(1)
	assert.Error(t, err)
}

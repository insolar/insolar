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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
)

func TestMain(m *testing.M) {
	absPath, err := filepath.Abs(dbDirPath)
	if err != nil {
		os.Exit(1)
	}

	if err = os.RemoveAll(absPath); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestGetRecordNotFound(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	rec, err := ledger.GetRecordByKey(record.Key{Hash: []byte("NotFoundRecord")})
	assert.Equal(t, err, ErrNotFound)
	assert.Nil(t, rec)
}

func TestSetRawRecord(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	// prepare record and it's raw representation
	var passRecPulse0 record.LockUnlockRequest
	raw, err := record.EncodeToRaw(&passRecPulse0)
	assert.Nil(t, err)
	key := record.Key{
		Pulse: 0,
		Hash:  raw.Hash(),
	}

	// record should not exists
	rec, err := ledger.GetRecordByKey(key)
	assert.Equal(t, err, ErrNotFound)
	assert.Nil(t, rec)

	// put record in storage by key
	id, err := ledger.setRawRecordByKey(key, raw)
	assert.Nil(t, err)
	// fmt.Printf("saved by id %x\n", id)

	// get record from storage by key
	gotrec, err := ledger.GetRecordByKey(key)
	assert.Nil(t, err)
	assert.Equal(t, &passRecPulse0, gotrec)

	// get record from storage by id
	gotrec, err = ledger.GetRecord(id)
	assert.Nil(t, err)
	assert.Equal(t, &passRecPulse0, gotrec)
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
	rec, err := ledger.GetRecord(idPulse1)
	assert.Nil(t, rec)
	assert.Equal(t, ErrNotFound, err)

	gotid, err := ledger.SetRecord(passRecPulse1)
	assert.Nil(t, err)
	assert.Equal(t, idPulse1, gotid)

	gotrec, err := ledger.GetRecord(gotid)
	assert.Nil(t, err)
	assert.Equal(t, passRecPulse1, gotrec)

	// check is record IDs in different pulses are not the same
	pulse0 := record.PulseNum(0)
	idPulse0 := pulse0.ID(gotrec)

	idPulse0Hex := fmt.Sprintf("%x", idPulse0)
	idPulse1Hex := fmt.Sprintf("%x", idPulse1)
	assert.NotEqual(t, idPulse1Hex, idPulse0Hex, "got hash")

}

func TestGetIndexOnEmptyDataReturnsNotFound(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	idx, isFound := ledger.GetIndex(record.ID{1})
	assert.Equal(t, isFound, false)
	assert.Nil(t, idx)
}

func TestSetIndexStoresDataInDB(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	idx := index.Lifeline{
		LatestStateID: record.ID{1, 2, 3},
		AppendIDs:     []record.ID{{1}, {2}, {3}},
	}
	err = ledger.SetIndex(record.ID{0}, &idx)
	assert.Nil(t, err)

	storedIndex, isFound := ledger.GetIndex(record.ID{0})
	assert.Equal(t, isFound, true)
	assert.Equal(t, *storedIndex, idx)
}

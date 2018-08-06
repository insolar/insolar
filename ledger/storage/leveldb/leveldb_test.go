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
	"os"
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/stretchr/testify/assert"
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

	idx, isFound := ledger.GetClassIndex(record.ID{1})
	assert.Equal(t, isFound, false)
	assert.Nil(t, idx)
}

func TestSetClassIndexStoresDataInDB(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	idx := index.ClassLifeline{
		LatestStateID: record.ID{1, 2, 3},
		MigrationIDs:  []record.ID{{1}, {2}, {3}},
	}
	err = ledger.SetClassIndex(record.ID{0}, &idx)
	assert.Nil(t, err)

	storedIndex, isFound := ledger.GetClassIndex(record.ID{0})
	assert.Equal(t, isFound, true)
	assert.Equal(t, *storedIndex, idx)
}

func TestGetObjectIndexOnEmptyDataReturnsNotFound(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	idx, isFound := ledger.GetObjectIndex(record.ID{1})
	assert.Equal(t, isFound, false)
	assert.Nil(t, idx)
}

func TestSetObjectIndexStoresDataInDB(t *testing.T) {
	ledger, err := InitDB()
	assert.Nil(t, err)
	defer ledger.Close()

	idx := index.ObjectLifeline{
		ClassID:       record.ID{5, 6},
		LatestStateID: record.ID{1, 2, 3},
		AppendIDs:     []record.ID{{1}, {2}, {3}},
	}
	err = ledger.SetObjectIndex(record.ID{0}, &idx)
	assert.Nil(t, err)

	storedIndex, isFound := ledger.GetObjectIndex(record.ID{0})
	assert.Equal(t, isFound, true)
	assert.Equal(t, *storedIndex, idx)
}

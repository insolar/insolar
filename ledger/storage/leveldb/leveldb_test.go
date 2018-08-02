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
	if err == nil {
		os.RemoveAll(absPath)
	}
	os.Exit(m.Run())
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
		LatestStateID:   record.ID{1, 2, 3},
		LatestStateType: 1,
		AppendIDs:       []record.ID{{1}, {2}, {3}},
	}
	err = ledger.SetIndex(record.ID{0}, idx)
	assert.Nil(t, err)

	storedIndex, isFound := ledger.GetIndex(record.ID{0})
	assert.Equal(t, isFound, true)
	assert.Equal(t, *storedIndex, idx)
}

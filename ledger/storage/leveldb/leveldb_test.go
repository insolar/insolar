package leveldb

import (
	"os"
	"testing"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/insolar/insolar/ledger/record"
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
	assert.Equal(t, err, nil)
	defer ledger.Close()

	index, isFound := ledger.GetIndex(record.ID{1})
	assert.Equal(t, isFound, false)
	assert.Equal(t, index, record.LifelineIndex{})
}

func TestSetIndexStoresDataInDB(t *testing.T) {
	ledger, err := InitDB()
	assert.Equal(t, err, nil)
	defer ledger.Close()

	index := record.LifelineIndex{
		LatestStateID: record.ID{1, 2, 3},
		LatestStateType: 1,
		AppendIDs: []record.ID{{1}, {2}, {3}},
	}
	err = ledger.SetIndex(record.ID{0}, index)
	assert.Equal(t, err, nil)

	storedIndex, isFound := ledger.GetIndex(record.ID{0})
	assert.Equal(t, isFound, true)
	assert.Equal(t, storedIndex, index)
}

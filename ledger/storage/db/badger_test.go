package db

import (
	"io/ioutil"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/google/gofuzz"
	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testKey struct {
	id    []byte
	scope Scope
}

func (k testKey) Scope() Scope {
	return k.scope
}

func (k testKey) ID() []byte {
	return k.id
}

func TestBadgerDB_Get(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	assert.NoError(t, err)

	db, err := NewBadgerDB(configuration.Ledger{Storage: configuration.Storage{DataDirectory: tmpdir}})
	require.NoError(t, err)

	var (
		key           testKey
		expectedValue []byte
	)
	f := fuzz.New().NilChance(0)
	f.Fuzz(&key)
	f.Fuzz(&expectedValue)
	err = db.backend.Update(func(txn *badger.Txn) error {
		return txn.Set(append(key.Scope().Bytes(), key.ID()...), expectedValue)
	})
	require.NoError(t, err)
	value, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

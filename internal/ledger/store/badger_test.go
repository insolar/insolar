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

package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testBadgerKey struct {
	id    []byte
	scope Scope
}

func (k testBadgerKey) Scope() Scope {
	return k.scope
}

func (k testBadgerKey) ID() []byte {
	return k.id
}

func TestBadgerDB_Get(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	db, err := NewBadgerDB(tmpdir)
	require.NoError(t, err)

	var (
		key           testBadgerKey
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

func TestBadgerDB_Set(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	db, err := NewBadgerDB(tmpdir)
	require.NoError(t, err)

	var (
		key           testBadgerKey
		expectedValue []byte
		value         []byte
	)
	f := fuzz.New().NilChance(0)
	f.Fuzz(&key)
	f.Fuzz(&expectedValue)
	err = db.Set(key, expectedValue)
	assert.NoError(t, err)

	err = db.backend.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(key.Scope().Bytes(), key.ID()...))
		require.NoError(t, err)
		value, err = item.ValueCopy(nil)
		require.NoError(t, err)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestBadgerDB_NewIterator(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	db, err := NewBadgerDB(tmpdir)
	require.NoError(t, err)

	var (
		commonScope      Scope
		commonPrefix     []byte
		randKeys         []testBadgerKey
		unexpectedValues [][]byte
		keys             []testBadgerKey
		expectedValues   [][]byte
		expected         map[string][]byte
	)

	expected = make(map[string][]byte)

	const (
		ArrayLength = 3
	)

	fuzz.New().NilChance(0).Fuzz(&commonScope)
	fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Fuzz(&commonPrefix)

	ff := fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Funcs(
		func(key *testBadgerKey, c fuzz.Continue) {
			var id []byte
			c.Fuzz(&id)
			*key = testBadgerKey{id: id, scope: commonScope}
		},
	)
	ff.Fuzz(&randKeys)
	ff.Fuzz(&unexpectedValues)

	f := fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Funcs(
		func(key *testBadgerKey, c fuzz.Continue) {
			var id []byte
			c.Fuzz(&id)
			*key = testBadgerKey{id: append(commonPrefix, id...), scope: commonScope}
		},
	)
	f.Fuzz(&keys)
	f.Fuzz(&expectedValues)

	err = db.backend.Update(func(txn *badger.Txn) error {
		for i := 0; i < ArrayLength; i++ {
			key := append(randKeys[i].Scope().Bytes(), randKeys[i].ID()...)
			err = txn.Set(key, unexpectedValues[i])
			if err != nil {
				return err
			}
		}
		for i := 0; i < ArrayLength; i++ {
			key := append(keys[i].Scope().Bytes(), keys[i].ID()...)
			err = txn.Set(key, expectedValues[i])
			if err != nil {
				return err
			}
			expected[string(keys[i].ID())] = expectedValues[i]
		}
		return nil
	})
	require.NoError(t, err)

	// test logic
	it := db.NewIterator(commonScope)
	defer it.Close()
	it.Seek(commonPrefix)
	for it.Next() {
		require.ElementsMatch(t, expected[string(it.Key())], it.Value())
	}
}

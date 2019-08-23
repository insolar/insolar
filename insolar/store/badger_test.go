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
	"bytes"
	"crypto/rand"
	"io/ioutil"
	rand2 "math/rand"
	"os"
	"sort"
	"testing"

	"github.com/dgraph-io/badger"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

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

	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := NewBadgerDB(ops)
	defer db.Stop(ctx)
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

	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := NewBadgerDB(ops)
	defer db.Stop(ctx)
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

func TestBadgerDB_Delete(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := NewBadgerDB(ops)
	defer db.Stop(ctx)
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

	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err = db.backend.View(func(txn *badger.Txn) error {
		item, err := txn.Get(fullKey)
		require.NoError(t, err)
		value, err = item.ValueCopy(nil)
		require.NoError(t, err)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)

	err = db.backend.Update(func(txn *badger.Txn) error {
		err := txn.Delete(fullKey)
		require.NoError(t, err)
		return nil
	})
	assert.NoError(t, err)

	err = db.backend.View(func(txn *badger.Txn) error {
		_, err := txn.Get(fullKey)
		require.EqualError(t, err, badger.ErrKeyNotFound.Error())
		return nil
	})

}

func TestBadgerDB_NewIterator(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := NewBadgerDB(ops)
	defer db.Stop(ctx)
	require.NoError(t, err)

	type kv struct {
		k testBadgerKey
		v []byte
	}

	var (
		commonScope  Scope
		commonPrefix []byte

		expected   []kv
		unexpected []kv
	)

	const (
		ArrayLength = 100
	)

	fuzz.New().NilChance(0).Fuzz(&commonScope)
	fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Fuzz(&commonPrefix)

	f := fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Funcs(
		func(key *testBadgerKey, c fuzz.Continue) {
			var id []byte
			c.Fuzz(&id)
			key.id = append(commonPrefix, id...)
			key.id[0] = commonPrefix[0] + 1
			key.scope = commonScope
		},
		func(pair *kv, c fuzz.Continue) {
			c.Fuzz(&pair.k)
			c.Fuzz(&pair.v)
		},
	)
	f.Fuzz(&unexpected)

	f = fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Funcs(
		func(key *testBadgerKey, c fuzz.Continue) {
			var id []byte
			c.Fuzz(&id)
			key.id = append(commonPrefix, id...)
			key.scope = commonScope
		},
		func(pair *kv, c fuzz.Continue) {
			c.Fuzz(&pair.k)
			c.Fuzz(&pair.v)
		},
	)
	f.Fuzz(&expected)

	sort.Slice(expected, func(i, j int) bool {
		return bytes.Compare(expected[i].k.ID(), expected[j].k.ID()) == -1
	})

	err = db.backend.Update(func(txn *badger.Txn) error {
		for i := 0; i < ArrayLength; i++ {
			key := append(unexpected[i].k.Scope().Bytes(), unexpected[i].k.ID()...)
			err = txn.Set(key, unexpected[i].v)
			if err != nil {
				return err
			}
		}
		for i := 0; i < ArrayLength; i++ {
			key := append(expected[i].k.Scope().Bytes(), expected[i].k.ID()...)
			err = txn.Set(key, expected[i].v)
			if err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	// test logic
	pivot := testBadgerKey{id: commonPrefix, scope: commonScope}
	it := db.NewIterator(pivot, false)
	defer it.Close()
	i := 0
	for it.Next() && i < len(expected) {
		require.Equal(t, expected[i].k.ID(), it.Key())
		val, err := it.Value()
		require.NoError(t, err)
		require.Equal(t, expected[i].v, val)
		i++
	}
	require.Equal(t, len(expected), i)
}

func TestBadgerDB_NewReverseIterator(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := NewBadgerDB(ops)
	defer db.Stop(ctx)
	require.NoError(t, err)

	type kv struct {
		k testBadgerKey
		v []byte
	}

	var (
		commonScope  Scope
		commonPrefix []byte

		expected   []kv
		unexpected []kv
	)

	const (
		ArrayLength  = 100
		ReverseOrder = true
	)

	fuzz.New().NilChance(0).Fuzz(&commonScope)
	fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Fuzz(&commonPrefix)

	f := fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Funcs(
		func(key *testBadgerKey, c fuzz.Continue) {
			var id []byte
			c.Fuzz(&id)
			key.id = append(commonPrefix, id...)
			key.id[0] = commonPrefix[0] + 1
			key.scope = commonScope
		},
		func(pair *kv, c fuzz.Continue) {
			c.Fuzz(&pair.k)
			c.Fuzz(&pair.v)
		},
	)
	f.Fuzz(&unexpected)

	f = fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Funcs(
		func(key *testBadgerKey, c fuzz.Continue) {
			var id []byte
			c.Fuzz(&id)
			key.id = append(commonPrefix, id...)
			key.scope = commonScope
		},
		func(pair *kv, c fuzz.Continue) {
			c.Fuzz(&pair.k)
			c.Fuzz(&pair.v)
		},
	)
	f.Fuzz(&expected)

	sort.Slice(expected, func(i, j int) bool {
		return bytes.Compare(expected[i].k.ID(), expected[j].k.ID()) == -1
	})

	err = db.backend.Update(func(txn *badger.Txn) error {
		for i := 0; i < ArrayLength; i++ {
			key := append(unexpected[i].k.Scope().Bytes(), unexpected[i].k.ID()...)
			err = txn.Set(key, unexpected[i].v)
			if err != nil {
				return err
			}
		}
		for i := 0; i < ArrayLength; i++ {
			key := append(expected[i].k.Scope().Bytes(), expected[i].k.ID()...)
			err = txn.Set(key, expected[i].v)
			if err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	// test logic
	prefix := fillPrefix(commonPrefix, ArrayLength*2)
	pivot := testBadgerKey{id: prefix, scope: commonScope}
	it := db.NewIterator(pivot, ReverseOrder)
	defer it.Close()
	i := 0
	for it.Next() && i < len(expected) {
		require.Equal(t, expected[len(expected)-i-1].k.ID(), it.Key())
		val, err := it.Value()
		require.NoError(t, err)
		require.Equal(t, expected[len(expected)-i-1].v, val)
		i++
	}
	require.Equal(t, len(expected), i)
}

func TestBadgerDB_SimpleReverse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := NewBadgerDB(ops)
	defer db.Stop(ctx)
	require.NoError(t, err)

	count := 100
	length := 10
	prefixes := make([][]byte, count)
	keys := make([][]byte, count)
	for i := 0; i < count; i++ {
		prefixes[i] = make([]byte, length)
		keys[i] = make([]byte, length)
		_, err = rand.Read(prefixes[i])
		require.NoError(t, err)
		_, err = rand.Read(keys[i])
		require.NoError(t, err)
		keys[i][0] = 0xFF
		keys[i] = append(prefixes[i], keys[i]...)
		err = db.Set(testBadgerKey{keys[i], ScopeRecord}, nil)
		require.NoError(t, err)
	}

	t.Run("ASC iteration", func(t *testing.T) {
		asc := make([][]byte, count)
		copy(asc, keys)
		sort.Slice(keys, func(i, j int) bool {
			return bytes.Compare(keys[i], keys[j]) == -1
		})
		sort.Slice(prefixes, func(i, j int) bool {
			return bytes.Compare(prefixes[i], prefixes[j]) == -1
		})

		seek := rand2.Intn(count)
		pivot := testBadgerKey{id: prefixes[seek], scope: ScopeRecord}
		it := db.NewIterator(pivot, false)
		defer it.Close()
		var actual [][]byte
		for it.Next() {
			actual = append(actual, it.Key())
		}
		require.Equal(t, count-seek, len(actual))
		require.Equal(t, keys[seek:], actual)
	})

	t.Run("DESC iteration", func(t *testing.T) {
		desc := make([][]byte, count)
		copy(desc, keys)
		sort.Slice(keys, func(i, j int) bool {
			return bytes.Compare(keys[i], keys[j]) >= 0
		})
		sort.Slice(prefixes, func(i, j int) bool {
			return bytes.Compare(prefixes[i], prefixes[j]) >= 0
		})

		seek := rand2.Intn(count)
		prefix := fillPrefix(prefixes[seek], length*2)
		pivot := testBadgerKey{id: prefix, scope: ScopeRecord}
		it := db.NewIterator(pivot, true)
		defer it.Close()
		var actual [][]byte
		for it.Next() {
			actual = append(actual, it.Key())
		}
		require.Equal(t, count-seek, len(actual))
		require.Equal(t, keys[seek:], actual)
	})
}

func fillPrefix(prefix []byte, keyLen int) []byte {
	rest := keyLen - len(prefix)
	filler := make([]byte, rest)
	for i := range filler {
		filler[i] = 0xFF
	}
	return bytes.Join([][]byte{prefix, filler}, nil)
}

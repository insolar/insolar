//
// Copyright (c) 2019. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
// Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
// Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
// Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
// Vestibulum commodo. Ut rhoncus gravida arcu.
//

package store

import (
	"bytes"
	"sort"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testMockKey struct {
	id    []byte
	scope Scope
}

func (k testMockKey) Scope() Scope {
	return k.scope
}

func (k testMockKey) ID() []byte {
	return k.id
}

func TestMockDB_Get(t *testing.T) {
	t.Parallel()

	db := NewMemoryMockDB()

	var (
		key           testMockKey
		expectedValue []byte
	)
	f := fuzz.New().NilChance(0)
	f.Fuzz(&key)
	f.Fuzz(&expectedValue)
	db.backend[string(append(key.Scope().Bytes(), key.ID()...))] = expectedValue
	value, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestNewMockDB_Get_ValueImmutable(t *testing.T) {
	db := NewMemoryMockDB()
	key := testMockKey{
		id:    []byte{1},
		scope: 0,
	}
	db.backend[string(append(key.Scope().Bytes(), key.ID()...))] = []byte{1, 2, 3}
	value, err := db.Get(key)
	assert.NoError(t, err)
	value[0] = 42
	sameValue, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, sameValue)
}

func TestMockDB_Set(t *testing.T) {
	t.Parallel()

	db := NewMemoryMockDB()

	var (
		key           testMockKey
		expectedValue []byte
	)
	f := fuzz.New().NilChance(0)
	f.Fuzz(&key)
	f.Fuzz(&expectedValue)
	err := db.Set(key, expectedValue)
	assert.NoError(t, err)

	value := db.backend[string(append(key.Scope().Bytes(), key.ID()...))]
	assert.Equal(t, expectedValue, value)
}

func TestMockDB_NewIterator(t *testing.T) {
	t.Parallel()

	db := NewMemoryMockDB()

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
		ArrayLength = 1000
	)

	fuzz.New().NilChance(0).Fuzz(&commonScope)
	fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Fuzz(&commonPrefix)

	f := fuzz.New().NilChance(0).NumElements(ArrayLength, ArrayLength).Funcs(
		func(key *testBadgerKey, c fuzz.Continue) {
			for {
				c.Fuzz(&key.id)
				// To ensure that unexpected keys will be started with prefix that less than expected keys
				if bytes.Compare(key.id, commonPrefix) == -1 {
					break
				}
			}
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

	err := error(nil)
	for i := 0; i < ArrayLength; i++ {
		err = db.Set(unexpected[i].k, unexpected[i].v)
		if err != nil {
			break
		}
	}
	for i := 0; i < ArrayLength; i++ {
		err = db.Set(expected[i].k, expected[i].v)
		if err != nil {
			break
		}
	}

	require.NoError(t, err)

	// test logic
	it := db.NewIterator(commonScope)
	defer it.Close()
	it.Seek(commonPrefix)
	i := 0
	for it.Next() && i < len(expected) {
		require.ElementsMatch(t, expected[i].k.ID(), it.Key())
		val, err := it.Value()
		require.NoError(t, err)
		require.ElementsMatch(t, expected[i].v, val)
		i++
	}
	require.Equal(t, len(expected), i)
}

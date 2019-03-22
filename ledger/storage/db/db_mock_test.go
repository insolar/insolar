/*
 *    Copyright 2019 Insolar Technologies
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

package db

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
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

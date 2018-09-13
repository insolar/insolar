/*
 *    Copyright 2018 Insolar
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

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMapStorage(t *testing.T) {
	mapStorage := NewMapStorage()

	assert.Equal(t, &MapStorage{
		storage: make(map[string]interface{}),
		keys:    []string{},
	}, mapStorage)
}

func TestMapStorage_Set(t *testing.T) {
	value := "some value"
	mapStorage := NewMapStorage()

	key, err := mapStorage.Set(value)

	assert.NoError(t, err)
	assert.Len(t, mapStorage.storage, 1)
	assert.Equal(t, value, mapStorage.storage[key])
	assert.Equal(t, []string{key}, mapStorage.keys)
}

func TestMapStorage_Get(t *testing.T) {
	value := "some value"
	mapStorage := NewMapStorage()
	key, _ := mapStorage.Set(value)

	storedValue, err := mapStorage.Get(key)

	assert.NoError(t, err)
	assert.Equal(t, value, storedValue)
}

func TestMapStorage_Get_Error(t *testing.T) {
	key := "some key"
	mapStorage := NewMapStorage()

	storedValue, err := mapStorage.Get(key)

	assert.EqualError(t, err, "object with record some key does not exist")
	assert.Nil(t, storedValue)
}

func TestMapStorage_GetKeys(t *testing.T) {
	value := "some value"
	mapStorage := NewMapStorage()
	keyFirst, _ := mapStorage.Set(value)
	keySecond, _ := mapStorage.Set(value)

	keys := mapStorage.GetKeys()

	assert.Equal(t, []string{keyFirst, keySecond}, keys)
}

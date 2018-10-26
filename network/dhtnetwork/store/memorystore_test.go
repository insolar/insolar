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

package store

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryStore(t *testing.T) {
	s := NewMemoryStore()

	assert.Equal(t, s, &memoryStore{
		mutex:        &sync.RWMutex{},
		data:         make(map[string][]byte),
		replicateMap: make(map[string]time.Time),
		expireMap:    make(map[string]time.Time),
	})
}

func TestMemoryStore_Store(t *testing.T) {
	s := newMemoryStore()

	data := []byte("some data")
	key := NewKey(data)

	replicationTime := time.Now().Add(time.Second * 1337)
	expirationTime := time.Now().Add(time.Second * 42)

	s.Store(key, data, replicationTime, expirationTime, true)

	assert.Len(t, s.data, 1)
	assert.Equal(t, s.data[key.String()], data)

	assert.Len(t, s.replicateMap, 1)
	assert.Equal(t, s.replicateMap[key.String()], replicationTime)

	assert.Len(t, s.expireMap, 1)
	assert.Equal(t, s.expireMap[key.String()], expirationTime)
}

func TestMemoryStore_Retrieve(t *testing.T) {
	s := NewMemoryStore()

	data := []byte("some data")
	key := NewKey(data)

	res, found := s.Retrieve(key)
	assert.Nil(t, res)
	assert.False(t, found)

	s.Store(key, data, time.Now(), time.Now(), true)

	res, found = s.Retrieve(key)
	assert.Equal(t, res, data)
	assert.True(t, found)
}

func TestMemoryStore_Delete(t *testing.T) {
	s := newMemoryStore()

	data := []byte("some data")
	key := NewKey(data)

	s.Store(key, data, time.Now(), time.Now(), true)

	s.Delete(key)

	assert.Len(t, s.data, 0)
	assert.Len(t, s.replicateMap, 0)
	assert.Len(t, s.expireMap, 0)
}

func TestMemoryStore_GetKeysReadyToReplicate(t *testing.T) {
	s := newMemoryStore()

	now := time.Now()

	data1 := []byte("some data1")
	key1 := NewKey(data1)
	s.Store(key1, data1, now.Add(-20*time.Second), now, true)

	data2 := []byte("some data2")
	key2 := NewKey(data2)
	s.Store(key2, data2, now.Add(-1*time.Nanosecond), now, true)

	data3 := []byte("some data3")
	key3 := NewKey(data3)
	s.Store(key3, data3, now.Add(10*time.Minute), now, true)

	data4 := []byte("some data4")
	key4 := NewKey(data4)
	s.Store(key4, data4, now.Add(20*time.Second), now, true)

	keys := s.GetKeysReadyToReplicate()

	assert.Len(t, keys, 2)
	assert.Contains(t, keys, key1)
	assert.Contains(t, keys, key2)
}

func TestMemoryStore_ExpireKeys(t *testing.T) {
	s := newMemoryStore()

	now := time.Now()

	data1 := []byte("some data1")
	key1 := NewKey(data1)
	s.Store(key1, data1, now, now.Add(-20*time.Second), true)

	data2 := []byte("some data2")
	key2 := NewKey(data2)
	s.Store(key2, data2, now, now.Add(-1*time.Nanosecond), true)

	data3 := []byte("some data3")
	key3 := NewKey(data3)
	s.Store(key3, data3, now, now.Add(10*time.Minute), true)

	data4 := []byte("some data4")
	key4 := NewKey(data4)
	s.Store(key4, data4, now, now.Add(20*time.Second), true)

	s.ExpireKeys()

	assert.Len(t, s.data, 2)
	assert.Len(t, s.replicateMap, 2)
	assert.Len(t, s.expireMap, 2)

	assert.Equal(t, s.data, map[string][]byte{
		key3.String(): data3,
		key4.String(): data4,
	})
}

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

package routing

import (
	"testing"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/stretchr/testify/assert"
)

func TestNewHashTable(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(id.GetRandomKey())
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
}

func TestHashTable_Lock(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(id.GetRandomKey())
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
	ht.Lock()
	ht.Unlock()
}

func TestHashTable_Unlock(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(id.GetRandomKey())
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
	ht.Lock()
	ht.Unlock()
}

func TestHashTable_ResetRefreshTimeForBucket(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(id.GetRandomKey())
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)

	time1 := ht.GetRefreshTimeForBucket(0)
	time.Sleep(time.Millisecond * 2000)
	ht.ResetRefreshTimeForBucket(0)
	time2 := ht.GetRefreshTimeForBucket(0)
	assert.NotEqual(t, time1, time2)
}

func TestHashTable_TotalHosts(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(id.GetRandomKey())
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)

	assert.Equal(t, 0, ht.TotalHosts())
}

func TestHashTable_GetRandomIDFromBucket(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(id.GetRandomKey())
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
	rand := ht.GetRandomIDFromBucket(0)
	assert.Equal(t, 20, len(rand))
}

func TestHashTable_GetTotalHostsInBucket(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash(id.GetRandomKey())
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)

	assert.Equal(t, 0, ht.GetTotalHostsInBucket(0))
}

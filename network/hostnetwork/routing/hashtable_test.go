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

package routing

import (
	"strconv"
	"testing"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/stretchr/testify/assert"
)

func TestNewHashTable(t *testing.T) {
	id1, _ := id.NewID()
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
}

func TestHashTable_Lock(t *testing.T) {
	id1, _ := id.NewID()
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
	ht.Lock()
	ht.Unlock()
}

func TestHashTable_Unlock(t *testing.T) {
	id1, _ := id.NewID()
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
	ht.Lock()
	ht.Unlock()
}

func TestHashTable_ResetRefreshTimeForBucket(t *testing.T) {
	id1, _ := id.NewID()
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
	id1, _ := id.NewID()
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)

	assert.Equal(t, 0, ht.TotalHosts())
}

func TestHashTable_GetRandomIDFromBucket(t *testing.T) {
	id1, _ := id.NewID()
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)
	rand := ht.GetRandomIDFromBucket(0)
	assert.Equal(t, 20, len(rand))
}

func TestHashTable_GetTotalHostsInBucket(t *testing.T) {
	id1, _ := id.NewID()
	addr, _ := host.NewAddress("127.0.0.1:3000")
	ht, err := NewHashTable(id1, addr)
	assert.NoError(t, err)
	assert.NotNil(t, ht)

	assert.Equal(t, 0, ht.GetTotalHostsInBucket(0))
}

func TestHashTable_GetHosts(t *testing.T) {
	prefix := "127.0.0.1:"
	port := 3000

	origin, _ := id.NewID()
	addr, _ := host.NewAddress(prefix + strconv.Itoa(port))
	ht, _ := NewHashTable(origin, addr)
	port++

	for i := 0; i < 10; i++ {
		id1, _ := id.NewID()
		addr, _ = host.NewAddress(prefix + strconv.Itoa(port))
		h := NewRouteHost(&host.Host{ID: id1, Address: addr})

		index := GetBucketIndexFromDifferingBit(origin, id1)
		bucket := ht.RoutingTable[index]
		bucket = append(bucket, h)
		ht.RoutingTable[index] = bucket
	}

	assert.Equal(t, 10, len(ht.GetHosts(100)))
	assert.Equal(t, 4, len(ht.GetHosts(4)))
}

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

	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/id"
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

func TestHashTable_MarkHostAsSeen(t *testing.T) {
	prefix := "127.0.0.1:"
	port := 3000

	origin, _ := id.NewID()
	addr, _ := host.NewAddress(prefix + strconv.Itoa(port))
	ht, _ := NewHashTable(origin, addr)
	port++

	id1, _ := id.NewID()
	addr, _ = host.NewAddress(prefix + strconv.Itoa(port))
	h := NewRouteHost(&host.Host{ID: id1, Address: addr})

	index := GetBucketIndexFromDifferingBit(origin, id1)
	bucket := ht.RoutingTable[index]
	bucket = append(bucket, h)
	ht.RoutingTable[index] = bucket

	ht.MarkHostAsSeen(h.ID.Bytes())
}

func TestHashTable_DoesHostExistInBucket(t *testing.T) {
	prefix := "127.0.0.1:"
	port := 3000

	origin, _ := id.NewID()
	addr, _ := host.NewAddress(prefix + strconv.Itoa(port))
	ht, _ := NewHashTable(origin, addr)
	port++

	id1, _ := id.NewID()
	addr, _ = host.NewAddress(prefix + strconv.Itoa(port))
	h := NewRouteHost(&host.Host{ID: id1, Address: addr})

	index := GetBucketIndexFromDifferingBit(origin, id1)
	assert.False(t, ht.DoesHostExistInBucket(index, h.ID.Bytes()))

	bucket := ht.RoutingTable[index]
	bucket = append(bucket, h)
	ht.RoutingTable[index] = bucket

	assert.True(t, ht.DoesHostExistInBucket(index, h.ID.Bytes()))
}

func TestHashTable_GetClosestContacts(t *testing.T) {
	prefix := "127.0.0.1:"
	port := 3000

	origin, _ := id.NewID()
	addr, _ := host.NewAddress(prefix + strconv.Itoa(port))
	ht, _ := NewHashTable(origin, addr)
	port++

	id1, _ := id.NewID()
	addr, _ = host.NewAddress(prefix + strconv.Itoa(port))
	h := NewRouteHost(&host.Host{ID: id1, Address: addr})

	index := GetBucketIndexFromDifferingBit(origin, id1)
	bucket := ht.RoutingTable[index]
	bucket = append(bucket, h)
	ht.RoutingTable[index] = bucket

	assert.NotNil(t, ht.GetClosestContacts(1, h.ID.Bytes(), nil))
}

func TestHashTable_GetAllHostsInBucketCloserThan(t *testing.T) {
	prefix := "127.0.0.1:"
	port := 3000

	origin, _ := id.NewID()
	addr, _ := host.NewAddress(prefix + strconv.Itoa(port))
	ht, _ := NewHashTable(origin, addr)
	port++

	id1, _ := id.NewID()
	addr, _ = host.NewAddress(prefix + strconv.Itoa(port))
	h := NewRouteHost(&host.Host{ID: id1, Address: addr})

	index := GetBucketIndexFromDifferingBit(origin, id1)
	bucket := ht.RoutingTable[index]
	bucket = append(bucket, h)
	ht.RoutingTable[index] = bucket

	ht.GetAllHostsInBucketCloserThan(index, h.ID.Bytes())
}

func TestGetBucketIndexFromDifferingBit(t *testing.T) {
	// binary: 1....0
	var id1 id.ID = []byte{
		0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	// binary: 0....0
	var id2 id.ID = []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	// binary: 0....1
	var id3 id.ID = []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}

	// symmetry check
	for i := 0; i < 10; i++ {
		first, _ := id.NewID()
		second, _ := id.NewID()
		assert.Equal(t,
			GetBucketIndexFromDifferingBit(first, second),
			GetBucketIndexFromDifferingBit(second, first))
	}

	assert.Equal(t, 159, GetBucketIndexFromDifferingBit(id1, id2))
	assert.Equal(t, 159, GetBucketIndexFromDifferingBit(id1, id3))
	assert.Equal(t, 0, GetBucketIndexFromDifferingBit(id2, id3))
	assert.Equal(t, 0, GetBucketIndexFromDifferingBit(id1, id1))
}

func TestHashTable_GetClosestContacts2(t *testing.T) {
	var id1 id.ID = []byte{
		0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	var id2 id.ID = []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xFF}
	var id3 id.ID = []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}
	var id4 id.ID = []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2}
	var origin id.ID = []byte{
		0x0, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
		0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0}

	prefix := "127.0.0.1:"
	port := 3000
	genAddress := func() *host.Address {
		address, _ := host.NewAddress(prefix + strconv.Itoa(port))
		port++
		return address
	}

	ht, _ := NewHashTable(origin, genAddress())

	appendId := func(identifier id.ID) {
		h := NewRouteHost(&host.Host{ID: identifier, Address: genAddress()})
		index := GetBucketIndexFromDifferingBit(origin, identifier)
		bucket := ht.RoutingTable[index]
		bucket = append(bucket, h)
		ht.RoutingTable[index] = bucket
	}

	appendId(id1)
	appendId(id2)
	appendId(id3)

	assert.Equal(t, id3, ht.GetClosestContacts(1, id4, nil).hosts[0].ID)
	assert.Equal(t, id3, ht.GetClosestContacts(2, id4, nil).hosts[0].ID)
	assert.Equal(t, id2, ht.GetClosestContacts(2, id4, nil).hosts[1].ID)
}

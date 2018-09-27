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
	"bytes"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
)

// IterateType is type of iteration.
type IterateType int

const (
	// IterateStore is iteration type for Store requests.
	IterateStore = IterateType(iota)

	// IterateBootstrap is iteration type for Bootstrap requests.
	IterateBootstrap

	// IterateFindHost is iteration type for FindHost requests.
	IterateFindHost

	// IterateFindValue is iteration type for FindValue requests.
	IterateFindValue
)

const (
	// ParallelCalls is a small number representing the degree of parallelism in network calls.
	ParallelCalls = 3

	// KeyBitSize is the size in bits of the keys used to identify hosts and store and
	// retrieve data; in basic Kademlia this is 160, the length of a SHA1.
	KeyBitSize = 160

	// MaxContactsInBucket the maximum number of contacts stored in a bucket.
	MaxContactsInBucket = 20
)

// HashTable represents the hash-table state.
type HashTable struct {
	// The local host.
	Origin *host.Host

	// Routing table a list of all known hosts in the network
	// Hosts within buckets are sorted by least recently seen e.g.
	// [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
	//  ^                                                           ^
	//  └ Least recently seen                    Most recently seen ┘
	RoutingTable [][]*RouteHost // 160x20

	mutex *sync.Mutex

	refreshMap [KeyBitSize]time.Time

	rand *rand.Rand
}

// NewHashTable creates new HashTable.
func NewHashTable(id id.ID, address *host.Address) (*HashTable, error) {
	if id.Bytes() == nil {
		return nil, errors.New("id required")
	}

	ht := &HashTable{
		mutex: &sync.Mutex{},
		Origin: &host.Host{
			ID:      id,
			Address: address,
		},
	}

	ht.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < KeyBitSize; i++ {
		ht.ResetRefreshTimeForBucket(i)
	}

	for i := 0; i < KeyBitSize; i++ {
		ht.RoutingTable = append(ht.RoutingTable, []*RouteHost{})
	}

	return ht, nil
}

// Lock locks internal table mutex.
func (ht *HashTable) Lock() {
	ht.mutex.Lock()
}

// Unlock unlocks internal table mutex.
func (ht *HashTable) Unlock() {
	ht.mutex.Unlock()
}

// ResetRefreshTimeForBucket resets refresh timer for given bucket.
func (ht *HashTable) ResetRefreshTimeForBucket(bucket int) {
	ht.Lock()
	defer ht.Unlock()

	ht.refreshMap[bucket] = time.Now()
}

// GetRefreshTimeForBucket returns Time when given bucket must be refreshed.
func (ht *HashTable) GetRefreshTimeForBucket(bucket int) time.Time {
	ht.Lock()
	defer ht.Unlock()

	return ht.refreshMap[bucket]
}

// MarkHostAsSeen marks given Host as seen.
func (ht *HashTable) MarkHostAsSeen(host []byte) {
	ht.Lock()
	defer ht.Unlock()

	index := GetBucketIndexFromDifferingBit(ht.Origin.ID.Bytes(), host)
	bucket := ht.RoutingTable[index]
	hostIndex := -1
	for i, v := range bucket {
		if bytes.Equal(v.ID.Bytes(), host) {
			hostIndex = i
			break
		}
	}
	if hostIndex == -1 {
		panic(errors.New("tried to mark nonexistent host as seen"))
	}

	n := bucket[hostIndex]
	bucket = append(bucket[:hostIndex], bucket[hostIndex+1:]...)
	bucket = append(bucket, n)
	ht.RoutingTable[index] = bucket
}

// DoesHostExistInBucket checks if given Host exists in given bucket.
func (ht *HashTable) DoesHostExistInBucket(bucket int, host []byte) bool {
	ht.Lock()
	defer ht.Unlock()

	for _, v := range ht.RoutingTable[bucket] {
		if bytes.Equal(v.ID.Bytes(), host) {
			return true
		}
	}
	return false
}

// GetClosestContacts returns RouteSet with num closest Hosts to target.
func (ht *HashTable) GetClosestContacts(num int, target []byte, ignoredHosts []*host.Host) *RouteSet {
	ht.Lock()
	defer ht.Unlock()
	// First we need to build the list of adjacent indices to our target
	// in order
	index := GetBucketIndexFromDifferingBit(ht.Origin.ID.Bytes(), target)
	indexList := []int{index}
	i := index - 1
	j := index + 1
	for len(indexList) < KeyBitSize {
		if j < KeyBitSize {
			indexList = append(indexList, j)
		}
		if i >= 0 {
			indexList = append(indexList, i)
		}
		i--
		j++
	}

	routeSet := NewRouteSet()
	leftToAdd := num

	// Next we select ParallelCalls contacts and add them to the route set
	ht.selectParallelCalls(leftToAdd, indexList, ignoredHosts, routeSet)

	sort.Sort(routeSet)

	return routeSet
}

func (ht *HashTable) selectParallelCalls(
	leftToAdd int,
	indexList []int,
	ignoredHosts []*host.Host,
	routeSet *RouteSet,
) {
	var index int
	for leftToAdd > 0 && len(indexList) > 0 {
		index, indexList = indexList[0], indexList[1:]
		bucketContacts := len(ht.RoutingTable[index])
		for i := 0; i < bucketContacts; i++ {
			ignored := false
			for j := 0; j < len(ignoredHosts); j++ {
				if ht.RoutingTable[index][i].ID.Equal(ignoredHosts[j].ID.Bytes()) {
					ignored = true
				}
			}
			if !ignored {
				routeSet.Append(ht.RoutingTable[index][i])
				leftToAdd--
				if leftToAdd == 0 {
					break
				}
			}
		}
	}
}

// GetAllHostsInBucketCloserThan returns all hosts from given bucket that are closer to id then our host.
func (ht *HashTable) GetAllHostsInBucketCloserThan(bucket int, id []byte) [][]byte {
	b := ht.RoutingTable[bucket]
	var hosts [][]byte
	for _, v := range b {
		d1 := ht.getDistance(id, ht.Origin.ID.Bytes())
		d2 := ht.getDistance(id, v.ID.Bytes())

		result := d1.Sub(d1, d2)
		if result.Sign() > -1 {
			hosts = append(hosts, v.ID.Bytes())
		}
	}

	return hosts
}

// GetTotalHostsInBucket returns number of Hosts in bucket.
func (ht *HashTable) GetTotalHostsInBucket(bucket int) int {
	ht.Lock()
	defer ht.Unlock()
	return len(ht.RoutingTable[bucket])
}

func (ht *HashTable) getDistance(id1, id2 []byte) *big.Int {
	var dst [MaxContactsInBucket]byte
	for i := 0; i < MaxContactsInBucket; i++ {
		dst[i] = id1[i] ^ id2[i]
	}
	ret := big.NewInt(0)
	return ret.SetBytes(dst[:])
}

// GetRandomIDFromBucket returns random host ID from given bucket.
func (ht *HashTable) GetRandomIDFromBucket(bucket int) []byte {
	ht.Lock()
	defer ht.Unlock()
	// Set the new requestID to to be equal in every byte up to
	// the byte of the first differing bit in the bucket

	byteIndex := bucket / 8
	var id1 []byte
	for i := 0; i < byteIndex; i++ {
		id1 = append(id1, ht.Origin.ID.Bytes()[i])
	}
	differingBitStart := bucket % 8

	var firstByte byte
	// check each bit from left to right in order
	for i := 0; i < 8; i++ {
		// Set the value of the bit to be the same as the requestID
		// up to the differing bit. Then begin randomizing
		var bit bool
		if i < differingBitStart {
			bit = hasBit(ht.Origin.ID.Bytes()[byteIndex], uint8(i))
		} else {
			bit = ht.rand.Intn(2) == 1
		}

		if bit {
			firstByte += byte(math.Pow(2, float64(7-i)))
		}
	}

	id1 = append(id1, firstByte)

	// Randomize each remaining byte
	for i := byteIndex + 1; i < 20; i++ {
		randomByte := byte(ht.rand.Intn(256))
		id1 = append(id1, randomByte)
	}

	return id1
}

// GetBucketIndexFromDifferingBit returns appropriate bucket number for two host IDs.
func GetBucketIndexFromDifferingBit(id1, id2 []byte) int {
	// Look at each byte from left to right
	for j := 0; j < len(id1); j++ {
		// xor the byte
		xor := id1[j] ^ id2[j]

		// check each bit on the xored result from left to right in order
		for i := 0; i < 8; i++ {
			if hasBit(xor, uint8(i)) {
				byteIndex := j * 8
				bitIndex := i
				return KeyBitSize - (byteIndex + bitIndex) - 1
			}
		}
	}

	// the ids must be the same
	// this should only happen during bootstrapping
	return 0
}

// TotalHosts returns total number of hosts in HashTable.
func (ht *HashTable) TotalHosts() int {
	ht.Lock()
	defer ht.Unlock()
	return ht.totalHosts()
}

func (ht *HashTable) totalHosts() int {
	var total int
	for _, v := range ht.RoutingTable {
		total += len(v)
	}
	return total
}

// GetHosts returns hostsNumber (or less) hosts from the HashTable, starting from the first bucket
func (ht *HashTable) GetHosts(hostsNumber int) []host.Host {
	ht.Lock()
	defer ht.Unlock()

	result := make([]host.Host, 0)

Loop:
	for _, bucket := range ht.RoutingTable {
		for _, routeHost := range bucket {
			result = append(result, *routeHost.Host)
			if len(result) == hostsNumber {
				break Loop
			}
		}
	}
	return result
}

// GetMulticastHosts returns snapshot of all hosts from the HashTable
func (ht *HashTable) GetMulticastHosts() []*RouteHost {
	ht.Lock()
	defer ht.Unlock()

	result := make([]*RouteHost, ht.totalHosts())
	index := 0
	for _, bucket := range ht.RoutingTable {
		copy(result[index:], bucket)
		index += len(bucket)
	}
	return result
}

// hasBit is a Simple helper function to determine the value of a particular
// bit in a byte by index
// Example:
// number:  1
// bits:    00000001
// pos:     01234567
func hasBit(n byte, pos uint8) bool {
	pos = 7 - pos
	val := n & (1 << pos)
	return val > 0
}

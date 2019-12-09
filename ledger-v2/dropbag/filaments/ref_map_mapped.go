//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package filaments

import (
	"math"
	"reflect"

	"github.com/insolar/insolar/ledger-v2/unsafekit"
	"github.com/insolar/insolar/reference"
)

type MappedRefMap struct {
	valueType        reflect.Type
	hashSeed         uint32
	expectedKeyCount int
	loadedKeyCount   int

	buckets []mappedRefBucket
	unsafekit.KeepAliveList
}

func (p *MappedRefMap) AddBucket(bucketIndex int, locator int64) {
	if len(p.buckets) != bucketIndex {
		panic("illegal value")
	}
	p.buckets = append(p.buckets, mappedRefBucket{locator: locator})
}

var emptyBucketMarker = make([]bucketKey, 0, 0)

func (p *MappedRefMap) LoadBucket(bucketIndex, bucketKeyCount int, chunks [][]byte) {
	bucket := &p.buckets[bucketIndex]
	switch {
	case p.expectedKeyCount < p.loadedKeyCount+bucketKeyCount:
		panic("illegal value")
	case bucket.keysL0 != nil:
		panic("illegal state")
	case bucketKeyCount > 0:
		//
	case len(chunks) > 0:
		panic("illegal value - unaligned chunk")
	default:
		bucket.keysL0 = make([]bucketKey, 0, 0) // empty bucket marker
		return
	}
	_ = bucket.loadBucket(bucketKeyCount, chunks)
	p.loadedKeyCount += bucketKeyCount
}

type mappedRefBucket struct {
	locator int64

	keysL0 []bucketKey
	keysL1 [][]bucketKey
}

func (p *mappedRefBucket) loadBucket(bucketKeyCount int, chunks [][]byte) error {
	p.keysL0 = make([]bucketKey, 0, bucketKeyCount)
	bucketKeyL0size := int(bucketKeyType.Size())
	keysL1count := 0

	var lastChunk []byte
	chunkIndex := 0

	remainingL0keys := bucketKeyCount
	for remainingL0keys > 0 {
		if len(lastChunk) == 0 {
			lastChunk = chunks[chunkIndex]
			chunkIndex++
		}

		keyChunk := lastChunk
		chunkKeys := len(lastChunk) / bucketKeyL0size

		switch {
		case chunkKeys > remainingL0keys:
			chunkKeys = remainingL0keys
			remainingL0keys = 0
			end := chunkKeys * bucketKeyL0size
			keyChunk = lastChunk[:end]
			lastChunk = lastChunk[end:]
		case len(lastChunk)%bucketKeyL0size == 0:
			remainingL0keys -= chunkKeys
			lastChunk = nil
		default:
			panic("illegal value - unaligned chunk")
		}

		keysL1count += p.addKeys(chunkKeys, keyChunk)
	}

	return p.loadBucketValues(keysL1count, lastChunk, chunks[chunkIndex:])
}

func (p *mappedRefBucket) loadBucketValues(keysL1count int, lastChunk []byte, chunks [][]byte) error {
	remainingL1keys := keysL1count

	chunkIndex := 0
	for remainingL1keys > 0 {
		if len(lastChunk) == 0 {
			lastChunk = chunks[chunkIndex]
			chunkIndex++
		}

	}
	//keysL1count := 0
	//
	//lastChunk := chunks[0]
	//chunkIndex := 0
	//
	//remainingL0keys := bucketKeyCount
	//for remainingL0keys > 0 {
	//	if len(lastChunk) == 0 {
	//		chunkIndex++
	//		lastChunk = chunks[chunkIndex]
	//	}
	//
	//	keyChunk := lastChunk
	//	chunkKeys := len(lastChunk) / bucketKeyL0size
	//
	//	switch {
	//	case chunkKeys > remainingL0keys:
	//		chunkKeys = remainingL0keys
	//		remainingL0keys = 0
	//		end := chunkKeys * bucketKeyL0size
	//		keyChunk = lastChunk[:end]
	//		lastChunk = lastChunk[end:]
	//	case len(lastChunk) % bucketKeyL0size == 0:
	//		remainingL0keys -= chunkKeys
	//		lastChunk = nil
	//	default:
	//		panic("illegal value - unaligned chunk")
	//	}
	//
	//	keysL1count += p.addKeys(chunkKeys, keyChunk)
	//}
	//

	//if len(lastChunk) != 0 && chunkIndex != len(chunks) {
	//	panic("illegal value - unaligned chunks")
	//}
	//
	return nil
}

func (p *mappedRefBucket) addKeys(keys int, chunk []byte) int {
	return 0
}

type bucketIndexArity uint8

const (
	specialOne bucketIndexArity = iota
	specialOneZero
	specialOneSelf
	specialMany
)

const countBits = 3 * 8
const positionMask = math.MaxUint64 >> countBits
const countMask = 1<<countBits - 1
const countShift = 64 - countBits

const arityBits = 2
const arityMask = 1<<arityBits - 1
const arityShift = 64 - arityBits

const countMaskWithArity = 1<<(countBits-arityBits) - 1

type bucketValue uint64

func (v bucketValue) getPos() int64 {
	return int64(v & positionMask)
}

func (v bucketValue) getCount() int {
	return int(v>>countShift) & countMask
}

func (v bucketValue) getCountWithArity() (int, bucketIndexArity) {
	return int(v>>countShift) & countMaskWithArity, bucketIndexArity(v>>arityShift) & arityMask
}

var bucketKeyType = reflect.TypeOf(bucketKey{})

type bucketKey struct {
	local reference.Local
	index bucketValue
}

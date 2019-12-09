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
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/reference"
)

type MappedRefMap struct {
	valueType        reflect.Type
	hashSeed         uint32
	expectedKeyCount int
	loadedKeyCount   int

	buckets []mappedRefBucket
	//	unsafekit.KeepAliveList
}

func (p *MappedRefMap) AddBucket(bucketIndex int, locator int64) {
	if len(p.buckets) != bucketIndex {
		panic("illegal value")
	}
	p.buckets = append(p.buckets, mappedRefBucket{locator: locator})
}

const MinKeyBucketBatchSize = 16

var emptyBucketMarker = make([][]bucketKey, 0, 0)

func (p *MappedRefMap) LoadBucket(bucketIndex, bucketKeyCount int, locator int64, chunks []longbits.ByteString) error {
	if bucketIndex < 0 || bucketIndex >= len(p.buckets) {
		panic("illegal value") // TODO error
	}

	bucket := &p.buckets[bucketIndex]
	switch {
	case p.expectedKeyCount < p.loadedKeyCount+bucketKeyCount:
		panic("illegal value") // TODO error
	case bucket.keysL0 != nil:
		panic("illegal state") // TODO error
	case bucketKeyCount > 0:
		// break
	case len(chunks) > 0:
		panic("illegal value - unaligned chunk") // TODO error
	default:
		bucket.locator = locator
		bucket.keysL0 = emptyBucketMarker
		return nil
	}

	loader := newBucketKeyLoader(chunks)
	keysL0, err := loader.loadKeys(bucketKeyCount)
	if err != nil {
		return err
	}

	countL1, countNV := countKeysOfL1(keysL0)
	if p.expectedKeyCount != countL1+countNV {
		panic("illegal state") // TODO error
	}

	keysL1, err := loader.loadKeys(countL1)
	if err != nil {
		return err
	}

	bucket.locator = locator
	bucket.keysL0 = keysL0
	bucket.keysL1 = keysL1
	p.loadedKeyCount += bucketKeyCount
	return nil
}

func countKeysOfL1(keysL0 [][]bucketKey) (countL1, countNV int) {
	for _, bkr := range keysL0 {
		for _, kb := range bkr {
			switch c, a := kb.index.getCountWithArity(); a {
			case specialOne:
				countL1++
			case specialOneZero, specialOneSelf:
				countNV++
			case specialMany:
				countL1 += c
			default:
				panic("unexpected")
			}
		}
	}
	return
}

type mappedRefBucket struct {
	locator int64

	keysL0 [][]bucketKey
	keysL1 [][]bucketKey
}

var bucketKeyTypeSlice = unsafekit.MustMMapSliceType(reflect.TypeOf([]bucketKey(nil)), false)
var bucketKeyType = bucketKeyTypeSlice.Elem()

type bucketKey struct {
	local reference.Local
	index bucketValue
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

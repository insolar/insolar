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

package refmap

import (
	"math"
	"math/bits"
	"reflect"
	"runtime"
	"sort"

	"github.com/insolar/insolar/ledger-v2/unsafekit"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/reference"
)

const MinKeyBucketBatchSize = 8 // min batch size for multi-batch buckets (except the last batch)

const tooBigUnsortedBucketSize = 32 // default size limit to force use of hashmap on an unsorted bucket. Excess impacts GC performance.
const tooBigSortedBucketSize = 1024 // default size limit to force use of hashmap on a sorted bucket. Excess impacts GC performance.

func NewReadOnlyMapToInt63(expectedKeyCount int, bucketCount int) ReadOnlyMapToInt63 {
	switch {
	case expectedKeyCount < 0:
		panic("illegal state")
	case bucketCount > 0:
		//
	case expectedKeyCount > 0:
		panic("illegal state")
	}
	return ReadOnlyMapToInt63{expectedKeyCount: expectedKeyCount, buckets: make([]mappedBucket, bucketCount)}
}

// implements READ-ONLY map[Holder]uint63 with external lazy load & memory-mapping support
type ReadOnlyMapToInt63 struct {
	hashSeed         uint32
	expectedKeyCount int
	loadedKeyCount   int
	bigBucketMinSize int

	//	unsafekit.KeepAliveList

	// must be rarely used as it impacts GC due references
	// per-bucket maps are required to avoid copy of keys inside a map for a composite key
	bigBuckets map[ /* bucketIndex */ uint32]bigBucketMap
	buckets    []mappedBucket

	sortedBuckets bool
}

type bigBucketMap map[ /* keyL0*/ longbits.ByteString]bucketValueSelector

func (p *ReadOnlyMapToInt63) SetBigBucketMinSize(bigBucketSize int) {
	switch {
	case p.loadedKeyCount > 0:
		panic("illegal state")
	case p.bigBucketMinSize != 0:
		panic("illegal state")
	case bigBucketSize <= MinKeyBucketBatchSize:
		panic("illegal state")
	}
	p.bigBucketMinSize = bigBucketSize
}

func (p *ReadOnlyMapToInt63) SetHashSeed(hashSeed uint32) {
	switch {
	case p.loadedKeyCount > 0:
		panic("illegal state")
	case p.bigBucketMinSize != 0:
		panic("illegal state")
	}
	p.hashSeed = hashSeed
}

func (p *ReadOnlyMapToInt63) GetHashSeed() uint32 {
	return p.hashSeed
}

func (p *ReadOnlyMapToInt63) SetSortedBuckets(sortedBuckets bool) {
	switch {
	case p.loadedKeyCount > 0:
		panic("illegal state")
	case p.bigBucketMinSize != 0:
		panic("illegal state")
	}
	p.sortedBuckets = sortedBuckets
}

func (p *ReadOnlyMapToInt63) SetLocator(bucketIndex int, locator int64) {
	if locator < 0 {
		panic("illegal value")
	}
	p.buckets[bucketIndex].locator = locator
}

func (p *ReadOnlyMapToInt63) GetLocator(bucketIndex int) int64 {
	return p.buckets[bucketIndex].locator
}

// result = ( N, true   ) - item was found, N = map[ref], N = [0..maxInt64]
// result = ( <0, false ) - item was not found
// result = ( B>=0, false) - item presence is unknown, bucket is missing, B is bucket number
//
func (p *ReadOnlyMapToInt63) GetValueOrBucket(ref reference.Holder) ( /* map value or missing bucket index */ int64, bool) {
	localRef := ref.GetLocal()
	s := unsafekit.WrapLocalRef(localRef)
	hashLocal := hash32(s, p.hashSeed) // localRef must be kept alive

	indexL0 := hashLocal % uint32(len(p.buckets)) // TODO use bitmask
	bucket := &p.buckets[indexL0]
	if !bucket.isLoaded() {
		return int64(indexL0), false
	}

	selectorL1 := bucketValueSelector(0)
	{
		ok := false
		if !bucket.isBigBucket() {
			largeBucket := p.bigBuckets[indexL0]
			if p.bigBucketMinSize == 0 || len(largeBucket) < p.bigBucketMinSize {
				panic("illegal state")
			}
			selectorL1, ok = largeBucket[s] // localRef must be kept alive for this operation
		} else {
			selectorL1, ok = bucket.findKeyL0(localRef, p.sortedBuckets)
		}
		runtime.KeepAlive(localRef)
		if !ok {
			return -1, false
		}
	}

	if selectorL1.isLeaf() {
		if ref.GetBase().Equal(*localRef) {
			return selectorL1.asValue(), true
		}
		return -1, false
	}

	if v, ok := bucket.findValue(selectorL1, ref.GetBase(), p.sortedBuckets); ok {
		return v.asValue(), ok
	} else {
		return -1, false
	}
}

var emptyBucketMarker = make([][]bucketKey, 0, 0)

func (p *ReadOnlyMapToInt63) LoadBucket(bucketIndex, bucketKeyL0Count, bucketKeyL1Count int, chunks []longbits.ByteString) error {
	if bucketIndex < 0 || bucketIndex >= len(p.buckets) {
		panic("illegal value") // TODO error
	}

	bucket := &p.buckets[bucketIndex]
	switch {
	case bucketKeyL1Count > bucketKeyL0Count:
		panic("illegal value")
	case p.expectedKeyCount < p.loadedKeyCount+bucketKeyL0Count:
		panic("illegal value") // TODO error
	case bucket.keysL0 != nil || bucket.pageBitsL0 == tooBigBatch:
		panic("illegal state") // TODO error
	case bucketKeyL0Count > 0:
		// break
	case len(chunks) > 0:
		panic("illegal value - unaligned chunk") // TODO error
	default:
		bucket.keysL0 = emptyBucketMarker
		return nil
	}

	switch {
	case p.bigBucketMinSize > 0:
		//
	case p.sortedBuckets:
		p.bigBucketMinSize = tooBigSortedBucketSize
	default:
		p.bigBucketMinSize = tooBigUnsortedBucketSize
	}

	loader := newBucketKeyLoader(chunks)

	if bucketKeyL0Count >= p.bigBucketMinSize {
		switch m, err := loader.loadKeysL0AsMap(bucketKeyL0Count, bucketKeyL1Count); {
		case err != nil:
			return err
		case p.bigBuckets == nil:
			p.bigBuckets = make(map[ /* bucketIndex */ uint32]bigBucketMap)
			fallthrough
		default:
			p.bigBuckets[uint32(bucketIndex)] = m
			bucket.pageBitsL0 = tooBigBatch
			p.loadedKeyCount += bucketKeyL0Count
			return nil
		}
	}

	keysL0, err := loader.loadKeys(bucketKeyL0Count)
	if err != nil {
		return err
	}

	keysL1, err := loader.loadKeys(bucketKeyL1Count)
	if err != nil {
		return err
	}

	bucket.keysL0 = keysL0
	bucket.pageBitsL0 = makePageBits(keysL0)
	bucket.keysL1 = keysL1
	bucket.pageBitsL1 = makePageBits(keysL1)
	p.loadedKeyCount += bucketKeyL0Count
	return nil
}

func makePageBits(keys [][]bucketKey) uint8 {
	if len(keys) < 2 {
		return 0
	}
	n := uint(len(keys[0]))
	if n >= MinKeyBucketBatchSize && args.IsPowerOfTwo(n) {
		if b := uint8(bits.Len(n)); b <= maxPowerOfTwoForBatch {
			return b
		}
	}
	return nonPowerOfTwoBatchAlignment
}

type KeyIndexGetterFunc func(index int) *bucketKey

type mappedBucket struct {
	locator int64

	keysL0 [][]bucketKey
	keysL1 [][]bucketKey

	pageBitsL0, pageBitsL1 uint8
}

func (b mappedBucket) isLoaded() bool {
	return b.keysL0 != nil || b.pageBitsL0 == tooBigBatch
}

func (b mappedBucket) isBigBucket() bool {
	return b.pageBitsL0 == tooBigBatch
}

const nonPowerOfTwoBatchAlignment = math.MaxUint8
const tooBigBatch = nonPowerOfTwoBatchAlignment - 1
const maxPowerOfTwoForBatch = 0x1F // also works as a hint for compiler optimization for shift ops

func getKeyCount(keys [][]bucketKey) int {
	switch batchCount := len(keys); batchCount {
	case 0:
		return 0
	case 1:
		return len(keys[0])
	default:
		return len(keys[0])*(batchCount-1) + len(keys[batchCount-1])
	}
}

type indexerFunc func(int) (pageNum, inPageIndex int)

func getKeyIndexer(keys [][]bucketKey, pageBits uint8) indexerFunc {
	switch pageBits {
	case 0:
		return func(index int) (page, pos int) {
			return 0, index
		}
	case nonPowerOfTwoBatchAlignment:
		base := len(keys[0])
		return func(index int) (page, pos int) {
			return index / base, index % base
		}
	default:
		b := pageBits & maxPowerOfTwoForBatch // hint for compiler optimization for shift ops
		if b != pageBits {
			panic("illegal state")
		}
		mask := 1<<b - 1
		return func(index int) (page, pos int) {
			return index >> b, index &^ mask
		}
	}
}

func findKey(ref *reference.Local, keys [][]bucketKey, pageBits uint8, sorted bool) (bucketValueSelector, bool) {
	if !sorted {
		for _, bkr := range keys {
			for _, br := range bkr {
				if ref.Equal(br.local) {
					return br.value, true
				}
			}
		}
		return 0, false
	}
	return findInSortedBatches(ref, keys, getKeyIndexer(keys, pageBits), getKeyCount(keys))
}

func findInSortedBatches(ref *reference.Local, keys [][]bucketKey, indexFn indexerFunc, count int) (bucketValueSelector, bool) {
	pos := sort.Search(count, func(i int) bool {
		pgn, pgi := indexFn(i)
		v0 := &keys[pgn][pgi]
		return ref.Compare(v0.local) >= 0
	})

	if pos < count {
		pgn, pgi := indexFn(pos)
		if v0 := &keys[pgn][pgi]; ref.Equal(v0.local) {
			return v0.value, true
		}
	}
	return 0, false
}

func (b mappedBucket) findKeyL0(ref *reference.Local, sorted bool) (bucketValueSelector, bool) {
	return findKey(ref, b.keysL0, b.pageBitsL0, sorted)
}

func (b mappedBucket) findValue(selectorL1 bucketValueSelector, ref *reference.Local, sorted bool) (bucketValueSelector, bool) {
	switch posL1, countL1 := selectorL1.asPosAndCount(); {
	case posL1 > math.MaxInt32:
		panic("illegal value")
	case countL1 == 1: // main case
		indexFn := getKeyIndexer(b.keysL1, b.pageBitsL1)
		pgn, pgi := indexFn(posL1)
		br := &b.keysL1[pgn][pgi]
		return br.value, ref.Equal(br.local)

	case countL1 <= 0:
		panic("illegal state")
	default:
		indexFn := getKeyIndexer(b.keysL1, b.pageBitsL1)
		return _findValue(ref, b.keysL1, indexFn, posL1, countL1, sorted)
	}
}

func _findValue(ref *reference.Local, keysL1 [][]bucketKey, indexFn indexerFunc, posL1, countL1 int, sorted bool) (bucketValueSelector, bool) {

	switch firstPgN, firstIdx := indexFn(posL1); {
	case !sorted:
		if firstIdx != 0 {
			firstPage := keysL1[firstPgN]
			for i, max := firstIdx, len(firstPage); i < max; i++ {
				br := &firstPage[i]
				if ref.Equal(br.local) {
					return br.value, true
				}
				if countL1--; countL1 <= 0 {
					return 0, false
				}
			}
			firstPgN++
		}
		for _, bkr := range keysL1[firstPgN:] {
			for _, br := range bkr {
				if ref.Equal(br.local) {
					return br.value, true
				}
				if countL1--; countL1 <= 0 {
					return 0, false
				}
			}
		}
		return 0, false

	case firstIdx != 0:
		prevFn := indexFn
		indexFn = func(i int) (pageNum, inPageIndex int) {
			return prevFn(i + firstIdx)
		}
		fallthrough
	default:
		return findInSortedBatches(ref, keysL1[firstPgN:], indexFn, countL1)
	}
}

var bucketKeyTypeSlice = unsafekit.MustMMapSliceType(reflect.TypeOf([]bucketKey(nil)), false)
var bucketKeyType = bucketKeyTypeSlice.Elem()

type bucketKey struct {
	local reference.Local
	value bucketValueSelector
}

const bucketKeySelectorFlag = 1 << 63
const countBits = 32
const positionMask = math.MaxUint64 >> countBits
const countMask = 1<<(countBits-1) - 1
const countShift = 64 - countBits

type bucketValueSelector int64 // positive values only

func (v bucketValueSelector) isLeaf() bool {
	return uint64(v)&bucketKeySelectorFlag == 0
}

func (v bucketValueSelector) asPosAndCount() (int, int) {
	return int(v) & positionMask, int(v>>countShift) & countMask
}

func (v bucketValueSelector) asValue() int64 {
	return int64(v)
}

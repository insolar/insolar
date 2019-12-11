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
	"math/bits"
	"sort"

	"github.com/insolar/insolar/reference"
)

type KeyWriterFunc func(*reference.Local, BucketState) error
type KeyResolveFunc func(id uint32) *reference.Local
type KeyBucketFunc func(count int, writeFn func(KeyWriterFunc) error) error
type KeyBucketFactoryFunc func(resolveFn KeyResolveFunc, bucketCount, bucketSize uint32) KeyBucketFunc

type WriteBucketer struct {
	MaxPerBucket int
	KeySorter    func(i, j *reference.Local) bool
	Output       KeyBucketFactoryFunc
}

func (p WriteBucketer) ProcessMap(m *UpdateableKeyMap) error {
	if p.MaxPerBucket < 16 {
		panic("illegal argument")
	}

	keyCount := uint32(m.InternedKeyCount())
	if keyCount == 0 {
		p.Output(m.GetInterned, 0, 0)
		return nil
	}

	bucketCount := (keyCount + uint32(p.MaxPerBucket) - 1) / uint32(p.MaxPerBucket)
	bitsPerBucket := uint8(bits.Len(uint(bucketCount)))
	bucketCount = 1 << bitsPerBucket
	bucketMask := bucketCount - 1

	bucketSize := (keyCount + bucketMask) / bucketCount
	bucketSize += bucketSize>>2 + 1
	keyCount = bucketSize * bucketCount

	buckets := make([]uint32, keyCount)
	var overflow map[uint32][]uint32

	for j, page := range m.buckets {
		pageOffset := j * 1 << m.pageBits
		for i := range page {
			mapBucket := &page[i]
			if mapBucket.IsEmpty() {
				continue
			}

			mapBucketIndex := uint32(i + pageOffset)
			bucketBase := (mapBucket.refHash & bucketMask) * bucketSize

			counter := buckets[bucketBase] + 1
			buckets[bucketBase] = counter
			if counter < bucketSize {
				buckets[bucketBase+counter] = mapBucketIndex
				continue
			}

			if overflow == nil {
				overflow = make(map[uint32][]uint32)
			}
			overflow[bucketBase] = append(overflow[bucketBase], mapBucketIndex)
		}
	}

	writeFn := p.Output(m.GetInterned, bucketCount, bucketSize)
	if p.KeySorter != nil {
		sorter := refMapKeySorter{p.KeySorter, writeFn, nil}
		writeFn = sorter.WriteSorted
	}

	for bucketBase := uint32(0); bucketBase < uint32(len(buckets)); bucketBase += bucketSize {
		counter := buckets[bucketBase]
		if counter == 0 {
			continue
		}

		overflowCount := uint32(0)
		bucketMain := buckets[bucketBase+1:]
		if counter < bucketSize {
			bucketMain = bucketMain[:counter]
		} else {
			bucketMain = bucketMain[:bucketSize-1]
			overflowCount = counter - bucketSize
		}

		var bucketOverflow []uint32
		if overflowCount > 0 {
			bucketOverflow = overflow[bucketBase]
			delete(overflow, bucketBase)
			if uint32(len(bucketOverflow)) != overflowCount {
				panic("illegal state")
			}
		}

		if err := writeFn(int(counter), func(fn KeyWriterFunc) error {
			for _, bn := range bucketMain {
				bucket := m.getBucket(bn)
				if err := fn(bucket.localRef, bucket.state); err != nil {
					return err
				}
			}
			for _, bn := range bucketOverflow {
				bucket := m.getBucket(bn)
				if err := fn(bucket.localRef, bucket.state); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

type refMapKeySorter struct {
	keySorter func(i, j *reference.Local) bool
	output    KeyBucketFunc
	items     []keyBucketItem
}

type keyBucketItem struct {
	*reference.Local
	BucketState
}

func (p refMapKeySorter) WriteSorted(count int, fn func(KeyWriterFunc) error) error {

	items := make([]keyBucketItem, 0, count)

	if err := fn(func(local *reference.Local, state BucketState) error {
		items = append(items, keyBucketItem{local, state})
		return nil
	}); err != nil {
		return err
	}

	p.items = items
	sort.Sort(p)
	p.items = nil

	return p.output(count, func(fn KeyWriterFunc) error {
		for _, it := range items {
			if err := fn(it.Local, it.BucketState); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p refMapKeySorter) Len() int {
	return len(p.items)
}

func (p refMapKeySorter) Less(i, j int) bool {
	return p.keySorter(p.items[i].Local, p.items[j].Local)
}

func (p refMapKeySorter) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

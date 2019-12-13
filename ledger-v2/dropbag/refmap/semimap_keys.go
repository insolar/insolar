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
	"runtime"

	"github.com/insolar/insolar/ledger-v2/fastrand"

	"github.com/insolar/insolar/ledger-v2/unsafekit"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/reference"
)

const MinBucketPageSize = 16
const MaxBucketPageSize = 65536

func NewUpdateableKeyMap() UpdateableKeyMap {
	return NewUpdateableKeyMapExt(128)
}

func NewUpdateableKeyMapExt(pageSize int) UpdateableKeyMap {
	switch {
	case pageSize < MinBucketPageSize:
		panic("illegal value")
	case pageSize > MaxBucketPageSize:
		panic("illegal value")
	}

	pageBits := uint8(bits.Len(uint(pageSize - 1)))
	alignedPageSize := uint32(1) << pageBits

	buckets := [][]refMapBucket{make([]refMapBucket, 1, alignedPageSize)}
	buckets[0][0] = refMapBucket{localRef: emptyLocalRef}

	return UpdateableKeyMap{
		map0:     map[longbits.ByteString]uint32{emptyLocalRefKey: 0},
		buckets:  buckets,
		pageBits: pageBits,
		pageMask: alignedPageSize - 1,
		hashSeed: fastrand.Uint32(),
	}
}

var emptyLocalRef = reference.EmptyLocal()
var emptyLocalRefKey = unsafekit.WrapLocalRef(emptyLocalRef)

type UpdateableKeyMap struct {
	hashSeed uint32

	map0    map[longbits.ByteString]uint32
	buckets [][]refMapBucket

	pageMask uint32
	pageBits uint8
}

type BucketState uint32

type BucketValueSelector struct {
	ValueSelector
	State BucketState
}

type ValueSelector struct {
	LocalId uint32 // BucketId
	BaseId  uint32
}

func (m *UpdateableKeyMap) GetHashSeed() uint32 {
	return m.hashSeed
}

func (m *UpdateableKeyMap) SetHashSeed(hashSeed uint32) {
	switch {
	case len(m.map0) > 1 || len(m.buckets) > 1:
		panic("illegal state")
	}
	m.hashSeed = hashSeed
}

func (m *UpdateableKeyMap) InternHolder(ref reference.Holder) reference.Holder {
	switch {
	case ref == nil:
		return nil
	case ref.IsEmpty():
		return reference.Empty()
	}

	p0, p1 := ref.GetLocal(), ref.GetBase()
	p0i, p1i := m.Intern(p0), p1
	if p0 == p1 {
		p1i = p0i
	} else {
		p1i = m.Intern(p1)
	}
	if p0 != p0i || p1 != p1i {
		return reference.NewNoCopy(p0i, p1i)
	}
	return ref
}

func (m *UpdateableKeyMap) InternedKeyCount() int {
	return len(m.map0)
}

func (m *UpdateableKeyMap) Intern(ref *reference.Local) *reference.Local {
	if ref == nil {
		return nil
	}
	_, r, _ := m.intern(ref)
	return r
}

func (m *UpdateableKeyMap) getBucket(bucketIndex uint32) *refMapBucket {
	return &m.buckets[bucketIndex>>m.pageBits][bucketIndex&m.pageMask]
}

func (m *UpdateableKeyMap) GetInterned(bucketIndex uint32) (*reference.Local, BucketState) {
	bucket := m.getBucket(bucketIndex)
	return bucket.localRef, bucket.state
}

func (m *UpdateableKeyMap) intern(ref *reference.Local) (uint32, *reference.Local, longbits.ByteString) {
	switch {
	case ref == nil:
		panic("illegal value")
	case ref.IsEmpty():
		return 0, reference.EmptyLocal(), ""
	}

	key := unsafekit.WrapLocalRef(ref) // MUST keep ref as long as key is in use
	if bucketIndex, ok := m.getIndexWithKey(ref, key); ok {
		return bucketIndex, m.getBucket(bucketIndex).localRef, key
	}

	pageIndex := len(m.buckets) - 1
	lastPage := &m.buckets[pageIndex]
	bucketIndex := uint32(pageIndex) << m.pageBits

	if pagePos := uint32(len(*lastPage)); pagePos > m.pageMask {
		m.buckets = append(m.buckets, make([]refMapBucket, 0, m.pageMask+1))
		pageIndex++
		lastPage = &m.buckets[pageIndex]
		bucketIndex += m.pageMask + 1
	} else {
		bucketIndex += pagePos
	}
	*lastPage = append(*lastPage, refMapBucket{localRef: ref})

	m.map0[key] = bucketIndex
	runtime.KeepAlive(ref)
	return bucketIndex, ref, key
}

func (m *UpdateableKeyMap) getIndexWithKey(ref *reference.Local, refKey longbits.ByteString) (uint32, bool) {
	bucketIndex, ok := m.map0[refKey]
	runtime.KeepAlive(ref) // make sure that (ref) stays while (refKey) is in use by mapaccess()
	return bucketIndex, ok
}

func (m *UpdateableKeyMap) getIndex(ref *reference.Local) (uint32, bool) {
	return m.getIndexWithKey(ref, unsafekit.WrapLocalRef(ref))
}

func (m *UpdateableKeyMap) Find(key reference.Holder) (BucketValueSelector, bool) {
	switch {
	case key.IsEmpty():
		return BucketValueSelector{}, false
	}

	var bucket *refMapBucket
	bucketIndex, ok := m.getIndex(key.GetLocal())
	if !ok {
		return BucketValueSelector{}, false
	}
	bucket = m.getBucket(bucketIndex)
	if bucket.IsEmpty() {
		return BucketValueSelector{}, false
	}

	if baseIndex, ok := m.getIndex(key.GetBase()); !ok {
		return BucketValueSelector{}, false
	} else {
		return BucketValueSelector{ValueSelector{bucketIndex, baseIndex}, bucket.state}, true
	}
}

func (m *UpdateableKeyMap) TryPut(key reference.Holder,
	valueFn func(internedKey reference.Holder, selector BucketValueSelector) BucketState,
) bool {
	switch {
	case valueFn == nil:
		panic("illegal value")
	case key.IsEmpty():
		return false
	}

	p0, p1 := key.GetLocal(), key.GetBase()
	bucket0, p0i, p0k := m.intern(p0)
	bucket1, p1i, _ := m.intern(p1)

	if p0 != p0i || p1 != p1i {
		key = reference.NewNoCopy(p0i, p1i)
	}

	bucket := m.getBucket(bucket0)
	prevState := bucket.state
	switch newState := valueFn(key, BucketValueSelector{ValueSelector{bucket0, bucket1}, prevState}); {
	case prevState == newState:
		return false
	case newState == 0:
		panic("illegal value")
	case prevState == 0 && bucket.refHash == 0:
		bucket.refHash = hash32(p0k, m.hashSeed)
		runtime.KeepAlive(p0) // ensures that (p0k) is ok
		fallthrough
	default:
		bucket.state = newState
		return true
	}
}

func (m *UpdateableKeyMap) TryTouch(key reference.Holder,
	valueFn func(selector BucketValueSelector) BucketState,
) bool {
	if valueFn == nil {
		panic("illegal value")
	}
	selector, ok := m.Find(key)
	if !ok {
		return false
	}

	bucket := m.getBucket(selector.LocalId)
	prevState := bucket.state
	switch newState := valueFn(selector); {
	case prevState == newState:
		return false
	default:
		bucket.state = newState
		return true
	}
}

type refMapBucket struct {
	localRef *reference.Local
	refHash  uint32
	state    BucketState
}

func (b refMapBucket) IsEmpty() bool {
	return b.state == 0
}

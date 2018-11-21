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

package packets

import (
	"io"

	"github.com/damnever/bitarray"
	"github.com/pkg/errors"
)

type DeviantBitSet struct {
	CompressedSet     bool
	HighBitLengthFlag bool
	LowBitLength      uint8
	HighBitLength     uint8
	Payload           []byte

	bucket []*BitSetBucket
	mapper BitSetMapper
}

func NewDeviantBitSet(buckets []*BitSetBucket, mapper BitSetMapper) (*DeviantBitSet, error) {
	bitset := &DeviantBitSet{
		bucket: buckets,
		mapper: mapper,
	}
	err := bitset.bucketToArray(buckets)
	if err != nil {
		return nil, err
	}
	return bitset, nil
}

func (dbs *DeviantBitSet) GetBuckets(mapper BitSetMapper) []*BitSetBucket {
	return dbs.bucket
}

func (dbs *DeviantBitSet) ApplyChanges(changes []*BitSetBucket) (BitSet, error) {
	for _, bucket := range changes {
		dbs.changeBucketState(bucket)
	}
	return dbs, nil
}

func (dbs *DeviantBitSet) Serialize() ([]byte, error) {
	return nil, nil
}

func (dbs *DeviantBitSet) Deserialize(data io.Reader) error {
	return nil
}

func (dbs *DeviantBitSet) changeBucketState(bucket *BitSetBucket) {
	for _, b := range dbs.bucket {
		if b.NodeID == bucket.NodeID {
			b.State = bucket.State
			return
		}
	}
	dbs.bucket = append(dbs.bucket, bucket)
}

func (dbs *DeviantBitSet) changeBitState(array *bitarray.BitArray, n int, state TriState) error {
	var err error
	switch state {
	case Legit:
		err = array.Clear(2*n, 2*n+1)
	case TimedOut:
		err = array.Clear(2*n, 2*n+1)
		if err != nil {
			return err
		}
		_, err = array.Put(2*n+1, 1)
	case Fraud:
		err = array.Clear(2*n, 2*n+1)
		if err != nil {
			return err
		}
		_, err = array.Put(2*n, 1)
	default:
		return errors.New("failed to change bit state: unknown state")
	}
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DeviantBitSet) bucketToArray(buckets []*BitSetBucket) error {
	array := bitarray.New(dbs.mapper.Length() * 2) // cuz stores 2 bits for 1 id
	for _, bucket := range buckets {
		n, err := dbs.mapper.RefToIndex(bucket.NodeID)
		if err != nil {
			return err
		}
		err = dbs.changeBitState(array, n, bucket.State)
		if err != nil {
			return err
		}
	}
	return nil
}

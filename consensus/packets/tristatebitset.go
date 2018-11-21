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
	"errors"
	"io"

	"github.com/damnever/bitarray"
)

const lastBitMask = 0x00000001

// TriStateBitSet bitset implementation.
type TriStateBitSet struct {
	CompressedSet     bool
	HighBitLengthFlag bool
	LowBitLength      uint8
	HighBitLength     uint8
	Payload           []byte

	cells  []*BitSetCell
	mapper BitSetMapper
}

// NewTriStateBitSet creates and returns a tristatebitset.
func NewTriStateBitSet(cells []*BitSetCell, mapper BitSetMapper) (*TriStateBitSet, error) {
	if (mapper == nil) || (cells == nil) {
		return nil, errors.New("failed to create tristatebitset")
	}
	bitset := &TriStateBitSet{
		cells:  cells,
		mapper: mapper,
	}
	err := bitset.bucketToArray(cells)
	if err != nil {
		return nil, err
	}
	return bitset, nil
}

func (dbs *TriStateBitSet) GetBuckets(mapper BitSetMapper) []*BitSetCell {
	return dbs.cells
}

func (dbs *TriStateBitSet) ApplyChanges(changes []*BitSetCell) (BitSet, error) {
	for _, cell := range changes {
		dbs.changeBucketState(cell)
	}
	return dbs, nil
}

func (dbs *TriStateBitSet) Serialize() ([]byte, error) {
	return nil, nil
}

func (dbs *TriStateBitSet) Deserialize(data io.Reader) error {
	return nil
}

func (dbs *TriStateBitSet) changeBucketState(cell *BitSetCell) error {
	n, err := dbs.mapper.RefToIndex(cell.NodeID)
	if err != nil {
		return err
	}
	dbs.cells[n] = cell
	return nil
}

func (dbs *TriStateBitSet) changeBitState(array *bitarray.BitArray, n int, state TriState) error {
	err := dbs.putLastBit(state, array, 2*n)
	if err != nil {
		return err
	}
	err = dbs.putLastBit(state>>1, array, 2*n+1)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *TriStateBitSet) putLastBit(state TriState, array *bitarray.BitArray, i int) error {
	bit := int(state & lastBitMask)
	_, err := array.Put(i, bit)
	return err
}

func (dbs *TriStateBitSet) bucketToArray(buckets []*BitSetCell) error {
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

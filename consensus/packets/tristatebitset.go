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

const lastBitMask = 0x00000001

// TriStateBitSet bitset implementation.
type TriStateBitSet struct {
	CompressedSet bool

	cells  []BitSetCell
	mapper BitSetMapper
}

// NewTriStateBitSet creates and returns a tristatebitset.
func NewTriStateBitSet(cells []*BitSetCell, mapper BitSetMapper) (*TriStateBitSet, error) {
	if (mapper == nil) || (cells == nil) {
		return nil, errors.New("[ NewTriStateBitSet ] failed to create tristatebitset")
	}
	bitset := &TriStateBitSet{
		cells:  make([]BitSetCell, mapper.Length()),
		mapper: mapper,
	}
	bitset.ApplyChanges(cells)
	return bitset, nil
}

func (dbs *TriStateBitSet) GetCells() []BitSetCell {
	return dbs.cells
}

func (dbs *TriStateBitSet) ApplyChanges(changes []*BitSetCell) {
	for _, cell := range changes {
		err := dbs.changeBucketState(cell)
		if err != nil {
			panic(err)
		}
	}
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
		return errors.Wrap(err, "[ changeBucketState ] failed to get index from ref")
	}
	dbs.cells[n] = *cell
	return nil
}

func putLastBit(array *bitarray.BitArray, state TriState, i int) error {
	bit := int(state & lastBitMask)
	_, err := array.Put(i, bit)
	return errors.Wrap(err, "[ putLastBit ] failed to put a bit ti bitset")
}

func changeBitState(array *bitarray.BitArray, i int, state TriState) error {
	err := putLastBit(array, state, 2*i)
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to put last bit")
	}
	err = putLastBit(array, state>>1, 2*i+1)
	if err != nil {
		return errors.Wrap(err, "[ changeBitState ] failed to put last bit")
	}
	return nil
}

func (dbs *TriStateBitSet) bucketToArray() (*bitarray.BitArray, error) {
	array := bitarray.New(dbs.mapper.Length() * 2) // cuz stores 2 bits for 1 id
	for _, bucket := range dbs.cells {
		n, err := dbs.mapper.RefToIndex(bucket.NodeID)
		if err != nil {
			return nil, errors.Wrap(err, "[ bucketToArray ] failed to get index from ref")
		}
		err = changeBitState(array, n, bucket.State)
		if err != nil {
			return nil, errors.Wrap(err, "[ bucketToArray ] failed to change bit state")
		}
	}
	return array, nil
}

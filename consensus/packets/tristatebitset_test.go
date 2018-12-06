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
	"bytes"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

const refsCount = 70

func initRefs() []core.RecordRef {
	result := make([]core.RecordRef, refsCount)
	for i := 0; i < refsCount; i++ {
		result[i] = testutils.RandomRef()
	}
	return result
}

func initBitCells(refs []core.RecordRef) []BitSetCell {
	result := make([]BitSetCell, refsCount)
	for i, ref := range refs {
		result[i] = BitSetCell{NodeID: ref, State: TimedOut}
	}
	return result
}

type BitSetMapperMock struct {
	refs []core.RecordRef
}

func (bsmm *BitSetMapperMock) IndexToRef(index int) (core.RecordRef, error) {
	if index > refsCount-1 {
		return testutils.RandomRef(), ErrBitSetOutOfRange
	}
	return bsmm.refs[index], nil
}

func (bsmm *BitSetMapperMock) RefToIndex(nodeID core.RecordRef) (int, error) {
	for i, ref := range bsmm.refs {
		if ref == nodeID {
			return i, nil
		}
	}
	return 0, ErrBitSetNodeIsMissing
}

func (bsmm *BitSetMapperMock) Length() int {
	return refsCount
}

func TestNewTriStateBitSet(t *testing.T) {
	_, err := NewBitSet(refsCount)
	assert.NoError(t, err)
}

func TestTriStateBitSet_GetBuckets(t *testing.T) {
	refs := initRefs()
	cells := initBitCells(refs)

	mapper := &BitSetMapperMock{refs: refs}
	bitset, _ := NewBitSet(len(cells))
	bitset.ApplyChanges(cells, mapper)
	newCells, err := bitset.GetCells(mapper)
	assert.NoError(t, err)
	assert.Equal(t, cells, newCells)
}

func TestTriStateBitSet_ApplyChanges(t *testing.T) {
	refs := initRefs()
	cells := initBitCells(refs)

	bitset, _ := NewBitSet(len(cells))

	mapper := &BitSetMapperMock{refs: refs}
	cells[refsCount-3].State = Fraud
	err := bitset.ApplyChanges(cells, mapper)
	assert.NoError(t, err)
	newCells1, err := bitset.GetCells(&BitSetMapperMock{refs: refs})
	assert.NoError(t, err)
	assert.Equal(t, cells, newCells1)
	cells[refsCount-4].State = Legit
	newCells2, err := bitset.GetCells(&BitSetMapperMock{refs: refs})
	assert.NoError(t, err)
	assert.NotEqual(t, cells, newCells2)
}

func TestBitSet(t *testing.T) {
	refs := initRefs()
	cells := initBitCells(refs)

	mapper := &BitSetMapperMock{refs: refs}
	bitset, _ := NewTriStateBitSet(len(cells))

	expected := bitset.array
	cells[refsCount-3].State = Fraud
	err := bitset.ApplyChanges(cells, mapper)
	assert.NoError(t, err)

	bitset2, _ := NewTriStateBitSet(len(cells))

	actual := bitset2.array
	cells[refsCount-3].State = Fraud
	err = bitset2.ApplyChanges(cells, mapper)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestBitArray(t *testing.T) {
	array := newBitArray(181)

	expected := uint8(1)
	err := array.put(1, 126)
	assert.NoError(t, err)
	err = array.put(1, 127)
	assert.NoError(t, err)
	err = array.put(1, 125)
	assert.NoError(t, err)
	err = array.put(0, 126)
	assert.NoError(t, err)
	actual, err := array.get(125)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	actual, err = array.get(127)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	actual, err = array.get(126)
	assert.NoError(t, err)
	expected = 0
	assert.Equal(t, expected, actual)
}

func TestTriStateBitSet_Serialize(t *testing.T) {
	refs := initRefs()
	cells := initBitCells(refs)

	mapper := &BitSetMapperMock{refs: refs}

	bitset, _ := NewTriStateBitSet(len(cells))
	bitset.ApplyChanges(cells, mapper)
	data, err := bitset.Serialize()
	assert.NoError(t, err)

	parsedBitSet, err := NewBitSet(len(cells))
	assert.NoError(t, err)
	parsedBitSet, err = DeserializeBitSet(bytes.NewReader(data))
	assert.NoError(t, err)

	expected, err := bitset.GetCells(mapper)
	assert.NoError(t, err)
	actual, err := parsedBitSet.GetCells(mapper)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTriStateBitSet_SerializeCompressed(t *testing.T) {
	refs := initRefs()
	cells := initBitCells(refs)

	mapper := &BitSetMapperMock{refs: refs}

	bitset, _ := NewTriStateBitSet(len(cells))
	bitset.CompressedSet = true
	err := bitset.ApplyChanges(cells, mapper)
	assert.NoError(t, err)
	data, err := bitset.Serialize()
	assert.NoError(t, err)

	parsedBitSet, err := NewBitSet(len(cells))
	assert.NoError(t, err)
	parsedBitSet, err = DeserializeBitSet(bytes.NewReader(data))
	assert.NoError(t, err)

	expected, err := bitset.GetCells(mapper)
	assert.NoError(t, err)
	actual, err := parsedBitSet.GetCells(mapper)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

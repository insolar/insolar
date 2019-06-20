//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package packets

import (
	"bytes"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func initRefs(count int) []insolar.Reference {
	result := make([]insolar.Reference, count)
	for i := 0; i < count; i++ {
		result[i] = testutils.RandomRef()
	}
	return result
}

func initBitCells(refs []insolar.Reference) []BitSetCell {
	result := make([]BitSetCell, len(refs))
	for i, ref := range refs {
		result[i] = BitSetCell{NodeID: ref, State: TimedOut}
	}
	return result
}

func initDiffBitCells(refs []insolar.Reference) []BitSetCell {
	result := make([]BitSetCell, len(refs))
	for i := 0; i < len(refs)/2+1; i++ {
		result[i] = BitSetCell{NodeID: refs[i], State: TimedOut}
	}
	for i := len(refs)/2 + 1; i < len(refs); i++ {
		result[i] = BitSetCell{NodeID: refs[i], State: Legit}
	}
	return result
}

type mapperMock struct {
	refs []insolar.Reference
}

func (bsmm *mapperMock) IndexToRef(index int) (insolar.Reference, error) {
	if index > len(bsmm.refs)-1 {
		return testutils.RandomRef(), ErrBitSetOutOfRange
	}
	return bsmm.refs[index], nil
}

func (bsmm *mapperMock) RefToIndex(nodeID insolar.Reference) (int, error) {
	for i, ref := range bsmm.refs {
		if ref == nodeID {
			return i, nil
		}
	}
	return 0, ErrBitSetNodeIsMissing
}

func (bsmm *mapperMock) Length() int {
	return len(bsmm.refs)
}

func TestBitSet_GetBuckets(t *testing.T) {
	count := 70
	refs := initRefs(count)
	cells := initBitCells(refs)

	mapper := &mapperMock{refs: refs}
	bitset, _ := NewBitSet(len(cells))
	err := bitset.ApplyChanges(cells, mapper)
	assert.NoError(t, err)
	newCells, err := bitset.GetCells(mapper)
	assert.NoError(t, err)
	assert.Equal(t, cells, newCells)
}

func TestBitSet_ApplyChanges(t *testing.T) {
	count := 65
	refs := initRefs(count)
	cells := initBitCells(refs)

	bitset, _ := NewBitSet(len(cells))

	mapper := &mapperMock{refs: refs}
	cells[count-3].State = Fraud
	err := bitset.ApplyChanges(cells, mapper)
	assert.NoError(t, err)
	newCells1, err := bitset.GetCells(&mapperMock{refs: refs})
	assert.NoError(t, err)
	assert.Equal(t, cells, newCells1)
	cells[count-4].State = Legit
	newCells2, err := bitset.GetCells(&mapperMock{refs: refs})
	assert.NoError(t, err)
	assert.NotEqual(t, cells, newCells2)
}

func Test_parseStatesFromByte(t *testing.T) {
	data := make([]BitSetState, statesInByte)
	data[0] = Fraud
	data[1] = Inconsistent
	data[2] = TimedOut
	data[3] = Legit

	b := writeStatesToByte(data)
	data2 := parseStatesFromByte(b)
	assert.Equal(t, data, data2[:])
}

func Test_parseStatesFromByte2(t *testing.T) {
	data := make([]BitSetState, 3)
	data[0] = Fraud
	data[1] = TimedOut
	data[2] = Legit

	b := writeStatesToByte(data)
	data2 := parseStatesFromByte(b)
	assert.Equal(t, data, data2[:3])
}

func TestBitArray_Serialize(t *testing.T) {
	ba := make(bitArray, 25)
	for i := 0; i < len(ba); i++ {
		ba[i] = BitSetState(i % statesInByte)
	}
	testSerializeBitarray(t, ba)
}

func TestBitArray_Serialize2(t *testing.T) {
	ba := make(bitArray, 25)
	ba[16] = Legit
	testSerializeBitarray(t, ba)
}

func TestBitArray_SerializeCompressed(t *testing.T) {
	ba := make(bitArray, 25)
	for i := 0; i < len(ba); i++ {
		ba[i] = BitSetState(i % statesInByte)
	}
	testSerializeBitarray(t, ba)
}

func TestBitArray_SerializeCompressed2(t *testing.T) {
	ba := make(bitArray, 25)
	ba[16] = Legit
	testSerializeBitarray(t, ba)
}

func testSerializeBitarray(t *testing.T, ba bitArray) {
	data, err := ba.serializeCompressed()
	assert.NoError(t, err)
	ba2, err := deserializeCompressed(bytes.NewReader(data), len(ba))
	assert.NoError(t, err)
	assert.Equal(t, ba, ba2)
}

func TestBitSet_Serialize(t *testing.T) {
	refs := initRefs(92)
	cells := initBitCells(refs)

	testSerializeDeserialize(t, refs, cells, false)
}

func TestBitSet_SerializeCompressed(t *testing.T) {
	refs := initRefs(44)
	cells := initBitCells(refs)

	testSerializeDeserialize(t, refs, cells, true)
}

func TestBitSet_ThousandStates_Compressed(t *testing.T) {
	refs := initRefs(1024)
	cells := initBitCells(refs)

	testSerializeDeserialize(t, refs, cells, true)
}

func TestBitSet_ThousandDiffStates_Compressed(t *testing.T) {
	refs := initRefs(1024)
	cells := initDiffBitCells(refs)

	testSerializeDeserialize(t, refs, cells, true)
}

func TestBitSet_ThousandStates(t *testing.T) {
	refs := initRefs(1024)
	cells := initBitCells(refs)

	testSerializeDeserialize(t, refs, cells, false)
}

func TestBitSet_ThousandDiffStates(t *testing.T) {
	refs := initRefs(1024)
	cells := initDiffBitCells(refs)

	testSerializeDeserialize(t, refs, cells, false)
}

func testSerializeDeserialize(t *testing.T, refs []insolar.Reference, cells []BitSetCell, compressed bool) {
	mapper := &mapperMock{refs: refs}

	bitset, _ := NewBitSetImpl(len(cells), compressed)
	err := bitset.ApplyChanges(cells, mapper)
	assert.NoError(t, err)
	data, err := bitset.Serialize()
	assert.NoError(t, err)

	parsedBitSet, err := DeserializeBitSet(bytes.NewReader(data))
	assert.NoError(t, err)

	expected, err := bitset.GetCells(mapper)
	assert.NoError(t, err)
	actual, err := parsedBitSet.GetCells(mapper)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package jet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

// helper for tests
func dbTreeForPulse(s *DBStore, pulse insolar.PulseNumber) *Tree {
	store := s
	serializedTree, err := store.db.Get(pulseKey(pulse))
	if err != nil {
		return nil
	}

	recovered := &Tree{}
	err = insolar.Deserialize(serializedTree, recovered)
	if err != nil {
		return nil
	}
	return recovered
}

func TestDBStorage_Empty(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	s := NewDBStore(db)

	all := s.All(ctx, insolar.FirstPulseNumber)
	require.Equal(t, 1, len(all), "should be just one jet ID")
	require.Equal(t, insolar.ZeroJetID, all[0], "JetID should be a zero on empty storage")
}

func TestDBStorage_UpdateJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	s := NewDBStore(db)

	var (
		expected = []insolar.JetID{insolar.ZeroJetID}
	)

	err := s.Update(ctx, 100, true, *insolar.NewJetID(0, nil))
	require.NoError(t, err)

	tree := dbTreeForPulse(s, 100)
	require.Equal(t, expected, tree.LeafIDs(), "actual tree in string form: %v", tree.String())
}

func TestDBStorage_SplitJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	s := NewDBStore(db)

	var (
		expectedLeft  = insolar.JetID{0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		expectedRight = insolar.JetID{0, 0, 0, 1, 1, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		expectedLeafs = Tree{Head: &jet{
			Actual: false,
			Left:   &jet{Actual: false},
			Right:  &jet{Actual: false},
		}}
	)

	root := insolar.NewJetID(0, nil)
	left, right, err := s.Split(ctx, 100, *root)
	require.NoError(t, err)
	assert.Equal(t, insolar.ZeroJetID, *root, "actual tree node in string form: %v", root.DebugString())
	assert.Equal(t, expectedLeft, left, "actual tree node in string form: %v", left.DebugString())
	assert.Equal(t, expectedRight, right, "actual tree node in string form: %v", right.DebugString())

	tree := dbTreeForPulse(s, 100)
	require.Equal(t, expectedLeafs, *tree, "actual tree in string form: %v", tree.String())
}

func TestDBStorage_CloneJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	s := NewDBStore(db)

	var (
		expectedZero = []insolar.JetID{insolar.ZeroJetID}
		expectedNil  []insolar.JetID
	)

	err := s.Update(ctx, 100, true, *insolar.NewJetID(0, nil))
	require.NoError(t, err)

	tree := dbTreeForPulse(s, 100)
	assert.Equal(t, expectedZero, tree.LeafIDs(), "actual tree in string form: %v", tree.String())

	err = s.Clone(ctx, 100, 101)
	require.NoError(t, err)

	tree = dbTreeForPulse(s, 101)
	assert.Equal(t, expectedNil, tree.LeafIDs(), "actual tree in string form: %v", tree.String())

	tree = dbTreeForPulse(s, 100)
	assert.Equal(t, expectedZero, tree.LeafIDs(), "actual tree in string form: %v", tree.String())
}

func TestDBStorage_ForID_Basic(t *testing.T) {
	ctx := inslogger.TestContext(t)

	pn := gen.PulseNumber()
	meaningfulBits := "01000011" + "11000011" + "010010"

	bits := parsePrefix(meaningfulBits)
	expectJetID := NewIDFromString(meaningfulBits)
	searchID := gen.ID()
	hash := searchID[insolar.RecordHashOffset:]
	hash = setBitsPrefix(hash, bits, len(meaningfulBits))
	copy(searchID[insolar.RecordHashOffset:], hash)

	for _, actuality := range []bool{true, false} {
		db := store.NewMemoryMockDB()
		s := NewDBStore(db)
		s.Update(ctx, pn, actuality, expectJetID)
		found, ok := s.ForID(ctx, pn, searchID)
		require.Equal(t, expectJetID, found, "got jet with exactly same prefix")
		require.Equal(t, actuality, ok, "jet should be in actuality state we defined in Update")
	}
}

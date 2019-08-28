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
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulse"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

// helper for tests
func dbTreeForPulse(s *DBStore, pulse insolar.PulseNumber) *Tree {
	store := s
	serializedTree, err := store.db.Get(pulseKey(pulse))
	if err != nil {
		return nil
	}

	recovered := &Tree{}
	err = recovered.Unmarshal(serializedTree)
	if err != nil {
		return nil
	}
	return recovered
}

func TestPulseKey(t *testing.T) {
	t.Parallel()

	expectedKey := pulseKey(insolar.GenesisPulse.PulseNumber)

	rawID := expectedKey.ID()

	actualKey := newPulseKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

func TestDBStore_TruncateHead(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	dbMock, err := store.NewBadgerDB(ops)

	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	dbStore := NewDBStore(dbMock)

	numElements := 10

	// it's used for writing pulses in random order to db
	indexes := make([]int, numElements)
	for i := 0; i < numElements; i++ {
		indexes[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(indexes), func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })

	startPulseNumber := insolar.GenesisPulse.PulseNumber
	for _, idx := range indexes {
		pulse := startPulseNumber + insolar.PulseNumber(idx)
		jetTree := NewTree(true)
		err := dbStore.set(pulse, jetTree)
		require.NoError(t, err)
	}

	for i := 0; i < numElements; i++ {
		tree := dbStore.get(startPulseNumber + insolar.PulseNumber(i))
		require.True(t, tree.Head.Actual)
	}

	numLeftElements := numElements / 2
	err = dbStore.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		tree := dbStore.get(startPulseNumber + insolar.PulseNumber(i))
		require.True(t, tree.Head.Actual)
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		tree := dbStore.get(startPulseNumber + insolar.PulseNumber(i))
		require.False(t, tree.Head.Actual)
	}

	// not existing record
	err = dbStore.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements+numElements*2))
	require.NoError(t, err)
}

func TestDBStorage_Empty(t *testing.T) {
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(ctx)
	s := NewDBStore(db)

	all := s.All(ctx, pulse.MinTimePulse)
	require.Equal(t, 1, len(all), "should be just one jet ID")
	require.Equal(t, insolar.ZeroJetID, all[0], "JetID should be a zero on empty storage")
}

func TestDBStorage_UpdateJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(ctx)
	s := NewDBStore(db)

	var (
		expected = []insolar.JetID{insolar.ZeroJetID}
	)

	err = s.Update(ctx, 100, true, *insolar.NewJetID(0, nil))
	require.NoError(t, err)

	tree := dbTreeForPulse(s, 100)
	require.Equal(t, expected, tree.LeafIDs(), "actual tree in string form: %v", tree.String())
}

func TestDBStorage_SplitJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(ctx)
	s := NewDBStore(db)

	lArray := []byte{0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	rArray := []byte{0, 0, 0, 1, 1, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var (
		expectedLeft  = insolar.JetID(*insolar.NewIDFromBytes(lArray))
		expectedRight = insolar.JetID(*insolar.NewIDFromBytes(rArray))
		expectedLeafs = Tree{Head: &Jet{
			Actual: false,
			Left:   &Jet{Actual: true},
			Right:  &Jet{Actual: true},
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

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(ctx)
	require.NoError(t, err)
	s := NewDBStore(db)

	var (
		expectedZero = []insolar.JetID{insolar.ZeroJetID}
		expectedNil  []insolar.JetID
	)

	err = s.Update(ctx, 100, true, *insolar.NewJetID(0, nil))
	require.NoError(t, err)

	tree := dbTreeForPulse(s, 100)
	assert.Equal(t, expectedZero, tree.LeafIDs(), "actual tree in string form: %v", tree.String())

	err = s.Clone(ctx, 100, 101, false)
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
	hash := searchID.Hash()
	hash = setBitsPrefix(hash, bits, len(meaningfulBits))
	searchID = *insolar.NewID(searchID.Pulse(), hash)

	for _, actuality := range []bool{true, false} {
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(ctx)
		s := NewDBStore(db)
		s.Update(ctx, pn, actuality, expectJetID)
		found, ok := s.ForID(ctx, pn, searchID)
		require.Equal(t, expectJetID, found, "got jet with exactly same prefix")
		require.Equal(t, actuality, ok, "jet should be in actuality state we defined in Update")
	}
}

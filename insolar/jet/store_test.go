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

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulse"
)

// helper for tests
func treeForPulse(s *Store, pulse insolar.PulseNumber) (*Tree, bool) {
	ltree, ok := s.trees[pulse]
	if !ok {
		return nil, false
	}
	return ltree.t, true
}

func TestJetStorage_Empty(t *testing.T) {
	ctx := inslogger.TestContext(t)
	s := NewStore()

	all := s.All(ctx, pulse.MinTimePulse)
	require.Equal(t, 1, len(all), "should be just one jet ID")
	require.Equal(t, insolar.ZeroJetID, all[0], "JetID should be a zero on empty storage")
}

func TestJetStorage_UpdateJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	s := NewStore()

	s.Update(ctx, 100, true, *insolar.NewJetID(0, nil))

	tree, _ := treeForPulse(s, 100)
	require.Equal(t, "root (level=0 actual=true)\n", tree.String())
}

func TestJetStorage_SplitJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	s := NewStore()

	zeroJet := insolar.ZeroJetID
	left, right, err := s.Split(ctx, 100, insolar.ZeroJetID)
	require.NoError(t, err)
	require.Equal(t, "[JET 0 -]", zeroJet.DebugString())
	require.Equal(t, "[JET 1 0]", left.DebugString())
	require.Equal(t, "[JET 1 1]", right.DebugString())

	tree, _ := treeForPulse(s, 100)
	require.Equal(t, "root (level=0 actual=false)\n 0 (level=1 actual=true)\n 1 (level=1 actual=true)\n", tree.String())
}

func TestJetStorage_CloneJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	s := NewStore()

	s.Update(ctx, 100, true, *insolar.NewJetID(0, nil))

	tree, _ := treeForPulse(s, 100)
	require.Equal(t, "root (level=0 actual=true)\n", tree.String())

	s.Clone(ctx, 100, 101, false)

	tree, _ = treeForPulse(s, 101)
	require.Equal(t, "root (level=0 actual=false)\n", tree.String())

	tree, _ = treeForPulse(s, 100)
	require.Equal(t, "root (level=0 actual=true)\n", tree.String())
}

func TestJetStorage_DeleteJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	s := NewStore()

	_, _, err := s.Split(ctx, 100, *insolar.NewJetID(0, nil))
	require.NoError(t, err)

	s.DeleteForPN(ctx, 100)

	_, ok := treeForPulse(s, 100)
	require.False(t, ok, "tree should be an empty")

	all := s.All(ctx, 100)
	require.Equal(t, 0, len(all), "should be just one jet ID")
}

func TestJetStorage_ForID_Basic(t *testing.T) {
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
		s := NewStore()
		s.Update(ctx, pn, actuality, expectJetID)
		found, ok := s.ForID(ctx, pn, searchID)
		require.Equal(t, expectJetID, found, "got jet with exactly same prefix")
		require.Equal(t, actuality, ok, "jet should be in actuality state we defined in Update")
	}
}

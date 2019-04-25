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
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
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

	all := s.All(ctx, gen.PulseNumber())
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

	left, right, err := s.Split(ctx, 100, *insolar.NewJetID(0, nil))
	require.NoError(t, err)
	require.Equal(t, "[JET 1 0]", left.DebugString())
	require.Equal(t, "[JET 1 1]", right.DebugString())

	tree, _ := treeForPulse(s, 100)
	require.Equal(t, "root (level=0 actual=false)\n 0 (level=1 actual=false)\n 1 (level=1 actual=false)\n", tree.String())
}

func TestJetStorage_CloneJetTree(t *testing.T) {
	ctx := inslogger.TestContext(t)
	s := NewStore()

	s.Update(ctx, 100, true, *insolar.NewJetID(0, nil))

	tree, _ := treeForPulse(s, 100)
	require.Equal(t, "root (level=0 actual=true)\n", tree.String())

	s.Clone(ctx, 100, 101)

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
	require.Equal(t, 1, len(all), "should be just one jet ID")
	require.Equal(t, insolar.ZeroJetID, all[0], "JetID should be a zero after tree removal")
}

func TestJetStorage_ForID_Basic(t *testing.T) {
	ctx := inslogger.TestContext(t)

	pn := gen.PulseNumber()
	meaningfulBits := "01000011" + "11000011" + "010010"

	bits := parsePrefix(meaningfulBits)
	expectJetID := NewIDFromString(meaningfulBits)
	// fmt.Printf("expectJetID:        %08b\n", expectJetID[:])
	searchID := gen.ID()
	hash := searchID[insolar.RecordHashOffset:]
	hash = setBitsPrefix(hash, bits, len(meaningfulBits))
	copy(searchID[insolar.RecordHashOffset:], hash)

	for _, actuality := range []bool{true, false} {
		s := NewStore()
		s.Update(ctx, pn, actuality, expectJetID)
		found, ok := s.ForID(ctx, pn, searchID)
		require.Equal(t, expectJetID, found, "got jet with exactly same prefix")
		require.Equal(t, actuality, ok, "jet should be in actuality state we defined in Update")
	}
}

func TestJetStorage_ForID_Fuzz(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pn := gen.PulseNumber()

	var jets = map[string]struct{}{}

	// findBestMatch returns best match JetID for provided recordID, searches JetID in `jets` set.
	findBestMatch := func(id insolar.ID) insolar.JetID {
		// search in substrings
		idbits := bitsToString(id.Hash())
		var stat = map[string]int{}
		for jetbits := range jets {
			stat[jetbits] = 0
			for i, b := range []byte(idbits) {
				if i >= len(jetbits) {
					break
				}
				if jetbits[i] != b {
					break
				}
				stat[jetbits]++
			}
		}

		var max int
		var found string
		for jetbits, matchlen := range stat {
			if matchlen > max {
				max = matchlen
				found = jetbits
			}
		}
		return NewIDFromString(found)
	}

	var searches []insolar.ID
	s := NewStore()
	// generate jet IDs, add them to jet store (actually to underlying jet tree)
	// fill searches list with ID with hashes what should match with generated Jet ID
	for i := 0; i < 100; i++ {
		jetID := gen.JetID()
		prefix, depth := jetID.Prefix(), jetID.Depth()

		searchID := gen.ID()
		hash := setBitsPrefix(searchID[insolar.RecordHashOffset:], prefix, int(depth))
		copy(searchID[insolar.RecordHashOffset:], hash)

		s.Update(ctx, pn, false, jetID)
		searches = append(searches, searchID)
	}
	// fill jets set with jets saved in store.
	for _, j := range s.All(ctx, pn) {
		prefix, depth := j.Prefix(), j.Depth()
		s := bitsToString(prefix)
		jets[s[:depth]] = struct{}{}
	}

	// check is ID match proper JetID
	for _, searchID := range searches {
		found, _ := s.ForID(ctx, pn, searchID)
		expect := findBestMatch(searchID)
		assertResult := assert.Equalf(t,
			expect, found,
			" expect  = %08b\n got     = %08b\n id hash = %08b\n",
			expect[insolar.JetPrefixOffset:],
			found[insolar.JetPrefixOffset:],
			searchID.Hash(),
		)
		if !assertResult {
			// print more info on fail for easy investigation what happened.
			fmt.Println("all jets:")
			var sorted []string
			for jetbits := range jets {
				sorted = append(sorted, jetbits)
			}
			sort.Strings(sorted)
			for _, j := range sorted {
				fmt.Printf("  %v\n", j)
			}
			t.Fail()
			return
		}
	}
}

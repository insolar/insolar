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

package executor

import (
	"bytes"
	"context"
	"log"
	"sort"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJetSplitter(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jc := jet.NewCoordinatorMock(t)
	// use real jet store
	js := jet.NewStore()
	da := drop.NewAccessorMock(t)
	splitter := NewJetSplitter(jc, js, js, da)
	require.NotNil(t, splitter, "jet splitter created")

	splitter.splitCount = 0

	me := gen.Reference()
	jc.MeMock.Return(me)

	// hasSplitIntention
	pn := gen.PulseNumber()
	// just avoid special pulses
	if pn < 60000 {
		pn += 60000
	}
	var (
		previous = pn
		current  = pn + 1
		newpulse = pn + 2
	)

	var splitID insolar.JetID
	da.ForPulseFunc = func(_ context.Context, jetID insolar.JetID, pn insolar.PulseNumber) (r drop.Drop, r1 error) {
		if pn != previous {
			return drop.Drop{}, errors.Errorf("unexpected pulse number %v, expects previous pulse %v", pn, previous)
		}
		split := false
		if splitID == jetID {
			split = true
		}
		return drop.Drop{Split: split}, nil
	}

	jets := []insolar.JetID{
		jet.NewIDFromString("0"),
		jet.NewIDFromString("10"),
		jet.NewIDFromString("11"),
	}

	// no filter for ID
	jc.LightExecutorForJetMock.Return(&me, nil)

	// update real jet store
	js.Update(ctx, current, true, jets...)

	t.Run("no_split", func(t *testing.T) {
		res, err := splitter.Do(ctx, previous, current, newpulse)
		require.NoError(t, err, "splitter method Do error check")
		require.Equal(t, len(jets), len(res), "compare jets count")
		var gotJets []insolar.JetID
		for _, info := range res {
			gotJets = append(gotJets, info.ID)
			assert.False(t, info.SplitPerformed, "split is not performed")
		}
		require.Equal(t, jsort(jets), jsort(gotJets), "compare results")
	})

	t.Run("with_split_intent", func(t *testing.T) {
		splitID = jet.NewIDFromString("11")

		res, err := splitter.Do(ctx, previous, current, newpulse)
		require.NoError(t, err, "splitter method Do error check")
		require.Equal(t, len(jets), len(res), "compare jets count")
		var gotJets []insolar.JetID
		for _, info := range res {
			gotJets = append(gotJets, info.ID)
			assert.False(t, info.SplitIntent, "no split")

			if info.ID != splitID {
				assert.False(t, info.SplitPerformed, "split is not performed")
			} else {
				assert.True(t, info.SplitPerformed, "split is performed")
			}

		}
		require.Equal(t, jsort(jets), jsort(gotJets), "compare results")
	})
}

func jsort(jets []insolar.JetID) []string {
	sort.Slice(jets, func(i, j int) bool {
		switch bytes.Compare(jets[i][:], jets[j][:]) {
		case -1:
			return true
		case 0, 1:
			return false
		default:
			log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
			return false
		}
	})
	result := make([]string, 0, len(jets))
	for _, j := range jets {
		result = append(result, j.DebugString())
	}
	return result
}

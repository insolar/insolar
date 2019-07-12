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
	"sort"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestJetSplitter(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jc := NewJetCalculatorMock(t)
	// use real jet store
	js := jet.NewStore()
	da := drop.NewAccessorMock(t)
	dm := drop.NewModifierMock(t)
	pc := pulse.NewCalculatorMock(t)
	splitter := NewJetSplitter(jc, js, js, da, dm, pc)
	require.NotNil(t, splitter, "jet splitter created")

	splitter.splitsLimit = 0

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

	pc.BackwardsMock.Return(insolar.Pulse{PulseNumber: previous}, nil)

	// beware splitID is a shared variable (check code below)
	var splitID insolar.JetID

	da.ForPulseFunc = func(_ context.Context, jetID insolar.JetID, pn insolar.PulseNumber) (drop.Drop, error) {
		if pn != previous {
			return drop.Drop{}, errors.Errorf("unexpected pulse number %v, expects previous pulse %v", pn, previous)
		}
		split := false
		if splitID == jetID {
			split = true
		}
		return drop.Drop{Split: split}, nil
	}
	dm.SetFunc = func(_ context.Context, d drop.Drop) error {
		if d.Pulse != current {
			return errors.Errorf("unexpected pulse number %v, expects current pulse %v", d.Pulse, current)
		}
		require.False(t, d.Split, "we have no jets with split intent although splitsLimit set to 0")
		return nil
	}

	jets := []insolar.JetID{
		jet.NewIDFromString("0"),
		jet.NewIDFromString("10"),
		jet.NewIDFromString("11"),
	}

	// no filter for ID
	jc.MineForPulseFunc = func(_ context.Context, pn insolar.PulseNumber) []insolar.JetID {
		return jets
	}

	// update real jet store
	err := js.Update(ctx, current, true, jets...)
	require.NoError(t, err, "jet store update")

	t.Run("no_split", func(t *testing.T) {
		gotJets, err := splitter.Do(ctx, current, newpulse)
		require.NoError(t, err, "splitter method Do error check")
		require.Equal(t, len(jets), len(gotJets), "compare jets count")
		require.Equal(t, jsort(jets), jsort(gotJets), "no splits")
	})

	t.Run("with_split_intent", func(t *testing.T) {
		splitID = jet.NewIDFromString("11")

		gotJets, err := splitter.Do(ctx, current, newpulse)
		require.NoError(t, err, "splitter method Do error check")
		require.Equal(t, len(jets)+1, len(gotJets), "compare jets count, expect one split")
		expectJets := make([]insolar.JetID, 0, len(jets))
		split0, split1 := jet.NewIDFromString("110"), jet.NewIDFromString("111")
		for _, id := range jets {
			if id == splitID {
				expectJets = append(expectJets, split0, split1)
				continue
			}
			expectJets = append(expectJets, id)
		}
		require.Equalf(t, jsort(expectJets), jsort(gotJets), "split %v is split to %v and %v",
			splitID.DebugString(), split0.DebugString(), split1.DebugString(),
		)
	})
}

func jsort(jets []insolar.JetID) []string {
	sort.Slice(jets, func(i, j int) bool {
		return bytes.Compare(jets[i][:], jets[j][:]) == -1
	})
	result := make([]string, 0, len(jets))
	for _, j := range jets {
		result = append(result, j.DebugString())
	}
	return result
}

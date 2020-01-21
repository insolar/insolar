// Copyright 2020 Insolar Network Ltd.
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

package executor

import (
	"context"
	"fmt"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func TestJetCalculator_New(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jc := jet.NewCoordinatorMock(t)
	js := jet.NewStorageMock(t)
	jetCalculator := NewJetCalculator(jc, js)
	require.NotNil(t, jetCalculator, "jet splitter created")

	me := gen.Reference()
	pn := gen.PulseNumber()
	jc.MeMock.Return(me)

	jc.LightExecutorForJetMock.Set(func(_ context.Context, jetID insolar.ID, p insolar.PulseNumber) (r *insolar.Reference, r1 error) {
		if p != pn {
			panic(fmt.Sprintf("pulse number %v is unexpected", p))
		}
		return &me, nil
	})

	var allJets []insolar.JetID
	js.AllMock.Set(func(_ context.Context, p insolar.PulseNumber) []insolar.JetID {
		if p == pn {
			return allJets
		}
		return nil
	})

	t.Run("empty_case", func(t *testing.T) {
		jets, err := jetCalculator.MineForPulse(ctx, pn)
		require.NoError(t, err)
		require.Nil(t, jets, "MineForPulse returns empty set of jets")
	})

	t.Run("one_element", func(t *testing.T) {
		allJets = []insolar.JetID{
			gen.JetID(),
		}
		jets, err := jetCalculator.MineForPulse(ctx, pn)
		require.NoError(t, err)
		require.NotNil(t, jets, "MineForPulse returns not empty set of jets")
		require.Equal(t, len(allJets), len(jets), "MineForPulse compare return values count")
		require.Equal(t, allJets, jets, "MineForPulse compare return values")
	})

	t.Run("multiple_elements", func(t *testing.T) {
		allJets = []insolar.JetID{
			jet.NewIDFromString("0"),
			jet.NewIDFromString("01"),
			jet.NewIDFromString("011"),
		}
		jets, err := jetCalculator.MineForPulse(ctx, pn)
		require.NoError(t, err)
		require.NotNil(t, jets, "MineForPulse returns not empty set of jets")
		require.Equal(t, len(allJets), len(jets), "MineForPulse compare return values count")
		require.Equal(t, allJets, jets, "MineForPulse compare return values")
	})

	t.Run("no nodes returns error", func(t *testing.T) {
		allJets = []insolar.JetID{
			jet.NewIDFromString("0"),
			jet.NewIDFromString("01"),
			jet.NewIDFromString("011"),
		}

		jc.LightExecutorForJetMock.Set(func(_ context.Context, jetID insolar.ID, p insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			if p != pn {
				panic(fmt.Sprintf("pulse number %v is unexpected", p))
			}
			if insolar.JetID(jetID) == allJets[1] {
				return nil, node.ErrNoNodes
			}
			return &me, nil
		})
		_, err := jetCalculator.MineForPulse(ctx, pn)
		require.Error(t, err)
	})
}

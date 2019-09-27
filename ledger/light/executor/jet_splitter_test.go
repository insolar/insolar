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
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type splitCase struct {
	name string
	cfg  configuration.JetSplit
	// represents expected values on every pulse for every jet
	pulses []map[insolar.JetID]jetConfig
}

type jetConfig struct {
	// how many record return record accessor
	records int
	// expected drop's threshold value
	dropThreshold int
	// is split expected
	hasSplit bool
}

var (
	initialDepth uint8 = 2
	// not limit depth by default
	defaultDepthLimit = initialDepth + 10
)

// initial jets
var (
	jet0  = jet.NewIDFromString("0")
	jet10 = jet.NewIDFromString("10")
	jet11 = jet.NewIDFromString("11")
)

// children jets
var (
	// left and right children for jet10
	jet10Left  = jet.NewIDFromString("100")
	jet10Right = jet.NewIDFromString("101")
)

var splitCases = []splitCase{
	{
		name: "no_split",
		cfg: configuration.JetSplit{
			ThresholdRecordsCount:  6,
			ThresholdOverflowCount: 0,
			DepthLimit:             defaultDepthLimit,
		},
		pulses: []map[insolar.JetID]jetConfig{
			{jet10: {5, 0, false}},
			{jet10: {3, 0, false}},
		},
	},
	{
		name: "split",
		cfg: configuration.JetSplit{
			ThresholdRecordsCount:  4,
			ThresholdOverflowCount: 0,
			DepthLimit:             defaultDepthLimit,
		},
		pulses: []map[insolar.JetID]jetConfig{
			{jet10: {5, 1, true}},
			{jet10Left: {3, 0, false}},
		},
	},
	{
		name: "split_with_overflow",
		cfg: configuration.JetSplit{
			ThresholdRecordsCount:  4,
			ThresholdOverflowCount: 1,
			DepthLimit:             defaultDepthLimit,
		},
		pulses: []map[insolar.JetID]jetConfig{
			{jet10: {5, 1, false}},
			{jet10: {5, 2, true}},
			{jet10Left: {5, 1, false}},
		},
	},
	{
		name: "no_split_with_overflow",
		cfg: configuration.JetSplit{
			ThresholdRecordsCount:  5,
			ThresholdOverflowCount: 1,
			DepthLimit:             defaultDepthLimit,
		},
		pulses: []map[insolar.JetID]jetConfig{
			{jet10: {5, 1, false}},
			{jet10: {4, 0, false}},
			{jet10: {5, 1, false}},
		},
	},
	{
		// expect here only one split has preformed
		name: "split_with_depth_limit",
		cfg: configuration.JetSplit{
			ThresholdRecordsCount:  4,
			ThresholdOverflowCount: 0,
			DepthLimit:             initialDepth + 1,
		},
		pulses: []map[insolar.JetID]jetConfig{
			{jet10: {5, 1, true}},
			{jet10Left: {5, 0, false}},
		},
	},
}

func TestJetSplitter(t *testing.T) {
	ctx := inslogger.TestContext(t)

	checkCase := func(t *testing.T, sc splitCase) {
		// real components
		jetStore := jet.NewStore()
		db := drop.NewStorageMemory()
		dropAccessor := db
		dropModifier := db
		// mocks
		jetCalc := NewJetCalculatorMock(t)
		collectionAccessor := object.NewRecordCollectionAccessorMock(t)
		pulseCalc := pulse.NewCalculatorMock(t)

		// create splitter
		splitter := NewJetSplitter(
			sc.cfg,
			jetCalc, jetStore, jetStore,
			dropAccessor, dropModifier,
			pulseCalc, collectionAccessor,
		)

		var initialPulse insolar.PulseNumber = 60000
		initialJets := []insolar.JetID{jet0, jet10, jet11}
		// initialize jet tree
		err := jetStore.Update(ctx, initialPulse, true, initialJets...)
		require.NoError(t, err, "jet store updated with initial jets")

		for i, jetsConfig := range sc.pulses {
			previous := initialPulse + insolar.PulseNumber(i) - 1
			ended := previous + 1
			newpulse := ended + 1

			pulseCalc.BackwardsMock.Return(insolar.Pulse{PulseNumber: previous}, nil)

			// jets state before possible split
			pulseStartedWithJets := jetStore.All(ctx, ended)

			collectionAccessor.ForPulseMock.Set(func(_ context.Context, jetID insolar.JetID, pn insolar.PulseNumber) []record.Material {
				jConf, ok := jetsConfig[jetID]
				if !ok {
					return nil
				}
				return make([]record.Material, jConf.records)
			})

			gotJets, err := splitter.Do(ctx, ended, newpulse, jetStore.All(ctx, ended), true)
			require.NoError(t, err, "splitter.Do performed")

			for jetID, jConf := range jetsConfig {
				require.Truef(t, jetInList(pulseStartedWithJets, jetID),
					"jet %v should be in jet-tree's leaves, got %v (+%v pulse)",
					jetID.DebugString(), jsort(pulseStartedWithJets), i)

				dropThreshold := splitter.getDropThreshold(ctx, jetID, ended)
				require.Equalf(t, jConf.dropThreshold, dropThreshold,
					"check drop.SplitThresholdExceeded for jet %v in +%v pulse", jetID.DebugString(), i)
			}

			var expectJets []insolar.JetID
			var splitJets []string
			for _, jetID := range pulseStartedWithJets {
				jConf := jetsConfig[jetID]
				if jConf.hasSplit {

					left, right := jet.Siblings(jetID)
					expectJets = append(expectJets, left, right)
					splitJets = append(splitJets, jetID.DebugString())
					continue
				}
				expectJets = append(expectJets, jetID)
			}
			jetsInfo := "jets should split " + strings.Join(splitJets, ", ")
			if len(splitJets) == 0 {
				jetsInfo = "no jets spit"
			}

			expectMsg := fmt.Sprintf("jet %v should split on +%v pulse", jetsInfo, i)
			require.Equal(t, jsort(expectJets), jsort(gotJets), expectMsg)

			for _, jetID := range pulseStartedWithJets {
				jConf := jetsConfig[jetID]
				block, err := dropAccessor.ForPulse(ctx, jetID, ended)
				require.NoErrorf(t, err,
					"should be drop for jet %v, on pulse +%v (%v)", jetID.DebugString(), i, ended)
				assert.Equalf(t, jConf.hasSplit, block.Split,
					"drop's split flag check for jet %v on pulse +%v", jetID.DebugString(), i)
			}
		}
	}

	for _, sc := range splitCases {
		t.Run(sc.name, func(t *testing.T) { checkCase(t, sc) })
	}
}

func jetInList(jets []insolar.JetID, jetID insolar.JetID) bool {
	for _, j := range jets {
		if j == jetID {
			return true
		}
	}
	return false
}

func jsort(jets []insolar.JetID) []string {
	sort.Slice(jets, func(i, j int) bool {
		return bytes.Compare(jets[i].Prefix(), jets[j].Prefix()) == -1
	})
	result := make([]string, 0, len(jets))
	for _, j := range jets {
		result = append(result, j.DebugString())
	}
	return result
}

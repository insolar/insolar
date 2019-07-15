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
	initialDepth      uint8 = 2
	defaultDepthLimit       = initialDepth + 2
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
	jet10_left  = jet.NewIDFromString("100")
	jet10_right = jet.NewIDFromString("101")
)

var cases = []splitCase{
	{
		name: "no_split",
		cfg: configuration.JetSplit{
			ThresholdRecordsCount:  5,
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
			{jet10: {3, 0, false}},
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
			{jet10: {5, 0, false}},
		},
	},
	{
		name: "no_split_with_overflow",
		cfg: configuration.JetSplit{
			ThresholdRecordsCount:  4,
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
			{jet10_left: {5, 0, false}},
		},
	},
}

func TestJetSplitter(t *testing.T) {
	ctx := inslogger.TestContext(t)

	// prepare initial pulses and jets
	var initalPulse insolar.PulseNumber = 60000
	previousPulse, endedPulse, newPulse := initalPulse, initalPulse+1, initalPulse+2
	initialJets := []insolar.JetID{jet0, jet10, jet11}

	checkCase := func(t *testing.T, sc splitCase) {
		// real components
		jetStore := jet.NewStore()
		err := jetStore.Update(ctx, endedPulse, true, initialJets...)
		require.NoError(t, err, "jet store updated with initial jets")
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

		// no filter for ID
		jetCalc.MineForPulseFunc = func(ctx context.Context, pn insolar.PulseNumber) []insolar.JetID {
			return jetStore.All(ctx, pn)
		}

		for i, jetsConfig := range sc.pulses {
			delta := insolar.PulseNumber(i)
			ended, newpulse := endedPulse+delta, newPulse+delta
			pulseCalc.BackwardsMock.Return(insolar.Pulse{PulseNumber: previousPulse + delta}, nil)

			pulseStartedWithJets := jetStore.All(ctx, ended)

			collectionAccessor.ForPulseFunc = func(_ context.Context, jetID insolar.JetID, pn insolar.PulseNumber) []record.Material {
				jConf, ok := jetsConfig[jetID]
				if !ok {
					return nil
				}
				return make([]record.Material, jConf.records)
			}

			gotJets, err := splitter.Do(ctx, ended, newpulse)
			require.NoError(t, err, "splitter.Do performed")

			for jetID, jConf := range jetsConfig {
				dropThreshold := splitter.getDropThreshold(ctx, jetID, ended)
				require.Equalf(t, jConf.dropThreshold, dropThreshold,
					"check drop.SplitThresholdExceeded for jet %v in +%v pulse", jetID.DebugString(), i)
			}

			var expectJets []insolar.JetID
			var splitedJets []string
			for _, jetID := range pulseStartedWithJets {
				jConf := jetsConfig[jetID]
				if jConf.hasSplit {
					left, right := jet.Siblings(jetID)
					expectJets = append(expectJets, left, right)
					splitedJets = append(splitedJets, jetID.DebugString())
					continue
				}
				expectJets = append(expectJets, jetID)
			}
			jetsInfo := "jets should split " + strings.Join(splitedJets, ", ")
			if len(splitedJets) == 0 {
				jetsInfo = "no jets spit"
			}

			expectMsg := fmt.Sprintf("jet %v should split on +%v pulse", jetsInfo, i)
			require.Equal(t, jsort(expectJets), jsort(gotJets), expectMsg)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { checkCase(t, c) })
	}
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

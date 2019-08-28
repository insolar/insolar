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
// +build slowtest

package integration_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
)

func Test_JetSplitEveryPulse(t *testing.T) {

	type splitCase struct {
		name               string
		cfg                configuration.JetSplit
		testPulsesQuantity int
	}

	var splitCases = []splitCase{
		{
			name: "splitEveryPulse",
			cfg: configuration.JetSplit{
				ThresholdRecordsCount:  0,
				ThresholdOverflowCount: 0,
				DepthLimit:             15,
			},
			testPulsesQuantity: 1,
		},
		{
			name: "splitEveryPulseLimitedDepth",
			cfg: configuration.JetSplit{
				ThresholdRecordsCount:  0,
				ThresholdOverflowCount: 0,
				DepthLimit:             3,
			},
			testPulsesQuantity: 5,
		},
	}

	t.Parallel()

	testCase := func(t *testing.T, sc splitCase) {

		var hotObjects = make(chan insolar.JetID)
		var replication = make(chan insolar.JetID)
		var hotObjectConfirm = make(chan insolar.JetID)

		var testPulsesQuantity = sc.testPulsesQuantity

		ctx := inslogger.TestContext(t)
		cfg := DefaultLightConfig()
		cfg.Ledger.JetSplit.DepthLimit = sc.cfg.DepthLimit
		cfg.Ledger.JetSplit.ThresholdOverflowCount = sc.cfg.ThresholdOverflowCount
		cfg.Ledger.JetSplit.ThresholdRecordsCount = sc.cfg.ThresholdRecordsCount

		s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {

			switch p := pl.(type) {
			case *payload.Replication:
				replication <- p.JetID

			case *payload.HotObjects:
				hotObjects <- p.JetID

			case *payload.GotHotConfirmation:
				hotObjectConfirm <- p.JetID
			}
			if meta.Receiver == NodeHeavy() {
				return DefaultHeavyResponse(pl)
			}
			return nil
		})
		require.NoError(t, err)
		defer s.Stop()

		calculateExpectedJets := func(jets []insolar.JetID, depthLimit uint8) []insolar.JetID {

			result := make([]insolar.JetID, 0, len(jets)*2)

			for _, jetID := range jets {
				if jetID.Depth() >= depthLimit {
					result = append(result, jetID)
				} else {
					jet1, jet2 := jet.Siblings(jetID)
					result = append(result, jet1, jet2)
				}
			}
			return result
		}

		// First pulse goes in storage then interrupts.
		s.SetPulse(ctx)

		{
			expectedJets := []insolar.JetID{insolar.ZeroJetID}

			for i := 0; i < testPulsesQuantity; i++ {

				s.SetPulse(ctx)

				previousPulseJets := expectedJets
				expectedJets = calculateExpectedJets(expectedJets, cfg.Ledger.JetSplit.DepthLimit)

				hotObjectsReceived := make(map[insolar.JetID]struct{})
				hotObjectsConfirmReceived := make(map[insolar.JetID]struct{})

				// collecting HO and HCO
				for range expectedJets {
					hotObjectsReceived[<-hotObjects] = struct{}{}
					hotObjectsConfirmReceived[<-hotObjectConfirm] = struct{}{}
				}

				for _, expectedJetId := range expectedJets {
					_, ok := hotObjectsReceived[expectedJetId]
					require.True(t, ok, "No expected jetId %s in hotObjectsReceived", expectedJetId.DebugString())

					_, ok = hotObjectsConfirmReceived[expectedJetId]
					require.True(t, ok, "No expected jetId %s in hotObjectsConfirmReceived", expectedJetId.DebugString())
				}

				// collecting Replication
				replicationObjectsReceived := make(map[insolar.JetID]struct{})
				for range previousPulseJets {
					replicationObjectsReceived[<-replication] = struct{}{}
				}

				for _, expectedJetId := range previousPulseJets {
					_, ok := replicationObjectsReceived[expectedJetId]
					require.True(t, ok, "No expected jetId %s in replicationObjectsReceived", expectedJetId.DebugString())
				}

				// check depthLimit works
				require.Equal(t, len(hotObjectsReceived), len(expectedJets))
			}
		}
	}

	for _, sc := range splitCases {
		t.Run(sc.name, func(t *testing.T) { testCase(t, sc) })
	}

}

func Test_JetSplitsWhenOverflows(t *testing.T) {
	t.Parallel()

	var hotObjects = make(chan insolar.JetID)
	var replication = make(chan insolar.JetID)
	var hotObjectConfirm = make(chan insolar.JetID)

	var pulsesQuantity = 10
	var recordsOnPulse = 3
	var splitOnPulse = 3
	var jetTree = jet.NewTree(true)

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	cfg.Ledger.JetSplit.DepthLimit = 5
	cfg.Ledger.JetSplit.ThresholdOverflowCount = 0
	cfg.Ledger.JetSplit.ThresholdRecordsCount = 2

	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
		switch p := pl.(type) {
		case *payload.Replication:
			replication <- p.JetID

		case *payload.HotObjects:
			hotObjects <- p.JetID

		case *payload.GotHotConfirmation:
			hotObjectConfirm <- p.JetID
		}
		if meta.Receiver == NodeHeavy() {
			return DefaultHeavyResponse(pl)
		}
		return nil
	})
	require.NoError(t, err)
	defer s.Stop()

	sendMessages := func(jetTree *jet.Tree) map[insolar.JetID]int {
		splittingJets := make(map[insolar.JetID]int)
		// Save code.
		for i := 0; i < recordsOnPulse; i++ {
			{
				p, _ := CallSetCode(ctx, s)
				RequireNotError(p)
				jetID, _ := jetTree.Find(p.(*payload.ID).ID)
				splittingJets[jetID]++
			}
		}
		return splittingJets
	}

	calculateExpectedJetsByTree := func(splittingJets map[insolar.JetID]int, jetTree *jet.Tree, depthLimit uint8, thresholdRecordsCount int) []insolar.JetID {

		for jetID, val := range splittingJets {
			if val >= thresholdRecordsCount && jetID.Depth() < depthLimit {
				_, _, _ = jetTree.Split(jetID)
			}
		}
		return jetTree.LeafIDs()
	}
	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	{
		expectedJets := []insolar.JetID{insolar.ZeroJetID}
		splittingJets := make(map[insolar.JetID]int)

		for i := 0; i < pulsesQuantity; i++ {

			s.SetPulse(ctx)

			// Saving previous for Replication check
			previousPulseJets := expectedJets

			// Starting to split test jet tree
			if i > splitOnPulse {
				expectedJets = calculateExpectedJetsByTree(splittingJets, jetTree, cfg.Ledger.JetSplit.DepthLimit, cfg.Ledger.JetSplit.ThresholdRecordsCount)
			}

			// Starting to send messages
			if i >= splitOnPulse {
				splittingJets = sendMessages(jetTree)
			}

			hotObjectsReceived := make(map[insolar.JetID]struct{})
			hotObjectsConfirmReceived := make(map[insolar.JetID]struct{})

			// collecting HO and HCO
			for range expectedJets {
				hotObjectsReceived[<-hotObjects] = struct{}{}
				hotObjectsConfirmReceived[<-hotObjectConfirm] = struct{}{}
			}

			for _, expectedJetId := range expectedJets {

				_, ok := hotObjectsReceived[expectedJetId]
				require.True(t, ok, "No expected jetId %s in hotObjectsReceived", expectedJetId.DebugString())

				_, ok = hotObjectsConfirmReceived[expectedJetId]
				require.True(t, ok, "No expected jetId %s in hotObjectsConfirmReceived", expectedJetId.DebugString())
			}

			// collecting Replication
			replicationObjectsReceived := make(map[insolar.JetID]struct{})
			for range previousPulseJets {
				replicationObjectsReceived[<-replication] = struct{}{}
			}

			for _, expectedJetId := range previousPulseJets {
				_, ok := replicationObjectsReceived[expectedJetId]
				require.True(t, ok, "No expected jetId %s in replicationObjectsReceived", expectedJetId.DebugString())
			}

		}
	}
}

func Test_LightStartsFromInitialState(t *testing.T) {
	t.Parallel()

	var hotObjects = make(chan insolar.JetID)
	var replication = make(chan insolar.JetID)
	var hotObjectConfirm = make(chan insolar.JetID)

	var initialSplits = 3
	var jetTree = jet.NewTree(true)

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	cfg.Ledger.JetSplit.DepthLimit = 5
	cfg.Ledger.JetSplit.ThresholdOverflowCount = 0
	cfg.Ledger.JetSplit.ThresholdRecordsCount = 2

	splitJetTree := func(jets []insolar.JetID, jetTree *jet.Tree, depthLimit uint8) []insolar.JetID {
		for _, jetID := range jets {
			if jetID.Depth() < depthLimit {
				_, _, _ = jetTree.Split(jetID)
			}
		}
		return jetTree.LeafIDs()
	}

	createDrops := func(jets []insolar.JetID) []drop.Drop {
		var drops []drop.Drop
		for _, jetID := range jets {
			drops = append(drops, drop.Drop{JetID: jetID, Pulse: insolar.FirstPulseNumber})
		}
		return drops
	}

	// Creating initial jet tree.
	initialJets := []insolar.JetID{insolar.ZeroJetID}
	for d := 0; d <= initialSplits; d++ {
		initialJets = splitJetTree(initialJets, jetTree, cfg.Ledger.JetSplit.DepthLimit)
	}

	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
		switch p := pl.(type) {
		case *payload.Replication:
			replication <- p.JetID

		case *payload.HotObjects:
			hotObjects <- p.JetID

		case *payload.GotHotConfirmation:
			hotObjectConfirm <- p.JetID
		}

		if meta.Receiver == NodeHeavy() {
			switch pl.(type) {
			case *payload.Replication, *payload.GotHotConfirmation:
				return nil
			case *payload.GetIndex:
				return []payload.Payload{&payload.Error{Code: payload.CodeNotFound}}
			case *payload.GetLightInitialState:
				return []payload.Payload{
					&payload.LightInitialState{
						NetworkStart: true,
						JetIDs:       initialJets,
						Pulse: pulse.PulseProto{
							PulseNumber: insolar.FirstPulseNumber,
						},
						Drops: createDrops(initialJets),
					},
				}
			}
		}
		return nil
	})
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	s.SetPulse(ctx)

	for i := 0; i < 10; i++ {
		p, _ := CallSetCode(ctx, s)
		RequireNotError(p)
	}

	hotObjectsReceived := make(map[insolar.JetID]struct{})
	hotObjectsConfirmReceived := make(map[insolar.JetID]struct{})
	replicationObjectsReceived := make(map[insolar.JetID]struct{})

	// collecting HO and HCO and Replication
	for range initialJets {
		hotObjectsReceived[<-hotObjects] = struct{}{}
		hotObjectsConfirmReceived[<-hotObjectConfirm] = struct{}{}

		replicationObjectsReceived[<-replication] = struct{}{}
	}

	for _, expectedJetId := range initialJets {
		_, ok := hotObjectsReceived[expectedJetId]
		require.True(t, ok, "No expected jetId %s in hotObjectsReceived", expectedJetId.DebugString())

		_, ok = hotObjectsConfirmReceived[expectedJetId]
		require.True(t, ok, "No expected jetId %s in hotObjectsConfirmReceived", expectedJetId.DebugString())

		_, ok = replicationObjectsReceived[expectedJetId]
		require.True(t, ok, "No expected jetId %s in replicationObjectsReceived", expectedJetId.DebugString())
	}

}

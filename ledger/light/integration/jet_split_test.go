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

package integration_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestJetSplit(t *testing.T) {
	t.Parallel()

	var hotObjects = make(chan insolar.JetID)
	var replication = make(chan insolar.JetID)
	var hotObjectConfirm = make(chan insolar.JetID)

	var testPulsesQuantity = 5

	// todo: InfoLevel
	ctx := inslogger.WithLoggerLevel(inslogger.TestContext(t), insolar.PanicLevel)
	cfg := DefaultLightConfig()
	cfg.Ledger.JetSplit.DepthLimit = 5
	cfg.Ledger.JetSplit.ThresholdOverflowCount = 0
	cfg.Ledger.JetSplit.ThresholdRecordsCount = 0

	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) {

		// fmt.Printf("type %T \n", pl)

		switch p := pl.(type) {
		case *payload.Replication:
			fmt.Printf("jetNum %s \n", p.JetID.DebugString())
			fmt.Printf("pulse index %s \n", p.Pulse.String())
			replication <- p.JetID

		case *payload.HotObjects:
			hotObjects <- p.JetID

		case *payload.GotHotConfirmation:
			hotObjectConfirm <- p.JetID
		}
	})

	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	fmt.Println("init pulse")

	{
		expectedJets := []insolar.JetID{insolar.ZeroJetID}

		for i := 0; i < testPulsesQuantity; i++ {

			fmt.Printf("\n\npulse: %d, %s\n", i, s.pulse.PulseNumber.String())
			s.Pulse(ctx)

			previousPulseJets := expectedJets
			expectedJets = calculateExpectedJets(expectedJets, cfg.Ledger.JetSplit.DepthLimit)

			hotObjectsReceived := make(map[insolar.JetID]struct{})
			hotObjectsConfirmReceived := make(map[insolar.JetID]struct{})

			// collecting HO and HCO
			for range expectedJets {
				hotObjectsReceived[<-hotObjects] = struct{}{}
				hotObjectsConfirmReceived[<-hotObjectConfirm] = struct{}{}
			}

			for kk := range hotObjectsReceived {
				fmt.Println("replication jet id: ", kk.DebugString())
			}

			for _, expectedJetId := range expectedJets {
				fmt.Println("exp jet id: ", expectedJetId.DebugString())

				_, ok := hotObjectsReceived[expectedJetId]
				require.True(t, ok, "No expected jetId in hotObjectsReceived")

				_, ok = hotObjectsConfirmReceived[expectedJetId]
				require.True(t, ok, "No expected jetId in hotObjectsConfirmReceived")
			}

			// collecting Replication
			replicationObjectsReceived := make(map[insolar.JetID]struct{})
			for range previousPulseJets {
				replicationObjectsReceived[<-replication] = struct{}{}
			}

			for _, expectedJetId := range previousPulseJets {
				fmt.Println("prev jet id: ", expectedJetId.DebugString())
				_, ok := replicationObjectsReceived[expectedJetId]
				require.True(t, ok, "No expected jetId in replicationObjectsReceived")
			}
		}
	}
}

func calculateExpectedJets(jetsMap []insolar.JetID, depthLimit uint8) []insolar.JetID {

	result := make([]insolar.JetID, 0, len(jetsMap)*2)

	for _, jetID := range jetsMap {
		if jetID.Depth() >= depthLimit {
			result = append(result, jetID)
		} else {
			jet1, jet2 := jet.Siblings(jetID)
			result = append(result, jet1, jet2)
		}

	}

	return result
}

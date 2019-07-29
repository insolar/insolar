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
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func Test_JetSplit(t *testing.T) {
	t.Parallel()

	var hotObjectsChan = make(chan insolar.JetID)
	var replicationChan = make(chan insolar.JetID)
	var hotObjectConfirmChan = make(chan insolar.JetID)

	var syncWg = new(sync.WaitGroup)
	var testPulsesQuantity = 10

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
			if p.JetID != insolar.ZeroJetID {
				replicationChan <- p.JetID
			}

		case *payload.HotObjects:
			fmt.Printf("jetNum %s \n", p.JetID.DebugString())
			fmt.Printf("len index %d \n", len(p.Indexes))

			if p.JetID != insolar.ZeroJetID {
				hotObjectsChan <- p.JetID
			}

		case *payload.GotHotConfirmation:
			if p.JetID != insolar.ZeroJetID {
				hotObjectConfirmChan <- p.JetID
			}
		}
	})

	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)
	fmt.Println("init pulses")

	syncWg.Add(1)

	go collectAndCheckData(
		ctx,
		t,
		testPulsesQuantity,
		cfg.Ledger.JetSplit.DepthLimit,
		s,
		syncWg,
		hotObjectsChan,
		replicationChan,
		hotObjectConfirmChan,
	)

	syncWg.Wait()
}

func collectAndCheckData(
	ctx context.Context,
	t *testing.T,
	testPulsesQuantity int,
	depthLimit uint8,
	s *Server,
	syncWg *sync.WaitGroup,
	jetsChannel chan insolar.JetID,
	replicationChannel chan insolar.JetID,
	hotObjectConfirmChan chan insolar.JetID,
) {
	defer syncWg.Done()

	var expectedJetsMap []insolar.JetID

	sendMessages := func() {
		// Creating root reason request.
		var reasonID insolar.ID
		{
			p, _ := callSetIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
			requireNotError(t, p)
			reasonID = p.(*payload.RequestInfo).RequestID
		}

		// Save and check code.
		{
			p, _ := callSetCode(ctx, t, s)
			requireNotError(t, p)
			// payloadId := p.(*payload.ID).ID
			// expectedIds = append(expectedIds, payloadId)
		}

		var expectedObjectID insolar.ID
		// Set, get request.
		{
			p, _ := callSetIncomingRequest(ctx, t, s, gen.ID(), reasonID, record.CTSaveAsChild)
			requireNotError(t, p)
			expectedObjectID = p.(*payload.RequestInfo).RequestID
			// expectedIds = append(expectedIds, expectedObjectID)
		}
		// Activate and check object.
		{
			p, _ := callActivateObject(ctx, t, s, expectedObjectID)
			requireNotError(t, p)

			// lifeline, material := requireGetObject(ctx, t, s, expectedObjectID)
			// expectedIds = append(expectedIds, *lifeline.LatestState)
			// require.Equal(t, &state, material.Virtual)
		}
		// Amend and check object.
		{
			p, _ := callSetIncomingRequest(ctx, t, s, expectedObjectID, reasonID, record.CTMethod)
			requireNotError(t, p)

			p, state := callAmendObject(ctx, t, s, expectedObjectID, p.(*payload.RequestInfo).RequestID)
			requireNotError(t, p)
			_, material := requireGetObject(ctx, t, s, expectedObjectID)
			require.Equal(t, &state, material.Virtual)

			// expectedLifeline = lifeline
			// expectedIds = append(expectedIds, *lifeline.LatestState)
		}
	}

	for i := 0; i < testPulsesQuantity; i++ {

		sendMessages()
		expectedJetsMap = calculateExpectedJets(expectedJetsMap, depthLimit)

		fmt.Println("\n\npulse: ", i)
		s.Pulse(ctx)

		hotObjects := make(map[insolar.JetID]struct{})
		replicationObjects := make(map[insolar.JetID]struct{})
		hotConfirmObjects := make(map[insolar.JetID]struct{})

		currentJetCount := len(expectedJetsMap)

		for k := 0; k < currentJetCount; k++ {
			hotObjects[<-jetsChannel] = struct{}{}
			replicationObjects[<-replicationChannel] = struct{}{}
			hotConfirmObjects[<-hotObjectConfirmChan] = struct{}{}
		}

		for kk := range hotObjects {
			fmt.Println("HO jet id: ", kk.DebugString())
		}

		for _, expectedJetId := range expectedJetsMap {
			fmt.Println("exp jet id: ", expectedJetId.DebugString())

			_, ok := hotObjects[expectedJetId]
			require.True(t, ok, "No jetId in HotObjects")

			_, ok = replicationObjects[expectedJetId]
			require.True(t, ok, "No jetId in Replication")

			_, ok = hotConfirmObjects[expectedJetId]
			require.True(t, ok, "No jetId in HotConfirmation")
		}

	}
}

func calculateExpectedJets(jetsMap []insolar.JetID, depthLimit uint8) []insolar.JetID {

	// creating first jet
	if len(jetsMap) == 0 {
		jetsMap = append(jetsMap, insolar.ZeroJetID)
	}

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

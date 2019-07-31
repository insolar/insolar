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
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
)

func Test_BootstrapCalls(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	t.Run("message before pulse received returns error", func(t *testing.T) {
		p, _ := CallSetCode(ctx, s)
		_, ok := p.(*payload.Error)
		assert.True(t, ok)
	})

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	t.Run("messages after two pulses return result", func(t *testing.T) {
		p, _ := CallSetCode(ctx, s)
		RequireNotError(p)
	})
}

func Test_LightReplication(t *testing.T) {
	t.Parallel()

	var secondPulseNumber = insolar.FirstPulseNumber + (PulseStep * 2)
	var expectedLifeline record.Lifeline
	var expectedObjectID insolar.ID

	var expectedIds []insolar.ID
	var receivedMessage = make(chan payload.Replication, 10)

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()

	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) {
		switch p := pl.(type) {
		case *payload.Replication:
			if p.Pulse == secondPulseNumber {
				go func() {
					receivedMessage <- *p
				}()
			}

		}
	})

	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	{
		// Creating root reason request.
		var reasonID insolar.ID
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		// Save and check code.
		{
			p, _ := CallSetCode(ctx, s)
			RequireNotError(p)
			payloadId := p.(*payload.ID).ID
			expectedIds = append(expectedIds, payloadId)
		}

		// Set, get request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), reasonID, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			expectedObjectID = rep.(*payload.RequestInfo).RequestID
			expectedIds = append(expectedIds, expectedObjectID)
		}
		// Activate and check object.
		{
			p, state := CallActivateObject(ctx, s, expectedObjectID)
			RequireNotError(p)

			lifeline, material := requireGetObject(ctx, t, s, expectedObjectID)
			expectedIds = append(expectedIds, *lifeline.LatestState)
			require.Equal(t, state, material.Virtual)
		}
		// Amend and check object.
		{
			msg, _ := MakeSetIncomingRequest(expectedObjectID, reasonID, false, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)

			p, state := CallAmendObject(ctx, s, expectedObjectID, rep.(*payload.RequestInfo).RequestID)
			RequireNotError(p)
			lifeline, material := requireGetObject(ctx, t, s, expectedObjectID)
			require.Equal(t, state, material.Virtual)

			expectedLifeline = lifeline
			expectedIds = append(expectedIds, *lifeline.LatestState)
		}
	}

	// Third pulse activate replication of second's pulse records
	s.SetPulse(ctx)

	{
		replicationPayload := <-receivedMessage

		var receivedLifeline record.Lifeline

		for _, recordIndex := range replicationPayload.Indexes {
			if recordIndex.ObjID == expectedObjectID {
				receivedLifeline = recordIndex.Lifeline
			}
		}

		replicatedIds := make(map[insolar.ID]struct{})

		require.Equal(t, 13, len(replicationPayload.Records))
		require.Equal(t, expectedLifeline, receivedLifeline)

		// testing payload
		cryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()

		for _, rec := range replicationPayload.Records {
			hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), rec.Virtual)
			id := insolar.NewID(secondPulseNumber, hash)
			replicatedIds[*id] = struct{}{}
		}

		for _, id := range expectedIds {
			_, ok := replicatedIds[id]
			require.True(t, ok, "No key in replicated data")
		}
	}

}

func Test_BasicOperations(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	runner := func(t *testing.T) {
		// Creating root reason request.
		var reasonID insolar.ID
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), true, true)
			rep := retryIfCancelled(func() payload.Payload {
				return SendMessage(ctx, s, &msg)
			})
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}
		// Save and check code.
		{
			var sent record.Virtual
			p := retryIfCancelled(func() payload.Payload {
				p, s := CallSetCode(ctx, s)
				sent = s
				return p
			})
			RequireNotError(p)

			p = CallGetCode(ctx, s, p.(*payload.ID).ID)
			RequireNotError(p)
			material := record.Material{}
			err := material.Unmarshal(p.(*payload.Code).Record)
			require.NoError(t, err)
			require.Equal(t, sent, material.Virtual)
		}
		var objectID insolar.ID
		// Set, get request.
		{
			msg, virtual := MakeSetIncomingRequest(gen.ID(), reasonID, true, true)
			p := retryIfCancelled(func() payload.Payload {
				return SendMessage(ctx, s, &msg)
			})
			RequireNotError(p)

			reqID := p.(*payload.RequestInfo).RequestID
			fmt.Println("asking for ", reqID.DebugString())
			p = CallGetRequest(ctx, s, p.(*payload.RequestInfo).RequestID)
			RequireNotError(p)
			require.Equal(t, virtual, p.(*payload.Request).Request)
			objectID = p.(*payload.Request).RequestID
		}
		// Activate and check object.
		{
			var state record.Virtual
			p := retryIfCancelled(func() payload.Payload {
				p, s := CallActivateObject(ctx, s, objectID)
				state = s
				return p
			})
			RequireNotError(p)
			_, material := requireGetObject(ctx, t, s, objectID)
			require.Equal(t, state, material.Virtual)
		}
		// Amend and check object.
		{
			msg, _ := MakeSetIncomingRequest(objectID, reasonID, false, true)
			p := retryIfCancelled(func() payload.Payload {
				return SendMessage(ctx, s, &msg)
			})
			RequireNotError(p)

			var state record.Virtual
			p = retryIfCancelled(func() payload.Payload {
				p, s := CallAmendObject(ctx, s, objectID, p.(*payload.RequestInfo).RequestID)
				state = s
				return p
			})
			RequireNotError(p)

			_, material := requireGetObject(ctx, t, s, objectID)
			require.Equal(t, state, material.Virtual)
		}
		// Deactivate and check object.
		{
			msg, _ := MakeSetIncomingRequest(objectID, reasonID, false, true)
			p := retryIfCancelled(func() payload.Payload {
				return SendMessage(ctx, s, &msg)
			})
			RequireNotError(p)

			retryIfCancelled(func() payload.Payload {
				p, _ := CallDeactivateObject(ctx, s, objectID, p.(*payload.RequestInfo).RequestID)
				return p
			})

			lifeline, _ := CallGetObject(ctx, s, objectID)
			_, ok := lifeline.(*payload.Error)
			assert.True(t, ok)
		}
	}

	t.Run("happy basic", runner)

	t.Run("happy concurrent", func(t *testing.T) {
		count := 100
		pulseAt := rand.Intn(count)
		var wg sync.WaitGroup
		wg.Add(count)
		for i := 0; i < count; i++ {
			if i == pulseAt {
				s.SetPulse(ctx)
			}
			i := i
			go func() {
				t.Run(fmt.Sprintf("iter %d", i), runner)
				wg.Done()
			}()
		}

		wg.Wait()
	})
}

func requireGetObject(ctx context.Context, t *testing.T, s *Server, objectID insolar.ID) (record.Lifeline, record.Material) {
	lifelinePL, statePL := CallGetObject(ctx, s, objectID)
	RequireNotError(lifelinePL)
	RequireNotError(statePL)

	lifeline := record.Lifeline{}
	err := lifeline.Unmarshal(lifelinePL.(*payload.Index).Index)
	require.NoError(t, err)

	state := record.Material{}
	err = state.Unmarshal(statePL.(*payload.State).Record)
	require.NoError(t, err)

	return lifeline, state
}

func retryIfCancelled(cb func() payload.Payload) payload.Payload {
	rep := cb()
	if err, ok := rep.(*payload.Error); ok {
		if err.Code == payload.CodeFlowCanceled {
			return retryIfCancelled(cb)
		}
	}

	return rep
}

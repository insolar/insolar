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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
)

func Test_LightReplication(t *testing.T) {
	t.Parallel()

	var secondPulseNumber = pulse.MinTimePulse + (PulseStep * 2)
	var expectedLifeline record.Lifeline
	var expectedObjectID insolar.ID

	var expectedIds []insolar.ID
	var receivedMessage = make(chan payload.Replication, 10)

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()

	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
		switch p := pl.(type) {
		case *payload.Replication:
			if p.Pulse == secondPulseNumber {
				go func() {
					receivedMessage <- *p
				}()
			}

		}
		if meta.Receiver == NodeHeavy() {
			return DefaultHeavyResponse(pl)
		}
		return nil
	})

	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	t.Run("messages after two pulses return result", func(t *testing.T) {
		p, _ := CallSetCode(ctx, s)
		RequireNotError(p)
		expectedIds = append(expectedIds, p.(*payload.ID).ID)
	})

	cryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()

	{
		var reasonID, lastFilament insolar.ID

		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
			expectedIds = append(expectedIds, reasonID)

			// Creating filament hash.
			{
				virtual := record.Wrap(&record.PendingFilament{
					RecordID:       reasonID,
					PreviousRecord: nil,
				})
				hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), virtual)
				id := *insolar.NewID(reasonID.Pulse(), hash)

				expectedIds = append(expectedIds, id)
			}
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
			msg, _ := MakeSetIncomingRequest(gen.ID(), reasonID, insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			expectedObjectID = rep.(*payload.RequestInfo).RequestID
			expectedIds = append(expectedIds, expectedObjectID)

			// Creating filament hash.
			{
				virtual := record.Wrap(&record.PendingFilament{
					RecordID:       expectedObjectID,
					PreviousRecord: nil,
				})
				hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), virtual)
				lastFilament = *insolar.NewID(reasonID.Pulse(), hash)

				expectedIds = append(expectedIds, lastFilament)
			}
		}

		// Activate object.
		{
			p, requestRec := CallActivateObject(ctx, s, expectedObjectID)
			RequireNotError(p)

			payloadId := p.(*payload.ResultInfo).ResultID
			expectedIds = append(expectedIds, payloadId)

			// Creating filament hash.
			{
				virtual := record.Wrap(&record.PendingFilament{
					RecordID:       payloadId,
					PreviousRecord: &lastFilament,
				})
				hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), virtual)
				lastFilament = *insolar.NewID(reasonID.Pulse(), hash)

				expectedIds = append(expectedIds, lastFilament)
			}

			// Create side effect hash.
			{
				hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), requestRec)
				id := *insolar.NewID(reasonID.Pulse(), hash)

				expectedIds = append(expectedIds, id)
			}

		}
		// Amend and check object.
		{
			msg, _ := MakeSetIncomingRequest(expectedObjectID, reasonID, insolar.ID{}, false, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)

			reqId := rep.(*payload.RequestInfo).RequestID
			expectedIds = append(expectedIds, reqId)

			// Create filament id.
			{
				virtual := record.Wrap(&record.PendingFilament{
					RecordID:       reqId,
					PreviousRecord: &lastFilament,
				})
				hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), virtual)
				lastFilament = *insolar.NewID(reasonID.Pulse(), hash)

				expectedIds = append(expectedIds, lastFilament)
			}

			p, amendRec := CallAmendObject(ctx, s, expectedObjectID, reqId)
			RequireNotError(p)

			reqId = p.(*payload.ResultInfo).ResultID
			expectedIds = append(expectedIds, reqId)

			// Create filament id.
			{
				virtual := record.Wrap(&record.PendingFilament{
					RecordID:       reqId,
					PreviousRecord: &lastFilament,
				})
				hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), virtual)
				id := *insolar.NewID(reasonID.Pulse(), hash)

				expectedIds = append(expectedIds, id)
			}

			// Create side effect hash.
			{
				hash := record.HashVirtual(cryptographyScheme.ReferenceHasher(), amendRec)
				id := *insolar.NewID(reasonID.Pulse(), hash)

				expectedIds = append(expectedIds, id)
			}

			lifeline, _ := requireGetObject(ctx, t, s, expectedObjectID)

			expectedLifeline = lifeline
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

		require.Equal(t, len(expectedIds), len(replicationPayload.Records))
		require.Equal(t, expectedLifeline, receivedLifeline)

		// testing payload
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

// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulse"

	"github.com/stretchr/testify/require"
)

var expectAbandoned int64

func Test_AbandonedNotification_WhenLightEmpty(t *testing.T) {
	// Configs.
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	cfg.LightChainLimit = 5

	// Responses from server.
	received := make(chan payload.AbandonedRequestsNotification)
	receivedConfirmations := make(chan payload.GotHotConfirmation)

	// Response from heavy.
	heavyResponse := func(pl payload.Payload) []payload.Payload {
		switch p := pl.(type) {
		case *payload.Replication, *payload.GotHotConfirmation:
			return nil
		case *payload.GetFilament: // Simulate heavy response when SetResult comes for filaments.
			virtual := record.Wrap(&record.PendingFilament{
				RecordID:       p.ObjectID,
				PreviousRecord: nil,
			})

			return []payload.Payload{&payload.FilamentSegment{
				ObjectID: p.ObjectID,
				Records: []record.CompositeFilamentRecord{
					{
						RecordID: p.ObjectID,
						MetaID:   p.StartFrom,
						Meta:     record.Material{Virtual: virtual},
						Record:   record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
					},
				},
			}}
		case *payload.GetLightInitialState:
			return []payload.Payload{&payload.LightInitialState{
				NetworkStart: true,
				JetIDs:       []insolar.JetID{insolar.ZeroJetID},
				Pulse: insolarPulse.PulseProto{
					PulseNumber: pulse.MinTimePulse,
				},
				Drops: []drop.Drop{
					{JetID: insolar.ZeroJetID, Pulse: pulse.MinTimePulse},
				},
				LightChainLimit: 5,
			}}
		case *payload.SearchIndex:
			return []payload.Payload{
				&payload.SearchIndexInfo{},
			}
		}

		panic(fmt.Sprintf("unexpected message to heavy %T", pl))
	}

	// Server init.
	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
		if notification, ok := pl.(*payload.AbandonedRequestsNotification); ok {
			atomic.AddInt64(&expectAbandoned, 1)
			received <- *notification
		}
		if confirmation, ok := pl.(*payload.GotHotConfirmation); ok {
			receivedConfirmations <- *confirmation
		}
		if meta.Receiver == NodeHeavy() {
			return heavyResponse(pl)
		}
		return nil
	})
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	t.Run("abandoned notification", func(t *testing.T) {
		// Set incoming request.
		msg, _ := MakeSetIncomingRequest(
			gen.ID(),
			gen.IDWithPulse(s.Pulse()),
			insolar.ID{},
			true,
			true,
			"",
		)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reqInfo := rep.(*payload.RequestInfo)
		objectID := reqInfo.ObjectID

		// Some pulses to reach the abandoned notification threshold.
		<-receivedConfirmations
		s.SetPulse(ctx)
		<-receivedConfirmations
		s.SetPulse(ctx)
		<-receivedConfirmations

		// Every pulse we must to send abandoned notifications (until it's processed).
		for i := 0; i < 100; i++ {
			s.SetPulse(ctx)
			<-receivedConfirmations

			notification := <-received
			require.Equal(t, objectID, notification.ObjectID)
		}

		requestID := reqInfo.RequestID

		// Set result -> close incoming request -> stop to send notifications.
		resMsg, _ := MakeSetResult(objectID, requestID)
		rep = SendMessage(ctx, s, &resMsg)
		RequireNotError(rep)

		// Checking for no notifications.
		for j := 0; j < 10; j++ {
			s.SetPulse(ctx)
			<-receivedConfirmations

			select {
			case _ = <-received:
				t.Error("unexpected abandoned notifications reply")
			default:
				// Do nothing. It's ok.
			}
		}
	})
}

func Test_AbandonedNotification_WhenLightInit(t *testing.T) {
	// Configs.
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	cfg.LightChainLimit = 5

	// Responses from server.
	received := make(chan payload.AbandonedRequestsNotification)
	receivedConfirmations := make(chan payload.GotHotConfirmation)

	// PulseNumber and ObjectID for mock heavy response
	pn := insolar.PulseNumber(pulse.MinTimePulse)
	objectID := gen.IDWithPulse(pn)

	// Response from heavy.
	heavyResponse := func(pl payload.Payload) []payload.Payload {
		switch pl.(type) {
		case *payload.Replication, *payload.GotHotConfirmation:
			return nil
		case *payload.GetLightInitialState:
			return []payload.Payload{&payload.LightInitialState{
				NetworkStart: true,
				JetIDs:       []insolar.JetID{insolar.ZeroJetID},
				Pulse: insolarPulse.PulseProto{
					PulseNumber: pulse.MinTimePulse,
				},
				Drops: []drop.Drop{
					{JetID: insolar.ZeroJetID, Pulse: pulse.MinTimePulse},
				},
				Indexes: []record.Index{
					{Lifeline: record.Lifeline{EarliestOpenRequest: &pn}, ObjID: objectID},
				},
				LightChainLimit: 5,
			}}
		}
		panic(fmt.Sprintf("unexpected message to heavy %T", pl))
	}

	// Server init.
	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
		if notification, ok := pl.(*payload.AbandonedRequestsNotification); ok {
			atomic.AddInt64(&expectAbandoned, 1)
			received <- *notification
		}
		if confirmation, ok := pl.(*payload.GotHotConfirmation); ok {
			receivedConfirmations <- *confirmation
		}
		if meta.Receiver == NodeHeavy() {
			return heavyResponse(pl)
		}
		return nil
	})
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	t.Run("abandoned notification from light start", func(t *testing.T) {
		// Some pulses to reach the abandoned notification threshold.
		<-receivedConfirmations
		s.SetPulse(ctx)
		<-receivedConfirmations
		s.SetPulse(ctx)
		<-receivedConfirmations

		// Every pulse we must to send abandoned notifications (until it's processed).
		for i := 0; i < 100; i++ {
			s.SetPulse(ctx)
			<-receivedConfirmations

			notification := <-received
			require.Equal(t, objectID, notification.ObjectID)
		}

	})
}

func Test_AbandonsMetricValue(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	cfg.LightChainLimit = 5
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	expectAbandonedValue := atomic.LoadInt64(&expectAbandoned)

	v := fetchMetricValue(
		s.metrics.Handler(),
		"insolar_requests_abandoned",
		float64(expectAbandonedValue),
		time.Second*5,
	)

	require.NoError(t, err, "fetch insolar_requests_abandoned metric value")
	// other tests could increment counter, so we expect at least expect value
	assert.GreaterOrEqualf(t, int64(v), expectAbandonedValue,
		"fetched insolar_requests_abandoned value equals or greater than calculated value")
}

func fetchMetricValue(
	h http.Handler,
	metricName string,
	expect float64,
	maxDuration time.Duration,
) float64 {
	tries := int64(5)
	var v float64
	for i := 0; i < int(tries); i++ {
		req, err := http.NewRequest("GET", "/metrics", nil)
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		v = insmetrics.SumMetricsValueByNamePrefix(rr.Body, metricName)
		if v > expect {
			break
		}
		time.Sleep(time.Duration(int64(maxDuration) / tries))
	}
	return v
}

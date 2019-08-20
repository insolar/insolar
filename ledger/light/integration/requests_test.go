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
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

// For better coverage of corner cases (pulse changing, messages from different pulses, etc)
// Server.SetPulse() should be put between logical ledger actions (set request, send message, set result, etc).
// But we can't cover all combinations here anyway. This should be done in unit tests.

func Test_IncomingRequest_Check(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	t.Run("registered is older than reason returns error", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()+1), insolar.ID{}, true, true)
		rep := SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeInvalidRequest)
	})

	t.Run("detached returns error", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		// Faking detached request.
		record.Unwrap(&msg.Request).(*record.IncomingRequest).ReturnMode = record.ReturnSaga
		rep := SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeInvalidRequest)
	})

	t.Run("registered API request appears in pendings", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reqInfo := rep.(*payload.RequestInfo)

		s.SetPulse(ctx)

		rep = CallGetPendings(ctx, s, reqInfo.RequestID)
		RequireNotError(rep)

		ids := rep.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, reqInfo.RequestID, ids.IDs[0])
	})

	t.Run("registered request appears in pendings", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		firstObjP := SendMessage(ctx, s, &msg)
		RequireNotError(firstObjP)
		reqInfo := firstObjP.(*payload.RequestInfo)

		s.SetPulse(ctx)

		msg, _ = MakeSetIncomingRequest(gen.ID(), reqInfo.RequestID, reqInfo.RequestID, true, false)
		secondObjP := SendMessage(ctx, s, &msg)
		RequireNotError(secondObjP)
		secondReqInfo := secondObjP.(*payload.RequestInfo)

		s.SetPulse(ctx)

		secondPendings := CallGetPendings(ctx, s, secondReqInfo.RequestID)
		RequireNotError(secondPendings)

		ids := secondPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, secondReqInfo.RequestID, ids.IDs[0])
	})

	t.Run("closed request does not appear in pendings", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reqInfo := rep.(*payload.RequestInfo)

		s.SetPulse(ctx)

		p, _ := CallActivateObject(ctx, s, reqInfo.RequestID)
		RequireNotError(p)

		s.SetPulse(ctx)

		p = CallGetPendings(ctx, s, reqInfo.RequestID)

		err := p.(*payload.Error)
		require.Equal(t, insolar.ErrNoPendingRequest.Error(), err.Text)
	})
}

func Test_IncomingRequest_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	t.Run("creation request duplicate found", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)

		// Set first request.
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		require.Nil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)

		s.SetPulse(ctx)

		// Try to set it again.
		rep = SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		require.NotNil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)

		// Check for result.
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(rep.(*payload.RequestInfo).Request)
		require.NoError(t, err)
		require.Equal(t, msg.Request, receivedDuplicate.Virtual)
	})

	t.Run("method request duplicate found", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		rootObject := rep.(*payload.RequestInfo).ObjectID
		reasonID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		msg, _ = MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep = SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		objectID := rep.(*payload.RequestInfo).ObjectID

		s.SetPulse(ctx)

		msg, _ = MakeSetIncomingRequest(objectID, reasonID, rootObject, false, false)

		// Set first request.
		rep = SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		require.Nil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)

		s.SetPulse(ctx)

		// Try to set it again.
		rep = SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		require.NotNil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)

		// Check for found duplicate.
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(rep.(*payload.RequestInfo).Request)
		require.NoError(t, err)
		require.Equal(t, msg.Request, receivedDuplicate.Virtual)
	})

	t.Run("method request duplicate with result found", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		rootObject := rep.(*payload.RequestInfo).ObjectID
		reasonID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		msg, _ = MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep = SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		objectID := rep.(*payload.RequestInfo).ObjectID

		s.SetPulse(ctx)

		requestMsg, _ := MakeSetIncomingRequest(objectID, reasonID, rootObject, false, false)

		// Set first request.
		rep = SendMessage(ctx, s, &requestMsg)
		RequireNotError(rep)
		require.Nil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)
		requestID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		// Set result.
		resMsg, resultVirtual := MakeSetResult(objectID, requestID)
		rep = SendMessage(ctx, s, &resMsg)
		RequireNotError(rep)

		s.SetPulse(ctx)

		// Try to set request again.
		rep = SendMessage(ctx, s, &requestMsg)
		RequireNotError(rep)
		requestInfo := rep.(*payload.RequestInfo)
		require.NotNil(t, requestInfo.Request)
		require.NotNil(t, requestInfo.Result)

		// Check for found duplicate.
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(requestInfo.Request)
		require.NoError(t, err)
		require.Equal(t, requestMsg.Request, receivedDuplicate.Virtual)

		// Check for result duplicate.
		receivedResult := record.Material{}
		err = receivedResult.Unmarshal(requestInfo.Result)
		require.NoError(t, err)
		require.Equal(t, resultVirtual, receivedResult.Virtual)
	})
}

func Test_OutgoingRequest_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	t.Run("method request duplicate found", func(t *testing.T) {
		args := make([]byte, 100)
		_, err := rand.Read(args)
		initReq := record.IncomingRequest{
			Object:    insolar.NewReference(gen.ID()),
			Arguments: args,
			CallType:  record.CTSaveAsChild,
			Reason:    *insolar.NewReference(*insolar.NewID(s.pulse.PulseNumber, []byte{1, 2, 3})),
			APINode:   gen.Reference(),
		}
		initReqMsg := &payload.SetIncomingRequest{
			Request: record.Wrap(&initReq),
		}

		// Set first request
		p := SendMessage(ctx, s, initReqMsg)
		RequireNotError(p)
		reqInfo := p.(*payload.RequestInfo)
		require.Nil(t, reqInfo.Request)
		require.Nil(t, reqInfo.Result)

		s.SetPulse(ctx)

		outgoingReq := record.OutgoingRequest{
			Object:   insolar.NewReference(reqInfo.RequestID),
			Reason:   *insolar.NewReference(reqInfo.RequestID),
			CallType: record.CTMethod,
			Caller:   *insolar.NewReference(reqInfo.RequestID),
		}
		outgoingReqMsg := &payload.SetOutgoingRequest{
			Request: record.Wrap(&outgoingReq),
		}

		// Set outgoing request
		outP := SendMessage(ctx, s, outgoingReqMsg)
		RequireNotError(outP)
		outReqInfo := outP.(*payload.RequestInfo)
		require.Nil(t, outReqInfo.Request)
		require.Nil(t, outReqInfo.Result)

		s.SetPulse(ctx)

		// Try to set an outgoing again
		outSecondP := SendMessage(ctx, s, outgoingReqMsg)
		RequireNotError(outSecondP)
		outReqSecondInfo := outSecondP.(*payload.RequestInfo)
		require.NotNil(t, outReqSecondInfo.Request)
		require.Nil(t, outReqSecondInfo.Result)

		// Check for the result
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(outReqSecondInfo.Request)
		require.NoError(t, err)
		require.Equal(t, &outgoingReq, record.Unwrap(&receivedDuplicate.Virtual))
	})
}

func Test_DetachedRequest_notification(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()

	received := make(chan payload.SagaCallAcceptNotification)
	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
		if notification, ok := pl.(*payload.SagaCallAcceptNotification); ok {
			received <- *notification
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

	t.Run("detached notification sent on detached reason close", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		objectID := rep.(*payload.RequestInfo).ObjectID

		s.SetPulse(ctx)

		msg, _ = MakeSetIncomingRequest(objectID, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true)
		rep = SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reasonID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		p, detachedRec := CallSetOutgoingRequest(ctx, s, objectID, reasonID, true)
		RequireNotError(p)
		detachedID := p.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		resMsg, _ := MakeSetResult(objectID, reasonID)
		rep = SendMessage(ctx, s, &resMsg)
		RequireNotError(rep)

		notification := <-received
		require.Equal(t, objectID, notification.ObjectID)
		require.Equal(t, detachedID, notification.DetachedRequestID)

		receivedRec := record.Virtual{}
		err := receivedRec.Unmarshal(notification.Request)
		require.NoError(t, err)
		require.Equal(t, detachedRec, receivedRec)
	})
}

func Test_Result_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)

	// Set request.
	rep := SendMessage(ctx, s, &msg)
	RequireNotError(rep)
	require.Nil(t, rep.(*payload.RequestInfo).Request)
	require.Nil(t, rep.(*payload.RequestInfo).Result)
	requestID := rep.(*payload.RequestInfo).RequestID
	objectID := requestID

	s.SetPulse(ctx)

	resMsg, resultVirtual := MakeSetResult(objectID, requestID)
	// Set result.
	rep = SendMessage(ctx, s, &resMsg)
	RequireNotError(rep)

	s.SetPulse(ctx)

	// Try to set it again.
	rep = SendMessage(ctx, s, &resMsg)
	RequireNotError(rep)

	resultInfo := rep.(*payload.ResultInfo)
	require.NotNil(t, resultInfo.Result)

	// Check duplicate.
	receivedResult := record.Material{}
	err = receivedResult.Unmarshal(resultInfo.Result)
	require.NoError(t, err)
	require.Equal(t, resultVirtual, receivedResult.Virtual)
}

func Test_IncomingRequest_ClosedReason(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	var reasonID insolar.ID

	t.Run("Incoming request can't be created w closed reason", func(t *testing.T) {

		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Closing request
		{
			objectID := reasonID

			resMsg, _ := MakeSetResult(objectID, reasonID)
			// Set result.
			rep := SendMessage(ctx, s, &resMsg)
			RequireNotError(rep)
		}

		s.SetPulse(ctx)

		// Creating incoming w closed reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), reasonID, reasonID, true, false)
			rep := SendMessage(ctx, s, &msg)
			RequireErrorCode(rep, payload.CodeReasonIsWrong)
		}
	})
}

func Test_OutgoingRequest_ClosedReason(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	var reasonID insolar.ID

	t.Run("Outgoing request can't be created w closed reason", func(t *testing.T) {
		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Closing request
		{
			resMsg, _ := MakeSetResult(reasonID, reasonID)
			// Set result.
			rep := SendMessage(ctx, s, &resMsg)
			RequireNotError(rep)
		}

		s.SetPulse(ctx)

		{
			pl, _ := MakeSetOutgoingRequest(reasonID, reasonID, false)
			rep := SendMessage(ctx, s, &pl)
			RequireErrorCode(rep, payload.CodeReasonNotFound)
		}
	})
}

func Test_Requests_OutgoingReason(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	var rootID, reasonID insolar.ID

	t.Run("Incoming/Outgoing request can't be created w outgoing reason", func(t *testing.T) {

		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			rootID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Creating outgoing
		{
			pl, _ := MakeSetOutgoingRequest(rootID, rootID, false)
			rep := SendMessage(ctx, s, &pl)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Creating wrong incoming
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), reasonID, rootID, true, false)
			rep := SendMessage(ctx, s, &msg)
			RequireErrorCode(rep, payload.CodeReasonIsWrong)
		}

		s.SetPulse(ctx)

		// Creating wrong outgoing
		{
			msg, _ := MakeSetOutgoingRequest(rootID, reasonID, false)
			rep := SendMessage(ctx, s, &msg)
			RequireErrorCode(rep, payload.CodeReasonNotFound)
		}
	})
}

func Test_OutgoingRequests_DifferentObjects(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	var rootID, rootID2 insolar.ID

	t.Run("Outgoing request can't be created w different from once object", func(t *testing.T) {

		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			rootID = rep.(*payload.RequestInfo).RequestID
		}
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			rootID2 = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Creating outgoing
		{
			pl, _ := MakeSetOutgoingRequest(rootID, rootID2, false)
			rep := SendMessage(ctx, s, &pl)
			RequireErrorCode(rep, payload.CodeReasonNotFound)
		}
	})
}

func Test_OutgoingDetached_InPendings(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	var rootID, secondReqId insolar.ID

	// Creating root reason request.
	{
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		rootID = rep.(*payload.RequestInfo).RequestID
	}

	t.Run("detached request not appears in pendings", func(t *testing.T) {

		s.SetPulse(ctx)

		// Creating outgoing
		pl, _ := MakeSetOutgoingRequest(rootID, rootID, true)
		rep := SendMessage(ctx, s, &pl)
		RequireNotError(rep)
		secondReqId = rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		firstPendings := CallGetPendings(ctx, s, rootID)
		RequireNotError(firstPendings)

		ids := firstPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.NotEqual(t, secondReqId, ids.IDs[0])
	})

	t.Run("detached request appears in pendings after closing root request", func(t *testing.T) {

		// Closing reason request
		{
			resMsg, _ := MakeSetResult(rootID, rootID)
			rep := SendMessage(ctx, s, &resMsg)
			RequireNotError(rep)
		}

		s.SetPulse(ctx)

		secondPendings := CallGetPendings(ctx, s, rootID)
		RequireNotError(secondPendings)

		ids := secondPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, secondReqId, ids.IDs[0])
	})
}

func Test_IncomingRequest_DifferentResults(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.SetPulse(ctx)

	var reasonID insolar.ID

	t.Run("Incoming request can't have several different results", func(t *testing.T) {
		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Closing request
		{
			resMsg, _ := MakeSetResult(reasonID, reasonID)
			rep := SendMessage(ctx, s, &resMsg)
			RequireNotError(rep)
		}

		s.SetPulse(ctx)

		{
			resMsg, _ := MakeSetResult(reasonID, reasonID)
			rep := SendMessage(ctx, s, &resMsg)
			RequireErrorCode(rep, payload.CodeRequestNotFound)
		}
	})
}

func Test_SetRequest_NoObjectReturnsError(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
		if meta.Receiver == NodeHeavy() {
			switch pl.(type) {
			case *payload.Replication, *payload.GotHotConfirmation:
				return nil
			case *payload.GetLightInitialState:
				return []payload.Payload{DefaultLightInitialState()}
			case *payload.GetIndex:
				return []payload.Payload{&payload.Error{Code: payload.CodeNotFound}}
			}
		}
		return nil
	})
	require.NoError(t, err)
	defer s.Stop()

	s.SetPulse(ctx)

	t.Run("incoming no object returns error", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.ID(), insolar.ID{}, false, true)
		rep := SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeNotFound)
	})

	t.Run("outgoing no object returns error", func(t *testing.T) {
		msg, _ := MakeSetOutgoingRequest(gen.ID(), gen.ID(), false)
		rep := SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeNotFound)
	})
}

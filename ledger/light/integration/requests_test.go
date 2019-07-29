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

func Test_IncomingRequests(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("pending was added", func(t *testing.T) {
		p, _ := callSetIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requireNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)
		p = fetchPendings(ctx, t, s, reqInfo.RequestID)
		requireNotError(t, p)

		ids := p.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, reqInfo.RequestID, ids.IDs[0])
	})

	t.Run("pending was added and closed", func(t *testing.T) {
		p, _ := callSetIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requireNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)

		p, _ = callActivateObject(ctx, t, s, reqInfo.RequestID)
		requireNotError(t, p)

		p = fetchPendings(ctx, t, s, reqInfo.RequestID)

		err := p.(*payload.Error)
		require.Equal(t, insolar.ErrNoPendingRequest.Error(), err.Text)
	})

	t.Run("reason on the another object", func(t *testing.T) {
		firstObjP, _ := callSetIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requireNotError(t, firstObjP)
		reqInfo := firstObjP.(*payload.RequestInfo)
		firstObjP, _ = callActivateObject(ctx, t, s, reqInfo.RequestID)
		requireNotError(t, firstObjP)

		secondObjP, _ := callSetIncomingRequest(ctx, t, s, gen.ID(), reqInfo.RequestID, record.CTSaveAsChild)
		requireNotError(t, secondObjP)
		secondReqInfo := secondObjP.(*payload.RequestInfo)
		secondPendings := fetchPendings(ctx, t, s, secondReqInfo.RequestID)
		requireNotError(t, secondPendings)

		ids := secondPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, secondReqInfo.RequestID, ids.IDs[0])
	})
}

func Test_OutgoingRequests(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()

	received := make(chan payload.SagaCallAcceptNotification)
	s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) {
		if notification, ok := pl.(*payload.SagaCallAcceptNotification); ok {
			received <- *notification
		}
	})
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("detached notification sent", func(t *testing.T) {
		p, _ := callSetIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requireNotError(t, p)
		objectID := p.(*payload.RequestInfo).ObjectID

		p, _ = callSetIncomingRequest(ctx, t, s, objectID, gen.ID(), record.CTMethod)
		requireNotError(t, p)
		reasonID := p.(*payload.RequestInfo).RequestID

		p, detachedRec := callSetOutgoingRequest(ctx, t, s, objectID, reasonID, true)
		requireNotError(t, p)
		detachedID := p.(*payload.RequestInfo).RequestID

		p, _ = callSetResult(ctx, t, s, objectID, reasonID)
		requireNotError(t, p)

		notification := <-received
		require.Equal(t, objectID, notification.ObjectID)
		require.Equal(t, detachedID, notification.DetachedRequestID)

		receivedRec := record.Virtual{}
		err := receivedRec.Unmarshal(notification.Request)
		require.NoError(t, err)
		require.Equal(t, detachedRec, receivedRec)
	})
}

func Test_DuplicatedRequests(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("try to register request twice. no result", func(t *testing.T) {
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
		p := sendMessage(ctx, t, s, initReqMsg)
		requireNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)
		require.Nil(t, reqInfo.Request)
		require.Nil(t, reqInfo.Result)

		// Try to set it again
		secondP := sendMessage(ctx, t, s, initReqMsg)
		requireNotError(t, secondP)
		reqInfo = secondP.(*payload.RequestInfo)
		require.NotNil(t, reqInfo.Request)
		require.Nil(t, reqInfo.Result)

		// Check for the result
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(reqInfo.Request)
		require.NoError(t, err)
		require.Equal(t, &initReq, record.Unwrap(&receivedDuplicate.Virtual))
	})

	t.Run("try to register outgoing request twice. no result", func(t *testing.T) {
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
		p := sendMessage(ctx, t, s, initReqMsg)
		requireNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)
		require.Nil(t, reqInfo.Request)
		require.Nil(t, reqInfo.Result)

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
		outP := sendMessage(ctx, t, s, outgoingReqMsg)
		requireNotError(t, outP)
		outReqInfo := p.(*payload.RequestInfo)
		require.Nil(t, outReqInfo.Request)
		require.Nil(t, outReqInfo.Result)

		// Try to set an outgoing again
		outSecondP := sendMessage(ctx, t, s, outgoingReqMsg)
		requireNotError(t, outSecondP)
		outReqSecondInfo := outSecondP.(*payload.RequestInfo)
		require.NotNil(t, outReqSecondInfo.Request)
		require.Nil(t, outReqSecondInfo.Result)

		// Check for the result
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(outReqSecondInfo.Request)
		require.NoError(t, err)
		require.Equal(t, &outgoingReq, record.Unwrap(&receivedDuplicate.Virtual))
	})

	t.Run("try to register request twice. when there is result", func(t *testing.T) {
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
		p := sendMessage(ctx, t, s, initReqMsg)
		requireNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)

		// Set result for the request
		p, _ = callActivateObject(ctx, t, s, reqInfo.RequestID)
		requireNotError(t, p)

		// Try to set it again
		secondP := sendMessage(ctx, t, s, initReqMsg)
		requireNotError(t, secondP)
		secondReqInfo := secondP.(*payload.RequestInfo)
		require.NotNil(t, secondReqInfo.Request)
		require.NotNil(t, secondReqInfo.Result)

		// Check for the request
		receivedDuplicateReq := record.Material{}
		err = receivedDuplicateReq.Unmarshal(secondReqInfo.Request)
		require.NoError(t, err)
		require.Equal(t, &initReq, record.Unwrap(&receivedDuplicateReq.Virtual))

		// Check for the result
		receivedDuplicateRes := record.Material{}
		err = receivedDuplicateRes.Unmarshal(secondReqInfo.Result)
		require.NoError(t, err)
		resultRecord := record.Unwrap(&receivedDuplicateRes.Virtual).(*record.Result)
		require.Equal(t, *insolar.NewReference(reqInfo.RequestID), resultRecord.Request)
		require.Equal(t, reqInfo.RequestID, resultRecord.Object)
	})

	t.Run("try to register result twice", func(t *testing.T) {
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
		p := sendMessage(ctx, t, s, initReqMsg)
		requireNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)

		// Set result for the request
		mem := make([]byte, 100)
		_, err = rand.Read(mem)
		require.NoError(t, err)
		rec := record.Wrap(&record.Activate{
			Request: *insolar.NewReference(reqInfo.RequestID),
			Memory:  mem,
		})
		buf, err := rec.Marshal()
		require.NoError(t, err)
		res := make([]byte, 100)
		_, err = rand.Read(res)
		require.NoError(t, err)
		resultRecord := record.Wrap(&record.Result{
			Request: *insolar.NewReference(reqInfo.RequestID),
			Object:  reqInfo.RequestID,
			Payload: res,
		})
		resBuf, err := resultRecord.Marshal()
		require.NoError(t, err)

		p = sendMessage(ctx, t, s, &payload.Activate{
			Record: buf,
			Result: resBuf,
		})
		requireNotError(t, p)

		// Try to set it again
		secondResP := sendMessage(ctx, t, s, &payload.Activate{
			Record: buf,
			Result: resBuf,
		})
		requireNotError(t, secondResP)
		secondReqInfo := secondResP.(*payload.ResultInfo)
		require.NotNil(t, secondReqInfo.Result)

		// Check for the result
		returnedResult := record.Material{}
		err = returnedResult.Unmarshal(secondReqInfo.Result)
		require.NoError(t, err)
		returnedRes := record.Unwrap(&returnedResult.Virtual).(*record.Result)
		require.Equal(t, *insolar.NewReference(reqInfo.RequestID), returnedRes.Request)
		require.Equal(t, reqInfo.RequestID, returnedRes.Object)
	})
}

func Test_CheckRequests(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("detached incoming fails with error", func(t *testing.T) {
		args := make([]byte, 100)
		_, err := rand.Read(args)
		require.NoError(t, err)
		initReq := record.IncomingRequest{
			Object:    insolar.NewReference(gen.ID()),
			Arguments: args,
			CallType:  record.CTSaveAsChild,
			Reason:    *insolar.NewReference(*insolar.NewID(s.pulse.PulseNumber, []byte{1, 2, 3})),
			APINode:   gen.Reference(),
			// Incoming can't be a detached request
			ReturnMode: record.ReturnSaga,
		}
		initReqMsg := &payload.SetIncomingRequest{
			Request: record.Wrap(&initReq),
		}

		// Set first request
		p := sendMessage(ctx, t, s, initReqMsg)
		errP, ok := p.(*payload.Error)
		require.Equal(t, true, ok)
		require.NotNil(t, errP)
	})
}

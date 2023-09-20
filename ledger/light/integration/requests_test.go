// +build slowtest

package integration_test

import (
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
)

// For better coverage of corner cases (pulse changing, messages from different pulses, etc)
// Server.SetPulse() should be put between logical ledger actions (set request, send message, set result, etc).
//
// Note, that we can't cover all combinations here anyway. This should be done in unit tests.

func Test_IncomingRequest_Check(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	t.Run("registered is older than reason returns error", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()+1), insolar.ID{}, true, true, "")
		rep := SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeRequestInvalid)
	})

	t.Run("registered API request appears in pendings", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true, "")
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reqInfo := rep.(*payload.RequestInfo)

		s.SetPulse(ctx)

		rep = CallGetPendings(ctx, s, reqInfo.RequestID, 1)
		RequireNotError(rep)

		ids := rep.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, reqInfo.RequestID, ids.IDs[0])
	})

	t.Run("registered request appears in pendings", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true, "")
		firstObjP := SendMessage(ctx, s, &msg)
		RequireNotError(firstObjP)
		reqInfo := firstObjP.(*payload.RequestInfo)

		s.SetPulse(ctx)

		msg, _ = MakeSetIncomingRequest(gen.ID(), reqInfo.RequestID, reqInfo.RequestID, true, false, "")
		secondObjP := SendMessage(ctx, s, &msg)
		RequireNotError(secondObjP)
		secondReqInfo := secondObjP.(*payload.RequestInfo)

		s.SetPulse(ctx)

		secondPendings := CallGetPendings(ctx, s, secondReqInfo.RequestID, 1)
		RequireNotError(secondPendings)

		ids := secondPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, secondReqInfo.RequestID, ids.IDs[0])
	})

	t.Run("closed request does not appear in pendings", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true, "")
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reqInfo := rep.(*payload.RequestInfo)

		s.SetPulse(ctx)

		p, _ := CallActivateObject(ctx, s, reqInfo.RequestID)
		RequireNotError(p)

		s.SetPulse(ctx)

		p = CallGetPendings(ctx, s, reqInfo.RequestID, 1)

		err := p.(*payload.Error)
		require.Equal(t, insolar.ErrNoPendingRequest.Error(), err.Text)
	})
}

func Test_IncomingRequest_Duplicate(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	t.Run("creation request duplicate found", func(t *testing.T) {
		msg, _ := MakeSetIncomingRequest(
			gen.ID(),
			gen.IDWithPulse(s.Pulse()),
			insolar.ID{},
			true,
			true,
			"reason",
		)
		// Create reason.
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reasonID := rep.(*payload.RequestInfo).RequestID

		msg, _ = MakeSetIncomingRequest(
			gen.ID(),
			reasonID,
			reasonID,
			true,
			false,
			"",
		)

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

		// Check for result.
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(rep.(*payload.RequestInfo).Request)
		require.NoError(t, err)
		require.Equal(t, msg.Request, receivedDuplicate.Virtual)
	})

	t.Run("outgoing request duplicate found", func(t *testing.T) {
		// Get reason object.
		reasonObject := CreateAndActivateObject(ctx, s, "object")

		// Make request on reason object that should be a reason for request on another object.
		msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "request")
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reasonID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		duplicate, _ := MakeSetOutgoingRequest(reasonObject, reasonID, false)

		// Set first request.
		rep = SendMessage(ctx, s, &duplicate)
		RequireNotError(rep)
		require.Nil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)

		s.SetPulse(ctx)

		// Try to set it again.
		rep = SendMessage(ctx, s, &duplicate)
		RequireNotError(rep)
		require.NotNil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)

		// Check for found duplicate.
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(rep.(*payload.RequestInfo).Request)
		require.NoError(t, err)
		require.Equal(t, duplicate.Request, receivedDuplicate.Virtual)
	})

	t.Run("incoming request duplicate with result found", func(t *testing.T) {
		// Get reason object.
		reasonObject := CreateAndActivateObject(ctx, s, "object")

		// Make reason request on reason object that should be a reason for request from another object.
		msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "reason")
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reasonID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		// Get second object.
		secondObject := CreateAndActivateObject(ctx, s, "object2")

		s.SetPulse(ctx)

		requestMsg, _ := MakeSetIncomingRequest(secondObject, reasonID, reasonObject, false, false, "")

		// Set first request on second object.
		rep = SendMessage(ctx, s, &requestMsg)
		RequireNotError(rep)
		require.Nil(t, rep.(*payload.RequestInfo).Request)
		require.Nil(t, rep.(*payload.RequestInfo).Result)
		requestID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		// Set result on second object..
		resMsg, resultVirtual := MakeSetResult(secondObject, requestID)
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
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	t.Run("method request duplicate found", func(t *testing.T) {
		// Get reason object.
		reasonObject := CreateAndActivateObject(ctx, s, "")

		s.SetPulse(ctx)

		msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "")
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reasonID := rep.(*payload.RequestInfo).RequestID

		outgoingReq := record.OutgoingRequest{
			Object:   insolar.NewReference(reasonObject),
			Reason:   *insolar.NewReference(reasonID),
			CallType: record.CTMethod,
			Caller:   *insolar.NewReference(reasonObject),
		}
		outgoingReqMsg := &payload.SetOutgoingRequest{
			Request: record.Wrap(&outgoingReq),
		}

		// Set outgoing request.
		outP := SendMessage(ctx, s, outgoingReqMsg)
		RequireNotError(outP)
		outReqInfo := outP.(*payload.RequestInfo)
		require.Nil(t, outReqInfo.Request)
		require.Nil(t, outReqInfo.Result)

		s.SetPulse(ctx)

		// Try to set an outgoing again.
		outSecondP := SendMessage(ctx, s, outgoingReqMsg)
		RequireNotError(outSecondP)
		outReqSecondInfo := outSecondP.(*payload.RequestInfo)
		require.NotNil(t, outReqSecondInfo.Request)
		require.Nil(t, outReqSecondInfo.Result)

		// Check for the result.
		receivedDuplicate := record.Material{}
		err = receivedDuplicate.Unmarshal(outReqSecondInfo.Request)
		require.NoError(t, err)
		require.Equal(t, &outgoingReq, record.Unwrap(&receivedDuplicate.Virtual))
	})
}

func Test_DetachedRequest_notification(t *testing.T) {
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

	t.Run("detached notification sent on detached reason close", func(t *testing.T) {
		// Get reason object.
		reasonObject := CreateAndActivateObject(ctx, s, "")

		s.SetPulse(ctx)

		msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "")
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)
		reasonID := rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		p, detachedRec := CallSetOutgoingRequest(ctx, s, reasonObject, reasonID, true)
		RequireNotError(p)
		detachedID := p.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		resMsg, _ := MakeSetResult(reasonObject, reasonID)
		rep = SendMessage(ctx, s, &resMsg)
		RequireNotError(rep)

		notification := <-received
		require.Equal(t, reasonObject, notification.ObjectID)
		require.Equal(t, detachedID, notification.DetachedRequestID)

		receivedRec := record.Virtual{}
		err := receivedRec.Unmarshal(notification.Request)
		require.NoError(t, err)
		require.Equal(t, detachedRec, receivedRec)
	})
}

func Test_Result_Duplicate(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	// Set request.
	msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true, "")
	rep := SendMessage(ctx, s, &msg)
	RequireNotError(rep)
	require.Nil(t, rep.(*payload.RequestInfo).Request)
	require.Nil(t, rep.(*payload.RequestInfo).Result)
	requestID := rep.(*payload.RequestInfo).RequestID
	objectID := requestID

	s.SetPulse(ctx)

	resMsg, _ := MakeSetResult(objectID, requestID)
	// Set result.
	rep = SendMessage(ctx, s, &resMsg)
	RequireNotError(rep)

	s.SetPulse(ctx)

	// Try to set it again.
	rep = SendMessage(ctx, s, &resMsg)
	RequireNotError(rep)
}

func Test_IncomingRequest_ClosedReason(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	var reasonID insolar.ID

	t.Run("Incoming request can't be created w closed reason", func(t *testing.T) {

		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true, "")
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Closing request.
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
			msg, _ := MakeSetIncomingRequest(gen.ID(), reasonID, reasonID, true, false, "")
			rep := SendMessage(ctx, s, &msg)
			RequireErrorCode(rep, payload.CodeReasonIsWrong)
		}
	})
}

func Test_IncomingRequest_ClosingWithOpenOutgoings(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	var reasonObject, reasonID insolar.ID
	t.Run("Incoming request can't be created w closed reason", func(t *testing.T) {

		// Get reason object.
		{
			reasonObject = CreateAndActivateObject(ctx, s, "")
		}

		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "")
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			// rootID = rep.(*payload.RequestInfo).RequestID
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Creating outgoing for request.
		{
			msg, _ := MakeSetOutgoingRequest(reasonObject, reasonID, false)
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
		}

		s.SetPulse(ctx)

		// Closing request.
		{
			resMsg, _ := MakeSetResult(reasonObject, reasonID)
			// Set result.
			rep := SendMessage(ctx, s, &resMsg)
			RequireErrorCode(rep, payload.CodeRequestNonClosedOutgoing)
		}
	})
}

func Test_IncomingRequest_ClosedReason_FromOtherObject(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	// If we need to test sequence of requests we can check `basic happy path` (program flow without pulse changing)
	// and `concurrent happy path` (program flow with random pulse changing):
	// 		t.Run("happy basic", func(t *testing.T) {...})
	// 		t.Run("happy concurrent", func(t *testing.T) {...})
	t.Run("detached incoming request from another object on closed reason", func(t *testing.T) {
		runner := func(t *testing.T) {
			var objectID insolar.ID  // Root reason object.
			var reasonID insolar.ID  // Root reason request.
			var anotherID insolar.ID // Another object.

			// Creating root reason object.
			{
				objectID = CreateAndActivateObject(ctx, s, "")
			}

			// Creating root reason request.
			{
				msg, _ := MakeSetIncomingRequest(
					objectID,
					gen.IDWithPulse(s.Pulse()),
					insolar.ID{},
					false,
					true,
					"reason",
				)
				rep := retryIfCancelled(func() payload.Payload {
					return SendMessage(ctx, s, &msg)
				})
				RequireNotError(rep)
				reasonID = rep.(*payload.RequestInfo).RequestID
			}

			// Creating detached outgoing request.
			{
				rep := retryIfCancelled(func() payload.Payload {
					p, _ := CallSetOutgoingRequest(ctx, s, objectID, reasonID, true)
					return p
				})
				RequireNotError(rep)
			}

			// Creating another object.
			{
				anotherID = CreateAndActivateObject(ctx, s, "")
			}

			// Closing reason request.
			{
				resMsg, _ := MakeSetResult(objectID, reasonID)
				rep := retryIfCancelled(func() payload.Payload {
					return SendMessage(ctx, s, &resMsg)
				})
				RequireNotError(rep)
			}

			// Creating request from another object with root object reason
			// in detached mode when reason closed already.
			{
				msg, _ := MakeSetIncomingRequestDetached(anotherID, reasonID, objectID, "request 2")
				rep := retryIfCancelled(func() payload.Payload {
					return SendMessage(ctx, s, &msg)
				})
				RequireNotError(rep)
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
					s.SetPulse(ctx) // Pulse changing.
				}
				i := i
				go func() {
					t.Run(fmt.Sprintf("iter %d", i), runner)
					wg.Done()
				}()
			}

			wg.Wait()
		})
	})
}

func Test_OutgoingRequest_ClosedReason(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	var reasonObject, reasonID insolar.ID

	t.Run("Outgoing request can't be created w closed reason", func(t *testing.T) {
		// Creating reason object.
		{
			reasonObject = CreateAndActivateObject(ctx, s, "")
		}

		// Creating reason request.
		{
			msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "")
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Closing request.
		{
			resMsg, _ := MakeSetResult(reasonObject, reasonID)
			// Set result.
			rep := SendMessage(ctx, s, &resMsg)
			RequireNotError(rep)
		}

		s.SetPulse(ctx)

		{
			pl, _ := MakeSetOutgoingRequest(reasonObject, reasonID, false)
			rep := SendMessage(ctx, s, &pl)
			RequireErrorCode(rep, payload.CodeReasonIsWrong)
		}
	})
}

func Test_Requests_OutgoingReason(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	var reasonObject, reasonID insolar.ID

	t.Run("Incoming/Outgoing request can't be created w outgoing reason", func(t *testing.T) {

		// Creating reason object.
		{
			reasonObject = CreateAndActivateObject(ctx, s, "")
		}

		// Creating reason request.
		{
			msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "")
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Creating outgoing.
		{
			pl, _ := MakeSetOutgoingRequest(reasonObject, reasonID, false)
			rep := SendMessage(ctx, s, &pl)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Creating wrong incoming.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), reasonID, reasonObject, true, false, "")
			rep := SendMessage(ctx, s, &msg)
			RequireErrorCode(rep, payload.CodeReasonIsWrong)
		}

		s.SetPulse(ctx)

		// Creating wrong outgoing.
		{
			msg, _ := MakeSetOutgoingRequest(reasonObject, reasonID, false)
			rep := SendMessage(ctx, s, &msg)
			RequireErrorCode(rep, payload.CodeReasonIsWrong)
		}
	})
}

func Test_OutgoingRequests_DifferentObjects(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	var rootID, rootID2 insolar.ID

	t.Run("Outgoing request can't be created w different from once object", func(t *testing.T) {

		// Creating root reason request.
		{
			rootID = CreateAndActivateObject(ctx, s, "")
		}

		{
			rootID2 = CreateAndActivateObject(ctx, s, "")
		}

		s.SetPulse(ctx)

		// Creating outgoing.
		{
			pl, _ := MakeSetOutgoingRequest(rootID, rootID2, false)
			rep := SendMessage(ctx, s, &pl)
			RequireErrorCode(rep, payload.CodeReasonIsWrong)
		}
	})
}

func Test_OutgoingDetached_InPendings(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	var outReqId, reasonObject, reasonID insolar.ID

	// Creating reason object.
	{
		reasonObject = CreateAndActivateObject(ctx, s, "")
	}

	s.SetPulse(ctx)
	{
		// Creating reason request.
		{
			msg, _ := MakeSetIncomingRequest(reasonObject, gen.IDWithPulse(s.Pulse()), insolar.ID{}, false, true, "")
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		// Creating outgoing.
		pl, _ := MakeSetOutgoingRequest(reasonObject, reasonID, true)
		rep := SendMessage(ctx, s, &pl)
		RequireNotError(rep)
		outReqId = rep.(*payload.RequestInfo).RequestID

		s.SetPulse(ctx)

		firstPendings := CallGetPendings(ctx, s, reasonObject, 1)
		RequireNotError(firstPendings)

		// detached request does not appears in pendings
		ids := firstPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.NotEqual(t, outReqId, ids.IDs[0])
	}

	// Detached request appears in pendings after closing root request.
	// Closing reason request.
	{
		{
			resMsg, _ := MakeSetResult(reasonObject, reasonID)
			rep := SendMessage(ctx, s, &resMsg)
			RequireNotError(rep)
		}

		s.SetPulse(ctx)

		secondPendings := CallGetPendings(ctx, s, reasonObject, 1)
		RequireNotError(secondPendings)

		ids := secondPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, outReqId, ids.IDs[0])
	}
}

func Test_IncomingRequest_DifferentResults(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop()

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	var reasonID insolar.ID

	t.Run("Incoming request can't have several different results", func(t *testing.T) {
		// Creating root reason request.
		{
			msg, _ := MakeSetIncomingRequest(gen.ID(), gen.IDWithPulse(s.Pulse()), insolar.ID{}, true, true, "")
			rep := SendMessage(ctx, s, &msg)
			RequireNotError(rep)
			reasonID = rep.(*payload.RequestInfo).RequestID
		}

		s.SetPulse(ctx)

		// Closing request.
		var originalResult record.Virtual
		{
			resMsg, virtual := MakeSetResult(reasonID, reasonID)
			rep := SendMessage(ctx, s, &resMsg)
			RequireNotError(rep)
			originalResult = virtual
		}

		s.SetPulse(ctx)

		{
			resMsg, _ := MakeSetResult(reasonID, reasonID)
			rep := SendMessage(ctx, s, &resMsg)
			res, ok := rep.(*payload.ErrorResultExists)
			require.True(t, ok, "returned ErrorResultExists")
			receivedResult := record.Material{}
			err := receivedResult.Unmarshal(res.Result)
			require.NoError(t, err)
			assert.Equal(t, originalResult, receivedResult.Virtual)
		}
	})
}

func Test_SetRequest_NoObjectReturnsError(t *testing.T) {
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
		msg, _ := MakeSetIncomingRequest(gen.ID(), gen.ID(), insolar.ID{}, false, true, "")
		rep := SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeNotFound)
	})

	t.Run("outgoing no object returns error", func(t *testing.T) {
		msg, _ := MakeSetOutgoingRequest(gen.ID(), gen.ID(), false)
		rep := SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeNotFound)
	})
}

func Test_SetRequest_LoopDetected(t *testing.T) {
	ctx := inslogger.TestContext(t)
	s, err := NewServer(ctx, DefaultLightConfig(), nil)
	require.NoError(t, err)

	s.SetPulse(ctx)

	t.Run("two requests with the same APIRequest trigger loop detection", func(t *testing.T) {
		objectID := CreateAndActivateObject(ctx, s, "object")

		s.SetPulse(ctx)

		msg, _ := MakeSetIncomingRequest(
			objectID,
			gen.IDWithPulse(s.Pulse()),
			insolar.ID{},
			false,
			true,
			"same request",
		)
		rep := SendMessage(ctx, s, &msg)
		RequireNotError(rep)

		msg, _ = MakeSetIncomingRequest(
			objectID,
			gen.IDWithPulse(s.Pulse()),
			insolar.ID{},
			false,
			true,
			"same request",
		)
		rep = SendMessage(ctx, s, &msg)
		RequireErrorCode(rep, payload.CodeLoopDetected)
	})
}

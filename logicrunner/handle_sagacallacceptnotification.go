package logicrunner

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
)

type HandleSagaCallAcceptNotification struct {
	dep  *Dependencies
	meta payload.Meta
}

func (h *HandleSagaCallAcceptNotification) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.SagaCallAcceptNotification{}
	err := msg.Unmarshal(h.meta.Payload)
	if err != nil {
		return err
	}

	outgoing := record.OutgoingRequest{}
	err = outgoing.Unmarshal(msg.Request)
	if err != nil {
		return err
	}

	// restore IncomingRequest by OutgoingRequest fields
	incoming := record.IncomingRequest{
		Caller:          outgoing.Caller,
		CallerPrototype: outgoing.CallerPrototype,
		Nonce:           outgoing.Nonce,

		Immutable: outgoing.Immutable,

		Object:    outgoing.Object,
		Prototype: outgoing.Prototype,
		Method:    outgoing.Method,
		Arguments: outgoing.Arguments,

		APIRequestID: outgoing.APIRequestID,
		Reason:       outgoing.Reason,

		// Saga calls are always asynchronous. We wait only for a confirmation
		// that the incoming request was registered by the second VE. This is
		// implemented in ContractRequester.CallMethod.
		ReturnMode: record.ReturnNoWait,
	}

	// Make a call to the second VE.
	callMsg := &message.CallMethod{IncomingRequest: incoming}
	cr := h.dep.lr.ContractRequester
	res, err := cr.CallMethod(ctx, callMsg)
	if err != nil {
		return err
	}

	// Register result of the outgoing method.
	outgoingReqRef := insolar.NewReference(msg.OutgoingReqID)
	reqResult := newRequestResult(res.(*reply.RegisterRequest).Request.Bytes(), outgoing.Caller)

	am := h.dep.lr.ArtifactManager
	return am.RegisterResult(ctx, *outgoingReqRef, reqResult)
}

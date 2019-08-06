package common

import (
	"context"

	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/logicrunner/transcript"
)

func BuildOutgoingSaveAsChildRequest(
	_ context.Context, current *transcript.Transcript, req rpctypes.UpSaveAsChildReq,
) *record.OutgoingRequest {

	current.Nonce++

	outgoing := record.OutgoingRequest{
		Caller:          req.Callee,
		CallerPrototype: req.CalleePrototype,
		Nonce:           current.Nonce,

		CallType:  record.CTSaveAsChild,
		Base:      &req.Parent,
		Prototype: &req.Prototype,
		Method:    req.ConstructorName,
		Arguments: req.ArgsSerialized,

		APIRequestID: current.Request.APIRequestID,
		Reason:       current.RequestRef,
	}

	return &outgoing
}

func BuildIncomingRequestFromOutgoing(outgoing *record.OutgoingRequest) *record.IncomingRequest {
	// Currently IncomingRequest and OutgoingRequest are almost exact copies of each other
	// thus the following code is a bit ugly. However this will change when we'll
	// figure out which fields are actually needed in OutgoingRequest and which are
	// not. Thus please keep the code the way it is for now, dont't introduce any
	// CommonRequestData structures or something like this.
	// This being said the implementation of Request interface differs for Incoming and
	// OutgoingRequest. See corresponding implementation of the interface methods.
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
	}

	if outgoing.ReturnMode == record.ReturnSaga {
		// We never wait for a result of saga call
		incoming.ReturnMode = record.ReturnNoWait
	} else {
		// If this is not a saga call just copy the ReturnMode
		incoming.ReturnMode = outgoing.ReturnMode
	}

	return &incoming
}

func BuildOutgoingRequest(
	_ context.Context, current *transcript.Transcript, req rpctypes.UpRouteReq,
) *record.OutgoingRequest {

	current.Nonce++

	outgoing := &record.OutgoingRequest{
		Caller:          req.Callee,
		CallerPrototype: req.CalleePrototype,
		Nonce:           current.Nonce,

		Immutable: req.Immutable,

		Object:    &req.Object,
		Prototype: &req.Prototype,
		Method:    req.Method,
		Arguments: req.Arguments,

		APIRequestID: current.Request.APIRequestID,
		Reason:       current.RequestRef,
	}

	if req.Saga {
		// OutgoingRequest with ReturnMode = ReturnSaga will be called by LME
		// when current object finishes the execution and validation.
		outgoing.ReturnMode = record.ReturnSaga
	} else if !req.Wait {
		outgoing.ReturnMode = record.ReturnNoWait
	}

	return outgoing
}

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

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

package outgoingsender

import (
	"context"
	"fmt"
	"strings"

	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/logicexecutor"

	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/go-actors/actor"
	"github.com/insolar/insolar/insolar/record"
)

var OutgoingRequestSenderDefaultQueueLimit = 1000

//go:generate minimock -i github.com/insolar/insolar/logicrunner/outgoingsender.OutgoingRequestSender -o ./ -s _mock.go -g

// OutgoingRequestSender is a type-safe wrapper for an actor implementation.
type OutgoingRequestSender interface {
	SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) (*insolar.Reference, insolar.Arguments, *record.IncomingRequest, error)
	SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest)
}

type outgoingRequestSender struct {
	as        actor.System
	senderPid actor.Pid
}

// Currently actor has only one state.
type outgoingSenderActorState struct {
	cr insolar.ContractRequester
	am artifacts.Client
}

type sendOutgoingResult struct {
	object   *insolar.Reference // only for CTSaveAsChild
	result   insolar.Arguments
	incoming *record.IncomingRequest // incoming request is used in a transcript
	err      error
}

type sendOutgoingRequestMessage struct {
	ctx              context.Context
	requestReference insolar.Reference       // registered request id
	outgoingRequest  *record.OutgoingRequest // outgoing request body
	resultChan       chan sendOutgoingResult // result that will be returned to the contract proxy
}

type sendAbandonedOutgoingRequestMessage struct {
	ctx              context.Context
	requestReference insolar.Reference       // registered request id
	outgoingRequest  *record.OutgoingRequest // outgoing request body
}

func NewOutgoingRequestSender(as actor.System, cr insolar.ContractRequester, am artifacts.Client) OutgoingRequestSender {
	pid := as.Spawn(func(system actor.System, pid actor.Pid) (actor.Actor, int) {
		state := newOutgoingSenderActorState(cr, am)
		queueLimit := OutgoingRequestSenderDefaultQueueLimit
		return state, queueLimit
	})

	return &outgoingRequestSender{
		as:        as,
		senderPid: pid,
	}
}

func (rs *outgoingRequestSender) SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) (*insolar.Reference, insolar.Arguments, *record.IncomingRequest, error) {
	resultChan := make(chan sendOutgoingResult, 1)
	msg := sendOutgoingRequestMessage{
		ctx:              ctx,
		requestReference: reqRef,
		outgoingRequest:  req,
		resultChan:       resultChan,
	}
	err := rs.as.Send(rs.senderPid, msg)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("SendOutgoingRequest failed: %v", err)
		return nil, insolar.Arguments{}, nil, err
	}

	res := <-resultChan
	return res.object, res.result, res.incoming, res.err
}

func (rs *outgoingRequestSender) SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) {
	msg := sendAbandonedOutgoingRequestMessage{
		ctx:              ctx,
		requestReference: reqRef,
		outgoingRequest:  req,
	}
	err := rs.as.Send(rs.senderPid, msg)
	if err != nil {
		// Actor's mailbox is most likely full. This is OK to lost an abandoned OutgoingRequest
		// in this case, LME will  re-send a corresponding notification anyway.
		inslogger.FromContext(ctx).Errorf("EnqueueAbandonedOutgoingRequest failed: %v", err)
	}
}

func newOutgoingSenderActorState(cr insolar.ContractRequester, am artifacts.Client) actor.Actor {
	return &outgoingSenderActorState{cr: cr, am: am}
}

func (a *outgoingSenderActorState) Receive(message actor.Message) (actor.Actor, error) {
	switch v := message.(type) {
	case sendOutgoingRequestMessage:
		// Currently it's possible to create an infinite number of goroutines here.
		// This number can be somehow limited in the future. The reason why a goroutine is
		// needed here is that an outgoing request can result in creating a new outgoing request
		// that can be directed to the same VE which would be waiting for a reply for a first
		// request, i.e. a deadlock situation.
		go func() {
			var res sendOutgoingResult
			res.object, res.result, res.incoming, res.err = a.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
			v.resultChan <- res
		}()
		return a, nil
	case sendAbandonedOutgoingRequestMessage:
		_, _, _, err := a.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
		// It's OK to just log an error,  LME will re-send a corresponding notification anyway.
		if err != nil {
			inslogger.FromContext(context.Background()).Errorf("abandonedOutgoingRequestActor: sendOutgoingRequest failed %v", err)
		}
		return a, nil
	default:
		inslogger.FromContext(context.Background()).Errorf("abandonedOutgoingRequestActor: unexpected message %v", v)
		return a, nil
	}
}

func (a *outgoingSenderActorState) sendOutgoingRequest(ctx context.Context, outgoingReqRef insolar.Reference, outgoing *record.OutgoingRequest) (*insolar.Reference, insolar.Arguments, *record.IncomingRequest, error) {
	var object *insolar.Reference
	var result insolar.Arguments

	incoming := common.BuildIncomingRequestFromOutgoing(outgoing)

	// Actually make a call.
	callMsg := &message.CallMethod{IncomingRequest: *incoming}
	res, err := a.cr.Call(ctx, callMsg)
	if err == nil {
		switch v := res.(type) {
		case *reply.CallMethod: // regular call
			object = v.Object // only for CTSaveAsChild
			result = v.Result
		case *reply.RegisterRequest: // no-wait call
			result = v.Request.Bytes()
		default:
			err = fmt.Errorf("sendOutgoingRequest: cr.Call returned unexpected type %T", v)
			return nil, result, nil, err
		}
	}

	// TODO: this is a part of horrible hack for making "index not found" error NOT system error. You MUST remove it in INS-3099
	if err != nil && !strings.Contains(err.Error(), "index not found") {
		return object, result, incoming, err
	}

	//  Register result of the outgoing method
	reqResult := logicexecutor.NewRequestResult(result, outgoing.Caller)
	registerResultErr := a.am.RegisterResult(ctx, outgoingReqRef, reqResult)

	// TODO: this is a part of horrible hack for making "index not found" error NOT system error. You MUST remove it in INS-3099
	if err != nil && strings.Contains(err.Error(), "index not found") {
		if registerResultErr != nil {
			inslogger.FromContext(ctx).Errorf("Failed to register result for request %s, error: %s", outgoingReqRef.String(), registerResultErr.Error())
		}
		return object, result, incoming, err
	}
	return object, result, incoming, registerResultErr
}

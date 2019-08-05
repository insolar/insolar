package logicrunner

import (
	"context"
	"strings"

	"github.com/insolar/insolar/logicrunner/artifacts"

	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/afiskon/go-actors/actor"
	"github.com/insolar/insolar/insolar/record"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.OutgoingRequestSender -o ./ -s _mock.go -g

// OutgoingRequestSender is a type-safe wrapper for an actor implementation.
type OutgoingRequestSender interface {
	SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) (insolar.Arguments, *record.IncomingRequest, error)
	SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest)
}

type outgoingRequestSender struct {
	senderPid actor.Pid
}

// Currently actor has only one state.
type outgoingSenderActorState struct {
	cr insolar.ContractRequester
	am artifacts.Client
}

type sendOutgoingResult struct {
	result   insolar.Arguments
	incoming *record.IncomingRequest // incoming request is used in a transcript
	err      error
}

type sendOutgoingRequestMessage struct {
	ctx              context.Context
	wait             bool
	requestReference insolar.Reference       // registered request id
	outgoingRequest  *record.OutgoingRequest // outgoing request body
	resultChan       chan sendOutgoingResult // result that will be returned to the contract proxy
}

type sendAbandonedOutgoingRequestMessage struct {
	ctx              context.Context
	requestReference insolar.Reference       // registered request id
	outgoingRequest  *record.OutgoingRequest // outgoing request body
}

func NewOutgoingRequestSender(cr insolar.ContractRequester, am artifacts.Client) OutgoingRequestSender {
	pid := GlobalActorSystem.Spawn(func(system actor.System, pid actor.Pid) (actor.Actor, int) {
		state := &outgoingSenderActorState{cr: cr, am: am}
		queueLimit := 1000
		return state, queueLimit
	})

	return &outgoingRequestSender{
		senderPid: pid,
	}
}

func (rs *outgoingRequestSender) SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) (insolar.Arguments, *record.IncomingRequest, error) {
	resultChan := make(chan sendOutgoingResult, 1)
	msg := sendOutgoingRequestMessage{
		ctx:              ctx,
		requestReference: reqRef,
		outgoingRequest:  req,
		resultChan:       resultChan,
	}
	err := GlobalActorSystem.Send(rs.senderPid, msg)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("SendOutgoingRequest failed: %v", err)
		return insolar.Arguments{}, nil, err
	}

	res := <-resultChan
	return res.result, res.incoming, res.err
}

func (rs *outgoingRequestSender) SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) {
	msg := sendAbandonedOutgoingRequestMessage{
		ctx:              ctx,
		requestReference: reqRef,
		outgoingRequest:  req,
	}
	err := GlobalActorSystem.Send(rs.senderPid, msg)
	if err != nil {
		// Actor's mailbox is most likely full. This is OK to lost an abandoned OutgoingRequest
		// in this case, LME will  re-send a corresponding notification anyway.
		inslogger.FromContext(ctx).Errorf("EnqueueAbandonedOutgoingRequest failed: %v", err)
	}
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
			res.result, res.incoming, res.err = a.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
			v.resultChan <- res
		}()
		return a, nil
	case sendAbandonedOutgoingRequestMessage:
		_, _, err := a.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
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

func (a *outgoingSenderActorState) sendOutgoingRequest(ctx context.Context, outgoingReqRef insolar.Reference, outgoing *record.OutgoingRequest) (insolar.Arguments, *record.IncomingRequest, error) {
	var result insolar.Arguments

	incoming := buildIncomingRequestFromOutgoing(outgoing)

	// Actually make a call.
	callMsg := &message.CallMethod{IncomingRequest: *incoming}
	res, err := a.cr.CallMethod(ctx, callMsg)
	if err == nil && (outgoing.ReturnMode == record.ReturnResult) {
		result = res.(*reply.CallMethod).Result
	}

	// TODO: this is a part of horrible hack for making "index not found" error NOT system error. You MUST remove it in INS-3099
	if err != nil && !strings.Contains(err.Error(), "index not found") {
		return result, incoming, err
	}

	//  Register result of the outgoing method
	reqResult := newRequestResult(result, outgoing.Caller)
	registerResultErr := a.am.RegisterResult(ctx, outgoingReqRef, reqResult)

	// TODO: this is a part of horrible hack for making "index not found" error NOT system error. You MUST remove it in INS-3099
	if err != nil && strings.Contains(err.Error(), "index not found") {
		if registerResultErr != nil {
			inslogger.FromContext(ctx).Errorf("Failed to register result for request %s, error: %s", outgoingReqRef.String(), registerResultErr.Error())
		}
		return result, incoming, err
	}
	return result, incoming, registerResultErr
}

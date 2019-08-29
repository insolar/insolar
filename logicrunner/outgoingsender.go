package logicrunner

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/insolar/go-actors/actor"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/requestresult"
)

var OutgoingRequestSenderDefaultQueueLimit = 1000
var OutgoingRequestSenderDefaultGoroutineLimit = int32(5000)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.OutgoingRequestSender -o ./ -s _mock.go -g

// OutgoingRequestSender is a type-safe wrapper for an actor implementation.
type OutgoingRequestSender interface {
	SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) (*insolar.Reference, insolar.Arguments, *record.IncomingRequest, error)
	SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest)
	Stop(ctx context.Context)
}

type outgoingRequestSender struct {
	as        actor.System
	senderPid actor.Pid
}

// Currently actor has only one state.
type outgoingSenderActorState struct {
	cr                            insolar.ContractRequester
	am                            artifacts.Client
	pa                            pulse.Accessor
	atomicRunningGoroutineCounter int32
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

func NewOutgoingRequestSender(as actor.System, cr insolar.ContractRequester, am artifacts.Client, pa pulse.Accessor) OutgoingRequestSender {
	pid := as.Spawn(func(system actor.System, pid actor.Pid) (actor.Actor, int) {
		state := newOutgoingSenderActorState(cr, am, pa)
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
		inslogger.FromContext(ctx).Errorf("SendAbandonedOutgoingRequest failed: %v", err)
	}
}

func (rs *outgoingRequestSender) Stop(_ context.Context) {
	rs.as.CloseAll()
}

func newOutgoingSenderActorState(cr insolar.ContractRequester, am artifacts.Client, pa pulse.Accessor) actor.Actor {
	return &outgoingSenderActorState{cr: cr, am: am, pa: pa}
}

func (a *outgoingSenderActorState) Receive(message actor.Message) (actor.Actor, error) {
	switch v := message.(type) {
	case sendOutgoingRequestMessage:
		if atomic.LoadInt32(&a.atomicRunningGoroutineCounter) >= OutgoingRequestSenderDefaultGoroutineLimit {
			var res sendOutgoingResult
			res.err = fmt.Errorf("OutgoingRequestActor: goroutine limit exceeded")
			v.resultChan <- res
			return a, nil
		}

		// The reason why a goroutine is needed here is that an outgoing request can result in
		// creating a new outgoing request that can be directed to the same VE which would be
		// waiting for a reply for a first request, i.e. a deadlock situation.
		// We limit the number of simultaneously running goroutines to prevent resource leakage.
		// It's OK to use atomics here because Receive is always executed by one goroutine. Thus
		// it's impossible to exceed the limit. It's possible that for a short period of time we'll
		// allow to create a little less goroutines that the limit says, but that's fine.
		atomic.AddInt32(&a.atomicRunningGoroutineCounter, 1)
		go func() {
			defer atomic.AddInt32(&a.atomicRunningGoroutineCounter, -1)

			var res sendOutgoingResult
			res.object, res.result, res.incoming, res.err = a.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
			v.resultChan <- res
		}()
		return a, nil
	case sendAbandonedOutgoingRequestMessage:
		_, _, _, err := a.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
		// It's OK to just log an error,  LME will re-send a corresponding notification anyway.
		if err != nil {
			inslogger.FromContext(context.Background()).Errorf("OutgoingRequestActor: sendOutgoingRequest failed %v", err)
		}
		return a, nil
	default:
		inslogger.FromContext(context.Background()).Errorf("OutgoingRequestActor: unexpected message %v", v)
		return a, nil
	}
}

func (a *outgoingSenderActorState) sendOutgoingRequest(ctx context.Context, outgoingReqRef insolar.Reference, outgoing *record.OutgoingRequest) (*insolar.Reference, insolar.Arguments, *record.IncomingRequest, error) {
	var object *insolar.Reference

	incoming := buildIncomingRequestFromOutgoing(outgoing)

	latestPulse, err := a.pa.Latest(ctx)
	if err != nil {
		err = errors.Wrapf(err, "sendOutgoingRequest: failed to get current pulse")
		return nil, nil, nil, err
	}
	// Actually make a call.
	callMsg := &payload.CallMethod{Request: incoming, PulseNumber: latestPulse.PulseNumber}
	res, _, err := a.cr.Call(ctx, callMsg)
	if err != nil {
		return nil, nil, nil, err
	}

	var result []byte

	switch v := res.(type) {
	case *reply.CallMethod: // regular call
		object = v.Object // only for CTSaveAsChild
		result = v.Result
	case *reply.RegisterRequest: // no-wait call
		result = v.Request.Bytes()
	default:
		err = fmt.Errorf("sendOutgoingRequest: cr.Call returned unexpected type %T", res)
		return nil, nil, nil, err
	}

	//  Register result of the outgoing method
	reqResult := requestresult.New(result, outgoing.Caller)
	err = a.am.RegisterResult(ctx, outgoingReqRef, reqResult)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can't register result")
	}

	return object, result, incoming, nil
}

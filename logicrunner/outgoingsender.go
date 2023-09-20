package logicrunner

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/go-actors/actor"
	aerr "github.com/insolar/go-actors/actor/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/requestresult"
)

var OutgoingRequestSenderDefaultQueueLimit = 1000
var OutgoingRequestSenderDefaultGoroutineLimit = int32(5000)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.OutgoingRequestSender -o ./ -s _mock.go -g

// OutgoingRequestSender is a type-safe wrapper for an actor implementation.
// Currently OutgoingRequestSender is implemented as a pair of actors. OutgoingSenderActor is
// responsible for sending regular outgoing requests and AbandonedSenderActor is responsible for
// sending abandoned requests. Generally we want to limit the number of outgoing requests, i.e. use
// some sort of queue. While this is easy for abandoned requests it's a bit tricky for regular requests
// (see comments below). Also generally speaking a synchronous abandoned request can create new outgoing requests
// which may cause a deadlock situation when a single actor is responsible for both types of messages. This is why two
// actors are used with two independent queues and their logic differs a little.
type OutgoingRequestSender interface {
	SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) (insolar.Arguments, *record.IncomingRequest, error)
	SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest)
	Stop(ctx context.Context)
}

type outgoingRequestSender struct {
	as                 actor.System
	outgoingSenderPid  actor.Pid
	abandonedSenderPid actor.Pid
}

type actorDeps struct {
	cr insolar.ContractRequester
	am artifacts.Client
	pa pulse.Accessor
}

type outgoingSenderActorState struct {
	deps                          actorDeps
	atomicRunningGoroutineCounter int32
}

type abandonedSenderActorState struct {
	deps actorDeps
}

type sendOutgoingResult struct {
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

type stopRequestSenderMessage struct {
	resultChan chan struct{}
}

func NewOutgoingRequestSender(as actor.System, cr insolar.ContractRequester, am artifacts.Client, pa pulse.Accessor) OutgoingRequestSender {
	outgoingSenderPid := as.Spawn(func(system actor.System, pid actor.Pid) (actor.Actor, int) {
		state := newOutgoingSenderActorState(cr, am, pa)
		queueLimit := OutgoingRequestSenderDefaultQueueLimit
		return state, queueLimit
	})

	abandonedSenderPid := as.Spawn(func(system actor.System, pid actor.Pid) (actor.Actor, int) {
		state := newAbandonedSenderActorState(cr, am, pa)
		queueLimit := OutgoingRequestSenderDefaultQueueLimit
		return state, queueLimit
	})

	return &outgoingRequestSender{
		as:                 as,
		outgoingSenderPid:  outgoingSenderPid,
		abandonedSenderPid: abandonedSenderPid,
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
	err := rs.as.Send(rs.outgoingSenderPid, msg)
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
	err := rs.as.Send(rs.abandonedSenderPid, msg)
	if err != nil {
		// Actor's mailbox is most likely full. This is OK to lost an abandoned OutgoingRequest
		// in this case, LME will  re-send a corresponding notification anyway.
		inslogger.FromContext(ctx).Errorf("SendAbandonedOutgoingRequest failed: %v", err)
	}
}

func (rs *outgoingRequestSender) Stop(ctx context.Context) {
	resultChanOutgoing := make(chan struct{}, 1)
	resultChanAbandoned := make(chan struct{}, 1)
	// We ignore both errors here because the only reason why SendPriority can fail
	// is that an actor doesn't exist or was already terminated. We don't expect either
	// situation here and there is no reasonable way to handle an error. If somehow
	// it happens Stop() will probably block forever and its OK (e.g. easy to debug using SIGABRT).
	rs.as.SendPriority(rs.outgoingSenderPid, stopRequestSenderMessage{ //nolint: errcheck
		resultChan: resultChanOutgoing,
	})
	rs.as.SendPriority(rs.abandonedSenderPid, stopRequestSenderMessage{ //nolint: errcheck
		resultChan: resultChanAbandoned,
	})

	// wait for a termination
	<-resultChanOutgoing
	<-resultChanAbandoned
}

func newOutgoingSenderActorState(cr insolar.ContractRequester, am artifacts.Client, pa pulse.Accessor) actor.Actor {
	return &outgoingSenderActorState{deps: actorDeps{cr: cr, am: am, pa: pa}}
}

func newAbandonedSenderActorState(cr insolar.ContractRequester, am artifacts.Client, pa pulse.Accessor) actor.Actor {
	return &abandonedSenderActorState{deps: actorDeps{cr: cr, am: am, pa: pa}}
}

func (a *outgoingSenderActorState) Receive(message actor.Message) (actor.Actor, error) {
	logger := inslogger.FromContext(context.Background()).WithField("actor", "outgoingSender")

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
		stats.Record(v.ctx, metrics.OutgoingSenderActorGoroutines.M(1))
		go func() {
			defer func() {
				atomic.AddInt32(&a.atomicRunningGoroutineCounter, -1)
				stats.Record(v.ctx, metrics.OutgoingSenderActorGoroutines.M(-1))
			}()

			var res sendOutgoingResult
			res.result, res.incoming, res.err = a.deps.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
			v.resultChan <- res
		}()
		return a, nil
	case stopRequestSenderMessage:
		v.resultChan <- struct{}{}
		return a, aerr.Terminate
	default:
		logger.Errorf("unexpected message %v", v)
		return a, nil
	}
}

func (a *abandonedSenderActorState) Receive(message actor.Message) (actor.Actor, error) {
	logger := inslogger.FromContext(context.Background()).WithField("actor", "abandonedSender")

	switch v := message.(type) {
	case sendAbandonedOutgoingRequestMessage:
		_, _, err := a.deps.sendOutgoingRequest(v.ctx, v.requestReference, v.outgoingRequest)
		// It's OK to just log an error,  LME will re-send a corresponding notification anyway.
		if err != nil {
			logger.Errorf("sendOutgoingRequest failed %v", err)
		}
		return a, nil
	case stopRequestSenderMessage:
		v.resultChan <- struct{}{}
		return a, aerr.Terminate
	default:
		logger.Errorf("unexpected message %v", v)
		return a, nil
	}
}

func (a *actorDeps) sendOutgoingRequest(ctx context.Context, outgoingReqRef insolar.Reference, outgoing *record.OutgoingRequest) (insolar.Arguments, *record.IncomingRequest, error) {
	incoming := buildIncomingRequestFromOutgoing(outgoing)

	latestPulse, err := a.pa.Latest(ctx)
	if err != nil {
		err = errors.Wrapf(err, "sendOutgoingRequest: failed to get current pulse")
		return nil, nil, err
	}

	inslogger.FromContext(ctx).Debug("sending incoming for outgoing request")

	// Actually make a call.
	callMsg := &payload.CallMethod{Request: incoming, PulseNumber: latestPulse.PulseNumber}
	res, _, err := a.cr.SendRequest(ctx, callMsg)
	if err != nil {
		return nil, nil, err
	}

	inslogger.FromContext(ctx).Debug("sent incoming for outgoing request")

	var result []byte

	switch v := res.(type) {
	case *reply.CallMethod: // regular call
		result = v.Result
	case *reply.RegisterRequest: // no-wait call
		result = v.Request.Bytes()
	default:
		err = fmt.Errorf("sendOutgoingRequest: cr.Call returned unexpected type %T", res)
		inslogger.FromContext(ctx).Error(err)
		return nil, nil, err
	}

	inslogger.FromContext(ctx).Debug("registering outgoing request result")

	//  Register result of the outgoing method
	reqResult := requestresult.New(result, outgoing.Caller)
	err = a.am.RegisterResult(ctx, outgoingReqRef, reqResult)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't register result")
	}

	inslogger.FromContext(ctx).Debug("registered outgoing request result")

	return result, incoming, nil
}

package logicrunner

import (
	"context"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/afiskon/go-actors/actor"
	"github.com/insolar/insolar/insolar/record"
)

// AbandonedOutgoingRequestSender is a type-safe wrapper for an actor implementation.
type AbandonedOutgoingRequestSender interface {
	EnqueueAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest)
}

type abandonedOutgoingRequestSender struct {
	senderPid actor.Pid
}

// Currently actor has only one state.
type abandonedOutgoingRequestActorState struct{}

// When actor receives this message it builds and sends a corresponding request.
type sendAbandonedOutgoingRequestMessage struct {
	requestReference insolar.Reference       // registered request id
	outgoingRequest  *record.OutgoingRequest // outgoing request body
}

func NewAbandonedOutgoingRequestSender() AbandonedOutgoingRequestSender {
	pid := GlobalActorSystem.Spawn(func(system actor.System, pid actor.Pid) (state actor.Actor, limit int) {
		return &abandonedOutgoingRequestActorState{}, 1000
	})

	return &abandonedOutgoingRequestSender{
		senderPid: pid,
	}
}

func (rs *abandonedOutgoingRequestSender) EnqueueAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) {
	msg := sendAbandonedOutgoingRequestMessage{
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
func (a *abandonedOutgoingRequestActorState) Receive(message actor.Message) (actor.Actor, error) {
	switch v := message.(type) {
	case sendAbandonedOutgoingRequestMessage:
		// TODO build and send an outgoing request from v.req
		return a, nil
	default:
		inslogger.FromContext(context.Background()).Errorf("abandonedOutgoingRequestActor: unexpected message %v", v)
		return a, nil
	}
}

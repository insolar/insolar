package logicrunner

import (
	"context"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/afiskon/go-actors/actor"
	"github.com/insolar/insolar/insolar/record"
)

// AALEKSEEV TODO use a pull here

//go:generate minimock -i github.com/insolar/insolar/logicrunner.OutgoingRequestSender -o ./ -s _mock.go -g

// OutgoingRequestSender is a type-safe wrapper for an actor implementation.
type OutgoingRequestSender interface {
	SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) error
	SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest)
}

type outgoingRequestSender struct {
	senderPid actor.Pid
}

// Currently actor has only one state.
type outgoingSenderActorState struct{}

type sendOutgoingRequestMessage struct {
	requestReference insolar.Reference       // registered request id
	outgoingRequest  *record.OutgoingRequest // outgoing request body
	resultChan       chan error              // result that will be returned to the contract proxy
}

type sendAbandonedOutgoingRequestMessage struct {
	requestReference insolar.Reference       // registered request id
	outgoingRequest  *record.OutgoingRequest // outgoing request body
}

func NewOutgoingRequestSender() OutgoingRequestSender {
	pid := GlobalActorSystem.Spawn(func(system actor.System, pid actor.Pid) (state actor.Actor, limit int) {
		return &outgoingSenderActorState{}, 1000
	})

	return &outgoingRequestSender{
		senderPid: pid,
	}
}

func (rs *outgoingRequestSender) SendOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) error {
	resultChan := make(chan error, 1)
	msg := sendOutgoingRequestMessage{
		requestReference: reqRef,
		outgoingRequest:  req,
		resultChan:       resultChan,
	}
	err := GlobalActorSystem.Send(rs.senderPid, msg)
	if err != nil {
		// Actor's mailbox is most likely full. This is OK to lost an abandoned OutgoingRequest
		// in this case, LME will  re-send a corresponding notification anyway.
		inslogger.FromContext(ctx).Errorf("SendOutgoingRequest failed: %v", err)
		return err
	}

	err = <-resultChan
	return err
}

func (rs *outgoingRequestSender) SendAbandonedOutgoingRequest(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) {
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

func (a *outgoingSenderActorState) Receive(message actor.Message) (actor.Actor, error) {
	switch v := message.(type) {
	case sendAbandonedOutgoingRequestMessage:
		err := a.sendOutgoingRequest(v.requestReference, v.outgoingRequest)
		// It's OK to just log an error,  LME will  re-send a corresponding notification anyway.
		if err != nil {
			inslogger.FromContext(context.Background()).Errorf("abandonedOutgoingRequestActor: sendOutgoingRequest failed %v", err)
		}
		return a, nil
	case sendOutgoingRequestMessage:
		err := a.sendOutgoingRequest(v.requestReference, v.outgoingRequest)
		v.resultChan <- err
		return a, nil
	default:
		inslogger.FromContext(context.Background()).Errorf("abandonedOutgoingRequestActor: unexpected message %v", v)
		return a, nil
	}
}

func (a *outgoingSenderActorState) sendOutgoingRequest(reqRef insolar.Reference, req *record.OutgoingRequest) error {
	// AALEKSEEV TODO move the logic here
	return nil
}

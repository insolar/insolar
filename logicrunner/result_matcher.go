package logicrunner

import (
	"sync"

	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/messagebus"
	"github.com/pkg/errors"
)

type resultsMatcher struct {
	lr                *LogicRunner
	lock              *sync.RWMutex
	executionNodes    map[insolar.Reference]insolar.Reference
	unwantedResponses map[insolar.Reference]message.ReturnResults
}

func newResultsMatcher(lr *LogicRunner) resultsMatcher {
	return resultsMatcher{
		lr:                lr,
		lock:              &sync.RWMutex{},
		executionNodes:    make(map[insolar.Reference]insolar.Reference),
		unwantedResponses: make(map[insolar.Reference]message.ReturnResults),
	}
}

func (rm *resultsMatcher) send(ctx context.Context, msg insolar.Message, receiver *insolar.Reference) {
	sender := messagebus.BuildSender(
		rm.lr.MessageBus.Send,
		messagebus.RetryIncorrectPulse(rm.lr.PulseAccessor),
		messagebus.RetryFlowCancelled(rm.lr.PulseAccessor),
	)
	_, err := sender(ctx, msg, &insolar.MessageSendOptions{
		Receiver: receiver,
	})
	if err != nil {
		inslogger.FromContext(ctx).Warn(errors.Wrap(err, "[ resultsMatcher::send ] Couldn't resend response"))
	}
}

func (rm *resultsMatcher) AddStillExecution(ctx context.Context, msg *message.StillExecuting) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	for _, reqRef := range msg.RequestRefs {
		if response, ok := rm.unwantedResponses[reqRef]; ok {
			inslogger.FromContext(ctx).Debug("[ resultsMatcher::AddStillExecution ] resend unwanted response ", reqRef)
			go rm.send(ctx, &response, &msg.Executor)
			delete(rm.executionNodes, reqRef)
			continue
		}
		rm.executionNodes[reqRef] = msg.Executor
	}
}

var flowCancelledError = &reply.Error{
	ErrType: reply.FlowCancelled,
}

func (rm *resultsMatcher) AddUnwantedResponse(ctx context.Context, msg insolar.Message) insolar.Reply {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	pulse, err := rm.lr.PulseAccessor.Latest(ctx)
	if err != nil {
		return flowCancelledError
	}
	response := msg.(*message.ReturnResults)
	node, err := rm.lr.JetCoordinator.VirtualExecutorForObject(ctx, *response.Target.Record(), pulse.PulseNumber)
	if err != nil {
		return flowCancelledError
	}
	if *node != rm.lr.JetCoordinator.Me() {
		return flowCancelledError
	}
	if node, ok := rm.executionNodes[response.Reason]; ok {
		inslogger.FromContext(ctx).Debug("[ resultsMatcher::AddUnwantedResponse ] resend unwanted response ", response.Reason)
		go rm.send(ctx, response, &node)
		delete(rm.unwantedResponses, response.Reason)
		return &reply.OK{}
	}

	newPulse, err := rm.lr.PulseAccessor.Latest(ctx)
	if err != nil {
		return flowCancelledError
	}
	if newPulse.PulseNumber == pulse.PulseNumber {
		rm.unwantedResponses[response.Reason] = *response
		return &reply.OK{}
	}
	return flowCancelledError
}

func (rm *resultsMatcher) Clear() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.executionNodes = make(map[insolar.Reference]insolar.Reference)
	rm.unwantedResponses = make(map[insolar.Reference]message.ReturnResults)
}

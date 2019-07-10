package logicrunner

import (
	"sync"

	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
)

type resultsMatcher struct {
	mb                insolar.MessageBus
	lock              *sync.RWMutex
	executionNodes    map[*insolar.Reference]*insolar.Reference
	unwantedResponses map[*insolar.Reference]*message.ReturnResults
}

func NewResultsMatcher(mb insolar.MessageBus) resultsMatcher {
	return resultsMatcher{
		mb:                mb,
		lock:              &sync.RWMutex{},
		executionNodes:    make(map[*insolar.Reference]*insolar.Reference),
		unwantedResponses: make(map[*insolar.Reference]*message.ReturnResults),
	}
}

func (rm *resultsMatcher) AddStillExecution(ctx context.Context, obj *insolar.Reference, node *insolar.Reference) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if response, ok := rm.unwantedResponses[obj]; ok {
		response.Target = *node
		// todo maybe call rm.mb.Send in goroutine
		// todo check errors? retry?
		rm.mb.Send(ctx, response, nil)
		return
	}
	rm.executionNodes[obj] = node
}

func (rm *resultsMatcher) AddUnwantedResponse(ctx context.Context, msg insolar.Message) {
	response, ok := msg.(*message.ReturnResults)
	if !ok {
		return //nil, errors.New("ReceiveResult() accepts only message.ReturnResults")
	}

	rm.lock.Lock()
	defer rm.lock.Unlock()
	if node, ok := rm.executionNodes[&response.RequestRef]; ok {
		response.Target = *node
		// todo maybe call rm.mb.Send in goroutine
		// todo check errors? retry?
		rm.mb.Send(ctx, response, nil)
		return
	}
	rm.unwantedResponses[&response.RequestRef] = response
}

func (rm *resultsMatcher) Clear() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.executionNodes = make(map[*insolar.Reference]*insolar.Reference)
	rm.unwantedResponses = make(map[*insolar.Reference]*message.ReturnResults)
}

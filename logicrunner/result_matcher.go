package logicrunner

import (
	"sync"

	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type resultsMatcher struct {
	// mb                insolar.MessageBus
	lr                *LogicRunner
	lock              *sync.RWMutex
	executionNodes    map[insolar.Reference]insolar.Reference
	unwantedResponses map[insolar.Reference]message.ReturnResults
}

func NewResultsMatcher(lr *LogicRunner) resultsMatcher {
	return resultsMatcher{
		lr: lr,
		// mb:                mb,
		lock:              &sync.RWMutex{},
		executionNodes:    make(map[insolar.Reference]insolar.Reference),
		unwantedResponses: make(map[insolar.Reference]message.ReturnResults),
	}
}

func (rm *resultsMatcher) AddStillExecution(ctx context.Context, msg *message.StillExecuting) {
	inslogger.FromContext(ctx).Warn("IP1: Receive StillExecution", msg.RequestRefs, "from", msg.Executor)
	rm.lock.Lock()
	defer rm.lock.Unlock()
	for _, reqRef := range msg.RequestRefs {
		if response, ok := rm.unwantedResponses[reqRef]; ok {
			// response.Target = *node
			// todo maybe call rm.mb.Send in goroutine
			// todo check errors? retry?
			inslogger.FromContext(ctx).Warn("IP1: Send StillExecution", reqRef)
			// go
			if rm.lr.MessageBus == nil {
				panic("NO MESSAGE BUS")
			}

			rm.lr.MessageBus.Send(ctx, &response, &insolar.MessageSendOptions{
				Receiver: &msg.Executor,
			})
			// j, _ := json.Marshal(rr)
			// fmt.Println("XZ2: ", err, string(j))
			continue
		}
		rm.executionNodes[reqRef] = msg.Executor
	}
}

func (rm *resultsMatcher) AddUnwantedResponse(ctx context.Context, msg insolar.Message) {
	response := msg.(*message.ReturnResults)
	inslogger.FromContext(ctx).Warn("IP1: Receive UnwantedResponse", response.Reason)

	if response.Reason.IsEmpty() {
		panic("EMPTY REASON")
	}

	rm.lock.Lock()
	defer rm.lock.Unlock()
	if node, ok := rm.executionNodes[response.Reason]; ok {
		// response.Target = *node
		// todo maybe call rm.mb.Send in goroutine
		// todo check errors? retry?
		// rm.mb.Send(ctx, response, nil)
		inslogger.FromContext(ctx).Warn("IP1: Send UnwantedResponse", response.Reason)
		// go
		if rm.lr.MessageBus == nil {
			panic("NO MESSAGE BUS")
		}
		rm.lr.MessageBus.Send(ctx, response, &insolar.MessageSendOptions{
			Receiver: &node,
		})
		return
	}
	rm.unwantedResponses[response.Reason] = *response
}

func (rm *resultsMatcher) Clear() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.executionNodes = make(map[insolar.Reference]insolar.Reference)
	rm.unwantedResponses = make(map[insolar.Reference]message.ReturnResults)
}

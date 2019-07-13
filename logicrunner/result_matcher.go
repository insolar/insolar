package logicrunner

import (
	"sync"

	"context"

	"fmt"

	"encoding/json"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
)

type resultsMatcher struct {
	mb                insolar.MessageBus
	lock              *sync.RWMutex
	executionNodes    map[insolar.Reference]insolar.Reference
	unwantedResponses map[insolar.Reference]message.ReturnResults
}

func NewResultsMatcher(mb insolar.MessageBus) resultsMatcher {
	return resultsMatcher{
		mb:                mb,
		lock:              &sync.RWMutex{},
		executionNodes:    make(map[insolar.Reference]insolar.Reference),
		unwantedResponses: make(map[insolar.Reference]message.ReturnResults),
	}
}

func (rm *resultsMatcher) AddStillExecution(ctx context.Context, msg *message.StillExecuting) {
	fmt.Println("IP1: Receive StillExecution", msg.RequestRefs)
	rm.lock.Lock()
	defer rm.lock.Unlock()
	for _, reqRef := range msg.RequestRefs {
		if response, ok := rm.unwantedResponses[reqRef]; ok {
			// response.Target = *node
			// todo maybe call rm.mb.Send in goroutine
			// todo check errors? retry?
			fmt.Println("IP1: Send StillExecution", reqRef)
			rr, err := rm.mb.Send(ctx, &response, &insolar.MessageSendOptions{
				Receiver: &msg.Executor,
			})
			j, _ := json.Marshal(rr)
			fmt.Println("XZ2: ", err, string(j))
			continue
		}
		rm.executionNodes[reqRef] = msg.Executor
	}
}

func (rm *resultsMatcher) AddUnwantedResponse(ctx context.Context, msg insolar.Message) {
	response := msg.(*message.ReturnResults)
	fmt.Println("IP1: Receive UnwantedResponse", response.RequestRef)

	rm.lock.Lock()
	defer rm.lock.Unlock()
	if node, ok := rm.executionNodes[response.RequestRef]; ok {
		// response.Target = *node
		// todo maybe call rm.mb.Send in goroutine
		// todo check errors? retry?
		// rm.mb.Send(ctx, response, nil)
		fmt.Println("IP1: Send UnwantedResponse", response.RequestRef)
		rr, err := rm.mb.Send(ctx, response, &insolar.MessageSendOptions{
			Receiver: &node,
		})
		j, _ := json.Marshal(rr)
		fmt.Println("XZ2: ", err, string(j))
		return
	}
	rm.unwantedResponses[response.RequestRef] = *response
}

func (rm *resultsMatcher) Clear() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.executionNodes = make(map[insolar.Reference]insolar.Reference)
	rm.unwantedResponses = make(map[insolar.Reference]message.ReturnResults)
}

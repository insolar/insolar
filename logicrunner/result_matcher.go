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

package logicrunner

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/messagebus"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ResultMatcher -o ./ -s _mock.go

type ResultMatcher interface {
	AddStillExecution(ctx context.Context, msg *message.StillExecuting)
	AddUnwantedResponse(ctx context.Context, msg *message.ReturnResults) error
	Clear()
}

type resultWithTraceID struct {
	traceID string
	result  message.ReturnResults
}

type resultsMatcher struct {
	lr                *LogicRunner
	lock              sync.RWMutex
	executionNodes    map[insolar.Reference]insolar.Reference
	unwantedResponses map[insolar.Reference]resultWithTraceID
}

func newResultsMatcher(lr *LogicRunner) *resultsMatcher {
	return &resultsMatcher{
		lr:                lr,
		lock:              sync.RWMutex{},
		executionNodes:    make(map[insolar.Reference]insolar.Reference),
		unwantedResponses: make(map[insolar.Reference]resultWithTraceID),
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
			ctx = inslogger.ContextWithTrace(ctx, response.traceID)
			inslogger.FromContext(ctx).Debug("[ resultsMatcher::AddStillExecution ] resend unwanted response ", reqRef)
			go rm.send(ctx, &response.result, &msg.Executor)
		}
		rm.executionNodes[reqRef] = msg.Executor
	}
}

func (rm *resultsMatcher) AddUnwantedResponse(ctx context.Context, msg *message.ReturnResults) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	object := *msg.Target.Record()

	err := rm.isStillExecutor(ctx, object)
	if err != nil {
		return err
	}

	if node, ok := rm.executionNodes[msg.Reason]; ok {
		inslogger.FromContext(ctx).Debug("[ resultsMatcher::AddUnwantedResponse ] resend unwanted response ", msg.Reason)
		go rm.send(ctx, msg, &node)
		delete(rm.unwantedResponses, msg.Reason)
		return nil
	}
	rm.unwantedResponses[msg.Reason] = resultWithTraceID{utils.TraceID(ctx), *msg}

	return rm.isStillExecutor(ctx, object)
}

// isStillExecutor is tmp solution. Needs to be moved on flow
func (rm *resultsMatcher) isStillExecutor(ctx context.Context, object insolar.ID) error {
	pulse, err := rm.lr.PulseAccessor.Latest(ctx)
	if err != nil {
		return flow.ErrCancelled
	}
	node, err := rm.lr.JetCoordinator.VirtualExecutorForObject(ctx, object, pulse.PulseNumber)
	if err != nil {
		return flow.ErrCancelled
	}
	if *node != rm.lr.JetCoordinator.Me() {
		return flow.ErrCancelled
	}
	return nil
}

func (rm *resultsMatcher) Clear() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.executionNodes = make(map[insolar.Reference]insolar.Reference)
	rm.unwantedResponses = make(map[insolar.Reference]resultWithTraceID)
}

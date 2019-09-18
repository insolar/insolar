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

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/metrics"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ResultMatcher -o ./ -s _mock.go -g

type ResultMatcher interface {
	AddStillExecution(ctx context.Context, msg payload.StillExecuting)
	AddUnwantedResponse(ctx context.Context, msg payload.ReturnResults)
	Clear(ctx context.Context)
}

type resultWithContext struct {
	ctx    context.Context
	result payload.ReturnResults
}

type resultsMatcher struct {
	sender        bus.Sender
	pulseAccessor pulse.Accessor

	lock              sync.Mutex
	executionNodes    map[insolar.Reference]insolar.Reference
	unwantedResponses map[insolar.Reference]resultWithContext
}

func newResultsMatcher(sender bus.Sender, pa pulse.Accessor) *resultsMatcher {
	return &resultsMatcher{
		sender:            sender,
		pulseAccessor:     pa,
		executionNodes:    make(map[insolar.Reference]insolar.Reference),
		unwantedResponses: make(map[insolar.Reference]resultWithContext),
	}
}

func (rm *resultsMatcher) AddStillExecution(ctx context.Context, msg payload.StillExecuting) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	inslogger.FromContext(ctx).Debug("got pendings confirmation")

	for _, reqRef := range msg.RequestRefs {
		if response, ok := rm.unwantedResponses[reqRef]; ok {
			ctx := response.ctx
			go rm.send(ctx, response.result, msg.Executor)
			delete(rm.unwantedResponses, reqRef)
		}
		rm.executionNodes[reqRef] = msg.Executor
	}
}

func (rm *resultsMatcher) AddUnwantedResponse(ctx context.Context, msg payload.ReturnResults) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	inslogger.FromContext(ctx).Debug("got unwanted response to request ", msg.RequestRef.String())

	if node, ok := rm.executionNodes[msg.Reason]; ok {
		go rm.send(ctx, msg, node)
		return
	}
	rm.unwantedResponses[msg.Reason] = resultWithContext{
		ctx:    ctx,
		result: msg,
	}
}

func (rm *resultsMatcher) Clear(ctx context.Context) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	rm.executionNodes = make(map[insolar.Reference]insolar.Reference)

	logger := inslogger.FromContext(ctx)
	stats.Record(ctx, metrics.ResultsMatcherDroppedResults.M(int64(len(rm.unwantedResponses))))
	for reqRef := range rm.unwantedResponses {
		logger.Warn("not claimed response to request ", reqRef.String(), ", not confirmed pending?")
	}
	rm.unwantedResponses = make(map[insolar.Reference]resultWithContext)
}

func (rm *resultsMatcher) send(ctx context.Context, msg payload.ReturnResults, receiver insolar.Reference) {
	stats.Record(ctx, metrics.ResultsMatcherSentResults.M(1))
	logger := inslogger.FromContext(ctx)

	logger.Debug("resending result of request ", msg.RequestRef.String(), " to ", receiver.String())

	msgData, err := payload.NewResultMessage(&msg)
	if err != nil {
		inslogger.FromContext(ctx).Debug("failed to serialize message")
		return
	}

	sender := bus.NewWaitOKWithRetrySender(rm.sender, rm.pulseAccessor, 1)
	sender.SendTarget(ctx, msgData, receiver)
}

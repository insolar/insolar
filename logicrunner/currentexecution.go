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
	"github.com/insolar/insolar/insolar/record"
)

type CurrentExecution struct {
	Context       context.Context
	LogicContext  *insolar.LogicCallContext
	RequestRef    *Ref
	Request       *record.Request
	RequesterNode *Ref
	SentResult    bool
	Nonce         uint64
	Deactivate    bool

	OutgoingRequests []OutgoingRequest
}

type OutgoingRequest struct {
	Request   record.Request
	NewObject *Ref
	Response  []byte
	Error     error
}

func (ce *CurrentExecution) AddOutgoingRequest(
	ctx context.Context, request record.Request, result []byte, newObject *Ref, err error,
) {
	rec := OutgoingRequest{
		Request: request,
		Response: result,
		NewObject: newObject,
		Error: err,
	}
	ce.OutgoingRequests = append(ce.OutgoingRequests, rec)
}

type CurrentExecutionList struct {
	lock       sync.RWMutex
	executions map[insolar.Reference]*CurrentExecution
}

func (ces *CurrentExecutionList) Get(requestRef insolar.Reference) *CurrentExecution {
	ces.lock.RLock()
	rv := ces.executions[requestRef]
	ces.lock.RUnlock()
	return rv
}

func (ces *CurrentExecutionList) Set(requestRef insolar.Reference, ce *CurrentExecution) {
	ces.lock.Lock()
	ces.executions[requestRef] = ce
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) Delete(requestRef insolar.Reference) {
	ces.lock.Lock()
	delete(ces.executions, requestRef)
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) GetByTraceID(traceid string) *CurrentExecution {
	ces.lock.RLock()
	for _, ce := range ces.executions {
		if ce.LogicContext.TraceID == traceid {
			ces.lock.RUnlock()
			return ce
		}
	}
	ces.lock.RUnlock()
	return nil
}

func (ces *CurrentExecutionList) GetMutable() *CurrentExecution {
	ces.lock.RLock()
	for _, ce := range ces.executions {
		if !ce.LogicContext.Immutable {
			ces.lock.RUnlock()
			return ce
		}
	}
	ces.lock.RUnlock()
	return nil
}

func (ces *CurrentExecutionList) Cleanup() {
	ces.lock.Lock()
	ces.executions = make(map[insolar.Reference]*CurrentExecution)
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) Length() int {
	ces.lock.RLock()
	rv := len(ces.executions)
	ces.lock.RUnlock()
	return rv
}

func (ces *CurrentExecutionList) Empty() bool {
	return ces.Length() == 0
}

type CurrentExecutionPredicate func(*CurrentExecution, interface{}) bool

func (ces *CurrentExecutionList) Check(predicate CurrentExecutionPredicate, args interface{}) bool {
	rv := true
	ces.lock.RLock()
	for _, current := range ces.executions {
		if !predicate(current, args) {
			rv = false
			break
		}
	}
	ces.lock.RUnlock()
	return rv
}

func NewCurrentExecutionList() *CurrentExecutionList {
	rv := &CurrentExecutionList{}
	rv.Cleanup()
	return rv
}

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
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type Transcript struct {
	State interface{} // Shows current execution status of task

	Context          context.Context
	LogicContext     *insolar.LogicCallContext
	Request          *record.Request
	RequestRef       *insolar.Reference
	RequesterNode    *insolar.Reference
	SentResult       bool
	Nonce            uint64
	Deactivate       bool
	OutgoingRequests []OutgoingRequest

	Parcel     insolar.Parcel
	FromLedger bool
}

func NewTranscript(ctx context.Context, parcel insolar.Parcel, requestRef *insolar.Reference,
	pulse *insolar.Pulse, callee insolar.Reference) *Transcript {

	msg := parcel.Message().(*message.CallMethod)

	logicalContext := &insolar.LogicCallContext{
		Mode:            insolar.ExecuteCallMode,
		Caller:          msg.GetCaller(),
		Callee:          &callee,
		Request:         requestRef,
		Time:            time.Now(), // TODO: probably we should take it earlier
		Pulse:           *pulse,
		TraceID:         inslogger.TraceID(ctx),
		CallerPrototype: &msg.CallerPrototype,
		Immutable:       msg.Immutable,
	}
	sender := parcel.GetSender()

	return &Transcript{
		Context:       ctx,
		LogicContext:  logicalContext,
		Request:       &msg.Request,
		RequestRef:    requestRef,
		RequesterNode: &sender,
		SentResult:    false,
		Nonce:         0,
		Deactivate:    false,

		Parcel:     parcel,
		FromLedger: false,
	}
}

type OutgoingRequest struct {
	Request   record.Request
	NewObject *Ref
	Response  []byte
	Error     error
}

func (ce *Transcript) AddOutgoingRequest(
	ctx context.Context, request record.Request, result []byte, newObject *Ref, err error,
) {
	rec := OutgoingRequest{
		Request:   request,
		Response:  result,
		NewObject: newObject,
		Error:     err,
	}
	ce.OutgoingRequests = append(ce.OutgoingRequests, rec)
}

type CurrentExecutionList struct {
	lock       sync.RWMutex
	executions map[insolar.Reference]*Transcript
}

func (ces *CurrentExecutionList) Get(requestRef insolar.Reference) *Transcript {
	ces.lock.RLock()
	rv := ces.executions[requestRef]
	ces.lock.RUnlock()
	return rv
}

func (ces *CurrentExecutionList) Set(requestRef insolar.Reference, ce *Transcript) {
	ces.lock.Lock()
	ces.executions[requestRef] = ce
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) Delete(requestRef insolar.Reference) {
	ces.lock.Lock()
	delete(ces.executions, requestRef)
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) GetByTraceID(traceid string) *Transcript {
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

func (ces *CurrentExecutionList) GetMutable() *Transcript {
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
	ces.executions = make(map[insolar.Reference]*Transcript)
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

type CurrentExecutionPredicate func(*Transcript, interface{}) bool

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

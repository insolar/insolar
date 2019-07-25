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
	"errors"
	"reflect"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type OutgoingRequest struct {
	Request   record.IncomingRequest
	NewObject *Ref
	Response  []byte
	Error     error
}

type Transcript struct {
	ObjectDescriptor artifacts.ObjectDescriptor
	Context          context.Context
	LogicContext     *insolar.LogicCallContext
	Request          *record.IncomingRequest
	RequestRef       insolar.Reference
	RequesterNode    *insolar.Reference
	Nonce            uint64
	Deactivate       bool
	OutgoingRequests []OutgoingRequest
	FromLedger       bool
}

func NewTranscript(
	ctx context.Context,
	requestRef insolar.Reference,
	request record.IncomingRequest,
) *Transcript {

	return &Transcript{
		Context:    ctx,
		Request:    &request,
		RequestRef: requestRef,
		Nonce:      0,
		Deactivate: false,

		FromLedger: false,
	}
}

func (t *Transcript) AddOutgoingRequest(
	ctx context.Context, request record.IncomingRequest, result []byte, newObject *Ref, err error,
) {
	rec := OutgoingRequest{
		Request:   request,
		Response:  result,
		NewObject: newObject,
		Error:     err,
	}
	t.OutgoingRequests = append(t.OutgoingRequests, rec)
}

func (t *Transcript) HasOutgoingRequest(
	ctx context.Context, request record.IncomingRequest,
) *OutgoingRequest {
	for i := range t.OutgoingRequests {
		if reflect.DeepEqual(t.OutgoingRequests[i].Request, request) {
			return &t.OutgoingRequests[i]
		}
	}
	return nil
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

func (ces *CurrentExecutionList) SetOnce(t *Transcript) error {
	ces.lock.Lock()
	defer ces.lock.Unlock()

	if _, has := ces.executions[t.RequestRef]; has {
		return errors.New("not setting, already in the set")
	}

	ces.executions[t.RequestRef] = t
	return nil
}

func (ces *CurrentExecutionList) Delete(requestRef insolar.Reference) {
	ces.lock.Lock()
	delete(ces.executions, requestRef)
	ces.lock.Unlock()
}

func (ces *CurrentExecutionList) GetByTraceID(traceid string) *Transcript {
	ces.lock.RLock()
	defer ces.lock.RUnlock()
	for _, ce := range ces.executions {
		if ce.LogicContext.TraceID == traceid {
			return ce
		}
	}
	return nil
}

func (ces *CurrentExecutionList) GetMutable() *Transcript {
	ces.lock.RLock()
	for _, ce := range ces.executions {
		if !ce.Request.Immutable {
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

func (ces *CurrentExecutionList) Has(requestRef insolar.Reference) bool {
	ces.lock.RLock()
	defer ces.lock.RUnlock()
	_, has := ces.executions[requestRef]
	return has
}

func (ces *CurrentExecutionList) GetAllRequestRefs() []insolar.Reference {
	ces.lock.RLock()
	defer ces.lock.RUnlock()
	out := make([]insolar.Reference, len(ces.executions))
	i := 0
	for key := range ces.executions {
		out[i] = key
		i++
	}
	return out
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

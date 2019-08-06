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

package currentexecutor

import (
	"errors"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/executionarchive"
	"github.com/insolar/insolar/logicrunner/transcript"
)

type CurrentExecutionList struct {
	lock       sync.RWMutex
	executions map[insolar.Reference]*transcript.Transcript
}

func (ces *CurrentExecutionList) Get(requestRef insolar.Reference) *transcript.Transcript {
	ces.lock.RLock()
	rv := ces.executions[requestRef]
	ces.lock.RUnlock()
	return rv
}

func (ces *CurrentExecutionList) SetOnce(t *transcript.Transcript) error {
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

func (ces *CurrentExecutionList) GetByTraceID(traceid string) *transcript.Transcript {
	ces.lock.RLock()
	defer ces.lock.RUnlock()
	for _, ce := range ces.executions {
		if ce.LogicContext.TraceID == traceid {
			return ce
		}
	}
	return nil
}

func (ces *CurrentExecutionList) GetMutable() *transcript.Transcript {
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
	ces.executions = make(map[insolar.Reference]*transcript.Transcript)
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

func (ces *CurrentExecutionList) Archive(archiver executionarchive.Archiver) {
	ces.lock.RLock()
	defer ces.lock.RUnlock()

	for _, current := range ces.executions {
		archiver.Archive(current)
	}
}

func NewCurrentExecutionList() *CurrentExecutionList {
	rv := &CurrentExecutionList{}
	rv.Cleanup()
	return rv
}

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

package currentexecution

import (
	"errors"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionarchive"
)

type List struct {
	lock       sync.RWMutex
	executions map[insolar.Reference]*common.Transcript
}

func (ces *List) Get(requestRef insolar.Reference) *common.Transcript {
	ces.lock.RLock()
	rv := ces.executions[requestRef]
	ces.lock.RUnlock()
	return rv
}

func (ces *List) SetOnce(t *common.Transcript) error {
	ces.lock.Lock()
	defer ces.lock.Unlock()

	if _, has := ces.executions[t.RequestRef]; has {
		return errors.New("not setting, already in the set")
	}

	ces.executions[t.RequestRef] = t
	return nil
}

func (ces *List) Delete(requestRef insolar.Reference) {
	ces.lock.Lock()
	defer ces.lock.Unlock()

	delete(ces.executions, requestRef)
}

func (ces *List) GetByTraceID(traceid string) *common.Transcript {
	ces.lock.RLock()
	defer ces.lock.RUnlock()
	for _, ce := range ces.executions {
		if ce.LogicContext.TraceID == traceid {
			return ce
		}
	}
	return nil
}

func (ces *List) GetMutable() *common.Transcript {
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

func (ces *List) Cleanup() {
	ces.lock.Lock()
	defer ces.lock.Unlock()

	ces.executions = make(map[insolar.Reference]*common.Transcript)
}

func (ces *List) Length() int {
	ces.lock.RLock()
	rv := len(ces.executions)
	ces.lock.RUnlock()
	return rv
}

func (ces *List) Empty() bool {
	return ces.Length() == 0
}

func (ces *List) Has(requestRef insolar.Reference) bool {
	ces.lock.RLock()
	defer ces.lock.RUnlock()
	_, has := ces.executions[requestRef]
	return has
}

func (ces *List) GetAllRequestRefs() []insolar.Reference {
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

func (ces *List) Archive(archiver executionarchive.Archiver) {
	ces.lock.RLock()
	defer ces.lock.RUnlock()

	for _, current := range ces.executions {
		archiver.Archive(current)
	}
}

func NewList() *List {
	rv := &List{}
	rv.Cleanup()
	return rv
}

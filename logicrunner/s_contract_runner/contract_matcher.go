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

package s_contract_runner

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CellID uuid.UUID

var (
	ErrUnmatchedResponse   = errors.New("no one awaits response")
	ErrAlreadyHaveResponse = errors.New("already got response")
)

type Matcher interface {
	// Feed informs running contract that we got a response
	// in case when no one awaits response - UnmatcherResponse error is returned
	Feed(id CellID, result interface{}) error

	// Register creates unique id and tells that we're awaiting for result on that id
	Register() CellID
	// Real (blocking) wait for result on
	Wait(id CellID) (interface{}, bool)
	// Exclude destroys waiting cell (if it exists or not)
	Exclude(id CellID)

	// Cleanup destroys all waiting cells and wake up all who awaits response
	Cleanup()
}

func newAwaitingCell(id CellID) *awaitingCell {
	return &awaitingCell{
		id:       id,
		response: make(chan interface{}, 1),
	}
}

type awaitingCell struct {
	id       CellID
	lock     sync.Mutex
	response chan interface{}
}

// NewResponseMatcher allocates basic result await service
func NewResponseMatcher() Matcher {
	return &responseMatcher{
		waitingList: make(map[CellID]*awaitingCell),
	}
}

type responseMatcher struct {
	lock        sync.Mutex
	waitingList map[CellID]*awaitingCell
}

func (r *responseMatcher) Feed(id CellID, result interface{}) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	rv, ok := r.waitingList[id]

	if !ok {
		return ErrUnmatchedResponse
	}

	select {
	case rv.response <- result:
		return nil
	default:
		return ErrAlreadyHaveResponse
	}
}

func (r *responseMatcher) Register() CellID {
	r.lock.Lock()
	defer r.lock.Unlock()

	// generate unique (non-used) id
	id := CellID(uuid.New())
	for {
		if _, ok := r.waitingList[id]; !ok {
			break
		}
	}

	r.waitingList[id] = newAwaitingCell(id)
	return id
}

func (r *responseMatcher) Wait(id CellID) (interface{}, bool) {
	r.lock.Lock()
	rv, ok := r.waitingList[id]
	r.lock.Unlock()

	if !ok {
		return nil, false
	}

	select {
	case rv, ok := <-rv.response:
		if !ok {
			return nil, false
		}
		r.Exclude(id)
		return rv, true
	default:
		return nil, false
	}
}

func (r *responseMatcher) Exclude(id CellID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	rv, ok := r.waitingList[id]
	if !ok {
		return
	}

	delete(r.waitingList, id)
	close(rv.response)
}

func (r *responseMatcher) Cleanup() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for key := range r.waitingList {
		r.Exclude(key)
	}
}

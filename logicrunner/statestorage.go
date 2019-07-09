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
	"sync"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.StateStorage -o ./ -s _mock.go
type StateStorage interface {
	sync.Locker

	GetObjectState(ref Ref) *ObjectState
	UpsertObjectState(ref Ref) *ObjectState
	MustObjectState(ref Ref) *ObjectState
	GetExecutionState(ref Ref) *ExecutionState
	DeleteObjectState(ref Ref)
	StateMap() *map[Ref]*ObjectState
}

type stateStorage struct {
	sync.RWMutex
	state map[Ref]*ObjectState // if object exists, we are validating or executing it right now
}

func NewStateStorage() StateStorage {
	ss := &stateStorage{
		state: make(map[Ref]*ObjectState),
	}
	return ss
}

func (ss *stateStorage) GetObjectState(ref Ref) *ObjectState {
	ss.RLock()
	res, ok := ss.state[ref]
	ss.RUnlock()
	if !ok {
		return nil
	}
	return res
}

func (ss *stateStorage) UpsertObjectState(ref Ref) *ObjectState {
	ss.RLock()
	if res, ok := ss.state[ref]; ok {
		ss.RUnlock()
		return res
	}
	ss.RUnlock()

	ss.Lock()
	defer ss.Unlock()
	if _, ok := ss.state[ref]; !ok {
		ss.state[ref] = &ObjectState{}
	}
	return ss.state[ref]
}

func (ss *stateStorage) MustObjectState(ref Ref) *ObjectState {
	res := ss.GetObjectState(ref)
	if res == nil {
		panic("No requested object state. ref: " + ref.String())
	}
	return res
}

func (ss *stateStorage) GetExecutionState(ref Ref) *ExecutionState {
	os := ss.GetObjectState(ref)
	if os == nil {
		return nil
	}

	os.Lock()
	defer os.Unlock()
	return os.ExecutionState
}

func (ss *stateStorage) DeleteObjectState(ref Ref) {
	delete(ss.state, ref)
}

func (ss *stateStorage) StateMap() *map[Ref]*ObjectState {
	return &ss.state
}

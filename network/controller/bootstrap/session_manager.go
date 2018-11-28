/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package bootstrap

import (
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/utils"
)

type SessionID uint64
type SessionState uint8

const (
	SessionStarted SessionState = iota + 1
)

type Session struct {
	NodeID core.RecordRef
	Cert   core.Certificate
	State  SessionState

	// TODO: expiry time
}

type SessionManager struct {
	sequence uint64
	lock     sync.RWMutex
	sessions map[SessionID]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{sessions: make(map[SessionID]*Session)}
}

func (sm *SessionManager) NewSession(ref core.RecordRef) *Session {
	id := utils.AtomicLoadAndIncrementUint64(&sm.sequence)
	result := &Session{NodeID: ref, State: SessionStarted}
	sm.lock.Lock()
	sm.sessions[SessionID(id)] = result
	sm.lock.Unlock()
	return result
}

func (sm *SessionManager) GetSession(id SessionID) *Session {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	return sm.sessions[id]
}

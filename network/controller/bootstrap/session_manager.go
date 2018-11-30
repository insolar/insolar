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
	"fmt"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

type SessionID uint64

//go:generate stringer -type=SessionState
type SessionState uint8

const (
	Authorized SessionState = iota + 1
	Challenge1
	Challenge2
)

type Session struct {
	NodeID core.RecordRef
	Cert   core.AuthorizationCertificate
	State  SessionState

	DiscoveryNonce Nonce

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

func (sm *SessionManager) NewSession(ref core.RecordRef, cert core.AuthorizationCertificate) SessionID {
	id := utils.AtomicLoadAndIncrementUint64(&sm.sequence)
	result := &Session{NodeID: ref, State: Authorized, Cert: cert}
	sessionID := SessionID(id)
	sm.lock.Lock()
	sm.sessions[sessionID] = result
	sm.lock.Unlock()
	return sessionID
}

func (sm *SessionManager) CheckSession(id SessionID, expected SessionState) error {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	_, err := sm.checkSession(id, expected)
	return err
}

func (sm *SessionManager) checkSession(id SessionID, expected SessionState) (*Session, error) {
	session := sm.sessions[id]
	if session == nil {
		return nil, errors.New(fmt.Sprintf("no such session ID: %d", id))
	}
	if session.State != expected {
		return nil, errors.New(fmt.Sprintf("session %d should have state %s but has %s", id, expected, session.State))
	}
	return session, nil
}

func (sm *SessionManager) SetDiscoveryNonce(id SessionID, discoveryNonce Nonce) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	session, err := sm.checkSession(id, Authorized)
	if err != nil {
		return err
	}
	session.DiscoveryNonce = discoveryNonce
	session.State = Challenge1
	return nil
}

func (sm *SessionManager) GetChallengeData(id SessionID) (core.AuthorizationCertificate, Nonce, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	session, err := sm.checkSession(id, Challenge1)
	if err != nil {
		return nil, nil, err
	}
	return session.Cert, session.DiscoveryNonce, nil
}

func (sm *SessionManager) ChallengePassed(id SessionID) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	session, err := sm.checkSession(id, Challenge1)
	if err != nil {
		return err
	}
	session.State = Challenge2
	return nil
}

func (sm *SessionManager) ReleaseSession(id SessionID) (*Session, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	session, err := sm.checkSession(id, Challenge2)
	if err != nil {
		return nil, err
	}
	delete(sm.sessions, id)
	return session, nil
}

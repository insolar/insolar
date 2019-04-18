//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package bootstrap

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/utils"
)

type SessionID uint64

//go:generate stringer -type=SessionState
type SessionState uint8

const (
	Authorized SessionState = iota + 1
	Challenge1
	Challenge2
)

const (
	stateRunning = uint32(iota + 1)
	stateIdle
)

type Session struct {
	NodeID insolar.Reference
	Cert   insolar.AuthorizationCertificate
	State  SessionState

	DiscoveryNonce Nonce

	Time time.Time
	TTL  time.Duration
}

func (s *Session) expirationTime() time.Time {
	return s.Time.Add(s.TTL)
}

type sessionWithID struct {
	*Session
	SessionID
}

type notification struct{}

type SessionManager interface {
	component.Starter
	component.Stopper

	NewSession(ref insolar.Reference, cert insolar.AuthorizationCertificate, ttl time.Duration) SessionID
	CheckSession(id SessionID, expected SessionState) error
	SetDiscoveryNonce(id SessionID, discoveryNonce Nonce) error
	GetChallengeData(id SessionID) (insolar.AuthorizationCertificate, Nonce, error)
	ChallengePassed(id SessionID) error
	ReleaseSession(id SessionID) (*Session, error)
	ProlongateSession(id SessionID, session *Session)
}

type sessionManager struct {
	sequence uint64
	lock     sync.RWMutex
	sessions map[SessionID]*Session
	state    uint32

	sessionsChangeNotification chan notification
	stopCleanupNotification    chan notification
}

func NewSessionManager() SessionManager {
	return &sessionManager{
		sessions:                   make(map[SessionID]*Session),
		sessionsChangeNotification: make(chan notification),
		stopCleanupNotification:    make(chan notification),
		state:                      stateIdle,
	}
}

func (sm *sessionManager) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[ sessionManager::Start ] start cleaning up sessions")

	if atomic.CompareAndSwapUint32(&sm.state, stateIdle, stateRunning) {
		go sm.cleanupExpiredSessions()
	} else {
		logger.Warn("[ sessionManager::Start ] Called twice")
	}

	return nil
}

func (sm *sessionManager) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[ sessionManager::Stop ] stop cleaning up sessions")

	if atomic.CompareAndSwapUint32(&sm.state, stateRunning, stateIdle) {
		sm.stopCleanupNotification <- notification{}
	} else {
		logger.Warn("[ sessionManager::Stop ] Called twice")
	}

	return nil
}

func (sm *sessionManager) NewSession(ref insolar.Reference, cert insolar.AuthorizationCertificate, ttl time.Duration) SessionID {
	id := utils.AtomicLoadAndIncrementUint64(&sm.sequence)
	session := &Session{
		NodeID: ref,
		State:  Authorized,
		Cert:   cert,
		Time:   time.Now(),
		TTL:    ttl,
	}
	sessionID := SessionID(id)
	sm.addSession(sessionID, session)
	return sessionID
}

func (sm *sessionManager) addSession(id SessionID, session *Session) {
	sm.lock.Lock()
	sm.sessions[id] = session
	sm.lock.Unlock()

	sm.sessionsChangeNotification <- notification{}
}

func (sm *sessionManager) CheckSession(id SessionID, expected SessionState) error {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	_, err := sm.checkSession(id, expected)
	return err
}

func (sm *sessionManager) checkSession(id SessionID, expected SessionState) (*Session, error) {
	session := sm.sessions[id]
	if session == nil {
		return nil, errors.New(fmt.Sprintf("no such session ID: %d", id))
	}
	if session.State != expected {
		return nil, errors.New(fmt.Sprintf("session %d should have state %s but has %s", id, expected, session.State))
	}
	return session, nil
}

func (sm *sessionManager) SetDiscoveryNonce(id SessionID, discoveryNonce Nonce) error {
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

func (sm *sessionManager) GetChallengeData(id SessionID) (insolar.AuthorizationCertificate, Nonce, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	session, err := sm.checkSession(id, Challenge1)
	if err != nil {
		return nil, nil, err
	}
	return session.Cert, session.DiscoveryNonce, nil
}

func (sm *sessionManager) ChallengePassed(id SessionID) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	session, err := sm.checkSession(id, Challenge1)
	if err != nil {
		return err
	}
	session.State = Challenge2
	return nil
}

func (sm *sessionManager) ReleaseSession(id SessionID) (*Session, error) {
	sm.lock.Lock()

	session, err := sm.checkSession(id, Challenge2)
	if err != nil {
		sm.lock.Unlock()
		return nil, err
	}
	delete(sm.sessions, id)
	sm.lock.Unlock()

	sm.sessionsChangeNotification <- notification{}

	return session, nil
}

func (sm *sessionManager) ProlongateSession(id SessionID, session *Session) {
	session.Time = time.Now()
	sm.addSession(id, session)
}

func (sm *sessionManager) cleanupExpiredSessions() {
	var sessionsByExpirationTime []*sessionWithID
	for {
		sm.lock.RLock()
		sessionsCount := len(sm.sessions)
		sm.lock.RUnlock()

		// We missed notification
		if sessionsCount != 0 && len(sessionsByExpirationTime) == 0 {
			sessionsByExpirationTime = sm.sortSessionsByExpirationTime()
		}

		// Session count is zero - wait for first session added.
		if len(sessionsByExpirationTime) == 0 {
			// Have no active sessions.
			// Block until sessions will be added or session manager begins to stop.
			select {
			case <-sm.sessionsChangeNotification:
				sessionsByExpirationTime = sm.sortSessionsByExpirationTime()
			case <-sm.stopCleanupNotification:
				return
			}

			// Check session instantly released concurrently
			if len(sessionsByExpirationTime) == 0 {
				continue
			}
		}

		// Get expiration time for next session and wait for it
		nextSessionToExpire := sessionsByExpirationTime[0]
		waitTime := time.Until(nextSessionToExpire.expirationTime())

		select {
		case <-sm.sessionsChangeNotification:
			// Handle new session. reorder expiration short list
			sessionsByExpirationTime = sm.sortSessionsByExpirationTime()

		case <-time.After(waitTime):
			// Move forward through sessions and check whether we should delete the session
			sessionsByExpirationTime = sm.expireSessions(sessionsByExpirationTime)

		case <-sm.stopCleanupNotification:
			return
		}
	}
}

func (sm *sessionManager) sortSessionsByExpirationTime() []*sessionWithID {
	sm.lock.RLock()

	// Read active session with their ids. We have to store them as a slice to keep ordering by expiration time.
	sessionsByExpirationTime := make([]*sessionWithID, 0, len(sm.sessions))
	for sessionID, session := range sm.sessions {
		sessionsByExpirationTime = append(sessionsByExpirationTime, &sessionWithID{
			SessionID: sessionID,
			Session:   session,
		})
	}

	sm.lock.RUnlock()

	sort.SliceStable(sessionsByExpirationTime, func(i, j int) bool {
		expirationTime1 := sessionsByExpirationTime[i].expirationTime()
		expirationTime2 := sessionsByExpirationTime[j].expirationTime()

		return expirationTime1.Before(expirationTime2)
	})

	return sessionsByExpirationTime
}

func (sm *sessionManager) expireSessions(sessionsByExpirationTime []*sessionWithID) []*sessionWithID {
	var shift int

	sm.lock.Lock()

	for i, session := range sessionsByExpirationTime {
		// Check when we have to stop expire
		if session.expirationTime().After(time.Now()) {
			break
		}

		delete(sm.sessions, session.SessionID)
		shift = i + 1
	}

	sm.lock.Unlock()

	return sessionsByExpirationTime[shift:]
}

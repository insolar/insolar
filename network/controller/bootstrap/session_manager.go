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
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
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

type SessionManager struct {
	sequence uint64
	lock     sync.RWMutex
	sessions map[SessionID]*Session

	newSessionNotification  chan notification
	stopCleanupNotification chan notification
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:                make(map[SessionID]*Session),
		newSessionNotification:  make(chan notification),
		stopCleanupNotification: make(chan notification),
	}
}

func (sm *SessionManager) Start(ctx context.Context) error {
	inslogger.FromContext(ctx).Debug("[ SessionManager::Start ] start cleaning up sessions")

	go sm.cleanupExpiredSessions()

	return nil
}

func (sm *SessionManager) Stop(ctx context.Context) error {
	inslogger.FromContext(ctx).Debug("[ SessionManager::Stop ] stop cleaning up sessions")

	sm.stopCleanupNotification <- notification{}

	return nil
}

func (sm *SessionManager) NewSession(ref core.RecordRef, cert core.AuthorizationCertificate, ttl time.Duration) SessionID {
	id := utils.AtomicLoadAndIncrementUint64(&sm.sequence)
	session := &Session{
		NodeID: ref,
		State:  Authorized,
		Cert:   cert,
		Time:   time.Now(),
		TTL:    ttl,
	}
	sessionID := SessionID(id)

	sm.lock.Lock()
	sm.sessions[sessionID] = session
	sm.lock.Unlock()

	sm.newSessionNotification <- notification{}

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

func (sm *SessionManager) cleanupExpiredSessions() {
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
			// Have no active sessions. Block till sessions will be added.
			<-sm.newSessionNotification

			sessionsByExpirationTime = sm.sortSessionsByExpirationTime()

			// Check session instantly released concurrently
			if len(sessionsByExpirationTime) == 0 {
				continue
			}
		}

		// Get expiration time for next session and wait for it
		nextSessionToExpire := sessionsByExpirationTime[0]
		waitTime := time.Until(nextSessionToExpire.expirationTime())

		select {
		case <-sm.newSessionNotification:
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

func (sm *SessionManager) sortSessionsByExpirationTime() []*sessionWithID {
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

func (sm *SessionManager) expireSessions(sessionsByExpirationTime []*sessionWithID) []*sessionWithID {
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

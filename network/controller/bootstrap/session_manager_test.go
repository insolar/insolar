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
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sessionMapLen(sm *SessionManager) int {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	return len(sm.sessions)
}

func sessionMapDelete(sm *SessionManager, id SessionID) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	delete(sm.sessions, id)
}

func TestSessionManager_CleanupSimple(t *testing.T) {
	sm := NewSessionManager()

	err := sm.Start(context.Background())
	require.NoError(t, err)

	sm.NewSession(core.RecordRef{}, nil, time.Second)
	require.Equal(t, sessionMapLen(sm), 1)

	time.Sleep(1500 * time.Millisecond)
	assert.Equal(t, sessionMapLen(sm), 0)
}

func TestSessionManager_CleanupConcurrent(t *testing.T) {
	sm := NewSessionManager()

	err := sm.Start(context.Background())
	require.NoError(t, err)

	id := sm.NewSession(core.RecordRef{}, nil, time.Second)
	require.Equal(t, sessionMapLen(sm), 1)

	// delete session here and check nothing happened
	sessionMapDelete(sm, id)

	time.Sleep(1500 * time.Millisecond)
	assert.Equal(t, sessionMapLen(sm), 0)
}

func TestSessionManager_CleanupOrder(t *testing.T) {
	sm := NewSessionManager()

	err := sm.Start(context.Background())
	require.NoError(t, err)

	sm.NewSession(core.RecordRef{}, nil, 2*time.Second)
	sm.NewSession(core.RecordRef{}, nil, 2*time.Second)
	sm.NewSession(core.RecordRef{}, nil, time.Second)
	require.Equal(t, sessionMapLen(sm), 3)

	time.Sleep(1500 * time.Millisecond)
	assert.Equal(t, sessionMapLen(sm), 2)
}

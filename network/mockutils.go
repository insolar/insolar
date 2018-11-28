/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package network

import (
	"context"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
)

// MockPulseManager mock struct to read and write pulse that implements core.PulseManager interface.
type MockPulseManager struct {
	currentPulse core.Pulse
	callback     func(core.Pulse)
	mutex        sync.Mutex
}

type MockLedger struct {
	pm MockPulseManager
}

// GetLocalStorage returns local storage to work with.
func (l *MockLedger) GetLocalStorage() core.LocalStorage {
	panic("implement me")
}

// GetArtifactManager returns artifact manager to work with.
func (l *MockLedger) GetArtifactManager() core.ArtifactManager {
	return nil
}

// GetJetCoordinator returns jet coordinator to work with.
func (l *MockLedger) GetJetCoordinator() core.JetCoordinator {
	return nil
}

// GetPulseManager returns pulse manager to work with.
func (l *MockLedger) GetPulseManager() core.PulseManager {
	return &l.pm
}

func (pm *MockPulseManager) Current(context.Context) (*core.Pulse, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	return &pm.currentPulse, nil
}

func (pm *MockPulseManager) Set(ctx context.Context, pulse core.Pulse, dry bool) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.currentPulse = pulse
	if pm.callback != nil {
		pm.callback(pulse)
	}
	return nil
}

func (pm *MockPulseManager) SetCallback(callback func(core.Pulse)) {
	pm.callback = callback
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}

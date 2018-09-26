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

package mockutils

import (
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
)

// MockServiceConfiguration function to generate mock configuration for a ServiceNetwork node.
func MockServiceConfiguration(host string, bootstrapHosts []string, nodeID string) (configuration.HostNetwork, configuration.NodeNetwork) {
	transport := configuration.Transport{Protocol: "UTP", Address: host, BehindNAT: false}
	h := configuration.HostNetwork{
		Transport:      transport,
		IsRelay:        false,
		BootstrapHosts: bootstrapHosts,
	}

	n := configuration.NodeNetwork{Node: &configuration.Node{ID: nodeID}}

	return h, n
}

// WaitTimeout function to wait on a wait group with a timeout.
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

// MockPulseManager mock struct to read and write pulse that implements core.PulseManager interface.
type MockPulseManager struct {
	currentPulse core.Pulse
	callback     func(core.Pulse)
}

func (pm *MockPulseManager) Current() (*core.Pulse, error) {
	return &pm.currentPulse, nil
}

func (pm *MockPulseManager) Set(pulse core.Pulse) error {
	pm.currentPulse = pulse
	if pm.callback != nil {
		pm.callback(pulse)
	}
	return nil
}

func (pm *MockPulseManager) SetCallback(callback func(core.Pulse)) {
	pm.callback = callback
}

// MockLedger mock struct that implements core.Ledger interface.
type MockLedger struct {
	PM core.PulseManager
}

func (l *MockLedger) GetArtifactManager() core.ArtifactManager {
	return nil
}

func (l *MockLedger) GetJetCoordinator() core.JetCoordinator {
	return nil
}

func (l *MockLedger) GetPulseManager() core.PulseManager {
	return l.PM
}

func (l *MockLedger) HandleEvent(core.Event) (core.Reaction, error) {
	return nil, nil
}

// GetDefaultCtx creates default context for the host handler.
func GetDefaultCtx(hostHandler hosthandler.HostHandler) hosthandler.Context {
	ctx, _ := hostnetwork.NewContextBuilder(hostHandler).SetDefaultHost().Build()
	return ctx
}

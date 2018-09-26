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

package hostnetwork

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
)

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

// GetDefaultCtx creates default context for the host handler.
func GetDefaultCtx(hostHandler hosthandler.HostHandler) hosthandler.Context {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		panic(err)
	}
	return ctx
}

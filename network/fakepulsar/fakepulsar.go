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

package fakepulsar

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

// onPulse is a callbaback for pulse recv.
type callbackOnPulse func(pulse core.Pulse)

// FakePulsar is a struct which uses at void network state.
type FakePulsar struct {
	onPulse   callbackOnPulse
	stop      chan bool
	timeoutMs int32 // ms
}

// NewFakePulsar creates and returns a new FakePulsar.
func NewFakePulsar(callback callbackOnPulse, timeoutMs int32) *FakePulsar {
	return &FakePulsar{
		onPulse:   callback,
		timeoutMs: timeoutMs,
		stop:      make(chan bool),
	}
}

// GetFakePulse creates and returns a fake pulse.
func (fp *FakePulsar) GetFakePulse() *core.Pulse {
	return fp.newPulse()
}

// Start starts sending a fake pulse.
func (fp *FakePulsar) Start(ctx context.Context) {
	go func(fp *FakePulsar) {
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(fp.timeoutMs)):
				{
					fp.onPulse(*fp.GetFakePulse())
				}
			case <-fp.stop:
				return
			}
		}

	}(fp)
}

// Stop sending a fake pulse.
func (fp *FakePulsar) Stop(ctx context.Context) {
	fp.stop <- true
}

func (fp *FakePulsar) newPulse() *core.Pulse {
	generator := entropygenerator.StandardEntropyGenerator{}
	return &core.Pulse{
		PulseNumber:     0,
		NextPulseNumber: 0,
		Entropy:         generator.GenerateEntropy(),
	}
}

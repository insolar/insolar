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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

const delta uint32 = 0
const prevPulseNumber = 0

// onPulse is a callbaback for pulse recv.
type callbackOnPulse func(pulse core.Pulse)

// FakePulsar is a struct which uses at void network state.
type FakePulsar struct {
	onPulse callbackOnPulse
	stop    chan bool
	timeout int32 // ms
}

// NewFakePulsar creates and returns a new FakePulsar.
func NewFakePulsar(callback callbackOnPulse, timeout int32) *FakePulsar {
	return &FakePulsar{
		onPulse: callback,
		timeout: timeout,
	}
}

// GetFakePulse creates and returns a fake pulse.
func (fp *FakePulsar) GetFakePulse() *core.Pulse {
	pulse := pulsar.NewPulse(delta, prevPulseNumber, &entropygenerator.StandardEntropyGenerator{})
	pulse.PulseNumber = 0
	pulse.NextPulseNumber = 0
	return pulse
}

// Start starts sending a fake pulse.
func (fp *FakePulsar) Start() {
	go func(stop chan bool) {
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(fp.timeout)):
				{
					fp.onPulse(*fp.GetFakePulse())
				}
			case <-stop:
				return
			}
		}

	}(fp.stop)
}

// Stop sending a fake pulse.
func (fp *FakePulsar) Stop() {
	fp.stop <- true
}

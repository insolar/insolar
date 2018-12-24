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
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

// onPulse is a callbaback for pulse recv.
// type callbackOnPulse func(ctx context.Context, pulse core.Pulse)

// FakePulsar is a struct which uses at void network state.
type FakePulsar struct {
	onPulse   network.PulseHandler
	stop      chan bool
	timeoutMs int32 // ms
	pulse     int

	mutex   sync.RWMutex
	running bool
}

// NewFakePulsar creates and returns a new FakePulsar.
func NewFakePulsar(callback network.PulseHandler, timeoutMs int32) *FakePulsar {
	return &FakePulsar{
		onPulse:   callback,
		timeoutMs: timeoutMs,
		stop:      make(chan bool, 1),
		running:   false,
	}
}

// Start starts sending a fake pulse.
func (fp *FakePulsar) Start(ctx context.Context) {
	fp.mutex.Lock()
	defer fp.mutex.Unlock()

	fp.running = true
	go func(fp *FakePulsar) {
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(fp.timeoutMs)):
				fp.onPulse.HandlePulse(ctx, *fp.newPulse())
			case <-fp.stop:
				return
			}
		}

	}(fp)
	log.Info("fake pulsar started")
}

// Stop sending a fake pulse.
func (fp *FakePulsar) Stop(ctx context.Context) {
	fp.mutex.Lock()
	defer fp.mutex.Unlock()

	log.Info("Fake pulsar going to stop")

	if fp.running {
		fp.stop <- true
		close(fp.stop)
		fp.running = false
	}
	log.Info("Fake pulsar stopped")
}

func (fp *FakePulsar) Stopped() bool {
	fp.mutex.RLock()
	defer fp.mutex.RUnlock()

	return !fp.running
}

func (fp *FakePulsar) newPulse() *core.Pulse {
	fp.pulse++
	// TODO: fair entropy
	// generator := entropygenerator.StandardEntropyGenerator{}
	return &core.Pulse{
		EpochPulseNumber: -1,
		PulseNumber:      core.PulseNumber(fp.pulse),
		NextPulseNumber:  core.PulseNumber(fp.pulse + 1),
		Entropy:          core.Entropy{},
	}
}

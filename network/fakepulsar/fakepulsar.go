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
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

// onPulse is a callbaback for pulse recv.
// type callbackOnPulse func(ctx context.Context, pulse core.Pulse)

// FakePulsar is a struct which uses at void network state.
type FakePulsar struct {
	onPulse network.PulseHandler
	stop    chan bool
	mutex   sync.RWMutex
	running bool

	firstPulseTime     time.Time
	pulseDuration      time.Duration
	pulseNumberDelta   core.PulseNumber
	currentPulseNumber core.PulseNumber
}

// NewFakePulsar creates and returns a new FakePulsar.
func NewFakePulsar(callback network.PulseHandler, pulseDuration time.Duration) *FakePulsar {
	return &FakePulsar{
		onPulse: callback,
		stop:    make(chan bool, 1),
		running: false,

		pulseDuration:    pulseDuration,
		pulseNumberDelta: core.PulseNumber(pulseDuration.Seconds()),
	}
}

// Start starts sending a fake pulse.
func (fp *FakePulsar) Start(ctx context.Context, firstPulseTime time.Time) {
	fp.mutex.Lock()
	defer fp.mutex.Unlock()

	logger := inslogger.FromContext(ctx)
	startTime := time.Now()

	fp.firstPulseTime = firstPulseTime

	currentPulseNumber, initialWaitDuration := fp.getCurrentPulseAndWaitTime(startTime)

	fp.running = true
	fp.currentPulseNumber = currentPulseNumber

	logger.Infof(
		"Fake pulsar is going to start, currentPulse: %d, next pulse scheduled for: %s",
		currentPulseNumber,
		startTime.Add(initialWaitDuration),
	)

	time.AfterFunc(initialWaitDuration, func() {
		fp.pulse(ctx)

		for {
			now := time.Now()
			_, pulseDuration := fp.getCurrentPulseAndWaitTime(now)

			logger.Debug("Pulse scheduled for: %s", now.Add(fp.pulseDuration))

			select {
			case <-time.After(pulseDuration):
				fp.pulse(ctx)
			case <-fp.stop:
				return
			}
		}
	})
}

func (fp *FakePulsar) getCurrentPulseAndWaitTime(target time.Time) (core.PulseNumber, time.Duration) {
	return getCurrentPulseAndWaitTime(target, fp.firstPulseTime, fp.pulseDuration)
}

func (fp *FakePulsar) pulse(ctx context.Context) {
	fp.currentPulseNumber += fp.pulseNumberDelta
	go fp.onPulse.HandlePulse(ctx, *fp.newPulse())
}

// Stop sending a fake pulse.
func (fp *FakePulsar) Stop(ctx context.Context) {
	fp.mutex.Lock()
	defer fp.mutex.Unlock()

	inslogger.FromContext(ctx).Info("Fake pulsar going to stop")

	if fp.running {
		fp.stop <- true
		close(fp.stop)
		fp.running = false
	}

	inslogger.FromContext(ctx).Info("Fake pulsar stopped")
}

func (fp *FakePulsar) Stopped() bool {
	fp.mutex.RLock()
	defer fp.mutex.RUnlock()

	return !fp.running
}

func (fp *FakePulsar) newPulse() *core.Pulse {
	return &core.Pulse{
		EpochPulseNumber: -1,
		PulseNumber:      core.PulseNumber(fp.currentPulseNumber),
		NextPulseNumber:  core.PulseNumber(fp.currentPulseNumber + fp.pulseNumberDelta),
		Entropy:          core.Entropy{},
	}
}

func getCurrentPulseAndWaitTime(target, firstPulseTime time.Time, pulseDuration time.Duration) (core.PulseNumber, time.Duration) {
	if target.Before(firstPulseTime) {
		panic(fmt.Sprintf("First pulse time `%s` is greater then target `%s`", firstPulseTime, target))
	}

	timeSinceFirstPulse := target.Sub(firstPulseTime)

	passedPulses := int64(timeSinceFirstPulse) / int64(pulseDuration)
	currentPulse := core.PulseNumber(passedPulses)

	passedPulsesDuration := time.Duration(int64(pulseDuration) * passedPulses)
	waitDuration := pulseDuration - (timeSinceFirstPulse - passedPulsesDuration)

	return currentPulse, waitDuration
}

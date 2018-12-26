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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
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
		stop:    make(chan bool),
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

	fp.running = true
	fp.firstPulseTime = firstPulseTime

	pulseInfo := fp.getPulseInfo()

	fp.currentPulseNumber = pulseInfo.currentPulseNumber

	logger.Infof(
		"Fake pulsar is going to start, currentPulse: %d, next pulse scheduled for: %s",
		pulseInfo.currentPulseNumber,
		time.Now().Add(pulseInfo.nextPulseAfter),
	)

	time.AfterFunc(pulseInfo.nextPulseAfter, func() {
		fp.pulse(ctx)
		for {
			pulseInfo := fp.getPulseInfo()

			logger.Debug("Pulse scheduled for: %s", time.Now().Add(pulseInfo.nextPulseAfter))

			select {
			case <-time.After(pulseInfo.nextPulseAfter):
				fp.pulse(ctx)
			case <-fp.stop:
				return
			}
		}
	})
}

func (fp *FakePulsar) getPulseInfo() pulseInfo {
	return calculatePulseInfo(time.Now(), fp.firstPulseTime, fp.pulseDuration)
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

func (fp *FakePulsar) newPulse() *core.Pulse {
	return &core.Pulse{
		EpochPulseNumber: -1,
		PulseNumber:      core.PulseNumber(fp.currentPulseNumber),
		NextPulseNumber:  core.PulseNumber(fp.currentPulseNumber + fp.pulseNumberDelta),
		Entropy:          core.Entropy{},
	}
}

type pulseInfo struct {
	currentPulseNumber core.PulseNumber
	nextPulseAfter     time.Duration
}

func calculatePulseInfo(targetTime, firstPulseTime time.Time, pulseDuration time.Duration) pulseInfo {
	if firstPulseTime.After(targetTime) {
		log.Warn("First pulse time `%s` is after then targetTime `%s`", firstPulseTime, targetTime)

		return pulseInfo{
			currentPulseNumber: core.PulseNumber(0),
			nextPulseAfter:     firstPulseTime.Sub(targetTime),
		}
	}

	timeSinceFirstPulse := targetTime.Sub(firstPulseTime)

	passedPulses := int64(timeSinceFirstPulse) / int64(pulseDuration)
	currentPulseNumber := core.PulseNumber(passedPulses)

	passedPulsesDuration := time.Duration(int64(pulseDuration) * passedPulses)
	nextPulseAfter := pulseDuration - (timeSinceFirstPulse - passedPulsesDuration)

	return pulseInfo{
		currentPulseNumber: currentPulseNumber,
		nextPulseAfter:     nextPulseAfter,
	}
}

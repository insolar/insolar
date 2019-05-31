//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package fakepulsar

import (
	"context"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

// onPulse is a callbaback for pulse recv.
// type callbackOnPulse func(ctx context.Context, pulse insolar.Pulse)

// FakePulsar is a struct which uses at void network state.
type FakePulsar struct {
	onPulse network.PulseHandler
	stop    chan bool
	mutex   sync.RWMutex
	running bool

	firstPulseTime     time.Time
	pulseDuration      time.Duration
	pulseNumberDelta   insolar.PulseNumber
	currentPulseNumber insolar.PulseNumber
}

// NewFakePulsar creates and returns a new FakePulsar.
func NewFakePulsar(callback network.PulseHandler, pulseDuration time.Duration) *FakePulsar {
	return &FakePulsar{
		onPulse: callback,
		stop:    make(chan bool),
		running: false,

		pulseDuration:    pulseDuration,
		pulseNumberDelta: insolar.PulseNumber(pulseDuration.Seconds()),
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

			logger.Debug("Pulse scheduled for: ", time.Now().Add(pulseInfo.nextPulseAfter))

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
	fp.mutex.Lock()
	defer fp.mutex.Unlock()

	if !fp.running {
		return
	}

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

func (fp *FakePulsar) newPulse() *insolar.Pulse {
	return &insolar.Pulse{
		PulseTimestamp:   time.Now().UnixNano(),
		PrevPulseNumber:  fp.currentPulseNumber - fp.pulseNumberDelta,
		PulseNumber:      fp.currentPulseNumber,
		NextPulseNumber:  fp.currentPulseNumber + fp.pulseNumberDelta,
		EpochPulseNumber: -1,
		Entropy:          insolar.Entropy{},
	}
}

type pulseInfo struct {
	currentPulseNumber insolar.PulseNumber
	nextPulseAfter     time.Duration
}

func calculatePulseInfo(targetTime, firstPulseTime time.Time, pulseDuration time.Duration) pulseInfo {
	if firstPulseTime.After(targetTime) {
		log.Warn("First pulse time `%s` is after then targetTime `%s`", firstPulseTime, targetTime)

		return pulseInfo{
			currentPulseNumber: insolar.PulseNumber(0),
			nextPulseAfter:     firstPulseTime.Sub(targetTime),
		}
	}

	timeSinceFirstPulse := targetTime.Sub(firstPulseTime)

	passedPulses := int64(timeSinceFirstPulse) / int64(pulseDuration)
	currentPulseNumber := insolar.PulseNumber(passedPulses)

	passedPulsesDuration := time.Duration(int64(pulseDuration) * passedPulses)
	nextPulseAfter := pulseDuration - (timeSinceFirstPulse - passedPulsesDuration)

	return pulseInfo{
		currentPulseNumber: currentPulseNumber,
		nextPulseAfter:     nextPulseAfter,
	}
}

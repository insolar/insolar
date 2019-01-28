/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package fakepulsar

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
)

// TODO: write adequate tests instead of this imitation of work
func TestGetFakePulse(t *testing.T) {
	handler := network.PulseHandlerMock{}
	handler.HandlePulseFunc = func(p context.Context, p1 core.Pulse) {}
	pulsar := NewFakePulsar(&handler, time.Second)
	pulse := pulsar.newPulse()
	assert.NotNil(t, pulse)
	pulsar.currentPulseNumber++
	pulse2 := pulsar.newPulse()
	assert.NotNil(t, pulse2)
	assert.NotEqual(t, pulse, pulse2)

	pulsar2 := NewFakePulsar(&handler, time.Second)
	pulsar2.currentPulseNumber = pulsar.currentPulseNumber
	pulse3 := pulsar2.newPulse()
	assert.NotNil(t, pulse3)
	assert.Equal(t, pulse3, pulse2)
}

func TestFakePulsar_Start(t *testing.T) {
	handler := network.PulseHandlerMock{}
	handler.HandlePulseFunc = func(p context.Context, p1 core.Pulse) {}
	pulsar := NewFakePulsar(&handler, time.Second)

	ctx := context.TODO()
	firstPulseTime := time.Now()

	pulsar.Start(ctx, firstPulseTime)
	workTime := time.Millisecond * 3500
	time.Sleep(workTime)

	pulsar.Stop(ctx)

	pulseInfo := calculatePulseInfo(firstPulseTime.Add(workTime), firstPulseTime, 1000*time.Millisecond)

	assert.Equal(t, core.PulseNumber(3), pulsar.currentPulseNumber)
	assert.Equal(t, pulsar.currentPulseNumber, pulseInfo.currentPulseNumber)
}

func TestCalculatePulseInfo(t *testing.T) {
	firstPulseTime := time.Now()
	timePassed := 3500 * time.Millisecond
	pulsarNow := firstPulseTime.Add(timePassed)

	pulseInfo := calculatePulseInfo(pulsarNow, firstPulseTime, time.Second)

	assert.Equal(t, core.PulseNumber(3), pulseInfo.currentPulseNumber)
	assert.Equal(t, pulseInfo.nextPulseAfter, 500*time.Millisecond)
}

func TestCalculatePulseInfo_FirstPulseInFuture(t *testing.T) {
	pulsarNow := time.Now()
	timePassed := 3500 * time.Millisecond
	firstPulseTime := pulsarNow.Add(timePassed)

	pulseInfo := calculatePulseInfo(pulsarNow, firstPulseTime, time.Second)

	assert.Equal(t, core.PulseNumber(0), pulseInfo.currentPulseNumber)
	assert.Equal(t, 3500*time.Millisecond, pulseInfo.nextPulseAfter)
}

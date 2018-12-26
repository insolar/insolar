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

func TestGetPassedPulseCountAndWaitTime(t *testing.T) {
	firstPulseTime := time.Now()
	timePassed := 3500 * time.Millisecond
	pulsarNow := firstPulseTime.Add(timePassed)

	pulseInfo := calculatePulseInfo(pulsarNow, firstPulseTime, time.Second)

	assert.Equal(t, core.PulseNumber(3), pulseInfo.currentPulseNumber)
	assert.Equal(t, pulseInfo.nextPulseAfter, 500*time.Millisecond)
}

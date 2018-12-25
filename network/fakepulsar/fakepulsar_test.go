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
	pulsar := NewFakePulsar(&handler, 1000)
	pulse := pulsar.newPulse()
	assert.NotNil(t, pulse)
	pulsar.pulseNum++
	pulse2 := pulsar.newPulse()
	assert.NotNil(t, pulse2)
	assert.NotEqual(t, pulse, pulse2)

	pulsar2 := NewFakePulsar(&handler, 1000)
	pulsar2.pulseNum = pulsar.pulseNum
	pulse3 := pulsar2.newPulse()
	assert.NotNil(t, pulse3)
	assert.Equal(t, pulse3, pulse2)
}

func TestFakePulsar_Start(t *testing.T) {
	handler := network.PulseHandlerMock{}
	handler.HandlePulseFunc = func(p context.Context, p1 core.Pulse) {}
	pulsar := NewFakePulsar(&handler, 1000)
	ctx := context.TODO()
	pulsar.Start(ctx)
	time.Sleep(time.Millisecond * 1100)
	pulsar.Stop(ctx)
}

func TestGetPassedPulseCountAndWaitTime(t *testing.T) {
	pulseCount := 5
	waitSec := int64(5)
	pulseTimeout := 12000

	firstPulseTime := time.Date(2018, 12, 25, 2, 10, 10, 0, time.Local)
	local := time.Date(2018, 12, 25, 2, 11, 15, 0, time.Local)

	count, waitTime := GetPassedPulseCountAndWaitTime(local.Unix(), firstPulseTime.Unix(), int32(pulseTimeout))

	assert.Equal(t, int64(pulseCount), count)
	assert.Equal(t, waitSec, waitTime)
}

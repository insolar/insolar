/*
 *    Copyright 2019 Insolar Technologies
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

package conveyor

import (
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

func TestConveyor_ChangePulse(t *testing.T) {
	conveyor, err := NewPulseConveyor()
	require.NoError(t, err)
	callback := mockCallback()
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	err = conveyor.PreparePulse(pulse, callback)
	require.NoError(t, err)

	callback.(*mockSyncDone).GetResult()

	err = conveyor.ActivatePulse()
	require.NoError(t, err)
}

func TestConveyor_ChangePulseMultipleTimes(t *testing.T) {
	conveyor, err := NewPulseConveyor()
	require.NoError(t, err)

	pulseNumber := testRealPulse + testPulseDelta
	for i := 0; i < 20; i++ {
		callback := mockCallback()
		pulseNumber += testPulseDelta
		pulse := insolar.Pulse{PulseNumber: pulseNumber, NextPulseNumber: pulseNumber + testPulseDelta}
		err = conveyor.PreparePulse(pulse, callback)
		require.NoError(t, err)

		callback.(*mockSyncDone).GetResult()

		err = conveyor.ActivatePulse()
		require.NoError(t, err)
	}
}

func TestConveyor_ChangePulseMultipleTimes_WithEvents(t *testing.T) {
	conveyor, err := NewPulseConveyor()
	require.NoError(t, err)

	pulseNumber := testRealPulse + testPulseDelta
	for i := 0; i < 100; i++ {

		go func() {
			for j := 0; j < 1; j++ {
				err = conveyor.SinkPush(pulseNumber, "TEST")
				require.NoError(t, err)

				err = conveyor.SinkPush(pulseNumber-testPulseDelta, "TEST")
				require.NoError(t, err)

				err = conveyor.SinkPush(pulseNumber+testPulseDelta, "TEST")
				require.NoError(t, err)

				err = conveyor.SinkPushAll(pulseNumber, []interface{}{"TEST", i * j})
				require.NoError(t, err)
			}
		}()

		go func() {
			for j := 0; j < 100; j++ {
				conveyor.GetState()
			}
		}()

		go func() {
			for j := 0; j < 100; j++ {
				conveyor.IsOperational()
			}
		}()

		callback := mockCallback()
		pulseNumber += testPulseDelta
		pulse := insolar.Pulse{PulseNumber: pulseNumber, NextPulseNumber: pulseNumber + testPulseDelta}
		err = conveyor.PreparePulse(pulse, callback)
		require.NoError(t, err)

		if i == 0 {
			require.Equal(t, 0, callback.(*mockSyncDone).GetResult())
		} else {
			require.Equal(t, 555, callback.(*mockSyncDone).GetResult())
		}

		err = conveyor.ActivatePulse()
		require.NoError(t, err)

		go func() {
			for j := 0; j < 10; j++ {
				require.NoError(t, conveyor.SinkPushAll(pulseNumber, []interface{}{"TEST", i}))
				require.NoError(t, conveyor.SinkPush(pulseNumber, "TEST"))
				require.NoError(t, conveyor.SinkPush(pulseNumber-testPulseDelta, "TEST"))
				err = conveyor.SinkPush(pulseNumber+testPulseDelta, "TEST")
				require.NoError(t, err)
			}
		}()
	}

	time.Sleep(time.Millisecond * 200)
}

// TODO: Add test on InitiateShutdown

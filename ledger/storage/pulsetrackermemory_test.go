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

package storage

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulseTrackerMemory_GetPulse(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[core.PulseNumber]*Pulse{},
	}
	existingPulse := &Pulse{Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber}}
	pulseTracker.memory[core.FirstPulseNumber] = existingPulse

	// Act
	pulse, err := pulseTracker.GetPulse(ctx, core.FirstPulseNumber)
	_, notFoundErr := pulseTracker.GetPulse(ctx, core.FirstPulseNumber+1)

	// Assert
	require.NoError(t, err)
	require.Equal(t, existingPulse, pulse)
	require.Equal(t, ErrPulseNotFound, notFoundErr)
}

func TestPulseTrackerMemory_GetPreviousPulse(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[core.PulseNumber]*Pulse{},
	}
	firstPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber},
	}
	secondPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 1},
		Prev:  &firstPulse.Pulse.PulseNumber,
	}
	badPrev := core.PulseNumber(42)
	thirdPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 2},
		Prev:  &badPrev,
	}
	pulseTracker.memory[core.FirstPulseNumber] = firstPulse
	pulseTracker.memory[core.FirstPulseNumber+1] = secondPulse
	pulseTracker.memory[core.FirstPulseNumber+2] = thirdPulse

	// Act
	pulse, err := pulseTracker.GetPreviousPulse(ctx, core.FirstPulseNumber+1)
	_, prevPulseErr := pulseTracker.GetPreviousPulse(ctx, core.FirstPulseNumber)
	_, badPrevErr := pulseTracker.GetPreviousPulse(ctx, core.FirstPulseNumber+2)
	_, notFoundErr := pulseTracker.GetPreviousPulse(ctx, 42)

	// Assert
	require.NoError(t, err)
	require.Equal(t, firstPulse, pulse)
	require.Equal(t, ErrPrevPulseNotFound, prevPulseErr)
	require.Equal(t, ErrPulseNotFound, badPrevErr)
	require.Equal(t, ErrPulseNotFound, notFoundErr)
}

func TestPulseTrackerMemory_GetNthPrevPulse(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[core.PulseNumber]*Pulse{},
	}
	firstPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber},
	}
	secondPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 1},
		Prev:  &firstPulse.Pulse.PulseNumber,
	}
	thirdPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 2},
		Prev:  &secondPulse.Pulse.PulseNumber,
	}
	forth := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 3},
		Prev:  &thirdPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[core.FirstPulseNumber] = firstPulse
	pulseTracker.memory[core.FirstPulseNumber+1] = secondPulse
	pulseTracker.memory[core.FirstPulseNumber+2] = thirdPulse
	pulseTracker.memory[core.FirstPulseNumber+3] = forth

	// Act and Assert
	targetPulse, err := pulseTracker.GetNthPrevPulse(ctx, 0, core.FirstPulseNumber+3)
	require.NoError(t, err)
	require.Equal(t, forth, targetPulse)

	prev1pulse, err := pulseTracker.GetNthPrevPulse(ctx, 1, core.FirstPulseNumber+3)
	require.NoError(t, err)
	require.Equal(t, thirdPulse, prev1pulse)

	prev2pulse, err := pulseTracker.GetNthPrevPulse(ctx, 2, core.FirstPulseNumber+3)
	require.NoError(t, err)
	require.Equal(t, secondPulse, prev2pulse)

	prev3pulse, err := pulseTracker.GetNthPrevPulse(ctx, 3, core.FirstPulseNumber+3)
	require.NoError(t, err)
	require.Equal(t, firstPulse, prev3pulse)

	_, err = pulseTracker.GetNthPrevPulse(ctx, 4, core.FirstPulseNumber+3)
	require.Equal(t, ErrPrevPulseNotFound, err)
}

func TestPulseTrackerMemory_GetLatestPulse(t *testing.T) {
	assert.True(t, false)
}

func TestPulseTrackerMemory_AddPulse(t *testing.T) {
	assert.True(t, false)
}

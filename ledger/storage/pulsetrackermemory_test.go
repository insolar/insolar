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
	assert.Equal(t, existingPulse, pulse)
	assert.Equal(t, ErrPulseNotFound, notFoundErr)
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
	assert.Equal(t, firstPulse, pulse)
	assert.Equal(t, ErrPrevPulseNotFound, prevPulseErr)
	assert.Equal(t, ErrPulseNotFound, badPrevErr)
	assert.Equal(t, ErrPulseNotFound, notFoundErr)
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
	fourth := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 3},
		Prev:  &thirdPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[core.FirstPulseNumber] = firstPulse
	pulseTracker.memory[core.FirstPulseNumber+1] = secondPulse
	pulseTracker.memory[core.FirstPulseNumber+2] = thirdPulse
	pulseTracker.memory[core.FirstPulseNumber+3] = fourth

	// Act and Assert
	targetPulse, err := pulseTracker.GetNthPrevPulse(ctx, 0, core.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, fourth, targetPulse)

	prev1pulse, err := pulseTracker.GetNthPrevPulse(ctx, 1, core.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, thirdPulse, prev1pulse)

	prev2pulse, err := pulseTracker.GetNthPrevPulse(ctx, 2, core.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, secondPulse, prev2pulse)

	prev3pulse, err := pulseTracker.GetNthPrevPulse(ctx, 3, core.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, firstPulse, prev3pulse)

	_, err = pulseTracker.GetNthPrevPulse(ctx, 4, core.FirstPulseNumber+3)
	assert.Equal(t, ErrPrevPulseNotFound, err)
}

func TestPulseTrackerMemory_GetLatestPulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[core.PulseNumber]*Pulse{},
	}

	// Check empty storage
	_, err := pulseTracker.GetLatestPulse(ctx)
	assert.Equal(t, ErrEmptyLatestPulse, err)

	// Check correct pulseNumber, but empty storage
	pulseTracker.latestPulse = 1
	_, err = pulseTracker.GetLatestPulse(ctx)
	assert.Equal(t, ErrPulseNotFound, err)

	// Add and check first pulse
	// latest = first
	firstPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber},
	}
	pulseTracker.memory[core.FirstPulseNumber] = firstPulse
	pulseTracker.latestPulse = core.FirstPulseNumber
	pulse, err := pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, firstPulse, pulse)

	// Add and check second pulse
	// latest = second
	secondPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 1},
		Prev:  &firstPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[core.FirstPulseNumber+1] = secondPulse
	pulseTracker.latestPulse = core.FirstPulseNumber + 1
	pulse, err = pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, secondPulse, pulse)

	// Add and check third pulse
	// latest != third, latest = second because third pulseNumber smaller than second
	thirdPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: 42},
		Prev:  &secondPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[42] = thirdPulse
	pulse, err = pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, secondPulse, pulse)

	// Add and check fourth pulse
	// latest = fourth
	fourthPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 3},
		Prev:  &thirdPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[core.FirstPulseNumber+3] = fourthPulse
	pulseTracker.latestPulse = core.FirstPulseNumber + 3
	pulse, err = pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, fourthPulse, pulse)
}

func TestPulseTrackerMemory_AddPulse_FailFirstCheck(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[core.PulseNumber]*Pulse{},
	}

	// Check pulse smaller than current
	firstPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber},
	}
	pulseTracker.latestPulse = core.FirstPulseNumber + 1
	err := pulseTracker.AddPulse(ctx, firstPulse.Pulse)

	assert.Equal(t, ErrLesserPulse, err)

	// Check pulse equal with current
	pulseTracker.latestPulse = core.FirstPulseNumber
	err = pulseTracker.AddPulse(ctx, firstPulse.Pulse)

	assert.Equal(t, ErrOverride, err)
}

func TestPulseTrackerMemory_AddPulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	// Arrange
	// Check new pulse adding
	pulseTracker := &pulseTrackerMemory{
		memory: map[core.PulseNumber]*Pulse{},
	}
	firstPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber},
	}
	prevPN := core.PulseNumber(0)
	firstPulse.Prev = &prevPN
	firstPulse.SerialNumber = 1

	// Act
	err := pulseTracker.AddPulse(ctx, firstPulse.Pulse)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, firstPulse.Pulse, pulseTracker.memory[core.FirstPulseNumber].Pulse)
	assert.Equal(t, firstPulse.Pulse.PulseNumber, pulseTracker.latestPulse)
	assert.Equal(t, firstPulse, pulseTracker.memory[core.FirstPulseNumber])

	// Arrange
	// Check pulse adding to non-empty storage
	secondPulse := &Pulse{
		Pulse: core.Pulse{PulseNumber: core.FirstPulseNumber + 1},
	}
	secondPulse.Prev = &firstPulse.Pulse.PulseNumber
	secondPulse.SerialNumber = firstPulse.SerialNumber + 1
	firstPulse.Next = &secondPulse.Pulse.PulseNumber

	// Act
	err = pulseTracker.AddPulse(ctx, secondPulse.Pulse)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, secondPulse.Pulse, pulseTracker.memory[core.FirstPulseNumber+1].Pulse)
	assert.Equal(t, secondPulse.Pulse.PulseNumber, pulseTracker.latestPulse)
	assert.Equal(t, *secondPulse.Prev, pulseTracker.memory[core.FirstPulseNumber].Pulse.PulseNumber)
	assert.Equal(t, secondPulse.SerialNumber, pulseTracker.memory[core.FirstPulseNumber+1].SerialNumber)
	assert.Equal(t, *firstPulse.Next, pulseTracker.memory[core.FirstPulseNumber+1].Pulse.PulseNumber)
}

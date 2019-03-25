//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package storage

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulseTrackerMemory_GetPulse(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[insolar.PulseNumber]Pulse{},
	}
	existingPulse := Pulse{Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber}}
	existingPulse.SerialNumber = 1
	pulseTracker.memory[insolar.FirstPulseNumber] = existingPulse

	// Act
	pulse, err := pulseTracker.GetPulse(ctx, insolar.FirstPulseNumber)
	_, notFoundErr := pulseTracker.GetPulse(ctx, insolar.FirstPulseNumber+1)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, existingPulse, *pulse)
	assert.Equal(t, insolar.ErrNotFound, notFoundErr)
}

func TestPulseTrackerMemory_GetPreviousPulse(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[insolar.PulseNumber]Pulse{},
	}
	firstPulse := Pulse{
		Pulse:        insolar.Pulse{PulseNumber: insolar.FirstPulseNumber},
		SerialNumber: 1,
	}
	secondPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1},
		Prev:  &firstPulse.Pulse.PulseNumber,
	}
	badPrev := insolar.PulseNumber(42)
	thirdPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 2},
		Prev:  &badPrev,
	}
	pulseTracker.memory[insolar.FirstPulseNumber] = firstPulse
	pulseTracker.memory[insolar.FirstPulseNumber+1] = secondPulse
	pulseTracker.memory[insolar.FirstPulseNumber+2] = thirdPulse

	// Act
	pulse, err := pulseTracker.GetPreviousPulse(ctx, insolar.FirstPulseNumber+1)
	_, prevPulseErr := pulseTracker.GetPreviousPulse(ctx, insolar.FirstPulseNumber)
	_, badPrevErr := pulseTracker.GetPreviousPulse(ctx, insolar.FirstPulseNumber+2)
	_, notFoundErr := pulseTracker.GetPreviousPulse(ctx, 42)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, firstPulse, *pulse)
	assert.Equal(t, ErrPrevPulse, prevPulseErr)
	assert.Equal(t, insolar.ErrNotFound, badPrevErr)
	assert.Equal(t, insolar.ErrNotFound, notFoundErr)
}

func TestPulseTrackerMemory_GetNthPrevPulse(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[insolar.PulseNumber]Pulse{},
	}
	firstPulse := Pulse{
		Pulse:        insolar.Pulse{PulseNumber: insolar.FirstPulseNumber},
		SerialNumber: 1,
	}
	secondPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1},
		Prev:  &firstPulse.Pulse.PulseNumber,
	}
	thirdPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 2},
		Prev:  &secondPulse.Pulse.PulseNumber,
	}
	fourth := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 3},
		Prev:  &thirdPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[insolar.FirstPulseNumber] = firstPulse
	pulseTracker.memory[insolar.FirstPulseNumber+1] = secondPulse
	pulseTracker.memory[insolar.FirstPulseNumber+2] = thirdPulse
	pulseTracker.memory[insolar.FirstPulseNumber+3] = fourth

	// Act and Assert
	targetPulse, err := pulseTracker.GetNthPrevPulse(ctx, 0, insolar.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, fourth, *targetPulse)

	prev1pulse, err := pulseTracker.GetNthPrevPulse(ctx, 1, insolar.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, thirdPulse, *prev1pulse)

	prev2pulse, err := pulseTracker.GetNthPrevPulse(ctx, 2, insolar.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, secondPulse, *prev2pulse)

	prev3pulse, err := pulseTracker.GetNthPrevPulse(ctx, 3, insolar.FirstPulseNumber+3)
	require.NoError(t, err)
	assert.Equal(t, firstPulse, *prev3pulse)

	_, err = pulseTracker.GetNthPrevPulse(ctx, 4, insolar.FirstPulseNumber+3)
	assert.Equal(t, ErrPrevPulse, err)
}

func TestPulseTrackerMemory_GetLatestPulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[insolar.PulseNumber]Pulse{},
	}

	// Check empty storage
	_, err := pulseTracker.GetLatestPulse(ctx)
	assert.Equal(t, insolar.ErrNotFound, err)

	// Check correct pulseNumber, but empty storage
	pulseTracker.latestPulse = 1
	_, err = pulseTracker.GetLatestPulse(ctx)
	assert.Equal(t, insolar.ErrNotFound, err)

	// Add and check first pulse
	// latest = first
	firstPulse := Pulse{
		Pulse:        insolar.Pulse{PulseNumber: insolar.FirstPulseNumber},
		SerialNumber: 1,
	}
	pulseTracker.memory[insolar.FirstPulseNumber] = firstPulse
	pulseTracker.latestPulse = insolar.FirstPulseNumber
	pulse, err := pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, firstPulse, *pulse)

	// Add and check second pulse
	// latest = second
	secondPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1},
		Prev:  &firstPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[insolar.FirstPulseNumber+1] = secondPulse
	pulseTracker.latestPulse = insolar.FirstPulseNumber + 1
	pulse, err = pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, secondPulse, *pulse)

	// Add and check third pulse
	// latest != third, latest = second because third pulseNumber smaller than second
	thirdPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: 42},
		Prev:  &secondPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[42] = thirdPulse
	pulse, err = pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, secondPulse, *pulse)

	// Add and check fourth pulse
	// latest = fourth
	fourthPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 3},
		Prev:  &thirdPulse.Pulse.PulseNumber,
	}
	pulseTracker.memory[insolar.FirstPulseNumber+3] = fourthPulse
	pulseTracker.latestPulse = insolar.FirstPulseNumber + 3
	pulse, err = pulseTracker.GetLatestPulse(ctx)
	require.NoError(t, err)
	assert.Equal(t, fourthPulse, *pulse)
}

func TestPulseTrackerMemory_AddPulse_FailFirstCheck(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	pulseTracker := &pulseTrackerMemory{
		memory: map[insolar.PulseNumber]Pulse{},
	}

	// Check pulse smaller than current
	firstPulse := &Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber},
	}
	pulseTracker.latestPulse = insolar.FirstPulseNumber + 1
	err := pulseTracker.AddPulse(ctx, firstPulse.Pulse)

	assert.Equal(t, ErrBadPulse, err)

	// Check pulse equal with current
	pulseTracker.latestPulse = insolar.FirstPulseNumber
	err = pulseTracker.AddPulse(ctx, firstPulse.Pulse)

	assert.Equal(t, ErrBadPulse, err)
}

func TestPulseTrackerMemory_AddPulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	// Arrange
	// Check new pulse adding
	pulseTracker := &pulseTrackerMemory{
		memory: map[insolar.PulseNumber]Pulse{},
	}
	firstPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber},
	}
	prevPN := insolar.PulseNumber(0)
	firstPulse.Prev = &prevPN
	firstPulse.SerialNumber = 1

	// Act
	err := pulseTracker.AddPulse(ctx, firstPulse.Pulse)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, firstPulse.Pulse, pulseTracker.memory[insolar.FirstPulseNumber].Pulse)
	assert.Equal(t, firstPulse.Pulse.PulseNumber, pulseTracker.latestPulse)
	assert.Equal(t, firstPulse, pulseTracker.memory[insolar.FirstPulseNumber])

	// Arrange
	// Check pulse adding to non-empty storage
	secondPulse := &Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1},
	}
	secondPulse.Prev = &firstPulse.Pulse.PulseNumber
	secondPulse.SerialNumber = firstPulse.SerialNumber + 1
	firstPulse.Next = &secondPulse.Pulse.PulseNumber

	// Act
	err = pulseTracker.AddPulse(ctx, secondPulse.Pulse)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, secondPulse.Pulse, pulseTracker.memory[insolar.FirstPulseNumber+1].Pulse)
	assert.Equal(t, secondPulse.Pulse.PulseNumber, pulseTracker.latestPulse)
	assert.Equal(t, *secondPulse.Prev, pulseTracker.memory[insolar.FirstPulseNumber].Pulse.PulseNumber)
	assert.Equal(t, secondPulse.SerialNumber, pulseTracker.memory[insolar.FirstPulseNumber+1].SerialNumber)
	assert.Equal(t, *firstPulse.Next, pulseTracker.memory[insolar.FirstPulseNumber+1].Pulse.PulseNumber)

	// Check pulse from the past for non-empty storage
	pastPulse := &Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber - 1},
	}
	err = pulseTracker.AddPulse(ctx, pastPulse.Pulse)
	require.Equal(t, ErrBadPulse, err)
}

func TestPulseTrackerMemory_DeletePulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	pulseTracker := &pulseTrackerMemory{
		memory: map[insolar.PulseNumber]Pulse{},
	}
	firstPulse := Pulse{
		Pulse: insolar.Pulse{PulseNumber: insolar.FirstPulseNumber},
	}
	assert.Equal(t, 0, len(pulseTracker.memory))

	// Check deleting from empty storage
	err := pulseTracker.DeletePulse(ctx, insolar.FirstPulseNumber)
	require.NoError(t, err)
	assert.Equal(t, 0, len(pulseTracker.memory))

	// Add pulse to storage
	pulseTracker.memory[insolar.FirstPulseNumber] = firstPulse
	assert.Equal(t, 1, len(pulseTracker.memory))
	assert.Equal(t, firstPulse.Pulse, pulseTracker.memory[insolar.FirstPulseNumber].Pulse)

	// Check deleting from non-empty storage
	err = pulseTracker.DeletePulse(ctx, insolar.FirstPulseNumber)
	require.NoError(t, err)
	assert.Equal(t, 0, len(pulseTracker.memory))
}

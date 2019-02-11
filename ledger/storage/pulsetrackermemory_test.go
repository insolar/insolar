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

	assert.Equal(t, notFoundErr, ErrNotFound)
}

func TestPulseTrackerMemory_GetPreviousPulse(t *testing.T) {
	assert.True(t, false)
}

func TestPulseTrackerMemory_GetNthPrevPulse(t *testing.T) {
	assert.True(t, false)
}

func TestPulseTrackerMemory_GetLatestPulse(t *testing.T) {
	assert.True(t, false)
}

func TestPulseTrackerMemory_AddPulse(t *testing.T) {
	assert.True(t, false)
}

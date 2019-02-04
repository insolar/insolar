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

package artifactmanager

import (
	"context"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewNodeCalculatorConcrete(t *testing.T) {
	t.Parallel()
	// Act
	calc := NewNodeCalculatorConcrete(12)

	// Assert
	require.NotNil(t, calc)
	require.Equal(t, 12, calc.LightChainLimit)
}

func TestNodeCalculatorConcrete_IsBeyondLimit_ProblemsWithTracker(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseMock.Return(nil, errors.New("it's expected"))
	calc := NewNodeCalculatorConcrete(12)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.NotNil(t, err)
	require.Equal(t, false, res)
}

func TestNodeCalculatorConcrete_IsBeyondLimit_ProblemsWithTracker_SecondCall(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{}, nil
		}

		return nil, errors.New("it's expected")
	}
	calc := NewNodeCalculatorConcrete(12)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.NotNil(t, err)
	require.Equal(t, false, res)
}

func TestNodeCalculatorConcrete_IsBeyondLimit_OutsideOfLightChainLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{SerialNumber: 50}, nil
		}

		return &storage.Pulse{SerialNumber: 24}, nil
	}
	calc := NewNodeCalculatorConcrete(25)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, true, res)
}

func TestNodeCalculatorConcrete_IsBeyondLimit_InsideOfLightChainLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{SerialNumber: 50}, nil
		}

		return &storage.Pulse{SerialNumber: 34}, nil
	}
	calc := NewNodeCalculatorConcrete(25)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, false, res)
}

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
// +build slowtest

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/testutils"
)

func TestVirtual_BasicOperations(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultVMConfig()

	t.Run("happy path", func(t *testing.T) {
		expectedRes := struct {
			blip string
		}{
			blip: "blop",
		}

		mle := testutils.NewMachineLogicExecutorMock(t).CallMethodMock.Set(
			func(_ context.Context, _ *insolar.LogicCallContext, _ insolar.Reference, _ []byte, _ string, _ insolar.Arguments) ([]byte, insolar.Arguments, error) {
				return insolar.MustSerialize(expectedRes), insolar.MustSerialize(expectedRes), nil
			},
		)

		mm := machinesmanager.NewMachinesManagerMock(t).GetExecutorMock.Set(
			func(_ insolar.MachineType) (insolar.MachineLogicExecutor, error) {
				return mle, nil
			},
		).RegisterExecutorMock.Set(
			func(_ insolar.MachineType, _ insolar.MachineLogicExecutor) error {
				return nil
			},
		)

		s, err := NewServer(t, ctx, cfg, nil, mm)
		require.NoError(t, err)
		defer s.Stop(ctx)

		// First pulse goes in storage then interrupts.
		s.IncrementPulse(ctx)

		codeID, err := s.mimic.AddCode(ctx, []byte{})
		require.NoError(t, err)

		prototypeID, err := s.mimic.AddObject(ctx, *codeID, true, []byte{})
		require.NoError(t, err)

		objectID, err := s.mimic.AddObject(ctx, *prototypeID, false, []byte{})
		require.NoError(t, err)

		objectRef := insolar.NewReference(*objectID)

		res, requestRef, err := CallContract(
			s, objectRef, "good.CallMethod", nil, s.pulse.PulseNumber,
		)
		require.NoError(t, err)

		assert.NotEmpty(t, requestRef)
		assert.Equal(t, &reply.CallMethod{
			Object: objectRef,
			Result: insolar.MustSerialize(expectedRes),
		}, res)
	})
}

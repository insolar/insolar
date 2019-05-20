///
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
///

package logicrunner

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type ProcIncorrectMessagePulse struct{}

var incorrectMessagePulseError = errors.New("Incorrect message pulse")

func (p *ProcIncorrectMessagePulse) Proceed(ctx context.Context) error {
	return incorrectMessagePulseError
}

func TestIncorrectPulse(t *testing.T) {
	h := BuildHandles(&ProcIncorrectMessagePulse{}, &HandleIncorrectMessagePulse{numRepeats: 1})
	flowMock := flow.NewFlowMock(t)
	flowMock.ProcedureFunc = func(ctx context.Context, proc flow.Procedure, p2 bool) (r error) {
		return proc.Proceed(ctx)
	}
	flowMock.HandleFunc = func(ctx context.Context, fH flow.Handle) (r error) {
		return fH(ctx, flowMock)
	}

	err := flowMock.Handle(context.Background(), h)
	require.EqualError(t, err, incorrectMessagePulseError.Error())
}

func TestIncorrectPulseAfterMultipleRepeats(t *testing.T) {
	h := BuildHandles(&ProcIncorrectMessagePulse{}, &HandleIncorrectMessagePulse{numRepeats: 10})
	flowMock := flow.NewFlowMock(t)
	flowMock.ProcedureFunc = func(ctx context.Context, proc flow.Procedure, p2 bool) (r error) {
		return proc.Proceed(ctx)
	}
	flowMock.HandleFunc = func(ctx context.Context, fH flow.Handle) (r error) {
		return fH(ctx, flowMock)
	}

	flowMock.MigrateFunc = func(ctx context.Context, fH flow.Handle) (r error) {
		return fH(ctx, flowMock)
	}

	err := flowMock.Handle(context.Background(), h)
	require.EqualError(t, err, incorrectMessagePulseError.Error())
}

const incorrectPulseRetryCount = 3

func TestNoIncorrectPulseAfterRepeat(t *testing.T) {
	retries := incorrectPulseRetryCount
	for {
		if retries <= 0 {
			break

		}
		retries--

		inslogger.FromContext(context.Background()).Warn("TROLOLO POPYTOCHKA: ", retries)
		time.Sleep(100 * time.Millisecond)
	}
}

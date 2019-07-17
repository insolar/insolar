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

package logicrunner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func InitBroker(_ *testing.T, ctx context.Context, count int, broker *ExecutionBroker, withMocks bool) {
	if withMocks {
		rem := broker.requestsExecutor.(*RequestsExecutorMock)
		rem.ExecuteAndSaveMock.Return(nil, nil)
		rem.SendReplyMock.Return()
	}

	for i := 0; i < count; i++ {
		reqRef := gen.Reference()
		broker.Put(ctx, false, &Transcript{
			LogicContext: &insolar.LogicCallContext{},
			Context:      ctx,
			RequestRef:   reqRef,
			Request:      &record.IncomingRequest{},
		})
	}
}

func newExecutionBroker(
	t *testing.T,
	ctx context.Context,
	count int,
	list *CurrentExecutionList,
	pending *insolar.PendingState,
) *ExecutionBroker {
	re := NewRequestsExecutorMock(t)
	mb := testutils.NewMessageBusMock(t)
	jc := jet.NewCoordinatorMock(t)
	ps := pulse.NewAccessorMock(t)
	pm := &publisherMock{}

	lr := LogicRunner{
		RequestsExecutor: NewRequestsExecutorMock(t),
		StateStorage:     NewStateStorage(pm, re, mb, jc, ps),
	}

	objectRef := gen.Reference()
	broker := lr.StateStorage.UpsertExecutionState(objectRef)

	InitBroker(t, ctx, count, broker, true)
	if list != nil {
		broker.currentList = list
	}
	if pending != nil {
		broker.executionState.pending = *pending
	}

	return broker
}

func TestExecutionState_OnPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	list := NewCurrentExecutionList()
	requestRef := gen.Reference()
	list.SetTranscript(&Transcript{RequestRef: requestRef})

	inPending := insolar.InPending

	table := []struct {
		name             string
		broker           *ExecutionBroker
		meNext           bool
		numberOfMessages int
		checkES          func(t *testing.T, es *ExecutionState, broker *ExecutionBroker)
	}{
		{
			name: "blank execution state",
		},
		{
			name:             "we have queue",
			broker:           newExecutionBroker(t, ctx, 1, nil, nil),
			numberOfMessages: 1,
			checkES: func(t *testing.T, es *ExecutionState, broker *ExecutionBroker) {
				require.Equal(t, 0, broker.mutable.Length())
			},
		},
		{
			name:             "we have queue, we are next",
			meNext:           true,
			broker:           newExecutionBroker(t, ctx, 1, nil, nil),
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState, broker *ExecutionBroker) {
				require.Equal(t, 1, broker.mutable.Length())
			},
		},
		{
			name:             "running something without queue, pending execution",
			broker:           newExecutionBroker(t, ctx, 0, list, nil),
			numberOfMessages: 2,
			checkES: func(t *testing.T, es *ExecutionState, broker *ExecutionBroker) {
				require.Equal(t, 0, broker.mutable.Length())
				require.Equal(t, insolar.InPending, es.pending)
			},
		},
		{
			name:             "running something without queue, we're next",
			broker:           newExecutionBroker(t, ctx, 0, list, nil),
			meNext:           true,
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState, broker *ExecutionBroker) {
				require.Equal(t, 0, broker.mutable.Length())
			},
		},
		{
			name:             "in not confirmed pending and no queue, still message",
			broker:           newExecutionBroker(t, ctx, 0, nil, &inPending),
			numberOfMessages: 1,
			checkES: func(t *testing.T, es *ExecutionState, broker *ExecutionBroker) {
				require.Equal(t, insolar.NotPending, es.pending)
			},
		},
		{
			name:             "in not confirmed pending and no queue, we're next",
			broker:           newExecutionBroker(t, ctx, 0, nil, &inPending),
			meNext:           true,
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState, broker *ExecutionBroker) {
				require.Equal(t, insolar.NotPending, es.pending)
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)

			messages := test.broker.OnPulse(ctx, test.meNext)
			require.Equal(t, test.numberOfMessages, len(messages))
			if test.checkES != nil {
				test.checkES(t, &test.broker.executionState, test.broker)
			}
		})
	}
}

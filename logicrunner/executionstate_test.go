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
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func InitBroker(_ *testing.T, ctx context.Context, count int, state *ExecutionState, withMocks bool) {
	if withMocks {
		rem := state.Broker.logicRunner.RequestsExecutor.(*RequestsExecutorMock)
		rem.ExecuteAndSaveMock.Return(nil, nil)
		rem.SendReplyMock.Return()
	}

	for i := 0; i < count; i++ {
		reqRef := gen.Reference()
		state.Broker.Put(ctx, false, &Transcript{
			LogicContext: &insolar.LogicCallContext{},
			Context:      ctx,
			RequestRef:   &reqRef,
			Request:      &record.IncomingRequest{},
		})
	}

}

func newExecutionStateLength(t *testing.T, ctx context.Context, count int, list *CurrentExecutionList,
	pending *message.PendingState) *ExecutionState {

	es := NewExecutionState(gen.Reference())
	es.Broker.logicRunner = &LogicRunner{}
	es.Broker.logicRunner.RequestsExecutor = NewRequestsExecutorMock(t)
	InitBroker(t, ctx, count, es, true)
	if list != nil {
		es.Broker.currentList = list
	}
	if pending != nil {
		es.pending = *pending
	}
	return es
}

func TestExecutionState_OnPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	list := NewCurrentExecutionList()
	requestRef := gen.Reference()
	list.SetTranscript(&Transcript{RequestRef: &requestRef})

	inPending := message.InPending

	table := []struct {
		name             string
		es               *ExecutionState
		meNext           bool
		numberOfMessages int
		checkES          func(t *testing.T, es *ExecutionState)
	}{
		{
			name: "blank execution state",
		},
		{
			name:             "we have queue",
			es:               newExecutionStateLength(t, ctx, 1, nil, nil),
			numberOfMessages: 1,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, es.Broker.mutable.Length(), 0)
			},
		},
		{
			name:             "we have queue, we are next",
			meNext:           true,
			es:               newExecutionStateLength(t, ctx, 1, nil, nil),
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, es.Broker.mutable.Length(), 1)
			},
		},
		{
			name:             "running something without queue, pending execution",
			es:               newExecutionStateLength(t, ctx, 0, list, nil),
			numberOfMessages: 2,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, es.Broker.mutable.Length(), 0)
				require.Equal(t, message.InPending, es.pending)
			},
		},
		{
			name:             "running something without queue, we're next",
			es:               newExecutionStateLength(t, ctx, 0, list, nil),
			meNext:           true,
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, es.Broker.mutable.Length(), 0)
			},
		},
		{
			name:             "in not confirmed pending and no queue, still message",
			es:               newExecutionStateLength(t, ctx, 0, nil, &inPending),
			numberOfMessages: 1,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, message.NotPending, es.pending)
			},
		},
		{
			name:             "in not confirmed pending and no queue, we're next",
			es:               newExecutionStateLength(t, ctx, 0, nil, &inPending),
			meNext:           true,
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, message.NotPending, es.pending)
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			messages := test.es.OnPulse(ctx, test.meNext)
			require.Equal(t, test.numberOfMessages, len(messages))
			if test.checkES != nil {
				test.checkES(t, test.es)
			}
		})
	}
}

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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestExecutionState_OnPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	list := NewCurrentExecutionList()
	list.Set(testutils.RandomRef(), &CurrentExecution{})

	table := []struct {
		name             string
		es               ExecutionState
		meNext           bool
		numberOfMessages int
		checkES          func(t *testing.T, es *ExecutionState)
	}{
		{
			name: "blank execution state",
		},
		{
			name:             "we have queue",
			es:               ExecutionState{Queue: []ExecutionQueueElement{{ctx: ctx}}},
			numberOfMessages: 1,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Len(t, es.Queue, 0)
			},
		},
		{
			name:             "we have queue, we are next",
			meNext:           true,
			es:               ExecutionState{Queue: []ExecutionQueueElement{{ctx: ctx}}},
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Len(t, es.Queue, 1)
			},
		},
		{
			name:             "running something without queue, pending execution",
			es:               ExecutionState{CurrentList: list},
			numberOfMessages: 2,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Len(t, es.Queue, 0)
				require.Equal(t, message.InPending, es.pending)
			},
		},
		{
			name:             "running something without queue, we're next",
			es:               ExecutionState{CurrentList: list},
			meNext: true,
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Len(t, es.Queue, 0)
			},
		},
		{
			name:             "in not confirmed pending and no queue, still message",
			es:               ExecutionState{pending: message.InPending},
			numberOfMessages: 1,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, message.NotPending, es.pending)
			},
		},
		{
			name:             "in not confirmed pending and no queue, we're next",
			es:               ExecutionState{pending: message.InPending},
			meNext: true,
			numberOfMessages: 0,
			checkES: func(t *testing.T, es *ExecutionState) {
				require.Equal(t, message.NotPending, es.pending)
			},
		},

	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			es := &test.es
			if es.CurrentList == nil {
				es.CurrentList = NewCurrentExecutionList()
			}
			messages := es.OnPulse(ctx, test.meNext)
			require.Equal(t, test.numberOfMessages, len(messages))
			if test.checkES != nil {
				test.checkES(t, es)
			}
		})
	}
}

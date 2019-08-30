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
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

func TestRequestsFetcher_New(t *testing.T) {
	rf := NewRequestsFetcher(gen.Reference(), nil, nil, nil)
	require.NotNil(t, rf)
}

func TestRequestsFetcher_FetchPendings(t *testing.T) {
	objectRef := gen.Reference()

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (insolar.Reference, artifacts.Client, ExecutionBrokerI)
	}{
		{
			name: "success",
			mocks: func(t minimock.Tester) (insolar.Reference, artifacts.Client, ExecutionBrokerI) {
				obj := gen.Reference()
				reqRef := gen.Reference()
				am := artifacts.NewClientMock(t).
					GetPendingsMock.Return([]insolar.Reference{reqRef}, nil).
					GetAbandonedRequestMock.Return(&record.IncomingRequest{Object: &objectRef}, nil)
				broker := NewExecutionBrokerIMock(t).
					IsKnownRequestMock.Return(false).
					AddRequestsFromLedgerMock.Return()

				return obj, am, broker
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)
			defer mc.Finish()
			defer mc.Wait(1 * time.Minute)

			obj, am, br := test.mocks(mc)
			rf := NewRequestsFetcher(obj, am, br, nil)
			rf.FetchPendings(ctx)
		})
	}
}

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
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
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
	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (insolar.Reference, artifacts.Client, ExecutionBrokerI)
	}{
		{
			name: "success",
			mocks: func(t minimock.Tester) (insolar.Reference, artifacts.Client, ExecutionBrokerI) {

				requestRef := gen.RecordReference()
				incoming := genIncomingRequest()

				am := artifacts.NewClientMock(t).
					GetPendingsMock.Return([]insolar.Reference{requestRef}, nil).
					GetAbandonedRequestMock.Return(incoming, nil)

				broker := NewExecutionBrokerIMock(t).
					IsKnownRequestMock.Return(false).
					AddRequestsFromLedgerMock.Return()

				return *incoming.Object, am, broker
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

func isKnownRequest(reference insolar.Reference, array []insolar.Reference) bool {
	for i := range array {
		if array[i].Equal(reference) {
			return false
		}
	}

	panic("should be known request")
}

func TestRequestsFetcher_Limits(t *testing.T) {
	ctx := inslogger.TestContext(t)

	mc := minimock.NewController(t)
	defer mc.Wait(1 * time.Second)

	objectRef := gen.Reference()
	pendings := gen.UniqueRecordReferences(20)

	am := artifacts.NewClientMock(mc).
		GetPendingsMock.Inspect(func(_ context.Context, objectRef insolar.Reference) {
		t.Logf("incoming obj %s", objectRef)
	}).Return(pendings, nil).
		GetAbandonedRequestMock.Set(
		func(ctx context.Context, objectRef insolar.Reference, requestRef insolar.Reference) (record.Request, error) {
			t.Logf("GetRequest %s", requestRef)
			return genIncomingRequest(), nil
		})

	eb := NewExecutionBrokerIMock(mc).
		IsKnownRequestMock.Return(false).
		AddRequestsFromLedgerMock.Return()

	fetcher := NewRequestsFetcher(objectRef, am, eb, nil)

	fetcherOriginal := fetcher.(*requestFetcher)
	err := fetcherOriginal.fetch(ctx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(len(pendings)), am.GetAbandonedRequestAfterCounter())
}

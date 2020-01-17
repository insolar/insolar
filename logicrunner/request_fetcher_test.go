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
	"sync"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/testutils"
)

func TestRequestsFetcher_New(t *testing.T) {
	defer testutils.LeakTester(t)

	rf := NewRequestsFetcher(gen.Reference(), nil, nil)
	require.NotNil(t, rf)
}

func TestRequestsFetcher_FetchPendings(t *testing.T) {
	defer testutils.LeakTester(t)

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (insolar.Reference, artifacts.Client)
	}{
		{
			name: "success",
			mocks: func(t minimock.Tester) (insolar.Reference, artifacts.Client) {

				requestRef := gen.RecordReference()
				incoming := genIncomingRequest()
				count := 0
				am := artifacts.NewClientMock(t).
					GetPendingsMock.Set(func(ctx context.Context, objectRef insolar.Reference, skip []insolar.ID) (ra1 []insolar.Reference, err error) {
					if count > 0 {
						return nil, insolar.ErrNoPendingRequest
					}
					count++
					return []insolar.Reference{requestRef}, nil
				}).
					GetRequestMock.Return(incoming, nil)

				return *incoming.Object, am
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			defer mc.Finish()
			defer mc.Wait(1 * time.Minute)

			obj, am := test.mocks(mc)
			rf := NewRequestsFetcher(obj, am, nil)
			feed := rf.FetchPendings(ctx)
			wg := &sync.WaitGroup{}
			wg.Add(1)
			var result []*common.Transcript
			go func() {
				for res := range feed {
					result = append(result, res)
				}
				wg.Done()
			}()
			wg.Wait()
			assert.Equal(t, 2, len(result))
			assert.NotEqual(t, nil, result[0])
			assert.Equal(t, (*common.Transcript)(nil), result[1])
		})
	}
}

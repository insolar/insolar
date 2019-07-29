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

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestHandleCall_CheckExecutionLoop(t *testing.T) {
	obj := gen.Reference()

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (*HandleCall, *record.IncomingRequest)
		loop  bool
	}{
		{
			name: "loop detected",
			loop: true,
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							GetExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									CheckExecutionLoopMock.Return(true),
							),
					},
				}
				req := &record.IncomingRequest{
					Object: &obj,
				}
				return h, req
			},
		},
		{
			name: "no loop, broker check",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							GetExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									CheckExecutionLoopMock.Return(false),
							),
					},
				}
				req := &record.IncomingRequest{
					Object: &obj,
				}
				return h, req
			},
		},
		{
			name: "no loop, not executing",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							GetExecutionStateMock.Expect(obj).
							Return( nil ),
					},
				}
				req := &record.IncomingRequest{
					Object: &obj,
				}
				return h, req
			},
		},
		{
			name: "no loop, nil object",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{},
				}
				req := &record.IncomingRequest{}
				return h, req
			},
		},
		{
			name: "no loop, constructor",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{},
				}
				req := &record.IncomingRequest{
					CallType: record.CTSaveAsChild,
				}
				return h, req
			},
		},
		{
			name: "no loop, no wait call",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{},
				}
				req := &record.IncomingRequest{
					ReturnMode: record.ReturnNoWait,
				}
				return h, req
			},
		},

	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			h, req := test.mocks(mc)
			loop := h.checkExecutionLoop(ctx, *req)
			require.Equal(t, test.loop, loop)

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}

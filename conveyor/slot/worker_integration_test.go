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

package slot

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/adapter/adapterstorage"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/matrix"
	"github.com/insolar/insolar/conveyor/handler"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

const maxState = fsm.StateID(1000)

// TODO: make mock stateMachine for integration tests using conveyor/generator package
type mockStateMachineSet struct {
	stateMachine matrix.StateMachine
}

func (s *mockStateMachineSet) GetStateMachineByID(id fsm.ID) matrix.StateMachine {
	return s.stateMachine
}

type mockStateMachineHolder struct{}

func (m *mockStateMachineHolder) makeSetAccessor() matrix.SetAccessor {
	return &mockStateMachineSet{
		stateMachine: m.GetStateMachinesByType(),
	}
}

func (m *mockStateMachineHolder) GetFutureConfig() matrix.SetAccessor {
	return m.makeSetAccessor()
}

func (m *mockStateMachineHolder) GetPresentConfig() matrix.SetAccessor {
	return m.makeSetAccessor()
}

func (m *mockStateMachineHolder) GetPastConfig() matrix.SetAccessor {
	return m.makeSetAccessor()
}

func (m *mockStateMachineHolder) GetInitialStateMachine() matrix.StateMachine {
	return m.GetStateMachinesByType()
}

func (m *mockStateMachineHolder) GetStateMachinesByType() matrix.StateMachine {

	sm := matrix.NewStateMachineMock(&testing.T{})
	sm.GetMigrationHandlerFunc = func(s fsm.StateID) (r handler.MigrationHandler) {
		return func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
			if s > maxState {
				s /= 2
			}
			return element.GetElementID(), fsm.NewElementState(fsm.ID(s%3), s+1), nil
		}
	}

	sm.GetTransitionHandlerFunc = func(s fsm.StateID) (r handler.TransitHandler) {
		return func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
			if s > maxState {
				s /= 2
			}
			return element.GetElementID(), fsm.NewElementState(fsm.ID(s%3), s+1), nil
		}
	}

	sm.GetResponseHandlerFunc = func(s fsm.StateID) (r handler.AdapterResponseHandler) {
		return func(element fsm.SlotElementHelper, response interface{}) (interface{}, fsm.ElementState, error) {
			if s > maxState {
				s /= 2
			}
			return element.GetPayload(), fsm.NewElementState(fsm.ID(s%3), s+1), nil
		}
	}

	return sm
}

func mockHandlerStorage() matrix.StateMachineHolder {
	return &mockStateMachineHolder{}
}

func initComponents(t *testing.T) {
	pc := testutils.NewPlatformCryptographyScheme()
	ledgerMock := testutils.NewLedgerLogicMock(t)
	ledgerMock.GetCodeFunc = func(p context.Context, p1 insolar.Parcel) (r insolar.Reply, r1 error) {
		return &reply.Code{}, nil
	}

	cm := &component.Manager{}
	ctx := context.TODO()

	components := adapterstorage.GetAllProcessors()
	components = append(components, pc, ledgerMock)
	cm.Inject(components...)
	err := cm.Init(ctx)
	if err != nil {
		t.Error("ComponentManager init failed", err)
	}
}

func setup() {
	HandlerStorage = mockHandlerStorage()
}

func testMainWrapper(m *testing.M) int {
	setup()
	code := m.Run()
	return code
}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func Test_run(t *testing.T) {
	initComponents(t)

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			for i := 1; i < 8000; i++ {
				state := fsm.StateID(i)
				if state > maxState {
					state /= maxState
				}
				element, err := slot.createElement(HandlerStorage.GetInitialStateMachine(), fsm.StateID(state), queue.OutputElement{})
				require.NoError(t, err)
				require.NotNil(t, element)
			}

			go func() {
				for i := 1; i < 10; i++ {
					resp := adapter.NewAdapterResponse(0, uint32(i), 0, 0)
					_ = slot.responseQueue.SinkPush(resp)
					time.Sleep(time.Millisecond * 50)
				}
			}()

			go func() {
				time.Sleep(time.Millisecond * 400)
				_ = slot.inputQueue.PushSignal(PendingPulseSignal, mockCallback())
			}()

			go func() {
				time.Sleep(time.Millisecond * 600)
				_ = slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback())
			}()

			go func() {
				time.Sleep(time.Millisecond * 800)
				_ = slot.inputQueue.PushSignal(CancelSignal, mockCallback())
			}()

			worker.run()

		})
	}
}

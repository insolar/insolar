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

package statestorage

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/executionarchive"
	"github.com/insolar/insolar/logicrunner/executionbroker"
)

type StateStorageSuite struct{ suite.Suite }

func TestStateStorage(t *testing.T) { suite.Run(t, new(StateStorageSuite)) }

func (s *StateStorageSuite) generateContext() context.Context {
	return inslogger.TestContext(s.T())
}

func (s *StateStorageSuite) generatePulse() insolar.Pulse {
	return insolar.Pulse{PulseNumber: gen.PulseNumber()}
}

func (s *StateStorageSuite) TestOnPulse() {
	mc := minimock.NewController(s.T())
	defer mc.Wait(2 * time.Second)

	ctx := s.generateContext()
	pulse := s.generatePulse()
	objectRef := gen.Reference()

	jc := jet.NewCoordinatorMock(mc).
		IsMeAuthorizedNowMock.Return(false, nil)

	ss := NewStateStorage(nil, nil, nil, jc, nil, nil, nil)
	rawStateStorage := ss.(*stateStorage)

	{ // empty state storage
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 0)
		s.Len(rawStateStorage.state, 0)
	}

	{ // state storage with empty object
		rawStateStorage.state[objectRef] = &ObjectState{}
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 0)
		s.Len(rawStateStorage.state, 0)
	}

	{ // state storage with empty execution archive
		rawStateStorage.state[objectRef] = &ObjectState{
			ExecutionArchive: executionarchive.NewExecutionArchiveMock(mc).
				OnPulseMock.Return(nil).
				IsEmptyMock.Return(true),
		}
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 0)
		s.Len(rawStateStorage.state, 0)
	}

	{ // state storage with non-empty execution archive
		rawStateStorage.state[objectRef] = &ObjectState{
			ExecutionArchive: executionarchive.NewExecutionArchiveMock(mc).
				OnPulseMock.Return([]insolar.Message{&message.StillExecuting{}}).
				IsEmptyMock.Return(false),
		}
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 1)
		s.Len(rawStateStorage.state, 1)

		delete(rawStateStorage.state, objectRef)
	}

	{ // state storage with execution archive and execution broker
		rawStateStorage.state[objectRef] = &ObjectState{
			ExecutionBroker: executionbroker.NewBrokerIMock(mc).
				OnPulseMock.Return([]insolar.Message{&message.ExecutorResults{}}),
			ExecutionArchive: executionarchive.NewExecutionArchiveMock(mc).
				OnPulseMock.Return(nil).
				IsEmptyMock.Return(true),
		}
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 1)
		s.Len(rawStateStorage.state, 0)
	}

	{ // state storage with multiple objects
		rawStateStorage.state[objectRef] = &ObjectState{
			ExecutionBroker: executionbroker.NewBrokerIMock(mc).
				OnPulseMock.Return([]insolar.Message{&message.ExecutorResults{}}),
			ExecutionArchive: executionarchive.NewExecutionArchiveMock(mc).
				OnPulseMock.Return([]insolar.Message{&message.StillExecuting{}}).
				IsEmptyMock.Return(false),
		}
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 2)
		s.Len(rawStateStorage.state, 1)

		delete(rawStateStorage.state, objectRef)
	}

	jc.IsMeAuthorizedNowMock.Return(false, nil)

	{ // state storage with multiple objects
		objectRef1 := gen.Reference()
		objectRef2 := gen.Reference()
		rawStateStorage.state[objectRef1] = &ObjectState{
			ExecutionBroker: executionbroker.NewBrokerIMock(mc).
				OnPulseMock.Return([]insolar.Message{&message.ExecutorResults{}}),
			ExecutionArchive: executionarchive.NewExecutionArchiveMock(mc).
				OnPulseMock.Return(nil).
				IsEmptyMock.Return(true),
		}
		rawStateStorage.state[objectRef2] = &ObjectState{
			ExecutionBroker: executionbroker.NewBrokerIMock(mc).
				OnPulseMock.Return([]insolar.Message{&message.ExecutorResults{}}),
			ExecutionArchive: executionarchive.NewExecutionArchiveMock(mc).
				OnPulseMock.Return([]insolar.Message{&message.StillExecuting{}}).
				IsEmptyMock.Return(false),
		}
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 3)
		s.Len(rawStateStorage.state, 1)
		s.NotNil(rawStateStorage.state[objectRef2].ExecutionArchive)
		s.Nil(rawStateStorage.state[objectRef2].ExecutionBroker)

		delete(rawStateStorage.state, objectRef2)
	}
}

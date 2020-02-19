// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/executionregistry"
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
	defer mc.Finish()

	ctx := s.generateContext()
	pulse := s.generatePulse()
	objectRef := gen.Reference()

	ss := NewStateStorage(nil, nil, nil, nil, nil, nil, nil, nil)
	rawStateStorage := ss.(*stateStorage)

	{ // empty state storage
		msgs := ss.OnPulse(ctx, pulse)
		s.Len(msgs, 0)
		s.Len(rawStateStorage.brokers, 0)
		s.Len(rawStateStorage.registries, 0)
	}

	{ // state storage with empty execution registry
		rawStateStorage.registries[objectRef] = executionregistry.NewExecutionRegistryMock(mc).
			OnPulseMock.Return(nil).
			IsEmptyMock.Return(true)
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 0)
		s.Len(rawStateStorage.brokers, 0)
		s.Len(rawStateStorage.registries, 0)
	}

	{ // state storage with non-empty execution registry
		rawStateStorage.registries[objectRef] = executionregistry.NewExecutionRegistryMock(mc).
			OnPulseMock.Return([]payload.Payload{&payload.StillExecuting{}}).
			IsEmptyMock.Return(false)
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 1)
		s.Len(rawStateStorage.brokers, 0)
		s.Len(rawStateStorage.registries, 1)

		delete(rawStateStorage.registries, objectRef)
	}

	{ // state storage with execution registry and execution broker
		rawStateStorage.registries[objectRef] = executionregistry.NewExecutionRegistryMock(mc).
			OnPulseMock.Return(nil).
			IsEmptyMock.Return(true)
		rawStateStorage.brokers[objectRef] = NewExecutionBrokerIMock(mc).
			OnPulseMock.Return([]payload.Payload{&payload.ExecutorResults{}})
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 1)
		s.Len(rawStateStorage.brokers, 0)
		s.Len(rawStateStorage.registries, 0)
	}

	{ // state storage with multiple objects
		rawStateStorage.brokers[objectRef] = NewExecutionBrokerIMock(mc).
			OnPulseMock.Return([]payload.Payload{&payload.ExecutorResults{}})
		rawStateStorage.registries[objectRef] = executionregistry.NewExecutionRegistryMock(mc).
			OnPulseMock.Return([]payload.Payload{&payload.StillExecuting{}}).
			IsEmptyMock.Return(false)
		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs[objectRef], 2)
		s.Len(rawStateStorage.brokers, 0)
		s.Len(rawStateStorage.registries, 1)

		delete(rawStateStorage.registries, objectRef)
	}

	{ // state storage with multiple objects
		objectRef1 := gen.Reference()
		objectRef2 := gen.Reference()

		rawStateStorage.brokers[objectRef1] = NewExecutionBrokerIMock(mc).
			OnPulseMock.Return([]payload.Payload{&payload.ExecutorResults{}})
		rawStateStorage.registries[objectRef1] = executionregistry.NewExecutionRegistryMock(mc).
			OnPulseMock.Return(nil).
			IsEmptyMock.Return(true)

		rawStateStorage.brokers[objectRef2] = NewExecutionBrokerIMock(mc).
			OnPulseMock.Return([]payload.Payload{&payload.ExecutorResults{}})
		rawStateStorage.registries[objectRef2] = executionregistry.NewExecutionRegistryMock(mc).
			OnPulseMock.Return([]payload.Payload{&payload.StillExecuting{}}).
			IsEmptyMock.Return(false)

		msgs := rawStateStorage.OnPulse(ctx, pulse)
		s.Len(msgs, 2)
		s.Len(msgs[objectRef1], 1)
		s.Len(msgs[objectRef2], 2)
		s.Len(rawStateStorage.brokers, 0)
		s.Len(rawStateStorage.registries, 1)
		s.NotNil(rawStateStorage.registries[objectRef2])
		s.Nil(rawStateStorage.brokers[objectRef2])

		delete(rawStateStorage.registries, objectRef2)
	}
}

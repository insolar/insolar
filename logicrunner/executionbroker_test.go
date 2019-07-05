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
	"math"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

// wait is Exponential retries waiting function
// example usage: require.True(wait(func))
func wait(check func() bool) bool {
	for i := 0; i < 16; i++ {
		time.Sleep(time.Millisecond * time.Duration(math.Pow(2, float64(i))))
		if check() {
			return true
		}
	}
	return false
}

type TranscriptDequeueSuite struct{ suite.Suite }

func TestTranscriptDequeue(t *testing.T) { suite.Run(t, new(TranscriptDequeueSuite)) }

func (s *TranscriptDequeueSuite) TestBasic() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	d.Push(&Transcript{Nonce: 1}, &Transcript{Nonce: 2})

	tr := d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(1), tr.Nonce)

	d.Prepend(&Transcript{Nonce: 3}, &Transcript{Nonce: 4})

	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(3), tr.Nonce)

	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(4), tr.Nonce)

	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(2), tr.Nonce)

	s.Nil(d.Pop())
}

func (s *TranscriptDequeueSuite) TestRotate() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	d.Push(&Transcript{Nonce: 1}, &Transcript{Nonce: 2})

	rotated := d.Rotate()
	s.Require().Len(rotated, 2)

	s.Nil(d.Pop())

	rotated = d.Rotate()
	s.Require().Len(rotated, 0)

	s.Nil(d.Pop())
}

func (s *TranscriptDequeueSuite) TestHasFromLedger() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	d.Prepend(&Transcript{Nonce: 3}, &Transcript{Nonce: 4})
	s.False(d.HasFromLedger() != nil)

	d.Push(&Transcript{FromLedger: true})
	s.True(d.HasFromLedger() != nil)

	d.Push(&Transcript{FromLedger: true})
	s.True(d.HasFromLedger() != nil)
}

func (s *TranscriptDequeueSuite) TestPopByReference() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref := gen.Reference()
	badRef := gen.Reference()

	d.Prepend(
		&Transcript{Nonce: 3, RequestRef: &badRef},
		&Transcript{Nonce: 4, RequestRef: &ref},
		&Transcript{Nonce: 5, RequestRef: &badRef},
	)

	tr := d.PopByReference(&ref)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref.Bytes())

	tr = d.PopByReference(&ref)
	s.Nil(tr)

	s.Equal(d.Len(), 2)
}

func (s *TranscriptDequeueSuite) TestTake() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	for i := 0; i < 15; i++ {
		d.Push(&Transcript{Nonce: uint64(i)})
	}

	trs := d.Take(0)
	s.Require().NotNil(d)
	s.Len(trs, 0)

	trs = d.Take(10)
	s.Require().NotNil(d)
	s.Len(trs, 10)

	trs = d.Take(10)
	s.Require().NotNil(d)
	s.Len(trs, 5)

	trs = d.Take(10)
	s.Require().NotNil(d)
	s.Len(trs, 0)
}

type ExecutionBrokerSuite struct {
	suite.Suite
	Context    context.Context
	Controller *minimock.Controller
}

func TestExecutionBroker(t *testing.T) { suite.Run(t, new(ExecutionBrokerSuite)) }

func (s *ExecutionBrokerSuite) BeforeTest(suiteName, testName string) {
	s.Context = inslogger.TestContext(s.T())
	s.Controller = minimock.NewController(s.T())
}

func waitOnChannel(channel chan struct{}) bool {
	select {
	case <-channel:
		return true
	case <-time.After(1 * time.Minute):
		return false
	}
}

func channelIsEmpty(channel chan struct{}) bool {
	select {
	case <-channel:
		return false
	default:
		return true
	}
}

func (s *ExecutionBrokerSuite) prepareLogicRunner(t *testing.T) *LogicRunner {
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	// initialize mocks
	am := artifacts.NewClientMock(s.Controller)
	dc := artifacts.NewDescriptorsCacheMock(s.Controller)
	mm := &mmanager{}
	re := NewRequestsExecutorMock(s.Controller)
	mb := testutils.NewMessageBusMock(s.Controller)
	jc := jet.NewCoordinatorMock(s.Controller)
	ps := pulse.NewAccessorMock(s.Controller)
	nn := network.NewNodeNetworkMock(s.Controller)

	// initialize lr
	lr.ArtifactManager = am
	lr.DescriptorsCache = dc
	lr.MessageBus = mb
	lr.MachinesManager = mm
	lr.JetCoordinator = jc
	lr.PulseAccessor = ps
	lr.NodeNetwork = nn
	lr.RequestsExecutor = re

	return lr
}

func (s *ExecutionBrokerSuite) TestPut() {
	lr := s.prepareLogicRunner(s.T())

	waitChannel := make(chan struct{})
	rem := lr.RequestsExecutor.(*RequestsExecutorMock)
	rem.ExecuteAndSaveMock.Set(func(_ context.Context, t *Transcript) (r insolar.Reply, r1 error) {
		if !t.Request.Immutable {
			waitChannel <- struct{}{}
		}
		return nil, nil
	})
	rem.SendReplyMock.Return()

	es := NewExecutionState(gen.Reference())
	es.RegisterLogicRunner(lr)
	b := es.Broker
	es.pending = message.NotPending

	processGoroutineExits := func() bool {
		b.processLock.Lock()
		defer b.processLock.Unlock()
		return b.processActive == false
	}

	reqRef1 := gen.Reference()
	tr := &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef1,
		Request:      &record.IncomingRequest{},
	}

	b.Put(s.Context, false, tr)
	s.Len(b.mutable.queue, 1)

	reqRef2 := gen.Reference()
	tr = &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef2,
		Request:      &record.IncomingRequest{},
	}

	b.Put(s.Context, true, tr)
	s.True(waitOnChannel(waitChannel), "failed to wait until put triggers start of queue processor")
	s.True(waitOnChannel(waitChannel), "failed to wait until queue processor'll finish processing")
	s.Require().True(wait(processGoroutineExits))
	s.Empty(waitChannel)

	s.Len(b.mutable.queue, 0)
	s.Len(b.immutable.queue, 0)
	s.Len(b.finished.queue, 2)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestPrepend() {
	lr := s.prepareLogicRunner(s.T())

	waitChannel := make(chan struct{})
	rem := lr.RequestsExecutor.(*RequestsExecutorMock)
	rem.ExecuteAndSaveMock.Set(func(_ context.Context, t *Transcript) (r insolar.Reply, r1 error) {
		if !t.Request.Immutable {
			waitChannel <- struct{}{}
		}
		return nil, nil
	})
	rem.SendReplyMock.Return()

	es := NewExecutionState(gen.Reference())
	es.RegisterLogicRunner(lr)
	b := es.Broker
	es.pending = message.NotPending

	processGoroutineExits := func() bool {
		b.processLock.Lock()
		defer b.processLock.Unlock()
		return b.processActive == false
	}

	reqRef1 := gen.Reference()
	tr := &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef1,
		Request:      &record.IncomingRequest{},
	}
	b.Prepend(s.Context, false, tr)
	s.Len(b.mutable.queue, 1)

	reqRef2 := gen.Reference()
	tr = &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef2,
		Request:      &record.IncomingRequest{},
	}
	b.Prepend(s.Context, true, tr)
	s.Require().True(waitOnChannel(waitChannel), "failed to wait until put triggers start of queue processor")
	s.Require().True(waitOnChannel(waitChannel), "failed to wait until queue processor'll finish processing")
	s.Require().True(wait(processGoroutineExits))
	s.Require().Empty(waitChannel)

	s.Len(b.mutable.queue, 0)
	s.Len(b.immutable.queue, 0)
	s.Len(b.finished.queue, 2)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestImmutable() {
	lr := s.prepareLogicRunner(s.T())

	waitMutableChannel := make(chan struct{})
	waitImmutableChannel := make(chan struct{})

	rem := lr.RequestsExecutor.(*RequestsExecutorMock)
	rem.ExecuteAndSaveMock.Set(func(_ context.Context, t *Transcript) (r insolar.Reply, r1 error) {
		if !t.Request.Immutable {
			waitMutableChannel <- struct{}{}
		} else {
			waitImmutableChannel <- struct{}{}
		}
		return nil, nil
	})
	rem.SendReplyMock.Return()

	es := NewExecutionState(gen.Reference())
	es.RegisterLogicRunner(lr)
	b := es.Broker
	es.pending = message.NotPending

	processGoroutineExits := func() bool {
		b.processLock.Lock()
		defer b.processLock.Unlock()
		return b.processActive == false
	}

	reqRef1 := gen.Reference()
	tr := &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef1,
		Request:      &record.IncomingRequest{Immutable: true},
	}

	b.Prepend(s.Context, false, tr)
	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))
	s.Require().Empty(waitMutableChannel)

	reqRef2 := gen.Reference()
	tr = &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef2,
		Request:      &record.IncomingRequest{Immutable: true},
	}

	b.Prepend(s.Context, true, tr)
	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))
	s.Require().Empty(waitMutableChannel)

	reqRef3 := gen.Reference()
	tr = &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef3,
		Request:      &record.IncomingRequest{Immutable: true},
	}

	// we can't process messages, do not do it
	es.Lock()
	es.pending = message.InPending
	es.Unlock()

	b.Prepend(s.Context, false, tr)
	s.Require().True(wait(func() bool { return b.immutable.Len() == 1 }), "failed to wait until immutable was put")
	s.Require().True(wait(processGoroutineExits))

	reqRef4 := gen.Reference()
	tr = &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef4,
		Request:      &record.IncomingRequest{Immutable: true},
	}

	b.Prepend(s.Context, true, tr)

	s.Require().True(wait(func() bool { return b.immutable.Len() == 2 }), "failed to wait until immutable was put")
	s.Require().True(wait(processGoroutineExits))
	s.Require().Empty(waitMutableChannel)

	b.StartProcessorIfNeeded(s.Context)
	s.Require().True(wait(processGoroutineExits))
	s.Empty(waitMutableChannel)
	s.Empty(waitImmutableChannel)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 2)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestRotate() {
	lr := s.prepareLogicRunner(s.T())

	rem := lr.RequestsExecutor.(*RequestsExecutorMock)
	rem.ExecuteAndSaveMock.Return(nil, nil)
	rem.SendReplyMock.Return()

	es := NewExecutionState(gen.Reference())
	es.RegisterLogicRunner(lr)
	b := es.Broker
	es.pending = message.NotPending

	for i := 0; i < 4; i++ {
		b.mutableLock.Lock()
		b.immutable.Push(&Transcript{})
		b.mutable.Push(&Transcript{})
		b.mutableLock.Unlock()
	}

	rotationResults := b.Rotate(10)
	s.Len(b.immutable.queue, 0)
	s.Len(b.mutable.queue, 0)
	s.Len(b.finished.queue, 0)
	s.Len(rotationResults.Requests, 8)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 4; i++ {
		b.mutableLock.Lock()
		b.immutable.Push(&Transcript{})
		b.mutableLock.Unlock()
	}

	rotationResults = b.Rotate(10)
	s.Len(b.immutable.queue, 0)
	s.Len(b.mutable.queue, 0)
	s.Len(b.finished.queue, 0)
	s.Len(rotationResults.Requests, 4)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 5; i++ {
		b.mutable.Push(&Transcript{})
		b.immutable.Push(&Transcript{})
	}

	rotationResults = b.Rotate(10)
	s.Len(b.immutable.queue, 0)
	s.Len(b.mutable.queue, 0)
	s.Len(b.finished.queue, 0)
	s.Len(rotationResults.Requests, 10)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 6; i++ {
		b.mutable.Push(&Transcript{})
		b.immutable.Push(&Transcript{})
	}

	rotationResults = b.Rotate(10)
	s.Len(b.immutable.queue, 0)
	s.Len(b.mutable.queue, 0)
	s.Len(b.finished.queue, 0)
	s.Len(rotationResults.Requests, 10)
	s.Len(rotationResults.Finished, 0)
	s.True(rotationResults.LedgerHasMoreRequests)
}

func (s *ExecutionBrokerSuite) TestDeduplication() {
	lr := s.prepareLogicRunner(s.T())

	rem := lr.RequestsExecutor.(*RequestsExecutorMock)
	rem.ExecuteAndSaveMock.Return(nil, nil)
	rem.SendReplyMock.Return()

	es := NewExecutionState(gen.Reference())
	es.RegisterLogicRunner(lr)
	b := es.Broker
	es.pending = message.InPending

	reqRef1 := gen.Reference()
	b.Put(s.Context, false, &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef1,
		Request:      &record.IncomingRequest{},
	}) // no duplication
	s.Len(b.mutable.queue, 1)

	b.Put(s.Context, false, &Transcript{
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   &reqRef1,
		Request:      &record.IncomingRequest{},
	}) // duplication
	s.Len(b.mutable.queue, 1)
}

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

	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar/bus"
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

type publisherMock struct{}

func (p *publisherMock) Publish(topic string, messages ...*wmMessage.Message) error {
	return nil
}

func (p *publisherMock) Close() error {
	return nil
}

// wait is Exponential retries waiting function
// example usage: require.True(wait(func))
func wait(check func(...interface{}) bool, args ...interface{}) bool {
	for i := 0; i < 16; i++ {
		time.Sleep(time.Millisecond * time.Duration(math.Pow(2, float64(i))))
		if check(args...) {
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

	// [] + [1, 2]
	d.Push(&Transcript{Nonce: 1}, &Transcript{Nonce: 2})

	// 1, [2]
	tr := d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(1), tr.Nonce)

	// [3, 4] + [2]
	d.Prepend(&Transcript{Nonce: 3}, &Transcript{Nonce: 4})

	// 3, [4, 2]
	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(3), tr.Nonce)

	// 4, [2]
	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(4), tr.Nonce)

	// 2, []
	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(2), tr.Nonce)

	// nil, []
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

	ref1, ref2, ref3 := gen.Reference(), gen.Reference(), gen.Reference()

	d.Prepend(
		&Transcript{Nonce: 3, RequestRef: ref1},
		&Transcript{Nonce: 4, RequestRef: ref2},
		&Transcript{Nonce: 5, RequestRef: ref3},
	)

	tr := d.PopByReference(ref2)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref2.Bytes())

	tr = d.PopByReference(ref2)
	s.Nil(tr)

	s.Nil(d.first.prev)
	s.Equal(d.first.next, d.last)
	s.Nil(d.last.next)
	s.Equal(d.last.prev, d.first)
	s.Equal(d.first.value.Nonce, uint64(3))
	s.Equal(d.last.value.Nonce, uint64(5))

	s.Equal(d.Length(), 2)
}

func (s *TranscriptDequeueSuite) TestPopByReferenceHead() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref1, ref2, ref3 := gen.Reference(), gen.Reference(), gen.Reference()
	el1 := &Transcript{Nonce: 3, RequestRef: ref1}
	el2 := &Transcript{Nonce: 4, RequestRef: ref2}
	el3 := &Transcript{Nonce: 5, RequestRef: ref3}
	d.Prepend(el1, el2, el3)

	tr := d.PopByReference(ref1)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref1.Bytes())

	s.Nil(d.first.prev)
	s.Equal(d.first.next, d.last)
	s.Nil(d.last.next)
	s.Equal(d.last.prev, d.first)
	s.Equal(d.first.value.Nonce, uint64(4))
	s.Equal(d.last.value.Nonce, uint64(5))

	s.Equal(d.Length(), 2)
}

func (s *TranscriptDequeueSuite) TestPopByReferenceTail() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref1, ref2, ref3 := gen.Reference(), gen.Reference(), gen.Reference()
	el1 := &Transcript{Nonce: 3, RequestRef: ref1}
	el2 := &Transcript{Nonce: 4, RequestRef: ref2}
	el3 := &Transcript{Nonce: 5, RequestRef: ref3}
	d.Prepend(el1, el2, el3)

	tr := d.PopByReference(ref3)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref3.Bytes())

	s.Nil(d.first.prev)
	s.Equal(d.first.next, d.last)
	s.Nil(d.last.next)
	s.Equal(d.last.prev, d.first)
	s.Equal(d.first.value.Nonce, uint64(3))
	s.Equal(d.last.value.Nonce, uint64(4))

	s.Equal(d.Length(), 2)
}

func (s *TranscriptDequeueSuite) TestPopByReferenceOneElement() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref1 := gen.Reference()

	d.Prepend(
		&Transcript{Nonce: 3, RequestRef: ref1},
	)

	tr := d.PopByReference(ref1)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref1.Bytes())

	s.Nil(d.first)
	s.Nil(d.last)
	s.Equal(d.Length(), 0)
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
	sender := bus.NewSenderMock(s.Controller)
	pm := &publisherMock{}
	lr, _ := NewLogicRunner(&configuration.LogicRunner{}, pm, sender)

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

	_ = lr.Init(s.Context)

	return lr
}

func immutableCount(args ...interface{}) bool {
	broker := args[0].(*ExecutionBroker)
	count := args[1].(int)

	return broker.immutable.Length() >= count
}

func finishedCount(args ...interface{}) bool {
	broker := args[0].(*ExecutionBroker)
	count := args[1].(int)

	return broker.finished.Length() >= count
}

func processorStatus(args ...interface{}) bool {
	broker := args[0].(*ExecutionBroker)
	status := args[1].(bool)

	return broker.isActiveProcessor() == status
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

	objectRef := gen.Reference()
	b := lr.StateStorage.UpsertExecutionState(objectRef)
	b.executionState.pending = message.NotPending

	tr := NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{})

	b.Put(s.Context, false, tr)
	s.Equal(b.mutable.Length(), 1)

	tr = NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{})

	b.Put(s.Context, true, tr)
	s.True(waitOnChannel(waitChannel), "failed to wait until put triggers start of queue processor")
	s.True(waitOnChannel(waitChannel), "failed to wait until queue processor'll finish processing")
	s.Require().True(wait(processorStatus, b, false))
	s.Empty(waitChannel)

	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.finished.Length(), 2)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)

	_ = lr.Stop(s.Context)
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

	objectRef := gen.Reference()
	b := lr.StateStorage.UpsertExecutionState(objectRef)
	b.executionState.pending = message.NotPending

	reqRef1 := gen.Reference()
	tr := NewTranscript(s.Context, reqRef1, record.IncomingRequest{})
	b.Prepend(s.Context, false, tr)
	s.Equal(b.mutable.Length(), 1)

	reqRef2 := gen.Reference()
	tr = NewTranscript(s.Context, reqRef2, record.IncomingRequest{})
	b.Prepend(s.Context, true, tr)
	s.Require().True(waitOnChannel(waitChannel), "failed to wait until put triggers start of queue processor")
	s.Require().True(waitOnChannel(waitChannel), "failed to wait until queue processor'll finish processing")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitChannel)

	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.finished.Length(), 2)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)

	_ = lr.Stop(s.Context)
}

func (s *ExecutionBrokerSuite) TestImmutable_NotPending() {
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

	objectRef := gen.Reference()
	b := lr.StateStorage.UpsertExecutionState(objectRef)
	b.executionState.pending = message.NotPending

	tr := NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, false, tr)
	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitMutableChannel)

	tr = NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, true, tr)
	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitMutableChannel)

	_ = lr.Stop(s.Context)
}

func (s *ExecutionBrokerSuite) TestImmutable_InPending() {
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

	objectRef := gen.Reference()
	b := lr.StateStorage.UpsertExecutionState(objectRef)
	b.executionState.pending = message.InPending

	tr := NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, false, tr)
	s.Require().True(wait(immutableCount, b, 1), "failed to wait until immutable was put")
	s.Require().True(wait(processorStatus, b, false))

	tr = NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, true, tr)

	s.Require().True(wait(immutableCount, b, 2), "failed to wait until immutable was put")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitMutableChannel)

	b.StartProcessorIfNeeded(s.Context)
	s.Require().True(wait(processorStatus, b, false))
	s.Empty(waitMutableChannel)
	s.Empty(waitImmutableChannel)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 2)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 0)

	_ = lr.Stop(s.Context)
}

func (s *ExecutionBrokerSuite) TestRotate() {
	lr := s.prepareLogicRunner(s.T())

	rem := lr.RequestsExecutor.(*RequestsExecutorMock)
	rem.ExecuteAndSaveMock.Return(nil, nil)
	rem.SendReplyMock.Return()

	objectRef := gen.Reference()
	b := lr.StateStorage.UpsertExecutionState(objectRef)
	b.executionState.pending = message.NotPending

	for i := 0; i < 4; i++ {
		b.stateLock.Lock()
		b.immutable.Push(&Transcript{})
		b.mutable.Push(&Transcript{})
		b.stateLock.Unlock()
	}

	rotationResults := b.Rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 8)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 4; i++ {
		b.stateLock.Lock()
		b.immutable.Push(&Transcript{})
		b.stateLock.Unlock()
	}

	rotationResults = b.Rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 4)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 5; i++ {
		b.mutable.Push(&Transcript{})
		b.immutable.Push(&Transcript{})
	}

	rotationResults = b.Rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 10)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 6; i++ {
		b.mutable.Push(&Transcript{})
		b.immutable.Push(&Transcript{})
	}

	rotationResults = b.Rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 10)
	s.Len(rotationResults.Finished, 0)
	s.True(rotationResults.LedgerHasMoreRequests)

	_ = lr.Stop(s.Context)
}

func (s *ExecutionBrokerSuite) TestDeduplication() {
	lr := s.prepareLogicRunner(s.T())

	rem := lr.RequestsExecutor.(*RequestsExecutorMock)
	rem.ExecuteAndSaveMock.Return(nil, nil)
	rem.SendReplyMock.Return()

	objectRef := gen.Reference()
	b := lr.StateStorage.UpsertExecutionState(objectRef)
	b.executionState.pending = message.InPending

	reqRef1 := gen.Reference()
	b.Put(s.Context, false, NewTranscript(s.Context, reqRef1, record.IncomingRequest{})) // no duplication
	s.Equal(b.mutable.Length(), 1)

	b.Put(s.Context, false, NewTranscript(s.Context, reqRef1, record.IncomingRequest{})) // duplication
	s.Equal(b.mutable.Length(), 1)

	_ = lr.Stop(s.Context)
}

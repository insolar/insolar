package logicrunner

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
)

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

type ExecutionBrokerSuite struct{ suite.Suite }

func TestExecutionBroker(t *testing.T) { suite.Run(t, new(ExecutionBrokerSuite)) }

// wait is Exponential retries waiting function
// example usage: require.True(wait(func))
func wait(check func() bool) bool {
	for i := 0; i < 15; i++ {
		time.Sleep(time.Millisecond * time.Duration(math.Pow(2, float64(i))))
		if check() {
			return true
		}
	}
	return false
}

func (s *ExecutionBrokerSuite) TestPut() {
	ctx := context.TODO()
	allowProcessingFunc := func(_ context.Context) error { return nil }
	executeFunc := func(_ context.Context, _ *Transcript, _ interface{}) error { return nil }

	b := NewExecutionBroker(allowProcessingFunc, executeFunc, nil)
	processGoroutineExits := func() bool {
		b.processLock.Lock()
		defer b.processLock.Unlock()
		return b.processActive == false
	}

	tr := &Transcript{LogicContext: &insolar.LogicCallContext{Immutable: false}}
	b.Put(ctx, false, tr)
	b.processLock.Lock()
	s.False(b.processActive) // this flag should be disbaled
	b.processLock.Unlock()

	s.Len(b.mutable.queue, 1)

	b.Put(ctx, true, tr)
	b.processLock.Lock()
	s.True(b.processActive) // this flag should be enabled
	b.processLock.Unlock()

	finishProcessing := func() bool { return b.finished.Len() == 2 }
	s.Require().True(wait(finishProcessing), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))

	s.Len(b.mutable.queue, 0)
	s.Len(b.immutable.queue, 0)
	s.Len(b.finished.queue, 2)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestPrepend() {
	ctx := context.TODO()
	allowProcessingFunc := func(_ context.Context) error { return nil }
	executeFunc := func(_ context.Context, _ *Transcript, _ interface{}) error { return nil }

	b := NewExecutionBroker(allowProcessingFunc, executeFunc, nil)
	processGoroutineExits := func() bool {
		b.processLock.Lock()
		defer b.processLock.Unlock()
		return b.processActive == false
	}

	tr := &Transcript{LogicContext: &insolar.LogicCallContext{Immutable: false}}
	b.Prepend(ctx, false, tr)
	b.processLock.Lock()
	s.False(b.processActive) // this flag should be disbaled
	b.processLock.Unlock()

	s.Len(b.mutable.queue, 1)

	b.Prepend(ctx, true, tr)
	b.processLock.Lock()
	s.True(b.processActive) // this flag should be enabled
	b.processLock.Unlock()

	finishProcessing := func() bool { return b.finished.Len() == 2 }
	s.Require().True(wait(finishProcessing), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))
	s.Len(b.mutable.queue, 0)
	s.Len(b.immutable.queue, 0)
	s.Len(b.finished.queue, 2)

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestImmutable() {

	ctx := context.TODO()
	forbidProcessingFunc := func(_ context.Context) error { return ErrRetryLater }
	allowProcessingFunc := func(_ context.Context) error { return nil }
	executeFunc := func(_ context.Context, _ *Transcript, _ interface{}) error { return nil }

	b := NewExecutionBroker(allowProcessingFunc, executeFunc, nil)
	processGoroutineExits := func() bool {
		b.processLock.Lock()
		defer b.processLock.Unlock()
		return b.processActive == false
	}

	tr := &Transcript{LogicContext: &insolar.LogicCallContext{Immutable: true}}
	b.Prepend(ctx, false, tr)
	b.processLock.Lock()
	s.False(b.processActive) // we're not starting processor, but job is started
	b.processLock.Unlock()

	finishProcessing := func() bool { return b.finished.Len() == 1 }
	s.Require().True(wait(finishProcessing), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))

	b.Prepend(ctx, true, tr)
	b.processLock.Lock()
	s.True(b.processActive) // we're not starting processor, but job is started
	b.processLock.Unlock()
	finishProcessing = func() bool { return b.finished.Len() == 2 }
	s.Require().True(wait(finishProcessing), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))

	// we can't process messages, do not do it
	b.checkFunc = forbidProcessingFunc
	b.Prepend(ctx, false, tr)
	b.processLock.Lock()
	s.False(b.processActive) // we're not starting processor, but job is started
	b.processLock.Unlock()

	finishProcessing = func() bool { return b.immutable.Len() == 1 }
	s.Require().True(wait(finishProcessing), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))

	b.Prepend(ctx, true, tr)
	b.processLock.Lock()
	s.True(b.processActive) // we're not starting processor, but job is started
	b.processLock.Unlock()
	finishProcessing = func() bool { return b.immutable.Len() == 2 }
	s.Require().True(wait(finishProcessing), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))

	b.StartProcessorIfNeeded(ctx)
	b.processLock.Lock()
	s.True(b.processActive) // we're not starting processor, but job is started
	b.processLock.Unlock()
	finishProcessing = func() bool {
		b.processLock.Lock()
		defer b.processLock.Unlock()
		return b.processActive == false
	}
	s.Require().True(wait(finishProcessing), "failed to wait while processing is finished")
	s.Require().True(wait(processGoroutineExits))

	rotationResults := b.Rotate(10)
	s.Len(rotationResults.Requests, 2)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestRotate() {
	b := NewExecutionBroker(nil, nil, nil)

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

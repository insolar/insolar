package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexLifelineAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexLifelineAccessorMock implements github.com/insolar/insolar/ledger/object.IndexLifelineAccessor
type IndexLifelineAccessorMock struct {
	t minimock.Tester

	LifelineForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)
	LifelineForIDCounter    uint64
	LifelineForIDPreCounter uint64
	LifelineForIDMock       mIndexLifelineAccessorMockLifelineForID
}

// NewIndexLifelineAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.IndexLifelineAccessor
func NewIndexLifelineAccessorMock(t minimock.Tester) *IndexLifelineAccessorMock {
	m := &IndexLifelineAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LifelineForIDMock = mIndexLifelineAccessorMockLifelineForID{mock: m}

	return m
}

type mIndexLifelineAccessorMockLifelineForID struct {
	mock              *IndexLifelineAccessorMock
	mainExpectation   *IndexLifelineAccessorMockLifelineForIDExpectation
	expectationSeries []*IndexLifelineAccessorMockLifelineForIDExpectation
}

type IndexLifelineAccessorMockLifelineForIDExpectation struct {
	input  *IndexLifelineAccessorMockLifelineForIDInput
	result *IndexLifelineAccessorMockLifelineForIDResult
}

type IndexLifelineAccessorMockLifelineForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type IndexLifelineAccessorMockLifelineForIDResult struct {
	r  Lifeline
	r1 error
}

// Expect specifies that invocation of IndexLifelineAccessor.LifelineForID is expected from 1 to Infinity times
func (m *mIndexLifelineAccessorMockLifelineForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mIndexLifelineAccessorMockLifelineForID {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLifelineAccessorMockLifelineForIDExpectation{}
	}
	m.mainExpectation.input = &IndexLifelineAccessorMockLifelineForIDInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of IndexLifelineAccessor.LifelineForID
func (m *mIndexLifelineAccessorMockLifelineForID) Return(r Lifeline, r1 error) *IndexLifelineAccessorMock {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLifelineAccessorMockLifelineForIDExpectation{}
	}
	m.mainExpectation.result = &IndexLifelineAccessorMockLifelineForIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexLifelineAccessor.LifelineForID is expected once
func (m *mIndexLifelineAccessorMockLifelineForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *IndexLifelineAccessorMockLifelineForIDExpectation {
	m.mock.LifelineForIDFunc = nil
	m.mainExpectation = nil

	expectation := &IndexLifelineAccessorMockLifelineForIDExpectation{}
	expectation.input = &IndexLifelineAccessorMockLifelineForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexLifelineAccessorMockLifelineForIDExpectation) Return(r Lifeline, r1 error) {
	e.result = &IndexLifelineAccessorMockLifelineForIDResult{r, r1}
}

// Set uses given function f as a mock of IndexLifelineAccessor.LifelineForID method
func (m *mIndexLifelineAccessorMockLifelineForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)) *IndexLifelineAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LifelineForIDFunc = f
	return m.mock
}

// LifelineForID implements github.com/insolar/insolar/ledger/object.IndexLifelineAccessor interface
func (m *IndexLifelineAccessorMock) LifelineForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.LifelineForIDPreCounter, 1)
	defer atomic.AddUint64(&m.LifelineForIDCounter, 1)

	if len(m.LifelineForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LifelineForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexLifelineAccessorMock.LifelineForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.LifelineForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexLifelineAccessorMockLifelineForIDInput{p, p1, p2}, "IndexLifelineAccessor.LifelineForID got unexpected parameters")

		result := m.LifelineForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexLifelineAccessorMock.LifelineForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDMock.mainExpectation != nil {

		input := m.LifelineForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexLifelineAccessorMockLifelineForIDInput{p, p1, p2}, "IndexLifelineAccessor.LifelineForID got unexpected parameters")
		}

		result := m.LifelineForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexLifelineAccessorMock.LifelineForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDFunc == nil {
		m.t.Fatalf("Unexpected call to IndexLifelineAccessorMock.LifelineForID. %v %v %v", p, p1, p2)
		return
	}

	return m.LifelineForIDFunc(p, p1, p2)
}

// LifelineForIDMinimockCounter returns a count of IndexLifelineAccessorMock.LifelineForIDFunc invocations
func (m *IndexLifelineAccessorMock) LifelineForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDCounter)
}

// LifelineForIDMinimockPreCounter returns the value of IndexLifelineAccessorMock.LifelineForID invocations
func (m *IndexLifelineAccessorMock) LifelineForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDPreCounter)
}

// LifelineForIDFinished returns true if mock invocations count is ok
func (m *IndexLifelineAccessorMock) LifelineForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LifelineForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LifelineForIDCounter) == uint64(len(m.LifelineForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LifelineForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LifelineForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LifelineForIDFunc != nil {
		return atomic.LoadUint64(&m.LifelineForIDCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexLifelineAccessorMock) ValidateCallCounters() {

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to IndexLifelineAccessorMock.LifelineForID")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexLifelineAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexLifelineAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexLifelineAccessorMock) MinimockFinish() {

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to IndexLifelineAccessorMock.LifelineForID")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexLifelineAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexLifelineAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.LifelineForIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.LifelineForIDFinished() {
				m.t.Error("Expected call to IndexLifelineAccessorMock.LifelineForID")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *IndexLifelineAccessorMock) AllMocksCalled() bool {

	if !m.LifelineForIDFinished() {
		return false
	}

	return true
}

package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexAccessorMock implements github.com/insolar/insolar/ledger/object.IndexAccessor
type IndexAccessorMock struct {
	t minimock.Tester

	LifelineForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)
	LifelineForIDCounter    uint64
	LifelineForIDPreCounter uint64
	LifelineForIDMock       mIndexAccessorMockLifelineForID
}

// NewIndexAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.IndexAccessor
func NewIndexAccessorMock(t minimock.Tester) *IndexAccessorMock {
	m := &IndexAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LifelineForIDMock = mIndexAccessorMockLifelineForID{mock: m}

	return m
}

type mIndexAccessorMockLifelineForID struct {
	mock              *IndexAccessorMock
	mainExpectation   *IndexAccessorMockLifelineForIDExpectation
	expectationSeries []*IndexAccessorMockLifelineForIDExpectation
}

type IndexAccessorMockLifelineForIDExpectation struct {
	input  *IndexAccessorMockLifelineForIDInput
	result *IndexAccessorMockLifelineForIDResult
}

type IndexAccessorMockLifelineForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type IndexAccessorMockLifelineForIDResult struct {
	r  Lifeline
	r1 error
}

// Expect specifies that invocation of IndexAccessor.LifelineForID is expected from 1 to Infinity times
func (m *mIndexAccessorMockLifelineForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mIndexAccessorMockLifelineForID {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockLifelineForIDExpectation{}
	}
	m.mainExpectation.input = &IndexAccessorMockLifelineForIDInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of IndexAccessor.LifelineForID
func (m *mIndexAccessorMockLifelineForID) Return(r Lifeline, r1 error) *IndexAccessorMock {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockLifelineForIDExpectation{}
	}
	m.mainExpectation.result = &IndexAccessorMockLifelineForIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexAccessor.LifelineForID is expected once
func (m *mIndexAccessorMockLifelineForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *IndexAccessorMockLifelineForIDExpectation {
	m.mock.LifelineForIDFunc = nil
	m.mainExpectation = nil

	expectation := &IndexAccessorMockLifelineForIDExpectation{}
	expectation.input = &IndexAccessorMockLifelineForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexAccessorMockLifelineForIDExpectation) Return(r Lifeline, r1 error) {
	e.result = &IndexAccessorMockLifelineForIDResult{r, r1}
}

// Set uses given function f as a mock of IndexAccessor.LifelineForID method
func (m *mIndexAccessorMockLifelineForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)) *IndexAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LifelineForIDFunc = f
	return m.mock
}

// LifelineForID implements github.com/insolar/insolar/ledger/object.IndexAccessor interface
func (m *IndexAccessorMock) LifelineForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.LifelineForIDPreCounter, 1)
	defer atomic.AddUint64(&m.LifelineForIDCounter, 1)

	if len(m.LifelineForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LifelineForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexAccessorMock.LifelineForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.LifelineForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexAccessorMockLifelineForIDInput{p, p1, p2}, "IndexAccessor.LifelineForID got unexpected parameters")

		result := m.LifelineForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.LifelineForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDMock.mainExpectation != nil {

		input := m.LifelineForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexAccessorMockLifelineForIDInput{p, p1, p2}, "IndexAccessor.LifelineForID got unexpected parameters")
		}

		result := m.LifelineForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.LifelineForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDFunc == nil {
		m.t.Fatalf("Unexpected call to IndexAccessorMock.LifelineForID. %v %v %v", p, p1, p2)
		return
	}

	return m.LifelineForIDFunc(p, p1, p2)
}

// LifelineForIDMinimockCounter returns a count of IndexAccessorMock.LifelineForIDFunc invocations
func (m *IndexAccessorMock) LifelineForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDCounter)
}

// LifelineForIDMinimockPreCounter returns the value of IndexAccessorMock.LifelineForID invocations
func (m *IndexAccessorMock) LifelineForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDPreCounter)
}

// LifelineForIDFinished returns true if mock invocations count is ok
func (m *IndexAccessorMock) LifelineForIDFinished() bool {
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
func (m *IndexAccessorMock) ValidateCallCounters() {

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.LifelineForID")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexAccessorMock) MinimockFinish() {

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.LifelineForID")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexAccessorMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IndexAccessorMock.LifelineForID")
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
func (m *IndexAccessorMock) AllMocksCalled() bool {

	if !m.LifelineForIDFinished() {
		return false
	}

	return true
}

package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PendingFilamentStateAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// PendingFilamentStateAccessorMock implements github.com/insolar/insolar/ledger/object.PendingFilamentStateAccessor
type PendingFilamentStateAccessorMock struct {
	t minimock.Tester

	WaitForRefreshFunc func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r <-chan struct {
	}, r1 error)
	WaitForRefreshCounter    uint64
	WaitForRefreshPreCounter uint64
	WaitForRefreshMock       mPendingFilamentStateAccessorMockWaitForRefresh
}

// NewPendingFilamentStateAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.PendingFilamentStateAccessor
func NewPendingFilamentStateAccessorMock(t minimock.Tester) *PendingFilamentStateAccessorMock {
	m := &PendingFilamentStateAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.WaitForRefreshMock = mPendingFilamentStateAccessorMockWaitForRefresh{mock: m}

	return m
}

type mPendingFilamentStateAccessorMockWaitForRefresh struct {
	mock              *PendingFilamentStateAccessorMock
	mainExpectation   *PendingFilamentStateAccessorMockWaitForRefreshExpectation
	expectationSeries []*PendingFilamentStateAccessorMockWaitForRefreshExpectation
}

type PendingFilamentStateAccessorMockWaitForRefreshExpectation struct {
	input  *PendingFilamentStateAccessorMockWaitForRefreshInput
	result *PendingFilamentStateAccessorMockWaitForRefreshResult
}

type PendingFilamentStateAccessorMockWaitForRefreshInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type PendingFilamentStateAccessorMockWaitForRefreshResult struct {
	r <-chan struct {
	}
	r1 error
}

// Expect specifies that invocation of PendingFilamentStateAccessor.WaitForRefresh is expected from 1 to Infinity times
func (m *mPendingFilamentStateAccessorMockWaitForRefresh) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mPendingFilamentStateAccessorMockWaitForRefresh {
	m.mock.WaitForRefreshFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingFilamentStateAccessorMockWaitForRefreshExpectation{}
	}
	m.mainExpectation.input = &PendingFilamentStateAccessorMockWaitForRefreshInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of PendingFilamentStateAccessor.WaitForRefresh
func (m *mPendingFilamentStateAccessorMockWaitForRefresh) Return(r <-chan struct {
}, r1 error) *PendingFilamentStateAccessorMock {
	m.mock.WaitForRefreshFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingFilamentStateAccessorMockWaitForRefreshExpectation{}
	}
	m.mainExpectation.result = &PendingFilamentStateAccessorMockWaitForRefreshResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of PendingFilamentStateAccessor.WaitForRefresh is expected once
func (m *mPendingFilamentStateAccessorMockWaitForRefresh) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *PendingFilamentStateAccessorMockWaitForRefreshExpectation {
	m.mock.WaitForRefreshFunc = nil
	m.mainExpectation = nil

	expectation := &PendingFilamentStateAccessorMockWaitForRefreshExpectation{}
	expectation.input = &PendingFilamentStateAccessorMockWaitForRefreshInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingFilamentStateAccessorMockWaitForRefreshExpectation) Return(r <-chan struct {
}, r1 error) {
	e.result = &PendingFilamentStateAccessorMockWaitForRefreshResult{r, r1}
}

// Set uses given function f as a mock of PendingFilamentStateAccessor.WaitForRefresh method
func (m *mPendingFilamentStateAccessorMockWaitForRefresh) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r <-chan struct {
}, r1 error)) *PendingFilamentStateAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WaitForRefreshFunc = f
	return m.mock
}

// WaitForRefresh implements github.com/insolar/insolar/ledger/object.PendingFilamentStateAccessor interface
func (m *PendingFilamentStateAccessorMock) WaitForRefresh(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r <-chan struct {
}, r1 error) {
	counter := atomic.AddUint64(&m.WaitForRefreshPreCounter, 1)
	defer atomic.AddUint64(&m.WaitForRefreshCounter, 1)

	if len(m.WaitForRefreshMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WaitForRefreshMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingFilamentStateAccessorMock.WaitForRefresh. %v %v %v", p, p1, p2)
			return
		}

		input := m.WaitForRefreshMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingFilamentStateAccessorMockWaitForRefreshInput{p, p1, p2}, "PendingFilamentStateAccessor.WaitForRefresh got unexpected parameters")

		result := m.WaitForRefreshMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingFilamentStateAccessorMock.WaitForRefresh")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WaitForRefreshMock.mainExpectation != nil {

		input := m.WaitForRefreshMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingFilamentStateAccessorMockWaitForRefreshInput{p, p1, p2}, "PendingFilamentStateAccessor.WaitForRefresh got unexpected parameters")
		}

		result := m.WaitForRefreshMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingFilamentStateAccessorMock.WaitForRefresh")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WaitForRefreshFunc == nil {
		m.t.Fatalf("Unexpected call to PendingFilamentStateAccessorMock.WaitForRefresh. %v %v %v", p, p1, p2)
		return
	}

	return m.WaitForRefreshFunc(p, p1, p2)
}

// WaitForRefreshMinimockCounter returns a count of PendingFilamentStateAccessorMock.WaitForRefreshFunc invocations
func (m *PendingFilamentStateAccessorMock) WaitForRefreshMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WaitForRefreshCounter)
}

// WaitForRefreshMinimockPreCounter returns the value of PendingFilamentStateAccessorMock.WaitForRefresh invocations
func (m *PendingFilamentStateAccessorMock) WaitForRefreshMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WaitForRefreshPreCounter)
}

// WaitForRefreshFinished returns true if mock invocations count is ok
func (m *PendingFilamentStateAccessorMock) WaitForRefreshFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.WaitForRefreshMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.WaitForRefreshCounter) == uint64(len(m.WaitForRefreshMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.WaitForRefreshMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.WaitForRefreshCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.WaitForRefreshFunc != nil {
		return atomic.LoadUint64(&m.WaitForRefreshCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingFilamentStateAccessorMock) ValidateCallCounters() {

	if !m.WaitForRefreshFinished() {
		m.t.Fatal("Expected call to PendingFilamentStateAccessorMock.WaitForRefresh")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingFilamentStateAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PendingFilamentStateAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PendingFilamentStateAccessorMock) MinimockFinish() {

	if !m.WaitForRefreshFinished() {
		m.t.Fatal("Expected call to PendingFilamentStateAccessorMock.WaitForRefresh")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PendingFilamentStateAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *PendingFilamentStateAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.WaitForRefreshFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.WaitForRefreshFinished() {
				m.t.Error("Expected call to PendingFilamentStateAccessorMock.WaitForRefresh")
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
func (m *PendingFilamentStateAccessorMock) AllMocksCalled() bool {

	if !m.WaitForRefreshFinished() {
		return false
	}

	return true
}

package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FilamentCalculator" can be found in github.com/insolar/insolar/ledger/light/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//FilamentCalculatorMock implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator
type FilamentCalculatorMock struct {
	t minimock.Tester

	PendingRequestsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []insolar.ID, r1 error)
	PendingRequestsCounter    uint64
	PendingRequestsPreCounter uint64
	PendingRequestsMock       mFilamentCalculatorMockPendingRequests

	RequestsFunc       func(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.PulseNumber) (r []record.CompositeFilamentRecord, r1 error)
	RequestsCounter    uint64
	RequestsPreCounter uint64
	RequestsMock       mFilamentCalculatorMockRequests
}

//NewFilamentCalculatorMock returns a mock for github.com/insolar/insolar/ledger/light/executor.FilamentCalculator
func NewFilamentCalculatorMock(t minimock.Tester) *FilamentCalculatorMock {
	m := &FilamentCalculatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PendingRequestsMock = mFilamentCalculatorMockPendingRequests{mock: m}
	m.RequestsMock = mFilamentCalculatorMockRequests{mock: m}

	return m
}

type mFilamentCalculatorMockPendingRequests struct {
	mock              *FilamentCalculatorMock
	mainExpectation   *FilamentCalculatorMockPendingRequestsExpectation
	expectationSeries []*FilamentCalculatorMockPendingRequestsExpectation
}

type FilamentCalculatorMockPendingRequestsExpectation struct {
	input  *FilamentCalculatorMockPendingRequestsInput
	result *FilamentCalculatorMockPendingRequestsResult
}

type FilamentCalculatorMockPendingRequestsInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type FilamentCalculatorMockPendingRequestsResult struct {
	r  []insolar.ID
	r1 error
}

//Expect specifies that invocation of FilamentCalculator.PendingRequests is expected from 1 to Infinity times
func (m *mFilamentCalculatorMockPendingRequests) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mFilamentCalculatorMockPendingRequests {
	m.mock.PendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockPendingRequestsExpectation{}
	}
	m.mainExpectation.input = &FilamentCalculatorMockPendingRequestsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of FilamentCalculator.PendingRequests
func (m *mFilamentCalculatorMockPendingRequests) Return(r []insolar.ID, r1 error) *FilamentCalculatorMock {
	m.mock.PendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockPendingRequestsExpectation{}
	}
	m.mainExpectation.result = &FilamentCalculatorMockPendingRequestsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCalculator.PendingRequests is expected once
func (m *mFilamentCalculatorMockPendingRequests) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *FilamentCalculatorMockPendingRequestsExpectation {
	m.mock.PendingRequestsFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCalculatorMockPendingRequestsExpectation{}
	expectation.input = &FilamentCalculatorMockPendingRequestsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCalculatorMockPendingRequestsExpectation) Return(r []insolar.ID, r1 error) {
	e.result = &FilamentCalculatorMockPendingRequestsResult{r, r1}
}

//Set uses given function f as a mock of FilamentCalculator.PendingRequests method
func (m *mFilamentCalculatorMockPendingRequests) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []insolar.ID, r1 error)) *FilamentCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PendingRequestsFunc = f
	return m.mock
}

//PendingRequests implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator interface
func (m *FilamentCalculatorMock) PendingRequests(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.PendingRequestsPreCounter, 1)
	defer atomic.AddUint64(&m.PendingRequestsCounter, 1)

	if len(m.PendingRequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PendingRequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCalculatorMock.PendingRequests. %v %v %v", p, p1, p2)
			return
		}

		input := m.PendingRequestsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCalculatorMockPendingRequestsInput{p, p1, p2}, "FilamentCalculator.PendingRequests got unexpected parameters")

		result := m.PendingRequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.PendingRequests")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PendingRequestsMock.mainExpectation != nil {

		input := m.PendingRequestsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCalculatorMockPendingRequestsInput{p, p1, p2}, "FilamentCalculator.PendingRequests got unexpected parameters")
		}

		result := m.PendingRequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.PendingRequests")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PendingRequestsFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCalculatorMock.PendingRequests. %v %v %v", p, p1, p2)
		return
	}

	return m.PendingRequestsFunc(p, p1, p2)
}

//PendingRequestsMinimockCounter returns a count of FilamentCalculatorMock.PendingRequestsFunc invocations
func (m *FilamentCalculatorMock) PendingRequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PendingRequestsCounter)
}

//PendingRequestsMinimockPreCounter returns the value of FilamentCalculatorMock.PendingRequests invocations
func (m *FilamentCalculatorMock) PendingRequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PendingRequestsPreCounter)
}

//PendingRequestsFinished returns true if mock invocations count is ok
func (m *FilamentCalculatorMock) PendingRequestsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PendingRequestsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PendingRequestsCounter) == uint64(len(m.PendingRequestsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PendingRequestsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PendingRequestsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PendingRequestsFunc != nil {
		return atomic.LoadUint64(&m.PendingRequestsCounter) > 0
	}

	return true
}

type mFilamentCalculatorMockRequests struct {
	mock              *FilamentCalculatorMock
	mainExpectation   *FilamentCalculatorMockRequestsExpectation
	expectationSeries []*FilamentCalculatorMockRequestsExpectation
}

type FilamentCalculatorMockRequestsExpectation struct {
	input  *FilamentCalculatorMockRequestsInput
	result *FilamentCalculatorMockRequestsResult
}

type FilamentCalculatorMockRequestsInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.ID
	p3 insolar.PulseNumber
	p4 insolar.PulseNumber
}

type FilamentCalculatorMockRequestsResult struct {
	r  []record.CompositeFilamentRecord
	r1 error
}

//Expect specifies that invocation of FilamentCalculator.Requests is expected from 1 to Infinity times
func (m *mFilamentCalculatorMockRequests) Expect(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.PulseNumber) *mFilamentCalculatorMockRequests {
	m.mock.RequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockRequestsExpectation{}
	}
	m.mainExpectation.input = &FilamentCalculatorMockRequestsInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of FilamentCalculator.Requests
func (m *mFilamentCalculatorMockRequests) Return(r []record.CompositeFilamentRecord, r1 error) *FilamentCalculatorMock {
	m.mock.RequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockRequestsExpectation{}
	}
	m.mainExpectation.result = &FilamentCalculatorMockRequestsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCalculator.Requests is expected once
func (m *mFilamentCalculatorMockRequests) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.PulseNumber) *FilamentCalculatorMockRequestsExpectation {
	m.mock.RequestsFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCalculatorMockRequestsExpectation{}
	expectation.input = &FilamentCalculatorMockRequestsInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCalculatorMockRequestsExpectation) Return(r []record.CompositeFilamentRecord, r1 error) {
	e.result = &FilamentCalculatorMockRequestsResult{r, r1}
}

//Set uses given function f as a mock of FilamentCalculator.Requests method
func (m *mFilamentCalculatorMockRequests) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.PulseNumber) (r []record.CompositeFilamentRecord, r1 error)) *FilamentCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RequestsFunc = f
	return m.mock
}

//Requests implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator interface
func (m *FilamentCalculatorMock) Requests(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.PulseNumber) (r []record.CompositeFilamentRecord, r1 error) {
	counter := atomic.AddUint64(&m.RequestsPreCounter, 1)
	defer atomic.AddUint64(&m.RequestsCounter, 1)

	if len(m.RequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCalculatorMock.Requests. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.RequestsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCalculatorMockRequestsInput{p, p1, p2, p3, p4}, "FilamentCalculator.Requests got unexpected parameters")

		result := m.RequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.Requests")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RequestsMock.mainExpectation != nil {

		input := m.RequestsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCalculatorMockRequestsInput{p, p1, p2, p3, p4}, "FilamentCalculator.Requests got unexpected parameters")
		}

		result := m.RequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.Requests")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RequestsFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCalculatorMock.Requests. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.RequestsFunc(p, p1, p2, p3, p4)
}

//RequestsMinimockCounter returns a count of FilamentCalculatorMock.RequestsFunc invocations
func (m *FilamentCalculatorMock) RequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RequestsCounter)
}

//RequestsMinimockPreCounter returns the value of FilamentCalculatorMock.Requests invocations
func (m *FilamentCalculatorMock) RequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RequestsPreCounter)
}

//RequestsFinished returns true if mock invocations count is ok
func (m *FilamentCalculatorMock) RequestsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RequestsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RequestsCounter) == uint64(len(m.RequestsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RequestsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RequestsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RequestsFunc != nil {
		return atomic.LoadUint64(&m.RequestsCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCalculatorMock) ValidateCallCounters() {

	if !m.PendingRequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.PendingRequests")
	}

	if !m.RequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.Requests")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCalculatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FilamentCalculatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FilamentCalculatorMock) MinimockFinish() {

	if !m.PendingRequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.PendingRequests")
	}

	if !m.RequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.Requests")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FilamentCalculatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FilamentCalculatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.PendingRequestsFinished()
		ok = ok && m.RequestsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.PendingRequestsFinished() {
				m.t.Error("Expected call to FilamentCalculatorMock.PendingRequests")
			}

			if !m.RequestsFinished() {
				m.t.Error("Expected call to FilamentCalculatorMock.Requests")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *FilamentCalculatorMock) AllMocksCalled() bool {

	if !m.PendingRequestsFinished() {
		return false
	}

	if !m.RequestsFinished() {
		return false
	}

	return true
}

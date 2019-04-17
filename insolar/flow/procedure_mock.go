package flow

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Procedure" can be found in github.com/insolar/insolar/insolar/flow
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProcedureMock implements github.com/insolar/insolar/insolar/flow.Procedure
type ProcedureMock struct {
	t minimock.Tester

	ProceedFunc       func(p context.Context) (r error)
	ProceedCounter    uint64
	ProceedPreCounter uint64
	ProceedMock       mProcedureMockProceed
}

//NewProcedureMock returns a mock for github.com/insolar/insolar/insolar/flow.Procedure
func NewProcedureMock(t minimock.Tester) *ProcedureMock {
	m := &ProcedureMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ProceedMock = mProcedureMockProceed{mock: m}

	return m
}

type mProcedureMockProceed struct {
	mock              *ProcedureMock
	mainExpectation   *ProcedureMockProceedExpectation
	expectationSeries []*ProcedureMockProceedExpectation
}

type ProcedureMockProceedExpectation struct {
	input  *ProcedureMockProceedInput
	result *ProcedureMockProceedResult
}

type ProcedureMockProceedInput struct {
	p context.Context
}

type ProcedureMockProceedResult struct {
	r error
}

//Expect specifies that invocation of Procedure.Proceed is expected from 1 to Infinity times
func (m *mProcedureMockProceed) Expect(p context.Context) *mProcedureMockProceed {
	m.mock.ProceedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProcedureMockProceedExpectation{}
	}
	m.mainExpectation.input = &ProcedureMockProceedInput{p}
	return m
}

//Return specifies results of invocation of Procedure.Proceed
func (m *mProcedureMockProceed) Return(r error) *ProcedureMock {
	m.mock.ProceedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProcedureMockProceedExpectation{}
	}
	m.mainExpectation.result = &ProcedureMockProceedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Procedure.Proceed is expected once
func (m *mProcedureMockProceed) ExpectOnce(p context.Context) *ProcedureMockProceedExpectation {
	m.mock.ProceedFunc = nil
	m.mainExpectation = nil

	expectation := &ProcedureMockProceedExpectation{}
	expectation.input = &ProcedureMockProceedInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProcedureMockProceedExpectation) Return(r error) {
	e.result = &ProcedureMockProceedResult{r}
}

//Set uses given function f as a mock of Procedure.Proceed method
func (m *mProcedureMockProceed) Set(f func(p context.Context) (r error)) *ProcedureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ProceedFunc = f
	return m.mock
}

//Proceed implements github.com/insolar/insolar/insolar/flow.Procedure interface
func (m *ProcedureMock) Proceed(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.ProceedPreCounter, 1)
	defer atomic.AddUint64(&m.ProceedCounter, 1)

	if len(m.ProceedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ProceedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProcedureMock.Proceed. %v", p)
			return
		}

		input := m.ProceedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProcedureMockProceedInput{p}, "Procedure.Proceed got unexpected parameters")

		result := m.ProceedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProcedureMock.Proceed")
			return
		}

		r = result.r

		return
	}

	if m.ProceedMock.mainExpectation != nil {

		input := m.ProceedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProcedureMockProceedInput{p}, "Procedure.Proceed got unexpected parameters")
		}

		result := m.ProceedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProcedureMock.Proceed")
		}

		r = result.r

		return
	}

	if m.ProceedFunc == nil {
		m.t.Fatalf("Unexpected call to ProcedureMock.Proceed. %v", p)
		return
	}

	return m.ProceedFunc(p)
}

//ProceedMinimockCounter returns a count of ProcedureMock.ProceedFunc invocations
func (m *ProcedureMock) ProceedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ProceedCounter)
}

//ProceedMinimockPreCounter returns the value of ProcedureMock.Proceed invocations
func (m *ProcedureMock) ProceedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ProceedPreCounter)
}

//ProceedFinished returns true if mock invocations count is ok
func (m *ProcedureMock) ProceedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ProceedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ProceedCounter) == uint64(len(m.ProceedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ProceedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ProceedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ProceedFunc != nil {
		return atomic.LoadUint64(&m.ProceedCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProcedureMock) ValidateCallCounters() {

	if !m.ProceedFinished() {
		m.t.Fatal("Expected call to ProcedureMock.Proceed")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProcedureMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ProcedureMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ProcedureMock) MinimockFinish() {

	if !m.ProceedFinished() {
		m.t.Fatal("Expected call to ProcedureMock.Proceed")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ProcedureMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ProcedureMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ProceedFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ProceedFinished() {
				m.t.Error("Expected call to ProcedureMock.Proceed")
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
func (m *ProcedureMock) AllMocksCalled() bool {

	if !m.ProceedFinished() {
		return false
	}

	return true
}

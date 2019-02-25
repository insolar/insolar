package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "TerminationHandler" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//TerminationHandlerMock implements github.com/insolar/insolar/core.TerminationHandler
type TerminationHandlerMock struct {
	t minimock.Tester

	AbortFunc       func(p string)
	AbortCounter    uint64
	AbortPreCounter uint64
	AbortMock       mTerminationHandlerMockAbort
}

//NewTerminationHandlerMock returns a mock for github.com/insolar/insolar/core.TerminationHandler
func NewTerminationHandlerMock(t minimock.Tester) *TerminationHandlerMock {
	m := &TerminationHandlerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AbortMock = mTerminationHandlerMockAbort{mock: m}

	return m
}

type mTerminationHandlerMockAbort struct {
	mock              *TerminationHandlerMock
	mainExpectation   *TerminationHandlerMockAbortExpectation
	expectationSeries []*TerminationHandlerMockAbortExpectation
}

type TerminationHandlerMockAbortExpectation struct {
	input *TerminationHandlerMockAbortInput
}

type TerminationHandlerMockAbortInput struct {
	p string
}

//Expect specifies that invocation of TerminationHandler.Abort is expected from 1 to Infinity times
func (m *mTerminationHandlerMockAbort) Expect(p string) *mTerminationHandlerMockAbort {
	m.mock.AbortFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockAbortExpectation{}
	}
	m.mainExpectation.input = &TerminationHandlerMockAbortInput{p}
	return m
}

//Return specifies results of invocation of TerminationHandler.Abort
func (m *mTerminationHandlerMockAbort) Return() *TerminationHandlerMock {
	m.mock.AbortFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockAbortExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TerminationHandler.Abort is expected once
func (m *mTerminationHandlerMockAbort) ExpectOnce(p string) *TerminationHandlerMockAbortExpectation {
	m.mock.AbortFunc = nil
	m.mainExpectation = nil

	expectation := &TerminationHandlerMockAbortExpectation{}
	expectation.input = &TerminationHandlerMockAbortInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TerminationHandler.Abort method
func (m *mTerminationHandlerMockAbort) Set(f func(p string)) *TerminationHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AbortFunc = f
	return m.mock
}

//Abort implements github.com/insolar/insolar/core.TerminationHandler interface
func (m *TerminationHandlerMock) Abort(p string) {
	counter := atomic.AddUint64(&m.AbortPreCounter, 1)
	defer atomic.AddUint64(&m.AbortCounter, 1)

	if len(m.AbortMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AbortMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TerminationHandlerMock.Abort. %v", p)
			return
		}

		input := m.AbortMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TerminationHandlerMockAbortInput{p}, "TerminationHandler.Abort got unexpected parameters")

		return
	}

	if m.AbortMock.mainExpectation != nil {

		input := m.AbortMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TerminationHandlerMockAbortInput{p}, "TerminationHandler.Abort got unexpected parameters")
		}

		return
	}

	if m.AbortFunc == nil {
		m.t.Fatalf("Unexpected call to TerminationHandlerMock.Abort. %v", p)
		return
	}

	m.AbortFunc(p)
}

//AbortMinimockCounter returns a count of TerminationHandlerMock.AbortFunc invocations
func (m *TerminationHandlerMock) AbortMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AbortCounter)
}

//AbortMinimockPreCounter returns the value of TerminationHandlerMock.Abort invocations
func (m *TerminationHandlerMock) AbortMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AbortPreCounter)
}

//AbortFinished returns true if mock invocations count is ok
func (m *TerminationHandlerMock) AbortFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AbortMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AbortCounter) == uint64(len(m.AbortMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AbortMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AbortCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AbortFunc != nil {
		return atomic.LoadUint64(&m.AbortCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TerminationHandlerMock) ValidateCallCounters() {

	if !m.AbortFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.Abort")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TerminationHandlerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TerminationHandlerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TerminationHandlerMock) MinimockFinish() {

	if !m.AbortFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.Abort")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TerminationHandlerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TerminationHandlerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AbortFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AbortFinished() {
				m.t.Error("Expected call to TerminationHandlerMock.Abort")
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
func (m *TerminationHandlerMock) AllMocksCalled() bool {

	if !m.AbortFinished() {
		return false
	}

	return true
}

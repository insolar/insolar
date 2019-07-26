package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StateIniter" can be found in github.com/insolar/insolar/ledger/light/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//StateIniterMock implements github.com/insolar/insolar/ledger/light/executor.StateIniter
type StateIniterMock struct {
	t minimock.Tester

	PrepareStateFunc       func(p context.Context, p1 insolar.PulseNumber) (r error)
	PrepareStateCounter    uint64
	PrepareStatePreCounter uint64
	PrepareStateMock       mStateIniterMockPrepareState
}

//NewStateIniterMock returns a mock for github.com/insolar/insolar/ledger/light/executor.StateIniter
func NewStateIniterMock(t minimock.Tester) *StateIniterMock {
	m := &StateIniterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PrepareStateMock = mStateIniterMockPrepareState{mock: m}

	return m
}

type mStateIniterMockPrepareState struct {
	mock              *StateIniterMock
	mainExpectation   *StateIniterMockPrepareStateExpectation
	expectationSeries []*StateIniterMockPrepareStateExpectation
}

type StateIniterMockPrepareStateExpectation struct {
	input  *StateIniterMockPrepareStateInput
	result *StateIniterMockPrepareStateResult
}

type StateIniterMockPrepareStateInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type StateIniterMockPrepareStateResult struct {
	r error
}

//Expect specifies that invocation of StateIniter.PrepareState is expected from 1 to Infinity times
func (m *mStateIniterMockPrepareState) Expect(p context.Context, p1 insolar.PulseNumber) *mStateIniterMockPrepareState {
	m.mock.PrepareStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateIniterMockPrepareStateExpectation{}
	}
	m.mainExpectation.input = &StateIniterMockPrepareStateInput{p, p1}
	return m
}

//Return specifies results of invocation of StateIniter.PrepareState
func (m *mStateIniterMockPrepareState) Return(r error) *StateIniterMock {
	m.mock.PrepareStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateIniterMockPrepareStateExpectation{}
	}
	m.mainExpectation.result = &StateIniterMockPrepareStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateIniter.PrepareState is expected once
func (m *mStateIniterMockPrepareState) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *StateIniterMockPrepareStateExpectation {
	m.mock.PrepareStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateIniterMockPrepareStateExpectation{}
	expectation.input = &StateIniterMockPrepareStateInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateIniterMockPrepareStateExpectation) Return(r error) {
	e.result = &StateIniterMockPrepareStateResult{r}
}

//Set uses given function f as a mock of StateIniter.PrepareState method
func (m *mStateIniterMockPrepareState) Set(f func(p context.Context, p1 insolar.PulseNumber) (r error)) *StateIniterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PrepareStateFunc = f
	return m.mock
}

//PrepareState implements github.com/insolar/insolar/ledger/light/executor.StateIniter interface
func (m *StateIniterMock) PrepareState(p context.Context, p1 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.PrepareStatePreCounter, 1)
	defer atomic.AddUint64(&m.PrepareStateCounter, 1)

	if len(m.PrepareStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PrepareStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateIniterMock.PrepareState. %v %v", p, p1)
			return
		}

		input := m.PrepareStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateIniterMockPrepareStateInput{p, p1}, "StateIniter.PrepareState got unexpected parameters")

		result := m.PrepareStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateIniterMock.PrepareState")
			return
		}

		r = result.r

		return
	}

	if m.PrepareStateMock.mainExpectation != nil {

		input := m.PrepareStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateIniterMockPrepareStateInput{p, p1}, "StateIniter.PrepareState got unexpected parameters")
		}

		result := m.PrepareStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateIniterMock.PrepareState")
		}

		r = result.r

		return
	}

	if m.PrepareStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateIniterMock.PrepareState. %v %v", p, p1)
		return
	}

	return m.PrepareStateFunc(p, p1)
}

//PrepareStateMinimockCounter returns a count of StateIniterMock.PrepareStateFunc invocations
func (m *StateIniterMock) PrepareStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PrepareStateCounter)
}

//PrepareStateMinimockPreCounter returns the value of StateIniterMock.PrepareState invocations
func (m *StateIniterMock) PrepareStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PrepareStatePreCounter)
}

//PrepareStateFinished returns true if mock invocations count is ok
func (m *StateIniterMock) PrepareStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PrepareStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PrepareStateCounter) == uint64(len(m.PrepareStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PrepareStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PrepareStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PrepareStateFunc != nil {
		return atomic.LoadUint64(&m.PrepareStateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateIniterMock) ValidateCallCounters() {

	if !m.PrepareStateFinished() {
		m.t.Fatal("Expected call to StateIniterMock.PrepareState")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateIniterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StateIniterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StateIniterMock) MinimockFinish() {

	if !m.PrepareStateFinished() {
		m.t.Fatal("Expected call to StateIniterMock.PrepareState")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StateIniterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StateIniterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.PrepareStateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.PrepareStateFinished() {
				m.t.Error("Expected call to StateIniterMock.PrepareState")
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
func (m *StateIniterMock) AllMocksCalled() bool {

	if !m.PrepareStateFinished() {
		return false
	}

	return true
}

package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseHandler" can be found in github.com/insolar/insolar/network
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseHandlerMock implements github.com/insolar/insolar/network.PulseHandler
type PulseHandlerMock struct {
	t minimock.Tester

	HandlePulseFunc       func(p context.Context, p1 core.Pulse)
	HandlePulseCounter    uint64
	HandlePulsePreCounter uint64
	HandlePulseMock       mPulseHandlerMockHandlePulse
}

//NewPulseHandlerMock returns a mock for github.com/insolar/insolar/network.PulseHandler
func NewPulseHandlerMock(t minimock.Tester) *PulseHandlerMock {
	m := &PulseHandlerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.HandlePulseMock = mPulseHandlerMockHandlePulse{mock: m}

	return m
}

type mPulseHandlerMockHandlePulse struct {
	mock              *PulseHandlerMock
	mainExpectation   *PulseHandlerMockHandlePulseExpectation
	expectationSeries []*PulseHandlerMockHandlePulseExpectation
}

type PulseHandlerMockHandlePulseExpectation struct {
	input *PulseHandlerMockHandlePulseInput
}

type PulseHandlerMockHandlePulseInput struct {
	p  context.Context
	p1 core.Pulse
}

//Expect specifies that invocation of PulseHandler.HandlePulse is expected from 1 to Infinity times
func (m *mPulseHandlerMockHandlePulse) Expect(p context.Context, p1 core.Pulse) *mPulseHandlerMockHandlePulse {
	m.mock.HandlePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseHandlerMockHandlePulseExpectation{}
	}
	m.mainExpectation.input = &PulseHandlerMockHandlePulseInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseHandler.HandlePulse
func (m *mPulseHandlerMockHandlePulse) Return() *PulseHandlerMock {
	m.mock.HandlePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseHandlerMockHandlePulseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PulseHandler.HandlePulse is expected once
func (m *mPulseHandlerMockHandlePulse) ExpectOnce(p context.Context, p1 core.Pulse) *PulseHandlerMockHandlePulseExpectation {
	m.mock.HandlePulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulseHandlerMockHandlePulseExpectation{}
	expectation.input = &PulseHandlerMockHandlePulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PulseHandler.HandlePulse method
func (m *mPulseHandlerMockHandlePulse) Set(f func(p context.Context, p1 core.Pulse)) *PulseHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HandlePulseFunc = f
	return m.mock
}

//HandlePulse implements github.com/insolar/insolar/network.PulseHandler interface
func (m *PulseHandlerMock) HandlePulse(p context.Context, p1 core.Pulse) {
	counter := atomic.AddUint64(&m.HandlePulsePreCounter, 1)
	defer atomic.AddUint64(&m.HandlePulseCounter, 1)

	if len(m.HandlePulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HandlePulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseHandlerMock.HandlePulse. %v %v", p, p1)
			return
		}

		input := m.HandlePulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseHandlerMockHandlePulseInput{p, p1}, "PulseHandler.HandlePulse got unexpected parameters")

		return
	}

	if m.HandlePulseMock.mainExpectation != nil {

		input := m.HandlePulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseHandlerMockHandlePulseInput{p, p1}, "PulseHandler.HandlePulse got unexpected parameters")
		}

		return
	}

	if m.HandlePulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulseHandlerMock.HandlePulse. %v %v", p, p1)
		return
	}

	m.HandlePulseFunc(p, p1)
}

//HandlePulseMinimockCounter returns a count of PulseHandlerMock.HandlePulseFunc invocations
func (m *PulseHandlerMock) HandlePulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HandlePulseCounter)
}

//HandlePulseMinimockPreCounter returns the value of PulseHandlerMock.HandlePulse invocations
func (m *PulseHandlerMock) HandlePulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HandlePulsePreCounter)
}

//HandlePulseFinished returns true if mock invocations count is ok
func (m *PulseHandlerMock) HandlePulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HandlePulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HandlePulseCounter) == uint64(len(m.HandlePulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HandlePulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HandlePulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HandlePulseFunc != nil {
		return atomic.LoadUint64(&m.HandlePulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseHandlerMock) ValidateCallCounters() {

	if !m.HandlePulseFinished() {
		m.t.Fatal("Expected call to PulseHandlerMock.HandlePulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseHandlerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseHandlerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseHandlerMock) MinimockFinish() {

	if !m.HandlePulseFinished() {
		m.t.Fatal("Expected call to PulseHandlerMock.HandlePulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseHandlerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseHandlerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.HandlePulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.HandlePulseFinished() {
				m.t.Error("Expected call to PulseHandlerMock.HandlePulse")
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
func (m *PulseHandlerMock) AllMocksCalled() bool {

	if !m.HandlePulseFinished() {
		return false
	}

	return true
}

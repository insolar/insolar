package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseHandler" can be found in github.com/insolar/insolar/network
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"
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
	mock             *PulseHandlerMock
	mockExpectations *PulseHandlerMockHandlePulseParams
}

//PulseHandlerMockHandlePulseParams represents input parameters of the PulseHandler.HandlePulse
type PulseHandlerMockHandlePulseParams struct {
	p  context.Context
	p1 core.Pulse
}

//Expect sets up expected params for the PulseHandler.HandlePulse
func (m *mPulseHandlerMockHandlePulse) Expect(p context.Context, p1 core.Pulse) *mPulseHandlerMockHandlePulse {
	m.mockExpectations = &PulseHandlerMockHandlePulseParams{p, p1}
	return m
}

//Return sets up a mock for PulseHandler.HandlePulse to return Return's arguments
func (m *mPulseHandlerMockHandlePulse) Return() *PulseHandlerMock {
	m.mock.HandlePulseFunc = func(p context.Context, p1 core.Pulse) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of PulseHandler.HandlePulse method
func (m *mPulseHandlerMockHandlePulse) Set(f func(p context.Context, p1 core.Pulse)) *PulseHandlerMock {
	m.mock.HandlePulseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//HandlePulse implements github.com/insolar/insolar/network.PulseHandler interface
func (m *PulseHandlerMock) HandlePulse(p context.Context, p1 core.Pulse) {
	atomic.AddUint64(&m.HandlePulsePreCounter, 1)
	defer atomic.AddUint64(&m.HandlePulseCounter, 1)

	if m.HandlePulseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.HandlePulseMock.mockExpectations, PulseHandlerMockHandlePulseParams{p, p1},
			"PulseHandler.HandlePulse got unexpected parameters")

		if m.HandlePulseFunc == nil {

			m.t.Fatal("No results are set for the PulseHandlerMock.HandlePulse")

			return
		}
	}

	if m.HandlePulseFunc == nil {
		m.t.Fatal("Unexpected call to PulseHandlerMock.HandlePulse")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseHandlerMock) ValidateCallCounters() {

	if m.HandlePulseFunc != nil && atomic.LoadUint64(&m.HandlePulseCounter) == 0 {
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

	if m.HandlePulseFunc != nil && atomic.LoadUint64(&m.HandlePulseCounter) == 0 {
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
		ok = ok && (m.HandlePulseFunc == nil || atomic.LoadUint64(&m.HandlePulseCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.HandlePulseFunc != nil && atomic.LoadUint64(&m.HandlePulseCounter) == 0 {
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

	if m.HandlePulseFunc != nil && atomic.LoadUint64(&m.HandlePulseCounter) == 0 {
		return false
	}

	return true
}

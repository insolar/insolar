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

	OnPulseFunc       func(p context.Context, p1 core.Pulse)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mPulseHandlerMockOnPulse
}

//NewPulseHandlerMock returns a mock for github.com/insolar/insolar/network.PulseHandler
func NewPulseHandlerMock(t minimock.Tester) *PulseHandlerMock {
	m := &PulseHandlerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.OnPulseMock = mPulseHandlerMockOnPulse{mock: m}

	return m
}

type mPulseHandlerMockOnPulse struct {
	mock             *PulseHandlerMock
	mockExpectations *PulseHandlerMockOnPulseParams
}

//PulseHandlerMockOnPulseParams represents input parameters of the PulseHandler.OnPulse
type PulseHandlerMockOnPulseParams struct {
	p  context.Context
	p1 core.Pulse
}

//Expect sets up expected params for the PulseHandler.OnPulse
func (m *mPulseHandlerMockOnPulse) Expect(p context.Context, p1 core.Pulse) *mPulseHandlerMockOnPulse {
	m.mockExpectations = &PulseHandlerMockOnPulseParams{p, p1}
	return m
}

//Return sets up a mock for PulseHandler.OnPulse to return Return's arguments
func (m *mPulseHandlerMockOnPulse) Return() *PulseHandlerMock {
	m.mock.OnPulseFunc = func(p context.Context, p1 core.Pulse) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of PulseHandler.OnPulse method
func (m *mPulseHandlerMockOnPulse) Set(f func(p context.Context, p1 core.Pulse)) *PulseHandlerMock {
	m.mock.OnPulseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/network.PulseHandler interface
func (m *PulseHandlerMock) OnPulse(p context.Context, p1 core.Pulse) {
	atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if m.OnPulseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.OnPulseMock.mockExpectations, PulseHandlerMockOnPulseParams{p, p1},
			"PulseHandler.OnPulse got unexpected parameters")

		if m.OnPulseFunc == nil {

			m.t.Fatal("No results are set for the PulseHandlerMock.OnPulse")

			return
		}
	}

	if m.OnPulseFunc == nil {
		m.t.Fatal("Unexpected call to PulseHandlerMock.OnPulse")
		return
	}

	m.OnPulseFunc(p, p1)
}

//OnPulseMinimockCounter returns a count of PulseHandlerMock.OnPulseFunc invocations
func (m *PulseHandlerMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of PulseHandlerMock.OnPulse invocations
func (m *PulseHandlerMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseHandlerMock) ValidateCallCounters() {

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		m.t.Fatal("Expected call to PulseHandlerMock.OnPulse")
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

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		m.t.Fatal("Expected call to PulseHandlerMock.OnPulse")
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
		ok = ok && (m.OnPulseFunc == nil || atomic.LoadUint64(&m.OnPulseCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
				m.t.Error("Expected call to PulseHandlerMock.OnPulse")
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

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		return false
	}

	return true
}

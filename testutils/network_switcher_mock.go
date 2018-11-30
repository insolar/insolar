package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NetworkSwitcher" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//NetworkSwitcherMock implements github.com/insolar/insolar/core.NetworkSwitcher
type NetworkSwitcherMock struct {
	t minimock.Tester

	GetStateFunc       func() (r core.NetworkState)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mNetworkSwitcherMockGetState

	OnPulseFunc       func(p context.Context, p1 core.Pulse) (r error)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mNetworkSwitcherMockOnPulse
}

//NewNetworkSwitcherMock returns a mock for github.com/insolar/insolar/core.NetworkSwitcher
func NewNetworkSwitcherMock(t minimock.Tester) *NetworkSwitcherMock {
	m := &NetworkSwitcherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetStateMock = mNetworkSwitcherMockGetState{mock: m}
	m.OnPulseMock = mNetworkSwitcherMockOnPulse{mock: m}

	return m
}

type mNetworkSwitcherMockGetState struct {
	mock *NetworkSwitcherMock
}

//Return sets up a mock for NetworkSwitcher.GetState to return Return's arguments
func (m *mNetworkSwitcherMockGetState) Return(r core.NetworkState) *NetworkSwitcherMock {
	m.mock.GetStateFunc = func() core.NetworkState {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NetworkSwitcher.GetState method
func (m *mNetworkSwitcherMockGetState) Set(f func() (r core.NetworkState)) *NetworkSwitcherMock {
	m.mock.GetStateFunc = f

	return m.mock
}

//GetState implements github.com/insolar/insolar/core.NetworkSwitcher interface
func (m *NetworkSwitcherMock) GetState() (r core.NetworkState) {
	atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if m.GetStateFunc == nil {
		m.t.Fatal("Unexpected call to NetworkSwitcherMock.GetState")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of NetworkSwitcherMock.GetStateFunc invocations
func (m *NetworkSwitcherMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of NetworkSwitcherMock.GetState invocations
func (m *NetworkSwitcherMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

type mNetworkSwitcherMockOnPulse struct {
	mock             *NetworkSwitcherMock
	mockExpectations *NetworkSwitcherMockOnPulseParams
}

//NetworkSwitcherMockOnPulseParams represents input parameters of the NetworkSwitcher.OnPulse
type NetworkSwitcherMockOnPulseParams struct {
	p  context.Context
	p1 core.Pulse
}

//Expect sets up expected params for the NetworkSwitcher.OnPulse
func (m *mNetworkSwitcherMockOnPulse) Expect(p context.Context, p1 core.Pulse) *mNetworkSwitcherMockOnPulse {
	m.mockExpectations = &NetworkSwitcherMockOnPulseParams{p, p1}
	return m
}

//Return sets up a mock for NetworkSwitcher.OnPulse to return Return's arguments
func (m *mNetworkSwitcherMockOnPulse) Return(r error) *NetworkSwitcherMock {
	m.mock.OnPulseFunc = func(p context.Context, p1 core.Pulse) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NetworkSwitcher.OnPulse method
func (m *mNetworkSwitcherMockOnPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *NetworkSwitcherMock {
	m.mock.OnPulseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/core.NetworkSwitcher interface
func (m *NetworkSwitcherMock) OnPulse(p context.Context, p1 core.Pulse) (r error) {
	atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if m.OnPulseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.OnPulseMock.mockExpectations, NetworkSwitcherMockOnPulseParams{p, p1},
			"NetworkSwitcher.OnPulse got unexpected parameters")

		if m.OnPulseFunc == nil {

			m.t.Fatal("No results are set for the NetworkSwitcherMock.OnPulse")

			return
		}
	}

	if m.OnPulseFunc == nil {
		m.t.Fatal("Unexpected call to NetworkSwitcherMock.OnPulse")
		return
	}

	return m.OnPulseFunc(p, p1)
}

//OnPulseMinimockCounter returns a count of NetworkSwitcherMock.OnPulseFunc invocations
func (m *NetworkSwitcherMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of NetworkSwitcherMock.OnPulse invocations
func (m *NetworkSwitcherMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkSwitcherMock) ValidateCallCounters() {

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		m.t.Fatal("Expected call to NetworkSwitcherMock.GetState")
	}

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		m.t.Fatal("Expected call to NetworkSwitcherMock.OnPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkSwitcherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NetworkSwitcherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NetworkSwitcherMock) MinimockFinish() {

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		m.t.Fatal("Expected call to NetworkSwitcherMock.GetState")
	}

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		m.t.Fatal("Expected call to NetworkSwitcherMock.OnPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NetworkSwitcherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NetworkSwitcherMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetStateFunc == nil || atomic.LoadUint64(&m.GetStateCounter) > 0)
		ok = ok && (m.OnPulseFunc == nil || atomic.LoadUint64(&m.OnPulseCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
				m.t.Error("Expected call to NetworkSwitcherMock.GetState")
			}

			if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
				m.t.Error("Expected call to NetworkSwitcherMock.OnPulse")
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
func (m *NetworkSwitcherMock) AllMocksCalled() bool {

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		return false
	}

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		return false
	}

	return true
}

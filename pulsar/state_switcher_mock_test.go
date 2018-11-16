package pulsar

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StateSwitcher" can be found in github.com/insolar/insolar/pulsar
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//StateSwitcherMock implements github.com/insolar/insolar/pulsar.StateSwitcher
type StateSwitcherMock struct {
	t minimock.Tester

	GetStateFunc       func() (r State)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mStateSwitcherMockGetState

	SetPulsarFunc       func(p *Pulsar)
	SetPulsarCounter    uint64
	SetPulsarPreCounter uint64
	SetPulsarMock       mStateSwitcherMockSetPulsar

	SwitchToStateFunc       func(p context.Context, p1 State, p2 interface{})
	SwitchToStateCounter    uint64
	SwitchToStatePreCounter uint64
	SwitchToStateMock       mStateSwitcherMockSwitchToState

	setStateFunc       func(p State)
	setStateCounter    uint64
	setStatePreCounter uint64
	setStateMock       mStateSwitcherMocksetState
}

//NewStateSwitcherMock returns a mock for github.com/insolar/insolar/pulsar.StateSwitcher
func NewStateSwitcherMock(t minimock.Tester) *StateSwitcherMock {
	m := &StateSwitcherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetStateMock = mStateSwitcherMockGetState{mock: m}
	m.SetPulsarMock = mStateSwitcherMockSetPulsar{mock: m}
	m.SwitchToStateMock = mStateSwitcherMockSwitchToState{mock: m}
	m.setStateMock = mStateSwitcherMocksetState{mock: m}

	return m
}

type mStateSwitcherMockGetState struct {
	mock *StateSwitcherMock
}

//Return sets up a mock for StateSwitcher.GetState to return Return's arguments
func (m *mStateSwitcherMockGetState) Return(r State) *StateSwitcherMock {
	m.mock.GetStateFunc = func() State {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of StateSwitcher.GetState method
func (m *mStateSwitcherMockGetState) Set(f func() (r State)) *StateSwitcherMock {
	m.mock.GetStateFunc = f

	return m.mock
}

//GetState implements github.com/insolar/insolar/pulsar.StateSwitcher interface
func (m *StateSwitcherMock) GetState() (r State) {
	atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if m.GetStateFunc == nil {
		m.t.Fatal("Unexpected call to StateSwitcherMock.GetState")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of StateSwitcherMock.GetStateFunc invocations
func (m *StateSwitcherMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of StateSwitcherMock.GetState invocations
func (m *StateSwitcherMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

type mStateSwitcherMockSetPulsar struct {
	mock             *StateSwitcherMock
	mockExpectations *StateSwitcherMockSetPulsarParams
}

//StateSwitcherMockSetPulsarParams represents input parameters of the StateSwitcher.SetPulsar
type StateSwitcherMockSetPulsarParams struct {
	p *Pulsar
}

//Expect sets up expected params for the StateSwitcher.SetPulsar
func (m *mStateSwitcherMockSetPulsar) Expect(p *Pulsar) *mStateSwitcherMockSetPulsar {
	m.mockExpectations = &StateSwitcherMockSetPulsarParams{p}
	return m
}

//Return sets up a mock for StateSwitcher.SetPulsar to return Return's arguments
func (m *mStateSwitcherMockSetPulsar) Return() *StateSwitcherMock {
	m.mock.SetPulsarFunc = func(p *Pulsar) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of StateSwitcher.SetPulsar method
func (m *mStateSwitcherMockSetPulsar) Set(f func(p *Pulsar)) *StateSwitcherMock {
	m.mock.SetPulsarFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetPulsar implements github.com/insolar/insolar/pulsar.StateSwitcher interface
func (m *StateSwitcherMock) SetPulsar(p *Pulsar) {
	atomic.AddUint64(&m.SetPulsarPreCounter, 1)
	defer atomic.AddUint64(&m.SetPulsarCounter, 1)

	if m.SetPulsarMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetPulsarMock.mockExpectations, StateSwitcherMockSetPulsarParams{p},
			"StateSwitcher.SetPulsar got unexpected parameters")

		if m.SetPulsarFunc == nil {

			m.t.Fatal("No results are set for the StateSwitcherMock.SetPulsar")

			return
		}
	}

	if m.SetPulsarFunc == nil {
		m.t.Fatal("Unexpected call to StateSwitcherMock.SetPulsar")
		return
	}

	m.SetPulsarFunc(p)
}

//SetPulsarMinimockCounter returns a count of StateSwitcherMock.SetPulsarFunc invocations
func (m *StateSwitcherMock) SetPulsarMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetPulsarCounter)
}

//SetPulsarMinimockPreCounter returns the value of StateSwitcherMock.SetPulsar invocations
func (m *StateSwitcherMock) SetPulsarMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPulsarPreCounter)
}

type mStateSwitcherMockSwitchToState struct {
	mock             *StateSwitcherMock
	mockExpectations *StateSwitcherMockSwitchToStateParams
}

//StateSwitcherMockSwitchToStateParams represents input parameters of the StateSwitcher.SwitchToState
type StateSwitcherMockSwitchToStateParams struct {
	p  context.Context
	p1 State
	p2 interface{}
}

//Expect sets up expected params for the StateSwitcher.SwitchToState
func (m *mStateSwitcherMockSwitchToState) Expect(p context.Context, p1 State, p2 interface{}) *mStateSwitcherMockSwitchToState {
	m.mockExpectations = &StateSwitcherMockSwitchToStateParams{p, p1, p2}
	return m
}

//Return sets up a mock for StateSwitcher.SwitchToState to return Return's arguments
func (m *mStateSwitcherMockSwitchToState) Return() *StateSwitcherMock {
	m.mock.SwitchToStateFunc = func(p context.Context, p1 State, p2 interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of StateSwitcher.SwitchToState method
func (m *mStateSwitcherMockSwitchToState) Set(f func(p context.Context, p1 State, p2 interface{})) *StateSwitcherMock {
	m.mock.SwitchToStateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SwitchToState implements github.com/insolar/insolar/pulsar.StateSwitcher interface
func (m *StateSwitcherMock) SwitchToState(p context.Context, p1 State, p2 interface{}) {
	atomic.AddUint64(&m.SwitchToStatePreCounter, 1)
	defer atomic.AddUint64(&m.SwitchToStateCounter, 1)

	if m.SwitchToStateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SwitchToStateMock.mockExpectations, StateSwitcherMockSwitchToStateParams{p, p1, p2},
			"StateSwitcher.SwitchToState got unexpected parameters")

		if m.SwitchToStateFunc == nil {

			m.t.Fatal("No results are set for the StateSwitcherMock.SwitchToState")

			return
		}
	}

	if m.SwitchToStateFunc == nil {
		m.t.Fatal("Unexpected call to StateSwitcherMock.SwitchToState")
		return
	}

	m.SwitchToStateFunc(p, p1, p2)
}

//SwitchToStateMinimockCounter returns a count of StateSwitcherMock.SwitchToStateFunc invocations
func (m *StateSwitcherMock) SwitchToStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SwitchToStateCounter)
}

//SwitchToStateMinimockPreCounter returns the value of StateSwitcherMock.SwitchToState invocations
func (m *StateSwitcherMock) SwitchToStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SwitchToStatePreCounter)
}

type mStateSwitcherMocksetState struct {
	mock             *StateSwitcherMock
	mockExpectations *StateSwitcherMocksetStateParams
}

//StateSwitcherMocksetStateParams represents input parameters of the StateSwitcher.setState
type StateSwitcherMocksetStateParams struct {
	p State
}

//Expect sets up expected params for the StateSwitcher.setState
func (m *mStateSwitcherMocksetState) Expect(p State) *mStateSwitcherMocksetState {
	m.mockExpectations = &StateSwitcherMocksetStateParams{p}
	return m
}

//Return sets up a mock for StateSwitcher.setState to return Return's arguments
func (m *mStateSwitcherMocksetState) Return() *StateSwitcherMock {
	m.mock.setStateFunc = func(p State) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of StateSwitcher.setState method
func (m *mStateSwitcherMocksetState) Set(f func(p State)) *StateSwitcherMock {
	m.mock.setStateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//setState implements github.com/insolar/insolar/pulsar.StateSwitcher interface
func (m *StateSwitcherMock) setState(p State) {
	atomic.AddUint64(&m.setStatePreCounter, 1)
	defer atomic.AddUint64(&m.setStateCounter, 1)

	if m.setStateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.setStateMock.mockExpectations, StateSwitcherMocksetStateParams{p},
			"StateSwitcher.setState got unexpected parameters")

		if m.setStateFunc == nil {

			m.t.Fatal("No results are set for the StateSwitcherMock.setState")

			return
		}
	}

	if m.setStateFunc == nil {
		m.t.Fatal("Unexpected call to StateSwitcherMock.setState")
		return
	}

	m.setStateFunc(p)
}

//setStateMinimockCounter returns a count of StateSwitcherMock.setStateFunc invocations
func (m *StateSwitcherMock) setStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.setStateCounter)
}

//setStateMinimockPreCounter returns the value of StateSwitcherMock.setState invocations
func (m *StateSwitcherMock) setStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.setStatePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateSwitcherMock) ValidateCallCounters() {

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.GetState")
	}

	if m.SetPulsarFunc != nil && atomic.LoadUint64(&m.SetPulsarCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.SetPulsar")
	}

	if m.SwitchToStateFunc != nil && atomic.LoadUint64(&m.SwitchToStateCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.SwitchToState")
	}

	if m.setStateFunc != nil && atomic.LoadUint64(&m.setStateCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.setState")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
//noinspection GoDeprecation
func (m *StateSwitcherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
//noinspection GoDeprecation
func (m *StateSwitcherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StateSwitcherMock) MinimockFinish() {

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.GetState")
	}

	if m.SetPulsarFunc != nil && atomic.LoadUint64(&m.SetPulsarCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.SetPulsar")
	}

	if m.SwitchToStateFunc != nil && atomic.LoadUint64(&m.SwitchToStateCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.SwitchToState")
	}

	if m.setStateFunc != nil && atomic.LoadUint64(&m.setStateCounter) == 0 {
		m.t.Fatal("Expected call to StateSwitcherMock.setState")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StateSwitcherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StateSwitcherMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetStateFunc == nil || atomic.LoadUint64(&m.GetStateCounter) > 0)
		ok = ok && (m.SetPulsarFunc == nil || atomic.LoadUint64(&m.SetPulsarCounter) > 0)
		ok = ok && (m.SwitchToStateFunc == nil || atomic.LoadUint64(&m.SwitchToStateCounter) > 0)
		ok = ok && (m.setStateFunc == nil || atomic.LoadUint64(&m.setStateCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
				m.t.Error("Expected call to StateSwitcherMock.GetState")
			}

			if m.SetPulsarFunc != nil && atomic.LoadUint64(&m.SetPulsarCounter) == 0 {
				m.t.Error("Expected call to StateSwitcherMock.SetPulsar")
			}

			if m.SwitchToStateFunc != nil && atomic.LoadUint64(&m.SwitchToStateCounter) == 0 {
				m.t.Error("Expected call to StateSwitcherMock.SwitchToState")
			}

			if m.setStateFunc != nil && atomic.LoadUint64(&m.setStateCounter) == 0 {
				m.t.Error("Expected call to StateSwitcherMock.setState")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. require.True(mock.AllMocksCalled())
func (m *StateSwitcherMock) AllMocksCalled() bool {

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		return false
	}

	if m.SetPulsarFunc != nil && atomic.LoadUint64(&m.SetPulsarCounter) == 0 {
		return false
	}

	if m.SwitchToStateFunc != nil && atomic.LoadUint64(&m.SwitchToStateCounter) == 0 {
		return false
	}

	if m.setStateFunc != nil && atomic.LoadUint64(&m.setStateCounter) == 0 {
		return false
	}

	return true
}

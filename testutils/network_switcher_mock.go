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
	mock              *NetworkSwitcherMock
	mainExpectation   *NetworkSwitcherMockGetStateExpectation
	expectationSeries []*NetworkSwitcherMockGetStateExpectation
}

type NetworkSwitcherMockGetStateExpectation struct {
	result *NetworkSwitcherMockGetStateResult
}

type NetworkSwitcherMockGetStateResult struct {
	r core.NetworkState
}

//Expect specifies that invocation of NetworkSwitcher.GetState is expected from 1 to Infinity times
func (m *mNetworkSwitcherMockGetState) Expect() *mNetworkSwitcherMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkSwitcherMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkSwitcher.GetState
func (m *mNetworkSwitcherMockGetState) Return(r core.NetworkState) *NetworkSwitcherMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkSwitcherMockGetStateExpectation{}
	}
	m.mainExpectation.result = &NetworkSwitcherMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkSwitcher.GetState is expected once
func (m *mNetworkSwitcherMockGetState) ExpectOnce() *NetworkSwitcherMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkSwitcherMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkSwitcherMockGetStateExpectation) Return(r core.NetworkState) {
	e.result = &NetworkSwitcherMockGetStateResult{r}
}

//Set uses given function f as a mock of NetworkSwitcher.GetState method
func (m *mNetworkSwitcherMockGetState) Set(f func() (r core.NetworkState)) *NetworkSwitcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/core.NetworkSwitcher interface
func (m *NetworkSwitcherMock) GetState() (r core.NetworkState) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkSwitcherMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkSwitcherMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkSwitcherMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkSwitcherMock.GetState.")
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

//GetStateFinished returns true if mock invocations count is ok
func (m *NetworkSwitcherMock) GetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStateCounter) == uint64(len(m.GetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStateFunc != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	return true
}

type mNetworkSwitcherMockOnPulse struct {
	mock              *NetworkSwitcherMock
	mainExpectation   *NetworkSwitcherMockOnPulseExpectation
	expectationSeries []*NetworkSwitcherMockOnPulseExpectation
}

type NetworkSwitcherMockOnPulseExpectation struct {
	input  *NetworkSwitcherMockOnPulseInput
	result *NetworkSwitcherMockOnPulseResult
}

type NetworkSwitcherMockOnPulseInput struct {
	p  context.Context
	p1 core.Pulse
}

type NetworkSwitcherMockOnPulseResult struct {
	r error
}

//Expect specifies that invocation of NetworkSwitcher.OnPulse is expected from 1 to Infinity times
func (m *mNetworkSwitcherMockOnPulse) Expect(p context.Context, p1 core.Pulse) *mNetworkSwitcherMockOnPulse {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkSwitcherMockOnPulseExpectation{}
	}
	m.mainExpectation.input = &NetworkSwitcherMockOnPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of NetworkSwitcher.OnPulse
func (m *mNetworkSwitcherMockOnPulse) Return(r error) *NetworkSwitcherMock {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkSwitcherMockOnPulseExpectation{}
	}
	m.mainExpectation.result = &NetworkSwitcherMockOnPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkSwitcher.OnPulse is expected once
func (m *mNetworkSwitcherMockOnPulse) ExpectOnce(p context.Context, p1 core.Pulse) *NetworkSwitcherMockOnPulseExpectation {
	m.mock.OnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkSwitcherMockOnPulseExpectation{}
	expectation.input = &NetworkSwitcherMockOnPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkSwitcherMockOnPulseExpectation) Return(r error) {
	e.result = &NetworkSwitcherMockOnPulseResult{r}
}

//Set uses given function f as a mock of NetworkSwitcher.OnPulse method
func (m *mNetworkSwitcherMockOnPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *NetworkSwitcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/core.NetworkSwitcher interface
func (m *NetworkSwitcherMock) OnPulse(p context.Context, p1 core.Pulse) (r error) {
	counter := atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if len(m.OnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkSwitcherMock.OnPulse. %v %v", p, p1)
			return
		}

		input := m.OnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkSwitcherMockOnPulseInput{p, p1}, "NetworkSwitcher.OnPulse got unexpected parameters")

		result := m.OnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkSwitcherMock.OnPulse")
			return
		}

		r = result.r

		return
	}

	if m.OnPulseMock.mainExpectation != nil {

		input := m.OnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkSwitcherMockOnPulseInput{p, p1}, "NetworkSwitcher.OnPulse got unexpected parameters")
		}

		result := m.OnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkSwitcherMock.OnPulse")
		}

		r = result.r

		return
	}

	if m.OnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkSwitcherMock.OnPulse. %v %v", p, p1)
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

//OnPulseFinished returns true if mock invocations count is ok
func (m *NetworkSwitcherMock) OnPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OnPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OnPulseCounter) == uint64(len(m.OnPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OnPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OnPulseFunc != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkSwitcherMock) ValidateCallCounters() {

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NetworkSwitcherMock.GetState")
	}

	if !m.OnPulseFinished() {
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

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NetworkSwitcherMock.GetState")
	}

	if !m.OnPulseFinished() {
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
		ok = ok && m.GetStateFinished()
		ok = ok && m.OnPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetStateFinished() {
				m.t.Error("Expected call to NetworkSwitcherMock.GetState")
			}

			if !m.OnPulseFinished() {
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

	if !m.GetStateFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	return true
}

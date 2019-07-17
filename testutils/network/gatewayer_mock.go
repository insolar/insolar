package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Gatewayer" can be found in github.com/insolar/insolar/network
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	network "github.com/insolar/insolar/network"

	testify_assert "github.com/stretchr/testify/assert"
)

//GatewayerMock implements github.com/insolar/insolar/network.Gatewayer
type GatewayerMock struct {
	t minimock.Tester

	GatewayFunc       func() (r network.Gateway)
	GatewayCounter    uint64
	GatewayPreCounter uint64
	GatewayMock       mGatewayerMockGateway

	SwitchStateFunc       func(p context.Context, p1 insolar.NetworkState)
	SwitchStateCounter    uint64
	SwitchStatePreCounter uint64
	SwitchStateMock       mGatewayerMockSwitchState
}

//NewGatewayerMock returns a mock for github.com/insolar/insolar/network.Gatewayer
func NewGatewayerMock(t minimock.Tester) *GatewayerMock {
	m := &GatewayerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GatewayMock = mGatewayerMockGateway{mock: m}
	m.SwitchStateMock = mGatewayerMockSwitchState{mock: m}

	return m
}

type mGatewayerMockGateway struct {
	mock              *GatewayerMock
	mainExpectation   *GatewayerMockGatewayExpectation
	expectationSeries []*GatewayerMockGatewayExpectation
}

type GatewayerMockGatewayExpectation struct {
	result *GatewayerMockGatewayResult
}

type GatewayerMockGatewayResult struct {
	r network.Gateway
}

//Expect specifies that invocation of Gatewayer.Gateway is expected from 1 to Infinity times
func (m *mGatewayerMockGateway) Expect() *mGatewayerMockGateway {
	m.mock.GatewayFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GatewayerMockGatewayExpectation{}
	}

	return m
}

//Return specifies results of invocation of Gatewayer.Gateway
func (m *mGatewayerMockGateway) Return(r network.Gateway) *GatewayerMock {
	m.mock.GatewayFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GatewayerMockGatewayExpectation{}
	}
	m.mainExpectation.result = &GatewayerMockGatewayResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Gatewayer.Gateway is expected once
func (m *mGatewayerMockGateway) ExpectOnce() *GatewayerMockGatewayExpectation {
	m.mock.GatewayFunc = nil
	m.mainExpectation = nil

	expectation := &GatewayerMockGatewayExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GatewayerMockGatewayExpectation) Return(r network.Gateway) {
	e.result = &GatewayerMockGatewayResult{r}
}

//Set uses given function f as a mock of Gatewayer.Gateway method
func (m *mGatewayerMockGateway) Set(f func() (r network.Gateway)) *GatewayerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GatewayFunc = f
	return m.mock
}

//Gateway implements github.com/insolar/insolar/network.Gatewayer interface
func (m *GatewayerMock) Gateway() (r network.Gateway) {
	counter := atomic.AddUint64(&m.GatewayPreCounter, 1)
	defer atomic.AddUint64(&m.GatewayCounter, 1)

	if len(m.GatewayMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GatewayMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GatewayerMock.Gateway.")
			return
		}

		result := m.GatewayMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GatewayerMock.Gateway")
			return
		}

		r = result.r

		return
	}

	if m.GatewayMock.mainExpectation != nil {

		result := m.GatewayMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GatewayerMock.Gateway")
		}

		r = result.r

		return
	}

	if m.GatewayFunc == nil {
		m.t.Fatalf("Unexpected call to GatewayerMock.Gateway.")
		return
	}

	return m.GatewayFunc()
}

//GatewayMinimockCounter returns a count of GatewayerMock.GatewayFunc invocations
func (m *GatewayerMock) GatewayMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GatewayCounter)
}

//GatewayMinimockPreCounter returns the value of GatewayerMock.Gateway invocations
func (m *GatewayerMock) GatewayMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GatewayPreCounter)
}

//GatewayFinished returns true if mock invocations count is ok
func (m *GatewayerMock) GatewayFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GatewayMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GatewayCounter) == uint64(len(m.GatewayMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GatewayMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GatewayCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GatewayFunc != nil {
		return atomic.LoadUint64(&m.GatewayCounter) > 0
	}

	return true
}

type mGatewayerMockSwitchState struct {
	mock              *GatewayerMock
	mainExpectation   *GatewayerMockSwitchStateExpectation
	expectationSeries []*GatewayerMockSwitchStateExpectation
}

type GatewayerMockSwitchStateExpectation struct {
	input *GatewayerMockSwitchStateInput
}

type GatewayerMockSwitchStateInput struct {
	p  context.Context
	p1 insolar.NetworkState
}

//Expect specifies that invocation of Gatewayer.SwitchState is expected from 1 to Infinity times
func (m *mGatewayerMockSwitchState) Expect(p context.Context, p1 insolar.NetworkState) *mGatewayerMockSwitchState {
	m.mock.SwitchStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GatewayerMockSwitchStateExpectation{}
	}
	m.mainExpectation.input = &GatewayerMockSwitchStateInput{p, p1}
	return m
}

//Return specifies results of invocation of Gatewayer.SwitchState
func (m *mGatewayerMockSwitchState) Return() *GatewayerMock {
	m.mock.SwitchStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GatewayerMockSwitchStateExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Gatewayer.SwitchState is expected once
func (m *mGatewayerMockSwitchState) ExpectOnce(p context.Context, p1 insolar.NetworkState) *GatewayerMockSwitchStateExpectation {
	m.mock.SwitchStateFunc = nil
	m.mainExpectation = nil

	expectation := &GatewayerMockSwitchStateExpectation{}
	expectation.input = &GatewayerMockSwitchStateInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Gatewayer.SwitchState method
func (m *mGatewayerMockSwitchState) Set(f func(p context.Context, p1 insolar.NetworkState)) *GatewayerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SwitchStateFunc = f
	return m.mock
}

//SwitchState implements github.com/insolar/insolar/network.Gatewayer interface
func (m *GatewayerMock) SwitchState(p context.Context, p1 insolar.NetworkState) {
	counter := atomic.AddUint64(&m.SwitchStatePreCounter, 1)
	defer atomic.AddUint64(&m.SwitchStateCounter, 1)

	if len(m.SwitchStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SwitchStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GatewayerMock.SwitchState. %v %v", p, p1)
			return
		}

		input := m.SwitchStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GatewayerMockSwitchStateInput{p, p1}, "Gatewayer.SwitchState got unexpected parameters")

		return
	}

	if m.SwitchStateMock.mainExpectation != nil {

		input := m.SwitchStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GatewayerMockSwitchStateInput{p, p1}, "Gatewayer.SwitchState got unexpected parameters")
		}

		return
	}

	if m.SwitchStateFunc == nil {
		m.t.Fatalf("Unexpected call to GatewayerMock.SwitchState. %v %v", p, p1)
		return
	}

	m.SwitchStateFunc(p, p1)
}

//SwitchStateMinimockCounter returns a count of GatewayerMock.SwitchStateFunc invocations
func (m *GatewayerMock) SwitchStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SwitchStateCounter)
}

//SwitchStateMinimockPreCounter returns the value of GatewayerMock.SwitchState invocations
func (m *GatewayerMock) SwitchStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SwitchStatePreCounter)
}

//SwitchStateFinished returns true if mock invocations count is ok
func (m *GatewayerMock) SwitchStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SwitchStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SwitchStateCounter) == uint64(len(m.SwitchStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SwitchStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SwitchStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SwitchStateFunc != nil {
		return atomic.LoadUint64(&m.SwitchStateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *GatewayerMock) ValidateCallCounters() {

	if !m.GatewayFinished() {
		m.t.Fatal("Expected call to GatewayerMock.Gateway")
	}

	if !m.SwitchStateFinished() {
		m.t.Fatal("Expected call to GatewayerMock.SwitchState")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *GatewayerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *GatewayerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *GatewayerMock) MinimockFinish() {

	if !m.GatewayFinished() {
		m.t.Fatal("Expected call to GatewayerMock.Gateway")
	}

	if !m.SwitchStateFinished() {
		m.t.Fatal("Expected call to GatewayerMock.SwitchState")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *GatewayerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *GatewayerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GatewayFinished()
		ok = ok && m.SwitchStateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GatewayFinished() {
				m.t.Error("Expected call to GatewayerMock.Gateway")
			}

			if !m.SwitchStateFinished() {
				m.t.Error("Expected call to GatewayerMock.SwitchState")
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
func (m *GatewayerMock) AllMocksCalled() bool {

	if !m.GatewayFinished() {
		return false
	}

	if !m.SwitchStateFinished() {
		return false
	}

	return true
}

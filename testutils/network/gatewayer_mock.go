package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Gatewayer" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
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

	SetGatewayFunc       func(p network.Gateway)
	SetGatewayCounter    uint64
	SetGatewayPreCounter uint64
	SetGatewayMock       mGatewayerMockSetGateway
}

//NewGatewayerMock returns a mock for github.com/insolar/insolar/network.Gatewayer
func NewGatewayerMock(t minimock.Tester) *GatewayerMock {
	m := &GatewayerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GatewayMock = mGatewayerMockGateway{mock: m}
	m.SetGatewayMock = mGatewayerMockSetGateway{mock: m}

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

type mGatewayerMockSetGateway struct {
	mock              *GatewayerMock
	mainExpectation   *GatewayerMockSetGatewayExpectation
	expectationSeries []*GatewayerMockSetGatewayExpectation
}

type GatewayerMockSetGatewayExpectation struct {
	input *GatewayerMockSetGatewayInput
}

type GatewayerMockSetGatewayInput struct {
	p network.Gateway
}

//Expect specifies that invocation of Gatewayer.SetGateway is expected from 1 to Infinity times
func (m *mGatewayerMockSetGateway) Expect(p network.Gateway) *mGatewayerMockSetGateway {
	m.mock.SetGatewayFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GatewayerMockSetGatewayExpectation{}
	}
	m.mainExpectation.input = &GatewayerMockSetGatewayInput{p}
	return m
}

//Return specifies results of invocation of Gatewayer.SetGateway
func (m *mGatewayerMockSetGateway) Return() *GatewayerMock {
	m.mock.SetGatewayFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GatewayerMockSetGatewayExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Gatewayer.SetGateway is expected once
func (m *mGatewayerMockSetGateway) ExpectOnce(p network.Gateway) *GatewayerMockSetGatewayExpectation {
	m.mock.SetGatewayFunc = nil
	m.mainExpectation = nil

	expectation := &GatewayerMockSetGatewayExpectation{}
	expectation.input = &GatewayerMockSetGatewayInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Gatewayer.SetGateway method
func (m *mGatewayerMockSetGateway) Set(f func(p network.Gateway)) *GatewayerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetGatewayFunc = f
	return m.mock
}

//SetGateway implements github.com/insolar/insolar/network.Gatewayer interface
func (m *GatewayerMock) SetGateway(p network.Gateway) {
	counter := atomic.AddUint64(&m.SetGatewayPreCounter, 1)
	defer atomic.AddUint64(&m.SetGatewayCounter, 1)

	if len(m.SetGatewayMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetGatewayMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GatewayerMock.SetGateway. %v", p)
			return
		}

		input := m.SetGatewayMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GatewayerMockSetGatewayInput{p}, "Gatewayer.SetGateway got unexpected parameters")

		return
	}

	if m.SetGatewayMock.mainExpectation != nil {

		input := m.SetGatewayMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GatewayerMockSetGatewayInput{p}, "Gatewayer.SetGateway got unexpected parameters")
		}

		return
	}

	if m.SetGatewayFunc == nil {
		m.t.Fatalf("Unexpected call to GatewayerMock.SetGateway. %v", p)
		return
	}

	m.SetGatewayFunc(p)
}

//SetGatewayMinimockCounter returns a count of GatewayerMock.SetGatewayFunc invocations
func (m *GatewayerMock) SetGatewayMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetGatewayCounter)
}

//SetGatewayMinimockPreCounter returns the value of GatewayerMock.SetGateway invocations
func (m *GatewayerMock) SetGatewayMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetGatewayPreCounter)
}

//SetGatewayFinished returns true if mock invocations count is ok
func (m *GatewayerMock) SetGatewayFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetGatewayMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetGatewayCounter) == uint64(len(m.SetGatewayMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetGatewayMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetGatewayCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetGatewayFunc != nil {
		return atomic.LoadUint64(&m.SetGatewayCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *GatewayerMock) ValidateCallCounters() {

	if !m.GatewayFinished() {
		m.t.Fatal("Expected call to GatewayerMock.Gateway")
	}

	if !m.SetGatewayFinished() {
		m.t.Fatal("Expected call to GatewayerMock.SetGateway")
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

	if !m.SetGatewayFinished() {
		m.t.Fatal("Expected call to GatewayerMock.SetGateway")
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
		ok = ok && m.SetGatewayFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GatewayFinished() {
				m.t.Error("Expected call to GatewayerMock.Gateway")
			}

			if !m.SetGatewayFinished() {
				m.t.Error("Expected call to GatewayerMock.SetGateway")
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

	if !m.SetGatewayFinished() {
		return false
	}

	return true
}

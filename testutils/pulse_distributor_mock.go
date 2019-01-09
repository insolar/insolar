package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseDistributor" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseDistributorMock implements github.com/insolar/insolar/core.PulseDistributor
type PulseDistributorMock struct {
	t minimock.Tester

	CloseFunc       func()
	CloseCounter    uint64
	ClosePreCounter uint64
	CloseMock       mPulseDistributorMockClose

	DistributeFunc       func(p context.Context, p1 *core.Pulse)
	DistributeCounter    uint64
	DistributePreCounter uint64
	DistributeMock       mPulseDistributorMockDistribute
}

//NewPulseDistributorMock returns a mock for github.com/insolar/insolar/core.PulseDistributor
func NewPulseDistributorMock(t minimock.Tester) *PulseDistributorMock {
	m := &PulseDistributorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseMock = mPulseDistributorMockClose{mock: m}
	m.DistributeMock = mPulseDistributorMockDistribute{mock: m}

	return m
}

type mPulseDistributorMockClose struct {
	mock              *PulseDistributorMock
	mainExpectation   *PulseDistributorMockCloseExpectation
	expectationSeries []*PulseDistributorMockCloseExpectation
}

type PulseDistributorMockCloseExpectation struct {
}

//Expect specifies that invocation of PulseDistributor.Close is expected from 1 to Infinity times
func (m *mPulseDistributorMockClose) Expect() *mPulseDistributorMockClose {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseDistributorMockCloseExpectation{}
	}

	return m
}

//Return specifies results of invocation of PulseDistributor.Close
func (m *mPulseDistributorMockClose) Return() *PulseDistributorMock {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseDistributorMockCloseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PulseDistributor.Close is expected once
func (m *mPulseDistributorMockClose) ExpectOnce() *PulseDistributorMockCloseExpectation {
	m.mock.CloseFunc = nil
	m.mainExpectation = nil

	expectation := &PulseDistributorMockCloseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PulseDistributor.Close method
func (m *mPulseDistributorMockClose) Set(f func()) *PulseDistributorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloseFunc = f
	return m.mock
}

//Close implements github.com/insolar/insolar/core.PulseDistributor interface
func (m *PulseDistributorMock) Close() {
	counter := atomic.AddUint64(&m.ClosePreCounter, 1)
	defer atomic.AddUint64(&m.CloseCounter, 1)

	if len(m.CloseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseDistributorMock.Close.")
			return
		}

		return
	}

	if m.CloseMock.mainExpectation != nil {

		return
	}

	if m.CloseFunc == nil {
		m.t.Fatalf("Unexpected call to PulseDistributorMock.Close.")
		return
	}

	m.CloseFunc()
}

//CloseMinimockCounter returns a count of PulseDistributorMock.CloseFunc invocations
func (m *PulseDistributorMock) CloseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseCounter)
}

//CloseMinimockPreCounter returns the value of PulseDistributorMock.Close invocations
func (m *PulseDistributorMock) CloseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClosePreCounter)
}

//CloseFinished returns true if mock invocations count is ok
func (m *PulseDistributorMock) CloseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloseCounter) == uint64(len(m.CloseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloseFunc != nil {
		return atomic.LoadUint64(&m.CloseCounter) > 0
	}

	return true
}

type mPulseDistributorMockDistribute struct {
	mock              *PulseDistributorMock
	mainExpectation   *PulseDistributorMockDistributeExpectation
	expectationSeries []*PulseDistributorMockDistributeExpectation
}

type PulseDistributorMockDistributeExpectation struct {
	input *PulseDistributorMockDistributeInput
}

type PulseDistributorMockDistributeInput struct {
	p  context.Context
	p1 *core.Pulse
}

//Expect specifies that invocation of PulseDistributor.Distribute is expected from 1 to Infinity times
func (m *mPulseDistributorMockDistribute) Expect(p context.Context, p1 *core.Pulse) *mPulseDistributorMockDistribute {
	m.mock.DistributeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseDistributorMockDistributeExpectation{}
	}
	m.mainExpectation.input = &PulseDistributorMockDistributeInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseDistributor.Distribute
func (m *mPulseDistributorMockDistribute) Return() *PulseDistributorMock {
	m.mock.DistributeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseDistributorMockDistributeExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PulseDistributor.Distribute is expected once
func (m *mPulseDistributorMockDistribute) ExpectOnce(p context.Context, p1 *core.Pulse) *PulseDistributorMockDistributeExpectation {
	m.mock.DistributeFunc = nil
	m.mainExpectation = nil

	expectation := &PulseDistributorMockDistributeExpectation{}
	expectation.input = &PulseDistributorMockDistributeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PulseDistributor.Distribute method
func (m *mPulseDistributorMockDistribute) Set(f func(p context.Context, p1 *core.Pulse)) *PulseDistributorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DistributeFunc = f
	return m.mock
}

//Distribute implements github.com/insolar/insolar/core.PulseDistributor interface
func (m *PulseDistributorMock) Distribute(p context.Context, p1 *core.Pulse) {
	counter := atomic.AddUint64(&m.DistributePreCounter, 1)
	defer atomic.AddUint64(&m.DistributeCounter, 1)

	if len(m.DistributeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DistributeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseDistributorMock.Distribute. %v %v", p, p1)
			return
		}

		input := m.DistributeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseDistributorMockDistributeInput{p, p1}, "PulseDistributor.Distribute got unexpected parameters")

		return
	}

	if m.DistributeMock.mainExpectation != nil {

		input := m.DistributeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseDistributorMockDistributeInput{p, p1}, "PulseDistributor.Distribute got unexpected parameters")
		}

		return
	}

	if m.DistributeFunc == nil {
		m.t.Fatalf("Unexpected call to PulseDistributorMock.Distribute. %v %v", p, p1)
		return
	}

	m.DistributeFunc(p, p1)
}

//DistributeMinimockCounter returns a count of PulseDistributorMock.DistributeFunc invocations
func (m *PulseDistributorMock) DistributeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DistributeCounter)
}

//DistributeMinimockPreCounter returns the value of PulseDistributorMock.Distribute invocations
func (m *PulseDistributorMock) DistributeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DistributePreCounter)
}

//DistributeFinished returns true if mock invocations count is ok
func (m *PulseDistributorMock) DistributeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DistributeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DistributeCounter) == uint64(len(m.DistributeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DistributeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DistributeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DistributeFunc != nil {
		return atomic.LoadUint64(&m.DistributeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseDistributorMock) ValidateCallCounters() {

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to PulseDistributorMock.Close")
	}

	if !m.DistributeFinished() {
		m.t.Fatal("Expected call to PulseDistributorMock.Distribute")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseDistributorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseDistributorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseDistributorMock) MinimockFinish() {

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to PulseDistributorMock.Close")
	}

	if !m.DistributeFinished() {
		m.t.Fatal("Expected call to PulseDistributorMock.Distribute")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseDistributorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseDistributorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloseFinished()
		ok = ok && m.DistributeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloseFinished() {
				m.t.Error("Expected call to PulseDistributorMock.Close")
			}

			if !m.DistributeFinished() {
				m.t.Error("Expected call to PulseDistributorMock.Distribute")
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
func (m *PulseDistributorMock) AllMocksCalled() bool {

	if !m.CloseFinished() {
		return false
	}

	if !m.DistributeFinished() {
		return false
	}

	return true
}

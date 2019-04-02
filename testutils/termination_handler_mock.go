package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "TerminationHandler" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//TerminationHandlerMock implements github.com/insolar/insolar/insolar.TerminationHandler
type TerminationHandlerMock struct {
	t minimock.Tester

	AbortFunc       func()
	AbortCounter    uint64
	AbortPreCounter uint64
	AbortMock       mTerminationHandlerMockAbort

	LeaveFunc       func(p context.Context, p1 insolar.PulseNumber)
	LeaveCounter    uint64
	LeavePreCounter uint64
	LeaveMock       mTerminationHandlerMockLeave

	OnLeaveApprovedFunc       func(p context.Context)
	OnLeaveApprovedCounter    uint64
	OnLeaveApprovedPreCounter uint64
	OnLeaveApprovedMock       mTerminationHandlerMockOnLeaveApproved
}

//NewTerminationHandlerMock returns a mock for github.com/insolar/insolar/insolar.TerminationHandler
func NewTerminationHandlerMock(t minimock.Tester) *TerminationHandlerMock {
	m := &TerminationHandlerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AbortMock = mTerminationHandlerMockAbort{mock: m}
	m.LeaveMock = mTerminationHandlerMockLeave{mock: m}
	m.OnLeaveApprovedMock = mTerminationHandlerMockOnLeaveApproved{mock: m}

	return m
}

type mTerminationHandlerMockAbort struct {
	mock              *TerminationHandlerMock
	mainExpectation   *TerminationHandlerMockAbortExpectation
	expectationSeries []*TerminationHandlerMockAbortExpectation
}

type TerminationHandlerMockAbortExpectation struct {
}

//Expect specifies that invocation of TerminationHandler.Abort is expected from 1 to Infinity times
func (m *mTerminationHandlerMockAbort) Expect() *mTerminationHandlerMockAbort {
	m.mock.AbortFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockAbortExpectation{}
	}

	return m
}

//Return specifies results of invocation of TerminationHandler.Abort
func (m *mTerminationHandlerMockAbort) Return() *TerminationHandlerMock {
	m.mock.AbortFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockAbortExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TerminationHandler.Abort is expected once
func (m *mTerminationHandlerMockAbort) ExpectOnce() *TerminationHandlerMockAbortExpectation {
	m.mock.AbortFunc = nil
	m.mainExpectation = nil

	expectation := &TerminationHandlerMockAbortExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TerminationHandler.Abort method
func (m *mTerminationHandlerMockAbort) Set(f func()) *TerminationHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AbortFunc = f
	return m.mock
}

//Abort implements github.com/insolar/insolar/insolar.TerminationHandler interface
func (m *TerminationHandlerMock) Abort() {
	counter := atomic.AddUint64(&m.AbortPreCounter, 1)
	defer atomic.AddUint64(&m.AbortCounter, 1)

	if len(m.AbortMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AbortMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TerminationHandlerMock.Abort.")
			return
		}

		return
	}

	if m.AbortMock.mainExpectation != nil {

		return
	}

	if m.AbortFunc == nil {
		m.t.Fatalf("Unexpected call to TerminationHandlerMock.Abort.")
		return
	}

	m.AbortFunc()
}

//AbortMinimockCounter returns a count of TerminationHandlerMock.AbortFunc invocations
func (m *TerminationHandlerMock) AbortMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AbortCounter)
}

//AbortMinimockPreCounter returns the value of TerminationHandlerMock.Abort invocations
func (m *TerminationHandlerMock) AbortMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AbortPreCounter)
}

//AbortFinished returns true if mock invocations count is ok
func (m *TerminationHandlerMock) AbortFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AbortMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AbortCounter) == uint64(len(m.AbortMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AbortMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AbortCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AbortFunc != nil {
		return atomic.LoadUint64(&m.AbortCounter) > 0
	}

	return true
}

type mTerminationHandlerMockLeave struct {
	mock              *TerminationHandlerMock
	mainExpectation   *TerminationHandlerMockLeaveExpectation
	expectationSeries []*TerminationHandlerMockLeaveExpectation
}

type TerminationHandlerMockLeaveExpectation struct {
	input *TerminationHandlerMockLeaveInput
}

type TerminationHandlerMockLeaveInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of TerminationHandler.Leave is expected from 1 to Infinity times
func (m *mTerminationHandlerMockLeave) Expect(p context.Context, p1 insolar.PulseNumber) *mTerminationHandlerMockLeave {
	m.mock.LeaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockLeaveExpectation{}
	}
	m.mainExpectation.input = &TerminationHandlerMockLeaveInput{p, p1}
	return m
}

//Return specifies results of invocation of TerminationHandler.Leave
func (m *mTerminationHandlerMockLeave) Return() *TerminationHandlerMock {
	m.mock.LeaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockLeaveExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TerminationHandler.Leave is expected once
func (m *mTerminationHandlerMockLeave) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *TerminationHandlerMockLeaveExpectation {
	m.mock.LeaveFunc = nil
	m.mainExpectation = nil

	expectation := &TerminationHandlerMockLeaveExpectation{}
	expectation.input = &TerminationHandlerMockLeaveInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TerminationHandler.Leave method
func (m *mTerminationHandlerMockLeave) Set(f func(p context.Context, p1 insolar.PulseNumber)) *TerminationHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LeaveFunc = f
	return m.mock
}

//Leave implements github.com/insolar/insolar/insolar.TerminationHandler interface
func (m *TerminationHandlerMock) Leave(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.LeavePreCounter, 1)
	defer atomic.AddUint64(&m.LeaveCounter, 1)

	if len(m.LeaveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LeaveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TerminationHandlerMock.Leave. %v %v", p, p1)
			return
		}

		input := m.LeaveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TerminationHandlerMockLeaveInput{p, p1}, "TerminationHandler.Leave got unexpected parameters")

		return
	}

	if m.LeaveMock.mainExpectation != nil {

		input := m.LeaveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TerminationHandlerMockLeaveInput{p, p1}, "TerminationHandler.Leave got unexpected parameters")
		}

		return
	}

	if m.LeaveFunc == nil {
		m.t.Fatalf("Unexpected call to TerminationHandlerMock.Leave. %v %v", p, p1)
		return
	}

	m.LeaveFunc(p, p1)
}

//LeaveMinimockCounter returns a count of TerminationHandlerMock.LeaveFunc invocations
func (m *TerminationHandlerMock) LeaveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LeaveCounter)
}

//LeaveMinimockPreCounter returns the value of TerminationHandlerMock.Leave invocations
func (m *TerminationHandlerMock) LeaveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LeavePreCounter)
}

//LeaveFinished returns true if mock invocations count is ok
func (m *TerminationHandlerMock) LeaveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LeaveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LeaveCounter) == uint64(len(m.LeaveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LeaveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LeaveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LeaveFunc != nil {
		return atomic.LoadUint64(&m.LeaveCounter) > 0
	}

	return true
}

type mTerminationHandlerMockOnLeaveApproved struct {
	mock              *TerminationHandlerMock
	mainExpectation   *TerminationHandlerMockOnLeaveApprovedExpectation
	expectationSeries []*TerminationHandlerMockOnLeaveApprovedExpectation
}

type TerminationHandlerMockOnLeaveApprovedExpectation struct {
	input *TerminationHandlerMockOnLeaveApprovedInput
}

type TerminationHandlerMockOnLeaveApprovedInput struct {
	p context.Context
}

//Expect specifies that invocation of TerminationHandler.OnLeaveApproved is expected from 1 to Infinity times
func (m *mTerminationHandlerMockOnLeaveApproved) Expect(p context.Context) *mTerminationHandlerMockOnLeaveApproved {
	m.mock.OnLeaveApprovedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockOnLeaveApprovedExpectation{}
	}
	m.mainExpectation.input = &TerminationHandlerMockOnLeaveApprovedInput{p}
	return m
}

//Return specifies results of invocation of TerminationHandler.OnLeaveApproved
func (m *mTerminationHandlerMockOnLeaveApproved) Return() *TerminationHandlerMock {
	m.mock.OnLeaveApprovedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TerminationHandlerMockOnLeaveApprovedExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TerminationHandler.OnLeaveApproved is expected once
func (m *mTerminationHandlerMockOnLeaveApproved) ExpectOnce(p context.Context) *TerminationHandlerMockOnLeaveApprovedExpectation {
	m.mock.OnLeaveApprovedFunc = nil
	m.mainExpectation = nil

	expectation := &TerminationHandlerMockOnLeaveApprovedExpectation{}
	expectation.input = &TerminationHandlerMockOnLeaveApprovedInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TerminationHandler.OnLeaveApproved method
func (m *mTerminationHandlerMockOnLeaveApproved) Set(f func(p context.Context)) *TerminationHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnLeaveApprovedFunc = f
	return m.mock
}

//OnLeaveApproved implements github.com/insolar/insolar/insolar.TerminationHandler interface
func (m *TerminationHandlerMock) OnLeaveApproved(p context.Context) {
	counter := atomic.AddUint64(&m.OnLeaveApprovedPreCounter, 1)
	defer atomic.AddUint64(&m.OnLeaveApprovedCounter, 1)

	if len(m.OnLeaveApprovedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnLeaveApprovedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TerminationHandlerMock.OnLeaveApproved. %v", p)
			return
		}

		input := m.OnLeaveApprovedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TerminationHandlerMockOnLeaveApprovedInput{p}, "TerminationHandler.OnLeaveApproved got unexpected parameters")

		return
	}

	if m.OnLeaveApprovedMock.mainExpectation != nil {

		input := m.OnLeaveApprovedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TerminationHandlerMockOnLeaveApprovedInput{p}, "TerminationHandler.OnLeaveApproved got unexpected parameters")
		}

		return
	}

	if m.OnLeaveApprovedFunc == nil {
		m.t.Fatalf("Unexpected call to TerminationHandlerMock.OnLeaveApproved. %v", p)
		return
	}

	m.OnLeaveApprovedFunc(p)
}

//OnLeaveApprovedMinimockCounter returns a count of TerminationHandlerMock.OnLeaveApprovedFunc invocations
func (m *TerminationHandlerMock) OnLeaveApprovedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnLeaveApprovedCounter)
}

//OnLeaveApprovedMinimockPreCounter returns the value of TerminationHandlerMock.OnLeaveApproved invocations
func (m *TerminationHandlerMock) OnLeaveApprovedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnLeaveApprovedPreCounter)
}

//OnLeaveApprovedFinished returns true if mock invocations count is ok
func (m *TerminationHandlerMock) OnLeaveApprovedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OnLeaveApprovedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OnLeaveApprovedCounter) == uint64(len(m.OnLeaveApprovedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OnLeaveApprovedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OnLeaveApprovedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OnLeaveApprovedFunc != nil {
		return atomic.LoadUint64(&m.OnLeaveApprovedCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TerminationHandlerMock) ValidateCallCounters() {

	if !m.AbortFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.Abort")
	}

	if !m.LeaveFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.Leave")
	}

	if !m.OnLeaveApprovedFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.OnLeaveApproved")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TerminationHandlerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TerminationHandlerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TerminationHandlerMock) MinimockFinish() {

	if !m.AbortFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.Abort")
	}

	if !m.LeaveFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.Leave")
	}

	if !m.OnLeaveApprovedFinished() {
		m.t.Fatal("Expected call to TerminationHandlerMock.OnLeaveApproved")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TerminationHandlerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TerminationHandlerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AbortFinished()
		ok = ok && m.LeaveFinished()
		ok = ok && m.OnLeaveApprovedFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AbortFinished() {
				m.t.Error("Expected call to TerminationHandlerMock.Abort")
			}

			if !m.LeaveFinished() {
				m.t.Error("Expected call to TerminationHandlerMock.Leave")
			}

			if !m.OnLeaveApprovedFinished() {
				m.t.Error("Expected call to TerminationHandlerMock.OnLeaveApproved")
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
func (m *TerminationHandlerMock) AllMocksCalled() bool {

	if !m.AbortFinished() {
		return false
	}

	if !m.LeaveFinished() {
		return false
	}

	if !m.OnLeaveApprovedFinished() {
		return false
	}

	return true
}

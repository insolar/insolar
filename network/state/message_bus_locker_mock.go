package state

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "messageBusLocker" can be found in github.com/insolar/insolar/network/state
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//messageBusLockerMock implements github.com/insolar/insolar/network/state.messageBusLocker
type messageBusLockerMock struct {
	t minimock.Tester

	AcquireFunc       func(p context.Context)
	AcquireCounter    uint64
	AcquirePreCounter uint64
	AcquireMock       mmessageBusLockerMockAcquire

	ReleaseFunc       func(p context.Context)
	ReleaseCounter    uint64
	ReleasePreCounter uint64
	ReleaseMock       mmessageBusLockerMockRelease
}

//NewmessageBusLockerMock returns a mock for github.com/insolar/insolar/network/state.messageBusLocker
func NewmessageBusLockerMock(t minimock.Tester) *messageBusLockerMock {
	m := &messageBusLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AcquireMock = mmessageBusLockerMockAcquire{mock: m}
	m.ReleaseMock = mmessageBusLockerMockRelease{mock: m}

	return m
}

type mmessageBusLockerMockAcquire struct {
	mock              *messageBusLockerMock
	mainExpectation   *messageBusLockerMockAcquireExpectation
	expectationSeries []*messageBusLockerMockAcquireExpectation
}

type messageBusLockerMockAcquireExpectation struct {
	input *messageBusLockerMockAcquireInput
}

type messageBusLockerMockAcquireInput struct {
	p context.Context
}

//Expect specifies that invocation of messageBusLocker.Acquire is expected from 1 to Infinity times
func (m *mmessageBusLockerMockAcquire) Expect(p context.Context) *mmessageBusLockerMockAcquire {
	m.mock.AcquireFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockAcquireExpectation{}
	}
	m.mainExpectation.input = &messageBusLockerMockAcquireInput{p}
	return m
}

//Return specifies results of invocation of messageBusLocker.Acquire
func (m *mmessageBusLockerMockAcquire) Return() *messageBusLockerMock {
	m.mock.AcquireFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockAcquireExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of messageBusLocker.Acquire is expected once
func (m *mmessageBusLockerMockAcquire) ExpectOnce(p context.Context) *messageBusLockerMockAcquireExpectation {
	m.mock.AcquireFunc = nil
	m.mainExpectation = nil

	expectation := &messageBusLockerMockAcquireExpectation{}
	expectation.input = &messageBusLockerMockAcquireInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of messageBusLocker.Acquire method
func (m *mmessageBusLockerMockAcquire) Set(f func(p context.Context)) *messageBusLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AcquireFunc = f
	return m.mock
}

//Acquire implements github.com/insolar/insolar/network/state.messageBusLocker interface
func (m *messageBusLockerMock) Acquire(p context.Context) {
	counter := atomic.AddUint64(&m.AcquirePreCounter, 1)
	defer atomic.AddUint64(&m.AcquireCounter, 1)

	if len(m.AcquireMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AcquireMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to messageBusLockerMock.Acquire. %v", p)
			return
		}

		input := m.AcquireMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, messageBusLockerMockAcquireInput{p}, "messageBusLocker.Acquire got unexpected parameters")

		return
	}

	if m.AcquireMock.mainExpectation != nil {

		input := m.AcquireMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, messageBusLockerMockAcquireInput{p}, "messageBusLocker.Acquire got unexpected parameters")
		}

		return
	}

	if m.AcquireFunc == nil {
		m.t.Fatalf("Unexpected call to messageBusLockerMock.Acquire. %v", p)
		return
	}

	m.AcquireFunc(p)
}

//AcquireMinimockCounter returns a count of messageBusLockerMock.AcquireFunc invocations
func (m *messageBusLockerMock) AcquireMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AcquireCounter)
}

//AcquireMinimockPreCounter returns the value of messageBusLockerMock.Acquire invocations
func (m *messageBusLockerMock) AcquireMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AcquirePreCounter)
}

//AcquireFinished returns true if mock invocations count is ok
func (m *messageBusLockerMock) AcquireFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AcquireMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AcquireCounter) == uint64(len(m.AcquireMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AcquireMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AcquireCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AcquireFunc != nil {
		return atomic.LoadUint64(&m.AcquireCounter) > 0
	}

	return true
}

type mmessageBusLockerMockRelease struct {
	mock              *messageBusLockerMock
	mainExpectation   *messageBusLockerMockReleaseExpectation
	expectationSeries []*messageBusLockerMockReleaseExpectation
}

type messageBusLockerMockReleaseExpectation struct {
	input *messageBusLockerMockReleaseInput
}

type messageBusLockerMockReleaseInput struct {
	p context.Context
}

//Expect specifies that invocation of messageBusLocker.Release is expected from 1 to Infinity times
func (m *mmessageBusLockerMockRelease) Expect(p context.Context) *mmessageBusLockerMockRelease {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockReleaseExpectation{}
	}
	m.mainExpectation.input = &messageBusLockerMockReleaseInput{p}
	return m
}

//Return specifies results of invocation of messageBusLocker.Release
func (m *mmessageBusLockerMockRelease) Return() *messageBusLockerMock {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockReleaseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of messageBusLocker.Release is expected once
func (m *mmessageBusLockerMockRelease) ExpectOnce(p context.Context) *messageBusLockerMockReleaseExpectation {
	m.mock.ReleaseFunc = nil
	m.mainExpectation = nil

	expectation := &messageBusLockerMockReleaseExpectation{}
	expectation.input = &messageBusLockerMockReleaseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of messageBusLocker.Release method
func (m *mmessageBusLockerMockRelease) Set(f func(p context.Context)) *messageBusLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReleaseFunc = f
	return m.mock
}

//Release implements github.com/insolar/insolar/network/state.messageBusLocker interface
func (m *messageBusLockerMock) Release(p context.Context) {
	counter := atomic.AddUint64(&m.ReleasePreCounter, 1)
	defer atomic.AddUint64(&m.ReleaseCounter, 1)

	if len(m.ReleaseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReleaseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to messageBusLockerMock.Release. %v", p)
			return
		}

		input := m.ReleaseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, messageBusLockerMockReleaseInput{p}, "messageBusLocker.Release got unexpected parameters")

		return
	}

	if m.ReleaseMock.mainExpectation != nil {

		input := m.ReleaseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, messageBusLockerMockReleaseInput{p}, "messageBusLocker.Release got unexpected parameters")
		}

		return
	}

	if m.ReleaseFunc == nil {
		m.t.Fatalf("Unexpected call to messageBusLockerMock.Release. %v", p)
		return
	}

	m.ReleaseFunc(p)
}

//ReleaseMinimockCounter returns a count of messageBusLockerMock.ReleaseFunc invocations
func (m *messageBusLockerMock) ReleaseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseCounter)
}

//ReleaseMinimockPreCounter returns the value of messageBusLockerMock.Release invocations
func (m *messageBusLockerMock) ReleaseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReleasePreCounter)
}

//ReleaseFinished returns true if mock invocations count is ok
func (m *messageBusLockerMock) ReleaseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReleaseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReleaseCounter) == uint64(len(m.ReleaseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReleaseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReleaseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReleaseFunc != nil {
		return atomic.LoadUint64(&m.ReleaseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *messageBusLockerMock) ValidateCallCounters() {

	if !m.AcquireFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Acquire")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Release")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *messageBusLockerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *messageBusLockerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *messageBusLockerMock) MinimockFinish() {

	if !m.AcquireFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Acquire")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Release")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *messageBusLockerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *messageBusLockerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AcquireFinished()
		ok = ok && m.ReleaseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AcquireFinished() {
				m.t.Error("Expected call to messageBusLockerMock.Acquire")
			}

			if !m.ReleaseFinished() {
				m.t.Error("Expected call to messageBusLockerMock.Release")
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
func (m *messageBusLockerMock) AllMocksCalled() bool {

	if !m.AcquireFinished() {
		return false
	}

	if !m.ReleaseFinished() {
		return false
	}

	return true
}

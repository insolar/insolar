package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NetworkLocker" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//NetworkLockerMock implements github.com/insolar/insolar/core.NetworkLocker
type NetworkLockerMock struct {
	t minimock.Tester

	AcquireGlobalLockFunc       func(p context.Context)
	AcquireGlobalLockCounter    uint64
	AcquireGlobalLockPreCounter uint64
	AcquireGlobalLockMock       mNetworkLockerMockAcquireGlobalLock

	ReleaseGlobalLockFunc       func(p context.Context)
	ReleaseGlobalLockCounter    uint64
	ReleaseGlobalLockPreCounter uint64
	ReleaseGlobalLockMock       mNetworkLockerMockReleaseGlobalLock
}

//NewNetworkLockerMock returns a mock for github.com/insolar/insolar/core.NetworkLocker
func NewNetworkLockerMock(t minimock.Tester) *NetworkLockerMock {
	m := &NetworkLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AcquireGlobalLockMock = mNetworkLockerMockAcquireGlobalLock{mock: m}
	m.ReleaseGlobalLockMock = mNetworkLockerMockReleaseGlobalLock{mock: m}

	return m
}

type mNetworkLockerMockAcquireGlobalLock struct {
	mock              *NetworkLockerMock
	mainExpectation   *NetworkLockerMockAcquireGlobalLockExpectation
	expectationSeries []*NetworkLockerMockAcquireGlobalLockExpectation
}

type NetworkLockerMockAcquireGlobalLockExpectation struct {
	input *NetworkLockerMockAcquireGlobalLockInput
}

type NetworkLockerMockAcquireGlobalLockInput struct {
	p context.Context
}

//Expect specifies that invocation of NetworkLocker.AcquireGlobalLock is expected from 1 to Infinity times
func (m *mNetworkLockerMockAcquireGlobalLock) Expect(p context.Context) *mNetworkLockerMockAcquireGlobalLock {
	m.mock.AcquireGlobalLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkLockerMockAcquireGlobalLockExpectation{}
	}
	m.mainExpectation.input = &NetworkLockerMockAcquireGlobalLockInput{p}
	return m
}

//Return specifies results of invocation of NetworkLocker.AcquireGlobalLock
func (m *mNetworkLockerMockAcquireGlobalLock) Return() *NetworkLockerMock {
	m.mock.AcquireGlobalLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkLockerMockAcquireGlobalLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NetworkLocker.AcquireGlobalLock is expected once
func (m *mNetworkLockerMockAcquireGlobalLock) ExpectOnce(p context.Context) *NetworkLockerMockAcquireGlobalLockExpectation {
	m.mock.AcquireGlobalLockFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkLockerMockAcquireGlobalLockExpectation{}
	expectation.input = &NetworkLockerMockAcquireGlobalLockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NetworkLocker.AcquireGlobalLock method
func (m *mNetworkLockerMockAcquireGlobalLock) Set(f func(p context.Context)) *NetworkLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AcquireGlobalLockFunc = f
	return m.mock
}

//AcquireGlobalLock implements github.com/insolar/insolar/core.NetworkLocker interface
func (m *NetworkLockerMock) AcquireGlobalLock(p context.Context) {
	counter := atomic.AddUint64(&m.AcquireGlobalLockPreCounter, 1)
	defer atomic.AddUint64(&m.AcquireGlobalLockCounter, 1)

	if len(m.AcquireGlobalLockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AcquireGlobalLockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkLockerMock.AcquireGlobalLock. %v", p)
			return
		}

		input := m.AcquireGlobalLockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkLockerMockAcquireGlobalLockInput{p}, "NetworkLocker.AcquireGlobalLock got unexpected parameters")

		return
	}

	if m.AcquireGlobalLockMock.mainExpectation != nil {

		input := m.AcquireGlobalLockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkLockerMockAcquireGlobalLockInput{p}, "NetworkLocker.AcquireGlobalLock got unexpected parameters")
		}

		return
	}

	if m.AcquireGlobalLockFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkLockerMock.AcquireGlobalLock. %v", p)
		return
	}

	m.AcquireGlobalLockFunc(p)
}

//AcquireGlobalLockMinimockCounter returns a count of NetworkLockerMock.AcquireGlobalLockFunc invocations
func (m *NetworkLockerMock) AcquireGlobalLockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AcquireGlobalLockCounter)
}

//AcquireGlobalLockMinimockPreCounter returns the value of NetworkLockerMock.AcquireGlobalLock invocations
func (m *NetworkLockerMock) AcquireGlobalLockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AcquireGlobalLockPreCounter)
}

//AcquireGlobalLockFinished returns true if mock invocations count is ok
func (m *NetworkLockerMock) AcquireGlobalLockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AcquireGlobalLockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AcquireGlobalLockCounter) == uint64(len(m.AcquireGlobalLockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AcquireGlobalLockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AcquireGlobalLockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AcquireGlobalLockFunc != nil {
		return atomic.LoadUint64(&m.AcquireGlobalLockCounter) > 0
	}

	return true
}

type mNetworkLockerMockReleaseGlobalLock struct {
	mock              *NetworkLockerMock
	mainExpectation   *NetworkLockerMockReleaseGlobalLockExpectation
	expectationSeries []*NetworkLockerMockReleaseGlobalLockExpectation
}

type NetworkLockerMockReleaseGlobalLockExpectation struct {
	input *NetworkLockerMockReleaseGlobalLockInput
}

type NetworkLockerMockReleaseGlobalLockInput struct {
	p context.Context
}

//Expect specifies that invocation of NetworkLocker.ReleaseGlobalLock is expected from 1 to Infinity times
func (m *mNetworkLockerMockReleaseGlobalLock) Expect(p context.Context) *mNetworkLockerMockReleaseGlobalLock {
	m.mock.ReleaseGlobalLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkLockerMockReleaseGlobalLockExpectation{}
	}
	m.mainExpectation.input = &NetworkLockerMockReleaseGlobalLockInput{p}
	return m
}

//Return specifies results of invocation of NetworkLocker.ReleaseGlobalLock
func (m *mNetworkLockerMockReleaseGlobalLock) Return() *NetworkLockerMock {
	m.mock.ReleaseGlobalLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkLockerMockReleaseGlobalLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NetworkLocker.ReleaseGlobalLock is expected once
func (m *mNetworkLockerMockReleaseGlobalLock) ExpectOnce(p context.Context) *NetworkLockerMockReleaseGlobalLockExpectation {
	m.mock.ReleaseGlobalLockFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkLockerMockReleaseGlobalLockExpectation{}
	expectation.input = &NetworkLockerMockReleaseGlobalLockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NetworkLocker.ReleaseGlobalLock method
func (m *mNetworkLockerMockReleaseGlobalLock) Set(f func(p context.Context)) *NetworkLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReleaseGlobalLockFunc = f
	return m.mock
}

//ReleaseGlobalLock implements github.com/insolar/insolar/core.NetworkLocker interface
func (m *NetworkLockerMock) ReleaseGlobalLock(p context.Context) {
	counter := atomic.AddUint64(&m.ReleaseGlobalLockPreCounter, 1)
	defer atomic.AddUint64(&m.ReleaseGlobalLockCounter, 1)

	if len(m.ReleaseGlobalLockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReleaseGlobalLockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkLockerMock.ReleaseGlobalLock. %v", p)
			return
		}

		input := m.ReleaseGlobalLockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkLockerMockReleaseGlobalLockInput{p}, "NetworkLocker.ReleaseGlobalLock got unexpected parameters")

		return
	}

	if m.ReleaseGlobalLockMock.mainExpectation != nil {

		input := m.ReleaseGlobalLockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkLockerMockReleaseGlobalLockInput{p}, "NetworkLocker.ReleaseGlobalLock got unexpected parameters")
		}

		return
	}

	if m.ReleaseGlobalLockFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkLockerMock.ReleaseGlobalLock. %v", p)
		return
	}

	m.ReleaseGlobalLockFunc(p)
}

//ReleaseGlobalLockMinimockCounter returns a count of NetworkLockerMock.ReleaseGlobalLockFunc invocations
func (m *NetworkLockerMock) ReleaseGlobalLockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseGlobalLockCounter)
}

//ReleaseGlobalLockMinimockPreCounter returns the value of NetworkLockerMock.ReleaseGlobalLock invocations
func (m *NetworkLockerMock) ReleaseGlobalLockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseGlobalLockPreCounter)
}

//ReleaseGlobalLockFinished returns true if mock invocations count is ok
func (m *NetworkLockerMock) ReleaseGlobalLockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReleaseGlobalLockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReleaseGlobalLockCounter) == uint64(len(m.ReleaseGlobalLockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReleaseGlobalLockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReleaseGlobalLockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReleaseGlobalLockFunc != nil {
		return atomic.LoadUint64(&m.ReleaseGlobalLockCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkLockerMock) ValidateCallCounters() {

	if !m.AcquireGlobalLockFinished() {
		m.t.Fatal("Expected call to NetworkLockerMock.AcquireGlobalLock")
	}

	if !m.ReleaseGlobalLockFinished() {
		m.t.Fatal("Expected call to NetworkLockerMock.ReleaseGlobalLock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkLockerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NetworkLockerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NetworkLockerMock) MinimockFinish() {

	if !m.AcquireGlobalLockFinished() {
		m.t.Fatal("Expected call to NetworkLockerMock.AcquireGlobalLock")
	}

	if !m.ReleaseGlobalLockFinished() {
		m.t.Fatal("Expected call to NetworkLockerMock.ReleaseGlobalLock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NetworkLockerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NetworkLockerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AcquireGlobalLockFinished()
		ok = ok && m.ReleaseGlobalLockFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AcquireGlobalLockFinished() {
				m.t.Error("Expected call to NetworkLockerMock.AcquireGlobalLock")
			}

			if !m.ReleaseGlobalLockFinished() {
				m.t.Error("Expected call to NetworkLockerMock.ReleaseGlobalLock")
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
func (m *NetworkLockerMock) AllMocksCalled() bool {

	if !m.AcquireGlobalLockFinished() {
		return false
	}

	if !m.ReleaseGlobalLockFinished() {
		return false
	}

	return true
}

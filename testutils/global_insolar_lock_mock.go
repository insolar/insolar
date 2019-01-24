package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "GlobalInsolarLock" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//GlobalInsolarLockMock implements github.com/insolar/insolar/core.GlobalInsolarLock
type GlobalInsolarLockMock struct {
	t minimock.Tester

	AcquireFunc       func(p context.Context)
	AcquireCounter    uint64
	AcquirePreCounter uint64
	AcquireMock       mGlobalInsolarLockMockAcquire

	ReleaseFunc       func(p context.Context)
	ReleaseCounter    uint64
	ReleasePreCounter uint64
	ReleaseMock       mGlobalInsolarLockMockRelease
}

//NewGlobalInsolarLockMock returns a mock for github.com/insolar/insolar/core.GlobalInsolarLock
func NewGlobalInsolarLockMock(t minimock.Tester) *GlobalInsolarLockMock {
	m := &GlobalInsolarLockMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AcquireMock = mGlobalInsolarLockMockAcquire{mock: m}
	m.ReleaseMock = mGlobalInsolarLockMockRelease{mock: m}

	return m
}

type mGlobalInsolarLockMockAcquire struct {
	mock              *GlobalInsolarLockMock
	mainExpectation   *GlobalInsolarLockMockAcquireExpectation
	expectationSeries []*GlobalInsolarLockMockAcquireExpectation
}

type GlobalInsolarLockMockAcquireExpectation struct {
	input *GlobalInsolarLockMockAcquireInput
}

type GlobalInsolarLockMockAcquireInput struct {
	p context.Context
}

//Expect specifies that invocation of GlobalInsolarLock.Acquire is expected from 1 to Infinity times
func (m *mGlobalInsolarLockMockAcquire) Expect(p context.Context) *mGlobalInsolarLockMockAcquire {
	m.mock.AcquireFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobalInsolarLockMockAcquireExpectation{}
	}
	m.mainExpectation.input = &GlobalInsolarLockMockAcquireInput{p}
	return m
}

//Return specifies results of invocation of GlobalInsolarLock.Acquire
func (m *mGlobalInsolarLockMockAcquire) Return() *GlobalInsolarLockMock {
	m.mock.AcquireFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobalInsolarLockMockAcquireExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of GlobalInsolarLock.Acquire is expected once
func (m *mGlobalInsolarLockMockAcquire) ExpectOnce(p context.Context) *GlobalInsolarLockMockAcquireExpectation {
	m.mock.AcquireFunc = nil
	m.mainExpectation = nil

	expectation := &GlobalInsolarLockMockAcquireExpectation{}
	expectation.input = &GlobalInsolarLockMockAcquireInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of GlobalInsolarLock.Acquire method
func (m *mGlobalInsolarLockMockAcquire) Set(f func(p context.Context)) *GlobalInsolarLockMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AcquireFunc = f
	return m.mock
}

//Acquire implements github.com/insolar/insolar/core.GlobalInsolarLock interface
func (m *GlobalInsolarLockMock) Acquire(p context.Context) {
	counter := atomic.AddUint64(&m.AcquirePreCounter, 1)
	defer atomic.AddUint64(&m.AcquireCounter, 1)

	if len(m.AcquireMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AcquireMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobalInsolarLockMock.Acquire. %v", p)
			return
		}

		input := m.AcquireMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GlobalInsolarLockMockAcquireInput{p}, "GlobalInsolarLock.Acquire got unexpected parameters")

		return
	}

	if m.AcquireMock.mainExpectation != nil {

		input := m.AcquireMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GlobalInsolarLockMockAcquireInput{p}, "GlobalInsolarLock.Acquire got unexpected parameters")
		}

		return
	}

	if m.AcquireFunc == nil {
		m.t.Fatalf("Unexpected call to GlobalInsolarLockMock.Acquire. %v", p)
		return
	}

	m.AcquireFunc(p)
}

//AcquireMinimockCounter returns a count of GlobalInsolarLockMock.AcquireFunc invocations
func (m *GlobalInsolarLockMock) AcquireMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AcquireCounter)
}

//AcquireMinimockPreCounter returns the value of GlobalInsolarLockMock.Acquire invocations
func (m *GlobalInsolarLockMock) AcquireMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AcquirePreCounter)
}

//AcquireFinished returns true if mock invocations count is ok
func (m *GlobalInsolarLockMock) AcquireFinished() bool {
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

type mGlobalInsolarLockMockRelease struct {
	mock              *GlobalInsolarLockMock
	mainExpectation   *GlobalInsolarLockMockReleaseExpectation
	expectationSeries []*GlobalInsolarLockMockReleaseExpectation
}

type GlobalInsolarLockMockReleaseExpectation struct {
	input *GlobalInsolarLockMockReleaseInput
}

type GlobalInsolarLockMockReleaseInput struct {
	p context.Context
}

//Expect specifies that invocation of GlobalInsolarLock.Release is expected from 1 to Infinity times
func (m *mGlobalInsolarLockMockRelease) Expect(p context.Context) *mGlobalInsolarLockMockRelease {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobalInsolarLockMockReleaseExpectation{}
	}
	m.mainExpectation.input = &GlobalInsolarLockMockReleaseInput{p}
	return m
}

//Return specifies results of invocation of GlobalInsolarLock.Release
func (m *mGlobalInsolarLockMockRelease) Return() *GlobalInsolarLockMock {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobalInsolarLockMockReleaseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of GlobalInsolarLock.Release is expected once
func (m *mGlobalInsolarLockMockRelease) ExpectOnce(p context.Context) *GlobalInsolarLockMockReleaseExpectation {
	m.mock.ReleaseFunc = nil
	m.mainExpectation = nil

	expectation := &GlobalInsolarLockMockReleaseExpectation{}
	expectation.input = &GlobalInsolarLockMockReleaseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of GlobalInsolarLock.Release method
func (m *mGlobalInsolarLockMockRelease) Set(f func(p context.Context)) *GlobalInsolarLockMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReleaseFunc = f
	return m.mock
}

//Release implements github.com/insolar/insolar/core.GlobalInsolarLock interface
func (m *GlobalInsolarLockMock) Release(p context.Context) {
	counter := atomic.AddUint64(&m.ReleasePreCounter, 1)
	defer atomic.AddUint64(&m.ReleaseCounter, 1)

	if len(m.ReleaseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReleaseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobalInsolarLockMock.Release. %v", p)
			return
		}

		input := m.ReleaseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GlobalInsolarLockMockReleaseInput{p}, "GlobalInsolarLock.Release got unexpected parameters")

		return
	}

	if m.ReleaseMock.mainExpectation != nil {

		input := m.ReleaseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GlobalInsolarLockMockReleaseInput{p}, "GlobalInsolarLock.Release got unexpected parameters")
		}

		return
	}

	if m.ReleaseFunc == nil {
		m.t.Fatalf("Unexpected call to GlobalInsolarLockMock.Release. %v", p)
		return
	}

	m.ReleaseFunc(p)
}

//ReleaseMinimockCounter returns a count of GlobalInsolarLockMock.ReleaseFunc invocations
func (m *GlobalInsolarLockMock) ReleaseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseCounter)
}

//ReleaseMinimockPreCounter returns the value of GlobalInsolarLockMock.Release invocations
func (m *GlobalInsolarLockMock) ReleaseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReleasePreCounter)
}

//ReleaseFinished returns true if mock invocations count is ok
func (m *GlobalInsolarLockMock) ReleaseFinished() bool {
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
func (m *GlobalInsolarLockMock) ValidateCallCounters() {

	if !m.AcquireFinished() {
		m.t.Fatal("Expected call to GlobalInsolarLockMock.Acquire")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to GlobalInsolarLockMock.Release")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *GlobalInsolarLockMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *GlobalInsolarLockMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *GlobalInsolarLockMock) MinimockFinish() {

	if !m.AcquireFinished() {
		m.t.Fatal("Expected call to GlobalInsolarLockMock.Acquire")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to GlobalInsolarLockMock.Release")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *GlobalInsolarLockMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *GlobalInsolarLockMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to GlobalInsolarLockMock.Acquire")
			}

			if !m.ReleaseFinished() {
				m.t.Error("Expected call to GlobalInsolarLockMock.Release")
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
func (m *GlobalInsolarLockMock) AllMocksCalled() bool {

	if !m.AcquireFinished() {
		return false
	}

	if !m.ReleaseFinished() {
		return false
	}

	return true
}

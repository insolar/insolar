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
	mock             *GlobalInsolarLockMock
	mockExpectations *GlobalInsolarLockMockAcquireParams
}

//GlobalInsolarLockMockAcquireParams represents input parameters of the GlobalInsolarLock.Acquire
type GlobalInsolarLockMockAcquireParams struct {
	p context.Context
}

//Expect sets up expected params for the GlobalInsolarLock.Acquire
func (m *mGlobalInsolarLockMockAcquire) Expect(p context.Context) *mGlobalInsolarLockMockAcquire {
	m.mockExpectations = &GlobalInsolarLockMockAcquireParams{p}
	return m
}

//Return sets up a mock for GlobalInsolarLock.Acquire to return Return's arguments
func (m *mGlobalInsolarLockMockAcquire) Return() *GlobalInsolarLockMock {
	m.mock.AcquireFunc = func(p context.Context) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of GlobalInsolarLock.Acquire method
func (m *mGlobalInsolarLockMockAcquire) Set(f func(p context.Context)) *GlobalInsolarLockMock {
	m.mock.AcquireFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Acquire implements github.com/insolar/insolar/core.GlobalInsolarLock interface
func (m *GlobalInsolarLockMock) Acquire(p context.Context) {
	atomic.AddUint64(&m.AcquirePreCounter, 1)
	defer atomic.AddUint64(&m.AcquireCounter, 1)

	if m.AcquireMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.AcquireMock.mockExpectations, GlobalInsolarLockMockAcquireParams{p},
			"GlobalInsolarLock.Acquire got unexpected parameters")

		if m.AcquireFunc == nil {

			m.t.Fatal("No results are set for the GlobalInsolarLockMock.Acquire")

			return
		}
	}

	if m.AcquireFunc == nil {
		m.t.Fatal("Unexpected call to GlobalInsolarLockMock.Acquire")
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

type mGlobalInsolarLockMockRelease struct {
	mock             *GlobalInsolarLockMock
	mockExpectations *GlobalInsolarLockMockReleaseParams
}

//GlobalInsolarLockMockReleaseParams represents input parameters of the GlobalInsolarLock.Release
type GlobalInsolarLockMockReleaseParams struct {
	p context.Context
}

//Expect sets up expected params for the GlobalInsolarLock.Release
func (m *mGlobalInsolarLockMockRelease) Expect(p context.Context) *mGlobalInsolarLockMockRelease {
	m.mockExpectations = &GlobalInsolarLockMockReleaseParams{p}
	return m
}

//Return sets up a mock for GlobalInsolarLock.Release to return Return's arguments
func (m *mGlobalInsolarLockMockRelease) Return() *GlobalInsolarLockMock {
	m.mock.ReleaseFunc = func(p context.Context) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of GlobalInsolarLock.Release method
func (m *mGlobalInsolarLockMockRelease) Set(f func(p context.Context)) *GlobalInsolarLockMock {
	m.mock.ReleaseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Release implements github.com/insolar/insolar/core.GlobalInsolarLock interface
func (m *GlobalInsolarLockMock) Release(p context.Context) {
	atomic.AddUint64(&m.ReleasePreCounter, 1)
	defer atomic.AddUint64(&m.ReleaseCounter, 1)

	if m.ReleaseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ReleaseMock.mockExpectations, GlobalInsolarLockMockReleaseParams{p},
			"GlobalInsolarLock.Release got unexpected parameters")

		if m.ReleaseFunc == nil {

			m.t.Fatal("No results are set for the GlobalInsolarLockMock.Release")

			return
		}
	}

	if m.ReleaseFunc == nil {
		m.t.Fatal("Unexpected call to GlobalInsolarLockMock.Release")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *GlobalInsolarLockMock) ValidateCallCounters() {

	if m.AcquireFunc != nil && atomic.LoadUint64(&m.AcquireCounter) == 0 {
		m.t.Fatal("Expected call to GlobalInsolarLockMock.Acquire")
	}

	if m.ReleaseFunc != nil && atomic.LoadUint64(&m.ReleaseCounter) == 0 {
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

	if m.AcquireFunc != nil && atomic.LoadUint64(&m.AcquireCounter) == 0 {
		m.t.Fatal("Expected call to GlobalInsolarLockMock.Acquire")
	}

	if m.ReleaseFunc != nil && atomic.LoadUint64(&m.ReleaseCounter) == 0 {
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
		ok = ok && (m.AcquireFunc == nil || atomic.LoadUint64(&m.AcquireCounter) > 0)
		ok = ok && (m.ReleaseFunc == nil || atomic.LoadUint64(&m.ReleaseCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.AcquireFunc != nil && atomic.LoadUint64(&m.AcquireCounter) == 0 {
				m.t.Error("Expected call to GlobalInsolarLockMock.Acquire")
			}

			if m.ReleaseFunc != nil && atomic.LoadUint64(&m.ReleaseCounter) == 0 {
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

	if m.AcquireFunc != nil && atomic.LoadUint64(&m.AcquireCounter) == 0 {
		return false
	}

	if m.ReleaseFunc != nil && atomic.LoadUint64(&m.ReleaseCounter) == 0 {
		return false
	}

	return true
}

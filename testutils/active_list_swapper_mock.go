package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ActiveListSwapper" can be found in github.com/insolar/insolar/ledger/pulsemanager
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//ActiveListSwapperMock implements github.com/insolar/insolar/ledger/pulsemanager.ActiveListSwapper
type ActiveListSwapperMock struct {
	t minimock.Tester

	MoveSyncToActiveFunc       func()
	MoveSyncToActiveCounter    uint64
	MoveSyncToActivePreCounter uint64
	MoveSyncToActiveMock       mActiveListSwapperMockMoveSyncToActive
}

//NewActiveListSwapperMock returns a mock for github.com/insolar/insolar/ledger/pulsemanager.ActiveListSwapper
func NewActiveListSwapperMock(t minimock.Tester) *ActiveListSwapperMock {
	m := &ActiveListSwapperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.MoveSyncToActiveMock = mActiveListSwapperMockMoveSyncToActive{mock: m}

	return m
}

type mActiveListSwapperMockMoveSyncToActive struct {
	mock *ActiveListSwapperMock
}

//Return sets up a mock for ActiveListSwapper.MoveSyncToActive to return Return's arguments
func (m *mActiveListSwapperMockMoveSyncToActive) Return() *ActiveListSwapperMock {
	m.mock.MoveSyncToActiveFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of ActiveListSwapper.MoveSyncToActive method
func (m *mActiveListSwapperMockMoveSyncToActive) Set(f func()) *ActiveListSwapperMock {
	m.mock.MoveSyncToActiveFunc = f

	return m.mock
}

//MoveSyncToActive implements github.com/insolar/insolar/ledger/pulsemanager.ActiveListSwapper interface
func (m *ActiveListSwapperMock) MoveSyncToActive() {
	atomic.AddUint64(&m.MoveSyncToActivePreCounter, 1)
	defer atomic.AddUint64(&m.MoveSyncToActiveCounter, 1)

	if m.MoveSyncToActiveFunc == nil {
		m.t.Fatal("Unexpected call to ActiveListSwapperMock.MoveSyncToActive")
		return
	}

	m.MoveSyncToActiveFunc()
}

//MoveSyncToActiveMinimockCounter returns a count of ActiveListSwapperMock.MoveSyncToActiveFunc invocations
func (m *ActiveListSwapperMock) MoveSyncToActiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MoveSyncToActiveCounter)
}

//MoveSyncToActiveMinimockPreCounter returns the value of ActiveListSwapperMock.MoveSyncToActive invocations
func (m *ActiveListSwapperMock) MoveSyncToActiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MoveSyncToActivePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ActiveListSwapperMock) ValidateCallCounters() {

	if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
		m.t.Fatal("Expected call to ActiveListSwapperMock.MoveSyncToActive")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ActiveListSwapperMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ActiveListSwapperMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ActiveListSwapperMock) MinimockFinish() {

	if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
		m.t.Fatal("Expected call to ActiveListSwapperMock.MoveSyncToActive")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ActiveListSwapperMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ActiveListSwapperMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.MoveSyncToActiveFunc == nil || atomic.LoadUint64(&m.MoveSyncToActiveCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
				m.t.Error("Expected call to ActiveListSwapperMock.MoveSyncToActive")
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
func (m *ActiveListSwapperMock) AllMocksCalled() bool {

	if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
		return false
	}

	return true
}

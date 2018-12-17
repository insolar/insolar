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

	LockFunc       func(p context.Context)
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mmessageBusLockerMockLock

	UnlockFunc       func(p context.Context)
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mmessageBusLockerMockUnlock
}

//NewmessageBusLockerMock returns a mock for github.com/insolar/insolar/network/state.messageBusLocker
func NewmessageBusLockerMock(t minimock.Tester) *messageBusLockerMock {
	m := &messageBusLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LockMock = mmessageBusLockerMockLock{mock: m}
	m.UnlockMock = mmessageBusLockerMockUnlock{mock: m}

	return m
}

type mmessageBusLockerMockLock struct {
	mock             *messageBusLockerMock
	mockExpectations *messageBusLockerMockLockParams
}

//messageBusLockerMockLockParams represents input parameters of the messageBusLocker.Lock
type messageBusLockerMockLockParams struct {
	p context.Context
}

//Expect sets up expected params for the messageBusLocker.Lock
func (m *mmessageBusLockerMockLock) Expect(p context.Context) *mmessageBusLockerMockLock {
	m.mockExpectations = &messageBusLockerMockLockParams{p}
	return m
}

//Return sets up a mock for messageBusLocker.Lock to return Return's arguments
func (m *mmessageBusLockerMockLock) Return() *messageBusLockerMock {
	m.mock.LockFunc = func(p context.Context) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of messageBusLocker.Lock method
func (m *mmessageBusLockerMockLock) Set(f func(p context.Context)) *messageBusLockerMock {
	m.mock.LockFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Lock implements github.com/insolar/insolar/network/state.messageBusLocker interface
func (m *messageBusLockerMock) Lock(p context.Context) {
	atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if m.LockMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LockMock.mockExpectations, messageBusLockerMockLockParams{p},
			"messageBusLocker.Lock got unexpected parameters")

		if m.LockFunc == nil {

			m.t.Fatal("No results are set for the messageBusLockerMock.Lock")

			return
		}
	}

	if m.LockFunc == nil {
		m.t.Fatal("Unexpected call to messageBusLockerMock.Lock")
		return
	}

	m.LockFunc(p)
}

//LockMinimockCounter returns a count of messageBusLockerMock.LockFunc invocations
func (m *messageBusLockerMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of messageBusLockerMock.Lock invocations
func (m *messageBusLockerMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

type mmessageBusLockerMockUnlock struct {
	mock             *messageBusLockerMock
	mockExpectations *messageBusLockerMockUnlockParams
}

//messageBusLockerMockUnlockParams represents input parameters of the messageBusLocker.Unlock
type messageBusLockerMockUnlockParams struct {
	p context.Context
}

//Expect sets up expected params for the messageBusLocker.Unlock
func (m *mmessageBusLockerMockUnlock) Expect(p context.Context) *mmessageBusLockerMockUnlock {
	m.mockExpectations = &messageBusLockerMockUnlockParams{p}
	return m
}

//Return sets up a mock for messageBusLocker.Unlock to return Return's arguments
func (m *mmessageBusLockerMockUnlock) Return() *messageBusLockerMock {
	m.mock.UnlockFunc = func(p context.Context) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of messageBusLocker.Unlock method
func (m *mmessageBusLockerMockUnlock) Set(f func(p context.Context)) *messageBusLockerMock {
	m.mock.UnlockFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Unlock implements github.com/insolar/insolar/network/state.messageBusLocker interface
func (m *messageBusLockerMock) Unlock(p context.Context) {
	atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if m.UnlockMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.UnlockMock.mockExpectations, messageBusLockerMockUnlockParams{p},
			"messageBusLocker.Unlock got unexpected parameters")

		if m.UnlockFunc == nil {

			m.t.Fatal("No results are set for the messageBusLockerMock.Unlock")

			return
		}
	}

	if m.UnlockFunc == nil {
		m.t.Fatal("Unexpected call to messageBusLockerMock.Unlock")
		return
	}

	m.UnlockFunc(p)
}

//UnlockMinimockCounter returns a count of messageBusLockerMock.UnlockFunc invocations
func (m *messageBusLockerMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of messageBusLockerMock.Unlock invocations
func (m *messageBusLockerMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *messageBusLockerMock) ValidateCallCounters() {

	if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
		m.t.Fatal("Expected call to messageBusLockerMock.Lock")
	}

	if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
		m.t.Fatal("Expected call to messageBusLockerMock.Unlock")
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

	if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
		m.t.Fatal("Expected call to messageBusLockerMock.Lock")
	}

	if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
		m.t.Fatal("Expected call to messageBusLockerMock.Unlock")
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
		ok = ok && (m.LockFunc == nil || atomic.LoadUint64(&m.LockCounter) > 0)
		ok = ok && (m.UnlockFunc == nil || atomic.LoadUint64(&m.UnlockCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
				m.t.Error("Expected call to messageBusLockerMock.Lock")
			}

			if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
				m.t.Error("Expected call to messageBusLockerMock.Unlock")
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

	if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
		return false
	}

	if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
		return false
	}

	return true
}

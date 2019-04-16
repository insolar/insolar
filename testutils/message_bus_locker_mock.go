package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MessageBusLocker" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//MessageBusLockerMock implements github.com/insolar/insolar/insolar.MessageBusLocker
type MessageBusLockerMock struct {
	t minimock.Tester

	LockFunc       func(p context.Context)
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mMessageBusLockerMockLock

	UnlockFunc       func(p context.Context)
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mMessageBusLockerMockUnlock
}

//NewMessageBusLockerMock returns a mock for github.com/insolar/insolar/insolar.MessageBusLocker
func NewMessageBusLockerMock(t minimock.Tester) *MessageBusLockerMock {
	m := &MessageBusLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LockMock = mMessageBusLockerMockLock{mock: m}
	m.UnlockMock = mMessageBusLockerMockUnlock{mock: m}

	return m
}

type mMessageBusLockerMockLock struct {
	mock              *MessageBusLockerMock
	mainExpectation   *MessageBusLockerMockLockExpectation
	expectationSeries []*MessageBusLockerMockLockExpectation
}

type MessageBusLockerMockLockExpectation struct {
	input *MessageBusLockerMockLockInput
}

type MessageBusLockerMockLockInput struct {
	p context.Context
}

//Expect specifies that invocation of MessageBusLocker.Lock is expected from 1 to Infinity times
func (m *mMessageBusLockerMockLock) Expect(p context.Context) *mMessageBusLockerMockLock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusLockerMockLockExpectation{}
	}
	m.mainExpectation.input = &MessageBusLockerMockLockInput{p}
	return m
}

//Return specifies results of invocation of MessageBusLocker.Lock
func (m *mMessageBusLockerMockLock) Return() *MessageBusLockerMock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusLockerMockLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of MessageBusLocker.Lock is expected once
func (m *mMessageBusLockerMockLock) ExpectOnce(p context.Context) *MessageBusLockerMockLockExpectation {
	m.mock.LockFunc = nil
	m.mainExpectation = nil

	expectation := &MessageBusLockerMockLockExpectation{}
	expectation.input = &MessageBusLockerMockLockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of MessageBusLocker.Lock method
func (m *mMessageBusLockerMockLock) Set(f func(p context.Context)) *MessageBusLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LockFunc = f
	return m.mock
}

//Lock implements github.com/insolar/insolar/insolar.MessageBusLocker interface
func (m *MessageBusLockerMock) Lock(p context.Context) {
	counter := atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if len(m.LockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MessageBusLockerMock.Lock. %v", p)
			return
		}

		input := m.LockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MessageBusLockerMockLockInput{p}, "MessageBusLocker.Lock got unexpected parameters")

		return
	}

	if m.LockMock.mainExpectation != nil {

		input := m.LockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MessageBusLockerMockLockInput{p}, "MessageBusLocker.Lock got unexpected parameters")
		}

		return
	}

	if m.LockFunc == nil {
		m.t.Fatalf("Unexpected call to MessageBusLockerMock.Lock. %v", p)
		return
	}

	m.LockFunc(p)
}

//LockMinimockCounter returns a count of MessageBusLockerMock.LockFunc invocations
func (m *MessageBusLockerMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of MessageBusLockerMock.Lock invocations
func (m *MessageBusLockerMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

//LockFinished returns true if mock invocations count is ok
func (m *MessageBusLockerMock) LockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LockCounter) == uint64(len(m.LockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LockFunc != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	return true
}

type mMessageBusLockerMockUnlock struct {
	mock              *MessageBusLockerMock
	mainExpectation   *MessageBusLockerMockUnlockExpectation
	expectationSeries []*MessageBusLockerMockUnlockExpectation
}

type MessageBusLockerMockUnlockExpectation struct {
	input *MessageBusLockerMockUnlockInput
}

type MessageBusLockerMockUnlockInput struct {
	p context.Context
}

//Expect specifies that invocation of MessageBusLocker.Unlock is expected from 1 to Infinity times
func (m *mMessageBusLockerMockUnlock) Expect(p context.Context) *mMessageBusLockerMockUnlock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusLockerMockUnlockExpectation{}
	}
	m.mainExpectation.input = &MessageBusLockerMockUnlockInput{p}
	return m
}

//Return specifies results of invocation of MessageBusLocker.Unlock
func (m *mMessageBusLockerMockUnlock) Return() *MessageBusLockerMock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusLockerMockUnlockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of MessageBusLocker.Unlock is expected once
func (m *mMessageBusLockerMockUnlock) ExpectOnce(p context.Context) *MessageBusLockerMockUnlockExpectation {
	m.mock.UnlockFunc = nil
	m.mainExpectation = nil

	expectation := &MessageBusLockerMockUnlockExpectation{}
	expectation.input = &MessageBusLockerMockUnlockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of MessageBusLocker.Unlock method
func (m *mMessageBusLockerMockUnlock) Set(f func(p context.Context)) *MessageBusLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnlockFunc = f
	return m.mock
}

//Unlock implements github.com/insolar/insolar/insolar.MessageBusLocker interface
func (m *MessageBusLockerMock) Unlock(p context.Context) {
	counter := atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if len(m.UnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MessageBusLockerMock.Unlock. %v", p)
			return
		}

		input := m.UnlockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MessageBusLockerMockUnlockInput{p}, "MessageBusLocker.Unlock got unexpected parameters")

		return
	}

	if m.UnlockMock.mainExpectation != nil {

		input := m.UnlockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MessageBusLockerMockUnlockInput{p}, "MessageBusLocker.Unlock got unexpected parameters")
		}

		return
	}

	if m.UnlockFunc == nil {
		m.t.Fatalf("Unexpected call to MessageBusLockerMock.Unlock. %v", p)
		return
	}

	m.UnlockFunc(p)
}

//UnlockMinimockCounter returns a count of MessageBusLockerMock.UnlockFunc invocations
func (m *MessageBusLockerMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of MessageBusLockerMock.Unlock invocations
func (m *MessageBusLockerMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//UnlockFinished returns true if mock invocations count is ok
func (m *MessageBusLockerMock) UnlockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UnlockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UnlockCounter) == uint64(len(m.UnlockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UnlockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UnlockFunc != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MessageBusLockerMock) ValidateCallCounters() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to MessageBusLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to MessageBusLockerMock.Unlock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MessageBusLockerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MessageBusLockerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MessageBusLockerMock) MinimockFinish() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to MessageBusLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to MessageBusLockerMock.Unlock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MessageBusLockerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MessageBusLockerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.LockFinished()
		ok = ok && m.UnlockFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.LockFinished() {
				m.t.Error("Expected call to MessageBusLockerMock.Lock")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to MessageBusLockerMock.Unlock")
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
func (m *MessageBusLockerMock) AllMocksCalled() bool {

	if !m.LockFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	return true
}

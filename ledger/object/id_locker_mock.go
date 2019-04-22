package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IDLocker" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IDLockerMock implements github.com/insolar/insolar/ledger/object.IDLocker
type IDLockerMock struct {
	t minimock.Tester

	LockFunc       func(p *insolar.ID)
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mIDLockerMockLock

	UnlockFunc       func(p *insolar.ID)
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mIDLockerMockUnlock
}

//NewIDLockerMock returns a mock for github.com/insolar/insolar/ledger/object.IDLocker
func NewIDLockerMock(t minimock.Tester) *IDLockerMock {
	m := &IDLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LockMock = mIDLockerMockLock{mock: m}
	m.UnlockMock = mIDLockerMockUnlock{mock: m}

	return m
}

type mIDLockerMockLock struct {
	mock              *IDLockerMock
	mainExpectation   *IDLockerMockLockExpectation
	expectationSeries []*IDLockerMockLockExpectation
}

type IDLockerMockLockExpectation struct {
	input *IDLockerMockLockInput
}

type IDLockerMockLockInput struct {
	p *insolar.ID
}

//Expect specifies that invocation of IDLocker.Lock is expected from 1 to Infinity times
func (m *mIDLockerMockLock) Expect(p *insolar.ID) *mIDLockerMockLock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockLockExpectation{}
	}
	m.mainExpectation.input = &IDLockerMockLockInput{p}
	return m
}

//Return specifies results of invocation of IDLocker.Lock
func (m *mIDLockerMockLock) Return() *IDLockerMock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of IDLocker.Lock is expected once
func (m *mIDLockerMockLock) ExpectOnce(p *insolar.ID) *IDLockerMockLockExpectation {
	m.mock.LockFunc = nil
	m.mainExpectation = nil

	expectation := &IDLockerMockLockExpectation{}
	expectation.input = &IDLockerMockLockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of IDLocker.Lock method
func (m *mIDLockerMockLock) Set(f func(p *insolar.ID)) *IDLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LockFunc = f
	return m.mock
}

//Lock implements github.com/insolar/insolar/ledger/object.IDLocker interface
func (m *IDLockerMock) Lock(p *insolar.ID) {
	counter := atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if len(m.LockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IDLockerMock.Lock. %v", p)
			return
		}

		input := m.LockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IDLockerMockLockInput{p}, "IDLocker.Lock got unexpected parameters")

		return
	}

	if m.LockMock.mainExpectation != nil {

		input := m.LockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IDLockerMockLockInput{p}, "IDLocker.Lock got unexpected parameters")
		}

		return
	}

	if m.LockFunc == nil {
		m.t.Fatalf("Unexpected call to IDLockerMock.Lock. %v", p)
		return
	}

	m.LockFunc(p)
}

//LockMinimockCounter returns a count of IDLockerMock.LockFunc invocations
func (m *IDLockerMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of IDLockerMock.Lock invocations
func (m *IDLockerMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

//LockFinished returns true if mock invocations count is ok
func (m *IDLockerMock) LockFinished() bool {
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

type mIDLockerMockUnlock struct {
	mock              *IDLockerMock
	mainExpectation   *IDLockerMockUnlockExpectation
	expectationSeries []*IDLockerMockUnlockExpectation
}

type IDLockerMockUnlockExpectation struct {
	input *IDLockerMockUnlockInput
}

type IDLockerMockUnlockInput struct {
	p *insolar.ID
}

//Expect specifies that invocation of IDLocker.Unlock is expected from 1 to Infinity times
func (m *mIDLockerMockUnlock) Expect(p *insolar.ID) *mIDLockerMockUnlock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockUnlockExpectation{}
	}
	m.mainExpectation.input = &IDLockerMockUnlockInput{p}
	return m
}

//Return specifies results of invocation of IDLocker.Unlock
func (m *mIDLockerMockUnlock) Return() *IDLockerMock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockUnlockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of IDLocker.Unlock is expected once
func (m *mIDLockerMockUnlock) ExpectOnce(p *insolar.ID) *IDLockerMockUnlockExpectation {
	m.mock.UnlockFunc = nil
	m.mainExpectation = nil

	expectation := &IDLockerMockUnlockExpectation{}
	expectation.input = &IDLockerMockUnlockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of IDLocker.Unlock method
func (m *mIDLockerMockUnlock) Set(f func(p *insolar.ID)) *IDLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnlockFunc = f
	return m.mock
}

//Unlock implements github.com/insolar/insolar/ledger/object.IDLocker interface
func (m *IDLockerMock) Unlock(p *insolar.ID) {
	counter := atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if len(m.UnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IDLockerMock.Unlock. %v", p)
			return
		}

		input := m.UnlockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IDLockerMockUnlockInput{p}, "IDLocker.Unlock got unexpected parameters")

		return
	}

	if m.UnlockMock.mainExpectation != nil {

		input := m.UnlockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IDLockerMockUnlockInput{p}, "IDLocker.Unlock got unexpected parameters")
		}

		return
	}

	if m.UnlockFunc == nil {
		m.t.Fatalf("Unexpected call to IDLockerMock.Unlock. %v", p)
		return
	}

	m.UnlockFunc(p)
}

//UnlockMinimockCounter returns a count of IDLockerMock.UnlockFunc invocations
func (m *IDLockerMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of IDLockerMock.Unlock invocations
func (m *IDLockerMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//UnlockFinished returns true if mock invocations count is ok
func (m *IDLockerMock) UnlockFinished() bool {
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
func (m *IDLockerMock) ValidateCallCounters() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.Unlock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IDLockerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IDLockerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IDLockerMock) MinimockFinish() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.Unlock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IDLockerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IDLockerMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IDLockerMock.Lock")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to IDLockerMock.Unlock")
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
func (m *IDLockerMock) AllMocksCalled() bool {

	if !m.LockFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	return true
}

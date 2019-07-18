package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexLocker" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexLockerMock implements github.com/insolar/insolar/ledger/object.IndexLocker
type IndexLockerMock struct {
	t minimock.Tester

	LockFunc       func(p insolar.ID)
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mIndexLockerMockLock

	UnlockFunc       func(p insolar.ID)
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mIndexLockerMockUnlock
}

//NewIndexLockerMock returns a mock for github.com/insolar/insolar/ledger/object.IndexLocker
func NewIndexLockerMock(t minimock.Tester) *IndexLockerMock {
	m := &IndexLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LockMock = mIndexLockerMockLock{mock: m}
	m.UnlockMock = mIndexLockerMockUnlock{mock: m}

	return m
}

type mIndexLockerMockLock struct {
	mock              *IndexLockerMock
	mainExpectation   *IndexLockerMockLockExpectation
	expectationSeries []*IndexLockerMockLockExpectation
}

type IndexLockerMockLockExpectation struct {
	input *IndexLockerMockLockInput
}

type IndexLockerMockLockInput struct {
	p insolar.ID
}

//Expect specifies that invocation of IndexLocker.Lock is expected from 1 to Infinity times
func (m *mIndexLockerMockLock) Expect(p insolar.ID) *mIndexLockerMockLock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLockerMockLockExpectation{}
	}
	m.mainExpectation.input = &IndexLockerMockLockInput{p}
	return m
}

//Return specifies results of invocation of IndexLocker.Lock
func (m *mIndexLockerMockLock) Return() *IndexLockerMock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLockerMockLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of IndexLocker.Lock is expected once
func (m *mIndexLockerMockLock) ExpectOnce(p insolar.ID) *IndexLockerMockLockExpectation {
	m.mock.LockFunc = nil
	m.mainExpectation = nil

	expectation := &IndexLockerMockLockExpectation{}
	expectation.input = &IndexLockerMockLockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of IndexLocker.Lock method
func (m *mIndexLockerMockLock) Set(f func(p insolar.ID)) *IndexLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LockFunc = f
	return m.mock
}

//Lock implements github.com/insolar/insolar/ledger/object.IndexLocker interface
func (m *IndexLockerMock) Lock(p insolar.ID) {
	counter := atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if len(m.LockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexLockerMock.Lock. %v", p)
			return
		}

		input := m.LockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexLockerMockLockInput{p}, "IndexLocker.Lock got unexpected parameters")

		return
	}

	if m.LockMock.mainExpectation != nil {

		input := m.LockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexLockerMockLockInput{p}, "IndexLocker.Lock got unexpected parameters")
		}

		return
	}

	if m.LockFunc == nil {
		m.t.Fatalf("Unexpected call to IndexLockerMock.Lock. %v", p)
		return
	}

	m.LockFunc(p)
}

//LockMinimockCounter returns a count of IndexLockerMock.LockFunc invocations
func (m *IndexLockerMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of IndexLockerMock.Lock invocations
func (m *IndexLockerMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

//LockFinished returns true if mock invocations count is ok
func (m *IndexLockerMock) LockFinished() bool {
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

type mIndexLockerMockUnlock struct {
	mock              *IndexLockerMock
	mainExpectation   *IndexLockerMockUnlockExpectation
	expectationSeries []*IndexLockerMockUnlockExpectation
}

type IndexLockerMockUnlockExpectation struct {
	input *IndexLockerMockUnlockInput
}

type IndexLockerMockUnlockInput struct {
	p insolar.ID
}

//Expect specifies that invocation of IndexLocker.Unlock is expected from 1 to Infinity times
func (m *mIndexLockerMockUnlock) Expect(p insolar.ID) *mIndexLockerMockUnlock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLockerMockUnlockExpectation{}
	}
	m.mainExpectation.input = &IndexLockerMockUnlockInput{p}
	return m
}

//Return specifies results of invocation of IndexLocker.Unlock
func (m *mIndexLockerMockUnlock) Return() *IndexLockerMock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLockerMockUnlockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of IndexLocker.Unlock is expected once
func (m *mIndexLockerMockUnlock) ExpectOnce(p insolar.ID) *IndexLockerMockUnlockExpectation {
	m.mock.UnlockFunc = nil
	m.mainExpectation = nil

	expectation := &IndexLockerMockUnlockExpectation{}
	expectation.input = &IndexLockerMockUnlockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of IndexLocker.Unlock method
func (m *mIndexLockerMockUnlock) Set(f func(p insolar.ID)) *IndexLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnlockFunc = f
	return m.mock
}

//Unlock implements github.com/insolar/insolar/ledger/object.IndexLocker interface
func (m *IndexLockerMock) Unlock(p insolar.ID) {
	counter := atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if len(m.UnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexLockerMock.Unlock. %v", p)
			return
		}

		input := m.UnlockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexLockerMockUnlockInput{p}, "IndexLocker.Unlock got unexpected parameters")

		return
	}

	if m.UnlockMock.mainExpectation != nil {

		input := m.UnlockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexLockerMockUnlockInput{p}, "IndexLocker.Unlock got unexpected parameters")
		}

		return
	}

	if m.UnlockFunc == nil {
		m.t.Fatalf("Unexpected call to IndexLockerMock.Unlock. %v", p)
		return
	}

	m.UnlockFunc(p)
}

//UnlockMinimockCounter returns a count of IndexLockerMock.UnlockFunc invocations
func (m *IndexLockerMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of IndexLockerMock.Unlock invocations
func (m *IndexLockerMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//UnlockFinished returns true if mock invocations count is ok
func (m *IndexLockerMock) UnlockFinished() bool {
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
func (m *IndexLockerMock) ValidateCallCounters() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to IndexLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to IndexLockerMock.Unlock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexLockerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexLockerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexLockerMock) MinimockFinish() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to IndexLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to IndexLockerMock.Unlock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexLockerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexLockerMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IndexLockerMock.Lock")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to IndexLockerMock.Unlock")
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
func (m *IndexLockerMock) AllMocksCalled() bool {

	if !m.LockFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	return true
}

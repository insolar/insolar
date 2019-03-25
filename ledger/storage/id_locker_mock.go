package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IDLocker" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IDLockerMock implements github.com/insolar/insolar/ledger/storage.IDLocker
type IDLockerMock struct {
	t minimock.Tester

	LockFunc       func(p *insolar.ID)
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mIDLockerMockLock

	RLockFunc       func(p *insolar.ID)
	RLockCounter    uint64
	RLockPreCounter uint64
	RLockMock       mIDLockerMockRLock

	RUnlockFunc       func(p *insolar.ID)
	RUnlockCounter    uint64
	RUnlockPreCounter uint64
	RUnlockMock       mIDLockerMockRUnlock

	UnlockFunc       func(p *insolar.ID)
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mIDLockerMockUnlock
}

//NewIDLockerMock returns a mock for github.com/insolar/insolar/ledger/storage.IDLocker
func NewIDLockerMock(t minimock.Tester) *IDLockerMock {
	m := &IDLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LockMock = mIDLockerMockLock{mock: m}
	m.RLockMock = mIDLockerMockRLock{mock: m}
	m.RUnlockMock = mIDLockerMockRUnlock{mock: m}
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

//Lock implements github.com/insolar/insolar/ledger/storage.IDLocker interface
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

type mIDLockerMockRLock struct {
	mock              *IDLockerMock
	mainExpectation   *IDLockerMockRLockExpectation
	expectationSeries []*IDLockerMockRLockExpectation
}

type IDLockerMockRLockExpectation struct {
	input *IDLockerMockRLockInput
}

type IDLockerMockRLockInput struct {
	p *insolar.ID
}

//Expect specifies that invocation of IDLocker.RLock is expected from 1 to Infinity times
func (m *mIDLockerMockRLock) Expect(p *insolar.ID) *mIDLockerMockRLock {
	m.mock.RLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockRLockExpectation{}
	}
	m.mainExpectation.input = &IDLockerMockRLockInput{p}
	return m
}

//Return specifies results of invocation of IDLocker.RLock
func (m *mIDLockerMockRLock) Return() *IDLockerMock {
	m.mock.RLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockRLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of IDLocker.RLock is expected once
func (m *mIDLockerMockRLock) ExpectOnce(p *insolar.ID) *IDLockerMockRLockExpectation {
	m.mock.RLockFunc = nil
	m.mainExpectation = nil

	expectation := &IDLockerMockRLockExpectation{}
	expectation.input = &IDLockerMockRLockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of IDLocker.RLock method
func (m *mIDLockerMockRLock) Set(f func(p *insolar.ID)) *IDLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RLockFunc = f
	return m.mock
}

//RLock implements github.com/insolar/insolar/ledger/storage.IDLocker interface
func (m *IDLockerMock) RLock(p *insolar.ID) {
	counter := atomic.AddUint64(&m.RLockPreCounter, 1)
	defer atomic.AddUint64(&m.RLockCounter, 1)

	if len(m.RLockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RLockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IDLockerMock.RLock. %v", p)
			return
		}

		input := m.RLockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IDLockerMockRLockInput{p}, "IDLocker.RLock got unexpected parameters")

		return
	}

	if m.RLockMock.mainExpectation != nil {

		input := m.RLockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IDLockerMockRLockInput{p}, "IDLocker.RLock got unexpected parameters")
		}

		return
	}

	if m.RLockFunc == nil {
		m.t.Fatalf("Unexpected call to IDLockerMock.RLock. %v", p)
		return
	}

	m.RLockFunc(p)
}

//RLockMinimockCounter returns a count of IDLockerMock.RLockFunc invocations
func (m *IDLockerMock) RLockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RLockCounter)
}

//RLockMinimockPreCounter returns the value of IDLockerMock.RLock invocations
func (m *IDLockerMock) RLockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RLockPreCounter)
}

//RLockFinished returns true if mock invocations count is ok
func (m *IDLockerMock) RLockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RLockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RLockCounter) == uint64(len(m.RLockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RLockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RLockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RLockFunc != nil {
		return atomic.LoadUint64(&m.RLockCounter) > 0
	}

	return true
}

type mIDLockerMockRUnlock struct {
	mock              *IDLockerMock
	mainExpectation   *IDLockerMockRUnlockExpectation
	expectationSeries []*IDLockerMockRUnlockExpectation
}

type IDLockerMockRUnlockExpectation struct {
	input *IDLockerMockRUnlockInput
}

type IDLockerMockRUnlockInput struct {
	p *insolar.ID
}

//Expect specifies that invocation of IDLocker.RUnlock is expected from 1 to Infinity times
func (m *mIDLockerMockRUnlock) Expect(p *insolar.ID) *mIDLockerMockRUnlock {
	m.mock.RUnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockRUnlockExpectation{}
	}
	m.mainExpectation.input = &IDLockerMockRUnlockInput{p}
	return m
}

//Return specifies results of invocation of IDLocker.RUnlock
func (m *mIDLockerMockRUnlock) Return() *IDLockerMock {
	m.mock.RUnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IDLockerMockRUnlockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of IDLocker.RUnlock is expected once
func (m *mIDLockerMockRUnlock) ExpectOnce(p *insolar.ID) *IDLockerMockRUnlockExpectation {
	m.mock.RUnlockFunc = nil
	m.mainExpectation = nil

	expectation := &IDLockerMockRUnlockExpectation{}
	expectation.input = &IDLockerMockRUnlockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of IDLocker.RUnlock method
func (m *mIDLockerMockRUnlock) Set(f func(p *insolar.ID)) *IDLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RUnlockFunc = f
	return m.mock
}

//RUnlock implements github.com/insolar/insolar/ledger/storage.IDLocker interface
func (m *IDLockerMock) RUnlock(p *insolar.ID) {
	counter := atomic.AddUint64(&m.RUnlockPreCounter, 1)
	defer atomic.AddUint64(&m.RUnlockCounter, 1)

	if len(m.RUnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RUnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IDLockerMock.RUnlock. %v", p)
			return
		}

		input := m.RUnlockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IDLockerMockRUnlockInput{p}, "IDLocker.RUnlock got unexpected parameters")

		return
	}

	if m.RUnlockMock.mainExpectation != nil {

		input := m.RUnlockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IDLockerMockRUnlockInput{p}, "IDLocker.RUnlock got unexpected parameters")
		}

		return
	}

	if m.RUnlockFunc == nil {
		m.t.Fatalf("Unexpected call to IDLockerMock.RUnlock. %v", p)
		return
	}

	m.RUnlockFunc(p)
}

//RUnlockMinimockCounter returns a count of IDLockerMock.RUnlockFunc invocations
func (m *IDLockerMock) RUnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RUnlockCounter)
}

//RUnlockMinimockPreCounter returns the value of IDLockerMock.RUnlock invocations
func (m *IDLockerMock) RUnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RUnlockPreCounter)
}

//RUnlockFinished returns true if mock invocations count is ok
func (m *IDLockerMock) RUnlockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RUnlockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RUnlockCounter) == uint64(len(m.RUnlockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RUnlockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RUnlockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RUnlockFunc != nil {
		return atomic.LoadUint64(&m.RUnlockCounter) > 0
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

//Unlock implements github.com/insolar/insolar/ledger/storage.IDLocker interface
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

	if !m.RLockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.RLock")
	}

	if !m.RUnlockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.RUnlock")
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

	if !m.RLockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.RLock")
	}

	if !m.RUnlockFinished() {
		m.t.Fatal("Expected call to IDLockerMock.RUnlock")
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
		ok = ok && m.RLockFinished()
		ok = ok && m.RUnlockFinished()
		ok = ok && m.UnlockFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.LockFinished() {
				m.t.Error("Expected call to IDLockerMock.Lock")
			}

			if !m.RLockFinished() {
				m.t.Error("Expected call to IDLockerMock.RLock")
			}

			if !m.RUnlockFinished() {
				m.t.Error("Expected call to IDLockerMock.RUnlock")
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

	if !m.RLockFinished() {
		return false
	}

	if !m.RUnlockFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	return true
}

package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PublicKeyStore" can be found in github.com/insolar/insolar/network/consensus/common/cryptkit
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//PublicKeyStoreMock implements github.com/insolar/insolar/network/consensus/common/cryptkit.PublicKeyStore
type PublicKeyStoreMock struct {
	t minimock.Tester

	PublicKeyStoreFunc       func()
	PublicKeyStoreCounter    uint64
	PublicKeyStorePreCounter uint64
	PublicKeyStoreMock       mPublicKeyStoreMockPublicKeyStore
}

//NewPublicKeyStoreMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptkit.PublicKeyStore
func NewPublicKeyStoreMock(t minimock.Tester) *PublicKeyStoreMock {
	m := &PublicKeyStoreMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PublicKeyStoreMock = mPublicKeyStoreMockPublicKeyStore{mock: m}

	return m
}

type mPublicKeyStoreMockPublicKeyStore struct {
	mock              *PublicKeyStoreMock
	mainExpectation   *PublicKeyStoreMockPublicKeyStoreExpectation
	expectationSeries []*PublicKeyStoreMockPublicKeyStoreExpectation
}

type PublicKeyStoreMockPublicKeyStoreExpectation struct {
}

//Expect specifies that invocation of PublicKeyStore.PublicKeyStore is expected from 1 to Infinity times
func (m *mPublicKeyStoreMockPublicKeyStore) Expect() *mPublicKeyStoreMockPublicKeyStore {
	m.mock.PublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublicKeyStoreMockPublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of PublicKeyStore.PublicKeyStore
func (m *mPublicKeyStoreMockPublicKeyStore) Return() *PublicKeyStoreMock {
	m.mock.PublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublicKeyStoreMockPublicKeyStoreExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PublicKeyStore.PublicKeyStore is expected once
func (m *mPublicKeyStoreMockPublicKeyStore) ExpectOnce() *PublicKeyStoreMockPublicKeyStoreExpectation {
	m.mock.PublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &PublicKeyStoreMockPublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PublicKeyStore.PublicKeyStore method
func (m *mPublicKeyStoreMockPublicKeyStore) Set(f func()) *PublicKeyStoreMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PublicKeyStoreFunc = f
	return m.mock
}

//PublicKeyStore implements github.com/insolar/insolar/network/consensus/common/cryptkit.PublicKeyStore interface
func (m *PublicKeyStoreMock) PublicKeyStore() {
	counter := atomic.AddUint64(&m.PublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.PublicKeyStoreCounter, 1)

	if len(m.PublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PublicKeyStoreMock.PublicKeyStore.")
			return
		}

		return
	}

	if m.PublicKeyStoreMock.mainExpectation != nil {

		return
	}

	if m.PublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to PublicKeyStoreMock.PublicKeyStore.")
		return
	}

	m.PublicKeyStoreFunc()
}

//PublicKeyStoreMinimockCounter returns a count of PublicKeyStoreMock.PublicKeyStoreFunc invocations
func (m *PublicKeyStoreMock) PublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PublicKeyStoreCounter)
}

//PublicKeyStoreMinimockPreCounter returns the value of PublicKeyStoreMock.PublicKeyStore invocations
func (m *PublicKeyStoreMock) PublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PublicKeyStorePreCounter)
}

//PublicKeyStoreFinished returns true if mock invocations count is ok
func (m *PublicKeyStoreMock) PublicKeyStoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PublicKeyStoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PublicKeyStoreCounter) == uint64(len(m.PublicKeyStoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PublicKeyStoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PublicKeyStoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PublicKeyStoreFunc != nil {
		return atomic.LoadUint64(&m.PublicKeyStoreCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PublicKeyStoreMock) ValidateCallCounters() {

	if !m.PublicKeyStoreFinished() {
		m.t.Fatal("Expected call to PublicKeyStoreMock.PublicKeyStore")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PublicKeyStoreMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PublicKeyStoreMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PublicKeyStoreMock) MinimockFinish() {

	if !m.PublicKeyStoreFinished() {
		m.t.Fatal("Expected call to PublicKeyStoreMock.PublicKeyStore")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PublicKeyStoreMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PublicKeyStoreMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.PublicKeyStoreFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.PublicKeyStoreFinished() {
				m.t.Error("Expected call to PublicKeyStoreMock.PublicKeyStore")
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
func (m *PublicKeyStoreMock) AllMocksCalled() bool {

	if !m.PublicKeyStoreFinished() {
		return false
	}

	return true
}

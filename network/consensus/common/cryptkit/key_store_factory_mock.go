package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "KeyStoreFactory" can be found in github.com/insolar/insolar/network/consensus/common/cryptkit
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//KeyStoreFactoryMock implements github.com/insolar/insolar/network/consensus/common/cryptkit.KeyStoreFactory
type KeyStoreFactoryMock struct {
	t minimock.Tester

	GetPublicKeyStoreFunc       func(p SignatureKeyHolder) (r PublicKeyStore)
	GetPublicKeyStoreCounter    uint64
	GetPublicKeyStorePreCounter uint64
	GetPublicKeyStoreMock       mKeyStoreFactoryMockGetPublicKeyStore
}

//NewKeyStoreFactoryMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptkit.KeyStoreFactory
func NewKeyStoreFactoryMock(t minimock.Tester) *KeyStoreFactoryMock {
	m := &KeyStoreFactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetPublicKeyStoreMock = mKeyStoreFactoryMockGetPublicKeyStore{mock: m}

	return m
}

type mKeyStoreFactoryMockGetPublicKeyStore struct {
	mock              *KeyStoreFactoryMock
	mainExpectation   *KeyStoreFactoryMockGetPublicKeyStoreExpectation
	expectationSeries []*KeyStoreFactoryMockGetPublicKeyStoreExpectation
}

type KeyStoreFactoryMockGetPublicKeyStoreExpectation struct {
	input  *KeyStoreFactoryMockGetPublicKeyStoreInput
	result *KeyStoreFactoryMockGetPublicKeyStoreResult
}

type KeyStoreFactoryMockGetPublicKeyStoreInput struct {
	p SignatureKeyHolder
}

type KeyStoreFactoryMockGetPublicKeyStoreResult struct {
	r PublicKeyStore
}

//Expect specifies that invocation of KeyStoreFactory.CreatePublicKeyStore is expected from 1 to Infinity times
func (m *mKeyStoreFactoryMockGetPublicKeyStore) Expect(p SignatureKeyHolder) *mKeyStoreFactoryMockGetPublicKeyStore {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyStoreFactoryMockGetPublicKeyStoreExpectation{}
	}
	m.mainExpectation.input = &KeyStoreFactoryMockGetPublicKeyStoreInput{p}
	return m
}

//Return specifies results of invocation of KeyStoreFactory.CreatePublicKeyStore
func (m *mKeyStoreFactoryMockGetPublicKeyStore) Return(r PublicKeyStore) *KeyStoreFactoryMock {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyStoreFactoryMockGetPublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &KeyStoreFactoryMockGetPublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyStoreFactory.CreatePublicKeyStore is expected once
func (m *mKeyStoreFactoryMockGetPublicKeyStore) ExpectOnce(p SignatureKeyHolder) *KeyStoreFactoryMockGetPublicKeyStoreExpectation {
	m.mock.GetPublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &KeyStoreFactoryMockGetPublicKeyStoreExpectation{}
	expectation.input = &KeyStoreFactoryMockGetPublicKeyStoreInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyStoreFactoryMockGetPublicKeyStoreExpectation) Return(r PublicKeyStore) {
	e.result = &KeyStoreFactoryMockGetPublicKeyStoreResult{r}
}

//Set uses given function f as a mock of KeyStoreFactory.CreatePublicKeyStore method
func (m *mKeyStoreFactoryMockGetPublicKeyStore) Set(f func(p SignatureKeyHolder) (r PublicKeyStore)) *KeyStoreFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyStoreFunc = f
	return m.mock
}

//CreatePublicKeyStore implements github.com/insolar/insolar/network/consensus/common/cryptkit.KeyStoreFactory interface
func (m *KeyStoreFactoryMock) CreatePublicKeyStore(p SignatureKeyHolder) (r PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetPublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyStoreCounter, 1)

	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyStoreFactoryMock.CreatePublicKeyStore. %v", p)
			return
		}

		input := m.GetPublicKeyStoreMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyStoreFactoryMockGetPublicKeyStoreInput{p}, "KeyStoreFactory.CreatePublicKeyStore got unexpected parameters")

		result := m.GetPublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyStoreFactoryMock.CreatePublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreMock.mainExpectation != nil {

		input := m.GetPublicKeyStoreMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyStoreFactoryMockGetPublicKeyStoreInput{p}, "KeyStoreFactory.CreatePublicKeyStore got unexpected parameters")
		}

		result := m.GetPublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyStoreFactoryMock.CreatePublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to KeyStoreFactoryMock.CreatePublicKeyStore. %v", p)
		return
	}

	return m.GetPublicKeyStoreFunc(p)
}

//GetPublicKeyStoreMinimockCounter returns a count of KeyStoreFactoryMock.GetPublicKeyStoreFunc invocations
func (m *KeyStoreFactoryMock) GetPublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStoreCounter)
}

//GetPublicKeyStoreMinimockPreCounter returns the value of KeyStoreFactoryMock.CreatePublicKeyStore invocations
func (m *KeyStoreFactoryMock) GetPublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStorePreCounter)
}

//GetPublicKeyStoreFinished returns true if mock invocations count is ok
func (m *KeyStoreFactoryMock) GetPublicKeyStoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) == uint64(len(m.GetPublicKeyStoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPublicKeyStoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPublicKeyStoreFunc != nil {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *KeyStoreFactoryMock) ValidateCallCounters() {

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to KeyStoreFactoryMock.CreatePublicKeyStore")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *KeyStoreFactoryMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *KeyStoreFactoryMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *KeyStoreFactoryMock) MinimockFinish() {

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to KeyStoreFactoryMock.CreatePublicKeyStore")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *KeyStoreFactoryMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *KeyStoreFactoryMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetPublicKeyStoreFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetPublicKeyStoreFinished() {
				m.t.Error("Expected call to KeyStoreFactoryMock.CreatePublicKeyStore")
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
func (m *KeyStoreFactoryMock) AllMocksCalled() bool {

	if !m.GetPublicKeyStoreFinished() {
		return false
	}

	return true
}

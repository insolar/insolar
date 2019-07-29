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

	CreatePublicKeyStoreFunc       func(p SignatureKeyHolder) (r PublicKeyStore)
	CreatePublicKeyStoreCounter    uint64
	CreatePublicKeyStorePreCounter uint64
	CreatePublicKeyStoreMock       mKeyStoreFactoryMockCreatePublicKeyStore
}

//NewKeyStoreFactoryMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptkit.KeyStoreFactory
func NewKeyStoreFactoryMock(t minimock.Tester) *KeyStoreFactoryMock {
	m := &KeyStoreFactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreatePublicKeyStoreMock = mKeyStoreFactoryMockCreatePublicKeyStore{mock: m}

	return m
}

type mKeyStoreFactoryMockCreatePublicKeyStore struct {
	mock              *KeyStoreFactoryMock
	mainExpectation   *KeyStoreFactoryMockCreatePublicKeyStoreExpectation
	expectationSeries []*KeyStoreFactoryMockCreatePublicKeyStoreExpectation
}

type KeyStoreFactoryMockCreatePublicKeyStoreExpectation struct {
	input  *KeyStoreFactoryMockCreatePublicKeyStoreInput
	result *KeyStoreFactoryMockCreatePublicKeyStoreResult
}

type KeyStoreFactoryMockCreatePublicKeyStoreInput struct {
	p SignatureKeyHolder
}

type KeyStoreFactoryMockCreatePublicKeyStoreResult struct {
	r PublicKeyStore
}

//Expect specifies that invocation of KeyStoreFactory.CreatePublicKeyStore is expected from 1 to Infinity times
func (m *mKeyStoreFactoryMockCreatePublicKeyStore) Expect(p SignatureKeyHolder) *mKeyStoreFactoryMockCreatePublicKeyStore {
	m.mock.CreatePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyStoreFactoryMockCreatePublicKeyStoreExpectation{}
	}
	m.mainExpectation.input = &KeyStoreFactoryMockCreatePublicKeyStoreInput{p}
	return m
}

//Return specifies results of invocation of KeyStoreFactory.CreatePublicKeyStore
func (m *mKeyStoreFactoryMockCreatePublicKeyStore) Return(r PublicKeyStore) *KeyStoreFactoryMock {
	m.mock.CreatePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyStoreFactoryMockCreatePublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &KeyStoreFactoryMockCreatePublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyStoreFactory.CreatePublicKeyStore is expected once
func (m *mKeyStoreFactoryMockCreatePublicKeyStore) ExpectOnce(p SignatureKeyHolder) *KeyStoreFactoryMockCreatePublicKeyStoreExpectation {
	m.mock.CreatePublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &KeyStoreFactoryMockCreatePublicKeyStoreExpectation{}
	expectation.input = &KeyStoreFactoryMockCreatePublicKeyStoreInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyStoreFactoryMockCreatePublicKeyStoreExpectation) Return(r PublicKeyStore) {
	e.result = &KeyStoreFactoryMockCreatePublicKeyStoreResult{r}
}

//Set uses given function f as a mock of KeyStoreFactory.CreatePublicKeyStore method
func (m *mKeyStoreFactoryMockCreatePublicKeyStore) Set(f func(p SignatureKeyHolder) (r PublicKeyStore)) *KeyStoreFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreatePublicKeyStoreFunc = f
	return m.mock
}

//CreatePublicKeyStore implements github.com/insolar/insolar/network/consensus/common/cryptkit.KeyStoreFactory interface
func (m *KeyStoreFactoryMock) CreatePublicKeyStore(p SignatureKeyHolder) (r PublicKeyStore) {
	counter := atomic.AddUint64(&m.CreatePublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.CreatePublicKeyStoreCounter, 1)

	if len(m.CreatePublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreatePublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyStoreFactoryMock.CreatePublicKeyStore. %v", p)
			return
		}

		input := m.CreatePublicKeyStoreMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyStoreFactoryMockCreatePublicKeyStoreInput{p}, "KeyStoreFactory.CreatePublicKeyStore got unexpected parameters")

		result := m.CreatePublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyStoreFactoryMock.CreatePublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.CreatePublicKeyStoreMock.mainExpectation != nil {

		input := m.CreatePublicKeyStoreMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyStoreFactoryMockCreatePublicKeyStoreInput{p}, "KeyStoreFactory.CreatePublicKeyStore got unexpected parameters")
		}

		result := m.CreatePublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyStoreFactoryMock.CreatePublicKeyStore")
		}

		r = result.r

		return
	}

	if m.CreatePublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to KeyStoreFactoryMock.CreatePublicKeyStore. %v", p)
		return
	}

	return m.CreatePublicKeyStoreFunc(p)
}

//CreatePublicKeyStoreMinimockCounter returns a count of KeyStoreFactoryMock.CreatePublicKeyStoreFunc invocations
func (m *KeyStoreFactoryMock) CreatePublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreatePublicKeyStoreCounter)
}

//CreatePublicKeyStoreMinimockPreCounter returns the value of KeyStoreFactoryMock.CreatePublicKeyStore invocations
func (m *KeyStoreFactoryMock) CreatePublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreatePublicKeyStorePreCounter)
}

//CreatePublicKeyStoreFinished returns true if mock invocations count is ok
func (m *KeyStoreFactoryMock) CreatePublicKeyStoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreatePublicKeyStoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreatePublicKeyStoreCounter) == uint64(len(m.CreatePublicKeyStoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreatePublicKeyStoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreatePublicKeyStoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreatePublicKeyStoreFunc != nil {
		return atomic.LoadUint64(&m.CreatePublicKeyStoreCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *KeyStoreFactoryMock) ValidateCallCounters() {

	if !m.CreatePublicKeyStoreFinished() {
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

	if !m.CreatePublicKeyStoreFinished() {
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
		ok = ok && m.CreatePublicKeyStoreFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CreatePublicKeyStoreFinished() {
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

	if !m.CreatePublicKeyStoreFinished() {
		return false
	}

	return true
}

package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SignatureVerifierFactory" can be found in github.com/insolar/insolar/network/consensus/common/cryptkit
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//SignatureVerifierFactoryMock implements github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureVerifierFactory
type SignatureVerifierFactoryMock struct {
	t minimock.Tester

	CreateSignatureVerifierWithPKSFunc       func(p PublicKeyStore) (r SignatureVerifier)
	CreateSignatureVerifierWithPKSCounter    uint64
	CreateSignatureVerifierWithPKSPreCounter uint64
	CreateSignatureVerifierWithPKSMock       mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS
}

//NewSignatureVerifierFactoryMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureVerifierFactory
func NewSignatureVerifierFactoryMock(t minimock.Tester) *SignatureVerifierFactoryMock {
	m := &SignatureVerifierFactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateSignatureVerifierWithPKSMock = mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS{mock: m}

	return m
}

type mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS struct {
	mock              *SignatureVerifierFactoryMock
	mainExpectation   *SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation
	expectationSeries []*SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation
}

type SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation struct {
	input  *SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSInput
	result *SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSResult
}

type SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSInput struct {
	p PublicKeyStore
}

type SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSResult struct {
	r SignatureVerifier
}

//Expect specifies that invocation of SignatureVerifierFactory.CreateSignatureVerifierWithPKS is expected from 1 to Infinity times
func (m *mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS) Expect(p PublicKeyStore) *mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS {
	m.mock.CreateSignatureVerifierWithPKSFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation{}
	}
	m.mainExpectation.input = &SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSInput{p}
	return m
}

//Return specifies results of invocation of SignatureVerifierFactory.CreateSignatureVerifierWithPKS
func (m *mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS) Return(r SignatureVerifier) *SignatureVerifierFactoryMock {
	m.mock.CreateSignatureVerifierWithPKSFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation{}
	}
	m.mainExpectation.result = &SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureVerifierFactory.CreateSignatureVerifierWithPKS is expected once
func (m *mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS) ExpectOnce(p PublicKeyStore) *SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation {
	m.mock.CreateSignatureVerifierWithPKSFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation{}
	expectation.input = &SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSExpectation) Return(r SignatureVerifier) {
	e.result = &SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSResult{r}
}

//Set uses given function f as a mock of SignatureVerifierFactory.CreateSignatureVerifierWithPKS method
func (m *mSignatureVerifierFactoryMockCreateSignatureVerifierWithPKS) Set(f func(p PublicKeyStore) (r SignatureVerifier)) *SignatureVerifierFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateSignatureVerifierWithPKSFunc = f
	return m.mock
}

//CreateSignatureVerifierWithPKS implements github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureVerifierFactory interface
func (m *SignatureVerifierFactoryMock) CreateSignatureVerifierWithPKS(p PublicKeyStore) (r SignatureVerifier) {
	counter := atomic.AddUint64(&m.CreateSignatureVerifierWithPKSPreCounter, 1)
	defer atomic.AddUint64(&m.CreateSignatureVerifierWithPKSCounter, 1)

	if len(m.CreateSignatureVerifierWithPKSMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateSignatureVerifierWithPKSMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS. %v", p)
			return
		}

		input := m.CreateSignatureVerifierWithPKSMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSInput{p}, "SignatureVerifierFactory.CreateSignatureVerifierWithPKS got unexpected parameters")

		result := m.CreateSignatureVerifierWithPKSMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS")
			return
		}

		r = result.r

		return
	}

	if m.CreateSignatureVerifierWithPKSMock.mainExpectation != nil {

		input := m.CreateSignatureVerifierWithPKSMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureVerifierFactoryMockCreateSignatureVerifierWithPKSInput{p}, "SignatureVerifierFactory.CreateSignatureVerifierWithPKS got unexpected parameters")
		}

		result := m.CreateSignatureVerifierWithPKSMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS")
		}

		r = result.r

		return
	}

	if m.CreateSignatureVerifierWithPKSFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS. %v", p)
		return
	}

	return m.CreateSignatureVerifierWithPKSFunc(p)
}

//CreateSignatureVerifierWithPKSMinimockCounter returns a count of SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKSFunc invocations
func (m *SignatureVerifierFactoryMock) CreateSignatureVerifierWithPKSMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateSignatureVerifierWithPKSCounter)
}

//CreateSignatureVerifierWithPKSMinimockPreCounter returns the value of SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS invocations
func (m *SignatureVerifierFactoryMock) CreateSignatureVerifierWithPKSMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateSignatureVerifierWithPKSPreCounter)
}

//CreateSignatureVerifierWithPKSFinished returns true if mock invocations count is ok
func (m *SignatureVerifierFactoryMock) CreateSignatureVerifierWithPKSFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreateSignatureVerifierWithPKSMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreateSignatureVerifierWithPKSCounter) == uint64(len(m.CreateSignatureVerifierWithPKSMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreateSignatureVerifierWithPKSMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreateSignatureVerifierWithPKSCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreateSignatureVerifierWithPKSFunc != nil {
		return atomic.LoadUint64(&m.CreateSignatureVerifierWithPKSCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SignatureVerifierFactoryMock) ValidateCallCounters() {

	if !m.CreateSignatureVerifierWithPKSFinished() {
		m.t.Fatal("Expected call to SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SignatureVerifierFactoryMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SignatureVerifierFactoryMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SignatureVerifierFactoryMock) MinimockFinish() {

	if !m.CreateSignatureVerifierWithPKSFinished() {
		m.t.Fatal("Expected call to SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SignatureVerifierFactoryMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SignatureVerifierFactoryMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CreateSignatureVerifierWithPKSFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CreateSignatureVerifierWithPKSFinished() {
				m.t.Error("Expected call to SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS")
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
func (m *SignatureVerifierFactoryMock) AllMocksCalled() bool {

	if !m.CreateSignatureVerifierWithPKSFinished() {
		return false
	}

	return true
}

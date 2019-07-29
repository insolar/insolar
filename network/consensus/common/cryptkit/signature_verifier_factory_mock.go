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

	GetSignatureVerifierWithPKSFunc       func(p PublicKeyStore) (r SignatureVerifier)
	GetSignatureVerifierWithPKSCounter    uint64
	GetSignatureVerifierWithPKSPreCounter uint64
	GetSignatureVerifierWithPKSMock       mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS
}

//NewSignatureVerifierFactoryMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureVerifierFactory
func NewSignatureVerifierFactoryMock(t minimock.Tester) *SignatureVerifierFactoryMock {
	m := &SignatureVerifierFactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetSignatureVerifierWithPKSMock = mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS{mock: m}

	return m
}

type mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS struct {
	mock              *SignatureVerifierFactoryMock
	mainExpectation   *SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation
	expectationSeries []*SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation
}

type SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation struct {
	input  *SignatureVerifierFactoryMockGetSignatureVerifierWithPKSInput
	result *SignatureVerifierFactoryMockGetSignatureVerifierWithPKSResult
}

type SignatureVerifierFactoryMockGetSignatureVerifierWithPKSInput struct {
	p PublicKeyStore
}

type SignatureVerifierFactoryMockGetSignatureVerifierWithPKSResult struct {
	r SignatureVerifier
}

//Expect specifies that invocation of SignatureVerifierFactory.CreateSignatureVerifierWithPKS is expected from 1 to Infinity times
func (m *mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS) Expect(p PublicKeyStore) *mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS {
	m.mock.GetSignatureVerifierWithPKSFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation{}
	}
	m.mainExpectation.input = &SignatureVerifierFactoryMockGetSignatureVerifierWithPKSInput{p}
	return m
}

//Return specifies results of invocation of SignatureVerifierFactory.CreateSignatureVerifierWithPKS
func (m *mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS) Return(r SignatureVerifier) *SignatureVerifierFactoryMock {
	m.mock.GetSignatureVerifierWithPKSFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation{}
	}
	m.mainExpectation.result = &SignatureVerifierFactoryMockGetSignatureVerifierWithPKSResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureVerifierFactory.CreateSignatureVerifierWithPKS is expected once
func (m *mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS) ExpectOnce(p PublicKeyStore) *SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation {
	m.mock.GetSignatureVerifierWithPKSFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation{}
	expectation.input = &SignatureVerifierFactoryMockGetSignatureVerifierWithPKSInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureVerifierFactoryMockGetSignatureVerifierWithPKSExpectation) Return(r SignatureVerifier) {
	e.result = &SignatureVerifierFactoryMockGetSignatureVerifierWithPKSResult{r}
}

//Set uses given function f as a mock of SignatureVerifierFactory.CreateSignatureVerifierWithPKS method
func (m *mSignatureVerifierFactoryMockGetSignatureVerifierWithPKS) Set(f func(p PublicKeyStore) (r SignatureVerifier)) *SignatureVerifierFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureVerifierWithPKSFunc = f
	return m.mock
}

//CreateSignatureVerifierWithPKS implements github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureVerifierFactory interface
func (m *SignatureVerifierFactoryMock) CreateSignatureVerifierWithPKS(p PublicKeyStore) (r SignatureVerifier) {
	counter := atomic.AddUint64(&m.GetSignatureVerifierWithPKSPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureVerifierWithPKSCounter, 1)

	if len(m.GetSignatureVerifierWithPKSMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureVerifierWithPKSMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS. %v", p)
			return
		}

		input := m.GetSignatureVerifierWithPKSMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureVerifierFactoryMockGetSignatureVerifierWithPKSInput{p}, "SignatureVerifierFactory.CreateSignatureVerifierWithPKS got unexpected parameters")

		result := m.GetSignatureVerifierWithPKSMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierWithPKSMock.mainExpectation != nil {

		input := m.GetSignatureVerifierWithPKSMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureVerifierFactoryMockGetSignatureVerifierWithPKSInput{p}, "SignatureVerifierFactory.CreateSignatureVerifierWithPKS got unexpected parameters")
		}

		result := m.GetSignatureVerifierWithPKSMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS")
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierWithPKSFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS. %v", p)
		return
	}

	return m.GetSignatureVerifierWithPKSFunc(p)
}

//GetSignatureVerifierWithPKSMinimockCounter returns a count of SignatureVerifierFactoryMock.GetSignatureVerifierWithPKSFunc invocations
func (m *SignatureVerifierFactoryMock) GetSignatureVerifierWithPKSMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierWithPKSCounter)
}

//GetSignatureVerifierWithPKSMinimockPreCounter returns the value of SignatureVerifierFactoryMock.CreateSignatureVerifierWithPKS invocations
func (m *SignatureVerifierFactoryMock) GetSignatureVerifierWithPKSMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierWithPKSPreCounter)
}

//GetSignatureVerifierWithPKSFinished returns true if mock invocations count is ok
func (m *SignatureVerifierFactoryMock) GetSignatureVerifierWithPKSFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureVerifierWithPKSMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureVerifierWithPKSCounter) == uint64(len(m.GetSignatureVerifierWithPKSMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureVerifierWithPKSMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureVerifierWithPKSCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureVerifierWithPKSFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureVerifierWithPKSCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SignatureVerifierFactoryMock) ValidateCallCounters() {

	if !m.GetSignatureVerifierWithPKSFinished() {
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

	if !m.GetSignatureVerifierWithPKSFinished() {
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
		ok = ok && m.GetSignatureVerifierWithPKSFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetSignatureVerifierWithPKSFinished() {
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

	if !m.GetSignatureVerifierWithPKSFinished() {
		return false
	}

	return true
}

package cryptography_containers

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DigestSigner" can be found in github.com/insolar/insolar/network/consensus/common
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//DigestSignerMock implements github.com/insolar/insolar/network/consensus/common.DigestSigner
type DigestSignerMock struct {
	t minimock.Tester

	GetSignMethodFunc       func() (r SignMethod)
	GetSignMethodCounter    uint64
	GetSignMethodPreCounter uint64
	GetSignMethodMock       mDigestSignerMockGetSignMethod

	SignDigestFunc       func(p Digest) (r Signature)
	SignDigestCounter    uint64
	SignDigestPreCounter uint64
	SignDigestMock       mDigestSignerMockSignDigest
}

//NewDigestSignerMock returns a mock for github.com/insolar/insolar/network/consensus/common.DigestSigner
func NewDigestSignerMock(t minimock.Tester) *DigestSignerMock {
	m := &DigestSignerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetSignMethodMock = mDigestSignerMockGetSignMethod{mock: m}
	m.SignDigestMock = mDigestSignerMockSignDigest{mock: m}

	return m
}

type mDigestSignerMockGetSignMethod struct {
	mock              *DigestSignerMock
	mainExpectation   *DigestSignerMockGetSignMethodExpectation
	expectationSeries []*DigestSignerMockGetSignMethodExpectation
}

type DigestSignerMockGetSignMethodExpectation struct {
	result *DigestSignerMockGetSignMethodResult
}

type DigestSignerMockGetSignMethodResult struct {
	r SignMethod
}

//Expect specifies that invocation of DigestSigner.GetSignMethod is expected from 1 to Infinity times
func (m *mDigestSignerMockGetSignMethod) Expect() *mDigestSignerMockGetSignMethod {
	m.mock.GetSignMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestSignerMockGetSignMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of DigestSigner.GetSignMethod
func (m *mDigestSignerMockGetSignMethod) Return(r SignMethod) *DigestSignerMock {
	m.mock.GetSignMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestSignerMockGetSignMethodExpectation{}
	}
	m.mainExpectation.result = &DigestSignerMockGetSignMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestSigner.GetSignMethod is expected once
func (m *mDigestSignerMockGetSignMethod) ExpectOnce() *DigestSignerMockGetSignMethodExpectation {
	m.mock.GetSignMethodFunc = nil
	m.mainExpectation = nil

	expectation := &DigestSignerMockGetSignMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestSignerMockGetSignMethodExpectation) Return(r SignMethod) {
	e.result = &DigestSignerMockGetSignMethodResult{r}
}

//Set uses given function f as a mock of DigestSigner.GetSignMethod method
func (m *mDigestSignerMockGetSignMethod) Set(f func() (r SignMethod)) *DigestSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignMethodFunc = f
	return m.mock
}

//GetSignMethod implements github.com/insolar/insolar/network/consensus/common.DigestSigner interface
func (m *DigestSignerMock) GetSignMethod() (r SignMethod) {
	counter := atomic.AddUint64(&m.GetSignMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignMethodCounter, 1)

	if len(m.GetSignMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestSignerMock.GetSignMethod.")
			return
		}

		result := m.GetSignMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestSignerMock.GetSignMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignMethodMock.mainExpectation != nil {

		result := m.GetSignMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestSignerMock.GetSignMethod")
		}

		r = result.r

		return
	}

	if m.GetSignMethodFunc == nil {
		m.t.Fatalf("Unexpected call to DigestSignerMock.GetSignMethod.")
		return
	}

	return m.GetSignMethodFunc()
}

//GetSignMethodMinimockCounter returns a count of DigestSignerMock.GetSignMethodFunc invocations
func (m *DigestSignerMock) GetSignMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignMethodCounter)
}

//GetSignMethodMinimockPreCounter returns the value of DigestSignerMock.GetSignMethod invocations
func (m *DigestSignerMock) GetSignMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignMethodPreCounter)
}

//GetSignMethodFinished returns true if mock invocations count is ok
func (m *DigestSignerMock) GetSignMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignMethodCounter) == uint64(len(m.GetSignMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignMethodFunc != nil {
		return atomic.LoadUint64(&m.GetSignMethodCounter) > 0
	}

	return true
}

type mDigestSignerMockSignDigest struct {
	mock              *DigestSignerMock
	mainExpectation   *DigestSignerMockSignDigestExpectation
	expectationSeries []*DigestSignerMockSignDigestExpectation
}

type DigestSignerMockSignDigestExpectation struct {
	input  *DigestSignerMockSignDigestInput
	result *DigestSignerMockSignDigestResult
}

type DigestSignerMockSignDigestInput struct {
	p Digest
}

type DigestSignerMockSignDigestResult struct {
	r Signature
}

//Expect specifies that invocation of DigestSigner.SignDigest is expected from 1 to Infinity times
func (m *mDigestSignerMockSignDigest) Expect(p Digest) *mDigestSignerMockSignDigest {
	m.mock.SignDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestSignerMockSignDigestExpectation{}
	}
	m.mainExpectation.input = &DigestSignerMockSignDigestInput{p}
	return m
}

//Return specifies results of invocation of DigestSigner.SignDigest
func (m *mDigestSignerMockSignDigest) Return(r Signature) *DigestSignerMock {
	m.mock.SignDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestSignerMockSignDigestExpectation{}
	}
	m.mainExpectation.result = &DigestSignerMockSignDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestSigner.SignDigest is expected once
func (m *mDigestSignerMockSignDigest) ExpectOnce(p Digest) *DigestSignerMockSignDigestExpectation {
	m.mock.SignDigestFunc = nil
	m.mainExpectation = nil

	expectation := &DigestSignerMockSignDigestExpectation{}
	expectation.input = &DigestSignerMockSignDigestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestSignerMockSignDigestExpectation) Return(r Signature) {
	e.result = &DigestSignerMockSignDigestResult{r}
}

//Set uses given function f as a mock of DigestSigner.SignDigest method
func (m *mDigestSignerMockSignDigest) Set(f func(p Digest) (r Signature)) *DigestSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignDigestFunc = f
	return m.mock
}

//SignDigest implements github.com/insolar/insolar/network/consensus/common.DigestSigner interface
func (m *DigestSignerMock) SignDigest(p Digest) (r Signature) {
	counter := atomic.AddUint64(&m.SignDigestPreCounter, 1)
	defer atomic.AddUint64(&m.SignDigestCounter, 1)

	if len(m.SignDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestSignerMock.SignDigest. %v", p)
			return
		}

		input := m.SignDigestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DigestSignerMockSignDigestInput{p}, "DigestSigner.SignDigest got unexpected parameters")

		result := m.SignDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestSignerMock.SignDigest")
			return
		}

		r = result.r

		return
	}

	if m.SignDigestMock.mainExpectation != nil {

		input := m.SignDigestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DigestSignerMockSignDigestInput{p}, "DigestSigner.SignDigest got unexpected parameters")
		}

		result := m.SignDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestSignerMock.SignDigest")
		}

		r = result.r

		return
	}

	if m.SignDigestFunc == nil {
		m.t.Fatalf("Unexpected call to DigestSignerMock.SignDigest. %v", p)
		return
	}

	return m.SignDigestFunc(p)
}

//SignDigestMinimockCounter returns a count of DigestSignerMock.SignDigestFunc invocations
func (m *DigestSignerMock) SignDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignDigestCounter)
}

//SignDigestMinimockPreCounter returns the value of DigestSignerMock.SignDigest invocations
func (m *DigestSignerMock) SignDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignDigestPreCounter)
}

//SignDigestFinished returns true if mock invocations count is ok
func (m *DigestSignerMock) SignDigestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SignDigestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SignDigestCounter) == uint64(len(m.SignDigestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SignDigestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SignDigestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SignDigestFunc != nil {
		return atomic.LoadUint64(&m.SignDigestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DigestSignerMock) ValidateCallCounters() {

	if !m.GetSignMethodFinished() {
		m.t.Fatal("Expected call to DigestSignerMock.GetSignMethod")
	}

	if !m.SignDigestFinished() {
		m.t.Fatal("Expected call to DigestSignerMock.SignDigest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DigestSignerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DigestSignerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DigestSignerMock) MinimockFinish() {

	if !m.GetSignMethodFinished() {
		m.t.Fatal("Expected call to DigestSignerMock.GetSignMethod")
	}

	if !m.SignDigestFinished() {
		m.t.Fatal("Expected call to DigestSignerMock.SignDigest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DigestSignerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DigestSignerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetSignMethodFinished()
		ok = ok && m.SignDigestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetSignMethodFinished() {
				m.t.Error("Expected call to DigestSignerMock.GetSignMethod")
			}

			if !m.SignDigestFinished() {
				m.t.Error("Expected call to DigestSignerMock.SignDigest")
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
func (m *DigestSignerMock) AllMocksCalled() bool {

	if !m.GetSignMethodFinished() {
		return false
	}

	if !m.SignDigestFinished() {
		return false
	}

	return true
}

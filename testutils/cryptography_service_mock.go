package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CryptographyService" can be found in github.com/insolar/insolar/core
*/
import (
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//CryptographyServiceMock implements github.com/insolar/insolar/core.CryptographyService
type CryptographyServiceMock struct {
	t minimock.Tester

	GetPublicKeyFunc       func() (r crypto.PublicKey, r1 error)
	GetPublicKeyCounter    uint64
	GetPublicKeyPreCounter uint64
	GetPublicKeyMock       mCryptographyServiceMockGetPublicKey

	SignFunc       func(p []byte) (r *core.Signature, r1 error)
	SignCounter    uint64
	SignPreCounter uint64
	SignMock       mCryptographyServiceMockSign

	VerifyFunc       func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool)
	VerifyCounter    uint64
	VerifyPreCounter uint64
	VerifyMock       mCryptographyServiceMockVerify
}

//NewCryptographyServiceMock returns a mock for github.com/insolar/insolar/core.CryptographyService
func NewCryptographyServiceMock(t minimock.Tester) *CryptographyServiceMock {
	m := &CryptographyServiceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetPublicKeyMock = mCryptographyServiceMockGetPublicKey{mock: m}
	m.SignMock = mCryptographyServiceMockSign{mock: m}
	m.VerifyMock = mCryptographyServiceMockVerify{mock: m}

	return m
}

type mCryptographyServiceMockGetPublicKey struct {
	mock              *CryptographyServiceMock
	mainExpectation   *CryptographyServiceMockGetPublicKeyExpectation
	expectationSeries []*CryptographyServiceMockGetPublicKeyExpectation
}

type CryptographyServiceMockGetPublicKeyExpectation struct {
	result *CryptographyServiceMockGetPublicKeyResult
}

type CryptographyServiceMockGetPublicKeyResult struct {
	r  crypto.PublicKey
	r1 error
}

//Expect specifies that invocation of CryptographyService.GetPublicKey is expected from 1 to Infinity times
func (m *mCryptographyServiceMockGetPublicKey) Expect() *mCryptographyServiceMockGetPublicKey {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CryptographyServiceMockGetPublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of CryptographyService.GetPublicKey
func (m *mCryptographyServiceMockGetPublicKey) Return(r crypto.PublicKey, r1 error) *CryptographyServiceMock {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CryptographyServiceMockGetPublicKeyExpectation{}
	}
	m.mainExpectation.result = &CryptographyServiceMockGetPublicKeyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CryptographyService.GetPublicKey is expected once
func (m *mCryptographyServiceMockGetPublicKey) ExpectOnce() *CryptographyServiceMockGetPublicKeyExpectation {
	m.mock.GetPublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &CryptographyServiceMockGetPublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CryptographyServiceMockGetPublicKeyExpectation) Return(r crypto.PublicKey, r1 error) {
	e.result = &CryptographyServiceMockGetPublicKeyResult{r, r1}
}

//Set uses given function f as a mock of CryptographyService.GetPublicKey method
func (m *mCryptographyServiceMockGetPublicKey) Set(f func() (r crypto.PublicKey, r1 error)) *CryptographyServiceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyFunc = f
	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/core.CryptographyService interface
func (m *CryptographyServiceMock) GetPublicKey() (r crypto.PublicKey, r1 error) {
	counter := atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if len(m.GetPublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CryptographyServiceMock.GetPublicKey.")
			return
		}

		result := m.GetPublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CryptographyServiceMock.GetPublicKey")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPublicKeyMock.mainExpectation != nil {

		result := m.GetPublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CryptographyServiceMock.GetPublicKey")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to CryptographyServiceMock.GetPublicKey.")
		return
	}

	return m.GetPublicKeyFunc()
}

//GetPublicKeyMinimockCounter returns a count of CryptographyServiceMock.GetPublicKeyFunc invocations
func (m *CryptographyServiceMock) GetPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyCounter)
}

//GetPublicKeyMinimockPreCounter returns the value of CryptographyServiceMock.GetPublicKey invocations
func (m *CryptographyServiceMock) GetPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyPreCounter)
}

//GetPublicKeyFinished returns true if mock invocations count is ok
func (m *CryptographyServiceMock) GetPublicKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPublicKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPublicKeyCounter) == uint64(len(m.GetPublicKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPublicKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPublicKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPublicKeyFunc != nil {
		return atomic.LoadUint64(&m.GetPublicKeyCounter) > 0
	}

	return true
}

type mCryptographyServiceMockSign struct {
	mock              *CryptographyServiceMock
	mainExpectation   *CryptographyServiceMockSignExpectation
	expectationSeries []*CryptographyServiceMockSignExpectation
}

type CryptographyServiceMockSignExpectation struct {
	input  *CryptographyServiceMockSignInput
	result *CryptographyServiceMockSignResult
}

type CryptographyServiceMockSignInput struct {
	p []byte
}

type CryptographyServiceMockSignResult struct {
	r  *core.Signature
	r1 error
}

//Expect specifies that invocation of CryptographyService.Sign is expected from 1 to Infinity times
func (m *mCryptographyServiceMockSign) Expect(p []byte) *mCryptographyServiceMockSign {
	m.mock.SignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CryptographyServiceMockSignExpectation{}
	}
	m.mainExpectation.input = &CryptographyServiceMockSignInput{p}
	return m
}

//Return specifies results of invocation of CryptographyService.Sign
func (m *mCryptographyServiceMockSign) Return(r *core.Signature, r1 error) *CryptographyServiceMock {
	m.mock.SignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CryptographyServiceMockSignExpectation{}
	}
	m.mainExpectation.result = &CryptographyServiceMockSignResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CryptographyService.Sign is expected once
func (m *mCryptographyServiceMockSign) ExpectOnce(p []byte) *CryptographyServiceMockSignExpectation {
	m.mock.SignFunc = nil
	m.mainExpectation = nil

	expectation := &CryptographyServiceMockSignExpectation{}
	expectation.input = &CryptographyServiceMockSignInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CryptographyServiceMockSignExpectation) Return(r *core.Signature, r1 error) {
	e.result = &CryptographyServiceMockSignResult{r, r1}
}

//Set uses given function f as a mock of CryptographyService.Sign method
func (m *mCryptographyServiceMockSign) Set(f func(p []byte) (r *core.Signature, r1 error)) *CryptographyServiceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignFunc = f
	return m.mock
}

//Sign implements github.com/insolar/insolar/core.CryptographyService interface
func (m *CryptographyServiceMock) Sign(p []byte) (r *core.Signature, r1 error) {
	counter := atomic.AddUint64(&m.SignPreCounter, 1)
	defer atomic.AddUint64(&m.SignCounter, 1)

	if len(m.SignMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CryptographyServiceMock.Sign. %v", p)
			return
		}

		input := m.SignMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CryptographyServiceMockSignInput{p}, "CryptographyService.Sign got unexpected parameters")

		result := m.SignMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CryptographyServiceMock.Sign")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SignMock.mainExpectation != nil {

		input := m.SignMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CryptographyServiceMockSignInput{p}, "CryptographyService.Sign got unexpected parameters")
		}

		result := m.SignMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CryptographyServiceMock.Sign")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SignFunc == nil {
		m.t.Fatalf("Unexpected call to CryptographyServiceMock.Sign. %v", p)
		return
	}

	return m.SignFunc(p)
}

//SignMinimockCounter returns a count of CryptographyServiceMock.SignFunc invocations
func (m *CryptographyServiceMock) SignMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignCounter)
}

//SignMinimockPreCounter returns the value of CryptographyServiceMock.Sign invocations
func (m *CryptographyServiceMock) SignMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignPreCounter)
}

//SignFinished returns true if mock invocations count is ok
func (m *CryptographyServiceMock) SignFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SignMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SignCounter) == uint64(len(m.SignMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SignMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SignCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SignFunc != nil {
		return atomic.LoadUint64(&m.SignCounter) > 0
	}

	return true
}

type mCryptographyServiceMockVerify struct {
	mock              *CryptographyServiceMock
	mainExpectation   *CryptographyServiceMockVerifyExpectation
	expectationSeries []*CryptographyServiceMockVerifyExpectation
}

type CryptographyServiceMockVerifyExpectation struct {
	input  *CryptographyServiceMockVerifyInput
	result *CryptographyServiceMockVerifyResult
}

type CryptographyServiceMockVerifyInput struct {
	p  crypto.PublicKey
	p1 core.Signature
	p2 []byte
}

type CryptographyServiceMockVerifyResult struct {
	r bool
}

//Expect specifies that invocation of CryptographyService.Verify is expected from 1 to Infinity times
func (m *mCryptographyServiceMockVerify) Expect(p crypto.PublicKey, p1 core.Signature, p2 []byte) *mCryptographyServiceMockVerify {
	m.mock.VerifyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CryptographyServiceMockVerifyExpectation{}
	}
	m.mainExpectation.input = &CryptographyServiceMockVerifyInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of CryptographyService.Verify
func (m *mCryptographyServiceMockVerify) Return(r bool) *CryptographyServiceMock {
	m.mock.VerifyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CryptographyServiceMockVerifyExpectation{}
	}
	m.mainExpectation.result = &CryptographyServiceMockVerifyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CryptographyService.Verify is expected once
func (m *mCryptographyServiceMockVerify) ExpectOnce(p crypto.PublicKey, p1 core.Signature, p2 []byte) *CryptographyServiceMockVerifyExpectation {
	m.mock.VerifyFunc = nil
	m.mainExpectation = nil

	expectation := &CryptographyServiceMockVerifyExpectation{}
	expectation.input = &CryptographyServiceMockVerifyInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CryptographyServiceMockVerifyExpectation) Return(r bool) {
	e.result = &CryptographyServiceMockVerifyResult{r}
}

//Set uses given function f as a mock of CryptographyService.Verify method
func (m *mCryptographyServiceMockVerify) Set(f func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool)) *CryptographyServiceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VerifyFunc = f
	return m.mock
}

//Verify implements github.com/insolar/insolar/core.CryptographyService interface
func (m *CryptographyServiceMock) Verify(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool) {
	counter := atomic.AddUint64(&m.VerifyPreCounter, 1)
	defer atomic.AddUint64(&m.VerifyCounter, 1)

	if len(m.VerifyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VerifyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CryptographyServiceMock.Verify. %v %v %v", p, p1, p2)
			return
		}

		input := m.VerifyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CryptographyServiceMockVerifyInput{p, p1, p2}, "CryptographyService.Verify got unexpected parameters")

		result := m.VerifyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CryptographyServiceMock.Verify")
			return
		}

		r = result.r

		return
	}

	if m.VerifyMock.mainExpectation != nil {

		input := m.VerifyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CryptographyServiceMockVerifyInput{p, p1, p2}, "CryptographyService.Verify got unexpected parameters")
		}

		result := m.VerifyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CryptographyServiceMock.Verify")
		}

		r = result.r

		return
	}

	if m.VerifyFunc == nil {
		m.t.Fatalf("Unexpected call to CryptographyServiceMock.Verify. %v %v %v", p, p1, p2)
		return
	}

	return m.VerifyFunc(p, p1, p2)
}

//VerifyMinimockCounter returns a count of CryptographyServiceMock.VerifyFunc invocations
func (m *CryptographyServiceMock) VerifyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyCounter)
}

//VerifyMinimockPreCounter returns the value of CryptographyServiceMock.Verify invocations
func (m *CryptographyServiceMock) VerifyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyPreCounter)
}

//VerifyFinished returns true if mock invocations count is ok
func (m *CryptographyServiceMock) VerifyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.VerifyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.VerifyCounter) == uint64(len(m.VerifyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.VerifyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.VerifyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.VerifyFunc != nil {
		return atomic.LoadUint64(&m.VerifyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CryptographyServiceMock) ValidateCallCounters() {

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to CryptographyServiceMock.GetPublicKey")
	}

	if !m.SignFinished() {
		m.t.Fatal("Expected call to CryptographyServiceMock.Sign")
	}

	if !m.VerifyFinished() {
		m.t.Fatal("Expected call to CryptographyServiceMock.Verify")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CryptographyServiceMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CryptographyServiceMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CryptographyServiceMock) MinimockFinish() {

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to CryptographyServiceMock.GetPublicKey")
	}

	if !m.SignFinished() {
		m.t.Fatal("Expected call to CryptographyServiceMock.Sign")
	}

	if !m.VerifyFinished() {
		m.t.Fatal("Expected call to CryptographyServiceMock.Verify")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CryptographyServiceMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CryptographyServiceMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetPublicKeyFinished()
		ok = ok && m.SignFinished()
		ok = ok && m.VerifyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetPublicKeyFinished() {
				m.t.Error("Expected call to CryptographyServiceMock.GetPublicKey")
			}

			if !m.SignFinished() {
				m.t.Error("Expected call to CryptographyServiceMock.Sign")
			}

			if !m.VerifyFinished() {
				m.t.Error("Expected call to CryptographyServiceMock.Verify")
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
func (m *CryptographyServiceMock) AllMocksCalled() bool {

	if !m.GetPublicKeyFinished() {
		return false
	}

	if !m.SignFinished() {
		return false
	}

	if !m.VerifyFinished() {
		return false
	}

	return true
}

package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CertificateHolder" can be found in github.com/insolar/insolar/network/consensus/common/cryptkit
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//CertificateHolderMock implements github.com/insolar/insolar/network/consensus/common/cryptkit.CertificateHolder
type CertificateHolderMock struct {
	t minimock.Tester

	GetPublicKeyFunc       func() (r SignatureKeyHolder)
	GetPublicKeyCounter    uint64
	GetPublicKeyPreCounter uint64
	GetPublicKeyMock       mCertificateHolderMockGetPublicKey

	IsValidForHostAddressFunc       func(p string) (r bool)
	IsValidForHostAddressCounter    uint64
	IsValidForHostAddressPreCounter uint64
	IsValidForHostAddressMock       mCertificateHolderMockIsValidForHostAddress
}

//NewCertificateHolderMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptkit.CertificateHolder
func NewCertificateHolderMock(t minimock.Tester) *CertificateHolderMock {
	m := &CertificateHolderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetPublicKeyMock = mCertificateHolderMockGetPublicKey{mock: m}
	m.IsValidForHostAddressMock = mCertificateHolderMockIsValidForHostAddress{mock: m}

	return m
}

type mCertificateHolderMockGetPublicKey struct {
	mock              *CertificateHolderMock
	mainExpectation   *CertificateHolderMockGetPublicKeyExpectation
	expectationSeries []*CertificateHolderMockGetPublicKeyExpectation
}

type CertificateHolderMockGetPublicKeyExpectation struct {
	result *CertificateHolderMockGetPublicKeyResult
}

type CertificateHolderMockGetPublicKeyResult struct {
	r SignatureKeyHolder
}

//Expect specifies that invocation of CertificateHolder.GetPublicKey is expected from 1 to Infinity times
func (m *mCertificateHolderMockGetPublicKey) Expect() *mCertificateHolderMockGetPublicKey {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateHolderMockGetPublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of CertificateHolder.GetPublicKey
func (m *mCertificateHolderMockGetPublicKey) Return(r SignatureKeyHolder) *CertificateHolderMock {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateHolderMockGetPublicKeyExpectation{}
	}
	m.mainExpectation.result = &CertificateHolderMockGetPublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CertificateHolder.GetPublicKey is expected once
func (m *mCertificateHolderMockGetPublicKey) ExpectOnce() *CertificateHolderMockGetPublicKeyExpectation {
	m.mock.GetPublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateHolderMockGetPublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateHolderMockGetPublicKeyExpectation) Return(r SignatureKeyHolder) {
	e.result = &CertificateHolderMockGetPublicKeyResult{r}
}

//Set uses given function f as a mock of CertificateHolder.GetPublicKey method
func (m *mCertificateHolderMockGetPublicKey) Set(f func() (r SignatureKeyHolder)) *CertificateHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyFunc = f
	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/network/consensus/common/cryptkit.CertificateHolder interface
func (m *CertificateHolderMock) GetPublicKey() (r SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if len(m.GetPublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateHolderMock.GetPublicKey.")
			return
		}

		result := m.GetPublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateHolderMock.GetPublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyMock.mainExpectation != nil {

		result := m.GetPublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateHolderMock.GetPublicKey")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateHolderMock.GetPublicKey.")
		return
	}

	return m.GetPublicKeyFunc()
}

//GetPublicKeyMinimockCounter returns a count of CertificateHolderMock.GetPublicKeyFunc invocations
func (m *CertificateHolderMock) GetPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyCounter)
}

//GetPublicKeyMinimockPreCounter returns the value of CertificateHolderMock.GetPublicKey invocations
func (m *CertificateHolderMock) GetPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyPreCounter)
}

//GetPublicKeyFinished returns true if mock invocations count is ok
func (m *CertificateHolderMock) GetPublicKeyFinished() bool {
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

type mCertificateHolderMockIsValidForHostAddress struct {
	mock              *CertificateHolderMock
	mainExpectation   *CertificateHolderMockIsValidForHostAddressExpectation
	expectationSeries []*CertificateHolderMockIsValidForHostAddressExpectation
}

type CertificateHolderMockIsValidForHostAddressExpectation struct {
	input  *CertificateHolderMockIsValidForHostAddressInput
	result *CertificateHolderMockIsValidForHostAddressResult
}

type CertificateHolderMockIsValidForHostAddressInput struct {
	p string
}

type CertificateHolderMockIsValidForHostAddressResult struct {
	r bool
}

//Expect specifies that invocation of CertificateHolder.IsValidForHostAddress is expected from 1 to Infinity times
func (m *mCertificateHolderMockIsValidForHostAddress) Expect(p string) *mCertificateHolderMockIsValidForHostAddress {
	m.mock.IsValidForHostAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateHolderMockIsValidForHostAddressExpectation{}
	}
	m.mainExpectation.input = &CertificateHolderMockIsValidForHostAddressInput{p}
	return m
}

//Return specifies results of invocation of CertificateHolder.IsValidForHostAddress
func (m *mCertificateHolderMockIsValidForHostAddress) Return(r bool) *CertificateHolderMock {
	m.mock.IsValidForHostAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateHolderMockIsValidForHostAddressExpectation{}
	}
	m.mainExpectation.result = &CertificateHolderMockIsValidForHostAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CertificateHolder.IsValidForHostAddress is expected once
func (m *mCertificateHolderMockIsValidForHostAddress) ExpectOnce(p string) *CertificateHolderMockIsValidForHostAddressExpectation {
	m.mock.IsValidForHostAddressFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateHolderMockIsValidForHostAddressExpectation{}
	expectation.input = &CertificateHolderMockIsValidForHostAddressInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateHolderMockIsValidForHostAddressExpectation) Return(r bool) {
	e.result = &CertificateHolderMockIsValidForHostAddressResult{r}
}

//Set uses given function f as a mock of CertificateHolder.IsValidForHostAddress method
func (m *mCertificateHolderMockIsValidForHostAddress) Set(f func(p string) (r bool)) *CertificateHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsValidForHostAddressFunc = f
	return m.mock
}

//IsValidForHostAddress implements github.com/insolar/insolar/network/consensus/common/cryptkit.CertificateHolder interface
func (m *CertificateHolderMock) IsValidForHostAddress(p string) (r bool) {
	counter := atomic.AddUint64(&m.IsValidForHostAddressPreCounter, 1)
	defer atomic.AddUint64(&m.IsValidForHostAddressCounter, 1)

	if len(m.IsValidForHostAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsValidForHostAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateHolderMock.IsValidForHostAddress. %v", p)
			return
		}

		input := m.IsValidForHostAddressMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CertificateHolderMockIsValidForHostAddressInput{p}, "CertificateHolder.IsValidForHostAddress got unexpected parameters")

		result := m.IsValidForHostAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateHolderMock.IsValidForHostAddress")
			return
		}

		r = result.r

		return
	}

	if m.IsValidForHostAddressMock.mainExpectation != nil {

		input := m.IsValidForHostAddressMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CertificateHolderMockIsValidForHostAddressInput{p}, "CertificateHolder.IsValidForHostAddress got unexpected parameters")
		}

		result := m.IsValidForHostAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateHolderMock.IsValidForHostAddress")
		}

		r = result.r

		return
	}

	if m.IsValidForHostAddressFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateHolderMock.IsValidForHostAddress. %v", p)
		return
	}

	return m.IsValidForHostAddressFunc(p)
}

//IsValidForHostAddressMinimockCounter returns a count of CertificateHolderMock.IsValidForHostAddressFunc invocations
func (m *CertificateHolderMock) IsValidForHostAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidForHostAddressCounter)
}

//IsValidForHostAddressMinimockPreCounter returns the value of CertificateHolderMock.IsValidForHostAddress invocations
func (m *CertificateHolderMock) IsValidForHostAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidForHostAddressPreCounter)
}

//IsValidForHostAddressFinished returns true if mock invocations count is ok
func (m *CertificateHolderMock) IsValidForHostAddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsValidForHostAddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsValidForHostAddressCounter) == uint64(len(m.IsValidForHostAddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsValidForHostAddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsValidForHostAddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsValidForHostAddressFunc != nil {
		return atomic.LoadUint64(&m.IsValidForHostAddressCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateHolderMock) ValidateCallCounters() {

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to CertificateHolderMock.GetPublicKey")
	}

	if !m.IsValidForHostAddressFinished() {
		m.t.Fatal("Expected call to CertificateHolderMock.IsValidForHostAddress")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateHolderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CertificateHolderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CertificateHolderMock) MinimockFinish() {

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to CertificateHolderMock.GetPublicKey")
	}

	if !m.IsValidForHostAddressFinished() {
		m.t.Fatal("Expected call to CertificateHolderMock.IsValidForHostAddress")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CertificateHolderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CertificateHolderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetPublicKeyFinished()
		ok = ok && m.IsValidForHostAddressFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetPublicKeyFinished() {
				m.t.Error("Expected call to CertificateHolderMock.GetPublicKey")
			}

			if !m.IsValidForHostAddressFinished() {
				m.t.Error("Expected call to CertificateHolderMock.IsValidForHostAddress")
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
func (m *CertificateHolderMock) AllMocksCalled() bool {

	if !m.GetPublicKeyFinished() {
		return false
	}

	if !m.IsValidForHostAddressFinished() {
		return false
	}

	return true
}

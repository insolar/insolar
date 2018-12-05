package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CertificateManager" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
)

//CertificateManagerMock implements github.com/insolar/insolar/core.CertificateManager
type CertificateManagerMock struct {
	t minimock.Tester

	GetCertificateFunc       func() (r core.Certificate)
	GetCertificateCounter    uint64
	GetCertificatePreCounter uint64
	GetCertificateMock       mCertificateManagerMockGetCertificate
}

//NewCertificateManagerMock returns a mock for github.com/insolar/insolar/core.CertificateManager
func NewCertificateManagerMock(t minimock.Tester) *CertificateManagerMock {
	m := &CertificateManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCertificateMock = mCertificateManagerMockGetCertificate{mock: m}

	return m
}

type mCertificateManagerMockGetCertificate struct {
	mock              *CertificateManagerMock
	mainExpectation   *CertificateManagerMockGetCertificateExpectation
	expectationSeries []*CertificateManagerMockGetCertificateExpectation
}

type CertificateManagerMockGetCertificateExpectation struct {
	result *CertificateManagerMockGetCertificateResult
}

type CertificateManagerMockGetCertificateResult struct {
	r core.Certificate
}

//Expect specifies that invocation of CertificateManager.GetCertificate is expected from 1 to Infinity times
func (m *mCertificateManagerMockGetCertificate) Expect() *mCertificateManagerMockGetCertificate {
	m.mock.GetCertificateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateManagerMockGetCertificateExpectation{}
	}

	return m
}

//Return specifies results of invocation of CertificateManager.GetCertificate
func (m *mCertificateManagerMockGetCertificate) Return(r core.Certificate) *CertificateManagerMock {
	m.mock.GetCertificateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateManagerMockGetCertificateExpectation{}
	}
	m.mainExpectation.result = &CertificateManagerMockGetCertificateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CertificateManager.GetCertificate is expected once
func (m *mCertificateManagerMockGetCertificate) ExpectOnce() *CertificateManagerMockGetCertificateExpectation {
	m.mock.GetCertificateFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateManagerMockGetCertificateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateManagerMockGetCertificateExpectation) Return(r core.Certificate) {
	e.result = &CertificateManagerMockGetCertificateResult{r}
}

//Set uses given function f as a mock of CertificateManager.GetCertificate method
func (m *mCertificateManagerMockGetCertificate) Set(f func() (r core.Certificate)) *CertificateManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCertificateFunc = f
	return m.mock
}

//GetCertificate implements github.com/insolar/insolar/core.CertificateManager interface
func (m *CertificateManagerMock) GetCertificate() (r core.Certificate) {
	counter := atomic.AddUint64(&m.GetCertificatePreCounter, 1)
	defer atomic.AddUint64(&m.GetCertificateCounter, 1)

	if len(m.GetCertificateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCertificateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateManagerMock.GetCertificate.")
			return
		}

		result := m.GetCertificateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateManagerMock.GetCertificate")
			return
		}

		r = result.r

		return
	}

	if m.GetCertificateMock.mainExpectation != nil {

		result := m.GetCertificateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateManagerMock.GetCertificate")
		}

		r = result.r

		return
	}

	if m.GetCertificateFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateManagerMock.GetCertificate.")
		return
	}

	return m.GetCertificateFunc()
}

//GetCertificateMinimockCounter returns a count of CertificateManagerMock.GetCertificateFunc invocations
func (m *CertificateManagerMock) GetCertificateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCertificateCounter)
}

//GetCertificateMinimockPreCounter returns the value of CertificateManagerMock.GetCertificate invocations
func (m *CertificateManagerMock) GetCertificateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCertificatePreCounter)
}

//GetCertificateFinished returns true if mock invocations count is ok
func (m *CertificateManagerMock) GetCertificateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCertificateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCertificateCounter) == uint64(len(m.GetCertificateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCertificateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCertificateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCertificateFunc != nil {
		return atomic.LoadUint64(&m.GetCertificateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateManagerMock) ValidateCallCounters() {

	if !m.GetCertificateFinished() {
		m.t.Fatal("Expected call to CertificateManagerMock.GetCertificate")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CertificateManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CertificateManagerMock) MinimockFinish() {

	if !m.GetCertificateFinished() {
		m.t.Fatal("Expected call to CertificateManagerMock.GetCertificate")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CertificateManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CertificateManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetCertificateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetCertificateFinished() {
				m.t.Error("Expected call to CertificateManagerMock.GetCertificate")
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
func (m *CertificateManagerMock) AllMocksCalled() bool {

	if !m.GetCertificateFinished() {
		return false
	}

	return true
}

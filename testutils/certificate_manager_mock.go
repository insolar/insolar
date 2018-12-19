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

	testify_assert "github.com/stretchr/testify/assert"
)

//CertificateManagerMock implements github.com/insolar/insolar/core.CertificateManager
type CertificateManagerMock struct {
	t minimock.Tester

	GetCertificateFunc       func() (r core.Certificate)
	GetCertificateCounter    uint64
	GetCertificatePreCounter uint64
	GetCertificateMock       mCertificateManagerMockGetCertificate

	NewUnsignedCertificateFunc       func(p string, p1 string, p2 string) (r core.Certificate, r1 error)
	NewUnsignedCertificateCounter    uint64
	NewUnsignedCertificatePreCounter uint64
	NewUnsignedCertificateMock       mCertificateManagerMockNewUnsignedCertificate

	VerifyAuthorizationCertificateFunc       func(p core.AuthorizationCertificate) (r bool, r1 error)
	VerifyAuthorizationCertificateCounter    uint64
	VerifyAuthorizationCertificatePreCounter uint64
	VerifyAuthorizationCertificateMock       mCertificateManagerMockVerifyAuthorizationCertificate
}

//NewCertificateManagerMock returns a mock for github.com/insolar/insolar/core.CertificateManager
func NewCertificateManagerMock(t minimock.Tester) *CertificateManagerMock {
	m := &CertificateManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCertificateMock = mCertificateManagerMockGetCertificate{mock: m}
	m.NewUnsignedCertificateMock = mCertificateManagerMockNewUnsignedCertificate{mock: m}
	m.VerifyAuthorizationCertificateMock = mCertificateManagerMockVerifyAuthorizationCertificate{mock: m}

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

type mCertificateManagerMockNewUnsignedCertificate struct {
	mock              *CertificateManagerMock
	mainExpectation   *CertificateManagerMockNewUnsignedCertificateExpectation
	expectationSeries []*CertificateManagerMockNewUnsignedCertificateExpectation
}

type CertificateManagerMockNewUnsignedCertificateExpectation struct {
	input  *CertificateManagerMockNewUnsignedCertificateInput
	result *CertificateManagerMockNewUnsignedCertificateResult
}

type CertificateManagerMockNewUnsignedCertificateInput struct {
	p  string
	p1 string
	p2 string
}

type CertificateManagerMockNewUnsignedCertificateResult struct {
	r  core.Certificate
	r1 error
}

//Expect specifies that invocation of CertificateManager.NewUnsignedCertificate is expected from 1 to Infinity times
func (m *mCertificateManagerMockNewUnsignedCertificate) Expect(p string, p1 string, p2 string) *mCertificateManagerMockNewUnsignedCertificate {
	m.mock.NewUnsignedCertificateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateManagerMockNewUnsignedCertificateExpectation{}
	}
	m.mainExpectation.input = &CertificateManagerMockNewUnsignedCertificateInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of CertificateManager.NewUnsignedCertificate
func (m *mCertificateManagerMockNewUnsignedCertificate) Return(r core.Certificate, r1 error) *CertificateManagerMock {
	m.mock.NewUnsignedCertificateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateManagerMockNewUnsignedCertificateExpectation{}
	}
	m.mainExpectation.result = &CertificateManagerMockNewUnsignedCertificateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CertificateManager.NewUnsignedCertificate is expected once
func (m *mCertificateManagerMockNewUnsignedCertificate) ExpectOnce(p string, p1 string, p2 string) *CertificateManagerMockNewUnsignedCertificateExpectation {
	m.mock.NewUnsignedCertificateFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateManagerMockNewUnsignedCertificateExpectation{}
	expectation.input = &CertificateManagerMockNewUnsignedCertificateInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateManagerMockNewUnsignedCertificateExpectation) Return(r core.Certificate, r1 error) {
	e.result = &CertificateManagerMockNewUnsignedCertificateResult{r, r1}
}

//Set uses given function f as a mock of CertificateManager.NewUnsignedCertificate method
func (m *mCertificateManagerMockNewUnsignedCertificate) Set(f func(p string, p1 string, p2 string) (r core.Certificate, r1 error)) *CertificateManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NewUnsignedCertificateFunc = f
	return m.mock
}

//NewUnsignedCertificate implements github.com/insolar/insolar/core.CertificateManager interface
func (m *CertificateManagerMock) NewUnsignedCertificate(p string, p1 string, p2 string) (r core.Certificate, r1 error) {
	counter := atomic.AddUint64(&m.NewUnsignedCertificatePreCounter, 1)
	defer atomic.AddUint64(&m.NewUnsignedCertificateCounter, 1)

	if len(m.NewUnsignedCertificateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NewUnsignedCertificateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateManagerMock.NewUnsignedCertificate. %v %v %v", p, p1, p2)
			return
		}

		input := m.NewUnsignedCertificateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CertificateManagerMockNewUnsignedCertificateInput{p, p1, p2}, "CertificateManager.NewUnsignedCertificate got unexpected parameters")

		result := m.NewUnsignedCertificateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateManagerMock.NewUnsignedCertificate")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewUnsignedCertificateMock.mainExpectation != nil {

		input := m.NewUnsignedCertificateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CertificateManagerMockNewUnsignedCertificateInput{p, p1, p2}, "CertificateManager.NewUnsignedCertificate got unexpected parameters")
		}

		result := m.NewUnsignedCertificateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateManagerMock.NewUnsignedCertificate")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewUnsignedCertificateFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateManagerMock.NewUnsignedCertificate. %v %v %v", p, p1, p2)
		return
	}

	return m.NewUnsignedCertificateFunc(p, p1, p2)
}

//NewUnsignedCertificateMinimockCounter returns a count of CertificateManagerMock.NewUnsignedCertificateFunc invocations
func (m *CertificateManagerMock) NewUnsignedCertificateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewUnsignedCertificateCounter)
}

//NewUnsignedCertificateMinimockPreCounter returns the value of CertificateManagerMock.NewUnsignedCertificate invocations
func (m *CertificateManagerMock) NewUnsignedCertificateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewUnsignedCertificatePreCounter)
}

//NewUnsignedCertificateFinished returns true if mock invocations count is ok
func (m *CertificateManagerMock) NewUnsignedCertificateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NewUnsignedCertificateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NewUnsignedCertificateCounter) == uint64(len(m.NewUnsignedCertificateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NewUnsignedCertificateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NewUnsignedCertificateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NewUnsignedCertificateFunc != nil {
		return atomic.LoadUint64(&m.NewUnsignedCertificateCounter) > 0
	}

	return true
}

type mCertificateManagerMockVerifyAuthorizationCertificate struct {
	mock              *CertificateManagerMock
	mainExpectation   *CertificateManagerMockVerifyAuthorizationCertificateExpectation
	expectationSeries []*CertificateManagerMockVerifyAuthorizationCertificateExpectation
}

type CertificateManagerMockVerifyAuthorizationCertificateExpectation struct {
	input  *CertificateManagerMockVerifyAuthorizationCertificateInput
	result *CertificateManagerMockVerifyAuthorizationCertificateResult
}

type CertificateManagerMockVerifyAuthorizationCertificateInput struct {
	p core.AuthorizationCertificate
}

type CertificateManagerMockVerifyAuthorizationCertificateResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of CertificateManager.VerifyAuthorizationCertificate is expected from 1 to Infinity times
func (m *mCertificateManagerMockVerifyAuthorizationCertificate) Expect(p core.AuthorizationCertificate) *mCertificateManagerMockVerifyAuthorizationCertificate {
	m.mock.VerifyAuthorizationCertificateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateManagerMockVerifyAuthorizationCertificateExpectation{}
	}
	m.mainExpectation.input = &CertificateManagerMockVerifyAuthorizationCertificateInput{p}
	return m
}

//Return specifies results of invocation of CertificateManager.VerifyAuthorizationCertificate
func (m *mCertificateManagerMockVerifyAuthorizationCertificate) Return(r bool, r1 error) *CertificateManagerMock {
	m.mock.VerifyAuthorizationCertificateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateManagerMockVerifyAuthorizationCertificateExpectation{}
	}
	m.mainExpectation.result = &CertificateManagerMockVerifyAuthorizationCertificateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CertificateManager.VerifyAuthorizationCertificate is expected once
func (m *mCertificateManagerMockVerifyAuthorizationCertificate) ExpectOnce(p core.AuthorizationCertificate) *CertificateManagerMockVerifyAuthorizationCertificateExpectation {
	m.mock.VerifyAuthorizationCertificateFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateManagerMockVerifyAuthorizationCertificateExpectation{}
	expectation.input = &CertificateManagerMockVerifyAuthorizationCertificateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateManagerMockVerifyAuthorizationCertificateExpectation) Return(r bool, r1 error) {
	e.result = &CertificateManagerMockVerifyAuthorizationCertificateResult{r, r1}
}

//Set uses given function f as a mock of CertificateManager.VerifyAuthorizationCertificate method
func (m *mCertificateManagerMockVerifyAuthorizationCertificate) Set(f func(p core.AuthorizationCertificate) (r bool, r1 error)) *CertificateManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VerifyAuthorizationCertificateFunc = f
	return m.mock
}

//VerifyAuthorizationCertificate implements github.com/insolar/insolar/core.CertificateManager interface
func (m *CertificateManagerMock) VerifyAuthorizationCertificate(p core.AuthorizationCertificate) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.VerifyAuthorizationCertificatePreCounter, 1)
	defer atomic.AddUint64(&m.VerifyAuthorizationCertificateCounter, 1)

	if len(m.VerifyAuthorizationCertificateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VerifyAuthorizationCertificateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateManagerMock.VerifyAuthorizationCertificate. %v", p)
			return
		}

		input := m.VerifyAuthorizationCertificateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CertificateManagerMockVerifyAuthorizationCertificateInput{p}, "CertificateManager.VerifyAuthorizationCertificate got unexpected parameters")

		result := m.VerifyAuthorizationCertificateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateManagerMock.VerifyAuthorizationCertificate")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VerifyAuthorizationCertificateMock.mainExpectation != nil {

		input := m.VerifyAuthorizationCertificateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CertificateManagerMockVerifyAuthorizationCertificateInput{p}, "CertificateManager.VerifyAuthorizationCertificate got unexpected parameters")
		}

		result := m.VerifyAuthorizationCertificateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateManagerMock.VerifyAuthorizationCertificate")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VerifyAuthorizationCertificateFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateManagerMock.VerifyAuthorizationCertificate. %v", p)
		return
	}

	return m.VerifyAuthorizationCertificateFunc(p)
}

//VerifyAuthorizationCertificateMinimockCounter returns a count of CertificateManagerMock.VerifyAuthorizationCertificateFunc invocations
func (m *CertificateManagerMock) VerifyAuthorizationCertificateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter)
}

//VerifyAuthorizationCertificateMinimockPreCounter returns the value of CertificateManagerMock.VerifyAuthorizationCertificate invocations
func (m *CertificateManagerMock) VerifyAuthorizationCertificateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyAuthorizationCertificatePreCounter)
}

//VerifyAuthorizationCertificateFinished returns true if mock invocations count is ok
func (m *CertificateManagerMock) VerifyAuthorizationCertificateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.VerifyAuthorizationCertificateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == uint64(len(m.VerifyAuthorizationCertificateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.VerifyAuthorizationCertificateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.VerifyAuthorizationCertificateFunc != nil {
		return atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateManagerMock) ValidateCallCounters() {

	if !m.GetCertificateFinished() {
		m.t.Fatal("Expected call to CertificateManagerMock.GetCertificate")
	}

	if !m.NewUnsignedCertificateFinished() {
		m.t.Fatal("Expected call to CertificateManagerMock.NewUnsignedCertificate")
	}

	if !m.VerifyAuthorizationCertificateFinished() {
		m.t.Fatal("Expected call to CertificateManagerMock.VerifyAuthorizationCertificate")
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

	if !m.NewUnsignedCertificateFinished() {
		m.t.Fatal("Expected call to CertificateManagerMock.NewUnsignedCertificate")
	}

	if !m.VerifyAuthorizationCertificateFinished() {
		m.t.Fatal("Expected call to CertificateManagerMock.VerifyAuthorizationCertificate")
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
		ok = ok && m.NewUnsignedCertificateFinished()
		ok = ok && m.VerifyAuthorizationCertificateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetCertificateFinished() {
				m.t.Error("Expected call to CertificateManagerMock.GetCertificate")
			}

			if !m.NewUnsignedCertificateFinished() {
				m.t.Error("Expected call to CertificateManagerMock.NewUnsignedCertificate")
			}

			if !m.VerifyAuthorizationCertificateFinished() {
				m.t.Error("Expected call to CertificateManagerMock.VerifyAuthorizationCertificate")
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

	if !m.NewUnsignedCertificateFinished() {
		return false
	}

	if !m.VerifyAuthorizationCertificateFinished() {
		return false
	}

	return true
}

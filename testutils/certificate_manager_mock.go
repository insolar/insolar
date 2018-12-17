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
	mock *CertificateManagerMock
}

//Return sets up a mock for CertificateManager.GetCertificate to return Return's arguments
func (m *mCertificateManagerMockGetCertificate) Return(r core.Certificate) *CertificateManagerMock {
	m.mock.GetCertificateFunc = func() core.Certificate {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of CertificateManager.GetCertificate method
func (m *mCertificateManagerMockGetCertificate) Set(f func() (r core.Certificate)) *CertificateManagerMock {
	m.mock.GetCertificateFunc = f

	return m.mock
}

//GetCertificate implements github.com/insolar/insolar/core.CertificateManager interface
func (m *CertificateManagerMock) GetCertificate() (r core.Certificate) {
	atomic.AddUint64(&m.GetCertificatePreCounter, 1)
	defer atomic.AddUint64(&m.GetCertificateCounter, 1)

	if m.GetCertificateFunc == nil {
		m.t.Fatal("Unexpected call to CertificateManagerMock.GetCertificate")
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

type mCertificateManagerMockNewUnsignedCertificate struct {
	mock             *CertificateManagerMock
	mockExpectations *CertificateManagerMockNewUnsignedCertificateParams
}

//CertificateManagerMockNewUnsignedCertificateParams represents input parameters of the CertificateManager.NewUnsignedCertificate
type CertificateManagerMockNewUnsignedCertificateParams struct {
	p  string
	p1 string
	p2 string
}

//Expect sets up expected params for the CertificateManager.NewUnsignedCertificate
func (m *mCertificateManagerMockNewUnsignedCertificate) Expect(p string, p1 string, p2 string) *mCertificateManagerMockNewUnsignedCertificate {
	m.mockExpectations = &CertificateManagerMockNewUnsignedCertificateParams{p, p1, p2}
	return m
}

//Return sets up a mock for CertificateManager.NewUnsignedCertificate to return Return's arguments
func (m *mCertificateManagerMockNewUnsignedCertificate) Return(r core.Certificate, r1 error) *CertificateManagerMock {
	m.mock.NewUnsignedCertificateFunc = func(p string, p1 string, p2 string) (core.Certificate, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of CertificateManager.NewUnsignedCertificate method
func (m *mCertificateManagerMockNewUnsignedCertificate) Set(f func(p string, p1 string, p2 string) (r core.Certificate, r1 error)) *CertificateManagerMock {
	m.mock.NewUnsignedCertificateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NewUnsignedCertificate implements github.com/insolar/insolar/core.CertificateManager interface
func (m *CertificateManagerMock) NewUnsignedCertificate(p string, p1 string, p2 string) (r core.Certificate, r1 error) {
	atomic.AddUint64(&m.NewUnsignedCertificatePreCounter, 1)
	defer atomic.AddUint64(&m.NewUnsignedCertificateCounter, 1)

	if m.NewUnsignedCertificateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NewUnsignedCertificateMock.mockExpectations, CertificateManagerMockNewUnsignedCertificateParams{p, p1, p2},
			"CertificateManager.NewUnsignedCertificate got unexpected parameters")

		if m.NewUnsignedCertificateFunc == nil {

			m.t.Fatal("No results are set for the CertificateManagerMock.NewUnsignedCertificate")

			return
		}
	}

	if m.NewUnsignedCertificateFunc == nil {
		m.t.Fatal("Unexpected call to CertificateManagerMock.NewUnsignedCertificate")
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

type mCertificateManagerMockVerifyAuthorizationCertificate struct {
	mock             *CertificateManagerMock
	mockExpectations *CertificateManagerMockVerifyAuthorizationCertificateParams
}

//CertificateManagerMockVerifyAuthorizationCertificateParams represents input parameters of the CertificateManager.VerifyAuthorizationCertificate
type CertificateManagerMockVerifyAuthorizationCertificateParams struct {
	p core.AuthorizationCertificate
}

//Expect sets up expected params for the CertificateManager.VerifyAuthorizationCertificate
func (m *mCertificateManagerMockVerifyAuthorizationCertificate) Expect(p core.AuthorizationCertificate) *mCertificateManagerMockVerifyAuthorizationCertificate {
	m.mockExpectations = &CertificateManagerMockVerifyAuthorizationCertificateParams{p}
	return m
}

//Return sets up a mock for CertificateManager.VerifyAuthorizationCertificate to return Return's arguments
func (m *mCertificateManagerMockVerifyAuthorizationCertificate) Return(r bool, r1 error) *CertificateManagerMock {
	m.mock.VerifyAuthorizationCertificateFunc = func(p core.AuthorizationCertificate) (bool, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of CertificateManager.VerifyAuthorizationCertificate method
func (m *mCertificateManagerMockVerifyAuthorizationCertificate) Set(f func(p core.AuthorizationCertificate) (r bool, r1 error)) *CertificateManagerMock {
	m.mock.VerifyAuthorizationCertificateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//VerifyAuthorizationCertificate implements github.com/insolar/insolar/core.CertificateManager interface
func (m *CertificateManagerMock) VerifyAuthorizationCertificate(p core.AuthorizationCertificate) (r bool, r1 error) {
	atomic.AddUint64(&m.VerifyAuthorizationCertificatePreCounter, 1)
	defer atomic.AddUint64(&m.VerifyAuthorizationCertificateCounter, 1)

	if m.VerifyAuthorizationCertificateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.VerifyAuthorizationCertificateMock.mockExpectations, CertificateManagerMockVerifyAuthorizationCertificateParams{p},
			"CertificateManager.VerifyAuthorizationCertificate got unexpected parameters")

		if m.VerifyAuthorizationCertificateFunc == nil {

			m.t.Fatal("No results are set for the CertificateManagerMock.VerifyAuthorizationCertificate")

			return
		}
	}

	if m.VerifyAuthorizationCertificateFunc == nil {
		m.t.Fatal("Unexpected call to CertificateManagerMock.VerifyAuthorizationCertificate")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateManagerMock) ValidateCallCounters() {

	if m.GetCertificateFunc != nil && atomic.LoadUint64(&m.GetCertificateCounter) == 0 {
		m.t.Fatal("Expected call to CertificateManagerMock.GetCertificate")
	}

	if m.NewUnsignedCertificateFunc != nil && atomic.LoadUint64(&m.NewUnsignedCertificateCounter) == 0 {
		m.t.Fatal("Expected call to CertificateManagerMock.NewUnsignedCertificate")
	}

	if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
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

	if m.GetCertificateFunc != nil && atomic.LoadUint64(&m.GetCertificateCounter) == 0 {
		m.t.Fatal("Expected call to CertificateManagerMock.GetCertificate")
	}

	if m.NewUnsignedCertificateFunc != nil && atomic.LoadUint64(&m.NewUnsignedCertificateCounter) == 0 {
		m.t.Fatal("Expected call to CertificateManagerMock.NewUnsignedCertificate")
	}

	if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
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
		ok = ok && (m.GetCertificateFunc == nil || atomic.LoadUint64(&m.GetCertificateCounter) > 0)
		ok = ok && (m.NewUnsignedCertificateFunc == nil || atomic.LoadUint64(&m.NewUnsignedCertificateCounter) > 0)
		ok = ok && (m.VerifyAuthorizationCertificateFunc == nil || atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetCertificateFunc != nil && atomic.LoadUint64(&m.GetCertificateCounter) == 0 {
				m.t.Error("Expected call to CertificateManagerMock.GetCertificate")
			}

			if m.NewUnsignedCertificateFunc != nil && atomic.LoadUint64(&m.NewUnsignedCertificateCounter) == 0 {
				m.t.Error("Expected call to CertificateManagerMock.NewUnsignedCertificate")
			}

			if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
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

	if m.GetCertificateFunc != nil && atomic.LoadUint64(&m.GetCertificateCounter) == 0 {
		return false
	}

	if m.NewUnsignedCertificateFunc != nil && atomic.LoadUint64(&m.NewUnsignedCertificateCounter) == 0 {
		return false
	}

	if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
		return false
	}

	return true
}

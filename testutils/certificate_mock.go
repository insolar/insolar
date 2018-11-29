package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Certificate" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//CertificateMock implements github.com/insolar/insolar/core.Certificate
type CertificateMock struct {
	t minimock.Tester

	GetBootstrapNodesFunc       func() (r []core.BootstrapNode)
	GetBootstrapNodesCounter    uint64
	GetBootstrapNodesPreCounter uint64
	GetBootstrapNodesMock       mCertificateMockGetBootstrapNodes

	GetRoleFunc       func() (r core.NodeRole)
	GetRoleCounter    uint64
	GetRolePreCounter uint64
	GetRoleMock       mCertificateMockGetRole

	GetRootDomainReferenceFunc       func() (r *core.RecordRef)
	GetRootDomainReferenceCounter    uint64
	GetRootDomainReferencePreCounter uint64
	GetRootDomainReferenceMock       mCertificateMockGetRootDomainReference

	SetRootDomainReferenceFunc       func(p *core.RecordRef)
	SetRootDomainReferenceCounter    uint64
	SetRootDomainReferencePreCounter uint64
	SetRootDomainReferenceMock       mCertificateMockSetRootDomainReference
}

//NewCertificateMock returns a mock for github.com/insolar/insolar/core.Certificate
func NewCertificateMock(t minimock.Tester) *CertificateMock {
	m := &CertificateMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetBootstrapNodesMock = mCertificateMockGetBootstrapNodes{mock: m}
	m.GetRoleMock = mCertificateMockGetRole{mock: m}
	m.GetRootDomainReferenceMock = mCertificateMockGetRootDomainReference{mock: m}
	m.SetRootDomainReferenceMock = mCertificateMockSetRootDomainReference{mock: m}

	return m
}

type mCertificateMockGetBootstrapNodes struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetBootstrapNodes to return Return's arguments
func (m *mCertificateMockGetBootstrapNodes) Return(r []core.BootstrapNode) *CertificateMock {
	m.mock.GetBootstrapNodesFunc = func() []core.BootstrapNode {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetBootstrapNodes method
func (m *mCertificateMockGetBootstrapNodes) Set(f func() (r []core.BootstrapNode)) *CertificateMock {
	m.mock.GetBootstrapNodesFunc = f

	return m.mock
}

//GetBootstrapNodes implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetBootstrapNodes() (r []core.BootstrapNode) {
	atomic.AddUint64(&m.GetBootstrapNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetBootstrapNodesCounter, 1)

	if m.GetBootstrapNodesFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetBootstrapNodes")
		return
	}

	return m.GetBootstrapNodesFunc()
}

//GetBootstrapNodesMinimockCounter returns a count of CertificateMock.GetBootstrapNodesFunc invocations
func (m *CertificateMock) GetBootstrapNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetBootstrapNodesCounter)
}

//GetBootstrapNodesMinimockPreCounter returns the value of CertificateMock.GetBootstrapNodes invocations
func (m *CertificateMock) GetBootstrapNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetBootstrapNodesPreCounter)
}

type mCertificateMockGetRole struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetRole to return Return's arguments
func (m *mCertificateMockGetRole) Return(r core.NodeRole) *CertificateMock {
	m.mock.GetRoleFunc = func() core.NodeRole {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetRole method
func (m *mCertificateMockGetRole) Set(f func() (r core.NodeRole)) *CertificateMock {
	m.mock.GetRoleFunc = f

	return m.mock
}

//GetRole implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetRole() (r core.NodeRole) {
	atomic.AddUint64(&m.GetRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetRoleCounter, 1)

	if m.GetRoleFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetRole")
		return
	}

	return m.GetRoleFunc()
}

//GetRoleMinimockCounter returns a count of CertificateMock.GetRoleFunc invocations
func (m *CertificateMock) GetRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRoleCounter)
}

//GetRoleMinimockPreCounter returns the value of CertificateMock.GetRole invocations
func (m *CertificateMock) GetRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRolePreCounter)
}

type mCertificateMockGetRootDomainReference struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetRootDomainReference to return Return's arguments
func (m *mCertificateMockGetRootDomainReference) Return(r *core.RecordRef) *CertificateMock {
	m.mock.GetRootDomainReferenceFunc = func() *core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetRootDomainReference method
func (m *mCertificateMockGetRootDomainReference) Set(f func() (r *core.RecordRef)) *CertificateMock {
	m.mock.GetRootDomainReferenceFunc = f

	return m.mock
}

//GetRootDomainReference implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetRootDomainReference() (r *core.RecordRef) {
	atomic.AddUint64(&m.GetRootDomainReferencePreCounter, 1)
	defer atomic.AddUint64(&m.GetRootDomainReferenceCounter, 1)

	if m.GetRootDomainReferenceFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetRootDomainReference")
		return
	}

	return m.GetRootDomainReferenceFunc()
}

//GetRootDomainReferenceMinimockCounter returns a count of CertificateMock.GetRootDomainReferenceFunc invocations
func (m *CertificateMock) GetRootDomainReferenceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRootDomainReferenceCounter)
}

//GetRootDomainReferenceMinimockPreCounter returns the value of CertificateMock.GetRootDomainReference invocations
func (m *CertificateMock) GetRootDomainReferenceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRootDomainReferencePreCounter)
}

type mCertificateMockSetRootDomainReference struct {
	mock             *CertificateMock
	mockExpectations *CertificateMockSetRootDomainReferenceParams
}

//CertificateMockSetRootDomainReferenceParams represents input parameters of the Certificate.SetRootDomainReference
type CertificateMockSetRootDomainReferenceParams struct {
	p *core.RecordRef
}

//Expect sets up expected params for the Certificate.SetRootDomainReference
func (m *mCertificateMockSetRootDomainReference) Expect(p *core.RecordRef) *mCertificateMockSetRootDomainReference {
	m.mockExpectations = &CertificateMockSetRootDomainReferenceParams{p}
	return m
}

//Return sets up a mock for Certificate.SetRootDomainReference to return Return's arguments
func (m *mCertificateMockSetRootDomainReference) Return() *CertificateMock {
	m.mock.SetRootDomainReferenceFunc = func(p *core.RecordRef) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.SetRootDomainReference method
func (m *mCertificateMockSetRootDomainReference) Set(f func(p *core.RecordRef)) *CertificateMock {
	m.mock.SetRootDomainReferenceFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetRootDomainReference implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) SetRootDomainReference(p *core.RecordRef) {
	atomic.AddUint64(&m.SetRootDomainReferencePreCounter, 1)
	defer atomic.AddUint64(&m.SetRootDomainReferenceCounter, 1)

	if m.SetRootDomainReferenceMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetRootDomainReferenceMock.mockExpectations, CertificateMockSetRootDomainReferenceParams{p},
			"Certificate.SetRootDomainReference got unexpected parameters")

		if m.SetRootDomainReferenceFunc == nil {

			m.t.Fatal("No results are set for the CertificateMock.SetRootDomainReference")

			return
		}
	}

	if m.SetRootDomainReferenceFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.SetRootDomainReference")
		return
	}

	m.SetRootDomainReferenceFunc(p)
}

//SetRootDomainReferenceMinimockCounter returns a count of CertificateMock.SetRootDomainReferenceFunc invocations
func (m *CertificateMock) SetRootDomainReferenceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRootDomainReferenceCounter)
}

//SetRootDomainReferenceMinimockPreCounter returns the value of CertificateMock.SetRootDomainReference invocations
func (m *CertificateMock) SetRootDomainReferenceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRootDomainReferencePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateMock) ValidateCallCounters() {

	if m.GetBootstrapNodesFunc != nil && atomic.LoadUint64(&m.GetBootstrapNodesCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetBootstrapNodes")
	}

	if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRole")
	}

	if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRootDomainReference")
	}

	if m.SetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.SetRootDomainReferenceCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.SetRootDomainReference")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CertificateMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CertificateMock) MinimockFinish() {

	if m.GetBootstrapNodesFunc != nil && atomic.LoadUint64(&m.GetBootstrapNodesCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetBootstrapNodes")
	}

	if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRole")
	}

	if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRootDomainReference")
	}

	if m.SetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.SetRootDomainReferenceCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.SetRootDomainReference")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CertificateMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CertificateMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetBootstrapNodesFunc == nil || atomic.LoadUint64(&m.GetBootstrapNodesCounter) > 0)
		ok = ok && (m.GetRoleFunc == nil || atomic.LoadUint64(&m.GetRoleCounter) > 0)
		ok = ok && (m.GetRootDomainReferenceFunc == nil || atomic.LoadUint64(&m.GetRootDomainReferenceCounter) > 0)
		ok = ok && (m.SetRootDomainReferenceFunc == nil || atomic.LoadUint64(&m.SetRootDomainReferenceCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetBootstrapNodesFunc != nil && atomic.LoadUint64(&m.GetBootstrapNodesCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetBootstrapNodes")
			}

			if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetRole")
			}

			if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetRootDomainReference")
			}

			if m.SetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.SetRootDomainReferenceCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.SetRootDomainReference")
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
func (m *CertificateMock) AllMocksCalled() bool {

	if m.GetBootstrapNodesFunc != nil && atomic.LoadUint64(&m.GetBootstrapNodesCounter) == 0 {
		return false
	}

	if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
		return false
	}

	if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
		return false
	}

	if m.SetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.SetRootDomainReferenceCounter) == 0 {
		return false
	}

	return true
}

package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Certificate" can be found in github.com/insolar/insolar/core
*/
import (
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//CertificateMock implements github.com/insolar/insolar/core.Certificate
type CertificateMock struct {
	t minimock.Tester

	GetDiscoveryNodesFunc       func() (r []core.DiscoveryNode)
	GetDiscoveryNodesCounter    uint64
	GetDiscoveryNodesPreCounter uint64
	GetDiscoveryNodesMock       mCertificateMockGetDiscoveryNodes

	GetNodeRefFunc       func() (r *core.RecordRef)
	GetNodeRefCounter    uint64
	GetNodeRefPreCounter uint64
	GetNodeRefMock       mCertificateMockGetNodeRef

	GetNodeSignFunc       func(p *core.RecordRef) (r []byte, r1 error)
	GetNodeSignCounter    uint64
	GetNodeSignPreCounter uint64
	GetNodeSignMock       mCertificateMockGetNodeSign

	GetPublicKeyFunc       func() (r crypto.PublicKey)
	GetPublicKeyCounter    uint64
	GetPublicKeyPreCounter uint64
	GetPublicKeyMock       mCertificateMockGetPublicKey

	GetRoleFunc       func() (r core.StaticRole)
	GetRoleCounter    uint64
	GetRolePreCounter uint64
	GetRoleMock       mCertificateMockGetRole

	GetRootDomainReferenceFunc       func() (r *core.RecordRef)
	GetRootDomainReferenceCounter    uint64
	GetRootDomainReferencePreCounter uint64
	GetRootDomainReferenceMock       mCertificateMockGetRootDomainReference

	NewCertForHostFunc       func(p string, p1 string, p2 string) (r core.Certificate, r1 error)
	NewCertForHostCounter    uint64
	NewCertForHostPreCounter uint64
	NewCertForHostMock       mCertificateMockNewCertForHost

	VerifyAuthorizationCertificateFunc       func(p core.AuthorizationCertificate) (r bool, r1 error)
	VerifyAuthorizationCertificateCounter    uint64
	VerifyAuthorizationCertificatePreCounter uint64
	VerifyAuthorizationCertificateMock       mCertificateMockVerifyAuthorizationCertificate
}

//NewCertificateMock returns a mock for github.com/insolar/insolar/core.Certificate
func NewCertificateMock(t minimock.Tester) *CertificateMock {
	m := &CertificateMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetDiscoveryNodesMock = mCertificateMockGetDiscoveryNodes{mock: m}
	m.GetNodeRefMock = mCertificateMockGetNodeRef{mock: m}
	m.GetNodeSignMock = mCertificateMockGetNodeSign{mock: m}
	m.GetPublicKeyMock = mCertificateMockGetPublicKey{mock: m}
	m.GetRoleMock = mCertificateMockGetRole{mock: m}
	m.GetRootDomainReferenceMock = mCertificateMockGetRootDomainReference{mock: m}
	m.NewCertForHostMock = mCertificateMockNewCertForHost{mock: m}
	m.VerifyAuthorizationCertificateMock = mCertificateMockVerifyAuthorizationCertificate{mock: m}

	return m
}

type mCertificateMockGetDiscoveryNodes struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetDiscoveryNodes to return Return's arguments
func (m *mCertificateMockGetDiscoveryNodes) Return(r []core.DiscoveryNode) *CertificateMock {
	m.mock.GetDiscoveryNodesFunc = func() []core.DiscoveryNode {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetDiscoveryNodes method
func (m *mCertificateMockGetDiscoveryNodes) Set(f func() (r []core.DiscoveryNode)) *CertificateMock {
	m.mock.GetDiscoveryNodesFunc = f

	return m.mock
}

//GetDiscoveryNodes implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetDiscoveryNodes() (r []core.DiscoveryNode) {
	atomic.AddUint64(&m.GetDiscoveryNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetDiscoveryNodesCounter, 1)

	if m.GetDiscoveryNodesFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetDiscoveryNodes")
		return
	}

	return m.GetDiscoveryNodesFunc()
}

//GetDiscoveryNodesMinimockCounter returns a count of CertificateMock.GetDiscoveryNodesFunc invocations
func (m *CertificateMock) GetDiscoveryNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDiscoveryNodesCounter)
}

//GetDiscoveryNodesMinimockPreCounter returns the value of CertificateMock.GetDiscoveryNodes invocations
func (m *CertificateMock) GetDiscoveryNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDiscoveryNodesPreCounter)
}

type mCertificateMockGetNodeRef struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetNodeRef to return Return's arguments
func (m *mCertificateMockGetNodeRef) Return(r *core.RecordRef) *CertificateMock {
	m.mock.GetNodeRefFunc = func() *core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetNodeRef method
func (m *mCertificateMockGetNodeRef) Set(f func() (r *core.RecordRef)) *CertificateMock {
	m.mock.GetNodeRefFunc = f

	return m.mock
}

//GetNodeRef implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetNodeRef() (r *core.RecordRef) {
	atomic.AddUint64(&m.GetNodeRefPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeRefCounter, 1)

	if m.GetNodeRefFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetNodeRef")
		return
	}

	return m.GetNodeRefFunc()
}

//GetNodeRefMinimockCounter returns a count of CertificateMock.GetNodeRefFunc invocations
func (m *CertificateMock) GetNodeRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefCounter)
}

//GetNodeRefMinimockPreCounter returns the value of CertificateMock.GetNodeRef invocations
func (m *CertificateMock) GetNodeRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefPreCounter)
}

type mCertificateMockGetNodeSign struct {
	mock             *CertificateMock
	mockExpectations *CertificateMockGetNodeSignParams
}

//CertificateMockGetNodeSignParams represents input parameters of the Certificate.GetNodeSign
type CertificateMockGetNodeSignParams struct {
	p *core.RecordRef
}

//Expect sets up expected params for the Certificate.GetNodeSign
func (m *mCertificateMockGetNodeSign) Expect(p *core.RecordRef) *mCertificateMockGetNodeSign {
	m.mockExpectations = &CertificateMockGetNodeSignParams{p}
	return m
}

//Return sets up a mock for Certificate.GetNodeSign to return Return's arguments
func (m *mCertificateMockGetNodeSign) Return(r []byte, r1 error) *CertificateMock {
	m.mock.GetNodeSignFunc = func(p *core.RecordRef) ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetNodeSign method
func (m *mCertificateMockGetNodeSign) Set(f func(p *core.RecordRef) (r []byte, r1 error)) *CertificateMock {
	m.mock.GetNodeSignFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetNodeSign implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetNodeSign(p *core.RecordRef) (r []byte, r1 error) {
	atomic.AddUint64(&m.GetNodeSignPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeSignCounter, 1)

	if m.GetNodeSignMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetNodeSignMock.mockExpectations, CertificateMockGetNodeSignParams{p},
			"Certificate.GetNodeSign got unexpected parameters")

		if m.GetNodeSignFunc == nil {

			m.t.Fatal("No results are set for the CertificateMock.GetNodeSign")

			return
		}
	}

	if m.GetNodeSignFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetNodeSign")
		return
	}

	return m.GetNodeSignFunc(p)
}

//GetNodeSignMinimockCounter returns a count of CertificateMock.GetNodeSignFunc invocations
func (m *CertificateMock) GetNodeSignMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeSignCounter)
}

//GetNodeSignMinimockPreCounter returns the value of CertificateMock.GetNodeSign invocations
func (m *CertificateMock) GetNodeSignMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeSignPreCounter)
}

type mCertificateMockGetPublicKey struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetPublicKey to return Return's arguments
func (m *mCertificateMockGetPublicKey) Return(r crypto.PublicKey) *CertificateMock {
	m.mock.GetPublicKeyFunc = func() crypto.PublicKey {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetPublicKey method
func (m *mCertificateMockGetPublicKey) Set(f func() (r crypto.PublicKey)) *CertificateMock {
	m.mock.GetPublicKeyFunc = f

	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetPublicKey() (r crypto.PublicKey) {
	atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if m.GetPublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetPublicKey")
		return
	}

	return m.GetPublicKeyFunc()
}

//GetPublicKeyMinimockCounter returns a count of CertificateMock.GetPublicKeyFunc invocations
func (m *CertificateMock) GetPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyCounter)
}

//GetPublicKeyMinimockPreCounter returns the value of CertificateMock.GetPublicKey invocations
func (m *CertificateMock) GetPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyPreCounter)
}

type mCertificateMockGetRole struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetRole to return Return's arguments
func (m *mCertificateMockGetRole) Return(r core.StaticRole) *CertificateMock {
	m.mock.GetRoleFunc = func() core.StaticRole {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetRole method
func (m *mCertificateMockGetRole) Set(f func() (r core.StaticRole)) *CertificateMock {
	m.mock.GetRoleFunc = f

	return m.mock
}

//GetRole implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetRole() (r core.StaticRole) {
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

type mCertificateMockNewCertForHost struct {
	mock             *CertificateMock
	mockExpectations *CertificateMockNewCertForHostParams
}

//CertificateMockNewCertForHostParams represents input parameters of the Certificate.NewCertForHost
type CertificateMockNewCertForHostParams struct {
	p  string
	p1 string
	p2 string
}

//Expect sets up expected params for the Certificate.NewCertForHost
func (m *mCertificateMockNewCertForHost) Expect(p string, p1 string, p2 string) *mCertificateMockNewCertForHost {
	m.mockExpectations = &CertificateMockNewCertForHostParams{p, p1, p2}
	return m
}

//Return sets up a mock for Certificate.NewCertForHost to return Return's arguments
func (m *mCertificateMockNewCertForHost) Return(r core.Certificate, r1 error) *CertificateMock {
	m.mock.NewCertForHostFunc = func(p string, p1 string, p2 string) (core.Certificate, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.NewCertForHost method
func (m *mCertificateMockNewCertForHost) Set(f func(p string, p1 string, p2 string) (r core.Certificate, r1 error)) *CertificateMock {
	m.mock.NewCertForHostFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NewCertForHost implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) NewCertForHost(p string, p1 string, p2 string) (r core.Certificate, r1 error) {
	atomic.AddUint64(&m.NewCertForHostPreCounter, 1)
	defer atomic.AddUint64(&m.NewCertForHostCounter, 1)

	if m.NewCertForHostMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NewCertForHostMock.mockExpectations, CertificateMockNewCertForHostParams{p, p1, p2},
			"Certificate.NewCertForHost got unexpected parameters")

		if m.NewCertForHostFunc == nil {

			m.t.Fatal("No results are set for the CertificateMock.NewCertForHost")

			return
		}
	}

	if m.NewCertForHostFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.NewCertForHost")
		return
	}

	return m.NewCertForHostFunc(p, p1, p2)
}

//NewCertForHostMinimockCounter returns a count of CertificateMock.NewCertForHostFunc invocations
func (m *CertificateMock) NewCertForHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewCertForHostCounter)
}

//NewCertForHostMinimockPreCounter returns the value of CertificateMock.NewCertForHost invocations
func (m *CertificateMock) NewCertForHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewCertForHostPreCounter)
}

type mCertificateMockVerifyAuthorizationCertificate struct {
	mock             *CertificateMock
	mockExpectations *CertificateMockVerifyAuthorizationCertificateParams
}

//CertificateMockVerifyAuthorizationCertificateParams represents input parameters of the Certificate.VerifyAuthorizationCertificate
type CertificateMockVerifyAuthorizationCertificateParams struct {
	p core.AuthorizationCertificate
}

//Expect sets up expected params for the Certificate.VerifyAuthorizationCertificate
func (m *mCertificateMockVerifyAuthorizationCertificate) Expect(p core.AuthorizationCertificate) *mCertificateMockVerifyAuthorizationCertificate {
	m.mockExpectations = &CertificateMockVerifyAuthorizationCertificateParams{p}
	return m
}

//Return sets up a mock for Certificate.VerifyAuthorizationCertificate to return Return's arguments
func (m *mCertificateMockVerifyAuthorizationCertificate) Return(r bool, r1 error) *CertificateMock {
	m.mock.VerifyAuthorizationCertificateFunc = func(p core.AuthorizationCertificate) (bool, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.VerifyAuthorizationCertificate method
func (m *mCertificateMockVerifyAuthorizationCertificate) Set(f func(p core.AuthorizationCertificate) (r bool, r1 error)) *CertificateMock {
	m.mock.VerifyAuthorizationCertificateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//VerifyAuthorizationCertificate implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) VerifyAuthorizationCertificate(p core.AuthorizationCertificate) (r bool, r1 error) {
	atomic.AddUint64(&m.VerifyAuthorizationCertificatePreCounter, 1)
	defer atomic.AddUint64(&m.VerifyAuthorizationCertificateCounter, 1)

	if m.VerifyAuthorizationCertificateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.VerifyAuthorizationCertificateMock.mockExpectations, CertificateMockVerifyAuthorizationCertificateParams{p},
			"Certificate.VerifyAuthorizationCertificate got unexpected parameters")

		if m.VerifyAuthorizationCertificateFunc == nil {

			m.t.Fatal("No results are set for the CertificateMock.VerifyAuthorizationCertificate")

			return
		}
	}

	if m.VerifyAuthorizationCertificateFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.VerifyAuthorizationCertificate")
		return
	}

	return m.VerifyAuthorizationCertificateFunc(p)
}

//VerifyAuthorizationCertificateMinimockCounter returns a count of CertificateMock.VerifyAuthorizationCertificateFunc invocations
func (m *CertificateMock) VerifyAuthorizationCertificateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter)
}

//VerifyAuthorizationCertificateMinimockPreCounter returns the value of CertificateMock.VerifyAuthorizationCertificate invocations
func (m *CertificateMock) VerifyAuthorizationCertificateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyAuthorizationCertificatePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateMock) ValidateCallCounters() {

	if m.GetDiscoveryNodesFunc != nil && atomic.LoadUint64(&m.GetDiscoveryNodesCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetDiscoveryNodes")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetNodeRef")
	}

	if m.GetNodeSignFunc != nil && atomic.LoadUint64(&m.GetNodeSignCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetNodeSign")
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetPublicKey")
	}

	if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRole")
	}

	if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRootDomainReference")
	}

	if m.NewCertForHostFunc != nil && atomic.LoadUint64(&m.NewCertForHostCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.NewCertForHost")
	}

	if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.VerifyAuthorizationCertificate")
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

	if m.GetDiscoveryNodesFunc != nil && atomic.LoadUint64(&m.GetDiscoveryNodesCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetDiscoveryNodes")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetNodeRef")
	}

	if m.GetNodeSignFunc != nil && atomic.LoadUint64(&m.GetNodeSignCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetNodeSign")
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetPublicKey")
	}

	if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRole")
	}

	if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetRootDomainReference")
	}

	if m.NewCertForHostFunc != nil && atomic.LoadUint64(&m.NewCertForHostCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.NewCertForHost")
	}

	if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.VerifyAuthorizationCertificate")
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
		ok = ok && (m.GetDiscoveryNodesFunc == nil || atomic.LoadUint64(&m.GetDiscoveryNodesCounter) > 0)
		ok = ok && (m.GetNodeRefFunc == nil || atomic.LoadUint64(&m.GetNodeRefCounter) > 0)
		ok = ok && (m.GetNodeSignFunc == nil || atomic.LoadUint64(&m.GetNodeSignCounter) > 0)
		ok = ok && (m.GetPublicKeyFunc == nil || atomic.LoadUint64(&m.GetPublicKeyCounter) > 0)
		ok = ok && (m.GetRoleFunc == nil || atomic.LoadUint64(&m.GetRoleCounter) > 0)
		ok = ok && (m.GetRootDomainReferenceFunc == nil || atomic.LoadUint64(&m.GetRootDomainReferenceCounter) > 0)
		ok = ok && (m.NewCertForHostFunc == nil || atomic.LoadUint64(&m.NewCertForHostCounter) > 0)
		ok = ok && (m.VerifyAuthorizationCertificateFunc == nil || atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetDiscoveryNodesFunc != nil && atomic.LoadUint64(&m.GetDiscoveryNodesCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetDiscoveryNodes")
			}

			if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetNodeRef")
			}

			if m.GetNodeSignFunc != nil && atomic.LoadUint64(&m.GetNodeSignCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetNodeSign")
			}

			if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetPublicKey")
			}

			if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetRole")
			}

			if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetRootDomainReference")
			}

			if m.NewCertForHostFunc != nil && atomic.LoadUint64(&m.NewCertForHostCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.NewCertForHost")
			}

			if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.VerifyAuthorizationCertificate")
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

	if m.GetDiscoveryNodesFunc != nil && atomic.LoadUint64(&m.GetDiscoveryNodesCounter) == 0 {
		return false
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		return false
	}

	if m.GetNodeSignFunc != nil && atomic.LoadUint64(&m.GetNodeSignCounter) == 0 {
		return false
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		return false
	}

	if m.GetRoleFunc != nil && atomic.LoadUint64(&m.GetRoleCounter) == 0 {
		return false
	}

	if m.GetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == 0 {
		return false
	}

	if m.NewCertForHostFunc != nil && atomic.LoadUint64(&m.NewCertForHostCounter) == 0 {
		return false
	}

	if m.VerifyAuthorizationCertificateFunc != nil && atomic.LoadUint64(&m.VerifyAuthorizationCertificateCounter) == 0 {
		return false
	}

	return true
}

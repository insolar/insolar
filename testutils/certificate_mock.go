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

	GetNodeSignFunc       func(p string) (r []byte, r1 error)
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

	SerializeFunc       func() (r []byte, r1 error)
	SerializeCounter    uint64
	SerializePreCounter uint64
	SerializeMock       mCertificateMockSerialize

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

	m.GetDiscoveryNodesMock = mCertificateMockGetDiscoveryNodes{mock: m}
	m.GetNodeRefMock = mCertificateMockGetNodeRef{mock: m}
	m.GetNodeSignMock = mCertificateMockGetNodeSign{mock: m}
	m.GetPublicKeyMock = mCertificateMockGetPublicKey{mock: m}
	m.GetRoleMock = mCertificateMockGetRole{mock: m}
	m.GetRootDomainReferenceMock = mCertificateMockGetRootDomainReference{mock: m}
	m.SerializeMock = mCertificateMockSerialize{mock: m}
	m.SetRootDomainReferenceMock = mCertificateMockSetRootDomainReference{mock: m}

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
	p string
}

//Expect sets up expected params for the Certificate.GetNodeSign
func (m *mCertificateMockGetNodeSign) Expect(p string) *mCertificateMockGetNodeSign {
	m.mockExpectations = &CertificateMockGetNodeSignParams{p}
	return m
}

//Return sets up a mock for Certificate.GetNodeSign to return Return's arguments
func (m *mCertificateMockGetNodeSign) Return(r []byte, r1 error) *CertificateMock {
	m.mock.GetNodeSignFunc = func(p string) ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetNodeSign method
func (m *mCertificateMockGetNodeSign) Set(f func(p string) (r []byte, r1 error)) *CertificateMock {
	m.mock.GetNodeSignFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetNodeSign implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetNodeSign(p string) (r []byte, r1 error) {
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

type mCertificateMockSerialize struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.Serialize to return Return's arguments
func (m *mCertificateMockSerialize) Return(r []byte, r1 error) *CertificateMock {
	m.mock.SerializeFunc = func() ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.Serialize method
func (m *mCertificateMockSerialize) Set(f func() (r []byte, r1 error)) *CertificateMock {
	m.mock.SerializeFunc = f

	return m.mock
}

//Serialize implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) Serialize() (r []byte, r1 error) {
	atomic.AddUint64(&m.SerializePreCounter, 1)
	defer atomic.AddUint64(&m.SerializeCounter, 1)

	if m.SerializeFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.Serialize")
		return
	}

	return m.SerializeFunc()
}

//SerializeMinimockCounter returns a count of CertificateMock.SerializeFunc invocations
func (m *CertificateMock) SerializeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SerializeCounter)
}

//SerializeMinimockPreCounter returns the value of CertificateMock.Serialize invocations
func (m *CertificateMock) SerializeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SerializePreCounter)
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

	if m.SerializeFunc != nil && atomic.LoadUint64(&m.SerializeCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.Serialize")
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

	if m.SerializeFunc != nil && atomic.LoadUint64(&m.SerializeCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.Serialize")
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
		ok = ok && (m.GetDiscoveryNodesFunc == nil || atomic.LoadUint64(&m.GetDiscoveryNodesCounter) > 0)
		ok = ok && (m.GetNodeRefFunc == nil || atomic.LoadUint64(&m.GetNodeRefCounter) > 0)
		ok = ok && (m.GetNodeSignFunc == nil || atomic.LoadUint64(&m.GetNodeSignCounter) > 0)
		ok = ok && (m.GetPublicKeyFunc == nil || atomic.LoadUint64(&m.GetPublicKeyCounter) > 0)
		ok = ok && (m.GetRoleFunc == nil || atomic.LoadUint64(&m.GetRoleCounter) > 0)
		ok = ok && (m.GetRootDomainReferenceFunc == nil || atomic.LoadUint64(&m.GetRootDomainReferenceCounter) > 0)
		ok = ok && (m.SerializeFunc == nil || atomic.LoadUint64(&m.SerializeCounter) > 0)
		ok = ok && (m.SetRootDomainReferenceFunc == nil || atomic.LoadUint64(&m.SetRootDomainReferenceCounter) > 0)

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

			if m.SerializeFunc != nil && atomic.LoadUint64(&m.SerializeCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.Serialize")
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

	if m.SerializeFunc != nil && atomic.LoadUint64(&m.SerializeCounter) == 0 {
		return false
	}

	if m.SetRootDomainReferenceFunc != nil && atomic.LoadUint64(&m.SetRootDomainReferenceCounter) == 0 {
		return false
	}

	return true
}

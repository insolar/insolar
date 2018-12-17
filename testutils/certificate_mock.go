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
)

//CertificateMock implements github.com/insolar/insolar/core.Certificate
type CertificateMock struct {
	t minimock.Tester

	GetDiscoveryNodesFunc       func() (r []core.DiscoveryNode)
	GetDiscoveryNodesCounter    uint64
	GetDiscoveryNodesPreCounter uint64
	GetDiscoveryNodesMock       mCertificateMockGetDiscoveryNodes

	GetDiscoverySignsFunc       func() (r map[*core.RecordRef][]byte)
	GetDiscoverySignsCounter    uint64
	GetDiscoverySignsPreCounter uint64
	GetDiscoverySignsMock       mCertificateMockGetDiscoverySigns

	GetNodeRefFunc       func() (r *core.RecordRef)
	GetNodeRefCounter    uint64
	GetNodeRefPreCounter uint64
	GetNodeRefMock       mCertificateMockGetNodeRef

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

	SerializeNodePartFunc       func() (r []byte)
	SerializeNodePartCounter    uint64
	SerializeNodePartPreCounter uint64
	SerializeNodePartMock       mCertificateMockSerializeNodePart
}

//NewCertificateMock returns a mock for github.com/insolar/insolar/core.Certificate
func NewCertificateMock(t minimock.Tester) *CertificateMock {
	m := &CertificateMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetDiscoveryNodesMock = mCertificateMockGetDiscoveryNodes{mock: m}
	m.GetDiscoverySignsMock = mCertificateMockGetDiscoverySigns{mock: m}
	m.GetNodeRefMock = mCertificateMockGetNodeRef{mock: m}
	m.GetPublicKeyMock = mCertificateMockGetPublicKey{mock: m}
	m.GetRoleMock = mCertificateMockGetRole{mock: m}
	m.GetRootDomainReferenceMock = mCertificateMockGetRootDomainReference{mock: m}
	m.SerializeNodePartMock = mCertificateMockSerializeNodePart{mock: m}

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

type mCertificateMockGetDiscoverySigns struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.GetDiscoverySigns to return Return's arguments
func (m *mCertificateMockGetDiscoverySigns) Return(r map[*core.RecordRef][]byte) *CertificateMock {
	m.mock.GetDiscoverySignsFunc = func() map[*core.RecordRef][]byte {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.GetDiscoverySigns method
func (m *mCertificateMockGetDiscoverySigns) Set(f func() (r map[*core.RecordRef][]byte)) *CertificateMock {
	m.mock.GetDiscoverySignsFunc = f

	return m.mock
}

//GetDiscoverySigns implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetDiscoverySigns() (r map[*core.RecordRef][]byte) {
	atomic.AddUint64(&m.GetDiscoverySignsPreCounter, 1)
	defer atomic.AddUint64(&m.GetDiscoverySignsCounter, 1)

	if m.GetDiscoverySignsFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.GetDiscoverySigns")
		return
	}

	return m.GetDiscoverySignsFunc()
}

//GetDiscoverySignsMinimockCounter returns a count of CertificateMock.GetDiscoverySignsFunc invocations
func (m *CertificateMock) GetDiscoverySignsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDiscoverySignsCounter)
}

//GetDiscoverySignsMinimockPreCounter returns the value of CertificateMock.GetDiscoverySigns invocations
func (m *CertificateMock) GetDiscoverySignsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDiscoverySignsPreCounter)
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

type mCertificateMockSerializeNodePart struct {
	mock *CertificateMock
}

//Return sets up a mock for Certificate.SerializeNodePart to return Return's arguments
func (m *mCertificateMockSerializeNodePart) Return(r []byte) *CertificateMock {
	m.mock.SerializeNodePartFunc = func() []byte {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Certificate.SerializeNodePart method
func (m *mCertificateMockSerializeNodePart) Set(f func() (r []byte)) *CertificateMock {
	m.mock.SerializeNodePartFunc = f

	return m.mock
}

//SerializeNodePart implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) SerializeNodePart() (r []byte) {
	atomic.AddUint64(&m.SerializeNodePartPreCounter, 1)
	defer atomic.AddUint64(&m.SerializeNodePartCounter, 1)

	if m.SerializeNodePartFunc == nil {
		m.t.Fatal("Unexpected call to CertificateMock.SerializeNodePart")
		return
	}

	return m.SerializeNodePartFunc()
}

//SerializeNodePartMinimockCounter returns a count of CertificateMock.SerializeNodePartFunc invocations
func (m *CertificateMock) SerializeNodePartMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SerializeNodePartCounter)
}

//SerializeNodePartMinimockPreCounter returns the value of CertificateMock.SerializeNodePart invocations
func (m *CertificateMock) SerializeNodePartMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SerializeNodePartPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateMock) ValidateCallCounters() {

	if m.GetDiscoveryNodesFunc != nil && atomic.LoadUint64(&m.GetDiscoveryNodesCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetDiscoveryNodes")
	}

	if m.GetDiscoverySignsFunc != nil && atomic.LoadUint64(&m.GetDiscoverySignsCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetDiscoverySigns")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetNodeRef")
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

	if m.SerializeNodePartFunc != nil && atomic.LoadUint64(&m.SerializeNodePartCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.SerializeNodePart")
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

	if m.GetDiscoverySignsFunc != nil && atomic.LoadUint64(&m.GetDiscoverySignsCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetDiscoverySigns")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.GetNodeRef")
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

	if m.SerializeNodePartFunc != nil && atomic.LoadUint64(&m.SerializeNodePartCounter) == 0 {
		m.t.Fatal("Expected call to CertificateMock.SerializeNodePart")
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
		ok = ok && (m.GetDiscoverySignsFunc == nil || atomic.LoadUint64(&m.GetDiscoverySignsCounter) > 0)
		ok = ok && (m.GetNodeRefFunc == nil || atomic.LoadUint64(&m.GetNodeRefCounter) > 0)
		ok = ok && (m.GetPublicKeyFunc == nil || atomic.LoadUint64(&m.GetPublicKeyCounter) > 0)
		ok = ok && (m.GetRoleFunc == nil || atomic.LoadUint64(&m.GetRoleCounter) > 0)
		ok = ok && (m.GetRootDomainReferenceFunc == nil || atomic.LoadUint64(&m.GetRootDomainReferenceCounter) > 0)
		ok = ok && (m.SerializeNodePartFunc == nil || atomic.LoadUint64(&m.SerializeNodePartCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetDiscoveryNodesFunc != nil && atomic.LoadUint64(&m.GetDiscoveryNodesCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetDiscoveryNodes")
			}

			if m.GetDiscoverySignsFunc != nil && atomic.LoadUint64(&m.GetDiscoverySignsCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetDiscoverySigns")
			}

			if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.GetNodeRef")
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

			if m.SerializeNodePartFunc != nil && atomic.LoadUint64(&m.SerializeNodePartCounter) == 0 {
				m.t.Error("Expected call to CertificateMock.SerializeNodePart")
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

	if m.GetDiscoverySignsFunc != nil && atomic.LoadUint64(&m.GetDiscoverySignsCounter) == 0 {
		return false
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
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

	if m.SerializeNodePartFunc != nil && atomic.LoadUint64(&m.SerializeNodePartCounter) == 0 {
		return false
	}

	return true
}

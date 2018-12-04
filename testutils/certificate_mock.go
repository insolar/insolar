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

	return m
}

type mCertificateMockGetDiscoveryNodes struct {
	mock              *CertificateMock
	mainExpectation   *CertificateMockGetDiscoveryNodesExpectation
	expectationSeries []*CertificateMockGetDiscoveryNodesExpectation
}

type CertificateMockGetDiscoveryNodesExpectation struct {
	result *CertificateMockGetDiscoveryNodesResult
}

type CertificateMockGetDiscoveryNodesResult struct {
	r []core.DiscoveryNode
}

//Expect specifies that invocation of Certificate.GetDiscoveryNodes is expected from 1 to Infinity times
func (m *mCertificateMockGetDiscoveryNodes) Expect() *mCertificateMockGetDiscoveryNodes {
	m.mock.GetDiscoveryNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetDiscoveryNodesExpectation{}
	}

	return m
}

//Return specifies results of invocation of Certificate.GetDiscoveryNodes
func (m *mCertificateMockGetDiscoveryNodes) Return(r []core.DiscoveryNode) *CertificateMock {
	m.mock.GetDiscoveryNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetDiscoveryNodesExpectation{}
	}
	m.mainExpectation.result = &CertificateMockGetDiscoveryNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Certificate.GetDiscoveryNodes is expected once
func (m *mCertificateMockGetDiscoveryNodes) ExpectOnce() *CertificateMockGetDiscoveryNodesExpectation {
	m.mock.GetDiscoveryNodesFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateMockGetDiscoveryNodesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateMockGetDiscoveryNodesExpectation) Return(r []core.DiscoveryNode) {
	e.result = &CertificateMockGetDiscoveryNodesResult{r}
}

//Set uses given function f as a mock of Certificate.GetDiscoveryNodes method
func (m *mCertificateMockGetDiscoveryNodes) Set(f func() (r []core.DiscoveryNode)) *CertificateMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDiscoveryNodesFunc = f
	return m.mock
}

//GetDiscoveryNodes implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetDiscoveryNodes() (r []core.DiscoveryNode) {
	counter := atomic.AddUint64(&m.GetDiscoveryNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetDiscoveryNodesCounter, 1)

	if len(m.GetDiscoveryNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDiscoveryNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateMock.GetDiscoveryNodes.")
			return
		}

		result := m.GetDiscoveryNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetDiscoveryNodes")
			return
		}

		r = result.r

		return
	}

	if m.GetDiscoveryNodesMock.mainExpectation != nil {

		result := m.GetDiscoveryNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetDiscoveryNodes")
		}

		r = result.r

		return
	}

	if m.GetDiscoveryNodesFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateMock.GetDiscoveryNodes.")
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

//GetDiscoveryNodesFinished returns true if mock invocations count is ok
func (m *CertificateMock) GetDiscoveryNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDiscoveryNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDiscoveryNodesCounter) == uint64(len(m.GetDiscoveryNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDiscoveryNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDiscoveryNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDiscoveryNodesFunc != nil {
		return atomic.LoadUint64(&m.GetDiscoveryNodesCounter) > 0
	}

	return true
}

type mCertificateMockGetNodeRef struct {
	mock              *CertificateMock
	mainExpectation   *CertificateMockGetNodeRefExpectation
	expectationSeries []*CertificateMockGetNodeRefExpectation
}

type CertificateMockGetNodeRefExpectation struct {
	result *CertificateMockGetNodeRefResult
}

type CertificateMockGetNodeRefResult struct {
	r *core.RecordRef
}

//Expect specifies that invocation of Certificate.GetNodeRef is expected from 1 to Infinity times
func (m *mCertificateMockGetNodeRef) Expect() *mCertificateMockGetNodeRef {
	m.mock.GetNodeRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetNodeRefExpectation{}
	}

	return m
}

//Return specifies results of invocation of Certificate.GetNodeRef
func (m *mCertificateMockGetNodeRef) Return(r *core.RecordRef) *CertificateMock {
	m.mock.GetNodeRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetNodeRefExpectation{}
	}
	m.mainExpectation.result = &CertificateMockGetNodeRefResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Certificate.GetNodeRef is expected once
func (m *mCertificateMockGetNodeRef) ExpectOnce() *CertificateMockGetNodeRefExpectation {
	m.mock.GetNodeRefFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateMockGetNodeRefExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateMockGetNodeRefExpectation) Return(r *core.RecordRef) {
	e.result = &CertificateMockGetNodeRefResult{r}
}

//Set uses given function f as a mock of Certificate.GetNodeRef method
func (m *mCertificateMockGetNodeRef) Set(f func() (r *core.RecordRef)) *CertificateMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeRefFunc = f
	return m.mock
}

//GetNodeRef implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetNodeRef() (r *core.RecordRef) {
	counter := atomic.AddUint64(&m.GetNodeRefPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeRefCounter, 1)

	if len(m.GetNodeRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateMock.GetNodeRef.")
			return
		}

		result := m.GetNodeRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetNodeRef")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeRefMock.mainExpectation != nil {

		result := m.GetNodeRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetNodeRef")
		}

		r = result.r

		return
	}

	if m.GetNodeRefFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateMock.GetNodeRef.")
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

//GetNodeRefFinished returns true if mock invocations count is ok
func (m *CertificateMock) GetNodeRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeRefCounter) == uint64(len(m.GetNodeRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeRefFunc != nil {
		return atomic.LoadUint64(&m.GetNodeRefCounter) > 0
	}

	return true
}

type mCertificateMockGetNodeSign struct {
	mock              *CertificateMock
	mainExpectation   *CertificateMockGetNodeSignExpectation
	expectationSeries []*CertificateMockGetNodeSignExpectation
}

type CertificateMockGetNodeSignExpectation struct {
	input  *CertificateMockGetNodeSignInput
	result *CertificateMockGetNodeSignResult
}

type CertificateMockGetNodeSignInput struct {
	p *core.RecordRef
}

type CertificateMockGetNodeSignResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of Certificate.GetNodeSign is expected from 1 to Infinity times
func (m *mCertificateMockGetNodeSign) Expect(p *core.RecordRef) *mCertificateMockGetNodeSign {
	m.mock.GetNodeSignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetNodeSignExpectation{}
	}
	m.mainExpectation.input = &CertificateMockGetNodeSignInput{p}
	return m
}

//Return specifies results of invocation of Certificate.GetNodeSign
func (m *mCertificateMockGetNodeSign) Return(r []byte, r1 error) *CertificateMock {
	m.mock.GetNodeSignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetNodeSignExpectation{}
	}
	m.mainExpectation.result = &CertificateMockGetNodeSignResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Certificate.GetNodeSign is expected once
func (m *mCertificateMockGetNodeSign) ExpectOnce(p *core.RecordRef) *CertificateMockGetNodeSignExpectation {
	m.mock.GetNodeSignFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateMockGetNodeSignExpectation{}
	expectation.input = &CertificateMockGetNodeSignInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateMockGetNodeSignExpectation) Return(r []byte, r1 error) {
	e.result = &CertificateMockGetNodeSignResult{r, r1}
}

//Set uses given function f as a mock of Certificate.GetNodeSign method
func (m *mCertificateMockGetNodeSign) Set(f func(p *core.RecordRef) (r []byte, r1 error)) *CertificateMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeSignFunc = f
	return m.mock
}

//GetNodeSign implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetNodeSign(p *core.RecordRef) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.GetNodeSignPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeSignCounter, 1)

	if len(m.GetNodeSignMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeSignMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateMock.GetNodeSign. %v", p)
			return
		}

		input := m.GetNodeSignMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CertificateMockGetNodeSignInput{p}, "Certificate.GetNodeSign got unexpected parameters")

		result := m.GetNodeSignMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetNodeSign")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetNodeSignMock.mainExpectation != nil {

		input := m.GetNodeSignMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CertificateMockGetNodeSignInput{p}, "Certificate.GetNodeSign got unexpected parameters")
		}

		result := m.GetNodeSignMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetNodeSign")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetNodeSignFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateMock.GetNodeSign. %v", p)
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

//GetNodeSignFinished returns true if mock invocations count is ok
func (m *CertificateMock) GetNodeSignFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeSignMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeSignCounter) == uint64(len(m.GetNodeSignMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeSignMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeSignCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeSignFunc != nil {
		return atomic.LoadUint64(&m.GetNodeSignCounter) > 0
	}

	return true
}

type mCertificateMockGetPublicKey struct {
	mock              *CertificateMock
	mainExpectation   *CertificateMockGetPublicKeyExpectation
	expectationSeries []*CertificateMockGetPublicKeyExpectation
}

type CertificateMockGetPublicKeyExpectation struct {
	result *CertificateMockGetPublicKeyResult
}

type CertificateMockGetPublicKeyResult struct {
	r crypto.PublicKey
}

//Expect specifies that invocation of Certificate.GetPublicKey is expected from 1 to Infinity times
func (m *mCertificateMockGetPublicKey) Expect() *mCertificateMockGetPublicKey {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetPublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of Certificate.GetPublicKey
func (m *mCertificateMockGetPublicKey) Return(r crypto.PublicKey) *CertificateMock {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetPublicKeyExpectation{}
	}
	m.mainExpectation.result = &CertificateMockGetPublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Certificate.GetPublicKey is expected once
func (m *mCertificateMockGetPublicKey) ExpectOnce() *CertificateMockGetPublicKeyExpectation {
	m.mock.GetPublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateMockGetPublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateMockGetPublicKeyExpectation) Return(r crypto.PublicKey) {
	e.result = &CertificateMockGetPublicKeyResult{r}
}

//Set uses given function f as a mock of Certificate.GetPublicKey method
func (m *mCertificateMockGetPublicKey) Set(f func() (r crypto.PublicKey)) *CertificateMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyFunc = f
	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetPublicKey() (r crypto.PublicKey) {
	counter := atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if len(m.GetPublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateMock.GetPublicKey.")
			return
		}

		result := m.GetPublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetPublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyMock.mainExpectation != nil {

		result := m.GetPublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetPublicKey")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateMock.GetPublicKey.")
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

//GetPublicKeyFinished returns true if mock invocations count is ok
func (m *CertificateMock) GetPublicKeyFinished() bool {
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

type mCertificateMockGetRole struct {
	mock              *CertificateMock
	mainExpectation   *CertificateMockGetRoleExpectation
	expectationSeries []*CertificateMockGetRoleExpectation
}

type CertificateMockGetRoleExpectation struct {
	result *CertificateMockGetRoleResult
}

type CertificateMockGetRoleResult struct {
	r core.StaticRole
}

//Expect specifies that invocation of Certificate.GetRole is expected from 1 to Infinity times
func (m *mCertificateMockGetRole) Expect() *mCertificateMockGetRole {
	m.mock.GetRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of Certificate.GetRole
func (m *mCertificateMockGetRole) Return(r core.StaticRole) *CertificateMock {
	m.mock.GetRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetRoleExpectation{}
	}
	m.mainExpectation.result = &CertificateMockGetRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Certificate.GetRole is expected once
func (m *mCertificateMockGetRole) ExpectOnce() *CertificateMockGetRoleExpectation {
	m.mock.GetRoleFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateMockGetRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateMockGetRoleExpectation) Return(r core.StaticRole) {
	e.result = &CertificateMockGetRoleResult{r}
}

//Set uses given function f as a mock of Certificate.GetRole method
func (m *mCertificateMockGetRole) Set(f func() (r core.StaticRole)) *CertificateMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRoleFunc = f
	return m.mock
}

//GetRole implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetRole() (r core.StaticRole) {
	counter := atomic.AddUint64(&m.GetRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetRoleCounter, 1)

	if len(m.GetRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateMock.GetRole.")
			return
		}

		result := m.GetRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetRole")
			return
		}

		r = result.r

		return
	}

	if m.GetRoleMock.mainExpectation != nil {

		result := m.GetRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetRole")
		}

		r = result.r

		return
	}

	if m.GetRoleFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateMock.GetRole.")
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

//GetRoleFinished returns true if mock invocations count is ok
func (m *CertificateMock) GetRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRoleCounter) == uint64(len(m.GetRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRoleFunc != nil {
		return atomic.LoadUint64(&m.GetRoleCounter) > 0
	}

	return true
}

type mCertificateMockGetRootDomainReference struct {
	mock              *CertificateMock
	mainExpectation   *CertificateMockGetRootDomainReferenceExpectation
	expectationSeries []*CertificateMockGetRootDomainReferenceExpectation
}

type CertificateMockGetRootDomainReferenceExpectation struct {
	result *CertificateMockGetRootDomainReferenceResult
}

type CertificateMockGetRootDomainReferenceResult struct {
	r *core.RecordRef
}

//Expect specifies that invocation of Certificate.GetRootDomainReference is expected from 1 to Infinity times
func (m *mCertificateMockGetRootDomainReference) Expect() *mCertificateMockGetRootDomainReference {
	m.mock.GetRootDomainReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetRootDomainReferenceExpectation{}
	}

	return m
}

//Return specifies results of invocation of Certificate.GetRootDomainReference
func (m *mCertificateMockGetRootDomainReference) Return(r *core.RecordRef) *CertificateMock {
	m.mock.GetRootDomainReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockGetRootDomainReferenceExpectation{}
	}
	m.mainExpectation.result = &CertificateMockGetRootDomainReferenceResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Certificate.GetRootDomainReference is expected once
func (m *mCertificateMockGetRootDomainReference) ExpectOnce() *CertificateMockGetRootDomainReferenceExpectation {
	m.mock.GetRootDomainReferenceFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateMockGetRootDomainReferenceExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateMockGetRootDomainReferenceExpectation) Return(r *core.RecordRef) {
	e.result = &CertificateMockGetRootDomainReferenceResult{r}
}

//Set uses given function f as a mock of Certificate.GetRootDomainReference method
func (m *mCertificateMockGetRootDomainReference) Set(f func() (r *core.RecordRef)) *CertificateMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRootDomainReferenceFunc = f
	return m.mock
}

//GetRootDomainReference implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) GetRootDomainReference() (r *core.RecordRef) {
	counter := atomic.AddUint64(&m.GetRootDomainReferencePreCounter, 1)
	defer atomic.AddUint64(&m.GetRootDomainReferenceCounter, 1)

	if len(m.GetRootDomainReferenceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRootDomainReferenceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateMock.GetRootDomainReference.")
			return
		}

		result := m.GetRootDomainReferenceMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetRootDomainReference")
			return
		}

		r = result.r

		return
	}

	if m.GetRootDomainReferenceMock.mainExpectation != nil {

		result := m.GetRootDomainReferenceMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.GetRootDomainReference")
		}

		r = result.r

		return
	}

	if m.GetRootDomainReferenceFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateMock.GetRootDomainReference.")
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

//GetRootDomainReferenceFinished returns true if mock invocations count is ok
func (m *CertificateMock) GetRootDomainReferenceFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRootDomainReferenceMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRootDomainReferenceCounter) == uint64(len(m.GetRootDomainReferenceMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRootDomainReferenceMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRootDomainReferenceCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRootDomainReferenceFunc != nil {
		return atomic.LoadUint64(&m.GetRootDomainReferenceCounter) > 0
	}

	return true
}

type mCertificateMockNewCertForHost struct {
	mock              *CertificateMock
	mainExpectation   *CertificateMockNewCertForHostExpectation
	expectationSeries []*CertificateMockNewCertForHostExpectation
}

type CertificateMockNewCertForHostExpectation struct {
	input  *CertificateMockNewCertForHostInput
	result *CertificateMockNewCertForHostResult
}

type CertificateMockNewCertForHostInput struct {
	p  string
	p1 string
	p2 string
}

type CertificateMockNewCertForHostResult struct {
	r  core.Certificate
	r1 error
}

//Expect specifies that invocation of Certificate.NewCertForHost is expected from 1 to Infinity times
func (m *mCertificateMockNewCertForHost) Expect(p string, p1 string, p2 string) *mCertificateMockNewCertForHost {
	m.mock.NewCertForHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockNewCertForHostExpectation{}
	}
	m.mainExpectation.input = &CertificateMockNewCertForHostInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Certificate.NewCertForHost
func (m *mCertificateMockNewCertForHost) Return(r core.Certificate, r1 error) *CertificateMock {
	m.mock.NewCertForHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CertificateMockNewCertForHostExpectation{}
	}
	m.mainExpectation.result = &CertificateMockNewCertForHostResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Certificate.NewCertForHost is expected once
func (m *mCertificateMockNewCertForHost) ExpectOnce(p string, p1 string, p2 string) *CertificateMockNewCertForHostExpectation {
	m.mock.NewCertForHostFunc = nil
	m.mainExpectation = nil

	expectation := &CertificateMockNewCertForHostExpectation{}
	expectation.input = &CertificateMockNewCertForHostInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CertificateMockNewCertForHostExpectation) Return(r core.Certificate, r1 error) {
	e.result = &CertificateMockNewCertForHostResult{r, r1}
}

//Set uses given function f as a mock of Certificate.NewCertForHost method
func (m *mCertificateMockNewCertForHost) Set(f func(p string, p1 string, p2 string) (r core.Certificate, r1 error)) *CertificateMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NewCertForHostFunc = f
	return m.mock
}

//NewCertForHost implements github.com/insolar/insolar/core.Certificate interface
func (m *CertificateMock) NewCertForHost(p string, p1 string, p2 string) (r core.Certificate, r1 error) {
	counter := atomic.AddUint64(&m.NewCertForHostPreCounter, 1)
	defer atomic.AddUint64(&m.NewCertForHostCounter, 1)

	if len(m.NewCertForHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NewCertForHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CertificateMock.NewCertForHost. %v %v %v", p, p1, p2)
			return
		}

		input := m.NewCertForHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CertificateMockNewCertForHostInput{p, p1, p2}, "Certificate.NewCertForHost got unexpected parameters")

		result := m.NewCertForHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.NewCertForHost")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewCertForHostMock.mainExpectation != nil {

		input := m.NewCertForHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CertificateMockNewCertForHostInput{p, p1, p2}, "Certificate.NewCertForHost got unexpected parameters")
		}

		result := m.NewCertForHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CertificateMock.NewCertForHost")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewCertForHostFunc == nil {
		m.t.Fatalf("Unexpected call to CertificateMock.NewCertForHost. %v %v %v", p, p1, p2)
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

//NewCertForHostFinished returns true if mock invocations count is ok
func (m *CertificateMock) NewCertForHostFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NewCertForHostMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NewCertForHostCounter) == uint64(len(m.NewCertForHostMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NewCertForHostMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NewCertForHostCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NewCertForHostFunc != nil {
		return atomic.LoadUint64(&m.NewCertForHostCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CertificateMock) ValidateCallCounters() {

	if !m.GetDiscoveryNodesFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetDiscoveryNodes")
	}

	if !m.GetNodeRefFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetNodeRef")
	}

	if !m.GetNodeSignFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetNodeSign")
	}

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetPublicKey")
	}

	if !m.GetRoleFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetRole")
	}

	if !m.GetRootDomainReferenceFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetRootDomainReference")
	}

	if !m.NewCertForHostFinished() {
		m.t.Fatal("Expected call to CertificateMock.NewCertForHost")
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

	if !m.GetDiscoveryNodesFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetDiscoveryNodes")
	}

	if !m.GetNodeRefFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetNodeRef")
	}

	if !m.GetNodeSignFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetNodeSign")
	}

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetPublicKey")
	}

	if !m.GetRoleFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetRole")
	}

	if !m.GetRootDomainReferenceFinished() {
		m.t.Fatal("Expected call to CertificateMock.GetRootDomainReference")
	}

	if !m.NewCertForHostFinished() {
		m.t.Fatal("Expected call to CertificateMock.NewCertForHost")
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
		ok = ok && m.GetDiscoveryNodesFinished()
		ok = ok && m.GetNodeRefFinished()
		ok = ok && m.GetNodeSignFinished()
		ok = ok && m.GetPublicKeyFinished()
		ok = ok && m.GetRoleFinished()
		ok = ok && m.GetRootDomainReferenceFinished()
		ok = ok && m.NewCertForHostFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetDiscoveryNodesFinished() {
				m.t.Error("Expected call to CertificateMock.GetDiscoveryNodes")
			}

			if !m.GetNodeRefFinished() {
				m.t.Error("Expected call to CertificateMock.GetNodeRef")
			}

			if !m.GetNodeSignFinished() {
				m.t.Error("Expected call to CertificateMock.GetNodeSign")
			}

			if !m.GetPublicKeyFinished() {
				m.t.Error("Expected call to CertificateMock.GetPublicKey")
			}

			if !m.GetRoleFinished() {
				m.t.Error("Expected call to CertificateMock.GetRole")
			}

			if !m.GetRootDomainReferenceFinished() {
				m.t.Error("Expected call to CertificateMock.GetRootDomainReference")
			}

			if !m.NewCertForHostFinished() {
				m.t.Error("Expected call to CertificateMock.NewCertForHost")
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

	if !m.GetDiscoveryNodesFinished() {
		return false
	}

	if !m.GetNodeRefFinished() {
		return false
	}

	if !m.GetNodeSignFinished() {
		return false
	}

	if !m.GetPublicKeyFinished() {
		return false
	}

	if !m.GetRoleFinished() {
		return false
	}

	if !m.GetRootDomainReferenceFinished() {
		return false
	}

	if !m.NewCertForHostFinished() {
		return false
	}

	return true
}

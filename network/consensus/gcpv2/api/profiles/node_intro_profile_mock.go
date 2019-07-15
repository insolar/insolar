package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeIntroProfile" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	endpoints "github.com/insolar/insolar/network/consensus/common/endpoints"
	member "github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeIntroProfileMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile
type NodeIntroProfileMock struct {
	t minimock.Tester

	GetAnnouncementSignatureFunc       func() (r cryptkit.SignatureHolder)
	GetAnnouncementSignatureCounter    uint64
	GetAnnouncementSignaturePreCounter uint64
	GetAnnouncementSignatureMock       mNodeIntroProfileMockGetAnnouncementSignature

	GetDefaultEndpointFunc       func() (r endpoints.Outbound)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mNodeIntroProfileMockGetDefaultEndpoint

	GetIntroductionFunc       func() (r NodeIntroduction)
	GetIntroductionCounter    uint64
	GetIntroductionPreCounter uint64
	GetIntroductionMock       mNodeIntroProfileMockGetIntroduction

	GetNodePublicKeyFunc       func() (r cryptkit.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mNodeIntroProfileMockGetNodePublicKey

	GetNodePublicKeyStoreFunc       func() (r cryptkit.PublicKeyStore)
	GetNodePublicKeyStoreCounter    uint64
	GetNodePublicKeyStorePreCounter uint64
	GetNodePublicKeyStoreMock       mNodeIntroProfileMockGetNodePublicKeyStore

	GetPrimaryRoleFunc       func() (r member.PrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mNodeIntroProfileMockGetPrimaryRole

	GetShortNodeIDFunc       func() (r insolar.ShortNodeID)
	GetShortNodeIDCounter    uint64
	GetShortNodeIDPreCounter uint64
	GetShortNodeIDMock       mNodeIntroProfileMockGetShortNodeID

	GetSpecialRolesFunc       func() (r member.SpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mNodeIntroProfileMockGetSpecialRoles

	GetStartPowerFunc       func() (r member.Power)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mNodeIntroProfileMockGetStartPower

	HasIntroductionFunc       func() (r bool)
	HasIntroductionCounter    uint64
	HasIntroductionPreCounter uint64
	HasIntroductionMock       mNodeIntroProfileMockHasIntroduction

	IsAcceptableHostFunc       func(p endpoints.Inbound) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mNodeIntroProfileMockIsAcceptableHost
}

//NewNodeIntroProfileMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile
func NewNodeIntroProfileMock(t minimock.Tester) *NodeIntroProfileMock {
	m := &NodeIntroProfileMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAnnouncementSignatureMock = mNodeIntroProfileMockGetAnnouncementSignature{mock: m}
	m.GetDefaultEndpointMock = mNodeIntroProfileMockGetDefaultEndpoint{mock: m}
	m.GetIntroductionMock = mNodeIntroProfileMockGetIntroduction{mock: m}
	m.GetNodePublicKeyMock = mNodeIntroProfileMockGetNodePublicKey{mock: m}
	m.GetNodePublicKeyStoreMock = mNodeIntroProfileMockGetNodePublicKeyStore{mock: m}
	m.GetPrimaryRoleMock = mNodeIntroProfileMockGetPrimaryRole{mock: m}
	m.GetShortNodeIDMock = mNodeIntroProfileMockGetShortNodeID{mock: m}
	m.GetSpecialRolesMock = mNodeIntroProfileMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mNodeIntroProfileMockGetStartPower{mock: m}
	m.HasIntroductionMock = mNodeIntroProfileMockHasIntroduction{mock: m}
	m.IsAcceptableHostMock = mNodeIntroProfileMockIsAcceptableHost{mock: m}

	return m
}

type mNodeIntroProfileMockGetAnnouncementSignature struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetAnnouncementSignatureExpectation
	expectationSeries []*NodeIntroProfileMockGetAnnouncementSignatureExpectation
}

type NodeIntroProfileMockGetAnnouncementSignatureExpectation struct {
	result *NodeIntroProfileMockGetAnnouncementSignatureResult
}

type NodeIntroProfileMockGetAnnouncementSignatureResult struct {
	r cryptkit.SignatureHolder
}

//Expect specifies that invocation of NodeIntroProfile.GetAnnouncementSignature is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetAnnouncementSignature) Expect() *mNodeIntroProfileMockGetAnnouncementSignature {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetAnnouncementSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetAnnouncementSignature
func (m *mNodeIntroProfileMockGetAnnouncementSignature) Return(r cryptkit.SignatureHolder) *NodeIntroProfileMock {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetAnnouncementSignatureExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetAnnouncementSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetAnnouncementSignature is expected once
func (m *mNodeIntroProfileMockGetAnnouncementSignature) ExpectOnce() *NodeIntroProfileMockGetAnnouncementSignatureExpectation {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetAnnouncementSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetAnnouncementSignatureExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &NodeIntroProfileMockGetAnnouncementSignatureResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetAnnouncementSignature method
func (m *mNodeIntroProfileMockGetAnnouncementSignature) Set(f func() (r cryptkit.SignatureHolder)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAnnouncementSignatureFunc = f
	return m.mock
}

//GetAnnouncementSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetAnnouncementSignature() (r cryptkit.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetAnnouncementSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetAnnouncementSignatureCounter, 1)

	if len(m.GetAnnouncementSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAnnouncementSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetAnnouncementSignature.")
			return
		}

		result := m.GetAnnouncementSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetAnnouncementSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureMock.mainExpectation != nil {

		result := m.GetAnnouncementSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetAnnouncementSignature")
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetAnnouncementSignature.")
		return
	}

	return m.GetAnnouncementSignatureFunc()
}

//GetAnnouncementSignatureMinimockCounter returns a count of NodeIntroProfileMock.GetAnnouncementSignatureFunc invocations
func (m *NodeIntroProfileMock) GetAnnouncementSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter)
}

//GetAnnouncementSignatureMinimockPreCounter returns the value of NodeIntroProfileMock.GetAnnouncementSignature invocations
func (m *NodeIntroProfileMock) GetAnnouncementSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignaturePreCounter)
}

//GetAnnouncementSignatureFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetAnnouncementSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetAnnouncementSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter) == uint64(len(m.GetAnnouncementSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetAnnouncementSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetAnnouncementSignatureFunc != nil {
		return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetDefaultEndpoint struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetDefaultEndpointExpectation
	expectationSeries []*NodeIntroProfileMockGetDefaultEndpointExpectation
}

type NodeIntroProfileMockGetDefaultEndpointExpectation struct {
	result *NodeIntroProfileMockGetDefaultEndpointResult
}

type NodeIntroProfileMockGetDefaultEndpointResult struct {
	r endpoints.Outbound
}

//Expect specifies that invocation of NodeIntroProfile.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetDefaultEndpoint) Expect() *mNodeIntroProfileMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetDefaultEndpoint
func (m *mNodeIntroProfileMockGetDefaultEndpoint) Return(r endpoints.Outbound) *NodeIntroProfileMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetDefaultEndpoint is expected once
func (m *mNodeIntroProfileMockGetDefaultEndpoint) ExpectOnce() *NodeIntroProfileMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetDefaultEndpointExpectation) Return(r endpoints.Outbound) {
	e.result = &NodeIntroProfileMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetDefaultEndpoint method
func (m *mNodeIntroProfileMockGetDefaultEndpoint) Set(f func() (r endpoints.Outbound)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetDefaultEndpoint() (r endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of NodeIntroProfileMock.GetDefaultEndpointFunc invocations
func (m *NodeIntroProfileMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of NodeIntroProfileMock.GetDefaultEndpoint invocations
func (m *NodeIntroProfileMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetDefaultEndpointFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDefaultEndpointCounter) == uint64(len(m.GetDefaultEndpointMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDefaultEndpointMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDefaultEndpointCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDefaultEndpointFunc != nil {
		return atomic.LoadUint64(&m.GetDefaultEndpointCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetIntroduction struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetIntroductionExpectation
	expectationSeries []*NodeIntroProfileMockGetIntroductionExpectation
}

type NodeIntroProfileMockGetIntroductionExpectation struct {
	result *NodeIntroProfileMockGetIntroductionResult
}

type NodeIntroProfileMockGetIntroductionResult struct {
	r NodeIntroduction
}

//Expect specifies that invocation of NodeIntroProfile.GetIntroduction is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetIntroduction) Expect() *mNodeIntroProfileMockGetIntroduction {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetIntroduction
func (m *mNodeIntroProfileMockGetIntroduction) Return(r NodeIntroduction) *NodeIntroProfileMock {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetIntroductionExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetIntroduction is expected once
func (m *mNodeIntroProfileMockGetIntroduction) ExpectOnce() *NodeIntroProfileMockGetIntroductionExpectation {
	m.mock.GetIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetIntroductionExpectation) Return(r NodeIntroduction) {
	e.result = &NodeIntroProfileMockGetIntroductionResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetIntroduction method
func (m *mNodeIntroProfileMockGetIntroduction) Set(f func() (r NodeIntroduction)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIntroductionFunc = f
	return m.mock
}

//GetIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetIntroduction() (r NodeIntroduction) {
	counter := atomic.AddUint64(&m.GetIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.GetIntroductionCounter, 1)

	if len(m.GetIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetIntroduction.")
			return
		}

		result := m.GetIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.GetIntroductionMock.mainExpectation != nil {

		result := m.GetIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetIntroduction")
		}

		r = result.r

		return
	}

	if m.GetIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetIntroduction.")
		return
	}

	return m.GetIntroductionFunc()
}

//GetIntroductionMinimockCounter returns a count of NodeIntroProfileMock.GetIntroductionFunc invocations
func (m *NodeIntroProfileMock) GetIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionCounter)
}

//GetIntroductionMinimockPreCounter returns the value of NodeIntroProfileMock.GetIntroduction invocations
func (m *NodeIntroProfileMock) GetIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionPreCounter)
}

//GetIntroductionFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetIntroductionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIntroductionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIntroductionCounter) == uint64(len(m.GetIntroductionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIntroductionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIntroductionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIntroductionFunc != nil {
		return atomic.LoadUint64(&m.GetIntroductionCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetNodePublicKey struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetNodePublicKeyExpectation
	expectationSeries []*NodeIntroProfileMockGetNodePublicKeyExpectation
}

type NodeIntroProfileMockGetNodePublicKeyExpectation struct {
	result *NodeIntroProfileMockGetNodePublicKeyResult
}

type NodeIntroProfileMockGetNodePublicKeyResult struct {
	r cryptkit.SignatureKeyHolder
}

//Expect specifies that invocation of NodeIntroProfile.GetNodePublicKey is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetNodePublicKey) Expect() *mNodeIntroProfileMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetNodePublicKey
func (m *mNodeIntroProfileMockGetNodePublicKey) Return(r cryptkit.SignatureKeyHolder) *NodeIntroProfileMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetNodePublicKey is expected once
func (m *mNodeIntroProfileMockGetNodePublicKey) ExpectOnce() *NodeIntroProfileMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetNodePublicKeyExpectation) Return(r cryptkit.SignatureKeyHolder) {
	e.result = &NodeIntroProfileMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetNodePublicKey method
func (m *mNodeIntroProfileMockGetNodePublicKey) Set(f func() (r cryptkit.SignatureKeyHolder)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetNodePublicKey() (r cryptkit.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of NodeIntroProfileMock.GetNodePublicKeyFunc invocations
func (m *NodeIntroProfileMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of NodeIntroProfileMock.GetNodePublicKey invocations
func (m *NodeIntroProfileMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetNodePublicKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodePublicKeyCounter) == uint64(len(m.GetNodePublicKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodePublicKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodePublicKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodePublicKeyFunc != nil {
		return atomic.LoadUint64(&m.GetNodePublicKeyCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetNodePublicKeyStore struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetNodePublicKeyStoreExpectation
	expectationSeries []*NodeIntroProfileMockGetNodePublicKeyStoreExpectation
}

type NodeIntroProfileMockGetNodePublicKeyStoreExpectation struct {
	result *NodeIntroProfileMockGetNodePublicKeyStoreResult
}

type NodeIntroProfileMockGetNodePublicKeyStoreResult struct {
	r cryptkit.PublicKeyStore
}

//Expect specifies that invocation of NodeIntroProfile.GetNodePublicKeyStore is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetNodePublicKeyStore) Expect() *mNodeIntroProfileMockGetNodePublicKeyStore {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetNodePublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetNodePublicKeyStore
func (m *mNodeIntroProfileMockGetNodePublicKeyStore) Return(r cryptkit.PublicKeyStore) *NodeIntroProfileMock {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetNodePublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetNodePublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetNodePublicKeyStore is expected once
func (m *mNodeIntroProfileMockGetNodePublicKeyStore) ExpectOnce() *NodeIntroProfileMockGetNodePublicKeyStoreExpectation {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetNodePublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetNodePublicKeyStoreExpectation) Return(r cryptkit.PublicKeyStore) {
	e.result = &NodeIntroProfileMockGetNodePublicKeyStoreResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetNodePublicKeyStore method
func (m *mNodeIntroProfileMockGetNodePublicKeyStore) Set(f func() (r cryptkit.PublicKeyStore)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyStoreFunc = f
	return m.mock
}

//GetNodePublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetNodePublicKeyStore() (r cryptkit.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyStoreCounter, 1)

	if len(m.GetNodePublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetNodePublicKeyStore.")
			return
		}

		result := m.GetNodePublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetNodePublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreMock.mainExpectation != nil {

		result := m.GetNodePublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetNodePublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetNodePublicKeyStore.")
		return
	}

	return m.GetNodePublicKeyStoreFunc()
}

//GetNodePublicKeyStoreMinimockCounter returns a count of NodeIntroProfileMock.GetNodePublicKeyStoreFunc invocations
func (m *NodeIntroProfileMock) GetNodePublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter)
}

//GetNodePublicKeyStoreMinimockPreCounter returns the value of NodeIntroProfileMock.GetNodePublicKeyStore invocations
func (m *NodeIntroProfileMock) GetNodePublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStorePreCounter)
}

//GetNodePublicKeyStoreFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetNodePublicKeyStoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodePublicKeyStoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter) == uint64(len(m.GetNodePublicKeyStoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodePublicKeyStoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodePublicKeyStoreFunc != nil {
		return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetPrimaryRole struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetPrimaryRoleExpectation
	expectationSeries []*NodeIntroProfileMockGetPrimaryRoleExpectation
}

type NodeIntroProfileMockGetPrimaryRoleExpectation struct {
	result *NodeIntroProfileMockGetPrimaryRoleResult
}

type NodeIntroProfileMockGetPrimaryRoleResult struct {
	r member.PrimaryRole
}

//Expect specifies that invocation of NodeIntroProfile.GetPrimaryRole is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetPrimaryRole) Expect() *mNodeIntroProfileMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetPrimaryRole
func (m *mNodeIntroProfileMockGetPrimaryRole) Return(r member.PrimaryRole) *NodeIntroProfileMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetPrimaryRole is expected once
func (m *mNodeIntroProfileMockGetPrimaryRole) ExpectOnce() *NodeIntroProfileMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetPrimaryRoleExpectation) Return(r member.PrimaryRole) {
	e.result = &NodeIntroProfileMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetPrimaryRole method
func (m *mNodeIntroProfileMockGetPrimaryRole) Set(f func() (r member.PrimaryRole)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetPrimaryRole() (r member.PrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of NodeIntroProfileMock.GetPrimaryRoleFunc invocations
func (m *NodeIntroProfileMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of NodeIntroProfileMock.GetPrimaryRole invocations
func (m *NodeIntroProfileMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetPrimaryRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPrimaryRoleCounter) == uint64(len(m.GetPrimaryRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPrimaryRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPrimaryRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPrimaryRoleFunc != nil {
		return atomic.LoadUint64(&m.GetPrimaryRoleCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetShortNodeID struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetShortNodeIDExpectation
	expectationSeries []*NodeIntroProfileMockGetShortNodeIDExpectation
}

type NodeIntroProfileMockGetShortNodeIDExpectation struct {
	result *NodeIntroProfileMockGetShortNodeIDResult
}

type NodeIntroProfileMockGetShortNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of NodeIntroProfile.GetShortNodeID is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetShortNodeID) Expect() *mNodeIntroProfileMockGetShortNodeID {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetShortNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetShortNodeID
func (m *mNodeIntroProfileMockGetShortNodeID) Return(r insolar.ShortNodeID) *NodeIntroProfileMock {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetShortNodeIDExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetShortNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetShortNodeID is expected once
func (m *mNodeIntroProfileMockGetShortNodeID) ExpectOnce() *NodeIntroProfileMockGetShortNodeIDExpectation {
	m.mock.GetShortNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetShortNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetShortNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &NodeIntroProfileMockGetShortNodeIDResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetShortNodeID method
func (m *mNodeIntroProfileMockGetShortNodeID) Set(f func() (r insolar.ShortNodeID)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetShortNodeIDFunc = f
	return m.mock
}

//GetShortNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetShortNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetShortNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetShortNodeIDCounter, 1)

	if len(m.GetShortNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetShortNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetShortNodeID.")
			return
		}

		result := m.GetShortNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetShortNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDMock.mainExpectation != nil {

		result := m.GetShortNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetShortNodeID")
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetShortNodeID.")
		return
	}

	return m.GetShortNodeIDFunc()
}

//GetShortNodeIDMinimockCounter returns a count of NodeIntroProfileMock.GetShortNodeIDFunc invocations
func (m *NodeIntroProfileMock) GetShortNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDCounter)
}

//GetShortNodeIDMinimockPreCounter returns the value of NodeIntroProfileMock.GetShortNodeID invocations
func (m *NodeIntroProfileMock) GetShortNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDPreCounter)
}

//GetShortNodeIDFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetShortNodeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetShortNodeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetShortNodeIDCounter) == uint64(len(m.GetShortNodeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetShortNodeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetShortNodeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetShortNodeIDFunc != nil {
		return atomic.LoadUint64(&m.GetShortNodeIDCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetSpecialRoles struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetSpecialRolesExpectation
	expectationSeries []*NodeIntroProfileMockGetSpecialRolesExpectation
}

type NodeIntroProfileMockGetSpecialRolesExpectation struct {
	result *NodeIntroProfileMockGetSpecialRolesResult
}

type NodeIntroProfileMockGetSpecialRolesResult struct {
	r member.SpecialRole
}

//Expect specifies that invocation of NodeIntroProfile.GetSpecialRoles is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetSpecialRoles) Expect() *mNodeIntroProfileMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetSpecialRoles
func (m *mNodeIntroProfileMockGetSpecialRoles) Return(r member.SpecialRole) *NodeIntroProfileMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetSpecialRoles is expected once
func (m *mNodeIntroProfileMockGetSpecialRoles) ExpectOnce() *NodeIntroProfileMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetSpecialRolesExpectation) Return(r member.SpecialRole) {
	e.result = &NodeIntroProfileMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetSpecialRoles method
func (m *mNodeIntroProfileMockGetSpecialRoles) Set(f func() (r member.SpecialRole)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetSpecialRoles() (r member.SpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of NodeIntroProfileMock.GetSpecialRolesFunc invocations
func (m *NodeIntroProfileMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of NodeIntroProfileMock.GetSpecialRoles invocations
func (m *NodeIntroProfileMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetSpecialRolesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSpecialRolesCounter) == uint64(len(m.GetSpecialRolesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSpecialRolesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSpecialRolesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSpecialRolesFunc != nil {
		return atomic.LoadUint64(&m.GetSpecialRolesCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockGetStartPower struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockGetStartPowerExpectation
	expectationSeries []*NodeIntroProfileMockGetStartPowerExpectation
}

type NodeIntroProfileMockGetStartPowerExpectation struct {
	result *NodeIntroProfileMockGetStartPowerResult
}

type NodeIntroProfileMockGetStartPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of NodeIntroProfile.GetStartPower is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockGetStartPower) Expect() *mNodeIntroProfileMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.GetStartPower
func (m *mNodeIntroProfileMockGetStartPower) Return(r member.Power) *NodeIntroProfileMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.GetStartPower is expected once
func (m *mNodeIntroProfileMockGetStartPower) ExpectOnce() *NodeIntroProfileMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockGetStartPowerExpectation) Return(r member.Power) {
	e.result = &NodeIntroProfileMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.GetStartPower method
func (m *mNodeIntroProfileMockGetStartPower) Set(f func() (r member.Power)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) GetStartPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of NodeIntroProfileMock.GetStartPowerFunc invocations
func (m *NodeIntroProfileMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of NodeIntroProfileMock.GetStartPower invocations
func (m *NodeIntroProfileMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) GetStartPowerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStartPowerCounter) == uint64(len(m.GetStartPowerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStartPowerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStartPowerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStartPowerFunc != nil {
		return atomic.LoadUint64(&m.GetStartPowerCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockHasIntroduction struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockHasIntroductionExpectation
	expectationSeries []*NodeIntroProfileMockHasIntroductionExpectation
}

type NodeIntroProfileMockHasIntroductionExpectation struct {
	result *NodeIntroProfileMockHasIntroductionResult
}

type NodeIntroProfileMockHasIntroductionResult struct {
	r bool
}

//Expect specifies that invocation of NodeIntroProfile.HasIntroduction is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockHasIntroduction) Expect() *mNodeIntroProfileMockHasIntroduction {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockHasIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeIntroProfile.HasIntroduction
func (m *mNodeIntroProfileMockHasIntroduction) Return(r bool) *NodeIntroProfileMock {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockHasIntroductionExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockHasIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.HasIntroduction is expected once
func (m *mNodeIntroProfileMockHasIntroduction) ExpectOnce() *NodeIntroProfileMockHasIntroductionExpectation {
	m.mock.HasIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockHasIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockHasIntroductionExpectation) Return(r bool) {
	e.result = &NodeIntroProfileMockHasIntroductionResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.HasIntroduction method
func (m *mNodeIntroProfileMockHasIntroduction) Set(f func() (r bool)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasIntroductionFunc = f
	return m.mock
}

//HasIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) HasIntroduction() (r bool) {
	counter := atomic.AddUint64(&m.HasIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.HasIntroductionCounter, 1)

	if len(m.HasIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.HasIntroduction.")
			return
		}

		result := m.HasIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.HasIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.HasIntroductionMock.mainExpectation != nil {

		result := m.HasIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.HasIntroduction")
		}

		r = result.r

		return
	}

	if m.HasIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.HasIntroduction.")
		return
	}

	return m.HasIntroductionFunc()
}

//HasIntroductionMinimockCounter returns a count of NodeIntroProfileMock.HasIntroductionFunc invocations
func (m *NodeIntroProfileMock) HasIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionCounter)
}

//HasIntroductionMinimockPreCounter returns the value of NodeIntroProfileMock.HasIntroduction invocations
func (m *NodeIntroProfileMock) HasIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionPreCounter)
}

//HasIntroductionFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) HasIntroductionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HasIntroductionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HasIntroductionCounter) == uint64(len(m.HasIntroductionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HasIntroductionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HasIntroductionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HasIntroductionFunc != nil {
		return atomic.LoadUint64(&m.HasIntroductionCounter) > 0
	}

	return true
}

type mNodeIntroProfileMockIsAcceptableHost struct {
	mock              *NodeIntroProfileMock
	mainExpectation   *NodeIntroProfileMockIsAcceptableHostExpectation
	expectationSeries []*NodeIntroProfileMockIsAcceptableHostExpectation
}

type NodeIntroProfileMockIsAcceptableHostExpectation struct {
	input  *NodeIntroProfileMockIsAcceptableHostInput
	result *NodeIntroProfileMockIsAcceptableHostResult
}

type NodeIntroProfileMockIsAcceptableHostInput struct {
	p endpoints.Inbound
}

type NodeIntroProfileMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of NodeIntroProfile.IsAcceptableHost is expected from 1 to Infinity times
func (m *mNodeIntroProfileMockIsAcceptableHost) Expect(p endpoints.Inbound) *mNodeIntroProfileMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &NodeIntroProfileMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of NodeIntroProfile.IsAcceptableHost
func (m *mNodeIntroProfileMockIsAcceptableHost) Return(r bool) *NodeIntroProfileMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeIntroProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &NodeIntroProfileMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeIntroProfile.IsAcceptableHost is expected once
func (m *mNodeIntroProfileMockIsAcceptableHost) ExpectOnce(p endpoints.Inbound) *NodeIntroProfileMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &NodeIntroProfileMockIsAcceptableHostExpectation{}
	expectation.input = &NodeIntroProfileMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeIntroProfileMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &NodeIntroProfileMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of NodeIntroProfile.IsAcceptableHost method
func (m *mNodeIntroProfileMockIsAcceptableHost) Set(f func(p endpoints.Inbound) (r bool)) *NodeIntroProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile interface
func (m *NodeIntroProfileMock) IsAcceptableHost(p endpoints.Inbound) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeIntroProfileMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeIntroProfileMockIsAcceptableHostInput{p}, "NodeIntroProfile.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeIntroProfileMockIsAcceptableHostInput{p}, "NodeIntroProfile.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeIntroProfileMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to NodeIntroProfileMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of NodeIntroProfileMock.IsAcceptableHostFunc invocations
func (m *NodeIntroProfileMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of NodeIntroProfileMock.IsAcceptableHost invocations
func (m *NodeIntroProfileMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *NodeIntroProfileMock) IsAcceptableHostFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsAcceptableHostCounter) == uint64(len(m.IsAcceptableHostMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsAcceptableHostMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsAcceptableHostCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsAcceptableHostFunc != nil {
		return atomic.LoadUint64(&m.IsAcceptableHostCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeIntroProfileMock) ValidateCallCounters() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetAnnouncementSignature")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetDefaultEndpoint")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetNodePublicKey")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetNodePublicKeyStore")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetPrimaryRole")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetShortNodeID")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.IsAcceptableHost")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeIntroProfileMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeIntroProfileMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeIntroProfileMock) MinimockFinish() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetAnnouncementSignature")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetDefaultEndpoint")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetNodePublicKey")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetNodePublicKeyStore")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetPrimaryRole")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetShortNodeID")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to NodeIntroProfileMock.IsAcceptableHost")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeIntroProfileMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeIntroProfileMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetAnnouncementSignatureFinished()
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetIntroductionFinished()
		ok = ok && m.GetNodePublicKeyFinished()
		ok = ok && m.GetNodePublicKeyStoreFinished()
		ok = ok && m.GetPrimaryRoleFinished()
		ok = ok && m.GetShortNodeIDFinished()
		ok = ok && m.GetSpecialRolesFinished()
		ok = ok && m.GetStartPowerFinished()
		ok = ok && m.HasIntroductionFinished()
		ok = ok && m.IsAcceptableHostFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAnnouncementSignatureFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetAnnouncementSignature")
			}

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetDefaultEndpoint")
			}

			if !m.GetIntroductionFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetIntroduction")
			}

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetNodePublicKey")
			}

			if !m.GetNodePublicKeyStoreFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetNodePublicKeyStore")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetPrimaryRole")
			}

			if !m.GetShortNodeIDFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetShortNodeID")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetSpecialRoles")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.GetStartPower")
			}

			if !m.HasIntroductionFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.HasIntroduction")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to NodeIntroProfileMock.IsAcceptableHost")
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
func (m *NodeIntroProfileMock) AllMocksCalled() bool {

	if !m.GetAnnouncementSignatureFinished() {
		return false
	}

	if !m.GetDefaultEndpointFinished() {
		return false
	}

	if !m.GetIntroductionFinished() {
		return false
	}

	if !m.GetNodePublicKeyFinished() {
		return false
	}

	if !m.GetNodePublicKeyStoreFinished() {
		return false
	}

	if !m.GetPrimaryRoleFinished() {
		return false
	}

	if !m.GetShortNodeIDFinished() {
		return false
	}

	if !m.GetSpecialRolesFinished() {
		return false
	}

	if !m.GetStartPowerFinished() {
		return false
	}

	if !m.HasIntroductionFinished() {
		return false
	}

	if !m.IsAcceptableHostFinished() {
		return false
	}

	return true
}

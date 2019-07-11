package api

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeProfile" can be found in github.com/insolar/insolar/network/consensus/gcpv2/common
*/
import (
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/network/consensus/common"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeProfileMock implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile
type NodeProfileMock struct {
	t minimock.Tester

	GetAnnouncementSignatureFunc       func() (r cryptography_containers.SignatureHolder)
	GetAnnouncementSignatureCounter    uint64
	GetAnnouncementSignaturePreCounter uint64
	GetAnnouncementSignatureMock       mNodeProfileMockGetAnnouncementSignature

	GetDeclaredPowerFunc       func() (r MemberPower)
	GetDeclaredPowerCounter    uint64
	GetDeclaredPowerPreCounter uint64
	GetDeclaredPowerMock       mNodeProfileMockGetDeclaredPower

	GetDefaultEndpointFunc       func() (r endpoints.NodeEndpoint)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mNodeProfileMockGetDefaultEndpoint

	GetIndexFunc       func() (r int)
	GetIndexCounter    uint64
	GetIndexPreCounter uint64
	GetIndexMock       mNodeProfileMockGetIndex

	GetIntroductionFunc       func() (r NodeIntroduction)
	GetIntroductionCounter    uint64
	GetIntroductionPreCounter uint64
	GetIntroductionMock       mNodeProfileMockGetIntroduction

	GetNodePublicKeyFunc       func() (r cryptography_containers.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mNodeProfileMockGetNodePublicKey

	GetNodePublicKeyStoreFunc       func() (r cryptography_containers.PublicKeyStore)
	GetNodePublicKeyStoreCounter    uint64
	GetNodePublicKeyStorePreCounter uint64
	GetNodePublicKeyStoreMock       mNodeProfileMockGetNodePublicKeyStore

	GetOpModeFunc       func() (r MemberOpMode)
	GetOpModeCounter    uint64
	GetOpModePreCounter uint64
	GetOpModeMock       mNodeProfileMockGetOpMode

	GetPrimaryRoleFunc       func() (r NodePrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mNodeProfileMockGetPrimaryRole

	GetShortNodeIDFunc       func() (r common.ShortNodeID)
	GetShortNodeIDCounter    uint64
	GetShortNodeIDPreCounter uint64
	GetShortNodeIDMock       mNodeProfileMockGetShortNodeID

	GetSignatureVerifierFunc       func() (r cryptography_containers.SignatureVerifier)
	GetSignatureVerifierCounter    uint64
	GetSignatureVerifierPreCounter uint64
	GetSignatureVerifierMock       mNodeProfileMockGetSignatureVerifier

	GetSpecialRolesFunc       func() (r NodeSpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mNodeProfileMockGetSpecialRoles

	GetStartPowerFunc       func() (r MemberPower)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mNodeProfileMockGetStartPower

	HasIntroductionFunc       func() (r bool)
	HasIntroductionCounter    uint64
	HasIntroductionPreCounter uint64
	HasIntroductionMock       mNodeProfileMockHasIntroduction

	IsAcceptableHostFunc       func(p endpoints.HostIdentityHolder) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mNodeProfileMockIsAcceptableHost

	IsJoinerFunc       func() (r bool)
	IsJoinerCounter    uint64
	IsJoinerPreCounter uint64
	IsJoinerMock       mNodeProfileMockIsJoiner
}

//NewNodeProfileMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile
func NewNodeProfileMock(t minimock.Tester) *NodeProfileMock {
	m := &NodeProfileMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAnnouncementSignatureMock = mNodeProfileMockGetAnnouncementSignature{mock: m}
	m.GetDeclaredPowerMock = mNodeProfileMockGetDeclaredPower{mock: m}
	m.GetDefaultEndpointMock = mNodeProfileMockGetDefaultEndpoint{mock: m}
	m.GetIndexMock = mNodeProfileMockGetIndex{mock: m}
	m.GetIntroductionMock = mNodeProfileMockGetIntroduction{mock: m}
	m.GetNodePublicKeyMock = mNodeProfileMockGetNodePublicKey{mock: m}
	m.GetNodePublicKeyStoreMock = mNodeProfileMockGetNodePublicKeyStore{mock: m}
	m.GetOpModeMock = mNodeProfileMockGetOpMode{mock: m}
	m.GetPrimaryRoleMock = mNodeProfileMockGetPrimaryRole{mock: m}
	m.GetShortNodeIDMock = mNodeProfileMockGetShortNodeID{mock: m}
	m.GetSignatureVerifierMock = mNodeProfileMockGetSignatureVerifier{mock: m}
	m.GetSpecialRolesMock = mNodeProfileMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mNodeProfileMockGetStartPower{mock: m}
	m.HasIntroductionMock = mNodeProfileMockHasIntroduction{mock: m}
	m.IsAcceptableHostMock = mNodeProfileMockIsAcceptableHost{mock: m}
	m.IsJoinerMock = mNodeProfileMockIsJoiner{mock: m}

	return m
}

type mNodeProfileMockGetAnnouncementSignature struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetAnnouncementSignatureExpectation
	expectationSeries []*NodeProfileMockGetAnnouncementSignatureExpectation
}

type NodeProfileMockGetAnnouncementSignatureExpectation struct {
	result *NodeProfileMockGetAnnouncementSignatureResult
}

type NodeProfileMockGetAnnouncementSignatureResult struct {
	r cryptography_containers.SignatureHolder
}

//Expect specifies that invocation of NodeProfile.GetAnnouncementSignature is expected from 1 to Infinity times
func (m *mNodeProfileMockGetAnnouncementSignature) Expect() *mNodeProfileMockGetAnnouncementSignature {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetAnnouncementSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetAnnouncementSignature
func (m *mNodeProfileMockGetAnnouncementSignature) Return(r cryptography_containers.SignatureHolder) *NodeProfileMock {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetAnnouncementSignatureExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetAnnouncementSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetAnnouncementSignature is expected once
func (m *mNodeProfileMockGetAnnouncementSignature) ExpectOnce() *NodeProfileMockGetAnnouncementSignatureExpectation {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetAnnouncementSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetAnnouncementSignatureExpectation) Return(r cryptography_containers.SignatureHolder) {
	e.result = &NodeProfileMockGetAnnouncementSignatureResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetAnnouncementSignature method
func (m *mNodeProfileMockGetAnnouncementSignature) Set(f func() (r cryptography_containers.SignatureHolder)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAnnouncementSignatureFunc = f
	return m.mock
}

//GetAnnouncementSignature implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetAnnouncementSignature() (r cryptography_containers.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetAnnouncementSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetAnnouncementSignatureCounter, 1)

	if len(m.GetAnnouncementSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAnnouncementSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetAnnouncementSignature.")
			return
		}

		result := m.GetAnnouncementSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetAnnouncementSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureMock.mainExpectation != nil {

		result := m.GetAnnouncementSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetAnnouncementSignature")
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetAnnouncementSignature.")
		return
	}

	return m.GetAnnouncementSignatureFunc()
}

//GetAnnouncementSignatureMinimockCounter returns a count of NodeProfileMock.GetAnnouncementSignatureFunc invocations
func (m *NodeProfileMock) GetAnnouncementSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter)
}

//GetAnnouncementSignatureMinimockPreCounter returns the value of NodeProfileMock.GetAnnouncementSignature invocations
func (m *NodeProfileMock) GetAnnouncementSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignaturePreCounter)
}

//GetAnnouncementSignatureFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetAnnouncementSignatureFinished() bool {
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

type mNodeProfileMockGetDeclaredPower struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetDeclaredPowerExpectation
	expectationSeries []*NodeProfileMockGetDeclaredPowerExpectation
}

type NodeProfileMockGetDeclaredPowerExpectation struct {
	result *NodeProfileMockGetDeclaredPowerResult
}

type NodeProfileMockGetDeclaredPowerResult struct {
	r MemberPower
}

//Expect specifies that invocation of NodeProfile.GetDeclaredPower is expected from 1 to Infinity times
func (m *mNodeProfileMockGetDeclaredPower) Expect() *mNodeProfileMockGetDeclaredPower {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetDeclaredPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetDeclaredPower
func (m *mNodeProfileMockGetDeclaredPower) Return(r MemberPower) *NodeProfileMock {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetDeclaredPowerExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetDeclaredPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetDeclaredPower is expected once
func (m *mNodeProfileMockGetDeclaredPower) ExpectOnce() *NodeProfileMockGetDeclaredPowerExpectation {
	m.mock.GetDeclaredPowerFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetDeclaredPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetDeclaredPowerExpectation) Return(r MemberPower) {
	e.result = &NodeProfileMockGetDeclaredPowerResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetDeclaredPower method
func (m *mNodeProfileMockGetDeclaredPower) Set(f func() (r MemberPower)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDeclaredPowerFunc = f
	return m.mock
}

//GetDeclaredPower implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetDeclaredPower() (r MemberPower) {
	counter := atomic.AddUint64(&m.GetDeclaredPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetDeclaredPowerCounter, 1)

	if len(m.GetDeclaredPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDeclaredPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetDeclaredPower.")
			return
		}

		result := m.GetDeclaredPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetDeclaredPower")
			return
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerMock.mainExpectation != nil {

		result := m.GetDeclaredPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetDeclaredPower")
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetDeclaredPower.")
		return
	}

	return m.GetDeclaredPowerFunc()
}

//GetDeclaredPowerMinimockCounter returns a count of NodeProfileMock.GetDeclaredPowerFunc invocations
func (m *NodeProfileMock) GetDeclaredPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerCounter)
}

//GetDeclaredPowerMinimockPreCounter returns the value of NodeProfileMock.GetDeclaredPower invocations
func (m *NodeProfileMock) GetDeclaredPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerPreCounter)
}

//GetDeclaredPowerFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetDeclaredPowerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDeclaredPowerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDeclaredPowerCounter) == uint64(len(m.GetDeclaredPowerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDeclaredPowerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDeclaredPowerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDeclaredPowerFunc != nil {
		return atomic.LoadUint64(&m.GetDeclaredPowerCounter) > 0
	}

	return true
}

type mNodeProfileMockGetDefaultEndpoint struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetDefaultEndpointExpectation
	expectationSeries []*NodeProfileMockGetDefaultEndpointExpectation
}

type NodeProfileMockGetDefaultEndpointExpectation struct {
	result *NodeProfileMockGetDefaultEndpointResult
}

type NodeProfileMockGetDefaultEndpointResult struct {
	r endpoints.NodeEndpoint
}

//Expect specifies that invocation of NodeProfile.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mNodeProfileMockGetDefaultEndpoint) Expect() *mNodeProfileMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetDefaultEndpoint
func (m *mNodeProfileMockGetDefaultEndpoint) Return(r endpoints.NodeEndpoint) *NodeProfileMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetDefaultEndpoint is expected once
func (m *mNodeProfileMockGetDefaultEndpoint) ExpectOnce() *NodeProfileMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetDefaultEndpointExpectation) Return(r endpoints.NodeEndpoint) {
	e.result = &NodeProfileMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetDefaultEndpoint method
func (m *mNodeProfileMockGetDefaultEndpoint) Set(f func() (r endpoints.NodeEndpoint)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetDefaultEndpoint() (r endpoints.NodeEndpoint) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of NodeProfileMock.GetDefaultEndpointFunc invocations
func (m *NodeProfileMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of NodeProfileMock.GetDefaultEndpoint invocations
func (m *NodeProfileMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetDefaultEndpointFinished() bool {
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

type mNodeProfileMockGetIndex struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetIndexExpectation
	expectationSeries []*NodeProfileMockGetIndexExpectation
}

type NodeProfileMockGetIndexExpectation struct {
	result *NodeProfileMockGetIndexResult
}

type NodeProfileMockGetIndexResult struct {
	r int
}

//Expect specifies that invocation of NodeProfile.GetIndex is expected from 1 to Infinity times
func (m *mNodeProfileMockGetIndex) Expect() *mNodeProfileMockGetIndex {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetIndexExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetIndex
func (m *mNodeProfileMockGetIndex) Return(r int) *NodeProfileMock {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetIndexExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetIndex is expected once
func (m *mNodeProfileMockGetIndex) ExpectOnce() *NodeProfileMockGetIndexExpectation {
	m.mock.GetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetIndexExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetIndexExpectation) Return(r int) {
	e.result = &NodeProfileMockGetIndexResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetIndex method
func (m *mNodeProfileMockGetIndex) Set(f func() (r int)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIndexFunc = f
	return m.mock
}

//GetIndex implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetIndex() (r int) {
	counter := atomic.AddUint64(&m.GetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.GetIndexCounter, 1)

	if len(m.GetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetIndex.")
			return
		}

		result := m.GetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetIndex")
			return
		}

		r = result.r

		return
	}

	if m.GetIndexMock.mainExpectation != nil {

		result := m.GetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetIndex")
		}

		r = result.r

		return
	}

	if m.GetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetIndex.")
		return
	}

	return m.GetIndexFunc()
}

//GetIndexMinimockCounter returns a count of NodeProfileMock.GetIndexFunc invocations
func (m *NodeProfileMock) GetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexCounter)
}

//GetIndexMinimockPreCounter returns the value of NodeProfileMock.GetIndex invocations
func (m *NodeProfileMock) GetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexPreCounter)
}

//GetIndexFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIndexCounter) == uint64(len(m.GetIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIndexFunc != nil {
		return atomic.LoadUint64(&m.GetIndexCounter) > 0
	}

	return true
}

type mNodeProfileMockGetIntroduction struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetIntroductionExpectation
	expectationSeries []*NodeProfileMockGetIntroductionExpectation
}

type NodeProfileMockGetIntroductionExpectation struct {
	result *NodeProfileMockGetIntroductionResult
}

type NodeProfileMockGetIntroductionResult struct {
	r NodeIntroduction
}

//Expect specifies that invocation of NodeProfile.GetIntroduction is expected from 1 to Infinity times
func (m *mNodeProfileMockGetIntroduction) Expect() *mNodeProfileMockGetIntroduction {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetIntroduction
func (m *mNodeProfileMockGetIntroduction) Return(r NodeIntroduction) *NodeProfileMock {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetIntroductionExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetIntroduction is expected once
func (m *mNodeProfileMockGetIntroduction) ExpectOnce() *NodeProfileMockGetIntroductionExpectation {
	m.mock.GetIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetIntroductionExpectation) Return(r NodeIntroduction) {
	e.result = &NodeProfileMockGetIntroductionResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetIntroduction method
func (m *mNodeProfileMockGetIntroduction) Set(f func() (r NodeIntroduction)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIntroductionFunc = f
	return m.mock
}

//GetIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetIntroduction() (r NodeIntroduction) {
	counter := atomic.AddUint64(&m.GetIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.GetIntroductionCounter, 1)

	if len(m.GetIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetIntroduction.")
			return
		}

		result := m.GetIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.GetIntroductionMock.mainExpectation != nil {

		result := m.GetIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetIntroduction")
		}

		r = result.r

		return
	}

	if m.GetIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetIntroduction.")
		return
	}

	return m.GetIntroductionFunc()
}

//GetIntroductionMinimockCounter returns a count of NodeProfileMock.GetIntroductionFunc invocations
func (m *NodeProfileMock) GetIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionCounter)
}

//GetIntroductionMinimockPreCounter returns the value of NodeProfileMock.GetIntroduction invocations
func (m *NodeProfileMock) GetIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionPreCounter)
}

//GetIntroductionFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetIntroductionFinished() bool {
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

type mNodeProfileMockGetNodePublicKey struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetNodePublicKeyExpectation
	expectationSeries []*NodeProfileMockGetNodePublicKeyExpectation
}

type NodeProfileMockGetNodePublicKeyExpectation struct {
	result *NodeProfileMockGetNodePublicKeyResult
}

type NodeProfileMockGetNodePublicKeyResult struct {
	r cryptography_containers.SignatureKeyHolder
}

//Expect specifies that invocation of NodeProfile.GetNodePublicKey is expected from 1 to Infinity times
func (m *mNodeProfileMockGetNodePublicKey) Expect() *mNodeProfileMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetNodePublicKey
func (m *mNodeProfileMockGetNodePublicKey) Return(r cryptography_containers.SignatureKeyHolder) *NodeProfileMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetNodePublicKey is expected once
func (m *mNodeProfileMockGetNodePublicKey) ExpectOnce() *NodeProfileMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetNodePublicKeyExpectation) Return(r cryptography_containers.SignatureKeyHolder) {
	e.result = &NodeProfileMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetNodePublicKey method
func (m *mNodeProfileMockGetNodePublicKey) Set(f func() (r cryptography_containers.SignatureKeyHolder)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetNodePublicKey() (r cryptography_containers.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of NodeProfileMock.GetNodePublicKeyFunc invocations
func (m *NodeProfileMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of NodeProfileMock.GetNodePublicKey invocations
func (m *NodeProfileMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetNodePublicKeyFinished() bool {
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

type mNodeProfileMockGetNodePublicKeyStore struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetNodePublicKeyStoreExpectation
	expectationSeries []*NodeProfileMockGetNodePublicKeyStoreExpectation
}

type NodeProfileMockGetNodePublicKeyStoreExpectation struct {
	result *NodeProfileMockGetNodePublicKeyStoreResult
}

type NodeProfileMockGetNodePublicKeyStoreResult struct {
	r cryptography_containers.PublicKeyStore
}

//Expect specifies that invocation of NodeProfile.GetNodePublicKeyStore is expected from 1 to Infinity times
func (m *mNodeProfileMockGetNodePublicKeyStore) Expect() *mNodeProfileMockGetNodePublicKeyStore {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetNodePublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetNodePublicKeyStore
func (m *mNodeProfileMockGetNodePublicKeyStore) Return(r cryptography_containers.PublicKeyStore) *NodeProfileMock {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetNodePublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetNodePublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetNodePublicKeyStore is expected once
func (m *mNodeProfileMockGetNodePublicKeyStore) ExpectOnce() *NodeProfileMockGetNodePublicKeyStoreExpectation {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetNodePublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetNodePublicKeyStoreExpectation) Return(r cryptography_containers.PublicKeyStore) {
	e.result = &NodeProfileMockGetNodePublicKeyStoreResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetNodePublicKeyStore method
func (m *mNodeProfileMockGetNodePublicKeyStore) Set(f func() (r cryptography_containers.PublicKeyStore)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyStoreFunc = f
	return m.mock
}

//GetNodePublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetNodePublicKeyStore() (r cryptography_containers.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyStoreCounter, 1)

	if len(m.GetNodePublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetNodePublicKeyStore.")
			return
		}

		result := m.GetNodePublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetNodePublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreMock.mainExpectation != nil {

		result := m.GetNodePublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetNodePublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetNodePublicKeyStore.")
		return
	}

	return m.GetNodePublicKeyStoreFunc()
}

//GetNodePublicKeyStoreMinimockCounter returns a count of NodeProfileMock.GetNodePublicKeyStoreFunc invocations
func (m *NodeProfileMock) GetNodePublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter)
}

//GetNodePublicKeyStoreMinimockPreCounter returns the value of NodeProfileMock.GetNodePublicKeyStore invocations
func (m *NodeProfileMock) GetNodePublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStorePreCounter)
}

//GetNodePublicKeyStoreFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetNodePublicKeyStoreFinished() bool {
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

type mNodeProfileMockGetOpMode struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetOpModeExpectation
	expectationSeries []*NodeProfileMockGetOpModeExpectation
}

type NodeProfileMockGetOpModeExpectation struct {
	result *NodeProfileMockGetOpModeResult
}

type NodeProfileMockGetOpModeResult struct {
	r MemberOpMode
}

//Expect specifies that invocation of NodeProfile.GetOpMode is expected from 1 to Infinity times
func (m *mNodeProfileMockGetOpMode) Expect() *mNodeProfileMockGetOpMode {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetOpModeExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetOpMode
func (m *mNodeProfileMockGetOpMode) Return(r MemberOpMode) *NodeProfileMock {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetOpModeExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetOpModeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetOpMode is expected once
func (m *mNodeProfileMockGetOpMode) ExpectOnce() *NodeProfileMockGetOpModeExpectation {
	m.mock.GetOpModeFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetOpModeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetOpModeExpectation) Return(r MemberOpMode) {
	e.result = &NodeProfileMockGetOpModeResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetOpMode method
func (m *mNodeProfileMockGetOpMode) Set(f func() (r MemberOpMode)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOpModeFunc = f
	return m.mock
}

//GetOpMode implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetOpMode() (r MemberOpMode) {
	counter := atomic.AddUint64(&m.GetOpModePreCounter, 1)
	defer atomic.AddUint64(&m.GetOpModeCounter, 1)

	if len(m.GetOpModeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOpModeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetOpMode.")
			return
		}

		result := m.GetOpModeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetOpMode")
			return
		}

		r = result.r

		return
	}

	if m.GetOpModeMock.mainExpectation != nil {

		result := m.GetOpModeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetOpMode")
		}

		r = result.r

		return
	}

	if m.GetOpModeFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetOpMode.")
		return
	}

	return m.GetOpModeFunc()
}

//GetOpModeMinimockCounter returns a count of NodeProfileMock.GetOpModeFunc invocations
func (m *NodeProfileMock) GetOpModeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModeCounter)
}

//GetOpModeMinimockPreCounter returns the value of NodeProfileMock.GetOpMode invocations
func (m *NodeProfileMock) GetOpModeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModePreCounter)
}

//GetOpModeFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetOpModeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOpModeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOpModeCounter) == uint64(len(m.GetOpModeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOpModeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOpModeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOpModeFunc != nil {
		return atomic.LoadUint64(&m.GetOpModeCounter) > 0
	}

	return true
}

type mNodeProfileMockGetPrimaryRole struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetPrimaryRoleExpectation
	expectationSeries []*NodeProfileMockGetPrimaryRoleExpectation
}

type NodeProfileMockGetPrimaryRoleExpectation struct {
	result *NodeProfileMockGetPrimaryRoleResult
}

type NodeProfileMockGetPrimaryRoleResult struct {
	r NodePrimaryRole
}

//Expect specifies that invocation of NodeProfile.GetPrimaryRole is expected from 1 to Infinity times
func (m *mNodeProfileMockGetPrimaryRole) Expect() *mNodeProfileMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetPrimaryRole
func (m *mNodeProfileMockGetPrimaryRole) Return(r NodePrimaryRole) *NodeProfileMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetPrimaryRole is expected once
func (m *mNodeProfileMockGetPrimaryRole) ExpectOnce() *NodeProfileMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetPrimaryRoleExpectation) Return(r NodePrimaryRole) {
	e.result = &NodeProfileMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetPrimaryRole method
func (m *mNodeProfileMockGetPrimaryRole) Set(f func() (r NodePrimaryRole)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetPrimaryRole() (r NodePrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of NodeProfileMock.GetPrimaryRoleFunc invocations
func (m *NodeProfileMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of NodeProfileMock.GetPrimaryRole invocations
func (m *NodeProfileMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetPrimaryRoleFinished() bool {
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

type mNodeProfileMockGetShortNodeID struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetShortNodeIDExpectation
	expectationSeries []*NodeProfileMockGetShortNodeIDExpectation
}

type NodeProfileMockGetShortNodeIDExpectation struct {
	result *NodeProfileMockGetShortNodeIDResult
}

type NodeProfileMockGetShortNodeIDResult struct {
	r common.ShortNodeID
}

//Expect specifies that invocation of NodeProfile.GetShortNodeID is expected from 1 to Infinity times
func (m *mNodeProfileMockGetShortNodeID) Expect() *mNodeProfileMockGetShortNodeID {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetShortNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetShortNodeID
func (m *mNodeProfileMockGetShortNodeID) Return(r common.ShortNodeID) *NodeProfileMock {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetShortNodeIDExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetShortNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetShortNodeID is expected once
func (m *mNodeProfileMockGetShortNodeID) ExpectOnce() *NodeProfileMockGetShortNodeIDExpectation {
	m.mock.GetShortNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetShortNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetShortNodeIDExpectation) Return(r common.ShortNodeID) {
	e.result = &NodeProfileMockGetShortNodeIDResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetShortNodeID method
func (m *mNodeProfileMockGetShortNodeID) Set(f func() (r common.ShortNodeID)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetShortNodeIDFunc = f
	return m.mock
}

//GetShortNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetShortNodeID() (r common.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetShortNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetShortNodeIDCounter, 1)

	if len(m.GetShortNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetShortNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetShortNodeID.")
			return
		}

		result := m.GetShortNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetShortNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDMock.mainExpectation != nil {

		result := m.GetShortNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetShortNodeID")
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetShortNodeID.")
		return
	}

	return m.GetShortNodeIDFunc()
}

//GetShortNodeIDMinimockCounter returns a count of NodeProfileMock.GetShortNodeIDFunc invocations
func (m *NodeProfileMock) GetShortNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDCounter)
}

//GetShortNodeIDMinimockPreCounter returns the value of NodeProfileMock.GetShortNodeID invocations
func (m *NodeProfileMock) GetShortNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDPreCounter)
}

//GetShortNodeIDFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetShortNodeIDFinished() bool {
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

type mNodeProfileMockGetSignatureVerifier struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetSignatureVerifierExpectation
	expectationSeries []*NodeProfileMockGetSignatureVerifierExpectation
}

type NodeProfileMockGetSignatureVerifierExpectation struct {
	result *NodeProfileMockGetSignatureVerifierResult
}

type NodeProfileMockGetSignatureVerifierResult struct {
	r cryptography_containers.SignatureVerifier
}

//Expect specifies that invocation of NodeProfile.GetSignatureVerifier is expected from 1 to Infinity times
func (m *mNodeProfileMockGetSignatureVerifier) Expect() *mNodeProfileMockGetSignatureVerifier {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetSignatureVerifierExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetSignatureVerifier
func (m *mNodeProfileMockGetSignatureVerifier) Return(r cryptography_containers.SignatureVerifier) *NodeProfileMock {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetSignatureVerifierExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetSignatureVerifierResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetSignatureVerifier is expected once
func (m *mNodeProfileMockGetSignatureVerifier) ExpectOnce() *NodeProfileMockGetSignatureVerifierExpectation {
	m.mock.GetSignatureVerifierFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetSignatureVerifierExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetSignatureVerifierExpectation) Return(r cryptography_containers.SignatureVerifier) {
	e.result = &NodeProfileMockGetSignatureVerifierResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetSignatureVerifier method
func (m *mNodeProfileMockGetSignatureVerifier) Set(f func() (r cryptography_containers.SignatureVerifier)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureVerifierFunc = f
	return m.mock
}

//GetSignatureVerifier implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetSignatureVerifier() (r cryptography_containers.SignatureVerifier) {
	counter := atomic.AddUint64(&m.GetSignatureVerifierPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureVerifierCounter, 1)

	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureVerifierMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetSignatureVerifier.")
			return
		}

		result := m.GetSignatureVerifierMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetSignatureVerifier")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierMock.mainExpectation != nil {

		result := m.GetSignatureVerifierMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetSignatureVerifier")
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetSignatureVerifier.")
		return
	}

	return m.GetSignatureVerifierFunc()
}

//GetSignatureVerifierMinimockCounter returns a count of NodeProfileMock.GetSignatureVerifierFunc invocations
func (m *NodeProfileMock) GetSignatureVerifierMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierCounter)
}

//GetSignatureVerifierMinimockPreCounter returns the value of NodeProfileMock.GetSignatureVerifier invocations
func (m *NodeProfileMock) GetSignatureVerifierMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierPreCounter)
}

//GetSignatureVerifierFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetSignatureVerifierFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureVerifierCounter) == uint64(len(m.GetSignatureVerifierMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureVerifierMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureVerifierCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureVerifierFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureVerifierCounter) > 0
	}

	return true
}

type mNodeProfileMockGetSpecialRoles struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetSpecialRolesExpectation
	expectationSeries []*NodeProfileMockGetSpecialRolesExpectation
}

type NodeProfileMockGetSpecialRolesExpectation struct {
	result *NodeProfileMockGetSpecialRolesResult
}

type NodeProfileMockGetSpecialRolesResult struct {
	r NodeSpecialRole
}

//Expect specifies that invocation of NodeProfile.GetSpecialRoles is expected from 1 to Infinity times
func (m *mNodeProfileMockGetSpecialRoles) Expect() *mNodeProfileMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetSpecialRoles
func (m *mNodeProfileMockGetSpecialRoles) Return(r NodeSpecialRole) *NodeProfileMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetSpecialRoles is expected once
func (m *mNodeProfileMockGetSpecialRoles) ExpectOnce() *NodeProfileMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetSpecialRolesExpectation) Return(r NodeSpecialRole) {
	e.result = &NodeProfileMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetSpecialRoles method
func (m *mNodeProfileMockGetSpecialRoles) Set(f func() (r NodeSpecialRole)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetSpecialRoles() (r NodeSpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of NodeProfileMock.GetSpecialRolesFunc invocations
func (m *NodeProfileMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of NodeProfileMock.GetSpecialRoles invocations
func (m *NodeProfileMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetSpecialRolesFinished() bool {
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

type mNodeProfileMockGetStartPower struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockGetStartPowerExpectation
	expectationSeries []*NodeProfileMockGetStartPowerExpectation
}

type NodeProfileMockGetStartPowerExpectation struct {
	result *NodeProfileMockGetStartPowerResult
}

type NodeProfileMockGetStartPowerResult struct {
	r MemberPower
}

//Expect specifies that invocation of NodeProfile.GetStartPower is expected from 1 to Infinity times
func (m *mNodeProfileMockGetStartPower) Expect() *mNodeProfileMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.GetStartPower
func (m *mNodeProfileMockGetStartPower) Return(r MemberPower) *NodeProfileMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.GetStartPower is expected once
func (m *mNodeProfileMockGetStartPower) ExpectOnce() *NodeProfileMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockGetStartPowerExpectation) Return(r MemberPower) {
	e.result = &NodeProfileMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of NodeProfile.GetStartPower method
func (m *mNodeProfileMockGetStartPower) Set(f func() (r MemberPower)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) GetStartPower() (r MemberPower) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of NodeProfileMock.GetStartPowerFunc invocations
func (m *NodeProfileMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of NodeProfileMock.GetStartPower invocations
func (m *NodeProfileMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) GetStartPowerFinished() bool {
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

type mNodeProfileMockHasIntroduction struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockHasIntroductionExpectation
	expectationSeries []*NodeProfileMockHasIntroductionExpectation
}

type NodeProfileMockHasIntroductionExpectation struct {
	result *NodeProfileMockHasIntroductionResult
}

type NodeProfileMockHasIntroductionResult struct {
	r bool
}

//Expect specifies that invocation of NodeProfile.HasIntroduction is expected from 1 to Infinity times
func (m *mNodeProfileMockHasIntroduction) Expect() *mNodeProfileMockHasIntroduction {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockHasIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.HasIntroduction
func (m *mNodeProfileMockHasIntroduction) Return(r bool) *NodeProfileMock {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockHasIntroductionExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockHasIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.HasIntroduction is expected once
func (m *mNodeProfileMockHasIntroduction) ExpectOnce() *NodeProfileMockHasIntroductionExpectation {
	m.mock.HasIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockHasIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockHasIntroductionExpectation) Return(r bool) {
	e.result = &NodeProfileMockHasIntroductionResult{r}
}

//Set uses given function f as a mock of NodeProfile.HasIntroduction method
func (m *mNodeProfileMockHasIntroduction) Set(f func() (r bool)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasIntroductionFunc = f
	return m.mock
}

//HasIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) HasIntroduction() (r bool) {
	counter := atomic.AddUint64(&m.HasIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.HasIntroductionCounter, 1)

	if len(m.HasIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.HasIntroduction.")
			return
		}

		result := m.HasIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.HasIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.HasIntroductionMock.mainExpectation != nil {

		result := m.HasIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.HasIntroduction")
		}

		r = result.r

		return
	}

	if m.HasIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.HasIntroduction.")
		return
	}

	return m.HasIntroductionFunc()
}

//HasIntroductionMinimockCounter returns a count of NodeProfileMock.HasIntroductionFunc invocations
func (m *NodeProfileMock) HasIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionCounter)
}

//HasIntroductionMinimockPreCounter returns the value of NodeProfileMock.HasIntroduction invocations
func (m *NodeProfileMock) HasIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionPreCounter)
}

//HasIntroductionFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) HasIntroductionFinished() bool {
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

type mNodeProfileMockIsAcceptableHost struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockIsAcceptableHostExpectation
	expectationSeries []*NodeProfileMockIsAcceptableHostExpectation
}

type NodeProfileMockIsAcceptableHostExpectation struct {
	input  *NodeProfileMockIsAcceptableHostInput
	result *NodeProfileMockIsAcceptableHostResult
}

type NodeProfileMockIsAcceptableHostInput struct {
	p endpoints.HostIdentityHolder
}

type NodeProfileMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of NodeProfile.IsAcceptableHost is expected from 1 to Infinity times
func (m *mNodeProfileMockIsAcceptableHost) Expect(p endpoints.HostIdentityHolder) *mNodeProfileMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &NodeProfileMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of NodeProfile.IsAcceptableHost
func (m *mNodeProfileMockIsAcceptableHost) Return(r bool) *NodeProfileMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.IsAcceptableHost is expected once
func (m *mNodeProfileMockIsAcceptableHost) ExpectOnce(p endpoints.HostIdentityHolder) *NodeProfileMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockIsAcceptableHostExpectation{}
	expectation.input = &NodeProfileMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &NodeProfileMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of NodeProfile.IsAcceptableHost method
func (m *mNodeProfileMockIsAcceptableHost) Set(f func(p endpoints.HostIdentityHolder) (r bool)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) IsAcceptableHost(p endpoints.HostIdentityHolder) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeProfileMockIsAcceptableHostInput{p}, "NodeProfile.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeProfileMockIsAcceptableHostInput{p}, "NodeProfile.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of NodeProfileMock.IsAcceptableHostFunc invocations
func (m *NodeProfileMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of NodeProfileMock.IsAcceptableHost invocations
func (m *NodeProfileMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) IsAcceptableHostFinished() bool {
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

type mNodeProfileMockIsJoiner struct {
	mock              *NodeProfileMock
	mainExpectation   *NodeProfileMockIsJoinerExpectation
	expectationSeries []*NodeProfileMockIsJoinerExpectation
}

type NodeProfileMockIsJoinerExpectation struct {
	result *NodeProfileMockIsJoinerResult
}

type NodeProfileMockIsJoinerResult struct {
	r bool
}

//Expect specifies that invocation of NodeProfile.IsJoiner is expected from 1 to Infinity times
func (m *mNodeProfileMockIsJoiner) Expect() *mNodeProfileMockIsJoiner {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockIsJoinerExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeProfile.IsJoiner
func (m *mNodeProfileMockIsJoiner) Return(r bool) *NodeProfileMock {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeProfileMockIsJoinerExpectation{}
	}
	m.mainExpectation.result = &NodeProfileMockIsJoinerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeProfile.IsJoiner is expected once
func (m *mNodeProfileMockIsJoiner) ExpectOnce() *NodeProfileMockIsJoinerExpectation {
	m.mock.IsJoinerFunc = nil
	m.mainExpectation = nil

	expectation := &NodeProfileMockIsJoinerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeProfileMockIsJoinerExpectation) Return(r bool) {
	e.result = &NodeProfileMockIsJoinerResult{r}
}

//Set uses given function f as a mock of NodeProfile.IsJoiner method
func (m *mNodeProfileMockIsJoiner) Set(f func() (r bool)) *NodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsJoinerFunc = f
	return m.mock
}

//IsJoiner implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeProfile interface
func (m *NodeProfileMock) IsJoiner() (r bool) {
	counter := atomic.AddUint64(&m.IsJoinerPreCounter, 1)
	defer atomic.AddUint64(&m.IsJoinerCounter, 1)

	if len(m.IsJoinerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsJoinerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeProfileMock.IsJoiner.")
			return
		}

		result := m.IsJoinerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.IsJoiner")
			return
		}

		r = result.r

		return
	}

	if m.IsJoinerMock.mainExpectation != nil {

		result := m.IsJoinerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeProfileMock.IsJoiner")
		}

		r = result.r

		return
	}

	if m.IsJoinerFunc == nil {
		m.t.Fatalf("Unexpected call to NodeProfileMock.IsJoiner.")
		return
	}

	return m.IsJoinerFunc()
}

//IsJoinerMinimockCounter returns a count of NodeProfileMock.IsJoinerFunc invocations
func (m *NodeProfileMock) IsJoinerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerCounter)
}

//IsJoinerMinimockPreCounter returns the value of NodeProfileMock.IsJoiner invocations
func (m *NodeProfileMock) IsJoinerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerPreCounter)
}

//IsJoinerFinished returns true if mock invocations count is ok
func (m *NodeProfileMock) IsJoinerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsJoinerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsJoinerCounter) == uint64(len(m.IsJoinerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsJoinerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsJoinerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsJoinerFunc != nil {
		return atomic.LoadUint64(&m.IsJoinerCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeProfileMock) ValidateCallCounters() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetNodePublicKey")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetNodePublicKeyStore")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetOpMode")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetPrimaryRole")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetShortNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.IsAcceptableHost")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.IsJoiner")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeProfileMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeProfileMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeProfileMock) MinimockFinish() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetNodePublicKey")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetNodePublicKeyStore")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetOpMode")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetPrimaryRole")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetShortNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.IsAcceptableHost")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to NodeProfileMock.IsJoiner")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeProfileMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeProfileMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetAnnouncementSignatureFinished()
		ok = ok && m.GetDeclaredPowerFinished()
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetIndexFinished()
		ok = ok && m.GetIntroductionFinished()
		ok = ok && m.GetNodePublicKeyFinished()
		ok = ok && m.GetNodePublicKeyStoreFinished()
		ok = ok && m.GetOpModeFinished()
		ok = ok && m.GetPrimaryRoleFinished()
		ok = ok && m.GetShortNodeIDFinished()
		ok = ok && m.GetSignatureVerifierFinished()
		ok = ok && m.GetSpecialRolesFinished()
		ok = ok && m.GetStartPowerFinished()
		ok = ok && m.HasIntroductionFinished()
		ok = ok && m.IsAcceptableHostFinished()
		ok = ok && m.IsJoinerFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAnnouncementSignatureFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetAnnouncementSignature")
			}

			if !m.GetDeclaredPowerFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetDeclaredPower")
			}

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetDefaultEndpoint")
			}

			if !m.GetIndexFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetIndex")
			}

			if !m.GetIntroductionFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetIntroduction")
			}

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetNodePublicKey")
			}

			if !m.GetNodePublicKeyStoreFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetNodePublicKeyStore")
			}

			if !m.GetOpModeFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetOpMode")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetPrimaryRole")
			}

			if !m.GetShortNodeIDFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetShortNodeID")
			}

			if !m.GetSignatureVerifierFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetSignatureVerifier")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetSpecialRoles")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to NodeProfileMock.GetStartPower")
			}

			if !m.HasIntroductionFinished() {
				m.t.Error("Expected call to NodeProfileMock.HasIntroduction")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to NodeProfileMock.IsAcceptableHost")
			}

			if !m.IsJoinerFinished() {
				m.t.Error("Expected call to NodeProfileMock.IsJoiner")
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
func (m *NodeProfileMock) AllMocksCalled() bool {

	if !m.GetAnnouncementSignatureFinished() {
		return false
	}

	if !m.GetDeclaredPowerFinished() {
		return false
	}

	if !m.GetDefaultEndpointFinished() {
		return false
	}

	if !m.GetIndexFinished() {
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

	if !m.GetOpModeFinished() {
		return false
	}

	if !m.GetPrimaryRoleFinished() {
		return false
	}

	if !m.GetShortNodeIDFinished() {
		return false
	}

	if !m.GetSignatureVerifierFinished() {
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

	if !m.IsJoinerFinished() {
		return false
	}

	return true
}

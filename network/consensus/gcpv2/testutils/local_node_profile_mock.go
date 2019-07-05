package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LocalNodeProfile" can be found in github.com/insolar/insolar/network/consensus/gcpv2/common
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	common "github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	testify_assert "github.com/stretchr/testify/assert"
)

//LocalNodeProfileMock implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile
type LocalNodeProfileMock struct {
	t minimock.Tester

	GetAnnouncementSignatureFunc       func() (r common.SignatureHolder)
	GetAnnouncementSignatureCounter    uint64
	GetAnnouncementSignaturePreCounter uint64
	GetAnnouncementSignatureMock       mLocalNodeProfileMockGetAnnouncementSignature

	GetDeclaredPowerFunc       func() (r common2.MemberPower)
	GetDeclaredPowerCounter    uint64
	GetDeclaredPowerPreCounter uint64
	GetDeclaredPowerMock       mLocalNodeProfileMockGetDeclaredPower

	GetDefaultEndpointFunc       func() (r common.NodeEndpoint)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mLocalNodeProfileMockGetDefaultEndpoint

	GetIndexFunc       func() (r int)
	GetIndexCounter    uint64
	GetIndexPreCounter uint64
	GetIndexMock       mLocalNodeProfileMockGetIndex

	GetIntroductionFunc       func() (r common2.NodeIntroduction)
	GetIntroductionCounter    uint64
	GetIntroductionPreCounter uint64
	GetIntroductionMock       mLocalNodeProfileMockGetIntroduction

	GetNodePublicKeyFunc       func() (r common.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mLocalNodeProfileMockGetNodePublicKey

	GetNodePublicKeyStoreFunc       func() (r common.PublicKeyStore)
	GetNodePublicKeyStoreCounter    uint64
	GetNodePublicKeyStorePreCounter uint64
	GetNodePublicKeyStoreMock       mLocalNodeProfileMockGetNodePublicKeyStore

	GetPrimaryRoleFunc       func() (r common2.NodePrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mLocalNodeProfileMockGetPrimaryRole

	GetShortNodeIDFunc       func() (r common.ShortNodeID)
	GetShortNodeIDCounter    uint64
	GetShortNodeIDPreCounter uint64
	GetShortNodeIDMock       mLocalNodeProfileMockGetShortNodeID

	GetSignatureVerifierFunc       func() (r common.SignatureVerifier)
	GetSignatureVerifierCounter    uint64
	GetSignatureVerifierPreCounter uint64
	GetSignatureVerifierMock       mLocalNodeProfileMockGetSignatureVerifier

	GetSpecialRolesFunc       func() (r common2.NodeSpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mLocalNodeProfileMockGetSpecialRoles

	GetStartPowerFunc       func() (r common2.MemberPower)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mLocalNodeProfileMockGetStartPower

	GetStateFunc       func() (r common2.MembershipState)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mLocalNodeProfileMockGetState

	HasIntroductionFunc       func() (r bool)
	HasIntroductionCounter    uint64
	HasIntroductionPreCounter uint64
	HasIntroductionMock       mLocalNodeProfileMockHasIntroduction

	IsAcceptableHostFunc       func(p common.HostIdentityHolder) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mLocalNodeProfileMockIsAcceptableHost

	LocalNodeProfileFunc       func()
	LocalNodeProfileCounter    uint64
	LocalNodeProfilePreCounter uint64
	LocalNodeProfileMock       mLocalNodeProfileMockLocalNodeProfile
}

//NewLocalNodeProfileMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile
func NewLocalNodeProfileMock(t minimock.Tester) *LocalNodeProfileMock {
	m := &LocalNodeProfileMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAnnouncementSignatureMock = mLocalNodeProfileMockGetAnnouncementSignature{mock: m}
	m.GetDeclaredPowerMock = mLocalNodeProfileMockGetDeclaredPower{mock: m}
	m.GetDefaultEndpointMock = mLocalNodeProfileMockGetDefaultEndpoint{mock: m}
	m.GetIndexMock = mLocalNodeProfileMockGetIndex{mock: m}
	m.GetIntroductionMock = mLocalNodeProfileMockGetIntroduction{mock: m}
	m.GetNodePublicKeyMock = mLocalNodeProfileMockGetNodePublicKey{mock: m}
	m.GetNodePublicKeyStoreMock = mLocalNodeProfileMockGetNodePublicKeyStore{mock: m}
	m.GetPrimaryRoleMock = mLocalNodeProfileMockGetPrimaryRole{mock: m}
	m.GetShortNodeIDMock = mLocalNodeProfileMockGetShortNodeID{mock: m}
	m.GetSignatureVerifierMock = mLocalNodeProfileMockGetSignatureVerifier{mock: m}
	m.GetSpecialRolesMock = mLocalNodeProfileMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mLocalNodeProfileMockGetStartPower{mock: m}
	m.GetStateMock = mLocalNodeProfileMockGetState{mock: m}
	m.HasIntroductionMock = mLocalNodeProfileMockHasIntroduction{mock: m}
	m.IsAcceptableHostMock = mLocalNodeProfileMockIsAcceptableHost{mock: m}
	m.LocalNodeProfileMock = mLocalNodeProfileMockLocalNodeProfile{mock: m}

	return m
}

type mLocalNodeProfileMockGetAnnouncementSignature struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetAnnouncementSignatureExpectation
	expectationSeries []*LocalNodeProfileMockGetAnnouncementSignatureExpectation
}

type LocalNodeProfileMockGetAnnouncementSignatureExpectation struct {
	result *LocalNodeProfileMockGetAnnouncementSignatureResult
}

type LocalNodeProfileMockGetAnnouncementSignatureResult struct {
	r common.SignatureHolder
}

//Expect specifies that invocation of LocalNodeProfile.GetAnnouncementSignature is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetAnnouncementSignature) Expect() *mLocalNodeProfileMockGetAnnouncementSignature {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetAnnouncementSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetAnnouncementSignature
func (m *mLocalNodeProfileMockGetAnnouncementSignature) Return(r common.SignatureHolder) *LocalNodeProfileMock {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetAnnouncementSignatureExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetAnnouncementSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetAnnouncementSignature is expected once
func (m *mLocalNodeProfileMockGetAnnouncementSignature) ExpectOnce() *LocalNodeProfileMockGetAnnouncementSignatureExpectation {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetAnnouncementSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetAnnouncementSignatureExpectation) Return(r common.SignatureHolder) {
	e.result = &LocalNodeProfileMockGetAnnouncementSignatureResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetAnnouncementSignature method
func (m *mLocalNodeProfileMockGetAnnouncementSignature) Set(f func() (r common.SignatureHolder)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAnnouncementSignatureFunc = f
	return m.mock
}

//GetAnnouncementSignature implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetAnnouncementSignature() (r common.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetAnnouncementSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetAnnouncementSignatureCounter, 1)

	if len(m.GetAnnouncementSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAnnouncementSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetAnnouncementSignature.")
			return
		}

		result := m.GetAnnouncementSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetAnnouncementSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureMock.mainExpectation != nil {

		result := m.GetAnnouncementSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetAnnouncementSignature")
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetAnnouncementSignature.")
		return
	}

	return m.GetAnnouncementSignatureFunc()
}

//GetAnnouncementSignatureMinimockCounter returns a count of LocalNodeProfileMock.GetAnnouncementSignatureFunc invocations
func (m *LocalNodeProfileMock) GetAnnouncementSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter)
}

//GetAnnouncementSignatureMinimockPreCounter returns the value of LocalNodeProfileMock.GetAnnouncementSignature invocations
func (m *LocalNodeProfileMock) GetAnnouncementSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignaturePreCounter)
}

//GetAnnouncementSignatureFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetAnnouncementSignatureFinished() bool {
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

type mLocalNodeProfileMockGetDeclaredPower struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetDeclaredPowerExpectation
	expectationSeries []*LocalNodeProfileMockGetDeclaredPowerExpectation
}

type LocalNodeProfileMockGetDeclaredPowerExpectation struct {
	result *LocalNodeProfileMockGetDeclaredPowerResult
}

type LocalNodeProfileMockGetDeclaredPowerResult struct {
	r common2.MemberPower
}

//Expect specifies that invocation of LocalNodeProfile.GetDeclaredPower is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetDeclaredPower) Expect() *mLocalNodeProfileMockGetDeclaredPower {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetDeclaredPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetDeclaredPower
func (m *mLocalNodeProfileMockGetDeclaredPower) Return(r common2.MemberPower) *LocalNodeProfileMock {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetDeclaredPowerExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetDeclaredPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetDeclaredPower is expected once
func (m *mLocalNodeProfileMockGetDeclaredPower) ExpectOnce() *LocalNodeProfileMockGetDeclaredPowerExpectation {
	m.mock.GetDeclaredPowerFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetDeclaredPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetDeclaredPowerExpectation) Return(r common2.MemberPower) {
	e.result = &LocalNodeProfileMockGetDeclaredPowerResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetDeclaredPower method
func (m *mLocalNodeProfileMockGetDeclaredPower) Set(f func() (r common2.MemberPower)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDeclaredPowerFunc = f
	return m.mock
}

//GetDeclaredPower implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetDeclaredPower() (r common2.MemberPower) {
	counter := atomic.AddUint64(&m.GetDeclaredPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetDeclaredPowerCounter, 1)

	if len(m.GetDeclaredPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDeclaredPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetDeclaredPower.")
			return
		}

		result := m.GetDeclaredPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetDeclaredPower")
			return
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerMock.mainExpectation != nil {

		result := m.GetDeclaredPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetDeclaredPower")
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetDeclaredPower.")
		return
	}

	return m.GetDeclaredPowerFunc()
}

//GetDeclaredPowerMinimockCounter returns a count of LocalNodeProfileMock.GetDeclaredPowerFunc invocations
func (m *LocalNodeProfileMock) GetDeclaredPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerCounter)
}

//GetDeclaredPowerMinimockPreCounter returns the value of LocalNodeProfileMock.GetDeclaredPower invocations
func (m *LocalNodeProfileMock) GetDeclaredPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerPreCounter)
}

//GetDeclaredPowerFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetDeclaredPowerFinished() bool {
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

type mLocalNodeProfileMockGetDefaultEndpoint struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetDefaultEndpointExpectation
	expectationSeries []*LocalNodeProfileMockGetDefaultEndpointExpectation
}

type LocalNodeProfileMockGetDefaultEndpointExpectation struct {
	result *LocalNodeProfileMockGetDefaultEndpointResult
}

type LocalNodeProfileMockGetDefaultEndpointResult struct {
	r common.NodeEndpoint
}

//Expect specifies that invocation of LocalNodeProfile.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetDefaultEndpoint) Expect() *mLocalNodeProfileMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetDefaultEndpoint
func (m *mLocalNodeProfileMockGetDefaultEndpoint) Return(r common.NodeEndpoint) *LocalNodeProfileMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetDefaultEndpoint is expected once
func (m *mLocalNodeProfileMockGetDefaultEndpoint) ExpectOnce() *LocalNodeProfileMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetDefaultEndpointExpectation) Return(r common.NodeEndpoint) {
	e.result = &LocalNodeProfileMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetDefaultEndpoint method
func (m *mLocalNodeProfileMockGetDefaultEndpoint) Set(f func() (r common.NodeEndpoint)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetDefaultEndpoint() (r common.NodeEndpoint) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of LocalNodeProfileMock.GetDefaultEndpointFunc invocations
func (m *LocalNodeProfileMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of LocalNodeProfileMock.GetDefaultEndpoint invocations
func (m *LocalNodeProfileMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetDefaultEndpointFinished() bool {
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

type mLocalNodeProfileMockGetIndex struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetIndexExpectation
	expectationSeries []*LocalNodeProfileMockGetIndexExpectation
}

type LocalNodeProfileMockGetIndexExpectation struct {
	result *LocalNodeProfileMockGetIndexResult
}

type LocalNodeProfileMockGetIndexResult struct {
	r int
}

//Expect specifies that invocation of LocalNodeProfile.GetIndex is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetIndex) Expect() *mLocalNodeProfileMockGetIndex {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetIndexExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetIndex
func (m *mLocalNodeProfileMockGetIndex) Return(r int) *LocalNodeProfileMock {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetIndexExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetIndex is expected once
func (m *mLocalNodeProfileMockGetIndex) ExpectOnce() *LocalNodeProfileMockGetIndexExpectation {
	m.mock.GetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetIndexExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetIndexExpectation) Return(r int) {
	e.result = &LocalNodeProfileMockGetIndexResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetIndex method
func (m *mLocalNodeProfileMockGetIndex) Set(f func() (r int)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIndexFunc = f
	return m.mock
}

//GetIndex implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetIndex() (r int) {
	counter := atomic.AddUint64(&m.GetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.GetIndexCounter, 1)

	if len(m.GetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetIndex.")
			return
		}

		result := m.GetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetIndex")
			return
		}

		r = result.r

		return
	}

	if m.GetIndexMock.mainExpectation != nil {

		result := m.GetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetIndex")
		}

		r = result.r

		return
	}

	if m.GetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetIndex.")
		return
	}

	return m.GetIndexFunc()
}

//GetIndexMinimockCounter returns a count of LocalNodeProfileMock.GetIndexFunc invocations
func (m *LocalNodeProfileMock) GetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexCounter)
}

//GetIndexMinimockPreCounter returns the value of LocalNodeProfileMock.GetIndex invocations
func (m *LocalNodeProfileMock) GetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexPreCounter)
}

//GetIndexFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetIndexFinished() bool {
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

type mLocalNodeProfileMockGetIntroduction struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetIntroductionExpectation
	expectationSeries []*LocalNodeProfileMockGetIntroductionExpectation
}

type LocalNodeProfileMockGetIntroductionExpectation struct {
	result *LocalNodeProfileMockGetIntroductionResult
}

type LocalNodeProfileMockGetIntroductionResult struct {
	r common2.NodeIntroduction
}

//Expect specifies that invocation of LocalNodeProfile.GetIntroduction is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetIntroduction) Expect() *mLocalNodeProfileMockGetIntroduction {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetIntroduction
func (m *mLocalNodeProfileMockGetIntroduction) Return(r common2.NodeIntroduction) *LocalNodeProfileMock {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetIntroductionExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetIntroduction is expected once
func (m *mLocalNodeProfileMockGetIntroduction) ExpectOnce() *LocalNodeProfileMockGetIntroductionExpectation {
	m.mock.GetIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetIntroductionExpectation) Return(r common2.NodeIntroduction) {
	e.result = &LocalNodeProfileMockGetIntroductionResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetIntroduction method
func (m *mLocalNodeProfileMockGetIntroduction) Set(f func() (r common2.NodeIntroduction)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIntroductionFunc = f
	return m.mock
}

//GetIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetIntroduction() (r common2.NodeIntroduction) {
	counter := atomic.AddUint64(&m.GetIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.GetIntroductionCounter, 1)

	if len(m.GetIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetIntroduction.")
			return
		}

		result := m.GetIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.GetIntroductionMock.mainExpectation != nil {

		result := m.GetIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetIntroduction")
		}

		r = result.r

		return
	}

	if m.GetIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetIntroduction.")
		return
	}

	return m.GetIntroductionFunc()
}

//GetIntroductionMinimockCounter returns a count of LocalNodeProfileMock.GetIntroductionFunc invocations
func (m *LocalNodeProfileMock) GetIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionCounter)
}

//GetIntroductionMinimockPreCounter returns the value of LocalNodeProfileMock.GetIntroduction invocations
func (m *LocalNodeProfileMock) GetIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionPreCounter)
}

//GetIntroductionFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetIntroductionFinished() bool {
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

type mLocalNodeProfileMockGetNodePublicKey struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetNodePublicKeyExpectation
	expectationSeries []*LocalNodeProfileMockGetNodePublicKeyExpectation
}

type LocalNodeProfileMockGetNodePublicKeyExpectation struct {
	result *LocalNodeProfileMockGetNodePublicKeyResult
}

type LocalNodeProfileMockGetNodePublicKeyResult struct {
	r common.SignatureKeyHolder
}

//Expect specifies that invocation of LocalNodeProfile.GetNodePublicKey is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetNodePublicKey) Expect() *mLocalNodeProfileMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetNodePublicKey
func (m *mLocalNodeProfileMockGetNodePublicKey) Return(r common.SignatureKeyHolder) *LocalNodeProfileMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetNodePublicKey is expected once
func (m *mLocalNodeProfileMockGetNodePublicKey) ExpectOnce() *LocalNodeProfileMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetNodePublicKeyExpectation) Return(r common.SignatureKeyHolder) {
	e.result = &LocalNodeProfileMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetNodePublicKey method
func (m *mLocalNodeProfileMockGetNodePublicKey) Set(f func() (r common.SignatureKeyHolder)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetNodePublicKey() (r common.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of LocalNodeProfileMock.GetNodePublicKeyFunc invocations
func (m *LocalNodeProfileMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of LocalNodeProfileMock.GetNodePublicKey invocations
func (m *LocalNodeProfileMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetNodePublicKeyFinished() bool {
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

type mLocalNodeProfileMockGetNodePublicKeyStore struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetNodePublicKeyStoreExpectation
	expectationSeries []*LocalNodeProfileMockGetNodePublicKeyStoreExpectation
}

type LocalNodeProfileMockGetNodePublicKeyStoreExpectation struct {
	result *LocalNodeProfileMockGetNodePublicKeyStoreResult
}

type LocalNodeProfileMockGetNodePublicKeyStoreResult struct {
	r common.PublicKeyStore
}

//Expect specifies that invocation of LocalNodeProfile.GetNodePublicKeyStore is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetNodePublicKeyStore) Expect() *mLocalNodeProfileMockGetNodePublicKeyStore {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetNodePublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetNodePublicKeyStore
func (m *mLocalNodeProfileMockGetNodePublicKeyStore) Return(r common.PublicKeyStore) *LocalNodeProfileMock {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetNodePublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetNodePublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetNodePublicKeyStore is expected once
func (m *mLocalNodeProfileMockGetNodePublicKeyStore) ExpectOnce() *LocalNodeProfileMockGetNodePublicKeyStoreExpectation {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetNodePublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetNodePublicKeyStoreExpectation) Return(r common.PublicKeyStore) {
	e.result = &LocalNodeProfileMockGetNodePublicKeyStoreResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetNodePublicKeyStore method
func (m *mLocalNodeProfileMockGetNodePublicKeyStore) Set(f func() (r common.PublicKeyStore)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyStoreFunc = f
	return m.mock
}

//GetNodePublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetNodePublicKeyStore() (r common.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyStoreCounter, 1)

	if len(m.GetNodePublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetNodePublicKeyStore.")
			return
		}

		result := m.GetNodePublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetNodePublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreMock.mainExpectation != nil {

		result := m.GetNodePublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetNodePublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetNodePublicKeyStore.")
		return
	}

	return m.GetNodePublicKeyStoreFunc()
}

//GetNodePublicKeyStoreMinimockCounter returns a count of LocalNodeProfileMock.GetNodePublicKeyStoreFunc invocations
func (m *LocalNodeProfileMock) GetNodePublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter)
}

//GetNodePublicKeyStoreMinimockPreCounter returns the value of LocalNodeProfileMock.GetNodePublicKeyStore invocations
func (m *LocalNodeProfileMock) GetNodePublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStorePreCounter)
}

//GetNodePublicKeyStoreFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetNodePublicKeyStoreFinished() bool {
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

type mLocalNodeProfileMockGetPrimaryRole struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetPrimaryRoleExpectation
	expectationSeries []*LocalNodeProfileMockGetPrimaryRoleExpectation
}

type LocalNodeProfileMockGetPrimaryRoleExpectation struct {
	result *LocalNodeProfileMockGetPrimaryRoleResult
}

type LocalNodeProfileMockGetPrimaryRoleResult struct {
	r common2.NodePrimaryRole
}

//Expect specifies that invocation of LocalNodeProfile.GetPrimaryRole is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetPrimaryRole) Expect() *mLocalNodeProfileMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetPrimaryRole
func (m *mLocalNodeProfileMockGetPrimaryRole) Return(r common2.NodePrimaryRole) *LocalNodeProfileMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetPrimaryRole is expected once
func (m *mLocalNodeProfileMockGetPrimaryRole) ExpectOnce() *LocalNodeProfileMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetPrimaryRoleExpectation) Return(r common2.NodePrimaryRole) {
	e.result = &LocalNodeProfileMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetPrimaryRole method
func (m *mLocalNodeProfileMockGetPrimaryRole) Set(f func() (r common2.NodePrimaryRole)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetPrimaryRole() (r common2.NodePrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of LocalNodeProfileMock.GetPrimaryRoleFunc invocations
func (m *LocalNodeProfileMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of LocalNodeProfileMock.GetPrimaryRole invocations
func (m *LocalNodeProfileMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetPrimaryRoleFinished() bool {
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

type mLocalNodeProfileMockGetShortNodeID struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetShortNodeIDExpectation
	expectationSeries []*LocalNodeProfileMockGetShortNodeIDExpectation
}

type LocalNodeProfileMockGetShortNodeIDExpectation struct {
	result *LocalNodeProfileMockGetShortNodeIDResult
}

type LocalNodeProfileMockGetShortNodeIDResult struct {
	r common.ShortNodeID
}

//Expect specifies that invocation of LocalNodeProfile.GetShortNodeID is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetShortNodeID) Expect() *mLocalNodeProfileMockGetShortNodeID {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetShortNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetShortNodeID
func (m *mLocalNodeProfileMockGetShortNodeID) Return(r common.ShortNodeID) *LocalNodeProfileMock {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetShortNodeIDExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetShortNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetShortNodeID is expected once
func (m *mLocalNodeProfileMockGetShortNodeID) ExpectOnce() *LocalNodeProfileMockGetShortNodeIDExpectation {
	m.mock.GetShortNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetShortNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetShortNodeIDExpectation) Return(r common.ShortNodeID) {
	e.result = &LocalNodeProfileMockGetShortNodeIDResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetShortNodeID method
func (m *mLocalNodeProfileMockGetShortNodeID) Set(f func() (r common.ShortNodeID)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetShortNodeIDFunc = f
	return m.mock
}

//GetShortNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetShortNodeID() (r common.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetShortNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetShortNodeIDCounter, 1)

	if len(m.GetShortNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetShortNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetShortNodeID.")
			return
		}

		result := m.GetShortNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetShortNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDMock.mainExpectation != nil {

		result := m.GetShortNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetShortNodeID")
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetShortNodeID.")
		return
	}

	return m.GetShortNodeIDFunc()
}

//GetShortNodeIDMinimockCounter returns a count of LocalNodeProfileMock.GetShortNodeIDFunc invocations
func (m *LocalNodeProfileMock) GetShortNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDCounter)
}

//GetShortNodeIDMinimockPreCounter returns the value of LocalNodeProfileMock.GetShortNodeID invocations
func (m *LocalNodeProfileMock) GetShortNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDPreCounter)
}

//GetShortNodeIDFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetShortNodeIDFinished() bool {
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

type mLocalNodeProfileMockGetSignatureVerifier struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetSignatureVerifierExpectation
	expectationSeries []*LocalNodeProfileMockGetSignatureVerifierExpectation
}

type LocalNodeProfileMockGetSignatureVerifierExpectation struct {
	result *LocalNodeProfileMockGetSignatureVerifierResult
}

type LocalNodeProfileMockGetSignatureVerifierResult struct {
	r common.SignatureVerifier
}

//Expect specifies that invocation of LocalNodeProfile.GetSignatureVerifier is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetSignatureVerifier) Expect() *mLocalNodeProfileMockGetSignatureVerifier {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetSignatureVerifierExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetSignatureVerifier
func (m *mLocalNodeProfileMockGetSignatureVerifier) Return(r common.SignatureVerifier) *LocalNodeProfileMock {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetSignatureVerifierExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetSignatureVerifierResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetSignatureVerifier is expected once
func (m *mLocalNodeProfileMockGetSignatureVerifier) ExpectOnce() *LocalNodeProfileMockGetSignatureVerifierExpectation {
	m.mock.GetSignatureVerifierFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetSignatureVerifierExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetSignatureVerifierExpectation) Return(r common.SignatureVerifier) {
	e.result = &LocalNodeProfileMockGetSignatureVerifierResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetSignatureVerifier method
func (m *mLocalNodeProfileMockGetSignatureVerifier) Set(f func() (r common.SignatureVerifier)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureVerifierFunc = f
	return m.mock
}

//GetSignatureVerifier implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetSignatureVerifier() (r common.SignatureVerifier) {
	counter := atomic.AddUint64(&m.GetSignatureVerifierPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureVerifierCounter, 1)

	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureVerifierMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetSignatureVerifier.")
			return
		}

		result := m.GetSignatureVerifierMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetSignatureVerifier")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierMock.mainExpectation != nil {

		result := m.GetSignatureVerifierMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetSignatureVerifier")
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetSignatureVerifier.")
		return
	}

	return m.GetSignatureVerifierFunc()
}

//GetSignatureVerifierMinimockCounter returns a count of LocalNodeProfileMock.GetSignatureVerifierFunc invocations
func (m *LocalNodeProfileMock) GetSignatureVerifierMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierCounter)
}

//GetSignatureVerifierMinimockPreCounter returns the value of LocalNodeProfileMock.GetSignatureVerifier invocations
func (m *LocalNodeProfileMock) GetSignatureVerifierMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierPreCounter)
}

//GetSignatureVerifierFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetSignatureVerifierFinished() bool {
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

type mLocalNodeProfileMockGetSpecialRoles struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetSpecialRolesExpectation
	expectationSeries []*LocalNodeProfileMockGetSpecialRolesExpectation
}

type LocalNodeProfileMockGetSpecialRolesExpectation struct {
	result *LocalNodeProfileMockGetSpecialRolesResult
}

type LocalNodeProfileMockGetSpecialRolesResult struct {
	r common2.NodeSpecialRole
}

//Expect specifies that invocation of LocalNodeProfile.GetSpecialRoles is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetSpecialRoles) Expect() *mLocalNodeProfileMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetSpecialRoles
func (m *mLocalNodeProfileMockGetSpecialRoles) Return(r common2.NodeSpecialRole) *LocalNodeProfileMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetSpecialRoles is expected once
func (m *mLocalNodeProfileMockGetSpecialRoles) ExpectOnce() *LocalNodeProfileMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetSpecialRolesExpectation) Return(r common2.NodeSpecialRole) {
	e.result = &LocalNodeProfileMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetSpecialRoles method
func (m *mLocalNodeProfileMockGetSpecialRoles) Set(f func() (r common2.NodeSpecialRole)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetSpecialRoles() (r common2.NodeSpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of LocalNodeProfileMock.GetSpecialRolesFunc invocations
func (m *LocalNodeProfileMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of LocalNodeProfileMock.GetSpecialRoles invocations
func (m *LocalNodeProfileMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetSpecialRolesFinished() bool {
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

type mLocalNodeProfileMockGetStartPower struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetStartPowerExpectation
	expectationSeries []*LocalNodeProfileMockGetStartPowerExpectation
}

type LocalNodeProfileMockGetStartPowerExpectation struct {
	result *LocalNodeProfileMockGetStartPowerResult
}

type LocalNodeProfileMockGetStartPowerResult struct {
	r common2.MemberPower
}

//Expect specifies that invocation of LocalNodeProfile.GetStartPower is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetStartPower) Expect() *mLocalNodeProfileMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetStartPower
func (m *mLocalNodeProfileMockGetStartPower) Return(r common2.MemberPower) *LocalNodeProfileMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetStartPower is expected once
func (m *mLocalNodeProfileMockGetStartPower) ExpectOnce() *LocalNodeProfileMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetStartPowerExpectation) Return(r common2.MemberPower) {
	e.result = &LocalNodeProfileMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetStartPower method
func (m *mLocalNodeProfileMockGetStartPower) Set(f func() (r common2.MemberPower)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetStartPower() (r common2.MemberPower) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of LocalNodeProfileMock.GetStartPowerFunc invocations
func (m *LocalNodeProfileMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of LocalNodeProfileMock.GetStartPower invocations
func (m *LocalNodeProfileMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetStartPowerFinished() bool {
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

type mLocalNodeProfileMockGetState struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockGetStateExpectation
	expectationSeries []*LocalNodeProfileMockGetStateExpectation
}

type LocalNodeProfileMockGetStateExpectation struct {
	result *LocalNodeProfileMockGetStateResult
}

type LocalNodeProfileMockGetStateResult struct {
	r common2.MembershipState
}

//Expect specifies that invocation of LocalNodeProfile.GetState is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockGetState) Expect() *mLocalNodeProfileMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.GetState
func (m *mLocalNodeProfileMockGetState) Return(r common2.MembershipState) *LocalNodeProfileMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockGetStateExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.GetState is expected once
func (m *mLocalNodeProfileMockGetState) ExpectOnce() *LocalNodeProfileMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockGetStateExpectation) Return(r common2.MembershipState) {
	e.result = &LocalNodeProfileMockGetStateResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.GetState method
func (m *mLocalNodeProfileMockGetState) Set(f func() (r common2.MembershipState)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) GetState() (r common2.MembershipState) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of LocalNodeProfileMock.GetStateFunc invocations
func (m *LocalNodeProfileMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of LocalNodeProfileMock.GetState invocations
func (m *LocalNodeProfileMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) GetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStateCounter) == uint64(len(m.GetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStateFunc != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	return true
}

type mLocalNodeProfileMockHasIntroduction struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockHasIntroductionExpectation
	expectationSeries []*LocalNodeProfileMockHasIntroductionExpectation
}

type LocalNodeProfileMockHasIntroductionExpectation struct {
	result *LocalNodeProfileMockHasIntroductionResult
}

type LocalNodeProfileMockHasIntroductionResult struct {
	r bool
}

//Expect specifies that invocation of LocalNodeProfile.HasIntroduction is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockHasIntroduction) Expect() *mLocalNodeProfileMockHasIntroduction {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockHasIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.HasIntroduction
func (m *mLocalNodeProfileMockHasIntroduction) Return(r bool) *LocalNodeProfileMock {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockHasIntroductionExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockHasIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.HasIntroduction is expected once
func (m *mLocalNodeProfileMockHasIntroduction) ExpectOnce() *LocalNodeProfileMockHasIntroductionExpectation {
	m.mock.HasIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockHasIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockHasIntroductionExpectation) Return(r bool) {
	e.result = &LocalNodeProfileMockHasIntroductionResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.HasIntroduction method
func (m *mLocalNodeProfileMockHasIntroduction) Set(f func() (r bool)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasIntroductionFunc = f
	return m.mock
}

//HasIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) HasIntroduction() (r bool) {
	counter := atomic.AddUint64(&m.HasIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.HasIntroductionCounter, 1)

	if len(m.HasIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.HasIntroduction.")
			return
		}

		result := m.HasIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.HasIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.HasIntroductionMock.mainExpectation != nil {

		result := m.HasIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.HasIntroduction")
		}

		r = result.r

		return
	}

	if m.HasIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.HasIntroduction.")
		return
	}

	return m.HasIntroductionFunc()
}

//HasIntroductionMinimockCounter returns a count of LocalNodeProfileMock.HasIntroductionFunc invocations
func (m *LocalNodeProfileMock) HasIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionCounter)
}

//HasIntroductionMinimockPreCounter returns the value of LocalNodeProfileMock.HasIntroduction invocations
func (m *LocalNodeProfileMock) HasIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionPreCounter)
}

//HasIntroductionFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) HasIntroductionFinished() bool {
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

type mLocalNodeProfileMockIsAcceptableHost struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockIsAcceptableHostExpectation
	expectationSeries []*LocalNodeProfileMockIsAcceptableHostExpectation
}

type LocalNodeProfileMockIsAcceptableHostExpectation struct {
	input  *LocalNodeProfileMockIsAcceptableHostInput
	result *LocalNodeProfileMockIsAcceptableHostResult
}

type LocalNodeProfileMockIsAcceptableHostInput struct {
	p common.HostIdentityHolder
}

type LocalNodeProfileMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of LocalNodeProfile.IsAcceptableHost is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockIsAcceptableHost) Expect(p common.HostIdentityHolder) *mLocalNodeProfileMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &LocalNodeProfileMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of LocalNodeProfile.IsAcceptableHost
func (m *mLocalNodeProfileMockIsAcceptableHost) Return(r bool) *LocalNodeProfileMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &LocalNodeProfileMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.IsAcceptableHost is expected once
func (m *mLocalNodeProfileMockIsAcceptableHost) ExpectOnce(p common.HostIdentityHolder) *LocalNodeProfileMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockIsAcceptableHostExpectation{}
	expectation.input = &LocalNodeProfileMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeProfileMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &LocalNodeProfileMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of LocalNodeProfile.IsAcceptableHost method
func (m *mLocalNodeProfileMockIsAcceptableHost) Set(f func(p common.HostIdentityHolder) (r bool)) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) IsAcceptableHost(p common.HostIdentityHolder) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LocalNodeProfileMockIsAcceptableHostInput{p}, "LocalNodeProfile.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LocalNodeProfileMockIsAcceptableHostInput{p}, "LocalNodeProfile.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeProfileMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of LocalNodeProfileMock.IsAcceptableHostFunc invocations
func (m *LocalNodeProfileMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of LocalNodeProfileMock.IsAcceptableHost invocations
func (m *LocalNodeProfileMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) IsAcceptableHostFinished() bool {
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

type mLocalNodeProfileMockLocalNodeProfile struct {
	mock              *LocalNodeProfileMock
	mainExpectation   *LocalNodeProfileMockLocalNodeProfileExpectation
	expectationSeries []*LocalNodeProfileMockLocalNodeProfileExpectation
}

type LocalNodeProfileMockLocalNodeProfileExpectation struct {
}

//Expect specifies that invocation of LocalNodeProfile.LocalNodeProfile is expected from 1 to Infinity times
func (m *mLocalNodeProfileMockLocalNodeProfile) Expect() *mLocalNodeProfileMockLocalNodeProfile {
	m.mock.LocalNodeProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockLocalNodeProfileExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNodeProfile.LocalNodeProfile
func (m *mLocalNodeProfileMockLocalNodeProfile) Return() *LocalNodeProfileMock {
	m.mock.LocalNodeProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeProfileMockLocalNodeProfileExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of LocalNodeProfile.LocalNodeProfile is expected once
func (m *mLocalNodeProfileMockLocalNodeProfile) ExpectOnce() *LocalNodeProfileMockLocalNodeProfileExpectation {
	m.mock.LocalNodeProfileFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeProfileMockLocalNodeProfileExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of LocalNodeProfile.LocalNodeProfile method
func (m *mLocalNodeProfileMockLocalNodeProfile) Set(f func()) *LocalNodeProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LocalNodeProfileFunc = f
	return m.mock
}

//LocalNodeProfile implements github.com/insolar/insolar/network/consensus/gcpv2/common.LocalNodeProfile interface
func (m *LocalNodeProfileMock) LocalNodeProfile() {
	counter := atomic.AddUint64(&m.LocalNodeProfilePreCounter, 1)
	defer atomic.AddUint64(&m.LocalNodeProfileCounter, 1)

	if len(m.LocalNodeProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LocalNodeProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeProfileMock.LocalNodeProfile.")
			return
		}

		return
	}

	if m.LocalNodeProfileMock.mainExpectation != nil {

		return
	}

	if m.LocalNodeProfileFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeProfileMock.LocalNodeProfile.")
		return
	}

	m.LocalNodeProfileFunc()
}

//LocalNodeProfileMinimockCounter returns a count of LocalNodeProfileMock.LocalNodeProfileFunc invocations
func (m *LocalNodeProfileMock) LocalNodeProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LocalNodeProfileCounter)
}

//LocalNodeProfileMinimockPreCounter returns the value of LocalNodeProfileMock.LocalNodeProfile invocations
func (m *LocalNodeProfileMock) LocalNodeProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LocalNodeProfilePreCounter)
}

//LocalNodeProfileFinished returns true if mock invocations count is ok
func (m *LocalNodeProfileMock) LocalNodeProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LocalNodeProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LocalNodeProfileCounter) == uint64(len(m.LocalNodeProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LocalNodeProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LocalNodeProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LocalNodeProfileFunc != nil {
		return atomic.LoadUint64(&m.LocalNodeProfileCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LocalNodeProfileMock) ValidateCallCounters() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetNodePublicKey")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetNodePublicKeyStore")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetPrimaryRole")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetShortNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetStartPower")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetState")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.IsAcceptableHost")
	}

	if !m.LocalNodeProfileFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.LocalNodeProfile")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LocalNodeProfileMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LocalNodeProfileMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LocalNodeProfileMock) MinimockFinish() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetNodePublicKey")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetNodePublicKeyStore")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetPrimaryRole")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetShortNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetStartPower")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.GetState")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.IsAcceptableHost")
	}

	if !m.LocalNodeProfileFinished() {
		m.t.Fatal("Expected call to LocalNodeProfileMock.LocalNodeProfile")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LocalNodeProfileMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LocalNodeProfileMock) MinimockWait(timeout time.Duration) {
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
		ok = ok && m.GetPrimaryRoleFinished()
		ok = ok && m.GetShortNodeIDFinished()
		ok = ok && m.GetSignatureVerifierFinished()
		ok = ok && m.GetSpecialRolesFinished()
		ok = ok && m.GetStartPowerFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.HasIntroductionFinished()
		ok = ok && m.IsAcceptableHostFinished()
		ok = ok && m.LocalNodeProfileFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAnnouncementSignatureFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetAnnouncementSignature")
			}

			if !m.GetDeclaredPowerFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetDeclaredPower")
			}

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetDefaultEndpoint")
			}

			if !m.GetIndexFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetIndex")
			}

			if !m.GetIntroductionFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetIntroduction")
			}

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetNodePublicKey")
			}

			if !m.GetNodePublicKeyStoreFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetNodePublicKeyStore")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetPrimaryRole")
			}

			if !m.GetShortNodeIDFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetShortNodeID")
			}

			if !m.GetSignatureVerifierFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetSignatureVerifier")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetSpecialRoles")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetStartPower")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.GetState")
			}

			if !m.HasIntroductionFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.HasIntroduction")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.IsAcceptableHost")
			}

			if !m.LocalNodeProfileFinished() {
				m.t.Error("Expected call to LocalNodeProfileMock.LocalNodeProfile")
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
func (m *LocalNodeProfileMock) AllMocksCalled() bool {

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

	if !m.GetStateFinished() {
		return false
	}

	if !m.HasIntroductionFinished() {
		return false
	}

	if !m.IsAcceptableHostFinished() {
		return false
	}

	if !m.LocalNodeProfileFinished() {
		return false
	}

	return true
}

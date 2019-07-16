package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ActiveNode" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
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

//ActiveNodeMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode
type ActiveNodeMock struct {
	t minimock.Tester

	GetAnnouncementSignatureFunc       func() (r cryptkit.SignatureHolder)
	GetAnnouncementSignatureCounter    uint64
	GetAnnouncementSignaturePreCounter uint64
	GetAnnouncementSignatureMock       mActiveNodeMockGetAnnouncementSignature

	GetDeclaredPowerFunc       func() (r member.Power)
	GetDeclaredPowerCounter    uint64
	GetDeclaredPowerPreCounter uint64
	GetDeclaredPowerMock       mActiveNodeMockGetDeclaredPower

	GetDefaultEndpointFunc       func() (r endpoints.Outbound)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mActiveNodeMockGetDefaultEndpoint

	GetIndexFunc       func() (r member.Index)
	GetIndexCounter    uint64
	GetIndexPreCounter uint64
	GetIndexMock       mActiveNodeMockGetIndex

	GetIntroductionFunc       func() (r NodeIntroduction)
	GetIntroductionCounter    uint64
	GetIntroductionPreCounter uint64
	GetIntroductionMock       mActiveNodeMockGetIntroduction

	GetNodePublicKeyFunc       func() (r cryptkit.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mActiveNodeMockGetNodePublicKey

	GetOpModeFunc       func() (r member.OpMode)
	GetOpModeCounter    uint64
	GetOpModePreCounter uint64
	GetOpModeMock       mActiveNodeMockGetOpMode

	GetPrimaryRoleFunc       func() (r member.PrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mActiveNodeMockGetPrimaryRole

	GetPublicKeyStoreFunc       func() (r cryptkit.PublicKeyStore)
	GetPublicKeyStoreCounter    uint64
	GetPublicKeyStorePreCounter uint64
	GetPublicKeyStoreMock       mActiveNodeMockGetPublicKeyStore

	GetShortNodeIDFunc       func() (r insolar.ShortNodeID)
	GetShortNodeIDCounter    uint64
	GetShortNodeIDPreCounter uint64
	GetShortNodeIDMock       mActiveNodeMockGetShortNodeID

	GetSignatureVerifierFunc       func() (r cryptkit.SignatureVerifier)
	GetSignatureVerifierCounter    uint64
	GetSignatureVerifierPreCounter uint64
	GetSignatureVerifierMock       mActiveNodeMockGetSignatureVerifier

	GetSpecialRolesFunc       func() (r member.SpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mActiveNodeMockGetSpecialRoles

	GetStartPowerFunc       func() (r member.Power)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mActiveNodeMockGetStartPower

	HasIntroductionFunc       func() (r bool)
	HasIntroductionCounter    uint64
	HasIntroductionPreCounter uint64
	HasIntroductionMock       mActiveNodeMockHasIntroduction

	IsAcceptableHostFunc       func(p endpoints.Inbound) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mActiveNodeMockIsAcceptableHost

	IsJoinerFunc       func() (r bool)
	IsJoinerCounter    uint64
	IsJoinerPreCounter uint64
	IsJoinerMock       mActiveNodeMockIsJoiner
}

func (m *ActiveNodeMock) GetStaticNodeID() insolar.ShortNodeID {
	return m.GetNodeID()
}

func (m *ActiveNodeMock) GetStatic() StaticProfile {
	return m
}

//NewActiveNodeMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode
func NewActiveNodeMock(t minimock.Tester) *ActiveNodeMock {
	m := &ActiveNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAnnouncementSignatureMock = mActiveNodeMockGetAnnouncementSignature{mock: m}
	m.GetDeclaredPowerMock = mActiveNodeMockGetDeclaredPower{mock: m}
	m.GetDefaultEndpointMock = mActiveNodeMockGetDefaultEndpoint{mock: m}
	m.GetIndexMock = mActiveNodeMockGetIndex{mock: m}
	m.GetIntroductionMock = mActiveNodeMockGetIntroduction{mock: m}
	m.GetNodePublicKeyMock = mActiveNodeMockGetNodePublicKey{mock: m}
	m.GetOpModeMock = mActiveNodeMockGetOpMode{mock: m}
	m.GetPrimaryRoleMock = mActiveNodeMockGetPrimaryRole{mock: m}
	m.GetPublicKeyStoreMock = mActiveNodeMockGetPublicKeyStore{mock: m}
	m.GetShortNodeIDMock = mActiveNodeMockGetShortNodeID{mock: m}
	m.GetSignatureVerifierMock = mActiveNodeMockGetSignatureVerifier{mock: m}
	m.GetSpecialRolesMock = mActiveNodeMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mActiveNodeMockGetStartPower{mock: m}
	m.HasIntroductionMock = mActiveNodeMockHasIntroduction{mock: m}
	m.IsAcceptableHostMock = mActiveNodeMockIsAcceptableHost{mock: m}
	m.IsJoinerMock = mActiveNodeMockIsJoiner{mock: m}

	return m
}

type mActiveNodeMockGetAnnouncementSignature struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetAnnouncementSignatureExpectation
	expectationSeries []*ActiveNodeMockGetAnnouncementSignatureExpectation
}

type ActiveNodeMockGetAnnouncementSignatureExpectation struct {
	result *ActiveNodeMockGetAnnouncementSignatureResult
}

type ActiveNodeMockGetAnnouncementSignatureResult struct {
	r cryptkit.SignatureHolder
}

//Expect specifies that invocation of ActiveNode.GetAnnouncementSignature is expected from 1 to Infinity times
func (m *mActiveNodeMockGetAnnouncementSignature) Expect() *mActiveNodeMockGetAnnouncementSignature {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetAnnouncementSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetAnnouncementSignature
func (m *mActiveNodeMockGetAnnouncementSignature) Return(r cryptkit.SignatureHolder) *ActiveNodeMock {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetAnnouncementSignatureExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetAnnouncementSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetAnnouncementSignature is expected once
func (m *mActiveNodeMockGetAnnouncementSignature) ExpectOnce() *ActiveNodeMockGetAnnouncementSignatureExpectation {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetAnnouncementSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetAnnouncementSignatureExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &ActiveNodeMockGetAnnouncementSignatureResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetAnnouncementSignature method
func (m *mActiveNodeMockGetAnnouncementSignature) Set(f func() (r cryptkit.SignatureHolder)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAnnouncementSignatureFunc = f
	return m.mock
}

//GetAnnouncementSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetAnnouncementSignature() (r cryptkit.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetAnnouncementSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetAnnouncementSignatureCounter, 1)

	if len(m.GetAnnouncementSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAnnouncementSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetAnnouncementSignature.")
			return
		}

		result := m.GetAnnouncementSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetAnnouncementSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureMock.mainExpectation != nil {

		result := m.GetAnnouncementSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetAnnouncementSignature")
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetAnnouncementSignature.")
		return
	}

	return m.GetAnnouncementSignatureFunc()
}

//GetAnnouncementSignatureMinimockCounter returns a count of ActiveNodeMock.GetAnnouncementSignatureFunc invocations
func (m *ActiveNodeMock) GetAnnouncementSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter)
}

//GetAnnouncementSignatureMinimockPreCounter returns the value of ActiveNodeMock.GetAnnouncementSignature invocations
func (m *ActiveNodeMock) GetAnnouncementSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignaturePreCounter)
}

//GetAnnouncementSignatureFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetAnnouncementSignatureFinished() bool {
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

type mActiveNodeMockGetDeclaredPower struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetDeclaredPowerExpectation
	expectationSeries []*ActiveNodeMockGetDeclaredPowerExpectation
}

type ActiveNodeMockGetDeclaredPowerExpectation struct {
	result *ActiveNodeMockGetDeclaredPowerResult
}

type ActiveNodeMockGetDeclaredPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of ActiveNode.GetDeclaredPower is expected from 1 to Infinity times
func (m *mActiveNodeMockGetDeclaredPower) Expect() *mActiveNodeMockGetDeclaredPower {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetDeclaredPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetDeclaredPower
func (m *mActiveNodeMockGetDeclaredPower) Return(r member.Power) *ActiveNodeMock {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetDeclaredPowerExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetDeclaredPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetDeclaredPower is expected once
func (m *mActiveNodeMockGetDeclaredPower) ExpectOnce() *ActiveNodeMockGetDeclaredPowerExpectation {
	m.mock.GetDeclaredPowerFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetDeclaredPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetDeclaredPowerExpectation) Return(r member.Power) {
	e.result = &ActiveNodeMockGetDeclaredPowerResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetDeclaredPower method
func (m *mActiveNodeMockGetDeclaredPower) Set(f func() (r member.Power)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDeclaredPowerFunc = f
	return m.mock
}

//GetDeclaredPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetDeclaredPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetDeclaredPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetDeclaredPowerCounter, 1)

	if len(m.GetDeclaredPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDeclaredPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetDeclaredPower.")
			return
		}

		result := m.GetDeclaredPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetDeclaredPower")
			return
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerMock.mainExpectation != nil {

		result := m.GetDeclaredPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetDeclaredPower")
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetDeclaredPower.")
		return
	}

	return m.GetDeclaredPowerFunc()
}

//GetDeclaredPowerMinimockCounter returns a count of ActiveNodeMock.GetDeclaredPowerFunc invocations
func (m *ActiveNodeMock) GetDeclaredPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerCounter)
}

//GetDeclaredPowerMinimockPreCounter returns the value of ActiveNodeMock.GetDeclaredPower invocations
func (m *ActiveNodeMock) GetDeclaredPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerPreCounter)
}

//GetDeclaredPowerFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetDeclaredPowerFinished() bool {
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

type mActiveNodeMockGetDefaultEndpoint struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetDefaultEndpointExpectation
	expectationSeries []*ActiveNodeMockGetDefaultEndpointExpectation
}

type ActiveNodeMockGetDefaultEndpointExpectation struct {
	result *ActiveNodeMockGetDefaultEndpointResult
}

type ActiveNodeMockGetDefaultEndpointResult struct {
	r endpoints.Outbound
}

//Expect specifies that invocation of ActiveNode.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mActiveNodeMockGetDefaultEndpoint) Expect() *mActiveNodeMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetDefaultEndpoint
func (m *mActiveNodeMockGetDefaultEndpoint) Return(r endpoints.Outbound) *ActiveNodeMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetDefaultEndpoint is expected once
func (m *mActiveNodeMockGetDefaultEndpoint) ExpectOnce() *ActiveNodeMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetDefaultEndpointExpectation) Return(r endpoints.Outbound) {
	e.result = &ActiveNodeMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetDefaultEndpoint method
func (m *mActiveNodeMockGetDefaultEndpoint) Set(f func() (r endpoints.Outbound)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetDefaultEndpoint() (r endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of ActiveNodeMock.GetDefaultEndpointFunc invocations
func (m *ActiveNodeMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of ActiveNodeMock.GetDefaultEndpoint invocations
func (m *ActiveNodeMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetDefaultEndpointFinished() bool {
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

type mActiveNodeMockGetIndex struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetIndexExpectation
	expectationSeries []*ActiveNodeMockGetIndexExpectation
}

type ActiveNodeMockGetIndexExpectation struct {
	result *ActiveNodeMockGetIndexResult
}

type ActiveNodeMockGetIndexResult struct {
	r member.Index
}

//Expect specifies that invocation of ActiveNode.GetIndex is expected from 1 to Infinity times
func (m *mActiveNodeMockGetIndex) Expect() *mActiveNodeMockGetIndex {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetIndexExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetIndex
func (m *mActiveNodeMockGetIndex) Return(r member.Index) *ActiveNodeMock {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetIndexExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetIndex is expected once
func (m *mActiveNodeMockGetIndex) ExpectOnce() *ActiveNodeMockGetIndexExpectation {
	m.mock.GetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetIndexExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetIndexExpectation) Return(r member.Index) {
	e.result = &ActiveNodeMockGetIndexResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetIndex method
func (m *mActiveNodeMockGetIndex) Set(f func() (r member.Index)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIndexFunc = f
	return m.mock
}

//GetIndex implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetIndex() (r member.Index) {
	counter := atomic.AddUint64(&m.GetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.GetIndexCounter, 1)

	if len(m.GetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetIndex.")
			return
		}

		result := m.GetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetIndex")
			return
		}

		r = result.r

		return
	}

	if m.GetIndexMock.mainExpectation != nil {

		result := m.GetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetIndex")
		}

		r = result.r

		return
	}

	if m.GetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetIndex.")
		return
	}

	return m.GetIndexFunc()
}

//GetIndexMinimockCounter returns a count of ActiveNodeMock.GetIndexFunc invocations
func (m *ActiveNodeMock) GetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexCounter)
}

//GetIndexMinimockPreCounter returns the value of ActiveNodeMock.GetIndex invocations
func (m *ActiveNodeMock) GetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexPreCounter)
}

//GetIndexFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetIndexFinished() bool {
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

type mActiveNodeMockGetIntroduction struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetIntroductionExpectation
	expectationSeries []*ActiveNodeMockGetIntroductionExpectation
}

type ActiveNodeMockGetIntroductionExpectation struct {
	result *ActiveNodeMockGetIntroductionResult
}

type ActiveNodeMockGetIntroductionResult struct {
	r NodeIntroduction
}

//Expect specifies that invocation of ActiveNode.GetIntroduction is expected from 1 to Infinity times
func (m *mActiveNodeMockGetIntroduction) Expect() *mActiveNodeMockGetIntroduction {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetIntroduction
func (m *mActiveNodeMockGetIntroduction) Return(r NodeIntroduction) *ActiveNodeMock {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetIntroductionExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetIntroduction is expected once
func (m *mActiveNodeMockGetIntroduction) ExpectOnce() *ActiveNodeMockGetIntroductionExpectation {
	m.mock.GetIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetIntroductionExpectation) Return(r NodeIntroduction) {
	e.result = &ActiveNodeMockGetIntroductionResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetIntroduction method
func (m *mActiveNodeMockGetIntroduction) Set(f func() (r NodeIntroduction)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIntroductionFunc = f
	return m.mock
}

//GetIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetIntroduction() (r NodeIntroduction) {
	counter := atomic.AddUint64(&m.GetIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.GetIntroductionCounter, 1)

	if len(m.GetIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetIntroduction.")
			return
		}

		result := m.GetIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.GetIntroductionMock.mainExpectation != nil {

		result := m.GetIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetIntroduction")
		}

		r = result.r

		return
	}

	if m.GetIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetIntroduction.")
		return
	}

	return m.GetIntroductionFunc()
}

//GetIntroductionMinimockCounter returns a count of ActiveNodeMock.GetIntroductionFunc invocations
func (m *ActiveNodeMock) GetIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionCounter)
}

//GetIntroductionMinimockPreCounter returns the value of ActiveNodeMock.GetIntroduction invocations
func (m *ActiveNodeMock) GetIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionPreCounter)
}

//GetIntroductionFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetIntroductionFinished() bool {
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

type mActiveNodeMockGetNodePublicKey struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetNodePublicKeyExpectation
	expectationSeries []*ActiveNodeMockGetNodePublicKeyExpectation
}

type ActiveNodeMockGetNodePublicKeyExpectation struct {
	result *ActiveNodeMockGetNodePublicKeyResult
}

type ActiveNodeMockGetNodePublicKeyResult struct {
	r cryptkit.SignatureKeyHolder
}

//Expect specifies that invocation of ActiveNode.GetNodePublicKey is expected from 1 to Infinity times
func (m *mActiveNodeMockGetNodePublicKey) Expect() *mActiveNodeMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetNodePublicKey
func (m *mActiveNodeMockGetNodePublicKey) Return(r cryptkit.SignatureKeyHolder) *ActiveNodeMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetNodePublicKey is expected once
func (m *mActiveNodeMockGetNodePublicKey) ExpectOnce() *ActiveNodeMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetNodePublicKeyExpectation) Return(r cryptkit.SignatureKeyHolder) {
	e.result = &ActiveNodeMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetNodePublicKey method
func (m *mActiveNodeMockGetNodePublicKey) Set(f func() (r cryptkit.SignatureKeyHolder)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetNodePublicKey() (r cryptkit.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of ActiveNodeMock.GetNodePublicKeyFunc invocations
func (m *ActiveNodeMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of ActiveNodeMock.GetNodePublicKey invocations
func (m *ActiveNodeMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetNodePublicKeyFinished() bool {
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

type mActiveNodeMockGetOpMode struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetOpModeExpectation
	expectationSeries []*ActiveNodeMockGetOpModeExpectation
}

type ActiveNodeMockGetOpModeExpectation struct {
	result *ActiveNodeMockGetOpModeResult
}

type ActiveNodeMockGetOpModeResult struct {
	r member.OpMode
}

//Expect specifies that invocation of ActiveNode.GetOpMode is expected from 1 to Infinity times
func (m *mActiveNodeMockGetOpMode) Expect() *mActiveNodeMockGetOpMode {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetOpModeExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetOpMode
func (m *mActiveNodeMockGetOpMode) Return(r member.OpMode) *ActiveNodeMock {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetOpModeExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetOpModeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetOpMode is expected once
func (m *mActiveNodeMockGetOpMode) ExpectOnce() *ActiveNodeMockGetOpModeExpectation {
	m.mock.GetOpModeFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetOpModeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetOpModeExpectation) Return(r member.OpMode) {
	e.result = &ActiveNodeMockGetOpModeResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetOpMode method
func (m *mActiveNodeMockGetOpMode) Set(f func() (r member.OpMode)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOpModeFunc = f
	return m.mock
}

//GetOpMode implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetOpMode() (r member.OpMode) {
	counter := atomic.AddUint64(&m.GetOpModePreCounter, 1)
	defer atomic.AddUint64(&m.GetOpModeCounter, 1)

	if len(m.GetOpModeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOpModeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetOpMode.")
			return
		}

		result := m.GetOpModeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetOpMode")
			return
		}

		r = result.r

		return
	}

	if m.GetOpModeMock.mainExpectation != nil {

		result := m.GetOpModeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetOpMode")
		}

		r = result.r

		return
	}

	if m.GetOpModeFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetOpMode.")
		return
	}

	return m.GetOpModeFunc()
}

//GetOpModeMinimockCounter returns a count of ActiveNodeMock.GetOpModeFunc invocations
func (m *ActiveNodeMock) GetOpModeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModeCounter)
}

//GetOpModeMinimockPreCounter returns the value of ActiveNodeMock.GetOpMode invocations
func (m *ActiveNodeMock) GetOpModeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModePreCounter)
}

//GetOpModeFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetOpModeFinished() bool {
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

type mActiveNodeMockGetPrimaryRole struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetPrimaryRoleExpectation
	expectationSeries []*ActiveNodeMockGetPrimaryRoleExpectation
}

type ActiveNodeMockGetPrimaryRoleExpectation struct {
	result *ActiveNodeMockGetPrimaryRoleResult
}

type ActiveNodeMockGetPrimaryRoleResult struct {
	r member.PrimaryRole
}

//Expect specifies that invocation of ActiveNode.GetPrimaryRole is expected from 1 to Infinity times
func (m *mActiveNodeMockGetPrimaryRole) Expect() *mActiveNodeMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetPrimaryRole
func (m *mActiveNodeMockGetPrimaryRole) Return(r member.PrimaryRole) *ActiveNodeMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetPrimaryRole is expected once
func (m *mActiveNodeMockGetPrimaryRole) ExpectOnce() *ActiveNodeMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetPrimaryRoleExpectation) Return(r member.PrimaryRole) {
	e.result = &ActiveNodeMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetPrimaryRole method
func (m *mActiveNodeMockGetPrimaryRole) Set(f func() (r member.PrimaryRole)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetPrimaryRole() (r member.PrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of ActiveNodeMock.GetPrimaryRoleFunc invocations
func (m *ActiveNodeMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of ActiveNodeMock.GetPrimaryRole invocations
func (m *ActiveNodeMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetPrimaryRoleFinished() bool {
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

type mActiveNodeMockGetPublicKeyStore struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetPublicKeyStoreExpectation
	expectationSeries []*ActiveNodeMockGetPublicKeyStoreExpectation
}

type ActiveNodeMockGetPublicKeyStoreExpectation struct {
	result *ActiveNodeMockGetPublicKeyStoreResult
}

type ActiveNodeMockGetPublicKeyStoreResult struct {
	r cryptkit.PublicKeyStore
}

//Expect specifies that invocation of ActiveNode.GetPublicKeyStore is expected from 1 to Infinity times
func (m *mActiveNodeMockGetPublicKeyStore) Expect() *mActiveNodeMockGetPublicKeyStore {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetPublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetPublicKeyStore
func (m *mActiveNodeMockGetPublicKeyStore) Return(r cryptkit.PublicKeyStore) *ActiveNodeMock {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetPublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetPublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetPublicKeyStore is expected once
func (m *mActiveNodeMockGetPublicKeyStore) ExpectOnce() *ActiveNodeMockGetPublicKeyStoreExpectation {
	m.mock.GetPublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetPublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetPublicKeyStoreExpectation) Return(r cryptkit.PublicKeyStore) {
	e.result = &ActiveNodeMockGetPublicKeyStoreResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetPublicKeyStore method
func (m *mActiveNodeMockGetPublicKeyStore) Set(f func() (r cryptkit.PublicKeyStore)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyStoreFunc = f
	return m.mock
}

//GetPublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetPublicKeyStore() (r cryptkit.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetPublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyStoreCounter, 1)

	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetPublicKeyStore.")
			return
		}

		result := m.GetPublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetPublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreMock.mainExpectation != nil {

		result := m.GetPublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetPublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetPublicKeyStore.")
		return
	}

	return m.GetPublicKeyStoreFunc()
}

//GetPublicKeyStoreMinimockCounter returns a count of ActiveNodeMock.GetPublicKeyStoreFunc invocations
func (m *ActiveNodeMock) GetPublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStoreCounter)
}

//GetPublicKeyStoreMinimockPreCounter returns the value of ActiveNodeMock.GetPublicKeyStore invocations
func (m *ActiveNodeMock) GetPublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStorePreCounter)
}

//GetPublicKeyStoreFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetPublicKeyStoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) == uint64(len(m.GetPublicKeyStoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPublicKeyStoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPublicKeyStoreFunc != nil {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) > 0
	}

	return true
}

type mActiveNodeMockGetShortNodeID struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetShortNodeIDExpectation
	expectationSeries []*ActiveNodeMockGetShortNodeIDExpectation
}

type ActiveNodeMockGetShortNodeIDExpectation struct {
	result *ActiveNodeMockGetShortNodeIDResult
}

type ActiveNodeMockGetShortNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of ActiveNode.GetNodeID is expected from 1 to Infinity times
func (m *mActiveNodeMockGetShortNodeID) Expect() *mActiveNodeMockGetShortNodeID {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetShortNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetNodeID
func (m *mActiveNodeMockGetShortNodeID) Return(r insolar.ShortNodeID) *ActiveNodeMock {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetShortNodeIDExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetShortNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetNodeID is expected once
func (m *mActiveNodeMockGetShortNodeID) ExpectOnce() *ActiveNodeMockGetShortNodeIDExpectation {
	m.mock.GetShortNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetShortNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetShortNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &ActiveNodeMockGetShortNodeIDResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetNodeID method
func (m *mActiveNodeMockGetShortNodeID) Set(f func() (r insolar.ShortNodeID)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetShortNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetShortNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetShortNodeIDCounter, 1)

	if len(m.GetShortNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetShortNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetNodeID.")
			return
		}

		result := m.GetShortNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDMock.mainExpectation != nil {

		result := m.GetShortNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetNodeID.")
		return
	}

	return m.GetShortNodeIDFunc()
}

//GetShortNodeIDMinimockCounter returns a count of ActiveNodeMock.GetShortNodeIDFunc invocations
func (m *ActiveNodeMock) GetShortNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDCounter)
}

//GetShortNodeIDMinimockPreCounter returns the value of ActiveNodeMock.GetNodeID invocations
func (m *ActiveNodeMock) GetShortNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDPreCounter)
}

//GetShortNodeIDFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetShortNodeIDFinished() bool {
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

type mActiveNodeMockGetSignatureVerifier struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetSignatureVerifierExpectation
	expectationSeries []*ActiveNodeMockGetSignatureVerifierExpectation
}

type ActiveNodeMockGetSignatureVerifierExpectation struct {
	result *ActiveNodeMockGetSignatureVerifierResult
}

type ActiveNodeMockGetSignatureVerifierResult struct {
	r cryptkit.SignatureVerifier
}

//Expect specifies that invocation of ActiveNode.GetSignatureVerifier is expected from 1 to Infinity times
func (m *mActiveNodeMockGetSignatureVerifier) Expect() *mActiveNodeMockGetSignatureVerifier {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetSignatureVerifierExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetSignatureVerifier
func (m *mActiveNodeMockGetSignatureVerifier) Return(r cryptkit.SignatureVerifier) *ActiveNodeMock {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetSignatureVerifierExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetSignatureVerifierResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetSignatureVerifier is expected once
func (m *mActiveNodeMockGetSignatureVerifier) ExpectOnce() *ActiveNodeMockGetSignatureVerifierExpectation {
	m.mock.GetSignatureVerifierFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetSignatureVerifierExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetSignatureVerifierExpectation) Return(r cryptkit.SignatureVerifier) {
	e.result = &ActiveNodeMockGetSignatureVerifierResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetSignatureVerifier method
func (m *mActiveNodeMockGetSignatureVerifier) Set(f func() (r cryptkit.SignatureVerifier)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureVerifierFunc = f
	return m.mock
}

//GetSignatureVerifier implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetSignatureVerifier() (r cryptkit.SignatureVerifier) {
	counter := atomic.AddUint64(&m.GetSignatureVerifierPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureVerifierCounter, 1)

	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureVerifierMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetSignatureVerifier.")
			return
		}

		result := m.GetSignatureVerifierMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetSignatureVerifier")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierMock.mainExpectation != nil {

		result := m.GetSignatureVerifierMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetSignatureVerifier")
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetSignatureVerifier.")
		return
	}

	return m.GetSignatureVerifierFunc()
}

//GetSignatureVerifierMinimockCounter returns a count of ActiveNodeMock.GetSignatureVerifierFunc invocations
func (m *ActiveNodeMock) GetSignatureVerifierMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierCounter)
}

//GetSignatureVerifierMinimockPreCounter returns the value of ActiveNodeMock.GetSignatureVerifier invocations
func (m *ActiveNodeMock) GetSignatureVerifierMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierPreCounter)
}

//GetSignatureVerifierFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetSignatureVerifierFinished() bool {
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

type mActiveNodeMockGetSpecialRoles struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetSpecialRolesExpectation
	expectationSeries []*ActiveNodeMockGetSpecialRolesExpectation
}

type ActiveNodeMockGetSpecialRolesExpectation struct {
	result *ActiveNodeMockGetSpecialRolesResult
}

type ActiveNodeMockGetSpecialRolesResult struct {
	r member.SpecialRole
}

//Expect specifies that invocation of ActiveNode.GetSpecialRoles is expected from 1 to Infinity times
func (m *mActiveNodeMockGetSpecialRoles) Expect() *mActiveNodeMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetSpecialRoles
func (m *mActiveNodeMockGetSpecialRoles) Return(r member.SpecialRole) *ActiveNodeMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetSpecialRoles is expected once
func (m *mActiveNodeMockGetSpecialRoles) ExpectOnce() *ActiveNodeMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetSpecialRolesExpectation) Return(r member.SpecialRole) {
	e.result = &ActiveNodeMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetSpecialRoles method
func (m *mActiveNodeMockGetSpecialRoles) Set(f func() (r member.SpecialRole)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetSpecialRoles() (r member.SpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of ActiveNodeMock.GetSpecialRolesFunc invocations
func (m *ActiveNodeMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of ActiveNodeMock.GetSpecialRoles invocations
func (m *ActiveNodeMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetSpecialRolesFinished() bool {
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

type mActiveNodeMockGetStartPower struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockGetStartPowerExpectation
	expectationSeries []*ActiveNodeMockGetStartPowerExpectation
}

type ActiveNodeMockGetStartPowerExpectation struct {
	result *ActiveNodeMockGetStartPowerResult
}

type ActiveNodeMockGetStartPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of ActiveNode.GetStartPower is expected from 1 to Infinity times
func (m *mActiveNodeMockGetStartPower) Expect() *mActiveNodeMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.GetStartPower
func (m *mActiveNodeMockGetStartPower) Return(r member.Power) *ActiveNodeMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.GetStartPower is expected once
func (m *mActiveNodeMockGetStartPower) ExpectOnce() *ActiveNodeMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockGetStartPowerExpectation) Return(r member.Power) {
	e.result = &ActiveNodeMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of ActiveNode.GetStartPower method
func (m *mActiveNodeMockGetStartPower) Set(f func() (r member.Power)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) GetStartPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of ActiveNodeMock.GetStartPowerFunc invocations
func (m *ActiveNodeMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of ActiveNodeMock.GetStartPower invocations
func (m *ActiveNodeMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) GetStartPowerFinished() bool {
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

type mActiveNodeMockHasIntroduction struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockHasIntroductionExpectation
	expectationSeries []*ActiveNodeMockHasIntroductionExpectation
}

type ActiveNodeMockHasIntroductionExpectation struct {
	result *ActiveNodeMockHasIntroductionResult
}

type ActiveNodeMockHasIntroductionResult struct {
	r bool
}

//Expect specifies that invocation of ActiveNode.HasIntroduction is expected from 1 to Infinity times
func (m *mActiveNodeMockHasIntroduction) Expect() *mActiveNodeMockHasIntroduction {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockHasIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.HasIntroduction
func (m *mActiveNodeMockHasIntroduction) Return(r bool) *ActiveNodeMock {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockHasIntroductionExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockHasIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.HasIntroduction is expected once
func (m *mActiveNodeMockHasIntroduction) ExpectOnce() *ActiveNodeMockHasIntroductionExpectation {
	m.mock.HasIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockHasIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockHasIntroductionExpectation) Return(r bool) {
	e.result = &ActiveNodeMockHasIntroductionResult{r}
}

//Set uses given function f as a mock of ActiveNode.HasIntroduction method
func (m *mActiveNodeMockHasIntroduction) Set(f func() (r bool)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasIntroductionFunc = f
	return m.mock
}

//HasIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) HasIntroduction() (r bool) {
	counter := atomic.AddUint64(&m.HasIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.HasIntroductionCounter, 1)

	if len(m.HasIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.HasIntroduction.")
			return
		}

		result := m.HasIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.HasIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.HasIntroductionMock.mainExpectation != nil {

		result := m.HasIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.HasIntroduction")
		}

		r = result.r

		return
	}

	if m.HasIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.HasIntroduction.")
		return
	}

	return m.HasIntroductionFunc()
}

//HasIntroductionMinimockCounter returns a count of ActiveNodeMock.HasIntroductionFunc invocations
func (m *ActiveNodeMock) HasIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionCounter)
}

//HasIntroductionMinimockPreCounter returns the value of ActiveNodeMock.HasIntroduction invocations
func (m *ActiveNodeMock) HasIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionPreCounter)
}

//HasIntroductionFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) HasIntroductionFinished() bool {
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

type mActiveNodeMockIsAcceptableHost struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockIsAcceptableHostExpectation
	expectationSeries []*ActiveNodeMockIsAcceptableHostExpectation
}

type ActiveNodeMockIsAcceptableHostExpectation struct {
	input  *ActiveNodeMockIsAcceptableHostInput
	result *ActiveNodeMockIsAcceptableHostResult
}

type ActiveNodeMockIsAcceptableHostInput struct {
	p endpoints.Inbound
}

type ActiveNodeMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of ActiveNode.IsAcceptableHost is expected from 1 to Infinity times
func (m *mActiveNodeMockIsAcceptableHost) Expect(p endpoints.Inbound) *mActiveNodeMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &ActiveNodeMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of ActiveNode.IsAcceptableHost
func (m *mActiveNodeMockIsAcceptableHost) Return(r bool) *ActiveNodeMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.IsAcceptableHost is expected once
func (m *mActiveNodeMockIsAcceptableHost) ExpectOnce(p endpoints.Inbound) *ActiveNodeMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockIsAcceptableHostExpectation{}
	expectation.input = &ActiveNodeMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &ActiveNodeMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of ActiveNode.IsAcceptableHost method
func (m *mActiveNodeMockIsAcceptableHost) Set(f func(p endpoints.Inbound) (r bool)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) IsAcceptableHost(p endpoints.Inbound) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ActiveNodeMockIsAcceptableHostInput{p}, "ActiveNode.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ActiveNodeMockIsAcceptableHostInput{p}, "ActiveNode.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of ActiveNodeMock.IsAcceptableHostFunc invocations
func (m *ActiveNodeMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of ActiveNodeMock.IsAcceptableHost invocations
func (m *ActiveNodeMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) IsAcceptableHostFinished() bool {
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

type mActiveNodeMockIsJoiner struct {
	mock              *ActiveNodeMock
	mainExpectation   *ActiveNodeMockIsJoinerExpectation
	expectationSeries []*ActiveNodeMockIsJoinerExpectation
}

type ActiveNodeMockIsJoinerExpectation struct {
	result *ActiveNodeMockIsJoinerResult
}

type ActiveNodeMockIsJoinerResult struct {
	r bool
}

//Expect specifies that invocation of ActiveNode.IsJoiner is expected from 1 to Infinity times
func (m *mActiveNodeMockIsJoiner) Expect() *mActiveNodeMockIsJoiner {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockIsJoinerExpectation{}
	}

	return m
}

//Return specifies results of invocation of ActiveNode.IsJoiner
func (m *mActiveNodeMockIsJoiner) Return(r bool) *ActiveNodeMock {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodeMockIsJoinerExpectation{}
	}
	m.mainExpectation.result = &ActiveNodeMockIsJoinerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNode.IsJoiner is expected once
func (m *mActiveNodeMockIsJoiner) ExpectOnce() *ActiveNodeMockIsJoinerExpectation {
	m.mock.IsJoinerFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodeMockIsJoinerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodeMockIsJoinerExpectation) Return(r bool) {
	e.result = &ActiveNodeMockIsJoinerResult{r}
}

//Set uses given function f as a mock of ActiveNode.IsJoiner method
func (m *mActiveNodeMockIsJoiner) Set(f func() (r bool)) *ActiveNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsJoinerFunc = f
	return m.mock
}

//IsJoiner implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode interface
func (m *ActiveNodeMock) IsJoiner() (r bool) {
	counter := atomic.AddUint64(&m.IsJoinerPreCounter, 1)
	defer atomic.AddUint64(&m.IsJoinerCounter, 1)

	if len(m.IsJoinerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsJoinerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodeMock.IsJoiner.")
			return
		}

		result := m.IsJoinerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.IsJoiner")
			return
		}

		r = result.r

		return
	}

	if m.IsJoinerMock.mainExpectation != nil {

		result := m.IsJoinerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodeMock.IsJoiner")
		}

		r = result.r

		return
	}

	if m.IsJoinerFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodeMock.IsJoiner.")
		return
	}

	return m.IsJoinerFunc()
}

//IsJoinerMinimockCounter returns a count of ActiveNodeMock.IsJoinerFunc invocations
func (m *ActiveNodeMock) IsJoinerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerCounter)
}

//IsJoinerMinimockPreCounter returns the value of ActiveNodeMock.IsJoiner invocations
func (m *ActiveNodeMock) IsJoinerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerPreCounter)
}

//IsJoinerFinished returns true if mock invocations count is ok
func (m *ActiveNodeMock) IsJoinerFinished() bool {
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
func (m *ActiveNodeMock) ValidateCallCounters() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetNodePublicKey")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetOpMode")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetPrimaryRole")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetPublicKeyStore")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.IsAcceptableHost")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.IsJoiner")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ActiveNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ActiveNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ActiveNodeMock) MinimockFinish() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetNodePublicKey")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetOpMode")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetPrimaryRole")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetPublicKeyStore")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.IsAcceptableHost")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to ActiveNodeMock.IsJoiner")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ActiveNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ActiveNodeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetAnnouncementSignatureFinished()
		ok = ok && m.GetDeclaredPowerFinished()
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetIndexFinished()
		ok = ok && m.GetIntroductionFinished()
		ok = ok && m.GetNodePublicKeyFinished()
		ok = ok && m.GetOpModeFinished()
		ok = ok && m.GetPrimaryRoleFinished()
		ok = ok && m.GetPublicKeyStoreFinished()
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
				m.t.Error("Expected call to ActiveNodeMock.GetAnnouncementSignature")
			}

			if !m.GetDeclaredPowerFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetDeclaredPower")
			}

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetDefaultEndpoint")
			}

			if !m.GetIndexFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetIndex")
			}

			if !m.GetIntroductionFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetIntroduction")
			}

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetNodePublicKey")
			}

			if !m.GetOpModeFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetOpMode")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetPrimaryRole")
			}

			if !m.GetPublicKeyStoreFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetPublicKeyStore")
			}

			if !m.GetShortNodeIDFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetNodeID")
			}

			if !m.GetSignatureVerifierFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetSignatureVerifier")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetSpecialRoles")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to ActiveNodeMock.GetStartPower")
			}

			if !m.HasIntroductionFinished() {
				m.t.Error("Expected call to ActiveNodeMock.HasIntroduction")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to ActiveNodeMock.IsAcceptableHost")
			}

			if !m.IsJoinerFinished() {
				m.t.Error("Expected call to ActiveNodeMock.IsJoiner")
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
func (m *ActiveNodeMock) AllMocksCalled() bool {

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

	if !m.GetOpModeFinished() {
		return false
	}

	if !m.GetPrimaryRoleFinished() {
		return false
	}

	if !m.GetPublicKeyStoreFinished() {
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

package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LocalNode" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
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

//LocalNodeMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode
type LocalNodeMock struct {
	t minimock.Tester

	GetAnnouncementSignatureFunc       func() (r cryptkit.SignatureHolder)
	GetAnnouncementSignatureCounter    uint64
	GetAnnouncementSignaturePreCounter uint64
	GetAnnouncementSignatureMock       mLocalNodeMockGetAnnouncementSignature

	GetDeclaredPowerFunc       func() (r member.Power)
	GetDeclaredPowerCounter    uint64
	GetDeclaredPowerPreCounter uint64
	GetDeclaredPowerMock       mLocalNodeMockGetDeclaredPower

	GetDefaultEndpointFunc       func() (r endpoints.Outbound)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mLocalNodeMockGetDefaultEndpoint

	GetIndexFunc       func() (r member.Index)
	GetIndexCounter    uint64
	GetIndexPreCounter uint64
	GetIndexMock       mLocalNodeMockGetIndex

	GetIntroductionFunc       func() (r NodeIntroduction)
	GetIntroductionCounter    uint64
	GetIntroductionPreCounter uint64
	GetIntroductionMock       mLocalNodeMockGetIntroduction

	GetNodePublicKeyFunc       func() (r cryptkit.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mLocalNodeMockGetNodePublicKey

	GetOpModeFunc       func() (r member.OpMode)
	GetOpModeCounter    uint64
	GetOpModePreCounter uint64
	GetOpModeMock       mLocalNodeMockGetOpMode

	GetPrimaryRoleFunc       func() (r member.PrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mLocalNodeMockGetPrimaryRole

	GetPublicKeyStoreFunc       func() (r cryptkit.PublicKeyStore)
	GetPublicKeyStoreCounter    uint64
	GetPublicKeyStorePreCounter uint64
	GetPublicKeyStoreMock       mLocalNodeMockGetPublicKeyStore

	GetShortNodeIDFunc       func() (r insolar.ShortNodeID)
	GetShortNodeIDCounter    uint64
	GetShortNodeIDPreCounter uint64
	GetShortNodeIDMock       mLocalNodeMockGetShortNodeID

	GetSignatureVerifierFunc       func() (r cryptkit.SignatureVerifier)
	GetSignatureVerifierCounter    uint64
	GetSignatureVerifierPreCounter uint64
	GetSignatureVerifierMock       mLocalNodeMockGetSignatureVerifier

	GetSpecialRolesFunc       func() (r member.SpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mLocalNodeMockGetSpecialRoles

	GetStartPowerFunc       func() (r member.Power)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mLocalNodeMockGetStartPower

	HasIntroductionFunc       func() (r bool)
	HasIntroductionCounter    uint64
	HasIntroductionPreCounter uint64
	HasIntroductionMock       mLocalNodeMockHasIntroduction

	IsAcceptableHostFunc       func(p endpoints.Inbound) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mLocalNodeMockIsAcceptableHost

	IsJoinerFunc       func() (r bool)
	IsJoinerCounter    uint64
	IsJoinerPreCounter uint64
	IsJoinerMock       mLocalNodeMockIsJoiner

	LocalNodeProfileFunc       func()
	LocalNodeProfileCounter    uint64
	LocalNodeProfilePreCounter uint64
	LocalNodeProfileMock       mLocalNodeMockLocalNodeProfile
}

func (m *LocalNodeMock) GetStaticNodeID() insolar.ShortNodeID {
	return m.GetNodeID()
}

func (m *LocalNodeMock) GetStatic() StaticProfile {
	return m
}

//NewLocalNodeMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode
func NewLocalNodeMock(t minimock.Tester) *LocalNodeMock {
	m := &LocalNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAnnouncementSignatureMock = mLocalNodeMockGetAnnouncementSignature{mock: m}
	m.GetDeclaredPowerMock = mLocalNodeMockGetDeclaredPower{mock: m}
	m.GetDefaultEndpointMock = mLocalNodeMockGetDefaultEndpoint{mock: m}
	m.GetIndexMock = mLocalNodeMockGetIndex{mock: m}
	m.GetIntroductionMock = mLocalNodeMockGetIntroduction{mock: m}
	m.GetNodePublicKeyMock = mLocalNodeMockGetNodePublicKey{mock: m}
	m.GetOpModeMock = mLocalNodeMockGetOpMode{mock: m}
	m.GetPrimaryRoleMock = mLocalNodeMockGetPrimaryRole{mock: m}
	m.GetPublicKeyStoreMock = mLocalNodeMockGetPublicKeyStore{mock: m}
	m.GetShortNodeIDMock = mLocalNodeMockGetShortNodeID{mock: m}
	m.GetSignatureVerifierMock = mLocalNodeMockGetSignatureVerifier{mock: m}
	m.GetSpecialRolesMock = mLocalNodeMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mLocalNodeMockGetStartPower{mock: m}
	m.HasIntroductionMock = mLocalNodeMockHasIntroduction{mock: m}
	m.IsAcceptableHostMock = mLocalNodeMockIsAcceptableHost{mock: m}
	m.IsJoinerMock = mLocalNodeMockIsJoiner{mock: m}
	m.LocalNodeProfileMock = mLocalNodeMockLocalNodeProfile{mock: m}

	return m
}

type mLocalNodeMockGetAnnouncementSignature struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetAnnouncementSignatureExpectation
	expectationSeries []*LocalNodeMockGetAnnouncementSignatureExpectation
}

type LocalNodeMockGetAnnouncementSignatureExpectation struct {
	result *LocalNodeMockGetAnnouncementSignatureResult
}

type LocalNodeMockGetAnnouncementSignatureResult struct {
	r cryptkit.SignatureHolder
}

//Expect specifies that invocation of LocalNode.GetAnnouncementSignature is expected from 1 to Infinity times
func (m *mLocalNodeMockGetAnnouncementSignature) Expect() *mLocalNodeMockGetAnnouncementSignature {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetAnnouncementSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetAnnouncementSignature
func (m *mLocalNodeMockGetAnnouncementSignature) Return(r cryptkit.SignatureHolder) *LocalNodeMock {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetAnnouncementSignatureExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetAnnouncementSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetAnnouncementSignature is expected once
func (m *mLocalNodeMockGetAnnouncementSignature) ExpectOnce() *LocalNodeMockGetAnnouncementSignatureExpectation {
	m.mock.GetAnnouncementSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetAnnouncementSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetAnnouncementSignatureExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &LocalNodeMockGetAnnouncementSignatureResult{r}
}

//Set uses given function f as a mock of LocalNode.GetAnnouncementSignature method
func (m *mLocalNodeMockGetAnnouncementSignature) Set(f func() (r cryptkit.SignatureHolder)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAnnouncementSignatureFunc = f
	return m.mock
}

//GetAnnouncementSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetAnnouncementSignature() (r cryptkit.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetAnnouncementSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetAnnouncementSignatureCounter, 1)

	if len(m.GetAnnouncementSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAnnouncementSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetAnnouncementSignature.")
			return
		}

		result := m.GetAnnouncementSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetAnnouncementSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureMock.mainExpectation != nil {

		result := m.GetAnnouncementSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetAnnouncementSignature")
		}

		r = result.r

		return
	}

	if m.GetAnnouncementSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetAnnouncementSignature.")
		return
	}

	return m.GetAnnouncementSignatureFunc()
}

//GetAnnouncementSignatureMinimockCounter returns a count of LocalNodeMock.GetAnnouncementSignatureFunc invocations
func (m *LocalNodeMock) GetAnnouncementSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignatureCounter)
}

//GetAnnouncementSignatureMinimockPreCounter returns the value of LocalNodeMock.GetAnnouncementSignature invocations
func (m *LocalNodeMock) GetAnnouncementSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnouncementSignaturePreCounter)
}

//GetAnnouncementSignatureFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetAnnouncementSignatureFinished() bool {
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

type mLocalNodeMockGetDeclaredPower struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetDeclaredPowerExpectation
	expectationSeries []*LocalNodeMockGetDeclaredPowerExpectation
}

type LocalNodeMockGetDeclaredPowerExpectation struct {
	result *LocalNodeMockGetDeclaredPowerResult
}

type LocalNodeMockGetDeclaredPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of LocalNode.GetDeclaredPower is expected from 1 to Infinity times
func (m *mLocalNodeMockGetDeclaredPower) Expect() *mLocalNodeMockGetDeclaredPower {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetDeclaredPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetDeclaredPower
func (m *mLocalNodeMockGetDeclaredPower) Return(r member.Power) *LocalNodeMock {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetDeclaredPowerExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetDeclaredPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetDeclaredPower is expected once
func (m *mLocalNodeMockGetDeclaredPower) ExpectOnce() *LocalNodeMockGetDeclaredPowerExpectation {
	m.mock.GetDeclaredPowerFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetDeclaredPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetDeclaredPowerExpectation) Return(r member.Power) {
	e.result = &LocalNodeMockGetDeclaredPowerResult{r}
}

//Set uses given function f as a mock of LocalNode.GetDeclaredPower method
func (m *mLocalNodeMockGetDeclaredPower) Set(f func() (r member.Power)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDeclaredPowerFunc = f
	return m.mock
}

//GetDeclaredPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetDeclaredPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetDeclaredPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetDeclaredPowerCounter, 1)

	if len(m.GetDeclaredPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDeclaredPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetDeclaredPower.")
			return
		}

		result := m.GetDeclaredPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetDeclaredPower")
			return
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerMock.mainExpectation != nil {

		result := m.GetDeclaredPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetDeclaredPower")
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetDeclaredPower.")
		return
	}

	return m.GetDeclaredPowerFunc()
}

//GetDeclaredPowerMinimockCounter returns a count of LocalNodeMock.GetDeclaredPowerFunc invocations
func (m *LocalNodeMock) GetDeclaredPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerCounter)
}

//GetDeclaredPowerMinimockPreCounter returns the value of LocalNodeMock.GetDeclaredPower invocations
func (m *LocalNodeMock) GetDeclaredPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerPreCounter)
}

//GetDeclaredPowerFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetDeclaredPowerFinished() bool {
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

type mLocalNodeMockGetDefaultEndpoint struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetDefaultEndpointExpectation
	expectationSeries []*LocalNodeMockGetDefaultEndpointExpectation
}

type LocalNodeMockGetDefaultEndpointExpectation struct {
	result *LocalNodeMockGetDefaultEndpointResult
}

type LocalNodeMockGetDefaultEndpointResult struct {
	r endpoints.Outbound
}

//Expect specifies that invocation of LocalNode.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mLocalNodeMockGetDefaultEndpoint) Expect() *mLocalNodeMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetDefaultEndpoint
func (m *mLocalNodeMockGetDefaultEndpoint) Return(r endpoints.Outbound) *LocalNodeMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetDefaultEndpoint is expected once
func (m *mLocalNodeMockGetDefaultEndpoint) ExpectOnce() *LocalNodeMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetDefaultEndpointExpectation) Return(r endpoints.Outbound) {
	e.result = &LocalNodeMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of LocalNode.GetDefaultEndpoint method
func (m *mLocalNodeMockGetDefaultEndpoint) Set(f func() (r endpoints.Outbound)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetDefaultEndpoint() (r endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of LocalNodeMock.GetDefaultEndpointFunc invocations
func (m *LocalNodeMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of LocalNodeMock.GetDefaultEndpoint invocations
func (m *LocalNodeMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetDefaultEndpointFinished() bool {
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

type mLocalNodeMockGetIndex struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetIndexExpectation
	expectationSeries []*LocalNodeMockGetIndexExpectation
}

type LocalNodeMockGetIndexExpectation struct {
	result *LocalNodeMockGetIndexResult
}

type LocalNodeMockGetIndexResult struct {
	r member.Index
}

//Expect specifies that invocation of LocalNode.GetIndex is expected from 1 to Infinity times
func (m *mLocalNodeMockGetIndex) Expect() *mLocalNodeMockGetIndex {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetIndexExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetIndex
func (m *mLocalNodeMockGetIndex) Return(r member.Index) *LocalNodeMock {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetIndexExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetIndex is expected once
func (m *mLocalNodeMockGetIndex) ExpectOnce() *LocalNodeMockGetIndexExpectation {
	m.mock.GetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetIndexExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetIndexExpectation) Return(r member.Index) {
	e.result = &LocalNodeMockGetIndexResult{r}
}

//Set uses given function f as a mock of LocalNode.GetIndex method
func (m *mLocalNodeMockGetIndex) Set(f func() (r member.Index)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIndexFunc = f
	return m.mock
}

//GetIndex implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetIndex() (r member.Index) {
	counter := atomic.AddUint64(&m.GetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.GetIndexCounter, 1)

	if len(m.GetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetIndex.")
			return
		}

		result := m.GetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetIndex")
			return
		}

		r = result.r

		return
	}

	if m.GetIndexMock.mainExpectation != nil {

		result := m.GetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetIndex")
		}

		r = result.r

		return
	}

	if m.GetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetIndex.")
		return
	}

	return m.GetIndexFunc()
}

//GetIndexMinimockCounter returns a count of LocalNodeMock.GetIndexFunc invocations
func (m *LocalNodeMock) GetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexCounter)
}

//GetIndexMinimockPreCounter returns the value of LocalNodeMock.GetIndex invocations
func (m *LocalNodeMock) GetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexPreCounter)
}

//GetIndexFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetIndexFinished() bool {
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

type mLocalNodeMockGetIntroduction struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetIntroductionExpectation
	expectationSeries []*LocalNodeMockGetIntroductionExpectation
}

type LocalNodeMockGetIntroductionExpectation struct {
	result *LocalNodeMockGetIntroductionResult
}

type LocalNodeMockGetIntroductionResult struct {
	r NodeIntroduction
}

//Expect specifies that invocation of LocalNode.GetIntroduction is expected from 1 to Infinity times
func (m *mLocalNodeMockGetIntroduction) Expect() *mLocalNodeMockGetIntroduction {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetIntroduction
func (m *mLocalNodeMockGetIntroduction) Return(r NodeIntroduction) *LocalNodeMock {
	m.mock.GetIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetIntroductionExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetIntroduction is expected once
func (m *mLocalNodeMockGetIntroduction) ExpectOnce() *LocalNodeMockGetIntroductionExpectation {
	m.mock.GetIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetIntroductionExpectation) Return(r NodeIntroduction) {
	e.result = &LocalNodeMockGetIntroductionResult{r}
}

//Set uses given function f as a mock of LocalNode.GetIntroduction method
func (m *mLocalNodeMockGetIntroduction) Set(f func() (r NodeIntroduction)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIntroductionFunc = f
	return m.mock
}

//GetIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetIntroduction() (r NodeIntroduction) {
	counter := atomic.AddUint64(&m.GetIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.GetIntroductionCounter, 1)

	if len(m.GetIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetIntroduction.")
			return
		}

		result := m.GetIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.GetIntroductionMock.mainExpectation != nil {

		result := m.GetIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetIntroduction")
		}

		r = result.r

		return
	}

	if m.GetIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetIntroduction.")
		return
	}

	return m.GetIntroductionFunc()
}

//GetIntroductionMinimockCounter returns a count of LocalNodeMock.GetIntroductionFunc invocations
func (m *LocalNodeMock) GetIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionCounter)
}

//GetIntroductionMinimockPreCounter returns the value of LocalNodeMock.GetIntroduction invocations
func (m *LocalNodeMock) GetIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroductionPreCounter)
}

//GetIntroductionFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetIntroductionFinished() bool {
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

type mLocalNodeMockGetNodePublicKey struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetNodePublicKeyExpectation
	expectationSeries []*LocalNodeMockGetNodePublicKeyExpectation
}

type LocalNodeMockGetNodePublicKeyExpectation struct {
	result *LocalNodeMockGetNodePublicKeyResult
}

type LocalNodeMockGetNodePublicKeyResult struct {
	r cryptkit.SignatureKeyHolder
}

//Expect specifies that invocation of LocalNode.GetNodePublicKey is expected from 1 to Infinity times
func (m *mLocalNodeMockGetNodePublicKey) Expect() *mLocalNodeMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetNodePublicKey
func (m *mLocalNodeMockGetNodePublicKey) Return(r cryptkit.SignatureKeyHolder) *LocalNodeMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetNodePublicKey is expected once
func (m *mLocalNodeMockGetNodePublicKey) ExpectOnce() *LocalNodeMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetNodePublicKeyExpectation) Return(r cryptkit.SignatureKeyHolder) {
	e.result = &LocalNodeMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of LocalNode.GetNodePublicKey method
func (m *mLocalNodeMockGetNodePublicKey) Set(f func() (r cryptkit.SignatureKeyHolder)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetNodePublicKey() (r cryptkit.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of LocalNodeMock.GetNodePublicKeyFunc invocations
func (m *LocalNodeMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of LocalNodeMock.GetNodePublicKey invocations
func (m *LocalNodeMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetNodePublicKeyFinished() bool {
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

type mLocalNodeMockGetOpMode struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetOpModeExpectation
	expectationSeries []*LocalNodeMockGetOpModeExpectation
}

type LocalNodeMockGetOpModeExpectation struct {
	result *LocalNodeMockGetOpModeResult
}

type LocalNodeMockGetOpModeResult struct {
	r member.OpMode
}

//Expect specifies that invocation of LocalNode.GetOpMode is expected from 1 to Infinity times
func (m *mLocalNodeMockGetOpMode) Expect() *mLocalNodeMockGetOpMode {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetOpModeExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetOpMode
func (m *mLocalNodeMockGetOpMode) Return(r member.OpMode) *LocalNodeMock {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetOpModeExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetOpModeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetOpMode is expected once
func (m *mLocalNodeMockGetOpMode) ExpectOnce() *LocalNodeMockGetOpModeExpectation {
	m.mock.GetOpModeFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetOpModeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetOpModeExpectation) Return(r member.OpMode) {
	e.result = &LocalNodeMockGetOpModeResult{r}
}

//Set uses given function f as a mock of LocalNode.GetOpMode method
func (m *mLocalNodeMockGetOpMode) Set(f func() (r member.OpMode)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOpModeFunc = f
	return m.mock
}

//GetOpMode implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetOpMode() (r member.OpMode) {
	counter := atomic.AddUint64(&m.GetOpModePreCounter, 1)
	defer atomic.AddUint64(&m.GetOpModeCounter, 1)

	if len(m.GetOpModeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOpModeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetOpMode.")
			return
		}

		result := m.GetOpModeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetOpMode")
			return
		}

		r = result.r

		return
	}

	if m.GetOpModeMock.mainExpectation != nil {

		result := m.GetOpModeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetOpMode")
		}

		r = result.r

		return
	}

	if m.GetOpModeFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetOpMode.")
		return
	}

	return m.GetOpModeFunc()
}

//GetOpModeMinimockCounter returns a count of LocalNodeMock.GetOpModeFunc invocations
func (m *LocalNodeMock) GetOpModeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModeCounter)
}

//GetOpModeMinimockPreCounter returns the value of LocalNodeMock.GetOpMode invocations
func (m *LocalNodeMock) GetOpModeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModePreCounter)
}

//GetOpModeFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetOpModeFinished() bool {
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

type mLocalNodeMockGetPrimaryRole struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetPrimaryRoleExpectation
	expectationSeries []*LocalNodeMockGetPrimaryRoleExpectation
}

type LocalNodeMockGetPrimaryRoleExpectation struct {
	result *LocalNodeMockGetPrimaryRoleResult
}

type LocalNodeMockGetPrimaryRoleResult struct {
	r member.PrimaryRole
}

//Expect specifies that invocation of LocalNode.GetPrimaryRole is expected from 1 to Infinity times
func (m *mLocalNodeMockGetPrimaryRole) Expect() *mLocalNodeMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetPrimaryRole
func (m *mLocalNodeMockGetPrimaryRole) Return(r member.PrimaryRole) *LocalNodeMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetPrimaryRole is expected once
func (m *mLocalNodeMockGetPrimaryRole) ExpectOnce() *LocalNodeMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetPrimaryRoleExpectation) Return(r member.PrimaryRole) {
	e.result = &LocalNodeMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of LocalNode.GetPrimaryRole method
func (m *mLocalNodeMockGetPrimaryRole) Set(f func() (r member.PrimaryRole)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetPrimaryRole() (r member.PrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of LocalNodeMock.GetPrimaryRoleFunc invocations
func (m *LocalNodeMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of LocalNodeMock.GetPrimaryRole invocations
func (m *LocalNodeMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetPrimaryRoleFinished() bool {
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

type mLocalNodeMockGetPublicKeyStore struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetPublicKeyStoreExpectation
	expectationSeries []*LocalNodeMockGetPublicKeyStoreExpectation
}

type LocalNodeMockGetPublicKeyStoreExpectation struct {
	result *LocalNodeMockGetPublicKeyStoreResult
}

type LocalNodeMockGetPublicKeyStoreResult struct {
	r cryptkit.PublicKeyStore
}

//Expect specifies that invocation of LocalNode.GetPublicKeyStore is expected from 1 to Infinity times
func (m *mLocalNodeMockGetPublicKeyStore) Expect() *mLocalNodeMockGetPublicKeyStore {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetPublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetPublicKeyStore
func (m *mLocalNodeMockGetPublicKeyStore) Return(r cryptkit.PublicKeyStore) *LocalNodeMock {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetPublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetPublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetPublicKeyStore is expected once
func (m *mLocalNodeMockGetPublicKeyStore) ExpectOnce() *LocalNodeMockGetPublicKeyStoreExpectation {
	m.mock.GetPublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetPublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetPublicKeyStoreExpectation) Return(r cryptkit.PublicKeyStore) {
	e.result = &LocalNodeMockGetPublicKeyStoreResult{r}
}

//Set uses given function f as a mock of LocalNode.GetPublicKeyStore method
func (m *mLocalNodeMockGetPublicKeyStore) Set(f func() (r cryptkit.PublicKeyStore)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyStoreFunc = f
	return m.mock
}

//GetPublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetPublicKeyStore() (r cryptkit.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetPublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyStoreCounter, 1)

	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetPublicKeyStore.")
			return
		}

		result := m.GetPublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetPublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreMock.mainExpectation != nil {

		result := m.GetPublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetPublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetPublicKeyStore.")
		return
	}

	return m.GetPublicKeyStoreFunc()
}

//GetPublicKeyStoreMinimockCounter returns a count of LocalNodeMock.GetPublicKeyStoreFunc invocations
func (m *LocalNodeMock) GetPublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStoreCounter)
}

//GetPublicKeyStoreMinimockPreCounter returns the value of LocalNodeMock.GetPublicKeyStore invocations
func (m *LocalNodeMock) GetPublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStorePreCounter)
}

//GetPublicKeyStoreFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetPublicKeyStoreFinished() bool {
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

type mLocalNodeMockGetShortNodeID struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetShortNodeIDExpectation
	expectationSeries []*LocalNodeMockGetShortNodeIDExpectation
}

type LocalNodeMockGetShortNodeIDExpectation struct {
	result *LocalNodeMockGetShortNodeIDResult
}

type LocalNodeMockGetShortNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of LocalNode.GetNodeID is expected from 1 to Infinity times
func (m *mLocalNodeMockGetShortNodeID) Expect() *mLocalNodeMockGetShortNodeID {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetShortNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetNodeID
func (m *mLocalNodeMockGetShortNodeID) Return(r insolar.ShortNodeID) *LocalNodeMock {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetShortNodeIDExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetShortNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetNodeID is expected once
func (m *mLocalNodeMockGetShortNodeID) ExpectOnce() *LocalNodeMockGetShortNodeIDExpectation {
	m.mock.GetShortNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetShortNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetShortNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &LocalNodeMockGetShortNodeIDResult{r}
}

//Set uses given function f as a mock of LocalNode.GetNodeID method
func (m *mLocalNodeMockGetShortNodeID) Set(f func() (r insolar.ShortNodeID)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetShortNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetShortNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetShortNodeIDCounter, 1)

	if len(m.GetShortNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetShortNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetNodeID.")
			return
		}

		result := m.GetShortNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDMock.mainExpectation != nil {

		result := m.GetShortNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetNodeID.")
		return
	}

	return m.GetShortNodeIDFunc()
}

//GetShortNodeIDMinimockCounter returns a count of LocalNodeMock.GetShortNodeIDFunc invocations
func (m *LocalNodeMock) GetShortNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDCounter)
}

//GetShortNodeIDMinimockPreCounter returns the value of LocalNodeMock.GetNodeID invocations
func (m *LocalNodeMock) GetShortNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDPreCounter)
}

//GetShortNodeIDFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetShortNodeIDFinished() bool {
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

type mLocalNodeMockGetSignatureVerifier struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetSignatureVerifierExpectation
	expectationSeries []*LocalNodeMockGetSignatureVerifierExpectation
}

type LocalNodeMockGetSignatureVerifierExpectation struct {
	result *LocalNodeMockGetSignatureVerifierResult
}

type LocalNodeMockGetSignatureVerifierResult struct {
	r cryptkit.SignatureVerifier
}

//Expect specifies that invocation of LocalNode.GetSignatureVerifier is expected from 1 to Infinity times
func (m *mLocalNodeMockGetSignatureVerifier) Expect() *mLocalNodeMockGetSignatureVerifier {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetSignatureVerifierExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetSignatureVerifier
func (m *mLocalNodeMockGetSignatureVerifier) Return(r cryptkit.SignatureVerifier) *LocalNodeMock {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetSignatureVerifierExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetSignatureVerifierResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetSignatureVerifier is expected once
func (m *mLocalNodeMockGetSignatureVerifier) ExpectOnce() *LocalNodeMockGetSignatureVerifierExpectation {
	m.mock.GetSignatureVerifierFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetSignatureVerifierExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetSignatureVerifierExpectation) Return(r cryptkit.SignatureVerifier) {
	e.result = &LocalNodeMockGetSignatureVerifierResult{r}
}

//Set uses given function f as a mock of LocalNode.GetSignatureVerifier method
func (m *mLocalNodeMockGetSignatureVerifier) Set(f func() (r cryptkit.SignatureVerifier)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureVerifierFunc = f
	return m.mock
}

//GetSignatureVerifier implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetSignatureVerifier() (r cryptkit.SignatureVerifier) {
	counter := atomic.AddUint64(&m.GetSignatureVerifierPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureVerifierCounter, 1)

	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureVerifierMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetSignatureVerifier.")
			return
		}

		result := m.GetSignatureVerifierMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetSignatureVerifier")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierMock.mainExpectation != nil {

		result := m.GetSignatureVerifierMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetSignatureVerifier")
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetSignatureVerifier.")
		return
	}

	return m.GetSignatureVerifierFunc()
}

//GetSignatureVerifierMinimockCounter returns a count of LocalNodeMock.GetSignatureVerifierFunc invocations
func (m *LocalNodeMock) GetSignatureVerifierMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierCounter)
}

//GetSignatureVerifierMinimockPreCounter returns the value of LocalNodeMock.GetSignatureVerifier invocations
func (m *LocalNodeMock) GetSignatureVerifierMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierPreCounter)
}

//GetSignatureVerifierFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetSignatureVerifierFinished() bool {
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

type mLocalNodeMockGetSpecialRoles struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetSpecialRolesExpectation
	expectationSeries []*LocalNodeMockGetSpecialRolesExpectation
}

type LocalNodeMockGetSpecialRolesExpectation struct {
	result *LocalNodeMockGetSpecialRolesResult
}

type LocalNodeMockGetSpecialRolesResult struct {
	r member.SpecialRole
}

//Expect specifies that invocation of LocalNode.GetSpecialRoles is expected from 1 to Infinity times
func (m *mLocalNodeMockGetSpecialRoles) Expect() *mLocalNodeMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetSpecialRoles
func (m *mLocalNodeMockGetSpecialRoles) Return(r member.SpecialRole) *LocalNodeMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetSpecialRoles is expected once
func (m *mLocalNodeMockGetSpecialRoles) ExpectOnce() *LocalNodeMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetSpecialRolesExpectation) Return(r member.SpecialRole) {
	e.result = &LocalNodeMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of LocalNode.GetSpecialRoles method
func (m *mLocalNodeMockGetSpecialRoles) Set(f func() (r member.SpecialRole)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetSpecialRoles() (r member.SpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of LocalNodeMock.GetSpecialRolesFunc invocations
func (m *LocalNodeMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of LocalNodeMock.GetSpecialRoles invocations
func (m *LocalNodeMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetSpecialRolesFinished() bool {
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

type mLocalNodeMockGetStartPower struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetStartPowerExpectation
	expectationSeries []*LocalNodeMockGetStartPowerExpectation
}

type LocalNodeMockGetStartPowerExpectation struct {
	result *LocalNodeMockGetStartPowerResult
}

type LocalNodeMockGetStartPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of LocalNode.GetStartPower is expected from 1 to Infinity times
func (m *mLocalNodeMockGetStartPower) Expect() *mLocalNodeMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetStartPower
func (m *mLocalNodeMockGetStartPower) Return(r member.Power) *LocalNodeMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetStartPower is expected once
func (m *mLocalNodeMockGetStartPower) ExpectOnce() *LocalNodeMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetStartPowerExpectation) Return(r member.Power) {
	e.result = &LocalNodeMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of LocalNode.GetStartPower method
func (m *mLocalNodeMockGetStartPower) Set(f func() (r member.Power)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetStartPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of LocalNodeMock.GetStartPowerFunc invocations
func (m *LocalNodeMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of LocalNodeMock.GetStartPower invocations
func (m *LocalNodeMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetStartPowerFinished() bool {
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

type mLocalNodeMockHasIntroduction struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockHasIntroductionExpectation
	expectationSeries []*LocalNodeMockHasIntroductionExpectation
}

type LocalNodeMockHasIntroductionExpectation struct {
	result *LocalNodeMockHasIntroductionResult
}

type LocalNodeMockHasIntroductionResult struct {
	r bool
}

//Expect specifies that invocation of LocalNode.HasIntroduction is expected from 1 to Infinity times
func (m *mLocalNodeMockHasIntroduction) Expect() *mLocalNodeMockHasIntroduction {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockHasIntroductionExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.HasIntroduction
func (m *mLocalNodeMockHasIntroduction) Return(r bool) *LocalNodeMock {
	m.mock.HasIntroductionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockHasIntroductionExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockHasIntroductionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.HasIntroduction is expected once
func (m *mLocalNodeMockHasIntroduction) ExpectOnce() *LocalNodeMockHasIntroductionExpectation {
	m.mock.HasIntroductionFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockHasIntroductionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockHasIntroductionExpectation) Return(r bool) {
	e.result = &LocalNodeMockHasIntroductionResult{r}
}

//Set uses given function f as a mock of LocalNode.HasIntroduction method
func (m *mLocalNodeMockHasIntroduction) Set(f func() (r bool)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasIntroductionFunc = f
	return m.mock
}

//HasIntroduction implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) HasIntroduction() (r bool) {
	counter := atomic.AddUint64(&m.HasIntroductionPreCounter, 1)
	defer atomic.AddUint64(&m.HasIntroductionCounter, 1)

	if len(m.HasIntroductionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasIntroductionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.HasIntroduction.")
			return
		}

		result := m.HasIntroductionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.HasIntroduction")
			return
		}

		r = result.r

		return
	}

	if m.HasIntroductionMock.mainExpectation != nil {

		result := m.HasIntroductionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.HasIntroduction")
		}

		r = result.r

		return
	}

	if m.HasIntroductionFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.HasIntroduction.")
		return
	}

	return m.HasIntroductionFunc()
}

//HasIntroductionMinimockCounter returns a count of LocalNodeMock.HasIntroductionFunc invocations
func (m *LocalNodeMock) HasIntroductionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionCounter)
}

//HasIntroductionMinimockPreCounter returns the value of LocalNodeMock.HasIntroduction invocations
func (m *LocalNodeMock) HasIntroductionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasIntroductionPreCounter)
}

//HasIntroductionFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) HasIntroductionFinished() bool {
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

type mLocalNodeMockIsAcceptableHost struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockIsAcceptableHostExpectation
	expectationSeries []*LocalNodeMockIsAcceptableHostExpectation
}

type LocalNodeMockIsAcceptableHostExpectation struct {
	input  *LocalNodeMockIsAcceptableHostInput
	result *LocalNodeMockIsAcceptableHostResult
}

type LocalNodeMockIsAcceptableHostInput struct {
	p endpoints.Inbound
}

type LocalNodeMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of LocalNode.IsAcceptableHost is expected from 1 to Infinity times
func (m *mLocalNodeMockIsAcceptableHost) Expect(p endpoints.Inbound) *mLocalNodeMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &LocalNodeMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of LocalNode.IsAcceptableHost
func (m *mLocalNodeMockIsAcceptableHost) Return(r bool) *LocalNodeMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.IsAcceptableHost is expected once
func (m *mLocalNodeMockIsAcceptableHost) ExpectOnce(p endpoints.Inbound) *LocalNodeMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockIsAcceptableHostExpectation{}
	expectation.input = &LocalNodeMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &LocalNodeMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of LocalNode.IsAcceptableHost method
func (m *mLocalNodeMockIsAcceptableHost) Set(f func(p endpoints.Inbound) (r bool)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) IsAcceptableHost(p endpoints.Inbound) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LocalNodeMockIsAcceptableHostInput{p}, "LocalNode.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LocalNodeMockIsAcceptableHostInput{p}, "LocalNode.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of LocalNodeMock.IsAcceptableHostFunc invocations
func (m *LocalNodeMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of LocalNodeMock.IsAcceptableHost invocations
func (m *LocalNodeMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) IsAcceptableHostFinished() bool {
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

type mLocalNodeMockIsJoiner struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockIsJoinerExpectation
	expectationSeries []*LocalNodeMockIsJoinerExpectation
}

type LocalNodeMockIsJoinerExpectation struct {
	result *LocalNodeMockIsJoinerResult
}

type LocalNodeMockIsJoinerResult struct {
	r bool
}

//Expect specifies that invocation of LocalNode.IsJoiner is expected from 1 to Infinity times
func (m *mLocalNodeMockIsJoiner) Expect() *mLocalNodeMockIsJoiner {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockIsJoinerExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.IsJoiner
func (m *mLocalNodeMockIsJoiner) Return(r bool) *LocalNodeMock {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockIsJoinerExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockIsJoinerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.IsJoiner is expected once
func (m *mLocalNodeMockIsJoiner) ExpectOnce() *LocalNodeMockIsJoinerExpectation {
	m.mock.IsJoinerFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockIsJoinerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockIsJoinerExpectation) Return(r bool) {
	e.result = &LocalNodeMockIsJoinerResult{r}
}

//Set uses given function f as a mock of LocalNode.IsJoiner method
func (m *mLocalNodeMockIsJoiner) Set(f func() (r bool)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsJoinerFunc = f
	return m.mock
}

//IsJoiner implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) IsJoiner() (r bool) {
	counter := atomic.AddUint64(&m.IsJoinerPreCounter, 1)
	defer atomic.AddUint64(&m.IsJoinerCounter, 1)

	if len(m.IsJoinerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsJoinerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.IsJoiner.")
			return
		}

		result := m.IsJoinerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.IsJoiner")
			return
		}

		r = result.r

		return
	}

	if m.IsJoinerMock.mainExpectation != nil {

		result := m.IsJoinerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.IsJoiner")
		}

		r = result.r

		return
	}

	if m.IsJoinerFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.IsJoiner.")
		return
	}

	return m.IsJoinerFunc()
}

//IsJoinerMinimockCounter returns a count of LocalNodeMock.IsJoinerFunc invocations
func (m *LocalNodeMock) IsJoinerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerCounter)
}

//IsJoinerMinimockPreCounter returns the value of LocalNodeMock.IsJoiner invocations
func (m *LocalNodeMock) IsJoinerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerPreCounter)
}

//IsJoinerFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) IsJoinerFinished() bool {
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

type mLocalNodeMockLocalNodeProfile struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockLocalNodeProfileExpectation
	expectationSeries []*LocalNodeMockLocalNodeProfileExpectation
}

type LocalNodeMockLocalNodeProfileExpectation struct {
}

//Expect specifies that invocation of LocalNode.LocalNodeProfile is expected from 1 to Infinity times
func (m *mLocalNodeMockLocalNodeProfile) Expect() *mLocalNodeMockLocalNodeProfile {
	m.mock.LocalNodeProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockLocalNodeProfileExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.LocalNodeProfile
func (m *mLocalNodeMockLocalNodeProfile) Return() *LocalNodeMock {
	m.mock.LocalNodeProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockLocalNodeProfileExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.LocalNodeProfile is expected once
func (m *mLocalNodeMockLocalNodeProfile) ExpectOnce() *LocalNodeMockLocalNodeProfileExpectation {
	m.mock.LocalNodeProfileFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockLocalNodeProfileExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of LocalNode.LocalNodeProfile method
func (m *mLocalNodeMockLocalNodeProfile) Set(f func()) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LocalNodeProfileFunc = f
	return m.mock
}

//LocalNodeProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) LocalNodeProfile() {
	counter := atomic.AddUint64(&m.LocalNodeProfilePreCounter, 1)
	defer atomic.AddUint64(&m.LocalNodeProfileCounter, 1)

	if len(m.LocalNodeProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LocalNodeProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.LocalNodeProfile.")
			return
		}

		return
	}

	if m.LocalNodeProfileMock.mainExpectation != nil {

		return
	}

	if m.LocalNodeProfileFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.LocalNodeProfile.")
		return
	}

	m.LocalNodeProfileFunc()
}

//LocalNodeProfileMinimockCounter returns a count of LocalNodeMock.LocalNodeProfileFunc invocations
func (m *LocalNodeMock) LocalNodeProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LocalNodeProfileCounter)
}

//LocalNodeProfileMinimockPreCounter returns the value of LocalNodeMock.LocalNodeProfile invocations
func (m *LocalNodeMock) LocalNodeProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LocalNodeProfilePreCounter)
}

//LocalNodeProfileFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) LocalNodeProfileFinished() bool {
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
func (m *LocalNodeMock) ValidateCallCounters() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetNodePublicKey")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetOpMode")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetPrimaryRole")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetPublicKeyStore")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.IsAcceptableHost")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.IsJoiner")
	}

	if !m.LocalNodeProfileFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.LocalNodeProfile")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LocalNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LocalNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LocalNodeMock) MinimockFinish() {

	if !m.GetAnnouncementSignatureFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetAnnouncementSignature")
	}

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetDeclaredPower")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetDefaultEndpoint")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetIndex")
	}

	if !m.GetIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetIntroduction")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetNodePublicKey")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetOpMode")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetPrimaryRole")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetPublicKeyStore")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetNodeID")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetSignatureVerifier")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetStartPower")
	}

	if !m.HasIntroductionFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.HasIntroduction")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.IsAcceptableHost")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.IsJoiner")
	}

	if !m.LocalNodeProfileFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.LocalNodeProfile")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LocalNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LocalNodeMock) MinimockWait(timeout time.Duration) {
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
		ok = ok && m.LocalNodeProfileFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAnnouncementSignatureFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetAnnouncementSignature")
			}

			if !m.GetDeclaredPowerFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetDeclaredPower")
			}

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetDefaultEndpoint")
			}

			if !m.GetIndexFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetIndex")
			}

			if !m.GetIntroductionFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetIntroduction")
			}

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetNodePublicKey")
			}

			if !m.GetOpModeFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetOpMode")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetPrimaryRole")
			}

			if !m.GetPublicKeyStoreFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetPublicKeyStore")
			}

			if !m.GetShortNodeIDFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetNodeID")
			}

			if !m.GetSignatureVerifierFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetSignatureVerifier")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetSpecialRoles")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetStartPower")
			}

			if !m.HasIntroductionFinished() {
				m.t.Error("Expected call to LocalNodeMock.HasIntroduction")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to LocalNodeMock.IsAcceptableHost")
			}

			if !m.IsJoinerFinished() {
				m.t.Error("Expected call to LocalNodeMock.IsJoiner")
			}

			if !m.LocalNodeProfileFinished() {
				m.t.Error("Expected call to LocalNodeMock.LocalNodeProfile")
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
func (m *LocalNodeMock) AllMocksCalled() bool {

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

	if !m.LocalNodeProfileFinished() {
		return false
	}

	return true
}

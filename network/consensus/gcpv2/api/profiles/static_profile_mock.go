package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StaticProfile" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
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

//StaticProfileMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile
type StaticProfileMock struct {
	t minimock.Tester

	GetBriefIntroSignedDigestFunc       func() (r cryptkit.SignedDigestHolder)
	GetBriefIntroSignedDigestCounter    uint64
	GetBriefIntroSignedDigestPreCounter uint64
	GetBriefIntroSignedDigestMock       mStaticProfileMockGetBriefIntroSignedDigest

	GetDefaultEndpointFunc       func() (r endpoints.Outbound)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mStaticProfileMockGetDefaultEndpoint

	GetExtensionFunc       func() (r StaticProfileExtension)
	GetExtensionCounter    uint64
	GetExtensionPreCounter uint64
	GetExtensionMock       mStaticProfileMockGetExtension

	GetNodePublicKeyFunc       func() (r cryptkit.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mStaticProfileMockGetNodePublicKey

	GetPrimaryRoleFunc       func() (r member.PrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mStaticProfileMockGetPrimaryRole

	GetPublicKeyStoreFunc       func() (r cryptkit.PublicKeyStore)
	GetPublicKeyStoreCounter    uint64
	GetPublicKeyStorePreCounter uint64
	GetPublicKeyStoreMock       mStaticProfileMockGetPublicKeyStore

	GetSpecialRolesFunc       func() (r member.SpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mStaticProfileMockGetSpecialRoles

	GetStartPowerFunc       func() (r member.Power)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mStaticProfileMockGetStartPower

	GetStaticNodeIDFunc       func() (r insolar.ShortNodeID)
	GetStaticNodeIDCounter    uint64
	GetStaticNodeIDPreCounter uint64
	GetStaticNodeIDMock       mStaticProfileMockGetStaticNodeID

	IsAcceptableHostFunc       func(p endpoints.Inbound) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mStaticProfileMockIsAcceptableHost
}

//NewStaticProfileMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile
func NewStaticProfileMock(t minimock.Tester) *StaticProfileMock {
	m := &StaticProfileMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetBriefIntroSignedDigestMock = mStaticProfileMockGetBriefIntroSignedDigest{mock: m}
	m.GetDefaultEndpointMock = mStaticProfileMockGetDefaultEndpoint{mock: m}
	m.GetExtensionMock = mStaticProfileMockGetExtension{mock: m}
	m.GetNodePublicKeyMock = mStaticProfileMockGetNodePublicKey{mock: m}
	m.GetPrimaryRoleMock = mStaticProfileMockGetPrimaryRole{mock: m}
	m.GetPublicKeyStoreMock = mStaticProfileMockGetPublicKeyStore{mock: m}
	m.GetSpecialRolesMock = mStaticProfileMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mStaticProfileMockGetStartPower{mock: m}
	m.GetStaticNodeIDMock = mStaticProfileMockGetStaticNodeID{mock: m}
	m.IsAcceptableHostMock = mStaticProfileMockIsAcceptableHost{mock: m}

	return m
}

type mStaticProfileMockGetBriefIntroSignedDigest struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetBriefIntroSignedDigestExpectation
	expectationSeries []*StaticProfileMockGetBriefIntroSignedDigestExpectation
}

type StaticProfileMockGetBriefIntroSignedDigestExpectation struct {
	result *StaticProfileMockGetBriefIntroSignedDigestResult
}

type StaticProfileMockGetBriefIntroSignedDigestResult struct {
	r cryptkit.SignedDigestHolder
}

//Expect specifies that invocation of StaticProfile.GetBriefIntroSignedDigest is expected from 1 to Infinity times
func (m *mStaticProfileMockGetBriefIntroSignedDigest) Expect() *mStaticProfileMockGetBriefIntroSignedDigest {
	m.mock.GetBriefIntroSignedDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetBriefIntroSignedDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetBriefIntroSignedDigest
func (m *mStaticProfileMockGetBriefIntroSignedDigest) Return(r cryptkit.SignedDigestHolder) *StaticProfileMock {
	m.mock.GetBriefIntroSignedDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetBriefIntroSignedDigestExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetBriefIntroSignedDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetBriefIntroSignedDigest is expected once
func (m *mStaticProfileMockGetBriefIntroSignedDigest) ExpectOnce() *StaticProfileMockGetBriefIntroSignedDigestExpectation {
	m.mock.GetBriefIntroSignedDigestFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetBriefIntroSignedDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetBriefIntroSignedDigestExpectation) Return(r cryptkit.SignedDigestHolder) {
	e.result = &StaticProfileMockGetBriefIntroSignedDigestResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetBriefIntroSignedDigest method
func (m *mStaticProfileMockGetBriefIntroSignedDigest) Set(f func() (r cryptkit.SignedDigestHolder)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetBriefIntroSignedDigestFunc = f
	return m.mock
}

//GetBriefIntroSignedDigest implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetBriefIntroSignedDigest() (r cryptkit.SignedDigestHolder) {
	counter := atomic.AddUint64(&m.GetBriefIntroSignedDigestPreCounter, 1)
	defer atomic.AddUint64(&m.GetBriefIntroSignedDigestCounter, 1)

	if len(m.GetBriefIntroSignedDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetBriefIntroSignedDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetBriefIntroSignedDigest.")
			return
		}

		result := m.GetBriefIntroSignedDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetBriefIntroSignedDigest")
			return
		}

		r = result.r

		return
	}

	if m.GetBriefIntroSignedDigestMock.mainExpectation != nil {

		result := m.GetBriefIntroSignedDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetBriefIntroSignedDigest")
		}

		r = result.r

		return
	}

	if m.GetBriefIntroSignedDigestFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetBriefIntroSignedDigest.")
		return
	}

	return m.GetBriefIntroSignedDigestFunc()
}

//GetBriefIntroSignedDigestMinimockCounter returns a count of StaticProfileMock.GetBriefIntroSignedDigestFunc invocations
func (m *StaticProfileMock) GetBriefIntroSignedDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefIntroSignedDigestCounter)
}

//GetBriefIntroSignedDigestMinimockPreCounter returns the value of StaticProfileMock.GetBriefIntroSignedDigest invocations
func (m *StaticProfileMock) GetBriefIntroSignedDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefIntroSignedDigestPreCounter)
}

//GetBriefIntroSignedDigestFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetBriefIntroSignedDigestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetBriefIntroSignedDigestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetBriefIntroSignedDigestCounter) == uint64(len(m.GetBriefIntroSignedDigestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetBriefIntroSignedDigestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetBriefIntroSignedDigestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetBriefIntroSignedDigestFunc != nil {
		return atomic.LoadUint64(&m.GetBriefIntroSignedDigestCounter) > 0
	}

	return true
}

type mStaticProfileMockGetDefaultEndpoint struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetDefaultEndpointExpectation
	expectationSeries []*StaticProfileMockGetDefaultEndpointExpectation
}

type StaticProfileMockGetDefaultEndpointExpectation struct {
	result *StaticProfileMockGetDefaultEndpointResult
}

type StaticProfileMockGetDefaultEndpointResult struct {
	r endpoints.Outbound
}

//Expect specifies that invocation of StaticProfile.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mStaticProfileMockGetDefaultEndpoint) Expect() *mStaticProfileMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetDefaultEndpoint
func (m *mStaticProfileMockGetDefaultEndpoint) Return(r endpoints.Outbound) *StaticProfileMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetDefaultEndpoint is expected once
func (m *mStaticProfileMockGetDefaultEndpoint) ExpectOnce() *StaticProfileMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetDefaultEndpointExpectation) Return(r endpoints.Outbound) {
	e.result = &StaticProfileMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetDefaultEndpoint method
func (m *mStaticProfileMockGetDefaultEndpoint) Set(f func() (r endpoints.Outbound)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetDefaultEndpoint() (r endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of StaticProfileMock.GetDefaultEndpointFunc invocations
func (m *StaticProfileMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of StaticProfileMock.GetDefaultEndpoint invocations
func (m *StaticProfileMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetDefaultEndpointFinished() bool {
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

type mStaticProfileMockGetExtension struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetExtensionExpectation
	expectationSeries []*StaticProfileMockGetExtensionExpectation
}

type StaticProfileMockGetExtensionExpectation struct {
	result *StaticProfileMockGetExtensionResult
}

type StaticProfileMockGetExtensionResult struct {
	r StaticProfileExtension
}

//Expect specifies that invocation of StaticProfile.GetExtension is expected from 1 to Infinity times
func (m *mStaticProfileMockGetExtension) Expect() *mStaticProfileMockGetExtension {
	m.mock.GetExtensionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetExtensionExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetExtension
func (m *mStaticProfileMockGetExtension) Return(r StaticProfileExtension) *StaticProfileMock {
	m.mock.GetExtensionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetExtensionExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetExtensionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetExtension is expected once
func (m *mStaticProfileMockGetExtension) ExpectOnce() *StaticProfileMockGetExtensionExpectation {
	m.mock.GetExtensionFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetExtensionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetExtensionExpectation) Return(r StaticProfileExtension) {
	e.result = &StaticProfileMockGetExtensionResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetExtension method
func (m *mStaticProfileMockGetExtension) Set(f func() (r StaticProfileExtension)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExtensionFunc = f
	return m.mock
}

//GetExtension implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetExtension() (r StaticProfileExtension) {
	counter := atomic.AddUint64(&m.GetExtensionPreCounter, 1)
	defer atomic.AddUint64(&m.GetExtensionCounter, 1)

	if len(m.GetExtensionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExtensionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetExtension.")
			return
		}

		result := m.GetExtensionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetExtension")
			return
		}

		r = result.r

		return
	}

	if m.GetExtensionMock.mainExpectation != nil {

		result := m.GetExtensionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetExtension")
		}

		r = result.r

		return
	}

	if m.GetExtensionFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetExtension.")
		return
	}

	return m.GetExtensionFunc()
}

//GetExtensionMinimockCounter returns a count of StaticProfileMock.GetExtensionFunc invocations
func (m *StaticProfileMock) GetExtensionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtensionCounter)
}

//GetExtensionMinimockPreCounter returns the value of StaticProfileMock.GetExtension invocations
func (m *StaticProfileMock) GetExtensionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtensionPreCounter)
}

//GetExtensionFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetExtensionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetExtensionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetExtensionCounter) == uint64(len(m.GetExtensionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetExtensionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetExtensionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetExtensionFunc != nil {
		return atomic.LoadUint64(&m.GetExtensionCounter) > 0
	}

	return true
}

type mStaticProfileMockGetNodePublicKey struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetNodePublicKeyExpectation
	expectationSeries []*StaticProfileMockGetNodePublicKeyExpectation
}

type StaticProfileMockGetNodePublicKeyExpectation struct {
	result *StaticProfileMockGetNodePublicKeyResult
}

type StaticProfileMockGetNodePublicKeyResult struct {
	r cryptkit.SignatureKeyHolder
}

//Expect specifies that invocation of StaticProfile.GetNodePublicKey is expected from 1 to Infinity times
func (m *mStaticProfileMockGetNodePublicKey) Expect() *mStaticProfileMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetNodePublicKey
func (m *mStaticProfileMockGetNodePublicKey) Return(r cryptkit.SignatureKeyHolder) *StaticProfileMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetNodePublicKey is expected once
func (m *mStaticProfileMockGetNodePublicKey) ExpectOnce() *StaticProfileMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetNodePublicKeyExpectation) Return(r cryptkit.SignatureKeyHolder) {
	e.result = &StaticProfileMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetNodePublicKey method
func (m *mStaticProfileMockGetNodePublicKey) Set(f func() (r cryptkit.SignatureKeyHolder)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetNodePublicKey() (r cryptkit.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of StaticProfileMock.GetNodePublicKeyFunc invocations
func (m *StaticProfileMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of StaticProfileMock.GetNodePublicKey invocations
func (m *StaticProfileMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetNodePublicKeyFinished() bool {
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

type mStaticProfileMockGetPrimaryRole struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetPrimaryRoleExpectation
	expectationSeries []*StaticProfileMockGetPrimaryRoleExpectation
}

type StaticProfileMockGetPrimaryRoleExpectation struct {
	result *StaticProfileMockGetPrimaryRoleResult
}

type StaticProfileMockGetPrimaryRoleResult struct {
	r member.PrimaryRole
}

//Expect specifies that invocation of StaticProfile.GetPrimaryRole is expected from 1 to Infinity times
func (m *mStaticProfileMockGetPrimaryRole) Expect() *mStaticProfileMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetPrimaryRole
func (m *mStaticProfileMockGetPrimaryRole) Return(r member.PrimaryRole) *StaticProfileMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetPrimaryRole is expected once
func (m *mStaticProfileMockGetPrimaryRole) ExpectOnce() *StaticProfileMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetPrimaryRoleExpectation) Return(r member.PrimaryRole) {
	e.result = &StaticProfileMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetPrimaryRole method
func (m *mStaticProfileMockGetPrimaryRole) Set(f func() (r member.PrimaryRole)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetPrimaryRole() (r member.PrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of StaticProfileMock.GetPrimaryRoleFunc invocations
func (m *StaticProfileMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of StaticProfileMock.GetPrimaryRole invocations
func (m *StaticProfileMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetPrimaryRoleFinished() bool {
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

type mStaticProfileMockGetPublicKeyStore struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetPublicKeyStoreExpectation
	expectationSeries []*StaticProfileMockGetPublicKeyStoreExpectation
}

type StaticProfileMockGetPublicKeyStoreExpectation struct {
	result *StaticProfileMockGetPublicKeyStoreResult
}

type StaticProfileMockGetPublicKeyStoreResult struct {
	r cryptkit.PublicKeyStore
}

//Expect specifies that invocation of StaticProfile.GetPublicKeyStore is expected from 1 to Infinity times
func (m *mStaticProfileMockGetPublicKeyStore) Expect() *mStaticProfileMockGetPublicKeyStore {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetPublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetPublicKeyStore
func (m *mStaticProfileMockGetPublicKeyStore) Return(r cryptkit.PublicKeyStore) *StaticProfileMock {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetPublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetPublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetPublicKeyStore is expected once
func (m *mStaticProfileMockGetPublicKeyStore) ExpectOnce() *StaticProfileMockGetPublicKeyStoreExpectation {
	m.mock.GetPublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetPublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetPublicKeyStoreExpectation) Return(r cryptkit.PublicKeyStore) {
	e.result = &StaticProfileMockGetPublicKeyStoreResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetPublicKeyStore method
func (m *mStaticProfileMockGetPublicKeyStore) Set(f func() (r cryptkit.PublicKeyStore)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyStoreFunc = f
	return m.mock
}

//GetPublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetPublicKeyStore() (r cryptkit.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetPublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyStoreCounter, 1)

	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetPublicKeyStore.")
			return
		}

		result := m.GetPublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetPublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreMock.mainExpectation != nil {

		result := m.GetPublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetPublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetPublicKeyStore.")
		return
	}

	return m.GetPublicKeyStoreFunc()
}

//GetPublicKeyStoreMinimockCounter returns a count of StaticProfileMock.GetPublicKeyStoreFunc invocations
func (m *StaticProfileMock) GetPublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStoreCounter)
}

//GetPublicKeyStoreMinimockPreCounter returns the value of StaticProfileMock.GetPublicKeyStore invocations
func (m *StaticProfileMock) GetPublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStorePreCounter)
}

//GetPublicKeyStoreFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetPublicKeyStoreFinished() bool {
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

type mStaticProfileMockGetSpecialRoles struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetSpecialRolesExpectation
	expectationSeries []*StaticProfileMockGetSpecialRolesExpectation
}

type StaticProfileMockGetSpecialRolesExpectation struct {
	result *StaticProfileMockGetSpecialRolesResult
}

type StaticProfileMockGetSpecialRolesResult struct {
	r member.SpecialRole
}

//Expect specifies that invocation of StaticProfile.GetSpecialRoles is expected from 1 to Infinity times
func (m *mStaticProfileMockGetSpecialRoles) Expect() *mStaticProfileMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetSpecialRoles
func (m *mStaticProfileMockGetSpecialRoles) Return(r member.SpecialRole) *StaticProfileMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetSpecialRoles is expected once
func (m *mStaticProfileMockGetSpecialRoles) ExpectOnce() *StaticProfileMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetSpecialRolesExpectation) Return(r member.SpecialRole) {
	e.result = &StaticProfileMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetSpecialRoles method
func (m *mStaticProfileMockGetSpecialRoles) Set(f func() (r member.SpecialRole)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetSpecialRoles() (r member.SpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of StaticProfileMock.GetSpecialRolesFunc invocations
func (m *StaticProfileMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of StaticProfileMock.GetSpecialRoles invocations
func (m *StaticProfileMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetSpecialRolesFinished() bool {
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

type mStaticProfileMockGetStartPower struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetStartPowerExpectation
	expectationSeries []*StaticProfileMockGetStartPowerExpectation
}

type StaticProfileMockGetStartPowerExpectation struct {
	result *StaticProfileMockGetStartPowerResult
}

type StaticProfileMockGetStartPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of StaticProfile.GetStartPower is expected from 1 to Infinity times
func (m *mStaticProfileMockGetStartPower) Expect() *mStaticProfileMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetStartPower
func (m *mStaticProfileMockGetStartPower) Return(r member.Power) *StaticProfileMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetStartPower is expected once
func (m *mStaticProfileMockGetStartPower) ExpectOnce() *StaticProfileMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetStartPowerExpectation) Return(r member.Power) {
	e.result = &StaticProfileMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetStartPower method
func (m *mStaticProfileMockGetStartPower) Set(f func() (r member.Power)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetStartPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of StaticProfileMock.GetStartPowerFunc invocations
func (m *StaticProfileMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of StaticProfileMock.GetStartPower invocations
func (m *StaticProfileMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetStartPowerFinished() bool {
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

type mStaticProfileMockGetStaticNodeID struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockGetStaticNodeIDExpectation
	expectationSeries []*StaticProfileMockGetStaticNodeIDExpectation
}

type StaticProfileMockGetStaticNodeIDExpectation struct {
	result *StaticProfileMockGetStaticNodeIDResult
}

type StaticProfileMockGetStaticNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of StaticProfile.GetStaticNodeID is expected from 1 to Infinity times
func (m *mStaticProfileMockGetStaticNodeID) Expect() *mStaticProfileMockGetStaticNodeID {
	m.mock.GetStaticNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetStaticNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfile.GetStaticNodeID
func (m *mStaticProfileMockGetStaticNodeID) Return(r insolar.ShortNodeID) *StaticProfileMock {
	m.mock.GetStaticNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockGetStaticNodeIDExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockGetStaticNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.GetStaticNodeID is expected once
func (m *mStaticProfileMockGetStaticNodeID) ExpectOnce() *StaticProfileMockGetStaticNodeIDExpectation {
	m.mock.GetStaticNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockGetStaticNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockGetStaticNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &StaticProfileMockGetStaticNodeIDResult{r}
}

//Set uses given function f as a mock of StaticProfile.GetStaticNodeID method
func (m *mStaticProfileMockGetStaticNodeID) Set(f func() (r insolar.ShortNodeID)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStaticNodeIDFunc = f
	return m.mock
}

//GetStaticNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) GetStaticNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetStaticNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetStaticNodeIDCounter, 1)

	if len(m.GetStaticNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStaticNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.GetStaticNodeID.")
			return
		}

		result := m.GetStaticNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetStaticNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetStaticNodeIDMock.mainExpectation != nil {

		result := m.GetStaticNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.GetStaticNodeID")
		}

		r = result.r

		return
	}

	if m.GetStaticNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.GetStaticNodeID.")
		return
	}

	return m.GetStaticNodeIDFunc()
}

//GetStaticNodeIDMinimockCounter returns a count of StaticProfileMock.GetStaticNodeIDFunc invocations
func (m *StaticProfileMock) GetStaticNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticNodeIDCounter)
}

//GetStaticNodeIDMinimockPreCounter returns the value of StaticProfileMock.GetStaticNodeID invocations
func (m *StaticProfileMock) GetStaticNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticNodeIDPreCounter)
}

//GetStaticNodeIDFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) GetStaticNodeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStaticNodeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStaticNodeIDCounter) == uint64(len(m.GetStaticNodeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStaticNodeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStaticNodeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStaticNodeIDFunc != nil {
		return atomic.LoadUint64(&m.GetStaticNodeIDCounter) > 0
	}

	return true
}

type mStaticProfileMockIsAcceptableHost struct {
	mock              *StaticProfileMock
	mainExpectation   *StaticProfileMockIsAcceptableHostExpectation
	expectationSeries []*StaticProfileMockIsAcceptableHostExpectation
}

type StaticProfileMockIsAcceptableHostExpectation struct {
	input  *StaticProfileMockIsAcceptableHostInput
	result *StaticProfileMockIsAcceptableHostResult
}

type StaticProfileMockIsAcceptableHostInput struct {
	p endpoints.Inbound
}

type StaticProfileMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of StaticProfile.IsAcceptableHost is expected from 1 to Infinity times
func (m *mStaticProfileMockIsAcceptableHost) Expect(p endpoints.Inbound) *mStaticProfileMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &StaticProfileMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of StaticProfile.IsAcceptableHost
func (m *mStaticProfileMockIsAcceptableHost) Return(r bool) *StaticProfileMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &StaticProfileMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfile.IsAcceptableHost is expected once
func (m *mStaticProfileMockIsAcceptableHost) ExpectOnce(p endpoints.Inbound) *StaticProfileMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileMockIsAcceptableHostExpectation{}
	expectation.input = &StaticProfileMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &StaticProfileMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of StaticProfile.IsAcceptableHost method
func (m *mStaticProfileMockIsAcceptableHost) Set(f func(p endpoints.Inbound) (r bool)) *StaticProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile interface
func (m *StaticProfileMock) IsAcceptableHost(p endpoints.Inbound) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StaticProfileMockIsAcceptableHostInput{p}, "StaticProfile.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StaticProfileMockIsAcceptableHostInput{p}, "StaticProfile.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of StaticProfileMock.IsAcceptableHostFunc invocations
func (m *StaticProfileMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of StaticProfileMock.IsAcceptableHost invocations
func (m *StaticProfileMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *StaticProfileMock) IsAcceptableHostFinished() bool {
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
func (m *StaticProfileMock) ValidateCallCounters() {

	if !m.GetBriefIntroSignedDigestFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetBriefIntroSignedDigest")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetDefaultEndpoint")
	}

	if !m.GetExtensionFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetExtension")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetNodePublicKey")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetPrimaryRole")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetPublicKeyStore")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetStartPower")
	}

	if !m.GetStaticNodeIDFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetStaticNodeID")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.IsAcceptableHost")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StaticProfileMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StaticProfileMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StaticProfileMock) MinimockFinish() {

	if !m.GetBriefIntroSignedDigestFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetBriefIntroSignedDigest")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetDefaultEndpoint")
	}

	if !m.GetExtensionFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetExtension")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetNodePublicKey")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetPrimaryRole")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetPublicKeyStore")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetStartPower")
	}

	if !m.GetStaticNodeIDFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.GetStaticNodeID")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to StaticProfileMock.IsAcceptableHost")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StaticProfileMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StaticProfileMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetBriefIntroSignedDigestFinished()
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetExtensionFinished()
		ok = ok && m.GetNodePublicKeyFinished()
		ok = ok && m.GetPrimaryRoleFinished()
		ok = ok && m.GetPublicKeyStoreFinished()
		ok = ok && m.GetSpecialRolesFinished()
		ok = ok && m.GetStartPowerFinished()
		ok = ok && m.GetStaticNodeIDFinished()
		ok = ok && m.IsAcceptableHostFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetBriefIntroSignedDigestFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetBriefIntroSignedDigest")
			}

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetDefaultEndpoint")
			}

			if !m.GetExtensionFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetExtension")
			}

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetNodePublicKey")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetPrimaryRole")
			}

			if !m.GetPublicKeyStoreFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetPublicKeyStore")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetSpecialRoles")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetStartPower")
			}

			if !m.GetStaticNodeIDFinished() {
				m.t.Error("Expected call to StaticProfileMock.GetStaticNodeID")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to StaticProfileMock.IsAcceptableHost")
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
func (m *StaticProfileMock) AllMocksCalled() bool {

	if !m.GetBriefIntroSignedDigestFinished() {
		return false
	}

	if !m.GetDefaultEndpointFinished() {
		return false
	}

	if !m.GetExtensionFinished() {
		return false
	}

	if !m.GetNodePublicKeyFinished() {
		return false
	}

	if !m.GetPrimaryRoleFinished() {
		return false
	}

	if !m.GetPublicKeyStoreFinished() {
		return false
	}

	if !m.GetSpecialRolesFinished() {
		return false
	}

	if !m.GetStartPowerFinished() {
		return false
	}

	if !m.GetStaticNodeIDFinished() {
		return false
	}

	if !m.IsAcceptableHostFinished() {
		return false
	}

	return true
}

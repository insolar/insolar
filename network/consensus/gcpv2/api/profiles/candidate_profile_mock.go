package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CandidateProfile" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	endpoints "github.com/insolar/insolar/network/consensus/common/endpoints"
	pulse "github.com/insolar/insolar/network/consensus/common/pulse"
	member "github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

//CandidateProfileMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile
type CandidateProfileMock struct {
	t minimock.Tester

	GetBriefIntroSignedDigestFunc       func() (r cryptkit.SignedDigestHolder)
	GetBriefIntroSignedDigestCounter    uint64
	GetBriefIntroSignedDigestPreCounter uint64
	GetBriefIntroSignedDigestMock       mCandidateProfileMockGetBriefIntroSignedDigest

	GetDefaultEndpointFunc       func() (r endpoints.Outbound)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mCandidateProfileMockGetDefaultEndpoint

	GetExtraEndpointsFunc       func() (r []endpoints.Outbound)
	GetExtraEndpointsCounter    uint64
	GetExtraEndpointsPreCounter uint64
	GetExtraEndpointsMock       mCandidateProfileMockGetExtraEndpoints

	GetIssuedAtPulseFunc       func() (r pulse.Number)
	GetIssuedAtPulseCounter    uint64
	GetIssuedAtPulsePreCounter uint64
	GetIssuedAtPulseMock       mCandidateProfileMockGetIssuedAtPulse

	GetIssuedAtTimeFunc       func() (r time.Time)
	GetIssuedAtTimeCounter    uint64
	GetIssuedAtTimePreCounter uint64
	GetIssuedAtTimeMock       mCandidateProfileMockGetIssuedAtTime

	GetIssuerIDFunc       func() (r insolar.ShortNodeID)
	GetIssuerIDCounter    uint64
	GetIssuerIDPreCounter uint64
	GetIssuerIDMock       mCandidateProfileMockGetIssuerID

	GetIssuerSignatureFunc       func() (r cryptkit.SignatureHolder)
	GetIssuerSignatureCounter    uint64
	GetIssuerSignaturePreCounter uint64
	GetIssuerSignatureMock       mCandidateProfileMockGetIssuerSignature

	GetNodePublicKeyFunc       func() (r cryptkit.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mCandidateProfileMockGetNodePublicKey

	GetPowerLevelsFunc       func() (r member.PowerSet)
	GetPowerLevelsCounter    uint64
	GetPowerLevelsPreCounter uint64
	GetPowerLevelsMock       mCandidateProfileMockGetPowerLevels

	GetPrimaryRoleFunc       func() (r member.PrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mCandidateProfileMockGetPrimaryRole

	GetReferenceFunc       func() (r insolar.Reference)
	GetReferenceCounter    uint64
	GetReferencePreCounter uint64
	GetReferenceMock       mCandidateProfileMockGetReference

	GetSpecialRolesFunc       func() (r member.SpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mCandidateProfileMockGetSpecialRoles

	GetStartPowerFunc       func() (r member.Power)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mCandidateProfileMockGetStartPower

	GetStaticNodeIDFunc       func() (r insolar.ShortNodeID)
	GetStaticNodeIDCounter    uint64
	GetStaticNodeIDPreCounter uint64
	GetStaticNodeIDMock       mCandidateProfileMockGetStaticNodeID
}

//NewCandidateProfileMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile
func NewCandidateProfileMock(t minimock.Tester) *CandidateProfileMock {
	m := &CandidateProfileMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetBriefIntroSignedDigestMock = mCandidateProfileMockGetBriefIntroSignedDigest{mock: m}
	m.GetDefaultEndpointMock = mCandidateProfileMockGetDefaultEndpoint{mock: m}
	m.GetExtraEndpointsMock = mCandidateProfileMockGetExtraEndpoints{mock: m}
	m.GetIssuedAtPulseMock = mCandidateProfileMockGetIssuedAtPulse{mock: m}
	m.GetIssuedAtTimeMock = mCandidateProfileMockGetIssuedAtTime{mock: m}
	m.GetIssuerIDMock = mCandidateProfileMockGetIssuerID{mock: m}
	m.GetIssuerSignatureMock = mCandidateProfileMockGetIssuerSignature{mock: m}
	m.GetNodePublicKeyMock = mCandidateProfileMockGetNodePublicKey{mock: m}
	m.GetPowerLevelsMock = mCandidateProfileMockGetPowerLevels{mock: m}
	m.GetPrimaryRoleMock = mCandidateProfileMockGetPrimaryRole{mock: m}
	m.GetReferenceMock = mCandidateProfileMockGetReference{mock: m}
	m.GetSpecialRolesMock = mCandidateProfileMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mCandidateProfileMockGetStartPower{mock: m}
	m.GetStaticNodeIDMock = mCandidateProfileMockGetStaticNodeID{mock: m}

	return m
}

type mCandidateProfileMockGetBriefIntroSignedDigest struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetBriefIntroSignedDigestExpectation
	expectationSeries []*CandidateProfileMockGetBriefIntroSignedDigestExpectation
}

type CandidateProfileMockGetBriefIntroSignedDigestExpectation struct {
	result *CandidateProfileMockGetBriefIntroSignedDigestResult
}

type CandidateProfileMockGetBriefIntroSignedDigestResult struct {
	r cryptkit.SignedDigestHolder
}

//Expect specifies that invocation of CandidateProfile.GetBriefIntroSignedDigest is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetBriefIntroSignedDigest) Expect() *mCandidateProfileMockGetBriefIntroSignedDigest {
	m.mock.GetBriefIntroSignedDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetBriefIntroSignedDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetBriefIntroSignedDigest
func (m *mCandidateProfileMockGetBriefIntroSignedDigest) Return(r cryptkit.SignedDigestHolder) *CandidateProfileMock {
	m.mock.GetBriefIntroSignedDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetBriefIntroSignedDigestExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetBriefIntroSignedDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetBriefIntroSignedDigest is expected once
func (m *mCandidateProfileMockGetBriefIntroSignedDigest) ExpectOnce() *CandidateProfileMockGetBriefIntroSignedDigestExpectation {
	m.mock.GetBriefIntroSignedDigestFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetBriefIntroSignedDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetBriefIntroSignedDigestExpectation) Return(r cryptkit.SignedDigestHolder) {
	e.result = &CandidateProfileMockGetBriefIntroSignedDigestResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetBriefIntroSignedDigest method
func (m *mCandidateProfileMockGetBriefIntroSignedDigest) Set(f func() (r cryptkit.SignedDigestHolder)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetBriefIntroSignedDigestFunc = f
	return m.mock
}

//GetBriefIntroSignedDigest implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetBriefIntroSignedDigest() (r cryptkit.SignedDigestHolder) {
	counter := atomic.AddUint64(&m.GetBriefIntroSignedDigestPreCounter, 1)
	defer atomic.AddUint64(&m.GetBriefIntroSignedDigestCounter, 1)

	if len(m.GetBriefIntroSignedDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetBriefIntroSignedDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetBriefIntroSignedDigest.")
			return
		}

		result := m.GetBriefIntroSignedDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetBriefIntroSignedDigest")
			return
		}

		r = result.r

		return
	}

	if m.GetBriefIntroSignedDigestMock.mainExpectation != nil {

		result := m.GetBriefIntroSignedDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetBriefIntroSignedDigest")
		}

		r = result.r

		return
	}

	if m.GetBriefIntroSignedDigestFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetBriefIntroSignedDigest.")
		return
	}

	return m.GetBriefIntroSignedDigestFunc()
}

//GetBriefIntroSignedDigestMinimockCounter returns a count of CandidateProfileMock.GetBriefIntroSignedDigestFunc invocations
func (m *CandidateProfileMock) GetBriefIntroSignedDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefIntroSignedDigestCounter)
}

//GetBriefIntroSignedDigestMinimockPreCounter returns the value of CandidateProfileMock.GetBriefIntroSignedDigest invocations
func (m *CandidateProfileMock) GetBriefIntroSignedDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefIntroSignedDigestPreCounter)
}

//GetBriefIntroSignedDigestFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetBriefIntroSignedDigestFinished() bool {
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

type mCandidateProfileMockGetDefaultEndpoint struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetDefaultEndpointExpectation
	expectationSeries []*CandidateProfileMockGetDefaultEndpointExpectation
}

type CandidateProfileMockGetDefaultEndpointExpectation struct {
	result *CandidateProfileMockGetDefaultEndpointResult
}

type CandidateProfileMockGetDefaultEndpointResult struct {
	r endpoints.Outbound
}

//Expect specifies that invocation of CandidateProfile.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetDefaultEndpoint) Expect() *mCandidateProfileMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetDefaultEndpoint
func (m *mCandidateProfileMockGetDefaultEndpoint) Return(r endpoints.Outbound) *CandidateProfileMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetDefaultEndpoint is expected once
func (m *mCandidateProfileMockGetDefaultEndpoint) ExpectOnce() *CandidateProfileMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetDefaultEndpointExpectation) Return(r endpoints.Outbound) {
	e.result = &CandidateProfileMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetDefaultEndpoint method
func (m *mCandidateProfileMockGetDefaultEndpoint) Set(f func() (r endpoints.Outbound)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetDefaultEndpoint() (r endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of CandidateProfileMock.GetDefaultEndpointFunc invocations
func (m *CandidateProfileMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of CandidateProfileMock.GetDefaultEndpoint invocations
func (m *CandidateProfileMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetDefaultEndpointFinished() bool {
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

type mCandidateProfileMockGetExtraEndpoints struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetExtraEndpointsExpectation
	expectationSeries []*CandidateProfileMockGetExtraEndpointsExpectation
}

type CandidateProfileMockGetExtraEndpointsExpectation struct {
	result *CandidateProfileMockGetExtraEndpointsResult
}

type CandidateProfileMockGetExtraEndpointsResult struct {
	r []endpoints.Outbound
}

//Expect specifies that invocation of CandidateProfile.GetExtraEndpoints is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetExtraEndpoints) Expect() *mCandidateProfileMockGetExtraEndpoints {
	m.mock.GetExtraEndpointsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetExtraEndpointsExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetExtraEndpoints
func (m *mCandidateProfileMockGetExtraEndpoints) Return(r []endpoints.Outbound) *CandidateProfileMock {
	m.mock.GetExtraEndpointsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetExtraEndpointsExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetExtraEndpointsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetExtraEndpoints is expected once
func (m *mCandidateProfileMockGetExtraEndpoints) ExpectOnce() *CandidateProfileMockGetExtraEndpointsExpectation {
	m.mock.GetExtraEndpointsFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetExtraEndpointsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetExtraEndpointsExpectation) Return(r []endpoints.Outbound) {
	e.result = &CandidateProfileMockGetExtraEndpointsResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetExtraEndpoints method
func (m *mCandidateProfileMockGetExtraEndpoints) Set(f func() (r []endpoints.Outbound)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExtraEndpointsFunc = f
	return m.mock
}

//GetExtraEndpoints implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetExtraEndpoints() (r []endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetExtraEndpointsPreCounter, 1)
	defer atomic.AddUint64(&m.GetExtraEndpointsCounter, 1)

	if len(m.GetExtraEndpointsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExtraEndpointsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetExtraEndpoints.")
			return
		}

		result := m.GetExtraEndpointsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetExtraEndpoints")
			return
		}

		r = result.r

		return
	}

	if m.GetExtraEndpointsMock.mainExpectation != nil {

		result := m.GetExtraEndpointsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetExtraEndpoints")
		}

		r = result.r

		return
	}

	if m.GetExtraEndpointsFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetExtraEndpoints.")
		return
	}

	return m.GetExtraEndpointsFunc()
}

//GetExtraEndpointsMinimockCounter returns a count of CandidateProfileMock.GetExtraEndpointsFunc invocations
func (m *CandidateProfileMock) GetExtraEndpointsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtraEndpointsCounter)
}

//GetExtraEndpointsMinimockPreCounter returns the value of CandidateProfileMock.GetExtraEndpoints invocations
func (m *CandidateProfileMock) GetExtraEndpointsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtraEndpointsPreCounter)
}

//GetExtraEndpointsFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetExtraEndpointsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetExtraEndpointsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetExtraEndpointsCounter) == uint64(len(m.GetExtraEndpointsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetExtraEndpointsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetExtraEndpointsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetExtraEndpointsFunc != nil {
		return atomic.LoadUint64(&m.GetExtraEndpointsCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetIssuedAtPulse struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetIssuedAtPulseExpectation
	expectationSeries []*CandidateProfileMockGetIssuedAtPulseExpectation
}

type CandidateProfileMockGetIssuedAtPulseExpectation struct {
	result *CandidateProfileMockGetIssuedAtPulseResult
}

type CandidateProfileMockGetIssuedAtPulseResult struct {
	r pulse.Number
}

//Expect specifies that invocation of CandidateProfile.GetIssuedAtPulse is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetIssuedAtPulse) Expect() *mCandidateProfileMockGetIssuedAtPulse {
	m.mock.GetIssuedAtPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuedAtPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetIssuedAtPulse
func (m *mCandidateProfileMockGetIssuedAtPulse) Return(r pulse.Number) *CandidateProfileMock {
	m.mock.GetIssuedAtPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuedAtPulseExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetIssuedAtPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetIssuedAtPulse is expected once
func (m *mCandidateProfileMockGetIssuedAtPulse) ExpectOnce() *CandidateProfileMockGetIssuedAtPulseExpectation {
	m.mock.GetIssuedAtPulseFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetIssuedAtPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetIssuedAtPulseExpectation) Return(r pulse.Number) {
	e.result = &CandidateProfileMockGetIssuedAtPulseResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetIssuedAtPulse method
func (m *mCandidateProfileMockGetIssuedAtPulse) Set(f func() (r pulse.Number)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtPulseFunc = f
	return m.mock
}

//GetIssuedAtPulse implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetIssuedAtPulse() (r pulse.Number) {
	counter := atomic.AddUint64(&m.GetIssuedAtPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuedAtPulseCounter, 1)

	if len(m.GetIssuedAtPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuedAtPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuedAtPulse.")
			return
		}

		result := m.GetIssuedAtPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuedAtPulse")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuedAtPulseMock.mainExpectation != nil {

		result := m.GetIssuedAtPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuedAtPulse")
		}

		r = result.r

		return
	}

	if m.GetIssuedAtPulseFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuedAtPulse.")
		return
	}

	return m.GetIssuedAtPulseFunc()
}

//GetIssuedAtPulseMinimockCounter returns a count of CandidateProfileMock.GetIssuedAtPulseFunc invocations
func (m *CandidateProfileMock) GetIssuedAtPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtPulseCounter)
}

//GetIssuedAtPulseMinimockPreCounter returns the value of CandidateProfileMock.GetIssuedAtPulse invocations
func (m *CandidateProfileMock) GetIssuedAtPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtPulsePreCounter)
}

//GetIssuedAtPulseFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetIssuedAtPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIssuedAtPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIssuedAtPulseCounter) == uint64(len(m.GetIssuedAtPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIssuedAtPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIssuedAtPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIssuedAtPulseFunc != nil {
		return atomic.LoadUint64(&m.GetIssuedAtPulseCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetIssuedAtTime struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetIssuedAtTimeExpectation
	expectationSeries []*CandidateProfileMockGetIssuedAtTimeExpectation
}

type CandidateProfileMockGetIssuedAtTimeExpectation struct {
	result *CandidateProfileMockGetIssuedAtTimeResult
}

type CandidateProfileMockGetIssuedAtTimeResult struct {
	r time.Time
}

//Expect specifies that invocation of CandidateProfile.GetIssuedAtTime is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetIssuedAtTime) Expect() *mCandidateProfileMockGetIssuedAtTime {
	m.mock.GetIssuedAtTimeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuedAtTimeExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetIssuedAtTime
func (m *mCandidateProfileMockGetIssuedAtTime) Return(r time.Time) *CandidateProfileMock {
	m.mock.GetIssuedAtTimeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuedAtTimeExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetIssuedAtTimeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetIssuedAtTime is expected once
func (m *mCandidateProfileMockGetIssuedAtTime) ExpectOnce() *CandidateProfileMockGetIssuedAtTimeExpectation {
	m.mock.GetIssuedAtTimeFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetIssuedAtTimeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetIssuedAtTimeExpectation) Return(r time.Time) {
	e.result = &CandidateProfileMockGetIssuedAtTimeResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetIssuedAtTime method
func (m *mCandidateProfileMockGetIssuedAtTime) Set(f func() (r time.Time)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtTimeFunc = f
	return m.mock
}

//GetIssuedAtTime implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetIssuedAtTime() (r time.Time) {
	counter := atomic.AddUint64(&m.GetIssuedAtTimePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuedAtTimeCounter, 1)

	if len(m.GetIssuedAtTimeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuedAtTimeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuedAtTime.")
			return
		}

		result := m.GetIssuedAtTimeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuedAtTime")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuedAtTimeMock.mainExpectation != nil {

		result := m.GetIssuedAtTimeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuedAtTime")
		}

		r = result.r

		return
	}

	if m.GetIssuedAtTimeFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuedAtTime.")
		return
	}

	return m.GetIssuedAtTimeFunc()
}

//GetIssuedAtTimeMinimockCounter returns a count of CandidateProfileMock.GetIssuedAtTimeFunc invocations
func (m *CandidateProfileMock) GetIssuedAtTimeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtTimeCounter)
}

//GetIssuedAtTimeMinimockPreCounter returns the value of CandidateProfileMock.GetIssuedAtTime invocations
func (m *CandidateProfileMock) GetIssuedAtTimeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtTimePreCounter)
}

//GetIssuedAtTimeFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetIssuedAtTimeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIssuedAtTimeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIssuedAtTimeCounter) == uint64(len(m.GetIssuedAtTimeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIssuedAtTimeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIssuedAtTimeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIssuedAtTimeFunc != nil {
		return atomic.LoadUint64(&m.GetIssuedAtTimeCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetIssuerID struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetIssuerIDExpectation
	expectationSeries []*CandidateProfileMockGetIssuerIDExpectation
}

type CandidateProfileMockGetIssuerIDExpectation struct {
	result *CandidateProfileMockGetIssuerIDResult
}

type CandidateProfileMockGetIssuerIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of CandidateProfile.GetIssuerID is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetIssuerID) Expect() *mCandidateProfileMockGetIssuerID {
	m.mock.GetIssuerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetIssuerID
func (m *mCandidateProfileMockGetIssuerID) Return(r insolar.ShortNodeID) *CandidateProfileMock {
	m.mock.GetIssuerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuerIDExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetIssuerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetIssuerID is expected once
func (m *mCandidateProfileMockGetIssuerID) ExpectOnce() *CandidateProfileMockGetIssuerIDExpectation {
	m.mock.GetIssuerIDFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetIssuerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetIssuerIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &CandidateProfileMockGetIssuerIDResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetIssuerID method
func (m *mCandidateProfileMockGetIssuerID) Set(f func() (r insolar.ShortNodeID)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerIDFunc = f
	return m.mock
}

//GetIssuerID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetIssuerID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetIssuerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuerIDCounter, 1)

	if len(m.GetIssuerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuerID.")
			return
		}

		result := m.GetIssuerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuerID")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuerIDMock.mainExpectation != nil {

		result := m.GetIssuerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuerID")
		}

		r = result.r

		return
	}

	if m.GetIssuerIDFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuerID.")
		return
	}

	return m.GetIssuerIDFunc()
}

//GetIssuerIDMinimockCounter returns a count of CandidateProfileMock.GetIssuerIDFunc invocations
func (m *CandidateProfileMock) GetIssuerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerIDCounter)
}

//GetIssuerIDMinimockPreCounter returns the value of CandidateProfileMock.GetIssuerID invocations
func (m *CandidateProfileMock) GetIssuerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerIDPreCounter)
}

//GetIssuerIDFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetIssuerIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIssuerIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIssuerIDCounter) == uint64(len(m.GetIssuerIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIssuerIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIssuerIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIssuerIDFunc != nil {
		return atomic.LoadUint64(&m.GetIssuerIDCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetIssuerSignature struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetIssuerSignatureExpectation
	expectationSeries []*CandidateProfileMockGetIssuerSignatureExpectation
}

type CandidateProfileMockGetIssuerSignatureExpectation struct {
	result *CandidateProfileMockGetIssuerSignatureResult
}

type CandidateProfileMockGetIssuerSignatureResult struct {
	r cryptkit.SignatureHolder
}

//Expect specifies that invocation of CandidateProfile.GetIssuerSignature is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetIssuerSignature) Expect() *mCandidateProfileMockGetIssuerSignature {
	m.mock.GetIssuerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuerSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetIssuerSignature
func (m *mCandidateProfileMockGetIssuerSignature) Return(r cryptkit.SignatureHolder) *CandidateProfileMock {
	m.mock.GetIssuerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetIssuerSignatureExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetIssuerSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetIssuerSignature is expected once
func (m *mCandidateProfileMockGetIssuerSignature) ExpectOnce() *CandidateProfileMockGetIssuerSignatureExpectation {
	m.mock.GetIssuerSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetIssuerSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetIssuerSignatureExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &CandidateProfileMockGetIssuerSignatureResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetIssuerSignature method
func (m *mCandidateProfileMockGetIssuerSignature) Set(f func() (r cryptkit.SignatureHolder)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerSignatureFunc = f
	return m.mock
}

//GetIssuerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetIssuerSignature() (r cryptkit.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetIssuerSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuerSignatureCounter, 1)

	if len(m.GetIssuerSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuerSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuerSignature.")
			return
		}

		result := m.GetIssuerSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuerSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuerSignatureMock.mainExpectation != nil {

		result := m.GetIssuerSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetIssuerSignature")
		}

		r = result.r

		return
	}

	if m.GetIssuerSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetIssuerSignature.")
		return
	}

	return m.GetIssuerSignatureFunc()
}

//GetIssuerSignatureMinimockCounter returns a count of CandidateProfileMock.GetIssuerSignatureFunc invocations
func (m *CandidateProfileMock) GetIssuerSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerSignatureCounter)
}

//GetIssuerSignatureMinimockPreCounter returns the value of CandidateProfileMock.GetIssuerSignature invocations
func (m *CandidateProfileMock) GetIssuerSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerSignaturePreCounter)
}

//GetIssuerSignatureFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetIssuerSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIssuerSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIssuerSignatureCounter) == uint64(len(m.GetIssuerSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIssuerSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIssuerSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIssuerSignatureFunc != nil {
		return atomic.LoadUint64(&m.GetIssuerSignatureCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetNodePublicKey struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetNodePublicKeyExpectation
	expectationSeries []*CandidateProfileMockGetNodePublicKeyExpectation
}

type CandidateProfileMockGetNodePublicKeyExpectation struct {
	result *CandidateProfileMockGetNodePublicKeyResult
}

type CandidateProfileMockGetNodePublicKeyResult struct {
	r cryptkit.SignatureKeyHolder
}

//Expect specifies that invocation of CandidateProfile.GetNodePublicKey is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetNodePublicKey) Expect() *mCandidateProfileMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetNodePublicKey
func (m *mCandidateProfileMockGetNodePublicKey) Return(r cryptkit.SignatureKeyHolder) *CandidateProfileMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetNodePublicKey is expected once
func (m *mCandidateProfileMockGetNodePublicKey) ExpectOnce() *CandidateProfileMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetNodePublicKeyExpectation) Return(r cryptkit.SignatureKeyHolder) {
	e.result = &CandidateProfileMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetNodePublicKey method
func (m *mCandidateProfileMockGetNodePublicKey) Set(f func() (r cryptkit.SignatureKeyHolder)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetNodePublicKey() (r cryptkit.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of CandidateProfileMock.GetNodePublicKeyFunc invocations
func (m *CandidateProfileMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of CandidateProfileMock.GetNodePublicKey invocations
func (m *CandidateProfileMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetNodePublicKeyFinished() bool {
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

type mCandidateProfileMockGetPowerLevels struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetPowerLevelsExpectation
	expectationSeries []*CandidateProfileMockGetPowerLevelsExpectation
}

type CandidateProfileMockGetPowerLevelsExpectation struct {
	result *CandidateProfileMockGetPowerLevelsResult
}

type CandidateProfileMockGetPowerLevelsResult struct {
	r member.PowerSet
}

//Expect specifies that invocation of CandidateProfile.GetPowerLevels is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetPowerLevels) Expect() *mCandidateProfileMockGetPowerLevels {
	m.mock.GetPowerLevelsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetPowerLevelsExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetPowerLevels
func (m *mCandidateProfileMockGetPowerLevels) Return(r member.PowerSet) *CandidateProfileMock {
	m.mock.GetPowerLevelsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetPowerLevelsExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetPowerLevelsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetPowerLevels is expected once
func (m *mCandidateProfileMockGetPowerLevels) ExpectOnce() *CandidateProfileMockGetPowerLevelsExpectation {
	m.mock.GetPowerLevelsFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetPowerLevelsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetPowerLevelsExpectation) Return(r member.PowerSet) {
	e.result = &CandidateProfileMockGetPowerLevelsResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetPowerLevels method
func (m *mCandidateProfileMockGetPowerLevels) Set(f func() (r member.PowerSet)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPowerLevelsFunc = f
	return m.mock
}

//GetPowerLevels implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetPowerLevels() (r member.PowerSet) {
	counter := atomic.AddUint64(&m.GetPowerLevelsPreCounter, 1)
	defer atomic.AddUint64(&m.GetPowerLevelsCounter, 1)

	if len(m.GetPowerLevelsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPowerLevelsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetPowerLevels.")
			return
		}

		result := m.GetPowerLevelsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetPowerLevels")
			return
		}

		r = result.r

		return
	}

	if m.GetPowerLevelsMock.mainExpectation != nil {

		result := m.GetPowerLevelsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetPowerLevels")
		}

		r = result.r

		return
	}

	if m.GetPowerLevelsFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetPowerLevels.")
		return
	}

	return m.GetPowerLevelsFunc()
}

//GetPowerLevelsMinimockCounter returns a count of CandidateProfileMock.GetPowerLevelsFunc invocations
func (m *CandidateProfileMock) GetPowerLevelsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPowerLevelsCounter)
}

//GetPowerLevelsMinimockPreCounter returns the value of CandidateProfileMock.GetPowerLevels invocations
func (m *CandidateProfileMock) GetPowerLevelsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPowerLevelsPreCounter)
}

//GetPowerLevelsFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetPowerLevelsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPowerLevelsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPowerLevelsCounter) == uint64(len(m.GetPowerLevelsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPowerLevelsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPowerLevelsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPowerLevelsFunc != nil {
		return atomic.LoadUint64(&m.GetPowerLevelsCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetPrimaryRole struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetPrimaryRoleExpectation
	expectationSeries []*CandidateProfileMockGetPrimaryRoleExpectation
}

type CandidateProfileMockGetPrimaryRoleExpectation struct {
	result *CandidateProfileMockGetPrimaryRoleResult
}

type CandidateProfileMockGetPrimaryRoleResult struct {
	r member.PrimaryRole
}

//Expect specifies that invocation of CandidateProfile.GetPrimaryRole is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetPrimaryRole) Expect() *mCandidateProfileMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetPrimaryRole
func (m *mCandidateProfileMockGetPrimaryRole) Return(r member.PrimaryRole) *CandidateProfileMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetPrimaryRole is expected once
func (m *mCandidateProfileMockGetPrimaryRole) ExpectOnce() *CandidateProfileMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetPrimaryRoleExpectation) Return(r member.PrimaryRole) {
	e.result = &CandidateProfileMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetPrimaryRole method
func (m *mCandidateProfileMockGetPrimaryRole) Set(f func() (r member.PrimaryRole)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetPrimaryRole() (r member.PrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of CandidateProfileMock.GetPrimaryRoleFunc invocations
func (m *CandidateProfileMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of CandidateProfileMock.GetPrimaryRole invocations
func (m *CandidateProfileMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetPrimaryRoleFinished() bool {
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

type mCandidateProfileMockGetReference struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetReferenceExpectation
	expectationSeries []*CandidateProfileMockGetReferenceExpectation
}

type CandidateProfileMockGetReferenceExpectation struct {
	result *CandidateProfileMockGetReferenceResult
}

type CandidateProfileMockGetReferenceResult struct {
	r insolar.Reference
}

//Expect specifies that invocation of CandidateProfile.GetReference is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetReference) Expect() *mCandidateProfileMockGetReference {
	m.mock.GetReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetReferenceExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetReference
func (m *mCandidateProfileMockGetReference) Return(r insolar.Reference) *CandidateProfileMock {
	m.mock.GetReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetReferenceExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetReferenceResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetReference is expected once
func (m *mCandidateProfileMockGetReference) ExpectOnce() *CandidateProfileMockGetReferenceExpectation {
	m.mock.GetReferenceFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetReferenceExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetReferenceExpectation) Return(r insolar.Reference) {
	e.result = &CandidateProfileMockGetReferenceResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetReference method
func (m *mCandidateProfileMockGetReference) Set(f func() (r insolar.Reference)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetReferenceFunc = f
	return m.mock
}

//GetReference implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetReference() (r insolar.Reference) {
	counter := atomic.AddUint64(&m.GetReferencePreCounter, 1)
	defer atomic.AddUint64(&m.GetReferenceCounter, 1)

	if len(m.GetReferenceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetReferenceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetReference.")
			return
		}

		result := m.GetReferenceMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetReference")
			return
		}

		r = result.r

		return
	}

	if m.GetReferenceMock.mainExpectation != nil {

		result := m.GetReferenceMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetReference")
		}

		r = result.r

		return
	}

	if m.GetReferenceFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetReference.")
		return
	}

	return m.GetReferenceFunc()
}

//GetReferenceMinimockCounter returns a count of CandidateProfileMock.GetReferenceFunc invocations
func (m *CandidateProfileMock) GetReferenceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetReferenceCounter)
}

//GetReferenceMinimockPreCounter returns the value of CandidateProfileMock.GetReference invocations
func (m *CandidateProfileMock) GetReferenceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetReferencePreCounter)
}

//GetReferenceFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetReferenceFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetReferenceMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetReferenceCounter) == uint64(len(m.GetReferenceMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetReferenceMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetReferenceCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetReferenceFunc != nil {
		return atomic.LoadUint64(&m.GetReferenceCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetSpecialRoles struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetSpecialRolesExpectation
	expectationSeries []*CandidateProfileMockGetSpecialRolesExpectation
}

type CandidateProfileMockGetSpecialRolesExpectation struct {
	result *CandidateProfileMockGetSpecialRolesResult
}

type CandidateProfileMockGetSpecialRolesResult struct {
	r member.SpecialRole
}

//Expect specifies that invocation of CandidateProfile.GetSpecialRoles is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetSpecialRoles) Expect() *mCandidateProfileMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetSpecialRoles
func (m *mCandidateProfileMockGetSpecialRoles) Return(r member.SpecialRole) *CandidateProfileMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetSpecialRoles is expected once
func (m *mCandidateProfileMockGetSpecialRoles) ExpectOnce() *CandidateProfileMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetSpecialRolesExpectation) Return(r member.SpecialRole) {
	e.result = &CandidateProfileMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetSpecialRoles method
func (m *mCandidateProfileMockGetSpecialRoles) Set(f func() (r member.SpecialRole)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetSpecialRoles() (r member.SpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of CandidateProfileMock.GetSpecialRolesFunc invocations
func (m *CandidateProfileMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of CandidateProfileMock.GetSpecialRoles invocations
func (m *CandidateProfileMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetSpecialRolesFinished() bool {
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

type mCandidateProfileMockGetStartPower struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetStartPowerExpectation
	expectationSeries []*CandidateProfileMockGetStartPowerExpectation
}

type CandidateProfileMockGetStartPowerExpectation struct {
	result *CandidateProfileMockGetStartPowerResult
}

type CandidateProfileMockGetStartPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of CandidateProfile.GetStartPower is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetStartPower) Expect() *mCandidateProfileMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetStartPower
func (m *mCandidateProfileMockGetStartPower) Return(r member.Power) *CandidateProfileMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetStartPower is expected once
func (m *mCandidateProfileMockGetStartPower) ExpectOnce() *CandidateProfileMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetStartPowerExpectation) Return(r member.Power) {
	e.result = &CandidateProfileMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetStartPower method
func (m *mCandidateProfileMockGetStartPower) Set(f func() (r member.Power)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetStartPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of CandidateProfileMock.GetStartPowerFunc invocations
func (m *CandidateProfileMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of CandidateProfileMock.GetStartPower invocations
func (m *CandidateProfileMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetStartPowerFinished() bool {
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

type mCandidateProfileMockGetStaticNodeID struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetStaticNodeIDExpectation
	expectationSeries []*CandidateProfileMockGetStaticNodeIDExpectation
}

type CandidateProfileMockGetStaticNodeIDExpectation struct {
	result *CandidateProfileMockGetStaticNodeIDResult
}

type CandidateProfileMockGetStaticNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of CandidateProfile.GetStaticNodeID is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetStaticNodeID) Expect() *mCandidateProfileMockGetStaticNodeID {
	m.mock.GetStaticNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetStaticNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetStaticNodeID
func (m *mCandidateProfileMockGetStaticNodeID) Return(r insolar.ShortNodeID) *CandidateProfileMock {
	m.mock.GetStaticNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetStaticNodeIDExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetStaticNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetStaticNodeID is expected once
func (m *mCandidateProfileMockGetStaticNodeID) ExpectOnce() *CandidateProfileMockGetStaticNodeIDExpectation {
	m.mock.GetStaticNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetStaticNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetStaticNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &CandidateProfileMockGetStaticNodeIDResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetStaticNodeID method
func (m *mCandidateProfileMockGetStaticNodeID) Set(f func() (r insolar.ShortNodeID)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStaticNodeIDFunc = f
	return m.mock
}

//GetStaticNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile interface
func (m *CandidateProfileMock) GetStaticNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetStaticNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetStaticNodeIDCounter, 1)

	if len(m.GetStaticNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStaticNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetStaticNodeID.")
			return
		}

		result := m.GetStaticNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetStaticNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetStaticNodeIDMock.mainExpectation != nil {

		result := m.GetStaticNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetStaticNodeID")
		}

		r = result.r

		return
	}

	if m.GetStaticNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetStaticNodeID.")
		return
	}

	return m.GetStaticNodeIDFunc()
}

//GetStaticNodeIDMinimockCounter returns a count of CandidateProfileMock.GetStaticNodeIDFunc invocations
func (m *CandidateProfileMock) GetStaticNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticNodeIDCounter)
}

//GetStaticNodeIDMinimockPreCounter returns the value of CandidateProfileMock.GetStaticNodeID invocations
func (m *CandidateProfileMock) GetStaticNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticNodeIDPreCounter)
}

//GetStaticNodeIDFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetStaticNodeIDFinished() bool {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CandidateProfileMock) ValidateCallCounters() {

	if !m.GetBriefIntroSignedDigestFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetBriefIntroSignedDigest")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetDefaultEndpoint")
	}

	if !m.GetExtraEndpointsFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetExtraEndpoints")
	}

	if !m.GetIssuedAtPulseFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuedAtPulse")
	}

	if !m.GetIssuedAtTimeFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuedAtTime")
	}

	if !m.GetIssuerIDFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuerID")
	}

	if !m.GetIssuerSignatureFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuerSignature")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodePublicKey")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetPowerLevels")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetPrimaryRole")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetReference")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetStartPower")
	}

	if !m.GetStaticNodeIDFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetStaticNodeID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CandidateProfileMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CandidateProfileMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CandidateProfileMock) MinimockFinish() {

	if !m.GetBriefIntroSignedDigestFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetBriefIntroSignedDigest")
	}

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetDefaultEndpoint")
	}

	if !m.GetExtraEndpointsFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetExtraEndpoints")
	}

	if !m.GetIssuedAtPulseFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuedAtPulse")
	}

	if !m.GetIssuedAtTimeFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuedAtTime")
	}

	if !m.GetIssuerIDFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuerID")
	}

	if !m.GetIssuerSignatureFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetIssuerSignature")
	}

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodePublicKey")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetPowerLevels")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetPrimaryRole")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetReference")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetSpecialRoles")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetStartPower")
	}

	if !m.GetStaticNodeIDFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetStaticNodeID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CandidateProfileMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CandidateProfileMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetBriefIntroSignedDigestFinished()
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetExtraEndpointsFinished()
		ok = ok && m.GetIssuedAtPulseFinished()
		ok = ok && m.GetIssuedAtTimeFinished()
		ok = ok && m.GetIssuerIDFinished()
		ok = ok && m.GetIssuerSignatureFinished()
		ok = ok && m.GetNodePublicKeyFinished()
		ok = ok && m.GetPowerLevelsFinished()
		ok = ok && m.GetPrimaryRoleFinished()
		ok = ok && m.GetReferenceFinished()
		ok = ok && m.GetSpecialRolesFinished()
		ok = ok && m.GetStartPowerFinished()
		ok = ok && m.GetStaticNodeIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetBriefIntroSignedDigestFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetBriefIntroSignedDigest")
			}

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetDefaultEndpoint")
			}

			if !m.GetExtraEndpointsFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetExtraEndpoints")
			}

			if !m.GetIssuedAtPulseFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetIssuedAtPulse")
			}

			if !m.GetIssuedAtTimeFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetIssuedAtTime")
			}

			if !m.GetIssuerIDFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetIssuerID")
			}

			if !m.GetIssuerSignatureFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetIssuerSignature")
			}

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetNodePublicKey")
			}

			if !m.GetPowerLevelsFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetPowerLevels")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetPrimaryRole")
			}

			if !m.GetReferenceFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetReference")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetSpecialRoles")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetStartPower")
			}

			if !m.GetStaticNodeIDFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetStaticNodeID")
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
func (m *CandidateProfileMock) AllMocksCalled() bool {

	if !m.GetBriefIntroSignedDigestFinished() {
		return false
	}

	if !m.GetDefaultEndpointFinished() {
		return false
	}

	if !m.GetExtraEndpointsFinished() {
		return false
	}

	if !m.GetIssuedAtPulseFinished() {
		return false
	}

	if !m.GetIssuedAtTimeFinished() {
		return false
	}

	if !m.GetIssuerIDFinished() {
		return false
	}

	if !m.GetIssuerSignatureFinished() {
		return false
	}

	if !m.GetNodePublicKeyFinished() {
		return false
	}

	if !m.GetPowerLevelsFinished() {
		return false
	}

	if !m.GetPrimaryRoleFinished() {
		return false
	}

	if !m.GetReferenceFinished() {
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

	return true
}

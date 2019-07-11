package gcp_types

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CandidateProfile" can be found in github.com/insolar/insolar/network/consensus/gcpv2/gcp_types
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	cryptography_containers "github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	endpoints "github.com/insolar/insolar/network/consensus/common/endpoints"
	pulse_data "github.com/insolar/insolar/network/consensus/common/pulse_data"
)

//CandidateProfileMock implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile
type CandidateProfileMock struct {
	t minimock.Tester

	GetExtraEndpointsFunc       func() (r []endpoints.NodeEndpoint)
	GetExtraEndpointsCounter    uint64
	GetExtraEndpointsPreCounter uint64
	GetExtraEndpointsMock       mCandidateProfileMockGetExtraEndpoints

	GetIssuedAtPulseFunc       func() (r pulse_data.PulseNumber)
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

	GetIssuerSignatureFunc       func() (r cryptography_containers.SignatureHolder)
	GetIssuerSignatureCounter    uint64
	GetIssuerSignaturePreCounter uint64
	GetIssuerSignatureMock       mCandidateProfileMockGetIssuerSignature

	GetJoinerSignatureFunc       func() (r cryptography_containers.SignatureHolder)
	GetJoinerSignatureCounter    uint64
	GetJoinerSignaturePreCounter uint64
	GetJoinerSignatureMock       mCandidateProfileMockGetJoinerSignature

	GetNodeEndpointFunc       func() (r endpoints.NodeEndpoint)
	GetNodeEndpointCounter    uint64
	GetNodeEndpointPreCounter uint64
	GetNodeEndpointMock       mCandidateProfileMockGetNodeEndpoint

	GetNodeIDFunc       func() (r insolar.ShortNodeID)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mCandidateProfileMockGetNodeID

	GetNodePKFunc       func() (r cryptography_containers.SignatureKeyHolder)
	GetNodePKCounter    uint64
	GetNodePKPreCounter uint64
	GetNodePKMock       mCandidateProfileMockGetNodePK

	GetNodePrimaryRoleFunc       func() (r NodePrimaryRole)
	GetNodePrimaryRoleCounter    uint64
	GetNodePrimaryRolePreCounter uint64
	GetNodePrimaryRoleMock       mCandidateProfileMockGetNodePrimaryRole

	GetNodeSpecialRolesFunc       func() (r NodeSpecialRole)
	GetNodeSpecialRolesCounter    uint64
	GetNodeSpecialRolesPreCounter uint64
	GetNodeSpecialRolesMock       mCandidateProfileMockGetNodeSpecialRoles

	GetPowerLevelsFunc       func() (r MemberPowerSet)
	GetPowerLevelsCounter    uint64
	GetPowerLevelsPreCounter uint64
	GetPowerLevelsMock       mCandidateProfileMockGetPowerLevels

	GetReferenceFunc       func() (r insolar.Reference)
	GetReferenceCounter    uint64
	GetReferencePreCounter uint64
	GetReferenceMock       mCandidateProfileMockGetReference

	GetStartPowerFunc       func() (r MemberPower)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mCandidateProfileMockGetStartPower
}

//NewCandidateProfileMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile
func NewCandidateProfileMock(t minimock.Tester) *CandidateProfileMock {
	m := &CandidateProfileMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetExtraEndpointsMock = mCandidateProfileMockGetExtraEndpoints{mock: m}
	m.GetIssuedAtPulseMock = mCandidateProfileMockGetIssuedAtPulse{mock: m}
	m.GetIssuedAtTimeMock = mCandidateProfileMockGetIssuedAtTime{mock: m}
	m.GetIssuerIDMock = mCandidateProfileMockGetIssuerID{mock: m}
	m.GetIssuerSignatureMock = mCandidateProfileMockGetIssuerSignature{mock: m}
	m.GetJoinerSignatureMock = mCandidateProfileMockGetJoinerSignature{mock: m}
	m.GetNodeEndpointMock = mCandidateProfileMockGetNodeEndpoint{mock: m}
	m.GetNodeIDMock = mCandidateProfileMockGetNodeID{mock: m}
	m.GetNodePKMock = mCandidateProfileMockGetNodePK{mock: m}
	m.GetNodePrimaryRoleMock = mCandidateProfileMockGetNodePrimaryRole{mock: m}
	m.GetNodeSpecialRolesMock = mCandidateProfileMockGetNodeSpecialRoles{mock: m}
	m.GetPowerLevelsMock = mCandidateProfileMockGetPowerLevels{mock: m}
	m.GetReferenceMock = mCandidateProfileMockGetReference{mock: m}
	m.GetStartPowerMock = mCandidateProfileMockGetStartPower{mock: m}

	return m
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
	r []endpoints.NodeEndpoint
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
func (m *mCandidateProfileMockGetExtraEndpoints) Return(r []endpoints.NodeEndpoint) *CandidateProfileMock {
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

func (e *CandidateProfileMockGetExtraEndpointsExpectation) Return(r []endpoints.NodeEndpoint) {
	e.result = &CandidateProfileMockGetExtraEndpointsResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetExtraEndpoints method
func (m *mCandidateProfileMockGetExtraEndpoints) Set(f func() (r []endpoints.NodeEndpoint)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExtraEndpointsFunc = f
	return m.mock
}

//GetExtraEndpoints implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetExtraEndpoints() (r []endpoints.NodeEndpoint) {
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
	r pulse_data.PulseNumber
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
func (m *mCandidateProfileMockGetIssuedAtPulse) Return(r pulse_data.PulseNumber) *CandidateProfileMock {
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

func (e *CandidateProfileMockGetIssuedAtPulseExpectation) Return(r pulse_data.PulseNumber) {
	e.result = &CandidateProfileMockGetIssuedAtPulseResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetIssuedAtPulse method
func (m *mCandidateProfileMockGetIssuedAtPulse) Set(f func() (r pulse_data.PulseNumber)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtPulseFunc = f
	return m.mock
}

//GetIssuedAtPulse implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetIssuedAtPulse() (r pulse_data.PulseNumber) {
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

//GetIssuedAtTime implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
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

//GetIssuerID implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
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
	r cryptography_containers.SignatureHolder
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
func (m *mCandidateProfileMockGetIssuerSignature) Return(r cryptography_containers.SignatureHolder) *CandidateProfileMock {
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

func (e *CandidateProfileMockGetIssuerSignatureExpectation) Return(r cryptography_containers.SignatureHolder) {
	e.result = &CandidateProfileMockGetIssuerSignatureResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetIssuerSignature method
func (m *mCandidateProfileMockGetIssuerSignature) Set(f func() (r cryptography_containers.SignatureHolder)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerSignatureFunc = f
	return m.mock
}

//GetIssuerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetIssuerSignature() (r cryptography_containers.SignatureHolder) {
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

type mCandidateProfileMockGetJoinerSignature struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetJoinerSignatureExpectation
	expectationSeries []*CandidateProfileMockGetJoinerSignatureExpectation
}

type CandidateProfileMockGetJoinerSignatureExpectation struct {
	result *CandidateProfileMockGetJoinerSignatureResult
}

type CandidateProfileMockGetJoinerSignatureResult struct {
	r cryptography_containers.SignatureHolder
}

//Expect specifies that invocation of CandidateProfile.GetJoinerSignature is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetJoinerSignature) Expect() *mCandidateProfileMockGetJoinerSignature {
	m.mock.GetJoinerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetJoinerSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetJoinerSignature
func (m *mCandidateProfileMockGetJoinerSignature) Return(r cryptography_containers.SignatureHolder) *CandidateProfileMock {
	m.mock.GetJoinerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetJoinerSignatureExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetJoinerSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetJoinerSignature is expected once
func (m *mCandidateProfileMockGetJoinerSignature) ExpectOnce() *CandidateProfileMockGetJoinerSignatureExpectation {
	m.mock.GetJoinerSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetJoinerSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetJoinerSignatureExpectation) Return(r cryptography_containers.SignatureHolder) {
	e.result = &CandidateProfileMockGetJoinerSignatureResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetJoinerSignature method
func (m *mCandidateProfileMockGetJoinerSignature) Set(f func() (r cryptography_containers.SignatureHolder)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetJoinerSignatureFunc = f
	return m.mock
}

//GetJoinerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetJoinerSignature() (r cryptography_containers.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetJoinerSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetJoinerSignatureCounter, 1)

	if len(m.GetJoinerSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetJoinerSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetJoinerSignature.")
			return
		}

		result := m.GetJoinerSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetJoinerSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetJoinerSignatureMock.mainExpectation != nil {

		result := m.GetJoinerSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetJoinerSignature")
		}

		r = result.r

		return
	}

	if m.GetJoinerSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetJoinerSignature.")
		return
	}

	return m.GetJoinerSignatureFunc()
}

//GetJoinerSignatureMinimockCounter returns a count of CandidateProfileMock.GetJoinerSignatureFunc invocations
func (m *CandidateProfileMock) GetJoinerSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetJoinerSignatureCounter)
}

//GetJoinerSignatureMinimockPreCounter returns the value of CandidateProfileMock.GetJoinerSignature invocations
func (m *CandidateProfileMock) GetJoinerSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetJoinerSignaturePreCounter)
}

//GetJoinerSignatureFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetJoinerSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetJoinerSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetJoinerSignatureCounter) == uint64(len(m.GetJoinerSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetJoinerSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetJoinerSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetJoinerSignatureFunc != nil {
		return atomic.LoadUint64(&m.GetJoinerSignatureCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetNodeEndpoint struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetNodeEndpointExpectation
	expectationSeries []*CandidateProfileMockGetNodeEndpointExpectation
}

type CandidateProfileMockGetNodeEndpointExpectation struct {
	result *CandidateProfileMockGetNodeEndpointResult
}

type CandidateProfileMockGetNodeEndpointResult struct {
	r endpoints.NodeEndpoint
}

//Expect specifies that invocation of CandidateProfile.GetNodeEndpoint is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetNodeEndpoint) Expect() *mCandidateProfileMockGetNodeEndpoint {
	m.mock.GetNodeEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodeEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetNodeEndpoint
func (m *mCandidateProfileMockGetNodeEndpoint) Return(r endpoints.NodeEndpoint) *CandidateProfileMock {
	m.mock.GetNodeEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodeEndpointExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetNodeEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetNodeEndpoint is expected once
func (m *mCandidateProfileMockGetNodeEndpoint) ExpectOnce() *CandidateProfileMockGetNodeEndpointExpectation {
	m.mock.GetNodeEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetNodeEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetNodeEndpointExpectation) Return(r endpoints.NodeEndpoint) {
	e.result = &CandidateProfileMockGetNodeEndpointResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetNodeEndpoint method
func (m *mCandidateProfileMockGetNodeEndpoint) Set(f func() (r endpoints.NodeEndpoint)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeEndpointFunc = f
	return m.mock
}

//GetNodeEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetNodeEndpoint() (r endpoints.NodeEndpoint) {
	counter := atomic.AddUint64(&m.GetNodeEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeEndpointCounter, 1)

	if len(m.GetNodeEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodeEndpoint.")
			return
		}

		result := m.GetNodeEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodeEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeEndpointMock.mainExpectation != nil {

		result := m.GetNodeEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodeEndpoint")
		}

		r = result.r

		return
	}

	if m.GetNodeEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodeEndpoint.")
		return
	}

	return m.GetNodeEndpointFunc()
}

//GetNodeEndpointMinimockCounter returns a count of CandidateProfileMock.GetNodeEndpointFunc invocations
func (m *CandidateProfileMock) GetNodeEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeEndpointCounter)
}

//GetNodeEndpointMinimockPreCounter returns the value of CandidateProfileMock.GetNodeEndpoint invocations
func (m *CandidateProfileMock) GetNodeEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeEndpointPreCounter)
}

//GetNodeEndpointFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetNodeEndpointFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeEndpointMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeEndpointCounter) == uint64(len(m.GetNodeEndpointMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeEndpointMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeEndpointCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeEndpointFunc != nil {
		return atomic.LoadUint64(&m.GetNodeEndpointCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetNodeID struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetNodeIDExpectation
	expectationSeries []*CandidateProfileMockGetNodeIDExpectation
}

type CandidateProfileMockGetNodeIDExpectation struct {
	result *CandidateProfileMockGetNodeIDResult
}

type CandidateProfileMockGetNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of CandidateProfile.GetNodeID is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetNodeID) Expect() *mCandidateProfileMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetNodeID
func (m *mCandidateProfileMockGetNodeID) Return(r insolar.ShortNodeID) *CandidateProfileMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetNodeID is expected once
func (m *mCandidateProfileMockGetNodeID) ExpectOnce() *CandidateProfileMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &CandidateProfileMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetNodeID method
func (m *mCandidateProfileMockGetNodeID) Set(f func() (r insolar.ShortNodeID)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of CandidateProfileMock.GetNodeIDFunc invocations
func (m *CandidateProfileMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of CandidateProfileMock.GetNodeID invocations
func (m *CandidateProfileMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetNodeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeIDCounter) == uint64(len(m.GetNodeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeIDFunc != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetNodePK struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetNodePKExpectation
	expectationSeries []*CandidateProfileMockGetNodePKExpectation
}

type CandidateProfileMockGetNodePKExpectation struct {
	result *CandidateProfileMockGetNodePKResult
}

type CandidateProfileMockGetNodePKResult struct {
	r cryptography_containers.SignatureKeyHolder
}

//Expect specifies that invocation of CandidateProfile.GetNodePK is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetNodePK) Expect() *mCandidateProfileMockGetNodePK {
	m.mock.GetNodePKFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodePKExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetNodePK
func (m *mCandidateProfileMockGetNodePK) Return(r cryptography_containers.SignatureKeyHolder) *CandidateProfileMock {
	m.mock.GetNodePKFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodePKExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetNodePKResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetNodePK is expected once
func (m *mCandidateProfileMockGetNodePK) ExpectOnce() *CandidateProfileMockGetNodePKExpectation {
	m.mock.GetNodePKFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetNodePKExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetNodePKExpectation) Return(r cryptography_containers.SignatureKeyHolder) {
	e.result = &CandidateProfileMockGetNodePKResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetNodePK method
func (m *mCandidateProfileMockGetNodePK) Set(f func() (r cryptography_containers.SignatureKeyHolder)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePKFunc = f
	return m.mock
}

//GetNodePK implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetNodePK() (r cryptography_containers.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePKPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePKCounter, 1)

	if len(m.GetNodePKMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePKMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodePK.")
			return
		}

		result := m.GetNodePKMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodePK")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePKMock.mainExpectation != nil {

		result := m.GetNodePKMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodePK")
		}

		r = result.r

		return
	}

	if m.GetNodePKFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodePK.")
		return
	}

	return m.GetNodePKFunc()
}

//GetNodePKMinimockCounter returns a count of CandidateProfileMock.GetNodePKFunc invocations
func (m *CandidateProfileMock) GetNodePKMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePKCounter)
}

//GetNodePKMinimockPreCounter returns the value of CandidateProfileMock.GetNodePK invocations
func (m *CandidateProfileMock) GetNodePKMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePKPreCounter)
}

//GetNodePKFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetNodePKFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodePKMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodePKCounter) == uint64(len(m.GetNodePKMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodePKMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodePKCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodePKFunc != nil {
		return atomic.LoadUint64(&m.GetNodePKCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetNodePrimaryRole struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetNodePrimaryRoleExpectation
	expectationSeries []*CandidateProfileMockGetNodePrimaryRoleExpectation
}

type CandidateProfileMockGetNodePrimaryRoleExpectation struct {
	result *CandidateProfileMockGetNodePrimaryRoleResult
}

type CandidateProfileMockGetNodePrimaryRoleResult struct {
	r NodePrimaryRole
}

//Expect specifies that invocation of CandidateProfile.GetNodePrimaryRole is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetNodePrimaryRole) Expect() *mCandidateProfileMockGetNodePrimaryRole {
	m.mock.GetNodePrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodePrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetNodePrimaryRole
func (m *mCandidateProfileMockGetNodePrimaryRole) Return(r NodePrimaryRole) *CandidateProfileMock {
	m.mock.GetNodePrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodePrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetNodePrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetNodePrimaryRole is expected once
func (m *mCandidateProfileMockGetNodePrimaryRole) ExpectOnce() *CandidateProfileMockGetNodePrimaryRoleExpectation {
	m.mock.GetNodePrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetNodePrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetNodePrimaryRoleExpectation) Return(r NodePrimaryRole) {
	e.result = &CandidateProfileMockGetNodePrimaryRoleResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetNodePrimaryRole method
func (m *mCandidateProfileMockGetNodePrimaryRole) Set(f func() (r NodePrimaryRole)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePrimaryRoleFunc = f
	return m.mock
}

//GetNodePrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetNodePrimaryRole() (r NodePrimaryRole) {
	counter := atomic.AddUint64(&m.GetNodePrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePrimaryRoleCounter, 1)

	if len(m.GetNodePrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodePrimaryRole.")
			return
		}

		result := m.GetNodePrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodePrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePrimaryRoleMock.mainExpectation != nil {

		result := m.GetNodePrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodePrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetNodePrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodePrimaryRole.")
		return
	}

	return m.GetNodePrimaryRoleFunc()
}

//GetNodePrimaryRoleMinimockCounter returns a count of CandidateProfileMock.GetNodePrimaryRoleFunc invocations
func (m *CandidateProfileMock) GetNodePrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePrimaryRoleCounter)
}

//GetNodePrimaryRoleMinimockPreCounter returns the value of CandidateProfileMock.GetNodePrimaryRole invocations
func (m *CandidateProfileMock) GetNodePrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePrimaryRolePreCounter)
}

//GetNodePrimaryRoleFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetNodePrimaryRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodePrimaryRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodePrimaryRoleCounter) == uint64(len(m.GetNodePrimaryRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodePrimaryRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodePrimaryRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodePrimaryRoleFunc != nil {
		return atomic.LoadUint64(&m.GetNodePrimaryRoleCounter) > 0
	}

	return true
}

type mCandidateProfileMockGetNodeSpecialRoles struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetNodeSpecialRolesExpectation
	expectationSeries []*CandidateProfileMockGetNodeSpecialRolesExpectation
}

type CandidateProfileMockGetNodeSpecialRolesExpectation struct {
	result *CandidateProfileMockGetNodeSpecialRolesResult
}

type CandidateProfileMockGetNodeSpecialRolesResult struct {
	r NodeSpecialRole
}

//Expect specifies that invocation of CandidateProfile.GetNodeSpecialRoles is expected from 1 to Infinity times
func (m *mCandidateProfileMockGetNodeSpecialRoles) Expect() *mCandidateProfileMockGetNodeSpecialRoles {
	m.mock.GetNodeSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodeSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateProfile.GetNodeSpecialRoles
func (m *mCandidateProfileMockGetNodeSpecialRoles) Return(r NodeSpecialRole) *CandidateProfileMock {
	m.mock.GetNodeSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateProfileMockGetNodeSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &CandidateProfileMockGetNodeSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateProfile.GetNodeSpecialRoles is expected once
func (m *mCandidateProfileMockGetNodeSpecialRoles) ExpectOnce() *CandidateProfileMockGetNodeSpecialRolesExpectation {
	m.mock.GetNodeSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateProfileMockGetNodeSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateProfileMockGetNodeSpecialRolesExpectation) Return(r NodeSpecialRole) {
	e.result = &CandidateProfileMockGetNodeSpecialRolesResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetNodeSpecialRoles method
func (m *mCandidateProfileMockGetNodeSpecialRoles) Set(f func() (r NodeSpecialRole)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeSpecialRolesFunc = f
	return m.mock
}

//GetNodeSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetNodeSpecialRoles() (r NodeSpecialRole) {
	counter := atomic.AddUint64(&m.GetNodeSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeSpecialRolesCounter, 1)

	if len(m.GetNodeSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodeSpecialRoles.")
			return
		}

		result := m.GetNodeSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodeSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeSpecialRolesMock.mainExpectation != nil {

		result := m.GetNodeSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateProfileMock.GetNodeSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetNodeSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateProfileMock.GetNodeSpecialRoles.")
		return
	}

	return m.GetNodeSpecialRolesFunc()
}

//GetNodeSpecialRolesMinimockCounter returns a count of CandidateProfileMock.GetNodeSpecialRolesFunc invocations
func (m *CandidateProfileMock) GetNodeSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeSpecialRolesCounter)
}

//GetNodeSpecialRolesMinimockPreCounter returns the value of CandidateProfileMock.GetNodeSpecialRoles invocations
func (m *CandidateProfileMock) GetNodeSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeSpecialRolesPreCounter)
}

//GetNodeSpecialRolesFinished returns true if mock invocations count is ok
func (m *CandidateProfileMock) GetNodeSpecialRolesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeSpecialRolesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeSpecialRolesCounter) == uint64(len(m.GetNodeSpecialRolesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeSpecialRolesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeSpecialRolesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeSpecialRolesFunc != nil {
		return atomic.LoadUint64(&m.GetNodeSpecialRolesCounter) > 0
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
	r MemberPowerSet
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
func (m *mCandidateProfileMockGetPowerLevels) Return(r MemberPowerSet) *CandidateProfileMock {
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

func (e *CandidateProfileMockGetPowerLevelsExpectation) Return(r MemberPowerSet) {
	e.result = &CandidateProfileMockGetPowerLevelsResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetPowerLevels method
func (m *mCandidateProfileMockGetPowerLevels) Set(f func() (r MemberPowerSet)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPowerLevelsFunc = f
	return m.mock
}

//GetPowerLevels implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetPowerLevels() (r MemberPowerSet) {
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

//GetReference implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
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

type mCandidateProfileMockGetStartPower struct {
	mock              *CandidateProfileMock
	mainExpectation   *CandidateProfileMockGetStartPowerExpectation
	expectationSeries []*CandidateProfileMockGetStartPowerExpectation
}

type CandidateProfileMockGetStartPowerExpectation struct {
	result *CandidateProfileMockGetStartPowerResult
}

type CandidateProfileMockGetStartPowerResult struct {
	r MemberPower
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
func (m *mCandidateProfileMockGetStartPower) Return(r MemberPower) *CandidateProfileMock {
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

func (e *CandidateProfileMockGetStartPowerExpectation) Return(r MemberPower) {
	e.result = &CandidateProfileMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of CandidateProfile.GetStartPower method
func (m *mCandidateProfileMockGetStartPower) Set(f func() (r MemberPower)) *CandidateProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/gcp_types.CandidateProfile interface
func (m *CandidateProfileMock) GetStartPower() (r MemberPower) {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CandidateProfileMock) ValidateCallCounters() {

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

	if !m.GetJoinerSignatureFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetJoinerSignature")
	}

	if !m.GetNodeEndpointFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodeEndpoint")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodeID")
	}

	if !m.GetNodePKFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodePK")
	}

	if !m.GetNodePrimaryRoleFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodePrimaryRole")
	}

	if !m.GetNodeSpecialRolesFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodeSpecialRoles")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetPowerLevels")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetReference")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetStartPower")
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

	if !m.GetJoinerSignatureFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetJoinerSignature")
	}

	if !m.GetNodeEndpointFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodeEndpoint")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodeID")
	}

	if !m.GetNodePKFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodePK")
	}

	if !m.GetNodePrimaryRoleFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodePrimaryRole")
	}

	if !m.GetNodeSpecialRolesFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetNodeSpecialRoles")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetPowerLevels")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetReference")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to CandidateProfileMock.GetStartPower")
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
		ok = ok && m.GetExtraEndpointsFinished()
		ok = ok && m.GetIssuedAtPulseFinished()
		ok = ok && m.GetIssuedAtTimeFinished()
		ok = ok && m.GetIssuerIDFinished()
		ok = ok && m.GetIssuerSignatureFinished()
		ok = ok && m.GetJoinerSignatureFinished()
		ok = ok && m.GetNodeEndpointFinished()
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetNodePKFinished()
		ok = ok && m.GetNodePrimaryRoleFinished()
		ok = ok && m.GetNodeSpecialRolesFinished()
		ok = ok && m.GetPowerLevelsFinished()
		ok = ok && m.GetReferenceFinished()
		ok = ok && m.GetStartPowerFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

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

			if !m.GetJoinerSignatureFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetJoinerSignature")
			}

			if !m.GetNodeEndpointFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetNodeEndpoint")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetNodeID")
			}

			if !m.GetNodePKFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetNodePK")
			}

			if !m.GetNodePrimaryRoleFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetNodePrimaryRole")
			}

			if !m.GetNodeSpecialRolesFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetNodeSpecialRoles")
			}

			if !m.GetPowerLevelsFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetPowerLevels")
			}

			if !m.GetReferenceFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetReference")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to CandidateProfileMock.GetStartPower")
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

	if !m.GetJoinerSignatureFinished() {
		return false
	}

	if !m.GetNodeEndpointFinished() {
		return false
	}

	if !m.GetNodeIDFinished() {
		return false
	}

	if !m.GetNodePKFinished() {
		return false
	}

	if !m.GetNodePrimaryRoleFinished() {
		return false
	}

	if !m.GetNodeSpecialRolesFinished() {
		return false
	}

	if !m.GetPowerLevelsFinished() {
		return false
	}

	if !m.GetReferenceFinished() {
		return false
	}

	if !m.GetStartPowerFinished() {
		return false
	}

	return true
}

package packets

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FullIntroductionReader" can be found in github.com/insolar/insolar/network/consensus/gcpv2/packets
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	common "github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

//FullIntroductionReaderMock implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader
type FullIntroductionReaderMock struct {
	t minimock.Tester

	GetExtraEndpointsFunc       func() (r []common.NodeEndpoint)
	GetExtraEndpointsCounter    uint64
	GetExtraEndpointsPreCounter uint64
	GetExtraEndpointsMock       mFullIntroductionReaderMockGetExtraEndpoints

	GetIssuedAtPulseFunc       func() (r common.PulseNumber)
	GetIssuedAtPulseCounter    uint64
	GetIssuedAtPulsePreCounter uint64
	GetIssuedAtPulseMock       mFullIntroductionReaderMockGetIssuedAtPulse

	GetIssuedAtTimeFunc       func() (r time.Time)
	GetIssuedAtTimeCounter    uint64
	GetIssuedAtTimePreCounter uint64
	GetIssuedAtTimeMock       mFullIntroductionReaderMockGetIssuedAtTime

	GetIssuerIDFunc       func() (r common.ShortNodeID)
	GetIssuerIDCounter    uint64
	GetIssuerIDPreCounter uint64
	GetIssuerIDMock       mFullIntroductionReaderMockGetIssuerID

	GetIssuerSignatureFunc       func() (r common.SignatureHolder)
	GetIssuerSignatureCounter    uint64
	GetIssuerSignaturePreCounter uint64
	GetIssuerSignatureMock       mFullIntroductionReaderMockGetIssuerSignature

	GetJoinerSignatureFunc       func() (r common.SignatureHolder)
	GetJoinerSignatureCounter    uint64
	GetJoinerSignaturePreCounter uint64
	GetJoinerSignatureMock       mFullIntroductionReaderMockGetJoinerSignature

	GetNodeEndpointFunc       func() (r common.NodeEndpoint)
	GetNodeEndpointCounter    uint64
	GetNodeEndpointPreCounter uint64
	GetNodeEndpointMock       mFullIntroductionReaderMockGetNodeEndpoint

	GetNodeIDFunc       func() (r common.ShortNodeID)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mFullIntroductionReaderMockGetNodeID

	GetNodePKFunc       func() (r common.SignatureKeyHolder)
	GetNodePKCounter    uint64
	GetNodePKPreCounter uint64
	GetNodePKMock       mFullIntroductionReaderMockGetNodePK

	GetNodePrimaryRoleFunc       func() (r common2.NodePrimaryRole)
	GetNodePrimaryRoleCounter    uint64
	GetNodePrimaryRolePreCounter uint64
	GetNodePrimaryRoleMock       mFullIntroductionReaderMockGetNodePrimaryRole

	GetNodeSpecialRolesFunc       func() (r common2.NodeSpecialRole)
	GetNodeSpecialRolesCounter    uint64
	GetNodeSpecialRolesPreCounter uint64
	GetNodeSpecialRolesMock       mFullIntroductionReaderMockGetNodeSpecialRoles

	GetPowerLevelsFunc       func() (r common2.MemberPowerSet)
	GetPowerLevelsCounter    uint64
	GetPowerLevelsPreCounter uint64
	GetPowerLevelsMock       mFullIntroductionReaderMockGetPowerLevels

	GetReferenceFunc       func() (r insolar.Reference)
	GetReferenceCounter    uint64
	GetReferencePreCounter uint64
	GetReferenceMock       mFullIntroductionReaderMockGetReference

	GetStartPowerFunc       func() (r common2.MemberPower)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mFullIntroductionReaderMockGetStartPower
}

//NewFullIntroductionReaderMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader
func NewFullIntroductionReaderMock(t minimock.Tester) *FullIntroductionReaderMock {
	m := &FullIntroductionReaderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetExtraEndpointsMock = mFullIntroductionReaderMockGetExtraEndpoints{mock: m}
	m.GetIssuedAtPulseMock = mFullIntroductionReaderMockGetIssuedAtPulse{mock: m}
	m.GetIssuedAtTimeMock = mFullIntroductionReaderMockGetIssuedAtTime{mock: m}
	m.GetIssuerIDMock = mFullIntroductionReaderMockGetIssuerID{mock: m}
	m.GetIssuerSignatureMock = mFullIntroductionReaderMockGetIssuerSignature{mock: m}
	m.GetJoinerSignatureMock = mFullIntroductionReaderMockGetJoinerSignature{mock: m}
	m.GetNodeEndpointMock = mFullIntroductionReaderMockGetNodeEndpoint{mock: m}
	m.GetNodeIDMock = mFullIntroductionReaderMockGetNodeID{mock: m}
	m.GetNodePKMock = mFullIntroductionReaderMockGetNodePK{mock: m}
	m.GetNodePrimaryRoleMock = mFullIntroductionReaderMockGetNodePrimaryRole{mock: m}
	m.GetNodeSpecialRolesMock = mFullIntroductionReaderMockGetNodeSpecialRoles{mock: m}
	m.GetPowerLevelsMock = mFullIntroductionReaderMockGetPowerLevels{mock: m}
	m.GetReferenceMock = mFullIntroductionReaderMockGetReference{mock: m}
	m.GetStartPowerMock = mFullIntroductionReaderMockGetStartPower{mock: m}

	return m
}

type mFullIntroductionReaderMockGetExtraEndpoints struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetExtraEndpointsExpectation
	expectationSeries []*FullIntroductionReaderMockGetExtraEndpointsExpectation
}

type FullIntroductionReaderMockGetExtraEndpointsExpectation struct {
	result *FullIntroductionReaderMockGetExtraEndpointsResult
}

type FullIntroductionReaderMockGetExtraEndpointsResult struct {
	r []common.NodeEndpoint
}

//Expect specifies that invocation of FullIntroductionReader.GetExtraEndpoints is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetExtraEndpoints) Expect() *mFullIntroductionReaderMockGetExtraEndpoints {
	m.mock.GetExtraEndpointsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetExtraEndpointsExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetExtraEndpoints
func (m *mFullIntroductionReaderMockGetExtraEndpoints) Return(r []common.NodeEndpoint) *FullIntroductionReaderMock {
	m.mock.GetExtraEndpointsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetExtraEndpointsExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetExtraEndpointsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetExtraEndpoints is expected once
func (m *mFullIntroductionReaderMockGetExtraEndpoints) ExpectOnce() *FullIntroductionReaderMockGetExtraEndpointsExpectation {
	m.mock.GetExtraEndpointsFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetExtraEndpointsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetExtraEndpointsExpectation) Return(r []common.NodeEndpoint) {
	e.result = &FullIntroductionReaderMockGetExtraEndpointsResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetExtraEndpoints method
func (m *mFullIntroductionReaderMockGetExtraEndpoints) Set(f func() (r []common.NodeEndpoint)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExtraEndpointsFunc = f
	return m.mock
}

//GetExtraEndpoints implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetExtraEndpoints() (r []common.NodeEndpoint) {
	counter := atomic.AddUint64(&m.GetExtraEndpointsPreCounter, 1)
	defer atomic.AddUint64(&m.GetExtraEndpointsCounter, 1)

	if len(m.GetExtraEndpointsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExtraEndpointsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetExtraEndpoints.")
			return
		}

		result := m.GetExtraEndpointsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetExtraEndpoints")
			return
		}

		r = result.r

		return
	}

	if m.GetExtraEndpointsMock.mainExpectation != nil {

		result := m.GetExtraEndpointsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetExtraEndpoints")
		}

		r = result.r

		return
	}

	if m.GetExtraEndpointsFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetExtraEndpoints.")
		return
	}

	return m.GetExtraEndpointsFunc()
}

//GetExtraEndpointsMinimockCounter returns a count of FullIntroductionReaderMock.GetExtraEndpointsFunc invocations
func (m *FullIntroductionReaderMock) GetExtraEndpointsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtraEndpointsCounter)
}

//GetExtraEndpointsMinimockPreCounter returns the value of FullIntroductionReaderMock.GetExtraEndpoints invocations
func (m *FullIntroductionReaderMock) GetExtraEndpointsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtraEndpointsPreCounter)
}

//GetExtraEndpointsFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetExtraEndpointsFinished() bool {
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

type mFullIntroductionReaderMockGetIssuedAtPulse struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetIssuedAtPulseExpectation
	expectationSeries []*FullIntroductionReaderMockGetIssuedAtPulseExpectation
}

type FullIntroductionReaderMockGetIssuedAtPulseExpectation struct {
	result *FullIntroductionReaderMockGetIssuedAtPulseResult
}

type FullIntroductionReaderMockGetIssuedAtPulseResult struct {
	r common.PulseNumber
}

//Expect specifies that invocation of FullIntroductionReader.GetIssuedAtPulse is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetIssuedAtPulse) Expect() *mFullIntroductionReaderMockGetIssuedAtPulse {
	m.mock.GetIssuedAtPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuedAtPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetIssuedAtPulse
func (m *mFullIntroductionReaderMockGetIssuedAtPulse) Return(r common.PulseNumber) *FullIntroductionReaderMock {
	m.mock.GetIssuedAtPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuedAtPulseExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetIssuedAtPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetIssuedAtPulse is expected once
func (m *mFullIntroductionReaderMockGetIssuedAtPulse) ExpectOnce() *FullIntroductionReaderMockGetIssuedAtPulseExpectation {
	m.mock.GetIssuedAtPulseFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetIssuedAtPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetIssuedAtPulseExpectation) Return(r common.PulseNumber) {
	e.result = &FullIntroductionReaderMockGetIssuedAtPulseResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetIssuedAtPulse method
func (m *mFullIntroductionReaderMockGetIssuedAtPulse) Set(f func() (r common.PulseNumber)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtPulseFunc = f
	return m.mock
}

//GetIssuedAtPulse implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetIssuedAtPulse() (r common.PulseNumber) {
	counter := atomic.AddUint64(&m.GetIssuedAtPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuedAtPulseCounter, 1)

	if len(m.GetIssuedAtPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuedAtPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuedAtPulse.")
			return
		}

		result := m.GetIssuedAtPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuedAtPulse")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuedAtPulseMock.mainExpectation != nil {

		result := m.GetIssuedAtPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuedAtPulse")
		}

		r = result.r

		return
	}

	if m.GetIssuedAtPulseFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuedAtPulse.")
		return
	}

	return m.GetIssuedAtPulseFunc()
}

//GetIssuedAtPulseMinimockCounter returns a count of FullIntroductionReaderMock.GetIssuedAtPulseFunc invocations
func (m *FullIntroductionReaderMock) GetIssuedAtPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtPulseCounter)
}

//GetIssuedAtPulseMinimockPreCounter returns the value of FullIntroductionReaderMock.GetIssuedAtPulse invocations
func (m *FullIntroductionReaderMock) GetIssuedAtPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtPulsePreCounter)
}

//GetIssuedAtPulseFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetIssuedAtPulseFinished() bool {
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

type mFullIntroductionReaderMockGetIssuedAtTime struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetIssuedAtTimeExpectation
	expectationSeries []*FullIntroductionReaderMockGetIssuedAtTimeExpectation
}

type FullIntroductionReaderMockGetIssuedAtTimeExpectation struct {
	result *FullIntroductionReaderMockGetIssuedAtTimeResult
}

type FullIntroductionReaderMockGetIssuedAtTimeResult struct {
	r time.Time
}

//Expect specifies that invocation of FullIntroductionReader.GetIssuedAtTime is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetIssuedAtTime) Expect() *mFullIntroductionReaderMockGetIssuedAtTime {
	m.mock.GetIssuedAtTimeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuedAtTimeExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetIssuedAtTime
func (m *mFullIntroductionReaderMockGetIssuedAtTime) Return(r time.Time) *FullIntroductionReaderMock {
	m.mock.GetIssuedAtTimeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuedAtTimeExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetIssuedAtTimeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetIssuedAtTime is expected once
func (m *mFullIntroductionReaderMockGetIssuedAtTime) ExpectOnce() *FullIntroductionReaderMockGetIssuedAtTimeExpectation {
	m.mock.GetIssuedAtTimeFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetIssuedAtTimeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetIssuedAtTimeExpectation) Return(r time.Time) {
	e.result = &FullIntroductionReaderMockGetIssuedAtTimeResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetIssuedAtTime method
func (m *mFullIntroductionReaderMockGetIssuedAtTime) Set(f func() (r time.Time)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtTimeFunc = f
	return m.mock
}

//GetIssuedAtTime implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetIssuedAtTime() (r time.Time) {
	counter := atomic.AddUint64(&m.GetIssuedAtTimePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuedAtTimeCounter, 1)

	if len(m.GetIssuedAtTimeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuedAtTimeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuedAtTime.")
			return
		}

		result := m.GetIssuedAtTimeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuedAtTime")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuedAtTimeMock.mainExpectation != nil {

		result := m.GetIssuedAtTimeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuedAtTime")
		}

		r = result.r

		return
	}

	if m.GetIssuedAtTimeFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuedAtTime.")
		return
	}

	return m.GetIssuedAtTimeFunc()
}

//GetIssuedAtTimeMinimockCounter returns a count of FullIntroductionReaderMock.GetIssuedAtTimeFunc invocations
func (m *FullIntroductionReaderMock) GetIssuedAtTimeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtTimeCounter)
}

//GetIssuedAtTimeMinimockPreCounter returns the value of FullIntroductionReaderMock.GetIssuedAtTime invocations
func (m *FullIntroductionReaderMock) GetIssuedAtTimeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtTimePreCounter)
}

//GetIssuedAtTimeFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetIssuedAtTimeFinished() bool {
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

type mFullIntroductionReaderMockGetIssuerID struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetIssuerIDExpectation
	expectationSeries []*FullIntroductionReaderMockGetIssuerIDExpectation
}

type FullIntroductionReaderMockGetIssuerIDExpectation struct {
	result *FullIntroductionReaderMockGetIssuerIDResult
}

type FullIntroductionReaderMockGetIssuerIDResult struct {
	r common.ShortNodeID
}

//Expect specifies that invocation of FullIntroductionReader.GetIssuerID is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetIssuerID) Expect() *mFullIntroductionReaderMockGetIssuerID {
	m.mock.GetIssuerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetIssuerID
func (m *mFullIntroductionReaderMockGetIssuerID) Return(r common.ShortNodeID) *FullIntroductionReaderMock {
	m.mock.GetIssuerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuerIDExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetIssuerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetIssuerID is expected once
func (m *mFullIntroductionReaderMockGetIssuerID) ExpectOnce() *FullIntroductionReaderMockGetIssuerIDExpectation {
	m.mock.GetIssuerIDFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetIssuerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetIssuerIDExpectation) Return(r common.ShortNodeID) {
	e.result = &FullIntroductionReaderMockGetIssuerIDResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetIssuerID method
func (m *mFullIntroductionReaderMockGetIssuerID) Set(f func() (r common.ShortNodeID)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerIDFunc = f
	return m.mock
}

//GetIssuerID implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetIssuerID() (r common.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetIssuerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuerIDCounter, 1)

	if len(m.GetIssuerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuerID.")
			return
		}

		result := m.GetIssuerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuerID")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuerIDMock.mainExpectation != nil {

		result := m.GetIssuerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuerID")
		}

		r = result.r

		return
	}

	if m.GetIssuerIDFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuerID.")
		return
	}

	return m.GetIssuerIDFunc()
}

//GetIssuerIDMinimockCounter returns a count of FullIntroductionReaderMock.GetIssuerIDFunc invocations
func (m *FullIntroductionReaderMock) GetIssuerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerIDCounter)
}

//GetIssuerIDMinimockPreCounter returns the value of FullIntroductionReaderMock.GetIssuerID invocations
func (m *FullIntroductionReaderMock) GetIssuerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerIDPreCounter)
}

//GetIssuerIDFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetIssuerIDFinished() bool {
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

type mFullIntroductionReaderMockGetIssuerSignature struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetIssuerSignatureExpectation
	expectationSeries []*FullIntroductionReaderMockGetIssuerSignatureExpectation
}

type FullIntroductionReaderMockGetIssuerSignatureExpectation struct {
	result *FullIntroductionReaderMockGetIssuerSignatureResult
}

type FullIntroductionReaderMockGetIssuerSignatureResult struct {
	r common.SignatureHolder
}

//Expect specifies that invocation of FullIntroductionReader.GetIssuerSignature is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetIssuerSignature) Expect() *mFullIntroductionReaderMockGetIssuerSignature {
	m.mock.GetIssuerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuerSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetIssuerSignature
func (m *mFullIntroductionReaderMockGetIssuerSignature) Return(r common.SignatureHolder) *FullIntroductionReaderMock {
	m.mock.GetIssuerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetIssuerSignatureExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetIssuerSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetIssuerSignature is expected once
func (m *mFullIntroductionReaderMockGetIssuerSignature) ExpectOnce() *FullIntroductionReaderMockGetIssuerSignatureExpectation {
	m.mock.GetIssuerSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetIssuerSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetIssuerSignatureExpectation) Return(r common.SignatureHolder) {
	e.result = &FullIntroductionReaderMockGetIssuerSignatureResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetIssuerSignature method
func (m *mFullIntroductionReaderMockGetIssuerSignature) Set(f func() (r common.SignatureHolder)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerSignatureFunc = f
	return m.mock
}

//GetIssuerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetIssuerSignature() (r common.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetIssuerSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuerSignatureCounter, 1)

	if len(m.GetIssuerSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuerSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuerSignature.")
			return
		}

		result := m.GetIssuerSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuerSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuerSignatureMock.mainExpectation != nil {

		result := m.GetIssuerSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetIssuerSignature")
		}

		r = result.r

		return
	}

	if m.GetIssuerSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetIssuerSignature.")
		return
	}

	return m.GetIssuerSignatureFunc()
}

//GetIssuerSignatureMinimockCounter returns a count of FullIntroductionReaderMock.GetIssuerSignatureFunc invocations
func (m *FullIntroductionReaderMock) GetIssuerSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerSignatureCounter)
}

//GetIssuerSignatureMinimockPreCounter returns the value of FullIntroductionReaderMock.GetIssuerSignature invocations
func (m *FullIntroductionReaderMock) GetIssuerSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerSignaturePreCounter)
}

//GetIssuerSignatureFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetIssuerSignatureFinished() bool {
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

type mFullIntroductionReaderMockGetJoinerSignature struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetJoinerSignatureExpectation
	expectationSeries []*FullIntroductionReaderMockGetJoinerSignatureExpectation
}

type FullIntroductionReaderMockGetJoinerSignatureExpectation struct {
	result *FullIntroductionReaderMockGetJoinerSignatureResult
}

type FullIntroductionReaderMockGetJoinerSignatureResult struct {
	r common.SignatureHolder
}

//Expect specifies that invocation of FullIntroductionReader.GetJoinerSignature is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetJoinerSignature) Expect() *mFullIntroductionReaderMockGetJoinerSignature {
	m.mock.GetJoinerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetJoinerSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetJoinerSignature
func (m *mFullIntroductionReaderMockGetJoinerSignature) Return(r common.SignatureHolder) *FullIntroductionReaderMock {
	m.mock.GetJoinerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetJoinerSignatureExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetJoinerSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetJoinerSignature is expected once
func (m *mFullIntroductionReaderMockGetJoinerSignature) ExpectOnce() *FullIntroductionReaderMockGetJoinerSignatureExpectation {
	m.mock.GetJoinerSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetJoinerSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetJoinerSignatureExpectation) Return(r common.SignatureHolder) {
	e.result = &FullIntroductionReaderMockGetJoinerSignatureResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetJoinerSignature method
func (m *mFullIntroductionReaderMockGetJoinerSignature) Set(f func() (r common.SignatureHolder)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetJoinerSignatureFunc = f
	return m.mock
}

//GetJoinerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetJoinerSignature() (r common.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetJoinerSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetJoinerSignatureCounter, 1)

	if len(m.GetJoinerSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetJoinerSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetJoinerSignature.")
			return
		}

		result := m.GetJoinerSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetJoinerSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetJoinerSignatureMock.mainExpectation != nil {

		result := m.GetJoinerSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetJoinerSignature")
		}

		r = result.r

		return
	}

	if m.GetJoinerSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetJoinerSignature.")
		return
	}

	return m.GetJoinerSignatureFunc()
}

//GetJoinerSignatureMinimockCounter returns a count of FullIntroductionReaderMock.GetJoinerSignatureFunc invocations
func (m *FullIntroductionReaderMock) GetJoinerSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetJoinerSignatureCounter)
}

//GetJoinerSignatureMinimockPreCounter returns the value of FullIntroductionReaderMock.GetJoinerSignature invocations
func (m *FullIntroductionReaderMock) GetJoinerSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetJoinerSignaturePreCounter)
}

//GetJoinerSignatureFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetJoinerSignatureFinished() bool {
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

type mFullIntroductionReaderMockGetNodeEndpoint struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetNodeEndpointExpectation
	expectationSeries []*FullIntroductionReaderMockGetNodeEndpointExpectation
}

type FullIntroductionReaderMockGetNodeEndpointExpectation struct {
	result *FullIntroductionReaderMockGetNodeEndpointResult
}

type FullIntroductionReaderMockGetNodeEndpointResult struct {
	r common.NodeEndpoint
}

//Expect specifies that invocation of FullIntroductionReader.GetNodeEndpoint is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetNodeEndpoint) Expect() *mFullIntroductionReaderMockGetNodeEndpoint {
	m.mock.GetNodeEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodeEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetNodeEndpoint
func (m *mFullIntroductionReaderMockGetNodeEndpoint) Return(r common.NodeEndpoint) *FullIntroductionReaderMock {
	m.mock.GetNodeEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodeEndpointExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetNodeEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetNodeEndpoint is expected once
func (m *mFullIntroductionReaderMockGetNodeEndpoint) ExpectOnce() *FullIntroductionReaderMockGetNodeEndpointExpectation {
	m.mock.GetNodeEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetNodeEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetNodeEndpointExpectation) Return(r common.NodeEndpoint) {
	e.result = &FullIntroductionReaderMockGetNodeEndpointResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetNodeEndpoint method
func (m *mFullIntroductionReaderMockGetNodeEndpoint) Set(f func() (r common.NodeEndpoint)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeEndpointFunc = f
	return m.mock
}

//GetNodeEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetNodeEndpoint() (r common.NodeEndpoint) {
	counter := atomic.AddUint64(&m.GetNodeEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeEndpointCounter, 1)

	if len(m.GetNodeEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodeEndpoint.")
			return
		}

		result := m.GetNodeEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodeEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeEndpointMock.mainExpectation != nil {

		result := m.GetNodeEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodeEndpoint")
		}

		r = result.r

		return
	}

	if m.GetNodeEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodeEndpoint.")
		return
	}

	return m.GetNodeEndpointFunc()
}

//GetNodeEndpointMinimockCounter returns a count of FullIntroductionReaderMock.GetNodeEndpointFunc invocations
func (m *FullIntroductionReaderMock) GetNodeEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeEndpointCounter)
}

//GetNodeEndpointMinimockPreCounter returns the value of FullIntroductionReaderMock.GetNodeEndpoint invocations
func (m *FullIntroductionReaderMock) GetNodeEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeEndpointPreCounter)
}

//GetNodeEndpointFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetNodeEndpointFinished() bool {
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

type mFullIntroductionReaderMockGetNodeID struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetNodeIDExpectation
	expectationSeries []*FullIntroductionReaderMockGetNodeIDExpectation
}

type FullIntroductionReaderMockGetNodeIDExpectation struct {
	result *FullIntroductionReaderMockGetNodeIDResult
}

type FullIntroductionReaderMockGetNodeIDResult struct {
	r common.ShortNodeID
}

//Expect specifies that invocation of FullIntroductionReader.GetNodeID is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetNodeID) Expect() *mFullIntroductionReaderMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetNodeID
func (m *mFullIntroductionReaderMockGetNodeID) Return(r common.ShortNodeID) *FullIntroductionReaderMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetNodeID is expected once
func (m *mFullIntroductionReaderMockGetNodeID) ExpectOnce() *FullIntroductionReaderMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetNodeIDExpectation) Return(r common.ShortNodeID) {
	e.result = &FullIntroductionReaderMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetNodeID method
func (m *mFullIntroductionReaderMockGetNodeID) Set(f func() (r common.ShortNodeID)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetNodeID() (r common.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of FullIntroductionReaderMock.GetNodeIDFunc invocations
func (m *FullIntroductionReaderMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of FullIntroductionReaderMock.GetNodeID invocations
func (m *FullIntroductionReaderMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetNodeIDFinished() bool {
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

type mFullIntroductionReaderMockGetNodePK struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetNodePKExpectation
	expectationSeries []*FullIntroductionReaderMockGetNodePKExpectation
}

type FullIntroductionReaderMockGetNodePKExpectation struct {
	result *FullIntroductionReaderMockGetNodePKResult
}

type FullIntroductionReaderMockGetNodePKResult struct {
	r common.SignatureKeyHolder
}

//Expect specifies that invocation of FullIntroductionReader.GetNodePK is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetNodePK) Expect() *mFullIntroductionReaderMockGetNodePK {
	m.mock.GetNodePKFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodePKExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetNodePK
func (m *mFullIntroductionReaderMockGetNodePK) Return(r common.SignatureKeyHolder) *FullIntroductionReaderMock {
	m.mock.GetNodePKFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodePKExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetNodePKResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetNodePK is expected once
func (m *mFullIntroductionReaderMockGetNodePK) ExpectOnce() *FullIntroductionReaderMockGetNodePKExpectation {
	m.mock.GetNodePKFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetNodePKExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetNodePKExpectation) Return(r common.SignatureKeyHolder) {
	e.result = &FullIntroductionReaderMockGetNodePKResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetNodePK method
func (m *mFullIntroductionReaderMockGetNodePK) Set(f func() (r common.SignatureKeyHolder)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePKFunc = f
	return m.mock
}

//GetNodePK implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetNodePK() (r common.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePKPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePKCounter, 1)

	if len(m.GetNodePKMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePKMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodePK.")
			return
		}

		result := m.GetNodePKMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodePK")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePKMock.mainExpectation != nil {

		result := m.GetNodePKMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodePK")
		}

		r = result.r

		return
	}

	if m.GetNodePKFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodePK.")
		return
	}

	return m.GetNodePKFunc()
}

//GetNodePKMinimockCounter returns a count of FullIntroductionReaderMock.GetNodePKFunc invocations
func (m *FullIntroductionReaderMock) GetNodePKMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePKCounter)
}

//GetNodePKMinimockPreCounter returns the value of FullIntroductionReaderMock.GetNodePK invocations
func (m *FullIntroductionReaderMock) GetNodePKMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePKPreCounter)
}

//GetNodePKFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetNodePKFinished() bool {
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

type mFullIntroductionReaderMockGetNodePrimaryRole struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetNodePrimaryRoleExpectation
	expectationSeries []*FullIntroductionReaderMockGetNodePrimaryRoleExpectation
}

type FullIntroductionReaderMockGetNodePrimaryRoleExpectation struct {
	result *FullIntroductionReaderMockGetNodePrimaryRoleResult
}

type FullIntroductionReaderMockGetNodePrimaryRoleResult struct {
	r common2.NodePrimaryRole
}

//Expect specifies that invocation of FullIntroductionReader.GetNodePrimaryRole is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetNodePrimaryRole) Expect() *mFullIntroductionReaderMockGetNodePrimaryRole {
	m.mock.GetNodePrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodePrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetNodePrimaryRole
func (m *mFullIntroductionReaderMockGetNodePrimaryRole) Return(r common2.NodePrimaryRole) *FullIntroductionReaderMock {
	m.mock.GetNodePrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodePrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetNodePrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetNodePrimaryRole is expected once
func (m *mFullIntroductionReaderMockGetNodePrimaryRole) ExpectOnce() *FullIntroductionReaderMockGetNodePrimaryRoleExpectation {
	m.mock.GetNodePrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetNodePrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetNodePrimaryRoleExpectation) Return(r common2.NodePrimaryRole) {
	e.result = &FullIntroductionReaderMockGetNodePrimaryRoleResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetNodePrimaryRole method
func (m *mFullIntroductionReaderMockGetNodePrimaryRole) Set(f func() (r common2.NodePrimaryRole)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePrimaryRoleFunc = f
	return m.mock
}

//GetNodePrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetNodePrimaryRole() (r common2.NodePrimaryRole) {
	counter := atomic.AddUint64(&m.GetNodePrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePrimaryRoleCounter, 1)

	if len(m.GetNodePrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodePrimaryRole.")
			return
		}

		result := m.GetNodePrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodePrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePrimaryRoleMock.mainExpectation != nil {

		result := m.GetNodePrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodePrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetNodePrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodePrimaryRole.")
		return
	}

	return m.GetNodePrimaryRoleFunc()
}

//GetNodePrimaryRoleMinimockCounter returns a count of FullIntroductionReaderMock.GetNodePrimaryRoleFunc invocations
func (m *FullIntroductionReaderMock) GetNodePrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePrimaryRoleCounter)
}

//GetNodePrimaryRoleMinimockPreCounter returns the value of FullIntroductionReaderMock.GetNodePrimaryRole invocations
func (m *FullIntroductionReaderMock) GetNodePrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePrimaryRolePreCounter)
}

//GetNodePrimaryRoleFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetNodePrimaryRoleFinished() bool {
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

type mFullIntroductionReaderMockGetNodeSpecialRoles struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetNodeSpecialRolesExpectation
	expectationSeries []*FullIntroductionReaderMockGetNodeSpecialRolesExpectation
}

type FullIntroductionReaderMockGetNodeSpecialRolesExpectation struct {
	result *FullIntroductionReaderMockGetNodeSpecialRolesResult
}

type FullIntroductionReaderMockGetNodeSpecialRolesResult struct {
	r common2.NodeSpecialRole
}

//Expect specifies that invocation of FullIntroductionReader.GetNodeSpecialRoles is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetNodeSpecialRoles) Expect() *mFullIntroductionReaderMockGetNodeSpecialRoles {
	m.mock.GetNodeSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodeSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetNodeSpecialRoles
func (m *mFullIntroductionReaderMockGetNodeSpecialRoles) Return(r common2.NodeSpecialRole) *FullIntroductionReaderMock {
	m.mock.GetNodeSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodeSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetNodeSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetNodeSpecialRoles is expected once
func (m *mFullIntroductionReaderMockGetNodeSpecialRoles) ExpectOnce() *FullIntroductionReaderMockGetNodeSpecialRolesExpectation {
	m.mock.GetNodeSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetNodeSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetNodeSpecialRolesExpectation) Return(r common2.NodeSpecialRole) {
	e.result = &FullIntroductionReaderMockGetNodeSpecialRolesResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetNodeSpecialRoles method
func (m *mFullIntroductionReaderMockGetNodeSpecialRoles) Set(f func() (r common2.NodeSpecialRole)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeSpecialRolesFunc = f
	return m.mock
}

//GetNodeSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetNodeSpecialRoles() (r common2.NodeSpecialRole) {
	counter := atomic.AddUint64(&m.GetNodeSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeSpecialRolesCounter, 1)

	if len(m.GetNodeSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodeSpecialRoles.")
			return
		}

		result := m.GetNodeSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodeSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeSpecialRolesMock.mainExpectation != nil {

		result := m.GetNodeSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodeSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetNodeSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodeSpecialRoles.")
		return
	}

	return m.GetNodeSpecialRolesFunc()
}

//GetNodeSpecialRolesMinimockCounter returns a count of FullIntroductionReaderMock.GetNodeSpecialRolesFunc invocations
func (m *FullIntroductionReaderMock) GetNodeSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeSpecialRolesCounter)
}

//GetNodeSpecialRolesMinimockPreCounter returns the value of FullIntroductionReaderMock.GetNodeSpecialRoles invocations
func (m *FullIntroductionReaderMock) GetNodeSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeSpecialRolesPreCounter)
}

//GetNodeSpecialRolesFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetNodeSpecialRolesFinished() bool {
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

type mFullIntroductionReaderMockGetPowerLevels struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetPowerLevelsExpectation
	expectationSeries []*FullIntroductionReaderMockGetPowerLevelsExpectation
}

type FullIntroductionReaderMockGetPowerLevelsExpectation struct {
	result *FullIntroductionReaderMockGetPowerLevelsResult
}

type FullIntroductionReaderMockGetPowerLevelsResult struct {
	r common2.MemberPowerSet
}

//Expect specifies that invocation of FullIntroductionReader.GetPowerLevels is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetPowerLevels) Expect() *mFullIntroductionReaderMockGetPowerLevels {
	m.mock.GetPowerLevelsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetPowerLevelsExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetPowerLevels
func (m *mFullIntroductionReaderMockGetPowerLevels) Return(r common2.MemberPowerSet) *FullIntroductionReaderMock {
	m.mock.GetPowerLevelsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetPowerLevelsExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetPowerLevelsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetPowerLevels is expected once
func (m *mFullIntroductionReaderMockGetPowerLevels) ExpectOnce() *FullIntroductionReaderMockGetPowerLevelsExpectation {
	m.mock.GetPowerLevelsFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetPowerLevelsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetPowerLevelsExpectation) Return(r common2.MemberPowerSet) {
	e.result = &FullIntroductionReaderMockGetPowerLevelsResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetPowerLevels method
func (m *mFullIntroductionReaderMockGetPowerLevels) Set(f func() (r common2.MemberPowerSet)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPowerLevelsFunc = f
	return m.mock
}

//GetPowerLevels implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetPowerLevels() (r common2.MemberPowerSet) {
	counter := atomic.AddUint64(&m.GetPowerLevelsPreCounter, 1)
	defer atomic.AddUint64(&m.GetPowerLevelsCounter, 1)

	if len(m.GetPowerLevelsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPowerLevelsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetPowerLevels.")
			return
		}

		result := m.GetPowerLevelsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetPowerLevels")
			return
		}

		r = result.r

		return
	}

	if m.GetPowerLevelsMock.mainExpectation != nil {

		result := m.GetPowerLevelsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetPowerLevels")
		}

		r = result.r

		return
	}

	if m.GetPowerLevelsFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetPowerLevels.")
		return
	}

	return m.GetPowerLevelsFunc()
}

//GetPowerLevelsMinimockCounter returns a count of FullIntroductionReaderMock.GetPowerLevelsFunc invocations
func (m *FullIntroductionReaderMock) GetPowerLevelsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPowerLevelsCounter)
}

//GetPowerLevelsMinimockPreCounter returns the value of FullIntroductionReaderMock.GetPowerLevels invocations
func (m *FullIntroductionReaderMock) GetPowerLevelsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPowerLevelsPreCounter)
}

//GetPowerLevelsFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetPowerLevelsFinished() bool {
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

type mFullIntroductionReaderMockGetReference struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetReferenceExpectation
	expectationSeries []*FullIntroductionReaderMockGetReferenceExpectation
}

type FullIntroductionReaderMockGetReferenceExpectation struct {
	result *FullIntroductionReaderMockGetReferenceResult
}

type FullIntroductionReaderMockGetReferenceResult struct {
	r insolar.Reference
}

//Expect specifies that invocation of FullIntroductionReader.GetReference is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetReference) Expect() *mFullIntroductionReaderMockGetReference {
	m.mock.GetReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetReferenceExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetReference
func (m *mFullIntroductionReaderMockGetReference) Return(r insolar.Reference) *FullIntroductionReaderMock {
	m.mock.GetReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetReferenceExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetReferenceResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetReference is expected once
func (m *mFullIntroductionReaderMockGetReference) ExpectOnce() *FullIntroductionReaderMockGetReferenceExpectation {
	m.mock.GetReferenceFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetReferenceExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetReferenceExpectation) Return(r insolar.Reference) {
	e.result = &FullIntroductionReaderMockGetReferenceResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetReference method
func (m *mFullIntroductionReaderMockGetReference) Set(f func() (r insolar.Reference)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetReferenceFunc = f
	return m.mock
}

//GetReference implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetReference() (r insolar.Reference) {
	counter := atomic.AddUint64(&m.GetReferencePreCounter, 1)
	defer atomic.AddUint64(&m.GetReferenceCounter, 1)

	if len(m.GetReferenceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetReferenceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetReference.")
			return
		}

		result := m.GetReferenceMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetReference")
			return
		}

		r = result.r

		return
	}

	if m.GetReferenceMock.mainExpectation != nil {

		result := m.GetReferenceMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetReference")
		}

		r = result.r

		return
	}

	if m.GetReferenceFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetReference.")
		return
	}

	return m.GetReferenceFunc()
}

//GetReferenceMinimockCounter returns a count of FullIntroductionReaderMock.GetReferenceFunc invocations
func (m *FullIntroductionReaderMock) GetReferenceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetReferenceCounter)
}

//GetReferenceMinimockPreCounter returns the value of FullIntroductionReaderMock.GetReference invocations
func (m *FullIntroductionReaderMock) GetReferenceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetReferencePreCounter)
}

//GetReferenceFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetReferenceFinished() bool {
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

type mFullIntroductionReaderMockGetStartPower struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetStartPowerExpectation
	expectationSeries []*FullIntroductionReaderMockGetStartPowerExpectation
}

type FullIntroductionReaderMockGetStartPowerExpectation struct {
	result *FullIntroductionReaderMockGetStartPowerResult
}

type FullIntroductionReaderMockGetStartPowerResult struct {
	r common2.MemberPower
}

//Expect specifies that invocation of FullIntroductionReader.GetStartPower is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetStartPower) Expect() *mFullIntroductionReaderMockGetStartPower {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetStartPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetStartPower
func (m *mFullIntroductionReaderMockGetStartPower) Return(r common2.MemberPower) *FullIntroductionReaderMock {
	m.mock.GetStartPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetStartPowerExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetStartPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetStartPower is expected once
func (m *mFullIntroductionReaderMockGetStartPower) ExpectOnce() *FullIntroductionReaderMockGetStartPowerExpectation {
	m.mock.GetStartPowerFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetStartPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetStartPowerExpectation) Return(r common2.MemberPower) {
	e.result = &FullIntroductionReaderMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetStartPower method
func (m *mFullIntroductionReaderMockGetStartPower) Set(f func() (r common2.MemberPower)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetStartPower() (r common2.MemberPower) {
	counter := atomic.AddUint64(&m.GetStartPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetStartPowerCounter, 1)

	if len(m.GetStartPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStartPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetStartPower.")
			return
		}

		result := m.GetStartPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetStartPower")
			return
		}

		r = result.r

		return
	}

	if m.GetStartPowerMock.mainExpectation != nil {

		result := m.GetStartPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetStartPower")
		}

		r = result.r

		return
	}

	if m.GetStartPowerFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetStartPower.")
		return
	}

	return m.GetStartPowerFunc()
}

//GetStartPowerMinimockCounter returns a count of FullIntroductionReaderMock.GetStartPowerFunc invocations
func (m *FullIntroductionReaderMock) GetStartPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerCounter)
}

//GetStartPowerMinimockPreCounter returns the value of FullIntroductionReaderMock.GetStartPower invocations
func (m *FullIntroductionReaderMock) GetStartPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStartPowerPreCounter)
}

//GetStartPowerFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetStartPowerFinished() bool {
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
func (m *FullIntroductionReaderMock) ValidateCallCounters() {

	if !m.GetExtraEndpointsFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetExtraEndpoints")
	}

	if !m.GetIssuedAtPulseFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuedAtPulse")
	}

	if !m.GetIssuedAtTimeFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuedAtTime")
	}

	if !m.GetIssuerIDFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuerID")
	}

	if !m.GetIssuerSignatureFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuerSignature")
	}

	if !m.GetJoinerSignatureFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetJoinerSignature")
	}

	if !m.GetNodeEndpointFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodeEndpoint")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodeID")
	}

	if !m.GetNodePKFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodePK")
	}

	if !m.GetNodePrimaryRoleFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodePrimaryRole")
	}

	if !m.GetNodeSpecialRolesFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodeSpecialRoles")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetPowerLevels")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetReference")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetStartPower")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FullIntroductionReaderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FullIntroductionReaderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FullIntroductionReaderMock) MinimockFinish() {

	if !m.GetExtraEndpointsFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetExtraEndpoints")
	}

	if !m.GetIssuedAtPulseFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuedAtPulse")
	}

	if !m.GetIssuedAtTimeFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuedAtTime")
	}

	if !m.GetIssuerIDFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuerID")
	}

	if !m.GetIssuerSignatureFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetIssuerSignature")
	}

	if !m.GetJoinerSignatureFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetJoinerSignature")
	}

	if !m.GetNodeEndpointFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodeEndpoint")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodeID")
	}

	if !m.GetNodePKFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodePK")
	}

	if !m.GetNodePrimaryRoleFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodePrimaryRole")
	}

	if !m.GetNodeSpecialRolesFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodeSpecialRoles")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetPowerLevels")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetReference")
	}

	if !m.GetStartPowerFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetStartPower")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FullIntroductionReaderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FullIntroductionReaderMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to FullIntroductionReaderMock.GetExtraEndpoints")
			}

			if !m.GetIssuedAtPulseFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetIssuedAtPulse")
			}

			if !m.GetIssuedAtTimeFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetIssuedAtTime")
			}

			if !m.GetIssuerIDFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetIssuerID")
			}

			if !m.GetIssuerSignatureFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetIssuerSignature")
			}

			if !m.GetJoinerSignatureFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetJoinerSignature")
			}

			if !m.GetNodeEndpointFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetNodeEndpoint")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetNodeID")
			}

			if !m.GetNodePKFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetNodePK")
			}

			if !m.GetNodePrimaryRoleFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetNodePrimaryRole")
			}

			if !m.GetNodeSpecialRolesFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetNodeSpecialRoles")
			}

			if !m.GetPowerLevelsFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetPowerLevels")
			}

			if !m.GetReferenceFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetReference")
			}

			if !m.GetStartPowerFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetStartPower")
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
func (m *FullIntroductionReaderMock) AllMocksCalled() bool {

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

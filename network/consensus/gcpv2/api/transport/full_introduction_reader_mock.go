package transport

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FullIntroductionReader" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/transport
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

//FullIntroductionReaderMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader
type FullIntroductionReaderMock struct {
	t minimock.Tester

	GetDefaultEndpointFunc       func() (r endpoints.Outbound)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mFullIntroductionReaderMockGetDefaultEndpoint

	GetExtraEndpointsFunc       func() (r []endpoints.Outbound)
	GetExtraEndpointsCounter    uint64
	GetExtraEndpointsPreCounter uint64
	GetExtraEndpointsMock       mFullIntroductionReaderMockGetExtraEndpoints

	GetIssuedAtPulseFunc       func() (r pulse.Number)
	GetIssuedAtPulseCounter    uint64
	GetIssuedAtPulsePreCounter uint64
	GetIssuedAtPulseMock       mFullIntroductionReaderMockGetIssuedAtPulse

	GetIssuedAtTimeFunc       func() (r time.Time)
	GetIssuedAtTimeCounter    uint64
	GetIssuedAtTimePreCounter uint64
	GetIssuedAtTimeMock       mFullIntroductionReaderMockGetIssuedAtTime

	GetIssuerIDFunc       func() (r insolar.ShortNodeID)
	GetIssuerIDCounter    uint64
	GetIssuerIDPreCounter uint64
	GetIssuerIDMock       mFullIntroductionReaderMockGetIssuerID

	GetIssuerSignatureFunc       func() (r cryptkit.SignatureHolder)
	GetIssuerSignatureCounter    uint64
	GetIssuerSignaturePreCounter uint64
	GetIssuerSignatureMock       mFullIntroductionReaderMockGetIssuerSignature

	GetJoinerSignatureFunc       func() (r cryptkit.SignatureHolder)
	GetJoinerSignatureCounter    uint64
	GetJoinerSignaturePreCounter uint64
	GetJoinerSignatureMock       mFullIntroductionReaderMockGetJoinerSignature

	GetNodePublicKeyFunc       func() (r cryptkit.SignatureKeyHolder)
	GetNodePublicKeyCounter    uint64
	GetNodePublicKeyPreCounter uint64
	GetNodePublicKeyMock       mFullIntroductionReaderMockGetNodePublicKey

	GetPowerLevelsFunc       func() (r member.PowerSet)
	GetPowerLevelsCounter    uint64
	GetPowerLevelsPreCounter uint64
	GetPowerLevelsMock       mFullIntroductionReaderMockGetPowerLevels

	GetPrimaryRoleFunc       func() (r member.PrimaryRole)
	GetPrimaryRoleCounter    uint64
	GetPrimaryRolePreCounter uint64
	GetPrimaryRoleMock       mFullIntroductionReaderMockGetPrimaryRole

	GetReferenceFunc       func() (r insolar.Reference)
	GetReferenceCounter    uint64
	GetReferencePreCounter uint64
	GetReferenceMock       mFullIntroductionReaderMockGetReference

	GetShortNodeIDFunc       func() (r insolar.ShortNodeID)
	GetShortNodeIDCounter    uint64
	GetShortNodeIDPreCounter uint64
	GetShortNodeIDMock       mFullIntroductionReaderMockGetShortNodeID

	GetSpecialRolesFunc       func() (r member.SpecialRole)
	GetSpecialRolesCounter    uint64
	GetSpecialRolesPreCounter uint64
	GetSpecialRolesMock       mFullIntroductionReaderMockGetSpecialRoles

	GetStartPowerFunc       func() (r member.Power)
	GetStartPowerCounter    uint64
	GetStartPowerPreCounter uint64
	GetStartPowerMock       mFullIntroductionReaderMockGetStartPower
}

//NewFullIntroductionReaderMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader
func NewFullIntroductionReaderMock(t minimock.Tester) *FullIntroductionReaderMock {
	m := &FullIntroductionReaderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetDefaultEndpointMock = mFullIntroductionReaderMockGetDefaultEndpoint{mock: m}
	m.GetExtraEndpointsMock = mFullIntroductionReaderMockGetExtraEndpoints{mock: m}
	m.GetIssuedAtPulseMock = mFullIntroductionReaderMockGetIssuedAtPulse{mock: m}
	m.GetIssuedAtTimeMock = mFullIntroductionReaderMockGetIssuedAtTime{mock: m}
	m.GetIssuerIDMock = mFullIntroductionReaderMockGetIssuerID{mock: m}
	m.GetIssuerSignatureMock = mFullIntroductionReaderMockGetIssuerSignature{mock: m}
	m.GetJoinerSignatureMock = mFullIntroductionReaderMockGetJoinerSignature{mock: m}
	m.GetNodePublicKeyMock = mFullIntroductionReaderMockGetNodePublicKey{mock: m}
	m.GetPowerLevelsMock = mFullIntroductionReaderMockGetPowerLevels{mock: m}
	m.GetPrimaryRoleMock = mFullIntroductionReaderMockGetPrimaryRole{mock: m}
	m.GetReferenceMock = mFullIntroductionReaderMockGetReference{mock: m}
	m.GetShortNodeIDMock = mFullIntroductionReaderMockGetShortNodeID{mock: m}
	m.GetSpecialRolesMock = mFullIntroductionReaderMockGetSpecialRoles{mock: m}
	m.GetStartPowerMock = mFullIntroductionReaderMockGetStartPower{mock: m}

	return m
}

type mFullIntroductionReaderMockGetDefaultEndpoint struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetDefaultEndpointExpectation
	expectationSeries []*FullIntroductionReaderMockGetDefaultEndpointExpectation
}

type FullIntroductionReaderMockGetDefaultEndpointExpectation struct {
	result *FullIntroductionReaderMockGetDefaultEndpointResult
}

type FullIntroductionReaderMockGetDefaultEndpointResult struct {
	r endpoints.Outbound
}

//Expect specifies that invocation of FullIntroductionReader.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetDefaultEndpoint) Expect() *mFullIntroductionReaderMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetDefaultEndpoint
func (m *mFullIntroductionReaderMockGetDefaultEndpoint) Return(r endpoints.Outbound) *FullIntroductionReaderMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetDefaultEndpoint is expected once
func (m *mFullIntroductionReaderMockGetDefaultEndpoint) ExpectOnce() *FullIntroductionReaderMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetDefaultEndpointExpectation) Return(r endpoints.Outbound) {
	e.result = &FullIntroductionReaderMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetDefaultEndpoint method
func (m *mFullIntroductionReaderMockGetDefaultEndpoint) Set(f func() (r endpoints.Outbound)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetDefaultEndpoint() (r endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of FullIntroductionReaderMock.GetDefaultEndpointFunc invocations
func (m *FullIntroductionReaderMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of FullIntroductionReaderMock.GetDefaultEndpoint invocations
func (m *FullIntroductionReaderMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetDefaultEndpointFinished() bool {
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

type mFullIntroductionReaderMockGetExtraEndpoints struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetExtraEndpointsExpectation
	expectationSeries []*FullIntroductionReaderMockGetExtraEndpointsExpectation
}

type FullIntroductionReaderMockGetExtraEndpointsExpectation struct {
	result *FullIntroductionReaderMockGetExtraEndpointsResult
}

type FullIntroductionReaderMockGetExtraEndpointsResult struct {
	r []endpoints.Outbound
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
func (m *mFullIntroductionReaderMockGetExtraEndpoints) Return(r []endpoints.Outbound) *FullIntroductionReaderMock {
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

func (e *FullIntroductionReaderMockGetExtraEndpointsExpectation) Return(r []endpoints.Outbound) {
	e.result = &FullIntroductionReaderMockGetExtraEndpointsResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetExtraEndpoints method
func (m *mFullIntroductionReaderMockGetExtraEndpoints) Set(f func() (r []endpoints.Outbound)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExtraEndpointsFunc = f
	return m.mock
}

//GetExtraEndpoints implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetExtraEndpoints() (r []endpoints.Outbound) {
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
	r pulse.Number
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
func (m *mFullIntroductionReaderMockGetIssuedAtPulse) Return(r pulse.Number) *FullIntroductionReaderMock {
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

func (e *FullIntroductionReaderMockGetIssuedAtPulseExpectation) Return(r pulse.Number) {
	e.result = &FullIntroductionReaderMockGetIssuedAtPulseResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetIssuedAtPulse method
func (m *mFullIntroductionReaderMockGetIssuedAtPulse) Set(f func() (r pulse.Number)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtPulseFunc = f
	return m.mock
}

//GetIssuedAtPulse implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetIssuedAtPulse() (r pulse.Number) {
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

//GetIssuedAtTime implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
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
	r insolar.ShortNodeID
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
func (m *mFullIntroductionReaderMockGetIssuerID) Return(r insolar.ShortNodeID) *FullIntroductionReaderMock {
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

func (e *FullIntroductionReaderMockGetIssuerIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &FullIntroductionReaderMockGetIssuerIDResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetIssuerID method
func (m *mFullIntroductionReaderMockGetIssuerID) Set(f func() (r insolar.ShortNodeID)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerIDFunc = f
	return m.mock
}

//GetIssuerID implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetIssuerID() (r insolar.ShortNodeID) {
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
	r cryptkit.SignatureHolder
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
func (m *mFullIntroductionReaderMockGetIssuerSignature) Return(r cryptkit.SignatureHolder) *FullIntroductionReaderMock {
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

func (e *FullIntroductionReaderMockGetIssuerSignatureExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &FullIntroductionReaderMockGetIssuerSignatureResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetIssuerSignature method
func (m *mFullIntroductionReaderMockGetIssuerSignature) Set(f func() (r cryptkit.SignatureHolder)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerSignatureFunc = f
	return m.mock
}

//GetIssuerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetIssuerSignature() (r cryptkit.SignatureHolder) {
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
	r cryptkit.SignatureHolder
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
func (m *mFullIntroductionReaderMockGetJoinerSignature) Return(r cryptkit.SignatureHolder) *FullIntroductionReaderMock {
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

func (e *FullIntroductionReaderMockGetJoinerSignatureExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &FullIntroductionReaderMockGetJoinerSignatureResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetJoinerSignature method
func (m *mFullIntroductionReaderMockGetJoinerSignature) Set(f func() (r cryptkit.SignatureHolder)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetJoinerSignatureFunc = f
	return m.mock
}

//GetJoinerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetJoinerSignature() (r cryptkit.SignatureHolder) {
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

type mFullIntroductionReaderMockGetNodePublicKey struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetNodePublicKeyExpectation
	expectationSeries []*FullIntroductionReaderMockGetNodePublicKeyExpectation
}

type FullIntroductionReaderMockGetNodePublicKeyExpectation struct {
	result *FullIntroductionReaderMockGetNodePublicKeyResult
}

type FullIntroductionReaderMockGetNodePublicKeyResult struct {
	r cryptkit.SignatureKeyHolder
}

//Expect specifies that invocation of FullIntroductionReader.GetNodePublicKey is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetNodePublicKey) Expect() *mFullIntroductionReaderMockGetNodePublicKey {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodePublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetNodePublicKey
func (m *mFullIntroductionReaderMockGetNodePublicKey) Return(r cryptkit.SignatureKeyHolder) *FullIntroductionReaderMock {
	m.mock.GetNodePublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetNodePublicKeyExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetNodePublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetNodePublicKey is expected once
func (m *mFullIntroductionReaderMockGetNodePublicKey) ExpectOnce() *FullIntroductionReaderMockGetNodePublicKeyExpectation {
	m.mock.GetNodePublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetNodePublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetNodePublicKeyExpectation) Return(r cryptkit.SignatureKeyHolder) {
	e.result = &FullIntroductionReaderMockGetNodePublicKeyResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetNodePublicKey method
func (m *mFullIntroductionReaderMockGetNodePublicKey) Set(f func() (r cryptkit.SignatureKeyHolder)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyFunc = f
	return m.mock
}

//GetNodePublicKey implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetNodePublicKey() (r cryptkit.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyCounter, 1)

	if len(m.GetNodePublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodePublicKey.")
			return
		}

		result := m.GetNodePublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodePublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyMock.mainExpectation != nil {

		result := m.GetNodePublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetNodePublicKey")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetNodePublicKey.")
		return
	}

	return m.GetNodePublicKeyFunc()
}

//GetNodePublicKeyMinimockCounter returns a count of FullIntroductionReaderMock.GetNodePublicKeyFunc invocations
func (m *FullIntroductionReaderMock) GetNodePublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyCounter)
}

//GetNodePublicKeyMinimockPreCounter returns the value of FullIntroductionReaderMock.GetNodePublicKey invocations
func (m *FullIntroductionReaderMock) GetNodePublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyPreCounter)
}

//GetNodePublicKeyFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetNodePublicKeyFinished() bool {
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

type mFullIntroductionReaderMockGetPowerLevels struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetPowerLevelsExpectation
	expectationSeries []*FullIntroductionReaderMockGetPowerLevelsExpectation
}

type FullIntroductionReaderMockGetPowerLevelsExpectation struct {
	result *FullIntroductionReaderMockGetPowerLevelsResult
}

type FullIntroductionReaderMockGetPowerLevelsResult struct {
	r member.PowerSet
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
func (m *mFullIntroductionReaderMockGetPowerLevels) Return(r member.PowerSet) *FullIntroductionReaderMock {
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

func (e *FullIntroductionReaderMockGetPowerLevelsExpectation) Return(r member.PowerSet) {
	e.result = &FullIntroductionReaderMockGetPowerLevelsResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetPowerLevels method
func (m *mFullIntroductionReaderMockGetPowerLevels) Set(f func() (r member.PowerSet)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPowerLevelsFunc = f
	return m.mock
}

//GetPowerLevels implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetPowerLevels() (r member.PowerSet) {
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

type mFullIntroductionReaderMockGetPrimaryRole struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetPrimaryRoleExpectation
	expectationSeries []*FullIntroductionReaderMockGetPrimaryRoleExpectation
}

type FullIntroductionReaderMockGetPrimaryRoleExpectation struct {
	result *FullIntroductionReaderMockGetPrimaryRoleResult
}

type FullIntroductionReaderMockGetPrimaryRoleResult struct {
	r member.PrimaryRole
}

//Expect specifies that invocation of FullIntroductionReader.GetPrimaryRole is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetPrimaryRole) Expect() *mFullIntroductionReaderMockGetPrimaryRole {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetPrimaryRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetPrimaryRole
func (m *mFullIntroductionReaderMockGetPrimaryRole) Return(r member.PrimaryRole) *FullIntroductionReaderMock {
	m.mock.GetPrimaryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetPrimaryRoleExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetPrimaryRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetPrimaryRole is expected once
func (m *mFullIntroductionReaderMockGetPrimaryRole) ExpectOnce() *FullIntroductionReaderMockGetPrimaryRoleExpectation {
	m.mock.GetPrimaryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetPrimaryRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetPrimaryRoleExpectation) Return(r member.PrimaryRole) {
	e.result = &FullIntroductionReaderMockGetPrimaryRoleResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetPrimaryRole method
func (m *mFullIntroductionReaderMockGetPrimaryRole) Set(f func() (r member.PrimaryRole)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimaryRoleFunc = f
	return m.mock
}

//GetPrimaryRole implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetPrimaryRole() (r member.PrimaryRole) {
	counter := atomic.AddUint64(&m.GetPrimaryRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimaryRoleCounter, 1)

	if len(m.GetPrimaryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimaryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetPrimaryRole.")
			return
		}

		result := m.GetPrimaryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetPrimaryRole")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleMock.mainExpectation != nil {

		result := m.GetPrimaryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetPrimaryRole")
		}

		r = result.r

		return
	}

	if m.GetPrimaryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetPrimaryRole.")
		return
	}

	return m.GetPrimaryRoleFunc()
}

//GetPrimaryRoleMinimockCounter returns a count of FullIntroductionReaderMock.GetPrimaryRoleFunc invocations
func (m *FullIntroductionReaderMock) GetPrimaryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRoleCounter)
}

//GetPrimaryRoleMinimockPreCounter returns the value of FullIntroductionReaderMock.GetPrimaryRole invocations
func (m *FullIntroductionReaderMock) GetPrimaryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimaryRolePreCounter)
}

//GetPrimaryRoleFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetPrimaryRoleFinished() bool {
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

//GetReference implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
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

type mFullIntroductionReaderMockGetShortNodeID struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetShortNodeIDExpectation
	expectationSeries []*FullIntroductionReaderMockGetShortNodeIDExpectation
}

type FullIntroductionReaderMockGetShortNodeIDExpectation struct {
	result *FullIntroductionReaderMockGetShortNodeIDResult
}

type FullIntroductionReaderMockGetShortNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of FullIntroductionReader.GetShortNodeID is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetShortNodeID) Expect() *mFullIntroductionReaderMockGetShortNodeID {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetShortNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetShortNodeID
func (m *mFullIntroductionReaderMockGetShortNodeID) Return(r insolar.ShortNodeID) *FullIntroductionReaderMock {
	m.mock.GetShortNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetShortNodeIDExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetShortNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetShortNodeID is expected once
func (m *mFullIntroductionReaderMockGetShortNodeID) ExpectOnce() *FullIntroductionReaderMockGetShortNodeIDExpectation {
	m.mock.GetShortNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetShortNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetShortNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &FullIntroductionReaderMockGetShortNodeIDResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetShortNodeID method
func (m *mFullIntroductionReaderMockGetShortNodeID) Set(f func() (r insolar.ShortNodeID)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetShortNodeIDFunc = f
	return m.mock
}

//GetShortNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetShortNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetShortNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetShortNodeIDCounter, 1)

	if len(m.GetShortNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetShortNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetShortNodeID.")
			return
		}

		result := m.GetShortNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetShortNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDMock.mainExpectation != nil {

		result := m.GetShortNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetShortNodeID")
		}

		r = result.r

		return
	}

	if m.GetShortNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetShortNodeID.")
		return
	}

	return m.GetShortNodeIDFunc()
}

//GetShortNodeIDMinimockCounter returns a count of FullIntroductionReaderMock.GetShortNodeIDFunc invocations
func (m *FullIntroductionReaderMock) GetShortNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDCounter)
}

//GetShortNodeIDMinimockPreCounter returns the value of FullIntroductionReaderMock.GetShortNodeID invocations
func (m *FullIntroductionReaderMock) GetShortNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetShortNodeIDPreCounter)
}

//GetShortNodeIDFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetShortNodeIDFinished() bool {
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

type mFullIntroductionReaderMockGetSpecialRoles struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetSpecialRolesExpectation
	expectationSeries []*FullIntroductionReaderMockGetSpecialRolesExpectation
}

type FullIntroductionReaderMockGetSpecialRolesExpectation struct {
	result *FullIntroductionReaderMockGetSpecialRolesResult
}

type FullIntroductionReaderMockGetSpecialRolesResult struct {
	r member.SpecialRole
}

//Expect specifies that invocation of FullIntroductionReader.GetSpecialRoles is expected from 1 to Infinity times
func (m *mFullIntroductionReaderMockGetSpecialRoles) Expect() *mFullIntroductionReaderMockGetSpecialRoles {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetSpecialRolesExpectation{}
	}

	return m
}

//Return specifies results of invocation of FullIntroductionReader.GetSpecialRoles
func (m *mFullIntroductionReaderMockGetSpecialRoles) Return(r member.SpecialRole) *FullIntroductionReaderMock {
	m.mock.GetSpecialRolesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FullIntroductionReaderMockGetSpecialRolesExpectation{}
	}
	m.mainExpectation.result = &FullIntroductionReaderMockGetSpecialRolesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FullIntroductionReader.GetSpecialRoles is expected once
func (m *mFullIntroductionReaderMockGetSpecialRoles) ExpectOnce() *FullIntroductionReaderMockGetSpecialRolesExpectation {
	m.mock.GetSpecialRolesFunc = nil
	m.mainExpectation = nil

	expectation := &FullIntroductionReaderMockGetSpecialRolesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FullIntroductionReaderMockGetSpecialRolesExpectation) Return(r member.SpecialRole) {
	e.result = &FullIntroductionReaderMockGetSpecialRolesResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetSpecialRoles method
func (m *mFullIntroductionReaderMockGetSpecialRoles) Set(f func() (r member.SpecialRole)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSpecialRolesFunc = f
	return m.mock
}

//GetSpecialRoles implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetSpecialRoles() (r member.SpecialRole) {
	counter := atomic.AddUint64(&m.GetSpecialRolesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSpecialRolesCounter, 1)

	if len(m.GetSpecialRolesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSpecialRolesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetSpecialRoles.")
			return
		}

		result := m.GetSpecialRolesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetSpecialRoles")
			return
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesMock.mainExpectation != nil {

		result := m.GetSpecialRolesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FullIntroductionReaderMock.GetSpecialRoles")
		}

		r = result.r

		return
	}

	if m.GetSpecialRolesFunc == nil {
		m.t.Fatalf("Unexpected call to FullIntroductionReaderMock.GetSpecialRoles.")
		return
	}

	return m.GetSpecialRolesFunc()
}

//GetSpecialRolesMinimockCounter returns a count of FullIntroductionReaderMock.GetSpecialRolesFunc invocations
func (m *FullIntroductionReaderMock) GetSpecialRolesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesCounter)
}

//GetSpecialRolesMinimockPreCounter returns the value of FullIntroductionReaderMock.GetSpecialRoles invocations
func (m *FullIntroductionReaderMock) GetSpecialRolesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSpecialRolesPreCounter)
}

//GetSpecialRolesFinished returns true if mock invocations count is ok
func (m *FullIntroductionReaderMock) GetSpecialRolesFinished() bool {
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

type mFullIntroductionReaderMockGetStartPower struct {
	mock              *FullIntroductionReaderMock
	mainExpectation   *FullIntroductionReaderMockGetStartPowerExpectation
	expectationSeries []*FullIntroductionReaderMockGetStartPowerExpectation
}

type FullIntroductionReaderMockGetStartPowerExpectation struct {
	result *FullIntroductionReaderMockGetStartPowerResult
}

type FullIntroductionReaderMockGetStartPowerResult struct {
	r member.Power
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
func (m *mFullIntroductionReaderMockGetStartPower) Return(r member.Power) *FullIntroductionReaderMock {
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

func (e *FullIntroductionReaderMockGetStartPowerExpectation) Return(r member.Power) {
	e.result = &FullIntroductionReaderMockGetStartPowerResult{r}
}

//Set uses given function f as a mock of FullIntroductionReader.GetStartPower method
func (m *mFullIntroductionReaderMockGetStartPower) Set(f func() (r member.Power)) *FullIntroductionReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStartPowerFunc = f
	return m.mock
}

//GetStartPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader interface
func (m *FullIntroductionReaderMock) GetStartPower() (r member.Power) {
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

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetDefaultEndpoint")
	}

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

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodePublicKey")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetPowerLevels")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetPrimaryRole")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetReference")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetShortNodeID")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetSpecialRoles")
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

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetDefaultEndpoint")
	}

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

	if !m.GetNodePublicKeyFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetNodePublicKey")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetPowerLevels")
	}

	if !m.GetPrimaryRoleFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetPrimaryRole")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetReference")
	}

	if !m.GetShortNodeIDFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetShortNodeID")
	}

	if !m.GetSpecialRolesFinished() {
		m.t.Fatal("Expected call to FullIntroductionReaderMock.GetSpecialRoles")
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
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetExtraEndpointsFinished()
		ok = ok && m.GetIssuedAtPulseFinished()
		ok = ok && m.GetIssuedAtTimeFinished()
		ok = ok && m.GetIssuerIDFinished()
		ok = ok && m.GetIssuerSignatureFinished()
		ok = ok && m.GetJoinerSignatureFinished()
		ok = ok && m.GetNodePublicKeyFinished()
		ok = ok && m.GetPowerLevelsFinished()
		ok = ok && m.GetPrimaryRoleFinished()
		ok = ok && m.GetReferenceFinished()
		ok = ok && m.GetShortNodeIDFinished()
		ok = ok && m.GetSpecialRolesFinished()
		ok = ok && m.GetStartPowerFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetDefaultEndpoint")
			}

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

			if !m.GetNodePublicKeyFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetNodePublicKey")
			}

			if !m.GetPowerLevelsFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetPowerLevels")
			}

			if !m.GetPrimaryRoleFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetPrimaryRole")
			}

			if !m.GetReferenceFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetReference")
			}

			if !m.GetShortNodeIDFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetShortNodeID")
			}

			if !m.GetSpecialRolesFinished() {
				m.t.Error("Expected call to FullIntroductionReaderMock.GetSpecialRoles")
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

	if !m.GetJoinerSignatureFinished() {
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

	if !m.GetShortNodeIDFinished() {
		return false
	}

	if !m.GetSpecialRolesFinished() {
		return false
	}

	if !m.GetStartPowerFinished() {
		return false
	}

	return true
}

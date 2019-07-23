package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StaticProfileExtension" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
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

//StaticProfileExtensionMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension
type StaticProfileExtensionMock struct {
	t minimock.Tester

	GetExtraEndpointsFunc       func() (r []endpoints.Outbound)
	GetExtraEndpointsCounter    uint64
	GetExtraEndpointsPreCounter uint64
	GetExtraEndpointsMock       mStaticProfileExtensionMockGetExtraEndpoints

	GetIntroducedNodeIDFunc       func() (r insolar.ShortNodeID)
	GetIntroducedNodeIDCounter    uint64
	GetIntroducedNodeIDPreCounter uint64
	GetIntroducedNodeIDMock       mStaticProfileExtensionMockGetIntroducedNodeID

	GetIssuedAtPulseFunc       func() (r pulse.Number)
	GetIssuedAtPulseCounter    uint64
	GetIssuedAtPulsePreCounter uint64
	GetIssuedAtPulseMock       mStaticProfileExtensionMockGetIssuedAtPulse

	GetIssuedAtTimeFunc       func() (r time.Time)
	GetIssuedAtTimeCounter    uint64
	GetIssuedAtTimePreCounter uint64
	GetIssuedAtTimeMock       mStaticProfileExtensionMockGetIssuedAtTime

	GetIssuerIDFunc       func() (r insolar.ShortNodeID)
	GetIssuerIDCounter    uint64
	GetIssuerIDPreCounter uint64
	GetIssuerIDMock       mStaticProfileExtensionMockGetIssuerID

	GetIssuerSignatureFunc       func() (r cryptkit.SignatureHolder)
	GetIssuerSignatureCounter    uint64
	GetIssuerSignaturePreCounter uint64
	GetIssuerSignatureMock       mStaticProfileExtensionMockGetIssuerSignature

	GetPowerLevelsFunc       func() (r member.PowerSet)
	GetPowerLevelsCounter    uint64
	GetPowerLevelsPreCounter uint64
	GetPowerLevelsMock       mStaticProfileExtensionMockGetPowerLevels

	GetReferenceFunc       func() (r insolar.Reference)
	GetReferenceCounter    uint64
	GetReferencePreCounter uint64
	GetReferenceMock       mStaticProfileExtensionMockGetReference
}

//NewStaticProfileExtensionMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension
func NewStaticProfileExtensionMock(t minimock.Tester) *StaticProfileExtensionMock {
	m := &StaticProfileExtensionMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetExtraEndpointsMock = mStaticProfileExtensionMockGetExtraEndpoints{mock: m}
	m.GetIntroducedNodeIDMock = mStaticProfileExtensionMockGetIntroducedNodeID{mock: m}
	m.GetIssuedAtPulseMock = mStaticProfileExtensionMockGetIssuedAtPulse{mock: m}
	m.GetIssuedAtTimeMock = mStaticProfileExtensionMockGetIssuedAtTime{mock: m}
	m.GetIssuerIDMock = mStaticProfileExtensionMockGetIssuerID{mock: m}
	m.GetIssuerSignatureMock = mStaticProfileExtensionMockGetIssuerSignature{mock: m}
	m.GetPowerLevelsMock = mStaticProfileExtensionMockGetPowerLevels{mock: m}
	m.GetReferenceMock = mStaticProfileExtensionMockGetReference{mock: m}

	return m
}

type mStaticProfileExtensionMockGetExtraEndpoints struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetExtraEndpointsExpectation
	expectationSeries []*StaticProfileExtensionMockGetExtraEndpointsExpectation
}

type StaticProfileExtensionMockGetExtraEndpointsExpectation struct {
	result *StaticProfileExtensionMockGetExtraEndpointsResult
}

type StaticProfileExtensionMockGetExtraEndpointsResult struct {
	r []endpoints.Outbound
}

//Expect specifies that invocation of StaticProfileExtension.GetExtraEndpoints is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetExtraEndpoints) Expect() *mStaticProfileExtensionMockGetExtraEndpoints {
	m.mock.GetExtraEndpointsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetExtraEndpointsExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetExtraEndpoints
func (m *mStaticProfileExtensionMockGetExtraEndpoints) Return(r []endpoints.Outbound) *StaticProfileExtensionMock {
	m.mock.GetExtraEndpointsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetExtraEndpointsExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetExtraEndpointsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetExtraEndpoints is expected once
func (m *mStaticProfileExtensionMockGetExtraEndpoints) ExpectOnce() *StaticProfileExtensionMockGetExtraEndpointsExpectation {
	m.mock.GetExtraEndpointsFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetExtraEndpointsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetExtraEndpointsExpectation) Return(r []endpoints.Outbound) {
	e.result = &StaticProfileExtensionMockGetExtraEndpointsResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetExtraEndpoints method
func (m *mStaticProfileExtensionMockGetExtraEndpoints) Set(f func() (r []endpoints.Outbound)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExtraEndpointsFunc = f
	return m.mock
}

//GetExtraEndpoints implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetExtraEndpoints() (r []endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetExtraEndpointsPreCounter, 1)
	defer atomic.AddUint64(&m.GetExtraEndpointsCounter, 1)

	if len(m.GetExtraEndpointsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExtraEndpointsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetExtraEndpoints.")
			return
		}

		result := m.GetExtraEndpointsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetExtraEndpoints")
			return
		}

		r = result.r

		return
	}

	if m.GetExtraEndpointsMock.mainExpectation != nil {

		result := m.GetExtraEndpointsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetExtraEndpoints")
		}

		r = result.r

		return
	}

	if m.GetExtraEndpointsFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetExtraEndpoints.")
		return
	}

	return m.GetExtraEndpointsFunc()
}

//GetExtraEndpointsMinimockCounter returns a count of StaticProfileExtensionMock.GetExtraEndpointsFunc invocations
func (m *StaticProfileExtensionMock) GetExtraEndpointsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtraEndpointsCounter)
}

//GetExtraEndpointsMinimockPreCounter returns the value of StaticProfileExtensionMock.GetExtraEndpoints invocations
func (m *StaticProfileExtensionMock) GetExtraEndpointsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExtraEndpointsPreCounter)
}

//GetExtraEndpointsFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetExtraEndpointsFinished() bool {
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

type mStaticProfileExtensionMockGetIntroducedNodeID struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetIntroducedNodeIDExpectation
	expectationSeries []*StaticProfileExtensionMockGetIntroducedNodeIDExpectation
}

type StaticProfileExtensionMockGetIntroducedNodeIDExpectation struct {
	result *StaticProfileExtensionMockGetIntroducedNodeIDResult
}

type StaticProfileExtensionMockGetIntroducedNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of StaticProfileExtension.GetIntroducedNodeID is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetIntroducedNodeID) Expect() *mStaticProfileExtensionMockGetIntroducedNodeID {
	m.mock.GetIntroducedNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIntroducedNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetIntroducedNodeID
func (m *mStaticProfileExtensionMockGetIntroducedNodeID) Return(r insolar.ShortNodeID) *StaticProfileExtensionMock {
	m.mock.GetIntroducedNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIntroducedNodeIDExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetIntroducedNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetIntroducedNodeID is expected once
func (m *mStaticProfileExtensionMockGetIntroducedNodeID) ExpectOnce() *StaticProfileExtensionMockGetIntroducedNodeIDExpectation {
	m.mock.GetIntroducedNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetIntroducedNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetIntroducedNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &StaticProfileExtensionMockGetIntroducedNodeIDResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetIntroducedNodeID method
func (m *mStaticProfileExtensionMockGetIntroducedNodeID) Set(f func() (r insolar.ShortNodeID)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIntroducedNodeIDFunc = f
	return m.mock
}

//GetIntroducedNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetIntroducedNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetIntroducedNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetIntroducedNodeIDCounter, 1)

	if len(m.GetIntroducedNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIntroducedNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIntroducedNodeID.")
			return
		}

		result := m.GetIntroducedNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIntroducedNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetIntroducedNodeIDMock.mainExpectation != nil {

		result := m.GetIntroducedNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIntroducedNodeID")
		}

		r = result.r

		return
	}

	if m.GetIntroducedNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIntroducedNodeID.")
		return
	}

	return m.GetIntroducedNodeIDFunc()
}

//GetIntroducedNodeIDMinimockCounter returns a count of StaticProfileExtensionMock.GetIntroducedNodeIDFunc invocations
func (m *StaticProfileExtensionMock) GetIntroducedNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroducedNodeIDCounter)
}

//GetIntroducedNodeIDMinimockPreCounter returns the value of StaticProfileExtensionMock.GetIntroducedNodeID invocations
func (m *StaticProfileExtensionMock) GetIntroducedNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIntroducedNodeIDPreCounter)
}

//GetIntroducedNodeIDFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetIntroducedNodeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIntroducedNodeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIntroducedNodeIDCounter) == uint64(len(m.GetIntroducedNodeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIntroducedNodeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIntroducedNodeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIntroducedNodeIDFunc != nil {
		return atomic.LoadUint64(&m.GetIntroducedNodeIDCounter) > 0
	}

	return true
}

type mStaticProfileExtensionMockGetIssuedAtPulse struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetIssuedAtPulseExpectation
	expectationSeries []*StaticProfileExtensionMockGetIssuedAtPulseExpectation
}

type StaticProfileExtensionMockGetIssuedAtPulseExpectation struct {
	result *StaticProfileExtensionMockGetIssuedAtPulseResult
}

type StaticProfileExtensionMockGetIssuedAtPulseResult struct {
	r pulse.Number
}

//Expect specifies that invocation of StaticProfileExtension.GetIssuedAtPulse is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetIssuedAtPulse) Expect() *mStaticProfileExtensionMockGetIssuedAtPulse {
	m.mock.GetIssuedAtPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuedAtPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetIssuedAtPulse
func (m *mStaticProfileExtensionMockGetIssuedAtPulse) Return(r pulse.Number) *StaticProfileExtensionMock {
	m.mock.GetIssuedAtPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuedAtPulseExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetIssuedAtPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetIssuedAtPulse is expected once
func (m *mStaticProfileExtensionMockGetIssuedAtPulse) ExpectOnce() *StaticProfileExtensionMockGetIssuedAtPulseExpectation {
	m.mock.GetIssuedAtPulseFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetIssuedAtPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetIssuedAtPulseExpectation) Return(r pulse.Number) {
	e.result = &StaticProfileExtensionMockGetIssuedAtPulseResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetIssuedAtPulse method
func (m *mStaticProfileExtensionMockGetIssuedAtPulse) Set(f func() (r pulse.Number)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtPulseFunc = f
	return m.mock
}

//GetIssuedAtPulse implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetIssuedAtPulse() (r pulse.Number) {
	counter := atomic.AddUint64(&m.GetIssuedAtPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuedAtPulseCounter, 1)

	if len(m.GetIssuedAtPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuedAtPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuedAtPulse.")
			return
		}

		result := m.GetIssuedAtPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuedAtPulse")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuedAtPulseMock.mainExpectation != nil {

		result := m.GetIssuedAtPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuedAtPulse")
		}

		r = result.r

		return
	}

	if m.GetIssuedAtPulseFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuedAtPulse.")
		return
	}

	return m.GetIssuedAtPulseFunc()
}

//GetIssuedAtPulseMinimockCounter returns a count of StaticProfileExtensionMock.GetIssuedAtPulseFunc invocations
func (m *StaticProfileExtensionMock) GetIssuedAtPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtPulseCounter)
}

//GetIssuedAtPulseMinimockPreCounter returns the value of StaticProfileExtensionMock.GetIssuedAtPulse invocations
func (m *StaticProfileExtensionMock) GetIssuedAtPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtPulsePreCounter)
}

//GetIssuedAtPulseFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetIssuedAtPulseFinished() bool {
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

type mStaticProfileExtensionMockGetIssuedAtTime struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetIssuedAtTimeExpectation
	expectationSeries []*StaticProfileExtensionMockGetIssuedAtTimeExpectation
}

type StaticProfileExtensionMockGetIssuedAtTimeExpectation struct {
	result *StaticProfileExtensionMockGetIssuedAtTimeResult
}

type StaticProfileExtensionMockGetIssuedAtTimeResult struct {
	r time.Time
}

//Expect specifies that invocation of StaticProfileExtension.GetIssuedAtTime is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetIssuedAtTime) Expect() *mStaticProfileExtensionMockGetIssuedAtTime {
	m.mock.GetIssuedAtTimeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuedAtTimeExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetIssuedAtTime
func (m *mStaticProfileExtensionMockGetIssuedAtTime) Return(r time.Time) *StaticProfileExtensionMock {
	m.mock.GetIssuedAtTimeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuedAtTimeExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetIssuedAtTimeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetIssuedAtTime is expected once
func (m *mStaticProfileExtensionMockGetIssuedAtTime) ExpectOnce() *StaticProfileExtensionMockGetIssuedAtTimeExpectation {
	m.mock.GetIssuedAtTimeFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetIssuedAtTimeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetIssuedAtTimeExpectation) Return(r time.Time) {
	e.result = &StaticProfileExtensionMockGetIssuedAtTimeResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetIssuedAtTime method
func (m *mStaticProfileExtensionMockGetIssuedAtTime) Set(f func() (r time.Time)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuedAtTimeFunc = f
	return m.mock
}

//GetIssuedAtTime implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetIssuedAtTime() (r time.Time) {
	counter := atomic.AddUint64(&m.GetIssuedAtTimePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuedAtTimeCounter, 1)

	if len(m.GetIssuedAtTimeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuedAtTimeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuedAtTime.")
			return
		}

		result := m.GetIssuedAtTimeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuedAtTime")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuedAtTimeMock.mainExpectation != nil {

		result := m.GetIssuedAtTimeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuedAtTime")
		}

		r = result.r

		return
	}

	if m.GetIssuedAtTimeFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuedAtTime.")
		return
	}

	return m.GetIssuedAtTimeFunc()
}

//GetIssuedAtTimeMinimockCounter returns a count of StaticProfileExtensionMock.GetIssuedAtTimeFunc invocations
func (m *StaticProfileExtensionMock) GetIssuedAtTimeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtTimeCounter)
}

//GetIssuedAtTimeMinimockPreCounter returns the value of StaticProfileExtensionMock.GetIssuedAtTime invocations
func (m *StaticProfileExtensionMock) GetIssuedAtTimeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuedAtTimePreCounter)
}

//GetIssuedAtTimeFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetIssuedAtTimeFinished() bool {
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

type mStaticProfileExtensionMockGetIssuerID struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetIssuerIDExpectation
	expectationSeries []*StaticProfileExtensionMockGetIssuerIDExpectation
}

type StaticProfileExtensionMockGetIssuerIDExpectation struct {
	result *StaticProfileExtensionMockGetIssuerIDResult
}

type StaticProfileExtensionMockGetIssuerIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of StaticProfileExtension.GetIssuerID is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetIssuerID) Expect() *mStaticProfileExtensionMockGetIssuerID {
	m.mock.GetIssuerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetIssuerID
func (m *mStaticProfileExtensionMockGetIssuerID) Return(r insolar.ShortNodeID) *StaticProfileExtensionMock {
	m.mock.GetIssuerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuerIDExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetIssuerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetIssuerID is expected once
func (m *mStaticProfileExtensionMockGetIssuerID) ExpectOnce() *StaticProfileExtensionMockGetIssuerIDExpectation {
	m.mock.GetIssuerIDFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetIssuerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetIssuerIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &StaticProfileExtensionMockGetIssuerIDResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetIssuerID method
func (m *mStaticProfileExtensionMockGetIssuerID) Set(f func() (r insolar.ShortNodeID)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerIDFunc = f
	return m.mock
}

//GetIssuerID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetIssuerID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetIssuerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuerIDCounter, 1)

	if len(m.GetIssuerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuerID.")
			return
		}

		result := m.GetIssuerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuerID")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuerIDMock.mainExpectation != nil {

		result := m.GetIssuerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuerID")
		}

		r = result.r

		return
	}

	if m.GetIssuerIDFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuerID.")
		return
	}

	return m.GetIssuerIDFunc()
}

//GetIssuerIDMinimockCounter returns a count of StaticProfileExtensionMock.GetIssuerIDFunc invocations
func (m *StaticProfileExtensionMock) GetIssuerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerIDCounter)
}

//GetIssuerIDMinimockPreCounter returns the value of StaticProfileExtensionMock.GetIssuerID invocations
func (m *StaticProfileExtensionMock) GetIssuerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerIDPreCounter)
}

//GetIssuerIDFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetIssuerIDFinished() bool {
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

type mStaticProfileExtensionMockGetIssuerSignature struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetIssuerSignatureExpectation
	expectationSeries []*StaticProfileExtensionMockGetIssuerSignatureExpectation
}

type StaticProfileExtensionMockGetIssuerSignatureExpectation struct {
	result *StaticProfileExtensionMockGetIssuerSignatureResult
}

type StaticProfileExtensionMockGetIssuerSignatureResult struct {
	r cryptkit.SignatureHolder
}

//Expect specifies that invocation of StaticProfileExtension.GetIssuerSignature is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetIssuerSignature) Expect() *mStaticProfileExtensionMockGetIssuerSignature {
	m.mock.GetIssuerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuerSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetIssuerSignature
func (m *mStaticProfileExtensionMockGetIssuerSignature) Return(r cryptkit.SignatureHolder) *StaticProfileExtensionMock {
	m.mock.GetIssuerSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetIssuerSignatureExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetIssuerSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetIssuerSignature is expected once
func (m *mStaticProfileExtensionMockGetIssuerSignature) ExpectOnce() *StaticProfileExtensionMockGetIssuerSignatureExpectation {
	m.mock.GetIssuerSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetIssuerSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetIssuerSignatureExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &StaticProfileExtensionMockGetIssuerSignatureResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetIssuerSignature method
func (m *mStaticProfileExtensionMockGetIssuerSignature) Set(f func() (r cryptkit.SignatureHolder)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIssuerSignatureFunc = f
	return m.mock
}

//GetIssuerSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetIssuerSignature() (r cryptkit.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetIssuerSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetIssuerSignatureCounter, 1)

	if len(m.GetIssuerSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIssuerSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuerSignature.")
			return
		}

		result := m.GetIssuerSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuerSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetIssuerSignatureMock.mainExpectation != nil {

		result := m.GetIssuerSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetIssuerSignature")
		}

		r = result.r

		return
	}

	if m.GetIssuerSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetIssuerSignature.")
		return
	}

	return m.GetIssuerSignatureFunc()
}

//GetIssuerSignatureMinimockCounter returns a count of StaticProfileExtensionMock.GetIssuerSignatureFunc invocations
func (m *StaticProfileExtensionMock) GetIssuerSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerSignatureCounter)
}

//GetIssuerSignatureMinimockPreCounter returns the value of StaticProfileExtensionMock.GetIssuerSignature invocations
func (m *StaticProfileExtensionMock) GetIssuerSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIssuerSignaturePreCounter)
}

//GetIssuerSignatureFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetIssuerSignatureFinished() bool {
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

type mStaticProfileExtensionMockGetPowerLevels struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetPowerLevelsExpectation
	expectationSeries []*StaticProfileExtensionMockGetPowerLevelsExpectation
}

type StaticProfileExtensionMockGetPowerLevelsExpectation struct {
	result *StaticProfileExtensionMockGetPowerLevelsResult
}

type StaticProfileExtensionMockGetPowerLevelsResult struct {
	r member.PowerSet
}

//Expect specifies that invocation of StaticProfileExtension.GetPowerLevels is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetPowerLevels) Expect() *mStaticProfileExtensionMockGetPowerLevels {
	m.mock.GetPowerLevelsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetPowerLevelsExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetPowerLevels
func (m *mStaticProfileExtensionMockGetPowerLevels) Return(r member.PowerSet) *StaticProfileExtensionMock {
	m.mock.GetPowerLevelsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetPowerLevelsExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetPowerLevelsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetPowerLevels is expected once
func (m *mStaticProfileExtensionMockGetPowerLevels) ExpectOnce() *StaticProfileExtensionMockGetPowerLevelsExpectation {
	m.mock.GetPowerLevelsFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetPowerLevelsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetPowerLevelsExpectation) Return(r member.PowerSet) {
	e.result = &StaticProfileExtensionMockGetPowerLevelsResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetPowerLevels method
func (m *mStaticProfileExtensionMockGetPowerLevels) Set(f func() (r member.PowerSet)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPowerLevelsFunc = f
	return m.mock
}

//GetPowerLevels implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetPowerLevels() (r member.PowerSet) {
	counter := atomic.AddUint64(&m.GetPowerLevelsPreCounter, 1)
	defer atomic.AddUint64(&m.GetPowerLevelsCounter, 1)

	if len(m.GetPowerLevelsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPowerLevelsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetPowerLevels.")
			return
		}

		result := m.GetPowerLevelsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetPowerLevels")
			return
		}

		r = result.r

		return
	}

	if m.GetPowerLevelsMock.mainExpectation != nil {

		result := m.GetPowerLevelsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetPowerLevels")
		}

		r = result.r

		return
	}

	if m.GetPowerLevelsFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetPowerLevels.")
		return
	}

	return m.GetPowerLevelsFunc()
}

//GetPowerLevelsMinimockCounter returns a count of StaticProfileExtensionMock.GetPowerLevelsFunc invocations
func (m *StaticProfileExtensionMock) GetPowerLevelsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPowerLevelsCounter)
}

//GetPowerLevelsMinimockPreCounter returns the value of StaticProfileExtensionMock.GetPowerLevels invocations
func (m *StaticProfileExtensionMock) GetPowerLevelsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPowerLevelsPreCounter)
}

//GetPowerLevelsFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetPowerLevelsFinished() bool {
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

type mStaticProfileExtensionMockGetReference struct {
	mock              *StaticProfileExtensionMock
	mainExpectation   *StaticProfileExtensionMockGetReferenceExpectation
	expectationSeries []*StaticProfileExtensionMockGetReferenceExpectation
}

type StaticProfileExtensionMockGetReferenceExpectation struct {
	result *StaticProfileExtensionMockGetReferenceResult
}

type StaticProfileExtensionMockGetReferenceResult struct {
	r insolar.Reference
}

//Expect specifies that invocation of StaticProfileExtension.GetReference is expected from 1 to Infinity times
func (m *mStaticProfileExtensionMockGetReference) Expect() *mStaticProfileExtensionMockGetReference {
	m.mock.GetReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetReferenceExpectation{}
	}

	return m
}

//Return specifies results of invocation of StaticProfileExtension.GetReference
func (m *mStaticProfileExtensionMockGetReference) Return(r insolar.Reference) *StaticProfileExtensionMock {
	m.mock.GetReferenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StaticProfileExtensionMockGetReferenceExpectation{}
	}
	m.mainExpectation.result = &StaticProfileExtensionMockGetReferenceResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StaticProfileExtension.GetReference is expected once
func (m *mStaticProfileExtensionMockGetReference) ExpectOnce() *StaticProfileExtensionMockGetReferenceExpectation {
	m.mock.GetReferenceFunc = nil
	m.mainExpectation = nil

	expectation := &StaticProfileExtensionMockGetReferenceExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StaticProfileExtensionMockGetReferenceExpectation) Return(r insolar.Reference) {
	e.result = &StaticProfileExtensionMockGetReferenceResult{r}
}

//Set uses given function f as a mock of StaticProfileExtension.GetReference method
func (m *mStaticProfileExtensionMockGetReference) Set(f func() (r insolar.Reference)) *StaticProfileExtensionMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetReferenceFunc = f
	return m.mock
}

//GetReference implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension interface
func (m *StaticProfileExtensionMock) GetReference() (r insolar.Reference) {
	counter := atomic.AddUint64(&m.GetReferencePreCounter, 1)
	defer atomic.AddUint64(&m.GetReferenceCounter, 1)

	if len(m.GetReferenceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetReferenceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetReference.")
			return
		}

		result := m.GetReferenceMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetReference")
			return
		}

		r = result.r

		return
	}

	if m.GetReferenceMock.mainExpectation != nil {

		result := m.GetReferenceMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StaticProfileExtensionMock.GetReference")
		}

		r = result.r

		return
	}

	if m.GetReferenceFunc == nil {
		m.t.Fatalf("Unexpected call to StaticProfileExtensionMock.GetReference.")
		return
	}

	return m.GetReferenceFunc()
}

//GetReferenceMinimockCounter returns a count of StaticProfileExtensionMock.GetReferenceFunc invocations
func (m *StaticProfileExtensionMock) GetReferenceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetReferenceCounter)
}

//GetReferenceMinimockPreCounter returns the value of StaticProfileExtensionMock.GetReference invocations
func (m *StaticProfileExtensionMock) GetReferenceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetReferencePreCounter)
}

//GetReferenceFinished returns true if mock invocations count is ok
func (m *StaticProfileExtensionMock) GetReferenceFinished() bool {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StaticProfileExtensionMock) ValidateCallCounters() {

	if !m.GetExtraEndpointsFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetExtraEndpoints")
	}

	if !m.GetIntroducedNodeIDFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIntroducedNodeID")
	}

	if !m.GetIssuedAtPulseFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuedAtPulse")
	}

	if !m.GetIssuedAtTimeFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuedAtTime")
	}

	if !m.GetIssuerIDFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuerID")
	}

	if !m.GetIssuerSignatureFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuerSignature")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetPowerLevels")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetReference")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StaticProfileExtensionMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StaticProfileExtensionMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StaticProfileExtensionMock) MinimockFinish() {

	if !m.GetExtraEndpointsFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetExtraEndpoints")
	}

	if !m.GetIntroducedNodeIDFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIntroducedNodeID")
	}

	if !m.GetIssuedAtPulseFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuedAtPulse")
	}

	if !m.GetIssuedAtTimeFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuedAtTime")
	}

	if !m.GetIssuerIDFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuerID")
	}

	if !m.GetIssuerSignatureFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetIssuerSignature")
	}

	if !m.GetPowerLevelsFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetPowerLevels")
	}

	if !m.GetReferenceFinished() {
		m.t.Fatal("Expected call to StaticProfileExtensionMock.GetReference")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StaticProfileExtensionMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StaticProfileExtensionMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetExtraEndpointsFinished()
		ok = ok && m.GetIntroducedNodeIDFinished()
		ok = ok && m.GetIssuedAtPulseFinished()
		ok = ok && m.GetIssuedAtTimeFinished()
		ok = ok && m.GetIssuerIDFinished()
		ok = ok && m.GetIssuerSignatureFinished()
		ok = ok && m.GetPowerLevelsFinished()
		ok = ok && m.GetReferenceFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetExtraEndpointsFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetExtraEndpoints")
			}

			if !m.GetIntroducedNodeIDFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetIntroducedNodeID")
			}

			if !m.GetIssuedAtPulseFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetIssuedAtPulse")
			}

			if !m.GetIssuedAtTimeFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetIssuedAtTime")
			}

			if !m.GetIssuerIDFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetIssuerID")
			}

			if !m.GetIssuerSignatureFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetIssuerSignature")
			}

			if !m.GetPowerLevelsFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetPowerLevels")
			}

			if !m.GetReferenceFinished() {
				m.t.Error("Expected call to StaticProfileExtensionMock.GetReference")
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
func (m *StaticProfileExtensionMock) AllMocksCalled() bool {

	if !m.GetExtraEndpointsFinished() {
		return false
	}

	if !m.GetIntroducedNodeIDFinished() {
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

	if !m.GetPowerLevelsFinished() {
		return false
	}

	if !m.GetReferenceFinished() {
		return false
	}

	return true
}

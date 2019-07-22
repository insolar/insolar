package census

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "VersionedRegistries" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/census
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	pulse "github.com/insolar/insolar/network/consensus/common/pulse"

	testify_assert "github.com/stretchr/testify/assert"
)

//VersionedRegistriesMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.VersionedRegistries
type VersionedRegistriesMock struct {
	t minimock.Tester

	CommitNextPulseFunc       func(p pulse.Data, p1 OnlinePopulation) (r VersionedRegistries)
	CommitNextPulseCounter    uint64
	CommitNextPulsePreCounter uint64
	CommitNextPulseMock       mVersionedRegistriesMockCommitNextPulse

	GetMandateRegistryFunc       func() (r MandateRegistry)
	GetMandateRegistryCounter    uint64
	GetMandateRegistryPreCounter uint64
	GetMandateRegistryMock       mVersionedRegistriesMockGetMandateRegistry

	GetMisbehaviorRegistryFunc       func() (r MisbehaviorRegistry)
	GetMisbehaviorRegistryCounter    uint64
	GetMisbehaviorRegistryPreCounter uint64
	GetMisbehaviorRegistryMock       mVersionedRegistriesMockGetMisbehaviorRegistry

	GetOfflinePopulationFunc       func() (r OfflinePopulation)
	GetOfflinePopulationCounter    uint64
	GetOfflinePopulationPreCounter uint64
	GetOfflinePopulationMock       mVersionedRegistriesMockGetOfflinePopulation

	GetVersionPulseDataFunc       func() (r pulse.Data)
	GetVersionPulseDataCounter    uint64
	GetVersionPulseDataPreCounter uint64
	GetVersionPulseDataMock       mVersionedRegistriesMockGetVersionPulseData
}

//NewVersionedRegistriesMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/census.VersionedRegistries
func NewVersionedRegistriesMock(t minimock.Tester) *VersionedRegistriesMock {
	m := &VersionedRegistriesMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CommitNextPulseMock = mVersionedRegistriesMockCommitNextPulse{mock: m}
	m.GetMandateRegistryMock = mVersionedRegistriesMockGetMandateRegistry{mock: m}
	m.GetMisbehaviorRegistryMock = mVersionedRegistriesMockGetMisbehaviorRegistry{mock: m}
	m.GetOfflinePopulationMock = mVersionedRegistriesMockGetOfflinePopulation{mock: m}
	m.GetVersionPulseDataMock = mVersionedRegistriesMockGetVersionPulseData{mock: m}

	return m
}

type mVersionedRegistriesMockCommitNextPulse struct {
	mock              *VersionedRegistriesMock
	mainExpectation   *VersionedRegistriesMockCommitNextPulseExpectation
	expectationSeries []*VersionedRegistriesMockCommitNextPulseExpectation
}

type VersionedRegistriesMockCommitNextPulseExpectation struct {
	input  *VersionedRegistriesMockCommitNextPulseInput
	result *VersionedRegistriesMockCommitNextPulseResult
}

type VersionedRegistriesMockCommitNextPulseInput struct {
	p  pulse.Data
	p1 OnlinePopulation
}

type VersionedRegistriesMockCommitNextPulseResult struct {
	r VersionedRegistries
}

//Expect specifies that invocation of VersionedRegistries.CommitNextPulse is expected from 1 to Infinity times
func (m *mVersionedRegistriesMockCommitNextPulse) Expect(p pulse.Data, p1 OnlinePopulation) *mVersionedRegistriesMockCommitNextPulse {
	m.mock.CommitNextPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockCommitNextPulseExpectation{}
	}
	m.mainExpectation.input = &VersionedRegistriesMockCommitNextPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of VersionedRegistries.CommitNextPulse
func (m *mVersionedRegistriesMockCommitNextPulse) Return(r VersionedRegistries) *VersionedRegistriesMock {
	m.mock.CommitNextPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockCommitNextPulseExpectation{}
	}
	m.mainExpectation.result = &VersionedRegistriesMockCommitNextPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of VersionedRegistries.CommitNextPulse is expected once
func (m *mVersionedRegistriesMockCommitNextPulse) ExpectOnce(p pulse.Data, p1 OnlinePopulation) *VersionedRegistriesMockCommitNextPulseExpectation {
	m.mock.CommitNextPulseFunc = nil
	m.mainExpectation = nil

	expectation := &VersionedRegistriesMockCommitNextPulseExpectation{}
	expectation.input = &VersionedRegistriesMockCommitNextPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *VersionedRegistriesMockCommitNextPulseExpectation) Return(r VersionedRegistries) {
	e.result = &VersionedRegistriesMockCommitNextPulseResult{r}
}

//Set uses given function f as a mock of VersionedRegistries.CommitNextPulse method
func (m *mVersionedRegistriesMockCommitNextPulse) Set(f func(p pulse.Data, p1 OnlinePopulation) (r VersionedRegistries)) *VersionedRegistriesMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CommitNextPulseFunc = f
	return m.mock
}

//CommitNextPulse implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.VersionedRegistries interface
func (m *VersionedRegistriesMock) CommitNextPulse(p pulse.Data, p1 OnlinePopulation) (r VersionedRegistries) {
	counter := atomic.AddUint64(&m.CommitNextPulsePreCounter, 1)
	defer atomic.AddUint64(&m.CommitNextPulseCounter, 1)

	if len(m.CommitNextPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CommitNextPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to VersionedRegistriesMock.CommitNextPulse. %v %v", p, p1)
			return
		}

		input := m.CommitNextPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, VersionedRegistriesMockCommitNextPulseInput{p, p1}, "VersionedRegistries.CommitNextPulse got unexpected parameters")

		result := m.CommitNextPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.CommitNextPulse")
			return
		}

		r = result.r

		return
	}

	if m.CommitNextPulseMock.mainExpectation != nil {

		input := m.CommitNextPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, VersionedRegistriesMockCommitNextPulseInput{p, p1}, "VersionedRegistries.CommitNextPulse got unexpected parameters")
		}

		result := m.CommitNextPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.CommitNextPulse")
		}

		r = result.r

		return
	}

	if m.CommitNextPulseFunc == nil {
		m.t.Fatalf("Unexpected call to VersionedRegistriesMock.CommitNextPulse. %v %v", p, p1)
		return
	}

	return m.CommitNextPulseFunc(p, p1)
}

//CommitNextPulseMinimockCounter returns a count of VersionedRegistriesMock.CommitNextPulseFunc invocations
func (m *VersionedRegistriesMock) CommitNextPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CommitNextPulseCounter)
}

//CommitNextPulseMinimockPreCounter returns the value of VersionedRegistriesMock.CommitNextPulse invocations
func (m *VersionedRegistriesMock) CommitNextPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CommitNextPulsePreCounter)
}

//CommitNextPulseFinished returns true if mock invocations count is ok
func (m *VersionedRegistriesMock) CommitNextPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CommitNextPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CommitNextPulseCounter) == uint64(len(m.CommitNextPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CommitNextPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CommitNextPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CommitNextPulseFunc != nil {
		return atomic.LoadUint64(&m.CommitNextPulseCounter) > 0
	}

	return true
}

type mVersionedRegistriesMockGetMandateRegistry struct {
	mock              *VersionedRegistriesMock
	mainExpectation   *VersionedRegistriesMockGetMandateRegistryExpectation
	expectationSeries []*VersionedRegistriesMockGetMandateRegistryExpectation
}

type VersionedRegistriesMockGetMandateRegistryExpectation struct {
	result *VersionedRegistriesMockGetMandateRegistryResult
}

type VersionedRegistriesMockGetMandateRegistryResult struct {
	r MandateRegistry
}

//Expect specifies that invocation of VersionedRegistries.GetMandateRegistry is expected from 1 to Infinity times
func (m *mVersionedRegistriesMockGetMandateRegistry) Expect() *mVersionedRegistriesMockGetMandateRegistry {
	m.mock.GetMandateRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetMandateRegistryExpectation{}
	}

	return m
}

//Return specifies results of invocation of VersionedRegistries.GetMandateRegistry
func (m *mVersionedRegistriesMockGetMandateRegistry) Return(r MandateRegistry) *VersionedRegistriesMock {
	m.mock.GetMandateRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetMandateRegistryExpectation{}
	}
	m.mainExpectation.result = &VersionedRegistriesMockGetMandateRegistryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of VersionedRegistries.GetMandateRegistry is expected once
func (m *mVersionedRegistriesMockGetMandateRegistry) ExpectOnce() *VersionedRegistriesMockGetMandateRegistryExpectation {
	m.mock.GetMandateRegistryFunc = nil
	m.mainExpectation = nil

	expectation := &VersionedRegistriesMockGetMandateRegistryExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *VersionedRegistriesMockGetMandateRegistryExpectation) Return(r MandateRegistry) {
	e.result = &VersionedRegistriesMockGetMandateRegistryResult{r}
}

//Set uses given function f as a mock of VersionedRegistries.GetMandateRegistry method
func (m *mVersionedRegistriesMockGetMandateRegistry) Set(f func() (r MandateRegistry)) *VersionedRegistriesMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMandateRegistryFunc = f
	return m.mock
}

//GetMandateRegistry implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.VersionedRegistries interface
func (m *VersionedRegistriesMock) GetMandateRegistry() (r MandateRegistry) {
	counter := atomic.AddUint64(&m.GetMandateRegistryPreCounter, 1)
	defer atomic.AddUint64(&m.GetMandateRegistryCounter, 1)

	if len(m.GetMandateRegistryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMandateRegistryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetMandateRegistry.")
			return
		}

		result := m.GetMandateRegistryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetMandateRegistry")
			return
		}

		r = result.r

		return
	}

	if m.GetMandateRegistryMock.mainExpectation != nil {

		result := m.GetMandateRegistryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetMandateRegistry")
		}

		r = result.r

		return
	}

	if m.GetMandateRegistryFunc == nil {
		m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetMandateRegistry.")
		return
	}

	return m.GetMandateRegistryFunc()
}

//GetMandateRegistryMinimockCounter returns a count of VersionedRegistriesMock.GetMandateRegistryFunc invocations
func (m *VersionedRegistriesMock) GetMandateRegistryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMandateRegistryCounter)
}

//GetMandateRegistryMinimockPreCounter returns the value of VersionedRegistriesMock.GetMandateRegistry invocations
func (m *VersionedRegistriesMock) GetMandateRegistryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMandateRegistryPreCounter)
}

//GetMandateRegistryFinished returns true if mock invocations count is ok
func (m *VersionedRegistriesMock) GetMandateRegistryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMandateRegistryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetMandateRegistryCounter) == uint64(len(m.GetMandateRegistryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMandateRegistryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetMandateRegistryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetMandateRegistryFunc != nil {
		return atomic.LoadUint64(&m.GetMandateRegistryCounter) > 0
	}

	return true
}

type mVersionedRegistriesMockGetMisbehaviorRegistry struct {
	mock              *VersionedRegistriesMock
	mainExpectation   *VersionedRegistriesMockGetMisbehaviorRegistryExpectation
	expectationSeries []*VersionedRegistriesMockGetMisbehaviorRegistryExpectation
}

type VersionedRegistriesMockGetMisbehaviorRegistryExpectation struct {
	result *VersionedRegistriesMockGetMisbehaviorRegistryResult
}

type VersionedRegistriesMockGetMisbehaviorRegistryResult struct {
	r MisbehaviorRegistry
}

//Expect specifies that invocation of VersionedRegistries.GetMisbehaviorRegistry is expected from 1 to Infinity times
func (m *mVersionedRegistriesMockGetMisbehaviorRegistry) Expect() *mVersionedRegistriesMockGetMisbehaviorRegistry {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetMisbehaviorRegistryExpectation{}
	}

	return m
}

//Return specifies results of invocation of VersionedRegistries.GetMisbehaviorRegistry
func (m *mVersionedRegistriesMockGetMisbehaviorRegistry) Return(r MisbehaviorRegistry) *VersionedRegistriesMock {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetMisbehaviorRegistryExpectation{}
	}
	m.mainExpectation.result = &VersionedRegistriesMockGetMisbehaviorRegistryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of VersionedRegistries.GetMisbehaviorRegistry is expected once
func (m *mVersionedRegistriesMockGetMisbehaviorRegistry) ExpectOnce() *VersionedRegistriesMockGetMisbehaviorRegistryExpectation {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.mainExpectation = nil

	expectation := &VersionedRegistriesMockGetMisbehaviorRegistryExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *VersionedRegistriesMockGetMisbehaviorRegistryExpectation) Return(r MisbehaviorRegistry) {
	e.result = &VersionedRegistriesMockGetMisbehaviorRegistryResult{r}
}

//Set uses given function f as a mock of VersionedRegistries.GetMisbehaviorRegistry method
func (m *mVersionedRegistriesMockGetMisbehaviorRegistry) Set(f func() (r MisbehaviorRegistry)) *VersionedRegistriesMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMisbehaviorRegistryFunc = f
	return m.mock
}

//GetMisbehaviorRegistry implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.VersionedRegistries interface
func (m *VersionedRegistriesMock) GetMisbehaviorRegistry() (r MisbehaviorRegistry) {
	counter := atomic.AddUint64(&m.GetMisbehaviorRegistryPreCounter, 1)
	defer atomic.AddUint64(&m.GetMisbehaviorRegistryCounter, 1)

	if len(m.GetMisbehaviorRegistryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMisbehaviorRegistryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetMisbehaviorRegistry.")
			return
		}

		result := m.GetMisbehaviorRegistryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetMisbehaviorRegistry")
			return
		}

		r = result.r

		return
	}

	if m.GetMisbehaviorRegistryMock.mainExpectation != nil {

		result := m.GetMisbehaviorRegistryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetMisbehaviorRegistry")
		}

		r = result.r

		return
	}

	if m.GetMisbehaviorRegistryFunc == nil {
		m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetMisbehaviorRegistry.")
		return
	}

	return m.GetMisbehaviorRegistryFunc()
}

//GetMisbehaviorRegistryMinimockCounter returns a count of VersionedRegistriesMock.GetMisbehaviorRegistryFunc invocations
func (m *VersionedRegistriesMock) GetMisbehaviorRegistryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMisbehaviorRegistryCounter)
}

//GetMisbehaviorRegistryMinimockPreCounter returns the value of VersionedRegistriesMock.GetMisbehaviorRegistry invocations
func (m *VersionedRegistriesMock) GetMisbehaviorRegistryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMisbehaviorRegistryPreCounter)
}

//GetMisbehaviorRegistryFinished returns true if mock invocations count is ok
func (m *VersionedRegistriesMock) GetMisbehaviorRegistryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMisbehaviorRegistryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetMisbehaviorRegistryCounter) == uint64(len(m.GetMisbehaviorRegistryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMisbehaviorRegistryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetMisbehaviorRegistryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetMisbehaviorRegistryFunc != nil {
		return atomic.LoadUint64(&m.GetMisbehaviorRegistryCounter) > 0
	}

	return true
}

type mVersionedRegistriesMockGetOfflinePopulation struct {
	mock              *VersionedRegistriesMock
	mainExpectation   *VersionedRegistriesMockGetOfflinePopulationExpectation
	expectationSeries []*VersionedRegistriesMockGetOfflinePopulationExpectation
}

type VersionedRegistriesMockGetOfflinePopulationExpectation struct {
	result *VersionedRegistriesMockGetOfflinePopulationResult
}

type VersionedRegistriesMockGetOfflinePopulationResult struct {
	r OfflinePopulation
}

//Expect specifies that invocation of VersionedRegistries.GetOfflinePopulation is expected from 1 to Infinity times
func (m *mVersionedRegistriesMockGetOfflinePopulation) Expect() *mVersionedRegistriesMockGetOfflinePopulation {
	m.mock.GetOfflinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetOfflinePopulationExpectation{}
	}

	return m
}

//Return specifies results of invocation of VersionedRegistries.GetOfflinePopulation
func (m *mVersionedRegistriesMockGetOfflinePopulation) Return(r OfflinePopulation) *VersionedRegistriesMock {
	m.mock.GetOfflinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetOfflinePopulationExpectation{}
	}
	m.mainExpectation.result = &VersionedRegistriesMockGetOfflinePopulationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of VersionedRegistries.GetOfflinePopulation is expected once
func (m *mVersionedRegistriesMockGetOfflinePopulation) ExpectOnce() *VersionedRegistriesMockGetOfflinePopulationExpectation {
	m.mock.GetOfflinePopulationFunc = nil
	m.mainExpectation = nil

	expectation := &VersionedRegistriesMockGetOfflinePopulationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *VersionedRegistriesMockGetOfflinePopulationExpectation) Return(r OfflinePopulation) {
	e.result = &VersionedRegistriesMockGetOfflinePopulationResult{r}
}

//Set uses given function f as a mock of VersionedRegistries.GetOfflinePopulation method
func (m *mVersionedRegistriesMockGetOfflinePopulation) Set(f func() (r OfflinePopulation)) *VersionedRegistriesMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOfflinePopulationFunc = f
	return m.mock
}

//GetOfflinePopulation implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.VersionedRegistries interface
func (m *VersionedRegistriesMock) GetOfflinePopulation() (r OfflinePopulation) {
	counter := atomic.AddUint64(&m.GetOfflinePopulationPreCounter, 1)
	defer atomic.AddUint64(&m.GetOfflinePopulationCounter, 1)

	if len(m.GetOfflinePopulationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOfflinePopulationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetOfflinePopulation.")
			return
		}

		result := m.GetOfflinePopulationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetOfflinePopulation")
			return
		}

		r = result.r

		return
	}

	if m.GetOfflinePopulationMock.mainExpectation != nil {

		result := m.GetOfflinePopulationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetOfflinePopulation")
		}

		r = result.r

		return
	}

	if m.GetOfflinePopulationFunc == nil {
		m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetOfflinePopulation.")
		return
	}

	return m.GetOfflinePopulationFunc()
}

//GetOfflinePopulationMinimockCounter returns a count of VersionedRegistriesMock.GetOfflinePopulationFunc invocations
func (m *VersionedRegistriesMock) GetOfflinePopulationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOfflinePopulationCounter)
}

//GetOfflinePopulationMinimockPreCounter returns the value of VersionedRegistriesMock.GetOfflinePopulation invocations
func (m *VersionedRegistriesMock) GetOfflinePopulationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOfflinePopulationPreCounter)
}

//GetOfflinePopulationFinished returns true if mock invocations count is ok
func (m *VersionedRegistriesMock) GetOfflinePopulationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOfflinePopulationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOfflinePopulationCounter) == uint64(len(m.GetOfflinePopulationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOfflinePopulationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOfflinePopulationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOfflinePopulationFunc != nil {
		return atomic.LoadUint64(&m.GetOfflinePopulationCounter) > 0
	}

	return true
}

type mVersionedRegistriesMockGetVersionPulseData struct {
	mock              *VersionedRegistriesMock
	mainExpectation   *VersionedRegistriesMockGetVersionPulseDataExpectation
	expectationSeries []*VersionedRegistriesMockGetVersionPulseDataExpectation
}

type VersionedRegistriesMockGetVersionPulseDataExpectation struct {
	result *VersionedRegistriesMockGetVersionPulseDataResult
}

type VersionedRegistriesMockGetVersionPulseDataResult struct {
	r pulse.Data
}

//Expect specifies that invocation of VersionedRegistries.GetVersionPulseData is expected from 1 to Infinity times
func (m *mVersionedRegistriesMockGetVersionPulseData) Expect() *mVersionedRegistriesMockGetVersionPulseData {
	m.mock.GetVersionPulseDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetVersionPulseDataExpectation{}
	}

	return m
}

//Return specifies results of invocation of VersionedRegistries.GetVersionPulseData
func (m *mVersionedRegistriesMockGetVersionPulseData) Return(r pulse.Data) *VersionedRegistriesMock {
	m.mock.GetVersionPulseDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &VersionedRegistriesMockGetVersionPulseDataExpectation{}
	}
	m.mainExpectation.result = &VersionedRegistriesMockGetVersionPulseDataResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of VersionedRegistries.GetVersionPulseData is expected once
func (m *mVersionedRegistriesMockGetVersionPulseData) ExpectOnce() *VersionedRegistriesMockGetVersionPulseDataExpectation {
	m.mock.GetVersionPulseDataFunc = nil
	m.mainExpectation = nil

	expectation := &VersionedRegistriesMockGetVersionPulseDataExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *VersionedRegistriesMockGetVersionPulseDataExpectation) Return(r pulse.Data) {
	e.result = &VersionedRegistriesMockGetVersionPulseDataResult{r}
}

//Set uses given function f as a mock of VersionedRegistries.GetVersionPulseData method
func (m *mVersionedRegistriesMockGetVersionPulseData) Set(f func() (r pulse.Data)) *VersionedRegistriesMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetVersionPulseDataFunc = f
	return m.mock
}

//GetVersionPulseData implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.VersionedRegistries interface
func (m *VersionedRegistriesMock) GetVersionPulseData() (r pulse.Data) {
	counter := atomic.AddUint64(&m.GetVersionPulseDataPreCounter, 1)
	defer atomic.AddUint64(&m.GetVersionPulseDataCounter, 1)

	if len(m.GetVersionPulseDataMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetVersionPulseDataMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetVersionPulseData.")
			return
		}

		result := m.GetVersionPulseDataMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetVersionPulseData")
			return
		}

		r = result.r

		return
	}

	if m.GetVersionPulseDataMock.mainExpectation != nil {

		result := m.GetVersionPulseDataMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the VersionedRegistriesMock.GetVersionPulseData")
		}

		r = result.r

		return
	}

	if m.GetVersionPulseDataFunc == nil {
		m.t.Fatalf("Unexpected call to VersionedRegistriesMock.GetVersionPulseData.")
		return
	}

	return m.GetVersionPulseDataFunc()
}

//GetVersionPulseDataMinimockCounter returns a count of VersionedRegistriesMock.GetVersionPulseDataFunc invocations
func (m *VersionedRegistriesMock) GetVersionPulseDataMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetVersionPulseDataCounter)
}

//GetVersionPulseDataMinimockPreCounter returns the value of VersionedRegistriesMock.GetVersionPulseData invocations
func (m *VersionedRegistriesMock) GetVersionPulseDataMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetVersionPulseDataPreCounter)
}

//GetVersionPulseDataFinished returns true if mock invocations count is ok
func (m *VersionedRegistriesMock) GetVersionPulseDataFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetVersionPulseDataMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetVersionPulseDataCounter) == uint64(len(m.GetVersionPulseDataMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetVersionPulseDataMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetVersionPulseDataCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetVersionPulseDataFunc != nil {
		return atomic.LoadUint64(&m.GetVersionPulseDataCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *VersionedRegistriesMock) ValidateCallCounters() {

	if !m.CommitNextPulseFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.CommitNextPulse")
	}

	if !m.GetMandateRegistryFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetMandateRegistry")
	}

	if !m.GetMisbehaviorRegistryFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetMisbehaviorRegistry")
	}

	if !m.GetOfflinePopulationFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetOfflinePopulation")
	}

	if !m.GetVersionPulseDataFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetVersionPulseData")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *VersionedRegistriesMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *VersionedRegistriesMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *VersionedRegistriesMock) MinimockFinish() {

	if !m.CommitNextPulseFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.CommitNextPulse")
	}

	if !m.GetMandateRegistryFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetMandateRegistry")
	}

	if !m.GetMisbehaviorRegistryFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetMisbehaviorRegistry")
	}

	if !m.GetOfflinePopulationFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetOfflinePopulation")
	}

	if !m.GetVersionPulseDataFinished() {
		m.t.Fatal("Expected call to VersionedRegistriesMock.GetVersionPulseData")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *VersionedRegistriesMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *VersionedRegistriesMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CommitNextPulseFinished()
		ok = ok && m.GetMandateRegistryFinished()
		ok = ok && m.GetMisbehaviorRegistryFinished()
		ok = ok && m.GetOfflinePopulationFinished()
		ok = ok && m.GetVersionPulseDataFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CommitNextPulseFinished() {
				m.t.Error("Expected call to VersionedRegistriesMock.CommitNextPulse")
			}

			if !m.GetMandateRegistryFinished() {
				m.t.Error("Expected call to VersionedRegistriesMock.GetMandateRegistry")
			}

			if !m.GetMisbehaviorRegistryFinished() {
				m.t.Error("Expected call to VersionedRegistriesMock.GetMisbehaviorRegistry")
			}

			if !m.GetOfflinePopulationFinished() {
				m.t.Error("Expected call to VersionedRegistriesMock.GetOfflinePopulation")
			}

			if !m.GetVersionPulseDataFinished() {
				m.t.Error("Expected call to VersionedRegistriesMock.GetVersionPulseData")
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
func (m *VersionedRegistriesMock) AllMocksCalled() bool {

	if !m.CommitNextPulseFinished() {
		return false
	}

	if !m.GetMandateRegistryFinished() {
		return false
	}

	if !m.GetMisbehaviorRegistryFinished() {
		return false
	}

	if !m.GetOfflinePopulationFinished() {
		return false
	}

	if !m.GetVersionPulseDataFinished() {
		return false
	}

	return true
}

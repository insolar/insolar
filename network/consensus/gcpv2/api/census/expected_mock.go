package census

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Expected" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/census
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	pulse "github.com/insolar/insolar/network/consensus/common/pulse"
	profiles "github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	proofs "github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"

	testify_assert "github.com/stretchr/testify/assert"
)

//ExpectedMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected
type ExpectedMock struct {
	t minimock.Tester

	CreateBuilderFunc       func(p pulse.Number) (r Builder)
	CreateBuilderCounter    uint64
	CreateBuilderPreCounter uint64
	CreateBuilderMock       mExpectedMockCreateBuilder

	GetCensusStateFunc       func() (r State)
	GetCensusStateCounter    uint64
	GetCensusStatePreCounter uint64
	GetCensusStateMock       mExpectedMockGetCensusState

	GetCloudStateHashFunc       func() (r proofs.CloudStateHash)
	GetCloudStateHashCounter    uint64
	GetCloudStateHashPreCounter uint64
	GetCloudStateHashMock       mExpectedMockGetCloudStateHash

	GetEvictedPopulationFunc       func() (r EvictedPopulation)
	GetEvictedPopulationCounter    uint64
	GetEvictedPopulationPreCounter uint64
	GetEvictedPopulationMock       mExpectedMockGetEvictedPopulation

	GetExpectedPulseNumberFunc       func() (r pulse.Number)
	GetExpectedPulseNumberCounter    uint64
	GetExpectedPulseNumberPreCounter uint64
	GetExpectedPulseNumberMock       mExpectedMockGetExpectedPulseNumber

	GetGlobulaStateHashFunc       func() (r proofs.GlobulaStateHash)
	GetGlobulaStateHashCounter    uint64
	GetGlobulaStateHashPreCounter uint64
	GetGlobulaStateHashMock       mExpectedMockGetGlobulaStateHash

	GetMandateRegistryFunc       func() (r MandateRegistry)
	GetMandateRegistryCounter    uint64
	GetMandateRegistryPreCounter uint64
	GetMandateRegistryMock       mExpectedMockGetMandateRegistry

	GetMisbehaviorRegistryFunc       func() (r MisbehaviorRegistry)
	GetMisbehaviorRegistryCounter    uint64
	GetMisbehaviorRegistryPreCounter uint64
	GetMisbehaviorRegistryMock       mExpectedMockGetMisbehaviorRegistry

	GetOfflinePopulationFunc       func() (r OfflinePopulation)
	GetOfflinePopulationCounter    uint64
	GetOfflinePopulationPreCounter uint64
	GetOfflinePopulationMock       mExpectedMockGetOfflinePopulation

	GetOnlinePopulationFunc       func() (r OnlinePopulation)
	GetOnlinePopulationCounter    uint64
	GetOnlinePopulationPreCounter uint64
	GetOnlinePopulationMock       mExpectedMockGetOnlinePopulation

	GetPreviousFunc       func() (r Active)
	GetPreviousCounter    uint64
	GetPreviousPreCounter uint64
	GetPreviousMock       mExpectedMockGetPrevious

	GetProfileFactoryFunc       func(p cryptkit.KeyStoreFactory) (r profiles.Factory)
	GetProfileFactoryCounter    uint64
	GetProfileFactoryPreCounter uint64
	GetProfileFactoryMock       mExpectedMockGetProfileFactory

	GetPulseNumberFunc       func() (r pulse.Number)
	GetPulseNumberCounter    uint64
	GetPulseNumberPreCounter uint64
	GetPulseNumberMock       mExpectedMockGetPulseNumber

	IsActiveFunc       func() (r bool)
	IsActiveCounter    uint64
	IsActivePreCounter uint64
	IsActiveMock       mExpectedMockIsActive

	MakeActiveFunc       func(p pulse.Data) (r Active)
	MakeActiveCounter    uint64
	MakeActivePreCounter uint64
	MakeActiveMock       mExpectedMockMakeActive
}

//NewExpectedMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected
func NewExpectedMock(t minimock.Tester) *ExpectedMock {
	m := &ExpectedMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateBuilderMock = mExpectedMockCreateBuilder{mock: m}
	m.GetCensusStateMock = mExpectedMockGetCensusState{mock: m}
	m.GetCloudStateHashMock = mExpectedMockGetCloudStateHash{mock: m}
	m.GetEvictedPopulationMock = mExpectedMockGetEvictedPopulation{mock: m}
	m.GetExpectedPulseNumberMock = mExpectedMockGetExpectedPulseNumber{mock: m}
	m.GetGlobulaStateHashMock = mExpectedMockGetGlobulaStateHash{mock: m}
	m.GetMandateRegistryMock = mExpectedMockGetMandateRegistry{mock: m}
	m.GetMisbehaviorRegistryMock = mExpectedMockGetMisbehaviorRegistry{mock: m}
	m.GetOfflinePopulationMock = mExpectedMockGetOfflinePopulation{mock: m}
	m.GetOnlinePopulationMock = mExpectedMockGetOnlinePopulation{mock: m}
	m.GetPreviousMock = mExpectedMockGetPrevious{mock: m}
	m.GetProfileFactoryMock = mExpectedMockGetProfileFactory{mock: m}
	m.GetPulseNumberMock = mExpectedMockGetPulseNumber{mock: m}
	m.IsActiveMock = mExpectedMockIsActive{mock: m}
	m.MakeActiveMock = mExpectedMockMakeActive{mock: m}

	return m
}

type mExpectedMockCreateBuilder struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockCreateBuilderExpectation
	expectationSeries []*ExpectedMockCreateBuilderExpectation
}

type ExpectedMockCreateBuilderExpectation struct {
	input  *ExpectedMockCreateBuilderInput
	result *ExpectedMockCreateBuilderResult
}

type ExpectedMockCreateBuilderInput struct {
	p pulse.Number
}

type ExpectedMockCreateBuilderResult struct {
	r Builder
}

//Expect specifies that invocation of Expected.CreateBuilder is expected from 1 to Infinity times
func (m *mExpectedMockCreateBuilder) Expect(p pulse.Number) *mExpectedMockCreateBuilder {
	m.mock.CreateBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockCreateBuilderExpectation{}
	}
	m.mainExpectation.input = &ExpectedMockCreateBuilderInput{p}
	return m
}

//Return specifies results of invocation of Expected.CreateBuilder
func (m *mExpectedMockCreateBuilder) Return(r Builder) *ExpectedMock {
	m.mock.CreateBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockCreateBuilderExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockCreateBuilderResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.CreateBuilder is expected once
func (m *mExpectedMockCreateBuilder) ExpectOnce(p pulse.Number) *ExpectedMockCreateBuilderExpectation {
	m.mock.CreateBuilderFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockCreateBuilderExpectation{}
	expectation.input = &ExpectedMockCreateBuilderInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockCreateBuilderExpectation) Return(r Builder) {
	e.result = &ExpectedMockCreateBuilderResult{r}
}

//Set uses given function f as a mock of Expected.CreateBuilder method
func (m *mExpectedMockCreateBuilder) Set(f func(p pulse.Number) (r Builder)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateBuilderFunc = f
	return m.mock
}

//CreateBuilder implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) CreateBuilder(p pulse.Number) (r Builder) {
	counter := atomic.AddUint64(&m.CreateBuilderPreCounter, 1)
	defer atomic.AddUint64(&m.CreateBuilderCounter, 1)

	if len(m.CreateBuilderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateBuilderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.CreateBuilder. %v", p)
			return
		}

		input := m.CreateBuilderMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExpectedMockCreateBuilderInput{p}, "Expected.CreateBuilder got unexpected parameters")

		result := m.CreateBuilderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.CreateBuilder")
			return
		}

		r = result.r

		return
	}

	if m.CreateBuilderMock.mainExpectation != nil {

		input := m.CreateBuilderMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExpectedMockCreateBuilderInput{p}, "Expected.CreateBuilder got unexpected parameters")
		}

		result := m.CreateBuilderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.CreateBuilder")
		}

		r = result.r

		return
	}

	if m.CreateBuilderFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.CreateBuilder. %v", p)
		return
	}

	return m.CreateBuilderFunc(p)
}

//CreateBuilderMinimockCounter returns a count of ExpectedMock.CreateBuilderFunc invocations
func (m *ExpectedMock) CreateBuilderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateBuilderCounter)
}

//CreateBuilderMinimockPreCounter returns the value of ExpectedMock.CreateBuilder invocations
func (m *ExpectedMock) CreateBuilderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateBuilderPreCounter)
}

//CreateBuilderFinished returns true if mock invocations count is ok
func (m *ExpectedMock) CreateBuilderFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreateBuilderMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreateBuilderCounter) == uint64(len(m.CreateBuilderMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreateBuilderMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreateBuilderCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreateBuilderFunc != nil {
		return atomic.LoadUint64(&m.CreateBuilderCounter) > 0
	}

	return true
}

type mExpectedMockGetCensusState struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetCensusStateExpectation
	expectationSeries []*ExpectedMockGetCensusStateExpectation
}

type ExpectedMockGetCensusStateExpectation struct {
	result *ExpectedMockGetCensusStateResult
}

type ExpectedMockGetCensusStateResult struct {
	r State
}

//Expect specifies that invocation of Expected.GetCensusState is expected from 1 to Infinity times
func (m *mExpectedMockGetCensusState) Expect() *mExpectedMockGetCensusState {
	m.mock.GetCensusStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetCensusStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetCensusState
func (m *mExpectedMockGetCensusState) Return(r State) *ExpectedMock {
	m.mock.GetCensusStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetCensusStateExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetCensusStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetCensusState is expected once
func (m *mExpectedMockGetCensusState) ExpectOnce() *ExpectedMockGetCensusStateExpectation {
	m.mock.GetCensusStateFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetCensusStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetCensusStateExpectation) Return(r State) {
	e.result = &ExpectedMockGetCensusStateResult{r}
}

//Set uses given function f as a mock of Expected.GetCensusState method
func (m *mExpectedMockGetCensusState) Set(f func() (r State)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCensusStateFunc = f
	return m.mock
}

//GetCensusState implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetCensusState() (r State) {
	counter := atomic.AddUint64(&m.GetCensusStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetCensusStateCounter, 1)

	if len(m.GetCensusStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCensusStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetCensusState.")
			return
		}

		result := m.GetCensusStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetCensusState")
			return
		}

		r = result.r

		return
	}

	if m.GetCensusStateMock.mainExpectation != nil {

		result := m.GetCensusStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetCensusState")
		}

		r = result.r

		return
	}

	if m.GetCensusStateFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetCensusState.")
		return
	}

	return m.GetCensusStateFunc()
}

//GetCensusStateMinimockCounter returns a count of ExpectedMock.GetCensusStateFunc invocations
func (m *ExpectedMock) GetCensusStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCensusStateCounter)
}

//GetCensusStateMinimockPreCounter returns the value of ExpectedMock.GetCensusState invocations
func (m *ExpectedMock) GetCensusStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCensusStatePreCounter)
}

//GetCensusStateFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetCensusStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCensusStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCensusStateCounter) == uint64(len(m.GetCensusStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCensusStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCensusStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCensusStateFunc != nil {
		return atomic.LoadUint64(&m.GetCensusStateCounter) > 0
	}

	return true
}

type mExpectedMockGetCloudStateHash struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetCloudStateHashExpectation
	expectationSeries []*ExpectedMockGetCloudStateHashExpectation
}

type ExpectedMockGetCloudStateHashExpectation struct {
	result *ExpectedMockGetCloudStateHashResult
}

type ExpectedMockGetCloudStateHashResult struct {
	r proofs.CloudStateHash
}

//Expect specifies that invocation of Expected.GetCloudStateHash is expected from 1 to Infinity times
func (m *mExpectedMockGetCloudStateHash) Expect() *mExpectedMockGetCloudStateHash {
	m.mock.GetCloudStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetCloudStateHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetCloudStateHash
func (m *mExpectedMockGetCloudStateHash) Return(r proofs.CloudStateHash) *ExpectedMock {
	m.mock.GetCloudStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetCloudStateHashExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetCloudStateHashResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetCloudStateHash is expected once
func (m *mExpectedMockGetCloudStateHash) ExpectOnce() *ExpectedMockGetCloudStateHashExpectation {
	m.mock.GetCloudStateHashFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetCloudStateHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetCloudStateHashExpectation) Return(r proofs.CloudStateHash) {
	e.result = &ExpectedMockGetCloudStateHashResult{r}
}

//Set uses given function f as a mock of Expected.GetCloudStateHash method
func (m *mExpectedMockGetCloudStateHash) Set(f func() (r proofs.CloudStateHash)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCloudStateHashFunc = f
	return m.mock
}

//GetCloudStateHash implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetCloudStateHash() (r proofs.CloudStateHash) {
	counter := atomic.AddUint64(&m.GetCloudStateHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetCloudStateHashCounter, 1)

	if len(m.GetCloudStateHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCloudStateHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetCloudStateHash.")
			return
		}

		result := m.GetCloudStateHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetCloudStateHash")
			return
		}

		r = result.r

		return
	}

	if m.GetCloudStateHashMock.mainExpectation != nil {

		result := m.GetCloudStateHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetCloudStateHash")
		}

		r = result.r

		return
	}

	if m.GetCloudStateHashFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetCloudStateHash.")
		return
	}

	return m.GetCloudStateHashFunc()
}

//GetCloudStateHashMinimockCounter returns a count of ExpectedMock.GetCloudStateHashFunc invocations
func (m *ExpectedMock) GetCloudStateHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudStateHashCounter)
}

//GetCloudStateHashMinimockPreCounter returns the value of ExpectedMock.GetCloudStateHash invocations
func (m *ExpectedMock) GetCloudStateHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudStateHashPreCounter)
}

//GetCloudStateHashFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetCloudStateHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCloudStateHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCloudStateHashCounter) == uint64(len(m.GetCloudStateHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCloudStateHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCloudStateHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCloudStateHashFunc != nil {
		return atomic.LoadUint64(&m.GetCloudStateHashCounter) > 0
	}

	return true
}

type mExpectedMockGetEvictedPopulation struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetEvictedPopulationExpectation
	expectationSeries []*ExpectedMockGetEvictedPopulationExpectation
}

type ExpectedMockGetEvictedPopulationExpectation struct {
	result *ExpectedMockGetEvictedPopulationResult
}

type ExpectedMockGetEvictedPopulationResult struct {
	r EvictedPopulation
}

//Expect specifies that invocation of Expected.GetEvictedPopulation is expected from 1 to Infinity times
func (m *mExpectedMockGetEvictedPopulation) Expect() *mExpectedMockGetEvictedPopulation {
	m.mock.GetEvictedPopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetEvictedPopulationExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetEvictedPopulation
func (m *mExpectedMockGetEvictedPopulation) Return(r EvictedPopulation) *ExpectedMock {
	m.mock.GetEvictedPopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetEvictedPopulationExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetEvictedPopulationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetEvictedPopulation is expected once
func (m *mExpectedMockGetEvictedPopulation) ExpectOnce() *ExpectedMockGetEvictedPopulationExpectation {
	m.mock.GetEvictedPopulationFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetEvictedPopulationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetEvictedPopulationExpectation) Return(r EvictedPopulation) {
	e.result = &ExpectedMockGetEvictedPopulationResult{r}
}

//Set uses given function f as a mock of Expected.GetEvictedPopulation method
func (m *mExpectedMockGetEvictedPopulation) Set(f func() (r EvictedPopulation)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetEvictedPopulationFunc = f
	return m.mock
}

//GetEvictedPopulation implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetEvictedPopulation() (r EvictedPopulation) {
	counter := atomic.AddUint64(&m.GetEvictedPopulationPreCounter, 1)
	defer atomic.AddUint64(&m.GetEvictedPopulationCounter, 1)

	if len(m.GetEvictedPopulationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetEvictedPopulationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetEvictedPopulation.")
			return
		}

		result := m.GetEvictedPopulationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetEvictedPopulation")
			return
		}

		r = result.r

		return
	}

	if m.GetEvictedPopulationMock.mainExpectation != nil {

		result := m.GetEvictedPopulationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetEvictedPopulation")
		}

		r = result.r

		return
	}

	if m.GetEvictedPopulationFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetEvictedPopulation.")
		return
	}

	return m.GetEvictedPopulationFunc()
}

//GetEvictedPopulationMinimockCounter returns a count of ExpectedMock.GetEvictedPopulationFunc invocations
func (m *ExpectedMock) GetEvictedPopulationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetEvictedPopulationCounter)
}

//GetEvictedPopulationMinimockPreCounter returns the value of ExpectedMock.GetEvictedPopulation invocations
func (m *ExpectedMock) GetEvictedPopulationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetEvictedPopulationPreCounter)
}

//GetEvictedPopulationFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetEvictedPopulationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetEvictedPopulationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetEvictedPopulationCounter) == uint64(len(m.GetEvictedPopulationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetEvictedPopulationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetEvictedPopulationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetEvictedPopulationFunc != nil {
		return atomic.LoadUint64(&m.GetEvictedPopulationCounter) > 0
	}

	return true
}

type mExpectedMockGetExpectedPulseNumber struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetExpectedPulseNumberExpectation
	expectationSeries []*ExpectedMockGetExpectedPulseNumberExpectation
}

type ExpectedMockGetExpectedPulseNumberExpectation struct {
	result *ExpectedMockGetExpectedPulseNumberResult
}

type ExpectedMockGetExpectedPulseNumberResult struct {
	r pulse.Number
}

//Expect specifies that invocation of Expected.GetExpectedPulseNumber is expected from 1 to Infinity times
func (m *mExpectedMockGetExpectedPulseNumber) Expect() *mExpectedMockGetExpectedPulseNumber {
	m.mock.GetExpectedPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetExpectedPulseNumberExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetExpectedPulseNumber
func (m *mExpectedMockGetExpectedPulseNumber) Return(r pulse.Number) *ExpectedMock {
	m.mock.GetExpectedPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetExpectedPulseNumberExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetExpectedPulseNumberResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetExpectedPulseNumber is expected once
func (m *mExpectedMockGetExpectedPulseNumber) ExpectOnce() *ExpectedMockGetExpectedPulseNumberExpectation {
	m.mock.GetExpectedPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetExpectedPulseNumberExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetExpectedPulseNumberExpectation) Return(r pulse.Number) {
	e.result = &ExpectedMockGetExpectedPulseNumberResult{r}
}

//Set uses given function f as a mock of Expected.GetExpectedPulseNumber method
func (m *mExpectedMockGetExpectedPulseNumber) Set(f func() (r pulse.Number)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExpectedPulseNumberFunc = f
	return m.mock
}

//GetExpectedPulseNumber implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetExpectedPulseNumber() (r pulse.Number) {
	counter := atomic.AddUint64(&m.GetExpectedPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.GetExpectedPulseNumberCounter, 1)

	if len(m.GetExpectedPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExpectedPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetExpectedPulseNumber.")
			return
		}

		result := m.GetExpectedPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetExpectedPulseNumber")
			return
		}

		r = result.r

		return
	}

	if m.GetExpectedPulseNumberMock.mainExpectation != nil {

		result := m.GetExpectedPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetExpectedPulseNumber")
		}

		r = result.r

		return
	}

	if m.GetExpectedPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetExpectedPulseNumber.")
		return
	}

	return m.GetExpectedPulseNumberFunc()
}

//GetExpectedPulseNumberMinimockCounter returns a count of ExpectedMock.GetExpectedPulseNumberFunc invocations
func (m *ExpectedMock) GetExpectedPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExpectedPulseNumberCounter)
}

//GetExpectedPulseNumberMinimockPreCounter returns the value of ExpectedMock.GetExpectedPulseNumber invocations
func (m *ExpectedMock) GetExpectedPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExpectedPulseNumberPreCounter)
}

//GetExpectedPulseNumberFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetExpectedPulseNumberFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetExpectedPulseNumberMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetExpectedPulseNumberCounter) == uint64(len(m.GetExpectedPulseNumberMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetExpectedPulseNumberMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetExpectedPulseNumberCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetExpectedPulseNumberFunc != nil {
		return atomic.LoadUint64(&m.GetExpectedPulseNumberCounter) > 0
	}

	return true
}

type mExpectedMockGetGlobulaStateHash struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetGlobulaStateHashExpectation
	expectationSeries []*ExpectedMockGetGlobulaStateHashExpectation
}

type ExpectedMockGetGlobulaStateHashExpectation struct {
	result *ExpectedMockGetGlobulaStateHashResult
}

type ExpectedMockGetGlobulaStateHashResult struct {
	r proofs.GlobulaStateHash
}

//Expect specifies that invocation of Expected.GetGlobulaStateHash is expected from 1 to Infinity times
func (m *mExpectedMockGetGlobulaStateHash) Expect() *mExpectedMockGetGlobulaStateHash {
	m.mock.GetGlobulaStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetGlobulaStateHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetGlobulaStateHash
func (m *mExpectedMockGetGlobulaStateHash) Return(r proofs.GlobulaStateHash) *ExpectedMock {
	m.mock.GetGlobulaStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetGlobulaStateHashExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetGlobulaStateHashResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetGlobulaStateHash is expected once
func (m *mExpectedMockGetGlobulaStateHash) ExpectOnce() *ExpectedMockGetGlobulaStateHashExpectation {
	m.mock.GetGlobulaStateHashFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetGlobulaStateHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetGlobulaStateHashExpectation) Return(r proofs.GlobulaStateHash) {
	e.result = &ExpectedMockGetGlobulaStateHashResult{r}
}

//Set uses given function f as a mock of Expected.GetGlobulaStateHash method
func (m *mExpectedMockGetGlobulaStateHash) Set(f func() (r proofs.GlobulaStateHash)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetGlobulaStateHashFunc = f
	return m.mock
}

//GetGlobulaStateHash implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetGlobulaStateHash() (r proofs.GlobulaStateHash) {
	counter := atomic.AddUint64(&m.GetGlobulaStateHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobulaStateHashCounter, 1)

	if len(m.GetGlobulaStateHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetGlobulaStateHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetGlobulaStateHash.")
			return
		}

		result := m.GetGlobulaStateHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetGlobulaStateHash")
			return
		}

		r = result.r

		return
	}

	if m.GetGlobulaStateHashMock.mainExpectation != nil {

		result := m.GetGlobulaStateHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetGlobulaStateHash")
		}

		r = result.r

		return
	}

	if m.GetGlobulaStateHashFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetGlobulaStateHash.")
		return
	}

	return m.GetGlobulaStateHashFunc()
}

//GetGlobulaStateHashMinimockCounter returns a count of ExpectedMock.GetGlobulaStateHashFunc invocations
func (m *ExpectedMock) GetGlobulaStateHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobulaStateHashCounter)
}

//GetGlobulaStateHashMinimockPreCounter returns the value of ExpectedMock.GetGlobulaStateHash invocations
func (m *ExpectedMock) GetGlobulaStateHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobulaStateHashPreCounter)
}

//GetGlobulaStateHashFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetGlobulaStateHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetGlobulaStateHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetGlobulaStateHashCounter) == uint64(len(m.GetGlobulaStateHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetGlobulaStateHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetGlobulaStateHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetGlobulaStateHashFunc != nil {
		return atomic.LoadUint64(&m.GetGlobulaStateHashCounter) > 0
	}

	return true
}

type mExpectedMockGetMandateRegistry struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetMandateRegistryExpectation
	expectationSeries []*ExpectedMockGetMandateRegistryExpectation
}

type ExpectedMockGetMandateRegistryExpectation struct {
	result *ExpectedMockGetMandateRegistryResult
}

type ExpectedMockGetMandateRegistryResult struct {
	r MandateRegistry
}

//Expect specifies that invocation of Expected.GetMandateRegistry is expected from 1 to Infinity times
func (m *mExpectedMockGetMandateRegistry) Expect() *mExpectedMockGetMandateRegistry {
	m.mock.GetMandateRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetMandateRegistryExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetMandateRegistry
func (m *mExpectedMockGetMandateRegistry) Return(r MandateRegistry) *ExpectedMock {
	m.mock.GetMandateRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetMandateRegistryExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetMandateRegistryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetMandateRegistry is expected once
func (m *mExpectedMockGetMandateRegistry) ExpectOnce() *ExpectedMockGetMandateRegistryExpectation {
	m.mock.GetMandateRegistryFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetMandateRegistryExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetMandateRegistryExpectation) Return(r MandateRegistry) {
	e.result = &ExpectedMockGetMandateRegistryResult{r}
}

//Set uses given function f as a mock of Expected.GetMandateRegistry method
func (m *mExpectedMockGetMandateRegistry) Set(f func() (r MandateRegistry)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMandateRegistryFunc = f
	return m.mock
}

//GetMandateRegistry implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetMandateRegistry() (r MandateRegistry) {
	counter := atomic.AddUint64(&m.GetMandateRegistryPreCounter, 1)
	defer atomic.AddUint64(&m.GetMandateRegistryCounter, 1)

	if len(m.GetMandateRegistryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMandateRegistryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetMandateRegistry.")
			return
		}

		result := m.GetMandateRegistryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetMandateRegistry")
			return
		}

		r = result.r

		return
	}

	if m.GetMandateRegistryMock.mainExpectation != nil {

		result := m.GetMandateRegistryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetMandateRegistry")
		}

		r = result.r

		return
	}

	if m.GetMandateRegistryFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetMandateRegistry.")
		return
	}

	return m.GetMandateRegistryFunc()
}

//GetMandateRegistryMinimockCounter returns a count of ExpectedMock.GetMandateRegistryFunc invocations
func (m *ExpectedMock) GetMandateRegistryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMandateRegistryCounter)
}

//GetMandateRegistryMinimockPreCounter returns the value of ExpectedMock.GetMandateRegistry invocations
func (m *ExpectedMock) GetMandateRegistryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMandateRegistryPreCounter)
}

//GetMandateRegistryFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetMandateRegistryFinished() bool {
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

type mExpectedMockGetMisbehaviorRegistry struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetMisbehaviorRegistryExpectation
	expectationSeries []*ExpectedMockGetMisbehaviorRegistryExpectation
}

type ExpectedMockGetMisbehaviorRegistryExpectation struct {
	result *ExpectedMockGetMisbehaviorRegistryResult
}

type ExpectedMockGetMisbehaviorRegistryResult struct {
	r MisbehaviorRegistry
}

//Expect specifies that invocation of Expected.GetMisbehaviorRegistry is expected from 1 to Infinity times
func (m *mExpectedMockGetMisbehaviorRegistry) Expect() *mExpectedMockGetMisbehaviorRegistry {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetMisbehaviorRegistryExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetMisbehaviorRegistry
func (m *mExpectedMockGetMisbehaviorRegistry) Return(r MisbehaviorRegistry) *ExpectedMock {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetMisbehaviorRegistryExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetMisbehaviorRegistryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetMisbehaviorRegistry is expected once
func (m *mExpectedMockGetMisbehaviorRegistry) ExpectOnce() *ExpectedMockGetMisbehaviorRegistryExpectation {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetMisbehaviorRegistryExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetMisbehaviorRegistryExpectation) Return(r MisbehaviorRegistry) {
	e.result = &ExpectedMockGetMisbehaviorRegistryResult{r}
}

//Set uses given function f as a mock of Expected.GetMisbehaviorRegistry method
func (m *mExpectedMockGetMisbehaviorRegistry) Set(f func() (r MisbehaviorRegistry)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMisbehaviorRegistryFunc = f
	return m.mock
}

//GetMisbehaviorRegistry implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetMisbehaviorRegistry() (r MisbehaviorRegistry) {
	counter := atomic.AddUint64(&m.GetMisbehaviorRegistryPreCounter, 1)
	defer atomic.AddUint64(&m.GetMisbehaviorRegistryCounter, 1)

	if len(m.GetMisbehaviorRegistryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMisbehaviorRegistryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetMisbehaviorRegistry.")
			return
		}

		result := m.GetMisbehaviorRegistryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetMisbehaviorRegistry")
			return
		}

		r = result.r

		return
	}

	if m.GetMisbehaviorRegistryMock.mainExpectation != nil {

		result := m.GetMisbehaviorRegistryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetMisbehaviorRegistry")
		}

		r = result.r

		return
	}

	if m.GetMisbehaviorRegistryFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetMisbehaviorRegistry.")
		return
	}

	return m.GetMisbehaviorRegistryFunc()
}

//GetMisbehaviorRegistryMinimockCounter returns a count of ExpectedMock.GetMisbehaviorRegistryFunc invocations
func (m *ExpectedMock) GetMisbehaviorRegistryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMisbehaviorRegistryCounter)
}

//GetMisbehaviorRegistryMinimockPreCounter returns the value of ExpectedMock.GetMisbehaviorRegistry invocations
func (m *ExpectedMock) GetMisbehaviorRegistryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMisbehaviorRegistryPreCounter)
}

//GetMisbehaviorRegistryFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetMisbehaviorRegistryFinished() bool {
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

type mExpectedMockGetOfflinePopulation struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetOfflinePopulationExpectation
	expectationSeries []*ExpectedMockGetOfflinePopulationExpectation
}

type ExpectedMockGetOfflinePopulationExpectation struct {
	result *ExpectedMockGetOfflinePopulationResult
}

type ExpectedMockGetOfflinePopulationResult struct {
	r OfflinePopulation
}

//Expect specifies that invocation of Expected.GetOfflinePopulation is expected from 1 to Infinity times
func (m *mExpectedMockGetOfflinePopulation) Expect() *mExpectedMockGetOfflinePopulation {
	m.mock.GetOfflinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetOfflinePopulationExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetOfflinePopulation
func (m *mExpectedMockGetOfflinePopulation) Return(r OfflinePopulation) *ExpectedMock {
	m.mock.GetOfflinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetOfflinePopulationExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetOfflinePopulationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetOfflinePopulation is expected once
func (m *mExpectedMockGetOfflinePopulation) ExpectOnce() *ExpectedMockGetOfflinePopulationExpectation {
	m.mock.GetOfflinePopulationFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetOfflinePopulationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetOfflinePopulationExpectation) Return(r OfflinePopulation) {
	e.result = &ExpectedMockGetOfflinePopulationResult{r}
}

//Set uses given function f as a mock of Expected.GetOfflinePopulation method
func (m *mExpectedMockGetOfflinePopulation) Set(f func() (r OfflinePopulation)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOfflinePopulationFunc = f
	return m.mock
}

//GetOfflinePopulation implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetOfflinePopulation() (r OfflinePopulation) {
	counter := atomic.AddUint64(&m.GetOfflinePopulationPreCounter, 1)
	defer atomic.AddUint64(&m.GetOfflinePopulationCounter, 1)

	if len(m.GetOfflinePopulationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOfflinePopulationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetOfflinePopulation.")
			return
		}

		result := m.GetOfflinePopulationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetOfflinePopulation")
			return
		}

		r = result.r

		return
	}

	if m.GetOfflinePopulationMock.mainExpectation != nil {

		result := m.GetOfflinePopulationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetOfflinePopulation")
		}

		r = result.r

		return
	}

	if m.GetOfflinePopulationFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetOfflinePopulation.")
		return
	}

	return m.GetOfflinePopulationFunc()
}

//GetOfflinePopulationMinimockCounter returns a count of ExpectedMock.GetOfflinePopulationFunc invocations
func (m *ExpectedMock) GetOfflinePopulationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOfflinePopulationCounter)
}

//GetOfflinePopulationMinimockPreCounter returns the value of ExpectedMock.GetOfflinePopulation invocations
func (m *ExpectedMock) GetOfflinePopulationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOfflinePopulationPreCounter)
}

//GetOfflinePopulationFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetOfflinePopulationFinished() bool {
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

type mExpectedMockGetOnlinePopulation struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetOnlinePopulationExpectation
	expectationSeries []*ExpectedMockGetOnlinePopulationExpectation
}

type ExpectedMockGetOnlinePopulationExpectation struct {
	result *ExpectedMockGetOnlinePopulationResult
}

type ExpectedMockGetOnlinePopulationResult struct {
	r OnlinePopulation
}

//Expect specifies that invocation of Expected.GetOnlinePopulation is expected from 1 to Infinity times
func (m *mExpectedMockGetOnlinePopulation) Expect() *mExpectedMockGetOnlinePopulation {
	m.mock.GetOnlinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetOnlinePopulationExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetOnlinePopulation
func (m *mExpectedMockGetOnlinePopulation) Return(r OnlinePopulation) *ExpectedMock {
	m.mock.GetOnlinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetOnlinePopulationExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetOnlinePopulationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetOnlinePopulation is expected once
func (m *mExpectedMockGetOnlinePopulation) ExpectOnce() *ExpectedMockGetOnlinePopulationExpectation {
	m.mock.GetOnlinePopulationFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetOnlinePopulationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetOnlinePopulationExpectation) Return(r OnlinePopulation) {
	e.result = &ExpectedMockGetOnlinePopulationResult{r}
}

//Set uses given function f as a mock of Expected.GetOnlinePopulation method
func (m *mExpectedMockGetOnlinePopulation) Set(f func() (r OnlinePopulation)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOnlinePopulationFunc = f
	return m.mock
}

//GetOnlinePopulation implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetOnlinePopulation() (r OnlinePopulation) {
	counter := atomic.AddUint64(&m.GetOnlinePopulationPreCounter, 1)
	defer atomic.AddUint64(&m.GetOnlinePopulationCounter, 1)

	if len(m.GetOnlinePopulationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOnlinePopulationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetOnlinePopulation.")
			return
		}

		result := m.GetOnlinePopulationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetOnlinePopulation")
			return
		}

		r = result.r

		return
	}

	if m.GetOnlinePopulationMock.mainExpectation != nil {

		result := m.GetOnlinePopulationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetOnlinePopulation")
		}

		r = result.r

		return
	}

	if m.GetOnlinePopulationFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetOnlinePopulation.")
		return
	}

	return m.GetOnlinePopulationFunc()
}

//GetOnlinePopulationMinimockCounter returns a count of ExpectedMock.GetOnlinePopulationFunc invocations
func (m *ExpectedMock) GetOnlinePopulationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOnlinePopulationCounter)
}

//GetOnlinePopulationMinimockPreCounter returns the value of ExpectedMock.GetOnlinePopulation invocations
func (m *ExpectedMock) GetOnlinePopulationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOnlinePopulationPreCounter)
}

//GetOnlinePopulationFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetOnlinePopulationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOnlinePopulationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOnlinePopulationCounter) == uint64(len(m.GetOnlinePopulationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOnlinePopulationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOnlinePopulationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOnlinePopulationFunc != nil {
		return atomic.LoadUint64(&m.GetOnlinePopulationCounter) > 0
	}

	return true
}

type mExpectedMockGetPrevious struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetPreviousExpectation
	expectationSeries []*ExpectedMockGetPreviousExpectation
}

type ExpectedMockGetPreviousExpectation struct {
	result *ExpectedMockGetPreviousResult
}

type ExpectedMockGetPreviousResult struct {
	r Active
}

//Expect specifies that invocation of Expected.GetPrevious is expected from 1 to Infinity times
func (m *mExpectedMockGetPrevious) Expect() *mExpectedMockGetPrevious {
	m.mock.GetPreviousFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetPreviousExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetPrevious
func (m *mExpectedMockGetPrevious) Return(r Active) *ExpectedMock {
	m.mock.GetPreviousFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetPreviousExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetPreviousResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetPrevious is expected once
func (m *mExpectedMockGetPrevious) ExpectOnce() *ExpectedMockGetPreviousExpectation {
	m.mock.GetPreviousFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetPreviousExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetPreviousExpectation) Return(r Active) {
	e.result = &ExpectedMockGetPreviousResult{r}
}

//Set uses given function f as a mock of Expected.GetPrevious method
func (m *mExpectedMockGetPrevious) Set(f func() (r Active)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPreviousFunc = f
	return m.mock
}

//GetPrevious implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetPrevious() (r Active) {
	counter := atomic.AddUint64(&m.GetPreviousPreCounter, 1)
	defer atomic.AddUint64(&m.GetPreviousCounter, 1)

	if len(m.GetPreviousMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPreviousMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetPrevious.")
			return
		}

		result := m.GetPreviousMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetPrevious")
			return
		}

		r = result.r

		return
	}

	if m.GetPreviousMock.mainExpectation != nil {

		result := m.GetPreviousMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetPrevious")
		}

		r = result.r

		return
	}

	if m.GetPreviousFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetPrevious.")
		return
	}

	return m.GetPreviousFunc()
}

//GetPreviousMinimockCounter returns a count of ExpectedMock.GetPreviousFunc invocations
func (m *ExpectedMock) GetPreviousMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreviousCounter)
}

//GetPreviousMinimockPreCounter returns the value of ExpectedMock.GetPrevious invocations
func (m *ExpectedMock) GetPreviousMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreviousPreCounter)
}

//GetPreviousFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetPreviousFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPreviousMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPreviousCounter) == uint64(len(m.GetPreviousMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPreviousMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPreviousCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPreviousFunc != nil {
		return atomic.LoadUint64(&m.GetPreviousCounter) > 0
	}

	return true
}

type mExpectedMockGetProfileFactory struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetProfileFactoryExpectation
	expectationSeries []*ExpectedMockGetProfileFactoryExpectation
}

type ExpectedMockGetProfileFactoryExpectation struct {
	input  *ExpectedMockGetProfileFactoryInput
	result *ExpectedMockGetProfileFactoryResult
}

type ExpectedMockGetProfileFactoryInput struct {
	p cryptkit.KeyStoreFactory
}

type ExpectedMockGetProfileFactoryResult struct {
	r profiles.Factory
}

//Expect specifies that invocation of Expected.GetProfileFactory is expected from 1 to Infinity times
func (m *mExpectedMockGetProfileFactory) Expect(p cryptkit.KeyStoreFactory) *mExpectedMockGetProfileFactory {
	m.mock.GetProfileFactoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetProfileFactoryExpectation{}
	}
	m.mainExpectation.input = &ExpectedMockGetProfileFactoryInput{p}
	return m
}

//Return specifies results of invocation of Expected.GetProfileFactory
func (m *mExpectedMockGetProfileFactory) Return(r profiles.Factory) *ExpectedMock {
	m.mock.GetProfileFactoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetProfileFactoryExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetProfileFactoryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetProfileFactory is expected once
func (m *mExpectedMockGetProfileFactory) ExpectOnce(p cryptkit.KeyStoreFactory) *ExpectedMockGetProfileFactoryExpectation {
	m.mock.GetProfileFactoryFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetProfileFactoryExpectation{}
	expectation.input = &ExpectedMockGetProfileFactoryInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetProfileFactoryExpectation) Return(r profiles.Factory) {
	e.result = &ExpectedMockGetProfileFactoryResult{r}
}

//Set uses given function f as a mock of Expected.GetProfileFactory method
func (m *mExpectedMockGetProfileFactory) Set(f func(p cryptkit.KeyStoreFactory) (r profiles.Factory)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetProfileFactoryFunc = f
	return m.mock
}

//GetProfileFactory implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetProfileFactory(p cryptkit.KeyStoreFactory) (r profiles.Factory) {
	counter := atomic.AddUint64(&m.GetProfileFactoryPreCounter, 1)
	defer atomic.AddUint64(&m.GetProfileFactoryCounter, 1)

	if len(m.GetProfileFactoryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetProfileFactoryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetProfileFactory. %v", p)
			return
		}

		input := m.GetProfileFactoryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExpectedMockGetProfileFactoryInput{p}, "Expected.GetProfileFactory got unexpected parameters")

		result := m.GetProfileFactoryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetProfileFactory")
			return
		}

		r = result.r

		return
	}

	if m.GetProfileFactoryMock.mainExpectation != nil {

		input := m.GetProfileFactoryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExpectedMockGetProfileFactoryInput{p}, "Expected.GetProfileFactory got unexpected parameters")
		}

		result := m.GetProfileFactoryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetProfileFactory")
		}

		r = result.r

		return
	}

	if m.GetProfileFactoryFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetProfileFactory. %v", p)
		return
	}

	return m.GetProfileFactoryFunc(p)
}

//GetProfileFactoryMinimockCounter returns a count of ExpectedMock.GetProfileFactoryFunc invocations
func (m *ExpectedMock) GetProfileFactoryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetProfileFactoryCounter)
}

//GetProfileFactoryMinimockPreCounter returns the value of ExpectedMock.GetProfileFactory invocations
func (m *ExpectedMock) GetProfileFactoryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetProfileFactoryPreCounter)
}

//GetProfileFactoryFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetProfileFactoryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetProfileFactoryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetProfileFactoryCounter) == uint64(len(m.GetProfileFactoryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetProfileFactoryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetProfileFactoryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetProfileFactoryFunc != nil {
		return atomic.LoadUint64(&m.GetProfileFactoryCounter) > 0
	}

	return true
}

type mExpectedMockGetPulseNumber struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockGetPulseNumberExpectation
	expectationSeries []*ExpectedMockGetPulseNumberExpectation
}

type ExpectedMockGetPulseNumberExpectation struct {
	result *ExpectedMockGetPulseNumberResult
}

type ExpectedMockGetPulseNumberResult struct {
	r pulse.Number
}

//Expect specifies that invocation of Expected.GetPulseNumber is expected from 1 to Infinity times
func (m *mExpectedMockGetPulseNumber) Expect() *mExpectedMockGetPulseNumber {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetPulseNumberExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.GetPulseNumber
func (m *mExpectedMockGetPulseNumber) Return(r pulse.Number) *ExpectedMock {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockGetPulseNumberExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockGetPulseNumberResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.GetPulseNumber is expected once
func (m *mExpectedMockGetPulseNumber) ExpectOnce() *ExpectedMockGetPulseNumberExpectation {
	m.mock.GetPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockGetPulseNumberExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockGetPulseNumberExpectation) Return(r pulse.Number) {
	e.result = &ExpectedMockGetPulseNumberResult{r}
}

//Set uses given function f as a mock of Expected.GetPulseNumber method
func (m *mExpectedMockGetPulseNumber) Set(f func() (r pulse.Number)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseNumberFunc = f
	return m.mock
}

//GetPulseNumber implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) GetPulseNumber() (r pulse.Number) {
	counter := atomic.AddUint64(&m.GetPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseNumberCounter, 1)

	if len(m.GetPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.GetPulseNumber.")
			return
		}

		result := m.GetPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetPulseNumber")
			return
		}

		r = result.r

		return
	}

	if m.GetPulseNumberMock.mainExpectation != nil {

		result := m.GetPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.GetPulseNumber")
		}

		r = result.r

		return
	}

	if m.GetPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.GetPulseNumber.")
		return
	}

	return m.GetPulseNumberFunc()
}

//GetPulseNumberMinimockCounter returns a count of ExpectedMock.GetPulseNumberFunc invocations
func (m *ExpectedMock) GetPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberCounter)
}

//GetPulseNumberMinimockPreCounter returns the value of ExpectedMock.GetPulseNumber invocations
func (m *ExpectedMock) GetPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberPreCounter)
}

//GetPulseNumberFinished returns true if mock invocations count is ok
func (m *ExpectedMock) GetPulseNumberFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulseNumberMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) == uint64(len(m.GetPulseNumberMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulseNumberMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulseNumberFunc != nil {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) > 0
	}

	return true
}

type mExpectedMockIsActive struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockIsActiveExpectation
	expectationSeries []*ExpectedMockIsActiveExpectation
}

type ExpectedMockIsActiveExpectation struct {
	result *ExpectedMockIsActiveResult
}

type ExpectedMockIsActiveResult struct {
	r bool
}

//Expect specifies that invocation of Expected.IsActive is expected from 1 to Infinity times
func (m *mExpectedMockIsActive) Expect() *mExpectedMockIsActive {
	m.mock.IsActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockIsActiveExpectation{}
	}

	return m
}

//Return specifies results of invocation of Expected.IsActive
func (m *mExpectedMockIsActive) Return(r bool) *ExpectedMock {
	m.mock.IsActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockIsActiveExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockIsActiveResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.IsActive is expected once
func (m *mExpectedMockIsActive) ExpectOnce() *ExpectedMockIsActiveExpectation {
	m.mock.IsActiveFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockIsActiveExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockIsActiveExpectation) Return(r bool) {
	e.result = &ExpectedMockIsActiveResult{r}
}

//Set uses given function f as a mock of Expected.IsActive method
func (m *mExpectedMockIsActive) Set(f func() (r bool)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsActiveFunc = f
	return m.mock
}

//IsActive implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) IsActive() (r bool) {
	counter := atomic.AddUint64(&m.IsActivePreCounter, 1)
	defer atomic.AddUint64(&m.IsActiveCounter, 1)

	if len(m.IsActiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsActiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.IsActive.")
			return
		}

		result := m.IsActiveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.IsActive")
			return
		}

		r = result.r

		return
	}

	if m.IsActiveMock.mainExpectation != nil {

		result := m.IsActiveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.IsActive")
		}

		r = result.r

		return
	}

	if m.IsActiveFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.IsActive.")
		return
	}

	return m.IsActiveFunc()
}

//IsActiveMinimockCounter returns a count of ExpectedMock.IsActiveFunc invocations
func (m *ExpectedMock) IsActiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsActiveCounter)
}

//IsActiveMinimockPreCounter returns the value of ExpectedMock.IsActive invocations
func (m *ExpectedMock) IsActiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsActivePreCounter)
}

//IsActiveFinished returns true if mock invocations count is ok
func (m *ExpectedMock) IsActiveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsActiveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsActiveCounter) == uint64(len(m.IsActiveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsActiveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsActiveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsActiveFunc != nil {
		return atomic.LoadUint64(&m.IsActiveCounter) > 0
	}

	return true
}

type mExpectedMockMakeActive struct {
	mock              *ExpectedMock
	mainExpectation   *ExpectedMockMakeActiveExpectation
	expectationSeries []*ExpectedMockMakeActiveExpectation
}

type ExpectedMockMakeActiveExpectation struct {
	input  *ExpectedMockMakeActiveInput
	result *ExpectedMockMakeActiveResult
}

type ExpectedMockMakeActiveInput struct {
	p pulse.Data
}

type ExpectedMockMakeActiveResult struct {
	r Active
}

//Expect specifies that invocation of Expected.MakeActive is expected from 1 to Infinity times
func (m *mExpectedMockMakeActive) Expect(p pulse.Data) *mExpectedMockMakeActive {
	m.mock.MakeActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockMakeActiveExpectation{}
	}
	m.mainExpectation.input = &ExpectedMockMakeActiveInput{p}
	return m
}

//Return specifies results of invocation of Expected.MakeActive
func (m *mExpectedMockMakeActive) Return(r Active) *ExpectedMock {
	m.mock.MakeActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExpectedMockMakeActiveExpectation{}
	}
	m.mainExpectation.result = &ExpectedMockMakeActiveResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Expected.MakeActive is expected once
func (m *mExpectedMockMakeActive) ExpectOnce(p pulse.Data) *ExpectedMockMakeActiveExpectation {
	m.mock.MakeActiveFunc = nil
	m.mainExpectation = nil

	expectation := &ExpectedMockMakeActiveExpectation{}
	expectation.input = &ExpectedMockMakeActiveInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExpectedMockMakeActiveExpectation) Return(r Active) {
	e.result = &ExpectedMockMakeActiveResult{r}
}

//Set uses given function f as a mock of Expected.MakeActive method
func (m *mExpectedMockMakeActive) Set(f func(p pulse.Data) (r Active)) *ExpectedMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MakeActiveFunc = f
	return m.mock
}

//MakeActive implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected interface
func (m *ExpectedMock) MakeActive(p pulse.Data) (r Active) {
	counter := atomic.AddUint64(&m.MakeActivePreCounter, 1)
	defer atomic.AddUint64(&m.MakeActiveCounter, 1)

	if len(m.MakeActiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MakeActiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExpectedMock.MakeActive. %v", p)
			return
		}

		input := m.MakeActiveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExpectedMockMakeActiveInput{p}, "Expected.MakeActive got unexpected parameters")

		result := m.MakeActiveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.MakeActive")
			return
		}

		r = result.r

		return
	}

	if m.MakeActiveMock.mainExpectation != nil {

		input := m.MakeActiveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExpectedMockMakeActiveInput{p}, "Expected.MakeActive got unexpected parameters")
		}

		result := m.MakeActiveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExpectedMock.MakeActive")
		}

		r = result.r

		return
	}

	if m.MakeActiveFunc == nil {
		m.t.Fatalf("Unexpected call to ExpectedMock.MakeActive. %v", p)
		return
	}

	return m.MakeActiveFunc(p)
}

//MakeActiveMinimockCounter returns a count of ExpectedMock.MakeActiveFunc invocations
func (m *ExpectedMock) MakeActiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MakeActiveCounter)
}

//MakeActiveMinimockPreCounter returns the value of ExpectedMock.MakeActive invocations
func (m *ExpectedMock) MakeActiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MakeActivePreCounter)
}

//MakeActiveFinished returns true if mock invocations count is ok
func (m *ExpectedMock) MakeActiveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MakeActiveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MakeActiveCounter) == uint64(len(m.MakeActiveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MakeActiveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MakeActiveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MakeActiveFunc != nil {
		return atomic.LoadUint64(&m.MakeActiveCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExpectedMock) ValidateCallCounters() {

	if !m.CreateBuilderFinished() {
		m.t.Fatal("Expected call to ExpectedMock.CreateBuilder")
	}

	if !m.GetCensusStateFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetCensusState")
	}

	if !m.GetCloudStateHashFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetCloudStateHash")
	}

	if !m.GetEvictedPopulationFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetEvictedPopulation")
	}

	if !m.GetExpectedPulseNumberFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetExpectedPulseNumber")
	}

	if !m.GetGlobulaStateHashFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetGlobulaStateHash")
	}

	if !m.GetMandateRegistryFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetMandateRegistry")
	}

	if !m.GetMisbehaviorRegistryFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetMisbehaviorRegistry")
	}

	if !m.GetOfflinePopulationFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetOfflinePopulation")
	}

	if !m.GetOnlinePopulationFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetOnlinePopulation")
	}

	if !m.GetPreviousFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetPrevious")
	}

	if !m.GetProfileFactoryFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetProfileFactory")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetPulseNumber")
	}

	if !m.IsActiveFinished() {
		m.t.Fatal("Expected call to ExpectedMock.IsActive")
	}

	if !m.MakeActiveFinished() {
		m.t.Fatal("Expected call to ExpectedMock.MakeActive")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExpectedMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ExpectedMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ExpectedMock) MinimockFinish() {

	if !m.CreateBuilderFinished() {
		m.t.Fatal("Expected call to ExpectedMock.CreateBuilder")
	}

	if !m.GetCensusStateFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetCensusState")
	}

	if !m.GetCloudStateHashFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetCloudStateHash")
	}

	if !m.GetEvictedPopulationFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetEvictedPopulation")
	}

	if !m.GetExpectedPulseNumberFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetExpectedPulseNumber")
	}

	if !m.GetGlobulaStateHashFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetGlobulaStateHash")
	}

	if !m.GetMandateRegistryFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetMandateRegistry")
	}

	if !m.GetMisbehaviorRegistryFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetMisbehaviorRegistry")
	}

	if !m.GetOfflinePopulationFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetOfflinePopulation")
	}

	if !m.GetOnlinePopulationFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetOnlinePopulation")
	}

	if !m.GetPreviousFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetPrevious")
	}

	if !m.GetProfileFactoryFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetProfileFactory")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to ExpectedMock.GetPulseNumber")
	}

	if !m.IsActiveFinished() {
		m.t.Fatal("Expected call to ExpectedMock.IsActive")
	}

	if !m.MakeActiveFinished() {
		m.t.Fatal("Expected call to ExpectedMock.MakeActive")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ExpectedMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ExpectedMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CreateBuilderFinished()
		ok = ok && m.GetCensusStateFinished()
		ok = ok && m.GetCloudStateHashFinished()
		ok = ok && m.GetEvictedPopulationFinished()
		ok = ok && m.GetExpectedPulseNumberFinished()
		ok = ok && m.GetGlobulaStateHashFinished()
		ok = ok && m.GetMandateRegistryFinished()
		ok = ok && m.GetMisbehaviorRegistryFinished()
		ok = ok && m.GetOfflinePopulationFinished()
		ok = ok && m.GetOnlinePopulationFinished()
		ok = ok && m.GetPreviousFinished()
		ok = ok && m.GetProfileFactoryFinished()
		ok = ok && m.GetPulseNumberFinished()
		ok = ok && m.IsActiveFinished()
		ok = ok && m.MakeActiveFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CreateBuilderFinished() {
				m.t.Error("Expected call to ExpectedMock.CreateBuilder")
			}

			if !m.GetCensusStateFinished() {
				m.t.Error("Expected call to ExpectedMock.GetCensusState")
			}

			if !m.GetCloudStateHashFinished() {
				m.t.Error("Expected call to ExpectedMock.GetCloudStateHash")
			}

			if !m.GetEvictedPopulationFinished() {
				m.t.Error("Expected call to ExpectedMock.GetEvictedPopulation")
			}

			if !m.GetExpectedPulseNumberFinished() {
				m.t.Error("Expected call to ExpectedMock.GetExpectedPulseNumber")
			}

			if !m.GetGlobulaStateHashFinished() {
				m.t.Error("Expected call to ExpectedMock.GetGlobulaStateHash")
			}

			if !m.GetMandateRegistryFinished() {
				m.t.Error("Expected call to ExpectedMock.GetMandateRegistry")
			}

			if !m.GetMisbehaviorRegistryFinished() {
				m.t.Error("Expected call to ExpectedMock.GetMisbehaviorRegistry")
			}

			if !m.GetOfflinePopulationFinished() {
				m.t.Error("Expected call to ExpectedMock.GetOfflinePopulation")
			}

			if !m.GetOnlinePopulationFinished() {
				m.t.Error("Expected call to ExpectedMock.GetOnlinePopulation")
			}

			if !m.GetPreviousFinished() {
				m.t.Error("Expected call to ExpectedMock.GetPrevious")
			}

			if !m.GetProfileFactoryFinished() {
				m.t.Error("Expected call to ExpectedMock.GetProfileFactory")
			}

			if !m.GetPulseNumberFinished() {
				m.t.Error("Expected call to ExpectedMock.GetPulseNumber")
			}

			if !m.IsActiveFinished() {
				m.t.Error("Expected call to ExpectedMock.IsActive")
			}

			if !m.MakeActiveFinished() {
				m.t.Error("Expected call to ExpectedMock.MakeActive")
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
func (m *ExpectedMock) AllMocksCalled() bool {

	if !m.CreateBuilderFinished() {
		return false
	}

	if !m.GetCensusStateFinished() {
		return false
	}

	if !m.GetCloudStateHashFinished() {
		return false
	}

	if !m.GetEvictedPopulationFinished() {
		return false
	}

	if !m.GetExpectedPulseNumberFinished() {
		return false
	}

	if !m.GetGlobulaStateHashFinished() {
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

	if !m.GetOnlinePopulationFinished() {
		return false
	}

	if !m.GetPreviousFinished() {
		return false
	}

	if !m.GetProfileFactoryFinished() {
		return false
	}

	if !m.GetPulseNumberFinished() {
		return false
	}

	if !m.IsActiveFinished() {
		return false
	}

	if !m.MakeActiveFinished() {
		return false
	}

	return true
}

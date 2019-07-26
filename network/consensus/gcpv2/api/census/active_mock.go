package census

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Active" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/census
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	pulse "github.com/insolar/insolar/network/consensus/common/pulse"
	profiles "github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	proofs "github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"

	testify_assert "github.com/stretchr/testify/assert"
)

//ActiveMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active
type ActiveMock struct {
	t minimock.Tester

	CreateBuilderFunc       func(p context.Context, p1 pulse.Number) (r Builder)
	CreateBuilderCounter    uint64
	CreateBuilderPreCounter uint64
	CreateBuilderMock       mActiveMockCreateBuilder

	GetCensusStateFunc       func() (r State)
	GetCensusStateCounter    uint64
	GetCensusStatePreCounter uint64
	GetCensusStateMock       mActiveMockGetCensusState

	GetCloudStateHashFunc       func() (r proofs.CloudStateHash)
	GetCloudStateHashCounter    uint64
	GetCloudStateHashPreCounter uint64
	GetCloudStateHashMock       mActiveMockGetCloudStateHash

	GetEvictedPopulationFunc       func() (r EvictedPopulation)
	GetEvictedPopulationCounter    uint64
	GetEvictedPopulationPreCounter uint64
	GetEvictedPopulationMock       mActiveMockGetEvictedPopulation

	GetExpectedPulseNumberFunc       func() (r pulse.Number)
	GetExpectedPulseNumberCounter    uint64
	GetExpectedPulseNumberPreCounter uint64
	GetExpectedPulseNumberMock       mActiveMockGetExpectedPulseNumber

	GetGlobulaStateHashFunc       func() (r proofs.GlobulaStateHash)
	GetGlobulaStateHashCounter    uint64
	GetGlobulaStateHashPreCounter uint64
	GetGlobulaStateHashMock       mActiveMockGetGlobulaStateHash

	GetMandateRegistryFunc       func() (r MandateRegistry)
	GetMandateRegistryCounter    uint64
	GetMandateRegistryPreCounter uint64
	GetMandateRegistryMock       mActiveMockGetMandateRegistry

	GetMisbehaviorRegistryFunc       func() (r MisbehaviorRegistry)
	GetMisbehaviorRegistryCounter    uint64
	GetMisbehaviorRegistryPreCounter uint64
	GetMisbehaviorRegistryMock       mActiveMockGetMisbehaviorRegistry

	GetOfflinePopulationFunc       func() (r OfflinePopulation)
	GetOfflinePopulationCounter    uint64
	GetOfflinePopulationPreCounter uint64
	GetOfflinePopulationMock       mActiveMockGetOfflinePopulation

	GetOnlinePopulationFunc       func() (r OnlinePopulation)
	GetOnlinePopulationCounter    uint64
	GetOnlinePopulationPreCounter uint64
	GetOnlinePopulationMock       mActiveMockGetOnlinePopulation

	GetProfileFactoryFunc       func(p cryptkit.KeyStoreFactory) (r profiles.Factory)
	GetProfileFactoryCounter    uint64
	GetProfileFactoryPreCounter uint64
	GetProfileFactoryMock       mActiveMockGetProfileFactory

	GetPulseDataFunc       func() (r pulse.Data)
	GetPulseDataCounter    uint64
	GetPulseDataPreCounter uint64
	GetPulseDataMock       mActiveMockGetPulseData

	GetPulseNumberFunc       func() (r pulse.Number)
	GetPulseNumberCounter    uint64
	GetPulseNumberPreCounter uint64
	GetPulseNumberMock       mActiveMockGetPulseNumber

	IsActiveFunc       func() (r bool)
	IsActiveCounter    uint64
	IsActivePreCounter uint64
	IsActiveMock       mActiveMockIsActive
}

//NewActiveMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active
func NewActiveMock(t minimock.Tester) *ActiveMock {
	m := &ActiveMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateBuilderMock = mActiveMockCreateBuilder{mock: m}
	m.GetCensusStateMock = mActiveMockGetCensusState{mock: m}
	m.GetCloudStateHashMock = mActiveMockGetCloudStateHash{mock: m}
	m.GetEvictedPopulationMock = mActiveMockGetEvictedPopulation{mock: m}
	m.GetExpectedPulseNumberMock = mActiveMockGetExpectedPulseNumber{mock: m}
	m.GetGlobulaStateHashMock = mActiveMockGetGlobulaStateHash{mock: m}
	m.GetMandateRegistryMock = mActiveMockGetMandateRegistry{mock: m}
	m.GetMisbehaviorRegistryMock = mActiveMockGetMisbehaviorRegistry{mock: m}
	m.GetOfflinePopulationMock = mActiveMockGetOfflinePopulation{mock: m}
	m.GetOnlinePopulationMock = mActiveMockGetOnlinePopulation{mock: m}
	m.GetProfileFactoryMock = mActiveMockGetProfileFactory{mock: m}
	m.GetPulseDataMock = mActiveMockGetPulseData{mock: m}
	m.GetPulseNumberMock = mActiveMockGetPulseNumber{mock: m}
	m.IsActiveMock = mActiveMockIsActive{mock: m}

	return m
}

type mActiveMockCreateBuilder struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockCreateBuilderExpectation
	expectationSeries []*ActiveMockCreateBuilderExpectation
}

type ActiveMockCreateBuilderExpectation struct {
	input  *ActiveMockCreateBuilderInput
	result *ActiveMockCreateBuilderResult
}

type ActiveMockCreateBuilderInput struct {
	p  context.Context
	p1 pulse.Number
}

type ActiveMockCreateBuilderResult struct {
	r Builder
}

//Expect specifies that invocation of Active.CreateBuilder is expected from 1 to Infinity times
func (m *mActiveMockCreateBuilder) Expect(p context.Context, p1 pulse.Number) *mActiveMockCreateBuilder {
	m.mock.CreateBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockCreateBuilderExpectation{}
	}
	m.mainExpectation.input = &ActiveMockCreateBuilderInput{p, p1}
	return m
}

//Return specifies results of invocation of Active.CreateBuilder
func (m *mActiveMockCreateBuilder) Return(r Builder) *ActiveMock {
	m.mock.CreateBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockCreateBuilderExpectation{}
	}
	m.mainExpectation.result = &ActiveMockCreateBuilderResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.CreateBuilder is expected once
func (m *mActiveMockCreateBuilder) ExpectOnce(p context.Context, p1 pulse.Number) *ActiveMockCreateBuilderExpectation {
	m.mock.CreateBuilderFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockCreateBuilderExpectation{}
	expectation.input = &ActiveMockCreateBuilderInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockCreateBuilderExpectation) Return(r Builder) {
	e.result = &ActiveMockCreateBuilderResult{r}
}

//Set uses given function f as a mock of Active.CreateBuilder method
func (m *mActiveMockCreateBuilder) Set(f func(p context.Context, p1 pulse.Number) (r Builder)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateBuilderFunc = f
	return m.mock
}

//CreateBuilder implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) CreateBuilder(p context.Context, p1 pulse.Number) (r Builder) {
	counter := atomic.AddUint64(&m.CreateBuilderPreCounter, 1)
	defer atomic.AddUint64(&m.CreateBuilderCounter, 1)

	if len(m.CreateBuilderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateBuilderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.CreateBuilder. %v %v", p, p1)
			return
		}

		input := m.CreateBuilderMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ActiveMockCreateBuilderInput{p, p1}, "Active.CreateBuilder got unexpected parameters")

		result := m.CreateBuilderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.CreateBuilder")
			return
		}

		r = result.r

		return
	}

	if m.CreateBuilderMock.mainExpectation != nil {

		input := m.CreateBuilderMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ActiveMockCreateBuilderInput{p, p1}, "Active.CreateBuilder got unexpected parameters")
		}

		result := m.CreateBuilderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.CreateBuilder")
		}

		r = result.r

		return
	}

	if m.CreateBuilderFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.CreateBuilder. %v %v", p, p1)
		return
	}

	return m.CreateBuilderFunc(p, p1)
}

//CreateBuilderMinimockCounter returns a count of ActiveMock.CreateBuilderFunc invocations
func (m *ActiveMock) CreateBuilderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateBuilderCounter)
}

//CreateBuilderMinimockPreCounter returns the value of ActiveMock.CreateBuilder invocations
func (m *ActiveMock) CreateBuilderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateBuilderPreCounter)
}

//CreateBuilderFinished returns true if mock invocations count is ok
func (m *ActiveMock) CreateBuilderFinished() bool {
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

type mActiveMockGetCensusState struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetCensusStateExpectation
	expectationSeries []*ActiveMockGetCensusStateExpectation
}

type ActiveMockGetCensusStateExpectation struct {
	result *ActiveMockGetCensusStateResult
}

type ActiveMockGetCensusStateResult struct {
	r State
}

//Expect specifies that invocation of Active.GetCensusState is expected from 1 to Infinity times
func (m *mActiveMockGetCensusState) Expect() *mActiveMockGetCensusState {
	m.mock.GetCensusStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetCensusStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetCensusState
func (m *mActiveMockGetCensusState) Return(r State) *ActiveMock {
	m.mock.GetCensusStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetCensusStateExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetCensusStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetCensusState is expected once
func (m *mActiveMockGetCensusState) ExpectOnce() *ActiveMockGetCensusStateExpectation {
	m.mock.GetCensusStateFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetCensusStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetCensusStateExpectation) Return(r State) {
	e.result = &ActiveMockGetCensusStateResult{r}
}

//Set uses given function f as a mock of Active.GetCensusState method
func (m *mActiveMockGetCensusState) Set(f func() (r State)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCensusStateFunc = f
	return m.mock
}

//GetCensusState implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetCensusState() (r State) {
	counter := atomic.AddUint64(&m.GetCensusStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetCensusStateCounter, 1)

	if len(m.GetCensusStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCensusStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetCensusState.")
			return
		}

		result := m.GetCensusStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetCensusState")
			return
		}

		r = result.r

		return
	}

	if m.GetCensusStateMock.mainExpectation != nil {

		result := m.GetCensusStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetCensusState")
		}

		r = result.r

		return
	}

	if m.GetCensusStateFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetCensusState.")
		return
	}

	return m.GetCensusStateFunc()
}

//GetCensusStateMinimockCounter returns a count of ActiveMock.GetCensusStateFunc invocations
func (m *ActiveMock) GetCensusStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCensusStateCounter)
}

//GetCensusStateMinimockPreCounter returns the value of ActiveMock.GetCensusState invocations
func (m *ActiveMock) GetCensusStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCensusStatePreCounter)
}

//GetCensusStateFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetCensusStateFinished() bool {
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

type mActiveMockGetCloudStateHash struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetCloudStateHashExpectation
	expectationSeries []*ActiveMockGetCloudStateHashExpectation
}

type ActiveMockGetCloudStateHashExpectation struct {
	result *ActiveMockGetCloudStateHashResult
}

type ActiveMockGetCloudStateHashResult struct {
	r proofs.CloudStateHash
}

//Expect specifies that invocation of Active.GetCloudStateHash is expected from 1 to Infinity times
func (m *mActiveMockGetCloudStateHash) Expect() *mActiveMockGetCloudStateHash {
	m.mock.GetCloudStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetCloudStateHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetCloudStateHash
func (m *mActiveMockGetCloudStateHash) Return(r proofs.CloudStateHash) *ActiveMock {
	m.mock.GetCloudStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetCloudStateHashExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetCloudStateHashResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetCloudStateHash is expected once
func (m *mActiveMockGetCloudStateHash) ExpectOnce() *ActiveMockGetCloudStateHashExpectation {
	m.mock.GetCloudStateHashFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetCloudStateHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetCloudStateHashExpectation) Return(r proofs.CloudStateHash) {
	e.result = &ActiveMockGetCloudStateHashResult{r}
}

//Set uses given function f as a mock of Active.GetCloudStateHash method
func (m *mActiveMockGetCloudStateHash) Set(f func() (r proofs.CloudStateHash)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCloudStateHashFunc = f
	return m.mock
}

//GetCloudStateHash implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetCloudStateHash() (r proofs.CloudStateHash) {
	counter := atomic.AddUint64(&m.GetCloudStateHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetCloudStateHashCounter, 1)

	if len(m.GetCloudStateHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCloudStateHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetCloudStateHash.")
			return
		}

		result := m.GetCloudStateHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetCloudStateHash")
			return
		}

		r = result.r

		return
	}

	if m.GetCloudStateHashMock.mainExpectation != nil {

		result := m.GetCloudStateHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetCloudStateHash")
		}

		r = result.r

		return
	}

	if m.GetCloudStateHashFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetCloudStateHash.")
		return
	}

	return m.GetCloudStateHashFunc()
}

//GetCloudStateHashMinimockCounter returns a count of ActiveMock.GetCloudStateHashFunc invocations
func (m *ActiveMock) GetCloudStateHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudStateHashCounter)
}

//GetCloudStateHashMinimockPreCounter returns the value of ActiveMock.GetCloudStateHash invocations
func (m *ActiveMock) GetCloudStateHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudStateHashPreCounter)
}

//GetCloudStateHashFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetCloudStateHashFinished() bool {
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

type mActiveMockGetEvictedPopulation struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetEvictedPopulationExpectation
	expectationSeries []*ActiveMockGetEvictedPopulationExpectation
}

type ActiveMockGetEvictedPopulationExpectation struct {
	result *ActiveMockGetEvictedPopulationResult
}

type ActiveMockGetEvictedPopulationResult struct {
	r EvictedPopulation
}

//Expect specifies that invocation of Active.GetEvictedPopulation is expected from 1 to Infinity times
func (m *mActiveMockGetEvictedPopulation) Expect() *mActiveMockGetEvictedPopulation {
	m.mock.GetEvictedPopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetEvictedPopulationExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetEvictedPopulation
func (m *mActiveMockGetEvictedPopulation) Return(r EvictedPopulation) *ActiveMock {
	m.mock.GetEvictedPopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetEvictedPopulationExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetEvictedPopulationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetEvictedPopulation is expected once
func (m *mActiveMockGetEvictedPopulation) ExpectOnce() *ActiveMockGetEvictedPopulationExpectation {
	m.mock.GetEvictedPopulationFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetEvictedPopulationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetEvictedPopulationExpectation) Return(r EvictedPopulation) {
	e.result = &ActiveMockGetEvictedPopulationResult{r}
}

//Set uses given function f as a mock of Active.GetEvictedPopulation method
func (m *mActiveMockGetEvictedPopulation) Set(f func() (r EvictedPopulation)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetEvictedPopulationFunc = f
	return m.mock
}

//GetEvictedPopulation implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetEvictedPopulation() (r EvictedPopulation) {
	counter := atomic.AddUint64(&m.GetEvictedPopulationPreCounter, 1)
	defer atomic.AddUint64(&m.GetEvictedPopulationCounter, 1)

	if len(m.GetEvictedPopulationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetEvictedPopulationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetEvictedPopulation.")
			return
		}

		result := m.GetEvictedPopulationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetEvictedPopulation")
			return
		}

		r = result.r

		return
	}

	if m.GetEvictedPopulationMock.mainExpectation != nil {

		result := m.GetEvictedPopulationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetEvictedPopulation")
		}

		r = result.r

		return
	}

	if m.GetEvictedPopulationFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetEvictedPopulation.")
		return
	}

	return m.GetEvictedPopulationFunc()
}

//GetEvictedPopulationMinimockCounter returns a count of ActiveMock.GetEvictedPopulationFunc invocations
func (m *ActiveMock) GetEvictedPopulationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetEvictedPopulationCounter)
}

//GetEvictedPopulationMinimockPreCounter returns the value of ActiveMock.GetEvictedPopulation invocations
func (m *ActiveMock) GetEvictedPopulationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetEvictedPopulationPreCounter)
}

//GetEvictedPopulationFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetEvictedPopulationFinished() bool {
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

type mActiveMockGetExpectedPulseNumber struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetExpectedPulseNumberExpectation
	expectationSeries []*ActiveMockGetExpectedPulseNumberExpectation
}

type ActiveMockGetExpectedPulseNumberExpectation struct {
	result *ActiveMockGetExpectedPulseNumberResult
}

type ActiveMockGetExpectedPulseNumberResult struct {
	r pulse.Number
}

//Expect specifies that invocation of Active.GetExpectedPulseNumber is expected from 1 to Infinity times
func (m *mActiveMockGetExpectedPulseNumber) Expect() *mActiveMockGetExpectedPulseNumber {
	m.mock.GetExpectedPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetExpectedPulseNumberExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetExpectedPulseNumber
func (m *mActiveMockGetExpectedPulseNumber) Return(r pulse.Number) *ActiveMock {
	m.mock.GetExpectedPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetExpectedPulseNumberExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetExpectedPulseNumberResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetExpectedPulseNumber is expected once
func (m *mActiveMockGetExpectedPulseNumber) ExpectOnce() *ActiveMockGetExpectedPulseNumberExpectation {
	m.mock.GetExpectedPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetExpectedPulseNumberExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetExpectedPulseNumberExpectation) Return(r pulse.Number) {
	e.result = &ActiveMockGetExpectedPulseNumberResult{r}
}

//Set uses given function f as a mock of Active.GetExpectedPulseNumber method
func (m *mActiveMockGetExpectedPulseNumber) Set(f func() (r pulse.Number)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExpectedPulseNumberFunc = f
	return m.mock
}

//GetExpectedPulseNumber implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetExpectedPulseNumber() (r pulse.Number) {
	counter := atomic.AddUint64(&m.GetExpectedPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.GetExpectedPulseNumberCounter, 1)

	if len(m.GetExpectedPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExpectedPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetExpectedPulseNumber.")
			return
		}

		result := m.GetExpectedPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetExpectedPulseNumber")
			return
		}

		r = result.r

		return
	}

	if m.GetExpectedPulseNumberMock.mainExpectation != nil {

		result := m.GetExpectedPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetExpectedPulseNumber")
		}

		r = result.r

		return
	}

	if m.GetExpectedPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetExpectedPulseNumber.")
		return
	}

	return m.GetExpectedPulseNumberFunc()
}

//GetExpectedPulseNumberMinimockCounter returns a count of ActiveMock.GetExpectedPulseNumberFunc invocations
func (m *ActiveMock) GetExpectedPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExpectedPulseNumberCounter)
}

//GetExpectedPulseNumberMinimockPreCounter returns the value of ActiveMock.GetExpectedPulseNumber invocations
func (m *ActiveMock) GetExpectedPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExpectedPulseNumberPreCounter)
}

//GetExpectedPulseNumberFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetExpectedPulseNumberFinished() bool {
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

type mActiveMockGetGlobulaStateHash struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetGlobulaStateHashExpectation
	expectationSeries []*ActiveMockGetGlobulaStateHashExpectation
}

type ActiveMockGetGlobulaStateHashExpectation struct {
	result *ActiveMockGetGlobulaStateHashResult
}

type ActiveMockGetGlobulaStateHashResult struct {
	r proofs.GlobulaStateHash
}

//Expect specifies that invocation of Active.GetGlobulaStateHash is expected from 1 to Infinity times
func (m *mActiveMockGetGlobulaStateHash) Expect() *mActiveMockGetGlobulaStateHash {
	m.mock.GetGlobulaStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetGlobulaStateHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetGlobulaStateHash
func (m *mActiveMockGetGlobulaStateHash) Return(r proofs.GlobulaStateHash) *ActiveMock {
	m.mock.GetGlobulaStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetGlobulaStateHashExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetGlobulaStateHashResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetGlobulaStateHash is expected once
func (m *mActiveMockGetGlobulaStateHash) ExpectOnce() *ActiveMockGetGlobulaStateHashExpectation {
	m.mock.GetGlobulaStateHashFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetGlobulaStateHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetGlobulaStateHashExpectation) Return(r proofs.GlobulaStateHash) {
	e.result = &ActiveMockGetGlobulaStateHashResult{r}
}

//Set uses given function f as a mock of Active.GetGlobulaStateHash method
func (m *mActiveMockGetGlobulaStateHash) Set(f func() (r proofs.GlobulaStateHash)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetGlobulaStateHashFunc = f
	return m.mock
}

//GetGlobulaStateHash implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetGlobulaStateHash() (r proofs.GlobulaStateHash) {
	counter := atomic.AddUint64(&m.GetGlobulaStateHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobulaStateHashCounter, 1)

	if len(m.GetGlobulaStateHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetGlobulaStateHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetGlobulaStateHash.")
			return
		}

		result := m.GetGlobulaStateHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetGlobulaStateHash")
			return
		}

		r = result.r

		return
	}

	if m.GetGlobulaStateHashMock.mainExpectation != nil {

		result := m.GetGlobulaStateHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetGlobulaStateHash")
		}

		r = result.r

		return
	}

	if m.GetGlobulaStateHashFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetGlobulaStateHash.")
		return
	}

	return m.GetGlobulaStateHashFunc()
}

//GetGlobulaStateHashMinimockCounter returns a count of ActiveMock.GetGlobulaStateHashFunc invocations
func (m *ActiveMock) GetGlobulaStateHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobulaStateHashCounter)
}

//GetGlobulaStateHashMinimockPreCounter returns the value of ActiveMock.GetGlobulaStateHash invocations
func (m *ActiveMock) GetGlobulaStateHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobulaStateHashPreCounter)
}

//GetGlobulaStateHashFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetGlobulaStateHashFinished() bool {
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

type mActiveMockGetMandateRegistry struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetMandateRegistryExpectation
	expectationSeries []*ActiveMockGetMandateRegistryExpectation
}

type ActiveMockGetMandateRegistryExpectation struct {
	result *ActiveMockGetMandateRegistryResult
}

type ActiveMockGetMandateRegistryResult struct {
	r MandateRegistry
}

//Expect specifies that invocation of Active.GetMandateRegistry is expected from 1 to Infinity times
func (m *mActiveMockGetMandateRegistry) Expect() *mActiveMockGetMandateRegistry {
	m.mock.GetMandateRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetMandateRegistryExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetMandateRegistry
func (m *mActiveMockGetMandateRegistry) Return(r MandateRegistry) *ActiveMock {
	m.mock.GetMandateRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetMandateRegistryExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetMandateRegistryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetMandateRegistry is expected once
func (m *mActiveMockGetMandateRegistry) ExpectOnce() *ActiveMockGetMandateRegistryExpectation {
	m.mock.GetMandateRegistryFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetMandateRegistryExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetMandateRegistryExpectation) Return(r MandateRegistry) {
	e.result = &ActiveMockGetMandateRegistryResult{r}
}

//Set uses given function f as a mock of Active.GetMandateRegistry method
func (m *mActiveMockGetMandateRegistry) Set(f func() (r MandateRegistry)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMandateRegistryFunc = f
	return m.mock
}

//GetMandateRegistry implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetMandateRegistry() (r MandateRegistry) {
	counter := atomic.AddUint64(&m.GetMandateRegistryPreCounter, 1)
	defer atomic.AddUint64(&m.GetMandateRegistryCounter, 1)

	if len(m.GetMandateRegistryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMandateRegistryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetMandateRegistry.")
			return
		}

		result := m.GetMandateRegistryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetMandateRegistry")
			return
		}

		r = result.r

		return
	}

	if m.GetMandateRegistryMock.mainExpectation != nil {

		result := m.GetMandateRegistryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetMandateRegistry")
		}

		r = result.r

		return
	}

	if m.GetMandateRegistryFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetMandateRegistry.")
		return
	}

	return m.GetMandateRegistryFunc()
}

//GetMandateRegistryMinimockCounter returns a count of ActiveMock.GetMandateRegistryFunc invocations
func (m *ActiveMock) GetMandateRegistryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMandateRegistryCounter)
}

//GetMandateRegistryMinimockPreCounter returns the value of ActiveMock.GetMandateRegistry invocations
func (m *ActiveMock) GetMandateRegistryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMandateRegistryPreCounter)
}

//GetMandateRegistryFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetMandateRegistryFinished() bool {
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

type mActiveMockGetMisbehaviorRegistry struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetMisbehaviorRegistryExpectation
	expectationSeries []*ActiveMockGetMisbehaviorRegistryExpectation
}

type ActiveMockGetMisbehaviorRegistryExpectation struct {
	result *ActiveMockGetMisbehaviorRegistryResult
}

type ActiveMockGetMisbehaviorRegistryResult struct {
	r MisbehaviorRegistry
}

//Expect specifies that invocation of Active.GetMisbehaviorRegistry is expected from 1 to Infinity times
func (m *mActiveMockGetMisbehaviorRegistry) Expect() *mActiveMockGetMisbehaviorRegistry {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetMisbehaviorRegistryExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetMisbehaviorRegistry
func (m *mActiveMockGetMisbehaviorRegistry) Return(r MisbehaviorRegistry) *ActiveMock {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetMisbehaviorRegistryExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetMisbehaviorRegistryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetMisbehaviorRegistry is expected once
func (m *mActiveMockGetMisbehaviorRegistry) ExpectOnce() *ActiveMockGetMisbehaviorRegistryExpectation {
	m.mock.GetMisbehaviorRegistryFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetMisbehaviorRegistryExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetMisbehaviorRegistryExpectation) Return(r MisbehaviorRegistry) {
	e.result = &ActiveMockGetMisbehaviorRegistryResult{r}
}

//Set uses given function f as a mock of Active.GetMisbehaviorRegistry method
func (m *mActiveMockGetMisbehaviorRegistry) Set(f func() (r MisbehaviorRegistry)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMisbehaviorRegistryFunc = f
	return m.mock
}

//GetMisbehaviorRegistry implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetMisbehaviorRegistry() (r MisbehaviorRegistry) {
	counter := atomic.AddUint64(&m.GetMisbehaviorRegistryPreCounter, 1)
	defer atomic.AddUint64(&m.GetMisbehaviorRegistryCounter, 1)

	if len(m.GetMisbehaviorRegistryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMisbehaviorRegistryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetMisbehaviorRegistry.")
			return
		}

		result := m.GetMisbehaviorRegistryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetMisbehaviorRegistry")
			return
		}

		r = result.r

		return
	}

	if m.GetMisbehaviorRegistryMock.mainExpectation != nil {

		result := m.GetMisbehaviorRegistryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetMisbehaviorRegistry")
		}

		r = result.r

		return
	}

	if m.GetMisbehaviorRegistryFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetMisbehaviorRegistry.")
		return
	}

	return m.GetMisbehaviorRegistryFunc()
}

//GetMisbehaviorRegistryMinimockCounter returns a count of ActiveMock.GetMisbehaviorRegistryFunc invocations
func (m *ActiveMock) GetMisbehaviorRegistryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMisbehaviorRegistryCounter)
}

//GetMisbehaviorRegistryMinimockPreCounter returns the value of ActiveMock.GetMisbehaviorRegistry invocations
func (m *ActiveMock) GetMisbehaviorRegistryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMisbehaviorRegistryPreCounter)
}

//GetMisbehaviorRegistryFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetMisbehaviorRegistryFinished() bool {
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

type mActiveMockGetOfflinePopulation struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetOfflinePopulationExpectation
	expectationSeries []*ActiveMockGetOfflinePopulationExpectation
}

type ActiveMockGetOfflinePopulationExpectation struct {
	result *ActiveMockGetOfflinePopulationResult
}

type ActiveMockGetOfflinePopulationResult struct {
	r OfflinePopulation
}

//Expect specifies that invocation of Active.GetOfflinePopulation is expected from 1 to Infinity times
func (m *mActiveMockGetOfflinePopulation) Expect() *mActiveMockGetOfflinePopulation {
	m.mock.GetOfflinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetOfflinePopulationExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetOfflinePopulation
func (m *mActiveMockGetOfflinePopulation) Return(r OfflinePopulation) *ActiveMock {
	m.mock.GetOfflinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetOfflinePopulationExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetOfflinePopulationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetOfflinePopulation is expected once
func (m *mActiveMockGetOfflinePopulation) ExpectOnce() *ActiveMockGetOfflinePopulationExpectation {
	m.mock.GetOfflinePopulationFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetOfflinePopulationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetOfflinePopulationExpectation) Return(r OfflinePopulation) {
	e.result = &ActiveMockGetOfflinePopulationResult{r}
}

//Set uses given function f as a mock of Active.GetOfflinePopulation method
func (m *mActiveMockGetOfflinePopulation) Set(f func() (r OfflinePopulation)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOfflinePopulationFunc = f
	return m.mock
}

//GetOfflinePopulation implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetOfflinePopulation() (r OfflinePopulation) {
	counter := atomic.AddUint64(&m.GetOfflinePopulationPreCounter, 1)
	defer atomic.AddUint64(&m.GetOfflinePopulationCounter, 1)

	if len(m.GetOfflinePopulationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOfflinePopulationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetOfflinePopulation.")
			return
		}

		result := m.GetOfflinePopulationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetOfflinePopulation")
			return
		}

		r = result.r

		return
	}

	if m.GetOfflinePopulationMock.mainExpectation != nil {

		result := m.GetOfflinePopulationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetOfflinePopulation")
		}

		r = result.r

		return
	}

	if m.GetOfflinePopulationFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetOfflinePopulation.")
		return
	}

	return m.GetOfflinePopulationFunc()
}

//GetOfflinePopulationMinimockCounter returns a count of ActiveMock.GetOfflinePopulationFunc invocations
func (m *ActiveMock) GetOfflinePopulationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOfflinePopulationCounter)
}

//GetOfflinePopulationMinimockPreCounter returns the value of ActiveMock.GetOfflinePopulation invocations
func (m *ActiveMock) GetOfflinePopulationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOfflinePopulationPreCounter)
}

//GetOfflinePopulationFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetOfflinePopulationFinished() bool {
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

type mActiveMockGetOnlinePopulation struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetOnlinePopulationExpectation
	expectationSeries []*ActiveMockGetOnlinePopulationExpectation
}

type ActiveMockGetOnlinePopulationExpectation struct {
	result *ActiveMockGetOnlinePopulationResult
}

type ActiveMockGetOnlinePopulationResult struct {
	r OnlinePopulation
}

//Expect specifies that invocation of Active.GetOnlinePopulation is expected from 1 to Infinity times
func (m *mActiveMockGetOnlinePopulation) Expect() *mActiveMockGetOnlinePopulation {
	m.mock.GetOnlinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetOnlinePopulationExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetOnlinePopulation
func (m *mActiveMockGetOnlinePopulation) Return(r OnlinePopulation) *ActiveMock {
	m.mock.GetOnlinePopulationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetOnlinePopulationExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetOnlinePopulationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetOnlinePopulation is expected once
func (m *mActiveMockGetOnlinePopulation) ExpectOnce() *ActiveMockGetOnlinePopulationExpectation {
	m.mock.GetOnlinePopulationFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetOnlinePopulationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetOnlinePopulationExpectation) Return(r OnlinePopulation) {
	e.result = &ActiveMockGetOnlinePopulationResult{r}
}

//Set uses given function f as a mock of Active.GetOnlinePopulation method
func (m *mActiveMockGetOnlinePopulation) Set(f func() (r OnlinePopulation)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOnlinePopulationFunc = f
	return m.mock
}

//GetOnlinePopulation implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetOnlinePopulation() (r OnlinePopulation) {
	counter := atomic.AddUint64(&m.GetOnlinePopulationPreCounter, 1)
	defer atomic.AddUint64(&m.GetOnlinePopulationCounter, 1)

	if len(m.GetOnlinePopulationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOnlinePopulationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetOnlinePopulation.")
			return
		}

		result := m.GetOnlinePopulationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetOnlinePopulation")
			return
		}

		r = result.r

		return
	}

	if m.GetOnlinePopulationMock.mainExpectation != nil {

		result := m.GetOnlinePopulationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetOnlinePopulation")
		}

		r = result.r

		return
	}

	if m.GetOnlinePopulationFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetOnlinePopulation.")
		return
	}

	return m.GetOnlinePopulationFunc()
}

//GetOnlinePopulationMinimockCounter returns a count of ActiveMock.GetOnlinePopulationFunc invocations
func (m *ActiveMock) GetOnlinePopulationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOnlinePopulationCounter)
}

//GetOnlinePopulationMinimockPreCounter returns the value of ActiveMock.GetOnlinePopulation invocations
func (m *ActiveMock) GetOnlinePopulationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOnlinePopulationPreCounter)
}

//GetOnlinePopulationFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetOnlinePopulationFinished() bool {
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

type mActiveMockGetProfileFactory struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetProfileFactoryExpectation
	expectationSeries []*ActiveMockGetProfileFactoryExpectation
}

type ActiveMockGetProfileFactoryExpectation struct {
	input  *ActiveMockGetProfileFactoryInput
	result *ActiveMockGetProfileFactoryResult
}

type ActiveMockGetProfileFactoryInput struct {
	p cryptkit.KeyStoreFactory
}

type ActiveMockGetProfileFactoryResult struct {
	r profiles.Factory
}

//Expect specifies that invocation of Active.GetProfileFactory is expected from 1 to Infinity times
func (m *mActiveMockGetProfileFactory) Expect(p cryptkit.KeyStoreFactory) *mActiveMockGetProfileFactory {
	m.mock.GetProfileFactoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetProfileFactoryExpectation{}
	}
	m.mainExpectation.input = &ActiveMockGetProfileFactoryInput{p}
	return m
}

//Return specifies results of invocation of Active.GetProfileFactory
func (m *mActiveMockGetProfileFactory) Return(r profiles.Factory) *ActiveMock {
	m.mock.GetProfileFactoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetProfileFactoryExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetProfileFactoryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetProfileFactory is expected once
func (m *mActiveMockGetProfileFactory) ExpectOnce(p cryptkit.KeyStoreFactory) *ActiveMockGetProfileFactoryExpectation {
	m.mock.GetProfileFactoryFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetProfileFactoryExpectation{}
	expectation.input = &ActiveMockGetProfileFactoryInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetProfileFactoryExpectation) Return(r profiles.Factory) {
	e.result = &ActiveMockGetProfileFactoryResult{r}
}

//Set uses given function f as a mock of Active.GetProfileFactory method
func (m *mActiveMockGetProfileFactory) Set(f func(p cryptkit.KeyStoreFactory) (r profiles.Factory)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetProfileFactoryFunc = f
	return m.mock
}

//GetProfileFactory implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetProfileFactory(p cryptkit.KeyStoreFactory) (r profiles.Factory) {
	counter := atomic.AddUint64(&m.GetProfileFactoryPreCounter, 1)
	defer atomic.AddUint64(&m.GetProfileFactoryCounter, 1)

	if len(m.GetProfileFactoryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetProfileFactoryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetProfileFactory. %v", p)
			return
		}

		input := m.GetProfileFactoryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ActiveMockGetProfileFactoryInput{p}, "Active.GetProfileFactory got unexpected parameters")

		result := m.GetProfileFactoryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetProfileFactory")
			return
		}

		r = result.r

		return
	}

	if m.GetProfileFactoryMock.mainExpectation != nil {

		input := m.GetProfileFactoryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ActiveMockGetProfileFactoryInput{p}, "Active.GetProfileFactory got unexpected parameters")
		}

		result := m.GetProfileFactoryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetProfileFactory")
		}

		r = result.r

		return
	}

	if m.GetProfileFactoryFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetProfileFactory. %v", p)
		return
	}

	return m.GetProfileFactoryFunc(p)
}

//GetProfileFactoryMinimockCounter returns a count of ActiveMock.GetProfileFactoryFunc invocations
func (m *ActiveMock) GetProfileFactoryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetProfileFactoryCounter)
}

//GetProfileFactoryMinimockPreCounter returns the value of ActiveMock.GetProfileFactory invocations
func (m *ActiveMock) GetProfileFactoryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetProfileFactoryPreCounter)
}

//GetProfileFactoryFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetProfileFactoryFinished() bool {
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

type mActiveMockGetPulseData struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetPulseDataExpectation
	expectationSeries []*ActiveMockGetPulseDataExpectation
}

type ActiveMockGetPulseDataExpectation struct {
	result *ActiveMockGetPulseDataResult
}

type ActiveMockGetPulseDataResult struct {
	r pulse.Data
}

//Expect specifies that invocation of Active.GetPulseData is expected from 1 to Infinity times
func (m *mActiveMockGetPulseData) Expect() *mActiveMockGetPulseData {
	m.mock.GetPulseDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetPulseDataExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetPulseData
func (m *mActiveMockGetPulseData) Return(r pulse.Data) *ActiveMock {
	m.mock.GetPulseDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetPulseDataExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetPulseDataResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetPulseData is expected once
func (m *mActiveMockGetPulseData) ExpectOnce() *ActiveMockGetPulseDataExpectation {
	m.mock.GetPulseDataFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetPulseDataExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetPulseDataExpectation) Return(r pulse.Data) {
	e.result = &ActiveMockGetPulseDataResult{r}
}

//Set uses given function f as a mock of Active.GetPulseData method
func (m *mActiveMockGetPulseData) Set(f func() (r pulse.Data)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseDataFunc = f
	return m.mock
}

//GetPulseData implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetPulseData() (r pulse.Data) {
	counter := atomic.AddUint64(&m.GetPulseDataPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseDataCounter, 1)

	if len(m.GetPulseDataMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseDataMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetPulseData.")
			return
		}

		result := m.GetPulseDataMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetPulseData")
			return
		}

		r = result.r

		return
	}

	if m.GetPulseDataMock.mainExpectation != nil {

		result := m.GetPulseDataMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetPulseData")
		}

		r = result.r

		return
	}

	if m.GetPulseDataFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetPulseData.")
		return
	}

	return m.GetPulseDataFunc()
}

//GetPulseDataMinimockCounter returns a count of ActiveMock.GetPulseDataFunc invocations
func (m *ActiveMock) GetPulseDataMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseDataCounter)
}

//GetPulseDataMinimockPreCounter returns the value of ActiveMock.GetPulseData invocations
func (m *ActiveMock) GetPulseDataMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseDataPreCounter)
}

//GetPulseDataFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetPulseDataFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulseDataMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulseDataCounter) == uint64(len(m.GetPulseDataMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulseDataMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulseDataCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulseDataFunc != nil {
		return atomic.LoadUint64(&m.GetPulseDataCounter) > 0
	}

	return true
}

type mActiveMockGetPulseNumber struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockGetPulseNumberExpectation
	expectationSeries []*ActiveMockGetPulseNumberExpectation
}

type ActiveMockGetPulseNumberExpectation struct {
	result *ActiveMockGetPulseNumberResult
}

type ActiveMockGetPulseNumberResult struct {
	r pulse.Number
}

//Expect specifies that invocation of Active.GetPulseNumber is expected from 1 to Infinity times
func (m *mActiveMockGetPulseNumber) Expect() *mActiveMockGetPulseNumber {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetPulseNumberExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.GetPulseNumber
func (m *mActiveMockGetPulseNumber) Return(r pulse.Number) *ActiveMock {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockGetPulseNumberExpectation{}
	}
	m.mainExpectation.result = &ActiveMockGetPulseNumberResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.GetPulseNumber is expected once
func (m *mActiveMockGetPulseNumber) ExpectOnce() *ActiveMockGetPulseNumberExpectation {
	m.mock.GetPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockGetPulseNumberExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockGetPulseNumberExpectation) Return(r pulse.Number) {
	e.result = &ActiveMockGetPulseNumberResult{r}
}

//Set uses given function f as a mock of Active.GetPulseNumber method
func (m *mActiveMockGetPulseNumber) Set(f func() (r pulse.Number)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseNumberFunc = f
	return m.mock
}

//GetPulseNumber implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) GetPulseNumber() (r pulse.Number) {
	counter := atomic.AddUint64(&m.GetPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseNumberCounter, 1)

	if len(m.GetPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.GetPulseNumber.")
			return
		}

		result := m.GetPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetPulseNumber")
			return
		}

		r = result.r

		return
	}

	if m.GetPulseNumberMock.mainExpectation != nil {

		result := m.GetPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.GetPulseNumber")
		}

		r = result.r

		return
	}

	if m.GetPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.GetPulseNumber.")
		return
	}

	return m.GetPulseNumberFunc()
}

//GetPulseNumberMinimockCounter returns a count of ActiveMock.GetPulseNumberFunc invocations
func (m *ActiveMock) GetPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberCounter)
}

//GetPulseNumberMinimockPreCounter returns the value of ActiveMock.GetPulseNumber invocations
func (m *ActiveMock) GetPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberPreCounter)
}

//GetPulseNumberFinished returns true if mock invocations count is ok
func (m *ActiveMock) GetPulseNumberFinished() bool {
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

type mActiveMockIsActive struct {
	mock              *ActiveMock
	mainExpectation   *ActiveMockIsActiveExpectation
	expectationSeries []*ActiveMockIsActiveExpectation
}

type ActiveMockIsActiveExpectation struct {
	result *ActiveMockIsActiveResult
}

type ActiveMockIsActiveResult struct {
	r bool
}

//Expect specifies that invocation of Active.IsActive is expected from 1 to Infinity times
func (m *mActiveMockIsActive) Expect() *mActiveMockIsActive {
	m.mock.IsActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockIsActiveExpectation{}
	}

	return m
}

//Return specifies results of invocation of Active.IsActive
func (m *mActiveMockIsActive) Return(r bool) *ActiveMock {
	m.mock.IsActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveMockIsActiveExpectation{}
	}
	m.mainExpectation.result = &ActiveMockIsActiveResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Active.IsActive is expected once
func (m *mActiveMockIsActive) ExpectOnce() *ActiveMockIsActiveExpectation {
	m.mock.IsActiveFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveMockIsActiveExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveMockIsActiveExpectation) Return(r bool) {
	e.result = &ActiveMockIsActiveResult{r}
}

//Set uses given function f as a mock of Active.IsActive method
func (m *mActiveMockIsActive) Set(f func() (r bool)) *ActiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsActiveFunc = f
	return m.mock
}

//IsActive implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active interface
func (m *ActiveMock) IsActive() (r bool) {
	counter := atomic.AddUint64(&m.IsActivePreCounter, 1)
	defer atomic.AddUint64(&m.IsActiveCounter, 1)

	if len(m.IsActiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsActiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveMock.IsActive.")
			return
		}

		result := m.IsActiveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.IsActive")
			return
		}

		r = result.r

		return
	}

	if m.IsActiveMock.mainExpectation != nil {

		result := m.IsActiveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveMock.IsActive")
		}

		r = result.r

		return
	}

	if m.IsActiveFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveMock.IsActive.")
		return
	}

	return m.IsActiveFunc()
}

//IsActiveMinimockCounter returns a count of ActiveMock.IsActiveFunc invocations
func (m *ActiveMock) IsActiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsActiveCounter)
}

//IsActiveMinimockPreCounter returns the value of ActiveMock.IsActive invocations
func (m *ActiveMock) IsActiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsActivePreCounter)
}

//IsActiveFinished returns true if mock invocations count is ok
func (m *ActiveMock) IsActiveFinished() bool {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ActiveMock) ValidateCallCounters() {

	if !m.CreateBuilderFinished() {
		m.t.Fatal("Expected call to ActiveMock.CreateBuilder")
	}

	if !m.GetCensusStateFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetCensusState")
	}

	if !m.GetCloudStateHashFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetCloudStateHash")
	}

	if !m.GetEvictedPopulationFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetEvictedPopulation")
	}

	if !m.GetExpectedPulseNumberFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetExpectedPulseNumber")
	}

	if !m.GetGlobulaStateHashFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetGlobulaStateHash")
	}

	if !m.GetMandateRegistryFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetMandateRegistry")
	}

	if !m.GetMisbehaviorRegistryFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetMisbehaviorRegistry")
	}

	if !m.GetOfflinePopulationFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetOfflinePopulation")
	}

	if !m.GetOnlinePopulationFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetOnlinePopulation")
	}

	if !m.GetProfileFactoryFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetProfileFactory")
	}

	if !m.GetPulseDataFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetPulseData")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetPulseNumber")
	}

	if !m.IsActiveFinished() {
		m.t.Fatal("Expected call to ActiveMock.IsActive")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ActiveMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ActiveMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ActiveMock) MinimockFinish() {

	if !m.CreateBuilderFinished() {
		m.t.Fatal("Expected call to ActiveMock.CreateBuilder")
	}

	if !m.GetCensusStateFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetCensusState")
	}

	if !m.GetCloudStateHashFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetCloudStateHash")
	}

	if !m.GetEvictedPopulationFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetEvictedPopulation")
	}

	if !m.GetExpectedPulseNumberFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetExpectedPulseNumber")
	}

	if !m.GetGlobulaStateHashFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetGlobulaStateHash")
	}

	if !m.GetMandateRegistryFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetMandateRegistry")
	}

	if !m.GetMisbehaviorRegistryFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetMisbehaviorRegistry")
	}

	if !m.GetOfflinePopulationFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetOfflinePopulation")
	}

	if !m.GetOnlinePopulationFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetOnlinePopulation")
	}

	if !m.GetProfileFactoryFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetProfileFactory")
	}

	if !m.GetPulseDataFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetPulseData")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to ActiveMock.GetPulseNumber")
	}

	if !m.IsActiveFinished() {
		m.t.Fatal("Expected call to ActiveMock.IsActive")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ActiveMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ActiveMock) MinimockWait(timeout time.Duration) {
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
		ok = ok && m.GetProfileFactoryFinished()
		ok = ok && m.GetPulseDataFinished()
		ok = ok && m.GetPulseNumberFinished()
		ok = ok && m.IsActiveFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CreateBuilderFinished() {
				m.t.Error("Expected call to ActiveMock.CreateBuilder")
			}

			if !m.GetCensusStateFinished() {
				m.t.Error("Expected call to ActiveMock.GetCensusState")
			}

			if !m.GetCloudStateHashFinished() {
				m.t.Error("Expected call to ActiveMock.GetCloudStateHash")
			}

			if !m.GetEvictedPopulationFinished() {
				m.t.Error("Expected call to ActiveMock.GetEvictedPopulation")
			}

			if !m.GetExpectedPulseNumberFinished() {
				m.t.Error("Expected call to ActiveMock.GetExpectedPulseNumber")
			}

			if !m.GetGlobulaStateHashFinished() {
				m.t.Error("Expected call to ActiveMock.GetGlobulaStateHash")
			}

			if !m.GetMandateRegistryFinished() {
				m.t.Error("Expected call to ActiveMock.GetMandateRegistry")
			}

			if !m.GetMisbehaviorRegistryFinished() {
				m.t.Error("Expected call to ActiveMock.GetMisbehaviorRegistry")
			}

			if !m.GetOfflinePopulationFinished() {
				m.t.Error("Expected call to ActiveMock.GetOfflinePopulation")
			}

			if !m.GetOnlinePopulationFinished() {
				m.t.Error("Expected call to ActiveMock.GetOnlinePopulation")
			}

			if !m.GetProfileFactoryFinished() {
				m.t.Error("Expected call to ActiveMock.GetProfileFactory")
			}

			if !m.GetPulseDataFinished() {
				m.t.Error("Expected call to ActiveMock.GetPulseData")
			}

			if !m.GetPulseNumberFinished() {
				m.t.Error("Expected call to ActiveMock.GetPulseNumber")
			}

			if !m.IsActiveFinished() {
				m.t.Error("Expected call to ActiveMock.IsActive")
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
func (m *ActiveMock) AllMocksCalled() bool {

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

	if !m.GetProfileFactoryFinished() {
		return false
	}

	if !m.GetPulseDataFinished() {
		return false
	}

	if !m.GetPulseNumberFinished() {
		return false
	}

	if !m.IsActiveFinished() {
		return false
	}

	return true
}

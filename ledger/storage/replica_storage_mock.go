package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ReplicaStorage" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ReplicaStorageMock implements github.com/insolar/insolar/ledger/storage.ReplicaStorage
type ReplicaStorageMock struct {
	t minimock.Tester

	GetAllNonEmptySyncClientJetsFunc       func(p context.Context) (r map[core.RecordID][]core.PulseNumber, r1 error)
	GetAllNonEmptySyncClientJetsCounter    uint64
	GetAllNonEmptySyncClientJetsPreCounter uint64
	GetAllNonEmptySyncClientJetsMock       mReplicaStorageMockGetAllNonEmptySyncClientJets

	GetAllSyncClientJetsFunc       func(p context.Context) (r map[core.RecordID][]core.PulseNumber, r1 error)
	GetAllSyncClientJetsCounter    uint64
	GetAllSyncClientJetsPreCounter uint64
	GetAllSyncClientJetsMock       mReplicaStorageMockGetAllSyncClientJets

	GetHeavySyncedPulseFunc       func(p context.Context, p1 core.RecordID) (r core.PulseNumber, r1 error)
	GetHeavySyncedPulseCounter    uint64
	GetHeavySyncedPulsePreCounter uint64
	GetHeavySyncedPulseMock       mReplicaStorageMockGetHeavySyncedPulse

	GetSyncClientJetPulsesFunc       func(p context.Context, p1 core.RecordID) (r []core.PulseNumber, r1 error)
	GetSyncClientJetPulsesCounter    uint64
	GetSyncClientJetPulsesPreCounter uint64
	GetSyncClientJetPulsesMock       mReplicaStorageMockGetSyncClientJetPulses

	SetHeavySyncedPulseFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)
	SetHeavySyncedPulseCounter    uint64
	SetHeavySyncedPulsePreCounter uint64
	SetHeavySyncedPulseMock       mReplicaStorageMockSetHeavySyncedPulse

	SetSyncClientJetPulsesFunc       func(p context.Context, p1 core.RecordID, p2 []core.PulseNumber) (r error)
	SetSyncClientJetPulsesCounter    uint64
	SetSyncClientJetPulsesPreCounter uint64
	SetSyncClientJetPulsesMock       mReplicaStorageMockSetSyncClientJetPulses
}

//NewReplicaStorageMock returns a mock for github.com/insolar/insolar/ledger/storage.ReplicaStorage
func NewReplicaStorageMock(t minimock.Tester) *ReplicaStorageMock {
	m := &ReplicaStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAllNonEmptySyncClientJetsMock = mReplicaStorageMockGetAllNonEmptySyncClientJets{mock: m}
	m.GetAllSyncClientJetsMock = mReplicaStorageMockGetAllSyncClientJets{mock: m}
	m.GetHeavySyncedPulseMock = mReplicaStorageMockGetHeavySyncedPulse{mock: m}
	m.GetSyncClientJetPulsesMock = mReplicaStorageMockGetSyncClientJetPulses{mock: m}
	m.SetHeavySyncedPulseMock = mReplicaStorageMockSetHeavySyncedPulse{mock: m}
	m.SetSyncClientJetPulsesMock = mReplicaStorageMockSetSyncClientJetPulses{mock: m}

	return m
}

type mReplicaStorageMockGetAllNonEmptySyncClientJets struct {
	mock              *ReplicaStorageMock
	mainExpectation   *ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation
	expectationSeries []*ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation
}

type ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation struct {
	input  *ReplicaStorageMockGetAllNonEmptySyncClientJetsInput
	result *ReplicaStorageMockGetAllNonEmptySyncClientJetsResult
}

type ReplicaStorageMockGetAllNonEmptySyncClientJetsInput struct {
	p context.Context
}

type ReplicaStorageMockGetAllNonEmptySyncClientJetsResult struct {
	r  map[core.RecordID][]core.PulseNumber
	r1 error
}

//Expect specifies that invocation of ReplicaStorage.GetAllNonEmptySyncClientJets is expected from 1 to Infinity times
func (m *mReplicaStorageMockGetAllNonEmptySyncClientJets) Expect(p context.Context) *mReplicaStorageMockGetAllNonEmptySyncClientJets {
	m.mock.GetAllNonEmptySyncClientJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation{}
	}
	m.mainExpectation.input = &ReplicaStorageMockGetAllNonEmptySyncClientJetsInput{p}
	return m
}

//Return specifies results of invocation of ReplicaStorage.GetAllNonEmptySyncClientJets
func (m *mReplicaStorageMockGetAllNonEmptySyncClientJets) Return(r map[core.RecordID][]core.PulseNumber, r1 error) *ReplicaStorageMock {
	m.mock.GetAllNonEmptySyncClientJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation{}
	}
	m.mainExpectation.result = &ReplicaStorageMockGetAllNonEmptySyncClientJetsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ReplicaStorage.GetAllNonEmptySyncClientJets is expected once
func (m *mReplicaStorageMockGetAllNonEmptySyncClientJets) ExpectOnce(p context.Context) *ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation {
	m.mock.GetAllNonEmptySyncClientJetsFunc = nil
	m.mainExpectation = nil

	expectation := &ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation{}
	expectation.input = &ReplicaStorageMockGetAllNonEmptySyncClientJetsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReplicaStorageMockGetAllNonEmptySyncClientJetsExpectation) Return(r map[core.RecordID][]core.PulseNumber, r1 error) {
	e.result = &ReplicaStorageMockGetAllNonEmptySyncClientJetsResult{r, r1}
}

//Set uses given function f as a mock of ReplicaStorage.GetAllNonEmptySyncClientJets method
func (m *mReplicaStorageMockGetAllNonEmptySyncClientJets) Set(f func(p context.Context) (r map[core.RecordID][]core.PulseNumber, r1 error)) *ReplicaStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAllNonEmptySyncClientJetsFunc = f
	return m.mock
}

//GetAllNonEmptySyncClientJets implements github.com/insolar/insolar/ledger/storage.ReplicaStorage interface
func (m *ReplicaStorageMock) GetAllNonEmptySyncClientJets(p context.Context) (r map[core.RecordID][]core.PulseNumber, r1 error) {
	counter := atomic.AddUint64(&m.GetAllNonEmptySyncClientJetsPreCounter, 1)
	defer atomic.AddUint64(&m.GetAllNonEmptySyncClientJetsCounter, 1)

	if len(m.GetAllNonEmptySyncClientJetsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAllNonEmptySyncClientJetsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetAllNonEmptySyncClientJets. %v", p)
			return
		}

		input := m.GetAllNonEmptySyncClientJetsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ReplicaStorageMockGetAllNonEmptySyncClientJetsInput{p}, "ReplicaStorage.GetAllNonEmptySyncClientJets got unexpected parameters")

		result := m.GetAllNonEmptySyncClientJetsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetAllNonEmptySyncClientJets")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetAllNonEmptySyncClientJetsMock.mainExpectation != nil {

		input := m.GetAllNonEmptySyncClientJetsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ReplicaStorageMockGetAllNonEmptySyncClientJetsInput{p}, "ReplicaStorage.GetAllNonEmptySyncClientJets got unexpected parameters")
		}

		result := m.GetAllNonEmptySyncClientJetsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetAllNonEmptySyncClientJets")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetAllNonEmptySyncClientJetsFunc == nil {
		m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetAllNonEmptySyncClientJets. %v", p)
		return
	}

	return m.GetAllNonEmptySyncClientJetsFunc(p)
}

//GetAllNonEmptySyncClientJetsMinimockCounter returns a count of ReplicaStorageMock.GetAllNonEmptySyncClientJetsFunc invocations
func (m *ReplicaStorageMock) GetAllNonEmptySyncClientJetsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAllNonEmptySyncClientJetsCounter)
}

//GetAllNonEmptySyncClientJetsMinimockPreCounter returns the value of ReplicaStorageMock.GetAllNonEmptySyncClientJets invocations
func (m *ReplicaStorageMock) GetAllNonEmptySyncClientJetsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAllNonEmptySyncClientJetsPreCounter)
}

//GetAllNonEmptySyncClientJetsFinished returns true if mock invocations count is ok
func (m *ReplicaStorageMock) GetAllNonEmptySyncClientJetsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetAllNonEmptySyncClientJetsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetAllNonEmptySyncClientJetsCounter) == uint64(len(m.GetAllNonEmptySyncClientJetsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetAllNonEmptySyncClientJetsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetAllNonEmptySyncClientJetsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetAllNonEmptySyncClientJetsFunc != nil {
		return atomic.LoadUint64(&m.GetAllNonEmptySyncClientJetsCounter) > 0
	}

	return true
}

type mReplicaStorageMockGetAllSyncClientJets struct {
	mock              *ReplicaStorageMock
	mainExpectation   *ReplicaStorageMockGetAllSyncClientJetsExpectation
	expectationSeries []*ReplicaStorageMockGetAllSyncClientJetsExpectation
}

type ReplicaStorageMockGetAllSyncClientJetsExpectation struct {
	input  *ReplicaStorageMockGetAllSyncClientJetsInput
	result *ReplicaStorageMockGetAllSyncClientJetsResult
}

type ReplicaStorageMockGetAllSyncClientJetsInput struct {
	p context.Context
}

type ReplicaStorageMockGetAllSyncClientJetsResult struct {
	r  map[core.RecordID][]core.PulseNumber
	r1 error
}

//Expect specifies that invocation of ReplicaStorage.GetAllSyncClientJets is expected from 1 to Infinity times
func (m *mReplicaStorageMockGetAllSyncClientJets) Expect(p context.Context) *mReplicaStorageMockGetAllSyncClientJets {
	m.mock.GetAllSyncClientJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetAllSyncClientJetsExpectation{}
	}
	m.mainExpectation.input = &ReplicaStorageMockGetAllSyncClientJetsInput{p}
	return m
}

//Return specifies results of invocation of ReplicaStorage.GetAllSyncClientJets
func (m *mReplicaStorageMockGetAllSyncClientJets) Return(r map[core.RecordID][]core.PulseNumber, r1 error) *ReplicaStorageMock {
	m.mock.GetAllSyncClientJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetAllSyncClientJetsExpectation{}
	}
	m.mainExpectation.result = &ReplicaStorageMockGetAllSyncClientJetsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ReplicaStorage.GetAllSyncClientJets is expected once
func (m *mReplicaStorageMockGetAllSyncClientJets) ExpectOnce(p context.Context) *ReplicaStorageMockGetAllSyncClientJetsExpectation {
	m.mock.GetAllSyncClientJetsFunc = nil
	m.mainExpectation = nil

	expectation := &ReplicaStorageMockGetAllSyncClientJetsExpectation{}
	expectation.input = &ReplicaStorageMockGetAllSyncClientJetsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReplicaStorageMockGetAllSyncClientJetsExpectation) Return(r map[core.RecordID][]core.PulseNumber, r1 error) {
	e.result = &ReplicaStorageMockGetAllSyncClientJetsResult{r, r1}
}

//Set uses given function f as a mock of ReplicaStorage.GetAllSyncClientJets method
func (m *mReplicaStorageMockGetAllSyncClientJets) Set(f func(p context.Context) (r map[core.RecordID][]core.PulseNumber, r1 error)) *ReplicaStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAllSyncClientJetsFunc = f
	return m.mock
}

//GetAllSyncClientJets implements github.com/insolar/insolar/ledger/storage.ReplicaStorage interface
func (m *ReplicaStorageMock) GetAllSyncClientJets(p context.Context) (r map[core.RecordID][]core.PulseNumber, r1 error) {
	counter := atomic.AddUint64(&m.GetAllSyncClientJetsPreCounter, 1)
	defer atomic.AddUint64(&m.GetAllSyncClientJetsCounter, 1)

	if len(m.GetAllSyncClientJetsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAllSyncClientJetsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetAllSyncClientJets. %v", p)
			return
		}

		input := m.GetAllSyncClientJetsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ReplicaStorageMockGetAllSyncClientJetsInput{p}, "ReplicaStorage.GetAllSyncClientJets got unexpected parameters")

		result := m.GetAllSyncClientJetsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetAllSyncClientJets")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetAllSyncClientJetsMock.mainExpectation != nil {

		input := m.GetAllSyncClientJetsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ReplicaStorageMockGetAllSyncClientJetsInput{p}, "ReplicaStorage.GetAllSyncClientJets got unexpected parameters")
		}

		result := m.GetAllSyncClientJetsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetAllSyncClientJets")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetAllSyncClientJetsFunc == nil {
		m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetAllSyncClientJets. %v", p)
		return
	}

	return m.GetAllSyncClientJetsFunc(p)
}

//GetAllSyncClientJetsMinimockCounter returns a count of ReplicaStorageMock.GetAllSyncClientJetsFunc invocations
func (m *ReplicaStorageMock) GetAllSyncClientJetsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAllSyncClientJetsCounter)
}

//GetAllSyncClientJetsMinimockPreCounter returns the value of ReplicaStorageMock.GetAllSyncClientJets invocations
func (m *ReplicaStorageMock) GetAllSyncClientJetsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAllSyncClientJetsPreCounter)
}

//GetAllSyncClientJetsFinished returns true if mock invocations count is ok
func (m *ReplicaStorageMock) GetAllSyncClientJetsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetAllSyncClientJetsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetAllSyncClientJetsCounter) == uint64(len(m.GetAllSyncClientJetsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetAllSyncClientJetsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetAllSyncClientJetsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetAllSyncClientJetsFunc != nil {
		return atomic.LoadUint64(&m.GetAllSyncClientJetsCounter) > 0
	}

	return true
}

type mReplicaStorageMockGetHeavySyncedPulse struct {
	mock              *ReplicaStorageMock
	mainExpectation   *ReplicaStorageMockGetHeavySyncedPulseExpectation
	expectationSeries []*ReplicaStorageMockGetHeavySyncedPulseExpectation
}

type ReplicaStorageMockGetHeavySyncedPulseExpectation struct {
	input  *ReplicaStorageMockGetHeavySyncedPulseInput
	result *ReplicaStorageMockGetHeavySyncedPulseResult
}

type ReplicaStorageMockGetHeavySyncedPulseInput struct {
	p  context.Context
	p1 core.RecordID
}

type ReplicaStorageMockGetHeavySyncedPulseResult struct {
	r  core.PulseNumber
	r1 error
}

//Expect specifies that invocation of ReplicaStorage.GetHeavySyncedPulse is expected from 1 to Infinity times
func (m *mReplicaStorageMockGetHeavySyncedPulse) Expect(p context.Context, p1 core.RecordID) *mReplicaStorageMockGetHeavySyncedPulse {
	m.mock.GetHeavySyncedPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetHeavySyncedPulseExpectation{}
	}
	m.mainExpectation.input = &ReplicaStorageMockGetHeavySyncedPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of ReplicaStorage.GetHeavySyncedPulse
func (m *mReplicaStorageMockGetHeavySyncedPulse) Return(r core.PulseNumber, r1 error) *ReplicaStorageMock {
	m.mock.GetHeavySyncedPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetHeavySyncedPulseExpectation{}
	}
	m.mainExpectation.result = &ReplicaStorageMockGetHeavySyncedPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ReplicaStorage.GetHeavySyncedPulse is expected once
func (m *mReplicaStorageMockGetHeavySyncedPulse) ExpectOnce(p context.Context, p1 core.RecordID) *ReplicaStorageMockGetHeavySyncedPulseExpectation {
	m.mock.GetHeavySyncedPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ReplicaStorageMockGetHeavySyncedPulseExpectation{}
	expectation.input = &ReplicaStorageMockGetHeavySyncedPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReplicaStorageMockGetHeavySyncedPulseExpectation) Return(r core.PulseNumber, r1 error) {
	e.result = &ReplicaStorageMockGetHeavySyncedPulseResult{r, r1}
}

//Set uses given function f as a mock of ReplicaStorage.GetHeavySyncedPulse method
func (m *mReplicaStorageMockGetHeavySyncedPulse) Set(f func(p context.Context, p1 core.RecordID) (r core.PulseNumber, r1 error)) *ReplicaStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetHeavySyncedPulseFunc = f
	return m.mock
}

//GetHeavySyncedPulse implements github.com/insolar/insolar/ledger/storage.ReplicaStorage interface
func (m *ReplicaStorageMock) GetHeavySyncedPulse(p context.Context, p1 core.RecordID) (r core.PulseNumber, r1 error) {
	counter := atomic.AddUint64(&m.GetHeavySyncedPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetHeavySyncedPulseCounter, 1)

	if len(m.GetHeavySyncedPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetHeavySyncedPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetHeavySyncedPulse. %v %v", p, p1)
			return
		}

		input := m.GetHeavySyncedPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ReplicaStorageMockGetHeavySyncedPulseInput{p, p1}, "ReplicaStorage.GetHeavySyncedPulse got unexpected parameters")

		result := m.GetHeavySyncedPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetHeavySyncedPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetHeavySyncedPulseMock.mainExpectation != nil {

		input := m.GetHeavySyncedPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ReplicaStorageMockGetHeavySyncedPulseInput{p, p1}, "ReplicaStorage.GetHeavySyncedPulse got unexpected parameters")
		}

		result := m.GetHeavySyncedPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetHeavySyncedPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetHeavySyncedPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetHeavySyncedPulse. %v %v", p, p1)
		return
	}

	return m.GetHeavySyncedPulseFunc(p, p1)
}

//GetHeavySyncedPulseMinimockCounter returns a count of ReplicaStorageMock.GetHeavySyncedPulseFunc invocations
func (m *ReplicaStorageMock) GetHeavySyncedPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHeavySyncedPulseCounter)
}

//GetHeavySyncedPulseMinimockPreCounter returns the value of ReplicaStorageMock.GetHeavySyncedPulse invocations
func (m *ReplicaStorageMock) GetHeavySyncedPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHeavySyncedPulsePreCounter)
}

//GetHeavySyncedPulseFinished returns true if mock invocations count is ok
func (m *ReplicaStorageMock) GetHeavySyncedPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetHeavySyncedPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetHeavySyncedPulseCounter) == uint64(len(m.GetHeavySyncedPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetHeavySyncedPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetHeavySyncedPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetHeavySyncedPulseFunc != nil {
		return atomic.LoadUint64(&m.GetHeavySyncedPulseCounter) > 0
	}

	return true
}

type mReplicaStorageMockGetSyncClientJetPulses struct {
	mock              *ReplicaStorageMock
	mainExpectation   *ReplicaStorageMockGetSyncClientJetPulsesExpectation
	expectationSeries []*ReplicaStorageMockGetSyncClientJetPulsesExpectation
}

type ReplicaStorageMockGetSyncClientJetPulsesExpectation struct {
	input  *ReplicaStorageMockGetSyncClientJetPulsesInput
	result *ReplicaStorageMockGetSyncClientJetPulsesResult
}

type ReplicaStorageMockGetSyncClientJetPulsesInput struct {
	p  context.Context
	p1 core.RecordID
}

type ReplicaStorageMockGetSyncClientJetPulsesResult struct {
	r  []core.PulseNumber
	r1 error
}

//Expect specifies that invocation of ReplicaStorage.GetSyncClientJetPulses is expected from 1 to Infinity times
func (m *mReplicaStorageMockGetSyncClientJetPulses) Expect(p context.Context, p1 core.RecordID) *mReplicaStorageMockGetSyncClientJetPulses {
	m.mock.GetSyncClientJetPulsesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetSyncClientJetPulsesExpectation{}
	}
	m.mainExpectation.input = &ReplicaStorageMockGetSyncClientJetPulsesInput{p, p1}
	return m
}

//Return specifies results of invocation of ReplicaStorage.GetSyncClientJetPulses
func (m *mReplicaStorageMockGetSyncClientJetPulses) Return(r []core.PulseNumber, r1 error) *ReplicaStorageMock {
	m.mock.GetSyncClientJetPulsesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockGetSyncClientJetPulsesExpectation{}
	}
	m.mainExpectation.result = &ReplicaStorageMockGetSyncClientJetPulsesResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ReplicaStorage.GetSyncClientJetPulses is expected once
func (m *mReplicaStorageMockGetSyncClientJetPulses) ExpectOnce(p context.Context, p1 core.RecordID) *ReplicaStorageMockGetSyncClientJetPulsesExpectation {
	m.mock.GetSyncClientJetPulsesFunc = nil
	m.mainExpectation = nil

	expectation := &ReplicaStorageMockGetSyncClientJetPulsesExpectation{}
	expectation.input = &ReplicaStorageMockGetSyncClientJetPulsesInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReplicaStorageMockGetSyncClientJetPulsesExpectation) Return(r []core.PulseNumber, r1 error) {
	e.result = &ReplicaStorageMockGetSyncClientJetPulsesResult{r, r1}
}

//Set uses given function f as a mock of ReplicaStorage.GetSyncClientJetPulses method
func (m *mReplicaStorageMockGetSyncClientJetPulses) Set(f func(p context.Context, p1 core.RecordID) (r []core.PulseNumber, r1 error)) *ReplicaStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSyncClientJetPulsesFunc = f
	return m.mock
}

//GetSyncClientJetPulses implements github.com/insolar/insolar/ledger/storage.ReplicaStorage interface
func (m *ReplicaStorageMock) GetSyncClientJetPulses(p context.Context, p1 core.RecordID) (r []core.PulseNumber, r1 error) {
	counter := atomic.AddUint64(&m.GetSyncClientJetPulsesPreCounter, 1)
	defer atomic.AddUint64(&m.GetSyncClientJetPulsesCounter, 1)

	if len(m.GetSyncClientJetPulsesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSyncClientJetPulsesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetSyncClientJetPulses. %v %v", p, p1)
			return
		}

		input := m.GetSyncClientJetPulsesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ReplicaStorageMockGetSyncClientJetPulsesInput{p, p1}, "ReplicaStorage.GetSyncClientJetPulses got unexpected parameters")

		result := m.GetSyncClientJetPulsesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetSyncClientJetPulses")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetSyncClientJetPulsesMock.mainExpectation != nil {

		input := m.GetSyncClientJetPulsesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ReplicaStorageMockGetSyncClientJetPulsesInput{p, p1}, "ReplicaStorage.GetSyncClientJetPulses got unexpected parameters")
		}

		result := m.GetSyncClientJetPulsesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.GetSyncClientJetPulses")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetSyncClientJetPulsesFunc == nil {
		m.t.Fatalf("Unexpected call to ReplicaStorageMock.GetSyncClientJetPulses. %v %v", p, p1)
		return
	}

	return m.GetSyncClientJetPulsesFunc(p, p1)
}

//GetSyncClientJetPulsesMinimockCounter returns a count of ReplicaStorageMock.GetSyncClientJetPulsesFunc invocations
func (m *ReplicaStorageMock) GetSyncClientJetPulsesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSyncClientJetPulsesCounter)
}

//GetSyncClientJetPulsesMinimockPreCounter returns the value of ReplicaStorageMock.GetSyncClientJetPulses invocations
func (m *ReplicaStorageMock) GetSyncClientJetPulsesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSyncClientJetPulsesPreCounter)
}

//GetSyncClientJetPulsesFinished returns true if mock invocations count is ok
func (m *ReplicaStorageMock) GetSyncClientJetPulsesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSyncClientJetPulsesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSyncClientJetPulsesCounter) == uint64(len(m.GetSyncClientJetPulsesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSyncClientJetPulsesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSyncClientJetPulsesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSyncClientJetPulsesFunc != nil {
		return atomic.LoadUint64(&m.GetSyncClientJetPulsesCounter) > 0
	}

	return true
}

type mReplicaStorageMockSetHeavySyncedPulse struct {
	mock              *ReplicaStorageMock
	mainExpectation   *ReplicaStorageMockSetHeavySyncedPulseExpectation
	expectationSeries []*ReplicaStorageMockSetHeavySyncedPulseExpectation
}

type ReplicaStorageMockSetHeavySyncedPulseExpectation struct {
	input  *ReplicaStorageMockSetHeavySyncedPulseInput
	result *ReplicaStorageMockSetHeavySyncedPulseResult
}

type ReplicaStorageMockSetHeavySyncedPulseInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

type ReplicaStorageMockSetHeavySyncedPulseResult struct {
	r error
}

//Expect specifies that invocation of ReplicaStorage.SetHeavySyncedPulse is expected from 1 to Infinity times
func (m *mReplicaStorageMockSetHeavySyncedPulse) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mReplicaStorageMockSetHeavySyncedPulse {
	m.mock.SetHeavySyncedPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockSetHeavySyncedPulseExpectation{}
	}
	m.mainExpectation.input = &ReplicaStorageMockSetHeavySyncedPulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ReplicaStorage.SetHeavySyncedPulse
func (m *mReplicaStorageMockSetHeavySyncedPulse) Return(r error) *ReplicaStorageMock {
	m.mock.SetHeavySyncedPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockSetHeavySyncedPulseExpectation{}
	}
	m.mainExpectation.result = &ReplicaStorageMockSetHeavySyncedPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ReplicaStorage.SetHeavySyncedPulse is expected once
func (m *mReplicaStorageMockSetHeavySyncedPulse) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *ReplicaStorageMockSetHeavySyncedPulseExpectation {
	m.mock.SetHeavySyncedPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ReplicaStorageMockSetHeavySyncedPulseExpectation{}
	expectation.input = &ReplicaStorageMockSetHeavySyncedPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReplicaStorageMockSetHeavySyncedPulseExpectation) Return(r error) {
	e.result = &ReplicaStorageMockSetHeavySyncedPulseResult{r}
}

//Set uses given function f as a mock of ReplicaStorage.SetHeavySyncedPulse method
func (m *mReplicaStorageMockSetHeavySyncedPulse) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)) *ReplicaStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetHeavySyncedPulseFunc = f
	return m.mock
}

//SetHeavySyncedPulse implements github.com/insolar/insolar/ledger/storage.ReplicaStorage interface
func (m *ReplicaStorageMock) SetHeavySyncedPulse(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.SetHeavySyncedPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetHeavySyncedPulseCounter, 1)

	if len(m.SetHeavySyncedPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetHeavySyncedPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReplicaStorageMock.SetHeavySyncedPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetHeavySyncedPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ReplicaStorageMockSetHeavySyncedPulseInput{p, p1, p2}, "ReplicaStorage.SetHeavySyncedPulse got unexpected parameters")

		result := m.SetHeavySyncedPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.SetHeavySyncedPulse")
			return
		}

		r = result.r

		return
	}

	if m.SetHeavySyncedPulseMock.mainExpectation != nil {

		input := m.SetHeavySyncedPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ReplicaStorageMockSetHeavySyncedPulseInput{p, p1, p2}, "ReplicaStorage.SetHeavySyncedPulse got unexpected parameters")
		}

		result := m.SetHeavySyncedPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.SetHeavySyncedPulse")
		}

		r = result.r

		return
	}

	if m.SetHeavySyncedPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ReplicaStorageMock.SetHeavySyncedPulse. %v %v %v", p, p1, p2)
		return
	}

	return m.SetHeavySyncedPulseFunc(p, p1, p2)
}

//SetHeavySyncedPulseMinimockCounter returns a count of ReplicaStorageMock.SetHeavySyncedPulseFunc invocations
func (m *ReplicaStorageMock) SetHeavySyncedPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetHeavySyncedPulseCounter)
}

//SetHeavySyncedPulseMinimockPreCounter returns the value of ReplicaStorageMock.SetHeavySyncedPulse invocations
func (m *ReplicaStorageMock) SetHeavySyncedPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetHeavySyncedPulsePreCounter)
}

//SetHeavySyncedPulseFinished returns true if mock invocations count is ok
func (m *ReplicaStorageMock) SetHeavySyncedPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetHeavySyncedPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetHeavySyncedPulseCounter) == uint64(len(m.SetHeavySyncedPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetHeavySyncedPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetHeavySyncedPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetHeavySyncedPulseFunc != nil {
		return atomic.LoadUint64(&m.SetHeavySyncedPulseCounter) > 0
	}

	return true
}

type mReplicaStorageMockSetSyncClientJetPulses struct {
	mock              *ReplicaStorageMock
	mainExpectation   *ReplicaStorageMockSetSyncClientJetPulsesExpectation
	expectationSeries []*ReplicaStorageMockSetSyncClientJetPulsesExpectation
}

type ReplicaStorageMockSetSyncClientJetPulsesExpectation struct {
	input  *ReplicaStorageMockSetSyncClientJetPulsesInput
	result *ReplicaStorageMockSetSyncClientJetPulsesResult
}

type ReplicaStorageMockSetSyncClientJetPulsesInput struct {
	p  context.Context
	p1 core.RecordID
	p2 []core.PulseNumber
}

type ReplicaStorageMockSetSyncClientJetPulsesResult struct {
	r error
}

//Expect specifies that invocation of ReplicaStorage.SetSyncClientJetPulses is expected from 1 to Infinity times
func (m *mReplicaStorageMockSetSyncClientJetPulses) Expect(p context.Context, p1 core.RecordID, p2 []core.PulseNumber) *mReplicaStorageMockSetSyncClientJetPulses {
	m.mock.SetSyncClientJetPulsesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockSetSyncClientJetPulsesExpectation{}
	}
	m.mainExpectation.input = &ReplicaStorageMockSetSyncClientJetPulsesInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ReplicaStorage.SetSyncClientJetPulses
func (m *mReplicaStorageMockSetSyncClientJetPulses) Return(r error) *ReplicaStorageMock {
	m.mock.SetSyncClientJetPulsesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReplicaStorageMockSetSyncClientJetPulsesExpectation{}
	}
	m.mainExpectation.result = &ReplicaStorageMockSetSyncClientJetPulsesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ReplicaStorage.SetSyncClientJetPulses is expected once
func (m *mReplicaStorageMockSetSyncClientJetPulses) ExpectOnce(p context.Context, p1 core.RecordID, p2 []core.PulseNumber) *ReplicaStorageMockSetSyncClientJetPulsesExpectation {
	m.mock.SetSyncClientJetPulsesFunc = nil
	m.mainExpectation = nil

	expectation := &ReplicaStorageMockSetSyncClientJetPulsesExpectation{}
	expectation.input = &ReplicaStorageMockSetSyncClientJetPulsesInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReplicaStorageMockSetSyncClientJetPulsesExpectation) Return(r error) {
	e.result = &ReplicaStorageMockSetSyncClientJetPulsesResult{r}
}

//Set uses given function f as a mock of ReplicaStorage.SetSyncClientJetPulses method
func (m *mReplicaStorageMockSetSyncClientJetPulses) Set(f func(p context.Context, p1 core.RecordID, p2 []core.PulseNumber) (r error)) *ReplicaStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetSyncClientJetPulsesFunc = f
	return m.mock
}

//SetSyncClientJetPulses implements github.com/insolar/insolar/ledger/storage.ReplicaStorage interface
func (m *ReplicaStorageMock) SetSyncClientJetPulses(p context.Context, p1 core.RecordID, p2 []core.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.SetSyncClientJetPulsesPreCounter, 1)
	defer atomic.AddUint64(&m.SetSyncClientJetPulsesCounter, 1)

	if len(m.SetSyncClientJetPulsesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetSyncClientJetPulsesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReplicaStorageMock.SetSyncClientJetPulses. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetSyncClientJetPulsesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ReplicaStorageMockSetSyncClientJetPulsesInput{p, p1, p2}, "ReplicaStorage.SetSyncClientJetPulses got unexpected parameters")

		result := m.SetSyncClientJetPulsesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.SetSyncClientJetPulses")
			return
		}

		r = result.r

		return
	}

	if m.SetSyncClientJetPulsesMock.mainExpectation != nil {

		input := m.SetSyncClientJetPulsesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ReplicaStorageMockSetSyncClientJetPulsesInput{p, p1, p2}, "ReplicaStorage.SetSyncClientJetPulses got unexpected parameters")
		}

		result := m.SetSyncClientJetPulsesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReplicaStorageMock.SetSyncClientJetPulses")
		}

		r = result.r

		return
	}

	if m.SetSyncClientJetPulsesFunc == nil {
		m.t.Fatalf("Unexpected call to ReplicaStorageMock.SetSyncClientJetPulses. %v %v %v", p, p1, p2)
		return
	}

	return m.SetSyncClientJetPulsesFunc(p, p1, p2)
}

//SetSyncClientJetPulsesMinimockCounter returns a count of ReplicaStorageMock.SetSyncClientJetPulsesFunc invocations
func (m *ReplicaStorageMock) SetSyncClientJetPulsesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetSyncClientJetPulsesCounter)
}

//SetSyncClientJetPulsesMinimockPreCounter returns the value of ReplicaStorageMock.SetSyncClientJetPulses invocations
func (m *ReplicaStorageMock) SetSyncClientJetPulsesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetSyncClientJetPulsesPreCounter)
}

//SetSyncClientJetPulsesFinished returns true if mock invocations count is ok
func (m *ReplicaStorageMock) SetSyncClientJetPulsesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetSyncClientJetPulsesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetSyncClientJetPulsesCounter) == uint64(len(m.SetSyncClientJetPulsesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetSyncClientJetPulsesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetSyncClientJetPulsesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetSyncClientJetPulsesFunc != nil {
		return atomic.LoadUint64(&m.SetSyncClientJetPulsesCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ReplicaStorageMock) ValidateCallCounters() {

	if !m.GetAllNonEmptySyncClientJetsFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetAllNonEmptySyncClientJets")
	}

	if !m.GetAllSyncClientJetsFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetAllSyncClientJets")
	}

	if !m.GetHeavySyncedPulseFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetHeavySyncedPulse")
	}

	if !m.GetSyncClientJetPulsesFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetSyncClientJetPulses")
	}

	if !m.SetHeavySyncedPulseFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.SetHeavySyncedPulse")
	}

	if !m.SetSyncClientJetPulsesFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.SetSyncClientJetPulses")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ReplicaStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ReplicaStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ReplicaStorageMock) MinimockFinish() {

	if !m.GetAllNonEmptySyncClientJetsFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetAllNonEmptySyncClientJets")
	}

	if !m.GetAllSyncClientJetsFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetAllSyncClientJets")
	}

	if !m.GetHeavySyncedPulseFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetHeavySyncedPulse")
	}

	if !m.GetSyncClientJetPulsesFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.GetSyncClientJetPulses")
	}

	if !m.SetHeavySyncedPulseFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.SetHeavySyncedPulse")
	}

	if !m.SetSyncClientJetPulsesFinished() {
		m.t.Fatal("Expected call to ReplicaStorageMock.SetSyncClientJetPulses")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ReplicaStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ReplicaStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetAllNonEmptySyncClientJetsFinished()
		ok = ok && m.GetAllSyncClientJetsFinished()
		ok = ok && m.GetHeavySyncedPulseFinished()
		ok = ok && m.GetSyncClientJetPulsesFinished()
		ok = ok && m.SetHeavySyncedPulseFinished()
		ok = ok && m.SetSyncClientJetPulsesFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAllNonEmptySyncClientJetsFinished() {
				m.t.Error("Expected call to ReplicaStorageMock.GetAllNonEmptySyncClientJets")
			}

			if !m.GetAllSyncClientJetsFinished() {
				m.t.Error("Expected call to ReplicaStorageMock.GetAllSyncClientJets")
			}

			if !m.GetHeavySyncedPulseFinished() {
				m.t.Error("Expected call to ReplicaStorageMock.GetHeavySyncedPulse")
			}

			if !m.GetSyncClientJetPulsesFinished() {
				m.t.Error("Expected call to ReplicaStorageMock.GetSyncClientJetPulses")
			}

			if !m.SetHeavySyncedPulseFinished() {
				m.t.Error("Expected call to ReplicaStorageMock.SetHeavySyncedPulse")
			}

			if !m.SetSyncClientJetPulsesFinished() {
				m.t.Error("Expected call to ReplicaStorageMock.SetSyncClientJetPulses")
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
func (m *ReplicaStorageMock) AllMocksCalled() bool {

	if !m.GetAllNonEmptySyncClientJetsFinished() {
		return false
	}

	if !m.GetAllSyncClientJetsFinished() {
		return false
	}

	if !m.GetHeavySyncedPulseFinished() {
		return false
	}

	if !m.GetSyncClientJetPulsesFinished() {
		return false
	}

	if !m.SetHeavySyncedPulseFinished() {
		return false
	}

	if !m.SetSyncClientJetPulsesFinished() {
		return false
	}

	return true
}

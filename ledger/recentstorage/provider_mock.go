package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Provider" can be found in github.com/insolar/insolar/ledger/recentstorage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProviderMock implements github.com/insolar/insolar/ledger/recentstorage.Provider
type ProviderMock struct {
	t minimock.Tester

	CloneIndexStorageFunc       func(p context.Context, p1 core.RecordID, p2 core.RecordID)
	CloneIndexStorageCounter    uint64
	CloneIndexStoragePreCounter uint64
	CloneIndexStorageMock       mProviderMockCloneIndexStorage

	ClonePendingStorageFunc       func(p context.Context, p1 core.RecordID, p2 core.RecordID)
	ClonePendingStorageCounter    uint64
	ClonePendingStoragePreCounter uint64
	ClonePendingStorageMock       mProviderMockClonePendingStorage

	DecreaseIndexesTTLFunc       func(p context.Context) (r map[core.RecordID][]core.RecordID)
	DecreaseIndexesTTLCounter    uint64
	DecreaseIndexesTTLPreCounter uint64
	DecreaseIndexesTTLMock       mProviderMockDecreaseIndexesTTL

	GetIndexStorageFunc       func(p context.Context, p1 core.RecordID) (r RecentIndexStorage)
	GetIndexStorageCounter    uint64
	GetIndexStoragePreCounter uint64
	GetIndexStorageMock       mProviderMockGetIndexStorage

	GetPendingStorageFunc       func(p context.Context, p1 core.RecordID) (r PendingStorage)
	GetPendingStorageCounter    uint64
	GetPendingStoragePreCounter uint64
	GetPendingStorageMock       mProviderMockGetPendingStorage

	RemoveIndexStorageFunc       func(p context.Context, p1 core.RecordID)
	RemoveIndexStorageCounter    uint64
	RemoveIndexStoragePreCounter uint64
	RemoveIndexStorageMock       mProviderMockRemoveIndexStorage

	RemovePendingStorageFunc       func(p context.Context, p1 core.RecordID)
	RemovePendingStorageCounter    uint64
	RemovePendingStoragePreCounter uint64
	RemovePendingStorageMock       mProviderMockRemovePendingStorage
}

//NewProviderMock returns a mock for github.com/insolar/insolar/ledger/recentstorage.Provider
func NewProviderMock(t minimock.Tester) *ProviderMock {
	m := &ProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloneIndexStorageMock = mProviderMockCloneIndexStorage{mock: m}
	m.ClonePendingStorageMock = mProviderMockClonePendingStorage{mock: m}
	m.DecreaseIndexesTTLMock = mProviderMockDecreaseIndexesTTL{mock: m}
	m.GetIndexStorageMock = mProviderMockGetIndexStorage{mock: m}
	m.GetPendingStorageMock = mProviderMockGetPendingStorage{mock: m}
	m.RemoveIndexStorageMock = mProviderMockRemoveIndexStorage{mock: m}
	m.RemovePendingStorageMock = mProviderMockRemovePendingStorage{mock: m}

	return m
}

type mProviderMockCloneIndexStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockCloneIndexStorageExpectation
	expectationSeries []*ProviderMockCloneIndexStorageExpectation
}

type ProviderMockCloneIndexStorageExpectation struct {
	input *ProviderMockCloneIndexStorageInput
}

type ProviderMockCloneIndexStorageInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.RecordID
}

//Expect specifies that invocation of Provider.CloneIndexStorage is expected from 1 to Infinity times
func (m *mProviderMockCloneIndexStorage) Expect(p context.Context, p1 core.RecordID, p2 core.RecordID) *mProviderMockCloneIndexStorage {
	m.mock.CloneIndexStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCloneIndexStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockCloneIndexStorageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Provider.CloneIndexStorage
func (m *mProviderMockCloneIndexStorage) Return() *ProviderMock {
	m.mock.CloneIndexStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCloneIndexStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.CloneIndexStorage is expected once
func (m *mProviderMockCloneIndexStorage) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.RecordID) *ProviderMockCloneIndexStorageExpectation {
	m.mock.CloneIndexStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockCloneIndexStorageExpectation{}
	expectation.input = &ProviderMockCloneIndexStorageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.CloneIndexStorage method
func (m *mProviderMockCloneIndexStorage) Set(f func(p context.Context, p1 core.RecordID, p2 core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloneIndexStorageFunc = f
	return m.mock
}

//CloneIndexStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) CloneIndexStorage(p context.Context, p1 core.RecordID, p2 core.RecordID) {
	counter := atomic.AddUint64(&m.CloneIndexStoragePreCounter, 1)
	defer atomic.AddUint64(&m.CloneIndexStorageCounter, 1)

	if len(m.CloneIndexStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloneIndexStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.CloneIndexStorage. %v %v %v", p, p1, p2)
			return
		}

		input := m.CloneIndexStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockCloneIndexStorageInput{p, p1, p2}, "Provider.CloneIndexStorage got unexpected parameters")

		return
	}

	if m.CloneIndexStorageMock.mainExpectation != nil {

		input := m.CloneIndexStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockCloneIndexStorageInput{p, p1, p2}, "Provider.CloneIndexStorage got unexpected parameters")
		}

		return
	}

	if m.CloneIndexStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.CloneIndexStorage. %v %v %v", p, p1, p2)
		return
	}

	m.CloneIndexStorageFunc(p, p1, p2)
}

//CloneIndexStorageMinimockCounter returns a count of ProviderMock.CloneIndexStorageFunc invocations
func (m *ProviderMock) CloneIndexStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloneIndexStorageCounter)
}

//CloneIndexStorageMinimockPreCounter returns the value of ProviderMock.CloneIndexStorage invocations
func (m *ProviderMock) CloneIndexStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CloneIndexStoragePreCounter)
}

//CloneIndexStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) CloneIndexStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloneIndexStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloneIndexStorageCounter) == uint64(len(m.CloneIndexStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloneIndexStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloneIndexStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloneIndexStorageFunc != nil {
		return atomic.LoadUint64(&m.CloneIndexStorageCounter) > 0
	}

	return true
}

type mProviderMockClonePendingStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockClonePendingStorageExpectation
	expectationSeries []*ProviderMockClonePendingStorageExpectation
}

type ProviderMockClonePendingStorageExpectation struct {
	input *ProviderMockClonePendingStorageInput
}

type ProviderMockClonePendingStorageInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.RecordID
}

//Expect specifies that invocation of Provider.ClonePendingStorage is expected from 1 to Infinity times
func (m *mProviderMockClonePendingStorage) Expect(p context.Context, p1 core.RecordID, p2 core.RecordID) *mProviderMockClonePendingStorage {
	m.mock.ClonePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockClonePendingStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockClonePendingStorageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Provider.ClonePendingStorage
func (m *mProviderMockClonePendingStorage) Return() *ProviderMock {
	m.mock.ClonePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockClonePendingStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.ClonePendingStorage is expected once
func (m *mProviderMockClonePendingStorage) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.RecordID) *ProviderMockClonePendingStorageExpectation {
	m.mock.ClonePendingStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockClonePendingStorageExpectation{}
	expectation.input = &ProviderMockClonePendingStorageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.ClonePendingStorage method
func (m *mProviderMockClonePendingStorage) Set(f func(p context.Context, p1 core.RecordID, p2 core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ClonePendingStorageFunc = f
	return m.mock
}

//ClonePendingStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) ClonePendingStorage(p context.Context, p1 core.RecordID, p2 core.RecordID) {
	counter := atomic.AddUint64(&m.ClonePendingStoragePreCounter, 1)
	defer atomic.AddUint64(&m.ClonePendingStorageCounter, 1)

	if len(m.ClonePendingStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ClonePendingStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.ClonePendingStorage. %v %v %v", p, p1, p2)
			return
		}

		input := m.ClonePendingStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockClonePendingStorageInput{p, p1, p2}, "Provider.ClonePendingStorage got unexpected parameters")

		return
	}

	if m.ClonePendingStorageMock.mainExpectation != nil {

		input := m.ClonePendingStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockClonePendingStorageInput{p, p1, p2}, "Provider.ClonePendingStorage got unexpected parameters")
		}

		return
	}

	if m.ClonePendingStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.ClonePendingStorage. %v %v %v", p, p1, p2)
		return
	}

	m.ClonePendingStorageFunc(p, p1, p2)
}

//ClonePendingStorageMinimockCounter returns a count of ProviderMock.ClonePendingStorageFunc invocations
func (m *ProviderMock) ClonePendingStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ClonePendingStorageCounter)
}

//ClonePendingStorageMinimockPreCounter returns the value of ProviderMock.ClonePendingStorage invocations
func (m *ProviderMock) ClonePendingStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClonePendingStoragePreCounter)
}

//ClonePendingStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) ClonePendingStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ClonePendingStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ClonePendingStorageCounter) == uint64(len(m.ClonePendingStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ClonePendingStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ClonePendingStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ClonePendingStorageFunc != nil {
		return atomic.LoadUint64(&m.ClonePendingStorageCounter) > 0
	}

	return true
}

type mProviderMockDecreaseIndexesTTL struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockDecreaseIndexesTTLExpectation
	expectationSeries []*ProviderMockDecreaseIndexesTTLExpectation
}

type ProviderMockDecreaseIndexesTTLExpectation struct {
	input  *ProviderMockDecreaseIndexesTTLInput
	result *ProviderMockDecreaseIndexesTTLResult
}

type ProviderMockDecreaseIndexesTTLInput struct {
	p context.Context
}

type ProviderMockDecreaseIndexesTTLResult struct {
	r map[core.RecordID][]core.RecordID
}

//Expect specifies that invocation of Provider.DecreaseIndexesTTL is expected from 1 to Infinity times
func (m *mProviderMockDecreaseIndexesTTL) Expect(p context.Context) *mProviderMockDecreaseIndexesTTL {
	m.mock.DecreaseIndexesTTLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockDecreaseIndexesTTLExpectation{}
	}
	m.mainExpectation.input = &ProviderMockDecreaseIndexesTTLInput{p}
	return m
}

//Return specifies results of invocation of Provider.DecreaseIndexesTTL
func (m *mProviderMockDecreaseIndexesTTL) Return(r map[core.RecordID][]core.RecordID) *ProviderMock {
	m.mock.DecreaseIndexesTTLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockDecreaseIndexesTTLExpectation{}
	}
	m.mainExpectation.result = &ProviderMockDecreaseIndexesTTLResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.DecreaseIndexesTTL is expected once
func (m *mProviderMockDecreaseIndexesTTL) ExpectOnce(p context.Context) *ProviderMockDecreaseIndexesTTLExpectation {
	m.mock.DecreaseIndexesTTLFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockDecreaseIndexesTTLExpectation{}
	expectation.input = &ProviderMockDecreaseIndexesTTLInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockDecreaseIndexesTTLExpectation) Return(r map[core.RecordID][]core.RecordID) {
	e.result = &ProviderMockDecreaseIndexesTTLResult{r}
}

//Set uses given function f as a mock of Provider.DecreaseIndexesTTL method
func (m *mProviderMockDecreaseIndexesTTL) Set(f func(p context.Context) (r map[core.RecordID][]core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DecreaseIndexesTTLFunc = f
	return m.mock
}

//DecreaseIndexesTTL implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) DecreaseIndexesTTL(p context.Context) (r map[core.RecordID][]core.RecordID) {
	counter := atomic.AddUint64(&m.DecreaseIndexesTTLPreCounter, 1)
	defer atomic.AddUint64(&m.DecreaseIndexesTTLCounter, 1)

	if len(m.DecreaseIndexesTTLMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DecreaseIndexesTTLMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.DecreaseIndexesTTL. %v", p)
			return
		}

		input := m.DecreaseIndexesTTLMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockDecreaseIndexesTTLInput{p}, "Provider.DecreaseIndexesTTL got unexpected parameters")

		result := m.DecreaseIndexesTTLMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.DecreaseIndexesTTL")
			return
		}

		r = result.r

		return
	}

	if m.DecreaseIndexesTTLMock.mainExpectation != nil {

		input := m.DecreaseIndexesTTLMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockDecreaseIndexesTTLInput{p}, "Provider.DecreaseIndexesTTL got unexpected parameters")
		}

		result := m.DecreaseIndexesTTLMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.DecreaseIndexesTTL")
		}

		r = result.r

		return
	}

	if m.DecreaseIndexesTTLFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.DecreaseIndexesTTL. %v", p)
		return
	}

	return m.DecreaseIndexesTTLFunc(p)
}

//DecreaseIndexesTTLMinimockCounter returns a count of ProviderMock.DecreaseIndexesTTLFunc invocations
func (m *ProviderMock) DecreaseIndexesTTLMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DecreaseIndexesTTLCounter)
}

//DecreaseIndexesTTLMinimockPreCounter returns the value of ProviderMock.DecreaseIndexesTTL invocations
func (m *ProviderMock) DecreaseIndexesTTLMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DecreaseIndexesTTLPreCounter)
}

//DecreaseIndexesTTLFinished returns true if mock invocations count is ok
func (m *ProviderMock) DecreaseIndexesTTLFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DecreaseIndexesTTLMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DecreaseIndexesTTLCounter) == uint64(len(m.DecreaseIndexesTTLMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DecreaseIndexesTTLMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DecreaseIndexesTTLCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DecreaseIndexesTTLFunc != nil {
		return atomic.LoadUint64(&m.DecreaseIndexesTTLCounter) > 0
	}

	return true
}

type mProviderMockGetIndexStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockGetIndexStorageExpectation
	expectationSeries []*ProviderMockGetIndexStorageExpectation
}

type ProviderMockGetIndexStorageExpectation struct {
	input  *ProviderMockGetIndexStorageInput
	result *ProviderMockGetIndexStorageResult
}

type ProviderMockGetIndexStorageInput struct {
	p  context.Context
	p1 core.RecordID
}

type ProviderMockGetIndexStorageResult struct {
	r RecentIndexStorage
}

//Expect specifies that invocation of Provider.GetIndexStorage is expected from 1 to Infinity times
func (m *mProviderMockGetIndexStorage) Expect(p context.Context, p1 core.RecordID) *mProviderMockGetIndexStorage {
	m.mock.GetIndexStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetIndexStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockGetIndexStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.GetIndexStorage
func (m *mProviderMockGetIndexStorage) Return(r RecentIndexStorage) *ProviderMock {
	m.mock.GetIndexStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetIndexStorageExpectation{}
	}
	m.mainExpectation.result = &ProviderMockGetIndexStorageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.GetIndexStorage is expected once
func (m *mProviderMockGetIndexStorage) ExpectOnce(p context.Context, p1 core.RecordID) *ProviderMockGetIndexStorageExpectation {
	m.mock.GetIndexStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockGetIndexStorageExpectation{}
	expectation.input = &ProviderMockGetIndexStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockGetIndexStorageExpectation) Return(r RecentIndexStorage) {
	e.result = &ProviderMockGetIndexStorageResult{r}
}

//Set uses given function f as a mock of Provider.GetIndexStorage method
func (m *mProviderMockGetIndexStorage) Set(f func(p context.Context, p1 core.RecordID) (r RecentIndexStorage)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIndexStorageFunc = f
	return m.mock
}

//GetIndexStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) GetIndexStorage(p context.Context, p1 core.RecordID) (r RecentIndexStorage) {
	counter := atomic.AddUint64(&m.GetIndexStoragePreCounter, 1)
	defer atomic.AddUint64(&m.GetIndexStorageCounter, 1)

	if len(m.GetIndexStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIndexStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.GetIndexStorage. %v %v", p, p1)
			return
		}

		input := m.GetIndexStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockGetIndexStorageInput{p, p1}, "Provider.GetIndexStorage got unexpected parameters")

		result := m.GetIndexStorageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetIndexStorage")
			return
		}

		r = result.r

		return
	}

	if m.GetIndexStorageMock.mainExpectation != nil {

		input := m.GetIndexStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockGetIndexStorageInput{p, p1}, "Provider.GetIndexStorage got unexpected parameters")
		}

		result := m.GetIndexStorageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetIndexStorage")
		}

		r = result.r

		return
	}

	if m.GetIndexStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.GetIndexStorage. %v %v", p, p1)
		return
	}

	return m.GetIndexStorageFunc(p, p1)
}

//GetIndexStorageMinimockCounter returns a count of ProviderMock.GetIndexStorageFunc invocations
func (m *ProviderMock) GetIndexStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexStorageCounter)
}

//GetIndexStorageMinimockPreCounter returns the value of ProviderMock.GetIndexStorage invocations
func (m *ProviderMock) GetIndexStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexStoragePreCounter)
}

//GetIndexStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) GetIndexStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIndexStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIndexStorageCounter) == uint64(len(m.GetIndexStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIndexStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIndexStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIndexStorageFunc != nil {
		return atomic.LoadUint64(&m.GetIndexStorageCounter) > 0
	}

	return true
}

type mProviderMockGetPendingStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockGetPendingStorageExpectation
	expectationSeries []*ProviderMockGetPendingStorageExpectation
}

type ProviderMockGetPendingStorageExpectation struct {
	input  *ProviderMockGetPendingStorageInput
	result *ProviderMockGetPendingStorageResult
}

type ProviderMockGetPendingStorageInput struct {
	p  context.Context
	p1 core.RecordID
}

type ProviderMockGetPendingStorageResult struct {
	r PendingStorage
}

//Expect specifies that invocation of Provider.GetPendingStorage is expected from 1 to Infinity times
func (m *mProviderMockGetPendingStorage) Expect(p context.Context, p1 core.RecordID) *mProviderMockGetPendingStorage {
	m.mock.GetPendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetPendingStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockGetPendingStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.GetPendingStorage
func (m *mProviderMockGetPendingStorage) Return(r PendingStorage) *ProviderMock {
	m.mock.GetPendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetPendingStorageExpectation{}
	}
	m.mainExpectation.result = &ProviderMockGetPendingStorageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.GetPendingStorage is expected once
func (m *mProviderMockGetPendingStorage) ExpectOnce(p context.Context, p1 core.RecordID) *ProviderMockGetPendingStorageExpectation {
	m.mock.GetPendingStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockGetPendingStorageExpectation{}
	expectation.input = &ProviderMockGetPendingStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockGetPendingStorageExpectation) Return(r PendingStorage) {
	e.result = &ProviderMockGetPendingStorageResult{r}
}

//Set uses given function f as a mock of Provider.GetPendingStorage method
func (m *mProviderMockGetPendingStorage) Set(f func(p context.Context, p1 core.RecordID) (r PendingStorage)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPendingStorageFunc = f
	return m.mock
}

//GetPendingStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) GetPendingStorage(p context.Context, p1 core.RecordID) (r PendingStorage) {
	counter := atomic.AddUint64(&m.GetPendingStoragePreCounter, 1)
	defer atomic.AddUint64(&m.GetPendingStorageCounter, 1)

	if len(m.GetPendingStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPendingStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.GetPendingStorage. %v %v", p, p1)
			return
		}

		input := m.GetPendingStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockGetPendingStorageInput{p, p1}, "Provider.GetPendingStorage got unexpected parameters")

		result := m.GetPendingStorageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetPendingStorage")
			return
		}

		r = result.r

		return
	}

	if m.GetPendingStorageMock.mainExpectation != nil {

		input := m.GetPendingStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockGetPendingStorageInput{p, p1}, "Provider.GetPendingStorage got unexpected parameters")
		}

		result := m.GetPendingStorageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetPendingStorage")
		}

		r = result.r

		return
	}

	if m.GetPendingStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.GetPendingStorage. %v %v", p, p1)
		return
	}

	return m.GetPendingStorageFunc(p, p1)
}

//GetPendingStorageMinimockCounter returns a count of ProviderMock.GetPendingStorageFunc invocations
func (m *ProviderMock) GetPendingStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPendingStorageCounter)
}

//GetPendingStorageMinimockPreCounter returns the value of ProviderMock.GetPendingStorage invocations
func (m *ProviderMock) GetPendingStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPendingStoragePreCounter)
}

//GetPendingStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) GetPendingStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPendingStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPendingStorageCounter) == uint64(len(m.GetPendingStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPendingStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPendingStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPendingStorageFunc != nil {
		return atomic.LoadUint64(&m.GetPendingStorageCounter) > 0
	}

	return true
}

type mProviderMockRemoveIndexStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockRemoveIndexStorageExpectation
	expectationSeries []*ProviderMockRemoveIndexStorageExpectation
}

type ProviderMockRemoveIndexStorageExpectation struct {
	input *ProviderMockRemoveIndexStorageInput
}

type ProviderMockRemoveIndexStorageInput struct {
	p  context.Context
	p1 core.RecordID
}

//Expect specifies that invocation of Provider.RemoveIndexStorage is expected from 1 to Infinity times
func (m *mProviderMockRemoveIndexStorage) Expect(p context.Context, p1 core.RecordID) *mProviderMockRemoveIndexStorage {
	m.mock.RemoveIndexStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemoveIndexStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockRemoveIndexStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.RemoveIndexStorage
func (m *mProviderMockRemoveIndexStorage) Return() *ProviderMock {
	m.mock.RemoveIndexStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemoveIndexStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.RemoveIndexStorage is expected once
func (m *mProviderMockRemoveIndexStorage) ExpectOnce(p context.Context, p1 core.RecordID) *ProviderMockRemoveIndexStorageExpectation {
	m.mock.RemoveIndexStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockRemoveIndexStorageExpectation{}
	expectation.input = &ProviderMockRemoveIndexStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.RemoveIndexStorage method
func (m *mProviderMockRemoveIndexStorage) Set(f func(p context.Context, p1 core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveIndexStorageFunc = f
	return m.mock
}

//RemoveIndexStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) RemoveIndexStorage(p context.Context, p1 core.RecordID) {
	counter := atomic.AddUint64(&m.RemoveIndexStoragePreCounter, 1)
	defer atomic.AddUint64(&m.RemoveIndexStorageCounter, 1)

	if len(m.RemoveIndexStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveIndexStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.RemoveIndexStorage. %v %v", p, p1)
			return
		}

		input := m.RemoveIndexStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockRemoveIndexStorageInput{p, p1}, "Provider.RemoveIndexStorage got unexpected parameters")

		return
	}

	if m.RemoveIndexStorageMock.mainExpectation != nil {

		input := m.RemoveIndexStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockRemoveIndexStorageInput{p, p1}, "Provider.RemoveIndexStorage got unexpected parameters")
		}

		return
	}

	if m.RemoveIndexStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.RemoveIndexStorage. %v %v", p, p1)
		return
	}

	m.RemoveIndexStorageFunc(p, p1)
}

//RemoveIndexStorageMinimockCounter returns a count of ProviderMock.RemoveIndexStorageFunc invocations
func (m *ProviderMock) RemoveIndexStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveIndexStorageCounter)
}

//RemoveIndexStorageMinimockPreCounter returns the value of ProviderMock.RemoveIndexStorage invocations
func (m *ProviderMock) RemoveIndexStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveIndexStoragePreCounter)
}

//RemoveIndexStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) RemoveIndexStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveIndexStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveIndexStorageCounter) == uint64(len(m.RemoveIndexStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveIndexStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveIndexStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveIndexStorageFunc != nil {
		return atomic.LoadUint64(&m.RemoveIndexStorageCounter) > 0
	}

	return true
}

type mProviderMockRemovePendingStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockRemovePendingStorageExpectation
	expectationSeries []*ProviderMockRemovePendingStorageExpectation
}

type ProviderMockRemovePendingStorageExpectation struct {
	input *ProviderMockRemovePendingStorageInput
}

type ProviderMockRemovePendingStorageInput struct {
	p  context.Context
	p1 core.RecordID
}

//Expect specifies that invocation of Provider.RemovePendingStorage is expected from 1 to Infinity times
func (m *mProviderMockRemovePendingStorage) Expect(p context.Context, p1 core.RecordID) *mProviderMockRemovePendingStorage {
	m.mock.RemovePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemovePendingStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockRemovePendingStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.RemovePendingStorage
func (m *mProviderMockRemovePendingStorage) Return() *ProviderMock {
	m.mock.RemovePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemovePendingStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.RemovePendingStorage is expected once
func (m *mProviderMockRemovePendingStorage) ExpectOnce(p context.Context, p1 core.RecordID) *ProviderMockRemovePendingStorageExpectation {
	m.mock.RemovePendingStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockRemovePendingStorageExpectation{}
	expectation.input = &ProviderMockRemovePendingStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.RemovePendingStorage method
func (m *mProviderMockRemovePendingStorage) Set(f func(p context.Context, p1 core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemovePendingStorageFunc = f
	return m.mock
}

//RemovePendingStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) RemovePendingStorage(p context.Context, p1 core.RecordID) {
	counter := atomic.AddUint64(&m.RemovePendingStoragePreCounter, 1)
	defer atomic.AddUint64(&m.RemovePendingStorageCounter, 1)

	if len(m.RemovePendingStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemovePendingStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.RemovePendingStorage. %v %v", p, p1)
			return
		}

		input := m.RemovePendingStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockRemovePendingStorageInput{p, p1}, "Provider.RemovePendingStorage got unexpected parameters")

		return
	}

	if m.RemovePendingStorageMock.mainExpectation != nil {

		input := m.RemovePendingStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockRemovePendingStorageInput{p, p1}, "Provider.RemovePendingStorage got unexpected parameters")
		}

		return
	}

	if m.RemovePendingStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.RemovePendingStorage. %v %v", p, p1)
		return
	}

	m.RemovePendingStorageFunc(p, p1)
}

//RemovePendingStorageMinimockCounter returns a count of ProviderMock.RemovePendingStorageFunc invocations
func (m *ProviderMock) RemovePendingStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingStorageCounter)
}

//RemovePendingStorageMinimockPreCounter returns the value of ProviderMock.RemovePendingStorage invocations
func (m *ProviderMock) RemovePendingStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingStoragePreCounter)
}

//RemovePendingStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) RemovePendingStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemovePendingStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemovePendingStorageCounter) == uint64(len(m.RemovePendingStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemovePendingStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemovePendingStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemovePendingStorageFunc != nil {
		return atomic.LoadUint64(&m.RemovePendingStorageCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) ValidateCallCounters() {

	if !m.CloneIndexStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.CloneIndexStorage")
	}

	if !m.ClonePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.ClonePendingStorage")
	}

	if !m.DecreaseIndexesTTLFinished() {
		m.t.Fatal("Expected call to ProviderMock.DecreaseIndexesTTL")
	}

	if !m.GetIndexStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetIndexStorage")
	}

	if !m.GetPendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetPendingStorage")
	}

	if !m.RemoveIndexStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemoveIndexStorage")
	}

	if !m.RemovePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemovePendingStorage")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ProviderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ProviderMock) MinimockFinish() {

	if !m.CloneIndexStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.CloneIndexStorage")
	}

	if !m.ClonePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.ClonePendingStorage")
	}

	if !m.DecreaseIndexesTTLFinished() {
		m.t.Fatal("Expected call to ProviderMock.DecreaseIndexesTTL")
	}

	if !m.GetIndexStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetIndexStorage")
	}

	if !m.GetPendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetPendingStorage")
	}

	if !m.RemoveIndexStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemoveIndexStorage")
	}

	if !m.RemovePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemovePendingStorage")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ProviderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ProviderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloneIndexStorageFinished()
		ok = ok && m.ClonePendingStorageFinished()
		ok = ok && m.DecreaseIndexesTTLFinished()
		ok = ok && m.GetIndexStorageFinished()
		ok = ok && m.GetPendingStorageFinished()
		ok = ok && m.RemoveIndexStorageFinished()
		ok = ok && m.RemovePendingStorageFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloneIndexStorageFinished() {
				m.t.Error("Expected call to ProviderMock.CloneIndexStorage")
			}

			if !m.ClonePendingStorageFinished() {
				m.t.Error("Expected call to ProviderMock.ClonePendingStorage")
			}

			if !m.DecreaseIndexesTTLFinished() {
				m.t.Error("Expected call to ProviderMock.DecreaseIndexesTTL")
			}

			if !m.GetIndexStorageFinished() {
				m.t.Error("Expected call to ProviderMock.GetIndexStorage")
			}

			if !m.GetPendingStorageFinished() {
				m.t.Error("Expected call to ProviderMock.GetPendingStorage")
			}

			if !m.RemoveIndexStorageFinished() {
				m.t.Error("Expected call to ProviderMock.RemoveIndexStorage")
			}

			if !m.RemovePendingStorageFinished() {
				m.t.Error("Expected call to ProviderMock.RemovePendingStorage")
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
func (m *ProviderMock) AllMocksCalled() bool {

	if !m.CloneIndexStorageFinished() {
		return false
	}

	if !m.ClonePendingStorageFinished() {
		return false
	}

	if !m.DecreaseIndexesTTLFinished() {
		return false
	}

	if !m.GetIndexStorageFinished() {
		return false
	}

	if !m.GetPendingStorageFinished() {
		return false
	}

	if !m.RemoveIndexStorageFinished() {
		return false
	}

	if !m.RemovePendingStorageFinished() {
		return false
	}

	return true
}

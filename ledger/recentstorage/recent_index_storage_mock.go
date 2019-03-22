package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecentIndexStorage" can be found in github.com/insolar/insolar/ledger/recentstorage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//RecentIndexStorageMock implements github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage
type RecentIndexStorageMock struct {
	t minimock.Tester

	AddObjectFunc       func(p context.Context, p1 insolar.ID)
	AddObjectCounter    uint64
	AddObjectPreCounter uint64
	AddObjectMock       mRecentIndexStorageMockAddObject

	AddObjectWithTLLFunc       func(p context.Context, p1 insolar.ID, p2 int)
	AddObjectWithTLLCounter    uint64
	AddObjectWithTLLPreCounter uint64
	AddObjectWithTLLMock       mRecentIndexStorageMockAddObjectWithTLL

	DecreaseIndexTTLFunc       func(p context.Context) (r []insolar.ID)
	DecreaseIndexTTLCounter    uint64
	DecreaseIndexTTLPreCounter uint64
	DecreaseIndexTTLMock       mRecentIndexStorageMockDecreaseIndexTTL

	FilterNotExistWithLockFunc       func(p context.Context, p1 []insolar.ID, p2 func(p []insolar.ID))
	FilterNotExistWithLockCounter    uint64
	FilterNotExistWithLockPreCounter uint64
	FilterNotExistWithLockMock       mRecentIndexStorageMockFilterNotExistWithLock

	GetObjectsFunc       func() (r map[insolar.ID]int)
	GetObjectsCounter    uint64
	GetObjectsPreCounter uint64
	GetObjectsMock       mRecentIndexStorageMockGetObjects
}

//NewRecentIndexStorageMock returns a mock for github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage
func NewRecentIndexStorageMock(t minimock.Tester) *RecentIndexStorageMock {
	m := &RecentIndexStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddObjectMock = mRecentIndexStorageMockAddObject{mock: m}
	m.AddObjectWithTLLMock = mRecentIndexStorageMockAddObjectWithTLL{mock: m}
	m.DecreaseIndexTTLMock = mRecentIndexStorageMockDecreaseIndexTTL{mock: m}
	m.FilterNotExistWithLockMock = mRecentIndexStorageMockFilterNotExistWithLock{mock: m}
	m.GetObjectsMock = mRecentIndexStorageMockGetObjects{mock: m}

	return m
}

type mRecentIndexStorageMockAddObject struct {
	mock              *RecentIndexStorageMock
	mainExpectation   *RecentIndexStorageMockAddObjectExpectation
	expectationSeries []*RecentIndexStorageMockAddObjectExpectation
}

type RecentIndexStorageMockAddObjectExpectation struct {
	input *RecentIndexStorageMockAddObjectInput
}

type RecentIndexStorageMockAddObjectInput struct {
	p  context.Context
	p1 insolar.ID
}

//Expect specifies that invocation of RecentIndexStorage.AddObject is expected from 1 to Infinity times
func (m *mRecentIndexStorageMockAddObject) Expect(p context.Context, p1 insolar.ID) *mRecentIndexStorageMockAddObject {
	m.mock.AddObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockAddObjectExpectation{}
	}
	m.mainExpectation.input = &RecentIndexStorageMockAddObjectInput{p, p1}
	return m
}

//Return specifies results of invocation of RecentIndexStorage.AddObject
func (m *mRecentIndexStorageMockAddObject) Return() *RecentIndexStorageMock {
	m.mock.AddObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockAddObjectExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentIndexStorage.AddObject is expected once
func (m *mRecentIndexStorageMockAddObject) ExpectOnce(p context.Context, p1 insolar.ID) *RecentIndexStorageMockAddObjectExpectation {
	m.mock.AddObjectFunc = nil
	m.mainExpectation = nil

	expectation := &RecentIndexStorageMockAddObjectExpectation{}
	expectation.input = &RecentIndexStorageMockAddObjectInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentIndexStorage.AddObject method
func (m *mRecentIndexStorageMockAddObject) Set(f func(p context.Context, p1 insolar.ID)) *RecentIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddObjectFunc = f
	return m.mock
}

//AddObject implements github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage interface
func (m *RecentIndexStorageMock) AddObject(p context.Context, p1 insolar.ID) {
	counter := atomic.AddUint64(&m.AddObjectPreCounter, 1)
	defer atomic.AddUint64(&m.AddObjectCounter, 1)

	if len(m.AddObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentIndexStorageMock.AddObject. %v %v", p, p1)
			return
		}

		input := m.AddObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentIndexStorageMockAddObjectInput{p, p1}, "RecentIndexStorage.AddObject got unexpected parameters")

		return
	}

	if m.AddObjectMock.mainExpectation != nil {

		input := m.AddObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentIndexStorageMockAddObjectInput{p, p1}, "RecentIndexStorage.AddObject got unexpected parameters")
		}

		return
	}

	if m.AddObjectFunc == nil {
		m.t.Fatalf("Unexpected call to RecentIndexStorageMock.AddObject. %v %v", p, p1)
		return
	}

	m.AddObjectFunc(p, p1)
}

//AddObjectMinimockCounter returns a count of RecentIndexStorageMock.AddObjectFunc invocations
func (m *RecentIndexStorageMock) AddObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectCounter)
}

//AddObjectMinimockPreCounter returns the value of RecentIndexStorageMock.AddObject invocations
func (m *RecentIndexStorageMock) AddObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectPreCounter)
}

//AddObjectFinished returns true if mock invocations count is ok
func (m *RecentIndexStorageMock) AddObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddObjectCounter) == uint64(len(m.AddObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddObjectFunc != nil {
		return atomic.LoadUint64(&m.AddObjectCounter) > 0
	}

	return true
}

type mRecentIndexStorageMockAddObjectWithTLL struct {
	mock              *RecentIndexStorageMock
	mainExpectation   *RecentIndexStorageMockAddObjectWithTLLExpectation
	expectationSeries []*RecentIndexStorageMockAddObjectWithTLLExpectation
}

type RecentIndexStorageMockAddObjectWithTLLExpectation struct {
	input *RecentIndexStorageMockAddObjectWithTLLInput
}

type RecentIndexStorageMockAddObjectWithTLLInput struct {
	p  context.Context
	p1 insolar.ID
	p2 int
}

//Expect specifies that invocation of RecentIndexStorage.AddObjectWithTLL is expected from 1 to Infinity times
func (m *mRecentIndexStorageMockAddObjectWithTLL) Expect(p context.Context, p1 insolar.ID, p2 int) *mRecentIndexStorageMockAddObjectWithTLL {
	m.mock.AddObjectWithTLLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockAddObjectWithTLLExpectation{}
	}
	m.mainExpectation.input = &RecentIndexStorageMockAddObjectWithTLLInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RecentIndexStorage.AddObjectWithTLL
func (m *mRecentIndexStorageMockAddObjectWithTLL) Return() *RecentIndexStorageMock {
	m.mock.AddObjectWithTLLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockAddObjectWithTLLExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentIndexStorage.AddObjectWithTLL is expected once
func (m *mRecentIndexStorageMockAddObjectWithTLL) ExpectOnce(p context.Context, p1 insolar.ID, p2 int) *RecentIndexStorageMockAddObjectWithTLLExpectation {
	m.mock.AddObjectWithTLLFunc = nil
	m.mainExpectation = nil

	expectation := &RecentIndexStorageMockAddObjectWithTLLExpectation{}
	expectation.input = &RecentIndexStorageMockAddObjectWithTLLInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentIndexStorage.AddObjectWithTLL method
func (m *mRecentIndexStorageMockAddObjectWithTLL) Set(f func(p context.Context, p1 insolar.ID, p2 int)) *RecentIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddObjectWithTLLFunc = f
	return m.mock
}

//AddObjectWithTLL implements github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage interface
func (m *RecentIndexStorageMock) AddObjectWithTLL(p context.Context, p1 insolar.ID, p2 int) {
	counter := atomic.AddUint64(&m.AddObjectWithTLLPreCounter, 1)
	defer atomic.AddUint64(&m.AddObjectWithTLLCounter, 1)

	if len(m.AddObjectWithTLLMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddObjectWithTLLMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentIndexStorageMock.AddObjectWithTLL. %v %v %v", p, p1, p2)
			return
		}

		input := m.AddObjectWithTLLMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentIndexStorageMockAddObjectWithTLLInput{p, p1, p2}, "RecentIndexStorage.AddObjectWithTLL got unexpected parameters")

		return
	}

	if m.AddObjectWithTLLMock.mainExpectation != nil {

		input := m.AddObjectWithTLLMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentIndexStorageMockAddObjectWithTLLInput{p, p1, p2}, "RecentIndexStorage.AddObjectWithTLL got unexpected parameters")
		}

		return
	}

	if m.AddObjectWithTLLFunc == nil {
		m.t.Fatalf("Unexpected call to RecentIndexStorageMock.AddObjectWithTLL. %v %v %v", p, p1, p2)
		return
	}

	m.AddObjectWithTLLFunc(p, p1, p2)
}

//AddObjectWithTLLMinimockCounter returns a count of RecentIndexStorageMock.AddObjectWithTLLFunc invocations
func (m *RecentIndexStorageMock) AddObjectWithTLLMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectWithTLLCounter)
}

//AddObjectWithTLLMinimockPreCounter returns the value of RecentIndexStorageMock.AddObjectWithTLL invocations
func (m *RecentIndexStorageMock) AddObjectWithTLLMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectWithTLLPreCounter)
}

//AddObjectWithTLLFinished returns true if mock invocations count is ok
func (m *RecentIndexStorageMock) AddObjectWithTLLFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddObjectWithTLLMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddObjectWithTLLCounter) == uint64(len(m.AddObjectWithTLLMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddObjectWithTLLMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddObjectWithTLLCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddObjectWithTLLFunc != nil {
		return atomic.LoadUint64(&m.AddObjectWithTLLCounter) > 0
	}

	return true
}

type mRecentIndexStorageMockDecreaseIndexTTL struct {
	mock              *RecentIndexStorageMock
	mainExpectation   *RecentIndexStorageMockDecreaseIndexTTLExpectation
	expectationSeries []*RecentIndexStorageMockDecreaseIndexTTLExpectation
}

type RecentIndexStorageMockDecreaseIndexTTLExpectation struct {
	input  *RecentIndexStorageMockDecreaseIndexTTLInput
	result *RecentIndexStorageMockDecreaseIndexTTLResult
}

type RecentIndexStorageMockDecreaseIndexTTLInput struct {
	p context.Context
}

type RecentIndexStorageMockDecreaseIndexTTLResult struct {
	r []insolar.ID
}

//Expect specifies that invocation of RecentIndexStorage.DecreaseIndexTTL is expected from 1 to Infinity times
func (m *mRecentIndexStorageMockDecreaseIndexTTL) Expect(p context.Context) *mRecentIndexStorageMockDecreaseIndexTTL {
	m.mock.DecreaseIndexTTLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockDecreaseIndexTTLExpectation{}
	}
	m.mainExpectation.input = &RecentIndexStorageMockDecreaseIndexTTLInput{p}
	return m
}

//Return specifies results of invocation of RecentIndexStorage.DecreaseIndexTTL
func (m *mRecentIndexStorageMockDecreaseIndexTTL) Return(r []insolar.ID) *RecentIndexStorageMock {
	m.mock.DecreaseIndexTTLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockDecreaseIndexTTLExpectation{}
	}
	m.mainExpectation.result = &RecentIndexStorageMockDecreaseIndexTTLResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecentIndexStorage.DecreaseIndexTTL is expected once
func (m *mRecentIndexStorageMockDecreaseIndexTTL) ExpectOnce(p context.Context) *RecentIndexStorageMockDecreaseIndexTTLExpectation {
	m.mock.DecreaseIndexTTLFunc = nil
	m.mainExpectation = nil

	expectation := &RecentIndexStorageMockDecreaseIndexTTLExpectation{}
	expectation.input = &RecentIndexStorageMockDecreaseIndexTTLInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecentIndexStorageMockDecreaseIndexTTLExpectation) Return(r []insolar.ID) {
	e.result = &RecentIndexStorageMockDecreaseIndexTTLResult{r}
}

//Set uses given function f as a mock of RecentIndexStorage.DecreaseIndexTTL method
func (m *mRecentIndexStorageMockDecreaseIndexTTL) Set(f func(p context.Context) (r []insolar.ID)) *RecentIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DecreaseIndexTTLFunc = f
	return m.mock
}

//DecreaseIndexTTL implements github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage interface
func (m *RecentIndexStorageMock) DecreaseIndexTTL(p context.Context) (r []insolar.ID) {
	counter := atomic.AddUint64(&m.DecreaseIndexTTLPreCounter, 1)
	defer atomic.AddUint64(&m.DecreaseIndexTTLCounter, 1)

	if len(m.DecreaseIndexTTLMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DecreaseIndexTTLMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentIndexStorageMock.DecreaseIndexTTL. %v", p)
			return
		}

		input := m.DecreaseIndexTTLMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentIndexStorageMockDecreaseIndexTTLInput{p}, "RecentIndexStorage.DecreaseIndexTTL got unexpected parameters")

		result := m.DecreaseIndexTTLMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecentIndexStorageMock.DecreaseIndexTTL")
			return
		}

		r = result.r

		return
	}

	if m.DecreaseIndexTTLMock.mainExpectation != nil {

		input := m.DecreaseIndexTTLMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentIndexStorageMockDecreaseIndexTTLInput{p}, "RecentIndexStorage.DecreaseIndexTTL got unexpected parameters")
		}

		result := m.DecreaseIndexTTLMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecentIndexStorageMock.DecreaseIndexTTL")
		}

		r = result.r

		return
	}

	if m.DecreaseIndexTTLFunc == nil {
		m.t.Fatalf("Unexpected call to RecentIndexStorageMock.DecreaseIndexTTL. %v", p)
		return
	}

	return m.DecreaseIndexTTLFunc(p)
}

//DecreaseIndexTTLMinimockCounter returns a count of RecentIndexStorageMock.DecreaseIndexTTLFunc invocations
func (m *RecentIndexStorageMock) DecreaseIndexTTLMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DecreaseIndexTTLCounter)
}

//DecreaseIndexTTLMinimockPreCounter returns the value of RecentIndexStorageMock.DecreaseIndexTTL invocations
func (m *RecentIndexStorageMock) DecreaseIndexTTLMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DecreaseIndexTTLPreCounter)
}

//DecreaseIndexTTLFinished returns true if mock invocations count is ok
func (m *RecentIndexStorageMock) DecreaseIndexTTLFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DecreaseIndexTTLMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DecreaseIndexTTLCounter) == uint64(len(m.DecreaseIndexTTLMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DecreaseIndexTTLMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DecreaseIndexTTLCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DecreaseIndexTTLFunc != nil {
		return atomic.LoadUint64(&m.DecreaseIndexTTLCounter) > 0
	}

	return true
}

type mRecentIndexStorageMockFilterNotExistWithLock struct {
	mock              *RecentIndexStorageMock
	mainExpectation   *RecentIndexStorageMockFilterNotExistWithLockExpectation
	expectationSeries []*RecentIndexStorageMockFilterNotExistWithLockExpectation
}

type RecentIndexStorageMockFilterNotExistWithLockExpectation struct {
	input *RecentIndexStorageMockFilterNotExistWithLockInput
}

type RecentIndexStorageMockFilterNotExistWithLockInput struct {
	p  context.Context
	p1 []insolar.ID
	p2 func(p []insolar.ID)
}

//Expect specifies that invocation of RecentIndexStorage.FilterNotExistWithLock is expected from 1 to Infinity times
func (m *mRecentIndexStorageMockFilterNotExistWithLock) Expect(p context.Context, p1 []insolar.ID, p2 func(p []insolar.ID)) *mRecentIndexStorageMockFilterNotExistWithLock {
	m.mock.FilterNotExistWithLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockFilterNotExistWithLockExpectation{}
	}
	m.mainExpectation.input = &RecentIndexStorageMockFilterNotExistWithLockInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RecentIndexStorage.FilterNotExistWithLock
func (m *mRecentIndexStorageMockFilterNotExistWithLock) Return() *RecentIndexStorageMock {
	m.mock.FilterNotExistWithLockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockFilterNotExistWithLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentIndexStorage.FilterNotExistWithLock is expected once
func (m *mRecentIndexStorageMockFilterNotExistWithLock) ExpectOnce(p context.Context, p1 []insolar.ID, p2 func(p []insolar.ID)) *RecentIndexStorageMockFilterNotExistWithLockExpectation {
	m.mock.FilterNotExistWithLockFunc = nil
	m.mainExpectation = nil

	expectation := &RecentIndexStorageMockFilterNotExistWithLockExpectation{}
	expectation.input = &RecentIndexStorageMockFilterNotExistWithLockInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentIndexStorage.FilterNotExistWithLock method
func (m *mRecentIndexStorageMockFilterNotExistWithLock) Set(f func(p context.Context, p1 []insolar.ID, p2 func(p []insolar.ID))) *RecentIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FilterNotExistWithLockFunc = f
	return m.mock
}

//FilterNotExistWithLock implements github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage interface
func (m *RecentIndexStorageMock) FilterNotExistWithLock(p context.Context, p1 []insolar.ID, p2 func(p []insolar.ID)) {
	counter := atomic.AddUint64(&m.FilterNotExistWithLockPreCounter, 1)
	defer atomic.AddUint64(&m.FilterNotExistWithLockCounter, 1)

	if len(m.FilterNotExistWithLockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FilterNotExistWithLockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentIndexStorageMock.FilterNotExistWithLock. %v %v %v", p, p1, p2)
			return
		}

		input := m.FilterNotExistWithLockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentIndexStorageMockFilterNotExistWithLockInput{p, p1, p2}, "RecentIndexStorage.FilterNotExistWithLock got unexpected parameters")

		return
	}

	if m.FilterNotExistWithLockMock.mainExpectation != nil {

		input := m.FilterNotExistWithLockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentIndexStorageMockFilterNotExistWithLockInput{p, p1, p2}, "RecentIndexStorage.FilterNotExistWithLock got unexpected parameters")
		}

		return
	}

	if m.FilterNotExistWithLockFunc == nil {
		m.t.Fatalf("Unexpected call to RecentIndexStorageMock.FilterNotExistWithLock. %v %v %v", p, p1, p2)
		return
	}

	m.FilterNotExistWithLockFunc(p, p1, p2)
}

//FilterNotExistWithLockMinimockCounter returns a count of RecentIndexStorageMock.FilterNotExistWithLockFunc invocations
func (m *RecentIndexStorageMock) FilterNotExistWithLockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FilterNotExistWithLockCounter)
}

//FilterNotExistWithLockMinimockPreCounter returns the value of RecentIndexStorageMock.FilterNotExistWithLock invocations
func (m *RecentIndexStorageMock) FilterNotExistWithLockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FilterNotExistWithLockPreCounter)
}

//FilterNotExistWithLockFinished returns true if mock invocations count is ok
func (m *RecentIndexStorageMock) FilterNotExistWithLockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FilterNotExistWithLockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FilterNotExistWithLockCounter) == uint64(len(m.FilterNotExistWithLockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FilterNotExistWithLockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FilterNotExistWithLockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FilterNotExistWithLockFunc != nil {
		return atomic.LoadUint64(&m.FilterNotExistWithLockCounter) > 0
	}

	return true
}

type mRecentIndexStorageMockGetObjects struct {
	mock              *RecentIndexStorageMock
	mainExpectation   *RecentIndexStorageMockGetObjectsExpectation
	expectationSeries []*RecentIndexStorageMockGetObjectsExpectation
}

type RecentIndexStorageMockGetObjectsExpectation struct {
	result *RecentIndexStorageMockGetObjectsResult
}

type RecentIndexStorageMockGetObjectsResult struct {
	r map[insolar.ID]int
}

//Expect specifies that invocation of RecentIndexStorage.GetObjects is expected from 1 to Infinity times
func (m *mRecentIndexStorageMockGetObjects) Expect() *mRecentIndexStorageMockGetObjects {
	m.mock.GetObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockGetObjectsExpectation{}
	}

	return m
}

//Return specifies results of invocation of RecentIndexStorage.GetObjects
func (m *mRecentIndexStorageMockGetObjects) Return(r map[insolar.ID]int) *RecentIndexStorageMock {
	m.mock.GetObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentIndexStorageMockGetObjectsExpectation{}
	}
	m.mainExpectation.result = &RecentIndexStorageMockGetObjectsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecentIndexStorage.GetObjects is expected once
func (m *mRecentIndexStorageMockGetObjects) ExpectOnce() *RecentIndexStorageMockGetObjectsExpectation {
	m.mock.GetObjectsFunc = nil
	m.mainExpectation = nil

	expectation := &RecentIndexStorageMockGetObjectsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecentIndexStorageMockGetObjectsExpectation) Return(r map[insolar.ID]int) {
	e.result = &RecentIndexStorageMockGetObjectsResult{r}
}

//Set uses given function f as a mock of RecentIndexStorage.GetObjects method
func (m *mRecentIndexStorageMockGetObjects) Set(f func() (r map[insolar.ID]int)) *RecentIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectsFunc = f
	return m.mock
}

//GetObjects implements github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage interface
func (m *RecentIndexStorageMock) GetObjects() (r map[insolar.ID]int) {
	counter := atomic.AddUint64(&m.GetObjectsPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectsCounter, 1)

	if len(m.GetObjectsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentIndexStorageMock.GetObjects.")
			return
		}

		result := m.GetObjectsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecentIndexStorageMock.GetObjects")
			return
		}

		r = result.r

		return
	}

	if m.GetObjectsMock.mainExpectation != nil {

		result := m.GetObjectsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecentIndexStorageMock.GetObjects")
		}

		r = result.r

		return
	}

	if m.GetObjectsFunc == nil {
		m.t.Fatalf("Unexpected call to RecentIndexStorageMock.GetObjects.")
		return
	}

	return m.GetObjectsFunc()
}

//GetObjectsMinimockCounter returns a count of RecentIndexStorageMock.GetObjectsFunc invocations
func (m *RecentIndexStorageMock) GetObjectsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectsCounter)
}

//GetObjectsMinimockPreCounter returns the value of RecentIndexStorageMock.GetObjects invocations
func (m *RecentIndexStorageMock) GetObjectsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectsPreCounter)
}

//GetObjectsFinished returns true if mock invocations count is ok
func (m *RecentIndexStorageMock) GetObjectsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetObjectsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetObjectsCounter) == uint64(len(m.GetObjectsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetObjectsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetObjectsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetObjectsFunc != nil {
		return atomic.LoadUint64(&m.GetObjectsCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecentIndexStorageMock) ValidateCallCounters() {

	if !m.AddObjectFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.AddObject")
	}

	if !m.AddObjectWithTLLFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.AddObjectWithTLL")
	}

	if !m.DecreaseIndexTTLFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.DecreaseIndexTTL")
	}

	if !m.FilterNotExistWithLockFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.FilterNotExistWithLock")
	}

	if !m.GetObjectsFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.GetObjects")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecentIndexStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecentIndexStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecentIndexStorageMock) MinimockFinish() {

	if !m.AddObjectFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.AddObject")
	}

	if !m.AddObjectWithTLLFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.AddObjectWithTLL")
	}

	if !m.DecreaseIndexTTLFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.DecreaseIndexTTL")
	}

	if !m.FilterNotExistWithLockFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.FilterNotExistWithLock")
	}

	if !m.GetObjectsFinished() {
		m.t.Fatal("Expected call to RecentIndexStorageMock.GetObjects")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecentIndexStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RecentIndexStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddObjectFinished()
		ok = ok && m.AddObjectWithTLLFinished()
		ok = ok && m.DecreaseIndexTTLFinished()
		ok = ok && m.FilterNotExistWithLockFinished()
		ok = ok && m.GetObjectsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddObjectFinished() {
				m.t.Error("Expected call to RecentIndexStorageMock.AddObject")
			}

			if !m.AddObjectWithTLLFinished() {
				m.t.Error("Expected call to RecentIndexStorageMock.AddObjectWithTLL")
			}

			if !m.DecreaseIndexTTLFinished() {
				m.t.Error("Expected call to RecentIndexStorageMock.DecreaseIndexTTL")
			}

			if !m.FilterNotExistWithLockFinished() {
				m.t.Error("Expected call to RecentIndexStorageMock.FilterNotExistWithLock")
			}

			if !m.GetObjectsFinished() {
				m.t.Error("Expected call to RecentIndexStorageMock.GetObjects")
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
func (m *RecentIndexStorageMock) AllMocksCalled() bool {

	if !m.AddObjectFinished() {
		return false
	}

	if !m.AddObjectWithTLLFinished() {
		return false
	}

	if !m.DecreaseIndexTTLFinished() {
		return false
	}

	if !m.FilterNotExistWithLockFinished() {
		return false
	}

	if !m.GetObjectsFinished() {
		return false
	}

	return true
}

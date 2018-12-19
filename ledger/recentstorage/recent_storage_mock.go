package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecentStorage" can be found in github.com/insolar/insolar/ledger/recentstorage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//RecentStorageMock implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage
type RecentStorageMock struct {
	t minimock.Tester

	AddObjectFunc       func(p core.RecordID, p1 bool)
	AddObjectCounter    uint64
	AddObjectPreCounter uint64
	AddObjectMock       mRecentStorageMockAddObject

	AddObjectWithTLLFunc       func(p core.RecordID, p1 int, p2 bool)
	AddObjectWithTLLCounter    uint64
	AddObjectWithTLLPreCounter uint64
	AddObjectWithTLLMock       mRecentStorageMockAddObjectWithTLL

	AddPendingRequestFunc       func(p core.RecordID)
	AddPendingRequestCounter    uint64
	AddPendingRequestPreCounter uint64
	AddPendingRequestMock       mRecentStorageMockAddPendingRequest

	ClearObjectsFunc       func()
	ClearObjectsCounter    uint64
	ClearObjectsPreCounter uint64
	ClearObjectsMock       mRecentStorageMockClearObjects

	ClearZeroTTLObjectsFunc       func()
	ClearZeroTTLObjectsCounter    uint64
	ClearZeroTTLObjectsPreCounter uint64
	ClearZeroTTLObjectsMock       mRecentStorageMockClearZeroTTLObjects

	GetObjectsFunc       func() (r map[core.RecordID]int)
	GetObjectsCounter    uint64
	GetObjectsPreCounter uint64
	GetObjectsMock       mRecentStorageMockGetObjects

	GetRequestsFunc       func() (r []core.RecordID)
	GetRequestsCounter    uint64
	GetRequestsPreCounter uint64
	GetRequestsMock       mRecentStorageMockGetRequests

	IsMineFunc       func(p core.RecordID) (r bool)
	IsMineCounter    uint64
	IsMinePreCounter uint64
	IsMineMock       mRecentStorageMockIsMine

	RemovePendingRequestFunc       func(p core.RecordID)
	RemovePendingRequestCounter    uint64
	RemovePendingRequestPreCounter uint64
	RemovePendingRequestMock       mRecentStorageMockRemovePendingRequest
}

//NewRecentStorageMock returns a mock for github.com/insolar/insolar/ledger/recentstorage.RecentStorage
func NewRecentStorageMock(t minimock.Tester) *RecentStorageMock {
	m := &RecentStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddObjectMock = mRecentStorageMockAddObject{mock: m}
	m.AddObjectWithTLLMock = mRecentStorageMockAddObjectWithTLL{mock: m}
	m.AddPendingRequestMock = mRecentStorageMockAddPendingRequest{mock: m}
	m.ClearObjectsMock = mRecentStorageMockClearObjects{mock: m}
	m.ClearZeroTTLObjectsMock = mRecentStorageMockClearZeroTTLObjects{mock: m}
	m.GetObjectsMock = mRecentStorageMockGetObjects{mock: m}
	m.GetRequestsMock = mRecentStorageMockGetRequests{mock: m}
	m.IsMineMock = mRecentStorageMockIsMine{mock: m}
	m.RemovePendingRequestMock = mRecentStorageMockRemovePendingRequest{mock: m}

	return m
}

type mRecentStorageMockAddObject struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockAddObjectExpectation
	expectationSeries []*RecentStorageMockAddObjectExpectation
}

type RecentStorageMockAddObjectExpectation struct {
	input *RecentStorageMockAddObjectInput
}

type RecentStorageMockAddObjectInput struct {
	p  core.RecordID
	p1 bool
}

//Expect specifies that invocation of RecentStorage.AddObject is expected from 1 to Infinity times
func (m *mRecentStorageMockAddObject) Expect(p core.RecordID, p1 bool) *mRecentStorageMockAddObject {
	m.mock.AddObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockAddObjectExpectation{}
	}
	m.mainExpectation.input = &RecentStorageMockAddObjectInput{p, p1}
	return m
}

//Return specifies results of invocation of RecentStorage.AddObject
func (m *mRecentStorageMockAddObject) Return() *RecentStorageMock {
	m.mock.AddObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockAddObjectExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.AddObject is expected once
func (m *mRecentStorageMockAddObject) ExpectOnce(p core.RecordID, p1 bool) *RecentStorageMockAddObjectExpectation {
	m.mock.AddObjectFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockAddObjectExpectation{}
	expectation.input = &RecentStorageMockAddObjectInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentStorage.AddObject method
func (m *mRecentStorageMockAddObject) Set(f func(p core.RecordID, p1 bool)) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddObjectFunc = f
	return m.mock
}

//AddObject implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) AddObject(p core.RecordID, p1 bool) {
	counter := atomic.AddUint64(&m.AddObjectPreCounter, 1)
	defer atomic.AddUint64(&m.AddObjectCounter, 1)

	if len(m.AddObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.AddObject. %v %v", p, p1)
			return
		}

		input := m.AddObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentStorageMockAddObjectInput{p, p1}, "RecentStorage.AddObject got unexpected parameters")

		return
	}

	if m.AddObjectMock.mainExpectation != nil {

		input := m.AddObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentStorageMockAddObjectInput{p, p1}, "RecentStorage.AddObject got unexpected parameters")
		}

		return
	}

	if m.AddObjectFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.AddObject. %v %v", p, p1)
		return
	}

	m.AddObjectFunc(p, p1)
}

//AddObjectMinimockCounter returns a count of RecentStorageMock.AddObjectFunc invocations
func (m *RecentStorageMock) AddObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectCounter)
}

//AddObjectMinimockPreCounter returns the value of RecentStorageMock.AddObject invocations
func (m *RecentStorageMock) AddObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectPreCounter)
}

//AddObjectFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) AddObjectFinished() bool {
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

type mRecentStorageMockAddObjectWithTLL struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockAddObjectWithTLLExpectation
	expectationSeries []*RecentStorageMockAddObjectWithTLLExpectation
}

type RecentStorageMockAddObjectWithTLLExpectation struct {
	input *RecentStorageMockAddObjectWithTLLInput
}

type RecentStorageMockAddObjectWithTLLInput struct {
	p  core.RecordID
	p1 int
	p2 bool
}

//Expect specifies that invocation of RecentStorage.AddObjectWithTLL is expected from 1 to Infinity times
func (m *mRecentStorageMockAddObjectWithTLL) Expect(p core.RecordID, p1 int, p2 bool) *mRecentStorageMockAddObjectWithTLL {
	m.mock.AddObjectWithTLLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockAddObjectWithTLLExpectation{}
	}
	m.mainExpectation.input = &RecentStorageMockAddObjectWithTLLInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RecentStorage.AddObjectWithTLL
func (m *mRecentStorageMockAddObjectWithTLL) Return() *RecentStorageMock {
	m.mock.AddObjectWithTLLFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockAddObjectWithTLLExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.AddObjectWithTLL is expected once
func (m *mRecentStorageMockAddObjectWithTLL) ExpectOnce(p core.RecordID, p1 int, p2 bool) *RecentStorageMockAddObjectWithTLLExpectation {
	m.mock.AddObjectWithTLLFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockAddObjectWithTLLExpectation{}
	expectation.input = &RecentStorageMockAddObjectWithTLLInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentStorage.AddObjectWithTLL method
func (m *mRecentStorageMockAddObjectWithTLL) Set(f func(p core.RecordID, p1 int, p2 bool)) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddObjectWithTLLFunc = f
	return m.mock
}

//AddObjectWithTLL implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) AddObjectWithTLL(p core.RecordID, p1 int, p2 bool) {
	counter := atomic.AddUint64(&m.AddObjectWithTLLPreCounter, 1)
	defer atomic.AddUint64(&m.AddObjectWithTLLCounter, 1)

	if len(m.AddObjectWithTLLMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddObjectWithTLLMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.AddObjectWithTLL. %v %v %v", p, p1, p2)
			return
		}

		input := m.AddObjectWithTLLMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentStorageMockAddObjectWithTLLInput{p, p1, p2}, "RecentStorage.AddObjectWithTLL got unexpected parameters")

		return
	}

	if m.AddObjectWithTLLMock.mainExpectation != nil {

		input := m.AddObjectWithTLLMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentStorageMockAddObjectWithTLLInput{p, p1, p2}, "RecentStorage.AddObjectWithTLL got unexpected parameters")
		}

		return
	}

	if m.AddObjectWithTLLFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.AddObjectWithTLL. %v %v %v", p, p1, p2)
		return
	}

	m.AddObjectWithTLLFunc(p, p1, p2)
}

//AddObjectWithTLLMinimockCounter returns a count of RecentStorageMock.AddObjectWithTLLFunc invocations
func (m *RecentStorageMock) AddObjectWithTLLMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectWithTLLCounter)
}

//AddObjectWithTLLMinimockPreCounter returns the value of RecentStorageMock.AddObjectWithTLL invocations
func (m *RecentStorageMock) AddObjectWithTLLMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddObjectWithTLLPreCounter)
}

//AddObjectWithTLLFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) AddObjectWithTLLFinished() bool {
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

type mRecentStorageMockAddPendingRequest struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockAddPendingRequestExpectation
	expectationSeries []*RecentStorageMockAddPendingRequestExpectation
}

type RecentStorageMockAddPendingRequestExpectation struct {
	input *RecentStorageMockAddPendingRequestInput
}

type RecentStorageMockAddPendingRequestInput struct {
	p core.RecordID
}

//Expect specifies that invocation of RecentStorage.AddPendingRequest is expected from 1 to Infinity times
func (m *mRecentStorageMockAddPendingRequest) Expect(p core.RecordID) *mRecentStorageMockAddPendingRequest {
	m.mock.AddPendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockAddPendingRequestExpectation{}
	}
	m.mainExpectation.input = &RecentStorageMockAddPendingRequestInput{p}
	return m
}

//Return specifies results of invocation of RecentStorage.AddPendingRequest
func (m *mRecentStorageMockAddPendingRequest) Return() *RecentStorageMock {
	m.mock.AddPendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockAddPendingRequestExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.AddPendingRequest is expected once
func (m *mRecentStorageMockAddPendingRequest) ExpectOnce(p core.RecordID) *RecentStorageMockAddPendingRequestExpectation {
	m.mock.AddPendingRequestFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockAddPendingRequestExpectation{}
	expectation.input = &RecentStorageMockAddPendingRequestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentStorage.AddPendingRequest method
func (m *mRecentStorageMockAddPendingRequest) Set(f func(p core.RecordID)) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddPendingRequestFunc = f
	return m.mock
}

//AddPendingRequest implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) AddPendingRequest(p core.RecordID) {
	counter := atomic.AddUint64(&m.AddPendingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.AddPendingRequestCounter, 1)

	if len(m.AddPendingRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddPendingRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.AddPendingRequest. %v", p)
			return
		}

		input := m.AddPendingRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentStorageMockAddPendingRequestInput{p}, "RecentStorage.AddPendingRequest got unexpected parameters")

		return
	}

	if m.AddPendingRequestMock.mainExpectation != nil {

		input := m.AddPendingRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentStorageMockAddPendingRequestInput{p}, "RecentStorage.AddPendingRequest got unexpected parameters")
		}

		return
	}

	if m.AddPendingRequestFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.AddPendingRequest. %v", p)
		return
	}

	m.AddPendingRequestFunc(p)
}

//AddPendingRequestMinimockCounter returns a count of RecentStorageMock.AddPendingRequestFunc invocations
func (m *RecentStorageMock) AddPendingRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddPendingRequestCounter)
}

//AddPendingRequestMinimockPreCounter returns the value of RecentStorageMock.AddPendingRequest invocations
func (m *RecentStorageMock) AddPendingRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddPendingRequestPreCounter)
}

//AddPendingRequestFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) AddPendingRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddPendingRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddPendingRequestCounter) == uint64(len(m.AddPendingRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddPendingRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddPendingRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddPendingRequestFunc != nil {
		return atomic.LoadUint64(&m.AddPendingRequestCounter) > 0
	}

	return true
}

type mRecentStorageMockClearObjects struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockClearObjectsExpectation
	expectationSeries []*RecentStorageMockClearObjectsExpectation
}

type RecentStorageMockClearObjectsExpectation struct {
}

//Expect specifies that invocation of RecentStorage.ClearObjects is expected from 1 to Infinity times
func (m *mRecentStorageMockClearObjects) Expect() *mRecentStorageMockClearObjects {
	m.mock.ClearObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockClearObjectsExpectation{}
	}

	return m
}

//Return specifies results of invocation of RecentStorage.ClearObjects
func (m *mRecentStorageMockClearObjects) Return() *RecentStorageMock {
	m.mock.ClearObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockClearObjectsExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.ClearObjects is expected once
func (m *mRecentStorageMockClearObjects) ExpectOnce() *RecentStorageMockClearObjectsExpectation {
	m.mock.ClearObjectsFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockClearObjectsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentStorage.ClearObjects method
func (m *mRecentStorageMockClearObjects) Set(f func()) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ClearObjectsFunc = f
	return m.mock
}

//ClearObjects implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) ClearObjects() {
	counter := atomic.AddUint64(&m.ClearObjectsPreCounter, 1)
	defer atomic.AddUint64(&m.ClearObjectsCounter, 1)

	if len(m.ClearObjectsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ClearObjectsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.ClearObjects.")
			return
		}

		return
	}

	if m.ClearObjectsMock.mainExpectation != nil {

		return
	}

	if m.ClearObjectsFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.ClearObjects.")
		return
	}

	m.ClearObjectsFunc()
}

//ClearObjectsMinimockCounter returns a count of RecentStorageMock.ClearObjectsFunc invocations
func (m *RecentStorageMock) ClearObjectsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ClearObjectsCounter)
}

//ClearObjectsMinimockPreCounter returns the value of RecentStorageMock.ClearObjects invocations
func (m *RecentStorageMock) ClearObjectsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClearObjectsPreCounter)
}

//ClearObjectsFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) ClearObjectsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ClearObjectsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ClearObjectsCounter) == uint64(len(m.ClearObjectsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ClearObjectsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ClearObjectsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ClearObjectsFunc != nil {
		return atomic.LoadUint64(&m.ClearObjectsCounter) > 0
	}

	return true
}

type mRecentStorageMockClearZeroTTLObjects struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockClearZeroTTLObjectsExpectation
	expectationSeries []*RecentStorageMockClearZeroTTLObjectsExpectation
}

type RecentStorageMockClearZeroTTLObjectsExpectation struct {
}

//Expect specifies that invocation of RecentStorage.ClearZeroTTLObjects is expected from 1 to Infinity times
func (m *mRecentStorageMockClearZeroTTLObjects) Expect() *mRecentStorageMockClearZeroTTLObjects {
	m.mock.ClearZeroTTLObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockClearZeroTTLObjectsExpectation{}
	}

	return m
}

//Return specifies results of invocation of RecentStorage.ClearZeroTTLObjects
func (m *mRecentStorageMockClearZeroTTLObjects) Return() *RecentStorageMock {
	m.mock.ClearZeroTTLObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockClearZeroTTLObjectsExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.ClearZeroTTLObjects is expected once
func (m *mRecentStorageMockClearZeroTTLObjects) ExpectOnce() *RecentStorageMockClearZeroTTLObjectsExpectation {
	m.mock.ClearZeroTTLObjectsFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockClearZeroTTLObjectsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentStorage.ClearZeroTTLObjects method
func (m *mRecentStorageMockClearZeroTTLObjects) Set(f func()) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ClearZeroTTLObjectsFunc = f
	return m.mock
}

//ClearZeroTTLObjects implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) ClearZeroTTLObjects() {
	counter := atomic.AddUint64(&m.ClearZeroTTLObjectsPreCounter, 1)
	defer atomic.AddUint64(&m.ClearZeroTTLObjectsCounter, 1)

	if len(m.ClearZeroTTLObjectsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ClearZeroTTLObjectsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.ClearZeroTTLObjects.")
			return
		}

		return
	}

	if m.ClearZeroTTLObjectsMock.mainExpectation != nil {

		return
	}

	if m.ClearZeroTTLObjectsFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.ClearZeroTTLObjects.")
		return
	}

	m.ClearZeroTTLObjectsFunc()
}

//ClearZeroTTLObjectsMinimockCounter returns a count of RecentStorageMock.ClearZeroTTLObjectsFunc invocations
func (m *RecentStorageMock) ClearZeroTTLObjectsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter)
}

//ClearZeroTTLObjectsMinimockPreCounter returns the value of RecentStorageMock.ClearZeroTTLObjects invocations
func (m *RecentStorageMock) ClearZeroTTLObjectsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClearZeroTTLObjectsPreCounter)
}

//ClearZeroTTLObjectsFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) ClearZeroTTLObjectsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ClearZeroTTLObjectsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) == uint64(len(m.ClearZeroTTLObjectsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ClearZeroTTLObjectsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ClearZeroTTLObjectsFunc != nil {
		return atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) > 0
	}

	return true
}

type mRecentStorageMockGetObjects struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockGetObjectsExpectation
	expectationSeries []*RecentStorageMockGetObjectsExpectation
}

type RecentStorageMockGetObjectsExpectation struct {
	result *RecentStorageMockGetObjectsResult
}

type RecentStorageMockGetObjectsResult struct {
	r map[core.RecordID]int
}

//Expect specifies that invocation of RecentStorage.GetObjects is expected from 1 to Infinity times
func (m *mRecentStorageMockGetObjects) Expect() *mRecentStorageMockGetObjects {
	m.mock.GetObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockGetObjectsExpectation{}
	}

	return m
}

//Return specifies results of invocation of RecentStorage.GetObjects
func (m *mRecentStorageMockGetObjects) Return(r map[core.RecordID]int) *RecentStorageMock {
	m.mock.GetObjectsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockGetObjectsExpectation{}
	}
	m.mainExpectation.result = &RecentStorageMockGetObjectsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.GetObjects is expected once
func (m *mRecentStorageMockGetObjects) ExpectOnce() *RecentStorageMockGetObjectsExpectation {
	m.mock.GetObjectsFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockGetObjectsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecentStorageMockGetObjectsExpectation) Return(r map[core.RecordID]int) {
	e.result = &RecentStorageMockGetObjectsResult{r}
}

//Set uses given function f as a mock of RecentStorage.GetObjects method
func (m *mRecentStorageMockGetObjects) Set(f func() (r map[core.RecordID]int)) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectsFunc = f
	return m.mock
}

//GetObjects implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) GetObjects() (r map[core.RecordID]int) {
	counter := atomic.AddUint64(&m.GetObjectsPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectsCounter, 1)

	if len(m.GetObjectsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.GetObjects.")
			return
		}

		result := m.GetObjectsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecentStorageMock.GetObjects")
			return
		}

		r = result.r

		return
	}

	if m.GetObjectsMock.mainExpectation != nil {

		result := m.GetObjectsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecentStorageMock.GetObjects")
		}

		r = result.r

		return
	}

	if m.GetObjectsFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.GetObjects.")
		return
	}

	return m.GetObjectsFunc()
}

//GetObjectsMinimockCounter returns a count of RecentStorageMock.GetObjectsFunc invocations
func (m *RecentStorageMock) GetObjectsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectsCounter)
}

//GetObjectsMinimockPreCounter returns the value of RecentStorageMock.GetObjects invocations
func (m *RecentStorageMock) GetObjectsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectsPreCounter)
}

//GetObjectsFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) GetObjectsFinished() bool {
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

type mRecentStorageMockGetRequests struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockGetRequestsExpectation
	expectationSeries []*RecentStorageMockGetRequestsExpectation
}

type RecentStorageMockGetRequestsExpectation struct {
	result *RecentStorageMockGetRequestsResult
}

type RecentStorageMockGetRequestsResult struct {
	r []core.RecordID
}

//Expect specifies that invocation of RecentStorage.GetRequests is expected from 1 to Infinity times
func (m *mRecentStorageMockGetRequests) Expect() *mRecentStorageMockGetRequests {
	m.mock.GetRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockGetRequestsExpectation{}
	}

	return m
}

//Return specifies results of invocation of RecentStorage.GetRequests
func (m *mRecentStorageMockGetRequests) Return(r []core.RecordID) *RecentStorageMock {
	m.mock.GetRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockGetRequestsExpectation{}
	}
	m.mainExpectation.result = &RecentStorageMockGetRequestsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.GetRequests is expected once
func (m *mRecentStorageMockGetRequests) ExpectOnce() *RecentStorageMockGetRequestsExpectation {
	m.mock.GetRequestsFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockGetRequestsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecentStorageMockGetRequestsExpectation) Return(r []core.RecordID) {
	e.result = &RecentStorageMockGetRequestsResult{r}
}

//Set uses given function f as a mock of RecentStorage.GetRequests method
func (m *mRecentStorageMockGetRequests) Set(f func() (r []core.RecordID)) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRequestsFunc = f
	return m.mock
}

//GetRequests implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) GetRequests() (r []core.RecordID) {
	counter := atomic.AddUint64(&m.GetRequestsPreCounter, 1)
	defer atomic.AddUint64(&m.GetRequestsCounter, 1)

	if len(m.GetRequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.GetRequests.")
			return
		}

		result := m.GetRequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecentStorageMock.GetRequests")
			return
		}

		r = result.r

		return
	}

	if m.GetRequestsMock.mainExpectation != nil {

		result := m.GetRequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecentStorageMock.GetRequests")
		}

		r = result.r

		return
	}

	if m.GetRequestsFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.GetRequests.")
		return
	}

	return m.GetRequestsFunc()
}

//GetRequestsMinimockCounter returns a count of RecentStorageMock.GetRequestsFunc invocations
func (m *RecentStorageMock) GetRequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRequestsCounter)
}

//GetRequestsMinimockPreCounter returns the value of RecentStorageMock.GetRequests invocations
func (m *RecentStorageMock) GetRequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRequestsPreCounter)
}

//GetRequestsFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) GetRequestsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRequestsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRequestsCounter) == uint64(len(m.GetRequestsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRequestsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRequestsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRequestsFunc != nil {
		return atomic.LoadUint64(&m.GetRequestsCounter) > 0
	}

	return true
}

type mRecentStorageMockIsMine struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockIsMineExpectation
	expectationSeries []*RecentStorageMockIsMineExpectation
}

type RecentStorageMockIsMineExpectation struct {
	input  *RecentStorageMockIsMineInput
	result *RecentStorageMockIsMineResult
}

type RecentStorageMockIsMineInput struct {
	p core.RecordID
}

type RecentStorageMockIsMineResult struct {
	r bool
}

//Expect specifies that invocation of RecentStorage.IsMine is expected from 1 to Infinity times
func (m *mRecentStorageMockIsMine) Expect(p core.RecordID) *mRecentStorageMockIsMine {
	m.mock.IsMineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockIsMineExpectation{}
	}
	m.mainExpectation.input = &RecentStorageMockIsMineInput{p}
	return m
}

//Return specifies results of invocation of RecentStorage.IsMine
func (m *mRecentStorageMockIsMine) Return(r bool) *RecentStorageMock {
	m.mock.IsMineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockIsMineExpectation{}
	}
	m.mainExpectation.result = &RecentStorageMockIsMineResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.IsMine is expected once
func (m *mRecentStorageMockIsMine) ExpectOnce(p core.RecordID) *RecentStorageMockIsMineExpectation {
	m.mock.IsMineFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockIsMineExpectation{}
	expectation.input = &RecentStorageMockIsMineInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecentStorageMockIsMineExpectation) Return(r bool) {
	e.result = &RecentStorageMockIsMineResult{r}
}

//Set uses given function f as a mock of RecentStorage.IsMine method
func (m *mRecentStorageMockIsMine) Set(f func(p core.RecordID) (r bool)) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsMineFunc = f
	return m.mock
}

//IsMine implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) IsMine(p core.RecordID) (r bool) {
	counter := atomic.AddUint64(&m.IsMinePreCounter, 1)
	defer atomic.AddUint64(&m.IsMineCounter, 1)

	if len(m.IsMineMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsMineMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.IsMine. %v", p)
			return
		}

		input := m.IsMineMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentStorageMockIsMineInput{p}, "RecentStorage.IsMine got unexpected parameters")

		result := m.IsMineMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecentStorageMock.IsMine")
			return
		}

		r = result.r

		return
	}

	if m.IsMineMock.mainExpectation != nil {

		input := m.IsMineMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentStorageMockIsMineInput{p}, "RecentStorage.IsMine got unexpected parameters")
		}

		result := m.IsMineMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecentStorageMock.IsMine")
		}

		r = result.r

		return
	}

	if m.IsMineFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.IsMine. %v", p)
		return
	}

	return m.IsMineFunc(p)
}

//IsMineMinimockCounter returns a count of RecentStorageMock.IsMineFunc invocations
func (m *RecentStorageMock) IsMineMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsMineCounter)
}

//IsMineMinimockPreCounter returns the value of RecentStorageMock.IsMine invocations
func (m *RecentStorageMock) IsMineMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsMinePreCounter)
}

//IsMineFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) IsMineFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsMineMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsMineCounter) == uint64(len(m.IsMineMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsMineMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsMineCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsMineFunc != nil {
		return atomic.LoadUint64(&m.IsMineCounter) > 0
	}

	return true
}

type mRecentStorageMockRemovePendingRequest struct {
	mock              *RecentStorageMock
	mainExpectation   *RecentStorageMockRemovePendingRequestExpectation
	expectationSeries []*RecentStorageMockRemovePendingRequestExpectation
}

type RecentStorageMockRemovePendingRequestExpectation struct {
	input *RecentStorageMockRemovePendingRequestInput
}

type RecentStorageMockRemovePendingRequestInput struct {
	p core.RecordID
}

//Expect specifies that invocation of RecentStorage.RemovePendingRequest is expected from 1 to Infinity times
func (m *mRecentStorageMockRemovePendingRequest) Expect(p core.RecordID) *mRecentStorageMockRemovePendingRequest {
	m.mock.RemovePendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockRemovePendingRequestExpectation{}
	}
	m.mainExpectation.input = &RecentStorageMockRemovePendingRequestInput{p}
	return m
}

//Return specifies results of invocation of RecentStorage.RemovePendingRequest
func (m *mRecentStorageMockRemovePendingRequest) Return() *RecentStorageMock {
	m.mock.RemovePendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecentStorageMockRemovePendingRequestExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecentStorage.RemovePendingRequest is expected once
func (m *mRecentStorageMockRemovePendingRequest) ExpectOnce(p core.RecordID) *RecentStorageMockRemovePendingRequestExpectation {
	m.mock.RemovePendingRequestFunc = nil
	m.mainExpectation = nil

	expectation := &RecentStorageMockRemovePendingRequestExpectation{}
	expectation.input = &RecentStorageMockRemovePendingRequestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecentStorage.RemovePendingRequest method
func (m *mRecentStorageMockRemovePendingRequest) Set(f func(p core.RecordID)) *RecentStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemovePendingRequestFunc = f
	return m.mock
}

//RemovePendingRequest implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) RemovePendingRequest(p core.RecordID) {
	counter := atomic.AddUint64(&m.RemovePendingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RemovePendingRequestCounter, 1)

	if len(m.RemovePendingRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemovePendingRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecentStorageMock.RemovePendingRequest. %v", p)
			return
		}

		input := m.RemovePendingRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecentStorageMockRemovePendingRequestInput{p}, "RecentStorage.RemovePendingRequest got unexpected parameters")

		return
	}

	if m.RemovePendingRequestMock.mainExpectation != nil {

		input := m.RemovePendingRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecentStorageMockRemovePendingRequestInput{p}, "RecentStorage.RemovePendingRequest got unexpected parameters")
		}

		return
	}

	if m.RemovePendingRequestFunc == nil {
		m.t.Fatalf("Unexpected call to RecentStorageMock.RemovePendingRequest. %v", p)
		return
	}

	m.RemovePendingRequestFunc(p)
}

//RemovePendingRequestMinimockCounter returns a count of RecentStorageMock.RemovePendingRequestFunc invocations
func (m *RecentStorageMock) RemovePendingRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingRequestCounter)
}

//RemovePendingRequestMinimockPreCounter returns the value of RecentStorageMock.RemovePendingRequest invocations
func (m *RecentStorageMock) RemovePendingRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingRequestPreCounter)
}

//RemovePendingRequestFinished returns true if mock invocations count is ok
func (m *RecentStorageMock) RemovePendingRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemovePendingRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemovePendingRequestCounter) == uint64(len(m.RemovePendingRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemovePendingRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemovePendingRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemovePendingRequestFunc != nil {
		return atomic.LoadUint64(&m.RemovePendingRequestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecentStorageMock) ValidateCallCounters() {

	if !m.AddObjectFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.AddObject")
	}

	if !m.AddObjectWithTLLFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.AddObjectWithTLL")
	}

	if !m.AddPendingRequestFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.AddPendingRequest")
	}

	if !m.ClearObjectsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.ClearObjects")
	}

	if !m.ClearZeroTTLObjectsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.ClearZeroTTLObjects")
	}

	if !m.GetObjectsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.GetObjects")
	}

	if !m.GetRequestsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.GetRequests")
	}

	if !m.IsMineFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.IsMine")
	}

	if !m.RemovePendingRequestFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.RemovePendingRequest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecentStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecentStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecentStorageMock) MinimockFinish() {

	if !m.AddObjectFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.AddObject")
	}

	if !m.AddObjectWithTLLFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.AddObjectWithTLL")
	}

	if !m.AddPendingRequestFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.AddPendingRequest")
	}

	if !m.ClearObjectsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.ClearObjects")
	}

	if !m.ClearZeroTTLObjectsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.ClearZeroTTLObjects")
	}

	if !m.GetObjectsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.GetObjects")
	}

	if !m.GetRequestsFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.GetRequests")
	}

	if !m.IsMineFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.IsMine")
	}

	if !m.RemovePendingRequestFinished() {
		m.t.Fatal("Expected call to RecentStorageMock.RemovePendingRequest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecentStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RecentStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddObjectFinished()
		ok = ok && m.AddObjectWithTLLFinished()
		ok = ok && m.AddPendingRequestFinished()
		ok = ok && m.ClearObjectsFinished()
		ok = ok && m.ClearZeroTTLObjectsFinished()
		ok = ok && m.GetObjectsFinished()
		ok = ok && m.GetRequestsFinished()
		ok = ok && m.IsMineFinished()
		ok = ok && m.RemovePendingRequestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddObjectFinished() {
				m.t.Error("Expected call to RecentStorageMock.AddObject")
			}

			if !m.AddObjectWithTLLFinished() {
				m.t.Error("Expected call to RecentStorageMock.AddObjectWithTLL")
			}

			if !m.AddPendingRequestFinished() {
				m.t.Error("Expected call to RecentStorageMock.AddPendingRequest")
			}

			if !m.ClearObjectsFinished() {
				m.t.Error("Expected call to RecentStorageMock.ClearObjects")
			}

			if !m.ClearZeroTTLObjectsFinished() {
				m.t.Error("Expected call to RecentStorageMock.ClearZeroTTLObjects")
			}

			if !m.GetObjectsFinished() {
				m.t.Error("Expected call to RecentStorageMock.GetObjects")
			}

			if !m.GetRequestsFinished() {
				m.t.Error("Expected call to RecentStorageMock.GetRequests")
			}

			if !m.IsMineFinished() {
				m.t.Error("Expected call to RecentStorageMock.IsMine")
			}

			if !m.RemovePendingRequestFinished() {
				m.t.Error("Expected call to RecentStorageMock.RemovePendingRequest")
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
func (m *RecentStorageMock) AllMocksCalled() bool {

	if !m.AddObjectFinished() {
		return false
	}

	if !m.AddObjectWithTLLFinished() {
		return false
	}

	if !m.AddPendingRequestFinished() {
		return false
	}

	if !m.ClearObjectsFinished() {
		return false
	}

	if !m.ClearZeroTTLObjectsFinished() {
		return false
	}

	if !m.GetObjectsFinished() {
		return false
	}

	if !m.GetRequestsFinished() {
		return false
	}

	if !m.IsMineFinished() {
		return false
	}

	if !m.RemovePendingRequestFinished() {
		return false
	}

	return true
}

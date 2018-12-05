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
	mock             *RecentStorageMock
	mockExpectations *RecentStorageMockAddObjectParams
}

//RecentStorageMockAddObjectParams represents input parameters of the RecentStorage.AddObject
type RecentStorageMockAddObjectParams struct {
	p  core.RecordID
	p1 bool
}

//Expect sets up expected params for the RecentStorage.AddObject
func (m *mRecentStorageMockAddObject) Expect(p core.RecordID, p1 bool) *mRecentStorageMockAddObject {
	m.mockExpectations = &RecentStorageMockAddObjectParams{p, p1}
	return m
}

//Return sets up a mock for RecentStorage.AddObject to return Return's arguments
func (m *mRecentStorageMockAddObject) Return() *RecentStorageMock {
	m.mock.AddObjectFunc = func(p core.RecordID, p1 bool) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.AddObject method
func (m *mRecentStorageMockAddObject) Set(f func(p core.RecordID, p1 bool)) *RecentStorageMock {
	m.mock.AddObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//AddObject implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) AddObject(p core.RecordID, p1 bool) {
	atomic.AddUint64(&m.AddObjectPreCounter, 1)
	defer atomic.AddUint64(&m.AddObjectCounter, 1)

	if m.AddObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.AddObjectMock.mockExpectations, RecentStorageMockAddObjectParams{p, p1},
			"RecentStorage.AddObject got unexpected parameters")

		if m.AddObjectFunc == nil {

			m.t.Fatal("No results are set for the RecentStorageMock.AddObject")

			return
		}
	}

	if m.AddObjectFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.AddObject")
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

type mRecentStorageMockAddObjectWithTLL struct {
	mock             *RecentStorageMock
	mockExpectations *RecentStorageMockAddObjectWithTLLParams
}

//RecentStorageMockAddObjectWithTLLParams represents input parameters of the RecentStorage.AddObjectWithTLL
type RecentStorageMockAddObjectWithTLLParams struct {
	p  core.RecordID
	p1 int
	p2 bool
}

//Expect sets up expected params for the RecentStorage.AddObjectWithTLL
func (m *mRecentStorageMockAddObjectWithTLL) Expect(p core.RecordID, p1 int, p2 bool) *mRecentStorageMockAddObjectWithTLL {
	m.mockExpectations = &RecentStorageMockAddObjectWithTLLParams{p, p1, p2}
	return m
}

//Return sets up a mock for RecentStorage.AddObjectWithTLL to return Return's arguments
func (m *mRecentStorageMockAddObjectWithTLL) Return() *RecentStorageMock {
	m.mock.AddObjectWithTLLFunc = func(p core.RecordID, p1 int, p2 bool) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.AddObjectWithTLL method
func (m *mRecentStorageMockAddObjectWithTLL) Set(f func(p core.RecordID, p1 int, p2 bool)) *RecentStorageMock {
	m.mock.AddObjectWithTLLFunc = f
	m.mockExpectations = nil
	return m.mock
}

//AddObjectWithTLL implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) AddObjectWithTLL(p core.RecordID, p1 int, p2 bool) {
	atomic.AddUint64(&m.AddObjectWithTLLPreCounter, 1)
	defer atomic.AddUint64(&m.AddObjectWithTLLCounter, 1)

	if m.AddObjectWithTLLMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.AddObjectWithTLLMock.mockExpectations, RecentStorageMockAddObjectWithTLLParams{p, p1, p2},
			"RecentStorage.AddObjectWithTLL got unexpected parameters")

		if m.AddObjectWithTLLFunc == nil {

			m.t.Fatal("No results are set for the RecentStorageMock.AddObjectWithTLL")

			return
		}
	}

	if m.AddObjectWithTLLFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.AddObjectWithTLL")
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

type mRecentStorageMockAddPendingRequest struct {
	mock             *RecentStorageMock
	mockExpectations *RecentStorageMockAddPendingRequestParams
}

//RecentStorageMockAddPendingRequestParams represents input parameters of the RecentStorage.AddPendingRequest
type RecentStorageMockAddPendingRequestParams struct {
	p core.RecordID
}

//Expect sets up expected params for the RecentStorage.AddPendingRequest
func (m *mRecentStorageMockAddPendingRequest) Expect(p core.RecordID) *mRecentStorageMockAddPendingRequest {
	m.mockExpectations = &RecentStorageMockAddPendingRequestParams{p}
	return m
}

//Return sets up a mock for RecentStorage.AddPendingRequest to return Return's arguments
func (m *mRecentStorageMockAddPendingRequest) Return() *RecentStorageMock {
	m.mock.AddPendingRequestFunc = func(p core.RecordID) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.AddPendingRequest method
func (m *mRecentStorageMockAddPendingRequest) Set(f func(p core.RecordID)) *RecentStorageMock {
	m.mock.AddPendingRequestFunc = f
	m.mockExpectations = nil
	return m.mock
}

//AddPendingRequest implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) AddPendingRequest(p core.RecordID) {
	atomic.AddUint64(&m.AddPendingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.AddPendingRequestCounter, 1)

	if m.AddPendingRequestMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.AddPendingRequestMock.mockExpectations, RecentStorageMockAddPendingRequestParams{p},
			"RecentStorage.AddPendingRequest got unexpected parameters")

		if m.AddPendingRequestFunc == nil {

			m.t.Fatal("No results are set for the RecentStorageMock.AddPendingRequest")

			return
		}
	}

	if m.AddPendingRequestFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.AddPendingRequest")
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

type mRecentStorageMockClearObjects struct {
	mock *RecentStorageMock
}

//Return sets up a mock for RecentStorage.ClearObjects to return Return's arguments
func (m *mRecentStorageMockClearObjects) Return() *RecentStorageMock {
	m.mock.ClearObjectsFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.ClearObjects method
func (m *mRecentStorageMockClearObjects) Set(f func()) *RecentStorageMock {
	m.mock.ClearObjectsFunc = f

	return m.mock
}

//ClearObjects implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) ClearObjects() {
	atomic.AddUint64(&m.ClearObjectsPreCounter, 1)
	defer atomic.AddUint64(&m.ClearObjectsCounter, 1)

	if m.ClearObjectsFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.ClearObjects")
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

type mRecentStorageMockClearZeroTTLObjects struct {
	mock *RecentStorageMock
}

//Return sets up a mock for RecentStorage.ClearZeroTTLObjects to return Return's arguments
func (m *mRecentStorageMockClearZeroTTLObjects) Return() *RecentStorageMock {
	m.mock.ClearZeroTTLObjectsFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.ClearZeroTTLObjects method
func (m *mRecentStorageMockClearZeroTTLObjects) Set(f func()) *RecentStorageMock {
	m.mock.ClearZeroTTLObjectsFunc = f

	return m.mock
}

//ClearZeroTTLObjects implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) ClearZeroTTLObjects() {
	atomic.AddUint64(&m.ClearZeroTTLObjectsPreCounter, 1)
	defer atomic.AddUint64(&m.ClearZeroTTLObjectsCounter, 1)

	if m.ClearZeroTTLObjectsFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.ClearZeroTTLObjects")
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

type mRecentStorageMockGetObjects struct {
	mock *RecentStorageMock
}

//Return sets up a mock for RecentStorage.GetObjects to return Return's arguments
func (m *mRecentStorageMockGetObjects) Return(r map[core.RecordID]int) *RecentStorageMock {
	m.mock.GetObjectsFunc = func() map[core.RecordID]int {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.GetObjects method
func (m *mRecentStorageMockGetObjects) Set(f func() (r map[core.RecordID]int)) *RecentStorageMock {
	m.mock.GetObjectsFunc = f

	return m.mock
}

//GetObjects implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) GetObjects() (r map[core.RecordID]int) {
	atomic.AddUint64(&m.GetObjectsPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectsCounter, 1)

	if m.GetObjectsFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.GetObjects")
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

type mRecentStorageMockGetRequests struct {
	mock *RecentStorageMock
}

//Return sets up a mock for RecentStorage.GetRequests to return Return's arguments
func (m *mRecentStorageMockGetRequests) Return(r []core.RecordID) *RecentStorageMock {
	m.mock.GetRequestsFunc = func() []core.RecordID {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.GetRequests method
func (m *mRecentStorageMockGetRequests) Set(f func() (r []core.RecordID)) *RecentStorageMock {
	m.mock.GetRequestsFunc = f

	return m.mock
}

//GetRequests implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) GetRequests() (r []core.RecordID) {
	atomic.AddUint64(&m.GetRequestsPreCounter, 1)
	defer atomic.AddUint64(&m.GetRequestsCounter, 1)

	if m.GetRequestsFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.GetRequests")
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

type mRecentStorageMockIsMine struct {
	mock             *RecentStorageMock
	mockExpectations *RecentStorageMockIsMineParams
}

//RecentStorageMockIsMineParams represents input parameters of the RecentStorage.IsMine
type RecentStorageMockIsMineParams struct {
	p core.RecordID
}

//Expect sets up expected params for the RecentStorage.IsMine
func (m *mRecentStorageMockIsMine) Expect(p core.RecordID) *mRecentStorageMockIsMine {
	m.mockExpectations = &RecentStorageMockIsMineParams{p}
	return m
}

//Return sets up a mock for RecentStorage.IsMine to return Return's arguments
func (m *mRecentStorageMockIsMine) Return(r bool) *RecentStorageMock {
	m.mock.IsMineFunc = func(p core.RecordID) bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.IsMine method
func (m *mRecentStorageMockIsMine) Set(f func(p core.RecordID) (r bool)) *RecentStorageMock {
	m.mock.IsMineFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IsMine implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) IsMine(p core.RecordID) (r bool) {
	atomic.AddUint64(&m.IsMinePreCounter, 1)
	defer atomic.AddUint64(&m.IsMineCounter, 1)

	if m.IsMineMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IsMineMock.mockExpectations, RecentStorageMockIsMineParams{p},
			"RecentStorage.IsMine got unexpected parameters")

		if m.IsMineFunc == nil {

			m.t.Fatal("No results are set for the RecentStorageMock.IsMine")

			return
		}
	}

	if m.IsMineFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.IsMine")
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

type mRecentStorageMockRemovePendingRequest struct {
	mock             *RecentStorageMock
	mockExpectations *RecentStorageMockRemovePendingRequestParams
}

//RecentStorageMockRemovePendingRequestParams represents input parameters of the RecentStorage.RemovePendingRequest
type RecentStorageMockRemovePendingRequestParams struct {
	p core.RecordID
}

//Expect sets up expected params for the RecentStorage.RemovePendingRequest
func (m *mRecentStorageMockRemovePendingRequest) Expect(p core.RecordID) *mRecentStorageMockRemovePendingRequest {
	m.mockExpectations = &RecentStorageMockRemovePendingRequestParams{p}
	return m
}

//Return sets up a mock for RecentStorage.RemovePendingRequest to return Return's arguments
func (m *mRecentStorageMockRemovePendingRequest) Return() *RecentStorageMock {
	m.mock.RemovePendingRequestFunc = func(p core.RecordID) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RecentStorage.RemovePendingRequest method
func (m *mRecentStorageMockRemovePendingRequest) Set(f func(p core.RecordID)) *RecentStorageMock {
	m.mock.RemovePendingRequestFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RemovePendingRequest implements github.com/insolar/insolar/ledger/recentstorage.RecentStorage interface
func (m *RecentStorageMock) RemovePendingRequest(p core.RecordID) {
	atomic.AddUint64(&m.RemovePendingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RemovePendingRequestCounter, 1)

	if m.RemovePendingRequestMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RemovePendingRequestMock.mockExpectations, RecentStorageMockRemovePendingRequestParams{p},
			"RecentStorage.RemovePendingRequest got unexpected parameters")

		if m.RemovePendingRequestFunc == nil {

			m.t.Fatal("No results are set for the RecentStorageMock.RemovePendingRequest")

			return
		}
	}

	if m.RemovePendingRequestFunc == nil {
		m.t.Fatal("Unexpected call to RecentStorageMock.RemovePendingRequest")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecentStorageMock) ValidateCallCounters() {

	if m.AddObjectFunc != nil && atomic.LoadUint64(&m.AddObjectCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.AddObject")
	}

	if m.AddObjectWithTLLFunc != nil && atomic.LoadUint64(&m.AddObjectWithTLLCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.AddObjectWithTLL")
	}

	if m.AddPendingRequestFunc != nil && atomic.LoadUint64(&m.AddPendingRequestCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.AddPendingRequest")
	}

	if m.ClearObjectsFunc != nil && atomic.LoadUint64(&m.ClearObjectsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.ClearObjects")
	}

	if m.ClearZeroTTLObjectsFunc != nil && atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.ClearZeroTTLObjects")
	}

	if m.GetObjectsFunc != nil && atomic.LoadUint64(&m.GetObjectsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.GetObjects")
	}

	if m.GetRequestsFunc != nil && atomic.LoadUint64(&m.GetRequestsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.GetRequests")
	}

	if m.IsMineFunc != nil && atomic.LoadUint64(&m.IsMineCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.IsMine")
	}

	if m.RemovePendingRequestFunc != nil && atomic.LoadUint64(&m.RemovePendingRequestCounter) == 0 {
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

	if m.AddObjectFunc != nil && atomic.LoadUint64(&m.AddObjectCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.AddObject")
	}

	if m.AddObjectWithTLLFunc != nil && atomic.LoadUint64(&m.AddObjectWithTLLCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.AddObjectWithTLL")
	}

	if m.AddPendingRequestFunc != nil && atomic.LoadUint64(&m.AddPendingRequestCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.AddPendingRequest")
	}

	if m.ClearObjectsFunc != nil && atomic.LoadUint64(&m.ClearObjectsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.ClearObjects")
	}

	if m.ClearZeroTTLObjectsFunc != nil && atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.ClearZeroTTLObjects")
	}

	if m.GetObjectsFunc != nil && atomic.LoadUint64(&m.GetObjectsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.GetObjects")
	}

	if m.GetRequestsFunc != nil && atomic.LoadUint64(&m.GetRequestsCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.GetRequests")
	}

	if m.IsMineFunc != nil && atomic.LoadUint64(&m.IsMineCounter) == 0 {
		m.t.Fatal("Expected call to RecentStorageMock.IsMine")
	}

	if m.RemovePendingRequestFunc != nil && atomic.LoadUint64(&m.RemovePendingRequestCounter) == 0 {
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
		ok = ok && (m.AddObjectFunc == nil || atomic.LoadUint64(&m.AddObjectCounter) > 0)
		ok = ok && (m.AddObjectWithTLLFunc == nil || atomic.LoadUint64(&m.AddObjectWithTLLCounter) > 0)
		ok = ok && (m.AddPendingRequestFunc == nil || atomic.LoadUint64(&m.AddPendingRequestCounter) > 0)
		ok = ok && (m.ClearObjectsFunc == nil || atomic.LoadUint64(&m.ClearObjectsCounter) > 0)
		ok = ok && (m.ClearZeroTTLObjectsFunc == nil || atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) > 0)
		ok = ok && (m.GetObjectsFunc == nil || atomic.LoadUint64(&m.GetObjectsCounter) > 0)
		ok = ok && (m.GetRequestsFunc == nil || atomic.LoadUint64(&m.GetRequestsCounter) > 0)
		ok = ok && (m.IsMineFunc == nil || atomic.LoadUint64(&m.IsMineCounter) > 0)
		ok = ok && (m.RemovePendingRequestFunc == nil || atomic.LoadUint64(&m.RemovePendingRequestCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.AddObjectFunc != nil && atomic.LoadUint64(&m.AddObjectCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.AddObject")
			}

			if m.AddObjectWithTLLFunc != nil && atomic.LoadUint64(&m.AddObjectWithTLLCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.AddObjectWithTLL")
			}

			if m.AddPendingRequestFunc != nil && atomic.LoadUint64(&m.AddPendingRequestCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.AddPendingRequest")
			}

			if m.ClearObjectsFunc != nil && atomic.LoadUint64(&m.ClearObjectsCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.ClearObjects")
			}

			if m.ClearZeroTTLObjectsFunc != nil && atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.ClearZeroTTLObjects")
			}

			if m.GetObjectsFunc != nil && atomic.LoadUint64(&m.GetObjectsCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.GetObjects")
			}

			if m.GetRequestsFunc != nil && atomic.LoadUint64(&m.GetRequestsCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.GetRequests")
			}

			if m.IsMineFunc != nil && atomic.LoadUint64(&m.IsMineCounter) == 0 {
				m.t.Error("Expected call to RecentStorageMock.IsMine")
			}

			if m.RemovePendingRequestFunc != nil && atomic.LoadUint64(&m.RemovePendingRequestCounter) == 0 {
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

	if m.AddObjectFunc != nil && atomic.LoadUint64(&m.AddObjectCounter) == 0 {
		return false
	}

	if m.AddObjectWithTLLFunc != nil && atomic.LoadUint64(&m.AddObjectWithTLLCounter) == 0 {
		return false
	}

	if m.AddPendingRequestFunc != nil && atomic.LoadUint64(&m.AddPendingRequestCounter) == 0 {
		return false
	}

	if m.ClearObjectsFunc != nil && atomic.LoadUint64(&m.ClearObjectsCounter) == 0 {
		return false
	}

	if m.ClearZeroTTLObjectsFunc != nil && atomic.LoadUint64(&m.ClearZeroTTLObjectsCounter) == 0 {
		return false
	}

	if m.GetObjectsFunc != nil && atomic.LoadUint64(&m.GetObjectsCounter) == 0 {
		return false
	}

	if m.GetRequestsFunc != nil && atomic.LoadUint64(&m.GetRequestsCounter) == 0 {
		return false
	}

	if m.IsMineFunc != nil && atomic.LoadUint64(&m.IsMineCounter) == 0 {
		return false
	}

	if m.RemovePendingRequestFunc != nil && atomic.LoadUint64(&m.RemovePendingRequestCounter) == 0 {
		return false
	}

	return true
}

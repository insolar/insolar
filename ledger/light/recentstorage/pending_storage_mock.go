package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PendingStorage" can be found in github.com/insolar/insolar/ledger/light/recentstorage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//PendingStorageMock implements github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage
type PendingStorageMock struct {
	t minimock.Tester

	AddPendingRequestFunc       func(p context.Context, p1 insolar.ID, p2 insolar.ID)
	AddPendingRequestCounter    uint64
	AddPendingRequestPreCounter uint64
	AddPendingRequestMock       mPendingStorageMockAddPendingRequest

	GetRequestsFunc       func() (r map[insolar.ID]PendingObjectContext)
	GetRequestsCounter    uint64
	GetRequestsPreCounter uint64
	GetRequestsMock       mPendingStorageMockGetRequests

	GetRequestsForObjectFunc       func(p insolar.ID) (r []insolar.ID)
	GetRequestsForObjectCounter    uint64
	GetRequestsForObjectPreCounter uint64
	GetRequestsForObjectMock       mPendingStorageMockGetRequestsForObject

	RemovePendingRequestFunc       func(p context.Context, p1 insolar.ID, p2 insolar.ID)
	RemovePendingRequestCounter    uint64
	RemovePendingRequestPreCounter uint64
	RemovePendingRequestMock       mPendingStorageMockRemovePendingRequest

	SetContextToObjectFunc       func(p context.Context, p1 insolar.ID, p2 PendingObjectContext)
	SetContextToObjectCounter    uint64
	SetContextToObjectPreCounter uint64
	SetContextToObjectMock       mPendingStorageMockSetContextToObject
}

//NewPendingStorageMock returns a mock for github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage
func NewPendingStorageMock(t minimock.Tester) *PendingStorageMock {
	m := &PendingStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddPendingRequestMock = mPendingStorageMockAddPendingRequest{mock: m}
	m.GetRequestsMock = mPendingStorageMockGetRequests{mock: m}
	m.GetRequestsForObjectMock = mPendingStorageMockGetRequestsForObject{mock: m}
	m.RemovePendingRequestMock = mPendingStorageMockRemovePendingRequest{mock: m}
	m.SetContextToObjectMock = mPendingStorageMockSetContextToObject{mock: m}

	return m
}

type mPendingStorageMockAddPendingRequest struct {
	mock              *PendingStorageMock
	mainExpectation   *PendingStorageMockAddPendingRequestExpectation
	expectationSeries []*PendingStorageMockAddPendingRequestExpectation
}

type PendingStorageMockAddPendingRequestExpectation struct {
	input *PendingStorageMockAddPendingRequestInput
}

type PendingStorageMockAddPendingRequestInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.ID
}

//Expect specifies that invocation of PendingStorage.AddPendingRequest is expected from 1 to Infinity times
func (m *mPendingStorageMockAddPendingRequest) Expect(p context.Context, p1 insolar.ID, p2 insolar.ID) *mPendingStorageMockAddPendingRequest {
	m.mock.AddPendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockAddPendingRequestExpectation{}
	}
	m.mainExpectation.input = &PendingStorageMockAddPendingRequestInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PendingStorage.AddPendingRequest
func (m *mPendingStorageMockAddPendingRequest) Return() *PendingStorageMock {
	m.mock.AddPendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockAddPendingRequestExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PendingStorage.AddPendingRequest is expected once
func (m *mPendingStorageMockAddPendingRequest) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.ID) *PendingStorageMockAddPendingRequestExpectation {
	m.mock.AddPendingRequestFunc = nil
	m.mainExpectation = nil

	expectation := &PendingStorageMockAddPendingRequestExpectation{}
	expectation.input = &PendingStorageMockAddPendingRequestInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PendingStorage.AddPendingRequest method
func (m *mPendingStorageMockAddPendingRequest) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.ID)) *PendingStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddPendingRequestFunc = f
	return m.mock
}

//AddPendingRequest implements github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage interface
func (m *PendingStorageMock) AddPendingRequest(p context.Context, p1 insolar.ID, p2 insolar.ID) {
	counter := atomic.AddUint64(&m.AddPendingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.AddPendingRequestCounter, 1)

	if len(m.AddPendingRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddPendingRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingStorageMock.AddPendingRequest. %v %v %v", p, p1, p2)
			return
		}

		input := m.AddPendingRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingStorageMockAddPendingRequestInput{p, p1, p2}, "PendingStorage.AddPendingRequest got unexpected parameters")

		return
	}

	if m.AddPendingRequestMock.mainExpectation != nil {

		input := m.AddPendingRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingStorageMockAddPendingRequestInput{p, p1, p2}, "PendingStorage.AddPendingRequest got unexpected parameters")
		}

		return
	}

	if m.AddPendingRequestFunc == nil {
		m.t.Fatalf("Unexpected call to PendingStorageMock.AddPendingRequest. %v %v %v", p, p1, p2)
		return
	}

	m.AddPendingRequestFunc(p, p1, p2)
}

//AddPendingRequestMinimockCounter returns a count of PendingStorageMock.AddPendingRequestFunc invocations
func (m *PendingStorageMock) AddPendingRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddPendingRequestCounter)
}

//AddPendingRequestMinimockPreCounter returns the value of PendingStorageMock.AddPendingRequest invocations
func (m *PendingStorageMock) AddPendingRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddPendingRequestPreCounter)
}

//AddPendingRequestFinished returns true if mock invocations count is ok
func (m *PendingStorageMock) AddPendingRequestFinished() bool {
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

type mPendingStorageMockGetRequests struct {
	mock              *PendingStorageMock
	mainExpectation   *PendingStorageMockGetRequestsExpectation
	expectationSeries []*PendingStorageMockGetRequestsExpectation
}

type PendingStorageMockGetRequestsExpectation struct {
	result *PendingStorageMockGetRequestsResult
}

type PendingStorageMockGetRequestsResult struct {
	r map[insolar.ID]PendingObjectContext
}

//Expect specifies that invocation of PendingStorage.GetRequests is expected from 1 to Infinity times
func (m *mPendingStorageMockGetRequests) Expect() *mPendingStorageMockGetRequests {
	m.mock.GetRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockGetRequestsExpectation{}
	}

	return m
}

//Return specifies results of invocation of PendingStorage.GetRequests
func (m *mPendingStorageMockGetRequests) Return(r map[insolar.ID]PendingObjectContext) *PendingStorageMock {
	m.mock.GetRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockGetRequestsExpectation{}
	}
	m.mainExpectation.result = &PendingStorageMockGetRequestsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingStorage.GetRequests is expected once
func (m *mPendingStorageMockGetRequests) ExpectOnce() *PendingStorageMockGetRequestsExpectation {
	m.mock.GetRequestsFunc = nil
	m.mainExpectation = nil

	expectation := &PendingStorageMockGetRequestsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingStorageMockGetRequestsExpectation) Return(r map[insolar.ID]PendingObjectContext) {
	e.result = &PendingStorageMockGetRequestsResult{r}
}

//Set uses given function f as a mock of PendingStorage.GetRequests method
func (m *mPendingStorageMockGetRequests) Set(f func() (r map[insolar.ID]PendingObjectContext)) *PendingStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRequestsFunc = f
	return m.mock
}

//GetRequests implements github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage interface
func (m *PendingStorageMock) GetRequests() (r map[insolar.ID]PendingObjectContext) {
	counter := atomic.AddUint64(&m.GetRequestsPreCounter, 1)
	defer atomic.AddUint64(&m.GetRequestsCounter, 1)

	if len(m.GetRequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingStorageMock.GetRequests.")
			return
		}

		result := m.GetRequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingStorageMock.GetRequests")
			return
		}

		r = result.r

		return
	}

	if m.GetRequestsMock.mainExpectation != nil {

		result := m.GetRequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingStorageMock.GetRequests")
		}

		r = result.r

		return
	}

	if m.GetRequestsFunc == nil {
		m.t.Fatalf("Unexpected call to PendingStorageMock.GetRequests.")
		return
	}

	return m.GetRequestsFunc()
}

//GetRequestsMinimockCounter returns a count of PendingStorageMock.GetRequestsFunc invocations
func (m *PendingStorageMock) GetRequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRequestsCounter)
}

//GetRequestsMinimockPreCounter returns the value of PendingStorageMock.GetRequests invocations
func (m *PendingStorageMock) GetRequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRequestsPreCounter)
}

//GetRequestsFinished returns true if mock invocations count is ok
func (m *PendingStorageMock) GetRequestsFinished() bool {
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

type mPendingStorageMockGetRequestsForObject struct {
	mock              *PendingStorageMock
	mainExpectation   *PendingStorageMockGetRequestsForObjectExpectation
	expectationSeries []*PendingStorageMockGetRequestsForObjectExpectation
}

type PendingStorageMockGetRequestsForObjectExpectation struct {
	input  *PendingStorageMockGetRequestsForObjectInput
	result *PendingStorageMockGetRequestsForObjectResult
}

type PendingStorageMockGetRequestsForObjectInput struct {
	p insolar.ID
}

type PendingStorageMockGetRequestsForObjectResult struct {
	r []insolar.ID
}

//Expect specifies that invocation of PendingStorage.GetRequestsForObject is expected from 1 to Infinity times
func (m *mPendingStorageMockGetRequestsForObject) Expect(p insolar.ID) *mPendingStorageMockGetRequestsForObject {
	m.mock.GetRequestsForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockGetRequestsForObjectExpectation{}
	}
	m.mainExpectation.input = &PendingStorageMockGetRequestsForObjectInput{p}
	return m
}

//Return specifies results of invocation of PendingStorage.GetRequestsForObject
func (m *mPendingStorageMockGetRequestsForObject) Return(r []insolar.ID) *PendingStorageMock {
	m.mock.GetRequestsForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockGetRequestsForObjectExpectation{}
	}
	m.mainExpectation.result = &PendingStorageMockGetRequestsForObjectResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingStorage.GetRequestsForObject is expected once
func (m *mPendingStorageMockGetRequestsForObject) ExpectOnce(p insolar.ID) *PendingStorageMockGetRequestsForObjectExpectation {
	m.mock.GetRequestsForObjectFunc = nil
	m.mainExpectation = nil

	expectation := &PendingStorageMockGetRequestsForObjectExpectation{}
	expectation.input = &PendingStorageMockGetRequestsForObjectInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingStorageMockGetRequestsForObjectExpectation) Return(r []insolar.ID) {
	e.result = &PendingStorageMockGetRequestsForObjectResult{r}
}

//Set uses given function f as a mock of PendingStorage.GetRequestsForObject method
func (m *mPendingStorageMockGetRequestsForObject) Set(f func(p insolar.ID) (r []insolar.ID)) *PendingStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRequestsForObjectFunc = f
	return m.mock
}

//GetRequestsForObject implements github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage interface
func (m *PendingStorageMock) GetRequestsForObject(p insolar.ID) (r []insolar.ID) {
	counter := atomic.AddUint64(&m.GetRequestsForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.GetRequestsForObjectCounter, 1)

	if len(m.GetRequestsForObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRequestsForObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingStorageMock.GetRequestsForObject. %v", p)
			return
		}

		input := m.GetRequestsForObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingStorageMockGetRequestsForObjectInput{p}, "PendingStorage.GetRequestsForObject got unexpected parameters")

		result := m.GetRequestsForObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingStorageMock.GetRequestsForObject")
			return
		}

		r = result.r

		return
	}

	if m.GetRequestsForObjectMock.mainExpectation != nil {

		input := m.GetRequestsForObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingStorageMockGetRequestsForObjectInput{p}, "PendingStorage.GetRequestsForObject got unexpected parameters")
		}

		result := m.GetRequestsForObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingStorageMock.GetRequestsForObject")
		}

		r = result.r

		return
	}

	if m.GetRequestsForObjectFunc == nil {
		m.t.Fatalf("Unexpected call to PendingStorageMock.GetRequestsForObject. %v", p)
		return
	}

	return m.GetRequestsForObjectFunc(p)
}

//GetRequestsForObjectMinimockCounter returns a count of PendingStorageMock.GetRequestsForObjectFunc invocations
func (m *PendingStorageMock) GetRequestsForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRequestsForObjectCounter)
}

//GetRequestsForObjectMinimockPreCounter returns the value of PendingStorageMock.GetRequestsForObject invocations
func (m *PendingStorageMock) GetRequestsForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRequestsForObjectPreCounter)
}

//GetRequestsForObjectFinished returns true if mock invocations count is ok
func (m *PendingStorageMock) GetRequestsForObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRequestsForObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRequestsForObjectCounter) == uint64(len(m.GetRequestsForObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRequestsForObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRequestsForObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRequestsForObjectFunc != nil {
		return atomic.LoadUint64(&m.GetRequestsForObjectCounter) > 0
	}

	return true
}

type mPendingStorageMockRemovePendingRequest struct {
	mock              *PendingStorageMock
	mainExpectation   *PendingStorageMockRemovePendingRequestExpectation
	expectationSeries []*PendingStorageMockRemovePendingRequestExpectation
}

type PendingStorageMockRemovePendingRequestExpectation struct {
	input *PendingStorageMockRemovePendingRequestInput
}

type PendingStorageMockRemovePendingRequestInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.ID
}

//Expect specifies that invocation of PendingStorage.RemovePendingRequest is expected from 1 to Infinity times
func (m *mPendingStorageMockRemovePendingRequest) Expect(p context.Context, p1 insolar.ID, p2 insolar.ID) *mPendingStorageMockRemovePendingRequest {
	m.mock.RemovePendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockRemovePendingRequestExpectation{}
	}
	m.mainExpectation.input = &PendingStorageMockRemovePendingRequestInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PendingStorage.RemovePendingRequest
func (m *mPendingStorageMockRemovePendingRequest) Return() *PendingStorageMock {
	m.mock.RemovePendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockRemovePendingRequestExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PendingStorage.RemovePendingRequest is expected once
func (m *mPendingStorageMockRemovePendingRequest) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.ID) *PendingStorageMockRemovePendingRequestExpectation {
	m.mock.RemovePendingRequestFunc = nil
	m.mainExpectation = nil

	expectation := &PendingStorageMockRemovePendingRequestExpectation{}
	expectation.input = &PendingStorageMockRemovePendingRequestInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PendingStorage.RemovePendingRequest method
func (m *mPendingStorageMockRemovePendingRequest) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.ID)) *PendingStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemovePendingRequestFunc = f
	return m.mock
}

//RemovePendingRequest implements github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage interface
func (m *PendingStorageMock) RemovePendingRequest(p context.Context, p1 insolar.ID, p2 insolar.ID) {
	counter := atomic.AddUint64(&m.RemovePendingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RemovePendingRequestCounter, 1)

	if len(m.RemovePendingRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemovePendingRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingStorageMock.RemovePendingRequest. %v %v %v", p, p1, p2)
			return
		}

		input := m.RemovePendingRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingStorageMockRemovePendingRequestInput{p, p1, p2}, "PendingStorage.RemovePendingRequest got unexpected parameters")

		return
	}

	if m.RemovePendingRequestMock.mainExpectation != nil {

		input := m.RemovePendingRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingStorageMockRemovePendingRequestInput{p, p1, p2}, "PendingStorage.RemovePendingRequest got unexpected parameters")
		}

		return
	}

	if m.RemovePendingRequestFunc == nil {
		m.t.Fatalf("Unexpected call to PendingStorageMock.RemovePendingRequest. %v %v %v", p, p1, p2)
		return
	}

	m.RemovePendingRequestFunc(p, p1, p2)
}

//RemovePendingRequestMinimockCounter returns a count of PendingStorageMock.RemovePendingRequestFunc invocations
func (m *PendingStorageMock) RemovePendingRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingRequestCounter)
}

//RemovePendingRequestMinimockPreCounter returns the value of PendingStorageMock.RemovePendingRequest invocations
func (m *PendingStorageMock) RemovePendingRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingRequestPreCounter)
}

//RemovePendingRequestFinished returns true if mock invocations count is ok
func (m *PendingStorageMock) RemovePendingRequestFinished() bool {
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

type mPendingStorageMockSetContextToObject struct {
	mock              *PendingStorageMock
	mainExpectation   *PendingStorageMockSetContextToObjectExpectation
	expectationSeries []*PendingStorageMockSetContextToObjectExpectation
}

type PendingStorageMockSetContextToObjectExpectation struct {
	input *PendingStorageMockSetContextToObjectInput
}

type PendingStorageMockSetContextToObjectInput struct {
	p  context.Context
	p1 insolar.ID
	p2 PendingObjectContext
}

//Expect specifies that invocation of PendingStorage.SetContextToObject is expected from 1 to Infinity times
func (m *mPendingStorageMockSetContextToObject) Expect(p context.Context, p1 insolar.ID, p2 PendingObjectContext) *mPendingStorageMockSetContextToObject {
	m.mock.SetContextToObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockSetContextToObjectExpectation{}
	}
	m.mainExpectation.input = &PendingStorageMockSetContextToObjectInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PendingStorage.SetContextToObject
func (m *mPendingStorageMockSetContextToObject) Return() *PendingStorageMock {
	m.mock.SetContextToObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingStorageMockSetContextToObjectExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PendingStorage.SetContextToObject is expected once
func (m *mPendingStorageMockSetContextToObject) ExpectOnce(p context.Context, p1 insolar.ID, p2 PendingObjectContext) *PendingStorageMockSetContextToObjectExpectation {
	m.mock.SetContextToObjectFunc = nil
	m.mainExpectation = nil

	expectation := &PendingStorageMockSetContextToObjectExpectation{}
	expectation.input = &PendingStorageMockSetContextToObjectInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PendingStorage.SetContextToObject method
func (m *mPendingStorageMockSetContextToObject) Set(f func(p context.Context, p1 insolar.ID, p2 PendingObjectContext)) *PendingStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetContextToObjectFunc = f
	return m.mock
}

//SetContextToObject implements github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage interface
func (m *PendingStorageMock) SetContextToObject(p context.Context, p1 insolar.ID, p2 PendingObjectContext) {
	counter := atomic.AddUint64(&m.SetContextToObjectPreCounter, 1)
	defer atomic.AddUint64(&m.SetContextToObjectCounter, 1)

	if len(m.SetContextToObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetContextToObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingStorageMock.SetContextToObject. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetContextToObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingStorageMockSetContextToObjectInput{p, p1, p2}, "PendingStorage.SetContextToObject got unexpected parameters")

		return
	}

	if m.SetContextToObjectMock.mainExpectation != nil {

		input := m.SetContextToObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingStorageMockSetContextToObjectInput{p, p1, p2}, "PendingStorage.SetContextToObject got unexpected parameters")
		}

		return
	}

	if m.SetContextToObjectFunc == nil {
		m.t.Fatalf("Unexpected call to PendingStorageMock.SetContextToObject. %v %v %v", p, p1, p2)
		return
	}

	m.SetContextToObjectFunc(p, p1, p2)
}

//SetContextToObjectMinimockCounter returns a count of PendingStorageMock.SetContextToObjectFunc invocations
func (m *PendingStorageMock) SetContextToObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetContextToObjectCounter)
}

//SetContextToObjectMinimockPreCounter returns the value of PendingStorageMock.SetContextToObject invocations
func (m *PendingStorageMock) SetContextToObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetContextToObjectPreCounter)
}

//SetContextToObjectFinished returns true if mock invocations count is ok
func (m *PendingStorageMock) SetContextToObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetContextToObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetContextToObjectCounter) == uint64(len(m.SetContextToObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetContextToObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetContextToObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetContextToObjectFunc != nil {
		return atomic.LoadUint64(&m.SetContextToObjectCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingStorageMock) ValidateCallCounters() {

	if !m.AddPendingRequestFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.AddPendingRequest")
	}

	if !m.GetRequestsFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.GetRequests")
	}

	if !m.GetRequestsForObjectFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.GetRequestsForObject")
	}

	if !m.RemovePendingRequestFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.RemovePendingRequest")
	}

	if !m.SetContextToObjectFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.SetContextToObject")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PendingStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PendingStorageMock) MinimockFinish() {

	if !m.AddPendingRequestFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.AddPendingRequest")
	}

	if !m.GetRequestsFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.GetRequests")
	}

	if !m.GetRequestsForObjectFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.GetRequestsForObject")
	}

	if !m.RemovePendingRequestFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.RemovePendingRequest")
	}

	if !m.SetContextToObjectFinished() {
		m.t.Fatal("Expected call to PendingStorageMock.SetContextToObject")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PendingStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PendingStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddPendingRequestFinished()
		ok = ok && m.GetRequestsFinished()
		ok = ok && m.GetRequestsForObjectFinished()
		ok = ok && m.RemovePendingRequestFinished()
		ok = ok && m.SetContextToObjectFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddPendingRequestFinished() {
				m.t.Error("Expected call to PendingStorageMock.AddPendingRequest")
			}

			if !m.GetRequestsFinished() {
				m.t.Error("Expected call to PendingStorageMock.GetRequests")
			}

			if !m.GetRequestsForObjectFinished() {
				m.t.Error("Expected call to PendingStorageMock.GetRequestsForObject")
			}

			if !m.RemovePendingRequestFinished() {
				m.t.Error("Expected call to PendingStorageMock.RemovePendingRequest")
			}

			if !m.SetContextToObjectFinished() {
				m.t.Error("Expected call to PendingStorageMock.SetContextToObject")
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
func (m *PendingStorageMock) AllMocksCalled() bool {

	if !m.AddPendingRequestFinished() {
		return false
	}

	if !m.GetRequestsFinished() {
		return false
	}

	if !m.GetRequestsForObjectFinished() {
		return false
	}

	if !m.RemovePendingRequestFinished() {
		return false
	}

	if !m.SetContextToObjectFinished() {
		return false
	}

	return true
}

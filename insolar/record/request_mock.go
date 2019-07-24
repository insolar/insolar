package record

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Request" can be found in github.com/insolar/insolar/insolar/record
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
)

//RequestMock implements github.com/insolar/insolar/insolar/record.Request
type RequestMock struct {
	t minimock.Tester

	AffinityRefFunc       func() (r *insolar.Reference)
	AffinityRefCounter    uint64
	AffinityRefPreCounter uint64
	AffinityRefMock       mRequestMockAffinityRef

	GetCallTypeFunc       func() (r CallType)
	GetCallTypeCounter    uint64
	GetCallTypePreCounter uint64
	GetCallTypeMock       mRequestMockGetCallType

	IsAPIRequestFunc       func() (r bool)
	IsAPIRequestCounter    uint64
	IsAPIRequestPreCounter uint64
	IsAPIRequestMock       mRequestMockIsAPIRequest

	IsCreationRequestFunc       func() (r bool)
	IsCreationRequestCounter    uint64
	IsCreationRequestPreCounter uint64
	IsCreationRequestMock       mRequestMockIsCreationRequest

	IsDetachedFunc       func() (r bool)
	IsDetachedCounter    uint64
	IsDetachedPreCounter uint64
	IsDetachedMock       mRequestMockIsDetached

	IsEmptyAPINodeFunc       func() (r bool)
	IsEmptyAPINodeCounter    uint64
	IsEmptyAPINodePreCounter uint64
	IsEmptyAPINodeMock       mRequestMockIsEmptyAPINode

	ReasonRefFunc       func() (r insolar.Reference)
	ReasonRefCounter    uint64
	ReasonRefPreCounter uint64
	ReasonRefMock       mRequestMockReasonRef
}

//NewRequestMock returns a mock for github.com/insolar/insolar/insolar/record.Request
func NewRequestMock(t minimock.Tester) *RequestMock {
	m := &RequestMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AffinityRefMock = mRequestMockAffinityRef{mock: m}
	m.GetCallTypeMock = mRequestMockGetCallType{mock: m}
	m.IsAPIRequestMock = mRequestMockIsAPIRequest{mock: m}
	m.IsCreationRequestMock = mRequestMockIsCreationRequest{mock: m}
	m.IsDetachedMock = mRequestMockIsDetached{mock: m}
	m.IsEmptyAPINodeMock = mRequestMockIsEmptyAPINode{mock: m}
	m.ReasonRefMock = mRequestMockReasonRef{mock: m}

	return m
}

type mRequestMockAffinityRef struct {
	mock              *RequestMock
	mainExpectation   *RequestMockAffinityRefExpectation
	expectationSeries []*RequestMockAffinityRefExpectation
}

type RequestMockAffinityRefExpectation struct {
	result *RequestMockAffinityRefResult
}

type RequestMockAffinityRefResult struct {
	r *insolar.Reference
}

//Expect specifies that invocation of Request.AffinityRef is expected from 1 to Infinity times
func (m *mRequestMockAffinityRef) Expect() *mRequestMockAffinityRef {
	m.mock.AffinityRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockAffinityRefExpectation{}
	}

	return m
}

//Return specifies results of invocation of Request.AffinityRef
func (m *mRequestMockAffinityRef) Return(r *insolar.Reference) *RequestMock {
	m.mock.AffinityRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockAffinityRefExpectation{}
	}
	m.mainExpectation.result = &RequestMockAffinityRefResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Request.AffinityRef is expected once
func (m *mRequestMockAffinityRef) ExpectOnce() *RequestMockAffinityRefExpectation {
	m.mock.AffinityRefFunc = nil
	m.mainExpectation = nil

	expectation := &RequestMockAffinityRefExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestMockAffinityRefExpectation) Return(r *insolar.Reference) {
	e.result = &RequestMockAffinityRefResult{r}
}

//Set uses given function f as a mock of Request.AffinityRef method
func (m *mRequestMockAffinityRef) Set(f func() (r *insolar.Reference)) *RequestMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AffinityRefFunc = f
	return m.mock
}

//AffinityRef implements github.com/insolar/insolar/insolar/record.Request interface
func (m *RequestMock) AffinityRef() (r *insolar.Reference) {
	counter := atomic.AddUint64(&m.AffinityRefPreCounter, 1)
	defer atomic.AddUint64(&m.AffinityRefCounter, 1)

	if len(m.AffinityRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AffinityRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestMock.AffinityRef.")
			return
		}

		result := m.AffinityRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.AffinityRef")
			return
		}

		r = result.r

		return
	}

	if m.AffinityRefMock.mainExpectation != nil {

		result := m.AffinityRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.AffinityRef")
		}

		r = result.r

		return
	}

	if m.AffinityRefFunc == nil {
		m.t.Fatalf("Unexpected call to RequestMock.AffinityRef.")
		return
	}

	return m.AffinityRefFunc()
}

//AffinityRefMinimockCounter returns a count of RequestMock.AffinityRefFunc invocations
func (m *RequestMock) AffinityRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AffinityRefCounter)
}

//AffinityRefMinimockPreCounter returns the value of RequestMock.AffinityRef invocations
func (m *RequestMock) AffinityRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AffinityRefPreCounter)
}

//AffinityRefFinished returns true if mock invocations count is ok
func (m *RequestMock) AffinityRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AffinityRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AffinityRefCounter) == uint64(len(m.AffinityRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AffinityRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AffinityRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AffinityRefFunc != nil {
		return atomic.LoadUint64(&m.AffinityRefCounter) > 0
	}

	return true
}

type mRequestMockGetCallType struct {
	mock              *RequestMock
	mainExpectation   *RequestMockGetCallTypeExpectation
	expectationSeries []*RequestMockGetCallTypeExpectation
}

type RequestMockGetCallTypeExpectation struct {
	result *RequestMockGetCallTypeResult
}

type RequestMockGetCallTypeResult struct {
	r CallType
}

//Expect specifies that invocation of Request.GetCallType is expected from 1 to Infinity times
func (m *mRequestMockGetCallType) Expect() *mRequestMockGetCallType {
	m.mock.GetCallTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockGetCallTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Request.GetCallType
func (m *mRequestMockGetCallType) Return(r CallType) *RequestMock {
	m.mock.GetCallTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockGetCallTypeExpectation{}
	}
	m.mainExpectation.result = &RequestMockGetCallTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Request.GetCallType is expected once
func (m *mRequestMockGetCallType) ExpectOnce() *RequestMockGetCallTypeExpectation {
	m.mock.GetCallTypeFunc = nil
	m.mainExpectation = nil

	expectation := &RequestMockGetCallTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestMockGetCallTypeExpectation) Return(r CallType) {
	e.result = &RequestMockGetCallTypeResult{r}
}

//Set uses given function f as a mock of Request.GetCallType method
func (m *mRequestMockGetCallType) Set(f func() (r CallType)) *RequestMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCallTypeFunc = f
	return m.mock
}

//GetCallType implements github.com/insolar/insolar/insolar/record.Request interface
func (m *RequestMock) GetCallType() (r CallType) {
	counter := atomic.AddUint64(&m.GetCallTypePreCounter, 1)
	defer atomic.AddUint64(&m.GetCallTypeCounter, 1)

	if len(m.GetCallTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCallTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestMock.GetCallType.")
			return
		}

		result := m.GetCallTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.GetCallType")
			return
		}

		r = result.r

		return
	}

	if m.GetCallTypeMock.mainExpectation != nil {

		result := m.GetCallTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.GetCallType")
		}

		r = result.r

		return
	}

	if m.GetCallTypeFunc == nil {
		m.t.Fatalf("Unexpected call to RequestMock.GetCallType.")
		return
	}

	return m.GetCallTypeFunc()
}

//GetCallTypeMinimockCounter returns a count of RequestMock.GetCallTypeFunc invocations
func (m *RequestMock) GetCallTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCallTypeCounter)
}

//GetCallTypeMinimockPreCounter returns the value of RequestMock.GetCallType invocations
func (m *RequestMock) GetCallTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCallTypePreCounter)
}

//GetCallTypeFinished returns true if mock invocations count is ok
func (m *RequestMock) GetCallTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCallTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCallTypeCounter) == uint64(len(m.GetCallTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCallTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCallTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCallTypeFunc != nil {
		return atomic.LoadUint64(&m.GetCallTypeCounter) > 0
	}

	return true
}

type mRequestMockIsAPIRequest struct {
	mock              *RequestMock
	mainExpectation   *RequestMockIsAPIRequestExpectation
	expectationSeries []*RequestMockIsAPIRequestExpectation
}

type RequestMockIsAPIRequestExpectation struct {
	result *RequestMockIsAPIRequestResult
}

type RequestMockIsAPIRequestResult struct {
	r bool
}

//Expect specifies that invocation of Request.IsAPIRequest is expected from 1 to Infinity times
func (m *mRequestMockIsAPIRequest) Expect() *mRequestMockIsAPIRequest {
	m.mock.IsAPIRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsAPIRequestExpectation{}
	}

	return m
}

//Return specifies results of invocation of Request.IsAPIRequest
func (m *mRequestMockIsAPIRequest) Return(r bool) *RequestMock {
	m.mock.IsAPIRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsAPIRequestExpectation{}
	}
	m.mainExpectation.result = &RequestMockIsAPIRequestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Request.IsAPIRequest is expected once
func (m *mRequestMockIsAPIRequest) ExpectOnce() *RequestMockIsAPIRequestExpectation {
	m.mock.IsAPIRequestFunc = nil
	m.mainExpectation = nil

	expectation := &RequestMockIsAPIRequestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestMockIsAPIRequestExpectation) Return(r bool) {
	e.result = &RequestMockIsAPIRequestResult{r}
}

//Set uses given function f as a mock of Request.IsAPIRequest method
func (m *mRequestMockIsAPIRequest) Set(f func() (r bool)) *RequestMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAPIRequestFunc = f
	return m.mock
}

//IsAPIRequest implements github.com/insolar/insolar/insolar/record.Request interface
func (m *RequestMock) IsAPIRequest() (r bool) {
	counter := atomic.AddUint64(&m.IsAPIRequestPreCounter, 1)
	defer atomic.AddUint64(&m.IsAPIRequestCounter, 1)

	if len(m.IsAPIRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAPIRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestMock.IsAPIRequest.")
			return
		}

		result := m.IsAPIRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsAPIRequest")
			return
		}

		r = result.r

		return
	}

	if m.IsAPIRequestMock.mainExpectation != nil {

		result := m.IsAPIRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsAPIRequest")
		}

		r = result.r

		return
	}

	if m.IsAPIRequestFunc == nil {
		m.t.Fatalf("Unexpected call to RequestMock.IsAPIRequest.")
		return
	}

	return m.IsAPIRequestFunc()
}

//IsAPIRequestMinimockCounter returns a count of RequestMock.IsAPIRequestFunc invocations
func (m *RequestMock) IsAPIRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAPIRequestCounter)
}

//IsAPIRequestMinimockPreCounter returns the value of RequestMock.IsAPIRequest invocations
func (m *RequestMock) IsAPIRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAPIRequestPreCounter)
}

//IsAPIRequestFinished returns true if mock invocations count is ok
func (m *RequestMock) IsAPIRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsAPIRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsAPIRequestCounter) == uint64(len(m.IsAPIRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsAPIRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsAPIRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsAPIRequestFunc != nil {
		return atomic.LoadUint64(&m.IsAPIRequestCounter) > 0
	}

	return true
}

type mRequestMockIsCreationRequest struct {
	mock              *RequestMock
	mainExpectation   *RequestMockIsCreationRequestExpectation
	expectationSeries []*RequestMockIsCreationRequestExpectation
}

type RequestMockIsCreationRequestExpectation struct {
	result *RequestMockIsCreationRequestResult
}

type RequestMockIsCreationRequestResult struct {
	r bool
}

//Expect specifies that invocation of Request.IsCreationRequest is expected from 1 to Infinity times
func (m *mRequestMockIsCreationRequest) Expect() *mRequestMockIsCreationRequest {
	m.mock.IsCreationRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsCreationRequestExpectation{}
	}

	return m
}

//Return specifies results of invocation of Request.IsCreationRequest
func (m *mRequestMockIsCreationRequest) Return(r bool) *RequestMock {
	m.mock.IsCreationRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsCreationRequestExpectation{}
	}
	m.mainExpectation.result = &RequestMockIsCreationRequestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Request.IsCreationRequest is expected once
func (m *mRequestMockIsCreationRequest) ExpectOnce() *RequestMockIsCreationRequestExpectation {
	m.mock.IsCreationRequestFunc = nil
	m.mainExpectation = nil

	expectation := &RequestMockIsCreationRequestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestMockIsCreationRequestExpectation) Return(r bool) {
	e.result = &RequestMockIsCreationRequestResult{r}
}

//Set uses given function f as a mock of Request.IsCreationRequest method
func (m *mRequestMockIsCreationRequest) Set(f func() (r bool)) *RequestMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsCreationRequestFunc = f
	return m.mock
}

//IsCreationRequest implements github.com/insolar/insolar/insolar/record.Request interface
func (m *RequestMock) IsCreationRequest() (r bool) {
	counter := atomic.AddUint64(&m.IsCreationRequestPreCounter, 1)
	defer atomic.AddUint64(&m.IsCreationRequestCounter, 1)

	if len(m.IsCreationRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsCreationRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestMock.IsCreationRequest.")
			return
		}

		result := m.IsCreationRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsCreationRequest")
			return
		}

		r = result.r

		return
	}

	if m.IsCreationRequestMock.mainExpectation != nil {

		result := m.IsCreationRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsCreationRequest")
		}

		r = result.r

		return
	}

	if m.IsCreationRequestFunc == nil {
		m.t.Fatalf("Unexpected call to RequestMock.IsCreationRequest.")
		return
	}

	return m.IsCreationRequestFunc()
}

//IsCreationRequestMinimockCounter returns a count of RequestMock.IsCreationRequestFunc invocations
func (m *RequestMock) IsCreationRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsCreationRequestCounter)
}

//IsCreationRequestMinimockPreCounter returns the value of RequestMock.IsCreationRequest invocations
func (m *RequestMock) IsCreationRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsCreationRequestPreCounter)
}

//IsCreationRequestFinished returns true if mock invocations count is ok
func (m *RequestMock) IsCreationRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsCreationRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsCreationRequestCounter) == uint64(len(m.IsCreationRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsCreationRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsCreationRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsCreationRequestFunc != nil {
		return atomic.LoadUint64(&m.IsCreationRequestCounter) > 0
	}

	return true
}

type mRequestMockIsDetached struct {
	mock              *RequestMock
	mainExpectation   *RequestMockIsDetachedExpectation
	expectationSeries []*RequestMockIsDetachedExpectation
}

type RequestMockIsDetachedExpectation struct {
	result *RequestMockIsDetachedResult
}

type RequestMockIsDetachedResult struct {
	r bool
}

//Expect specifies that invocation of Request.IsDetached is expected from 1 to Infinity times
func (m *mRequestMockIsDetached) Expect() *mRequestMockIsDetached {
	m.mock.IsDetachedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsDetachedExpectation{}
	}

	return m
}

//Return specifies results of invocation of Request.IsDetached
func (m *mRequestMockIsDetached) Return(r bool) *RequestMock {
	m.mock.IsDetachedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsDetachedExpectation{}
	}
	m.mainExpectation.result = &RequestMockIsDetachedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Request.IsDetached is expected once
func (m *mRequestMockIsDetached) ExpectOnce() *RequestMockIsDetachedExpectation {
	m.mock.IsDetachedFunc = nil
	m.mainExpectation = nil

	expectation := &RequestMockIsDetachedExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestMockIsDetachedExpectation) Return(r bool) {
	e.result = &RequestMockIsDetachedResult{r}
}

//Set uses given function f as a mock of Request.IsDetached method
func (m *mRequestMockIsDetached) Set(f func() (r bool)) *RequestMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsDetachedFunc = f
	return m.mock
}

//IsDetached implements github.com/insolar/insolar/insolar/record.Request interface
func (m *RequestMock) IsDetached() (r bool) {
	counter := atomic.AddUint64(&m.IsDetachedPreCounter, 1)
	defer atomic.AddUint64(&m.IsDetachedCounter, 1)

	if len(m.IsDetachedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsDetachedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestMock.IsDetached.")
			return
		}

		result := m.IsDetachedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsDetached")
			return
		}

		r = result.r

		return
	}

	if m.IsDetachedMock.mainExpectation != nil {

		result := m.IsDetachedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsDetached")
		}

		r = result.r

		return
	}

	if m.IsDetachedFunc == nil {
		m.t.Fatalf("Unexpected call to RequestMock.IsDetached.")
		return
	}

	return m.IsDetachedFunc()
}

//IsDetachedMinimockCounter returns a count of RequestMock.IsDetachedFunc invocations
func (m *RequestMock) IsDetachedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsDetachedCounter)
}

//IsDetachedMinimockPreCounter returns the value of RequestMock.IsDetached invocations
func (m *RequestMock) IsDetachedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsDetachedPreCounter)
}

//IsDetachedFinished returns true if mock invocations count is ok
func (m *RequestMock) IsDetachedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsDetachedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsDetachedCounter) == uint64(len(m.IsDetachedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsDetachedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsDetachedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsDetachedFunc != nil {
		return atomic.LoadUint64(&m.IsDetachedCounter) > 0
	}

	return true
}

type mRequestMockIsEmptyAPINode struct {
	mock              *RequestMock
	mainExpectation   *RequestMockIsEmptyAPINodeExpectation
	expectationSeries []*RequestMockIsEmptyAPINodeExpectation
}

type RequestMockIsEmptyAPINodeExpectation struct {
	result *RequestMockIsEmptyAPINodeResult
}

type RequestMockIsEmptyAPINodeResult struct {
	r bool
}

//Expect specifies that invocation of Request.IsEmptyAPINode is expected from 1 to Infinity times
func (m *mRequestMockIsEmptyAPINode) Expect() *mRequestMockIsEmptyAPINode {
	m.mock.IsEmptyAPINodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsEmptyAPINodeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Request.IsEmptyAPINode
func (m *mRequestMockIsEmptyAPINode) Return(r bool) *RequestMock {
	m.mock.IsEmptyAPINodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockIsEmptyAPINodeExpectation{}
	}
	m.mainExpectation.result = &RequestMockIsEmptyAPINodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Request.IsEmptyAPINode is expected once
func (m *mRequestMockIsEmptyAPINode) ExpectOnce() *RequestMockIsEmptyAPINodeExpectation {
	m.mock.IsEmptyAPINodeFunc = nil
	m.mainExpectation = nil

	expectation := &RequestMockIsEmptyAPINodeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestMockIsEmptyAPINodeExpectation) Return(r bool) {
	e.result = &RequestMockIsEmptyAPINodeResult{r}
}

//Set uses given function f as a mock of Request.IsEmptyAPINode method
func (m *mRequestMockIsEmptyAPINode) Set(f func() (r bool)) *RequestMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsEmptyAPINodeFunc = f
	return m.mock
}

//IsEmptyAPINode implements github.com/insolar/insolar/insolar/record.Request interface
func (m *RequestMock) IsEmptyAPINode() (r bool) {
	counter := atomic.AddUint64(&m.IsEmptyAPINodePreCounter, 1)
	defer atomic.AddUint64(&m.IsEmptyAPINodeCounter, 1)

	if len(m.IsEmptyAPINodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsEmptyAPINodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestMock.IsEmptyAPINode.")
			return
		}

		result := m.IsEmptyAPINodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsEmptyAPINode")
			return
		}

		r = result.r

		return
	}

	if m.IsEmptyAPINodeMock.mainExpectation != nil {

		result := m.IsEmptyAPINodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.IsEmptyAPINode")
		}

		r = result.r

		return
	}

	if m.IsEmptyAPINodeFunc == nil {
		m.t.Fatalf("Unexpected call to RequestMock.IsEmptyAPINode.")
		return
	}

	return m.IsEmptyAPINodeFunc()
}

//IsEmptyAPINodeMinimockCounter returns a count of RequestMock.IsEmptyAPINodeFunc invocations
func (m *RequestMock) IsEmptyAPINodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsEmptyAPINodeCounter)
}

//IsEmptyAPINodeMinimockPreCounter returns the value of RequestMock.IsEmptyAPINode invocations
func (m *RequestMock) IsEmptyAPINodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsEmptyAPINodePreCounter)
}

//IsEmptyAPINodeFinished returns true if mock invocations count is ok
func (m *RequestMock) IsEmptyAPINodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsEmptyAPINodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsEmptyAPINodeCounter) == uint64(len(m.IsEmptyAPINodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsEmptyAPINodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsEmptyAPINodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsEmptyAPINodeFunc != nil {
		return atomic.LoadUint64(&m.IsEmptyAPINodeCounter) > 0
	}

	return true
}

type mRequestMockReasonRef struct {
	mock              *RequestMock
	mainExpectation   *RequestMockReasonRefExpectation
	expectationSeries []*RequestMockReasonRefExpectation
}

type RequestMockReasonRefExpectation struct {
	result *RequestMockReasonRefResult
}

type RequestMockReasonRefResult struct {
	r insolar.Reference
}

//Expect specifies that invocation of Request.ReasonRef is expected from 1 to Infinity times
func (m *mRequestMockReasonRef) Expect() *mRequestMockReasonRef {
	m.mock.ReasonRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockReasonRefExpectation{}
	}

	return m
}

//Return specifies results of invocation of Request.ReasonRef
func (m *mRequestMockReasonRef) Return(r insolar.Reference) *RequestMock {
	m.mock.ReasonRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestMockReasonRefExpectation{}
	}
	m.mainExpectation.result = &RequestMockReasonRefResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Request.ReasonRef is expected once
func (m *mRequestMockReasonRef) ExpectOnce() *RequestMockReasonRefExpectation {
	m.mock.ReasonRefFunc = nil
	m.mainExpectation = nil

	expectation := &RequestMockReasonRefExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestMockReasonRefExpectation) Return(r insolar.Reference) {
	e.result = &RequestMockReasonRefResult{r}
}

//Set uses given function f as a mock of Request.ReasonRef method
func (m *mRequestMockReasonRef) Set(f func() (r insolar.Reference)) *RequestMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReasonRefFunc = f
	return m.mock
}

//ReasonRef implements github.com/insolar/insolar/insolar/record.Request interface
func (m *RequestMock) ReasonRef() (r insolar.Reference) {
	counter := atomic.AddUint64(&m.ReasonRefPreCounter, 1)
	defer atomic.AddUint64(&m.ReasonRefCounter, 1)

	if len(m.ReasonRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReasonRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestMock.ReasonRef.")
			return
		}

		result := m.ReasonRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.ReasonRef")
			return
		}

		r = result.r

		return
	}

	if m.ReasonRefMock.mainExpectation != nil {

		result := m.ReasonRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestMock.ReasonRef")
		}

		r = result.r

		return
	}

	if m.ReasonRefFunc == nil {
		m.t.Fatalf("Unexpected call to RequestMock.ReasonRef.")
		return
	}

	return m.ReasonRefFunc()
}

//ReasonRefMinimockCounter returns a count of RequestMock.ReasonRefFunc invocations
func (m *RequestMock) ReasonRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReasonRefCounter)
}

//ReasonRefMinimockPreCounter returns the value of RequestMock.ReasonRef invocations
func (m *RequestMock) ReasonRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReasonRefPreCounter)
}

//ReasonRefFinished returns true if mock invocations count is ok
func (m *RequestMock) ReasonRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReasonRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReasonRefCounter) == uint64(len(m.ReasonRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReasonRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReasonRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReasonRefFunc != nil {
		return atomic.LoadUint64(&m.ReasonRefCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestMock) ValidateCallCounters() {

	if !m.AffinityRefFinished() {
		m.t.Fatal("Expected call to RequestMock.AffinityRef")
	}

	if !m.GetCallTypeFinished() {
		m.t.Fatal("Expected call to RequestMock.GetCallType")
	}

	if !m.IsAPIRequestFinished() {
		m.t.Fatal("Expected call to RequestMock.IsAPIRequest")
	}

	if !m.IsCreationRequestFinished() {
		m.t.Fatal("Expected call to RequestMock.IsCreationRequest")
	}

	if !m.IsDetachedFinished() {
		m.t.Fatal("Expected call to RequestMock.IsDetached")
	}

	if !m.IsEmptyAPINodeFinished() {
		m.t.Fatal("Expected call to RequestMock.IsEmptyAPINode")
	}

	if !m.ReasonRefFinished() {
		m.t.Fatal("Expected call to RequestMock.ReasonRef")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RequestMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RequestMock) MinimockFinish() {

	if !m.AffinityRefFinished() {
		m.t.Fatal("Expected call to RequestMock.AffinityRef")
	}

	if !m.GetCallTypeFinished() {
		m.t.Fatal("Expected call to RequestMock.GetCallType")
	}

	if !m.IsAPIRequestFinished() {
		m.t.Fatal("Expected call to RequestMock.IsAPIRequest")
	}

	if !m.IsCreationRequestFinished() {
		m.t.Fatal("Expected call to RequestMock.IsCreationRequest")
	}

	if !m.IsDetachedFinished() {
		m.t.Fatal("Expected call to RequestMock.IsDetached")
	}

	if !m.IsEmptyAPINodeFinished() {
		m.t.Fatal("Expected call to RequestMock.IsEmptyAPINode")
	}

	if !m.ReasonRefFinished() {
		m.t.Fatal("Expected call to RequestMock.ReasonRef")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RequestMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RequestMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AffinityRefFinished()
		ok = ok && m.GetCallTypeFinished()
		ok = ok && m.IsAPIRequestFinished()
		ok = ok && m.IsCreationRequestFinished()
		ok = ok && m.IsDetachedFinished()
		ok = ok && m.IsEmptyAPINodeFinished()
		ok = ok && m.ReasonRefFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AffinityRefFinished() {
				m.t.Error("Expected call to RequestMock.AffinityRef")
			}

			if !m.GetCallTypeFinished() {
				m.t.Error("Expected call to RequestMock.GetCallType")
			}

			if !m.IsAPIRequestFinished() {
				m.t.Error("Expected call to RequestMock.IsAPIRequest")
			}

			if !m.IsCreationRequestFinished() {
				m.t.Error("Expected call to RequestMock.IsCreationRequest")
			}

			if !m.IsDetachedFinished() {
				m.t.Error("Expected call to RequestMock.IsDetached")
			}

			if !m.IsEmptyAPINodeFinished() {
				m.t.Error("Expected call to RequestMock.IsEmptyAPINode")
			}

			if !m.ReasonRefFinished() {
				m.t.Error("Expected call to RequestMock.ReasonRef")
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
func (m *RequestMock) AllMocksCalled() bool {

	if !m.AffinityRefFinished() {
		return false
	}

	if !m.GetCallTypeFinished() {
		return false
	}

	if !m.IsAPIRequestFinished() {
		return false
	}

	if !m.IsCreationRequestFinished() {
		return false
	}

	if !m.IsDetachedFinished() {
		return false
	}

	if !m.IsEmptyAPINodeFinished() {
		return false
	}

	if !m.ReasonRefFinished() {
		return false
	}

	return true
}

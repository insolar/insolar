package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RequestChecker" can be found in github.com/insolar/insolar/ledger/light/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//RequestCheckerMock implements github.com/insolar/insolar/ledger/light/executor.RequestChecker
type RequestCheckerMock struct {
	t minimock.Tester

	CheckRequestFunc       func(p context.Context, p1 insolar.ID, p2 record.Request) (r error)
	CheckRequestCounter    uint64
	CheckRequestPreCounter uint64
	CheckRequestMock       mRequestCheckerMockCheckRequest
}

//NewRequestCheckerMock returns a mock for github.com/insolar/insolar/ledger/light/executor.RequestChecker
func NewRequestCheckerMock(t minimock.Tester) *RequestCheckerMock {
	m := &RequestCheckerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CheckRequestMock = mRequestCheckerMockCheckRequest{mock: m}

	return m
}

type mRequestCheckerMockCheckRequest struct {
	mock              *RequestCheckerMock
	mainExpectation   *RequestCheckerMockCheckRequestExpectation
	expectationSeries []*RequestCheckerMockCheckRequestExpectation
}

type RequestCheckerMockCheckRequestExpectation struct {
	input  *RequestCheckerMockCheckRequestInput
	result *RequestCheckerMockCheckRequestResult
}

type RequestCheckerMockCheckRequestInput struct {
	p  context.Context
	p1 insolar.ID
	p2 record.Request
}

type RequestCheckerMockCheckRequestResult struct {
	r error
}

//Expect specifies that invocation of RequestChecker.CheckRequest is expected from 1 to Infinity times
func (m *mRequestCheckerMockCheckRequest) Expect(p context.Context, p1 insolar.ID, p2 record.Request) *mRequestCheckerMockCheckRequest {
	m.mock.CheckRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestCheckerMockCheckRequestExpectation{}
	}
	m.mainExpectation.input = &RequestCheckerMockCheckRequestInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RequestChecker.CheckRequest
func (m *mRequestCheckerMockCheckRequest) Return(r error) *RequestCheckerMock {
	m.mock.CheckRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestCheckerMockCheckRequestExpectation{}
	}
	m.mainExpectation.result = &RequestCheckerMockCheckRequestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RequestChecker.CheckRequest is expected once
func (m *mRequestCheckerMockCheckRequest) ExpectOnce(p context.Context, p1 insolar.ID, p2 record.Request) *RequestCheckerMockCheckRequestExpectation {
	m.mock.CheckRequestFunc = nil
	m.mainExpectation = nil

	expectation := &RequestCheckerMockCheckRequestExpectation{}
	expectation.input = &RequestCheckerMockCheckRequestInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestCheckerMockCheckRequestExpectation) Return(r error) {
	e.result = &RequestCheckerMockCheckRequestResult{r}
}

//Set uses given function f as a mock of RequestChecker.CheckRequest method
func (m *mRequestCheckerMockCheckRequest) Set(f func(p context.Context, p1 insolar.ID, p2 record.Request) (r error)) *RequestCheckerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CheckRequestFunc = f
	return m.mock
}

//CheckRequest implements github.com/insolar/insolar/ledger/light/executor.RequestChecker interface
func (m *RequestCheckerMock) CheckRequest(p context.Context, p1 insolar.ID, p2 record.Request) (r error) {
	counter := atomic.AddUint64(&m.CheckRequestPreCounter, 1)
	defer atomic.AddUint64(&m.CheckRequestCounter, 1)

	if len(m.CheckRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CheckRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestCheckerMock.CheckRequest. %v %v %v", p, p1, p2)
			return
		}

		input := m.CheckRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RequestCheckerMockCheckRequestInput{p, p1, p2}, "RequestChecker.CheckRequest got unexpected parameters")

		result := m.CheckRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestCheckerMock.CheckRequest")
			return
		}

		r = result.r

		return
	}

	if m.CheckRequestMock.mainExpectation != nil {

		input := m.CheckRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RequestCheckerMockCheckRequestInput{p, p1, p2}, "RequestChecker.CheckRequest got unexpected parameters")
		}

		result := m.CheckRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestCheckerMock.CheckRequest")
		}

		r = result.r

		return
	}

	if m.CheckRequestFunc == nil {
		m.t.Fatalf("Unexpected call to RequestCheckerMock.CheckRequest. %v %v %v", p, p1, p2)
		return
	}

	return m.CheckRequestFunc(p, p1, p2)
}

//CheckRequestMinimockCounter returns a count of RequestCheckerMock.CheckRequestFunc invocations
func (m *RequestCheckerMock) CheckRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CheckRequestCounter)
}

//CheckRequestMinimockPreCounter returns the value of RequestCheckerMock.CheckRequest invocations
func (m *RequestCheckerMock) CheckRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CheckRequestPreCounter)
}

//CheckRequestFinished returns true if mock invocations count is ok
func (m *RequestCheckerMock) CheckRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CheckRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CheckRequestCounter) == uint64(len(m.CheckRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CheckRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CheckRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CheckRequestFunc != nil {
		return atomic.LoadUint64(&m.CheckRequestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestCheckerMock) ValidateCallCounters() {

	if !m.CheckRequestFinished() {
		m.t.Fatal("Expected call to RequestCheckerMock.CheckRequest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestCheckerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RequestCheckerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RequestCheckerMock) MinimockFinish() {

	if !m.CheckRequestFinished() {
		m.t.Fatal("Expected call to RequestCheckerMock.CheckRequest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RequestCheckerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RequestCheckerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CheckRequestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CheckRequestFinished() {
				m.t.Error("Expected call to RequestCheckerMock.CheckRequest")
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
func (m *RequestCheckerMock) AllMocksCalled() bool {

	if !m.CheckRequestFinished() {
		return false
	}

	return true
}

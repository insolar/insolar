package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RequestsFetcher" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//RequestsFetcherMock implements github.com/insolar/insolar/logicrunner.RequestsFetcher
type RequestsFetcherMock struct {
	t minimock.Tester

	AbortFunc       func(p context.Context)
	AbortCounter    uint64
	AbortPreCounter uint64
	AbortMock       mRequestsFetcherMockAbort

	FetchPendingsFunc       func(p context.Context)
	FetchPendingsCounter    uint64
	FetchPendingsPreCounter uint64
	FetchPendingsMock       mRequestsFetcherMockFetchPendings
}

//NewRequestsFetcherMock returns a mock for github.com/insolar/insolar/logicrunner.RequestsFetcher
func NewRequestsFetcherMock(t minimock.Tester) *RequestsFetcherMock {
	m := &RequestsFetcherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AbortMock = mRequestsFetcherMockAbort{mock: m}
	m.FetchPendingsMock = mRequestsFetcherMockFetchPendings{mock: m}

	return m
}

type mRequestsFetcherMockAbort struct {
	mock              *RequestsFetcherMock
	mainExpectation   *RequestsFetcherMockAbortExpectation
	expectationSeries []*RequestsFetcherMockAbortExpectation
}

type RequestsFetcherMockAbortExpectation struct {
	input *RequestsFetcherMockAbortInput
}

type RequestsFetcherMockAbortInput struct {
	p context.Context
}

//Expect specifies that invocation of RequestsFetcher.Abort is expected from 1 to Infinity times
func (m *mRequestsFetcherMockAbort) Expect(p context.Context) *mRequestsFetcherMockAbort {
	m.mock.AbortFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsFetcherMockAbortExpectation{}
	}
	m.mainExpectation.input = &RequestsFetcherMockAbortInput{p}
	return m
}

//Return specifies results of invocation of RequestsFetcher.Abort
func (m *mRequestsFetcherMockAbort) Return() *RequestsFetcherMock {
	m.mock.AbortFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsFetcherMockAbortExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RequestsFetcher.Abort is expected once
func (m *mRequestsFetcherMockAbort) ExpectOnce(p context.Context) *RequestsFetcherMockAbortExpectation {
	m.mock.AbortFunc = nil
	m.mainExpectation = nil

	expectation := &RequestsFetcherMockAbortExpectation{}
	expectation.input = &RequestsFetcherMockAbortInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RequestsFetcher.Abort method
func (m *mRequestsFetcherMockAbort) Set(f func(p context.Context)) *RequestsFetcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AbortFunc = f
	return m.mock
}

//Abort implements github.com/insolar/insolar/logicrunner.RequestsFetcher interface
func (m *RequestsFetcherMock) Abort(p context.Context) {
	counter := atomic.AddUint64(&m.AbortPreCounter, 1)
	defer atomic.AddUint64(&m.AbortCounter, 1)

	if len(m.AbortMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AbortMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestsFetcherMock.Abort. %v", p)
			return
		}

		input := m.AbortMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RequestsFetcherMockAbortInput{p}, "RequestsFetcher.Abort got unexpected parameters")

		return
	}

	if m.AbortMock.mainExpectation != nil {

		input := m.AbortMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RequestsFetcherMockAbortInput{p}, "RequestsFetcher.Abort got unexpected parameters")
		}

		return
	}

	if m.AbortFunc == nil {
		m.t.Fatalf("Unexpected call to RequestsFetcherMock.Abort. %v", p)
		return
	}

	m.AbortFunc(p)
}

//AbortMinimockCounter returns a count of RequestsFetcherMock.AbortFunc invocations
func (m *RequestsFetcherMock) AbortMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AbortCounter)
}

//AbortMinimockPreCounter returns the value of RequestsFetcherMock.Abort invocations
func (m *RequestsFetcherMock) AbortMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AbortPreCounter)
}

//AbortFinished returns true if mock invocations count is ok
func (m *RequestsFetcherMock) AbortFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AbortMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AbortCounter) == uint64(len(m.AbortMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AbortMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AbortCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AbortFunc != nil {
		return atomic.LoadUint64(&m.AbortCounter) > 0
	}

	return true
}

type mRequestsFetcherMockFetchPendings struct {
	mock              *RequestsFetcherMock
	mainExpectation   *RequestsFetcherMockFetchPendingsExpectation
	expectationSeries []*RequestsFetcherMockFetchPendingsExpectation
}

type RequestsFetcherMockFetchPendingsExpectation struct {
	input *RequestsFetcherMockFetchPendingsInput
}

type RequestsFetcherMockFetchPendingsInput struct {
	p context.Context
}

//Expect specifies that invocation of RequestsFetcher.FetchPendings is expected from 1 to Infinity times
func (m *mRequestsFetcherMockFetchPendings) Expect(p context.Context) *mRequestsFetcherMockFetchPendings {
	m.mock.FetchPendingsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsFetcherMockFetchPendingsExpectation{}
	}
	m.mainExpectation.input = &RequestsFetcherMockFetchPendingsInput{p}
	return m
}

//Return specifies results of invocation of RequestsFetcher.FetchPendings
func (m *mRequestsFetcherMockFetchPendings) Return() *RequestsFetcherMock {
	m.mock.FetchPendingsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsFetcherMockFetchPendingsExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RequestsFetcher.FetchPendings is expected once
func (m *mRequestsFetcherMockFetchPendings) ExpectOnce(p context.Context) *RequestsFetcherMockFetchPendingsExpectation {
	m.mock.FetchPendingsFunc = nil
	m.mainExpectation = nil

	expectation := &RequestsFetcherMockFetchPendingsExpectation{}
	expectation.input = &RequestsFetcherMockFetchPendingsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RequestsFetcher.FetchPendings method
func (m *mRequestsFetcherMockFetchPendings) Set(f func(p context.Context)) *RequestsFetcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FetchPendingsFunc = f
	return m.mock
}

//FetchPendings implements github.com/insolar/insolar/logicrunner.RequestsFetcher interface
func (m *RequestsFetcherMock) FetchPendings(p context.Context) {
	counter := atomic.AddUint64(&m.FetchPendingsPreCounter, 1)
	defer atomic.AddUint64(&m.FetchPendingsCounter, 1)

	if len(m.FetchPendingsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FetchPendingsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestsFetcherMock.FetchPendings. %v", p)
			return
		}

		input := m.FetchPendingsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RequestsFetcherMockFetchPendingsInput{p}, "RequestsFetcher.FetchPendings got unexpected parameters")

		return
	}

	if m.FetchPendingsMock.mainExpectation != nil {

		input := m.FetchPendingsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RequestsFetcherMockFetchPendingsInput{p}, "RequestsFetcher.FetchPendings got unexpected parameters")
		}

		return
	}

	if m.FetchPendingsFunc == nil {
		m.t.Fatalf("Unexpected call to RequestsFetcherMock.FetchPendings. %v", p)
		return
	}

	m.FetchPendingsFunc(p)
}

//FetchPendingsMinimockCounter returns a count of RequestsFetcherMock.FetchPendingsFunc invocations
func (m *RequestsFetcherMock) FetchPendingsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FetchPendingsCounter)
}

//FetchPendingsMinimockPreCounter returns the value of RequestsFetcherMock.FetchPendings invocations
func (m *RequestsFetcherMock) FetchPendingsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FetchPendingsPreCounter)
}

//FetchPendingsFinished returns true if mock invocations count is ok
func (m *RequestsFetcherMock) FetchPendingsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FetchPendingsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FetchPendingsCounter) == uint64(len(m.FetchPendingsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FetchPendingsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FetchPendingsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FetchPendingsFunc != nil {
		return atomic.LoadUint64(&m.FetchPendingsCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestsFetcherMock) ValidateCallCounters() {

	if !m.AbortFinished() {
		m.t.Fatal("Expected call to RequestsFetcherMock.Abort")
	}

	if !m.FetchPendingsFinished() {
		m.t.Fatal("Expected call to RequestsFetcherMock.FetchPendings")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestsFetcherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RequestsFetcherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RequestsFetcherMock) MinimockFinish() {

	if !m.AbortFinished() {
		m.t.Fatal("Expected call to RequestsFetcherMock.Abort")
	}

	if !m.FetchPendingsFinished() {
		m.t.Fatal("Expected call to RequestsFetcherMock.FetchPendings")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RequestsFetcherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RequestsFetcherMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AbortFinished()
		ok = ok && m.FetchPendingsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AbortFinished() {
				m.t.Error("Expected call to RequestsFetcherMock.Abort")
			}

			if !m.FetchPendingsFinished() {
				m.t.Error("Expected call to RequestsFetcherMock.FetchPendings")
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
func (m *RequestsFetcherMock) AllMocksCalled() bool {

	if !m.AbortFinished() {
		return false
	}

	if !m.FetchPendingsFinished() {
		return false
	}

	return true
}

package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ResultMatcher" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	message "github.com/insolar/insolar/insolar/message"

	testify_assert "github.com/stretchr/testify/assert"
)

//ResultMatcherMock implements github.com/insolar/insolar/logicrunner.ResultMatcher
type ResultMatcherMock struct {
	t minimock.Tester

	AddStillExecutionFunc       func(p context.Context, p1 *message.StillExecuting)
	AddStillExecutionCounter    uint64
	AddStillExecutionPreCounter uint64
	AddStillExecutionMock       mResultMatcherMockAddStillExecution

	AddUnwantedResponseFunc       func(p context.Context, p1 *message.ReturnResults) (r error)
	AddUnwantedResponseCounter    uint64
	AddUnwantedResponsePreCounter uint64
	AddUnwantedResponseMock       mResultMatcherMockAddUnwantedResponse

	ClearFunc       func()
	ClearCounter    uint64
	ClearPreCounter uint64
	ClearMock       mResultMatcherMockClear
}

//NewResultMatcherMock returns a mock for github.com/insolar/insolar/logicrunner.ResultMatcher
func NewResultMatcherMock(t minimock.Tester) *ResultMatcherMock {
	m := &ResultMatcherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddStillExecutionMock = mResultMatcherMockAddStillExecution{mock: m}
	m.AddUnwantedResponseMock = mResultMatcherMockAddUnwantedResponse{mock: m}
	m.ClearMock = mResultMatcherMockClear{mock: m}

	return m
}

type mResultMatcherMockAddStillExecution struct {
	mock              *ResultMatcherMock
	mainExpectation   *ResultMatcherMockAddStillExecutionExpectation
	expectationSeries []*ResultMatcherMockAddStillExecutionExpectation
}

type ResultMatcherMockAddStillExecutionExpectation struct {
	input *ResultMatcherMockAddStillExecutionInput
}

type ResultMatcherMockAddStillExecutionInput struct {
	p  context.Context
	p1 *message.StillExecuting
}

//Expect specifies that invocation of ResultMatcher.AddStillExecution is expected from 1 to Infinity times
func (m *mResultMatcherMockAddStillExecution) Expect(p context.Context, p1 *message.StillExecuting) *mResultMatcherMockAddStillExecution {
	m.mock.AddStillExecutionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResultMatcherMockAddStillExecutionExpectation{}
	}
	m.mainExpectation.input = &ResultMatcherMockAddStillExecutionInput{p, p1}
	return m
}

//Return specifies results of invocation of ResultMatcher.AddStillExecution
func (m *mResultMatcherMockAddStillExecution) Return() *ResultMatcherMock {
	m.mock.AddStillExecutionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResultMatcherMockAddStillExecutionExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ResultMatcher.AddStillExecution is expected once
func (m *mResultMatcherMockAddStillExecution) ExpectOnce(p context.Context, p1 *message.StillExecuting) *ResultMatcherMockAddStillExecutionExpectation {
	m.mock.AddStillExecutionFunc = nil
	m.mainExpectation = nil

	expectation := &ResultMatcherMockAddStillExecutionExpectation{}
	expectation.input = &ResultMatcherMockAddStillExecutionInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ResultMatcher.AddStillExecution method
func (m *mResultMatcherMockAddStillExecution) Set(f func(p context.Context, p1 *message.StillExecuting)) *ResultMatcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddStillExecutionFunc = f
	return m.mock
}

//AddStillExecution implements github.com/insolar/insolar/logicrunner.ResultMatcher interface
func (m *ResultMatcherMock) AddStillExecution(p context.Context, p1 *message.StillExecuting) {
	counter := atomic.AddUint64(&m.AddStillExecutionPreCounter, 1)
	defer atomic.AddUint64(&m.AddStillExecutionCounter, 1)

	if len(m.AddStillExecutionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddStillExecutionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ResultMatcherMock.AddStillExecution. %v %v", p, p1)
			return
		}

		input := m.AddStillExecutionMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ResultMatcherMockAddStillExecutionInput{p, p1}, "ResultMatcher.AddStillExecution got unexpected parameters")

		return
	}

	if m.AddStillExecutionMock.mainExpectation != nil {

		input := m.AddStillExecutionMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ResultMatcherMockAddStillExecutionInput{p, p1}, "ResultMatcher.AddStillExecution got unexpected parameters")
		}

		return
	}

	if m.AddStillExecutionFunc == nil {
		m.t.Fatalf("Unexpected call to ResultMatcherMock.AddStillExecution. %v %v", p, p1)
		return
	}

	m.AddStillExecutionFunc(p, p1)
}

//AddStillExecutionMinimockCounter returns a count of ResultMatcherMock.AddStillExecutionFunc invocations
func (m *ResultMatcherMock) AddStillExecutionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddStillExecutionCounter)
}

//AddStillExecutionMinimockPreCounter returns the value of ResultMatcherMock.AddStillExecution invocations
func (m *ResultMatcherMock) AddStillExecutionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddStillExecutionPreCounter)
}

//AddStillExecutionFinished returns true if mock invocations count is ok
func (m *ResultMatcherMock) AddStillExecutionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddStillExecutionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddStillExecutionCounter) == uint64(len(m.AddStillExecutionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddStillExecutionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddStillExecutionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddStillExecutionFunc != nil {
		return atomic.LoadUint64(&m.AddStillExecutionCounter) > 0
	}

	return true
}

type mResultMatcherMockAddUnwantedResponse struct {
	mock              *ResultMatcherMock
	mainExpectation   *ResultMatcherMockAddUnwantedResponseExpectation
	expectationSeries []*ResultMatcherMockAddUnwantedResponseExpectation
}

type ResultMatcherMockAddUnwantedResponseExpectation struct {
	input  *ResultMatcherMockAddUnwantedResponseInput
	result *ResultMatcherMockAddUnwantedResponseResult
}

type ResultMatcherMockAddUnwantedResponseInput struct {
	p  context.Context
	p1 *message.ReturnResults
}

type ResultMatcherMockAddUnwantedResponseResult struct {
	r error
}

//Expect specifies that invocation of ResultMatcher.AddUnwantedResponse is expected from 1 to Infinity times
func (m *mResultMatcherMockAddUnwantedResponse) Expect(p context.Context, p1 *message.ReturnResults) *mResultMatcherMockAddUnwantedResponse {
	m.mock.AddUnwantedResponseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResultMatcherMockAddUnwantedResponseExpectation{}
	}
	m.mainExpectation.input = &ResultMatcherMockAddUnwantedResponseInput{p, p1}
	return m
}

//Return specifies results of invocation of ResultMatcher.AddUnwantedResponse
func (m *mResultMatcherMockAddUnwantedResponse) Return(r error) *ResultMatcherMock {
	m.mock.AddUnwantedResponseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResultMatcherMockAddUnwantedResponseExpectation{}
	}
	m.mainExpectation.result = &ResultMatcherMockAddUnwantedResponseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ResultMatcher.AddUnwantedResponse is expected once
func (m *mResultMatcherMockAddUnwantedResponse) ExpectOnce(p context.Context, p1 *message.ReturnResults) *ResultMatcherMockAddUnwantedResponseExpectation {
	m.mock.AddUnwantedResponseFunc = nil
	m.mainExpectation = nil

	expectation := &ResultMatcherMockAddUnwantedResponseExpectation{}
	expectation.input = &ResultMatcherMockAddUnwantedResponseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ResultMatcherMockAddUnwantedResponseExpectation) Return(r error) {
	e.result = &ResultMatcherMockAddUnwantedResponseResult{r}
}

//Set uses given function f as a mock of ResultMatcher.AddUnwantedResponse method
func (m *mResultMatcherMockAddUnwantedResponse) Set(f func(p context.Context, p1 *message.ReturnResults) (r error)) *ResultMatcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddUnwantedResponseFunc = f
	return m.mock
}

//AddUnwantedResponse implements github.com/insolar/insolar/logicrunner.ResultMatcher interface
func (m *ResultMatcherMock) AddUnwantedResponse(p context.Context, p1 *message.ReturnResults) (r error) {
	counter := atomic.AddUint64(&m.AddUnwantedResponsePreCounter, 1)
	defer atomic.AddUint64(&m.AddUnwantedResponseCounter, 1)

	if len(m.AddUnwantedResponseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddUnwantedResponseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ResultMatcherMock.AddUnwantedResponse. %v %v", p, p1)
			return
		}

		input := m.AddUnwantedResponseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ResultMatcherMockAddUnwantedResponseInput{p, p1}, "ResultMatcher.AddUnwantedResponse got unexpected parameters")

		result := m.AddUnwantedResponseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ResultMatcherMock.AddUnwantedResponse")
			return
		}

		r = result.r

		return
	}

	if m.AddUnwantedResponseMock.mainExpectation != nil {

		input := m.AddUnwantedResponseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ResultMatcherMockAddUnwantedResponseInput{p, p1}, "ResultMatcher.AddUnwantedResponse got unexpected parameters")
		}

		result := m.AddUnwantedResponseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ResultMatcherMock.AddUnwantedResponse")
		}

		r = result.r

		return
	}

	if m.AddUnwantedResponseFunc == nil {
		m.t.Fatalf("Unexpected call to ResultMatcherMock.AddUnwantedResponse. %v %v", p, p1)
		return
	}

	return m.AddUnwantedResponseFunc(p, p1)
}

//AddUnwantedResponseMinimockCounter returns a count of ResultMatcherMock.AddUnwantedResponseFunc invocations
func (m *ResultMatcherMock) AddUnwantedResponseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddUnwantedResponseCounter)
}

//AddUnwantedResponseMinimockPreCounter returns the value of ResultMatcherMock.AddUnwantedResponse invocations
func (m *ResultMatcherMock) AddUnwantedResponseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddUnwantedResponsePreCounter)
}

//AddUnwantedResponseFinished returns true if mock invocations count is ok
func (m *ResultMatcherMock) AddUnwantedResponseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddUnwantedResponseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddUnwantedResponseCounter) == uint64(len(m.AddUnwantedResponseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddUnwantedResponseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddUnwantedResponseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddUnwantedResponseFunc != nil {
		return atomic.LoadUint64(&m.AddUnwantedResponseCounter) > 0
	}

	return true
}

type mResultMatcherMockClear struct {
	mock              *ResultMatcherMock
	mainExpectation   *ResultMatcherMockClearExpectation
	expectationSeries []*ResultMatcherMockClearExpectation
}

type ResultMatcherMockClearExpectation struct {
}

//Expect specifies that invocation of ResultMatcher.Clear is expected from 1 to Infinity times
func (m *mResultMatcherMockClear) Expect() *mResultMatcherMockClear {
	m.mock.ClearFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResultMatcherMockClearExpectation{}
	}

	return m
}

//Return specifies results of invocation of ResultMatcher.Clear
func (m *mResultMatcherMockClear) Return() *ResultMatcherMock {
	m.mock.ClearFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResultMatcherMockClearExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ResultMatcher.Clear is expected once
func (m *mResultMatcherMockClear) ExpectOnce() *ResultMatcherMockClearExpectation {
	m.mock.ClearFunc = nil
	m.mainExpectation = nil

	expectation := &ResultMatcherMockClearExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ResultMatcher.Clear method
func (m *mResultMatcherMockClear) Set(f func()) *ResultMatcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ClearFunc = f
	return m.mock
}

//Clear implements github.com/insolar/insolar/logicrunner.ResultMatcher interface
func (m *ResultMatcherMock) Clear() {
	counter := atomic.AddUint64(&m.ClearPreCounter, 1)
	defer atomic.AddUint64(&m.ClearCounter, 1)

	if len(m.ClearMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ClearMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ResultMatcherMock.Clear.")
			return
		}

		return
	}

	if m.ClearMock.mainExpectation != nil {

		return
	}

	if m.ClearFunc == nil {
		m.t.Fatalf("Unexpected call to ResultMatcherMock.Clear.")
		return
	}

	m.ClearFunc()
}

//ClearMinimockCounter returns a count of ResultMatcherMock.ClearFunc invocations
func (m *ResultMatcherMock) ClearMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ClearCounter)
}

//ClearMinimockPreCounter returns the value of ResultMatcherMock.Clear invocations
func (m *ResultMatcherMock) ClearMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClearPreCounter)
}

//ClearFinished returns true if mock invocations count is ok
func (m *ResultMatcherMock) ClearFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ClearMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ClearCounter) == uint64(len(m.ClearMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ClearMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ClearCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ClearFunc != nil {
		return atomic.LoadUint64(&m.ClearCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ResultMatcherMock) ValidateCallCounters() {

	if !m.AddStillExecutionFinished() {
		m.t.Fatal("Expected call to ResultMatcherMock.AddStillExecution")
	}

	if !m.AddUnwantedResponseFinished() {
		m.t.Fatal("Expected call to ResultMatcherMock.AddUnwantedResponse")
	}

	if !m.ClearFinished() {
		m.t.Fatal("Expected call to ResultMatcherMock.Clear")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ResultMatcherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ResultMatcherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ResultMatcherMock) MinimockFinish() {

	if !m.AddStillExecutionFinished() {
		m.t.Fatal("Expected call to ResultMatcherMock.AddStillExecution")
	}

	if !m.AddUnwantedResponseFinished() {
		m.t.Fatal("Expected call to ResultMatcherMock.AddUnwantedResponse")
	}

	if !m.ClearFinished() {
		m.t.Fatal("Expected call to ResultMatcherMock.Clear")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ResultMatcherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ResultMatcherMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddStillExecutionFinished()
		ok = ok && m.AddUnwantedResponseFinished()
		ok = ok && m.ClearFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddStillExecutionFinished() {
				m.t.Error("Expected call to ResultMatcherMock.AddStillExecution")
			}

			if !m.AddUnwantedResponseFinished() {
				m.t.Error("Expected call to ResultMatcherMock.AddUnwantedResponse")
			}

			if !m.ClearFinished() {
				m.t.Error("Expected call to ResultMatcherMock.Clear")
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
func (m *ResultMatcherMock) AllMocksCalled() bool {

	if !m.AddStillExecutionFinished() {
		return false
	}

	if !m.AddUnwantedResponseFinished() {
		return false
	}

	if !m.ClearFinished() {
		return false
	}

	return true
}

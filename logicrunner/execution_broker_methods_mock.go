package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ExecutionBrokerMethods" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//ExecutionBrokerMethodsMock implements github.com/insolar/insolar/logicrunner.ExecutionBrokerMethods
type ExecutionBrokerMethodsMock struct {
	t minimock.Tester

	CheckFunc       func(p context.Context) (r error)
	CheckCounter    uint64
	CheckPreCounter uint64
	CheckMock       mExecutionBrokerMethodsMockCheck

	ExecuteFunc       func(p context.Context, p1 *Transcript) (r error)
	ExecuteCounter    uint64
	ExecutePreCounter uint64
	ExecuteMock       mExecutionBrokerMethodsMockExecute
}

//NewExecutionBrokerMethodsMock returns a mock for github.com/insolar/insolar/logicrunner.ExecutionBrokerMethods
func NewExecutionBrokerMethodsMock(t minimock.Tester) *ExecutionBrokerMethodsMock {
	m := &ExecutionBrokerMethodsMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CheckMock = mExecutionBrokerMethodsMockCheck{mock: m}
	m.ExecuteMock = mExecutionBrokerMethodsMockExecute{mock: m}

	return m
}

type mExecutionBrokerMethodsMockCheck struct {
	mock              *ExecutionBrokerMethodsMock
	mainExpectation   *ExecutionBrokerMethodsMockCheckExpectation
	expectationSeries []*ExecutionBrokerMethodsMockCheckExpectation
}

type ExecutionBrokerMethodsMockCheckExpectation struct {
	input  *ExecutionBrokerMethodsMockCheckInput
	result *ExecutionBrokerMethodsMockCheckResult
}

type ExecutionBrokerMethodsMockCheckInput struct {
	p context.Context
}

type ExecutionBrokerMethodsMockCheckResult struct {
	r error
}

//Expect specifies that invocation of ExecutionBrokerMethods.Check is expected from 1 to Infinity times
func (m *mExecutionBrokerMethodsMockCheck) Expect(p context.Context) *mExecutionBrokerMethodsMockCheck {
	m.mock.CheckFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerMethodsMockCheckExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerMethodsMockCheckInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerMethods.Check
func (m *mExecutionBrokerMethodsMockCheck) Return(r error) *ExecutionBrokerMethodsMock {
	m.mock.CheckFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerMethodsMockCheckExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerMethodsMockCheckResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerMethods.Check is expected once
func (m *mExecutionBrokerMethodsMockCheck) ExpectOnce(p context.Context) *ExecutionBrokerMethodsMockCheckExpectation {
	m.mock.CheckFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerMethodsMockCheckExpectation{}
	expectation.input = &ExecutionBrokerMethodsMockCheckInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerMethodsMockCheckExpectation) Return(r error) {
	e.result = &ExecutionBrokerMethodsMockCheckResult{r}
}

//Set uses given function f as a mock of ExecutionBrokerMethods.Check method
func (m *mExecutionBrokerMethodsMockCheck) Set(f func(p context.Context) (r error)) *ExecutionBrokerMethodsMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CheckFunc = f
	return m.mock
}

//Check implements github.com/insolar/insolar/logicrunner.ExecutionBrokerMethods interface
func (m *ExecutionBrokerMethodsMock) Check(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.CheckPreCounter, 1)
	defer atomic.AddUint64(&m.CheckCounter, 1)

	if len(m.CheckMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CheckMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerMethodsMock.Check. %v", p)
			return
		}

		input := m.CheckMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerMethodsMockCheckInput{p}, "ExecutionBrokerMethods.Check got unexpected parameters")

		result := m.CheckMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerMethodsMock.Check")
			return
		}

		r = result.r

		return
	}

	if m.CheckMock.mainExpectation != nil {

		input := m.CheckMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerMethodsMockCheckInput{p}, "ExecutionBrokerMethods.Check got unexpected parameters")
		}

		result := m.CheckMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerMethodsMock.Check")
		}

		r = result.r

		return
	}

	if m.CheckFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerMethodsMock.Check. %v", p)
		return
	}

	return m.CheckFunc(p)
}

//CheckMinimockCounter returns a count of ExecutionBrokerMethodsMock.CheckFunc invocations
func (m *ExecutionBrokerMethodsMock) CheckMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CheckCounter)
}

//CheckMinimockPreCounter returns the value of ExecutionBrokerMethodsMock.Check invocations
func (m *ExecutionBrokerMethodsMock) CheckMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CheckPreCounter)
}

//CheckFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerMethodsMock) CheckFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CheckMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CheckCounter) == uint64(len(m.CheckMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CheckMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CheckCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CheckFunc != nil {
		return atomic.LoadUint64(&m.CheckCounter) > 0
	}

	return true
}

type mExecutionBrokerMethodsMockExecute struct {
	mock              *ExecutionBrokerMethodsMock
	mainExpectation   *ExecutionBrokerMethodsMockExecuteExpectation
	expectationSeries []*ExecutionBrokerMethodsMockExecuteExpectation
}

type ExecutionBrokerMethodsMockExecuteExpectation struct {
	input  *ExecutionBrokerMethodsMockExecuteInput
	result *ExecutionBrokerMethodsMockExecuteResult
}

type ExecutionBrokerMethodsMockExecuteInput struct {
	p  context.Context
	p1 *Transcript
}

type ExecutionBrokerMethodsMockExecuteResult struct {
	r error
}

//Expect specifies that invocation of ExecutionBrokerMethods.Execute is expected from 1 to Infinity times
func (m *mExecutionBrokerMethodsMockExecute) Expect(p context.Context, p1 *Transcript) *mExecutionBrokerMethodsMockExecute {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerMethodsMockExecuteExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerMethodsMockExecuteInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerMethods.Execute
func (m *mExecutionBrokerMethodsMockExecute) Return(r error) *ExecutionBrokerMethodsMock {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerMethodsMockExecuteExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerMethodsMockExecuteResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerMethods.Execute is expected once
func (m *mExecutionBrokerMethodsMockExecute) ExpectOnce(p context.Context, p1 *Transcript) *ExecutionBrokerMethodsMockExecuteExpectation {
	m.mock.ExecuteFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerMethodsMockExecuteExpectation{}
	expectation.input = &ExecutionBrokerMethodsMockExecuteInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerMethodsMockExecuteExpectation) Return(r error) {
	e.result = &ExecutionBrokerMethodsMockExecuteResult{r}
}

//Set uses given function f as a mock of ExecutionBrokerMethods.Execute method
func (m *mExecutionBrokerMethodsMockExecute) Set(f func(p context.Context, p1 *Transcript) (r error)) *ExecutionBrokerMethodsMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecuteFunc = f
	return m.mock
}

//Execute implements github.com/insolar/insolar/logicrunner.ExecutionBrokerMethods interface
func (m *ExecutionBrokerMethodsMock) Execute(p context.Context, p1 *Transcript) (r error) {
	counter := atomic.AddUint64(&m.ExecutePreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteCounter, 1)

	if len(m.ExecuteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecuteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerMethodsMock.Execute. %v %v", p, p1)
			return
		}

		input := m.ExecuteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerMethodsMockExecuteInput{p, p1}, "ExecutionBrokerMethods.Execute got unexpected parameters")

		result := m.ExecuteMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerMethodsMock.Execute")
			return
		}

		r = result.r

		return
	}

	if m.ExecuteMock.mainExpectation != nil {

		input := m.ExecuteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerMethodsMockExecuteInput{p, p1}, "ExecutionBrokerMethods.Execute got unexpected parameters")
		}

		result := m.ExecuteMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerMethodsMock.Execute")
		}

		r = result.r

		return
	}

	if m.ExecuteFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerMethodsMock.Execute. %v %v", p, p1)
		return
	}

	return m.ExecuteFunc(p, p1)
}

//ExecuteMinimockCounter returns a count of ExecutionBrokerMethodsMock.ExecuteFunc invocations
func (m *ExecutionBrokerMethodsMock) ExecuteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteCounter)
}

//ExecuteMinimockPreCounter returns the value of ExecutionBrokerMethodsMock.Execute invocations
func (m *ExecutionBrokerMethodsMock) ExecuteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecutePreCounter)
}

//ExecuteFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerMethodsMock) ExecuteFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExecuteMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExecuteCounter) == uint64(len(m.ExecuteMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExecuteMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExecuteCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExecuteFunc != nil {
		return atomic.LoadUint64(&m.ExecuteCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExecutionBrokerMethodsMock) ValidateCallCounters() {

	if !m.CheckFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerMethodsMock.Check")
	}

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerMethodsMock.Execute")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExecutionBrokerMethodsMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ExecutionBrokerMethodsMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ExecutionBrokerMethodsMock) MinimockFinish() {

	if !m.CheckFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerMethodsMock.Check")
	}

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerMethodsMock.Execute")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ExecutionBrokerMethodsMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ExecutionBrokerMethodsMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CheckFinished()
		ok = ok && m.ExecuteFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CheckFinished() {
				m.t.Error("Expected call to ExecutionBrokerMethodsMock.Check")
			}

			if !m.ExecuteFinished() {
				m.t.Error("Expected call to ExecutionBrokerMethodsMock.Execute")
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
func (m *ExecutionBrokerMethodsMock) AllMocksCalled() bool {

	if !m.CheckFinished() {
		return false
	}

	if !m.ExecuteFinished() {
		return false
	}

	return true
}

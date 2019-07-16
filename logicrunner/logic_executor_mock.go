package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LogicExecutor" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	artifacts "github.com/insolar/insolar/logicrunner/artifacts"

	testify_assert "github.com/stretchr/testify/assert"
)

//LogicExecutorMock implements github.com/insolar/insolar/logicrunner.LogicExecutor
type LogicExecutorMock struct {
	t minimock.Tester

	ExecuteFunc       func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)
	ExecuteCounter    uint64
	ExecutePreCounter uint64
	ExecuteMock       mLogicExecutorMockExecute

	ExecuteConstructorFunc       func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)
	ExecuteConstructorCounter    uint64
	ExecuteConstructorPreCounter uint64
	ExecuteConstructorMock       mLogicExecutorMockExecuteConstructor

	ExecuteMethodFunc       func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)
	ExecuteMethodCounter    uint64
	ExecuteMethodPreCounter uint64
	ExecuteMethodMock       mLogicExecutorMockExecuteMethod
}

//NewLogicExecutorMock returns a mock for github.com/insolar/insolar/logicrunner.LogicExecutor
func NewLogicExecutorMock(t minimock.Tester) *LogicExecutorMock {
	m := &LogicExecutorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExecuteMock = mLogicExecutorMockExecute{mock: m}
	m.ExecuteConstructorMock = mLogicExecutorMockExecuteConstructor{mock: m}
	m.ExecuteMethodMock = mLogicExecutorMockExecuteMethod{mock: m}

	return m
}

type mLogicExecutorMockExecute struct {
	mock              *LogicExecutorMock
	mainExpectation   *LogicExecutorMockExecuteExpectation
	expectationSeries []*LogicExecutorMockExecuteExpectation
}

type LogicExecutorMockExecuteExpectation struct {
	input  *LogicExecutorMockExecuteInput
	result *LogicExecutorMockExecuteResult
}

type LogicExecutorMockExecuteInput struct {
	p  context.Context
	p1 *Transcript
}

type LogicExecutorMockExecuteResult struct {
	r  artifacts.RequestResult
	r1 error
}

//Expect specifies that invocation of LogicExecutor.Execute is expected from 1 to Infinity times
func (m *mLogicExecutorMockExecute) Expect(p context.Context, p1 *Transcript) *mLogicExecutorMockExecute {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicExecutorMockExecuteExpectation{}
	}
	m.mainExpectation.input = &LogicExecutorMockExecuteInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicExecutor.Execute
func (m *mLogicExecutorMockExecute) Return(r artifacts.RequestResult, r1 error) *LogicExecutorMock {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicExecutorMockExecuteExpectation{}
	}
	m.mainExpectation.result = &LogicExecutorMockExecuteResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicExecutor.Execute is expected once
func (m *mLogicExecutorMockExecute) ExpectOnce(p context.Context, p1 *Transcript) *LogicExecutorMockExecuteExpectation {
	m.mock.ExecuteFunc = nil
	m.mainExpectation = nil

	expectation := &LogicExecutorMockExecuteExpectation{}
	expectation.input = &LogicExecutorMockExecuteInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicExecutorMockExecuteExpectation) Return(r artifacts.RequestResult, r1 error) {
	e.result = &LogicExecutorMockExecuteResult{r, r1}
}

//Set uses given function f as a mock of LogicExecutor.Execute method
func (m *mLogicExecutorMockExecute) Set(f func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)) *LogicExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecuteFunc = f
	return m.mock
}

//Execute implements github.com/insolar/insolar/logicrunner.LogicExecutor interface
func (m *LogicExecutorMock) Execute(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error) {
	counter := atomic.AddUint64(&m.ExecutePreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteCounter, 1)

	if len(m.ExecuteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecuteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicExecutorMock.Execute. %v %v", p, p1)
			return
		}

		input := m.ExecuteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicExecutorMockExecuteInput{p, p1}, "LogicExecutor.Execute got unexpected parameters")

		result := m.ExecuteMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicExecutorMock.Execute")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteMock.mainExpectation != nil {

		input := m.ExecuteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicExecutorMockExecuteInput{p, p1}, "LogicExecutor.Execute got unexpected parameters")
		}

		result := m.ExecuteMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicExecutorMock.Execute")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteFunc == nil {
		m.t.Fatalf("Unexpected call to LogicExecutorMock.Execute. %v %v", p, p1)
		return
	}

	return m.ExecuteFunc(p, p1)
}

//ExecuteMinimockCounter returns a count of LogicExecutorMock.ExecuteFunc invocations
func (m *LogicExecutorMock) ExecuteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteCounter)
}

//ExecuteMinimockPreCounter returns the value of LogicExecutorMock.Execute invocations
func (m *LogicExecutorMock) ExecuteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecutePreCounter)
}

//ExecuteFinished returns true if mock invocations count is ok
func (m *LogicExecutorMock) ExecuteFinished() bool {
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

type mLogicExecutorMockExecuteConstructor struct {
	mock              *LogicExecutorMock
	mainExpectation   *LogicExecutorMockExecuteConstructorExpectation
	expectationSeries []*LogicExecutorMockExecuteConstructorExpectation
}

type LogicExecutorMockExecuteConstructorExpectation struct {
	input  *LogicExecutorMockExecuteConstructorInput
	result *LogicExecutorMockExecuteConstructorResult
}

type LogicExecutorMockExecuteConstructorInput struct {
	p  context.Context
	p1 *Transcript
}

type LogicExecutorMockExecuteConstructorResult struct {
	r  artifacts.RequestResult
	r1 error
}

//Expect specifies that invocation of LogicExecutor.ExecuteConstructor is expected from 1 to Infinity times
func (m *mLogicExecutorMockExecuteConstructor) Expect(p context.Context, p1 *Transcript) *mLogicExecutorMockExecuteConstructor {
	m.mock.ExecuteConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicExecutorMockExecuteConstructorExpectation{}
	}
	m.mainExpectation.input = &LogicExecutorMockExecuteConstructorInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicExecutor.ExecuteConstructor
func (m *mLogicExecutorMockExecuteConstructor) Return(r artifacts.RequestResult, r1 error) *LogicExecutorMock {
	m.mock.ExecuteConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicExecutorMockExecuteConstructorExpectation{}
	}
	m.mainExpectation.result = &LogicExecutorMockExecuteConstructorResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicExecutor.ExecuteConstructor is expected once
func (m *mLogicExecutorMockExecuteConstructor) ExpectOnce(p context.Context, p1 *Transcript) *LogicExecutorMockExecuteConstructorExpectation {
	m.mock.ExecuteConstructorFunc = nil
	m.mainExpectation = nil

	expectation := &LogicExecutorMockExecuteConstructorExpectation{}
	expectation.input = &LogicExecutorMockExecuteConstructorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicExecutorMockExecuteConstructorExpectation) Return(r artifacts.RequestResult, r1 error) {
	e.result = &LogicExecutorMockExecuteConstructorResult{r, r1}
}

//Set uses given function f as a mock of LogicExecutor.ExecuteConstructor method
func (m *mLogicExecutorMockExecuteConstructor) Set(f func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)) *LogicExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecuteConstructorFunc = f
	return m.mock
}

//ExecuteConstructor implements github.com/insolar/insolar/logicrunner.LogicExecutor interface
func (m *LogicExecutorMock) ExecuteConstructor(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error) {
	counter := atomic.AddUint64(&m.ExecuteConstructorPreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteConstructorCounter, 1)

	if len(m.ExecuteConstructorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecuteConstructorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicExecutorMock.ExecuteConstructor. %v %v", p, p1)
			return
		}

		input := m.ExecuteConstructorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicExecutorMockExecuteConstructorInput{p, p1}, "LogicExecutor.ExecuteConstructor got unexpected parameters")

		result := m.ExecuteConstructorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicExecutorMock.ExecuteConstructor")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteConstructorMock.mainExpectation != nil {

		input := m.ExecuteConstructorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicExecutorMockExecuteConstructorInput{p, p1}, "LogicExecutor.ExecuteConstructor got unexpected parameters")
		}

		result := m.ExecuteConstructorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicExecutorMock.ExecuteConstructor")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteConstructorFunc == nil {
		m.t.Fatalf("Unexpected call to LogicExecutorMock.ExecuteConstructor. %v %v", p, p1)
		return
	}

	return m.ExecuteConstructorFunc(p, p1)
}

//ExecuteConstructorMinimockCounter returns a count of LogicExecutorMock.ExecuteConstructorFunc invocations
func (m *LogicExecutorMock) ExecuteConstructorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteConstructorCounter)
}

//ExecuteConstructorMinimockPreCounter returns the value of LogicExecutorMock.ExecuteConstructor invocations
func (m *LogicExecutorMock) ExecuteConstructorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteConstructorPreCounter)
}

//ExecuteConstructorFinished returns true if mock invocations count is ok
func (m *LogicExecutorMock) ExecuteConstructorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExecuteConstructorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExecuteConstructorCounter) == uint64(len(m.ExecuteConstructorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExecuteConstructorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExecuteConstructorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExecuteConstructorFunc != nil {
		return atomic.LoadUint64(&m.ExecuteConstructorCounter) > 0
	}

	return true
}

type mLogicExecutorMockExecuteMethod struct {
	mock              *LogicExecutorMock
	mainExpectation   *LogicExecutorMockExecuteMethodExpectation
	expectationSeries []*LogicExecutorMockExecuteMethodExpectation
}

type LogicExecutorMockExecuteMethodExpectation struct {
	input  *LogicExecutorMockExecuteMethodInput
	result *LogicExecutorMockExecuteMethodResult
}

type LogicExecutorMockExecuteMethodInput struct {
	p  context.Context
	p1 *Transcript
}

type LogicExecutorMockExecuteMethodResult struct {
	r  artifacts.RequestResult
	r1 error
}

//Expect specifies that invocation of LogicExecutor.ExecuteMethod is expected from 1 to Infinity times
func (m *mLogicExecutorMockExecuteMethod) Expect(p context.Context, p1 *Transcript) *mLogicExecutorMockExecuteMethod {
	m.mock.ExecuteMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicExecutorMockExecuteMethodExpectation{}
	}
	m.mainExpectation.input = &LogicExecutorMockExecuteMethodInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicExecutor.ExecuteMethod
func (m *mLogicExecutorMockExecuteMethod) Return(r artifacts.RequestResult, r1 error) *LogicExecutorMock {
	m.mock.ExecuteMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicExecutorMockExecuteMethodExpectation{}
	}
	m.mainExpectation.result = &LogicExecutorMockExecuteMethodResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicExecutor.ExecuteMethod is expected once
func (m *mLogicExecutorMockExecuteMethod) ExpectOnce(p context.Context, p1 *Transcript) *LogicExecutorMockExecuteMethodExpectation {
	m.mock.ExecuteMethodFunc = nil
	m.mainExpectation = nil

	expectation := &LogicExecutorMockExecuteMethodExpectation{}
	expectation.input = &LogicExecutorMockExecuteMethodInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicExecutorMockExecuteMethodExpectation) Return(r artifacts.RequestResult, r1 error) {
	e.result = &LogicExecutorMockExecuteMethodResult{r, r1}
}

//Set uses given function f as a mock of LogicExecutor.ExecuteMethod method
func (m *mLogicExecutorMockExecuteMethod) Set(f func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)) *LogicExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecuteMethodFunc = f
	return m.mock
}

//ExecuteMethod implements github.com/insolar/insolar/logicrunner.LogicExecutor interface
func (m *LogicExecutorMock) ExecuteMethod(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error) {
	counter := atomic.AddUint64(&m.ExecuteMethodPreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteMethodCounter, 1)

	if len(m.ExecuteMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecuteMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicExecutorMock.ExecuteMethod. %v %v", p, p1)
			return
		}

		input := m.ExecuteMethodMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicExecutorMockExecuteMethodInput{p, p1}, "LogicExecutor.ExecuteMethod got unexpected parameters")

		result := m.ExecuteMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicExecutorMock.ExecuteMethod")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteMethodMock.mainExpectation != nil {

		input := m.ExecuteMethodMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicExecutorMockExecuteMethodInput{p, p1}, "LogicExecutor.ExecuteMethod got unexpected parameters")
		}

		result := m.ExecuteMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicExecutorMock.ExecuteMethod")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteMethodFunc == nil {
		m.t.Fatalf("Unexpected call to LogicExecutorMock.ExecuteMethod. %v %v", p, p1)
		return
	}

	return m.ExecuteMethodFunc(p, p1)
}

//ExecuteMethodMinimockCounter returns a count of LogicExecutorMock.ExecuteMethodFunc invocations
func (m *LogicExecutorMock) ExecuteMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteMethodCounter)
}

//ExecuteMethodMinimockPreCounter returns the value of LogicExecutorMock.ExecuteMethod invocations
func (m *LogicExecutorMock) ExecuteMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteMethodPreCounter)
}

//ExecuteMethodFinished returns true if mock invocations count is ok
func (m *LogicExecutorMock) ExecuteMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExecuteMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExecuteMethodCounter) == uint64(len(m.ExecuteMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExecuteMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExecuteMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExecuteMethodFunc != nil {
		return atomic.LoadUint64(&m.ExecuteMethodCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicExecutorMock) ValidateCallCounters() {

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to LogicExecutorMock.Execute")
	}

	if !m.ExecuteConstructorFinished() {
		m.t.Fatal("Expected call to LogicExecutorMock.ExecuteConstructor")
	}

	if !m.ExecuteMethodFinished() {
		m.t.Fatal("Expected call to LogicExecutorMock.ExecuteMethod")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicExecutorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LogicExecutorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LogicExecutorMock) MinimockFinish() {

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to LogicExecutorMock.Execute")
	}

	if !m.ExecuteConstructorFinished() {
		m.t.Fatal("Expected call to LogicExecutorMock.ExecuteConstructor")
	}

	if !m.ExecuteMethodFinished() {
		m.t.Fatal("Expected call to LogicExecutorMock.ExecuteMethod")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LogicExecutorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LogicExecutorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ExecuteFinished()
		ok = ok && m.ExecuteConstructorFinished()
		ok = ok && m.ExecuteMethodFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ExecuteFinished() {
				m.t.Error("Expected call to LogicExecutorMock.Execute")
			}

			if !m.ExecuteConstructorFinished() {
				m.t.Error("Expected call to LogicExecutorMock.ExecuteConstructor")
			}

			if !m.ExecuteMethodFinished() {
				m.t.Error("Expected call to LogicExecutorMock.ExecuteMethod")
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
func (m *LogicExecutorMock) AllMocksCalled() bool {

	if !m.ExecuteFinished() {
		return false
	}

	if !m.ExecuteConstructorFinished() {
		return false
	}

	if !m.ExecuteMethodFinished() {
		return false
	}

	return true
}

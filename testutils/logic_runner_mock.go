package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LogicRunner" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//LogicRunnerMock implements github.com/insolar/insolar/core.LogicRunner
type LogicRunnerMock struct {
	t minimock.Tester

	ExecuteFunc       func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)
	ExecuteCounter    uint64
	ExecutePreCounter uint64
	ExecuteMock       mLogicRunnerMockExecute

	ExecutorResultsFunc       func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)
	ExecutorResultsCounter    uint64
	ExecutorResultsPreCounter uint64
	ExecutorResultsMock       mLogicRunnerMockExecutorResults

	OnPulseFunc       func(p context.Context, p1 core.Pulse) (r error)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mLogicRunnerMockOnPulse

	ProcessValidationResultsFunc       func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)
	ProcessValidationResultsCounter    uint64
	ProcessValidationResultsPreCounter uint64
	ProcessValidationResultsMock       mLogicRunnerMockProcessValidationResults

	ValidateFunc       func(p context.Context, p1 core.RecordRef, p2 core.Pulse, p3 core.CaseBind) (r int, r1 error)
	ValidateCounter    uint64
	ValidatePreCounter uint64
	ValidateMock       mLogicRunnerMockValidate

	ValidateCaseBindFunc       func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)
	ValidateCaseBindCounter    uint64
	ValidateCaseBindPreCounter uint64
	ValidateCaseBindMock       mLogicRunnerMockValidateCaseBind
}

//NewLogicRunnerMock returns a mock for github.com/insolar/insolar/core.LogicRunner
func NewLogicRunnerMock(t minimock.Tester) *LogicRunnerMock {
	m := &LogicRunnerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExecuteMock = mLogicRunnerMockExecute{mock: m}
	m.ExecutorResultsMock = mLogicRunnerMockExecutorResults{mock: m}
	m.OnPulseMock = mLogicRunnerMockOnPulse{mock: m}
	m.ProcessValidationResultsMock = mLogicRunnerMockProcessValidationResults{mock: m}
	m.ValidateMock = mLogicRunnerMockValidate{mock: m}
	m.ValidateCaseBindMock = mLogicRunnerMockValidateCaseBind{mock: m}

	return m
}

type mLogicRunnerMockExecute struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockExecuteExpectation
	expectationSeries []*LogicRunnerMockExecuteExpectation
}

type LogicRunnerMockExecuteExpectation struct {
	input  *LogicRunnerMockExecuteInput
	result *LogicRunnerMockExecuteResult
}

type LogicRunnerMockExecuteInput struct {
	p  context.Context
	p1 core.Parcel
}

type LogicRunnerMockExecuteResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of LogicRunner.Execute is expected from 1 to Infinity times
func (m *mLogicRunnerMockExecute) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockExecute {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockExecuteExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockExecuteInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.Execute
func (m *mLogicRunnerMockExecute) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockExecuteExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockExecuteResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.Execute is expected once
func (m *mLogicRunnerMockExecute) ExpectOnce(p context.Context, p1 core.Parcel) *LogicRunnerMockExecuteExpectation {
	m.mock.ExecuteFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockExecuteExpectation{}
	expectation.input = &LogicRunnerMockExecuteInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockExecuteExpectation) Return(r core.Reply, r1 error) {
	e.result = &LogicRunnerMockExecuteResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.Execute method
func (m *mLogicRunnerMockExecute) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecuteFunc = f
	return m.mock
}

//Execute implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) Execute(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.ExecutePreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteCounter, 1)

	if len(m.ExecuteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecuteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.Execute. %v %v", p, p1)
			return
		}

		input := m.ExecuteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockExecuteInput{p, p1}, "LogicRunner.Execute got unexpected parameters")

		result := m.ExecuteMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.Execute")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteMock.mainExpectation != nil {

		input := m.ExecuteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockExecuteInput{p, p1}, "LogicRunner.Execute got unexpected parameters")
		}

		result := m.ExecuteMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.Execute")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.Execute. %v %v", p, p1)
		return
	}

	return m.ExecuteFunc(p, p1)
}

//ExecuteMinimockCounter returns a count of LogicRunnerMock.ExecuteFunc invocations
func (m *LogicRunnerMock) ExecuteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteCounter)
}

//ExecuteMinimockPreCounter returns the value of LogicRunnerMock.Execute invocations
func (m *LogicRunnerMock) ExecuteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecutePreCounter)
}

//ExecuteFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) ExecuteFinished() bool {
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

type mLogicRunnerMockExecutorResults struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockExecutorResultsExpectation
	expectationSeries []*LogicRunnerMockExecutorResultsExpectation
}

type LogicRunnerMockExecutorResultsExpectation struct {
	input  *LogicRunnerMockExecutorResultsInput
	result *LogicRunnerMockExecutorResultsResult
}

type LogicRunnerMockExecutorResultsInput struct {
	p  context.Context
	p1 core.Parcel
}

type LogicRunnerMockExecutorResultsResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of LogicRunner.ExecutorResults is expected from 1 to Infinity times
func (m *mLogicRunnerMockExecutorResults) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockExecutorResults {
	m.mock.ExecutorResultsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockExecutorResultsExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockExecutorResultsInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.ExecutorResults
func (m *mLogicRunnerMockExecutorResults) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ExecutorResultsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockExecutorResultsExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockExecutorResultsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.ExecutorResults is expected once
func (m *mLogicRunnerMockExecutorResults) ExpectOnce(p context.Context, p1 core.Parcel) *LogicRunnerMockExecutorResultsExpectation {
	m.mock.ExecutorResultsFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockExecutorResultsExpectation{}
	expectation.input = &LogicRunnerMockExecutorResultsInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockExecutorResultsExpectation) Return(r core.Reply, r1 error) {
	e.result = &LogicRunnerMockExecutorResultsResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.ExecutorResults method
func (m *mLogicRunnerMockExecutorResults) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecutorResultsFunc = f
	return m.mock
}

//ExecutorResults implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) ExecutorResults(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.ExecutorResultsPreCounter, 1)
	defer atomic.AddUint64(&m.ExecutorResultsCounter, 1)

	if len(m.ExecutorResultsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecutorResultsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.ExecutorResults. %v %v", p, p1)
			return
		}

		input := m.ExecutorResultsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockExecutorResultsInput{p, p1}, "LogicRunner.ExecutorResults got unexpected parameters")

		result := m.ExecutorResultsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.ExecutorResults")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecutorResultsMock.mainExpectation != nil {

		input := m.ExecutorResultsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockExecutorResultsInput{p, p1}, "LogicRunner.ExecutorResults got unexpected parameters")
		}

		result := m.ExecutorResultsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.ExecutorResults")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecutorResultsFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.ExecutorResults. %v %v", p, p1)
		return
	}

	return m.ExecutorResultsFunc(p, p1)
}

//ExecutorResultsMinimockCounter returns a count of LogicRunnerMock.ExecutorResultsFunc invocations
func (m *LogicRunnerMock) ExecutorResultsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecutorResultsCounter)
}

//ExecutorResultsMinimockPreCounter returns the value of LogicRunnerMock.ExecutorResults invocations
func (m *LogicRunnerMock) ExecutorResultsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecutorResultsPreCounter)
}

//ExecutorResultsFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) ExecutorResultsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExecutorResultsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExecutorResultsCounter) == uint64(len(m.ExecutorResultsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExecutorResultsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExecutorResultsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExecutorResultsFunc != nil {
		return atomic.LoadUint64(&m.ExecutorResultsCounter) > 0
	}

	return true
}

type mLogicRunnerMockOnPulse struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockOnPulseExpectation
	expectationSeries []*LogicRunnerMockOnPulseExpectation
}

type LogicRunnerMockOnPulseExpectation struct {
	input  *LogicRunnerMockOnPulseInput
	result *LogicRunnerMockOnPulseResult
}

type LogicRunnerMockOnPulseInput struct {
	p  context.Context
	p1 core.Pulse
}

type LogicRunnerMockOnPulseResult struct {
	r error
}

//Expect specifies that invocation of LogicRunner.OnPulse is expected from 1 to Infinity times
func (m *mLogicRunnerMockOnPulse) Expect(p context.Context, p1 core.Pulse) *mLogicRunnerMockOnPulse {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockOnPulseExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockOnPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.OnPulse
func (m *mLogicRunnerMockOnPulse) Return(r error) *LogicRunnerMock {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockOnPulseExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockOnPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.OnPulse is expected once
func (m *mLogicRunnerMockOnPulse) ExpectOnce(p context.Context, p1 core.Pulse) *LogicRunnerMockOnPulseExpectation {
	m.mock.OnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockOnPulseExpectation{}
	expectation.input = &LogicRunnerMockOnPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockOnPulseExpectation) Return(r error) {
	e.result = &LogicRunnerMockOnPulseResult{r}
}

//Set uses given function f as a mock of LogicRunner.OnPulse method
func (m *mLogicRunnerMockOnPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) OnPulse(p context.Context, p1 core.Pulse) (r error) {
	counter := atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if len(m.OnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.OnPulse. %v %v", p, p1)
			return
		}

		input := m.OnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockOnPulseInput{p, p1}, "LogicRunner.OnPulse got unexpected parameters")

		result := m.OnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.OnPulse")
			return
		}

		r = result.r

		return
	}

	if m.OnPulseMock.mainExpectation != nil {

		input := m.OnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockOnPulseInput{p, p1}, "LogicRunner.OnPulse got unexpected parameters")
		}

		result := m.OnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.OnPulse")
		}

		r = result.r

		return
	}

	if m.OnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.OnPulse. %v %v", p, p1)
		return
	}

	return m.OnPulseFunc(p, p1)
}

//OnPulseMinimockCounter returns a count of LogicRunnerMock.OnPulseFunc invocations
func (m *LogicRunnerMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of LogicRunnerMock.OnPulse invocations
func (m *LogicRunnerMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//OnPulseFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) OnPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OnPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OnPulseCounter) == uint64(len(m.OnPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OnPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OnPulseFunc != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	return true
}

type mLogicRunnerMockProcessValidationResults struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockProcessValidationResultsExpectation
	expectationSeries []*LogicRunnerMockProcessValidationResultsExpectation
}

type LogicRunnerMockProcessValidationResultsExpectation struct {
	input  *LogicRunnerMockProcessValidationResultsInput
	result *LogicRunnerMockProcessValidationResultsResult
}

type LogicRunnerMockProcessValidationResultsInput struct {
	p  context.Context
	p1 core.Parcel
}

type LogicRunnerMockProcessValidationResultsResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of LogicRunner.ProcessValidationResults is expected from 1 to Infinity times
func (m *mLogicRunnerMockProcessValidationResults) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockProcessValidationResults {
	m.mock.ProcessValidationResultsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockProcessValidationResultsExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockProcessValidationResultsInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.ProcessValidationResults
func (m *mLogicRunnerMockProcessValidationResults) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ProcessValidationResultsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockProcessValidationResultsExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockProcessValidationResultsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.ProcessValidationResults is expected once
func (m *mLogicRunnerMockProcessValidationResults) ExpectOnce(p context.Context, p1 core.Parcel) *LogicRunnerMockProcessValidationResultsExpectation {
	m.mock.ProcessValidationResultsFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockProcessValidationResultsExpectation{}
	expectation.input = &LogicRunnerMockProcessValidationResultsInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockProcessValidationResultsExpectation) Return(r core.Reply, r1 error) {
	e.result = &LogicRunnerMockProcessValidationResultsResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.ProcessValidationResults method
func (m *mLogicRunnerMockProcessValidationResults) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ProcessValidationResultsFunc = f
	return m.mock
}

//ProcessValidationResults implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) ProcessValidationResults(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.ProcessValidationResultsPreCounter, 1)
	defer atomic.AddUint64(&m.ProcessValidationResultsCounter, 1)

	if len(m.ProcessValidationResultsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ProcessValidationResultsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.ProcessValidationResults. %v %v", p, p1)
			return
		}

		input := m.ProcessValidationResultsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockProcessValidationResultsInput{p, p1}, "LogicRunner.ProcessValidationResults got unexpected parameters")

		result := m.ProcessValidationResultsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.ProcessValidationResults")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ProcessValidationResultsMock.mainExpectation != nil {

		input := m.ProcessValidationResultsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockProcessValidationResultsInput{p, p1}, "LogicRunner.ProcessValidationResults got unexpected parameters")
		}

		result := m.ProcessValidationResultsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.ProcessValidationResults")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ProcessValidationResultsFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.ProcessValidationResults. %v %v", p, p1)
		return
	}

	return m.ProcessValidationResultsFunc(p, p1)
}

//ProcessValidationResultsMinimockCounter returns a count of LogicRunnerMock.ProcessValidationResultsFunc invocations
func (m *LogicRunnerMock) ProcessValidationResultsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ProcessValidationResultsCounter)
}

//ProcessValidationResultsMinimockPreCounter returns the value of LogicRunnerMock.ProcessValidationResults invocations
func (m *LogicRunnerMock) ProcessValidationResultsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ProcessValidationResultsPreCounter)
}

//ProcessValidationResultsFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) ProcessValidationResultsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ProcessValidationResultsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ProcessValidationResultsCounter) == uint64(len(m.ProcessValidationResultsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ProcessValidationResultsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ProcessValidationResultsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ProcessValidationResultsFunc != nil {
		return atomic.LoadUint64(&m.ProcessValidationResultsCounter) > 0
	}

	return true
}

type mLogicRunnerMockValidate struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockValidateExpectation
	expectationSeries []*LogicRunnerMockValidateExpectation
}

type LogicRunnerMockValidateExpectation struct {
	input  *LogicRunnerMockValidateInput
	result *LogicRunnerMockValidateResult
}

type LogicRunnerMockValidateInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.Pulse
	p3 core.CaseBind
}

type LogicRunnerMockValidateResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of LogicRunner.Validate is expected from 1 to Infinity times
func (m *mLogicRunnerMockValidate) Expect(p context.Context, p1 core.RecordRef, p2 core.Pulse, p3 core.CaseBind) *mLogicRunnerMockValidate {
	m.mock.ValidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockValidateExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockValidateInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of LogicRunner.Validate
func (m *mLogicRunnerMockValidate) Return(r int, r1 error) *LogicRunnerMock {
	m.mock.ValidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockValidateExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockValidateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.Validate is expected once
func (m *mLogicRunnerMockValidate) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.Pulse, p3 core.CaseBind) *LogicRunnerMockValidateExpectation {
	m.mock.ValidateFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockValidateExpectation{}
	expectation.input = &LogicRunnerMockValidateInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockValidateExpectation) Return(r int, r1 error) {
	e.result = &LogicRunnerMockValidateResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.Validate method
func (m *mLogicRunnerMockValidate) Set(f func(p context.Context, p1 core.RecordRef, p2 core.Pulse, p3 core.CaseBind) (r int, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ValidateFunc = f
	return m.mock
}

//Validate implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) Validate(p context.Context, p1 core.RecordRef, p2 core.Pulse, p3 core.CaseBind) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ValidatePreCounter, 1)
	defer atomic.AddUint64(&m.ValidateCounter, 1)

	if len(m.ValidateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ValidateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.Validate. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.ValidateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockValidateInput{p, p1, p2, p3}, "LogicRunner.Validate got unexpected parameters")

		result := m.ValidateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.Validate")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateMock.mainExpectation != nil {

		input := m.ValidateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockValidateInput{p, p1, p2, p3}, "LogicRunner.Validate got unexpected parameters")
		}

		result := m.ValidateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.Validate")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.Validate. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.ValidateFunc(p, p1, p2, p3)
}

//ValidateMinimockCounter returns a count of LogicRunnerMock.ValidateFunc invocations
func (m *LogicRunnerMock) ValidateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ValidateCounter)
}

//ValidateMinimockPreCounter returns the value of LogicRunnerMock.Validate invocations
func (m *LogicRunnerMock) ValidateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ValidatePreCounter)
}

//ValidateFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) ValidateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ValidateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ValidateCounter) == uint64(len(m.ValidateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ValidateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ValidateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ValidateFunc != nil {
		return atomic.LoadUint64(&m.ValidateCounter) > 0
	}

	return true
}

type mLogicRunnerMockValidateCaseBind struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockValidateCaseBindExpectation
	expectationSeries []*LogicRunnerMockValidateCaseBindExpectation
}

type LogicRunnerMockValidateCaseBindExpectation struct {
	input  *LogicRunnerMockValidateCaseBindInput
	result *LogicRunnerMockValidateCaseBindResult
}

type LogicRunnerMockValidateCaseBindInput struct {
	p  context.Context
	p1 core.Parcel
}

type LogicRunnerMockValidateCaseBindResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of LogicRunner.ValidateCaseBind is expected from 1 to Infinity times
func (m *mLogicRunnerMockValidateCaseBind) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockValidateCaseBind {
	m.mock.ValidateCaseBindFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockValidateCaseBindExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockValidateCaseBindInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.ValidateCaseBind
func (m *mLogicRunnerMockValidateCaseBind) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ValidateCaseBindFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockValidateCaseBindExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockValidateCaseBindResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.ValidateCaseBind is expected once
func (m *mLogicRunnerMockValidateCaseBind) ExpectOnce(p context.Context, p1 core.Parcel) *LogicRunnerMockValidateCaseBindExpectation {
	m.mock.ValidateCaseBindFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockValidateCaseBindExpectation{}
	expectation.input = &LogicRunnerMockValidateCaseBindInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockValidateCaseBindExpectation) Return(r core.Reply, r1 error) {
	e.result = &LogicRunnerMockValidateCaseBindResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.ValidateCaseBind method
func (m *mLogicRunnerMockValidateCaseBind) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ValidateCaseBindFunc = f
	return m.mock
}

//ValidateCaseBind implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) ValidateCaseBind(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.ValidateCaseBindPreCounter, 1)
	defer atomic.AddUint64(&m.ValidateCaseBindCounter, 1)

	if len(m.ValidateCaseBindMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ValidateCaseBindMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.ValidateCaseBind. %v %v", p, p1)
			return
		}

		input := m.ValidateCaseBindMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockValidateCaseBindInput{p, p1}, "LogicRunner.ValidateCaseBind got unexpected parameters")

		result := m.ValidateCaseBindMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.ValidateCaseBind")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateCaseBindMock.mainExpectation != nil {

		input := m.ValidateCaseBindMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockValidateCaseBindInput{p, p1}, "LogicRunner.ValidateCaseBind got unexpected parameters")
		}

		result := m.ValidateCaseBindMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.ValidateCaseBind")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateCaseBindFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.ValidateCaseBind. %v %v", p, p1)
		return
	}

	return m.ValidateCaseBindFunc(p, p1)
}

//ValidateCaseBindMinimockCounter returns a count of LogicRunnerMock.ValidateCaseBindFunc invocations
func (m *LogicRunnerMock) ValidateCaseBindMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ValidateCaseBindCounter)
}

//ValidateCaseBindMinimockPreCounter returns the value of LogicRunnerMock.ValidateCaseBind invocations
func (m *LogicRunnerMock) ValidateCaseBindMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ValidateCaseBindPreCounter)
}

//ValidateCaseBindFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) ValidateCaseBindFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ValidateCaseBindMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ValidateCaseBindCounter) == uint64(len(m.ValidateCaseBindMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ValidateCaseBindMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ValidateCaseBindCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ValidateCaseBindFunc != nil {
		return atomic.LoadUint64(&m.ValidateCaseBindCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicRunnerMock) ValidateCallCounters() {

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.Execute")
	}

	if !m.ExecutorResultsFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.ExecutorResults")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
	}

	if !m.ProcessValidationResultsFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.ProcessValidationResults")
	}

	if !m.ValidateFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.Validate")
	}

	if !m.ValidateCaseBindFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.ValidateCaseBind")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicRunnerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LogicRunnerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LogicRunnerMock) MinimockFinish() {

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.Execute")
	}

	if !m.ExecutorResultsFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.ExecutorResults")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
	}

	if !m.ProcessValidationResultsFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.ProcessValidationResults")
	}

	if !m.ValidateFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.Validate")
	}

	if !m.ValidateCaseBindFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.ValidateCaseBind")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LogicRunnerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LogicRunnerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ExecuteFinished()
		ok = ok && m.ExecutorResultsFinished()
		ok = ok && m.OnPulseFinished()
		ok = ok && m.ProcessValidationResultsFinished()
		ok = ok && m.ValidateFinished()
		ok = ok && m.ValidateCaseBindFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ExecuteFinished() {
				m.t.Error("Expected call to LogicRunnerMock.Execute")
			}

			if !m.ExecutorResultsFinished() {
				m.t.Error("Expected call to LogicRunnerMock.ExecutorResults")
			}

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to LogicRunnerMock.OnPulse")
			}

			if !m.ProcessValidationResultsFinished() {
				m.t.Error("Expected call to LogicRunnerMock.ProcessValidationResults")
			}

			if !m.ValidateFinished() {
				m.t.Error("Expected call to LogicRunnerMock.Validate")
			}

			if !m.ValidateCaseBindFinished() {
				m.t.Error("Expected call to LogicRunnerMock.ValidateCaseBind")
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
func (m *LogicRunnerMock) AllMocksCalled() bool {

	if !m.ExecuteFinished() {
		return false
	}

	if !m.ExecutorResultsFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	if !m.ProcessValidationResultsFinished() {
		return false
	}

	if !m.ValidateFinished() {
		return false
	}

	if !m.ValidateCaseBindFinished() {
		return false
	}

	return true
}

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

	HandleExecutorResultsMessageFunc       func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)
	HandleExecutorResultsMessageCounter    uint64
	HandleExecutorResultsMessagePreCounter uint64
	HandleExecutorResultsMessageMock       mLogicRunnerMockHandleExecutorResultsMessage

	HandleValidateCaseBindMessageFunc       func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)
	HandleValidateCaseBindMessageCounter    uint64
	HandleValidateCaseBindMessagePreCounter uint64
	HandleValidateCaseBindMessageMock       mLogicRunnerMockHandleValidateCaseBindMessage

	HandleValidationResultsMessageFunc       func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)
	HandleValidationResultsMessageCounter    uint64
	HandleValidationResultsMessagePreCounter uint64
	HandleValidationResultsMessageMock       mLogicRunnerMockHandleValidationResultsMessage

	OnPulseFunc       func(p context.Context, p1 core.Pulse) (r error)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mLogicRunnerMockOnPulse
}

//NewLogicRunnerMock returns a mock for github.com/insolar/insolar/core.LogicRunner
func NewLogicRunnerMock(t minimock.Tester) *LogicRunnerMock {
	m := &LogicRunnerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExecuteMock = mLogicRunnerMockExecute{mock: m}
	m.HandleExecutorResultsMessageMock = mLogicRunnerMockHandleExecutorResultsMessage{mock: m}
	m.HandleValidateCaseBindMessageMock = mLogicRunnerMockHandleValidateCaseBindMessage{mock: m}
	m.HandleValidationResultsMessageMock = mLogicRunnerMockHandleValidationResultsMessage{mock: m}
	m.OnPulseMock = mLogicRunnerMockOnPulse{mock: m}

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

type mLogicRunnerMockHandleExecutorResultsMessage struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockHandleExecutorResultsMessageExpectation
	expectationSeries []*LogicRunnerMockHandleExecutorResultsMessageExpectation
}

type LogicRunnerMockHandleExecutorResultsMessageExpectation struct {
	input  *LogicRunnerMockHandleExecutorResultsMessageInput
	result *LogicRunnerMockHandleExecutorResultsMessageResult
}

type LogicRunnerMockHandleExecutorResultsMessageInput struct {
	p  context.Context
	p1 core.Parcel
}

type LogicRunnerMockHandleExecutorResultsMessageResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of LogicRunner.HandleExecutorResultsMessage is expected from 1 to Infinity times
func (m *mLogicRunnerMockHandleExecutorResultsMessage) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockHandleExecutorResultsMessage {
	m.mock.HandleExecutorResultsMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockHandleExecutorResultsMessageExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockHandleExecutorResultsMessageInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.HandleExecutorResultsMessage
func (m *mLogicRunnerMockHandleExecutorResultsMessage) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.HandleExecutorResultsMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockHandleExecutorResultsMessageExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockHandleExecutorResultsMessageResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.HandleExecutorResultsMessage is expected once
func (m *mLogicRunnerMockHandleExecutorResultsMessage) ExpectOnce(p context.Context, p1 core.Parcel) *LogicRunnerMockHandleExecutorResultsMessageExpectation {
	m.mock.HandleExecutorResultsMessageFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockHandleExecutorResultsMessageExpectation{}
	expectation.input = &LogicRunnerMockHandleExecutorResultsMessageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockHandleExecutorResultsMessageExpectation) Return(r core.Reply, r1 error) {
	e.result = &LogicRunnerMockHandleExecutorResultsMessageResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.HandleExecutorResultsMessage method
func (m *mLogicRunnerMockHandleExecutorResultsMessage) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HandleExecutorResultsMessageFunc = f
	return m.mock
}

//HandleExecutorResultsMessage implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) HandleExecutorResultsMessage(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.HandleExecutorResultsMessagePreCounter, 1)
	defer atomic.AddUint64(&m.HandleExecutorResultsMessageCounter, 1)

	if len(m.HandleExecutorResultsMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HandleExecutorResultsMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.HandleExecutorResultsMessage. %v %v", p, p1)
			return
		}

		input := m.HandleExecutorResultsMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockHandleExecutorResultsMessageInput{p, p1}, "LogicRunner.HandleExecutorResultsMessage got unexpected parameters")

		result := m.HandleExecutorResultsMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.HandleExecutorResultsMessage")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HandleExecutorResultsMessageMock.mainExpectation != nil {

		input := m.HandleExecutorResultsMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockHandleExecutorResultsMessageInput{p, p1}, "LogicRunner.HandleExecutorResultsMessage got unexpected parameters")
		}

		result := m.HandleExecutorResultsMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.HandleExecutorResultsMessage")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HandleExecutorResultsMessageFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.HandleExecutorResultsMessage. %v %v", p, p1)
		return
	}

	return m.HandleExecutorResultsMessageFunc(p, p1)
}

//HandleExecutorResultsMessageMinimockCounter returns a count of LogicRunnerMock.HandleExecutorResultsMessageFunc invocations
func (m *LogicRunnerMock) HandleExecutorResultsMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HandleExecutorResultsMessageCounter)
}

//HandleExecutorResultsMessageMinimockPreCounter returns the value of LogicRunnerMock.HandleExecutorResultsMessage invocations
func (m *LogicRunnerMock) HandleExecutorResultsMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HandleExecutorResultsMessagePreCounter)
}

//HandleExecutorResultsMessageFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) HandleExecutorResultsMessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HandleExecutorResultsMessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HandleExecutorResultsMessageCounter) == uint64(len(m.HandleExecutorResultsMessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HandleExecutorResultsMessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HandleExecutorResultsMessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HandleExecutorResultsMessageFunc != nil {
		return atomic.LoadUint64(&m.HandleExecutorResultsMessageCounter) > 0
	}

	return true
}

type mLogicRunnerMockHandleValidateCaseBindMessage struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockHandleValidateCaseBindMessageExpectation
	expectationSeries []*LogicRunnerMockHandleValidateCaseBindMessageExpectation
}

type LogicRunnerMockHandleValidateCaseBindMessageExpectation struct {
	input  *LogicRunnerMockHandleValidateCaseBindMessageInput
	result *LogicRunnerMockHandleValidateCaseBindMessageResult
}

type LogicRunnerMockHandleValidateCaseBindMessageInput struct {
	p  context.Context
	p1 core.Parcel
}

type LogicRunnerMockHandleValidateCaseBindMessageResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of LogicRunner.HandleValidateCaseBindMessage is expected from 1 to Infinity times
func (m *mLogicRunnerMockHandleValidateCaseBindMessage) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockHandleValidateCaseBindMessage {
	m.mock.HandleValidateCaseBindMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockHandleValidateCaseBindMessageExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockHandleValidateCaseBindMessageInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.HandleValidateCaseBindMessage
func (m *mLogicRunnerMockHandleValidateCaseBindMessage) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.HandleValidateCaseBindMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockHandleValidateCaseBindMessageExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockHandleValidateCaseBindMessageResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.HandleValidateCaseBindMessage is expected once
func (m *mLogicRunnerMockHandleValidateCaseBindMessage) ExpectOnce(p context.Context, p1 core.Parcel) *LogicRunnerMockHandleValidateCaseBindMessageExpectation {
	m.mock.HandleValidateCaseBindMessageFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockHandleValidateCaseBindMessageExpectation{}
	expectation.input = &LogicRunnerMockHandleValidateCaseBindMessageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockHandleValidateCaseBindMessageExpectation) Return(r core.Reply, r1 error) {
	e.result = &LogicRunnerMockHandleValidateCaseBindMessageResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.HandleValidateCaseBindMessage method
func (m *mLogicRunnerMockHandleValidateCaseBindMessage) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HandleValidateCaseBindMessageFunc = f
	return m.mock
}

//HandleValidateCaseBindMessage implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) HandleValidateCaseBindMessage(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.HandleValidateCaseBindMessagePreCounter, 1)
	defer atomic.AddUint64(&m.HandleValidateCaseBindMessageCounter, 1)

	if len(m.HandleValidateCaseBindMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HandleValidateCaseBindMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.HandleValidateCaseBindMessage. %v %v", p, p1)
			return
		}

		input := m.HandleValidateCaseBindMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockHandleValidateCaseBindMessageInput{p, p1}, "LogicRunner.HandleValidateCaseBindMessage got unexpected parameters")

		result := m.HandleValidateCaseBindMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.HandleValidateCaseBindMessage")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HandleValidateCaseBindMessageMock.mainExpectation != nil {

		input := m.HandleValidateCaseBindMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockHandleValidateCaseBindMessageInput{p, p1}, "LogicRunner.HandleValidateCaseBindMessage got unexpected parameters")
		}

		result := m.HandleValidateCaseBindMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.HandleValidateCaseBindMessage")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HandleValidateCaseBindMessageFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.HandleValidateCaseBindMessage. %v %v", p, p1)
		return
	}

	return m.HandleValidateCaseBindMessageFunc(p, p1)
}

//HandleValidateCaseBindMessageMinimockCounter returns a count of LogicRunnerMock.HandleValidateCaseBindMessageFunc invocations
func (m *LogicRunnerMock) HandleValidateCaseBindMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HandleValidateCaseBindMessageCounter)
}

//HandleValidateCaseBindMessageMinimockPreCounter returns the value of LogicRunnerMock.HandleValidateCaseBindMessage invocations
func (m *LogicRunnerMock) HandleValidateCaseBindMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HandleValidateCaseBindMessagePreCounter)
}

//HandleValidateCaseBindMessageFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) HandleValidateCaseBindMessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HandleValidateCaseBindMessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HandleValidateCaseBindMessageCounter) == uint64(len(m.HandleValidateCaseBindMessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HandleValidateCaseBindMessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HandleValidateCaseBindMessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HandleValidateCaseBindMessageFunc != nil {
		return atomic.LoadUint64(&m.HandleValidateCaseBindMessageCounter) > 0
	}

	return true
}

type mLogicRunnerMockHandleValidationResultsMessage struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockHandleValidationResultsMessageExpectation
	expectationSeries []*LogicRunnerMockHandleValidationResultsMessageExpectation
}

type LogicRunnerMockHandleValidationResultsMessageExpectation struct {
	input  *LogicRunnerMockHandleValidationResultsMessageInput
	result *LogicRunnerMockHandleValidationResultsMessageResult
}

type LogicRunnerMockHandleValidationResultsMessageInput struct {
	p  context.Context
	p1 core.Parcel
}

type LogicRunnerMockHandleValidationResultsMessageResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of LogicRunner.HandleValidationResultsMessage is expected from 1 to Infinity times
func (m *mLogicRunnerMockHandleValidationResultsMessage) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockHandleValidationResultsMessage {
	m.mock.HandleValidationResultsMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockHandleValidationResultsMessageExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockHandleValidationResultsMessageInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.HandleValidationResultsMessage
func (m *mLogicRunnerMockHandleValidationResultsMessage) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.HandleValidationResultsMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockHandleValidationResultsMessageExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockHandleValidationResultsMessageResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.HandleValidationResultsMessage is expected once
func (m *mLogicRunnerMockHandleValidationResultsMessage) ExpectOnce(p context.Context, p1 core.Parcel) *LogicRunnerMockHandleValidationResultsMessageExpectation {
	m.mock.HandleValidationResultsMessageFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockHandleValidationResultsMessageExpectation{}
	expectation.input = &LogicRunnerMockHandleValidationResultsMessageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockHandleValidationResultsMessageExpectation) Return(r core.Reply, r1 error) {
	e.result = &LogicRunnerMockHandleValidationResultsMessageResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.HandleValidationResultsMessage method
func (m *mLogicRunnerMockHandleValidationResultsMessage) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HandleValidationResultsMessageFunc = f
	return m.mock
}

//HandleValidationResultsMessage implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) HandleValidationResultsMessage(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.HandleValidationResultsMessagePreCounter, 1)
	defer atomic.AddUint64(&m.HandleValidationResultsMessageCounter, 1)

	if len(m.HandleValidationResultsMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HandleValidationResultsMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.HandleValidationResultsMessage. %v %v", p, p1)
			return
		}

		input := m.HandleValidationResultsMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockHandleValidationResultsMessageInput{p, p1}, "LogicRunner.HandleValidationResultsMessage got unexpected parameters")

		result := m.HandleValidationResultsMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.HandleValidationResultsMessage")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HandleValidationResultsMessageMock.mainExpectation != nil {

		input := m.HandleValidationResultsMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockHandleValidationResultsMessageInput{p, p1}, "LogicRunner.HandleValidationResultsMessage got unexpected parameters")
		}

		result := m.HandleValidationResultsMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.HandleValidationResultsMessage")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HandleValidationResultsMessageFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.HandleValidationResultsMessage. %v %v", p, p1)
		return
	}

	return m.HandleValidationResultsMessageFunc(p, p1)
}

//HandleValidationResultsMessageMinimockCounter returns a count of LogicRunnerMock.HandleValidationResultsMessageFunc invocations
func (m *LogicRunnerMock) HandleValidationResultsMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HandleValidationResultsMessageCounter)
}

//HandleValidationResultsMessageMinimockPreCounter returns the value of LogicRunnerMock.HandleValidationResultsMessage invocations
func (m *LogicRunnerMock) HandleValidationResultsMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HandleValidationResultsMessagePreCounter)
}

//HandleValidationResultsMessageFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) HandleValidationResultsMessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HandleValidationResultsMessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HandleValidationResultsMessageCounter) == uint64(len(m.HandleValidationResultsMessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HandleValidationResultsMessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HandleValidationResultsMessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HandleValidationResultsMessageFunc != nil {
		return atomic.LoadUint64(&m.HandleValidationResultsMessageCounter) > 0
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicRunnerMock) ValidateCallCounters() {

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.Execute")
	}

	if !m.HandleExecutorResultsMessageFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.HandleExecutorResultsMessage")
	}

	if !m.HandleValidateCaseBindMessageFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.HandleValidateCaseBindMessage")
	}

	if !m.HandleValidationResultsMessageFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.HandleValidationResultsMessage")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
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

	if !m.HandleExecutorResultsMessageFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.HandleExecutorResultsMessage")
	}

	if !m.HandleValidateCaseBindMessageFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.HandleValidateCaseBindMessage")
	}

	if !m.HandleValidationResultsMessageFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.HandleValidationResultsMessage")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
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
		ok = ok && m.HandleExecutorResultsMessageFinished()
		ok = ok && m.HandleValidateCaseBindMessageFinished()
		ok = ok && m.HandleValidationResultsMessageFinished()
		ok = ok && m.OnPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ExecuteFinished() {
				m.t.Error("Expected call to LogicRunnerMock.Execute")
			}

			if !m.HandleExecutorResultsMessageFinished() {
				m.t.Error("Expected call to LogicRunnerMock.HandleExecutorResultsMessage")
			}

			if !m.HandleValidateCaseBindMessageFinished() {
				m.t.Error("Expected call to LogicRunnerMock.HandleValidateCaseBindMessage")
			}

			if !m.HandleValidationResultsMessageFinished() {
				m.t.Error("Expected call to LogicRunnerMock.HandleValidationResultsMessage")
			}

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to LogicRunnerMock.OnPulse")
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

	if !m.HandleExecutorResultsMessageFinished() {
		return false
	}

	if !m.HandleValidateCaseBindMessageFinished() {
		return false
	}

	if !m.HandleValidationResultsMessageFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	return true
}

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
	m.ValidateCaseBindMock = mLogicRunnerMockValidateCaseBind{mock: m}

	return m
}

type mLogicRunnerMockExecute struct {
	mock             *LogicRunnerMock
	mockExpectations *LogicRunnerMockExecuteParams
}

//LogicRunnerMockExecuteParams represents input parameters of the LogicRunner.Execute
type LogicRunnerMockExecuteParams struct {
	p  context.Context
	p1 core.Parcel
}

//Expect sets up expected params for the LogicRunner.Execute
func (m *mLogicRunnerMockExecute) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockExecute {
	m.mockExpectations = &LogicRunnerMockExecuteParams{p, p1}
	return m
}

//Return sets up a mock for LogicRunner.Execute to return Return's arguments
func (m *mLogicRunnerMockExecute) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ExecuteFunc = func(p context.Context, p1 core.Parcel) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of LogicRunner.Execute method
func (m *mLogicRunnerMockExecute) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mock.ExecuteFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Execute implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) Execute(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.ExecutePreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteCounter, 1)

	if m.ExecuteMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExecuteMock.mockExpectations, LogicRunnerMockExecuteParams{p, p1},
			"LogicRunner.Execute got unexpected parameters")

		if m.ExecuteFunc == nil {

			m.t.Fatal("No results are set for the LogicRunnerMock.Execute")

			return
		}
	}

	if m.ExecuteFunc == nil {
		m.t.Fatal("Unexpected call to LogicRunnerMock.Execute")
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

type mLogicRunnerMockExecutorResults struct {
	mock             *LogicRunnerMock
	mockExpectations *LogicRunnerMockExecutorResultsParams
}

//LogicRunnerMockExecutorResultsParams represents input parameters of the LogicRunner.ExecutorResults
type LogicRunnerMockExecutorResultsParams struct {
	p  context.Context
	p1 core.Parcel
}

//Expect sets up expected params for the LogicRunner.ExecutorResults
func (m *mLogicRunnerMockExecutorResults) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockExecutorResults {
	m.mockExpectations = &LogicRunnerMockExecutorResultsParams{p, p1}
	return m
}

//Return sets up a mock for LogicRunner.ExecutorResults to return Return's arguments
func (m *mLogicRunnerMockExecutorResults) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ExecutorResultsFunc = func(p context.Context, p1 core.Parcel) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of LogicRunner.ExecutorResults method
func (m *mLogicRunnerMockExecutorResults) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mock.ExecutorResultsFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ExecutorResults implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) ExecutorResults(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.ExecutorResultsPreCounter, 1)
	defer atomic.AddUint64(&m.ExecutorResultsCounter, 1)

	if m.ExecutorResultsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExecutorResultsMock.mockExpectations, LogicRunnerMockExecutorResultsParams{p, p1},
			"LogicRunner.ExecutorResults got unexpected parameters")

		if m.ExecutorResultsFunc == nil {

			m.t.Fatal("No results are set for the LogicRunnerMock.ExecutorResults")

			return
		}
	}

	if m.ExecutorResultsFunc == nil {
		m.t.Fatal("Unexpected call to LogicRunnerMock.ExecutorResults")
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

type mLogicRunnerMockOnPulse struct {
	mock             *LogicRunnerMock
	mockExpectations *LogicRunnerMockOnPulseParams
}

//LogicRunnerMockOnPulseParams represents input parameters of the LogicRunner.OnPulse
type LogicRunnerMockOnPulseParams struct {
	p  context.Context
	p1 core.Pulse
}

//Expect sets up expected params for the LogicRunner.OnPulse
func (m *mLogicRunnerMockOnPulse) Expect(p context.Context, p1 core.Pulse) *mLogicRunnerMockOnPulse {
	m.mockExpectations = &LogicRunnerMockOnPulseParams{p, p1}
	return m
}

//Return sets up a mock for LogicRunner.OnPulse to return Return's arguments
func (m *mLogicRunnerMockOnPulse) Return(r error) *LogicRunnerMock {
	m.mock.OnPulseFunc = func(p context.Context, p1 core.Pulse) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of LogicRunner.OnPulse method
func (m *mLogicRunnerMockOnPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *LogicRunnerMock {
	m.mock.OnPulseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) OnPulse(p context.Context, p1 core.Pulse) (r error) {
	atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if m.OnPulseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.OnPulseMock.mockExpectations, LogicRunnerMockOnPulseParams{p, p1},
			"LogicRunner.OnPulse got unexpected parameters")

		if m.OnPulseFunc == nil {

			m.t.Fatal("No results are set for the LogicRunnerMock.OnPulse")

			return
		}
	}

	if m.OnPulseFunc == nil {
		m.t.Fatal("Unexpected call to LogicRunnerMock.OnPulse")
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

type mLogicRunnerMockProcessValidationResults struct {
	mock             *LogicRunnerMock
	mockExpectations *LogicRunnerMockProcessValidationResultsParams
}

//LogicRunnerMockProcessValidationResultsParams represents input parameters of the LogicRunner.ProcessValidationResults
type LogicRunnerMockProcessValidationResultsParams struct {
	p  context.Context
	p1 core.Parcel
}

//Expect sets up expected params for the LogicRunner.ProcessValidationResults
func (m *mLogicRunnerMockProcessValidationResults) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockProcessValidationResults {
	m.mockExpectations = &LogicRunnerMockProcessValidationResultsParams{p, p1}
	return m
}

//Return sets up a mock for LogicRunner.ProcessValidationResults to return Return's arguments
func (m *mLogicRunnerMockProcessValidationResults) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ProcessValidationResultsFunc = func(p context.Context, p1 core.Parcel) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of LogicRunner.ProcessValidationResults method
func (m *mLogicRunnerMockProcessValidationResults) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mock.ProcessValidationResultsFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ProcessValidationResults implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) ProcessValidationResults(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.ProcessValidationResultsPreCounter, 1)
	defer atomic.AddUint64(&m.ProcessValidationResultsCounter, 1)

	if m.ProcessValidationResultsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ProcessValidationResultsMock.mockExpectations, LogicRunnerMockProcessValidationResultsParams{p, p1},
			"LogicRunner.ProcessValidationResults got unexpected parameters")

		if m.ProcessValidationResultsFunc == nil {

			m.t.Fatal("No results are set for the LogicRunnerMock.ProcessValidationResults")

			return
		}
	}

	if m.ProcessValidationResultsFunc == nil {
		m.t.Fatal("Unexpected call to LogicRunnerMock.ProcessValidationResults")
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

type mLogicRunnerMockValidateCaseBind struct {
	mock             *LogicRunnerMock
	mockExpectations *LogicRunnerMockValidateCaseBindParams
}

//LogicRunnerMockValidateCaseBindParams represents input parameters of the LogicRunner.ValidateCaseBind
type LogicRunnerMockValidateCaseBindParams struct {
	p  context.Context
	p1 core.Parcel
}

//Expect sets up expected params for the LogicRunner.ValidateCaseBind
func (m *mLogicRunnerMockValidateCaseBind) Expect(p context.Context, p1 core.Parcel) *mLogicRunnerMockValidateCaseBind {
	m.mockExpectations = &LogicRunnerMockValidateCaseBindParams{p, p1}
	return m
}

//Return sets up a mock for LogicRunner.ValidateCaseBind to return Return's arguments
func (m *mLogicRunnerMockValidateCaseBind) Return(r core.Reply, r1 error) *LogicRunnerMock {
	m.mock.ValidateCaseBindFunc = func(p context.Context, p1 core.Parcel) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of LogicRunner.ValidateCaseBind method
func (m *mLogicRunnerMockValidateCaseBind) Set(f func(p context.Context, p1 core.Parcel) (r core.Reply, r1 error)) *LogicRunnerMock {
	m.mock.ValidateCaseBindFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ValidateCaseBind implements github.com/insolar/insolar/core.LogicRunner interface
func (m *LogicRunnerMock) ValidateCaseBind(p context.Context, p1 core.Parcel) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.ValidateCaseBindPreCounter, 1)
	defer atomic.AddUint64(&m.ValidateCaseBindCounter, 1)

	if m.ValidateCaseBindMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ValidateCaseBindMock.mockExpectations, LogicRunnerMockValidateCaseBindParams{p, p1},
			"LogicRunner.ValidateCaseBind got unexpected parameters")

		if m.ValidateCaseBindFunc == nil {

			m.t.Fatal("No results are set for the LogicRunnerMock.ValidateCaseBind")

			return
		}
	}

	if m.ValidateCaseBindFunc == nil {
		m.t.Fatal("Unexpected call to LogicRunnerMock.ValidateCaseBind")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicRunnerMock) ValidateCallCounters() {

	if m.ExecuteFunc != nil && atomic.LoadUint64(&m.ExecuteCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.Execute")
	}

	if m.ExecutorResultsFunc != nil && atomic.LoadUint64(&m.ExecutorResultsCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.ExecutorResults")
	}

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
	}

	if m.ProcessValidationResultsFunc != nil && atomic.LoadUint64(&m.ProcessValidationResultsCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.ProcessValidationResults")
	}

	if m.ValidateCaseBindFunc != nil && atomic.LoadUint64(&m.ValidateCaseBindCounter) == 0 {
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

	if m.ExecuteFunc != nil && atomic.LoadUint64(&m.ExecuteCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.Execute")
	}

	if m.ExecutorResultsFunc != nil && atomic.LoadUint64(&m.ExecutorResultsCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.ExecutorResults")
	}

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
	}

	if m.ProcessValidationResultsFunc != nil && atomic.LoadUint64(&m.ProcessValidationResultsCounter) == 0 {
		m.t.Fatal("Expected call to LogicRunnerMock.ProcessValidationResults")
	}

	if m.ValidateCaseBindFunc != nil && atomic.LoadUint64(&m.ValidateCaseBindCounter) == 0 {
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
		ok = ok && (m.ExecuteFunc == nil || atomic.LoadUint64(&m.ExecuteCounter) > 0)
		ok = ok && (m.ExecutorResultsFunc == nil || atomic.LoadUint64(&m.ExecutorResultsCounter) > 0)
		ok = ok && (m.OnPulseFunc == nil || atomic.LoadUint64(&m.OnPulseCounter) > 0)
		ok = ok && (m.ProcessValidationResultsFunc == nil || atomic.LoadUint64(&m.ProcessValidationResultsCounter) > 0)
		ok = ok && (m.ValidateCaseBindFunc == nil || atomic.LoadUint64(&m.ValidateCaseBindCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ExecuteFunc != nil && atomic.LoadUint64(&m.ExecuteCounter) == 0 {
				m.t.Error("Expected call to LogicRunnerMock.Execute")
			}

			if m.ExecutorResultsFunc != nil && atomic.LoadUint64(&m.ExecutorResultsCounter) == 0 {
				m.t.Error("Expected call to LogicRunnerMock.ExecutorResults")
			}

			if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
				m.t.Error("Expected call to LogicRunnerMock.OnPulse")
			}

			if m.ProcessValidationResultsFunc != nil && atomic.LoadUint64(&m.ProcessValidationResultsCounter) == 0 {
				m.t.Error("Expected call to LogicRunnerMock.ProcessValidationResults")
			}

			if m.ValidateCaseBindFunc != nil && atomic.LoadUint64(&m.ValidateCaseBindCounter) == 0 {
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

	if m.ExecuteFunc != nil && atomic.LoadUint64(&m.ExecuteCounter) == 0 {
		return false
	}

	if m.ExecutorResultsFunc != nil && atomic.LoadUint64(&m.ExecutorResultsCounter) == 0 {
		return false
	}

	if m.OnPulseFunc != nil && atomic.LoadUint64(&m.OnPulseCounter) == 0 {
		return false
	}

	if m.ProcessValidationResultsFunc != nil && atomic.LoadUint64(&m.ProcessValidationResultsCounter) == 0 {
		return false
	}

	if m.ValidateCaseBindFunc != nil && atomic.LoadUint64(&m.ValidateCaseBindCounter) == 0 {
		return false
	}

	return true
}

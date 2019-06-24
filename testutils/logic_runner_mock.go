package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LogicRunner" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//LogicRunnerMock implements github.com/insolar/insolar/insolar.LogicRunner
type LogicRunnerMock struct {
	t minimock.Tester

	GetExecutorFunc       func(p insolar.MachineType) (r insolar.MachineLogicExecutor, r1 error)
	GetExecutorCounter    uint64
	GetExecutorPreCounter uint64
	GetExecutorMock       mLogicRunnerMockGetExecutor

	OnPulseFunc       func(p context.Context, p1 insolar.Pulse) (r error)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mLogicRunnerMockOnPulse

	RegisterExecutorFunc       func(p insolar.MachineType, p1 insolar.MachineLogicExecutor) (r error)
	RegisterExecutorCounter    uint64
	RegisterExecutorPreCounter uint64
	RegisterExecutorMock       mLogicRunnerMockRegisterExecutor
}

//NewLogicRunnerMock returns a mock for github.com/insolar/insolar/insolar.LogicRunner
func NewLogicRunnerMock(t minimock.Tester) *LogicRunnerMock {
	m := &LogicRunnerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetExecutorMock = mLogicRunnerMockGetExecutor{mock: m}
	m.OnPulseMock = mLogicRunnerMockOnPulse{mock: m}
	m.RegisterExecutorMock = mLogicRunnerMockRegisterExecutor{mock: m}

	return m
}

type mLogicRunnerMockGetExecutor struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockGetExecutorExpectation
	expectationSeries []*LogicRunnerMockGetExecutorExpectation
}

type LogicRunnerMockGetExecutorExpectation struct {
	input  *LogicRunnerMockGetExecutorInput
	result *LogicRunnerMockGetExecutorResult
}

type LogicRunnerMockGetExecutorInput struct {
	p insolar.MachineType
}

type LogicRunnerMockGetExecutorResult struct {
	r  insolar.MachineLogicExecutor
	r1 error
}

//Expect specifies that invocation of LogicRunner.GetExecutor is expected from 1 to Infinity times
func (m *mLogicRunnerMockGetExecutor) Expect(p insolar.MachineType) *mLogicRunnerMockGetExecutor {
	m.mock.GetExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockGetExecutorExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockGetExecutorInput{p}
	return m
}

//Return specifies results of invocation of LogicRunner.GetExecutor
func (m *mLogicRunnerMockGetExecutor) Return(r insolar.MachineLogicExecutor, r1 error) *LogicRunnerMock {
	m.mock.GetExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockGetExecutorExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockGetExecutorResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.GetExecutor is expected once
func (m *mLogicRunnerMockGetExecutor) ExpectOnce(p insolar.MachineType) *LogicRunnerMockGetExecutorExpectation {
	m.mock.GetExecutorFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockGetExecutorExpectation{}
	expectation.input = &LogicRunnerMockGetExecutorInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockGetExecutorExpectation) Return(r insolar.MachineLogicExecutor, r1 error) {
	e.result = &LogicRunnerMockGetExecutorResult{r, r1}
}

//Set uses given function f as a mock of LogicRunner.GetExecutor method
func (m *mLogicRunnerMockGetExecutor) Set(f func(p insolar.MachineType) (r insolar.MachineLogicExecutor, r1 error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExecutorFunc = f
	return m.mock
}

//GetExecutor implements github.com/insolar/insolar/insolar.LogicRunner interface
func (m *LogicRunnerMock) GetExecutor(p insolar.MachineType) (r insolar.MachineLogicExecutor, r1 error) {
	counter := atomic.AddUint64(&m.GetExecutorPreCounter, 1)
	defer atomic.AddUint64(&m.GetExecutorCounter, 1)

	if len(m.GetExecutorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExecutorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.GetExecutor. %v", p)
			return
		}

		input := m.GetExecutorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockGetExecutorInput{p}, "LogicRunner.GetExecutor got unexpected parameters")

		result := m.GetExecutorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.GetExecutor")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetExecutorMock.mainExpectation != nil {

		input := m.GetExecutorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockGetExecutorInput{p}, "LogicRunner.GetExecutor got unexpected parameters")
		}

		result := m.GetExecutorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.GetExecutor")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetExecutorFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.GetExecutor. %v", p)
		return
	}

	return m.GetExecutorFunc(p)
}

//GetExecutorMinimockCounter returns a count of LogicRunnerMock.GetExecutorFunc invocations
func (m *LogicRunnerMock) GetExecutorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutorCounter)
}

//GetExecutorMinimockPreCounter returns the value of LogicRunnerMock.GetExecutor invocations
func (m *LogicRunnerMock) GetExecutorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutorPreCounter)
}

//GetExecutorFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) GetExecutorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetExecutorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetExecutorCounter) == uint64(len(m.GetExecutorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetExecutorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetExecutorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetExecutorFunc != nil {
		return atomic.LoadUint64(&m.GetExecutorCounter) > 0
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
	p1 insolar.Pulse
}

type LogicRunnerMockOnPulseResult struct {
	r error
}

//Expect specifies that invocation of LogicRunner.OnPulse is expected from 1 to Infinity times
func (m *mLogicRunnerMockOnPulse) Expect(p context.Context, p1 insolar.Pulse) *mLogicRunnerMockOnPulse {
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
func (m *mLogicRunnerMockOnPulse) ExpectOnce(p context.Context, p1 insolar.Pulse) *LogicRunnerMockOnPulseExpectation {
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
func (m *mLogicRunnerMockOnPulse) Set(f func(p context.Context, p1 insolar.Pulse) (r error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/insolar.LogicRunner interface
func (m *LogicRunnerMock) OnPulse(p context.Context, p1 insolar.Pulse) (r error) {
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

type mLogicRunnerMockRegisterExecutor struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockRegisterExecutorExpectation
	expectationSeries []*LogicRunnerMockRegisterExecutorExpectation
}

type LogicRunnerMockRegisterExecutorExpectation struct {
	input  *LogicRunnerMockRegisterExecutorInput
	result *LogicRunnerMockRegisterExecutorResult
}

type LogicRunnerMockRegisterExecutorInput struct {
	p  insolar.MachineType
	p1 insolar.MachineLogicExecutor
}

type LogicRunnerMockRegisterExecutorResult struct {
	r error
}

//Expect specifies that invocation of LogicRunner.RegisterExecutor is expected from 1 to Infinity times
func (m *mLogicRunnerMockRegisterExecutor) Expect(p insolar.MachineType, p1 insolar.MachineLogicExecutor) *mLogicRunnerMockRegisterExecutor {
	m.mock.RegisterExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockRegisterExecutorExpectation{}
	}
	m.mainExpectation.input = &LogicRunnerMockRegisterExecutorInput{p, p1}
	return m
}

//Return specifies results of invocation of LogicRunner.RegisterExecutor
func (m *mLogicRunnerMockRegisterExecutor) Return(r error) *LogicRunnerMock {
	m.mock.RegisterExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockRegisterExecutorExpectation{}
	}
	m.mainExpectation.result = &LogicRunnerMockRegisterExecutorResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.RegisterExecutor is expected once
func (m *mLogicRunnerMockRegisterExecutor) ExpectOnce(p insolar.MachineType, p1 insolar.MachineLogicExecutor) *LogicRunnerMockRegisterExecutorExpectation {
	m.mock.RegisterExecutorFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockRegisterExecutorExpectation{}
	expectation.input = &LogicRunnerMockRegisterExecutorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LogicRunnerMockRegisterExecutorExpectation) Return(r error) {
	e.result = &LogicRunnerMockRegisterExecutorResult{r}
}

//Set uses given function f as a mock of LogicRunner.RegisterExecutor method
func (m *mLogicRunnerMockRegisterExecutor) Set(f func(p insolar.MachineType, p1 insolar.MachineLogicExecutor) (r error)) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterExecutorFunc = f
	return m.mock
}

//RegisterExecutor implements github.com/insolar/insolar/insolar.LogicRunner interface
func (m *LogicRunnerMock) RegisterExecutor(p insolar.MachineType, p1 insolar.MachineLogicExecutor) (r error) {
	counter := atomic.AddUint64(&m.RegisterExecutorPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterExecutorCounter, 1)

	if len(m.RegisterExecutorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterExecutorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.RegisterExecutor. %v %v", p, p1)
			return
		}

		input := m.RegisterExecutorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LogicRunnerMockRegisterExecutorInput{p, p1}, "LogicRunner.RegisterExecutor got unexpected parameters")

		result := m.RegisterExecutorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.RegisterExecutor")
			return
		}

		r = result.r

		return
	}

	if m.RegisterExecutorMock.mainExpectation != nil {

		input := m.RegisterExecutorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LogicRunnerMockRegisterExecutorInput{p, p1}, "LogicRunner.RegisterExecutor got unexpected parameters")
		}

		result := m.RegisterExecutorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LogicRunnerMock.RegisterExecutor")
		}

		r = result.r

		return
	}

	if m.RegisterExecutorFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.RegisterExecutor. %v %v", p, p1)
		return
	}

	return m.RegisterExecutorFunc(p, p1)
}

//RegisterExecutorMinimockCounter returns a count of LogicRunnerMock.RegisterExecutorFunc invocations
func (m *LogicRunnerMock) RegisterExecutorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterExecutorCounter)
}

//RegisterExecutorMinimockPreCounter returns the value of LogicRunnerMock.RegisterExecutor invocations
func (m *LogicRunnerMock) RegisterExecutorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterExecutorPreCounter)
}

//RegisterExecutorFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) RegisterExecutorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterExecutorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterExecutorCounter) == uint64(len(m.RegisterExecutorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterExecutorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterExecutorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterExecutorFunc != nil {
		return atomic.LoadUint64(&m.RegisterExecutorCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicRunnerMock) ValidateCallCounters() {

	if !m.GetExecutorFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.GetExecutor")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
	}

	if !m.RegisterExecutorFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.RegisterExecutor")
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

	if !m.GetExecutorFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.GetExecutor")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.OnPulse")
	}

	if !m.RegisterExecutorFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.RegisterExecutor")
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
		ok = ok && m.GetExecutorFinished()
		ok = ok && m.OnPulseFinished()
		ok = ok && m.RegisterExecutorFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetExecutorFinished() {
				m.t.Error("Expected call to LogicRunnerMock.GetExecutor")
			}

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to LogicRunnerMock.OnPulse")
			}

			if !m.RegisterExecutorFinished() {
				m.t.Error("Expected call to LogicRunnerMock.RegisterExecutor")
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

	if !m.GetExecutorFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	if !m.RegisterExecutorFinished() {
		return false
	}

	return true
}

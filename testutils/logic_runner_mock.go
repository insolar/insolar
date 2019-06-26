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

	LRIFunc       func()
	LRICounter    uint64
	LRIPreCounter uint64
	LRIMock       mLogicRunnerMockLRI

	OnPulseFunc       func(p context.Context, p1 insolar.Pulse) (r error)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mLogicRunnerMockOnPulse
}

//NewLogicRunnerMock returns a mock for github.com/insolar/insolar/insolar.LogicRunner
func NewLogicRunnerMock(t minimock.Tester) *LogicRunnerMock {
	m := &LogicRunnerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LRIMock = mLogicRunnerMockLRI{mock: m}
	m.OnPulseMock = mLogicRunnerMockOnPulse{mock: m}

	return m
}

type mLogicRunnerMockLRI struct {
	mock              *LogicRunnerMock
	mainExpectation   *LogicRunnerMockLRIExpectation
	expectationSeries []*LogicRunnerMockLRIExpectation
}

type LogicRunnerMockLRIExpectation struct {
}

//Expect specifies that invocation of LogicRunner.LRI is expected from 1 to Infinity times
func (m *mLogicRunnerMockLRI) Expect() *mLogicRunnerMockLRI {
	m.mock.LRIFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockLRIExpectation{}
	}

	return m
}

//Return specifies results of invocation of LogicRunner.LRI
func (m *mLogicRunnerMockLRI) Return() *LogicRunnerMock {
	m.mock.LRIFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LogicRunnerMockLRIExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of LogicRunner.LRI is expected once
func (m *mLogicRunnerMockLRI) ExpectOnce() *LogicRunnerMockLRIExpectation {
	m.mock.LRIFunc = nil
	m.mainExpectation = nil

	expectation := &LogicRunnerMockLRIExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of LogicRunner.LRI method
func (m *mLogicRunnerMockLRI) Set(f func()) *LogicRunnerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LRIFunc = f
	return m.mock
}

//LRI implements github.com/insolar/insolar/insolar.LogicRunner interface
func (m *LogicRunnerMock) LRI() {
	counter := atomic.AddUint64(&m.LRIPreCounter, 1)
	defer atomic.AddUint64(&m.LRICounter, 1)

	if len(m.LRIMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LRIMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LogicRunnerMock.LRI.")
			return
		}

		return
	}

	if m.LRIMock.mainExpectation != nil {

		return
	}

	if m.LRIFunc == nil {
		m.t.Fatalf("Unexpected call to LogicRunnerMock.LRI.")
		return
	}

	m.LRIFunc()
}

//LRIMinimockCounter returns a count of LogicRunnerMock.LRIFunc invocations
func (m *LogicRunnerMock) LRIMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LRICounter)
}

//LRIMinimockPreCounter returns the value of LogicRunnerMock.LRI invocations
func (m *LogicRunnerMock) LRIMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LRIPreCounter)
}

//LRIFinished returns true if mock invocations count is ok
func (m *LogicRunnerMock) LRIFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LRIMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LRICounter) == uint64(len(m.LRIMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LRIMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LRICounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LRIFunc != nil {
		return atomic.LoadUint64(&m.LRICounter) > 0
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LogicRunnerMock) ValidateCallCounters() {

	if !m.LRIFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.LRI")
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

	if !m.LRIFinished() {
		m.t.Fatal("Expected call to LogicRunnerMock.LRI")
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
		ok = ok && m.LRIFinished()
		ok = ok && m.OnPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.LRIFinished() {
				m.t.Error("Expected call to LogicRunnerMock.LRI")
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

	if !m.LRIFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	return true
}

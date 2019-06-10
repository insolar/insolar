package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetCalculator" can be found in github.com/insolar/insolar/ledger/light/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetCalculatorMock implements github.com/insolar/insolar/ledger/light/executor.JetCalculator
type JetCalculatorMock struct {
	t minimock.Tester

	MineForPulseFunc       func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)
	MineForPulseCounter    uint64
	MineForPulsePreCounter uint64
	MineForPulseMock       mJetCalculatorMockMineForPulse
}

//NewJetCalculatorMock returns a mock for github.com/insolar/insolar/ledger/light/executor.JetCalculator
func NewJetCalculatorMock(t minimock.Tester) *JetCalculatorMock {
	m := &JetCalculatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.MineForPulseMock = mJetCalculatorMockMineForPulse{mock: m}

	return m
}

type mJetCalculatorMockMineForPulse struct {
	mock              *JetCalculatorMock
	mainExpectation   *JetCalculatorMockMineForPulseExpectation
	expectationSeries []*JetCalculatorMockMineForPulseExpectation
}

type JetCalculatorMockMineForPulseExpectation struct {
	input  *JetCalculatorMockMineForPulseInput
	result *JetCalculatorMockMineForPulseResult
}

type JetCalculatorMockMineForPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type JetCalculatorMockMineForPulseResult struct {
	r []insolar.JetID
}

//Expect specifies that invocation of JetCalculator.MineForPulse is expected from 1 to Infinity times
func (m *mJetCalculatorMockMineForPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mJetCalculatorMockMineForPulse {
	m.mock.MineForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCalculatorMockMineForPulseExpectation{}
	}
	m.mainExpectation.input = &JetCalculatorMockMineForPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of JetCalculator.MineForPulse
func (m *mJetCalculatorMockMineForPulse) Return(r []insolar.JetID) *JetCalculatorMock {
	m.mock.MineForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCalculatorMockMineForPulseExpectation{}
	}
	m.mainExpectation.result = &JetCalculatorMockMineForPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetCalculator.MineForPulse is expected once
func (m *mJetCalculatorMockMineForPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *JetCalculatorMockMineForPulseExpectation {
	m.mock.MineForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &JetCalculatorMockMineForPulseExpectation{}
	expectation.input = &JetCalculatorMockMineForPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetCalculatorMockMineForPulseExpectation) Return(r []insolar.JetID) {
	e.result = &JetCalculatorMockMineForPulseResult{r}
}

//Set uses given function f as a mock of JetCalculator.MineForPulse method
func (m *mJetCalculatorMockMineForPulse) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)) *JetCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MineForPulseFunc = f
	return m.mock
}

//MineForPulse implements github.com/insolar/insolar/ledger/light/executor.JetCalculator interface
func (m *JetCalculatorMock) MineForPulse(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID) {
	counter := atomic.AddUint64(&m.MineForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.MineForPulseCounter, 1)

	if len(m.MineForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MineForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetCalculatorMock.MineForPulse. %v %v", p, p1)
			return
		}

		input := m.MineForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetCalculatorMockMineForPulseInput{p, p1}, "JetCalculator.MineForPulse got unexpected parameters")

		result := m.MineForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetCalculatorMock.MineForPulse")
			return
		}

		r = result.r

		return
	}

	if m.MineForPulseMock.mainExpectation != nil {

		input := m.MineForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetCalculatorMockMineForPulseInput{p, p1}, "JetCalculator.MineForPulse got unexpected parameters")
		}

		result := m.MineForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetCalculatorMock.MineForPulse")
		}

		r = result.r

		return
	}

	if m.MineForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to JetCalculatorMock.MineForPulse. %v %v", p, p1)
		return
	}

	return m.MineForPulseFunc(p, p1)
}

//MineForPulseMinimockCounter returns a count of JetCalculatorMock.MineForPulseFunc invocations
func (m *JetCalculatorMock) MineForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MineForPulseCounter)
}

//MineForPulseMinimockPreCounter returns the value of JetCalculatorMock.MineForPulse invocations
func (m *JetCalculatorMock) MineForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MineForPulsePreCounter)
}

//MineForPulseFinished returns true if mock invocations count is ok
func (m *JetCalculatorMock) MineForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MineForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MineForPulseCounter) == uint64(len(m.MineForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MineForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MineForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MineForPulseFunc != nil {
		return atomic.LoadUint64(&m.MineForPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetCalculatorMock) ValidateCallCounters() {

	if !m.MineForPulseFinished() {
		m.t.Fatal("Expected call to JetCalculatorMock.MineForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetCalculatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetCalculatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetCalculatorMock) MinimockFinish() {

	if !m.MineForPulseFinished() {
		m.t.Fatal("Expected call to JetCalculatorMock.MineForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetCalculatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetCalculatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.MineForPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.MineForPulseFinished() {
				m.t.Error("Expected call to JetCalculatorMock.MineForPulse")
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
func (m *JetCalculatorMock) AllMocksCalled() bool {

	if !m.MineForPulseFinished() {
		return false
	}

	return true
}

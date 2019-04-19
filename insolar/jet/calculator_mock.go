package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Calculator" can be found in github.com/insolar/insolar/insolar/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//CalculatorMock implements github.com/insolar/insolar/insolar/jet.Calculator
type CalculatorMock struct {
	t minimock.Tester

	MineForPulseFunc       func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)
	MineForPulseCounter    uint64
	MineForPulsePreCounter uint64
	MineForPulseMock       mCalculatorMockMineForPulse
}

//NewCalculatorMock returns a mock for github.com/insolar/insolar/insolar/jet.Calculator
func NewCalculatorMock(t minimock.Tester) *CalculatorMock {
	m := &CalculatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.MineForPulseMock = mCalculatorMockMineForPulse{mock: m}

	return m
}

type mCalculatorMockMineForPulse struct {
	mock              *CalculatorMock
	mainExpectation   *CalculatorMockMineForPulseExpectation
	expectationSeries []*CalculatorMockMineForPulseExpectation
}

type CalculatorMockMineForPulseExpectation struct {
	input  *CalculatorMockMineForPulseInput
	result *CalculatorMockMineForPulseResult
}

type CalculatorMockMineForPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type CalculatorMockMineForPulseResult struct {
	r []insolar.JetID
}

//Expect specifies that invocation of Calculator.MineForPulse is expected from 1 to Infinity times
func (m *mCalculatorMockMineForPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mCalculatorMockMineForPulse {
	m.mock.MineForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockMineForPulseExpectation{}
	}
	m.mainExpectation.input = &CalculatorMockMineForPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of Calculator.MineForPulse
func (m *mCalculatorMockMineForPulse) Return(r []insolar.JetID) *CalculatorMock {
	m.mock.MineForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockMineForPulseExpectation{}
	}
	m.mainExpectation.result = &CalculatorMockMineForPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Calculator.MineForPulse is expected once
func (m *mCalculatorMockMineForPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *CalculatorMockMineForPulseExpectation {
	m.mock.MineForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &CalculatorMockMineForPulseExpectation{}
	expectation.input = &CalculatorMockMineForPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CalculatorMockMineForPulseExpectation) Return(r []insolar.JetID) {
	e.result = &CalculatorMockMineForPulseResult{r}
}

//Set uses given function f as a mock of Calculator.MineForPulse method
func (m *mCalculatorMockMineForPulse) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)) *CalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MineForPulseFunc = f
	return m.mock
}

//MineForPulse implements github.com/insolar/insolar/insolar/jet.Calculator interface
func (m *CalculatorMock) MineForPulse(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID) {
	counter := atomic.AddUint64(&m.MineForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.MineForPulseCounter, 1)

	if len(m.MineForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MineForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CalculatorMock.MineForPulse. %v %v", p, p1)
			return
		}

		input := m.MineForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CalculatorMockMineForPulseInput{p, p1}, "Calculator.MineForPulse got unexpected parameters")

		result := m.MineForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.MineForPulse")
			return
		}

		r = result.r

		return
	}

	if m.MineForPulseMock.mainExpectation != nil {

		input := m.MineForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CalculatorMockMineForPulseInput{p, p1}, "Calculator.MineForPulse got unexpected parameters")
		}

		result := m.MineForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.MineForPulse")
		}

		r = result.r

		return
	}

	if m.MineForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to CalculatorMock.MineForPulse. %v %v", p, p1)
		return
	}

	return m.MineForPulseFunc(p, p1)
}

//MineForPulseMinimockCounter returns a count of CalculatorMock.MineForPulseFunc invocations
func (m *CalculatorMock) MineForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MineForPulseCounter)
}

//MineForPulseMinimockPreCounter returns the value of CalculatorMock.MineForPulse invocations
func (m *CalculatorMock) MineForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MineForPulsePreCounter)
}

//MineForPulseFinished returns true if mock invocations count is ok
func (m *CalculatorMock) MineForPulseFinished() bool {
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
func (m *CalculatorMock) ValidateCallCounters() {

	if !m.MineForPulseFinished() {
		m.t.Fatal("Expected call to CalculatorMock.MineForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CalculatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CalculatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CalculatorMock) MinimockFinish() {

	if !m.MineForPulseFinished() {
		m.t.Fatal("Expected call to CalculatorMock.MineForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CalculatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CalculatorMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to CalculatorMock.MineForPulse")
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
func (m *CalculatorMock) AllMocksCalled() bool {

	if !m.MineForPulseFinished() {
		return false
	}

	return true
}

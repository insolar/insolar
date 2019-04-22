package pulse

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Calculator" can be found in github.com/insolar/insolar/insolar/pulse
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//CalculatorMock implements github.com/insolar/insolar/insolar/pulse.Calculator
type CalculatorMock struct {
	t minimock.Tester

	BackwardsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error)
	BackwardsCounter    uint64
	BackwardsPreCounter uint64
	BackwardsMock       mCalculatorMockBackwards

	ForwardsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error)
	ForwardsCounter    uint64
	ForwardsPreCounter uint64
	ForwardsMock       mCalculatorMockForwards
}

//NewCalculatorMock returns a mock for github.com/insolar/insolar/insolar/pulse.Calculator
func NewCalculatorMock(t minimock.Tester) *CalculatorMock {
	m := &CalculatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BackwardsMock = mCalculatorMockBackwards{mock: m}
	m.ForwardsMock = mCalculatorMockForwards{mock: m}

	return m
}

type mCalculatorMockBackwards struct {
	mock              *CalculatorMock
	mainExpectation   *CalculatorMockBackwardsExpectation
	expectationSeries []*CalculatorMockBackwardsExpectation
}

type CalculatorMockBackwardsExpectation struct {
	input  *CalculatorMockBackwardsInput
	result *CalculatorMockBackwardsResult
}

type CalculatorMockBackwardsInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 int
}

type CalculatorMockBackwardsResult struct {
	r  insolar.Pulse
	r1 error
}

//Expect specifies that invocation of Calculator.Backwards is expected from 1 to Infinity times
func (m *mCalculatorMockBackwards) Expect(p context.Context, p1 insolar.PulseNumber, p2 int) *mCalculatorMockBackwards {
	m.mock.BackwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockBackwardsExpectation{}
	}
	m.mainExpectation.input = &CalculatorMockBackwardsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Calculator.Backwards
func (m *mCalculatorMockBackwards) Return(r insolar.Pulse, r1 error) *CalculatorMock {
	m.mock.BackwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockBackwardsExpectation{}
	}
	m.mainExpectation.result = &CalculatorMockBackwardsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Calculator.Backwards is expected once
func (m *mCalculatorMockBackwards) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 int) *CalculatorMockBackwardsExpectation {
	m.mock.BackwardsFunc = nil
	m.mainExpectation = nil

	expectation := &CalculatorMockBackwardsExpectation{}
	expectation.input = &CalculatorMockBackwardsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CalculatorMockBackwardsExpectation) Return(r insolar.Pulse, r1 error) {
	e.result = &CalculatorMockBackwardsResult{r, r1}
}

//Set uses given function f as a mock of Calculator.Backwards method
func (m *mCalculatorMockBackwards) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error)) *CalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BackwardsFunc = f
	return m.mock
}

//Backwards implements github.com/insolar/insolar/insolar/pulse.Calculator interface
func (m *CalculatorMock) Backwards(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.BackwardsPreCounter, 1)
	defer atomic.AddUint64(&m.BackwardsCounter, 1)

	if len(m.BackwardsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BackwardsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CalculatorMock.Backwards. %v %v %v", p, p1, p2)
			return
		}

		input := m.BackwardsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CalculatorMockBackwardsInput{p, p1, p2}, "Calculator.Backwards got unexpected parameters")

		result := m.BackwardsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.Backwards")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BackwardsMock.mainExpectation != nil {

		input := m.BackwardsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CalculatorMockBackwardsInput{p, p1, p2}, "Calculator.Backwards got unexpected parameters")
		}

		result := m.BackwardsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.Backwards")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BackwardsFunc == nil {
		m.t.Fatalf("Unexpected call to CalculatorMock.Backwards. %v %v %v", p, p1, p2)
		return
	}

	return m.BackwardsFunc(p, p1, p2)
}

//BackwardsMinimockCounter returns a count of CalculatorMock.BackwardsFunc invocations
func (m *CalculatorMock) BackwardsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BackwardsCounter)
}

//BackwardsMinimockPreCounter returns the value of CalculatorMock.Backwards invocations
func (m *CalculatorMock) BackwardsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BackwardsPreCounter)
}

//BackwardsFinished returns true if mock invocations count is ok
func (m *CalculatorMock) BackwardsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.BackwardsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.BackwardsCounter) == uint64(len(m.BackwardsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.BackwardsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.BackwardsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.BackwardsFunc != nil {
		return atomic.LoadUint64(&m.BackwardsCounter) > 0
	}

	return true
}

type mCalculatorMockForwards struct {
	mock              *CalculatorMock
	mainExpectation   *CalculatorMockForwardsExpectation
	expectationSeries []*CalculatorMockForwardsExpectation
}

type CalculatorMockForwardsExpectation struct {
	input  *CalculatorMockForwardsInput
	result *CalculatorMockForwardsResult
}

type CalculatorMockForwardsInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 int
}

type CalculatorMockForwardsResult struct {
	r  insolar.Pulse
	r1 error
}

//Expect specifies that invocation of Calculator.Forwards is expected from 1 to Infinity times
func (m *mCalculatorMockForwards) Expect(p context.Context, p1 insolar.PulseNumber, p2 int) *mCalculatorMockForwards {
	m.mock.ForwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockForwardsExpectation{}
	}
	m.mainExpectation.input = &CalculatorMockForwardsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Calculator.Forwards
func (m *mCalculatorMockForwards) Return(r insolar.Pulse, r1 error) *CalculatorMock {
	m.mock.ForwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockForwardsExpectation{}
	}
	m.mainExpectation.result = &CalculatorMockForwardsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Calculator.Forwards is expected once
func (m *mCalculatorMockForwards) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 int) *CalculatorMockForwardsExpectation {
	m.mock.ForwardsFunc = nil
	m.mainExpectation = nil

	expectation := &CalculatorMockForwardsExpectation{}
	expectation.input = &CalculatorMockForwardsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CalculatorMockForwardsExpectation) Return(r insolar.Pulse, r1 error) {
	e.result = &CalculatorMockForwardsResult{r, r1}
}

//Set uses given function f as a mock of Calculator.Forwards method
func (m *mCalculatorMockForwards) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error)) *CalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForwardsFunc = f
	return m.mock
}

//Forwards implements github.com/insolar/insolar/insolar/pulse.Calculator interface
func (m *CalculatorMock) Forwards(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.ForwardsPreCounter, 1)
	defer atomic.AddUint64(&m.ForwardsCounter, 1)

	if len(m.ForwardsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForwardsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CalculatorMock.Forwards. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForwardsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CalculatorMockForwardsInput{p, p1, p2}, "Calculator.Forwards got unexpected parameters")

		result := m.ForwardsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.Forwards")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForwardsMock.mainExpectation != nil {

		input := m.ForwardsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CalculatorMockForwardsInput{p, p1, p2}, "Calculator.Forwards got unexpected parameters")
		}

		result := m.ForwardsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.Forwards")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForwardsFunc == nil {
		m.t.Fatalf("Unexpected call to CalculatorMock.Forwards. %v %v %v", p, p1, p2)
		return
	}

	return m.ForwardsFunc(p, p1, p2)
}

//ForwardsMinimockCounter returns a count of CalculatorMock.ForwardsFunc invocations
func (m *CalculatorMock) ForwardsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForwardsCounter)
}

//ForwardsMinimockPreCounter returns the value of CalculatorMock.Forwards invocations
func (m *CalculatorMock) ForwardsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForwardsPreCounter)
}

//ForwardsFinished returns true if mock invocations count is ok
func (m *CalculatorMock) ForwardsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForwardsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForwardsCounter) == uint64(len(m.ForwardsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForwardsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForwardsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForwardsFunc != nil {
		return atomic.LoadUint64(&m.ForwardsCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CalculatorMock) ValidateCallCounters() {

	if !m.BackwardsFinished() {
		m.t.Fatal("Expected call to CalculatorMock.Backwards")
	}

	if !m.ForwardsFinished() {
		m.t.Fatal("Expected call to CalculatorMock.Forwards")
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

	if !m.BackwardsFinished() {
		m.t.Fatal("Expected call to CalculatorMock.Backwards")
	}

	if !m.ForwardsFinished() {
		m.t.Fatal("Expected call to CalculatorMock.Forwards")
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
		ok = ok && m.BackwardsFinished()
		ok = ok && m.ForwardsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.BackwardsFinished() {
				m.t.Error("Expected call to CalculatorMock.Backwards")
			}

			if !m.ForwardsFinished() {
				m.t.Error("Expected call to CalculatorMock.Forwards")
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

	if !m.BackwardsFinished() {
		return false
	}

	if !m.ForwardsFinished() {
		return false
	}

	return true
}

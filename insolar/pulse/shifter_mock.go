package pulse

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Shifter" can be found in github.com/insolar/insolar/insolar/pulse
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ShifterMock implements github.com/insolar/insolar/insolar/pulse.Shifter
type ShifterMock struct {
	t minimock.Tester

	ShiftFunc       func(p context.Context, p1 insolar.PulseNumber) (r error)
	ShiftCounter    uint64
	ShiftPreCounter uint64
	ShiftMock       mShifterMockShift
}

//NewShifterMock returns a mock for github.com/insolar/insolar/insolar/pulse.Shifter
func NewShifterMock(t minimock.Tester) *ShifterMock {
	m := &ShifterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ShiftMock = mShifterMockShift{mock: m}

	return m
}

type mShifterMockShift struct {
	mock              *ShifterMock
	mainExpectation   *ShifterMockShiftExpectation
	expectationSeries []*ShifterMockShiftExpectation
}

type ShifterMockShiftExpectation struct {
	input  *ShifterMockShiftInput
	result *ShifterMockShiftResult
}

type ShifterMockShiftInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type ShifterMockShiftResult struct {
	r error
}

//Expect specifies that invocation of Shifter.Shift is expected from 1 to Infinity times
func (m *mShifterMockShift) Expect(p context.Context, p1 insolar.PulseNumber) *mShifterMockShift {
	m.mock.ShiftFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ShifterMockShiftExpectation{}
	}
	m.mainExpectation.input = &ShifterMockShiftInput{p, p1}
	return m
}

//Return specifies results of invocation of Shifter.Shift
func (m *mShifterMockShift) Return(r error) *ShifterMock {
	m.mock.ShiftFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ShifterMockShiftExpectation{}
	}
	m.mainExpectation.result = &ShifterMockShiftResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Shifter.Shift is expected once
func (m *mShifterMockShift) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *ShifterMockShiftExpectation {
	m.mock.ShiftFunc = nil
	m.mainExpectation = nil

	expectation := &ShifterMockShiftExpectation{}
	expectation.input = &ShifterMockShiftInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ShifterMockShiftExpectation) Return(r error) {
	e.result = &ShifterMockShiftResult{r}
}

//Set uses given function f as a mock of Shifter.Shift method
func (m *mShifterMockShift) Set(f func(p context.Context, p1 insolar.PulseNumber) (r error)) *ShifterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ShiftFunc = f
	return m.mock
}

//Shift implements github.com/insolar/insolar/insolar/pulse.Shifter interface
func (m *ShifterMock) Shift(p context.Context, p1 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.ShiftPreCounter, 1)
	defer atomic.AddUint64(&m.ShiftCounter, 1)

	if len(m.ShiftMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ShiftMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ShifterMock.Shift. %v %v", p, p1)
			return
		}

		input := m.ShiftMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ShifterMockShiftInput{p, p1}, "Shifter.Shift got unexpected parameters")

		result := m.ShiftMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ShifterMock.Shift")
			return
		}

		r = result.r

		return
	}

	if m.ShiftMock.mainExpectation != nil {

		input := m.ShiftMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ShifterMockShiftInput{p, p1}, "Shifter.Shift got unexpected parameters")
		}

		result := m.ShiftMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ShifterMock.Shift")
		}

		r = result.r

		return
	}

	if m.ShiftFunc == nil {
		m.t.Fatalf("Unexpected call to ShifterMock.Shift. %v %v", p, p1)
		return
	}

	return m.ShiftFunc(p, p1)
}

//ShiftMinimockCounter returns a count of ShifterMock.ShiftFunc invocations
func (m *ShifterMock) ShiftMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ShiftCounter)
}

//ShiftMinimockPreCounter returns the value of ShifterMock.Shift invocations
func (m *ShifterMock) ShiftMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ShiftPreCounter)
}

//ShiftFinished returns true if mock invocations count is ok
func (m *ShifterMock) ShiftFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ShiftMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ShiftCounter) == uint64(len(m.ShiftMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ShiftMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ShiftCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ShiftFunc != nil {
		return atomic.LoadUint64(&m.ShiftCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ShifterMock) ValidateCallCounters() {

	if !m.ShiftFinished() {
		m.t.Fatal("Expected call to ShifterMock.Shift")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ShifterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ShifterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ShifterMock) MinimockFinish() {

	if !m.ShiftFinished() {
		m.t.Fatal("Expected call to ShifterMock.Shift")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ShifterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ShifterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ShiftFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ShiftFinished() {
				m.t.Error("Expected call to ShifterMock.Shift")
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
func (m *ShifterMock) AllMocksCalled() bool {

	if !m.ShiftFinished() {
		return false
	}

	return true
}

package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseCalculator" can be found in github.com/insolar/insolar/network/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseCalculatorMock implements github.com/insolar/insolar/network/storage.PulseCalculator
type PulseCalculatorMock struct {
	t minimock.Tester

	BackwardsFunc       func(p context.Context, p1 core.PulseNumber, p2 int) (r core.Pulse, r1 error)
	BackwardsCounter    uint64
	BackwardsPreCounter uint64
	BackwardsMock       mPulseCalculatorMockBackwards

	ForwardsFunc       func(p context.Context, p1 core.PulseNumber, p2 int) (r core.Pulse, r1 error)
	ForwardsCounter    uint64
	ForwardsPreCounter uint64
	ForwardsMock       mPulseCalculatorMockForwards
}

//NewPulseCalculatorMock returns a mock for github.com/insolar/insolar/network/storage.PulseCalculator
func NewPulseCalculatorMock(t minimock.Tester) *PulseCalculatorMock {
	m := &PulseCalculatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BackwardsMock = mPulseCalculatorMockBackwards{mock: m}
	m.ForwardsMock = mPulseCalculatorMockForwards{mock: m}

	return m
}

type mPulseCalculatorMockBackwards struct {
	mock              *PulseCalculatorMock
	mainExpectation   *PulseCalculatorMockBackwardsExpectation
	expectationSeries []*PulseCalculatorMockBackwardsExpectation
}

type PulseCalculatorMockBackwardsExpectation struct {
	input  *PulseCalculatorMockBackwardsInput
	result *PulseCalculatorMockBackwardsResult
}

type PulseCalculatorMockBackwardsInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 int
}

type PulseCalculatorMockBackwardsResult struct {
	r  core.Pulse
	r1 error
}

//Expect specifies that invocation of PulseCalculator.Backwards is expected from 1 to Infinity times
func (m *mPulseCalculatorMockBackwards) Expect(p context.Context, p1 core.PulseNumber, p2 int) *mPulseCalculatorMockBackwards {
	m.mock.BackwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseCalculatorMockBackwardsExpectation{}
	}
	m.mainExpectation.input = &PulseCalculatorMockBackwardsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PulseCalculator.Backwards
func (m *mPulseCalculatorMockBackwards) Return(r core.Pulse, r1 error) *PulseCalculatorMock {
	m.mock.BackwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseCalculatorMockBackwardsExpectation{}
	}
	m.mainExpectation.result = &PulseCalculatorMockBackwardsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseCalculator.Backwards is expected once
func (m *mPulseCalculatorMockBackwards) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 int) *PulseCalculatorMockBackwardsExpectation {
	m.mock.BackwardsFunc = nil
	m.mainExpectation = nil

	expectation := &PulseCalculatorMockBackwardsExpectation{}
	expectation.input = &PulseCalculatorMockBackwardsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseCalculatorMockBackwardsExpectation) Return(r core.Pulse, r1 error) {
	e.result = &PulseCalculatorMockBackwardsResult{r, r1}
}

//Set uses given function f as a mock of PulseCalculator.Backwards method
func (m *mPulseCalculatorMockBackwards) Set(f func(p context.Context, p1 core.PulseNumber, p2 int) (r core.Pulse, r1 error)) *PulseCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BackwardsFunc = f
	return m.mock
}

//Backwards implements github.com/insolar/insolar/network/storage.PulseCalculator interface
func (m *PulseCalculatorMock) Backwards(p context.Context, p1 core.PulseNumber, p2 int) (r core.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.BackwardsPreCounter, 1)
	defer atomic.AddUint64(&m.BackwardsCounter, 1)

	if len(m.BackwardsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BackwardsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseCalculatorMock.Backwards. %v %v %v", p, p1, p2)
			return
		}

		input := m.BackwardsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseCalculatorMockBackwardsInput{p, p1, p2}, "PulseCalculator.Backwards got unexpected parameters")

		result := m.BackwardsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseCalculatorMock.Backwards")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BackwardsMock.mainExpectation != nil {

		input := m.BackwardsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseCalculatorMockBackwardsInput{p, p1, p2}, "PulseCalculator.Backwards got unexpected parameters")
		}

		result := m.BackwardsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseCalculatorMock.Backwards")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BackwardsFunc == nil {
		m.t.Fatalf("Unexpected call to PulseCalculatorMock.Backwards. %v %v %v", p, p1, p2)
		return
	}

	return m.BackwardsFunc(p, p1, p2)
}

//BackwardsMinimockCounter returns a count of PulseCalculatorMock.BackwardsFunc invocations
func (m *PulseCalculatorMock) BackwardsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BackwardsCounter)
}

//BackwardsMinimockPreCounter returns the value of PulseCalculatorMock.Backwards invocations
func (m *PulseCalculatorMock) BackwardsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BackwardsPreCounter)
}

//BackwardsFinished returns true if mock invocations count is ok
func (m *PulseCalculatorMock) BackwardsFinished() bool {
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

type mPulseCalculatorMockForwards struct {
	mock              *PulseCalculatorMock
	mainExpectation   *PulseCalculatorMockForwardsExpectation
	expectationSeries []*PulseCalculatorMockForwardsExpectation
}

type PulseCalculatorMockForwardsExpectation struct {
	input  *PulseCalculatorMockForwardsInput
	result *PulseCalculatorMockForwardsResult
}

type PulseCalculatorMockForwardsInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 int
}

type PulseCalculatorMockForwardsResult struct {
	r  core.Pulse
	r1 error
}

//Expect specifies that invocation of PulseCalculator.Forwards is expected from 1 to Infinity times
func (m *mPulseCalculatorMockForwards) Expect(p context.Context, p1 core.PulseNumber, p2 int) *mPulseCalculatorMockForwards {
	m.mock.ForwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseCalculatorMockForwardsExpectation{}
	}
	m.mainExpectation.input = &PulseCalculatorMockForwardsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PulseCalculator.Forwards
func (m *mPulseCalculatorMockForwards) Return(r core.Pulse, r1 error) *PulseCalculatorMock {
	m.mock.ForwardsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseCalculatorMockForwardsExpectation{}
	}
	m.mainExpectation.result = &PulseCalculatorMockForwardsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseCalculator.Forwards is expected once
func (m *mPulseCalculatorMockForwards) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 int) *PulseCalculatorMockForwardsExpectation {
	m.mock.ForwardsFunc = nil
	m.mainExpectation = nil

	expectation := &PulseCalculatorMockForwardsExpectation{}
	expectation.input = &PulseCalculatorMockForwardsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseCalculatorMockForwardsExpectation) Return(r core.Pulse, r1 error) {
	e.result = &PulseCalculatorMockForwardsResult{r, r1}
}

//Set uses given function f as a mock of PulseCalculator.Forwards method
func (m *mPulseCalculatorMockForwards) Set(f func(p context.Context, p1 core.PulseNumber, p2 int) (r core.Pulse, r1 error)) *PulseCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForwardsFunc = f
	return m.mock
}

//Forwards implements github.com/insolar/insolar/network/storage.PulseCalculator interface
func (m *PulseCalculatorMock) Forwards(p context.Context, p1 core.PulseNumber, p2 int) (r core.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.ForwardsPreCounter, 1)
	defer atomic.AddUint64(&m.ForwardsCounter, 1)

	if len(m.ForwardsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForwardsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseCalculatorMock.Forwards. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForwardsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseCalculatorMockForwardsInput{p, p1, p2}, "PulseCalculator.Forwards got unexpected parameters")

		result := m.ForwardsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseCalculatorMock.Forwards")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForwardsMock.mainExpectation != nil {

		input := m.ForwardsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseCalculatorMockForwardsInput{p, p1, p2}, "PulseCalculator.Forwards got unexpected parameters")
		}

		result := m.ForwardsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseCalculatorMock.Forwards")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForwardsFunc == nil {
		m.t.Fatalf("Unexpected call to PulseCalculatorMock.Forwards. %v %v %v", p, p1, p2)
		return
	}

	return m.ForwardsFunc(p, p1, p2)
}

//ForwardsMinimockCounter returns a count of PulseCalculatorMock.ForwardsFunc invocations
func (m *PulseCalculatorMock) ForwardsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForwardsCounter)
}

//ForwardsMinimockPreCounter returns the value of PulseCalculatorMock.Forwards invocations
func (m *PulseCalculatorMock) ForwardsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForwardsPreCounter)
}

//ForwardsFinished returns true if mock invocations count is ok
func (m *PulseCalculatorMock) ForwardsFinished() bool {
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
func (m *PulseCalculatorMock) ValidateCallCounters() {

	if !m.BackwardsFinished() {
		m.t.Fatal("Expected call to PulseCalculatorMock.Backwards")
	}

	if !m.ForwardsFinished() {
		m.t.Fatal("Expected call to PulseCalculatorMock.Forwards")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseCalculatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseCalculatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseCalculatorMock) MinimockFinish() {

	if !m.BackwardsFinished() {
		m.t.Fatal("Expected call to PulseCalculatorMock.Backwards")
	}

	if !m.ForwardsFinished() {
		m.t.Fatal("Expected call to PulseCalculatorMock.Forwards")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseCalculatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseCalculatorMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to PulseCalculatorMock.Backwards")
			}

			if !m.ForwardsFinished() {
				m.t.Error("Expected call to PulseCalculatorMock.Forwards")
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
func (m *PulseCalculatorMock) AllMocksCalled() bool {

	if !m.BackwardsFinished() {
		return false
	}

	if !m.ForwardsFinished() {
		return false
	}

	return true
}

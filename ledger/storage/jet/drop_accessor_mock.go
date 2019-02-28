package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DropAccessor" can be found in github.com/insolar/insolar/ledger/storage/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//DropAccessorMock implements github.com/insolar/insolar/ledger/storage/jet.DropAccessor
type DropAccessorMock struct {
	t minimock.Tester

	ForPulseFunc       func(p context.Context, p1 core.JetID, p2 core.PulseNumber) (r JetDrop, r1 error)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mDropAccessorMockForPulse
}

//NewDropAccessorMock returns a mock for github.com/insolar/insolar/ledger/storage/jet.DropAccessor
func NewDropAccessorMock(t minimock.Tester) *DropAccessorMock {
	m := &DropAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseMock = mDropAccessorMockForPulse{mock: m}

	return m
}

type mDropAccessorMockForPulse struct {
	mock              *DropAccessorMock
	mainExpectation   *DropAccessorMockForPulseExpectation
	expectationSeries []*DropAccessorMockForPulseExpectation
}

type DropAccessorMockForPulseExpectation struct {
	input  *DropAccessorMockForPulseInput
	result *DropAccessorMockForPulseResult
}

type DropAccessorMockForPulseInput struct {
	p  context.Context
	p1 core.JetID
	p2 core.PulseNumber
}

type DropAccessorMockForPulseResult struct {
	r  JetDrop
	r1 error
}

//Expect specifies that invocation of DropAccessor.ForPulse is expected from 1 to Infinity times
func (m *mDropAccessorMockForPulse) Expect(p context.Context, p1 core.JetID, p2 core.PulseNumber) *mDropAccessorMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.input = &DropAccessorMockForPulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DropAccessor.ForPulse
func (m *mDropAccessorMockForPulse) Return(r JetDrop, r1 error) *DropAccessorMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.result = &DropAccessorMockForPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DropAccessor.ForPulse is expected once
func (m *mDropAccessorMockForPulse) ExpectOnce(p context.Context, p1 core.JetID, p2 core.PulseNumber) *DropAccessorMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &DropAccessorMockForPulseExpectation{}
	expectation.input = &DropAccessorMockForPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropAccessorMockForPulseExpectation) Return(r JetDrop, r1 error) {
	e.result = &DropAccessorMockForPulseResult{r, r1}
}

//Set uses given function f as a mock of DropAccessor.ForPulse method
func (m *mDropAccessorMockForPulse) Set(f func(p context.Context, p1 core.JetID, p2 core.PulseNumber) (r JetDrop, r1 error)) *DropAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/storage/jet.DropAccessor interface
func (m *DropAccessorMock) ForPulse(p context.Context, p1 core.JetID, p2 core.PulseNumber) (r JetDrop, r1 error) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropAccessorMock.ForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropAccessorMockForPulseInput{p, p1, p2}, "DropAccessor.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropAccessorMock.ForPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropAccessorMockForPulseInput{p, p1, p2}, "DropAccessor.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropAccessorMock.ForPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to DropAccessorMock.ForPulse. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPulseFunc(p, p1, p2)
}

//ForPulseMinimockCounter returns a count of DropAccessorMock.ForPulseFunc invocations
func (m *DropAccessorMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

//ForPulseMinimockPreCounter returns the value of DropAccessorMock.ForPulse invocations
func (m *DropAccessorMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

//ForPulseFinished returns true if mock invocations count is ok
func (m *DropAccessorMock) ForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseCounter) == uint64(len(m.ForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseFunc != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DropAccessorMock) ValidateCallCounters() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to DropAccessorMock.ForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DropAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DropAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DropAccessorMock) MinimockFinish() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to DropAccessorMock.ForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DropAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DropAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseFinished() {
				m.t.Error("Expected call to DropAccessorMock.ForPulse")
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
func (m *DropAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseFinished() {
		return false
	}

	return true
}

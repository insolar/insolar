package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseManager" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseManagerMock implements github.com/insolar/insolar/core.PulseManager
type PulseManagerMock struct {
	t minimock.Tester

	SetFunc       func(p context.Context, p1 core.Pulse, p2 bool) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mPulseManagerMockSet
}

//NewPulseManagerMock returns a mock for github.com/insolar/insolar/core.PulseManager
func NewPulseManagerMock(t minimock.Tester) *PulseManagerMock {
	m := &PulseManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetMock = mPulseManagerMockSet{mock: m}

	return m
}

type mPulseManagerMockSet struct {
	mock              *PulseManagerMock
	mainExpectation   *PulseManagerMockSetExpectation
	expectationSeries []*PulseManagerMockSetExpectation
}

type PulseManagerMockSetExpectation struct {
	input  *PulseManagerMockSetInput
	result *PulseManagerMockSetResult
}

type PulseManagerMockSetInput struct {
	p  context.Context
	p1 core.Pulse
	p2 bool
}

type PulseManagerMockSetResult struct {
	r error
}

//Expect specifies that invocation of PulseManager.Set is expected from 1 to Infinity times
func (m *mPulseManagerMockSet) Expect(p context.Context, p1 core.Pulse, p2 bool) *mPulseManagerMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseManagerMockSetExpectation{}
	}
	m.mainExpectation.input = &PulseManagerMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PulseManager.Set
func (m *mPulseManagerMockSet) Return(r error) *PulseManagerMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseManagerMockSetExpectation{}
	}
	m.mainExpectation.result = &PulseManagerMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseManager.Set is expected once
func (m *mPulseManagerMockSet) ExpectOnce(p context.Context, p1 core.Pulse, p2 bool) *PulseManagerMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &PulseManagerMockSetExpectation{}
	expectation.input = &PulseManagerMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseManagerMockSetExpectation) Return(r error) {
	e.result = &PulseManagerMockSetResult{r}
}

//Set uses given function f as a mock of PulseManager.Set method
func (m *mPulseManagerMockSet) Set(f func(p context.Context, p1 core.Pulse, p2 bool) (r error)) *PulseManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/core.PulseManager interface
func (m *PulseManagerMock) Set(p context.Context, p1 core.Pulse, p2 bool) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseManagerMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseManagerMockSetInput{p, p1, p2}, "PulseManager.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseManagerMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseManagerMockSetInput{p, p1, p2}, "PulseManager.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseManagerMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to PulseManagerMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of PulseManagerMock.SetFunc invocations
func (m *PulseManagerMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of PulseManagerMock.Set invocations
func (m *PulseManagerMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *PulseManagerMock) SetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCounter) == uint64(len(m.SetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFunc != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseManagerMock) ValidateCallCounters() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to PulseManagerMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseManagerMock) MinimockFinish() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to PulseManagerMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetFinished() {
				m.t.Error("Expected call to PulseManagerMock.Set")
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
func (m *PulseManagerMock) AllMocksCalled() bool {

	if !m.SetFinished() {
		return false
	}

	return true
}

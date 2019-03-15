package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseAppender" can be found in github.com/insolar/insolar/network/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseAppenderMock implements github.com/insolar/insolar/network/storage.PulseAppender
type PulseAppenderMock struct {
	t minimock.Tester

	AppendFunc       func(p context.Context, p1 core.Pulse) (r error)
	AppendCounter    uint64
	AppendPreCounter uint64
	AppendMock       mPulseAppenderMockAppend
}

//NewPulseAppenderMock returns a mock for github.com/insolar/insolar/network/storage.PulseAppender
func NewPulseAppenderMock(t minimock.Tester) *PulseAppenderMock {
	m := &PulseAppenderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AppendMock = mPulseAppenderMockAppend{mock: m}

	return m
}

type mPulseAppenderMockAppend struct {
	mock              *PulseAppenderMock
	mainExpectation   *PulseAppenderMockAppendExpectation
	expectationSeries []*PulseAppenderMockAppendExpectation
}

type PulseAppenderMockAppendExpectation struct {
	input  *PulseAppenderMockAppendInput
	result *PulseAppenderMockAppendResult
}

type PulseAppenderMockAppendInput struct {
	p  context.Context
	p1 core.Pulse
}

type PulseAppenderMockAppendResult struct {
	r error
}

//Expect specifies that invocation of PulseAppender.Append is expected from 1 to Infinity times
func (m *mPulseAppenderMockAppend) Expect(p context.Context, p1 core.Pulse) *mPulseAppenderMockAppend {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseAppenderMockAppendExpectation{}
	}
	m.mainExpectation.input = &PulseAppenderMockAppendInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseAppender.Append
func (m *mPulseAppenderMockAppend) Return(r error) *PulseAppenderMock {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseAppenderMockAppendExpectation{}
	}
	m.mainExpectation.result = &PulseAppenderMockAppendResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseAppender.Append is expected once
func (m *mPulseAppenderMockAppend) ExpectOnce(p context.Context, p1 core.Pulse) *PulseAppenderMockAppendExpectation {
	m.mock.AppendFunc = nil
	m.mainExpectation = nil

	expectation := &PulseAppenderMockAppendExpectation{}
	expectation.input = &PulseAppenderMockAppendInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseAppenderMockAppendExpectation) Return(r error) {
	e.result = &PulseAppenderMockAppendResult{r}
}

//Set uses given function f as a mock of PulseAppender.Append method
func (m *mPulseAppenderMockAppend) Set(f func(p context.Context, p1 core.Pulse) (r error)) *PulseAppenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AppendFunc = f
	return m.mock
}

//Append implements github.com/insolar/insolar/network/storage.PulseAppender interface
func (m *PulseAppenderMock) Append(p context.Context, p1 core.Pulse) (r error) {
	counter := atomic.AddUint64(&m.AppendPreCounter, 1)
	defer atomic.AddUint64(&m.AppendCounter, 1)

	if len(m.AppendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AppendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseAppenderMock.Append. %v %v", p, p1)
			return
		}

		input := m.AppendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseAppenderMockAppendInput{p, p1}, "PulseAppender.Append got unexpected parameters")

		result := m.AppendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseAppenderMock.Append")
			return
		}

		r = result.r

		return
	}

	if m.AppendMock.mainExpectation != nil {

		input := m.AppendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseAppenderMockAppendInput{p, p1}, "PulseAppender.Append got unexpected parameters")
		}

		result := m.AppendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseAppenderMock.Append")
		}

		r = result.r

		return
	}

	if m.AppendFunc == nil {
		m.t.Fatalf("Unexpected call to PulseAppenderMock.Append. %v %v", p, p1)
		return
	}

	return m.AppendFunc(p, p1)
}

//AppendMinimockCounter returns a count of PulseAppenderMock.AppendFunc invocations
func (m *PulseAppenderMock) AppendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AppendCounter)
}

//AppendMinimockPreCounter returns the value of PulseAppenderMock.Append invocations
func (m *PulseAppenderMock) AppendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AppendPreCounter)
}

//AppendFinished returns true if mock invocations count is ok
func (m *PulseAppenderMock) AppendFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AppendMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AppendCounter) == uint64(len(m.AppendMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AppendMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AppendCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AppendFunc != nil {
		return atomic.LoadUint64(&m.AppendCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseAppenderMock) ValidateCallCounters() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to PulseAppenderMock.Append")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseAppenderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseAppenderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseAppenderMock) MinimockFinish() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to PulseAppenderMock.Append")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseAppenderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseAppenderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AppendFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AppendFinished() {
				m.t.Error("Expected call to PulseAppenderMock.Append")
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
func (m *PulseAppenderMock) AllMocksCalled() bool {

	if !m.AppendFinished() {
		return false
	}

	return true
}

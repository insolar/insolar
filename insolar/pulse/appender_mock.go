package pulse

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Appender" can be found in github.com/insolar/insolar/insolar/pulse
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//AppenderMock implements github.com/insolar/insolar/insolar/pulse.Appender
type AppenderMock struct {
	t minimock.Tester

	AppendFunc       func(p context.Context, p1 insolar.Pulse) (r error)
	AppendCounter    uint64
	AppendPreCounter uint64
	AppendMock       mAppenderMockAppend
}

//NewAppenderMock returns a mock for github.com/insolar/insolar/insolar/pulse.Appender
func NewAppenderMock(t minimock.Tester) *AppenderMock {
	m := &AppenderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AppendMock = mAppenderMockAppend{mock: m}

	return m
}

type mAppenderMockAppend struct {
	mock              *AppenderMock
	mainExpectation   *AppenderMockAppendExpectation
	expectationSeries []*AppenderMockAppendExpectation
}

type AppenderMockAppendExpectation struct {
	input  *AppenderMockAppendInput
	result *AppenderMockAppendResult
}

type AppenderMockAppendInput struct {
	p  context.Context
	p1 insolar.Pulse
}

type AppenderMockAppendResult struct {
	r error
}

//Expect specifies that invocation of Appender.Append is expected from 1 to Infinity times
func (m *mAppenderMockAppend) Expect(p context.Context, p1 insolar.Pulse) *mAppenderMockAppend {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AppenderMockAppendExpectation{}
	}
	m.mainExpectation.input = &AppenderMockAppendInput{p, p1}
	return m
}

//Return specifies results of invocation of Appender.Append
func (m *mAppenderMockAppend) Return(r error) *AppenderMock {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AppenderMockAppendExpectation{}
	}
	m.mainExpectation.result = &AppenderMockAppendResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Appender.Append is expected once
func (m *mAppenderMockAppend) ExpectOnce(p context.Context, p1 insolar.Pulse) *AppenderMockAppendExpectation {
	m.mock.AppendFunc = nil
	m.mainExpectation = nil

	expectation := &AppenderMockAppendExpectation{}
	expectation.input = &AppenderMockAppendInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AppenderMockAppendExpectation) Return(r error) {
	e.result = &AppenderMockAppendResult{r}
}

//Set uses given function f as a mock of Appender.Append method
func (m *mAppenderMockAppend) Set(f func(p context.Context, p1 insolar.Pulse) (r error)) *AppenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AppendFunc = f
	return m.mock
}

//Append implements github.com/insolar/insolar/insolar/pulse.Appender interface
func (m *AppenderMock) Append(p context.Context, p1 insolar.Pulse) (r error) {
	counter := atomic.AddUint64(&m.AppendPreCounter, 1)
	defer atomic.AddUint64(&m.AppendCounter, 1)

	if len(m.AppendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AppendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AppenderMock.Append. %v %v", p, p1)
			return
		}

		input := m.AppendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AppenderMockAppendInput{p, p1}, "Appender.Append got unexpected parameters")

		result := m.AppendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AppenderMock.Append")
			return
		}

		r = result.r

		return
	}

	if m.AppendMock.mainExpectation != nil {

		input := m.AppendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AppenderMockAppendInput{p, p1}, "Appender.Append got unexpected parameters")
		}

		result := m.AppendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AppenderMock.Append")
		}

		r = result.r

		return
	}

	if m.AppendFunc == nil {
		m.t.Fatalf("Unexpected call to AppenderMock.Append. %v %v", p, p1)
		return
	}

	return m.AppendFunc(p, p1)
}

//AppendMinimockCounter returns a count of AppenderMock.AppendFunc invocations
func (m *AppenderMock) AppendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AppendCounter)
}

//AppendMinimockPreCounter returns the value of AppenderMock.Append invocations
func (m *AppenderMock) AppendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AppendPreCounter)
}

//AppendFinished returns true if mock invocations count is ok
func (m *AppenderMock) AppendFinished() bool {
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
func (m *AppenderMock) ValidateCallCounters() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to AppenderMock.Append")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AppenderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AppenderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AppenderMock) MinimockFinish() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to AppenderMock.Append")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AppenderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AppenderMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to AppenderMock.Append")
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
func (m *AppenderMock) AllMocksCalled() bool {

	if !m.AppendFinished() {
		return false
	}

	return true
}

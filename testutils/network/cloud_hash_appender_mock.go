package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CloudHashAppender" can be found in github.com/insolar/insolar/network/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//CloudHashAppenderMock implements github.com/insolar/insolar/network/storage.CloudHashAppender
type CloudHashAppenderMock struct {
	t minimock.Tester

	AppendFunc       func(p context.Context, p1 insolar.PulseNumber, p2 []byte) (r error)
	AppendCounter    uint64
	AppendPreCounter uint64
	AppendMock       mCloudHashAppenderMockAppend
}

//NewCloudHashAppenderMock returns a mock for github.com/insolar/insolar/network/storage.CloudHashAppender
func NewCloudHashAppenderMock(t minimock.Tester) *CloudHashAppenderMock {
	m := &CloudHashAppenderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AppendMock = mCloudHashAppenderMockAppend{mock: m}

	return m
}

type mCloudHashAppenderMockAppend struct {
	mock              *CloudHashAppenderMock
	mainExpectation   *CloudHashAppenderMockAppendExpectation
	expectationSeries []*CloudHashAppenderMockAppendExpectation
}

type CloudHashAppenderMockAppendExpectation struct {
	input  *CloudHashAppenderMockAppendInput
	result *CloudHashAppenderMockAppendResult
}

type CloudHashAppenderMockAppendInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 []byte
}

type CloudHashAppenderMockAppendResult struct {
	r error
}

//Expect specifies that invocation of CloudHashAppender.Append is expected from 1 to Infinity times
func (m *mCloudHashAppenderMockAppend) Expect(p context.Context, p1 insolar.PulseNumber, p2 []byte) *mCloudHashAppenderMockAppend {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudHashAppenderMockAppendExpectation{}
	}
	m.mainExpectation.input = &CloudHashAppenderMockAppendInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of CloudHashAppender.Append
func (m *mCloudHashAppenderMockAppend) Return(r error) *CloudHashAppenderMock {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudHashAppenderMockAppendExpectation{}
	}
	m.mainExpectation.result = &CloudHashAppenderMockAppendResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudHashAppender.Append is expected once
func (m *mCloudHashAppenderMockAppend) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 []byte) *CloudHashAppenderMockAppendExpectation {
	m.mock.AppendFunc = nil
	m.mainExpectation = nil

	expectation := &CloudHashAppenderMockAppendExpectation{}
	expectation.input = &CloudHashAppenderMockAppendInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudHashAppenderMockAppendExpectation) Return(r error) {
	e.result = &CloudHashAppenderMockAppendResult{r}
}

//Set uses given function f as a mock of CloudHashAppender.Append method
func (m *mCloudHashAppenderMockAppend) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 []byte) (r error)) *CloudHashAppenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AppendFunc = f
	return m.mock
}

//Append implements github.com/insolar/insolar/network/storage.CloudHashAppender interface
func (m *CloudHashAppenderMock) Append(p context.Context, p1 insolar.PulseNumber, p2 []byte) (r error) {
	counter := atomic.AddUint64(&m.AppendPreCounter, 1)
	defer atomic.AddUint64(&m.AppendCounter, 1)

	if len(m.AppendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AppendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudHashAppenderMock.Append. %v %v %v", p, p1, p2)
			return
		}

		input := m.AppendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CloudHashAppenderMockAppendInput{p, p1, p2}, "CloudHashAppender.Append got unexpected parameters")

		result := m.AppendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudHashAppenderMock.Append")
			return
		}

		r = result.r

		return
	}

	if m.AppendMock.mainExpectation != nil {

		input := m.AppendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CloudHashAppenderMockAppendInput{p, p1, p2}, "CloudHashAppender.Append got unexpected parameters")
		}

		result := m.AppendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudHashAppenderMock.Append")
		}

		r = result.r

		return
	}

	if m.AppendFunc == nil {
		m.t.Fatalf("Unexpected call to CloudHashAppenderMock.Append. %v %v %v", p, p1, p2)
		return
	}

	return m.AppendFunc(p, p1, p2)
}

//AppendMinimockCounter returns a count of CloudHashAppenderMock.AppendFunc invocations
func (m *CloudHashAppenderMock) AppendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AppendCounter)
}

//AppendMinimockPreCounter returns the value of CloudHashAppenderMock.Append invocations
func (m *CloudHashAppenderMock) AppendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AppendPreCounter)
}

//AppendFinished returns true if mock invocations count is ok
func (m *CloudHashAppenderMock) AppendFinished() bool {
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
func (m *CloudHashAppenderMock) ValidateCallCounters() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to CloudHashAppenderMock.Append")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CloudHashAppenderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CloudHashAppenderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CloudHashAppenderMock) MinimockFinish() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to CloudHashAppenderMock.Append")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CloudHashAppenderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CloudHashAppenderMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to CloudHashAppenderMock.Append")
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
func (m *CloudHashAppenderMock) AllMocksCalled() bool {

	if !m.AppendFinished() {
		return false
	}

	return true
}

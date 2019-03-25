package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SnapshotAppender" can be found in github.com/insolar/insolar/network/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	network "github.com/insolar/insolar/network"

	testify_assert "github.com/stretchr/testify/assert"
)

//SnapshotAppenderMock implements github.com/insolar/insolar/network/storage.SnapshotAppender
type SnapshotAppenderMock struct {
	t minimock.Tester

	AppendFunc       func(p context.Context, p1 insolar.PulseNumber, p2 network.Snapshot) (r error)
	AppendCounter    uint64
	AppendPreCounter uint64
	AppendMock       mSnapshotAppenderMockAppend
}

//NewSnapshotAppenderMock returns a mock for github.com/insolar/insolar/network/storage.SnapshotAppender
func NewSnapshotAppenderMock(t minimock.Tester) *SnapshotAppenderMock {
	m := &SnapshotAppenderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AppendMock = mSnapshotAppenderMockAppend{mock: m}

	return m
}

type mSnapshotAppenderMockAppend struct {
	mock              *SnapshotAppenderMock
	mainExpectation   *SnapshotAppenderMockAppendExpectation
	expectationSeries []*SnapshotAppenderMockAppendExpectation
}

type SnapshotAppenderMockAppendExpectation struct {
	input  *SnapshotAppenderMockAppendInput
	result *SnapshotAppenderMockAppendResult
}

type SnapshotAppenderMockAppendInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 network.Snapshot
}

type SnapshotAppenderMockAppendResult struct {
	r error
}

//Expect specifies that invocation of SnapshotAppender.Append is expected from 1 to Infinity times
func (m *mSnapshotAppenderMockAppend) Expect(p context.Context, p1 insolar.PulseNumber, p2 network.Snapshot) *mSnapshotAppenderMockAppend {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SnapshotAppenderMockAppendExpectation{}
	}
	m.mainExpectation.input = &SnapshotAppenderMockAppendInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of SnapshotAppender.Append
func (m *mSnapshotAppenderMockAppend) Return(r error) *SnapshotAppenderMock {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SnapshotAppenderMockAppendExpectation{}
	}
	m.mainExpectation.result = &SnapshotAppenderMockAppendResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SnapshotAppender.Append is expected once
func (m *mSnapshotAppenderMockAppend) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 network.Snapshot) *SnapshotAppenderMockAppendExpectation {
	m.mock.AppendFunc = nil
	m.mainExpectation = nil

	expectation := &SnapshotAppenderMockAppendExpectation{}
	expectation.input = &SnapshotAppenderMockAppendInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SnapshotAppenderMockAppendExpectation) Return(r error) {
	e.result = &SnapshotAppenderMockAppendResult{r}
}

//Set uses given function f as a mock of SnapshotAppender.Append method
func (m *mSnapshotAppenderMockAppend) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 network.Snapshot) (r error)) *SnapshotAppenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AppendFunc = f
	return m.mock
}

//Append implements github.com/insolar/insolar/network/storage.SnapshotAppender interface
func (m *SnapshotAppenderMock) Append(p context.Context, p1 insolar.PulseNumber, p2 network.Snapshot) (r error) {
	counter := atomic.AddUint64(&m.AppendPreCounter, 1)
	defer atomic.AddUint64(&m.AppendCounter, 1)

	if len(m.AppendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AppendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SnapshotAppenderMock.Append. %v %v %v", p, p1, p2)
			return
		}

		input := m.AppendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SnapshotAppenderMockAppendInput{p, p1, p2}, "SnapshotAppender.Append got unexpected parameters")

		result := m.AppendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SnapshotAppenderMock.Append")
			return
		}

		r = result.r

		return
	}

	if m.AppendMock.mainExpectation != nil {

		input := m.AppendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SnapshotAppenderMockAppendInput{p, p1, p2}, "SnapshotAppender.Append got unexpected parameters")
		}

		result := m.AppendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SnapshotAppenderMock.Append")
		}

		r = result.r

		return
	}

	if m.AppendFunc == nil {
		m.t.Fatalf("Unexpected call to SnapshotAppenderMock.Append. %v %v %v", p, p1, p2)
		return
	}

	return m.AppendFunc(p, p1, p2)
}

//AppendMinimockCounter returns a count of SnapshotAppenderMock.AppendFunc invocations
func (m *SnapshotAppenderMock) AppendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AppendCounter)
}

//AppendMinimockPreCounter returns the value of SnapshotAppenderMock.Append invocations
func (m *SnapshotAppenderMock) AppendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AppendPreCounter)
}

//AppendFinished returns true if mock invocations count is ok
func (m *SnapshotAppenderMock) AppendFinished() bool {
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
func (m *SnapshotAppenderMock) ValidateCallCounters() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to SnapshotAppenderMock.Append")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SnapshotAppenderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SnapshotAppenderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SnapshotAppenderMock) MinimockFinish() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to SnapshotAppenderMock.Append")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SnapshotAppenderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SnapshotAppenderMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to SnapshotAppenderMock.Append")
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
func (m *SnapshotAppenderMock) AllMocksCalled() bool {

	if !m.AppendFinished() {
		return false
	}

	return true
}

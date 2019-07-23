package replica

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Target" can be found in github.com/insolar/insolar/ledger/heavy/replica
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//TargetMock implements github.com/insolar/insolar/ledger/heavy/replica.Target
type TargetMock struct {
	t minimock.Tester

	NotifyFunc       func(p context.Context, p1 insolar.PulseNumber) (r error)
	NotifyCounter    uint64
	NotifyPreCounter uint64
	NotifyMock       mTargetMockNotify
}

//NewTargetMock returns a mock for github.com/insolar/insolar/ledger/heavy/replica.Target
func NewTargetMock(t minimock.Tester) *TargetMock {
	m := &TargetMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.NotifyMock = mTargetMockNotify{mock: m}

	return m
}

type mTargetMockNotify struct {
	mock              *TargetMock
	mainExpectation   *TargetMockNotifyExpectation
	expectationSeries []*TargetMockNotifyExpectation
}

type TargetMockNotifyExpectation struct {
	input  *TargetMockNotifyInput
	result *TargetMockNotifyResult
}

type TargetMockNotifyInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type TargetMockNotifyResult struct {
	r error
}

//Expect specifies that invocation of Target.Notify is expected from 1 to Infinity times
func (m *mTargetMockNotify) Expect(p context.Context, p1 insolar.PulseNumber) *mTargetMockNotify {
	m.mock.NotifyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TargetMockNotifyExpectation{}
	}
	m.mainExpectation.input = &TargetMockNotifyInput{p, p1}
	return m
}

//Return specifies results of invocation of Target.Notify
func (m *mTargetMockNotify) Return(r error) *TargetMock {
	m.mock.NotifyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TargetMockNotifyExpectation{}
	}
	m.mainExpectation.result = &TargetMockNotifyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Target.Notify is expected once
func (m *mTargetMockNotify) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *TargetMockNotifyExpectation {
	m.mock.NotifyFunc = nil
	m.mainExpectation = nil

	expectation := &TargetMockNotifyExpectation{}
	expectation.input = &TargetMockNotifyInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TargetMockNotifyExpectation) Return(r error) {
	e.result = &TargetMockNotifyResult{r}
}

//Set uses given function f as a mock of Target.Notify method
func (m *mTargetMockNotify) Set(f func(p context.Context, p1 insolar.PulseNumber) (r error)) *TargetMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NotifyFunc = f
	return m.mock
}

//Notify implements github.com/insolar/insolar/ledger/heavy/replica.Target interface
func (m *TargetMock) Notify(p context.Context, p1 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.NotifyPreCounter, 1)
	defer atomic.AddUint64(&m.NotifyCounter, 1)

	if len(m.NotifyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NotifyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TargetMock.Notify. %v %v", p, p1)
			return
		}

		input := m.NotifyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TargetMockNotifyInput{p, p1}, "Target.Notify got unexpected parameters")

		result := m.NotifyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TargetMock.Notify")
			return
		}

		r = result.r

		return
	}

	if m.NotifyMock.mainExpectation != nil {

		input := m.NotifyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TargetMockNotifyInput{p, p1}, "Target.Notify got unexpected parameters")
		}

		result := m.NotifyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TargetMock.Notify")
		}

		r = result.r

		return
	}

	if m.NotifyFunc == nil {
		m.t.Fatalf("Unexpected call to TargetMock.Notify. %v %v", p, p1)
		return
	}

	return m.NotifyFunc(p, p1)
}

//NotifyMinimockCounter returns a count of TargetMock.NotifyFunc invocations
func (m *TargetMock) NotifyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NotifyCounter)
}

//NotifyMinimockPreCounter returns the value of TargetMock.Notify invocations
func (m *TargetMock) NotifyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NotifyPreCounter)
}

//NotifyFinished returns true if mock invocations count is ok
func (m *TargetMock) NotifyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NotifyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NotifyCounter) == uint64(len(m.NotifyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NotifyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NotifyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NotifyFunc != nil {
		return atomic.LoadUint64(&m.NotifyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TargetMock) ValidateCallCounters() {

	if !m.NotifyFinished() {
		m.t.Fatal("Expected call to TargetMock.Notify")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TargetMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TargetMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TargetMock) MinimockFinish() {

	if !m.NotifyFinished() {
		m.t.Fatal("Expected call to TargetMock.Notify")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TargetMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TargetMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.NotifyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.NotifyFinished() {
				m.t.Error("Expected call to TargetMock.Notify")
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
func (m *TargetMock) AllMocksCalled() bool {

	if !m.NotifyFinished() {
		return false
	}

	return true
}

package replication

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Cleaner" can be found in github.com/insolar/insolar/ledger/light/replication
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//CleanerMock implements github.com/insolar/insolar/ledger/light/replication.Cleaner
type CleanerMock struct {
	t minimock.Tester

	NotifyAboutPulseFunc       func(p context.Context, p1 insolar.PulseNumber)
	NotifyAboutPulseCounter    uint64
	NotifyAboutPulsePreCounter uint64
	NotifyAboutPulseMock       mCleanerMockNotifyAboutPulse
}

//NewCleanerMock returns a mock for github.com/insolar/insolar/ledger/light/replication.Cleaner
func NewCleanerMock(t minimock.Tester) *CleanerMock {
	m := &CleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.NotifyAboutPulseMock = mCleanerMockNotifyAboutPulse{mock: m}

	return m
}

type mCleanerMockNotifyAboutPulse struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockNotifyAboutPulseExpectation
	expectationSeries []*CleanerMockNotifyAboutPulseExpectation
}

type CleanerMockNotifyAboutPulseExpectation struct {
	input *CleanerMockNotifyAboutPulseInput
}

type CleanerMockNotifyAboutPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of Cleaner.NotifyAboutPulse is expected from 1 to Infinity times
func (m *mCleanerMockNotifyAboutPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mCleanerMockNotifyAboutPulse {
	m.mock.NotifyAboutPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockNotifyAboutPulseExpectation{}
	}
	m.mainExpectation.input = &CleanerMockNotifyAboutPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of Cleaner.NotifyAboutPulse
func (m *mCleanerMockNotifyAboutPulse) Return() *CleanerMock {
	m.mock.NotifyAboutPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockNotifyAboutPulseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Cleaner.NotifyAboutPulse is expected once
func (m *mCleanerMockNotifyAboutPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *CleanerMockNotifyAboutPulseExpectation {
	m.mock.NotifyAboutPulseFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockNotifyAboutPulseExpectation{}
	expectation.input = &CleanerMockNotifyAboutPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Cleaner.NotifyAboutPulse method
func (m *mCleanerMockNotifyAboutPulse) Set(f func(p context.Context, p1 insolar.PulseNumber)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NotifyAboutPulseFunc = f
	return m.mock
}

//NotifyAboutPulse implements github.com/insolar/insolar/ledger/light/replication.Cleaner interface
func (m *CleanerMock) NotifyAboutPulse(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.NotifyAboutPulsePreCounter, 1)
	defer atomic.AddUint64(&m.NotifyAboutPulseCounter, 1)

	if len(m.NotifyAboutPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NotifyAboutPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.NotifyAboutPulse. %v %v", p, p1)
			return
		}

		input := m.NotifyAboutPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockNotifyAboutPulseInput{p, p1}, "Cleaner.NotifyAboutPulse got unexpected parameters")

		return
	}

	if m.NotifyAboutPulseMock.mainExpectation != nil {

		input := m.NotifyAboutPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockNotifyAboutPulseInput{p, p1}, "Cleaner.NotifyAboutPulse got unexpected parameters")
		}

		return
	}

	if m.NotifyAboutPulseFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.NotifyAboutPulse. %v %v", p, p1)
		return
	}

	m.NotifyAboutPulseFunc(p, p1)
}

//NotifyAboutPulseMinimockCounter returns a count of CleanerMock.NotifyAboutPulseFunc invocations
func (m *CleanerMock) NotifyAboutPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NotifyAboutPulseCounter)
}

//NotifyAboutPulseMinimockPreCounter returns the value of CleanerMock.NotifyAboutPulse invocations
func (m *CleanerMock) NotifyAboutPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NotifyAboutPulsePreCounter)
}

//NotifyAboutPulseFinished returns true if mock invocations count is ok
func (m *CleanerMock) NotifyAboutPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NotifyAboutPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NotifyAboutPulseCounter) == uint64(len(m.NotifyAboutPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NotifyAboutPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NotifyAboutPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NotifyAboutPulseFunc != nil {
		return atomic.LoadUint64(&m.NotifyAboutPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CleanerMock) ValidateCallCounters() {

	if !m.NotifyAboutPulseFinished() {
		m.t.Fatal("Expected call to CleanerMock.NotifyAboutPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CleanerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CleanerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CleanerMock) MinimockFinish() {

	if !m.NotifyAboutPulseFinished() {
		m.t.Fatal("Expected call to CleanerMock.NotifyAboutPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.NotifyAboutPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.NotifyAboutPulseFinished() {
				m.t.Error("Expected call to CleanerMock.NotifyAboutPulse")
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
func (m *CleanerMock) AllMocksCalled() bool {

	if !m.NotifyAboutPulseFinished() {
		return false
	}

	return true
}

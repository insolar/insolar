package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordCleaner" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//RecordCleanerMock implements github.com/insolar/insolar/ledger/object.RecordCleaner
type RecordCleanerMock struct {
	t minimock.Tester

	DeleteForPNFunc       func(p context.Context, p1 insolar.PulseNumber)
	DeleteForPNCounter    uint64
	DeleteForPNPreCounter uint64
	DeleteForPNMock       mRecordCleanerMockDeleteForPN
}

//NewRecordCleanerMock returns a mock for github.com/insolar/insolar/ledger/object.RecordCleaner
func NewRecordCleanerMock(t minimock.Tester) *RecordCleanerMock {
	m := &RecordCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteForPNMock = mRecordCleanerMockDeleteForPN{mock: m}

	return m
}

type mRecordCleanerMockDeleteForPN struct {
	mock              *RecordCleanerMock
	mainExpectation   *RecordCleanerMockDeleteForPNExpectation
	expectationSeries []*RecordCleanerMockDeleteForPNExpectation
}

type RecordCleanerMockDeleteForPNExpectation struct {
	input *RecordCleanerMockDeleteForPNInput
}

type RecordCleanerMockDeleteForPNInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of RecordCleaner.DeleteForPN is expected from 1 to Infinity times
func (m *mRecordCleanerMockDeleteForPN) Expect(p context.Context, p1 insolar.PulseNumber) *mRecordCleanerMockDeleteForPN {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCleanerMockDeleteForPNExpectation{}
	}
	m.mainExpectation.input = &RecordCleanerMockDeleteForPNInput{p, p1}
	return m
}

//Return specifies results of invocation of RecordCleaner.DeleteForPN
func (m *mRecordCleanerMockDeleteForPN) Return() *RecordCleanerMock {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCleanerMockDeleteForPNExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecordCleaner.DeleteForPN is expected once
func (m *mRecordCleanerMockDeleteForPN) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *RecordCleanerMockDeleteForPNExpectation {
	m.mock.DeleteForPNFunc = nil
	m.mainExpectation = nil

	expectation := &RecordCleanerMockDeleteForPNExpectation{}
	expectation.input = &RecordCleanerMockDeleteForPNInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecordCleaner.DeleteForPN method
func (m *mRecordCleanerMockDeleteForPN) Set(f func(p context.Context, p1 insolar.PulseNumber)) *RecordCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteForPNFunc = f
	return m.mock
}

//DeleteForPN implements github.com/insolar/insolar/ledger/object.RecordCleaner interface
func (m *RecordCleanerMock) DeleteForPN(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.DeleteForPNPreCounter, 1)
	defer atomic.AddUint64(&m.DeleteForPNCounter, 1)

	if len(m.DeleteForPNMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteForPNMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordCleanerMock.DeleteForPN. %v %v", p, p1)
			return
		}

		input := m.DeleteForPNMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordCleanerMockDeleteForPNInput{p, p1}, "RecordCleaner.DeleteForPN got unexpected parameters")

		return
	}

	if m.DeleteForPNMock.mainExpectation != nil {

		input := m.DeleteForPNMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordCleanerMockDeleteForPNInput{p, p1}, "RecordCleaner.DeleteForPN got unexpected parameters")
		}

		return
	}

	if m.DeleteForPNFunc == nil {
		m.t.Fatalf("Unexpected call to RecordCleanerMock.DeleteForPN. %v %v", p, p1)
		return
	}

	m.DeleteForPNFunc(p, p1)
}

//DeleteForPNMinimockCounter returns a count of RecordCleanerMock.DeleteForPNFunc invocations
func (m *RecordCleanerMock) DeleteForPNMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNCounter)
}

//DeleteForPNMinimockPreCounter returns the value of RecordCleanerMock.DeleteForPN invocations
func (m *RecordCleanerMock) DeleteForPNMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNPreCounter)
}

//DeleteForPNFinished returns true if mock invocations count is ok
func (m *RecordCleanerMock) DeleteForPNFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeleteForPNMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeleteForPNCounter) == uint64(len(m.DeleteForPNMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeleteForPNMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeleteForPNCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeleteForPNFunc != nil {
		return atomic.LoadUint64(&m.DeleteForPNCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordCleanerMock) ValidateCallCounters() {

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to RecordCleanerMock.DeleteForPN")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordCleanerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecordCleanerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecordCleanerMock) MinimockFinish() {

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to RecordCleanerMock.DeleteForPN")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecordCleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RecordCleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.DeleteForPNFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeleteForPNFinished() {
				m.t.Error("Expected call to RecordCleanerMock.DeleteForPN")
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
func (m *RecordCleanerMock) AllMocksCalled() bool {

	if !m.DeleteForPNFinished() {
		return false
	}

	return true
}

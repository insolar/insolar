package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordCleaner" can be found in github.com/insolar/insolar/ledger/storage/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//RecordCleanerMock implements github.com/insolar/insolar/ledger/storage/object.RecordCleaner
type RecordCleanerMock struct {
	t minimock.Tester

	RemoveFunc       func(p context.Context, p1 insolar.PulseNumber)
	RemoveCounter    uint64
	RemovePreCounter uint64
	RemoveMock       mRecordCleanerMockRemove
}

//NewRecordCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage/object.RecordCleaner
func NewRecordCleanerMock(t minimock.Tester) *RecordCleanerMock {
	m := &RecordCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveMock = mRecordCleanerMockRemove{mock: m}

	return m
}

type mRecordCleanerMockRemove struct {
	mock              *RecordCleanerMock
	mainExpectation   *RecordCleanerMockRemoveExpectation
	expectationSeries []*RecordCleanerMockRemoveExpectation
}

type RecordCleanerMockRemoveExpectation struct {
	input *RecordCleanerMockRemoveInput
}

type RecordCleanerMockRemoveInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of RecordCleaner.Remove is expected from 1 to Infinity times
func (m *mRecordCleanerMockRemove) Expect(p context.Context, p1 insolar.PulseNumber) *mRecordCleanerMockRemove {
	m.mock.RemoveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCleanerMockRemoveExpectation{}
	}
	m.mainExpectation.input = &RecordCleanerMockRemoveInput{p, p1}
	return m
}

//Return specifies results of invocation of RecordCleaner.Remove
func (m *mRecordCleanerMockRemove) Return() *RecordCleanerMock {
	m.mock.RemoveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCleanerMockRemoveExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecordCleaner.Remove is expected once
func (m *mRecordCleanerMockRemove) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *RecordCleanerMockRemoveExpectation {
	m.mock.RemoveFunc = nil
	m.mainExpectation = nil

	expectation := &RecordCleanerMockRemoveExpectation{}
	expectation.input = &RecordCleanerMockRemoveInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecordCleaner.Remove method
func (m *mRecordCleanerMockRemove) Set(f func(p context.Context, p1 insolar.PulseNumber)) *RecordCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveFunc = f
	return m.mock
}

//Remove implements github.com/insolar/insolar/ledger/storage/object.RecordCleaner interface
func (m *RecordCleanerMock) Remove(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.RemovePreCounter, 1)
	defer atomic.AddUint64(&m.RemoveCounter, 1)

	if len(m.RemoveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordCleanerMock.Remove. %v %v", p, p1)
			return
		}

		input := m.RemoveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordCleanerMockRemoveInput{p, p1}, "RecordCleaner.Remove got unexpected parameters")

		return
	}

	if m.RemoveMock.mainExpectation != nil {

		input := m.RemoveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordCleanerMockRemoveInput{p, p1}, "RecordCleaner.Remove got unexpected parameters")
		}

		return
	}

	if m.RemoveFunc == nil {
		m.t.Fatalf("Unexpected call to RecordCleanerMock.Remove. %v %v", p, p1)
		return
	}

	m.RemoveFunc(p, p1)
}

//RemoveMinimockCounter returns a count of RecordCleanerMock.RemoveFunc invocations
func (m *RecordCleanerMock) RemoveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveCounter)
}

//RemoveMinimockPreCounter returns the value of RecordCleanerMock.Remove invocations
func (m *RecordCleanerMock) RemoveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePreCounter)
}

//RemoveFinished returns true if mock invocations count is ok
func (m *RecordCleanerMock) RemoveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveCounter) == uint64(len(m.RemoveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveFunc != nil {
		return atomic.LoadUint64(&m.RemoveCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordCleanerMock) ValidateCallCounters() {

	if !m.RemoveFinished() {
		m.t.Fatal("Expected call to RecordCleanerMock.Remove")
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

	if !m.RemoveFinished() {
		m.t.Fatal("Expected call to RecordCleanerMock.Remove")
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
		ok = ok && m.RemoveFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveFinished() {
				m.t.Error("Expected call to RecordCleanerMock.Remove")
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

	if !m.RemoveFinished() {
		return false
	}

	return true
}

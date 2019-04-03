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

	RemoveUntilFunc       func(p context.Context, p1 insolar.PulseNumber)
	RemoveUntilCounter    uint64
	RemoveUntilPreCounter uint64
	RemoveUntilMock       mRecordCleanerMockRemoveUntil
}

//NewRecordCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage/object.RecordCleaner
func NewRecordCleanerMock(t minimock.Tester) *RecordCleanerMock {
	m := &RecordCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveUntilMock = mRecordCleanerMockRemoveUntil{mock: m}

	return m
}

type mRecordCleanerMockRemoveUntil struct {
	mock              *RecordCleanerMock
	mainExpectation   *RecordCleanerMockRemoveUntilExpectation
	expectationSeries []*RecordCleanerMockRemoveUntilExpectation
}

type RecordCleanerMockRemoveUntilExpectation struct {
	input *RecordCleanerMockRemoveUntilInput
}

type RecordCleanerMockRemoveUntilInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of RecordCleaner.RemoveUntil is expected from 1 to Infinity times
func (m *mRecordCleanerMockRemoveUntil) Expect(p context.Context, p1 insolar.PulseNumber) *mRecordCleanerMockRemoveUntil {
	m.mock.RemoveUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCleanerMockRemoveUntilExpectation{}
	}
	m.mainExpectation.input = &RecordCleanerMockRemoveUntilInput{p, p1}
	return m
}

//Return specifies results of invocation of RecordCleaner.RemoveUntil
func (m *mRecordCleanerMockRemoveUntil) Return() *RecordCleanerMock {
	m.mock.RemoveUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCleanerMockRemoveUntilExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RecordCleaner.RemoveUntil is expected once
func (m *mRecordCleanerMockRemoveUntil) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *RecordCleanerMockRemoveUntilExpectation {
	m.mock.RemoveUntilFunc = nil
	m.mainExpectation = nil

	expectation := &RecordCleanerMockRemoveUntilExpectation{}
	expectation.input = &RecordCleanerMockRemoveUntilInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RecordCleaner.RemoveUntil method
func (m *mRecordCleanerMockRemoveUntil) Set(f func(p context.Context, p1 insolar.PulseNumber)) *RecordCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveUntilFunc = f
	return m.mock
}

//RemoveUntil implements github.com/insolar/insolar/ledger/storage/object.RecordCleaner interface
func (m *RecordCleanerMock) RemoveUntil(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.RemoveUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveUntilCounter, 1)

	if len(m.RemoveUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordCleanerMock.RemoveUntil. %v %v", p, p1)
			return
		}

		input := m.RemoveUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordCleanerMockRemoveUntilInput{p, p1}, "RecordCleaner.RemoveUntil got unexpected parameters")

		return
	}

	if m.RemoveUntilMock.mainExpectation != nil {

		input := m.RemoveUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordCleanerMockRemoveUntilInput{p, p1}, "RecordCleaner.RemoveUntil got unexpected parameters")
		}

		return
	}

	if m.RemoveUntilFunc == nil {
		m.t.Fatalf("Unexpected call to RecordCleanerMock.RemoveUntil. %v %v", p, p1)
		return
	}

	m.RemoveUntilFunc(p, p1)
}

//RemoveUntilMinimockCounter returns a count of RecordCleanerMock.RemoveUntilFunc invocations
func (m *RecordCleanerMock) RemoveUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveUntilCounter)
}

//RemoveUntilMinimockPreCounter returns the value of RecordCleanerMock.RemoveUntil invocations
func (m *RecordCleanerMock) RemoveUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveUntilPreCounter)
}

//RemoveUntilFinished returns true if mock invocations count is ok
func (m *RecordCleanerMock) RemoveUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveUntilCounter) == uint64(len(m.RemoveUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveUntilCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordCleanerMock) ValidateCallCounters() {

	if !m.RemoveUntilFinished() {
		m.t.Fatal("Expected call to RecordCleanerMock.RemoveUntil")
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

	if !m.RemoveUntilFinished() {
		m.t.Fatal("Expected call to RecordCleanerMock.RemoveUntil")
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
		ok = ok && m.RemoveUntilFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveUntilFinished() {
				m.t.Error("Expected call to RecordCleanerMock.RemoveUntil")
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

	if !m.RemoveUntilFinished() {
		return false
	}

	return true
}

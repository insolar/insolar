package drop

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Cleaner" can be found in github.com/insolar/insolar/ledger/storage/drop
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//CleanerMock implements github.com/insolar/insolar/ledger/storage/drop.Cleaner
type CleanerMock struct {
	t minimock.Tester

	DeleteFunc       func(p core.PulseNumber)
	DeleteCounter    uint64
	DeletePreCounter uint64
	DeleteMock       mCleanerMockDelete
}

//NewCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage/drop.Cleaner
func NewCleanerMock(t minimock.Tester) *CleanerMock {
	m := &CleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteMock = mCleanerMockDelete{mock: m}

	return m
}

type mCleanerMockDelete struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockDeleteExpectation
	expectationSeries []*CleanerMockDeleteExpectation
}

type CleanerMockDeleteExpectation struct {
	input *CleanerMockDeleteInput
}

type CleanerMockDeleteInput struct {
	p core.PulseNumber
}

//Expect specifies that invocation of Cleaner.Delete is expected from 1 to Infinity times
func (m *mCleanerMockDelete) Expect(p core.PulseNumber) *mCleanerMockDelete {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockDeleteExpectation{}
	}
	m.mainExpectation.input = &CleanerMockDeleteInput{p}
	return m
}

//Return specifies results of invocation of Cleaner.Delete
func (m *mCleanerMockDelete) Return() *CleanerMock {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockDeleteExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Cleaner.Delete is expected once
func (m *mCleanerMockDelete) ExpectOnce(p core.PulseNumber) *CleanerMockDeleteExpectation {
	m.mock.DeleteFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockDeleteExpectation{}
	expectation.input = &CleanerMockDeleteInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Cleaner.Delete method
func (m *mCleanerMockDelete) Set(f func(p core.PulseNumber)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteFunc = f
	return m.mock
}

//Delete implements github.com/insolar/insolar/ledger/storage/drop.Cleaner interface
func (m *CleanerMock) Delete(p core.PulseNumber) {
	counter := atomic.AddUint64(&m.DeletePreCounter, 1)
	defer atomic.AddUint64(&m.DeleteCounter, 1)

	if len(m.DeleteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.Delete. %v", p)
			return
		}

		input := m.DeleteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockDeleteInput{p}, "Cleaner.Delete got unexpected parameters")

		return
	}

	if m.DeleteMock.mainExpectation != nil {

		input := m.DeleteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockDeleteInput{p}, "Cleaner.Delete got unexpected parameters")
		}

		return
	}

	if m.DeleteFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.Delete. %v", p)
		return
	}

	m.DeleteFunc(p)
}

//DeleteMinimockCounter returns a count of CleanerMock.DeleteFunc invocations
func (m *CleanerMock) DeleteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteCounter)
}

//DeleteMinimockPreCounter returns the value of CleanerMock.Delete invocations
func (m *CleanerMock) DeleteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeletePreCounter)
}

//DeleteFinished returns true if mock invocations count is ok
func (m *CleanerMock) DeleteFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeleteMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeleteCounter) == uint64(len(m.DeleteMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeleteMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeleteCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeleteFunc != nil {
		return atomic.LoadUint64(&m.DeleteCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CleanerMock) ValidateCallCounters() {

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to CleanerMock.Delete")
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

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to CleanerMock.Delete")
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
		ok = ok && m.DeleteFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeleteFinished() {
				m.t.Error("Expected call to CleanerMock.Delete")
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

	if !m.DeleteFinished() {
		return false
	}

	return true
}

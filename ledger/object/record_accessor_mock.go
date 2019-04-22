package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//RecordAccessorMock implements github.com/insolar/insolar/ledger/object.RecordAccessor
type RecordAccessorMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.ID) (r record.MaterialRecord, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mRecordAccessorMockForID
}

//NewRecordAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.RecordAccessor
func NewRecordAccessorMock(t minimock.Tester) *RecordAccessorMock {
	m := &RecordAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mRecordAccessorMockForID{mock: m}

	return m
}

type mRecordAccessorMockForID struct {
	mock              *RecordAccessorMock
	mainExpectation   *RecordAccessorMockForIDExpectation
	expectationSeries []*RecordAccessorMockForIDExpectation
}

type RecordAccessorMockForIDExpectation struct {
	input  *RecordAccessorMockForIDInput
	result *RecordAccessorMockForIDResult
}

type RecordAccessorMockForIDInput struct {
	p  context.Context
	p1 insolar.ID
}

type RecordAccessorMockForIDResult struct {
	r  record.MaterialRecord
	r1 error
}

//Expect specifies that invocation of RecordAccessor.ForID is expected from 1 to Infinity times
func (m *mRecordAccessorMockForID) Expect(p context.Context, p1 insolar.ID) *mRecordAccessorMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordAccessorMockForIDExpectation{}
	}
	m.mainExpectation.input = &RecordAccessorMockForIDInput{p, p1}
	return m
}

//Return specifies results of invocation of RecordAccessor.ForID
func (m *mRecordAccessorMockForID) Return(r record.MaterialRecord, r1 error) *RecordAccessorMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordAccessorMockForIDExpectation{}
	}
	m.mainExpectation.result = &RecordAccessorMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RecordAccessor.ForID is expected once
func (m *mRecordAccessorMockForID) ExpectOnce(p context.Context, p1 insolar.ID) *RecordAccessorMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &RecordAccessorMockForIDExpectation{}
	expectation.input = &RecordAccessorMockForIDInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordAccessorMockForIDExpectation) Return(r record.MaterialRecord, r1 error) {
	e.result = &RecordAccessorMockForIDResult{r, r1}
}

//Set uses given function f as a mock of RecordAccessor.ForID method
func (m *mRecordAccessorMockForID) Set(f func(p context.Context, p1 insolar.ID) (r record.MaterialRecord, r1 error)) *RecordAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.RecordAccessor interface
func (m *RecordAccessorMock) ForID(p context.Context, p1 insolar.ID) (r record.MaterialRecord, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordAccessorMock.ForID. %v %v", p, p1)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordAccessorMockForIDInput{p, p1}, "RecordAccessor.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordAccessorMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordAccessorMockForIDInput{p, p1}, "RecordAccessor.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordAccessorMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to RecordAccessorMock.ForID. %v %v", p, p1)
		return
	}

	return m.ForIDFunc(p, p1)
}

//ForIDMinimockCounter returns a count of RecordAccessorMock.ForIDFunc invocations
func (m *RecordAccessorMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of RecordAccessorMock.ForID invocations
func (m *RecordAccessorMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *RecordAccessorMock) ForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForIDCounter) == uint64(len(m.ForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForIDFunc != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordAccessorMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to RecordAccessorMock.ForID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecordAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecordAccessorMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to RecordAccessorMock.ForID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecordAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RecordAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to RecordAccessorMock.ForID")
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
func (m *RecordAccessorMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	return true
}

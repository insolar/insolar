package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordCollectionAccessor" can be found in github.com/insolar/insolar/ledger/object
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

//RecordCollectionAccessorMock implements github.com/insolar/insolar/ledger/object.RecordCollectionAccessor
type RecordCollectionAccessorMock struct {
	t minimock.Tester

	ForPulseFunc       func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []record.MaterialRecord)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mRecordCollectionAccessorMockForPulse
}

//NewRecordCollectionAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.RecordCollectionAccessor
func NewRecordCollectionAccessorMock(t minimock.Tester) *RecordCollectionAccessorMock {
	m := &RecordCollectionAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseMock = mRecordCollectionAccessorMockForPulse{mock: m}

	return m
}

type mRecordCollectionAccessorMockForPulse struct {
	mock              *RecordCollectionAccessorMock
	mainExpectation   *RecordCollectionAccessorMockForPulseExpectation
	expectationSeries []*RecordCollectionAccessorMockForPulseExpectation
}

type RecordCollectionAccessorMockForPulseExpectation struct {
	input  *RecordCollectionAccessorMockForPulseInput
	result *RecordCollectionAccessorMockForPulseResult
}

type RecordCollectionAccessorMockForPulseInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 insolar.PulseNumber
}

type RecordCollectionAccessorMockForPulseResult struct {
	r []record.MaterialRecord
}

//Expect specifies that invocation of RecordCollectionAccessor.ForPulse is expected from 1 to Infinity times
func (m *mRecordCollectionAccessorMockForPulse) Expect(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *mRecordCollectionAccessorMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCollectionAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.input = &RecordCollectionAccessorMockForPulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RecordCollectionAccessor.ForPulse
func (m *mRecordCollectionAccessorMockForPulse) Return(r []record.MaterialRecord) *RecordCollectionAccessorMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordCollectionAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.result = &RecordCollectionAccessorMockForPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecordCollectionAccessor.ForPulse is expected once
func (m *mRecordCollectionAccessorMockForPulse) ExpectOnce(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *RecordCollectionAccessorMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &RecordCollectionAccessorMockForPulseExpectation{}
	expectation.input = &RecordCollectionAccessorMockForPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordCollectionAccessorMockForPulseExpectation) Return(r []record.MaterialRecord) {
	e.result = &RecordCollectionAccessorMockForPulseResult{r}
}

//Set uses given function f as a mock of RecordCollectionAccessor.ForPulse method
func (m *mRecordCollectionAccessorMockForPulse) Set(f func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []record.MaterialRecord)) *RecordCollectionAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/object.RecordCollectionAccessor interface
func (m *RecordCollectionAccessorMock) ForPulse(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []record.MaterialRecord) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordCollectionAccessorMock.ForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordCollectionAccessorMockForPulseInput{p, p1, p2}, "RecordCollectionAccessor.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordCollectionAccessorMock.ForPulse")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordCollectionAccessorMockForPulseInput{p, p1, p2}, "RecordCollectionAccessor.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordCollectionAccessorMock.ForPulse")
		}

		r = result.r

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to RecordCollectionAccessorMock.ForPulse. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPulseFunc(p, p1, p2)
}

//ForPulseMinimockCounter returns a count of RecordCollectionAccessorMock.ForPulseFunc invocations
func (m *RecordCollectionAccessorMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

//ForPulseMinimockPreCounter returns the value of RecordCollectionAccessorMock.ForPulse invocations
func (m *RecordCollectionAccessorMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

//ForPulseFinished returns true if mock invocations count is ok
func (m *RecordCollectionAccessorMock) ForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseCounter) == uint64(len(m.ForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseFunc != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordCollectionAccessorMock) ValidateCallCounters() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to RecordCollectionAccessorMock.ForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordCollectionAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecordCollectionAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecordCollectionAccessorMock) MinimockFinish() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to RecordCollectionAccessorMock.ForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecordCollectionAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RecordCollectionAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseFinished() {
				m.t.Error("Expected call to RecordCollectionAccessorMock.ForPulse")
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
func (m *RecordCollectionAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseFinished() {
		return false
	}

	return true
}

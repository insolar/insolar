package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordPositionModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//RecordPositionModifierMock implements github.com/insolar/insolar/ledger/object.RecordPositionModifier
type RecordPositionModifierMock struct {
	t minimock.Tester

	IncrementPositionFunc       func(p insolar.ID) (r error)
	IncrementPositionCounter    uint64
	IncrementPositionPreCounter uint64
	IncrementPositionMock       mRecordPositionModifierMockIncrementPosition
}

//NewRecordPositionModifierMock returns a mock for github.com/insolar/insolar/ledger/object.RecordPositionModifier
func NewRecordPositionModifierMock(t minimock.Tester) *RecordPositionModifierMock {
	m := &RecordPositionModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IncrementPositionMock = mRecordPositionModifierMockIncrementPosition{mock: m}

	return m
}

type mRecordPositionModifierMockIncrementPosition struct {
	mock              *RecordPositionModifierMock
	mainExpectation   *RecordPositionModifierMockIncrementPositionExpectation
	expectationSeries []*RecordPositionModifierMockIncrementPositionExpectation
}

type RecordPositionModifierMockIncrementPositionExpectation struct {
	input  *RecordPositionModifierMockIncrementPositionInput
	result *RecordPositionModifierMockIncrementPositionResult
}

type RecordPositionModifierMockIncrementPositionInput struct {
	p insolar.ID
}

type RecordPositionModifierMockIncrementPositionResult struct {
	r error
}

//Expect specifies that invocation of RecordPositionModifier.IncrementPosition is expected from 1 to Infinity times
func (m *mRecordPositionModifierMockIncrementPosition) Expect(p insolar.ID) *mRecordPositionModifierMockIncrementPosition {
	m.mock.IncrementPositionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordPositionModifierMockIncrementPositionExpectation{}
	}
	m.mainExpectation.input = &RecordPositionModifierMockIncrementPositionInput{p}
	return m
}

//Return specifies results of invocation of RecordPositionModifier.IncrementPosition
func (m *mRecordPositionModifierMockIncrementPosition) Return(r error) *RecordPositionModifierMock {
	m.mock.IncrementPositionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordPositionModifierMockIncrementPositionExpectation{}
	}
	m.mainExpectation.result = &RecordPositionModifierMockIncrementPositionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecordPositionModifier.IncrementPosition is expected once
func (m *mRecordPositionModifierMockIncrementPosition) ExpectOnce(p insolar.ID) *RecordPositionModifierMockIncrementPositionExpectation {
	m.mock.IncrementPositionFunc = nil
	m.mainExpectation = nil

	expectation := &RecordPositionModifierMockIncrementPositionExpectation{}
	expectation.input = &RecordPositionModifierMockIncrementPositionInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordPositionModifierMockIncrementPositionExpectation) Return(r error) {
	e.result = &RecordPositionModifierMockIncrementPositionResult{r}
}

//Set uses given function f as a mock of RecordPositionModifier.IncrementPosition method
func (m *mRecordPositionModifierMockIncrementPosition) Set(f func(p insolar.ID) (r error)) *RecordPositionModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IncrementPositionFunc = f
	return m.mock
}

//IncrementPosition implements github.com/insolar/insolar/ledger/object.RecordPositionModifier interface
func (m *RecordPositionModifierMock) IncrementPosition(p insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.IncrementPositionPreCounter, 1)
	defer atomic.AddUint64(&m.IncrementPositionCounter, 1)

	if len(m.IncrementPositionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IncrementPositionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordPositionModifierMock.IncrementPosition. %v", p)
			return
		}

		input := m.IncrementPositionMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordPositionModifierMockIncrementPositionInput{p}, "RecordPositionModifier.IncrementPosition got unexpected parameters")

		result := m.IncrementPositionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordPositionModifierMock.IncrementPosition")
			return
		}

		r = result.r

		return
	}

	if m.IncrementPositionMock.mainExpectation != nil {

		input := m.IncrementPositionMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordPositionModifierMockIncrementPositionInput{p}, "RecordPositionModifier.IncrementPosition got unexpected parameters")
		}

		result := m.IncrementPositionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordPositionModifierMock.IncrementPosition")
		}

		r = result.r

		return
	}

	if m.IncrementPositionFunc == nil {
		m.t.Fatalf("Unexpected call to RecordPositionModifierMock.IncrementPosition. %v", p)
		return
	}

	return m.IncrementPositionFunc(p)
}

//IncrementPositionMinimockCounter returns a count of RecordPositionModifierMock.IncrementPositionFunc invocations
func (m *RecordPositionModifierMock) IncrementPositionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IncrementPositionCounter)
}

//IncrementPositionMinimockPreCounter returns the value of RecordPositionModifierMock.IncrementPosition invocations
func (m *RecordPositionModifierMock) IncrementPositionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IncrementPositionPreCounter)
}

//IncrementPositionFinished returns true if mock invocations count is ok
func (m *RecordPositionModifierMock) IncrementPositionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IncrementPositionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IncrementPositionCounter) == uint64(len(m.IncrementPositionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IncrementPositionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IncrementPositionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IncrementPositionFunc != nil {
		return atomic.LoadUint64(&m.IncrementPositionCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordPositionModifierMock) ValidateCallCounters() {

	if !m.IncrementPositionFinished() {
		m.t.Fatal("Expected call to RecordPositionModifierMock.IncrementPosition")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordPositionModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecordPositionModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecordPositionModifierMock) MinimockFinish() {

	if !m.IncrementPositionFinished() {
		m.t.Fatal("Expected call to RecordPositionModifierMock.IncrementPosition")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecordPositionModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RecordPositionModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.IncrementPositionFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IncrementPositionFinished() {
				m.t.Error("Expected call to RecordPositionModifierMock.IncrementPosition")
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
func (m *RecordPositionModifierMock) AllMocksCalled() bool {

	if !m.IncrementPositionFinished() {
		return false
	}

	return true
}

package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordPositionAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// RecordPositionAccessorMock implements github.com/insolar/insolar/ledger/object.RecordPositionAccessor
type RecordPositionAccessorMock struct {
	t minimock.Tester

	AtPositionFunc       func(p insolar.PulseNumber, p1 uint32) (r insolar.ID, r1 error)
	AtPositionCounter    uint64
	AtPositionPreCounter uint64
	AtPositionMock       mRecordPositionAccessorMockAtPosition

	LastKnownPositionFunc       func(p insolar.PulseNumber) (r uint32, r1 error)
	LastKnownPositionCounter    uint64
	LastKnownPositionPreCounter uint64
	LastKnownPositionMock       mRecordPositionAccessorMockLastKnownPosition
}

// NewRecordPositionAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.RecordPositionAccessor
func NewRecordPositionAccessorMock(t minimock.Tester) *RecordPositionAccessorMock {
	m := &RecordPositionAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AtPositionMock = mRecordPositionAccessorMockAtPosition{mock: m}
	m.LastKnownPositionMock = mRecordPositionAccessorMockLastKnownPosition{mock: m}

	return m
}

type mRecordPositionAccessorMockAtPosition struct {
	mock              *RecordPositionAccessorMock
	mainExpectation   *RecordPositionAccessorMockAtPositionExpectation
	expectationSeries []*RecordPositionAccessorMockAtPositionExpectation
}

type RecordPositionAccessorMockAtPositionExpectation struct {
	input  *RecordPositionAccessorMockAtPositionInput
	result *RecordPositionAccessorMockAtPositionResult
}

type RecordPositionAccessorMockAtPositionInput struct {
	p  insolar.PulseNumber
	p1 uint32
}

type RecordPositionAccessorMockAtPositionResult struct {
	r  insolar.ID
	r1 error
}

// Expect specifies that invocation of RecordPositionAccessor.AtPosition is expected from 1 to Infinity times
func (m *mRecordPositionAccessorMockAtPosition) Expect(p insolar.PulseNumber, p1 uint32) *mRecordPositionAccessorMockAtPosition {
	m.mock.AtPositionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordPositionAccessorMockAtPositionExpectation{}
	}
	m.mainExpectation.input = &RecordPositionAccessorMockAtPositionInput{p, p1}
	return m
}

// Return specifies results of invocation of RecordPositionAccessor.AtPosition
func (m *mRecordPositionAccessorMockAtPosition) Return(r insolar.ID, r1 error) *RecordPositionAccessorMock {
	m.mock.AtPositionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordPositionAccessorMockAtPositionExpectation{}
	}
	m.mainExpectation.result = &RecordPositionAccessorMockAtPositionResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of RecordPositionAccessor.AtPosition is expected once
func (m *mRecordPositionAccessorMockAtPosition) ExpectOnce(p insolar.PulseNumber, p1 uint32) *RecordPositionAccessorMockAtPositionExpectation {
	m.mock.AtPositionFunc = nil
	m.mainExpectation = nil

	expectation := &RecordPositionAccessorMockAtPositionExpectation{}
	expectation.input = &RecordPositionAccessorMockAtPositionInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordPositionAccessorMockAtPositionExpectation) Return(r insolar.ID, r1 error) {
	e.result = &RecordPositionAccessorMockAtPositionResult{r, r1}
}

// Set uses given function f as a mock of RecordPositionAccessor.AtPosition method
func (m *mRecordPositionAccessorMockAtPosition) Set(f func(p insolar.PulseNumber, p1 uint32) (r insolar.ID, r1 error)) *RecordPositionAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AtPositionFunc = f
	return m.mock
}

// AtPosition implements github.com/insolar/insolar/ledger/object.RecordPositionAccessor interface
func (m *RecordPositionAccessorMock) AtPosition(p insolar.PulseNumber, p1 uint32) (r insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.AtPositionPreCounter, 1)
	defer atomic.AddUint64(&m.AtPositionCounter, 1)

	if len(m.AtPositionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AtPositionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordPositionAccessorMock.AtPosition. %v %v", p, p1)
			return
		}

		input := m.AtPositionMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordPositionAccessorMockAtPositionInput{p, p1}, "RecordPositionAccessor.AtPosition got unexpected parameters")

		result := m.AtPositionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordPositionAccessorMock.AtPosition")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AtPositionMock.mainExpectation != nil {

		input := m.AtPositionMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordPositionAccessorMockAtPositionInput{p, p1}, "RecordPositionAccessor.AtPosition got unexpected parameters")
		}

		result := m.AtPositionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordPositionAccessorMock.AtPosition")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AtPositionFunc == nil {
		m.t.Fatalf("Unexpected call to RecordPositionAccessorMock.AtPosition. %v %v", p, p1)
		return
	}

	return m.AtPositionFunc(p, p1)
}

// AtPositionMinimockCounter returns a count of RecordPositionAccessorMock.AtPositionFunc invocations
func (m *RecordPositionAccessorMock) AtPositionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AtPositionCounter)
}

// AtPositionMinimockPreCounter returns the value of RecordPositionAccessorMock.AtPosition invocations
func (m *RecordPositionAccessorMock) AtPositionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AtPositionPreCounter)
}

// AtPositionFinished returns true if mock invocations count is ok
func (m *RecordPositionAccessorMock) AtPositionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AtPositionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AtPositionCounter) == uint64(len(m.AtPositionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AtPositionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AtPositionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AtPositionFunc != nil {
		return atomic.LoadUint64(&m.AtPositionCounter) > 0
	}

	return true
}

type mRecordPositionAccessorMockLastKnownPosition struct {
	mock              *RecordPositionAccessorMock
	mainExpectation   *RecordPositionAccessorMockLastKnownPositionExpectation
	expectationSeries []*RecordPositionAccessorMockLastKnownPositionExpectation
}

type RecordPositionAccessorMockLastKnownPositionExpectation struct {
	input  *RecordPositionAccessorMockLastKnownPositionInput
	result *RecordPositionAccessorMockLastKnownPositionResult
}

type RecordPositionAccessorMockLastKnownPositionInput struct {
	p insolar.PulseNumber
}

type RecordPositionAccessorMockLastKnownPositionResult struct {
	r  uint32
	r1 error
}

// Expect specifies that invocation of RecordPositionAccessor.LastKnownPosition is expected from 1 to Infinity times
func (m *mRecordPositionAccessorMockLastKnownPosition) Expect(p insolar.PulseNumber) *mRecordPositionAccessorMockLastKnownPosition {
	m.mock.LastKnownPositionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordPositionAccessorMockLastKnownPositionExpectation{}
	}
	m.mainExpectation.input = &RecordPositionAccessorMockLastKnownPositionInput{p}
	return m
}

// Return specifies results of invocation of RecordPositionAccessor.LastKnownPosition
func (m *mRecordPositionAccessorMockLastKnownPosition) Return(r uint32, r1 error) *RecordPositionAccessorMock {
	m.mock.LastKnownPositionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordPositionAccessorMockLastKnownPositionExpectation{}
	}
	m.mainExpectation.result = &RecordPositionAccessorMockLastKnownPositionResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of RecordPositionAccessor.LastKnownPosition is expected once
func (m *mRecordPositionAccessorMockLastKnownPosition) ExpectOnce(p insolar.PulseNumber) *RecordPositionAccessorMockLastKnownPositionExpectation {
	m.mock.LastKnownPositionFunc = nil
	m.mainExpectation = nil

	expectation := &RecordPositionAccessorMockLastKnownPositionExpectation{}
	expectation.input = &RecordPositionAccessorMockLastKnownPositionInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordPositionAccessorMockLastKnownPositionExpectation) Return(r uint32, r1 error) {
	e.result = &RecordPositionAccessorMockLastKnownPositionResult{r, r1}
}

// Set uses given function f as a mock of RecordPositionAccessor.LastKnownPosition method
func (m *mRecordPositionAccessorMockLastKnownPosition) Set(f func(p insolar.PulseNumber) (r uint32, r1 error)) *RecordPositionAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LastKnownPositionFunc = f
	return m.mock
}

// LastKnownPosition implements github.com/insolar/insolar/ledger/object.RecordPositionAccessor interface
func (m *RecordPositionAccessorMock) LastKnownPosition(p insolar.PulseNumber) (r uint32, r1 error) {
	counter := atomic.AddUint64(&m.LastKnownPositionPreCounter, 1)
	defer atomic.AddUint64(&m.LastKnownPositionCounter, 1)

	if len(m.LastKnownPositionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LastKnownPositionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordPositionAccessorMock.LastKnownPosition. %v", p)
			return
		}

		input := m.LastKnownPositionMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordPositionAccessorMockLastKnownPositionInput{p}, "RecordPositionAccessor.LastKnownPosition got unexpected parameters")

		result := m.LastKnownPositionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordPositionAccessorMock.LastKnownPosition")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LastKnownPositionMock.mainExpectation != nil {

		input := m.LastKnownPositionMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordPositionAccessorMockLastKnownPositionInput{p}, "RecordPositionAccessor.LastKnownPosition got unexpected parameters")
		}

		result := m.LastKnownPositionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordPositionAccessorMock.LastKnownPosition")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LastKnownPositionFunc == nil {
		m.t.Fatalf("Unexpected call to RecordPositionAccessorMock.LastKnownPosition. %v", p)
		return
	}

	return m.LastKnownPositionFunc(p)
}

// LastKnownPositionMinimockCounter returns a count of RecordPositionAccessorMock.LastKnownPositionFunc invocations
func (m *RecordPositionAccessorMock) LastKnownPositionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LastKnownPositionCounter)
}

// LastKnownPositionMinimockPreCounter returns the value of RecordPositionAccessorMock.LastKnownPosition invocations
func (m *RecordPositionAccessorMock) LastKnownPositionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LastKnownPositionPreCounter)
}

// LastKnownPositionFinished returns true if mock invocations count is ok
func (m *RecordPositionAccessorMock) LastKnownPositionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LastKnownPositionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LastKnownPositionCounter) == uint64(len(m.LastKnownPositionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LastKnownPositionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LastKnownPositionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LastKnownPositionFunc != nil {
		return atomic.LoadUint64(&m.LastKnownPositionCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordPositionAccessorMock) ValidateCallCounters() {

	if !m.AtPositionFinished() {
		m.t.Fatal("Expected call to RecordPositionAccessorMock.AtPosition")
	}

	if !m.LastKnownPositionFinished() {
		m.t.Fatal("Expected call to RecordPositionAccessorMock.LastKnownPosition")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordPositionAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecordPositionAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecordPositionAccessorMock) MinimockFinish() {

	if !m.AtPositionFinished() {
		m.t.Fatal("Expected call to RecordPositionAccessorMock.AtPosition")
	}

	if !m.LastKnownPositionFinished() {
		m.t.Fatal("Expected call to RecordPositionAccessorMock.LastKnownPosition")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecordPositionAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *RecordPositionAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AtPositionFinished()
		ok = ok && m.LastKnownPositionFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AtPositionFinished() {
				m.t.Error("Expected call to RecordPositionAccessorMock.AtPosition")
			}

			if !m.LastKnownPositionFinished() {
				m.t.Error("Expected call to RecordPositionAccessorMock.LastKnownPosition")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *RecordPositionAccessorMock) AllMocksCalled() bool {

	if !m.AtPositionFinished() {
		return false
	}

	if !m.LastKnownPositionFinished() {
		return false
	}

	return true
}

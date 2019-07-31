package object

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i github.com/insolar/insolar/ledger/object.RecordPositionAccessor -o ./record_position_accessor_mock.go

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
)

// RecordPositionAccessorMock implements RecordPositionAccessor
type RecordPositionAccessorMock struct {
	t minimock.Tester

	funcAtPosition          func(pn insolar.PulseNumber, position uint32) (i1 insolar.ID, err error)
	inspectFuncAtPosition   func(pn insolar.PulseNumber, position uint32)
	afterAtPositionCounter  uint64
	beforeAtPositionCounter uint64
	AtPositionMock          mRecordPositionAccessorMockAtPosition

	funcLastKnownPosition          func(pn insolar.PulseNumber) (u1 uint32, err error)
	inspectFuncLastKnownPosition   func(pn insolar.PulseNumber)
	afterLastKnownPositionCounter  uint64
	beforeLastKnownPositionCounter uint64
	LastKnownPositionMock          mRecordPositionAccessorMockLastKnownPosition
}

// NewRecordPositionAccessorMock returns a mock for RecordPositionAccessor
func NewRecordPositionAccessorMock(t minimock.Tester) *RecordPositionAccessorMock {
	m := &RecordPositionAccessorMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AtPositionMock = mRecordPositionAccessorMockAtPosition{mock: m}
	m.AtPositionMock.callArgs = []*RecordPositionAccessorMockAtPositionParams{}

	m.LastKnownPositionMock = mRecordPositionAccessorMockLastKnownPosition{mock: m}
	m.LastKnownPositionMock.callArgs = []*RecordPositionAccessorMockLastKnownPositionParams{}

	return m
}

type mRecordPositionAccessorMockAtPosition struct {
	mock               *RecordPositionAccessorMock
	defaultExpectation *RecordPositionAccessorMockAtPositionExpectation
	expectations       []*RecordPositionAccessorMockAtPositionExpectation

	callArgs []*RecordPositionAccessorMockAtPositionParams
	mutex    sync.RWMutex
}

// RecordPositionAccessorMockAtPositionExpectation specifies expectation struct of the RecordPositionAccessor.AtPosition
type RecordPositionAccessorMockAtPositionExpectation struct {
	mock    *RecordPositionAccessorMock
	params  *RecordPositionAccessorMockAtPositionParams
	results *RecordPositionAccessorMockAtPositionResults
	Counter uint64
}

// RecordPositionAccessorMockAtPositionParams contains parameters of the RecordPositionAccessor.AtPosition
type RecordPositionAccessorMockAtPositionParams struct {
	pn       insolar.PulseNumber
	position uint32
}

// RecordPositionAccessorMockAtPositionResults contains results of the RecordPositionAccessor.AtPosition
type RecordPositionAccessorMockAtPositionResults struct {
	i1  insolar.ID
	err error
}

// Expect sets up expected params for RecordPositionAccessor.AtPosition
func (mmAtPosition *mRecordPositionAccessorMockAtPosition) Expect(pn insolar.PulseNumber, position uint32) *mRecordPositionAccessorMockAtPosition {
	if mmAtPosition.mock.funcAtPosition != nil {
		mmAtPosition.mock.t.Fatalf("RecordPositionAccessorMock.AtPosition mock is already set by Set")
	}

	if mmAtPosition.defaultExpectation == nil {
		mmAtPosition.defaultExpectation = &RecordPositionAccessorMockAtPositionExpectation{}
	}

	mmAtPosition.defaultExpectation.params = &RecordPositionAccessorMockAtPositionParams{pn, position}
	for _, e := range mmAtPosition.expectations {
		if minimock.Equal(e.params, mmAtPosition.defaultExpectation.params) {
			mmAtPosition.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmAtPosition.defaultExpectation.params)
		}
	}

	return mmAtPosition
}

// Inspect accepts an inspector function that has same arguments as the RecordPositionAccessor.AtPosition
func (mmAtPosition *mRecordPositionAccessorMockAtPosition) Inspect(f func(pn insolar.PulseNumber, position uint32)) *mRecordPositionAccessorMockAtPosition {
	if mmAtPosition.mock.inspectFuncAtPosition != nil {
		mmAtPosition.mock.t.Fatalf("Inspect function is already set for RecordPositionAccessorMock.AtPosition")
	}

	mmAtPosition.mock.inspectFuncAtPosition = f

	return mmAtPosition
}

// Return sets up results that will be returned by RecordPositionAccessor.AtPosition
func (mmAtPosition *mRecordPositionAccessorMockAtPosition) Return(i1 insolar.ID, err error) *RecordPositionAccessorMock {
	if mmAtPosition.mock.funcAtPosition != nil {
		mmAtPosition.mock.t.Fatalf("RecordPositionAccessorMock.AtPosition mock is already set by Set")
	}

	if mmAtPosition.defaultExpectation == nil {
		mmAtPosition.defaultExpectation = &RecordPositionAccessorMockAtPositionExpectation{mock: mmAtPosition.mock}
	}
	mmAtPosition.defaultExpectation.results = &RecordPositionAccessorMockAtPositionResults{i1, err}
	return mmAtPosition.mock
}

//Set uses given function f to mock the RecordPositionAccessor.AtPosition method
func (mmAtPosition *mRecordPositionAccessorMockAtPosition) Set(f func(pn insolar.PulseNumber, position uint32) (i1 insolar.ID, err error)) *RecordPositionAccessorMock {
	if mmAtPosition.defaultExpectation != nil {
		mmAtPosition.mock.t.Fatalf("Default expectation is already set for the RecordPositionAccessor.AtPosition method")
	}

	if len(mmAtPosition.expectations) > 0 {
		mmAtPosition.mock.t.Fatalf("Some expectations are already set for the RecordPositionAccessor.AtPosition method")
	}

	mmAtPosition.mock.funcAtPosition = f
	return mmAtPosition.mock
}

// When sets expectation for the RecordPositionAccessor.AtPosition which will trigger the result defined by the following
// Then helper
func (mmAtPosition *mRecordPositionAccessorMockAtPosition) When(pn insolar.PulseNumber, position uint32) *RecordPositionAccessorMockAtPositionExpectation {
	if mmAtPosition.mock.funcAtPosition != nil {
		mmAtPosition.mock.t.Fatalf("RecordPositionAccessorMock.AtPosition mock is already set by Set")
	}

	expectation := &RecordPositionAccessorMockAtPositionExpectation{
		mock:   mmAtPosition.mock,
		params: &RecordPositionAccessorMockAtPositionParams{pn, position},
	}
	mmAtPosition.expectations = append(mmAtPosition.expectations, expectation)
	return expectation
}

// Then sets up RecordPositionAccessor.AtPosition return parameters for the expectation previously defined by the When method
func (e *RecordPositionAccessorMockAtPositionExpectation) Then(i1 insolar.ID, err error) *RecordPositionAccessorMock {
	e.results = &RecordPositionAccessorMockAtPositionResults{i1, err}
	return e.mock
}

// AtPosition implements RecordPositionAccessor
func (mmAtPosition *RecordPositionAccessorMock) AtPosition(pn insolar.PulseNumber, position uint32) (i1 insolar.ID, err error) {
	mm_atomic.AddUint64(&mmAtPosition.beforeAtPositionCounter, 1)
	defer mm_atomic.AddUint64(&mmAtPosition.afterAtPositionCounter, 1)

	if mmAtPosition.inspectFuncAtPosition != nil {
		mmAtPosition.inspectFuncAtPosition(pn, position)
	}

	params := &RecordPositionAccessorMockAtPositionParams{pn, position}

	// Record call args
	mmAtPosition.AtPositionMock.mutex.Lock()
	mmAtPosition.AtPositionMock.callArgs = append(mmAtPosition.AtPositionMock.callArgs, params)
	mmAtPosition.AtPositionMock.mutex.Unlock()

	for _, e := range mmAtPosition.AtPositionMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.i1, e.results.err
		}
	}

	if mmAtPosition.AtPositionMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmAtPosition.AtPositionMock.defaultExpectation.Counter, 1)
		want := mmAtPosition.AtPositionMock.defaultExpectation.params
		got := RecordPositionAccessorMockAtPositionParams{pn, position}
		if want != nil && !minimock.Equal(*want, got) {
			mmAtPosition.t.Errorf("RecordPositionAccessorMock.AtPosition got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmAtPosition.AtPositionMock.defaultExpectation.results
		if results == nil {
			mmAtPosition.t.Fatal("No results are set for the RecordPositionAccessorMock.AtPosition")
		}
		return (*results).i1, (*results).err
	}
	if mmAtPosition.funcAtPosition != nil {
		return mmAtPosition.funcAtPosition(pn, position)
	}
	mmAtPosition.t.Fatalf("Unexpected call to RecordPositionAccessorMock.AtPosition. %v %v", pn, position)
	return
}

// AtPositionAfterCounter returns a count of finished RecordPositionAccessorMock.AtPosition invocations
func (mmAtPosition *RecordPositionAccessorMock) AtPositionAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAtPosition.afterAtPositionCounter)
}

// AtPositionBeforeCounter returns a count of RecordPositionAccessorMock.AtPosition invocations
func (mmAtPosition *RecordPositionAccessorMock) AtPositionBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAtPosition.beforeAtPositionCounter)
}

// Calls returns a list of arguments used in each call to RecordPositionAccessorMock.AtPosition.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmAtPosition *mRecordPositionAccessorMockAtPosition) Calls() []*RecordPositionAccessorMockAtPositionParams {
	mmAtPosition.mutex.RLock()

	argCopy := make([]*RecordPositionAccessorMockAtPositionParams, len(mmAtPosition.callArgs))
	copy(argCopy, mmAtPosition.callArgs)

	mmAtPosition.mutex.RUnlock()

	return argCopy
}

// MinimockAtPositionDone returns true if the count of the AtPosition invocations corresponds
// the number of defined expectations
func (m *RecordPositionAccessorMock) MinimockAtPositionDone() bool {
	for _, e := range m.AtPositionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AtPositionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAtPositionCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAtPosition != nil && mm_atomic.LoadUint64(&m.afterAtPositionCounter) < 1 {
		return false
	}
	return true
}

// MinimockAtPositionInspect logs each unmet expectation
func (m *RecordPositionAccessorMock) MinimockAtPositionInspect() {
	for _, e := range m.AtPositionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RecordPositionAccessorMock.AtPosition with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AtPositionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAtPositionCounter) < 1 {
		if m.AtPositionMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RecordPositionAccessorMock.AtPosition")
		} else {
			m.t.Errorf("Expected call to RecordPositionAccessorMock.AtPosition with params: %#v", *m.AtPositionMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAtPosition != nil && mm_atomic.LoadUint64(&m.afterAtPositionCounter) < 1 {
		m.t.Error("Expected call to RecordPositionAccessorMock.AtPosition")
	}
}

type mRecordPositionAccessorMockLastKnownPosition struct {
	mock               *RecordPositionAccessorMock
	defaultExpectation *RecordPositionAccessorMockLastKnownPositionExpectation
	expectations       []*RecordPositionAccessorMockLastKnownPositionExpectation

	callArgs []*RecordPositionAccessorMockLastKnownPositionParams
	mutex    sync.RWMutex
}

// RecordPositionAccessorMockLastKnownPositionExpectation specifies expectation struct of the RecordPositionAccessor.LastKnownPosition
type RecordPositionAccessorMockLastKnownPositionExpectation struct {
	mock    *RecordPositionAccessorMock
	params  *RecordPositionAccessorMockLastKnownPositionParams
	results *RecordPositionAccessorMockLastKnownPositionResults
	Counter uint64
}

// RecordPositionAccessorMockLastKnownPositionParams contains parameters of the RecordPositionAccessor.LastKnownPosition
type RecordPositionAccessorMockLastKnownPositionParams struct {
	pn insolar.PulseNumber
}

// RecordPositionAccessorMockLastKnownPositionResults contains results of the RecordPositionAccessor.LastKnownPosition
type RecordPositionAccessorMockLastKnownPositionResults struct {
	u1  uint32
	err error
}

// Expect sets up expected params for RecordPositionAccessor.LastKnownPosition
func (mmLastKnownPosition *mRecordPositionAccessorMockLastKnownPosition) Expect(pn insolar.PulseNumber) *mRecordPositionAccessorMockLastKnownPosition {
	if mmLastKnownPosition.mock.funcLastKnownPosition != nil {
		mmLastKnownPosition.mock.t.Fatalf("RecordPositionAccessorMock.LastKnownPosition mock is already set by Set")
	}

	if mmLastKnownPosition.defaultExpectation == nil {
		mmLastKnownPosition.defaultExpectation = &RecordPositionAccessorMockLastKnownPositionExpectation{}
	}

	mmLastKnownPosition.defaultExpectation.params = &RecordPositionAccessorMockLastKnownPositionParams{pn}
	for _, e := range mmLastKnownPosition.expectations {
		if minimock.Equal(e.params, mmLastKnownPosition.defaultExpectation.params) {
			mmLastKnownPosition.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmLastKnownPosition.defaultExpectation.params)
		}
	}

	return mmLastKnownPosition
}

// Inspect accepts an inspector function that has same arguments as the RecordPositionAccessor.LastKnownPosition
func (mmLastKnownPosition *mRecordPositionAccessorMockLastKnownPosition) Inspect(f func(pn insolar.PulseNumber)) *mRecordPositionAccessorMockLastKnownPosition {
	if mmLastKnownPosition.mock.inspectFuncLastKnownPosition != nil {
		mmLastKnownPosition.mock.t.Fatalf("Inspect function is already set for RecordPositionAccessorMock.LastKnownPosition")
	}

	mmLastKnownPosition.mock.inspectFuncLastKnownPosition = f

	return mmLastKnownPosition
}

// Return sets up results that will be returned by RecordPositionAccessor.LastKnownPosition
func (mmLastKnownPosition *mRecordPositionAccessorMockLastKnownPosition) Return(u1 uint32, err error) *RecordPositionAccessorMock {
	if mmLastKnownPosition.mock.funcLastKnownPosition != nil {
		mmLastKnownPosition.mock.t.Fatalf("RecordPositionAccessorMock.LastKnownPosition mock is already set by Set")
	}

	if mmLastKnownPosition.defaultExpectation == nil {
		mmLastKnownPosition.defaultExpectation = &RecordPositionAccessorMockLastKnownPositionExpectation{mock: mmLastKnownPosition.mock}
	}
	mmLastKnownPosition.defaultExpectation.results = &RecordPositionAccessorMockLastKnownPositionResults{u1, err}
	return mmLastKnownPosition.mock
}

//Set uses given function f to mock the RecordPositionAccessor.LastKnownPosition method
func (mmLastKnownPosition *mRecordPositionAccessorMockLastKnownPosition) Set(f func(pn insolar.PulseNumber) (u1 uint32, err error)) *RecordPositionAccessorMock {
	if mmLastKnownPosition.defaultExpectation != nil {
		mmLastKnownPosition.mock.t.Fatalf("Default expectation is already set for the RecordPositionAccessor.LastKnownPosition method")
	}

	if len(mmLastKnownPosition.expectations) > 0 {
		mmLastKnownPosition.mock.t.Fatalf("Some expectations are already set for the RecordPositionAccessor.LastKnownPosition method")
	}

	mmLastKnownPosition.mock.funcLastKnownPosition = f
	return mmLastKnownPosition.mock
}

// When sets expectation for the RecordPositionAccessor.LastKnownPosition which will trigger the result defined by the following
// Then helper
func (mmLastKnownPosition *mRecordPositionAccessorMockLastKnownPosition) When(pn insolar.PulseNumber) *RecordPositionAccessorMockLastKnownPositionExpectation {
	if mmLastKnownPosition.mock.funcLastKnownPosition != nil {
		mmLastKnownPosition.mock.t.Fatalf("RecordPositionAccessorMock.LastKnownPosition mock is already set by Set")
	}

	expectation := &RecordPositionAccessorMockLastKnownPositionExpectation{
		mock:   mmLastKnownPosition.mock,
		params: &RecordPositionAccessorMockLastKnownPositionParams{pn},
	}
	mmLastKnownPosition.expectations = append(mmLastKnownPosition.expectations, expectation)
	return expectation
}

// Then sets up RecordPositionAccessor.LastKnownPosition return parameters for the expectation previously defined by the When method
func (e *RecordPositionAccessorMockLastKnownPositionExpectation) Then(u1 uint32, err error) *RecordPositionAccessorMock {
	e.results = &RecordPositionAccessorMockLastKnownPositionResults{u1, err}
	return e.mock
}

// LastKnownPosition implements RecordPositionAccessor
func (mmLastKnownPosition *RecordPositionAccessorMock) LastKnownPosition(pn insolar.PulseNumber) (u1 uint32, err error) {
	mm_atomic.AddUint64(&mmLastKnownPosition.beforeLastKnownPositionCounter, 1)
	defer mm_atomic.AddUint64(&mmLastKnownPosition.afterLastKnownPositionCounter, 1)

	if mmLastKnownPosition.inspectFuncLastKnownPosition != nil {
		mmLastKnownPosition.inspectFuncLastKnownPosition(pn)
	}

	params := &RecordPositionAccessorMockLastKnownPositionParams{pn}

	// Record call args
	mmLastKnownPosition.LastKnownPositionMock.mutex.Lock()
	mmLastKnownPosition.LastKnownPositionMock.callArgs = append(mmLastKnownPosition.LastKnownPositionMock.callArgs, params)
	mmLastKnownPosition.LastKnownPositionMock.mutex.Unlock()

	for _, e := range mmLastKnownPosition.LastKnownPositionMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.u1, e.results.err
		}
	}

	if mmLastKnownPosition.LastKnownPositionMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmLastKnownPosition.LastKnownPositionMock.defaultExpectation.Counter, 1)
		want := mmLastKnownPosition.LastKnownPositionMock.defaultExpectation.params
		got := RecordPositionAccessorMockLastKnownPositionParams{pn}
		if want != nil && !minimock.Equal(*want, got) {
			mmLastKnownPosition.t.Errorf("RecordPositionAccessorMock.LastKnownPosition got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmLastKnownPosition.LastKnownPositionMock.defaultExpectation.results
		if results == nil {
			mmLastKnownPosition.t.Fatal("No results are set for the RecordPositionAccessorMock.LastKnownPosition")
		}
		return (*results).u1, (*results).err
	}
	if mmLastKnownPosition.funcLastKnownPosition != nil {
		return mmLastKnownPosition.funcLastKnownPosition(pn)
	}
	mmLastKnownPosition.t.Fatalf("Unexpected call to RecordPositionAccessorMock.LastKnownPosition. %v", pn)
	return
}

// LastKnownPositionAfterCounter returns a count of finished RecordPositionAccessorMock.LastKnownPosition invocations
func (mmLastKnownPosition *RecordPositionAccessorMock) LastKnownPositionAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmLastKnownPosition.afterLastKnownPositionCounter)
}

// LastKnownPositionBeforeCounter returns a count of RecordPositionAccessorMock.LastKnownPosition invocations
func (mmLastKnownPosition *RecordPositionAccessorMock) LastKnownPositionBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmLastKnownPosition.beforeLastKnownPositionCounter)
}

// Calls returns a list of arguments used in each call to RecordPositionAccessorMock.LastKnownPosition.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmLastKnownPosition *mRecordPositionAccessorMockLastKnownPosition) Calls() []*RecordPositionAccessorMockLastKnownPositionParams {
	mmLastKnownPosition.mutex.RLock()

	argCopy := make([]*RecordPositionAccessorMockLastKnownPositionParams, len(mmLastKnownPosition.callArgs))
	copy(argCopy, mmLastKnownPosition.callArgs)

	mmLastKnownPosition.mutex.RUnlock()

	return argCopy
}

// MinimockLastKnownPositionDone returns true if the count of the LastKnownPosition invocations corresponds
// the number of defined expectations
func (m *RecordPositionAccessorMock) MinimockLastKnownPositionDone() bool {
	for _, e := range m.LastKnownPositionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LastKnownPositionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterLastKnownPositionCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLastKnownPosition != nil && mm_atomic.LoadUint64(&m.afterLastKnownPositionCounter) < 1 {
		return false
	}
	return true
}

// MinimockLastKnownPositionInspect logs each unmet expectation
func (m *RecordPositionAccessorMock) MinimockLastKnownPositionInspect() {
	for _, e := range m.LastKnownPositionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RecordPositionAccessorMock.LastKnownPosition with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LastKnownPositionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterLastKnownPositionCounter) < 1 {
		if m.LastKnownPositionMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RecordPositionAccessorMock.LastKnownPosition")
		} else {
			m.t.Errorf("Expected call to RecordPositionAccessorMock.LastKnownPosition with params: %#v", *m.LastKnownPositionMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLastKnownPosition != nil && mm_atomic.LoadUint64(&m.afterLastKnownPositionCounter) < 1 {
		m.t.Error("Expected call to RecordPositionAccessorMock.LastKnownPosition")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *RecordPositionAccessorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockAtPositionInspect()

		m.MinimockLastKnownPositionInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *RecordPositionAccessorMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *RecordPositionAccessorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAtPositionDone() &&
		m.MinimockLastKnownPositionDone()
}

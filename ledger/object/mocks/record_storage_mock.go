package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordStorage" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

// RecordStorageMock implements github.com/insolar/insolar/ledger/object.RecordStorage
type RecordStorageMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.ID) (r record.Material, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mRecordStorageMockForID

	SetFunc       func(p context.Context, p1 insolar.ID, p2 record.Material) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mRecordStorageMockSet
}

// NewRecordStorageMock returns a mock for github.com/insolar/insolar/ledger/object.RecordStorage
func NewRecordStorageMock(t minimock.Tester) *RecordStorageMock {
	m := &RecordStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mRecordStorageMockForID{mock: m}
	m.SetMock = mRecordStorageMockSet{mock: m}

	return m
}

type mRecordStorageMockForID struct {
	mock              *RecordStorageMock
	mainExpectation   *RecordStorageMockForIDExpectation
	expectationSeries []*RecordStorageMockForIDExpectation
}

type RecordStorageMockForIDExpectation struct {
	input  *RecordStorageMockForIDInput
	result *RecordStorageMockForIDResult
}

type RecordStorageMockForIDInput struct {
	p  context.Context
	p1 insolar.ID
}

type RecordStorageMockForIDResult struct {
	r  record.Material
	r1 error
}

// Expect specifies that invocation of RecordStorage.ForID is expected from 1 to Infinity times
func (m *mRecordStorageMockForID) Expect(p context.Context, p1 insolar.ID) *mRecordStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordStorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &RecordStorageMockForIDInput{p, p1}
	return m
}

// Return specifies results of invocation of RecordStorage.ForID
func (m *mRecordStorageMockForID) Return(r record.Material, r1 error) *RecordStorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordStorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &RecordStorageMockForIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of RecordStorage.ForID is expected once
func (m *mRecordStorageMockForID) ExpectOnce(p context.Context, p1 insolar.ID) *RecordStorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &RecordStorageMockForIDExpectation{}
	expectation.input = &RecordStorageMockForIDInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordStorageMockForIDExpectation) Return(r record.Material, r1 error) {
	e.result = &RecordStorageMockForIDResult{r, r1}
}

// Set uses given function f as a mock of RecordStorage.ForID method
func (m *mRecordStorageMockForID) Set(f func(p context.Context, p1 insolar.ID) (r record.Material, r1 error)) *RecordStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

// ForID implements github.com/insolar/insolar/ledger/object.RecordStorage interface
func (m *RecordStorageMock) ForID(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordStorageMock.ForID. %v %v", p, p1)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordStorageMockForIDInput{p, p1}, "RecordStorage.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordStorageMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordStorageMockForIDInput{p, p1}, "RecordStorage.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordStorageMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to RecordStorageMock.ForID. %v %v", p, p1)
		return
	}

	return m.ForIDFunc(p, p1)
}

// ForIDMinimockCounter returns a count of RecordStorageMock.ForIDFunc invocations
func (m *RecordStorageMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

// ForIDMinimockPreCounter returns the value of RecordStorageMock.ForID invocations
func (m *RecordStorageMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

// ForIDFinished returns true if mock invocations count is ok
func (m *RecordStorageMock) ForIDFinished() bool {
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

type mRecordStorageMockSet struct {
	mock              *RecordStorageMock
	mainExpectation   *RecordStorageMockSetExpectation
	expectationSeries []*RecordStorageMockSetExpectation
}

type RecordStorageMockSetExpectation struct {
	input  *RecordStorageMockSetInput
	result *RecordStorageMockSetResult
}

type RecordStorageMockSetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 record.Material
}

type RecordStorageMockSetResult struct {
	r error
}

// Expect specifies that invocation of RecordStorage.Set is expected from 1 to Infinity times
func (m *mRecordStorageMockSet) Expect(p context.Context, p1 insolar.ID, p2 record.Material) *mRecordStorageMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordStorageMockSetExpectation{}
	}
	m.mainExpectation.input = &RecordStorageMockSetInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of RecordStorage.Set
func (m *mRecordStorageMockSet) Return(r error) *RecordStorageMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordStorageMockSetExpectation{}
	}
	m.mainExpectation.result = &RecordStorageMockSetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of RecordStorage.Set is expected once
func (m *mRecordStorageMockSet) ExpectOnce(p context.Context, p1 insolar.ID, p2 record.Material) *RecordStorageMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &RecordStorageMockSetExpectation{}
	expectation.input = &RecordStorageMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordStorageMockSetExpectation) Return(r error) {
	e.result = &RecordStorageMockSetResult{r}
}

// Set uses given function f as a mock of RecordStorage.Set method
func (m *mRecordStorageMockSet) Set(f func(p context.Context, p1 insolar.ID, p2 record.Material) (r error)) *RecordStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

// Set implements github.com/insolar/insolar/ledger/object.RecordStorage interface
func (m *RecordStorageMock) Set(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordStorageMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordStorageMockSetInput{p, p1, p2}, "RecordStorage.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordStorageMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordStorageMockSetInput{p, p1, p2}, "RecordStorage.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordStorageMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to RecordStorageMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

// SetMinimockCounter returns a count of RecordStorageMock.SetFunc invocations
func (m *RecordStorageMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

// SetMinimockPreCounter returns the value of RecordStorageMock.Set invocations
func (m *RecordStorageMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

// SetFinished returns true if mock invocations count is ok
func (m *RecordStorageMock) SetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCounter) == uint64(len(m.SetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFunc != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordStorageMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to RecordStorageMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to RecordStorageMock.Set")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordStorageMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecordStorageMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecordStorageMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to RecordStorageMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to RecordStorageMock.Set")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecordStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *RecordStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to RecordStorageMock.ForID")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to RecordStorageMock.Set")
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
func (m *RecordStorageMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}

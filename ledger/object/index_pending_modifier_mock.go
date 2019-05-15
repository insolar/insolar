package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexPendingModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexPendingModifierMock implements github.com/insolar/insolar/ledger/object.IndexPendingModifier
type IndexPendingModifierMock struct {
	t minimock.Tester

	SetRequestFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)
	SetRequestCounter    uint64
	SetRequestPreCounter uint64
	SetRequestMock       mIndexPendingModifierMockSetRequest

	SetResultRecordFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)
	SetResultRecordCounter    uint64
	SetResultRecordPreCounter uint64
	SetResultRecordMock       mIndexPendingModifierMockSetResultRecord
}

// NewIndexPendingModifierMock returns a mock for github.com/insolar/insolar/ledger/object.IndexPendingModifier
func NewIndexPendingModifierMock(t minimock.Tester) *IndexPendingModifierMock {
	m := &IndexPendingModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetRequestMock = mIndexPendingModifierMockSetRequest{mock: m}
	m.SetResultRecordMock = mIndexPendingModifierMockSetResultRecord{mock: m}

	return m
}

type mIndexPendingModifierMockSetRequest struct {
	mock              *IndexPendingModifierMock
	mainExpectation   *IndexPendingModifierMockSetRequestExpectation
	expectationSeries []*IndexPendingModifierMockSetRequestExpectation
}

type IndexPendingModifierMockSetRequestExpectation struct {
	input  *IndexPendingModifierMockSetRequestInput
	result *IndexPendingModifierMockSetRequestResult
}

type IndexPendingModifierMockSetRequestInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
}

type IndexPendingModifierMockSetRequestResult struct {
	r error
}

// Expect specifies that invocation of IndexPendingModifier.SetRequest is expected from 1 to Infinity times
func (m *mIndexPendingModifierMockSetRequest) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *mIndexPendingModifierMockSetRequest {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexPendingModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.input = &IndexPendingModifierMockSetRequestInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of IndexPendingModifier.SetRequest
func (m *mIndexPendingModifierMockSetRequest) Return(r error) *IndexPendingModifierMock {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexPendingModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.result = &IndexPendingModifierMockSetRequestResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexPendingModifier.SetRequest is expected once
func (m *mIndexPendingModifierMockSetRequest) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *IndexPendingModifierMockSetRequestExpectation {
	m.mock.SetRequestFunc = nil
	m.mainExpectation = nil

	expectation := &IndexPendingModifierMockSetRequestExpectation{}
	expectation.input = &IndexPendingModifierMockSetRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexPendingModifierMockSetRequestExpectation) Return(r error) {
	e.result = &IndexPendingModifierMockSetRequestResult{r}
}

// Set uses given function f as a mock of IndexPendingModifier.SetRequest method
func (m *mIndexPendingModifierMockSetRequest) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)) *IndexPendingModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetRequestFunc = f
	return m.mock
}

// SetRequest implements github.com/insolar/insolar/ledger/object.IndexPendingModifier interface
func (m *IndexPendingModifierMock) SetRequest(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SetRequestCounter, 1)

	if len(m.SetRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexPendingModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexPendingModifierMockSetRequestInput{p, p1, p2, p3}, "IndexPendingModifier.SetRequest got unexpected parameters")

		result := m.SetRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexPendingModifierMock.SetRequest")
			return
		}

		r = result.r

		return
	}

	if m.SetRequestMock.mainExpectation != nil {

		input := m.SetRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexPendingModifierMockSetRequestInput{p, p1, p2, p3}, "IndexPendingModifier.SetRequest got unexpected parameters")
		}

		result := m.SetRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexPendingModifierMock.SetRequest")
		}

		r = result.r

		return
	}

	if m.SetRequestFunc == nil {
		m.t.Fatalf("Unexpected call to IndexPendingModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetRequestFunc(p, p1, p2, p3)
}

// SetRequestMinimockCounter returns a count of IndexPendingModifierMock.SetRequestFunc invocations
func (m *IndexPendingModifierMock) SetRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestCounter)
}

// SetRequestMinimockPreCounter returns the value of IndexPendingModifierMock.SetRequest invocations
func (m *IndexPendingModifierMock) SetRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestPreCounter)
}

// SetRequestFinished returns true if mock invocations count is ok
func (m *IndexPendingModifierMock) SetRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetRequestCounter) == uint64(len(m.SetRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetRequestFunc != nil {
		return atomic.LoadUint64(&m.SetRequestCounter) > 0
	}

	return true
}

type mIndexPendingModifierMockSetResultRecord struct {
	mock              *IndexPendingModifierMock
	mainExpectation   *IndexPendingModifierMockSetResultRecordExpectation
	expectationSeries []*IndexPendingModifierMockSetResultRecordExpectation
}

type IndexPendingModifierMockSetResultRecordExpectation struct {
	input  *IndexPendingModifierMockSetResultRecordInput
	result *IndexPendingModifierMockSetResultRecordResult
}

type IndexPendingModifierMockSetResultRecordInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
}

type IndexPendingModifierMockSetResultRecordResult struct {
	r error
}

// Expect specifies that invocation of IndexPendingModifier.SetResultRecord is expected from 1 to Infinity times
func (m *mIndexPendingModifierMockSetResultRecord) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *mIndexPendingModifierMockSetResultRecord {
	m.mock.SetResultRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexPendingModifierMockSetResultRecordExpectation{}
	}
	m.mainExpectation.input = &IndexPendingModifierMockSetResultRecordInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of IndexPendingModifier.SetResultRecord
func (m *mIndexPendingModifierMockSetResultRecord) Return(r error) *IndexPendingModifierMock {
	m.mock.SetResultRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexPendingModifierMockSetResultRecordExpectation{}
	}
	m.mainExpectation.result = &IndexPendingModifierMockSetResultRecordResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexPendingModifier.SetResultRecord is expected once
func (m *mIndexPendingModifierMockSetResultRecord) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *IndexPendingModifierMockSetResultRecordExpectation {
	m.mock.SetResultRecordFunc = nil
	m.mainExpectation = nil

	expectation := &IndexPendingModifierMockSetResultRecordExpectation{}
	expectation.input = &IndexPendingModifierMockSetResultRecordInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexPendingModifierMockSetResultRecordExpectation) Return(r error) {
	e.result = &IndexPendingModifierMockSetResultRecordResult{r}
}

// Set uses given function f as a mock of IndexPendingModifier.SetResultRecord method
func (m *mIndexPendingModifierMockSetResultRecord) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)) *IndexPendingModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetResultRecordFunc = f
	return m.mock
}

// SetResultRecord implements github.com/insolar/insolar/ledger/object.IndexPendingModifier interface
func (m *IndexPendingModifierMock) SetResultRecord(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetResultRecordPreCounter, 1)
	defer atomic.AddUint64(&m.SetResultRecordCounter, 1)

	if len(m.SetResultRecordMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetResultRecordMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexPendingModifierMock.SetResultRecord. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetResultRecordMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexPendingModifierMockSetResultRecordInput{p, p1, p2, p3}, "IndexPendingModifier.SetResultRecord got unexpected parameters")

		result := m.SetResultRecordMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexPendingModifierMock.SetResultRecord")
			return
		}

		r = result.r

		return
	}

	if m.SetResultRecordMock.mainExpectation != nil {

		input := m.SetResultRecordMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexPendingModifierMockSetResultRecordInput{p, p1, p2, p3}, "IndexPendingModifier.SetResultRecord got unexpected parameters")
		}

		result := m.SetResultRecordMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexPendingModifierMock.SetResultRecord")
		}

		r = result.r

		return
	}

	if m.SetResultRecordFunc == nil {
		m.t.Fatalf("Unexpected call to IndexPendingModifierMock.SetResultRecord. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetResultRecordFunc(p, p1, p2, p3)
}

// SetResultRecordMinimockCounter returns a count of IndexPendingModifierMock.SetResultRecordFunc invocations
func (m *IndexPendingModifierMock) SetResultRecordMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultRecordCounter)
}

// SetResultRecordMinimockPreCounter returns the value of IndexPendingModifierMock.SetResultRecord invocations
func (m *IndexPendingModifierMock) SetResultRecordMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultRecordPreCounter)
}

// SetResultRecordFinished returns true if mock invocations count is ok
func (m *IndexPendingModifierMock) SetResultRecordFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetResultRecordMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetResultRecordCounter) == uint64(len(m.SetResultRecordMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetResultRecordMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetResultRecordCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetResultRecordFunc != nil {
		return atomic.LoadUint64(&m.SetResultRecordCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexPendingModifierMock) ValidateCallCounters() {

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to IndexPendingModifierMock.SetRequest")
	}

	if !m.SetResultRecordFinished() {
		m.t.Fatal("Expected call to IndexPendingModifierMock.SetResultRecord")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexPendingModifierMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexPendingModifierMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexPendingModifierMock) MinimockFinish() {

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to IndexPendingModifierMock.SetRequest")
	}

	if !m.SetResultRecordFinished() {
		m.t.Fatal("Expected call to IndexPendingModifierMock.SetResultRecord")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexPendingModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexPendingModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetRequestFinished()
		ok = ok && m.SetResultRecordFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetRequestFinished() {
				m.t.Error("Expected call to IndexPendingModifierMock.SetRequest")
			}

			if !m.SetResultRecordFinished() {
				m.t.Error("Expected call to IndexPendingModifierMock.SetResultRecord")
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
func (m *IndexPendingModifierMock) AllMocksCalled() bool {

	if !m.SetRequestFinished() {
		return false
	}

	if !m.SetResultRecordFinished() {
		return false
	}

	return true
}

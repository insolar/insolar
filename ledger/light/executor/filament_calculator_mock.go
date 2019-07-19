package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FilamentCalculator" can be found in github.com/insolar/insolar/ledger/light/executor
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

//FilamentCalculatorMock implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator
type FilamentCalculatorMock struct {
	t minimock.Tester

	FindRecordFunc       func(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.ID) (r record.CompositeFilamentRecord, r1 error)
	FindRecordCounter    uint64
	FindRecordPreCounter uint64
	FindRecordMock       mFilamentCalculatorMockFindRecord

	PendingRequestsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.CompositeFilamentRecord, r1 error)
	PendingRequestsCounter    uint64
	PendingRequestsPreCounter uint64
	PendingRequestsMock       mFilamentCalculatorMockPendingRequests

	RequestDuplicateFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error)
	RequestDuplicateCounter    uint64
	RequestDuplicatePreCounter uint64
	RequestDuplicateMock       mFilamentCalculatorMockRequestDuplicate

	RequestsFunc       func(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber) (r []record.CompositeFilamentRecord, r1 error)
	RequestsCounter    uint64
	RequestsPreCounter uint64
	RequestsMock       mFilamentCalculatorMockRequests

	ResultDuplicateFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Result) (r *record.CompositeFilamentRecord, r1 error)
	ResultDuplicateCounter    uint64
	ResultDuplicatePreCounter uint64
	ResultDuplicateMock       mFilamentCalculatorMockResultDuplicate
}

//NewFilamentCalculatorMock returns a mock for github.com/insolar/insolar/ledger/light/executor.FilamentCalculator
func NewFilamentCalculatorMock(t minimock.Tester) *FilamentCalculatorMock {
	m := &FilamentCalculatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FindRecordMock = mFilamentCalculatorMockFindRecord{mock: m}
	m.PendingRequestsMock = mFilamentCalculatorMockPendingRequests{mock: m}
	m.RequestDuplicateMock = mFilamentCalculatorMockRequestDuplicate{mock: m}
	m.RequestsMock = mFilamentCalculatorMockRequests{mock: m}
	m.ResultDuplicateMock = mFilamentCalculatorMockResultDuplicate{mock: m}

	return m
}

type mFilamentCalculatorMockFindRecord struct {
	mock              *FilamentCalculatorMock
	mainExpectation   *FilamentCalculatorMockFindRecordExpectation
	expectationSeries []*FilamentCalculatorMockFindRecordExpectation
}

type FilamentCalculatorMockFindRecordExpectation struct {
	input  *FilamentCalculatorMockFindRecordInput
	result *FilamentCalculatorMockFindRecordResult
}

type FilamentCalculatorMockFindRecordInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.ID
	p3 insolar.ID
}

type FilamentCalculatorMockFindRecordResult struct {
	r  record.CompositeFilamentRecord
	r1 error
}

//Expect specifies that invocation of FilamentCalculator.FindRecord is expected from 1 to Infinity times
func (m *mFilamentCalculatorMockFindRecord) Expect(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.ID) *mFilamentCalculatorMockFindRecord {
	m.mock.FindRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockFindRecordExpectation{}
	}
	m.mainExpectation.input = &FilamentCalculatorMockFindRecordInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of FilamentCalculator.FindRecord
func (m *mFilamentCalculatorMockFindRecord) Return(r record.CompositeFilamentRecord, r1 error) *FilamentCalculatorMock {
	m.mock.FindRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockFindRecordExpectation{}
	}
	m.mainExpectation.result = &FilamentCalculatorMockFindRecordResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCalculator.FindRecord is expected once
func (m *mFilamentCalculatorMockFindRecord) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.ID) *FilamentCalculatorMockFindRecordExpectation {
	m.mock.FindRecordFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCalculatorMockFindRecordExpectation{}
	expectation.input = &FilamentCalculatorMockFindRecordInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCalculatorMockFindRecordExpectation) Return(r record.CompositeFilamentRecord, r1 error) {
	e.result = &FilamentCalculatorMockFindRecordResult{r, r1}
}

//Set uses given function f as a mock of FilamentCalculator.FindRecord method
func (m *mFilamentCalculatorMockFindRecord) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.ID) (r record.CompositeFilamentRecord, r1 error)) *FilamentCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FindRecordFunc = f
	return m.mock
}

//FindRecord implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator interface
func (m *FilamentCalculatorMock) FindRecord(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.ID) (r record.CompositeFilamentRecord, r1 error) {
	counter := atomic.AddUint64(&m.FindRecordPreCounter, 1)
	defer atomic.AddUint64(&m.FindRecordCounter, 1)

	if len(m.FindRecordMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FindRecordMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCalculatorMock.FindRecord. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.FindRecordMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCalculatorMockFindRecordInput{p, p1, p2, p3}, "FilamentCalculator.FindRecord got unexpected parameters")

		result := m.FindRecordMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.FindRecord")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FindRecordMock.mainExpectation != nil {

		input := m.FindRecordMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCalculatorMockFindRecordInput{p, p1, p2, p3}, "FilamentCalculator.FindRecord got unexpected parameters")
		}

		result := m.FindRecordMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.FindRecord")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FindRecordFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCalculatorMock.FindRecord. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.FindRecordFunc(p, p1, p2, p3)
}

//FindRecordMinimockCounter returns a count of FilamentCalculatorMock.FindRecordFunc invocations
func (m *FilamentCalculatorMock) FindRecordMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FindRecordCounter)
}

//FindRecordMinimockPreCounter returns the value of FilamentCalculatorMock.FindRecord invocations
func (m *FilamentCalculatorMock) FindRecordMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FindRecordPreCounter)
}

//FindRecordFinished returns true if mock invocations count is ok
func (m *FilamentCalculatorMock) FindRecordFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FindRecordMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FindRecordCounter) == uint64(len(m.FindRecordMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FindRecordMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FindRecordCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FindRecordFunc != nil {
		return atomic.LoadUint64(&m.FindRecordCounter) > 0
	}

	return true
}

type mFilamentCalculatorMockPendingRequests struct {
	mock              *FilamentCalculatorMock
	mainExpectation   *FilamentCalculatorMockPendingRequestsExpectation
	expectationSeries []*FilamentCalculatorMockPendingRequestsExpectation
}

type FilamentCalculatorMockPendingRequestsExpectation struct {
	input  *FilamentCalculatorMockPendingRequestsInput
	result *FilamentCalculatorMockPendingRequestsResult
}

type FilamentCalculatorMockPendingRequestsInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type FilamentCalculatorMockPendingRequestsResult struct {
	r  []record.CompositeFilamentRecord
	r1 error
}

//Expect specifies that invocation of FilamentCalculator.PendingRequests is expected from 1 to Infinity times
func (m *mFilamentCalculatorMockPendingRequests) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mFilamentCalculatorMockPendingRequests {
	m.mock.PendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockPendingRequestsExpectation{}
	}
	m.mainExpectation.input = &FilamentCalculatorMockPendingRequestsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of FilamentCalculator.PendingRequests
func (m *mFilamentCalculatorMockPendingRequests) Return(r []record.CompositeFilamentRecord, r1 error) *FilamentCalculatorMock {
	m.mock.PendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockPendingRequestsExpectation{}
	}
	m.mainExpectation.result = &FilamentCalculatorMockPendingRequestsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCalculator.PendingRequests is expected once
func (m *mFilamentCalculatorMockPendingRequests) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *FilamentCalculatorMockPendingRequestsExpectation {
	m.mock.PendingRequestsFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCalculatorMockPendingRequestsExpectation{}
	expectation.input = &FilamentCalculatorMockPendingRequestsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCalculatorMockPendingRequestsExpectation) Return(r []record.CompositeFilamentRecord, r1 error) {
	e.result = &FilamentCalculatorMockPendingRequestsResult{r, r1}
}

//Set uses given function f as a mock of FilamentCalculator.PendingRequests method
func (m *mFilamentCalculatorMockPendingRequests) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.CompositeFilamentRecord, r1 error)) *FilamentCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PendingRequestsFunc = f
	return m.mock
}

//PendingRequests implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator interface
func (m *FilamentCalculatorMock) PendingRequests(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.CompositeFilamentRecord, r1 error) {
	counter := atomic.AddUint64(&m.PendingRequestsPreCounter, 1)
	defer atomic.AddUint64(&m.PendingRequestsCounter, 1)

	if len(m.PendingRequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PendingRequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCalculatorMock.PendingRequests. %v %v %v", p, p1, p2)
			return
		}

		input := m.PendingRequestsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCalculatorMockPendingRequestsInput{p, p1, p2}, "FilamentCalculator.PendingRequests got unexpected parameters")

		result := m.PendingRequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.PendingRequests")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PendingRequestsMock.mainExpectation != nil {

		input := m.PendingRequestsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCalculatorMockPendingRequestsInput{p, p1, p2}, "FilamentCalculator.PendingRequests got unexpected parameters")
		}

		result := m.PendingRequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.PendingRequests")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PendingRequestsFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCalculatorMock.PendingRequests. %v %v %v", p, p1, p2)
		return
	}

	return m.PendingRequestsFunc(p, p1, p2)
}

//PendingRequestsMinimockCounter returns a count of FilamentCalculatorMock.PendingRequestsFunc invocations
func (m *FilamentCalculatorMock) PendingRequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PendingRequestsCounter)
}

//PendingRequestsMinimockPreCounter returns the value of FilamentCalculatorMock.PendingRequests invocations
func (m *FilamentCalculatorMock) PendingRequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PendingRequestsPreCounter)
}

//PendingRequestsFinished returns true if mock invocations count is ok
func (m *FilamentCalculatorMock) PendingRequestsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PendingRequestsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PendingRequestsCounter) == uint64(len(m.PendingRequestsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PendingRequestsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PendingRequestsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PendingRequestsFunc != nil {
		return atomic.LoadUint64(&m.PendingRequestsCounter) > 0
	}

	return true
}

type mFilamentCalculatorMockRequestDuplicate struct {
	mock              *FilamentCalculatorMock
	mainExpectation   *FilamentCalculatorMockRequestDuplicateExpectation
	expectationSeries []*FilamentCalculatorMockRequestDuplicateExpectation
}

type FilamentCalculatorMockRequestDuplicateExpectation struct {
	input  *FilamentCalculatorMockRequestDuplicateInput
	result *FilamentCalculatorMockRequestDuplicateResult
}

type FilamentCalculatorMockRequestDuplicateInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
	p4 record.Request
}

type FilamentCalculatorMockRequestDuplicateResult struct {
	r  *record.CompositeFilamentRecord
	r1 *record.CompositeFilamentRecord
	r2 error
}

//Expect specifies that invocation of FilamentCalculator.RequestDuplicate is expected from 1 to Infinity times
func (m *mFilamentCalculatorMockRequestDuplicate) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Request) *mFilamentCalculatorMockRequestDuplicate {
	m.mock.RequestDuplicateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockRequestDuplicateExpectation{}
	}
	m.mainExpectation.input = &FilamentCalculatorMockRequestDuplicateInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of FilamentCalculator.RequestDuplicate
func (m *mFilamentCalculatorMockRequestDuplicate) Return(r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) *FilamentCalculatorMock {
	m.mock.RequestDuplicateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockRequestDuplicateExpectation{}
	}
	m.mainExpectation.result = &FilamentCalculatorMockRequestDuplicateResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCalculator.RequestDuplicate is expected once
func (m *mFilamentCalculatorMockRequestDuplicate) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Request) *FilamentCalculatorMockRequestDuplicateExpectation {
	m.mock.RequestDuplicateFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCalculatorMockRequestDuplicateExpectation{}
	expectation.input = &FilamentCalculatorMockRequestDuplicateInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCalculatorMockRequestDuplicateExpectation) Return(r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) {
	e.result = &FilamentCalculatorMockRequestDuplicateResult{r, r1, r2}
}

//Set uses given function f as a mock of FilamentCalculator.RequestDuplicate method
func (m *mFilamentCalculatorMockRequestDuplicate) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error)) *FilamentCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RequestDuplicateFunc = f
	return m.mock
}

//RequestDuplicate implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator interface
func (m *FilamentCalculatorMock) RequestDuplicate(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) {
	counter := atomic.AddUint64(&m.RequestDuplicatePreCounter, 1)
	defer atomic.AddUint64(&m.RequestDuplicateCounter, 1)

	if len(m.RequestDuplicateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RequestDuplicateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCalculatorMock.RequestDuplicate. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.RequestDuplicateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCalculatorMockRequestDuplicateInput{p, p1, p2, p3, p4}, "FilamentCalculator.RequestDuplicate got unexpected parameters")

		result := m.RequestDuplicateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.RequestDuplicate")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.RequestDuplicateMock.mainExpectation != nil {

		input := m.RequestDuplicateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCalculatorMockRequestDuplicateInput{p, p1, p2, p3, p4}, "FilamentCalculator.RequestDuplicate got unexpected parameters")
		}

		result := m.RequestDuplicateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.RequestDuplicate")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.RequestDuplicateFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCalculatorMock.RequestDuplicate. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.RequestDuplicateFunc(p, p1, p2, p3, p4)
}

//RequestDuplicateMinimockCounter returns a count of FilamentCalculatorMock.RequestDuplicateFunc invocations
func (m *FilamentCalculatorMock) RequestDuplicateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RequestDuplicateCounter)
}

//RequestDuplicateMinimockPreCounter returns the value of FilamentCalculatorMock.RequestDuplicate invocations
func (m *FilamentCalculatorMock) RequestDuplicateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RequestDuplicatePreCounter)
}

//RequestDuplicateFinished returns true if mock invocations count is ok
func (m *FilamentCalculatorMock) RequestDuplicateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RequestDuplicateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RequestDuplicateCounter) == uint64(len(m.RequestDuplicateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RequestDuplicateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RequestDuplicateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RequestDuplicateFunc != nil {
		return atomic.LoadUint64(&m.RequestDuplicateCounter) > 0
	}

	return true
}

type mFilamentCalculatorMockRequests struct {
	mock              *FilamentCalculatorMock
	mainExpectation   *FilamentCalculatorMockRequestsExpectation
	expectationSeries []*FilamentCalculatorMockRequestsExpectation
}

type FilamentCalculatorMockRequestsExpectation struct {
	input  *FilamentCalculatorMockRequestsInput
	result *FilamentCalculatorMockRequestsResult
}

type FilamentCalculatorMockRequestsInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.ID
	p3 insolar.PulseNumber
}

type FilamentCalculatorMockRequestsResult struct {
	r  []record.CompositeFilamentRecord
	r1 error
}

//Expect specifies that invocation of FilamentCalculator.Requests is expected from 1 to Infinity times
func (m *mFilamentCalculatorMockRequests) Expect(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber) *mFilamentCalculatorMockRequests {
	m.mock.RequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockRequestsExpectation{}
	}
	m.mainExpectation.input = &FilamentCalculatorMockRequestsInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of FilamentCalculator.Requests
func (m *mFilamentCalculatorMockRequests) Return(r []record.CompositeFilamentRecord, r1 error) *FilamentCalculatorMock {
	m.mock.RequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockRequestsExpectation{}
	}
	m.mainExpectation.result = &FilamentCalculatorMockRequestsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCalculator.Requests is expected once
func (m *mFilamentCalculatorMockRequests) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber) *FilamentCalculatorMockRequestsExpectation {
	m.mock.RequestsFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCalculatorMockRequestsExpectation{}
	expectation.input = &FilamentCalculatorMockRequestsInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCalculatorMockRequestsExpectation) Return(r []record.CompositeFilamentRecord, r1 error) {
	e.result = &FilamentCalculatorMockRequestsResult{r, r1}
}

//Set uses given function f as a mock of FilamentCalculator.Requests method
func (m *mFilamentCalculatorMockRequests) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber) (r []record.CompositeFilamentRecord, r1 error)) *FilamentCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RequestsFunc = f
	return m.mock
}

//Requests implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator interface
func (m *FilamentCalculatorMock) Requests(p context.Context, p1 insolar.ID, p2 insolar.ID, p3 insolar.PulseNumber) (r []record.CompositeFilamentRecord, r1 error) {
	counter := atomic.AddUint64(&m.RequestsPreCounter, 1)
	defer atomic.AddUint64(&m.RequestsCounter, 1)

	if len(m.RequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCalculatorMock.Requests. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RequestsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCalculatorMockRequestsInput{p, p1, p2, p3}, "FilamentCalculator.Requests got unexpected parameters")

		result := m.RequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.Requests")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RequestsMock.mainExpectation != nil {

		input := m.RequestsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCalculatorMockRequestsInput{p, p1, p2, p3}, "FilamentCalculator.Requests got unexpected parameters")
		}

		result := m.RequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.Requests")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RequestsFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCalculatorMock.Requests. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RequestsFunc(p, p1, p2, p3)
}

//RequestsMinimockCounter returns a count of FilamentCalculatorMock.RequestsFunc invocations
func (m *FilamentCalculatorMock) RequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RequestsCounter)
}

//RequestsMinimockPreCounter returns the value of FilamentCalculatorMock.Requests invocations
func (m *FilamentCalculatorMock) RequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RequestsPreCounter)
}

//RequestsFinished returns true if mock invocations count is ok
func (m *FilamentCalculatorMock) RequestsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RequestsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RequestsCounter) == uint64(len(m.RequestsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RequestsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RequestsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RequestsFunc != nil {
		return atomic.LoadUint64(&m.RequestsCounter) > 0
	}

	return true
}

type mFilamentCalculatorMockResultDuplicate struct {
	mock              *FilamentCalculatorMock
	mainExpectation   *FilamentCalculatorMockResultDuplicateExpectation
	expectationSeries []*FilamentCalculatorMockResultDuplicateExpectation
}

type FilamentCalculatorMockResultDuplicateExpectation struct {
	input  *FilamentCalculatorMockResultDuplicateInput
	result *FilamentCalculatorMockResultDuplicateResult
}

type FilamentCalculatorMockResultDuplicateInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
	p4 record.Result
}

type FilamentCalculatorMockResultDuplicateResult struct {
	r  *record.CompositeFilamentRecord
	r1 error
}

//Expect specifies that invocation of FilamentCalculator.ResultDuplicate is expected from 1 to Infinity times
func (m *mFilamentCalculatorMockResultDuplicate) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Result) *mFilamentCalculatorMockResultDuplicate {
	m.mock.ResultDuplicateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockResultDuplicateExpectation{}
	}
	m.mainExpectation.input = &FilamentCalculatorMockResultDuplicateInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of FilamentCalculator.ResultDuplicate
func (m *mFilamentCalculatorMockResultDuplicate) Return(r *record.CompositeFilamentRecord, r1 error) *FilamentCalculatorMock {
	m.mock.ResultDuplicateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCalculatorMockResultDuplicateExpectation{}
	}
	m.mainExpectation.result = &FilamentCalculatorMockResultDuplicateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCalculator.ResultDuplicate is expected once
func (m *mFilamentCalculatorMockResultDuplicate) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Result) *FilamentCalculatorMockResultDuplicateExpectation {
	m.mock.ResultDuplicateFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCalculatorMockResultDuplicateExpectation{}
	expectation.input = &FilamentCalculatorMockResultDuplicateInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCalculatorMockResultDuplicateExpectation) Return(r *record.CompositeFilamentRecord, r1 error) {
	e.result = &FilamentCalculatorMockResultDuplicateResult{r, r1}
}

//Set uses given function f as a mock of FilamentCalculator.ResultDuplicate method
func (m *mFilamentCalculatorMockResultDuplicate) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Result) (r *record.CompositeFilamentRecord, r1 error)) *FilamentCalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResultDuplicateFunc = f
	return m.mock
}

//ResultDuplicate implements github.com/insolar/insolar/ledger/light/executor.FilamentCalculator interface
func (m *FilamentCalculatorMock) ResultDuplicate(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Result) (r *record.CompositeFilamentRecord, r1 error) {
	counter := atomic.AddUint64(&m.ResultDuplicatePreCounter, 1)
	defer atomic.AddUint64(&m.ResultDuplicateCounter, 1)

	if len(m.ResultDuplicateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ResultDuplicateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCalculatorMock.ResultDuplicate. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.ResultDuplicateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCalculatorMockResultDuplicateInput{p, p1, p2, p3, p4}, "FilamentCalculator.ResultDuplicate got unexpected parameters")

		result := m.ResultDuplicateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.ResultDuplicate")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResultDuplicateMock.mainExpectation != nil {

		input := m.ResultDuplicateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCalculatorMockResultDuplicateInput{p, p1, p2, p3, p4}, "FilamentCalculator.ResultDuplicate got unexpected parameters")
		}

		result := m.ResultDuplicateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCalculatorMock.ResultDuplicate")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResultDuplicateFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCalculatorMock.ResultDuplicate. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.ResultDuplicateFunc(p, p1, p2, p3, p4)
}

//ResultDuplicateMinimockCounter returns a count of FilamentCalculatorMock.ResultDuplicateFunc invocations
func (m *FilamentCalculatorMock) ResultDuplicateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResultDuplicateCounter)
}

//ResultDuplicateMinimockPreCounter returns the value of FilamentCalculatorMock.ResultDuplicate invocations
func (m *FilamentCalculatorMock) ResultDuplicateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResultDuplicatePreCounter)
}

//ResultDuplicateFinished returns true if mock invocations count is ok
func (m *FilamentCalculatorMock) ResultDuplicateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ResultDuplicateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ResultDuplicateCounter) == uint64(len(m.ResultDuplicateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ResultDuplicateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ResultDuplicateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ResultDuplicateFunc != nil {
		return atomic.LoadUint64(&m.ResultDuplicateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCalculatorMock) ValidateCallCounters() {

	if !m.FindRecordFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.FindRecord")
	}

	if !m.PendingRequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.PendingRequests")
	}

	if !m.RequestDuplicateFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.RequestDuplicate")
	}

	if !m.RequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.Requests")
	}

	if !m.ResultDuplicateFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.ResultDuplicate")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCalculatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FilamentCalculatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FilamentCalculatorMock) MinimockFinish() {

	if !m.FindRecordFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.FindRecord")
	}

	if !m.PendingRequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.PendingRequests")
	}

	if !m.RequestDuplicateFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.RequestDuplicate")
	}

	if !m.RequestsFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.Requests")
	}

	if !m.ResultDuplicateFinished() {
		m.t.Fatal("Expected call to FilamentCalculatorMock.ResultDuplicate")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FilamentCalculatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FilamentCalculatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FindRecordFinished()
		ok = ok && m.PendingRequestsFinished()
		ok = ok && m.RequestDuplicateFinished()
		ok = ok && m.RequestsFinished()
		ok = ok && m.ResultDuplicateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FindRecordFinished() {
				m.t.Error("Expected call to FilamentCalculatorMock.FindRecord")
			}

			if !m.PendingRequestsFinished() {
				m.t.Error("Expected call to FilamentCalculatorMock.PendingRequests")
			}

			if !m.RequestDuplicateFinished() {
				m.t.Error("Expected call to FilamentCalculatorMock.RequestDuplicate")
			}

			if !m.RequestsFinished() {
				m.t.Error("Expected call to FilamentCalculatorMock.Requests")
			}

			if !m.ResultDuplicateFinished() {
				m.t.Error("Expected call to FilamentCalculatorMock.ResultDuplicate")
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
func (m *FilamentCalculatorMock) AllMocksCalled() bool {

	if !m.FindRecordFinished() {
		return false
	}

	if !m.PendingRequestsFinished() {
		return false
	}

	if !m.RequestDuplicateFinished() {
		return false
	}

	if !m.RequestsFinished() {
		return false
	}

	if !m.ResultDuplicateFinished() {
		return false
	}

	return true
}

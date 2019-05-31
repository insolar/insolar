package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PendingAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"

	testify_assert "github.com/stretchr/testify/assert"
)

// PendingAccessorMock implements github.com/insolar/insolar/ledger/object.PendingAccessor
type PendingAccessorMock struct {
	t minimock.Tester

	MetaForObjIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r object.PendingMeta, r1 error)
	MetaForObjIDCounter    uint64
	MetaForObjIDPreCounter uint64
	MetaForObjIDMock       mPendingAccessorMockMetaForObjID

	OpenRequestsForObjIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []record.Request, r1 error)
	OpenRequestsForObjIDCounter    uint64
	OpenRequestsForObjIDPreCounter uint64
	OpenRequestsForObjIDMock       mPendingAccessorMockOpenRequestsForObjID

	RecordsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.Virtual, r1 error)
	RecordsCounter    uint64
	RecordsPreCounter uint64
	RecordsMock       mPendingAccessorMockRecords
}

// NewPendingAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.PendingAccessor
func NewPendingAccessorMock(t minimock.Tester) *PendingAccessorMock {
	m := &PendingAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.MetaForObjIDMock = mPendingAccessorMockMetaForObjID{mock: m}
	m.OpenRequestsForObjIDMock = mPendingAccessorMockOpenRequestsForObjID{mock: m}
	m.RecordsMock = mPendingAccessorMockRecords{mock: m}

	return m
}

type mPendingAccessorMockMetaForObjID struct {
	mock              *PendingAccessorMock
	mainExpectation   *PendingAccessorMockMetaForObjIDExpectation
	expectationSeries []*PendingAccessorMockMetaForObjIDExpectation
}

type PendingAccessorMockMetaForObjIDExpectation struct {
	input  *PendingAccessorMockMetaForObjIDInput
	result *PendingAccessorMockMetaForObjIDResult
}

type PendingAccessorMockMetaForObjIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type PendingAccessorMockMetaForObjIDResult struct {
	r  object.PendingMeta
	r1 error
}

// Expect specifies that invocation of PendingAccessor.MetaForObjID is expected from 1 to Infinity times
func (m *mPendingAccessorMockMetaForObjID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mPendingAccessorMockMetaForObjID {
	m.mock.MetaForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockMetaForObjIDExpectation{}
	}
	m.mainExpectation.input = &PendingAccessorMockMetaForObjIDInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of PendingAccessor.MetaForObjID
func (m *mPendingAccessorMockMetaForObjID) Return(r object.PendingMeta, r1 error) *PendingAccessorMock {
	m.mock.MetaForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockMetaForObjIDExpectation{}
	}
	m.mainExpectation.result = &PendingAccessorMockMetaForObjIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of PendingAccessor.MetaForObjID is expected once
func (m *mPendingAccessorMockMetaForObjID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *PendingAccessorMockMetaForObjIDExpectation {
	m.mock.MetaForObjIDFunc = nil
	m.mainExpectation = nil

	expectation := &PendingAccessorMockMetaForObjIDExpectation{}
	expectation.input = &PendingAccessorMockMetaForObjIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingAccessorMockMetaForObjIDExpectation) Return(r object.PendingMeta, r1 error) {
	e.result = &PendingAccessorMockMetaForObjIDResult{r, r1}
}

// Set uses given function f as a mock of PendingAccessor.MetaForObjID method
func (m *mPendingAccessorMockMetaForObjID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r object.PendingMeta, r1 error)) *PendingAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MetaForObjIDFunc = f
	return m.mock
}

// MetaForObjID implements github.com/insolar/insolar/ledger/object.PendingAccessor interface
func (m *PendingAccessorMock) MetaForObjID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r object.PendingMeta, r1 error) {
	counter := atomic.AddUint64(&m.MetaForObjIDPreCounter, 1)
	defer atomic.AddUint64(&m.MetaForObjIDCounter, 1)

	if len(m.MetaForObjIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MetaForObjIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingAccessorMock.MetaForObjID. %v %v %v", p, p1, p2)
			return
		}

		input := m.MetaForObjIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingAccessorMockMetaForObjIDInput{p, p1, p2}, "PendingAccessor.MetaForObjID got unexpected parameters")

		result := m.MetaForObjIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.MetaForObjID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.MetaForObjIDMock.mainExpectation != nil {

		input := m.MetaForObjIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingAccessorMockMetaForObjIDInput{p, p1, p2}, "PendingAccessor.MetaForObjID got unexpected parameters")
		}

		result := m.MetaForObjIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.MetaForObjID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.MetaForObjIDFunc == nil {
		m.t.Fatalf("Unexpected call to PendingAccessorMock.MetaForObjID. %v %v %v", p, p1, p2)
		return
	}

	return m.MetaForObjIDFunc(p, p1, p2)
}

// MetaForObjIDMinimockCounter returns a count of PendingAccessorMock.MetaForObjIDFunc invocations
func (m *PendingAccessorMock) MetaForObjIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MetaForObjIDCounter)
}

// MetaForObjIDMinimockPreCounter returns the value of PendingAccessorMock.MetaForObjID invocations
func (m *PendingAccessorMock) MetaForObjIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MetaForObjIDPreCounter)
}

// MetaForObjIDFinished returns true if mock invocations count is ok
func (m *PendingAccessorMock) MetaForObjIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MetaForObjIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MetaForObjIDCounter) == uint64(len(m.MetaForObjIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MetaForObjIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MetaForObjIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MetaForObjIDFunc != nil {
		return atomic.LoadUint64(&m.MetaForObjIDCounter) > 0
	}

	return true
}

type mPendingAccessorMockOpenRequestsForObjID struct {
	mock              *PendingAccessorMock
	mainExpectation   *PendingAccessorMockOpenRequestsForObjIDExpectation
	expectationSeries []*PendingAccessorMockOpenRequestsForObjIDExpectation
}

type PendingAccessorMockOpenRequestsForObjIDExpectation struct {
	input  *PendingAccessorMockOpenRequestsForObjIDInput
	result *PendingAccessorMockOpenRequestsForObjIDResult
}

type PendingAccessorMockOpenRequestsForObjIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 int
}

type PendingAccessorMockOpenRequestsForObjIDResult struct {
	r  []record.Request
	r1 error
}

// Expect specifies that invocation of PendingAccessor.OpenRequestsForObjID is expected from 1 to Infinity times
func (m *mPendingAccessorMockOpenRequestsForObjID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) *mPendingAccessorMockOpenRequestsForObjID {
	m.mock.OpenRequestsForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockOpenRequestsForObjIDExpectation{}
	}
	m.mainExpectation.input = &PendingAccessorMockOpenRequestsForObjIDInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of PendingAccessor.OpenRequestsForObjID
func (m *mPendingAccessorMockOpenRequestsForObjID) Return(r []record.Request, r1 error) *PendingAccessorMock {
	m.mock.OpenRequestsForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockOpenRequestsForObjIDExpectation{}
	}
	m.mainExpectation.result = &PendingAccessorMockOpenRequestsForObjIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of PendingAccessor.OpenRequestsForObjID is expected once
func (m *mPendingAccessorMockOpenRequestsForObjID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) *PendingAccessorMockOpenRequestsForObjIDExpectation {
	m.mock.OpenRequestsForObjIDFunc = nil
	m.mainExpectation = nil

	expectation := &PendingAccessorMockOpenRequestsForObjIDExpectation{}
	expectation.input = &PendingAccessorMockOpenRequestsForObjIDInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingAccessorMockOpenRequestsForObjIDExpectation) Return(r []record.Request, r1 error) {
	e.result = &PendingAccessorMockOpenRequestsForObjIDResult{r, r1}
}

// Set uses given function f as a mock of PendingAccessor.OpenRequestsForObjID method
func (m *mPendingAccessorMockOpenRequestsForObjID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []record.Request, r1 error)) *PendingAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OpenRequestsForObjIDFunc = f
	return m.mock
}

// OpenRequestsForObjID implements github.com/insolar/insolar/ledger/object.PendingAccessor interface
func (m *PendingAccessorMock) OpenRequestsForObjID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []record.Request, r1 error) {
	counter := atomic.AddUint64(&m.OpenRequestsForObjIDPreCounter, 1)
	defer atomic.AddUint64(&m.OpenRequestsForObjIDCounter, 1)

	if len(m.OpenRequestsForObjIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OpenRequestsForObjIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingAccessorMock.OpenRequestsForObjID. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.OpenRequestsForObjIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingAccessorMockOpenRequestsForObjIDInput{p, p1, p2, p3}, "PendingAccessor.OpenRequestsForObjID got unexpected parameters")

		result := m.OpenRequestsForObjIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.OpenRequestsForObjID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.OpenRequestsForObjIDMock.mainExpectation != nil {

		input := m.OpenRequestsForObjIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingAccessorMockOpenRequestsForObjIDInput{p, p1, p2, p3}, "PendingAccessor.OpenRequestsForObjID got unexpected parameters")
		}

		result := m.OpenRequestsForObjIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.OpenRequestsForObjID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.OpenRequestsForObjIDFunc == nil {
		m.t.Fatalf("Unexpected call to PendingAccessorMock.OpenRequestsForObjID. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.OpenRequestsForObjIDFunc(p, p1, p2, p3)
}

// OpenRequestsForObjIDMinimockCounter returns a count of PendingAccessorMock.OpenRequestsForObjIDFunc invocations
func (m *PendingAccessorMock) OpenRequestsForObjIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OpenRequestsForObjIDCounter)
}

// OpenRequestsForObjIDMinimockPreCounter returns the value of PendingAccessorMock.OpenRequestsForObjID invocations
func (m *PendingAccessorMock) OpenRequestsForObjIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OpenRequestsForObjIDPreCounter)
}

// OpenRequestsForObjIDFinished returns true if mock invocations count is ok
func (m *PendingAccessorMock) OpenRequestsForObjIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OpenRequestsForObjIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OpenRequestsForObjIDCounter) == uint64(len(m.OpenRequestsForObjIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OpenRequestsForObjIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OpenRequestsForObjIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OpenRequestsForObjIDFunc != nil {
		return atomic.LoadUint64(&m.OpenRequestsForObjIDCounter) > 0
	}

	return true
}

type mPendingAccessorMockRecords struct {
	mock              *PendingAccessorMock
	mainExpectation   *PendingAccessorMockRecordsExpectation
	expectationSeries []*PendingAccessorMockRecordsExpectation
}

type PendingAccessorMockRecordsExpectation struct {
	input  *PendingAccessorMockRecordsInput
	result *PendingAccessorMockRecordsResult
}

type PendingAccessorMockRecordsInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type PendingAccessorMockRecordsResult struct {
	r  []record.Virtual
	r1 error
}

// Expect specifies that invocation of PendingAccessor.Records is expected from 1 to Infinity times
func (m *mPendingAccessorMockRecords) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mPendingAccessorMockRecords {
	m.mock.RecordsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockRecordsExpectation{}
	}
	m.mainExpectation.input = &PendingAccessorMockRecordsInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of PendingAccessor.Records
func (m *mPendingAccessorMockRecords) Return(r []record.Virtual, r1 error) *PendingAccessorMock {
	m.mock.RecordsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockRecordsExpectation{}
	}
	m.mainExpectation.result = &PendingAccessorMockRecordsResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of PendingAccessor.Records is expected once
func (m *mPendingAccessorMockRecords) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *PendingAccessorMockRecordsExpectation {
	m.mock.RecordsFunc = nil
	m.mainExpectation = nil

	expectation := &PendingAccessorMockRecordsExpectation{}
	expectation.input = &PendingAccessorMockRecordsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingAccessorMockRecordsExpectation) Return(r []record.Virtual, r1 error) {
	e.result = &PendingAccessorMockRecordsResult{r, r1}
}

// Set uses given function f as a mock of PendingAccessor.Records method
func (m *mPendingAccessorMockRecords) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.Virtual, r1 error)) *PendingAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RecordsFunc = f
	return m.mock
}

// Records implements github.com/insolar/insolar/ledger/object.PendingAccessor interface
func (m *PendingAccessorMock) Records(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.Virtual, r1 error) {
	counter := atomic.AddUint64(&m.RecordsPreCounter, 1)
	defer atomic.AddUint64(&m.RecordsCounter, 1)

	if len(m.RecordsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RecordsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingAccessorMock.Records. %v %v %v", p, p1, p2)
			return
		}

		input := m.RecordsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingAccessorMockRecordsInput{p, p1, p2}, "PendingAccessor.Records got unexpected parameters")

		result := m.RecordsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.Records")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RecordsMock.mainExpectation != nil {

		input := m.RecordsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingAccessorMockRecordsInput{p, p1, p2}, "PendingAccessor.Records got unexpected parameters")
		}

		result := m.RecordsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.Records")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RecordsFunc == nil {
		m.t.Fatalf("Unexpected call to PendingAccessorMock.Records. %v %v %v", p, p1, p2)
		return
	}

	return m.RecordsFunc(p, p1, p2)
}

// RecordsMinimockCounter returns a count of PendingAccessorMock.RecordsFunc invocations
func (m *PendingAccessorMock) RecordsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RecordsCounter)
}

// RecordsMinimockPreCounter returns the value of PendingAccessorMock.Records invocations
func (m *PendingAccessorMock) RecordsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RecordsPreCounter)
}

// RecordsFinished returns true if mock invocations count is ok
func (m *PendingAccessorMock) RecordsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RecordsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RecordsCounter) == uint64(len(m.RecordsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RecordsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RecordsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RecordsFunc != nil {
		return atomic.LoadUint64(&m.RecordsCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingAccessorMock) ValidateCallCounters() {

	if !m.MetaForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.MetaForObjID")
	}

	if !m.OpenRequestsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsForObjID")
	}

	if !m.RecordsFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.Records")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PendingAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PendingAccessorMock) MinimockFinish() {

	if !m.MetaForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.MetaForObjID")
	}

	if !m.OpenRequestsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsForObjID")
	}

	if !m.RecordsFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.Records")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PendingAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *PendingAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.MetaForObjIDFinished()
		ok = ok && m.OpenRequestsForObjIDFinished()
		ok = ok && m.RecordsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.MetaForObjIDFinished() {
				m.t.Error("Expected call to PendingAccessorMock.MetaForObjID")
			}

			if !m.OpenRequestsForObjIDFinished() {
				m.t.Error("Expected call to PendingAccessorMock.OpenRequestsForObjID")
			}

			if !m.RecordsFinished() {
				m.t.Error("Expected call to PendingAccessorMock.Records")
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
func (m *PendingAccessorMock) AllMocksCalled() bool {

	if !m.MetaForObjIDFinished() {
		return false
	}

	if !m.OpenRequestsForObjIDFinished() {
		return false
	}

	if !m.RecordsFinished() {
		return false
	}

	return true
}

package object

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

	testify_assert "github.com/stretchr/testify/assert"
)

//PendingAccessorMock implements github.com/insolar/insolar/ledger/object.PendingAccessor
type PendingAccessorMock struct {
	t minimock.Tester

	OpenRequestsForObjIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []record.Request, r1 error)
	OpenRequestsForObjIDCounter    uint64
	OpenRequestsForObjIDPreCounter uint64
	OpenRequestsForObjIDMock       mPendingAccessorMockOpenRequestsForObjID

	OpenRequestsIDsForObjIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []insolar.ID, r1 error)
	OpenRequestsIDsForObjIDCounter    uint64
	OpenRequestsIDsForObjIDPreCounter uint64
	OpenRequestsIDsForObjIDMock       mPendingAccessorMockOpenRequestsIDsForObjID

	RecordsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.CompositeFilamentRecord, r1 error)
	RecordsCounter    uint64
	RecordsPreCounter uint64
	RecordsMock       mPendingAccessorMockRecords
}

//NewPendingAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.PendingAccessor
func NewPendingAccessorMock(t minimock.Tester) *PendingAccessorMock {
	m := &PendingAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.OpenRequestsForObjIDMock = mPendingAccessorMockOpenRequestsForObjID{mock: m}
	m.OpenRequestsIDsForObjIDMock = mPendingAccessorMockOpenRequestsIDsForObjID{mock: m}
	m.RecordsMock = mPendingAccessorMockRecords{mock: m}

	return m
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

//Expect specifies that invocation of PendingAccessor.OpenRequestsForObjID is expected from 1 to Infinity times
func (m *mPendingAccessorMockOpenRequestsForObjID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) *mPendingAccessorMockOpenRequestsForObjID {
	m.mock.OpenRequestsForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockOpenRequestsForObjIDExpectation{}
	}
	m.mainExpectation.input = &PendingAccessorMockOpenRequestsForObjIDInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of PendingAccessor.OpenRequestsForObjID
func (m *mPendingAccessorMockOpenRequestsForObjID) Return(r []record.Request, r1 error) *PendingAccessorMock {
	m.mock.OpenRequestsForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockOpenRequestsForObjIDExpectation{}
	}
	m.mainExpectation.result = &PendingAccessorMockOpenRequestsForObjIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingAccessor.OpenRequestsForObjID is expected once
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

//Set uses given function f as a mock of PendingAccessor.OpenRequestsForObjID method
func (m *mPendingAccessorMockOpenRequestsForObjID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []record.Request, r1 error)) *PendingAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OpenRequestsForObjIDFunc = f
	return m.mock
}

//OpenRequestsForObjID implements github.com/insolar/insolar/ledger/object.PendingAccessor interface
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

//OpenRequestsForObjIDMinimockCounter returns a count of PendingAccessorMock.OpenRequestsForObjIDFunc invocations
func (m *PendingAccessorMock) OpenRequestsForObjIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OpenRequestsForObjIDCounter)
}

//OpenRequestsForObjIDMinimockPreCounter returns the value of PendingAccessorMock.OpenRequestsForObjID invocations
func (m *PendingAccessorMock) OpenRequestsForObjIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OpenRequestsForObjIDPreCounter)
}

//OpenRequestsForObjIDFinished returns true if mock invocations count is ok
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

type mPendingAccessorMockOpenRequestsIDsForObjID struct {
	mock              *PendingAccessorMock
	mainExpectation   *PendingAccessorMockOpenRequestsIDsForObjIDExpectation
	expectationSeries []*PendingAccessorMockOpenRequestsIDsForObjIDExpectation
}

type PendingAccessorMockOpenRequestsIDsForObjIDExpectation struct {
	input  *PendingAccessorMockOpenRequestsIDsForObjIDInput
	result *PendingAccessorMockOpenRequestsIDsForObjIDResult
}

type PendingAccessorMockOpenRequestsIDsForObjIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 int
}

type PendingAccessorMockOpenRequestsIDsForObjIDResult struct {
	r  []insolar.ID
	r1 error
}

// Expect specifies that invocation of PendingAccessor.OpenRequestsIDsForObjID is expected from 1 to Infinity times
func (m *mPendingAccessorMockOpenRequestsIDsForObjID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) *mPendingAccessorMockOpenRequestsIDsForObjID {
	m.mock.OpenRequestsIDsForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockOpenRequestsIDsForObjIDExpectation{}
	}
	m.mainExpectation.input = &PendingAccessorMockOpenRequestsIDsForObjIDInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of PendingAccessor.OpenRequestsIDsForObjID
func (m *mPendingAccessorMockOpenRequestsIDsForObjID) Return(r []insolar.ID, r1 error) *PendingAccessorMock {
	m.mock.OpenRequestsIDsForObjIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockOpenRequestsIDsForObjIDExpectation{}
	}
	m.mainExpectation.result = &PendingAccessorMockOpenRequestsIDsForObjIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of PendingAccessor.OpenRequestsIDsForObjID is expected once
func (m *mPendingAccessorMockOpenRequestsIDsForObjID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) *PendingAccessorMockOpenRequestsIDsForObjIDExpectation {
	m.mock.OpenRequestsIDsForObjIDFunc = nil
	m.mainExpectation = nil

	expectation := &PendingAccessorMockOpenRequestsIDsForObjIDExpectation{}
	expectation.input = &PendingAccessorMockOpenRequestsIDsForObjIDInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingAccessorMockOpenRequestsIDsForObjIDExpectation) Return(r []insolar.ID, r1 error) {
	e.result = &PendingAccessorMockOpenRequestsIDsForObjIDResult{r, r1}
}

// Set uses given function f as a mock of PendingAccessor.OpenRequestsIDsForObjID method
func (m *mPendingAccessorMockOpenRequestsIDsForObjID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []insolar.ID, r1 error)) *PendingAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OpenRequestsIDsForObjIDFunc = f
	return m.mock
}

// OpenRequestsIDsForObjID implements github.com/insolar/insolar/ledger/object.PendingAccessor interface
func (m *PendingAccessorMock) OpenRequestsIDsForObjID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.OpenRequestsIDsForObjIDPreCounter, 1)
	defer atomic.AddUint64(&m.OpenRequestsIDsForObjIDCounter, 1)

	if len(m.OpenRequestsIDsForObjIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OpenRequestsIDsForObjIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingAccessorMock.OpenRequestsIDsForObjID. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.OpenRequestsIDsForObjIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingAccessorMockOpenRequestsIDsForObjIDInput{p, p1, p2, p3}, "PendingAccessor.OpenRequestsIDsForObjID got unexpected parameters")

		result := m.OpenRequestsIDsForObjIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.OpenRequestsIDsForObjID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.OpenRequestsIDsForObjIDMock.mainExpectation != nil {

		input := m.OpenRequestsIDsForObjIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingAccessorMockOpenRequestsIDsForObjIDInput{p, p1, p2, p3}, "PendingAccessor.OpenRequestsIDsForObjID got unexpected parameters")
		}

		result := m.OpenRequestsIDsForObjIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.OpenRequestsIDsForObjID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.OpenRequestsIDsForObjIDFunc == nil {
		m.t.Fatalf("Unexpected call to PendingAccessorMock.OpenRequestsIDsForObjID. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.OpenRequestsIDsForObjIDFunc(p, p1, p2, p3)
}

// OpenRequestsIDsForObjIDMinimockCounter returns a count of PendingAccessorMock.OpenRequestsIDsForObjIDFunc invocations
func (m *PendingAccessorMock) OpenRequestsIDsForObjIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OpenRequestsIDsForObjIDCounter)
}

// OpenRequestsIDsForObjIDMinimockPreCounter returns the value of PendingAccessorMock.OpenRequestsIDsForObjID invocations
func (m *PendingAccessorMock) OpenRequestsIDsForObjIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OpenRequestsIDsForObjIDPreCounter)
}

// OpenRequestsIDsForObjIDFinished returns true if mock invocations count is ok
func (m *PendingAccessorMock) OpenRequestsIDsForObjIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OpenRequestsIDsForObjIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OpenRequestsIDsForObjIDCounter) == uint64(len(m.OpenRequestsIDsForObjIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OpenRequestsIDsForObjIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OpenRequestsIDsForObjIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OpenRequestsIDsForObjIDFunc != nil {
		return atomic.LoadUint64(&m.OpenRequestsIDsForObjIDCounter) > 0
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
	r  []record.CompositeFilamentRecord
	r1 error
}

//Expect specifies that invocation of PendingAccessor.Records is expected from 1 to Infinity times
func (m *mPendingAccessorMockRecords) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mPendingAccessorMockRecords {
	m.mock.RecordsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockRecordsExpectation{}
	}
	m.mainExpectation.input = &PendingAccessorMockRecordsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PendingAccessor.Records
func (m *mPendingAccessorMockRecords) Return(r []record.CompositeFilamentRecord, r1 error) *PendingAccessorMock {
	m.mock.RecordsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockRecordsExpectation{}
	}
	m.mainExpectation.result = &PendingAccessorMockRecordsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingAccessor.Records is expected once
func (m *mPendingAccessorMockRecords) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *PendingAccessorMockRecordsExpectation {
	m.mock.RecordsFunc = nil
	m.mainExpectation = nil

	expectation := &PendingAccessorMockRecordsExpectation{}
	expectation.input = &PendingAccessorMockRecordsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingAccessorMockRecordsExpectation) Return(r []record.CompositeFilamentRecord, r1 error) {
	e.result = &PendingAccessorMockRecordsResult{r, r1}
}

//Set uses given function f as a mock of PendingAccessor.Records method
func (m *mPendingAccessorMockRecords) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.CompositeFilamentRecord, r1 error)) *PendingAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RecordsFunc = f
	return m.mock
}

//Records implements github.com/insolar/insolar/ledger/object.PendingAccessor interface
func (m *PendingAccessorMock) Records(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r []record.CompositeFilamentRecord, r1 error) {
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

//RecordsMinimockCounter returns a count of PendingAccessorMock.RecordsFunc invocations
func (m *PendingAccessorMock) RecordsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RecordsCounter)
}

//RecordsMinimockPreCounter returns the value of PendingAccessorMock.Records invocations
func (m *PendingAccessorMock) RecordsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RecordsPreCounter)
}

//RecordsFinished returns true if mock invocations count is ok
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingAccessorMock) ValidateCallCounters() {

	if !m.OpenRequestsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsForObjID")
	}

	if !m.OpenRequestsIDsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsIDsForObjID")
	}

	if !m.RecordsFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.Records")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PendingAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PendingAccessorMock) MinimockFinish() {

	if !m.OpenRequestsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsForObjID")
	}

	if !m.OpenRequestsIDsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsIDsForObjID")
	}

	if !m.RecordsFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.Records")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PendingAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PendingAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.OpenRequestsForObjIDFinished()
		ok = ok && m.OpenRequestsIDsForObjIDFinished()
		ok = ok && m.RecordsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.OpenRequestsForObjIDFinished() {
				m.t.Error("Expected call to PendingAccessorMock.OpenRequestsForObjID")
			}

			if !m.OpenRequestsIDsForObjIDFinished() {
				m.t.Error("Expected call to PendingAccessorMock.OpenRequestsIDsForObjID")
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

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *PendingAccessorMock) AllMocksCalled() bool {

	if !m.OpenRequestsForObjIDFinished() {
		return false
	}

	if !m.OpenRequestsIDsForObjIDFinished() {
		return false
	}

	if !m.RecordsFinished() {
		return false
	}

	return true
}

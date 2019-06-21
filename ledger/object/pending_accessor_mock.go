package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PendingAccessor" can be found in github.com/insolar/insolar/ledger/object
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

//PendingAccessorMock implements github.com/insolar/insolar/ledger/object.PendingAccessor
type PendingAccessorMock struct {
	t minimock.Tester

	FirstPendingFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r *record.PendingFilament, r1 error)
	FirstPendingCounter    uint64
	FirstPendingPreCounter uint64
	FirstPendingMock       mPendingAccessorMockFirstPending

	OpenRequestsForObjIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 int) (r []record.Request, r1 error)
	OpenRequestsForObjIDCounter    uint64
	OpenRequestsForObjIDPreCounter uint64
	OpenRequestsForObjIDMock       mPendingAccessorMockOpenRequestsForObjID

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

	m.FirstPendingMock = mPendingAccessorMockFirstPending{mock: m}
	m.OpenRequestsForObjIDMock = mPendingAccessorMockOpenRequestsForObjID{mock: m}
	m.RecordsMock = mPendingAccessorMockRecords{mock: m}

	return m
}

type mPendingAccessorMockFirstPending struct {
	mock              *PendingAccessorMock
	mainExpectation   *PendingAccessorMockFirstPendingExpectation
	expectationSeries []*PendingAccessorMockFirstPendingExpectation
}

type PendingAccessorMockFirstPendingExpectation struct {
	input  *PendingAccessorMockFirstPendingInput
	result *PendingAccessorMockFirstPendingResult
}

type PendingAccessorMockFirstPendingInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type PendingAccessorMockFirstPendingResult struct {
	r  *record.PendingFilament
	r1 error
}

//Expect specifies that invocation of PendingAccessor.FirstPending is expected from 1 to Infinity times
func (m *mPendingAccessorMockFirstPending) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mPendingAccessorMockFirstPending {
	m.mock.FirstPendingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockFirstPendingExpectation{}
	}
	m.mainExpectation.input = &PendingAccessorMockFirstPendingInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PendingAccessor.FirstPending
func (m *mPendingAccessorMockFirstPending) Return(r *record.PendingFilament, r1 error) *PendingAccessorMock {
	m.mock.FirstPendingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingAccessorMockFirstPendingExpectation{}
	}
	m.mainExpectation.result = &PendingAccessorMockFirstPendingResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingAccessor.FirstPending is expected once
func (m *mPendingAccessorMockFirstPending) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *PendingAccessorMockFirstPendingExpectation {
	m.mock.FirstPendingFunc = nil
	m.mainExpectation = nil

	expectation := &PendingAccessorMockFirstPendingExpectation{}
	expectation.input = &PendingAccessorMockFirstPendingInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingAccessorMockFirstPendingExpectation) Return(r *record.PendingFilament, r1 error) {
	e.result = &PendingAccessorMockFirstPendingResult{r, r1}
}

//Set uses given function f as a mock of PendingAccessor.FirstPending method
func (m *mPendingAccessorMockFirstPending) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r *record.PendingFilament, r1 error)) *PendingAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FirstPendingFunc = f
	return m.mock
}

//FirstPending implements github.com/insolar/insolar/ledger/object.PendingAccessor interface
func (m *PendingAccessorMock) FirstPending(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r *record.PendingFilament, r1 error) {
	counter := atomic.AddUint64(&m.FirstPendingPreCounter, 1)
	defer atomic.AddUint64(&m.FirstPendingCounter, 1)

	if len(m.FirstPendingMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FirstPendingMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingAccessorMock.FirstPending. %v %v %v", p, p1, p2)
			return
		}

		input := m.FirstPendingMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingAccessorMockFirstPendingInput{p, p1, p2}, "PendingAccessor.FirstPending got unexpected parameters")

		result := m.FirstPendingMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.FirstPending")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FirstPendingMock.mainExpectation != nil {

		input := m.FirstPendingMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingAccessorMockFirstPendingInput{p, p1, p2}, "PendingAccessor.FirstPending got unexpected parameters")
		}

		result := m.FirstPendingMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingAccessorMock.FirstPending")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FirstPendingFunc == nil {
		m.t.Fatalf("Unexpected call to PendingAccessorMock.FirstPending. %v %v %v", p, p1, p2)
		return
	}

	return m.FirstPendingFunc(p, p1, p2)
}

//FirstPendingMinimockCounter returns a count of PendingAccessorMock.FirstPendingFunc invocations
func (m *PendingAccessorMock) FirstPendingMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FirstPendingCounter)
}

//FirstPendingMinimockPreCounter returns the value of PendingAccessorMock.FirstPending invocations
func (m *PendingAccessorMock) FirstPendingMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FirstPendingPreCounter)
}

//FirstPendingFinished returns true if mock invocations count is ok
func (m *PendingAccessorMock) FirstPendingFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FirstPendingMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FirstPendingCounter) == uint64(len(m.FirstPendingMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FirstPendingMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FirstPendingCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FirstPendingFunc != nil {
		return atomic.LoadUint64(&m.FirstPendingCounter) > 0
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

	if !m.FirstPendingFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.FirstPending")
	}

	if !m.OpenRequestsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsForObjID")
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

	if !m.FirstPendingFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.FirstPending")
	}

	if !m.OpenRequestsForObjIDFinished() {
		m.t.Fatal("Expected call to PendingAccessorMock.OpenRequestsForObjID")
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
		ok = ok && m.FirstPendingFinished()
		ok = ok && m.OpenRequestsForObjIDFinished()
		ok = ok && m.RecordsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FirstPendingFinished() {
				m.t.Error("Expected call to PendingAccessorMock.FirstPending")
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

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *PendingAccessorMock) AllMocksCalled() bool {

	if !m.FirstPendingFinished() {
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

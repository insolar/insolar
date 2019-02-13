package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Cleaner" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	recentstorage "github.com/insolar/insolar/ledger/recentstorage"

	testify_assert "github.com/stretchr/testify/assert"
)

//CleanerMock implements github.com/insolar/insolar/ledger/storage.Cleaner
type CleanerMock struct {
	t minimock.Tester

	CleanJetIndexesFunc       func(p context.Context, p1 core.RecordID, p2 recentstorage.RecentIndexStorage, p3 []core.RecordID) (r RmStat, r1 error)
	CleanJetIndexesCounter    uint64
	CleanJetIndexesPreCounter uint64
	CleanJetIndexesMock       mCleanerMockCleanJetIndexes

	CleanJetRecordsUntilPulseFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r map[string]RmStat, r1 error)
	CleanJetRecordsUntilPulseCounter    uint64
	CleanJetRecordsUntilPulsePreCounter uint64
	CleanJetRecordsUntilPulseMock       mCleanerMockCleanJetRecordsUntilPulse
}

//NewCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage.Cleaner
func NewCleanerMock(t minimock.Tester) *CleanerMock {
	m := &CleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CleanJetIndexesMock = mCleanerMockCleanJetIndexes{mock: m}
	m.CleanJetRecordsUntilPulseMock = mCleanerMockCleanJetRecordsUntilPulse{mock: m}

	return m
}

type mCleanerMockCleanJetIndexes struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockCleanJetIndexesExpectation
	expectationSeries []*CleanerMockCleanJetIndexesExpectation
}

type CleanerMockCleanJetIndexesExpectation struct {
	input  *CleanerMockCleanJetIndexesInput
	result *CleanerMockCleanJetIndexesResult
}

type CleanerMockCleanJetIndexesInput struct {
	p  context.Context
	p1 core.RecordID
	p2 recentstorage.RecentIndexStorage
	p3 []core.RecordID
}

type CleanerMockCleanJetIndexesResult struct {
	r  RmStat
	r1 error
}

//Expect specifies that invocation of Cleaner.CleanJetIndexes is expected from 1 to Infinity times
func (m *mCleanerMockCleanJetIndexes) Expect(p context.Context, p1 core.RecordID, p2 recentstorage.RecentIndexStorage, p3 []core.RecordID) *mCleanerMockCleanJetIndexes {
	m.mock.CleanJetIndexesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockCleanJetIndexesExpectation{}
	}
	m.mainExpectation.input = &CleanerMockCleanJetIndexesInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Cleaner.CleanJetIndexes
func (m *mCleanerMockCleanJetIndexes) Return(r RmStat, r1 error) *CleanerMock {
	m.mock.CleanJetIndexesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockCleanJetIndexesExpectation{}
	}
	m.mainExpectation.result = &CleanerMockCleanJetIndexesResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Cleaner.CleanJetIndexes is expected once
func (m *mCleanerMockCleanJetIndexes) ExpectOnce(p context.Context, p1 core.RecordID, p2 recentstorage.RecentIndexStorage, p3 []core.RecordID) *CleanerMockCleanJetIndexesExpectation {
	m.mock.CleanJetIndexesFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockCleanJetIndexesExpectation{}
	expectation.input = &CleanerMockCleanJetIndexesInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CleanerMockCleanJetIndexesExpectation) Return(r RmStat, r1 error) {
	e.result = &CleanerMockCleanJetIndexesResult{r, r1}
}

//Set uses given function f as a mock of Cleaner.CleanJetIndexes method
func (m *mCleanerMockCleanJetIndexes) Set(f func(p context.Context, p1 core.RecordID, p2 recentstorage.RecentIndexStorage, p3 []core.RecordID) (r RmStat, r1 error)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CleanJetIndexesFunc = f
	return m.mock
}

//CleanJetIndexes implements github.com/insolar/insolar/ledger/storage.Cleaner interface
func (m *CleanerMock) CleanJetIndexes(p context.Context, p1 core.RecordID, p2 recentstorage.RecentIndexStorage, p3 []core.RecordID) (r RmStat, r1 error) {
	counter := atomic.AddUint64(&m.CleanJetIndexesPreCounter, 1)
	defer atomic.AddUint64(&m.CleanJetIndexesCounter, 1)

	if len(m.CleanJetIndexesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CleanJetIndexesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.CleanJetIndexes. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.CleanJetIndexesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockCleanJetIndexesInput{p, p1, p2, p3}, "Cleaner.CleanJetIndexes got unexpected parameters")

		result := m.CleanJetIndexesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.CleanJetIndexes")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CleanJetIndexesMock.mainExpectation != nil {

		input := m.CleanJetIndexesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockCleanJetIndexesInput{p, p1, p2, p3}, "Cleaner.CleanJetIndexes got unexpected parameters")
		}

		result := m.CleanJetIndexesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.CleanJetIndexes")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CleanJetIndexesFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.CleanJetIndexes. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.CleanJetIndexesFunc(p, p1, p2, p3)
}

//CleanJetIndexesMinimockCounter returns a count of CleanerMock.CleanJetIndexesFunc invocations
func (m *CleanerMock) CleanJetIndexesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CleanJetIndexesCounter)
}

//CleanJetIndexesMinimockPreCounter returns the value of CleanerMock.CleanJetIndexes invocations
func (m *CleanerMock) CleanJetIndexesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CleanJetIndexesPreCounter)
}

//CleanJetIndexesFinished returns true if mock invocations count is ok
func (m *CleanerMock) CleanJetIndexesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CleanJetIndexesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CleanJetIndexesCounter) == uint64(len(m.CleanJetIndexesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CleanJetIndexesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CleanJetIndexesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CleanJetIndexesFunc != nil {
		return atomic.LoadUint64(&m.CleanJetIndexesCounter) > 0
	}

	return true
}

type mCleanerMockCleanJetRecordsUntilPulse struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockCleanJetRecordsUntilPulseExpectation
	expectationSeries []*CleanerMockCleanJetRecordsUntilPulseExpectation
}

type CleanerMockCleanJetRecordsUntilPulseExpectation struct {
	input  *CleanerMockCleanJetRecordsUntilPulseInput
	result *CleanerMockCleanJetRecordsUntilPulseResult
}

type CleanerMockCleanJetRecordsUntilPulseInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

type CleanerMockCleanJetRecordsUntilPulseResult struct {
	r  map[string]RmStat
	r1 error
}

//Expect specifies that invocation of Cleaner.CleanJetRecordsUntilPulse is expected from 1 to Infinity times
func (m *mCleanerMockCleanJetRecordsUntilPulse) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mCleanerMockCleanJetRecordsUntilPulse {
	m.mock.CleanJetRecordsUntilPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockCleanJetRecordsUntilPulseExpectation{}
	}
	m.mainExpectation.input = &CleanerMockCleanJetRecordsUntilPulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Cleaner.CleanJetRecordsUntilPulse
func (m *mCleanerMockCleanJetRecordsUntilPulse) Return(r map[string]RmStat, r1 error) *CleanerMock {
	m.mock.CleanJetRecordsUntilPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockCleanJetRecordsUntilPulseExpectation{}
	}
	m.mainExpectation.result = &CleanerMockCleanJetRecordsUntilPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Cleaner.CleanJetRecordsUntilPulse is expected once
func (m *mCleanerMockCleanJetRecordsUntilPulse) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *CleanerMockCleanJetRecordsUntilPulseExpectation {
	m.mock.CleanJetRecordsUntilPulseFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockCleanJetRecordsUntilPulseExpectation{}
	expectation.input = &CleanerMockCleanJetRecordsUntilPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CleanerMockCleanJetRecordsUntilPulseExpectation) Return(r map[string]RmStat, r1 error) {
	e.result = &CleanerMockCleanJetRecordsUntilPulseResult{r, r1}
}

//Set uses given function f as a mock of Cleaner.CleanJetRecordsUntilPulse method
func (m *mCleanerMockCleanJetRecordsUntilPulse) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r map[string]RmStat, r1 error)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CleanJetRecordsUntilPulseFunc = f
	return m.mock
}

//CleanJetRecordsUntilPulse implements github.com/insolar/insolar/ledger/storage.Cleaner interface
func (m *CleanerMock) CleanJetRecordsUntilPulse(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r map[string]RmStat, r1 error) {
	counter := atomic.AddUint64(&m.CleanJetRecordsUntilPulsePreCounter, 1)
	defer atomic.AddUint64(&m.CleanJetRecordsUntilPulseCounter, 1)

	if len(m.CleanJetRecordsUntilPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CleanJetRecordsUntilPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.CleanJetRecordsUntilPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.CleanJetRecordsUntilPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockCleanJetRecordsUntilPulseInput{p, p1, p2}, "Cleaner.CleanJetRecordsUntilPulse got unexpected parameters")

		result := m.CleanJetRecordsUntilPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.CleanJetRecordsUntilPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CleanJetRecordsUntilPulseMock.mainExpectation != nil {

		input := m.CleanJetRecordsUntilPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockCleanJetRecordsUntilPulseInput{p, p1, p2}, "Cleaner.CleanJetRecordsUntilPulse got unexpected parameters")
		}

		result := m.CleanJetRecordsUntilPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.CleanJetRecordsUntilPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CleanJetRecordsUntilPulseFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.CleanJetRecordsUntilPulse. %v %v %v", p, p1, p2)
		return
	}

	return m.CleanJetRecordsUntilPulseFunc(p, p1, p2)
}

//CleanJetRecordsUntilPulseMinimockCounter returns a count of CleanerMock.CleanJetRecordsUntilPulseFunc invocations
func (m *CleanerMock) CleanJetRecordsUntilPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CleanJetRecordsUntilPulseCounter)
}

//CleanJetRecordsUntilPulseMinimockPreCounter returns the value of CleanerMock.CleanJetRecordsUntilPulse invocations
func (m *CleanerMock) CleanJetRecordsUntilPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CleanJetRecordsUntilPulsePreCounter)
}

//CleanJetRecordsUntilPulseFinished returns true if mock invocations count is ok
func (m *CleanerMock) CleanJetRecordsUntilPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CleanJetRecordsUntilPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CleanJetRecordsUntilPulseCounter) == uint64(len(m.CleanJetRecordsUntilPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CleanJetRecordsUntilPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CleanJetRecordsUntilPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CleanJetRecordsUntilPulseFunc != nil {
		return atomic.LoadUint64(&m.CleanJetRecordsUntilPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CleanerMock) ValidateCallCounters() {

	if !m.CleanJetIndexesFinished() {
		m.t.Fatal("Expected call to CleanerMock.CleanJetIndexes")
	}

	if !m.CleanJetRecordsUntilPulseFinished() {
		m.t.Fatal("Expected call to CleanerMock.CleanJetRecordsUntilPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CleanerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CleanerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CleanerMock) MinimockFinish() {

	if !m.CleanJetIndexesFinished() {
		m.t.Fatal("Expected call to CleanerMock.CleanJetIndexes")
	}

	if !m.CleanJetRecordsUntilPulseFinished() {
		m.t.Fatal("Expected call to CleanerMock.CleanJetRecordsUntilPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CleanJetIndexesFinished()
		ok = ok && m.CleanJetRecordsUntilPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CleanJetIndexesFinished() {
				m.t.Error("Expected call to CleanerMock.CleanJetIndexes")
			}

			if !m.CleanJetRecordsUntilPulseFinished() {
				m.t.Error("Expected call to CleanerMock.CleanJetRecordsUntilPulse")
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
func (m *CleanerMock) AllMocksCalled() bool {

	if !m.CleanJetIndexesFinished() {
		return false
	}

	if !m.CleanJetRecordsUntilPulseFinished() {
		return false
	}

	return true
}

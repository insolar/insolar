package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Cleaner" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/recentstorage"

	testify_assert "github.com/stretchr/testify/assert"
)

// CleanerMock implements github.com/insolar/insolar/ledger/storage.Cleaner
type CleanerMock struct {
	t minimock.Tester

	RemoveAllForJetUntilPulseFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r map[string]int, r1 error)
	RemoveAllForJetUntilPulseCounter    uint64
	RemoveAllForJetUntilPulsePreCounter uint64
	RemoveAllForJetUntilPulseMock       mCleanerMockRemoveAllForJetUntilPulse

	RemoveJetBlobsUntilFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r int, r1 error)
	RemoveJetBlobsUntilCounter    uint64
	RemoveJetBlobsUntilPreCounter uint64
	RemoveJetBlobsUntilMock       mCleanerMockRemoveJetBlobsUntil

	RemoveJetDropsUntilFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r int, r1 error)
	RemoveJetDropsUntilCounter    uint64
	RemoveJetDropsUntilPreCounter uint64
	RemoveJetDropsUntilMock       mCleanerMockRemoveJetDropsUntil

	RemoveJetIndexesUntilFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r int, r1 error)
	RemoveJetIndexesUntilCounter    uint64
	RemoveJetIndexesUntilPreCounter uint64
	RemoveJetIndexesUntilMock       mCleanerMockRemoveJetIndexesUntil

	RemoveJetRecordsUntilFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r int, r1 error)
	RemoveJetRecordsUntilCounter    uint64
	RemoveJetRecordsUntilPreCounter uint64
	RemoveJetRecordsUntilMock       mCleanerMockRemoveJetRecordsUntil
}

// NewCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage.Cleaner
func NewCleanerMock(t minimock.Tester) *CleanerMock {
	m := &CleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveAllForJetUntilPulseMock = mCleanerMockRemoveAllForJetUntilPulse{mock: m}
	m.RemoveJetBlobsUntilMock = mCleanerMockRemoveJetBlobsUntil{mock: m}
	m.RemoveJetDropsUntilMock = mCleanerMockRemoveJetDropsUntil{mock: m}
	m.RemoveJetIndexesUntilMock = mCleanerMockRemoveJetIndexesUntil{mock: m}
	m.RemoveJetRecordsUntilMock = mCleanerMockRemoveJetRecordsUntil{mock: m}

	return m
}

type mCleanerMockRemoveAllForJetUntilPulse struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockRemoveAllForJetUntilPulseExpectation
	expectationSeries []*CleanerMockRemoveAllForJetUntilPulseExpectation
}

type CleanerMockRemoveAllForJetUntilPulseExpectation struct {
	input  *CleanerMockRemoveAllForJetUntilPulseInput
	result *CleanerMockRemoveAllForJetUntilPulseResult
}

type CleanerMockRemoveAllForJetUntilPulseInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 recentstorage.RecentStorage
}

type CleanerMockRemoveAllForJetUntilPulseResult struct {
	r  map[string]int
	r1 error
}

// Expect specifies that invocation of Cleaner.RemoveAllForJetUntilPulse is expected from 1 to Infinity times
func (m *mCleanerMockRemoveAllForJetUntilPulse) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) *mCleanerMockRemoveAllForJetUntilPulse {
	m.mock.RemoveAllForJetUntilPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveAllForJetUntilPulseExpectation{}
	}
	m.mainExpectation.input = &CleanerMockRemoveAllForJetUntilPulseInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of Cleaner.RemoveAllForJetUntilPulse
func (m *mCleanerMockRemoveAllForJetUntilPulse) Return(r map[string]int, r1 error) *CleanerMock {
	m.mock.RemoveAllForJetUntilPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveAllForJetUntilPulseExpectation{}
	}
	m.mainExpectation.result = &CleanerMockRemoveAllForJetUntilPulseResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of Cleaner.RemoveAllForJetUntilPulse is expected once
func (m *mCleanerMockRemoveAllForJetUntilPulse) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) *CleanerMockRemoveAllForJetUntilPulseExpectation {
	m.mock.RemoveAllForJetUntilPulseFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockRemoveAllForJetUntilPulseExpectation{}
	expectation.input = &CleanerMockRemoveAllForJetUntilPulseInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CleanerMockRemoveAllForJetUntilPulseExpectation) Return(r map[string]int, r1 error) {
	e.result = &CleanerMockRemoveAllForJetUntilPulseResult{r, r1}
}

// Set uses given function f as a mock of Cleaner.RemoveAllForJetUntilPulse method
func (m *mCleanerMockRemoveAllForJetUntilPulse) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r map[string]int, r1 error)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveAllForJetUntilPulseFunc = f
	return m.mock
}

// RemoveAllForJetUntilPulse implements github.com/insolar/insolar/ledger/storage.Cleaner interface
func (m *CleanerMock) RemoveAllForJetUntilPulse(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r map[string]int, r1 error) {
	counter := atomic.AddUint64(&m.RemoveAllForJetUntilPulsePreCounter, 1)
	defer atomic.AddUint64(&m.RemoveAllForJetUntilPulseCounter, 1)

	if len(m.RemoveAllForJetUntilPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveAllForJetUntilPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.RemoveAllForJetUntilPulse. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RemoveAllForJetUntilPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockRemoveAllForJetUntilPulseInput{p, p1, p2, p3}, "Cleaner.RemoveAllForJetUntilPulse got unexpected parameters")

		result := m.RemoveAllForJetUntilPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveAllForJetUntilPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveAllForJetUntilPulseMock.mainExpectation != nil {

		input := m.RemoveAllForJetUntilPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockRemoveAllForJetUntilPulseInput{p, p1, p2, p3}, "Cleaner.RemoveAllForJetUntilPulse got unexpected parameters")
		}

		result := m.RemoveAllForJetUntilPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveAllForJetUntilPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveAllForJetUntilPulseFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.RemoveAllForJetUntilPulse. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RemoveAllForJetUntilPulseFunc(p, p1, p2, p3)
}

// RemoveAllForJetUntilPulseMinimockCounter returns a count of CleanerMock.RemoveAllForJetUntilPulseFunc invocations
func (m *CleanerMock) RemoveAllForJetUntilPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllForJetUntilPulseCounter)
}

// RemoveAllForJetUntilPulseMinimockPreCounter returns the value of CleanerMock.RemoveAllForJetUntilPulse invocations
func (m *CleanerMock) RemoveAllForJetUntilPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllForJetUntilPulsePreCounter)
}

// RemoveAllForJetUntilPulseFinished returns true if mock invocations count is ok
func (m *CleanerMock) RemoveAllForJetUntilPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveAllForJetUntilPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveAllForJetUntilPulseCounter) == uint64(len(m.RemoveAllForJetUntilPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveAllForJetUntilPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveAllForJetUntilPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveAllForJetUntilPulseFunc != nil {
		return atomic.LoadUint64(&m.RemoveAllForJetUntilPulseCounter) > 0
	}

	return true
}

type mCleanerMockRemoveJetBlobsUntil struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockRemoveJetBlobsUntilExpectation
	expectationSeries []*CleanerMockRemoveJetBlobsUntilExpectation
}

type CleanerMockRemoveJetBlobsUntilExpectation struct {
	input  *CleanerMockRemoveJetBlobsUntilInput
	result *CleanerMockRemoveJetBlobsUntilResult
}

type CleanerMockRemoveJetBlobsUntilInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

type CleanerMockRemoveJetBlobsUntilResult struct {
	r  int
	r1 error
}

// Expect specifies that invocation of Cleaner.RemoveJetBlobsUntil is expected from 1 to Infinity times
func (m *mCleanerMockRemoveJetBlobsUntil) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mCleanerMockRemoveJetBlobsUntil {
	m.mock.RemoveJetBlobsUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetBlobsUntilExpectation{}
	}
	m.mainExpectation.input = &CleanerMockRemoveJetBlobsUntilInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of Cleaner.RemoveJetBlobsUntil
func (m *mCleanerMockRemoveJetBlobsUntil) Return(r int, r1 error) *CleanerMock {
	m.mock.RemoveJetBlobsUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetBlobsUntilExpectation{}
	}
	m.mainExpectation.result = &CleanerMockRemoveJetBlobsUntilResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of Cleaner.RemoveJetBlobsUntil is expected once
func (m *mCleanerMockRemoveJetBlobsUntil) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *CleanerMockRemoveJetBlobsUntilExpectation {
	m.mock.RemoveJetBlobsUntilFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockRemoveJetBlobsUntilExpectation{}
	expectation.input = &CleanerMockRemoveJetBlobsUntilInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CleanerMockRemoveJetBlobsUntilExpectation) Return(r int, r1 error) {
	e.result = &CleanerMockRemoveJetBlobsUntilResult{r, r1}
}

// Set uses given function f as a mock of Cleaner.RemoveJetBlobsUntil method
func (m *mCleanerMockRemoveJetBlobsUntil) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r int, r1 error)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveJetBlobsUntilFunc = f
	return m.mock
}

// RemoveJetBlobsUntil implements github.com/insolar/insolar/ledger/storage.Cleaner interface
func (m *CleanerMock) RemoveJetBlobsUntil(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r int, r1 error) {
	counter := atomic.AddUint64(&m.RemoveJetBlobsUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveJetBlobsUntilCounter, 1)

	if len(m.RemoveJetBlobsUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveJetBlobsUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetBlobsUntil. %v %v %v", p, p1, p2)
			return
		}

		input := m.RemoveJetBlobsUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockRemoveJetBlobsUntilInput{p, p1, p2}, "Cleaner.RemoveJetBlobsUntil got unexpected parameters")

		result := m.RemoveJetBlobsUntilMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetBlobsUntil")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetBlobsUntilMock.mainExpectation != nil {

		input := m.RemoveJetBlobsUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockRemoveJetBlobsUntilInput{p, p1, p2}, "Cleaner.RemoveJetBlobsUntil got unexpected parameters")
		}

		result := m.RemoveJetBlobsUntilMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetBlobsUntil")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetBlobsUntilFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetBlobsUntil. %v %v %v", p, p1, p2)
		return
	}

	return m.RemoveJetBlobsUntilFunc(p, p1, p2)
}

// RemoveJetBlobsUntilMinimockCounter returns a count of CleanerMock.RemoveJetBlobsUntilFunc invocations
func (m *CleanerMock) RemoveJetBlobsUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetBlobsUntilCounter)
}

// RemoveJetBlobsUntilMinimockPreCounter returns the value of CleanerMock.RemoveJetBlobsUntil invocations
func (m *CleanerMock) RemoveJetBlobsUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetBlobsUntilPreCounter)
}

// RemoveJetBlobsUntilFinished returns true if mock invocations count is ok
func (m *CleanerMock) RemoveJetBlobsUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveJetBlobsUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveJetBlobsUntilCounter) == uint64(len(m.RemoveJetBlobsUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveJetBlobsUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveJetBlobsUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveJetBlobsUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveJetBlobsUntilCounter) > 0
	}

	return true
}

type mCleanerMockRemoveJetDropsUntil struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockRemoveJetDropsUntilExpectation
	expectationSeries []*CleanerMockRemoveJetDropsUntilExpectation
}

type CleanerMockRemoveJetDropsUntilExpectation struct {
	input  *CleanerMockRemoveJetDropsUntilInput
	result *CleanerMockRemoveJetDropsUntilResult
}

type CleanerMockRemoveJetDropsUntilInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

type CleanerMockRemoveJetDropsUntilResult struct {
	r  int
	r1 error
}

// Expect specifies that invocation of Cleaner.RemoveJetDropsUntil is expected from 1 to Infinity times
func (m *mCleanerMockRemoveJetDropsUntil) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mCleanerMockRemoveJetDropsUntil {
	m.mock.RemoveJetDropsUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetDropsUntilExpectation{}
	}
	m.mainExpectation.input = &CleanerMockRemoveJetDropsUntilInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of Cleaner.RemoveJetDropsUntil
func (m *mCleanerMockRemoveJetDropsUntil) Return(r int, r1 error) *CleanerMock {
	m.mock.RemoveJetDropsUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetDropsUntilExpectation{}
	}
	m.mainExpectation.result = &CleanerMockRemoveJetDropsUntilResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of Cleaner.RemoveJetDropsUntil is expected once
func (m *mCleanerMockRemoveJetDropsUntil) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *CleanerMockRemoveJetDropsUntilExpectation {
	m.mock.RemoveJetDropsUntilFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockRemoveJetDropsUntilExpectation{}
	expectation.input = &CleanerMockRemoveJetDropsUntilInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CleanerMockRemoveJetDropsUntilExpectation) Return(r int, r1 error) {
	e.result = &CleanerMockRemoveJetDropsUntilResult{r, r1}
}

// Set uses given function f as a mock of Cleaner.RemoveJetDropsUntil method
func (m *mCleanerMockRemoveJetDropsUntil) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r int, r1 error)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveJetDropsUntilFunc = f
	return m.mock
}

// RemoveJetDropsUntil implements github.com/insolar/insolar/ledger/storage.Cleaner interface
func (m *CleanerMock) RemoveJetDropsUntil(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r int, r1 error) {
	counter := atomic.AddUint64(&m.RemoveJetDropsUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveJetDropsUntilCounter, 1)

	if len(m.RemoveJetDropsUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveJetDropsUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetDropsUntil. %v %v %v", p, p1, p2)
			return
		}

		input := m.RemoveJetDropsUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockRemoveJetDropsUntilInput{p, p1, p2}, "Cleaner.RemoveJetDropsUntil got unexpected parameters")

		result := m.RemoveJetDropsUntilMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetDropsUntil")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetDropsUntilMock.mainExpectation != nil {

		input := m.RemoveJetDropsUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockRemoveJetDropsUntilInput{p, p1, p2}, "Cleaner.RemoveJetDropsUntil got unexpected parameters")
		}

		result := m.RemoveJetDropsUntilMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetDropsUntil")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetDropsUntilFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetDropsUntil. %v %v %v", p, p1, p2)
		return
	}

	return m.RemoveJetDropsUntilFunc(p, p1, p2)
}

// RemoveJetDropsUntilMinimockCounter returns a count of CleanerMock.RemoveJetDropsUntilFunc invocations
func (m *CleanerMock) RemoveJetDropsUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetDropsUntilCounter)
}

// RemoveJetDropsUntilMinimockPreCounter returns the value of CleanerMock.RemoveJetDropsUntil invocations
func (m *CleanerMock) RemoveJetDropsUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetDropsUntilPreCounter)
}

// RemoveJetDropsUntilFinished returns true if mock invocations count is ok
func (m *CleanerMock) RemoveJetDropsUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveJetDropsUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveJetDropsUntilCounter) == uint64(len(m.RemoveJetDropsUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveJetDropsUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveJetDropsUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveJetDropsUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveJetDropsUntilCounter) > 0
	}

	return true
}

type mCleanerMockRemoveJetIndexesUntil struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockRemoveJetIndexesUntilExpectation
	expectationSeries []*CleanerMockRemoveJetIndexesUntilExpectation
}

type CleanerMockRemoveJetIndexesUntilExpectation struct {
	input  *CleanerMockRemoveJetIndexesUntilInput
	result *CleanerMockRemoveJetIndexesUntilResult
}

type CleanerMockRemoveJetIndexesUntilInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 recentstorage.RecentStorage
}

type CleanerMockRemoveJetIndexesUntilResult struct {
	r  int
	r1 error
}

// Expect specifies that invocation of Cleaner.RemoveJetIndexesUntil is expected from 1 to Infinity times
func (m *mCleanerMockRemoveJetIndexesUntil) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) *mCleanerMockRemoveJetIndexesUntil {
	m.mock.RemoveJetIndexesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetIndexesUntilExpectation{}
	}
	m.mainExpectation.input = &CleanerMockRemoveJetIndexesUntilInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of Cleaner.RemoveJetIndexesUntil
func (m *mCleanerMockRemoveJetIndexesUntil) Return(r int, r1 error) *CleanerMock {
	m.mock.RemoveJetIndexesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetIndexesUntilExpectation{}
	}
	m.mainExpectation.result = &CleanerMockRemoveJetIndexesUntilResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of Cleaner.RemoveJetIndexesUntil is expected once
func (m *mCleanerMockRemoveJetIndexesUntil) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) *CleanerMockRemoveJetIndexesUntilExpectation {
	m.mock.RemoveJetIndexesUntilFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockRemoveJetIndexesUntilExpectation{}
	expectation.input = &CleanerMockRemoveJetIndexesUntilInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CleanerMockRemoveJetIndexesUntilExpectation) Return(r int, r1 error) {
	e.result = &CleanerMockRemoveJetIndexesUntilResult{r, r1}
}

// Set uses given function f as a mock of Cleaner.RemoveJetIndexesUntil method
func (m *mCleanerMockRemoveJetIndexesUntil) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r int, r1 error)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveJetIndexesUntilFunc = f
	return m.mock
}

// RemoveJetIndexesUntil implements github.com/insolar/insolar/ledger/storage.Cleaner interface
func (m *CleanerMock) RemoveJetIndexesUntil(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r int, r1 error) {
	counter := atomic.AddUint64(&m.RemoveJetIndexesUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveJetIndexesUntilCounter, 1)

	if len(m.RemoveJetIndexesUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveJetIndexesUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetIndexesUntil. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RemoveJetIndexesUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockRemoveJetIndexesUntilInput{p, p1, p2, p3}, "Cleaner.RemoveJetIndexesUntil got unexpected parameters")

		result := m.RemoveJetIndexesUntilMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetIndexesUntil")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetIndexesUntilMock.mainExpectation != nil {

		input := m.RemoveJetIndexesUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockRemoveJetIndexesUntilInput{p, p1, p2, p3}, "Cleaner.RemoveJetIndexesUntil got unexpected parameters")
		}

		result := m.RemoveJetIndexesUntilMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetIndexesUntil")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetIndexesUntilFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetIndexesUntil. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RemoveJetIndexesUntilFunc(p, p1, p2, p3)
}

// RemoveJetIndexesUntilMinimockCounter returns a count of CleanerMock.RemoveJetIndexesUntilFunc invocations
func (m *CleanerMock) RemoveJetIndexesUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetIndexesUntilCounter)
}

// RemoveJetIndexesUntilMinimockPreCounter returns the value of CleanerMock.RemoveJetIndexesUntil invocations
func (m *CleanerMock) RemoveJetIndexesUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetIndexesUntilPreCounter)
}

// RemoveJetIndexesUntilFinished returns true if mock invocations count is ok
func (m *CleanerMock) RemoveJetIndexesUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveJetIndexesUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveJetIndexesUntilCounter) == uint64(len(m.RemoveJetIndexesUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveJetIndexesUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveJetIndexesUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveJetIndexesUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveJetIndexesUntilCounter) > 0
	}

	return true
}

type mCleanerMockRemoveJetRecordsUntil struct {
	mock              *CleanerMock
	mainExpectation   *CleanerMockRemoveJetRecordsUntilExpectation
	expectationSeries []*CleanerMockRemoveJetRecordsUntilExpectation
}

type CleanerMockRemoveJetRecordsUntilExpectation struct {
	input  *CleanerMockRemoveJetRecordsUntilInput
	result *CleanerMockRemoveJetRecordsUntilResult
}

type CleanerMockRemoveJetRecordsUntilInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 recentstorage.RecentStorage
}

type CleanerMockRemoveJetRecordsUntilResult struct {
	r  int
	r1 error
}

// Expect specifies that invocation of Cleaner.RemoveJetRecordsUntil is expected from 1 to Infinity times
func (m *mCleanerMockRemoveJetRecordsUntil) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) *mCleanerMockRemoveJetRecordsUntil {
	m.mock.RemoveJetRecordsUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetRecordsUntilExpectation{}
	}
	m.mainExpectation.input = &CleanerMockRemoveJetRecordsUntilInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of Cleaner.RemoveJetRecordsUntil
func (m *mCleanerMockRemoveJetRecordsUntil) Return(r int, r1 error) *CleanerMock {
	m.mock.RemoveJetRecordsUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CleanerMockRemoveJetRecordsUntilExpectation{}
	}
	m.mainExpectation.result = &CleanerMockRemoveJetRecordsUntilResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of Cleaner.RemoveJetRecordsUntil is expected once
func (m *mCleanerMockRemoveJetRecordsUntil) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) *CleanerMockRemoveJetRecordsUntilExpectation {
	m.mock.RemoveJetRecordsUntilFunc = nil
	m.mainExpectation = nil

	expectation := &CleanerMockRemoveJetRecordsUntilExpectation{}
	expectation.input = &CleanerMockRemoveJetRecordsUntilInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CleanerMockRemoveJetRecordsUntilExpectation) Return(r int, r1 error) {
	e.result = &CleanerMockRemoveJetRecordsUntilResult{r, r1}
}

// Set uses given function f as a mock of Cleaner.RemoveJetRecordsUntil method
func (m *mCleanerMockRemoveJetRecordsUntil) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r int, r1 error)) *CleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveJetRecordsUntilFunc = f
	return m.mock
}

// RemoveJetRecordsUntil implements github.com/insolar/insolar/ledger/storage.Cleaner interface
func (m *CleanerMock) RemoveJetRecordsUntil(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 recentstorage.RecentStorage) (r int, r1 error) {
	counter := atomic.AddUint64(&m.RemoveJetRecordsUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveJetRecordsUntilCounter, 1)

	if len(m.RemoveJetRecordsUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveJetRecordsUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetRecordsUntil. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RemoveJetRecordsUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CleanerMockRemoveJetRecordsUntilInput{p, p1, p2, p3}, "Cleaner.RemoveJetRecordsUntil got unexpected parameters")

		result := m.RemoveJetRecordsUntilMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetRecordsUntil")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetRecordsUntilMock.mainExpectation != nil {

		input := m.RemoveJetRecordsUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CleanerMockRemoveJetRecordsUntilInput{p, p1, p2, p3}, "Cleaner.RemoveJetRecordsUntil got unexpected parameters")
		}

		result := m.RemoveJetRecordsUntilMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CleanerMock.RemoveJetRecordsUntil")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RemoveJetRecordsUntilFunc == nil {
		m.t.Fatalf("Unexpected call to CleanerMock.RemoveJetRecordsUntil. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RemoveJetRecordsUntilFunc(p, p1, p2, p3)
}

// RemoveJetRecordsUntilMinimockCounter returns a count of CleanerMock.RemoveJetRecordsUntilFunc invocations
func (m *CleanerMock) RemoveJetRecordsUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetRecordsUntilCounter)
}

// RemoveJetRecordsUntilMinimockPreCounter returns the value of CleanerMock.RemoveJetRecordsUntil invocations
func (m *CleanerMock) RemoveJetRecordsUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJetRecordsUntilPreCounter)
}

// RemoveJetRecordsUntilFinished returns true if mock invocations count is ok
func (m *CleanerMock) RemoveJetRecordsUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveJetRecordsUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveJetRecordsUntilCounter) == uint64(len(m.RemoveJetRecordsUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveJetRecordsUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveJetRecordsUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveJetRecordsUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveJetRecordsUntilCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CleanerMock) ValidateCallCounters() {

	if !m.RemoveAllForJetUntilPulseFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveAllForJetUntilPulse")
	}

	if !m.RemoveJetBlobsUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetBlobsUntil")
	}

	if !m.RemoveJetDropsUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetDropsUntil")
	}

	if !m.RemoveJetIndexesUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetIndexesUntil")
	}

	if !m.RemoveJetRecordsUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetRecordsUntil")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CleanerMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CleanerMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CleanerMock) MinimockFinish() {

	if !m.RemoveAllForJetUntilPulseFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveAllForJetUntilPulse")
	}

	if !m.RemoveJetBlobsUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetBlobsUntil")
	}

	if !m.RemoveJetDropsUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetDropsUntil")
	}

	if !m.RemoveJetIndexesUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetIndexesUntil")
	}

	if !m.RemoveJetRecordsUntilFinished() {
		m.t.Fatal("Expected call to CleanerMock.RemoveJetRecordsUntil")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *CleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RemoveAllForJetUntilPulseFinished()
		ok = ok && m.RemoveJetBlobsUntilFinished()
		ok = ok && m.RemoveJetDropsUntilFinished()
		ok = ok && m.RemoveJetIndexesUntilFinished()
		ok = ok && m.RemoveJetRecordsUntilFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveAllForJetUntilPulseFinished() {
				m.t.Error("Expected call to CleanerMock.RemoveAllForJetUntilPulse")
			}

			if !m.RemoveJetBlobsUntilFinished() {
				m.t.Error("Expected call to CleanerMock.RemoveJetBlobsUntil")
			}

			if !m.RemoveJetDropsUntilFinished() {
				m.t.Error("Expected call to CleanerMock.RemoveJetDropsUntil")
			}

			if !m.RemoveJetIndexesUntilFinished() {
				m.t.Error("Expected call to CleanerMock.RemoveJetIndexesUntil")
			}

			if !m.RemoveJetRecordsUntilFinished() {
				m.t.Error("Expected call to CleanerMock.RemoveJetRecordsUntil")
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
func (m *CleanerMock) AllMocksCalled() bool {

	if !m.RemoveAllForJetUntilPulseFinished() {
		return false
	}

	if !m.RemoveJetBlobsUntilFinished() {
		return false
	}

	if !m.RemoveJetDropsUntilFinished() {
		return false
	}

	if !m.RemoveJetIndexesUntilFinished() {
		return false
	}

	if !m.RemoveJetRecordsUntilFinished() {
		return false
	}

	return true
}

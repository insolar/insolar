package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PendingModifier" can be found in github.com/insolar/insolar/ledger/object
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

// PendingModifierMock implements github.com/insolar/insolar/ledger/object.PendingModifier
type PendingModifierMock struct {
	t minimock.Tester

	SetRecordFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Virtual) (r error)
	SetRecordCounter    uint64
	SetRecordPreCounter uint64
	SetRecordMock       mPendingModifierMockSetRecord
}

// NewPendingModifierMock returns a mock for github.com/insolar/insolar/ledger/object.PendingModifier
func NewPendingModifierMock(t minimock.Tester) *PendingModifierMock {
	m := &PendingModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetRecordMock = mPendingModifierMockSetRecord{mock: m}

	return m
}

type mPendingModifierMockSetRecord struct {
	mock              *PendingModifierMock
	mainExpectation   *PendingModifierMockSetRecordExpectation
	expectationSeries []*PendingModifierMockSetRecordExpectation
}

type PendingModifierMockSetRecordExpectation struct {
	input  *PendingModifierMockSetRecordInput
	result *PendingModifierMockSetRecordResult
}

type PendingModifierMockSetRecordInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 record.Virtual
}

type PendingModifierMockSetRecordResult struct {
	r error
}

// Expect specifies that invocation of PendingModifier.SetRecord is expected from 1 to Infinity times
func (m *mPendingModifierMockSetRecord) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Virtual) *mPendingModifierMockSetRecord {
	m.mock.SetRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetRecordExpectation{}
	}
	m.mainExpectation.input = &PendingModifierMockSetRecordInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of PendingModifier.SetRecord
func (m *mPendingModifierMockSetRecord) Return(r error) *PendingModifierMock {
	m.mock.SetRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetRecordExpectation{}
	}
	m.mainExpectation.result = &PendingModifierMockSetRecordResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of PendingModifier.SetRecord is expected once
func (m *mPendingModifierMockSetRecord) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Virtual) *PendingModifierMockSetRecordExpectation {
	m.mock.SetRecordFunc = nil
	m.mainExpectation = nil

	expectation := &PendingModifierMockSetRecordExpectation{}
	expectation.input = &PendingModifierMockSetRecordInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingModifierMockSetRecordExpectation) Return(r error) {
	e.result = &PendingModifierMockSetRecordResult{r}
}

// Set uses given function f as a mock of PendingModifier.SetRecord method
func (m *mPendingModifierMockSetRecord) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Virtual) (r error)) *PendingModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetRecordFunc = f
	return m.mock
}

// SetRecord implements github.com/insolar/insolar/ledger/object.PendingModifier interface
func (m *PendingModifierMock) SetRecord(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Virtual) (r error) {
	counter := atomic.AddUint64(&m.SetRecordPreCounter, 1)
	defer atomic.AddUint64(&m.SetRecordCounter, 1)

	if len(m.SetRecordMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetRecordMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingModifierMock.SetRecord. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetRecordMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingModifierMockSetRecordInput{p, p1, p2, p3}, "PendingModifier.SetRecord got unexpected parameters")

		result := m.SetRecordMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetRecord")
			return
		}

		r = result.r

		return
	}

	if m.SetRecordMock.mainExpectation != nil {

		input := m.SetRecordMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingModifierMockSetRecordInput{p, p1, p2, p3}, "PendingModifier.SetRecord got unexpected parameters")
		}

		result := m.SetRecordMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetRecord")
		}

		r = result.r

		return
	}

	if m.SetRecordFunc == nil {
		m.t.Fatalf("Unexpected call to PendingModifierMock.SetRecord. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetRecordFunc(p, p1, p2, p3)
}

// SetRecordMinimockCounter returns a count of PendingModifierMock.SetRecordFunc invocations
func (m *PendingModifierMock) SetRecordMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRecordCounter)
}

// SetRecordMinimockPreCounter returns the value of PendingModifierMock.SetRecord invocations
func (m *PendingModifierMock) SetRecordMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRecordPreCounter)
}

// SetRecordFinished returns true if mock invocations count is ok
func (m *PendingModifierMock) SetRecordFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetRecordMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetRecordCounter) == uint64(len(m.SetRecordMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetRecordMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetRecordCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetRecordFunc != nil {
		return atomic.LoadUint64(&m.SetRecordCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingModifierMock) ValidateCallCounters() {

	if !m.SetRecordFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetRecord")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingModifierMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PendingModifierMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PendingModifierMock) MinimockFinish() {

	if !m.SetRecordFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetRecord")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PendingModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *PendingModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetRecordFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetRecordFinished() {
				m.t.Error("Expected call to PendingModifierMock.SetRecord")
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
func (m *PendingModifierMock) AllMocksCalled() bool {

	if !m.SetRecordFinished() {
		return false
	}

	return true
}

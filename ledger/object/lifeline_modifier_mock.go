package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// LifelineModifierMock implements github.com/insolar/insolar/ledger/object.LifelineModifier
type LifelineModifierMock struct {
	t minimock.Tester

	SetFunc       func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mLifelineModifierMockSet
}

// NewLifelineModifierMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineModifier
func NewLifelineModifierMock(t minimock.Tester) *LifelineModifierMock {
	m := &LifelineModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetMock = mLifelineModifierMockSet{mock: m}

	return m
}

type mLifelineModifierMockSet struct {
	mock              *LifelineModifierMock
	mainExpectation   *LifelineModifierMockSetExpectation
	expectationSeries []*LifelineModifierMockSetExpectation
}

type LifelineModifierMockSetExpectation struct {
	input  *LifelineModifierMockSetInput
	result *LifelineModifierMockSetResult
}

type LifelineModifierMockSetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 Lifeline
}

type LifelineModifierMockSetResult struct {
	r error
}

// Expect specifies that invocation of LifelineModifier.Set is expected from 1 to Infinity times
func (m *mLifelineModifierMockSet) Expect(p context.Context, p1 insolar.ID, p2 Lifeline) *mLifelineModifierMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineModifierMockSetExpectation{}
	}
	m.mainExpectation.input = &LifelineModifierMockSetInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of LifelineModifier.Set
func (m *mLifelineModifierMockSet) Return(r error) *LifelineModifierMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineModifierMockSetExpectation{}
	}
	m.mainExpectation.result = &LifelineModifierMockSetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of LifelineModifier.Set is expected once
func (m *mLifelineModifierMockSet) ExpectOnce(p context.Context, p1 insolar.ID, p2 Lifeline) *LifelineModifierMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineModifierMockSetExpectation{}
	expectation.input = &LifelineModifierMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineModifierMockSetExpectation) Return(r error) {
	e.result = &LifelineModifierMockSetResult{r}
}

// Set uses given function f as a mock of LifelineModifier.Set method
func (m *mLifelineModifierMockSet) Set(f func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)) *LifelineModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

// Set implements github.com/insolar/insolar/ledger/object.LifelineModifier interface
func (m *LifelineModifierMock) Set(p context.Context, p1 insolar.ID, p2 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineModifierMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineModifierMockSetInput{p, p1, p2}, "LifelineModifier.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineModifierMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineModifierMockSetInput{p, p1, p2}, "LifelineModifier.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineModifierMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineModifierMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

// SetMinimockCounter returns a count of LifelineModifierMock.SetFunc invocations
func (m *LifelineModifierMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

// SetMinimockPreCounter returns the value of LifelineModifierMock.Set invocations
func (m *LifelineModifierMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

// SetFinished returns true if mock invocations count is ok
func (m *LifelineModifierMock) SetFinished() bool {
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
func (m *LifelineModifierMock) ValidateCallCounters() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to LifelineModifierMock.Set")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineModifierMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LifelineModifierMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LifelineModifierMock) MinimockFinish() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to LifelineModifierMock.Set")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LifelineModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *LifelineModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetFinished() {
				m.t.Error("Expected call to LifelineModifierMock.Set")
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
func (m *LifelineModifierMock) AllMocksCalled() bool {

	if !m.SetFinished() {
		return false
	}

	return true
}

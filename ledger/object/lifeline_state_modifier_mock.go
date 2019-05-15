package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineStateModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// LifelineStateModifierMock implements github.com/insolar/insolar/ledger/object.LifelineStateModifier
type LifelineStateModifierMock struct {
	t minimock.Tester

	SetLifelineUsageFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)
	SetLifelineUsageCounter    uint64
	SetLifelineUsagePreCounter uint64
	SetLifelineUsageMock       mLifelineStateModifierMockSetLifelineUsage
}

// NewLifelineStateModifierMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineStateModifier
func NewLifelineStateModifierMock(t minimock.Tester) *LifelineStateModifierMock {
	m := &LifelineStateModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetLifelineUsageMock = mLifelineStateModifierMockSetLifelineUsage{mock: m}

	return m
}

type mLifelineStateModifierMockSetLifelineUsage struct {
	mock              *LifelineStateModifierMock
	mainExpectation   *LifelineStateModifierMockSetLifelineUsageExpectation
	expectationSeries []*LifelineStateModifierMockSetLifelineUsageExpectation
}

type LifelineStateModifierMockSetLifelineUsageExpectation struct {
	input  *LifelineStateModifierMockSetLifelineUsageInput
	result *LifelineStateModifierMockSetLifelineUsageResult
}

type LifelineStateModifierMockSetLifelineUsageInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type LifelineStateModifierMockSetLifelineUsageResult struct {
	r error
}

// Expect specifies that invocation of LifelineStateModifier.SetLifelineUsage is expected from 1 to Infinity times
func (m *mLifelineStateModifierMockSetLifelineUsage) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mLifelineStateModifierMockSetLifelineUsage {
	m.mock.SetLifelineUsageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineStateModifierMockSetLifelineUsageExpectation{}
	}
	m.mainExpectation.input = &LifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of LifelineStateModifier.SetLifelineUsage
func (m *mLifelineStateModifierMockSetLifelineUsage) Return(r error) *LifelineStateModifierMock {
	m.mock.SetLifelineUsageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineStateModifierMockSetLifelineUsageExpectation{}
	}
	m.mainExpectation.result = &LifelineStateModifierMockSetLifelineUsageResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of LifelineStateModifier.SetLifelineUsage is expected once
func (m *mLifelineStateModifierMockSetLifelineUsage) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *LifelineStateModifierMockSetLifelineUsageExpectation {
	m.mock.SetLifelineUsageFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineStateModifierMockSetLifelineUsageExpectation{}
	expectation.input = &LifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineStateModifierMockSetLifelineUsageExpectation) Return(r error) {
	e.result = &LifelineStateModifierMockSetLifelineUsageResult{r}
}

// Set uses given function f as a mock of LifelineStateModifier.SetLifelineUsage method
func (m *mLifelineStateModifierMockSetLifelineUsage) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)) *LifelineStateModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLifelineUsageFunc = f
	return m.mock
}

// SetLifelineUsage implements github.com/insolar/insolar/ledger/object.LifelineStateModifier interface
func (m *LifelineStateModifierMock) SetLifelineUsage(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetLifelineUsagePreCounter, 1)
	defer atomic.AddUint64(&m.SetLifelineUsageCounter, 1)

	if len(m.SetLifelineUsageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLifelineUsageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineStateModifierMock.SetLifelineUsage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetLifelineUsageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}, "LifelineStateModifier.SetLifelineUsage got unexpected parameters")

		result := m.SetLifelineUsageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineStateModifierMock.SetLifelineUsage")
			return
		}

		r = result.r

		return
	}

	if m.SetLifelineUsageMock.mainExpectation != nil {

		input := m.SetLifelineUsageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}, "LifelineStateModifier.SetLifelineUsage got unexpected parameters")
		}

		result := m.SetLifelineUsageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineStateModifierMock.SetLifelineUsage")
		}

		r = result.r

		return
	}

	if m.SetLifelineUsageFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineStateModifierMock.SetLifelineUsage. %v %v %v", p, p1, p2)
		return
	}

	return m.SetLifelineUsageFunc(p, p1, p2)
}

// SetLifelineUsageMinimockCounter returns a count of LifelineStateModifierMock.SetLifelineUsageFunc invocations
func (m *LifelineStateModifierMock) SetLifelineUsageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineUsageCounter)
}

// SetLifelineUsageMinimockPreCounter returns the value of LifelineStateModifierMock.SetLifelineUsage invocations
func (m *LifelineStateModifierMock) SetLifelineUsageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineUsagePreCounter)
}

// SetLifelineUsageFinished returns true if mock invocations count is ok
func (m *LifelineStateModifierMock) SetLifelineUsageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetLifelineUsageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetLifelineUsageCounter) == uint64(len(m.SetLifelineUsageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetLifelineUsageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetLifelineUsageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetLifelineUsageFunc != nil {
		return atomic.LoadUint64(&m.SetLifelineUsageCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineStateModifierMock) ValidateCallCounters() {

	if !m.SetLifelineUsageFinished() {
		m.t.Fatal("Expected call to LifelineStateModifierMock.SetLifelineUsage")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineStateModifierMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LifelineStateModifierMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LifelineStateModifierMock) MinimockFinish() {

	if !m.SetLifelineUsageFinished() {
		m.t.Fatal("Expected call to LifelineStateModifierMock.SetLifelineUsage")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LifelineStateModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *LifelineStateModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetLifelineUsageFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetLifelineUsageFinished() {
				m.t.Error("Expected call to LifelineStateModifierMock.SetLifelineUsage")
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
func (m *LifelineStateModifierMock) AllMocksCalled() bool {

	if !m.SetLifelineUsageFinished() {
		return false
	}

	return true
}

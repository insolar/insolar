package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexLifelineStateModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexLifelineStateModifierMock implements github.com/insolar/insolar/ledger/object.IndexLifelineStateModifier
type IndexLifelineStateModifierMock struct {
	t minimock.Tester

	SetLifelineUsageFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)
	SetLifelineUsageCounter    uint64
	SetLifelineUsagePreCounter uint64
	SetLifelineUsageMock       mIndexLifelineStateModifierMockSetLifelineUsage
}

// NewIndexLifelineStateModifierMock returns a mock for github.com/insolar/insolar/ledger/object.IndexLifelineStateModifier
func NewIndexLifelineStateModifierMock(t minimock.Tester) *IndexLifelineStateModifierMock {
	m := &IndexLifelineStateModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetLifelineUsageMock = mIndexLifelineStateModifierMockSetLifelineUsage{mock: m}

	return m
}

type mIndexLifelineStateModifierMockSetLifelineUsage struct {
	mock              *IndexLifelineStateModifierMock
	mainExpectation   *IndexLifelineStateModifierMockSetLifelineUsageExpectation
	expectationSeries []*IndexLifelineStateModifierMockSetLifelineUsageExpectation
}

type IndexLifelineStateModifierMockSetLifelineUsageExpectation struct {
	input  *IndexLifelineStateModifierMockSetLifelineUsageInput
	result *IndexLifelineStateModifierMockSetLifelineUsageResult
}

type IndexLifelineStateModifierMockSetLifelineUsageInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type IndexLifelineStateModifierMockSetLifelineUsageResult struct {
	r error
}

// Expect specifies that invocation of IndexLifelineStateModifier.SetLifelineUsage is expected from 1 to Infinity times
func (m *mIndexLifelineStateModifierMockSetLifelineUsage) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mIndexLifelineStateModifierMockSetLifelineUsage {
	m.mock.SetLifelineUsageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLifelineStateModifierMockSetLifelineUsageExpectation{}
	}
	m.mainExpectation.input = &IndexLifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of IndexLifelineStateModifier.SetLifelineUsage
func (m *mIndexLifelineStateModifierMockSetLifelineUsage) Return(r error) *IndexLifelineStateModifierMock {
	m.mock.SetLifelineUsageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLifelineStateModifierMockSetLifelineUsageExpectation{}
	}
	m.mainExpectation.result = &IndexLifelineStateModifierMockSetLifelineUsageResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexLifelineStateModifier.SetLifelineUsage is expected once
func (m *mIndexLifelineStateModifierMockSetLifelineUsage) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *IndexLifelineStateModifierMockSetLifelineUsageExpectation {
	m.mock.SetLifelineUsageFunc = nil
	m.mainExpectation = nil

	expectation := &IndexLifelineStateModifierMockSetLifelineUsageExpectation{}
	expectation.input = &IndexLifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexLifelineStateModifierMockSetLifelineUsageExpectation) Return(r error) {
	e.result = &IndexLifelineStateModifierMockSetLifelineUsageResult{r}
}

// Set uses given function f as a mock of IndexLifelineStateModifier.SetLifelineUsage method
func (m *mIndexLifelineStateModifierMockSetLifelineUsage) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)) *IndexLifelineStateModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLifelineUsageFunc = f
	return m.mock
}

// SetLifelineUsage implements github.com/insolar/insolar/ledger/object.IndexLifelineStateModifier interface
func (m *IndexLifelineStateModifierMock) SetLifelineUsage(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetLifelineUsagePreCounter, 1)
	defer atomic.AddUint64(&m.SetLifelineUsageCounter, 1)

	if len(m.SetLifelineUsageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLifelineUsageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexLifelineStateModifierMock.SetLifelineUsage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetLifelineUsageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexLifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}, "IndexLifelineStateModifier.SetLifelineUsage got unexpected parameters")

		result := m.SetLifelineUsageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexLifelineStateModifierMock.SetLifelineUsage")
			return
		}

		r = result.r

		return
	}

	if m.SetLifelineUsageMock.mainExpectation != nil {

		input := m.SetLifelineUsageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexLifelineStateModifierMockSetLifelineUsageInput{p, p1, p2}, "IndexLifelineStateModifier.SetLifelineUsage got unexpected parameters")
		}

		result := m.SetLifelineUsageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexLifelineStateModifierMock.SetLifelineUsage")
		}

		r = result.r

		return
	}

	if m.SetLifelineUsageFunc == nil {
		m.t.Fatalf("Unexpected call to IndexLifelineStateModifierMock.SetLifelineUsage. %v %v %v", p, p1, p2)
		return
	}

	return m.SetLifelineUsageFunc(p, p1, p2)
}

// SetLifelineUsageMinimockCounter returns a count of IndexLifelineStateModifierMock.SetLifelineUsageFunc invocations
func (m *IndexLifelineStateModifierMock) SetLifelineUsageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineUsageCounter)
}

// SetLifelineUsageMinimockPreCounter returns the value of IndexLifelineStateModifierMock.SetLifelineUsage invocations
func (m *IndexLifelineStateModifierMock) SetLifelineUsageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineUsagePreCounter)
}

// SetLifelineUsageFinished returns true if mock invocations count is ok
func (m *IndexLifelineStateModifierMock) SetLifelineUsageFinished() bool {
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
func (m *IndexLifelineStateModifierMock) ValidateCallCounters() {

	if !m.SetLifelineUsageFinished() {
		m.t.Fatal("Expected call to IndexLifelineStateModifierMock.SetLifelineUsage")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexLifelineStateModifierMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexLifelineStateModifierMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexLifelineStateModifierMock) MinimockFinish() {

	if !m.SetLifelineUsageFinished() {
		m.t.Fatal("Expected call to IndexLifelineStateModifierMock.SetLifelineUsage")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexLifelineStateModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexLifelineStateModifierMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IndexLifelineStateModifierMock.SetLifelineUsage")
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
func (m *IndexLifelineStateModifierMock) AllMocksCalled() bool {

	if !m.SetLifelineUsageFinished() {
		return false
	}

	return true
}

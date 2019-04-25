package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ExtendedLifelineModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// ExtendedLifelineModifierMock implements github.com/insolar/insolar/ledger/object.ExtendedLifelineModifier
type ExtendedLifelineModifierMock struct {
	t minimock.Tester

	SetUsageForPulseFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber)
	SetUsageForPulseCounter    uint64
	SetUsageForPulsePreCounter uint64
	SetUsageForPulseMock       mExtendedLifelineModifierMockSetUsageForPulse

	SetWithMetaFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) (r error)
	SetWithMetaCounter    uint64
	SetWithMetaPreCounter uint64
	SetWithMetaMock       mExtendedLifelineModifierMockSetWithMeta
}

// NewExtendedLifelineModifierMock returns a mock for github.com/insolar/insolar/ledger/object.ExtendedLifelineModifier
func NewExtendedLifelineModifierMock(t minimock.Tester) *ExtendedLifelineModifierMock {
	m := &ExtendedLifelineModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetUsageForPulseMock = mExtendedLifelineModifierMockSetUsageForPulse{mock: m}
	m.SetWithMetaMock = mExtendedLifelineModifierMockSetWithMeta{mock: m}

	return m
}

type mExtendedLifelineModifierMockSetUsageForPulse struct {
	mock              *ExtendedLifelineModifierMock
	mainExpectation   *ExtendedLifelineModifierMockSetUsageForPulseExpectation
	expectationSeries []*ExtendedLifelineModifierMockSetUsageForPulseExpectation
}

type ExtendedLifelineModifierMockSetUsageForPulseExpectation struct {
	input *ExtendedLifelineModifierMockSetUsageForPulseInput
}

type ExtendedLifelineModifierMockSetUsageForPulseInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

// Expect specifies that invocation of ExtendedLifelineModifier.SetUsageForPulse is expected from 1 to Infinity times
func (m *mExtendedLifelineModifierMockSetUsageForPulse) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mExtendedLifelineModifierMockSetUsageForPulse {
	m.mock.SetUsageForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedLifelineModifierMockSetUsageForPulseExpectation{}
	}
	m.mainExpectation.input = &ExtendedLifelineModifierMockSetUsageForPulseInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of ExtendedLifelineModifier.SetUsageForPulse
func (m *mExtendedLifelineModifierMockSetUsageForPulse) Return() *ExtendedLifelineModifierMock {
	m.mock.SetUsageForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedLifelineModifierMockSetUsageForPulseExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of ExtendedLifelineModifier.SetUsageForPulse is expected once
func (m *mExtendedLifelineModifierMockSetUsageForPulse) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *ExtendedLifelineModifierMockSetUsageForPulseExpectation {
	m.mock.SetUsageForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ExtendedLifelineModifierMockSetUsageForPulseExpectation{}
	expectation.input = &ExtendedLifelineModifierMockSetUsageForPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of ExtendedLifelineModifier.SetUsageForPulse method
func (m *mExtendedLifelineModifierMockSetUsageForPulse) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber)) *ExtendedLifelineModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetUsageForPulseFunc = f
	return m.mock
}

// SetUsageForPulse implements github.com/insolar/insolar/ledger/object.ExtendedLifelineModifier interface
func (m *ExtendedLifelineModifierMock) SetUsageForPulse(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.SetUsageForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetUsageForPulseCounter, 1)

	if len(m.SetUsageForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetUsageForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExtendedLifelineModifierMock.SetUsageForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetUsageForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExtendedLifelineModifierMockSetUsageForPulseInput{p, p1, p2}, "ExtendedLifelineModifier.SetUsageForPulse got unexpected parameters")

		return
	}

	if m.SetUsageForPulseMock.mainExpectation != nil {

		input := m.SetUsageForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExtendedLifelineModifierMockSetUsageForPulseInput{p, p1, p2}, "ExtendedLifelineModifier.SetUsageForPulse got unexpected parameters")
		}

		return
	}

	if m.SetUsageForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ExtendedLifelineModifierMock.SetUsageForPulse. %v %v %v", p, p1, p2)
		return
	}

	m.SetUsageForPulseFunc(p, p1, p2)
}

// SetUsageForPulseMinimockCounter returns a count of ExtendedLifelineModifierMock.SetUsageForPulseFunc invocations
func (m *ExtendedLifelineModifierMock) SetUsageForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetUsageForPulseCounter)
}

// SetUsageForPulseMinimockPreCounter returns the value of ExtendedLifelineModifierMock.SetUsageForPulse invocations
func (m *ExtendedLifelineModifierMock) SetUsageForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetUsageForPulsePreCounter)
}

// SetUsageForPulseFinished returns true if mock invocations count is ok
func (m *ExtendedLifelineModifierMock) SetUsageForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetUsageForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetUsageForPulseCounter) == uint64(len(m.SetUsageForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetUsageForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetUsageForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetUsageForPulseFunc != nil {
		return atomic.LoadUint64(&m.SetUsageForPulseCounter) > 0
	}

	return true
}

type mExtendedLifelineModifierMockSetWithMeta struct {
	mock              *ExtendedLifelineModifierMock
	mainExpectation   *ExtendedLifelineModifierMockSetWithMetaExpectation
	expectationSeries []*ExtendedLifelineModifierMockSetWithMetaExpectation
}

type ExtendedLifelineModifierMockSetWithMetaExpectation struct {
	input  *ExtendedLifelineModifierMockSetWithMetaInput
	result *ExtendedLifelineModifierMockSetWithMetaResult
}

type ExtendedLifelineModifierMockSetWithMetaInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 Lifeline
}

type ExtendedLifelineModifierMockSetWithMetaResult struct {
	r error
}

// Expect specifies that invocation of ExtendedLifelineModifier.SetWithMeta is expected from 1 to Infinity times
func (m *mExtendedLifelineModifierMockSetWithMeta) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) *mExtendedLifelineModifierMockSetWithMeta {
	m.mock.SetWithMetaFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedLifelineModifierMockSetWithMetaExpectation{}
	}
	m.mainExpectation.input = &ExtendedLifelineModifierMockSetWithMetaInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of ExtendedLifelineModifier.SetWithMeta
func (m *mExtendedLifelineModifierMockSetWithMeta) Return(r error) *ExtendedLifelineModifierMock {
	m.mock.SetWithMetaFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedLifelineModifierMockSetWithMetaExpectation{}
	}
	m.mainExpectation.result = &ExtendedLifelineModifierMockSetWithMetaResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of ExtendedLifelineModifier.SetWithMeta is expected once
func (m *mExtendedLifelineModifierMockSetWithMeta) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) *ExtendedLifelineModifierMockSetWithMetaExpectation {
	m.mock.SetWithMetaFunc = nil
	m.mainExpectation = nil

	expectation := &ExtendedLifelineModifierMockSetWithMetaExpectation{}
	expectation.input = &ExtendedLifelineModifierMockSetWithMetaInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExtendedLifelineModifierMockSetWithMetaExpectation) Return(r error) {
	e.result = &ExtendedLifelineModifierMockSetWithMetaResult{r}
}

// Set uses given function f as a mock of ExtendedLifelineModifier.SetWithMeta method
func (m *mExtendedLifelineModifierMockSetWithMeta) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) (r error)) *ExtendedLifelineModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetWithMetaFunc = f
	return m.mock
}

// SetWithMeta implements github.com/insolar/insolar/ledger/object.ExtendedLifelineModifier interface
func (m *ExtendedLifelineModifierMock) SetWithMeta(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetWithMetaPreCounter, 1)
	defer atomic.AddUint64(&m.SetWithMetaCounter, 1)

	if len(m.SetWithMetaMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetWithMetaMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExtendedLifelineModifierMock.SetWithMeta. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetWithMetaMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExtendedLifelineModifierMockSetWithMetaInput{p, p1, p2, p3}, "ExtendedLifelineModifier.SetWithMeta got unexpected parameters")

		result := m.SetWithMetaMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExtendedLifelineModifierMock.SetWithMeta")
			return
		}

		r = result.r

		return
	}

	if m.SetWithMetaMock.mainExpectation != nil {

		input := m.SetWithMetaMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExtendedLifelineModifierMockSetWithMetaInput{p, p1, p2, p3}, "ExtendedLifelineModifier.SetWithMeta got unexpected parameters")
		}

		result := m.SetWithMetaMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExtendedLifelineModifierMock.SetWithMeta")
		}

		r = result.r

		return
	}

	if m.SetWithMetaFunc == nil {
		m.t.Fatalf("Unexpected call to ExtendedLifelineModifierMock.SetWithMeta. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetWithMetaFunc(p, p1, p2, p3)
}

// SetWithMetaMinimockCounter returns a count of ExtendedLifelineModifierMock.SetWithMetaFunc invocations
func (m *ExtendedLifelineModifierMock) SetWithMetaMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetWithMetaCounter)
}

// SetWithMetaMinimockPreCounter returns the value of ExtendedLifelineModifierMock.SetWithMeta invocations
func (m *ExtendedLifelineModifierMock) SetWithMetaMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetWithMetaPreCounter)
}

// SetWithMetaFinished returns true if mock invocations count is ok
func (m *ExtendedLifelineModifierMock) SetWithMetaFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetWithMetaMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetWithMetaCounter) == uint64(len(m.SetWithMetaMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetWithMetaMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetWithMetaCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetWithMetaFunc != nil {
		return atomic.LoadUint64(&m.SetWithMetaCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExtendedLifelineModifierMock) ValidateCallCounters() {

	if !m.SetUsageForPulseFinished() {
		m.t.Fatal("Expected call to ExtendedLifelineModifierMock.SetUsageForPulse")
	}

	if !m.SetWithMetaFinished() {
		m.t.Fatal("Expected call to ExtendedLifelineModifierMock.SetWithMeta")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExtendedLifelineModifierMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ExtendedLifelineModifierMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ExtendedLifelineModifierMock) MinimockFinish() {

	if !m.SetUsageForPulseFinished() {
		m.t.Fatal("Expected call to ExtendedLifelineModifierMock.SetUsageForPulse")
	}

	if !m.SetWithMetaFinished() {
		m.t.Fatal("Expected call to ExtendedLifelineModifierMock.SetWithMeta")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ExtendedLifelineModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *ExtendedLifelineModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetUsageForPulseFinished()
		ok = ok && m.SetWithMetaFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetUsageForPulseFinished() {
				m.t.Error("Expected call to ExtendedLifelineModifierMock.SetUsageForPulse")
			}

			if !m.SetWithMetaFinished() {
				m.t.Error("Expected call to ExtendedLifelineModifierMock.SetWithMeta")
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
func (m *ExtendedLifelineModifierMock) AllMocksCalled() bool {

	if !m.SetUsageForPulseFinished() {
		return false
	}

	if !m.SetWithMetaFinished() {
		return false
	}

	return true
}

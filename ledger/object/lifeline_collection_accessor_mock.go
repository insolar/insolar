package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineCollectionAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// LifelineCollectionAccessorMock implements github.com/insolar/insolar/ledger/object.LifelineCollectionAccessor
type LifelineCollectionAccessorMock struct {
	t minimock.Tester

	ForJetFunc       func(p context.Context, p1 insolar.JetID) (r map[insolar.ID]LifelineMeta)
	ForJetCounter    uint64
	ForJetPreCounter uint64
	ForJetMock       mLifelineCollectionAccessorMockForJet

	ForPulseAndJetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r map[insolar.ID]Lifeline)
	ForPulseAndJetCounter    uint64
	ForPulseAndJetPreCounter uint64
	ForPulseAndJetMock       mLifelineCollectionAccessorMockForPulseAndJet
}

// NewLifelineCollectionAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineCollectionAccessor
func NewLifelineCollectionAccessorMock(t minimock.Tester) *LifelineCollectionAccessorMock {
	m := &LifelineCollectionAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForJetMock = mLifelineCollectionAccessorMockForJet{mock: m}
	m.ForPulseAndJetMock = mLifelineCollectionAccessorMockForPulseAndJet{mock: m}

	return m
}

type mLifelineCollectionAccessorMockForJet struct {
	mock              *LifelineCollectionAccessorMock
	mainExpectation   *LifelineCollectionAccessorMockForJetExpectation
	expectationSeries []*LifelineCollectionAccessorMockForJetExpectation
}

type LifelineCollectionAccessorMockForJetExpectation struct {
	input  *LifelineCollectionAccessorMockForJetInput
	result *LifelineCollectionAccessorMockForJetResult
}

type LifelineCollectionAccessorMockForJetInput struct {
	p  context.Context
	p1 insolar.JetID
}

type LifelineCollectionAccessorMockForJetResult struct {
	r map[insolar.ID]LifelineMeta
}

// Expect specifies that invocation of LifelineCollectionAccessor.ForJet is expected from 1 to Infinity times
func (m *mLifelineCollectionAccessorMockForJet) Expect(p context.Context, p1 insolar.JetID) *mLifelineCollectionAccessorMockForJet {
	m.mock.ForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineCollectionAccessorMockForJetExpectation{}
	}
	m.mainExpectation.input = &LifelineCollectionAccessorMockForJetInput{p, p1}
	return m
}

// Return specifies results of invocation of LifelineCollectionAccessor.ForJet
func (m *mLifelineCollectionAccessorMockForJet) Return(r map[insolar.ID]LifelineMeta) *LifelineCollectionAccessorMock {
	m.mock.ForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineCollectionAccessorMockForJetExpectation{}
	}
	m.mainExpectation.result = &LifelineCollectionAccessorMockForJetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of LifelineCollectionAccessor.ForJet is expected once
func (m *mLifelineCollectionAccessorMockForJet) ExpectOnce(p context.Context, p1 insolar.JetID) *LifelineCollectionAccessorMockForJetExpectation {
	m.mock.ForJetFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineCollectionAccessorMockForJetExpectation{}
	expectation.input = &LifelineCollectionAccessorMockForJetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineCollectionAccessorMockForJetExpectation) Return(r map[insolar.ID]LifelineMeta) {
	e.result = &LifelineCollectionAccessorMockForJetResult{r}
}

// Set uses given function f as a mock of LifelineCollectionAccessor.ForJet method
func (m *mLifelineCollectionAccessorMockForJet) Set(f func(p context.Context, p1 insolar.JetID) (r map[insolar.ID]LifelineMeta)) *LifelineCollectionAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForJetFunc = f
	return m.mock
}

// ForJet implements github.com/insolar/insolar/ledger/object.LifelineCollectionAccessor interface
func (m *LifelineCollectionAccessorMock) ForJet(p context.Context, p1 insolar.JetID) (r map[insolar.ID]LifelineMeta) {
	counter := atomic.AddUint64(&m.ForJetPreCounter, 1)
	defer atomic.AddUint64(&m.ForJetCounter, 1)

	if len(m.ForJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineCollectionAccessorMock.ForJet. %v %v", p, p1)
			return
		}

		input := m.ForJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineCollectionAccessorMockForJetInput{p, p1}, "LifelineCollectionAccessor.ForJet got unexpected parameters")

		result := m.ForJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineCollectionAccessorMock.ForJet")
			return
		}

		r = result.r

		return
	}

	if m.ForJetMock.mainExpectation != nil {

		input := m.ForJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineCollectionAccessorMockForJetInput{p, p1}, "LifelineCollectionAccessor.ForJet got unexpected parameters")
		}

		result := m.ForJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineCollectionAccessorMock.ForJet")
		}

		r = result.r

		return
	}

	if m.ForJetFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineCollectionAccessorMock.ForJet. %v %v", p, p1)
		return
	}

	return m.ForJetFunc(p, p1)
}

// ForJetMinimockCounter returns a count of LifelineCollectionAccessorMock.ForJetFunc invocations
func (m *LifelineCollectionAccessorMock) ForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForJetCounter)
}

// ForJetMinimockPreCounter returns the value of LifelineCollectionAccessorMock.ForJet invocations
func (m *LifelineCollectionAccessorMock) ForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForJetPreCounter)
}

// ForJetFinished returns true if mock invocations count is ok
func (m *LifelineCollectionAccessorMock) ForJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForJetCounter) == uint64(len(m.ForJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForJetFunc != nil {
		return atomic.LoadUint64(&m.ForJetCounter) > 0
	}

	return true
}

type mLifelineCollectionAccessorMockForPulseAndJet struct {
	mock              *LifelineCollectionAccessorMock
	mainExpectation   *LifelineCollectionAccessorMockForPulseAndJetExpectation
	expectationSeries []*LifelineCollectionAccessorMockForPulseAndJetExpectation
}

type LifelineCollectionAccessorMockForPulseAndJetExpectation struct {
	input  *LifelineCollectionAccessorMockForPulseAndJetInput
	result *LifelineCollectionAccessorMockForPulseAndJetResult
}

type LifelineCollectionAccessorMockForPulseAndJetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type LifelineCollectionAccessorMockForPulseAndJetResult struct {
	r map[insolar.ID]Lifeline
}

// Expect specifies that invocation of LifelineCollectionAccessor.ForPulseAndJet is expected from 1 to Infinity times
func (m *mLifelineCollectionAccessorMockForPulseAndJet) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mLifelineCollectionAccessorMockForPulseAndJet {
	m.mock.ForPulseAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineCollectionAccessorMockForPulseAndJetExpectation{}
	}
	m.mainExpectation.input = &LifelineCollectionAccessorMockForPulseAndJetInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of LifelineCollectionAccessor.ForPulseAndJet
func (m *mLifelineCollectionAccessorMockForPulseAndJet) Return(r map[insolar.ID]Lifeline) *LifelineCollectionAccessorMock {
	m.mock.ForPulseAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineCollectionAccessorMockForPulseAndJetExpectation{}
	}
	m.mainExpectation.result = &LifelineCollectionAccessorMockForPulseAndJetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of LifelineCollectionAccessor.ForPulseAndJet is expected once
func (m *mLifelineCollectionAccessorMockForPulseAndJet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *LifelineCollectionAccessorMockForPulseAndJetExpectation {
	m.mock.ForPulseAndJetFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineCollectionAccessorMockForPulseAndJetExpectation{}
	expectation.input = &LifelineCollectionAccessorMockForPulseAndJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineCollectionAccessorMockForPulseAndJetExpectation) Return(r map[insolar.ID]Lifeline) {
	e.result = &LifelineCollectionAccessorMockForPulseAndJetResult{r}
}

// Set uses given function f as a mock of LifelineCollectionAccessor.ForPulseAndJet method
func (m *mLifelineCollectionAccessorMockForPulseAndJet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r map[insolar.ID]Lifeline)) *LifelineCollectionAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseAndJetFunc = f
	return m.mock
}

// ForPulseAndJet implements github.com/insolar/insolar/ledger/object.LifelineCollectionAccessor interface
func (m *LifelineCollectionAccessorMock) ForPulseAndJet(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r map[insolar.ID]Lifeline) {
	counter := atomic.AddUint64(&m.ForPulseAndJetPreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseAndJetCounter, 1)

	if len(m.ForPulseAndJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseAndJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineCollectionAccessorMock.ForPulseAndJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPulseAndJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineCollectionAccessorMockForPulseAndJetInput{p, p1, p2}, "LifelineCollectionAccessor.ForPulseAndJet got unexpected parameters")

		result := m.ForPulseAndJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineCollectionAccessorMock.ForPulseAndJet")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseAndJetMock.mainExpectation != nil {

		input := m.ForPulseAndJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineCollectionAccessorMockForPulseAndJetInput{p, p1, p2}, "LifelineCollectionAccessor.ForPulseAndJet got unexpected parameters")
		}

		result := m.ForPulseAndJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineCollectionAccessorMock.ForPulseAndJet")
		}

		r = result.r

		return
	}

	if m.ForPulseAndJetFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineCollectionAccessorMock.ForPulseAndJet. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPulseAndJetFunc(p, p1, p2)
}

// ForPulseAndJetMinimockCounter returns a count of LifelineCollectionAccessorMock.ForPulseAndJetFunc invocations
func (m *LifelineCollectionAccessorMock) ForPulseAndJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseAndJetCounter)
}

// ForPulseAndJetMinimockPreCounter returns the value of LifelineCollectionAccessorMock.ForPulseAndJet invocations
func (m *LifelineCollectionAccessorMock) ForPulseAndJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseAndJetPreCounter)
}

// ForPulseAndJetFinished returns true if mock invocations count is ok
func (m *LifelineCollectionAccessorMock) ForPulseAndJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseAndJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) == uint64(len(m.ForPulseAndJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseAndJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseAndJetFunc != nil {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineCollectionAccessorMock) ValidateCallCounters() {

	if !m.ForJetFinished() {
		m.t.Fatal("Expected call to LifelineCollectionAccessorMock.ForJet")
	}

	if !m.ForPulseAndJetFinished() {
		m.t.Fatal("Expected call to LifelineCollectionAccessorMock.ForPulseAndJet")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineCollectionAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LifelineCollectionAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LifelineCollectionAccessorMock) MinimockFinish() {

	if !m.ForJetFinished() {
		m.t.Fatal("Expected call to LifelineCollectionAccessorMock.ForJet")
	}

	if !m.ForPulseAndJetFinished() {
		m.t.Fatal("Expected call to LifelineCollectionAccessorMock.ForPulseAndJet")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LifelineCollectionAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *LifelineCollectionAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForJetFinished()
		ok = ok && m.ForPulseAndJetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForJetFinished() {
				m.t.Error("Expected call to LifelineCollectionAccessorMock.ForJet")
			}

			if !m.ForPulseAndJetFinished() {
				m.t.Error("Expected call to LifelineCollectionAccessorMock.ForPulseAndJet")
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
func (m *LifelineCollectionAccessorMock) AllMocksCalled() bool {

	if !m.ForJetFinished() {
		return false
	}

	if !m.ForPulseAndJetFinished() {
		return false
	}

	return true
}

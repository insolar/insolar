package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexAccessorMock implements github.com/insolar/insolar/ledger/object.IndexAccessor
type IndexAccessorMock struct {
	t minimock.Tester

	ForPNAndJetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []FilamentIndex)
	ForPNAndJetCounter    uint64
	ForPNAndJetPreCounter uint64
	ForPNAndJetMock       mIndexAccessorMockForPNAndJet

	IndexFunc       func(p insolar.PulseNumber, p1 insolar.ID) (r *LockedIndex)
	IndexCounter    uint64
	IndexPreCounter uint64
	IndexMock       mIndexAccessorMockIndex
}

// NewIndexAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.IndexAccessor
func NewIndexAccessorMock(t minimock.Tester) *IndexAccessorMock {
	m := &IndexAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPNAndJetMock = mIndexAccessorMockForPNAndJet{mock: m}
	m.IndexMock = mIndexAccessorMockIndex{mock: m}

	return m
}

type mIndexAccessorMockForPNAndJet struct {
	mock              *IndexAccessorMock
	mainExpectation   *IndexAccessorMockForPNAndJetExpectation
	expectationSeries []*IndexAccessorMockForPNAndJetExpectation
}

type IndexAccessorMockForPNAndJetExpectation struct {
	input  *IndexAccessorMockForPNAndJetInput
	result *IndexAccessorMockForPNAndJetResult
}

type IndexAccessorMockForPNAndJetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type IndexAccessorMockForPNAndJetResult struct {
	r []FilamentIndex
}

// Expect specifies that invocation of IndexAccessor.ForPNAndJet is expected from 1 to Infinity times
func (m *mIndexAccessorMockForPNAndJet) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mIndexAccessorMockForPNAndJet {
	m.mock.ForPNAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockForPNAndJetExpectation{}
	}
	m.mainExpectation.input = &IndexAccessorMockForPNAndJetInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of IndexAccessor.ForPNAndJet
func (m *mIndexAccessorMockForPNAndJet) Return(r []FilamentIndex) *IndexAccessorMock {
	m.mock.ForPNAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockForPNAndJetExpectation{}
	}
	m.mainExpectation.result = &IndexAccessorMockForPNAndJetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexAccessor.ForPNAndJet is expected once
func (m *mIndexAccessorMockForPNAndJet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *IndexAccessorMockForPNAndJetExpectation {
	m.mock.ForPNAndJetFunc = nil
	m.mainExpectation = nil

	expectation := &IndexAccessorMockForPNAndJetExpectation{}
	expectation.input = &IndexAccessorMockForPNAndJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexAccessorMockForPNAndJetExpectation) Return(r []FilamentIndex) {
	e.result = &IndexAccessorMockForPNAndJetResult{r}
}

// Set uses given function f as a mock of IndexAccessor.ForPNAndJet method
func (m *mIndexAccessorMockForPNAndJet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []FilamentIndex)) *IndexAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPNAndJetFunc = f
	return m.mock
}

// ForPNAndJet implements github.com/insolar/insolar/ledger/object.IndexAccessor interface
func (m *IndexAccessorMock) ForPNAndJet(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []FilamentIndex) {
	counter := atomic.AddUint64(&m.ForPNAndJetPreCounter, 1)
	defer atomic.AddUint64(&m.ForPNAndJetCounter, 1)

	if len(m.ForPNAndJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPNAndJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexAccessorMock.ForPNAndJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPNAndJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexAccessorMockForPNAndJetInput{p, p1, p2}, "IndexAccessor.ForPNAndJet got unexpected parameters")

		result := m.ForPNAndJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.ForPNAndJet")
			return
		}

		r = result.r

		return
	}

	if m.ForPNAndJetMock.mainExpectation != nil {

		input := m.ForPNAndJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexAccessorMockForPNAndJetInput{p, p1, p2}, "IndexAccessor.ForPNAndJet got unexpected parameters")
		}

		result := m.ForPNAndJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.ForPNAndJet")
		}

		r = result.r

		return
	}

	if m.ForPNAndJetFunc == nil {
		m.t.Fatalf("Unexpected call to IndexAccessorMock.ForPNAndJet. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPNAndJetFunc(p, p1, p2)
}

// ForPNAndJetMinimockCounter returns a count of IndexAccessorMock.ForPNAndJetFunc invocations
func (m *IndexAccessorMock) ForPNAndJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNAndJetCounter)
}

// ForPNAndJetMinimockPreCounter returns the value of IndexAccessorMock.ForPNAndJet invocations
func (m *IndexAccessorMock) ForPNAndJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNAndJetPreCounter)
}

// ForPNAndJetFinished returns true if mock invocations count is ok
func (m *IndexAccessorMock) ForPNAndJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPNAndJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPNAndJetCounter) == uint64(len(m.ForPNAndJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPNAndJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPNAndJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPNAndJetFunc != nil {
		return atomic.LoadUint64(&m.ForPNAndJetCounter) > 0
	}

	return true
}

type mIndexAccessorMockIndex struct {
	mock              *IndexAccessorMock
	mainExpectation   *IndexAccessorMockIndexExpectation
	expectationSeries []*IndexAccessorMockIndexExpectation
}

type IndexAccessorMockIndexExpectation struct {
	input  *IndexAccessorMockIndexInput
	result *IndexAccessorMockIndexResult
}

type IndexAccessorMockIndexInput struct {
	p  insolar.PulseNumber
	p1 insolar.ID
}

type IndexAccessorMockIndexResult struct {
	r *LockedIndex
}

// Expect specifies that invocation of IndexAccessor.Index is expected from 1 to Infinity times
func (m *mIndexAccessorMockIndex) Expect(p insolar.PulseNumber, p1 insolar.ID) *mIndexAccessorMockIndex {
	m.mock.IndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockIndexExpectation{}
	}
	m.mainExpectation.input = &IndexAccessorMockIndexInput{p, p1}
	return m
}

// Return specifies results of invocation of IndexAccessor.Index
func (m *mIndexAccessorMockIndex) Return(r *LockedIndex) *IndexAccessorMock {
	m.mock.IndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockIndexExpectation{}
	}
	m.mainExpectation.result = &IndexAccessorMockIndexResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexAccessor.Index is expected once
func (m *mIndexAccessorMockIndex) ExpectOnce(p insolar.PulseNumber, p1 insolar.ID) *IndexAccessorMockIndexExpectation {
	m.mock.IndexFunc = nil
	m.mainExpectation = nil

	expectation := &IndexAccessorMockIndexExpectation{}
	expectation.input = &IndexAccessorMockIndexInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexAccessorMockIndexExpectation) Return(r *LockedIndex) {
	e.result = &IndexAccessorMockIndexResult{r}
}

// Set uses given function f as a mock of IndexAccessor.Index method
func (m *mIndexAccessorMockIndex) Set(f func(p insolar.PulseNumber, p1 insolar.ID) (r *LockedIndex)) *IndexAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IndexFunc = f
	return m.mock
}

// Index implements github.com/insolar/insolar/ledger/object.IndexAccessor interface
func (m *IndexAccessorMock) Index(p insolar.PulseNumber, p1 insolar.ID) (r *LockedIndex) {
	counter := atomic.AddUint64(&m.IndexPreCounter, 1)
	defer atomic.AddUint64(&m.IndexCounter, 1)

	if len(m.IndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexAccessorMock.Index. %v %v", p, p1)
			return
		}

		input := m.IndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexAccessorMockIndexInput{p, p1}, "IndexAccessor.Index got unexpected parameters")

		result := m.IndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.Index")
			return
		}

		r = result.r

		return
	}

	if m.IndexMock.mainExpectation != nil {

		input := m.IndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexAccessorMockIndexInput{p, p1}, "IndexAccessor.Index got unexpected parameters")
		}

		result := m.IndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.Index")
		}

		r = result.r

		return
	}

	if m.IndexFunc == nil {
		m.t.Fatalf("Unexpected call to IndexAccessorMock.Index. %v %v", p, p1)
		return
	}

	return m.IndexFunc(p, p1)
}

// IndexMinimockCounter returns a count of IndexAccessorMock.IndexFunc invocations
func (m *IndexAccessorMock) IndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IndexCounter)
}

// IndexMinimockPreCounter returns the value of IndexAccessorMock.Index invocations
func (m *IndexAccessorMock) IndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IndexPreCounter)
}

// IndexFinished returns true if mock invocations count is ok
func (m *IndexAccessorMock) IndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IndexCounter) == uint64(len(m.IndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IndexFunc != nil {
		return atomic.LoadUint64(&m.IndexCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexAccessorMock) ValidateCallCounters() {

	if !m.ForPNAndJetFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.ForPNAndJet")
	}

	if !m.IndexFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.Index")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexAccessorMock) MinimockFinish() {

	if !m.ForPNAndJetFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.ForPNAndJet")
	}

	if !m.IndexFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.Index")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPNAndJetFinished()
		ok = ok && m.IndexFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPNAndJetFinished() {
				m.t.Error("Expected call to IndexAccessorMock.ForPNAndJet")
			}

			if !m.IndexFinished() {
				m.t.Error("Expected call to IndexAccessorMock.Index")
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
func (m *IndexAccessorMock) AllMocksCalled() bool {

	if !m.ForPNAndJetFinished() {
		return false
	}

	if !m.IndexFinished() {
		return false
	}

	return true
}

package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexBucketAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexBucketAccessorMock implements github.com/insolar/insolar/ledger/object.IndexBucketAccessor
type IndexBucketAccessorMock struct {
	t minimock.Tester

	ForPNAndJetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []IndexBucket)
	ForPNAndJetCounter    uint64
	ForPNAndJetPreCounter uint64
	ForPNAndJetMock       mIndexBucketAccessorMockForPNAndJet
}

// NewIndexBucketAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.IndexBucketAccessor
func NewIndexBucketAccessorMock(t minimock.Tester) *IndexBucketAccessorMock {
	m := &IndexBucketAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPNAndJetMock = mIndexBucketAccessorMockForPNAndJet{mock: m}

	return m
}

type mIndexBucketAccessorMockForPNAndJet struct {
	mock              *IndexBucketAccessorMock
	mainExpectation   *IndexBucketAccessorMockForPNAndJetExpectation
	expectationSeries []*IndexBucketAccessorMockForPNAndJetExpectation
}

type IndexBucketAccessorMockForPNAndJetExpectation struct {
	input  *IndexBucketAccessorMockForPNAndJetInput
	result *IndexBucketAccessorMockForPNAndJetResult
}

type IndexBucketAccessorMockForPNAndJetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type IndexBucketAccessorMockForPNAndJetResult struct {
	r []IndexBucket
}

// Expect specifies that invocation of IndexBucketAccessor.ForPNAndJet is expected from 1 to Infinity times
func (m *mIndexBucketAccessorMockForPNAndJet) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mIndexBucketAccessorMockForPNAndJet {
	m.mock.ForPNAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexBucketAccessorMockForPNAndJetExpectation{}
	}
	m.mainExpectation.input = &IndexBucketAccessorMockForPNAndJetInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of IndexBucketAccessor.ForPNAndJet
func (m *mIndexBucketAccessorMockForPNAndJet) Return(r []IndexBucket) *IndexBucketAccessorMock {
	m.mock.ForPNAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexBucketAccessorMockForPNAndJetExpectation{}
	}
	m.mainExpectation.result = &IndexBucketAccessorMockForPNAndJetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexBucketAccessor.ForPNAndJet is expected once
func (m *mIndexBucketAccessorMockForPNAndJet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *IndexBucketAccessorMockForPNAndJetExpectation {
	m.mock.ForPNAndJetFunc = nil
	m.mainExpectation = nil

	expectation := &IndexBucketAccessorMockForPNAndJetExpectation{}
	expectation.input = &IndexBucketAccessorMockForPNAndJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexBucketAccessorMockForPNAndJetExpectation) Return(r []IndexBucket) {
	e.result = &IndexBucketAccessorMockForPNAndJetResult{r}
}

// Set uses given function f as a mock of IndexBucketAccessor.ForPNAndJet method
func (m *mIndexBucketAccessorMockForPNAndJet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []IndexBucket)) *IndexBucketAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPNAndJetFunc = f
	return m.mock
}

// ForPNAndJet implements github.com/insolar/insolar/ledger/object.IndexBucketAccessor interface
func (m *IndexBucketAccessorMock) ForPNAndJet(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []IndexBucket) {
	counter := atomic.AddUint64(&m.ForPNAndJetPreCounter, 1)
	defer atomic.AddUint64(&m.ForPNAndJetCounter, 1)

	if len(m.ForPNAndJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPNAndJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexBucketAccessorMock.ForPNAndJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPNAndJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexBucketAccessorMockForPNAndJetInput{p, p1, p2}, "IndexBucketAccessor.ForPNAndJet got unexpected parameters")

		result := m.ForPNAndJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexBucketAccessorMock.ForPNAndJet")
			return
		}

		r = result.r

		return
	}

	if m.ForPNAndJetMock.mainExpectation != nil {

		input := m.ForPNAndJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexBucketAccessorMockForPNAndJetInput{p, p1, p2}, "IndexBucketAccessor.ForPNAndJet got unexpected parameters")
		}

		result := m.ForPNAndJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexBucketAccessorMock.ForPNAndJet")
		}

		r = result.r

		return
	}

	if m.ForPNAndJetFunc == nil {
		m.t.Fatalf("Unexpected call to IndexBucketAccessorMock.ForPNAndJet. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPNAndJetFunc(p, p1, p2)
}

// ForPNAndJetMinimockCounter returns a count of IndexBucketAccessorMock.ForPNAndJetFunc invocations
func (m *IndexBucketAccessorMock) ForPNAndJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNAndJetCounter)
}

// ForPNAndJetMinimockPreCounter returns the value of IndexBucketAccessorMock.ForPNAndJet invocations
func (m *IndexBucketAccessorMock) ForPNAndJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNAndJetPreCounter)
}

// ForPNAndJetFinished returns true if mock invocations count is ok
func (m *IndexBucketAccessorMock) ForPNAndJetFinished() bool {
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

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexBucketAccessorMock) ValidateCallCounters() {

	if !m.ForPNAndJetFinished() {
		m.t.Fatal("Expected call to IndexBucketAccessorMock.ForPNAndJet")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexBucketAccessorMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexBucketAccessorMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexBucketAccessorMock) MinimockFinish() {

	if !m.ForPNAndJetFinished() {
		m.t.Fatal("Expected call to IndexBucketAccessorMock.ForPNAndJet")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexBucketAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexBucketAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPNAndJetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPNAndJetFinished() {
				m.t.Error("Expected call to IndexBucketAccessorMock.ForPNAndJet")
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
func (m *IndexBucketAccessorMock) AllMocksCalled() bool {

	if !m.ForPNAndJetFinished() {
		return false
	}

	return true
}

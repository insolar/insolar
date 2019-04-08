package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexCleaner" can be found in github.com/insolar/insolar/ledger/storage/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexCleanerMock implements github.com/insolar/insolar/ledger/storage/object.IndexCleaner
type IndexCleanerMock struct {
	t minimock.Tester

	RemoveWithIDsFunc func(p context.Context, p1 map[insolar.ID]struct {
	})
	RemoveWithIDsCounter    uint64
	RemoveWithIDsPreCounter uint64
	RemoveWithIDsMock       mIndexCleanerMockRemoveWithIDs
}

// NewIndexCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage/object.IndexCleaner
func NewIndexCleanerMock(t minimock.Tester) *IndexCleanerMock {
	m := &IndexCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveWithIDsMock = mIndexCleanerMockRemoveWithIDs{mock: m}

	return m
}

type mIndexCleanerMockRemoveWithIDs struct {
	mock              *IndexCleanerMock
	mainExpectation   *IndexCleanerMockRemoveWithIDsExpectation
	expectationSeries []*IndexCleanerMockRemoveWithIDsExpectation
}

type IndexCleanerMockRemoveWithIDsExpectation struct {
	input *IndexCleanerMockRemoveWithIDsInput
}

type IndexCleanerMockRemoveWithIDsInput struct {
	p  context.Context
	p1 map[insolar.ID]struct {
	}
}

// Expect specifies that invocation of IndexCleaner.RemoveWithIDs is expected from 1 to Infinity times
func (m *mIndexCleanerMockRemoveWithIDs) Expect(p context.Context, p1 map[insolar.ID]struct {
}) *mIndexCleanerMockRemoveWithIDs {
	m.mock.RemoveWithIDsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCleanerMockRemoveWithIDsExpectation{}
	}
	m.mainExpectation.input = &IndexCleanerMockRemoveWithIDsInput{p, p1}
	return m
}

// Return specifies results of invocation of IndexCleaner.RemoveWithIDs
func (m *mIndexCleanerMockRemoveWithIDs) Return() *IndexCleanerMock {
	m.mock.RemoveWithIDsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCleanerMockRemoveWithIDsExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of IndexCleaner.RemoveWithIDs is expected once
func (m *mIndexCleanerMockRemoveWithIDs) ExpectOnce(p context.Context, p1 map[insolar.ID]struct {
}) *IndexCleanerMockRemoveWithIDsExpectation {
	m.mock.RemoveWithIDsFunc = nil
	m.mainExpectation = nil

	expectation := &IndexCleanerMockRemoveWithIDsExpectation{}
	expectation.input = &IndexCleanerMockRemoveWithIDsInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of IndexCleaner.RemoveWithIDs method
func (m *mIndexCleanerMockRemoveWithIDs) Set(f func(p context.Context, p1 map[insolar.ID]struct {
})) *IndexCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveWithIDsFunc = f
	return m.mock
}

// RemoveWithIDs implements github.com/insolar/insolar/ledger/storage/object.IndexCleaner interface
func (m *IndexCleanerMock) RemoveWithIDs(p context.Context, p1 map[insolar.ID]struct {
}) {
	counter := atomic.AddUint64(&m.RemoveWithIDsPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveWithIDsCounter, 1)

	if len(m.RemoveWithIDsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveWithIDsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexCleanerMock.RemoveWithIDs. %v %v", p, p1)
			return
		}

		input := m.RemoveWithIDsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexCleanerMockRemoveWithIDsInput{p, p1}, "IndexCleaner.RemoveWithIDs got unexpected parameters")

		return
	}

	if m.RemoveWithIDsMock.mainExpectation != nil {

		input := m.RemoveWithIDsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexCleanerMockRemoveWithIDsInput{p, p1}, "IndexCleaner.RemoveWithIDs got unexpected parameters")
		}

		return
	}

	if m.RemoveWithIDsFunc == nil {
		m.t.Fatalf("Unexpected call to IndexCleanerMock.RemoveWithIDs. %v %v", p, p1)
		return
	}

	m.RemoveWithIDsFunc(p, p1)
}

// RemoveWithIDsMinimockCounter returns a count of IndexCleanerMock.RemoveWithIDsFunc invocations
func (m *IndexCleanerMock) RemoveWithIDsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveWithIDsCounter)
}

// RemoveWithIDsMinimockPreCounter returns the value of IndexCleanerMock.RemoveWithIDs invocations
func (m *IndexCleanerMock) RemoveWithIDsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveWithIDsPreCounter)
}

// RemoveWithIDsFinished returns true if mock invocations count is ok
func (m *IndexCleanerMock) RemoveWithIDsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveWithIDsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveWithIDsCounter) == uint64(len(m.RemoveWithIDsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveWithIDsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveWithIDsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveWithIDsFunc != nil {
		return atomic.LoadUint64(&m.RemoveWithIDsCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCleanerMock) ValidateCallCounters() {

	if !m.RemoveWithIDsFinished() {
		m.t.Fatal("Expected call to IndexCleanerMock.RemoveWithIDs")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCleanerMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexCleanerMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexCleanerMock) MinimockFinish() {

	if !m.RemoveWithIDsFinished() {
		m.t.Fatal("Expected call to IndexCleanerMock.RemoveWithIDs")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexCleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexCleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RemoveWithIDsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveWithIDsFinished() {
				m.t.Error("Expected call to IndexCleanerMock.RemoveWithIDs")
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
func (m *IndexCleanerMock) AllMocksCalled() bool {

	if !m.RemoveWithIDsFinished() {
		return false
	}

	return true
}

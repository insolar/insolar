package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineStorage" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// LifelineStorageMock implements github.com/insolar/insolar/ledger/object.LifelineStorage
type LifelineStorageMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.ID) (r Lifeline, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mLifelineStorageMockForID

	SetFunc       func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mLifelineStorageMockSet
}

// NewLifelineStorageMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineStorage
func NewLifelineStorageMock(t minimock.Tester) *LifelineStorageMock {
	m := &LifelineStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mLifelineStorageMockForID{mock: m}
	m.SetMock = mLifelineStorageMockSet{mock: m}

	return m
}

type mLifelineStorageMockForID struct {
	mock              *LifelineStorageMock
	mainExpectation   *LifelineStorageMockForIDExpectation
	expectationSeries []*LifelineStorageMockForIDExpectation
}

type LifelineStorageMockForIDExpectation struct {
	input  *LifelineStorageMockForIDInput
	result *LifelineStorageMockForIDResult
}

type LifelineStorageMockForIDInput struct {
	p  context.Context
	p1 insolar.ID
}

type LifelineStorageMockForIDResult struct {
	r  Lifeline
	r1 error
}

// Expect specifies that invocation of LifelineStorage.ForID is expected from 1 to Infinity times
func (m *mLifelineStorageMockForID) Expect(p context.Context, p1 insolar.ID) *mLifelineStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineStorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &LifelineStorageMockForIDInput{p, p1}
	return m
}

// Return specifies results of invocation of LifelineStorage.ForID
func (m *mLifelineStorageMockForID) Return(r Lifeline, r1 error) *LifelineStorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineStorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &LifelineStorageMockForIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of LifelineStorage.ForID is expected once
func (m *mLifelineStorageMockForID) ExpectOnce(p context.Context, p1 insolar.ID) *LifelineStorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineStorageMockForIDExpectation{}
	expectation.input = &LifelineStorageMockForIDInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineStorageMockForIDExpectation) Return(r Lifeline, r1 error) {
	e.result = &LifelineStorageMockForIDResult{r, r1}
}

// Set uses given function f as a mock of LifelineStorage.ForID method
func (m *mLifelineStorageMockForID) Set(f func(p context.Context, p1 insolar.ID) (r Lifeline, r1 error)) *LifelineStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

// ForID implements github.com/insolar/insolar/ledger/object.LifelineStorage interface
func (m *LifelineStorageMock) ForID(p context.Context, p1 insolar.ID) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineStorageMock.ForID. %v %v", p, p1)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineStorageMockForIDInput{p, p1}, "LifelineStorage.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineStorageMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineStorageMockForIDInput{p, p1}, "LifelineStorage.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineStorageMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineStorageMock.ForID. %v %v", p, p1)
		return
	}

	return m.ForIDFunc(p, p1)
}

// ForIDMinimockCounter returns a count of LifelineStorageMock.ForIDFunc invocations
func (m *LifelineStorageMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

// ForIDMinimockPreCounter returns the value of LifelineStorageMock.ForID invocations
func (m *LifelineStorageMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

// ForIDFinished returns true if mock invocations count is ok
func (m *LifelineStorageMock) ForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForIDCounter) == uint64(len(m.ForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForIDFunc != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	return true
}

type mLifelineStorageMockSet struct {
	mock              *LifelineStorageMock
	mainExpectation   *LifelineStorageMockSetExpectation
	expectationSeries []*LifelineStorageMockSetExpectation
}

type LifelineStorageMockSetExpectation struct {
	input  *LifelineStorageMockSetInput
	result *LifelineStorageMockSetResult
}

type LifelineStorageMockSetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 Lifeline
}

type LifelineStorageMockSetResult struct {
	r error
}

// Expect specifies that invocation of LifelineStorage.Set is expected from 1 to Infinity times
func (m *mLifelineStorageMockSet) Expect(p context.Context, p1 insolar.ID, p2 Lifeline) *mLifelineStorageMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineStorageMockSetExpectation{}
	}
	m.mainExpectation.input = &LifelineStorageMockSetInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of LifelineStorage.Set
func (m *mLifelineStorageMockSet) Return(r error) *LifelineStorageMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineStorageMockSetExpectation{}
	}
	m.mainExpectation.result = &LifelineStorageMockSetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of LifelineStorage.Set is expected once
func (m *mLifelineStorageMockSet) ExpectOnce(p context.Context, p1 insolar.ID, p2 Lifeline) *LifelineStorageMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineStorageMockSetExpectation{}
	expectation.input = &LifelineStorageMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineStorageMockSetExpectation) Return(r error) {
	e.result = &LifelineStorageMockSetResult{r}
}

// Set uses given function f as a mock of LifelineStorage.Set method
func (m *mLifelineStorageMockSet) Set(f func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)) *LifelineStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

// Set implements github.com/insolar/insolar/ledger/object.LifelineStorage interface
func (m *LifelineStorageMock) Set(p context.Context, p1 insolar.ID, p2 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineStorageMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineStorageMockSetInput{p, p1, p2}, "LifelineStorage.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineStorageMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineStorageMockSetInput{p, p1, p2}, "LifelineStorage.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineStorageMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineStorageMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

// SetMinimockCounter returns a count of LifelineStorageMock.SetFunc invocations
func (m *LifelineStorageMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

// SetMinimockPreCounter returns the value of LifelineStorageMock.Set invocations
func (m *LifelineStorageMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

// SetFinished returns true if mock invocations count is ok
func (m *LifelineStorageMock) SetFinished() bool {
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
func (m *LifelineStorageMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to LifelineStorageMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to LifelineStorageMock.Set")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineStorageMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LifelineStorageMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LifelineStorageMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to LifelineStorageMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to LifelineStorageMock.Set")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LifelineStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *LifelineStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to LifelineStorageMock.ForID")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to LifelineStorageMock.Set")
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
func (m *LifelineStorageMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}

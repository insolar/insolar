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

//IndexAccessorMock implements github.com/insolar/insolar/ledger/object.IndexAccessor
type IndexAccessorMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r FilamentIndex, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mIndexAccessorMockForID

	ForPulseFunc       func(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mIndexAccessorMockForPulse
}

//NewIndexAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.IndexAccessor
func NewIndexAccessorMock(t minimock.Tester) *IndexAccessorMock {
	m := &IndexAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mIndexAccessorMockForID{mock: m}
	m.ForPulseMock = mIndexAccessorMockForPulse{mock: m}

	return m
}

type mIndexAccessorMockForID struct {
	mock              *IndexAccessorMock
	mainExpectation   *IndexAccessorMockForIDExpectation
	expectationSeries []*IndexAccessorMockForIDExpectation
}

type IndexAccessorMockForIDExpectation struct {
	input  *IndexAccessorMockForIDInput
	result *IndexAccessorMockForIDResult
}

type IndexAccessorMockForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type IndexAccessorMockForIDResult struct {
	r  FilamentIndex
	r1 error
}

//Expect specifies that invocation of IndexAccessor.ForID is expected from 1 to Infinity times
func (m *mIndexAccessorMockForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mIndexAccessorMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockForIDExpectation{}
	}
	m.mainExpectation.input = &IndexAccessorMockForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexAccessor.ForID
func (m *mIndexAccessorMockForID) Return(r FilamentIndex, r1 error) *IndexAccessorMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockForIDExpectation{}
	}
	m.mainExpectation.result = &IndexAccessorMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexAccessor.ForID is expected once
func (m *mIndexAccessorMockForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *IndexAccessorMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &IndexAccessorMockForIDExpectation{}
	expectation.input = &IndexAccessorMockForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexAccessorMockForIDExpectation) Return(r FilamentIndex, r1 error) {
	e.result = &IndexAccessorMockForIDResult{r, r1}
}

//Set uses given function f as a mock of IndexAccessor.ForID method
func (m *mIndexAccessorMockForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r FilamentIndex, r1 error)) *IndexAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.IndexAccessor interface
func (m *IndexAccessorMock) ForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r FilamentIndex, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexAccessorMock.ForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexAccessorMockForIDInput{p, p1, p2}, "IndexAccessor.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexAccessorMockForIDInput{p, p1, p2}, "IndexAccessor.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to IndexAccessorMock.ForID. %v %v %v", p, p1, p2)
		return
	}

	return m.ForIDFunc(p, p1, p2)
}

//ForIDMinimockCounter returns a count of IndexAccessorMock.ForIDFunc invocations
func (m *IndexAccessorMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of IndexAccessorMock.ForID invocations
func (m *IndexAccessorMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *IndexAccessorMock) ForIDFinished() bool {
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

type mIndexAccessorMockForPulse struct {
	mock              *IndexAccessorMock
	mainExpectation   *IndexAccessorMockForPulseExpectation
	expectationSeries []*IndexAccessorMockForPulseExpectation
}

type IndexAccessorMockForPulseExpectation struct {
	input  *IndexAccessorMockForPulseInput
	result *IndexAccessorMockForPulseResult
}

type IndexAccessorMockForPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type IndexAccessorMockForPulseResult struct {
	r []FilamentIndex
}

// Expect specifies that invocation of IndexAccessor.ForPulse is expected from 1 to Infinity times
func (m *mIndexAccessorMockForPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mIndexAccessorMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.input = &IndexAccessorMockForPulseInput{p, p1}
	return m
}

// Return specifies results of invocation of IndexAccessor.ForPulse
func (m *mIndexAccessorMockForPulse) Return(r []FilamentIndex) *IndexAccessorMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.result = &IndexAccessorMockForPulseResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of IndexAccessor.ForPulse is expected once
func (m *mIndexAccessorMockForPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *IndexAccessorMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &IndexAccessorMockForPulseExpectation{}
	expectation.input = &IndexAccessorMockForPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexAccessorMockForPulseExpectation) Return(r []FilamentIndex) {
	e.result = &IndexAccessorMockForPulseResult{r}
}

// Set uses given function f as a mock of IndexAccessor.ForPulse method
func (m *mIndexAccessorMockForPulse) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex)) *IndexAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

// ForPulse implements github.com/insolar/insolar/ledger/object.IndexAccessor interface
func (m *IndexAccessorMock) ForPulse(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexAccessorMock.ForPulse. %v %v", p, p1)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexAccessorMockForPulseInput{p, p1}, "IndexAccessor.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.ForPulse")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexAccessorMockForPulseInput{p, p1}, "IndexAccessor.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexAccessorMock.ForPulse")
		}

		r = result.r

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to IndexAccessorMock.ForPulse. %v %v", p, p1)
		return
	}

	return m.ForPulseFunc(p, p1)
}

// ForPulseMinimockCounter returns a count of IndexAccessorMock.ForPulseFunc invocations
func (m *IndexAccessorMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

// ForPulseMinimockPreCounter returns the value of IndexAccessorMock.ForPulse invocations
func (m *IndexAccessorMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

// ForPulseFinished returns true if mock invocations count is ok
func (m *IndexAccessorMock) ForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseCounter) == uint64(len(m.ForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseFunc != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexAccessorMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.ForID")
	}

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.ForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexAccessorMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.ForID")
	}

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to IndexAccessorMock.ForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()
		ok = ok && m.ForPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to IndexAccessorMock.ForID")
			}

			if !m.ForPulseFinished() {
				m.t.Error("Expected call to IndexAccessorMock.ForPulse")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *IndexAccessorMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.ForPulseFinished() {
		return false
	}

	return true
}

package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineIndex" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//LifelineIndexMock implements github.com/insolar/insolar/ledger/object.LifelineIndex
type LifelineIndexMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mLifelineIndexMockForID

	SetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mLifelineIndexMockSet
}

//NewLifelineIndexMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineIndex
func NewLifelineIndexMock(t minimock.Tester) *LifelineIndexMock {
	m := &LifelineIndexMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mLifelineIndexMockForID{mock: m}
	m.SetMock = mLifelineIndexMockSet{mock: m}

	return m
}

type mLifelineIndexMockForID struct {
	mock              *LifelineIndexMock
	mainExpectation   *LifelineIndexMockForIDExpectation
	expectationSeries []*LifelineIndexMockForIDExpectation
}

type LifelineIndexMockForIDExpectation struct {
	input  *LifelineIndexMockForIDInput
	result *LifelineIndexMockForIDResult
}

type LifelineIndexMockForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type LifelineIndexMockForIDResult struct {
	r  Lifeline
	r1 error
}

//Expect specifies that invocation of LifelineIndex.ForID is expected from 1 to Infinity times
func (m *mLifelineIndexMockForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mLifelineIndexMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockForIDExpectation{}
	}
	m.mainExpectation.input = &LifelineIndexMockForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of LifelineIndex.ForID
func (m *mLifelineIndexMockForID) Return(r Lifeline, r1 error) *LifelineIndexMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockForIDExpectation{}
	}
	m.mainExpectation.result = &LifelineIndexMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LifelineIndex.ForID is expected once
func (m *mLifelineIndexMockForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *LifelineIndexMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineIndexMockForIDExpectation{}
	expectation.input = &LifelineIndexMockForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineIndexMockForIDExpectation) Return(r Lifeline, r1 error) {
	e.result = &LifelineIndexMockForIDResult{r, r1}
}

//Set uses given function f as a mock of LifelineIndex.ForID method
func (m *mLifelineIndexMockForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)) *LifelineIndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.LifelineIndex interface
func (m *LifelineIndexMock) ForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineIndexMock.ForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineIndexMockForIDInput{p, p1, p2}, "LifelineIndex.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineIndexMockForIDInput{p, p1, p2}, "LifelineIndex.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineIndexMock.ForID. %v %v %v", p, p1, p2)
		return
	}

	return m.ForIDFunc(p, p1, p2)
}

//ForIDMinimockCounter returns a count of LifelineIndexMock.ForIDFunc invocations
func (m *LifelineIndexMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of LifelineIndexMock.ForID invocations
func (m *LifelineIndexMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *LifelineIndexMock) ForIDFinished() bool {
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

type mLifelineIndexMockSet struct {
	mock              *LifelineIndexMock
	mainExpectation   *LifelineIndexMockSetExpectation
	expectationSeries []*LifelineIndexMockSetExpectation
}

type LifelineIndexMockSetExpectation struct {
	input  *LifelineIndexMockSetInput
	result *LifelineIndexMockSetResult
}

type LifelineIndexMockSetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 Lifeline
}

type LifelineIndexMockSetResult struct {
	r error
}

//Expect specifies that invocation of LifelineIndex.Set is expected from 1 to Infinity times
func (m *mLifelineIndexMockSet) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *mLifelineIndexMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockSetExpectation{}
	}
	m.mainExpectation.input = &LifelineIndexMockSetInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of LifelineIndex.Set
func (m *mLifelineIndexMockSet) Return(r error) *LifelineIndexMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockSetExpectation{}
	}
	m.mainExpectation.result = &LifelineIndexMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LifelineIndex.Set is expected once
func (m *mLifelineIndexMockSet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *LifelineIndexMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineIndexMockSetExpectation{}
	expectation.input = &LifelineIndexMockSetInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineIndexMockSetExpectation) Return(r error) {
	e.result = &LifelineIndexMockSetResult{r}
}

//Set uses given function f as a mock of LifelineIndex.Set method
func (m *mLifelineIndexMockSet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)) *LifelineIndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/object.LifelineIndex interface
func (m *LifelineIndexMock) Set(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineIndexMock.Set. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineIndexMockSetInput{p, p1, p2, p3}, "LifelineIndex.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineIndexMockSetInput{p, p1, p2, p3}, "LifelineIndex.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineIndexMock.Set. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetFunc(p, p1, p2, p3)
}

//SetMinimockCounter returns a count of LifelineIndexMock.SetFunc invocations
func (m *LifelineIndexMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of LifelineIndexMock.Set invocations
func (m *LifelineIndexMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *LifelineIndexMock) SetFinished() bool {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineIndexMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineIndexMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LifelineIndexMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LifelineIndexMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LifelineIndexMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LifelineIndexMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to LifelineIndexMock.ForID")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to LifelineIndexMock.Set")
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
func (m *LifelineIndexMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}

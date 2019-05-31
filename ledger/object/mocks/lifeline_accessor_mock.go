package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	object "github.com/insolar/insolar/ledger/object"

	testify_assert "github.com/stretchr/testify/assert"
)

//LifelineAccessorMock implements github.com/insolar/insolar/ledger/object.LifelineAccessor
type LifelineAccessorMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r object.Lifeline, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mLifelineAccessorMockForID
}

//NewLifelineAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineAccessor
func NewLifelineAccessorMock(t minimock.Tester) *LifelineAccessorMock {
	m := &LifelineAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mLifelineAccessorMockForID{mock: m}

	return m
}

type mLifelineAccessorMockForID struct {
	mock              *LifelineAccessorMock
	mainExpectation   *LifelineAccessorMockForIDExpectation
	expectationSeries []*LifelineAccessorMockForIDExpectation
}

type LifelineAccessorMockForIDExpectation struct {
	input  *LifelineAccessorMockForIDInput
	result *LifelineAccessorMockForIDResult
}

type LifelineAccessorMockForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type LifelineAccessorMockForIDResult struct {
	r  object.Lifeline
	r1 error
}

//Expect specifies that invocation of LifelineAccessor.ForID is expected from 1 to Infinity times
func (m *mLifelineAccessorMockForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mLifelineAccessorMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineAccessorMockForIDExpectation{}
	}
	m.mainExpectation.input = &LifelineAccessorMockForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of LifelineAccessor.ForID
func (m *mLifelineAccessorMockForID) Return(r object.Lifeline, r1 error) *LifelineAccessorMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineAccessorMockForIDExpectation{}
	}
	m.mainExpectation.result = &LifelineAccessorMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LifelineAccessor.ForID is expected once
func (m *mLifelineAccessorMockForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *LifelineAccessorMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineAccessorMockForIDExpectation{}
	expectation.input = &LifelineAccessorMockForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineAccessorMockForIDExpectation) Return(r object.Lifeline, r1 error) {
	e.result = &LifelineAccessorMockForIDResult{r, r1}
}

//Set uses given function f as a mock of LifelineAccessor.ForID method
func (m *mLifelineAccessorMockForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r object.Lifeline, r1 error)) *LifelineAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.LifelineAccessor interface
func (m *LifelineAccessorMock) ForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r object.Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineAccessorMock.ForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineAccessorMockForIDInput{p, p1, p2}, "LifelineAccessor.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineAccessorMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineAccessorMockForIDInput{p, p1, p2}, "LifelineAccessor.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineAccessorMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineAccessorMock.ForID. %v %v %v", p, p1, p2)
		return
	}

	return m.ForIDFunc(p, p1, p2)
}

//ForIDMinimockCounter returns a count of LifelineAccessorMock.ForIDFunc invocations
func (m *LifelineAccessorMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of LifelineAccessorMock.ForID invocations
func (m *LifelineAccessorMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *LifelineAccessorMock) ForIDFinished() bool {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineAccessorMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to LifelineAccessorMock.ForID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LifelineAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LifelineAccessorMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to LifelineAccessorMock.ForID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LifelineAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LifelineAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to LifelineAccessorMock.ForID")
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
func (m *LifelineAccessorMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	return true
}

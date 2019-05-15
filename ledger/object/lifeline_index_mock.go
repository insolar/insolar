package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineIndex" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//LifelineIndexMock implements github.com/insolar/insolar/ledger/object.LifelineIndex
type LifelineIndexMock struct {
	t minimock.Tester

	LifelineForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)
	LifelineForIDCounter    uint64
	LifelineForIDPreCounter uint64
	LifelineForIDMock       mLifelineIndexMockLifelineForID

	SetLifelineFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)
	SetLifelineCounter    uint64
	SetLifelinePreCounter uint64
	SetLifelineMock       mLifelineIndexMockSetLifeline
}

//NewLifelineIndexMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineIndex
func NewLifelineIndexMock(t minimock.Tester) *LifelineIndexMock {
	m := &LifelineIndexMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LifelineForIDMock = mLifelineIndexMockLifelineForID{mock: m}
	m.SetLifelineMock = mLifelineIndexMockSetLifeline{mock: m}

	return m
}

type mLifelineIndexMockLifelineForID struct {
	mock              *LifelineIndexMock
	mainExpectation   *LifelineIndexMockLifelineForIDExpectation
	expectationSeries []*LifelineIndexMockLifelineForIDExpectation
}

type LifelineIndexMockLifelineForIDExpectation struct {
	input  *LifelineIndexMockLifelineForIDInput
	result *LifelineIndexMockLifelineForIDResult
}

type LifelineIndexMockLifelineForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type LifelineIndexMockLifelineForIDResult struct {
	r  Lifeline
	r1 error
}

//Expect specifies that invocation of LifelineIndex.LifelineForID is expected from 1 to Infinity times
func (m *mLifelineIndexMockLifelineForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mLifelineIndexMockLifelineForID {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockLifelineForIDExpectation{}
	}
	m.mainExpectation.input = &LifelineIndexMockLifelineForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of LifelineIndex.LifelineForID
func (m *mLifelineIndexMockLifelineForID) Return(r Lifeline, r1 error) *LifelineIndexMock {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockLifelineForIDExpectation{}
	}
	m.mainExpectation.result = &LifelineIndexMockLifelineForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LifelineIndex.LifelineForID is expected once
func (m *mLifelineIndexMockLifelineForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *LifelineIndexMockLifelineForIDExpectation {
	m.mock.LifelineForIDFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineIndexMockLifelineForIDExpectation{}
	expectation.input = &LifelineIndexMockLifelineForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineIndexMockLifelineForIDExpectation) Return(r Lifeline, r1 error) {
	e.result = &LifelineIndexMockLifelineForIDResult{r, r1}
}

//Set uses given function f as a mock of LifelineIndex.LifelineForID method
func (m *mLifelineIndexMockLifelineForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)) *LifelineIndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LifelineForIDFunc = f
	return m.mock
}

//LifelineForID implements github.com/insolar/insolar/ledger/object.LifelineIndex interface
func (m *LifelineIndexMock) LifelineForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.LifelineForIDPreCounter, 1)
	defer atomic.AddUint64(&m.LifelineForIDCounter, 1)

	if len(m.LifelineForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LifelineForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineIndexMock.LifelineForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.LifelineForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineIndexMockLifelineForIDInput{p, p1, p2}, "LifelineIndex.LifelineForID got unexpected parameters")

		result := m.LifelineForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.LifelineForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDMock.mainExpectation != nil {

		input := m.LifelineForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineIndexMockLifelineForIDInput{p, p1, p2}, "LifelineIndex.LifelineForID got unexpected parameters")
		}

		result := m.LifelineForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.LifelineForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineIndexMock.LifelineForID. %v %v %v", p, p1, p2)
		return
	}

	return m.LifelineForIDFunc(p, p1, p2)
}

//LifelineForIDMinimockCounter returns a count of LifelineIndexMock.LifelineForIDFunc invocations
func (m *LifelineIndexMock) LifelineForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDCounter)
}

//LifelineForIDMinimockPreCounter returns the value of LifelineIndexMock.LifelineForID invocations
func (m *LifelineIndexMock) LifelineForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDPreCounter)
}

//LifelineForIDFinished returns true if mock invocations count is ok
func (m *LifelineIndexMock) LifelineForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LifelineForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LifelineForIDCounter) == uint64(len(m.LifelineForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LifelineForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LifelineForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LifelineForIDFunc != nil {
		return atomic.LoadUint64(&m.LifelineForIDCounter) > 0
	}

	return true
}

type mLifelineIndexMockSetLifeline struct {
	mock              *LifelineIndexMock
	mainExpectation   *LifelineIndexMockSetLifelineExpectation
	expectationSeries []*LifelineIndexMockSetLifelineExpectation
}

type LifelineIndexMockSetLifelineExpectation struct {
	input  *LifelineIndexMockSetLifelineInput
	result *LifelineIndexMockSetLifelineResult
}

type LifelineIndexMockSetLifelineInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 Lifeline
}

type LifelineIndexMockSetLifelineResult struct {
	r error
}

//Expect specifies that invocation of LifelineIndex.SetLifeline is expected from 1 to Infinity times
func (m *mLifelineIndexMockSetLifeline) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *mLifelineIndexMockSetLifeline {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockSetLifelineExpectation{}
	}
	m.mainExpectation.input = &LifelineIndexMockSetLifelineInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of LifelineIndex.SetLifeline
func (m *mLifelineIndexMockSetLifeline) Return(r error) *LifelineIndexMock {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineIndexMockSetLifelineExpectation{}
	}
	m.mainExpectation.result = &LifelineIndexMockSetLifelineResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LifelineIndex.SetLifeline is expected once
func (m *mLifelineIndexMockSetLifeline) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *LifelineIndexMockSetLifelineExpectation {
	m.mock.SetLifelineFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineIndexMockSetLifelineExpectation{}
	expectation.input = &LifelineIndexMockSetLifelineInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LifelineIndexMockSetLifelineExpectation) Return(r error) {
	e.result = &LifelineIndexMockSetLifelineResult{r}
}

//Set uses given function f as a mock of LifelineIndex.SetLifeline method
func (m *mLifelineIndexMockSetLifeline) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)) *LifelineIndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLifelineFunc = f
	return m.mock
}

//SetLifeline implements github.com/insolar/insolar/ledger/object.LifelineIndex interface
func (m *LifelineIndexMock) SetLifeline(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetLifelinePreCounter, 1)
	defer atomic.AddUint64(&m.SetLifelineCounter, 1)

	if len(m.SetLifelineMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLifelineMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineIndexMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetLifelineMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineIndexMockSetLifelineInput{p, p1, p2, p3}, "LifelineIndex.SetLifeline got unexpected parameters")

		result := m.SetLifelineMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.SetLifeline")
			return
		}

		r = result.r

		return
	}

	if m.SetLifelineMock.mainExpectation != nil {

		input := m.SetLifelineMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineIndexMockSetLifelineInput{p, p1, p2, p3}, "LifelineIndex.SetLifeline got unexpected parameters")
		}

		result := m.SetLifelineMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LifelineIndexMock.SetLifeline")
		}

		r = result.r

		return
	}

	if m.SetLifelineFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineIndexMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetLifelineFunc(p, p1, p2, p3)
}

//SetLifelineMinimockCounter returns a count of LifelineIndexMock.SetLifelineFunc invocations
func (m *LifelineIndexMock) SetLifelineMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineCounter)
}

//SetLifelineMinimockPreCounter returns the value of LifelineIndexMock.SetLifeline invocations
func (m *LifelineIndexMock) SetLifelineMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelinePreCounter)
}

//SetLifelineFinished returns true if mock invocations count is ok
func (m *LifelineIndexMock) SetLifelineFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetLifelineMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetLifelineCounter) == uint64(len(m.SetLifelineMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetLifelineMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetLifelineCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetLifelineFunc != nil {
		return atomic.LoadUint64(&m.SetLifelineCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineIndexMock) ValidateCallCounters() {

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.LifelineForID")
	}

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.SetLifeline")
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

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.LifelineForID")
	}

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to LifelineIndexMock.SetLifeline")
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
		ok = ok && m.LifelineForIDFinished()
		ok = ok && m.SetLifelineFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.LifelineForIDFinished() {
				m.t.Error("Expected call to LifelineIndexMock.LifelineForID")
			}

			if !m.SetLifelineFinished() {
				m.t.Error("Expected call to LifelineIndexMock.SetLifeline")
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

	if !m.LifelineForIDFinished() {
		return false
	}

	if !m.SetLifelineFinished() {
		return false
	}

	return true
}

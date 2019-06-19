package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexHeavyModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexHeavyModifierMock implements github.com/insolar/insolar/ledger/object.IndexHeavyModifier
type IndexHeavyModifierMock struct {
	t minimock.Tester

	SetIndexFunc       func(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) (r error)
	SetIndexCounter    uint64
	SetIndexPreCounter uint64
	SetIndexMock       mIndexHeavyModifierMockSetIndex
}

//NewIndexHeavyModifierMock returns a mock for github.com/insolar/insolar/ledger/object.IndexHeavyModifier
func NewIndexHeavyModifierMock(t minimock.Tester) *IndexHeavyModifierMock {
	m := &IndexHeavyModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetIndexMock = mIndexHeavyModifierMockSetIndex{mock: m}

	return m
}

type mIndexHeavyModifierMockSetIndex struct {
	mock              *IndexHeavyModifierMock
	mainExpectation   *IndexHeavyModifierMockSetIndexExpectation
	expectationSeries []*IndexHeavyModifierMockSetIndexExpectation
}

type IndexHeavyModifierMockSetIndexExpectation struct {
	input  *IndexHeavyModifierMockSetIndexInput
	result *IndexHeavyModifierMockSetIndexResult
}

type IndexHeavyModifierMockSetIndexInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 FilamentIndex
}

type IndexHeavyModifierMockSetIndexResult struct {
	r error
}

//Expect specifies that invocation of IndexHeavyModifier.SetIndex is expected from 1 to Infinity times
func (m *mIndexHeavyModifierMockSetIndex) Expect(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) *mIndexHeavyModifierMockSetIndex {
	m.mock.SetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexHeavyModifierMockSetIndexExpectation{}
	}
	m.mainExpectation.input = &IndexHeavyModifierMockSetIndexInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexHeavyModifier.SetIndex
func (m *mIndexHeavyModifierMockSetIndex) Return(r error) *IndexHeavyModifierMock {
	m.mock.SetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexHeavyModifierMockSetIndexExpectation{}
	}
	m.mainExpectation.result = &IndexHeavyModifierMockSetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexHeavyModifier.SetIndex is expected once
func (m *mIndexHeavyModifierMockSetIndex) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) *IndexHeavyModifierMockSetIndexExpectation {
	m.mock.SetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &IndexHeavyModifierMockSetIndexExpectation{}
	expectation.input = &IndexHeavyModifierMockSetIndexInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexHeavyModifierMockSetIndexExpectation) Return(r error) {
	e.result = &IndexHeavyModifierMockSetIndexResult{r}
}

//Set uses given function f as a mock of IndexHeavyModifier.SetIndex method
func (m *mIndexHeavyModifierMockSetIndex) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) (r error)) *IndexHeavyModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetIndexFunc = f
	return m.mock
}

//SetIndex implements github.com/insolar/insolar/ledger/object.IndexHeavyModifier interface
func (m *IndexHeavyModifierMock) SetIndex(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) (r error) {
	counter := atomic.AddUint64(&m.SetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.SetIndexCounter, 1)

	if len(m.SetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexHeavyModifierMock.SetIndex. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexHeavyModifierMockSetIndexInput{p, p1, p2}, "IndexHeavyModifier.SetIndex got unexpected parameters")

		result := m.SetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexHeavyModifierMock.SetIndex")
			return
		}

		r = result.r

		return
	}

	if m.SetIndexMock.mainExpectation != nil {

		input := m.SetIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexHeavyModifierMockSetIndexInput{p, p1, p2}, "IndexHeavyModifier.SetIndex got unexpected parameters")
		}

		result := m.SetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexHeavyModifierMock.SetIndex")
		}

		r = result.r

		return
	}

	if m.SetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to IndexHeavyModifierMock.SetIndex. %v %v %v", p, p1, p2)
		return
	}

	return m.SetIndexFunc(p, p1, p2)
}

//SetIndexMinimockCounter returns a count of IndexHeavyModifierMock.SetIndexFunc invocations
func (m *IndexHeavyModifierMock) SetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetIndexCounter)
}

//SetIndexMinimockPreCounter returns the value of IndexHeavyModifierMock.SetIndex invocations
func (m *IndexHeavyModifierMock) SetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetIndexPreCounter)
}

//SetIndexFinished returns true if mock invocations count is ok
func (m *IndexHeavyModifierMock) SetIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetIndexCounter) == uint64(len(m.SetIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetIndexFunc != nil {
		return atomic.LoadUint64(&m.SetIndexCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexHeavyModifierMock) ValidateCallCounters() {

	if !m.SetIndexFinished() {
		m.t.Fatal("Expected call to IndexHeavyModifierMock.SetIndex")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexHeavyModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexHeavyModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexHeavyModifierMock) MinimockFinish() {

	if !m.SetIndexFinished() {
		m.t.Fatal("Expected call to IndexHeavyModifierMock.SetIndex")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexHeavyModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexHeavyModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetIndexFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetIndexFinished() {
				m.t.Error("Expected call to IndexHeavyModifierMock.SetIndex")
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
func (m *IndexHeavyModifierMock) AllMocksCalled() bool {

	if !m.SetIndexFinished() {
		return false
	}

	return true
}

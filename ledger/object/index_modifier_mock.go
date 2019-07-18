package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexModifierMock implements github.com/insolar/insolar/ledger/object.IndexModifier
type IndexModifierMock struct {
	t minimock.Tester

	SetIndexFunc       func(p context.Context, p1 insolar.PulseNumber, p2 record.Index) (r error)
	SetIndexCounter    uint64
	SetIndexPreCounter uint64
	SetIndexMock       mIndexModifierMockSetIndex
}

//NewIndexModifierMock returns a mock for github.com/insolar/insolar/ledger/object.IndexModifier
func NewIndexModifierMock(t minimock.Tester) *IndexModifierMock {
	m := &IndexModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetIndexMock = mIndexModifierMockSetIndex{mock: m}

	return m
}

type mIndexModifierMockSetIndex struct {
	mock              *IndexModifierMock
	mainExpectation   *IndexModifierMockSetIndexExpectation
	expectationSeries []*IndexModifierMockSetIndexExpectation
}

type IndexModifierMockSetIndexExpectation struct {
	input  *IndexModifierMockSetIndexInput
	result *IndexModifierMockSetIndexResult
}

type IndexModifierMockSetIndexInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 record.Index
}

type IndexModifierMockSetIndexResult struct {
	r error
}

//Expect specifies that invocation of IndexModifier.SetIndex is expected from 1 to Infinity times
func (m *mIndexModifierMockSetIndex) Expect(p context.Context, p1 insolar.PulseNumber, p2 record.Index) *mIndexModifierMockSetIndex {
	m.mock.SetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetIndexExpectation{}
	}
	m.mainExpectation.input = &IndexModifierMockSetIndexInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexModifier.SetIndex
func (m *mIndexModifierMockSetIndex) Return(r error) *IndexModifierMock {
	m.mock.SetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetIndexExpectation{}
	}
	m.mainExpectation.result = &IndexModifierMockSetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexModifier.SetIndex is expected once
func (m *mIndexModifierMockSetIndex) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 record.Index) *IndexModifierMockSetIndexExpectation {
	m.mock.SetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &IndexModifierMockSetIndexExpectation{}
	expectation.input = &IndexModifierMockSetIndexInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexModifierMockSetIndexExpectation) Return(r error) {
	e.result = &IndexModifierMockSetIndexResult{r}
}

//Set uses given function f as a mock of IndexModifier.SetIndex method
func (m *mIndexModifierMockSetIndex) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 record.Index) (r error)) *IndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetIndexFunc = f
	return m.mock
}

//SetIndex implements github.com/insolar/insolar/ledger/object.IndexModifier interface
func (m *IndexModifierMock) SetIndex(p context.Context, p1 insolar.PulseNumber, p2 record.Index) (r error) {
	counter := atomic.AddUint64(&m.SetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.SetIndexCounter, 1)

	if len(m.SetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexModifierMock.SetIndex. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexModifierMockSetIndexInput{p, p1, p2}, "IndexModifier.SetIndex got unexpected parameters")

		result := m.SetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetIndex")
			return
		}

		r = result.r

		return
	}

	if m.SetIndexMock.mainExpectation != nil {

		input := m.SetIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexModifierMockSetIndexInput{p, p1, p2}, "IndexModifier.SetIndex got unexpected parameters")
		}

		result := m.SetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetIndex")
		}

		r = result.r

		return
	}

	if m.SetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to IndexModifierMock.SetIndex. %v %v %v", p, p1, p2)
		return
	}

	return m.SetIndexFunc(p, p1, p2)
}

//SetIndexMinimockCounter returns a count of IndexModifierMock.SetIndexFunc invocations
func (m *IndexModifierMock) SetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetIndexCounter)
}

//SetIndexMinimockPreCounter returns the value of IndexModifierMock.SetIndex invocations
func (m *IndexModifierMock) SetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetIndexPreCounter)
}

//SetIndexFinished returns true if mock invocations count is ok
func (m *IndexModifierMock) SetIndexFinished() bool {
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
func (m *IndexModifierMock) ValidateCallCounters() {

	if !m.SetIndexFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetIndex")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexModifierMock) MinimockFinish() {

	if !m.SetIndexFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetIndex")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexModifierMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IndexModifierMock.SetIndex")
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
func (m *IndexModifierMock) AllMocksCalled() bool {

	if !m.SetIndexFinished() {
		return false
	}

	return true
}

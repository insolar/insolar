package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexSaver" can be found in github.com/insolar/insolar/ledger/storage/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexSaverMock implements github.com/insolar/insolar/ledger/storage/object.IndexSaver
type IndexSaverMock struct {
	t minimock.Tester

	SaveIndexFromHeavyFunc       func(p context.Context, p1 insolar.ID, p2 insolar.Reference, p3 *insolar.Reference) (r Lifeline, r1 error)
	SaveIndexFromHeavyCounter    uint64
	SaveIndexFromHeavyPreCounter uint64
	SaveIndexFromHeavyMock       mIndexSaverMockSaveIndexFromHeavy
}

//NewIndexSaverMock returns a mock for github.com/insolar/insolar/ledger/storage/object.IndexSaver
func NewIndexSaverMock(t minimock.Tester) *IndexSaverMock {
	m := &IndexSaverMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SaveIndexFromHeavyMock = mIndexSaverMockSaveIndexFromHeavy{mock: m}

	return m
}

type mIndexSaverMockSaveIndexFromHeavy struct {
	mock              *IndexSaverMock
	mainExpectation   *IndexSaverMockSaveIndexFromHeavyExpectation
	expectationSeries []*IndexSaverMockSaveIndexFromHeavyExpectation
}

type IndexSaverMockSaveIndexFromHeavyExpectation struct {
	input  *IndexSaverMockSaveIndexFromHeavyInput
	result *IndexSaverMockSaveIndexFromHeavyResult
}

type IndexSaverMockSaveIndexFromHeavyInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.Reference
	p3 *insolar.Reference
}

type IndexSaverMockSaveIndexFromHeavyResult struct {
	r  Lifeline
	r1 error
}

//Expect specifies that invocation of IndexSaver.SaveIndexFromHeavy is expected from 1 to Infinity times
func (m *mIndexSaverMockSaveIndexFromHeavy) Expect(p context.Context, p1 insolar.ID, p2 insolar.Reference, p3 *insolar.Reference) *mIndexSaverMockSaveIndexFromHeavy {
	m.mock.SaveIndexFromHeavyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexSaverMockSaveIndexFromHeavyExpectation{}
	}
	m.mainExpectation.input = &IndexSaverMockSaveIndexFromHeavyInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of IndexSaver.SaveIndexFromHeavy
func (m *mIndexSaverMockSaveIndexFromHeavy) Return(r Lifeline, r1 error) *IndexSaverMock {
	m.mock.SaveIndexFromHeavyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexSaverMockSaveIndexFromHeavyExpectation{}
	}
	m.mainExpectation.result = &IndexSaverMockSaveIndexFromHeavyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexSaver.SaveIndexFromHeavy is expected once
func (m *mIndexSaverMockSaveIndexFromHeavy) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.Reference, p3 *insolar.Reference) *IndexSaverMockSaveIndexFromHeavyExpectation {
	m.mock.SaveIndexFromHeavyFunc = nil
	m.mainExpectation = nil

	expectation := &IndexSaverMockSaveIndexFromHeavyExpectation{}
	expectation.input = &IndexSaverMockSaveIndexFromHeavyInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexSaverMockSaveIndexFromHeavyExpectation) Return(r Lifeline, r1 error) {
	e.result = &IndexSaverMockSaveIndexFromHeavyResult{r, r1}
}

//Set uses given function f as a mock of IndexSaver.SaveIndexFromHeavy method
func (m *mIndexSaverMockSaveIndexFromHeavy) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.Reference, p3 *insolar.Reference) (r Lifeline, r1 error)) *IndexSaverMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SaveIndexFromHeavyFunc = f
	return m.mock
}

//SaveIndexFromHeavy implements github.com/insolar/insolar/ledger/storage/object.IndexSaver interface
func (m *IndexSaverMock) SaveIndexFromHeavy(p context.Context, p1 insolar.ID, p2 insolar.Reference, p3 *insolar.Reference) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.SaveIndexFromHeavyPreCounter, 1)
	defer atomic.AddUint64(&m.SaveIndexFromHeavyCounter, 1)

	if len(m.SaveIndexFromHeavyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SaveIndexFromHeavyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexSaverMock.SaveIndexFromHeavy. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SaveIndexFromHeavyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexSaverMockSaveIndexFromHeavyInput{p, p1, p2, p3}, "IndexSaver.SaveIndexFromHeavy got unexpected parameters")

		result := m.SaveIndexFromHeavyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexSaverMock.SaveIndexFromHeavy")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SaveIndexFromHeavyMock.mainExpectation != nil {

		input := m.SaveIndexFromHeavyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexSaverMockSaveIndexFromHeavyInput{p, p1, p2, p3}, "IndexSaver.SaveIndexFromHeavy got unexpected parameters")
		}

		result := m.SaveIndexFromHeavyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexSaverMock.SaveIndexFromHeavy")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SaveIndexFromHeavyFunc == nil {
		m.t.Fatalf("Unexpected call to IndexSaverMock.SaveIndexFromHeavy. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SaveIndexFromHeavyFunc(p, p1, p2, p3)
}

//SaveIndexFromHeavyMinimockCounter returns a count of IndexSaverMock.SaveIndexFromHeavyFunc invocations
func (m *IndexSaverMock) SaveIndexFromHeavyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SaveIndexFromHeavyCounter)
}

//SaveIndexFromHeavyMinimockPreCounter returns the value of IndexSaverMock.SaveIndexFromHeavy invocations
func (m *IndexSaverMock) SaveIndexFromHeavyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SaveIndexFromHeavyPreCounter)
}

//SaveIndexFromHeavyFinished returns true if mock invocations count is ok
func (m *IndexSaverMock) SaveIndexFromHeavyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SaveIndexFromHeavyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SaveIndexFromHeavyCounter) == uint64(len(m.SaveIndexFromHeavyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SaveIndexFromHeavyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SaveIndexFromHeavyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SaveIndexFromHeavyFunc != nil {
		return atomic.LoadUint64(&m.SaveIndexFromHeavyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexSaverMock) ValidateCallCounters() {

	if !m.SaveIndexFromHeavyFinished() {
		m.t.Fatal("Expected call to IndexSaverMock.SaveIndexFromHeavy")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexSaverMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexSaverMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexSaverMock) MinimockFinish() {

	if !m.SaveIndexFromHeavyFinished() {
		m.t.Fatal("Expected call to IndexSaverMock.SaveIndexFromHeavy")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexSaverMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexSaverMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SaveIndexFromHeavyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SaveIndexFromHeavyFinished() {
				m.t.Error("Expected call to IndexSaverMock.SaveIndexFromHeavy")
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
func (m *IndexSaverMock) AllMocksCalled() bool {

	if !m.SaveIndexFromHeavyFinished() {
		return false
	}

	return true
}

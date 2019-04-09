package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexCollectionAccessor" can be found in github.com/insolar/insolar/ledger/storage/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexCollectionAccessorMock implements github.com/insolar/insolar/ledger/storage/object.IndexCollectionAccessor
type IndexCollectionAccessorMock struct {
	t minimock.Tester

	ForPulseAndJetFunc       func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r map[insolar.ID]Lifeline)
	ForPulseAndJetCounter    uint64
	ForPulseAndJetPreCounter uint64
	ForPulseAndJetMock       mIndexCollectionAccessorMockForPulseAndJet
}

//NewIndexCollectionAccessorMock returns a mock for github.com/insolar/insolar/ledger/storage/object.IndexCollectionAccessor
func NewIndexCollectionAccessorMock(t minimock.Tester) *IndexCollectionAccessorMock {
	m := &IndexCollectionAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseAndJetMock = mIndexCollectionAccessorMockForPulseAndJet{mock: m}

	return m
}

type mIndexCollectionAccessorMockForPulseAndJet struct {
	mock              *IndexCollectionAccessorMock
	mainExpectation   *IndexCollectionAccessorMockForPulseAndJetExpectation
	expectationSeries []*IndexCollectionAccessorMockForPulseAndJetExpectation
}

type IndexCollectionAccessorMockForPulseAndJetExpectation struct {
	input  *IndexCollectionAccessorMockForPulseAndJetInput
	result *IndexCollectionAccessorMockForPulseAndJetResult
}

type IndexCollectionAccessorMockForPulseAndJetInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 insolar.PulseNumber
}

type IndexCollectionAccessorMockForPulseAndJetResult struct {
	r map[insolar.ID]Lifeline
}

//Expect specifies that invocation of IndexCollectionAccessor.ForPulseAndJet is expected from 1 to Infinity times
func (m *mIndexCollectionAccessorMockForPulseAndJet) Expect(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *mIndexCollectionAccessorMockForPulseAndJet {
	m.mock.ForPulseAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCollectionAccessorMockForPulseAndJetExpectation{}
	}
	m.mainExpectation.input = &IndexCollectionAccessorMockForPulseAndJetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexCollectionAccessor.ForPulseAndJet
func (m *mIndexCollectionAccessorMockForPulseAndJet) Return(r map[insolar.ID]Lifeline) *IndexCollectionAccessorMock {
	m.mock.ForPulseAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCollectionAccessorMockForPulseAndJetExpectation{}
	}
	m.mainExpectation.result = &IndexCollectionAccessorMockForPulseAndJetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexCollectionAccessor.ForPulseAndJet is expected once
func (m *mIndexCollectionAccessorMockForPulseAndJet) ExpectOnce(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *IndexCollectionAccessorMockForPulseAndJetExpectation {
	m.mock.ForPulseAndJetFunc = nil
	m.mainExpectation = nil

	expectation := &IndexCollectionAccessorMockForPulseAndJetExpectation{}
	expectation.input = &IndexCollectionAccessorMockForPulseAndJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexCollectionAccessorMockForPulseAndJetExpectation) Return(r map[insolar.ID]Lifeline) {
	e.result = &IndexCollectionAccessorMockForPulseAndJetResult{r}
}

//Set uses given function f as a mock of IndexCollectionAccessor.ForPulseAndJet method
func (m *mIndexCollectionAccessorMockForPulseAndJet) Set(f func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r map[insolar.ID]Lifeline)) *IndexCollectionAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseAndJetFunc = f
	return m.mock
}

//ForPulseAndJet implements github.com/insolar/insolar/ledger/storage/object.IndexCollectionAccessor interface
func (m *IndexCollectionAccessorMock) ForPulseAndJet(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r map[insolar.ID]Lifeline) {
	counter := atomic.AddUint64(&m.ForPulseAndJetPreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseAndJetCounter, 1)

	if len(m.ForPulseAndJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseAndJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexCollectionAccessorMock.ForPulseAndJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPulseAndJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexCollectionAccessorMockForPulseAndJetInput{p, p1, p2}, "IndexCollectionAccessor.ForPulseAndJet got unexpected parameters")

		result := m.ForPulseAndJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexCollectionAccessorMock.ForPulseAndJet")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseAndJetMock.mainExpectation != nil {

		input := m.ForPulseAndJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexCollectionAccessorMockForPulseAndJetInput{p, p1, p2}, "IndexCollectionAccessor.ForPulseAndJet got unexpected parameters")
		}

		result := m.ForPulseAndJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexCollectionAccessorMock.ForPulseAndJet")
		}

		r = result.r

		return
	}

	if m.ForPulseAndJetFunc == nil {
		m.t.Fatalf("Unexpected call to IndexCollectionAccessorMock.ForPulseAndJet. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPulseAndJetFunc(p, p1, p2)
}

//ForPulseAndJetMinimockCounter returns a count of IndexCollectionAccessorMock.ForPulseAndJetFunc invocations
func (m *IndexCollectionAccessorMock) ForPulseAndJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseAndJetCounter)
}

//ForPulseAndJetMinimockPreCounter returns the value of IndexCollectionAccessorMock.ForPulseAndJet invocations
func (m *IndexCollectionAccessorMock) ForPulseAndJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseAndJetPreCounter)
}

//ForPulseAndJetFinished returns true if mock invocations count is ok
func (m *IndexCollectionAccessorMock) ForPulseAndJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseAndJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) == uint64(len(m.ForPulseAndJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseAndJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseAndJetFunc != nil {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCollectionAccessorMock) ValidateCallCounters() {

	if !m.ForPulseAndJetFinished() {
		m.t.Fatal("Expected call to IndexCollectionAccessorMock.ForPulseAndJet")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCollectionAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexCollectionAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexCollectionAccessorMock) MinimockFinish() {

	if !m.ForPulseAndJetFinished() {
		m.t.Fatal("Expected call to IndexCollectionAccessorMock.ForPulseAndJet")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexCollectionAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexCollectionAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseAndJetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseAndJetFinished() {
				m.t.Error("Expected call to IndexCollectionAccessorMock.ForPulseAndJet")
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
func (m *IndexCollectionAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseAndJetFinished() {
		return false
	}

	return true
}

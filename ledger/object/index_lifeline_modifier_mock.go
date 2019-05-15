package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexLifelineModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexLifelineModifierMock implements github.com/insolar/insolar/ledger/object.IndexLifelineModifier
type IndexLifelineModifierMock struct {
	t minimock.Tester

	SetLifelineFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)
	SetLifelineCounter    uint64
	SetLifelinePreCounter uint64
	SetLifelineMock       mIndexLifelineModifierMockSetLifeline
}

//NewIndexLifelineModifierMock returns a mock for github.com/insolar/insolar/ledger/object.IndexLifelineModifier
func NewIndexLifelineModifierMock(t minimock.Tester) *IndexLifelineModifierMock {
	m := &IndexLifelineModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetLifelineMock = mIndexLifelineModifierMockSetLifeline{mock: m}

	return m
}

type mIndexLifelineModifierMockSetLifeline struct {
	mock              *IndexLifelineModifierMock
	mainExpectation   *IndexLifelineModifierMockSetLifelineExpectation
	expectationSeries []*IndexLifelineModifierMockSetLifelineExpectation
}

type IndexLifelineModifierMockSetLifelineExpectation struct {
	input  *IndexLifelineModifierMockSetLifelineInput
	result *IndexLifelineModifierMockSetLifelineResult
}

type IndexLifelineModifierMockSetLifelineInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 Lifeline
}

type IndexLifelineModifierMockSetLifelineResult struct {
	r error
}

//Expect specifies that invocation of IndexLifelineModifier.SetLifeline is expected from 1 to Infinity times
func (m *mIndexLifelineModifierMockSetLifeline) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *mIndexLifelineModifierMockSetLifeline {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLifelineModifierMockSetLifelineExpectation{}
	}
	m.mainExpectation.input = &IndexLifelineModifierMockSetLifelineInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of IndexLifelineModifier.SetLifeline
func (m *mIndexLifelineModifierMockSetLifeline) Return(r error) *IndexLifelineModifierMock {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexLifelineModifierMockSetLifelineExpectation{}
	}
	m.mainExpectation.result = &IndexLifelineModifierMockSetLifelineResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexLifelineModifier.SetLifeline is expected once
func (m *mIndexLifelineModifierMockSetLifeline) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *IndexLifelineModifierMockSetLifelineExpectation {
	m.mock.SetLifelineFunc = nil
	m.mainExpectation = nil

	expectation := &IndexLifelineModifierMockSetLifelineExpectation{}
	expectation.input = &IndexLifelineModifierMockSetLifelineInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexLifelineModifierMockSetLifelineExpectation) Return(r error) {
	e.result = &IndexLifelineModifierMockSetLifelineResult{r}
}

//Set uses given function f as a mock of IndexLifelineModifier.SetLifeline method
func (m *mIndexLifelineModifierMockSetLifeline) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)) *IndexLifelineModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLifelineFunc = f
	return m.mock
}

//SetLifeline implements github.com/insolar/insolar/ledger/object.IndexLifelineModifier interface
func (m *IndexLifelineModifierMock) SetLifeline(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetLifelinePreCounter, 1)
	defer atomic.AddUint64(&m.SetLifelineCounter, 1)

	if len(m.SetLifelineMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLifelineMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexLifelineModifierMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetLifelineMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexLifelineModifierMockSetLifelineInput{p, p1, p2, p3}, "IndexLifelineModifier.SetLifeline got unexpected parameters")

		result := m.SetLifelineMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexLifelineModifierMock.SetLifeline")
			return
		}

		r = result.r

		return
	}

	if m.SetLifelineMock.mainExpectation != nil {

		input := m.SetLifelineMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexLifelineModifierMockSetLifelineInput{p, p1, p2, p3}, "IndexLifelineModifier.SetLifeline got unexpected parameters")
		}

		result := m.SetLifelineMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexLifelineModifierMock.SetLifeline")
		}

		r = result.r

		return
	}

	if m.SetLifelineFunc == nil {
		m.t.Fatalf("Unexpected call to IndexLifelineModifierMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetLifelineFunc(p, p1, p2, p3)
}

//SetLifelineMinimockCounter returns a count of IndexLifelineModifierMock.SetLifelineFunc invocations
func (m *IndexLifelineModifierMock) SetLifelineMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineCounter)
}

//SetLifelineMinimockPreCounter returns the value of IndexLifelineModifierMock.SetLifeline invocations
func (m *IndexLifelineModifierMock) SetLifelineMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelinePreCounter)
}

//SetLifelineFinished returns true if mock invocations count is ok
func (m *IndexLifelineModifierMock) SetLifelineFinished() bool {
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
func (m *IndexLifelineModifierMock) ValidateCallCounters() {

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to IndexLifelineModifierMock.SetLifeline")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexLifelineModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexLifelineModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexLifelineModifierMock) MinimockFinish() {

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to IndexLifelineModifierMock.SetLifeline")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexLifelineModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexLifelineModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetLifelineFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetLifelineFinished() {
				m.t.Error("Expected call to IndexLifelineModifierMock.SetLifeline")
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
func (m *IndexLifelineModifierMock) AllMocksCalled() bool {

	if !m.SetLifelineFinished() {
		return false
	}

	return true
}

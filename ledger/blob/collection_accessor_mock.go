package blob

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CollectionAccessor" can be found in github.com/insolar/insolar/ledger/blob
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//CollectionAccessorMock implements github.com/insolar/insolar/ledger/blob.CollectionAccessor
type CollectionAccessorMock struct {
	t minimock.Tester

	ForPulseFunc       func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []Blob)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mCollectionAccessorMockForPulse
}

//NewCollectionAccessorMock returns a mock for github.com/insolar/insolar/ledger/blob.CollectionAccessor
func NewCollectionAccessorMock(t minimock.Tester) *CollectionAccessorMock {
	m := &CollectionAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseMock = mCollectionAccessorMockForPulse{mock: m}

	return m
}

type mCollectionAccessorMockForPulse struct {
	mock              *CollectionAccessorMock
	mainExpectation   *CollectionAccessorMockForPulseExpectation
	expectationSeries []*CollectionAccessorMockForPulseExpectation
}

type CollectionAccessorMockForPulseExpectation struct {
	input  *CollectionAccessorMockForPulseInput
	result *CollectionAccessorMockForPulseResult
}

type CollectionAccessorMockForPulseInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 insolar.PulseNumber
}

type CollectionAccessorMockForPulseResult struct {
	r []Blob
}

//Expect specifies that invocation of CollectionAccessor.ForPulse is expected from 1 to Infinity times
func (m *mCollectionAccessorMockForPulse) Expect(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *mCollectionAccessorMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CollectionAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.input = &CollectionAccessorMockForPulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of CollectionAccessor.ForPulse
func (m *mCollectionAccessorMockForPulse) Return(r []Blob) *CollectionAccessorMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CollectionAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.result = &CollectionAccessorMockForPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CollectionAccessor.ForPulse is expected once
func (m *mCollectionAccessorMockForPulse) ExpectOnce(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *CollectionAccessorMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &CollectionAccessorMockForPulseExpectation{}
	expectation.input = &CollectionAccessorMockForPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CollectionAccessorMockForPulseExpectation) Return(r []Blob) {
	e.result = &CollectionAccessorMockForPulseResult{r}
}

//Set uses given function f as a mock of CollectionAccessor.ForPulse method
func (m *mCollectionAccessorMockForPulse) Set(f func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []Blob)) *CollectionAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/blob.CollectionAccessor interface
func (m *CollectionAccessorMock) ForPulse(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []Blob) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CollectionAccessorMock.ForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CollectionAccessorMockForPulseInput{p, p1, p2}, "CollectionAccessor.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CollectionAccessorMock.ForPulse")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CollectionAccessorMockForPulseInput{p, p1, p2}, "CollectionAccessor.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CollectionAccessorMock.ForPulse")
		}

		r = result.r

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to CollectionAccessorMock.ForPulse. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPulseFunc(p, p1, p2)
}

//ForPulseMinimockCounter returns a count of CollectionAccessorMock.ForPulseFunc invocations
func (m *CollectionAccessorMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

//ForPulseMinimockPreCounter returns the value of CollectionAccessorMock.ForPulse invocations
func (m *CollectionAccessorMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

//ForPulseFinished returns true if mock invocations count is ok
func (m *CollectionAccessorMock) ForPulseFinished() bool {
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
func (m *CollectionAccessorMock) ValidateCallCounters() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to CollectionAccessorMock.ForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CollectionAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CollectionAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CollectionAccessorMock) MinimockFinish() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to CollectionAccessorMock.ForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CollectionAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CollectionAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseFinished() {
				m.t.Error("Expected call to CollectionAccessorMock.ForPulse")
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
func (m *CollectionAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseFinished() {
		return false
	}

	return true
}

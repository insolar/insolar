package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexBucketAccessor" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexBucketAccessorMock implements github.com/insolar/insolar/ledger/object.IndexBucketAccessor
type IndexBucketAccessorMock struct {
	t minimock.Tester

	ForPulseFunc       func(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mIndexBucketAccessorMockForPulse
}

//NewIndexBucketAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.IndexBucketAccessor
func NewIndexBucketAccessorMock(t minimock.Tester) *IndexBucketAccessorMock {
	m := &IndexBucketAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseMock = mIndexBucketAccessorMockForPulse{mock: m}

	return m
}

type mIndexBucketAccessorMockForPulse struct {
	mock              *IndexBucketAccessorMock
	mainExpectation   *IndexBucketAccessorMockForPulseExpectation
	expectationSeries []*IndexBucketAccessorMockForPulseExpectation
}

type IndexBucketAccessorMockForPulseExpectation struct {
	input  *IndexBucketAccessorMockForPulseInput
	result *IndexBucketAccessorMockForPulseResult
}

type IndexBucketAccessorMockForPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type IndexBucketAccessorMockForPulseResult struct {
	r []FilamentIndex
}

//Expect specifies that invocation of IndexBucketAccessor.ForPulse is expected from 1 to Infinity times
func (m *mIndexBucketAccessorMockForPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mIndexBucketAccessorMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexBucketAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.input = &IndexBucketAccessorMockForPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of IndexBucketAccessor.ForPulse
func (m *mIndexBucketAccessorMockForPulse) Return(r []FilamentIndex) *IndexBucketAccessorMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexBucketAccessorMockForPulseExpectation{}
	}
	m.mainExpectation.result = &IndexBucketAccessorMockForPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexBucketAccessor.ForPulse is expected once
func (m *mIndexBucketAccessorMockForPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *IndexBucketAccessorMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &IndexBucketAccessorMockForPulseExpectation{}
	expectation.input = &IndexBucketAccessorMockForPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexBucketAccessorMockForPulseExpectation) Return(r []FilamentIndex) {
	e.result = &IndexBucketAccessorMockForPulseResult{r}
}

//Set uses given function f as a mock of IndexBucketAccessor.ForPulse method
func (m *mIndexBucketAccessorMockForPulse) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex)) *IndexBucketAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/object.IndexBucketAccessor interface
func (m *IndexBucketAccessorMock) ForPulse(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexBucketAccessorMock.ForPulse. %v %v", p, p1)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexBucketAccessorMockForPulseInput{p, p1}, "IndexBucketAccessor.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexBucketAccessorMock.ForPulse")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexBucketAccessorMockForPulseInput{p, p1}, "IndexBucketAccessor.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexBucketAccessorMock.ForPulse")
		}

		r = result.r

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to IndexBucketAccessorMock.ForPulse. %v %v", p, p1)
		return
	}

	return m.ForPulseFunc(p, p1)
}

//ForPulseMinimockCounter returns a count of IndexBucketAccessorMock.ForPulseFunc invocations
func (m *IndexBucketAccessorMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

//ForPulseMinimockPreCounter returns the value of IndexBucketAccessorMock.ForPulse invocations
func (m *IndexBucketAccessorMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

//ForPulseFinished returns true if mock invocations count is ok
func (m *IndexBucketAccessorMock) ForPulseFinished() bool {
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
func (m *IndexBucketAccessorMock) ValidateCallCounters() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to IndexBucketAccessorMock.ForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexBucketAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexBucketAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexBucketAccessorMock) MinimockFinish() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to IndexBucketAccessorMock.ForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexBucketAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexBucketAccessorMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IndexBucketAccessorMock.ForPulse")
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
func (m *IndexBucketAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseFinished() {
		return false
	}

	return true
}

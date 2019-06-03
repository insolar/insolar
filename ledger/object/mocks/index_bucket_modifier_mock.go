package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexBucketModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/object"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexBucketModifierMock implements github.com/insolar/insolar/ledger/object.IndexBucketModifier
type IndexBucketModifierMock struct {
	t minimock.Tester

	SetBucketFunc       func(p context.Context, p1 insolar.PulseNumber, p2 object.FilamentIndex) (r error)
	SetBucketCounter    uint64
	SetBucketPreCounter uint64
	SetBucketMock       mIndexBucketModifierMockSetBucket
}

//NewIndexBucketModifierMock returns a mock for github.com/insolar/insolar/ledger/object.IndexBucketModifier
func NewIndexBucketModifierMock(t minimock.Tester) *IndexBucketModifierMock {
	m := &IndexBucketModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetBucketMock = mIndexBucketModifierMockSetBucket{mock: m}

	return m
}

type mIndexBucketModifierMockSetBucket struct {
	mock              *IndexBucketModifierMock
	mainExpectation   *IndexBucketModifierMockSetBucketExpectation
	expectationSeries []*IndexBucketModifierMockSetBucketExpectation
}

type IndexBucketModifierMockSetBucketExpectation struct {
	input  *IndexBucketModifierMockSetBucketInput
	result *IndexBucketModifierMockSetBucketResult
}

type IndexBucketModifierMockSetBucketInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 object.FilamentIndex
}

type IndexBucketModifierMockSetBucketResult struct {
	r error
}

//Expect specifies that invocation of IndexBucketModifier.SetBucket is expected from 1 to Infinity times
func (m *mIndexBucketModifierMockSetBucket) Expect(p context.Context, p1 insolar.PulseNumber, p2 object.FilamentIndex) *mIndexBucketModifierMockSetBucket {
	m.mock.SetBucketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexBucketModifierMockSetBucketExpectation{}
	}
	m.mainExpectation.input = &IndexBucketModifierMockSetBucketInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexBucketModifier.SetBucket
func (m *mIndexBucketModifierMockSetBucket) Return(r error) *IndexBucketModifierMock {
	m.mock.SetBucketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexBucketModifierMockSetBucketExpectation{}
	}
	m.mainExpectation.result = &IndexBucketModifierMockSetBucketResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexBucketModifier.SetBucket is expected once
func (m *mIndexBucketModifierMockSetBucket) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 object.FilamentIndex) *IndexBucketModifierMockSetBucketExpectation {
	m.mock.SetBucketFunc = nil
	m.mainExpectation = nil

	expectation := &IndexBucketModifierMockSetBucketExpectation{}
	expectation.input = &IndexBucketModifierMockSetBucketInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexBucketModifierMockSetBucketExpectation) Return(r error) {
	e.result = &IndexBucketModifierMockSetBucketResult{r}
}

//Set uses given function f as a mock of IndexBucketModifier.SetBucket method
func (m *mIndexBucketModifierMockSetBucket) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 object.FilamentIndex) (r error)) *IndexBucketModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetBucketFunc = f
	return m.mock
}

//SetBucket implements github.com/insolar/insolar/ledger/object.IndexBucketModifier interface
func (m *IndexBucketModifierMock) SetBucket(p context.Context, p1 insolar.PulseNumber, p2 object.FilamentIndex) (r error) {
	counter := atomic.AddUint64(&m.SetBucketPreCounter, 1)
	defer atomic.AddUint64(&m.SetBucketCounter, 1)

	if len(m.SetBucketMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetBucketMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexBucketModifierMock.SetBucket. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetBucketMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexBucketModifierMockSetBucketInput{p, p1, p2}, "IndexBucketModifier.SetBucket got unexpected parameters")

		result := m.SetBucketMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexBucketModifierMock.SetBucket")
			return
		}

		r = result.r

		return
	}

	if m.SetBucketMock.mainExpectation != nil {

		input := m.SetBucketMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexBucketModifierMockSetBucketInput{p, p1, p2}, "IndexBucketModifier.SetBucket got unexpected parameters")
		}

		result := m.SetBucketMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexBucketModifierMock.SetBucket")
		}

		r = result.r

		return
	}

	if m.SetBucketFunc == nil {
		m.t.Fatalf("Unexpected call to IndexBucketModifierMock.SetBucket. %v %v %v", p, p1, p2)
		return
	}

	return m.SetBucketFunc(p, p1, p2)
}

//SetBucketMinimockCounter returns a count of IndexBucketModifierMock.SetBucketFunc invocations
func (m *IndexBucketModifierMock) SetBucketMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetBucketCounter)
}

//SetBucketMinimockPreCounter returns the value of IndexBucketModifierMock.SetBucket invocations
func (m *IndexBucketModifierMock) SetBucketMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetBucketPreCounter)
}

//SetBucketFinished returns true if mock invocations count is ok
func (m *IndexBucketModifierMock) SetBucketFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetBucketMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetBucketCounter) == uint64(len(m.SetBucketMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetBucketMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetBucketCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetBucketFunc != nil {
		return atomic.LoadUint64(&m.SetBucketCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexBucketModifierMock) ValidateCallCounters() {

	if !m.SetBucketFinished() {
		m.t.Fatal("Expected call to IndexBucketModifierMock.SetBucket")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexBucketModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexBucketModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexBucketModifierMock) MinimockFinish() {

	if !m.SetBucketFinished() {
		m.t.Fatal("Expected call to IndexBucketModifierMock.SetBucket")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexBucketModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexBucketModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetBucketFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetBucketFinished() {
				m.t.Error("Expected call to IndexBucketModifierMock.SetBucket")
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
func (m *IndexBucketModifierMock) AllMocksCalled() bool {

	if !m.SetBucketFinished() {
		return false
	}

	return true
}

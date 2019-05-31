package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ObjectIndexModifier" can be found in github.com/insolar/insolar/ledger/object
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

// ObjectIndexModifierMock implements github.com/insolar/insolar/ledger/object.ObjectIndexModifier
type ObjectIndexModifierMock struct {
	t minimock.Tester

	SetObjectIndexFunc       func(p context.Context, p1 insolar.PulseNumber, p2 object.ObjectIndex) (r error)
	SetObjectIndexCounter    uint64
	SetObjectIndexPreCounter uint64
	SetObjectIndexMock       mObjectIndexModifierMockSetObjectIndex
}

// NewObjectIndexModifierMock returns a mock for github.com/insolar/insolar/ledger/object.ObjectIndexModifier
func NewObjectIndexModifierMock(t minimock.Tester) *ObjectIndexModifierMock {
	m := &ObjectIndexModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetObjectIndexMock = mObjectIndexModifierMockSetObjectIndex{mock: m}

	return m
}

type mObjectIndexModifierMockSetObjectIndex struct {
	mock              *ObjectIndexModifierMock
	mainExpectation   *ObjectIndexModifierMockSetObjectIndexExpectation
	expectationSeries []*ObjectIndexModifierMockSetObjectIndexExpectation
}

type ObjectIndexModifierMockSetObjectIndexExpectation struct {
	input  *ObjectIndexModifierMockSetObjectIndexInput
	result *ObjectIndexModifierMockSetObjectIndexResult
}

type ObjectIndexModifierMockSetObjectIndexInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 object.ObjectIndex
}

type ObjectIndexModifierMockSetObjectIndexResult struct {
	r error
}

// Expect specifies that invocation of ObjectIndexModifier.SetObjectIndex is expected from 1 to Infinity times
func (m *mObjectIndexModifierMockSetObjectIndex) Expect(p context.Context, p1 insolar.PulseNumber, p2 object.ObjectIndex) *mObjectIndexModifierMockSetObjectIndex {
	m.mock.SetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectIndexModifierMockSetObjectIndexExpectation{}
	}
	m.mainExpectation.input = &ObjectIndexModifierMockSetObjectIndexInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of ObjectIndexModifier.SetObjectIndex
func (m *mObjectIndexModifierMockSetObjectIndex) Return(r error) *ObjectIndexModifierMock {
	m.mock.SetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectIndexModifierMockSetObjectIndexExpectation{}
	}
	m.mainExpectation.result = &ObjectIndexModifierMockSetObjectIndexResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of ObjectIndexModifier.SetObjectIndex is expected once
func (m *mObjectIndexModifierMockSetObjectIndex) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 object.ObjectIndex) *ObjectIndexModifierMockSetObjectIndexExpectation {
	m.mock.SetObjectIndexFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectIndexModifierMockSetObjectIndexExpectation{}
	expectation.input = &ObjectIndexModifierMockSetObjectIndexInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectIndexModifierMockSetObjectIndexExpectation) Return(r error) {
	e.result = &ObjectIndexModifierMockSetObjectIndexResult{r}
}

// Set uses given function f as a mock of ObjectIndexModifier.SetObjectIndex method
func (m *mObjectIndexModifierMockSetObjectIndex) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 object.ObjectIndex) (r error)) *ObjectIndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetObjectIndexFunc = f
	return m.mock
}

// SetObjectIndex implements github.com/insolar/insolar/ledger/object.ObjectIndexModifier interface
func (m *ObjectIndexModifierMock) SetObjectIndex(p context.Context, p1 insolar.PulseNumber, p2 object.ObjectIndex) (r error) {
	counter := atomic.AddUint64(&m.SetObjectIndexPreCounter, 1)
	defer atomic.AddUint64(&m.SetObjectIndexCounter, 1)

	if len(m.SetObjectIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetObjectIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectIndexModifierMock.SetObjectIndex. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetObjectIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectIndexModifierMockSetObjectIndexInput{p, p1, p2}, "ObjectIndexModifier.SetObjectIndex got unexpected parameters")

		result := m.SetObjectIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectIndexModifierMock.SetObjectIndex")
			return
		}

		r = result.r

		return
	}

	if m.SetObjectIndexMock.mainExpectation != nil {

		input := m.SetObjectIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectIndexModifierMockSetObjectIndexInput{p, p1, p2}, "ObjectIndexModifier.SetObjectIndex got unexpected parameters")
		}

		result := m.SetObjectIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectIndexModifierMock.SetObjectIndex")
		}

		r = result.r

		return
	}

	if m.SetObjectIndexFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectIndexModifierMock.SetObjectIndex. %v %v %v", p, p1, p2)
		return
	}

	return m.SetObjectIndexFunc(p, p1, p2)
}

// SetObjectIndexMinimockCounter returns a count of ObjectIndexModifierMock.SetObjectIndexFunc invocations
func (m *ObjectIndexModifierMock) SetObjectIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetObjectIndexCounter)
}

// SetObjectIndexMinimockPreCounter returns the value of ObjectIndexModifierMock.SetObjectIndex invocations
func (m *ObjectIndexModifierMock) SetObjectIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetObjectIndexPreCounter)
}

// SetObjectIndexFinished returns true if mock invocations count is ok
func (m *ObjectIndexModifierMock) SetObjectIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetObjectIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetObjectIndexCounter) == uint64(len(m.SetObjectIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetObjectIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetObjectIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetObjectIndexFunc != nil {
		return atomic.LoadUint64(&m.SetObjectIndexCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectIndexModifierMock) ValidateCallCounters() {

	if !m.SetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectIndexModifierMock.SetObjectIndex")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectIndexModifierMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ObjectIndexModifierMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ObjectIndexModifierMock) MinimockFinish() {

	if !m.SetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectIndexModifierMock.SetObjectIndex")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ObjectIndexModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *ObjectIndexModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetObjectIndexFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetObjectIndexFinished() {
				m.t.Error("Expected call to ObjectIndexModifierMock.SetObjectIndex")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *ObjectIndexModifierMock) AllMocksCalled() bool {

	if !m.SetObjectIndexFinished() {
		return false
	}

	return true
}

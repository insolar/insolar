package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ObjectIndexAccessor" can be found in github.com/insolar/insolar/ledger/object
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

// ObjectIndexAccessorMock implements github.com/insolar/insolar/ledger/object.ObjectIndexAccessor
type ObjectIndexAccessorMock struct {
	t minimock.Tester

	ForPNAndJetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []object.ObjectIndex)
	ForPNAndJetCounter    uint64
	ForPNAndJetPreCounter uint64
	ForPNAndJetMock       mObjectIndexAccessorMockForPNAndJet
}

// NewObjectIndexAccessorMock returns a mock for github.com/insolar/insolar/ledger/object.ObjectIndexAccessor
func NewObjectIndexAccessorMock(t minimock.Tester) *ObjectIndexAccessorMock {
	m := &ObjectIndexAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPNAndJetMock = mObjectIndexAccessorMockForPNAndJet{mock: m}

	return m
}

type mObjectIndexAccessorMockForPNAndJet struct {
	mock              *ObjectIndexAccessorMock
	mainExpectation   *ObjectIndexAccessorMockForPNAndJetExpectation
	expectationSeries []*ObjectIndexAccessorMockForPNAndJetExpectation
}

type ObjectIndexAccessorMockForPNAndJetExpectation struct {
	input  *ObjectIndexAccessorMockForPNAndJetInput
	result *ObjectIndexAccessorMockForPNAndJetResult
}

type ObjectIndexAccessorMockForPNAndJetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type ObjectIndexAccessorMockForPNAndJetResult struct {
	r []object.ObjectIndex
}

// Expect specifies that invocation of ObjectIndexAccessor.ForPNAndJet is expected from 1 to Infinity times
func (m *mObjectIndexAccessorMockForPNAndJet) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mObjectIndexAccessorMockForPNAndJet {
	m.mock.ForPNAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectIndexAccessorMockForPNAndJetExpectation{}
	}
	m.mainExpectation.input = &ObjectIndexAccessorMockForPNAndJetInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of ObjectIndexAccessor.ForPNAndJet
func (m *mObjectIndexAccessorMockForPNAndJet) Return(r []object.ObjectIndex) *ObjectIndexAccessorMock {
	m.mock.ForPNAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectIndexAccessorMockForPNAndJetExpectation{}
	}
	m.mainExpectation.result = &ObjectIndexAccessorMockForPNAndJetResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of ObjectIndexAccessor.ForPNAndJet is expected once
func (m *mObjectIndexAccessorMockForPNAndJet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *ObjectIndexAccessorMockForPNAndJetExpectation {
	m.mock.ForPNAndJetFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectIndexAccessorMockForPNAndJetExpectation{}
	expectation.input = &ObjectIndexAccessorMockForPNAndJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectIndexAccessorMockForPNAndJetExpectation) Return(r []object.ObjectIndex) {
	e.result = &ObjectIndexAccessorMockForPNAndJetResult{r}
}

// Set uses given function f as a mock of ObjectIndexAccessor.ForPNAndJet method
func (m *mObjectIndexAccessorMockForPNAndJet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []object.ObjectIndex)) *ObjectIndexAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPNAndJetFunc = f
	return m.mock
}

// ForPNAndJet implements github.com/insolar/insolar/ledger/object.ObjectIndexAccessor interface
func (m *ObjectIndexAccessorMock) ForPNAndJet(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r []object.ObjectIndex) {
	counter := atomic.AddUint64(&m.ForPNAndJetPreCounter, 1)
	defer atomic.AddUint64(&m.ForPNAndJetCounter, 1)

	if len(m.ForPNAndJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPNAndJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectIndexAccessorMock.ForPNAndJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPNAndJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectIndexAccessorMockForPNAndJetInput{p, p1, p2}, "ObjectIndexAccessor.ForPNAndJet got unexpected parameters")

		result := m.ForPNAndJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectIndexAccessorMock.ForPNAndJet")
			return
		}

		r = result.r

		return
	}

	if m.ForPNAndJetMock.mainExpectation != nil {

		input := m.ForPNAndJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectIndexAccessorMockForPNAndJetInput{p, p1, p2}, "ObjectIndexAccessor.ForPNAndJet got unexpected parameters")
		}

		result := m.ForPNAndJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectIndexAccessorMock.ForPNAndJet")
		}

		r = result.r

		return
	}

	if m.ForPNAndJetFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectIndexAccessorMock.ForPNAndJet. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPNAndJetFunc(p, p1, p2)
}

// ForPNAndJetMinimockCounter returns a count of ObjectIndexAccessorMock.ForPNAndJetFunc invocations
func (m *ObjectIndexAccessorMock) ForPNAndJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNAndJetCounter)
}

// ForPNAndJetMinimockPreCounter returns the value of ObjectIndexAccessorMock.ForPNAndJet invocations
func (m *ObjectIndexAccessorMock) ForPNAndJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNAndJetPreCounter)
}

//ForPNAndJetFinished returns true if mock invocations count is ok
func (m *ObjectIndexAccessorMock) ForPNAndJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPNAndJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPNAndJetCounter) == uint64(len(m.ForPNAndJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPNAndJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPNAndJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPNAndJetFunc != nil {
		return atomic.LoadUint64(&m.ForPNAndJetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectIndexAccessorMock) ValidateCallCounters() {

	if !m.ForPNAndJetFinished() {
		m.t.Fatal("Expected call to ObjectIndexAccessorMock.ForPNAndJet")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectIndexAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ObjectIndexAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ObjectIndexAccessorMock) MinimockFinish() {

	if !m.ForPNAndJetFinished() {
		m.t.Fatal("Expected call to ObjectIndexAccessorMock.ForPNAndJet")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ObjectIndexAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ObjectIndexAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPNAndJetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPNAndJetFinished() {
				m.t.Error("Expected call to ObjectIndexAccessorMock.ForPNAndJet")
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
func (m *ObjectIndexAccessorMock) AllMocksCalled() bool {

	if !m.ForPNAndJetFinished() {
		return false
	}

	return true
}

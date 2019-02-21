package nodes

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Setter" can be found in github.com/insolar/insolar/ledger/storage/nodes
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//SetterMock implements github.com/insolar/insolar/ledger/storage/nodes.Setter
type SetterMock struct {
	t minimock.Tester

	RemoveActiveNodesUntilFunc       func(p core.PulseNumber)
	RemoveActiveNodesUntilCounter    uint64
	RemoveActiveNodesUntilPreCounter uint64
	RemoveActiveNodesUntilMock       mSetterMockRemoveActiveNodesUntil

	SetFunc       func(p core.PulseNumber, p1 []core.Node) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mSetterMockSet
}

//NewSetterMock returns a mock for github.com/insolar/insolar/ledger/storage/nodes.Setter
func NewSetterMock(t minimock.Tester) *SetterMock {
	m := &SetterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveActiveNodesUntilMock = mSetterMockRemoveActiveNodesUntil{mock: m}
	m.SetMock = mSetterMockSet{mock: m}

	return m
}

type mSetterMockRemoveActiveNodesUntil struct {
	mock              *SetterMock
	mainExpectation   *SetterMockRemoveActiveNodesUntilExpectation
	expectationSeries []*SetterMockRemoveActiveNodesUntilExpectation
}

type SetterMockRemoveActiveNodesUntilExpectation struct {
	input *SetterMockRemoveActiveNodesUntilInput
}

type SetterMockRemoveActiveNodesUntilInput struct {
	p core.PulseNumber
}

//Expect specifies that invocation of Setter.RemoveActiveNodesUntil is expected from 1 to Infinity times
func (m *mSetterMockRemoveActiveNodesUntil) Expect(p core.PulseNumber) *mSetterMockRemoveActiveNodesUntil {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SetterMockRemoveActiveNodesUntilExpectation{}
	}
	m.mainExpectation.input = &SetterMockRemoveActiveNodesUntilInput{p}
	return m
}

//Return specifies results of invocation of Setter.RemoveActiveNodesUntil
func (m *mSetterMockRemoveActiveNodesUntil) Return() *SetterMock {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SetterMockRemoveActiveNodesUntilExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Setter.RemoveActiveNodesUntil is expected once
func (m *mSetterMockRemoveActiveNodesUntil) ExpectOnce(p core.PulseNumber) *SetterMockRemoveActiveNodesUntilExpectation {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.mainExpectation = nil

	expectation := &SetterMockRemoveActiveNodesUntilExpectation{}
	expectation.input = &SetterMockRemoveActiveNodesUntilInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Setter.RemoveActiveNodesUntil method
func (m *mSetterMockRemoveActiveNodesUntil) Set(f func(p core.PulseNumber)) *SetterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveActiveNodesUntilFunc = f
	return m.mock
}

//RemoveActiveNodesUntil implements github.com/insolar/insolar/ledger/storage/nodes.Setter interface
func (m *SetterMock) RemoveActiveNodesUntil(p core.PulseNumber) {
	counter := atomic.AddUint64(&m.RemoveActiveNodesUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveActiveNodesUntilCounter, 1)

	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SetterMock.RemoveActiveNodesUntil. %v", p)
			return
		}

		input := m.RemoveActiveNodesUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SetterMockRemoveActiveNodesUntilInput{p}, "Setter.RemoveActiveNodesUntil got unexpected parameters")

		return
	}

	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {

		input := m.RemoveActiveNodesUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SetterMockRemoveActiveNodesUntilInput{p}, "Setter.RemoveActiveNodesUntil got unexpected parameters")
		}

		return
	}

	if m.RemoveActiveNodesUntilFunc == nil {
		m.t.Fatalf("Unexpected call to SetterMock.RemoveActiveNodesUntil. %v", p)
		return
	}

	m.RemoveActiveNodesUntilFunc(p)
}

//RemoveActiveNodesUntilMinimockCounter returns a count of SetterMock.RemoveActiveNodesUntilFunc invocations
func (m *SetterMock) RemoveActiveNodesUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter)
}

//RemoveActiveNodesUntilMinimockPreCounter returns the value of SetterMock.RemoveActiveNodesUntil invocations
func (m *SetterMock) RemoveActiveNodesUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilPreCounter)
}

//RemoveActiveNodesUntilFinished returns true if mock invocations count is ok
func (m *SetterMock) RemoveActiveNodesUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) == uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	return true
}

type mSetterMockSet struct {
	mock              *SetterMock
	mainExpectation   *SetterMockSetExpectation
	expectationSeries []*SetterMockSetExpectation
}

type SetterMockSetExpectation struct {
	input  *SetterMockSetInput
	result *SetterMockSetResult
}

type SetterMockSetInput struct {
	p  core.PulseNumber
	p1 []core.Node
}

type SetterMockSetResult struct {
	r error
}

//Expect specifies that invocation of Setter.Set is expected from 1 to Infinity times
func (m *mSetterMockSet) Expect(p core.PulseNumber, p1 []core.Node) *mSetterMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SetterMockSetExpectation{}
	}
	m.mainExpectation.input = &SetterMockSetInput{p, p1}
	return m
}

//Return specifies results of invocation of Setter.Set
func (m *mSetterMockSet) Return(r error) *SetterMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SetterMockSetExpectation{}
	}
	m.mainExpectation.result = &SetterMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Setter.Set is expected once
func (m *mSetterMockSet) ExpectOnce(p core.PulseNumber, p1 []core.Node) *SetterMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &SetterMockSetExpectation{}
	expectation.input = &SetterMockSetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SetterMockSetExpectation) Return(r error) {
	e.result = &SetterMockSetResult{r}
}

//Set uses given function f as a mock of Setter.Set method
func (m *mSetterMockSet) Set(f func(p core.PulseNumber, p1 []core.Node) (r error)) *SetterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/storage/nodes.Setter interface
func (m *SetterMock) Set(p core.PulseNumber, p1 []core.Node) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SetterMock.Set. %v %v", p, p1)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SetterMockSetInput{p, p1}, "Setter.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SetterMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SetterMockSetInput{p, p1}, "Setter.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SetterMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to SetterMock.Set. %v %v", p, p1)
		return
	}

	return m.SetFunc(p, p1)
}

//SetMinimockCounter returns a count of SetterMock.SetFunc invocations
func (m *SetterMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of SetterMock.Set invocations
func (m *SetterMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *SetterMock) SetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCounter) == uint64(len(m.SetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFunc != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SetterMock) ValidateCallCounters() {

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to SetterMock.RemoveActiveNodesUntil")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to SetterMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SetterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SetterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SetterMock) MinimockFinish() {

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to SetterMock.RemoveActiveNodesUntil")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to SetterMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SetterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SetterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RemoveActiveNodesUntilFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveActiveNodesUntilFinished() {
				m.t.Error("Expected call to SetterMock.RemoveActiveNodesUntil")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to SetterMock.Set")
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
func (m *SetterMock) AllMocksCalled() bool {

	if !m.RemoveActiveNodesUntilFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}

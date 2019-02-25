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

	DeleteFunc       func(p core.PulseNumber)
	DeleteCounter    uint64
	DeletePreCounter uint64
	DeleteMock       mSetterMockDelete

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

	m.DeleteMock = mSetterMockDelete{mock: m}
	m.SetMock = mSetterMockSet{mock: m}

	return m
}

type mSetterMockDelete struct {
	mock              *SetterMock
	mainExpectation   *SetterMockDeleteExpectation
	expectationSeries []*SetterMockDeleteExpectation
}

type SetterMockDeleteExpectation struct {
	input *SetterMockDeleteInput
}

type SetterMockDeleteInput struct {
	p core.PulseNumber
}

//Expect specifies that invocation of Setter.Delete is expected from 1 to Infinity times
func (m *mSetterMockDelete) Expect(p core.PulseNumber) *mSetterMockDelete {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SetterMockDeleteExpectation{}
	}
	m.mainExpectation.input = &SetterMockDeleteInput{p}
	return m
}

//Return specifies results of invocation of Setter.Delete
func (m *mSetterMockDelete) Return() *SetterMock {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SetterMockDeleteExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Setter.Delete is expected once
func (m *mSetterMockDelete) ExpectOnce(p core.PulseNumber) *SetterMockDeleteExpectation {
	m.mock.DeleteFunc = nil
	m.mainExpectation = nil

	expectation := &SetterMockDeleteExpectation{}
	expectation.input = &SetterMockDeleteInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Setter.Delete method
func (m *mSetterMockDelete) Set(f func(p core.PulseNumber)) *SetterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteFunc = f
	return m.mock
}

//Delete implements github.com/insolar/insolar/ledger/storage/nodes.Setter interface
func (m *SetterMock) Delete(p core.PulseNumber) {
	counter := atomic.AddUint64(&m.DeletePreCounter, 1)
	defer atomic.AddUint64(&m.DeleteCounter, 1)

	if len(m.DeleteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SetterMock.Delete. %v", p)
			return
		}

		input := m.DeleteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SetterMockDeleteInput{p}, "Setter.Delete got unexpected parameters")

		return
	}

	if m.DeleteMock.mainExpectation != nil {

		input := m.DeleteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SetterMockDeleteInput{p}, "Setter.Delete got unexpected parameters")
		}

		return
	}

	if m.DeleteFunc == nil {
		m.t.Fatalf("Unexpected call to SetterMock.Delete. %v", p)
		return
	}

	m.DeleteFunc(p)
}

//DeleteMinimockCounter returns a count of SetterMock.DeleteFunc invocations
func (m *SetterMock) DeleteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteCounter)
}

//DeleteMinimockPreCounter returns the value of SetterMock.Delete invocations
func (m *SetterMock) DeleteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeletePreCounter)
}

//DeleteFinished returns true if mock invocations count is ok
func (m *SetterMock) DeleteFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeleteMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeleteCounter) == uint64(len(m.DeleteMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeleteMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeleteCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeleteFunc != nil {
		return atomic.LoadUint64(&m.DeleteCounter) > 0
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

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to SetterMock.Delete")
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

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to SetterMock.Delete")
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
		ok = ok && m.DeleteFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeleteFinished() {
				m.t.Error("Expected call to SetterMock.Delete")
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

	if !m.DeleteFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}

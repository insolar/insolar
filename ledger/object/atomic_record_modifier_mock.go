package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "AtomicRecordModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//AtomicRecordModifierMock implements github.com/insolar/insolar/ledger/object.AtomicRecordModifier
type AtomicRecordModifierMock struct {
	t minimock.Tester

	SetAtomicFunc       func(p context.Context, p1 ...record.Material) (r error)
	SetAtomicCounter    uint64
	SetAtomicPreCounter uint64
	SetAtomicMock       mAtomicRecordModifierMockSetAtomic
}

//NewAtomicRecordModifierMock returns a mock for github.com/insolar/insolar/ledger/object.AtomicRecordModifier
func NewAtomicRecordModifierMock(t minimock.Tester) *AtomicRecordModifierMock {
	m := &AtomicRecordModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetAtomicMock = mAtomicRecordModifierMockSetAtomic{mock: m}

	return m
}

type mAtomicRecordModifierMockSetAtomic struct {
	mock              *AtomicRecordModifierMock
	mainExpectation   *AtomicRecordModifierMockSetAtomicExpectation
	expectationSeries []*AtomicRecordModifierMockSetAtomicExpectation
}

type AtomicRecordModifierMockSetAtomicExpectation struct {
	input  *AtomicRecordModifierMockSetAtomicInput
	result *AtomicRecordModifierMockSetAtomicResult
}

type AtomicRecordModifierMockSetAtomicInput struct {
	p  context.Context
	p1 []record.Material
}

type AtomicRecordModifierMockSetAtomicResult struct {
	r error
}

//Expect specifies that invocation of AtomicRecordModifier.SetAtomic is expected from 1 to Infinity times
func (m *mAtomicRecordModifierMockSetAtomic) Expect(p context.Context, p1 ...record.Material) *mAtomicRecordModifierMockSetAtomic {
	m.mock.SetAtomicFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AtomicRecordModifierMockSetAtomicExpectation{}
	}
	m.mainExpectation.input = &AtomicRecordModifierMockSetAtomicInput{p, p1}
	return m
}

//Return specifies results of invocation of AtomicRecordModifier.SetAtomic
func (m *mAtomicRecordModifierMockSetAtomic) Return(r error) *AtomicRecordModifierMock {
	m.mock.SetAtomicFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AtomicRecordModifierMockSetAtomicExpectation{}
	}
	m.mainExpectation.result = &AtomicRecordModifierMockSetAtomicResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of AtomicRecordModifier.SetAtomic is expected once
func (m *mAtomicRecordModifierMockSetAtomic) ExpectOnce(p context.Context, p1 ...record.Material) *AtomicRecordModifierMockSetAtomicExpectation {
	m.mock.SetAtomicFunc = nil
	m.mainExpectation = nil

	expectation := &AtomicRecordModifierMockSetAtomicExpectation{}
	expectation.input = &AtomicRecordModifierMockSetAtomicInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AtomicRecordModifierMockSetAtomicExpectation) Return(r error) {
	e.result = &AtomicRecordModifierMockSetAtomicResult{r}
}

//Set uses given function f as a mock of AtomicRecordModifier.SetAtomic method
func (m *mAtomicRecordModifierMockSetAtomic) Set(f func(p context.Context, p1 ...record.Material) (r error)) *AtomicRecordModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetAtomicFunc = f
	return m.mock
}

//SetAtomic implements github.com/insolar/insolar/ledger/object.AtomicRecordModifier interface
func (m *AtomicRecordModifierMock) SetAtomic(p context.Context, p1 ...record.Material) (r error) {
	counter := atomic.AddUint64(&m.SetAtomicPreCounter, 1)
	defer atomic.AddUint64(&m.SetAtomicCounter, 1)

	if len(m.SetAtomicMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetAtomicMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AtomicRecordModifierMock.SetAtomic. %v %v", p, p1)
			return
		}

		input := m.SetAtomicMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AtomicRecordModifierMockSetAtomicInput{p, p1}, "AtomicRecordModifier.SetAtomic got unexpected parameters")

		result := m.SetAtomicMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AtomicRecordModifierMock.SetAtomic")
			return
		}

		r = result.r

		return
	}

	if m.SetAtomicMock.mainExpectation != nil {

		input := m.SetAtomicMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AtomicRecordModifierMockSetAtomicInput{p, p1}, "AtomicRecordModifier.SetAtomic got unexpected parameters")
		}

		result := m.SetAtomicMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AtomicRecordModifierMock.SetAtomic")
		}

		r = result.r

		return
	}

	if m.SetAtomicFunc == nil {
		m.t.Fatalf("Unexpected call to AtomicRecordModifierMock.SetAtomic. %v %v", p, p1)
		return
	}

	return m.SetAtomicFunc(p, p1...)
}

//SetAtomicMinimockCounter returns a count of AtomicRecordModifierMock.SetAtomicFunc invocations
func (m *AtomicRecordModifierMock) SetAtomicMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetAtomicCounter)
}

//SetAtomicMinimockPreCounter returns the value of AtomicRecordModifierMock.SetAtomic invocations
func (m *AtomicRecordModifierMock) SetAtomicMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetAtomicPreCounter)
}

//SetAtomicFinished returns true if mock invocations count is ok
func (m *AtomicRecordModifierMock) SetAtomicFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetAtomicMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetAtomicCounter) == uint64(len(m.SetAtomicMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetAtomicMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetAtomicCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetAtomicFunc != nil {
		return atomic.LoadUint64(&m.SetAtomicCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AtomicRecordModifierMock) ValidateCallCounters() {

	if !m.SetAtomicFinished() {
		m.t.Fatal("Expected call to AtomicRecordModifierMock.SetAtomic")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AtomicRecordModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AtomicRecordModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AtomicRecordModifierMock) MinimockFinish() {

	if !m.SetAtomicFinished() {
		m.t.Fatal("Expected call to AtomicRecordModifierMock.SetAtomic")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AtomicRecordModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AtomicRecordModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetAtomicFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetAtomicFinished() {
				m.t.Error("Expected call to AtomicRecordModifierMock.SetAtomic")
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
func (m *AtomicRecordModifierMock) AllMocksCalled() bool {

	if !m.SetAtomicFinished() {
		return false
	}

	return true
}

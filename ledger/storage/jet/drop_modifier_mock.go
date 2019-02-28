package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DropModifier" can be found in github.com/insolar/insolar/ledger/storage/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	storage "github.com/insolar/insolar/ledger/storage"

	testify_assert "github.com/stretchr/testify/assert"
)

//DropModifierMock implements github.com/insolar/insolar/ledger/storage/jet.DropModifier
type DropModifierMock struct {
	t minimock.Tester

	SetFunc       func(p context.Context, p1 storage.JetID, p2 JetDrop) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mDropModifierMockSet
}

//NewDropModifierMock returns a mock for github.com/insolar/insolar/ledger/storage/jet.DropModifier
func NewDropModifierMock(t minimock.Tester) *DropModifierMock {
	m := &DropModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetMock = mDropModifierMockSet{mock: m}

	return m
}

type mDropModifierMockSet struct {
	mock              *DropModifierMock
	mainExpectation   *DropModifierMockSetExpectation
	expectationSeries []*DropModifierMockSetExpectation
}

type DropModifierMockSetExpectation struct {
	input  *DropModifierMockSetInput
	result *DropModifierMockSetResult
}

type DropModifierMockSetInput struct {
	p  context.Context
	p1 storage.JetID
	p2 JetDrop
}

type DropModifierMockSetResult struct {
	r error
}

//Expect specifies that invocation of DropModifier.Set is expected from 1 to Infinity times
func (m *mDropModifierMockSet) Expect(p context.Context, p1 storage.JetID, p2 JetDrop) *mDropModifierMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropModifierMockSetExpectation{}
	}
	m.mainExpectation.input = &DropModifierMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DropModifier.Set
func (m *mDropModifierMockSet) Return(r error) *DropModifierMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropModifierMockSetExpectation{}
	}
	m.mainExpectation.result = &DropModifierMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DropModifier.Set is expected once
func (m *mDropModifierMockSet) ExpectOnce(p context.Context, p1 storage.JetID, p2 JetDrop) *DropModifierMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &DropModifierMockSetExpectation{}
	expectation.input = &DropModifierMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropModifierMockSetExpectation) Return(r error) {
	e.result = &DropModifierMockSetResult{r}
}

//Set uses given function f as a mock of DropModifier.Set method
func (m *mDropModifierMockSet) Set(f func(p context.Context, p1 storage.JetID, p2 JetDrop) (r error)) *DropModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/storage/jet.DropModifier interface
func (m *DropModifierMock) Set(p context.Context, p1 storage.JetID, p2 JetDrop) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropModifierMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropModifierMockSetInput{p, p1, p2}, "DropModifier.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropModifierMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropModifierMockSetInput{p, p1, p2}, "DropModifier.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropModifierMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to DropModifierMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of DropModifierMock.SetFunc invocations
func (m *DropModifierMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of DropModifierMock.Set invocations
func (m *DropModifierMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *DropModifierMock) SetFinished() bool {
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
func (m *DropModifierMock) ValidateCallCounters() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to DropModifierMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DropModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DropModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DropModifierMock) MinimockFinish() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to DropModifierMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DropModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DropModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetFinished() {
				m.t.Error("Expected call to DropModifierMock.Set")
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
func (m *DropModifierMock) AllMocksCalled() bool {

	if !m.SetFinished() {
		return false
	}

	return true
}

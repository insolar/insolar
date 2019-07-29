package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MemoryIndexModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//MemoryIndexModifierMock implements github.com/insolar/insolar/ledger/object.MemoryIndexModifier
type MemoryIndexModifierMock struct {
	t minimock.Tester

	SetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 record.Index)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mMemoryIndexModifierMockSet
}

//NewMemoryIndexModifierMock returns a mock for github.com/insolar/insolar/ledger/object.MemoryIndexModifier
func NewMemoryIndexModifierMock(t minimock.Tester) *MemoryIndexModifierMock {
	m := &MemoryIndexModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetMock = mMemoryIndexModifierMockSet{mock: m}

	return m
}

type mMemoryIndexModifierMockSet struct {
	mock              *MemoryIndexModifierMock
	mainExpectation   *MemoryIndexModifierMockSetExpectation
	expectationSeries []*MemoryIndexModifierMockSetExpectation
}

type MemoryIndexModifierMockSetExpectation struct {
	input *MemoryIndexModifierMockSetInput
}

type MemoryIndexModifierMockSetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 record.Index
}

//Expect specifies that invocation of MemoryIndexModifier.Set is expected from 1 to Infinity times
func (m *mMemoryIndexModifierMockSet) Expect(p context.Context, p1 insolar.PulseNumber, p2 record.Index) *mMemoryIndexModifierMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexModifierMockSetExpectation{}
	}
	m.mainExpectation.input = &MemoryIndexModifierMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of MemoryIndexModifier.Set
func (m *mMemoryIndexModifierMockSet) Return() *MemoryIndexModifierMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexModifierMockSetExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of MemoryIndexModifier.Set is expected once
func (m *mMemoryIndexModifierMockSet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 record.Index) *MemoryIndexModifierMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &MemoryIndexModifierMockSetExpectation{}
	expectation.input = &MemoryIndexModifierMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of MemoryIndexModifier.Set method
func (m *mMemoryIndexModifierMockSet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 record.Index)) *MemoryIndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/object.MemoryIndexModifier interface
func (m *MemoryIndexModifierMock) Set(p context.Context, p1 insolar.PulseNumber, p2 record.Index) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemoryIndexModifierMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MemoryIndexModifierMockSetInput{p, p1, p2}, "MemoryIndexModifier.Set got unexpected parameters")

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MemoryIndexModifierMockSetInput{p, p1, p2}, "MemoryIndexModifier.Set got unexpected parameters")
		}

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to MemoryIndexModifierMock.Set. %v %v %v", p, p1, p2)
		return
	}

	m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of MemoryIndexModifierMock.SetFunc invocations
func (m *MemoryIndexModifierMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of MemoryIndexModifierMock.Set invocations
func (m *MemoryIndexModifierMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *MemoryIndexModifierMock) SetFinished() bool {
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
func (m *MemoryIndexModifierMock) ValidateCallCounters() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to MemoryIndexModifierMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MemoryIndexModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MemoryIndexModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MemoryIndexModifierMock) MinimockFinish() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to MemoryIndexModifierMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MemoryIndexModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MemoryIndexModifierMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to MemoryIndexModifierMock.Set")
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
func (m *MemoryIndexModifierMock) AllMocksCalled() bool {

	if !m.SetFinished() {
		return false
	}

	return true
}

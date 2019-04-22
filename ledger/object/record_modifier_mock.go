package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RecordModifier" can be found in github.com/insolar/insolar/ledger/object
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

//RecordModifierMock implements github.com/insolar/insolar/ledger/object.RecordModifier
type RecordModifierMock struct {
	t minimock.Tester

	SetFunc       func(p context.Context, p1 insolar.ID, p2 record.MaterialRecord) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mRecordModifierMockSet
}

//NewRecordModifierMock returns a mock for github.com/insolar/insolar/ledger/object.RecordModifier
func NewRecordModifierMock(t minimock.Tester) *RecordModifierMock {
	m := &RecordModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetMock = mRecordModifierMockSet{mock: m}

	return m
}

type mRecordModifierMockSet struct {
	mock              *RecordModifierMock
	mainExpectation   *RecordModifierMockSetExpectation
	expectationSeries []*RecordModifierMockSetExpectation
}

type RecordModifierMockSetExpectation struct {
	input  *RecordModifierMockSetInput
	result *RecordModifierMockSetResult
}

type RecordModifierMockSetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 record.MaterialRecord
}

type RecordModifierMockSetResult struct {
	r error
}

//Expect specifies that invocation of RecordModifier.Set is expected from 1 to Infinity times
func (m *mRecordModifierMockSet) Expect(p context.Context, p1 insolar.ID, p2 record.MaterialRecord) *mRecordModifierMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordModifierMockSetExpectation{}
	}
	m.mainExpectation.input = &RecordModifierMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RecordModifier.Set
func (m *mRecordModifierMockSet) Return(r error) *RecordModifierMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RecordModifierMockSetExpectation{}
	}
	m.mainExpectation.result = &RecordModifierMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RecordModifier.Set is expected once
func (m *mRecordModifierMockSet) ExpectOnce(p context.Context, p1 insolar.ID, p2 record.MaterialRecord) *RecordModifierMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &RecordModifierMockSetExpectation{}
	expectation.input = &RecordModifierMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RecordModifierMockSetExpectation) Return(r error) {
	e.result = &RecordModifierMockSetResult{r}
}

//Set uses given function f as a mock of RecordModifier.Set method
func (m *mRecordModifierMockSet) Set(f func(p context.Context, p1 insolar.ID, p2 record.MaterialRecord) (r error)) *RecordModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/object.RecordModifier interface
func (m *RecordModifierMock) Set(p context.Context, p1 insolar.ID, p2 record.MaterialRecord) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RecordModifierMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RecordModifierMockSetInput{p, p1, p2}, "RecordModifier.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RecordModifierMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RecordModifierMockSetInput{p, p1, p2}, "RecordModifier.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RecordModifierMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to RecordModifierMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of RecordModifierMock.SetFunc invocations
func (m *RecordModifierMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of RecordModifierMock.Set invocations
func (m *RecordModifierMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *RecordModifierMock) SetFinished() bool {
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
func (m *RecordModifierMock) ValidateCallCounters() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to RecordModifierMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RecordModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RecordModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RecordModifierMock) MinimockFinish() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to RecordModifierMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RecordModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RecordModifierMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to RecordModifierMock.Set")
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
func (m *RecordModifierMock) AllMocksCalled() bool {

	if !m.SetFinished() {
		return false
	}

	return true
}

package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexModifier" can be found in github.com/insolar/insolar/ledger/storage/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexModifierMock implements github.com/insolar/insolar/ledger/storage/object.IndexModifier
type IndexModifierMock struct {
	t minimock.Tester

	SetFunc       func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mIndexModifierMockSet

	UpdateUsagePulseFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)
	UpdateUsagePulseCounter    uint64
	UpdateUsagePulsePreCounter uint64
	UpdateUsagePulseMock       mIndexModifierMockUpdateUsagePulse
}

//NewIndexModifierMock returns a mock for github.com/insolar/insolar/ledger/storage/object.IndexModifier
func NewIndexModifierMock(t minimock.Tester) *IndexModifierMock {
	m := &IndexModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetMock = mIndexModifierMockSet{mock: m}
	m.UpdateUsagePulseMock = mIndexModifierMockUpdateUsagePulse{mock: m}

	return m
}

type mIndexModifierMockSet struct {
	mock              *IndexModifierMock
	mainExpectation   *IndexModifierMockSetExpectation
	expectationSeries []*IndexModifierMockSetExpectation
}

type IndexModifierMockSetExpectation struct {
	input  *IndexModifierMockSetInput
	result *IndexModifierMockSetResult
}

type IndexModifierMockSetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 Lifeline
}

type IndexModifierMockSetResult struct {
	r error
}

//Expect specifies that invocation of IndexModifier.Set is expected from 1 to Infinity times
func (m *mIndexModifierMockSet) Expect(p context.Context, p1 insolar.ID, p2 Lifeline) *mIndexModifierMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetExpectation{}
	}
	m.mainExpectation.input = &IndexModifierMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexModifier.Set
func (m *mIndexModifierMockSet) Return(r error) *IndexModifierMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetExpectation{}
	}
	m.mainExpectation.result = &IndexModifierMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexModifier.Set is expected once
func (m *mIndexModifierMockSet) ExpectOnce(p context.Context, p1 insolar.ID, p2 Lifeline) *IndexModifierMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &IndexModifierMockSetExpectation{}
	expectation.input = &IndexModifierMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexModifierMockSetExpectation) Return(r error) {
	e.result = &IndexModifierMockSetResult{r}
}

//Set uses given function f as a mock of IndexModifier.Set method
func (m *mIndexModifierMockSet) Set(f func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)) *IndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/storage/object.IndexModifier interface
func (m *IndexModifierMock) Set(p context.Context, p1 insolar.ID, p2 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexModifierMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexModifierMockSetInput{p, p1, p2}, "IndexModifier.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexModifierMockSetInput{p, p1, p2}, "IndexModifier.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to IndexModifierMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of IndexModifierMock.SetFunc invocations
func (m *IndexModifierMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of IndexModifierMock.Set invocations
func (m *IndexModifierMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *IndexModifierMock) SetFinished() bool {
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

type mIndexModifierMockUpdateUsagePulse struct {
	mock              *IndexModifierMock
	mainExpectation   *IndexModifierMockUpdateUsagePulseExpectation
	expectationSeries []*IndexModifierMockUpdateUsagePulseExpectation
}

type IndexModifierMockUpdateUsagePulseExpectation struct {
	input  *IndexModifierMockUpdateUsagePulseInput
	result *IndexModifierMockUpdateUsagePulseResult
}

type IndexModifierMockUpdateUsagePulseInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type IndexModifierMockUpdateUsagePulseResult struct {
	r error
}

//Expect specifies that invocation of IndexModifier.UpdateUsagePulse is expected from 1 to Infinity times
func (m *mIndexModifierMockUpdateUsagePulse) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mIndexModifierMockUpdateUsagePulse {
	m.mock.UpdateUsagePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockUpdateUsagePulseExpectation{}
	}
	m.mainExpectation.input = &IndexModifierMockUpdateUsagePulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexModifier.UpdateUsagePulse
func (m *mIndexModifierMockUpdateUsagePulse) Return(r error) *IndexModifierMock {
	m.mock.UpdateUsagePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockUpdateUsagePulseExpectation{}
	}
	m.mainExpectation.result = &IndexModifierMockUpdateUsagePulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexModifier.UpdateUsagePulse is expected once
func (m *mIndexModifierMockUpdateUsagePulse) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *IndexModifierMockUpdateUsagePulseExpectation {
	m.mock.UpdateUsagePulseFunc = nil
	m.mainExpectation = nil

	expectation := &IndexModifierMockUpdateUsagePulseExpectation{}
	expectation.input = &IndexModifierMockUpdateUsagePulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexModifierMockUpdateUsagePulseExpectation) Return(r error) {
	e.result = &IndexModifierMockUpdateUsagePulseResult{r}
}

//Set uses given function f as a mock of IndexModifier.UpdateUsagePulse method
func (m *mIndexModifierMockUpdateUsagePulse) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)) *IndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateUsagePulseFunc = f
	return m.mock
}

//UpdateUsagePulse implements github.com/insolar/insolar/ledger/storage/object.IndexModifier interface
func (m *IndexModifierMock) UpdateUsagePulse(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.UpdateUsagePulsePreCounter, 1)
	defer atomic.AddUint64(&m.UpdateUsagePulseCounter, 1)

	if len(m.UpdateUsagePulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateUsagePulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexModifierMock.UpdateUsagePulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.UpdateUsagePulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexModifierMockUpdateUsagePulseInput{p, p1, p2}, "IndexModifier.UpdateUsagePulse got unexpected parameters")

		result := m.UpdateUsagePulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.UpdateUsagePulse")
			return
		}

		r = result.r

		return
	}

	if m.UpdateUsagePulseMock.mainExpectation != nil {

		input := m.UpdateUsagePulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexModifierMockUpdateUsagePulseInput{p, p1, p2}, "IndexModifier.UpdateUsagePulse got unexpected parameters")
		}

		result := m.UpdateUsagePulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.UpdateUsagePulse")
		}

		r = result.r

		return
	}

	if m.UpdateUsagePulseFunc == nil {
		m.t.Fatalf("Unexpected call to IndexModifierMock.UpdateUsagePulse. %v %v %v", p, p1, p2)
		return
	}

	return m.UpdateUsagePulseFunc(p, p1, p2)
}

//UpdateUsagePulseMinimockCounter returns a count of IndexModifierMock.UpdateUsagePulseFunc invocations
func (m *IndexModifierMock) UpdateUsagePulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateUsagePulseCounter)
}

//UpdateUsagePulseMinimockPreCounter returns the value of IndexModifierMock.UpdateUsagePulse invocations
func (m *IndexModifierMock) UpdateUsagePulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateUsagePulsePreCounter)
}

//UpdateUsagePulseFinished returns true if mock invocations count is ok
func (m *IndexModifierMock) UpdateUsagePulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpdateUsagePulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpdateUsagePulseCounter) == uint64(len(m.UpdateUsagePulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpdateUsagePulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpdateUsagePulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpdateUsagePulseFunc != nil {
		return atomic.LoadUint64(&m.UpdateUsagePulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexModifierMock) ValidateCallCounters() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.Set")
	}

	if !m.UpdateUsagePulseFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.UpdateUsagePulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexModifierMock) MinimockFinish() {

	if !m.SetFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.Set")
	}

	if !m.UpdateUsagePulseFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.UpdateUsagePulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetFinished()
		ok = ok && m.UpdateUsagePulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetFinished() {
				m.t.Error("Expected call to IndexModifierMock.Set")
			}

			if !m.UpdateUsagePulseFinished() {
				m.t.Error("Expected call to IndexModifierMock.UpdateUsagePulse")
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
func (m *IndexModifierMock) AllMocksCalled() bool {

	if !m.SetFinished() {
		return false
	}

	if !m.UpdateUsagePulseFinished() {
		return false
	}

	return true
}

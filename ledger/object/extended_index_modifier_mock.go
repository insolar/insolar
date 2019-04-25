package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ExtendedIndexModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ExtendedIndexModifierMock implements github.com/insolar/insolar/ledger/object.ExtendedIndexModifier
type ExtendedIndexModifierMock struct {
	t minimock.Tester

	SetUsageForPulseFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber)
	SetUsageForPulseCounter    uint64
	SetUsageForPulsePreCounter uint64
	SetUsageForPulseMock       mExtendedIndexModifierMockSetUsageForPulse

	SetWithMetaFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) (r error)
	SetWithMetaCounter    uint64
	SetWithMetaPreCounter uint64
	SetWithMetaMock       mExtendedIndexModifierMockSetWithMeta
}

//NewExtendedIndexModifierMock returns a mock for github.com/insolar/insolar/ledger/object.ExtendedIndexModifier
func NewExtendedIndexModifierMock(t minimock.Tester) *ExtendedIndexModifierMock {
	m := &ExtendedIndexModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetUsageForPulseMock = mExtendedIndexModifierMockSetUsageForPulse{mock: m}
	m.SetWithMetaMock = mExtendedIndexModifierMockSetWithMeta{mock: m}

	return m
}

type mExtendedIndexModifierMockSetUsageForPulse struct {
	mock              *ExtendedIndexModifierMock
	mainExpectation   *ExtendedIndexModifierMockSetUsageForPulseExpectation
	expectationSeries []*ExtendedIndexModifierMockSetUsageForPulseExpectation
}

type ExtendedIndexModifierMockSetUsageForPulseExpectation struct {
	input *ExtendedIndexModifierMockSetUsageForPulseInput
}

type ExtendedIndexModifierMockSetUsageForPulseInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

//Expect specifies that invocation of ExtendedIndexModifier.SetUsageForPulse is expected from 1 to Infinity times
func (m *mExtendedIndexModifierMockSetUsageForPulse) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mExtendedIndexModifierMockSetUsageForPulse {
	m.mock.SetUsageForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedIndexModifierMockSetUsageForPulseExpectation{}
	}
	m.mainExpectation.input = &ExtendedIndexModifierMockSetUsageForPulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ExtendedIndexModifier.SetUsageForPulse
func (m *mExtendedIndexModifierMockSetUsageForPulse) Return() *ExtendedIndexModifierMock {
	m.mock.SetUsageForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedIndexModifierMockSetUsageForPulseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExtendedIndexModifier.SetUsageForPulse is expected once
func (m *mExtendedIndexModifierMockSetUsageForPulse) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *ExtendedIndexModifierMockSetUsageForPulseExpectation {
	m.mock.SetUsageForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ExtendedIndexModifierMockSetUsageForPulseExpectation{}
	expectation.input = &ExtendedIndexModifierMockSetUsageForPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExtendedIndexModifier.SetUsageForPulse method
func (m *mExtendedIndexModifierMockSetUsageForPulse) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber)) *ExtendedIndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetUsageForPulseFunc = f
	return m.mock
}

//SetUsageForPulse implements github.com/insolar/insolar/ledger/object.ExtendedIndexModifier interface
func (m *ExtendedIndexModifierMock) SetUsageForPulse(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.SetUsageForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetUsageForPulseCounter, 1)

	if len(m.SetUsageForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetUsageForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExtendedIndexModifierMock.SetUsageForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetUsageForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExtendedIndexModifierMockSetUsageForPulseInput{p, p1, p2}, "ExtendedIndexModifier.SetUsageForPulse got unexpected parameters")

		return
	}

	if m.SetUsageForPulseMock.mainExpectation != nil {

		input := m.SetUsageForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExtendedIndexModifierMockSetUsageForPulseInput{p, p1, p2}, "ExtendedIndexModifier.SetUsageForPulse got unexpected parameters")
		}

		return
	}

	if m.SetUsageForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ExtendedIndexModifierMock.SetUsageForPulse. %v %v %v", p, p1, p2)
		return
	}

	m.SetUsageForPulseFunc(p, p1, p2)
}

//SetUsageForPulseMinimockCounter returns a count of ExtendedIndexModifierMock.SetUsageForPulseFunc invocations
func (m *ExtendedIndexModifierMock) SetUsageForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetUsageForPulseCounter)
}

//SetUsageForPulseMinimockPreCounter returns the value of ExtendedIndexModifierMock.SetUsageForPulse invocations
func (m *ExtendedIndexModifierMock) SetUsageForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetUsageForPulsePreCounter)
}

//SetUsageForPulseFinished returns true if mock invocations count is ok
func (m *ExtendedIndexModifierMock) SetUsageForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetUsageForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetUsageForPulseCounter) == uint64(len(m.SetUsageForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetUsageForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetUsageForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetUsageForPulseFunc != nil {
		return atomic.LoadUint64(&m.SetUsageForPulseCounter) > 0
	}

	return true
}

type mExtendedIndexModifierMockSetWithMeta struct {
	mock              *ExtendedIndexModifierMock
	mainExpectation   *ExtendedIndexModifierMockSetWithMetaExpectation
	expectationSeries []*ExtendedIndexModifierMockSetWithMetaExpectation
}

type ExtendedIndexModifierMockSetWithMetaExpectation struct {
	input  *ExtendedIndexModifierMockSetWithMetaInput
	result *ExtendedIndexModifierMockSetWithMetaResult
}

type ExtendedIndexModifierMockSetWithMetaInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 Lifeline
}

type ExtendedIndexModifierMockSetWithMetaResult struct {
	r error
}

//Expect specifies that invocation of ExtendedIndexModifier.SetWithMeta is expected from 1 to Infinity times
func (m *mExtendedIndexModifierMockSetWithMeta) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) *mExtendedIndexModifierMockSetWithMeta {
	m.mock.SetWithMetaFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedIndexModifierMockSetWithMetaExpectation{}
	}
	m.mainExpectation.input = &ExtendedIndexModifierMockSetWithMetaInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ExtendedIndexModifier.SetWithMeta
func (m *mExtendedIndexModifierMockSetWithMeta) Return(r error) *ExtendedIndexModifierMock {
	m.mock.SetWithMetaFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExtendedIndexModifierMockSetWithMetaExpectation{}
	}
	m.mainExpectation.result = &ExtendedIndexModifierMockSetWithMetaResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExtendedIndexModifier.SetWithMeta is expected once
func (m *mExtendedIndexModifierMockSetWithMeta) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) *ExtendedIndexModifierMockSetWithMetaExpectation {
	m.mock.SetWithMetaFunc = nil
	m.mainExpectation = nil

	expectation := &ExtendedIndexModifierMockSetWithMetaExpectation{}
	expectation.input = &ExtendedIndexModifierMockSetWithMetaInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExtendedIndexModifierMockSetWithMetaExpectation) Return(r error) {
	e.result = &ExtendedIndexModifierMockSetWithMetaResult{r}
}

//Set uses given function f as a mock of ExtendedIndexModifier.SetWithMeta method
func (m *mExtendedIndexModifierMockSetWithMeta) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) (r error)) *ExtendedIndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetWithMetaFunc = f
	return m.mock
}

//SetWithMeta implements github.com/insolar/insolar/ledger/object.ExtendedIndexModifier interface
func (m *ExtendedIndexModifierMock) SetWithMeta(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetWithMetaPreCounter, 1)
	defer atomic.AddUint64(&m.SetWithMetaCounter, 1)

	if len(m.SetWithMetaMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetWithMetaMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExtendedIndexModifierMock.SetWithMeta. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetWithMetaMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExtendedIndexModifierMockSetWithMetaInput{p, p1, p2, p3}, "ExtendedIndexModifier.SetWithMeta got unexpected parameters")

		result := m.SetWithMetaMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExtendedIndexModifierMock.SetWithMeta")
			return
		}

		r = result.r

		return
	}

	if m.SetWithMetaMock.mainExpectation != nil {

		input := m.SetWithMetaMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExtendedIndexModifierMockSetWithMetaInput{p, p1, p2, p3}, "ExtendedIndexModifier.SetWithMeta got unexpected parameters")
		}

		result := m.SetWithMetaMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExtendedIndexModifierMock.SetWithMeta")
		}

		r = result.r

		return
	}

	if m.SetWithMetaFunc == nil {
		m.t.Fatalf("Unexpected call to ExtendedIndexModifierMock.SetWithMeta. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetWithMetaFunc(p, p1, p2, p3)
}

//SetWithMetaMinimockCounter returns a count of ExtendedIndexModifierMock.SetWithMetaFunc invocations
func (m *ExtendedIndexModifierMock) SetWithMetaMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetWithMetaCounter)
}

//SetWithMetaMinimockPreCounter returns the value of ExtendedIndexModifierMock.SetWithMeta invocations
func (m *ExtendedIndexModifierMock) SetWithMetaMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetWithMetaPreCounter)
}

//SetWithMetaFinished returns true if mock invocations count is ok
func (m *ExtendedIndexModifierMock) SetWithMetaFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetWithMetaMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetWithMetaCounter) == uint64(len(m.SetWithMetaMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetWithMetaMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetWithMetaCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetWithMetaFunc != nil {
		return atomic.LoadUint64(&m.SetWithMetaCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExtendedIndexModifierMock) ValidateCallCounters() {

	if !m.SetUsageForPulseFinished() {
		m.t.Fatal("Expected call to ExtendedIndexModifierMock.SetUsageForPulse")
	}

	if !m.SetWithMetaFinished() {
		m.t.Fatal("Expected call to ExtendedIndexModifierMock.SetWithMeta")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExtendedIndexModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ExtendedIndexModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ExtendedIndexModifierMock) MinimockFinish() {

	if !m.SetUsageForPulseFinished() {
		m.t.Fatal("Expected call to ExtendedIndexModifierMock.SetUsageForPulse")
	}

	if !m.SetWithMetaFinished() {
		m.t.Fatal("Expected call to ExtendedIndexModifierMock.SetWithMeta")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ExtendedIndexModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ExtendedIndexModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetUsageForPulseFinished()
		ok = ok && m.SetWithMetaFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetUsageForPulseFinished() {
				m.t.Error("Expected call to ExtendedIndexModifierMock.SetUsageForPulse")
			}

			if !m.SetWithMetaFinished() {
				m.t.Error("Expected call to ExtendedIndexModifierMock.SetWithMeta")
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
func (m *ExtendedIndexModifierMock) AllMocksCalled() bool {

	if !m.SetUsageForPulseFinished() {
		return false
	}

	if !m.SetWithMetaFinished() {
		return false
	}

	return true
}

package drop

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Accessor" can be found in github.com/insolar/insolar/ledger/drop
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//AccessorMock implements github.com/insolar/insolar/ledger/drop.Accessor
type AccessorMock struct {
	t minimock.Tester

	ForPulseFunc       func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r Drop, r1 error)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mAccessorMockForPulse
}

//NewAccessorMock returns a mock for github.com/insolar/insolar/ledger/drop.Accessor
func NewAccessorMock(t minimock.Tester) *AccessorMock {
	m := &AccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseMock = mAccessorMockForPulse{mock: m}

	return m
}

type mAccessorMockForPulse struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockForPulseExpectation
	expectationSeries []*AccessorMockForPulseExpectation
}

type AccessorMockForPulseExpectation struct {
	input  *AccessorMockForPulseInput
	result *AccessorMockForPulseResult
}

type AccessorMockForPulseInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 insolar.PulseNumber
}

type AccessorMockForPulseResult struct {
	r  Drop
	r1 error
}

//Expect specifies that invocation of Accessor.ForPulse is expected from 1 to Infinity times
func (m *mAccessorMockForPulse) Expect(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *mAccessorMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockForPulseExpectation{}
	}
	m.mainExpectation.input = &AccessorMockForPulseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Accessor.ForPulse
func (m *mAccessorMockForPulse) Return(r Drop, r1 error) *AccessorMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockForPulseExpectation{}
	}
	m.mainExpectation.result = &AccessorMockForPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.ForPulse is expected once
func (m *mAccessorMockForPulse) ExpectOnce(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *AccessorMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockForPulseExpectation{}
	expectation.input = &AccessorMockForPulseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockForPulseExpectation) Return(r Drop, r1 error) {
	e.result = &AccessorMockForPulseResult{r, r1}
}

//Set uses given function f as a mock of Accessor.ForPulse method
func (m *mAccessorMockForPulse) Set(f func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r Drop, r1 error)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/drop.Accessor interface
func (m *AccessorMock) ForPulse(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r Drop, r1 error) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.ForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockForPulseInput{p, p1, p2}, "Accessor.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.ForPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockForPulseInput{p, p1, p2}, "Accessor.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.ForPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.ForPulse. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPulseFunc(p, p1, p2)
}

//ForPulseMinimockCounter returns a count of AccessorMock.ForPulseFunc invocations
func (m *AccessorMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

//ForPulseMinimockPreCounter returns the value of AccessorMock.ForPulse invocations
func (m *AccessorMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

//ForPulseFinished returns true if mock invocations count is ok
func (m *AccessorMock) ForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseCounter) == uint64(len(m.ForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseFunc != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) ValidateCallCounters() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to AccessorMock.ForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AccessorMock) MinimockFinish() {

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to AccessorMock.ForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseFinished() {
				m.t.Error("Expected call to AccessorMock.ForPulse")
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
func (m *AccessorMock) AllMocksCalled() bool {

	if !m.ForPulseFinished() {
		return false
	}

	return true
}

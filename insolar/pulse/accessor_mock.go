package pulse

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Accessor" can be found in github.com/insolar/insolar/insolar/pulse
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//AccessorMock implements github.com/insolar/insolar/insolar/pulse.Accessor
type AccessorMock struct {
	t minimock.Tester

	ForPulseNumberFunc       func(p context.Context, p1 insolar.PulseNumber) (r insolar.Pulse, r1 error)
	ForPulseNumberCounter    uint64
	ForPulseNumberPreCounter uint64
	ForPulseNumberMock       mAccessorMockForPulseNumber

	LatestFunc       func(p context.Context) (r insolar.Pulse, r1 error)
	LatestCounter    uint64
	LatestPreCounter uint64
	LatestMock       mAccessorMockLatest
}

//NewAccessorMock returns a mock for github.com/insolar/insolar/insolar/pulse.Accessor
func NewAccessorMock(t minimock.Tester) *AccessorMock {
	m := &AccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseNumberMock = mAccessorMockForPulseNumber{mock: m}
	m.LatestMock = mAccessorMockLatest{mock: m}

	return m
}

type mAccessorMockForPulseNumber struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockForPulseNumberExpectation
	expectationSeries []*AccessorMockForPulseNumberExpectation
}

type AccessorMockForPulseNumberExpectation struct {
	input  *AccessorMockForPulseNumberInput
	result *AccessorMockForPulseNumberResult
}

type AccessorMockForPulseNumberInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type AccessorMockForPulseNumberResult struct {
	r  insolar.Pulse
	r1 error
}

//Expect specifies that invocation of Accessor.ForPulseNumber is expected from 1 to Infinity times
func (m *mAccessorMockForPulseNumber) Expect(p context.Context, p1 insolar.PulseNumber) *mAccessorMockForPulseNumber {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.input = &AccessorMockForPulseNumberInput{p, p1}
	return m
}

//Return specifies results of invocation of Accessor.ForPulseNumber
func (m *mAccessorMockForPulseNumber) Return(r insolar.Pulse, r1 error) *AccessorMock {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.result = &AccessorMockForPulseNumberResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.ForPulseNumber is expected once
func (m *mAccessorMockForPulseNumber) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *AccessorMockForPulseNumberExpectation {
	m.mock.ForPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockForPulseNumberExpectation{}
	expectation.input = &AccessorMockForPulseNumberInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockForPulseNumberExpectation) Return(r insolar.Pulse, r1 error) {
	e.result = &AccessorMockForPulseNumberResult{r, r1}
}

//Set uses given function f as a mock of Accessor.ForPulseNumber method
func (m *mAccessorMockForPulseNumber) Set(f func(p context.Context, p1 insolar.PulseNumber) (r insolar.Pulse, r1 error)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseNumberFunc = f
	return m.mock
}

//ForPulseNumber implements github.com/insolar/insolar/insolar/pulse.Accessor interface
func (m *AccessorMock) ForPulseNumber(p context.Context, p1 insolar.PulseNumber) (r insolar.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.ForPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseNumberCounter, 1)

	if len(m.ForPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.ForPulseNumber. %v %v", p, p1)
			return
		}

		input := m.ForPulseNumberMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockForPulseNumberInput{p, p1}, "Accessor.ForPulseNumber got unexpected parameters")

		result := m.ForPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.ForPulseNumber")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberMock.mainExpectation != nil {

		input := m.ForPulseNumberMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockForPulseNumberInput{p, p1}, "Accessor.ForPulseNumber got unexpected parameters")
		}

		result := m.ForPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.ForPulseNumber")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.ForPulseNumber. %v %v", p, p1)
		return
	}

	return m.ForPulseNumberFunc(p, p1)
}

//ForPulseNumberMinimockCounter returns a count of AccessorMock.ForPulseNumberFunc invocations
func (m *AccessorMock) ForPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberCounter)
}

//ForPulseNumberMinimockPreCounter returns the value of AccessorMock.ForPulseNumber invocations
func (m *AccessorMock) ForPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberPreCounter)
}

//ForPulseNumberFinished returns true if mock invocations count is ok
func (m *AccessorMock) ForPulseNumberFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseNumberMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) == uint64(len(m.ForPulseNumberMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseNumberMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseNumberFunc != nil {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) > 0
	}

	return true
}

type mAccessorMockLatest struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockLatestExpectation
	expectationSeries []*AccessorMockLatestExpectation
}

type AccessorMockLatestExpectation struct {
	input  *AccessorMockLatestInput
	result *AccessorMockLatestResult
}

type AccessorMockLatestInput struct {
	p context.Context
}

type AccessorMockLatestResult struct {
	r  insolar.Pulse
	r1 error
}

//Expect specifies that invocation of Accessor.Latest is expected from 1 to Infinity times
func (m *mAccessorMockLatest) Expect(p context.Context) *mAccessorMockLatest {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockLatestExpectation{}
	}
	m.mainExpectation.input = &AccessorMockLatestInput{p}
	return m
}

//Return specifies results of invocation of Accessor.Latest
func (m *mAccessorMockLatest) Return(r insolar.Pulse, r1 error) *AccessorMock {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockLatestExpectation{}
	}
	m.mainExpectation.result = &AccessorMockLatestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.Latest is expected once
func (m *mAccessorMockLatest) ExpectOnce(p context.Context) *AccessorMockLatestExpectation {
	m.mock.LatestFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockLatestExpectation{}
	expectation.input = &AccessorMockLatestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockLatestExpectation) Return(r insolar.Pulse, r1 error) {
	e.result = &AccessorMockLatestResult{r, r1}
}

//Set uses given function f as a mock of Accessor.Latest method
func (m *mAccessorMockLatest) Set(f func(p context.Context) (r insolar.Pulse, r1 error)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LatestFunc = f
	return m.mock
}

//Latest implements github.com/insolar/insolar/insolar/pulse.Accessor interface
func (m *AccessorMock) Latest(p context.Context) (r insolar.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.LatestPreCounter, 1)
	defer atomic.AddUint64(&m.LatestCounter, 1)

	if len(m.LatestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LatestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.Latest. %v", p)
			return
		}

		input := m.LatestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockLatestInput{p}, "Accessor.Latest got unexpected parameters")

		result := m.LatestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.Latest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestMock.mainExpectation != nil {

		input := m.LatestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockLatestInput{p}, "Accessor.Latest got unexpected parameters")
		}

		result := m.LatestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.Latest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.Latest. %v", p)
		return
	}

	return m.LatestFunc(p)
}

//LatestMinimockCounter returns a count of AccessorMock.LatestFunc invocations
func (m *AccessorMock) LatestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LatestCounter)
}

//LatestMinimockPreCounter returns the value of AccessorMock.Latest invocations
func (m *AccessorMock) LatestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LatestPreCounter)
}

//LatestFinished returns true if mock invocations count is ok
func (m *AccessorMock) LatestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LatestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LatestCounter) == uint64(len(m.LatestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LatestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LatestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LatestFunc != nil {
		return atomic.LoadUint64(&m.LatestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) ValidateCallCounters() {

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to AccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to AccessorMock.Latest")
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

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to AccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to AccessorMock.Latest")
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
		ok = ok && m.ForPulseNumberFinished()
		ok = ok && m.LatestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseNumberFinished() {
				m.t.Error("Expected call to AccessorMock.ForPulseNumber")
			}

			if !m.LatestFinished() {
				m.t.Error("Expected call to AccessorMock.Latest")
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

	if !m.ForPulseNumberFinished() {
		return false
	}

	if !m.LatestFinished() {
		return false
	}

	return true
}

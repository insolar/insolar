package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseAccessor" can be found in github.com/insolar/insolar/network/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseAccessorMock implements github.com/insolar/insolar/network/storage.PulseAccessor
type PulseAccessorMock struct {
	t minimock.Tester

	ForPulseNumberFunc       func(p context.Context, p1 insolar.PulseNumber) (r insolar.Pulse, r1 error)
	ForPulseNumberCounter    uint64
	ForPulseNumberPreCounter uint64
	ForPulseNumberMock       mPulseAccessorMockForPulseNumber

	LatestFunc       func(p context.Context) (r insolar.Pulse, r1 error)
	LatestCounter    uint64
	LatestPreCounter uint64
	LatestMock       mPulseAccessorMockLatest
}

//NewPulseAccessorMock returns a mock for github.com/insolar/insolar/network/storage.PulseAccessor
func NewPulseAccessorMock(t minimock.Tester) *PulseAccessorMock {
	m := &PulseAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseNumberMock = mPulseAccessorMockForPulseNumber{mock: m}
	m.LatestMock = mPulseAccessorMockLatest{mock: m}

	return m
}

type mPulseAccessorMockForPulseNumber struct {
	mock              *PulseAccessorMock
	mainExpectation   *PulseAccessorMockForPulseNumberExpectation
	expectationSeries []*PulseAccessorMockForPulseNumberExpectation
}

type PulseAccessorMockForPulseNumberExpectation struct {
	input  *PulseAccessorMockForPulseNumberInput
	result *PulseAccessorMockForPulseNumberResult
}

type PulseAccessorMockForPulseNumberInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type PulseAccessorMockForPulseNumberResult struct {
	r  insolar.Pulse
	r1 error
}

//Expect specifies that invocation of PulseAccessor.ForPulseNumber is expected from 1 to Infinity times
func (m *mPulseAccessorMockForPulseNumber) Expect(p context.Context, p1 insolar.PulseNumber) *mPulseAccessorMockForPulseNumber {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseAccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.input = &PulseAccessorMockForPulseNumberInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseAccessor.ForPulseNumber
func (m *mPulseAccessorMockForPulseNumber) Return(r insolar.Pulse, r1 error) *PulseAccessorMock {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseAccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.result = &PulseAccessorMockForPulseNumberResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseAccessor.ForPulseNumber is expected once
func (m *mPulseAccessorMockForPulseNumber) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *PulseAccessorMockForPulseNumberExpectation {
	m.mock.ForPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &PulseAccessorMockForPulseNumberExpectation{}
	expectation.input = &PulseAccessorMockForPulseNumberInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseAccessorMockForPulseNumberExpectation) Return(r insolar.Pulse, r1 error) {
	e.result = &PulseAccessorMockForPulseNumberResult{r, r1}
}

//Set uses given function f as a mock of PulseAccessor.ForPulseNumber method
func (m *mPulseAccessorMockForPulseNumber) Set(f func(p context.Context, p1 insolar.PulseNumber) (r insolar.Pulse, r1 error)) *PulseAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseNumberFunc = f
	return m.mock
}

//ForPulseNumber implements github.com/insolar/insolar/network/storage.PulseAccessor interface
func (m *PulseAccessorMock) ForPulseNumber(p context.Context, p1 insolar.PulseNumber) (r insolar.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.ForPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseNumberCounter, 1)

	if len(m.ForPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseAccessorMock.ForPulseNumber. %v %v", p, p1)
			return
		}

		input := m.ForPulseNumberMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseAccessorMockForPulseNumberInput{p, p1}, "PulseAccessor.ForPulseNumber got unexpected parameters")

		result := m.ForPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseAccessorMock.ForPulseNumber")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberMock.mainExpectation != nil {

		input := m.ForPulseNumberMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseAccessorMockForPulseNumberInput{p, p1}, "PulseAccessor.ForPulseNumber got unexpected parameters")
		}

		result := m.ForPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseAccessorMock.ForPulseNumber")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to PulseAccessorMock.ForPulseNumber. %v %v", p, p1)
		return
	}

	return m.ForPulseNumberFunc(p, p1)
}

//ForPulseNumberMinimockCounter returns a count of PulseAccessorMock.ForPulseNumberFunc invocations
func (m *PulseAccessorMock) ForPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberCounter)
}

//ForPulseNumberMinimockPreCounter returns the value of PulseAccessorMock.ForPulseNumber invocations
func (m *PulseAccessorMock) ForPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberPreCounter)
}

//ForPulseNumberFinished returns true if mock invocations count is ok
func (m *PulseAccessorMock) ForPulseNumberFinished() bool {
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

type mPulseAccessorMockLatest struct {
	mock              *PulseAccessorMock
	mainExpectation   *PulseAccessorMockLatestExpectation
	expectationSeries []*PulseAccessorMockLatestExpectation
}

type PulseAccessorMockLatestExpectation struct {
	input  *PulseAccessorMockLatestInput
	result *PulseAccessorMockLatestResult
}

type PulseAccessorMockLatestInput struct {
	p context.Context
}

type PulseAccessorMockLatestResult struct {
	r  insolar.Pulse
	r1 error
}

//Expect specifies that invocation of PulseAccessor.Latest is expected from 1 to Infinity times
func (m *mPulseAccessorMockLatest) Expect(p context.Context) *mPulseAccessorMockLatest {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseAccessorMockLatestExpectation{}
	}
	m.mainExpectation.input = &PulseAccessorMockLatestInput{p}
	return m
}

//Return specifies results of invocation of PulseAccessor.Latest
func (m *mPulseAccessorMockLatest) Return(r insolar.Pulse, r1 error) *PulseAccessorMock {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseAccessorMockLatestExpectation{}
	}
	m.mainExpectation.result = &PulseAccessorMockLatestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseAccessor.Latest is expected once
func (m *mPulseAccessorMockLatest) ExpectOnce(p context.Context) *PulseAccessorMockLatestExpectation {
	m.mock.LatestFunc = nil
	m.mainExpectation = nil

	expectation := &PulseAccessorMockLatestExpectation{}
	expectation.input = &PulseAccessorMockLatestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseAccessorMockLatestExpectation) Return(r insolar.Pulse, r1 error) {
	e.result = &PulseAccessorMockLatestResult{r, r1}
}

//Set uses given function f as a mock of PulseAccessor.Latest method
func (m *mPulseAccessorMockLatest) Set(f func(p context.Context) (r insolar.Pulse, r1 error)) *PulseAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LatestFunc = f
	return m.mock
}

//Latest implements github.com/insolar/insolar/network/storage.PulseAccessor interface
func (m *PulseAccessorMock) Latest(p context.Context) (r insolar.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.LatestPreCounter, 1)
	defer atomic.AddUint64(&m.LatestCounter, 1)

	if len(m.LatestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LatestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseAccessorMock.Latest. %v", p)
			return
		}

		input := m.LatestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseAccessorMockLatestInput{p}, "PulseAccessor.Latest got unexpected parameters")

		result := m.LatestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseAccessorMock.Latest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestMock.mainExpectation != nil {

		input := m.LatestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseAccessorMockLatestInput{p}, "PulseAccessor.Latest got unexpected parameters")
		}

		result := m.LatestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseAccessorMock.Latest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestFunc == nil {
		m.t.Fatalf("Unexpected call to PulseAccessorMock.Latest. %v", p)
		return
	}

	return m.LatestFunc(p)
}

//LatestMinimockCounter returns a count of PulseAccessorMock.LatestFunc invocations
func (m *PulseAccessorMock) LatestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LatestCounter)
}

//LatestMinimockPreCounter returns the value of PulseAccessorMock.Latest invocations
func (m *PulseAccessorMock) LatestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LatestPreCounter)
}

//LatestFinished returns true if mock invocations count is ok
func (m *PulseAccessorMock) LatestFinished() bool {
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
func (m *PulseAccessorMock) ValidateCallCounters() {

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to PulseAccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to PulseAccessorMock.Latest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseAccessorMock) MinimockFinish() {

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to PulseAccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to PulseAccessorMock.Latest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseAccessorMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to PulseAccessorMock.ForPulseNumber")
			}

			if !m.LatestFinished() {
				m.t.Error("Expected call to PulseAccessorMock.Latest")
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
func (m *PulseAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseNumberFinished() {
		return false
	}

	if !m.LatestFinished() {
		return false
	}

	return true
}

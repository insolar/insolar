package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SwitcherWorkAround" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//SwitcherWorkAroundMock implements github.com/insolar/insolar/core.SwitcherWorkAround
type SwitcherWorkAroundMock struct {
	t minimock.Tester

	IsBootstrappedFunc       func() (r bool)
	IsBootstrappedCounter    uint64
	IsBootstrappedPreCounter uint64
	IsBootstrappedMock       mSwitcherWorkAroundMockIsBootstrapped

	SetIsBootstrappedFunc       func(p bool)
	SetIsBootstrappedCounter    uint64
	SetIsBootstrappedPreCounter uint64
	SetIsBootstrappedMock       mSwitcherWorkAroundMockSetIsBootstrapped
}

//NewSwitcherWorkAroundMock returns a mock for github.com/insolar/insolar/core.SwitcherWorkAround
func NewSwitcherWorkAroundMock(t minimock.Tester) *SwitcherWorkAroundMock {
	m := &SwitcherWorkAroundMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IsBootstrappedMock = mSwitcherWorkAroundMockIsBootstrapped{mock: m}
	m.SetIsBootstrappedMock = mSwitcherWorkAroundMockSetIsBootstrapped{mock: m}

	return m
}

type mSwitcherWorkAroundMockIsBootstrapped struct {
	mock              *SwitcherWorkAroundMock
	mainExpectation   *SwitcherWorkAroundMockIsBootstrappedExpectation
	expectationSeries []*SwitcherWorkAroundMockIsBootstrappedExpectation
}

type SwitcherWorkAroundMockIsBootstrappedExpectation struct {
	result *SwitcherWorkAroundMockIsBootstrappedResult
}

type SwitcherWorkAroundMockIsBootstrappedResult struct {
	r bool
}

//Expect specifies that invocation of SwitcherWorkAround.IsBootstrapped is expected from 1 to Infinity times
func (m *mSwitcherWorkAroundMockIsBootstrapped) Expect() *mSwitcherWorkAroundMockIsBootstrapped {
	m.mock.IsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SwitcherWorkAroundMockIsBootstrappedExpectation{}
	}

	return m
}

//Return specifies results of invocation of SwitcherWorkAround.IsBootstrapped
func (m *mSwitcherWorkAroundMockIsBootstrapped) Return(r bool) *SwitcherWorkAroundMock {
	m.mock.IsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SwitcherWorkAroundMockIsBootstrappedExpectation{}
	}
	m.mainExpectation.result = &SwitcherWorkAroundMockIsBootstrappedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SwitcherWorkAround.IsBootstrapped is expected once
func (m *mSwitcherWorkAroundMockIsBootstrapped) ExpectOnce() *SwitcherWorkAroundMockIsBootstrappedExpectation {
	m.mock.IsBootstrappedFunc = nil
	m.mainExpectation = nil

	expectation := &SwitcherWorkAroundMockIsBootstrappedExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SwitcherWorkAroundMockIsBootstrappedExpectation) Return(r bool) {
	e.result = &SwitcherWorkAroundMockIsBootstrappedResult{r}
}

//Set uses given function f as a mock of SwitcherWorkAround.IsBootstrapped method
func (m *mSwitcherWorkAroundMockIsBootstrapped) Set(f func() (r bool)) *SwitcherWorkAroundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsBootstrappedFunc = f
	return m.mock
}

//IsBootstrapped implements github.com/insolar/insolar/core.SwitcherWorkAround interface
func (m *SwitcherWorkAroundMock) IsBootstrapped() (r bool) {
	counter := atomic.AddUint64(&m.IsBootstrappedPreCounter, 1)
	defer atomic.AddUint64(&m.IsBootstrappedCounter, 1)

	if len(m.IsBootstrappedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsBootstrappedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SwitcherWorkAroundMock.IsBootstrapped.")
			return
		}

		result := m.IsBootstrappedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SwitcherWorkAroundMock.IsBootstrapped")
			return
		}

		r = result.r

		return
	}

	if m.IsBootstrappedMock.mainExpectation != nil {

		result := m.IsBootstrappedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SwitcherWorkAroundMock.IsBootstrapped")
		}

		r = result.r

		return
	}

	if m.IsBootstrappedFunc == nil {
		m.t.Fatalf("Unexpected call to SwitcherWorkAroundMock.IsBootstrapped.")
		return
	}

	return m.IsBootstrappedFunc()
}

//IsBootstrappedMinimockCounter returns a count of SwitcherWorkAroundMock.IsBootstrappedFunc invocations
func (m *SwitcherWorkAroundMock) IsBootstrappedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsBootstrappedCounter)
}

//IsBootstrappedMinimockPreCounter returns the value of SwitcherWorkAroundMock.IsBootstrapped invocations
func (m *SwitcherWorkAroundMock) IsBootstrappedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsBootstrappedPreCounter)
}

//IsBootstrappedFinished returns true if mock invocations count is ok
func (m *SwitcherWorkAroundMock) IsBootstrappedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsBootstrappedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsBootstrappedCounter) == uint64(len(m.IsBootstrappedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsBootstrappedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsBootstrappedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsBootstrappedFunc != nil {
		return atomic.LoadUint64(&m.IsBootstrappedCounter) > 0
	}

	return true
}

type mSwitcherWorkAroundMockSetIsBootstrapped struct {
	mock              *SwitcherWorkAroundMock
	mainExpectation   *SwitcherWorkAroundMockSetIsBootstrappedExpectation
	expectationSeries []*SwitcherWorkAroundMockSetIsBootstrappedExpectation
}

type SwitcherWorkAroundMockSetIsBootstrappedExpectation struct {
	input *SwitcherWorkAroundMockSetIsBootstrappedInput
}

type SwitcherWorkAroundMockSetIsBootstrappedInput struct {
	p bool
}

//Expect specifies that invocation of SwitcherWorkAround.SetIsBootstrapped is expected from 1 to Infinity times
func (m *mSwitcherWorkAroundMockSetIsBootstrapped) Expect(p bool) *mSwitcherWorkAroundMockSetIsBootstrapped {
	m.mock.SetIsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SwitcherWorkAroundMockSetIsBootstrappedExpectation{}
	}
	m.mainExpectation.input = &SwitcherWorkAroundMockSetIsBootstrappedInput{p}
	return m
}

//Return specifies results of invocation of SwitcherWorkAround.SetIsBootstrapped
func (m *mSwitcherWorkAroundMockSetIsBootstrapped) Return() *SwitcherWorkAroundMock {
	m.mock.SetIsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SwitcherWorkAroundMockSetIsBootstrappedExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of SwitcherWorkAround.SetIsBootstrapped is expected once
func (m *mSwitcherWorkAroundMockSetIsBootstrapped) ExpectOnce(p bool) *SwitcherWorkAroundMockSetIsBootstrappedExpectation {
	m.mock.SetIsBootstrappedFunc = nil
	m.mainExpectation = nil

	expectation := &SwitcherWorkAroundMockSetIsBootstrappedExpectation{}
	expectation.input = &SwitcherWorkAroundMockSetIsBootstrappedInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of SwitcherWorkAround.SetIsBootstrapped method
func (m *mSwitcherWorkAroundMockSetIsBootstrapped) Set(f func(p bool)) *SwitcherWorkAroundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetIsBootstrappedFunc = f
	return m.mock
}

//SetIsBootstrapped implements github.com/insolar/insolar/core.SwitcherWorkAround interface
func (m *SwitcherWorkAroundMock) SetIsBootstrapped(p bool) {
	counter := atomic.AddUint64(&m.SetIsBootstrappedPreCounter, 1)
	defer atomic.AddUint64(&m.SetIsBootstrappedCounter, 1)

	if len(m.SetIsBootstrappedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetIsBootstrappedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SwitcherWorkAroundMock.SetIsBootstrapped. %v", p)
			return
		}

		input := m.SetIsBootstrappedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SwitcherWorkAroundMockSetIsBootstrappedInput{p}, "SwitcherWorkAround.SetIsBootstrapped got unexpected parameters")

		return
	}

	if m.SetIsBootstrappedMock.mainExpectation != nil {

		input := m.SetIsBootstrappedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SwitcherWorkAroundMockSetIsBootstrappedInput{p}, "SwitcherWorkAround.SetIsBootstrapped got unexpected parameters")
		}

		return
	}

	if m.SetIsBootstrappedFunc == nil {
		m.t.Fatalf("Unexpected call to SwitcherWorkAroundMock.SetIsBootstrapped. %v", p)
		return
	}

	m.SetIsBootstrappedFunc(p)
}

//SetIsBootstrappedMinimockCounter returns a count of SwitcherWorkAroundMock.SetIsBootstrappedFunc invocations
func (m *SwitcherWorkAroundMock) SetIsBootstrappedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetIsBootstrappedCounter)
}

//SetIsBootstrappedMinimockPreCounter returns the value of SwitcherWorkAroundMock.SetIsBootstrapped invocations
func (m *SwitcherWorkAroundMock) SetIsBootstrappedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetIsBootstrappedPreCounter)
}

//SetIsBootstrappedFinished returns true if mock invocations count is ok
func (m *SwitcherWorkAroundMock) SetIsBootstrappedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetIsBootstrappedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetIsBootstrappedCounter) == uint64(len(m.SetIsBootstrappedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetIsBootstrappedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetIsBootstrappedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetIsBootstrappedFunc != nil {
		return atomic.LoadUint64(&m.SetIsBootstrappedCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SwitcherWorkAroundMock) ValidateCallCounters() {

	if !m.IsBootstrappedFinished() {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.IsBootstrapped")
	}

	if !m.SetIsBootstrappedFinished() {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.SetIsBootstrapped")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SwitcherWorkAroundMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SwitcherWorkAroundMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SwitcherWorkAroundMock) MinimockFinish() {

	if !m.IsBootstrappedFinished() {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.IsBootstrapped")
	}

	if !m.SetIsBootstrappedFinished() {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.SetIsBootstrapped")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SwitcherWorkAroundMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SwitcherWorkAroundMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.IsBootstrappedFinished()
		ok = ok && m.SetIsBootstrappedFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IsBootstrappedFinished() {
				m.t.Error("Expected call to SwitcherWorkAroundMock.IsBootstrapped")
			}

			if !m.SetIsBootstrappedFinished() {
				m.t.Error("Expected call to SwitcherWorkAroundMock.SetIsBootstrapped")
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
func (m *SwitcherWorkAroundMock) AllMocksCalled() bool {

	if !m.IsBootstrappedFinished() {
		return false
	}

	if !m.SetIsBootstrappedFinished() {
		return false
	}

	return true
}

package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseStorage" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseStorageMock implements github.com/insolar/insolar/core.PulseStorage
type PulseStorageMock struct {
	t minimock.Tester

	CurrentFunc       func(p context.Context) (r *core.Pulse, r1 error)
	CurrentCounter    uint64
	CurrentPreCounter uint64
	CurrentMock       mPulseStorageMockCurrent
}

//NewPulseStorageMock returns a mock for github.com/insolar/insolar/core.PulseStorage
func NewPulseStorageMock(t minimock.Tester) *PulseStorageMock {
	m := &PulseStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CurrentMock = mPulseStorageMockCurrent{mock: m}

	return m
}

type mPulseStorageMockCurrent struct {
	mock              *PulseStorageMock
	mainExpectation   *PulseStorageMockCurrentExpectation
	expectationSeries []*PulseStorageMockCurrentExpectation
}

type PulseStorageMockCurrentExpectation struct {
	input  *PulseStorageMockCurrentInput
	result *PulseStorageMockCurrentResult
}

type PulseStorageMockCurrentInput struct {
	p context.Context
}

type PulseStorageMockCurrentResult struct {
	r  *core.Pulse
	r1 error
}

//Expect specifies that invocation of PulseStorage.Current is expected from 1 to Infinity times
func (m *mPulseStorageMockCurrent) Expect(p context.Context) *mPulseStorageMockCurrent {
	m.mock.CurrentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseStorageMockCurrentExpectation{}
	}
	m.mainExpectation.input = &PulseStorageMockCurrentInput{p}
	return m
}

//Return specifies results of invocation of PulseStorage.Current
func (m *mPulseStorageMockCurrent) Return(r *core.Pulse, r1 error) *PulseStorageMock {
	m.mock.CurrentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseStorageMockCurrentExpectation{}
	}
	m.mainExpectation.result = &PulseStorageMockCurrentResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseStorage.Current is expected once
func (m *mPulseStorageMockCurrent) ExpectOnce(p context.Context) *PulseStorageMockCurrentExpectation {
	m.mock.CurrentFunc = nil
	m.mainExpectation = nil

	expectation := &PulseStorageMockCurrentExpectation{}
	expectation.input = &PulseStorageMockCurrentInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseStorageMockCurrentExpectation) Return(r *core.Pulse, r1 error) {
	e.result = &PulseStorageMockCurrentResult{r, r1}
}

//Set uses given function f as a mock of PulseStorage.Current method
func (m *mPulseStorageMockCurrent) Set(f func(p context.Context) (r *core.Pulse, r1 error)) *PulseStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CurrentFunc = f
	return m.mock
}

//Current implements github.com/insolar/insolar/core.PulseStorage interface
func (m *PulseStorageMock) Current(p context.Context) (r *core.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.CurrentPreCounter, 1)
	defer atomic.AddUint64(&m.CurrentCounter, 1)

	if len(m.CurrentMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CurrentMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseStorageMock.Current. %v", p)
			return
		}

		input := m.CurrentMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseStorageMockCurrentInput{p}, "PulseStorage.Current got unexpected parameters")

		result := m.CurrentMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseStorageMock.Current")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CurrentMock.mainExpectation != nil {

		input := m.CurrentMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseStorageMockCurrentInput{p}, "PulseStorage.Current got unexpected parameters")
		}

		result := m.CurrentMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseStorageMock.Current")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CurrentFunc == nil {
		m.t.Fatalf("Unexpected call to PulseStorageMock.Current. %v", p)
		return
	}

	return m.CurrentFunc(p)
}

//CurrentMinimockCounter returns a count of PulseStorageMock.CurrentFunc invocations
func (m *PulseStorageMock) CurrentMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CurrentCounter)
}

//CurrentMinimockPreCounter returns the value of PulseStorageMock.Current invocations
func (m *PulseStorageMock) CurrentMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CurrentPreCounter)
}

//CurrentFinished returns true if mock invocations count is ok
func (m *PulseStorageMock) CurrentFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CurrentMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CurrentCounter) == uint64(len(m.CurrentMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CurrentMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CurrentCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CurrentFunc != nil {
		return atomic.LoadUint64(&m.CurrentCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseStorageMock) ValidateCallCounters() {

	if !m.CurrentFinished() {
		m.t.Fatal("Expected call to PulseStorageMock.Current")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseStorageMock) MinimockFinish() {

	if !m.CurrentFinished() {
		m.t.Fatal("Expected call to PulseStorageMock.Current")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CurrentFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CurrentFinished() {
				m.t.Error("Expected call to PulseStorageMock.Current")
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
func (m *PulseStorageMock) AllMocksCalled() bool {

	if !m.CurrentFinished() {
		return false
	}

	return true
}

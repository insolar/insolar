package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MachinesManager" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//MachinesManagerMock implements github.com/insolar/insolar/logicrunner.MachinesManager
type MachinesManagerMock struct {
	t minimock.Tester

	GetExecutorFunc       func(p insolar.MachineType) (r insolar.MachineLogicExecutor, r1 error)
	GetExecutorCounter    uint64
	GetExecutorPreCounter uint64
	GetExecutorMock       mMachinesManagerMockGetExecutor

	RegisterExecutorFunc       func(p insolar.MachineType, p1 insolar.MachineLogicExecutor) (r error)
	RegisterExecutorCounter    uint64
	RegisterExecutorPreCounter uint64
	RegisterExecutorMock       mMachinesManagerMockRegisterExecutor
}

//NewMachinesManagerMock returns a mock for github.com/insolar/insolar/logicrunner.MachinesManager
func NewMachinesManagerMock(t minimock.Tester) *MachinesManagerMock {
	m := &MachinesManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetExecutorMock = mMachinesManagerMockGetExecutor{mock: m}
	m.RegisterExecutorMock = mMachinesManagerMockRegisterExecutor{mock: m}

	return m
}

type mMachinesManagerMockGetExecutor struct {
	mock              *MachinesManagerMock
	mainExpectation   *MachinesManagerMockGetExecutorExpectation
	expectationSeries []*MachinesManagerMockGetExecutorExpectation
}

type MachinesManagerMockGetExecutorExpectation struct {
	input  *MachinesManagerMockGetExecutorInput
	result *MachinesManagerMockGetExecutorResult
}

type MachinesManagerMockGetExecutorInput struct {
	p insolar.MachineType
}

type MachinesManagerMockGetExecutorResult struct {
	r  insolar.MachineLogicExecutor
	r1 error
}

//Expect specifies that invocation of MachinesManager.GetExecutor is expected from 1 to Infinity times
func (m *mMachinesManagerMockGetExecutor) Expect(p insolar.MachineType) *mMachinesManagerMockGetExecutor {
	m.mock.GetExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachinesManagerMockGetExecutorExpectation{}
	}
	m.mainExpectation.input = &MachinesManagerMockGetExecutorInput{p}
	return m
}

//Return specifies results of invocation of MachinesManager.GetExecutor
func (m *mMachinesManagerMockGetExecutor) Return(r insolar.MachineLogicExecutor, r1 error) *MachinesManagerMock {
	m.mock.GetExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachinesManagerMockGetExecutorExpectation{}
	}
	m.mainExpectation.result = &MachinesManagerMockGetExecutorResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of MachinesManager.GetExecutor is expected once
func (m *mMachinesManagerMockGetExecutor) ExpectOnce(p insolar.MachineType) *MachinesManagerMockGetExecutorExpectation {
	m.mock.GetExecutorFunc = nil
	m.mainExpectation = nil

	expectation := &MachinesManagerMockGetExecutorExpectation{}
	expectation.input = &MachinesManagerMockGetExecutorInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MachinesManagerMockGetExecutorExpectation) Return(r insolar.MachineLogicExecutor, r1 error) {
	e.result = &MachinesManagerMockGetExecutorResult{r, r1}
}

//Set uses given function f as a mock of MachinesManager.GetExecutor method
func (m *mMachinesManagerMockGetExecutor) Set(f func(p insolar.MachineType) (r insolar.MachineLogicExecutor, r1 error)) *MachinesManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExecutorFunc = f
	return m.mock
}

//GetExecutor implements github.com/insolar/insolar/logicrunner.MachinesManager interface
func (m *MachinesManagerMock) GetExecutor(p insolar.MachineType) (r insolar.MachineLogicExecutor, r1 error) {
	counter := atomic.AddUint64(&m.GetExecutorPreCounter, 1)
	defer atomic.AddUint64(&m.GetExecutorCounter, 1)

	if len(m.GetExecutorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExecutorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MachinesManagerMock.GetExecutor. %v", p)
			return
		}

		input := m.GetExecutorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MachinesManagerMockGetExecutorInput{p}, "MachinesManager.GetExecutor got unexpected parameters")

		result := m.GetExecutorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MachinesManagerMock.GetExecutor")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetExecutorMock.mainExpectation != nil {

		input := m.GetExecutorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MachinesManagerMockGetExecutorInput{p}, "MachinesManager.GetExecutor got unexpected parameters")
		}

		result := m.GetExecutorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MachinesManagerMock.GetExecutor")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetExecutorFunc == nil {
		m.t.Fatalf("Unexpected call to MachinesManagerMock.GetExecutor. %v", p)
		return
	}

	return m.GetExecutorFunc(p)
}

//GetExecutorMinimockCounter returns a count of MachinesManagerMock.GetExecutorFunc invocations
func (m *MachinesManagerMock) GetExecutorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutorCounter)
}

//GetExecutorMinimockPreCounter returns the value of MachinesManagerMock.GetExecutor invocations
func (m *MachinesManagerMock) GetExecutorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutorPreCounter)
}

//GetExecutorFinished returns true if mock invocations count is ok
func (m *MachinesManagerMock) GetExecutorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetExecutorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetExecutorCounter) == uint64(len(m.GetExecutorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetExecutorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetExecutorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetExecutorFunc != nil {
		return atomic.LoadUint64(&m.GetExecutorCounter) > 0
	}

	return true
}

type mMachinesManagerMockRegisterExecutor struct {
	mock              *MachinesManagerMock
	mainExpectation   *MachinesManagerMockRegisterExecutorExpectation
	expectationSeries []*MachinesManagerMockRegisterExecutorExpectation
}

type MachinesManagerMockRegisterExecutorExpectation struct {
	input  *MachinesManagerMockRegisterExecutorInput
	result *MachinesManagerMockRegisterExecutorResult
}

type MachinesManagerMockRegisterExecutorInput struct {
	p  insolar.MachineType
	p1 insolar.MachineLogicExecutor
}

type MachinesManagerMockRegisterExecutorResult struct {
	r error
}

//Expect specifies that invocation of MachinesManager.RegisterExecutor is expected from 1 to Infinity times
func (m *mMachinesManagerMockRegisterExecutor) Expect(p insolar.MachineType, p1 insolar.MachineLogicExecutor) *mMachinesManagerMockRegisterExecutor {
	m.mock.RegisterExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachinesManagerMockRegisterExecutorExpectation{}
	}
	m.mainExpectation.input = &MachinesManagerMockRegisterExecutorInput{p, p1}
	return m
}

//Return specifies results of invocation of MachinesManager.RegisterExecutor
func (m *mMachinesManagerMockRegisterExecutor) Return(r error) *MachinesManagerMock {
	m.mock.RegisterExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachinesManagerMockRegisterExecutorExpectation{}
	}
	m.mainExpectation.result = &MachinesManagerMockRegisterExecutorResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MachinesManager.RegisterExecutor is expected once
func (m *mMachinesManagerMockRegisterExecutor) ExpectOnce(p insolar.MachineType, p1 insolar.MachineLogicExecutor) *MachinesManagerMockRegisterExecutorExpectation {
	m.mock.RegisterExecutorFunc = nil
	m.mainExpectation = nil

	expectation := &MachinesManagerMockRegisterExecutorExpectation{}
	expectation.input = &MachinesManagerMockRegisterExecutorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MachinesManagerMockRegisterExecutorExpectation) Return(r error) {
	e.result = &MachinesManagerMockRegisterExecutorResult{r}
}

//Set uses given function f as a mock of MachinesManager.RegisterExecutor method
func (m *mMachinesManagerMockRegisterExecutor) Set(f func(p insolar.MachineType, p1 insolar.MachineLogicExecutor) (r error)) *MachinesManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterExecutorFunc = f
	return m.mock
}

//RegisterExecutor implements github.com/insolar/insolar/logicrunner.MachinesManager interface
func (m *MachinesManagerMock) RegisterExecutor(p insolar.MachineType, p1 insolar.MachineLogicExecutor) (r error) {
	counter := atomic.AddUint64(&m.RegisterExecutorPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterExecutorCounter, 1)

	if len(m.RegisterExecutorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterExecutorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MachinesManagerMock.RegisterExecutor. %v %v", p, p1)
			return
		}

		input := m.RegisterExecutorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MachinesManagerMockRegisterExecutorInput{p, p1}, "MachinesManager.RegisterExecutor got unexpected parameters")

		result := m.RegisterExecutorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MachinesManagerMock.RegisterExecutor")
			return
		}

		r = result.r

		return
	}

	if m.RegisterExecutorMock.mainExpectation != nil {

		input := m.RegisterExecutorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MachinesManagerMockRegisterExecutorInput{p, p1}, "MachinesManager.RegisterExecutor got unexpected parameters")
		}

		result := m.RegisterExecutorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MachinesManagerMock.RegisterExecutor")
		}

		r = result.r

		return
	}

	if m.RegisterExecutorFunc == nil {
		m.t.Fatalf("Unexpected call to MachinesManagerMock.RegisterExecutor. %v %v", p, p1)
		return
	}

	return m.RegisterExecutorFunc(p, p1)
}

//RegisterExecutorMinimockCounter returns a count of MachinesManagerMock.RegisterExecutorFunc invocations
func (m *MachinesManagerMock) RegisterExecutorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterExecutorCounter)
}

//RegisterExecutorMinimockPreCounter returns the value of MachinesManagerMock.RegisterExecutor invocations
func (m *MachinesManagerMock) RegisterExecutorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterExecutorPreCounter)
}

//RegisterExecutorFinished returns true if mock invocations count is ok
func (m *MachinesManagerMock) RegisterExecutorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterExecutorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterExecutorCounter) == uint64(len(m.RegisterExecutorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterExecutorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterExecutorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterExecutorFunc != nil {
		return atomic.LoadUint64(&m.RegisterExecutorCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MachinesManagerMock) ValidateCallCounters() {

	if !m.GetExecutorFinished() {
		m.t.Fatal("Expected call to MachinesManagerMock.GetExecutor")
	}

	if !m.RegisterExecutorFinished() {
		m.t.Fatal("Expected call to MachinesManagerMock.RegisterExecutor")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MachinesManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MachinesManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MachinesManagerMock) MinimockFinish() {

	if !m.GetExecutorFinished() {
		m.t.Fatal("Expected call to MachinesManagerMock.GetExecutor")
	}

	if !m.RegisterExecutorFinished() {
		m.t.Fatal("Expected call to MachinesManagerMock.RegisterExecutor")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MachinesManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MachinesManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetExecutorFinished()
		ok = ok && m.RegisterExecutorFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetExecutorFinished() {
				m.t.Error("Expected call to MachinesManagerMock.GetExecutor")
			}

			if !m.RegisterExecutorFinished() {
				m.t.Error("Expected call to MachinesManagerMock.RegisterExecutor")
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
func (m *MachinesManagerMock) AllMocksCalled() bool {

	if !m.GetExecutorFinished() {
		return false
	}

	if !m.RegisterExecutorFinished() {
		return false
	}

	return true
}

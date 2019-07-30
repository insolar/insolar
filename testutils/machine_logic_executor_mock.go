package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MachineLogicExecutor" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//MachineLogicExecutorMock implements github.com/insolar/insolar/insolar.MachineLogicExecutor
type MachineLogicExecutorMock struct {
	t minimock.Tester

	CallConstructorFunc       func(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 string, p4 insolar.Arguments) (r []byte, r1 string, r2 error)
	CallConstructorCounter    uint64
	CallConstructorPreCounter uint64
	CallConstructorMock       mMachineLogicExecutorMockCallConstructor

	CallMethodFunc       func(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 []byte, p4 string, p5 insolar.Arguments) (r []byte, r1 insolar.Arguments, r2 error)
	CallMethodCounter    uint64
	CallMethodPreCounter uint64
	CallMethodMock       mMachineLogicExecutorMockCallMethod
}

//NewMachineLogicExecutorMock returns a mock for github.com/insolar/insolar/insolar.MachineLogicExecutor
func NewMachineLogicExecutorMock(t minimock.Tester) *MachineLogicExecutorMock {
	m := &MachineLogicExecutorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CallConstructorMock = mMachineLogicExecutorMockCallConstructor{mock: m}
	m.CallMethodMock = mMachineLogicExecutorMockCallMethod{mock: m}

	return m
}

type mMachineLogicExecutorMockCallConstructor struct {
	mock              *MachineLogicExecutorMock
	mainExpectation   *MachineLogicExecutorMockCallConstructorExpectation
	expectationSeries []*MachineLogicExecutorMockCallConstructorExpectation
}

type MachineLogicExecutorMockCallConstructorExpectation struct {
	input  *MachineLogicExecutorMockCallConstructorInput
	result *MachineLogicExecutorMockCallConstructorResult
}

type MachineLogicExecutorMockCallConstructorInput struct {
	p  context.Context
	p1 *insolar.LogicCallContext
	p2 insolar.Reference
	p3 string
	p4 insolar.Arguments
}

type MachineLogicExecutorMockCallConstructorResult struct {
	r  []byte
	r1 string
	r2 error
}

//Expect specifies that invocation of MachineLogicExecutor.CallConstructor is expected from 1 to Infinity times
func (m *mMachineLogicExecutorMockCallConstructor) Expect(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 string, p4 insolar.Arguments) *mMachineLogicExecutorMockCallConstructor {
	m.mock.CallConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachineLogicExecutorMockCallConstructorExpectation{}
	}
	m.mainExpectation.input = &MachineLogicExecutorMockCallConstructorInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of MachineLogicExecutor.CallConstructor
func (m *mMachineLogicExecutorMockCallConstructor) Return(r []byte, r1 string, r2 error) *MachineLogicExecutorMock {
	m.mock.CallConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachineLogicExecutorMockCallConstructorExpectation{}
	}
	m.mainExpectation.result = &MachineLogicExecutorMockCallConstructorResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of MachineLogicExecutor.CallConstructor is expected once
func (m *mMachineLogicExecutorMockCallConstructor) ExpectOnce(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 string, p4 insolar.Arguments) *MachineLogicExecutorMockCallConstructorExpectation {
	m.mock.CallConstructorFunc = nil
	m.mainExpectation = nil

	expectation := &MachineLogicExecutorMockCallConstructorExpectation{}
	expectation.input = &MachineLogicExecutorMockCallConstructorInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MachineLogicExecutorMockCallConstructorExpectation) Return(r []byte, r1 string, r2 error) {
	e.result = &MachineLogicExecutorMockCallConstructorResult{r, r1, r2}
}

//Set uses given function f as a mock of MachineLogicExecutor.CallConstructor method
func (m *mMachineLogicExecutorMockCallConstructor) Set(f func(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 string, p4 insolar.Arguments) (r []byte, r1 string, r2 error)) *MachineLogicExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CallConstructorFunc = f
	return m.mock
}

//CallConstructor implements github.com/insolar/insolar/insolar.MachineLogicExecutor interface
func (m *MachineLogicExecutorMock) CallConstructor(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 string, p4 insolar.Arguments) (r []byte, r1 string, r2 error) {
	counter := atomic.AddUint64(&m.CallConstructorPreCounter, 1)
	defer atomic.AddUint64(&m.CallConstructorCounter, 1)

	if len(m.CallConstructorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CallConstructorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MachineLogicExecutorMock.CallConstructor. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.CallConstructorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MachineLogicExecutorMockCallConstructorInput{p, p1, p2, p3, p4}, "MachineLogicExecutor.CallConstructor got unexpected parameters")

		result := m.CallConstructorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MachineLogicExecutorMock.CallConstructor")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.CallConstructorMock.mainExpectation != nil {

		input := m.CallConstructorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MachineLogicExecutorMockCallConstructorInput{p, p1, p2, p3, p4}, "MachineLogicExecutor.CallConstructor got unexpected parameters")
		}

		result := m.CallConstructorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MachineLogicExecutorMock.CallConstructor")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.CallConstructorFunc == nil {
		m.t.Fatalf("Unexpected call to MachineLogicExecutorMock.CallConstructor. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.CallConstructorFunc(p, p1, p2, p3, p4)
}

//CallConstructorMinimockCounter returns a count of MachineLogicExecutorMock.CallConstructorFunc invocations
func (m *MachineLogicExecutorMock) CallConstructorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CallConstructorCounter)
}

//CallConstructorMinimockPreCounter returns the value of MachineLogicExecutorMock.CallConstructor invocations
func (m *MachineLogicExecutorMock) CallConstructorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CallConstructorPreCounter)
}

//CallConstructorFinished returns true if mock invocations count is ok
func (m *MachineLogicExecutorMock) CallConstructorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CallConstructorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CallConstructorCounter) == uint64(len(m.CallConstructorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CallConstructorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CallConstructorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CallConstructorFunc != nil {
		return atomic.LoadUint64(&m.CallConstructorCounter) > 0
	}

	return true
}

type mMachineLogicExecutorMockCallMethod struct {
	mock              *MachineLogicExecutorMock
	mainExpectation   *MachineLogicExecutorMockCallMethodExpectation
	expectationSeries []*MachineLogicExecutorMockCallMethodExpectation
}

type MachineLogicExecutorMockCallMethodExpectation struct {
	input  *MachineLogicExecutorMockCallMethodInput
	result *MachineLogicExecutorMockCallMethodResult
}

type MachineLogicExecutorMockCallMethodInput struct {
	p  context.Context
	p1 *insolar.LogicCallContext
	p2 insolar.Reference
	p3 []byte
	p4 string
	p5 insolar.Arguments
}

type MachineLogicExecutorMockCallMethodResult struct {
	r  []byte
	r1 insolar.Arguments
	r2 error
}

//Expect specifies that invocation of MachineLogicExecutor.CallMethod is expected from 1 to Infinity times
func (m *mMachineLogicExecutorMockCallMethod) Expect(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 []byte, p4 string, p5 insolar.Arguments) *mMachineLogicExecutorMockCallMethod {
	m.mock.CallMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachineLogicExecutorMockCallMethodExpectation{}
	}
	m.mainExpectation.input = &MachineLogicExecutorMockCallMethodInput{p, p1, p2, p3, p4, p5}
	return m
}

//Return specifies results of invocation of MachineLogicExecutor.CallMethod
func (m *mMachineLogicExecutorMockCallMethod) Return(r []byte, r1 insolar.Arguments, r2 error) *MachineLogicExecutorMock {
	m.mock.CallMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MachineLogicExecutorMockCallMethodExpectation{}
	}
	m.mainExpectation.result = &MachineLogicExecutorMockCallMethodResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of MachineLogicExecutor.CallMethod is expected once
func (m *mMachineLogicExecutorMockCallMethod) ExpectOnce(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 []byte, p4 string, p5 insolar.Arguments) *MachineLogicExecutorMockCallMethodExpectation {
	m.mock.CallMethodFunc = nil
	m.mainExpectation = nil

	expectation := &MachineLogicExecutorMockCallMethodExpectation{}
	expectation.input = &MachineLogicExecutorMockCallMethodInput{p, p1, p2, p3, p4, p5}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MachineLogicExecutorMockCallMethodExpectation) Return(r []byte, r1 insolar.Arguments, r2 error) {
	e.result = &MachineLogicExecutorMockCallMethodResult{r, r1, r2}
}

//Set uses given function f as a mock of MachineLogicExecutor.CallMethod method
func (m *mMachineLogicExecutorMockCallMethod) Set(f func(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 []byte, p4 string, p5 insolar.Arguments) (r []byte, r1 insolar.Arguments, r2 error)) *MachineLogicExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CallMethodFunc = f
	return m.mock
}

//CallMethod implements github.com/insolar/insolar/insolar.MachineLogicExecutor interface
func (m *MachineLogicExecutorMock) CallMethod(p context.Context, p1 *insolar.LogicCallContext, p2 insolar.Reference, p3 []byte, p4 string, p5 insolar.Arguments) (r []byte, r1 insolar.Arguments, r2 error) {
	counter := atomic.AddUint64(&m.CallMethodPreCounter, 1)
	defer atomic.AddUint64(&m.CallMethodCounter, 1)

	if len(m.CallMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CallMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MachineLogicExecutorMock.CallMethod. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
			return
		}

		input := m.CallMethodMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MachineLogicExecutorMockCallMethodInput{p, p1, p2, p3, p4, p5}, "MachineLogicExecutor.CallMethod got unexpected parameters")

		result := m.CallMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MachineLogicExecutorMock.CallMethod")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.CallMethodMock.mainExpectation != nil {

		input := m.CallMethodMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MachineLogicExecutorMockCallMethodInput{p, p1, p2, p3, p4, p5}, "MachineLogicExecutor.CallMethod got unexpected parameters")
		}

		result := m.CallMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MachineLogicExecutorMock.CallMethod")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.CallMethodFunc == nil {
		m.t.Fatalf("Unexpected call to MachineLogicExecutorMock.CallMethod. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
		return
	}

	return m.CallMethodFunc(p, p1, p2, p3, p4, p5)
}

//CallMethodMinimockCounter returns a count of MachineLogicExecutorMock.CallMethodFunc invocations
func (m *MachineLogicExecutorMock) CallMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CallMethodCounter)
}

//CallMethodMinimockPreCounter returns the value of MachineLogicExecutorMock.CallMethod invocations
func (m *MachineLogicExecutorMock) CallMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CallMethodPreCounter)
}

//CallMethodFinished returns true if mock invocations count is ok
func (m *MachineLogicExecutorMock) CallMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CallMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CallMethodCounter) == uint64(len(m.CallMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CallMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CallMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CallMethodFunc != nil {
		return atomic.LoadUint64(&m.CallMethodCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MachineLogicExecutorMock) ValidateCallCounters() {

	if !m.CallConstructorFinished() {
		m.t.Fatal("Expected call to MachineLogicExecutorMock.CallConstructor")
	}

	if !m.CallMethodFinished() {
		m.t.Fatal("Expected call to MachineLogicExecutorMock.CallMethod")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MachineLogicExecutorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MachineLogicExecutorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MachineLogicExecutorMock) MinimockFinish() {

	if !m.CallConstructorFinished() {
		m.t.Fatal("Expected call to MachineLogicExecutorMock.CallConstructor")
	}

	if !m.CallMethodFinished() {
		m.t.Fatal("Expected call to MachineLogicExecutorMock.CallMethod")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MachineLogicExecutorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MachineLogicExecutorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CallConstructorFinished()
		ok = ok && m.CallMethodFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CallConstructorFinished() {
				m.t.Error("Expected call to MachineLogicExecutorMock.CallConstructor")
			}

			if !m.CallMethodFinished() {
				m.t.Error("Expected call to MachineLogicExecutorMock.CallMethod")
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
func (m *MachineLogicExecutorMock) AllMocksCalled() bool {

	if !m.CallConstructorFinished() {
		return false
	}

	if !m.CallMethodFinished() {
		return false
	}

	return true
}

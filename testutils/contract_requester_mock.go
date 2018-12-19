package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ContractRequester" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ContractRequesterMock implements github.com/insolar/insolar/core.ContractRequester
type ContractRequesterMock struct {
	t minimock.Tester

	CallConstructorFunc       func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) (r *core.RecordRef, r1 error)
	CallConstructorCounter    uint64
	CallConstructorPreCounter uint64
	CallConstructorMock       mContractRequesterMockCallConstructor

	CallMethodFunc       func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error)
	CallMethodCounter    uint64
	CallMethodPreCounter uint64
	CallMethodMock       mContractRequesterMockCallMethod

	SendRequestFunc       func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (r core.Reply, r1 error)
	SendRequestCounter    uint64
	SendRequestPreCounter uint64
	SendRequestMock       mContractRequesterMockSendRequest
}

//NewContractRequesterMock returns a mock for github.com/insolar/insolar/core.ContractRequester
func NewContractRequesterMock(t minimock.Tester) *ContractRequesterMock {
	m := &ContractRequesterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CallConstructorMock = mContractRequesterMockCallConstructor{mock: m}
	m.CallMethodMock = mContractRequesterMockCallMethod{mock: m}
	m.SendRequestMock = mContractRequesterMockSendRequest{mock: m}

	return m
}

type mContractRequesterMockCallConstructor struct {
	mock              *ContractRequesterMock
	mainExpectation   *ContractRequesterMockCallConstructorExpectation
	expectationSeries []*ContractRequesterMockCallConstructorExpectation
}

type ContractRequesterMockCallConstructorExpectation struct {
	input  *ContractRequesterMockCallConstructorInput
	result *ContractRequesterMockCallConstructorResult
}

type ContractRequesterMockCallConstructorInput struct {
	p  context.Context
	p1 core.Message
	p2 bool
	p3 *core.RecordRef
	p4 *core.RecordRef
	p5 string
	p6 core.Arguments
	p7 int
}

type ContractRequesterMockCallConstructorResult struct {
	r  *core.RecordRef
	r1 error
}

//Expect specifies that invocation of ContractRequester.CallConstructor is expected from 1 to Infinity times
func (m *mContractRequesterMockCallConstructor) Expect(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) *mContractRequesterMockCallConstructor {
	m.mock.CallConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallConstructorExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockCallConstructorInput{p, p1, p2, p3, p4, p5, p6, p7}
	return m
}

//Return specifies results of invocation of ContractRequester.CallConstructor
func (m *mContractRequesterMockCallConstructor) Return(r *core.RecordRef, r1 error) *ContractRequesterMock {
	m.mock.CallConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallConstructorExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockCallConstructorResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.CallConstructor is expected once
func (m *mContractRequesterMockCallConstructor) ExpectOnce(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) *ContractRequesterMockCallConstructorExpectation {
	m.mock.CallConstructorFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockCallConstructorExpectation{}
	expectation.input = &ContractRequesterMockCallConstructorInput{p, p1, p2, p3, p4, p5, p6, p7}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockCallConstructorExpectation) Return(r *core.RecordRef, r1 error) {
	e.result = &ContractRequesterMockCallConstructorResult{r, r1}
}

//Set uses given function f as a mock of ContractRequester.CallConstructor method
func (m *mContractRequesterMockCallConstructor) Set(f func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) (r *core.RecordRef, r1 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CallConstructorFunc = f
	return m.mock
}

//CallConstructor implements github.com/insolar/insolar/core.ContractRequester interface
func (m *ContractRequesterMock) CallConstructor(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) (r *core.RecordRef, r1 error) {
	counter := atomic.AddUint64(&m.CallConstructorPreCounter, 1)
	defer atomic.AddUint64(&m.CallConstructorCounter, 1)

	if len(m.CallConstructorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CallConstructorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ContractRequesterMock.CallConstructor. %v %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6, p7)
			return
		}

		input := m.CallConstructorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ContractRequesterMockCallConstructorInput{p, p1, p2, p3, p4, p5, p6, p7}, "ContractRequester.CallConstructor got unexpected parameters")

		result := m.CallConstructorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.CallConstructor")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CallConstructorMock.mainExpectation != nil {

		input := m.CallConstructorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ContractRequesterMockCallConstructorInput{p, p1, p2, p3, p4, p5, p6, p7}, "ContractRequester.CallConstructor got unexpected parameters")
		}

		result := m.CallConstructorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.CallConstructor")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CallConstructorFunc == nil {
		m.t.Fatalf("Unexpected call to ContractRequesterMock.CallConstructor. %v %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6, p7)
		return
	}

	return m.CallConstructorFunc(p, p1, p2, p3, p4, p5, p6, p7)
}

//CallConstructorMinimockCounter returns a count of ContractRequesterMock.CallConstructorFunc invocations
func (m *ContractRequesterMock) CallConstructorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CallConstructorCounter)
}

//CallConstructorMinimockPreCounter returns the value of ContractRequesterMock.CallConstructor invocations
func (m *ContractRequesterMock) CallConstructorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CallConstructorPreCounter)
}

//CallConstructorFinished returns true if mock invocations count is ok
func (m *ContractRequesterMock) CallConstructorFinished() bool {
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

type mContractRequesterMockCallMethod struct {
	mock              *ContractRequesterMock
	mainExpectation   *ContractRequesterMockCallMethodExpectation
	expectationSeries []*ContractRequesterMockCallMethodExpectation
}

type ContractRequesterMockCallMethodExpectation struct {
	input  *ContractRequesterMockCallMethodInput
	result *ContractRequesterMockCallMethodResult
}

type ContractRequesterMockCallMethodInput struct {
	p  context.Context
	p1 core.Message
	p2 bool
	p3 *core.RecordRef
	p4 string
	p5 core.Arguments
	p6 *core.RecordRef
}

type ContractRequesterMockCallMethodResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of ContractRequester.CallMethod is expected from 1 to Infinity times
func (m *mContractRequesterMockCallMethod) Expect(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) *mContractRequesterMockCallMethod {
	m.mock.CallMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallMethodExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockCallMethodInput{p, p1, p2, p3, p4, p5, p6}
	return m
}

//Return specifies results of invocation of ContractRequester.CallMethod
func (m *mContractRequesterMockCallMethod) Return(r core.Reply, r1 error) *ContractRequesterMock {
	m.mock.CallMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallMethodExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockCallMethodResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.CallMethod is expected once
func (m *mContractRequesterMockCallMethod) ExpectOnce(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) *ContractRequesterMockCallMethodExpectation {
	m.mock.CallMethodFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockCallMethodExpectation{}
	expectation.input = &ContractRequesterMockCallMethodInput{p, p1, p2, p3, p4, p5, p6}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockCallMethodExpectation) Return(r core.Reply, r1 error) {
	e.result = &ContractRequesterMockCallMethodResult{r, r1}
}

//Set uses given function f as a mock of ContractRequester.CallMethod method
func (m *mContractRequesterMockCallMethod) Set(f func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CallMethodFunc = f
	return m.mock
}

//CallMethod implements github.com/insolar/insolar/core.ContractRequester interface
func (m *ContractRequesterMock) CallMethod(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.CallMethodPreCounter, 1)
	defer atomic.AddUint64(&m.CallMethodCounter, 1)

	if len(m.CallMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CallMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ContractRequesterMock.CallMethod. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
			return
		}

		input := m.CallMethodMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ContractRequesterMockCallMethodInput{p, p1, p2, p3, p4, p5, p6}, "ContractRequester.CallMethod got unexpected parameters")

		result := m.CallMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.CallMethod")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CallMethodMock.mainExpectation != nil {

		input := m.CallMethodMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ContractRequesterMockCallMethodInput{p, p1, p2, p3, p4, p5, p6}, "ContractRequester.CallMethod got unexpected parameters")
		}

		result := m.CallMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.CallMethod")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CallMethodFunc == nil {
		m.t.Fatalf("Unexpected call to ContractRequesterMock.CallMethod. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
		return
	}

	return m.CallMethodFunc(p, p1, p2, p3, p4, p5, p6)
}

//CallMethodMinimockCounter returns a count of ContractRequesterMock.CallMethodFunc invocations
func (m *ContractRequesterMock) CallMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CallMethodCounter)
}

//CallMethodMinimockPreCounter returns the value of ContractRequesterMock.CallMethod invocations
func (m *ContractRequesterMock) CallMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CallMethodPreCounter)
}

//CallMethodFinished returns true if mock invocations count is ok
func (m *ContractRequesterMock) CallMethodFinished() bool {
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

type mContractRequesterMockSendRequest struct {
	mock              *ContractRequesterMock
	mainExpectation   *ContractRequesterMockSendRequestExpectation
	expectationSeries []*ContractRequesterMockSendRequestExpectation
}

type ContractRequesterMockSendRequestExpectation struct {
	input  *ContractRequesterMockSendRequestInput
	result *ContractRequesterMockSendRequestResult
}

type ContractRequesterMockSendRequestInput struct {
	p  context.Context
	p1 *core.RecordRef
	p2 string
	p3 []interface{}
}

type ContractRequesterMockSendRequestResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of ContractRequester.SendRequest is expected from 1 to Infinity times
func (m *mContractRequesterMockSendRequest) Expect(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) *mContractRequesterMockSendRequest {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockSendRequestExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockSendRequestInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ContractRequester.SendRequest
func (m *mContractRequesterMockSendRequest) Return(r core.Reply, r1 error) *ContractRequesterMock {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockSendRequestExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockSendRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.SendRequest is expected once
func (m *mContractRequesterMockSendRequest) ExpectOnce(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) *ContractRequesterMockSendRequestExpectation {
	m.mock.SendRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockSendRequestExpectation{}
	expectation.input = &ContractRequesterMockSendRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockSendRequestExpectation) Return(r core.Reply, r1 error) {
	e.result = &ContractRequesterMockSendRequestResult{r, r1}
}

//Set uses given function f as a mock of ContractRequester.SendRequest method
func (m *mContractRequesterMockSendRequest) Set(f func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (r core.Reply, r1 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRequestFunc = f
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/core.ContractRequester interface
func (m *ContractRequesterMock) SendRequest(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.SendRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestCounter, 1)

	if len(m.SendRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ContractRequesterMock.SendRequest. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ContractRequesterMockSendRequestInput{p, p1, p2, p3}, "ContractRequester.SendRequest got unexpected parameters")

		result := m.SendRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.SendRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestMock.mainExpectation != nil {

		input := m.SendRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ContractRequesterMockSendRequestInput{p, p1, p2, p3}, "ContractRequester.SendRequest got unexpected parameters")
		}

		result := m.SendRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.SendRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ContractRequesterMock.SendRequest. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SendRequestFunc(p, p1, p2, p3)
}

//SendRequestMinimockCounter returns a count of ContractRequesterMock.SendRequestFunc invocations
func (m *ContractRequesterMock) SendRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestCounter)
}

//SendRequestMinimockPreCounter returns the value of ContractRequesterMock.SendRequest invocations
func (m *ContractRequesterMock) SendRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestPreCounter)
}

//SendRequestFinished returns true if mock invocations count is ok
func (m *ContractRequesterMock) SendRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendRequestCounter) == uint64(len(m.SendRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendRequestFunc != nil {
		return atomic.LoadUint64(&m.SendRequestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ContractRequesterMock) ValidateCallCounters() {

	if !m.CallConstructorFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallConstructor")
	}

	if !m.CallMethodFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallMethod")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.SendRequest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ContractRequesterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ContractRequesterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ContractRequesterMock) MinimockFinish() {

	if !m.CallConstructorFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallConstructor")
	}

	if !m.CallMethodFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallMethod")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.SendRequest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ContractRequesterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ContractRequesterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CallConstructorFinished()
		ok = ok && m.CallMethodFinished()
		ok = ok && m.SendRequestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CallConstructorFinished() {
				m.t.Error("Expected call to ContractRequesterMock.CallConstructor")
			}

			if !m.CallMethodFinished() {
				m.t.Error("Expected call to ContractRequesterMock.CallMethod")
			}

			if !m.SendRequestFinished() {
				m.t.Error("Expected call to ContractRequesterMock.SendRequest")
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
func (m *ContractRequesterMock) AllMocksCalled() bool {

	if !m.CallConstructorFinished() {
		return false
	}

	if !m.CallMethodFinished() {
		return false
	}

	if !m.SendRequestFinished() {
		return false
	}

	return true
}

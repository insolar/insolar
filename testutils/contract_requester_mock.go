package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ContractRequester" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ContractRequesterMock implements github.com/insolar/insolar/insolar.ContractRequester
type ContractRequesterMock struct {
	t minimock.Tester

	CallFunc       func(p context.Context, p1 insolar.Message) (r insolar.Reply, r1 error)
	CallCounter    uint64
	CallPreCounter uint64
	CallMock       mContractRequesterMockCall

	CallConstructorFunc       func(p context.Context, p1 insolar.Message) (r *insolar.Reference, r1 string, r2 error)
	CallConstructorCounter    uint64
	CallConstructorPreCounter uint64
	CallConstructorMock       mContractRequesterMockCallConstructor

	CallMethodFunc       func(p context.Context, p1 insolar.Message) (r insolar.Reply, r1 error)
	CallMethodCounter    uint64
	CallMethodPreCounter uint64
	CallMethodMock       mContractRequesterMockCallMethod

	SendRequestFunc       func(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}) (r insolar.Reply, r1 error)
	SendRequestCounter    uint64
	SendRequestPreCounter uint64
	SendRequestMock       mContractRequesterMockSendRequest

	SendRequestWithPulseFunc       func(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}, p4 insolar.PulseNumber) (r insolar.Reply, r1 error)
	SendRequestWithPulseCounter    uint64
	SendRequestWithPulsePreCounter uint64
	SendRequestWithPulseMock       mContractRequesterMockSendRequestWithPulse
}

//NewContractRequesterMock returns a mock for github.com/insolar/insolar/insolar.ContractRequester
func NewContractRequesterMock(t minimock.Tester) *ContractRequesterMock {
	m := &ContractRequesterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CallMock = mContractRequesterMockCall{mock: m}
	m.CallConstructorMock = mContractRequesterMockCallConstructor{mock: m}
	m.CallMethodMock = mContractRequesterMockCallMethod{mock: m}
	m.SendRequestMock = mContractRequesterMockSendRequest{mock: m}
	m.SendRequestWithPulseMock = mContractRequesterMockSendRequestWithPulse{mock: m}

	return m
}

type mContractRequesterMockCall struct {
	mock              *ContractRequesterMock
	mainExpectation   *ContractRequesterMockCallExpectation
	expectationSeries []*ContractRequesterMockCallExpectation
}

type ContractRequesterMockCallExpectation struct {
	input  *ContractRequesterMockCallInput
	result *ContractRequesterMockCallResult
}

type ContractRequesterMockCallInput struct {
	p  context.Context
	p1 insolar.Message
}

type ContractRequesterMockCallResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of ContractRequester.Call is expected from 1 to Infinity times
func (m *mContractRequesterMockCall) Expect(p context.Context, p1 insolar.Message) *mContractRequesterMockCall {
	m.mock.CallFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockCallInput{p, p1}
	return m
}

//Return specifies results of invocation of ContractRequester.Call
func (m *mContractRequesterMockCall) Return(r insolar.Reply, r1 error) *ContractRequesterMock {
	m.mock.CallFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockCallResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.Call is expected once
func (m *mContractRequesterMockCall) ExpectOnce(p context.Context, p1 insolar.Message) *ContractRequesterMockCallExpectation {
	m.mock.CallFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockCallExpectation{}
	expectation.input = &ContractRequesterMockCallInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockCallExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &ContractRequesterMockCallResult{r, r1}
}

//Set uses given function f as a mock of ContractRequester.Call method
func (m *mContractRequesterMockCall) Set(f func(p context.Context, p1 insolar.Message) (r insolar.Reply, r1 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CallFunc = f
	return m.mock
}

//Call implements github.com/insolar/insolar/insolar.ContractRequester interface
func (m *ContractRequesterMock) Call(p context.Context, p1 insolar.Message) (r insolar.Reply, r1 error) {
	counter := atomic.AddUint64(&m.CallPreCounter, 1)
	defer atomic.AddUint64(&m.CallCounter, 1)

	if len(m.CallMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CallMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ContractRequesterMock.Call. %v %v", p, p1)
			return
		}

		input := m.CallMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ContractRequesterMockCallInput{p, p1}, "ContractRequester.Call got unexpected parameters")

		result := m.CallMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.Call")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CallMock.mainExpectation != nil {

		input := m.CallMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ContractRequesterMockCallInput{p, p1}, "ContractRequester.Call got unexpected parameters")
		}

		result := m.CallMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.Call")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CallFunc == nil {
		m.t.Fatalf("Unexpected call to ContractRequesterMock.Call. %v %v", p, p1)
		return
	}

	return m.CallFunc(p, p1)
}

//CallMinimockCounter returns a count of ContractRequesterMock.CallFunc invocations
func (m *ContractRequesterMock) CallMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CallCounter)
}

//CallMinimockPreCounter returns the value of ContractRequesterMock.Call invocations
func (m *ContractRequesterMock) CallMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CallPreCounter)
}

//CallFinished returns true if mock invocations count is ok
func (m *ContractRequesterMock) CallFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CallMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CallCounter) == uint64(len(m.CallMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CallMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CallCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CallFunc != nil {
		return atomic.LoadUint64(&m.CallCounter) > 0
	}

	return true
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
	p1 insolar.Message
}

type ContractRequesterMockCallConstructorResult struct {
	r  *insolar.Reference
	r1 string
	r2 error
}

//Expect specifies that invocation of ContractRequester.CallConstructor is expected from 1 to Infinity times
func (m *mContractRequesterMockCallConstructor) Expect(p context.Context, p1 insolar.Message) *mContractRequesterMockCallConstructor {
	m.mock.CallConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallConstructorExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockCallConstructorInput{p, p1}
	return m
}

//Return specifies results of invocation of ContractRequester.CallConstructor
func (m *mContractRequesterMockCallConstructor) Return(r *insolar.Reference, r1 string, r2 error) *ContractRequesterMock {
	m.mock.CallConstructorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallConstructorExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockCallConstructorResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.CallConstructor is expected once
func (m *mContractRequesterMockCallConstructor) ExpectOnce(p context.Context, p1 insolar.Message) *ContractRequesterMockCallConstructorExpectation {
	m.mock.CallConstructorFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockCallConstructorExpectation{}
	expectation.input = &ContractRequesterMockCallConstructorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockCallConstructorExpectation) Return(r *insolar.Reference, r1 string, r2 error) {
	e.result = &ContractRequesterMockCallConstructorResult{r, r1, r2}
}

//Set uses given function f as a mock of ContractRequester.CallConstructor method
func (m *mContractRequesterMockCallConstructor) Set(f func(p context.Context, p1 insolar.Message) (r *insolar.Reference, r1 string, r2 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CallConstructorFunc = f
	return m.mock
}

//CallConstructor implements github.com/insolar/insolar/insolar.ContractRequester interface
func (m *ContractRequesterMock) CallConstructor(p context.Context, p1 insolar.Message) (r *insolar.Reference, r1 string, r2 error) {
	counter := atomic.AddUint64(&m.CallConstructorPreCounter, 1)
	defer atomic.AddUint64(&m.CallConstructorCounter, 1)

	if len(m.CallConstructorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CallConstructorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ContractRequesterMock.CallConstructor. %v %v", p, p1)
			return
		}

		input := m.CallConstructorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ContractRequesterMockCallConstructorInput{p, p1}, "ContractRequester.CallConstructor got unexpected parameters")

		result := m.CallConstructorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.CallConstructor")
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
			testify_assert.Equal(m.t, *input, ContractRequesterMockCallConstructorInput{p, p1}, "ContractRequester.CallConstructor got unexpected parameters")
		}

		result := m.CallConstructorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.CallConstructor")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.CallConstructorFunc == nil {
		m.t.Fatalf("Unexpected call to ContractRequesterMock.CallConstructor. %v %v", p, p1)
		return
	}

	return m.CallConstructorFunc(p, p1)
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
	p1 insolar.Message
}

type ContractRequesterMockCallMethodResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of ContractRequester.CallMethod is expected from 1 to Infinity times
func (m *mContractRequesterMockCallMethod) Expect(p context.Context, p1 insolar.Message) *mContractRequesterMockCallMethod {
	m.mock.CallMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallMethodExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockCallMethodInput{p, p1}
	return m
}

//Return specifies results of invocation of ContractRequester.CallMethod
func (m *mContractRequesterMockCallMethod) Return(r insolar.Reply, r1 error) *ContractRequesterMock {
	m.mock.CallMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockCallMethodExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockCallMethodResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.CallMethod is expected once
func (m *mContractRequesterMockCallMethod) ExpectOnce(p context.Context, p1 insolar.Message) *ContractRequesterMockCallMethodExpectation {
	m.mock.CallMethodFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockCallMethodExpectation{}
	expectation.input = &ContractRequesterMockCallMethodInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockCallMethodExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &ContractRequesterMockCallMethodResult{r, r1}
}

//Set uses given function f as a mock of ContractRequester.CallMethod method
func (m *mContractRequesterMockCallMethod) Set(f func(p context.Context, p1 insolar.Message) (r insolar.Reply, r1 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CallMethodFunc = f
	return m.mock
}

//CallMethod implements github.com/insolar/insolar/insolar.ContractRequester interface
func (m *ContractRequesterMock) CallMethod(p context.Context, p1 insolar.Message) (r insolar.Reply, r1 error) {
	counter := atomic.AddUint64(&m.CallMethodPreCounter, 1)
	defer atomic.AddUint64(&m.CallMethodCounter, 1)

	if len(m.CallMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CallMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ContractRequesterMock.CallMethod. %v %v", p, p1)
			return
		}

		input := m.CallMethodMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ContractRequesterMockCallMethodInput{p, p1}, "ContractRequester.CallMethod got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, ContractRequesterMockCallMethodInput{p, p1}, "ContractRequester.CallMethod got unexpected parameters")
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
		m.t.Fatalf("Unexpected call to ContractRequesterMock.CallMethod. %v %v", p, p1)
		return
	}

	return m.CallMethodFunc(p, p1)
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
	p1 *insolar.Reference
	p2 string
	p3 []interface{}
}

type ContractRequesterMockSendRequestResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of ContractRequester.SendRequest is expected from 1 to Infinity times
func (m *mContractRequesterMockSendRequest) Expect(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}) *mContractRequesterMockSendRequest {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockSendRequestExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockSendRequestInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ContractRequester.SendRequest
func (m *mContractRequesterMockSendRequest) Return(r insolar.Reply, r1 error) *ContractRequesterMock {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockSendRequestExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockSendRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.SendRequest is expected once
func (m *mContractRequesterMockSendRequest) ExpectOnce(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}) *ContractRequesterMockSendRequestExpectation {
	m.mock.SendRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockSendRequestExpectation{}
	expectation.input = &ContractRequesterMockSendRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockSendRequestExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &ContractRequesterMockSendRequestResult{r, r1}
}

//Set uses given function f as a mock of ContractRequester.SendRequest method
func (m *mContractRequesterMockSendRequest) Set(f func(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}) (r insolar.Reply, r1 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRequestFunc = f
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/insolar.ContractRequester interface
func (m *ContractRequesterMock) SendRequest(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}) (r insolar.Reply, r1 error) {
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

type mContractRequesterMockSendRequestWithPulse struct {
	mock              *ContractRequesterMock
	mainExpectation   *ContractRequesterMockSendRequestWithPulseExpectation
	expectationSeries []*ContractRequesterMockSendRequestWithPulseExpectation
}

type ContractRequesterMockSendRequestWithPulseExpectation struct {
	input  *ContractRequesterMockSendRequestWithPulseInput
	result *ContractRequesterMockSendRequestWithPulseResult
}

type ContractRequesterMockSendRequestWithPulseInput struct {
	p  context.Context
	p1 *insolar.Reference
	p2 string
	p3 []interface{}
	p4 insolar.PulseNumber
}

type ContractRequesterMockSendRequestWithPulseResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of ContractRequester.SendRequestWithPulse is expected from 1 to Infinity times
func (m *mContractRequesterMockSendRequestWithPulse) Expect(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}, p4 insolar.PulseNumber) *mContractRequesterMockSendRequestWithPulse {
	m.mock.SendRequestWithPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockSendRequestWithPulseExpectation{}
	}
	m.mainExpectation.input = &ContractRequesterMockSendRequestWithPulseInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of ContractRequester.SendRequestWithPulse
func (m *mContractRequesterMockSendRequestWithPulse) Return(r insolar.Reply, r1 error) *ContractRequesterMock {
	m.mock.SendRequestWithPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ContractRequesterMockSendRequestWithPulseExpectation{}
	}
	m.mainExpectation.result = &ContractRequesterMockSendRequestWithPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ContractRequester.SendRequestWithPulse is expected once
func (m *mContractRequesterMockSendRequestWithPulse) ExpectOnce(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}, p4 insolar.PulseNumber) *ContractRequesterMockSendRequestWithPulseExpectation {
	m.mock.SendRequestWithPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ContractRequesterMockSendRequestWithPulseExpectation{}
	expectation.input = &ContractRequesterMockSendRequestWithPulseInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ContractRequesterMockSendRequestWithPulseExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &ContractRequesterMockSendRequestWithPulseResult{r, r1}
}

//Set uses given function f as a mock of ContractRequester.SendRequestWithPulse method
func (m *mContractRequesterMockSendRequestWithPulse) Set(f func(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}, p4 insolar.PulseNumber) (r insolar.Reply, r1 error)) *ContractRequesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRequestWithPulseFunc = f
	return m.mock
}

//SendRequestWithPulse implements github.com/insolar/insolar/insolar.ContractRequester interface
func (m *ContractRequesterMock) SendRequestWithPulse(p context.Context, p1 *insolar.Reference, p2 string, p3 []interface{}, p4 insolar.PulseNumber) (r insolar.Reply, r1 error) {
	counter := atomic.AddUint64(&m.SendRequestWithPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestWithPulseCounter, 1)

	if len(m.SendRequestWithPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendRequestWithPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ContractRequesterMock.SendRequestWithPulse. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.SendRequestWithPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ContractRequesterMockSendRequestWithPulseInput{p, p1, p2, p3, p4}, "ContractRequester.SendRequestWithPulse got unexpected parameters")

		result := m.SendRequestWithPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.SendRequestWithPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestWithPulseMock.mainExpectation != nil {

		input := m.SendRequestWithPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ContractRequesterMockSendRequestWithPulseInput{p, p1, p2, p3, p4}, "ContractRequester.SendRequestWithPulse got unexpected parameters")
		}

		result := m.SendRequestWithPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ContractRequesterMock.SendRequestWithPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestWithPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ContractRequesterMock.SendRequestWithPulse. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.SendRequestWithPulseFunc(p, p1, p2, p3, p4)
}

//SendRequestWithPulseMinimockCounter returns a count of ContractRequesterMock.SendRequestWithPulseFunc invocations
func (m *ContractRequesterMock) SendRequestWithPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestWithPulseCounter)
}

//SendRequestWithPulseMinimockPreCounter returns the value of ContractRequesterMock.SendRequestWithPulse invocations
func (m *ContractRequesterMock) SendRequestWithPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestWithPulsePreCounter)
}

//SendRequestWithPulseFinished returns true if mock invocations count is ok
func (m *ContractRequesterMock) SendRequestWithPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendRequestWithPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendRequestWithPulseCounter) == uint64(len(m.SendRequestWithPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendRequestWithPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendRequestWithPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendRequestWithPulseFunc != nil {
		return atomic.LoadUint64(&m.SendRequestWithPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ContractRequesterMock) ValidateCallCounters() {

	if !m.CallFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.Call")
	}

	if !m.CallConstructorFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallConstructor")
	}

	if !m.CallMethodFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallMethod")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.SendRequest")
	}

	if !m.SendRequestWithPulseFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.SendRequestWithPulse")
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

	if !m.CallFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.Call")
	}

	if !m.CallConstructorFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallConstructor")
	}

	if !m.CallMethodFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.CallMethod")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.SendRequest")
	}

	if !m.SendRequestWithPulseFinished() {
		m.t.Fatal("Expected call to ContractRequesterMock.SendRequestWithPulse")
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
		ok = ok && m.CallFinished()
		ok = ok && m.CallConstructorFinished()
		ok = ok && m.CallMethodFinished()
		ok = ok && m.SendRequestFinished()
		ok = ok && m.SendRequestWithPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CallFinished() {
				m.t.Error("Expected call to ContractRequesterMock.Call")
			}

			if !m.CallConstructorFinished() {
				m.t.Error("Expected call to ContractRequesterMock.CallConstructor")
			}

			if !m.CallMethodFinished() {
				m.t.Error("Expected call to ContractRequesterMock.CallMethod")
			}

			if !m.SendRequestFinished() {
				m.t.Error("Expected call to ContractRequesterMock.SendRequest")
			}

			if !m.SendRequestWithPulseFinished() {
				m.t.Error("Expected call to ContractRequesterMock.SendRequestWithPulse")
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

	if !m.CallFinished() {
		return false
	}

	if !m.CallConstructorFinished() {
		return false
	}

	if !m.CallMethodFinished() {
		return false
	}

	if !m.SendRequestFinished() {
		return false
	}

	if !m.SendRequestWithPulseFinished() {
		return false
	}

	return true
}

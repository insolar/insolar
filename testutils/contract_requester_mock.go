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
	mock             *ContractRequesterMock
	mockExpectations *ContractRequesterMockCallConstructorParams
}

//ContractRequesterMockCallConstructorParams represents input parameters of the ContractRequester.CallConstructor
type ContractRequesterMockCallConstructorParams struct {
	p  context.Context
	p1 core.Message
	p2 bool
	p3 *core.RecordRef
	p4 *core.RecordRef
	p5 string
	p6 core.Arguments
	p7 int
}

//Expect sets up expected params for the ContractRequester.CallConstructor
func (m *mContractRequesterMockCallConstructor) Expect(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) *mContractRequesterMockCallConstructor {
	m.mockExpectations = &ContractRequesterMockCallConstructorParams{p, p1, p2, p3, p4, p5, p6, p7}
	return m
}

//Return sets up a mock for ContractRequester.CallConstructor to return Return's arguments
func (m *mContractRequesterMockCallConstructor) Return(r *core.RecordRef, r1 error) *ContractRequesterMock {
	m.mock.CallConstructorFunc = func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ContractRequester.CallConstructor method
func (m *mContractRequesterMockCallConstructor) Set(f func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) (r *core.RecordRef, r1 error)) *ContractRequesterMock {
	m.mock.CallConstructorFunc = f
	m.mockExpectations = nil
	return m.mock
}

//CallConstructor implements github.com/insolar/insolar/core.ContractRequester interface
func (m *ContractRequesterMock) CallConstructor(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 *core.RecordRef, p5 string, p6 core.Arguments, p7 int) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.CallConstructorPreCounter, 1)
	defer atomic.AddUint64(&m.CallConstructorCounter, 1)

	if m.CallConstructorMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.CallConstructorMock.mockExpectations, ContractRequesterMockCallConstructorParams{p, p1, p2, p3, p4, p5, p6, p7},
			"ContractRequester.CallConstructor got unexpected parameters")

		if m.CallConstructorFunc == nil {

			m.t.Fatal("No results are set for the ContractRequesterMock.CallConstructor")

			return
		}
	}

	if m.CallConstructorFunc == nil {
		m.t.Fatal("Unexpected call to ContractRequesterMock.CallConstructor")
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

type mContractRequesterMockCallMethod struct {
	mock             *ContractRequesterMock
	mockExpectations *ContractRequesterMockCallMethodParams
}

//ContractRequesterMockCallMethodParams represents input parameters of the ContractRequester.CallMethod
type ContractRequesterMockCallMethodParams struct {
	p  context.Context
	p1 core.Message
	p2 bool
	p3 *core.RecordRef
	p4 string
	p5 core.Arguments
	p6 *core.RecordRef
}

//Expect sets up expected params for the ContractRequester.CallMethod
func (m *mContractRequesterMockCallMethod) Expect(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) *mContractRequesterMockCallMethod {
	m.mockExpectations = &ContractRequesterMockCallMethodParams{p, p1, p2, p3, p4, p5, p6}
	return m
}

//Return sets up a mock for ContractRequester.CallMethod to return Return's arguments
func (m *mContractRequesterMockCallMethod) Return(r core.Reply, r1 error) *ContractRequesterMock {
	m.mock.CallMethodFunc = func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ContractRequester.CallMethod method
func (m *mContractRequesterMockCallMethod) Set(f func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error)) *ContractRequesterMock {
	m.mock.CallMethodFunc = f
	m.mockExpectations = nil
	return m.mock
}

//CallMethod implements github.com/insolar/insolar/core.ContractRequester interface
func (m *ContractRequesterMock) CallMethod(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.CallMethodPreCounter, 1)
	defer atomic.AddUint64(&m.CallMethodCounter, 1)

	if m.CallMethodMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.CallMethodMock.mockExpectations, ContractRequesterMockCallMethodParams{p, p1, p2, p3, p4, p5, p6},
			"ContractRequester.CallMethod got unexpected parameters")

		if m.CallMethodFunc == nil {

			m.t.Fatal("No results are set for the ContractRequesterMock.CallMethod")

			return
		}
	}

	if m.CallMethodFunc == nil {
		m.t.Fatal("Unexpected call to ContractRequesterMock.CallMethod")
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

type mContractRequesterMockSendRequest struct {
	mock             *ContractRequesterMock
	mockExpectations *ContractRequesterMockSendRequestParams
}

//ContractRequesterMockSendRequestParams represents input parameters of the ContractRequester.SendRequest
type ContractRequesterMockSendRequestParams struct {
	p  context.Context
	p1 *core.RecordRef
	p2 string
	p3 []interface{}
}

//Expect sets up expected params for the ContractRequester.SendRequest
func (m *mContractRequesterMockSendRequest) Expect(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) *mContractRequesterMockSendRequest {
	m.mockExpectations = &ContractRequesterMockSendRequestParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for ContractRequester.SendRequest to return Return's arguments
func (m *mContractRequesterMockSendRequest) Return(r core.Reply, r1 error) *ContractRequesterMock {
	m.mock.SendRequestFunc = func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ContractRequester.SendRequest method
func (m *mContractRequesterMockSendRequest) Set(f func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (r core.Reply, r1 error)) *ContractRequesterMock {
	m.mock.SendRequestFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/core.ContractRequester interface
func (m *ContractRequesterMock) SendRequest(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.SendRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestCounter, 1)

	if m.SendRequestMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendRequestMock.mockExpectations, ContractRequesterMockSendRequestParams{p, p1, p2, p3},
			"ContractRequester.SendRequest got unexpected parameters")

		if m.SendRequestFunc == nil {

			m.t.Fatal("No results are set for the ContractRequesterMock.SendRequest")

			return
		}
	}

	if m.SendRequestFunc == nil {
		m.t.Fatal("Unexpected call to ContractRequesterMock.SendRequest")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ContractRequesterMock) ValidateCallCounters() {

	if m.CallConstructorFunc != nil && atomic.LoadUint64(&m.CallConstructorCounter) == 0 {
		m.t.Fatal("Expected call to ContractRequesterMock.CallConstructor")
	}

	if m.CallMethodFunc != nil && atomic.LoadUint64(&m.CallMethodCounter) == 0 {
		m.t.Fatal("Expected call to ContractRequesterMock.CallMethod")
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
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

	if m.CallConstructorFunc != nil && atomic.LoadUint64(&m.CallConstructorCounter) == 0 {
		m.t.Fatal("Expected call to ContractRequesterMock.CallConstructor")
	}

	if m.CallMethodFunc != nil && atomic.LoadUint64(&m.CallMethodCounter) == 0 {
		m.t.Fatal("Expected call to ContractRequesterMock.CallMethod")
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
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
		ok = ok && (m.CallConstructorFunc == nil || atomic.LoadUint64(&m.CallConstructorCounter) > 0)
		ok = ok && (m.CallMethodFunc == nil || atomic.LoadUint64(&m.CallMethodCounter) > 0)
		ok = ok && (m.SendRequestFunc == nil || atomic.LoadUint64(&m.SendRequestCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CallConstructorFunc != nil && atomic.LoadUint64(&m.CallConstructorCounter) == 0 {
				m.t.Error("Expected call to ContractRequesterMock.CallConstructor")
			}

			if m.CallMethodFunc != nil && atomic.LoadUint64(&m.CallMethodCounter) == 0 {
				m.t.Error("Expected call to ContractRequesterMock.CallMethod")
			}

			if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
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

	if m.CallConstructorFunc != nil && atomic.LoadUint64(&m.CallConstructorCounter) == 0 {
		return false
	}

	if m.CallMethodFunc != nil && atomic.LoadUint64(&m.CallMethodCounter) == 0 {
		return false
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
		return false
	}

	return true
}

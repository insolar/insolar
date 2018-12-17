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

	CallContractFunc       func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error)
	CallContractCounter    uint64
	CallContractPreCounter uint64
	CallContractMock       mContractRequesterMockCallContract

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

	m.CallContractMock = mContractRequesterMockCallContract{mock: m}
	m.SendRequestMock = mContractRequesterMockSendRequest{mock: m}

	return m
}

type mContractRequesterMockCallContract struct {
	mock             *ContractRequesterMock
	mockExpectations *ContractRequesterMockCallContractParams
}

//ContractRequesterMockCallContractParams represents input parameters of the ContractRequester.CallContract
type ContractRequesterMockCallContractParams struct {
	p  context.Context
	p1 core.Message
	p2 bool
	p3 *core.RecordRef
	p4 string
	p5 core.Arguments
	p6 *core.RecordRef
}

//Expect sets up expected params for the ContractRequester.CallContract
func (m *mContractRequesterMockCallContract) Expect(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) *mContractRequesterMockCallContract {
	m.mockExpectations = &ContractRequesterMockCallContractParams{p, p1, p2, p3, p4, p5, p6}
	return m
}

//Return sets up a mock for ContractRequester.CallContract to return Return's arguments
func (m *mContractRequesterMockCallContract) Return(r core.Reply, r1 error) *ContractRequesterMock {
	m.mock.CallContractFunc = func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ContractRequester.CallContract method
func (m *mContractRequesterMockCallContract) Set(f func(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error)) *ContractRequesterMock {
	m.mock.CallContractFunc = f
	m.mockExpectations = nil
	return m.mock
}

//CallContract implements github.com/insolar/insolar/core.ContractRequester interface
func (m *ContractRequesterMock) CallContract(p context.Context, p1 core.Message, p2 bool, p3 *core.RecordRef, p4 string, p5 core.Arguments, p6 *core.RecordRef) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.CallContractPreCounter, 1)
	defer atomic.AddUint64(&m.CallContractCounter, 1)

	if m.CallContractMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.CallContractMock.mockExpectations, ContractRequesterMockCallContractParams{p, p1, p2, p3, p4, p5, p6},
			"ContractRequester.CallContract got unexpected parameters")

		if m.CallContractFunc == nil {

			m.t.Fatal("No results are set for the ContractRequesterMock.CallContract")

			return
		}
	}

	if m.CallContractFunc == nil {
		m.t.Fatal("Unexpected call to ContractRequesterMock.CallContract")
		return
	}

	return m.CallContractFunc(p, p1, p2, p3, p4, p5, p6)
}

//CallContractMinimockCounter returns a count of ContractRequesterMock.CallContractFunc invocations
func (m *ContractRequesterMock) CallContractMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CallContractCounter)
}

//CallContractMinimockPreCounter returns the value of ContractRequesterMock.CallContract invocations
func (m *ContractRequesterMock) CallContractMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CallContractPreCounter)
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

	if m.CallContractFunc != nil && atomic.LoadUint64(&m.CallContractCounter) == 0 {
		m.t.Fatal("Expected call to ContractRequesterMock.CallContract")
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

	if m.CallContractFunc != nil && atomic.LoadUint64(&m.CallContractCounter) == 0 {
		m.t.Fatal("Expected call to ContractRequesterMock.CallContract")
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
		ok = ok && (m.CallContractFunc == nil || atomic.LoadUint64(&m.CallContractCounter) > 0)
		ok = ok && (m.SendRequestFunc == nil || atomic.LoadUint64(&m.SendRequestCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CallContractFunc != nil && atomic.LoadUint64(&m.CallContractCounter) == 0 {
				m.t.Error("Expected call to ContractRequesterMock.CallContract")
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

	if m.CallContractFunc != nil && atomic.LoadUint64(&m.CallContractCounter) == 0 {
		return false
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
		return false
	}

	return true
}

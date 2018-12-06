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

	m.SendRequestMock = mContractRequesterMockSendRequest{mock: m}

	return m
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
		ok = ok && m.SendRequestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

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

	if !m.SendRequestFinished() {
		return false
	}

	return true
}

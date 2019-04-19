package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Network" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	message "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//NetworkMock implements github.com/insolar/insolar/insolar.Network
type NetworkMock struct {
	t minimock.Tester

	GetStateFunc       func() (r insolar.NetworkState)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mNetworkMockGetState

	LeaveFunc       func(p context.Context, p1 insolar.PulseNumber)
	LeaveCounter    uint64
	LeavePreCounter uint64
	LeaveMock       mNetworkMockLeave

	ProcessOutcomeFunc       func(p *message.Message) (r []*message.Message, r1 error)
	ProcessOutcomeCounter    uint64
	ProcessOutcomePreCounter uint64
	ProcessOutcomeMock       mNetworkMockProcessOutcome

	RemoteProcedureRegisterFunc       func(p string, p1 insolar.RemoteProcedure)
	RemoteProcedureRegisterCounter    uint64
	RemoteProcedureRegisterPreCounter uint64
	RemoteProcedureRegisterMock       mNetworkMockRemoteProcedureRegister

	SendCascadeMessageFunc       func(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error)
	SendCascadeMessageCounter    uint64
	SendCascadeMessagePreCounter uint64
	SendCascadeMessageMock       mNetworkMockSendCascadeMessage

	SendMessageFunc       func(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error)
	SendMessageCounter    uint64
	SendMessagePreCounter uint64
	SendMessageMock       mNetworkMockSendMessage
}

//NewNetworkMock returns a mock for github.com/insolar/insolar/insolar.Network
func NewNetworkMock(t minimock.Tester) *NetworkMock {
	m := &NetworkMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetStateMock = mNetworkMockGetState{mock: m}
	m.LeaveMock = mNetworkMockLeave{mock: m}
	m.ProcessOutcomeMock = mNetworkMockProcessOutcome{mock: m}
	m.RemoteProcedureRegisterMock = mNetworkMockRemoteProcedureRegister{mock: m}
	m.SendCascadeMessageMock = mNetworkMockSendCascadeMessage{mock: m}
	m.SendMessageMock = mNetworkMockSendMessage{mock: m}

	return m
}

type mNetworkMockGetState struct {
	mock              *NetworkMock
	mainExpectation   *NetworkMockGetStateExpectation
	expectationSeries []*NetworkMockGetStateExpectation
}

type NetworkMockGetStateExpectation struct {
	result *NetworkMockGetStateResult
}

type NetworkMockGetStateResult struct {
	r insolar.NetworkState
}

//Expect specifies that invocation of Network.GetState is expected from 1 to Infinity times
func (m *mNetworkMockGetState) Expect() *mNetworkMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of Network.GetState
func (m *mNetworkMockGetState) Return(r insolar.NetworkState) *NetworkMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockGetStateExpectation{}
	}
	m.mainExpectation.result = &NetworkMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Network.GetState is expected once
func (m *mNetworkMockGetState) ExpectOnce() *NetworkMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkMockGetStateExpectation) Return(r insolar.NetworkState) {
	e.result = &NetworkMockGetStateResult{r}
}

//Set uses given function f as a mock of Network.GetState method
func (m *mNetworkMockGetState) Set(f func() (r insolar.NetworkState)) *NetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/insolar.Network interface
func (m *NetworkMock) GetState() (r insolar.NetworkState) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of NetworkMock.GetStateFunc invocations
func (m *NetworkMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of NetworkMock.GetState invocations
func (m *NetworkMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *NetworkMock) GetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStateCounter) == uint64(len(m.GetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStateFunc != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	return true
}

type mNetworkMockLeave struct {
	mock              *NetworkMock
	mainExpectation   *NetworkMockLeaveExpectation
	expectationSeries []*NetworkMockLeaveExpectation
}

type NetworkMockLeaveExpectation struct {
	input *NetworkMockLeaveInput
}

type NetworkMockLeaveInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of Network.Leave is expected from 1 to Infinity times
func (m *mNetworkMockLeave) Expect(p context.Context, p1 insolar.PulseNumber) *mNetworkMockLeave {
	m.mock.LeaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockLeaveExpectation{}
	}
	m.mainExpectation.input = &NetworkMockLeaveInput{p, p1}
	return m
}

//Return specifies results of invocation of Network.Leave
func (m *mNetworkMockLeave) Return() *NetworkMock {
	m.mock.LeaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockLeaveExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Network.Leave is expected once
func (m *mNetworkMockLeave) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *NetworkMockLeaveExpectation {
	m.mock.LeaveFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkMockLeaveExpectation{}
	expectation.input = &NetworkMockLeaveInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Network.Leave method
func (m *mNetworkMockLeave) Set(f func(p context.Context, p1 insolar.PulseNumber)) *NetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LeaveFunc = f
	return m.mock
}

//Leave implements github.com/insolar/insolar/insolar.Network interface
func (m *NetworkMock) Leave(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.LeavePreCounter, 1)
	defer atomic.AddUint64(&m.LeaveCounter, 1)

	if len(m.LeaveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LeaveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkMock.Leave. %v %v", p, p1)
			return
		}

		input := m.LeaveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkMockLeaveInput{p, p1}, "Network.Leave got unexpected parameters")

		return
	}

	if m.LeaveMock.mainExpectation != nil {

		input := m.LeaveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkMockLeaveInput{p, p1}, "Network.Leave got unexpected parameters")
		}

		return
	}

	if m.LeaveFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkMock.Leave. %v %v", p, p1)
		return
	}

	m.LeaveFunc(p, p1)
}

//LeaveMinimockCounter returns a count of NetworkMock.LeaveFunc invocations
func (m *NetworkMock) LeaveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LeaveCounter)
}

//LeaveMinimockPreCounter returns the value of NetworkMock.Leave invocations
func (m *NetworkMock) LeaveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LeavePreCounter)
}

//LeaveFinished returns true if mock invocations count is ok
func (m *NetworkMock) LeaveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LeaveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LeaveCounter) == uint64(len(m.LeaveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LeaveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LeaveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LeaveFunc != nil {
		return atomic.LoadUint64(&m.LeaveCounter) > 0
	}

	return true
}

type mNetworkMockProcessOutcome struct {
	mock              *NetworkMock
	mainExpectation   *NetworkMockProcessOutcomeExpectation
	expectationSeries []*NetworkMockProcessOutcomeExpectation
}

type NetworkMockProcessOutcomeExpectation struct {
	input  *NetworkMockProcessOutcomeInput
	result *NetworkMockProcessOutcomeResult
}

type NetworkMockProcessOutcomeInput struct {
	p *message.Message
}

type NetworkMockProcessOutcomeResult struct {
	r  []*message.Message
	r1 error
}

//Expect specifies that invocation of Network.ProcessOutcome is expected from 1 to Infinity times
func (m *mNetworkMockProcessOutcome) Expect(p *message.Message) *mNetworkMockProcessOutcome {
	m.mock.ProcessOutcomeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockProcessOutcomeExpectation{}
	}
	m.mainExpectation.input = &NetworkMockProcessOutcomeInput{p}
	return m
}

//Return specifies results of invocation of Network.ProcessOutcome
func (m *mNetworkMockProcessOutcome) Return(r []*message.Message, r1 error) *NetworkMock {
	m.mock.ProcessOutcomeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockProcessOutcomeExpectation{}
	}
	m.mainExpectation.result = &NetworkMockProcessOutcomeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Network.ProcessOutcome is expected once
func (m *mNetworkMockProcessOutcome) ExpectOnce(p *message.Message) *NetworkMockProcessOutcomeExpectation {
	m.mock.ProcessOutcomeFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkMockProcessOutcomeExpectation{}
	expectation.input = &NetworkMockProcessOutcomeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkMockProcessOutcomeExpectation) Return(r []*message.Message, r1 error) {
	e.result = &NetworkMockProcessOutcomeResult{r, r1}
}

//Set uses given function f as a mock of Network.ProcessOutcome method
func (m *mNetworkMockProcessOutcome) Set(f func(p *message.Message) (r []*message.Message, r1 error)) *NetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ProcessOutcomeFunc = f
	return m.mock
}

//ProcessOutcome implements github.com/insolar/insolar/insolar.Network interface
func (m *NetworkMock) ProcessOutcome(p *message.Message) (r []*message.Message, r1 error) {
	counter := atomic.AddUint64(&m.ProcessOutcomePreCounter, 1)
	defer atomic.AddUint64(&m.ProcessOutcomeCounter, 1)

	if len(m.ProcessOutcomeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ProcessOutcomeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkMock.ProcessOutcome. %v", p)
			return
		}

		input := m.ProcessOutcomeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkMockProcessOutcomeInput{p}, "Network.ProcessOutcome got unexpected parameters")

		result := m.ProcessOutcomeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.ProcessOutcome")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ProcessOutcomeMock.mainExpectation != nil {

		input := m.ProcessOutcomeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkMockProcessOutcomeInput{p}, "Network.ProcessOutcome got unexpected parameters")
		}

		result := m.ProcessOutcomeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.ProcessOutcome")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ProcessOutcomeFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkMock.ProcessOutcome. %v", p)
		return
	}

	return m.ProcessOutcomeFunc(p)
}

//ProcessOutcomeMinimockCounter returns a count of NetworkMock.ProcessOutcomeFunc invocations
func (m *NetworkMock) ProcessOutcomeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ProcessOutcomeCounter)
}

//ProcessOutcomeMinimockPreCounter returns the value of NetworkMock.ProcessOutcome invocations
func (m *NetworkMock) ProcessOutcomeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ProcessOutcomePreCounter)
}

//ProcessOutcomeFinished returns true if mock invocations count is ok
func (m *NetworkMock) ProcessOutcomeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ProcessOutcomeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ProcessOutcomeCounter) == uint64(len(m.ProcessOutcomeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ProcessOutcomeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ProcessOutcomeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ProcessOutcomeFunc != nil {
		return atomic.LoadUint64(&m.ProcessOutcomeCounter) > 0
	}

	return true
}

type mNetworkMockRemoteProcedureRegister struct {
	mock              *NetworkMock
	mainExpectation   *NetworkMockRemoteProcedureRegisterExpectation
	expectationSeries []*NetworkMockRemoteProcedureRegisterExpectation
}

type NetworkMockRemoteProcedureRegisterExpectation struct {
	input *NetworkMockRemoteProcedureRegisterInput
}

type NetworkMockRemoteProcedureRegisterInput struct {
	p  string
	p1 insolar.RemoteProcedure
}

//Expect specifies that invocation of Network.RemoteProcedureRegister is expected from 1 to Infinity times
func (m *mNetworkMockRemoteProcedureRegister) Expect(p string, p1 insolar.RemoteProcedure) *mNetworkMockRemoteProcedureRegister {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockRemoteProcedureRegisterExpectation{}
	}
	m.mainExpectation.input = &NetworkMockRemoteProcedureRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of Network.RemoteProcedureRegister
func (m *mNetworkMockRemoteProcedureRegister) Return() *NetworkMock {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockRemoteProcedureRegisterExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Network.RemoteProcedureRegister is expected once
func (m *mNetworkMockRemoteProcedureRegister) ExpectOnce(p string, p1 insolar.RemoteProcedure) *NetworkMockRemoteProcedureRegisterExpectation {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkMockRemoteProcedureRegisterExpectation{}
	expectation.input = &NetworkMockRemoteProcedureRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Network.RemoteProcedureRegister method
func (m *mNetworkMockRemoteProcedureRegister) Set(f func(p string, p1 insolar.RemoteProcedure)) *NetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoteProcedureRegisterFunc = f
	return m.mock
}

//RemoteProcedureRegister implements github.com/insolar/insolar/insolar.Network interface
func (m *NetworkMock) RemoteProcedureRegister(p string, p1 insolar.RemoteProcedure) {
	counter := atomic.AddUint64(&m.RemoteProcedureRegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RemoteProcedureRegisterCounter, 1)

	if len(m.RemoteProcedureRegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoteProcedureRegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkMock.RemoteProcedureRegister. %v %v", p, p1)
			return
		}

		input := m.RemoteProcedureRegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkMockRemoteProcedureRegisterInput{p, p1}, "Network.RemoteProcedureRegister got unexpected parameters")

		return
	}

	if m.RemoteProcedureRegisterMock.mainExpectation != nil {

		input := m.RemoteProcedureRegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkMockRemoteProcedureRegisterInput{p, p1}, "Network.RemoteProcedureRegister got unexpected parameters")
		}

		return
	}

	if m.RemoteProcedureRegisterFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkMock.RemoteProcedureRegister. %v %v", p, p1)
		return
	}

	m.RemoteProcedureRegisterFunc(p, p1)
}

//RemoteProcedureRegisterMinimockCounter returns a count of NetworkMock.RemoteProcedureRegisterFunc invocations
func (m *NetworkMock) RemoteProcedureRegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter)
}

//RemoteProcedureRegisterMinimockPreCounter returns the value of NetworkMock.RemoteProcedureRegister invocations
func (m *NetworkMock) RemoteProcedureRegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoteProcedureRegisterPreCounter)
}

//RemoteProcedureRegisterFinished returns true if mock invocations count is ok
func (m *NetworkMock) RemoteProcedureRegisterFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoteProcedureRegisterMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter) == uint64(len(m.RemoteProcedureRegisterMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoteProcedureRegisterMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoteProcedureRegisterFunc != nil {
		return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter) > 0
	}

	return true
}

type mNetworkMockSendCascadeMessage struct {
	mock              *NetworkMock
	mainExpectation   *NetworkMockSendCascadeMessageExpectation
	expectationSeries []*NetworkMockSendCascadeMessageExpectation
}

type NetworkMockSendCascadeMessageExpectation struct {
	input  *NetworkMockSendCascadeMessageInput
	result *NetworkMockSendCascadeMessageResult
}

type NetworkMockSendCascadeMessageInput struct {
	p  insolar.Cascade
	p1 string
	p2 insolar.Parcel
}

type NetworkMockSendCascadeMessageResult struct {
	r error
}

//Expect specifies that invocation of Network.SendCascadeMessage is expected from 1 to Infinity times
func (m *mNetworkMockSendCascadeMessage) Expect(p insolar.Cascade, p1 string, p2 insolar.Parcel) *mNetworkMockSendCascadeMessage {
	m.mock.SendCascadeMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockSendCascadeMessageExpectation{}
	}
	m.mainExpectation.input = &NetworkMockSendCascadeMessageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Network.SendCascadeMessage
func (m *mNetworkMockSendCascadeMessage) Return(r error) *NetworkMock {
	m.mock.SendCascadeMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockSendCascadeMessageExpectation{}
	}
	m.mainExpectation.result = &NetworkMockSendCascadeMessageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Network.SendCascadeMessage is expected once
func (m *mNetworkMockSendCascadeMessage) ExpectOnce(p insolar.Cascade, p1 string, p2 insolar.Parcel) *NetworkMockSendCascadeMessageExpectation {
	m.mock.SendCascadeMessageFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkMockSendCascadeMessageExpectation{}
	expectation.input = &NetworkMockSendCascadeMessageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkMockSendCascadeMessageExpectation) Return(r error) {
	e.result = &NetworkMockSendCascadeMessageResult{r}
}

//Set uses given function f as a mock of Network.SendCascadeMessage method
func (m *mNetworkMockSendCascadeMessage) Set(f func(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error)) *NetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendCascadeMessageFunc = f
	return m.mock
}

//SendCascadeMessage implements github.com/insolar/insolar/insolar.Network interface
func (m *NetworkMock) SendCascadeMessage(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error) {
	counter := atomic.AddUint64(&m.SendCascadeMessagePreCounter, 1)
	defer atomic.AddUint64(&m.SendCascadeMessageCounter, 1)

	if len(m.SendCascadeMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendCascadeMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkMock.SendCascadeMessage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendCascadeMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkMockSendCascadeMessageInput{p, p1, p2}, "Network.SendCascadeMessage got unexpected parameters")

		result := m.SendCascadeMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.SendCascadeMessage")
			return
		}

		r = result.r

		return
	}

	if m.SendCascadeMessageMock.mainExpectation != nil {

		input := m.SendCascadeMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkMockSendCascadeMessageInput{p, p1, p2}, "Network.SendCascadeMessage got unexpected parameters")
		}

		result := m.SendCascadeMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.SendCascadeMessage")
		}

		r = result.r

		return
	}

	if m.SendCascadeMessageFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkMock.SendCascadeMessage. %v %v %v", p, p1, p2)
		return
	}

	return m.SendCascadeMessageFunc(p, p1, p2)
}

//SendCascadeMessageMinimockCounter returns a count of NetworkMock.SendCascadeMessageFunc invocations
func (m *NetworkMock) SendCascadeMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCascadeMessageCounter)
}

//SendCascadeMessageMinimockPreCounter returns the value of NetworkMock.SendCascadeMessage invocations
func (m *NetworkMock) SendCascadeMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendCascadeMessagePreCounter)
}

//SendCascadeMessageFinished returns true if mock invocations count is ok
func (m *NetworkMock) SendCascadeMessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendCascadeMessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendCascadeMessageCounter) == uint64(len(m.SendCascadeMessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendCascadeMessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendCascadeMessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendCascadeMessageFunc != nil {
		return atomic.LoadUint64(&m.SendCascadeMessageCounter) > 0
	}

	return true
}

type mNetworkMockSendMessage struct {
	mock              *NetworkMock
	mainExpectation   *NetworkMockSendMessageExpectation
	expectationSeries []*NetworkMockSendMessageExpectation
}

type NetworkMockSendMessageExpectation struct {
	input  *NetworkMockSendMessageInput
	result *NetworkMockSendMessageResult
}

type NetworkMockSendMessageInput struct {
	p  insolar.Reference
	p1 string
	p2 insolar.Parcel
}

type NetworkMockSendMessageResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of Network.SendMessage is expected from 1 to Infinity times
func (m *mNetworkMockSendMessage) Expect(p insolar.Reference, p1 string, p2 insolar.Parcel) *mNetworkMockSendMessage {
	m.mock.SendMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockSendMessageExpectation{}
	}
	m.mainExpectation.input = &NetworkMockSendMessageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Network.SendMessage
func (m *mNetworkMockSendMessage) Return(r []byte, r1 error) *NetworkMock {
	m.mock.SendMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkMockSendMessageExpectation{}
	}
	m.mainExpectation.result = &NetworkMockSendMessageResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Network.SendMessage is expected once
func (m *mNetworkMockSendMessage) ExpectOnce(p insolar.Reference, p1 string, p2 insolar.Parcel) *NetworkMockSendMessageExpectation {
	m.mock.SendMessageFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkMockSendMessageExpectation{}
	expectation.input = &NetworkMockSendMessageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkMockSendMessageExpectation) Return(r []byte, r1 error) {
	e.result = &NetworkMockSendMessageResult{r, r1}
}

//Set uses given function f as a mock of Network.SendMessage method
func (m *mNetworkMockSendMessage) Set(f func(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error)) *NetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendMessageFunc = f
	return m.mock
}

//SendMessage implements github.com/insolar/insolar/insolar.Network interface
func (m *NetworkMock) SendMessage(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.SendMessagePreCounter, 1)
	defer atomic.AddUint64(&m.SendMessageCounter, 1)

	if len(m.SendMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkMock.SendMessage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkMockSendMessageInput{p, p1, p2}, "Network.SendMessage got unexpected parameters")

		result := m.SendMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.SendMessage")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMessageMock.mainExpectation != nil {

		input := m.SendMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkMockSendMessageInput{p, p1, p2}, "Network.SendMessage got unexpected parameters")
		}

		result := m.SendMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkMock.SendMessage")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMessageFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkMock.SendMessage. %v %v %v", p, p1, p2)
		return
	}

	return m.SendMessageFunc(p, p1, p2)
}

//SendMessageMinimockCounter returns a count of NetworkMock.SendMessageFunc invocations
func (m *NetworkMock) SendMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessageCounter)
}

//SendMessageMinimockPreCounter returns the value of NetworkMock.SendMessage invocations
func (m *NetworkMock) SendMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessagePreCounter)
}

//SendMessageFinished returns true if mock invocations count is ok
func (m *NetworkMock) SendMessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendMessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendMessageCounter) == uint64(len(m.SendMessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendMessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendMessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendMessageFunc != nil {
		return atomic.LoadUint64(&m.SendMessageCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkMock) ValidateCallCounters() {

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NetworkMock.GetState")
	}

	if !m.LeaveFinished() {
		m.t.Fatal("Expected call to NetworkMock.Leave")
	}

	if !m.ProcessOutcomeFinished() {
		m.t.Fatal("Expected call to NetworkMock.ProcessOutcome")
	}

	if !m.RemoteProcedureRegisterFinished() {
		m.t.Fatal("Expected call to NetworkMock.RemoteProcedureRegister")
	}

	if !m.SendCascadeMessageFinished() {
		m.t.Fatal("Expected call to NetworkMock.SendCascadeMessage")
	}

	if !m.SendMessageFinished() {
		m.t.Fatal("Expected call to NetworkMock.SendMessage")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NetworkMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NetworkMock) MinimockFinish() {

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NetworkMock.GetState")
	}

	if !m.LeaveFinished() {
		m.t.Fatal("Expected call to NetworkMock.Leave")
	}

	if !m.ProcessOutcomeFinished() {
		m.t.Fatal("Expected call to NetworkMock.ProcessOutcome")
	}

	if !m.RemoteProcedureRegisterFinished() {
		m.t.Fatal("Expected call to NetworkMock.RemoteProcedureRegister")
	}

	if !m.SendCascadeMessageFinished() {
		m.t.Fatal("Expected call to NetworkMock.SendCascadeMessage")
	}

	if !m.SendMessageFinished() {
		m.t.Fatal("Expected call to NetworkMock.SendMessage")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NetworkMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NetworkMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetStateFinished()
		ok = ok && m.LeaveFinished()
		ok = ok && m.ProcessOutcomeFinished()
		ok = ok && m.RemoteProcedureRegisterFinished()
		ok = ok && m.SendCascadeMessageFinished()
		ok = ok && m.SendMessageFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetStateFinished() {
				m.t.Error("Expected call to NetworkMock.GetState")
			}

			if !m.LeaveFinished() {
				m.t.Error("Expected call to NetworkMock.Leave")
			}

			if !m.ProcessOutcomeFinished() {
				m.t.Error("Expected call to NetworkMock.ProcessOutcome")
			}

			if !m.RemoteProcedureRegisterFinished() {
				m.t.Error("Expected call to NetworkMock.RemoteProcedureRegister")
			}

			if !m.SendCascadeMessageFinished() {
				m.t.Error("Expected call to NetworkMock.SendCascadeMessage")
			}

			if !m.SendMessageFinished() {
				m.t.Error("Expected call to NetworkMock.SendMessage")
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
func (m *NetworkMock) AllMocksCalled() bool {

	if !m.GetStateFinished() {
		return false
	}

	if !m.LeaveFinished() {
		return false
	}

	if !m.ProcessOutcomeFinished() {
		return false
	}

	if !m.RemoteProcedureRegisterFinished() {
		return false
	}

	if !m.SendCascadeMessageFinished() {
		return false
	}

	if !m.SendMessageFinished() {
		return false
	}

	return true
}

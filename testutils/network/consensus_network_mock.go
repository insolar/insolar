package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ConsensusNetwork" can be found in github.com/insolar/insolar/network
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	network "github.com/insolar/insolar/network"
	types "github.com/insolar/insolar/network/transport/packet/types"

	testify_assert "github.com/stretchr/testify/assert"
)

//ConsensusNetworkMock implements github.com/insolar/insolar/network.ConsensusNetwork
type ConsensusNetworkMock struct {
	t minimock.Tester

	GetNodeIDFunc       func() (r core.RecordRef)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mConsensusNetworkMockGetNodeID

	NewRequestBuilderFunc       func() (r network.RequestBuilder)
	NewRequestBuilderCounter    uint64
	NewRequestBuilderPreCounter uint64
	NewRequestBuilderMock       mConsensusNetworkMockNewRequestBuilder

	PublicAddressFunc       func() (r string)
	PublicAddressCounter    uint64
	PublicAddressPreCounter uint64
	PublicAddressMock       mConsensusNetworkMockPublicAddress

	RegisterRequestHandlerFunc       func(p types.PacketType, p1 network.ConsensusRequestHandler)
	RegisterRequestHandlerCounter    uint64
	RegisterRequestHandlerPreCounter uint64
	RegisterRequestHandlerMock       mConsensusNetworkMockRegisterRequestHandler

	SendRequestFunc       func(p network.Request, p1 core.RecordRef) (r error)
	SendRequestCounter    uint64
	SendRequestPreCounter uint64
	SendRequestMock       mConsensusNetworkMockSendRequest

	StartFunc       func(p context.Context)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mConsensusNetworkMockStart

	StopFunc       func()
	StopCounter    uint64
	StopPreCounter uint64
	StopMock       mConsensusNetworkMockStop
}

//NewConsensusNetworkMock returns a mock for github.com/insolar/insolar/network.ConsensusNetwork
func NewConsensusNetworkMock(t minimock.Tester) *ConsensusNetworkMock {
	m := &ConsensusNetworkMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetNodeIDMock = mConsensusNetworkMockGetNodeID{mock: m}
	m.NewRequestBuilderMock = mConsensusNetworkMockNewRequestBuilder{mock: m}
	m.PublicAddressMock = mConsensusNetworkMockPublicAddress{mock: m}
	m.RegisterRequestHandlerMock = mConsensusNetworkMockRegisterRequestHandler{mock: m}
	m.SendRequestMock = mConsensusNetworkMockSendRequest{mock: m}
	m.StartMock = mConsensusNetworkMockStart{mock: m}
	m.StopMock = mConsensusNetworkMockStop{mock: m}

	return m
}

type mConsensusNetworkMockGetNodeID struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockGetNodeIDExpectation
	expectationSeries []*ConsensusNetworkMockGetNodeIDExpectation
}

type ConsensusNetworkMockGetNodeIDExpectation struct {
	result *ConsensusNetworkMockGetNodeIDResult
}

type ConsensusNetworkMockGetNodeIDResult struct {
	r core.RecordRef
}

//Expect specifies that invocation of ConsensusNetwork.GetNodeID is expected from 1 to Infinity times
func (m *mConsensusNetworkMockGetNodeID) Expect() *mConsensusNetworkMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of ConsensusNetwork.GetNodeID
func (m *mConsensusNetworkMockGetNodeID) Return(r core.RecordRef) *ConsensusNetworkMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &ConsensusNetworkMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.GetNodeID is expected once
func (m *mConsensusNetworkMockGetNodeID) ExpectOnce() *ConsensusNetworkMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConsensusNetworkMockGetNodeIDExpectation) Return(r core.RecordRef) {
	e.result = &ConsensusNetworkMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of ConsensusNetwork.GetNodeID method
func (m *mConsensusNetworkMockGetNodeID) Set(f func() (r core.RecordRef)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) GetNodeID() (r core.RecordRef) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of ConsensusNetworkMock.GetNodeIDFunc invocations
func (m *ConsensusNetworkMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of ConsensusNetworkMock.GetNodeID invocations
func (m *ConsensusNetworkMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) GetNodeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeIDCounter) == uint64(len(m.GetNodeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeIDFunc != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	return true
}

type mConsensusNetworkMockNewRequestBuilder struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockNewRequestBuilderExpectation
	expectationSeries []*ConsensusNetworkMockNewRequestBuilderExpectation
}

type ConsensusNetworkMockNewRequestBuilderExpectation struct {
	result *ConsensusNetworkMockNewRequestBuilderResult
}

type ConsensusNetworkMockNewRequestBuilderResult struct {
	r network.RequestBuilder
}

//Expect specifies that invocation of ConsensusNetwork.NewRequestBuilder is expected from 1 to Infinity times
func (m *mConsensusNetworkMockNewRequestBuilder) Expect() *mConsensusNetworkMockNewRequestBuilder {
	m.mock.NewRequestBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockNewRequestBuilderExpectation{}
	}

	return m
}

//Return specifies results of invocation of ConsensusNetwork.NewRequestBuilder
func (m *mConsensusNetworkMockNewRequestBuilder) Return(r network.RequestBuilder) *ConsensusNetworkMock {
	m.mock.NewRequestBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockNewRequestBuilderExpectation{}
	}
	m.mainExpectation.result = &ConsensusNetworkMockNewRequestBuilderResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.NewRequestBuilder is expected once
func (m *mConsensusNetworkMockNewRequestBuilder) ExpectOnce() *ConsensusNetworkMockNewRequestBuilderExpectation {
	m.mock.NewRequestBuilderFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockNewRequestBuilderExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConsensusNetworkMockNewRequestBuilderExpectation) Return(r network.RequestBuilder) {
	e.result = &ConsensusNetworkMockNewRequestBuilderResult{r}
}

//Set uses given function f as a mock of ConsensusNetwork.NewRequestBuilder method
func (m *mConsensusNetworkMockNewRequestBuilder) Set(f func() (r network.RequestBuilder)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NewRequestBuilderFunc = f
	return m.mock
}

//NewRequestBuilder implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) NewRequestBuilder() (r network.RequestBuilder) {
	counter := atomic.AddUint64(&m.NewRequestBuilderPreCounter, 1)
	defer atomic.AddUint64(&m.NewRequestBuilderCounter, 1)

	if len(m.NewRequestBuilderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NewRequestBuilderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.NewRequestBuilder.")
			return
		}

		result := m.NewRequestBuilderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.NewRequestBuilder")
			return
		}

		r = result.r

		return
	}

	if m.NewRequestBuilderMock.mainExpectation != nil {

		result := m.NewRequestBuilderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.NewRequestBuilder")
		}

		r = result.r

		return
	}

	if m.NewRequestBuilderFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.NewRequestBuilder.")
		return
	}

	return m.NewRequestBuilderFunc()
}

//NewRequestBuilderMinimockCounter returns a count of ConsensusNetworkMock.NewRequestBuilderFunc invocations
func (m *ConsensusNetworkMock) NewRequestBuilderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewRequestBuilderCounter)
}

//NewRequestBuilderMinimockPreCounter returns the value of ConsensusNetworkMock.NewRequestBuilder invocations
func (m *ConsensusNetworkMock) NewRequestBuilderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewRequestBuilderPreCounter)
}

//NewRequestBuilderFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) NewRequestBuilderFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NewRequestBuilderMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NewRequestBuilderCounter) == uint64(len(m.NewRequestBuilderMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NewRequestBuilderMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NewRequestBuilderCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NewRequestBuilderFunc != nil {
		return atomic.LoadUint64(&m.NewRequestBuilderCounter) > 0
	}

	return true
}

type mConsensusNetworkMockPublicAddress struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockPublicAddressExpectation
	expectationSeries []*ConsensusNetworkMockPublicAddressExpectation
}

type ConsensusNetworkMockPublicAddressExpectation struct {
	result *ConsensusNetworkMockPublicAddressResult
}

type ConsensusNetworkMockPublicAddressResult struct {
	r string
}

//Expect specifies that invocation of ConsensusNetwork.PublicAddress is expected from 1 to Infinity times
func (m *mConsensusNetworkMockPublicAddress) Expect() *mConsensusNetworkMockPublicAddress {
	m.mock.PublicAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockPublicAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of ConsensusNetwork.PublicAddress
func (m *mConsensusNetworkMockPublicAddress) Return(r string) *ConsensusNetworkMock {
	m.mock.PublicAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockPublicAddressExpectation{}
	}
	m.mainExpectation.result = &ConsensusNetworkMockPublicAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.PublicAddress is expected once
func (m *mConsensusNetworkMockPublicAddress) ExpectOnce() *ConsensusNetworkMockPublicAddressExpectation {
	m.mock.PublicAddressFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockPublicAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConsensusNetworkMockPublicAddressExpectation) Return(r string) {
	e.result = &ConsensusNetworkMockPublicAddressResult{r}
}

//Set uses given function f as a mock of ConsensusNetwork.PublicAddress method
func (m *mConsensusNetworkMockPublicAddress) Set(f func() (r string)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PublicAddressFunc = f
	return m.mock
}

//PublicAddress implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) PublicAddress() (r string) {
	counter := atomic.AddUint64(&m.PublicAddressPreCounter, 1)
	defer atomic.AddUint64(&m.PublicAddressCounter, 1)

	if len(m.PublicAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PublicAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.PublicAddress.")
			return
		}

		result := m.PublicAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.PublicAddress")
			return
		}

		r = result.r

		return
	}

	if m.PublicAddressMock.mainExpectation != nil {

		result := m.PublicAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.PublicAddress")
		}

		r = result.r

		return
	}

	if m.PublicAddressFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.PublicAddress.")
		return
	}

	return m.PublicAddressFunc()
}

//PublicAddressMinimockCounter returns a count of ConsensusNetworkMock.PublicAddressFunc invocations
func (m *ConsensusNetworkMock) PublicAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PublicAddressCounter)
}

//PublicAddressMinimockPreCounter returns the value of ConsensusNetworkMock.PublicAddress invocations
func (m *ConsensusNetworkMock) PublicAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PublicAddressPreCounter)
}

//PublicAddressFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) PublicAddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PublicAddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PublicAddressCounter) == uint64(len(m.PublicAddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PublicAddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PublicAddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PublicAddressFunc != nil {
		return atomic.LoadUint64(&m.PublicAddressCounter) > 0
	}

	return true
}

type mConsensusNetworkMockRegisterRequestHandler struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockRegisterRequestHandlerExpectation
	expectationSeries []*ConsensusNetworkMockRegisterRequestHandlerExpectation
}

type ConsensusNetworkMockRegisterRequestHandlerExpectation struct {
	input *ConsensusNetworkMockRegisterRequestHandlerInput
}

type ConsensusNetworkMockRegisterRequestHandlerInput struct {
	p  types.PacketType
	p1 network.ConsensusRequestHandler
}

//Expect specifies that invocation of ConsensusNetwork.RegisterRequestHandler is expected from 1 to Infinity times
func (m *mConsensusNetworkMockRegisterRequestHandler) Expect(p types.PacketType, p1 network.ConsensusRequestHandler) *mConsensusNetworkMockRegisterRequestHandler {
	m.mock.RegisterRequestHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockRegisterRequestHandlerExpectation{}
	}
	m.mainExpectation.input = &ConsensusNetworkMockRegisterRequestHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of ConsensusNetwork.RegisterRequestHandler
func (m *mConsensusNetworkMockRegisterRequestHandler) Return() *ConsensusNetworkMock {
	m.mock.RegisterRequestHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockRegisterRequestHandlerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.RegisterRequestHandler is expected once
func (m *mConsensusNetworkMockRegisterRequestHandler) ExpectOnce(p types.PacketType, p1 network.ConsensusRequestHandler) *ConsensusNetworkMockRegisterRequestHandlerExpectation {
	m.mock.RegisterRequestHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockRegisterRequestHandlerExpectation{}
	expectation.input = &ConsensusNetworkMockRegisterRequestHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ConsensusNetwork.RegisterRequestHandler method
func (m *mConsensusNetworkMockRegisterRequestHandler) Set(f func(p types.PacketType, p1 network.ConsensusRequestHandler)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterRequestHandlerFunc = f
	return m.mock
}

//RegisterRequestHandler implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) RegisterRequestHandler(p types.PacketType, p1 network.ConsensusRequestHandler) {
	counter := atomic.AddUint64(&m.RegisterRequestHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestHandlerCounter, 1)

	if len(m.RegisterRequestHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterRequestHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.RegisterRequestHandler. %v %v", p, p1)
			return
		}

		input := m.RegisterRequestHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConsensusNetworkMockRegisterRequestHandlerInput{p, p1}, "ConsensusNetwork.RegisterRequestHandler got unexpected parameters")

		return
	}

	if m.RegisterRequestHandlerMock.mainExpectation != nil {

		input := m.RegisterRequestHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConsensusNetworkMockRegisterRequestHandlerInput{p, p1}, "ConsensusNetwork.RegisterRequestHandler got unexpected parameters")
		}

		return
	}

	if m.RegisterRequestHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.RegisterRequestHandler. %v %v", p, p1)
		return
	}

	m.RegisterRequestHandlerFunc(p, p1)
}

//RegisterRequestHandlerMinimockCounter returns a count of ConsensusNetworkMock.RegisterRequestHandlerFunc invocations
func (m *ConsensusNetworkMock) RegisterRequestHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestHandlerCounter)
}

//RegisterRequestHandlerMinimockPreCounter returns the value of ConsensusNetworkMock.RegisterRequestHandler invocations
func (m *ConsensusNetworkMock) RegisterRequestHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestHandlerPreCounter)
}

//RegisterRequestHandlerFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) RegisterRequestHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterRequestHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == uint64(len(m.RegisterRequestHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterRequestHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterRequestHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterRequestHandlerFunc != nil {
		return atomic.LoadUint64(&m.RegisterRequestHandlerCounter) > 0
	}

	return true
}

type mConsensusNetworkMockSendRequest struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockSendRequestExpectation
	expectationSeries []*ConsensusNetworkMockSendRequestExpectation
}

type ConsensusNetworkMockSendRequestExpectation struct {
	input  *ConsensusNetworkMockSendRequestInput
	result *ConsensusNetworkMockSendRequestResult
}

type ConsensusNetworkMockSendRequestInput struct {
	p  network.Request
	p1 core.RecordRef
}

type ConsensusNetworkMockSendRequestResult struct {
	r error
}

//Expect specifies that invocation of ConsensusNetwork.SendRequest is expected from 1 to Infinity times
func (m *mConsensusNetworkMockSendRequest) Expect(p network.Request, p1 core.RecordRef) *mConsensusNetworkMockSendRequest {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockSendRequestExpectation{}
	}
	m.mainExpectation.input = &ConsensusNetworkMockSendRequestInput{p, p1}
	return m
}

//Return specifies results of invocation of ConsensusNetwork.SendRequest
func (m *mConsensusNetworkMockSendRequest) Return(r error) *ConsensusNetworkMock {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockSendRequestExpectation{}
	}
	m.mainExpectation.result = &ConsensusNetworkMockSendRequestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.SendRequest is expected once
func (m *mConsensusNetworkMockSendRequest) ExpectOnce(p network.Request, p1 core.RecordRef) *ConsensusNetworkMockSendRequestExpectation {
	m.mock.SendRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockSendRequestExpectation{}
	expectation.input = &ConsensusNetworkMockSendRequestInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConsensusNetworkMockSendRequestExpectation) Return(r error) {
	e.result = &ConsensusNetworkMockSendRequestResult{r}
}

//Set uses given function f as a mock of ConsensusNetwork.SendRequest method
func (m *mConsensusNetworkMockSendRequest) Set(f func(p network.Request, p1 core.RecordRef) (r error)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRequestFunc = f
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) SendRequest(p network.Request, p1 core.RecordRef) (r error) {
	counter := atomic.AddUint64(&m.SendRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestCounter, 1)

	if len(m.SendRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.SendRequest. %v %v", p, p1)
			return
		}

		input := m.SendRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConsensusNetworkMockSendRequestInput{p, p1}, "ConsensusNetwork.SendRequest got unexpected parameters")

		result := m.SendRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.SendRequest")
			return
		}

		r = result.r

		return
	}

	if m.SendRequestMock.mainExpectation != nil {

		input := m.SendRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConsensusNetworkMockSendRequestInput{p, p1}, "ConsensusNetwork.SendRequest got unexpected parameters")
		}

		result := m.SendRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.SendRequest")
		}

		r = result.r

		return
	}

	if m.SendRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.SendRequest. %v %v", p, p1)
		return
	}

	return m.SendRequestFunc(p, p1)
}

//SendRequestMinimockCounter returns a count of ConsensusNetworkMock.SendRequestFunc invocations
func (m *ConsensusNetworkMock) SendRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestCounter)
}

//SendRequestMinimockPreCounter returns the value of ConsensusNetworkMock.SendRequest invocations
func (m *ConsensusNetworkMock) SendRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestPreCounter)
}

//SendRequestFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) SendRequestFinished() bool {
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

type mConsensusNetworkMockStart struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockStartExpectation
	expectationSeries []*ConsensusNetworkMockStartExpectation
}

type ConsensusNetworkMockStartExpectation struct {
	input *ConsensusNetworkMockStartInput
}

type ConsensusNetworkMockStartInput struct {
	p context.Context
}

//Expect specifies that invocation of ConsensusNetwork.Start is expected from 1 to Infinity times
func (m *mConsensusNetworkMockStart) Expect(p context.Context) *mConsensusNetworkMockStart {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockStartExpectation{}
	}
	m.mainExpectation.input = &ConsensusNetworkMockStartInput{p}
	return m
}

//Return specifies results of invocation of ConsensusNetwork.Start
func (m *mConsensusNetworkMockStart) Return() *ConsensusNetworkMock {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockStartExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.Start is expected once
func (m *mConsensusNetworkMockStart) ExpectOnce(p context.Context) *ConsensusNetworkMockStartExpectation {
	m.mock.StartFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockStartExpectation{}
	expectation.input = &ConsensusNetworkMockStartInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ConsensusNetwork.Start method
func (m *mConsensusNetworkMockStart) Set(f func(p context.Context)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StartFunc = f
	return m.mock
}

//Start implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) Start(p context.Context) {
	counter := atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if len(m.StartMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StartMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.Start. %v", p)
			return
		}

		input := m.StartMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConsensusNetworkMockStartInput{p}, "ConsensusNetwork.Start got unexpected parameters")

		return
	}

	if m.StartMock.mainExpectation != nil {

		input := m.StartMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConsensusNetworkMockStartInput{p}, "ConsensusNetwork.Start got unexpected parameters")
		}

		return
	}

	if m.StartFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.Start. %v", p)
		return
	}

	m.StartFunc(p)
}

//StartMinimockCounter returns a count of ConsensusNetworkMock.StartFunc invocations
func (m *ConsensusNetworkMock) StartMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StartCounter)
}

//StartMinimockPreCounter returns the value of ConsensusNetworkMock.Start invocations
func (m *ConsensusNetworkMock) StartMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StartPreCounter)
}

//StartFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) StartFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StartMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StartCounter) == uint64(len(m.StartMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StartMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StartFunc != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	return true
}

type mConsensusNetworkMockStop struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockStopExpectation
	expectationSeries []*ConsensusNetworkMockStopExpectation
}

type ConsensusNetworkMockStopExpectation struct {
}

//Expect specifies that invocation of ConsensusNetwork.Stop is expected from 1 to Infinity times
func (m *mConsensusNetworkMockStop) Expect() *mConsensusNetworkMockStop {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockStopExpectation{}
	}

	return m
}

//Return specifies results of invocation of ConsensusNetwork.Stop
func (m *mConsensusNetworkMockStop) Return() *ConsensusNetworkMock {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockStopExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.Stop is expected once
func (m *mConsensusNetworkMockStop) ExpectOnce() *ConsensusNetworkMockStopExpectation {
	m.mock.StopFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockStopExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ConsensusNetwork.Stop method
func (m *mConsensusNetworkMockStop) Set(f func()) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StopFunc = f
	return m.mock
}

//Stop implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) Stop() {
	counter := atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if len(m.StopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.Stop.")
			return
		}

		return
	}

	if m.StopMock.mainExpectation != nil {

		return
	}

	if m.StopFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.Stop.")
		return
	}

	m.StopFunc()
}

//StopMinimockCounter returns a count of ConsensusNetworkMock.StopFunc invocations
func (m *ConsensusNetworkMock) StopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StopCounter)
}

//StopMinimockPreCounter returns the value of ConsensusNetworkMock.Stop invocations
func (m *ConsensusNetworkMock) StopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StopPreCounter)
}

//StopFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) StopFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StopMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StopCounter) == uint64(len(m.StopMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StopMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StopCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StopFunc != nil {
		return atomic.LoadUint64(&m.StopCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ConsensusNetworkMock) ValidateCallCounters() {

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.GetNodeID")
	}

	if !m.NewRequestBuilderFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.NewRequestBuilder")
	}

	if !m.PublicAddressFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.PublicAddress")
	}

	if !m.RegisterRequestHandlerFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.RegisterRequestHandler")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.SendRequest")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.Stop")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ConsensusNetworkMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ConsensusNetworkMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ConsensusNetworkMock) MinimockFinish() {

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.GetNodeID")
	}

	if !m.NewRequestBuilderFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.NewRequestBuilder")
	}

	if !m.PublicAddressFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.PublicAddress")
	}

	if !m.RegisterRequestHandlerFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.RegisterRequestHandler")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.SendRequest")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.Stop")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ConsensusNetworkMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ConsensusNetworkMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.NewRequestBuilderFinished()
		ok = ok && m.PublicAddressFinished()
		ok = ok && m.RegisterRequestHandlerFinished()
		ok = ok && m.SendRequestFinished()
		ok = ok && m.StartFinished()
		ok = ok && m.StopFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.GetNodeID")
			}

			if !m.NewRequestBuilderFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.NewRequestBuilder")
			}

			if !m.PublicAddressFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.PublicAddress")
			}

			if !m.RegisterRequestHandlerFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.RegisterRequestHandler")
			}

			if !m.SendRequestFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.SendRequest")
			}

			if !m.StartFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.Start")
			}

			if !m.StopFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.Stop")
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
func (m *ConsensusNetworkMock) AllMocksCalled() bool {

	if !m.GetNodeIDFinished() {
		return false
	}

	if !m.NewRequestBuilderFinished() {
		return false
	}

	if !m.PublicAddressFinished() {
		return false
	}

	if !m.RegisterRequestHandlerFinished() {
		return false
	}

	if !m.SendRequestFinished() {
		return false
	}

	if !m.StartFinished() {
		return false
	}

	if !m.StopFinished() {
		return false
	}

	return true
}

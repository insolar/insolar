package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "HostNetwork" can be found in github.com/insolar/insolar/network
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

//HostNetworkMock implements github.com/insolar/insolar/network.HostNetwork
type HostNetworkMock struct {
	t minimock.Tester

	BuildResponseFunc       func(p network.Request, p1 interface{}) (r network.Response)
	BuildResponseCounter    uint64
	BuildResponsePreCounter uint64
	BuildResponseMock       mHostNetworkMockBuildResponse

	GetNodeIDFunc       func() (r core.RecordRef)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mHostNetworkMockGetNodeID

	NewRequestBuilderFunc       func() (r network.RequestBuilder)
	NewRequestBuilderCounter    uint64
	NewRequestBuilderPreCounter uint64
	NewRequestBuilderMock       mHostNetworkMockNewRequestBuilder

	PublicAddressFunc       func() (r string)
	PublicAddressCounter    uint64
	PublicAddressPreCounter uint64
	PublicAddressMock       mHostNetworkMockPublicAddress

	RegisterRequestHandlerFunc       func(p types.PacketType, p1 network.RequestHandler)
	RegisterRequestHandlerCounter    uint64
	RegisterRequestHandlerPreCounter uint64
	RegisterRequestHandlerMock       mHostNetworkMockRegisterRequestHandler

	SendRequestFunc       func(p network.Request, p1 core.RecordRef) (r network.Future, r1 error)
	SendRequestCounter    uint64
	SendRequestPreCounter uint64
	SendRequestMock       mHostNetworkMockSendRequest

	StartFunc       func(p context.Context)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mHostNetworkMockStart

	StopFunc       func()
	StopCounter    uint64
	StopPreCounter uint64
	StopMock       mHostNetworkMockStop
}

//NewHostNetworkMock returns a mock for github.com/insolar/insolar/network.HostNetwork
func NewHostNetworkMock(t minimock.Tester) *HostNetworkMock {
	m := &HostNetworkMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BuildResponseMock = mHostNetworkMockBuildResponse{mock: m}
	m.GetNodeIDMock = mHostNetworkMockGetNodeID{mock: m}
	m.NewRequestBuilderMock = mHostNetworkMockNewRequestBuilder{mock: m}
	m.PublicAddressMock = mHostNetworkMockPublicAddress{mock: m}
	m.RegisterRequestHandlerMock = mHostNetworkMockRegisterRequestHandler{mock: m}
	m.SendRequestMock = mHostNetworkMockSendRequest{mock: m}
	m.StartMock = mHostNetworkMockStart{mock: m}
	m.StopMock = mHostNetworkMockStop{mock: m}

	return m
}

type mHostNetworkMockBuildResponse struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockBuildResponseExpectation
	expectationSeries []*HostNetworkMockBuildResponseExpectation
}

type HostNetworkMockBuildResponseExpectation struct {
	input  *HostNetworkMockBuildResponseInput
	result *HostNetworkMockBuildResponseResult
}

type HostNetworkMockBuildResponseInput struct {
	p  network.Request
	p1 interface{}
}

type HostNetworkMockBuildResponseResult struct {
	r network.Response
}

//Expect specifies that invocation of HostNetwork.BuildResponse is expected from 1 to Infinity times
func (m *mHostNetworkMockBuildResponse) Expect(p network.Request, p1 interface{}) *mHostNetworkMockBuildResponse {
	m.mock.BuildResponseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockBuildResponseExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockBuildResponseInput{p, p1}
	return m
}

//Return specifies results of invocation of HostNetwork.BuildResponse
func (m *mHostNetworkMockBuildResponse) Return(r network.Response) *HostNetworkMock {
	m.mock.BuildResponseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockBuildResponseExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockBuildResponseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.BuildResponse is expected once
func (m *mHostNetworkMockBuildResponse) ExpectOnce(p network.Request, p1 interface{}) *HostNetworkMockBuildResponseExpectation {
	m.mock.BuildResponseFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockBuildResponseExpectation{}
	expectation.input = &HostNetworkMockBuildResponseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockBuildResponseExpectation) Return(r network.Response) {
	e.result = &HostNetworkMockBuildResponseResult{r}
}

//Set uses given function f as a mock of HostNetwork.BuildResponse method
func (m *mHostNetworkMockBuildResponse) Set(f func(p network.Request, p1 interface{}) (r network.Response)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BuildResponseFunc = f
	return m.mock
}

//BuildResponse implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) BuildResponse(p network.Request, p1 interface{}) (r network.Response) {
	counter := atomic.AddUint64(&m.BuildResponsePreCounter, 1)
	defer atomic.AddUint64(&m.BuildResponseCounter, 1)

	if len(m.BuildResponseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BuildResponseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.BuildResponse. %v %v", p, p1)
			return
		}

		input := m.BuildResponseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockBuildResponseInput{p, p1}, "HostNetwork.BuildResponse got unexpected parameters")

		result := m.BuildResponseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.BuildResponse")
			return
		}

		r = result.r

		return
	}

	if m.BuildResponseMock.mainExpectation != nil {

		input := m.BuildResponseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockBuildResponseInput{p, p1}, "HostNetwork.BuildResponse got unexpected parameters")
		}

		result := m.BuildResponseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.BuildResponse")
		}

		r = result.r

		return
	}

	if m.BuildResponseFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.BuildResponse. %v %v", p, p1)
		return
	}

	return m.BuildResponseFunc(p, p1)
}

//BuildResponseMinimockCounter returns a count of HostNetworkMock.BuildResponseFunc invocations
func (m *HostNetworkMock) BuildResponseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BuildResponseCounter)
}

//BuildResponseMinimockPreCounter returns the value of HostNetworkMock.BuildResponse invocations
func (m *HostNetworkMock) BuildResponseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BuildResponsePreCounter)
}

//BuildResponseFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) BuildResponseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.BuildResponseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.BuildResponseCounter) == uint64(len(m.BuildResponseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.BuildResponseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.BuildResponseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.BuildResponseFunc != nil {
		return atomic.LoadUint64(&m.BuildResponseCounter) > 0
	}

	return true
}

type mHostNetworkMockGetNodeID struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockGetNodeIDExpectation
	expectationSeries []*HostNetworkMockGetNodeIDExpectation
}

type HostNetworkMockGetNodeIDExpectation struct {
	result *HostNetworkMockGetNodeIDResult
}

type HostNetworkMockGetNodeIDResult struct {
	r core.RecordRef
}

//Expect specifies that invocation of HostNetwork.GetNodeID is expected from 1 to Infinity times
func (m *mHostNetworkMockGetNodeID) Expect() *mHostNetworkMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of HostNetwork.GetNodeID
func (m *mHostNetworkMockGetNodeID) Return(r core.RecordRef) *HostNetworkMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.GetNodeID is expected once
func (m *mHostNetworkMockGetNodeID) ExpectOnce() *HostNetworkMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockGetNodeIDExpectation) Return(r core.RecordRef) {
	e.result = &HostNetworkMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of HostNetwork.GetNodeID method
func (m *mHostNetworkMockGetNodeID) Set(f func() (r core.RecordRef)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) GetNodeID() (r core.RecordRef) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of HostNetworkMock.GetNodeIDFunc invocations
func (m *HostNetworkMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of HostNetworkMock.GetNodeID invocations
func (m *HostNetworkMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) GetNodeIDFinished() bool {
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

type mHostNetworkMockNewRequestBuilder struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockNewRequestBuilderExpectation
	expectationSeries []*HostNetworkMockNewRequestBuilderExpectation
}

type HostNetworkMockNewRequestBuilderExpectation struct {
	result *HostNetworkMockNewRequestBuilderResult
}

type HostNetworkMockNewRequestBuilderResult struct {
	r network.RequestBuilder
}

//Expect specifies that invocation of HostNetwork.NewRequestBuilder is expected from 1 to Infinity times
func (m *mHostNetworkMockNewRequestBuilder) Expect() *mHostNetworkMockNewRequestBuilder {
	m.mock.NewRequestBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockNewRequestBuilderExpectation{}
	}

	return m
}

//Return specifies results of invocation of HostNetwork.NewRequestBuilder
func (m *mHostNetworkMockNewRequestBuilder) Return(r network.RequestBuilder) *HostNetworkMock {
	m.mock.NewRequestBuilderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockNewRequestBuilderExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockNewRequestBuilderResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.NewRequestBuilder is expected once
func (m *mHostNetworkMockNewRequestBuilder) ExpectOnce() *HostNetworkMockNewRequestBuilderExpectation {
	m.mock.NewRequestBuilderFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockNewRequestBuilderExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockNewRequestBuilderExpectation) Return(r network.RequestBuilder) {
	e.result = &HostNetworkMockNewRequestBuilderResult{r}
}

//Set uses given function f as a mock of HostNetwork.NewRequestBuilder method
func (m *mHostNetworkMockNewRequestBuilder) Set(f func() (r network.RequestBuilder)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NewRequestBuilderFunc = f
	return m.mock
}

//NewRequestBuilder implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) NewRequestBuilder() (r network.RequestBuilder) {
	counter := atomic.AddUint64(&m.NewRequestBuilderPreCounter, 1)
	defer atomic.AddUint64(&m.NewRequestBuilderCounter, 1)

	if len(m.NewRequestBuilderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NewRequestBuilderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.NewRequestBuilder.")
			return
		}

		result := m.NewRequestBuilderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.NewRequestBuilder")
			return
		}

		r = result.r

		return
	}

	if m.NewRequestBuilderMock.mainExpectation != nil {

		result := m.NewRequestBuilderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.NewRequestBuilder")
		}

		r = result.r

		return
	}

	if m.NewRequestBuilderFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.NewRequestBuilder.")
		return
	}

	return m.NewRequestBuilderFunc()
}

//NewRequestBuilderMinimockCounter returns a count of HostNetworkMock.NewRequestBuilderFunc invocations
func (m *HostNetworkMock) NewRequestBuilderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewRequestBuilderCounter)
}

//NewRequestBuilderMinimockPreCounter returns the value of HostNetworkMock.NewRequestBuilder invocations
func (m *HostNetworkMock) NewRequestBuilderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewRequestBuilderPreCounter)
}

//NewRequestBuilderFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) NewRequestBuilderFinished() bool {
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

type mHostNetworkMockPublicAddress struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockPublicAddressExpectation
	expectationSeries []*HostNetworkMockPublicAddressExpectation
}

type HostNetworkMockPublicAddressExpectation struct {
	result *HostNetworkMockPublicAddressResult
}

type HostNetworkMockPublicAddressResult struct {
	r string
}

//Expect specifies that invocation of HostNetwork.PublicAddress is expected from 1 to Infinity times
func (m *mHostNetworkMockPublicAddress) Expect() *mHostNetworkMockPublicAddress {
	m.mock.PublicAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockPublicAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of HostNetwork.PublicAddress
func (m *mHostNetworkMockPublicAddress) Return(r string) *HostNetworkMock {
	m.mock.PublicAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockPublicAddressExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockPublicAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.PublicAddress is expected once
func (m *mHostNetworkMockPublicAddress) ExpectOnce() *HostNetworkMockPublicAddressExpectation {
	m.mock.PublicAddressFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockPublicAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockPublicAddressExpectation) Return(r string) {
	e.result = &HostNetworkMockPublicAddressResult{r}
}

//Set uses given function f as a mock of HostNetwork.PublicAddress method
func (m *mHostNetworkMockPublicAddress) Set(f func() (r string)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PublicAddressFunc = f
	return m.mock
}

//PublicAddress implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) PublicAddress() (r string) {
	counter := atomic.AddUint64(&m.PublicAddressPreCounter, 1)
	defer atomic.AddUint64(&m.PublicAddressCounter, 1)

	if len(m.PublicAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PublicAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.PublicAddress.")
			return
		}

		result := m.PublicAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.PublicAddress")
			return
		}

		r = result.r

		return
	}

	if m.PublicAddressMock.mainExpectation != nil {

		result := m.PublicAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.PublicAddress")
		}

		r = result.r

		return
	}

	if m.PublicAddressFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.PublicAddress.")
		return
	}

	return m.PublicAddressFunc()
}

//PublicAddressMinimockCounter returns a count of HostNetworkMock.PublicAddressFunc invocations
func (m *HostNetworkMock) PublicAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PublicAddressCounter)
}

//PublicAddressMinimockPreCounter returns the value of HostNetworkMock.PublicAddress invocations
func (m *HostNetworkMock) PublicAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PublicAddressPreCounter)
}

//PublicAddressFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) PublicAddressFinished() bool {
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

type mHostNetworkMockRegisterRequestHandler struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockRegisterRequestHandlerExpectation
	expectationSeries []*HostNetworkMockRegisterRequestHandlerExpectation
}

type HostNetworkMockRegisterRequestHandlerExpectation struct {
	input *HostNetworkMockRegisterRequestHandlerInput
}

type HostNetworkMockRegisterRequestHandlerInput struct {
	p  types.PacketType
	p1 network.RequestHandler
}

//Expect specifies that invocation of HostNetwork.RegisterRequestHandler is expected from 1 to Infinity times
func (m *mHostNetworkMockRegisterRequestHandler) Expect(p types.PacketType, p1 network.RequestHandler) *mHostNetworkMockRegisterRequestHandler {
	m.mock.RegisterRequestHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockRegisterRequestHandlerExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockRegisterRequestHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of HostNetwork.RegisterRequestHandler
func (m *mHostNetworkMockRegisterRequestHandler) Return() *HostNetworkMock {
	m.mock.RegisterRequestHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockRegisterRequestHandlerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.RegisterRequestHandler is expected once
func (m *mHostNetworkMockRegisterRequestHandler) ExpectOnce(p types.PacketType, p1 network.RequestHandler) *HostNetworkMockRegisterRequestHandlerExpectation {
	m.mock.RegisterRequestHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockRegisterRequestHandlerExpectation{}
	expectation.input = &HostNetworkMockRegisterRequestHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of HostNetwork.RegisterRequestHandler method
func (m *mHostNetworkMockRegisterRequestHandler) Set(f func(p types.PacketType, p1 network.RequestHandler)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterRequestHandlerFunc = f
	return m.mock
}

//RegisterRequestHandler implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) RegisterRequestHandler(p types.PacketType, p1 network.RequestHandler) {
	counter := atomic.AddUint64(&m.RegisterRequestHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestHandlerCounter, 1)

	if len(m.RegisterRequestHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterRequestHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.RegisterRequestHandler. %v %v", p, p1)
			return
		}

		input := m.RegisterRequestHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockRegisterRequestHandlerInput{p, p1}, "HostNetwork.RegisterRequestHandler got unexpected parameters")

		return
	}

	if m.RegisterRequestHandlerMock.mainExpectation != nil {

		input := m.RegisterRequestHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockRegisterRequestHandlerInput{p, p1}, "HostNetwork.RegisterRequestHandler got unexpected parameters")
		}

		return
	}

	if m.RegisterRequestHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.RegisterRequestHandler. %v %v", p, p1)
		return
	}

	m.RegisterRequestHandlerFunc(p, p1)
}

//RegisterRequestHandlerMinimockCounter returns a count of HostNetworkMock.RegisterRequestHandlerFunc invocations
func (m *HostNetworkMock) RegisterRequestHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestHandlerCounter)
}

//RegisterRequestHandlerMinimockPreCounter returns the value of HostNetworkMock.RegisterRequestHandler invocations
func (m *HostNetworkMock) RegisterRequestHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestHandlerPreCounter)
}

//RegisterRequestHandlerFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) RegisterRequestHandlerFinished() bool {
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

type mHostNetworkMockSendRequest struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockSendRequestExpectation
	expectationSeries []*HostNetworkMockSendRequestExpectation
}

type HostNetworkMockSendRequestExpectation struct {
	input  *HostNetworkMockSendRequestInput
	result *HostNetworkMockSendRequestResult
}

type HostNetworkMockSendRequestInput struct {
	p  network.Request
	p1 core.RecordRef
}

type HostNetworkMockSendRequestResult struct {
	r  network.Future
	r1 error
}

//Expect specifies that invocation of HostNetwork.SendRequest is expected from 1 to Infinity times
func (m *mHostNetworkMockSendRequest) Expect(p network.Request, p1 core.RecordRef) *mHostNetworkMockSendRequest {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockSendRequestExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockSendRequestInput{p, p1}
	return m
}

//Return specifies results of invocation of HostNetwork.SendRequest
func (m *mHostNetworkMockSendRequest) Return(r network.Future, r1 error) *HostNetworkMock {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockSendRequestExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockSendRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.SendRequest is expected once
func (m *mHostNetworkMockSendRequest) ExpectOnce(p network.Request, p1 core.RecordRef) *HostNetworkMockSendRequestExpectation {
	m.mock.SendRequestFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockSendRequestExpectation{}
	expectation.input = &HostNetworkMockSendRequestInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockSendRequestExpectation) Return(r network.Future, r1 error) {
	e.result = &HostNetworkMockSendRequestResult{r, r1}
}

//Set uses given function f as a mock of HostNetwork.SendRequest method
func (m *mHostNetworkMockSendRequest) Set(f func(p network.Request, p1 core.RecordRef) (r network.Future, r1 error)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRequestFunc = f
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) SendRequest(p network.Request, p1 core.RecordRef) (r network.Future, r1 error) {
	counter := atomic.AddUint64(&m.SendRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestCounter, 1)

	if len(m.SendRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.SendRequest. %v %v", p, p1)
			return
		}

		input := m.SendRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockSendRequestInput{p, p1}, "HostNetwork.SendRequest got unexpected parameters")

		result := m.SendRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.SendRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestMock.mainExpectation != nil {

		input := m.SendRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockSendRequestInput{p, p1}, "HostNetwork.SendRequest got unexpected parameters")
		}

		result := m.SendRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.SendRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.SendRequest. %v %v", p, p1)
		return
	}

	return m.SendRequestFunc(p, p1)
}

//SendRequestMinimockCounter returns a count of HostNetworkMock.SendRequestFunc invocations
func (m *HostNetworkMock) SendRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestCounter)
}

//SendRequestMinimockPreCounter returns the value of HostNetworkMock.SendRequest invocations
func (m *HostNetworkMock) SendRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestPreCounter)
}

//SendRequestFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) SendRequestFinished() bool {
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

type mHostNetworkMockStart struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockStartExpectation
	expectationSeries []*HostNetworkMockStartExpectation
}

type HostNetworkMockStartExpectation struct {
	input *HostNetworkMockStartInput
}

type HostNetworkMockStartInput struct {
	p context.Context
}

//Expect specifies that invocation of HostNetwork.Start is expected from 1 to Infinity times
func (m *mHostNetworkMockStart) Expect(p context.Context) *mHostNetworkMockStart {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockStartExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockStartInput{p}
	return m
}

//Return specifies results of invocation of HostNetwork.Start
func (m *mHostNetworkMockStart) Return() *HostNetworkMock {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockStartExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.Start is expected once
func (m *mHostNetworkMockStart) ExpectOnce(p context.Context) *HostNetworkMockStartExpectation {
	m.mock.StartFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockStartExpectation{}
	expectation.input = &HostNetworkMockStartInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of HostNetwork.Start method
func (m *mHostNetworkMockStart) Set(f func(p context.Context)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StartFunc = f
	return m.mock
}

//Start implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) Start(p context.Context) {
	counter := atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if len(m.StartMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StartMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.Start. %v", p)
			return
		}

		input := m.StartMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockStartInput{p}, "HostNetwork.Start got unexpected parameters")

		return
	}

	if m.StartMock.mainExpectation != nil {

		input := m.StartMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockStartInput{p}, "HostNetwork.Start got unexpected parameters")
		}

		return
	}

	if m.StartFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.Start. %v", p)
		return
	}

	m.StartFunc(p)
}

//StartMinimockCounter returns a count of HostNetworkMock.StartFunc invocations
func (m *HostNetworkMock) StartMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StartCounter)
}

//StartMinimockPreCounter returns the value of HostNetworkMock.Start invocations
func (m *HostNetworkMock) StartMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StartPreCounter)
}

//StartFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) StartFinished() bool {
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

type mHostNetworkMockStop struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockStopExpectation
	expectationSeries []*HostNetworkMockStopExpectation
}

type HostNetworkMockStopExpectation struct {
}

//Expect specifies that invocation of HostNetwork.Stop is expected from 1 to Infinity times
func (m *mHostNetworkMockStop) Expect() *mHostNetworkMockStop {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockStopExpectation{}
	}

	return m
}

//Return specifies results of invocation of HostNetwork.Stop
func (m *mHostNetworkMockStop) Return() *HostNetworkMock {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockStopExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.Stop is expected once
func (m *mHostNetworkMockStop) ExpectOnce() *HostNetworkMockStopExpectation {
	m.mock.StopFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockStopExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of HostNetwork.Stop method
func (m *mHostNetworkMockStop) Set(f func()) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StopFunc = f
	return m.mock
}

//Stop implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) Stop() {
	counter := atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if len(m.StopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.Stop.")
			return
		}

		return
	}

	if m.StopMock.mainExpectation != nil {

		return
	}

	if m.StopFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.Stop.")
		return
	}

	m.StopFunc()
}

//StopMinimockCounter returns a count of HostNetworkMock.StopFunc invocations
func (m *HostNetworkMock) StopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StopCounter)
}

//StopMinimockPreCounter returns the value of HostNetworkMock.Stop invocations
func (m *HostNetworkMock) StopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StopPreCounter)
}

//StopFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) StopFinished() bool {
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
func (m *HostNetworkMock) ValidateCallCounters() {

	if !m.BuildResponseFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.BuildResponse")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.GetNodeID")
	}

	if !m.NewRequestBuilderFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.NewRequestBuilder")
	}

	if !m.PublicAddressFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.PublicAddress")
	}

	if !m.RegisterRequestHandlerFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.RegisterRequestHandler")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.SendRequest")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.Stop")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HostNetworkMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *HostNetworkMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *HostNetworkMock) MinimockFinish() {

	if !m.BuildResponseFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.BuildResponse")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.GetNodeID")
	}

	if !m.NewRequestBuilderFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.NewRequestBuilder")
	}

	if !m.PublicAddressFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.PublicAddress")
	}

	if !m.RegisterRequestHandlerFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.RegisterRequestHandler")
	}

	if !m.SendRequestFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.SendRequest")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.Stop")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *HostNetworkMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *HostNetworkMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.BuildResponseFinished()
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

			if !m.BuildResponseFinished() {
				m.t.Error("Expected call to HostNetworkMock.BuildResponse")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to HostNetworkMock.GetNodeID")
			}

			if !m.NewRequestBuilderFinished() {
				m.t.Error("Expected call to HostNetworkMock.NewRequestBuilder")
			}

			if !m.PublicAddressFinished() {
				m.t.Error("Expected call to HostNetworkMock.PublicAddress")
			}

			if !m.RegisterRequestHandlerFinished() {
				m.t.Error("Expected call to HostNetworkMock.RegisterRequestHandler")
			}

			if !m.SendRequestFinished() {
				m.t.Error("Expected call to HostNetworkMock.SendRequest")
			}

			if !m.StartFinished() {
				m.t.Error("Expected call to HostNetworkMock.Start")
			}

			if !m.StopFinished() {
				m.t.Error("Expected call to HostNetworkMock.Stop")
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
func (m *HostNetworkMock) AllMocksCalled() bool {

	if !m.BuildResponseFinished() {
		return false
	}

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

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
	insolar "github.com/insolar/insolar/insolar"
	network "github.com/insolar/insolar/network"
	host "github.com/insolar/insolar/network/hostnetwork/host"
	types "github.com/insolar/insolar/network/hostnetwork/packet/types"

	testify_assert "github.com/stretchr/testify/assert"
)

//HostNetworkMock implements github.com/insolar/insolar/network.HostNetwork
type HostNetworkMock struct {
	t minimock.Tester

	BuildResponseFunc       func(p context.Context, p1 network.Packet, p2 interface{}) (r network.Packet)
	BuildResponseCounter    uint64
	BuildResponsePreCounter uint64
	BuildResponseMock       mHostNetworkMockBuildResponse

	InitFunc       func(p context.Context) (r error)
	InitCounter    uint64
	InitPreCounter uint64
	InitMock       mHostNetworkMockInit

	PublicAddressFunc       func() (r string)
	PublicAddressCounter    uint64
	PublicAddressPreCounter uint64
	PublicAddressMock       mHostNetworkMockPublicAddress

	RegisterRequestHandlerFunc       func(p types.PacketType, p1 network.RequestHandler)
	RegisterRequestHandlerCounter    uint64
	RegisterRequestHandlerPreCounter uint64
	RegisterRequestHandlerMock       mHostNetworkMockRegisterRequestHandler

	SendRequestFunc       func(p context.Context, p1 types.PacketType, p2 interface{}, p3 insolar.Reference) (r network.Future, r1 error)
	SendRequestCounter    uint64
	SendRequestPreCounter uint64
	SendRequestMock       mHostNetworkMockSendRequest

	SendRequestToHostFunc       func(p context.Context, p1 types.PacketType, p2 interface{}, p3 *host.Host) (r network.Future, r1 error)
	SendRequestToHostCounter    uint64
	SendRequestToHostPreCounter uint64
	SendRequestToHostMock       mHostNetworkMockSendRequestToHost

	StartFunc       func(p context.Context) (r error)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mHostNetworkMockStart

	StopFunc       func(p context.Context) (r error)
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
	m.InitMock = mHostNetworkMockInit{mock: m}
	m.PublicAddressMock = mHostNetworkMockPublicAddress{mock: m}
	m.RegisterRequestHandlerMock = mHostNetworkMockRegisterRequestHandler{mock: m}
	m.SendRequestMock = mHostNetworkMockSendRequest{mock: m}
	m.SendRequestToHostMock = mHostNetworkMockSendRequestToHost{mock: m}
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
	p  context.Context
	p1 network.Packet
	p2 interface{}
}

type HostNetworkMockBuildResponseResult struct {
	r network.Packet
}

//Expect specifies that invocation of HostNetwork.BuildResponse is expected from 1 to Infinity times
func (m *mHostNetworkMockBuildResponse) Expect(p context.Context, p1 network.Packet, p2 interface{}) *mHostNetworkMockBuildResponse {
	m.mock.BuildResponseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockBuildResponseExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockBuildResponseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of HostNetwork.BuildResponse
func (m *mHostNetworkMockBuildResponse) Return(r network.Packet) *HostNetworkMock {
	m.mock.BuildResponseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockBuildResponseExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockBuildResponseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.BuildResponse is expected once
func (m *mHostNetworkMockBuildResponse) ExpectOnce(p context.Context, p1 network.Packet, p2 interface{}) *HostNetworkMockBuildResponseExpectation {
	m.mock.BuildResponseFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockBuildResponseExpectation{}
	expectation.input = &HostNetworkMockBuildResponseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockBuildResponseExpectation) Return(r network.Packet) {
	e.result = &HostNetworkMockBuildResponseResult{r}
}

//Set uses given function f as a mock of HostNetwork.BuildResponse method
func (m *mHostNetworkMockBuildResponse) Set(f func(p context.Context, p1 network.Packet, p2 interface{}) (r network.Packet)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BuildResponseFunc = f
	return m.mock
}

//BuildResponse implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) BuildResponse(p context.Context, p1 network.Packet, p2 interface{}) (r network.Packet) {
	counter := atomic.AddUint64(&m.BuildResponsePreCounter, 1)
	defer atomic.AddUint64(&m.BuildResponseCounter, 1)

	if len(m.BuildResponseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BuildResponseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.BuildResponse. %v %v %v", p, p1, p2)
			return
		}

		input := m.BuildResponseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockBuildResponseInput{p, p1, p2}, "HostNetwork.BuildResponse got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, HostNetworkMockBuildResponseInput{p, p1, p2}, "HostNetwork.BuildResponse got unexpected parameters")
		}

		result := m.BuildResponseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.BuildResponse")
		}

		r = result.r

		return
	}

	if m.BuildResponseFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.BuildResponse. %v %v %v", p, p1, p2)
		return
	}

	return m.BuildResponseFunc(p, p1, p2)
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

type mHostNetworkMockInit struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockInitExpectation
	expectationSeries []*HostNetworkMockInitExpectation
}

type HostNetworkMockInitExpectation struct {
	input  *HostNetworkMockInitInput
	result *HostNetworkMockInitResult
}

type HostNetworkMockInitInput struct {
	p context.Context
}

type HostNetworkMockInitResult struct {
	r error
}

//Expect specifies that invocation of HostNetwork.Init is expected from 1 to Infinity times
func (m *mHostNetworkMockInit) Expect(p context.Context) *mHostNetworkMockInit {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockInitExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockInitInput{p}
	return m
}

//Return specifies results of invocation of HostNetwork.Init
func (m *mHostNetworkMockInit) Return(r error) *HostNetworkMock {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockInitExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockInitResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.Init is expected once
func (m *mHostNetworkMockInit) ExpectOnce(p context.Context) *HostNetworkMockInitExpectation {
	m.mock.InitFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockInitExpectation{}
	expectation.input = &HostNetworkMockInitInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockInitExpectation) Return(r error) {
	e.result = &HostNetworkMockInitResult{r}
}

//Set uses given function f as a mock of HostNetwork.Init method
func (m *mHostNetworkMockInit) Set(f func(p context.Context) (r error)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InitFunc = f
	return m.mock
}

//Init implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) Init(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.InitPreCounter, 1)
	defer atomic.AddUint64(&m.InitCounter, 1)

	if len(m.InitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.Init. %v", p)
			return
		}

		input := m.InitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockInitInput{p}, "HostNetwork.Init got unexpected parameters")

		result := m.InitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.Init")
			return
		}

		r = result.r

		return
	}

	if m.InitMock.mainExpectation != nil {

		input := m.InitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockInitInput{p}, "HostNetwork.Init got unexpected parameters")
		}

		result := m.InitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.Init")
		}

		r = result.r

		return
	}

	if m.InitFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.Init. %v", p)
		return
	}

	return m.InitFunc(p)
}

//InitMinimockCounter returns a count of HostNetworkMock.InitFunc invocations
func (m *HostNetworkMock) InitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InitCounter)
}

//InitMinimockPreCounter returns the value of HostNetworkMock.Init invocations
func (m *HostNetworkMock) InitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InitPreCounter)
}

//InitFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) InitFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InitMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InitCounter) == uint64(len(m.InitMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InitMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InitCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InitFunc != nil {
		return atomic.LoadUint64(&m.InitCounter) > 0
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
	p  context.Context
	p1 types.PacketType
	p2 interface{}
	p3 insolar.Reference
}

type HostNetworkMockSendRequestResult struct {
	r  network.Future
	r1 error
}

//Expect specifies that invocation of HostNetwork.SendRequest is expected from 1 to Infinity times
func (m *mHostNetworkMockSendRequest) Expect(p context.Context, p1 types.PacketType, p2 interface{}, p3 insolar.Reference) *mHostNetworkMockSendRequest {
	m.mock.SendRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockSendRequestExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockSendRequestInput{p, p1, p2, p3}
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
func (m *mHostNetworkMockSendRequest) ExpectOnce(p context.Context, p1 types.PacketType, p2 interface{}, p3 insolar.Reference) *HostNetworkMockSendRequestExpectation {
	m.mock.SendRequestFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockSendRequestExpectation{}
	expectation.input = &HostNetworkMockSendRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockSendRequestExpectation) Return(r network.Future, r1 error) {
	e.result = &HostNetworkMockSendRequestResult{r, r1}
}

//Set uses given function f as a mock of HostNetwork.SendRequest method
func (m *mHostNetworkMockSendRequest) Set(f func(p context.Context, p1 types.PacketType, p2 interface{}, p3 insolar.Reference) (r network.Future, r1 error)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRequestFunc = f
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) SendRequest(p context.Context, p1 types.PacketType, p2 interface{}, p3 insolar.Reference) (r network.Future, r1 error) {
	counter := atomic.AddUint64(&m.SendRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestCounter, 1)

	if len(m.SendRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.SendRequest. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockSendRequestInput{p, p1, p2, p3}, "HostNetwork.SendRequest got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, HostNetworkMockSendRequestInput{p, p1, p2, p3}, "HostNetwork.SendRequest got unexpected parameters")
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
		m.t.Fatalf("Unexpected call to HostNetworkMock.SendRequest. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SendRequestFunc(p, p1, p2, p3)
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

type mHostNetworkMockSendRequestToHost struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockSendRequestToHostExpectation
	expectationSeries []*HostNetworkMockSendRequestToHostExpectation
}

type HostNetworkMockSendRequestToHostExpectation struct {
	input  *HostNetworkMockSendRequestToHostInput
	result *HostNetworkMockSendRequestToHostResult
}

type HostNetworkMockSendRequestToHostInput struct {
	p  context.Context
	p1 types.PacketType
	p2 interface{}
	p3 *host.Host
}

type HostNetworkMockSendRequestToHostResult struct {
	r  network.Future
	r1 error
}

//Expect specifies that invocation of HostNetwork.SendRequestToHost is expected from 1 to Infinity times
func (m *mHostNetworkMockSendRequestToHost) Expect(p context.Context, p1 types.PacketType, p2 interface{}, p3 *host.Host) *mHostNetworkMockSendRequestToHost {
	m.mock.SendRequestToHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockSendRequestToHostExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockSendRequestToHostInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of HostNetwork.SendRequestToHost
func (m *mHostNetworkMockSendRequestToHost) Return(r network.Future, r1 error) *HostNetworkMock {
	m.mock.SendRequestToHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockSendRequestToHostExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockSendRequestToHostResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.SendRequestToHost is expected once
func (m *mHostNetworkMockSendRequestToHost) ExpectOnce(p context.Context, p1 types.PacketType, p2 interface{}, p3 *host.Host) *HostNetworkMockSendRequestToHostExpectation {
	m.mock.SendRequestToHostFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockSendRequestToHostExpectation{}
	expectation.input = &HostNetworkMockSendRequestToHostInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockSendRequestToHostExpectation) Return(r network.Future, r1 error) {
	e.result = &HostNetworkMockSendRequestToHostResult{r, r1}
}

//Set uses given function f as a mock of HostNetwork.SendRequestToHost method
func (m *mHostNetworkMockSendRequestToHost) Set(f func(p context.Context, p1 types.PacketType, p2 interface{}, p3 *host.Host) (r network.Future, r1 error)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRequestToHostFunc = f
	return m.mock
}

//SendRequestToHost implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) SendRequestToHost(p context.Context, p1 types.PacketType, p2 interface{}, p3 *host.Host) (r network.Future, r1 error) {
	counter := atomic.AddUint64(&m.SendRequestToHostPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestToHostCounter, 1)

	if len(m.SendRequestToHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendRequestToHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.SendRequestToHost. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendRequestToHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockSendRequestToHostInput{p, p1, p2, p3}, "HostNetwork.SendRequestToHost got unexpected parameters")

		result := m.SendRequestToHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.SendRequestToHost")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestToHostMock.mainExpectation != nil {

		input := m.SendRequestToHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockSendRequestToHostInput{p, p1, p2, p3}, "HostNetwork.SendRequestToHost got unexpected parameters")
		}

		result := m.SendRequestToHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.SendRequestToHost")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRequestToHostFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.SendRequestToHost. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SendRequestToHostFunc(p, p1, p2, p3)
}

//SendRequestToHostMinimockCounter returns a count of HostNetworkMock.SendRequestToHostFunc invocations
func (m *HostNetworkMock) SendRequestToHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestToHostCounter)
}

//SendRequestToHostMinimockPreCounter returns the value of HostNetworkMock.SendRequestToHost invocations
func (m *HostNetworkMock) SendRequestToHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendRequestToHostPreCounter)
}

//SendRequestToHostFinished returns true if mock invocations count is ok
func (m *HostNetworkMock) SendRequestToHostFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendRequestToHostMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendRequestToHostCounter) == uint64(len(m.SendRequestToHostMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendRequestToHostMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendRequestToHostCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendRequestToHostFunc != nil {
		return atomic.LoadUint64(&m.SendRequestToHostCounter) > 0
	}

	return true
}

type mHostNetworkMockStart struct {
	mock              *HostNetworkMock
	mainExpectation   *HostNetworkMockStartExpectation
	expectationSeries []*HostNetworkMockStartExpectation
}

type HostNetworkMockStartExpectation struct {
	input  *HostNetworkMockStartInput
	result *HostNetworkMockStartResult
}

type HostNetworkMockStartInput struct {
	p context.Context
}

type HostNetworkMockStartResult struct {
	r error
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
func (m *mHostNetworkMockStart) Return(r error) *HostNetworkMock {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockStartExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockStartResult{r}
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

func (e *HostNetworkMockStartExpectation) Return(r error) {
	e.result = &HostNetworkMockStartResult{r}
}

//Set uses given function f as a mock of HostNetwork.Start method
func (m *mHostNetworkMockStart) Set(f func(p context.Context) (r error)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StartFunc = f
	return m.mock
}

//Start implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) Start(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if len(m.StartMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StartMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.Start. %v", p)
			return
		}

		input := m.StartMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockStartInput{p}, "HostNetwork.Start got unexpected parameters")

		result := m.StartMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.Start")
			return
		}

		r = result.r

		return
	}

	if m.StartMock.mainExpectation != nil {

		input := m.StartMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockStartInput{p}, "HostNetwork.Start got unexpected parameters")
		}

		result := m.StartMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.Start")
		}

		r = result.r

		return
	}

	if m.StartFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.Start. %v", p)
		return
	}

	return m.StartFunc(p)
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
	input  *HostNetworkMockStopInput
	result *HostNetworkMockStopResult
}

type HostNetworkMockStopInput struct {
	p context.Context
}

type HostNetworkMockStopResult struct {
	r error
}

//Expect specifies that invocation of HostNetwork.Stop is expected from 1 to Infinity times
func (m *mHostNetworkMockStop) Expect(p context.Context) *mHostNetworkMockStop {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockStopExpectation{}
	}
	m.mainExpectation.input = &HostNetworkMockStopInput{p}
	return m
}

//Return specifies results of invocation of HostNetwork.Stop
func (m *mHostNetworkMockStop) Return(r error) *HostNetworkMock {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostNetworkMockStopExpectation{}
	}
	m.mainExpectation.result = &HostNetworkMockStopResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostNetwork.Stop is expected once
func (m *mHostNetworkMockStop) ExpectOnce(p context.Context) *HostNetworkMockStopExpectation {
	m.mock.StopFunc = nil
	m.mainExpectation = nil

	expectation := &HostNetworkMockStopExpectation{}
	expectation.input = &HostNetworkMockStopInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostNetworkMockStopExpectation) Return(r error) {
	e.result = &HostNetworkMockStopResult{r}
}

//Set uses given function f as a mock of HostNetwork.Stop method
func (m *mHostNetworkMockStop) Set(f func(p context.Context) (r error)) *HostNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StopFunc = f
	return m.mock
}

//Stop implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) Stop(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if len(m.StopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostNetworkMock.Stop. %v", p)
			return
		}

		input := m.StopMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostNetworkMockStopInput{p}, "HostNetwork.Stop got unexpected parameters")

		result := m.StopMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.Stop")
			return
		}

		r = result.r

		return
	}

	if m.StopMock.mainExpectation != nil {

		input := m.StopMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostNetworkMockStopInput{p}, "HostNetwork.Stop got unexpected parameters")
		}

		result := m.StopMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostNetworkMock.Stop")
		}

		r = result.r

		return
	}

	if m.StopFunc == nil {
		m.t.Fatalf("Unexpected call to HostNetworkMock.Stop. %v", p)
		return
	}

	return m.StopFunc(p)
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

	if !m.InitFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.Init")
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

	if !m.SendRequestToHostFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.SendRequestToHost")
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

	if !m.InitFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.Init")
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

	if !m.SendRequestToHostFinished() {
		m.t.Fatal("Expected call to HostNetworkMock.SendRequestToHost")
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
		ok = ok && m.InitFinished()
		ok = ok && m.PublicAddressFinished()
		ok = ok && m.RegisterRequestHandlerFinished()
		ok = ok && m.SendRequestFinished()
		ok = ok && m.SendRequestToHostFinished()
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

			if !m.InitFinished() {
				m.t.Error("Expected call to HostNetworkMock.Init")
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

			if !m.SendRequestToHostFinished() {
				m.t.Error("Expected call to HostNetworkMock.SendRequestToHost")
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

	if !m.InitFinished() {
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

	if !m.SendRequestToHostFinished() {
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

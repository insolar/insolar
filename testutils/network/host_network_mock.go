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
	mock             *HostNetworkMock
	mockExpectations *HostNetworkMockBuildResponseParams
}

//HostNetworkMockBuildResponseParams represents input parameters of the HostNetwork.BuildResponse
type HostNetworkMockBuildResponseParams struct {
	p  network.Request
	p1 interface{}
}

//Expect sets up expected params for the HostNetwork.BuildResponse
func (m *mHostNetworkMockBuildResponse) Expect(p network.Request, p1 interface{}) *mHostNetworkMockBuildResponse {
	m.mockExpectations = &HostNetworkMockBuildResponseParams{p, p1}
	return m
}

//Return sets up a mock for HostNetwork.BuildResponse to return Return's arguments
func (m *mHostNetworkMockBuildResponse) Return(r network.Response) *HostNetworkMock {
	m.mock.BuildResponseFunc = func(p network.Request, p1 interface{}) network.Response {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.BuildResponse method
func (m *mHostNetworkMockBuildResponse) Set(f func(p network.Request, p1 interface{}) (r network.Response)) *HostNetworkMock {
	m.mock.BuildResponseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//BuildResponse implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) BuildResponse(p network.Request, p1 interface{}) (r network.Response) {
	atomic.AddUint64(&m.BuildResponsePreCounter, 1)
	defer atomic.AddUint64(&m.BuildResponseCounter, 1)

	if m.BuildResponseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.BuildResponseMock.mockExpectations, HostNetworkMockBuildResponseParams{p, p1},
			"HostNetwork.BuildResponse got unexpected parameters")

		if m.BuildResponseFunc == nil {

			m.t.Fatal("No results are set for the HostNetworkMock.BuildResponse")

			return
		}
	}

	if m.BuildResponseFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.BuildResponse")
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

type mHostNetworkMockGetNodeID struct {
	mock *HostNetworkMock
}

//Return sets up a mock for HostNetwork.GetNodeID to return Return's arguments
func (m *mHostNetworkMockGetNodeID) Return(r core.RecordRef) *HostNetworkMock {
	m.mock.GetNodeIDFunc = func() core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.GetNodeID method
func (m *mHostNetworkMockGetNodeID) Set(f func() (r core.RecordRef)) *HostNetworkMock {
	m.mock.GetNodeIDFunc = f

	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) GetNodeID() (r core.RecordRef) {
	atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if m.GetNodeIDFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.GetNodeID")
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

type mHostNetworkMockNewRequestBuilder struct {
	mock *HostNetworkMock
}

//Return sets up a mock for HostNetwork.NewRequestBuilder to return Return's arguments
func (m *mHostNetworkMockNewRequestBuilder) Return(r network.RequestBuilder) *HostNetworkMock {
	m.mock.NewRequestBuilderFunc = func() network.RequestBuilder {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.NewRequestBuilder method
func (m *mHostNetworkMockNewRequestBuilder) Set(f func() (r network.RequestBuilder)) *HostNetworkMock {
	m.mock.NewRequestBuilderFunc = f

	return m.mock
}

//NewRequestBuilder implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) NewRequestBuilder() (r network.RequestBuilder) {
	atomic.AddUint64(&m.NewRequestBuilderPreCounter, 1)
	defer atomic.AddUint64(&m.NewRequestBuilderCounter, 1)

	if m.NewRequestBuilderFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.NewRequestBuilder")
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

type mHostNetworkMockPublicAddress struct {
	mock *HostNetworkMock
}

//Return sets up a mock for HostNetwork.PublicAddress to return Return's arguments
func (m *mHostNetworkMockPublicAddress) Return(r string) *HostNetworkMock {
	m.mock.PublicAddressFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.PublicAddress method
func (m *mHostNetworkMockPublicAddress) Set(f func() (r string)) *HostNetworkMock {
	m.mock.PublicAddressFunc = f

	return m.mock
}

//PublicAddress implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) PublicAddress() (r string) {
	atomic.AddUint64(&m.PublicAddressPreCounter, 1)
	defer atomic.AddUint64(&m.PublicAddressCounter, 1)

	if m.PublicAddressFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.PublicAddress")
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

type mHostNetworkMockRegisterRequestHandler struct {
	mock             *HostNetworkMock
	mockExpectations *HostNetworkMockRegisterRequestHandlerParams
}

//HostNetworkMockRegisterRequestHandlerParams represents input parameters of the HostNetwork.RegisterRequestHandler
type HostNetworkMockRegisterRequestHandlerParams struct {
	p  types.PacketType
	p1 network.RequestHandler
}

//Expect sets up expected params for the HostNetwork.RegisterRequestHandler
func (m *mHostNetworkMockRegisterRequestHandler) Expect(p types.PacketType, p1 network.RequestHandler) *mHostNetworkMockRegisterRequestHandler {
	m.mockExpectations = &HostNetworkMockRegisterRequestHandlerParams{p, p1}
	return m
}

//Return sets up a mock for HostNetwork.RegisterRequestHandler to return Return's arguments
func (m *mHostNetworkMockRegisterRequestHandler) Return() *HostNetworkMock {
	m.mock.RegisterRequestHandlerFunc = func(p types.PacketType, p1 network.RequestHandler) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.RegisterRequestHandler method
func (m *mHostNetworkMockRegisterRequestHandler) Set(f func(p types.PacketType, p1 network.RequestHandler)) *HostNetworkMock {
	m.mock.RegisterRequestHandlerFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RegisterRequestHandler implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) RegisterRequestHandler(p types.PacketType, p1 network.RequestHandler) {
	atomic.AddUint64(&m.RegisterRequestHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestHandlerCounter, 1)

	if m.RegisterRequestHandlerMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RegisterRequestHandlerMock.mockExpectations, HostNetworkMockRegisterRequestHandlerParams{p, p1},
			"HostNetwork.RegisterRequestHandler got unexpected parameters")

		if m.RegisterRequestHandlerFunc == nil {

			m.t.Fatal("No results are set for the HostNetworkMock.RegisterRequestHandler")

			return
		}
	}

	if m.RegisterRequestHandlerFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.RegisterRequestHandler")
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

type mHostNetworkMockSendRequest struct {
	mock             *HostNetworkMock
	mockExpectations *HostNetworkMockSendRequestParams
}

//HostNetworkMockSendRequestParams represents input parameters of the HostNetwork.SendRequest
type HostNetworkMockSendRequestParams struct {
	p  network.Request
	p1 core.RecordRef
}

//Expect sets up expected params for the HostNetwork.SendRequest
func (m *mHostNetworkMockSendRequest) Expect(p network.Request, p1 core.RecordRef) *mHostNetworkMockSendRequest {
	m.mockExpectations = &HostNetworkMockSendRequestParams{p, p1}
	return m
}

//Return sets up a mock for HostNetwork.SendRequest to return Return's arguments
func (m *mHostNetworkMockSendRequest) Return(r network.Future, r1 error) *HostNetworkMock {
	m.mock.SendRequestFunc = func(p network.Request, p1 core.RecordRef) (network.Future, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.SendRequest method
func (m *mHostNetworkMockSendRequest) Set(f func(p network.Request, p1 core.RecordRef) (r network.Future, r1 error)) *HostNetworkMock {
	m.mock.SendRequestFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) SendRequest(p network.Request, p1 core.RecordRef) (r network.Future, r1 error) {
	atomic.AddUint64(&m.SendRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestCounter, 1)

	if m.SendRequestMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendRequestMock.mockExpectations, HostNetworkMockSendRequestParams{p, p1},
			"HostNetwork.SendRequest got unexpected parameters")

		if m.SendRequestFunc == nil {

			m.t.Fatal("No results are set for the HostNetworkMock.SendRequest")

			return
		}
	}

	if m.SendRequestFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.SendRequest")
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

type mHostNetworkMockStart struct {
	mock             *HostNetworkMock
	mockExpectations *HostNetworkMockStartParams
}

//HostNetworkMockStartParams represents input parameters of the HostNetwork.Start
type HostNetworkMockStartParams struct {
	p context.Context
}

//Expect sets up expected params for the HostNetwork.Start
func (m *mHostNetworkMockStart) Expect(p context.Context) *mHostNetworkMockStart {
	m.mockExpectations = &HostNetworkMockStartParams{p}
	return m
}

//Return sets up a mock for HostNetwork.Start to return Return's arguments
func (m *mHostNetworkMockStart) Return() *HostNetworkMock {
	m.mock.StartFunc = func(p context.Context) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.Start method
func (m *mHostNetworkMockStart) Set(f func(p context.Context)) *HostNetworkMock {
	m.mock.StartFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Start implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) Start(p context.Context) {
	atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if m.StartMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.StartMock.mockExpectations, HostNetworkMockStartParams{p},
			"HostNetwork.Start got unexpected parameters")

		if m.StartFunc == nil {

			m.t.Fatal("No results are set for the HostNetworkMock.Start")

			return
		}
	}

	if m.StartFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.Start")
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

type mHostNetworkMockStop struct {
	mock *HostNetworkMock
}

//Return sets up a mock for HostNetwork.Stop to return Return's arguments
func (m *mHostNetworkMockStop) Return() *HostNetworkMock {
	m.mock.StopFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of HostNetwork.Stop method
func (m *mHostNetworkMockStop) Set(f func()) *HostNetworkMock {
	m.mock.StopFunc = f

	return m.mock
}

//Stop implements github.com/insolar/insolar/network.HostNetwork interface
func (m *HostNetworkMock) Stop() {
	atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if m.StopFunc == nil {
		m.t.Fatal("Unexpected call to HostNetworkMock.Stop")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HostNetworkMock) ValidateCallCounters() {

	if m.BuildResponseFunc != nil && atomic.LoadUint64(&m.BuildResponseCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.BuildResponse")
	}

	if m.GetNodeIDFunc != nil && atomic.LoadUint64(&m.GetNodeIDCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.GetNodeID")
	}

	if m.NewRequestBuilderFunc != nil && atomic.LoadUint64(&m.NewRequestBuilderCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.NewRequestBuilder")
	}

	if m.PublicAddressFunc != nil && atomic.LoadUint64(&m.PublicAddressCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.PublicAddress")
	}

	if m.RegisterRequestHandlerFunc != nil && atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.RegisterRequestHandler")
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.SendRequest")
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.Start")
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
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

	if m.BuildResponseFunc != nil && atomic.LoadUint64(&m.BuildResponseCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.BuildResponse")
	}

	if m.GetNodeIDFunc != nil && atomic.LoadUint64(&m.GetNodeIDCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.GetNodeID")
	}

	if m.NewRequestBuilderFunc != nil && atomic.LoadUint64(&m.NewRequestBuilderCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.NewRequestBuilder")
	}

	if m.PublicAddressFunc != nil && atomic.LoadUint64(&m.PublicAddressCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.PublicAddress")
	}

	if m.RegisterRequestHandlerFunc != nil && atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.RegisterRequestHandler")
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.SendRequest")
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		m.t.Fatal("Expected call to HostNetworkMock.Start")
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
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
		ok = ok && (m.BuildResponseFunc == nil || atomic.LoadUint64(&m.BuildResponseCounter) > 0)
		ok = ok && (m.GetNodeIDFunc == nil || atomic.LoadUint64(&m.GetNodeIDCounter) > 0)
		ok = ok && (m.NewRequestBuilderFunc == nil || atomic.LoadUint64(&m.NewRequestBuilderCounter) > 0)
		ok = ok && (m.PublicAddressFunc == nil || atomic.LoadUint64(&m.PublicAddressCounter) > 0)
		ok = ok && (m.RegisterRequestHandlerFunc == nil || atomic.LoadUint64(&m.RegisterRequestHandlerCounter) > 0)
		ok = ok && (m.SendRequestFunc == nil || atomic.LoadUint64(&m.SendRequestCounter) > 0)
		ok = ok && (m.StartFunc == nil || atomic.LoadUint64(&m.StartCounter) > 0)
		ok = ok && (m.StopFunc == nil || atomic.LoadUint64(&m.StopCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.BuildResponseFunc != nil && atomic.LoadUint64(&m.BuildResponseCounter) == 0 {
				m.t.Error("Expected call to HostNetworkMock.BuildResponse")
			}

			if m.GetNodeIDFunc != nil && atomic.LoadUint64(&m.GetNodeIDCounter) == 0 {
				m.t.Error("Expected call to HostNetworkMock.GetNodeID")
			}

			if m.NewRequestBuilderFunc != nil && atomic.LoadUint64(&m.NewRequestBuilderCounter) == 0 {
				m.t.Error("Expected call to HostNetworkMock.NewRequestBuilder")
			}

			if m.PublicAddressFunc != nil && atomic.LoadUint64(&m.PublicAddressCounter) == 0 {
				m.t.Error("Expected call to HostNetworkMock.PublicAddress")
			}

			if m.RegisterRequestHandlerFunc != nil && atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == 0 {
				m.t.Error("Expected call to HostNetworkMock.RegisterRequestHandler")
			}

			if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
				m.t.Error("Expected call to HostNetworkMock.SendRequest")
			}

			if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
				m.t.Error("Expected call to HostNetworkMock.Start")
			}

			if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
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

	if m.BuildResponseFunc != nil && atomic.LoadUint64(&m.BuildResponseCounter) == 0 {
		return false
	}

	if m.GetNodeIDFunc != nil && atomic.LoadUint64(&m.GetNodeIDCounter) == 0 {
		return false
	}

	if m.NewRequestBuilderFunc != nil && atomic.LoadUint64(&m.NewRequestBuilderCounter) == 0 {
		return false
	}

	if m.PublicAddressFunc != nil && atomic.LoadUint64(&m.PublicAddressCounter) == 0 {
		return false
	}

	if m.RegisterRequestHandlerFunc != nil && atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == 0 {
		return false
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
		return false
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		return false
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
		return false
	}

	return true
}

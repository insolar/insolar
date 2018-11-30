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
	mock *ConsensusNetworkMock
}

//Return sets up a mock for ConsensusNetwork.GetNodeID to return Return's arguments
func (m *mConsensusNetworkMockGetNodeID) Return(r core.RecordRef) *ConsensusNetworkMock {
	m.mock.GetNodeIDFunc = func() core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ConsensusNetwork.GetNodeID method
func (m *mConsensusNetworkMockGetNodeID) Set(f func() (r core.RecordRef)) *ConsensusNetworkMock {
	m.mock.GetNodeIDFunc = f

	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) GetNodeID() (r core.RecordRef) {
	atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if m.GetNodeIDFunc == nil {
		m.t.Fatal("Unexpected call to ConsensusNetworkMock.GetNodeID")
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

type mConsensusNetworkMockNewRequestBuilder struct {
	mock *ConsensusNetworkMock
}

//Return sets up a mock for ConsensusNetwork.NewRequestBuilder to return Return's arguments
func (m *mConsensusNetworkMockNewRequestBuilder) Return(r network.RequestBuilder) *ConsensusNetworkMock {
	m.mock.NewRequestBuilderFunc = func() network.RequestBuilder {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ConsensusNetwork.NewRequestBuilder method
func (m *mConsensusNetworkMockNewRequestBuilder) Set(f func() (r network.RequestBuilder)) *ConsensusNetworkMock {
	m.mock.NewRequestBuilderFunc = f

	return m.mock
}

//NewRequestBuilder implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) NewRequestBuilder() (r network.RequestBuilder) {
	atomic.AddUint64(&m.NewRequestBuilderPreCounter, 1)
	defer atomic.AddUint64(&m.NewRequestBuilderCounter, 1)

	if m.NewRequestBuilderFunc == nil {
		m.t.Fatal("Unexpected call to ConsensusNetworkMock.NewRequestBuilder")
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

type mConsensusNetworkMockPublicAddress struct {
	mock *ConsensusNetworkMock
}

//Return sets up a mock for ConsensusNetwork.PublicAddress to return Return's arguments
func (m *mConsensusNetworkMockPublicAddress) Return(r string) *ConsensusNetworkMock {
	m.mock.PublicAddressFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ConsensusNetwork.PublicAddress method
func (m *mConsensusNetworkMockPublicAddress) Set(f func() (r string)) *ConsensusNetworkMock {
	m.mock.PublicAddressFunc = f

	return m.mock
}

//PublicAddress implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) PublicAddress() (r string) {
	atomic.AddUint64(&m.PublicAddressPreCounter, 1)
	defer atomic.AddUint64(&m.PublicAddressCounter, 1)

	if m.PublicAddressFunc == nil {
		m.t.Fatal("Unexpected call to ConsensusNetworkMock.PublicAddress")
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

type mConsensusNetworkMockRegisterRequestHandler struct {
	mock             *ConsensusNetworkMock
	mockExpectations *ConsensusNetworkMockRegisterRequestHandlerParams
}

//ConsensusNetworkMockRegisterRequestHandlerParams represents input parameters of the ConsensusNetwork.RegisterRequestHandler
type ConsensusNetworkMockRegisterRequestHandlerParams struct {
	p  types.PacketType
	p1 network.ConsensusRequestHandler
}

//Expect sets up expected params for the ConsensusNetwork.RegisterRequestHandler
func (m *mConsensusNetworkMockRegisterRequestHandler) Expect(p types.PacketType, p1 network.ConsensusRequestHandler) *mConsensusNetworkMockRegisterRequestHandler {
	m.mockExpectations = &ConsensusNetworkMockRegisterRequestHandlerParams{p, p1}
	return m
}

//Return sets up a mock for ConsensusNetwork.RegisterRequestHandler to return Return's arguments
func (m *mConsensusNetworkMockRegisterRequestHandler) Return() *ConsensusNetworkMock {
	m.mock.RegisterRequestHandlerFunc = func(p types.PacketType, p1 network.ConsensusRequestHandler) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of ConsensusNetwork.RegisterRequestHandler method
func (m *mConsensusNetworkMockRegisterRequestHandler) Set(f func(p types.PacketType, p1 network.ConsensusRequestHandler)) *ConsensusNetworkMock {
	m.mock.RegisterRequestHandlerFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RegisterRequestHandler implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) RegisterRequestHandler(p types.PacketType, p1 network.ConsensusRequestHandler) {
	atomic.AddUint64(&m.RegisterRequestHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestHandlerCounter, 1)

	if m.RegisterRequestHandlerMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RegisterRequestHandlerMock.mockExpectations, ConsensusNetworkMockRegisterRequestHandlerParams{p, p1},
			"ConsensusNetwork.RegisterRequestHandler got unexpected parameters")

		if m.RegisterRequestHandlerFunc == nil {

			m.t.Fatal("No results are set for the ConsensusNetworkMock.RegisterRequestHandler")

			return
		}
	}

	if m.RegisterRequestHandlerFunc == nil {
		m.t.Fatal("Unexpected call to ConsensusNetworkMock.RegisterRequestHandler")
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

type mConsensusNetworkMockSendRequest struct {
	mock             *ConsensusNetworkMock
	mockExpectations *ConsensusNetworkMockSendRequestParams
}

//ConsensusNetworkMockSendRequestParams represents input parameters of the ConsensusNetwork.SendRequest
type ConsensusNetworkMockSendRequestParams struct {
	p  network.Request
	p1 core.RecordRef
}

//Expect sets up expected params for the ConsensusNetwork.SendRequest
func (m *mConsensusNetworkMockSendRequest) Expect(p network.Request, p1 core.RecordRef) *mConsensusNetworkMockSendRequest {
	m.mockExpectations = &ConsensusNetworkMockSendRequestParams{p, p1}
	return m
}

//Return sets up a mock for ConsensusNetwork.SendRequest to return Return's arguments
func (m *mConsensusNetworkMockSendRequest) Return(r error) *ConsensusNetworkMock {
	m.mock.SendRequestFunc = func(p network.Request, p1 core.RecordRef) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ConsensusNetwork.SendRequest method
func (m *mConsensusNetworkMockSendRequest) Set(f func(p network.Request, p1 core.RecordRef) (r error)) *ConsensusNetworkMock {
	m.mock.SendRequestFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SendRequest implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) SendRequest(p network.Request, p1 core.RecordRef) (r error) {
	atomic.AddUint64(&m.SendRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SendRequestCounter, 1)

	if m.SendRequestMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendRequestMock.mockExpectations, ConsensusNetworkMockSendRequestParams{p, p1},
			"ConsensusNetwork.SendRequest got unexpected parameters")

		if m.SendRequestFunc == nil {

			m.t.Fatal("No results are set for the ConsensusNetworkMock.SendRequest")

			return
		}
	}

	if m.SendRequestFunc == nil {
		m.t.Fatal("Unexpected call to ConsensusNetworkMock.SendRequest")
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

type mConsensusNetworkMockStart struct {
	mock             *ConsensusNetworkMock
	mockExpectations *ConsensusNetworkMockStartParams
}

//ConsensusNetworkMockStartParams represents input parameters of the ConsensusNetwork.Start
type ConsensusNetworkMockStartParams struct {
	p context.Context
}

//Expect sets up expected params for the ConsensusNetwork.Start
func (m *mConsensusNetworkMockStart) Expect(p context.Context) *mConsensusNetworkMockStart {
	m.mockExpectations = &ConsensusNetworkMockStartParams{p}
	return m
}

//Return sets up a mock for ConsensusNetwork.Start to return Return's arguments
func (m *mConsensusNetworkMockStart) Return() *ConsensusNetworkMock {
	m.mock.StartFunc = func(p context.Context) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of ConsensusNetwork.Start method
func (m *mConsensusNetworkMockStart) Set(f func(p context.Context)) *ConsensusNetworkMock {
	m.mock.StartFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Start implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) Start(p context.Context) {
	atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if m.StartMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.StartMock.mockExpectations, ConsensusNetworkMockStartParams{p},
			"ConsensusNetwork.Start got unexpected parameters")

		if m.StartFunc == nil {

			m.t.Fatal("No results are set for the ConsensusNetworkMock.Start")

			return
		}
	}

	if m.StartFunc == nil {
		m.t.Fatal("Unexpected call to ConsensusNetworkMock.Start")
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

type mConsensusNetworkMockStop struct {
	mock *ConsensusNetworkMock
}

//Return sets up a mock for ConsensusNetwork.Stop to return Return's arguments
func (m *mConsensusNetworkMockStop) Return() *ConsensusNetworkMock {
	m.mock.StopFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of ConsensusNetwork.Stop method
func (m *mConsensusNetworkMockStop) Set(f func()) *ConsensusNetworkMock {
	m.mock.StopFunc = f

	return m.mock
}

//Stop implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) Stop() {
	atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if m.StopFunc == nil {
		m.t.Fatal("Unexpected call to ConsensusNetworkMock.Stop")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ConsensusNetworkMock) ValidateCallCounters() {

	if m.GetNodeIDFunc != nil && atomic.LoadUint64(&m.GetNodeIDCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.GetNodeID")
	}

	if m.NewRequestBuilderFunc != nil && atomic.LoadUint64(&m.NewRequestBuilderCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.NewRequestBuilder")
	}

	if m.PublicAddressFunc != nil && atomic.LoadUint64(&m.PublicAddressCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.PublicAddress")
	}

	if m.RegisterRequestHandlerFunc != nil && atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.RegisterRequestHandler")
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.SendRequest")
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.Start")
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
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

	if m.GetNodeIDFunc != nil && atomic.LoadUint64(&m.GetNodeIDCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.GetNodeID")
	}

	if m.NewRequestBuilderFunc != nil && atomic.LoadUint64(&m.NewRequestBuilderCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.NewRequestBuilder")
	}

	if m.PublicAddressFunc != nil && atomic.LoadUint64(&m.PublicAddressCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.PublicAddress")
	}

	if m.RegisterRequestHandlerFunc != nil && atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.RegisterRequestHandler")
	}

	if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.SendRequest")
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		m.t.Fatal("Expected call to ConsensusNetworkMock.Start")
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
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

			if m.GetNodeIDFunc != nil && atomic.LoadUint64(&m.GetNodeIDCounter) == 0 {
				m.t.Error("Expected call to ConsensusNetworkMock.GetNodeID")
			}

			if m.NewRequestBuilderFunc != nil && atomic.LoadUint64(&m.NewRequestBuilderCounter) == 0 {
				m.t.Error("Expected call to ConsensusNetworkMock.NewRequestBuilder")
			}

			if m.PublicAddressFunc != nil && atomic.LoadUint64(&m.PublicAddressCounter) == 0 {
				m.t.Error("Expected call to ConsensusNetworkMock.PublicAddress")
			}

			if m.RegisterRequestHandlerFunc != nil && atomic.LoadUint64(&m.RegisterRequestHandlerCounter) == 0 {
				m.t.Error("Expected call to ConsensusNetworkMock.RegisterRequestHandler")
			}

			if m.SendRequestFunc != nil && atomic.LoadUint64(&m.SendRequestCounter) == 0 {
				m.t.Error("Expected call to ConsensusNetworkMock.SendRequest")
			}

			if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
				m.t.Error("Expected call to ConsensusNetworkMock.Start")
			}

			if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
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

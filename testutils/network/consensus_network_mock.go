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
	packets "github.com/insolar/insolar/consensus/packets"
	core "github.com/insolar/insolar/core"
	network "github.com/insolar/insolar/network"

	testify_assert "github.com/stretchr/testify/assert"
)

//ConsensusNetworkMock implements github.com/insolar/insolar/network.ConsensusNetwork
type ConsensusNetworkMock struct {
	t minimock.Tester

	GetNodeIDFunc       func() (r core.RecordRef)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mConsensusNetworkMockGetNodeID

	PublicAddressFunc       func() (r string)
	PublicAddressCounter    uint64
	PublicAddressPreCounter uint64
	PublicAddressMock       mConsensusNetworkMockPublicAddress

	RegisterPacketHandlerFunc       func(p packets.PacketType, p1 network.ConsensusPacketHandler)
	RegisterPacketHandlerCounter    uint64
	RegisterPacketHandlerPreCounter uint64
	RegisterPacketHandlerMock       mConsensusNetworkMockRegisterPacketHandler

	SignAndSendPacketFunc       func(p packets.ConsensusPacket, p1 core.RecordRef, p2 core.CryptographyService) (r error)
	SignAndSendPacketCounter    uint64
	SignAndSendPacketPreCounter uint64
	SignAndSendPacketMock       mConsensusNetworkMockSignAndSendPacket

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
	m.PublicAddressMock = mConsensusNetworkMockPublicAddress{mock: m}
	m.RegisterPacketHandlerMock = mConsensusNetworkMockRegisterPacketHandler{mock: m}
	m.SignAndSendPacketMock = mConsensusNetworkMockSignAndSendPacket{mock: m}
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

type mConsensusNetworkMockRegisterPacketHandler struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockRegisterPacketHandlerExpectation
	expectationSeries []*ConsensusNetworkMockRegisterPacketHandlerExpectation
}

type ConsensusNetworkMockRegisterPacketHandlerExpectation struct {
	input *ConsensusNetworkMockRegisterPacketHandlerInput
}

type ConsensusNetworkMockRegisterPacketHandlerInput struct {
	p  packets.PacketType
	p1 network.ConsensusPacketHandler
}

//Expect specifies that invocation of ConsensusNetwork.RegisterPacketHandler is expected from 1 to Infinity times
func (m *mConsensusNetworkMockRegisterPacketHandler) Expect(p packets.PacketType, p1 network.ConsensusPacketHandler) *mConsensusNetworkMockRegisterPacketHandler {
	m.mock.RegisterPacketHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockRegisterPacketHandlerExpectation{}
	}
	m.mainExpectation.input = &ConsensusNetworkMockRegisterPacketHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of ConsensusNetwork.RegisterPacketHandler
func (m *mConsensusNetworkMockRegisterPacketHandler) Return() *ConsensusNetworkMock {
	m.mock.RegisterPacketHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockRegisterPacketHandlerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.RegisterPacketHandler is expected once
func (m *mConsensusNetworkMockRegisterPacketHandler) ExpectOnce(p packets.PacketType, p1 network.ConsensusPacketHandler) *ConsensusNetworkMockRegisterPacketHandlerExpectation {
	m.mock.RegisterPacketHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockRegisterPacketHandlerExpectation{}
	expectation.input = &ConsensusNetworkMockRegisterPacketHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ConsensusNetwork.RegisterPacketHandler method
func (m *mConsensusNetworkMockRegisterPacketHandler) Set(f func(p packets.PacketType, p1 network.ConsensusPacketHandler)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterPacketHandlerFunc = f
	return m.mock
}

//RegisterPacketHandler implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) RegisterPacketHandler(p packets.PacketType, p1 network.ConsensusPacketHandler) {
	counter := atomic.AddUint64(&m.RegisterPacketHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterPacketHandlerCounter, 1)

	if len(m.RegisterPacketHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterPacketHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.RegisterPacketHandler. %v %v", p, p1)
			return
		}

		input := m.RegisterPacketHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConsensusNetworkMockRegisterPacketHandlerInput{p, p1}, "ConsensusNetwork.RegisterPacketHandler got unexpected parameters")

		return
	}

	if m.RegisterPacketHandlerMock.mainExpectation != nil {

		input := m.RegisterPacketHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConsensusNetworkMockRegisterPacketHandlerInput{p, p1}, "ConsensusNetwork.RegisterPacketHandler got unexpected parameters")
		}

		return
	}

	if m.RegisterPacketHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.RegisterPacketHandler. %v %v", p, p1)
		return
	}

	m.RegisterPacketHandlerFunc(p, p1)
}

//RegisterPacketHandlerMinimockCounter returns a count of ConsensusNetworkMock.RegisterPacketHandlerFunc invocations
func (m *ConsensusNetworkMock) RegisterPacketHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterPacketHandlerCounter)
}

//RegisterPacketHandlerMinimockPreCounter returns the value of ConsensusNetworkMock.RegisterPacketHandler invocations
func (m *ConsensusNetworkMock) RegisterPacketHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterPacketHandlerPreCounter)
}

//RegisterPacketHandlerFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) RegisterPacketHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterPacketHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterPacketHandlerCounter) == uint64(len(m.RegisterPacketHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterPacketHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterPacketHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterPacketHandlerFunc != nil {
		return atomic.LoadUint64(&m.RegisterPacketHandlerCounter) > 0
	}

	return true
}

type mConsensusNetworkMockSignAndSendPacket struct {
	mock              *ConsensusNetworkMock
	mainExpectation   *ConsensusNetworkMockSignAndSendPacketExpectation
	expectationSeries []*ConsensusNetworkMockSignAndSendPacketExpectation
}

type ConsensusNetworkMockSignAndSendPacketExpectation struct {
	input  *ConsensusNetworkMockSignAndSendPacketInput
	result *ConsensusNetworkMockSignAndSendPacketResult
}

type ConsensusNetworkMockSignAndSendPacketInput struct {
	p  packets.ConsensusPacket
	p1 core.RecordRef
	p2 core.CryptographyService
}

type ConsensusNetworkMockSignAndSendPacketResult struct {
	r error
}

//Expect specifies that invocation of ConsensusNetwork.SignAndSendPacket is expected from 1 to Infinity times
func (m *mConsensusNetworkMockSignAndSendPacket) Expect(p packets.ConsensusPacket, p1 core.RecordRef, p2 core.CryptographyService) *mConsensusNetworkMockSignAndSendPacket {
	m.mock.SignAndSendPacketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockSignAndSendPacketExpectation{}
	}
	m.mainExpectation.input = &ConsensusNetworkMockSignAndSendPacketInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ConsensusNetwork.SignAndSendPacket
func (m *mConsensusNetworkMockSignAndSendPacket) Return(r error) *ConsensusNetworkMock {
	m.mock.SignAndSendPacketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConsensusNetworkMockSignAndSendPacketExpectation{}
	}
	m.mainExpectation.result = &ConsensusNetworkMockSignAndSendPacketResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ConsensusNetwork.SignAndSendPacket is expected once
func (m *mConsensusNetworkMockSignAndSendPacket) ExpectOnce(p packets.ConsensusPacket, p1 core.RecordRef, p2 core.CryptographyService) *ConsensusNetworkMockSignAndSendPacketExpectation {
	m.mock.SignAndSendPacketFunc = nil
	m.mainExpectation = nil

	expectation := &ConsensusNetworkMockSignAndSendPacketExpectation{}
	expectation.input = &ConsensusNetworkMockSignAndSendPacketInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConsensusNetworkMockSignAndSendPacketExpectation) Return(r error) {
	e.result = &ConsensusNetworkMockSignAndSendPacketResult{r}
}

//Set uses given function f as a mock of ConsensusNetwork.SignAndSendPacket method
func (m *mConsensusNetworkMockSignAndSendPacket) Set(f func(p packets.ConsensusPacket, p1 core.RecordRef, p2 core.CryptographyService) (r error)) *ConsensusNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignAndSendPacketFunc = f
	return m.mock
}

//SignAndSendPacket implements github.com/insolar/insolar/network.ConsensusNetwork interface
func (m *ConsensusNetworkMock) SignAndSendPacket(p packets.ConsensusPacket, p1 core.RecordRef, p2 core.CryptographyService) (r error) {
	counter := atomic.AddUint64(&m.SignAndSendPacketPreCounter, 1)
	defer atomic.AddUint64(&m.SignAndSendPacketCounter, 1)

	if len(m.SignAndSendPacketMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignAndSendPacketMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConsensusNetworkMock.SignAndSendPacket. %v %v %v", p, p1, p2)
			return
		}

		input := m.SignAndSendPacketMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConsensusNetworkMockSignAndSendPacketInput{p, p1, p2}, "ConsensusNetwork.SignAndSendPacket got unexpected parameters")

		result := m.SignAndSendPacketMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.SignAndSendPacket")
			return
		}

		r = result.r

		return
	}

	if m.SignAndSendPacketMock.mainExpectation != nil {

		input := m.SignAndSendPacketMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConsensusNetworkMockSignAndSendPacketInput{p, p1, p2}, "ConsensusNetwork.SignAndSendPacket got unexpected parameters")
		}

		result := m.SignAndSendPacketMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConsensusNetworkMock.SignAndSendPacket")
		}

		r = result.r

		return
	}

	if m.SignAndSendPacketFunc == nil {
		m.t.Fatalf("Unexpected call to ConsensusNetworkMock.SignAndSendPacket. %v %v %v", p, p1, p2)
		return
	}

	return m.SignAndSendPacketFunc(p, p1, p2)
}

//SignAndSendPacketMinimockCounter returns a count of ConsensusNetworkMock.SignAndSendPacketFunc invocations
func (m *ConsensusNetworkMock) SignAndSendPacketMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignAndSendPacketCounter)
}

//SignAndSendPacketMinimockPreCounter returns the value of ConsensusNetworkMock.SignAndSendPacket invocations
func (m *ConsensusNetworkMock) SignAndSendPacketMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignAndSendPacketPreCounter)
}

//SignAndSendPacketFinished returns true if mock invocations count is ok
func (m *ConsensusNetworkMock) SignAndSendPacketFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SignAndSendPacketMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SignAndSendPacketCounter) == uint64(len(m.SignAndSendPacketMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SignAndSendPacketMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SignAndSendPacketCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SignAndSendPacketFunc != nil {
		return atomic.LoadUint64(&m.SignAndSendPacketCounter) > 0
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

	if !m.PublicAddressFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.PublicAddress")
	}

	if !m.RegisterPacketHandlerFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.RegisterPacketHandler")
	}

	if !m.SignAndSendPacketFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.SignAndSendPacket")
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

	if !m.PublicAddressFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.PublicAddress")
	}

	if !m.RegisterPacketHandlerFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.RegisterPacketHandler")
	}

	if !m.SignAndSendPacketFinished() {
		m.t.Fatal("Expected call to ConsensusNetworkMock.SignAndSendPacket")
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
		ok = ok && m.PublicAddressFinished()
		ok = ok && m.RegisterPacketHandlerFinished()
		ok = ok && m.SignAndSendPacketFinished()
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

			if !m.PublicAddressFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.PublicAddress")
			}

			if !m.RegisterPacketHandlerFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.RegisterPacketHandler")
			}

			if !m.SignAndSendPacketFinished() {
				m.t.Error("Expected call to ConsensusNetworkMock.SignAndSendPacket")
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

	if !m.PublicAddressFinished() {
		return false
	}

	if !m.RegisterPacketHandlerFinished() {
		return false
	}

	if !m.SignAndSendPacketFinished() {
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

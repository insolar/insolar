package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RPCController" can be found in github.com/insolar/insolar/network/controller
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//RPCControllerMock implements github.com/insolar/insolar/network/controller.RPCController
type RPCControllerMock struct {
	t minimock.Tester

	IAmRPCControllerFunc       func()
	IAmRPCControllerCounter    uint64
	IAmRPCControllerPreCounter uint64
	IAmRPCControllerMock       mRPCControllerMockIAmRPCController

	InitFunc       func(p context.Context) (r error)
	InitCounter    uint64
	InitPreCounter uint64
	InitMock       mRPCControllerMockInit

	RemoteProcedureRegisterFunc       func(p string, p1 insolar.RemoteProcedure)
	RemoteProcedureRegisterCounter    uint64
	RemoteProcedureRegisterPreCounter uint64
	RemoteProcedureRegisterMock       mRPCControllerMockRemoteProcedureRegister

	SendBytesFunc       func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error)
	SendBytesCounter    uint64
	SendBytesPreCounter uint64
	SendBytesMock       mRPCControllerMockSendBytes

	SendCascadeMessageFunc       func(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error)
	SendCascadeMessageCounter    uint64
	SendCascadeMessagePreCounter uint64
	SendCascadeMessageMock       mRPCControllerMockSendCascadeMessage

	SendMessageFunc       func(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error)
	SendMessageCounter    uint64
	SendMessagePreCounter uint64
	SendMessageMock       mRPCControllerMockSendMessage
}

//NewRPCControllerMock returns a mock for github.com/insolar/insolar/network/controller.RPCController
func NewRPCControllerMock(t minimock.Tester) *RPCControllerMock {
	m := &RPCControllerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IAmRPCControllerMock = mRPCControllerMockIAmRPCController{mock: m}
	m.InitMock = mRPCControllerMockInit{mock: m}
	m.RemoteProcedureRegisterMock = mRPCControllerMockRemoteProcedureRegister{mock: m}
	m.SendBytesMock = mRPCControllerMockSendBytes{mock: m}
	m.SendCascadeMessageMock = mRPCControllerMockSendCascadeMessage{mock: m}
	m.SendMessageMock = mRPCControllerMockSendMessage{mock: m}

	return m
}

type mRPCControllerMockIAmRPCController struct {
	mock              *RPCControllerMock
	mainExpectation   *RPCControllerMockIAmRPCControllerExpectation
	expectationSeries []*RPCControllerMockIAmRPCControllerExpectation
}

type RPCControllerMockIAmRPCControllerExpectation struct {
}

//Expect specifies that invocation of RPCController.IAmRPCController is expected from 1 to Infinity times
func (m *mRPCControllerMockIAmRPCController) Expect() *mRPCControllerMockIAmRPCController {
	m.mock.IAmRPCControllerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockIAmRPCControllerExpectation{}
	}

	return m
}

//Return specifies results of invocation of RPCController.IAmRPCController
func (m *mRPCControllerMockIAmRPCController) Return() *RPCControllerMock {
	m.mock.IAmRPCControllerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockIAmRPCControllerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RPCController.IAmRPCController is expected once
func (m *mRPCControllerMockIAmRPCController) ExpectOnce() *RPCControllerMockIAmRPCControllerExpectation {
	m.mock.IAmRPCControllerFunc = nil
	m.mainExpectation = nil

	expectation := &RPCControllerMockIAmRPCControllerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RPCController.IAmRPCController method
func (m *mRPCControllerMockIAmRPCController) Set(f func()) *RPCControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IAmRPCControllerFunc = f
	return m.mock
}

//IAmRPCController implements github.com/insolar/insolar/network/controller.RPCController interface
func (m *RPCControllerMock) IAmRPCController() {
	counter := atomic.AddUint64(&m.IAmRPCControllerPreCounter, 1)
	defer atomic.AddUint64(&m.IAmRPCControllerCounter, 1)

	if len(m.IAmRPCControllerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IAmRPCControllerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RPCControllerMock.IAmRPCController.")
			return
		}

		return
	}

	if m.IAmRPCControllerMock.mainExpectation != nil {

		return
	}

	if m.IAmRPCControllerFunc == nil {
		m.t.Fatalf("Unexpected call to RPCControllerMock.IAmRPCController.")
		return
	}

	m.IAmRPCControllerFunc()
}

//IAmRPCControllerMinimockCounter returns a count of RPCControllerMock.IAmRPCControllerFunc invocations
func (m *RPCControllerMock) IAmRPCControllerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IAmRPCControllerCounter)
}

//IAmRPCControllerMinimockPreCounter returns the value of RPCControllerMock.IAmRPCController invocations
func (m *RPCControllerMock) IAmRPCControllerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IAmRPCControllerPreCounter)
}

//IAmRPCControllerFinished returns true if mock invocations count is ok
func (m *RPCControllerMock) IAmRPCControllerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IAmRPCControllerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IAmRPCControllerCounter) == uint64(len(m.IAmRPCControllerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IAmRPCControllerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IAmRPCControllerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IAmRPCControllerFunc != nil {
		return atomic.LoadUint64(&m.IAmRPCControllerCounter) > 0
	}

	return true
}

type mRPCControllerMockInit struct {
	mock              *RPCControllerMock
	mainExpectation   *RPCControllerMockInitExpectation
	expectationSeries []*RPCControllerMockInitExpectation
}

type RPCControllerMockInitExpectation struct {
	input  *RPCControllerMockInitInput
	result *RPCControllerMockInitResult
}

type RPCControllerMockInitInput struct {
	p context.Context
}

type RPCControllerMockInitResult struct {
	r error
}

//Expect specifies that invocation of RPCController.Init is expected from 1 to Infinity times
func (m *mRPCControllerMockInit) Expect(p context.Context) *mRPCControllerMockInit {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockInitExpectation{}
	}
	m.mainExpectation.input = &RPCControllerMockInitInput{p}
	return m
}

//Return specifies results of invocation of RPCController.Init
func (m *mRPCControllerMockInit) Return(r error) *RPCControllerMock {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockInitExpectation{}
	}
	m.mainExpectation.result = &RPCControllerMockInitResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RPCController.Init is expected once
func (m *mRPCControllerMockInit) ExpectOnce(p context.Context) *RPCControllerMockInitExpectation {
	m.mock.InitFunc = nil
	m.mainExpectation = nil

	expectation := &RPCControllerMockInitExpectation{}
	expectation.input = &RPCControllerMockInitInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RPCControllerMockInitExpectation) Return(r error) {
	e.result = &RPCControllerMockInitResult{r}
}

//Set uses given function f as a mock of RPCController.Init method
func (m *mRPCControllerMockInit) Set(f func(p context.Context) (r error)) *RPCControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InitFunc = f
	return m.mock
}

//Init implements github.com/insolar/insolar/network/controller.RPCController interface
func (m *RPCControllerMock) Init(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.InitPreCounter, 1)
	defer atomic.AddUint64(&m.InitCounter, 1)

	if len(m.InitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RPCControllerMock.Init. %v", p)
			return
		}

		input := m.InitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RPCControllerMockInitInput{p}, "RPCController.Init got unexpected parameters")

		result := m.InitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.Init")
			return
		}

		r = result.r

		return
	}

	if m.InitMock.mainExpectation != nil {

		input := m.InitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RPCControllerMockInitInput{p}, "RPCController.Init got unexpected parameters")
		}

		result := m.InitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.Init")
		}

		r = result.r

		return
	}

	if m.InitFunc == nil {
		m.t.Fatalf("Unexpected call to RPCControllerMock.Init. %v", p)
		return
	}

	return m.InitFunc(p)
}

//InitMinimockCounter returns a count of RPCControllerMock.InitFunc invocations
func (m *RPCControllerMock) InitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InitCounter)
}

//InitMinimockPreCounter returns the value of RPCControllerMock.Init invocations
func (m *RPCControllerMock) InitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InitPreCounter)
}

//InitFinished returns true if mock invocations count is ok
func (m *RPCControllerMock) InitFinished() bool {
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

type mRPCControllerMockRemoteProcedureRegister struct {
	mock              *RPCControllerMock
	mainExpectation   *RPCControllerMockRemoteProcedureRegisterExpectation
	expectationSeries []*RPCControllerMockRemoteProcedureRegisterExpectation
}

type RPCControllerMockRemoteProcedureRegisterExpectation struct {
	input *RPCControllerMockRemoteProcedureRegisterInput
}

type RPCControllerMockRemoteProcedureRegisterInput struct {
	p  string
	p1 insolar.RemoteProcedure
}

//Expect specifies that invocation of RPCController.RemoteProcedureRegister is expected from 1 to Infinity times
func (m *mRPCControllerMockRemoteProcedureRegister) Expect(p string, p1 insolar.RemoteProcedure) *mRPCControllerMockRemoteProcedureRegister {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockRemoteProcedureRegisterExpectation{}
	}
	m.mainExpectation.input = &RPCControllerMockRemoteProcedureRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of RPCController.RemoteProcedureRegister
func (m *mRPCControllerMockRemoteProcedureRegister) Return() *RPCControllerMock {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockRemoteProcedureRegisterExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RPCController.RemoteProcedureRegister is expected once
func (m *mRPCControllerMockRemoteProcedureRegister) ExpectOnce(p string, p1 insolar.RemoteProcedure) *RPCControllerMockRemoteProcedureRegisterExpectation {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.mainExpectation = nil

	expectation := &RPCControllerMockRemoteProcedureRegisterExpectation{}
	expectation.input = &RPCControllerMockRemoteProcedureRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RPCController.RemoteProcedureRegister method
func (m *mRPCControllerMockRemoteProcedureRegister) Set(f func(p string, p1 insolar.RemoteProcedure)) *RPCControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoteProcedureRegisterFunc = f
	return m.mock
}

//RemoteProcedureRegister implements github.com/insolar/insolar/network/controller.RPCController interface
func (m *RPCControllerMock) RemoteProcedureRegister(p string, p1 insolar.RemoteProcedure) {
	counter := atomic.AddUint64(&m.RemoteProcedureRegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RemoteProcedureRegisterCounter, 1)

	if len(m.RemoteProcedureRegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoteProcedureRegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RPCControllerMock.RemoteProcedureRegister. %v %v", p, p1)
			return
		}

		input := m.RemoteProcedureRegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RPCControllerMockRemoteProcedureRegisterInput{p, p1}, "RPCController.RemoteProcedureRegister got unexpected parameters")

		return
	}

	if m.RemoteProcedureRegisterMock.mainExpectation != nil {

		input := m.RemoteProcedureRegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RPCControllerMockRemoteProcedureRegisterInput{p, p1}, "RPCController.RemoteProcedureRegister got unexpected parameters")
		}

		return
	}

	if m.RemoteProcedureRegisterFunc == nil {
		m.t.Fatalf("Unexpected call to RPCControllerMock.RemoteProcedureRegister. %v %v", p, p1)
		return
	}

	m.RemoteProcedureRegisterFunc(p, p1)
}

//RemoteProcedureRegisterMinimockCounter returns a count of RPCControllerMock.RemoteProcedureRegisterFunc invocations
func (m *RPCControllerMock) RemoteProcedureRegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter)
}

//RemoteProcedureRegisterMinimockPreCounter returns the value of RPCControllerMock.RemoteProcedureRegister invocations
func (m *RPCControllerMock) RemoteProcedureRegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoteProcedureRegisterPreCounter)
}

//RemoteProcedureRegisterFinished returns true if mock invocations count is ok
func (m *RPCControllerMock) RemoteProcedureRegisterFinished() bool {
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

type mRPCControllerMockSendBytes struct {
	mock              *RPCControllerMock
	mainExpectation   *RPCControllerMockSendBytesExpectation
	expectationSeries []*RPCControllerMockSendBytesExpectation
}

type RPCControllerMockSendBytesExpectation struct {
	input  *RPCControllerMockSendBytesInput
	result *RPCControllerMockSendBytesResult
}

type RPCControllerMockSendBytesInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 string
	p3 []byte
}

type RPCControllerMockSendBytesResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of RPCController.SendBytes is expected from 1 to Infinity times
func (m *mRPCControllerMockSendBytes) Expect(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) *mRPCControllerMockSendBytes {
	m.mock.SendBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockSendBytesExpectation{}
	}
	m.mainExpectation.input = &RPCControllerMockSendBytesInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of RPCController.SendBytes
func (m *mRPCControllerMockSendBytes) Return(r []byte, r1 error) *RPCControllerMock {
	m.mock.SendBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockSendBytesExpectation{}
	}
	m.mainExpectation.result = &RPCControllerMockSendBytesResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RPCController.SendBytes is expected once
func (m *mRPCControllerMockSendBytes) ExpectOnce(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) *RPCControllerMockSendBytesExpectation {
	m.mock.SendBytesFunc = nil
	m.mainExpectation = nil

	expectation := &RPCControllerMockSendBytesExpectation{}
	expectation.input = &RPCControllerMockSendBytesInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RPCControllerMockSendBytesExpectation) Return(r []byte, r1 error) {
	e.result = &RPCControllerMockSendBytesResult{r, r1}
}

//Set uses given function f as a mock of RPCController.SendBytes method
func (m *mRPCControllerMockSendBytes) Set(f func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error)) *RPCControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendBytesFunc = f
	return m.mock
}

//SendBytes implements github.com/insolar/insolar/network/controller.RPCController interface
func (m *RPCControllerMock) SendBytes(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.SendBytesPreCounter, 1)
	defer atomic.AddUint64(&m.SendBytesCounter, 1)

	if len(m.SendBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RPCControllerMock.SendBytes. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendBytesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RPCControllerMockSendBytesInput{p, p1, p2, p3}, "RPCController.SendBytes got unexpected parameters")

		result := m.SendBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.SendBytes")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendBytesMock.mainExpectation != nil {

		input := m.SendBytesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RPCControllerMockSendBytesInput{p, p1, p2, p3}, "RPCController.SendBytes got unexpected parameters")
		}

		result := m.SendBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.SendBytes")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendBytesFunc == nil {
		m.t.Fatalf("Unexpected call to RPCControllerMock.SendBytes. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SendBytesFunc(p, p1, p2, p3)
}

//SendBytesMinimockCounter returns a count of RPCControllerMock.SendBytesFunc invocations
func (m *RPCControllerMock) SendBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendBytesCounter)
}

//SendBytesMinimockPreCounter returns the value of RPCControllerMock.SendBytes invocations
func (m *RPCControllerMock) SendBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendBytesPreCounter)
}

//SendBytesFinished returns true if mock invocations count is ok
func (m *RPCControllerMock) SendBytesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendBytesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendBytesCounter) == uint64(len(m.SendBytesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendBytesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendBytesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendBytesFunc != nil {
		return atomic.LoadUint64(&m.SendBytesCounter) > 0
	}

	return true
}

type mRPCControllerMockSendCascadeMessage struct {
	mock              *RPCControllerMock
	mainExpectation   *RPCControllerMockSendCascadeMessageExpectation
	expectationSeries []*RPCControllerMockSendCascadeMessageExpectation
}

type RPCControllerMockSendCascadeMessageExpectation struct {
	input  *RPCControllerMockSendCascadeMessageInput
	result *RPCControllerMockSendCascadeMessageResult
}

type RPCControllerMockSendCascadeMessageInput struct {
	p  insolar.Cascade
	p1 string
	p2 insolar.Parcel
}

type RPCControllerMockSendCascadeMessageResult struct {
	r error
}

//Expect specifies that invocation of RPCController.SendCascadeMessage is expected from 1 to Infinity times
func (m *mRPCControllerMockSendCascadeMessage) Expect(p insolar.Cascade, p1 string, p2 insolar.Parcel) *mRPCControllerMockSendCascadeMessage {
	m.mock.SendCascadeMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockSendCascadeMessageExpectation{}
	}
	m.mainExpectation.input = &RPCControllerMockSendCascadeMessageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RPCController.SendCascadeMessage
func (m *mRPCControllerMockSendCascadeMessage) Return(r error) *RPCControllerMock {
	m.mock.SendCascadeMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockSendCascadeMessageExpectation{}
	}
	m.mainExpectation.result = &RPCControllerMockSendCascadeMessageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of RPCController.SendCascadeMessage is expected once
func (m *mRPCControllerMockSendCascadeMessage) ExpectOnce(p insolar.Cascade, p1 string, p2 insolar.Parcel) *RPCControllerMockSendCascadeMessageExpectation {
	m.mock.SendCascadeMessageFunc = nil
	m.mainExpectation = nil

	expectation := &RPCControllerMockSendCascadeMessageExpectation{}
	expectation.input = &RPCControllerMockSendCascadeMessageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RPCControllerMockSendCascadeMessageExpectation) Return(r error) {
	e.result = &RPCControllerMockSendCascadeMessageResult{r}
}

//Set uses given function f as a mock of RPCController.SendCascadeMessage method
func (m *mRPCControllerMockSendCascadeMessage) Set(f func(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error)) *RPCControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendCascadeMessageFunc = f
	return m.mock
}

//SendCascadeMessage implements github.com/insolar/insolar/network/controller.RPCController interface
func (m *RPCControllerMock) SendCascadeMessage(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error) {
	counter := atomic.AddUint64(&m.SendCascadeMessagePreCounter, 1)
	defer atomic.AddUint64(&m.SendCascadeMessageCounter, 1)

	if len(m.SendCascadeMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendCascadeMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RPCControllerMock.SendCascadeMessage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendCascadeMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RPCControllerMockSendCascadeMessageInput{p, p1, p2}, "RPCController.SendCascadeMessage got unexpected parameters")

		result := m.SendCascadeMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.SendCascadeMessage")
			return
		}

		r = result.r

		return
	}

	if m.SendCascadeMessageMock.mainExpectation != nil {

		input := m.SendCascadeMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RPCControllerMockSendCascadeMessageInput{p, p1, p2}, "RPCController.SendCascadeMessage got unexpected parameters")
		}

		result := m.SendCascadeMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.SendCascadeMessage")
		}

		r = result.r

		return
	}

	if m.SendCascadeMessageFunc == nil {
		m.t.Fatalf("Unexpected call to RPCControllerMock.SendCascadeMessage. %v %v %v", p, p1, p2)
		return
	}

	return m.SendCascadeMessageFunc(p, p1, p2)
}

//SendCascadeMessageMinimockCounter returns a count of RPCControllerMock.SendCascadeMessageFunc invocations
func (m *RPCControllerMock) SendCascadeMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCascadeMessageCounter)
}

//SendCascadeMessageMinimockPreCounter returns the value of RPCControllerMock.SendCascadeMessage invocations
func (m *RPCControllerMock) SendCascadeMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendCascadeMessagePreCounter)
}

//SendCascadeMessageFinished returns true if mock invocations count is ok
func (m *RPCControllerMock) SendCascadeMessageFinished() bool {
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

type mRPCControllerMockSendMessage struct {
	mock              *RPCControllerMock
	mainExpectation   *RPCControllerMockSendMessageExpectation
	expectationSeries []*RPCControllerMockSendMessageExpectation
}

type RPCControllerMockSendMessageExpectation struct {
	input  *RPCControllerMockSendMessageInput
	result *RPCControllerMockSendMessageResult
}

type RPCControllerMockSendMessageInput struct {
	p  insolar.Reference
	p1 string
	p2 insolar.Parcel
}

type RPCControllerMockSendMessageResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of RPCController.SendMessage is expected from 1 to Infinity times
func (m *mRPCControllerMockSendMessage) Expect(p insolar.Reference, p1 string, p2 insolar.Parcel) *mRPCControllerMockSendMessage {
	m.mock.SendMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockSendMessageExpectation{}
	}
	m.mainExpectation.input = &RPCControllerMockSendMessageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RPCController.SendMessage
func (m *mRPCControllerMockSendMessage) Return(r []byte, r1 error) *RPCControllerMock {
	m.mock.SendMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RPCControllerMockSendMessageExpectation{}
	}
	m.mainExpectation.result = &RPCControllerMockSendMessageResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RPCController.SendMessage is expected once
func (m *mRPCControllerMockSendMessage) ExpectOnce(p insolar.Reference, p1 string, p2 insolar.Parcel) *RPCControllerMockSendMessageExpectation {
	m.mock.SendMessageFunc = nil
	m.mainExpectation = nil

	expectation := &RPCControllerMockSendMessageExpectation{}
	expectation.input = &RPCControllerMockSendMessageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RPCControllerMockSendMessageExpectation) Return(r []byte, r1 error) {
	e.result = &RPCControllerMockSendMessageResult{r, r1}
}

//Set uses given function f as a mock of RPCController.SendMessage method
func (m *mRPCControllerMockSendMessage) Set(f func(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error)) *RPCControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendMessageFunc = f
	return m.mock
}

//SendMessage implements github.com/insolar/insolar/network/controller.RPCController interface
func (m *RPCControllerMock) SendMessage(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.SendMessagePreCounter, 1)
	defer atomic.AddUint64(&m.SendMessageCounter, 1)

	if len(m.SendMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RPCControllerMock.SendMessage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RPCControllerMockSendMessageInput{p, p1, p2}, "RPCController.SendMessage got unexpected parameters")

		result := m.SendMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.SendMessage")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMessageMock.mainExpectation != nil {

		input := m.SendMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RPCControllerMockSendMessageInput{p, p1, p2}, "RPCController.SendMessage got unexpected parameters")
		}

		result := m.SendMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RPCControllerMock.SendMessage")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMessageFunc == nil {
		m.t.Fatalf("Unexpected call to RPCControllerMock.SendMessage. %v %v %v", p, p1, p2)
		return
	}

	return m.SendMessageFunc(p, p1, p2)
}

//SendMessageMinimockCounter returns a count of RPCControllerMock.SendMessageFunc invocations
func (m *RPCControllerMock) SendMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessageCounter)
}

//SendMessageMinimockPreCounter returns the value of RPCControllerMock.SendMessage invocations
func (m *RPCControllerMock) SendMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessagePreCounter)
}

//SendMessageFinished returns true if mock invocations count is ok
func (m *RPCControllerMock) SendMessageFinished() bool {
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
func (m *RPCControllerMock) ValidateCallCounters() {

	if !m.IAmRPCControllerFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.IAmRPCController")
	}

	if !m.InitFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.Init")
	}

	if !m.RemoteProcedureRegisterFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.RemoteProcedureRegister")
	}

	if !m.SendBytesFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.SendBytes")
	}

	if !m.SendCascadeMessageFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.SendCascadeMessage")
	}

	if !m.SendMessageFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.SendMessage")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RPCControllerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RPCControllerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RPCControllerMock) MinimockFinish() {

	if !m.IAmRPCControllerFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.IAmRPCController")
	}

	if !m.InitFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.Init")
	}

	if !m.RemoteProcedureRegisterFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.RemoteProcedureRegister")
	}

	if !m.SendBytesFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.SendBytes")
	}

	if !m.SendCascadeMessageFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.SendCascadeMessage")
	}

	if !m.SendMessageFinished() {
		m.t.Fatal("Expected call to RPCControllerMock.SendMessage")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RPCControllerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RPCControllerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.IAmRPCControllerFinished()
		ok = ok && m.InitFinished()
		ok = ok && m.RemoteProcedureRegisterFinished()
		ok = ok && m.SendBytesFinished()
		ok = ok && m.SendCascadeMessageFinished()
		ok = ok && m.SendMessageFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IAmRPCControllerFinished() {
				m.t.Error("Expected call to RPCControllerMock.IAmRPCController")
			}

			if !m.InitFinished() {
				m.t.Error("Expected call to RPCControllerMock.Init")
			}

			if !m.RemoteProcedureRegisterFinished() {
				m.t.Error("Expected call to RPCControllerMock.RemoteProcedureRegister")
			}

			if !m.SendBytesFinished() {
				m.t.Error("Expected call to RPCControllerMock.SendBytes")
			}

			if !m.SendCascadeMessageFinished() {
				m.t.Error("Expected call to RPCControllerMock.SendCascadeMessage")
			}

			if !m.SendMessageFinished() {
				m.t.Error("Expected call to RPCControllerMock.SendMessage")
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
func (m *RPCControllerMock) AllMocksCalled() bool {

	if !m.IAmRPCControllerFinished() {
		return false
	}

	if !m.InitFinished() {
		return false
	}

	if !m.RemoteProcedureRegisterFinished() {
		return false
	}

	if !m.SendBytesFinished() {
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

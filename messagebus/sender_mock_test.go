package messagebus

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "sender" can be found in github.com/insolar/insolar/messagebus
*/
import (
	context "context"
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	testify_assert "github.com/stretchr/testify/assert"
)

//senderMock implements github.com/insolar/insolar/messagebus.sender
type senderMock struct {
	t minimock.Tester

	CreateParcelFunc       func(p context.Context, p1 core.Message, p2 core.DelegationToken) (r core.Parcel, r1 error)
	CreateParcelCounter    uint64
	CreateParcelPreCounter uint64
	CreateParcelMock       msenderMockCreateParcel

	MustRegisterFunc       func(p core.MessageType, p1 core.MessageHandler)
	MustRegisterCounter    uint64
	MustRegisterPreCounter uint64
	MustRegisterMock       msenderMockMustRegister

	NewPlayerFunc       func(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error)
	NewPlayerCounter    uint64
	NewPlayerPreCounter uint64
	NewPlayerMock       msenderMockNewPlayer

	NewRecorderFunc       func(p context.Context) (r core.MessageBus, r1 error)
	NewRecorderCounter    uint64
	NewRecorderPreCounter uint64
	NewRecorderMock       msenderMockNewRecorder

	RegisterFunc       func(p core.MessageType, p1 core.MessageHandler) (r error)
	RegisterCounter    uint64
	RegisterPreCounter uint64
	RegisterMock       msenderMockRegister

	SendFunc       func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       msenderMockSend

	SendParcelFunc       func(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) (r core.Reply, r1 error)
	SendParcelCounter    uint64
	SendParcelPreCounter uint64
	SendParcelMock       msenderMockSendParcel

	WriteTapeFunc       func(p context.Context, p1 io.Writer) (r error)
	WriteTapeCounter    uint64
	WriteTapePreCounter uint64
	WriteTapeMock       msenderMockWriteTape
}

//NewsenderMock returns a mock for github.com/insolar/insolar/messagebus.sender
func NewsenderMock(t minimock.Tester) *senderMock {
	m := &senderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateParcelMock = msenderMockCreateParcel{mock: m}
	m.MustRegisterMock = msenderMockMustRegister{mock: m}
	m.NewPlayerMock = msenderMockNewPlayer{mock: m}
	m.NewRecorderMock = msenderMockNewRecorder{mock: m}
	m.RegisterMock = msenderMockRegister{mock: m}
	m.SendMock = msenderMockSend{mock: m}
	m.SendParcelMock = msenderMockSendParcel{mock: m}
	m.WriteTapeMock = msenderMockWriteTape{mock: m}

	return m
}

type msenderMockCreateParcel struct {
	mock             *senderMock
	mockExpectations *senderMockCreateParcelParams
}

//senderMockCreateParcelParams represents input parameters of the sender.CreateParcel
type senderMockCreateParcelParams struct {
	p  context.Context
	p1 core.Message
	p2 core.DelegationToken
}

//Expect sets up expected params for the sender.CreateParcel
func (m *msenderMockCreateParcel) Expect(p context.Context, p1 core.Message, p2 core.DelegationToken) *msenderMockCreateParcel {
	m.mockExpectations = &senderMockCreateParcelParams{p, p1, p2}
	return m
}

//Return sets up a mock for sender.CreateParcel to return Return's arguments
func (m *msenderMockCreateParcel) Return(r core.Parcel, r1 error) *senderMock {
	m.mock.CreateParcelFunc = func(p context.Context, p1 core.Message, p2 core.DelegationToken) (core.Parcel, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.CreateParcel method
func (m *msenderMockCreateParcel) Set(f func(p context.Context, p1 core.Message, p2 core.DelegationToken) (r core.Parcel, r1 error)) *senderMock {
	m.mock.CreateParcelFunc = f
	m.mockExpectations = nil
	return m.mock
}

//CreateParcel implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) CreateParcel(p context.Context, p1 core.Message, p2 core.DelegationToken) (r core.Parcel, r1 error) {
	atomic.AddUint64(&m.CreateParcelPreCounter, 1)
	defer atomic.AddUint64(&m.CreateParcelCounter, 1)

	if m.CreateParcelMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.CreateParcelMock.mockExpectations, senderMockCreateParcelParams{p, p1, p2},
			"sender.CreateParcel got unexpected parameters")

		if m.CreateParcelFunc == nil {

			m.t.Fatal("No results are set for the senderMock.CreateParcel")

			return
		}
	}

	if m.CreateParcelFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.CreateParcel")
		return
	}

	return m.CreateParcelFunc(p, p1, p2)
}

//CreateParcelMinimockCounter returns a count of senderMock.CreateParcelFunc invocations
func (m *senderMock) CreateParcelMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateParcelCounter)
}

//CreateParcelMinimockPreCounter returns the value of senderMock.CreateParcel invocations
func (m *senderMock) CreateParcelMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateParcelPreCounter)
}

type msenderMockMustRegister struct {
	mock             *senderMock
	mockExpectations *senderMockMustRegisterParams
}

//senderMockMustRegisterParams represents input parameters of the sender.MustRegister
type senderMockMustRegisterParams struct {
	p  core.MessageType
	p1 core.MessageHandler
}

//Expect sets up expected params for the sender.MustRegister
func (m *msenderMockMustRegister) Expect(p core.MessageType, p1 core.MessageHandler) *msenderMockMustRegister {
	m.mockExpectations = &senderMockMustRegisterParams{p, p1}
	return m
}

//Return sets up a mock for sender.MustRegister to return Return's arguments
func (m *msenderMockMustRegister) Return() *senderMock {
	m.mock.MustRegisterFunc = func(p core.MessageType, p1 core.MessageHandler) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of sender.MustRegister method
func (m *msenderMockMustRegister) Set(f func(p core.MessageType, p1 core.MessageHandler)) *senderMock {
	m.mock.MustRegisterFunc = f
	m.mockExpectations = nil
	return m.mock
}

//MustRegister implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) MustRegister(p core.MessageType, p1 core.MessageHandler) {
	atomic.AddUint64(&m.MustRegisterPreCounter, 1)
	defer atomic.AddUint64(&m.MustRegisterCounter, 1)

	if m.MustRegisterMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.MustRegisterMock.mockExpectations, senderMockMustRegisterParams{p, p1},
			"sender.MustRegister got unexpected parameters")

		if m.MustRegisterFunc == nil {

			m.t.Fatal("No results are set for the senderMock.MustRegister")

			return
		}
	}

	if m.MustRegisterFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.MustRegister")
		return
	}

	m.MustRegisterFunc(p, p1)
}

//MustRegisterMinimockCounter returns a count of senderMock.MustRegisterFunc invocations
func (m *senderMock) MustRegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustRegisterCounter)
}

//MustRegisterMinimockPreCounter returns the value of senderMock.MustRegister invocations
func (m *senderMock) MustRegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustRegisterPreCounter)
}

type msenderMockNewPlayer struct {
	mock             *senderMock
	mockExpectations *senderMockNewPlayerParams
}

//senderMockNewPlayerParams represents input parameters of the sender.NewPlayer
type senderMockNewPlayerParams struct {
	p  context.Context
	p1 io.Reader
}

//Expect sets up expected params for the sender.NewPlayer
func (m *msenderMockNewPlayer) Expect(p context.Context, p1 io.Reader) *msenderMockNewPlayer {
	m.mockExpectations = &senderMockNewPlayerParams{p, p1}
	return m
}

//Return sets up a mock for sender.NewPlayer to return Return's arguments
func (m *msenderMockNewPlayer) Return(r core.MessageBus, r1 error) *senderMock {
	m.mock.NewPlayerFunc = func(p context.Context, p1 io.Reader) (core.MessageBus, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.NewPlayer method
func (m *msenderMockNewPlayer) Set(f func(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error)) *senderMock {
	m.mock.NewPlayerFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NewPlayer implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) NewPlayer(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error) {
	atomic.AddUint64(&m.NewPlayerPreCounter, 1)
	defer atomic.AddUint64(&m.NewPlayerCounter, 1)

	if m.NewPlayerMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NewPlayerMock.mockExpectations, senderMockNewPlayerParams{p, p1},
			"sender.NewPlayer got unexpected parameters")

		if m.NewPlayerFunc == nil {

			m.t.Fatal("No results are set for the senderMock.NewPlayer")

			return
		}
	}

	if m.NewPlayerFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.NewPlayer")
		return
	}

	return m.NewPlayerFunc(p, p1)
}

//NewPlayerMinimockCounter returns a count of senderMock.NewPlayerFunc invocations
func (m *senderMock) NewPlayerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewPlayerCounter)
}

//NewPlayerMinimockPreCounter returns the value of senderMock.NewPlayer invocations
func (m *senderMock) NewPlayerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewPlayerPreCounter)
}

type msenderMockNewRecorder struct {
	mock             *senderMock
	mockExpectations *senderMockNewRecorderParams
}

//senderMockNewRecorderParams represents input parameters of the sender.NewRecorder
type senderMockNewRecorderParams struct {
	p context.Context
}

//Expect sets up expected params for the sender.NewRecorder
func (m *msenderMockNewRecorder) Expect(p context.Context) *msenderMockNewRecorder {
	m.mockExpectations = &senderMockNewRecorderParams{p}
	return m
}

//Return sets up a mock for sender.NewRecorder to return Return's arguments
func (m *msenderMockNewRecorder) Return(r core.MessageBus, r1 error) *senderMock {
	m.mock.NewRecorderFunc = func(p context.Context) (core.MessageBus, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.NewRecorder method
func (m *msenderMockNewRecorder) Set(f func(p context.Context) (r core.MessageBus, r1 error)) *senderMock {
	m.mock.NewRecorderFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NewRecorder implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) NewRecorder(p context.Context) (r core.MessageBus, r1 error) {
	atomic.AddUint64(&m.NewRecorderPreCounter, 1)
	defer atomic.AddUint64(&m.NewRecorderCounter, 1)

	if m.NewRecorderMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NewRecorderMock.mockExpectations, senderMockNewRecorderParams{p},
			"sender.NewRecorder got unexpected parameters")

		if m.NewRecorderFunc == nil {

			m.t.Fatal("No results are set for the senderMock.NewRecorder")

			return
		}
	}

	if m.NewRecorderFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.NewRecorder")
		return
	}

	return m.NewRecorderFunc(p)
}

//NewRecorderMinimockCounter returns a count of senderMock.NewRecorderFunc invocations
func (m *senderMock) NewRecorderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewRecorderCounter)
}

//NewRecorderMinimockPreCounter returns the value of senderMock.NewRecorder invocations
func (m *senderMock) NewRecorderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewRecorderPreCounter)
}

type msenderMockRegister struct {
	mock             *senderMock
	mockExpectations *senderMockRegisterParams
}

//senderMockRegisterParams represents input parameters of the sender.Register
type senderMockRegisterParams struct {
	p  core.MessageType
	p1 core.MessageHandler
}

//Expect sets up expected params for the sender.Register
func (m *msenderMockRegister) Expect(p core.MessageType, p1 core.MessageHandler) *msenderMockRegister {
	m.mockExpectations = &senderMockRegisterParams{p, p1}
	return m
}

//Return sets up a mock for sender.Register to return Return's arguments
func (m *msenderMockRegister) Return(r error) *senderMock {
	m.mock.RegisterFunc = func(p core.MessageType, p1 core.MessageHandler) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of sender.Register method
func (m *msenderMockRegister) Set(f func(p core.MessageType, p1 core.MessageHandler) (r error)) *senderMock {
	m.mock.RegisterFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Register implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) Register(p core.MessageType, p1 core.MessageHandler) (r error) {
	atomic.AddUint64(&m.RegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterCounter, 1)

	if m.RegisterMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RegisterMock.mockExpectations, senderMockRegisterParams{p, p1},
			"sender.Register got unexpected parameters")

		if m.RegisterFunc == nil {

			m.t.Fatal("No results are set for the senderMock.Register")

			return
		}
	}

	if m.RegisterFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.Register")
		return
	}

	return m.RegisterFunc(p, p1)
}

//RegisterMinimockCounter returns a count of senderMock.RegisterFunc invocations
func (m *senderMock) RegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterCounter)
}

//RegisterMinimockPreCounter returns the value of senderMock.Register invocations
func (m *senderMock) RegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterPreCounter)
}

type msenderMockSend struct {
	mock             *senderMock
	mockExpectations *senderMockSendParams
}

//senderMockSendParams represents input parameters of the sender.Send
type senderMockSendParams struct {
	p  context.Context
	p1 core.Message
	p2 *core.MessageSendOptions
}

//Expect sets up expected params for the sender.Send
func (m *msenderMockSend) Expect(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) *msenderMockSend {
	m.mockExpectations = &senderMockSendParams{p, p1, p2}
	return m
}

//Return sets up a mock for sender.Send to return Return's arguments
func (m *msenderMockSend) Return(r core.Reply, r1 error) *senderMock {
	m.mock.SendFunc = func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.Send method
func (m *msenderMockSend) Set(f func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error)) *senderMock {
	m.mock.SendFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Send implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) Send(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if m.SendMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendMock.mockExpectations, senderMockSendParams{p, p1, p2},
			"sender.Send got unexpected parameters")

		if m.SendFunc == nil {

			m.t.Fatal("No results are set for the senderMock.Send")

			return
		}
	}

	if m.SendFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.Send")
		return
	}

	return m.SendFunc(p, p1, p2)
}

//SendMinimockCounter returns a count of senderMock.SendFunc invocations
func (m *senderMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of senderMock.Send invocations
func (m *senderMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

type msenderMockSendParcel struct {
	mock             *senderMock
	mockExpectations *senderMockSendParcelParams
}

//senderMockSendParcelParams represents input parameters of the sender.SendParcel
type senderMockSendParcelParams struct {
	p  context.Context
	p1 core.Parcel
	p2 *core.MessageSendOptions
}

//Expect sets up expected params for the sender.SendParcel
func (m *msenderMockSendParcel) Expect(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) *msenderMockSendParcel {
	m.mockExpectations = &senderMockSendParcelParams{p, p1, p2}
	return m
}

//Return sets up a mock for sender.SendParcel to return Return's arguments
func (m *msenderMockSendParcel) Return(r core.Reply, r1 error) *senderMock {
	m.mock.SendParcelFunc = func(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.SendParcel method
func (m *msenderMockSendParcel) Set(f func(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) (r core.Reply, r1 error)) *senderMock {
	m.mock.SendParcelFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SendParcel implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) SendParcel(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.SendParcelPreCounter, 1)
	defer atomic.AddUint64(&m.SendParcelCounter, 1)

	if m.SendParcelMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendParcelMock.mockExpectations, senderMockSendParcelParams{p, p1, p2},
			"sender.SendParcel got unexpected parameters")

		if m.SendParcelFunc == nil {

			m.t.Fatal("No results are set for the senderMock.SendParcel")

			return
		}
	}

	if m.SendParcelFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.SendParcel")
		return
	}

	return m.SendParcelFunc(p, p1, p2)
}

//SendParcelMinimockCounter returns a count of senderMock.SendParcelFunc invocations
func (m *senderMock) SendParcelMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendParcelCounter)
}

//SendParcelMinimockPreCounter returns the value of senderMock.SendParcel invocations
func (m *senderMock) SendParcelMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendParcelPreCounter)
}

type msenderMockWriteTape struct {
	mock             *senderMock
	mockExpectations *senderMockWriteTapeParams
}

//senderMockWriteTapeParams represents input parameters of the sender.WriteTape
type senderMockWriteTapeParams struct {
	p  context.Context
	p1 io.Writer
}

//Expect sets up expected params for the sender.WriteTape
func (m *msenderMockWriteTape) Expect(p context.Context, p1 io.Writer) *msenderMockWriteTape {
	m.mockExpectations = &senderMockWriteTapeParams{p, p1}
	return m
}

//Return sets up a mock for sender.WriteTape to return Return's arguments
func (m *msenderMockWriteTape) Return(r error) *senderMock {
	m.mock.WriteTapeFunc = func(p context.Context, p1 io.Writer) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of sender.WriteTape method
func (m *msenderMockWriteTape) Set(f func(p context.Context, p1 io.Writer) (r error)) *senderMock {
	m.mock.WriteTapeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//WriteTape implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) WriteTape(p context.Context, p1 io.Writer) (r error) {
	atomic.AddUint64(&m.WriteTapePreCounter, 1)
	defer atomic.AddUint64(&m.WriteTapeCounter, 1)

	if m.WriteTapeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.WriteTapeMock.mockExpectations, senderMockWriteTapeParams{p, p1},
			"sender.WriteTape got unexpected parameters")

		if m.WriteTapeFunc == nil {

			m.t.Fatal("No results are set for the senderMock.WriteTape")

			return
		}
	}

	if m.WriteTapeFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.WriteTape")
		return
	}

	return m.WriteTapeFunc(p, p1)
}

//WriteTapeMinimockCounter returns a count of senderMock.WriteTapeFunc invocations
func (m *senderMock) WriteTapeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteTapeCounter)
}

//WriteTapeMinimockPreCounter returns the value of senderMock.WriteTape invocations
func (m *senderMock) WriteTapeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteTapePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *senderMock) ValidateCallCounters() {

	if m.CreateParcelFunc != nil && atomic.LoadUint64(&m.CreateParcelCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.CreateParcel")
	}

	if m.MustRegisterFunc != nil && atomic.LoadUint64(&m.MustRegisterCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.MustRegister")
	}

	if m.NewPlayerFunc != nil && atomic.LoadUint64(&m.NewPlayerCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.NewPlayer")
	}

	if m.NewRecorderFunc != nil && atomic.LoadUint64(&m.NewRecorderCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.NewRecorder")
	}

	if m.RegisterFunc != nil && atomic.LoadUint64(&m.RegisterCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.Register")
	}

	if m.SendFunc != nil && atomic.LoadUint64(&m.SendCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.Send")
	}

	if m.SendParcelFunc != nil && atomic.LoadUint64(&m.SendParcelCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.SendParcel")
	}

	if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.WriteTape")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *senderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *senderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *senderMock) MinimockFinish() {

	if m.CreateParcelFunc != nil && atomic.LoadUint64(&m.CreateParcelCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.CreateParcel")
	}

	if m.MustRegisterFunc != nil && atomic.LoadUint64(&m.MustRegisterCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.MustRegister")
	}

	if m.NewPlayerFunc != nil && atomic.LoadUint64(&m.NewPlayerCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.NewPlayer")
	}

	if m.NewRecorderFunc != nil && atomic.LoadUint64(&m.NewRecorderCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.NewRecorder")
	}

	if m.RegisterFunc != nil && atomic.LoadUint64(&m.RegisterCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.Register")
	}

	if m.SendFunc != nil && atomic.LoadUint64(&m.SendCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.Send")
	}

	if m.SendParcelFunc != nil && atomic.LoadUint64(&m.SendParcelCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.SendParcel")
	}

	if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.WriteTape")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *senderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *senderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.CreateParcelFunc == nil || atomic.LoadUint64(&m.CreateParcelCounter) > 0)
		ok = ok && (m.MustRegisterFunc == nil || atomic.LoadUint64(&m.MustRegisterCounter) > 0)
		ok = ok && (m.NewPlayerFunc == nil || atomic.LoadUint64(&m.NewPlayerCounter) > 0)
		ok = ok && (m.NewRecorderFunc == nil || atomic.LoadUint64(&m.NewRecorderCounter) > 0)
		ok = ok && (m.RegisterFunc == nil || atomic.LoadUint64(&m.RegisterCounter) > 0)
		ok = ok && (m.SendFunc == nil || atomic.LoadUint64(&m.SendCounter) > 0)
		ok = ok && (m.SendParcelFunc == nil || atomic.LoadUint64(&m.SendParcelCounter) > 0)
		ok = ok && (m.WriteTapeFunc == nil || atomic.LoadUint64(&m.WriteTapeCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CreateParcelFunc != nil && atomic.LoadUint64(&m.CreateParcelCounter) == 0 {
				m.t.Error("Expected call to senderMock.CreateParcel")
			}

			if m.MustRegisterFunc != nil && atomic.LoadUint64(&m.MustRegisterCounter) == 0 {
				m.t.Error("Expected call to senderMock.MustRegister")
			}

			if m.NewPlayerFunc != nil && atomic.LoadUint64(&m.NewPlayerCounter) == 0 {
				m.t.Error("Expected call to senderMock.NewPlayer")
			}

			if m.NewRecorderFunc != nil && atomic.LoadUint64(&m.NewRecorderCounter) == 0 {
				m.t.Error("Expected call to senderMock.NewRecorder")
			}

			if m.RegisterFunc != nil && atomic.LoadUint64(&m.RegisterCounter) == 0 {
				m.t.Error("Expected call to senderMock.Register")
			}

			if m.SendFunc != nil && atomic.LoadUint64(&m.SendCounter) == 0 {
				m.t.Error("Expected call to senderMock.Send")
			}

			if m.SendParcelFunc != nil && atomic.LoadUint64(&m.SendParcelCounter) == 0 {
				m.t.Error("Expected call to senderMock.SendParcel")
			}

			if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
				m.t.Error("Expected call to senderMock.WriteTape")
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
func (m *senderMock) AllMocksCalled() bool {

	if m.CreateParcelFunc != nil && atomic.LoadUint64(&m.CreateParcelCounter) == 0 {
		return false
	}

	if m.MustRegisterFunc != nil && atomic.LoadUint64(&m.MustRegisterCounter) == 0 {
		return false
	}

	if m.NewPlayerFunc != nil && atomic.LoadUint64(&m.NewPlayerCounter) == 0 {
		return false
	}

	if m.NewRecorderFunc != nil && atomic.LoadUint64(&m.NewRecorderCounter) == 0 {
		return false
	}

	if m.RegisterFunc != nil && atomic.LoadUint64(&m.RegisterCounter) == 0 {
		return false
	}

	if m.SendFunc != nil && atomic.LoadUint64(&m.SendCounter) == 0 {
		return false
	}

	if m.SendParcelFunc != nil && atomic.LoadUint64(&m.SendParcelCounter) == 0 {
		return false
	}

	if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
		return false
	}

	return true
}

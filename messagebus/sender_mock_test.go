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

	CreateSignedMessageFunc       func(p context.Context, p1 core.PulseNumber, p2 core.Message) (r core.SignedMessage, r1 error)
	CreateSignedMessageCounter    uint64
	CreateSignedMessagePreCounter uint64
	CreateSignedMessageMock       msenderMockCreateSignedMessage

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

	SendFunc       func(p context.Context, p1 core.Message) (r core.Reply, r1 error)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       msenderMockSend

	SendMessageFunc       func(p context.Context, p1 *core.Pulse, p2 core.SignedMessage) (r core.Reply, r1 error)
	SendMessageCounter    uint64
	SendMessagePreCounter uint64
	SendMessageMock       msenderMockSendMessage

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

	m.CreateSignedMessageMock = msenderMockCreateSignedMessage{mock: m}
	m.MustRegisterMock = msenderMockMustRegister{mock: m}
	m.NewPlayerMock = msenderMockNewPlayer{mock: m}
	m.NewRecorderMock = msenderMockNewRecorder{mock: m}
	m.RegisterMock = msenderMockRegister{mock: m}
	m.SendMock = msenderMockSend{mock: m}
	m.SendMessageMock = msenderMockSendMessage{mock: m}
	m.WriteTapeMock = msenderMockWriteTape{mock: m}

	return m
}

type msenderMockCreateSignedMessage struct {
	mock             *senderMock
	mockExpectations *senderMockCreateSignedMessageParams
}

//senderMockCreateSignedMessageParams represents input parameters of the sender.CreateSignedMessage
type senderMockCreateSignedMessageParams struct {
	p  context.Context
	p1 core.PulseNumber
	p2 core.Message
}

//Expect sets up expected params for the sender.CreateSignedMessage
func (m *msenderMockCreateSignedMessage) Expect(p context.Context, p1 core.PulseNumber, p2 core.Message) *msenderMockCreateSignedMessage {
	m.mockExpectations = &senderMockCreateSignedMessageParams{p, p1, p2}
	return m
}

//Return sets up a mock for sender.CreateSignedMessage to return Return's arguments
func (m *msenderMockCreateSignedMessage) Return(r core.SignedMessage, r1 error) *senderMock {
	m.mock.CreateSignedMessageFunc = func(p context.Context, p1 core.PulseNumber, p2 core.Message) (core.SignedMessage, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.CreateSignedMessage method
func (m *msenderMockCreateSignedMessage) Set(f func(p context.Context, p1 core.PulseNumber, p2 core.Message) (r core.SignedMessage, r1 error)) *senderMock {
	m.mock.CreateSignedMessageFunc = f
	m.mockExpectations = nil
	return m.mock
}

//CreateSignedMessage implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) CreateSignedMessage(p context.Context, p1 core.PulseNumber, p2 core.Message) (r core.SignedMessage, r1 error) {
	atomic.AddUint64(&m.CreateSignedMessagePreCounter, 1)
	defer atomic.AddUint64(&m.CreateSignedMessageCounter, 1)

	if m.CreateSignedMessageMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.CreateSignedMessageMock.mockExpectations, senderMockCreateSignedMessageParams{p, p1, p2},
			"sender.CreateSignedMessage got unexpected parameters")

		if m.CreateSignedMessageFunc == nil {

			m.t.Fatal("No results are set for the senderMock.CreateSignedMessage")

			return
		}
	}

	if m.CreateSignedMessageFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.CreateSignedMessage")
		return
	}

	return m.CreateSignedMessageFunc(p, p1, p2)
}

//CreateSignedMessageMinimockCounter returns a count of senderMock.CreateSignedMessageFunc invocations
func (m *senderMock) CreateSignedMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateSignedMessageCounter)
}

//CreateSignedMessageMinimockPreCounter returns the value of senderMock.CreateSignedMessage invocations
func (m *senderMock) CreateSignedMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateSignedMessagePreCounter)
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
}

//Expect sets up expected params for the sender.Send
func (m *msenderMockSend) Expect(p context.Context, p1 core.Message) *msenderMockSend {
	m.mockExpectations = &senderMockSendParams{p, p1}
	return m
}

//Return sets up a mock for sender.Send to return Return's arguments
func (m *msenderMockSend) Return(r core.Reply, r1 error) *senderMock {
	m.mock.SendFunc = func(p context.Context, p1 core.Message) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.Send method
func (m *msenderMockSend) Set(f func(p context.Context, p1 core.Message) (r core.Reply, r1 error)) *senderMock {
	m.mock.SendFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Send implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) Send(p context.Context, p1 core.Message) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if m.SendMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendMock.mockExpectations, senderMockSendParams{p, p1},
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

	return m.SendFunc(p, p1)
}

//SendMinimockCounter returns a count of senderMock.SendFunc invocations
func (m *senderMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of senderMock.Send invocations
func (m *senderMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

type msenderMockSendMessage struct {
	mock             *senderMock
	mockExpectations *senderMockSendMessageParams
}

//senderMockSendMessageParams represents input parameters of the sender.SendMessage
type senderMockSendMessageParams struct {
	p  context.Context
	p1 *core.Pulse
	p2 core.SignedMessage
}

//Expect sets up expected params for the sender.SendMessage
func (m *msenderMockSendMessage) Expect(p context.Context, p1 *core.Pulse, p2 core.SignedMessage) *msenderMockSendMessage {
	m.mockExpectations = &senderMockSendMessageParams{p, p1, p2}
	return m
}

//Return sets up a mock for sender.SendMessage to return Return's arguments
func (m *msenderMockSendMessage) Return(r core.Reply, r1 error) *senderMock {
	m.mock.SendMessageFunc = func(p context.Context, p1 *core.Pulse, p2 core.SignedMessage) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of sender.SendMessage method
func (m *msenderMockSendMessage) Set(f func(p context.Context, p1 *core.Pulse, p2 core.SignedMessage) (r core.Reply, r1 error)) *senderMock {
	m.mock.SendMessageFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SendMessage implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) SendMessage(p context.Context, p1 *core.Pulse, p2 core.SignedMessage) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.SendMessagePreCounter, 1)
	defer atomic.AddUint64(&m.SendMessageCounter, 1)

	if m.SendMessageMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendMessageMock.mockExpectations, senderMockSendMessageParams{p, p1, p2},
			"sender.SendMessage got unexpected parameters")

		if m.SendMessageFunc == nil {

			m.t.Fatal("No results are set for the senderMock.SendMessage")

			return
		}
	}

	if m.SendMessageFunc == nil {
		m.t.Fatal("Unexpected call to senderMock.SendMessage")
		return
	}

	return m.SendMessageFunc(p, p1, p2)
}

//SendMessageMinimockCounter returns a count of senderMock.SendMessageFunc invocations
func (m *senderMock) SendMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessageCounter)
}

//SendMessageMinimockPreCounter returns the value of senderMock.SendMessage invocations
func (m *senderMock) SendMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessagePreCounter)
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

	if m.CreateSignedMessageFunc != nil && atomic.LoadUint64(&m.CreateSignedMessageCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.CreateSignedMessage")
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

	if m.SendMessageFunc != nil && atomic.LoadUint64(&m.SendMessageCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.SendMessage")
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

	if m.CreateSignedMessageFunc != nil && atomic.LoadUint64(&m.CreateSignedMessageCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.CreateSignedMessage")
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

	if m.SendMessageFunc != nil && atomic.LoadUint64(&m.SendMessageCounter) == 0 {
		m.t.Fatal("Expected call to senderMock.SendMessage")
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
		ok = ok && (m.CreateSignedMessageFunc == nil || atomic.LoadUint64(&m.CreateSignedMessageCounter) > 0)
		ok = ok && (m.MustRegisterFunc == nil || atomic.LoadUint64(&m.MustRegisterCounter) > 0)
		ok = ok && (m.NewPlayerFunc == nil || atomic.LoadUint64(&m.NewPlayerCounter) > 0)
		ok = ok && (m.NewRecorderFunc == nil || atomic.LoadUint64(&m.NewRecorderCounter) > 0)
		ok = ok && (m.RegisterFunc == nil || atomic.LoadUint64(&m.RegisterCounter) > 0)
		ok = ok && (m.SendFunc == nil || atomic.LoadUint64(&m.SendCounter) > 0)
		ok = ok && (m.SendMessageFunc == nil || atomic.LoadUint64(&m.SendMessageCounter) > 0)
		ok = ok && (m.WriteTapeFunc == nil || atomic.LoadUint64(&m.WriteTapeCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CreateSignedMessageFunc != nil && atomic.LoadUint64(&m.CreateSignedMessageCounter) == 0 {
				m.t.Error("Expected call to senderMock.CreateSignedMessage")
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

			if m.SendMessageFunc != nil && atomic.LoadUint64(&m.SendMessageCounter) == 0 {
				m.t.Error("Expected call to senderMock.SendMessage")
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

	if m.CreateSignedMessageFunc != nil && atomic.LoadUint64(&m.CreateSignedMessageCounter) == 0 {
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

	if m.SendMessageFunc != nil && atomic.LoadUint64(&m.SendMessageCounter) == 0 {
		return false
	}

	if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
		return false
	}

	return true
}

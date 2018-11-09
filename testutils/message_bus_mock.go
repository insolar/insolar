package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MessageBus" can be found in github.com/insolar/insolar/core
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

//MessageBusMock implements github.com/insolar/insolar/core.MessageBus
type MessageBusMock struct {
	t minimock.Tester

	MustRegisterFunc       func(p core.MessageType, p1 core.MessageHandler)
	MustRegisterCounter    uint64
	MustRegisterPreCounter uint64
	MustRegisterMock       mMessageBusMockMustRegister

	NewPlayerFunc       func(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error)
	NewPlayerCounter    uint64
	NewPlayerPreCounter uint64
	NewPlayerMock       mMessageBusMockNewPlayer

	NewRecorderFunc       func(p context.Context) (r core.MessageBus, r1 error)
	NewRecorderCounter    uint64
	NewRecorderPreCounter uint64
	NewRecorderMock       mMessageBusMockNewRecorder

	RegisterFunc       func(p core.MessageType, p1 core.MessageHandler) (r error)
	RegisterCounter    uint64
	RegisterPreCounter uint64
	RegisterMock       mMessageBusMockRegister

	SendFunc       func(p context.Context, p1 core.Message) (r core.Reply, r1 error)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       mMessageBusMockSend

	WriteTapeFunc       func(p context.Context, p1 io.Writer) (r error)
	WriteTapeCounter    uint64
	WriteTapePreCounter uint64
	WriteTapeMock       mMessageBusMockWriteTape
}

//NewMessageBusMock returns a mock for github.com/insolar/insolar/core.MessageBus
func NewMessageBusMock(t minimock.Tester) *MessageBusMock {
	m := &MessageBusMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.MustRegisterMock = mMessageBusMockMustRegister{mock: m}
	m.NewPlayerMock = mMessageBusMockNewPlayer{mock: m}
	m.NewRecorderMock = mMessageBusMockNewRecorder{mock: m}
	m.RegisterMock = mMessageBusMockRegister{mock: m}
	m.SendMock = mMessageBusMockSend{mock: m}
	m.WriteTapeMock = mMessageBusMockWriteTape{mock: m}

	return m
}

type mMessageBusMockMustRegister struct {
	mock             *MessageBusMock
	mockExpectations *MessageBusMockMustRegisterParams
}

//MessageBusMockMustRegisterParams represents input parameters of the MessageBus.MustRegister
type MessageBusMockMustRegisterParams struct {
	p  core.MessageType
	p1 core.MessageHandler
}

//Expect sets up expected params for the MessageBus.MustRegister
func (m *mMessageBusMockMustRegister) Expect(p core.MessageType, p1 core.MessageHandler) *mMessageBusMockMustRegister {
	m.mockExpectations = &MessageBusMockMustRegisterParams{p, p1}
	return m
}

//Return sets up a mock for MessageBus.MustRegister to return Return's arguments
func (m *mMessageBusMockMustRegister) Return() *MessageBusMock {
	m.mock.MustRegisterFunc = func(p core.MessageType, p1 core.MessageHandler) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of MessageBus.MustRegister method
func (m *mMessageBusMockMustRegister) Set(f func(p core.MessageType, p1 core.MessageHandler)) *MessageBusMock {
	m.mock.MustRegisterFunc = f
	m.mockExpectations = nil
	return m.mock
}

//MustRegister implements github.com/insolar/insolar/core.MessageBus interface
func (m *MessageBusMock) MustRegister(p core.MessageType, p1 core.MessageHandler) {
	atomic.AddUint64(&m.MustRegisterPreCounter, 1)
	defer atomic.AddUint64(&m.MustRegisterCounter, 1)

	if m.MustRegisterMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.MustRegisterMock.mockExpectations, MessageBusMockMustRegisterParams{p, p1},
			"MessageBus.MustRegister got unexpected parameters")

		if m.MustRegisterFunc == nil {

			m.t.Fatal("No results are set for the MessageBusMock.MustRegister")

			return
		}
	}

	if m.MustRegisterFunc == nil {
		m.t.Fatal("Unexpected call to MessageBusMock.MustRegister")
		return
	}

	m.MustRegisterFunc(p, p1)
}

//MustRegisterMinimockCounter returns a count of MessageBusMock.MustRegisterFunc invocations
func (m *MessageBusMock) MustRegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustRegisterCounter)
}

//MustRegisterMinimockPreCounter returns the value of MessageBusMock.MustRegister invocations
func (m *MessageBusMock) MustRegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustRegisterPreCounter)
}

type mMessageBusMockNewPlayer struct {
	mock             *MessageBusMock
	mockExpectations *MessageBusMockNewPlayerParams
}

//MessageBusMockNewPlayerParams represents input parameters of the MessageBus.NewPlayer
type MessageBusMockNewPlayerParams struct {
	p  context.Context
	p1 io.Reader
}

//Expect sets up expected params for the MessageBus.NewPlayer
func (m *mMessageBusMockNewPlayer) Expect(p context.Context, p1 io.Reader) *mMessageBusMockNewPlayer {
	m.mockExpectations = &MessageBusMockNewPlayerParams{p, p1}
	return m
}

//Return sets up a mock for MessageBus.NewPlayer to return Return's arguments
func (m *mMessageBusMockNewPlayer) Return(r core.MessageBus, r1 error) *MessageBusMock {
	m.mock.NewPlayerFunc = func(p context.Context, p1 io.Reader) (core.MessageBus, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of MessageBus.NewPlayer method
func (m *mMessageBusMockNewPlayer) Set(f func(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error)) *MessageBusMock {
	m.mock.NewPlayerFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NewPlayer implements github.com/insolar/insolar/core.MessageBus interface
func (m *MessageBusMock) NewPlayer(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error) {
	atomic.AddUint64(&m.NewPlayerPreCounter, 1)
	defer atomic.AddUint64(&m.NewPlayerCounter, 1)

	if m.NewPlayerMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NewPlayerMock.mockExpectations, MessageBusMockNewPlayerParams{p, p1},
			"MessageBus.NewPlayer got unexpected parameters")

		if m.NewPlayerFunc == nil {

			m.t.Fatal("No results are set for the MessageBusMock.NewPlayer")

			return
		}
	}

	if m.NewPlayerFunc == nil {
		m.t.Fatal("Unexpected call to MessageBusMock.NewPlayer")
		return
	}

	return m.NewPlayerFunc(p, p1)
}

//NewPlayerMinimockCounter returns a count of MessageBusMock.NewPlayerFunc invocations
func (m *MessageBusMock) NewPlayerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewPlayerCounter)
}

//NewPlayerMinimockPreCounter returns the value of MessageBusMock.NewPlayer invocations
func (m *MessageBusMock) NewPlayerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewPlayerPreCounter)
}

type mMessageBusMockNewRecorder struct {
	mock             *MessageBusMock
	mockExpectations *MessageBusMockNewRecorderParams
}

//MessageBusMockNewRecorderParams represents input parameters of the MessageBus.NewRecorder
type MessageBusMockNewRecorderParams struct {
	p context.Context
}

//Expect sets up expected params for the MessageBus.NewRecorder
func (m *mMessageBusMockNewRecorder) Expect(p context.Context) *mMessageBusMockNewRecorder {
	m.mockExpectations = &MessageBusMockNewRecorderParams{p}
	return m
}

//Return sets up a mock for MessageBus.NewRecorder to return Return's arguments
func (m *mMessageBusMockNewRecorder) Return(r core.MessageBus, r1 error) *MessageBusMock {
	m.mock.NewRecorderFunc = func(p context.Context) (core.MessageBus, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of MessageBus.NewRecorder method
func (m *mMessageBusMockNewRecorder) Set(f func(p context.Context) (r core.MessageBus, r1 error)) *MessageBusMock {
	m.mock.NewRecorderFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NewRecorder implements github.com/insolar/insolar/core.MessageBus interface
func (m *MessageBusMock) NewRecorder(p context.Context) (r core.MessageBus, r1 error) {
	atomic.AddUint64(&m.NewRecorderPreCounter, 1)
	defer atomic.AddUint64(&m.NewRecorderCounter, 1)

	if m.NewRecorderMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NewRecorderMock.mockExpectations, MessageBusMockNewRecorderParams{p},
			"MessageBus.NewRecorder got unexpected parameters")

		if m.NewRecorderFunc == nil {

			m.t.Fatal("No results are set for the MessageBusMock.NewRecorder")

			return
		}
	}

	if m.NewRecorderFunc == nil {
		m.t.Fatal("Unexpected call to MessageBusMock.NewRecorder")
		return
	}

	return m.NewRecorderFunc(p)
}

//NewRecorderMinimockCounter returns a count of MessageBusMock.NewRecorderFunc invocations
func (m *MessageBusMock) NewRecorderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewRecorderCounter)
}

//NewRecorderMinimockPreCounter returns the value of MessageBusMock.NewRecorder invocations
func (m *MessageBusMock) NewRecorderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewRecorderPreCounter)
}

type mMessageBusMockRegister struct {
	mock             *MessageBusMock
	mockExpectations *MessageBusMockRegisterParams
}

//MessageBusMockRegisterParams represents input parameters of the MessageBus.Register
type MessageBusMockRegisterParams struct {
	p  core.MessageType
	p1 core.MessageHandler
}

//Expect sets up expected params for the MessageBus.Register
func (m *mMessageBusMockRegister) Expect(p core.MessageType, p1 core.MessageHandler) *mMessageBusMockRegister {
	m.mockExpectations = &MessageBusMockRegisterParams{p, p1}
	return m
}

//Return sets up a mock for MessageBus.Register to return Return's arguments
func (m *mMessageBusMockRegister) Return(r error) *MessageBusMock {
	m.mock.RegisterFunc = func(p core.MessageType, p1 core.MessageHandler) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of MessageBus.Register method
func (m *mMessageBusMockRegister) Set(f func(p core.MessageType, p1 core.MessageHandler) (r error)) *MessageBusMock {
	m.mock.RegisterFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Register implements github.com/insolar/insolar/core.MessageBus interface
func (m *MessageBusMock) Register(p core.MessageType, p1 core.MessageHandler) (r error) {
	atomic.AddUint64(&m.RegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterCounter, 1)

	if m.RegisterMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RegisterMock.mockExpectations, MessageBusMockRegisterParams{p, p1},
			"MessageBus.Register got unexpected parameters")

		if m.RegisterFunc == nil {

			m.t.Fatal("No results are set for the MessageBusMock.Register")

			return
		}
	}

	if m.RegisterFunc == nil {
		m.t.Fatal("Unexpected call to MessageBusMock.Register")
		return
	}

	return m.RegisterFunc(p, p1)
}

//RegisterMinimockCounter returns a count of MessageBusMock.RegisterFunc invocations
func (m *MessageBusMock) RegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterCounter)
}

//RegisterMinimockPreCounter returns the value of MessageBusMock.Register invocations
func (m *MessageBusMock) RegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterPreCounter)
}

type mMessageBusMockSend struct {
	mock             *MessageBusMock
	mockExpectations *MessageBusMockSendParams
}

//MessageBusMockSendParams represents input parameters of the MessageBus.Send
type MessageBusMockSendParams struct {
	p  context.Context
	p1 core.Message
}

//Expect sets up expected params for the MessageBus.Send
func (m *mMessageBusMockSend) Expect(p context.Context, p1 core.Message) *mMessageBusMockSend {
	m.mockExpectations = &MessageBusMockSendParams{p, p1}
	return m
}

//Return sets up a mock for MessageBus.Send to return Return's arguments
func (m *mMessageBusMockSend) Return(r core.Reply, r1 error) *MessageBusMock {
	m.mock.SendFunc = func(p context.Context, p1 core.Message) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of MessageBus.Send method
func (m *mMessageBusMockSend) Set(f func(p context.Context, p1 core.Message) (r core.Reply, r1 error)) *MessageBusMock {
	m.mock.SendFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Send implements github.com/insolar/insolar/core.MessageBus interface
func (m *MessageBusMock) Send(p context.Context, p1 core.Message) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if m.SendMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SendMock.mockExpectations, MessageBusMockSendParams{p, p1},
			"MessageBus.Send got unexpected parameters")

		if m.SendFunc == nil {

			m.t.Fatal("No results are set for the MessageBusMock.Send")

			return
		}
	}

	if m.SendFunc == nil {
		m.t.Fatal("Unexpected call to MessageBusMock.Send")
		return
	}

	return m.SendFunc(p, p1)
}

//SendMinimockCounter returns a count of MessageBusMock.SendFunc invocations
func (m *MessageBusMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of MessageBusMock.Send invocations
func (m *MessageBusMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

type mMessageBusMockWriteTape struct {
	mock             *MessageBusMock
	mockExpectations *MessageBusMockWriteTapeParams
}

//MessageBusMockWriteTapeParams represents input parameters of the MessageBus.WriteTape
type MessageBusMockWriteTapeParams struct {
	p  context.Context
	p1 io.Writer
}

//Expect sets up expected params for the MessageBus.WriteTape
func (m *mMessageBusMockWriteTape) Expect(p context.Context, p1 io.Writer) *mMessageBusMockWriteTape {
	m.mockExpectations = &MessageBusMockWriteTapeParams{p, p1}
	return m
}

//Return sets up a mock for MessageBus.WriteTape to return Return's arguments
func (m *mMessageBusMockWriteTape) Return(r error) *MessageBusMock {
	m.mock.WriteTapeFunc = func(p context.Context, p1 io.Writer) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of MessageBus.WriteTape method
func (m *mMessageBusMockWriteTape) Set(f func(p context.Context, p1 io.Writer) (r error)) *MessageBusMock {
	m.mock.WriteTapeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//WriteTape implements github.com/insolar/insolar/core.MessageBus interface
func (m *MessageBusMock) WriteTape(p context.Context, p1 io.Writer) (r error) {
	atomic.AddUint64(&m.WriteTapePreCounter, 1)
	defer atomic.AddUint64(&m.WriteTapeCounter, 1)

	if m.WriteTapeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.WriteTapeMock.mockExpectations, MessageBusMockWriteTapeParams{p, p1},
			"MessageBus.WriteTape got unexpected parameters")

		if m.WriteTapeFunc == nil {

			m.t.Fatal("No results are set for the MessageBusMock.WriteTape")

			return
		}
	}

	if m.WriteTapeFunc == nil {
		m.t.Fatal("Unexpected call to MessageBusMock.WriteTape")
		return
	}

	return m.WriteTapeFunc(p, p1)
}

//WriteTapeMinimockCounter returns a count of MessageBusMock.WriteTapeFunc invocations
func (m *MessageBusMock) WriteTapeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteTapeCounter)
}

//WriteTapeMinimockPreCounter returns the value of MessageBusMock.WriteTape invocations
func (m *MessageBusMock) WriteTapeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteTapePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MessageBusMock) ValidateCallCounters() {

	if m.MustRegisterFunc != nil && atomic.LoadUint64(&m.MustRegisterCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.MustRegister")
	}

	if m.NewPlayerFunc != nil && atomic.LoadUint64(&m.NewPlayerCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.NewPlayer")
	}

	if m.NewRecorderFunc != nil && atomic.LoadUint64(&m.NewRecorderCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.NewRecorder")
	}

	if m.RegisterFunc != nil && atomic.LoadUint64(&m.RegisterCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.Register")
	}

	if m.SendFunc != nil && atomic.LoadUint64(&m.SendCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.Send")
	}

	if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.WriteTape")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MessageBusMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MessageBusMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MessageBusMock) MinimockFinish() {

	if m.MustRegisterFunc != nil && atomic.LoadUint64(&m.MustRegisterCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.MustRegister")
	}

	if m.NewPlayerFunc != nil && atomic.LoadUint64(&m.NewPlayerCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.NewPlayer")
	}

	if m.NewRecorderFunc != nil && atomic.LoadUint64(&m.NewRecorderCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.NewRecorder")
	}

	if m.RegisterFunc != nil && atomic.LoadUint64(&m.RegisterCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.Register")
	}

	if m.SendFunc != nil && atomic.LoadUint64(&m.SendCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.Send")
	}

	if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
		m.t.Fatal("Expected call to MessageBusMock.WriteTape")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MessageBusMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MessageBusMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.MustRegisterFunc == nil || atomic.LoadUint64(&m.MustRegisterCounter) > 0)
		ok = ok && (m.NewPlayerFunc == nil || atomic.LoadUint64(&m.NewPlayerCounter) > 0)
		ok = ok && (m.NewRecorderFunc == nil || atomic.LoadUint64(&m.NewRecorderCounter) > 0)
		ok = ok && (m.RegisterFunc == nil || atomic.LoadUint64(&m.RegisterCounter) > 0)
		ok = ok && (m.SendFunc == nil || atomic.LoadUint64(&m.SendCounter) > 0)
		ok = ok && (m.WriteTapeFunc == nil || atomic.LoadUint64(&m.WriteTapeCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.MustRegisterFunc != nil && atomic.LoadUint64(&m.MustRegisterCounter) == 0 {
				m.t.Error("Expected call to MessageBusMock.MustRegister")
			}

			if m.NewPlayerFunc != nil && atomic.LoadUint64(&m.NewPlayerCounter) == 0 {
				m.t.Error("Expected call to MessageBusMock.NewPlayer")
			}

			if m.NewRecorderFunc != nil && atomic.LoadUint64(&m.NewRecorderCounter) == 0 {
				m.t.Error("Expected call to MessageBusMock.NewRecorder")
			}

			if m.RegisterFunc != nil && atomic.LoadUint64(&m.RegisterCounter) == 0 {
				m.t.Error("Expected call to MessageBusMock.Register")
			}

			if m.SendFunc != nil && atomic.LoadUint64(&m.SendCounter) == 0 {
				m.t.Error("Expected call to MessageBusMock.Send")
			}

			if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
				m.t.Error("Expected call to MessageBusMock.WriteTape")
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
func (m *MessageBusMock) AllMocksCalled() bool {

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

	if m.WriteTapeFunc != nil && atomic.LoadUint64(&m.WriteTapeCounter) == 0 {
		return false
	}

	return true
}

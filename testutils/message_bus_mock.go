package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MessageBus" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//MessageBusMock implements github.com/insolar/insolar/insolar.MessageBus
type MessageBusMock struct {
	t minimock.Tester

	MustRegisterFunc       func(p insolar.MessageType, p1 insolar.MessageHandler)
	MustRegisterCounter    uint64
	MustRegisterPreCounter uint64
	MustRegisterMock       mMessageBusMockMustRegister

	OnPulseFunc       func(p context.Context, p1 insolar.Pulse) (r error)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mMessageBusMockOnPulse

	RegisterFunc       func(p insolar.MessageType, p1 insolar.MessageHandler) (r error)
	RegisterCounter    uint64
	RegisterPreCounter uint64
	RegisterMock       mMessageBusMockRegister

	SendFunc       func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       mMessageBusMockSend
}

//NewMessageBusMock returns a mock for github.com/insolar/insolar/insolar.MessageBus
func NewMessageBusMock(t minimock.Tester) *MessageBusMock {
	m := &MessageBusMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.MustRegisterMock = mMessageBusMockMustRegister{mock: m}
	m.OnPulseMock = mMessageBusMockOnPulse{mock: m}
	m.RegisterMock = mMessageBusMockRegister{mock: m}
	m.SendMock = mMessageBusMockSend{mock: m}

	return m
}

type mMessageBusMockMustRegister struct {
	mock              *MessageBusMock
	mainExpectation   *MessageBusMockMustRegisterExpectation
	expectationSeries []*MessageBusMockMustRegisterExpectation
}

type MessageBusMockMustRegisterExpectation struct {
	input *MessageBusMockMustRegisterInput
}

type MessageBusMockMustRegisterInput struct {
	p  insolar.MessageType
	p1 insolar.MessageHandler
}

//Expect specifies that invocation of MessageBus.MustRegister is expected from 1 to Infinity times
func (m *mMessageBusMockMustRegister) Expect(p insolar.MessageType, p1 insolar.MessageHandler) *mMessageBusMockMustRegister {
	m.mock.MustRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockMustRegisterExpectation{}
	}
	m.mainExpectation.input = &MessageBusMockMustRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of MessageBus.MustRegister
func (m *mMessageBusMockMustRegister) Return() *MessageBusMock {
	m.mock.MustRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockMustRegisterExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of MessageBus.MustRegister is expected once
func (m *mMessageBusMockMustRegister) ExpectOnce(p insolar.MessageType, p1 insolar.MessageHandler) *MessageBusMockMustRegisterExpectation {
	m.mock.MustRegisterFunc = nil
	m.mainExpectation = nil

	expectation := &MessageBusMockMustRegisterExpectation{}
	expectation.input = &MessageBusMockMustRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of MessageBus.MustRegister method
func (m *mMessageBusMockMustRegister) Set(f func(p insolar.MessageType, p1 insolar.MessageHandler)) *MessageBusMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MustRegisterFunc = f
	return m.mock
}

//MustRegister implements github.com/insolar/insolar/insolar.MessageBus interface
func (m *MessageBusMock) MustRegister(p insolar.MessageType, p1 insolar.MessageHandler) {
	counter := atomic.AddUint64(&m.MustRegisterPreCounter, 1)
	defer atomic.AddUint64(&m.MustRegisterCounter, 1)

	if len(m.MustRegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MustRegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MessageBusMock.MustRegister. %v %v", p, p1)
			return
		}

		input := m.MustRegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MessageBusMockMustRegisterInput{p, p1}, "MessageBus.MustRegister got unexpected parameters")

		return
	}

	if m.MustRegisterMock.mainExpectation != nil {

		input := m.MustRegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MessageBusMockMustRegisterInput{p, p1}, "MessageBus.MustRegister got unexpected parameters")
		}

		return
	}

	if m.MustRegisterFunc == nil {
		m.t.Fatalf("Unexpected call to MessageBusMock.MustRegister. %v %v", p, p1)
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

//MustRegisterFinished returns true if mock invocations count is ok
func (m *MessageBusMock) MustRegisterFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MustRegisterMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MustRegisterCounter) == uint64(len(m.MustRegisterMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MustRegisterMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MustRegisterCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MustRegisterFunc != nil {
		return atomic.LoadUint64(&m.MustRegisterCounter) > 0
	}

	return true
}

type mMessageBusMockOnPulse struct {
	mock              *MessageBusMock
	mainExpectation   *MessageBusMockOnPulseExpectation
	expectationSeries []*MessageBusMockOnPulseExpectation
}

type MessageBusMockOnPulseExpectation struct {
	input  *MessageBusMockOnPulseInput
	result *MessageBusMockOnPulseResult
}

type MessageBusMockOnPulseInput struct {
	p  context.Context
	p1 insolar.Pulse
}

type MessageBusMockOnPulseResult struct {
	r error
}

//Expect specifies that invocation of MessageBus.OnPulse is expected from 1 to Infinity times
func (m *mMessageBusMockOnPulse) Expect(p context.Context, p1 insolar.Pulse) *mMessageBusMockOnPulse {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockOnPulseExpectation{}
	}
	m.mainExpectation.input = &MessageBusMockOnPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of MessageBus.OnPulse
func (m *mMessageBusMockOnPulse) Return(r error) *MessageBusMock {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockOnPulseExpectation{}
	}
	m.mainExpectation.result = &MessageBusMockOnPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MessageBus.OnPulse is expected once
func (m *mMessageBusMockOnPulse) ExpectOnce(p context.Context, p1 insolar.Pulse) *MessageBusMockOnPulseExpectation {
	m.mock.OnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &MessageBusMockOnPulseExpectation{}
	expectation.input = &MessageBusMockOnPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MessageBusMockOnPulseExpectation) Return(r error) {
	e.result = &MessageBusMockOnPulseResult{r}
}

//Set uses given function f as a mock of MessageBus.OnPulse method
func (m *mMessageBusMockOnPulse) Set(f func(p context.Context, p1 insolar.Pulse) (r error)) *MessageBusMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/insolar.MessageBus interface
func (m *MessageBusMock) OnPulse(p context.Context, p1 insolar.Pulse) (r error) {
	counter := atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if len(m.OnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MessageBusMock.OnPulse. %v %v", p, p1)
			return
		}

		input := m.OnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MessageBusMockOnPulseInput{p, p1}, "MessageBus.OnPulse got unexpected parameters")

		result := m.OnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MessageBusMock.OnPulse")
			return
		}

		r = result.r

		return
	}

	if m.OnPulseMock.mainExpectation != nil {

		input := m.OnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MessageBusMockOnPulseInput{p, p1}, "MessageBus.OnPulse got unexpected parameters")
		}

		result := m.OnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MessageBusMock.OnPulse")
		}

		r = result.r

		return
	}

	if m.OnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to MessageBusMock.OnPulse. %v %v", p, p1)
		return
	}

	return m.OnPulseFunc(p, p1)
}

//OnPulseMinimockCounter returns a count of MessageBusMock.OnPulseFunc invocations
func (m *MessageBusMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of MessageBusMock.OnPulse invocations
func (m *MessageBusMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//OnPulseFinished returns true if mock invocations count is ok
func (m *MessageBusMock) OnPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OnPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OnPulseCounter) == uint64(len(m.OnPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OnPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OnPulseFunc != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	return true
}

type mMessageBusMockRegister struct {
	mock              *MessageBusMock
	mainExpectation   *MessageBusMockRegisterExpectation
	expectationSeries []*MessageBusMockRegisterExpectation
}

type MessageBusMockRegisterExpectation struct {
	input  *MessageBusMockRegisterInput
	result *MessageBusMockRegisterResult
}

type MessageBusMockRegisterInput struct {
	p  insolar.MessageType
	p1 insolar.MessageHandler
}

type MessageBusMockRegisterResult struct {
	r error
}

//Expect specifies that invocation of MessageBus.Register is expected from 1 to Infinity times
func (m *mMessageBusMockRegister) Expect(p insolar.MessageType, p1 insolar.MessageHandler) *mMessageBusMockRegister {
	m.mock.RegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockRegisterExpectation{}
	}
	m.mainExpectation.input = &MessageBusMockRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of MessageBus.Register
func (m *mMessageBusMockRegister) Return(r error) *MessageBusMock {
	m.mock.RegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockRegisterExpectation{}
	}
	m.mainExpectation.result = &MessageBusMockRegisterResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MessageBus.Register is expected once
func (m *mMessageBusMockRegister) ExpectOnce(p insolar.MessageType, p1 insolar.MessageHandler) *MessageBusMockRegisterExpectation {
	m.mock.RegisterFunc = nil
	m.mainExpectation = nil

	expectation := &MessageBusMockRegisterExpectation{}
	expectation.input = &MessageBusMockRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MessageBusMockRegisterExpectation) Return(r error) {
	e.result = &MessageBusMockRegisterResult{r}
}

//Set uses given function f as a mock of MessageBus.Register method
func (m *mMessageBusMockRegister) Set(f func(p insolar.MessageType, p1 insolar.MessageHandler) (r error)) *MessageBusMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterFunc = f
	return m.mock
}

//Register implements github.com/insolar/insolar/insolar.MessageBus interface
func (m *MessageBusMock) Register(p insolar.MessageType, p1 insolar.MessageHandler) (r error) {
	counter := atomic.AddUint64(&m.RegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterCounter, 1)

	if len(m.RegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MessageBusMock.Register. %v %v", p, p1)
			return
		}

		input := m.RegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MessageBusMockRegisterInput{p, p1}, "MessageBus.Register got unexpected parameters")

		result := m.RegisterMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MessageBusMock.Register")
			return
		}

		r = result.r

		return
	}

	if m.RegisterMock.mainExpectation != nil {

		input := m.RegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MessageBusMockRegisterInput{p, p1}, "MessageBus.Register got unexpected parameters")
		}

		result := m.RegisterMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MessageBusMock.Register")
		}

		r = result.r

		return
	}

	if m.RegisterFunc == nil {
		m.t.Fatalf("Unexpected call to MessageBusMock.Register. %v %v", p, p1)
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

//RegisterFinished returns true if mock invocations count is ok
func (m *MessageBusMock) RegisterFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterCounter) == uint64(len(m.RegisterMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterFunc != nil {
		return atomic.LoadUint64(&m.RegisterCounter) > 0
	}

	return true
}

type mMessageBusMockSend struct {
	mock              *MessageBusMock
	mainExpectation   *MessageBusMockSendExpectation
	expectationSeries []*MessageBusMockSendExpectation
}

type MessageBusMockSendExpectation struct {
	input  *MessageBusMockSendInput
	result *MessageBusMockSendResult
}

type MessageBusMockSendInput struct {
	p  context.Context
	p1 insolar.Message
	p2 *insolar.MessageSendOptions
}

type MessageBusMockSendResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of MessageBus.Send is expected from 1 to Infinity times
func (m *mMessageBusMockSend) Expect(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) *mMessageBusMockSend {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockSendExpectation{}
	}
	m.mainExpectation.input = &MessageBusMockSendInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of MessageBus.Send
func (m *mMessageBusMockSend) Return(r insolar.Reply, r1 error) *MessageBusMock {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MessageBusMockSendExpectation{}
	}
	m.mainExpectation.result = &MessageBusMockSendResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of MessageBus.Send is expected once
func (m *mMessageBusMockSend) ExpectOnce(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) *MessageBusMockSendExpectation {
	m.mock.SendFunc = nil
	m.mainExpectation = nil

	expectation := &MessageBusMockSendExpectation{}
	expectation.input = &MessageBusMockSendInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MessageBusMockSendExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &MessageBusMockSendResult{r, r1}
}

//Set uses given function f as a mock of MessageBus.Send method
func (m *mMessageBusMockSend) Set(f func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error)) *MessageBusMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendFunc = f
	return m.mock
}

//Send implements github.com/insolar/insolar/insolar.MessageBus interface
func (m *MessageBusMock) Send(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
	counter := atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if len(m.SendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MessageBusMock.Send. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MessageBusMockSendInput{p, p1, p2}, "MessageBus.Send got unexpected parameters")

		result := m.SendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MessageBusMock.Send")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMock.mainExpectation != nil {

		input := m.SendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MessageBusMockSendInput{p, p1, p2}, "MessageBus.Send got unexpected parameters")
		}

		result := m.SendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MessageBusMock.Send")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendFunc == nil {
		m.t.Fatalf("Unexpected call to MessageBusMock.Send. %v %v %v", p, p1, p2)
		return
	}

	return m.SendFunc(p, p1, p2)
}

//SendMinimockCounter returns a count of MessageBusMock.SendFunc invocations
func (m *MessageBusMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of MessageBusMock.Send invocations
func (m *MessageBusMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

//SendFinished returns true if mock invocations count is ok
func (m *MessageBusMock) SendFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendCounter) == uint64(len(m.SendMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendFunc != nil {
		return atomic.LoadUint64(&m.SendCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MessageBusMock) ValidateCallCounters() {

	if !m.MustRegisterFinished() {
		m.t.Fatal("Expected call to MessageBusMock.MustRegister")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to MessageBusMock.OnPulse")
	}

	if !m.RegisterFinished() {
		m.t.Fatal("Expected call to MessageBusMock.Register")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to MessageBusMock.Send")
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

	if !m.MustRegisterFinished() {
		m.t.Fatal("Expected call to MessageBusMock.MustRegister")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to MessageBusMock.OnPulse")
	}

	if !m.RegisterFinished() {
		m.t.Fatal("Expected call to MessageBusMock.Register")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to MessageBusMock.Send")
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
		ok = ok && m.MustRegisterFinished()
		ok = ok && m.OnPulseFinished()
		ok = ok && m.RegisterFinished()
		ok = ok && m.SendFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.MustRegisterFinished() {
				m.t.Error("Expected call to MessageBusMock.MustRegister")
			}

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to MessageBusMock.OnPulse")
			}

			if !m.RegisterFinished() {
				m.t.Error("Expected call to MessageBusMock.Register")
			}

			if !m.SendFinished() {
				m.t.Error("Expected call to MessageBusMock.Send")
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

	if !m.MustRegisterFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	if !m.RegisterFinished() {
		return false
	}

	if !m.SendFinished() {
		return false
	}

	return true
}

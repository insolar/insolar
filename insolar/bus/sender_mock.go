package bus

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Sender" can be found in github.com/insolar/insolar/insolar/bus
*/
import (
	context "context"
	"sync/atomic"
	"time"

	message "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//SenderMock implements github.com/insolar/insolar/insolar/bus.Sender
type SenderMock struct {
	t minimock.Tester

	ReplyFunc       func(p context.Context, p1 *message.Message, p2 *message.Message)
	ReplyCounter    uint64
	ReplyPreCounter uint64
	ReplyMock       mSenderMockReply

	SendFunc       func(p context.Context, p1 *message.Message) (r <-chan *message.Message, r1 func())
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       mSenderMockSend
}

//NewSenderMock returns a mock for github.com/insolar/insolar/insolar/bus.Sender
func NewSenderMock(t minimock.Tester) *SenderMock {
	m := &SenderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ReplyMock = mSenderMockReply{mock: m}
	m.SendMock = mSenderMockSend{mock: m}

	return m
}

type mSenderMockReply struct {
	mock              *SenderMock
	mainExpectation   *SenderMockReplyExpectation
	expectationSeries []*SenderMockReplyExpectation
}

type SenderMockReplyExpectation struct {
	input *SenderMockReplyInput
}

type SenderMockReplyInput struct {
	p  context.Context
	p1 *message.Message
	p2 *message.Message
}

//Expect specifies that invocation of Sender.Reply is expected from 1 to Infinity times
func (m *mSenderMockReply) Expect(p context.Context, p1 *message.Message, p2 *message.Message) *mSenderMockReply {
	m.mock.ReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockReplyExpectation{}
	}
	m.mainExpectation.input = &SenderMockReplyInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Sender.Reply
func (m *mSenderMockReply) Return() *SenderMock {
	m.mock.ReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockReplyExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Sender.Reply is expected once
func (m *mSenderMockReply) ExpectOnce(p context.Context, p1 *message.Message, p2 *message.Message) *SenderMockReplyExpectation {
	m.mock.ReplyFunc = nil
	m.mainExpectation = nil

	expectation := &SenderMockReplyExpectation{}
	expectation.input = &SenderMockReplyInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Sender.Reply method
func (m *mSenderMockReply) Set(f func(p context.Context, p1 *message.Message, p2 *message.Message)) *SenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReplyFunc = f
	return m.mock
}

//Reply implements github.com/insolar/insolar/insolar/bus.Sender interface
func (m *SenderMock) Reply(p context.Context, p1 *message.Message, p2 *message.Message) {
	counter := atomic.AddUint64(&m.ReplyPreCounter, 1)
	defer atomic.AddUint64(&m.ReplyCounter, 1)

	if len(m.ReplyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReplyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SenderMock.Reply. %v %v %v", p, p1, p2)
			return
		}

		input := m.ReplyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SenderMockReplyInput{p, p1, p2}, "Sender.Reply got unexpected parameters")

		return
	}

	if m.ReplyMock.mainExpectation != nil {

		input := m.ReplyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SenderMockReplyInput{p, p1, p2}, "Sender.Reply got unexpected parameters")
		}

		return
	}

	if m.ReplyFunc == nil {
		m.t.Fatalf("Unexpected call to SenderMock.Reply. %v %v %v", p, p1, p2)
		return
	}

	m.ReplyFunc(p, p1, p2)
}

//ReplyMinimockCounter returns a count of SenderMock.ReplyFunc invocations
func (m *SenderMock) ReplyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReplyCounter)
}

//ReplyMinimockPreCounter returns the value of SenderMock.Reply invocations
func (m *SenderMock) ReplyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReplyPreCounter)
}

//ReplyFinished returns true if mock invocations count is ok
func (m *SenderMock) ReplyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReplyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReplyCounter) == uint64(len(m.ReplyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReplyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReplyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReplyFunc != nil {
		return atomic.LoadUint64(&m.ReplyCounter) > 0
	}

	return true
}

type mSenderMockSend struct {
	mock              *SenderMock
	mainExpectation   *SenderMockSendExpectation
	expectationSeries []*SenderMockSendExpectation
}

type SenderMockSendExpectation struct {
	input  *SenderMockSendInput
	result *SenderMockSendResult
}

type SenderMockSendInput struct {
	p  context.Context
	p1 *message.Message
}

type SenderMockSendResult struct {
	r  <-chan *message.Message
	r1 func()
}

//Expect specifies that invocation of Sender.Send is expected from 1 to Infinity times
func (m *mSenderMockSend) Expect(p context.Context, p1 *message.Message) *mSenderMockSend {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockSendExpectation{}
	}
	m.mainExpectation.input = &SenderMockSendInput{p, p1}
	return m
}

//Return specifies results of invocation of Sender.Send
func (m *mSenderMockSend) Return(r <-chan *message.Message, r1 func()) *SenderMock {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockSendExpectation{}
	}
	m.mainExpectation.result = &SenderMockSendResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Sender.Send is expected once
func (m *mSenderMockSend) ExpectOnce(p context.Context, p1 *message.Message) *SenderMockSendExpectation {
	m.mock.SendFunc = nil
	m.mainExpectation = nil

	expectation := &SenderMockSendExpectation{}
	expectation.input = &SenderMockSendInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SenderMockSendExpectation) Return(r <-chan *message.Message, r1 func()) {
	e.result = &SenderMockSendResult{r, r1}
}

//Set uses given function f as a mock of Sender.Send method
func (m *mSenderMockSend) Set(f func(p context.Context, p1 *message.Message) (r <-chan *message.Message, r1 func())) *SenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendFunc = f
	return m.mock
}

//Send implements github.com/insolar/insolar/insolar/bus.Sender interface
func (m *SenderMock) Send(p context.Context, p1 *message.Message) (r <-chan *message.Message, r1 func()) {
	counter := atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if len(m.SendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SenderMock.Send. %v %v", p, p1)
			return
		}

		input := m.SendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SenderMockSendInput{p, p1}, "Sender.Send got unexpected parameters")

		result := m.SendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SenderMock.Send")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMock.mainExpectation != nil {

		input := m.SendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SenderMockSendInput{p, p1}, "Sender.Send got unexpected parameters")
		}

		result := m.SendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SenderMock.Send")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendFunc == nil {
		m.t.Fatalf("Unexpected call to SenderMock.Send. %v %v", p, p1)
		return
	}

	return m.SendFunc(p, p1)
}

//SendMinimockCounter returns a count of SenderMock.SendFunc invocations
func (m *SenderMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of SenderMock.Send invocations
func (m *SenderMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

//SendFinished returns true if mock invocations count is ok
func (m *SenderMock) SendFinished() bool {
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
func (m *SenderMock) ValidateCallCounters() {

	if !m.ReplyFinished() {
		m.t.Fatal("Expected call to SenderMock.Reply")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to SenderMock.Send")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SenderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SenderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SenderMock) MinimockFinish() {

	if !m.ReplyFinished() {
		m.t.Fatal("Expected call to SenderMock.Reply")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to SenderMock.Send")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SenderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SenderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ReplyFinished()
		ok = ok && m.SendFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ReplyFinished() {
				m.t.Error("Expected call to SenderMock.Reply")
			}

			if !m.SendFinished() {
				m.t.Error("Expected call to SenderMock.Send")
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
func (m *SenderMock) AllMocksCalled() bool {

	if !m.ReplyFinished() {
		return false
	}

	if !m.SendFinished() {
		return false
	}

	return true
}

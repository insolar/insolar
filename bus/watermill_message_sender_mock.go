package bus

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "WatermillMessageSender" can be found in github.com/insolar/insolar/bus
*/
import (
	context "context"
	"sync/atomic"
	"time"

	message "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//WatermillMessageSenderMock implements github.com/insolar/insolar/bus.WatermillMessageSender
type WatermillMessageSenderMock struct {
	t minimock.Tester

	SendFunc       func(p context.Context, p1 *message.Message) (r <-chan *message.Message)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       mWatermillMessageSenderMockSend
}

//NewWatermillMessageSenderMock returns a mock for github.com/insolar/insolar/bus.WatermillMessageSender
func NewWatermillMessageSenderMock(t minimock.Tester) *WatermillMessageSenderMock {
	m := &WatermillMessageSenderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SendMock = mWatermillMessageSenderMockSend{mock: m}

	return m
}

type mWatermillMessageSenderMockSend struct {
	mock              *WatermillMessageSenderMock
	mainExpectation   *WatermillMessageSenderMockSendExpectation
	expectationSeries []*WatermillMessageSenderMockSendExpectation
}

type WatermillMessageSenderMockSendExpectation struct {
	input  *WatermillMessageSenderMockSendInput
	result *WatermillMessageSenderMockSendResult
}

type WatermillMessageSenderMockSendInput struct {
	p  context.Context
	p1 *message.Message
}

type WatermillMessageSenderMockSendResult struct {
	r <-chan *message.Message
}

//Expect specifies that invocation of WatermillMessageSender.Send is expected from 1 to Infinity times
func (m *mWatermillMessageSenderMockSend) Expect(p context.Context, p1 *message.Message) *mWatermillMessageSenderMockSend {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WatermillMessageSenderMockSendExpectation{}
	}
	m.mainExpectation.input = &WatermillMessageSenderMockSendInput{p, p1}
	return m
}

//Return specifies results of invocation of WatermillMessageSender.Send
func (m *mWatermillMessageSenderMockSend) Return(r <-chan *message.Message) *WatermillMessageSenderMock {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WatermillMessageSenderMockSendExpectation{}
	}
	m.mainExpectation.result = &WatermillMessageSenderMockSendResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of WatermillMessageSender.Send is expected once
func (m *mWatermillMessageSenderMockSend) ExpectOnce(p context.Context, p1 *message.Message) *WatermillMessageSenderMockSendExpectation {
	m.mock.SendFunc = nil
	m.mainExpectation = nil

	expectation := &WatermillMessageSenderMockSendExpectation{}
	expectation.input = &WatermillMessageSenderMockSendInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *WatermillMessageSenderMockSendExpectation) Return(r <-chan *message.Message) {
	e.result = &WatermillMessageSenderMockSendResult{r}
}

//Set uses given function f as a mock of WatermillMessageSender.Send method
func (m *mWatermillMessageSenderMockSend) Set(f func(p context.Context, p1 *message.Message) (r <-chan *message.Message)) *WatermillMessageSenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendFunc = f
	return m.mock
}

//Send implements github.com/insolar/insolar/bus.WatermillMessageSender interface
func (m *WatermillMessageSenderMock) Send(p context.Context, p1 *message.Message) (r <-chan *message.Message) {
	counter := atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if len(m.SendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to WatermillMessageSenderMock.Send. %v %v", p, p1)
			return
		}

		input := m.SendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, WatermillMessageSenderMockSendInput{p, p1}, "WatermillMessageSender.Send got unexpected parameters")

		result := m.SendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the WatermillMessageSenderMock.Send")
			return
		}

		r = result.r

		return
	}

	if m.SendMock.mainExpectation != nil {

		input := m.SendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, WatermillMessageSenderMockSendInput{p, p1}, "WatermillMessageSender.Send got unexpected parameters")
		}

		result := m.SendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the WatermillMessageSenderMock.Send")
		}

		r = result.r

		return
	}

	if m.SendFunc == nil {
		m.t.Fatalf("Unexpected call to WatermillMessageSenderMock.Send. %v %v", p, p1)
		return
	}

	return m.SendFunc(p, p1)
}

//SendMinimockCounter returns a count of WatermillMessageSenderMock.SendFunc invocations
func (m *WatermillMessageSenderMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of WatermillMessageSenderMock.Send invocations
func (m *WatermillMessageSenderMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

//SendFinished returns true if mock invocations count is ok
func (m *WatermillMessageSenderMock) SendFinished() bool {
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
func (m *WatermillMessageSenderMock) ValidateCallCounters() {

	if !m.SendFinished() {
		m.t.Fatal("Expected call to WatermillMessageSenderMock.Send")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *WatermillMessageSenderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *WatermillMessageSenderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *WatermillMessageSenderMock) MinimockFinish() {

	if !m.SendFinished() {
		m.t.Fatal("Expected call to WatermillMessageSenderMock.Send")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *WatermillMessageSenderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *WatermillMessageSenderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SendFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SendFinished() {
				m.t.Error("Expected call to WatermillMessageSenderMock.Send")
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
func (m *WatermillMessageSenderMock) AllMocksCalled() bool {

	if !m.SendFinished() {
		return false
	}

	return true
}

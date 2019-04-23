package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Bus" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	message "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//BusMock implements github.com/insolar/insolar/insolar.Bus
type BusMock struct {
	t minimock.Tester

	SendFunc       func(p context.Context, p1 *message.Message) (r <-chan *message.Message)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       mBusMockSend
}

//NewBusMock returns a mock for github.com/insolar/insolar/insolar.Bus
func NewBusMock(t minimock.Tester) *BusMock {
	m := &BusMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SendMock = mBusMockSend{mock: m}

	return m
}

type mBusMockSend struct {
	mock              *BusMock
	mainExpectation   *BusMockSendExpectation
	expectationSeries []*BusMockSendExpectation
}

type BusMockSendExpectation struct {
	input  *BusMockSendInput
	result *BusMockSendResult
}

type BusMockSendInput struct {
	p  context.Context
	p1 *message.Message
}

type BusMockSendResult struct {
	r <-chan *message.Message
}

//Expect specifies that invocation of Bus.Send is expected from 1 to Infinity times
func (m *mBusMockSend) Expect(p context.Context, p1 *message.Message) *mBusMockSend {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BusMockSendExpectation{}
	}
	m.mainExpectation.input = &BusMockSendInput{p, p1}
	return m
}

//Return specifies results of invocation of Bus.Send
func (m *mBusMockSend) Return(r <-chan *message.Message) *BusMock {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BusMockSendExpectation{}
	}
	m.mainExpectation.result = &BusMockSendResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Bus.Send is expected once
func (m *mBusMockSend) ExpectOnce(p context.Context, p1 *message.Message) *BusMockSendExpectation {
	m.mock.SendFunc = nil
	m.mainExpectation = nil

	expectation := &BusMockSendExpectation{}
	expectation.input = &BusMockSendInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BusMockSendExpectation) Return(r <-chan *message.Message) {
	e.result = &BusMockSendResult{r}
}

//Set uses given function f as a mock of Bus.Send method
func (m *mBusMockSend) Set(f func(p context.Context, p1 *message.Message) (r <-chan *message.Message)) *BusMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendFunc = f
	return m.mock
}

//Send implements github.com/insolar/insolar/insolar.Bus interface
func (m *BusMock) Send(p context.Context, p1 *message.Message) (r <-chan *message.Message) {
	counter := atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if len(m.SendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BusMock.Send. %v %v", p, p1)
			return
		}

		input := m.SendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, BusMockSendInput{p, p1}, "Bus.Send got unexpected parameters")

		result := m.SendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BusMock.Send")
			return
		}

		r = result.r

		return
	}

	if m.SendMock.mainExpectation != nil {

		input := m.SendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, BusMockSendInput{p, p1}, "Bus.Send got unexpected parameters")
		}

		result := m.SendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BusMock.Send")
		}

		r = result.r

		return
	}

	if m.SendFunc == nil {
		m.t.Fatalf("Unexpected call to BusMock.Send. %v %v", p, p1)
		return
	}

	return m.SendFunc(p, p1)
}

//SendMinimockCounter returns a count of BusMock.SendFunc invocations
func (m *BusMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of BusMock.Send invocations
func (m *BusMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

//SendFinished returns true if mock invocations count is ok
func (m *BusMock) SendFinished() bool {
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
func (m *BusMock) ValidateCallCounters() {

	if !m.SendFinished() {
		m.t.Fatal("Expected call to BusMock.Send")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BusMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *BusMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *BusMock) MinimockFinish() {

	if !m.SendFinished() {
		m.t.Fatal("Expected call to BusMock.Send")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *BusMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *BusMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to BusMock.Send")
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
func (m *BusMock) AllMocksCalled() bool {

	if !m.SendFinished() {
		return false
	}

	return true
}

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
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//SenderMock implements github.com/insolar/insolar/insolar/bus.Sender
type SenderMock struct {
	t minimock.Tester

	ReplyFunc       func(p context.Context, p1 *message.Message, p2 *message.Message)
	ReplyCounter    uint64
	ReplyPreCounter uint64
	ReplyMock       mSenderMockReply

	SendRoleFunc       func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func())
	SendRoleCounter    uint64
	SendRolePreCounter uint64
	SendRoleMock       mSenderMockSendRole

	SendTargetFunc       func(p context.Context, p1 *message.Message, p2 insolar.Reference) (r <-chan *message.Message, r1 func())
	SendTargetCounter    uint64
	SendTargetPreCounter uint64
	SendTargetMock       mSenderMockSendTarget
}

//NewSenderMock returns a mock for github.com/insolar/insolar/insolar/bus.Sender
func NewSenderMock(t minimock.Tester) *SenderMock {
	m := &SenderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ReplyMock = mSenderMockReply{mock: m}
	m.SendRoleMock = mSenderMockSendRole{mock: m}
	m.SendTargetMock = mSenderMockSendTarget{mock: m}

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

type mSenderMockSendRole struct {
	mock              *SenderMock
	mainExpectation   *SenderMockSendRoleExpectation
	expectationSeries []*SenderMockSendRoleExpectation
}

type SenderMockSendRoleExpectation struct {
	input  *SenderMockSendRoleInput
	result *SenderMockSendRoleResult
}

type SenderMockSendRoleInput struct {
	p  context.Context
	p1 *message.Message
	p2 insolar.DynamicRole
	p3 insolar.Reference
}

type SenderMockSendRoleResult struct {
	r  <-chan *message.Message
	r1 func()
}

//Expect specifies that invocation of Sender.SendRole is expected from 1 to Infinity times
func (m *mSenderMockSendRole) Expect(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) *mSenderMockSendRole {
	m.mock.SendRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockSendRoleExpectation{}
	}
	m.mainExpectation.input = &SenderMockSendRoleInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Sender.SendRole
func (m *mSenderMockSendRole) Return(r <-chan *message.Message, r1 func()) *SenderMock {
	m.mock.SendRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockSendRoleExpectation{}
	}
	m.mainExpectation.result = &SenderMockSendRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Sender.SendRole is expected once
func (m *mSenderMockSendRole) ExpectOnce(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) *SenderMockSendRoleExpectation {
	m.mock.SendRoleFunc = nil
	m.mainExpectation = nil

	expectation := &SenderMockSendRoleExpectation{}
	expectation.input = &SenderMockSendRoleInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SenderMockSendRoleExpectation) Return(r <-chan *message.Message, r1 func()) {
	e.result = &SenderMockSendRoleResult{r, r1}
}

//Set uses given function f as a mock of Sender.SendRole method
func (m *mSenderMockSendRole) Set(f func(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func())) *SenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendRoleFunc = f
	return m.mock
}

//SendRole implements github.com/insolar/insolar/insolar/bus.Sender interface
func (m *SenderMock) SendRole(p context.Context, p1 *message.Message, p2 insolar.DynamicRole, p3 insolar.Reference) (r <-chan *message.Message, r1 func()) {
	counter := atomic.AddUint64(&m.SendRolePreCounter, 1)
	defer atomic.AddUint64(&m.SendRoleCounter, 1)

	if len(m.SendRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SenderMock.SendRole. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SenderMockSendRoleInput{p, p1, p2, p3}, "Sender.SendRole got unexpected parameters")

		result := m.SendRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SenderMock.SendRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRoleMock.mainExpectation != nil {

		input := m.SendRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SenderMockSendRoleInput{p, p1, p2, p3}, "Sender.SendRole got unexpected parameters")
		}

		result := m.SendRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SenderMock.SendRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendRoleFunc == nil {
		m.t.Fatalf("Unexpected call to SenderMock.SendRole. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SendRoleFunc(p, p1, p2, p3)
}

//SendRoleMinimockCounter returns a count of SenderMock.SendRoleFunc invocations
func (m *SenderMock) SendRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendRoleCounter)
}

//SendRoleMinimockPreCounter returns the value of SenderMock.SendRole invocations
func (m *SenderMock) SendRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendRolePreCounter)
}

//SendRoleFinished returns true if mock invocations count is ok
func (m *SenderMock) SendRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendRoleCounter) == uint64(len(m.SendRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendRoleFunc != nil {
		return atomic.LoadUint64(&m.SendRoleCounter) > 0
	}

	return true
}

type mSenderMockSendTarget struct {
	mock              *SenderMock
	mainExpectation   *SenderMockSendTargetExpectation
	expectationSeries []*SenderMockSendTargetExpectation
}

type SenderMockSendTargetExpectation struct {
	input  *SenderMockSendTargetInput
	result *SenderMockSendTargetResult
}

type SenderMockSendTargetInput struct {
	p  context.Context
	p1 *message.Message
	p2 insolar.Reference
}

type SenderMockSendTargetResult struct {
	r  <-chan *message.Message
	r1 func()
}

//Expect specifies that invocation of Sender.SendTarget is expected from 1 to Infinity times
func (m *mSenderMockSendTarget) Expect(p context.Context, p1 *message.Message, p2 insolar.Reference) *mSenderMockSendTarget {
	m.mock.SendTargetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockSendTargetExpectation{}
	}
	m.mainExpectation.input = &SenderMockSendTargetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Sender.SendTarget
func (m *mSenderMockSendTarget) Return(r <-chan *message.Message, r1 func()) *SenderMock {
	m.mock.SendTargetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SenderMockSendTargetExpectation{}
	}
	m.mainExpectation.result = &SenderMockSendTargetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Sender.SendTarget is expected once
func (m *mSenderMockSendTarget) ExpectOnce(p context.Context, p1 *message.Message, p2 insolar.Reference) *SenderMockSendTargetExpectation {
	m.mock.SendTargetFunc = nil
	m.mainExpectation = nil

	expectation := &SenderMockSendTargetExpectation{}
	expectation.input = &SenderMockSendTargetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SenderMockSendTargetExpectation) Return(r <-chan *message.Message, r1 func()) {
	e.result = &SenderMockSendTargetResult{r, r1}
}

//Set uses given function f as a mock of Sender.SendTarget method
func (m *mSenderMockSendTarget) Set(f func(p context.Context, p1 *message.Message, p2 insolar.Reference) (r <-chan *message.Message, r1 func())) *SenderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendTargetFunc = f
	return m.mock
}

//SendTarget implements github.com/insolar/insolar/insolar/bus.Sender interface
func (m *SenderMock) SendTarget(p context.Context, p1 *message.Message, p2 insolar.Reference) (r <-chan *message.Message, r1 func()) {
	counter := atomic.AddUint64(&m.SendTargetPreCounter, 1)
	defer atomic.AddUint64(&m.SendTargetCounter, 1)

	if len(m.SendTargetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendTargetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SenderMock.SendTarget. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendTargetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SenderMockSendTargetInput{p, p1, p2}, "Sender.SendTarget got unexpected parameters")

		result := m.SendTargetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SenderMock.SendTarget")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendTargetMock.mainExpectation != nil {

		input := m.SendTargetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SenderMockSendTargetInput{p, p1, p2}, "Sender.SendTarget got unexpected parameters")
		}

		result := m.SendTargetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SenderMock.SendTarget")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendTargetFunc == nil {
		m.t.Fatalf("Unexpected call to SenderMock.SendTarget. %v %v %v", p, p1, p2)
		return
	}

	return m.SendTargetFunc(p, p1, p2)
}

//SendTargetMinimockCounter returns a count of SenderMock.SendTargetFunc invocations
func (m *SenderMock) SendTargetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendTargetCounter)
}

//SendTargetMinimockPreCounter returns the value of SenderMock.SendTarget invocations
func (m *SenderMock) SendTargetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendTargetPreCounter)
}

//SendTargetFinished returns true if mock invocations count is ok
func (m *SenderMock) SendTargetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendTargetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendTargetCounter) == uint64(len(m.SendTargetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendTargetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendTargetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendTargetFunc != nil {
		return atomic.LoadUint64(&m.SendTargetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SenderMock) ValidateCallCounters() {

	if !m.ReplyFinished() {
		m.t.Fatal("Expected call to SenderMock.Reply")
	}

	if !m.SendRoleFinished() {
		m.t.Fatal("Expected call to SenderMock.SendRole")
	}

	if !m.SendTargetFinished() {
		m.t.Fatal("Expected call to SenderMock.SendTarget")
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

	if !m.SendRoleFinished() {
		m.t.Fatal("Expected call to SenderMock.SendRole")
	}

	if !m.SendTargetFinished() {
		m.t.Fatal("Expected call to SenderMock.SendTarget")
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
		ok = ok && m.SendRoleFinished()
		ok = ok && m.SendTargetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ReplyFinished() {
				m.t.Error("Expected call to SenderMock.Reply")
			}

			if !m.SendRoleFinished() {
				m.t.Error("Expected call to SenderMock.SendRole")
			}

			if !m.SendTargetFinished() {
				m.t.Error("Expected call to SenderMock.SendTarget")
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

	if !m.SendRoleFinished() {
		return false
	}

	if !m.SendTargetFinished() {
		return false
	}

	return true
}

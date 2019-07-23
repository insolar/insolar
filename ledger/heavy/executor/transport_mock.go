package replica

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Transport" can be found in github.com/insolar/insolar/ledger/heavy/replica
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//TransportMock implements github.com/insolar/insolar/ledger/heavy/replica.Transport
type TransportMock struct {
	t minimock.Tester

	MeFunc       func() (r string)
	MeCounter    uint64
	MePreCounter uint64
	MeMock       mTransportMockMe

	RegisterFunc       func(p string, p1 Handle)
	RegisterCounter    uint64
	RegisterPreCounter uint64
	RegisterMock       mTransportMockRegister

	SendFunc       func(p context.Context, p1 string, p2 string, p3 []byte) (r []byte, r1 error)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       mTransportMockSend
}

//NewTransportMock returns a mock for github.com/insolar/insolar/ledger/heavy/replica.Transport
func NewTransportMock(t minimock.Tester) *TransportMock {
	m := &TransportMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.MeMock = mTransportMockMe{mock: m}
	m.RegisterMock = mTransportMockRegister{mock: m}
	m.SendMock = mTransportMockSend{mock: m}

	return m
}

type mTransportMockMe struct {
	mock              *TransportMock
	mainExpectation   *TransportMockMeExpectation
	expectationSeries []*TransportMockMeExpectation
}

type TransportMockMeExpectation struct {
	result *TransportMockMeResult
}

type TransportMockMeResult struct {
	r string
}

//Expect specifies that invocation of Transport.Me is expected from 1 to Infinity times
func (m *mTransportMockMe) Expect() *mTransportMockMe {
	m.mock.MeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TransportMockMeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Transport.Me
func (m *mTransportMockMe) Return(r string) *TransportMock {
	m.mock.MeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TransportMockMeExpectation{}
	}
	m.mainExpectation.result = &TransportMockMeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Transport.Me is expected once
func (m *mTransportMockMe) ExpectOnce() *TransportMockMeExpectation {
	m.mock.MeFunc = nil
	m.mainExpectation = nil

	expectation := &TransportMockMeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TransportMockMeExpectation) Return(r string) {
	e.result = &TransportMockMeResult{r}
}

//Set uses given function f as a mock of Transport.Me method
func (m *mTransportMockMe) Set(f func() (r string)) *TransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MeFunc = f
	return m.mock
}

//Me implements github.com/insolar/insolar/ledger/heavy/replica.Transport interface
func (m *TransportMock) Me() (r string) {
	counter := atomic.AddUint64(&m.MePreCounter, 1)
	defer atomic.AddUint64(&m.MeCounter, 1)

	if len(m.MeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TransportMock.Me.")
			return
		}

		result := m.MeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TransportMock.Me")
			return
		}

		r = result.r

		return
	}

	if m.MeMock.mainExpectation != nil {

		result := m.MeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TransportMock.Me")
		}

		r = result.r

		return
	}

	if m.MeFunc == nil {
		m.t.Fatalf("Unexpected call to TransportMock.Me.")
		return
	}

	return m.MeFunc()
}

//MeMinimockCounter returns a count of TransportMock.MeFunc invocations
func (m *TransportMock) MeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MeCounter)
}

//MeMinimockPreCounter returns the value of TransportMock.Me invocations
func (m *TransportMock) MeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MePreCounter)
}

//MeFinished returns true if mock invocations count is ok
func (m *TransportMock) MeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MeCounter) == uint64(len(m.MeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MeFunc != nil {
		return atomic.LoadUint64(&m.MeCounter) > 0
	}

	return true
}

type mTransportMockRegister struct {
	mock              *TransportMock
	mainExpectation   *TransportMockRegisterExpectation
	expectationSeries []*TransportMockRegisterExpectation
}

type TransportMockRegisterExpectation struct {
	input *TransportMockRegisterInput
}

type TransportMockRegisterInput struct {
	p  string
	p1 Handle
}

//Expect specifies that invocation of Transport.Register is expected from 1 to Infinity times
func (m *mTransportMockRegister) Expect(p string, p1 Handle) *mTransportMockRegister {
	m.mock.RegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TransportMockRegisterExpectation{}
	}
	m.mainExpectation.input = &TransportMockRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of Transport.Register
func (m *mTransportMockRegister) Return() *TransportMock {
	m.mock.RegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TransportMockRegisterExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Transport.Register is expected once
func (m *mTransportMockRegister) ExpectOnce(p string, p1 Handle) *TransportMockRegisterExpectation {
	m.mock.RegisterFunc = nil
	m.mainExpectation = nil

	expectation := &TransportMockRegisterExpectation{}
	expectation.input = &TransportMockRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Transport.Register method
func (m *mTransportMockRegister) Set(f func(p string, p1 Handle)) *TransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterFunc = f
	return m.mock
}

//Register implements github.com/insolar/insolar/ledger/heavy/replica.Transport interface
func (m *TransportMock) Register(p string, p1 Handle) {
	counter := atomic.AddUint64(&m.RegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterCounter, 1)

	if len(m.RegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TransportMock.Register. %v %v", p, p1)
			return
		}

		input := m.RegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TransportMockRegisterInput{p, p1}, "Transport.Register got unexpected parameters")

		return
	}

	if m.RegisterMock.mainExpectation != nil {

		input := m.RegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TransportMockRegisterInput{p, p1}, "Transport.Register got unexpected parameters")
		}

		return
	}

	if m.RegisterFunc == nil {
		m.t.Fatalf("Unexpected call to TransportMock.Register. %v %v", p, p1)
		return
	}

	m.RegisterFunc(p, p1)
}

//RegisterMinimockCounter returns a count of TransportMock.RegisterFunc invocations
func (m *TransportMock) RegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterCounter)
}

//RegisterMinimockPreCounter returns the value of TransportMock.Register invocations
func (m *TransportMock) RegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterPreCounter)
}

//RegisterFinished returns true if mock invocations count is ok
func (m *TransportMock) RegisterFinished() bool {
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

type mTransportMockSend struct {
	mock              *TransportMock
	mainExpectation   *TransportMockSendExpectation
	expectationSeries []*TransportMockSendExpectation
}

type TransportMockSendExpectation struct {
	input  *TransportMockSendInput
	result *TransportMockSendResult
}

type TransportMockSendInput struct {
	p  context.Context
	p1 string
	p2 string
	p3 []byte
}

type TransportMockSendResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of Transport.Send is expected from 1 to Infinity times
func (m *mTransportMockSend) Expect(p context.Context, p1 string, p2 string, p3 []byte) *mTransportMockSend {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TransportMockSendExpectation{}
	}
	m.mainExpectation.input = &TransportMockSendInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Transport.Send
func (m *mTransportMockSend) Return(r []byte, r1 error) *TransportMock {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TransportMockSendExpectation{}
	}
	m.mainExpectation.result = &TransportMockSendResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Transport.Send is expected once
func (m *mTransportMockSend) ExpectOnce(p context.Context, p1 string, p2 string, p3 []byte) *TransportMockSendExpectation {
	m.mock.SendFunc = nil
	m.mainExpectation = nil

	expectation := &TransportMockSendExpectation{}
	expectation.input = &TransportMockSendInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TransportMockSendExpectation) Return(r []byte, r1 error) {
	e.result = &TransportMockSendResult{r, r1}
}

//Set uses given function f as a mock of Transport.Send method
func (m *mTransportMockSend) Set(f func(p context.Context, p1 string, p2 string, p3 []byte) (r []byte, r1 error)) *TransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendFunc = f
	return m.mock
}

//Send implements github.com/insolar/insolar/ledger/heavy/replica.Transport interface
func (m *TransportMock) Send(p context.Context, p1 string, p2 string, p3 []byte) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if len(m.SendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TransportMock.Send. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TransportMockSendInput{p, p1, p2, p3}, "Transport.Send got unexpected parameters")

		result := m.SendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TransportMock.Send")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMock.mainExpectation != nil {

		input := m.SendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TransportMockSendInput{p, p1, p2, p3}, "Transport.Send got unexpected parameters")
		}

		result := m.SendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TransportMock.Send")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendFunc == nil {
		m.t.Fatalf("Unexpected call to TransportMock.Send. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SendFunc(p, p1, p2, p3)
}

//SendMinimockCounter returns a count of TransportMock.SendFunc invocations
func (m *TransportMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of TransportMock.Send invocations
func (m *TransportMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

//SendFinished returns true if mock invocations count is ok
func (m *TransportMock) SendFinished() bool {
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
func (m *TransportMock) ValidateCallCounters() {

	if !m.MeFinished() {
		m.t.Fatal("Expected call to TransportMock.Me")
	}

	if !m.RegisterFinished() {
		m.t.Fatal("Expected call to TransportMock.Register")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to TransportMock.Send")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TransportMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TransportMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TransportMock) MinimockFinish() {

	if !m.MeFinished() {
		m.t.Fatal("Expected call to TransportMock.Me")
	}

	if !m.RegisterFinished() {
		m.t.Fatal("Expected call to TransportMock.Register")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to TransportMock.Send")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TransportMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TransportMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.MeFinished()
		ok = ok && m.RegisterFinished()
		ok = ok && m.SendFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.MeFinished() {
				m.t.Error("Expected call to TransportMock.Me")
			}

			if !m.RegisterFinished() {
				m.t.Error("Expected call to TransportMock.Register")
			}

			if !m.SendFinished() {
				m.t.Error("Expected call to TransportMock.Send")
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
func (m *TransportMock) AllMocksCalled() bool {

	if !m.MeFinished() {
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

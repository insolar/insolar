package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StreamTransport" can be found in github.com/insolar/insolar/network/transport
*/
import (
	context "context"
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//StreamTransportMock implements github.com/insolar/insolar/network/transport.StreamTransport
type StreamTransportMock struct {
	t minimock.Tester

	AddressFunc       func() (r string)
	AddressCounter    uint64
	AddressPreCounter uint64
	AddressMock       mStreamTransportMockAddress

	DialFunc       func(p context.Context, p1 string) (r io.ReadWriteCloser, r1 error)
	DialCounter    uint64
	DialPreCounter uint64
	DialMock       mStreamTransportMockDial

	StartFunc       func(p context.Context) (r error)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mStreamTransportMockStart

	StopFunc       func(p context.Context) (r error)
	StopCounter    uint64
	StopPreCounter uint64
	StopMock       mStreamTransportMockStop
}

//NewStreamTransportMock returns a mock for github.com/insolar/insolar/network/transport.StreamTransport
func NewStreamTransportMock(t minimock.Tester) *StreamTransportMock {
	m := &StreamTransportMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddressMock = mStreamTransportMockAddress{mock: m}
	m.DialMock = mStreamTransportMockDial{mock: m}
	m.StartMock = mStreamTransportMockStart{mock: m}
	m.StopMock = mStreamTransportMockStop{mock: m}

	return m
}

type mStreamTransportMockAddress struct {
	mock              *StreamTransportMock
	mainExpectation   *StreamTransportMockAddressExpectation
	expectationSeries []*StreamTransportMockAddressExpectation
}

type StreamTransportMockAddressExpectation struct {
	result *StreamTransportMockAddressResult
}

type StreamTransportMockAddressResult struct {
	r string
}

//Expect specifies that invocation of StreamTransport.Address is expected from 1 to Infinity times
func (m *mStreamTransportMockAddress) Expect() *mStreamTransportMockAddress {
	m.mock.AddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of StreamTransport.Address
func (m *mStreamTransportMockAddress) Return(r string) *StreamTransportMock {
	m.mock.AddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockAddressExpectation{}
	}
	m.mainExpectation.result = &StreamTransportMockAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StreamTransport.Address is expected once
func (m *mStreamTransportMockAddress) ExpectOnce() *StreamTransportMockAddressExpectation {
	m.mock.AddressFunc = nil
	m.mainExpectation = nil

	expectation := &StreamTransportMockAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StreamTransportMockAddressExpectation) Return(r string) {
	e.result = &StreamTransportMockAddressResult{r}
}

//Set uses given function f as a mock of StreamTransport.Address method
func (m *mStreamTransportMockAddress) Set(f func() (r string)) *StreamTransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddressFunc = f
	return m.mock
}

//Address implements github.com/insolar/insolar/network/transport.StreamTransport interface
func (m *StreamTransportMock) Address() (r string) {
	counter := atomic.AddUint64(&m.AddressPreCounter, 1)
	defer atomic.AddUint64(&m.AddressCounter, 1)

	if len(m.AddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StreamTransportMock.Address.")
			return
		}

		result := m.AddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Address")
			return
		}

		r = result.r

		return
	}

	if m.AddressMock.mainExpectation != nil {

		result := m.AddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Address")
		}

		r = result.r

		return
	}

	if m.AddressFunc == nil {
		m.t.Fatalf("Unexpected call to StreamTransportMock.Address.")
		return
	}

	return m.AddressFunc()
}

//AddressMinimockCounter returns a count of StreamTransportMock.AddressFunc invocations
func (m *StreamTransportMock) AddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddressCounter)
}

//AddressMinimockPreCounter returns the value of StreamTransportMock.Address invocations
func (m *StreamTransportMock) AddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddressPreCounter)
}

//AddressFinished returns true if mock invocations count is ok
func (m *StreamTransportMock) AddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddressCounter) == uint64(len(m.AddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddressFunc != nil {
		return atomic.LoadUint64(&m.AddressCounter) > 0
	}

	return true
}

type mStreamTransportMockDial struct {
	mock              *StreamTransportMock
	mainExpectation   *StreamTransportMockDialExpectation
	expectationSeries []*StreamTransportMockDialExpectation
}

type StreamTransportMockDialExpectation struct {
	input  *StreamTransportMockDialInput
	result *StreamTransportMockDialResult
}

type StreamTransportMockDialInput struct {
	p  context.Context
	p1 string
}

type StreamTransportMockDialResult struct {
	r  io.ReadWriteCloser
	r1 error
}

//Expect specifies that invocation of StreamTransport.Dial is expected from 1 to Infinity times
func (m *mStreamTransportMockDial) Expect(p context.Context, p1 string) *mStreamTransportMockDial {
	m.mock.DialFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockDialExpectation{}
	}
	m.mainExpectation.input = &StreamTransportMockDialInput{p, p1}
	return m
}

//Return specifies results of invocation of StreamTransport.Dial
func (m *mStreamTransportMockDial) Return(r io.ReadWriteCloser, r1 error) *StreamTransportMock {
	m.mock.DialFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockDialExpectation{}
	}
	m.mainExpectation.result = &StreamTransportMockDialResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of StreamTransport.Dial is expected once
func (m *mStreamTransportMockDial) ExpectOnce(p context.Context, p1 string) *StreamTransportMockDialExpectation {
	m.mock.DialFunc = nil
	m.mainExpectation = nil

	expectation := &StreamTransportMockDialExpectation{}
	expectation.input = &StreamTransportMockDialInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StreamTransportMockDialExpectation) Return(r io.ReadWriteCloser, r1 error) {
	e.result = &StreamTransportMockDialResult{r, r1}
}

//Set uses given function f as a mock of StreamTransport.Dial method
func (m *mStreamTransportMockDial) Set(f func(p context.Context, p1 string) (r io.ReadWriteCloser, r1 error)) *StreamTransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DialFunc = f
	return m.mock
}

//Dial implements github.com/insolar/insolar/network/transport.StreamTransport interface
func (m *StreamTransportMock) Dial(p context.Context, p1 string) (r io.ReadWriteCloser, r1 error) {
	counter := atomic.AddUint64(&m.DialPreCounter, 1)
	defer atomic.AddUint64(&m.DialCounter, 1)

	if len(m.DialMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DialMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StreamTransportMock.Dial. %v %v", p, p1)
			return
		}

		input := m.DialMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StreamTransportMockDialInput{p, p1}, "StreamTransport.Dial got unexpected parameters")

		result := m.DialMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Dial")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DialMock.mainExpectation != nil {

		input := m.DialMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StreamTransportMockDialInput{p, p1}, "StreamTransport.Dial got unexpected parameters")
		}

		result := m.DialMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Dial")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DialFunc == nil {
		m.t.Fatalf("Unexpected call to StreamTransportMock.Dial. %v %v", p, p1)
		return
	}

	return m.DialFunc(p, p1)
}

//DialMinimockCounter returns a count of StreamTransportMock.DialFunc invocations
func (m *StreamTransportMock) DialMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DialCounter)
}

//DialMinimockPreCounter returns the value of StreamTransportMock.Dial invocations
func (m *StreamTransportMock) DialMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DialPreCounter)
}

//DialFinished returns true if mock invocations count is ok
func (m *StreamTransportMock) DialFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DialMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DialCounter) == uint64(len(m.DialMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DialMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DialCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DialFunc != nil {
		return atomic.LoadUint64(&m.DialCounter) > 0
	}

	return true
}

type mStreamTransportMockStart struct {
	mock              *StreamTransportMock
	mainExpectation   *StreamTransportMockStartExpectation
	expectationSeries []*StreamTransportMockStartExpectation
}

type StreamTransportMockStartExpectation struct {
	input  *StreamTransportMockStartInput
	result *StreamTransportMockStartResult
}

type StreamTransportMockStartInput struct {
	p context.Context
}

type StreamTransportMockStartResult struct {
	r error
}

//Expect specifies that invocation of StreamTransport.Start is expected from 1 to Infinity times
func (m *mStreamTransportMockStart) Expect(p context.Context) *mStreamTransportMockStart {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockStartExpectation{}
	}
	m.mainExpectation.input = &StreamTransportMockStartInput{p}
	return m
}

//Return specifies results of invocation of StreamTransport.Start
func (m *mStreamTransportMockStart) Return(r error) *StreamTransportMock {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockStartExpectation{}
	}
	m.mainExpectation.result = &StreamTransportMockStartResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StreamTransport.Start is expected once
func (m *mStreamTransportMockStart) ExpectOnce(p context.Context) *StreamTransportMockStartExpectation {
	m.mock.StartFunc = nil
	m.mainExpectation = nil

	expectation := &StreamTransportMockStartExpectation{}
	expectation.input = &StreamTransportMockStartInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StreamTransportMockStartExpectation) Return(r error) {
	e.result = &StreamTransportMockStartResult{r}
}

//Set uses given function f as a mock of StreamTransport.Start method
func (m *mStreamTransportMockStart) Set(f func(p context.Context) (r error)) *StreamTransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StartFunc = f
	return m.mock
}

//Start implements github.com/insolar/insolar/network/transport.StreamTransport interface
func (m *StreamTransportMock) Start(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if len(m.StartMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StartMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StreamTransportMock.Start. %v", p)
			return
		}

		input := m.StartMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StreamTransportMockStartInput{p}, "StreamTransport.Start got unexpected parameters")

		result := m.StartMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Start")
			return
		}

		r = result.r

		return
	}

	if m.StartMock.mainExpectation != nil {

		input := m.StartMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StreamTransportMockStartInput{p}, "StreamTransport.Start got unexpected parameters")
		}

		result := m.StartMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Start")
		}

		r = result.r

		return
	}

	if m.StartFunc == nil {
		m.t.Fatalf("Unexpected call to StreamTransportMock.Start. %v", p)
		return
	}

	return m.StartFunc(p)
}

//StartMinimockCounter returns a count of StreamTransportMock.StartFunc invocations
func (m *StreamTransportMock) StartMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StartCounter)
}

//StartMinimockPreCounter returns the value of StreamTransportMock.Start invocations
func (m *StreamTransportMock) StartMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StartPreCounter)
}

//StartFinished returns true if mock invocations count is ok
func (m *StreamTransportMock) StartFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StartMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StartCounter) == uint64(len(m.StartMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StartMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StartFunc != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	return true
}

type mStreamTransportMockStop struct {
	mock              *StreamTransportMock
	mainExpectation   *StreamTransportMockStopExpectation
	expectationSeries []*StreamTransportMockStopExpectation
}

type StreamTransportMockStopExpectation struct {
	input  *StreamTransportMockStopInput
	result *StreamTransportMockStopResult
}

type StreamTransportMockStopInput struct {
	p context.Context
}

type StreamTransportMockStopResult struct {
	r error
}

//Expect specifies that invocation of StreamTransport.Stop is expected from 1 to Infinity times
func (m *mStreamTransportMockStop) Expect(p context.Context) *mStreamTransportMockStop {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockStopExpectation{}
	}
	m.mainExpectation.input = &StreamTransportMockStopInput{p}
	return m
}

//Return specifies results of invocation of StreamTransport.Stop
func (m *mStreamTransportMockStop) Return(r error) *StreamTransportMock {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockStopExpectation{}
	}
	m.mainExpectation.result = &StreamTransportMockStopResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StreamTransport.Stop is expected once
func (m *mStreamTransportMockStop) ExpectOnce(p context.Context) *StreamTransportMockStopExpectation {
	m.mock.StopFunc = nil
	m.mainExpectation = nil

	expectation := &StreamTransportMockStopExpectation{}
	expectation.input = &StreamTransportMockStopInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StreamTransportMockStopExpectation) Return(r error) {
	e.result = &StreamTransportMockStopResult{r}
}

//Set uses given function f as a mock of StreamTransport.Stop method
func (m *mStreamTransportMockStop) Set(f func(p context.Context) (r error)) *StreamTransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StopFunc = f
	return m.mock
}

//Stop implements github.com/insolar/insolar/network/transport.StreamTransport interface
func (m *StreamTransportMock) Stop(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if len(m.StopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StreamTransportMock.Stop. %v", p)
			return
		}

		input := m.StopMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StreamTransportMockStopInput{p}, "StreamTransport.Stop got unexpected parameters")

		result := m.StopMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Stop")
			return
		}

		r = result.r

		return
	}

	if m.StopMock.mainExpectation != nil {

		input := m.StopMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StreamTransportMockStopInput{p}, "StreamTransport.Stop got unexpected parameters")
		}

		result := m.StopMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StreamTransportMock.Stop")
		}

		r = result.r

		return
	}

	if m.StopFunc == nil {
		m.t.Fatalf("Unexpected call to StreamTransportMock.Stop. %v", p)
		return
	}

	return m.StopFunc(p)
}

//StopMinimockCounter returns a count of StreamTransportMock.StopFunc invocations
func (m *StreamTransportMock) StopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StopCounter)
}

//StopMinimockPreCounter returns the value of StreamTransportMock.Stop invocations
func (m *StreamTransportMock) StopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StopPreCounter)
}

//StopFinished returns true if mock invocations count is ok
func (m *StreamTransportMock) StopFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StopMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StopCounter) == uint64(len(m.StopMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StopMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StopCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StopFunc != nil {
		return atomic.LoadUint64(&m.StopCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StreamTransportMock) ValidateCallCounters() {

	if !m.AddressFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Address")
	}

	if !m.DialFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Dial")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Stop")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StreamTransportMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StreamTransportMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StreamTransportMock) MinimockFinish() {

	if !m.AddressFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Address")
	}

	if !m.DialFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Dial")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Stop")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StreamTransportMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StreamTransportMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddressFinished()
		ok = ok && m.DialFinished()
		ok = ok && m.StartFinished()
		ok = ok && m.StopFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddressFinished() {
				m.t.Error("Expected call to StreamTransportMock.Address")
			}

			if !m.DialFinished() {
				m.t.Error("Expected call to StreamTransportMock.Dial")
			}

			if !m.StartFinished() {
				m.t.Error("Expected call to StreamTransportMock.Start")
			}

			if !m.StopFinished() {
				m.t.Error("Expected call to StreamTransportMock.Stop")
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
func (m *StreamTransportMock) AllMocksCalled() bool {

	if !m.AddressFinished() {
		return false
	}

	if !m.DialFinished() {
		return false
	}

	if !m.StartFinished() {
		return false
	}

	if !m.StopFinished() {
		return false
	}

	return true
}

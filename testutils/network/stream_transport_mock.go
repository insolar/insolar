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
	transport "github.com/insolar/insolar/network/transport"

	testify_assert "github.com/stretchr/testify/assert"
)

//StreamTransportMock implements github.com/insolar/insolar/network/transport.StreamTransport
type StreamTransportMock struct {
	t minimock.Tester

	DialFunc       func(p context.Context, p1 string) (r io.ReadWriteCloser, r1 error)
	DialCounter    uint64
	DialPreCounter uint64
	DialMock       mStreamTransportMockDial

	SetStreamHandlerFunc       func(p transport.StreamHandler)
	SetStreamHandlerCounter    uint64
	SetStreamHandlerPreCounter uint64
	SetStreamHandlerMock       mStreamTransportMockSetStreamHandler

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

	m.DialMock = mStreamTransportMockDial{mock: m}
	m.SetStreamHandlerMock = mStreamTransportMockSetStreamHandler{mock: m}
	m.StartMock = mStreamTransportMockStart{mock: m}
	m.StopMock = mStreamTransportMockStop{mock: m}

	return m
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

type mStreamTransportMockSetStreamHandler struct {
	mock              *StreamTransportMock
	mainExpectation   *StreamTransportMockSetStreamHandlerExpectation
	expectationSeries []*StreamTransportMockSetStreamHandlerExpectation
}

type StreamTransportMockSetStreamHandlerExpectation struct {
	input *StreamTransportMockSetStreamHandlerInput
}

type StreamTransportMockSetStreamHandlerInput struct {
	p transport.StreamHandler
}

//Expect specifies that invocation of StreamTransport.SetStreamHandler is expected from 1 to Infinity times
func (m *mStreamTransportMockSetStreamHandler) Expect(p transport.StreamHandler) *mStreamTransportMockSetStreamHandler {
	m.mock.SetStreamHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockSetStreamHandlerExpectation{}
	}
	m.mainExpectation.input = &StreamTransportMockSetStreamHandlerInput{p}
	return m
}

//Return specifies results of invocation of StreamTransport.SetStreamHandler
func (m *mStreamTransportMockSetStreamHandler) Return() *StreamTransportMock {
	m.mock.SetStreamHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StreamTransportMockSetStreamHandlerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of StreamTransport.SetStreamHandler is expected once
func (m *mStreamTransportMockSetStreamHandler) ExpectOnce(p transport.StreamHandler) *StreamTransportMockSetStreamHandlerExpectation {
	m.mock.SetStreamHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StreamTransportMockSetStreamHandlerExpectation{}
	expectation.input = &StreamTransportMockSetStreamHandlerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of StreamTransport.SetStreamHandler method
func (m *mStreamTransportMockSetStreamHandler) Set(f func(p transport.StreamHandler)) *StreamTransportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetStreamHandlerFunc = f
	return m.mock
}

//SetStreamHandler implements github.com/insolar/insolar/network/transport.StreamTransport interface
func (m *StreamTransportMock) SetStreamHandler(p transport.StreamHandler) {
	counter := atomic.AddUint64(&m.SetStreamHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.SetStreamHandlerCounter, 1)

	if len(m.SetStreamHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetStreamHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StreamTransportMock.SetStreamHandler. %v", p)
			return
		}

		input := m.SetStreamHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StreamTransportMockSetStreamHandlerInput{p}, "StreamTransport.SetStreamHandler got unexpected parameters")

		return
	}

	if m.SetStreamHandlerMock.mainExpectation != nil {

		input := m.SetStreamHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StreamTransportMockSetStreamHandlerInput{p}, "StreamTransport.SetStreamHandler got unexpected parameters")
		}

		return
	}

	if m.SetStreamHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StreamTransportMock.SetStreamHandler. %v", p)
		return
	}

	m.SetStreamHandlerFunc(p)
}

//SetStreamHandlerMinimockCounter returns a count of StreamTransportMock.SetStreamHandlerFunc invocations
func (m *StreamTransportMock) SetStreamHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetStreamHandlerCounter)
}

//SetStreamHandlerMinimockPreCounter returns the value of StreamTransportMock.SetStreamHandler invocations
func (m *StreamTransportMock) SetStreamHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetStreamHandlerPreCounter)
}

//SetStreamHandlerFinished returns true if mock invocations count is ok
func (m *StreamTransportMock) SetStreamHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetStreamHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetStreamHandlerCounter) == uint64(len(m.SetStreamHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetStreamHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetStreamHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetStreamHandlerFunc != nil {
		return atomic.LoadUint64(&m.SetStreamHandlerCounter) > 0
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

	if !m.DialFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Dial")
	}

	if !m.SetStreamHandlerFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.SetStreamHandler")
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

	if !m.DialFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.Dial")
	}

	if !m.SetStreamHandlerFinished() {
		m.t.Fatal("Expected call to StreamTransportMock.SetStreamHandler")
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
		ok = ok && m.DialFinished()
		ok = ok && m.SetStreamHandlerFinished()
		ok = ok && m.StartFinished()
		ok = ok && m.StopFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DialFinished() {
				m.t.Error("Expected call to StreamTransportMock.Dial")
			}

			if !m.SetStreamHandlerFinished() {
				m.t.Error("Expected call to StreamTransportMock.SetStreamHandler")
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

	if !m.DialFinished() {
		return false
	}

	if !m.SetStreamHandlerFinished() {
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

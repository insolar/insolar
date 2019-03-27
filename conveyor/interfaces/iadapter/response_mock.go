package iadapter

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Response" can be found in github.com/insolar/insolar/conveyor/interfaces/iadapter
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	adapterid "github.com/insolar/insolar/conveyor/adapter/adapterid"
)

//ResponseMock implements github.com/insolar/insolar/conveyor/interfaces/iadapter.Response
type ResponseMock struct {
	t minimock.Tester

	GetAdapterIDFunc       func() (r adapterid.ID)
	GetAdapterIDCounter    uint64
	GetAdapterIDPreCounter uint64
	GetAdapterIDMock       mResponseMockGetAdapterID

	GetElementIDFunc       func() (r uint32)
	GetElementIDCounter    uint64
	GetElementIDPreCounter uint64
	GetElementIDMock       mResponseMockGetElementID

	GetHandlerIDFunc       func() (r uint32)
	GetHandlerIDCounter    uint64
	GetHandlerIDPreCounter uint64
	GetHandlerIDMock       mResponseMockGetHandlerID

	GetRespPayloadFunc       func() (r interface{})
	GetRespPayloadCounter    uint64
	GetRespPayloadPreCounter uint64
	GetRespPayloadMock       mResponseMockGetRespPayload
}

//NewResponseMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/iadapter.Response
func NewResponseMock(t minimock.Tester) *ResponseMock {
	m := &ResponseMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAdapterIDMock = mResponseMockGetAdapterID{mock: m}
	m.GetElementIDMock = mResponseMockGetElementID{mock: m}
	m.GetHandlerIDMock = mResponseMockGetHandlerID{mock: m}
	m.GetRespPayloadMock = mResponseMockGetRespPayload{mock: m}

	return m
}

type mResponseMockGetAdapterID struct {
	mock              *ResponseMock
	mainExpectation   *ResponseMockGetAdapterIDExpectation
	expectationSeries []*ResponseMockGetAdapterIDExpectation
}

type ResponseMockGetAdapterIDExpectation struct {
	result *ResponseMockGetAdapterIDResult
}

type ResponseMockGetAdapterIDResult struct {
	r adapterid.ID
}

//Expect specifies that invocation of Response.GetAdapterID is expected from 1 to Infinity times
func (m *mResponseMockGetAdapterID) Expect() *mResponseMockGetAdapterID {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetAdapterIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of Response.GetAdapterID
func (m *mResponseMockGetAdapterID) Return(r adapterid.ID) *ResponseMock {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetAdapterIDExpectation{}
	}
	m.mainExpectation.result = &ResponseMockGetAdapterIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Response.GetAdapterID is expected once
func (m *mResponseMockGetAdapterID) ExpectOnce() *ResponseMockGetAdapterIDExpectation {
	m.mock.GetAdapterIDFunc = nil
	m.mainExpectation = nil

	expectation := &ResponseMockGetAdapterIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ResponseMockGetAdapterIDExpectation) Return(r adapterid.ID) {
	e.result = &ResponseMockGetAdapterIDResult{r}
}

//Set uses given function f as a mock of Response.GetAdapterID method
func (m *mResponseMockGetAdapterID) Set(f func() (r adapterid.ID)) *ResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAdapterIDFunc = f
	return m.mock
}

//GetAdapterID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.Response interface
func (m *ResponseMock) GetAdapterID() (r adapterid.ID) {
	counter := atomic.AddUint64(&m.GetAdapterIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetAdapterIDCounter, 1)

	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAdapterIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ResponseMock.GetAdapterID.")
			return
		}

		result := m.GetAdapterIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetAdapterID")
			return
		}

		r = result.r

		return
	}

	if m.GetAdapterIDMock.mainExpectation != nil {

		result := m.GetAdapterIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetAdapterID")
		}

		r = result.r

		return
	}

	if m.GetAdapterIDFunc == nil {
		m.t.Fatalf("Unexpected call to ResponseMock.GetAdapterID.")
		return
	}

	return m.GetAdapterIDFunc()
}

//GetAdapterIDMinimockCounter returns a count of ResponseMock.GetAdapterIDFunc invocations
func (m *ResponseMock) GetAdapterIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDCounter)
}

//GetAdapterIDMinimockPreCounter returns the value of ResponseMock.GetAdapterID invocations
func (m *ResponseMock) GetAdapterIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDPreCounter)
}

//GetAdapterIDFinished returns true if mock invocations count is ok
func (m *ResponseMock) GetAdapterIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetAdapterIDCounter) == uint64(len(m.GetAdapterIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetAdapterIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetAdapterIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetAdapterIDFunc != nil {
		return atomic.LoadUint64(&m.GetAdapterIDCounter) > 0
	}

	return true
}

type mResponseMockGetElementID struct {
	mock              *ResponseMock
	mainExpectation   *ResponseMockGetElementIDExpectation
	expectationSeries []*ResponseMockGetElementIDExpectation
}

type ResponseMockGetElementIDExpectation struct {
	result *ResponseMockGetElementIDResult
}

type ResponseMockGetElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of Response.GetElementID is expected from 1 to Infinity times
func (m *mResponseMockGetElementID) Expect() *mResponseMockGetElementID {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of Response.GetElementID
func (m *mResponseMockGetElementID) Return(r uint32) *ResponseMock {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetElementIDExpectation{}
	}
	m.mainExpectation.result = &ResponseMockGetElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Response.GetElementID is expected once
func (m *mResponseMockGetElementID) ExpectOnce() *ResponseMockGetElementIDExpectation {
	m.mock.GetElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &ResponseMockGetElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ResponseMockGetElementIDExpectation) Return(r uint32) {
	e.result = &ResponseMockGetElementIDResult{r}
}

//Set uses given function f as a mock of Response.GetElementID method
func (m *mResponseMockGetElementID) Set(f func() (r uint32)) *ResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetElementIDFunc = f
	return m.mock
}

//GetElementID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.Response interface
func (m *ResponseMock) GetElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetElementIDCounter, 1)

	if len(m.GetElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ResponseMock.GetElementID.")
			return
		}

		result := m.GetElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetElementIDMock.mainExpectation != nil {

		result := m.GetElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetElementID")
		}

		r = result.r

		return
	}

	if m.GetElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to ResponseMock.GetElementID.")
		return
	}

	return m.GetElementIDFunc()
}

//GetElementIDMinimockCounter returns a count of ResponseMock.GetElementIDFunc invocations
func (m *ResponseMock) GetElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDCounter)
}

//GetElementIDMinimockPreCounter returns the value of ResponseMock.GetElementID invocations
func (m *ResponseMock) GetElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDPreCounter)
}

//GetElementIDFinished returns true if mock invocations count is ok
func (m *ResponseMock) GetElementIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetElementIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetElementIDCounter) == uint64(len(m.GetElementIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetElementIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetElementIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetElementIDFunc != nil {
		return atomic.LoadUint64(&m.GetElementIDCounter) > 0
	}

	return true
}

type mResponseMockGetHandlerID struct {
	mock              *ResponseMock
	mainExpectation   *ResponseMockGetHandlerIDExpectation
	expectationSeries []*ResponseMockGetHandlerIDExpectation
}

type ResponseMockGetHandlerIDExpectation struct {
	result *ResponseMockGetHandlerIDResult
}

type ResponseMockGetHandlerIDResult struct {
	r uint32
}

//Expect specifies that invocation of Response.GetHandlerID is expected from 1 to Infinity times
func (m *mResponseMockGetHandlerID) Expect() *mResponseMockGetHandlerID {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetHandlerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of Response.GetHandlerID
func (m *mResponseMockGetHandlerID) Return(r uint32) *ResponseMock {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetHandlerIDExpectation{}
	}
	m.mainExpectation.result = &ResponseMockGetHandlerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Response.GetHandlerID is expected once
func (m *mResponseMockGetHandlerID) ExpectOnce() *ResponseMockGetHandlerIDExpectation {
	m.mock.GetHandlerIDFunc = nil
	m.mainExpectation = nil

	expectation := &ResponseMockGetHandlerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ResponseMockGetHandlerIDExpectation) Return(r uint32) {
	e.result = &ResponseMockGetHandlerIDResult{r}
}

//Set uses given function f as a mock of Response.GetHandlerID method
func (m *mResponseMockGetHandlerID) Set(f func() (r uint32)) *ResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetHandlerIDFunc = f
	return m.mock
}

//GetHandlerID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.Response interface
func (m *ResponseMock) GetHandlerID() (r uint32) {
	counter := atomic.AddUint64(&m.GetHandlerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetHandlerIDCounter, 1)

	if len(m.GetHandlerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetHandlerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ResponseMock.GetHandlerID.")
			return
		}

		result := m.GetHandlerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetHandlerID")
			return
		}

		r = result.r

		return
	}

	if m.GetHandlerIDMock.mainExpectation != nil {

		result := m.GetHandlerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetHandlerID")
		}

		r = result.r

		return
	}

	if m.GetHandlerIDFunc == nil {
		m.t.Fatalf("Unexpected call to ResponseMock.GetHandlerID.")
		return
	}

	return m.GetHandlerIDFunc()
}

//GetHandlerIDMinimockCounter returns a count of ResponseMock.GetHandlerIDFunc invocations
func (m *ResponseMock) GetHandlerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDCounter)
}

//GetHandlerIDMinimockPreCounter returns the value of ResponseMock.GetHandlerID invocations
func (m *ResponseMock) GetHandlerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDPreCounter)
}

//GetHandlerIDFinished returns true if mock invocations count is ok
func (m *ResponseMock) GetHandlerIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetHandlerIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetHandlerIDCounter) == uint64(len(m.GetHandlerIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetHandlerIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetHandlerIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetHandlerIDFunc != nil {
		return atomic.LoadUint64(&m.GetHandlerIDCounter) > 0
	}

	return true
}

type mResponseMockGetRespPayload struct {
	mock              *ResponseMock
	mainExpectation   *ResponseMockGetRespPayloadExpectation
	expectationSeries []*ResponseMockGetRespPayloadExpectation
}

type ResponseMockGetRespPayloadExpectation struct {
	result *ResponseMockGetRespPayloadResult
}

type ResponseMockGetRespPayloadResult struct {
	r interface{}
}

//Expect specifies that invocation of Response.GetRespPayload is expected from 1 to Infinity times
func (m *mResponseMockGetRespPayload) Expect() *mResponseMockGetRespPayload {
	m.mock.GetRespPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetRespPayloadExpectation{}
	}

	return m
}

//Return specifies results of invocation of Response.GetRespPayload
func (m *mResponseMockGetRespPayload) Return(r interface{}) *ResponseMock {
	m.mock.GetRespPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ResponseMockGetRespPayloadExpectation{}
	}
	m.mainExpectation.result = &ResponseMockGetRespPayloadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Response.GetRespPayload is expected once
func (m *mResponseMockGetRespPayload) ExpectOnce() *ResponseMockGetRespPayloadExpectation {
	m.mock.GetRespPayloadFunc = nil
	m.mainExpectation = nil

	expectation := &ResponseMockGetRespPayloadExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ResponseMockGetRespPayloadExpectation) Return(r interface{}) {
	e.result = &ResponseMockGetRespPayloadResult{r}
}

//Set uses given function f as a mock of Response.GetRespPayload method
func (m *mResponseMockGetRespPayload) Set(f func() (r interface{})) *ResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRespPayloadFunc = f
	return m.mock
}

//GetRespPayload implements github.com/insolar/insolar/conveyor/interfaces/iadapter.Response interface
func (m *ResponseMock) GetRespPayload() (r interface{}) {
	counter := atomic.AddUint64(&m.GetRespPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.GetRespPayloadCounter, 1)

	if len(m.GetRespPayloadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRespPayloadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ResponseMock.GetRespPayload.")
			return
		}

		result := m.GetRespPayloadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetRespPayload")
			return
		}

		r = result.r

		return
	}

	if m.GetRespPayloadMock.mainExpectation != nil {

		result := m.GetRespPayloadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ResponseMock.GetRespPayload")
		}

		r = result.r

		return
	}

	if m.GetRespPayloadFunc == nil {
		m.t.Fatalf("Unexpected call to ResponseMock.GetRespPayload.")
		return
	}

	return m.GetRespPayloadFunc()
}

//GetRespPayloadMinimockCounter returns a count of ResponseMock.GetRespPayloadFunc invocations
func (m *ResponseMock) GetRespPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRespPayloadCounter)
}

//GetRespPayloadMinimockPreCounter returns the value of ResponseMock.GetRespPayload invocations
func (m *ResponseMock) GetRespPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRespPayloadPreCounter)
}

//GetRespPayloadFinished returns true if mock invocations count is ok
func (m *ResponseMock) GetRespPayloadFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRespPayloadMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRespPayloadCounter) == uint64(len(m.GetRespPayloadMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRespPayloadMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRespPayloadCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRespPayloadFunc != nil {
		return atomic.LoadUint64(&m.GetRespPayloadCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ResponseMock) ValidateCallCounters() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetAdapterID")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetElementID")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetHandlerID")
	}

	if !m.GetRespPayloadFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetRespPayload")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ResponseMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ResponseMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ResponseMock) MinimockFinish() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetAdapterID")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetElementID")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetHandlerID")
	}

	if !m.GetRespPayloadFinished() {
		m.t.Fatal("Expected call to ResponseMock.GetRespPayload")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ResponseMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ResponseMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetAdapterIDFinished()
		ok = ok && m.GetElementIDFinished()
		ok = ok && m.GetHandlerIDFinished()
		ok = ok && m.GetRespPayloadFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAdapterIDFinished() {
				m.t.Error("Expected call to ResponseMock.GetAdapterID")
			}

			if !m.GetElementIDFinished() {
				m.t.Error("Expected call to ResponseMock.GetElementID")
			}

			if !m.GetHandlerIDFinished() {
				m.t.Error("Expected call to ResponseMock.GetHandlerID")
			}

			if !m.GetRespPayloadFinished() {
				m.t.Error("Expected call to ResponseMock.GetRespPayload")
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
func (m *ResponseMock) AllMocksCalled() bool {

	if !m.GetAdapterIDFinished() {
		return false
	}

	if !m.GetElementIDFinished() {
		return false
	}

	if !m.GetHandlerIDFinished() {
		return false
	}

	if !m.GetRespPayloadFinished() {
		return false
	}

	return true
}

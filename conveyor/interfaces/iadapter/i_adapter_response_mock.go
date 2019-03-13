package iadapter

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IAdapterResponse" can be found in github.com/insolar/insolar/conveyor/interfaces/iadapter
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//IAdapterResponseMock implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterResponse
type IAdapterResponseMock struct {
	t minimock.Tester

	GetAdapterIDFunc       func() (r uint32)
	GetAdapterIDCounter    uint64
	GetAdapterIDPreCounter uint64
	GetAdapterIDMock       mIAdapterResponseMockGetAdapterID

	GetElementIDFunc       func() (r uint32)
	GetElementIDCounter    uint64
	GetElementIDPreCounter uint64
	GetElementIDMock       mIAdapterResponseMockGetElementID

	GetHandlerIDFunc       func() (r uint32)
	GetHandlerIDCounter    uint64
	GetHandlerIDPreCounter uint64
	GetHandlerIDMock       mIAdapterResponseMockGetHandlerID

	GetRespPayloadFunc       func() (r interface{})
	GetRespPayloadCounter    uint64
	GetRespPayloadPreCounter uint64
	GetRespPayloadMock       mIAdapterResponseMockGetRespPayload
}

//NewIAdapterResponseMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterResponse
func NewIAdapterResponseMock(t minimock.Tester) *IAdapterResponseMock {
	m := &IAdapterResponseMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAdapterIDMock = mIAdapterResponseMockGetAdapterID{mock: m}
	m.GetElementIDMock = mIAdapterResponseMockGetElementID{mock: m}
	m.GetHandlerIDMock = mIAdapterResponseMockGetHandlerID{mock: m}
	m.GetRespPayloadMock = mIAdapterResponseMockGetRespPayload{mock: m}

	return m
}

type mIAdapterResponseMockGetAdapterID struct {
	mock              *IAdapterResponseMock
	mainExpectation   *IAdapterResponseMockGetAdapterIDExpectation
	expectationSeries []*IAdapterResponseMockGetAdapterIDExpectation
}

type IAdapterResponseMockGetAdapterIDExpectation struct {
	result *IAdapterResponseMockGetAdapterIDResult
}

type IAdapterResponseMockGetAdapterIDResult struct {
	r uint32
}

//Expect specifies that invocation of IAdapterResponse.GetAdapterID is expected from 1 to Infinity times
func (m *mIAdapterResponseMockGetAdapterID) Expect() *mIAdapterResponseMockGetAdapterID {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetAdapterIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterResponse.GetAdapterID
func (m *mIAdapterResponseMockGetAdapterID) Return(r uint32) *IAdapterResponseMock {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetAdapterIDExpectation{}
	}
	m.mainExpectation.result = &IAdapterResponseMockGetAdapterIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterResponse.GetAdapterID is expected once
func (m *mIAdapterResponseMockGetAdapterID) ExpectOnce() *IAdapterResponseMockGetAdapterIDExpectation {
	m.mock.GetAdapterIDFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterResponseMockGetAdapterIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterResponseMockGetAdapterIDExpectation) Return(r uint32) {
	e.result = &IAdapterResponseMockGetAdapterIDResult{r}
}

//Set uses given function f as a mock of IAdapterResponse.GetAdapterID method
func (m *mIAdapterResponseMockGetAdapterID) Set(f func() (r uint32)) *IAdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAdapterIDFunc = f
	return m.mock
}

//GetAdapterID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterResponse interface
func (m *IAdapterResponseMock) GetAdapterID() (r uint32) {
	counter := atomic.AddUint64(&m.GetAdapterIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetAdapterIDCounter, 1)

	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAdapterIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetAdapterID.")
			return
		}

		result := m.GetAdapterIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetAdapterID")
			return
		}

		r = result.r

		return
	}

	if m.GetAdapterIDMock.mainExpectation != nil {

		result := m.GetAdapterIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetAdapterID")
		}

		r = result.r

		return
	}

	if m.GetAdapterIDFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetAdapterID.")
		return
	}

	return m.GetAdapterIDFunc()
}

//GetAdapterIDMinimockCounter returns a count of IAdapterResponseMock.GetAdapterIDFunc invocations
func (m *IAdapterResponseMock) GetAdapterIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDCounter)
}

//GetAdapterIDMinimockPreCounter returns the value of IAdapterResponseMock.GetAdapterID invocations
func (m *IAdapterResponseMock) GetAdapterIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDPreCounter)
}

//GetAdapterIDFinished returns true if mock invocations count is ok
func (m *IAdapterResponseMock) GetAdapterIDFinished() bool {
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

type mIAdapterResponseMockGetElementID struct {
	mock              *IAdapterResponseMock
	mainExpectation   *IAdapterResponseMockGetElementIDExpectation
	expectationSeries []*IAdapterResponseMockGetElementIDExpectation
}

type IAdapterResponseMockGetElementIDExpectation struct {
	result *IAdapterResponseMockGetElementIDResult
}

type IAdapterResponseMockGetElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of IAdapterResponse.GetElementID is expected from 1 to Infinity times
func (m *mIAdapterResponseMockGetElementID) Expect() *mIAdapterResponseMockGetElementID {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterResponse.GetElementID
func (m *mIAdapterResponseMockGetElementID) Return(r uint32) *IAdapterResponseMock {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetElementIDExpectation{}
	}
	m.mainExpectation.result = &IAdapterResponseMockGetElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterResponse.GetElementID is expected once
func (m *mIAdapterResponseMockGetElementID) ExpectOnce() *IAdapterResponseMockGetElementIDExpectation {
	m.mock.GetElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterResponseMockGetElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterResponseMockGetElementIDExpectation) Return(r uint32) {
	e.result = &IAdapterResponseMockGetElementIDResult{r}
}

//Set uses given function f as a mock of IAdapterResponse.GetElementID method
func (m *mIAdapterResponseMockGetElementID) Set(f func() (r uint32)) *IAdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetElementIDFunc = f
	return m.mock
}

//GetElementID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterResponse interface
func (m *IAdapterResponseMock) GetElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetElementIDCounter, 1)

	if len(m.GetElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetElementID.")
			return
		}

		result := m.GetElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetElementIDMock.mainExpectation != nil {

		result := m.GetElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetElementID")
		}

		r = result.r

		return
	}

	if m.GetElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetElementID.")
		return
	}

	return m.GetElementIDFunc()
}

//GetElementIDMinimockCounter returns a count of IAdapterResponseMock.GetElementIDFunc invocations
func (m *IAdapterResponseMock) GetElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDCounter)
}

//GetElementIDMinimockPreCounter returns the value of IAdapterResponseMock.GetElementID invocations
func (m *IAdapterResponseMock) GetElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDPreCounter)
}

//GetElementIDFinished returns true if mock invocations count is ok
func (m *IAdapterResponseMock) GetElementIDFinished() bool {
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

type mIAdapterResponseMockGetHandlerID struct {
	mock              *IAdapterResponseMock
	mainExpectation   *IAdapterResponseMockGetHandlerIDExpectation
	expectationSeries []*IAdapterResponseMockGetHandlerIDExpectation
}

type IAdapterResponseMockGetHandlerIDExpectation struct {
	result *IAdapterResponseMockGetHandlerIDResult
}

type IAdapterResponseMockGetHandlerIDResult struct {
	r uint32
}

//Expect specifies that invocation of IAdapterResponse.GetHandlerID is expected from 1 to Infinity times
func (m *mIAdapterResponseMockGetHandlerID) Expect() *mIAdapterResponseMockGetHandlerID {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetHandlerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterResponse.GetHandlerID
func (m *mIAdapterResponseMockGetHandlerID) Return(r uint32) *IAdapterResponseMock {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetHandlerIDExpectation{}
	}
	m.mainExpectation.result = &IAdapterResponseMockGetHandlerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterResponse.GetHandlerID is expected once
func (m *mIAdapterResponseMockGetHandlerID) ExpectOnce() *IAdapterResponseMockGetHandlerIDExpectation {
	m.mock.GetHandlerIDFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterResponseMockGetHandlerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterResponseMockGetHandlerIDExpectation) Return(r uint32) {
	e.result = &IAdapterResponseMockGetHandlerIDResult{r}
}

//Set uses given function f as a mock of IAdapterResponse.GetHandlerID method
func (m *mIAdapterResponseMockGetHandlerID) Set(f func() (r uint32)) *IAdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetHandlerIDFunc = f
	return m.mock
}

//GetHandlerID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterResponse interface
func (m *IAdapterResponseMock) GetHandlerID() (r uint32) {
	counter := atomic.AddUint64(&m.GetHandlerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetHandlerIDCounter, 1)

	if len(m.GetHandlerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetHandlerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetHandlerID.")
			return
		}

		result := m.GetHandlerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetHandlerID")
			return
		}

		r = result.r

		return
	}

	if m.GetHandlerIDMock.mainExpectation != nil {

		result := m.GetHandlerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetHandlerID")
		}

		r = result.r

		return
	}

	if m.GetHandlerIDFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetHandlerID.")
		return
	}

	return m.GetHandlerIDFunc()
}

//GetHandlerIDMinimockCounter returns a count of IAdapterResponseMock.GetHandlerIDFunc invocations
func (m *IAdapterResponseMock) GetHandlerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDCounter)
}

//GetHandlerIDMinimockPreCounter returns the value of IAdapterResponseMock.GetHandlerID invocations
func (m *IAdapterResponseMock) GetHandlerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDPreCounter)
}

//GetHandlerIDFinished returns true if mock invocations count is ok
func (m *IAdapterResponseMock) GetHandlerIDFinished() bool {
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

type mIAdapterResponseMockGetRespPayload struct {
	mock              *IAdapterResponseMock
	mainExpectation   *IAdapterResponseMockGetRespPayloadExpectation
	expectationSeries []*IAdapterResponseMockGetRespPayloadExpectation
}

type IAdapterResponseMockGetRespPayloadExpectation struct {
	result *IAdapterResponseMockGetRespPayloadResult
}

type IAdapterResponseMockGetRespPayloadResult struct {
	r interface{}
}

//Expect specifies that invocation of IAdapterResponse.GetRespPayload is expected from 1 to Infinity times
func (m *mIAdapterResponseMockGetRespPayload) Expect() *mIAdapterResponseMockGetRespPayload {
	m.mock.GetRespPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetRespPayloadExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterResponse.GetRespPayload
func (m *mIAdapterResponseMockGetRespPayload) Return(r interface{}) *IAdapterResponseMock {
	m.mock.GetRespPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterResponseMockGetRespPayloadExpectation{}
	}
	m.mainExpectation.result = &IAdapterResponseMockGetRespPayloadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterResponse.GetRespPayload is expected once
func (m *mIAdapterResponseMockGetRespPayload) ExpectOnce() *IAdapterResponseMockGetRespPayloadExpectation {
	m.mock.GetRespPayloadFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterResponseMockGetRespPayloadExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterResponseMockGetRespPayloadExpectation) Return(r interface{}) {
	e.result = &IAdapterResponseMockGetRespPayloadResult{r}
}

//Set uses given function f as a mock of IAdapterResponse.GetRespPayload method
func (m *mIAdapterResponseMockGetRespPayload) Set(f func() (r interface{})) *IAdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRespPayloadFunc = f
	return m.mock
}

//GetRespPayload implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterResponse interface
func (m *IAdapterResponseMock) GetRespPayload() (r interface{}) {
	counter := atomic.AddUint64(&m.GetRespPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.GetRespPayloadCounter, 1)

	if len(m.GetRespPayloadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRespPayloadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetRespPayload.")
			return
		}

		result := m.GetRespPayloadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetRespPayload")
			return
		}

		r = result.r

		return
	}

	if m.GetRespPayloadMock.mainExpectation != nil {

		result := m.GetRespPayloadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterResponseMock.GetRespPayload")
		}

		r = result.r

		return
	}

	if m.GetRespPayloadFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterResponseMock.GetRespPayload.")
		return
	}

	return m.GetRespPayloadFunc()
}

//GetRespPayloadMinimockCounter returns a count of IAdapterResponseMock.GetRespPayloadFunc invocations
func (m *IAdapterResponseMock) GetRespPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRespPayloadCounter)
}

//GetRespPayloadMinimockPreCounter returns the value of IAdapterResponseMock.GetRespPayload invocations
func (m *IAdapterResponseMock) GetRespPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRespPayloadPreCounter)
}

//GetRespPayloadFinished returns true if mock invocations count is ok
func (m *IAdapterResponseMock) GetRespPayloadFinished() bool {
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
func (m *IAdapterResponseMock) ValidateCallCounters() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetAdapterID")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetElementID")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetHandlerID")
	}

	if !m.GetRespPayloadFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetRespPayload")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IAdapterResponseMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IAdapterResponseMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IAdapterResponseMock) MinimockFinish() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetAdapterID")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetElementID")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetHandlerID")
	}

	if !m.GetRespPayloadFinished() {
		m.t.Fatal("Expected call to IAdapterResponseMock.GetRespPayload")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IAdapterResponseMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IAdapterResponseMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IAdapterResponseMock.GetAdapterID")
			}

			if !m.GetElementIDFinished() {
				m.t.Error("Expected call to IAdapterResponseMock.GetElementID")
			}

			if !m.GetHandlerIDFinished() {
				m.t.Error("Expected call to IAdapterResponseMock.GetHandlerID")
			}

			if !m.GetRespPayloadFinished() {
				m.t.Error("Expected call to IAdapterResponseMock.GetRespPayload")
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
func (m *IAdapterResponseMock) AllMocksCalled() bool {

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

package iadapter

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IAdapterNestedEvent" can be found in github.com/insolar/insolar/conveyor/interfaces/iadapter
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//IAdapterNestedEventMock implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterNestedEvent
type IAdapterNestedEventMock struct {
	t minimock.Tester

	GetAdapterIDFunc       func() (r uint32)
	GetAdapterIDCounter    uint64
	GetAdapterIDPreCounter uint64
	GetAdapterIDMock       mIAdapterNestedEventMockGetAdapterID

	GetEventPayloadFunc       func() (r interface{})
	GetEventPayloadCounter    uint64
	GetEventPayloadPreCounter uint64
	GetEventPayloadMock       mIAdapterNestedEventMockGetEventPayload

	GetHandlerIDFunc       func() (r uint32)
	GetHandlerIDCounter    uint64
	GetHandlerIDPreCounter uint64
	GetHandlerIDMock       mIAdapterNestedEventMockGetHandlerID

	GetParentElementIDFunc       func() (r uint32)
	GetParentElementIDCounter    uint64
	GetParentElementIDPreCounter uint64
	GetParentElementIDMock       mIAdapterNestedEventMockGetParentElementID
}

//NewIAdapterNestedEventMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterNestedEvent
func NewIAdapterNestedEventMock(t minimock.Tester) *IAdapterNestedEventMock {
	m := &IAdapterNestedEventMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAdapterIDMock = mIAdapterNestedEventMockGetAdapterID{mock: m}
	m.GetEventPayloadMock = mIAdapterNestedEventMockGetEventPayload{mock: m}
	m.GetHandlerIDMock = mIAdapterNestedEventMockGetHandlerID{mock: m}
	m.GetParentElementIDMock = mIAdapterNestedEventMockGetParentElementID{mock: m}

	return m
}

type mIAdapterNestedEventMockGetAdapterID struct {
	mock              *IAdapterNestedEventMock
	mainExpectation   *IAdapterNestedEventMockGetAdapterIDExpectation
	expectationSeries []*IAdapterNestedEventMockGetAdapterIDExpectation
}

type IAdapterNestedEventMockGetAdapterIDExpectation struct {
	result *IAdapterNestedEventMockGetAdapterIDResult
}

type IAdapterNestedEventMockGetAdapterIDResult struct {
	r uint32
}

//Expect specifies that invocation of IAdapterNestedEvent.GetAdapterID is expected from 1 to Infinity times
func (m *mIAdapterNestedEventMockGetAdapterID) Expect() *mIAdapterNestedEventMockGetAdapterID {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetAdapterIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterNestedEvent.GetAdapterID
func (m *mIAdapterNestedEventMockGetAdapterID) Return(r uint32) *IAdapterNestedEventMock {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetAdapterIDExpectation{}
	}
	m.mainExpectation.result = &IAdapterNestedEventMockGetAdapterIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterNestedEvent.GetAdapterID is expected once
func (m *mIAdapterNestedEventMockGetAdapterID) ExpectOnce() *IAdapterNestedEventMockGetAdapterIDExpectation {
	m.mock.GetAdapterIDFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterNestedEventMockGetAdapterIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterNestedEventMockGetAdapterIDExpectation) Return(r uint32) {
	e.result = &IAdapterNestedEventMockGetAdapterIDResult{r}
}

//Set uses given function f as a mock of IAdapterNestedEvent.GetAdapterID method
func (m *mIAdapterNestedEventMockGetAdapterID) Set(f func() (r uint32)) *IAdapterNestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAdapterIDFunc = f
	return m.mock
}

//GetAdapterID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterNestedEvent interface
func (m *IAdapterNestedEventMock) GetAdapterID() (r uint32) {
	counter := atomic.AddUint64(&m.GetAdapterIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetAdapterIDCounter, 1)

	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAdapterIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetAdapterID.")
			return
		}

		result := m.GetAdapterIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetAdapterID")
			return
		}

		r = result.r

		return
	}

	if m.GetAdapterIDMock.mainExpectation != nil {

		result := m.GetAdapterIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetAdapterID")
		}

		r = result.r

		return
	}

	if m.GetAdapterIDFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetAdapterID.")
		return
	}

	return m.GetAdapterIDFunc()
}

//GetAdapterIDMinimockCounter returns a count of IAdapterNestedEventMock.GetAdapterIDFunc invocations
func (m *IAdapterNestedEventMock) GetAdapterIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDCounter)
}

//GetAdapterIDMinimockPreCounter returns the value of IAdapterNestedEventMock.GetAdapterID invocations
func (m *IAdapterNestedEventMock) GetAdapterIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDPreCounter)
}

//GetAdapterIDFinished returns true if mock invocations count is ok
func (m *IAdapterNestedEventMock) GetAdapterIDFinished() bool {
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

type mIAdapterNestedEventMockGetEventPayload struct {
	mock              *IAdapterNestedEventMock
	mainExpectation   *IAdapterNestedEventMockGetEventPayloadExpectation
	expectationSeries []*IAdapterNestedEventMockGetEventPayloadExpectation
}

type IAdapterNestedEventMockGetEventPayloadExpectation struct {
	result *IAdapterNestedEventMockGetEventPayloadResult
}

type IAdapterNestedEventMockGetEventPayloadResult struct {
	r interface{}
}

//Expect specifies that invocation of IAdapterNestedEvent.GetEventPayload is expected from 1 to Infinity times
func (m *mIAdapterNestedEventMockGetEventPayload) Expect() *mIAdapterNestedEventMockGetEventPayload {
	m.mock.GetEventPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetEventPayloadExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterNestedEvent.GetEventPayload
func (m *mIAdapterNestedEventMockGetEventPayload) Return(r interface{}) *IAdapterNestedEventMock {
	m.mock.GetEventPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetEventPayloadExpectation{}
	}
	m.mainExpectation.result = &IAdapterNestedEventMockGetEventPayloadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterNestedEvent.GetEventPayload is expected once
func (m *mIAdapterNestedEventMockGetEventPayload) ExpectOnce() *IAdapterNestedEventMockGetEventPayloadExpectation {
	m.mock.GetEventPayloadFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterNestedEventMockGetEventPayloadExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterNestedEventMockGetEventPayloadExpectation) Return(r interface{}) {
	e.result = &IAdapterNestedEventMockGetEventPayloadResult{r}
}

//Set uses given function f as a mock of IAdapterNestedEvent.GetEventPayload method
func (m *mIAdapterNestedEventMockGetEventPayload) Set(f func() (r interface{})) *IAdapterNestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetEventPayloadFunc = f
	return m.mock
}

//GetEventPayload implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterNestedEvent interface
func (m *IAdapterNestedEventMock) GetEventPayload() (r interface{}) {
	counter := atomic.AddUint64(&m.GetEventPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.GetEventPayloadCounter, 1)

	if len(m.GetEventPayloadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetEventPayloadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetEventPayload.")
			return
		}

		result := m.GetEventPayloadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetEventPayload")
			return
		}

		r = result.r

		return
	}

	if m.GetEventPayloadMock.mainExpectation != nil {

		result := m.GetEventPayloadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetEventPayload")
		}

		r = result.r

		return
	}

	if m.GetEventPayloadFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetEventPayload.")
		return
	}

	return m.GetEventPayloadFunc()
}

//GetEventPayloadMinimockCounter returns a count of IAdapterNestedEventMock.GetEventPayloadFunc invocations
func (m *IAdapterNestedEventMock) GetEventPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetEventPayloadCounter)
}

//GetEventPayloadMinimockPreCounter returns the value of IAdapterNestedEventMock.GetEventPayload invocations
func (m *IAdapterNestedEventMock) GetEventPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetEventPayloadPreCounter)
}

//GetEventPayloadFinished returns true if mock invocations count is ok
func (m *IAdapterNestedEventMock) GetEventPayloadFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetEventPayloadMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetEventPayloadCounter) == uint64(len(m.GetEventPayloadMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetEventPayloadMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetEventPayloadCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetEventPayloadFunc != nil {
		return atomic.LoadUint64(&m.GetEventPayloadCounter) > 0
	}

	return true
}

type mIAdapterNestedEventMockGetHandlerID struct {
	mock              *IAdapterNestedEventMock
	mainExpectation   *IAdapterNestedEventMockGetHandlerIDExpectation
	expectationSeries []*IAdapterNestedEventMockGetHandlerIDExpectation
}

type IAdapterNestedEventMockGetHandlerIDExpectation struct {
	result *IAdapterNestedEventMockGetHandlerIDResult
}

type IAdapterNestedEventMockGetHandlerIDResult struct {
	r uint32
}

//Expect specifies that invocation of IAdapterNestedEvent.GetHandlerID is expected from 1 to Infinity times
func (m *mIAdapterNestedEventMockGetHandlerID) Expect() *mIAdapterNestedEventMockGetHandlerID {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetHandlerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterNestedEvent.GetHandlerID
func (m *mIAdapterNestedEventMockGetHandlerID) Return(r uint32) *IAdapterNestedEventMock {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetHandlerIDExpectation{}
	}
	m.mainExpectation.result = &IAdapterNestedEventMockGetHandlerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterNestedEvent.GetHandlerID is expected once
func (m *mIAdapterNestedEventMockGetHandlerID) ExpectOnce() *IAdapterNestedEventMockGetHandlerIDExpectation {
	m.mock.GetHandlerIDFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterNestedEventMockGetHandlerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterNestedEventMockGetHandlerIDExpectation) Return(r uint32) {
	e.result = &IAdapterNestedEventMockGetHandlerIDResult{r}
}

//Set uses given function f as a mock of IAdapterNestedEvent.GetHandlerID method
func (m *mIAdapterNestedEventMockGetHandlerID) Set(f func() (r uint32)) *IAdapterNestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetHandlerIDFunc = f
	return m.mock
}

//GetHandlerID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterNestedEvent interface
func (m *IAdapterNestedEventMock) GetHandlerID() (r uint32) {
	counter := atomic.AddUint64(&m.GetHandlerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetHandlerIDCounter, 1)

	if len(m.GetHandlerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetHandlerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetHandlerID.")
			return
		}

		result := m.GetHandlerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetHandlerID")
			return
		}

		r = result.r

		return
	}

	if m.GetHandlerIDMock.mainExpectation != nil {

		result := m.GetHandlerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetHandlerID")
		}

		r = result.r

		return
	}

	if m.GetHandlerIDFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetHandlerID.")
		return
	}

	return m.GetHandlerIDFunc()
}

//GetHandlerIDMinimockCounter returns a count of IAdapterNestedEventMock.GetHandlerIDFunc invocations
func (m *IAdapterNestedEventMock) GetHandlerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDCounter)
}

//GetHandlerIDMinimockPreCounter returns the value of IAdapterNestedEventMock.GetHandlerID invocations
func (m *IAdapterNestedEventMock) GetHandlerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDPreCounter)
}

//GetHandlerIDFinished returns true if mock invocations count is ok
func (m *IAdapterNestedEventMock) GetHandlerIDFinished() bool {
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

type mIAdapterNestedEventMockGetParentElementID struct {
	mock              *IAdapterNestedEventMock
	mainExpectation   *IAdapterNestedEventMockGetParentElementIDExpectation
	expectationSeries []*IAdapterNestedEventMockGetParentElementIDExpectation
}

type IAdapterNestedEventMockGetParentElementIDExpectation struct {
	result *IAdapterNestedEventMockGetParentElementIDResult
}

type IAdapterNestedEventMockGetParentElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of IAdapterNestedEvent.GetParentElementID is expected from 1 to Infinity times
func (m *mIAdapterNestedEventMockGetParentElementID) Expect() *mIAdapterNestedEventMockGetParentElementID {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetParentElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of IAdapterNestedEvent.GetParentElementID
func (m *mIAdapterNestedEventMockGetParentElementID) Return(r uint32) *IAdapterNestedEventMock {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IAdapterNestedEventMockGetParentElementIDExpectation{}
	}
	m.mainExpectation.result = &IAdapterNestedEventMockGetParentElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IAdapterNestedEvent.GetParentElementID is expected once
func (m *mIAdapterNestedEventMockGetParentElementID) ExpectOnce() *IAdapterNestedEventMockGetParentElementIDExpectation {
	m.mock.GetParentElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &IAdapterNestedEventMockGetParentElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IAdapterNestedEventMockGetParentElementIDExpectation) Return(r uint32) {
	e.result = &IAdapterNestedEventMockGetParentElementIDResult{r}
}

//Set uses given function f as a mock of IAdapterNestedEvent.GetParentElementID method
func (m *mIAdapterNestedEventMockGetParentElementID) Set(f func() (r uint32)) *IAdapterNestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetParentElementIDFunc = f
	return m.mock
}

//GetParentElementID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.IAdapterNestedEvent interface
func (m *IAdapterNestedEventMock) GetParentElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetParentElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetParentElementIDCounter, 1)

	if len(m.GetParentElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetParentElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetParentElementID.")
			return
		}

		result := m.GetParentElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetParentElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetParentElementIDMock.mainExpectation != nil {

		result := m.GetParentElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IAdapterNestedEventMock.GetParentElementID")
		}

		r = result.r

		return
	}

	if m.GetParentElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to IAdapterNestedEventMock.GetParentElementID.")
		return
	}

	return m.GetParentElementIDFunc()
}

//GetParentElementIDMinimockCounter returns a count of IAdapterNestedEventMock.GetParentElementIDFunc invocations
func (m *IAdapterNestedEventMock) GetParentElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDCounter)
}

//GetParentElementIDMinimockPreCounter returns the value of IAdapterNestedEventMock.GetParentElementID invocations
func (m *IAdapterNestedEventMock) GetParentElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDPreCounter)
}

//GetParentElementIDFinished returns true if mock invocations count is ok
func (m *IAdapterNestedEventMock) GetParentElementIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetParentElementIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetParentElementIDCounter) == uint64(len(m.GetParentElementIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetParentElementIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetParentElementIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetParentElementIDFunc != nil {
		return atomic.LoadUint64(&m.GetParentElementIDCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IAdapterNestedEventMock) ValidateCallCounters() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetAdapterID")
	}

	if !m.GetEventPayloadFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetEventPayload")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetHandlerID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetParentElementID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IAdapterNestedEventMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IAdapterNestedEventMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IAdapterNestedEventMock) MinimockFinish() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetAdapterID")
	}

	if !m.GetEventPayloadFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetEventPayload")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetHandlerID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to IAdapterNestedEventMock.GetParentElementID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IAdapterNestedEventMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IAdapterNestedEventMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetAdapterIDFinished()
		ok = ok && m.GetEventPayloadFinished()
		ok = ok && m.GetHandlerIDFinished()
		ok = ok && m.GetParentElementIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAdapterIDFinished() {
				m.t.Error("Expected call to IAdapterNestedEventMock.GetAdapterID")
			}

			if !m.GetEventPayloadFinished() {
				m.t.Error("Expected call to IAdapterNestedEventMock.GetEventPayload")
			}

			if !m.GetHandlerIDFinished() {
				m.t.Error("Expected call to IAdapterNestedEventMock.GetHandlerID")
			}

			if !m.GetParentElementIDFinished() {
				m.t.Error("Expected call to IAdapterNestedEventMock.GetParentElementID")
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
func (m *IAdapterNestedEventMock) AllMocksCalled() bool {

	if !m.GetAdapterIDFinished() {
		return false
	}

	if !m.GetEventPayloadFinished() {
		return false
	}

	if !m.GetHandlerIDFinished() {
		return false
	}

	if !m.GetParentElementIDFinished() {
		return false
	}

	return true
}

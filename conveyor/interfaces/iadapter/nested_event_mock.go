package iadapter

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NestedEvent" can be found in github.com/insolar/insolar/conveyor/interfaces/iadapter
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//NestedEventMock implements github.com/insolar/insolar/conveyor/interfaces/iadapter.NestedEvent
type NestedEventMock struct {
	t minimock.Tester

	GetAdapterIDFunc       func() (r uint32)
	GetAdapterIDCounter    uint64
	GetAdapterIDPreCounter uint64
	GetAdapterIDMock       mNestedEventMockGetAdapterID

	GetEventPayloadFunc       func() (r interface{})
	GetEventPayloadCounter    uint64
	GetEventPayloadPreCounter uint64
	GetEventPayloadMock       mNestedEventMockGetEventPayload

	GetHandlerIDFunc       func() (r uint32)
	GetHandlerIDCounter    uint64
	GetHandlerIDPreCounter uint64
	GetHandlerIDMock       mNestedEventMockGetHandlerID

	GetParentElementIDFunc       func() (r uint32)
	GetParentElementIDCounter    uint64
	GetParentElementIDPreCounter uint64
	GetParentElementIDMock       mNestedEventMockGetParentElementID
}

//NewNestedEventMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/iadapter.NestedEvent
func NewNestedEventMock(t minimock.Tester) *NestedEventMock {
	m := &NestedEventMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAdapterIDMock = mNestedEventMockGetAdapterID{mock: m}
	m.GetEventPayloadMock = mNestedEventMockGetEventPayload{mock: m}
	m.GetHandlerIDMock = mNestedEventMockGetHandlerID{mock: m}
	m.GetParentElementIDMock = mNestedEventMockGetParentElementID{mock: m}

	return m
}

type mNestedEventMockGetAdapterID struct {
	mock              *NestedEventMock
	mainExpectation   *NestedEventMockGetAdapterIDExpectation
	expectationSeries []*NestedEventMockGetAdapterIDExpectation
}

type NestedEventMockGetAdapterIDExpectation struct {
	result *NestedEventMockGetAdapterIDResult
}

type NestedEventMockGetAdapterIDResult struct {
	r uint32
}

//Expect specifies that invocation of NestedEvent.GetAdapterID is expected from 1 to Infinity times
func (m *mNestedEventMockGetAdapterID) Expect() *mNestedEventMockGetAdapterID {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetAdapterIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NestedEvent.GetAdapterID
func (m *mNestedEventMockGetAdapterID) Return(r uint32) *NestedEventMock {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetAdapterIDExpectation{}
	}
	m.mainExpectation.result = &NestedEventMockGetAdapterIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NestedEvent.GetAdapterID is expected once
func (m *mNestedEventMockGetAdapterID) ExpectOnce() *NestedEventMockGetAdapterIDExpectation {
	m.mock.GetAdapterIDFunc = nil
	m.mainExpectation = nil

	expectation := &NestedEventMockGetAdapterIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NestedEventMockGetAdapterIDExpectation) Return(r uint32) {
	e.result = &NestedEventMockGetAdapterIDResult{r}
}

//Set uses given function f as a mock of NestedEvent.GetAdapterID method
func (m *mNestedEventMockGetAdapterID) Set(f func() (r uint32)) *NestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAdapterIDFunc = f
	return m.mock
}

//GetAdapterID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.NestedEvent interface
func (m *NestedEventMock) GetAdapterID() (r uint32) {
	counter := atomic.AddUint64(&m.GetAdapterIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetAdapterIDCounter, 1)

	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAdapterIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NestedEventMock.GetAdapterID.")
			return
		}

		result := m.GetAdapterIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetAdapterID")
			return
		}

		r = result.r

		return
	}

	if m.GetAdapterIDMock.mainExpectation != nil {

		result := m.GetAdapterIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetAdapterID")
		}

		r = result.r

		return
	}

	if m.GetAdapterIDFunc == nil {
		m.t.Fatalf("Unexpected call to NestedEventMock.GetAdapterID.")
		return
	}

	return m.GetAdapterIDFunc()
}

//GetAdapterIDMinimockCounter returns a count of NestedEventMock.GetAdapterIDFunc invocations
func (m *NestedEventMock) GetAdapterIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDCounter)
}

//GetAdapterIDMinimockPreCounter returns the value of NestedEventMock.GetAdapterID invocations
func (m *NestedEventMock) GetAdapterIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDPreCounter)
}

//GetAdapterIDFinished returns true if mock invocations count is ok
func (m *NestedEventMock) GetAdapterIDFinished() bool {
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

type mNestedEventMockGetEventPayload struct {
	mock              *NestedEventMock
	mainExpectation   *NestedEventMockGetEventPayloadExpectation
	expectationSeries []*NestedEventMockGetEventPayloadExpectation
}

type NestedEventMockGetEventPayloadExpectation struct {
	result *NestedEventMockGetEventPayloadResult
}

type NestedEventMockGetEventPayloadResult struct {
	r interface{}
}

//Expect specifies that invocation of NestedEvent.GetEventPayload is expected from 1 to Infinity times
func (m *mNestedEventMockGetEventPayload) Expect() *mNestedEventMockGetEventPayload {
	m.mock.GetEventPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetEventPayloadExpectation{}
	}

	return m
}

//Return specifies results of invocation of NestedEvent.GetEventPayload
func (m *mNestedEventMockGetEventPayload) Return(r interface{}) *NestedEventMock {
	m.mock.GetEventPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetEventPayloadExpectation{}
	}
	m.mainExpectation.result = &NestedEventMockGetEventPayloadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NestedEvent.GetEventPayload is expected once
func (m *mNestedEventMockGetEventPayload) ExpectOnce() *NestedEventMockGetEventPayloadExpectation {
	m.mock.GetEventPayloadFunc = nil
	m.mainExpectation = nil

	expectation := &NestedEventMockGetEventPayloadExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NestedEventMockGetEventPayloadExpectation) Return(r interface{}) {
	e.result = &NestedEventMockGetEventPayloadResult{r}
}

//Set uses given function f as a mock of NestedEvent.GetEventPayload method
func (m *mNestedEventMockGetEventPayload) Set(f func() (r interface{})) *NestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetEventPayloadFunc = f
	return m.mock
}

//GetEventPayload implements github.com/insolar/insolar/conveyor/interfaces/iadapter.NestedEvent interface
func (m *NestedEventMock) GetEventPayload() (r interface{}) {
	counter := atomic.AddUint64(&m.GetEventPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.GetEventPayloadCounter, 1)

	if len(m.GetEventPayloadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetEventPayloadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NestedEventMock.GetEventPayload.")
			return
		}

		result := m.GetEventPayloadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetEventPayload")
			return
		}

		r = result.r

		return
	}

	if m.GetEventPayloadMock.mainExpectation != nil {

		result := m.GetEventPayloadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetEventPayload")
		}

		r = result.r

		return
	}

	if m.GetEventPayloadFunc == nil {
		m.t.Fatalf("Unexpected call to NestedEventMock.GetEventPayload.")
		return
	}

	return m.GetEventPayloadFunc()
}

//GetEventPayloadMinimockCounter returns a count of NestedEventMock.GetEventPayloadFunc invocations
func (m *NestedEventMock) GetEventPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetEventPayloadCounter)
}

//GetEventPayloadMinimockPreCounter returns the value of NestedEventMock.GetEventPayload invocations
func (m *NestedEventMock) GetEventPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetEventPayloadPreCounter)
}

//GetEventPayloadFinished returns true if mock invocations count is ok
func (m *NestedEventMock) GetEventPayloadFinished() bool {
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

type mNestedEventMockGetHandlerID struct {
	mock              *NestedEventMock
	mainExpectation   *NestedEventMockGetHandlerIDExpectation
	expectationSeries []*NestedEventMockGetHandlerIDExpectation
}

type NestedEventMockGetHandlerIDExpectation struct {
	result *NestedEventMockGetHandlerIDResult
}

type NestedEventMockGetHandlerIDResult struct {
	r uint32
}

//Expect specifies that invocation of NestedEvent.GetHandlerID is expected from 1 to Infinity times
func (m *mNestedEventMockGetHandlerID) Expect() *mNestedEventMockGetHandlerID {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetHandlerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NestedEvent.GetHandlerID
func (m *mNestedEventMockGetHandlerID) Return(r uint32) *NestedEventMock {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetHandlerIDExpectation{}
	}
	m.mainExpectation.result = &NestedEventMockGetHandlerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NestedEvent.GetHandlerID is expected once
func (m *mNestedEventMockGetHandlerID) ExpectOnce() *NestedEventMockGetHandlerIDExpectation {
	m.mock.GetHandlerIDFunc = nil
	m.mainExpectation = nil

	expectation := &NestedEventMockGetHandlerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NestedEventMockGetHandlerIDExpectation) Return(r uint32) {
	e.result = &NestedEventMockGetHandlerIDResult{r}
}

//Set uses given function f as a mock of NestedEvent.GetHandlerID method
func (m *mNestedEventMockGetHandlerID) Set(f func() (r uint32)) *NestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetHandlerIDFunc = f
	return m.mock
}

//GetHandlerID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.NestedEvent interface
func (m *NestedEventMock) GetHandlerID() (r uint32) {
	counter := atomic.AddUint64(&m.GetHandlerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetHandlerIDCounter, 1)

	if len(m.GetHandlerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetHandlerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NestedEventMock.GetHandlerID.")
			return
		}

		result := m.GetHandlerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetHandlerID")
			return
		}

		r = result.r

		return
	}

	if m.GetHandlerIDMock.mainExpectation != nil {

		result := m.GetHandlerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetHandlerID")
		}

		r = result.r

		return
	}

	if m.GetHandlerIDFunc == nil {
		m.t.Fatalf("Unexpected call to NestedEventMock.GetHandlerID.")
		return
	}

	return m.GetHandlerIDFunc()
}

//GetHandlerIDMinimockCounter returns a count of NestedEventMock.GetHandlerIDFunc invocations
func (m *NestedEventMock) GetHandlerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDCounter)
}

//GetHandlerIDMinimockPreCounter returns the value of NestedEventMock.GetHandlerID invocations
func (m *NestedEventMock) GetHandlerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDPreCounter)
}

//GetHandlerIDFinished returns true if mock invocations count is ok
func (m *NestedEventMock) GetHandlerIDFinished() bool {
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

type mNestedEventMockGetParentElementID struct {
	mock              *NestedEventMock
	mainExpectation   *NestedEventMockGetParentElementIDExpectation
	expectationSeries []*NestedEventMockGetParentElementIDExpectation
}

type NestedEventMockGetParentElementIDExpectation struct {
	result *NestedEventMockGetParentElementIDResult
}

type NestedEventMockGetParentElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of NestedEvent.GetParentElementID is expected from 1 to Infinity times
func (m *mNestedEventMockGetParentElementID) Expect() *mNestedEventMockGetParentElementID {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetParentElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NestedEvent.GetParentElementID
func (m *mNestedEventMockGetParentElementID) Return(r uint32) *NestedEventMock {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NestedEventMockGetParentElementIDExpectation{}
	}
	m.mainExpectation.result = &NestedEventMockGetParentElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NestedEvent.GetParentElementID is expected once
func (m *mNestedEventMockGetParentElementID) ExpectOnce() *NestedEventMockGetParentElementIDExpectation {
	m.mock.GetParentElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &NestedEventMockGetParentElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NestedEventMockGetParentElementIDExpectation) Return(r uint32) {
	e.result = &NestedEventMockGetParentElementIDResult{r}
}

//Set uses given function f as a mock of NestedEvent.GetParentElementID method
func (m *mNestedEventMockGetParentElementID) Set(f func() (r uint32)) *NestedEventMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetParentElementIDFunc = f
	return m.mock
}

//GetParentElementID implements github.com/insolar/insolar/conveyor/interfaces/iadapter.NestedEvent interface
func (m *NestedEventMock) GetParentElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetParentElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetParentElementIDCounter, 1)

	if len(m.GetParentElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetParentElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NestedEventMock.GetParentElementID.")
			return
		}

		result := m.GetParentElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetParentElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetParentElementIDMock.mainExpectation != nil {

		result := m.GetParentElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NestedEventMock.GetParentElementID")
		}

		r = result.r

		return
	}

	if m.GetParentElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to NestedEventMock.GetParentElementID.")
		return
	}

	return m.GetParentElementIDFunc()
}

//GetParentElementIDMinimockCounter returns a count of NestedEventMock.GetParentElementIDFunc invocations
func (m *NestedEventMock) GetParentElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDCounter)
}

//GetParentElementIDMinimockPreCounter returns the value of NestedEventMock.GetParentElementID invocations
func (m *NestedEventMock) GetParentElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDPreCounter)
}

//GetParentElementIDFinished returns true if mock invocations count is ok
func (m *NestedEventMock) GetParentElementIDFinished() bool {
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
func (m *NestedEventMock) ValidateCallCounters() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetAdapterID")
	}

	if !m.GetEventPayloadFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetEventPayload")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetHandlerID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetParentElementID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NestedEventMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NestedEventMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NestedEventMock) MinimockFinish() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetAdapterID")
	}

	if !m.GetEventPayloadFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetEventPayload")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetHandlerID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to NestedEventMock.GetParentElementID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NestedEventMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NestedEventMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to NestedEventMock.GetAdapterID")
			}

			if !m.GetEventPayloadFinished() {
				m.t.Error("Expected call to NestedEventMock.GetEventPayload")
			}

			if !m.GetHandlerIDFinished() {
				m.t.Error("Expected call to NestedEventMock.GetHandlerID")
			}

			if !m.GetParentElementIDFinished() {
				m.t.Error("Expected call to NestedEventMock.GetParentElementID")
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
func (m *NestedEventMock) AllMocksCalled() bool {

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

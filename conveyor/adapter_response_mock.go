package conveyor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "AdapterResponse" can be found in github.com/insolar/insolar/conveyor
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	adapterid "github.com/insolar/insolar/conveyor/adapter/adapterid"
)

//AdapterResponseMock implements github.com/insolar/insolar/conveyor.AdapterResponse
type AdapterResponseMock struct {
	t minimock.Tester

	GetAdapterIDFunc       func() (r adapterid.ID)
	GetAdapterIDCounter    uint64
	GetAdapterIDPreCounter uint64
	GetAdapterIDMock       mAdapterResponseMockGetAdapterID

	GetElementIDFunc       func() (r uint32)
	GetElementIDCounter    uint64
	GetElementIDPreCounter uint64
	GetElementIDMock       mAdapterResponseMockGetElementID

	GetHandlerIDFunc       func() (r uint32)
	GetHandlerIDCounter    uint64
	GetHandlerIDPreCounter uint64
	GetHandlerIDMock       mAdapterResponseMockGetHandlerID

	GetRespPayloadFunc       func() (r interface{})
	GetRespPayloadCounter    uint64
	GetRespPayloadPreCounter uint64
	GetRespPayloadMock       mAdapterResponseMockGetRespPayload
}

//NewAdapterResponseMock returns a mock for github.com/insolar/insolar/conveyor.AdapterResponse
func NewAdapterResponseMock(t minimock.Tester) *AdapterResponseMock {
	m := &AdapterResponseMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAdapterIDMock = mAdapterResponseMockGetAdapterID{mock: m}
	m.GetElementIDMock = mAdapterResponseMockGetElementID{mock: m}
	m.GetHandlerIDMock = mAdapterResponseMockGetHandlerID{mock: m}
	m.GetRespPayloadMock = mAdapterResponseMockGetRespPayload{mock: m}

	return m
}

type mAdapterResponseMockGetAdapterID struct {
	mock              *AdapterResponseMock
	mainExpectation   *AdapterResponseMockGetAdapterIDExpectation
	expectationSeries []*AdapterResponseMockGetAdapterIDExpectation
}

type AdapterResponseMockGetAdapterIDExpectation struct {
	result *AdapterResponseMockGetAdapterIDResult
}

type AdapterResponseMockGetAdapterIDResult struct {
	r adapterid.ID
}

//Expect specifies that invocation of AdapterResponse.GetAdapterID is expected from 1 to Infinity times
func (m *mAdapterResponseMockGetAdapterID) Expect() *mAdapterResponseMockGetAdapterID {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetAdapterIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of AdapterResponse.GetAdapterID
func (m *mAdapterResponseMockGetAdapterID) Return(r adapterid.ID) *AdapterResponseMock {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetAdapterIDExpectation{}
	}
	m.mainExpectation.result = &AdapterResponseMockGetAdapterIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of AdapterResponse.GetAdapterID is expected once
func (m *mAdapterResponseMockGetAdapterID) ExpectOnce() *AdapterResponseMockGetAdapterIDExpectation {
	m.mock.GetAdapterIDFunc = nil
	m.mainExpectation = nil

	expectation := &AdapterResponseMockGetAdapterIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AdapterResponseMockGetAdapterIDExpectation) Return(r adapterid.ID) {
	e.result = &AdapterResponseMockGetAdapterIDResult{r}
}

//Set uses given function f as a mock of AdapterResponse.GetAdapterID method
func (m *mAdapterResponseMockGetAdapterID) Set(f func() (r adapterid.ID)) *AdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAdapterIDFunc = f
	return m.mock
}

//GetAdapterID implements github.com/insolar/insolar/conveyor.AdapterResponse interface
func (m *AdapterResponseMock) GetAdapterID() (r adapterid.ID) {
	counter := atomic.AddUint64(&m.GetAdapterIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetAdapterIDCounter, 1)

	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAdapterIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AdapterResponseMock.GetAdapterID.")
			return
		}

		result := m.GetAdapterIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetAdapterID")
			return
		}

		r = result.r

		return
	}

	if m.GetAdapterIDMock.mainExpectation != nil {

		result := m.GetAdapterIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetAdapterID")
		}

		r = result.r

		return
	}

	if m.GetAdapterIDFunc == nil {
		m.t.Fatalf("Unexpected call to AdapterResponseMock.GetAdapterID.")
		return
	}

	return m.GetAdapterIDFunc()
}

//GetAdapterIDMinimockCounter returns a count of AdapterResponseMock.GetAdapterIDFunc invocations
func (m *AdapterResponseMock) GetAdapterIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDCounter)
}

//GetAdapterIDMinimockPreCounter returns the value of AdapterResponseMock.GetAdapterID invocations
func (m *AdapterResponseMock) GetAdapterIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDPreCounter)
}

//GetAdapterIDFinished returns true if mock invocations count is ok
func (m *AdapterResponseMock) GetAdapterIDFinished() bool {
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

type mAdapterResponseMockGetElementID struct {
	mock              *AdapterResponseMock
	mainExpectation   *AdapterResponseMockGetElementIDExpectation
	expectationSeries []*AdapterResponseMockGetElementIDExpectation
}

type AdapterResponseMockGetElementIDExpectation struct {
	result *AdapterResponseMockGetElementIDResult
}

type AdapterResponseMockGetElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of AdapterResponse.GetElementID is expected from 1 to Infinity times
func (m *mAdapterResponseMockGetElementID) Expect() *mAdapterResponseMockGetElementID {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of AdapterResponse.GetElementID
func (m *mAdapterResponseMockGetElementID) Return(r uint32) *AdapterResponseMock {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetElementIDExpectation{}
	}
	m.mainExpectation.result = &AdapterResponseMockGetElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of AdapterResponse.GetElementID is expected once
func (m *mAdapterResponseMockGetElementID) ExpectOnce() *AdapterResponseMockGetElementIDExpectation {
	m.mock.GetElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &AdapterResponseMockGetElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AdapterResponseMockGetElementIDExpectation) Return(r uint32) {
	e.result = &AdapterResponseMockGetElementIDResult{r}
}

//Set uses given function f as a mock of AdapterResponse.GetElementID method
func (m *mAdapterResponseMockGetElementID) Set(f func() (r uint32)) *AdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetElementIDFunc = f
	return m.mock
}

//GetElementID implements github.com/insolar/insolar/conveyor.AdapterResponse interface
func (m *AdapterResponseMock) GetElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetElementIDCounter, 1)

	if len(m.GetElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AdapterResponseMock.GetElementID.")
			return
		}

		result := m.GetElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetElementIDMock.mainExpectation != nil {

		result := m.GetElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetElementID")
		}

		r = result.r

		return
	}

	if m.GetElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to AdapterResponseMock.GetElementID.")
		return
	}

	return m.GetElementIDFunc()
}

//GetElementIDMinimockCounter returns a count of AdapterResponseMock.GetElementIDFunc invocations
func (m *AdapterResponseMock) GetElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDCounter)
}

//GetElementIDMinimockPreCounter returns the value of AdapterResponseMock.GetElementID invocations
func (m *AdapterResponseMock) GetElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDPreCounter)
}

//GetElementIDFinished returns true if mock invocations count is ok
func (m *AdapterResponseMock) GetElementIDFinished() bool {
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

type mAdapterResponseMockGetHandlerID struct {
	mock              *AdapterResponseMock
	mainExpectation   *AdapterResponseMockGetHandlerIDExpectation
	expectationSeries []*AdapterResponseMockGetHandlerIDExpectation
}

type AdapterResponseMockGetHandlerIDExpectation struct {
	result *AdapterResponseMockGetHandlerIDResult
}

type AdapterResponseMockGetHandlerIDResult struct {
	r uint32
}

//Expect specifies that invocation of AdapterResponse.GetHandlerID is expected from 1 to Infinity times
func (m *mAdapterResponseMockGetHandlerID) Expect() *mAdapterResponseMockGetHandlerID {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetHandlerIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of AdapterResponse.GetHandlerID
func (m *mAdapterResponseMockGetHandlerID) Return(r uint32) *AdapterResponseMock {
	m.mock.GetHandlerIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetHandlerIDExpectation{}
	}
	m.mainExpectation.result = &AdapterResponseMockGetHandlerIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of AdapterResponse.GetHandlerID is expected once
func (m *mAdapterResponseMockGetHandlerID) ExpectOnce() *AdapterResponseMockGetHandlerIDExpectation {
	m.mock.GetHandlerIDFunc = nil
	m.mainExpectation = nil

	expectation := &AdapterResponseMockGetHandlerIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AdapterResponseMockGetHandlerIDExpectation) Return(r uint32) {
	e.result = &AdapterResponseMockGetHandlerIDResult{r}
}

//Set uses given function f as a mock of AdapterResponse.GetHandlerID method
func (m *mAdapterResponseMockGetHandlerID) Set(f func() (r uint32)) *AdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetHandlerIDFunc = f
	return m.mock
}

//GetHandlerID implements github.com/insolar/insolar/conveyor.AdapterResponse interface
func (m *AdapterResponseMock) GetHandlerID() (r uint32) {
	counter := atomic.AddUint64(&m.GetHandlerIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetHandlerIDCounter, 1)

	if len(m.GetHandlerIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetHandlerIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AdapterResponseMock.GetHandlerID.")
			return
		}

		result := m.GetHandlerIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetHandlerID")
			return
		}

		r = result.r

		return
	}

	if m.GetHandlerIDMock.mainExpectation != nil {

		result := m.GetHandlerIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetHandlerID")
		}

		r = result.r

		return
	}

	if m.GetHandlerIDFunc == nil {
		m.t.Fatalf("Unexpected call to AdapterResponseMock.GetHandlerID.")
		return
	}

	return m.GetHandlerIDFunc()
}

//GetHandlerIDMinimockCounter returns a count of AdapterResponseMock.GetHandlerIDFunc invocations
func (m *AdapterResponseMock) GetHandlerIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDCounter)
}

//GetHandlerIDMinimockPreCounter returns the value of AdapterResponseMock.GetHandlerID invocations
func (m *AdapterResponseMock) GetHandlerIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHandlerIDPreCounter)
}

//GetHandlerIDFinished returns true if mock invocations count is ok
func (m *AdapterResponseMock) GetHandlerIDFinished() bool {
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

type mAdapterResponseMockGetRespPayload struct {
	mock              *AdapterResponseMock
	mainExpectation   *AdapterResponseMockGetRespPayloadExpectation
	expectationSeries []*AdapterResponseMockGetRespPayloadExpectation
}

type AdapterResponseMockGetRespPayloadExpectation struct {
	result *AdapterResponseMockGetRespPayloadResult
}

type AdapterResponseMockGetRespPayloadResult struct {
	r interface{}
}

//Expect specifies that invocation of AdapterResponse.GetRespPayload is expected from 1 to Infinity times
func (m *mAdapterResponseMockGetRespPayload) Expect() *mAdapterResponseMockGetRespPayload {
	m.mock.GetRespPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetRespPayloadExpectation{}
	}

	return m
}

//Return specifies results of invocation of AdapterResponse.GetRespPayload
func (m *mAdapterResponseMockGetRespPayload) Return(r interface{}) *AdapterResponseMock {
	m.mock.GetRespPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AdapterResponseMockGetRespPayloadExpectation{}
	}
	m.mainExpectation.result = &AdapterResponseMockGetRespPayloadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of AdapterResponse.GetRespPayload is expected once
func (m *mAdapterResponseMockGetRespPayload) ExpectOnce() *AdapterResponseMockGetRespPayloadExpectation {
	m.mock.GetRespPayloadFunc = nil
	m.mainExpectation = nil

	expectation := &AdapterResponseMockGetRespPayloadExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AdapterResponseMockGetRespPayloadExpectation) Return(r interface{}) {
	e.result = &AdapterResponseMockGetRespPayloadResult{r}
}

//Set uses given function f as a mock of AdapterResponse.GetRespPayload method
func (m *mAdapterResponseMockGetRespPayload) Set(f func() (r interface{})) *AdapterResponseMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRespPayloadFunc = f
	return m.mock
}

//GetRespPayload implements github.com/insolar/insolar/conveyor.AdapterResponse interface
func (m *AdapterResponseMock) GetRespPayload() (r interface{}) {
	counter := atomic.AddUint64(&m.GetRespPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.GetRespPayloadCounter, 1)

	if len(m.GetRespPayloadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRespPayloadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AdapterResponseMock.GetRespPayload.")
			return
		}

		result := m.GetRespPayloadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetRespPayload")
			return
		}

		r = result.r

		return
	}

	if m.GetRespPayloadMock.mainExpectation != nil {

		result := m.GetRespPayloadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AdapterResponseMock.GetRespPayload")
		}

		r = result.r

		return
	}

	if m.GetRespPayloadFunc == nil {
		m.t.Fatalf("Unexpected call to AdapterResponseMock.GetRespPayload.")
		return
	}

	return m.GetRespPayloadFunc()
}

//GetRespPayloadMinimockCounter returns a count of AdapterResponseMock.GetRespPayloadFunc invocations
func (m *AdapterResponseMock) GetRespPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRespPayloadCounter)
}

//GetRespPayloadMinimockPreCounter returns the value of AdapterResponseMock.GetRespPayload invocations
func (m *AdapterResponseMock) GetRespPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRespPayloadPreCounter)
}

//GetRespPayloadFinished returns true if mock invocations count is ok
func (m *AdapterResponseMock) GetRespPayloadFinished() bool {
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
func (m *AdapterResponseMock) ValidateCallCounters() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetAdapterID")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetElementID")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetHandlerID")
	}

	if !m.GetRespPayloadFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetRespPayload")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AdapterResponseMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AdapterResponseMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AdapterResponseMock) MinimockFinish() {

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetAdapterID")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetElementID")
	}

	if !m.GetHandlerIDFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetHandlerID")
	}

	if !m.GetRespPayloadFinished() {
		m.t.Fatal("Expected call to AdapterResponseMock.GetRespPayload")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AdapterResponseMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AdapterResponseMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to AdapterResponseMock.GetAdapterID")
			}

			if !m.GetElementIDFinished() {
				m.t.Error("Expected call to AdapterResponseMock.GetElementID")
			}

			if !m.GetHandlerIDFinished() {
				m.t.Error("Expected call to AdapterResponseMock.GetHandlerID")
			}

			if !m.GetRespPayloadFinished() {
				m.t.Error("Expected call to AdapterResponseMock.GetRespPayload")
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
func (m *AdapterResponseMock) AllMocksCalled() bool {

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

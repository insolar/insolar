package fsm

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SlotElementReadOnly" can be found in github.com/insolar/insolar/conveyor/fsm
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//SlotElementReadOnlyMock implements github.com/insolar/insolar/conveyor/fsm.SlotElementReadOnly
type SlotElementReadOnlyMock struct {
	t minimock.Tester

	GetElementIDFunc       func() (r uint32)
	GetElementIDCounter    uint64
	GetElementIDPreCounter uint64
	GetElementIDMock       mSlotElementReadOnlyMockGetElementID

	GetNodeIDFunc       func() (r uint32)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mSlotElementReadOnlyMockGetNodeID

	GetStateFunc       func() (r StateID)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mSlotElementReadOnlyMockGetState

	GetTypeFunc       func() (r ID)
	GetTypeCounter    uint64
	GetTypePreCounter uint64
	GetTypeMock       mSlotElementReadOnlyMockGetType
}

//NewSlotElementReadOnlyMock returns a mock for github.com/insolar/insolar/conveyor/fsm.SlotElementReadOnly
func NewSlotElementReadOnlyMock(t minimock.Tester) *SlotElementReadOnlyMock {
	m := &SlotElementReadOnlyMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetElementIDMock = mSlotElementReadOnlyMockGetElementID{mock: m}
	m.GetNodeIDMock = mSlotElementReadOnlyMockGetNodeID{mock: m}
	m.GetStateMock = mSlotElementReadOnlyMockGetState{mock: m}
	m.GetTypeMock = mSlotElementReadOnlyMockGetType{mock: m}

	return m
}

type mSlotElementReadOnlyMockGetElementID struct {
	mock              *SlotElementReadOnlyMock
	mainExpectation   *SlotElementReadOnlyMockGetElementIDExpectation
	expectationSeries []*SlotElementReadOnlyMockGetElementIDExpectation
}

type SlotElementReadOnlyMockGetElementIDExpectation struct {
	result *SlotElementReadOnlyMockGetElementIDResult
}

type SlotElementReadOnlyMockGetElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementReadOnly.GetElementID is expected from 1 to Infinity times
func (m *mSlotElementReadOnlyMockGetElementID) Expect() *mSlotElementReadOnlyMockGetElementID {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementReadOnly.GetElementID
func (m *mSlotElementReadOnlyMockGetElementID) Return(r uint32) *SlotElementReadOnlyMock {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetElementIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementReadOnlyMockGetElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementReadOnly.GetElementID is expected once
func (m *mSlotElementReadOnlyMockGetElementID) ExpectOnce() *SlotElementReadOnlyMockGetElementIDExpectation {
	m.mock.GetElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementReadOnlyMockGetElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementReadOnlyMockGetElementIDExpectation) Return(r uint32) {
	e.result = &SlotElementReadOnlyMockGetElementIDResult{r}
}

//Set uses given function f as a mock of SlotElementReadOnly.GetElementID method
func (m *mSlotElementReadOnlyMockGetElementID) Set(f func() (r uint32)) *SlotElementReadOnlyMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetElementIDFunc = f
	return m.mock
}

//GetElementID implements github.com/insolar/insolar/conveyor/fsm.SlotElementReadOnly interface
func (m *SlotElementReadOnlyMock) GetElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetElementIDCounter, 1)

	if len(m.GetElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetElementID.")
			return
		}

		result := m.GetElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetElementIDMock.mainExpectation != nil {

		result := m.GetElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetElementID")
		}

		r = result.r

		return
	}

	if m.GetElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetElementID.")
		return
	}

	return m.GetElementIDFunc()
}

//GetElementIDMinimockCounter returns a count of SlotElementReadOnlyMock.GetElementIDFunc invocations
func (m *SlotElementReadOnlyMock) GetElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDCounter)
}

//GetElementIDMinimockPreCounter returns the value of SlotElementReadOnlyMock.GetElementID invocations
func (m *SlotElementReadOnlyMock) GetElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDPreCounter)
}

//GetElementIDFinished returns true if mock invocations count is ok
func (m *SlotElementReadOnlyMock) GetElementIDFinished() bool {
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

type mSlotElementReadOnlyMockGetNodeID struct {
	mock              *SlotElementReadOnlyMock
	mainExpectation   *SlotElementReadOnlyMockGetNodeIDExpectation
	expectationSeries []*SlotElementReadOnlyMockGetNodeIDExpectation
}

type SlotElementReadOnlyMockGetNodeIDExpectation struct {
	result *SlotElementReadOnlyMockGetNodeIDResult
}

type SlotElementReadOnlyMockGetNodeIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementReadOnly.GetNodeID is expected from 1 to Infinity times
func (m *mSlotElementReadOnlyMockGetNodeID) Expect() *mSlotElementReadOnlyMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementReadOnly.GetNodeID
func (m *mSlotElementReadOnlyMockGetNodeID) Return(r uint32) *SlotElementReadOnlyMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementReadOnlyMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementReadOnly.GetNodeID is expected once
func (m *mSlotElementReadOnlyMockGetNodeID) ExpectOnce() *SlotElementReadOnlyMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementReadOnlyMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementReadOnlyMockGetNodeIDExpectation) Return(r uint32) {
	e.result = &SlotElementReadOnlyMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of SlotElementReadOnly.GetNodeID method
func (m *mSlotElementReadOnlyMockGetNodeID) Set(f func() (r uint32)) *SlotElementReadOnlyMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/conveyor/fsm.SlotElementReadOnly interface
func (m *SlotElementReadOnlyMock) GetNodeID() (r uint32) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of SlotElementReadOnlyMock.GetNodeIDFunc invocations
func (m *SlotElementReadOnlyMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of SlotElementReadOnlyMock.GetNodeID invocations
func (m *SlotElementReadOnlyMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *SlotElementReadOnlyMock) GetNodeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeIDCounter) == uint64(len(m.GetNodeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeIDFunc != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	return true
}

type mSlotElementReadOnlyMockGetState struct {
	mock              *SlotElementReadOnlyMock
	mainExpectation   *SlotElementReadOnlyMockGetStateExpectation
	expectationSeries []*SlotElementReadOnlyMockGetStateExpectation
}

type SlotElementReadOnlyMockGetStateExpectation struct {
	result *SlotElementReadOnlyMockGetStateResult
}

type SlotElementReadOnlyMockGetStateResult struct {
	r StateID
}

//Expect specifies that invocation of SlotElementReadOnly.GetState is expected from 1 to Infinity times
func (m *mSlotElementReadOnlyMockGetState) Expect() *mSlotElementReadOnlyMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementReadOnly.GetState
func (m *mSlotElementReadOnlyMockGetState) Return(r StateID) *SlotElementReadOnlyMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetStateExpectation{}
	}
	m.mainExpectation.result = &SlotElementReadOnlyMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementReadOnly.GetState is expected once
func (m *mSlotElementReadOnlyMockGetState) ExpectOnce() *SlotElementReadOnlyMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementReadOnlyMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementReadOnlyMockGetStateExpectation) Return(r StateID) {
	e.result = &SlotElementReadOnlyMockGetStateResult{r}
}

//Set uses given function f as a mock of SlotElementReadOnly.GetState method
func (m *mSlotElementReadOnlyMockGetState) Set(f func() (r StateID)) *SlotElementReadOnlyMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/conveyor/fsm.SlotElementReadOnly interface
func (m *SlotElementReadOnlyMock) GetState() (r StateID) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of SlotElementReadOnlyMock.GetStateFunc invocations
func (m *SlotElementReadOnlyMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of SlotElementReadOnlyMock.GetState invocations
func (m *SlotElementReadOnlyMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *SlotElementReadOnlyMock) GetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStateCounter) == uint64(len(m.GetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStateFunc != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	return true
}

type mSlotElementReadOnlyMockGetType struct {
	mock              *SlotElementReadOnlyMock
	mainExpectation   *SlotElementReadOnlyMockGetTypeExpectation
	expectationSeries []*SlotElementReadOnlyMockGetTypeExpectation
}

type SlotElementReadOnlyMockGetTypeExpectation struct {
	result *SlotElementReadOnlyMockGetTypeResult
}

type SlotElementReadOnlyMockGetTypeResult struct {
	r ID
}

//Expect specifies that invocation of SlotElementReadOnly.GetType is expected from 1 to Infinity times
func (m *mSlotElementReadOnlyMockGetType) Expect() *mSlotElementReadOnlyMockGetType {
	m.mock.GetTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementReadOnly.GetType
func (m *mSlotElementReadOnlyMockGetType) Return(r ID) *SlotElementReadOnlyMock {
	m.mock.GetTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementReadOnlyMockGetTypeExpectation{}
	}
	m.mainExpectation.result = &SlotElementReadOnlyMockGetTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementReadOnly.GetType is expected once
func (m *mSlotElementReadOnlyMockGetType) ExpectOnce() *SlotElementReadOnlyMockGetTypeExpectation {
	m.mock.GetTypeFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementReadOnlyMockGetTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementReadOnlyMockGetTypeExpectation) Return(r ID) {
	e.result = &SlotElementReadOnlyMockGetTypeResult{r}
}

//Set uses given function f as a mock of SlotElementReadOnly.GetType method
func (m *mSlotElementReadOnlyMockGetType) Set(f func() (r ID)) *SlotElementReadOnlyMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTypeFunc = f
	return m.mock
}

//GetType implements github.com/insolar/insolar/conveyor/fsm.SlotElementReadOnly interface
func (m *SlotElementReadOnlyMock) GetType() (r ID) {
	counter := atomic.AddUint64(&m.GetTypePreCounter, 1)
	defer atomic.AddUint64(&m.GetTypeCounter, 1)

	if len(m.GetTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetType.")
			return
		}

		result := m.GetTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetType")
			return
		}

		r = result.r

		return
	}

	if m.GetTypeMock.mainExpectation != nil {

		result := m.GetTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementReadOnlyMock.GetType")
		}

		r = result.r

		return
	}

	if m.GetTypeFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementReadOnlyMock.GetType.")
		return
	}

	return m.GetTypeFunc()
}

//GetTypeMinimockCounter returns a count of SlotElementReadOnlyMock.GetTypeFunc invocations
func (m *SlotElementReadOnlyMock) GetTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypeCounter)
}

//GetTypeMinimockPreCounter returns the value of SlotElementReadOnlyMock.GetType invocations
func (m *SlotElementReadOnlyMock) GetTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypePreCounter)
}

//GetTypeFinished returns true if mock invocations count is ok
func (m *SlotElementReadOnlyMock) GetTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetTypeCounter) == uint64(len(m.GetTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetTypeFunc != nil {
		return atomic.LoadUint64(&m.GetTypeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotElementReadOnlyMock) ValidateCallCounters() {

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetElementID")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetNodeID")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetState")
	}

	if !m.GetTypeFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetType")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotElementReadOnlyMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SlotElementReadOnlyMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SlotElementReadOnlyMock) MinimockFinish() {

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetElementID")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetNodeID")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetState")
	}

	if !m.GetTypeFinished() {
		m.t.Fatal("Expected call to SlotElementReadOnlyMock.GetType")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SlotElementReadOnlyMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SlotElementReadOnlyMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetElementIDFinished()
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.GetTypeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetElementIDFinished() {
				m.t.Error("Expected call to SlotElementReadOnlyMock.GetElementID")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to SlotElementReadOnlyMock.GetNodeID")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to SlotElementReadOnlyMock.GetState")
			}

			if !m.GetTypeFinished() {
				m.t.Error("Expected call to SlotElementReadOnlyMock.GetType")
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
func (m *SlotElementReadOnlyMock) AllMocksCalled() bool {

	if !m.GetElementIDFinished() {
		return false
	}

	if !m.GetNodeIDFinished() {
		return false
	}

	if !m.GetStateFinished() {
		return false
	}

	if !m.GetTypeFinished() {
		return false
	}

	return true
}

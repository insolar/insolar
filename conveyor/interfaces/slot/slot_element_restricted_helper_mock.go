package slot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SlotElementRestrictedHelper" can be found in github.com/insolar/insolar/conveyor/interfaces/slot
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	fsm "github.com/insolar/insolar/conveyor/interfaces/fsm"
)

//SlotElementRestrictedHelperMock implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper
type SlotElementRestrictedHelperMock struct {
	t minimock.Tester

	GetElementIDFunc       func() (r uint32)
	GetElementIDCounter    uint64
	GetElementIDPreCounter uint64
	GetElementIDMock       mSlotElementRestrictedHelperMockGetElementID

	GetInputEventFunc       func() (r interface{})
	GetInputEventCounter    uint64
	GetInputEventPreCounter uint64
	GetInputEventMock       mSlotElementRestrictedHelperMockGetInputEvent

	GetNodeIDFunc       func() (r uint32)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mSlotElementRestrictedHelperMockGetNodeID

	GetParentElementIDFunc       func() (r uint32)
	GetParentElementIDCounter    uint64
	GetParentElementIDPreCounter uint64
	GetParentElementIDMock       mSlotElementRestrictedHelperMockGetParentElementID

	GetPayloadFunc       func() (r interface{})
	GetPayloadCounter    uint64
	GetPayloadPreCounter uint64
	GetPayloadMock       mSlotElementRestrictedHelperMockGetPayload

	GetStateFunc       func() (r fsm.StateID)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mSlotElementRestrictedHelperMockGetState

	GetTypeFunc       func() (r fsm.ID)
	GetTypeCounter    uint64
	GetTypePreCounter uint64
	GetTypeMock       mSlotElementRestrictedHelperMockGetType

	LeaveSequenceFunc       func()
	LeaveSequenceCounter    uint64
	LeaveSequencePreCounter uint64
	LeaveSequenceMock       mSlotElementRestrictedHelperMockLeaveSequence

	ReactivateFunc       func()
	ReactivateCounter    uint64
	ReactivatePreCounter uint64
	ReactivateMock       mSlotElementRestrictedHelperMockReactivate
}

//NewSlotElementRestrictedHelperMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper
func NewSlotElementRestrictedHelperMock(t minimock.Tester) *SlotElementRestrictedHelperMock {
	m := &SlotElementRestrictedHelperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetElementIDMock = mSlotElementRestrictedHelperMockGetElementID{mock: m}
	m.GetInputEventMock = mSlotElementRestrictedHelperMockGetInputEvent{mock: m}
	m.GetNodeIDMock = mSlotElementRestrictedHelperMockGetNodeID{mock: m}
	m.GetParentElementIDMock = mSlotElementRestrictedHelperMockGetParentElementID{mock: m}
	m.GetPayloadMock = mSlotElementRestrictedHelperMockGetPayload{mock: m}
	m.GetStateMock = mSlotElementRestrictedHelperMockGetState{mock: m}
	m.GetTypeMock = mSlotElementRestrictedHelperMockGetType{mock: m}
	m.LeaveSequenceMock = mSlotElementRestrictedHelperMockLeaveSequence{mock: m}
	m.ReactivateMock = mSlotElementRestrictedHelperMockReactivate{mock: m}

	return m
}

type mSlotElementRestrictedHelperMockGetElementID struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockGetElementIDExpectation
	expectationSeries []*SlotElementRestrictedHelperMockGetElementIDExpectation
}

type SlotElementRestrictedHelperMockGetElementIDExpectation struct {
	result *SlotElementRestrictedHelperMockGetElementIDResult
}

type SlotElementRestrictedHelperMockGetElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementRestrictedHelper.GetElementID is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockGetElementID) Expect() *mSlotElementRestrictedHelperMockGetElementID {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.GetElementID
func (m *mSlotElementRestrictedHelperMockGetElementID) Return(r uint32) *SlotElementRestrictedHelperMock {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetElementIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementRestrictedHelperMockGetElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.GetElementID is expected once
func (m *mSlotElementRestrictedHelperMockGetElementID) ExpectOnce() *SlotElementRestrictedHelperMockGetElementIDExpectation {
	m.mock.GetElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockGetElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementRestrictedHelperMockGetElementIDExpectation) Return(r uint32) {
	e.result = &SlotElementRestrictedHelperMockGetElementIDResult{r}
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.GetElementID method
func (m *mSlotElementRestrictedHelperMockGetElementID) Set(f func() (r uint32)) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetElementIDFunc = f
	return m.mock
}

//GetElementID implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) GetElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetElementIDCounter, 1)

	if len(m.GetElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetElementID.")
			return
		}

		result := m.GetElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetElementIDMock.mainExpectation != nil {

		result := m.GetElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetElementID")
		}

		r = result.r

		return
	}

	if m.GetElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetElementID.")
		return
	}

	return m.GetElementIDFunc()
}

//GetElementIDMinimockCounter returns a count of SlotElementRestrictedHelperMock.GetElementIDFunc invocations
func (m *SlotElementRestrictedHelperMock) GetElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDCounter)
}

//GetElementIDMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.GetElementID invocations
func (m *SlotElementRestrictedHelperMock) GetElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDPreCounter)
}

//GetElementIDFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) GetElementIDFinished() bool {
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

type mSlotElementRestrictedHelperMockGetInputEvent struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockGetInputEventExpectation
	expectationSeries []*SlotElementRestrictedHelperMockGetInputEventExpectation
}

type SlotElementRestrictedHelperMockGetInputEventExpectation struct {
	result *SlotElementRestrictedHelperMockGetInputEventResult
}

type SlotElementRestrictedHelperMockGetInputEventResult struct {
	r interface{}
}

//Expect specifies that invocation of SlotElementRestrictedHelper.GetInputEvent is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockGetInputEvent) Expect() *mSlotElementRestrictedHelperMockGetInputEvent {
	m.mock.GetInputEventFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetInputEventExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.GetInputEvent
func (m *mSlotElementRestrictedHelperMockGetInputEvent) Return(r interface{}) *SlotElementRestrictedHelperMock {
	m.mock.GetInputEventFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetInputEventExpectation{}
	}
	m.mainExpectation.result = &SlotElementRestrictedHelperMockGetInputEventResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.GetInputEvent is expected once
func (m *mSlotElementRestrictedHelperMockGetInputEvent) ExpectOnce() *SlotElementRestrictedHelperMockGetInputEventExpectation {
	m.mock.GetInputEventFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockGetInputEventExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementRestrictedHelperMockGetInputEventExpectation) Return(r interface{}) {
	e.result = &SlotElementRestrictedHelperMockGetInputEventResult{r}
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.GetInputEvent method
func (m *mSlotElementRestrictedHelperMockGetInputEvent) Set(f func() (r interface{})) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetInputEventFunc = f
	return m.mock
}

//GetInputEvent implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) GetInputEvent() (r interface{}) {
	counter := atomic.AddUint64(&m.GetInputEventPreCounter, 1)
	defer atomic.AddUint64(&m.GetInputEventCounter, 1)

	if len(m.GetInputEventMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetInputEventMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetInputEvent.")
			return
		}

		result := m.GetInputEventMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetInputEvent")
			return
		}

		r = result.r

		return
	}

	if m.GetInputEventMock.mainExpectation != nil {

		result := m.GetInputEventMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetInputEvent")
		}

		r = result.r

		return
	}

	if m.GetInputEventFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetInputEvent.")
		return
	}

	return m.GetInputEventFunc()
}

//GetInputEventMinimockCounter returns a count of SlotElementRestrictedHelperMock.GetInputEventFunc invocations
func (m *SlotElementRestrictedHelperMock) GetInputEventMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetInputEventCounter)
}

//GetInputEventMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.GetInputEvent invocations
func (m *SlotElementRestrictedHelperMock) GetInputEventMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetInputEventPreCounter)
}

//GetInputEventFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) GetInputEventFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetInputEventMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetInputEventCounter) == uint64(len(m.GetInputEventMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetInputEventMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetInputEventCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetInputEventFunc != nil {
		return atomic.LoadUint64(&m.GetInputEventCounter) > 0
	}

	return true
}

type mSlotElementRestrictedHelperMockGetNodeID struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockGetNodeIDExpectation
	expectationSeries []*SlotElementRestrictedHelperMockGetNodeIDExpectation
}

type SlotElementRestrictedHelperMockGetNodeIDExpectation struct {
	result *SlotElementRestrictedHelperMockGetNodeIDResult
}

type SlotElementRestrictedHelperMockGetNodeIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementRestrictedHelper.GetNodeID is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockGetNodeID) Expect() *mSlotElementRestrictedHelperMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.GetNodeID
func (m *mSlotElementRestrictedHelperMockGetNodeID) Return(r uint32) *SlotElementRestrictedHelperMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementRestrictedHelperMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.GetNodeID is expected once
func (m *mSlotElementRestrictedHelperMockGetNodeID) ExpectOnce() *SlotElementRestrictedHelperMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementRestrictedHelperMockGetNodeIDExpectation) Return(r uint32) {
	e.result = &SlotElementRestrictedHelperMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.GetNodeID method
func (m *mSlotElementRestrictedHelperMockGetNodeID) Set(f func() (r uint32)) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) GetNodeID() (r uint32) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of SlotElementRestrictedHelperMock.GetNodeIDFunc invocations
func (m *SlotElementRestrictedHelperMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.GetNodeID invocations
func (m *SlotElementRestrictedHelperMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) GetNodeIDFinished() bool {
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

type mSlotElementRestrictedHelperMockGetParentElementID struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockGetParentElementIDExpectation
	expectationSeries []*SlotElementRestrictedHelperMockGetParentElementIDExpectation
}

type SlotElementRestrictedHelperMockGetParentElementIDExpectation struct {
	result *SlotElementRestrictedHelperMockGetParentElementIDResult
}

type SlotElementRestrictedHelperMockGetParentElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementRestrictedHelper.GetParentElementID is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockGetParentElementID) Expect() *mSlotElementRestrictedHelperMockGetParentElementID {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetParentElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.GetParentElementID
func (m *mSlotElementRestrictedHelperMockGetParentElementID) Return(r uint32) *SlotElementRestrictedHelperMock {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetParentElementIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementRestrictedHelperMockGetParentElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.GetParentElementID is expected once
func (m *mSlotElementRestrictedHelperMockGetParentElementID) ExpectOnce() *SlotElementRestrictedHelperMockGetParentElementIDExpectation {
	m.mock.GetParentElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockGetParentElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementRestrictedHelperMockGetParentElementIDExpectation) Return(r uint32) {
	e.result = &SlotElementRestrictedHelperMockGetParentElementIDResult{r}
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.GetParentElementID method
func (m *mSlotElementRestrictedHelperMockGetParentElementID) Set(f func() (r uint32)) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetParentElementIDFunc = f
	return m.mock
}

//GetParentElementID implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) GetParentElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetParentElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetParentElementIDCounter, 1)

	if len(m.GetParentElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetParentElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetParentElementID.")
			return
		}

		result := m.GetParentElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetParentElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetParentElementIDMock.mainExpectation != nil {

		result := m.GetParentElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetParentElementID")
		}

		r = result.r

		return
	}

	if m.GetParentElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetParentElementID.")
		return
	}

	return m.GetParentElementIDFunc()
}

//GetParentElementIDMinimockCounter returns a count of SlotElementRestrictedHelperMock.GetParentElementIDFunc invocations
func (m *SlotElementRestrictedHelperMock) GetParentElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDCounter)
}

//GetParentElementIDMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.GetParentElementID invocations
func (m *SlotElementRestrictedHelperMock) GetParentElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDPreCounter)
}

//GetParentElementIDFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) GetParentElementIDFinished() bool {
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

type mSlotElementRestrictedHelperMockGetPayload struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockGetPayloadExpectation
	expectationSeries []*SlotElementRestrictedHelperMockGetPayloadExpectation
}

type SlotElementRestrictedHelperMockGetPayloadExpectation struct {
	result *SlotElementRestrictedHelperMockGetPayloadResult
}

type SlotElementRestrictedHelperMockGetPayloadResult struct {
	r interface{}
}

//Expect specifies that invocation of SlotElementRestrictedHelper.GetPayload is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockGetPayload) Expect() *mSlotElementRestrictedHelperMockGetPayload {
	m.mock.GetPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetPayloadExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.GetPayload
func (m *mSlotElementRestrictedHelperMockGetPayload) Return(r interface{}) *SlotElementRestrictedHelperMock {
	m.mock.GetPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetPayloadExpectation{}
	}
	m.mainExpectation.result = &SlotElementRestrictedHelperMockGetPayloadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.GetPayload is expected once
func (m *mSlotElementRestrictedHelperMockGetPayload) ExpectOnce() *SlotElementRestrictedHelperMockGetPayloadExpectation {
	m.mock.GetPayloadFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockGetPayloadExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementRestrictedHelperMockGetPayloadExpectation) Return(r interface{}) {
	e.result = &SlotElementRestrictedHelperMockGetPayloadResult{r}
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.GetPayload method
func (m *mSlotElementRestrictedHelperMockGetPayload) Set(f func() (r interface{})) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPayloadFunc = f
	return m.mock
}

//GetPayload implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) GetPayload() (r interface{}) {
	counter := atomic.AddUint64(&m.GetPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.GetPayloadCounter, 1)

	if len(m.GetPayloadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPayloadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetPayload.")
			return
		}

		result := m.GetPayloadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetPayload")
			return
		}

		r = result.r

		return
	}

	if m.GetPayloadMock.mainExpectation != nil {

		result := m.GetPayloadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetPayload")
		}

		r = result.r

		return
	}

	if m.GetPayloadFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetPayload.")
		return
	}

	return m.GetPayloadFunc()
}

//GetPayloadMinimockCounter returns a count of SlotElementRestrictedHelperMock.GetPayloadFunc invocations
func (m *SlotElementRestrictedHelperMock) GetPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPayloadCounter)
}

//GetPayloadMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.GetPayload invocations
func (m *SlotElementRestrictedHelperMock) GetPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPayloadPreCounter)
}

//GetPayloadFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) GetPayloadFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPayloadMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPayloadCounter) == uint64(len(m.GetPayloadMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPayloadMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPayloadCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPayloadFunc != nil {
		return atomic.LoadUint64(&m.GetPayloadCounter) > 0
	}

	return true
}

type mSlotElementRestrictedHelperMockGetState struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockGetStateExpectation
	expectationSeries []*SlotElementRestrictedHelperMockGetStateExpectation
}

type SlotElementRestrictedHelperMockGetStateExpectation struct {
	result *SlotElementRestrictedHelperMockGetStateResult
}

type SlotElementRestrictedHelperMockGetStateResult struct {
	r fsm.StateID
}

//Expect specifies that invocation of SlotElementRestrictedHelper.GetState is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockGetState) Expect() *mSlotElementRestrictedHelperMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.GetState
func (m *mSlotElementRestrictedHelperMockGetState) Return(r fsm.StateID) *SlotElementRestrictedHelperMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetStateExpectation{}
	}
	m.mainExpectation.result = &SlotElementRestrictedHelperMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.GetState is expected once
func (m *mSlotElementRestrictedHelperMockGetState) ExpectOnce() *SlotElementRestrictedHelperMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementRestrictedHelperMockGetStateExpectation) Return(r fsm.StateID) {
	e.result = &SlotElementRestrictedHelperMockGetStateResult{r}
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.GetState method
func (m *mSlotElementRestrictedHelperMockGetState) Set(f func() (r fsm.StateID)) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) GetState() (r fsm.StateID) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of SlotElementRestrictedHelperMock.GetStateFunc invocations
func (m *SlotElementRestrictedHelperMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.GetState invocations
func (m *SlotElementRestrictedHelperMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) GetStateFinished() bool {
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

type mSlotElementRestrictedHelperMockGetType struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockGetTypeExpectation
	expectationSeries []*SlotElementRestrictedHelperMockGetTypeExpectation
}

type SlotElementRestrictedHelperMockGetTypeExpectation struct {
	result *SlotElementRestrictedHelperMockGetTypeResult
}

type SlotElementRestrictedHelperMockGetTypeResult struct {
	r fsm.ID
}

//Expect specifies that invocation of SlotElementRestrictedHelper.GetType is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockGetType) Expect() *mSlotElementRestrictedHelperMockGetType {
	m.mock.GetTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.GetType
func (m *mSlotElementRestrictedHelperMockGetType) Return(r fsm.ID) *SlotElementRestrictedHelperMock {
	m.mock.GetTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockGetTypeExpectation{}
	}
	m.mainExpectation.result = &SlotElementRestrictedHelperMockGetTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.GetType is expected once
func (m *mSlotElementRestrictedHelperMockGetType) ExpectOnce() *SlotElementRestrictedHelperMockGetTypeExpectation {
	m.mock.GetTypeFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockGetTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementRestrictedHelperMockGetTypeExpectation) Return(r fsm.ID) {
	e.result = &SlotElementRestrictedHelperMockGetTypeResult{r}
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.GetType method
func (m *mSlotElementRestrictedHelperMockGetType) Set(f func() (r fsm.ID)) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTypeFunc = f
	return m.mock
}

//GetType implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) GetType() (r fsm.ID) {
	counter := atomic.AddUint64(&m.GetTypePreCounter, 1)
	defer atomic.AddUint64(&m.GetTypeCounter, 1)

	if len(m.GetTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetType.")
			return
		}

		result := m.GetTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetType")
			return
		}

		r = result.r

		return
	}

	if m.GetTypeMock.mainExpectation != nil {

		result := m.GetTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementRestrictedHelperMock.GetType")
		}

		r = result.r

		return
	}

	if m.GetTypeFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.GetType.")
		return
	}

	return m.GetTypeFunc()
}

//GetTypeMinimockCounter returns a count of SlotElementRestrictedHelperMock.GetTypeFunc invocations
func (m *SlotElementRestrictedHelperMock) GetTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypeCounter)
}

//GetTypeMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.GetType invocations
func (m *SlotElementRestrictedHelperMock) GetTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypePreCounter)
}

//GetTypeFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) GetTypeFinished() bool {
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

type mSlotElementRestrictedHelperMockLeaveSequence struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockLeaveSequenceExpectation
	expectationSeries []*SlotElementRestrictedHelperMockLeaveSequenceExpectation
}

type SlotElementRestrictedHelperMockLeaveSequenceExpectation struct {
}

//Expect specifies that invocation of SlotElementRestrictedHelper.LeaveSequence is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockLeaveSequence) Expect() *mSlotElementRestrictedHelperMockLeaveSequence {
	m.mock.LeaveSequenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockLeaveSequenceExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.LeaveSequence
func (m *mSlotElementRestrictedHelperMockLeaveSequence) Return() *SlotElementRestrictedHelperMock {
	m.mock.LeaveSequenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockLeaveSequenceExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.LeaveSequence is expected once
func (m *mSlotElementRestrictedHelperMockLeaveSequence) ExpectOnce() *SlotElementRestrictedHelperMockLeaveSequenceExpectation {
	m.mock.LeaveSequenceFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockLeaveSequenceExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.LeaveSequence method
func (m *mSlotElementRestrictedHelperMockLeaveSequence) Set(f func()) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LeaveSequenceFunc = f
	return m.mock
}

//LeaveSequence implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) LeaveSequence() {
	counter := atomic.AddUint64(&m.LeaveSequencePreCounter, 1)
	defer atomic.AddUint64(&m.LeaveSequenceCounter, 1)

	if len(m.LeaveSequenceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LeaveSequenceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.LeaveSequence.")
			return
		}

		return
	}

	if m.LeaveSequenceMock.mainExpectation != nil {

		return
	}

	if m.LeaveSequenceFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.LeaveSequence.")
		return
	}

	m.LeaveSequenceFunc()
}

//LeaveSequenceMinimockCounter returns a count of SlotElementRestrictedHelperMock.LeaveSequenceFunc invocations
func (m *SlotElementRestrictedHelperMock) LeaveSequenceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LeaveSequenceCounter)
}

//LeaveSequenceMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.LeaveSequence invocations
func (m *SlotElementRestrictedHelperMock) LeaveSequenceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LeaveSequencePreCounter)
}

//LeaveSequenceFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) LeaveSequenceFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LeaveSequenceMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LeaveSequenceCounter) == uint64(len(m.LeaveSequenceMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LeaveSequenceMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LeaveSequenceCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LeaveSequenceFunc != nil {
		return atomic.LoadUint64(&m.LeaveSequenceCounter) > 0
	}

	return true
}

type mSlotElementRestrictedHelperMockReactivate struct {
	mock              *SlotElementRestrictedHelperMock
	mainExpectation   *SlotElementRestrictedHelperMockReactivateExpectation
	expectationSeries []*SlotElementRestrictedHelperMockReactivateExpectation
}

type SlotElementRestrictedHelperMockReactivateExpectation struct {
}

//Expect specifies that invocation of SlotElementRestrictedHelper.Reactivate is expected from 1 to Infinity times
func (m *mSlotElementRestrictedHelperMockReactivate) Expect() *mSlotElementRestrictedHelperMockReactivate {
	m.mock.ReactivateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockReactivateExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementRestrictedHelper.Reactivate
func (m *mSlotElementRestrictedHelperMockReactivate) Return() *SlotElementRestrictedHelperMock {
	m.mock.ReactivateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementRestrictedHelperMockReactivateExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementRestrictedHelper.Reactivate is expected once
func (m *mSlotElementRestrictedHelperMockReactivate) ExpectOnce() *SlotElementRestrictedHelperMockReactivateExpectation {
	m.mock.ReactivateFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementRestrictedHelperMockReactivateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of SlotElementRestrictedHelper.Reactivate method
func (m *mSlotElementRestrictedHelperMockReactivate) Set(f func()) *SlotElementRestrictedHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReactivateFunc = f
	return m.mock
}

//Reactivate implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementRestrictedHelper interface
func (m *SlotElementRestrictedHelperMock) Reactivate() {
	counter := atomic.AddUint64(&m.ReactivatePreCounter, 1)
	defer atomic.AddUint64(&m.ReactivateCounter, 1)

	if len(m.ReactivateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReactivateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.Reactivate.")
			return
		}

		return
	}

	if m.ReactivateMock.mainExpectation != nil {

		return
	}

	if m.ReactivateFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementRestrictedHelperMock.Reactivate.")
		return
	}

	m.ReactivateFunc()
}

//ReactivateMinimockCounter returns a count of SlotElementRestrictedHelperMock.ReactivateFunc invocations
func (m *SlotElementRestrictedHelperMock) ReactivateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReactivateCounter)
}

//ReactivateMinimockPreCounter returns the value of SlotElementRestrictedHelperMock.Reactivate invocations
func (m *SlotElementRestrictedHelperMock) ReactivateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReactivatePreCounter)
}

//ReactivateFinished returns true if mock invocations count is ok
func (m *SlotElementRestrictedHelperMock) ReactivateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReactivateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReactivateCounter) == uint64(len(m.ReactivateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReactivateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReactivateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReactivateFunc != nil {
		return atomic.LoadUint64(&m.ReactivateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotElementRestrictedHelperMock) ValidateCallCounters() {

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetElementID")
	}

	if !m.GetInputEventFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetInputEvent")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetNodeID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetParentElementID")
	}

	if !m.GetPayloadFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetPayload")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetState")
	}

	if !m.GetTypeFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetType")
	}

	if !m.LeaveSequenceFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.LeaveSequence")
	}

	if !m.ReactivateFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.Reactivate")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotElementRestrictedHelperMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SlotElementRestrictedHelperMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SlotElementRestrictedHelperMock) MinimockFinish() {

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetElementID")
	}

	if !m.GetInputEventFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetInputEvent")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetNodeID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetParentElementID")
	}

	if !m.GetPayloadFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetPayload")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetState")
	}

	if !m.GetTypeFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.GetType")
	}

	if !m.LeaveSequenceFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.LeaveSequence")
	}

	if !m.ReactivateFinished() {
		m.t.Fatal("Expected call to SlotElementRestrictedHelperMock.Reactivate")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SlotElementRestrictedHelperMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SlotElementRestrictedHelperMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetElementIDFinished()
		ok = ok && m.GetInputEventFinished()
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetParentElementIDFinished()
		ok = ok && m.GetPayloadFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.GetTypeFinished()
		ok = ok && m.LeaveSequenceFinished()
		ok = ok && m.ReactivateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetElementIDFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.GetElementID")
			}

			if !m.GetInputEventFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.GetInputEvent")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.GetNodeID")
			}

			if !m.GetParentElementIDFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.GetParentElementID")
			}

			if !m.GetPayloadFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.GetPayload")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.GetState")
			}

			if !m.GetTypeFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.GetType")
			}

			if !m.LeaveSequenceFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.LeaveSequence")
			}

			if !m.ReactivateFinished() {
				m.t.Error("Expected call to SlotElementRestrictedHelperMock.Reactivate")
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
func (m *SlotElementRestrictedHelperMock) AllMocksCalled() bool {

	if !m.GetElementIDFinished() {
		return false
	}

	if !m.GetInputEventFinished() {
		return false
	}

	if !m.GetNodeIDFinished() {
		return false
	}

	if !m.GetParentElementIDFinished() {
		return false
	}

	if !m.GetPayloadFinished() {
		return false
	}

	if !m.GetStateFinished() {
		return false
	}

	if !m.GetTypeFinished() {
		return false
	}

	if !m.LeaveSequenceFinished() {
		return false
	}

	if !m.ReactivateFinished() {
		return false
	}

	return true
}

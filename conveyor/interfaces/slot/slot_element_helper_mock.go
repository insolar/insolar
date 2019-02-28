package slot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SlotElementHelper" can be found in github.com/insolar/insolar/conveyor/interfaces/slot
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//SlotElementHelperMock implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper
type SlotElementHelperMock struct {
	t minimock.Tester

	DeactivateTillFunc       func(p ReactivateMode)
	DeactivateTillCounter    uint64
	DeactivateTillPreCounter uint64
	DeactivateTillMock       mSlotElementHelperMockDeactivateTill

	GetElementIDFunc       func() (r uint32)
	GetElementIDCounter    uint64
	GetElementIDPreCounter uint64
	GetElementIDMock       mSlotElementHelperMockGetElementID

	GetInputEventFunc       func() (r interface{})
	GetInputEventCounter    uint64
	GetInputEventPreCounter uint64
	GetInputEventMock       mSlotElementHelperMockGetInputEvent

	GetNodeIDFunc       func() (r uint32)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mSlotElementHelperMockGetNodeID

	GetParentElementIDFunc       func() (r uint32)
	GetParentElementIDCounter    uint64
	GetParentElementIDPreCounter uint64
	GetParentElementIDMock       mSlotElementHelperMockGetParentElementID

	GetPayloadFunc       func() (r interface{})
	GetPayloadCounter    uint64
	GetPayloadPreCounter uint64
	GetPayloadMock       mSlotElementHelperMockGetPayload

	GetStateFunc       func() (r int)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mSlotElementHelperMockGetState

	GetTypeFunc       func() (r int)
	GetTypeCounter    uint64
	GetTypePreCounter uint64
	GetTypeMock       mSlotElementHelperMockGetType

	InformParentFunc       func(p interface{}) (r bool)
	InformParentCounter    uint64
	InformParentPreCounter uint64
	InformParentMock       mSlotElementHelperMockInformParent

	LeaveSequenceFunc       func()
	LeaveSequenceCounter    uint64
	LeaveSequencePreCounter uint64
	LeaveSequenceMock       mSlotElementHelperMockLeaveSequence

	ReactivateFunc       func()
	ReactivateCounter    uint64
	ReactivatePreCounter uint64
	ReactivateMock       mSlotElementHelperMockReactivate

	SendTaskFunc       func(p uint32, p1 interface{}, p2 uint32) (r error)
	SendTaskCounter    uint64
	SendTaskPreCounter uint64
	SendTaskMock       mSlotElementHelperMockSendTask
}

//NewSlotElementHelperMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper
func NewSlotElementHelperMock(t minimock.Tester) *SlotElementHelperMock {
	m := &SlotElementHelperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeactivateTillMock = mSlotElementHelperMockDeactivateTill{mock: m}
	m.GetElementIDMock = mSlotElementHelperMockGetElementID{mock: m}
	m.GetInputEventMock = mSlotElementHelperMockGetInputEvent{mock: m}
	m.GetNodeIDMock = mSlotElementHelperMockGetNodeID{mock: m}
	m.GetParentElementIDMock = mSlotElementHelperMockGetParentElementID{mock: m}
	m.GetPayloadMock = mSlotElementHelperMockGetPayload{mock: m}
	m.GetStateMock = mSlotElementHelperMockGetState{mock: m}
	m.GetTypeMock = mSlotElementHelperMockGetType{mock: m}
	m.InformParentMock = mSlotElementHelperMockInformParent{mock: m}
	m.LeaveSequenceMock = mSlotElementHelperMockLeaveSequence{mock: m}
	m.ReactivateMock = mSlotElementHelperMockReactivate{mock: m}
	m.SendTaskMock = mSlotElementHelperMockSendTask{mock: m}

	return m
}

type mSlotElementHelperMockDeactivateTill struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockDeactivateTillExpectation
	expectationSeries []*SlotElementHelperMockDeactivateTillExpectation
}

type SlotElementHelperMockDeactivateTillExpectation struct {
	input *SlotElementHelperMockDeactivateTillInput
}

type SlotElementHelperMockDeactivateTillInput struct {
	p ReactivateMode
}

//Expect specifies that invocation of SlotElementHelper.DeactivateTill is expected from 1 to Infinity times
func (m *mSlotElementHelperMockDeactivateTill) Expect(p ReactivateMode) *mSlotElementHelperMockDeactivateTill {
	m.mock.DeactivateTillFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockDeactivateTillExpectation{}
	}
	m.mainExpectation.input = &SlotElementHelperMockDeactivateTillInput{p}
	return m
}

//Return specifies results of invocation of SlotElementHelper.DeactivateTill
func (m *mSlotElementHelperMockDeactivateTill) Return() *SlotElementHelperMock {
	m.mock.DeactivateTillFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockDeactivateTillExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.DeactivateTill is expected once
func (m *mSlotElementHelperMockDeactivateTill) ExpectOnce(p ReactivateMode) *SlotElementHelperMockDeactivateTillExpectation {
	m.mock.DeactivateTillFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockDeactivateTillExpectation{}
	expectation.input = &SlotElementHelperMockDeactivateTillInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of SlotElementHelper.DeactivateTill method
func (m *mSlotElementHelperMockDeactivateTill) Set(f func(p ReactivateMode)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeactivateTillFunc = f
	return m.mock
}

//DeactivateTill implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) DeactivateTill(p ReactivateMode) {
	counter := atomic.AddUint64(&m.DeactivateTillPreCounter, 1)
	defer atomic.AddUint64(&m.DeactivateTillCounter, 1)

	if len(m.DeactivateTillMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeactivateTillMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.DeactivateTill. %v", p)
			return
		}

		input := m.DeactivateTillMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SlotElementHelperMockDeactivateTillInput{p}, "SlotElementHelper.DeactivateTill got unexpected parameters")

		return
	}

	if m.DeactivateTillMock.mainExpectation != nil {

		input := m.DeactivateTillMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SlotElementHelperMockDeactivateTillInput{p}, "SlotElementHelper.DeactivateTill got unexpected parameters")
		}

		return
	}

	if m.DeactivateTillFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.DeactivateTill. %v", p)
		return
	}

	m.DeactivateTillFunc(p)
}

//DeactivateTillMinimockCounter returns a count of SlotElementHelperMock.DeactivateTillFunc invocations
func (m *SlotElementHelperMock) DeactivateTillMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateTillCounter)
}

//DeactivateTillMinimockPreCounter returns the value of SlotElementHelperMock.DeactivateTill invocations
func (m *SlotElementHelperMock) DeactivateTillMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateTillPreCounter)
}

//DeactivateTillFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) DeactivateTillFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeactivateTillMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeactivateTillCounter) == uint64(len(m.DeactivateTillMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeactivateTillMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeactivateTillCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeactivateTillFunc != nil {
		return atomic.LoadUint64(&m.DeactivateTillCounter) > 0
	}

	return true
}

type mSlotElementHelperMockGetElementID struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockGetElementIDExpectation
	expectationSeries []*SlotElementHelperMockGetElementIDExpectation
}

type SlotElementHelperMockGetElementIDExpectation struct {
	result *SlotElementHelperMockGetElementIDResult
}

type SlotElementHelperMockGetElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementHelper.GetElementID is expected from 1 to Infinity times
func (m *mSlotElementHelperMockGetElementID) Expect() *mSlotElementHelperMockGetElementID {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.GetElementID
func (m *mSlotElementHelperMockGetElementID) Return(r uint32) *SlotElementHelperMock {
	m.mock.GetElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetElementIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockGetElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.GetElementID is expected once
func (m *mSlotElementHelperMockGetElementID) ExpectOnce() *SlotElementHelperMockGetElementIDExpectation {
	m.mock.GetElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockGetElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockGetElementIDExpectation) Return(r uint32) {
	e.result = &SlotElementHelperMockGetElementIDResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.GetElementID method
func (m *mSlotElementHelperMockGetElementID) Set(f func() (r uint32)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetElementIDFunc = f
	return m.mock
}

//GetElementID implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) GetElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetElementIDCounter, 1)

	if len(m.GetElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetElementID.")
			return
		}

		result := m.GetElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetElementIDMock.mainExpectation != nil {

		result := m.GetElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetElementID")
		}

		r = result.r

		return
	}

	if m.GetElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetElementID.")
		return
	}

	return m.GetElementIDFunc()
}

//GetElementIDMinimockCounter returns a count of SlotElementHelperMock.GetElementIDFunc invocations
func (m *SlotElementHelperMock) GetElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDCounter)
}

//GetElementIDMinimockPreCounter returns the value of SlotElementHelperMock.GetElementID invocations
func (m *SlotElementHelperMock) GetElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetElementIDPreCounter)
}

//GetElementIDFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) GetElementIDFinished() bool {
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

type mSlotElementHelperMockGetInputEvent struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockGetInputEventExpectation
	expectationSeries []*SlotElementHelperMockGetInputEventExpectation
}

type SlotElementHelperMockGetInputEventExpectation struct {
	result *SlotElementHelperMockGetInputEventResult
}

type SlotElementHelperMockGetInputEventResult struct {
	r interface{}
}

//Expect specifies that invocation of SlotElementHelper.GetInputEvent is expected from 1 to Infinity times
func (m *mSlotElementHelperMockGetInputEvent) Expect() *mSlotElementHelperMockGetInputEvent {
	m.mock.GetInputEventFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetInputEventExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.GetInputEvent
func (m *mSlotElementHelperMockGetInputEvent) Return(r interface{}) *SlotElementHelperMock {
	m.mock.GetInputEventFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetInputEventExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockGetInputEventResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.GetInputEvent is expected once
func (m *mSlotElementHelperMockGetInputEvent) ExpectOnce() *SlotElementHelperMockGetInputEventExpectation {
	m.mock.GetInputEventFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockGetInputEventExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockGetInputEventExpectation) Return(r interface{}) {
	e.result = &SlotElementHelperMockGetInputEventResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.GetInputEvent method
func (m *mSlotElementHelperMockGetInputEvent) Set(f func() (r interface{})) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetInputEventFunc = f
	return m.mock
}

//GetInputEvent implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) GetInputEvent() (r interface{}) {
	counter := atomic.AddUint64(&m.GetInputEventPreCounter, 1)
	defer atomic.AddUint64(&m.GetInputEventCounter, 1)

	if len(m.GetInputEventMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetInputEventMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetInputEvent.")
			return
		}

		result := m.GetInputEventMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetInputEvent")
			return
		}

		r = result.r

		return
	}

	if m.GetInputEventMock.mainExpectation != nil {

		result := m.GetInputEventMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetInputEvent")
		}

		r = result.r

		return
	}

	if m.GetInputEventFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetInputEvent.")
		return
	}

	return m.GetInputEventFunc()
}

//GetInputEventMinimockCounter returns a count of SlotElementHelperMock.GetInputEventFunc invocations
func (m *SlotElementHelperMock) GetInputEventMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetInputEventCounter)
}

//GetInputEventMinimockPreCounter returns the value of SlotElementHelperMock.GetInputEvent invocations
func (m *SlotElementHelperMock) GetInputEventMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetInputEventPreCounter)
}

//GetInputEventFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) GetInputEventFinished() bool {
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

type mSlotElementHelperMockGetNodeID struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockGetNodeIDExpectation
	expectationSeries []*SlotElementHelperMockGetNodeIDExpectation
}

type SlotElementHelperMockGetNodeIDExpectation struct {
	result *SlotElementHelperMockGetNodeIDResult
}

type SlotElementHelperMockGetNodeIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementHelper.GetNodeID is expected from 1 to Infinity times
func (m *mSlotElementHelperMockGetNodeID) Expect() *mSlotElementHelperMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.GetNodeID
func (m *mSlotElementHelperMockGetNodeID) Return(r uint32) *SlotElementHelperMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.GetNodeID is expected once
func (m *mSlotElementHelperMockGetNodeID) ExpectOnce() *SlotElementHelperMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockGetNodeIDExpectation) Return(r uint32) {
	e.result = &SlotElementHelperMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.GetNodeID method
func (m *mSlotElementHelperMockGetNodeID) Set(f func() (r uint32)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) GetNodeID() (r uint32) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of SlotElementHelperMock.GetNodeIDFunc invocations
func (m *SlotElementHelperMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of SlotElementHelperMock.GetNodeID invocations
func (m *SlotElementHelperMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) GetNodeIDFinished() bool {
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

type mSlotElementHelperMockGetParentElementID struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockGetParentElementIDExpectation
	expectationSeries []*SlotElementHelperMockGetParentElementIDExpectation
}

type SlotElementHelperMockGetParentElementIDExpectation struct {
	result *SlotElementHelperMockGetParentElementIDResult
}

type SlotElementHelperMockGetParentElementIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotElementHelper.GetParentElementID is expected from 1 to Infinity times
func (m *mSlotElementHelperMockGetParentElementID) Expect() *mSlotElementHelperMockGetParentElementID {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetParentElementIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.GetParentElementID
func (m *mSlotElementHelperMockGetParentElementID) Return(r uint32) *SlotElementHelperMock {
	m.mock.GetParentElementIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetParentElementIDExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockGetParentElementIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.GetParentElementID is expected once
func (m *mSlotElementHelperMockGetParentElementID) ExpectOnce() *SlotElementHelperMockGetParentElementIDExpectation {
	m.mock.GetParentElementIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockGetParentElementIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockGetParentElementIDExpectation) Return(r uint32) {
	e.result = &SlotElementHelperMockGetParentElementIDResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.GetParentElementID method
func (m *mSlotElementHelperMockGetParentElementID) Set(f func() (r uint32)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetParentElementIDFunc = f
	return m.mock
}

//GetParentElementID implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) GetParentElementID() (r uint32) {
	counter := atomic.AddUint64(&m.GetParentElementIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetParentElementIDCounter, 1)

	if len(m.GetParentElementIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetParentElementIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetParentElementID.")
			return
		}

		result := m.GetParentElementIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetParentElementID")
			return
		}

		r = result.r

		return
	}

	if m.GetParentElementIDMock.mainExpectation != nil {

		result := m.GetParentElementIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetParentElementID")
		}

		r = result.r

		return
	}

	if m.GetParentElementIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetParentElementID.")
		return
	}

	return m.GetParentElementIDFunc()
}

//GetParentElementIDMinimockCounter returns a count of SlotElementHelperMock.GetParentElementIDFunc invocations
func (m *SlotElementHelperMock) GetParentElementIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDCounter)
}

//GetParentElementIDMinimockPreCounter returns the value of SlotElementHelperMock.GetParentElementID invocations
func (m *SlotElementHelperMock) GetParentElementIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetParentElementIDPreCounter)
}

//GetParentElementIDFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) GetParentElementIDFinished() bool {
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

type mSlotElementHelperMockGetPayload struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockGetPayloadExpectation
	expectationSeries []*SlotElementHelperMockGetPayloadExpectation
}

type SlotElementHelperMockGetPayloadExpectation struct {
	result *SlotElementHelperMockGetPayloadResult
}

type SlotElementHelperMockGetPayloadResult struct {
	r interface{}
}

//Expect specifies that invocation of SlotElementHelper.GetPayload is expected from 1 to Infinity times
func (m *mSlotElementHelperMockGetPayload) Expect() *mSlotElementHelperMockGetPayload {
	m.mock.GetPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetPayloadExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.GetPayload
func (m *mSlotElementHelperMockGetPayload) Return(r interface{}) *SlotElementHelperMock {
	m.mock.GetPayloadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetPayloadExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockGetPayloadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.GetPayload is expected once
func (m *mSlotElementHelperMockGetPayload) ExpectOnce() *SlotElementHelperMockGetPayloadExpectation {
	m.mock.GetPayloadFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockGetPayloadExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockGetPayloadExpectation) Return(r interface{}) {
	e.result = &SlotElementHelperMockGetPayloadResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.GetPayload method
func (m *mSlotElementHelperMockGetPayload) Set(f func() (r interface{})) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPayloadFunc = f
	return m.mock
}

//GetPayload implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) GetPayload() (r interface{}) {
	counter := atomic.AddUint64(&m.GetPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.GetPayloadCounter, 1)

	if len(m.GetPayloadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPayloadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetPayload.")
			return
		}

		result := m.GetPayloadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetPayload")
			return
		}

		r = result.r

		return
	}

	if m.GetPayloadMock.mainExpectation != nil {

		result := m.GetPayloadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetPayload")
		}

		r = result.r

		return
	}

	if m.GetPayloadFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetPayload.")
		return
	}

	return m.GetPayloadFunc()
}

//GetPayloadMinimockCounter returns a count of SlotElementHelperMock.GetPayloadFunc invocations
func (m *SlotElementHelperMock) GetPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPayloadCounter)
}

//GetPayloadMinimockPreCounter returns the value of SlotElementHelperMock.GetPayload invocations
func (m *SlotElementHelperMock) GetPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPayloadPreCounter)
}

//GetPayloadFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) GetPayloadFinished() bool {
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

type mSlotElementHelperMockGetState struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockGetStateExpectation
	expectationSeries []*SlotElementHelperMockGetStateExpectation
}

type SlotElementHelperMockGetStateExpectation struct {
	result *SlotElementHelperMockGetStateResult
}

type SlotElementHelperMockGetStateResult struct {
	r int
}

//Expect specifies that invocation of SlotElementHelper.GetState is expected from 1 to Infinity times
func (m *mSlotElementHelperMockGetState) Expect() *mSlotElementHelperMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.GetState
func (m *mSlotElementHelperMockGetState) Return(r int) *SlotElementHelperMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetStateExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.GetState is expected once
func (m *mSlotElementHelperMockGetState) ExpectOnce() *SlotElementHelperMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockGetStateExpectation) Return(r int) {
	e.result = &SlotElementHelperMockGetStateResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.GetState method
func (m *mSlotElementHelperMockGetState) Set(f func() (r int)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) GetState() (r int) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of SlotElementHelperMock.GetStateFunc invocations
func (m *SlotElementHelperMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of SlotElementHelperMock.GetState invocations
func (m *SlotElementHelperMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) GetStateFinished() bool {
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

type mSlotElementHelperMockGetType struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockGetTypeExpectation
	expectationSeries []*SlotElementHelperMockGetTypeExpectation
}

type SlotElementHelperMockGetTypeExpectation struct {
	result *SlotElementHelperMockGetTypeResult
}

type SlotElementHelperMockGetTypeResult struct {
	r int
}

//Expect specifies that invocation of SlotElementHelper.GetType is expected from 1 to Infinity times
func (m *mSlotElementHelperMockGetType) Expect() *mSlotElementHelperMockGetType {
	m.mock.GetTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.GetType
func (m *mSlotElementHelperMockGetType) Return(r int) *SlotElementHelperMock {
	m.mock.GetTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockGetTypeExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockGetTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.GetType is expected once
func (m *mSlotElementHelperMockGetType) ExpectOnce() *SlotElementHelperMockGetTypeExpectation {
	m.mock.GetTypeFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockGetTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockGetTypeExpectation) Return(r int) {
	e.result = &SlotElementHelperMockGetTypeResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.GetType method
func (m *mSlotElementHelperMockGetType) Set(f func() (r int)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTypeFunc = f
	return m.mock
}

//GetType implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) GetType() (r int) {
	counter := atomic.AddUint64(&m.GetTypePreCounter, 1)
	defer atomic.AddUint64(&m.GetTypeCounter, 1)

	if len(m.GetTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetType.")
			return
		}

		result := m.GetTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetType")
			return
		}

		r = result.r

		return
	}

	if m.GetTypeMock.mainExpectation != nil {

		result := m.GetTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.GetType")
		}

		r = result.r

		return
	}

	if m.GetTypeFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.GetType.")
		return
	}

	return m.GetTypeFunc()
}

//GetTypeMinimockCounter returns a count of SlotElementHelperMock.GetTypeFunc invocations
func (m *SlotElementHelperMock) GetTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypeCounter)
}

//GetTypeMinimockPreCounter returns the value of SlotElementHelperMock.GetType invocations
func (m *SlotElementHelperMock) GetTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypePreCounter)
}

//GetTypeFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) GetTypeFinished() bool {
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

type mSlotElementHelperMockInformParent struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockInformParentExpectation
	expectationSeries []*SlotElementHelperMockInformParentExpectation
}

type SlotElementHelperMockInformParentExpectation struct {
	input  *SlotElementHelperMockInformParentInput
	result *SlotElementHelperMockInformParentResult
}

type SlotElementHelperMockInformParentInput struct {
	p interface{}
}

type SlotElementHelperMockInformParentResult struct {
	r bool
}

//Expect specifies that invocation of SlotElementHelper.InformParent is expected from 1 to Infinity times
func (m *mSlotElementHelperMockInformParent) Expect(p interface{}) *mSlotElementHelperMockInformParent {
	m.mock.InformParentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockInformParentExpectation{}
	}
	m.mainExpectation.input = &SlotElementHelperMockInformParentInput{p}
	return m
}

//Return specifies results of invocation of SlotElementHelper.InformParent
func (m *mSlotElementHelperMockInformParent) Return(r bool) *SlotElementHelperMock {
	m.mock.InformParentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockInformParentExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockInformParentResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.InformParent is expected once
func (m *mSlotElementHelperMockInformParent) ExpectOnce(p interface{}) *SlotElementHelperMockInformParentExpectation {
	m.mock.InformParentFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockInformParentExpectation{}
	expectation.input = &SlotElementHelperMockInformParentInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockInformParentExpectation) Return(r bool) {
	e.result = &SlotElementHelperMockInformParentResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.InformParent method
func (m *mSlotElementHelperMockInformParent) Set(f func(p interface{}) (r bool)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InformParentFunc = f
	return m.mock
}

//InformParent implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) InformParent(p interface{}) (r bool) {
	counter := atomic.AddUint64(&m.InformParentPreCounter, 1)
	defer atomic.AddUint64(&m.InformParentCounter, 1)

	if len(m.InformParentMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InformParentMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.InformParent. %v", p)
			return
		}

		input := m.InformParentMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SlotElementHelperMockInformParentInput{p}, "SlotElementHelper.InformParent got unexpected parameters")

		result := m.InformParentMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.InformParent")
			return
		}

		r = result.r

		return
	}

	if m.InformParentMock.mainExpectation != nil {

		input := m.InformParentMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SlotElementHelperMockInformParentInput{p}, "SlotElementHelper.InformParent got unexpected parameters")
		}

		result := m.InformParentMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.InformParent")
		}

		r = result.r

		return
	}

	if m.InformParentFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.InformParent. %v", p)
		return
	}

	return m.InformParentFunc(p)
}

//InformParentMinimockCounter returns a count of SlotElementHelperMock.InformParentFunc invocations
func (m *SlotElementHelperMock) InformParentMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InformParentCounter)
}

//InformParentMinimockPreCounter returns the value of SlotElementHelperMock.InformParent invocations
func (m *SlotElementHelperMock) InformParentMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InformParentPreCounter)
}

//InformParentFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) InformParentFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InformParentMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InformParentCounter) == uint64(len(m.InformParentMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InformParentMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InformParentCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InformParentFunc != nil {
		return atomic.LoadUint64(&m.InformParentCounter) > 0
	}

	return true
}

type mSlotElementHelperMockLeaveSequence struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockLeaveSequenceExpectation
	expectationSeries []*SlotElementHelperMockLeaveSequenceExpectation
}

type SlotElementHelperMockLeaveSequenceExpectation struct {
}

//Expect specifies that invocation of SlotElementHelper.LeaveSequence is expected from 1 to Infinity times
func (m *mSlotElementHelperMockLeaveSequence) Expect() *mSlotElementHelperMockLeaveSequence {
	m.mock.LeaveSequenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockLeaveSequenceExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.LeaveSequence
func (m *mSlotElementHelperMockLeaveSequence) Return() *SlotElementHelperMock {
	m.mock.LeaveSequenceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockLeaveSequenceExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.LeaveSequence is expected once
func (m *mSlotElementHelperMockLeaveSequence) ExpectOnce() *SlotElementHelperMockLeaveSequenceExpectation {
	m.mock.LeaveSequenceFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockLeaveSequenceExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of SlotElementHelper.LeaveSequence method
func (m *mSlotElementHelperMockLeaveSequence) Set(f func()) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LeaveSequenceFunc = f
	return m.mock
}

//LeaveSequence implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) LeaveSequence() {
	counter := atomic.AddUint64(&m.LeaveSequencePreCounter, 1)
	defer atomic.AddUint64(&m.LeaveSequenceCounter, 1)

	if len(m.LeaveSequenceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LeaveSequenceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.LeaveSequence.")
			return
		}

		return
	}

	if m.LeaveSequenceMock.mainExpectation != nil {

		return
	}

	if m.LeaveSequenceFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.LeaveSequence.")
		return
	}

	m.LeaveSequenceFunc()
}

//LeaveSequenceMinimockCounter returns a count of SlotElementHelperMock.LeaveSequenceFunc invocations
func (m *SlotElementHelperMock) LeaveSequenceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LeaveSequenceCounter)
}

//LeaveSequenceMinimockPreCounter returns the value of SlotElementHelperMock.LeaveSequence invocations
func (m *SlotElementHelperMock) LeaveSequenceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LeaveSequencePreCounter)
}

//LeaveSequenceFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) LeaveSequenceFinished() bool {
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

type mSlotElementHelperMockReactivate struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockReactivateExpectation
	expectationSeries []*SlotElementHelperMockReactivateExpectation
}

type SlotElementHelperMockReactivateExpectation struct {
}

//Expect specifies that invocation of SlotElementHelper.Reactivate is expected from 1 to Infinity times
func (m *mSlotElementHelperMockReactivate) Expect() *mSlotElementHelperMockReactivate {
	m.mock.ReactivateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockReactivateExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotElementHelper.Reactivate
func (m *mSlotElementHelperMockReactivate) Return() *SlotElementHelperMock {
	m.mock.ReactivateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockReactivateExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.Reactivate is expected once
func (m *mSlotElementHelperMockReactivate) ExpectOnce() *SlotElementHelperMockReactivateExpectation {
	m.mock.ReactivateFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockReactivateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of SlotElementHelper.Reactivate method
func (m *mSlotElementHelperMockReactivate) Set(f func()) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReactivateFunc = f
	return m.mock
}

//Reactivate implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) Reactivate() {
	counter := atomic.AddUint64(&m.ReactivatePreCounter, 1)
	defer atomic.AddUint64(&m.ReactivateCounter, 1)

	if len(m.ReactivateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReactivateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.Reactivate.")
			return
		}

		return
	}

	if m.ReactivateMock.mainExpectation != nil {

		return
	}

	if m.ReactivateFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.Reactivate.")
		return
	}

	m.ReactivateFunc()
}

//ReactivateMinimockCounter returns a count of SlotElementHelperMock.ReactivateFunc invocations
func (m *SlotElementHelperMock) ReactivateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReactivateCounter)
}

//ReactivateMinimockPreCounter returns the value of SlotElementHelperMock.Reactivate invocations
func (m *SlotElementHelperMock) ReactivateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReactivatePreCounter)
}

//ReactivateFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) ReactivateFinished() bool {
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

type mSlotElementHelperMockSendTask struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockSendTaskExpectation
	expectationSeries []*SlotElementHelperMockSendTaskExpectation
}

type SlotElementHelperMockSendTaskExpectation struct {
	input  *SlotElementHelperMockSendTaskInput
	result *SlotElementHelperMockSendTaskResult
}

type SlotElementHelperMockSendTaskInput struct {
	p  uint32
	p1 interface{}
	p2 uint32
}

type SlotElementHelperMockSendTaskResult struct {
	r error
}

//Expect specifies that invocation of SlotElementHelper.SendTask is expected from 1 to Infinity times
func (m *mSlotElementHelperMockSendTask) Expect(p uint32, p1 interface{}, p2 uint32) *mSlotElementHelperMockSendTask {
	m.mock.SendTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockSendTaskExpectation{}
	}
	m.mainExpectation.input = &SlotElementHelperMockSendTaskInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of SlotElementHelper.SendTask
func (m *mSlotElementHelperMockSendTask) Return(r error) *SlotElementHelperMock {
	m.mock.SendTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockSendTaskExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockSendTaskResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.SendTask is expected once
func (m *mSlotElementHelperMockSendTask) ExpectOnce(p uint32, p1 interface{}, p2 uint32) *SlotElementHelperMockSendTaskExpectation {
	m.mock.SendTaskFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockSendTaskExpectation{}
	expectation.input = &SlotElementHelperMockSendTaskInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockSendTaskExpectation) Return(r error) {
	e.result = &SlotElementHelperMockSendTaskResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.SendTask method
func (m *mSlotElementHelperMockSendTask) Set(f func(p uint32, p1 interface{}, p2 uint32) (r error)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendTaskFunc = f
	return m.mock
}

//SendTask implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotElementHelper interface
func (m *SlotElementHelperMock) SendTask(p uint32, p1 interface{}, p2 uint32) (r error) {
	counter := atomic.AddUint64(&m.SendTaskPreCounter, 1)
	defer atomic.AddUint64(&m.SendTaskCounter, 1)

	if len(m.SendTaskMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendTaskMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.SendTask. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendTaskMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SlotElementHelperMockSendTaskInput{p, p1, p2}, "SlotElementHelper.SendTask got unexpected parameters")

		result := m.SendTaskMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.SendTask")
			return
		}

		r = result.r

		return
	}

	if m.SendTaskMock.mainExpectation != nil {

		input := m.SendTaskMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SlotElementHelperMockSendTaskInput{p, p1, p2}, "SlotElementHelper.SendTask got unexpected parameters")
		}

		result := m.SendTaskMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.SendTask")
		}

		r = result.r

		return
	}

	if m.SendTaskFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.SendTask. %v %v %v", p, p1, p2)
		return
	}

	return m.SendTaskFunc(p, p1, p2)
}

//SendTaskMinimockCounter returns a count of SlotElementHelperMock.SendTaskFunc invocations
func (m *SlotElementHelperMock) SendTaskMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendTaskCounter)
}

//SendTaskMinimockPreCounter returns the value of SlotElementHelperMock.SendTask invocations
func (m *SlotElementHelperMock) SendTaskMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendTaskPreCounter)
}

//SendTaskFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) SendTaskFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendTaskMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendTaskCounter) == uint64(len(m.SendTaskMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendTaskMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendTaskCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendTaskFunc != nil {
		return atomic.LoadUint64(&m.SendTaskCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotElementHelperMock) ValidateCallCounters() {

	if !m.DeactivateTillFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.DeactivateTill")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetElementID")
	}

	if !m.GetInputEventFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetInputEvent")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetNodeID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetParentElementID")
	}

	if !m.GetPayloadFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetPayload")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetState")
	}

	if !m.GetTypeFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetType")
	}

	if !m.InformParentFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.InformParent")
	}

	if !m.LeaveSequenceFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.LeaveSequence")
	}

	if !m.ReactivateFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.Reactivate")
	}

	if !m.SendTaskFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.SendTask")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotElementHelperMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SlotElementHelperMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SlotElementHelperMock) MinimockFinish() {

	if !m.DeactivateTillFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.DeactivateTill")
	}

	if !m.GetElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetElementID")
	}

	if !m.GetInputEventFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetInputEvent")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetNodeID")
	}

	if !m.GetParentElementIDFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetParentElementID")
	}

	if !m.GetPayloadFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetPayload")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetState")
	}

	if !m.GetTypeFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.GetType")
	}

	if !m.InformParentFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.InformParent")
	}

	if !m.LeaveSequenceFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.LeaveSequence")
	}

	if !m.ReactivateFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.Reactivate")
	}

	if !m.SendTaskFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.SendTask")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SlotElementHelperMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SlotElementHelperMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.DeactivateTillFinished()
		ok = ok && m.GetElementIDFinished()
		ok = ok && m.GetInputEventFinished()
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetParentElementIDFinished()
		ok = ok && m.GetPayloadFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.GetTypeFinished()
		ok = ok && m.InformParentFinished()
		ok = ok && m.LeaveSequenceFinished()
		ok = ok && m.ReactivateFinished()
		ok = ok && m.SendTaskFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeactivateTillFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.DeactivateTill")
			}

			if !m.GetElementIDFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.GetElementID")
			}

			if !m.GetInputEventFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.GetInputEvent")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.GetNodeID")
			}

			if !m.GetParentElementIDFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.GetParentElementID")
			}

			if !m.GetPayloadFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.GetPayload")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.GetState")
			}

			if !m.GetTypeFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.GetType")
			}

			if !m.InformParentFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.InformParent")
			}

			if !m.LeaveSequenceFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.LeaveSequence")
			}

			if !m.ReactivateFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.Reactivate")
			}

			if !m.SendTaskFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.SendTask")
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
func (m *SlotElementHelperMock) AllMocksCalled() bool {

	if !m.DeactivateTillFinished() {
		return false
	}

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

	if !m.InformParentFinished() {
		return false
	}

	if !m.LeaveSequenceFinished() {
		return false
	}

	if !m.ReactivateFinished() {
		return false
	}

	if !m.SendTaskFinished() {
		return false
	}

	return true
}

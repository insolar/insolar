/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package slot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SlotElementHelper" can be found in github.com/insolar/insolar/conveyor/interfaces
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//SlotElementHelperMock implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper
type SlotElementHelperMock struct {
	t minimock.Tester

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

	LeaveSequenceFunc       func()
	LeaveSequenceCounter    uint64
	LeaveSequencePreCounter uint64
	LeaveSequenceMock       mSlotElementHelperMockLeaveSequence

	ReactivateFunc       func()
	ReactivateCounter    uint64
	ReactivatePreCounter uint64
	ReactivateMock       mSlotElementHelperMockReactivate

	deactivateTillFunc       func(p reactivateMode)
	deactivateTillCounter    uint64
	deactivateTillPreCounter uint64
	deactivateTillMock       mSlotElementHelperMockdeactivateTill

	informParentFunc       func(p interface{}) (r bool)
	informParentCounter    uint64
	informParentPreCounter uint64
	informParentMock       mSlotElementHelperMockinformParent

	sendTaskFunc       func(p uint32, p1 interface{}, p2 uint32) (r error)
	sendTaskCounter    uint64
	sendTaskPreCounter uint64
	sendTaskMock       mSlotElementHelperMocksendTask
}

//NewSlotElementHelperMock returns a mock for github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper
func NewSlotElementHelperMock(t minimock.Tester) *SlotElementHelperMock {
	m := &SlotElementHelperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetElementIDMock = mSlotElementHelperMockGetElementID{mock: m}
	m.GetInputEventMock = mSlotElementHelperMockGetInputEvent{mock: m}
	m.GetNodeIDMock = mSlotElementHelperMockGetNodeID{mock: m}
	m.GetParentElementIDMock = mSlotElementHelperMockGetParentElementID{mock: m}
	m.GetPayloadMock = mSlotElementHelperMockGetPayload{mock: m}
	m.GetStateMock = mSlotElementHelperMockGetState{mock: m}
	m.GetTypeMock = mSlotElementHelperMockGetType{mock: m}
	m.LeaveSequenceMock = mSlotElementHelperMockLeaveSequence{mock: m}
	m.ReactivateMock = mSlotElementHelperMockReactivate{mock: m}
	m.deactivateTillMock = mSlotElementHelperMockdeactivateTill{mock: m}
	m.informParentMock = mSlotElementHelperMockinformParent{mock: m}
	m.sendTaskMock = mSlotElementHelperMocksendTask{mock: m}

	return m
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

//GetElementID implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//GetInputEvent implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//GetNodeID implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//GetParentElementID implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//GetPayload implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//GetState implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//GetType implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//LeaveSequence implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

//Reactivate implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
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

type mSlotElementHelperMockdeactivateTill struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockdeactivateTillExpectation
	expectationSeries []*SlotElementHelperMockdeactivateTillExpectation
}

type SlotElementHelperMockdeactivateTillExpectation struct {
	input *SlotElementHelperMockdeactivateTillInput
}

type SlotElementHelperMockdeactivateTillInput struct {
	p reactivateMode
}

//Expect specifies that invocation of SlotElementHelper.deactivateTill is expected from 1 to Infinity times
func (m *mSlotElementHelperMockdeactivateTill) Expect(p reactivateMode) *mSlotElementHelperMockdeactivateTill {
	m.mock.deactivateTillFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockdeactivateTillExpectation{}
	}
	m.mainExpectation.input = &SlotElementHelperMockdeactivateTillInput{p}
	return m
}

//Return specifies results of invocation of SlotElementHelper.deactivateTill
func (m *mSlotElementHelperMockdeactivateTill) Return() *SlotElementHelperMock {
	m.mock.deactivateTillFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockdeactivateTillExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.deactivateTill is expected once
func (m *mSlotElementHelperMockdeactivateTill) ExpectOnce(p reactivateMode) *SlotElementHelperMockdeactivateTillExpectation {
	m.mock.deactivateTillFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockdeactivateTillExpectation{}
	expectation.input = &SlotElementHelperMockdeactivateTillInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of SlotElementHelper.deactivateTill method
func (m *mSlotElementHelperMockdeactivateTill) Set(f func(p reactivateMode)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.deactivateTillFunc = f
	return m.mock
}

//deactivateTill implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
func (m *SlotElementHelperMock) deactivateTill(p reactivateMode) {
	counter := atomic.AddUint64(&m.deactivateTillPreCounter, 1)
	defer atomic.AddUint64(&m.deactivateTillCounter, 1)

	if len(m.deactivateTillMock.expectationSeries) > 0 {
		if counter > uint64(len(m.deactivateTillMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.deactivateTill. %v", p)
			return
		}

		input := m.deactivateTillMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SlotElementHelperMockdeactivateTillInput{p}, "SlotElementHelper.deactivateTill got unexpected parameters")

		return
	}

	if m.deactivateTillMock.mainExpectation != nil {

		input := m.deactivateTillMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SlotElementHelperMockdeactivateTillInput{p}, "SlotElementHelper.deactivateTill got unexpected parameters")
		}

		return
	}

	if m.deactivateTillFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.deactivateTill. %v", p)
		return
	}

	m.deactivateTillFunc(p)
}

//deactivateTillMinimockCounter returns a count of SlotElementHelperMock.deactivateTillFunc invocations
func (m *SlotElementHelperMock) deactivateTillMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.deactivateTillCounter)
}

//deactivateTillMinimockPreCounter returns the value of SlotElementHelperMock.deactivateTill invocations
func (m *SlotElementHelperMock) deactivateTillMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.deactivateTillPreCounter)
}

//deactivateTillFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) deactivateTillFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.deactivateTillMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.deactivateTillCounter) == uint64(len(m.deactivateTillMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.deactivateTillMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.deactivateTillCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.deactivateTillFunc != nil {
		return atomic.LoadUint64(&m.deactivateTillCounter) > 0
	}

	return true
}

type mSlotElementHelperMockinformParent struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMockinformParentExpectation
	expectationSeries []*SlotElementHelperMockinformParentExpectation
}

type SlotElementHelperMockinformParentExpectation struct {
	input  *SlotElementHelperMockinformParentInput
	result *SlotElementHelperMockinformParentResult
}

type SlotElementHelperMockinformParentInput struct {
	p interface{}
}

type SlotElementHelperMockinformParentResult struct {
	r bool
}

//Expect specifies that invocation of SlotElementHelper.informParent is expected from 1 to Infinity times
func (m *mSlotElementHelperMockinformParent) Expect(p interface{}) *mSlotElementHelperMockinformParent {
	m.mock.informParentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockinformParentExpectation{}
	}
	m.mainExpectation.input = &SlotElementHelperMockinformParentInput{p}
	return m
}

//Return specifies results of invocation of SlotElementHelper.informParent
func (m *mSlotElementHelperMockinformParent) Return(r bool) *SlotElementHelperMock {
	m.mock.informParentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMockinformParentExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMockinformParentResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.informParent is expected once
func (m *mSlotElementHelperMockinformParent) ExpectOnce(p interface{}) *SlotElementHelperMockinformParentExpectation {
	m.mock.informParentFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMockinformParentExpectation{}
	expectation.input = &SlotElementHelperMockinformParentInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMockinformParentExpectation) Return(r bool) {
	e.result = &SlotElementHelperMockinformParentResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.informParent method
func (m *mSlotElementHelperMockinformParent) Set(f func(p interface{}) (r bool)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.informParentFunc = f
	return m.mock
}

//informParent implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
func (m *SlotElementHelperMock) informParent(p interface{}) (r bool) {
	counter := atomic.AddUint64(&m.informParentPreCounter, 1)
	defer atomic.AddUint64(&m.informParentCounter, 1)

	if len(m.informParentMock.expectationSeries) > 0 {
		if counter > uint64(len(m.informParentMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.informParent. %v", p)
			return
		}

		input := m.informParentMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SlotElementHelperMockinformParentInput{p}, "SlotElementHelper.informParent got unexpected parameters")

		result := m.informParentMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.informParent")
			return
		}

		r = result.r

		return
	}

	if m.informParentMock.mainExpectation != nil {

		input := m.informParentMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SlotElementHelperMockinformParentInput{p}, "SlotElementHelper.informParent got unexpected parameters")
		}

		result := m.informParentMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.informParent")
		}

		r = result.r

		return
	}

	if m.informParentFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.informParent. %v", p)
		return
	}

	return m.informParentFunc(p)
}

//informParentMinimockCounter returns a count of SlotElementHelperMock.informParentFunc invocations
func (m *SlotElementHelperMock) informParentMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.informParentCounter)
}

//informParentMinimockPreCounter returns the value of SlotElementHelperMock.informParent invocations
func (m *SlotElementHelperMock) informParentMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.informParentPreCounter)
}

//informParentFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) informParentFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.informParentMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.informParentCounter) == uint64(len(m.informParentMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.informParentMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.informParentCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.informParentFunc != nil {
		return atomic.LoadUint64(&m.informParentCounter) > 0
	}

	return true
}

type mSlotElementHelperMocksendTask struct {
	mock              *SlotElementHelperMock
	mainExpectation   *SlotElementHelperMocksendTaskExpectation
	expectationSeries []*SlotElementHelperMocksendTaskExpectation
}

type SlotElementHelperMocksendTaskExpectation struct {
	input  *SlotElementHelperMocksendTaskInput
	result *SlotElementHelperMocksendTaskResult
}

type SlotElementHelperMocksendTaskInput struct {
	p  uint32
	p1 interface{}
	p2 uint32
}

type SlotElementHelperMocksendTaskResult struct {
	r error
}

//Expect specifies that invocation of SlotElementHelper.sendTask is expected from 1 to Infinity times
func (m *mSlotElementHelperMocksendTask) Expect(p uint32, p1 interface{}, p2 uint32) *mSlotElementHelperMocksendTask {
	m.mock.sendTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMocksendTaskExpectation{}
	}
	m.mainExpectation.input = &SlotElementHelperMocksendTaskInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of SlotElementHelper.sendTask
func (m *mSlotElementHelperMocksendTask) Return(r error) *SlotElementHelperMock {
	m.mock.sendTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotElementHelperMocksendTaskExpectation{}
	}
	m.mainExpectation.result = &SlotElementHelperMocksendTaskResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotElementHelper.sendTask is expected once
func (m *mSlotElementHelperMocksendTask) ExpectOnce(p uint32, p1 interface{}, p2 uint32) *SlotElementHelperMocksendTaskExpectation {
	m.mock.sendTaskFunc = nil
	m.mainExpectation = nil

	expectation := &SlotElementHelperMocksendTaskExpectation{}
	expectation.input = &SlotElementHelperMocksendTaskInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotElementHelperMocksendTaskExpectation) Return(r error) {
	e.result = &SlotElementHelperMocksendTaskResult{r}
}

//Set uses given function f as a mock of SlotElementHelper.sendTask method
func (m *mSlotElementHelperMocksendTask) Set(f func(p uint32, p1 interface{}, p2 uint32) (r error)) *SlotElementHelperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.sendTaskFunc = f
	return m.mock
}

//sendTask implements github.com/insolar/insolar/conveyor/interfaces.SlotElementHelper interface
func (m *SlotElementHelperMock) sendTask(p uint32, p1 interface{}, p2 uint32) (r error) {
	counter := atomic.AddUint64(&m.sendTaskPreCounter, 1)
	defer atomic.AddUint64(&m.sendTaskCounter, 1)

	if len(m.sendTaskMock.expectationSeries) > 0 {
		if counter > uint64(len(m.sendTaskMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotElementHelperMock.sendTask. %v %v %v", p, p1, p2)
			return
		}

		input := m.sendTaskMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SlotElementHelperMocksendTaskInput{p, p1, p2}, "SlotElementHelper.sendTask got unexpected parameters")

		result := m.sendTaskMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.sendTask")
			return
		}

		r = result.r

		return
	}

	if m.sendTaskMock.mainExpectation != nil {

		input := m.sendTaskMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SlotElementHelperMocksendTaskInput{p, p1, p2}, "SlotElementHelper.sendTask got unexpected parameters")
		}

		result := m.sendTaskMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotElementHelperMock.sendTask")
		}

		r = result.r

		return
	}

	if m.sendTaskFunc == nil {
		m.t.Fatalf("Unexpected call to SlotElementHelperMock.sendTask. %v %v %v", p, p1, p2)
		return
	}

	return m.sendTaskFunc(p, p1, p2)
}

//sendTaskMinimockCounter returns a count of SlotElementHelperMock.sendTaskFunc invocations
func (m *SlotElementHelperMock) sendTaskMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.sendTaskCounter)
}

//sendTaskMinimockPreCounter returns the value of SlotElementHelperMock.sendTask invocations
func (m *SlotElementHelperMock) sendTaskMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.sendTaskPreCounter)
}

//sendTaskFinished returns true if mock invocations count is ok
func (m *SlotElementHelperMock) sendTaskFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.sendTaskMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.sendTaskCounter) == uint64(len(m.sendTaskMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.sendTaskMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.sendTaskCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.sendTaskFunc != nil {
		return atomic.LoadUint64(&m.sendTaskCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotElementHelperMock) ValidateCallCounters() {

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

	if !m.LeaveSequenceFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.LeaveSequence")
	}

	if !m.ReactivateFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.Reactivate")
	}

	if !m.deactivateTillFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.deactivateTill")
	}

	if !m.informParentFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.informParent")
	}

	if !m.sendTaskFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.sendTask")
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

	if !m.LeaveSequenceFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.LeaveSequence")
	}

	if !m.ReactivateFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.Reactivate")
	}

	if !m.deactivateTillFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.deactivateTill")
	}

	if !m.informParentFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.informParent")
	}

	if !m.sendTaskFinished() {
		m.t.Fatal("Expected call to SlotElementHelperMock.sendTask")
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
		ok = ok && m.GetElementIDFinished()
		ok = ok && m.GetInputEventFinished()
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetParentElementIDFinished()
		ok = ok && m.GetPayloadFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.GetTypeFinished()
		ok = ok && m.LeaveSequenceFinished()
		ok = ok && m.ReactivateFinished()
		ok = ok && m.deactivateTillFinished()
		ok = ok && m.informParentFinished()
		ok = ok && m.sendTaskFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

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

			if !m.LeaveSequenceFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.LeaveSequence")
			}

			if !m.ReactivateFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.Reactivate")
			}

			if !m.deactivateTillFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.deactivateTill")
			}

			if !m.informParentFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.informParent")
			}

			if !m.sendTaskFinished() {
				m.t.Error("Expected call to SlotElementHelperMock.sendTask")
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

	if !m.deactivateTillFinished() {
		return false
	}

	if !m.informParentFinished() {
		return false
	}

	if !m.sendTaskFinished() {
		return false
	}

	return true
}

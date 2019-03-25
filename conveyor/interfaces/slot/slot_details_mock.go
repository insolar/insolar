package slot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SlotDetails" can be found in github.com/insolar/insolar/conveyor/interfaces/slot
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
)

//SlotDetailsMock implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotDetails
type SlotDetailsMock struct {
	t minimock.Tester

	GetNodeDataFunc       func() (r interface{})
	GetNodeDataCounter    uint64
	GetNodeDataPreCounter uint64
	GetNodeDataMock       mSlotDetailsMockGetNodeData

	GetNodeIDFunc       func() (r uint32)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mSlotDetailsMockGetNodeID

	GetPulseDataFunc       func() (r insolar.Pulse)
	GetPulseDataCounter    uint64
	GetPulseDataPreCounter uint64
	GetPulseDataMock       mSlotDetailsMockGetPulseData

	GetPulseNumberFunc       func() (r insolar.PulseNumber)
	GetPulseNumberCounter    uint64
	GetPulseNumberPreCounter uint64
	GetPulseNumberMock       mSlotDetailsMockGetPulseNumber
}

//NewSlotDetailsMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/slot.SlotDetails
func NewSlotDetailsMock(t minimock.Tester) *SlotDetailsMock {
	m := &SlotDetailsMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetNodeDataMock = mSlotDetailsMockGetNodeData{mock: m}
	m.GetNodeIDMock = mSlotDetailsMockGetNodeID{mock: m}
	m.GetPulseDataMock = mSlotDetailsMockGetPulseData{mock: m}
	m.GetPulseNumberMock = mSlotDetailsMockGetPulseNumber{mock: m}

	return m
}

type mSlotDetailsMockGetNodeData struct {
	mock              *SlotDetailsMock
	mainExpectation   *SlotDetailsMockGetNodeDataExpectation
	expectationSeries []*SlotDetailsMockGetNodeDataExpectation
}

type SlotDetailsMockGetNodeDataExpectation struct {
	result *SlotDetailsMockGetNodeDataResult
}

type SlotDetailsMockGetNodeDataResult struct {
	r interface{}
}

//Expect specifies that invocation of SlotDetails.GetNodeData is expected from 1 to Infinity times
func (m *mSlotDetailsMockGetNodeData) Expect() *mSlotDetailsMockGetNodeData {
	m.mock.GetNodeDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetNodeDataExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotDetails.GetNodeData
func (m *mSlotDetailsMockGetNodeData) Return(r interface{}) *SlotDetailsMock {
	m.mock.GetNodeDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetNodeDataExpectation{}
	}
	m.mainExpectation.result = &SlotDetailsMockGetNodeDataResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotDetails.GetNodeData is expected once
func (m *mSlotDetailsMockGetNodeData) ExpectOnce() *SlotDetailsMockGetNodeDataExpectation {
	m.mock.GetNodeDataFunc = nil
	m.mainExpectation = nil

	expectation := &SlotDetailsMockGetNodeDataExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotDetailsMockGetNodeDataExpectation) Return(r interface{}) {
	e.result = &SlotDetailsMockGetNodeDataResult{r}
}

//Set uses given function f as a mock of SlotDetails.GetNodeData method
func (m *mSlotDetailsMockGetNodeData) Set(f func() (r interface{})) *SlotDetailsMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeDataFunc = f
	return m.mock
}

//GetNodeData implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotDetails interface
func (m *SlotDetailsMock) GetNodeData() (r interface{}) {
	counter := atomic.AddUint64(&m.GetNodeDataPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeDataCounter, 1)

	if len(m.GetNodeDataMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeDataMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotDetailsMock.GetNodeData.")
			return
		}

		result := m.GetNodeDataMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetNodeData")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeDataMock.mainExpectation != nil {

		result := m.GetNodeDataMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetNodeData")
		}

		r = result.r

		return
	}

	if m.GetNodeDataFunc == nil {
		m.t.Fatalf("Unexpected call to SlotDetailsMock.GetNodeData.")
		return
	}

	return m.GetNodeDataFunc()
}

//GetNodeDataMinimockCounter returns a count of SlotDetailsMock.GetNodeDataFunc invocations
func (m *SlotDetailsMock) GetNodeDataMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeDataCounter)
}

//GetNodeDataMinimockPreCounter returns the value of SlotDetailsMock.GetNodeData invocations
func (m *SlotDetailsMock) GetNodeDataMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeDataPreCounter)
}

//GetNodeDataFinished returns true if mock invocations count is ok
func (m *SlotDetailsMock) GetNodeDataFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeDataMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeDataCounter) == uint64(len(m.GetNodeDataMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeDataMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeDataCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeDataFunc != nil {
		return atomic.LoadUint64(&m.GetNodeDataCounter) > 0
	}

	return true
}

type mSlotDetailsMockGetNodeID struct {
	mock              *SlotDetailsMock
	mainExpectation   *SlotDetailsMockGetNodeIDExpectation
	expectationSeries []*SlotDetailsMockGetNodeIDExpectation
}

type SlotDetailsMockGetNodeIDExpectation struct {
	result *SlotDetailsMockGetNodeIDResult
}

type SlotDetailsMockGetNodeIDResult struct {
	r uint32
}

//Expect specifies that invocation of SlotDetails.GetNodeID is expected from 1 to Infinity times
func (m *mSlotDetailsMockGetNodeID) Expect() *mSlotDetailsMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotDetails.GetNodeID
func (m *mSlotDetailsMockGetNodeID) Return(r uint32) *SlotDetailsMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &SlotDetailsMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotDetails.GetNodeID is expected once
func (m *mSlotDetailsMockGetNodeID) ExpectOnce() *SlotDetailsMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &SlotDetailsMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotDetailsMockGetNodeIDExpectation) Return(r uint32) {
	e.result = &SlotDetailsMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of SlotDetails.GetNodeID method
func (m *mSlotDetailsMockGetNodeID) Set(f func() (r uint32)) *SlotDetailsMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotDetails interface
func (m *SlotDetailsMock) GetNodeID() (r uint32) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotDetailsMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to SlotDetailsMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of SlotDetailsMock.GetNodeIDFunc invocations
func (m *SlotDetailsMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of SlotDetailsMock.GetNodeID invocations
func (m *SlotDetailsMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *SlotDetailsMock) GetNodeIDFinished() bool {
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

type mSlotDetailsMockGetPulseData struct {
	mock              *SlotDetailsMock
	mainExpectation   *SlotDetailsMockGetPulseDataExpectation
	expectationSeries []*SlotDetailsMockGetPulseDataExpectation
}

type SlotDetailsMockGetPulseDataExpectation struct {
	result *SlotDetailsMockGetPulseDataResult
}

type SlotDetailsMockGetPulseDataResult struct {
	r insolar.Pulse
}

//Expect specifies that invocation of SlotDetails.GetPulseData is expected from 1 to Infinity times
func (m *mSlotDetailsMockGetPulseData) Expect() *mSlotDetailsMockGetPulseData {
	m.mock.GetPulseDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetPulseDataExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotDetails.GetPulseData
func (m *mSlotDetailsMockGetPulseData) Return(r insolar.Pulse) *SlotDetailsMock {
	m.mock.GetPulseDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetPulseDataExpectation{}
	}
	m.mainExpectation.result = &SlotDetailsMockGetPulseDataResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotDetails.GetPulseData is expected once
func (m *mSlotDetailsMockGetPulseData) ExpectOnce() *SlotDetailsMockGetPulseDataExpectation {
	m.mock.GetPulseDataFunc = nil
	m.mainExpectation = nil

	expectation := &SlotDetailsMockGetPulseDataExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotDetailsMockGetPulseDataExpectation) Return(r insolar.Pulse) {
	e.result = &SlotDetailsMockGetPulseDataResult{r}
}

//Set uses given function f as a mock of SlotDetails.GetPulseData method
func (m *mSlotDetailsMockGetPulseData) Set(f func() (r insolar.Pulse)) *SlotDetailsMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseDataFunc = f
	return m.mock
}

//GetPulseData implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotDetails interface
func (m *SlotDetailsMock) GetPulseData() (r insolar.Pulse) {
	counter := atomic.AddUint64(&m.GetPulseDataPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseDataCounter, 1)

	if len(m.GetPulseDataMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseDataMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotDetailsMock.GetPulseData.")
			return
		}

		result := m.GetPulseDataMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetPulseData")
			return
		}

		r = result.r

		return
	}

	if m.GetPulseDataMock.mainExpectation != nil {

		result := m.GetPulseDataMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetPulseData")
		}

		r = result.r

		return
	}

	if m.GetPulseDataFunc == nil {
		m.t.Fatalf("Unexpected call to SlotDetailsMock.GetPulseData.")
		return
	}

	return m.GetPulseDataFunc()
}

//GetPulseDataMinimockCounter returns a count of SlotDetailsMock.GetPulseDataFunc invocations
func (m *SlotDetailsMock) GetPulseDataMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseDataCounter)
}

//GetPulseDataMinimockPreCounter returns the value of SlotDetailsMock.GetPulseData invocations
func (m *SlotDetailsMock) GetPulseDataMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseDataPreCounter)
}

//GetPulseDataFinished returns true if mock invocations count is ok
func (m *SlotDetailsMock) GetPulseDataFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulseDataMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulseDataCounter) == uint64(len(m.GetPulseDataMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulseDataMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulseDataCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulseDataFunc != nil {
		return atomic.LoadUint64(&m.GetPulseDataCounter) > 0
	}

	return true
}

type mSlotDetailsMockGetPulseNumber struct {
	mock              *SlotDetailsMock
	mainExpectation   *SlotDetailsMockGetPulseNumberExpectation
	expectationSeries []*SlotDetailsMockGetPulseNumberExpectation
}

type SlotDetailsMockGetPulseNumberExpectation struct {
	result *SlotDetailsMockGetPulseNumberResult
}

type SlotDetailsMockGetPulseNumberResult struct {
	r insolar.PulseNumber
}

//Expect specifies that invocation of SlotDetails.GetPulseNumber is expected from 1 to Infinity times
func (m *mSlotDetailsMockGetPulseNumber) Expect() *mSlotDetailsMockGetPulseNumber {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetPulseNumberExpectation{}
	}

	return m
}

//Return specifies results of invocation of SlotDetails.GetPulseNumber
func (m *mSlotDetailsMockGetPulseNumber) Return(r insolar.PulseNumber) *SlotDetailsMock {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SlotDetailsMockGetPulseNumberExpectation{}
	}
	m.mainExpectation.result = &SlotDetailsMockGetPulseNumberResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SlotDetails.GetPulseNumber is expected once
func (m *mSlotDetailsMockGetPulseNumber) ExpectOnce() *SlotDetailsMockGetPulseNumberExpectation {
	m.mock.GetPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &SlotDetailsMockGetPulseNumberExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SlotDetailsMockGetPulseNumberExpectation) Return(r insolar.PulseNumber) {
	e.result = &SlotDetailsMockGetPulseNumberResult{r}
}

//Set uses given function f as a mock of SlotDetails.GetPulseNumber method
func (m *mSlotDetailsMockGetPulseNumber) Set(f func() (r insolar.PulseNumber)) *SlotDetailsMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseNumberFunc = f
	return m.mock
}

//GetPulseNumber implements github.com/insolar/insolar/conveyor/interfaces/slot.SlotDetails interface
func (m *SlotDetailsMock) GetPulseNumber() (r insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.GetPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseNumberCounter, 1)

	if len(m.GetPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SlotDetailsMock.GetPulseNumber.")
			return
		}

		result := m.GetPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetPulseNumber")
			return
		}

		r = result.r

		return
	}

	if m.GetPulseNumberMock.mainExpectation != nil {

		result := m.GetPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SlotDetailsMock.GetPulseNumber")
		}

		r = result.r

		return
	}

	if m.GetPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to SlotDetailsMock.GetPulseNumber.")
		return
	}

	return m.GetPulseNumberFunc()
}

//GetPulseNumberMinimockCounter returns a count of SlotDetailsMock.GetPulseNumberFunc invocations
func (m *SlotDetailsMock) GetPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberCounter)
}

//GetPulseNumberMinimockPreCounter returns the value of SlotDetailsMock.GetPulseNumber invocations
func (m *SlotDetailsMock) GetPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberPreCounter)
}

//GetPulseNumberFinished returns true if mock invocations count is ok
func (m *SlotDetailsMock) GetPulseNumberFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulseNumberMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) == uint64(len(m.GetPulseNumberMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulseNumberMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulseNumberFunc != nil {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotDetailsMock) ValidateCallCounters() {

	if !m.GetNodeDataFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetNodeData")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetNodeID")
	}

	if !m.GetPulseDataFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetPulseData")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetPulseNumber")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SlotDetailsMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SlotDetailsMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SlotDetailsMock) MinimockFinish() {

	if !m.GetNodeDataFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetNodeData")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetNodeID")
	}

	if !m.GetPulseDataFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetPulseData")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to SlotDetailsMock.GetPulseNumber")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SlotDetailsMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SlotDetailsMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetNodeDataFinished()
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetPulseDataFinished()
		ok = ok && m.GetPulseNumberFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetNodeDataFinished() {
				m.t.Error("Expected call to SlotDetailsMock.GetNodeData")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to SlotDetailsMock.GetNodeID")
			}

			if !m.GetPulseDataFinished() {
				m.t.Error("Expected call to SlotDetailsMock.GetPulseData")
			}

			if !m.GetPulseNumberFinished() {
				m.t.Error("Expected call to SlotDetailsMock.GetPulseNumber")
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
func (m *SlotDetailsMock) AllMocksCalled() bool {

	if !m.GetNodeDataFinished() {
		return false
	}

	if !m.GetNodeIDFinished() {
		return false
	}

	if !m.GetPulseDataFinished() {
		return false
	}

	if !m.GetPulseNumberFinished() {
		return false
	}

	return true
}

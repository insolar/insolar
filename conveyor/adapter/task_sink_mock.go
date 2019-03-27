package adapter

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "TaskSink" can be found in github.com/insolar/insolar/conveyor/adapter
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	adapterid "github.com/insolar/insolar/conveyor/adapter/adapterid"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//TaskSinkMock implements github.com/insolar/insolar/conveyor/adapter.TaskSink
type TaskSinkMock struct {
	t minimock.Tester

	CancelElementTasksFunc       func(p insolar.PulseNumber, p1 uint32)
	CancelElementTasksCounter    uint64
	CancelElementTasksPreCounter uint64
	CancelElementTasksMock       mTaskSinkMockCancelElementTasks

	CancelPulseTasksFunc       func(p insolar.PulseNumber)
	CancelPulseTasksCounter    uint64
	CancelPulseTasksPreCounter uint64
	CancelPulseTasksMock       mTaskSinkMockCancelPulseTasks

	FlushNodeTasksFunc       func(p uint32)
	FlushNodeTasksCounter    uint64
	FlushNodeTasksPreCounter uint64
	FlushNodeTasksMock       mTaskSinkMockFlushNodeTasks

	FlushPulseTasksFunc       func(p insolar.PulseNumber)
	FlushPulseTasksCounter    uint64
	FlushPulseTasksPreCounter uint64
	FlushPulseTasksMock       mTaskSinkMockFlushPulseTasks

	GetAdapterIDFunc       func() (r adapterid.ID)
	GetAdapterIDCounter    uint64
	GetAdapterIDPreCounter uint64
	GetAdapterIDMock       mTaskSinkMockGetAdapterID

	PushTaskFunc       func(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) (r error)
	PushTaskCounter    uint64
	PushTaskPreCounter uint64
	PushTaskMock       mTaskSinkMockPushTask
}

//NewTaskSinkMock returns a mock for github.com/insolar/insolar/conveyor/adapter.TaskSink
func NewTaskSinkMock(t minimock.Tester) *TaskSinkMock {
	m := &TaskSinkMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CancelElementTasksMock = mTaskSinkMockCancelElementTasks{mock: m}
	m.CancelPulseTasksMock = mTaskSinkMockCancelPulseTasks{mock: m}
	m.FlushNodeTasksMock = mTaskSinkMockFlushNodeTasks{mock: m}
	m.FlushPulseTasksMock = mTaskSinkMockFlushPulseTasks{mock: m}
	m.GetAdapterIDMock = mTaskSinkMockGetAdapterID{mock: m}
	m.PushTaskMock = mTaskSinkMockPushTask{mock: m}

	return m
}

type mTaskSinkMockCancelElementTasks struct {
	mock              *TaskSinkMock
	mainExpectation   *TaskSinkMockCancelElementTasksExpectation
	expectationSeries []*TaskSinkMockCancelElementTasksExpectation
}

type TaskSinkMockCancelElementTasksExpectation struct {
	input *TaskSinkMockCancelElementTasksInput
}

type TaskSinkMockCancelElementTasksInput struct {
	p  insolar.PulseNumber
	p1 uint32
}

//Expect specifies that invocation of TaskSink.CancelElementTasks is expected from 1 to Infinity times
func (m *mTaskSinkMockCancelElementTasks) Expect(p insolar.PulseNumber, p1 uint32) *mTaskSinkMockCancelElementTasks {
	m.mock.CancelElementTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockCancelElementTasksExpectation{}
	}
	m.mainExpectation.input = &TaskSinkMockCancelElementTasksInput{p, p1}
	return m
}

//Return specifies results of invocation of TaskSink.CancelElementTasks
func (m *mTaskSinkMockCancelElementTasks) Return() *TaskSinkMock {
	m.mock.CancelElementTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockCancelElementTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TaskSink.CancelElementTasks is expected once
func (m *mTaskSinkMockCancelElementTasks) ExpectOnce(p insolar.PulseNumber, p1 uint32) *TaskSinkMockCancelElementTasksExpectation {
	m.mock.CancelElementTasksFunc = nil
	m.mainExpectation = nil

	expectation := &TaskSinkMockCancelElementTasksExpectation{}
	expectation.input = &TaskSinkMockCancelElementTasksInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TaskSink.CancelElementTasks method
func (m *mTaskSinkMockCancelElementTasks) Set(f func(p insolar.PulseNumber, p1 uint32)) *TaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CancelElementTasksFunc = f
	return m.mock
}

//CancelElementTasks implements github.com/insolar/insolar/conveyor/adapter.TaskSink interface
func (m *TaskSinkMock) CancelElementTasks(p insolar.PulseNumber, p1 uint32) {
	counter := atomic.AddUint64(&m.CancelElementTasksPreCounter, 1)
	defer atomic.AddUint64(&m.CancelElementTasksCounter, 1)

	if len(m.CancelElementTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CancelElementTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskSinkMock.CancelElementTasks. %v %v", p, p1)
			return
		}

		input := m.CancelElementTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskSinkMockCancelElementTasksInput{p, p1}, "TaskSink.CancelElementTasks got unexpected parameters")

		return
	}

	if m.CancelElementTasksMock.mainExpectation != nil {

		input := m.CancelElementTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskSinkMockCancelElementTasksInput{p, p1}, "TaskSink.CancelElementTasks got unexpected parameters")
		}

		return
	}

	if m.CancelElementTasksFunc == nil {
		m.t.Fatalf("Unexpected call to TaskSinkMock.CancelElementTasks. %v %v", p, p1)
		return
	}

	m.CancelElementTasksFunc(p, p1)
}

//CancelElementTasksMinimockCounter returns a count of TaskSinkMock.CancelElementTasksFunc invocations
func (m *TaskSinkMock) CancelElementTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CancelElementTasksCounter)
}

//CancelElementTasksMinimockPreCounter returns the value of TaskSinkMock.CancelElementTasks invocations
func (m *TaskSinkMock) CancelElementTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CancelElementTasksPreCounter)
}

//CancelElementTasksFinished returns true if mock invocations count is ok
func (m *TaskSinkMock) CancelElementTasksFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CancelElementTasksMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CancelElementTasksCounter) == uint64(len(m.CancelElementTasksMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CancelElementTasksMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CancelElementTasksCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CancelElementTasksFunc != nil {
		return atomic.LoadUint64(&m.CancelElementTasksCounter) > 0
	}

	return true
}

type mTaskSinkMockCancelPulseTasks struct {
	mock              *TaskSinkMock
	mainExpectation   *TaskSinkMockCancelPulseTasksExpectation
	expectationSeries []*TaskSinkMockCancelPulseTasksExpectation
}

type TaskSinkMockCancelPulseTasksExpectation struct {
	input *TaskSinkMockCancelPulseTasksInput
}

type TaskSinkMockCancelPulseTasksInput struct {
	p insolar.PulseNumber
}

//Expect specifies that invocation of TaskSink.CancelPulseTasks is expected from 1 to Infinity times
func (m *mTaskSinkMockCancelPulseTasks) Expect(p insolar.PulseNumber) *mTaskSinkMockCancelPulseTasks {
	m.mock.CancelPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockCancelPulseTasksExpectation{}
	}
	m.mainExpectation.input = &TaskSinkMockCancelPulseTasksInput{p}
	return m
}

//Return specifies results of invocation of TaskSink.CancelPulseTasks
func (m *mTaskSinkMockCancelPulseTasks) Return() *TaskSinkMock {
	m.mock.CancelPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockCancelPulseTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TaskSink.CancelPulseTasks is expected once
func (m *mTaskSinkMockCancelPulseTasks) ExpectOnce(p insolar.PulseNumber) *TaskSinkMockCancelPulseTasksExpectation {
	m.mock.CancelPulseTasksFunc = nil
	m.mainExpectation = nil

	expectation := &TaskSinkMockCancelPulseTasksExpectation{}
	expectation.input = &TaskSinkMockCancelPulseTasksInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TaskSink.CancelPulseTasks method
func (m *mTaskSinkMockCancelPulseTasks) Set(f func(p insolar.PulseNumber)) *TaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CancelPulseTasksFunc = f
	return m.mock
}

//CancelPulseTasks implements github.com/insolar/insolar/conveyor/adapter.TaskSink interface
func (m *TaskSinkMock) CancelPulseTasks(p insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.CancelPulseTasksPreCounter, 1)
	defer atomic.AddUint64(&m.CancelPulseTasksCounter, 1)

	if len(m.CancelPulseTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CancelPulseTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskSinkMock.CancelPulseTasks. %v", p)
			return
		}

		input := m.CancelPulseTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskSinkMockCancelPulseTasksInput{p}, "TaskSink.CancelPulseTasks got unexpected parameters")

		return
	}

	if m.CancelPulseTasksMock.mainExpectation != nil {

		input := m.CancelPulseTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskSinkMockCancelPulseTasksInput{p}, "TaskSink.CancelPulseTasks got unexpected parameters")
		}

		return
	}

	if m.CancelPulseTasksFunc == nil {
		m.t.Fatalf("Unexpected call to TaskSinkMock.CancelPulseTasks. %v", p)
		return
	}

	m.CancelPulseTasksFunc(p)
}

//CancelPulseTasksMinimockCounter returns a count of TaskSinkMock.CancelPulseTasksFunc invocations
func (m *TaskSinkMock) CancelPulseTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CancelPulseTasksCounter)
}

//CancelPulseTasksMinimockPreCounter returns the value of TaskSinkMock.CancelPulseTasks invocations
func (m *TaskSinkMock) CancelPulseTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CancelPulseTasksPreCounter)
}

//CancelPulseTasksFinished returns true if mock invocations count is ok
func (m *TaskSinkMock) CancelPulseTasksFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CancelPulseTasksMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CancelPulseTasksCounter) == uint64(len(m.CancelPulseTasksMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CancelPulseTasksMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CancelPulseTasksCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CancelPulseTasksFunc != nil {
		return atomic.LoadUint64(&m.CancelPulseTasksCounter) > 0
	}

	return true
}

type mTaskSinkMockFlushNodeTasks struct {
	mock              *TaskSinkMock
	mainExpectation   *TaskSinkMockFlushNodeTasksExpectation
	expectationSeries []*TaskSinkMockFlushNodeTasksExpectation
}

type TaskSinkMockFlushNodeTasksExpectation struct {
	input *TaskSinkMockFlushNodeTasksInput
}

type TaskSinkMockFlushNodeTasksInput struct {
	p uint32
}

//Expect specifies that invocation of TaskSink.FlushNodeTasks is expected from 1 to Infinity times
func (m *mTaskSinkMockFlushNodeTasks) Expect(p uint32) *mTaskSinkMockFlushNodeTasks {
	m.mock.FlushNodeTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockFlushNodeTasksExpectation{}
	}
	m.mainExpectation.input = &TaskSinkMockFlushNodeTasksInput{p}
	return m
}

//Return specifies results of invocation of TaskSink.FlushNodeTasks
func (m *mTaskSinkMockFlushNodeTasks) Return() *TaskSinkMock {
	m.mock.FlushNodeTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockFlushNodeTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TaskSink.FlushNodeTasks is expected once
func (m *mTaskSinkMockFlushNodeTasks) ExpectOnce(p uint32) *TaskSinkMockFlushNodeTasksExpectation {
	m.mock.FlushNodeTasksFunc = nil
	m.mainExpectation = nil

	expectation := &TaskSinkMockFlushNodeTasksExpectation{}
	expectation.input = &TaskSinkMockFlushNodeTasksInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TaskSink.FlushNodeTasks method
func (m *mTaskSinkMockFlushNodeTasks) Set(f func(p uint32)) *TaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FlushNodeTasksFunc = f
	return m.mock
}

//FlushNodeTasks implements github.com/insolar/insolar/conveyor/adapter.TaskSink interface
func (m *TaskSinkMock) FlushNodeTasks(p uint32) {
	counter := atomic.AddUint64(&m.FlushNodeTasksPreCounter, 1)
	defer atomic.AddUint64(&m.FlushNodeTasksCounter, 1)

	if len(m.FlushNodeTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FlushNodeTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskSinkMock.FlushNodeTasks. %v", p)
			return
		}

		input := m.FlushNodeTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskSinkMockFlushNodeTasksInput{p}, "TaskSink.FlushNodeTasks got unexpected parameters")

		return
	}

	if m.FlushNodeTasksMock.mainExpectation != nil {

		input := m.FlushNodeTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskSinkMockFlushNodeTasksInput{p}, "TaskSink.FlushNodeTasks got unexpected parameters")
		}

		return
	}

	if m.FlushNodeTasksFunc == nil {
		m.t.Fatalf("Unexpected call to TaskSinkMock.FlushNodeTasks. %v", p)
		return
	}

	m.FlushNodeTasksFunc(p)
}

//FlushNodeTasksMinimockCounter returns a count of TaskSinkMock.FlushNodeTasksFunc invocations
func (m *TaskSinkMock) FlushNodeTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FlushNodeTasksCounter)
}

//FlushNodeTasksMinimockPreCounter returns the value of TaskSinkMock.FlushNodeTasks invocations
func (m *TaskSinkMock) FlushNodeTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FlushNodeTasksPreCounter)
}

//FlushNodeTasksFinished returns true if mock invocations count is ok
func (m *TaskSinkMock) FlushNodeTasksFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FlushNodeTasksMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FlushNodeTasksCounter) == uint64(len(m.FlushNodeTasksMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FlushNodeTasksMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FlushNodeTasksCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FlushNodeTasksFunc != nil {
		return atomic.LoadUint64(&m.FlushNodeTasksCounter) > 0
	}

	return true
}

type mTaskSinkMockFlushPulseTasks struct {
	mock              *TaskSinkMock
	mainExpectation   *TaskSinkMockFlushPulseTasksExpectation
	expectationSeries []*TaskSinkMockFlushPulseTasksExpectation
}

type TaskSinkMockFlushPulseTasksExpectation struct {
	input *TaskSinkMockFlushPulseTasksInput
}

type TaskSinkMockFlushPulseTasksInput struct {
	p insolar.PulseNumber
}

//Expect specifies that invocation of TaskSink.FlushPulseTasks is expected from 1 to Infinity times
func (m *mTaskSinkMockFlushPulseTasks) Expect(p insolar.PulseNumber) *mTaskSinkMockFlushPulseTasks {
	m.mock.FlushPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockFlushPulseTasksExpectation{}
	}
	m.mainExpectation.input = &TaskSinkMockFlushPulseTasksInput{p}
	return m
}

//Return specifies results of invocation of TaskSink.FlushPulseTasks
func (m *mTaskSinkMockFlushPulseTasks) Return() *TaskSinkMock {
	m.mock.FlushPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockFlushPulseTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TaskSink.FlushPulseTasks is expected once
func (m *mTaskSinkMockFlushPulseTasks) ExpectOnce(p insolar.PulseNumber) *TaskSinkMockFlushPulseTasksExpectation {
	m.mock.FlushPulseTasksFunc = nil
	m.mainExpectation = nil

	expectation := &TaskSinkMockFlushPulseTasksExpectation{}
	expectation.input = &TaskSinkMockFlushPulseTasksInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TaskSink.FlushPulseTasks method
func (m *mTaskSinkMockFlushPulseTasks) Set(f func(p insolar.PulseNumber)) *TaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FlushPulseTasksFunc = f
	return m.mock
}

//FlushPulseTasks implements github.com/insolar/insolar/conveyor/adapter.TaskSink interface
func (m *TaskSinkMock) FlushPulseTasks(p insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.FlushPulseTasksPreCounter, 1)
	defer atomic.AddUint64(&m.FlushPulseTasksCounter, 1)

	if len(m.FlushPulseTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FlushPulseTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskSinkMock.FlushPulseTasks. %v", p)
			return
		}

		input := m.FlushPulseTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskSinkMockFlushPulseTasksInput{p}, "TaskSink.FlushPulseTasks got unexpected parameters")

		return
	}

	if m.FlushPulseTasksMock.mainExpectation != nil {

		input := m.FlushPulseTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskSinkMockFlushPulseTasksInput{p}, "TaskSink.FlushPulseTasks got unexpected parameters")
		}

		return
	}

	if m.FlushPulseTasksFunc == nil {
		m.t.Fatalf("Unexpected call to TaskSinkMock.FlushPulseTasks. %v", p)
		return
	}

	m.FlushPulseTasksFunc(p)
}

//FlushPulseTasksMinimockCounter returns a count of TaskSinkMock.FlushPulseTasksFunc invocations
func (m *TaskSinkMock) FlushPulseTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FlushPulseTasksCounter)
}

//FlushPulseTasksMinimockPreCounter returns the value of TaskSinkMock.FlushPulseTasks invocations
func (m *TaskSinkMock) FlushPulseTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FlushPulseTasksPreCounter)
}

//FlushPulseTasksFinished returns true if mock invocations count is ok
func (m *TaskSinkMock) FlushPulseTasksFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FlushPulseTasksMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FlushPulseTasksCounter) == uint64(len(m.FlushPulseTasksMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FlushPulseTasksMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FlushPulseTasksCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FlushPulseTasksFunc != nil {
		return atomic.LoadUint64(&m.FlushPulseTasksCounter) > 0
	}

	return true
}

type mTaskSinkMockGetAdapterID struct {
	mock              *TaskSinkMock
	mainExpectation   *TaskSinkMockGetAdapterIDExpectation
	expectationSeries []*TaskSinkMockGetAdapterIDExpectation
}

type TaskSinkMockGetAdapterIDExpectation struct {
	result *TaskSinkMockGetAdapterIDResult
}

type TaskSinkMockGetAdapterIDResult struct {
	r adapterid.ID
}

//Expect specifies that invocation of TaskSink.GetAdapterID is expected from 1 to Infinity times
func (m *mTaskSinkMockGetAdapterID) Expect() *mTaskSinkMockGetAdapterID {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockGetAdapterIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of TaskSink.GetAdapterID
func (m *mTaskSinkMockGetAdapterID) Return(r adapterid.ID) *TaskSinkMock {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockGetAdapterIDExpectation{}
	}
	m.mainExpectation.result = &TaskSinkMockGetAdapterIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of TaskSink.GetAdapterID is expected once
func (m *mTaskSinkMockGetAdapterID) ExpectOnce() *TaskSinkMockGetAdapterIDExpectation {
	m.mock.GetAdapterIDFunc = nil
	m.mainExpectation = nil

	expectation := &TaskSinkMockGetAdapterIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TaskSinkMockGetAdapterIDExpectation) Return(r adapterid.ID) {
	e.result = &TaskSinkMockGetAdapterIDResult{r}
}

//Set uses given function f as a mock of TaskSink.GetAdapterID method
func (m *mTaskSinkMockGetAdapterID) Set(f func() (r adapterid.ID)) *TaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAdapterIDFunc = f
	return m.mock
}

//GetAdapterID implements github.com/insolar/insolar/conveyor/adapter.TaskSink interface
func (m *TaskSinkMock) GetAdapterID() (r adapterid.ID) {
	counter := atomic.AddUint64(&m.GetAdapterIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetAdapterIDCounter, 1)

	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAdapterIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskSinkMock.GetAdapterID.")
			return
		}

		result := m.GetAdapterIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TaskSinkMock.GetAdapterID")
			return
		}

		r = result.r

		return
	}

	if m.GetAdapterIDMock.mainExpectation != nil {

		result := m.GetAdapterIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TaskSinkMock.GetAdapterID")
		}

		r = result.r

		return
	}

	if m.GetAdapterIDFunc == nil {
		m.t.Fatalf("Unexpected call to TaskSinkMock.GetAdapterID.")
		return
	}

	return m.GetAdapterIDFunc()
}

//GetAdapterIDMinimockCounter returns a count of TaskSinkMock.GetAdapterIDFunc invocations
func (m *TaskSinkMock) GetAdapterIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDCounter)
}

//GetAdapterIDMinimockPreCounter returns the value of TaskSinkMock.GetAdapterID invocations
func (m *TaskSinkMock) GetAdapterIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDPreCounter)
}

//GetAdapterIDFinished returns true if mock invocations count is ok
func (m *TaskSinkMock) GetAdapterIDFinished() bool {
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

type mTaskSinkMockPushTask struct {
	mock              *TaskSinkMock
	mainExpectation   *TaskSinkMockPushTaskExpectation
	expectationSeries []*TaskSinkMockPushTaskExpectation
}

type TaskSinkMockPushTaskExpectation struct {
	input  *TaskSinkMockPushTaskInput
	result *TaskSinkMockPushTaskResult
}

type TaskSinkMockPushTaskInput struct {
	p  AdapterToSlotResponseSink
	p1 uint32
	p2 uint32
	p3 interface{}
}

type TaskSinkMockPushTaskResult struct {
	r error
}

//Expect specifies that invocation of TaskSink.PushTask is expected from 1 to Infinity times
func (m *mTaskSinkMockPushTask) Expect(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) *mTaskSinkMockPushTask {
	m.mock.PushTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockPushTaskExpectation{}
	}
	m.mainExpectation.input = &TaskSinkMockPushTaskInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of TaskSink.PushTask
func (m *mTaskSinkMockPushTask) Return(r error) *TaskSinkMock {
	m.mock.PushTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskSinkMockPushTaskExpectation{}
	}
	m.mainExpectation.result = &TaskSinkMockPushTaskResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of TaskSink.PushTask is expected once
func (m *mTaskSinkMockPushTask) ExpectOnce(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) *TaskSinkMockPushTaskExpectation {
	m.mock.PushTaskFunc = nil
	m.mainExpectation = nil

	expectation := &TaskSinkMockPushTaskExpectation{}
	expectation.input = &TaskSinkMockPushTaskInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TaskSinkMockPushTaskExpectation) Return(r error) {
	e.result = &TaskSinkMockPushTaskResult{r}
}

//Set uses given function f as a mock of TaskSink.PushTask method
func (m *mTaskSinkMockPushTask) Set(f func(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) (r error)) *TaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PushTaskFunc = f
	return m.mock
}

//PushTask implements github.com/insolar/insolar/conveyor/adapter.TaskSink interface
func (m *TaskSinkMock) PushTask(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) (r error) {
	counter := atomic.AddUint64(&m.PushTaskPreCounter, 1)
	defer atomic.AddUint64(&m.PushTaskCounter, 1)

	if len(m.PushTaskMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PushTaskMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskSinkMock.PushTask. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.PushTaskMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskSinkMockPushTaskInput{p, p1, p2, p3}, "TaskSink.PushTask got unexpected parameters")

		result := m.PushTaskMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TaskSinkMock.PushTask")
			return
		}

		r = result.r

		return
	}

	if m.PushTaskMock.mainExpectation != nil {

		input := m.PushTaskMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskSinkMockPushTaskInput{p, p1, p2, p3}, "TaskSink.PushTask got unexpected parameters")
		}

		result := m.PushTaskMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TaskSinkMock.PushTask")
		}

		r = result.r

		return
	}

	if m.PushTaskFunc == nil {
		m.t.Fatalf("Unexpected call to TaskSinkMock.PushTask. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.PushTaskFunc(p, p1, p2, p3)
}

//PushTaskMinimockCounter returns a count of TaskSinkMock.PushTaskFunc invocations
func (m *TaskSinkMock) PushTaskMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PushTaskCounter)
}

//PushTaskMinimockPreCounter returns the value of TaskSinkMock.PushTask invocations
func (m *TaskSinkMock) PushTaskMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PushTaskPreCounter)
}

//PushTaskFinished returns true if mock invocations count is ok
func (m *TaskSinkMock) PushTaskFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PushTaskMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PushTaskCounter) == uint64(len(m.PushTaskMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PushTaskMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PushTaskCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PushTaskFunc != nil {
		return atomic.LoadUint64(&m.PushTaskCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TaskSinkMock) ValidateCallCounters() {

	if !m.CancelElementTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.CancelElementTasks")
	}

	if !m.CancelPulseTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.CancelPulseTasks")
	}

	if !m.FlushNodeTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.FlushNodeTasks")
	}

	if !m.FlushPulseTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.FlushPulseTasks")
	}

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.GetAdapterID")
	}

	if !m.PushTaskFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.PushTask")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TaskSinkMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TaskSinkMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TaskSinkMock) MinimockFinish() {

	if !m.CancelElementTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.CancelElementTasks")
	}

	if !m.CancelPulseTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.CancelPulseTasks")
	}

	if !m.FlushNodeTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.FlushNodeTasks")
	}

	if !m.FlushPulseTasksFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.FlushPulseTasks")
	}

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.GetAdapterID")
	}

	if !m.PushTaskFinished() {
		m.t.Fatal("Expected call to TaskSinkMock.PushTask")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TaskSinkMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TaskSinkMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CancelElementTasksFinished()
		ok = ok && m.CancelPulseTasksFinished()
		ok = ok && m.FlushNodeTasksFinished()
		ok = ok && m.FlushPulseTasksFinished()
		ok = ok && m.GetAdapterIDFinished()
		ok = ok && m.PushTaskFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CancelElementTasksFinished() {
				m.t.Error("Expected call to TaskSinkMock.CancelElementTasks")
			}

			if !m.CancelPulseTasksFinished() {
				m.t.Error("Expected call to TaskSinkMock.CancelPulseTasks")
			}

			if !m.FlushNodeTasksFinished() {
				m.t.Error("Expected call to TaskSinkMock.FlushNodeTasks")
			}

			if !m.FlushPulseTasksFinished() {
				m.t.Error("Expected call to TaskSinkMock.FlushPulseTasks")
			}

			if !m.GetAdapterIDFinished() {
				m.t.Error("Expected call to TaskSinkMock.GetAdapterID")
			}

			if !m.PushTaskFinished() {
				m.t.Error("Expected call to TaskSinkMock.PushTask")
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
func (m *TaskSinkMock) AllMocksCalled() bool {

	if !m.CancelElementTasksFinished() {
		return false
	}

	if !m.CancelPulseTasksFinished() {
		return false
	}

	if !m.FlushNodeTasksFinished() {
		return false
	}

	if !m.FlushPulseTasksFinished() {
		return false
	}

	if !m.GetAdapterIDFinished() {
		return false
	}

	if !m.PushTaskFinished() {
		return false
	}

	return true
}

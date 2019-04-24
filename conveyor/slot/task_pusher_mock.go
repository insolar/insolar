package slot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "TaskPusher" can be found in github.com/insolar/insolar/conveyor/slot
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	queue "github.com/insolar/insolar/conveyor/queue"

	testify_assert "github.com/stretchr/testify/assert"
)

//TaskPusherMock implements github.com/insolar/insolar/conveyor/slot.TaskPusher
type TaskPusherMock struct {
	t minimock.Tester

	PushSignalFunc       func(p uint32, p1 queue.SyncDone) (r error)
	PushSignalCounter    uint64
	PushSignalPreCounter uint64
	PushSignalMock       mTaskPusherMockPushSignal

	SinkPushFunc       func(p interface{}) (r error)
	SinkPushCounter    uint64
	SinkPushPreCounter uint64
	SinkPushMock       mTaskPusherMockSinkPush

	SinkPushAllFunc       func(p []interface{}) (r error)
	SinkPushAllCounter    uint64
	SinkPushAllPreCounter uint64
	SinkPushAllMock       mTaskPusherMockSinkPushAll
}

//NewTaskPusherMock returns a mock for github.com/insolar/insolar/conveyor/slot.TaskPusher
func NewTaskPusherMock(t minimock.Tester) *TaskPusherMock {
	m := &TaskPusherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PushSignalMock = mTaskPusherMockPushSignal{mock: m}
	m.SinkPushMock = mTaskPusherMockSinkPush{mock: m}
	m.SinkPushAllMock = mTaskPusherMockSinkPushAll{mock: m}

	return m
}

type mTaskPusherMockPushSignal struct {
	mock              *TaskPusherMock
	mainExpectation   *TaskPusherMockPushSignalExpectation
	expectationSeries []*TaskPusherMockPushSignalExpectation
}

type TaskPusherMockPushSignalExpectation struct {
	input  *TaskPusherMockPushSignalInput
	result *TaskPusherMockPushSignalResult
}

type TaskPusherMockPushSignalInput struct {
	p  uint32
	p1 queue.SyncDone
}

type TaskPusherMockPushSignalResult struct {
	r error
}

//Expect specifies that invocation of TaskPusher.PushSignal is expected from 1 to Infinity times
func (m *mTaskPusherMockPushSignal) Expect(p uint32, p1 queue.SyncDone) *mTaskPusherMockPushSignal {
	m.mock.PushSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskPusherMockPushSignalExpectation{}
	}
	m.mainExpectation.input = &TaskPusherMockPushSignalInput{p, p1}
	return m
}

//Return specifies results of invocation of TaskPusher.PushSignal
func (m *mTaskPusherMockPushSignal) Return(r error) *TaskPusherMock {
	m.mock.PushSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskPusherMockPushSignalExpectation{}
	}
	m.mainExpectation.result = &TaskPusherMockPushSignalResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of TaskPusher.PushSignal is expected once
func (m *mTaskPusherMockPushSignal) ExpectOnce(p uint32, p1 queue.SyncDone) *TaskPusherMockPushSignalExpectation {
	m.mock.PushSignalFunc = nil
	m.mainExpectation = nil

	expectation := &TaskPusherMockPushSignalExpectation{}
	expectation.input = &TaskPusherMockPushSignalInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TaskPusherMockPushSignalExpectation) Return(r error) {
	e.result = &TaskPusherMockPushSignalResult{r}
}

//Set uses given function f as a mock of TaskPusher.PushSignal method
func (m *mTaskPusherMockPushSignal) Set(f func(p uint32, p1 queue.SyncDone) (r error)) *TaskPusherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PushSignalFunc = f
	return m.mock
}

//PushSignal implements github.com/insolar/insolar/conveyor/slot.TaskPusher interface
func (m *TaskPusherMock) PushSignal(p uint32, p1 queue.SyncDone) (r error) {
	counter := atomic.AddUint64(&m.PushSignalPreCounter, 1)
	defer atomic.AddUint64(&m.PushSignalCounter, 1)

	if len(m.PushSignalMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PushSignalMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskPusherMock.PushSignal. %v %v", p, p1)
			return
		}

		input := m.PushSignalMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskPusherMockPushSignalInput{p, p1}, "TaskPusher.PushSignal got unexpected parameters")

		result := m.PushSignalMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TaskPusherMock.PushSignal")
			return
		}

		r = result.r

		return
	}

	if m.PushSignalMock.mainExpectation != nil {

		input := m.PushSignalMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskPusherMockPushSignalInput{p, p1}, "TaskPusher.PushSignal got unexpected parameters")
		}

		result := m.PushSignalMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TaskPusherMock.PushSignal")
		}

		r = result.r

		return
	}

	if m.PushSignalFunc == nil {
		m.t.Fatalf("Unexpected call to TaskPusherMock.PushSignal. %v %v", p, p1)
		return
	}

	return m.PushSignalFunc(p, p1)
}

//PushSignalMinimockCounter returns a count of TaskPusherMock.PushSignalFunc invocations
func (m *TaskPusherMock) PushSignalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PushSignalCounter)
}

//PushSignalMinimockPreCounter returns the value of TaskPusherMock.PushSignal invocations
func (m *TaskPusherMock) PushSignalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PushSignalPreCounter)
}

//PushSignalFinished returns true if mock invocations count is ok
func (m *TaskPusherMock) PushSignalFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PushSignalMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PushSignalCounter) == uint64(len(m.PushSignalMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PushSignalMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PushSignalCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PushSignalFunc != nil {
		return atomic.LoadUint64(&m.PushSignalCounter) > 0
	}

	return true
}

type mTaskPusherMockSinkPush struct {
	mock              *TaskPusherMock
	mainExpectation   *TaskPusherMockSinkPushExpectation
	expectationSeries []*TaskPusherMockSinkPushExpectation
}

type TaskPusherMockSinkPushExpectation struct {
	input  *TaskPusherMockSinkPushInput
	result *TaskPusherMockSinkPushResult
}

type TaskPusherMockSinkPushInput struct {
	p interface{}
}

type TaskPusherMockSinkPushResult struct {
	r error
}

//Expect specifies that invocation of TaskPusher.SinkPush is expected from 1 to Infinity times
func (m *mTaskPusherMockSinkPush) Expect(p interface{}) *mTaskPusherMockSinkPush {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskPusherMockSinkPushExpectation{}
	}
	m.mainExpectation.input = &TaskPusherMockSinkPushInput{p}
	return m
}

//Return specifies results of invocation of TaskPusher.SinkPush
func (m *mTaskPusherMockSinkPush) Return(r error) *TaskPusherMock {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskPusherMockSinkPushExpectation{}
	}
	m.mainExpectation.result = &TaskPusherMockSinkPushResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of TaskPusher.SinkPush is expected once
func (m *mTaskPusherMockSinkPush) ExpectOnce(p interface{}) *TaskPusherMockSinkPushExpectation {
	m.mock.SinkPushFunc = nil
	m.mainExpectation = nil

	expectation := &TaskPusherMockSinkPushExpectation{}
	expectation.input = &TaskPusherMockSinkPushInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TaskPusherMockSinkPushExpectation) Return(r error) {
	e.result = &TaskPusherMockSinkPushResult{r}
}

//Set uses given function f as a mock of TaskPusher.SinkPush method
func (m *mTaskPusherMockSinkPush) Set(f func(p interface{}) (r error)) *TaskPusherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushFunc = f
	return m.mock
}

//SinkPush implements github.com/insolar/insolar/conveyor/slot.TaskPusher interface
func (m *TaskPusherMock) SinkPush(p interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushCounter, 1)

	if len(m.SinkPushMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskPusherMock.SinkPush. %v", p)
			return
		}

		input := m.SinkPushMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskPusherMockSinkPushInput{p}, "TaskPusher.SinkPush got unexpected parameters")

		result := m.SinkPushMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TaskPusherMock.SinkPush")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushMock.mainExpectation != nil {

		input := m.SinkPushMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskPusherMockSinkPushInput{p}, "TaskPusher.SinkPush got unexpected parameters")
		}

		result := m.SinkPushMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TaskPusherMock.SinkPush")
		}

		r = result.r

		return
	}

	if m.SinkPushFunc == nil {
		m.t.Fatalf("Unexpected call to TaskPusherMock.SinkPush. %v", p)
		return
	}

	return m.SinkPushFunc(p)
}

//SinkPushMinimockCounter returns a count of TaskPusherMock.SinkPushFunc invocations
func (m *TaskPusherMock) SinkPushMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushCounter)
}

//SinkPushMinimockPreCounter returns the value of TaskPusherMock.SinkPush invocations
func (m *TaskPusherMock) SinkPushMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushPreCounter)
}

//SinkPushFinished returns true if mock invocations count is ok
func (m *TaskPusherMock) SinkPushFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SinkPushMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SinkPushCounter) == uint64(len(m.SinkPushMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SinkPushMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SinkPushCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SinkPushFunc != nil {
		return atomic.LoadUint64(&m.SinkPushCounter) > 0
	}

	return true
}

type mTaskPusherMockSinkPushAll struct {
	mock              *TaskPusherMock
	mainExpectation   *TaskPusherMockSinkPushAllExpectation
	expectationSeries []*TaskPusherMockSinkPushAllExpectation
}

type TaskPusherMockSinkPushAllExpectation struct {
	input  *TaskPusherMockSinkPushAllInput
	result *TaskPusherMockSinkPushAllResult
}

type TaskPusherMockSinkPushAllInput struct {
	p []interface{}
}

type TaskPusherMockSinkPushAllResult struct {
	r error
}

//Expect specifies that invocation of TaskPusher.SinkPushAll is expected from 1 to Infinity times
func (m *mTaskPusherMockSinkPushAll) Expect(p []interface{}) *mTaskPusherMockSinkPushAll {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskPusherMockSinkPushAllExpectation{}
	}
	m.mainExpectation.input = &TaskPusherMockSinkPushAllInput{p}
	return m
}

//Return specifies results of invocation of TaskPusher.SinkPushAll
func (m *mTaskPusherMockSinkPushAll) Return(r error) *TaskPusherMock {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TaskPusherMockSinkPushAllExpectation{}
	}
	m.mainExpectation.result = &TaskPusherMockSinkPushAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of TaskPusher.SinkPushAll is expected once
func (m *mTaskPusherMockSinkPushAll) ExpectOnce(p []interface{}) *TaskPusherMockSinkPushAllExpectation {
	m.mock.SinkPushAllFunc = nil
	m.mainExpectation = nil

	expectation := &TaskPusherMockSinkPushAllExpectation{}
	expectation.input = &TaskPusherMockSinkPushAllInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TaskPusherMockSinkPushAllExpectation) Return(r error) {
	e.result = &TaskPusherMockSinkPushAllResult{r}
}

//Set uses given function f as a mock of TaskPusher.SinkPushAll method
func (m *mTaskPusherMockSinkPushAll) Set(f func(p []interface{}) (r error)) *TaskPusherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushAllFunc = f
	return m.mock
}

//SinkPushAll implements github.com/insolar/insolar/conveyor/slot.TaskPusher interface
func (m *TaskPusherMock) SinkPushAll(p []interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushAllPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushAllCounter, 1)

	if len(m.SinkPushAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TaskPusherMock.SinkPushAll. %v", p)
			return
		}

		input := m.SinkPushAllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TaskPusherMockSinkPushAllInput{p}, "TaskPusher.SinkPushAll got unexpected parameters")

		result := m.SinkPushAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TaskPusherMock.SinkPushAll")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushAllMock.mainExpectation != nil {

		input := m.SinkPushAllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TaskPusherMockSinkPushAllInput{p}, "TaskPusher.SinkPushAll got unexpected parameters")
		}

		result := m.SinkPushAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TaskPusherMock.SinkPushAll")
		}

		r = result.r

		return
	}

	if m.SinkPushAllFunc == nil {
		m.t.Fatalf("Unexpected call to TaskPusherMock.SinkPushAll. %v", p)
		return
	}

	return m.SinkPushAllFunc(p)
}

//SinkPushAllMinimockCounter returns a count of TaskPusherMock.SinkPushAllFunc invocations
func (m *TaskPusherMock) SinkPushAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllCounter)
}

//SinkPushAllMinimockPreCounter returns the value of TaskPusherMock.SinkPushAll invocations
func (m *TaskPusherMock) SinkPushAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllPreCounter)
}

//SinkPushAllFinished returns true if mock invocations count is ok
func (m *TaskPusherMock) SinkPushAllFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SinkPushAllMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SinkPushAllCounter) == uint64(len(m.SinkPushAllMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SinkPushAllMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SinkPushAllCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SinkPushAllFunc != nil {
		return atomic.LoadUint64(&m.SinkPushAllCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TaskPusherMock) ValidateCallCounters() {

	if !m.PushSignalFinished() {
		m.t.Fatal("Expected call to TaskPusherMock.PushSignal")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to TaskPusherMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to TaskPusherMock.SinkPushAll")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TaskPusherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TaskPusherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TaskPusherMock) MinimockFinish() {

	if !m.PushSignalFinished() {
		m.t.Fatal("Expected call to TaskPusherMock.PushSignal")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to TaskPusherMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to TaskPusherMock.SinkPushAll")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TaskPusherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TaskPusherMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.PushSignalFinished()
		ok = ok && m.SinkPushFinished()
		ok = ok && m.SinkPushAllFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.PushSignalFinished() {
				m.t.Error("Expected call to TaskPusherMock.PushSignal")
			}

			if !m.SinkPushFinished() {
				m.t.Error("Expected call to TaskPusherMock.SinkPush")
			}

			if !m.SinkPushAllFinished() {
				m.t.Error("Expected call to TaskPusherMock.SinkPushAll")
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
func (m *TaskPusherMock) AllMocksCalled() bool {

	if !m.PushSignalFinished() {
		return false
	}

	if !m.SinkPushFinished() {
		return false
	}

	if !m.SinkPushAllFinished() {
		return false
	}

	return true
}

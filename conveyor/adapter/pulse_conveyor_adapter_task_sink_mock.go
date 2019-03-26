package adapter

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseConveyorAdapterTaskSink" can be found in github.com/insolar/insolar/conveyor/adapter
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseConveyorAdapterTaskSinkMock implements github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink
type PulseConveyorAdapterTaskSinkMock struct {
	t minimock.Tester

	CancelElementTasksFunc       func(p insolar.PulseNumber, p1 uint32)
	CancelElementTasksCounter    uint64
	CancelElementTasksPreCounter uint64
	CancelElementTasksMock       mPulseConveyorAdapterTaskSinkMockCancelElementTasks

	CancelPulseTasksFunc       func(p insolar.PulseNumber)
	CancelPulseTasksCounter    uint64
	CancelPulseTasksPreCounter uint64
	CancelPulseTasksMock       mPulseConveyorAdapterTaskSinkMockCancelPulseTasks

	FlushNodeTasksFunc       func(p uint32)
	FlushNodeTasksCounter    uint64
	FlushNodeTasksPreCounter uint64
	FlushNodeTasksMock       mPulseConveyorAdapterTaskSinkMockFlushNodeTasks

	FlushPulseTasksFunc       func(p insolar.PulseNumber)
	FlushPulseTasksCounter    uint64
	FlushPulseTasksPreCounter uint64
	FlushPulseTasksMock       mPulseConveyorAdapterTaskSinkMockFlushPulseTasks

	GetAdapterIDFunc       func() (r uint32)
	GetAdapterIDCounter    uint64
	GetAdapterIDPreCounter uint64
	GetAdapterIDMock       mPulseConveyorAdapterTaskSinkMockGetAdapterID

	PushTaskFunc       func(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) (r error)
	PushTaskCounter    uint64
	PushTaskPreCounter uint64
	PushTaskMock       mPulseConveyorAdapterTaskSinkMockPushTask
}

//NewPulseConveyorAdapterTaskSinkMock returns a mock for github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink
func NewPulseConveyorAdapterTaskSinkMock(t minimock.Tester) *PulseConveyorAdapterTaskSinkMock {
	m := &PulseConveyorAdapterTaskSinkMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CancelElementTasksMock = mPulseConveyorAdapterTaskSinkMockCancelElementTasks{mock: m}
	m.CancelPulseTasksMock = mPulseConveyorAdapterTaskSinkMockCancelPulseTasks{mock: m}
	m.FlushNodeTasksMock = mPulseConveyorAdapterTaskSinkMockFlushNodeTasks{mock: m}
	m.FlushPulseTasksMock = mPulseConveyorAdapterTaskSinkMockFlushPulseTasks{mock: m}
	m.GetAdapterIDMock = mPulseConveyorAdapterTaskSinkMockGetAdapterID{mock: m}
	m.PushTaskMock = mPulseConveyorAdapterTaskSinkMockPushTask{mock: m}

	return m
}

type mPulseConveyorAdapterTaskSinkMockCancelElementTasks struct {
	mock              *PulseConveyorAdapterTaskSinkMock
	mainExpectation   *PulseConveyorAdapterTaskSinkMockCancelElementTasksExpectation
	expectationSeries []*PulseConveyorAdapterTaskSinkMockCancelElementTasksExpectation
}

type PulseConveyorAdapterTaskSinkMockCancelElementTasksExpectation struct {
	input *PulseConveyorAdapterTaskSinkMockCancelElementTasksInput
}

type PulseConveyorAdapterTaskSinkMockCancelElementTasksInput struct {
	p  insolar.PulseNumber
	p1 uint32
}

//Expect specifies that invocation of PulseConveyorAdapterTaskSink.CancelElementTasks is expected from 1 to Infinity times
func (m *mPulseConveyorAdapterTaskSinkMockCancelElementTasks) Expect(p insolar.PulseNumber, p1 uint32) *mPulseConveyorAdapterTaskSinkMockCancelElementTasks {
	m.mock.CancelElementTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockCancelElementTasksExpectation{}
	}
	m.mainExpectation.input = &PulseConveyorAdapterTaskSinkMockCancelElementTasksInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseConveyorAdapterTaskSink.CancelElementTasks
func (m *mPulseConveyorAdapterTaskSinkMockCancelElementTasks) Return() *PulseConveyorAdapterTaskSinkMock {
	m.mock.CancelElementTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockCancelElementTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PulseConveyorAdapterTaskSink.CancelElementTasks is expected once
func (m *mPulseConveyorAdapterTaskSinkMockCancelElementTasks) ExpectOnce(p insolar.PulseNumber, p1 uint32) *PulseConveyorAdapterTaskSinkMockCancelElementTasksExpectation {
	m.mock.CancelElementTasksFunc = nil
	m.mainExpectation = nil

	expectation := &PulseConveyorAdapterTaskSinkMockCancelElementTasksExpectation{}
	expectation.input = &PulseConveyorAdapterTaskSinkMockCancelElementTasksInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PulseConveyorAdapterTaskSink.CancelElementTasks method
func (m *mPulseConveyorAdapterTaskSinkMockCancelElementTasks) Set(f func(p insolar.PulseNumber, p1 uint32)) *PulseConveyorAdapterTaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CancelElementTasksFunc = f
	return m.mock
}

//CancelElementTasks implements github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink interface
func (m *PulseConveyorAdapterTaskSinkMock) CancelElementTasks(p insolar.PulseNumber, p1 uint32) {
	counter := atomic.AddUint64(&m.CancelElementTasksPreCounter, 1)
	defer atomic.AddUint64(&m.CancelElementTasksCounter, 1)

	if len(m.CancelElementTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CancelElementTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.CancelElementTasks. %v %v", p, p1)
			return
		}

		input := m.CancelElementTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockCancelElementTasksInput{p, p1}, "PulseConveyorAdapterTaskSink.CancelElementTasks got unexpected parameters")

		return
	}

	if m.CancelElementTasksMock.mainExpectation != nil {

		input := m.CancelElementTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockCancelElementTasksInput{p, p1}, "PulseConveyorAdapterTaskSink.CancelElementTasks got unexpected parameters")
		}

		return
	}

	if m.CancelElementTasksFunc == nil {
		m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.CancelElementTasks. %v %v", p, p1)
		return
	}

	m.CancelElementTasksFunc(p, p1)
}

//CancelElementTasksMinimockCounter returns a count of PulseConveyorAdapterTaskSinkMock.CancelElementTasksFunc invocations
func (m *PulseConveyorAdapterTaskSinkMock) CancelElementTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CancelElementTasksCounter)
}

//CancelElementTasksMinimockPreCounter returns the value of PulseConveyorAdapterTaskSinkMock.CancelElementTasks invocations
func (m *PulseConveyorAdapterTaskSinkMock) CancelElementTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CancelElementTasksPreCounter)
}

//CancelElementTasksFinished returns true if mock invocations count is ok
func (m *PulseConveyorAdapterTaskSinkMock) CancelElementTasksFinished() bool {
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

type mPulseConveyorAdapterTaskSinkMockCancelPulseTasks struct {
	mock              *PulseConveyorAdapterTaskSinkMock
	mainExpectation   *PulseConveyorAdapterTaskSinkMockCancelPulseTasksExpectation
	expectationSeries []*PulseConveyorAdapterTaskSinkMockCancelPulseTasksExpectation
}

type PulseConveyorAdapterTaskSinkMockCancelPulseTasksExpectation struct {
	input *PulseConveyorAdapterTaskSinkMockCancelPulseTasksInput
}

type PulseConveyorAdapterTaskSinkMockCancelPulseTasksInput struct {
	p insolar.PulseNumber
}

//Expect specifies that invocation of PulseConveyorAdapterTaskSink.CancelPulseTasks is expected from 1 to Infinity times
func (m *mPulseConveyorAdapterTaskSinkMockCancelPulseTasks) Expect(p insolar.PulseNumber) *mPulseConveyorAdapterTaskSinkMockCancelPulseTasks {
	m.mock.CancelPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockCancelPulseTasksExpectation{}
	}
	m.mainExpectation.input = &PulseConveyorAdapterTaskSinkMockCancelPulseTasksInput{p}
	return m
}

//Return specifies results of invocation of PulseConveyorAdapterTaskSink.CancelPulseTasks
func (m *mPulseConveyorAdapterTaskSinkMockCancelPulseTasks) Return() *PulseConveyorAdapterTaskSinkMock {
	m.mock.CancelPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockCancelPulseTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PulseConveyorAdapterTaskSink.CancelPulseTasks is expected once
func (m *mPulseConveyorAdapterTaskSinkMockCancelPulseTasks) ExpectOnce(p insolar.PulseNumber) *PulseConveyorAdapterTaskSinkMockCancelPulseTasksExpectation {
	m.mock.CancelPulseTasksFunc = nil
	m.mainExpectation = nil

	expectation := &PulseConveyorAdapterTaskSinkMockCancelPulseTasksExpectation{}
	expectation.input = &PulseConveyorAdapterTaskSinkMockCancelPulseTasksInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PulseConveyorAdapterTaskSink.CancelPulseTasks method
func (m *mPulseConveyorAdapterTaskSinkMockCancelPulseTasks) Set(f func(p insolar.PulseNumber)) *PulseConveyorAdapterTaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CancelPulseTasksFunc = f
	return m.mock
}

//CancelPulseTasks implements github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink interface
func (m *PulseConveyorAdapterTaskSinkMock) CancelPulseTasks(p insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.CancelPulseTasksPreCounter, 1)
	defer atomic.AddUint64(&m.CancelPulseTasksCounter, 1)

	if len(m.CancelPulseTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CancelPulseTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.CancelPulseTasks. %v", p)
			return
		}

		input := m.CancelPulseTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockCancelPulseTasksInput{p}, "PulseConveyorAdapterTaskSink.CancelPulseTasks got unexpected parameters")

		return
	}

	if m.CancelPulseTasksMock.mainExpectation != nil {

		input := m.CancelPulseTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockCancelPulseTasksInput{p}, "PulseConveyorAdapterTaskSink.CancelPulseTasks got unexpected parameters")
		}

		return
	}

	if m.CancelPulseTasksFunc == nil {
		m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.CancelPulseTasks. %v", p)
		return
	}

	m.CancelPulseTasksFunc(p)
}

//CancelPulseTasksMinimockCounter returns a count of PulseConveyorAdapterTaskSinkMock.CancelPulseTasksFunc invocations
func (m *PulseConveyorAdapterTaskSinkMock) CancelPulseTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CancelPulseTasksCounter)
}

//CancelPulseTasksMinimockPreCounter returns the value of PulseConveyorAdapterTaskSinkMock.CancelPulseTasks invocations
func (m *PulseConveyorAdapterTaskSinkMock) CancelPulseTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CancelPulseTasksPreCounter)
}

//CancelPulseTasksFinished returns true if mock invocations count is ok
func (m *PulseConveyorAdapterTaskSinkMock) CancelPulseTasksFinished() bool {
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

type mPulseConveyorAdapterTaskSinkMockFlushNodeTasks struct {
	mock              *PulseConveyorAdapterTaskSinkMock
	mainExpectation   *PulseConveyorAdapterTaskSinkMockFlushNodeTasksExpectation
	expectationSeries []*PulseConveyorAdapterTaskSinkMockFlushNodeTasksExpectation
}

type PulseConveyorAdapterTaskSinkMockFlushNodeTasksExpectation struct {
	input *PulseConveyorAdapterTaskSinkMockFlushNodeTasksInput
}

type PulseConveyorAdapterTaskSinkMockFlushNodeTasksInput struct {
	p uint32
}

//Expect specifies that invocation of PulseConveyorAdapterTaskSink.FlushNodeTasks is expected from 1 to Infinity times
func (m *mPulseConveyorAdapterTaskSinkMockFlushNodeTasks) Expect(p uint32) *mPulseConveyorAdapterTaskSinkMockFlushNodeTasks {
	m.mock.FlushNodeTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockFlushNodeTasksExpectation{}
	}
	m.mainExpectation.input = &PulseConveyorAdapterTaskSinkMockFlushNodeTasksInput{p}
	return m
}

//Return specifies results of invocation of PulseConveyorAdapterTaskSink.FlushNodeTasks
func (m *mPulseConveyorAdapterTaskSinkMockFlushNodeTasks) Return() *PulseConveyorAdapterTaskSinkMock {
	m.mock.FlushNodeTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockFlushNodeTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PulseConveyorAdapterTaskSink.FlushNodeTasks is expected once
func (m *mPulseConveyorAdapterTaskSinkMockFlushNodeTasks) ExpectOnce(p uint32) *PulseConveyorAdapterTaskSinkMockFlushNodeTasksExpectation {
	m.mock.FlushNodeTasksFunc = nil
	m.mainExpectation = nil

	expectation := &PulseConveyorAdapterTaskSinkMockFlushNodeTasksExpectation{}
	expectation.input = &PulseConveyorAdapterTaskSinkMockFlushNodeTasksInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PulseConveyorAdapterTaskSink.FlushNodeTasks method
func (m *mPulseConveyorAdapterTaskSinkMockFlushNodeTasks) Set(f func(p uint32)) *PulseConveyorAdapterTaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FlushNodeTasksFunc = f
	return m.mock
}

//FlushNodeTasks implements github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink interface
func (m *PulseConveyorAdapterTaskSinkMock) FlushNodeTasks(p uint32) {
	counter := atomic.AddUint64(&m.FlushNodeTasksPreCounter, 1)
	defer atomic.AddUint64(&m.FlushNodeTasksCounter, 1)

	if len(m.FlushNodeTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FlushNodeTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.FlushNodeTasks. %v", p)
			return
		}

		input := m.FlushNodeTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockFlushNodeTasksInput{p}, "PulseConveyorAdapterTaskSink.FlushNodeTasks got unexpected parameters")

		return
	}

	if m.FlushNodeTasksMock.mainExpectation != nil {

		input := m.FlushNodeTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockFlushNodeTasksInput{p}, "PulseConveyorAdapterTaskSink.FlushNodeTasks got unexpected parameters")
		}

		return
	}

	if m.FlushNodeTasksFunc == nil {
		m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.FlushNodeTasks. %v", p)
		return
	}

	m.FlushNodeTasksFunc(p)
}

//FlushNodeTasksMinimockCounter returns a count of PulseConveyorAdapterTaskSinkMock.FlushNodeTasksFunc invocations
func (m *PulseConveyorAdapterTaskSinkMock) FlushNodeTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FlushNodeTasksCounter)
}

//FlushNodeTasksMinimockPreCounter returns the value of PulseConveyorAdapterTaskSinkMock.FlushNodeTasks invocations
func (m *PulseConveyorAdapterTaskSinkMock) FlushNodeTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FlushNodeTasksPreCounter)
}

//FlushNodeTasksFinished returns true if mock invocations count is ok
func (m *PulseConveyorAdapterTaskSinkMock) FlushNodeTasksFinished() bool {
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

type mPulseConveyorAdapterTaskSinkMockFlushPulseTasks struct {
	mock              *PulseConveyorAdapterTaskSinkMock
	mainExpectation   *PulseConveyorAdapterTaskSinkMockFlushPulseTasksExpectation
	expectationSeries []*PulseConveyorAdapterTaskSinkMockFlushPulseTasksExpectation
}

type PulseConveyorAdapterTaskSinkMockFlushPulseTasksExpectation struct {
	input *PulseConveyorAdapterTaskSinkMockFlushPulseTasksInput
}

type PulseConveyorAdapterTaskSinkMockFlushPulseTasksInput struct {
	p insolar.PulseNumber
}

//Expect specifies that invocation of PulseConveyorAdapterTaskSink.FlushPulseTasks is expected from 1 to Infinity times
func (m *mPulseConveyorAdapterTaskSinkMockFlushPulseTasks) Expect(p insolar.PulseNumber) *mPulseConveyorAdapterTaskSinkMockFlushPulseTasks {
	m.mock.FlushPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockFlushPulseTasksExpectation{}
	}
	m.mainExpectation.input = &PulseConveyorAdapterTaskSinkMockFlushPulseTasksInput{p}
	return m
}

//Return specifies results of invocation of PulseConveyorAdapterTaskSink.FlushPulseTasks
func (m *mPulseConveyorAdapterTaskSinkMockFlushPulseTasks) Return() *PulseConveyorAdapterTaskSinkMock {
	m.mock.FlushPulseTasksFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockFlushPulseTasksExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of PulseConveyorAdapterTaskSink.FlushPulseTasks is expected once
func (m *mPulseConveyorAdapterTaskSinkMockFlushPulseTasks) ExpectOnce(p insolar.PulseNumber) *PulseConveyorAdapterTaskSinkMockFlushPulseTasksExpectation {
	m.mock.FlushPulseTasksFunc = nil
	m.mainExpectation = nil

	expectation := &PulseConveyorAdapterTaskSinkMockFlushPulseTasksExpectation{}
	expectation.input = &PulseConveyorAdapterTaskSinkMockFlushPulseTasksInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of PulseConveyorAdapterTaskSink.FlushPulseTasks method
func (m *mPulseConveyorAdapterTaskSinkMockFlushPulseTasks) Set(f func(p insolar.PulseNumber)) *PulseConveyorAdapterTaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FlushPulseTasksFunc = f
	return m.mock
}

//FlushPulseTasks implements github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink interface
func (m *PulseConveyorAdapterTaskSinkMock) FlushPulseTasks(p insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.FlushPulseTasksPreCounter, 1)
	defer atomic.AddUint64(&m.FlushPulseTasksCounter, 1)

	if len(m.FlushPulseTasksMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FlushPulseTasksMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.FlushPulseTasks. %v", p)
			return
		}

		input := m.FlushPulseTasksMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockFlushPulseTasksInput{p}, "PulseConveyorAdapterTaskSink.FlushPulseTasks got unexpected parameters")

		return
	}

	if m.FlushPulseTasksMock.mainExpectation != nil {

		input := m.FlushPulseTasksMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockFlushPulseTasksInput{p}, "PulseConveyorAdapterTaskSink.FlushPulseTasks got unexpected parameters")
		}

		return
	}

	if m.FlushPulseTasksFunc == nil {
		m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.FlushPulseTasks. %v", p)
		return
	}

	m.FlushPulseTasksFunc(p)
}

//FlushPulseTasksMinimockCounter returns a count of PulseConveyorAdapterTaskSinkMock.FlushPulseTasksFunc invocations
func (m *PulseConveyorAdapterTaskSinkMock) FlushPulseTasksMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FlushPulseTasksCounter)
}

//FlushPulseTasksMinimockPreCounter returns the value of PulseConveyorAdapterTaskSinkMock.FlushPulseTasks invocations
func (m *PulseConveyorAdapterTaskSinkMock) FlushPulseTasksMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FlushPulseTasksPreCounter)
}

//FlushPulseTasksFinished returns true if mock invocations count is ok
func (m *PulseConveyorAdapterTaskSinkMock) FlushPulseTasksFinished() bool {
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

type mPulseConveyorAdapterTaskSinkMockGetAdapterID struct {
	mock              *PulseConveyorAdapterTaskSinkMock
	mainExpectation   *PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation
	expectationSeries []*PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation
}

type PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation struct {
	result *PulseConveyorAdapterTaskSinkMockGetAdapterIDResult
}

type PulseConveyorAdapterTaskSinkMockGetAdapterIDResult struct {
	r uint32
}

//Expect specifies that invocation of PulseConveyorAdapterTaskSink.GetAdapterID is expected from 1 to Infinity times
func (m *mPulseConveyorAdapterTaskSinkMockGetAdapterID) Expect() *mPulseConveyorAdapterTaskSinkMockGetAdapterID {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of PulseConveyorAdapterTaskSink.GetAdapterID
func (m *mPulseConveyorAdapterTaskSinkMockGetAdapterID) Return(r uint32) *PulseConveyorAdapterTaskSinkMock {
	m.mock.GetAdapterIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation{}
	}
	m.mainExpectation.result = &PulseConveyorAdapterTaskSinkMockGetAdapterIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseConveyorAdapterTaskSink.GetAdapterID is expected once
func (m *mPulseConveyorAdapterTaskSinkMockGetAdapterID) ExpectOnce() *PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation {
	m.mock.GetAdapterIDFunc = nil
	m.mainExpectation = nil

	expectation := &PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseConveyorAdapterTaskSinkMockGetAdapterIDExpectation) Return(r uint32) {
	e.result = &PulseConveyorAdapterTaskSinkMockGetAdapterIDResult{r}
}

//Set uses given function f as a mock of PulseConveyorAdapterTaskSink.GetAdapterID method
func (m *mPulseConveyorAdapterTaskSinkMockGetAdapterID) Set(f func() (r uint32)) *PulseConveyorAdapterTaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAdapterIDFunc = f
	return m.mock
}

//GetAdapterID implements github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink interface
func (m *PulseConveyorAdapterTaskSinkMock) GetAdapterID() (r uint32) {
	counter := atomic.AddUint64(&m.GetAdapterIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetAdapterIDCounter, 1)

	if len(m.GetAdapterIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAdapterIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.GetAdapterID.")
			return
		}

		result := m.GetAdapterIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseConveyorAdapterTaskSinkMock.GetAdapterID")
			return
		}

		r = result.r

		return
	}

	if m.GetAdapterIDMock.mainExpectation != nil {

		result := m.GetAdapterIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseConveyorAdapterTaskSinkMock.GetAdapterID")
		}

		r = result.r

		return
	}

	if m.GetAdapterIDFunc == nil {
		m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.GetAdapterID.")
		return
	}

	return m.GetAdapterIDFunc()
}

//GetAdapterIDMinimockCounter returns a count of PulseConveyorAdapterTaskSinkMock.GetAdapterIDFunc invocations
func (m *PulseConveyorAdapterTaskSinkMock) GetAdapterIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDCounter)
}

//GetAdapterIDMinimockPreCounter returns the value of PulseConveyorAdapterTaskSinkMock.GetAdapterID invocations
func (m *PulseConveyorAdapterTaskSinkMock) GetAdapterIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAdapterIDPreCounter)
}

//GetAdapterIDFinished returns true if mock invocations count is ok
func (m *PulseConveyorAdapterTaskSinkMock) GetAdapterIDFinished() bool {
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

type mPulseConveyorAdapterTaskSinkMockPushTask struct {
	mock              *PulseConveyorAdapterTaskSinkMock
	mainExpectation   *PulseConveyorAdapterTaskSinkMockPushTaskExpectation
	expectationSeries []*PulseConveyorAdapterTaskSinkMockPushTaskExpectation
}

type PulseConveyorAdapterTaskSinkMockPushTaskExpectation struct {
	input  *PulseConveyorAdapterTaskSinkMockPushTaskInput
	result *PulseConveyorAdapterTaskSinkMockPushTaskResult
}

type PulseConveyorAdapterTaskSinkMockPushTaskInput struct {
	p  AdapterToSlotResponseSink
	p1 uint32
	p2 uint32
	p3 interface{}
}

type PulseConveyorAdapterTaskSinkMockPushTaskResult struct {
	r error
}

//Expect specifies that invocation of PulseConveyorAdapterTaskSink.PushTask is expected from 1 to Infinity times
func (m *mPulseConveyorAdapterTaskSinkMockPushTask) Expect(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) *mPulseConveyorAdapterTaskSinkMockPushTask {
	m.mock.PushTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockPushTaskExpectation{}
	}
	m.mainExpectation.input = &PulseConveyorAdapterTaskSinkMockPushTaskInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of PulseConveyorAdapterTaskSink.PushTask
func (m *mPulseConveyorAdapterTaskSinkMockPushTask) Return(r error) *PulseConveyorAdapterTaskSinkMock {
	m.mock.PushTaskFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseConveyorAdapterTaskSinkMockPushTaskExpectation{}
	}
	m.mainExpectation.result = &PulseConveyorAdapterTaskSinkMockPushTaskResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseConveyorAdapterTaskSink.PushTask is expected once
func (m *mPulseConveyorAdapterTaskSinkMockPushTask) ExpectOnce(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) *PulseConveyorAdapterTaskSinkMockPushTaskExpectation {
	m.mock.PushTaskFunc = nil
	m.mainExpectation = nil

	expectation := &PulseConveyorAdapterTaskSinkMockPushTaskExpectation{}
	expectation.input = &PulseConveyorAdapterTaskSinkMockPushTaskInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseConveyorAdapterTaskSinkMockPushTaskExpectation) Return(r error) {
	e.result = &PulseConveyorAdapterTaskSinkMockPushTaskResult{r}
}

//Set uses given function f as a mock of PulseConveyorAdapterTaskSink.PushTask method
func (m *mPulseConveyorAdapterTaskSinkMockPushTask) Set(f func(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) (r error)) *PulseConveyorAdapterTaskSinkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PushTaskFunc = f
	return m.mock
}

//PushTask implements github.com/insolar/insolar/conveyor/adapter.PulseConveyorAdapterTaskSink interface
func (m *PulseConveyorAdapterTaskSinkMock) PushTask(p AdapterToSlotResponseSink, p1 uint32, p2 uint32, p3 interface{}) (r error) {
	counter := atomic.AddUint64(&m.PushTaskPreCounter, 1)
	defer atomic.AddUint64(&m.PushTaskCounter, 1)

	if len(m.PushTaskMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PushTaskMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.PushTask. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.PushTaskMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockPushTaskInput{p, p1, p2, p3}, "PulseConveyorAdapterTaskSink.PushTask got unexpected parameters")

		result := m.PushTaskMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseConveyorAdapterTaskSinkMock.PushTask")
			return
		}

		r = result.r

		return
	}

	if m.PushTaskMock.mainExpectation != nil {

		input := m.PushTaskMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseConveyorAdapterTaskSinkMockPushTaskInput{p, p1, p2, p3}, "PulseConveyorAdapterTaskSink.PushTask got unexpected parameters")
		}

		result := m.PushTaskMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseConveyorAdapterTaskSinkMock.PushTask")
		}

		r = result.r

		return
	}

	if m.PushTaskFunc == nil {
		m.t.Fatalf("Unexpected call to PulseConveyorAdapterTaskSinkMock.PushTask. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.PushTaskFunc(p, p1, p2, p3)
}

//PushTaskMinimockCounter returns a count of PulseConveyorAdapterTaskSinkMock.PushTaskFunc invocations
func (m *PulseConveyorAdapterTaskSinkMock) PushTaskMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PushTaskCounter)
}

//PushTaskMinimockPreCounter returns the value of PulseConveyorAdapterTaskSinkMock.PushTask invocations
func (m *PulseConveyorAdapterTaskSinkMock) PushTaskMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PushTaskPreCounter)
}

//PushTaskFinished returns true if mock invocations count is ok
func (m *PulseConveyorAdapterTaskSinkMock) PushTaskFinished() bool {
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
func (m *PulseConveyorAdapterTaskSinkMock) ValidateCallCounters() {

	if !m.CancelElementTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.CancelElementTasks")
	}

	if !m.CancelPulseTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.CancelPulseTasks")
	}

	if !m.FlushNodeTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.FlushNodeTasks")
	}

	if !m.FlushPulseTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.FlushPulseTasks")
	}

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.GetAdapterID")
	}

	if !m.PushTaskFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.PushTask")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseConveyorAdapterTaskSinkMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseConveyorAdapterTaskSinkMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseConveyorAdapterTaskSinkMock) MinimockFinish() {

	if !m.CancelElementTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.CancelElementTasks")
	}

	if !m.CancelPulseTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.CancelPulseTasks")
	}

	if !m.FlushNodeTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.FlushNodeTasks")
	}

	if !m.FlushPulseTasksFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.FlushPulseTasks")
	}

	if !m.GetAdapterIDFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.GetAdapterID")
	}

	if !m.PushTaskFinished() {
		m.t.Fatal("Expected call to PulseConveyorAdapterTaskSinkMock.PushTask")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseConveyorAdapterTaskSinkMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseConveyorAdapterTaskSinkMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to PulseConveyorAdapterTaskSinkMock.CancelElementTasks")
			}

			if !m.CancelPulseTasksFinished() {
				m.t.Error("Expected call to PulseConveyorAdapterTaskSinkMock.CancelPulseTasks")
			}

			if !m.FlushNodeTasksFinished() {
				m.t.Error("Expected call to PulseConveyorAdapterTaskSinkMock.FlushNodeTasks")
			}

			if !m.FlushPulseTasksFinished() {
				m.t.Error("Expected call to PulseConveyorAdapterTaskSinkMock.FlushPulseTasks")
			}

			if !m.GetAdapterIDFinished() {
				m.t.Error("Expected call to PulseConveyorAdapterTaskSinkMock.GetAdapterID")
			}

			if !m.PushTaskFinished() {
				m.t.Error("Expected call to PulseConveyorAdapterTaskSinkMock.PushTask")
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
func (m *PulseConveyorAdapterTaskSinkMock) AllMocksCalled() bool {

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

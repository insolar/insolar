package queue

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Queue" can be found in github.com/insolar/insolar/conveyor/queue
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//QueueMock implements github.com/insolar/insolar/conveyor/queue.Queue
type QueueMock struct {
	t minimock.Tester

	BlockAndRemoveAllFunc       func() (r []OutputElement)
	BlockAndRemoveAllCounter    uint64
	BlockAndRemoveAllPreCounter uint64
	BlockAndRemoveAllMock       mQueueMockBlockAndRemoveAll

	HasSignalFunc       func() (r bool)
	HasSignalCounter    uint64
	HasSignalPreCounter uint64
	HasSignalMock       mQueueMockHasSignal

	PushSignalFunc       func(p uint32, p1 SyncDone) (r error)
	PushSignalCounter    uint64
	PushSignalPreCounter uint64
	PushSignalMock       mQueueMockPushSignal

	RemoveAllFunc       func() (r []OutputElement)
	RemoveAllCounter    uint64
	RemoveAllPreCounter uint64
	RemoveAllMock       mQueueMockRemoveAll

	SinkPushFunc       func(p interface{}) (r error)
	SinkPushCounter    uint64
	SinkPushPreCounter uint64
	SinkPushMock       mQueueMockSinkPush

	SinkPushAllFunc       func(p []interface{}) (r error)
	SinkPushAllCounter    uint64
	SinkPushAllPreCounter uint64
	SinkPushAllMock       mQueueMockSinkPushAll

	UnblockFunc       func() (r bool)
	UnblockCounter    uint64
	UnblockPreCounter uint64
	UnblockMock       mQueueMockUnblock
}

//NewQueueMock returns a mock for github.com/insolar/insolar/conveyor/queue.Queue
func NewQueueMock(t minimock.Tester) *QueueMock {
	m := &QueueMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BlockAndRemoveAllMock = mQueueMockBlockAndRemoveAll{mock: m}
	m.HasSignalMock = mQueueMockHasSignal{mock: m}
	m.PushSignalMock = mQueueMockPushSignal{mock: m}
	m.RemoveAllMock = mQueueMockRemoveAll{mock: m}
	m.SinkPushMock = mQueueMockSinkPush{mock: m}
	m.SinkPushAllMock = mQueueMockSinkPushAll{mock: m}
	m.UnblockMock = mQueueMockUnblock{mock: m}

	return m
}

type mQueueMockBlockAndRemoveAll struct {
	mock              *QueueMock
	mainExpectation   *QueueMockBlockAndRemoveAllExpectation
	expectationSeries []*QueueMockBlockAndRemoveAllExpectation
}

type QueueMockBlockAndRemoveAllExpectation struct {
	result *QueueMockBlockAndRemoveAllResult
}

type QueueMockBlockAndRemoveAllResult struct {
	r []OutputElement
}

//Expect specifies that invocation of Queue.BlockAndRemoveAll is expected from 1 to Infinity times
func (m *mQueueMockBlockAndRemoveAll) Expect() *mQueueMockBlockAndRemoveAll {
	m.mock.BlockAndRemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockBlockAndRemoveAllExpectation{}
	}

	return m
}

//Return specifies results of invocation of Queue.BlockAndRemoveAll
func (m *mQueueMockBlockAndRemoveAll) Return(r []OutputElement) *QueueMock {
	m.mock.BlockAndRemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockBlockAndRemoveAllExpectation{}
	}
	m.mainExpectation.result = &QueueMockBlockAndRemoveAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Queue.BlockAndRemoveAll is expected once
func (m *mQueueMockBlockAndRemoveAll) ExpectOnce() *QueueMockBlockAndRemoveAllExpectation {
	m.mock.BlockAndRemoveAllFunc = nil
	m.mainExpectation = nil

	expectation := &QueueMockBlockAndRemoveAllExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *QueueMockBlockAndRemoveAllExpectation) Return(r []OutputElement) {
	e.result = &QueueMockBlockAndRemoveAllResult{r}
}

//Set uses given function f as a mock of Queue.BlockAndRemoveAll method
func (m *mQueueMockBlockAndRemoveAll) Set(f func() (r []OutputElement)) *QueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BlockAndRemoveAllFunc = f
	return m.mock
}

//BlockAndRemoveAll implements github.com/insolar/insolar/conveyor/queue.Queue interface
func (m *QueueMock) BlockAndRemoveAll() (r []OutputElement) {
	counter := atomic.AddUint64(&m.BlockAndRemoveAllPreCounter, 1)
	defer atomic.AddUint64(&m.BlockAndRemoveAllCounter, 1)

	if len(m.BlockAndRemoveAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BlockAndRemoveAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to QueueMock.BlockAndRemoveAll.")
			return
		}

		result := m.BlockAndRemoveAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.BlockAndRemoveAll")
			return
		}

		r = result.r

		return
	}

	if m.BlockAndRemoveAllMock.mainExpectation != nil {

		result := m.BlockAndRemoveAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.BlockAndRemoveAll")
		}

		r = result.r

		return
	}

	if m.BlockAndRemoveAllFunc == nil {
		m.t.Fatalf("Unexpected call to QueueMock.BlockAndRemoveAll.")
		return
	}

	return m.BlockAndRemoveAllFunc()
}

//BlockAndRemoveAllMinimockCounter returns a count of QueueMock.BlockAndRemoveAllFunc invocations
func (m *QueueMock) BlockAndRemoveAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BlockAndRemoveAllCounter)
}

//BlockAndRemoveAllMinimockPreCounter returns the value of QueueMock.BlockAndRemoveAll invocations
func (m *QueueMock) BlockAndRemoveAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BlockAndRemoveAllPreCounter)
}

//BlockAndRemoveAllFinished returns true if mock invocations count is ok
func (m *QueueMock) BlockAndRemoveAllFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.BlockAndRemoveAllMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.BlockAndRemoveAllCounter) == uint64(len(m.BlockAndRemoveAllMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.BlockAndRemoveAllMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.BlockAndRemoveAllCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.BlockAndRemoveAllFunc != nil {
		return atomic.LoadUint64(&m.BlockAndRemoveAllCounter) > 0
	}

	return true
}

type mQueueMockHasSignal struct {
	mock              *QueueMock
	mainExpectation   *QueueMockHasSignalExpectation
	expectationSeries []*QueueMockHasSignalExpectation
}

type QueueMockHasSignalExpectation struct {
	result *QueueMockHasSignalResult
}

type QueueMockHasSignalResult struct {
	r bool
}

//Expect specifies that invocation of Queue.HasSignal is expected from 1 to Infinity times
func (m *mQueueMockHasSignal) Expect() *mQueueMockHasSignal {
	m.mock.HasSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockHasSignalExpectation{}
	}

	return m
}

//Return specifies results of invocation of Queue.HasSignal
func (m *mQueueMockHasSignal) Return(r bool) *QueueMock {
	m.mock.HasSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockHasSignalExpectation{}
	}
	m.mainExpectation.result = &QueueMockHasSignalResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Queue.HasSignal is expected once
func (m *mQueueMockHasSignal) ExpectOnce() *QueueMockHasSignalExpectation {
	m.mock.HasSignalFunc = nil
	m.mainExpectation = nil

	expectation := &QueueMockHasSignalExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *QueueMockHasSignalExpectation) Return(r bool) {
	e.result = &QueueMockHasSignalResult{r}
}

//Set uses given function f as a mock of Queue.HasSignal method
func (m *mQueueMockHasSignal) Set(f func() (r bool)) *QueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasSignalFunc = f
	return m.mock
}

//HasSignal implements github.com/insolar/insolar/conveyor/queue.Queue interface
func (m *QueueMock) HasSignal() (r bool) {
	counter := atomic.AddUint64(&m.HasSignalPreCounter, 1)
	defer atomic.AddUint64(&m.HasSignalCounter, 1)

	if len(m.HasSignalMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasSignalMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to QueueMock.HasSignal.")
			return
		}

		result := m.HasSignalMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.HasSignal")
			return
		}

		r = result.r

		return
	}

	if m.HasSignalMock.mainExpectation != nil {

		result := m.HasSignalMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.HasSignal")
		}

		r = result.r

		return
	}

	if m.HasSignalFunc == nil {
		m.t.Fatalf("Unexpected call to QueueMock.HasSignal.")
		return
	}

	return m.HasSignalFunc()
}

//HasSignalMinimockCounter returns a count of QueueMock.HasSignalFunc invocations
func (m *QueueMock) HasSignalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasSignalCounter)
}

//HasSignalMinimockPreCounter returns the value of QueueMock.HasSignal invocations
func (m *QueueMock) HasSignalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasSignalPreCounter)
}

//HasSignalFinished returns true if mock invocations count is ok
func (m *QueueMock) HasSignalFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HasSignalMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HasSignalCounter) == uint64(len(m.HasSignalMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HasSignalMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HasSignalCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HasSignalFunc != nil {
		return atomic.LoadUint64(&m.HasSignalCounter) > 0
	}

	return true
}

type mQueueMockPushSignal struct {
	mock              *QueueMock
	mainExpectation   *QueueMockPushSignalExpectation
	expectationSeries []*QueueMockPushSignalExpectation
}

type QueueMockPushSignalExpectation struct {
	input  *QueueMockPushSignalInput
	result *QueueMockPushSignalResult
}

type QueueMockPushSignalInput struct {
	p  uint32
	p1 SyncDone
}

type QueueMockPushSignalResult struct {
	r error
}

//Expect specifies that invocation of Queue.PushSignal is expected from 1 to Infinity times
func (m *mQueueMockPushSignal) Expect(p uint32, p1 SyncDone) *mQueueMockPushSignal {
	m.mock.PushSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockPushSignalExpectation{}
	}
	m.mainExpectation.input = &QueueMockPushSignalInput{p, p1}
	return m
}

//Return specifies results of invocation of Queue.PushSignal
func (m *mQueueMockPushSignal) Return(r error) *QueueMock {
	m.mock.PushSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockPushSignalExpectation{}
	}
	m.mainExpectation.result = &QueueMockPushSignalResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Queue.PushSignal is expected once
func (m *mQueueMockPushSignal) ExpectOnce(p uint32, p1 SyncDone) *QueueMockPushSignalExpectation {
	m.mock.PushSignalFunc = nil
	m.mainExpectation = nil

	expectation := &QueueMockPushSignalExpectation{}
	expectation.input = &QueueMockPushSignalInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *QueueMockPushSignalExpectation) Return(r error) {
	e.result = &QueueMockPushSignalResult{r}
}

//Set uses given function f as a mock of Queue.PushSignal method
func (m *mQueueMockPushSignal) Set(f func(p uint32, p1 SyncDone) (r error)) *QueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PushSignalFunc = f
	return m.mock
}

//PushSignal implements github.com/insolar/insolar/conveyor/queue.Queue interface
func (m *QueueMock) PushSignal(p uint32, p1 SyncDone) (r error) {
	counter := atomic.AddUint64(&m.PushSignalPreCounter, 1)
	defer atomic.AddUint64(&m.PushSignalCounter, 1)

	if len(m.PushSignalMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PushSignalMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to QueueMock.PushSignal. %v %v", p, p1)
			return
		}

		input := m.PushSignalMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, QueueMockPushSignalInput{p, p1}, "Queue.PushSignal got unexpected parameters")

		result := m.PushSignalMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.PushSignal")
			return
		}

		r = result.r

		return
	}

	if m.PushSignalMock.mainExpectation != nil {

		input := m.PushSignalMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, QueueMockPushSignalInput{p, p1}, "Queue.PushSignal got unexpected parameters")
		}

		result := m.PushSignalMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.PushSignal")
		}

		r = result.r

		return
	}

	if m.PushSignalFunc == nil {
		m.t.Fatalf("Unexpected call to QueueMock.PushSignal. %v %v", p, p1)
		return
	}

	return m.PushSignalFunc(p, p1)
}

//PushSignalMinimockCounter returns a count of QueueMock.PushSignalFunc invocations
func (m *QueueMock) PushSignalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PushSignalCounter)
}

//PushSignalMinimockPreCounter returns the value of QueueMock.PushSignal invocations
func (m *QueueMock) PushSignalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PushSignalPreCounter)
}

//PushSignalFinished returns true if mock invocations count is ok
func (m *QueueMock) PushSignalFinished() bool {
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

type mQueueMockRemoveAll struct {
	mock              *QueueMock
	mainExpectation   *QueueMockRemoveAllExpectation
	expectationSeries []*QueueMockRemoveAllExpectation
}

type QueueMockRemoveAllExpectation struct {
	result *QueueMockRemoveAllResult
}

type QueueMockRemoveAllResult struct {
	r []OutputElement
}

//Expect specifies that invocation of Queue.RemoveAll is expected from 1 to Infinity times
func (m *mQueueMockRemoveAll) Expect() *mQueueMockRemoveAll {
	m.mock.RemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockRemoveAllExpectation{}
	}

	return m
}

//Return specifies results of invocation of Queue.RemoveAll
func (m *mQueueMockRemoveAll) Return(r []OutputElement) *QueueMock {
	m.mock.RemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockRemoveAllExpectation{}
	}
	m.mainExpectation.result = &QueueMockRemoveAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Queue.RemoveAll is expected once
func (m *mQueueMockRemoveAll) ExpectOnce() *QueueMockRemoveAllExpectation {
	m.mock.RemoveAllFunc = nil
	m.mainExpectation = nil

	expectation := &QueueMockRemoveAllExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *QueueMockRemoveAllExpectation) Return(r []OutputElement) {
	e.result = &QueueMockRemoveAllResult{r}
}

//Set uses given function f as a mock of Queue.RemoveAll method
func (m *mQueueMockRemoveAll) Set(f func() (r []OutputElement)) *QueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveAllFunc = f
	return m.mock
}

//RemoveAll implements github.com/insolar/insolar/conveyor/queue.Queue interface
func (m *QueueMock) RemoveAll() (r []OutputElement) {
	counter := atomic.AddUint64(&m.RemoveAllPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveAllCounter, 1)

	if len(m.RemoveAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to QueueMock.RemoveAll.")
			return
		}

		result := m.RemoveAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.RemoveAll")
			return
		}

		r = result.r

		return
	}

	if m.RemoveAllMock.mainExpectation != nil {

		result := m.RemoveAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.RemoveAll")
		}

		r = result.r

		return
	}

	if m.RemoveAllFunc == nil {
		m.t.Fatalf("Unexpected call to QueueMock.RemoveAll.")
		return
	}

	return m.RemoveAllFunc()
}

//RemoveAllMinimockCounter returns a count of QueueMock.RemoveAllFunc invocations
func (m *QueueMock) RemoveAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllCounter)
}

//RemoveAllMinimockPreCounter returns the value of QueueMock.RemoveAll invocations
func (m *QueueMock) RemoveAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllPreCounter)
}

//RemoveAllFinished returns true if mock invocations count is ok
func (m *QueueMock) RemoveAllFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveAllMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveAllCounter) == uint64(len(m.RemoveAllMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveAllMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveAllCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveAllFunc != nil {
		return atomic.LoadUint64(&m.RemoveAllCounter) > 0
	}

	return true
}

type mQueueMockSinkPush struct {
	mock              *QueueMock
	mainExpectation   *QueueMockSinkPushExpectation
	expectationSeries []*QueueMockSinkPushExpectation
}

type QueueMockSinkPushExpectation struct {
	input  *QueueMockSinkPushInput
	result *QueueMockSinkPushResult
}

type QueueMockSinkPushInput struct {
	p interface{}
}

type QueueMockSinkPushResult struct {
	r error
}

//Expect specifies that invocation of Queue.SinkPush is expected from 1 to Infinity times
func (m *mQueueMockSinkPush) Expect(p interface{}) *mQueueMockSinkPush {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockSinkPushExpectation{}
	}
	m.mainExpectation.input = &QueueMockSinkPushInput{p}
	return m
}

//Return specifies results of invocation of Queue.SinkPush
func (m *mQueueMockSinkPush) Return(r error) *QueueMock {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockSinkPushExpectation{}
	}
	m.mainExpectation.result = &QueueMockSinkPushResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Queue.SinkPush is expected once
func (m *mQueueMockSinkPush) ExpectOnce(p interface{}) *QueueMockSinkPushExpectation {
	m.mock.SinkPushFunc = nil
	m.mainExpectation = nil

	expectation := &QueueMockSinkPushExpectation{}
	expectation.input = &QueueMockSinkPushInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *QueueMockSinkPushExpectation) Return(r error) {
	e.result = &QueueMockSinkPushResult{r}
}

//Set uses given function f as a mock of Queue.SinkPush method
func (m *mQueueMockSinkPush) Set(f func(p interface{}) (r error)) *QueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushFunc = f
	return m.mock
}

//SinkPush implements github.com/insolar/insolar/conveyor/queue.Queue interface
func (m *QueueMock) SinkPush(p interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushCounter, 1)

	if len(m.SinkPushMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to QueueMock.SinkPush. %v", p)
			return
		}

		input := m.SinkPushMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, QueueMockSinkPushInput{p}, "Queue.SinkPush got unexpected parameters")

		result := m.SinkPushMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.SinkPush")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushMock.mainExpectation != nil {

		input := m.SinkPushMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, QueueMockSinkPushInput{p}, "Queue.SinkPush got unexpected parameters")
		}

		result := m.SinkPushMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.SinkPush")
		}

		r = result.r

		return
	}

	if m.SinkPushFunc == nil {
		m.t.Fatalf("Unexpected call to QueueMock.SinkPush. %v", p)
		return
	}

	return m.SinkPushFunc(p)
}

//SinkPushMinimockCounter returns a count of QueueMock.SinkPushFunc invocations
func (m *QueueMock) SinkPushMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushCounter)
}

//SinkPushMinimockPreCounter returns the value of QueueMock.SinkPush invocations
func (m *QueueMock) SinkPushMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushPreCounter)
}

//SinkPushFinished returns true if mock invocations count is ok
func (m *QueueMock) SinkPushFinished() bool {
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

type mQueueMockSinkPushAll struct {
	mock              *QueueMock
	mainExpectation   *QueueMockSinkPushAllExpectation
	expectationSeries []*QueueMockSinkPushAllExpectation
}

type QueueMockSinkPushAllExpectation struct {
	input  *QueueMockSinkPushAllInput
	result *QueueMockSinkPushAllResult
}

type QueueMockSinkPushAllInput struct {
	p []interface{}
}

type QueueMockSinkPushAllResult struct {
	r error
}

//Expect specifies that invocation of Queue.SinkPushAll is expected from 1 to Infinity times
func (m *mQueueMockSinkPushAll) Expect(p []interface{}) *mQueueMockSinkPushAll {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockSinkPushAllExpectation{}
	}
	m.mainExpectation.input = &QueueMockSinkPushAllInput{p}
	return m
}

//Return specifies results of invocation of Queue.SinkPushAll
func (m *mQueueMockSinkPushAll) Return(r error) *QueueMock {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockSinkPushAllExpectation{}
	}
	m.mainExpectation.result = &QueueMockSinkPushAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Queue.SinkPushAll is expected once
func (m *mQueueMockSinkPushAll) ExpectOnce(p []interface{}) *QueueMockSinkPushAllExpectation {
	m.mock.SinkPushAllFunc = nil
	m.mainExpectation = nil

	expectation := &QueueMockSinkPushAllExpectation{}
	expectation.input = &QueueMockSinkPushAllInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *QueueMockSinkPushAllExpectation) Return(r error) {
	e.result = &QueueMockSinkPushAllResult{r}
}

//Set uses given function f as a mock of Queue.SinkPushAll method
func (m *mQueueMockSinkPushAll) Set(f func(p []interface{}) (r error)) *QueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushAllFunc = f
	return m.mock
}

//SinkPushAll implements github.com/insolar/insolar/conveyor/queue.Queue interface
func (m *QueueMock) SinkPushAll(p []interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushAllPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushAllCounter, 1)

	if len(m.SinkPushAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to QueueMock.SinkPushAll. %v", p)
			return
		}

		input := m.SinkPushAllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, QueueMockSinkPushAllInput{p}, "Queue.SinkPushAll got unexpected parameters")

		result := m.SinkPushAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.SinkPushAll")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushAllMock.mainExpectation != nil {

		input := m.SinkPushAllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, QueueMockSinkPushAllInput{p}, "Queue.SinkPushAll got unexpected parameters")
		}

		result := m.SinkPushAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.SinkPushAll")
		}

		r = result.r

		return
	}

	if m.SinkPushAllFunc == nil {
		m.t.Fatalf("Unexpected call to QueueMock.SinkPushAll. %v", p)
		return
	}

	return m.SinkPushAllFunc(p)
}

//SinkPushAllMinimockCounter returns a count of QueueMock.SinkPushAllFunc invocations
func (m *QueueMock) SinkPushAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllCounter)
}

//SinkPushAllMinimockPreCounter returns the value of QueueMock.SinkPushAll invocations
func (m *QueueMock) SinkPushAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllPreCounter)
}

//SinkPushAllFinished returns true if mock invocations count is ok
func (m *QueueMock) SinkPushAllFinished() bool {
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

type mQueueMockUnblock struct {
	mock              *QueueMock
	mainExpectation   *QueueMockUnblockExpectation
	expectationSeries []*QueueMockUnblockExpectation
}

type QueueMockUnblockExpectation struct {
	result *QueueMockUnblockResult
}

type QueueMockUnblockResult struct {
	r bool
}

//Expect specifies that invocation of Queue.Unblock is expected from 1 to Infinity times
func (m *mQueueMockUnblock) Expect() *mQueueMockUnblock {
	m.mock.UnblockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockUnblockExpectation{}
	}

	return m
}

//Return specifies results of invocation of Queue.Unblock
func (m *mQueueMockUnblock) Return(r bool) *QueueMock {
	m.mock.UnblockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &QueueMockUnblockExpectation{}
	}
	m.mainExpectation.result = &QueueMockUnblockResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Queue.Unblock is expected once
func (m *mQueueMockUnblock) ExpectOnce() *QueueMockUnblockExpectation {
	m.mock.UnblockFunc = nil
	m.mainExpectation = nil

	expectation := &QueueMockUnblockExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *QueueMockUnblockExpectation) Return(r bool) {
	e.result = &QueueMockUnblockResult{r}
}

//Set uses given function f as a mock of Queue.Unblock method
func (m *mQueueMockUnblock) Set(f func() (r bool)) *QueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnblockFunc = f
	return m.mock
}

//Unblock implements github.com/insolar/insolar/conveyor/queue.Queue interface
func (m *QueueMock) Unblock() (r bool) {
	counter := atomic.AddUint64(&m.UnblockPreCounter, 1)
	defer atomic.AddUint64(&m.UnblockCounter, 1)

	if len(m.UnblockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnblockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to QueueMock.Unblock.")
			return
		}

		result := m.UnblockMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.Unblock")
			return
		}

		r = result.r

		return
	}

	if m.UnblockMock.mainExpectation != nil {

		result := m.UnblockMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the QueueMock.Unblock")
		}

		r = result.r

		return
	}

	if m.UnblockFunc == nil {
		m.t.Fatalf("Unexpected call to QueueMock.Unblock.")
		return
	}

	return m.UnblockFunc()
}

//UnblockMinimockCounter returns a count of QueueMock.UnblockFunc invocations
func (m *QueueMock) UnblockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnblockCounter)
}

//UnblockMinimockPreCounter returns the value of QueueMock.Unblock invocations
func (m *QueueMock) UnblockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnblockPreCounter)
}

//UnblockFinished returns true if mock invocations count is ok
func (m *QueueMock) UnblockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UnblockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UnblockCounter) == uint64(len(m.UnblockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UnblockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UnblockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UnblockFunc != nil {
		return atomic.LoadUint64(&m.UnblockCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *QueueMock) ValidateCallCounters() {

	if !m.BlockAndRemoveAllFinished() {
		m.t.Fatal("Expected call to QueueMock.BlockAndRemoveAll")
	}

	if !m.HasSignalFinished() {
		m.t.Fatal("Expected call to QueueMock.HasSignal")
	}

	if !m.PushSignalFinished() {
		m.t.Fatal("Expected call to QueueMock.PushSignal")
	}

	if !m.RemoveAllFinished() {
		m.t.Fatal("Expected call to QueueMock.RemoveAll")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to QueueMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to QueueMock.SinkPushAll")
	}

	if !m.UnblockFinished() {
		m.t.Fatal("Expected call to QueueMock.Unblock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *QueueMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *QueueMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *QueueMock) MinimockFinish() {

	if !m.BlockAndRemoveAllFinished() {
		m.t.Fatal("Expected call to QueueMock.BlockAndRemoveAll")
	}

	if !m.HasSignalFinished() {
		m.t.Fatal("Expected call to QueueMock.HasSignal")
	}

	if !m.PushSignalFinished() {
		m.t.Fatal("Expected call to QueueMock.PushSignal")
	}

	if !m.RemoveAllFinished() {
		m.t.Fatal("Expected call to QueueMock.RemoveAll")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to QueueMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to QueueMock.SinkPushAll")
	}

	if !m.UnblockFinished() {
		m.t.Fatal("Expected call to QueueMock.Unblock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *QueueMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *QueueMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.BlockAndRemoveAllFinished()
		ok = ok && m.HasSignalFinished()
		ok = ok && m.PushSignalFinished()
		ok = ok && m.RemoveAllFinished()
		ok = ok && m.SinkPushFinished()
		ok = ok && m.SinkPushAllFinished()
		ok = ok && m.UnblockFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.BlockAndRemoveAllFinished() {
				m.t.Error("Expected call to QueueMock.BlockAndRemoveAll")
			}

			if !m.HasSignalFinished() {
				m.t.Error("Expected call to QueueMock.HasSignal")
			}

			if !m.PushSignalFinished() {
				m.t.Error("Expected call to QueueMock.PushSignal")
			}

			if !m.RemoveAllFinished() {
				m.t.Error("Expected call to QueueMock.RemoveAll")
			}

			if !m.SinkPushFinished() {
				m.t.Error("Expected call to QueueMock.SinkPush")
			}

			if !m.SinkPushAllFinished() {
				m.t.Error("Expected call to QueueMock.SinkPushAll")
			}

			if !m.UnblockFinished() {
				m.t.Error("Expected call to QueueMock.Unblock")
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
func (m *QueueMock) AllMocksCalled() bool {

	if !m.BlockAndRemoveAllFinished() {
		return false
	}

	if !m.HasSignalFinished() {
		return false
	}

	if !m.PushSignalFinished() {
		return false
	}

	if !m.RemoveAllFinished() {
		return false
	}

	if !m.SinkPushFinished() {
		return false
	}

	if !m.SinkPushAllFinished() {
		return false
	}

	if !m.UnblockFinished() {
		return false
	}

	return true
}

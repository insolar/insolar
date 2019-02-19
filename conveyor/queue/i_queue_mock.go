package queue

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IQueue" can be found in github.com/insolar/insolar/conveyor/queue
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//IQueueMock implements github.com/insolar/insolar/conveyor/queue.IQueue
type IQueueMock struct {
	t minimock.Tester

	BlockAndRemoveAllFunc       func() (r []OutputElement)
	BlockAndRemoveAllCounter    uint64
	BlockAndRemoveAllPreCounter uint64
	BlockAndRemoveAllMock       mIQueueMockBlockAndRemoveAll

	HasSignalFunc       func() (r bool)
	HasSignalCounter    uint64
	HasSignalPreCounter uint64
	HasSignalMock       mIQueueMockHasSignal

	PushSignalFunc       func(p uint32, p1 SyncDone) (r error)
	PushSignalCounter    uint64
	PushSignalPreCounter uint64
	PushSignalMock       mIQueueMockPushSignal

	RemoveAllFunc       func() (r []OutputElement)
	RemoveAllCounter    uint64
	RemoveAllPreCounter uint64
	RemoveAllMock       mIQueueMockRemoveAll

	SinkPushFunc       func(p interface{}) (r error)
	SinkPushCounter    uint64
	SinkPushPreCounter uint64
	SinkPushMock       mIQueueMockSinkPush

	SinkPushAllFunc       func(p []interface{}) (r error)
	SinkPushAllCounter    uint64
	SinkPushAllPreCounter uint64
	SinkPushAllMock       mIQueueMockSinkPushAll

	UnblockFunc       func() (r bool)
	UnblockCounter    uint64
	UnblockPreCounter uint64
	UnblockMock       mIQueueMockUnblock
}

//NewIQueueMock returns a mock for github.com/insolar/insolar/conveyor/queue.IQueue
func NewIQueueMock(t minimock.Tester) *IQueueMock {
	m := &IQueueMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BlockAndRemoveAllMock = mIQueueMockBlockAndRemoveAll{mock: m}
	m.HasSignalMock = mIQueueMockHasSignal{mock: m}
	m.PushSignalMock = mIQueueMockPushSignal{mock: m}
	m.RemoveAllMock = mIQueueMockRemoveAll{mock: m}
	m.SinkPushMock = mIQueueMockSinkPush{mock: m}
	m.SinkPushAllMock = mIQueueMockSinkPushAll{mock: m}
	m.UnblockMock = mIQueueMockUnblock{mock: m}

	return m
}

type mIQueueMockBlockAndRemoveAll struct {
	mock              *IQueueMock
	mainExpectation   *IQueueMockBlockAndRemoveAllExpectation
	expectationSeries []*IQueueMockBlockAndRemoveAllExpectation
}

type IQueueMockBlockAndRemoveAllExpectation struct {
	result *IQueueMockBlockAndRemoveAllResult
}

type IQueueMockBlockAndRemoveAllResult struct {
	r []OutputElement
}

//Expect specifies that invocation of IQueue.BlockAndRemoveAll is expected from 1 to Infinity times
func (m *mIQueueMockBlockAndRemoveAll) Expect() *mIQueueMockBlockAndRemoveAll {
	m.mock.BlockAndRemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockBlockAndRemoveAllExpectation{}
	}

	return m
}

//Return specifies results of invocation of IQueue.BlockAndRemoveAll
func (m *mIQueueMockBlockAndRemoveAll) Return(r []OutputElement) *IQueueMock {
	m.mock.BlockAndRemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockBlockAndRemoveAllExpectation{}
	}
	m.mainExpectation.result = &IQueueMockBlockAndRemoveAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IQueue.BlockAndRemoveAll is expected once
func (m *mIQueueMockBlockAndRemoveAll) ExpectOnce() *IQueueMockBlockAndRemoveAllExpectation {
	m.mock.BlockAndRemoveAllFunc = nil
	m.mainExpectation = nil

	expectation := &IQueueMockBlockAndRemoveAllExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IQueueMockBlockAndRemoveAllExpectation) Return(r []OutputElement) {
	e.result = &IQueueMockBlockAndRemoveAllResult{r}
}

//Set uses given function f as a mock of IQueue.BlockAndRemoveAll method
func (m *mIQueueMockBlockAndRemoveAll) Set(f func() (r []OutputElement)) *IQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BlockAndRemoveAllFunc = f
	return m.mock
}

//BlockAndRemoveAll implements github.com/insolar/insolar/conveyor/queue.IQueue interface
func (m *IQueueMock) BlockAndRemoveAll() (r []OutputElement) {
	counter := atomic.AddUint64(&m.BlockAndRemoveAllPreCounter, 1)
	defer atomic.AddUint64(&m.BlockAndRemoveAllCounter, 1)

	if len(m.BlockAndRemoveAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BlockAndRemoveAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IQueueMock.BlockAndRemoveAll.")
			return
		}

		result := m.BlockAndRemoveAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.BlockAndRemoveAll")
			return
		}

		r = result.r

		return
	}

	if m.BlockAndRemoveAllMock.mainExpectation != nil {

		result := m.BlockAndRemoveAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.BlockAndRemoveAll")
		}

		r = result.r

		return
	}

	if m.BlockAndRemoveAllFunc == nil {
		m.t.Fatalf("Unexpected call to IQueueMock.BlockAndRemoveAll.")
		return
	}

	return m.BlockAndRemoveAllFunc()
}

//BlockAndRemoveAllMinimockCounter returns a count of IQueueMock.BlockAndRemoveAllFunc invocations
func (m *IQueueMock) BlockAndRemoveAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BlockAndRemoveAllCounter)
}

//BlockAndRemoveAllMinimockPreCounter returns the value of IQueueMock.BlockAndRemoveAll invocations
func (m *IQueueMock) BlockAndRemoveAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BlockAndRemoveAllPreCounter)
}

//BlockAndRemoveAllFinished returns true if mock invocations count is ok
func (m *IQueueMock) BlockAndRemoveAllFinished() bool {
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

type mIQueueMockHasSignal struct {
	mock              *IQueueMock
	mainExpectation   *IQueueMockHasSignalExpectation
	expectationSeries []*IQueueMockHasSignalExpectation
}

type IQueueMockHasSignalExpectation struct {
	result *IQueueMockHasSignalResult
}

type IQueueMockHasSignalResult struct {
	r bool
}

//Expect specifies that invocation of IQueue.HasSignal is expected from 1 to Infinity times
func (m *mIQueueMockHasSignal) Expect() *mIQueueMockHasSignal {
	m.mock.HasSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockHasSignalExpectation{}
	}

	return m
}

//Return specifies results of invocation of IQueue.HasSignal
func (m *mIQueueMockHasSignal) Return(r bool) *IQueueMock {
	m.mock.HasSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockHasSignalExpectation{}
	}
	m.mainExpectation.result = &IQueueMockHasSignalResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IQueue.HasSignal is expected once
func (m *mIQueueMockHasSignal) ExpectOnce() *IQueueMockHasSignalExpectation {
	m.mock.HasSignalFunc = nil
	m.mainExpectation = nil

	expectation := &IQueueMockHasSignalExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IQueueMockHasSignalExpectation) Return(r bool) {
	e.result = &IQueueMockHasSignalResult{r}
}

//Set uses given function f as a mock of IQueue.HasSignal method
func (m *mIQueueMockHasSignal) Set(f func() (r bool)) *IQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasSignalFunc = f
	return m.mock
}

//HasSignal implements github.com/insolar/insolar/conveyor/queue.IQueue interface
func (m *IQueueMock) HasSignal() (r bool) {
	counter := atomic.AddUint64(&m.HasSignalPreCounter, 1)
	defer atomic.AddUint64(&m.HasSignalCounter, 1)

	if len(m.HasSignalMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasSignalMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IQueueMock.HasSignal.")
			return
		}

		result := m.HasSignalMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.HasSignal")
			return
		}

		r = result.r

		return
	}

	if m.HasSignalMock.mainExpectation != nil {

		result := m.HasSignalMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.HasSignal")
		}

		r = result.r

		return
	}

	if m.HasSignalFunc == nil {
		m.t.Fatalf("Unexpected call to IQueueMock.HasSignal.")
		return
	}

	return m.HasSignalFunc()
}

//HasSignalMinimockCounter returns a count of IQueueMock.HasSignalFunc invocations
func (m *IQueueMock) HasSignalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasSignalCounter)
}

//HasSignalMinimockPreCounter returns the value of IQueueMock.HasSignal invocations
func (m *IQueueMock) HasSignalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasSignalPreCounter)
}

//HasSignalFinished returns true if mock invocations count is ok
func (m *IQueueMock) HasSignalFinished() bool {
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

type mIQueueMockPushSignal struct {
	mock              *IQueueMock
	mainExpectation   *IQueueMockPushSignalExpectation
	expectationSeries []*IQueueMockPushSignalExpectation
}

type IQueueMockPushSignalExpectation struct {
	input  *IQueueMockPushSignalInput
	result *IQueueMockPushSignalResult
}

type IQueueMockPushSignalInput struct {
	p  uint32
	p1 SyncDone
}

type IQueueMockPushSignalResult struct {
	r error
}

//Expect specifies that invocation of IQueue.PushSignal is expected from 1 to Infinity times
func (m *mIQueueMockPushSignal) Expect(p uint32, p1 SyncDone) *mIQueueMockPushSignal {
	m.mock.PushSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockPushSignalExpectation{}
	}
	m.mainExpectation.input = &IQueueMockPushSignalInput{p, p1}
	return m
}

//Return specifies results of invocation of IQueue.PushSignal
func (m *mIQueueMockPushSignal) Return(r error) *IQueueMock {
	m.mock.PushSignalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockPushSignalExpectation{}
	}
	m.mainExpectation.result = &IQueueMockPushSignalResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IQueue.PushSignal is expected once
func (m *mIQueueMockPushSignal) ExpectOnce(p uint32, p1 SyncDone) *IQueueMockPushSignalExpectation {
	m.mock.PushSignalFunc = nil
	m.mainExpectation = nil

	expectation := &IQueueMockPushSignalExpectation{}
	expectation.input = &IQueueMockPushSignalInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IQueueMockPushSignalExpectation) Return(r error) {
	e.result = &IQueueMockPushSignalResult{r}
}

//Set uses given function f as a mock of IQueue.PushSignal method
func (m *mIQueueMockPushSignal) Set(f func(p uint32, p1 SyncDone) (r error)) *IQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PushSignalFunc = f
	return m.mock
}

//PushSignal implements github.com/insolar/insolar/conveyor/queue.IQueue interface
func (m *IQueueMock) PushSignal(p uint32, p1 SyncDone) (r error) {
	counter := atomic.AddUint64(&m.PushSignalPreCounter, 1)
	defer atomic.AddUint64(&m.PushSignalCounter, 1)

	if len(m.PushSignalMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PushSignalMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IQueueMock.PushSignal. %v %v", p, p1)
			return
		}

		input := m.PushSignalMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IQueueMockPushSignalInput{p, p1}, "IQueue.PushSignal got unexpected parameters")

		result := m.PushSignalMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.PushSignal")
			return
		}

		r = result.r

		return
	}

	if m.PushSignalMock.mainExpectation != nil {

		input := m.PushSignalMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IQueueMockPushSignalInput{p, p1}, "IQueue.PushSignal got unexpected parameters")
		}

		result := m.PushSignalMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.PushSignal")
		}

		r = result.r

		return
	}

	if m.PushSignalFunc == nil {
		m.t.Fatalf("Unexpected call to IQueueMock.PushSignal. %v %v", p, p1)
		return
	}

	return m.PushSignalFunc(p, p1)
}

//PushSignalMinimockCounter returns a count of IQueueMock.PushSignalFunc invocations
func (m *IQueueMock) PushSignalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PushSignalCounter)
}

//PushSignalMinimockPreCounter returns the value of IQueueMock.PushSignal invocations
func (m *IQueueMock) PushSignalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PushSignalPreCounter)
}

//PushSignalFinished returns true if mock invocations count is ok
func (m *IQueueMock) PushSignalFinished() bool {
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

type mIQueueMockRemoveAll struct {
	mock              *IQueueMock
	mainExpectation   *IQueueMockRemoveAllExpectation
	expectationSeries []*IQueueMockRemoveAllExpectation
}

type IQueueMockRemoveAllExpectation struct {
	result *IQueueMockRemoveAllResult
}

type IQueueMockRemoveAllResult struct {
	r []OutputElement
}

//Expect specifies that invocation of IQueue.RemoveAll is expected from 1 to Infinity times
func (m *mIQueueMockRemoveAll) Expect() *mIQueueMockRemoveAll {
	m.mock.RemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockRemoveAllExpectation{}
	}

	return m
}

//Return specifies results of invocation of IQueue.RemoveAll
func (m *mIQueueMockRemoveAll) Return(r []OutputElement) *IQueueMock {
	m.mock.RemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockRemoveAllExpectation{}
	}
	m.mainExpectation.result = &IQueueMockRemoveAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IQueue.RemoveAll is expected once
func (m *mIQueueMockRemoveAll) ExpectOnce() *IQueueMockRemoveAllExpectation {
	m.mock.RemoveAllFunc = nil
	m.mainExpectation = nil

	expectation := &IQueueMockRemoveAllExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IQueueMockRemoveAllExpectation) Return(r []OutputElement) {
	e.result = &IQueueMockRemoveAllResult{r}
}

//Set uses given function f as a mock of IQueue.RemoveAll method
func (m *mIQueueMockRemoveAll) Set(f func() (r []OutputElement)) *IQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveAllFunc = f
	return m.mock
}

//RemoveAll implements github.com/insolar/insolar/conveyor/queue.IQueue interface
func (m *IQueueMock) RemoveAll() (r []OutputElement) {
	counter := atomic.AddUint64(&m.RemoveAllPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveAllCounter, 1)

	if len(m.RemoveAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IQueueMock.RemoveAll.")
			return
		}

		result := m.RemoveAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.RemoveAll")
			return
		}

		r = result.r

		return
	}

	if m.RemoveAllMock.mainExpectation != nil {

		result := m.RemoveAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.RemoveAll")
		}

		r = result.r

		return
	}

	if m.RemoveAllFunc == nil {
		m.t.Fatalf("Unexpected call to IQueueMock.RemoveAll.")
		return
	}

	return m.RemoveAllFunc()
}

//RemoveAllMinimockCounter returns a count of IQueueMock.RemoveAllFunc invocations
func (m *IQueueMock) RemoveAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllCounter)
}

//RemoveAllMinimockPreCounter returns the value of IQueueMock.RemoveAll invocations
func (m *IQueueMock) RemoveAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllPreCounter)
}

//RemoveAllFinished returns true if mock invocations count is ok
func (m *IQueueMock) RemoveAllFinished() bool {
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

type mIQueueMockSinkPush struct {
	mock              *IQueueMock
	mainExpectation   *IQueueMockSinkPushExpectation
	expectationSeries []*IQueueMockSinkPushExpectation
}

type IQueueMockSinkPushExpectation struct {
	input  *IQueueMockSinkPushInput
	result *IQueueMockSinkPushResult
}

type IQueueMockSinkPushInput struct {
	p interface{}
}

type IQueueMockSinkPushResult struct {
	r error
}

//Expect specifies that invocation of IQueue.SinkPush is expected from 1 to Infinity times
func (m *mIQueueMockSinkPush) Expect(p interface{}) *mIQueueMockSinkPush {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockSinkPushExpectation{}
	}
	m.mainExpectation.input = &IQueueMockSinkPushInput{p}
	return m
}

//Return specifies results of invocation of IQueue.SinkPush
func (m *mIQueueMockSinkPush) Return(r error) *IQueueMock {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockSinkPushExpectation{}
	}
	m.mainExpectation.result = &IQueueMockSinkPushResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IQueue.SinkPush is expected once
func (m *mIQueueMockSinkPush) ExpectOnce(p interface{}) *IQueueMockSinkPushExpectation {
	m.mock.SinkPushFunc = nil
	m.mainExpectation = nil

	expectation := &IQueueMockSinkPushExpectation{}
	expectation.input = &IQueueMockSinkPushInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IQueueMockSinkPushExpectation) Return(r error) {
	e.result = &IQueueMockSinkPushResult{r}
}

//Set uses given function f as a mock of IQueue.SinkPush method
func (m *mIQueueMockSinkPush) Set(f func(p interface{}) (r error)) *IQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushFunc = f
	return m.mock
}

//SinkPush implements github.com/insolar/insolar/conveyor/queue.IQueue interface
func (m *IQueueMock) SinkPush(p interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushCounter, 1)

	if len(m.SinkPushMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IQueueMock.SinkPush. %v", p)
			return
		}

		input := m.SinkPushMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IQueueMockSinkPushInput{p}, "IQueue.SinkPush got unexpected parameters")

		result := m.SinkPushMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.SinkPush")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushMock.mainExpectation != nil {

		input := m.SinkPushMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IQueueMockSinkPushInput{p}, "IQueue.SinkPush got unexpected parameters")
		}

		result := m.SinkPushMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.SinkPush")
		}

		r = result.r

		return
	}

	if m.SinkPushFunc == nil {
		m.t.Fatalf("Unexpected call to IQueueMock.SinkPush. %v", p)
		return
	}

	return m.SinkPushFunc(p)
}

//SinkPushMinimockCounter returns a count of IQueueMock.SinkPushFunc invocations
func (m *IQueueMock) SinkPushMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushCounter)
}

//SinkPushMinimockPreCounter returns the value of IQueueMock.SinkPush invocations
func (m *IQueueMock) SinkPushMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushPreCounter)
}

//SinkPushFinished returns true if mock invocations count is ok
func (m *IQueueMock) SinkPushFinished() bool {
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

type mIQueueMockSinkPushAll struct {
	mock              *IQueueMock
	mainExpectation   *IQueueMockSinkPushAllExpectation
	expectationSeries []*IQueueMockSinkPushAllExpectation
}

type IQueueMockSinkPushAllExpectation struct {
	input  *IQueueMockSinkPushAllInput
	result *IQueueMockSinkPushAllResult
}

type IQueueMockSinkPushAllInput struct {
	p []interface{}
}

type IQueueMockSinkPushAllResult struct {
	r error
}

//Expect specifies that invocation of IQueue.SinkPushAll is expected from 1 to Infinity times
func (m *mIQueueMockSinkPushAll) Expect(p []interface{}) *mIQueueMockSinkPushAll {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockSinkPushAllExpectation{}
	}
	m.mainExpectation.input = &IQueueMockSinkPushAllInput{p}
	return m
}

//Return specifies results of invocation of IQueue.SinkPushAll
func (m *mIQueueMockSinkPushAll) Return(r error) *IQueueMock {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockSinkPushAllExpectation{}
	}
	m.mainExpectation.result = &IQueueMockSinkPushAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IQueue.SinkPushAll is expected once
func (m *mIQueueMockSinkPushAll) ExpectOnce(p []interface{}) *IQueueMockSinkPushAllExpectation {
	m.mock.SinkPushAllFunc = nil
	m.mainExpectation = nil

	expectation := &IQueueMockSinkPushAllExpectation{}
	expectation.input = &IQueueMockSinkPushAllInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IQueueMockSinkPushAllExpectation) Return(r error) {
	e.result = &IQueueMockSinkPushAllResult{r}
}

//Set uses given function f as a mock of IQueue.SinkPushAll method
func (m *mIQueueMockSinkPushAll) Set(f func(p []interface{}) (r error)) *IQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushAllFunc = f
	return m.mock
}

//SinkPushAll implements github.com/insolar/insolar/conveyor/queue.IQueue interface
func (m *IQueueMock) SinkPushAll(p []interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushAllPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushAllCounter, 1)

	if len(m.SinkPushAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IQueueMock.SinkPushAll. %v", p)
			return
		}

		input := m.SinkPushAllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IQueueMockSinkPushAllInput{p}, "IQueue.SinkPushAll got unexpected parameters")

		result := m.SinkPushAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.SinkPushAll")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushAllMock.mainExpectation != nil {

		input := m.SinkPushAllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IQueueMockSinkPushAllInput{p}, "IQueue.SinkPushAll got unexpected parameters")
		}

		result := m.SinkPushAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.SinkPushAll")
		}

		r = result.r

		return
	}

	if m.SinkPushAllFunc == nil {
		m.t.Fatalf("Unexpected call to IQueueMock.SinkPushAll. %v", p)
		return
	}

	return m.SinkPushAllFunc(p)
}

//SinkPushAllMinimockCounter returns a count of IQueueMock.SinkPushAllFunc invocations
func (m *IQueueMock) SinkPushAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllCounter)
}

//SinkPushAllMinimockPreCounter returns the value of IQueueMock.SinkPushAll invocations
func (m *IQueueMock) SinkPushAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllPreCounter)
}

//SinkPushAllFinished returns true if mock invocations count is ok
func (m *IQueueMock) SinkPushAllFinished() bool {
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

type mIQueueMockUnblock struct {
	mock              *IQueueMock
	mainExpectation   *IQueueMockUnblockExpectation
	expectationSeries []*IQueueMockUnblockExpectation
}

type IQueueMockUnblockExpectation struct {
	result *IQueueMockUnblockResult
}

type IQueueMockUnblockResult struct {
	r bool
}

//Expect specifies that invocation of IQueue.Unblock is expected from 1 to Infinity times
func (m *mIQueueMockUnblock) Expect() *mIQueueMockUnblock {
	m.mock.UnblockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockUnblockExpectation{}
	}

	return m
}

//Return specifies results of invocation of IQueue.Unblock
func (m *mIQueueMockUnblock) Return(r bool) *IQueueMock {
	m.mock.UnblockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IQueueMockUnblockExpectation{}
	}
	m.mainExpectation.result = &IQueueMockUnblockResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IQueue.Unblock is expected once
func (m *mIQueueMockUnblock) ExpectOnce() *IQueueMockUnblockExpectation {
	m.mock.UnblockFunc = nil
	m.mainExpectation = nil

	expectation := &IQueueMockUnblockExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IQueueMockUnblockExpectation) Return(r bool) {
	e.result = &IQueueMockUnblockResult{r}
}

//Set uses given function f as a mock of IQueue.Unblock method
func (m *mIQueueMockUnblock) Set(f func() (r bool)) *IQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnblockFunc = f
	return m.mock
}

//Unblock implements github.com/insolar/insolar/conveyor/queue.IQueue interface
func (m *IQueueMock) Unblock() (r bool) {
	counter := atomic.AddUint64(&m.UnblockPreCounter, 1)
	defer atomic.AddUint64(&m.UnblockCounter, 1)

	if len(m.UnblockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnblockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IQueueMock.Unblock.")
			return
		}

		result := m.UnblockMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.Unblock")
			return
		}

		r = result.r

		return
	}

	if m.UnblockMock.mainExpectation != nil {

		result := m.UnblockMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IQueueMock.Unblock")
		}

		r = result.r

		return
	}

	if m.UnblockFunc == nil {
		m.t.Fatalf("Unexpected call to IQueueMock.Unblock.")
		return
	}

	return m.UnblockFunc()
}

//UnblockMinimockCounter returns a count of IQueueMock.UnblockFunc invocations
func (m *IQueueMock) UnblockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnblockCounter)
}

//UnblockMinimockPreCounter returns the value of IQueueMock.Unblock invocations
func (m *IQueueMock) UnblockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnblockPreCounter)
}

//UnblockFinished returns true if mock invocations count is ok
func (m *IQueueMock) UnblockFinished() bool {
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
func (m *IQueueMock) ValidateCallCounters() {

	if !m.BlockAndRemoveAllFinished() {
		m.t.Fatal("Expected call to IQueueMock.BlockAndRemoveAll")
	}

	if !m.HasSignalFinished() {
		m.t.Fatal("Expected call to IQueueMock.HasSignal")
	}

	if !m.PushSignalFinished() {
		m.t.Fatal("Expected call to IQueueMock.PushSignal")
	}

	if !m.RemoveAllFinished() {
		m.t.Fatal("Expected call to IQueueMock.RemoveAll")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to IQueueMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to IQueueMock.SinkPushAll")
	}

	if !m.UnblockFinished() {
		m.t.Fatal("Expected call to IQueueMock.Unblock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IQueueMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IQueueMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IQueueMock) MinimockFinish() {

	if !m.BlockAndRemoveAllFinished() {
		m.t.Fatal("Expected call to IQueueMock.BlockAndRemoveAll")
	}

	if !m.HasSignalFinished() {
		m.t.Fatal("Expected call to IQueueMock.HasSignal")
	}

	if !m.PushSignalFinished() {
		m.t.Fatal("Expected call to IQueueMock.PushSignal")
	}

	if !m.RemoveAllFinished() {
		m.t.Fatal("Expected call to IQueueMock.RemoveAll")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to IQueueMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to IQueueMock.SinkPushAll")
	}

	if !m.UnblockFinished() {
		m.t.Fatal("Expected call to IQueueMock.Unblock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IQueueMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IQueueMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to IQueueMock.BlockAndRemoveAll")
			}

			if !m.HasSignalFinished() {
				m.t.Error("Expected call to IQueueMock.HasSignal")
			}

			if !m.PushSignalFinished() {
				m.t.Error("Expected call to IQueueMock.PushSignal")
			}

			if !m.RemoveAllFinished() {
				m.t.Error("Expected call to IQueueMock.RemoveAll")
			}

			if !m.SinkPushFinished() {
				m.t.Error("Expected call to IQueueMock.SinkPush")
			}

			if !m.SinkPushAllFinished() {
				m.t.Error("Expected call to IQueueMock.SinkPushAll")
			}

			if !m.UnblockFinished() {
				m.t.Error("Expected call to IQueueMock.Unblock")
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
func (m *IQueueMock) AllMocksCalled() bool {

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

package conveyer

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NonBlockingQueue" can be found in github.com/insolar/insolar/conveyer
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//NonBlockingQueueMock implements github.com/insolar/insolar/conveyer.NonBlockingQueue
type NonBlockingQueueMock struct {
	t minimock.Tester

	RemoveAllFunc       func() (r []interface{})
	RemoveAllCounter    uint64
	RemoveAllPreCounter uint64
	RemoveAllMock       mNonBlockingQueueMockRemoveAll

	SinkPushFunc       func(p interface{}) (r bool)
	SinkPushCounter    uint64
	SinkPushPreCounter uint64
	SinkPushMock       mNonBlockingQueueMockSinkPush

	SinkPushAllFunc       func(p []interface{}) (r bool)
	SinkPushAllCounter    uint64
	SinkPushAllPreCounter uint64
	SinkPushAllMock       mNonBlockingQueueMockSinkPushAll
}

//NewNonBlockingQueueMock returns a mock for github.com/insolar/insolar/conveyer.NonBlockingQueue
func NewNonBlockingQueueMock(t minimock.Tester) *NonBlockingQueueMock {
	m := &NonBlockingQueueMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveAllMock = mNonBlockingQueueMockRemoveAll{mock: m}
	m.SinkPushMock = mNonBlockingQueueMockSinkPush{mock: m}
	m.SinkPushAllMock = mNonBlockingQueueMockSinkPushAll{mock: m}

	return m
}

type mNonBlockingQueueMockRemoveAll struct {
	mock              *NonBlockingQueueMock
	mainExpectation   *NonBlockingQueueMockRemoveAllExpectation
	expectationSeries []*NonBlockingQueueMockRemoveAllExpectation
}

type NonBlockingQueueMockRemoveAllExpectation struct {
	result *NonBlockingQueueMockRemoveAllResult
}

type NonBlockingQueueMockRemoveAllResult struct {
	r []interface{}
}

//Expect specifies that invocation of NonBlockingQueue.RemoveAll is expected from 1 to Infinity times
func (m *mNonBlockingQueueMockRemoveAll) Expect() *mNonBlockingQueueMockRemoveAll {
	m.mock.RemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NonBlockingQueueMockRemoveAllExpectation{}
	}

	return m
}

//Return specifies results of invocation of NonBlockingQueue.RemoveAll
func (m *mNonBlockingQueueMockRemoveAll) Return(r []interface{}) *NonBlockingQueueMock {
	m.mock.RemoveAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NonBlockingQueueMockRemoveAllExpectation{}
	}
	m.mainExpectation.result = &NonBlockingQueueMockRemoveAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NonBlockingQueue.RemoveAll is expected once
func (m *mNonBlockingQueueMockRemoveAll) ExpectOnce() *NonBlockingQueueMockRemoveAllExpectation {
	m.mock.RemoveAllFunc = nil
	m.mainExpectation = nil

	expectation := &NonBlockingQueueMockRemoveAllExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NonBlockingQueueMockRemoveAllExpectation) Return(r []interface{}) {
	e.result = &NonBlockingQueueMockRemoveAllResult{r}
}

//Set uses given function f as a mock of NonBlockingQueue.RemoveAll method
func (m *mNonBlockingQueueMockRemoveAll) Set(f func() (r []interface{})) *NonBlockingQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveAllFunc = f
	return m.mock
}

//RemoveAll implements github.com/insolar/insolar/conveyer.NonBlockingQueue interface
func (m *NonBlockingQueueMock) RemoveAll() (r []interface{}) {
	counter := atomic.AddUint64(&m.RemoveAllPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveAllCounter, 1)

	if len(m.RemoveAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NonBlockingQueueMock.RemoveAll.")
			return
		}

		result := m.RemoveAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NonBlockingQueueMock.RemoveAll")
			return
		}

		r = result.r

		return
	}

	if m.RemoveAllMock.mainExpectation != nil {

		result := m.RemoveAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NonBlockingQueueMock.RemoveAll")
		}

		r = result.r

		return
	}

	if m.RemoveAllFunc == nil {
		m.t.Fatalf("Unexpected call to NonBlockingQueueMock.RemoveAll.")
		return
	}

	return m.RemoveAllFunc()
}

//RemoveAllMinimockCounter returns a count of NonBlockingQueueMock.RemoveAllFunc invocations
func (m *NonBlockingQueueMock) RemoveAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllCounter)
}

//RemoveAllMinimockPreCounter returns the value of NonBlockingQueueMock.RemoveAll invocations
func (m *NonBlockingQueueMock) RemoveAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveAllPreCounter)
}

//RemoveAllFinished returns true if mock invocations count is ok
func (m *NonBlockingQueueMock) RemoveAllFinished() bool {
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

type mNonBlockingQueueMockSinkPush struct {
	mock              *NonBlockingQueueMock
	mainExpectation   *NonBlockingQueueMockSinkPushExpectation
	expectationSeries []*NonBlockingQueueMockSinkPushExpectation
}

type NonBlockingQueueMockSinkPushExpectation struct {
	input  *NonBlockingQueueMockSinkPushInput
	result *NonBlockingQueueMockSinkPushResult
}

type NonBlockingQueueMockSinkPushInput struct {
	p interface{}
}

type NonBlockingQueueMockSinkPushResult struct {
	r bool
}

//Expect specifies that invocation of NonBlockingQueue.SinkPush is expected from 1 to Infinity times
func (m *mNonBlockingQueueMockSinkPush) Expect(p interface{}) *mNonBlockingQueueMockSinkPush {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NonBlockingQueueMockSinkPushExpectation{}
	}
	m.mainExpectation.input = &NonBlockingQueueMockSinkPushInput{p}
	return m
}

//Return specifies results of invocation of NonBlockingQueue.SinkPush
func (m *mNonBlockingQueueMockSinkPush) Return(r bool) *NonBlockingQueueMock {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NonBlockingQueueMockSinkPushExpectation{}
	}
	m.mainExpectation.result = &NonBlockingQueueMockSinkPushResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NonBlockingQueue.SinkPush is expected once
func (m *mNonBlockingQueueMockSinkPush) ExpectOnce(p interface{}) *NonBlockingQueueMockSinkPushExpectation {
	m.mock.SinkPushFunc = nil
	m.mainExpectation = nil

	expectation := &NonBlockingQueueMockSinkPushExpectation{}
	expectation.input = &NonBlockingQueueMockSinkPushInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NonBlockingQueueMockSinkPushExpectation) Return(r bool) {
	e.result = &NonBlockingQueueMockSinkPushResult{r}
}

//Set uses given function f as a mock of NonBlockingQueue.SinkPush method
func (m *mNonBlockingQueueMockSinkPush) Set(f func(p interface{}) (r bool)) *NonBlockingQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushFunc = f
	return m.mock
}

//SinkPush implements github.com/insolar/insolar/conveyer.NonBlockingQueue interface
func (m *NonBlockingQueueMock) SinkPush(p interface{}) (r bool) {
	counter := atomic.AddUint64(&m.SinkPushPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushCounter, 1)

	if len(m.SinkPushMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NonBlockingQueueMock.SinkPush. %v", p)
			return
		}

		input := m.SinkPushMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NonBlockingQueueMockSinkPushInput{p}, "NonBlockingQueue.SinkPush got unexpected parameters")

		result := m.SinkPushMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NonBlockingQueueMock.SinkPush")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushMock.mainExpectation != nil {

		input := m.SinkPushMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NonBlockingQueueMockSinkPushInput{p}, "NonBlockingQueue.SinkPush got unexpected parameters")
		}

		result := m.SinkPushMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NonBlockingQueueMock.SinkPush")
		}

		r = result.r

		return
	}

	if m.SinkPushFunc == nil {
		m.t.Fatalf("Unexpected call to NonBlockingQueueMock.SinkPush. %v", p)
		return
	}

	return m.SinkPushFunc(p)
}

//SinkPushMinimockCounter returns a count of NonBlockingQueueMock.SinkPushFunc invocations
func (m *NonBlockingQueueMock) SinkPushMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushCounter)
}

//SinkPushMinimockPreCounter returns the value of NonBlockingQueueMock.SinkPush invocations
func (m *NonBlockingQueueMock) SinkPushMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushPreCounter)
}

//SinkPushFinished returns true if mock invocations count is ok
func (m *NonBlockingQueueMock) SinkPushFinished() bool {
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

type mNonBlockingQueueMockSinkPushAll struct {
	mock              *NonBlockingQueueMock
	mainExpectation   *NonBlockingQueueMockSinkPushAllExpectation
	expectationSeries []*NonBlockingQueueMockSinkPushAllExpectation
}

type NonBlockingQueueMockSinkPushAllExpectation struct {
	input  *NonBlockingQueueMockSinkPushAllInput
	result *NonBlockingQueueMockSinkPushAllResult
}

type NonBlockingQueueMockSinkPushAllInput struct {
	p []interface{}
}

type NonBlockingQueueMockSinkPushAllResult struct {
	r bool
}

//Expect specifies that invocation of NonBlockingQueue.SinkPushAll is expected from 1 to Infinity times
func (m *mNonBlockingQueueMockSinkPushAll) Expect(p []interface{}) *mNonBlockingQueueMockSinkPushAll {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NonBlockingQueueMockSinkPushAllExpectation{}
	}
	m.mainExpectation.input = &NonBlockingQueueMockSinkPushAllInput{p}
	return m
}

//Return specifies results of invocation of NonBlockingQueue.SinkPushAll
func (m *mNonBlockingQueueMockSinkPushAll) Return(r bool) *NonBlockingQueueMock {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NonBlockingQueueMockSinkPushAllExpectation{}
	}
	m.mainExpectation.result = &NonBlockingQueueMockSinkPushAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NonBlockingQueue.SinkPushAll is expected once
func (m *mNonBlockingQueueMockSinkPushAll) ExpectOnce(p []interface{}) *NonBlockingQueueMockSinkPushAllExpectation {
	m.mock.SinkPushAllFunc = nil
	m.mainExpectation = nil

	expectation := &NonBlockingQueueMockSinkPushAllExpectation{}
	expectation.input = &NonBlockingQueueMockSinkPushAllInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NonBlockingQueueMockSinkPushAllExpectation) Return(r bool) {
	e.result = &NonBlockingQueueMockSinkPushAllResult{r}
}

//Set uses given function f as a mock of NonBlockingQueue.SinkPushAll method
func (m *mNonBlockingQueueMockSinkPushAll) Set(f func(p []interface{}) (r bool)) *NonBlockingQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushAllFunc = f
	return m.mock
}

//SinkPushAll implements github.com/insolar/insolar/conveyer.NonBlockingQueue interface
func (m *NonBlockingQueueMock) SinkPushAll(p []interface{}) (r bool) {
	counter := atomic.AddUint64(&m.SinkPushAllPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushAllCounter, 1)

	if len(m.SinkPushAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NonBlockingQueueMock.SinkPushAll. %v", p)
			return
		}

		input := m.SinkPushAllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NonBlockingQueueMockSinkPushAllInput{p}, "NonBlockingQueue.SinkPushAll got unexpected parameters")

		result := m.SinkPushAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NonBlockingQueueMock.SinkPushAll")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushAllMock.mainExpectation != nil {

		input := m.SinkPushAllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NonBlockingQueueMockSinkPushAllInput{p}, "NonBlockingQueue.SinkPushAll got unexpected parameters")
		}

		result := m.SinkPushAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NonBlockingQueueMock.SinkPushAll")
		}

		r = result.r

		return
	}

	if m.SinkPushAllFunc == nil {
		m.t.Fatalf("Unexpected call to NonBlockingQueueMock.SinkPushAll. %v", p)
		return
	}

	return m.SinkPushAllFunc(p)
}

//SinkPushAllMinimockCounter returns a count of NonBlockingQueueMock.SinkPushAllFunc invocations
func (m *NonBlockingQueueMock) SinkPushAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllCounter)
}

//SinkPushAllMinimockPreCounter returns the value of NonBlockingQueueMock.SinkPushAll invocations
func (m *NonBlockingQueueMock) SinkPushAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllPreCounter)
}

//SinkPushAllFinished returns true if mock invocations count is ok
func (m *NonBlockingQueueMock) SinkPushAllFinished() bool {
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
func (m *NonBlockingQueueMock) ValidateCallCounters() {

	if !m.RemoveAllFinished() {
		m.t.Fatal("Expected call to NonBlockingQueueMock.RemoveAll")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to NonBlockingQueueMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to NonBlockingQueueMock.SinkPushAll")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NonBlockingQueueMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NonBlockingQueueMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NonBlockingQueueMock) MinimockFinish() {

	if !m.RemoveAllFinished() {
		m.t.Fatal("Expected call to NonBlockingQueueMock.RemoveAll")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to NonBlockingQueueMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to NonBlockingQueueMock.SinkPushAll")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NonBlockingQueueMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NonBlockingQueueMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RemoveAllFinished()
		ok = ok && m.SinkPushFinished()
		ok = ok && m.SinkPushAllFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveAllFinished() {
				m.t.Error("Expected call to NonBlockingQueueMock.RemoveAll")
			}

			if !m.SinkPushFinished() {
				m.t.Error("Expected call to NonBlockingQueueMock.SinkPush")
			}

			if !m.SinkPushAllFinished() {
				m.t.Error("Expected call to NonBlockingQueueMock.SinkPushAll")
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
func (m *NonBlockingQueueMock) AllMocksCalled() bool {

	if !m.RemoveAllFinished() {
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

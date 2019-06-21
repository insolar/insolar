package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FilamentCacheManager" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//FilamentCacheManagerMock implements github.com/insolar/insolar/ledger/object.FilamentCacheManager
type FilamentCacheManagerMock struct {
	t minimock.Tester

	GatherFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)
	GatherCounter    uint64
	GatherPreCounter uint64
	GatherMock       mFilamentCacheManagerMockGather

	SendAbandonedNotificationFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)
	SendAbandonedNotificationCounter    uint64
	SendAbandonedNotificationPreCounter uint64
	SendAbandonedNotificationMock       mFilamentCacheManagerMockSendAbandonedNotification
}

//NewFilamentCacheManagerMock returns a mock for github.com/insolar/insolar/ledger/object.FilamentCacheManager
func NewFilamentCacheManagerMock(t minimock.Tester) *FilamentCacheManagerMock {
	m := &FilamentCacheManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GatherMock = mFilamentCacheManagerMockGather{mock: m}
	m.SendAbandonedNotificationMock = mFilamentCacheManagerMockSendAbandonedNotification{mock: m}

	return m
}

type mFilamentCacheManagerMockGather struct {
	mock              *FilamentCacheManagerMock
	mainExpectation   *FilamentCacheManagerMockGatherExpectation
	expectationSeries []*FilamentCacheManagerMockGatherExpectation
}

type FilamentCacheManagerMockGatherExpectation struct {
	input  *FilamentCacheManagerMockGatherInput
	result *FilamentCacheManagerMockGatherResult
}

type FilamentCacheManagerMockGatherInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type FilamentCacheManagerMockGatherResult struct {
	r error
}

//Expect specifies that invocation of FilamentCacheManager.Gather is expected from 1 to Infinity times
func (m *mFilamentCacheManagerMockGather) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mFilamentCacheManagerMockGather {
	m.mock.GatherFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCacheManagerMockGatherExpectation{}
	}
	m.mainExpectation.input = &FilamentCacheManagerMockGatherInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of FilamentCacheManager.Gather
func (m *mFilamentCacheManagerMockGather) Return(r error) *FilamentCacheManagerMock {
	m.mock.GatherFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCacheManagerMockGatherExpectation{}
	}
	m.mainExpectation.result = &FilamentCacheManagerMockGatherResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCacheManager.Gather is expected once
func (m *mFilamentCacheManagerMockGather) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *FilamentCacheManagerMockGatherExpectation {
	m.mock.GatherFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCacheManagerMockGatherExpectation{}
	expectation.input = &FilamentCacheManagerMockGatherInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCacheManagerMockGatherExpectation) Return(r error) {
	e.result = &FilamentCacheManagerMockGatherResult{r}
}

//Set uses given function f as a mock of FilamentCacheManager.Gather method
func (m *mFilamentCacheManagerMockGather) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)) *FilamentCacheManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GatherFunc = f
	return m.mock
}

//Gather implements github.com/insolar/insolar/ledger/object.FilamentCacheManager interface
func (m *FilamentCacheManagerMock) Gather(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.GatherPreCounter, 1)
	defer atomic.AddUint64(&m.GatherCounter, 1)

	if len(m.GatherMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GatherMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCacheManagerMock.Gather. %v %v %v", p, p1, p2)
			return
		}

		input := m.GatherMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCacheManagerMockGatherInput{p, p1, p2}, "FilamentCacheManager.Gather got unexpected parameters")

		result := m.GatherMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCacheManagerMock.Gather")
			return
		}

		r = result.r

		return
	}

	if m.GatherMock.mainExpectation != nil {

		input := m.GatherMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCacheManagerMockGatherInput{p, p1, p2}, "FilamentCacheManager.Gather got unexpected parameters")
		}

		result := m.GatherMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCacheManagerMock.Gather")
		}

		r = result.r

		return
	}

	if m.GatherFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCacheManagerMock.Gather. %v %v %v", p, p1, p2)
		return
	}

	return m.GatherFunc(p, p1, p2)
}

//GatherMinimockCounter returns a count of FilamentCacheManagerMock.GatherFunc invocations
func (m *FilamentCacheManagerMock) GatherMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GatherCounter)
}

//GatherMinimockPreCounter returns the value of FilamentCacheManagerMock.Gather invocations
func (m *FilamentCacheManagerMock) GatherMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GatherPreCounter)
}

//GatherFinished returns true if mock invocations count is ok
func (m *FilamentCacheManagerMock) GatherFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GatherMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GatherCounter) == uint64(len(m.GatherMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GatherMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GatherCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GatherFunc != nil {
		return atomic.LoadUint64(&m.GatherCounter) > 0
	}

	return true
}

type mFilamentCacheManagerMockSendAbandonedNotification struct {
	mock              *FilamentCacheManagerMock
	mainExpectation   *FilamentCacheManagerMockSendAbandonedNotificationExpectation
	expectationSeries []*FilamentCacheManagerMockSendAbandonedNotificationExpectation
}

type FilamentCacheManagerMockSendAbandonedNotificationExpectation struct {
	input  *FilamentCacheManagerMockSendAbandonedNotificationInput
	result *FilamentCacheManagerMockSendAbandonedNotificationResult
}

type FilamentCacheManagerMockSendAbandonedNotificationInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type FilamentCacheManagerMockSendAbandonedNotificationResult struct {
	r error
}

//Expect specifies that invocation of FilamentCacheManager.SendAbandonedNotification is expected from 1 to Infinity times
func (m *mFilamentCacheManagerMockSendAbandonedNotification) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mFilamentCacheManagerMockSendAbandonedNotification {
	m.mock.SendAbandonedNotificationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCacheManagerMockSendAbandonedNotificationExpectation{}
	}
	m.mainExpectation.input = &FilamentCacheManagerMockSendAbandonedNotificationInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of FilamentCacheManager.SendAbandonedNotification
func (m *mFilamentCacheManagerMockSendAbandonedNotification) Return(r error) *FilamentCacheManagerMock {
	m.mock.SendAbandonedNotificationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCacheManagerMockSendAbandonedNotificationExpectation{}
	}
	m.mainExpectation.result = &FilamentCacheManagerMockSendAbandonedNotificationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCacheManager.SendAbandonedNotification is expected once
func (m *mFilamentCacheManagerMockSendAbandonedNotification) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *FilamentCacheManagerMockSendAbandonedNotificationExpectation {
	m.mock.SendAbandonedNotificationFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCacheManagerMockSendAbandonedNotificationExpectation{}
	expectation.input = &FilamentCacheManagerMockSendAbandonedNotificationInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentCacheManagerMockSendAbandonedNotificationExpectation) Return(r error) {
	e.result = &FilamentCacheManagerMockSendAbandonedNotificationResult{r}
}

//Set uses given function f as a mock of FilamentCacheManager.SendAbandonedNotification method
func (m *mFilamentCacheManagerMockSendAbandonedNotification) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)) *FilamentCacheManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendAbandonedNotificationFunc = f
	return m.mock
}

//SendAbandonedNotification implements github.com/insolar/insolar/ledger/object.FilamentCacheManager interface
func (m *FilamentCacheManagerMock) SendAbandonedNotification(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SendAbandonedNotificationPreCounter, 1)
	defer atomic.AddUint64(&m.SendAbandonedNotificationCounter, 1)

	if len(m.SendAbandonedNotificationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendAbandonedNotificationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCacheManagerMock.SendAbandonedNotification. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendAbandonedNotificationMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCacheManagerMockSendAbandonedNotificationInput{p, p1, p2}, "FilamentCacheManager.SendAbandonedNotification got unexpected parameters")

		result := m.SendAbandonedNotificationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCacheManagerMock.SendAbandonedNotification")
			return
		}

		r = result.r

		return
	}

	if m.SendAbandonedNotificationMock.mainExpectation != nil {

		input := m.SendAbandonedNotificationMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCacheManagerMockSendAbandonedNotificationInput{p, p1, p2}, "FilamentCacheManager.SendAbandonedNotification got unexpected parameters")
		}

		result := m.SendAbandonedNotificationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentCacheManagerMock.SendAbandonedNotification")
		}

		r = result.r

		return
	}

	if m.SendAbandonedNotificationFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCacheManagerMock.SendAbandonedNotification. %v %v %v", p, p1, p2)
		return
	}

	return m.SendAbandonedNotificationFunc(p, p1, p2)
}

//SendAbandonedNotificationMinimockCounter returns a count of FilamentCacheManagerMock.SendAbandonedNotificationFunc invocations
func (m *FilamentCacheManagerMock) SendAbandonedNotificationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendAbandonedNotificationCounter)
}

//SendAbandonedNotificationMinimockPreCounter returns the value of FilamentCacheManagerMock.SendAbandonedNotification invocations
func (m *FilamentCacheManagerMock) SendAbandonedNotificationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendAbandonedNotificationPreCounter)
}

//SendAbandonedNotificationFinished returns true if mock invocations count is ok
func (m *FilamentCacheManagerMock) SendAbandonedNotificationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendAbandonedNotificationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendAbandonedNotificationCounter) == uint64(len(m.SendAbandonedNotificationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendAbandonedNotificationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendAbandonedNotificationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendAbandonedNotificationFunc != nil {
		return atomic.LoadUint64(&m.SendAbandonedNotificationCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCacheManagerMock) ValidateCallCounters() {

	if !m.GatherFinished() {
		m.t.Fatal("Expected call to FilamentCacheManagerMock.Gather")
	}

	if !m.SendAbandonedNotificationFinished() {
		m.t.Fatal("Expected call to FilamentCacheManagerMock.SendAbandonedNotification")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCacheManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FilamentCacheManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FilamentCacheManagerMock) MinimockFinish() {

	if !m.GatherFinished() {
		m.t.Fatal("Expected call to FilamentCacheManagerMock.Gather")
	}

	if !m.SendAbandonedNotificationFinished() {
		m.t.Fatal("Expected call to FilamentCacheManagerMock.SendAbandonedNotification")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FilamentCacheManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FilamentCacheManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GatherFinished()
		ok = ok && m.SendAbandonedNotificationFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GatherFinished() {
				m.t.Error("Expected call to FilamentCacheManagerMock.Gather")
			}

			if !m.SendAbandonedNotificationFinished() {
				m.t.Error("Expected call to FilamentCacheManagerMock.SendAbandonedNotification")
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
func (m *FilamentCacheManagerMock) AllMocksCalled() bool {

	if !m.GatherFinished() {
		return false
	}

	if !m.SendAbandonedNotificationFinished() {
		return false
	}

	return true
}

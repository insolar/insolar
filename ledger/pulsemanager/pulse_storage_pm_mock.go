package pulsemanager

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "pulseStoragePm" can be found in github.com/insolar/insolar/ledger/pulsemanager
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//pulseStoragePmMock implements github.com/insolar/insolar/ledger/pulsemanager.pulseStoragePm
type pulseStoragePmMock struct {
	t minimock.Tester

	CurrentFunc       func(p context.Context) (r *core.Pulse, r1 error)
	CurrentCounter    uint64
	CurrentPreCounter uint64
	CurrentMock       mpulseStoragePmMockCurrent

	LockFunc       func()
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mpulseStoragePmMockLock

	UnlockFunc       func()
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mpulseStoragePmMockUnlock
}

//NewpulseStoragePmMock returns a mock for github.com/insolar/insolar/ledger/pulsemanager.pulseStoragePm
func NewpulseStoragePmMock(t minimock.Tester) *pulseStoragePmMock {
	m := &pulseStoragePmMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CurrentMock = mpulseStoragePmMockCurrent{mock: m}
	m.LockMock = mpulseStoragePmMockLock{mock: m}
	m.UnlockMock = mpulseStoragePmMockUnlock{mock: m}

	return m
}

type mpulseStoragePmMockCurrent struct {
	mock              *pulseStoragePmMock
	mainExpectation   *pulseStoragePmMockCurrentExpectation
	expectationSeries []*pulseStoragePmMockCurrentExpectation
}

type pulseStoragePmMockCurrentExpectation struct {
	input  *pulseStoragePmMockCurrentInput
	result *pulseStoragePmMockCurrentResult
}

type pulseStoragePmMockCurrentInput struct {
	p context.Context
}

type pulseStoragePmMockCurrentResult struct {
	r  *core.Pulse
	r1 error
}

//Expect specifies that invocation of pulseStoragePm.Current is expected from 1 to Infinity times
func (m *mpulseStoragePmMockCurrent) Expect(p context.Context) *mpulseStoragePmMockCurrent {
	m.mock.CurrentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &pulseStoragePmMockCurrentExpectation{}
	}
	m.mainExpectation.input = &pulseStoragePmMockCurrentInput{p}
	return m
}

//Return specifies results of invocation of pulseStoragePm.Current
func (m *mpulseStoragePmMockCurrent) Return(r *core.Pulse, r1 error) *pulseStoragePmMock {
	m.mock.CurrentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &pulseStoragePmMockCurrentExpectation{}
	}
	m.mainExpectation.result = &pulseStoragePmMockCurrentResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of pulseStoragePm.Current is expected once
func (m *mpulseStoragePmMockCurrent) ExpectOnce(p context.Context) *pulseStoragePmMockCurrentExpectation {
	m.mock.CurrentFunc = nil
	m.mainExpectation = nil

	expectation := &pulseStoragePmMockCurrentExpectation{}
	expectation.input = &pulseStoragePmMockCurrentInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *pulseStoragePmMockCurrentExpectation) Return(r *core.Pulse, r1 error) {
	e.result = &pulseStoragePmMockCurrentResult{r, r1}
}

//Set uses given function f as a mock of pulseStoragePm.Current method
func (m *mpulseStoragePmMockCurrent) Set(f func(p context.Context) (r *core.Pulse, r1 error)) *pulseStoragePmMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CurrentFunc = f
	return m.mock
}

//Current implements github.com/insolar/insolar/ledger/pulsemanager.pulseStoragePm interface
func (m *pulseStoragePmMock) Current(p context.Context) (r *core.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.CurrentPreCounter, 1)
	defer atomic.AddUint64(&m.CurrentCounter, 1)

	if len(m.CurrentMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CurrentMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to pulseStoragePmMock.Current. %v", p)
			return
		}

		input := m.CurrentMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, pulseStoragePmMockCurrentInput{p}, "pulseStoragePm.Current got unexpected parameters")

		result := m.CurrentMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the pulseStoragePmMock.Current")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CurrentMock.mainExpectation != nil {

		input := m.CurrentMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, pulseStoragePmMockCurrentInput{p}, "pulseStoragePm.Current got unexpected parameters")
		}

		result := m.CurrentMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the pulseStoragePmMock.Current")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CurrentFunc == nil {
		m.t.Fatalf("Unexpected call to pulseStoragePmMock.Current. %v", p)
		return
	}

	return m.CurrentFunc(p)
}

//CurrentMinimockCounter returns a count of pulseStoragePmMock.CurrentFunc invocations
func (m *pulseStoragePmMock) CurrentMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CurrentCounter)
}

//CurrentMinimockPreCounter returns the value of pulseStoragePmMock.Current invocations
func (m *pulseStoragePmMock) CurrentMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CurrentPreCounter)
}

//CurrentFinished returns true if mock invocations count is ok
func (m *pulseStoragePmMock) CurrentFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CurrentMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CurrentCounter) == uint64(len(m.CurrentMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CurrentMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CurrentCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CurrentFunc != nil {
		return atomic.LoadUint64(&m.CurrentCounter) > 0
	}

	return true
}

type mpulseStoragePmMockLock struct {
	mock              *pulseStoragePmMock
	mainExpectation   *pulseStoragePmMockLockExpectation
	expectationSeries []*pulseStoragePmMockLockExpectation
}

type pulseStoragePmMockLockExpectation struct {
}

//Expect specifies that invocation of pulseStoragePm.Lock is expected from 1 to Infinity times
func (m *mpulseStoragePmMockLock) Expect() *mpulseStoragePmMockLock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &pulseStoragePmMockLockExpectation{}
	}

	return m
}

//Return specifies results of invocation of pulseStoragePm.Lock
func (m *mpulseStoragePmMockLock) Return() *pulseStoragePmMock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &pulseStoragePmMockLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of pulseStoragePm.Lock is expected once
func (m *mpulseStoragePmMockLock) ExpectOnce() *pulseStoragePmMockLockExpectation {
	m.mock.LockFunc = nil
	m.mainExpectation = nil

	expectation := &pulseStoragePmMockLockExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of pulseStoragePm.Lock method
func (m *mpulseStoragePmMockLock) Set(f func()) *pulseStoragePmMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LockFunc = f
	return m.mock
}

//Lock implements github.com/insolar/insolar/ledger/pulsemanager.pulseStoragePm interface
func (m *pulseStoragePmMock) Lock() {
	counter := atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if len(m.LockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to pulseStoragePmMock.Lock.")
			return
		}

		return
	}

	if m.LockMock.mainExpectation != nil {

		return
	}

	if m.LockFunc == nil {
		m.t.Fatalf("Unexpected call to pulseStoragePmMock.Lock.")
		return
	}

	m.LockFunc()
}

//LockMinimockCounter returns a count of pulseStoragePmMock.LockFunc invocations
func (m *pulseStoragePmMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of pulseStoragePmMock.Lock invocations
func (m *pulseStoragePmMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

//LockFinished returns true if mock invocations count is ok
func (m *pulseStoragePmMock) LockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LockCounter) == uint64(len(m.LockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LockFunc != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	return true
}

type mpulseStoragePmMockUnlock struct {
	mock              *pulseStoragePmMock
	mainExpectation   *pulseStoragePmMockUnlockExpectation
	expectationSeries []*pulseStoragePmMockUnlockExpectation
}

type pulseStoragePmMockUnlockExpectation struct {
}

//Expect specifies that invocation of pulseStoragePm.Unlock is expected from 1 to Infinity times
func (m *mpulseStoragePmMockUnlock) Expect() *mpulseStoragePmMockUnlock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &pulseStoragePmMockUnlockExpectation{}
	}

	return m
}

//Return specifies results of invocation of pulseStoragePm.Unlock
func (m *mpulseStoragePmMockUnlock) Return() *pulseStoragePmMock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &pulseStoragePmMockUnlockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of pulseStoragePm.Unlock is expected once
func (m *mpulseStoragePmMockUnlock) ExpectOnce() *pulseStoragePmMockUnlockExpectation {
	m.mock.UnlockFunc = nil
	m.mainExpectation = nil

	expectation := &pulseStoragePmMockUnlockExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of pulseStoragePm.Unlock method
func (m *mpulseStoragePmMockUnlock) Set(f func()) *pulseStoragePmMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnlockFunc = f
	return m.mock
}

//Unlock implements github.com/insolar/insolar/ledger/pulsemanager.pulseStoragePm interface
func (m *pulseStoragePmMock) Unlock() {
	counter := atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if len(m.UnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to pulseStoragePmMock.Unlock.")
			return
		}

		return
	}

	if m.UnlockMock.mainExpectation != nil {

		return
	}

	if m.UnlockFunc == nil {
		m.t.Fatalf("Unexpected call to pulseStoragePmMock.Unlock.")
		return
	}

	m.UnlockFunc()
}

//UnlockMinimockCounter returns a count of pulseStoragePmMock.UnlockFunc invocations
func (m *pulseStoragePmMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of pulseStoragePmMock.Unlock invocations
func (m *pulseStoragePmMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//UnlockFinished returns true if mock invocations count is ok
func (m *pulseStoragePmMock) UnlockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UnlockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UnlockCounter) == uint64(len(m.UnlockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UnlockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UnlockFunc != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *pulseStoragePmMock) ValidateCallCounters() {

	if !m.CurrentFinished() {
		m.t.Fatal("Expected call to pulseStoragePmMock.Current")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to pulseStoragePmMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to pulseStoragePmMock.Unlock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *pulseStoragePmMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *pulseStoragePmMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *pulseStoragePmMock) MinimockFinish() {

	if !m.CurrentFinished() {
		m.t.Fatal("Expected call to pulseStoragePmMock.Current")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to pulseStoragePmMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to pulseStoragePmMock.Unlock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *pulseStoragePmMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *pulseStoragePmMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CurrentFinished()
		ok = ok && m.LockFinished()
		ok = ok && m.UnlockFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CurrentFinished() {
				m.t.Error("Expected call to pulseStoragePmMock.Current")
			}

			if !m.LockFinished() {
				m.t.Error("Expected call to pulseStoragePmMock.Lock")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to pulseStoragePmMock.Unlock")
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
func (m *pulseStoragePmMock) AllMocksCalled() bool {

	if !m.CurrentFinished() {
		return false
	}

	if !m.LockFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	return true
}

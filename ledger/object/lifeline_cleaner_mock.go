package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LifelineCleaner" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// LifelineCleanerMock implements github.com/insolar/insolar/ledger/object.LifelineCleaner
type LifelineCleanerMock struct {
	t minimock.Tester

	DeleteForPNFunc       func(p context.Context, p1 insolar.PulseNumber)
	DeleteForPNCounter    uint64
	DeleteForPNPreCounter uint64
	DeleteForPNMock       mLifelineCleanerMockDeleteForPN
}

// NewLifelineCleanerMock returns a mock for github.com/insolar/insolar/ledger/object.LifelineCleaner
func NewLifelineCleanerMock(t minimock.Tester) *LifelineCleanerMock {
	m := &LifelineCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteForPNMock = mLifelineCleanerMockDeleteForPN{mock: m}

	return m
}

type mLifelineCleanerMockDeleteForPN struct {
	mock              *LifelineCleanerMock
	mainExpectation   *LifelineCleanerMockDeleteForPNExpectation
	expectationSeries []*LifelineCleanerMockDeleteForPNExpectation
}

type LifelineCleanerMockDeleteForPNExpectation struct {
	input *LifelineCleanerMockDeleteForPNInput
}

type LifelineCleanerMockDeleteForPNInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

// Expect specifies that invocation of LifelineCleaner.DeleteForPN is expected from 1 to Infinity times
func (m *mLifelineCleanerMockDeleteForPN) Expect(p context.Context, p1 insolar.PulseNumber) *mLifelineCleanerMockDeleteForPN {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineCleanerMockDeleteForPNExpectation{}
	}
	m.mainExpectation.input = &LifelineCleanerMockDeleteForPNInput{p, p1}
	return m
}

// Return specifies results of invocation of LifelineCleaner.DeleteForPN
func (m *mLifelineCleanerMockDeleteForPN) Return() *LifelineCleanerMock {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LifelineCleanerMockDeleteForPNExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of LifelineCleaner.DeleteForPN is expected once
func (m *mLifelineCleanerMockDeleteForPN) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *LifelineCleanerMockDeleteForPNExpectation {
	m.mock.DeleteForPNFunc = nil
	m.mainExpectation = nil

	expectation := &LifelineCleanerMockDeleteForPNExpectation{}
	expectation.input = &LifelineCleanerMockDeleteForPNInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of LifelineCleaner.DeleteForPN method
func (m *mLifelineCleanerMockDeleteForPN) Set(f func(p context.Context, p1 insolar.PulseNumber)) *LifelineCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteForPNFunc = f
	return m.mock
}

// DeleteForPN implements github.com/insolar/insolar/ledger/object.LifelineCleaner interface
func (m *LifelineCleanerMock) DeleteForPN(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.DeleteForPNPreCounter, 1)
	defer atomic.AddUint64(&m.DeleteForPNCounter, 1)

	if len(m.DeleteForPNMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteForPNMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LifelineCleanerMock.DeleteForPN. %v %v", p, p1)
			return
		}

		input := m.DeleteForPNMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LifelineCleanerMockDeleteForPNInput{p, p1}, "LifelineCleaner.DeleteForPN got unexpected parameters")

		return
	}

	if m.DeleteForPNMock.mainExpectation != nil {

		input := m.DeleteForPNMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LifelineCleanerMockDeleteForPNInput{p, p1}, "LifelineCleaner.DeleteForPN got unexpected parameters")
		}

		return
	}

	if m.DeleteForPNFunc == nil {
		m.t.Fatalf("Unexpected call to LifelineCleanerMock.DeleteForPN. %v %v", p, p1)
		return
	}

	m.DeleteForPNFunc(p, p1)
}

// DeleteForPNMinimockCounter returns a count of LifelineCleanerMock.DeleteForPNFunc invocations
func (m *LifelineCleanerMock) DeleteForPNMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNCounter)
}

// DeleteForPNMinimockPreCounter returns the value of LifelineCleanerMock.DeleteForPN invocations
func (m *LifelineCleanerMock) DeleteForPNMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNPreCounter)
}

// DeleteForPNFinished returns true if mock invocations count is ok
func (m *LifelineCleanerMock) DeleteForPNFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeleteForPNMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeleteForPNCounter) == uint64(len(m.DeleteForPNMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeleteForPNMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeleteForPNCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeleteForPNFunc != nil {
		return atomic.LoadUint64(&m.DeleteForPNCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineCleanerMock) ValidateCallCounters() {

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to LifelineCleanerMock.DeleteForPN")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LifelineCleanerMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LifelineCleanerMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LifelineCleanerMock) MinimockFinish() {

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to LifelineCleanerMock.DeleteForPN")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LifelineCleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *LifelineCleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.DeleteForPNFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeleteForPNFinished() {
				m.t.Error("Expected call to LifelineCleanerMock.DeleteForPN")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *LifelineCleanerMock) AllMocksCalled() bool {

	if !m.DeleteForPNFinished() {
		return false
	}

	return true
}

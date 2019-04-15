package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexCleaner" can be found in github.com/insolar/insolar/ledger/storage/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexCleanerMock implements github.com/insolar/insolar/ledger/storage/object.IndexCleaner
type IndexCleanerMock struct {
	t minimock.Tester

	RemoveUntilFunc func(p context.Context, p1 insolar.PulseNumber, p2 map[insolar.ID]struct {
	})
	RemoveUntilCounter    uint64
	RemoveUntilPreCounter uint64
	RemoveUntilMock       mIndexCleanerMockRemoveUntil
}

//NewIndexCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage/object.IndexCleaner
func NewIndexCleanerMock(t minimock.Tester) *IndexCleanerMock {
	m := &IndexCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveUntilMock = mIndexCleanerMockRemoveUntil{mock: m}

	return m
}

type mIndexCleanerMockRemoveUntil struct {
	mock              *IndexCleanerMock
	mainExpectation   *IndexCleanerMockRemoveUntilExpectation
	expectationSeries []*IndexCleanerMockRemoveUntilExpectation
}

type IndexCleanerMockRemoveUntilExpectation struct {
	input *IndexCleanerMockRemoveUntilInput
}

type IndexCleanerMockRemoveUntilInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 map[insolar.ID]struct {
	}
}

// Expect specifies that invocation of IndexCleaner.RemoveForPulse is expected from 1 to Infinity times
func (m *mIndexCleanerMockRemoveUntil) Expect(p context.Context, p1 insolar.PulseNumber, p2 map[insolar.ID]struct {
}) *mIndexCleanerMockRemoveUntil {
	m.mock.RemoveUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCleanerMockRemoveUntilExpectation{}
	}
	m.mainExpectation.input = &IndexCleanerMockRemoveUntilInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of IndexCleaner.RemoveForPulse
func (m *mIndexCleanerMockRemoveUntil) Return() *IndexCleanerMock {
	m.mock.RemoveUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCleanerMockRemoveUntilExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of IndexCleaner.RemoveForPulse is expected once
func (m *mIndexCleanerMockRemoveUntil) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 map[insolar.ID]struct {
}) *IndexCleanerMockRemoveUntilExpectation {
	m.mock.RemoveUntilFunc = nil
	m.mainExpectation = nil

	expectation := &IndexCleanerMockRemoveUntilExpectation{}
	expectation.input = &IndexCleanerMockRemoveUntilInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of IndexCleaner.RemoveForPulse method
func (m *mIndexCleanerMockRemoveUntil) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 map[insolar.ID]struct {
})) *IndexCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveUntilFunc = f
	return m.mock
}

// RemoveForPulse implements github.com/insolar/insolar/ledger/storage/object.IndexCleaner interface
func (m *IndexCleanerMock) RemoveForPulse(ctx context.Context, pn insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.RemoveUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveUntilCounter, 1)

	if len(m.RemoveUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexCleanerMock.RemoveForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.RemoveUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexCleanerMockRemoveUntilInput{p, p1, p2}, "IndexCleaner.RemoveForPulse got unexpected parameters")

		return
	}

	if m.RemoveUntilMock.mainExpectation != nil {

		input := m.RemoveUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexCleanerMockRemoveUntilInput{p, p1, p2}, "IndexCleaner.RemoveForPulse got unexpected parameters")
		}

		return
	}

	if m.RemoveUntilFunc == nil {
		m.t.Fatalf("Unexpected call to IndexCleanerMock.RemoveForPulse. %v %v %v", p, p1, p2)
		return
	}

	m.RemoveUntilFunc(p, p1, p2)
}

// RemoveUntilMinimockCounter returns a count of IndexCleanerMock.RemoveUntilFunc invocations
func (m *IndexCleanerMock) RemoveUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveUntilCounter)
}

// RemoveUntilMinimockPreCounter returns the value of IndexCleanerMock.RemoveForPulse invocations
func (m *IndexCleanerMock) RemoveUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveUntilPreCounter)
}

// RemoveUntilFinished returns true if mock invocations count is ok
func (m *IndexCleanerMock) RemoveUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveUntilCounter) == uint64(len(m.RemoveUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveUntilCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCleanerMock) ValidateCallCounters() {

	if !m.RemoveUntilFinished() {
		m.t.Fatal("Expected call to IndexCleanerMock.RemoveForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCleanerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexCleanerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexCleanerMock) MinimockFinish() {

	if !m.RemoveUntilFinished() {
		m.t.Fatal("Expected call to IndexCleanerMock.RemoveForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexCleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexCleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RemoveUntilFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveUntilFinished() {
				m.t.Error("Expected call to IndexCleanerMock.RemoveForPulse")
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
func (m *IndexCleanerMock) AllMocksCalled() bool {

	if !m.RemoveUntilFinished() {
		return false
	}

	return true
}

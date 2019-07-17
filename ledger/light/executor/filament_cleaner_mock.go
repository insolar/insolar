package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FilamentCleaner" can be found in github.com/insolar/insolar/ledger/light/executor
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//FilamentCleanerMock implements github.com/insolar/insolar/ledger/light/executor.FilamentCleaner
type FilamentCleanerMock struct {
	t minimock.Tester

	ClearFunc       func(p insolar.ID)
	ClearCounter    uint64
	ClearPreCounter uint64
	ClearMock       mFilamentCleanerMockClear
}

//NewFilamentCleanerMock returns a mock for github.com/insolar/insolar/ledger/light/executor.FilamentCleaner
func NewFilamentCleanerMock(t minimock.Tester) *FilamentCleanerMock {
	m := &FilamentCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ClearMock = mFilamentCleanerMockClear{mock: m}

	return m
}

type mFilamentCleanerMockClear struct {
	mock              *FilamentCleanerMock
	mainExpectation   *FilamentCleanerMockClearExpectation
	expectationSeries []*FilamentCleanerMockClearExpectation
}

type FilamentCleanerMockClearExpectation struct {
	input *FilamentCleanerMockClearInput
}

type FilamentCleanerMockClearInput struct {
	p insolar.ID
}

//Expect specifies that invocation of FilamentCleaner.Clear is expected from 1 to Infinity times
func (m *mFilamentCleanerMockClear) Expect(p insolar.ID) *mFilamentCleanerMockClear {
	m.mock.ClearFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCleanerMockClearExpectation{}
	}
	m.mainExpectation.input = &FilamentCleanerMockClearInput{p}
	return m
}

//Return specifies results of invocation of FilamentCleaner.Clear
func (m *mFilamentCleanerMockClear) Return() *FilamentCleanerMock {
	m.mock.ClearFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentCleanerMockClearExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of FilamentCleaner.Clear is expected once
func (m *mFilamentCleanerMockClear) ExpectOnce(p insolar.ID) *FilamentCleanerMockClearExpectation {
	m.mock.ClearFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentCleanerMockClearExpectation{}
	expectation.input = &FilamentCleanerMockClearInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of FilamentCleaner.Clear method
func (m *mFilamentCleanerMockClear) Set(f func(p insolar.ID)) *FilamentCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ClearFunc = f
	return m.mock
}

//Clear implements github.com/insolar/insolar/ledger/light/executor.FilamentCleaner interface
func (m *FilamentCleanerMock) Clear(p insolar.ID) {
	counter := atomic.AddUint64(&m.ClearPreCounter, 1)
	defer atomic.AddUint64(&m.ClearCounter, 1)

	if len(m.ClearMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ClearMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentCleanerMock.Clear. %v", p)
			return
		}

		input := m.ClearMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentCleanerMockClearInput{p}, "FilamentCleaner.Clear got unexpected parameters")

		return
	}

	if m.ClearMock.mainExpectation != nil {

		input := m.ClearMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentCleanerMockClearInput{p}, "FilamentCleaner.Clear got unexpected parameters")
		}

		return
	}

	if m.ClearFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentCleanerMock.Clear. %v", p)
		return
	}

	m.ClearFunc(p)
}

//ClearMinimockCounter returns a count of FilamentCleanerMock.ClearFunc invocations
func (m *FilamentCleanerMock) ClearMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ClearCounter)
}

//ClearMinimockPreCounter returns the value of FilamentCleanerMock.Clear invocations
func (m *FilamentCleanerMock) ClearMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClearPreCounter)
}

//ClearFinished returns true if mock invocations count is ok
func (m *FilamentCleanerMock) ClearFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ClearMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ClearCounter) == uint64(len(m.ClearMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ClearMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ClearCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ClearFunc != nil {
		return atomic.LoadUint64(&m.ClearCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCleanerMock) ValidateCallCounters() {

	if !m.ClearFinished() {
		m.t.Fatal("Expected call to FilamentCleanerMock.Clear")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentCleanerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FilamentCleanerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FilamentCleanerMock) MinimockFinish() {

	if !m.ClearFinished() {
		m.t.Fatal("Expected call to FilamentCleanerMock.Clear")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FilamentCleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FilamentCleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ClearFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ClearFinished() {
				m.t.Error("Expected call to FilamentCleanerMock.Clear")
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
func (m *FilamentCleanerMock) AllMocksCalled() bool {

	if !m.ClearFinished() {
		return false
	}

	return true
}

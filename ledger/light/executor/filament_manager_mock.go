package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FilamentManager" can be found in github.com/insolar/insolar/ledger/light/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//FilamentManagerMock implements github.com/insolar/insolar/ledger/light/executor.FilamentManager
type FilamentManagerMock struct {
	t minimock.Tester

	SetResultFunc       func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (r *record.CompositeFilamentRecord, r1 error)
	SetResultCounter    uint64
	SetResultPreCounter uint64
	SetResultMock       mFilamentManagerMockSetResult
}

//NewFilamentManagerMock returns a mock for github.com/insolar/insolar/ledger/light/executor.FilamentManager
func NewFilamentManagerMock(t minimock.Tester) *FilamentManagerMock {
	m := &FilamentManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetResultMock = mFilamentManagerMockSetResult{mock: m}

	return m
}

type mFilamentManagerMockSetResult struct {
	mock              *FilamentManagerMock
	mainExpectation   *FilamentManagerMockSetResultExpectation
	expectationSeries []*FilamentManagerMockSetResultExpectation
}

type FilamentManagerMockSetResultExpectation struct {
	input  *FilamentManagerMockSetResultInput
	result *FilamentManagerMockSetResultResult
}

type FilamentManagerMockSetResultInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.JetID
	p3 record.Result
}

type FilamentManagerMockSetResultResult struct {
	r  *record.CompositeFilamentRecord
	r1 error
}

//Expect specifies that invocation of FilamentManager.SetResult is expected from 1 to Infinity times
func (m *mFilamentManagerMockSetResult) Expect(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) *mFilamentManagerMockSetResult {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentManagerMockSetResultExpectation{}
	}
	m.mainExpectation.input = &FilamentManagerMockSetResultInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of FilamentManager.SetResult
func (m *mFilamentManagerMockSetResult) Return(r *record.CompositeFilamentRecord, r1 error) *FilamentManagerMock {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentManagerMockSetResultExpectation{}
	}
	m.mainExpectation.result = &FilamentManagerMockSetResultResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentManager.SetResult is expected once
func (m *mFilamentManagerMockSetResult) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) *FilamentManagerMockSetResultExpectation {
	m.mock.SetResultFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentManagerMockSetResultExpectation{}
	expectation.input = &FilamentManagerMockSetResultInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentManagerMockSetResultExpectation) Return(r *record.CompositeFilamentRecord, r1 error) {
	e.result = &FilamentManagerMockSetResultResult{r, r1}
}

//Set uses given function f as a mock of FilamentManager.SetResult method
func (m *mFilamentManagerMockSetResult) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (r *record.CompositeFilamentRecord, r1 error)) *FilamentManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetResultFunc = f
	return m.mock
}

//SetResult implements github.com/insolar/insolar/ledger/light/executor.FilamentManager interface
func (m *FilamentManagerMock) SetResult(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (r *record.CompositeFilamentRecord, r1 error) {
	counter := atomic.AddUint64(&m.SetResultPreCounter, 1)
	defer atomic.AddUint64(&m.SetResultCounter, 1)

	if len(m.SetResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentManagerMock.SetResult. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentManagerMockSetResultInput{p, p1, p2, p3}, "FilamentManager.SetResult got unexpected parameters")

		result := m.SetResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentManagerMock.SetResult")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetResultMock.mainExpectation != nil {

		input := m.SetResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentManagerMockSetResultInput{p, p1, p2, p3}, "FilamentManager.SetResult got unexpected parameters")
		}

		result := m.SetResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentManagerMock.SetResult")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetResultFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentManagerMock.SetResult. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetResultFunc(p, p1, p2, p3)
}

//SetResultMinimockCounter returns a count of FilamentManagerMock.SetResultFunc invocations
func (m *FilamentManagerMock) SetResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultCounter)
}

//SetResultMinimockPreCounter returns the value of FilamentManagerMock.SetResult invocations
func (m *FilamentManagerMock) SetResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultPreCounter)
}

//SetResultFinished returns true if mock invocations count is ok
func (m *FilamentManagerMock) SetResultFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetResultMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetResultCounter) == uint64(len(m.SetResultMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetResultMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetResultCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetResultFunc != nil {
		return atomic.LoadUint64(&m.SetResultCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentManagerMock) ValidateCallCounters() {

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to FilamentManagerMock.SetResult")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FilamentManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FilamentManagerMock) MinimockFinish() {

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to FilamentManagerMock.SetResult")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FilamentManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FilamentManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetResultFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetResultFinished() {
				m.t.Error("Expected call to FilamentManagerMock.SetResult")
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
func (m *FilamentManagerMock) AllMocksCalled() bool {

	if !m.SetResultFinished() {
		return false
	}

	return true
}

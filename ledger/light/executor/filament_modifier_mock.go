package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FilamentModifier" can be found in github.com/insolar/insolar/ledger/light/executor
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

//FilamentModifierMock implements github.com/insolar/insolar/ledger/light/executor.FilamentModifier
type FilamentModifierMock struct {
	t minimock.Tester

	SetRequestFunc       func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error)
	SetRequestCounter    uint64
	SetRequestPreCounter uint64
	SetRequestMock       mFilamentModifierMockSetRequest

	SetResultFunc       func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (r error)
	SetResultCounter    uint64
	SetResultPreCounter uint64
	SetResultMock       mFilamentModifierMockSetResult
}

//NewFilamentModifierMock returns a mock for github.com/insolar/insolar/ledger/light/executor.FilamentModifier
func NewFilamentModifierMock(t minimock.Tester) *FilamentModifierMock {
	m := &FilamentModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetRequestMock = mFilamentModifierMockSetRequest{mock: m}
	m.SetResultMock = mFilamentModifierMockSetResult{mock: m}

	return m
}

type mFilamentModifierMockSetRequest struct {
	mock              *FilamentModifierMock
	mainExpectation   *FilamentModifierMockSetRequestExpectation
	expectationSeries []*FilamentModifierMockSetRequestExpectation
}

type FilamentModifierMockSetRequestExpectation struct {
	input  *FilamentModifierMockSetRequestInput
	result *FilamentModifierMockSetRequestResult
}

type FilamentModifierMockSetRequestInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.JetID
	p3 record.Request
}

type FilamentModifierMockSetRequestResult struct {
	r  *record.CompositeFilamentRecord
	r1 *record.CompositeFilamentRecord
	r2 error
}

//Expect specifies that invocation of FilamentModifier.SetRequest is expected from 1 to Infinity times
func (m *mFilamentModifierMockSetRequest) Expect(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Request) *mFilamentModifierMockSetRequest {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.input = &FilamentModifierMockSetRequestInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of FilamentModifier.SetRequest
func (m *mFilamentModifierMockSetRequest) Return(r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) *FilamentModifierMock {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.result = &FilamentModifierMockSetRequestResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentModifier.SetRequest is expected once
func (m *mFilamentModifierMockSetRequest) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Request) *FilamentModifierMockSetRequestExpectation {
	m.mock.SetRequestFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentModifierMockSetRequestExpectation{}
	expectation.input = &FilamentModifierMockSetRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentModifierMockSetRequestExpectation) Return(r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) {
	e.result = &FilamentModifierMockSetRequestResult{r, r1, r2}
}

//Set uses given function f as a mock of FilamentModifier.SetRequest method
func (m *mFilamentModifierMockSetRequest) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error)) *FilamentModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetRequestFunc = f
	return m.mock
}

//SetRequest implements github.com/insolar/insolar/ledger/light/executor.FilamentModifier interface
func (m *FilamentModifierMock) SetRequest(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) {
	counter := atomic.AddUint64(&m.SetRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SetRequestCounter, 1)

	if len(m.SetRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentModifierMockSetRequestInput{p, p1, p2, p3}, "FilamentModifier.SetRequest got unexpected parameters")

		result := m.SetRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentModifierMock.SetRequest")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.SetRequestMock.mainExpectation != nil {

		input := m.SetRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentModifierMockSetRequestInput{p, p1, p2, p3}, "FilamentModifier.SetRequest got unexpected parameters")
		}

		result := m.SetRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentModifierMock.SetRequest")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.SetRequestFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetRequestFunc(p, p1, p2, p3)
}

//SetRequestMinimockCounter returns a count of FilamentModifierMock.SetRequestFunc invocations
func (m *FilamentModifierMock) SetRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestCounter)
}

//SetRequestMinimockPreCounter returns the value of FilamentModifierMock.SetRequest invocations
func (m *FilamentModifierMock) SetRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestPreCounter)
}

//SetRequestFinished returns true if mock invocations count is ok
func (m *FilamentModifierMock) SetRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetRequestCounter) == uint64(len(m.SetRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetRequestFunc != nil {
		return atomic.LoadUint64(&m.SetRequestCounter) > 0
	}

	return true
}

type mFilamentModifierMockSetResult struct {
	mock              *FilamentModifierMock
	mainExpectation   *FilamentModifierMockSetResultExpectation
	expectationSeries []*FilamentModifierMockSetResultExpectation
}

type FilamentModifierMockSetResultExpectation struct {
	input  *FilamentModifierMockSetResultInput
	result *FilamentModifierMockSetResultResult
}

type FilamentModifierMockSetResultInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.JetID
	p3 record.Result
}

type FilamentModifierMockSetResultResult struct {
	r error
}

//Expect specifies that invocation of FilamentModifier.SetResult is expected from 1 to Infinity times
func (m *mFilamentModifierMockSetResult) Expect(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) *mFilamentModifierMockSetResult {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentModifierMockSetResultExpectation{}
	}
	m.mainExpectation.input = &FilamentModifierMockSetResultInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of FilamentModifier.SetResult
func (m *mFilamentModifierMockSetResult) Return(r error) *FilamentModifierMock {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FilamentModifierMockSetResultExpectation{}
	}
	m.mainExpectation.result = &FilamentModifierMockSetResultResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FilamentModifier.SetResult is expected once
func (m *mFilamentModifierMockSetResult) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) *FilamentModifierMockSetResultExpectation {
	m.mock.SetResultFunc = nil
	m.mainExpectation = nil

	expectation := &FilamentModifierMockSetResultExpectation{}
	expectation.input = &FilamentModifierMockSetResultInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FilamentModifierMockSetResultExpectation) Return(r error) {
	e.result = &FilamentModifierMockSetResultResult{r}
}

//Set uses given function f as a mock of FilamentModifier.SetResult method
func (m *mFilamentModifierMockSetResult) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (r error)) *FilamentModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetResultFunc = f
	return m.mock
}

//SetResult implements github.com/insolar/insolar/ledger/light/executor.FilamentModifier interface
func (m *FilamentModifierMock) SetResult(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (r error) {
	counter := atomic.AddUint64(&m.SetResultPreCounter, 1)
	defer atomic.AddUint64(&m.SetResultCounter, 1)

	if len(m.SetResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FilamentModifierMock.SetResult. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FilamentModifierMockSetResultInput{p, p1, p2, p3}, "FilamentModifier.SetResult got unexpected parameters")

		result := m.SetResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentModifierMock.SetResult")
			return
		}

		r = result.r

		return
	}

	if m.SetResultMock.mainExpectation != nil {

		input := m.SetResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FilamentModifierMockSetResultInput{p, p1, p2, p3}, "FilamentModifier.SetResult got unexpected parameters")
		}

		result := m.SetResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FilamentModifierMock.SetResult")
		}

		r = result.r

		return
	}

	if m.SetResultFunc == nil {
		m.t.Fatalf("Unexpected call to FilamentModifierMock.SetResult. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetResultFunc(p, p1, p2, p3)
}

//SetResultMinimockCounter returns a count of FilamentModifierMock.SetResultFunc invocations
func (m *FilamentModifierMock) SetResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultCounter)
}

//SetResultMinimockPreCounter returns the value of FilamentModifierMock.SetResult invocations
func (m *FilamentModifierMock) SetResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultPreCounter)
}

//SetResultFinished returns true if mock invocations count is ok
func (m *FilamentModifierMock) SetResultFinished() bool {
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
func (m *FilamentModifierMock) ValidateCallCounters() {

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to FilamentModifierMock.SetRequest")
	}

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to FilamentModifierMock.SetResult")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FilamentModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FilamentModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FilamentModifierMock) MinimockFinish() {

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to FilamentModifierMock.SetRequest")
	}

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to FilamentModifierMock.SetResult")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FilamentModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FilamentModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetRequestFinished()
		ok = ok && m.SetResultFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetRequestFinished() {
				m.t.Error("Expected call to FilamentModifierMock.SetRequest")
			}

			if !m.SetResultFinished() {
				m.t.Error("Expected call to FilamentModifierMock.SetResult")
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
func (m *FilamentModifierMock) AllMocksCalled() bool {

	if !m.SetRequestFinished() {
		return false
	}

	if !m.SetResultFinished() {
		return false
	}

	return true
}

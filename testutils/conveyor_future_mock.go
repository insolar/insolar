package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ConveyorFuture" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ConveyorFutureMock implements github.com/insolar/insolar/core.ConveyorFuture
type ConveyorFutureMock struct {
	t minimock.Tester

	CancelFunc       func()
	CancelCounter    uint64
	CancelPreCounter uint64
	CancelMock       mConveyorFutureMockCancel

	GetResultFunc       func(p time.Duration) (r core.Reply, r1 error)
	GetResultCounter    uint64
	GetResultPreCounter uint64
	GetResultMock       mConveyorFutureMockGetResult

	ResultFunc       func() (r <-chan core.Reply)
	ResultCounter    uint64
	ResultPreCounter uint64
	ResultMock       mConveyorFutureMockResult

	SetResultFunc       func(p core.Reply)
	SetResultCounter    uint64
	SetResultPreCounter uint64
	SetResultMock       mConveyorFutureMockSetResult
}

//NewConveyorFutureMock returns a mock for github.com/insolar/insolar/core.ConveyorFuture
func NewConveyorFutureMock(t minimock.Tester) *ConveyorFutureMock {
	m := &ConveyorFutureMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CancelMock = mConveyorFutureMockCancel{mock: m}
	m.GetResultMock = mConveyorFutureMockGetResult{mock: m}
	m.ResultMock = mConveyorFutureMockResult{mock: m}
	m.SetResultMock = mConveyorFutureMockSetResult{mock: m}

	return m
}

type mConveyorFutureMockCancel struct {
	mock              *ConveyorFutureMock
	mainExpectation   *ConveyorFutureMockCancelExpectation
	expectationSeries []*ConveyorFutureMockCancelExpectation
}

type ConveyorFutureMockCancelExpectation struct {
}

//Expect specifies that invocation of ConveyorFuture.Cancel is expected from 1 to Infinity times
func (m *mConveyorFutureMockCancel) Expect() *mConveyorFutureMockCancel {
	m.mock.CancelFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockCancelExpectation{}
	}

	return m
}

//Return specifies results of invocation of ConveyorFuture.Cancel
func (m *mConveyorFutureMockCancel) Return() *ConveyorFutureMock {
	m.mock.CancelFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockCancelExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ConveyorFuture.Cancel is expected once
func (m *mConveyorFutureMockCancel) ExpectOnce() *ConveyorFutureMockCancelExpectation {
	m.mock.CancelFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorFutureMockCancelExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ConveyorFuture.Cancel method
func (m *mConveyorFutureMockCancel) Set(f func()) *ConveyorFutureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CancelFunc = f
	return m.mock
}

//Cancel implements github.com/insolar/insolar/core.ConveyorFuture interface
func (m *ConveyorFutureMock) Cancel() {
	counter := atomic.AddUint64(&m.CancelPreCounter, 1)
	defer atomic.AddUint64(&m.CancelCounter, 1)

	if len(m.CancelMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CancelMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorFutureMock.Cancel.")
			return
		}

		return
	}

	if m.CancelMock.mainExpectation != nil {

		return
	}

	if m.CancelFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorFutureMock.Cancel.")
		return
	}

	m.CancelFunc()
}

//CancelMinimockCounter returns a count of ConveyorFutureMock.CancelFunc invocations
func (m *ConveyorFutureMock) CancelMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CancelCounter)
}

//CancelMinimockPreCounter returns the value of ConveyorFutureMock.Cancel invocations
func (m *ConveyorFutureMock) CancelMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CancelPreCounter)
}

//CancelFinished returns true if mock invocations count is ok
func (m *ConveyorFutureMock) CancelFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CancelMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CancelCounter) == uint64(len(m.CancelMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CancelMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CancelCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CancelFunc != nil {
		return atomic.LoadUint64(&m.CancelCounter) > 0
	}

	return true
}

type mConveyorFutureMockGetResult struct {
	mock              *ConveyorFutureMock
	mainExpectation   *ConveyorFutureMockGetResultExpectation
	expectationSeries []*ConveyorFutureMockGetResultExpectation
}

type ConveyorFutureMockGetResultExpectation struct {
	input  *ConveyorFutureMockGetResultInput
	result *ConveyorFutureMockGetResultResult
}

type ConveyorFutureMockGetResultInput struct {
	p time.Duration
}

type ConveyorFutureMockGetResultResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of ConveyorFuture.GetResult is expected from 1 to Infinity times
func (m *mConveyorFutureMockGetResult) Expect(p time.Duration) *mConveyorFutureMockGetResult {
	m.mock.GetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockGetResultExpectation{}
	}
	m.mainExpectation.input = &ConveyorFutureMockGetResultInput{p}
	return m
}

//Return specifies results of invocation of ConveyorFuture.GetResult
func (m *mConveyorFutureMockGetResult) Return(r core.Reply, r1 error) *ConveyorFutureMock {
	m.mock.GetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockGetResultExpectation{}
	}
	m.mainExpectation.result = &ConveyorFutureMockGetResultResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ConveyorFuture.GetResult is expected once
func (m *mConveyorFutureMockGetResult) ExpectOnce(p time.Duration) *ConveyorFutureMockGetResultExpectation {
	m.mock.GetResultFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorFutureMockGetResultExpectation{}
	expectation.input = &ConveyorFutureMockGetResultInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorFutureMockGetResultExpectation) Return(r core.Reply, r1 error) {
	e.result = &ConveyorFutureMockGetResultResult{r, r1}
}

//Set uses given function f as a mock of ConveyorFuture.GetResult method
func (m *mConveyorFutureMockGetResult) Set(f func(p time.Duration) (r core.Reply, r1 error)) *ConveyorFutureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetResultFunc = f
	return m.mock
}

//GetResult implements github.com/insolar/insolar/core.ConveyorFuture interface
func (m *ConveyorFutureMock) GetResult(p time.Duration) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.GetResultPreCounter, 1)
	defer atomic.AddUint64(&m.GetResultCounter, 1)

	if len(m.GetResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorFutureMock.GetResult. %v", p)
			return
		}

		input := m.GetResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConveyorFutureMockGetResultInput{p}, "ConveyorFuture.GetResult got unexpected parameters")

		result := m.GetResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorFutureMock.GetResult")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetResultMock.mainExpectation != nil {

		input := m.GetResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConveyorFutureMockGetResultInput{p}, "ConveyorFuture.GetResult got unexpected parameters")
		}

		result := m.GetResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorFutureMock.GetResult")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetResultFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorFutureMock.GetResult. %v", p)
		return
	}

	return m.GetResultFunc(p)
}

//GetResultMinimockCounter returns a count of ConveyorFutureMock.GetResultFunc invocations
func (m *ConveyorFutureMock) GetResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetResultCounter)
}

//GetResultMinimockPreCounter returns the value of ConveyorFutureMock.GetResult invocations
func (m *ConveyorFutureMock) GetResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetResultPreCounter)
}

//GetResultFinished returns true if mock invocations count is ok
func (m *ConveyorFutureMock) GetResultFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetResultMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetResultCounter) == uint64(len(m.GetResultMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetResultMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetResultCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetResultFunc != nil {
		return atomic.LoadUint64(&m.GetResultCounter) > 0
	}

	return true
}

type mConveyorFutureMockResult struct {
	mock              *ConveyorFutureMock
	mainExpectation   *ConveyorFutureMockResultExpectation
	expectationSeries []*ConveyorFutureMockResultExpectation
}

type ConveyorFutureMockResultExpectation struct {
	result *ConveyorFutureMockResultResult
}

type ConveyorFutureMockResultResult struct {
	r <-chan core.Reply
}

//Expect specifies that invocation of ConveyorFuture.Result is expected from 1 to Infinity times
func (m *mConveyorFutureMockResult) Expect() *mConveyorFutureMockResult {
	m.mock.ResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockResultExpectation{}
	}

	return m
}

//Return specifies results of invocation of ConveyorFuture.Result
func (m *mConveyorFutureMockResult) Return(r <-chan core.Reply) *ConveyorFutureMock {
	m.mock.ResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockResultExpectation{}
	}
	m.mainExpectation.result = &ConveyorFutureMockResultResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ConveyorFuture.Result is expected once
func (m *mConveyorFutureMockResult) ExpectOnce() *ConveyorFutureMockResultExpectation {
	m.mock.ResultFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorFutureMockResultExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorFutureMockResultExpectation) Return(r <-chan core.Reply) {
	e.result = &ConveyorFutureMockResultResult{r}
}

//Set uses given function f as a mock of ConveyorFuture.Result method
func (m *mConveyorFutureMockResult) Set(f func() (r <-chan core.Reply)) *ConveyorFutureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResultFunc = f
	return m.mock
}

//Result implements github.com/insolar/insolar/core.ConveyorFuture interface
func (m *ConveyorFutureMock) Result() (r <-chan core.Reply) {
	counter := atomic.AddUint64(&m.ResultPreCounter, 1)
	defer atomic.AddUint64(&m.ResultCounter, 1)

	if len(m.ResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorFutureMock.Result.")
			return
		}

		result := m.ResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorFutureMock.Result")
			return
		}

		r = result.r

		return
	}

	if m.ResultMock.mainExpectation != nil {

		result := m.ResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorFutureMock.Result")
		}

		r = result.r

		return
	}

	if m.ResultFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorFutureMock.Result.")
		return
	}

	return m.ResultFunc()
}

//ResultMinimockCounter returns a count of ConveyorFutureMock.ResultFunc invocations
func (m *ConveyorFutureMock) ResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResultCounter)
}

//ResultMinimockPreCounter returns the value of ConveyorFutureMock.Result invocations
func (m *ConveyorFutureMock) ResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResultPreCounter)
}

//ResultFinished returns true if mock invocations count is ok
func (m *ConveyorFutureMock) ResultFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ResultMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ResultCounter) == uint64(len(m.ResultMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ResultMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ResultCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ResultFunc != nil {
		return atomic.LoadUint64(&m.ResultCounter) > 0
	}

	return true
}

type mConveyorFutureMockSetResult struct {
	mock              *ConveyorFutureMock
	mainExpectation   *ConveyorFutureMockSetResultExpectation
	expectationSeries []*ConveyorFutureMockSetResultExpectation
}

type ConveyorFutureMockSetResultExpectation struct {
	input *ConveyorFutureMockSetResultInput
}

type ConveyorFutureMockSetResultInput struct {
	p core.Reply
}

//Expect specifies that invocation of ConveyorFuture.SetResult is expected from 1 to Infinity times
func (m *mConveyorFutureMockSetResult) Expect(p core.Reply) *mConveyorFutureMockSetResult {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockSetResultExpectation{}
	}
	m.mainExpectation.input = &ConveyorFutureMockSetResultInput{p}
	return m
}

//Return specifies results of invocation of ConveyorFuture.SetResult
func (m *mConveyorFutureMockSetResult) Return() *ConveyorFutureMock {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorFutureMockSetResultExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ConveyorFuture.SetResult is expected once
func (m *mConveyorFutureMockSetResult) ExpectOnce(p core.Reply) *ConveyorFutureMockSetResultExpectation {
	m.mock.SetResultFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorFutureMockSetResultExpectation{}
	expectation.input = &ConveyorFutureMockSetResultInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ConveyorFuture.SetResult method
func (m *mConveyorFutureMockSetResult) Set(f func(p core.Reply)) *ConveyorFutureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetResultFunc = f
	return m.mock
}

//SetResult implements github.com/insolar/insolar/core.ConveyorFuture interface
func (m *ConveyorFutureMock) SetResult(p core.Reply) {
	counter := atomic.AddUint64(&m.SetResultPreCounter, 1)
	defer atomic.AddUint64(&m.SetResultCounter, 1)

	if len(m.SetResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorFutureMock.SetResult. %v", p)
			return
		}

		input := m.SetResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConveyorFutureMockSetResultInput{p}, "ConveyorFuture.SetResult got unexpected parameters")

		return
	}

	if m.SetResultMock.mainExpectation != nil {

		input := m.SetResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConveyorFutureMockSetResultInput{p}, "ConveyorFuture.SetResult got unexpected parameters")
		}

		return
	}

	if m.SetResultFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorFutureMock.SetResult. %v", p)
		return
	}

	m.SetResultFunc(p)
}

//SetResultMinimockCounter returns a count of ConveyorFutureMock.SetResultFunc invocations
func (m *ConveyorFutureMock) SetResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultCounter)
}

//SetResultMinimockPreCounter returns the value of ConveyorFutureMock.SetResult invocations
func (m *ConveyorFutureMock) SetResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultPreCounter)
}

//SetResultFinished returns true if mock invocations count is ok
func (m *ConveyorFutureMock) SetResultFinished() bool {
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
func (m *ConveyorFutureMock) ValidateCallCounters() {

	if !m.CancelFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.Cancel")
	}

	if !m.GetResultFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.GetResult")
	}

	if !m.ResultFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.Result")
	}

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.SetResult")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ConveyorFutureMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ConveyorFutureMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ConveyorFutureMock) MinimockFinish() {

	if !m.CancelFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.Cancel")
	}

	if !m.GetResultFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.GetResult")
	}

	if !m.ResultFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.Result")
	}

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to ConveyorFutureMock.SetResult")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ConveyorFutureMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ConveyorFutureMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CancelFinished()
		ok = ok && m.GetResultFinished()
		ok = ok && m.ResultFinished()
		ok = ok && m.SetResultFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CancelFinished() {
				m.t.Error("Expected call to ConveyorFutureMock.Cancel")
			}

			if !m.GetResultFinished() {
				m.t.Error("Expected call to ConveyorFutureMock.GetResult")
			}

			if !m.ResultFinished() {
				m.t.Error("Expected call to ConveyorFutureMock.Result")
			}

			if !m.SetResultFinished() {
				m.t.Error("Expected call to ConveyorFutureMock.SetResult")
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
func (m *ConveyorFutureMock) AllMocksCalled() bool {

	if !m.CancelFinished() {
		return false
	}

	if !m.GetResultFinished() {
		return false
	}

	if !m.ResultFinished() {
		return false
	}

	if !m.SetResultFinished() {
		return false
	}

	return true
}

package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PendingModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//PendingModifierMock implements github.com/insolar/insolar/ledger/object.PendingModifier
type PendingModifierMock struct {
	t minimock.Tester

	RefreshStateFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)
	RefreshStateCounter    uint64
	RefreshStatePreCounter uint64
	RefreshStateMock       mPendingModifierMockRefreshState

	SetFilamentFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.PulseNumber, p4 []record.Virtual) (r error)
	SetFilamentCounter    uint64
	SetFilamentPreCounter uint64
	SetFilamentMock       mPendingModifierMockSetFilament

	SetRequestFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Request) (r error)
	SetRequestCounter    uint64
	SetRequestPreCounter uint64
	SetRequestMock       mPendingModifierMockSetRequest

	SetResultFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Result) (r error)
	SetResultCounter    uint64
	SetResultPreCounter uint64
	SetResultMock       mPendingModifierMockSetResult
}

//NewPendingModifierMock returns a mock for github.com/insolar/insolar/ledger/object.PendingModifier
func NewPendingModifierMock(t minimock.Tester) *PendingModifierMock {
	m := &PendingModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RefreshStateMock = mPendingModifierMockRefreshState{mock: m}
	m.SetFilamentMock = mPendingModifierMockSetFilament{mock: m}
	m.SetRequestMock = mPendingModifierMockSetRequest{mock: m}
	m.SetResultMock = mPendingModifierMockSetResult{mock: m}

	return m
}

type mPendingModifierMockRefreshState struct {
	mock              *PendingModifierMock
	mainExpectation   *PendingModifierMockRefreshStateExpectation
	expectationSeries []*PendingModifierMockRefreshStateExpectation
}

type PendingModifierMockRefreshStateExpectation struct {
	input  *PendingModifierMockRefreshStateInput
	result *PendingModifierMockRefreshStateResult
}

type PendingModifierMockRefreshStateInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type PendingModifierMockRefreshStateResult struct {
	r error
}

//Expect specifies that invocation of PendingModifier.RefreshState is expected from 1 to Infinity times
func (m *mPendingModifierMockRefreshState) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mPendingModifierMockRefreshState {
	m.mock.RefreshStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockRefreshStateExpectation{}
	}
	m.mainExpectation.input = &PendingModifierMockRefreshStateInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of PendingModifier.RefreshState
func (m *mPendingModifierMockRefreshState) Return(r error) *PendingModifierMock {
	m.mock.RefreshStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockRefreshStateExpectation{}
	}
	m.mainExpectation.result = &PendingModifierMockRefreshStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingModifier.RefreshState is expected once
func (m *mPendingModifierMockRefreshState) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *PendingModifierMockRefreshStateExpectation {
	m.mock.RefreshStateFunc = nil
	m.mainExpectation = nil

	expectation := &PendingModifierMockRefreshStateExpectation{}
	expectation.input = &PendingModifierMockRefreshStateInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingModifierMockRefreshStateExpectation) Return(r error) {
	e.result = &PendingModifierMockRefreshStateResult{r}
}

//Set uses given function f as a mock of PendingModifier.RefreshState method
func (m *mPendingModifierMockRefreshState) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error)) *PendingModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RefreshStateFunc = f
	return m.mock
}

//RefreshState implements github.com/insolar/insolar/ledger/object.PendingModifier interface
func (m *PendingModifierMock) RefreshState(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.RefreshStatePreCounter, 1)
	defer atomic.AddUint64(&m.RefreshStateCounter, 1)

	if len(m.RefreshStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RefreshStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingModifierMock.RefreshState. %v %v %v", p, p1, p2)
			return
		}

		input := m.RefreshStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingModifierMockRefreshStateInput{p, p1, p2}, "PendingModifier.RefreshState got unexpected parameters")

		result := m.RefreshStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.RefreshState")
			return
		}

		r = result.r

		return
	}

	if m.RefreshStateMock.mainExpectation != nil {

		input := m.RefreshStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingModifierMockRefreshStateInput{p, p1, p2}, "PendingModifier.RefreshState got unexpected parameters")
		}

		result := m.RefreshStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.RefreshState")
		}

		r = result.r

		return
	}

	if m.RefreshStateFunc == nil {
		m.t.Fatalf("Unexpected call to PendingModifierMock.RefreshState. %v %v %v", p, p1, p2)
		return
	}

	return m.RefreshStateFunc(p, p1, p2)
}

//RefreshStateMinimockCounter returns a count of PendingModifierMock.RefreshStateFunc invocations
func (m *PendingModifierMock) RefreshStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RefreshStateCounter)
}

//RefreshStateMinimockPreCounter returns the value of PendingModifierMock.RefreshState invocations
func (m *PendingModifierMock) RefreshStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RefreshStatePreCounter)
}

//RefreshStateFinished returns true if mock invocations count is ok
func (m *PendingModifierMock) RefreshStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RefreshStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RefreshStateCounter) == uint64(len(m.RefreshStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RefreshStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RefreshStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RefreshStateFunc != nil {
		return atomic.LoadUint64(&m.RefreshStateCounter) > 0
	}

	return true
}

type mPendingModifierMockSetFilament struct {
	mock              *PendingModifierMock
	mainExpectation   *PendingModifierMockSetFilamentExpectation
	expectationSeries []*PendingModifierMockSetFilamentExpectation
}

type PendingModifierMockSetFilamentExpectation struct {
	input  *PendingModifierMockSetFilamentInput
	result *PendingModifierMockSetFilamentResult
}

type PendingModifierMockSetFilamentInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.PulseNumber
	p4 []record.Virtual
}

type PendingModifierMockSetFilamentResult struct {
	r error
}

//Expect specifies that invocation of PendingModifier.SetFilament is expected from 1 to Infinity times
func (m *mPendingModifierMockSetFilament) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.PulseNumber, p4 []record.Virtual) *mPendingModifierMockSetFilament {
	m.mock.SetFilamentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetFilamentExpectation{}
	}
	m.mainExpectation.input = &PendingModifierMockSetFilamentInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of PendingModifier.SetFilament
func (m *mPendingModifierMockSetFilament) Return(r error) *PendingModifierMock {
	m.mock.SetFilamentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetFilamentExpectation{}
	}
	m.mainExpectation.result = &PendingModifierMockSetFilamentResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingModifier.SetFilament is expected once
func (m *mPendingModifierMockSetFilament) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.PulseNumber, p4 []record.Virtual) *PendingModifierMockSetFilamentExpectation {
	m.mock.SetFilamentFunc = nil
	m.mainExpectation = nil

	expectation := &PendingModifierMockSetFilamentExpectation{}
	expectation.input = &PendingModifierMockSetFilamentInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingModifierMockSetFilamentExpectation) Return(r error) {
	e.result = &PendingModifierMockSetFilamentResult{r}
}

//Set uses given function f as a mock of PendingModifier.SetFilament method
func (m *mPendingModifierMockSetFilament) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.PulseNumber, p4 []record.Virtual) (r error)) *PendingModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFilamentFunc = f
	return m.mock
}

//SetFilament implements github.com/insolar/insolar/ledger/object.PendingModifier interface
func (m *PendingModifierMock) SetFilament(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.PulseNumber, p4 []record.Virtual) (r error) {
	counter := atomic.AddUint64(&m.SetFilamentPreCounter, 1)
	defer atomic.AddUint64(&m.SetFilamentCounter, 1)

	if len(m.SetFilamentMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetFilamentMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingModifierMock.SetFilament. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.SetFilamentMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingModifierMockSetFilamentInput{p, p1, p2, p3, p4}, "PendingModifier.SetFilament got unexpected parameters")

		result := m.SetFilamentMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetFilament")
			return
		}

		r = result.r

		return
	}

	if m.SetFilamentMock.mainExpectation != nil {

		input := m.SetFilamentMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingModifierMockSetFilamentInput{p, p1, p2, p3, p4}, "PendingModifier.SetFilament got unexpected parameters")
		}

		result := m.SetFilamentMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetFilament")
		}

		r = result.r

		return
	}

	if m.SetFilamentFunc == nil {
		m.t.Fatalf("Unexpected call to PendingModifierMock.SetFilament. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.SetFilamentFunc(p, p1, p2, p3, p4)
}

//SetFilamentMinimockCounter returns a count of PendingModifierMock.SetFilamentFunc invocations
func (m *PendingModifierMock) SetFilamentMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetFilamentCounter)
}

//SetFilamentMinimockPreCounter returns the value of PendingModifierMock.SetFilament invocations
func (m *PendingModifierMock) SetFilamentMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetFilamentPreCounter)
}

//SetFilamentFinished returns true if mock invocations count is ok
func (m *PendingModifierMock) SetFilamentFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetFilamentMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetFilamentCounter) == uint64(len(m.SetFilamentMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetFilamentMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetFilamentCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFilamentFunc != nil {
		return atomic.LoadUint64(&m.SetFilamentCounter) > 0
	}

	return true
}

type mPendingModifierMockSetRequest struct {
	mock              *PendingModifierMock
	mainExpectation   *PendingModifierMockSetRequestExpectation
	expectationSeries []*PendingModifierMockSetRequestExpectation
}

type PendingModifierMockSetRequestExpectation struct {
	input  *PendingModifierMockSetRequestInput
	result *PendingModifierMockSetRequestResult
}

type PendingModifierMockSetRequestInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 record.Request
}

type PendingModifierMockSetRequestResult struct {
	r error
}

//Expect specifies that invocation of PendingModifier.SetRequest is expected from 1 to Infinity times
func (m *mPendingModifierMockSetRequest) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Request) *mPendingModifierMockSetRequest {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.input = &PendingModifierMockSetRequestInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of PendingModifier.SetRequest
func (m *mPendingModifierMockSetRequest) Return(r error) *PendingModifierMock {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.result = &PendingModifierMockSetRequestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingModifier.SetRequest is expected once
func (m *mPendingModifierMockSetRequest) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Request) *PendingModifierMockSetRequestExpectation {
	m.mock.SetRequestFunc = nil
	m.mainExpectation = nil

	expectation := &PendingModifierMockSetRequestExpectation{}
	expectation.input = &PendingModifierMockSetRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingModifierMockSetRequestExpectation) Return(r error) {
	e.result = &PendingModifierMockSetRequestResult{r}
}

//Set uses given function f as a mock of PendingModifier.SetRequest method
func (m *mPendingModifierMockSetRequest) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Request) (r error)) *PendingModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetRequestFunc = f
	return m.mock
}

//SetRequest implements github.com/insolar/insolar/ledger/object.PendingModifier interface
func (m *PendingModifierMock) SetRequest(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Request) (r error) {
	counter := atomic.AddUint64(&m.SetRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SetRequestCounter, 1)

	if len(m.SetRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingModifierMockSetRequestInput{p, p1, p2, p3}, "PendingModifier.SetRequest got unexpected parameters")

		result := m.SetRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetRequest")
			return
		}

		r = result.r

		return
	}

	if m.SetRequestMock.mainExpectation != nil {

		input := m.SetRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingModifierMockSetRequestInput{p, p1, p2, p3}, "PendingModifier.SetRequest got unexpected parameters")
		}

		result := m.SetRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetRequest")
		}

		r = result.r

		return
	}

	if m.SetRequestFunc == nil {
		m.t.Fatalf("Unexpected call to PendingModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetRequestFunc(p, p1, p2, p3)
}

//SetRequestMinimockCounter returns a count of PendingModifierMock.SetRequestFunc invocations
func (m *PendingModifierMock) SetRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestCounter)
}

//SetRequestMinimockPreCounter returns the value of PendingModifierMock.SetRequest invocations
func (m *PendingModifierMock) SetRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestPreCounter)
}

//SetRequestFinished returns true if mock invocations count is ok
func (m *PendingModifierMock) SetRequestFinished() bool {
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

type mPendingModifierMockSetResult struct {
	mock              *PendingModifierMock
	mainExpectation   *PendingModifierMockSetResultExpectation
	expectationSeries []*PendingModifierMockSetResultExpectation
}

type PendingModifierMockSetResultExpectation struct {
	input  *PendingModifierMockSetResultInput
	result *PendingModifierMockSetResultResult
}

type PendingModifierMockSetResultInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 record.Result
}

type PendingModifierMockSetResultResult struct {
	r error
}

//Expect specifies that invocation of PendingModifier.SetResult is expected from 1 to Infinity times
func (m *mPendingModifierMockSetResult) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Result) *mPendingModifierMockSetResult {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetResultExpectation{}
	}
	m.mainExpectation.input = &PendingModifierMockSetResultInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of PendingModifier.SetResult
func (m *mPendingModifierMockSetResult) Return(r error) *PendingModifierMock {
	m.mock.SetResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PendingModifierMockSetResultExpectation{}
	}
	m.mainExpectation.result = &PendingModifierMockSetResultResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PendingModifier.SetResult is expected once
func (m *mPendingModifierMockSetResult) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Result) *PendingModifierMockSetResultExpectation {
	m.mock.SetResultFunc = nil
	m.mainExpectation = nil

	expectation := &PendingModifierMockSetResultExpectation{}
	expectation.input = &PendingModifierMockSetResultInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PendingModifierMockSetResultExpectation) Return(r error) {
	e.result = &PendingModifierMockSetResultResult{r}
}

//Set uses given function f as a mock of PendingModifier.SetResult method
func (m *mPendingModifierMockSetResult) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Result) (r error)) *PendingModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetResultFunc = f
	return m.mock
}

//SetResult implements github.com/insolar/insolar/ledger/object.PendingModifier interface
func (m *PendingModifierMock) SetResult(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 record.Result) (r error) {
	counter := atomic.AddUint64(&m.SetResultPreCounter, 1)
	defer atomic.AddUint64(&m.SetResultCounter, 1)

	if len(m.SetResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PendingModifierMock.SetResult. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PendingModifierMockSetResultInput{p, p1, p2, p3}, "PendingModifier.SetResult got unexpected parameters")

		result := m.SetResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetResult")
			return
		}

		r = result.r

		return
	}

	if m.SetResultMock.mainExpectation != nil {

		input := m.SetResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PendingModifierMockSetResultInput{p, p1, p2, p3}, "PendingModifier.SetResult got unexpected parameters")
		}

		result := m.SetResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PendingModifierMock.SetResult")
		}

		r = result.r

		return
	}

	if m.SetResultFunc == nil {
		m.t.Fatalf("Unexpected call to PendingModifierMock.SetResult. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetResultFunc(p, p1, p2, p3)
}

//SetResultMinimockCounter returns a count of PendingModifierMock.SetResultFunc invocations
func (m *PendingModifierMock) SetResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultCounter)
}

//SetResultMinimockPreCounter returns the value of PendingModifierMock.SetResult invocations
func (m *PendingModifierMock) SetResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultPreCounter)
}

//SetResultFinished returns true if mock invocations count is ok
func (m *PendingModifierMock) SetResultFinished() bool {
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
func (m *PendingModifierMock) ValidateCallCounters() {

	if !m.RefreshStateFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.RefreshState")
	}

	if !m.SetFilamentFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetFilament")
	}

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetRequest")
	}

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetResult")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PendingModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PendingModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PendingModifierMock) MinimockFinish() {

	if !m.RefreshStateFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.RefreshState")
	}

	if !m.SetFilamentFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetFilament")
	}

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetRequest")
	}

	if !m.SetResultFinished() {
		m.t.Fatal("Expected call to PendingModifierMock.SetResult")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PendingModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PendingModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RefreshStateFinished()
		ok = ok && m.SetFilamentFinished()
		ok = ok && m.SetRequestFinished()
		ok = ok && m.SetResultFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RefreshStateFinished() {
				m.t.Error("Expected call to PendingModifierMock.RefreshState")
			}

			if !m.SetFilamentFinished() {
				m.t.Error("Expected call to PendingModifierMock.SetFilament")
			}

			if !m.SetRequestFinished() {
				m.t.Error("Expected call to PendingModifierMock.SetRequest")
			}

			if !m.SetResultFinished() {
				m.t.Error("Expected call to PendingModifierMock.SetResult")
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
func (m *PendingModifierMock) AllMocksCalled() bool {

	if !m.RefreshStateFinished() {
		return false
	}

	if !m.SetFilamentFinished() {
		return false
	}

	if !m.SetRequestFinished() {
		return false
	}

	if !m.SetResultFinished() {
		return false
	}

	return true
}

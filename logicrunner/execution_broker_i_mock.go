package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ExecutionBrokerI" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ExecutionBrokerIMock implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI
type ExecutionBrokerIMock struct {
	t minimock.Tester

	AbandonedRequestsOnLedgerFunc       func(p context.Context)
	AbandonedRequestsOnLedgerCounter    uint64
	AbandonedRequestsOnLedgerPreCounter uint64
	AbandonedRequestsOnLedgerMock       mExecutionBrokerIMockAbandonedRequestsOnLedger

	AddAdditionalRequestFromPrevExecutorFunc       func(p context.Context, p1 *Transcript)
	AddAdditionalRequestFromPrevExecutorCounter    uint64
	AddAdditionalRequestFromPrevExecutorPreCounter uint64
	AddAdditionalRequestFromPrevExecutorMock       mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor

	AddFreshRequestFunc       func(p context.Context, p1 *Transcript)
	AddFreshRequestCounter    uint64
	AddFreshRequestPreCounter uint64
	AddFreshRequestMock       mExecutionBrokerIMockAddFreshRequest

	AddRequestsFromLedgerFunc       func(p context.Context, p1 ...*Transcript)
	AddRequestsFromLedgerCounter    uint64
	AddRequestsFromLedgerPreCounter uint64
	AddRequestsFromLedgerMock       mExecutionBrokerIMockAddRequestsFromLedger

	AddRequestsFromPrevExecutorFunc       func(p context.Context, p1 ...*Transcript)
	AddRequestsFromPrevExecutorCounter    uint64
	AddRequestsFromPrevExecutorPreCounter uint64
	AddRequestsFromPrevExecutorMock       mExecutionBrokerIMockAddRequestsFromPrevExecutor

	CheckExecutionLoopFunc       func(p context.Context, p1 string) (r bool)
	CheckExecutionLoopCounter    uint64
	CheckExecutionLoopPreCounter uint64
	CheckExecutionLoopMock       mExecutionBrokerIMockCheckExecutionLoop

	FetchMoreRequestsFromLedgerFunc       func(p context.Context)
	FetchMoreRequestsFromLedgerCounter    uint64
	FetchMoreRequestsFromLedgerPreCounter uint64
	FetchMoreRequestsFromLedgerMock       mExecutionBrokerIMockFetchMoreRequestsFromLedger

	GetActiveTranscriptFunc       func(p insolar.Reference) (r *Transcript)
	GetActiveTranscriptCounter    uint64
	GetActiveTranscriptPreCounter uint64
	GetActiveTranscriptMock       mExecutionBrokerIMockGetActiveTranscript

	IsKnownRequestFunc       func(p context.Context, p1 insolar.Reference) (r bool)
	IsKnownRequestCounter    uint64
	IsKnownRequestPreCounter uint64
	IsKnownRequestMock       mExecutionBrokerIMockIsKnownRequest

	MoreRequestsOnLedgerFunc       func(p context.Context)
	MoreRequestsOnLedgerCounter    uint64
	MoreRequestsOnLedgerPreCounter uint64
	MoreRequestsOnLedgerMock       mExecutionBrokerIMockMoreRequestsOnLedger

	NoMoreRequestsOnLedgerFunc       func(p context.Context)
	NoMoreRequestsOnLedgerCounter    uint64
	NoMoreRequestsOnLedgerPreCounter uint64
	NoMoreRequestsOnLedgerMock       mExecutionBrokerIMockNoMoreRequestsOnLedger

	OnPulseFunc       func(p context.Context, p1 bool) (r bool, r1 []insolar.Message)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mExecutionBrokerIMockOnPulse

	PendingStateFunc       func() (r insolar.PendingState)
	PendingStateCounter    uint64
	PendingStatePreCounter uint64
	PendingStateMock       mExecutionBrokerIMockPendingState

	PrevExecutorFinishedPendingFunc       func(p context.Context) (r error)
	PrevExecutorFinishedPendingCounter    uint64
	PrevExecutorFinishedPendingPreCounter uint64
	PrevExecutorFinishedPendingMock       mExecutionBrokerIMockPrevExecutorFinishedPending

	PrevExecutorPendingResultFunc       func(p context.Context, p1 insolar.PendingState)
	PrevExecutorPendingResultCounter    uint64
	PrevExecutorPendingResultPreCounter uint64
	PrevExecutorPendingResultMock       mExecutionBrokerIMockPrevExecutorPendingResult

	PrevExecutorStillExecutingFunc       func(p context.Context)
	PrevExecutorStillExecutingCounter    uint64
	PrevExecutorStillExecutingPreCounter uint64
	PrevExecutorStillExecutingMock       mExecutionBrokerIMockPrevExecutorStillExecuting

	SetNotPendingFunc       func(p context.Context)
	SetNotPendingCounter    uint64
	SetNotPendingPreCounter uint64
	SetNotPendingMock       mExecutionBrokerIMockSetNotPending
}

//NewExecutionBrokerIMock returns a mock for github.com/insolar/insolar/logicrunner.ExecutionBrokerI
func NewExecutionBrokerIMock(t minimock.Tester) *ExecutionBrokerIMock {
	m := &ExecutionBrokerIMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AbandonedRequestsOnLedgerMock = mExecutionBrokerIMockAbandonedRequestsOnLedger{mock: m}
	m.AddAdditionalRequestFromPrevExecutorMock = mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor{mock: m}
	m.AddFreshRequestMock = mExecutionBrokerIMockAddFreshRequest{mock: m}
	m.AddRequestsFromLedgerMock = mExecutionBrokerIMockAddRequestsFromLedger{mock: m}
	m.AddRequestsFromPrevExecutorMock = mExecutionBrokerIMockAddRequestsFromPrevExecutor{mock: m}
	m.CheckExecutionLoopMock = mExecutionBrokerIMockCheckExecutionLoop{mock: m}
	m.FetchMoreRequestsFromLedgerMock = mExecutionBrokerIMockFetchMoreRequestsFromLedger{mock: m}
	m.GetActiveTranscriptMock = mExecutionBrokerIMockGetActiveTranscript{mock: m}
	m.IsKnownRequestMock = mExecutionBrokerIMockIsKnownRequest{mock: m}
	m.MoreRequestsOnLedgerMock = mExecutionBrokerIMockMoreRequestsOnLedger{mock: m}
	m.NoMoreRequestsOnLedgerMock = mExecutionBrokerIMockNoMoreRequestsOnLedger{mock: m}
	m.OnPulseMock = mExecutionBrokerIMockOnPulse{mock: m}
	m.PendingStateMock = mExecutionBrokerIMockPendingState{mock: m}
	m.PrevExecutorFinishedPendingMock = mExecutionBrokerIMockPrevExecutorFinishedPending{mock: m}
	m.PrevExecutorPendingResultMock = mExecutionBrokerIMockPrevExecutorPendingResult{mock: m}
	m.PrevExecutorStillExecutingMock = mExecutionBrokerIMockPrevExecutorStillExecuting{mock: m}
	m.SetNotPendingMock = mExecutionBrokerIMockSetNotPending{mock: m}

	return m
}

type mExecutionBrokerIMockAbandonedRequestsOnLedger struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockAbandonedRequestsOnLedgerExpectation
	expectationSeries []*ExecutionBrokerIMockAbandonedRequestsOnLedgerExpectation
}

type ExecutionBrokerIMockAbandonedRequestsOnLedgerExpectation struct {
	input *ExecutionBrokerIMockAbandonedRequestsOnLedgerInput
}

type ExecutionBrokerIMockAbandonedRequestsOnLedgerInput struct {
	p context.Context
}

//Expect specifies that invocation of ExecutionBrokerI.AbandonedRequestsOnLedger is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockAbandonedRequestsOnLedger) Expect(p context.Context) *mExecutionBrokerIMockAbandonedRequestsOnLedger {
	m.mock.AbandonedRequestsOnLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAbandonedRequestsOnLedgerExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockAbandonedRequestsOnLedgerInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.AbandonedRequestsOnLedger
func (m *mExecutionBrokerIMockAbandonedRequestsOnLedger) Return() *ExecutionBrokerIMock {
	m.mock.AbandonedRequestsOnLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAbandonedRequestsOnLedgerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.AbandonedRequestsOnLedger is expected once
func (m *mExecutionBrokerIMockAbandonedRequestsOnLedger) ExpectOnce(p context.Context) *ExecutionBrokerIMockAbandonedRequestsOnLedgerExpectation {
	m.mock.AbandonedRequestsOnLedgerFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockAbandonedRequestsOnLedgerExpectation{}
	expectation.input = &ExecutionBrokerIMockAbandonedRequestsOnLedgerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.AbandonedRequestsOnLedger method
func (m *mExecutionBrokerIMockAbandonedRequestsOnLedger) Set(f func(p context.Context)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AbandonedRequestsOnLedgerFunc = f
	return m.mock
}

//AbandonedRequestsOnLedger implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) AbandonedRequestsOnLedger(p context.Context) {
	counter := atomic.AddUint64(&m.AbandonedRequestsOnLedgerPreCounter, 1)
	defer atomic.AddUint64(&m.AbandonedRequestsOnLedgerCounter, 1)

	if len(m.AbandonedRequestsOnLedgerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AbandonedRequestsOnLedgerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AbandonedRequestsOnLedger. %v", p)
			return
		}

		input := m.AbandonedRequestsOnLedgerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAbandonedRequestsOnLedgerInput{p}, "ExecutionBrokerI.AbandonedRequestsOnLedger got unexpected parameters")

		return
	}

	if m.AbandonedRequestsOnLedgerMock.mainExpectation != nil {

		input := m.AbandonedRequestsOnLedgerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAbandonedRequestsOnLedgerInput{p}, "ExecutionBrokerI.AbandonedRequestsOnLedger got unexpected parameters")
		}

		return
	}

	if m.AbandonedRequestsOnLedgerFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AbandonedRequestsOnLedger. %v", p)
		return
	}

	m.AbandonedRequestsOnLedgerFunc(p)
}

//AbandonedRequestsOnLedgerMinimockCounter returns a count of ExecutionBrokerIMock.AbandonedRequestsOnLedgerFunc invocations
func (m *ExecutionBrokerIMock) AbandonedRequestsOnLedgerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AbandonedRequestsOnLedgerCounter)
}

//AbandonedRequestsOnLedgerMinimockPreCounter returns the value of ExecutionBrokerIMock.AbandonedRequestsOnLedger invocations
func (m *ExecutionBrokerIMock) AbandonedRequestsOnLedgerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AbandonedRequestsOnLedgerPreCounter)
}

//AbandonedRequestsOnLedgerFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) AbandonedRequestsOnLedgerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AbandonedRequestsOnLedgerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AbandonedRequestsOnLedgerCounter) == uint64(len(m.AbandonedRequestsOnLedgerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AbandonedRequestsOnLedgerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AbandonedRequestsOnLedgerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AbandonedRequestsOnLedgerFunc != nil {
		return atomic.LoadUint64(&m.AbandonedRequestsOnLedgerCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorExpectation
	expectationSeries []*ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorExpectation
}

type ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorExpectation struct {
	input *ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorInput
}

type ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorInput struct {
	p  context.Context
	p1 *Transcript
}

//Expect specifies that invocation of ExecutionBrokerI.AddAdditionalRequestFromPrevExecutor is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor) Expect(p context.Context, p1 *Transcript) *mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor {
	m.mock.AddAdditionalRequestFromPrevExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.AddAdditionalRequestFromPrevExecutor
func (m *mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor) Return() *ExecutionBrokerIMock {
	m.mock.AddAdditionalRequestFromPrevExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.AddAdditionalRequestFromPrevExecutor is expected once
func (m *mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor) ExpectOnce(p context.Context, p1 *Transcript) *ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorExpectation {
	m.mock.AddAdditionalRequestFromPrevExecutorFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorExpectation{}
	expectation.input = &ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.AddAdditionalRequestFromPrevExecutor method
func (m *mExecutionBrokerIMockAddAdditionalRequestFromPrevExecutor) Set(f func(p context.Context, p1 *Transcript)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddAdditionalRequestFromPrevExecutorFunc = f
	return m.mock
}

//AddAdditionalRequestFromPrevExecutor implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) AddAdditionalRequestFromPrevExecutor(p context.Context, p1 *Transcript) {
	counter := atomic.AddUint64(&m.AddAdditionalRequestFromPrevExecutorPreCounter, 1)
	defer atomic.AddUint64(&m.AddAdditionalRequestFromPrevExecutorCounter, 1)

	if len(m.AddAdditionalRequestFromPrevExecutorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddAdditionalRequestFromPrevExecutorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddAdditionalRequestFromPrevExecutor. %v %v", p, p1)
			return
		}

		input := m.AddAdditionalRequestFromPrevExecutorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorInput{p, p1}, "ExecutionBrokerI.AddAdditionalRequestFromPrevExecutor got unexpected parameters")

		return
	}

	if m.AddAdditionalRequestFromPrevExecutorMock.mainExpectation != nil {

		input := m.AddAdditionalRequestFromPrevExecutorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddAdditionalRequestFromPrevExecutorInput{p, p1}, "ExecutionBrokerI.AddAdditionalRequestFromPrevExecutor got unexpected parameters")
		}

		return
	}

	if m.AddAdditionalRequestFromPrevExecutorFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddAdditionalRequestFromPrevExecutor. %v %v", p, p1)
		return
	}

	m.AddAdditionalRequestFromPrevExecutorFunc(p, p1)
}

//AddAdditionalRequestFromPrevExecutorMinimockCounter returns a count of ExecutionBrokerIMock.AddAdditionalRequestFromPrevExecutorFunc invocations
func (m *ExecutionBrokerIMock) AddAdditionalRequestFromPrevExecutorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddAdditionalRequestFromPrevExecutorCounter)
}

//AddAdditionalRequestFromPrevExecutorMinimockPreCounter returns the value of ExecutionBrokerIMock.AddAdditionalRequestFromPrevExecutor invocations
func (m *ExecutionBrokerIMock) AddAdditionalRequestFromPrevExecutorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddAdditionalRequestFromPrevExecutorPreCounter)
}

//AddAdditionalRequestFromPrevExecutorFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) AddAdditionalRequestFromPrevExecutorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddAdditionalRequestFromPrevExecutorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddAdditionalRequestFromPrevExecutorCounter) == uint64(len(m.AddAdditionalRequestFromPrevExecutorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddAdditionalRequestFromPrevExecutorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddAdditionalRequestFromPrevExecutorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddAdditionalRequestFromPrevExecutorFunc != nil {
		return atomic.LoadUint64(&m.AddAdditionalRequestFromPrevExecutorCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockAddFreshRequest struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockAddFreshRequestExpectation
	expectationSeries []*ExecutionBrokerIMockAddFreshRequestExpectation
}

type ExecutionBrokerIMockAddFreshRequestExpectation struct {
	input *ExecutionBrokerIMockAddFreshRequestInput
}

type ExecutionBrokerIMockAddFreshRequestInput struct {
	p  context.Context
	p1 *Transcript
}

//Expect specifies that invocation of ExecutionBrokerI.AddFreshRequest is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockAddFreshRequest) Expect(p context.Context, p1 *Transcript) *mExecutionBrokerIMockAddFreshRequest {
	m.mock.AddFreshRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddFreshRequestExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockAddFreshRequestInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.AddFreshRequest
func (m *mExecutionBrokerIMockAddFreshRequest) Return() *ExecutionBrokerIMock {
	m.mock.AddFreshRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddFreshRequestExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.AddFreshRequest is expected once
func (m *mExecutionBrokerIMockAddFreshRequest) ExpectOnce(p context.Context, p1 *Transcript) *ExecutionBrokerIMockAddFreshRequestExpectation {
	m.mock.AddFreshRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockAddFreshRequestExpectation{}
	expectation.input = &ExecutionBrokerIMockAddFreshRequestInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.AddFreshRequest method
func (m *mExecutionBrokerIMockAddFreshRequest) Set(f func(p context.Context, p1 *Transcript)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddFreshRequestFunc = f
	return m.mock
}

//AddFreshRequest implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) AddFreshRequest(p context.Context, p1 *Transcript) {
	counter := atomic.AddUint64(&m.AddFreshRequestPreCounter, 1)
	defer atomic.AddUint64(&m.AddFreshRequestCounter, 1)

	if len(m.AddFreshRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddFreshRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddFreshRequest. %v %v", p, p1)
			return
		}

		input := m.AddFreshRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddFreshRequestInput{p, p1}, "ExecutionBrokerI.AddFreshRequest got unexpected parameters")

		return
	}

	if m.AddFreshRequestMock.mainExpectation != nil {

		input := m.AddFreshRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddFreshRequestInput{p, p1}, "ExecutionBrokerI.AddFreshRequest got unexpected parameters")
		}

		return
	}

	if m.AddFreshRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddFreshRequest. %v %v", p, p1)
		return
	}

	m.AddFreshRequestFunc(p, p1)
}

//AddFreshRequestMinimockCounter returns a count of ExecutionBrokerIMock.AddFreshRequestFunc invocations
func (m *ExecutionBrokerIMock) AddFreshRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddFreshRequestCounter)
}

//AddFreshRequestMinimockPreCounter returns the value of ExecutionBrokerIMock.AddFreshRequest invocations
func (m *ExecutionBrokerIMock) AddFreshRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddFreshRequestPreCounter)
}

//AddFreshRequestFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) AddFreshRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddFreshRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddFreshRequestCounter) == uint64(len(m.AddFreshRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddFreshRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddFreshRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddFreshRequestFunc != nil {
		return atomic.LoadUint64(&m.AddFreshRequestCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockAddRequestsFromLedger struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockAddRequestsFromLedgerExpectation
	expectationSeries []*ExecutionBrokerIMockAddRequestsFromLedgerExpectation
}

type ExecutionBrokerIMockAddRequestsFromLedgerExpectation struct {
	input *ExecutionBrokerIMockAddRequestsFromLedgerInput
}

type ExecutionBrokerIMockAddRequestsFromLedgerInput struct {
	p  context.Context
	p1 []*Transcript
}

//Expect specifies that invocation of ExecutionBrokerI.AddRequestsFromLedger is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockAddRequestsFromLedger) Expect(p context.Context, p1 ...*Transcript) *mExecutionBrokerIMockAddRequestsFromLedger {
	m.mock.AddRequestsFromLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddRequestsFromLedgerExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockAddRequestsFromLedgerInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.AddRequestsFromLedger
func (m *mExecutionBrokerIMockAddRequestsFromLedger) Return() *ExecutionBrokerIMock {
	m.mock.AddRequestsFromLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddRequestsFromLedgerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.AddRequestsFromLedger is expected once
func (m *mExecutionBrokerIMockAddRequestsFromLedger) ExpectOnce(p context.Context, p1 ...*Transcript) *ExecutionBrokerIMockAddRequestsFromLedgerExpectation {
	m.mock.AddRequestsFromLedgerFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockAddRequestsFromLedgerExpectation{}
	expectation.input = &ExecutionBrokerIMockAddRequestsFromLedgerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.AddRequestsFromLedger method
func (m *mExecutionBrokerIMockAddRequestsFromLedger) Set(f func(p context.Context, p1 ...*Transcript)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddRequestsFromLedgerFunc = f
	return m.mock
}

//AddRequestsFromLedger implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) AddRequestsFromLedger(p context.Context, p1 ...*Transcript) {
	counter := atomic.AddUint64(&m.AddRequestsFromLedgerPreCounter, 1)
	defer atomic.AddUint64(&m.AddRequestsFromLedgerCounter, 1)

	if len(m.AddRequestsFromLedgerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddRequestsFromLedgerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddRequestsFromLedger. %v %v", p, p1)
			return
		}

		input := m.AddRequestsFromLedgerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddRequestsFromLedgerInput{p, p1}, "ExecutionBrokerI.AddRequestsFromLedger got unexpected parameters")

		return
	}

	if m.AddRequestsFromLedgerMock.mainExpectation != nil {

		input := m.AddRequestsFromLedgerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddRequestsFromLedgerInput{p, p1}, "ExecutionBrokerI.AddRequestsFromLedger got unexpected parameters")
		}

		return
	}

	if m.AddRequestsFromLedgerFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddRequestsFromLedger. %v %v", p, p1)
		return
	}

	m.AddRequestsFromLedgerFunc(p, p1...)
}

//AddRequestsFromLedgerMinimockCounter returns a count of ExecutionBrokerIMock.AddRequestsFromLedgerFunc invocations
func (m *ExecutionBrokerIMock) AddRequestsFromLedgerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddRequestsFromLedgerCounter)
}

//AddRequestsFromLedgerMinimockPreCounter returns the value of ExecutionBrokerIMock.AddRequestsFromLedger invocations
func (m *ExecutionBrokerIMock) AddRequestsFromLedgerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddRequestsFromLedgerPreCounter)
}

//AddRequestsFromLedgerFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) AddRequestsFromLedgerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddRequestsFromLedgerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddRequestsFromLedgerCounter) == uint64(len(m.AddRequestsFromLedgerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddRequestsFromLedgerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddRequestsFromLedgerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddRequestsFromLedgerFunc != nil {
		return atomic.LoadUint64(&m.AddRequestsFromLedgerCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockAddRequestsFromPrevExecutor struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockAddRequestsFromPrevExecutorExpectation
	expectationSeries []*ExecutionBrokerIMockAddRequestsFromPrevExecutorExpectation
}

type ExecutionBrokerIMockAddRequestsFromPrevExecutorExpectation struct {
	input *ExecutionBrokerIMockAddRequestsFromPrevExecutorInput
}

type ExecutionBrokerIMockAddRequestsFromPrevExecutorInput struct {
	p  context.Context
	p1 []*Transcript
}

//Expect specifies that invocation of ExecutionBrokerI.AddRequestsFromPrevExecutor is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockAddRequestsFromPrevExecutor) Expect(p context.Context, p1 ...*Transcript) *mExecutionBrokerIMockAddRequestsFromPrevExecutor {
	m.mock.AddRequestsFromPrevExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddRequestsFromPrevExecutorExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockAddRequestsFromPrevExecutorInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.AddRequestsFromPrevExecutor
func (m *mExecutionBrokerIMockAddRequestsFromPrevExecutor) Return() *ExecutionBrokerIMock {
	m.mock.AddRequestsFromPrevExecutorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockAddRequestsFromPrevExecutorExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.AddRequestsFromPrevExecutor is expected once
func (m *mExecutionBrokerIMockAddRequestsFromPrevExecutor) ExpectOnce(p context.Context, p1 ...*Transcript) *ExecutionBrokerIMockAddRequestsFromPrevExecutorExpectation {
	m.mock.AddRequestsFromPrevExecutorFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockAddRequestsFromPrevExecutorExpectation{}
	expectation.input = &ExecutionBrokerIMockAddRequestsFromPrevExecutorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.AddRequestsFromPrevExecutor method
func (m *mExecutionBrokerIMockAddRequestsFromPrevExecutor) Set(f func(p context.Context, p1 ...*Transcript)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddRequestsFromPrevExecutorFunc = f
	return m.mock
}

//AddRequestsFromPrevExecutor implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) AddRequestsFromPrevExecutor(p context.Context, p1 ...*Transcript) {
	counter := atomic.AddUint64(&m.AddRequestsFromPrevExecutorPreCounter, 1)
	defer atomic.AddUint64(&m.AddRequestsFromPrevExecutorCounter, 1)

	if len(m.AddRequestsFromPrevExecutorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddRequestsFromPrevExecutorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddRequestsFromPrevExecutor. %v %v", p, p1)
			return
		}

		input := m.AddRequestsFromPrevExecutorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddRequestsFromPrevExecutorInput{p, p1}, "ExecutionBrokerI.AddRequestsFromPrevExecutor got unexpected parameters")

		return
	}

	if m.AddRequestsFromPrevExecutorMock.mainExpectation != nil {

		input := m.AddRequestsFromPrevExecutorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockAddRequestsFromPrevExecutorInput{p, p1}, "ExecutionBrokerI.AddRequestsFromPrevExecutor got unexpected parameters")
		}

		return
	}

	if m.AddRequestsFromPrevExecutorFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.AddRequestsFromPrevExecutor. %v %v", p, p1)
		return
	}

	m.AddRequestsFromPrevExecutorFunc(p, p1...)
}

//AddRequestsFromPrevExecutorMinimockCounter returns a count of ExecutionBrokerIMock.AddRequestsFromPrevExecutorFunc invocations
func (m *ExecutionBrokerIMock) AddRequestsFromPrevExecutorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddRequestsFromPrevExecutorCounter)
}

//AddRequestsFromPrevExecutorMinimockPreCounter returns the value of ExecutionBrokerIMock.AddRequestsFromPrevExecutor invocations
func (m *ExecutionBrokerIMock) AddRequestsFromPrevExecutorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddRequestsFromPrevExecutorPreCounter)
}

//AddRequestsFromPrevExecutorFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) AddRequestsFromPrevExecutorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddRequestsFromPrevExecutorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddRequestsFromPrevExecutorCounter) == uint64(len(m.AddRequestsFromPrevExecutorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddRequestsFromPrevExecutorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddRequestsFromPrevExecutorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddRequestsFromPrevExecutorFunc != nil {
		return atomic.LoadUint64(&m.AddRequestsFromPrevExecutorCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockCheckExecutionLoop struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockCheckExecutionLoopExpectation
	expectationSeries []*ExecutionBrokerIMockCheckExecutionLoopExpectation
}

type ExecutionBrokerIMockCheckExecutionLoopExpectation struct {
	input  *ExecutionBrokerIMockCheckExecutionLoopInput
	result *ExecutionBrokerIMockCheckExecutionLoopResult
}

type ExecutionBrokerIMockCheckExecutionLoopInput struct {
	p  context.Context
	p1 string
}

type ExecutionBrokerIMockCheckExecutionLoopResult struct {
	r bool
}

//Expect specifies that invocation of ExecutionBrokerI.CheckExecutionLoop is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockCheckExecutionLoop) Expect(p context.Context, p1 string) *mExecutionBrokerIMockCheckExecutionLoop {
	m.mock.CheckExecutionLoopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockCheckExecutionLoopExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockCheckExecutionLoopInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.CheckExecutionLoop
func (m *mExecutionBrokerIMockCheckExecutionLoop) Return(r bool) *ExecutionBrokerIMock {
	m.mock.CheckExecutionLoopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockCheckExecutionLoopExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerIMockCheckExecutionLoopResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.CheckExecutionLoop is expected once
func (m *mExecutionBrokerIMockCheckExecutionLoop) ExpectOnce(p context.Context, p1 string) *ExecutionBrokerIMockCheckExecutionLoopExpectation {
	m.mock.CheckExecutionLoopFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockCheckExecutionLoopExpectation{}
	expectation.input = &ExecutionBrokerIMockCheckExecutionLoopInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerIMockCheckExecutionLoopExpectation) Return(r bool) {
	e.result = &ExecutionBrokerIMockCheckExecutionLoopResult{r}
}

//Set uses given function f as a mock of ExecutionBrokerI.CheckExecutionLoop method
func (m *mExecutionBrokerIMockCheckExecutionLoop) Set(f func(p context.Context, p1 string) (r bool)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CheckExecutionLoopFunc = f
	return m.mock
}

//CheckExecutionLoop implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) CheckExecutionLoop(p context.Context, p1 string) (r bool) {
	counter := atomic.AddUint64(&m.CheckExecutionLoopPreCounter, 1)
	defer atomic.AddUint64(&m.CheckExecutionLoopCounter, 1)

	if len(m.CheckExecutionLoopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CheckExecutionLoopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.CheckExecutionLoop. %v %v", p, p1)
			return
		}

		input := m.CheckExecutionLoopMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockCheckExecutionLoopInput{p, p1}, "ExecutionBrokerI.CheckExecutionLoop got unexpected parameters")

		result := m.CheckExecutionLoopMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.CheckExecutionLoop")
			return
		}

		r = result.r

		return
	}

	if m.CheckExecutionLoopMock.mainExpectation != nil {

		input := m.CheckExecutionLoopMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockCheckExecutionLoopInput{p, p1}, "ExecutionBrokerI.CheckExecutionLoop got unexpected parameters")
		}

		result := m.CheckExecutionLoopMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.CheckExecutionLoop")
		}

		r = result.r

		return
	}

	if m.CheckExecutionLoopFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.CheckExecutionLoop. %v %v", p, p1)
		return
	}

	return m.CheckExecutionLoopFunc(p, p1)
}

//CheckExecutionLoopMinimockCounter returns a count of ExecutionBrokerIMock.CheckExecutionLoopFunc invocations
func (m *ExecutionBrokerIMock) CheckExecutionLoopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CheckExecutionLoopCounter)
}

//CheckExecutionLoopMinimockPreCounter returns the value of ExecutionBrokerIMock.CheckExecutionLoop invocations
func (m *ExecutionBrokerIMock) CheckExecutionLoopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CheckExecutionLoopPreCounter)
}

//CheckExecutionLoopFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) CheckExecutionLoopFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CheckExecutionLoopMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CheckExecutionLoopCounter) == uint64(len(m.CheckExecutionLoopMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CheckExecutionLoopMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CheckExecutionLoopCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CheckExecutionLoopFunc != nil {
		return atomic.LoadUint64(&m.CheckExecutionLoopCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockFetchMoreRequestsFromLedger struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockFetchMoreRequestsFromLedgerExpectation
	expectationSeries []*ExecutionBrokerIMockFetchMoreRequestsFromLedgerExpectation
}

type ExecutionBrokerIMockFetchMoreRequestsFromLedgerExpectation struct {
	input *ExecutionBrokerIMockFetchMoreRequestsFromLedgerInput
}

type ExecutionBrokerIMockFetchMoreRequestsFromLedgerInput struct {
	p context.Context
}

//Expect specifies that invocation of ExecutionBrokerI.FetchMoreRequestsFromLedger is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockFetchMoreRequestsFromLedger) Expect(p context.Context) *mExecutionBrokerIMockFetchMoreRequestsFromLedger {
	m.mock.FetchMoreRequestsFromLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockFetchMoreRequestsFromLedgerExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockFetchMoreRequestsFromLedgerInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.FetchMoreRequestsFromLedger
func (m *mExecutionBrokerIMockFetchMoreRequestsFromLedger) Return() *ExecutionBrokerIMock {
	m.mock.FetchMoreRequestsFromLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockFetchMoreRequestsFromLedgerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.FetchMoreRequestsFromLedger is expected once
func (m *mExecutionBrokerIMockFetchMoreRequestsFromLedger) ExpectOnce(p context.Context) *ExecutionBrokerIMockFetchMoreRequestsFromLedgerExpectation {
	m.mock.FetchMoreRequestsFromLedgerFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockFetchMoreRequestsFromLedgerExpectation{}
	expectation.input = &ExecutionBrokerIMockFetchMoreRequestsFromLedgerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.FetchMoreRequestsFromLedger method
func (m *mExecutionBrokerIMockFetchMoreRequestsFromLedger) Set(f func(p context.Context)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FetchMoreRequestsFromLedgerFunc = f
	return m.mock
}

//FetchMoreRequestsFromLedger implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) FetchMoreRequestsFromLedger(p context.Context) {
	counter := atomic.AddUint64(&m.FetchMoreRequestsFromLedgerPreCounter, 1)
	defer atomic.AddUint64(&m.FetchMoreRequestsFromLedgerCounter, 1)

	if len(m.FetchMoreRequestsFromLedgerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FetchMoreRequestsFromLedgerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger. %v", p)
			return
		}

		input := m.FetchMoreRequestsFromLedgerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockFetchMoreRequestsFromLedgerInput{p}, "ExecutionBrokerI.FetchMoreRequestsFromLedger got unexpected parameters")

		return
	}

	if m.FetchMoreRequestsFromLedgerMock.mainExpectation != nil {

		input := m.FetchMoreRequestsFromLedgerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockFetchMoreRequestsFromLedgerInput{p}, "ExecutionBrokerI.FetchMoreRequestsFromLedger got unexpected parameters")
		}

		return
	}

	if m.FetchMoreRequestsFromLedgerFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger. %v", p)
		return
	}

	m.FetchMoreRequestsFromLedgerFunc(p)
}

//FetchMoreRequestsFromLedgerMinimockCounter returns a count of ExecutionBrokerIMock.FetchMoreRequestsFromLedgerFunc invocations
func (m *ExecutionBrokerIMock) FetchMoreRequestsFromLedgerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FetchMoreRequestsFromLedgerCounter)
}

//FetchMoreRequestsFromLedgerMinimockPreCounter returns the value of ExecutionBrokerIMock.FetchMoreRequestsFromLedger invocations
func (m *ExecutionBrokerIMock) FetchMoreRequestsFromLedgerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FetchMoreRequestsFromLedgerPreCounter)
}

//FetchMoreRequestsFromLedgerFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) FetchMoreRequestsFromLedgerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FetchMoreRequestsFromLedgerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FetchMoreRequestsFromLedgerCounter) == uint64(len(m.FetchMoreRequestsFromLedgerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FetchMoreRequestsFromLedgerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FetchMoreRequestsFromLedgerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FetchMoreRequestsFromLedgerFunc != nil {
		return atomic.LoadUint64(&m.FetchMoreRequestsFromLedgerCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockGetActiveTranscript struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockGetActiveTranscriptExpectation
	expectationSeries []*ExecutionBrokerIMockGetActiveTranscriptExpectation
}

type ExecutionBrokerIMockGetActiveTranscriptExpectation struct {
	input  *ExecutionBrokerIMockGetActiveTranscriptInput
	result *ExecutionBrokerIMockGetActiveTranscriptResult
}

type ExecutionBrokerIMockGetActiveTranscriptInput struct {
	p insolar.Reference
}

type ExecutionBrokerIMockGetActiveTranscriptResult struct {
	r *Transcript
}

//Expect specifies that invocation of ExecutionBrokerI.GetActiveTranscript is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockGetActiveTranscript) Expect(p insolar.Reference) *mExecutionBrokerIMockGetActiveTranscript {
	m.mock.GetActiveTranscriptFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockGetActiveTranscriptExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockGetActiveTranscriptInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.GetActiveTranscript
func (m *mExecutionBrokerIMockGetActiveTranscript) Return(r *Transcript) *ExecutionBrokerIMock {
	m.mock.GetActiveTranscriptFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockGetActiveTranscriptExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerIMockGetActiveTranscriptResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.GetActiveTranscript is expected once
func (m *mExecutionBrokerIMockGetActiveTranscript) ExpectOnce(p insolar.Reference) *ExecutionBrokerIMockGetActiveTranscriptExpectation {
	m.mock.GetActiveTranscriptFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockGetActiveTranscriptExpectation{}
	expectation.input = &ExecutionBrokerIMockGetActiveTranscriptInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerIMockGetActiveTranscriptExpectation) Return(r *Transcript) {
	e.result = &ExecutionBrokerIMockGetActiveTranscriptResult{r}
}

//Set uses given function f as a mock of ExecutionBrokerI.GetActiveTranscript method
func (m *mExecutionBrokerIMockGetActiveTranscript) Set(f func(p insolar.Reference) (r *Transcript)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveTranscriptFunc = f
	return m.mock
}

//GetActiveTranscript implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) GetActiveTranscript(p insolar.Reference) (r *Transcript) {
	counter := atomic.AddUint64(&m.GetActiveTranscriptPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveTranscriptCounter, 1)

	if len(m.GetActiveTranscriptMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveTranscriptMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.GetActiveTranscript. %v", p)
			return
		}

		input := m.GetActiveTranscriptMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockGetActiveTranscriptInput{p}, "ExecutionBrokerI.GetActiveTranscript got unexpected parameters")

		result := m.GetActiveTranscriptMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.GetActiveTranscript")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveTranscriptMock.mainExpectation != nil {

		input := m.GetActiveTranscriptMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockGetActiveTranscriptInput{p}, "ExecutionBrokerI.GetActiveTranscript got unexpected parameters")
		}

		result := m.GetActiveTranscriptMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.GetActiveTranscript")
		}

		r = result.r

		return
	}

	if m.GetActiveTranscriptFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.GetActiveTranscript. %v", p)
		return
	}

	return m.GetActiveTranscriptFunc(p)
}

//GetActiveTranscriptMinimockCounter returns a count of ExecutionBrokerIMock.GetActiveTranscriptFunc invocations
func (m *ExecutionBrokerIMock) GetActiveTranscriptMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveTranscriptCounter)
}

//GetActiveTranscriptMinimockPreCounter returns the value of ExecutionBrokerIMock.GetActiveTranscript invocations
func (m *ExecutionBrokerIMock) GetActiveTranscriptMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveTranscriptPreCounter)
}

//GetActiveTranscriptFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) GetActiveTranscriptFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveTranscriptMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveTranscriptCounter) == uint64(len(m.GetActiveTranscriptMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveTranscriptMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveTranscriptCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveTranscriptFunc != nil {
		return atomic.LoadUint64(&m.GetActiveTranscriptCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockIsKnownRequest struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockIsKnownRequestExpectation
	expectationSeries []*ExecutionBrokerIMockIsKnownRequestExpectation
}

type ExecutionBrokerIMockIsKnownRequestExpectation struct {
	input  *ExecutionBrokerIMockIsKnownRequestInput
	result *ExecutionBrokerIMockIsKnownRequestResult
}

type ExecutionBrokerIMockIsKnownRequestInput struct {
	p  context.Context
	p1 insolar.Reference
}

type ExecutionBrokerIMockIsKnownRequestResult struct {
	r bool
}

//Expect specifies that invocation of ExecutionBrokerI.IsKnownRequest is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockIsKnownRequest) Expect(p context.Context, p1 insolar.Reference) *mExecutionBrokerIMockIsKnownRequest {
	m.mock.IsKnownRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockIsKnownRequestExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockIsKnownRequestInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.IsKnownRequest
func (m *mExecutionBrokerIMockIsKnownRequest) Return(r bool) *ExecutionBrokerIMock {
	m.mock.IsKnownRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockIsKnownRequestExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerIMockIsKnownRequestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.IsKnownRequest is expected once
func (m *mExecutionBrokerIMockIsKnownRequest) ExpectOnce(p context.Context, p1 insolar.Reference) *ExecutionBrokerIMockIsKnownRequestExpectation {
	m.mock.IsKnownRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockIsKnownRequestExpectation{}
	expectation.input = &ExecutionBrokerIMockIsKnownRequestInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerIMockIsKnownRequestExpectation) Return(r bool) {
	e.result = &ExecutionBrokerIMockIsKnownRequestResult{r}
}

//Set uses given function f as a mock of ExecutionBrokerI.IsKnownRequest method
func (m *mExecutionBrokerIMockIsKnownRequest) Set(f func(p context.Context, p1 insolar.Reference) (r bool)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsKnownRequestFunc = f
	return m.mock
}

//IsKnownRequest implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) IsKnownRequest(p context.Context, p1 insolar.Reference) (r bool) {
	counter := atomic.AddUint64(&m.IsKnownRequestPreCounter, 1)
	defer atomic.AddUint64(&m.IsKnownRequestCounter, 1)

	if len(m.IsKnownRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsKnownRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.IsKnownRequest. %v %v", p, p1)
			return
		}

		input := m.IsKnownRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockIsKnownRequestInput{p, p1}, "ExecutionBrokerI.IsKnownRequest got unexpected parameters")

		result := m.IsKnownRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.IsKnownRequest")
			return
		}

		r = result.r

		return
	}

	if m.IsKnownRequestMock.mainExpectation != nil {

		input := m.IsKnownRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockIsKnownRequestInput{p, p1}, "ExecutionBrokerI.IsKnownRequest got unexpected parameters")
		}

		result := m.IsKnownRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.IsKnownRequest")
		}

		r = result.r

		return
	}

	if m.IsKnownRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.IsKnownRequest. %v %v", p, p1)
		return
	}

	return m.IsKnownRequestFunc(p, p1)
}

//IsKnownRequestMinimockCounter returns a count of ExecutionBrokerIMock.IsKnownRequestFunc invocations
func (m *ExecutionBrokerIMock) IsKnownRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsKnownRequestCounter)
}

//IsKnownRequestMinimockPreCounter returns the value of ExecutionBrokerIMock.IsKnownRequest invocations
func (m *ExecutionBrokerIMock) IsKnownRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsKnownRequestPreCounter)
}

//IsKnownRequestFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) IsKnownRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsKnownRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsKnownRequestCounter) == uint64(len(m.IsKnownRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsKnownRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsKnownRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsKnownRequestFunc != nil {
		return atomic.LoadUint64(&m.IsKnownRequestCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockMoreRequestsOnLedger struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockMoreRequestsOnLedgerExpectation
	expectationSeries []*ExecutionBrokerIMockMoreRequestsOnLedgerExpectation
}

type ExecutionBrokerIMockMoreRequestsOnLedgerExpectation struct {
	input *ExecutionBrokerIMockMoreRequestsOnLedgerInput
}

type ExecutionBrokerIMockMoreRequestsOnLedgerInput struct {
	p context.Context
}

//Expect specifies that invocation of ExecutionBrokerI.MoreRequestsOnLedger is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockMoreRequestsOnLedger) Expect(p context.Context) *mExecutionBrokerIMockMoreRequestsOnLedger {
	m.mock.MoreRequestsOnLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockMoreRequestsOnLedgerExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockMoreRequestsOnLedgerInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.MoreRequestsOnLedger
func (m *mExecutionBrokerIMockMoreRequestsOnLedger) Return() *ExecutionBrokerIMock {
	m.mock.MoreRequestsOnLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockMoreRequestsOnLedgerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.MoreRequestsOnLedger is expected once
func (m *mExecutionBrokerIMockMoreRequestsOnLedger) ExpectOnce(p context.Context) *ExecutionBrokerIMockMoreRequestsOnLedgerExpectation {
	m.mock.MoreRequestsOnLedgerFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockMoreRequestsOnLedgerExpectation{}
	expectation.input = &ExecutionBrokerIMockMoreRequestsOnLedgerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.MoreRequestsOnLedger method
func (m *mExecutionBrokerIMockMoreRequestsOnLedger) Set(f func(p context.Context)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MoreRequestsOnLedgerFunc = f
	return m.mock
}

//MoreRequestsOnLedger implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) MoreRequestsOnLedger(p context.Context) {
	counter := atomic.AddUint64(&m.MoreRequestsOnLedgerPreCounter, 1)
	defer atomic.AddUint64(&m.MoreRequestsOnLedgerCounter, 1)

	if len(m.MoreRequestsOnLedgerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MoreRequestsOnLedgerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.MoreRequestsOnLedger. %v", p)
			return
		}

		input := m.MoreRequestsOnLedgerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockMoreRequestsOnLedgerInput{p}, "ExecutionBrokerI.MoreRequestsOnLedger got unexpected parameters")

		return
	}

	if m.MoreRequestsOnLedgerMock.mainExpectation != nil {

		input := m.MoreRequestsOnLedgerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockMoreRequestsOnLedgerInput{p}, "ExecutionBrokerI.MoreRequestsOnLedger got unexpected parameters")
		}

		return
	}

	if m.MoreRequestsOnLedgerFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.MoreRequestsOnLedger. %v", p)
		return
	}

	m.MoreRequestsOnLedgerFunc(p)
}

//MoreRequestsOnLedgerMinimockCounter returns a count of ExecutionBrokerIMock.MoreRequestsOnLedgerFunc invocations
func (m *ExecutionBrokerIMock) MoreRequestsOnLedgerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MoreRequestsOnLedgerCounter)
}

//MoreRequestsOnLedgerMinimockPreCounter returns the value of ExecutionBrokerIMock.MoreRequestsOnLedger invocations
func (m *ExecutionBrokerIMock) MoreRequestsOnLedgerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MoreRequestsOnLedgerPreCounter)
}

//MoreRequestsOnLedgerFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) MoreRequestsOnLedgerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MoreRequestsOnLedgerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MoreRequestsOnLedgerCounter) == uint64(len(m.MoreRequestsOnLedgerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MoreRequestsOnLedgerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MoreRequestsOnLedgerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MoreRequestsOnLedgerFunc != nil {
		return atomic.LoadUint64(&m.MoreRequestsOnLedgerCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockNoMoreRequestsOnLedger struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockNoMoreRequestsOnLedgerExpectation
	expectationSeries []*ExecutionBrokerIMockNoMoreRequestsOnLedgerExpectation
}

type ExecutionBrokerIMockNoMoreRequestsOnLedgerExpectation struct {
	input *ExecutionBrokerIMockNoMoreRequestsOnLedgerInput
}

type ExecutionBrokerIMockNoMoreRequestsOnLedgerInput struct {
	p context.Context
}

//Expect specifies that invocation of ExecutionBrokerI.NoMoreRequestsOnLedger is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockNoMoreRequestsOnLedger) Expect(p context.Context) *mExecutionBrokerIMockNoMoreRequestsOnLedger {
	m.mock.NoMoreRequestsOnLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockNoMoreRequestsOnLedgerExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockNoMoreRequestsOnLedgerInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.NoMoreRequestsOnLedger
func (m *mExecutionBrokerIMockNoMoreRequestsOnLedger) Return() *ExecutionBrokerIMock {
	m.mock.NoMoreRequestsOnLedgerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockNoMoreRequestsOnLedgerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.NoMoreRequestsOnLedger is expected once
func (m *mExecutionBrokerIMockNoMoreRequestsOnLedger) ExpectOnce(p context.Context) *ExecutionBrokerIMockNoMoreRequestsOnLedgerExpectation {
	m.mock.NoMoreRequestsOnLedgerFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockNoMoreRequestsOnLedgerExpectation{}
	expectation.input = &ExecutionBrokerIMockNoMoreRequestsOnLedgerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.NoMoreRequestsOnLedger method
func (m *mExecutionBrokerIMockNoMoreRequestsOnLedger) Set(f func(p context.Context)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NoMoreRequestsOnLedgerFunc = f
	return m.mock
}

//NoMoreRequestsOnLedger implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) NoMoreRequestsOnLedger(p context.Context) {
	counter := atomic.AddUint64(&m.NoMoreRequestsOnLedgerPreCounter, 1)
	defer atomic.AddUint64(&m.NoMoreRequestsOnLedgerCounter, 1)

	if len(m.NoMoreRequestsOnLedgerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NoMoreRequestsOnLedgerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.NoMoreRequestsOnLedger. %v", p)
			return
		}

		input := m.NoMoreRequestsOnLedgerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockNoMoreRequestsOnLedgerInput{p}, "ExecutionBrokerI.NoMoreRequestsOnLedger got unexpected parameters")

		return
	}

	if m.NoMoreRequestsOnLedgerMock.mainExpectation != nil {

		input := m.NoMoreRequestsOnLedgerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockNoMoreRequestsOnLedgerInput{p}, "ExecutionBrokerI.NoMoreRequestsOnLedger got unexpected parameters")
		}

		return
	}

	if m.NoMoreRequestsOnLedgerFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.NoMoreRequestsOnLedger. %v", p)
		return
	}

	m.NoMoreRequestsOnLedgerFunc(p)
}

//NoMoreRequestsOnLedgerMinimockCounter returns a count of ExecutionBrokerIMock.NoMoreRequestsOnLedgerFunc invocations
func (m *ExecutionBrokerIMock) NoMoreRequestsOnLedgerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NoMoreRequestsOnLedgerCounter)
}

//NoMoreRequestsOnLedgerMinimockPreCounter returns the value of ExecutionBrokerIMock.NoMoreRequestsOnLedger invocations
func (m *ExecutionBrokerIMock) NoMoreRequestsOnLedgerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NoMoreRequestsOnLedgerPreCounter)
}

//NoMoreRequestsOnLedgerFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) NoMoreRequestsOnLedgerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NoMoreRequestsOnLedgerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NoMoreRequestsOnLedgerCounter) == uint64(len(m.NoMoreRequestsOnLedgerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NoMoreRequestsOnLedgerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NoMoreRequestsOnLedgerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NoMoreRequestsOnLedgerFunc != nil {
		return atomic.LoadUint64(&m.NoMoreRequestsOnLedgerCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockOnPulse struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockOnPulseExpectation
	expectationSeries []*ExecutionBrokerIMockOnPulseExpectation
}

type ExecutionBrokerIMockOnPulseExpectation struct {
	input  *ExecutionBrokerIMockOnPulseInput
	result *ExecutionBrokerIMockOnPulseResult
}

type ExecutionBrokerIMockOnPulseInput struct {
	p  context.Context
	p1 bool
}

type ExecutionBrokerIMockOnPulseResult struct {
	r  bool
	r1 []insolar.Message
}

//Expect specifies that invocation of ExecutionBrokerI.OnPulse is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockOnPulse) Expect(p context.Context, p1 bool) *mExecutionBrokerIMockOnPulse {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockOnPulseExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockOnPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.OnPulse
func (m *mExecutionBrokerIMockOnPulse) Return(r bool, r1 []insolar.Message) *ExecutionBrokerIMock {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockOnPulseExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerIMockOnPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.OnPulse is expected once
func (m *mExecutionBrokerIMockOnPulse) ExpectOnce(p context.Context, p1 bool) *ExecutionBrokerIMockOnPulseExpectation {
	m.mock.OnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockOnPulseExpectation{}
	expectation.input = &ExecutionBrokerIMockOnPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerIMockOnPulseExpectation) Return(r bool, r1 []insolar.Message) {
	e.result = &ExecutionBrokerIMockOnPulseResult{r, r1}
}

//Set uses given function f as a mock of ExecutionBrokerI.OnPulse method
func (m *mExecutionBrokerIMockOnPulse) Set(f func(p context.Context, p1 bool) (r bool, r1 []insolar.Message)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) OnPulse(p context.Context, p1 bool) (r bool, r1 []insolar.Message) {
	counter := atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if len(m.OnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.OnPulse. %v %v", p, p1)
			return
		}

		input := m.OnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockOnPulseInput{p, p1}, "ExecutionBrokerI.OnPulse got unexpected parameters")

		result := m.OnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.OnPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.OnPulseMock.mainExpectation != nil {

		input := m.OnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockOnPulseInput{p, p1}, "ExecutionBrokerI.OnPulse got unexpected parameters")
		}

		result := m.OnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.OnPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.OnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.OnPulse. %v %v", p, p1)
		return
	}

	return m.OnPulseFunc(p, p1)
}

//OnPulseMinimockCounter returns a count of ExecutionBrokerIMock.OnPulseFunc invocations
func (m *ExecutionBrokerIMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of ExecutionBrokerIMock.OnPulse invocations
func (m *ExecutionBrokerIMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//OnPulseFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) OnPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OnPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OnPulseCounter) == uint64(len(m.OnPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OnPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OnPulseFunc != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockPendingState struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockPendingStateExpectation
	expectationSeries []*ExecutionBrokerIMockPendingStateExpectation
}

type ExecutionBrokerIMockPendingStateExpectation struct {
	result *ExecutionBrokerIMockPendingStateResult
}

type ExecutionBrokerIMockPendingStateResult struct {
	r insolar.PendingState
}

//Expect specifies that invocation of ExecutionBrokerI.PendingState is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockPendingState) Expect() *mExecutionBrokerIMockPendingState {
	m.mock.PendingStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPendingStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of ExecutionBrokerI.PendingState
func (m *mExecutionBrokerIMockPendingState) Return(r insolar.PendingState) *ExecutionBrokerIMock {
	m.mock.PendingStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPendingStateExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerIMockPendingStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.PendingState is expected once
func (m *mExecutionBrokerIMockPendingState) ExpectOnce() *ExecutionBrokerIMockPendingStateExpectation {
	m.mock.PendingStateFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockPendingStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerIMockPendingStateExpectation) Return(r insolar.PendingState) {
	e.result = &ExecutionBrokerIMockPendingStateResult{r}
}

//Set uses given function f as a mock of ExecutionBrokerI.PendingState method
func (m *mExecutionBrokerIMockPendingState) Set(f func() (r insolar.PendingState)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PendingStateFunc = f
	return m.mock
}

//PendingState implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) PendingState() (r insolar.PendingState) {
	counter := atomic.AddUint64(&m.PendingStatePreCounter, 1)
	defer atomic.AddUint64(&m.PendingStateCounter, 1)

	if len(m.PendingStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PendingStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PendingState.")
			return
		}

		result := m.PendingStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.PendingState")
			return
		}

		r = result.r

		return
	}

	if m.PendingStateMock.mainExpectation != nil {

		result := m.PendingStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.PendingState")
		}

		r = result.r

		return
	}

	if m.PendingStateFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PendingState.")
		return
	}

	return m.PendingStateFunc()
}

//PendingStateMinimockCounter returns a count of ExecutionBrokerIMock.PendingStateFunc invocations
func (m *ExecutionBrokerIMock) PendingStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PendingStateCounter)
}

//PendingStateMinimockPreCounter returns the value of ExecutionBrokerIMock.PendingState invocations
func (m *ExecutionBrokerIMock) PendingStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PendingStatePreCounter)
}

//PendingStateFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) PendingStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PendingStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PendingStateCounter) == uint64(len(m.PendingStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PendingStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PendingStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PendingStateFunc != nil {
		return atomic.LoadUint64(&m.PendingStateCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockPrevExecutorFinishedPending struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation
	expectationSeries []*ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation
}

type ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation struct {
	input  *ExecutionBrokerIMockPrevExecutorFinishedPendingInput
	result *ExecutionBrokerIMockPrevExecutorFinishedPendingResult
}

type ExecutionBrokerIMockPrevExecutorFinishedPendingInput struct {
	p context.Context
}

type ExecutionBrokerIMockPrevExecutorFinishedPendingResult struct {
	r error
}

//Expect specifies that invocation of ExecutionBrokerI.PrevExecutorFinishedPending is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockPrevExecutorFinishedPending) Expect(p context.Context) *mExecutionBrokerIMockPrevExecutorFinishedPending {
	m.mock.PrevExecutorFinishedPendingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockPrevExecutorFinishedPendingInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.PrevExecutorFinishedPending
func (m *mExecutionBrokerIMockPrevExecutorFinishedPending) Return(r error) *ExecutionBrokerIMock {
	m.mock.PrevExecutorFinishedPendingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation{}
	}
	m.mainExpectation.result = &ExecutionBrokerIMockPrevExecutorFinishedPendingResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.PrevExecutorFinishedPending is expected once
func (m *mExecutionBrokerIMockPrevExecutorFinishedPending) ExpectOnce(p context.Context) *ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation {
	m.mock.PrevExecutorFinishedPendingFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation{}
	expectation.input = &ExecutionBrokerIMockPrevExecutorFinishedPendingInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionBrokerIMockPrevExecutorFinishedPendingExpectation) Return(r error) {
	e.result = &ExecutionBrokerIMockPrevExecutorFinishedPendingResult{r}
}

//Set uses given function f as a mock of ExecutionBrokerI.PrevExecutorFinishedPending method
func (m *mExecutionBrokerIMockPrevExecutorFinishedPending) Set(f func(p context.Context) (r error)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PrevExecutorFinishedPendingFunc = f
	return m.mock
}

//PrevExecutorFinishedPending implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) PrevExecutorFinishedPending(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.PrevExecutorFinishedPendingPreCounter, 1)
	defer atomic.AddUint64(&m.PrevExecutorFinishedPendingCounter, 1)

	if len(m.PrevExecutorFinishedPendingMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PrevExecutorFinishedPendingMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PrevExecutorFinishedPending. %v", p)
			return
		}

		input := m.PrevExecutorFinishedPendingMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrevExecutorFinishedPendingInput{p}, "ExecutionBrokerI.PrevExecutorFinishedPending got unexpected parameters")

		result := m.PrevExecutorFinishedPendingMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.PrevExecutorFinishedPending")
			return
		}

		r = result.r

		return
	}

	if m.PrevExecutorFinishedPendingMock.mainExpectation != nil {

		input := m.PrevExecutorFinishedPendingMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrevExecutorFinishedPendingInput{p}, "ExecutionBrokerI.PrevExecutorFinishedPending got unexpected parameters")
		}

		result := m.PrevExecutorFinishedPendingMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionBrokerIMock.PrevExecutorFinishedPending")
		}

		r = result.r

		return
	}

	if m.PrevExecutorFinishedPendingFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PrevExecutorFinishedPending. %v", p)
		return
	}

	return m.PrevExecutorFinishedPendingFunc(p)
}

//PrevExecutorFinishedPendingMinimockCounter returns a count of ExecutionBrokerIMock.PrevExecutorFinishedPendingFunc invocations
func (m *ExecutionBrokerIMock) PrevExecutorFinishedPendingMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PrevExecutorFinishedPendingCounter)
}

//PrevExecutorFinishedPendingMinimockPreCounter returns the value of ExecutionBrokerIMock.PrevExecutorFinishedPending invocations
func (m *ExecutionBrokerIMock) PrevExecutorFinishedPendingMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PrevExecutorFinishedPendingPreCounter)
}

//PrevExecutorFinishedPendingFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) PrevExecutorFinishedPendingFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PrevExecutorFinishedPendingMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PrevExecutorFinishedPendingCounter) == uint64(len(m.PrevExecutorFinishedPendingMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PrevExecutorFinishedPendingMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PrevExecutorFinishedPendingCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PrevExecutorFinishedPendingFunc != nil {
		return atomic.LoadUint64(&m.PrevExecutorFinishedPendingCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockPrevExecutorPendingResult struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockPrevExecutorPendingResultExpectation
	expectationSeries []*ExecutionBrokerIMockPrevExecutorPendingResultExpectation
}

type ExecutionBrokerIMockPrevExecutorPendingResultExpectation struct {
	input *ExecutionBrokerIMockPrevExecutorPendingResultInput
}

type ExecutionBrokerIMockPrevExecutorPendingResultInput struct {
	p  context.Context
	p1 insolar.PendingState
}

//Expect specifies that invocation of ExecutionBrokerI.PrevExecutorPendingResult is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockPrevExecutorPendingResult) Expect(p context.Context, p1 insolar.PendingState) *mExecutionBrokerIMockPrevExecutorPendingResult {
	m.mock.PrevExecutorPendingResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrevExecutorPendingResultExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockPrevExecutorPendingResultInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.PrevExecutorPendingResult
func (m *mExecutionBrokerIMockPrevExecutorPendingResult) Return() *ExecutionBrokerIMock {
	m.mock.PrevExecutorPendingResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrevExecutorPendingResultExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.PrevExecutorPendingResult is expected once
func (m *mExecutionBrokerIMockPrevExecutorPendingResult) ExpectOnce(p context.Context, p1 insolar.PendingState) *ExecutionBrokerIMockPrevExecutorPendingResultExpectation {
	m.mock.PrevExecutorPendingResultFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockPrevExecutorPendingResultExpectation{}
	expectation.input = &ExecutionBrokerIMockPrevExecutorPendingResultInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.PrevExecutorPendingResult method
func (m *mExecutionBrokerIMockPrevExecutorPendingResult) Set(f func(p context.Context, p1 insolar.PendingState)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PrevExecutorPendingResultFunc = f
	return m.mock
}

//PrevExecutorPendingResult implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) PrevExecutorPendingResult(p context.Context, p1 insolar.PendingState) {
	counter := atomic.AddUint64(&m.PrevExecutorPendingResultPreCounter, 1)
	defer atomic.AddUint64(&m.PrevExecutorPendingResultCounter, 1)

	if len(m.PrevExecutorPendingResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PrevExecutorPendingResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PrevExecutorPendingResult. %v %v", p, p1)
			return
		}

		input := m.PrevExecutorPendingResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrevExecutorPendingResultInput{p, p1}, "ExecutionBrokerI.PrevExecutorPendingResult got unexpected parameters")

		return
	}

	if m.PrevExecutorPendingResultMock.mainExpectation != nil {

		input := m.PrevExecutorPendingResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrevExecutorPendingResultInput{p, p1}, "ExecutionBrokerI.PrevExecutorPendingResult got unexpected parameters")
		}

		return
	}

	if m.PrevExecutorPendingResultFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PrevExecutorPendingResult. %v %v", p, p1)
		return
	}

	m.PrevExecutorPendingResultFunc(p, p1)
}

//PrevExecutorPendingResultMinimockCounter returns a count of ExecutionBrokerIMock.PrevExecutorPendingResultFunc invocations
func (m *ExecutionBrokerIMock) PrevExecutorPendingResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PrevExecutorPendingResultCounter)
}

//PrevExecutorPendingResultMinimockPreCounter returns the value of ExecutionBrokerIMock.PrevExecutorPendingResult invocations
func (m *ExecutionBrokerIMock) PrevExecutorPendingResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PrevExecutorPendingResultPreCounter)
}

//PrevExecutorPendingResultFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) PrevExecutorPendingResultFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PrevExecutorPendingResultMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PrevExecutorPendingResultCounter) == uint64(len(m.PrevExecutorPendingResultMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PrevExecutorPendingResultMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PrevExecutorPendingResultCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PrevExecutorPendingResultFunc != nil {
		return atomic.LoadUint64(&m.PrevExecutorPendingResultCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockPrevExecutorStillExecuting struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockPrevExecutorStillExecutingExpectation
	expectationSeries []*ExecutionBrokerIMockPrevExecutorStillExecutingExpectation
}

type ExecutionBrokerIMockPrevExecutorStillExecutingExpectation struct {
	input *ExecutionBrokerIMockPrevExecutorStillExecutingInput
}

type ExecutionBrokerIMockPrevExecutorStillExecutingInput struct {
	p context.Context
}

//Expect specifies that invocation of ExecutionBrokerI.PrevExecutorStillExecuting is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockPrevExecutorStillExecuting) Expect(p context.Context) *mExecutionBrokerIMockPrevExecutorStillExecuting {
	m.mock.PrevExecutorStillExecutingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrevExecutorStillExecutingExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockPrevExecutorStillExecutingInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.PrevExecutorStillExecuting
func (m *mExecutionBrokerIMockPrevExecutorStillExecuting) Return() *ExecutionBrokerIMock {
	m.mock.PrevExecutorStillExecutingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrevExecutorStillExecutingExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.PrevExecutorStillExecuting is expected once
func (m *mExecutionBrokerIMockPrevExecutorStillExecuting) ExpectOnce(p context.Context) *ExecutionBrokerIMockPrevExecutorStillExecutingExpectation {
	m.mock.PrevExecutorStillExecutingFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockPrevExecutorStillExecutingExpectation{}
	expectation.input = &ExecutionBrokerIMockPrevExecutorStillExecutingInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.PrevExecutorStillExecuting method
func (m *mExecutionBrokerIMockPrevExecutorStillExecuting) Set(f func(p context.Context)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PrevExecutorStillExecutingFunc = f
	return m.mock
}

//PrevExecutorStillExecuting implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) PrevExecutorStillExecuting(p context.Context) {
	counter := atomic.AddUint64(&m.PrevExecutorStillExecutingPreCounter, 1)
	defer atomic.AddUint64(&m.PrevExecutorStillExecutingCounter, 1)

	if len(m.PrevExecutorStillExecutingMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PrevExecutorStillExecutingMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PrevExecutorStillExecuting. %v", p)
			return
		}

		input := m.PrevExecutorStillExecutingMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrevExecutorStillExecutingInput{p}, "ExecutionBrokerI.PrevExecutorStillExecuting got unexpected parameters")

		return
	}

	if m.PrevExecutorStillExecutingMock.mainExpectation != nil {

		input := m.PrevExecutorStillExecutingMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrevExecutorStillExecutingInput{p}, "ExecutionBrokerI.PrevExecutorStillExecuting got unexpected parameters")
		}

		return
	}

	if m.PrevExecutorStillExecutingFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.PrevExecutorStillExecuting. %v", p)
		return
	}

	m.PrevExecutorStillExecutingFunc(p)
}

//PrevExecutorStillExecutingMinimockCounter returns a count of ExecutionBrokerIMock.PrevExecutorStillExecutingFunc invocations
func (m *ExecutionBrokerIMock) PrevExecutorStillExecutingMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PrevExecutorStillExecutingCounter)
}

//PrevExecutorStillExecutingMinimockPreCounter returns the value of ExecutionBrokerIMock.PrevExecutorStillExecuting invocations
func (m *ExecutionBrokerIMock) PrevExecutorStillExecutingMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PrevExecutorStillExecutingPreCounter)
}

//PrevExecutorStillExecutingFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) PrevExecutorStillExecutingFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PrevExecutorStillExecutingMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PrevExecutorStillExecutingCounter) == uint64(len(m.PrevExecutorStillExecutingMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PrevExecutorStillExecutingMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PrevExecutorStillExecutingCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PrevExecutorStillExecutingFunc != nil {
		return atomic.LoadUint64(&m.PrevExecutorStillExecutingCounter) > 0
	}

	return true
}

type mExecutionBrokerIMockSetNotPending struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockSetNotPendingExpectation
	expectationSeries []*ExecutionBrokerIMockSetNotPendingExpectation
}

type ExecutionBrokerIMockSetNotPendingExpectation struct {
	input *ExecutionBrokerIMockSetNotPendingInput
}

type ExecutionBrokerIMockSetNotPendingInput struct {
	p context.Context
}

//Expect specifies that invocation of ExecutionBrokerI.SetNotPending is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockSetNotPending) Expect(p context.Context) *mExecutionBrokerIMockSetNotPending {
	m.mock.SetNotPendingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockSetNotPendingExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockSetNotPendingInput{p}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.SetNotPending
func (m *mExecutionBrokerIMockSetNotPending) Return() *ExecutionBrokerIMock {
	m.mock.SetNotPendingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockSetNotPendingExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.SetNotPending is expected once
func (m *mExecutionBrokerIMockSetNotPending) ExpectOnce(p context.Context) *ExecutionBrokerIMockSetNotPendingExpectation {
	m.mock.SetNotPendingFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockSetNotPendingExpectation{}
	expectation.input = &ExecutionBrokerIMockSetNotPendingInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.SetNotPending method
func (m *mExecutionBrokerIMockSetNotPending) Set(f func(p context.Context)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetNotPendingFunc = f
	return m.mock
}

//SetNotPending implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) SetNotPending(p context.Context) {
	counter := atomic.AddUint64(&m.SetNotPendingPreCounter, 1)
	defer atomic.AddUint64(&m.SetNotPendingCounter, 1)

	if len(m.SetNotPendingMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetNotPendingMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.SetNotPending. %v", p)
			return
		}

		input := m.SetNotPendingMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockSetNotPendingInput{p}, "ExecutionBrokerI.SetNotPending got unexpected parameters")

		return
	}

	if m.SetNotPendingMock.mainExpectation != nil {

		input := m.SetNotPendingMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockSetNotPendingInput{p}, "ExecutionBrokerI.SetNotPending got unexpected parameters")
		}

		return
	}

	if m.SetNotPendingFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.SetNotPending. %v", p)
		return
	}

	m.SetNotPendingFunc(p)
}

//SetNotPendingMinimockCounter returns a count of ExecutionBrokerIMock.SetNotPendingFunc invocations
func (m *ExecutionBrokerIMock) SetNotPendingMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetNotPendingCounter)
}

//SetNotPendingMinimockPreCounter returns the value of ExecutionBrokerIMock.SetNotPending invocations
func (m *ExecutionBrokerIMock) SetNotPendingMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetNotPendingPreCounter)
}

//SetNotPendingFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) SetNotPendingFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetNotPendingMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetNotPendingCounter) == uint64(len(m.SetNotPendingMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetNotPendingMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetNotPendingCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetNotPendingFunc != nil {
		return atomic.LoadUint64(&m.SetNotPendingCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExecutionBrokerIMock) ValidateCallCounters() {

	if !m.AbandonedRequestsOnLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AbandonedRequestsOnLedger")
	}

	if !m.AddAdditionalRequestFromPrevExecutorFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddAdditionalRequestFromPrevExecutor")
	}

	if !m.AddFreshRequestFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddFreshRequest")
	}

	if !m.AddRequestsFromLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddRequestsFromLedger")
	}

	if !m.AddRequestsFromPrevExecutorFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddRequestsFromPrevExecutor")
	}

	if !m.CheckExecutionLoopFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.CheckExecutionLoop")
	}

	if !m.FetchMoreRequestsFromLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger")
	}

	if !m.GetActiveTranscriptFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.GetActiveTranscript")
	}

	if !m.IsKnownRequestFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.IsKnownRequest")
	}

	if !m.MoreRequestsOnLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.MoreRequestsOnLedger")
	}

	if !m.NoMoreRequestsOnLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.NoMoreRequestsOnLedger")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.OnPulse")
	}

	if !m.PendingStateFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PendingState")
	}

	if !m.PrevExecutorFinishedPendingFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PrevExecutorFinishedPending")
	}

	if !m.PrevExecutorPendingResultFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PrevExecutorPendingResult")
	}

	if !m.PrevExecutorStillExecutingFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PrevExecutorStillExecuting")
	}

	if !m.SetNotPendingFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.SetNotPending")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExecutionBrokerIMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ExecutionBrokerIMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ExecutionBrokerIMock) MinimockFinish() {

	if !m.AbandonedRequestsOnLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AbandonedRequestsOnLedger")
	}

	if !m.AddAdditionalRequestFromPrevExecutorFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddAdditionalRequestFromPrevExecutor")
	}

	if !m.AddFreshRequestFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddFreshRequest")
	}

	if !m.AddRequestsFromLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddRequestsFromLedger")
	}

	if !m.AddRequestsFromPrevExecutorFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.AddRequestsFromPrevExecutor")
	}

	if !m.CheckExecutionLoopFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.CheckExecutionLoop")
	}

	if !m.FetchMoreRequestsFromLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger")
	}

	if !m.GetActiveTranscriptFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.GetActiveTranscript")
	}

	if !m.IsKnownRequestFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.IsKnownRequest")
	}

	if !m.MoreRequestsOnLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.MoreRequestsOnLedger")
	}

	if !m.NoMoreRequestsOnLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.NoMoreRequestsOnLedger")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.OnPulse")
	}

	if !m.PendingStateFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PendingState")
	}

	if !m.PrevExecutorFinishedPendingFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PrevExecutorFinishedPending")
	}

	if !m.PrevExecutorPendingResultFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PrevExecutorPendingResult")
	}

	if !m.PrevExecutorStillExecutingFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.PrevExecutorStillExecuting")
	}

	if !m.SetNotPendingFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.SetNotPending")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ExecutionBrokerIMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ExecutionBrokerIMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AbandonedRequestsOnLedgerFinished()
		ok = ok && m.AddAdditionalRequestFromPrevExecutorFinished()
		ok = ok && m.AddFreshRequestFinished()
		ok = ok && m.AddRequestsFromLedgerFinished()
		ok = ok && m.AddRequestsFromPrevExecutorFinished()
		ok = ok && m.CheckExecutionLoopFinished()
		ok = ok && m.FetchMoreRequestsFromLedgerFinished()
		ok = ok && m.GetActiveTranscriptFinished()
		ok = ok && m.IsKnownRequestFinished()
		ok = ok && m.MoreRequestsOnLedgerFinished()
		ok = ok && m.NoMoreRequestsOnLedgerFinished()
		ok = ok && m.OnPulseFinished()
		ok = ok && m.PendingStateFinished()
		ok = ok && m.PrevExecutorFinishedPendingFinished()
		ok = ok && m.PrevExecutorPendingResultFinished()
		ok = ok && m.PrevExecutorStillExecutingFinished()
		ok = ok && m.SetNotPendingFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AbandonedRequestsOnLedgerFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.AbandonedRequestsOnLedger")
			}

			if !m.AddAdditionalRequestFromPrevExecutorFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.AddAdditionalRequestFromPrevExecutor")
			}

			if !m.AddFreshRequestFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.AddFreshRequest")
			}

			if !m.AddRequestsFromLedgerFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.AddRequestsFromLedger")
			}

			if !m.AddRequestsFromPrevExecutorFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.AddRequestsFromPrevExecutor")
			}

			if !m.CheckExecutionLoopFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.CheckExecutionLoop")
			}

			if !m.FetchMoreRequestsFromLedgerFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger")
			}

			if !m.GetActiveTranscriptFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.GetActiveTranscript")
			}

			if !m.IsKnownRequestFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.IsKnownRequest")
			}

			if !m.MoreRequestsOnLedgerFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.MoreRequestsOnLedger")
			}

			if !m.NoMoreRequestsOnLedgerFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.NoMoreRequestsOnLedger")
			}

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.OnPulse")
			}

			if !m.PendingStateFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.PendingState")
			}

			if !m.PrevExecutorFinishedPendingFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.PrevExecutorFinishedPending")
			}

			if !m.PrevExecutorPendingResultFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.PrevExecutorPendingResult")
			}

			if !m.PrevExecutorStillExecutingFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.PrevExecutorStillExecuting")
			}

			if !m.SetNotPendingFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.SetNotPending")
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
func (m *ExecutionBrokerIMock) AllMocksCalled() bool {

	if !m.AbandonedRequestsOnLedgerFinished() {
		return false
	}

	if !m.AddAdditionalRequestFromPrevExecutorFinished() {
		return false
	}

	if !m.AddFreshRequestFinished() {
		return false
	}

	if !m.AddRequestsFromLedgerFinished() {
		return false
	}

	if !m.AddRequestsFromPrevExecutorFinished() {
		return false
	}

	if !m.CheckExecutionLoopFinished() {
		return false
	}

	if !m.FetchMoreRequestsFromLedgerFinished() {
		return false
	}

	if !m.GetActiveTranscriptFinished() {
		return false
	}

	if !m.IsKnownRequestFinished() {
		return false
	}

	if !m.MoreRequestsOnLedgerFinished() {
		return false
	}

	if !m.NoMoreRequestsOnLedgerFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	if !m.PendingStateFinished() {
		return false
	}

	if !m.PrevExecutorFinishedPendingFinished() {
		return false
	}

	if !m.PrevExecutorPendingResultFinished() {
		return false
	}

	if !m.PrevExecutorStillExecutingFinished() {
		return false
	}

	if !m.SetNotPendingFinished() {
		return false
	}

	return true
}

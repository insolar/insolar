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

	FetchMoreRequestsFromLedgerFunc       func(p context.Context)
	FetchMoreRequestsFromLedgerCounter    uint64
	FetchMoreRequestsFromLedgerPreCounter uint64
	FetchMoreRequestsFromLedgerMock       mExecutionBrokerIMockFetchMoreRequestsFromLedger

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

	PrependFunc       func(p context.Context, p1 bool, p2 ...*Transcript)
	PrependCounter    uint64
	PrependPreCounter uint64
	PrependMock       mExecutionBrokerIMockPrepend
}

//NewExecutionBrokerIMock returns a mock for github.com/insolar/insolar/logicrunner.ExecutionBrokerI
func NewExecutionBrokerIMock(t minimock.Tester) *ExecutionBrokerIMock {
	m := &ExecutionBrokerIMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FetchMoreRequestsFromLedgerMock = mExecutionBrokerIMockFetchMoreRequestsFromLedger{mock: m}
	m.IsKnownRequestMock = mExecutionBrokerIMockIsKnownRequest{mock: m}
	m.MoreRequestsOnLedgerMock = mExecutionBrokerIMockMoreRequestsOnLedger{mock: m}
	m.NoMoreRequestsOnLedgerMock = mExecutionBrokerIMockNoMoreRequestsOnLedger{mock: m}
	m.PrependMock = mExecutionBrokerIMockPrepend{mock: m}

	return m
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

type mExecutionBrokerIMockPrepend struct {
	mock              *ExecutionBrokerIMock
	mainExpectation   *ExecutionBrokerIMockPrependExpectation
	expectationSeries []*ExecutionBrokerIMockPrependExpectation
}

type ExecutionBrokerIMockPrependExpectation struct {
	input *ExecutionBrokerIMockPrependInput
}

type ExecutionBrokerIMockPrependInput struct {
	p  context.Context
	p1 bool
	p2 []*Transcript
}

//Expect specifies that invocation of ExecutionBrokerI.Prepend is expected from 1 to Infinity times
func (m *mExecutionBrokerIMockPrepend) Expect(p context.Context, p1 bool, p2 ...*Transcript) *mExecutionBrokerIMockPrepend {
	m.mock.PrependFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrependExpectation{}
	}
	m.mainExpectation.input = &ExecutionBrokerIMockPrependInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ExecutionBrokerI.Prepend
func (m *mExecutionBrokerIMockPrepend) Return() *ExecutionBrokerIMock {
	m.mock.PrependFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionBrokerIMockPrependExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionBrokerI.Prepend is expected once
func (m *mExecutionBrokerIMockPrepend) ExpectOnce(p context.Context, p1 bool, p2 ...*Transcript) *ExecutionBrokerIMockPrependExpectation {
	m.mock.PrependFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionBrokerIMockPrependExpectation{}
	expectation.input = &ExecutionBrokerIMockPrependInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionBrokerI.Prepend method
func (m *mExecutionBrokerIMockPrepend) Set(f func(p context.Context, p1 bool, p2 ...*Transcript)) *ExecutionBrokerIMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PrependFunc = f
	return m.mock
}

//Prepend implements github.com/insolar/insolar/logicrunner.ExecutionBrokerI interface
func (m *ExecutionBrokerIMock) Prepend(p context.Context, p1 bool, p2 ...*Transcript) {
	counter := atomic.AddUint64(&m.PrependPreCounter, 1)
	defer atomic.AddUint64(&m.PrependCounter, 1)

	if len(m.PrependMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PrependMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.Prepend. %v %v %v", p, p1, p2)
			return
		}

		input := m.PrependMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrependInput{p, p1, p2}, "ExecutionBrokerI.Prepend got unexpected parameters")

		return
	}

	if m.PrependMock.mainExpectation != nil {

		input := m.PrependMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionBrokerIMockPrependInput{p, p1, p2}, "ExecutionBrokerI.Prepend got unexpected parameters")
		}

		return
	}

	if m.PrependFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionBrokerIMock.Prepend. %v %v %v", p, p1, p2)
		return
	}

	m.PrependFunc(p, p1, p2...)
}

//PrependMinimockCounter returns a count of ExecutionBrokerIMock.PrependFunc invocations
func (m *ExecutionBrokerIMock) PrependMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PrependCounter)
}

//PrependMinimockPreCounter returns the value of ExecutionBrokerIMock.Prepend invocations
func (m *ExecutionBrokerIMock) PrependMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PrependPreCounter)
}

//PrependFinished returns true if mock invocations count is ok
func (m *ExecutionBrokerIMock) PrependFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PrependMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PrependCounter) == uint64(len(m.PrependMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PrependMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PrependCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PrependFunc != nil {
		return atomic.LoadUint64(&m.PrependCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExecutionBrokerIMock) ValidateCallCounters() {

	if !m.FetchMoreRequestsFromLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger")
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

	if !m.PrependFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.Prepend")
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

	if !m.FetchMoreRequestsFromLedgerFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger")
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

	if !m.PrependFinished() {
		m.t.Fatal("Expected call to ExecutionBrokerIMock.Prepend")
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
		ok = ok && m.FetchMoreRequestsFromLedgerFinished()
		ok = ok && m.IsKnownRequestFinished()
		ok = ok && m.MoreRequestsOnLedgerFinished()
		ok = ok && m.NoMoreRequestsOnLedgerFinished()
		ok = ok && m.PrependFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FetchMoreRequestsFromLedgerFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.FetchMoreRequestsFromLedger")
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

			if !m.PrependFinished() {
				m.t.Error("Expected call to ExecutionBrokerIMock.Prepend")
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

	if !m.FetchMoreRequestsFromLedgerFinished() {
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

	if !m.PrependFinished() {
		return false
	}

	return true
}

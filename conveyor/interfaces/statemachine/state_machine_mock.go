package statemachine

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StateMachine" can be found in github.com/insolar/insolar/conveyor/interfaces/statemachine
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//StateMachineMock implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine
type StateMachineMock struct {
	t minimock.Tester

	GetMigrationHandlerFunc       func(p StateID) (r MigrationHandler)
	GetMigrationHandlerCounter    uint64
	GetMigrationHandlerPreCounter uint64
	GetMigrationHandlerMock       mStateMachineMockGetMigrationHandler

	GetNestedHandlerFunc       func(p StateID) (r NestedHandler)
	GetNestedHandlerCounter    uint64
	GetNestedHandlerPreCounter uint64
	GetNestedHandlerMock       mStateMachineMockGetNestedHandler

	GetResponseErrorHandlerFunc       func(p StateID) (r ResponseErrorHandler)
	GetResponseErrorHandlerCounter    uint64
	GetResponseErrorHandlerPreCounter uint64
	GetResponseErrorHandlerMock       mStateMachineMockGetResponseErrorHandler

	GetResponseHandlerFunc       func(p StateID) (r AdapterResponseHandler)
	GetResponseHandlerCounter    uint64
	GetResponseHandlerPreCounter uint64
	GetResponseHandlerMock       mStateMachineMockGetResponseHandler

	GetTransitionErrorHandlerFunc       func(p StateID) (r TransitionErrorHandler)
	GetTransitionErrorHandlerCounter    uint64
	GetTransitionErrorHandlerPreCounter uint64
	GetTransitionErrorHandlerMock       mStateMachineMockGetTransitionErrorHandler

	GetTransitionHandlerFunc       func(p StateID) (r TransitHandler)
	GetTransitionHandlerCounter    uint64
	GetTransitionHandlerPreCounter uint64
	GetTransitionHandlerMock       mStateMachineMockGetTransitionHandler

	GetTypeIDFunc       func() (r ID)
	GetTypeIDCounter    uint64
	GetTypeIDPreCounter uint64
	GetTypeIDMock       mStateMachineMockGetTypeID
}

//NewStateMachineMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine
func NewStateMachineMock(t minimock.Tester) *StateMachineMock {
	m := &StateMachineMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetMigrationHandlerMock = mStateMachineMockGetMigrationHandler{mock: m}
	m.GetNestedHandlerMock = mStateMachineMockGetNestedHandler{mock: m}
	m.GetResponseErrorHandlerMock = mStateMachineMockGetResponseErrorHandler{mock: m}
	m.GetResponseHandlerMock = mStateMachineMockGetResponseHandler{mock: m}
	m.GetTransitionErrorHandlerMock = mStateMachineMockGetTransitionErrorHandler{mock: m}
	m.GetTransitionHandlerMock = mStateMachineMockGetTransitionHandler{mock: m}
	m.GetTypeIDMock = mStateMachineMockGetTypeID{mock: m}

	return m
}

type mStateMachineMockGetMigrationHandler struct {
	mock              *StateMachineMock
	mainExpectation   *StateMachineMockGetMigrationHandlerExpectation
	expectationSeries []*StateMachineMockGetMigrationHandlerExpectation
}

type StateMachineMockGetMigrationHandlerExpectation struct {
	input  *StateMachineMockGetMigrationHandlerInput
	result *StateMachineMockGetMigrationHandlerResult
}

type StateMachineMockGetMigrationHandlerInput struct {
	p StateID
}

type StateMachineMockGetMigrationHandlerResult struct {
	r MigrationHandler
}

//Expect specifies that invocation of StateMachine.GetMigrationHandler is expected from 1 to Infinity times
func (m *mStateMachineMockGetMigrationHandler) Expect(p StateID) *mStateMachineMockGetMigrationHandler {
	m.mock.GetMigrationHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetMigrationHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineMockGetMigrationHandlerInput{p}
	return m
}

//Return specifies results of invocation of StateMachine.GetMigrationHandler
func (m *mStateMachineMockGetMigrationHandler) Return(r MigrationHandler) *StateMachineMock {
	m.mock.GetMigrationHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetMigrationHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineMockGetMigrationHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachine.GetMigrationHandler is expected once
func (m *mStateMachineMockGetMigrationHandler) ExpectOnce(p StateID) *StateMachineMockGetMigrationHandlerExpectation {
	m.mock.GetMigrationHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineMockGetMigrationHandlerExpectation{}
	expectation.input = &StateMachineMockGetMigrationHandlerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineMockGetMigrationHandlerExpectation) Return(r MigrationHandler) {
	e.result = &StateMachineMockGetMigrationHandlerResult{r}
}

//Set uses given function f as a mock of StateMachine.GetMigrationHandler method
func (m *mStateMachineMockGetMigrationHandler) Set(f func(p StateID) (r MigrationHandler)) *StateMachineMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMigrationHandlerFunc = f
	return m.mock
}

//GetMigrationHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine interface
func (m *StateMachineMock) GetMigrationHandler(p StateID) (r MigrationHandler) {
	counter := atomic.AddUint64(&m.GetMigrationHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetMigrationHandlerCounter, 1)

	if len(m.GetMigrationHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMigrationHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineMock.GetMigrationHandler. %v", p)
			return
		}

		input := m.GetMigrationHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineMockGetMigrationHandlerInput{p}, "StateMachine.GetMigrationHandler got unexpected parameters")

		result := m.GetMigrationHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetMigrationHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetMigrationHandlerMock.mainExpectation != nil {

		input := m.GetMigrationHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineMockGetMigrationHandlerInput{p}, "StateMachine.GetMigrationHandler got unexpected parameters")
		}

		result := m.GetMigrationHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetMigrationHandler")
		}

		r = result.r

		return
	}

	if m.GetMigrationHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineMock.GetMigrationHandler. %v", p)
		return
	}

	return m.GetMigrationHandlerFunc(p)
}

//GetMigrationHandlerMinimockCounter returns a count of StateMachineMock.GetMigrationHandlerFunc invocations
func (m *StateMachineMock) GetMigrationHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMigrationHandlerCounter)
}

//GetMigrationHandlerMinimockPreCounter returns the value of StateMachineMock.GetMigrationHandler invocations
func (m *StateMachineMock) GetMigrationHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMigrationHandlerPreCounter)
}

//GetMigrationHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineMock) GetMigrationHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMigrationHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetMigrationHandlerCounter) == uint64(len(m.GetMigrationHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMigrationHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetMigrationHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetMigrationHandlerFunc != nil {
		return atomic.LoadUint64(&m.GetMigrationHandlerCounter) > 0
	}

	return true
}

type mStateMachineMockGetNestedHandler struct {
	mock              *StateMachineMock
	mainExpectation   *StateMachineMockGetNestedHandlerExpectation
	expectationSeries []*StateMachineMockGetNestedHandlerExpectation
}

type StateMachineMockGetNestedHandlerExpectation struct {
	input  *StateMachineMockGetNestedHandlerInput
	result *StateMachineMockGetNestedHandlerResult
}

type StateMachineMockGetNestedHandlerInput struct {
	p StateID
}

type StateMachineMockGetNestedHandlerResult struct {
	r NestedHandler
}

//Expect specifies that invocation of StateMachine.GetNestedHandler is expected from 1 to Infinity times
func (m *mStateMachineMockGetNestedHandler) Expect(p StateID) *mStateMachineMockGetNestedHandler {
	m.mock.GetNestedHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetNestedHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineMockGetNestedHandlerInput{p}
	return m
}

//Return specifies results of invocation of StateMachine.GetNestedHandler
func (m *mStateMachineMockGetNestedHandler) Return(r NestedHandler) *StateMachineMock {
	m.mock.GetNestedHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetNestedHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineMockGetNestedHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachine.GetNestedHandler is expected once
func (m *mStateMachineMockGetNestedHandler) ExpectOnce(p StateID) *StateMachineMockGetNestedHandlerExpectation {
	m.mock.GetNestedHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineMockGetNestedHandlerExpectation{}
	expectation.input = &StateMachineMockGetNestedHandlerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineMockGetNestedHandlerExpectation) Return(r NestedHandler) {
	e.result = &StateMachineMockGetNestedHandlerResult{r}
}

//Set uses given function f as a mock of StateMachine.GetNestedHandler method
func (m *mStateMachineMockGetNestedHandler) Set(f func(p StateID) (r NestedHandler)) *StateMachineMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNestedHandlerFunc = f
	return m.mock
}

//GetNestedHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine interface
func (m *StateMachineMock) GetNestedHandler(p StateID) (r NestedHandler) {
	counter := atomic.AddUint64(&m.GetNestedHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetNestedHandlerCounter, 1)

	if len(m.GetNestedHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNestedHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineMock.GetNestedHandler. %v", p)
			return
		}

		input := m.GetNestedHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineMockGetNestedHandlerInput{p}, "StateMachine.GetNestedHandler got unexpected parameters")

		result := m.GetNestedHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetNestedHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetNestedHandlerMock.mainExpectation != nil {

		input := m.GetNestedHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineMockGetNestedHandlerInput{p}, "StateMachine.GetNestedHandler got unexpected parameters")
		}

		result := m.GetNestedHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetNestedHandler")
		}

		r = result.r

		return
	}

	if m.GetNestedHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineMock.GetNestedHandler. %v", p)
		return
	}

	return m.GetNestedHandlerFunc(p)
}

//GetNestedHandlerMinimockCounter returns a count of StateMachineMock.GetNestedHandlerFunc invocations
func (m *StateMachineMock) GetNestedHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNestedHandlerCounter)
}

//GetNestedHandlerMinimockPreCounter returns the value of StateMachineMock.GetNestedHandler invocations
func (m *StateMachineMock) GetNestedHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNestedHandlerPreCounter)
}

//GetNestedHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineMock) GetNestedHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNestedHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNestedHandlerCounter) == uint64(len(m.GetNestedHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNestedHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNestedHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNestedHandlerFunc != nil {
		return atomic.LoadUint64(&m.GetNestedHandlerCounter) > 0
	}

	return true
}

type mStateMachineMockGetResponseErrorHandler struct {
	mock              *StateMachineMock
	mainExpectation   *StateMachineMockGetResponseErrorHandlerExpectation
	expectationSeries []*StateMachineMockGetResponseErrorHandlerExpectation
}

type StateMachineMockGetResponseErrorHandlerExpectation struct {
	input  *StateMachineMockGetResponseErrorHandlerInput
	result *StateMachineMockGetResponseErrorHandlerResult
}

type StateMachineMockGetResponseErrorHandlerInput struct {
	p StateID
}

type StateMachineMockGetResponseErrorHandlerResult struct {
	r ResponseErrorHandler
}

//Expect specifies that invocation of StateMachine.GetResponseErrorHandler is expected from 1 to Infinity times
func (m *mStateMachineMockGetResponseErrorHandler) Expect(p StateID) *mStateMachineMockGetResponseErrorHandler {
	m.mock.GetResponseErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetResponseErrorHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineMockGetResponseErrorHandlerInput{p}
	return m
}

//Return specifies results of invocation of StateMachine.GetResponseErrorHandler
func (m *mStateMachineMockGetResponseErrorHandler) Return(r ResponseErrorHandler) *StateMachineMock {
	m.mock.GetResponseErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetResponseErrorHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineMockGetResponseErrorHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachine.GetResponseErrorHandler is expected once
func (m *mStateMachineMockGetResponseErrorHandler) ExpectOnce(p StateID) *StateMachineMockGetResponseErrorHandlerExpectation {
	m.mock.GetResponseErrorHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineMockGetResponseErrorHandlerExpectation{}
	expectation.input = &StateMachineMockGetResponseErrorHandlerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineMockGetResponseErrorHandlerExpectation) Return(r ResponseErrorHandler) {
	e.result = &StateMachineMockGetResponseErrorHandlerResult{r}
}

//Set uses given function f as a mock of StateMachine.GetResponseErrorHandler method
func (m *mStateMachineMockGetResponseErrorHandler) Set(f func(p StateID) (r ResponseErrorHandler)) *StateMachineMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetResponseErrorHandlerFunc = f
	return m.mock
}

//GetResponseErrorHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine interface
func (m *StateMachineMock) GetResponseErrorHandler(p StateID) (r ResponseErrorHandler) {
	counter := atomic.AddUint64(&m.GetResponseErrorHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetResponseErrorHandlerCounter, 1)

	if len(m.GetResponseErrorHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetResponseErrorHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineMock.GetResponseErrorHandler. %v", p)
			return
		}

		input := m.GetResponseErrorHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineMockGetResponseErrorHandlerInput{p}, "StateMachine.GetResponseErrorHandler got unexpected parameters")

		result := m.GetResponseErrorHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetResponseErrorHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetResponseErrorHandlerMock.mainExpectation != nil {

		input := m.GetResponseErrorHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineMockGetResponseErrorHandlerInput{p}, "StateMachine.GetResponseErrorHandler got unexpected parameters")
		}

		result := m.GetResponseErrorHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetResponseErrorHandler")
		}

		r = result.r

		return
	}

	if m.GetResponseErrorHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineMock.GetResponseErrorHandler. %v", p)
		return
	}

	return m.GetResponseErrorHandlerFunc(p)
}

//GetResponseErrorHandlerMinimockCounter returns a count of StateMachineMock.GetResponseErrorHandlerFunc invocations
func (m *StateMachineMock) GetResponseErrorHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseErrorHandlerCounter)
}

//GetResponseErrorHandlerMinimockPreCounter returns the value of StateMachineMock.GetResponseErrorHandler invocations
func (m *StateMachineMock) GetResponseErrorHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseErrorHandlerPreCounter)
}

//GetResponseErrorHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineMock) GetResponseErrorHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetResponseErrorHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetResponseErrorHandlerCounter) == uint64(len(m.GetResponseErrorHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetResponseErrorHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetResponseErrorHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetResponseErrorHandlerFunc != nil {
		return atomic.LoadUint64(&m.GetResponseErrorHandlerCounter) > 0
	}

	return true
}

type mStateMachineMockGetResponseHandler struct {
	mock              *StateMachineMock
	mainExpectation   *StateMachineMockGetResponseHandlerExpectation
	expectationSeries []*StateMachineMockGetResponseHandlerExpectation
}

type StateMachineMockGetResponseHandlerExpectation struct {
	input  *StateMachineMockGetResponseHandlerInput
	result *StateMachineMockGetResponseHandlerResult
}

type StateMachineMockGetResponseHandlerInput struct {
	p StateID
}

type StateMachineMockGetResponseHandlerResult struct {
	r AdapterResponseHandler
}

//Expect specifies that invocation of StateMachine.GetResponseHandler is expected from 1 to Infinity times
func (m *mStateMachineMockGetResponseHandler) Expect(p StateID) *mStateMachineMockGetResponseHandler {
	m.mock.GetResponseHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetResponseHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineMockGetResponseHandlerInput{p}
	return m
}

//Return specifies results of invocation of StateMachine.GetResponseHandler
func (m *mStateMachineMockGetResponseHandler) Return(r AdapterResponseHandler) *StateMachineMock {
	m.mock.GetResponseHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetResponseHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineMockGetResponseHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachine.GetResponseHandler is expected once
func (m *mStateMachineMockGetResponseHandler) ExpectOnce(p StateID) *StateMachineMockGetResponseHandlerExpectation {
	m.mock.GetResponseHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineMockGetResponseHandlerExpectation{}
	expectation.input = &StateMachineMockGetResponseHandlerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineMockGetResponseHandlerExpectation) Return(r AdapterResponseHandler) {
	e.result = &StateMachineMockGetResponseHandlerResult{r}
}

//Set uses given function f as a mock of StateMachine.GetResponseHandler method
func (m *mStateMachineMockGetResponseHandler) Set(f func(p StateID) (r AdapterResponseHandler)) *StateMachineMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetResponseHandlerFunc = f
	return m.mock
}

//GetResponseHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine interface
func (m *StateMachineMock) GetResponseHandler(p StateID) (r AdapterResponseHandler) {
	counter := atomic.AddUint64(&m.GetResponseHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetResponseHandlerCounter, 1)

	if len(m.GetResponseHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetResponseHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineMock.GetResponseHandler. %v", p)
			return
		}

		input := m.GetResponseHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineMockGetResponseHandlerInput{p}, "StateMachine.GetResponseHandler got unexpected parameters")

		result := m.GetResponseHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetResponseHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetResponseHandlerMock.mainExpectation != nil {

		input := m.GetResponseHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineMockGetResponseHandlerInput{p}, "StateMachine.GetResponseHandler got unexpected parameters")
		}

		result := m.GetResponseHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetResponseHandler")
		}

		r = result.r

		return
	}

	if m.GetResponseHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineMock.GetResponseHandler. %v", p)
		return
	}

	return m.GetResponseHandlerFunc(p)
}

//GetResponseHandlerMinimockCounter returns a count of StateMachineMock.GetResponseHandlerFunc invocations
func (m *StateMachineMock) GetResponseHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseHandlerCounter)
}

//GetResponseHandlerMinimockPreCounter returns the value of StateMachineMock.GetResponseHandler invocations
func (m *StateMachineMock) GetResponseHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseHandlerPreCounter)
}

//GetResponseHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineMock) GetResponseHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetResponseHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetResponseHandlerCounter) == uint64(len(m.GetResponseHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetResponseHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetResponseHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetResponseHandlerFunc != nil {
		return atomic.LoadUint64(&m.GetResponseHandlerCounter) > 0
	}

	return true
}

type mStateMachineMockGetTransitionErrorHandler struct {
	mock              *StateMachineMock
	mainExpectation   *StateMachineMockGetTransitionErrorHandlerExpectation
	expectationSeries []*StateMachineMockGetTransitionErrorHandlerExpectation
}

type StateMachineMockGetTransitionErrorHandlerExpectation struct {
	input  *StateMachineMockGetTransitionErrorHandlerInput
	result *StateMachineMockGetTransitionErrorHandlerResult
}

type StateMachineMockGetTransitionErrorHandlerInput struct {
	p StateID
}

type StateMachineMockGetTransitionErrorHandlerResult struct {
	r TransitionErrorHandler
}

//Expect specifies that invocation of StateMachine.GetTransitionErrorHandler is expected from 1 to Infinity times
func (m *mStateMachineMockGetTransitionErrorHandler) Expect(p StateID) *mStateMachineMockGetTransitionErrorHandler {
	m.mock.GetTransitionErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetTransitionErrorHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineMockGetTransitionErrorHandlerInput{p}
	return m
}

//Return specifies results of invocation of StateMachine.GetTransitionErrorHandler
func (m *mStateMachineMockGetTransitionErrorHandler) Return(r TransitionErrorHandler) *StateMachineMock {
	m.mock.GetTransitionErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetTransitionErrorHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineMockGetTransitionErrorHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachine.GetTransitionErrorHandler is expected once
func (m *mStateMachineMockGetTransitionErrorHandler) ExpectOnce(p StateID) *StateMachineMockGetTransitionErrorHandlerExpectation {
	m.mock.GetTransitionErrorHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineMockGetTransitionErrorHandlerExpectation{}
	expectation.input = &StateMachineMockGetTransitionErrorHandlerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineMockGetTransitionErrorHandlerExpectation) Return(r TransitionErrorHandler) {
	e.result = &StateMachineMockGetTransitionErrorHandlerResult{r}
}

//Set uses given function f as a mock of StateMachine.GetTransitionErrorHandler method
func (m *mStateMachineMockGetTransitionErrorHandler) Set(f func(p StateID) (r TransitionErrorHandler)) *StateMachineMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTransitionErrorHandlerFunc = f
	return m.mock
}

//GetTransitionErrorHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine interface
func (m *StateMachineMock) GetTransitionErrorHandler(p StateID) (r TransitionErrorHandler) {
	counter := atomic.AddUint64(&m.GetTransitionErrorHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetTransitionErrorHandlerCounter, 1)

	if len(m.GetTransitionErrorHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTransitionErrorHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineMock.GetTransitionErrorHandler. %v", p)
			return
		}

		input := m.GetTransitionErrorHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineMockGetTransitionErrorHandlerInput{p}, "StateMachine.GetTransitionErrorHandler got unexpected parameters")

		result := m.GetTransitionErrorHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetTransitionErrorHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetTransitionErrorHandlerMock.mainExpectation != nil {

		input := m.GetTransitionErrorHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineMockGetTransitionErrorHandlerInput{p}, "StateMachine.GetTransitionErrorHandler got unexpected parameters")
		}

		result := m.GetTransitionErrorHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetTransitionErrorHandler")
		}

		r = result.r

		return
	}

	if m.GetTransitionErrorHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineMock.GetTransitionErrorHandler. %v", p)
		return
	}

	return m.GetTransitionErrorHandlerFunc(p)
}

//GetTransitionErrorHandlerMinimockCounter returns a count of StateMachineMock.GetTransitionErrorHandlerFunc invocations
func (m *StateMachineMock) GetTransitionErrorHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionErrorHandlerCounter)
}

//GetTransitionErrorHandlerMinimockPreCounter returns the value of StateMachineMock.GetTransitionErrorHandler invocations
func (m *StateMachineMock) GetTransitionErrorHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionErrorHandlerPreCounter)
}

//GetTransitionErrorHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineMock) GetTransitionErrorHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetTransitionErrorHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetTransitionErrorHandlerCounter) == uint64(len(m.GetTransitionErrorHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetTransitionErrorHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetTransitionErrorHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetTransitionErrorHandlerFunc != nil {
		return atomic.LoadUint64(&m.GetTransitionErrorHandlerCounter) > 0
	}

	return true
}

type mStateMachineMockGetTransitionHandler struct {
	mock              *StateMachineMock
	mainExpectation   *StateMachineMockGetTransitionHandlerExpectation
	expectationSeries []*StateMachineMockGetTransitionHandlerExpectation
}

type StateMachineMockGetTransitionHandlerExpectation struct {
	input  *StateMachineMockGetTransitionHandlerInput
	result *StateMachineMockGetTransitionHandlerResult
}

type StateMachineMockGetTransitionHandlerInput struct {
	p StateID
}

type StateMachineMockGetTransitionHandlerResult struct {
	r TransitHandler
}

//Expect specifies that invocation of StateMachine.GetTransitionHandler is expected from 1 to Infinity times
func (m *mStateMachineMockGetTransitionHandler) Expect(p StateID) *mStateMachineMockGetTransitionHandler {
	m.mock.GetTransitionHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetTransitionHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineMockGetTransitionHandlerInput{p}
	return m
}

//Return specifies results of invocation of StateMachine.GetTransitionHandler
func (m *mStateMachineMockGetTransitionHandler) Return(r TransitHandler) *StateMachineMock {
	m.mock.GetTransitionHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetTransitionHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineMockGetTransitionHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachine.GetTransitionHandler is expected once
func (m *mStateMachineMockGetTransitionHandler) ExpectOnce(p StateID) *StateMachineMockGetTransitionHandlerExpectation {
	m.mock.GetTransitionHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineMockGetTransitionHandlerExpectation{}
	expectation.input = &StateMachineMockGetTransitionHandlerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineMockGetTransitionHandlerExpectation) Return(r TransitHandler) {
	e.result = &StateMachineMockGetTransitionHandlerResult{r}
}

//Set uses given function f as a mock of StateMachine.GetTransitionHandler method
func (m *mStateMachineMockGetTransitionHandler) Set(f func(p StateID) (r TransitHandler)) *StateMachineMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTransitionHandlerFunc = f
	return m.mock
}

//GetTransitionHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine interface
func (m *StateMachineMock) GetTransitionHandler(p StateID) (r TransitHandler) {
	counter := atomic.AddUint64(&m.GetTransitionHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetTransitionHandlerCounter, 1)

	if len(m.GetTransitionHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTransitionHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineMock.GetTransitionHandler. %v", p)
			return
		}

		input := m.GetTransitionHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineMockGetTransitionHandlerInput{p}, "StateMachine.GetTransitionHandler got unexpected parameters")

		result := m.GetTransitionHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetTransitionHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetTransitionHandlerMock.mainExpectation != nil {

		input := m.GetTransitionHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineMockGetTransitionHandlerInput{p}, "StateMachine.GetTransitionHandler got unexpected parameters")
		}

		result := m.GetTransitionHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetTransitionHandler")
		}

		r = result.r

		return
	}

	if m.GetTransitionHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineMock.GetTransitionHandler. %v", p)
		return
	}

	return m.GetTransitionHandlerFunc(p)
}

//GetTransitionHandlerMinimockCounter returns a count of StateMachineMock.GetTransitionHandlerFunc invocations
func (m *StateMachineMock) GetTransitionHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionHandlerCounter)
}

//GetTransitionHandlerMinimockPreCounter returns the value of StateMachineMock.GetTransitionHandler invocations
func (m *StateMachineMock) GetTransitionHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionHandlerPreCounter)
}

//GetTransitionHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineMock) GetTransitionHandlerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetTransitionHandlerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetTransitionHandlerCounter) == uint64(len(m.GetTransitionHandlerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetTransitionHandlerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetTransitionHandlerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetTransitionHandlerFunc != nil {
		return atomic.LoadUint64(&m.GetTransitionHandlerCounter) > 0
	}

	return true
}

type mStateMachineMockGetTypeID struct {
	mock              *StateMachineMock
	mainExpectation   *StateMachineMockGetTypeIDExpectation
	expectationSeries []*StateMachineMockGetTypeIDExpectation
}

type StateMachineMockGetTypeIDExpectation struct {
	result *StateMachineMockGetTypeIDResult
}

type StateMachineMockGetTypeIDResult struct {
	r ID
}

//Expect specifies that invocation of StateMachine.GetTypeID is expected from 1 to Infinity times
func (m *mStateMachineMockGetTypeID) Expect() *mStateMachineMockGetTypeID {
	m.mock.GetTypeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetTypeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of StateMachine.GetTypeID
func (m *mStateMachineMockGetTypeID) Return(r ID) *StateMachineMock {
	m.mock.GetTypeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineMockGetTypeIDExpectation{}
	}
	m.mainExpectation.result = &StateMachineMockGetTypeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachine.GetTypeID is expected once
func (m *mStateMachineMockGetTypeID) ExpectOnce() *StateMachineMockGetTypeIDExpectation {
	m.mock.GetTypeIDFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineMockGetTypeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineMockGetTypeIDExpectation) Return(r ID) {
	e.result = &StateMachineMockGetTypeIDResult{r}
}

//Set uses given function f as a mock of StateMachine.GetTypeID method
func (m *mStateMachineMockGetTypeID) Set(f func() (r ID)) *StateMachineMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTypeIDFunc = f
	return m.mock
}

//GetTypeID implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine interface
func (m *StateMachineMock) GetTypeID() (r ID) {
	counter := atomic.AddUint64(&m.GetTypeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetTypeIDCounter, 1)

	if len(m.GetTypeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTypeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineMock.GetTypeID.")
			return
		}

		result := m.GetTypeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetTypeID")
			return
		}

		r = result.r

		return
	}

	if m.GetTypeIDMock.mainExpectation != nil {

		result := m.GetTypeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineMock.GetTypeID")
		}

		r = result.r

		return
	}

	if m.GetTypeIDFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineMock.GetTypeID.")
		return
	}

	return m.GetTypeIDFunc()
}

//GetTypeIDMinimockCounter returns a count of StateMachineMock.GetTypeIDFunc invocations
func (m *StateMachineMock) GetTypeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypeIDCounter)
}

//GetTypeIDMinimockPreCounter returns the value of StateMachineMock.GetTypeID invocations
func (m *StateMachineMock) GetTypeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypeIDPreCounter)
}

//GetTypeIDFinished returns true if mock invocations count is ok
func (m *StateMachineMock) GetTypeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetTypeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetTypeIDCounter) == uint64(len(m.GetTypeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetTypeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetTypeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetTypeIDFunc != nil {
		return atomic.LoadUint64(&m.GetTypeIDCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateMachineMock) ValidateCallCounters() {

	if !m.GetMigrationHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetMigrationHandler")
	}

	if !m.GetNestedHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetNestedHandler")
	}

	if !m.GetResponseErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetResponseErrorHandler")
	}

	if !m.GetResponseHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetResponseHandler")
	}

	if !m.GetTransitionErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetTransitionErrorHandler")
	}

	if !m.GetTransitionHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetTransitionHandler")
	}

	if !m.GetTypeIDFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetTypeID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateMachineMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StateMachineMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StateMachineMock) MinimockFinish() {

	if !m.GetMigrationHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetMigrationHandler")
	}

	if !m.GetNestedHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetNestedHandler")
	}

	if !m.GetResponseErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetResponseErrorHandler")
	}

	if !m.GetResponseHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetResponseHandler")
	}

	if !m.GetTransitionErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetTransitionErrorHandler")
	}

	if !m.GetTransitionHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetTransitionHandler")
	}

	if !m.GetTypeIDFinished() {
		m.t.Fatal("Expected call to StateMachineMock.GetTypeID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StateMachineMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StateMachineMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetMigrationHandlerFinished()
		ok = ok && m.GetNestedHandlerFinished()
		ok = ok && m.GetResponseErrorHandlerFinished()
		ok = ok && m.GetResponseHandlerFinished()
		ok = ok && m.GetTransitionErrorHandlerFinished()
		ok = ok && m.GetTransitionHandlerFinished()
		ok = ok && m.GetTypeIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetMigrationHandlerFinished() {
				m.t.Error("Expected call to StateMachineMock.GetMigrationHandler")
			}

			if !m.GetNestedHandlerFinished() {
				m.t.Error("Expected call to StateMachineMock.GetNestedHandler")
			}

			if !m.GetResponseErrorHandlerFinished() {
				m.t.Error("Expected call to StateMachineMock.GetResponseErrorHandler")
			}

			if !m.GetResponseHandlerFinished() {
				m.t.Error("Expected call to StateMachineMock.GetResponseHandler")
			}

			if !m.GetTransitionErrorHandlerFinished() {
				m.t.Error("Expected call to StateMachineMock.GetTransitionErrorHandler")
			}

			if !m.GetTransitionHandlerFinished() {
				m.t.Error("Expected call to StateMachineMock.GetTransitionHandler")
			}

			if !m.GetTypeIDFinished() {
				m.t.Error("Expected call to StateMachineMock.GetTypeID")
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
func (m *StateMachineMock) AllMocksCalled() bool {

	if !m.GetMigrationHandlerFinished() {
		return false
	}

	if !m.GetNestedHandlerFinished() {
		return false
	}

	if !m.GetResponseErrorHandlerFinished() {
		return false
	}

	if !m.GetResponseHandlerFinished() {
		return false
	}

	if !m.GetTransitionErrorHandlerFinished() {
		return false
	}

	if !m.GetTransitionHandlerFinished() {
		return false
	}

	if !m.GetTypeIDFinished() {
		return false
	}

	return true
}

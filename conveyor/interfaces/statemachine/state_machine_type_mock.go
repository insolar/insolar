package statemachine

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StateMachineType" can be found in github.com/insolar/insolar/conveyor/interfaces/statemachine
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//StateMachineTypeMock implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType
type StateMachineTypeMock struct {
	t minimock.Tester

	GetMigrationHandlerFunc       func(p SlotType, p1 uint32) (r MigrationHandler)
	GetMigrationHandlerCounter    uint64
	GetMigrationHandlerPreCounter uint64
	GetMigrationHandlerMock       mStateMachineTypeMockGetMigrationHandler

	GetNestedHandlerFunc       func(p SlotType, p1 uint32) (r NestedHandler)
	GetNestedHandlerCounter    uint64
	GetNestedHandlerPreCounter uint64
	GetNestedHandlerMock       mStateMachineTypeMockGetNestedHandler

	GetResponseErrorHandlerFunc       func(p SlotType, p1 uint32) (r ResponseErrorHandler)
	GetResponseErrorHandlerCounter    uint64
	GetResponseErrorHandlerPreCounter uint64
	GetResponseErrorHandlerMock       mStateMachineTypeMockGetResponseErrorHandler

	GetResponseHandlerFunc       func(p SlotType, p1 uint32) (r AdapterResponseHandler)
	GetResponseHandlerCounter    uint64
	GetResponseHandlerPreCounter uint64
	GetResponseHandlerMock       mStateMachineTypeMockGetResponseHandler

	GetTransitionErrorHandlerFunc       func(p SlotType, p1 uint32) (r TransitionErrorHandler)
	GetTransitionErrorHandlerCounter    uint64
	GetTransitionErrorHandlerPreCounter uint64
	GetTransitionErrorHandlerMock       mStateMachineTypeMockGetTransitionErrorHandler

	GetTransitionHandlerFunc       func(p SlotType, p1 uint32) (r TransitHandler)
	GetTransitionHandlerCounter    uint64
	GetTransitionHandlerPreCounter uint64
	GetTransitionHandlerMock       mStateMachineTypeMockGetTransitionHandler

	GetTypeIDFunc       func() (r int)
	GetTypeIDCounter    uint64
	GetTypeIDPreCounter uint64
	GetTypeIDMock       mStateMachineTypeMockGetTypeID
}

//NewStateMachineTypeMock returns a mock for github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType
func NewStateMachineTypeMock(t minimock.Tester) *StateMachineTypeMock {
	m := &StateMachineTypeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetMigrationHandlerMock = mStateMachineTypeMockGetMigrationHandler{mock: m}
	m.GetNestedHandlerMock = mStateMachineTypeMockGetNestedHandler{mock: m}
	m.GetResponseErrorHandlerMock = mStateMachineTypeMockGetResponseErrorHandler{mock: m}
	m.GetResponseHandlerMock = mStateMachineTypeMockGetResponseHandler{mock: m}
	m.GetTransitionErrorHandlerMock = mStateMachineTypeMockGetTransitionErrorHandler{mock: m}
	m.GetTransitionHandlerMock = mStateMachineTypeMockGetTransitionHandler{mock: m}
	m.GetTypeIDMock = mStateMachineTypeMockGetTypeID{mock: m}

	return m
}

type mStateMachineTypeMockGetMigrationHandler struct {
	mock              *StateMachineTypeMock
	mainExpectation   *StateMachineTypeMockGetMigrationHandlerExpectation
	expectationSeries []*StateMachineTypeMockGetMigrationHandlerExpectation
}

type StateMachineTypeMockGetMigrationHandlerExpectation struct {
	input  *StateMachineTypeMockGetMigrationHandlerInput
	result *StateMachineTypeMockGetMigrationHandlerResult
}

type StateMachineTypeMockGetMigrationHandlerInput struct {
	p  SlotType
	p1 uint32
}

type StateMachineTypeMockGetMigrationHandlerResult struct {
	r MigrationHandler
}

//Expect specifies that invocation of StateMachineType.GetMigrationHandler is expected from 1 to Infinity times
func (m *mStateMachineTypeMockGetMigrationHandler) Expect(p SlotType, p1 uint32) *mStateMachineTypeMockGetMigrationHandler {
	m.mock.GetMigrationHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetMigrationHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineTypeMockGetMigrationHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of StateMachineType.GetMigrationHandler
func (m *mStateMachineTypeMockGetMigrationHandler) Return(r MigrationHandler) *StateMachineTypeMock {
	m.mock.GetMigrationHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetMigrationHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineTypeMockGetMigrationHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachineType.GetMigrationHandler is expected once
func (m *mStateMachineTypeMockGetMigrationHandler) ExpectOnce(p SlotType, p1 uint32) *StateMachineTypeMockGetMigrationHandlerExpectation {
	m.mock.GetMigrationHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineTypeMockGetMigrationHandlerExpectation{}
	expectation.input = &StateMachineTypeMockGetMigrationHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineTypeMockGetMigrationHandlerExpectation) Return(r MigrationHandler) {
	e.result = &StateMachineTypeMockGetMigrationHandlerResult{r}
}

//Set uses given function f as a mock of StateMachineType.GetMigrationHandler method
func (m *mStateMachineTypeMockGetMigrationHandler) Set(f func(p SlotType, p1 uint32) (r MigrationHandler)) *StateMachineTypeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMigrationHandlerFunc = f
	return m.mock
}

//GetMigrationHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType interface
func (m *StateMachineTypeMock) GetMigrationHandler(p SlotType, p1 uint32) (r MigrationHandler) {
	counter := atomic.AddUint64(&m.GetMigrationHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetMigrationHandlerCounter, 1)

	if len(m.GetMigrationHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMigrationHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetMigrationHandler. %v %v", p, p1)
			return
		}

		input := m.GetMigrationHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineTypeMockGetMigrationHandlerInput{p, p1}, "StateMachineType.GetMigrationHandler got unexpected parameters")

		result := m.GetMigrationHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetMigrationHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetMigrationHandlerMock.mainExpectation != nil {

		input := m.GetMigrationHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineTypeMockGetMigrationHandlerInput{p, p1}, "StateMachineType.GetMigrationHandler got unexpected parameters")
		}

		result := m.GetMigrationHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetMigrationHandler")
		}

		r = result.r

		return
	}

	if m.GetMigrationHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetMigrationHandler. %v %v", p, p1)
		return
	}

	return m.GetMigrationHandlerFunc(p, p1)
}

//GetMigrationHandlerMinimockCounter returns a count of StateMachineTypeMock.GetMigrationHandlerFunc invocations
func (m *StateMachineTypeMock) GetMigrationHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMigrationHandlerCounter)
}

//GetMigrationHandlerMinimockPreCounter returns the value of StateMachineTypeMock.GetMigrationHandler invocations
func (m *StateMachineTypeMock) GetMigrationHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMigrationHandlerPreCounter)
}

//GetMigrationHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineTypeMock) GetMigrationHandlerFinished() bool {
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

type mStateMachineTypeMockGetNestedHandler struct {
	mock              *StateMachineTypeMock
	mainExpectation   *StateMachineTypeMockGetNestedHandlerExpectation
	expectationSeries []*StateMachineTypeMockGetNestedHandlerExpectation
}

type StateMachineTypeMockGetNestedHandlerExpectation struct {
	input  *StateMachineTypeMockGetNestedHandlerInput
	result *StateMachineTypeMockGetNestedHandlerResult
}

type StateMachineTypeMockGetNestedHandlerInput struct {
	p  SlotType
	p1 uint32
}

type StateMachineTypeMockGetNestedHandlerResult struct {
	r NestedHandler
}

//Expect specifies that invocation of StateMachineType.GetNestedHandler is expected from 1 to Infinity times
func (m *mStateMachineTypeMockGetNestedHandler) Expect(p SlotType, p1 uint32) *mStateMachineTypeMockGetNestedHandler {
	m.mock.GetNestedHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetNestedHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineTypeMockGetNestedHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of StateMachineType.GetNestedHandler
func (m *mStateMachineTypeMockGetNestedHandler) Return(r NestedHandler) *StateMachineTypeMock {
	m.mock.GetNestedHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetNestedHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineTypeMockGetNestedHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachineType.GetNestedHandler is expected once
func (m *mStateMachineTypeMockGetNestedHandler) ExpectOnce(p SlotType, p1 uint32) *StateMachineTypeMockGetNestedHandlerExpectation {
	m.mock.GetNestedHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineTypeMockGetNestedHandlerExpectation{}
	expectation.input = &StateMachineTypeMockGetNestedHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineTypeMockGetNestedHandlerExpectation) Return(r NestedHandler) {
	e.result = &StateMachineTypeMockGetNestedHandlerResult{r}
}

//Set uses given function f as a mock of StateMachineType.GetNestedHandler method
func (m *mStateMachineTypeMockGetNestedHandler) Set(f func(p SlotType, p1 uint32) (r NestedHandler)) *StateMachineTypeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNestedHandlerFunc = f
	return m.mock
}

//GetNestedHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType interface
func (m *StateMachineTypeMock) GetNestedHandler(p SlotType, p1 uint32) (r NestedHandler) {
	counter := atomic.AddUint64(&m.GetNestedHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetNestedHandlerCounter, 1)

	if len(m.GetNestedHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNestedHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetNestedHandler. %v %v", p, p1)
			return
		}

		input := m.GetNestedHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineTypeMockGetNestedHandlerInput{p, p1}, "StateMachineType.GetNestedHandler got unexpected parameters")

		result := m.GetNestedHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetNestedHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetNestedHandlerMock.mainExpectation != nil {

		input := m.GetNestedHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineTypeMockGetNestedHandlerInput{p, p1}, "StateMachineType.GetNestedHandler got unexpected parameters")
		}

		result := m.GetNestedHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetNestedHandler")
		}

		r = result.r

		return
	}

	if m.GetNestedHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetNestedHandler. %v %v", p, p1)
		return
	}

	return m.GetNestedHandlerFunc(p, p1)
}

//GetNestedHandlerMinimockCounter returns a count of StateMachineTypeMock.GetNestedHandlerFunc invocations
func (m *StateMachineTypeMock) GetNestedHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNestedHandlerCounter)
}

//GetNestedHandlerMinimockPreCounter returns the value of StateMachineTypeMock.GetNestedHandler invocations
func (m *StateMachineTypeMock) GetNestedHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNestedHandlerPreCounter)
}

//GetNestedHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineTypeMock) GetNestedHandlerFinished() bool {
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

type mStateMachineTypeMockGetResponseErrorHandler struct {
	mock              *StateMachineTypeMock
	mainExpectation   *StateMachineTypeMockGetResponseErrorHandlerExpectation
	expectationSeries []*StateMachineTypeMockGetResponseErrorHandlerExpectation
}

type StateMachineTypeMockGetResponseErrorHandlerExpectation struct {
	input  *StateMachineTypeMockGetResponseErrorHandlerInput
	result *StateMachineTypeMockGetResponseErrorHandlerResult
}

type StateMachineTypeMockGetResponseErrorHandlerInput struct {
	p  SlotType
	p1 uint32
}

type StateMachineTypeMockGetResponseErrorHandlerResult struct {
	r ResponseErrorHandler
}

//Expect specifies that invocation of StateMachineType.GetResponseErrorHandler is expected from 1 to Infinity times
func (m *mStateMachineTypeMockGetResponseErrorHandler) Expect(p SlotType, p1 uint32) *mStateMachineTypeMockGetResponseErrorHandler {
	m.mock.GetResponseErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetResponseErrorHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineTypeMockGetResponseErrorHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of StateMachineType.GetResponseErrorHandler
func (m *mStateMachineTypeMockGetResponseErrorHandler) Return(r ResponseErrorHandler) *StateMachineTypeMock {
	m.mock.GetResponseErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetResponseErrorHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineTypeMockGetResponseErrorHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachineType.GetResponseErrorHandler is expected once
func (m *mStateMachineTypeMockGetResponseErrorHandler) ExpectOnce(p SlotType, p1 uint32) *StateMachineTypeMockGetResponseErrorHandlerExpectation {
	m.mock.GetResponseErrorHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineTypeMockGetResponseErrorHandlerExpectation{}
	expectation.input = &StateMachineTypeMockGetResponseErrorHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineTypeMockGetResponseErrorHandlerExpectation) Return(r ResponseErrorHandler) {
	e.result = &StateMachineTypeMockGetResponseErrorHandlerResult{r}
}

//Set uses given function f as a mock of StateMachineType.GetResponseErrorHandler method
func (m *mStateMachineTypeMockGetResponseErrorHandler) Set(f func(p SlotType, p1 uint32) (r ResponseErrorHandler)) *StateMachineTypeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetResponseErrorHandlerFunc = f
	return m.mock
}

//GetResponseErrorHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType interface
func (m *StateMachineTypeMock) GetResponseErrorHandler(p SlotType, p1 uint32) (r ResponseErrorHandler) {
	counter := atomic.AddUint64(&m.GetResponseErrorHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetResponseErrorHandlerCounter, 1)

	if len(m.GetResponseErrorHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetResponseErrorHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetResponseErrorHandler. %v %v", p, p1)
			return
		}

		input := m.GetResponseErrorHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineTypeMockGetResponseErrorHandlerInput{p, p1}, "StateMachineType.GetResponseErrorHandler got unexpected parameters")

		result := m.GetResponseErrorHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetResponseErrorHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetResponseErrorHandlerMock.mainExpectation != nil {

		input := m.GetResponseErrorHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineTypeMockGetResponseErrorHandlerInput{p, p1}, "StateMachineType.GetResponseErrorHandler got unexpected parameters")
		}

		result := m.GetResponseErrorHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetResponseErrorHandler")
		}

		r = result.r

		return
	}

	if m.GetResponseErrorHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetResponseErrorHandler. %v %v", p, p1)
		return
	}

	return m.GetResponseErrorHandlerFunc(p, p1)
}

//GetResponseErrorHandlerMinimockCounter returns a count of StateMachineTypeMock.GetResponseErrorHandlerFunc invocations
func (m *StateMachineTypeMock) GetResponseErrorHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseErrorHandlerCounter)
}

//GetResponseErrorHandlerMinimockPreCounter returns the value of StateMachineTypeMock.GetResponseErrorHandler invocations
func (m *StateMachineTypeMock) GetResponseErrorHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseErrorHandlerPreCounter)
}

//GetResponseErrorHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineTypeMock) GetResponseErrorHandlerFinished() bool {
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

type mStateMachineTypeMockGetResponseHandler struct {
	mock              *StateMachineTypeMock
	mainExpectation   *StateMachineTypeMockGetResponseHandlerExpectation
	expectationSeries []*StateMachineTypeMockGetResponseHandlerExpectation
}

type StateMachineTypeMockGetResponseHandlerExpectation struct {
	input  *StateMachineTypeMockGetResponseHandlerInput
	result *StateMachineTypeMockGetResponseHandlerResult
}

type StateMachineTypeMockGetResponseHandlerInput struct {
	p  SlotType
	p1 uint32
}

type StateMachineTypeMockGetResponseHandlerResult struct {
	r AdapterResponseHandler
}

//Expect specifies that invocation of StateMachineType.GetResponseHandler is expected from 1 to Infinity times
func (m *mStateMachineTypeMockGetResponseHandler) Expect(p SlotType, p1 uint32) *mStateMachineTypeMockGetResponseHandler {
	m.mock.GetResponseHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetResponseHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineTypeMockGetResponseHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of StateMachineType.GetResponseHandler
func (m *mStateMachineTypeMockGetResponseHandler) Return(r AdapterResponseHandler) *StateMachineTypeMock {
	m.mock.GetResponseHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetResponseHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineTypeMockGetResponseHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachineType.GetResponseHandler is expected once
func (m *mStateMachineTypeMockGetResponseHandler) ExpectOnce(p SlotType, p1 uint32) *StateMachineTypeMockGetResponseHandlerExpectation {
	m.mock.GetResponseHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineTypeMockGetResponseHandlerExpectation{}
	expectation.input = &StateMachineTypeMockGetResponseHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineTypeMockGetResponseHandlerExpectation) Return(r AdapterResponseHandler) {
	e.result = &StateMachineTypeMockGetResponseHandlerResult{r}
}

//Set uses given function f as a mock of StateMachineType.GetResponseHandler method
func (m *mStateMachineTypeMockGetResponseHandler) Set(f func(p SlotType, p1 uint32) (r AdapterResponseHandler)) *StateMachineTypeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetResponseHandlerFunc = f
	return m.mock
}

//GetResponseHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType interface
func (m *StateMachineTypeMock) GetResponseHandler(p SlotType, p1 uint32) (r AdapterResponseHandler) {
	counter := atomic.AddUint64(&m.GetResponseHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetResponseHandlerCounter, 1)

	if len(m.GetResponseHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetResponseHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetResponseHandler. %v %v", p, p1)
			return
		}

		input := m.GetResponseHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineTypeMockGetResponseHandlerInput{p, p1}, "StateMachineType.GetResponseHandler got unexpected parameters")

		result := m.GetResponseHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetResponseHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetResponseHandlerMock.mainExpectation != nil {

		input := m.GetResponseHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineTypeMockGetResponseHandlerInput{p, p1}, "StateMachineType.GetResponseHandler got unexpected parameters")
		}

		result := m.GetResponseHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetResponseHandler")
		}

		r = result.r

		return
	}

	if m.GetResponseHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetResponseHandler. %v %v", p, p1)
		return
	}

	return m.GetResponseHandlerFunc(p, p1)
}

//GetResponseHandlerMinimockCounter returns a count of StateMachineTypeMock.GetResponseHandlerFunc invocations
func (m *StateMachineTypeMock) GetResponseHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseHandlerCounter)
}

//GetResponseHandlerMinimockPreCounter returns the value of StateMachineTypeMock.GetResponseHandler invocations
func (m *StateMachineTypeMock) GetResponseHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetResponseHandlerPreCounter)
}

//GetResponseHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineTypeMock) GetResponseHandlerFinished() bool {
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

type mStateMachineTypeMockGetTransitionErrorHandler struct {
	mock              *StateMachineTypeMock
	mainExpectation   *StateMachineTypeMockGetTransitionErrorHandlerExpectation
	expectationSeries []*StateMachineTypeMockGetTransitionErrorHandlerExpectation
}

type StateMachineTypeMockGetTransitionErrorHandlerExpectation struct {
	input  *StateMachineTypeMockGetTransitionErrorHandlerInput
	result *StateMachineTypeMockGetTransitionErrorHandlerResult
}

type StateMachineTypeMockGetTransitionErrorHandlerInput struct {
	p  SlotType
	p1 uint32
}

type StateMachineTypeMockGetTransitionErrorHandlerResult struct {
	r TransitionErrorHandler
}

//Expect specifies that invocation of StateMachineType.GetTransitionErrorHandler is expected from 1 to Infinity times
func (m *mStateMachineTypeMockGetTransitionErrorHandler) Expect(p SlotType, p1 uint32) *mStateMachineTypeMockGetTransitionErrorHandler {
	m.mock.GetTransitionErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetTransitionErrorHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineTypeMockGetTransitionErrorHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of StateMachineType.GetTransitionErrorHandler
func (m *mStateMachineTypeMockGetTransitionErrorHandler) Return(r TransitionErrorHandler) *StateMachineTypeMock {
	m.mock.GetTransitionErrorHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetTransitionErrorHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineTypeMockGetTransitionErrorHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachineType.GetTransitionErrorHandler is expected once
func (m *mStateMachineTypeMockGetTransitionErrorHandler) ExpectOnce(p SlotType, p1 uint32) *StateMachineTypeMockGetTransitionErrorHandlerExpectation {
	m.mock.GetTransitionErrorHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineTypeMockGetTransitionErrorHandlerExpectation{}
	expectation.input = &StateMachineTypeMockGetTransitionErrorHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineTypeMockGetTransitionErrorHandlerExpectation) Return(r TransitionErrorHandler) {
	e.result = &StateMachineTypeMockGetTransitionErrorHandlerResult{r}
}

//Set uses given function f as a mock of StateMachineType.GetTransitionErrorHandler method
func (m *mStateMachineTypeMockGetTransitionErrorHandler) Set(f func(p SlotType, p1 uint32) (r TransitionErrorHandler)) *StateMachineTypeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTransitionErrorHandlerFunc = f
	return m.mock
}

//GetTransitionErrorHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType interface
func (m *StateMachineTypeMock) GetTransitionErrorHandler(p SlotType, p1 uint32) (r TransitionErrorHandler) {
	counter := atomic.AddUint64(&m.GetTransitionErrorHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetTransitionErrorHandlerCounter, 1)

	if len(m.GetTransitionErrorHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTransitionErrorHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetTransitionErrorHandler. %v %v", p, p1)
			return
		}

		input := m.GetTransitionErrorHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineTypeMockGetTransitionErrorHandlerInput{p, p1}, "StateMachineType.GetTransitionErrorHandler got unexpected parameters")

		result := m.GetTransitionErrorHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetTransitionErrorHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetTransitionErrorHandlerMock.mainExpectation != nil {

		input := m.GetTransitionErrorHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineTypeMockGetTransitionErrorHandlerInput{p, p1}, "StateMachineType.GetTransitionErrorHandler got unexpected parameters")
		}

		result := m.GetTransitionErrorHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetTransitionErrorHandler")
		}

		r = result.r

		return
	}

	if m.GetTransitionErrorHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetTransitionErrorHandler. %v %v", p, p1)
		return
	}

	return m.GetTransitionErrorHandlerFunc(p, p1)
}

//GetTransitionErrorHandlerMinimockCounter returns a count of StateMachineTypeMock.GetTransitionErrorHandlerFunc invocations
func (m *StateMachineTypeMock) GetTransitionErrorHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionErrorHandlerCounter)
}

//GetTransitionErrorHandlerMinimockPreCounter returns the value of StateMachineTypeMock.GetTransitionErrorHandler invocations
func (m *StateMachineTypeMock) GetTransitionErrorHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionErrorHandlerPreCounter)
}

//GetTransitionErrorHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineTypeMock) GetTransitionErrorHandlerFinished() bool {
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

type mStateMachineTypeMockGetTransitionHandler struct {
	mock              *StateMachineTypeMock
	mainExpectation   *StateMachineTypeMockGetTransitionHandlerExpectation
	expectationSeries []*StateMachineTypeMockGetTransitionHandlerExpectation
}

type StateMachineTypeMockGetTransitionHandlerExpectation struct {
	input  *StateMachineTypeMockGetTransitionHandlerInput
	result *StateMachineTypeMockGetTransitionHandlerResult
}

type StateMachineTypeMockGetTransitionHandlerInput struct {
	p  SlotType
	p1 uint32
}

type StateMachineTypeMockGetTransitionHandlerResult struct {
	r TransitHandler
}

//Expect specifies that invocation of StateMachineType.GetTransitionHandler is expected from 1 to Infinity times
func (m *mStateMachineTypeMockGetTransitionHandler) Expect(p SlotType, p1 uint32) *mStateMachineTypeMockGetTransitionHandler {
	m.mock.GetTransitionHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetTransitionHandlerExpectation{}
	}
	m.mainExpectation.input = &StateMachineTypeMockGetTransitionHandlerInput{p, p1}
	return m
}

//Return specifies results of invocation of StateMachineType.GetTransitionHandler
func (m *mStateMachineTypeMockGetTransitionHandler) Return(r TransitHandler) *StateMachineTypeMock {
	m.mock.GetTransitionHandlerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetTransitionHandlerExpectation{}
	}
	m.mainExpectation.result = &StateMachineTypeMockGetTransitionHandlerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachineType.GetTransitionHandler is expected once
func (m *mStateMachineTypeMockGetTransitionHandler) ExpectOnce(p SlotType, p1 uint32) *StateMachineTypeMockGetTransitionHandlerExpectation {
	m.mock.GetTransitionHandlerFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineTypeMockGetTransitionHandlerExpectation{}
	expectation.input = &StateMachineTypeMockGetTransitionHandlerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineTypeMockGetTransitionHandlerExpectation) Return(r TransitHandler) {
	e.result = &StateMachineTypeMockGetTransitionHandlerResult{r}
}

//Set uses given function f as a mock of StateMachineType.GetTransitionHandler method
func (m *mStateMachineTypeMockGetTransitionHandler) Set(f func(p SlotType, p1 uint32) (r TransitHandler)) *StateMachineTypeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTransitionHandlerFunc = f
	return m.mock
}

//GetTransitionHandler implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType interface
func (m *StateMachineTypeMock) GetTransitionHandler(p SlotType, p1 uint32) (r TransitHandler) {
	counter := atomic.AddUint64(&m.GetTransitionHandlerPreCounter, 1)
	defer atomic.AddUint64(&m.GetTransitionHandlerCounter, 1)

	if len(m.GetTransitionHandlerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTransitionHandlerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetTransitionHandler. %v %v", p, p1)
			return
		}

		input := m.GetTransitionHandlerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateMachineTypeMockGetTransitionHandlerInput{p, p1}, "StateMachineType.GetTransitionHandler got unexpected parameters")

		result := m.GetTransitionHandlerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetTransitionHandler")
			return
		}

		r = result.r

		return
	}

	if m.GetTransitionHandlerMock.mainExpectation != nil {

		input := m.GetTransitionHandlerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateMachineTypeMockGetTransitionHandlerInput{p, p1}, "StateMachineType.GetTransitionHandler got unexpected parameters")
		}

		result := m.GetTransitionHandlerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetTransitionHandler")
		}

		r = result.r

		return
	}

	if m.GetTransitionHandlerFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetTransitionHandler. %v %v", p, p1)
		return
	}

	return m.GetTransitionHandlerFunc(p, p1)
}

//GetTransitionHandlerMinimockCounter returns a count of StateMachineTypeMock.GetTransitionHandlerFunc invocations
func (m *StateMachineTypeMock) GetTransitionHandlerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionHandlerCounter)
}

//GetTransitionHandlerMinimockPreCounter returns the value of StateMachineTypeMock.GetTransitionHandler invocations
func (m *StateMachineTypeMock) GetTransitionHandlerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransitionHandlerPreCounter)
}

//GetTransitionHandlerFinished returns true if mock invocations count is ok
func (m *StateMachineTypeMock) GetTransitionHandlerFinished() bool {
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

type mStateMachineTypeMockGetTypeID struct {
	mock              *StateMachineTypeMock
	mainExpectation   *StateMachineTypeMockGetTypeIDExpectation
	expectationSeries []*StateMachineTypeMockGetTypeIDExpectation
}

type StateMachineTypeMockGetTypeIDExpectation struct {
	result *StateMachineTypeMockGetTypeIDResult
}

type StateMachineTypeMockGetTypeIDResult struct {
	r int
}

//Expect specifies that invocation of StateMachineType.GetTypeID is expected from 1 to Infinity times
func (m *mStateMachineTypeMockGetTypeID) Expect() *mStateMachineTypeMockGetTypeID {
	m.mock.GetTypeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetTypeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of StateMachineType.GetTypeID
func (m *mStateMachineTypeMockGetTypeID) Return(r int) *StateMachineTypeMock {
	m.mock.GetTypeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateMachineTypeMockGetTypeIDExpectation{}
	}
	m.mainExpectation.result = &StateMachineTypeMockGetTypeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateMachineType.GetTypeID is expected once
func (m *mStateMachineTypeMockGetTypeID) ExpectOnce() *StateMachineTypeMockGetTypeIDExpectation {
	m.mock.GetTypeIDFunc = nil
	m.mainExpectation = nil

	expectation := &StateMachineTypeMockGetTypeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateMachineTypeMockGetTypeIDExpectation) Return(r int) {
	e.result = &StateMachineTypeMockGetTypeIDResult{r}
}

//Set uses given function f as a mock of StateMachineType.GetTypeID method
func (m *mStateMachineTypeMockGetTypeID) Set(f func() (r int)) *StateMachineTypeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTypeIDFunc = f
	return m.mock
}

//GetTypeID implements github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType interface
func (m *StateMachineTypeMock) GetTypeID() (r int) {
	counter := atomic.AddUint64(&m.GetTypeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetTypeIDCounter, 1)

	if len(m.GetTypeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTypeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetTypeID.")
			return
		}

		result := m.GetTypeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetTypeID")
			return
		}

		r = result.r

		return
	}

	if m.GetTypeIDMock.mainExpectation != nil {

		result := m.GetTypeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateMachineTypeMock.GetTypeID")
		}

		r = result.r

		return
	}

	if m.GetTypeIDFunc == nil {
		m.t.Fatalf("Unexpected call to StateMachineTypeMock.GetTypeID.")
		return
	}

	return m.GetTypeIDFunc()
}

//GetTypeIDMinimockCounter returns a count of StateMachineTypeMock.GetTypeIDFunc invocations
func (m *StateMachineTypeMock) GetTypeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypeIDCounter)
}

//GetTypeIDMinimockPreCounter returns the value of StateMachineTypeMock.GetTypeID invocations
func (m *StateMachineTypeMock) GetTypeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTypeIDPreCounter)
}

//GetTypeIDFinished returns true if mock invocations count is ok
func (m *StateMachineTypeMock) GetTypeIDFinished() bool {
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
func (m *StateMachineTypeMock) ValidateCallCounters() {

	if !m.GetMigrationHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetMigrationHandler")
	}

	if !m.GetNestedHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetNestedHandler")
	}

	if !m.GetResponseErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetResponseErrorHandler")
	}

	if !m.GetResponseHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetResponseHandler")
	}

	if !m.GetTransitionErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetTransitionErrorHandler")
	}

	if !m.GetTransitionHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetTransitionHandler")
	}

	if !m.GetTypeIDFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetTypeID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateMachineTypeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StateMachineTypeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StateMachineTypeMock) MinimockFinish() {

	if !m.GetMigrationHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetMigrationHandler")
	}

	if !m.GetNestedHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetNestedHandler")
	}

	if !m.GetResponseErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetResponseErrorHandler")
	}

	if !m.GetResponseHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetResponseHandler")
	}

	if !m.GetTransitionErrorHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetTransitionErrorHandler")
	}

	if !m.GetTransitionHandlerFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetTransitionHandler")
	}

	if !m.GetTypeIDFinished() {
		m.t.Fatal("Expected call to StateMachineTypeMock.GetTypeID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StateMachineTypeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StateMachineTypeMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to StateMachineTypeMock.GetMigrationHandler")
			}

			if !m.GetNestedHandlerFinished() {
				m.t.Error("Expected call to StateMachineTypeMock.GetNestedHandler")
			}

			if !m.GetResponseErrorHandlerFinished() {
				m.t.Error("Expected call to StateMachineTypeMock.GetResponseErrorHandler")
			}

			if !m.GetResponseHandlerFinished() {
				m.t.Error("Expected call to StateMachineTypeMock.GetResponseHandler")
			}

			if !m.GetTransitionErrorHandlerFinished() {
				m.t.Error("Expected call to StateMachineTypeMock.GetTransitionErrorHandler")
			}

			if !m.GetTransitionHandlerFinished() {
				m.t.Error("Expected call to StateMachineTypeMock.GetTransitionHandler")
			}

			if !m.GetTypeIDFinished() {
				m.t.Error("Expected call to StateMachineTypeMock.GetTypeID")
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
func (m *StateMachineTypeMock) AllMocksCalled() bool {

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

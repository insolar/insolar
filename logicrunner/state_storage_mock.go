package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StateStorage" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//StateStorageMock implements github.com/insolar/insolar/logicrunner.StateStorage
type StateStorageMock struct {
	t minimock.Tester

	DeleteObjectStateFunc       func(p insolar.Reference)
	DeleteObjectStateCounter    uint64
	DeleteObjectStatePreCounter uint64
	DeleteObjectStateMock       mStateStorageMockDeleteObjectState

	GetExecutionStateFunc       func(p insolar.Reference) (r *ExecutionBroker)
	GetExecutionStateCounter    uint64
	GetExecutionStatePreCounter uint64
	GetExecutionStateMock       mStateStorageMockGetExecutionState

	GetValidationStateFunc       func(p insolar.Reference) (r *ExecutionState)
	GetValidationStateCounter    uint64
	GetValidationStatePreCounter uint64
	GetValidationStateMock       mStateStorageMockGetValidationState

	LockFunc       func()
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mStateStorageMockLock

	StateMapFunc       func() (r *map[insolar.Reference]*ObjectState)
	StateMapCounter    uint64
	StateMapPreCounter uint64
	StateMapMock       mStateStorageMockStateMap

	UnlockFunc       func()
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mStateStorageMockUnlock

	UpsertExecutionStateFunc       func(p insolar.Reference) (r *ExecutionBroker)
	UpsertExecutionStateCounter    uint64
	UpsertExecutionStatePreCounter uint64
	UpsertExecutionStateMock       mStateStorageMockUpsertExecutionState

	UpsertValidationStateFunc       func(p insolar.Reference) (r *ExecutionState)
	UpsertValidationStateCounter    uint64
	UpsertValidationStatePreCounter uint64
	UpsertValidationStateMock       mStateStorageMockUpsertValidationState
}

//NewStateStorageMock returns a mock for github.com/insolar/insolar/logicrunner.StateStorage
func NewStateStorageMock(t minimock.Tester) *StateStorageMock {
	m := &StateStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteObjectStateMock = mStateStorageMockDeleteObjectState{mock: m}
	m.GetExecutionStateMock = mStateStorageMockGetExecutionState{mock: m}
	m.GetValidationStateMock = mStateStorageMockGetValidationState{mock: m}
	m.LockMock = mStateStorageMockLock{mock: m}
	m.StateMapMock = mStateStorageMockStateMap{mock: m}
	m.UnlockMock = mStateStorageMockUnlock{mock: m}
	m.UpsertExecutionStateMock = mStateStorageMockUpsertExecutionState{mock: m}
	m.UpsertValidationStateMock = mStateStorageMockUpsertValidationState{mock: m}

	return m
}

type mStateStorageMockDeleteObjectState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockDeleteObjectStateExpectation
	expectationSeries []*StateStorageMockDeleteObjectStateExpectation
}

type StateStorageMockDeleteObjectStateExpectation struct {
	input *StateStorageMockDeleteObjectStateInput
}

type StateStorageMockDeleteObjectStateInput struct {
	p insolar.Reference
}

//Expect specifies that invocation of StateStorage.DeleteObjectState is expected from 1 to Infinity times
func (m *mStateStorageMockDeleteObjectState) Expect(p insolar.Reference) *mStateStorageMockDeleteObjectState {
	m.mock.DeleteObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockDeleteObjectStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockDeleteObjectStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.DeleteObjectState
func (m *mStateStorageMockDeleteObjectState) Return() *StateStorageMock {
	m.mock.DeleteObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockDeleteObjectStateExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.DeleteObjectState is expected once
func (m *mStateStorageMockDeleteObjectState) ExpectOnce(p insolar.Reference) *StateStorageMockDeleteObjectStateExpectation {
	m.mock.DeleteObjectStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockDeleteObjectStateExpectation{}
	expectation.input = &StateStorageMockDeleteObjectStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of StateStorage.DeleteObjectState method
func (m *mStateStorageMockDeleteObjectState) Set(f func(p insolar.Reference)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteObjectStateFunc = f
	return m.mock
}

//DeleteObjectState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) DeleteObjectState(p insolar.Reference) {
	counter := atomic.AddUint64(&m.DeleteObjectStatePreCounter, 1)
	defer atomic.AddUint64(&m.DeleteObjectStateCounter, 1)

	if len(m.DeleteObjectStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteObjectStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.DeleteObjectState. %v", p)
			return
		}

		input := m.DeleteObjectStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockDeleteObjectStateInput{p}, "StateStorage.DeleteObjectState got unexpected parameters")

		return
	}

	if m.DeleteObjectStateMock.mainExpectation != nil {

		input := m.DeleteObjectStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockDeleteObjectStateInput{p}, "StateStorage.DeleteObjectState got unexpected parameters")
		}

		return
	}

	if m.DeleteObjectStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.DeleteObjectState. %v", p)
		return
	}

	m.DeleteObjectStateFunc(p)
}

//DeleteObjectStateMinimockCounter returns a count of StateStorageMock.DeleteObjectStateFunc invocations
func (m *StateStorageMock) DeleteObjectStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteObjectStateCounter)
}

//DeleteObjectStateMinimockPreCounter returns the value of StateStorageMock.DeleteObjectState invocations
func (m *StateStorageMock) DeleteObjectStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteObjectStatePreCounter)
}

//DeleteObjectStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) DeleteObjectStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeleteObjectStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeleteObjectStateCounter) == uint64(len(m.DeleteObjectStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeleteObjectStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeleteObjectStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeleteObjectStateFunc != nil {
		return atomic.LoadUint64(&m.DeleteObjectStateCounter) > 0
	}

	return true
}

type mStateStorageMockGetExecutionState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockGetExecutionStateExpectation
	expectationSeries []*StateStorageMockGetExecutionStateExpectation
}

type StateStorageMockGetExecutionStateExpectation struct {
	input  *StateStorageMockGetExecutionStateInput
	result *StateStorageMockGetExecutionStateResult
}

type StateStorageMockGetExecutionStateInput struct {
	p insolar.Reference
}

type StateStorageMockGetExecutionStateResult struct {
	r *ExecutionBroker
}

//Expect specifies that invocation of StateStorage.GetExecutionState is expected from 1 to Infinity times
func (m *mStateStorageMockGetExecutionState) Expect(p insolar.Reference) *mStateStorageMockGetExecutionState {
	m.mock.GetExecutionStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetExecutionStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockGetExecutionStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.GetExecutionState
func (m *mStateStorageMockGetExecutionState) Return(r *ExecutionBroker) *StateStorageMock {
	m.mock.GetExecutionStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetExecutionStateExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockGetExecutionStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.GetExecutionState is expected once
func (m *mStateStorageMockGetExecutionState) ExpectOnce(p insolar.Reference) *StateStorageMockGetExecutionStateExpectation {
	m.mock.GetExecutionStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockGetExecutionStateExpectation{}
	expectation.input = &StateStorageMockGetExecutionStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockGetExecutionStateExpectation) Return(r *ExecutionBroker) {
	e.result = &StateStorageMockGetExecutionStateResult{r}
}

//Set uses given function f as a mock of StateStorage.GetExecutionState method
func (m *mStateStorageMockGetExecutionState) Set(f func(p insolar.Reference) (r *ExecutionBroker)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExecutionStateFunc = f
	return m.mock
}

//GetExecutionState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) GetExecutionState(p insolar.Reference) (r *ExecutionBroker) {
	counter := atomic.AddUint64(&m.GetExecutionStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetExecutionStateCounter, 1)

	if len(m.GetExecutionStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExecutionStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.GetExecutionState. %v", p)
			return
		}

		input := m.GetExecutionStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockGetExecutionStateInput{p}, "StateStorage.GetExecutionState got unexpected parameters")

		result := m.GetExecutionStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetExecutionState")
			return
		}

		r = result.r

		return
	}

	if m.GetExecutionStateMock.mainExpectation != nil {

		input := m.GetExecutionStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockGetExecutionStateInput{p}, "StateStorage.GetExecutionState got unexpected parameters")
		}

		result := m.GetExecutionStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetExecutionState")
		}

		r = result.r

		return
	}

	if m.GetExecutionStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.GetExecutionState. %v", p)
		return
	}

	return m.GetExecutionStateFunc(p)
}

//GetExecutionStateMinimockCounter returns a count of StateStorageMock.GetExecutionStateFunc invocations
func (m *StateStorageMock) GetExecutionStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutionStateCounter)
}

//GetExecutionStateMinimockPreCounter returns the value of StateStorageMock.GetExecutionState invocations
func (m *StateStorageMock) GetExecutionStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutionStatePreCounter)
}

//GetExecutionStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) GetExecutionStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetExecutionStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetExecutionStateCounter) == uint64(len(m.GetExecutionStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetExecutionStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetExecutionStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetExecutionStateFunc != nil {
		return atomic.LoadUint64(&m.GetExecutionStateCounter) > 0
	}

	return true
}

type mStateStorageMockGetValidationState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockGetValidationStateExpectation
	expectationSeries []*StateStorageMockGetValidationStateExpectation
}

type StateStorageMockGetValidationStateExpectation struct {
	input  *StateStorageMockGetValidationStateInput
	result *StateStorageMockGetValidationStateResult
}

type StateStorageMockGetValidationStateInput struct {
	p insolar.Reference
}

type StateStorageMockGetValidationStateResult struct {
	r *ExecutionState
}

//Expect specifies that invocation of StateStorage.GetValidationState is expected from 1 to Infinity times
func (m *mStateStorageMockGetValidationState) Expect(p insolar.Reference) *mStateStorageMockGetValidationState {
	m.mock.GetValidationStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetValidationStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockGetValidationStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.GetValidationState
func (m *mStateStorageMockGetValidationState) Return(r *ExecutionState) *StateStorageMock {
	m.mock.GetValidationStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetValidationStateExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockGetValidationStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.GetValidationState is expected once
func (m *mStateStorageMockGetValidationState) ExpectOnce(p insolar.Reference) *StateStorageMockGetValidationStateExpectation {
	m.mock.GetValidationStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockGetValidationStateExpectation{}
	expectation.input = &StateStorageMockGetValidationStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockGetValidationStateExpectation) Return(r *ExecutionState) {
	e.result = &StateStorageMockGetValidationStateResult{r}
}

//Set uses given function f as a mock of StateStorage.GetValidationState method
func (m *mStateStorageMockGetValidationState) Set(f func(p insolar.Reference) (r *ExecutionState)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetValidationStateFunc = f
	return m.mock
}

//GetValidationState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) GetValidationState(p insolar.Reference) (r *ExecutionState) {
	counter := atomic.AddUint64(&m.GetValidationStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetValidationStateCounter, 1)

	if len(m.GetValidationStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetValidationStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.GetValidationState. %v", p)
			return
		}

		input := m.GetValidationStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockGetValidationStateInput{p}, "StateStorage.GetValidationState got unexpected parameters")

		result := m.GetValidationStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetValidationState")
			return
		}

		r = result.r

		return
	}

	if m.GetValidationStateMock.mainExpectation != nil {

		input := m.GetValidationStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockGetValidationStateInput{p}, "StateStorage.GetValidationState got unexpected parameters")
		}

		result := m.GetValidationStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetValidationState")
		}

		r = result.r

		return
	}

	if m.GetValidationStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.GetValidationState. %v", p)
		return
	}

	return m.GetValidationStateFunc(p)
}

//GetValidationStateMinimockCounter returns a count of StateStorageMock.GetValidationStateFunc invocations
func (m *StateStorageMock) GetValidationStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetValidationStateCounter)
}

//GetValidationStateMinimockPreCounter returns the value of StateStorageMock.GetValidationState invocations
func (m *StateStorageMock) GetValidationStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetValidationStatePreCounter)
}

//GetValidationStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) GetValidationStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetValidationStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetValidationStateCounter) == uint64(len(m.GetValidationStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetValidationStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetValidationStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetValidationStateFunc != nil {
		return atomic.LoadUint64(&m.GetValidationStateCounter) > 0
	}

	return true
}

type mStateStorageMockLock struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockLockExpectation
	expectationSeries []*StateStorageMockLockExpectation
}

type StateStorageMockLockExpectation struct {
}

//Expect specifies that invocation of StateStorage.Lock is expected from 1 to Infinity times
func (m *mStateStorageMockLock) Expect() *mStateStorageMockLock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockLockExpectation{}
	}

	return m
}

//Return specifies results of invocation of StateStorage.Lock
func (m *mStateStorageMockLock) Return() *StateStorageMock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.Lock is expected once
func (m *mStateStorageMockLock) ExpectOnce() *StateStorageMockLockExpectation {
	m.mock.LockFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockLockExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of StateStorage.Lock method
func (m *mStateStorageMockLock) Set(f func()) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LockFunc = f
	return m.mock
}

//Lock implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) Lock() {
	counter := atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if len(m.LockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.Lock.")
			return
		}

		return
	}

	if m.LockMock.mainExpectation != nil {

		return
	}

	if m.LockFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.Lock.")
		return
	}

	m.LockFunc()
}

//LockMinimockCounter returns a count of StateStorageMock.LockFunc invocations
func (m *StateStorageMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of StateStorageMock.Lock invocations
func (m *StateStorageMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

//LockFinished returns true if mock invocations count is ok
func (m *StateStorageMock) LockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LockCounter) == uint64(len(m.LockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LockFunc != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	return true
}

type mStateStorageMockStateMap struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockStateMapExpectation
	expectationSeries []*StateStorageMockStateMapExpectation
}

type StateStorageMockStateMapExpectation struct {
	result *StateStorageMockStateMapResult
}

type StateStorageMockStateMapResult struct {
	r *map[insolar.Reference]*ObjectState
}

//Expect specifies that invocation of StateStorage.StateMap is expected from 1 to Infinity times
func (m *mStateStorageMockStateMap) Expect() *mStateStorageMockStateMap {
	m.mock.StateMapFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockStateMapExpectation{}
	}

	return m
}

//Return specifies results of invocation of StateStorage.StateMap
func (m *mStateStorageMockStateMap) Return(r *map[insolar.Reference]*ObjectState) *StateStorageMock {
	m.mock.StateMapFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockStateMapExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockStateMapResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.StateMap is expected once
func (m *mStateStorageMockStateMap) ExpectOnce() *StateStorageMockStateMapExpectation {
	m.mock.StateMapFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockStateMapExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockStateMapExpectation) Return(r *map[insolar.Reference]*ObjectState) {
	e.result = &StateStorageMockStateMapResult{r}
}

//Set uses given function f as a mock of StateStorage.StateMap method
func (m *mStateStorageMockStateMap) Set(f func() (r *map[insolar.Reference]*ObjectState)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StateMapFunc = f
	return m.mock
}

//StateMap implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) StateMap() (r *map[insolar.Reference]*ObjectState) {
	counter := atomic.AddUint64(&m.StateMapPreCounter, 1)
	defer atomic.AddUint64(&m.StateMapCounter, 1)

	if len(m.StateMapMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StateMapMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.StateMap.")
			return
		}

		result := m.StateMapMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.StateMap")
			return
		}

		r = result.r

		return
	}

	if m.StateMapMock.mainExpectation != nil {

		result := m.StateMapMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.StateMap")
		}

		r = result.r

		return
	}

	if m.StateMapFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.StateMap.")
		return
	}

	return m.StateMapFunc()
}

//StateMapMinimockCounter returns a count of StateStorageMock.StateMapFunc invocations
func (m *StateStorageMock) StateMapMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StateMapCounter)
}

//StateMapMinimockPreCounter returns the value of StateStorageMock.StateMap invocations
func (m *StateStorageMock) StateMapMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StateMapPreCounter)
}

//StateMapFinished returns true if mock invocations count is ok
func (m *StateStorageMock) StateMapFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StateMapMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StateMapCounter) == uint64(len(m.StateMapMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StateMapMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StateMapCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StateMapFunc != nil {
		return atomic.LoadUint64(&m.StateMapCounter) > 0
	}

	return true
}

type mStateStorageMockUnlock struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockUnlockExpectation
	expectationSeries []*StateStorageMockUnlockExpectation
}

type StateStorageMockUnlockExpectation struct {
}

//Expect specifies that invocation of StateStorage.Unlock is expected from 1 to Infinity times
func (m *mStateStorageMockUnlock) Expect() *mStateStorageMockUnlock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUnlockExpectation{}
	}

	return m
}

//Return specifies results of invocation of StateStorage.Unlock
func (m *mStateStorageMockUnlock) Return() *StateStorageMock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUnlockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.Unlock is expected once
func (m *mStateStorageMockUnlock) ExpectOnce() *StateStorageMockUnlockExpectation {
	m.mock.UnlockFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockUnlockExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of StateStorage.Unlock method
func (m *mStateStorageMockUnlock) Set(f func()) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnlockFunc = f
	return m.mock
}

//Unlock implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) Unlock() {
	counter := atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if len(m.UnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.Unlock.")
			return
		}

		return
	}

	if m.UnlockMock.mainExpectation != nil {

		return
	}

	if m.UnlockFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.Unlock.")
		return
	}

	m.UnlockFunc()
}

//UnlockMinimockCounter returns a count of StateStorageMock.UnlockFunc invocations
func (m *StateStorageMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of StateStorageMock.Unlock invocations
func (m *StateStorageMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//UnlockFinished returns true if mock invocations count is ok
func (m *StateStorageMock) UnlockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UnlockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UnlockCounter) == uint64(len(m.UnlockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UnlockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UnlockFunc != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	return true
}

type mStateStorageMockUpsertExecutionState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockUpsertExecutionStateExpectation
	expectationSeries []*StateStorageMockUpsertExecutionStateExpectation
}

type StateStorageMockUpsertExecutionStateExpectation struct {
	input  *StateStorageMockUpsertExecutionStateInput
	result *StateStorageMockUpsertExecutionStateResult
}

type StateStorageMockUpsertExecutionStateInput struct {
	p insolar.Reference
}

type StateStorageMockUpsertExecutionStateResult struct {
	r *ExecutionBroker
}

//Expect specifies that invocation of StateStorage.UpsertExecutionState is expected from 1 to Infinity times
func (m *mStateStorageMockUpsertExecutionState) Expect(p insolar.Reference) *mStateStorageMockUpsertExecutionState {
	m.mock.UpsertExecutionStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUpsertExecutionStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockUpsertExecutionStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.UpsertExecutionState
func (m *mStateStorageMockUpsertExecutionState) Return(r *ExecutionBroker) *StateStorageMock {
	m.mock.UpsertExecutionStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUpsertExecutionStateExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockUpsertExecutionStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.UpsertExecutionState is expected once
func (m *mStateStorageMockUpsertExecutionState) ExpectOnce(p insolar.Reference) *StateStorageMockUpsertExecutionStateExpectation {
	m.mock.UpsertExecutionStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockUpsertExecutionStateExpectation{}
	expectation.input = &StateStorageMockUpsertExecutionStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockUpsertExecutionStateExpectation) Return(r *ExecutionBroker) {
	e.result = &StateStorageMockUpsertExecutionStateResult{r}
}

//Set uses given function f as a mock of StateStorage.UpsertExecutionState method
func (m *mStateStorageMockUpsertExecutionState) Set(f func(p insolar.Reference) (r *ExecutionBroker)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpsertExecutionStateFunc = f
	return m.mock
}

//UpsertExecutionState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) UpsertExecutionState(p insolar.Reference) (r *ExecutionBroker) {
	counter := atomic.AddUint64(&m.UpsertExecutionStatePreCounter, 1)
	defer atomic.AddUint64(&m.UpsertExecutionStateCounter, 1)

	if len(m.UpsertExecutionStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpsertExecutionStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.UpsertExecutionState. %v", p)
			return
		}

		input := m.UpsertExecutionStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockUpsertExecutionStateInput{p}, "StateStorage.UpsertExecutionState got unexpected parameters")

		result := m.UpsertExecutionStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.UpsertExecutionState")
			return
		}

		r = result.r

		return
	}

	if m.UpsertExecutionStateMock.mainExpectation != nil {

		input := m.UpsertExecutionStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockUpsertExecutionStateInput{p}, "StateStorage.UpsertExecutionState got unexpected parameters")
		}

		result := m.UpsertExecutionStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.UpsertExecutionState")
		}

		r = result.r

		return
	}

	if m.UpsertExecutionStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.UpsertExecutionState. %v", p)
		return
	}

	return m.UpsertExecutionStateFunc(p)
}

//UpsertExecutionStateMinimockCounter returns a count of StateStorageMock.UpsertExecutionStateFunc invocations
func (m *StateStorageMock) UpsertExecutionStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertExecutionStateCounter)
}

//UpsertExecutionStateMinimockPreCounter returns the value of StateStorageMock.UpsertExecutionState invocations
func (m *StateStorageMock) UpsertExecutionStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertExecutionStatePreCounter)
}

//UpsertExecutionStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) UpsertExecutionStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpsertExecutionStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpsertExecutionStateCounter) == uint64(len(m.UpsertExecutionStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpsertExecutionStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpsertExecutionStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpsertExecutionStateFunc != nil {
		return atomic.LoadUint64(&m.UpsertExecutionStateCounter) > 0
	}

	return true
}

type mStateStorageMockUpsertValidationState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockUpsertValidationStateExpectation
	expectationSeries []*StateStorageMockUpsertValidationStateExpectation
}

type StateStorageMockUpsertValidationStateExpectation struct {
	input  *StateStorageMockUpsertValidationStateInput
	result *StateStorageMockUpsertValidationStateResult
}

type StateStorageMockUpsertValidationStateInput struct {
	p insolar.Reference
}

type StateStorageMockUpsertValidationStateResult struct {
	r *ExecutionState
}

//Expect specifies that invocation of StateStorage.UpsertValidationState is expected from 1 to Infinity times
func (m *mStateStorageMockUpsertValidationState) Expect(p insolar.Reference) *mStateStorageMockUpsertValidationState {
	m.mock.UpsertValidationStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUpsertValidationStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockUpsertValidationStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.UpsertValidationState
func (m *mStateStorageMockUpsertValidationState) Return(r *ExecutionState) *StateStorageMock {
	m.mock.UpsertValidationStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUpsertValidationStateExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockUpsertValidationStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.UpsertValidationState is expected once
func (m *mStateStorageMockUpsertValidationState) ExpectOnce(p insolar.Reference) *StateStorageMockUpsertValidationStateExpectation {
	m.mock.UpsertValidationStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockUpsertValidationStateExpectation{}
	expectation.input = &StateStorageMockUpsertValidationStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockUpsertValidationStateExpectation) Return(r *ExecutionState) {
	e.result = &StateStorageMockUpsertValidationStateResult{r}
}

//Set uses given function f as a mock of StateStorage.UpsertValidationState method
func (m *mStateStorageMockUpsertValidationState) Set(f func(p insolar.Reference) (r *ExecutionState)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpsertValidationStateFunc = f
	return m.mock
}

//UpsertValidationState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) UpsertValidationState(p insolar.Reference) (r *ExecutionState) {
	counter := atomic.AddUint64(&m.UpsertValidationStatePreCounter, 1)
	defer atomic.AddUint64(&m.UpsertValidationStateCounter, 1)

	if len(m.UpsertValidationStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpsertValidationStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.UpsertValidationState. %v", p)
			return
		}

		input := m.UpsertValidationStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockUpsertValidationStateInput{p}, "StateStorage.UpsertValidationState got unexpected parameters")

		result := m.UpsertValidationStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.UpsertValidationState")
			return
		}

		r = result.r

		return
	}

	if m.UpsertValidationStateMock.mainExpectation != nil {

		input := m.UpsertValidationStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockUpsertValidationStateInput{p}, "StateStorage.UpsertValidationState got unexpected parameters")
		}

		result := m.UpsertValidationStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.UpsertValidationState")
		}

		r = result.r

		return
	}

	if m.UpsertValidationStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.UpsertValidationState. %v", p)
		return
	}

	return m.UpsertValidationStateFunc(p)
}

//UpsertValidationStateMinimockCounter returns a count of StateStorageMock.UpsertValidationStateFunc invocations
func (m *StateStorageMock) UpsertValidationStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertValidationStateCounter)
}

//UpsertValidationStateMinimockPreCounter returns the value of StateStorageMock.UpsertValidationState invocations
func (m *StateStorageMock) UpsertValidationStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertValidationStatePreCounter)
}

//UpsertValidationStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) UpsertValidationStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpsertValidationStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpsertValidationStateCounter) == uint64(len(m.UpsertValidationStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpsertValidationStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpsertValidationStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpsertValidationStateFunc != nil {
		return atomic.LoadUint64(&m.UpsertValidationStateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateStorageMock) ValidateCallCounters() {

	if !m.DeleteObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.DeleteObjectState")
	}

	if !m.GetExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetExecutionState")
	}

	if !m.GetValidationStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetValidationState")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Lock")
	}

	if !m.StateMapFinished() {
		m.t.Fatal("Expected call to StateStorageMock.StateMap")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Unlock")
	}

	if !m.UpsertExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertExecutionState")
	}

	if !m.UpsertValidationStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertValidationState")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StateStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StateStorageMock) MinimockFinish() {

	if !m.DeleteObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.DeleteObjectState")
	}

	if !m.GetExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetExecutionState")
	}

	if !m.GetValidationStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetValidationState")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Lock")
	}

	if !m.StateMapFinished() {
		m.t.Fatal("Expected call to StateStorageMock.StateMap")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Unlock")
	}

	if !m.UpsertExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertExecutionState")
	}

	if !m.UpsertValidationStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertValidationState")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StateStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StateStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.DeleteObjectStateFinished()
		ok = ok && m.GetExecutionStateFinished()
		ok = ok && m.GetValidationStateFinished()
		ok = ok && m.LockFinished()
		ok = ok && m.StateMapFinished()
		ok = ok && m.UnlockFinished()
		ok = ok && m.UpsertExecutionStateFinished()
		ok = ok && m.UpsertValidationStateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeleteObjectStateFinished() {
				m.t.Error("Expected call to StateStorageMock.DeleteObjectState")
			}

			if !m.GetExecutionStateFinished() {
				m.t.Error("Expected call to StateStorageMock.GetExecutionState")
			}

			if !m.GetValidationStateFinished() {
				m.t.Error("Expected call to StateStorageMock.GetValidationState")
			}

			if !m.LockFinished() {
				m.t.Error("Expected call to StateStorageMock.Lock")
			}

			if !m.StateMapFinished() {
				m.t.Error("Expected call to StateStorageMock.StateMap")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to StateStorageMock.Unlock")
			}

			if !m.UpsertExecutionStateFinished() {
				m.t.Error("Expected call to StateStorageMock.UpsertExecutionState")
			}

			if !m.UpsertValidationStateFinished() {
				m.t.Error("Expected call to StateStorageMock.UpsertValidationState")
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
func (m *StateStorageMock) AllMocksCalled() bool {

	if !m.DeleteObjectStateFinished() {
		return false
	}

	if !m.GetExecutionStateFinished() {
		return false
	}

	if !m.GetValidationStateFinished() {
		return false
	}

	if !m.LockFinished() {
		return false
	}

	if !m.StateMapFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	if !m.UpsertExecutionStateFinished() {
		return false
	}

	if !m.UpsertValidationStateFinished() {
		return false
	}

	return true
}

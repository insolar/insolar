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

	GetExecutionStateFunc       func(p insolar.Reference) (r *ExecutionState)
	GetExecutionStateCounter    uint64
	GetExecutionStatePreCounter uint64
	GetExecutionStateMock       mStateStorageMockGetExecutionState

	GetObjectStateFunc       func(p insolar.Reference) (r *ObjectState)
	GetObjectStateCounter    uint64
	GetObjectStatePreCounter uint64
	GetObjectStateMock       mStateStorageMockGetObjectState

	LockFunc       func()
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mStateStorageMockLock

	MustObjectStateFunc       func(p insolar.Reference) (r *ObjectState)
	MustObjectStateCounter    uint64
	MustObjectStatePreCounter uint64
	MustObjectStateMock       mStateStorageMockMustObjectState

	StateMapFunc       func() (r *map[insolar.Reference]*ObjectState)
	StateMapCounter    uint64
	StateMapPreCounter uint64
	StateMapMock       mStateStorageMockStateMap

	UnlockFunc       func()
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mStateStorageMockUnlock

	UpsertObjectStateFunc       func(p insolar.Reference) (r *ObjectState)
	UpsertObjectStateCounter    uint64
	UpsertObjectStatePreCounter uint64
	UpsertObjectStateMock       mStateStorageMockUpsertObjectState
}

//NewStateStorageMock returns a mock for github.com/insolar/insolar/logicrunner.StateStorage
func NewStateStorageMock(t minimock.Tester) *StateStorageMock {
	m := &StateStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteObjectStateMock = mStateStorageMockDeleteObjectState{mock: m}
	m.GetExecutionStateMock = mStateStorageMockGetExecutionState{mock: m}
	m.GetObjectStateMock = mStateStorageMockGetObjectState{mock: m}
	m.LockMock = mStateStorageMockLock{mock: m}
	m.MustObjectStateMock = mStateStorageMockMustObjectState{mock: m}
	m.StateMapMock = mStateStorageMockStateMap{mock: m}
	m.UnlockMock = mStateStorageMockUnlock{mock: m}
	m.UpsertObjectStateMock = mStateStorageMockUpsertObjectState{mock: m}

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
	r *ExecutionState
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
func (m *mStateStorageMockGetExecutionState) Return(r *ExecutionState) *StateStorageMock {
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

func (e *StateStorageMockGetExecutionStateExpectation) Return(r *ExecutionState) {
	e.result = &StateStorageMockGetExecutionStateResult{r}
}

//Set uses given function f as a mock of StateStorage.GetExecutionState method
func (m *mStateStorageMockGetExecutionState) Set(f func(p insolar.Reference) (r *ExecutionState)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExecutionStateFunc = f
	return m.mock
}

//GetExecutionState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) GetExecutionState(p insolar.Reference) (r *ExecutionState) {
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

type mStateStorageMockGetObjectState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockGetObjectStateExpectation
	expectationSeries []*StateStorageMockGetObjectStateExpectation
}

type StateStorageMockGetObjectStateExpectation struct {
	input  *StateStorageMockGetObjectStateInput
	result *StateStorageMockGetObjectStateResult
}

type StateStorageMockGetObjectStateInput struct {
	p insolar.Reference
}

type StateStorageMockGetObjectStateResult struct {
	r *ObjectState
}

//Expect specifies that invocation of StateStorage.GetObjectState is expected from 1 to Infinity times
func (m *mStateStorageMockGetObjectState) Expect(p insolar.Reference) *mStateStorageMockGetObjectState {
	m.mock.GetObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetObjectStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockGetObjectStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.GetObjectState
func (m *mStateStorageMockGetObjectState) Return(r *ObjectState) *StateStorageMock {
	m.mock.GetObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetObjectStateExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockGetObjectStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.GetObjectState is expected once
func (m *mStateStorageMockGetObjectState) ExpectOnce(p insolar.Reference) *StateStorageMockGetObjectStateExpectation {
	m.mock.GetObjectStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockGetObjectStateExpectation{}
	expectation.input = &StateStorageMockGetObjectStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockGetObjectStateExpectation) Return(r *ObjectState) {
	e.result = &StateStorageMockGetObjectStateResult{r}
}

//Set uses given function f as a mock of StateStorage.GetObjectState method
func (m *mStateStorageMockGetObjectState) Set(f func(p insolar.Reference) (r *ObjectState)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectStateFunc = f
	return m.mock
}

//GetObjectState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) GetObjectState(p insolar.Reference) (r *ObjectState) {
	counter := atomic.AddUint64(&m.GetObjectStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectStateCounter, 1)

	if len(m.GetObjectStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.GetObjectState. %v", p)
			return
		}

		input := m.GetObjectStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockGetObjectStateInput{p}, "StateStorage.GetObjectState got unexpected parameters")

		result := m.GetObjectStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetObjectState")
			return
		}

		r = result.r

		return
	}

	if m.GetObjectStateMock.mainExpectation != nil {

		input := m.GetObjectStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockGetObjectStateInput{p}, "StateStorage.GetObjectState got unexpected parameters")
		}

		result := m.GetObjectStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetObjectState")
		}

		r = result.r

		return
	}

	if m.GetObjectStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.GetObjectState. %v", p)
		return
	}

	return m.GetObjectStateFunc(p)
}

//GetObjectStateMinimockCounter returns a count of StateStorageMock.GetObjectStateFunc invocations
func (m *StateStorageMock) GetObjectStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectStateCounter)
}

//GetObjectStateMinimockPreCounter returns the value of StateStorageMock.GetObjectState invocations
func (m *StateStorageMock) GetObjectStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectStatePreCounter)
}

//GetObjectStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) GetObjectStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetObjectStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetObjectStateCounter) == uint64(len(m.GetObjectStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetObjectStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetObjectStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetObjectStateFunc != nil {
		return atomic.LoadUint64(&m.GetObjectStateCounter) > 0
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

type mStateStorageMockMustObjectState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockMustObjectStateExpectation
	expectationSeries []*StateStorageMockMustObjectStateExpectation
}

type StateStorageMockMustObjectStateExpectation struct {
	input  *StateStorageMockMustObjectStateInput
	result *StateStorageMockMustObjectStateResult
}

type StateStorageMockMustObjectStateInput struct {
	p insolar.Reference
}

type StateStorageMockMustObjectStateResult struct {
	r *ObjectState
}

//Expect specifies that invocation of StateStorage.MustObjectState is expected from 1 to Infinity times
func (m *mStateStorageMockMustObjectState) Expect(p insolar.Reference) *mStateStorageMockMustObjectState {
	m.mock.MustObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockMustObjectStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockMustObjectStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.MustObjectState
func (m *mStateStorageMockMustObjectState) Return(r *ObjectState) *StateStorageMock {
	m.mock.MustObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockMustObjectStateExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockMustObjectStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.MustObjectState is expected once
func (m *mStateStorageMockMustObjectState) ExpectOnce(p insolar.Reference) *StateStorageMockMustObjectStateExpectation {
	m.mock.MustObjectStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockMustObjectStateExpectation{}
	expectation.input = &StateStorageMockMustObjectStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockMustObjectStateExpectation) Return(r *ObjectState) {
	e.result = &StateStorageMockMustObjectStateResult{r}
}

//Set uses given function f as a mock of StateStorage.MustObjectState method
func (m *mStateStorageMockMustObjectState) Set(f func(p insolar.Reference) (r *ObjectState)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MustObjectStateFunc = f
	return m.mock
}

//MustObjectState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) MustObjectState(p insolar.Reference) (r *ObjectState) {
	counter := atomic.AddUint64(&m.MustObjectStatePreCounter, 1)
	defer atomic.AddUint64(&m.MustObjectStateCounter, 1)

	if len(m.MustObjectStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MustObjectStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.MustObjectState. %v", p)
			return
		}

		input := m.MustObjectStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockMustObjectStateInput{p}, "StateStorage.MustObjectState got unexpected parameters")

		result := m.MustObjectStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.MustObjectState")
			return
		}

		r = result.r

		return
	}

	if m.MustObjectStateMock.mainExpectation != nil {

		input := m.MustObjectStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockMustObjectStateInput{p}, "StateStorage.MustObjectState got unexpected parameters")
		}

		result := m.MustObjectStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.MustObjectState")
		}

		r = result.r

		return
	}

	if m.MustObjectStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.MustObjectState. %v", p)
		return
	}

	return m.MustObjectStateFunc(p)
}

//MustObjectStateMinimockCounter returns a count of StateStorageMock.MustObjectStateFunc invocations
func (m *StateStorageMock) MustObjectStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustObjectStateCounter)
}

//MustObjectStateMinimockPreCounter returns the value of StateStorageMock.MustObjectState invocations
func (m *StateStorageMock) MustObjectStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustObjectStatePreCounter)
}

//MustObjectStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) MustObjectStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MustObjectStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MustObjectStateCounter) == uint64(len(m.MustObjectStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MustObjectStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MustObjectStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MustObjectStateFunc != nil {
		return atomic.LoadUint64(&m.MustObjectStateCounter) > 0
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

type mStateStorageMockUpsertObjectState struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockUpsertObjectStateExpectation
	expectationSeries []*StateStorageMockUpsertObjectStateExpectation
}

type StateStorageMockUpsertObjectStateExpectation struct {
	input  *StateStorageMockUpsertObjectStateInput
	result *StateStorageMockUpsertObjectStateResult
}

type StateStorageMockUpsertObjectStateInput struct {
	p insolar.Reference
}

type StateStorageMockUpsertObjectStateResult struct {
	r *ObjectState
}

//Expect specifies that invocation of StateStorage.UpsertObjectState is expected from 1 to Infinity times
func (m *mStateStorageMockUpsertObjectState) Expect(p insolar.Reference) *mStateStorageMockUpsertObjectState {
	m.mock.UpsertObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUpsertObjectStateExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockUpsertObjectStateInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.UpsertObjectState
func (m *mStateStorageMockUpsertObjectState) Return(r *ObjectState) *StateStorageMock {
	m.mock.UpsertObjectStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockUpsertObjectStateExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockUpsertObjectStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.UpsertObjectState is expected once
func (m *mStateStorageMockUpsertObjectState) ExpectOnce(p insolar.Reference) *StateStorageMockUpsertObjectStateExpectation {
	m.mock.UpsertObjectStateFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockUpsertObjectStateExpectation{}
	expectation.input = &StateStorageMockUpsertObjectStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockUpsertObjectStateExpectation) Return(r *ObjectState) {
	e.result = &StateStorageMockUpsertObjectStateResult{r}
}

//Set uses given function f as a mock of StateStorage.UpsertObjectState method
func (m *mStateStorageMockUpsertObjectState) Set(f func(p insolar.Reference) (r *ObjectState)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpsertObjectStateFunc = f
	return m.mock
}

//UpsertObjectState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) UpsertObjectState(p insolar.Reference) (r *ObjectState) {
	counter := atomic.AddUint64(&m.UpsertObjectStatePreCounter, 1)
	defer atomic.AddUint64(&m.UpsertObjectStateCounter, 1)

	if len(m.UpsertObjectStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpsertObjectStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.UpsertObjectState. %v", p)
			return
		}

		input := m.UpsertObjectStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockUpsertObjectStateInput{p}, "StateStorage.UpsertObjectState got unexpected parameters")

		result := m.UpsertObjectStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.UpsertObjectState")
			return
		}

		r = result.r

		return
	}

	if m.UpsertObjectStateMock.mainExpectation != nil {

		input := m.UpsertObjectStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockUpsertObjectStateInput{p}, "StateStorage.UpsertObjectState got unexpected parameters")
		}

		result := m.UpsertObjectStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.UpsertObjectState")
		}

		r = result.r

		return
	}

	if m.UpsertObjectStateFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.UpsertObjectState. %v", p)
		return
	}

	return m.UpsertObjectStateFunc(p)
}

//UpsertObjectStateMinimockCounter returns a count of StateStorageMock.UpsertObjectStateFunc invocations
func (m *StateStorageMock) UpsertObjectStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertObjectStateCounter)
}

//UpsertObjectStateMinimockPreCounter returns the value of StateStorageMock.UpsertObjectState invocations
func (m *StateStorageMock) UpsertObjectStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertObjectStatePreCounter)
}

//UpsertObjectStateFinished returns true if mock invocations count is ok
func (m *StateStorageMock) UpsertObjectStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpsertObjectStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpsertObjectStateCounter) == uint64(len(m.UpsertObjectStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpsertObjectStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpsertObjectStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpsertObjectStateFunc != nil {
		return atomic.LoadUint64(&m.UpsertObjectStateCounter) > 0
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

	if !m.GetObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetObjectState")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Lock")
	}

	if !m.MustObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.MustObjectState")
	}

	if !m.StateMapFinished() {
		m.t.Fatal("Expected call to StateStorageMock.StateMap")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Unlock")
	}

	if !m.UpsertObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertObjectState")
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

	if !m.GetObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetObjectState")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Lock")
	}

	if !m.MustObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.MustObjectState")
	}

	if !m.StateMapFinished() {
		m.t.Fatal("Expected call to StateStorageMock.StateMap")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Unlock")
	}

	if !m.UpsertObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertObjectState")
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
		ok = ok && m.GetObjectStateFinished()
		ok = ok && m.LockFinished()
		ok = ok && m.MustObjectStateFinished()
		ok = ok && m.StateMapFinished()
		ok = ok && m.UnlockFinished()
		ok = ok && m.UpsertObjectStateFinished()

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

			if !m.GetObjectStateFinished() {
				m.t.Error("Expected call to StateStorageMock.GetObjectState")
			}

			if !m.LockFinished() {
				m.t.Error("Expected call to StateStorageMock.Lock")
			}

			if !m.MustObjectStateFinished() {
				m.t.Error("Expected call to StateStorageMock.MustObjectState")
			}

			if !m.StateMapFinished() {
				m.t.Error("Expected call to StateStorageMock.StateMap")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to StateStorageMock.Unlock")
			}

			if !m.UpsertObjectStateFinished() {
				m.t.Error("Expected call to StateStorageMock.UpsertObjectState")
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

	if !m.GetObjectStateFinished() {
		return false
	}

	if !m.LockFinished() {
		return false
	}

	if !m.MustObjectStateFinished() {
		return false
	}

	if !m.StateMapFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	if !m.UpsertObjectStateFinished() {
		return false
	}

	return true
}

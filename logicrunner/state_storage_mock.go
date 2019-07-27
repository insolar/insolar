package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "StateStorage" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
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

	GetExecutionArchiveFunc       func(p insolar.Reference) (r ExecutionArchive)
	GetExecutionArchiveCounter    uint64
	GetExecutionArchivePreCounter uint64
	GetExecutionArchiveMock       mStateStorageMockGetExecutionArchive

	GetExecutionStateFunc       func(p insolar.Reference) (r ExecutionBrokerI)
	GetExecutionStateCounter    uint64
	GetExecutionStatePreCounter uint64
	GetExecutionStateMock       mStateStorageMockGetExecutionState

	IsEmptyFunc       func() (r bool)
	IsEmptyCounter    uint64
	IsEmptyPreCounter uint64
	IsEmptyMock       mStateStorageMockIsEmpty

	LockFunc       func()
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mStateStorageMockLock

	OnPulseFunc       func(p context.Context, p1 insolar.Pulse) (r []insolar.Message)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mStateStorageMockOnPulse

	UnlockFunc       func()
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mStateStorageMockUnlock

	UpsertExecutionStateFunc       func(p insolar.Reference) (r ExecutionBrokerI)
	UpsertExecutionStateCounter    uint64
	UpsertExecutionStatePreCounter uint64
	UpsertExecutionStateMock       mStateStorageMockUpsertExecutionState
}

//NewStateStorageMock returns a mock for github.com/insolar/insolar/logicrunner.StateStorage
func NewStateStorageMock(t minimock.Tester) *StateStorageMock {
	m := &StateStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteObjectStateMock = mStateStorageMockDeleteObjectState{mock: m}
	m.GetExecutionArchiveMock = mStateStorageMockGetExecutionArchive{mock: m}
	m.GetExecutionStateMock = mStateStorageMockGetExecutionState{mock: m}
	m.IsEmptyMock = mStateStorageMockIsEmpty{mock: m}
	m.LockMock = mStateStorageMockLock{mock: m}
	m.OnPulseMock = mStateStorageMockOnPulse{mock: m}
	m.UnlockMock = mStateStorageMockUnlock{mock: m}
	m.UpsertExecutionStateMock = mStateStorageMockUpsertExecutionState{mock: m}

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

type mStateStorageMockGetExecutionArchive struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockGetExecutionArchiveExpectation
	expectationSeries []*StateStorageMockGetExecutionArchiveExpectation
}

type StateStorageMockGetExecutionArchiveExpectation struct {
	input  *StateStorageMockGetExecutionArchiveInput
	result *StateStorageMockGetExecutionArchiveResult
}

type StateStorageMockGetExecutionArchiveInput struct {
	p insolar.Reference
}

type StateStorageMockGetExecutionArchiveResult struct {
	r ExecutionArchive
}

//Expect specifies that invocation of StateStorage.GetExecutionArchive is expected from 1 to Infinity times
func (m *mStateStorageMockGetExecutionArchive) Expect(p insolar.Reference) *mStateStorageMockGetExecutionArchive {
	m.mock.GetExecutionArchiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetExecutionArchiveExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockGetExecutionArchiveInput{p}
	return m
}

//Return specifies results of invocation of StateStorage.GetExecutionArchive
func (m *mStateStorageMockGetExecutionArchive) Return(r ExecutionArchive) *StateStorageMock {
	m.mock.GetExecutionArchiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockGetExecutionArchiveExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockGetExecutionArchiveResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.GetExecutionArchive is expected once
func (m *mStateStorageMockGetExecutionArchive) ExpectOnce(p insolar.Reference) *StateStorageMockGetExecutionArchiveExpectation {
	m.mock.GetExecutionArchiveFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockGetExecutionArchiveExpectation{}
	expectation.input = &StateStorageMockGetExecutionArchiveInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockGetExecutionArchiveExpectation) Return(r ExecutionArchive) {
	e.result = &StateStorageMockGetExecutionArchiveResult{r}
}

//Set uses given function f as a mock of StateStorage.GetExecutionArchive method
func (m *mStateStorageMockGetExecutionArchive) Set(f func(p insolar.Reference) (r ExecutionArchive)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExecutionArchiveFunc = f
	return m.mock
}

//GetExecutionArchive implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) GetExecutionArchive(p insolar.Reference) (r ExecutionArchive) {
	counter := atomic.AddUint64(&m.GetExecutionArchivePreCounter, 1)
	defer atomic.AddUint64(&m.GetExecutionArchiveCounter, 1)

	if len(m.GetExecutionArchiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetExecutionArchiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.GetExecutionArchive. %v", p)
			return
		}

		input := m.GetExecutionArchiveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockGetExecutionArchiveInput{p}, "StateStorage.GetExecutionArchive got unexpected parameters")

		result := m.GetExecutionArchiveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetExecutionArchive")
			return
		}

		r = result.r

		return
	}

	if m.GetExecutionArchiveMock.mainExpectation != nil {

		input := m.GetExecutionArchiveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockGetExecutionArchiveInput{p}, "StateStorage.GetExecutionArchive got unexpected parameters")
		}

		result := m.GetExecutionArchiveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.GetExecutionArchive")
		}

		r = result.r

		return
	}

	if m.GetExecutionArchiveFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.GetExecutionArchive. %v", p)
		return
	}

	return m.GetExecutionArchiveFunc(p)
}

//GetExecutionArchiveMinimockCounter returns a count of StateStorageMock.GetExecutionArchiveFunc invocations
func (m *StateStorageMock) GetExecutionArchiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutionArchiveCounter)
}

//GetExecutionArchiveMinimockPreCounter returns the value of StateStorageMock.GetExecutionArchive invocations
func (m *StateStorageMock) GetExecutionArchiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetExecutionArchivePreCounter)
}

//GetExecutionArchiveFinished returns true if mock invocations count is ok
func (m *StateStorageMock) GetExecutionArchiveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetExecutionArchiveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetExecutionArchiveCounter) == uint64(len(m.GetExecutionArchiveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetExecutionArchiveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetExecutionArchiveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetExecutionArchiveFunc != nil {
		return atomic.LoadUint64(&m.GetExecutionArchiveCounter) > 0
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
	r ExecutionBrokerI
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
func (m *mStateStorageMockGetExecutionState) Return(r ExecutionBrokerI) *StateStorageMock {
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

func (e *StateStorageMockGetExecutionStateExpectation) Return(r ExecutionBrokerI) {
	e.result = &StateStorageMockGetExecutionStateResult{r}
}

//Set uses given function f as a mock of StateStorage.GetExecutionState method
func (m *mStateStorageMockGetExecutionState) Set(f func(p insolar.Reference) (r ExecutionBrokerI)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetExecutionStateFunc = f
	return m.mock
}

//GetExecutionState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) GetExecutionState(p insolar.Reference) (r ExecutionBrokerI) {
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

type mStateStorageMockIsEmpty struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockIsEmptyExpectation
	expectationSeries []*StateStorageMockIsEmptyExpectation
}

type StateStorageMockIsEmptyExpectation struct {
	result *StateStorageMockIsEmptyResult
}

type StateStorageMockIsEmptyResult struct {
	r bool
}

//Expect specifies that invocation of StateStorage.IsEmpty is expected from 1 to Infinity times
func (m *mStateStorageMockIsEmpty) Expect() *mStateStorageMockIsEmpty {
	m.mock.IsEmptyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockIsEmptyExpectation{}
	}

	return m
}

//Return specifies results of invocation of StateStorage.IsEmpty
func (m *mStateStorageMockIsEmpty) Return(r bool) *StateStorageMock {
	m.mock.IsEmptyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockIsEmptyExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockIsEmptyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.IsEmpty is expected once
func (m *mStateStorageMockIsEmpty) ExpectOnce() *StateStorageMockIsEmptyExpectation {
	m.mock.IsEmptyFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockIsEmptyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockIsEmptyExpectation) Return(r bool) {
	e.result = &StateStorageMockIsEmptyResult{r}
}

//Set uses given function f as a mock of StateStorage.IsEmpty method
func (m *mStateStorageMockIsEmpty) Set(f func() (r bool)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsEmptyFunc = f
	return m.mock
}

//IsEmpty implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) IsEmpty() (r bool) {
	counter := atomic.AddUint64(&m.IsEmptyPreCounter, 1)
	defer atomic.AddUint64(&m.IsEmptyCounter, 1)

	if len(m.IsEmptyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsEmptyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.IsEmpty.")
			return
		}

		result := m.IsEmptyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.IsEmpty")
			return
		}

		r = result.r

		return
	}

	if m.IsEmptyMock.mainExpectation != nil {

		result := m.IsEmptyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.IsEmpty")
		}

		r = result.r

		return
	}

	if m.IsEmptyFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.IsEmpty.")
		return
	}

	return m.IsEmptyFunc()
}

//IsEmptyMinimockCounter returns a count of StateStorageMock.IsEmptyFunc invocations
func (m *StateStorageMock) IsEmptyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsEmptyCounter)
}

//IsEmptyMinimockPreCounter returns the value of StateStorageMock.IsEmpty invocations
func (m *StateStorageMock) IsEmptyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsEmptyPreCounter)
}

//IsEmptyFinished returns true if mock invocations count is ok
func (m *StateStorageMock) IsEmptyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsEmptyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsEmptyCounter) == uint64(len(m.IsEmptyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsEmptyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsEmptyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsEmptyFunc != nil {
		return atomic.LoadUint64(&m.IsEmptyCounter) > 0
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

type mStateStorageMockOnPulse struct {
	mock              *StateStorageMock
	mainExpectation   *StateStorageMockOnPulseExpectation
	expectationSeries []*StateStorageMockOnPulseExpectation
}

type StateStorageMockOnPulseExpectation struct {
	input  *StateStorageMockOnPulseInput
	result *StateStorageMockOnPulseResult
}

type StateStorageMockOnPulseInput struct {
	p  context.Context
	p1 insolar.Pulse
}

type StateStorageMockOnPulseResult struct {
	r []insolar.Message
}

//Expect specifies that invocation of StateStorage.OnPulse is expected from 1 to Infinity times
func (m *mStateStorageMockOnPulse) Expect(p context.Context, p1 insolar.Pulse) *mStateStorageMockOnPulse {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockOnPulseExpectation{}
	}
	m.mainExpectation.input = &StateStorageMockOnPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of StateStorage.OnPulse
func (m *mStateStorageMockOnPulse) Return(r []insolar.Message) *StateStorageMock {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StateStorageMockOnPulseExpectation{}
	}
	m.mainExpectation.result = &StateStorageMockOnPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of StateStorage.OnPulse is expected once
func (m *mStateStorageMockOnPulse) ExpectOnce(p context.Context, p1 insolar.Pulse) *StateStorageMockOnPulseExpectation {
	m.mock.OnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &StateStorageMockOnPulseExpectation{}
	expectation.input = &StateStorageMockOnPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StateStorageMockOnPulseExpectation) Return(r []insolar.Message) {
	e.result = &StateStorageMockOnPulseResult{r}
}

//Set uses given function f as a mock of StateStorage.OnPulse method
func (m *mStateStorageMockOnPulse) Set(f func(p context.Context, p1 insolar.Pulse) (r []insolar.Message)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) OnPulse(p context.Context, p1 insolar.Pulse) (r []insolar.Message) {
	counter := atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if len(m.OnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StateStorageMock.OnPulse. %v %v", p, p1)
			return
		}

		input := m.OnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StateStorageMockOnPulseInput{p, p1}, "StateStorage.OnPulse got unexpected parameters")

		result := m.OnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.OnPulse")
			return
		}

		r = result.r

		return
	}

	if m.OnPulseMock.mainExpectation != nil {

		input := m.OnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StateStorageMockOnPulseInput{p, p1}, "StateStorage.OnPulse got unexpected parameters")
		}

		result := m.OnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StateStorageMock.OnPulse")
		}

		r = result.r

		return
	}

	if m.OnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to StateStorageMock.OnPulse. %v %v", p, p1)
		return
	}

	return m.OnPulseFunc(p, p1)
}

//OnPulseMinimockCounter returns a count of StateStorageMock.OnPulseFunc invocations
func (m *StateStorageMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of StateStorageMock.OnPulse invocations
func (m *StateStorageMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//OnPulseFinished returns true if mock invocations count is ok
func (m *StateStorageMock) OnPulseFinished() bool {
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
	r ExecutionBrokerI
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
func (m *mStateStorageMockUpsertExecutionState) Return(r ExecutionBrokerI) *StateStorageMock {
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

func (e *StateStorageMockUpsertExecutionStateExpectation) Return(r ExecutionBrokerI) {
	e.result = &StateStorageMockUpsertExecutionStateResult{r}
}

//Set uses given function f as a mock of StateStorage.UpsertExecutionState method
func (m *mStateStorageMockUpsertExecutionState) Set(f func(p insolar.Reference) (r ExecutionBrokerI)) *StateStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpsertExecutionStateFunc = f
	return m.mock
}

//UpsertExecutionState implements github.com/insolar/insolar/logicrunner.StateStorage interface
func (m *StateStorageMock) UpsertExecutionState(p insolar.Reference) (r ExecutionBrokerI) {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StateStorageMock) ValidateCallCounters() {

	if !m.DeleteObjectStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.DeleteObjectState")
	}

	if !m.GetExecutionArchiveFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetExecutionArchive")
	}

	if !m.GetExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetExecutionState")
	}

	if !m.IsEmptyFinished() {
		m.t.Fatal("Expected call to StateStorageMock.IsEmpty")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Lock")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to StateStorageMock.OnPulse")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Unlock")
	}

	if !m.UpsertExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertExecutionState")
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

	if !m.GetExecutionArchiveFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetExecutionArchive")
	}

	if !m.GetExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.GetExecutionState")
	}

	if !m.IsEmptyFinished() {
		m.t.Fatal("Expected call to StateStorageMock.IsEmpty")
	}

	if !m.LockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Lock")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to StateStorageMock.OnPulse")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to StateStorageMock.Unlock")
	}

	if !m.UpsertExecutionStateFinished() {
		m.t.Fatal("Expected call to StateStorageMock.UpsertExecutionState")
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
		ok = ok && m.GetExecutionArchiveFinished()
		ok = ok && m.GetExecutionStateFinished()
		ok = ok && m.IsEmptyFinished()
		ok = ok && m.LockFinished()
		ok = ok && m.OnPulseFinished()
		ok = ok && m.UnlockFinished()
		ok = ok && m.UpsertExecutionStateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeleteObjectStateFinished() {
				m.t.Error("Expected call to StateStorageMock.DeleteObjectState")
			}

			if !m.GetExecutionArchiveFinished() {
				m.t.Error("Expected call to StateStorageMock.GetExecutionArchive")
			}

			if !m.GetExecutionStateFinished() {
				m.t.Error("Expected call to StateStorageMock.GetExecutionState")
			}

			if !m.IsEmptyFinished() {
				m.t.Error("Expected call to StateStorageMock.IsEmpty")
			}

			if !m.LockFinished() {
				m.t.Error("Expected call to StateStorageMock.Lock")
			}

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to StateStorageMock.OnPulse")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to StateStorageMock.Unlock")
			}

			if !m.UpsertExecutionStateFinished() {
				m.t.Error("Expected call to StateStorageMock.UpsertExecutionState")
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

	if !m.GetExecutionArchiveFinished() {
		return false
	}

	if !m.GetExecutionStateFinished() {
		return false
	}

	if !m.IsEmptyFinished() {
		return false
	}

	if !m.LockFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	if !m.UpsertExecutionStateFinished() {
		return false
	}

	return true
}

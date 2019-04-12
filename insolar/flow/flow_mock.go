package flow

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Flow" can be found in github.com/insolar/insolar/insolar/flow
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//FlowMock implements github.com/insolar/insolar/insolar/flow.Flow
type FlowMock struct {
	t minimock.Tester

	ContinueFunc       func(p context.Context)
	ContinueCounter    uint64
	ContinuePreCounter uint64
	ContinueMock       mFlowMockContinue

	HandleFunc       func(p context.Context, p1 Handle) (r error)
	HandleCounter    uint64
	HandlePreCounter uint64
	HandleMock       mFlowMockHandle

	MigrateFunc       func(p context.Context, p1 Handle) (r error)
	MigrateCounter    uint64
	MigratePreCounter uint64
	MigrateMock       mFlowMockMigrate

	ProcedureFunc       func(p context.Context, p1 Procedure) (r error)
	ProcedureCounter    uint64
	ProcedurePreCounter uint64
	ProcedureMock       mFlowMockProcedure
}

//NewFlowMock returns a mock for github.com/insolar/insolar/insolar/flow.Flow
func NewFlowMock(t minimock.Tester) *FlowMock {
	m := &FlowMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ContinueMock = mFlowMockContinue{mock: m}
	m.HandleMock = mFlowMockHandle{mock: m}
	m.MigrateMock = mFlowMockMigrate{mock: m}
	m.ProcedureMock = mFlowMockProcedure{mock: m}

	return m
}

type mFlowMockContinue struct {
	mock              *FlowMock
	mainExpectation   *FlowMockContinueExpectation
	expectationSeries []*FlowMockContinueExpectation
}

type FlowMockContinueExpectation struct {
	input *FlowMockContinueInput
}

type FlowMockContinueInput struct {
	p context.Context
}

//Expect specifies that invocation of Flow.Continue is expected from 1 to Infinity times
func (m *mFlowMockContinue) Expect(p context.Context) *mFlowMockContinue {
	m.mock.ContinueFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockContinueExpectation{}
	}
	m.mainExpectation.input = &FlowMockContinueInput{p}
	return m
}

//Return specifies results of invocation of Flow.Continue
func (m *mFlowMockContinue) Return() *FlowMock {
	m.mock.ContinueFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockContinueExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Flow.Continue is expected once
func (m *mFlowMockContinue) ExpectOnce(p context.Context) *FlowMockContinueExpectation {
	m.mock.ContinueFunc = nil
	m.mainExpectation = nil

	expectation := &FlowMockContinueExpectation{}
	expectation.input = &FlowMockContinueInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Flow.Continue method
func (m *mFlowMockContinue) Set(f func(p context.Context)) *FlowMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ContinueFunc = f
	return m.mock
}

//Continue implements github.com/insolar/insolar/insolar/flow.Flow interface
func (m *FlowMock) Continue(p context.Context) {
	counter := atomic.AddUint64(&m.ContinuePreCounter, 1)
	defer atomic.AddUint64(&m.ContinueCounter, 1)

	if len(m.ContinueMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ContinueMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FlowMock.Continue. %v", p)
			return
		}

		input := m.ContinueMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FlowMockContinueInput{p}, "Flow.Continue got unexpected parameters")

		return
	}

	if m.ContinueMock.mainExpectation != nil {

		input := m.ContinueMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FlowMockContinueInput{p}, "Flow.Continue got unexpected parameters")
		}

		return
	}

	if m.ContinueFunc == nil {
		m.t.Fatalf("Unexpected call to FlowMock.Continue. %v", p)
		return
	}

	m.ContinueFunc(p)
}

//ContinueMinimockCounter returns a count of FlowMock.ContinueFunc invocations
func (m *FlowMock) ContinueMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ContinueCounter)
}

//ContinueMinimockPreCounter returns the value of FlowMock.Continue invocations
func (m *FlowMock) ContinueMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ContinuePreCounter)
}

//ContinueFinished returns true if mock invocations count is ok
func (m *FlowMock) ContinueFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ContinueMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ContinueCounter) == uint64(len(m.ContinueMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ContinueMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ContinueCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ContinueFunc != nil {
		return atomic.LoadUint64(&m.ContinueCounter) > 0
	}

	return true
}

type mFlowMockHandle struct {
	mock              *FlowMock
	mainExpectation   *FlowMockHandleExpectation
	expectationSeries []*FlowMockHandleExpectation
}

type FlowMockHandleExpectation struct {
	input  *FlowMockHandleInput
	result *FlowMockHandleResult
}

type FlowMockHandleInput struct {
	p  context.Context
	p1 Handle
}

type FlowMockHandleResult struct {
	r error
}

//Expect specifies that invocation of Flow.Handle is expected from 1 to Infinity times
func (m *mFlowMockHandle) Expect(p context.Context, p1 Handle) *mFlowMockHandle {
	m.mock.HandleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockHandleExpectation{}
	}
	m.mainExpectation.input = &FlowMockHandleInput{p, p1}
	return m
}

//Return specifies results of invocation of Flow.Handle
func (m *mFlowMockHandle) Return(r error) *FlowMock {
	m.mock.HandleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockHandleExpectation{}
	}
	m.mainExpectation.result = &FlowMockHandleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Flow.Handle is expected once
func (m *mFlowMockHandle) ExpectOnce(p context.Context, p1 Handle) *FlowMockHandleExpectation {
	m.mock.HandleFunc = nil
	m.mainExpectation = nil

	expectation := &FlowMockHandleExpectation{}
	expectation.input = &FlowMockHandleInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FlowMockHandleExpectation) Return(r error) {
	e.result = &FlowMockHandleResult{r}
}

//Set uses given function f as a mock of Flow.Handle method
func (m *mFlowMockHandle) Set(f func(p context.Context, p1 Handle) (r error)) *FlowMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HandleFunc = f
	return m.mock
}

//Handle implements github.com/insolar/insolar/insolar/flow.Flow interface
func (m *FlowMock) Handle(p context.Context, p1 Handle) (r error) {
	counter := atomic.AddUint64(&m.HandlePreCounter, 1)
	defer atomic.AddUint64(&m.HandleCounter, 1)

	if len(m.HandleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HandleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FlowMock.Handle. %v %v", p, p1)
			return
		}

		input := m.HandleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FlowMockHandleInput{p, p1}, "Flow.Handle got unexpected parameters")

		result := m.HandleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FlowMock.Handle")
			return
		}

		r = result.r

		return
	}

	if m.HandleMock.mainExpectation != nil {

		input := m.HandleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FlowMockHandleInput{p, p1}, "Flow.Handle got unexpected parameters")
		}

		result := m.HandleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FlowMock.Handle")
		}

		r = result.r

		return
	}

	if m.HandleFunc == nil {
		m.t.Fatalf("Unexpected call to FlowMock.Handle. %v %v", p, p1)
		return
	}

	return m.HandleFunc(p, p1)
}

//HandleMinimockCounter returns a count of FlowMock.HandleFunc invocations
func (m *FlowMock) HandleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HandleCounter)
}

//HandleMinimockPreCounter returns the value of FlowMock.Handle invocations
func (m *FlowMock) HandleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HandlePreCounter)
}

//HandleFinished returns true if mock invocations count is ok
func (m *FlowMock) HandleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HandleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HandleCounter) == uint64(len(m.HandleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HandleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HandleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HandleFunc != nil {
		return atomic.LoadUint64(&m.HandleCounter) > 0
	}

	return true
}

type mFlowMockMigrate struct {
	mock              *FlowMock
	mainExpectation   *FlowMockMigrateExpectation
	expectationSeries []*FlowMockMigrateExpectation
}

type FlowMockMigrateExpectation struct {
	input  *FlowMockMigrateInput
	result *FlowMockMigrateResult
}

type FlowMockMigrateInput struct {
	p  context.Context
	p1 Handle
}

type FlowMockMigrateResult struct {
	r error
}

//Expect specifies that invocation of Flow.Migrate is expected from 1 to Infinity times
func (m *mFlowMockMigrate) Expect(p context.Context, p1 Handle) *mFlowMockMigrate {
	m.mock.MigrateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockMigrateExpectation{}
	}
	m.mainExpectation.input = &FlowMockMigrateInput{p, p1}
	return m
}

//Return specifies results of invocation of Flow.Migrate
func (m *mFlowMockMigrate) Return(r error) *FlowMock {
	m.mock.MigrateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockMigrateExpectation{}
	}
	m.mainExpectation.result = &FlowMockMigrateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Flow.Migrate is expected once
func (m *mFlowMockMigrate) ExpectOnce(p context.Context, p1 Handle) *FlowMockMigrateExpectation {
	m.mock.MigrateFunc = nil
	m.mainExpectation = nil

	expectation := &FlowMockMigrateExpectation{}
	expectation.input = &FlowMockMigrateInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FlowMockMigrateExpectation) Return(r error) {
	e.result = &FlowMockMigrateResult{r}
}

//Set uses given function f as a mock of Flow.Migrate method
func (m *mFlowMockMigrate) Set(f func(p context.Context, p1 Handle) (r error)) *FlowMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MigrateFunc = f
	return m.mock
}

//Migrate implements github.com/insolar/insolar/insolar/flow.Flow interface
func (m *FlowMock) Migrate(p context.Context, p1 Handle) (r error) {
	counter := atomic.AddUint64(&m.MigratePreCounter, 1)
	defer atomic.AddUint64(&m.MigrateCounter, 1)

	if len(m.MigrateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MigrateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FlowMock.Migrate. %v %v", p, p1)
			return
		}

		input := m.MigrateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FlowMockMigrateInput{p, p1}, "Flow.Migrate got unexpected parameters")

		result := m.MigrateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FlowMock.Migrate")
			return
		}

		r = result.r

		return
	}

	if m.MigrateMock.mainExpectation != nil {

		input := m.MigrateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FlowMockMigrateInput{p, p1}, "Flow.Migrate got unexpected parameters")
		}

		result := m.MigrateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FlowMock.Migrate")
		}

		r = result.r

		return
	}

	if m.MigrateFunc == nil {
		m.t.Fatalf("Unexpected call to FlowMock.Migrate. %v %v", p, p1)
		return
	}

	return m.MigrateFunc(p, p1)
}

//MigrateMinimockCounter returns a count of FlowMock.MigrateFunc invocations
func (m *FlowMock) MigrateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MigrateCounter)
}

//MigrateMinimockPreCounter returns the value of FlowMock.Migrate invocations
func (m *FlowMock) MigrateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MigratePreCounter)
}

//MigrateFinished returns true if mock invocations count is ok
func (m *FlowMock) MigrateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MigrateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MigrateCounter) == uint64(len(m.MigrateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MigrateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MigrateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MigrateFunc != nil {
		return atomic.LoadUint64(&m.MigrateCounter) > 0
	}

	return true
}

type mFlowMockProcedure struct {
	mock              *FlowMock
	mainExpectation   *FlowMockProcedureExpectation
	expectationSeries []*FlowMockProcedureExpectation
}

type FlowMockProcedureExpectation struct {
	input  *FlowMockProcedureInput
	result *FlowMockProcedureResult
}

type FlowMockProcedureInput struct {
	p  context.Context
	p1 Procedure
}

type FlowMockProcedureResult struct {
	r error
}

//Expect specifies that invocation of Flow.Procedure is expected from 1 to Infinity times
func (m *mFlowMockProcedure) Expect(p context.Context, p1 Procedure) *mFlowMockProcedure {
	m.mock.ProcedureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockProcedureExpectation{}
	}
	m.mainExpectation.input = &FlowMockProcedureInput{p, p1}
	return m
}

//Return specifies results of invocation of Flow.Procedure
func (m *mFlowMockProcedure) Return(r error) *FlowMock {
	m.mock.ProcedureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FlowMockProcedureExpectation{}
	}
	m.mainExpectation.result = &FlowMockProcedureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Flow.Procedure is expected once
func (m *mFlowMockProcedure) ExpectOnce(p context.Context, p1 Procedure) *FlowMockProcedureExpectation {
	m.mock.ProcedureFunc = nil
	m.mainExpectation = nil

	expectation := &FlowMockProcedureExpectation{}
	expectation.input = &FlowMockProcedureInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FlowMockProcedureExpectation) Return(r error) {
	e.result = &FlowMockProcedureResult{r}
}

//Set uses given function f as a mock of Flow.Procedure method
func (m *mFlowMockProcedure) Set(f func(p context.Context, p1 Procedure) (r error)) *FlowMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ProcedureFunc = f
	return m.mock
}

//Procedure implements github.com/insolar/insolar/insolar/flow.Flow interface
func (m *FlowMock) Procedure(p context.Context, p1 Procedure) (r error) {
	counter := atomic.AddUint64(&m.ProcedurePreCounter, 1)
	defer atomic.AddUint64(&m.ProcedureCounter, 1)

	if len(m.ProcedureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ProcedureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FlowMock.Procedure. %v %v", p, p1)
			return
		}

		input := m.ProcedureMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FlowMockProcedureInput{p, p1}, "Flow.Procedure got unexpected parameters")

		result := m.ProcedureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FlowMock.Procedure")
			return
		}

		r = result.r

		return
	}

	if m.ProcedureMock.mainExpectation != nil {

		input := m.ProcedureMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FlowMockProcedureInput{p, p1}, "Flow.Procedure got unexpected parameters")
		}

		result := m.ProcedureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FlowMock.Procedure")
		}

		r = result.r

		return
	}

	if m.ProcedureFunc == nil {
		m.t.Fatalf("Unexpected call to FlowMock.Procedure. %v %v", p, p1)
		return
	}

	return m.ProcedureFunc(p, p1)
}

//ProcedureMinimockCounter returns a count of FlowMock.ProcedureFunc invocations
func (m *FlowMock) ProcedureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ProcedureCounter)
}

//ProcedureMinimockPreCounter returns the value of FlowMock.Procedure invocations
func (m *FlowMock) ProcedureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ProcedurePreCounter)
}

//ProcedureFinished returns true if mock invocations count is ok
func (m *FlowMock) ProcedureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ProcedureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ProcedureCounter) == uint64(len(m.ProcedureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ProcedureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ProcedureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ProcedureFunc != nil {
		return atomic.LoadUint64(&m.ProcedureCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FlowMock) ValidateCallCounters() {

	if !m.ContinueFinished() {
		m.t.Fatal("Expected call to FlowMock.Continue")
	}

	if !m.HandleFinished() {
		m.t.Fatal("Expected call to FlowMock.Handle")
	}

	if !m.MigrateFinished() {
		m.t.Fatal("Expected call to FlowMock.Migrate")
	}

	if !m.ProcedureFinished() {
		m.t.Fatal("Expected call to FlowMock.Procedure")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FlowMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FlowMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FlowMock) MinimockFinish() {

	if !m.ContinueFinished() {
		m.t.Fatal("Expected call to FlowMock.Continue")
	}

	if !m.HandleFinished() {
		m.t.Fatal("Expected call to FlowMock.Handle")
	}

	if !m.MigrateFinished() {
		m.t.Fatal("Expected call to FlowMock.Migrate")
	}

	if !m.ProcedureFinished() {
		m.t.Fatal("Expected call to FlowMock.Procedure")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FlowMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FlowMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ContinueFinished()
		ok = ok && m.HandleFinished()
		ok = ok && m.MigrateFinished()
		ok = ok && m.ProcedureFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ContinueFinished() {
				m.t.Error("Expected call to FlowMock.Continue")
			}

			if !m.HandleFinished() {
				m.t.Error("Expected call to FlowMock.Handle")
			}

			if !m.MigrateFinished() {
				m.t.Error("Expected call to FlowMock.Migrate")
			}

			if !m.ProcedureFinished() {
				m.t.Error("Expected call to FlowMock.Procedure")
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
func (m *FlowMock) AllMocksCalled() bool {

	if !m.ContinueFinished() {
		return false
	}

	if !m.HandleFinished() {
		return false
	}

	if !m.MigrateFinished() {
		return false
	}

	if !m.ProcedureFinished() {
		return false
	}

	return true
}

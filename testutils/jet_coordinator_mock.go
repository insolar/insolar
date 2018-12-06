package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetCoordinator" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetCoordinatorMock implements github.com/insolar/insolar/core.JetCoordinator
type JetCoordinatorMock struct {
	t minimock.Tester

	GetActiveNodesFunc       func(p core.PulseNumber) (r []core.Node, r1 error)
	GetActiveNodesCounter    uint64
	GetActiveNodesPreCounter uint64
	GetActiveNodesMock       mJetCoordinatorMockGetActiveNodes

	IsAuthorizedFunc       func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error)
	IsAuthorizedCounter    uint64
	IsAuthorizedPreCounter uint64
	IsAuthorizedMock       mJetCoordinatorMockIsAuthorized

	QueryRoleFunc       func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) (r []core.RecordRef, r1 error)
	QueryRoleCounter    uint64
	QueryRolePreCounter uint64
	QueryRoleMock       mJetCoordinatorMockQueryRole
}

//NewJetCoordinatorMock returns a mock for github.com/insolar/insolar/core.JetCoordinator
func NewJetCoordinatorMock(t minimock.Tester) *JetCoordinatorMock {
	m := &JetCoordinatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetActiveNodesMock = mJetCoordinatorMockGetActiveNodes{mock: m}
	m.IsAuthorizedMock = mJetCoordinatorMockIsAuthorized{mock: m}
	m.QueryRoleMock = mJetCoordinatorMockQueryRole{mock: m}

	return m
}

type mJetCoordinatorMockGetActiveNodes struct {
	mock              *JetCoordinatorMock
	mainExpectation   *JetCoordinatorMockGetActiveNodesExpectation
	expectationSeries []*JetCoordinatorMockGetActiveNodesExpectation
}

type JetCoordinatorMockGetActiveNodesExpectation struct {
	input  *JetCoordinatorMockGetActiveNodesInput
	result *JetCoordinatorMockGetActiveNodesResult
}

type JetCoordinatorMockGetActiveNodesInput struct {
	p core.PulseNumber
}

type JetCoordinatorMockGetActiveNodesResult struct {
	r  []core.Node
	r1 error
}

//Expect specifies that invocation of JetCoordinator.GetActiveNodes is expected from 1 to Infinity times
func (m *mJetCoordinatorMockGetActiveNodes) Expect(p core.PulseNumber) *mJetCoordinatorMockGetActiveNodes {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCoordinatorMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.input = &JetCoordinatorMockGetActiveNodesInput{p}
	return m
}

//Return specifies results of invocation of JetCoordinator.GetActiveNodes
func (m *mJetCoordinatorMockGetActiveNodes) Return(r []core.Node, r1 error) *JetCoordinatorMock {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCoordinatorMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &JetCoordinatorMockGetActiveNodesResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of JetCoordinator.GetActiveNodes is expected once
func (m *mJetCoordinatorMockGetActiveNodes) ExpectOnce(p core.PulseNumber) *JetCoordinatorMockGetActiveNodesExpectation {
	m.mock.GetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &JetCoordinatorMockGetActiveNodesExpectation{}
	expectation.input = &JetCoordinatorMockGetActiveNodesInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetCoordinatorMockGetActiveNodesExpectation) Return(r []core.Node, r1 error) {
	e.result = &JetCoordinatorMockGetActiveNodesResult{r, r1}
}

//Set uses given function f as a mock of JetCoordinator.GetActiveNodes method
func (m *mJetCoordinatorMockGetActiveNodes) Set(f func(p core.PulseNumber) (r []core.Node, r1 error)) *JetCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesFunc = f
	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) GetActiveNodes(p core.PulseNumber) (r []core.Node, r1 error) {
	counter := atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetCoordinatorMock.GetActiveNodes. %v", p)
			return
		}

		input := m.GetActiveNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetCoordinatorMockGetActiveNodesInput{p}, "JetCoordinator.GetActiveNodes got unexpected parameters")

		result := m.GetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetCoordinatorMock.GetActiveNodes")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesMock.mainExpectation != nil {

		input := m.GetActiveNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetCoordinatorMockGetActiveNodesInput{p}, "JetCoordinator.GetActiveNodes got unexpected parameters")
		}

		result := m.GetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetCoordinatorMock.GetActiveNodes")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to JetCoordinatorMock.GetActiveNodes. %v", p)
		return
	}

	return m.GetActiveNodesFunc(p)
}

//GetActiveNodesMinimockCounter returns a count of JetCoordinatorMock.GetActiveNodesFunc invocations
func (m *JetCoordinatorMock) GetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesCounter)
}

//GetActiveNodesMinimockPreCounter returns the value of JetCoordinatorMock.GetActiveNodes invocations
func (m *JetCoordinatorMock) GetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesPreCounter)
}

//GetActiveNodesFinished returns true if mock invocations count is ok
func (m *JetCoordinatorMock) GetActiveNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodesCounter) == uint64(len(m.GetActiveNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodesFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodesCounter) > 0
	}

	return true
}

type mJetCoordinatorMockIsAuthorized struct {
	mock              *JetCoordinatorMock
	mainExpectation   *JetCoordinatorMockIsAuthorizedExpectation
	expectationSeries []*JetCoordinatorMockIsAuthorizedExpectation
}

type JetCoordinatorMockIsAuthorizedExpectation struct {
	input  *JetCoordinatorMockIsAuthorizedInput
	result *JetCoordinatorMockIsAuthorizedResult
}

type JetCoordinatorMockIsAuthorizedInput struct {
	p  context.Context
	p1 core.DynamicRole
	p2 *core.RecordRef
	p3 core.PulseNumber
	p4 core.RecordRef
}

type JetCoordinatorMockIsAuthorizedResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of JetCoordinator.IsAuthorized is expected from 1 to Infinity times
func (m *mJetCoordinatorMockIsAuthorized) Expect(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) *mJetCoordinatorMockIsAuthorized {
	m.mock.IsAuthorizedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCoordinatorMockIsAuthorizedExpectation{}
	}
	m.mainExpectation.input = &JetCoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of JetCoordinator.IsAuthorized
func (m *mJetCoordinatorMockIsAuthorized) Return(r bool, r1 error) *JetCoordinatorMock {
	m.mock.IsAuthorizedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCoordinatorMockIsAuthorizedExpectation{}
	}
	m.mainExpectation.result = &JetCoordinatorMockIsAuthorizedResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of JetCoordinator.IsAuthorized is expected once
func (m *mJetCoordinatorMockIsAuthorized) ExpectOnce(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) *JetCoordinatorMockIsAuthorizedExpectation {
	m.mock.IsAuthorizedFunc = nil
	m.mainExpectation = nil

	expectation := &JetCoordinatorMockIsAuthorizedExpectation{}
	expectation.input = &JetCoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetCoordinatorMockIsAuthorizedExpectation) Return(r bool, r1 error) {
	e.result = &JetCoordinatorMockIsAuthorizedResult{r, r1}
}

//Set uses given function f as a mock of JetCoordinator.IsAuthorized method
func (m *mJetCoordinatorMockIsAuthorized) Set(f func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error)) *JetCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAuthorizedFunc = f
	return m.mock
}

//IsAuthorized implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) IsAuthorized(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.IsAuthorizedPreCounter, 1)
	defer atomic.AddUint64(&m.IsAuthorizedCounter, 1)

	if len(m.IsAuthorizedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAuthorizedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetCoordinatorMock.IsAuthorized. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.IsAuthorizedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetCoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}, "JetCoordinator.IsAuthorized got unexpected parameters")

		result := m.IsAuthorizedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetCoordinatorMock.IsAuthorized")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IsAuthorizedMock.mainExpectation != nil {

		input := m.IsAuthorizedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetCoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}, "JetCoordinator.IsAuthorized got unexpected parameters")
		}

		result := m.IsAuthorizedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetCoordinatorMock.IsAuthorized")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IsAuthorizedFunc == nil {
		m.t.Fatalf("Unexpected call to JetCoordinatorMock.IsAuthorized. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.IsAuthorizedFunc(p, p1, p2, p3, p4)
}

//IsAuthorizedMinimockCounter returns a count of JetCoordinatorMock.IsAuthorizedFunc invocations
func (m *JetCoordinatorMock) IsAuthorizedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAuthorizedCounter)
}

//IsAuthorizedMinimockPreCounter returns the value of JetCoordinatorMock.IsAuthorized invocations
func (m *JetCoordinatorMock) IsAuthorizedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAuthorizedPreCounter)
}

//IsAuthorizedFinished returns true if mock invocations count is ok
func (m *JetCoordinatorMock) IsAuthorizedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsAuthorizedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsAuthorizedCounter) == uint64(len(m.IsAuthorizedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsAuthorizedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsAuthorizedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsAuthorizedFunc != nil {
		return atomic.LoadUint64(&m.IsAuthorizedCounter) > 0
	}

	return true
}

type mJetCoordinatorMockQueryRole struct {
	mock              *JetCoordinatorMock
	mainExpectation   *JetCoordinatorMockQueryRoleExpectation
	expectationSeries []*JetCoordinatorMockQueryRoleExpectation
}

type JetCoordinatorMockQueryRoleExpectation struct {
	input  *JetCoordinatorMockQueryRoleInput
	result *JetCoordinatorMockQueryRoleResult
}

type JetCoordinatorMockQueryRoleInput struct {
	p  context.Context
	p1 core.DynamicRole
	p2 *core.RecordRef
	p3 core.PulseNumber
}

type JetCoordinatorMockQueryRoleResult struct {
	r  []core.RecordRef
	r1 error
}

//Expect specifies that invocation of JetCoordinator.QueryRole is expected from 1 to Infinity times
func (m *mJetCoordinatorMockQueryRole) Expect(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) *mJetCoordinatorMockQueryRole {
	m.mock.QueryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCoordinatorMockQueryRoleExpectation{}
	}
	m.mainExpectation.input = &JetCoordinatorMockQueryRoleInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of JetCoordinator.QueryRole
func (m *mJetCoordinatorMockQueryRole) Return(r []core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.QueryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetCoordinatorMockQueryRoleExpectation{}
	}
	m.mainExpectation.result = &JetCoordinatorMockQueryRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of JetCoordinator.QueryRole is expected once
func (m *mJetCoordinatorMockQueryRole) ExpectOnce(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) *JetCoordinatorMockQueryRoleExpectation {
	m.mock.QueryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &JetCoordinatorMockQueryRoleExpectation{}
	expectation.input = &JetCoordinatorMockQueryRoleInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetCoordinatorMockQueryRoleExpectation) Return(r []core.RecordRef, r1 error) {
	e.result = &JetCoordinatorMockQueryRoleResult{r, r1}
}

//Set uses given function f as a mock of JetCoordinator.QueryRole method
func (m *mJetCoordinatorMockQueryRole) Set(f func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) (r []core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.QueryRoleFunc = f
	return m.mock
}

//QueryRole implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) QueryRole(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) (r []core.RecordRef, r1 error) {
	counter := atomic.AddUint64(&m.QueryRolePreCounter, 1)
	defer atomic.AddUint64(&m.QueryRoleCounter, 1)

	if len(m.QueryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.QueryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetCoordinatorMock.QueryRole. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.QueryRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetCoordinatorMockQueryRoleInput{p, p1, p2, p3}, "JetCoordinator.QueryRole got unexpected parameters")

		result := m.QueryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetCoordinatorMock.QueryRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.QueryRoleMock.mainExpectation != nil {

		input := m.QueryRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetCoordinatorMockQueryRoleInput{p, p1, p2, p3}, "JetCoordinator.QueryRole got unexpected parameters")
		}

		result := m.QueryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetCoordinatorMock.QueryRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.QueryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to JetCoordinatorMock.QueryRole. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.QueryRoleFunc(p, p1, p2, p3)
}

//QueryRoleMinimockCounter returns a count of JetCoordinatorMock.QueryRoleFunc invocations
func (m *JetCoordinatorMock) QueryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.QueryRoleCounter)
}

//QueryRoleMinimockPreCounter returns the value of JetCoordinatorMock.QueryRole invocations
func (m *JetCoordinatorMock) QueryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.QueryRolePreCounter)
}

//QueryRoleFinished returns true if mock invocations count is ok
func (m *JetCoordinatorMock) QueryRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.QueryRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.QueryRoleCounter) == uint64(len(m.QueryRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.QueryRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.QueryRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.QueryRoleFunc != nil {
		return atomic.LoadUint64(&m.QueryRoleCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetCoordinatorMock) ValidateCallCounters() {

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to JetCoordinatorMock.GetActiveNodes")
	}

	if !m.IsAuthorizedFinished() {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsAuthorized")
	}

	if !m.QueryRoleFinished() {
		m.t.Fatal("Expected call to JetCoordinatorMock.QueryRole")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetCoordinatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetCoordinatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetCoordinatorMock) MinimockFinish() {

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to JetCoordinatorMock.GetActiveNodes")
	}

	if !m.IsAuthorizedFinished() {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsAuthorized")
	}

	if !m.QueryRoleFinished() {
		m.t.Fatal("Expected call to JetCoordinatorMock.QueryRole")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetCoordinatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetCoordinatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetActiveNodesFinished()
		ok = ok && m.IsAuthorizedFinished()
		ok = ok && m.QueryRoleFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetActiveNodesFinished() {
				m.t.Error("Expected call to JetCoordinatorMock.GetActiveNodes")
			}

			if !m.IsAuthorizedFinished() {
				m.t.Error("Expected call to JetCoordinatorMock.IsAuthorized")
			}

			if !m.QueryRoleFinished() {
				m.t.Error("Expected call to JetCoordinatorMock.QueryRole")
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
func (m *JetCoordinatorMock) AllMocksCalled() bool {

	if !m.GetActiveNodesFinished() {
		return false
	}

	if !m.IsAuthorizedFinished() {
		return false
	}

	if !m.QueryRoleFinished() {
		return false
	}

	return true
}

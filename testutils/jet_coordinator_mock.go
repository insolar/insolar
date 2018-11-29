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
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockGetActiveNodesParams
}

//JetCoordinatorMockGetActiveNodesParams represents input parameters of the JetCoordinator.GetActiveNodes
type JetCoordinatorMockGetActiveNodesParams struct {
	p core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.GetActiveNodes
func (m *mJetCoordinatorMockGetActiveNodes) Expect(p core.PulseNumber) *mJetCoordinatorMockGetActiveNodes {
	m.mockExpectations = &JetCoordinatorMockGetActiveNodesParams{p}
	return m
}

//Return sets up a mock for JetCoordinator.GetActiveNodes to return Return's arguments
func (m *mJetCoordinatorMockGetActiveNodes) Return(r []core.Node, r1 error) *JetCoordinatorMock {
	m.mock.GetActiveNodesFunc = func(p core.PulseNumber) ([]core.Node, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.GetActiveNodes method
func (m *mJetCoordinatorMockGetActiveNodes) Set(f func(p core.PulseNumber) (r []core.Node, r1 error)) *JetCoordinatorMock {
	m.mock.GetActiveNodesFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) GetActiveNodes(p core.PulseNumber) (r []core.Node, r1 error) {
	atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if m.GetActiveNodesMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetActiveNodesMock.mockExpectations, JetCoordinatorMockGetActiveNodesParams{p},
			"JetCoordinator.GetActiveNodes got unexpected parameters")

		if m.GetActiveNodesFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.GetActiveNodes")

			return
		}
	}

	if m.GetActiveNodesFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.GetActiveNodes")
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

type mJetCoordinatorMockIsAuthorized struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockIsAuthorizedParams
}

//JetCoordinatorMockIsAuthorizedParams represents input parameters of the JetCoordinator.IsAuthorized
type JetCoordinatorMockIsAuthorizedParams struct {
	p  context.Context
	p1 core.DynamicRole
	p2 *core.RecordRef
	p3 core.PulseNumber
	p4 core.RecordRef
}

//Expect sets up expected params for the JetCoordinator.IsAuthorized
func (m *mJetCoordinatorMockIsAuthorized) Expect(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) *mJetCoordinatorMockIsAuthorized {
	m.mockExpectations = &JetCoordinatorMockIsAuthorizedParams{p, p1, p2, p3, p4}
	return m
}

//Return sets up a mock for JetCoordinator.IsAuthorized to return Return's arguments
func (m *mJetCoordinatorMockIsAuthorized) Return(r bool, r1 error) *JetCoordinatorMock {
	m.mock.IsAuthorizedFunc = func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) (bool, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.IsAuthorized method
func (m *mJetCoordinatorMockIsAuthorized) Set(f func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error)) *JetCoordinatorMock {
	m.mock.IsAuthorizedFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IsAuthorized implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) IsAuthorized(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error) {
	atomic.AddUint64(&m.IsAuthorizedPreCounter, 1)
	defer atomic.AddUint64(&m.IsAuthorizedCounter, 1)

	if m.IsAuthorizedMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IsAuthorizedMock.mockExpectations, JetCoordinatorMockIsAuthorizedParams{p, p1, p2, p3, p4},
			"JetCoordinator.IsAuthorized got unexpected parameters")

		if m.IsAuthorizedFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.IsAuthorized")

			return
		}
	}

	if m.IsAuthorizedFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.IsAuthorized")
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

type mJetCoordinatorMockQueryRole struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockQueryRoleParams
}

//JetCoordinatorMockQueryRoleParams represents input parameters of the JetCoordinator.QueryRole
type JetCoordinatorMockQueryRoleParams struct {
	p  context.Context
	p1 core.DynamicRole
	p2 *core.RecordRef
	p3 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.QueryRole
func (m *mJetCoordinatorMockQueryRole) Expect(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) *mJetCoordinatorMockQueryRole {
	m.mockExpectations = &JetCoordinatorMockQueryRoleParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for JetCoordinator.QueryRole to return Return's arguments
func (m *mJetCoordinatorMockQueryRole) Return(r []core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.QueryRoleFunc = func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) ([]core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.QueryRole method
func (m *mJetCoordinatorMockQueryRole) Set(f func(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) (r []core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.QueryRoleFunc = f
	m.mockExpectations = nil
	return m.mock
}

//QueryRole implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) QueryRole(p context.Context, p1 core.DynamicRole, p2 *core.RecordRef, p3 core.PulseNumber) (r []core.RecordRef, r1 error) {
	atomic.AddUint64(&m.QueryRolePreCounter, 1)
	defer atomic.AddUint64(&m.QueryRoleCounter, 1)

	if m.QueryRoleMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.QueryRoleMock.mockExpectations, JetCoordinatorMockQueryRoleParams{p, p1, p2, p3},
			"JetCoordinator.QueryRole got unexpected parameters")

		if m.QueryRoleFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.QueryRole")

			return
		}
	}

	if m.QueryRoleFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.QueryRole")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetCoordinatorMock) ValidateCallCounters() {

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.GetActiveNodes")
	}

	if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsAuthorized")
	}

	if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
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

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.GetActiveNodes")
	}

	if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsAuthorized")
	}

	if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
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
		ok = ok && (m.GetActiveNodesFunc == nil || atomic.LoadUint64(&m.GetActiveNodesCounter) > 0)
		ok = ok && (m.IsAuthorizedFunc == nil || atomic.LoadUint64(&m.IsAuthorizedCounter) > 0)
		ok = ok && (m.QueryRoleFunc == nil || atomic.LoadUint64(&m.QueryRoleCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.GetActiveNodes")
			}

			if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.IsAuthorized")
			}

			if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
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

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		return false
	}

	if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
		return false
	}

	if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
		return false
	}

	return true
}

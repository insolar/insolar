package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeStorage" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeStorageMock implements github.com/insolar/insolar/ledger/storage.NodeStorage
type NodeStorageMock struct {
	t minimock.Tester

	GetActiveNodesFunc       func(p core.PulseNumber) (r []core.Node, r1 error)
	GetActiveNodesCounter    uint64
	GetActiveNodesPreCounter uint64
	GetActiveNodesMock       mNodeStorageMockGetActiveNodes

	GetActiveNodesByRoleFunc       func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error)
	GetActiveNodesByRoleCounter    uint64
	GetActiveNodesByRolePreCounter uint64
	GetActiveNodesByRoleMock       mNodeStorageMockGetActiveNodesByRole

	RemoveActiveNodesUntilFunc       func(p core.PulseNumber)
	RemoveActiveNodesUntilCounter    uint64
	RemoveActiveNodesUntilPreCounter uint64
	RemoveActiveNodesUntilMock       mNodeStorageMockRemoveActiveNodesUntil

	SetActiveNodesFunc       func(p core.PulseNumber, p1 []core.Node) (r error)
	SetActiveNodesCounter    uint64
	SetActiveNodesPreCounter uint64
	SetActiveNodesMock       mNodeStorageMockSetActiveNodes
}

//NewNodeStorageMock returns a mock for github.com/insolar/insolar/ledger/storage.NodeStorage
func NewNodeStorageMock(t minimock.Tester) *NodeStorageMock {
	m := &NodeStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetActiveNodesMock = mNodeStorageMockGetActiveNodes{mock: m}
	m.GetActiveNodesByRoleMock = mNodeStorageMockGetActiveNodesByRole{mock: m}
	m.RemoveActiveNodesUntilMock = mNodeStorageMockRemoveActiveNodesUntil{mock: m}
	m.SetActiveNodesMock = mNodeStorageMockSetActiveNodes{mock: m}

	return m
}

type mNodeStorageMockGetActiveNodes struct {
	mock              *NodeStorageMock
	mainExpectation   *NodeStorageMockGetActiveNodesExpectation
	expectationSeries []*NodeStorageMockGetActiveNodesExpectation
}

type NodeStorageMockGetActiveNodesExpectation struct {
	input  *NodeStorageMockGetActiveNodesInput
	result *NodeStorageMockGetActiveNodesResult
}

type NodeStorageMockGetActiveNodesInput struct {
	p core.PulseNumber
}

type NodeStorageMockGetActiveNodesResult struct {
	r  []core.Node
	r1 error
}

//Expect specifies that invocation of NodeStorage.GetActiveNodes is expected from 1 to Infinity times
func (m *mNodeStorageMockGetActiveNodes) Expect(p core.PulseNumber) *mNodeStorageMockGetActiveNodes {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.input = &NodeStorageMockGetActiveNodesInput{p}
	return m
}

//Return specifies results of invocation of NodeStorage.GetActiveNodes
func (m *mNodeStorageMockGetActiveNodes) Return(r []core.Node, r1 error) *NodeStorageMock {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &NodeStorageMockGetActiveNodesResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStorage.GetActiveNodes is expected once
func (m *mNodeStorageMockGetActiveNodes) ExpectOnce(p core.PulseNumber) *NodeStorageMockGetActiveNodesExpectation {
	m.mock.GetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStorageMockGetActiveNodesExpectation{}
	expectation.input = &NodeStorageMockGetActiveNodesInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStorageMockGetActiveNodesExpectation) Return(r []core.Node, r1 error) {
	e.result = &NodeStorageMockGetActiveNodesResult{r, r1}
}

//Set uses given function f as a mock of NodeStorage.GetActiveNodes method
func (m *mNodeStorageMockGetActiveNodes) Set(f func(p core.PulseNumber) (r []core.Node, r1 error)) *NodeStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesFunc = f
	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/ledger/storage.NodeStorage interface
func (m *NodeStorageMock) GetActiveNodes(p core.PulseNumber) (r []core.Node, r1 error) {
	counter := atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStorageMock.GetActiveNodes. %v", p)
			return
		}

		input := m.GetActiveNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStorageMockGetActiveNodesInput{p}, "NodeStorage.GetActiveNodes got unexpected parameters")

		result := m.GetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStorageMock.GetActiveNodes")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesMock.mainExpectation != nil {

		input := m.GetActiveNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStorageMockGetActiveNodesInput{p}, "NodeStorage.GetActiveNodes got unexpected parameters")
		}

		result := m.GetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStorageMock.GetActiveNodes")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStorageMock.GetActiveNodes. %v", p)
		return
	}

	return m.GetActiveNodesFunc(p)
}

//GetActiveNodesMinimockCounter returns a count of NodeStorageMock.GetActiveNodesFunc invocations
func (m *NodeStorageMock) GetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesCounter)
}

//GetActiveNodesMinimockPreCounter returns the value of NodeStorageMock.GetActiveNodes invocations
func (m *NodeStorageMock) GetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesPreCounter)
}

//GetActiveNodesFinished returns true if mock invocations count is ok
func (m *NodeStorageMock) GetActiveNodesFinished() bool {
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

type mNodeStorageMockGetActiveNodesByRole struct {
	mock              *NodeStorageMock
	mainExpectation   *NodeStorageMockGetActiveNodesByRoleExpectation
	expectationSeries []*NodeStorageMockGetActiveNodesByRoleExpectation
}

type NodeStorageMockGetActiveNodesByRoleExpectation struct {
	input  *NodeStorageMockGetActiveNodesByRoleInput
	result *NodeStorageMockGetActiveNodesByRoleResult
}

type NodeStorageMockGetActiveNodesByRoleInput struct {
	p  core.PulseNumber
	p1 core.StaticRole
}

type NodeStorageMockGetActiveNodesByRoleResult struct {
	r  []core.Node
	r1 error
}

//Expect specifies that invocation of NodeStorage.GetActiveNodesByRole is expected from 1 to Infinity times
func (m *mNodeStorageMockGetActiveNodesByRole) Expect(p core.PulseNumber, p1 core.StaticRole) *mNodeStorageMockGetActiveNodesByRole {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockGetActiveNodesByRoleExpectation{}
	}
	m.mainExpectation.input = &NodeStorageMockGetActiveNodesByRoleInput{p, p1}
	return m
}

//Return specifies results of invocation of NodeStorage.GetActiveNodesByRole
func (m *mNodeStorageMockGetActiveNodesByRole) Return(r []core.Node, r1 error) *NodeStorageMock {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockGetActiveNodesByRoleExpectation{}
	}
	m.mainExpectation.result = &NodeStorageMockGetActiveNodesByRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStorage.GetActiveNodesByRole is expected once
func (m *mNodeStorageMockGetActiveNodesByRole) ExpectOnce(p core.PulseNumber, p1 core.StaticRole) *NodeStorageMockGetActiveNodesByRoleExpectation {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStorageMockGetActiveNodesByRoleExpectation{}
	expectation.input = &NodeStorageMockGetActiveNodesByRoleInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStorageMockGetActiveNodesByRoleExpectation) Return(r []core.Node, r1 error) {
	e.result = &NodeStorageMockGetActiveNodesByRoleResult{r, r1}
}

//Set uses given function f as a mock of NodeStorage.GetActiveNodesByRole method
func (m *mNodeStorageMockGetActiveNodesByRole) Set(f func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error)) *NodeStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesByRoleFunc = f
	return m.mock
}

//GetActiveNodesByRole implements github.com/insolar/insolar/ledger/storage.NodeStorage interface
func (m *NodeStorageMock) GetActiveNodesByRole(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error) {
	counter := atomic.AddUint64(&m.GetActiveNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesByRoleCounter, 1)

	if len(m.GetActiveNodesByRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesByRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStorageMock.GetActiveNodesByRole. %v %v", p, p1)
			return
		}

		input := m.GetActiveNodesByRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStorageMockGetActiveNodesByRoleInput{p, p1}, "NodeStorage.GetActiveNodesByRole got unexpected parameters")

		result := m.GetActiveNodesByRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStorageMock.GetActiveNodesByRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesByRoleMock.mainExpectation != nil {

		input := m.GetActiveNodesByRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStorageMockGetActiveNodesByRoleInput{p, p1}, "NodeStorage.GetActiveNodesByRole got unexpected parameters")
		}

		result := m.GetActiveNodesByRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStorageMock.GetActiveNodesByRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesByRoleFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStorageMock.GetActiveNodesByRole. %v %v", p, p1)
		return
	}

	return m.GetActiveNodesByRoleFunc(p, p1)
}

//GetActiveNodesByRoleMinimockCounter returns a count of NodeStorageMock.GetActiveNodesByRoleFunc invocations
func (m *NodeStorageMock) GetActiveNodesByRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter)
}

//GetActiveNodesByRoleMinimockPreCounter returns the value of NodeStorageMock.GetActiveNodesByRole invocations
func (m *NodeStorageMock) GetActiveNodesByRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRolePreCounter)
}

//GetActiveNodesByRoleFinished returns true if mock invocations count is ok
func (m *NodeStorageMock) GetActiveNodesByRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodesByRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == uint64(len(m.GetActiveNodesByRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodesByRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodesByRoleFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) > 0
	}

	return true
}

type mNodeStorageMockRemoveActiveNodesUntil struct {
	mock              *NodeStorageMock
	mainExpectation   *NodeStorageMockRemoveActiveNodesUntilExpectation
	expectationSeries []*NodeStorageMockRemoveActiveNodesUntilExpectation
}

type NodeStorageMockRemoveActiveNodesUntilExpectation struct {
	input *NodeStorageMockRemoveActiveNodesUntilInput
}

type NodeStorageMockRemoveActiveNodesUntilInput struct {
	p core.PulseNumber
}

//Expect specifies that invocation of NodeStorage.RemoveActiveNodesUntil is expected from 1 to Infinity times
func (m *mNodeStorageMockRemoveActiveNodesUntil) Expect(p core.PulseNumber) *mNodeStorageMockRemoveActiveNodesUntil {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockRemoveActiveNodesUntilExpectation{}
	}
	m.mainExpectation.input = &NodeStorageMockRemoveActiveNodesUntilInput{p}
	return m
}

//Return specifies results of invocation of NodeStorage.RemoveActiveNodesUntil
func (m *mNodeStorageMockRemoveActiveNodesUntil) Return() *NodeStorageMock {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockRemoveActiveNodesUntilExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeStorage.RemoveActiveNodesUntil is expected once
func (m *mNodeStorageMockRemoveActiveNodesUntil) ExpectOnce(p core.PulseNumber) *NodeStorageMockRemoveActiveNodesUntilExpectation {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStorageMockRemoveActiveNodesUntilExpectation{}
	expectation.input = &NodeStorageMockRemoveActiveNodesUntilInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeStorage.RemoveActiveNodesUntil method
func (m *mNodeStorageMockRemoveActiveNodesUntil) Set(f func(p core.PulseNumber)) *NodeStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveActiveNodesUntilFunc = f
	return m.mock
}

//RemoveActiveNodesUntil implements github.com/insolar/insolar/ledger/storage.NodeStorage interface
func (m *NodeStorageMock) RemoveActiveNodesUntil(p core.PulseNumber) {
	counter := atomic.AddUint64(&m.RemoveActiveNodesUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveActiveNodesUntilCounter, 1)

	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStorageMock.RemoveActiveNodesUntil. %v", p)
			return
		}

		input := m.RemoveActiveNodesUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStorageMockRemoveActiveNodesUntilInput{p}, "NodeStorage.RemoveActiveNodesUntil got unexpected parameters")

		return
	}

	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {

		input := m.RemoveActiveNodesUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStorageMockRemoveActiveNodesUntilInput{p}, "NodeStorage.RemoveActiveNodesUntil got unexpected parameters")
		}

		return
	}

	if m.RemoveActiveNodesUntilFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStorageMock.RemoveActiveNodesUntil. %v", p)
		return
	}

	m.RemoveActiveNodesUntilFunc(p)
}

//RemoveActiveNodesUntilMinimockCounter returns a count of NodeStorageMock.RemoveActiveNodesUntilFunc invocations
func (m *NodeStorageMock) RemoveActiveNodesUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter)
}

//RemoveActiveNodesUntilMinimockPreCounter returns the value of NodeStorageMock.RemoveActiveNodesUntil invocations
func (m *NodeStorageMock) RemoveActiveNodesUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilPreCounter)
}

//RemoveActiveNodesUntilFinished returns true if mock invocations count is ok
func (m *NodeStorageMock) RemoveActiveNodesUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) == uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	return true
}

type mNodeStorageMockSetActiveNodes struct {
	mock              *NodeStorageMock
	mainExpectation   *NodeStorageMockSetActiveNodesExpectation
	expectationSeries []*NodeStorageMockSetActiveNodesExpectation
}

type NodeStorageMockSetActiveNodesExpectation struct {
	input  *NodeStorageMockSetActiveNodesInput
	result *NodeStorageMockSetActiveNodesResult
}

type NodeStorageMockSetActiveNodesInput struct {
	p  core.PulseNumber
	p1 []core.Node
}

type NodeStorageMockSetActiveNodesResult struct {
	r error
}

//Expect specifies that invocation of NodeStorage.SetActiveNodes is expected from 1 to Infinity times
func (m *mNodeStorageMockSetActiveNodes) Expect(p core.PulseNumber, p1 []core.Node) *mNodeStorageMockSetActiveNodes {
	m.mock.SetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockSetActiveNodesExpectation{}
	}
	m.mainExpectation.input = &NodeStorageMockSetActiveNodesInput{p, p1}
	return m
}

//Return specifies results of invocation of NodeStorage.SetActiveNodes
func (m *mNodeStorageMockSetActiveNodes) Return(r error) *NodeStorageMock {
	m.mock.SetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStorageMockSetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &NodeStorageMockSetActiveNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStorage.SetActiveNodes is expected once
func (m *mNodeStorageMockSetActiveNodes) ExpectOnce(p core.PulseNumber, p1 []core.Node) *NodeStorageMockSetActiveNodesExpectation {
	m.mock.SetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStorageMockSetActiveNodesExpectation{}
	expectation.input = &NodeStorageMockSetActiveNodesInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStorageMockSetActiveNodesExpectation) Return(r error) {
	e.result = &NodeStorageMockSetActiveNodesResult{r}
}

//Set uses given function f as a mock of NodeStorage.SetActiveNodes method
func (m *mNodeStorageMockSetActiveNodes) Set(f func(p core.PulseNumber, p1 []core.Node) (r error)) *NodeStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetActiveNodesFunc = f
	return m.mock
}

//SetActiveNodes implements github.com/insolar/insolar/ledger/storage.NodeStorage interface
func (m *NodeStorageMock) SetActiveNodes(p core.PulseNumber, p1 []core.Node) (r error) {
	counter := atomic.AddUint64(&m.SetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.SetActiveNodesCounter, 1)

	if len(m.SetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStorageMock.SetActiveNodes. %v %v", p, p1)
			return
		}

		input := m.SetActiveNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStorageMockSetActiveNodesInput{p, p1}, "NodeStorage.SetActiveNodes got unexpected parameters")

		result := m.SetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStorageMock.SetActiveNodes")
			return
		}

		r = result.r

		return
	}

	if m.SetActiveNodesMock.mainExpectation != nil {

		input := m.SetActiveNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStorageMockSetActiveNodesInput{p, p1}, "NodeStorage.SetActiveNodes got unexpected parameters")
		}

		result := m.SetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStorageMock.SetActiveNodes")
		}

		r = result.r

		return
	}

	if m.SetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStorageMock.SetActiveNodes. %v %v", p, p1)
		return
	}

	return m.SetActiveNodesFunc(p, p1)
}

//SetActiveNodesMinimockCounter returns a count of NodeStorageMock.SetActiveNodesFunc invocations
func (m *NodeStorageMock) SetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetActiveNodesCounter)
}

//SetActiveNodesMinimockPreCounter returns the value of NodeStorageMock.SetActiveNodes invocations
func (m *NodeStorageMock) SetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetActiveNodesPreCounter)
}

//SetActiveNodesFinished returns true if mock invocations count is ok
func (m *NodeStorageMock) SetActiveNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetActiveNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetActiveNodesCounter) == uint64(len(m.SetActiveNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetActiveNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetActiveNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetActiveNodesFunc != nil {
		return atomic.LoadUint64(&m.SetActiveNodesCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeStorageMock) ValidateCallCounters() {

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.GetActiveNodes")
	}

	if !m.GetActiveNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.GetActiveNodesByRole")
	}

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.RemoveActiveNodesUntil")
	}

	if !m.SetActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.SetActiveNodes")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeStorageMock) MinimockFinish() {

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.GetActiveNodes")
	}

	if !m.GetActiveNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.GetActiveNodesByRole")
	}

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.RemoveActiveNodesUntil")
	}

	if !m.SetActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeStorageMock.SetActiveNodes")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetActiveNodesFinished()
		ok = ok && m.GetActiveNodesByRoleFinished()
		ok = ok && m.RemoveActiveNodesUntilFinished()
		ok = ok && m.SetActiveNodesFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetActiveNodesFinished() {
				m.t.Error("Expected call to NodeStorageMock.GetActiveNodes")
			}

			if !m.GetActiveNodesByRoleFinished() {
				m.t.Error("Expected call to NodeStorageMock.GetActiveNodesByRole")
			}

			if !m.RemoveActiveNodesUntilFinished() {
				m.t.Error("Expected call to NodeStorageMock.RemoveActiveNodesUntil")
			}

			if !m.SetActiveNodesFinished() {
				m.t.Error("Expected call to NodeStorageMock.SetActiveNodes")
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
func (m *NodeStorageMock) AllMocksCalled() bool {

	if !m.GetActiveNodesFinished() {
		return false
	}

	if !m.GetActiveNodesByRoleFinished() {
		return false
	}

	if !m.RemoveActiveNodesUntilFinished() {
		return false
	}

	if !m.SetActiveNodesFinished() {
		return false
	}

	return true
}

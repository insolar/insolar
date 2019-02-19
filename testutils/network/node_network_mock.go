package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeNetwork" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeNetworkMock implements github.com/insolar/insolar/core.NodeNetwork
type NodeNetworkMock struct {
	t minimock.Tester

	GetOriginFunc       func() (r core.Node)
	GetOriginCounter    uint64
	GetOriginPreCounter uint64
	GetOriginMock       mNodeNetworkMockGetOrigin

	GetStateFunc       func() (r core.NodeNetworkState)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mNodeNetworkMockGetState

	GetWorkingNodeFunc       func(p core.RecordRef) (r core.Node)
	GetWorkingNodeCounter    uint64
	GetWorkingNodePreCounter uint64
	GetWorkingNodeMock       mNodeNetworkMockGetWorkingNode

	GetWorkingNodesFunc       func() (r []core.Node)
	GetWorkingNodesCounter    uint64
	GetWorkingNodesPreCounter uint64
	GetWorkingNodesMock       mNodeNetworkMockGetWorkingNodes

	GetWorkingNodesByRoleFunc       func(p core.DynamicRole) (r []core.RecordRef)
	GetWorkingNodesByRoleCounter    uint64
	GetWorkingNodesByRolePreCounter uint64
	GetWorkingNodesByRoleMock       mNodeNetworkMockGetWorkingNodesByRole
}

//NewNodeNetworkMock returns a mock for github.com/insolar/insolar/core.NodeNetwork
func NewNodeNetworkMock(t minimock.Tester) *NodeNetworkMock {
	m := &NodeNetworkMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetOriginMock = mNodeNetworkMockGetOrigin{mock: m}
	m.GetStateMock = mNodeNetworkMockGetState{mock: m}
	m.GetWorkingNodeMock = mNodeNetworkMockGetWorkingNode{mock: m}
	m.GetWorkingNodesMock = mNodeNetworkMockGetWorkingNodes{mock: m}
	m.GetWorkingNodesByRoleMock = mNodeNetworkMockGetWorkingNodesByRole{mock: m}

	return m
}

type mNodeNetworkMockGetOrigin struct {
	mock              *NodeNetworkMock
	mainExpectation   *NodeNetworkMockGetOriginExpectation
	expectationSeries []*NodeNetworkMockGetOriginExpectation
}

type NodeNetworkMockGetOriginExpectation struct {
	result *NodeNetworkMockGetOriginResult
}

type NodeNetworkMockGetOriginResult struct {
	r core.Node
}

//Expect specifies that invocation of NodeNetwork.GetOrigin is expected from 1 to Infinity times
func (m *mNodeNetworkMockGetOrigin) Expect() *mNodeNetworkMockGetOrigin {
	m.mock.GetOriginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetOriginExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeNetwork.GetOrigin
func (m *mNodeNetworkMockGetOrigin) Return(r core.Node) *NodeNetworkMock {
	m.mock.GetOriginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetOriginExpectation{}
	}
	m.mainExpectation.result = &NodeNetworkMockGetOriginResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeNetwork.GetOrigin is expected once
func (m *mNodeNetworkMockGetOrigin) ExpectOnce() *NodeNetworkMockGetOriginExpectation {
	m.mock.GetOriginFunc = nil
	m.mainExpectation = nil

	expectation := &NodeNetworkMockGetOriginExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeNetworkMockGetOriginExpectation) Return(r core.Node) {
	e.result = &NodeNetworkMockGetOriginResult{r}
}

//Set uses given function f as a mock of NodeNetwork.GetOrigin method
func (m *mNodeNetworkMockGetOrigin) Set(f func() (r core.Node)) *NodeNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOriginFunc = f
	return m.mock
}

//GetOrigin implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetOrigin() (r core.Node) {
	counter := atomic.AddUint64(&m.GetOriginPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginCounter, 1)

	if len(m.GetOriginMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOriginMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeNetworkMock.GetOrigin.")
			return
		}

		result := m.GetOriginMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetOrigin")
			return
		}

		r = result.r

		return
	}

	if m.GetOriginMock.mainExpectation != nil {

		result := m.GetOriginMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetOrigin")
		}

		r = result.r

		return
	}

	if m.GetOriginFunc == nil {
		m.t.Fatalf("Unexpected call to NodeNetworkMock.GetOrigin.")
		return
	}

	return m.GetOriginFunc()
}

//GetOriginMinimockCounter returns a count of NodeNetworkMock.GetOriginFunc invocations
func (m *NodeNetworkMock) GetOriginMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginCounter)
}

//GetOriginMinimockPreCounter returns the value of NodeNetworkMock.GetOrigin invocations
func (m *NodeNetworkMock) GetOriginMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginPreCounter)
}

//GetOriginFinished returns true if mock invocations count is ok
func (m *NodeNetworkMock) GetOriginFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOriginMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOriginCounter) == uint64(len(m.GetOriginMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOriginMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOriginCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOriginFunc != nil {
		return atomic.LoadUint64(&m.GetOriginCounter) > 0
	}

	return true
}

type mNodeNetworkMockGetState struct {
	mock              *NodeNetworkMock
	mainExpectation   *NodeNetworkMockGetStateExpectation
	expectationSeries []*NodeNetworkMockGetStateExpectation
}

type NodeNetworkMockGetStateExpectation struct {
	result *NodeNetworkMockGetStateResult
}

type NodeNetworkMockGetStateResult struct {
	r core.NodeNetworkState
}

//Expect specifies that invocation of NodeNetwork.GetState is expected from 1 to Infinity times
func (m *mNodeNetworkMockGetState) Expect() *mNodeNetworkMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeNetwork.GetState
func (m *mNodeNetworkMockGetState) Return(r core.NodeNetworkState) *NodeNetworkMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetStateExpectation{}
	}
	m.mainExpectation.result = &NodeNetworkMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeNetwork.GetState is expected once
func (m *mNodeNetworkMockGetState) ExpectOnce() *NodeNetworkMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &NodeNetworkMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeNetworkMockGetStateExpectation) Return(r core.NodeNetworkState) {
	e.result = &NodeNetworkMockGetStateResult{r}
}

//Set uses given function f as a mock of NodeNetwork.GetState method
func (m *mNodeNetworkMockGetState) Set(f func() (r core.NodeNetworkState)) *NodeNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetState() (r core.NodeNetworkState) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeNetworkMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to NodeNetworkMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of NodeNetworkMock.GetStateFunc invocations
func (m *NodeNetworkMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of NodeNetworkMock.GetState invocations
func (m *NodeNetworkMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *NodeNetworkMock) GetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStateCounter) == uint64(len(m.GetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStateFunc != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	return true
}

type mNodeNetworkMockGetWorkingNode struct {
	mock              *NodeNetworkMock
	mainExpectation   *NodeNetworkMockGetWorkingNodeExpectation
	expectationSeries []*NodeNetworkMockGetWorkingNodeExpectation
}

type NodeNetworkMockGetWorkingNodeExpectation struct {
	input  *NodeNetworkMockGetWorkingNodeInput
	result *NodeNetworkMockGetWorkingNodeResult
}

type NodeNetworkMockGetWorkingNodeInput struct {
	p core.RecordRef
}

type NodeNetworkMockGetWorkingNodeResult struct {
	r core.Node
}

//Expect specifies that invocation of NodeNetwork.GetWorkingNode is expected from 1 to Infinity times
func (m *mNodeNetworkMockGetWorkingNode) Expect(p core.RecordRef) *mNodeNetworkMockGetWorkingNode {
	m.mock.GetWorkingNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetWorkingNodeExpectation{}
	}
	m.mainExpectation.input = &NodeNetworkMockGetWorkingNodeInput{p}
	return m
}

//Return specifies results of invocation of NodeNetwork.GetWorkingNode
func (m *mNodeNetworkMockGetWorkingNode) Return(r core.Node) *NodeNetworkMock {
	m.mock.GetWorkingNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetWorkingNodeExpectation{}
	}
	m.mainExpectation.result = &NodeNetworkMockGetWorkingNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeNetwork.GetWorkingNode is expected once
func (m *mNodeNetworkMockGetWorkingNode) ExpectOnce(p core.RecordRef) *NodeNetworkMockGetWorkingNodeExpectation {
	m.mock.GetWorkingNodeFunc = nil
	m.mainExpectation = nil

	expectation := &NodeNetworkMockGetWorkingNodeExpectation{}
	expectation.input = &NodeNetworkMockGetWorkingNodeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeNetworkMockGetWorkingNodeExpectation) Return(r core.Node) {
	e.result = &NodeNetworkMockGetWorkingNodeResult{r}
}

//Set uses given function f as a mock of NodeNetwork.GetWorkingNode method
func (m *mNodeNetworkMockGetWorkingNode) Set(f func(p core.RecordRef) (r core.Node)) *NodeNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodeFunc = f
	return m.mock
}

//GetWorkingNode implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetWorkingNode(p core.RecordRef) (r core.Node) {
	counter := atomic.AddUint64(&m.GetWorkingNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodeCounter, 1)

	if len(m.GetWorkingNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeNetworkMock.GetWorkingNode. %v", p)
			return
		}

		input := m.GetWorkingNodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeNetworkMockGetWorkingNodeInput{p}, "NodeNetwork.GetWorkingNode got unexpected parameters")

		result := m.GetWorkingNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetWorkingNode")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodeMock.mainExpectation != nil {

		input := m.GetWorkingNodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeNetworkMockGetWorkingNodeInput{p}, "NodeNetwork.GetWorkingNode got unexpected parameters")
		}

		result := m.GetWorkingNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetWorkingNode")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodeFunc == nil {
		m.t.Fatalf("Unexpected call to NodeNetworkMock.GetWorkingNode. %v", p)
		return
	}

	return m.GetWorkingNodeFunc(p)
}

//GetWorkingNodeMinimockCounter returns a count of NodeNetworkMock.GetWorkingNodeFunc invocations
func (m *NodeNetworkMock) GetWorkingNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodeCounter)
}

//GetWorkingNodeMinimockPreCounter returns the value of NodeNetworkMock.GetWorkingNode invocations
func (m *NodeNetworkMock) GetWorkingNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodePreCounter)
}

//GetWorkingNodeFinished returns true if mock invocations count is ok
func (m *NodeNetworkMock) GetWorkingNodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetWorkingNodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetWorkingNodeCounter) == uint64(len(m.GetWorkingNodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetWorkingNodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetWorkingNodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetWorkingNodeFunc != nil {
		return atomic.LoadUint64(&m.GetWorkingNodeCounter) > 0
	}

	return true
}

type mNodeNetworkMockGetWorkingNodes struct {
	mock              *NodeNetworkMock
	mainExpectation   *NodeNetworkMockGetWorkingNodesExpectation
	expectationSeries []*NodeNetworkMockGetWorkingNodesExpectation
}

type NodeNetworkMockGetWorkingNodesExpectation struct {
	result *NodeNetworkMockGetWorkingNodesResult
}

type NodeNetworkMockGetWorkingNodesResult struct {
	r []core.Node
}

//Expect specifies that invocation of NodeNetwork.GetWorkingNodes is expected from 1 to Infinity times
func (m *mNodeNetworkMockGetWorkingNodes) Expect() *mNodeNetworkMockGetWorkingNodes {
	m.mock.GetWorkingNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetWorkingNodesExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeNetwork.GetWorkingNodes
func (m *mNodeNetworkMockGetWorkingNodes) Return(r []core.Node) *NodeNetworkMock {
	m.mock.GetWorkingNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetWorkingNodesExpectation{}
	}
	m.mainExpectation.result = &NodeNetworkMockGetWorkingNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeNetwork.GetWorkingNodes is expected once
func (m *mNodeNetworkMockGetWorkingNodes) ExpectOnce() *NodeNetworkMockGetWorkingNodesExpectation {
	m.mock.GetWorkingNodesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeNetworkMockGetWorkingNodesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeNetworkMockGetWorkingNodesExpectation) Return(r []core.Node) {
	e.result = &NodeNetworkMockGetWorkingNodesResult{r}
}

//Set uses given function f as a mock of NodeNetwork.GetWorkingNodes method
func (m *mNodeNetworkMockGetWorkingNodes) Set(f func() (r []core.Node)) *NodeNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodesFunc = f
	return m.mock
}

//GetWorkingNodes implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetWorkingNodes() (r []core.Node) {
	counter := atomic.AddUint64(&m.GetWorkingNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodesCounter, 1)

	if len(m.GetWorkingNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeNetworkMock.GetWorkingNodes.")
			return
		}

		result := m.GetWorkingNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetWorkingNodes")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesMock.mainExpectation != nil {

		result := m.GetWorkingNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetWorkingNodes")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeNetworkMock.GetWorkingNodes.")
		return
	}

	return m.GetWorkingNodesFunc()
}

//GetWorkingNodesMinimockCounter returns a count of NodeNetworkMock.GetWorkingNodesFunc invocations
func (m *NodeNetworkMock) GetWorkingNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesCounter)
}

//GetWorkingNodesMinimockPreCounter returns the value of NodeNetworkMock.GetWorkingNodes invocations
func (m *NodeNetworkMock) GetWorkingNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesPreCounter)
}

//GetWorkingNodesFinished returns true if mock invocations count is ok
func (m *NodeNetworkMock) GetWorkingNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetWorkingNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetWorkingNodesCounter) == uint64(len(m.GetWorkingNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetWorkingNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetWorkingNodesFunc != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesCounter) > 0
	}

	return true
}

type mNodeNetworkMockGetWorkingNodesByRole struct {
	mock              *NodeNetworkMock
	mainExpectation   *NodeNetworkMockGetWorkingNodesByRoleExpectation
	expectationSeries []*NodeNetworkMockGetWorkingNodesByRoleExpectation
}

type NodeNetworkMockGetWorkingNodesByRoleExpectation struct {
	input  *NodeNetworkMockGetWorkingNodesByRoleInput
	result *NodeNetworkMockGetWorkingNodesByRoleResult
}

type NodeNetworkMockGetWorkingNodesByRoleInput struct {
	p core.DynamicRole
}

type NodeNetworkMockGetWorkingNodesByRoleResult struct {
	r []core.RecordRef
}

//Expect specifies that invocation of NodeNetwork.GetWorkingNodesByRole is expected from 1 to Infinity times
func (m *mNodeNetworkMockGetWorkingNodesByRole) Expect(p core.DynamicRole) *mNodeNetworkMockGetWorkingNodesByRole {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetWorkingNodesByRoleExpectation{}
	}
	m.mainExpectation.input = &NodeNetworkMockGetWorkingNodesByRoleInput{p}
	return m
}

//Return specifies results of invocation of NodeNetwork.GetWorkingNodesByRole
func (m *mNodeNetworkMockGetWorkingNodesByRole) Return(r []core.RecordRef) *NodeNetworkMock {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeNetworkMockGetWorkingNodesByRoleExpectation{}
	}
	m.mainExpectation.result = &NodeNetworkMockGetWorkingNodesByRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeNetwork.GetWorkingNodesByRole is expected once
func (m *mNodeNetworkMockGetWorkingNodesByRole) ExpectOnce(p core.DynamicRole) *NodeNetworkMockGetWorkingNodesByRoleExpectation {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.mainExpectation = nil

	expectation := &NodeNetworkMockGetWorkingNodesByRoleExpectation{}
	expectation.input = &NodeNetworkMockGetWorkingNodesByRoleInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeNetworkMockGetWorkingNodesByRoleExpectation) Return(r []core.RecordRef) {
	e.result = &NodeNetworkMockGetWorkingNodesByRoleResult{r}
}

//Set uses given function f as a mock of NodeNetwork.GetWorkingNodesByRole method
func (m *mNodeNetworkMockGetWorkingNodesByRole) Set(f func(p core.DynamicRole) (r []core.RecordRef)) *NodeNetworkMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodesByRoleFunc = f
	return m.mock
}

//GetWorkingNodesByRole implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetWorkingNodesByRole(p core.DynamicRole) (r []core.RecordRef) {
	counter := atomic.AddUint64(&m.GetWorkingNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodesByRoleCounter, 1)

	if len(m.GetWorkingNodesByRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodesByRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeNetworkMock.GetWorkingNodesByRole. %v", p)
			return
		}

		input := m.GetWorkingNodesByRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeNetworkMockGetWorkingNodesByRoleInput{p}, "NodeNetwork.GetWorkingNodesByRole got unexpected parameters")

		result := m.GetWorkingNodesByRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetWorkingNodesByRole")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesByRoleMock.mainExpectation != nil {

		input := m.GetWorkingNodesByRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeNetworkMockGetWorkingNodesByRoleInput{p}, "NodeNetwork.GetWorkingNodesByRole got unexpected parameters")
		}

		result := m.GetWorkingNodesByRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeNetworkMock.GetWorkingNodesByRole")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesByRoleFunc == nil {
		m.t.Fatalf("Unexpected call to NodeNetworkMock.GetWorkingNodesByRole. %v", p)
		return
	}

	return m.GetWorkingNodesByRoleFunc(p)
}

//GetWorkingNodesByRoleMinimockCounter returns a count of NodeNetworkMock.GetWorkingNodesByRoleFunc invocations
func (m *NodeNetworkMock) GetWorkingNodesByRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter)
}

//GetWorkingNodesByRoleMinimockPreCounter returns the value of NodeNetworkMock.GetWorkingNodesByRole invocations
func (m *NodeNetworkMock) GetWorkingNodesByRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesByRolePreCounter)
}

//GetWorkingNodesByRoleFinished returns true if mock invocations count is ok
func (m *NodeNetworkMock) GetWorkingNodesByRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetWorkingNodesByRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter) == uint64(len(m.GetWorkingNodesByRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetWorkingNodesByRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetWorkingNodesByRoleFunc != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeNetworkMock) ValidateCallCounters() {

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetOrigin")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetState")
	}

	if !m.GetWorkingNodeFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetWorkingNode")
	}

	if !m.GetWorkingNodesFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetWorkingNodes")
	}

	if !m.GetWorkingNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetWorkingNodesByRole")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeNetworkMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeNetworkMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeNetworkMock) MinimockFinish() {

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetOrigin")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetState")
	}

	if !m.GetWorkingNodeFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetWorkingNode")
	}

	if !m.GetWorkingNodesFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetWorkingNodes")
	}

	if !m.GetWorkingNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeNetworkMock.GetWorkingNodesByRole")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeNetworkMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeNetworkMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetOriginFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.GetWorkingNodeFinished()
		ok = ok && m.GetWorkingNodesFinished()
		ok = ok && m.GetWorkingNodesByRoleFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetOriginFinished() {
				m.t.Error("Expected call to NodeNetworkMock.GetOrigin")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to NodeNetworkMock.GetState")
			}

			if !m.GetWorkingNodeFinished() {
				m.t.Error("Expected call to NodeNetworkMock.GetWorkingNode")
			}

			if !m.GetWorkingNodesFinished() {
				m.t.Error("Expected call to NodeNetworkMock.GetWorkingNodes")
			}

			if !m.GetWorkingNodesByRoleFinished() {
				m.t.Error("Expected call to NodeNetworkMock.GetWorkingNodesByRole")
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
func (m *NodeNetworkMock) AllMocksCalled() bool {

	if !m.GetOriginFinished() {
		return false
	}

	if !m.GetStateFinished() {
		return false
	}

	if !m.GetWorkingNodeFinished() {
		return false
	}

	if !m.GetWorkingNodesFinished() {
		return false
	}

	if !m.GetWorkingNodesByRoleFinished() {
		return false
	}

	return true
}

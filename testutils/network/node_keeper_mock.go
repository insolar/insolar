package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeKeeper" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	packets "github.com/insolar/insolar/consensus/packets"
	core "github.com/insolar/insolar/core"
	network "github.com/insolar/insolar/network"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeKeeperMock implements github.com/insolar/insolar/network.NodeKeeper
type NodeKeeperMock struct {
	t minimock.Tester

	AddActiveNodesFunc       func(p []core.Node)
	AddActiveNodesCounter    uint64
	AddActiveNodesPreCounter uint64
	AddActiveNodesMock       mNodeKeeperMockAddActiveNodes

	AddPendingClaimFunc       func(p packets.ReferendumClaim) (r bool)
	AddPendingClaimCounter    uint64
	AddPendingClaimPreCounter uint64
	AddPendingClaimMock       mNodeKeeperMockAddPendingClaim

	GetActiveNodeFunc       func(p core.RecordRef) (r core.Node)
	GetActiveNodeCounter    uint64
	GetActiveNodePreCounter uint64
	GetActiveNodeMock       mNodeKeeperMockGetActiveNode

	GetActiveNodeByShortIDFunc       func(p core.ShortNodeID) (r core.Node)
	GetActiveNodeByShortIDCounter    uint64
	GetActiveNodeByShortIDPreCounter uint64
	GetActiveNodeByShortIDMock       mNodeKeeperMockGetActiveNodeByShortID

	GetActiveNodesFunc       func() (r []core.Node)
	GetActiveNodesCounter    uint64
	GetActiveNodesPreCounter uint64
	GetActiveNodesMock       mNodeKeeperMockGetActiveNodes

	GetActiveNodesByRoleFunc       func(p core.DynamicRole) (r []core.RecordRef)
	GetActiveNodesByRoleCounter    uint64
	GetActiveNodesByRolePreCounter uint64
	GetActiveNodesByRoleMock       mNodeKeeperMockGetActiveNodesByRole

	GetClaimQueueFunc       func() (r network.ClaimQueue)
	GetClaimQueueCounter    uint64
	GetClaimQueuePreCounter uint64
	GetClaimQueueMock       mNodeKeeperMockGetClaimQueue

	GetCloudHashFunc       func() (r []byte)
	GetCloudHashCounter    uint64
	GetCloudHashPreCounter uint64
	GetCloudHashMock       mNodeKeeperMockGetCloudHash

	GetOriginFunc       func() (r core.Node)
	GetOriginCounter    uint64
	GetOriginPreCounter uint64
	GetOriginMock       mNodeKeeperMockGetOrigin

	GetOriginClaimFunc       func() (r *packets.NodeJoinClaim, r1 error)
	GetOriginClaimCounter    uint64
	GetOriginClaimPreCounter uint64
	GetOriginClaimMock       mNodeKeeperMockGetOriginClaim

	GetSparseUnsyncListFunc       func(p int) (r network.UnsyncList)
	GetSparseUnsyncListCounter    uint64
	GetSparseUnsyncListPreCounter uint64
	GetSparseUnsyncListMock       mNodeKeeperMockGetSparseUnsyncList

	GetStateFunc       func() (r network.NodeKeeperState)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mNodeKeeperMockGetState

	GetUnsyncListFunc       func() (r network.UnsyncList)
	GetUnsyncListCounter    uint64
	GetUnsyncListPreCounter uint64
	GetUnsyncListMock       mNodeKeeperMockGetUnsyncList

	IsBootstrappedFunc       func() (r bool)
	IsBootstrappedCounter    uint64
	IsBootstrappedPreCounter uint64
	IsBootstrappedMock       mNodeKeeperMockIsBootstrapped

	MoveSyncToActiveFunc       func()
	MoveSyncToActiveCounter    uint64
	MoveSyncToActivePreCounter uint64
	MoveSyncToActiveMock       mNodeKeeperMockMoveSyncToActive

	NodesJoinedDuringPreviousPulseFunc       func() (r bool)
	NodesJoinedDuringPreviousPulseCounter    uint64
	NodesJoinedDuringPreviousPulsePreCounter uint64
	NodesJoinedDuringPreviousPulseMock       mNodeKeeperMockNodesJoinedDuringPreviousPulse

	SetCloudHashFunc       func(p []byte)
	SetCloudHashCounter    uint64
	SetCloudHashPreCounter uint64
	SetCloudHashMock       mNodeKeeperMockSetCloudHash

	SetIsBootstrappedFunc       func(p bool)
	SetIsBootstrappedCounter    uint64
	SetIsBootstrappedPreCounter uint64
	SetIsBootstrappedMock       mNodeKeeperMockSetIsBootstrapped

	SetStateFunc       func(p network.NodeKeeperState)
	SetStateCounter    uint64
	SetStatePreCounter uint64
	SetStateMock       mNodeKeeperMockSetState

	SyncFunc       func(p network.UnsyncList)
	SyncCounter    uint64
	SyncPreCounter uint64
	SyncMock       mNodeKeeperMockSync
}

//NewNodeKeeperMock returns a mock for github.com/insolar/insolar/network.NodeKeeper
func NewNodeKeeperMock(t minimock.Tester) *NodeKeeperMock {
	m := &NodeKeeperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddActiveNodesMock = mNodeKeeperMockAddActiveNodes{mock: m}
	m.AddPendingClaimMock = mNodeKeeperMockAddPendingClaim{mock: m}
	m.GetActiveNodeMock = mNodeKeeperMockGetActiveNode{mock: m}
	m.GetActiveNodeByShortIDMock = mNodeKeeperMockGetActiveNodeByShortID{mock: m}
	m.GetActiveNodesMock = mNodeKeeperMockGetActiveNodes{mock: m}
	m.GetActiveNodesByRoleMock = mNodeKeeperMockGetActiveNodesByRole{mock: m}
	m.GetClaimQueueMock = mNodeKeeperMockGetClaimQueue{mock: m}
	m.GetCloudHashMock = mNodeKeeperMockGetCloudHash{mock: m}
	m.GetOriginMock = mNodeKeeperMockGetOrigin{mock: m}
	m.GetOriginClaimMock = mNodeKeeperMockGetOriginClaim{mock: m}
	m.GetSparseUnsyncListMock = mNodeKeeperMockGetSparseUnsyncList{mock: m}
	m.GetStateMock = mNodeKeeperMockGetState{mock: m}
	m.GetUnsyncListMock = mNodeKeeperMockGetUnsyncList{mock: m}
	m.IsBootstrappedMock = mNodeKeeperMockIsBootstrapped{mock: m}
	m.MoveSyncToActiveMock = mNodeKeeperMockMoveSyncToActive{mock: m}
	m.NodesJoinedDuringPreviousPulseMock = mNodeKeeperMockNodesJoinedDuringPreviousPulse{mock: m}
	m.SetCloudHashMock = mNodeKeeperMockSetCloudHash{mock: m}
	m.SetIsBootstrappedMock = mNodeKeeperMockSetIsBootstrapped{mock: m}
	m.SetStateMock = mNodeKeeperMockSetState{mock: m}
	m.SyncMock = mNodeKeeperMockSync{mock: m}

	return m
}

type mNodeKeeperMockAddActiveNodes struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockAddActiveNodesExpectation
	expectationSeries []*NodeKeeperMockAddActiveNodesExpectation
}

type NodeKeeperMockAddActiveNodesExpectation struct {
	input *NodeKeeperMockAddActiveNodesInput
}

type NodeKeeperMockAddActiveNodesInput struct {
	p []core.Node
}

//Expect specifies that invocation of NodeKeeper.AddActiveNodes is expected from 1 to Infinity times
func (m *mNodeKeeperMockAddActiveNodes) Expect(p []core.Node) *mNodeKeeperMockAddActiveNodes {
	m.mock.AddActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockAddActiveNodesExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockAddActiveNodesInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.AddActiveNodes
func (m *mNodeKeeperMockAddActiveNodes) Return() *NodeKeeperMock {
	m.mock.AddActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockAddActiveNodesExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.AddActiveNodes is expected once
func (m *mNodeKeeperMockAddActiveNodes) ExpectOnce(p []core.Node) *NodeKeeperMockAddActiveNodesExpectation {
	m.mock.AddActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockAddActiveNodesExpectation{}
	expectation.input = &NodeKeeperMockAddActiveNodesInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeKeeper.AddActiveNodes method
func (m *mNodeKeeperMockAddActiveNodes) Set(f func(p []core.Node)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddActiveNodesFunc = f
	return m.mock
}

//AddActiveNodes implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) AddActiveNodes(p []core.Node) {
	counter := atomic.AddUint64(&m.AddActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.AddActiveNodesCounter, 1)

	if len(m.AddActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.AddActiveNodes. %v", p)
			return
		}

		input := m.AddActiveNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockAddActiveNodesInput{p}, "NodeKeeper.AddActiveNodes got unexpected parameters")

		return
	}

	if m.AddActiveNodesMock.mainExpectation != nil {

		input := m.AddActiveNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockAddActiveNodesInput{p}, "NodeKeeper.AddActiveNodes got unexpected parameters")
		}

		return
	}

	if m.AddActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.AddActiveNodes. %v", p)
		return
	}

	m.AddActiveNodesFunc(p)
}

//AddActiveNodesMinimockCounter returns a count of NodeKeeperMock.AddActiveNodesFunc invocations
func (m *NodeKeeperMock) AddActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddActiveNodesCounter)
}

//AddActiveNodesMinimockPreCounter returns the value of NodeKeeperMock.AddActiveNodes invocations
func (m *NodeKeeperMock) AddActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddActiveNodesPreCounter)
}

//AddActiveNodesFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) AddActiveNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddActiveNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddActiveNodesCounter) == uint64(len(m.AddActiveNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddActiveNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddActiveNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddActiveNodesFunc != nil {
		return atomic.LoadUint64(&m.AddActiveNodesCounter) > 0
	}

	return true
}

type mNodeKeeperMockAddPendingClaim struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockAddPendingClaimExpectation
	expectationSeries []*NodeKeeperMockAddPendingClaimExpectation
}

type NodeKeeperMockAddPendingClaimExpectation struct {
	input  *NodeKeeperMockAddPendingClaimInput
	result *NodeKeeperMockAddPendingClaimResult
}

type NodeKeeperMockAddPendingClaimInput struct {
	p packets.ReferendumClaim
}

type NodeKeeperMockAddPendingClaimResult struct {
	r bool
}

//Expect specifies that invocation of NodeKeeper.AddPendingClaim is expected from 1 to Infinity times
func (m *mNodeKeeperMockAddPendingClaim) Expect(p packets.ReferendumClaim) *mNodeKeeperMockAddPendingClaim {
	m.mock.AddPendingClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockAddPendingClaimExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockAddPendingClaimInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.AddPendingClaim
func (m *mNodeKeeperMockAddPendingClaim) Return(r bool) *NodeKeeperMock {
	m.mock.AddPendingClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockAddPendingClaimExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockAddPendingClaimResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.AddPendingClaim is expected once
func (m *mNodeKeeperMockAddPendingClaim) ExpectOnce(p packets.ReferendumClaim) *NodeKeeperMockAddPendingClaimExpectation {
	m.mock.AddPendingClaimFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockAddPendingClaimExpectation{}
	expectation.input = &NodeKeeperMockAddPendingClaimInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockAddPendingClaimExpectation) Return(r bool) {
	e.result = &NodeKeeperMockAddPendingClaimResult{r}
}

//Set uses given function f as a mock of NodeKeeper.AddPendingClaim method
func (m *mNodeKeeperMockAddPendingClaim) Set(f func(p packets.ReferendumClaim) (r bool)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddPendingClaimFunc = f
	return m.mock
}

//AddPendingClaim implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) AddPendingClaim(p packets.ReferendumClaim) (r bool) {
	counter := atomic.AddUint64(&m.AddPendingClaimPreCounter, 1)
	defer atomic.AddUint64(&m.AddPendingClaimCounter, 1)

	if len(m.AddPendingClaimMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddPendingClaimMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.AddPendingClaim. %v", p)
			return
		}

		input := m.AddPendingClaimMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockAddPendingClaimInput{p}, "NodeKeeper.AddPendingClaim got unexpected parameters")

		result := m.AddPendingClaimMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.AddPendingClaim")
			return
		}

		r = result.r

		return
	}

	if m.AddPendingClaimMock.mainExpectation != nil {

		input := m.AddPendingClaimMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockAddPendingClaimInput{p}, "NodeKeeper.AddPendingClaim got unexpected parameters")
		}

		result := m.AddPendingClaimMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.AddPendingClaim")
		}

		r = result.r

		return
	}

	if m.AddPendingClaimFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.AddPendingClaim. %v", p)
		return
	}

	return m.AddPendingClaimFunc(p)
}

//AddPendingClaimMinimockCounter returns a count of NodeKeeperMock.AddPendingClaimFunc invocations
func (m *NodeKeeperMock) AddPendingClaimMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddPendingClaimCounter)
}

//AddPendingClaimMinimockPreCounter returns the value of NodeKeeperMock.AddPendingClaim invocations
func (m *NodeKeeperMock) AddPendingClaimMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddPendingClaimPreCounter)
}

//AddPendingClaimFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) AddPendingClaimFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddPendingClaimMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddPendingClaimCounter) == uint64(len(m.AddPendingClaimMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddPendingClaimMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddPendingClaimCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddPendingClaimFunc != nil {
		return atomic.LoadUint64(&m.AddPendingClaimCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetActiveNode struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetActiveNodeExpectation
	expectationSeries []*NodeKeeperMockGetActiveNodeExpectation
}

type NodeKeeperMockGetActiveNodeExpectation struct {
	input  *NodeKeeperMockGetActiveNodeInput
	result *NodeKeeperMockGetActiveNodeResult
}

type NodeKeeperMockGetActiveNodeInput struct {
	p core.RecordRef
}

type NodeKeeperMockGetActiveNodeResult struct {
	r core.Node
}

//Expect specifies that invocation of NodeKeeper.GetActiveNode is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetActiveNode) Expect(p core.RecordRef) *mNodeKeeperMockGetActiveNode {
	m.mock.GetActiveNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodeExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockGetActiveNodeInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.GetActiveNode
func (m *mNodeKeeperMockGetActiveNode) Return(r core.Node) *NodeKeeperMock {
	m.mock.GetActiveNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodeExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetActiveNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetActiveNode is expected once
func (m *mNodeKeeperMockGetActiveNode) ExpectOnce(p core.RecordRef) *NodeKeeperMockGetActiveNodeExpectation {
	m.mock.GetActiveNodeFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetActiveNodeExpectation{}
	expectation.input = &NodeKeeperMockGetActiveNodeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetActiveNodeExpectation) Return(r core.Node) {
	e.result = &NodeKeeperMockGetActiveNodeResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNode method
func (m *mNodeKeeperMockGetActiveNode) Set(f func(p core.RecordRef) (r core.Node)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodeFunc = f
	return m.mock
}

//GetActiveNode implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNode(p core.RecordRef) (r core.Node) {
	counter := atomic.AddUint64(&m.GetActiveNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeCounter, 1)

	if len(m.GetActiveNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNode. %v", p)
			return
		}

		input := m.GetActiveNodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockGetActiveNodeInput{p}, "NodeKeeper.GetActiveNode got unexpected parameters")

		result := m.GetActiveNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNode")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodeMock.mainExpectation != nil {

		input := m.GetActiveNodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockGetActiveNodeInput{p}, "NodeKeeper.GetActiveNode got unexpected parameters")
		}

		result := m.GetActiveNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNode")
		}

		r = result.r

		return
	}

	if m.GetActiveNodeFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNode. %v", p)
		return
	}

	return m.GetActiveNodeFunc(p)
}

//GetActiveNodeMinimockCounter returns a count of NodeKeeperMock.GetActiveNodeFunc invocations
func (m *NodeKeeperMock) GetActiveNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeCounter)
}

//GetActiveNodeMinimockPreCounter returns the value of NodeKeeperMock.GetActiveNode invocations
func (m *NodeKeeperMock) GetActiveNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodePreCounter)
}

//GetActiveNodeFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetActiveNodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodeCounter) == uint64(len(m.GetActiveNodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodeFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodeCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetActiveNodeByShortID struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetActiveNodeByShortIDExpectation
	expectationSeries []*NodeKeeperMockGetActiveNodeByShortIDExpectation
}

type NodeKeeperMockGetActiveNodeByShortIDExpectation struct {
	input  *NodeKeeperMockGetActiveNodeByShortIDInput
	result *NodeKeeperMockGetActiveNodeByShortIDResult
}

type NodeKeeperMockGetActiveNodeByShortIDInput struct {
	p core.ShortNodeID
}

type NodeKeeperMockGetActiveNodeByShortIDResult struct {
	r core.Node
}

//Expect specifies that invocation of NodeKeeper.GetActiveNodeByShortID is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetActiveNodeByShortID) Expect(p core.ShortNodeID) *mNodeKeeperMockGetActiveNodeByShortID {
	m.mock.GetActiveNodeByShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodeByShortIDExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockGetActiveNodeByShortIDInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.GetActiveNodeByShortID
func (m *mNodeKeeperMockGetActiveNodeByShortID) Return(r core.Node) *NodeKeeperMock {
	m.mock.GetActiveNodeByShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodeByShortIDExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetActiveNodeByShortIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetActiveNodeByShortID is expected once
func (m *mNodeKeeperMockGetActiveNodeByShortID) ExpectOnce(p core.ShortNodeID) *NodeKeeperMockGetActiveNodeByShortIDExpectation {
	m.mock.GetActiveNodeByShortIDFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetActiveNodeByShortIDExpectation{}
	expectation.input = &NodeKeeperMockGetActiveNodeByShortIDInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetActiveNodeByShortIDExpectation) Return(r core.Node) {
	e.result = &NodeKeeperMockGetActiveNodeByShortIDResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNodeByShortID method
func (m *mNodeKeeperMockGetActiveNodeByShortID) Set(f func(p core.ShortNodeID) (r core.Node)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodeByShortIDFunc = f
	return m.mock
}

//GetActiveNodeByShortID implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNodeByShortID(p core.ShortNodeID) (r core.Node) {
	counter := atomic.AddUint64(&m.GetActiveNodeByShortIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeByShortIDCounter, 1)

	if len(m.GetActiveNodeByShortIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodeByShortIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNodeByShortID. %v", p)
			return
		}

		input := m.GetActiveNodeByShortIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockGetActiveNodeByShortIDInput{p}, "NodeKeeper.GetActiveNodeByShortID got unexpected parameters")

		result := m.GetActiveNodeByShortIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodeByShortID")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodeByShortIDMock.mainExpectation != nil {

		input := m.GetActiveNodeByShortIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockGetActiveNodeByShortIDInput{p}, "NodeKeeper.GetActiveNodeByShortID got unexpected parameters")
		}

		result := m.GetActiveNodeByShortIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodeByShortID")
		}

		r = result.r

		return
	}

	if m.GetActiveNodeByShortIDFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNodeByShortID. %v", p)
		return
	}

	return m.GetActiveNodeByShortIDFunc(p)
}

//GetActiveNodeByShortIDMinimockCounter returns a count of NodeKeeperMock.GetActiveNodeByShortIDFunc invocations
func (m *NodeKeeperMock) GetActiveNodeByShortIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter)
}

//GetActiveNodeByShortIDMinimockPreCounter returns the value of NodeKeeperMock.GetActiveNodeByShortID invocations
func (m *NodeKeeperMock) GetActiveNodeByShortIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeByShortIDPreCounter)
}

//GetActiveNodeByShortIDFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetActiveNodeByShortIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodeByShortIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) == uint64(len(m.GetActiveNodeByShortIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodeByShortIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodeByShortIDFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetActiveNodes struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetActiveNodesExpectation
	expectationSeries []*NodeKeeperMockGetActiveNodesExpectation
}

type NodeKeeperMockGetActiveNodesExpectation struct {
	result *NodeKeeperMockGetActiveNodesResult
}

type NodeKeeperMockGetActiveNodesResult struct {
	r []core.Node
}

//Expect specifies that invocation of NodeKeeper.GetActiveNodes is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetActiveNodes) Expect() *mNodeKeeperMockGetActiveNodes {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodesExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetActiveNodes
func (m *mNodeKeeperMockGetActiveNodes) Return(r []core.Node) *NodeKeeperMock {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetActiveNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetActiveNodes is expected once
func (m *mNodeKeeperMockGetActiveNodes) ExpectOnce() *NodeKeeperMockGetActiveNodesExpectation {
	m.mock.GetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetActiveNodesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetActiveNodesExpectation) Return(r []core.Node) {
	e.result = &NodeKeeperMockGetActiveNodesResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNodes method
func (m *mNodeKeeperMockGetActiveNodes) Set(f func() (r []core.Node)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesFunc = f
	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNodes() (r []core.Node) {
	counter := atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNodes.")
			return
		}

		result := m.GetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodes")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodesMock.mainExpectation != nil {

		result := m.GetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodes")
		}

		r = result.r

		return
	}

	if m.GetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNodes.")
		return
	}

	return m.GetActiveNodesFunc()
}

//GetActiveNodesMinimockCounter returns a count of NodeKeeperMock.GetActiveNodesFunc invocations
func (m *NodeKeeperMock) GetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesCounter)
}

//GetActiveNodesMinimockPreCounter returns the value of NodeKeeperMock.GetActiveNodes invocations
func (m *NodeKeeperMock) GetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesPreCounter)
}

//GetActiveNodesFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetActiveNodesFinished() bool {
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

type mNodeKeeperMockGetActiveNodesByRole struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetActiveNodesByRoleExpectation
	expectationSeries []*NodeKeeperMockGetActiveNodesByRoleExpectation
}

type NodeKeeperMockGetActiveNodesByRoleExpectation struct {
	input  *NodeKeeperMockGetActiveNodesByRoleInput
	result *NodeKeeperMockGetActiveNodesByRoleResult
}

type NodeKeeperMockGetActiveNodesByRoleInput struct {
	p core.DynamicRole
}

type NodeKeeperMockGetActiveNodesByRoleResult struct {
	r []core.RecordRef
}

//Expect specifies that invocation of NodeKeeper.GetActiveNodesByRole is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetActiveNodesByRole) Expect(p core.DynamicRole) *mNodeKeeperMockGetActiveNodesByRole {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodesByRoleExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockGetActiveNodesByRoleInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.GetActiveNodesByRole
func (m *mNodeKeeperMockGetActiveNodesByRole) Return(r []core.RecordRef) *NodeKeeperMock {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetActiveNodesByRoleExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetActiveNodesByRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetActiveNodesByRole is expected once
func (m *mNodeKeeperMockGetActiveNodesByRole) ExpectOnce(p core.DynamicRole) *NodeKeeperMockGetActiveNodesByRoleExpectation {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetActiveNodesByRoleExpectation{}
	expectation.input = &NodeKeeperMockGetActiveNodesByRoleInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetActiveNodesByRoleExpectation) Return(r []core.RecordRef) {
	e.result = &NodeKeeperMockGetActiveNodesByRoleResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNodesByRole method
func (m *mNodeKeeperMockGetActiveNodesByRole) Set(f func(p core.DynamicRole) (r []core.RecordRef)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesByRoleFunc = f
	return m.mock
}

//GetActiveNodesByRole implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNodesByRole(p core.DynamicRole) (r []core.RecordRef) {
	counter := atomic.AddUint64(&m.GetActiveNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesByRoleCounter, 1)

	if len(m.GetActiveNodesByRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesByRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNodesByRole. %v", p)
			return
		}

		input := m.GetActiveNodesByRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockGetActiveNodesByRoleInput{p}, "NodeKeeper.GetActiveNodesByRole got unexpected parameters")

		result := m.GetActiveNodesByRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodesByRole")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodesByRoleMock.mainExpectation != nil {

		input := m.GetActiveNodesByRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockGetActiveNodesByRoleInput{p}, "NodeKeeper.GetActiveNodesByRole got unexpected parameters")
		}

		result := m.GetActiveNodesByRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodesByRole")
		}

		r = result.r

		return
	}

	if m.GetActiveNodesByRoleFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetActiveNodesByRole. %v", p)
		return
	}

	return m.GetActiveNodesByRoleFunc(p)
}

//GetActiveNodesByRoleMinimockCounter returns a count of NodeKeeperMock.GetActiveNodesByRoleFunc invocations
func (m *NodeKeeperMock) GetActiveNodesByRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter)
}

//GetActiveNodesByRoleMinimockPreCounter returns the value of NodeKeeperMock.GetActiveNodesByRole invocations
func (m *NodeKeeperMock) GetActiveNodesByRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRolePreCounter)
}

//GetActiveNodesByRoleFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetActiveNodesByRoleFinished() bool {
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

type mNodeKeeperMockGetClaimQueue struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetClaimQueueExpectation
	expectationSeries []*NodeKeeperMockGetClaimQueueExpectation
}

type NodeKeeperMockGetClaimQueueExpectation struct {
	result *NodeKeeperMockGetClaimQueueResult
}

type NodeKeeperMockGetClaimQueueResult struct {
	r network.ClaimQueue
}

//Expect specifies that invocation of NodeKeeper.GetClaimQueue is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetClaimQueue) Expect() *mNodeKeeperMockGetClaimQueue {
	m.mock.GetClaimQueueFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetClaimQueueExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetClaimQueue
func (m *mNodeKeeperMockGetClaimQueue) Return(r network.ClaimQueue) *NodeKeeperMock {
	m.mock.GetClaimQueueFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetClaimQueueExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetClaimQueueResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetClaimQueue is expected once
func (m *mNodeKeeperMockGetClaimQueue) ExpectOnce() *NodeKeeperMockGetClaimQueueExpectation {
	m.mock.GetClaimQueueFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetClaimQueueExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetClaimQueueExpectation) Return(r network.ClaimQueue) {
	e.result = &NodeKeeperMockGetClaimQueueResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetClaimQueue method
func (m *mNodeKeeperMockGetClaimQueue) Set(f func() (r network.ClaimQueue)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetClaimQueueFunc = f
	return m.mock
}

//GetClaimQueue implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetClaimQueue() (r network.ClaimQueue) {
	counter := atomic.AddUint64(&m.GetClaimQueuePreCounter, 1)
	defer atomic.AddUint64(&m.GetClaimQueueCounter, 1)

	if len(m.GetClaimQueueMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetClaimQueueMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetClaimQueue.")
			return
		}

		result := m.GetClaimQueueMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetClaimQueue")
			return
		}

		r = result.r

		return
	}

	if m.GetClaimQueueMock.mainExpectation != nil {

		result := m.GetClaimQueueMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetClaimQueue")
		}

		r = result.r

		return
	}

	if m.GetClaimQueueFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetClaimQueue.")
		return
	}

	return m.GetClaimQueueFunc()
}

//GetClaimQueueMinimockCounter returns a count of NodeKeeperMock.GetClaimQueueFunc invocations
func (m *NodeKeeperMock) GetClaimQueueMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetClaimQueueCounter)
}

//GetClaimQueueMinimockPreCounter returns the value of NodeKeeperMock.GetClaimQueue invocations
func (m *NodeKeeperMock) GetClaimQueueMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetClaimQueuePreCounter)
}

//GetClaimQueueFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetClaimQueueFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetClaimQueueMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetClaimQueueCounter) == uint64(len(m.GetClaimQueueMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetClaimQueueMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetClaimQueueCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetClaimQueueFunc != nil {
		return atomic.LoadUint64(&m.GetClaimQueueCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetCloudHash struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetCloudHashExpectation
	expectationSeries []*NodeKeeperMockGetCloudHashExpectation
}

type NodeKeeperMockGetCloudHashExpectation struct {
	result *NodeKeeperMockGetCloudHashResult
}

type NodeKeeperMockGetCloudHashResult struct {
	r []byte
}

//Expect specifies that invocation of NodeKeeper.GetCloudHash is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetCloudHash) Expect() *mNodeKeeperMockGetCloudHash {
	m.mock.GetCloudHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetCloudHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetCloudHash
func (m *mNodeKeeperMockGetCloudHash) Return(r []byte) *NodeKeeperMock {
	m.mock.GetCloudHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetCloudHashExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetCloudHashResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetCloudHash is expected once
func (m *mNodeKeeperMockGetCloudHash) ExpectOnce() *NodeKeeperMockGetCloudHashExpectation {
	m.mock.GetCloudHashFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetCloudHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetCloudHashExpectation) Return(r []byte) {
	e.result = &NodeKeeperMockGetCloudHashResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetCloudHash method
func (m *mNodeKeeperMockGetCloudHash) Set(f func() (r []byte)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCloudHashFunc = f
	return m.mock
}

//GetCloudHash implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetCloudHash() (r []byte) {
	counter := atomic.AddUint64(&m.GetCloudHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetCloudHashCounter, 1)

	if len(m.GetCloudHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCloudHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetCloudHash.")
			return
		}

		result := m.GetCloudHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetCloudHash")
			return
		}

		r = result.r

		return
	}

	if m.GetCloudHashMock.mainExpectation != nil {

		result := m.GetCloudHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetCloudHash")
		}

		r = result.r

		return
	}

	if m.GetCloudHashFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetCloudHash.")
		return
	}

	return m.GetCloudHashFunc()
}

//GetCloudHashMinimockCounter returns a count of NodeKeeperMock.GetCloudHashFunc invocations
func (m *NodeKeeperMock) GetCloudHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudHashCounter)
}

//GetCloudHashMinimockPreCounter returns the value of NodeKeeperMock.GetCloudHash invocations
func (m *NodeKeeperMock) GetCloudHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudHashPreCounter)
}

//GetCloudHashFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetCloudHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCloudHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCloudHashCounter) == uint64(len(m.GetCloudHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCloudHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCloudHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCloudHashFunc != nil {
		return atomic.LoadUint64(&m.GetCloudHashCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetOrigin struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetOriginExpectation
	expectationSeries []*NodeKeeperMockGetOriginExpectation
}

type NodeKeeperMockGetOriginExpectation struct {
	result *NodeKeeperMockGetOriginResult
}

type NodeKeeperMockGetOriginResult struct {
	r core.Node
}

//Expect specifies that invocation of NodeKeeper.GetOrigin is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetOrigin) Expect() *mNodeKeeperMockGetOrigin {
	m.mock.GetOriginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetOrigin
func (m *mNodeKeeperMockGetOrigin) Return(r core.Node) *NodeKeeperMock {
	m.mock.GetOriginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetOriginResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetOrigin is expected once
func (m *mNodeKeeperMockGetOrigin) ExpectOnce() *NodeKeeperMockGetOriginExpectation {
	m.mock.GetOriginFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetOriginExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetOriginExpectation) Return(r core.Node) {
	e.result = &NodeKeeperMockGetOriginResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetOrigin method
func (m *mNodeKeeperMockGetOrigin) Set(f func() (r core.Node)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOriginFunc = f
	return m.mock
}

//GetOrigin implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetOrigin() (r core.Node) {
	counter := atomic.AddUint64(&m.GetOriginPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginCounter, 1)

	if len(m.GetOriginMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOriginMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOrigin.")
			return
		}

		result := m.GetOriginMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOrigin")
			return
		}

		r = result.r

		return
	}

	if m.GetOriginMock.mainExpectation != nil {

		result := m.GetOriginMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOrigin")
		}

		r = result.r

		return
	}

	if m.GetOriginFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOrigin.")
		return
	}

	return m.GetOriginFunc()
}

//GetOriginMinimockCounter returns a count of NodeKeeperMock.GetOriginFunc invocations
func (m *NodeKeeperMock) GetOriginMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginCounter)
}

//GetOriginMinimockPreCounter returns the value of NodeKeeperMock.GetOrigin invocations
func (m *NodeKeeperMock) GetOriginMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginPreCounter)
}

//GetOriginFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetOriginFinished() bool {
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

type mNodeKeeperMockGetOriginClaim struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetOriginClaimExpectation
	expectationSeries []*NodeKeeperMockGetOriginClaimExpectation
}

type NodeKeeperMockGetOriginClaimExpectation struct {
	result *NodeKeeperMockGetOriginClaimResult
}

type NodeKeeperMockGetOriginClaimResult struct {
	r  *packets.NodeJoinClaim
	r1 error
}

//Expect specifies that invocation of NodeKeeper.GetOriginClaim is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetOriginClaim) Expect() *mNodeKeeperMockGetOriginClaim {
	m.mock.GetOriginClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginClaimExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetOriginClaim
func (m *mNodeKeeperMockGetOriginClaim) Return(r *packets.NodeJoinClaim, r1 error) *NodeKeeperMock {
	m.mock.GetOriginClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginClaimExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetOriginClaimResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetOriginClaim is expected once
func (m *mNodeKeeperMockGetOriginClaim) ExpectOnce() *NodeKeeperMockGetOriginClaimExpectation {
	m.mock.GetOriginClaimFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetOriginClaimExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetOriginClaimExpectation) Return(r *packets.NodeJoinClaim, r1 error) {
	e.result = &NodeKeeperMockGetOriginClaimResult{r, r1}
}

//Set uses given function f as a mock of NodeKeeper.GetOriginClaim method
func (m *mNodeKeeperMockGetOriginClaim) Set(f func() (r *packets.NodeJoinClaim, r1 error)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOriginClaimFunc = f
	return m.mock
}

//GetOriginClaim implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetOriginClaim() (r *packets.NodeJoinClaim, r1 error) {
	counter := atomic.AddUint64(&m.GetOriginClaimPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginClaimCounter, 1)

	if len(m.GetOriginClaimMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOriginClaimMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOriginClaim.")
			return
		}

		result := m.GetOriginClaimMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOriginClaim")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetOriginClaimMock.mainExpectation != nil {

		result := m.GetOriginClaimMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOriginClaim")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetOriginClaimFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOriginClaim.")
		return
	}

	return m.GetOriginClaimFunc()
}

//GetOriginClaimMinimockCounter returns a count of NodeKeeperMock.GetOriginClaimFunc invocations
func (m *NodeKeeperMock) GetOriginClaimMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginClaimCounter)
}

//GetOriginClaimMinimockPreCounter returns the value of NodeKeeperMock.GetOriginClaim invocations
func (m *NodeKeeperMock) GetOriginClaimMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginClaimPreCounter)
}

//GetOriginClaimFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetOriginClaimFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOriginClaimMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOriginClaimCounter) == uint64(len(m.GetOriginClaimMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOriginClaimMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOriginClaimCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOriginClaimFunc != nil {
		return atomic.LoadUint64(&m.GetOriginClaimCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetSparseUnsyncList struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetSparseUnsyncListExpectation
	expectationSeries []*NodeKeeperMockGetSparseUnsyncListExpectation
}

type NodeKeeperMockGetSparseUnsyncListExpectation struct {
	input  *NodeKeeperMockGetSparseUnsyncListInput
	result *NodeKeeperMockGetSparseUnsyncListResult
}

type NodeKeeperMockGetSparseUnsyncListInput struct {
	p int
}

type NodeKeeperMockGetSparseUnsyncListResult struct {
	r network.UnsyncList
}

//Expect specifies that invocation of NodeKeeper.GetSparseUnsyncList is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetSparseUnsyncList) Expect(p int) *mNodeKeeperMockGetSparseUnsyncList {
	m.mock.GetSparseUnsyncListFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetSparseUnsyncListExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockGetSparseUnsyncListInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.GetSparseUnsyncList
func (m *mNodeKeeperMockGetSparseUnsyncList) Return(r network.UnsyncList) *NodeKeeperMock {
	m.mock.GetSparseUnsyncListFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetSparseUnsyncListExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetSparseUnsyncListResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetSparseUnsyncList is expected once
func (m *mNodeKeeperMockGetSparseUnsyncList) ExpectOnce(p int) *NodeKeeperMockGetSparseUnsyncListExpectation {
	m.mock.GetSparseUnsyncListFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetSparseUnsyncListExpectation{}
	expectation.input = &NodeKeeperMockGetSparseUnsyncListInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetSparseUnsyncListExpectation) Return(r network.UnsyncList) {
	e.result = &NodeKeeperMockGetSparseUnsyncListResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetSparseUnsyncList method
func (m *mNodeKeeperMockGetSparseUnsyncList) Set(f func(p int) (r network.UnsyncList)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSparseUnsyncListFunc = f
	return m.mock
}

//GetSparseUnsyncList implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetSparseUnsyncList(p int) (r network.UnsyncList) {
	counter := atomic.AddUint64(&m.GetSparseUnsyncListPreCounter, 1)
	defer atomic.AddUint64(&m.GetSparseUnsyncListCounter, 1)

	if len(m.GetSparseUnsyncListMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSparseUnsyncListMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetSparseUnsyncList. %v", p)
			return
		}

		input := m.GetSparseUnsyncListMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockGetSparseUnsyncListInput{p}, "NodeKeeper.GetSparseUnsyncList got unexpected parameters")

		result := m.GetSparseUnsyncListMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetSparseUnsyncList")
			return
		}

		r = result.r

		return
	}

	if m.GetSparseUnsyncListMock.mainExpectation != nil {

		input := m.GetSparseUnsyncListMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockGetSparseUnsyncListInput{p}, "NodeKeeper.GetSparseUnsyncList got unexpected parameters")
		}

		result := m.GetSparseUnsyncListMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetSparseUnsyncList")
		}

		r = result.r

		return
	}

	if m.GetSparseUnsyncListFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetSparseUnsyncList. %v", p)
		return
	}

	return m.GetSparseUnsyncListFunc(p)
}

//GetSparseUnsyncListMinimockCounter returns a count of NodeKeeperMock.GetSparseUnsyncListFunc invocations
func (m *NodeKeeperMock) GetSparseUnsyncListMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSparseUnsyncListCounter)
}

//GetSparseUnsyncListMinimockPreCounter returns the value of NodeKeeperMock.GetSparseUnsyncList invocations
func (m *NodeKeeperMock) GetSparseUnsyncListMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSparseUnsyncListPreCounter)
}

//GetSparseUnsyncListFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetSparseUnsyncListFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSparseUnsyncListMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSparseUnsyncListCounter) == uint64(len(m.GetSparseUnsyncListMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSparseUnsyncListMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSparseUnsyncListCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSparseUnsyncListFunc != nil {
		return atomic.LoadUint64(&m.GetSparseUnsyncListCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetState struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetStateExpectation
	expectationSeries []*NodeKeeperMockGetStateExpectation
}

type NodeKeeperMockGetStateExpectation struct {
	result *NodeKeeperMockGetStateResult
}

type NodeKeeperMockGetStateResult struct {
	r network.NodeKeeperState
}

//Expect specifies that invocation of NodeKeeper.GetState is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetState) Expect() *mNodeKeeperMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetState
func (m *mNodeKeeperMockGetState) Return(r network.NodeKeeperState) *NodeKeeperMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetStateExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetState is expected once
func (m *mNodeKeeperMockGetState) ExpectOnce() *NodeKeeperMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetStateExpectation) Return(r network.NodeKeeperState) {
	e.result = &NodeKeeperMockGetStateResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetState method
func (m *mNodeKeeperMockGetState) Set(f func() (r network.NodeKeeperState)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetState() (r network.NodeKeeperState) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of NodeKeeperMock.GetStateFunc invocations
func (m *NodeKeeperMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of NodeKeeperMock.GetState invocations
func (m *NodeKeeperMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetStateFinished() bool {
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

type mNodeKeeperMockGetUnsyncList struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetUnsyncListExpectation
	expectationSeries []*NodeKeeperMockGetUnsyncListExpectation
}

type NodeKeeperMockGetUnsyncListExpectation struct {
	result *NodeKeeperMockGetUnsyncListResult
}

type NodeKeeperMockGetUnsyncListResult struct {
	r network.UnsyncList
}

//Expect specifies that invocation of NodeKeeper.GetUnsyncList is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetUnsyncList) Expect() *mNodeKeeperMockGetUnsyncList {
	m.mock.GetUnsyncListFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetUnsyncListExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetUnsyncList
func (m *mNodeKeeperMockGetUnsyncList) Return(r network.UnsyncList) *NodeKeeperMock {
	m.mock.GetUnsyncListFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetUnsyncListExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetUnsyncListResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetUnsyncList is expected once
func (m *mNodeKeeperMockGetUnsyncList) ExpectOnce() *NodeKeeperMockGetUnsyncListExpectation {
	m.mock.GetUnsyncListFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetUnsyncListExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetUnsyncListExpectation) Return(r network.UnsyncList) {
	e.result = &NodeKeeperMockGetUnsyncListResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetUnsyncList method
func (m *mNodeKeeperMockGetUnsyncList) Set(f func() (r network.UnsyncList)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetUnsyncListFunc = f
	return m.mock
}

//GetUnsyncList implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetUnsyncList() (r network.UnsyncList) {
	counter := atomic.AddUint64(&m.GetUnsyncListPreCounter, 1)
	defer atomic.AddUint64(&m.GetUnsyncListCounter, 1)

	if len(m.GetUnsyncListMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetUnsyncListMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetUnsyncList.")
			return
		}

		result := m.GetUnsyncListMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetUnsyncList")
			return
		}

		r = result.r

		return
	}

	if m.GetUnsyncListMock.mainExpectation != nil {

		result := m.GetUnsyncListMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetUnsyncList")
		}

		r = result.r

		return
	}

	if m.GetUnsyncListFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetUnsyncList.")
		return
	}

	return m.GetUnsyncListFunc()
}

//GetUnsyncListMinimockCounter returns a count of NodeKeeperMock.GetUnsyncListFunc invocations
func (m *NodeKeeperMock) GetUnsyncListMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetUnsyncListCounter)
}

//GetUnsyncListMinimockPreCounter returns the value of NodeKeeperMock.GetUnsyncList invocations
func (m *NodeKeeperMock) GetUnsyncListMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetUnsyncListPreCounter)
}

//GetUnsyncListFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetUnsyncListFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetUnsyncListMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetUnsyncListCounter) == uint64(len(m.GetUnsyncListMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetUnsyncListMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetUnsyncListCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetUnsyncListFunc != nil {
		return atomic.LoadUint64(&m.GetUnsyncListCounter) > 0
	}

	return true
}

type mNodeKeeperMockIsBootstrapped struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockIsBootstrappedExpectation
	expectationSeries []*NodeKeeperMockIsBootstrappedExpectation
}

type NodeKeeperMockIsBootstrappedExpectation struct {
	result *NodeKeeperMockIsBootstrappedResult
}

type NodeKeeperMockIsBootstrappedResult struct {
	r bool
}

//Expect specifies that invocation of NodeKeeper.IsBootstrapped is expected from 1 to Infinity times
func (m *mNodeKeeperMockIsBootstrapped) Expect() *mNodeKeeperMockIsBootstrapped {
	m.mock.IsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockIsBootstrappedExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.IsBootstrapped
func (m *mNodeKeeperMockIsBootstrapped) Return(r bool) *NodeKeeperMock {
	m.mock.IsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockIsBootstrappedExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockIsBootstrappedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.IsBootstrapped is expected once
func (m *mNodeKeeperMockIsBootstrapped) ExpectOnce() *NodeKeeperMockIsBootstrappedExpectation {
	m.mock.IsBootstrappedFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockIsBootstrappedExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockIsBootstrappedExpectation) Return(r bool) {
	e.result = &NodeKeeperMockIsBootstrappedResult{r}
}

//Set uses given function f as a mock of NodeKeeper.IsBootstrapped method
func (m *mNodeKeeperMockIsBootstrapped) Set(f func() (r bool)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsBootstrappedFunc = f
	return m.mock
}

//IsBootstrapped implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) IsBootstrapped() (r bool) {
	counter := atomic.AddUint64(&m.IsBootstrappedPreCounter, 1)
	defer atomic.AddUint64(&m.IsBootstrappedCounter, 1)

	if len(m.IsBootstrappedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsBootstrappedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.IsBootstrapped.")
			return
		}

		result := m.IsBootstrappedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.IsBootstrapped")
			return
		}

		r = result.r

		return
	}

	if m.IsBootstrappedMock.mainExpectation != nil {

		result := m.IsBootstrappedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.IsBootstrapped")
		}

		r = result.r

		return
	}

	if m.IsBootstrappedFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.IsBootstrapped.")
		return
	}

	return m.IsBootstrappedFunc()
}

//IsBootstrappedMinimockCounter returns a count of NodeKeeperMock.IsBootstrappedFunc invocations
func (m *NodeKeeperMock) IsBootstrappedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsBootstrappedCounter)
}

//IsBootstrappedMinimockPreCounter returns the value of NodeKeeperMock.IsBootstrapped invocations
func (m *NodeKeeperMock) IsBootstrappedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsBootstrappedPreCounter)
}

//IsBootstrappedFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) IsBootstrappedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsBootstrappedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsBootstrappedCounter) == uint64(len(m.IsBootstrappedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsBootstrappedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsBootstrappedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsBootstrappedFunc != nil {
		return atomic.LoadUint64(&m.IsBootstrappedCounter) > 0
	}

	return true
}

type mNodeKeeperMockMoveSyncToActive struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockMoveSyncToActiveExpectation
	expectationSeries []*NodeKeeperMockMoveSyncToActiveExpectation
}

type NodeKeeperMockMoveSyncToActiveExpectation struct {
}

//Expect specifies that invocation of NodeKeeper.MoveSyncToActive is expected from 1 to Infinity times
func (m *mNodeKeeperMockMoveSyncToActive) Expect() *mNodeKeeperMockMoveSyncToActive {
	m.mock.MoveSyncToActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockMoveSyncToActiveExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.MoveSyncToActive
func (m *mNodeKeeperMockMoveSyncToActive) Return() *NodeKeeperMock {
	m.mock.MoveSyncToActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockMoveSyncToActiveExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.MoveSyncToActive is expected once
func (m *mNodeKeeperMockMoveSyncToActive) ExpectOnce() *NodeKeeperMockMoveSyncToActiveExpectation {
	m.mock.MoveSyncToActiveFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockMoveSyncToActiveExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeKeeper.MoveSyncToActive method
func (m *mNodeKeeperMockMoveSyncToActive) Set(f func()) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MoveSyncToActiveFunc = f
	return m.mock
}

//MoveSyncToActive implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) MoveSyncToActive() {
	counter := atomic.AddUint64(&m.MoveSyncToActivePreCounter, 1)
	defer atomic.AddUint64(&m.MoveSyncToActiveCounter, 1)

	if len(m.MoveSyncToActiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MoveSyncToActiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.MoveSyncToActive.")
			return
		}

		return
	}

	if m.MoveSyncToActiveMock.mainExpectation != nil {

		return
	}

	if m.MoveSyncToActiveFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.MoveSyncToActive.")
		return
	}

	m.MoveSyncToActiveFunc()
}

//MoveSyncToActiveMinimockCounter returns a count of NodeKeeperMock.MoveSyncToActiveFunc invocations
func (m *NodeKeeperMock) MoveSyncToActiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MoveSyncToActiveCounter)
}

//MoveSyncToActiveMinimockPreCounter returns the value of NodeKeeperMock.MoveSyncToActive invocations
func (m *NodeKeeperMock) MoveSyncToActiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MoveSyncToActivePreCounter)
}

//MoveSyncToActiveFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) MoveSyncToActiveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MoveSyncToActiveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MoveSyncToActiveCounter) == uint64(len(m.MoveSyncToActiveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MoveSyncToActiveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MoveSyncToActiveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MoveSyncToActiveFunc != nil {
		return atomic.LoadUint64(&m.MoveSyncToActiveCounter) > 0
	}

	return true
}

type mNodeKeeperMockNodesJoinedDuringPreviousPulse struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation
	expectationSeries []*NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation
}

type NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation struct {
	result *NodeKeeperMockNodesJoinedDuringPreviousPulseResult
}

type NodeKeeperMockNodesJoinedDuringPreviousPulseResult struct {
	r bool
}

//Expect specifies that invocation of NodeKeeper.NodesJoinedDuringPreviousPulse is expected from 1 to Infinity times
func (m *mNodeKeeperMockNodesJoinedDuringPreviousPulse) Expect() *mNodeKeeperMockNodesJoinedDuringPreviousPulse {
	m.mock.NodesJoinedDuringPreviousPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.NodesJoinedDuringPreviousPulse
func (m *mNodeKeeperMockNodesJoinedDuringPreviousPulse) Return(r bool) *NodeKeeperMock {
	m.mock.NodesJoinedDuringPreviousPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockNodesJoinedDuringPreviousPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.NodesJoinedDuringPreviousPulse is expected once
func (m *mNodeKeeperMockNodesJoinedDuringPreviousPulse) ExpectOnce() *NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation {
	m.mock.NodesJoinedDuringPreviousPulseFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockNodesJoinedDuringPreviousPulseExpectation) Return(r bool) {
	e.result = &NodeKeeperMockNodesJoinedDuringPreviousPulseResult{r}
}

//Set uses given function f as a mock of NodeKeeper.NodesJoinedDuringPreviousPulse method
func (m *mNodeKeeperMockNodesJoinedDuringPreviousPulse) Set(f func() (r bool)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NodesJoinedDuringPreviousPulseFunc = f
	return m.mock
}

//NodesJoinedDuringPreviousPulse implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) NodesJoinedDuringPreviousPulse() (r bool) {
	counter := atomic.AddUint64(&m.NodesJoinedDuringPreviousPulsePreCounter, 1)
	defer atomic.AddUint64(&m.NodesJoinedDuringPreviousPulseCounter, 1)

	if len(m.NodesJoinedDuringPreviousPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NodesJoinedDuringPreviousPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse.")
			return
		}

		result := m.NodesJoinedDuringPreviousPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.NodesJoinedDuringPreviousPulse")
			return
		}

		r = result.r

		return
	}

	if m.NodesJoinedDuringPreviousPulseMock.mainExpectation != nil {

		result := m.NodesJoinedDuringPreviousPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.NodesJoinedDuringPreviousPulse")
		}

		r = result.r

		return
	}

	if m.NodesJoinedDuringPreviousPulseFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse.")
		return
	}

	return m.NodesJoinedDuringPreviousPulseFunc()
}

//NodesJoinedDuringPreviousPulseMinimockCounter returns a count of NodeKeeperMock.NodesJoinedDuringPreviousPulseFunc invocations
func (m *NodeKeeperMock) NodesJoinedDuringPreviousPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter)
}

//NodesJoinedDuringPreviousPulseMinimockPreCounter returns the value of NodeKeeperMock.NodesJoinedDuringPreviousPulse invocations
func (m *NodeKeeperMock) NodesJoinedDuringPreviousPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulsePreCounter)
}

//NodesJoinedDuringPreviousPulseFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) NodesJoinedDuringPreviousPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NodesJoinedDuringPreviousPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) == uint64(len(m.NodesJoinedDuringPreviousPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NodesJoinedDuringPreviousPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NodesJoinedDuringPreviousPulseFunc != nil {
		return atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) > 0
	}

	return true
}

type mNodeKeeperMockSetCloudHash struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockSetCloudHashExpectation
	expectationSeries []*NodeKeeperMockSetCloudHashExpectation
}

type NodeKeeperMockSetCloudHashExpectation struct {
	input *NodeKeeperMockSetCloudHashInput
}

type NodeKeeperMockSetCloudHashInput struct {
	p []byte
}

//Expect specifies that invocation of NodeKeeper.SetCloudHash is expected from 1 to Infinity times
func (m *mNodeKeeperMockSetCloudHash) Expect(p []byte) *mNodeKeeperMockSetCloudHash {
	m.mock.SetCloudHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetCloudHashExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockSetCloudHashInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.SetCloudHash
func (m *mNodeKeeperMockSetCloudHash) Return() *NodeKeeperMock {
	m.mock.SetCloudHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetCloudHashExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.SetCloudHash is expected once
func (m *mNodeKeeperMockSetCloudHash) ExpectOnce(p []byte) *NodeKeeperMockSetCloudHashExpectation {
	m.mock.SetCloudHashFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockSetCloudHashExpectation{}
	expectation.input = &NodeKeeperMockSetCloudHashInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeKeeper.SetCloudHash method
func (m *mNodeKeeperMockSetCloudHash) Set(f func(p []byte)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetCloudHashFunc = f
	return m.mock
}

//SetCloudHash implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) SetCloudHash(p []byte) {
	counter := atomic.AddUint64(&m.SetCloudHashPreCounter, 1)
	defer atomic.AddUint64(&m.SetCloudHashCounter, 1)

	if len(m.SetCloudHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetCloudHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.SetCloudHash. %v", p)
			return
		}

		input := m.SetCloudHashMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockSetCloudHashInput{p}, "NodeKeeper.SetCloudHash got unexpected parameters")

		return
	}

	if m.SetCloudHashMock.mainExpectation != nil {

		input := m.SetCloudHashMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockSetCloudHashInput{p}, "NodeKeeper.SetCloudHash got unexpected parameters")
		}

		return
	}

	if m.SetCloudHashFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.SetCloudHash. %v", p)
		return
	}

	m.SetCloudHashFunc(p)
}

//SetCloudHashMinimockCounter returns a count of NodeKeeperMock.SetCloudHashFunc invocations
func (m *NodeKeeperMock) SetCloudHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCloudHashCounter)
}

//SetCloudHashMinimockPreCounter returns the value of NodeKeeperMock.SetCloudHash invocations
func (m *NodeKeeperMock) SetCloudHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetCloudHashPreCounter)
}

//SetCloudHashFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) SetCloudHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetCloudHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCloudHashCounter) == uint64(len(m.SetCloudHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetCloudHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCloudHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetCloudHashFunc != nil {
		return atomic.LoadUint64(&m.SetCloudHashCounter) > 0
	}

	return true
}

type mNodeKeeperMockSetIsBootstrapped struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockSetIsBootstrappedExpectation
	expectationSeries []*NodeKeeperMockSetIsBootstrappedExpectation
}

type NodeKeeperMockSetIsBootstrappedExpectation struct {
	input *NodeKeeperMockSetIsBootstrappedInput
}

type NodeKeeperMockSetIsBootstrappedInput struct {
	p bool
}

//Expect specifies that invocation of NodeKeeper.SetIsBootstrapped is expected from 1 to Infinity times
func (m *mNodeKeeperMockSetIsBootstrapped) Expect(p bool) *mNodeKeeperMockSetIsBootstrapped {
	m.mock.SetIsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetIsBootstrappedExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockSetIsBootstrappedInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.SetIsBootstrapped
func (m *mNodeKeeperMockSetIsBootstrapped) Return() *NodeKeeperMock {
	m.mock.SetIsBootstrappedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetIsBootstrappedExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.SetIsBootstrapped is expected once
func (m *mNodeKeeperMockSetIsBootstrapped) ExpectOnce(p bool) *NodeKeeperMockSetIsBootstrappedExpectation {
	m.mock.SetIsBootstrappedFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockSetIsBootstrappedExpectation{}
	expectation.input = &NodeKeeperMockSetIsBootstrappedInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeKeeper.SetIsBootstrapped method
func (m *mNodeKeeperMockSetIsBootstrapped) Set(f func(p bool)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetIsBootstrappedFunc = f
	return m.mock
}

//SetIsBootstrapped implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) SetIsBootstrapped(p bool) {
	counter := atomic.AddUint64(&m.SetIsBootstrappedPreCounter, 1)
	defer atomic.AddUint64(&m.SetIsBootstrappedCounter, 1)

	if len(m.SetIsBootstrappedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetIsBootstrappedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.SetIsBootstrapped. %v", p)
			return
		}

		input := m.SetIsBootstrappedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockSetIsBootstrappedInput{p}, "NodeKeeper.SetIsBootstrapped got unexpected parameters")

		return
	}

	if m.SetIsBootstrappedMock.mainExpectation != nil {

		input := m.SetIsBootstrappedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockSetIsBootstrappedInput{p}, "NodeKeeper.SetIsBootstrapped got unexpected parameters")
		}

		return
	}

	if m.SetIsBootstrappedFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.SetIsBootstrapped. %v", p)
		return
	}

	m.SetIsBootstrappedFunc(p)
}

//SetIsBootstrappedMinimockCounter returns a count of NodeKeeperMock.SetIsBootstrappedFunc invocations
func (m *NodeKeeperMock) SetIsBootstrappedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetIsBootstrappedCounter)
}

//SetIsBootstrappedMinimockPreCounter returns the value of NodeKeeperMock.SetIsBootstrapped invocations
func (m *NodeKeeperMock) SetIsBootstrappedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetIsBootstrappedPreCounter)
}

//SetIsBootstrappedFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) SetIsBootstrappedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetIsBootstrappedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetIsBootstrappedCounter) == uint64(len(m.SetIsBootstrappedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetIsBootstrappedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetIsBootstrappedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetIsBootstrappedFunc != nil {
		return atomic.LoadUint64(&m.SetIsBootstrappedCounter) > 0
	}

	return true
}

type mNodeKeeperMockSetState struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockSetStateExpectation
	expectationSeries []*NodeKeeperMockSetStateExpectation
}

type NodeKeeperMockSetStateExpectation struct {
	input *NodeKeeperMockSetStateInput
}

type NodeKeeperMockSetStateInput struct {
	p network.NodeKeeperState
}

//Expect specifies that invocation of NodeKeeper.SetState is expected from 1 to Infinity times
func (m *mNodeKeeperMockSetState) Expect(p network.NodeKeeperState) *mNodeKeeperMockSetState {
	m.mock.SetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetStateExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockSetStateInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.SetState
func (m *mNodeKeeperMockSetState) Return() *NodeKeeperMock {
	m.mock.SetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetStateExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.SetState is expected once
func (m *mNodeKeeperMockSetState) ExpectOnce(p network.NodeKeeperState) *NodeKeeperMockSetStateExpectation {
	m.mock.SetStateFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockSetStateExpectation{}
	expectation.input = &NodeKeeperMockSetStateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeKeeper.SetState method
func (m *mNodeKeeperMockSetState) Set(f func(p network.NodeKeeperState)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetStateFunc = f
	return m.mock
}

//SetState implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) SetState(p network.NodeKeeperState) {
	counter := atomic.AddUint64(&m.SetStatePreCounter, 1)
	defer atomic.AddUint64(&m.SetStateCounter, 1)

	if len(m.SetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.SetState. %v", p)
			return
		}

		input := m.SetStateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockSetStateInput{p}, "NodeKeeper.SetState got unexpected parameters")

		return
	}

	if m.SetStateMock.mainExpectation != nil {

		input := m.SetStateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockSetStateInput{p}, "NodeKeeper.SetState got unexpected parameters")
		}

		return
	}

	if m.SetStateFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.SetState. %v", p)
		return
	}

	m.SetStateFunc(p)
}

//SetStateMinimockCounter returns a count of NodeKeeperMock.SetStateFunc invocations
func (m *NodeKeeperMock) SetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetStateCounter)
}

//SetStateMinimockPreCounter returns the value of NodeKeeperMock.SetState invocations
func (m *NodeKeeperMock) SetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetStatePreCounter)
}

//SetStateFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) SetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetStateCounter) == uint64(len(m.SetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetStateFunc != nil {
		return atomic.LoadUint64(&m.SetStateCounter) > 0
	}

	return true
}

type mNodeKeeperMockSync struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockSyncExpectation
	expectationSeries []*NodeKeeperMockSyncExpectation
}

type NodeKeeperMockSyncExpectation struct {
	input *NodeKeeperMockSyncInput
}

type NodeKeeperMockSyncInput struct {
	p network.UnsyncList
}

//Expect specifies that invocation of NodeKeeper.Sync is expected from 1 to Infinity times
func (m *mNodeKeeperMockSync) Expect(p network.UnsyncList) *mNodeKeeperMockSync {
	m.mock.SyncFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSyncExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockSyncInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.Sync
func (m *mNodeKeeperMockSync) Return() *NodeKeeperMock {
	m.mock.SyncFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSyncExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.Sync is expected once
func (m *mNodeKeeperMockSync) ExpectOnce(p network.UnsyncList) *NodeKeeperMockSyncExpectation {
	m.mock.SyncFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockSyncExpectation{}
	expectation.input = &NodeKeeperMockSyncInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeKeeper.Sync method
func (m *mNodeKeeperMockSync) Set(f func(p network.UnsyncList)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SyncFunc = f
	return m.mock
}

//Sync implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) Sync(p network.UnsyncList) {
	counter := atomic.AddUint64(&m.SyncPreCounter, 1)
	defer atomic.AddUint64(&m.SyncCounter, 1)

	if len(m.SyncMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SyncMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.Sync. %v", p)
			return
		}

		input := m.SyncMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockSyncInput{p}, "NodeKeeper.Sync got unexpected parameters")

		return
	}

	if m.SyncMock.mainExpectation != nil {

		input := m.SyncMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockSyncInput{p}, "NodeKeeper.Sync got unexpected parameters")
		}

		return
	}

	if m.SyncFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.Sync. %v", p)
		return
	}

	m.SyncFunc(p)
}

//SyncMinimockCounter returns a count of NodeKeeperMock.SyncFunc invocations
func (m *NodeKeeperMock) SyncMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SyncCounter)
}

//SyncMinimockPreCounter returns the value of NodeKeeperMock.Sync invocations
func (m *NodeKeeperMock) SyncMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SyncPreCounter)
}

//SyncFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) SyncFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SyncMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SyncCounter) == uint64(len(m.SyncMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SyncMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SyncCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SyncFunc != nil {
		return atomic.LoadUint64(&m.SyncCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeKeeperMock) ValidateCallCounters() {

	if !m.AddActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.AddActiveNodes")
	}

	if !m.AddPendingClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.AddPendingClaim")
	}

	if !m.GetActiveNodeFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNode")
	}

	if !m.GetActiveNodeByShortIDFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodeByShortID")
	}

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodes")
	}

	if !m.GetActiveNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodesByRole")
	}

	if !m.GetClaimQueueFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetClaimQueue")
	}

	if !m.GetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetCloudHash")
	}

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOrigin")
	}

	if !m.GetOriginClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginClaim")
	}

	if !m.GetSparseUnsyncListFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetSparseUnsyncList")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetState")
	}

	if !m.GetUnsyncListFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetUnsyncList")
	}

	if !m.IsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.IsBootstrapped")
	}

	if !m.MoveSyncToActiveFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.MoveSyncToActive")
	}

	if !m.NodesJoinedDuringPreviousPulseFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse")
	}

	if !m.SetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetCloudHash")
	}

	if !m.SetIsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetIsBootstrapped")
	}

	if !m.SetStateFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetState")
	}

	if !m.SyncFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.Sync")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeKeeperMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeKeeperMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeKeeperMock) MinimockFinish() {

	if !m.AddActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.AddActiveNodes")
	}

	if !m.AddPendingClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.AddPendingClaim")
	}

	if !m.GetActiveNodeFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNode")
	}

	if !m.GetActiveNodeByShortIDFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodeByShortID")
	}

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodes")
	}

	if !m.GetActiveNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodesByRole")
	}

	if !m.GetClaimQueueFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetClaimQueue")
	}

	if !m.GetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetCloudHash")
	}

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOrigin")
	}

	if !m.GetOriginClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginClaim")
	}

	if !m.GetSparseUnsyncListFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetSparseUnsyncList")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetState")
	}

	if !m.GetUnsyncListFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetUnsyncList")
	}

	if !m.IsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.IsBootstrapped")
	}

	if !m.MoveSyncToActiveFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.MoveSyncToActive")
	}

	if !m.NodesJoinedDuringPreviousPulseFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse")
	}

	if !m.SetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetCloudHash")
	}

	if !m.SetIsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetIsBootstrapped")
	}

	if !m.SetStateFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetState")
	}

	if !m.SyncFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.Sync")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeKeeperMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeKeeperMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddActiveNodesFinished()
		ok = ok && m.AddPendingClaimFinished()
		ok = ok && m.GetActiveNodeFinished()
		ok = ok && m.GetActiveNodeByShortIDFinished()
		ok = ok && m.GetActiveNodesFinished()
		ok = ok && m.GetActiveNodesByRoleFinished()
		ok = ok && m.GetClaimQueueFinished()
		ok = ok && m.GetCloudHashFinished()
		ok = ok && m.GetOriginFinished()
		ok = ok && m.GetOriginClaimFinished()
		ok = ok && m.GetSparseUnsyncListFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.GetUnsyncListFinished()
		ok = ok && m.IsBootstrappedFinished()
		ok = ok && m.MoveSyncToActiveFinished()
		ok = ok && m.NodesJoinedDuringPreviousPulseFinished()
		ok = ok && m.SetCloudHashFinished()
		ok = ok && m.SetIsBootstrappedFinished()
		ok = ok && m.SetStateFinished()
		ok = ok && m.SyncFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddActiveNodesFinished() {
				m.t.Error("Expected call to NodeKeeperMock.AddActiveNodes")
			}

			if !m.AddPendingClaimFinished() {
				m.t.Error("Expected call to NodeKeeperMock.AddPendingClaim")
			}

			if !m.GetActiveNodeFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNode")
			}

			if !m.GetActiveNodeByShortIDFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNodeByShortID")
			}

			if !m.GetActiveNodesFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNodes")
			}

			if !m.GetActiveNodesByRoleFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNodesByRole")
			}

			if !m.GetClaimQueueFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetClaimQueue")
			}

			if !m.GetCloudHashFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetCloudHash")
			}

			if !m.GetOriginFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetOrigin")
			}

			if !m.GetOriginClaimFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetOriginClaim")
			}

			if !m.GetSparseUnsyncListFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetSparseUnsyncList")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetState")
			}

			if !m.GetUnsyncListFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetUnsyncList")
			}

			if !m.IsBootstrappedFinished() {
				m.t.Error("Expected call to NodeKeeperMock.IsBootstrapped")
			}

			if !m.MoveSyncToActiveFinished() {
				m.t.Error("Expected call to NodeKeeperMock.MoveSyncToActive")
			}

			if !m.NodesJoinedDuringPreviousPulseFinished() {
				m.t.Error("Expected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse")
			}

			if !m.SetCloudHashFinished() {
				m.t.Error("Expected call to NodeKeeperMock.SetCloudHash")
			}

			if !m.SetIsBootstrappedFinished() {
				m.t.Error("Expected call to NodeKeeperMock.SetIsBootstrapped")
			}

			if !m.SetStateFinished() {
				m.t.Error("Expected call to NodeKeeperMock.SetState")
			}

			if !m.SyncFinished() {
				m.t.Error("Expected call to NodeKeeperMock.Sync")
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
func (m *NodeKeeperMock) AllMocksCalled() bool {

	if !m.AddActiveNodesFinished() {
		return false
	}

	if !m.AddPendingClaimFinished() {
		return false
	}

	if !m.GetActiveNodeFinished() {
		return false
	}

	if !m.GetActiveNodeByShortIDFinished() {
		return false
	}

	if !m.GetActiveNodesFinished() {
		return false
	}

	if !m.GetActiveNodesByRoleFinished() {
		return false
	}

	if !m.GetClaimQueueFinished() {
		return false
	}

	if !m.GetCloudHashFinished() {
		return false
	}

	if !m.GetOriginFinished() {
		return false
	}

	if !m.GetOriginClaimFinished() {
		return false
	}

	if !m.GetSparseUnsyncListFinished() {
		return false
	}

	if !m.GetStateFinished() {
		return false
	}

	if !m.GetUnsyncListFinished() {
		return false
	}

	if !m.IsBootstrappedFinished() {
		return false
	}

	if !m.MoveSyncToActiveFinished() {
		return false
	}

	if !m.NodesJoinedDuringPreviousPulseFinished() {
		return false
	}

	if !m.SetCloudHashFinished() {
		return false
	}

	if !m.SetIsBootstrappedFinished() {
		return false
	}

	if !m.SetStateFinished() {
		return false
	}

	if !m.SyncFinished() {
		return false
	}

	return true
}

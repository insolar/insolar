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
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
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

	GetOriginClaimFunc       func() (r *packets.NodeJoinClaim)
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

	SetOriginClaimFunc       func(p *packets.NodeJoinClaim)
	SetOriginClaimCounter    uint64
	SetOriginClaimPreCounter uint64
	SetOriginClaimMock       mNodeKeeperMockSetOriginClaim

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
	m.MoveSyncToActiveMock = mNodeKeeperMockMoveSyncToActive{mock: m}
	m.NodesJoinedDuringPreviousPulseMock = mNodeKeeperMockNodesJoinedDuringPreviousPulse{mock: m}
	m.SetCloudHashMock = mNodeKeeperMockSetCloudHash{mock: m}
	m.SetOriginClaimMock = mNodeKeeperMockSetOriginClaim{mock: m}
	m.SetStateMock = mNodeKeeperMockSetState{mock: m}
	m.SyncMock = mNodeKeeperMockSync{mock: m}

	return m
}

type mNodeKeeperMockAddActiveNodes struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockAddActiveNodesParams
}

//NodeKeeperMockAddActiveNodesParams represents input parameters of the NodeKeeper.AddActiveNodes
type NodeKeeperMockAddActiveNodesParams struct {
	p []core.Node
}

//Expect sets up expected params for the NodeKeeper.AddActiveNodes
func (m *mNodeKeeperMockAddActiveNodes) Expect(p []core.Node) *mNodeKeeperMockAddActiveNodes {
	m.mockExpectations = &NodeKeeperMockAddActiveNodesParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.AddActiveNodes to return Return's arguments
func (m *mNodeKeeperMockAddActiveNodes) Return() *NodeKeeperMock {
	m.mock.AddActiveNodesFunc = func(p []core.Node) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.AddActiveNodes method
func (m *mNodeKeeperMockAddActiveNodes) Set(f func(p []core.Node)) *NodeKeeperMock {
	m.mock.AddActiveNodesFunc = f
	m.mockExpectations = nil
	return m.mock
}

//AddActiveNodes implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) AddActiveNodes(p []core.Node) {
	atomic.AddUint64(&m.AddActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.AddActiveNodesCounter, 1)

	if m.AddActiveNodesMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.AddActiveNodesMock.mockExpectations, NodeKeeperMockAddActiveNodesParams{p},
			"NodeKeeper.AddActiveNodes got unexpected parameters")

		if m.AddActiveNodesFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.AddActiveNodes")

			return
		}
	}

	if m.AddActiveNodesFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.AddActiveNodes")
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

type mNodeKeeperMockAddPendingClaim struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockAddPendingClaimParams
}

//NodeKeeperMockAddPendingClaimParams represents input parameters of the NodeKeeper.AddPendingClaim
type NodeKeeperMockAddPendingClaimParams struct {
	p packets.ReferendumClaim
}

//Expect sets up expected params for the NodeKeeper.AddPendingClaim
func (m *mNodeKeeperMockAddPendingClaim) Expect(p packets.ReferendumClaim) *mNodeKeeperMockAddPendingClaim {
	m.mockExpectations = &NodeKeeperMockAddPendingClaimParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.AddPendingClaim to return Return's arguments
func (m *mNodeKeeperMockAddPendingClaim) Return(r bool) *NodeKeeperMock {
	m.mock.AddPendingClaimFunc = func(p packets.ReferendumClaim) bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.AddPendingClaim method
func (m *mNodeKeeperMockAddPendingClaim) Set(f func(p packets.ReferendumClaim) (r bool)) *NodeKeeperMock {
	m.mock.AddPendingClaimFunc = f
	m.mockExpectations = nil
	return m.mock
}

//AddPendingClaim implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) AddPendingClaim(p packets.ReferendumClaim) (r bool) {
	atomic.AddUint64(&m.AddPendingClaimPreCounter, 1)
	defer atomic.AddUint64(&m.AddPendingClaimCounter, 1)

	if m.AddPendingClaimMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.AddPendingClaimMock.mockExpectations, NodeKeeperMockAddPendingClaimParams{p},
			"NodeKeeper.AddPendingClaim got unexpected parameters")

		if m.AddPendingClaimFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.AddPendingClaim")

			return
		}
	}

	if m.AddPendingClaimFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.AddPendingClaim")
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

type mNodeKeeperMockGetActiveNode struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockGetActiveNodeParams
}

//NodeKeeperMockGetActiveNodeParams represents input parameters of the NodeKeeper.GetActiveNode
type NodeKeeperMockGetActiveNodeParams struct {
	p core.RecordRef
}

//Expect sets up expected params for the NodeKeeper.GetActiveNode
func (m *mNodeKeeperMockGetActiveNode) Expect(p core.RecordRef) *mNodeKeeperMockGetActiveNode {
	m.mockExpectations = &NodeKeeperMockGetActiveNodeParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.GetActiveNode to return Return's arguments
func (m *mNodeKeeperMockGetActiveNode) Return(r core.Node) *NodeKeeperMock {
	m.mock.GetActiveNodeFunc = func(p core.RecordRef) core.Node {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNode method
func (m *mNodeKeeperMockGetActiveNode) Set(f func(p core.RecordRef) (r core.Node)) *NodeKeeperMock {
	m.mock.GetActiveNodeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetActiveNode implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNode(p core.RecordRef) (r core.Node) {
	atomic.AddUint64(&m.GetActiveNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeCounter, 1)

	if m.GetActiveNodeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetActiveNodeMock.mockExpectations, NodeKeeperMockGetActiveNodeParams{p},
			"NodeKeeper.GetActiveNode got unexpected parameters")

		if m.GetActiveNodeFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNode")

			return
		}
	}

	if m.GetActiveNodeFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetActiveNode")
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

type mNodeKeeperMockGetActiveNodeByShortID struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockGetActiveNodeByShortIDParams
}

//NodeKeeperMockGetActiveNodeByShortIDParams represents input parameters of the NodeKeeper.GetActiveNodeByShortID
type NodeKeeperMockGetActiveNodeByShortIDParams struct {
	p core.ShortNodeID
}

//Expect sets up expected params for the NodeKeeper.GetActiveNodeByShortID
func (m *mNodeKeeperMockGetActiveNodeByShortID) Expect(p core.ShortNodeID) *mNodeKeeperMockGetActiveNodeByShortID {
	m.mockExpectations = &NodeKeeperMockGetActiveNodeByShortIDParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.GetActiveNodeByShortID to return Return's arguments
func (m *mNodeKeeperMockGetActiveNodeByShortID) Return(r core.Node) *NodeKeeperMock {
	m.mock.GetActiveNodeByShortIDFunc = func(p core.ShortNodeID) core.Node {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNodeByShortID method
func (m *mNodeKeeperMockGetActiveNodeByShortID) Set(f func(p core.ShortNodeID) (r core.Node)) *NodeKeeperMock {
	m.mock.GetActiveNodeByShortIDFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetActiveNodeByShortID implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNodeByShortID(p core.ShortNodeID) (r core.Node) {
	atomic.AddUint64(&m.GetActiveNodeByShortIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeByShortIDCounter, 1)

	if m.GetActiveNodeByShortIDMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetActiveNodeByShortIDMock.mockExpectations, NodeKeeperMockGetActiveNodeByShortIDParams{p},
			"NodeKeeper.GetActiveNodeByShortID got unexpected parameters")

		if m.GetActiveNodeByShortIDFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodeByShortID")

			return
		}
	}

	if m.GetActiveNodeByShortIDFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetActiveNodeByShortID")
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

type mNodeKeeperMockGetActiveNodes struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.GetActiveNodes to return Return's arguments
func (m *mNodeKeeperMockGetActiveNodes) Return(r []core.Node) *NodeKeeperMock {
	m.mock.GetActiveNodesFunc = func() []core.Node {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNodes method
func (m *mNodeKeeperMockGetActiveNodes) Set(f func() (r []core.Node)) *NodeKeeperMock {
	m.mock.GetActiveNodesFunc = f

	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNodes() (r []core.Node) {
	atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if m.GetActiveNodesFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetActiveNodes")
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

type mNodeKeeperMockGetActiveNodesByRole struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockGetActiveNodesByRoleParams
}

//NodeKeeperMockGetActiveNodesByRoleParams represents input parameters of the NodeKeeper.GetActiveNodesByRole
type NodeKeeperMockGetActiveNodesByRoleParams struct {
	p core.DynamicRole
}

//Expect sets up expected params for the NodeKeeper.GetActiveNodesByRole
func (m *mNodeKeeperMockGetActiveNodesByRole) Expect(p core.DynamicRole) *mNodeKeeperMockGetActiveNodesByRole {
	m.mockExpectations = &NodeKeeperMockGetActiveNodesByRoleParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.GetActiveNodesByRole to return Return's arguments
func (m *mNodeKeeperMockGetActiveNodesByRole) Return(r []core.RecordRef) *NodeKeeperMock {
	m.mock.GetActiveNodesByRoleFunc = func(p core.DynamicRole) []core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetActiveNodesByRole method
func (m *mNodeKeeperMockGetActiveNodesByRole) Set(f func(p core.DynamicRole) (r []core.RecordRef)) *NodeKeeperMock {
	m.mock.GetActiveNodesByRoleFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetActiveNodesByRole implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetActiveNodesByRole(p core.DynamicRole) (r []core.RecordRef) {
	atomic.AddUint64(&m.GetActiveNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesByRoleCounter, 1)

	if m.GetActiveNodesByRoleMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetActiveNodesByRoleMock.mockExpectations, NodeKeeperMockGetActiveNodesByRoleParams{p},
			"NodeKeeper.GetActiveNodesByRole got unexpected parameters")

		if m.GetActiveNodesByRoleFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.GetActiveNodesByRole")

			return
		}
	}

	if m.GetActiveNodesByRoleFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetActiveNodesByRole")
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

type mNodeKeeperMockGetClaimQueue struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.GetClaimQueue to return Return's arguments
func (m *mNodeKeeperMockGetClaimQueue) Return(r network.ClaimQueue) *NodeKeeperMock {
	m.mock.GetClaimQueueFunc = func() network.ClaimQueue {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetClaimQueue method
func (m *mNodeKeeperMockGetClaimQueue) Set(f func() (r network.ClaimQueue)) *NodeKeeperMock {
	m.mock.GetClaimQueueFunc = f

	return m.mock
}

//GetClaimQueue implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetClaimQueue() (r network.ClaimQueue) {
	atomic.AddUint64(&m.GetClaimQueuePreCounter, 1)
	defer atomic.AddUint64(&m.GetClaimQueueCounter, 1)

	if m.GetClaimQueueFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetClaimQueue")
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

type mNodeKeeperMockGetCloudHash struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.GetCloudHash to return Return's arguments
func (m *mNodeKeeperMockGetCloudHash) Return(r []byte) *NodeKeeperMock {
	m.mock.GetCloudHashFunc = func() []byte {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetCloudHash method
func (m *mNodeKeeperMockGetCloudHash) Set(f func() (r []byte)) *NodeKeeperMock {
	m.mock.GetCloudHashFunc = f

	return m.mock
}

//GetCloudHash implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetCloudHash() (r []byte) {
	atomic.AddUint64(&m.GetCloudHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetCloudHashCounter, 1)

	if m.GetCloudHashFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetCloudHash")
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

type mNodeKeeperMockGetOrigin struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.GetOrigin to return Return's arguments
func (m *mNodeKeeperMockGetOrigin) Return(r core.Node) *NodeKeeperMock {
	m.mock.GetOriginFunc = func() core.Node {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetOrigin method
func (m *mNodeKeeperMockGetOrigin) Set(f func() (r core.Node)) *NodeKeeperMock {
	m.mock.GetOriginFunc = f

	return m.mock
}

//GetOrigin implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetOrigin() (r core.Node) {
	atomic.AddUint64(&m.GetOriginPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginCounter, 1)

	if m.GetOriginFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetOrigin")
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

type mNodeKeeperMockGetOriginClaim struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.GetOriginClaim to return Return's arguments
func (m *mNodeKeeperMockGetOriginClaim) Return(r *packets.NodeJoinClaim) *NodeKeeperMock {
	m.mock.GetOriginClaimFunc = func() *packets.NodeJoinClaim {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetOriginClaim method
func (m *mNodeKeeperMockGetOriginClaim) Set(f func() (r *packets.NodeJoinClaim)) *NodeKeeperMock {
	m.mock.GetOriginClaimFunc = f

	return m.mock
}

//GetOriginClaim implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetOriginClaim() (r *packets.NodeJoinClaim) {
	atomic.AddUint64(&m.GetOriginClaimPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginClaimCounter, 1)

	if m.GetOriginClaimFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetOriginClaim")
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

type mNodeKeeperMockGetSparseUnsyncList struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockGetSparseUnsyncListParams
}

//NodeKeeperMockGetSparseUnsyncListParams represents input parameters of the NodeKeeper.GetSparseUnsyncList
type NodeKeeperMockGetSparseUnsyncListParams struct {
	p int
}

//Expect sets up expected params for the NodeKeeper.GetSparseUnsyncList
func (m *mNodeKeeperMockGetSparseUnsyncList) Expect(p int) *mNodeKeeperMockGetSparseUnsyncList {
	m.mockExpectations = &NodeKeeperMockGetSparseUnsyncListParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.GetSparseUnsyncList to return Return's arguments
func (m *mNodeKeeperMockGetSparseUnsyncList) Return(r network.UnsyncList) *NodeKeeperMock {
	m.mock.GetSparseUnsyncListFunc = func(p int) network.UnsyncList {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetSparseUnsyncList method
func (m *mNodeKeeperMockGetSparseUnsyncList) Set(f func(p int) (r network.UnsyncList)) *NodeKeeperMock {
	m.mock.GetSparseUnsyncListFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetSparseUnsyncList implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetSparseUnsyncList(p int) (r network.UnsyncList) {
	atomic.AddUint64(&m.GetSparseUnsyncListPreCounter, 1)
	defer atomic.AddUint64(&m.GetSparseUnsyncListCounter, 1)

	if m.GetSparseUnsyncListMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetSparseUnsyncListMock.mockExpectations, NodeKeeperMockGetSparseUnsyncListParams{p},
			"NodeKeeper.GetSparseUnsyncList got unexpected parameters")

		if m.GetSparseUnsyncListFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.GetSparseUnsyncList")

			return
		}
	}

	if m.GetSparseUnsyncListFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetSparseUnsyncList")
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

type mNodeKeeperMockGetState struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.GetState to return Return's arguments
func (m *mNodeKeeperMockGetState) Return(r network.NodeKeeperState) *NodeKeeperMock {
	m.mock.GetStateFunc = func() network.NodeKeeperState {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetState method
func (m *mNodeKeeperMockGetState) Set(f func() (r network.NodeKeeperState)) *NodeKeeperMock {
	m.mock.GetStateFunc = f

	return m.mock
}

//GetState implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetState() (r network.NodeKeeperState) {
	atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if m.GetStateFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetState")
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

type mNodeKeeperMockGetUnsyncList struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.GetUnsyncList to return Return's arguments
func (m *mNodeKeeperMockGetUnsyncList) Return(r network.UnsyncList) *NodeKeeperMock {
	m.mock.GetUnsyncListFunc = func() network.UnsyncList {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.GetUnsyncList method
func (m *mNodeKeeperMockGetUnsyncList) Set(f func() (r network.UnsyncList)) *NodeKeeperMock {
	m.mock.GetUnsyncListFunc = f

	return m.mock
}

//GetUnsyncList implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetUnsyncList() (r network.UnsyncList) {
	atomic.AddUint64(&m.GetUnsyncListPreCounter, 1)
	defer atomic.AddUint64(&m.GetUnsyncListCounter, 1)

	if m.GetUnsyncListFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.GetUnsyncList")
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

type mNodeKeeperMockMoveSyncToActive struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.MoveSyncToActive to return Return's arguments
func (m *mNodeKeeperMockMoveSyncToActive) Return() *NodeKeeperMock {
	m.mock.MoveSyncToActiveFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.MoveSyncToActive method
func (m *mNodeKeeperMockMoveSyncToActive) Set(f func()) *NodeKeeperMock {
	m.mock.MoveSyncToActiveFunc = f

	return m.mock
}

//MoveSyncToActive implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) MoveSyncToActive() {
	atomic.AddUint64(&m.MoveSyncToActivePreCounter, 1)
	defer atomic.AddUint64(&m.MoveSyncToActiveCounter, 1)

	if m.MoveSyncToActiveFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.MoveSyncToActive")
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

type mNodeKeeperMockNodesJoinedDuringPreviousPulse struct {
	mock *NodeKeeperMock
}

//Return sets up a mock for NodeKeeper.NodesJoinedDuringPreviousPulse to return Return's arguments
func (m *mNodeKeeperMockNodesJoinedDuringPreviousPulse) Return(r bool) *NodeKeeperMock {
	m.mock.NodesJoinedDuringPreviousPulseFunc = func() bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.NodesJoinedDuringPreviousPulse method
func (m *mNodeKeeperMockNodesJoinedDuringPreviousPulse) Set(f func() (r bool)) *NodeKeeperMock {
	m.mock.NodesJoinedDuringPreviousPulseFunc = f

	return m.mock
}

//NodesJoinedDuringPreviousPulse implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) NodesJoinedDuringPreviousPulse() (r bool) {
	atomic.AddUint64(&m.NodesJoinedDuringPreviousPulsePreCounter, 1)
	defer atomic.AddUint64(&m.NodesJoinedDuringPreviousPulseCounter, 1)

	if m.NodesJoinedDuringPreviousPulseFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse")
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

type mNodeKeeperMockSetCloudHash struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockSetCloudHashParams
}

//NodeKeeperMockSetCloudHashParams represents input parameters of the NodeKeeper.SetCloudHash
type NodeKeeperMockSetCloudHashParams struct {
	p []byte
}

//Expect sets up expected params for the NodeKeeper.SetCloudHash
func (m *mNodeKeeperMockSetCloudHash) Expect(p []byte) *mNodeKeeperMockSetCloudHash {
	m.mockExpectations = &NodeKeeperMockSetCloudHashParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.SetCloudHash to return Return's arguments
func (m *mNodeKeeperMockSetCloudHash) Return() *NodeKeeperMock {
	m.mock.SetCloudHashFunc = func(p []byte) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.SetCloudHash method
func (m *mNodeKeeperMockSetCloudHash) Set(f func(p []byte)) *NodeKeeperMock {
	m.mock.SetCloudHashFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetCloudHash implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) SetCloudHash(p []byte) {
	atomic.AddUint64(&m.SetCloudHashPreCounter, 1)
	defer atomic.AddUint64(&m.SetCloudHashCounter, 1)

	if m.SetCloudHashMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetCloudHashMock.mockExpectations, NodeKeeperMockSetCloudHashParams{p},
			"NodeKeeper.SetCloudHash got unexpected parameters")

		if m.SetCloudHashFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.SetCloudHash")

			return
		}
	}

	if m.SetCloudHashFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.SetCloudHash")
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

type mNodeKeeperMockSetOriginClaim struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockSetOriginClaimParams
}

//NodeKeeperMockSetOriginClaimParams represents input parameters of the NodeKeeper.SetOriginClaim
type NodeKeeperMockSetOriginClaimParams struct {
	p *packets.NodeJoinClaim
}

//Expect sets up expected params for the NodeKeeper.SetOriginClaim
func (m *mNodeKeeperMockSetOriginClaim) Expect(p *packets.NodeJoinClaim) *mNodeKeeperMockSetOriginClaim {
	m.mockExpectations = &NodeKeeperMockSetOriginClaimParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.SetOriginClaim to return Return's arguments
func (m *mNodeKeeperMockSetOriginClaim) Return() *NodeKeeperMock {
	m.mock.SetOriginClaimFunc = func(p *packets.NodeJoinClaim) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.SetOriginClaim method
func (m *mNodeKeeperMockSetOriginClaim) Set(f func(p *packets.NodeJoinClaim)) *NodeKeeperMock {
	m.mock.SetOriginClaimFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetOriginClaim implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) SetOriginClaim(p *packets.NodeJoinClaim) {
	atomic.AddUint64(&m.SetOriginClaimPreCounter, 1)
	defer atomic.AddUint64(&m.SetOriginClaimCounter, 1)

	if m.SetOriginClaimMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetOriginClaimMock.mockExpectations, NodeKeeperMockSetOriginClaimParams{p},
			"NodeKeeper.SetOriginClaim got unexpected parameters")

		if m.SetOriginClaimFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.SetOriginClaim")

			return
		}
	}

	if m.SetOriginClaimFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.SetOriginClaim")
		return
	}

	m.SetOriginClaimFunc(p)
}

//SetOriginClaimMinimockCounter returns a count of NodeKeeperMock.SetOriginClaimFunc invocations
func (m *NodeKeeperMock) SetOriginClaimMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetOriginClaimCounter)
}

//SetOriginClaimMinimockPreCounter returns the value of NodeKeeperMock.SetOriginClaim invocations
func (m *NodeKeeperMock) SetOriginClaimMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetOriginClaimPreCounter)
}

type mNodeKeeperMockSetState struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockSetStateParams
}

//NodeKeeperMockSetStateParams represents input parameters of the NodeKeeper.SetState
type NodeKeeperMockSetStateParams struct {
	p network.NodeKeeperState
}

//Expect sets up expected params for the NodeKeeper.SetState
func (m *mNodeKeeperMockSetState) Expect(p network.NodeKeeperState) *mNodeKeeperMockSetState {
	m.mockExpectations = &NodeKeeperMockSetStateParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.SetState to return Return's arguments
func (m *mNodeKeeperMockSetState) Return() *NodeKeeperMock {
	m.mock.SetStateFunc = func(p network.NodeKeeperState) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.SetState method
func (m *mNodeKeeperMockSetState) Set(f func(p network.NodeKeeperState)) *NodeKeeperMock {
	m.mock.SetStateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetState implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) SetState(p network.NodeKeeperState) {
	atomic.AddUint64(&m.SetStatePreCounter, 1)
	defer atomic.AddUint64(&m.SetStateCounter, 1)

	if m.SetStateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetStateMock.mockExpectations, NodeKeeperMockSetStateParams{p},
			"NodeKeeper.SetState got unexpected parameters")

		if m.SetStateFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.SetState")

			return
		}
	}

	if m.SetStateFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.SetState")
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

type mNodeKeeperMockSync struct {
	mock             *NodeKeeperMock
	mockExpectations *NodeKeeperMockSyncParams
}

//NodeKeeperMockSyncParams represents input parameters of the NodeKeeper.Sync
type NodeKeeperMockSyncParams struct {
	p network.UnsyncList
}

//Expect sets up expected params for the NodeKeeper.Sync
func (m *mNodeKeeperMockSync) Expect(p network.UnsyncList) *mNodeKeeperMockSync {
	m.mockExpectations = &NodeKeeperMockSyncParams{p}
	return m
}

//Return sets up a mock for NodeKeeper.Sync to return Return's arguments
func (m *mNodeKeeperMockSync) Return() *NodeKeeperMock {
	m.mock.SyncFunc = func(p network.UnsyncList) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of NodeKeeper.Sync method
func (m *mNodeKeeperMockSync) Set(f func(p network.UnsyncList)) *NodeKeeperMock {
	m.mock.SyncFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Sync implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) Sync(p network.UnsyncList) {
	atomic.AddUint64(&m.SyncPreCounter, 1)
	defer atomic.AddUint64(&m.SyncCounter, 1)

	if m.SyncMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SyncMock.mockExpectations, NodeKeeperMockSyncParams{p},
			"NodeKeeper.Sync got unexpected parameters")

		if m.SyncFunc == nil {

			m.t.Fatal("No results are set for the NodeKeeperMock.Sync")

			return
		}
	}

	if m.SyncFunc == nil {
		m.t.Fatal("Unexpected call to NodeKeeperMock.Sync")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeKeeperMock) ValidateCallCounters() {

	if m.AddActiveNodesFunc != nil && atomic.LoadUint64(&m.AddActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.AddActiveNodes")
	}

	if m.AddPendingClaimFunc != nil && atomic.LoadUint64(&m.AddPendingClaimCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.AddPendingClaim")
	}

	if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNode")
	}

	if m.GetActiveNodeByShortIDFunc != nil && atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodeByShortID")
	}

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodes")
	}

	if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodesByRole")
	}

	if m.GetClaimQueueFunc != nil && atomic.LoadUint64(&m.GetClaimQueueCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetClaimQueue")
	}

	if m.GetCloudHashFunc != nil && atomic.LoadUint64(&m.GetCloudHashCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetCloudHash")
	}

	if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOrigin")
	}

	if m.GetOriginClaimFunc != nil && atomic.LoadUint64(&m.GetOriginClaimCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginClaim")
	}

	if m.GetSparseUnsyncListFunc != nil && atomic.LoadUint64(&m.GetSparseUnsyncListCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetSparseUnsyncList")
	}

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetState")
	}

	if m.GetUnsyncListFunc != nil && atomic.LoadUint64(&m.GetUnsyncListCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetUnsyncList")
	}

	if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.MoveSyncToActive")
	}

	if m.NodesJoinedDuringPreviousPulseFunc != nil && atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse")
	}

	if m.SetCloudHashFunc != nil && atomic.LoadUint64(&m.SetCloudHashCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.SetCloudHash")
	}

	if m.SetOriginClaimFunc != nil && atomic.LoadUint64(&m.SetOriginClaimCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.SetOriginClaim")
	}

	if m.SetStateFunc != nil && atomic.LoadUint64(&m.SetStateCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.SetState")
	}

	if m.SyncFunc != nil && atomic.LoadUint64(&m.SyncCounter) == 0 {
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

	if m.AddActiveNodesFunc != nil && atomic.LoadUint64(&m.AddActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.AddActiveNodes")
	}

	if m.AddPendingClaimFunc != nil && atomic.LoadUint64(&m.AddPendingClaimCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.AddPendingClaim")
	}

	if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNode")
	}

	if m.GetActiveNodeByShortIDFunc != nil && atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodeByShortID")
	}

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodes")
	}

	if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetActiveNodesByRole")
	}

	if m.GetClaimQueueFunc != nil && atomic.LoadUint64(&m.GetClaimQueueCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetClaimQueue")
	}

	if m.GetCloudHashFunc != nil && atomic.LoadUint64(&m.GetCloudHashCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetCloudHash")
	}

	if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOrigin")
	}

	if m.GetOriginClaimFunc != nil && atomic.LoadUint64(&m.GetOriginClaimCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginClaim")
	}

	if m.GetSparseUnsyncListFunc != nil && atomic.LoadUint64(&m.GetSparseUnsyncListCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetSparseUnsyncList")
	}

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetState")
	}

	if m.GetUnsyncListFunc != nil && atomic.LoadUint64(&m.GetUnsyncListCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.GetUnsyncList")
	}

	if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.MoveSyncToActive")
	}

	if m.NodesJoinedDuringPreviousPulseFunc != nil && atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse")
	}

	if m.SetCloudHashFunc != nil && atomic.LoadUint64(&m.SetCloudHashCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.SetCloudHash")
	}

	if m.SetOriginClaimFunc != nil && atomic.LoadUint64(&m.SetOriginClaimCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.SetOriginClaim")
	}

	if m.SetStateFunc != nil && atomic.LoadUint64(&m.SetStateCounter) == 0 {
		m.t.Fatal("Expected call to NodeKeeperMock.SetState")
	}

	if m.SyncFunc != nil && atomic.LoadUint64(&m.SyncCounter) == 0 {
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
		ok = ok && (m.AddActiveNodesFunc == nil || atomic.LoadUint64(&m.AddActiveNodesCounter) > 0)
		ok = ok && (m.AddPendingClaimFunc == nil || atomic.LoadUint64(&m.AddPendingClaimCounter) > 0)
		ok = ok && (m.GetActiveNodeFunc == nil || atomic.LoadUint64(&m.GetActiveNodeCounter) > 0)
		ok = ok && (m.GetActiveNodeByShortIDFunc == nil || atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) > 0)
		ok = ok && (m.GetActiveNodesFunc == nil || atomic.LoadUint64(&m.GetActiveNodesCounter) > 0)
		ok = ok && (m.GetActiveNodesByRoleFunc == nil || atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) > 0)
		ok = ok && (m.GetClaimQueueFunc == nil || atomic.LoadUint64(&m.GetClaimQueueCounter) > 0)
		ok = ok && (m.GetCloudHashFunc == nil || atomic.LoadUint64(&m.GetCloudHashCounter) > 0)
		ok = ok && (m.GetOriginFunc == nil || atomic.LoadUint64(&m.GetOriginCounter) > 0)
		ok = ok && (m.GetOriginClaimFunc == nil || atomic.LoadUint64(&m.GetOriginClaimCounter) > 0)
		ok = ok && (m.GetSparseUnsyncListFunc == nil || atomic.LoadUint64(&m.GetSparseUnsyncListCounter) > 0)
		ok = ok && (m.GetStateFunc == nil || atomic.LoadUint64(&m.GetStateCounter) > 0)
		ok = ok && (m.GetUnsyncListFunc == nil || atomic.LoadUint64(&m.GetUnsyncListCounter) > 0)
		ok = ok && (m.MoveSyncToActiveFunc == nil || atomic.LoadUint64(&m.MoveSyncToActiveCounter) > 0)
		ok = ok && (m.NodesJoinedDuringPreviousPulseFunc == nil || atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) > 0)
		ok = ok && (m.SetCloudHashFunc == nil || atomic.LoadUint64(&m.SetCloudHashCounter) > 0)
		ok = ok && (m.SetOriginClaimFunc == nil || atomic.LoadUint64(&m.SetOriginClaimCounter) > 0)
		ok = ok && (m.SetStateFunc == nil || atomic.LoadUint64(&m.SetStateCounter) > 0)
		ok = ok && (m.SyncFunc == nil || atomic.LoadUint64(&m.SyncCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.AddActiveNodesFunc != nil && atomic.LoadUint64(&m.AddActiveNodesCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.AddActiveNodes")
			}

			if m.AddPendingClaimFunc != nil && atomic.LoadUint64(&m.AddPendingClaimCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.AddPendingClaim")
			}

			if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNode")
			}

			if m.GetActiveNodeByShortIDFunc != nil && atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNodeByShortID")
			}

			if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNodes")
			}

			if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetActiveNodesByRole")
			}

			if m.GetClaimQueueFunc != nil && atomic.LoadUint64(&m.GetClaimQueueCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetClaimQueue")
			}

			if m.GetCloudHashFunc != nil && atomic.LoadUint64(&m.GetCloudHashCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetCloudHash")
			}

			if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetOrigin")
			}

			if m.GetOriginClaimFunc != nil && atomic.LoadUint64(&m.GetOriginClaimCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetOriginClaim")
			}

			if m.GetSparseUnsyncListFunc != nil && atomic.LoadUint64(&m.GetSparseUnsyncListCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetSparseUnsyncList")
			}

			if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetState")
			}

			if m.GetUnsyncListFunc != nil && atomic.LoadUint64(&m.GetUnsyncListCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.GetUnsyncList")
			}

			if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.MoveSyncToActive")
			}

			if m.NodesJoinedDuringPreviousPulseFunc != nil && atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.NodesJoinedDuringPreviousPulse")
			}

			if m.SetCloudHashFunc != nil && atomic.LoadUint64(&m.SetCloudHashCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.SetCloudHash")
			}

			if m.SetOriginClaimFunc != nil && atomic.LoadUint64(&m.SetOriginClaimCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.SetOriginClaim")
			}

			if m.SetStateFunc != nil && atomic.LoadUint64(&m.SetStateCounter) == 0 {
				m.t.Error("Expected call to NodeKeeperMock.SetState")
			}

			if m.SyncFunc != nil && atomic.LoadUint64(&m.SyncCounter) == 0 {
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

	if m.AddActiveNodesFunc != nil && atomic.LoadUint64(&m.AddActiveNodesCounter) == 0 {
		return false
	}

	if m.AddPendingClaimFunc != nil && atomic.LoadUint64(&m.AddPendingClaimCounter) == 0 {
		return false
	}

	if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
		return false
	}

	if m.GetActiveNodeByShortIDFunc != nil && atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) == 0 {
		return false
	}

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		return false
	}

	if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
		return false
	}

	if m.GetClaimQueueFunc != nil && atomic.LoadUint64(&m.GetClaimQueueCounter) == 0 {
		return false
	}

	if m.GetCloudHashFunc != nil && atomic.LoadUint64(&m.GetCloudHashCounter) == 0 {
		return false
	}

	if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
		return false
	}

	if m.GetOriginClaimFunc != nil && atomic.LoadUint64(&m.GetOriginClaimCounter) == 0 {
		return false
	}

	if m.GetSparseUnsyncListFunc != nil && atomic.LoadUint64(&m.GetSparseUnsyncListCounter) == 0 {
		return false
	}

	if m.GetStateFunc != nil && atomic.LoadUint64(&m.GetStateCounter) == 0 {
		return false
	}

	if m.GetUnsyncListFunc != nil && atomic.LoadUint64(&m.GetUnsyncListCounter) == 0 {
		return false
	}

	if m.MoveSyncToActiveFunc != nil && atomic.LoadUint64(&m.MoveSyncToActiveCounter) == 0 {
		return false
	}

	if m.NodesJoinedDuringPreviousPulseFunc != nil && atomic.LoadUint64(&m.NodesJoinedDuringPreviousPulseCounter) == 0 {
		return false
	}

	if m.SetCloudHashFunc != nil && atomic.LoadUint64(&m.SetCloudHashCounter) == 0 {
		return false
	}

	if m.SetOriginClaimFunc != nil && atomic.LoadUint64(&m.SetOriginClaimCounter) == 0 {
		return false
	}

	if m.SetStateFunc != nil && atomic.LoadUint64(&m.SetStateCounter) == 0 {
		return false
	}

	if m.SyncFunc != nil && atomic.LoadUint64(&m.SyncCounter) == 0 {
		return false
	}

	return true
}

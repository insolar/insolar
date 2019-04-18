package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeKeeper" can be found in github.com/insolar/insolar/network
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	packets "github.com/insolar/insolar/consensus/packets"
	insolar "github.com/insolar/insolar/insolar"
	network "github.com/insolar/insolar/network"
	node "github.com/insolar/insolar/network/node"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeKeeperMock implements github.com/insolar/insolar/network.NodeKeeper
type NodeKeeperMock struct {
	t minimock.Tester

	GetAccessorFunc       func() (r network.Accessor)
	GetAccessorCounter    uint64
	GetAccessorPreCounter uint64
	GetAccessorMock       mNodeKeeperMockGetAccessor

	GetClaimQueueFunc       func() (r network.ClaimQueue)
	GetClaimQueueCounter    uint64
	GetClaimQueuePreCounter uint64
	GetClaimQueueMock       mNodeKeeperMockGetClaimQueue

	GetCloudHashFunc       func() (r []byte)
	GetCloudHashCounter    uint64
	GetCloudHashPreCounter uint64
	GetCloudHashMock       mNodeKeeperMockGetCloudHash

	GetConsensusInfoFunc       func() (r network.ConsensusInfo)
	GetConsensusInfoCounter    uint64
	GetConsensusInfoPreCounter uint64
	GetConsensusInfoMock       mNodeKeeperMockGetConsensusInfo

	GetOriginFunc       func() (r insolar.NetworkNode)
	GetOriginCounter    uint64
	GetOriginPreCounter uint64
	GetOriginMock       mNodeKeeperMockGetOrigin

	GetOriginAnnounceClaimFunc       func(p packets.BitSetMapper) (r *packets.NodeAnnounceClaim, r1 error)
	GetOriginAnnounceClaimCounter    uint64
	GetOriginAnnounceClaimPreCounter uint64
	GetOriginAnnounceClaimMock       mNodeKeeperMockGetOriginAnnounceClaim

	GetOriginJoinClaimFunc       func() (r *packets.NodeJoinClaim, r1 error)
	GetOriginJoinClaimCounter    uint64
	GetOriginJoinClaimPreCounter uint64
	GetOriginJoinClaimMock       mNodeKeeperMockGetOriginJoinClaim

	GetSnapshotCopyFunc       func() (r *node.Snapshot)
	GetSnapshotCopyCounter    uint64
	GetSnapshotCopyPreCounter uint64
	GetSnapshotCopyMock       mNodeKeeperMockGetSnapshotCopy

	GetWorkingNodeFunc       func(p insolar.Reference) (r insolar.NetworkNode)
	GetWorkingNodeCounter    uint64
	GetWorkingNodePreCounter uint64
	GetWorkingNodeMock       mNodeKeeperMockGetWorkingNode

	GetWorkingNodesFunc       func() (r []insolar.NetworkNode)
	GetWorkingNodesCounter    uint64
	GetWorkingNodesPreCounter uint64
	GetWorkingNodesMock       mNodeKeeperMockGetWorkingNodes

	GetWorkingNodesByRoleFunc       func(p insolar.DynamicRole) (r []insolar.Reference)
	GetWorkingNodesByRoleCounter    uint64
	GetWorkingNodesByRolePreCounter uint64
	GetWorkingNodesByRoleMock       mNodeKeeperMockGetWorkingNodesByRole

	IsBootstrappedFunc       func() (r bool)
	IsBootstrappedCounter    uint64
	IsBootstrappedPreCounter uint64
	IsBootstrappedMock       mNodeKeeperMockIsBootstrapped

	MoveSyncToActiveFunc       func(p context.Context, p1 insolar.PulseNumber) (r error)
	MoveSyncToActiveCounter    uint64
	MoveSyncToActivePreCounter uint64
	MoveSyncToActiveMock       mNodeKeeperMockMoveSyncToActive

	SetCloudHashFunc       func(p []byte)
	SetCloudHashCounter    uint64
	SetCloudHashPreCounter uint64
	SetCloudHashMock       mNodeKeeperMockSetCloudHash

	SetInitialSnapshotFunc       func(p []insolar.NetworkNode)
	SetInitialSnapshotCounter    uint64
	SetInitialSnapshotPreCounter uint64
	SetInitialSnapshotMock       mNodeKeeperMockSetInitialSnapshot

	SetIsBootstrappedFunc       func(p bool)
	SetIsBootstrappedCounter    uint64
	SetIsBootstrappedPreCounter uint64
	SetIsBootstrappedMock       mNodeKeeperMockSetIsBootstrapped

	SyncFunc       func(p context.Context, p1 []insolar.NetworkNode, p2 []packets.ReferendumClaim) (r error)
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

	m.GetAccessorMock = mNodeKeeperMockGetAccessor{mock: m}
	m.GetClaimQueueMock = mNodeKeeperMockGetClaimQueue{mock: m}
	m.GetCloudHashMock = mNodeKeeperMockGetCloudHash{mock: m}
	m.GetConsensusInfoMock = mNodeKeeperMockGetConsensusInfo{mock: m}
	m.GetOriginMock = mNodeKeeperMockGetOrigin{mock: m}
	m.GetOriginAnnounceClaimMock = mNodeKeeperMockGetOriginAnnounceClaim{mock: m}
	m.GetOriginJoinClaimMock = mNodeKeeperMockGetOriginJoinClaim{mock: m}
	m.GetSnapshotCopyMock = mNodeKeeperMockGetSnapshotCopy{mock: m}
	m.GetWorkingNodeMock = mNodeKeeperMockGetWorkingNode{mock: m}
	m.GetWorkingNodesMock = mNodeKeeperMockGetWorkingNodes{mock: m}
	m.GetWorkingNodesByRoleMock = mNodeKeeperMockGetWorkingNodesByRole{mock: m}
	m.IsBootstrappedMock = mNodeKeeperMockIsBootstrapped{mock: m}
	m.MoveSyncToActiveMock = mNodeKeeperMockMoveSyncToActive{mock: m}
	m.SetCloudHashMock = mNodeKeeperMockSetCloudHash{mock: m}
	m.SetInitialSnapshotMock = mNodeKeeperMockSetInitialSnapshot{mock: m}
	m.SetIsBootstrappedMock = mNodeKeeperMockSetIsBootstrapped{mock: m}
	m.SyncMock = mNodeKeeperMockSync{mock: m}

	return m
}

type mNodeKeeperMockGetAccessor struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetAccessorExpectation
	expectationSeries []*NodeKeeperMockGetAccessorExpectation
}

type NodeKeeperMockGetAccessorExpectation struct {
	result *NodeKeeperMockGetAccessorResult
}

type NodeKeeperMockGetAccessorResult struct {
	r network.Accessor
}

//Expect specifies that invocation of NodeKeeper.GetAccessor is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetAccessor) Expect() *mNodeKeeperMockGetAccessor {
	m.mock.GetAccessorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetAccessorExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetAccessor
func (m *mNodeKeeperMockGetAccessor) Return(r network.Accessor) *NodeKeeperMock {
	m.mock.GetAccessorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetAccessorExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetAccessorResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetAccessor is expected once
func (m *mNodeKeeperMockGetAccessor) ExpectOnce() *NodeKeeperMockGetAccessorExpectation {
	m.mock.GetAccessorFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetAccessorExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetAccessorExpectation) Return(r network.Accessor) {
	e.result = &NodeKeeperMockGetAccessorResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetAccessor method
func (m *mNodeKeeperMockGetAccessor) Set(f func() (r network.Accessor)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetAccessorFunc = f
	return m.mock
}

//GetAccessor implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetAccessor() (r network.Accessor) {
	counter := atomic.AddUint64(&m.GetAccessorPreCounter, 1)
	defer atomic.AddUint64(&m.GetAccessorCounter, 1)

	if len(m.GetAccessorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetAccessorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetAccessor.")
			return
		}

		result := m.GetAccessorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetAccessor")
			return
		}

		r = result.r

		return
	}

	if m.GetAccessorMock.mainExpectation != nil {

		result := m.GetAccessorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetAccessor")
		}

		r = result.r

		return
	}

	if m.GetAccessorFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetAccessor.")
		return
	}

	return m.GetAccessorFunc()
}

//GetAccessorMinimockCounter returns a count of NodeKeeperMock.GetAccessorFunc invocations
func (m *NodeKeeperMock) GetAccessorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAccessorCounter)
}

//GetAccessorMinimockPreCounter returns the value of NodeKeeperMock.GetAccessor invocations
func (m *NodeKeeperMock) GetAccessorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAccessorPreCounter)
}

//GetAccessorFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetAccessorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetAccessorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetAccessorCounter) == uint64(len(m.GetAccessorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetAccessorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetAccessorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetAccessorFunc != nil {
		return atomic.LoadUint64(&m.GetAccessorCounter) > 0
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

type mNodeKeeperMockGetConsensusInfo struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetConsensusInfoExpectation
	expectationSeries []*NodeKeeperMockGetConsensusInfoExpectation
}

type NodeKeeperMockGetConsensusInfoExpectation struct {
	result *NodeKeeperMockGetConsensusInfoResult
}

type NodeKeeperMockGetConsensusInfoResult struct {
	r network.ConsensusInfo
}

//Expect specifies that invocation of NodeKeeper.GetConsensusInfo is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetConsensusInfo) Expect() *mNodeKeeperMockGetConsensusInfo {
	m.mock.GetConsensusInfoFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetConsensusInfoExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetConsensusInfo
func (m *mNodeKeeperMockGetConsensusInfo) Return(r network.ConsensusInfo) *NodeKeeperMock {
	m.mock.GetConsensusInfoFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetConsensusInfoExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetConsensusInfoResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetConsensusInfo is expected once
func (m *mNodeKeeperMockGetConsensusInfo) ExpectOnce() *NodeKeeperMockGetConsensusInfoExpectation {
	m.mock.GetConsensusInfoFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetConsensusInfoExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetConsensusInfoExpectation) Return(r network.ConsensusInfo) {
	e.result = &NodeKeeperMockGetConsensusInfoResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetConsensusInfo method
func (m *mNodeKeeperMockGetConsensusInfo) Set(f func() (r network.ConsensusInfo)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetConsensusInfoFunc = f
	return m.mock
}

//GetConsensusInfo implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetConsensusInfo() (r network.ConsensusInfo) {
	counter := atomic.AddUint64(&m.GetConsensusInfoPreCounter, 1)
	defer atomic.AddUint64(&m.GetConsensusInfoCounter, 1)

	if len(m.GetConsensusInfoMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetConsensusInfoMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetConsensusInfo.")
			return
		}

		result := m.GetConsensusInfoMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetConsensusInfo")
			return
		}

		r = result.r

		return
	}

	if m.GetConsensusInfoMock.mainExpectation != nil {

		result := m.GetConsensusInfoMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetConsensusInfo")
		}

		r = result.r

		return
	}

	if m.GetConsensusInfoFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetConsensusInfo.")
		return
	}

	return m.GetConsensusInfoFunc()
}

//GetConsensusInfoMinimockCounter returns a count of NodeKeeperMock.GetConsensusInfoFunc invocations
func (m *NodeKeeperMock) GetConsensusInfoMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetConsensusInfoCounter)
}

//GetConsensusInfoMinimockPreCounter returns the value of NodeKeeperMock.GetConsensusInfo invocations
func (m *NodeKeeperMock) GetConsensusInfoMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetConsensusInfoPreCounter)
}

//GetConsensusInfoFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetConsensusInfoFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetConsensusInfoMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetConsensusInfoCounter) == uint64(len(m.GetConsensusInfoMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetConsensusInfoMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetConsensusInfoCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetConsensusInfoFunc != nil {
		return atomic.LoadUint64(&m.GetConsensusInfoCounter) > 0
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
	r insolar.NetworkNode
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
func (m *mNodeKeeperMockGetOrigin) Return(r insolar.NetworkNode) *NodeKeeperMock {
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

func (e *NodeKeeperMockGetOriginExpectation) Return(r insolar.NetworkNode) {
	e.result = &NodeKeeperMockGetOriginResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetOrigin method
func (m *mNodeKeeperMockGetOrigin) Set(f func() (r insolar.NetworkNode)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOriginFunc = f
	return m.mock
}

//GetOrigin implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetOrigin() (r insolar.NetworkNode) {
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

type mNodeKeeperMockGetOriginAnnounceClaim struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetOriginAnnounceClaimExpectation
	expectationSeries []*NodeKeeperMockGetOriginAnnounceClaimExpectation
}

type NodeKeeperMockGetOriginAnnounceClaimExpectation struct {
	input  *NodeKeeperMockGetOriginAnnounceClaimInput
	result *NodeKeeperMockGetOriginAnnounceClaimResult
}

type NodeKeeperMockGetOriginAnnounceClaimInput struct {
	p packets.BitSetMapper
}

type NodeKeeperMockGetOriginAnnounceClaimResult struct {
	r  *packets.NodeAnnounceClaim
	r1 error
}

//Expect specifies that invocation of NodeKeeper.GetOriginAnnounceClaim is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetOriginAnnounceClaim) Expect(p packets.BitSetMapper) *mNodeKeeperMockGetOriginAnnounceClaim {
	m.mock.GetOriginAnnounceClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginAnnounceClaimExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockGetOriginAnnounceClaimInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.GetOriginAnnounceClaim
func (m *mNodeKeeperMockGetOriginAnnounceClaim) Return(r *packets.NodeAnnounceClaim, r1 error) *NodeKeeperMock {
	m.mock.GetOriginAnnounceClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginAnnounceClaimExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetOriginAnnounceClaimResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetOriginAnnounceClaim is expected once
func (m *mNodeKeeperMockGetOriginAnnounceClaim) ExpectOnce(p packets.BitSetMapper) *NodeKeeperMockGetOriginAnnounceClaimExpectation {
	m.mock.GetOriginAnnounceClaimFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetOriginAnnounceClaimExpectation{}
	expectation.input = &NodeKeeperMockGetOriginAnnounceClaimInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetOriginAnnounceClaimExpectation) Return(r *packets.NodeAnnounceClaim, r1 error) {
	e.result = &NodeKeeperMockGetOriginAnnounceClaimResult{r, r1}
}

//Set uses given function f as a mock of NodeKeeper.GetOriginAnnounceClaim method
func (m *mNodeKeeperMockGetOriginAnnounceClaim) Set(f func(p packets.BitSetMapper) (r *packets.NodeAnnounceClaim, r1 error)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOriginAnnounceClaimFunc = f
	return m.mock
}

//GetOriginAnnounceClaim implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetOriginAnnounceClaim(p packets.BitSetMapper) (r *packets.NodeAnnounceClaim, r1 error) {
	counter := atomic.AddUint64(&m.GetOriginAnnounceClaimPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginAnnounceClaimCounter, 1)

	if len(m.GetOriginAnnounceClaimMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOriginAnnounceClaimMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOriginAnnounceClaim. %v", p)
			return
		}

		input := m.GetOriginAnnounceClaimMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockGetOriginAnnounceClaimInput{p}, "NodeKeeper.GetOriginAnnounceClaim got unexpected parameters")

		result := m.GetOriginAnnounceClaimMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOriginAnnounceClaim")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetOriginAnnounceClaimMock.mainExpectation != nil {

		input := m.GetOriginAnnounceClaimMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockGetOriginAnnounceClaimInput{p}, "NodeKeeper.GetOriginAnnounceClaim got unexpected parameters")
		}

		result := m.GetOriginAnnounceClaimMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOriginAnnounceClaim")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetOriginAnnounceClaimFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOriginAnnounceClaim. %v", p)
		return
	}

	return m.GetOriginAnnounceClaimFunc(p)
}

//GetOriginAnnounceClaimMinimockCounter returns a count of NodeKeeperMock.GetOriginAnnounceClaimFunc invocations
func (m *NodeKeeperMock) GetOriginAnnounceClaimMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginAnnounceClaimCounter)
}

//GetOriginAnnounceClaimMinimockPreCounter returns the value of NodeKeeperMock.GetOriginAnnounceClaim invocations
func (m *NodeKeeperMock) GetOriginAnnounceClaimMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginAnnounceClaimPreCounter)
}

//GetOriginAnnounceClaimFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetOriginAnnounceClaimFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOriginAnnounceClaimMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOriginAnnounceClaimCounter) == uint64(len(m.GetOriginAnnounceClaimMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOriginAnnounceClaimMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOriginAnnounceClaimCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOriginAnnounceClaimFunc != nil {
		return atomic.LoadUint64(&m.GetOriginAnnounceClaimCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetOriginJoinClaim struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetOriginJoinClaimExpectation
	expectationSeries []*NodeKeeperMockGetOriginJoinClaimExpectation
}

type NodeKeeperMockGetOriginJoinClaimExpectation struct {
	result *NodeKeeperMockGetOriginJoinClaimResult
}

type NodeKeeperMockGetOriginJoinClaimResult struct {
	r  *packets.NodeJoinClaim
	r1 error
}

//Expect specifies that invocation of NodeKeeper.GetOriginJoinClaim is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetOriginJoinClaim) Expect() *mNodeKeeperMockGetOriginJoinClaim {
	m.mock.GetOriginJoinClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginJoinClaimExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetOriginJoinClaim
func (m *mNodeKeeperMockGetOriginJoinClaim) Return(r *packets.NodeJoinClaim, r1 error) *NodeKeeperMock {
	m.mock.GetOriginJoinClaimFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetOriginJoinClaimExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetOriginJoinClaimResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetOriginJoinClaim is expected once
func (m *mNodeKeeperMockGetOriginJoinClaim) ExpectOnce() *NodeKeeperMockGetOriginJoinClaimExpectation {
	m.mock.GetOriginJoinClaimFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetOriginJoinClaimExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetOriginJoinClaimExpectation) Return(r *packets.NodeJoinClaim, r1 error) {
	e.result = &NodeKeeperMockGetOriginJoinClaimResult{r, r1}
}

//Set uses given function f as a mock of NodeKeeper.GetOriginJoinClaim method
func (m *mNodeKeeperMockGetOriginJoinClaim) Set(f func() (r *packets.NodeJoinClaim, r1 error)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOriginJoinClaimFunc = f
	return m.mock
}

//GetOriginJoinClaim implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetOriginJoinClaim() (r *packets.NodeJoinClaim, r1 error) {
	counter := atomic.AddUint64(&m.GetOriginJoinClaimPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginJoinClaimCounter, 1)

	if len(m.GetOriginJoinClaimMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOriginJoinClaimMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOriginJoinClaim.")
			return
		}

		result := m.GetOriginJoinClaimMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOriginJoinClaim")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetOriginJoinClaimMock.mainExpectation != nil {

		result := m.GetOriginJoinClaimMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetOriginJoinClaim")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetOriginJoinClaimFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetOriginJoinClaim.")
		return
	}

	return m.GetOriginJoinClaimFunc()
}

//GetOriginJoinClaimMinimockCounter returns a count of NodeKeeperMock.GetOriginJoinClaimFunc invocations
func (m *NodeKeeperMock) GetOriginJoinClaimMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginJoinClaimCounter)
}

//GetOriginJoinClaimMinimockPreCounter returns the value of NodeKeeperMock.GetOriginJoinClaim invocations
func (m *NodeKeeperMock) GetOriginJoinClaimMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginJoinClaimPreCounter)
}

//GetOriginJoinClaimFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetOriginJoinClaimFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOriginJoinClaimMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOriginJoinClaimCounter) == uint64(len(m.GetOriginJoinClaimMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOriginJoinClaimMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOriginJoinClaimCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOriginJoinClaimFunc != nil {
		return atomic.LoadUint64(&m.GetOriginJoinClaimCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetSnapshotCopy struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetSnapshotCopyExpectation
	expectationSeries []*NodeKeeperMockGetSnapshotCopyExpectation
}

type NodeKeeperMockGetSnapshotCopyExpectation struct {
	result *NodeKeeperMockGetSnapshotCopyResult
}

type NodeKeeperMockGetSnapshotCopyResult struct {
	r *node.Snapshot
}

//Expect specifies that invocation of NodeKeeper.GetSnapshotCopy is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetSnapshotCopy) Expect() *mNodeKeeperMockGetSnapshotCopy {
	m.mock.GetSnapshotCopyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetSnapshotCopyExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetSnapshotCopy
func (m *mNodeKeeperMockGetSnapshotCopy) Return(r *node.Snapshot) *NodeKeeperMock {
	m.mock.GetSnapshotCopyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetSnapshotCopyExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetSnapshotCopyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetSnapshotCopy is expected once
func (m *mNodeKeeperMockGetSnapshotCopy) ExpectOnce() *NodeKeeperMockGetSnapshotCopyExpectation {
	m.mock.GetSnapshotCopyFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetSnapshotCopyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetSnapshotCopyExpectation) Return(r *node.Snapshot) {
	e.result = &NodeKeeperMockGetSnapshotCopyResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetSnapshotCopy method
func (m *mNodeKeeperMockGetSnapshotCopy) Set(f func() (r *node.Snapshot)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSnapshotCopyFunc = f
	return m.mock
}

//GetSnapshotCopy implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetSnapshotCopy() (r *node.Snapshot) {
	counter := atomic.AddUint64(&m.GetSnapshotCopyPreCounter, 1)
	defer atomic.AddUint64(&m.GetSnapshotCopyCounter, 1)

	if len(m.GetSnapshotCopyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSnapshotCopyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetSnapshotCopy.")
			return
		}

		result := m.GetSnapshotCopyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetSnapshotCopy")
			return
		}

		r = result.r

		return
	}

	if m.GetSnapshotCopyMock.mainExpectation != nil {

		result := m.GetSnapshotCopyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetSnapshotCopy")
		}

		r = result.r

		return
	}

	if m.GetSnapshotCopyFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetSnapshotCopy.")
		return
	}

	return m.GetSnapshotCopyFunc()
}

//GetSnapshotCopyMinimockCounter returns a count of NodeKeeperMock.GetSnapshotCopyFunc invocations
func (m *NodeKeeperMock) GetSnapshotCopyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSnapshotCopyCounter)
}

//GetSnapshotCopyMinimockPreCounter returns the value of NodeKeeperMock.GetSnapshotCopy invocations
func (m *NodeKeeperMock) GetSnapshotCopyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSnapshotCopyPreCounter)
}

//GetSnapshotCopyFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetSnapshotCopyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSnapshotCopyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSnapshotCopyCounter) == uint64(len(m.GetSnapshotCopyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSnapshotCopyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSnapshotCopyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSnapshotCopyFunc != nil {
		return atomic.LoadUint64(&m.GetSnapshotCopyCounter) > 0
	}

	return true
}

type mNodeKeeperMockGetWorkingNode struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetWorkingNodeExpectation
	expectationSeries []*NodeKeeperMockGetWorkingNodeExpectation
}

type NodeKeeperMockGetWorkingNodeExpectation struct {
	input  *NodeKeeperMockGetWorkingNodeInput
	result *NodeKeeperMockGetWorkingNodeResult
}

type NodeKeeperMockGetWorkingNodeInput struct {
	p insolar.Reference
}

type NodeKeeperMockGetWorkingNodeResult struct {
	r insolar.NetworkNode
}

//Expect specifies that invocation of NodeKeeper.GetWorkingNode is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetWorkingNode) Expect(p insolar.Reference) *mNodeKeeperMockGetWorkingNode {
	m.mock.GetWorkingNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetWorkingNodeExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockGetWorkingNodeInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.GetWorkingNode
func (m *mNodeKeeperMockGetWorkingNode) Return(r insolar.NetworkNode) *NodeKeeperMock {
	m.mock.GetWorkingNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetWorkingNodeExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetWorkingNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetWorkingNode is expected once
func (m *mNodeKeeperMockGetWorkingNode) ExpectOnce(p insolar.Reference) *NodeKeeperMockGetWorkingNodeExpectation {
	m.mock.GetWorkingNodeFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetWorkingNodeExpectation{}
	expectation.input = &NodeKeeperMockGetWorkingNodeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetWorkingNodeExpectation) Return(r insolar.NetworkNode) {
	e.result = &NodeKeeperMockGetWorkingNodeResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetWorkingNode method
func (m *mNodeKeeperMockGetWorkingNode) Set(f func(p insolar.Reference) (r insolar.NetworkNode)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodeFunc = f
	return m.mock
}

//GetWorkingNode implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetWorkingNode(p insolar.Reference) (r insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetWorkingNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodeCounter, 1)

	if len(m.GetWorkingNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetWorkingNode. %v", p)
			return
		}

		input := m.GetWorkingNodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockGetWorkingNodeInput{p}, "NodeKeeper.GetWorkingNode got unexpected parameters")

		result := m.GetWorkingNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetWorkingNode")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodeMock.mainExpectation != nil {

		input := m.GetWorkingNodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockGetWorkingNodeInput{p}, "NodeKeeper.GetWorkingNode got unexpected parameters")
		}

		result := m.GetWorkingNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetWorkingNode")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodeFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetWorkingNode. %v", p)
		return
	}

	return m.GetWorkingNodeFunc(p)
}

//GetWorkingNodeMinimockCounter returns a count of NodeKeeperMock.GetWorkingNodeFunc invocations
func (m *NodeKeeperMock) GetWorkingNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodeCounter)
}

//GetWorkingNodeMinimockPreCounter returns the value of NodeKeeperMock.GetWorkingNode invocations
func (m *NodeKeeperMock) GetWorkingNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodePreCounter)
}

//GetWorkingNodeFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetWorkingNodeFinished() bool {
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

type mNodeKeeperMockGetWorkingNodes struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetWorkingNodesExpectation
	expectationSeries []*NodeKeeperMockGetWorkingNodesExpectation
}

type NodeKeeperMockGetWorkingNodesExpectation struct {
	result *NodeKeeperMockGetWorkingNodesResult
}

type NodeKeeperMockGetWorkingNodesResult struct {
	r []insolar.NetworkNode
}

//Expect specifies that invocation of NodeKeeper.GetWorkingNodes is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetWorkingNodes) Expect() *mNodeKeeperMockGetWorkingNodes {
	m.mock.GetWorkingNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetWorkingNodesExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeKeeper.GetWorkingNodes
func (m *mNodeKeeperMockGetWorkingNodes) Return(r []insolar.NetworkNode) *NodeKeeperMock {
	m.mock.GetWorkingNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetWorkingNodesExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetWorkingNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetWorkingNodes is expected once
func (m *mNodeKeeperMockGetWorkingNodes) ExpectOnce() *NodeKeeperMockGetWorkingNodesExpectation {
	m.mock.GetWorkingNodesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetWorkingNodesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetWorkingNodesExpectation) Return(r []insolar.NetworkNode) {
	e.result = &NodeKeeperMockGetWorkingNodesResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetWorkingNodes method
func (m *mNodeKeeperMockGetWorkingNodes) Set(f func() (r []insolar.NetworkNode)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodesFunc = f
	return m.mock
}

//GetWorkingNodes implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetWorkingNodes() (r []insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetWorkingNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodesCounter, 1)

	if len(m.GetWorkingNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetWorkingNodes.")
			return
		}

		result := m.GetWorkingNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetWorkingNodes")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesMock.mainExpectation != nil {

		result := m.GetWorkingNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetWorkingNodes")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetWorkingNodes.")
		return
	}

	return m.GetWorkingNodesFunc()
}

//GetWorkingNodesMinimockCounter returns a count of NodeKeeperMock.GetWorkingNodesFunc invocations
func (m *NodeKeeperMock) GetWorkingNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesCounter)
}

//GetWorkingNodesMinimockPreCounter returns the value of NodeKeeperMock.GetWorkingNodes invocations
func (m *NodeKeeperMock) GetWorkingNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesPreCounter)
}

//GetWorkingNodesFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetWorkingNodesFinished() bool {
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

type mNodeKeeperMockGetWorkingNodesByRole struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockGetWorkingNodesByRoleExpectation
	expectationSeries []*NodeKeeperMockGetWorkingNodesByRoleExpectation
}

type NodeKeeperMockGetWorkingNodesByRoleExpectation struct {
	input  *NodeKeeperMockGetWorkingNodesByRoleInput
	result *NodeKeeperMockGetWorkingNodesByRoleResult
}

type NodeKeeperMockGetWorkingNodesByRoleInput struct {
	p insolar.DynamicRole
}

type NodeKeeperMockGetWorkingNodesByRoleResult struct {
	r []insolar.Reference
}

//Expect specifies that invocation of NodeKeeper.GetWorkingNodesByRole is expected from 1 to Infinity times
func (m *mNodeKeeperMockGetWorkingNodesByRole) Expect(p insolar.DynamicRole) *mNodeKeeperMockGetWorkingNodesByRole {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetWorkingNodesByRoleExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockGetWorkingNodesByRoleInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.GetWorkingNodesByRole
func (m *mNodeKeeperMockGetWorkingNodesByRole) Return(r []insolar.Reference) *NodeKeeperMock {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockGetWorkingNodesByRoleExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockGetWorkingNodesByRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.GetWorkingNodesByRole is expected once
func (m *mNodeKeeperMockGetWorkingNodesByRole) ExpectOnce(p insolar.DynamicRole) *NodeKeeperMockGetWorkingNodesByRoleExpectation {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockGetWorkingNodesByRoleExpectation{}
	expectation.input = &NodeKeeperMockGetWorkingNodesByRoleInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockGetWorkingNodesByRoleExpectation) Return(r []insolar.Reference) {
	e.result = &NodeKeeperMockGetWorkingNodesByRoleResult{r}
}

//Set uses given function f as a mock of NodeKeeper.GetWorkingNodesByRole method
func (m *mNodeKeeperMockGetWorkingNodesByRole) Set(f func(p insolar.DynamicRole) (r []insolar.Reference)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodesByRoleFunc = f
	return m.mock
}

//GetWorkingNodesByRole implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) GetWorkingNodesByRole(p insolar.DynamicRole) (r []insolar.Reference) {
	counter := atomic.AddUint64(&m.GetWorkingNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodesByRoleCounter, 1)

	if len(m.GetWorkingNodesByRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodesByRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.GetWorkingNodesByRole. %v", p)
			return
		}

		input := m.GetWorkingNodesByRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockGetWorkingNodesByRoleInput{p}, "NodeKeeper.GetWorkingNodesByRole got unexpected parameters")

		result := m.GetWorkingNodesByRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetWorkingNodesByRole")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesByRoleMock.mainExpectation != nil {

		input := m.GetWorkingNodesByRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockGetWorkingNodesByRoleInput{p}, "NodeKeeper.GetWorkingNodesByRole got unexpected parameters")
		}

		result := m.GetWorkingNodesByRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.GetWorkingNodesByRole")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesByRoleFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.GetWorkingNodesByRole. %v", p)
		return
	}

	return m.GetWorkingNodesByRoleFunc(p)
}

//GetWorkingNodesByRoleMinimockCounter returns a count of NodeKeeperMock.GetWorkingNodesByRoleFunc invocations
func (m *NodeKeeperMock) GetWorkingNodesByRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter)
}

//GetWorkingNodesByRoleMinimockPreCounter returns the value of NodeKeeperMock.GetWorkingNodesByRole invocations
func (m *NodeKeeperMock) GetWorkingNodesByRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesByRolePreCounter)
}

//GetWorkingNodesByRoleFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) GetWorkingNodesByRoleFinished() bool {
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
	input  *NodeKeeperMockMoveSyncToActiveInput
	result *NodeKeeperMockMoveSyncToActiveResult
}

type NodeKeeperMockMoveSyncToActiveInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type NodeKeeperMockMoveSyncToActiveResult struct {
	r error
}

//Expect specifies that invocation of NodeKeeper.MoveSyncToActive is expected from 1 to Infinity times
func (m *mNodeKeeperMockMoveSyncToActive) Expect(p context.Context, p1 insolar.PulseNumber) *mNodeKeeperMockMoveSyncToActive {
	m.mock.MoveSyncToActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockMoveSyncToActiveExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockMoveSyncToActiveInput{p, p1}
	return m
}

//Return specifies results of invocation of NodeKeeper.MoveSyncToActive
func (m *mNodeKeeperMockMoveSyncToActive) Return(r error) *NodeKeeperMock {
	m.mock.MoveSyncToActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockMoveSyncToActiveExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockMoveSyncToActiveResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.MoveSyncToActive is expected once
func (m *mNodeKeeperMockMoveSyncToActive) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *NodeKeeperMockMoveSyncToActiveExpectation {
	m.mock.MoveSyncToActiveFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockMoveSyncToActiveExpectation{}
	expectation.input = &NodeKeeperMockMoveSyncToActiveInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockMoveSyncToActiveExpectation) Return(r error) {
	e.result = &NodeKeeperMockMoveSyncToActiveResult{r}
}

//Set uses given function f as a mock of NodeKeeper.MoveSyncToActive method
func (m *mNodeKeeperMockMoveSyncToActive) Set(f func(p context.Context, p1 insolar.PulseNumber) (r error)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MoveSyncToActiveFunc = f
	return m.mock
}

//MoveSyncToActive implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) MoveSyncToActive(p context.Context, p1 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.MoveSyncToActivePreCounter, 1)
	defer atomic.AddUint64(&m.MoveSyncToActiveCounter, 1)

	if len(m.MoveSyncToActiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MoveSyncToActiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.MoveSyncToActive. %v %v", p, p1)
			return
		}

		input := m.MoveSyncToActiveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockMoveSyncToActiveInput{p, p1}, "NodeKeeper.MoveSyncToActive got unexpected parameters")

		result := m.MoveSyncToActiveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.MoveSyncToActive")
			return
		}

		r = result.r

		return
	}

	if m.MoveSyncToActiveMock.mainExpectation != nil {

		input := m.MoveSyncToActiveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockMoveSyncToActiveInput{p, p1}, "NodeKeeper.MoveSyncToActive got unexpected parameters")
		}

		result := m.MoveSyncToActiveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.MoveSyncToActive")
		}

		r = result.r

		return
	}

	if m.MoveSyncToActiveFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.MoveSyncToActive. %v %v", p, p1)
		return
	}

	return m.MoveSyncToActiveFunc(p, p1)
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

type mNodeKeeperMockSetInitialSnapshot struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockSetInitialSnapshotExpectation
	expectationSeries []*NodeKeeperMockSetInitialSnapshotExpectation
}

type NodeKeeperMockSetInitialSnapshotExpectation struct {
	input *NodeKeeperMockSetInitialSnapshotInput
}

type NodeKeeperMockSetInitialSnapshotInput struct {
	p []insolar.NetworkNode
}

//Expect specifies that invocation of NodeKeeper.SetInitialSnapshot is expected from 1 to Infinity times
func (m *mNodeKeeperMockSetInitialSnapshot) Expect(p []insolar.NetworkNode) *mNodeKeeperMockSetInitialSnapshot {
	m.mock.SetInitialSnapshotFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetInitialSnapshotExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockSetInitialSnapshotInput{p}
	return m
}

//Return specifies results of invocation of NodeKeeper.SetInitialSnapshot
func (m *mNodeKeeperMockSetInitialSnapshot) Return() *NodeKeeperMock {
	m.mock.SetInitialSnapshotFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSetInitialSnapshotExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.SetInitialSnapshot is expected once
func (m *mNodeKeeperMockSetInitialSnapshot) ExpectOnce(p []insolar.NetworkNode) *NodeKeeperMockSetInitialSnapshotExpectation {
	m.mock.SetInitialSnapshotFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockSetInitialSnapshotExpectation{}
	expectation.input = &NodeKeeperMockSetInitialSnapshotInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of NodeKeeper.SetInitialSnapshot method
func (m *mNodeKeeperMockSetInitialSnapshot) Set(f func(p []insolar.NetworkNode)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetInitialSnapshotFunc = f
	return m.mock
}

//SetInitialSnapshot implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) SetInitialSnapshot(p []insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.SetInitialSnapshotPreCounter, 1)
	defer atomic.AddUint64(&m.SetInitialSnapshotCounter, 1)

	if len(m.SetInitialSnapshotMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetInitialSnapshotMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.SetInitialSnapshot. %v", p)
			return
		}

		input := m.SetInitialSnapshotMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockSetInitialSnapshotInput{p}, "NodeKeeper.SetInitialSnapshot got unexpected parameters")

		return
	}

	if m.SetInitialSnapshotMock.mainExpectation != nil {

		input := m.SetInitialSnapshotMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockSetInitialSnapshotInput{p}, "NodeKeeper.SetInitialSnapshot got unexpected parameters")
		}

		return
	}

	if m.SetInitialSnapshotFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.SetInitialSnapshot. %v", p)
		return
	}

	m.SetInitialSnapshotFunc(p)
}

//SetInitialSnapshotMinimockCounter returns a count of NodeKeeperMock.SetInitialSnapshotFunc invocations
func (m *NodeKeeperMock) SetInitialSnapshotMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetInitialSnapshotCounter)
}

//SetInitialSnapshotMinimockPreCounter returns the value of NodeKeeperMock.SetInitialSnapshot invocations
func (m *NodeKeeperMock) SetInitialSnapshotMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetInitialSnapshotPreCounter)
}

//SetInitialSnapshotFinished returns true if mock invocations count is ok
func (m *NodeKeeperMock) SetInitialSnapshotFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetInitialSnapshotMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetInitialSnapshotCounter) == uint64(len(m.SetInitialSnapshotMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetInitialSnapshotMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetInitialSnapshotCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetInitialSnapshotFunc != nil {
		return atomic.LoadUint64(&m.SetInitialSnapshotCounter) > 0
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

type mNodeKeeperMockSync struct {
	mock              *NodeKeeperMock
	mainExpectation   *NodeKeeperMockSyncExpectation
	expectationSeries []*NodeKeeperMockSyncExpectation
}

type NodeKeeperMockSyncExpectation struct {
	input  *NodeKeeperMockSyncInput
	result *NodeKeeperMockSyncResult
}

type NodeKeeperMockSyncInput struct {
	p  context.Context
	p1 []insolar.NetworkNode
	p2 []packets.ReferendumClaim
}

type NodeKeeperMockSyncResult struct {
	r error
}

//Expect specifies that invocation of NodeKeeper.Sync is expected from 1 to Infinity times
func (m *mNodeKeeperMockSync) Expect(p context.Context, p1 []insolar.NetworkNode, p2 []packets.ReferendumClaim) *mNodeKeeperMockSync {
	m.mock.SyncFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSyncExpectation{}
	}
	m.mainExpectation.input = &NodeKeeperMockSyncInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of NodeKeeper.Sync
func (m *mNodeKeeperMockSync) Return(r error) *NodeKeeperMock {
	m.mock.SyncFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeKeeperMockSyncExpectation{}
	}
	m.mainExpectation.result = &NodeKeeperMockSyncResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeKeeper.Sync is expected once
func (m *mNodeKeeperMockSync) ExpectOnce(p context.Context, p1 []insolar.NetworkNode, p2 []packets.ReferendumClaim) *NodeKeeperMockSyncExpectation {
	m.mock.SyncFunc = nil
	m.mainExpectation = nil

	expectation := &NodeKeeperMockSyncExpectation{}
	expectation.input = &NodeKeeperMockSyncInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeKeeperMockSyncExpectation) Return(r error) {
	e.result = &NodeKeeperMockSyncResult{r}
}

//Set uses given function f as a mock of NodeKeeper.Sync method
func (m *mNodeKeeperMockSync) Set(f func(p context.Context, p1 []insolar.NetworkNode, p2 []packets.ReferendumClaim) (r error)) *NodeKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SyncFunc = f
	return m.mock
}

//Sync implements github.com/insolar/insolar/network.NodeKeeper interface
func (m *NodeKeeperMock) Sync(p context.Context, p1 []insolar.NetworkNode, p2 []packets.ReferendumClaim) (r error) {
	counter := atomic.AddUint64(&m.SyncPreCounter, 1)
	defer atomic.AddUint64(&m.SyncCounter, 1)

	if len(m.SyncMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SyncMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeKeeperMock.Sync. %v %v %v", p, p1, p2)
			return
		}

		input := m.SyncMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeKeeperMockSyncInput{p, p1, p2}, "NodeKeeper.Sync got unexpected parameters")

		result := m.SyncMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.Sync")
			return
		}

		r = result.r

		return
	}

	if m.SyncMock.mainExpectation != nil {

		input := m.SyncMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeKeeperMockSyncInput{p, p1, p2}, "NodeKeeper.Sync got unexpected parameters")
		}

		result := m.SyncMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeKeeperMock.Sync")
		}

		r = result.r

		return
	}

	if m.SyncFunc == nil {
		m.t.Fatalf("Unexpected call to NodeKeeperMock.Sync. %v %v %v", p, p1, p2)
		return
	}

	return m.SyncFunc(p, p1, p2)
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

	if !m.GetAccessorFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetAccessor")
	}

	if !m.GetClaimQueueFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetClaimQueue")
	}

	if !m.GetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetCloudHash")
	}

	if !m.GetConsensusInfoFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetConsensusInfo")
	}

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOrigin")
	}

	if !m.GetOriginAnnounceClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginAnnounceClaim")
	}

	if !m.GetOriginJoinClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginJoinClaim")
	}

	if !m.GetSnapshotCopyFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetSnapshotCopy")
	}

	if !m.GetWorkingNodeFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetWorkingNode")
	}

	if !m.GetWorkingNodesFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetWorkingNodes")
	}

	if !m.GetWorkingNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetWorkingNodesByRole")
	}

	if !m.IsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.IsBootstrapped")
	}

	if !m.MoveSyncToActiveFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.MoveSyncToActive")
	}

	if !m.SetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetCloudHash")
	}

	if !m.SetInitialSnapshotFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetInitialSnapshot")
	}

	if !m.SetIsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetIsBootstrapped")
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

	if !m.GetAccessorFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetAccessor")
	}

	if !m.GetClaimQueueFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetClaimQueue")
	}

	if !m.GetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetCloudHash")
	}

	if !m.GetConsensusInfoFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetConsensusInfo")
	}

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOrigin")
	}

	if !m.GetOriginAnnounceClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginAnnounceClaim")
	}

	if !m.GetOriginJoinClaimFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetOriginJoinClaim")
	}

	if !m.GetSnapshotCopyFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetSnapshotCopy")
	}

	if !m.GetWorkingNodeFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetWorkingNode")
	}

	if !m.GetWorkingNodesFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetWorkingNodes")
	}

	if !m.GetWorkingNodesByRoleFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.GetWorkingNodesByRole")
	}

	if !m.IsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.IsBootstrapped")
	}

	if !m.MoveSyncToActiveFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.MoveSyncToActive")
	}

	if !m.SetCloudHashFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetCloudHash")
	}

	if !m.SetInitialSnapshotFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetInitialSnapshot")
	}

	if !m.SetIsBootstrappedFinished() {
		m.t.Fatal("Expected call to NodeKeeperMock.SetIsBootstrapped")
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
		ok = ok && m.GetAccessorFinished()
		ok = ok && m.GetClaimQueueFinished()
		ok = ok && m.GetCloudHashFinished()
		ok = ok && m.GetConsensusInfoFinished()
		ok = ok && m.GetOriginFinished()
		ok = ok && m.GetOriginAnnounceClaimFinished()
		ok = ok && m.GetOriginJoinClaimFinished()
		ok = ok && m.GetSnapshotCopyFinished()
		ok = ok && m.GetWorkingNodeFinished()
		ok = ok && m.GetWorkingNodesFinished()
		ok = ok && m.GetWorkingNodesByRoleFinished()
		ok = ok && m.IsBootstrappedFinished()
		ok = ok && m.MoveSyncToActiveFinished()
		ok = ok && m.SetCloudHashFinished()
		ok = ok && m.SetInitialSnapshotFinished()
		ok = ok && m.SetIsBootstrappedFinished()
		ok = ok && m.SyncFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetAccessorFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetAccessor")
			}

			if !m.GetClaimQueueFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetClaimQueue")
			}

			if !m.GetCloudHashFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetCloudHash")
			}

			if !m.GetConsensusInfoFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetConsensusInfo")
			}

			if !m.GetOriginFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetOrigin")
			}

			if !m.GetOriginAnnounceClaimFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetOriginAnnounceClaim")
			}

			if !m.GetOriginJoinClaimFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetOriginJoinClaim")
			}

			if !m.GetSnapshotCopyFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetSnapshotCopy")
			}

			if !m.GetWorkingNodeFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetWorkingNode")
			}

			if !m.GetWorkingNodesFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetWorkingNodes")
			}

			if !m.GetWorkingNodesByRoleFinished() {
				m.t.Error("Expected call to NodeKeeperMock.GetWorkingNodesByRole")
			}

			if !m.IsBootstrappedFinished() {
				m.t.Error("Expected call to NodeKeeperMock.IsBootstrapped")
			}

			if !m.MoveSyncToActiveFinished() {
				m.t.Error("Expected call to NodeKeeperMock.MoveSyncToActive")
			}

			if !m.SetCloudHashFinished() {
				m.t.Error("Expected call to NodeKeeperMock.SetCloudHash")
			}

			if !m.SetInitialSnapshotFinished() {
				m.t.Error("Expected call to NodeKeeperMock.SetInitialSnapshot")
			}

			if !m.SetIsBootstrappedFinished() {
				m.t.Error("Expected call to NodeKeeperMock.SetIsBootstrapped")
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

	if !m.GetAccessorFinished() {
		return false
	}

	if !m.GetClaimQueueFinished() {
		return false
	}

	if !m.GetCloudHashFinished() {
		return false
	}

	if !m.GetConsensusInfoFinished() {
		return false
	}

	if !m.GetOriginFinished() {
		return false
	}

	if !m.GetOriginAnnounceClaimFinished() {
		return false
	}

	if !m.GetOriginJoinClaimFinished() {
		return false
	}

	if !m.GetSnapshotCopyFinished() {
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

	if !m.IsBootstrappedFinished() {
		return false
	}

	if !m.MoveSyncToActiveFinished() {
		return false
	}

	if !m.SetCloudHashFinished() {
		return false
	}

	if !m.SetInitialSnapshotFinished() {
		return false
	}

	if !m.SetIsBootstrappedFinished() {
		return false
	}

	if !m.SyncFinished() {
		return false
	}

	return true
}

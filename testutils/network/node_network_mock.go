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

	GetActiveNodeFunc       func(p core.RecordRef) (r core.Node)
	GetActiveNodeCounter    uint64
	GetActiveNodePreCounter uint64
	GetActiveNodeMock       mNodeNetworkMockGetActiveNode

	GetActiveNodesFunc       func() (r []core.Node)
	GetActiveNodesCounter    uint64
	GetActiveNodesPreCounter uint64
	GetActiveNodesMock       mNodeNetworkMockGetActiveNodes

	GetActiveNodesByRoleFunc       func(p core.JetRole) (r []core.RecordRef)
	GetActiveNodesByRoleCounter    uint64
	GetActiveNodesByRolePreCounter uint64
	GetActiveNodesByRoleMock       mNodeNetworkMockGetActiveNodesByRole

	GetOriginFunc       func() (r core.Node)
	GetOriginCounter    uint64
	GetOriginPreCounter uint64
	GetOriginMock       mNodeNetworkMockGetOrigin
}

//NewNodeNetworkMock returns a mock for github.com/insolar/insolar/core.NodeNetwork
func NewNodeNetworkMock(t minimock.Tester) *NodeNetworkMock {
	m := &NodeNetworkMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetActiveNodeMock = mNodeNetworkMockGetActiveNode{mock: m}
	m.GetActiveNodesMock = mNodeNetworkMockGetActiveNodes{mock: m}
	m.GetActiveNodesByRoleMock = mNodeNetworkMockGetActiveNodesByRole{mock: m}
	m.GetOriginMock = mNodeNetworkMockGetOrigin{mock: m}

	return m
}

type mNodeNetworkMockGetActiveNode struct {
	mock             *NodeNetworkMock
	mockExpectations *NodeNetworkMockGetActiveNodeParams
}

//NodeNetworkMockGetActiveNodeParams represents input parameters of the NodeNetwork.GetActiveNode
type NodeNetworkMockGetActiveNodeParams struct {
	p core.RecordRef
}

//Expect sets up expected params for the NodeNetwork.GetActiveNode
func (m *mNodeNetworkMockGetActiveNode) Expect(p core.RecordRef) *mNodeNetworkMockGetActiveNode {
	m.mockExpectations = &NodeNetworkMockGetActiveNodeParams{p}
	return m
}

//Return sets up a mock for NodeNetwork.GetActiveNode to return Return's arguments
func (m *mNodeNetworkMockGetActiveNode) Return(r core.Node) *NodeNetworkMock {
	m.mock.GetActiveNodeFunc = func(p core.RecordRef) core.Node {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeNetwork.GetActiveNode method
func (m *mNodeNetworkMockGetActiveNode) Set(f func(p core.RecordRef) (r core.Node)) *NodeNetworkMock {
	m.mock.GetActiveNodeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetActiveNode implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetActiveNode(p core.RecordRef) (r core.Node) {
	atomic.AddUint64(&m.GetActiveNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeCounter, 1)

	if m.GetActiveNodeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetActiveNodeMock.mockExpectations, NodeNetworkMockGetActiveNodeParams{p},
			"NodeNetwork.GetActiveNode got unexpected parameters")

		if m.GetActiveNodeFunc == nil {

			m.t.Fatal("No results are set for the NodeNetworkMock.GetActiveNode")

			return
		}
	}

	if m.GetActiveNodeFunc == nil {
		m.t.Fatal("Unexpected call to NodeNetworkMock.GetActiveNode")
		return
	}

	return m.GetActiveNodeFunc(p)
}

//GetActiveNodeMinimockCounter returns a count of NodeNetworkMock.GetActiveNodeFunc invocations
func (m *NodeNetworkMock) GetActiveNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeCounter)
}

//GetActiveNodeMinimockPreCounter returns the value of NodeNetworkMock.GetActiveNode invocations
func (m *NodeNetworkMock) GetActiveNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodePreCounter)
}

type mNodeNetworkMockGetActiveNodes struct {
	mock *NodeNetworkMock
}

//Return sets up a mock for NodeNetwork.GetActiveNodes to return Return's arguments
func (m *mNodeNetworkMockGetActiveNodes) Return(r []core.Node) *NodeNetworkMock {
	m.mock.GetActiveNodesFunc = func() []core.Node {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeNetwork.GetActiveNodes method
func (m *mNodeNetworkMockGetActiveNodes) Set(f func() (r []core.Node)) *NodeNetworkMock {
	m.mock.GetActiveNodesFunc = f

	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetActiveNodes() (r []core.Node) {
	atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if m.GetActiveNodesFunc == nil {
		m.t.Fatal("Unexpected call to NodeNetworkMock.GetActiveNodes")
		return
	}

	return m.GetActiveNodesFunc()
}

//GetActiveNodesMinimockCounter returns a count of NodeNetworkMock.GetActiveNodesFunc invocations
func (m *NodeNetworkMock) GetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesCounter)
}

//GetActiveNodesMinimockPreCounter returns the value of NodeNetworkMock.GetActiveNodes invocations
func (m *NodeNetworkMock) GetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesPreCounter)
}

type mNodeNetworkMockGetActiveNodesByRole struct {
	mock             *NodeNetworkMock
	mockExpectations *NodeNetworkMockGetActiveNodesByRoleParams
}

//NodeNetworkMockGetActiveNodesByRoleParams represents input parameters of the NodeNetwork.GetActiveNodesByRole
type NodeNetworkMockGetActiveNodesByRoleParams struct {
	p core.JetRole
}

//Expect sets up expected params for the NodeNetwork.GetActiveNodesByRole
func (m *mNodeNetworkMockGetActiveNodesByRole) Expect(p core.JetRole) *mNodeNetworkMockGetActiveNodesByRole {
	m.mockExpectations = &NodeNetworkMockGetActiveNodesByRoleParams{p}
	return m
}

//Return sets up a mock for NodeNetwork.GetActiveNodesByRole to return Return's arguments
func (m *mNodeNetworkMockGetActiveNodesByRole) Return(r []core.RecordRef) *NodeNetworkMock {
	m.mock.GetActiveNodesByRoleFunc = func(p core.JetRole) []core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeNetwork.GetActiveNodesByRole method
func (m *mNodeNetworkMockGetActiveNodesByRole) Set(f func(p core.JetRole) (r []core.RecordRef)) *NodeNetworkMock {
	m.mock.GetActiveNodesByRoleFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetActiveNodesByRole implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetActiveNodesByRole(p core.JetRole) (r []core.RecordRef) {
	atomic.AddUint64(&m.GetActiveNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesByRoleCounter, 1)

	if m.GetActiveNodesByRoleMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetActiveNodesByRoleMock.mockExpectations, NodeNetworkMockGetActiveNodesByRoleParams{p},
			"NodeNetwork.GetActiveNodesByRole got unexpected parameters")

		if m.GetActiveNodesByRoleFunc == nil {

			m.t.Fatal("No results are set for the NodeNetworkMock.GetActiveNodesByRole")

			return
		}
	}

	if m.GetActiveNodesByRoleFunc == nil {
		m.t.Fatal("Unexpected call to NodeNetworkMock.GetActiveNodesByRole")
		return
	}

	return m.GetActiveNodesByRoleFunc(p)
}

//GetActiveNodesByRoleMinimockCounter returns a count of NodeNetworkMock.GetActiveNodesByRoleFunc invocations
func (m *NodeNetworkMock) GetActiveNodesByRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter)
}

//GetActiveNodesByRoleMinimockPreCounter returns the value of NodeNetworkMock.GetActiveNodesByRole invocations
func (m *NodeNetworkMock) GetActiveNodesByRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRolePreCounter)
}

type mNodeNetworkMockGetOrigin struct {
	mock *NodeNetworkMock
}

//Return sets up a mock for NodeNetwork.GetOrigin to return Return's arguments
func (m *mNodeNetworkMockGetOrigin) Return(r core.Node) *NodeNetworkMock {
	m.mock.GetOriginFunc = func() core.Node {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NodeNetwork.GetOrigin method
func (m *mNodeNetworkMockGetOrigin) Set(f func() (r core.Node)) *NodeNetworkMock {
	m.mock.GetOriginFunc = f

	return m.mock
}

//GetOrigin implements github.com/insolar/insolar/core.NodeNetwork interface
func (m *NodeNetworkMock) GetOrigin() (r core.Node) {
	atomic.AddUint64(&m.GetOriginPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginCounter, 1)

	if m.GetOriginFunc == nil {
		m.t.Fatal("Unexpected call to NodeNetworkMock.GetOrigin")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeNetworkMock) ValidateCallCounters() {

	if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetActiveNode")
	}

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetActiveNodes")
	}

	if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetActiveNodesByRole")
	}

	if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetOrigin")
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

	if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetActiveNode")
	}

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetActiveNodes")
	}

	if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetActiveNodesByRole")
	}

	if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
		m.t.Fatal("Expected call to NodeNetworkMock.GetOrigin")
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
		ok = ok && (m.GetActiveNodeFunc == nil || atomic.LoadUint64(&m.GetActiveNodeCounter) > 0)
		ok = ok && (m.GetActiveNodesFunc == nil || atomic.LoadUint64(&m.GetActiveNodesCounter) > 0)
		ok = ok && (m.GetActiveNodesByRoleFunc == nil || atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) > 0)
		ok = ok && (m.GetOriginFunc == nil || atomic.LoadUint64(&m.GetOriginCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
				m.t.Error("Expected call to NodeNetworkMock.GetActiveNode")
			}

			if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
				m.t.Error("Expected call to NodeNetworkMock.GetActiveNodes")
			}

			if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
				m.t.Error("Expected call to NodeNetworkMock.GetActiveNodesByRole")
			}

			if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
				m.t.Error("Expected call to NodeNetworkMock.GetOrigin")
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

	if m.GetActiveNodeFunc != nil && atomic.LoadUint64(&m.GetActiveNodeCounter) == 0 {
		return false
	}

	if m.GetActiveNodesFunc != nil && atomic.LoadUint64(&m.GetActiveNodesCounter) == 0 {
		return false
	}

	if m.GetActiveNodesByRoleFunc != nil && atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == 0 {
		return false
	}

	if m.GetOriginFunc != nil && atomic.LoadUint64(&m.GetOriginCounter) == 0 {
		return false
	}

	return true
}

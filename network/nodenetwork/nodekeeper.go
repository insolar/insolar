/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package nodenetwork

import (
	"bytes"
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/insolar/insolar/configuration"
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"

	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
)

// NewNodeNetwork create active node component
func NewNodeNetwork(configuration configuration.Configuration) (core.NodeNetwork, error) {
	origin, err := createOrigin(configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create origin node")
	}
	nodeKeeper := NewNodeKeeper(origin)

	if len(configuration.Host.BootstrapHosts) == 0 {
		log.Info("Bootstrap nodes are not set. Init zeronet.")
		nodeKeeper.AddActiveNodes([]core.Node{origin})
	}

	return nodeKeeper, nil
}

func createOrigin(configuration configuration.Configuration) (MutableNode, error) {
	nodeID := core.NewRefFromBase58(configuration.Node.Node.ID)
	publicAddress, err := resolveAddress(configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to resolve public address")
	}

	// TODO: get roles from certificate
	// TODO: pass public key
	return newMutableNode(
		nodeID,
		[]core.NodeRole{core.RoleVirtual, core.RoleHeavyMaterial, core.RoleLightMaterial},
		nil,
		0,
		publicAddress,
		version.Version,
	), nil
}

func resolveAddress(configuration configuration.Configuration) (string, error) {
	conn, address, err := transport.NewConnection(configuration.Host.Transport)
	if err != nil {
		return "", err
	}
	err = conn.Close()
	if err != nil {
		log.Warn(err)
	}
	return address, nil
}

// NewNodeKeeper create new NodeKeeper
func NewNodeKeeper(origin core.Node) network.NodeKeeper {
	return &nodekeeper{
		origin:       origin,
		state:        network.Undefined,
		claimQueue:   newClaimQueue(),
		active:       make(map[core.RecordRef]core.Node),
		indexNode:    make(map[core.NodeRole][]core.RecordRef),
		indexShortID: make(map[core.ShortNodeID]core.Node),
	}
}

type nodekeeper struct {
	origin      core.Node
	originClaim *consensus.NodeJoinClaim
	originLock  sync.RWMutex
	state       network.NodeKeeperState
	claimQueue  *claimQueue

	nodesJoinedDuringPrevPulse bool

	cloudHashLock sync.RWMutex
	cloudHash     []byte

	activeLock   sync.RWMutex
	active       map[core.RecordRef]core.Node
	indexNode    map[core.NodeRole][]core.RecordRef
	indexShortID map[core.ShortNodeID]core.Node

	lock   sync.Mutex
	unsync []network.NodeClaim
	sync   []network.NodeClaim
}

func (nk *nodekeeper) Start(ctx context.Context, components core.Components) error {
	return nil
}

func (nk *nodekeeper) Stop(ctx context.Context) error {
	return nil
}

func (nk *nodekeeper) GetOrigin() core.Node {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.origin
}

func (nk *nodekeeper) GetCloudHash() []byte {
	nk.cloudHashLock.RLock()
	defer nk.cloudHashLock.RUnlock()

	return nk.cloudHash
}

func (nk *nodekeeper) SetCloudHash(cloudHash []byte) {
	nk.cloudHashLock.Lock()
	defer nk.cloudHashLock.Unlock()

	nk.cloudHash = cloudHash
}

func (nk *nodekeeper) GetActiveNodes() []core.Node {
	nk.activeLock.RLock()
	result := make([]core.Node, len(nk.active))
	index := 0
	for _, node := range nk.active {
		result[index] = node
		index++
	}
	nk.activeLock.RUnlock()
	// Sort active nodes to return list with determinate order on every node.
	// If we have more than 10k nodes, we need to optimize this
	sort.Slice(result, func(i, j int) bool {
		return bytes.Compare(result[i].ID().Bytes(), result[j].ID().Bytes()) < 0
	})
	return result
}

func (nk *nodekeeper) GetActiveNodesByRole(role core.JetRole) []core.RecordRef {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	list, exists := nk.indexNode[jetRoleToNodeRole(role)]
	if !exists {
		return nil
	}
	result := make([]core.RecordRef, len(list))
	copy(result, list)
	return result
}

func (nk *nodekeeper) AddActiveNodes(nodes []core.Node) {
	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	activeNodes := make([]string, len(nodes))
	for i, node := range nodes {
		nk.addActiveNode(node)
		activeNodes[i] = node.ID().String()
	}
	log.Debugf("Added active nodes: %s", strings.Join(activeNodes, ", "))
}

func (nk *nodekeeper) GetActiveNode(ref core.RecordRef) core.Node {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.active[ref]
}

func (nk *nodekeeper) GetActiveNodeByShortID(shortID core.ShortNodeID) core.Node {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.indexShortID[shortID]
}

func (nk *nodekeeper) addActiveNode(node core.Node) {
	if node.ID().Equal(nk.origin.ID()) {
		nk.origin = node
		log.Infof("Added origin node %s to active list", nk.origin.ID())
	}
	nk.active[node.ID()] = node
	for _, role := range node.Roles() {
		list, ok := nk.indexNode[role]
		if !ok {
			list := make([]core.RecordRef, 0)
			nk.indexNode[role] = list
		}
		nk.indexNode[role] = append(list, node.ID())
	}
	nk.indexShortID[node.ShortID()] = node
}

func (nk *nodekeeper) delActiveNode(ref core.RecordRef) {
	if ref.Equal(nk.origin.ID()) {
		// we received acknowledge to leave, can gracefully stop
		// TODO: graceful stop instead of panic
		panic("Node leave acknowledged by network. Goodbye!")
	}
	active, ok := nk.active[ref]
	if !ok {
		return
	}
	delete(nk.active, ref)
	delete(nk.indexShortID, active.ShortID())
	// we should have one role in Node structure in future
	delete(nk.indexNode, active.Roles()[0])
}

func (nk *nodekeeper) SetState(state network.NodeKeeperState) {
	nk.state = state
}

func (nk *nodekeeper) GetState() network.NodeKeeperState {
	return nk.state
}

func (nk *nodekeeper) SetOriginClaim(claim *consensus.NodeJoinClaim) {
	nk.originLock.Lock()
	defer nk.originLock.Unlock()

	nk.originClaim = claim
}

func (nk *nodekeeper) GetOriginClaim() *consensus.NodeJoinClaim {
	nk.originLock.RLock()
	defer nk.originLock.RUnlock()

	return nk.originClaim
}

func (nk *nodekeeper) AddPendingClaim(claim consensus.ReferendumClaim) bool {
	nk.claimQueue.Push(claim)
	return true
}

func (nk *nodekeeper) GetClaimQueue() network.ClaimQueue {
	return nk.claimQueue
}

func (nk *nodekeeper) NodesJoinedDuringPreviousPulse() bool {
	return nk.nodesJoinedDuringPrevPulse
}

func (nk *nodekeeper) AddUnsyncClaims([]*network.NodeClaim) {
	log.Error("implement me!")
}

func (nk *nodekeeper) CalculateUnsyncMergedHash() []byte {
	log.Error("implement me!")
	return nil
}

func (nk *nodekeeper) Sync(deviant []core.Node) {
	nk.lock.Lock()
	defer nk.lock.Unlock()

	nk.sync = make([]network.NodeClaim, 0)
	for _, claim := range nk.unsync {
		isDeviant := false
		for _, d := range deviant {
			if d.ID().Equal(claim.Initiator) {
				isDeviant = true
				break
			}

			if !isDeviant {
				nk.sync = append(nk.sync, claim)
			}
		}
	}
	nk.unsync = make([]network.NodeClaim, 0)
}

func (nk *nodekeeper) MoveSyncToActive() {
	nk.lock.Lock()
	nk.activeLock.Lock()
	defer func() {
		nk.activeLock.Unlock()
		nk.lock.Unlock()
	}()

	nk.nodesJoinedDuringPrevPulse = false
	for _, nodeClaim := range nk.sync {
		switch nodeClaim.Claim.Type() {
		case consensus.TypeNodeJoinClaim:
			claim := nodeClaim.Claim.(*consensus.NodeJoinClaim)
			nk.addActiveNode(claim.Node())
			nk.nodesJoinedDuringPrevPulse = true
		case consensus.TypeNodeLeaveClaim:
			nk.delActiveNode(nodeClaim.Initiator)
		default:
			log.Error("unknown claim type: %d", nodeClaim.Claim.Type())
		}
	}

	nk.sync = make([]network.NodeClaim, 0)
}

func jetRoleToNodeRole(role core.JetRole) core.NodeRole {
	switch role {
	case core.RoleVirtualExecutor:
		return core.RoleVirtual
	case core.RoleVirtualValidator:
		return core.RoleVirtual
	case core.RoleLightExecutor:
		return core.RoleLightMaterial
	case core.RoleLightValidator:
		return core.RoleLightMaterial
	case core.RoleHeavyExecutor:
		return core.RoleHeavyMaterial
	default:
		return core.RoleUnknown
	}
}

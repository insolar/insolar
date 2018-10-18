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

package nodekeeper

import (
	"bytes"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/pkg/errors"
)

// NewActiveNodeComponent create active node component
func NewActiveNodeComponent(configuration configuration.Configuration) (core.ActiveNodeComponent, error) {
	nodeID := core.NewRefFromBase58(configuration.Node.Node.ID)
	nodeKeeper := NewNodeKeeper(nodeID)
	// TODO: get roles from certificate
	// TODO: pass public key
	if len(configuration.Host.BootstrapHosts) == 0 {
		publicAddress, err := resolveAddress(configuration)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create active node component")
		}

		log.Info("Bootstrap nodes is not set. Init zeronet.")
		nodeKeeper.AddActiveNodes([]*core.ActiveNode{&core.ActiveNode{
			NodeID:   nodeID,
			PulseNum: 0,
			State:    core.NodeActive,
			Roles:    []core.NodeRole{core.RoleVirtual, core.RoleHeavyMaterial, core.RoleLightMaterial},
			Address:  publicAddress,
			// PublicKey: ???
		}})
	}
	return nodeKeeper, nil
}

func resolveAddress(configuration configuration.Configuration) (string, error) {
	conn, address, err := transport.NewConnection(configuration.Host.Transport)
	defer func() { _ = conn.Close() }()
	if err != nil {
		return "", err
	}
	return address, nil
}

// NewNodeKeeper create new NodeKeeper
func NewNodeKeeper(nodeID core.RecordRef) consensus.NodeKeeper {
	return &nodekeeper{
		nodeID:      nodeID,
		state:       undefined,
		active:      make(map[core.RecordRef]*core.ActiveNode),
		sync:        make([]*core.ActiveNode, 0),
		unsync:      make([]*core.ActiveNode, 0),
		listWaiters: make([]chan *UnsyncList, 0),
		nodeWaiters: make(map[core.RecordRef]chan *core.ActiveNode),
		index:       make(map[core.NodeRole][]core.RecordRef),
	}
}

type nodekeeperState uint8

const (
	undefined = nodekeeperState(iota + 1)
	pulseSet
	synced
)

type nodekeeper struct {
	nodeID core.RecordRef
	self   *core.ActiveNode
	state  nodekeeperState
	pulse  core.PulseNumber

	activeLock sync.RWMutex
	active     map[core.RecordRef]*core.ActiveNode
	index      map[core.NodeRole][]core.RecordRef
	sync       []*core.ActiveNode

	unsyncLock  sync.Mutex
	unsync      []*core.ActiveNode
	unsyncList  *UnsyncList
	listWaiters []chan *UnsyncList
	nodeWaiters map[core.RecordRef]chan *core.ActiveNode
}

func (nk *nodekeeper) Start(components core.Components) error {
	return nil
}

func (nk *nodekeeper) Stop() error {
	return nil
}

func (nk *nodekeeper) GetID() core.RecordRef {
	return nk.nodeID
}

func (nk *nodekeeper) GetSelf() *core.ActiveNode {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.self
}

func (nk *nodekeeper) GetActiveNodes() []*core.ActiveNode {
	nk.activeLock.RLock()
	result := make([]*core.ActiveNode, len(nk.active))
	index := 0
	for _, node := range nk.active {
		result[index] = node
		index++
	}
	nk.activeLock.RUnlock()
	// Sort active nodes to return list with determinate order on every node.
	// If we have more than 10k nodes, we need to optimize this
	sort.Slice(result, func(i, j int) bool {
		return bytes.Compare(result[i].NodeID[:], result[j].NodeID[:]) < 0
	})
	return result
}

func (nk *nodekeeper) GetActiveNodesByRole(role core.JetRole) []core.RecordRef {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	list, exists := nk.index[calculateJetRole(role)]
	if !exists {
		return nil
	}
	result := make([]core.RecordRef, len(list))
	copy(result, list)
	return result
}

func (nk *nodekeeper) AddActiveNodes(nodes []*core.ActiveNode) {
	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	activeNodes := make([]string, len(nodes))
	for i, node := range nodes {
		if node.NodeID.Equal(nk.nodeID) {
			nk.self = node
			log.Infof("Added self node %s to active list", nk.nodeID)
		}
		nk.active[node.NodeID] = node
		activeNodes[i] = node.NodeID.String()

		for _, role := range node.Roles {
			list, ok := nk.index[role]
			if !ok {
				list := make([]core.RecordRef, 0)
				nk.index[role] = list
			}
			nk.index[role] = append(list, node.NodeID)
		}
	}
	log.Debugf("Added active nodes: %s", strings.Join(activeNodes, ", "))
}

func (nk *nodekeeper) GetActiveNode(ref core.RecordRef) *core.ActiveNode {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.active[ref]
}

func (nk *nodekeeper) SetPulse(number core.PulseNumber) (bool, consensus.UnsyncList) {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.state == undefined {
		return true, nk.collectUnsync(number)
	}

	if number <= nk.pulse {
		log.Warnf("NodeKeeper: ignored SetPulse call with number=%d while current=%d", uint32(number), uint32(nk.pulse))
		return false, nil
	}

	if nk.state == pulseSet {
		log.Warn("NodeKeeper: SetPulse called from `pulseSet` state")
		nk.activeLock.Lock()
		nk.syncUnsafe(nil)
		nk.activeLock.Unlock()
	}

	return true, nk.collectUnsync(number)
}

func (nk *nodekeeper) Sync(syncCandidates []*core.ActiveNode, number core.PulseNumber) {
	nk.unsyncLock.Lock()
	nk.activeLock.Lock()

	defer func() {
		nk.activeLock.Unlock()
		nk.unsyncLock.Unlock()
	}()

	if nk.state == synced || nk.state == undefined {
		log.Warn("NodeKeeper: ignored Sync call from `synced` or `undefined` state")
		return
	}

	if nk.pulse > number {
		log.Warnf("NodeKeeper: ignored Sync call because passed number %d is less than internal number %d",
			number, nk.pulse)
		return
	}

	var candidates string
	for _, node := range syncCandidates {
		candidates += node.NodeID.String() + ", "
	}
	log.Debugf("Moving unsync to sync: %s", candidates)

	nk.syncUnsafe(syncCandidates)
}

func (nk *nodekeeper) AddUnsync(nodeID core.RecordRef, roles []core.NodeRole, address string /*, publicKey *ecdsa.PublicKey*/) (chan *core.ActiveNode, error) {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.self == nil {
		return nil, errors.New("cannot add node to unsync list: current node is not active")
	}

	node := &core.ActiveNode{
		NodeID:   nodeID,
		PulseNum: nk.pulse,
		State:    core.NodeJoined,
		Roles:    roles,
		Address:  address,
		// PublicKey: publicKey,
	}

	nk.unsync = append(nk.unsync, node)
	ch := make(chan *core.ActiveNode, 1)
	nk.nodeWaiters[node.NodeID] = ch
	return ch, nil
}

func (nk *nodekeeper) GetUnsyncHolder(pulse core.PulseNumber, duration time.Duration) (consensus.UnsyncList, error) {
	nk.unsyncLock.Lock()
	currentPulse := nk.pulse

	if currentPulse == pulse {
		result := nk.unsyncList
		nk.unsyncLock.Unlock()
		return result, nil
	}
	if currentPulse > pulse || duration < 0 {
		nk.unsyncLock.Unlock()
		return nil, errors.Errorf("GetUnsyncHolder called with pulse %d, but current NodeKeeper pulse is %d",
			pulse, currentPulse)
	}
	ch := make(chan *UnsyncList, 1)
	nk.listWaiters = append(nk.listWaiters, ch)
	nk.unsyncLock.Unlock()
	var result *UnsyncList
	select {
	case data := <-ch:
		if data == nil {
			return nil, errors.New("GetUnsyncHolder: channel closed")
		}
		result = data
	case <-time.After(duration):
		return nil, errors.New("GetUnsyncHolder: timeout")
	}
	if result.GetPulse() != pulse {
		return nil, errors.Errorf("GetUnsyncHolder called with pulse %d, but current UnsyncHolder pulse is %d",
			pulse, result.GetPulse())
	}
	return result, nil
}

func (nk *nodekeeper) syncUnsafe(syncCandidates []*core.ActiveNode) {
	// sync -> active
	for _, node := range nk.sync {
		nk.active[node.NodeID] = node
	}
	// unsync -> sync
	nk.sync = syncCandidates

	// first notify all synced nodes that they have passed the consensus
	for _, node := range nk.sync {
		ch, exists := nk.nodeWaiters[node.NodeID]
		if !exists {
			return
		}
		ch <- node
		close(ch)
		delete(nk.nodeWaiters, node.NodeID)
	}
	// then notify all the others that they have not passed the consensus
	for _, ch := range nk.nodeWaiters {
		close(ch)
	}
	// drop old waiters map and create new
	nk.nodeWaiters = make(map[core.RecordRef]chan *core.ActiveNode)
	nk.state = synced
	log.Infof("Sync success for pulse %d", nk.pulse)
}

func (nk *nodekeeper) collectUnsync(number core.PulseNumber) consensus.UnsyncList {
	nk.pulse = number
	nk.state = pulseSet

	for _, node := range nk.unsync {
		node.PulseNum = nk.pulse
	}
	tmp := nk.unsync
	nk.unsync = make([]*core.ActiveNode, 0)
	nk.unsyncList = NewUnsyncHolder(nk.pulse, tmp)
	if len(nk.listWaiters) == 0 {
		return nk.unsyncList
	}
	// notify waiters that new unsync holder is available for read
	for _, ch := range nk.listWaiters {
		ch <- nk.unsyncList
		close(ch)
	}
	nk.listWaiters = make([]chan *UnsyncList, 0)
	return nk.unsyncList
}

func calculateJetRole(role core.JetRole) core.NodeRole {
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

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
	"time"

	"github.com/insolar/insolar/configuration"
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

func createOrigin(configuration configuration.Configuration) (mutableNode, error) {
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
	err2 := conn.Close()
	if err2 != nil {
		log.Warn(err2)
	}
	if err != nil {
		return "", err
	}
	return address, nil
}

// NewNodeKeeper create new NodeKeeper
func NewNodeKeeper(origin core.Node) network.NodeKeeper {
	return &nodekeeper{
		origin:       origin,
		state:        undefined,
		active:       make(map[core.RecordRef]core.Node),
		sync:         make([]core.Node, 0),
		unsync:       make([]mutableNode, 0),
		listWaiters:  make([]chan *UnsyncList, 0),
		nodeWaiters:  make(map[core.RecordRef]chan core.Node),
		indexNode:    make(map[core.NodeRole][]core.RecordRef),
		indexShortID: make(map[uint32]core.Node),
	}
}

type nodekeeperState uint8

const (
	undefined = nodekeeperState(iota + 1)
	pulseSet
	synced
)

type nodekeeper struct {
	origin core.Node
	state  nodekeeperState
	pulse  core.PulseNumber

	activeLock   sync.RWMutex
	active       map[core.RecordRef]core.Node
	indexNode    map[core.NodeRole][]core.RecordRef
	indexShortID map[uint32]core.Node
	sync         []core.Node

	unsyncLock  sync.Mutex
	unsync      []mutableNode
	unsyncList  *UnsyncList
	listWaiters []chan *UnsyncList
	nodeWaiters map[core.RecordRef]chan core.Node
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

func (nk *nodekeeper) GetActiveNodeByShortID(shortID uint32) core.Node {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.indexShortID[shortID]
}

func (nk *nodekeeper) SetPulse(number core.PulseNumber) (bool, network.UnsyncList) {
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

func (nk *nodekeeper) Sync(syncCandidates []core.Node, number core.PulseNumber) {
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
		candidates += node.ID().String() + ", "
	}
	log.Debugf("Moving unsync to sync: %s", candidates)

	nk.syncUnsafe(syncCandidates)
}

func (nk *nodekeeper) AddUnsync(nodeID core.RecordRef, roles []core.NodeRole, address string,
	version string /*, publicKey *ecdsa.PublicKey*/) (chan core.Node, error) {

	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	node := newMutableNode(
		nodeID,
		roles,
		nil, // TODO publicKey
		nk.pulse,
		address,
		version,
	)

	nk.unsync = append(nk.unsync, node)
	ch := make(chan core.Node, 1)
	nk.nodeWaiters[node.ID()] = ch
	return ch, nil
}

func (nk *nodekeeper) GetUnsyncHolder(pulse core.PulseNumber, duration time.Duration) (network.UnsyncList, error) {
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

func (nk *nodekeeper) syncUnsafe(syncCandidates []core.Node) {
	// sync -> active
	for _, node := range nk.sync {
		nk.addActiveNode(node)
	}
	// unsync -> sync
	nk.sync = syncCandidates

	// first notify all synced nodes that they have passed the consensus
	for _, node := range nk.sync {
		ch, exists := nk.nodeWaiters[node.ID()]
		if !exists {
			return
		}
		ch <- node
		close(ch)
		delete(nk.nodeWaiters, node.ID())
	}
	// then notify all the others that they have not passed the consensus
	for _, ch := range nk.nodeWaiters {
		close(ch)
	}
	// drop old waiters map and create new
	nk.nodeWaiters = make(map[core.RecordRef]chan core.Node)
	nk.state = synced
	log.Infof("Sync success for pulse %d", nk.pulse)
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

func (nk *nodekeeper) collectUnsync(number core.PulseNumber) network.UnsyncList {
	nk.pulse = number
	nk.state = pulseSet

	for _, node := range nk.unsync {
		node.SetPulse(nk.pulse)
	}
	tmp := nk.unsync
	nk.unsync = make([]mutableNode, 0)

	unsyncNodes := mutableNodes(tmp).Export()

	nk.unsyncList = NewUnsyncHolder(nk.pulse, unsyncNodes)
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

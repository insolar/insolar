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
	"sort"
	"strings"
	"sync"

	"github.com/insolar/insolar/configuration"
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	coreutils "github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
)

// NewNodeNetwork create active node component
func NewNodeNetwork(configuration configuration.HostNetwork, certificate core.Certificate) (core.NodeNetwork, error) {
	origin, err := createOrigin(configuration, certificate)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create origin node")
	}
	nodeKeeper := NewNodeKeeper(origin)
	if len(certificate.GetDiscoveryNodes()) == 0 || utils.OriginIsDiscovery(certificate) {
		nodeKeeper.AddActiveNodes([]core.Node{origin})
	}
	return nodeKeeper, nil
}

func createOrigin(configuration configuration.HostNetwork, certificate core.Certificate) (MutableNode, error) {
	publicAddress, err := resolveAddress(configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to resolve public address")
	}

	role := certificate.GetRole()
	if role == core.StaticRoleUnknown {
		log.Info("[ createOrigin ] Use core.StaticRoleLightMaterial, since no role in certificate")
		role = core.StaticRoleLightMaterial
	}

	// TODO: get roles from certificate
	return newMutableNode(
		*certificate.GetNodeRef(),
		role,
		certificate.GetPublicKey(),
		publicAddress,
		version.Version,
	), nil
}

func resolveAddress(configuration configuration.HostNetwork) (string, error) {
	conn, address, err := transport.NewConnection(configuration.Transport)
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
		indexNode:    make(map[core.StaticRole]*recordRefSet),
		indexShortID: make(map[core.ShortNodeID]core.Node),
	}
}

type nodekeeper struct {
	origin     core.Node
	originLock sync.RWMutex
	state      network.NodeKeeperState
	claimQueue *claimQueue

	nodesJoinedDuringPrevPulse bool

	cloudHashLock sync.RWMutex
	cloudHash     []byte

	activeLock   sync.RWMutex
	active       map[core.RecordRef]core.Node
	indexNode    map[core.StaticRole]*recordRefSet
	indexShortID map[core.ShortNodeID]core.Node

	sync     network.UnsyncList
	syncLock sync.Mutex

	isBootstrap     bool
	isBootstrapLock sync.RWMutex

	Cryptography core.CryptographyService `inject:""`
}

// TODO: remove this method when bootstrap mechanism completed
// IsBootstrapped method returns true when bootstrapNodes are connected to each other
func (nk *nodekeeper) IsBootstrapped() bool {
	nk.isBootstrapLock.RLock()
	defer nk.isBootstrapLock.RUnlock()

	return nk.isBootstrap
}

// TODO: remove this method when bootstrap mechanism completed
// SetIsBootstrapped method set is bootstrap completed
func (nk *nodekeeper) SetIsBootstrapped(isBootstrap bool) {
	nk.isBootstrapLock.Lock()
	defer nk.isBootstrapLock.Unlock()

	nk.isBootstrap = isBootstrap
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
		return result[i].ID().Compare(result[j].ID()) < 0
	})
	return result
}

func (nk *nodekeeper) GetActiveNodesByRole(role core.DynamicRole) []core.RecordRef {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	list, exists := nk.indexNode[jetRoleToNodeRole(role)]
	if !exists {
		return nil
	}
	return list.Collect()
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

	list, ok := nk.indexNode[node.Role()]
	if !ok {
		list = newRecordRefSet()
	}
	list.Add(node.ID())
	nk.indexNode[node.Role()] = list

	nk.indexShortID[node.ShortID()] = node
}

func (nk *nodekeeper) delActiveNode(ref core.RecordRef) {
	if ref.Equal(nk.origin.ID()) {
		// we received acknowledge to leave, can gracefully stop

		// graceful stop instead of panic
		err := coreutils.SendGracefulStopSignal()
		if err != nil {
			// we tried :(
			panic("Node leave acknowledged by network. Goodbye!")
		}
	}
	active, ok := nk.active[ref]
	if !ok {
		return
	}
	delete(nk.active, ref)
	delete(nk.indexShortID, active.ShortID())
	nk.indexNode[active.Role()].Remove(ref)
}

func (nk *nodekeeper) SetState(state network.NodeKeeperState) {
	nk.state = state
}

func (nk *nodekeeper) GetState() network.NodeKeeperState {
	return nk.state
}

func (nk *nodekeeper) GetOriginClaim() (*consensus.NodeJoinClaim, error) {
	nk.originLock.RLock()
	defer nk.originLock.RUnlock()

	return nk.nodeToClaim()
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

func (nk *nodekeeper) GetUnsyncList() network.UnsyncList {
	return newUnsyncList(nk.GetActiveNodes())
}

func (nk *nodekeeper) GetSparseUnsyncList(length int) network.UnsyncList {
	return newSparseUnsyncList(length)
}

func (nk *nodekeeper) Sync(list network.UnsyncList) {
	nk.syncLock.Lock()
	defer nk.syncLock.Unlock()

	nk.sync = list
}

func (nk *nodekeeper) MoveSyncToActive() {
	nk.activeLock.Lock()
	nk.syncLock.Lock()
	defer func() {
		nk.syncLock.Unlock()
		nk.activeLock.Unlock()
	}()

	sync := nk.sync.(*unsyncList)
	sync.mergeWith(sync.claims, nk.addActiveNode, nk.delActiveNode)
}

func (nk *nodekeeper) nodeToClaim() (*consensus.NodeJoinClaim, error) {
	key, err := nk.Cryptography.GetPublicKey()
	if err != nil {
		return nil, errors.Wrap(err, "[ nodeToClaim ] failed to get a public key")
	}
	keyProc := platformpolicy.NewKeyProcessor()
	exportedKey, err := keyProc.ExportPublicKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "[ nodeToClaim ] failed to export a public key")
	}
	var keyData [consensus.PublicKeyLength]byte
	copy(keyData[:], exportedKey[:consensus.PublicKeyLength])

	var s [consensus.SignatureLength]byte
	claim := consensus.NodeJoinClaim{
		ShortNodeID:             nk.origin.ShortID(),
		RelayNodeID:             nk.origin.ShortID(),
		ProtocolVersionAndFlags: 0,
		JoinsAfter:              0,
		NodeRoleRecID:           0, // TODO: how to get a role as int?
		NodeRef:                 nk.origin.ID(),
		NodePK:                  keyData,
		Signature:               s,
	}

	dataToSign, err := claim.SerializeWithoutSign()
	if err != nil {
		return nil, errors.Wrap(err, "[ nodeToClaim ] failed to serialize a claim")
	}
	sign, err := nk.sign(dataToSign)
	if err != nil {
		return nil, errors.Wrap(err, "[ nodeToClaim ] failed to sign a claim")
	}

	copy(claim.Signature[:], sign[:consensus.SignatureLength])
	return &claim, nil
}

func (nk *nodekeeper) sign(data []byte) ([]byte, error) {
	sign, err := nk.Cryptography.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ sign ] failed to sign a claim")
	}
	return sign.Bytes(), nil
}

func jetRoleToNodeRole(role core.DynamicRole) core.StaticRole {
	switch role {
	case core.DynamicRoleVirtualExecutor:
		return core.StaticRoleVirtual
	case core.DynamicRoleVirtualValidator:
		return core.StaticRoleVirtual
	case core.DynamicRoleLightExecutor:
		return core.StaticRoleLightMaterial
	case core.DynamicRoleLightValidator:
		return core.StaticRoleLightMaterial
	case core.DynamicRoleHeavyExecutor:
		return core.StaticRoleHeavyMaterial
	default:
		return core.StaticRoleUnknown
	}
}

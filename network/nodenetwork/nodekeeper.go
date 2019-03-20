/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package nodenetwork

import (
	"context"
	"github.com/insolar/insolar/network/node"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/configuration"
	consensusMetrics "github.com/insolar/insolar/consensus"
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// NewNodeNetwork create active node component
func NewNodeNetwork(configuration configuration.HostNetwork, certificate core.Certificate) (core.NodeNetwork, error) {
	origin, err := createOrigin(configuration, certificate)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create origin node")
	}
	nodeKeeper := NewNodeKeeper(origin)
	if !utils.OriginIsDiscovery(certificate) {
		origin.(node.MutableNode).SetState(core.NodePending)
	}
	return nodeKeeper, nil
}

func createOrigin(configuration configuration.HostNetwork, certificate core.Certificate) (core.Node, error) {
	publicAddress, err := resolveAddress(configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to resolve public address")
	}

	role := certificate.GetRole()
	if role == core.StaticRoleUnknown {
		log.Info("[ createOrigin ] Use core.StaticRoleLightMaterial, since no role in certificate")
		role = core.StaticRoleLightMaterial
	}

	return node.NewNode(
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
	nk := &nodekeeper{
		origin:        origin,
		claimQueue:    newClaimQueue(),
		consensusInfo: newConsensusInfo(),
		syncNodes:     make([]core.Node, 0),
		syncClaims:    make([]consensus.ReferendumClaim, 0),
	}
	nk.SetInitialSnapshot([]core.Node{})
	return nk
}

type nodekeeper struct {
	origin        core.Node
	claimQueue    *claimQueue
	consensusInfo *consensusInfo

	cloudHashLock sync.RWMutex
	cloudHash     []byte

	activeLock sync.RWMutex
	snapshot   *node.Snapshot
	accessor   *node.Accessor

	syncLock   sync.Mutex
	syncNodes  []core.Node
	syncClaims []consensus.ReferendumClaim

	isBootstrap     bool
	isBootstrapLock sync.RWMutex

	Cryptography core.CryptographyService `inject:""`
}

func (nk *nodekeeper) SetInitialSnapshot(nodes []core.Node) {
	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	nodesMap := make(map[core.RecordRef]core.Node)
	for _, node := range nodes {
		nodesMap[node.ID()] = node
	}
	nk.snapshot = node.NewSnapshot(core.FirstPulseNumber, nodesMap)
	nk.accessor = node.NewAccessor(nk.snapshot)
	nk.syncNodes = nk.accessor.GetActiveNodes()
}

func (nk *nodekeeper) GetAccessor() network.Accessor {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.accessor
}

func (nk *nodekeeper) GetConsensusInfo() network.ConsensusInfo {
	return nk.consensusInfo
}

func (nk *nodekeeper) GetWorkingNode(ref core.RecordRef) core.Node {
	return nk.GetAccessor().GetWorkingNode(ref)
}

func (nk *nodekeeper) GetWorkingNodesByRole(role core.DynamicRole) []core.RecordRef {
	return nk.GetAccessor().GetWorkingNodesByRole(role)
}

func (nk *nodekeeper) Wipe(isDiscovery bool) {
	log.Warn("don't use it in production")

	nk.isBootstrapLock.Lock()
	nk.isBootstrap = false
	nk.isBootstrapLock.Unlock()

	nk.consensusInfo.flush(false)

	nk.cloudHashLock.Lock()
	nk.cloudHash = nil
	nk.cloudHashLock.Unlock()

	nk.SetInitialSnapshot([]core.Node{})

	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	nk.claimQueue = newClaimQueue()
	nk.syncLock.Lock()
	nk.syncNodes = make([]core.Node, 0)
	nk.syncClaims = make([]consensus.ReferendumClaim, 0)
	if isDiscovery {
		nk.origin.(node.MutableNode).SetState(core.NodeReady)
	}
	nk.syncLock.Unlock()
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

func (nk *nodekeeper) GetWorkingNodes() []core.Node {
	return nk.GetAccessor().GetWorkingNodes()
}

func (nk *nodekeeper) GetOriginJoinClaim() (*consensus.NodeJoinClaim, error) {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.nodeToSignedClaim()
}

func (nk *nodekeeper) GetOriginAnnounceClaim(mapper consensus.BitSetMapper) (*consensus.NodeAnnounceClaim, error) {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.nodeToAnnounceClaim(mapper)
}

func (nk *nodekeeper) AddPendingClaim(claim consensus.ReferendumClaim) bool {
	nk.claimQueue.Push(claim)
	return true
}

func (nk *nodekeeper) GetClaimQueue() network.ClaimQueue {
	return nk.claimQueue
}

func (nk *nodekeeper) GetUnsyncList() network.UnsyncList {
	activeNodes := nk.GetAccessor().GetActiveNodes()
	return newUnsyncList(nk.origin, activeNodes, len(activeNodes))
}

func (nk *nodekeeper) GetSparseUnsyncList(length int) network.UnsyncList {
	return newUnsyncList(nk.origin, nil, length)
}

func (nk *nodekeeper) Sync(ctx context.Context, nodes []core.Node, claims []consensus.ReferendumClaim) error {
	nk.syncLock.Lock()
	defer nk.syncLock.Unlock()

	inslogger.FromContext(ctx).Debugf("Sync, nodes: %d, claims: %d", len(nodes), len(claims))
	nk.syncNodes = nodes
	nk.syncClaims = claims

	foundOrigin := false
	for _, node := range nodes {
		if node.ID().Equal(nk.origin.ID()) {
			foundOrigin = true
			nk.syncOrigin(node)
			nk.consensusInfo.SetIsJoiner(false)
		}
	}

	if nk.shouldExit(foundOrigin) {
		return errors.New("node leave acknowledged by network")
	}

	return nil
}

// syncOrigin synchronize data in origin node with node from active list in case when they are different objects
func (nk *nodekeeper) syncOrigin(n core.Node) {
	if nk.origin == n {
		return
	}
	mutableOrigin := nk.origin.(node.MutableNode)
	mutableOrigin.SetState(n.GetState())
	mutableOrigin.SetLeavingETA(n.LeavingETA())
	mutableOrigin.SetShortID(n.ShortID())
}

func (nk *nodekeeper) MoveSyncToActive(ctx context.Context) error {
	nk.activeLock.Lock()
	nk.syncLock.Lock()
	defer func() {
		nk.syncLock.Unlock()
		nk.activeLock.Unlock()
	}()

	mergeResult, err := GetMergedCopy(nk.syncNodes, nk.syncClaims)
	if err != nil {
		return errors.Wrap(err, "[ MoveSyncToActive ] Failed to calculate new active list")
	}
	inslogger.FromContext(ctx).Infof("[ MoveSyncToActive ] New active list confirmed. Active list size: %d -> %d",
		len(nk.accessor.GetActiveNodes()), len(mergeResult.ActiveList))

	nk.snapshot = node.NewSnapshot(core.PulseNumber(0), mergeResult.ActiveList)
	nk.accessor = node.NewAccessor(nk.snapshot)
	stats.Record(ctx, consensusMetrics.ActiveNodes.M(int64(len(nk.accessor.GetActiveNodes()))))
	nk.consensusInfo.flush(mergeResult.NodesJoinedDuringPrevPulse)
	return nil
}

func (nk *nodekeeper) shouldExit(foundOrigin bool) bool {
	return !foundOrigin && nk.origin.GetState() == core.NodeReady && len(nk.GetAccessor().GetActiveNodes()) != 0
}

func (nk *nodekeeper) nodeToSignedClaim() (*consensus.NodeJoinClaim, error) {
	claim, err := consensus.NodeToClaim(nk.origin)
	if err != nil {
		return nil, err
	}
	dataToSign, err := claim.SerializeRaw()
	log.Debugf("dataToSign len: %d", len(dataToSign))
	if err != nil {
		return nil, errors.Wrap(err, "[ nodeToSignedClaim ] failed to serialize a claim")
	}
	sign, err := nk.sign(dataToSign)
	log.Debugf("sign len: %d", len(sign))
	if err != nil {
		return nil, errors.Wrap(err, "[ nodeToSignedClaim ] failed to sign a claim")
	}
	copy(claim.Signature[:], sign[:consensus.SignatureLength])
	return claim, nil
}

func (nk *nodekeeper) nodeToAnnounceClaim(mapper consensus.BitSetMapper) (*consensus.NodeAnnounceClaim, error) {
	claim := consensus.NodeAnnounceClaim{}
	joinClaim, err := consensus.NodeToClaim(nk.origin)
	if err != nil {
		return nil, err
	}
	claim.NodeJoinClaim = *joinClaim
	claim.NodeCount = uint16(mapper.Length())
	announcerIndex, err := mapper.RefToIndex(nk.origin.ID())
	if err != nil {
		return nil, errors.Wrap(err, "[ nodeToAnnounceClaim ] failed to map origin node ID to bitset index")
	}
	claim.NodeAnnouncerIndex = uint16(announcerIndex)
	claim.BitSetMapper = mapper
	claim.SetCloudHash(nk.GetCloudHash())
	return &claim, nil
}

func (nk *nodekeeper) sign(data []byte) ([]byte, error) {
	sign, err := nk.Cryptography.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ sign ] failed to sign a claim")
	}
	return sign.Bytes(), nil
}

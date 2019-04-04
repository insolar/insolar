//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package nodenetwork

import (
	"context"
	"net"
	"sync"

	"github.com/insolar/insolar/network/node"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/configuration"
	consensusMetrics "github.com/insolar/insolar/consensus"
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// NewNodeNetwork create active node component
func NewNodeNetwork(configuration configuration.HostNetwork, certificate insolar.Certificate) (insolar.NodeNetwork, error) {
	origin, err := createOrigin(configuration, certificate)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create origin node")
	}
	nodeKeeper := NewNodeKeeper(origin)
	if !utils.OriginIsDiscovery(certificate) {
		origin.(node.MutableNode).SetState(insolar.NodePending)
	}
	return nodeKeeper, nil
}

func createOrigin(configuration configuration.HostNetwork, certificate insolar.Certificate) (insolar.NetworkNode, error) {
	publicAddress, err := resolveAddress(configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to resolve public address")
	}

	role := certificate.GetRole()
	if role == insolar.StaticRoleUnknown {
		log.Info("[ createOrigin ] Use insolar.StaticRoleLightMaterial, since no role in certificate")
		role = insolar.StaticRoleLightMaterial
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
	addr, err := net.ResolveTCPAddr("tcp", configuration.Transport.Address)
	if err != nil {
		return "", err
	}
	address, err := transport.Resolve(configuration.Transport.FixedPublicAddress, addr.String())
	if err != nil {
		return "", err
	}
	return address, nil
}

// NewNodeKeeper create new NodeKeeper
func NewNodeKeeper(origin insolar.NetworkNode) network.NodeKeeper {
	nk := &nodekeeper{
		origin:        origin,
		claimQueue:    newClaimQueue(),
		consensusInfo: newConsensusInfo(),
		syncNodes:     make([]insolar.NetworkNode, 0),
		syncClaims:    make([]consensus.ReferendumClaim, 0),
	}
	nk.SetInitialSnapshot([]insolar.NetworkNode{})
	return nk
}

type nodekeeper struct {
	origin        insolar.NetworkNode
	claimQueue    *claimQueue
	consensusInfo *consensusInfo

	cloudHashLock sync.RWMutex
	cloudHash     []byte

	activeLock sync.RWMutex
	snapshot   *node.Snapshot
	accessor   *node.Accessor

	syncLock   sync.Mutex
	syncNodes  []insolar.NetworkNode
	syncClaims []consensus.ReferendumClaim

	isBootstrap     bool
	isBootstrapLock sync.RWMutex

	Cryptography       insolar.CryptographyService `inject:""`
	TerminationHandler insolar.TerminationHandler  `inject:""`

	gateway   network.Gateway
	gatewayMu sync.RWMutex
}

func (nk *nodekeeper) Gateway() network.Gateway {
	nk.gatewayMu.RLock()
	defer nk.gatewayMu.RUnlock()
	return nk.gateway
}

func (nk *nodekeeper) SetGateway(g network.Gateway) {
	nk.gatewayMu.Lock()
	defer nk.gatewayMu.Unlock()
	nk.gateway = g
}

func (nk *nodekeeper) GetSnapshotCopy() *node.Snapshot {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.snapshot.Copy()
}

func (nk *nodekeeper) SetInitialSnapshot(nodes []insolar.NetworkNode) {
	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	nodesMap := make(map[insolar.Reference]insolar.NetworkNode)
	for _, node := range nodes {
		nodesMap[node.ID()] = node
	}
	nk.snapshot = node.NewSnapshot(insolar.FirstPulseNumber, nodesMap)
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

func (nk *nodekeeper) GetWorkingNode(ref insolar.Reference) insolar.NetworkNode {
	return nk.GetAccessor().GetWorkingNode(ref)
}

func (nk *nodekeeper) GetWorkingNodesByRole(role insolar.DynamicRole) []insolar.Reference {
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

	nk.SetInitialSnapshot([]insolar.NetworkNode{})

	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	nk.claimQueue = newClaimQueue()
	nk.syncLock.Lock()
	nk.syncNodes = make([]insolar.NetworkNode, 0)
	nk.syncClaims = make([]consensus.ReferendumClaim, 0)
	if isDiscovery {
		nk.origin.(node.MutableNode).SetState(insolar.NodeReady)
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

func (nk *nodekeeper) GetOrigin() insolar.NetworkNode {
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

func (nk *nodekeeper) GetWorkingNodes() []insolar.NetworkNode {
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

func (nk *nodekeeper) Sync(ctx context.Context, nodes []insolar.NetworkNode, claims []consensus.ReferendumClaim) error {
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
func (nk *nodekeeper) syncOrigin(n insolar.NetworkNode) {
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

	nk.snapshot = node.NewSnapshot(insolar.PulseNumber(0), mergeResult.ActiveList)
	nk.accessor = node.NewAccessor(nk.snapshot)
	stats.Record(ctx, consensusMetrics.ActiveNodes.M(int64(len(nk.accessor.GetActiveNodes()))))
	nk.consensusInfo.flush(mergeResult.NodesJoinedDuringPrevPulse)
	nk.gracefulStopIfNeeded(ctx)
	return nil
}

func (nk *nodekeeper) gracefulStopIfNeeded(ctx context.Context) {
	if nk.origin.GetState() == insolar.NodeLeaving {
		nk.TerminationHandler.OnLeaveApproved(ctx)
	}
}

func (nk *nodekeeper) shouldExit(foundOrigin bool) bool {
	return !foundOrigin && nk.origin.GetState() == insolar.NodeReady && len(nk.GetAccessor().GetActiveNodes()) != 0
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

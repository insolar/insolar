/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package servicenetwork

import (
	"context"
	"time"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
)

type nodeKeeperTestInterface interface {
	// Wipe all active nodes for test purposes
	Wipe(isDiscovery bool)
}

type nodeKeeperWrapper struct {
	original network.NodeKeeper
}

func (n *nodeKeeperWrapper) GetWorkingNode(ref core.RecordRef) core.Node {
	return n.original.GetWorkingNode(ref)
}

func (n *nodeKeeperWrapper) GetWorkingNodes() []core.Node {
	return n.original.GetWorkingNodes()
}

func (n *nodeKeeperWrapper) GetWorkingNodesByRole(role core.DynamicRole) []core.RecordRef {
	return n.original.GetWorkingNodesByRole(role)
}

func (n *nodeKeeperWrapper) Wipe(isDiscovery bool) {
	n.original.(nodeKeeperTestInterface).Wipe(isDiscovery)
}

func (n *nodeKeeperWrapper) AddTemporaryMapping(nodeID core.RecordRef, shortID core.ShortNodeID, address string) error {
	return n.original.AddTemporaryMapping(nodeID, shortID, address)
}

func (n *nodeKeeperWrapper) ResolveConsensus(shortID core.ShortNodeID) *host.Host {
	return n.original.ResolveConsensus(shortID)
}

func (n *nodeKeeperWrapper) ResolveConsensusRef(nodeID core.RecordRef) *host.Host {
	return n.original.ResolveConsensusRef(nodeID)
}

type phaseManagerWrapper struct {
	original phases.PhaseManager
	result   chan error
}

func (p *phaseManagerWrapper) OnPulse(ctx context.Context, pulse *core.Pulse, pulseStartTime time.Time) error {
	res := p.original.OnPulse(ctx, pulse, pulseStartTime)
	p.result <- res
	return res
}

func (n *nodeKeeperWrapper) GetOrigin() core.Node {
	return n.original.GetOrigin()
}

func (n *nodeKeeperWrapper) GetActiveNode(ref core.RecordRef) core.Node {
	return n.original.GetActiveNode(ref)
}

func (n *nodeKeeperWrapper) GetActiveNodes() []core.Node {
	tmp := n.original.GetActiveNodes()
	return tmp
}

func (n *nodeKeeperWrapper) GetCloudHash() []byte {
	return n.original.GetCloudHash()
}

func (n *nodeKeeperWrapper) IsBootstrapped() bool {
	return n.original.IsBootstrapped()
}

func (n *nodeKeeperWrapper) SetIsBootstrapped(isBootstrap bool) {
	n.original.SetIsBootstrapped(isBootstrap)
}

func (n *nodeKeeperWrapper) SetCloudHash(hash []byte) {
	n.original.SetCloudHash(hash)
}

func (n *nodeKeeperWrapper) AddActiveNodes(nodes []core.Node) {
	n.original.AddActiveNodes(nodes)
}

func (n *nodeKeeperWrapper) GetActiveNodeByShortID(shortID core.ShortNodeID) core.Node {
	return n.original.GetActiveNodeByShortID(shortID)
}

func (n *nodeKeeperWrapper) SetState(state core.NodeNetworkState) {
	n.original.SetState(state)
}

func (n *nodeKeeperWrapper) GetState() core.NodeNetworkState {
	return n.original.GetState()
}

func (n *nodeKeeperWrapper) GetOriginJoinClaim() (*consensus.NodeJoinClaim, error) {
	return n.original.GetOriginJoinClaim()
}

func (n *nodeKeeperWrapper) GetOriginAnnounceClaim(mapper consensus.BitSetMapper) (*consensus.NodeAnnounceClaim, error) {
	return n.original.GetOriginAnnounceClaim(mapper)
}

func (n *nodeKeeperWrapper) NodesJoinedDuringPreviousPulse() bool {
	return n.original.NodesJoinedDuringPreviousPulse()
}

func (n *nodeKeeperWrapper) AddPendingClaim(claim consensus.ReferendumClaim) bool {
	return n.original.AddPendingClaim(claim)
}

func (n *nodeKeeperWrapper) GetClaimQueue() network.ClaimQueue {
	return n.original.GetClaimQueue()
}

func (n *nodeKeeperWrapper) GetUnsyncList() network.UnsyncList {
	return n.original.GetUnsyncList()
}

func (n *nodeKeeperWrapper) GetSparseUnsyncList(length int) network.UnsyncList {
	return n.original.GetSparseUnsyncList(length)
}

func (n *nodeKeeperWrapper) Sync(ctx context.Context, nodes []core.Node, claims []consensus.ReferendumClaim) error {
	return n.original.Sync(ctx, nodes, claims)
}

func (n *nodeKeeperWrapper) MoveSyncToActive(ctx context.Context) error {
	return n.original.MoveSyncToActive(ctx)
}

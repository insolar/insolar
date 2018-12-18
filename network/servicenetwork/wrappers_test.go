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

package servicenetwork

import (
	"context"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
)

type nodeKeeperWrapper struct {
	original network.NodeKeeper
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

func (p *phaseManagerWrapper) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	res := p.original.OnPulse(ctx, pulse)
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
	tmp = tmp[:len(tmp)-2]
	return tmp
}

func (n *nodeKeeperWrapper) GetActiveNodesByRole(role core.DynamicRole) []core.RecordRef {
	return n.original.GetActiveNodesByRole(role)
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

func (n *nodeKeeperWrapper) SetState(state network.NodeKeeperState) {
	n.original.SetState(state)
}

func (n *nodeKeeperWrapper) GetState() network.NodeKeeperState {
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

func (n *nodeKeeperWrapper) Sync(list network.UnsyncList) {
	n.original.Sync(list)
}

func (n *nodeKeeperWrapper) MoveSyncToActive() error {
	return n.original.MoveSyncToActive()
}

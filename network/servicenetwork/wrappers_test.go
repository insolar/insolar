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

// +build networktest

package servicenetwork

import (
	"context"
	"time"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
)

type nodeKeeperWrapper struct {
	original network.NodeKeeper
}

func (n *nodeKeeperWrapper) GetSnapshotCopy() *node.Snapshot {
	return n.original.GetSnapshotCopy()
}

func (n *nodeKeeperWrapper) GetAccessor() network.Accessor {
	return n.original.GetAccessor()
}

func (n *nodeKeeperWrapper) GetConsensusInfo() network.ConsensusInfo {
	return n.original.GetConsensusInfo()
}

func (n *nodeKeeperWrapper) GetWorkingNode(ref insolar.Reference) insolar.NetworkNode {
	return n.original.GetWorkingNode(ref)
}

func (n *nodeKeeperWrapper) GetWorkingNodes() []insolar.NetworkNode {
	return n.original.GetWorkingNodes()
}

func (n *nodeKeeperWrapper) GetWorkingNodesByRole(role insolar.DynamicRole) []insolar.Reference {
	return n.original.GetWorkingNodesByRole(role)
}

func (n *nodeKeeperWrapper) GetOrigin() insolar.NetworkNode {
	return n.original.GetOrigin()
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

func (n *nodeKeeperWrapper) SetInitialSnapshot(nodes []insolar.NetworkNode) {
	n.original.SetInitialSnapshot(nodes)
}

func (n *nodeKeeperWrapper) GetOriginJoinClaim() (*consensus.NodeJoinClaim, error) {
	return n.original.GetOriginJoinClaim()
}

func (n *nodeKeeperWrapper) GetOriginAnnounceClaim(mapper consensus.BitSetMapper) (*consensus.NodeAnnounceClaim, error) {
	return n.original.GetOriginAnnounceClaim(mapper)
}

func (n *nodeKeeperWrapper) GetClaimQueue() network.ClaimQueue {
	return n.original.GetClaimQueue()
}

func (n *nodeKeeperWrapper) Sync(ctx context.Context, nodes []insolar.NetworkNode, claims []consensus.ReferendumClaim) error {
	return n.original.Sync(ctx, nodes, claims)
}

func (n *nodeKeeperWrapper) MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error {
	return n.original.MoveSyncToActive(ctx, number)
}

type phaseManagerWrapper struct {
	original phases.PhaseManager
	result   chan error
}

func (p *phaseManagerWrapper) OnPulse(ctx context.Context, pulse *insolar.Pulse, pulseStartTime time.Time) error {
	res := p.original.OnPulse(ctx, pulse, pulseStartTime)
	p.result <- res
	return res
}

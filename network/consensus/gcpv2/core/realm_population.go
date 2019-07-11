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

package core

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
)

type RealmPopulation interface {
	GetNodeCount() int
	GetOthersCount() int
	GetJoinersCount() int
	GetBftMajorityCount() int
	IsComplete() bool

	GetNodeAppearance(id insolar.ShortNodeID) *NodeAppearance
	GetActiveNodeAppearance(id insolar.ShortNodeID) *NodeAppearance
	GetJoinerNodeAppearance(id insolar.ShortNodeID) *NodeAppearance
	GetNodeAppearanceByIndex(idx int) *NodeAppearance

	GetShuffledOtherNodes() []*NodeAppearance /* no joiners or self included */
	GetIndexedNodes() []*NodeAppearance       /* no joiners included */

	GetSelf() *NodeAppearance

	CreateNodeAppearance(ctx context.Context, inp gcp_types.NodeProfile) *NodeAppearance

	AddToPurgatory(n *NodeAppearance) (*NodeAppearance, PurgatoryNodeState)
	AddToDynamics(n *NodeAppearance) (*NodeAppearance, []*NodeAppearance)

	SetOrUpdateVectorHelper(v *RealmVectorHelper) *RealmVectorHelper
}

type PurgatoryNodeState int

const PurgatoryDuplicatePK PurgatoryNodeState = -1
const PurgatoryExistingMember PurgatoryNodeState = -2

type NodeInitFunc func(ctx context.Context, n *NodeAppearance)

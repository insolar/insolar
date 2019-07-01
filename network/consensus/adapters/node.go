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

package adapters

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/utils"
)

type NodeIntroduction struct {
	node insolar.NetworkNode
}

func NewNodeIntroduction(node insolar.NetworkNode) *NodeIntroduction {
	return &NodeIntroduction{
		node: node,
	}
}

func (ni *NodeIntroduction) GetShortNodeID() common.ShortNodeID {
	return common.ShortNodeID(ni.node.ShortID())
}

func (ni *NodeIntroduction) GetClaimEvidence() common.SignedEvidenceHolder {
	// TODO: do something with sign
	return nil
}

type NodeIntroProfile struct {
	node        insolar.NetworkNode
	isDiscovery bool
}

func NewNodeIntroProfile(node insolar.NetworkNode, certificate insolar.Certificate) *NodeIntroProfile {
	return &NodeIntroProfile{
		node:        node,
		isDiscovery: utils.IsDiscovery(node.ID(), certificate),
	}
}

func (nip *NodeIntroProfile) GetPrimaryRole() common2.NodePrimaryRole {
	return StaticRoleToPrimaryRole(nip.node.Role())
}

func (nip *NodeIntroProfile) GetSpecialRoles() common2.NodeSpecialRole {
	if nip.isDiscovery {
		return common2.SpecialRoleDiscovery
	}

	return common2.SpecialRoleNoRole
}

func (nip *NodeIntroProfile) IsAllowedPower(p common2.MemberPower) bool {
	// TODO: do something with power
	return true
}

func (nip *NodeIntroProfile) GetIntroduction() common2.NodeIntroduction {
	return NewNodeIntroduction(nip.node)
}

func (nip *NodeIntroProfile) GetDefaultEndpoint() common.HostAddress {
	return common.HostAddress(nip.node.Address())
}

func (nip *NodeIntroProfile) GetNodePublicKeyStore() common.PublicKeyStore {
	ecdsaPublicKey := nip.node.PublicKey().(*ecdsa.PublicKey)
	return NewECDSAPublicKeyStore(ecdsaPublicKey)
}

func (nip *NodeIntroProfile) IsAcceptableHost(from common.HostIdentityHolder) bool {
	endpoint := nip.GetDefaultEndpoint()
	return endpoint.Equals(from.GetHostAddress())
}

func (nip *NodeIntroProfile) GetShortNodeID() common.ShortNodeID {
	return common.ShortNodeID(nip.node.ShortID())
}

func (nip *NodeIntroProfile) String() string {
	return fmt.Sprintf("{sid:%d, node:%s}", nip.node.ShortID(), nip.node.ID().String())
}

func NewNodeIntroProfileList(nodes []insolar.NetworkNode, certificate insolar.Certificate) []common2.NodeIntroProfile {
	intros := make([]common2.NodeIntroProfile, len(nodes))
	for i, n := range nodes {
		intros[i] = NewNodeIntroProfile(n, certificate)
	}

	return intros
}

func NewNetworkNode(profile common2.NodeProfile) insolar.NetworkNode {
	intro := common2.NodeIntroProfile(profile)
	nn := intro.GetIntroduction().(*NodeIntroduction).node

	networkNode := node.NewNode(
		nn.ID(),
		PrimaryRoleToStaticRole(profile.GetPrimaryRole()),
		nn.PublicKey(),
		profile.GetDefaultEndpoint().String(),
		nn.Version(),
	)

	mutableNode := networkNode.(node.MutableNode)

	mutableNode.SetShortID(nn.ShortID())
	mutableNode.SetState(MembershipStateToNodeState(profile.GetState()))

	return networkNode
}

func NewNetworkNodeList(profiles []common2.NodeProfile) []insolar.NetworkNode {
	networkNodes := make([]insolar.NetworkNode, len(profiles))
	for i, p := range profiles {
		networkNodes[i] = NewNetworkNode(p)
	}

	return networkNodes
}

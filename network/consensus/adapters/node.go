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
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensusv1/packets"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/utils"
)

type NodeIntroduction struct {
	shortID insolar.ShortNodeID
	ref     insolar.Reference
}

func NewNodeIntroduction(networkNode insolar.NetworkNode) *NodeIntroduction {
	return newNodeIntroduction(
		insolar.ShortNodeID(networkNode.ShortID()),
		networkNode.ID(),
	)
}

func newNodeIntroduction(shortID insolar.ShortNodeID, ref insolar.Reference) *NodeIntroduction {
	return &NodeIntroduction{
		shortID: shortID,
		ref:     ref,
	}
}

func (ni *NodeIntroduction) ConvertPowerRequest(request gcp_types.PowerRequest) gcp_types.MemberPower {
	if ok, cl := request.AsCapacityLevel(); ok {
		return gcp_types.MemberPowerOf(uint16(cl.DefaultPercent()))
	}
	_, pw := request.AsMemberPower()
	return pw
}

func (ni *NodeIntroduction) GetNodeReference() insolar.Reference {
	return ni.ref
}

func (ni *NodeIntroduction) IsAllowedPower(p gcp_types.MemberPower) bool {
	// TODO: do something with power
	return true
}

func (ni *NodeIntroduction) GetShortNodeID() insolar.ShortNodeID {
	return ni.shortID
}

type NodeIntroProfile struct {
	shortID     insolar.ShortNodeID
	primaryRole gcp_types.NodePrimaryRole
	specialRole gcp_types.NodeSpecialRole
	intro       gcp_types.NodeIntroduction
	endpoint    endpoints.NodeEndpoint
	store       cryptography_containers.PublicKeyStore
	keyHolder   cryptography_containers.SignatureKeyHolder

	signature cryptography_containers.SignatureHolder
}

func NewNodeIntroProfile(networkNode insolar.NetworkNode, certificate insolar.Certificate, keyProcessor insolar.KeyProcessor) *NodeIntroProfile {
	specialRole := gcp_types.SpecialRoleNone
	if utils.IsDiscovery(networkNode.ID(), certificate) {
		specialRole = gcp_types.SpecialRoleDiscovery
	}

	publicKey := networkNode.PublicKey().(*ecdsa.PublicKey)
	mutableNode := networkNode.(node.MutableNode)
	signature := mutableNode.GetSignature()

	return newNodeIntroProfile(
		insolar.ShortNodeID(networkNode.ShortID()),
		StaticRoleToPrimaryRole(networkNode.Role()),
		specialRole,
		NewNodeIntroduction(networkNode),
		NewNodeEndpoint(networkNode.Address()),
		NewECDSAPublicKeyStore(publicKey),
		NewECDSASignatureKeyHolder(publicKey, keyProcessor),
		cryptography_containers.NewSignature(
			long_bits.NewBits512FromBytes(signature.Bytes()),
			SHA3512Digest.SignedBy(SECP256r1Sign),
		).AsSignatureHolder(),
	)
}

func newNodeIntroProfile(
	shortID insolar.ShortNodeID,
	primaryRole gcp_types.NodePrimaryRole,
	specialRole gcp_types.NodeSpecialRole,
	intro gcp_types.NodeIntroduction,
	endpoint endpoints.NodeEndpoint,
	store cryptography_containers.PublicKeyStore,
	keyHolder cryptography_containers.SignatureKeyHolder,
	signature cryptography_containers.SignatureHolder,
) *NodeIntroProfile {
	return &NodeIntroProfile{
		shortID:     shortID,
		primaryRole: primaryRole,
		specialRole: specialRole,
		intro:       intro,
		endpoint:    endpoint,
		store:       store,
		keyHolder:   keyHolder,
		signature:   signature,
	}
}

func (nip *NodeIntroProfile) GetPrimaryRole() gcp_types.NodePrimaryRole {
	return nip.primaryRole
}

func (nip *NodeIntroProfile) GetSpecialRoles() gcp_types.NodeSpecialRole {
	return nip.specialRole
}

func (nip *NodeIntroProfile) HasIntroduction() bool {
	return nip.intro != nil
}

func (nip *NodeIntroProfile) GetIntroduction() gcp_types.NodeIntroduction {
	return nip.intro
}

func (nip *NodeIntroProfile) GetDefaultEndpoint() endpoints.NodeEndpoint {
	return nip.endpoint
}

func (nip *NodeIntroProfile) GetNodePublicKeyStore() cryptography_containers.PublicKeyStore {
	return nip.store
}

func (nip *NodeIntroProfile) GetNodePublicKey() cryptography_containers.SignatureKeyHolder {
	return nip.keyHolder
}

func (nip *NodeIntroProfile) GetStartPower() gcp_types.MemberPower {
	// TODO: get from certificate
	return 10
}

func (nip *NodeIntroProfile) IsAcceptableHost(from endpoints.HostIdentityHolder) bool {
	address := nip.endpoint.GetNameAddress()
	return address.Equals(from.GetHostAddress())
}

func (nip *NodeIntroProfile) GetShortNodeID() insolar.ShortNodeID {
	return nip.shortID
}

func (nip *NodeIntroProfile) GetAnnouncementSignature() cryptography_containers.SignatureHolder {
	return nip.signature
}

func (nip *NodeIntroProfile) String() string {
	return fmt.Sprintf("{sid:%d, node:%s}", nip.shortID, nip.intro.GetNodeReference().String())
}

type NodeEndpoint struct {
	name endpoints.HostAddress
	addr packets.NodeAddress
}

func NewNodeEndpoint(address string) *NodeEndpoint {
	addr, err := packets.NewNodeAddress(address)
	if err != nil {
		panic(err)
	}

	return &NodeEndpoint{
		name: endpoints.HostAddress(address),
		addr: addr,
	}
}

func (p *NodeEndpoint) GetEndpointType() endpoints.NodeEndpointType {
	return endpoints.NameEndpoint
}

func (*NodeEndpoint) GetRelayID() insolar.ShortNodeID {
	return 0
}

func (p *NodeEndpoint) GetNameAddress() endpoints.HostAddress {
	return p.name
}

func (p *NodeEndpoint) GetIPAddress() packets.NodeAddress {
	return p.addr
}

func NewNodeIntroProfileList(nodes []insolar.NetworkNode, certificate insolar.Certificate, keyProcessor insolar.KeyProcessor) []gcp_types.NodeIntroProfile {
	intros := make([]gcp_types.NodeIntroProfile, len(nodes))
	for i, n := range nodes {
		intros[i] = NewNodeIntroProfile(n, certificate, keyProcessor)
	}

	return intros
}

func NewNetworkNode(profile gcp_types.NodeProfile) insolar.NetworkNode {
	store := profile.GetNodePublicKeyStore()
	introduction := profile.GetIntroduction()

	networkNode := node.NewNode(
		introduction.GetNodeReference(),
		PrimaryRoleToStaticRole(profile.GetPrimaryRole()),
		store.(*ECDSAPublicKeyStore).publicKey,
		profile.GetDefaultEndpoint().GetNameAddress().String(),
		"",
	)

	mutableNode := networkNode.(node.MutableNode)

	mutableNode.SetShortID(insolar.ShortNodeID(profile.GetShortNodeID()))
	mutableNode.SetState(insolar.NodeReady)
	mutableNode.SetSignature(insolar.SignatureFromBytes(profile.GetAnnouncementSignature().AsBytes()))

	return networkNode
}

func NewNetworkNodeList(profiles []gcp_types.NodeProfile) []insolar.NetworkNode {
	networkNodes := make([]insolar.NetworkNode, len(profiles))
	for i, p := range profiles {
		networkNodes[i] = NewNetworkNode(p)
	}

	return networkNodes
}

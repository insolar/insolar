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
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/pulse"
)

type StaticProfileExtension struct {
	shortID   insolar.ShortNodeID
	ref       insolar.Reference
	signature cryptkit.SignatureHolder
}

func NewStaticProfileExtension(networkNode insolar.NetworkNode) *StaticProfileExtension {
	_, signature := networkNode.(node.MutableNode).GetSignature()

	return newStaticProfileExtension(
		networkNode.ShortID(),
		networkNode.ID(),
		cryptkit.NewSignature(
			longbits.NewBits512FromBytes(signature.Bytes()),
			SHA3512Digest.SignedBy(SECP256r1Sign),
		).AsSignatureHolder(),
	)
}

func newStaticProfileExtension(shortID insolar.ShortNodeID, ref insolar.Reference, signature cryptkit.SignatureHolder) *StaticProfileExtension {
	return &StaticProfileExtension{
		shortID:   shortID,
		ref:       ref,
		signature: signature,
	}
}

func (ni *StaticProfileExtension) GetPowerLevels() member.PowerSet {
	return member.PowerSet{0, 0, 0, 0xff}
}

func (ni *StaticProfileExtension) GetIntroducedNodeID() insolar.ShortNodeID {
	return ni.shortID
}

func (ni *StaticProfileExtension) GetExtraEndpoints() []endpoints.Outbound {
	return nil
}

func (ni *StaticProfileExtension) GetIssuedAtPulse() pulse.Number {
	return pulse.NewFirstEphemeralData().PulseNumber
}

func (ni *StaticProfileExtension) GetIssuedAtTime() time.Time {
	return time.Unix(int64(pulse.NewFirstEphemeralData().Timestamp), 0)
}

func (ni *StaticProfileExtension) GetIssuerID() insolar.ShortNodeID {
	return ni.shortID
}

func (ni *StaticProfileExtension) GetIssuerSignature() cryptkit.SignatureHolder {
	return ni.signature
}

func (ni *StaticProfileExtension) GetReference() insolar.Reference {
	return ni.ref
}

type StaticProfile struct {
	shortID     insolar.ShortNodeID
	primaryRole member.PrimaryRole
	specialRole member.SpecialRole
	intro       profiles.StaticProfileExtension
	endpoint    endpoints.Outbound
	store       cryptkit.PublicKeyStore
	keyHolder   cryptkit.SignatureKeyHolder

	signature cryptkit.SignedDigestHolder
}

func NewStaticProfile(networkNode insolar.NetworkNode, certificate insolar.Certificate, keyProcessor insolar.KeyProcessor) *StaticProfile {

	specialRole := member.SpecialRoleNone
	if network.IsDiscovery(networkNode.ID(), certificate) {
		specialRole = member.SpecialRoleDiscovery
	}

	publicKey := networkNode.PublicKey().(*ecdsa.PublicKey)
	mutableNode := networkNode.(node.MutableNode)
	digest, signature := mutableNode.GetSignature()

	return newStaticProfile(
		networkNode.ShortID(),
		StaticRoleToPrimaryRole(networkNode.Role()),
		specialRole,
		NewStaticProfileExtension(networkNode),
		NewOutbound(networkNode.Address()),
		NewECDSAPublicKeyStore(publicKey),
		NewECDSASignatureKeyHolder(publicKey, keyProcessor),
		cryptkit.NewSignedDigest(
			cryptkit.NewDigest(longbits.NewBits512FromBytes(digest), SHA3512Digest),
			cryptkit.NewSignature(longbits.NewBits512FromBytes(signature.Bytes()), SHA3512Digest.SignedBy(SECP256r1Sign)),
		).AsSignedDigestHolder(),
	)
}

func newStaticProfile(
	shortID insolar.ShortNodeID,
	primaryRole member.PrimaryRole,
	specialRole member.SpecialRole,
	intro profiles.StaticProfileExtension,
	endpoint endpoints.Outbound,
	store cryptkit.PublicKeyStore,
	keyHolder cryptkit.SignatureKeyHolder,
	signature cryptkit.SignedDigestHolder,
) *StaticProfile {
	return &StaticProfile{
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

func (sp *StaticProfile) GetPrimaryRole() member.PrimaryRole {
	return sp.primaryRole
}

func (sp *StaticProfile) GetSpecialRoles() member.SpecialRole {
	return sp.specialRole
}

func (sp *StaticProfile) GetExtension() profiles.StaticProfileExtension {
	return sp.intro
}

func (sp *StaticProfile) GetDefaultEndpoint() endpoints.Outbound {
	return sp.endpoint
}

func (sp *StaticProfile) GetPublicKeyStore() cryptkit.PublicKeyStore {
	return sp.store
}

func (sp *StaticProfile) GetNodePublicKey() cryptkit.SignatureKeyHolder {
	return sp.keyHolder
}

func (sp *StaticProfile) GetStartPower() member.Power {
	// TODO: get from certificate
	return 10
}

func (sp *StaticProfile) IsAcceptableHost(from endpoints.Inbound) bool {
	address := sp.endpoint.GetNameAddress()
	return address.Equals(from.GetNameAddress())
}

func (sp *StaticProfile) GetStaticNodeID() insolar.ShortNodeID {
	return sp.shortID
}

func (sp *StaticProfile) GetBriefIntroSignedDigest() cryptkit.SignedDigestHolder {
	return sp.signature
}

func (sp *StaticProfile) String() string {
	return fmt.Sprintf("{sid:%d, node:%s}", sp.shortID, sp.intro.GetReference().String())
}

type Outbound struct {
	name endpoints.Name
	addr endpoints.IPAddress
}

func NewOutbound(address string) *Outbound {
	addr, err := endpoints.NewIPAddress(address)
	if err != nil {
		panic(err)
	}

	return &Outbound{
		name: endpoints.Name(address),
		addr: addr,
	}
}

func (p *Outbound) CanAccept(connection endpoints.Inbound) bool {
	return true
}

func (p *Outbound) GetEndpointType() endpoints.NodeEndpointType {
	return endpoints.IPEndpoint
}

func (*Outbound) GetRelayID() insolar.ShortNodeID {
	return 0
}

func (p *Outbound) GetNameAddress() endpoints.Name {
	return p.name
}

func (p *Outbound) GetIPAddress() endpoints.IPAddress {
	return p.addr
}

func (p *Outbound) AsByteString() longbits.ByteString {
	return longbits.ByteString(p.addr.String())
}

func NewStaticProfileList(nodes []insolar.NetworkNode, certificate insolar.Certificate, keyProcessor insolar.KeyProcessor) []profiles.StaticProfile {
	intros := make([]profiles.StaticProfile, len(nodes))
	for i, n := range nodes {
		intros[i] = NewStaticProfile(n, certificate, keyProcessor)
	}

	profiles.SortStaticProfiles(intros, false)

	return intros
}

func NewNetworkNode(profile profiles.ActiveNode) insolar.NetworkNode {
	nip := profile.GetStatic()
	store := nip.GetPublicKeyStore()
	introduction := nip.GetExtension()

	networkNode := node.NewNode(
		introduction.GetReference(),
		PrimaryRoleToStaticRole(nip.GetPrimaryRole()),
		store.(*ECDSAPublicKeyStore).publicKey,
		nip.GetDefaultEndpoint().GetNameAddress().String(),
		"",
	)

	mutableNode := networkNode.(node.MutableNode)

	mutableNode.SetShortID(profile.GetNodeID())
	mutableNode.SetState(insolar.NodeReady)

	mutableNode.SetPower(insolar.Power(profile.GetDeclaredPower()))
	if profile.GetOpMode().IsPowerless() {
		mutableNode.SetPower(0)
	}

	sd := nip.GetBriefIntroSignedDigest()
	mutableNode.SetSignature(
		sd.GetDigestHolder().AsBytes(),
		insolar.SignatureFromBytes(sd.GetSignatureHolder().AsBytes()),
	)

	return networkNode
}

func NewNetworkNodeList(profiles []profiles.ActiveNode) []insolar.NetworkNode {
	networkNodes := make([]insolar.NetworkNode, len(profiles))
	for i, p := range profiles {
		networkNodes[i] = NewNetworkNode(p)
	}

	return networkNodes
}

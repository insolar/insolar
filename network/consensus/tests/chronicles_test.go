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

package tests

import (
	"fmt"
	"math"

	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"

	"github.com/insolar/insolar/network/consensusv1/packets"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/gcpv2/census"
)

func NewEmuChronicles(intros []gcp_types.NodeIntroProfile, localNodeIndex int, primingCloudStateHash gcp_types.CloudStateHash) api.ConsensusChronicles {
	pop := census.NewManyNodePopulation(intros[localNodeIndex], intros)
	chronicles := census.NewLocalChronicles()
	census.NewPrimingCensus(
		&pop,
		nil,
		&EmuVersionedRegistries{primingCloudStateHash: primingCloudStateHash},
	).SetAsActiveTo(chronicles)
	return chronicles
}

func NewEmuNodeIntros(names ...string) []gcp_types.NodeIntroProfile {
	r := make([]gcp_types.NodeIntroProfile, len(names))
	for i, n := range names {
		var sr gcp_types.NodeSpecialRole
		var pr gcp_types.NodePrimaryRole
		switch n[0] {
		case 'h':
			pr = gcp_types.PrimaryRoleHeavyMaterial
			sr = gcp_types.SpecialRoleDiscovery
		case 'l':
			pr = gcp_types.PrimaryRoleLightMaterial
		case 'v':
			pr = gcp_types.PrimaryRoleVirtual
		default:
			pr = gcp_types.PrimaryRoleNeutral
			sr = gcp_types.SpecialRoleDiscovery
		}
		r[i] = NewEmuNodeIntro(i, endpoints.HostAddress(n), pr, sr)
	}
	return r
}

type EmuVersionedRegistries struct {
	pd                    pulse_data.PulseData
	primingCloudStateHash gcp_types.CloudStateHash
}

func (c *EmuVersionedRegistries) GetConsensusConfiguration() api.ConsensusConfiguration {
	return c
}

func (c *EmuVersionedRegistries) GetPrimingCloudHash() gcp_types.CloudStateHash {
	return c.primingCloudStateHash
}

func (c *EmuVersionedRegistries) FindRegisteredProfile(identity endpoints.HostIdentityHolder) gcp_types.HostProfile {
	return NewEmuNodeIntro(-1, identity.GetHostAddress(),
		/* unused by HostProfile */ gcp_types.NodePrimaryRole(math.MaxUint8), 0)
}

func (c *EmuVersionedRegistries) AddReport(report gcp_types.MisbehaviorReport) {
}

func (c *EmuVersionedRegistries) CommitNextPulse(pd pulse_data.PulseData, population api.OnlinePopulation) api.VersionedRegistries {
	pd.EnsurePulseData()
	cp := *c
	cp.pd = pd
	return &cp
}

func (c *EmuVersionedRegistries) GetMisbehaviorRegistry() api.MisbehaviorRegistry {
	return c
}

func (c *EmuVersionedRegistries) GetMandateRegistry() api.MandateRegistry {
	return c
}

func (c *EmuVersionedRegistries) GetOfflinePopulation() api.OfflinePopulation {
	return c
}

func (c *EmuVersionedRegistries) GetVersionPulseData() pulse_data.PulseData {
	return c.pd
}

const ShortNodeIdOffset = 1000

func NewEmuNodeIntro(id int, s endpoints.HostAddress, pr gcp_types.NodePrimaryRole, sr gcp_types.NodeSpecialRole) gcp_types.NodeIntroProfile {
	return &emuNodeIntro{
		id: insolar.ShortNodeID(ShortNodeIdOffset + id),
		n:  &emuEndpoint{name: s},
		pr: pr,
		sr: sr,
	}
}

var _ endpoints.NodeEndpoint = &emuEndpoint{}

type emuEndpoint struct {
	name endpoints.HostAddress
}

func (p *emuEndpoint) GetIPAddress() packets.NodeAddress {
	panic("implement me")
}

func (p *emuEndpoint) GetEndpointType() endpoints.NodeEndpointType {
	return endpoints.NameEndpoint
}

func (*emuEndpoint) GetRelayID() insolar.ShortNodeID {
	return 0
}

func (p *emuEndpoint) GetNameAddress() endpoints.HostAddress {
	return p.name
}

type emuNodeIntro struct {
	n  endpoints.NodeEndpoint
	id insolar.ShortNodeID
	pr gcp_types.NodePrimaryRole
	sr gcp_types.NodeSpecialRole
}

func (c *emuNodeIntro) GetNodePublicKey() cryptography_containers.SignatureKeyHolder {
	v := &long_bits.Bits512{}
	long_bits.FillBitsWithStaticNoise(uint32(c.id), v[:])
	k := cryptography_containers.NewSignatureKey(v, "stub/stub", cryptography_containers.PublicAsymmetricKey)
	return &k
}

func (c *emuNodeIntro) GetStartPower() gcp_types.MemberPower {
	return 10
}

func (c *emuNodeIntro) GetNodeReference() insolar.Reference {
	panic("unsupported")
}

func (c *emuNodeIntro) HasIntroduction() bool {
	return true
}

func (c *emuNodeIntro) ConvertPowerRequest(request gcp_types.PowerRequest) gcp_types.MemberPower {
	if ok, cl := request.AsCapacityLevel(); ok {
		return gcp_types.MemberPowerOf(uint16(cl.DefaultPercent()))
	}
	_, pw := request.AsMemberPower()
	return pw
}

func (c *emuNodeIntro) GetPrimaryRole() gcp_types.NodePrimaryRole {
	return c.pr
}

func (c *emuNodeIntro) GetSpecialRoles() gcp_types.NodeSpecialRole {
	return c.sr
}

func (*emuNodeIntro) IsAllowedPower(p gcp_types.MemberPower) bool {
	return true
}

func (c *emuNodeIntro) GetAnnouncementSignature() cryptography_containers.SignatureHolder {
	return nil
}

func (c *emuNodeIntro) GetDefaultEndpoint() endpoints.NodeEndpoint {
	return c.n
}

func (*emuNodeIntro) GetNodePublicKeyStore() cryptography_containers.PublicKeyStore {
	return nil
}

func (c *emuNodeIntro) IsAcceptableHost(from endpoints.HostIdentityHolder) bool {
	addr := c.n.GetNameAddress()
	return addr.Equals(from.GetHostAddress())
}

func (c *emuNodeIntro) GetShortNodeID() insolar.ShortNodeID {
	return c.id
}

func (c *emuNodeIntro) GetIntroduction() gcp_types.NodeIntroduction {
	return c
}

func (c *emuNodeIntro) String() string {
	return fmt.Sprintf("{sid:%v, n:%v}", c.id, c.n)
}

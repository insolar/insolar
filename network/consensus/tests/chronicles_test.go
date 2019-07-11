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

	"github.com/insolar/insolar/network/consensusv1/packets"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/errors"
)

func NewEmuChronicles(intros []common2.NodeIntroProfile, localNodeIndex int, primingCloudStateHash common2.CloudStateHash) census.ConsensusChronicles {
	pop := census.NewManyNodePopulation(intros[localNodeIndex], intros)
	chronicles := census.NewLocalChronicles()
	census.NewPrimingCensus(
		&pop,
		nil,
		&EmuVersionedRegistries{primingCloudStateHash: primingCloudStateHash},
	).SetAsActiveTo(chronicles)
	return chronicles
}

func NewEmuNodeIntros(names ...string) []common2.NodeIntroProfile {
	r := make([]common2.NodeIntroProfile, len(names))
	for i, n := range names {
		var sr common2.NodeSpecialRole
		var pr common2.NodePrimaryRole
		switch n[0] {
		case 'h':
			pr = common2.PrimaryRoleHeavyMaterial
			sr = common2.SpecialRoleDiscovery
		case 'l':
			pr = common2.PrimaryRoleLightMaterial
		case 'v':
			pr = common2.PrimaryRoleVirtual
		default:
			pr = common2.PrimaryRoleNeutral
			sr = common2.SpecialRoleDiscovery
		}
		r[i] = NewEmuNodeIntro(i, common.HostAddress(n), pr, sr)
	}
	return r
}

type EmuVersionedRegistries struct {
	pd                    common.PulseData
	primingCloudStateHash common2.CloudStateHash
}

func (c *EmuVersionedRegistries) GetConsensusConfiguration() census.ConsensusConfiguration {
	return c
}

func (c *EmuVersionedRegistries) GetPrimingCloudHash() common2.CloudStateHash {
	return c.primingCloudStateHash
}

func (c *EmuVersionedRegistries) FindRegisteredProfile(identity common.HostIdentityHolder) common2.HostProfile {
	return NewEmuNodeIntro(-1, identity.GetHostAddress(),
		/* unused by HostProfile */ common2.NodePrimaryRole(math.MaxUint8), 0)
}

func (c *EmuVersionedRegistries) AddReport(report errors.MisbehaviorReport) {
}

func (c *EmuVersionedRegistries) CommitNextPulse(pd common.PulseData, population census.OnlinePopulation) census.VersionedRegistries {
	pd.EnsurePulseData()
	cp := *c
	cp.pd = pd
	return &cp
}

func (c *EmuVersionedRegistries) GetMisbehaviorRegistry() census.MisbehaviorRegistry {
	return c
}

func (c *EmuVersionedRegistries) GetMandateRegistry() census.MandateRegistry {
	return c
}

func (c *EmuVersionedRegistries) GetOfflinePopulation() census.OfflinePopulation {
	return c
}

func (c *EmuVersionedRegistries) GetVersionPulseData() common.PulseData {
	return c.pd
}

const ShortNodeIdOffset = 1000

func NewEmuNodeIntro(id int, s common.HostAddress, pr common2.NodePrimaryRole, sr common2.NodeSpecialRole) common2.NodeIntroProfile {
	return &emuNodeIntro{
		id: common.ShortNodeID(ShortNodeIdOffset + id),
		n:  &emuEndpoint{name: s},
		pr: pr,
		sr: sr,
	}
}

var _ common.NodeEndpoint = &emuEndpoint{}

type emuEndpoint struct {
	name common.HostAddress
}

func (p *emuEndpoint) GetIPAddress() packets.NodeAddress {
	panic("implement me")
}

func (p *emuEndpoint) GetEndpointType() common.NodeEndpointType {
	return common.NameEndpoint
}

func (*emuEndpoint) GetRelayID() common.ShortNodeID {
	return 0
}

func (p *emuEndpoint) GetNameAddress() common.HostAddress {
	return p.name
}

type emuNodeIntro struct {
	n  common.NodeEndpoint
	id common.ShortNodeID
	pr common2.NodePrimaryRole
	sr common2.NodeSpecialRole
}

func (c *emuNodeIntro) GetNodePublicKey() common.SignatureKeyHolder {
	v := &common.Bits512{}
	common.FillBitsWithStaticNoise(uint32(c.id), v[:])
	k := common.NewSignatureKey(v, "stub/stub", common.PublicAsymmetricKey)
	return &k
}

func (c *emuNodeIntro) GetStartPower() common2.MemberPower {
	return 10
}

func (c *emuNodeIntro) GetNodeReference() insolar.Reference {
	panic("unsupported")
}

func (c *emuNodeIntro) HasIntroduction() bool {
	return true
}

func (c *emuNodeIntro) ConvertPowerRequest(request common2.PowerRequest) common2.MemberPower {
	if ok, cl := request.AsCapacityLevel(); ok {
		return common2.MemberPowerOf(uint16(cl.DefaultPercent()))
	}
	_, pw := request.AsMemberPower()
	return pw
}

func (c *emuNodeIntro) GetPrimaryRole() common2.NodePrimaryRole {
	return c.pr
}

func (c *emuNodeIntro) GetSpecialRoles() common2.NodeSpecialRole {
	return c.sr
}

func (*emuNodeIntro) IsAllowedPower(p common2.MemberPower) bool {
	return true
}

func (c *emuNodeIntro) GetAnnouncementSignature() common.SignatureHolder {
	return nil
}

func (c *emuNodeIntro) GetDefaultEndpoint() common.NodeEndpoint {
	return c.n
}

func (*emuNodeIntro) GetNodePublicKeyStore() common.PublicKeyStore {
	return nil
}

func (c *emuNodeIntro) IsAcceptableHost(from common.HostIdentityHolder) bool {
	addr := c.n.GetNameAddress()
	return addr.Equals(from.GetHostAddress())
}

func (c *emuNodeIntro) GetShortNodeID() common.ShortNodeID {
	return c.id
}

func (c *emuNodeIntro) GetIntroduction() common2.NodeIntroduction {
	return c
}

func (c *emuNodeIntro) String() string {
	return fmt.Sprintf("{sid:%v, n:%v}", c.id, c.n)
}

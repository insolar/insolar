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
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
	"github.com/insolar/insolar/pulse"
)

func NewEmuChronicles(intros []profiles.StaticProfile, localNodeIndex int, asJoiner bool,
	primingCloudStateHash proofs.CloudStateHash) api.ConsensusChronicles {

	var localCensus *censusimpl.PrimingCensusTemplate
	registries := &EmuVersionedRegistries{primingCloudStateHash: primingCloudStateHash}

	if asJoiner {
		if len(intros) != 1 && localNodeIndex != 0 {
			panic("illegal state")
		}
		localCensus = censusimpl.NewPrimingCensusForJoiner(intros[localNodeIndex], registries, EmuDefaultCryptography, true)
	} else {
		localCensus = censusimpl.NewPrimingCensus(intros, intros[localNodeIndex], registries, EmuDefaultCryptography, true)
	}

	chronicles := censusimpl.NewLocalChronicles(profiles.NewSimpleProfileIntroFactory(EmuDefaultCryptography))
	localCensus.SetAsActiveTo(chronicles)
	return chronicles
}

func NewEmuNodeIntros(names ...string) []profiles.StaticProfile {
	r := make([]profiles.StaticProfile, len(names))
	for i, n := range names {
		r[i] = NewEmuNodeIntroByName(i, n)
	}
	return r
}

func NewEmuNodeIntroByName(id int, name string) *EmuNodeIntro {

	var sr member.SpecialRole
	var pr member.PrimaryRole
	switch name[0] {
	case 'H':
		pr = member.PrimaryRoleHeavyMaterial
		sr = member.SpecialRoleDiscovery
	case 'L':
		pr = member.PrimaryRoleLightMaterial
	case 'V':
		pr = member.PrimaryRoleVirtual
	default:
		pr = member.PrimaryRoleNeutral
		sr = member.SpecialRoleDiscovery
	}
	return NewEmuNodeIntro(id, endpoints.Name(name), pr, sr)
}

type EmuVersionedRegistries struct {
	pd                    pulse.Data
	primingCloudStateHash proofs.CloudStateHash
}

func (c *EmuVersionedRegistries) GetNearestValidPulseData() pulse.Data {
	panic("implement me")
}

func (c *EmuVersionedRegistries) GetCloudIdentity() cryptkit.DigestHolder {
	return c.primingCloudStateHash
}

func (c *EmuVersionedRegistries) GetConsensusConfiguration() census.ConsensusConfiguration {
	return c
}

func (c *EmuVersionedRegistries) GetPrimingCloudHash() proofs.CloudStateHash {
	return c.primingCloudStateHash
}

func (c *EmuVersionedRegistries) FindRegisteredProfile(identity endpoints.Inbound) profiles.Host {
	return NewEmuNodeIntro(-1, identity.GetNameAddress(),
		/* unused by HostProfile */ member.PrimaryRole(math.MaxUint8), 0)
}

func (c *EmuVersionedRegistries) AddReport(report misbehavior.Report) {
}

func (c *EmuVersionedRegistries) CommitNextPulse(pd pulse.Data, population census.OnlinePopulation) census.VersionedRegistries {
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

func (c *EmuVersionedRegistries) GetVersionPulseData() pulse.Data {
	return c.pd
}

const ShortNodeIdOffset = 1000

func NewEmuNodeIntro(id int, s endpoints.Name, pr member.PrimaryRole, sr member.SpecialRole) *EmuNodeIntro {
	return &EmuNodeIntro{
		id: insolar.ShortNodeID(ShortNodeIdOffset + id),
		n:  &emuEndpoint{name: s},
		pr: pr,
		sr: sr,
	}
}

var _ endpoints.Outbound = &emuEndpoint{}

type emuEndpoint struct {
	name endpoints.Name
}

func (p *emuEndpoint) CanAccept(connection endpoints.Inbound) bool {
	return p.name == connection.GetNameAddress()
}

func (p *emuEndpoint) AsByteString() longbits.ByteString {
	return longbits.ByteString(fmt.Sprintf("out:name:%s", p.name))
}

func (p *emuEndpoint) GetIPAddress() endpoints.IPAddress {
	panic("implement me")
}

func (p *emuEndpoint) GetEndpointType() endpoints.NodeEndpointType {
	return endpoints.NameEndpoint
}

func (*emuEndpoint) GetRelayID() insolar.ShortNodeID {
	return 0
}

func (p *emuEndpoint) GetNameAddress() endpoints.Name {
	return p.name
}

var _ profiles.StaticProfile = &EmuNodeIntro{}
var baseIssuedAtTime = time.Now()

type EmuNodeIntro struct {
	n  endpoints.Outbound
	id insolar.ShortNodeID
	pr member.PrimaryRole
	sr member.SpecialRole
}

func (c *EmuNodeIntro) GetBriefIntroSignedDigest() cryptkit.SignedDigestHolder {
	dd := longbits.NewBits64(uint64(1000000 + c.id))
	ds := longbits.NewBits64(uint64(1000000+c.id) << 32)

	return cryptkit.NewSignedDigest(
		cryptkit.NewDigest(&dd, "stubHash"),
		cryptkit.NewSignature(&ds, "stubSign")).AsSignedDigestHolder()
}

func (c *EmuNodeIntro) GetIssuedAtPulse() pulse.Number {
	return 0
}

func (c *EmuNodeIntro) GetIssuedAtTime() time.Time {
	return baseIssuedAtTime
}

func (c *EmuNodeIntro) GetPowerLevels() member.PowerSet {
	return member.PowerSet{0, 0, 0, 0xFF}
}

func (c *EmuNodeIntro) GetExtraEndpoints() []endpoints.Outbound {
	return nil
}

func (c *EmuNodeIntro) GetIssuerID() insolar.ShortNodeID {
	return 0
}

func (c *EmuNodeIntro) GetIssuerSignature() cryptkit.SignatureHolder {
	ds := longbits.NewBits64(uint64(5000000+c.id) << 32)

	return cryptkit.NewSignature(&ds, "stubSign").AsSignatureHolder()
}

func (c *EmuNodeIntro) GetNodePublicKey() cryptkit.SignatureKeyHolder {
	v := &longbits.Bits512{}
	longbits.FillBitsWithStaticNoise(uint32(c.id), v[:])
	k := cryptkit.NewSignatureKey(v, "stub/stub", cryptkit.PublicAsymmetricKey)
	return &k
}

func (c *EmuNodeIntro) GetStartPower() member.Power {
	return 10
}

func (c *EmuNodeIntro) GetReference() insolar.Reference {
	return *insolar.NewEmptyReference()
}

func (c *EmuNodeIntro) ConvertPowerRequest(request power.Request) member.Power {
	if ok, cl := request.AsCapacityLevel(); ok {
		return member.PowerOf(uint16(cl.DefaultPercent()))
	}
	_, pw := request.AsMemberPower()
	return pw
}

func (c *EmuNodeIntro) GetPrimaryRole() member.PrimaryRole {
	return c.pr
}

func (c *EmuNodeIntro) GetSpecialRoles() member.SpecialRole {
	return c.sr
}

func (*EmuNodeIntro) IsAllowedPower(p member.Power) bool {
	return true
}

func (c *EmuNodeIntro) GetDefaultEndpoint() endpoints.Outbound {
	return c.n
}

func (*EmuNodeIntro) GetPublicKeyStore() cryptkit.PublicKeyStore {
	return nil
}

func (c *EmuNodeIntro) IsAcceptableHost(from endpoints.Inbound) bool {
	addr := c.n.GetNameAddress()
	return addr.Equals(from.GetNameAddress())
}

func (c *EmuNodeIntro) GetStaticNodeID() insolar.ShortNodeID {
	return c.id
}

func (c *EmuNodeIntro) GetIntroducedNodeID() insolar.ShortNodeID {
	return c.id
}

func (c *EmuNodeIntro) GetExtension() profiles.StaticProfileExtension {
	return c
}

func (c *EmuNodeIntro) String() string {
	return fmt.Sprintf("{sid:%v, n:%v}", c.id, c.n)
}

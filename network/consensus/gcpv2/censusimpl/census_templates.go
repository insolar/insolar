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

package censusimpl

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

var _ localActiveCensus = &PrimingCensusTemplate{}

type copyToOnlinePopulation interface {
	copyToPopulation
	census.OnlinePopulation
}

func NewPrimingCensusForJoiner(localProfile profiles.StaticProfile, registries census.VersionedRegistries,
	vf cryptkit.SignatureVerifierFactory, allowEphemeral bool) *PrimingCensusTemplate {

	pop := NewJoinerPopulation(localProfile, vf)
	return newPrimingCensus(&pop, registries, allowEphemeral)
}

func NewPrimingCensus(intros []profiles.StaticProfile, localProfile profiles.StaticProfile, registries census.VersionedRegistries,
	vf cryptkit.SignatureVerifierFactory, allowEphemeral bool) *PrimingCensusTemplate {

	if len(intros) == 0 {
		panic("illegal state")
	}
	localID := localProfile.GetStaticNodeID()
	pop := NewManyNodePopulation(intros, localID, vf)
	return newPrimingCensus(&pop, registries, allowEphemeral)
}

func newPrimingCensus(pop copyToOnlinePopulation, registries census.VersionedRegistries, allowEphemeral bool) *PrimingCensusTemplate {

	var pd pulse.Data
	if allowEphemeral {
		pd.PulseEpoch = pulse.EphemeralPulseEpoch
	} else {
		pd = registries.GetVersionPulseData()
	}

	return &PrimingCensusTemplate{CensusTemplate{nil,
		pop, &evictedPopulation{}, pd, registries,
	}}
}

var _ census.Prime = &PrimingCensusTemplate{}

type PrimingCensusTemplate struct {
	CensusTemplate
}

func (c *PrimingCensusTemplate) onMadeActive() {
}

func (c *PrimingCensusTemplate) IsActive() bool {
	return c.chronicles.GetActiveCensus() == c
}

func (c *PrimingCensusTemplate) SetAsActiveTo(chronicles LocalConsensusChronicles) {
	if c.chronicles != nil {
		panic("illegal state")
	}
	lc := chronicles.(*localChronicles)
	c.chronicles = lc
	lc.makeActive(nil, c)
}

func (*PrimingCensusTemplate) GetCensusState() census.State {
	return census.PrimingCensus
}

func (c *PrimingCensusTemplate) GetExpectedPulseNumber() pulse.Number {
	switch {
	case c.pd.IsEmpty():
		return pulse.Unknown
	case c.pd.IsExpectedPulse():
		return c.pd.GetPulseNumber()
	}
	return c.pd.GetNextPulseNumber()
}

func (c *PrimingCensusTemplate) CreateBuilder(ctx context.Context, pn pulse.Number) census.Builder {

	if !pn.IsUnknown() && (!pn.IsTimePulse() || !c.GetExpectedPulseNumber().IsUnknownOrEqualTo(pn)) {
		panic("illegal value")
	}

	return newLocalCensusBuilder(ctx, c.chronicles, pn, c.online)
}

func (c *PrimingCensusTemplate) BuildCopy(pd pulse.Data, csh proofs.CloudStateHash, gsh proofs.GlobulaStateHash) census.Built {
	if csh == nil {
		panic("illegal value: CSH is nil")
	}
	if gsh == nil {
		panic("illegal value: GSH is nil")
	}

	if c.pd.PulseEpoch != pulse.EphemeralPulseEpoch && pd.IsFromEphemeral() {
		panic("illegal value")
	}
	if !c.GetExpectedPulseNumber().IsUnknownOrEqualTo(pd.PulseNumber) {
		panic("illegal value")
	}

	return &BuiltCensusTemplate{ExpectedCensusTemplate{
		c.chronicles, c.online, c.evicted, c.chronicles.active,
		csh, gsh, pd.PulseNumber,
	}}
}

func (c *PrimingCensusTemplate) GetPulseNumber() pulse.Number {
	return c.pd.GetPulseNumber()
}

func (c *PrimingCensusTemplate) GetPulseData() pulse.Data {
	return c.pd
}

func (*PrimingCensusTemplate) GetGlobulaStateHash() proofs.GlobulaStateHash {
	return nil
}

func (*PrimingCensusTemplate) GetCloudStateHash() proofs.CloudStateHash {
	return nil
}

func (c PrimingCensusTemplate) String() string {
	return fmt.Sprintf("priming %s", c.CensusTemplate.String())
}

type CensusTemplate struct {
	chronicles *localChronicles
	online     copyToOnlinePopulation
	evicted    census.EvictedPopulation
	pd         pulse.Data

	registries census.VersionedRegistries
}

func (c *CensusTemplate) GetNearestPulseData() (bool, pulse.Data) {
	return true, c.pd
}

func (c *CensusTemplate) GetProfileFactory(ksf cryptkit.KeyStoreFactory) profiles.Factory {
	return c.chronicles.profileFactory
}

func (c *CensusTemplate) setVersionedRegistries(vr census.VersionedRegistries) {
	if vr == nil {
		panic("versioned registries are nil")
	}
	c.registries = vr
}

func (c *CensusTemplate) getVersionedRegistries() census.VersionedRegistries {
	return c.registries
}

func (c *CensusTemplate) GetOnlinePopulation() census.OnlinePopulation {
	return c.online
}

func (c *CensusTemplate) GetEvictedPopulation() census.EvictedPopulation {
	return c.evicted
}

func (c *CensusTemplate) GetOfflinePopulation() census.OfflinePopulation {
	return c.registries.GetOfflinePopulation()
}

func (c *CensusTemplate) GetMisbehaviorRegistry() census.MisbehaviorRegistry {
	return c.registries.GetMisbehaviorRegistry()
}

func (c *CensusTemplate) GetMandateRegistry() census.MandateRegistry {
	return c.registries.GetMandateRegistry()
}

func (c CensusTemplate) String() string {
	return fmt.Sprintf("pd:%v evicted:%v online:[%v]", c.pd, c.evicted, c.online)
}

var _ census.Active = &ActiveCensusTemplate{}

type ActiveCensusTemplate struct {
	CensusTemplate
	activeRef *ActiveCensusTemplate // hack for stringer
	gsh       proofs.GlobulaStateHash
	csh       proofs.CloudStateHash
}

func (c *ActiveCensusTemplate) onMadeActive() {
	c.activeRef = c
}

func (c *ActiveCensusTemplate) IsActive() bool {
	return c.chronicles.GetActiveCensus() == c
}

func (c *ActiveCensusTemplate) GetExpectedPulseNumber() pulse.Number {
	return c.pd.GetNextPulseNumber()
}

func (c *ActiveCensusTemplate) CreateBuilder(ctx context.Context, pn pulse.Number) census.Builder {

	if !pn.IsUnknown() && (!pn.IsTimePulse() || !c.GetExpectedPulseNumber().IsUnknownOrEqualTo(pn)) {
		panic("illegal value")
	}

	return newLocalCensusBuilder(ctx, c.chronicles, pn, c.online)
}

func (*ActiveCensusTemplate) GetCensusState() census.State {
	return census.SealedCensus
}

func (c *ActiveCensusTemplate) GetPulseNumber() pulse.Number {
	return c.pd.PulseNumber
}

func (c *ActiveCensusTemplate) GetPulseData() pulse.Data {
	return c.pd
}

func (c *ActiveCensusTemplate) GetGlobulaStateHash() proofs.GlobulaStateHash {
	return c.gsh
}

func (c *ActiveCensusTemplate) GetCloudStateHash() proofs.CloudStateHash {
	return c.csh
}

func (c ActiveCensusTemplate) String() string {
	mode := "active"
	if c.activeRef != c.chronicles.GetActiveCensus() {
		mode = "ex-active"
	}
	return fmt.Sprintf("%s %s gsh:%v csh:%v", mode, c.CensusTemplate.String(), c.gsh, c.csh)
}

var _ census.Expected = &ExpectedCensusTemplate{}

type ExpectedCensusTemplate struct {
	chronicles *localChronicles
	online     copyToOnlinePopulation
	evicted    census.EvictedPopulation
	prev       census.Active
	gsh        proofs.GlobulaStateHash
	csh        proofs.CloudStateHash
	pn         pulse.Number
}

func (c *ExpectedCensusTemplate) Rebuild(pn pulse.Number) census.Built {

	if !pn.IsUnknown() && (!pn.IsTimePulse() || !c.GetExpectedPulseNumber().IsUnknownOrEqualTo(pn)) {
		panic("illegal value")
	}

	cp := BuiltCensusTemplate{*c}
	cp.expected.pn = pn
	return &cp
}

func (c *ExpectedCensusTemplate) GetNearestPulseData() (bool, pulse.Data) {
	return false, c.prev.GetPulseData()
}

func (c *ExpectedCensusTemplate) GetProfileFactory(ksf cryptkit.KeyStoreFactory) profiles.Factory {
	return c.chronicles.GetProfileFactory(ksf)
}

func (c *ExpectedCensusTemplate) GetEvictedPopulation() census.EvictedPopulation {
	return c.evicted
}

func (c *ExpectedCensusTemplate) GetExpectedPulseNumber() pulse.Number {
	return c.pn
}

func (c *ExpectedCensusTemplate) GetCensusState() census.State {
	return census.CompleteCensus
}

func (c *ExpectedCensusTemplate) GetPulseNumber() pulse.Number {
	return c.pn
}

func (c *ExpectedCensusTemplate) GetGlobulaStateHash() proofs.GlobulaStateHash {
	return c.gsh
}

func (c *ExpectedCensusTemplate) GetCloudStateHash() proofs.CloudStateHash {
	return c.csh
}

func (c *ExpectedCensusTemplate) GetOnlinePopulation() census.OnlinePopulation {
	return c.online
}

func (c *ExpectedCensusTemplate) GetOfflinePopulation() census.OfflinePopulation {
	// TODO Should be provided via relevant builder
	return c.prev.GetOfflinePopulation()
}

func (c *ExpectedCensusTemplate) GetMisbehaviorRegistry() census.MisbehaviorRegistry {
	// TODO Should be provided via relevant builder
	return c.prev.GetMisbehaviorRegistry()
}

func (c *ExpectedCensusTemplate) GetMandateRegistry() census.MandateRegistry {
	// TODO Should be provided via relevant builder
	return c.prev.GetMandateRegistry()
}

func (c *ExpectedCensusTemplate) CreateBuilder(ctx context.Context, pn pulse.Number) census.Builder {
	return newLocalCensusBuilder(ctx, c.chronicles, pn, c.online)
}

func (c *ExpectedCensusTemplate) GetPrevious() census.Active {
	return c.prev
}

func (c *ExpectedCensusTemplate) MakeActive(pd pulse.Data) census.Active {

	a := ActiveCensusTemplate{
		CensusTemplate{
			c.chronicles, c.online, c.evicted, pd, nil,
		}, nil,
		c.gsh,
		c.csh,
	}

	c.chronicles.makeActive(c, &a)
	return &a
}

func (c *ExpectedCensusTemplate) IsActive() bool {
	return false
}

func (c ExpectedCensusTemplate) String() string {
	return fmt.Sprintf("expected pn:%v evicted:%v online:[%v] gsh:%v csh:%v", c.pn, c.evicted, c.online, c.gsh, c.csh)
}

var _ census.Built = &BuiltCensusTemplate{}

type BuiltCensusTemplate struct {
	expected ExpectedCensusTemplate
}

func (b *BuiltCensusTemplate) Update(csh proofs.CloudStateHash, gsh proofs.GlobulaStateHash) census.Built {
	if b.expected.gsh != nil && gsh == nil {
		panic("illegal value")
	}
	if b.expected.csh != nil && csh == nil {
		panic("illegal value")
	}
	cp := *b
	cp.expected.gsh = gsh
	cp.expected.csh = csh
	return &cp
}

func (b *BuiltCensusTemplate) GetOnlinePopulation() census.OnlinePopulation {
	return b.expected.GetOnlinePopulation()
}

func (b *BuiltCensusTemplate) GetEvictedPopulation() census.EvictedPopulation {
	return b.expected.GetEvictedPopulation()
}

func (b *BuiltCensusTemplate) GetGlobulaStateHash() proofs.GlobulaStateHash {
	return b.expected.GetGlobulaStateHash()
}

func (b *BuiltCensusTemplate) GetCloudStateHash() proofs.CloudStateHash {
	return b.expected.GetCloudStateHash()
}

func (b *BuiltCensusTemplate) GetNearestPulseData() (bool, pulse.Data) {
	return b.expected.GetNearestPulseData()
}

func (b *BuiltCensusTemplate) MakeExpected() census.Expected {
	return b.expected.chronicles.makeExpected(&b.expected)
}

func (b BuiltCensusTemplate) String() string {
	return fmt.Sprintf("built-%s", b.expected.String())
}

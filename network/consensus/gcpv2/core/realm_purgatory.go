///
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
///

package core

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
	"sync"
)

func NewRealmPurgatory(population RealmPopulation, pf profiles.Factory, svf cryptkit.SignatureVerifierFactory,
	callback *nodeContext, postponedPacketFn PostponedPacketFunc) RealmPurgatory {
	return RealmPurgatory{
		population:        population,
		profileFactory:    pf,
		svFactory:         svf,
		callback:          callback,
		postponedPacketFn: postponedPacketFn,
	}
}

type RealmPurgatory struct {
	population        RealmPopulation
	svFactory         cryptkit.SignatureVerifierFactory
	profileFactory    profiles.Factory
	postponedPacketFn PostponedPacketFunc

	callback *nodeContext

	/* LOCK WARNING!
	This lock is engaged inside NodePhantom's lock.
	DO NOT call NodePhantom methods under this lock.
	*/
	rw sync.RWMutex

	phantomByID map[insolar.ShortNodeID]*NodePhantom

	//phantomByEP map[string]*NodePhantom
}

//type PurgatoryNodeState int
//
//const PurgatoryDuplicatePK PurgatoryNodeState = -1
//const PurgatoryExistingMember PurgatoryNodeState = -2

func (p *RealmPurgatory) GetPhantomNode(id insolar.ShortNodeID) *NodePhantom {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return p.phantomByID[id]
}

func (p *RealmPurgatory) getPhantomNode(id insolar.ShortNodeID) (*NodePhantom, bool) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	np, ok := p.phantomByID[id]
	return np, ok
}

type applyNodeIntroFunc func(ctx context.Context, brief profiles.BriefCandidateProfile, full profiles.CandidateProfile,
	announcerID, originID insolar.ShortNodeID) error

func (p *RealmPurgatory) _getApplyFn(id insolar.ShortNodeID, np *NodePhantom, isNotice bool) applyNodeIntroFunc {

	if np != nil {
		return np.ApplyNodeIntro
	}
	//NB! np == NIL - it means that phantom was moved to a normal population
	if isNotice {
		//we dont need to send notice to a NodeAppearance
		return func(_ context.Context, _ profiles.BriefCandidateProfile, _ profiles.CandidateProfile, _, _ insolar.ShortNodeID) error {
			return nil
		}
	}
	na := p.population.GetNodeAppearance(id)
	if na == nil {
		//nil entry in the purgatory MUST have a relevant NodeAppearance
		panic("illegal state")
	}
	return na.ApplyNodeIntro
}

func (p *RealmPurgatory) getApplyFn(id insolar.ShortNodeID, isNotice bool) applyNodeIntroFunc {

	np, ok := p.getPhantomNode(id)
	if ok {
		return p._getApplyFn(id, np, isNotice)
	}

	p.rw.Lock()
	defer p.rw.Unlock()
	if p.phantomByID == nil {
		p.phantomByID = make(map[insolar.ShortNodeID]*NodePhantom)
		//p.phantomByEP = make(map[string]*NodePhantom)
	} else {
		np, ok = p.phantomByID[id]
		if ok {
			return p._getApplyFn(id, np, isNotice)
		}
	}
	limiter := p.population.CreatePacketLimiter()
	np = NewNodePhantom(p, id, limiter)
	p.phantomByID[id] = np
	return np.ApplyNodeIntro
}

func (p *RealmPurgatory) ascendFromPurgatory(ctx context.Context, id insolar.ShortNodeID, nsp profiles.StaticProfile,
	sv cryptkit.SignatureVerifier, packets []PostponedPacket) {

	if sv == nil {
		sv = p.svFactory.GetSignatureVerifierWithPKS(nsp.GetPublicKeyStore())
	}
	np := censusimpl.NewJoinerProfile(nsp, sv, nsp.GetStartPower())
	na := p.population.CreateNodeAppearance(ctx, &np)

	p.rw.Lock()
	defer p.rw.Unlock()
	p.phantomByID[id] = nil //leave marker
	//delete(p.phantomByEP, ...)
	_, _ = p.population.AddToDynamics(na)
	go p.flushPostponedPackets(packets)
}

func (p *RealmPurgatory) NoticeFromNeighbourhood(ctx context.Context,
	joinerID insolar.ShortNodeID, announcerID, originID insolar.ShortNodeID) error {

	return p.getApplyFn(joinerID, true)(ctx, nil, nil, announcerID, originID)
}

func (p *RealmPurgatory) BriefSelfFromNeighbourhood(ctx context.Context,
	joinerID insolar.ShortNodeID, brief transport.BriefIntroductionReader, originID insolar.ShortNodeID) error {

	if brief == nil {
		panic("illegal value")
	}
	return p.getApplyFn(joinerID, false)(ctx, brief, nil, joinerID, insolar.AbsentShortNodeID)
}

func (p *RealmPurgatory) FromSelfIntroduction(ctx context.Context,
	joinerID insolar.ShortNodeID, brief transport.BriefIntroductionReader, full transport.FullIntroductionReader) error {

	if brief == nil && full == nil {
		panic("illegal value")
	}
	if brief != nil && full != nil {
		if !profiles.EqualStaticProfiles(full, brief) {
			return fmt.Errorf("deserialization error")
		}
		brief = nil
	}

	return p.getApplyFn(joinerID, false)(ctx, brief, full, joinerID, insolar.AbsentShortNodeID)
}

func (p *RealmPurgatory) FromMemberAnnouncement(ctx context.Context, joinerID insolar.ShortNodeID,
	brief transport.BriefIntroductionReader, full transport.FullIntroductionReader, announcerID insolar.ShortNodeID) error {

	if brief == nil && full == nil {
		panic("illegal value")
	}
	if brief != nil && full != nil {
		if !profiles.EqualStaticProfiles(full, brief) {
			return fmt.Errorf("deserialization error")
		}
		brief = nil
	}

	return p.getApplyFn(joinerID, false)(ctx, brief, full, announcerID, announcerID)
}

func (p *RealmPurgatory) sendPostponedPacket(_ context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags PacketVerifyFlags) {

	p.postponedPacketFn(packet, from, flags)
}

func (p *RealmPurgatory) GetProfileFactory() profiles.Factory {
	return p.profileFactory
}

func (p *RealmPurgatory) IsBriefAscensionAllowed() bool {
	return true // TODO using false will fail vector calculation, because only NodeAppearance can be there now
}

func (p *RealmPurgatory) flushPostponedPackets(packets []PostponedPacket) {
	for _, pp := range packets {
		p.postponedPacketFn(pp.Packet, pp.From, pp.VerifyFlags)
	}
}

func (p *RealmPurgatory) onBriefProfileCreated(n *NodePhantom) {
	p.callback.onPurgatoryNodeAdded(p.callback.GetPopulationVersion(), n)
}

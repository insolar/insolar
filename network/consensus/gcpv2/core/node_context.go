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
	"sync/atomic"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
)

func (p *nodeContext) initPrep(localNodeID insolar.ShortNodeID, signatureVerifierFactory cryptkit.SignatureVerifierFactory, capture misbehavior.ReportFunc) {
	p.localNodeID = localNodeID
	p.signatureVerifierFactory = signatureVerifierFactory
	p.fraudFactory = misbehavior.NewFraudFactory(capture)
	p.blameFactory = misbehavior.NewBlameFactory(capture)
}

func (p *nodeContext) initFull(localNodeID insolar.ShortNodeID, signatureVerifierFactory cryptkit.SignatureVerifierFactory,
	neighborTrustThreshold uint8, capture misbehavior.ReportFunc) {

	p.localNodeID = localNodeID
	p.signatureVerifierFactory = signatureVerifierFactory
	p.fraudFactory = misbehavior.NewFraudFactory(capture)
	p.blameFactory = misbehavior.NewBlameFactory(capture)
	p.nbTrustThreshold = neighborTrustThreshold
}

func (p *nodeContext) setNodeToPhaseCallback(phaseControllerCallback NodeUpdateCallback) {
	p.phaseControllerCallback = phaseControllerCallback
}

type NodeContextHolder *nodeContext

type nodeContext struct {
	fraudFactory            misbehavior.FraudFactory
	blameFactory            misbehavior.BlameFactory
	phaseControllerCallback NodeUpdateCallback

	populationVersion uint32 // atomic

	signatureVerifierFactory cryptkit.SignatureVerifierFactory
	nbTrustThreshold         uint8
	localNodeID              insolar.ShortNodeID
}

func (p *nodeContext) updatePopulationVersion() uint32 {
	return atomic.AddUint32(&p.populationVersion, 1)
}

func (p *nodeContext) GetPopulationVersion() uint32 {
	return atomic.LoadUint32(&p.populationVersion)
}

func (p *nodeContext) GetNeighbourhoodTrustThreshold() uint8 {
	if p.nbTrustThreshold == 0 {
		panic("illegal state: not allowed for PrepRealm")
	}
	return p.nbTrustThreshold
}

func (p *nodeContext) GetFraudFactory() misbehavior.FraudFactory {
	return p.fraudFactory
}

func (p *nodeContext) GetBlameFactory() misbehavior.BlameFactory {
	return p.blameFactory
}

func (p *nodeContext) GetSignatureVerifierFactory() cryptkit.SignatureVerifierFactory {
	return p.signatureVerifierFactory
}

func (p *nodeContext) onTrustUpdated(populationVersion uint32, n *NodeAppearance, before member.TrustLevel, after member.TrustLevel) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnTrustUpdated(populationVersion, n, before, after)
}

func (p *nodeContext) onNodeStateAssigned(populationVersion uint32, n *NodeAppearance) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnNodeStateAssigned(populationVersion, n)
}

func (p *nodeContext) onDynamicNodeAdded(populationVersion uint32, n *NodeAppearance, fullIntro bool) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnDynamicNodeAdded(populationVersion, n, fullIntro)
}

func (p *nodeContext) onPurgatoryNodeAdded(populationVersion uint32, n *NodePhantom) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnPurgatoryNodeAdded(populationVersion, n)
}

func (p *nodeContext) onCustomEvent(populationVersion uint32, n *NodeAppearance, event interface{}) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnCustomEvent(populationVersion, n, event)
}

func (p *nodeContext) onDynamicPopulationCompleted(populationVersion uint32, indexedCount int) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnDynamicPopulationCompleted(populationVersion, indexedCount)
}

func (p *nodeContext) GetLocalNodeID() insolar.ShortNodeID {
	return p.localNodeID
}

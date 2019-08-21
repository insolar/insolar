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

package population

import (
	"sync/atomic"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/pulse"
)

func NewHook(localNode profiles.ActiveNode, eventDispatcher EventDispatcher, hookCfg SharedNodeContext) Hook {

	if eventDispatcher == nil {
		eventDispatcher = NewPanicDispatcher("illegal state")
	}

	return Hook{
		internalPopulationEventDispatcher: eventDispatcher,
		local:                             localNode,
		config:                            hookConfig{hookCfg},
	}
}

var _ EventDispatcher = &Hook{}

type SharedNodeContext struct {
	FraudFactory     misbehavior.FraudFactory
	BlameFactory     misbehavior.BlameFactory
	Assistant        transport.CryptographyAssistant
	PulseData        pulse.DataHolder
	NbTrustThreshold uint8
	EphemeralMode    api.EphemeralMode
}

func NewSharedNodeContext(assistant transport.CryptographyAssistant, pdh pulse.DataHolder, nbTrustThreshold uint8,
	ephemeralMode api.EphemeralMode, capture misbehavior.ReportFunc) SharedNodeContext {
	return SharedNodeContext{
		misbehavior.NewFraudFactory(capture),
		misbehavior.NewBlameFactory(capture),
		assistant,
		pdh,
		nbTrustThreshold,
		ephemeralMode,
	}
}

func NewSharedNodeContextByPulseNumber(assistant transport.CryptographyAssistant, pn pulse.Number, nbTrustThreshold uint8,
	ephemeralMode api.EphemeralMode, capture misbehavior.ReportFunc) SharedNodeContext {
	return SharedNodeContext{
		misbehavior.NewFraudFactory(capture),
		misbehavior.NewBlameFactory(capture),
		assistant,
		pulseDataHolder{pn},
		nbTrustThreshold,
		ephemeralMode,
	}
}

var _ pulse.DataHolder = &pulseDataHolder{}

type pulseDataHolder struct {
	pn pulse.Number
}

func (p pulseDataHolder) GetPulseNumber() pulse.Number {
	return p.pn
}

func (p pulseDataHolder) GetPulseData() pulse.Data {
	panic("illegal state")
}

func (p pulseDataHolder) GetPulseDataDigest() cryptkit.DigestHolder {
	return nil
}

type hookConfig struct {
	SharedNodeContext
}

type Hook struct {
	internalPopulationEventDispatcher
	config            hookConfig
	populationVersion uint32 // atomic
	local             profiles.ActiveNode
}

func (p *Hook) GetPulseData() pulse.DataHolder {
	return p.config.PulseData
}

func (p *Hook) UpdatePopulationVersion() uint32 {
	return atomic.AddUint32(&p.populationVersion, 1)
}

func (p *Hook) GetPopulationVersion() uint32 {
	return atomic.LoadUint32(&p.populationVersion)
}

func (p *Hook) GetNeighbourhoodTrustThreshold() uint8 {
	if p.config.NbTrustThreshold == 0 {
		panic("illegal state: not allowed for PrepRealm")
	}
	return p.config.NbTrustThreshold
}

func (p *Hook) GetFraudFactory() misbehavior.FraudFactory {
	return p.config.FraudFactory
}

func (p *Hook) GetBlameFactory() misbehavior.BlameFactory {
	return p.config.BlameFactory
}

func (p *Hook) GetCryptographyAssistant() transport.CryptographyAssistant {
	return p.config.Assistant
}

func (p *Hook) GetLocalNodeID() insolar.ShortNodeID {
	return p.local.GetNodeID()
}

func (p *Hook) GetLocalProfile() profiles.ActiveNode {
	return p.local
}

func (p *Hook) GetEphemeralMode() api.EphemeralMode {
	return p.config.EphemeralMode
}

type EventClosureFunc func(EventDispatcher)
type EventDispatchFunc func(EventClosureFunc)

type EventDispatcher interface {
	internalPopulationEventDispatcher
}

type MemberPacketSender interface {
	transport.TargetProfile
	SetPacketSent(pt phases.PacketType) bool
}

type UpdateFlags uint32

const (
	FlagCreated   UpdateFlags = 1 << iota
	FlagFixedInit             // for indexed members of a fixed population
	FlagUpdatedProfile
	FlagAscent // for purgatory nodes
)

type internalPopulationEventDispatcher interface {
	OnTrustUpdated(populationVersion uint32, n *NodeAppearance, before member.TrustLevel, after member.TrustLevel, fullProfile bool)
	OnNodeStateAssigned(populationVersion uint32, n *NodeAppearance)
	OnDynamicNodeUpdate(populationVersion uint32, n *NodeAppearance, flags UpdateFlags)
	OnPurgatoryNodeUpdate(populationVersion uint32, n MemberPacketSender, flags UpdateFlags)
	OnCustomEvent(populationVersion uint32, n *NodeAppearance, event interface{})
	OnDynamicPopulationCompleted(populationVersion uint32, indexedCount int)
}

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
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

// hides embedded pointer from external access
type hLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

type coreRealm struct {
	/* Provided externally at construction. Don't need mutex */
	hLocker

	roundContext  context.Context
	strategy      RoundStrategy
	config        api.LocalNodeConfiguration
	initialCensus census.Operational

	/* Derived from the ones provided externally - set at init() or start(). Don't need mutex */
	signer          cryptkit.DigestSigner
	digest          transport.ConsensusDigestFactory
	verifierFactory transport.CryptographyFactory
	upstream        api.UpstreamController
	roundStartedAt  time.Time

	expectedPopulationSize uint16
	nbhSizes               transport.NeighbourhoodSizes

	self *NodeAppearance /* Special case - this field is set twice, by start() of PrepRealm and FullRealm */

	/*
		Other fields - need mutex during PrepRealm, unless accessed by start() of PrepRealm
		FullRealm doesnt need a lock to read them
	*/
	pulseData     pulse.Data
	originalPulse proofs.OriginalPulsarPacket
}

func (r *coreRealm) init(hLocker hLocker, strategy RoundStrategy, transport transport.Factory,
	config api.LocalNodeConfiguration, initialCensus census.Operational) {

	r.hLocker = hLocker

	r.strategy = strategy
	r.config = config
	r.initialCensus = initialCensus

	r.verifierFactory = transport.GetCryptographyFactory()
	r.digest = r.verifierFactory.GetDigestFactory()

	sks := config.GetSecretKeyStore()
	r.signer = r.verifierFactory.GetNodeSigner(sks)
}

func (r *coreRealm) initPopulation(powerRequest power.Request, nbhSizes transport.NeighbourhoodSizes) {

	r.nbhSizes = nbhSizes
	population := r.initialCensus.GetOnlinePopulation()

	/*
		Here we initialize self like for PrepRealm. This is not perfect, but it enables access to Self earlier.
	*/
	nodeContext := &nodeContext{}
	profile := population.GetLocalProfile()

	if profile.GetOpMode().IsEvicted() {
		/*
			Previous round has provided an incorrect population, as eviction of local node must force one-node population of self
		*/
		panic("illegal state")
	}

	r.self = NewNodeAppearanceAsSelf(profile, nodeContext)
	r.self.requestedPower = profile.GetStatic().GetIntroduction().ConvertPowerRequest(powerRequest)

	nodeContext.initPrep(r.verifierFactory,
		func(report misbehavior.Report) interface{} {
			r.initialCensus.GetMisbehaviorRegistry().AddReport(report)
			return nil
		})

	r.expectedPopulationSize = uint16(population.GetCount())
}

func (r *coreRealm) GetStrategy() RoundStrategy {
	return r.strategy
}

func (r *coreRealm) GetVerifierFactory() cryptkit.SignatureVerifierFactory {
	return r.verifierFactory
}

func (r *coreRealm) GetDigestFactory() transport.ConsensusDigestFactory {
	return r.digest
}

func (r *coreRealm) GetSigner() cryptkit.DigestSigner {
	return r.signer
}

func (r *coreRealm) GetSignatureVerifier(pks cryptkit.PublicKeyStore) cryptkit.SignatureVerifier {
	return r.verifierFactory.GetSignatureVerifierWithPKS(pks)
}

func (r *coreRealm) GetStartedAt() time.Time {
	return r.roundStartedAt
}

func (r *coreRealm) AdjustedAfter(d time.Duration) time.Duration {
	return time.Until(r.roundStartedAt.Add(d))
}

func (r *coreRealm) GetRoundContext() context.Context {
	return r.roundContext
}

func (r *coreRealm) GetLocalConfig() api.LocalNodeConfiguration {
	return r.config
}

func (r *coreRealm) IsJoiner() bool {
	return r.self.IsJoiner()
}

func (r *coreRealm) GetSelfNodeID() insolar.ShortNodeID {
	return r.self.profile.GetShortNodeID()
}

func (r *coreRealm) GetSelf() *NodeAppearance {
	return r.self
}

func (r *coreRealm) GetPrimingCloudHash() proofs.CloudStateHash {
	return r.initialCensus.GetMandateRegistry().GetPrimingCloudHash()
}

func (r *coreRealm) VerifyPacketAuthenticity(packet transport.PacketParser, from endpoints.Inbound, strictFrom bool) error {
	nr := r.initialCensus.GetOfflinePopulation().FindRegisteredProfile(from)
	if nr == nil {
		nr = r.initialCensus.GetMandateRegistry().FindRegisteredProfile(from)
		if nr == nil {
			return fmt.Errorf("unable to identify sender: %v", from)
		}
	}
	sf := r.verifierFactory.GetSignatureVerifierWithPKS(nr.GetPublicKeyStore())
	return VerifyPacketAuthenticityBy(packet, nr, sf, from, strictFrom)
}

func VerifyPacketRoute(ctx context.Context, packet transport.PacketParser, selfID insolar.ShortNodeID) (bool, error) {

	sid := packet.GetSourceID()
	if sid == selfID {
		return false, fmt.Errorf("loopback, SourceID(%v) == thisNodeID(%v)", sid, selfID)
	}

	rid := packet.GetReceiverID()
	if rid != selfID {
		return false, fmt.Errorf("receiverID(%v) != thisNodeID(%v)", rid, selfID)
	}

	tid := packet.GetTargetID()
	if tid != selfID {
		// Relaying
		if packet.IsRelayForbidden() {
			return false, fmt.Errorf("sender doesn't allow relaying for targetID(%v)", tid)
		}

		// TODO relay support
		err := fmt.Errorf("unsupported: relay is required for targetID(%v)", tid)
		inslogger.FromContext(ctx).Errorf(err.Error())
		// allow sender to be different from source
		return false, err
	}

	// sender must be source
	return packet.IsRelayForbidden(), nil
}

func VerifyPacketAuthenticityBy(packet transport.PacketParser, nr profiles.Host, sf cryptkit.SignatureVerifier,
	from endpoints.Inbound, strictFrom bool) error {

	if strictFrom && !nr.IsAcceptableHost(from) {
		return fmt.Errorf("host is not allowed by node registration: node=%v, host=%v", nr, from)
	}

	ps := packet.GetPacketSignature()
	if !ps.IsVerifiableBy(sf) {
		return fmt.Errorf("unable to verify packet signature from sender: %v", from)
	}
	if !ps.VerifyWith(sf) {
		return fmt.Errorf("packet signature doesn't match for sender: %v", from)
	}
	return nil
}

func LazyPacketParse(packet transport.PacketParser) (transport.PacketParser, error) {

	//this enables lazy parsing - packet is fully parsed AFTER validation, hence makes it less prone to exploits for non-members
	newPacket, err := packet.ParsePacketBody()
	if err != nil {
		return packet, err
	}
	if newPacket == nil {
		return packet, nil
	}
	return newPacket, nil
}

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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"sync"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
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
	initialCensus api.OperationalCensus

	/* Derived from the ones provided externally - set at init() or start(). Don't need mutex */
	signer          cryptography_containers.DigestSigner
	digest          cryptography_containers.DigestFactory
	verifierFactory api.TransportCryptographyFactory
	upstream        api.UpstreamPulseController
	roundStartedAt  time.Time

	self *NodeAppearance /* Special case - this field is set twice, by start() of PrepRealm and FullRealm */

	/*
		Other fields - need mutex during PrepRealm, unless accessed by start() of PrepRealm
		FullRealm doesnt need a lock to read them
	*/
	pulseData     pulse_data.PulseData
	originalPulse packets.OriginalPulsarPacket
}

func (r *coreRealm) init(hLocker hLocker, strategy RoundStrategy, transport api.TransportFactory,
	config api.LocalNodeConfiguration, initialCensus api.OperationalCensus, powerRequest gcp_types.PowerRequest) {

	r.hLocker = hLocker

	r.strategy = strategy
	r.config = config
	r.initialCensus = initialCensus

	r.verifierFactory = transport.GetCryptographyFactory()
	r.digest = r.verifierFactory.GetDigestFactory()

	sks := config.GetSecretKeyStore()
	r.signer = r.verifierFactory.GetNodeSigner(sks)

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
	r.self.requestedPower = profile.GetIntroduction().ConvertPowerRequest(powerRequest)

	nodeContext.initPrep(
		func(report gcp_types.MisbehaviorReport) interface{} {
			r.initialCensus.GetMisbehaviorRegistry().AddReport(report)
			return nil
		})
}

func (r *coreRealm) GetStrategy() RoundStrategy {
	return r.strategy
}

func (r *coreRealm) GetVerifierFactory() cryptography_containers.SignatureVerifierFactory {
	return r.verifierFactory
}

func (r *coreRealm) GetDigestFactory() cryptography_containers.DigestFactory {
	return r.digest
}

func (r *coreRealm) GetSigner() cryptography_containers.DigestSigner {
	return r.signer
}

func (r *coreRealm) GetSignatureVerifier(pks cryptography_containers.PublicKeyStore) cryptography_containers.SignatureVerifier {
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

func (r *coreRealm) GetPrimingCloudHash() gcp_types.CloudStateHash {
	return r.initialCensus.GetMandateRegistry().GetPrimingCloudHash()
}

func (r *coreRealm) VerifyPacketAuthenticity(packet packets.PacketParser, from endpoints.HostIdentityHolder, strictFrom bool) error {
	nr := r.initialCensus.GetOfflinePopulation().FindRegisteredProfile(from)
	if nr == nil {
		nr = r.initialCensus.GetMandateRegistry().FindRegisteredProfile(from)
		if nr == nil {
			return fmt.Errorf("unable to identify sender: %v", from)
		}
	}
	sf := r.verifierFactory.GetSignatureVerifierWithPKS(nr.GetNodePublicKeyStore())
	return VerifyPacketAuthenticityBy(packet, nr, sf, from, strictFrom)
}

func VerifyPacketAuthenticityBy(packet packets.PacketParser, nr gcp_types.HostProfile, sf cryptography_containers.SignatureVerifier,
	from endpoints.HostIdentityHolder, strictFrom bool) error {

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

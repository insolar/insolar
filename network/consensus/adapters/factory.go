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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"github.com/insolar/insolar/network/consensus/gcpv2/phases"
)

type ECDSASignatureVerifierFactory struct {
	digester *Sha3512Digester
	scheme   insolar.PlatformCryptographyScheme
}

func NewECDSASignatureVerifierFactory(
	digester *Sha3512Digester,
	scheme insolar.PlatformCryptographyScheme,
) *ECDSASignatureVerifierFactory {
	return &ECDSASignatureVerifierFactory{
		digester: digester,
		scheme:   scheme,
	}
}

func (vf *ECDSASignatureVerifierFactory) GetSignatureVerifierWithPKS(pks cryptography_containers.PublicKeyStore) cryptography_containers.SignatureVerifier {
	keyStore := pks.(*ECDSAPublicKeyStore)

	return NewECDSASignatureVerifier(
		vf.digester,
		vf.scheme,
		keyStore.publicKey,
	)
}

type DigestFactory struct {
	pcs insolar.PlatformCryptographyScheme
}

func NewDigestFactory(pcs insolar.PlatformCryptographyScheme) *DigestFactory {
	return &DigestFactory{
		pcs: pcs,
	}
}

func (df *DigestFactory) GetPacketDigester() cryptography_containers.DataDigester {
	return NewSha3512Digester(df.pcs)
}

func (df *DigestFactory) GetGshDigester() cryptography_containers.SequenceDigester {
	return &gshDigester{}
}

type TransportCryptographyFactory struct {
	verifierFactory *ECDSASignatureVerifierFactory
	digestFactory   *DigestFactory
	scheme          insolar.PlatformCryptographyScheme
}

func NewTransportCryptographyFactory(scheme insolar.PlatformCryptographyScheme) *TransportCryptographyFactory {
	return &TransportCryptographyFactory{
		verifierFactory: NewECDSASignatureVerifierFactory(
			NewSha3512Digester(scheme),
			scheme,
		),
		digestFactory: NewDigestFactory(scheme),
		scheme:        scheme,
	}
}

func (cf *TransportCryptographyFactory) GetSignatureVerifierWithPKS(pks cryptography_containers.PublicKeyStore) cryptography_containers.SignatureVerifier {
	return cf.verifierFactory.GetSignatureVerifierWithPKS(pks)
}

func (cf *TransportCryptographyFactory) GetDigestFactory() cryptography_containers.DigestFactory {
	return cf.digestFactory
}

func (cf *TransportCryptographyFactory) GetNodeSigner(sks cryptography_containers.SecretKeyStore) cryptography_containers.DigestSigner {
	isks := sks.(*ECDSASecretKeyStore)

	return NewECDSADigestSigner(isks.privateKey, cf.scheme)
}

func (cf *TransportCryptographyFactory) GetPublicKeyStore(skh cryptography_containers.SignatureKeyHolder) cryptography_containers.PublicKeyStore {
	kh := skh.(*ECDSASignatureKeyHolder)

	return NewECDSAPublicKeyStore(kh.publicKey)
}

type RoundStrategyFactory struct{}

func NewRoundStrategyFactory() *RoundStrategyFactory {
	return &RoundStrategyFactory{}
}

func (rsf *RoundStrategyFactory) CreateRoundStrategy(chronicle api.ConsensusChronicles, config api.LocalNodeConfiguration) core.RoundStrategy {
	return NewRoundStrategy(
		phases.NewRegularPhaseBundleByDefault(),
		chronicle,
		config,
	)
}

type TransportFactory struct {
	cryptographyFactory api.TransportCryptographyFactory
	packetBuilder       api.PacketBuilder
	packetSender        api.PacketSender
}

func NewTransportFactory(
	cryptographyFactory api.TransportCryptographyFactory,
	packetBuilder api.PacketBuilder,
	packetSender api.PacketSender,
) *TransportFactory {
	return &TransportFactory{
		cryptographyFactory: cryptographyFactory,
		packetBuilder:       packetBuilder,
		packetSender:        packetSender,
	}
}

func (tf *TransportFactory) GetPacketSender() api.PacketSender {
	return tf.packetSender
}

func (tf *TransportFactory) GetPacketBuilder(signer cryptography_containers.DigestSigner) api.PacketBuilder {
	return tf.packetBuilder
}

func (tf *TransportFactory) GetCryptographyFactory() api.TransportCryptographyFactory {
	return tf.cryptographyFactory
}

type NodeProfileFactory struct {
	keyProcessor insolar.KeyProcessor
}

func NewNodeProfileFactory(keyProcessor insolar.KeyProcessor) *NodeProfileFactory {
	return &NodeProfileFactory{
		keyProcessor: keyProcessor,
	}
}

func (npf *NodeProfileFactory) createProfile(candidate gcp_types.BriefCandidateProfile, signature cryptography_containers.SignatureHolder, intro gcp_types.NodeIntroduction) *NodeIntroProfile {
	keyHolder := candidate.GetNodePK()
	pk, err := npf.keyProcessor.ImportPublicKeyBinary(keyHolder.AsBytes())
	if err != nil {
		panic(err)
	}

	store := NewECDSAPublicKeyStore(pk.(*ecdsa.PublicKey))

	return newNodeIntroProfile(
		candidate.GetNodeID(),
		candidate.GetNodePrimaryRole(),
		candidate.GetNodeSpecialRoles(),
		intro,
		candidate.GetNodeEndpoint(),
		store,
		keyHolder,
		signature,
	)
}

func (npf *NodeProfileFactory) CreateBriefIntroProfile(candidate gcp_types.BriefCandidateProfile) gcp_types.NodeIntroProfile {
	return npf.createProfile(candidate, candidate.GetJoinerSignature(), nil)
}

func (npf *NodeProfileFactory) CreateFullIntroProfile(candidate gcp_types.CandidateProfile) gcp_types.NodeIntroProfile {
	intro := newNodeIntroduction(candidate.GetNodeID(), candidate.GetReference())

	return npf.createProfile(candidate, candidate.GetJoinerSignature(), intro)
}

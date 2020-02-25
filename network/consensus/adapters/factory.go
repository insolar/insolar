// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package adapters

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle"
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

func (vf *ECDSASignatureVerifierFactory) CreateSignatureVerifierWithPKS(pks cryptkit.PublicKeyStore) cryptkit.SignatureVerifier {
	keyStore := pks.(*ECDSAPublicKeyStore)

	return NewECDSASignatureVerifier(
		vf.digester,
		vf.scheme,
		keyStore.publicKey,
	)
}

type TransportCryptographyFactory struct {
	verifierFactory *ECDSASignatureVerifierFactory
	digestFactory   *ConsensusDigestFactory
	scheme          insolar.PlatformCryptographyScheme
}

func NewTransportCryptographyFactory(scheme insolar.PlatformCryptographyScheme) *TransportCryptographyFactory {
	return &TransportCryptographyFactory{
		verifierFactory: NewECDSASignatureVerifierFactory(
			NewSha3512Digester(scheme),
			scheme,
		),
		digestFactory: NewConsensusDigestFactory(scheme),
		scheme:        scheme,
	}
}

func (cf *TransportCryptographyFactory) CreateSignatureVerifierWithPKS(pks cryptkit.PublicKeyStore) cryptkit.SignatureVerifier {
	return cf.verifierFactory.CreateSignatureVerifierWithPKS(pks)
}

func (cf *TransportCryptographyFactory) GetDigestFactory() transport.ConsensusDigestFactory {
	return cf.digestFactory
}

func (cf *TransportCryptographyFactory) CreateNodeSigner(sks cryptkit.SecretKeyStore) cryptkit.DigestSigner {
	ks := sks.(*ECDSASecretKeyStore)

	return NewECDSADigestSigner(ks.privateKey, cf.scheme)
}

func (cf *TransportCryptographyFactory) CreatePublicKeyStore(skh cryptkit.SignatureKeyHolder) cryptkit.PublicKeyStore {
	kh := skh.(*ECDSASignatureKeyHolder)

	return NewECDSAPublicKeyStore(kh.publicKey)
}

type RoundStrategyFactory struct {
	bundleFactory core.PhaseControllersBundleFactory
}

func NewRoundStrategyFactory() *RoundStrategyFactory {
	return &RoundStrategyFactory{
		bundleFactory: phasebundle.NewStandardBundleFactoryDefault(),
	}
}

func (rsf *RoundStrategyFactory) CreateRoundStrategy(online census.OnlinePopulation, config api.LocalNodeConfiguration) (core.RoundStrategy, core.PhaseControllersBundle) {
	rs := NewRoundStrategy(config)
	pcb := rsf.bundleFactory.CreateControllersBundle(online, config)
	return rs, pcb

}

type TransportFactory struct {
	cryptographyFactory transport.CryptographyAssistant
	packetBuilder       transport.PacketBuilder
	packetSender        transport.PacketSender
}

func NewTransportFactory(
	cryptographyFactory transport.CryptographyAssistant,
	packetBuilder transport.PacketBuilder,
	packetSender transport.PacketSender,
) *TransportFactory {
	return &TransportFactory{
		cryptographyFactory: cryptographyFactory,
		packetBuilder:       packetBuilder,
		packetSender:        packetSender,
	}
}

func (tf *TransportFactory) GetPacketSender() transport.PacketSender {
	return tf.packetSender
}

func (tf *TransportFactory) GetPacketBuilder(signer cryptkit.DigestSigner) transport.PacketBuilder {
	return tf.packetBuilder
}

func (tf *TransportFactory) GetCryptographyFactory() transport.CryptographyAssistant {
	return tf.cryptographyFactory
}

type keyStoreFactory struct {
	keyProcessor insolar.KeyProcessor
}

func (p *keyStoreFactory) CreatePublicKeyStore(keyHolder cryptkit.SignatureKeyHolder) cryptkit.PublicKeyStore {
	pk, err := p.keyProcessor.ImportPublicKeyBinary(keyHolder.AsBytes())
	if err != nil {
		panic(err)
	}
	return NewECDSAPublicKeyStore(pk.(*ecdsa.PublicKey))
}

func NewNodeProfileFactory(keyProcessor insolar.KeyProcessor) profiles.Factory {
	return profiles.NewSimpleProfileIntroFactory(&keyStoreFactory{keyProcessor})
}

type ConsensusDigestFactory struct {
	scheme insolar.PlatformCryptographyScheme
}

func NewConsensusDigestFactory(scheme insolar.PlatformCryptographyScheme) *ConsensusDigestFactory {
	return &ConsensusDigestFactory{
		scheme: scheme,
	}
}

func (cdf *ConsensusDigestFactory) CreatePacketDigester() cryptkit.DataDigester {
	return NewSha3512Digester(cdf.scheme)
}

func (cdf *ConsensusDigestFactory) CreateSequenceDigester() cryptkit.SequenceDigester {
	return NewSequenceDigester(NewSha3512Digester(cdf.scheme))
}

func (cdf *ConsensusDigestFactory) CreateAnnouncementDigester() cryptkit.SequenceDigester {
	return NewSequenceDigester(NewSha3512Digester(cdf.scheme))
}

func (cdf *ConsensusDigestFactory) CreateGlobulaStateDigester() transport.StateDigester {
	return NewStateDigester(
		NewSequenceDigester(NewSha3512Digester(cdf.scheme)),
	)
}

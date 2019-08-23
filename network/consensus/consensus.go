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

package consensus

import (
	"context"
	"fmt"
	"reflect"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	transport2 "github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/serialization"
	"github.com/insolar/insolar/network/transport"
)

type packetProcessorSetter interface {
	SetPacketProcessor(adapters.PacketProcessor)
	SetPacketParserFactory(factory adapters.PacketParserFactory)
}

type Mode uint

const (
	ReadyNetwork = Mode(iota)
	Joiner
)

func New(ctx context.Context, dep Dep) Installer {
	dep.verify()

	constructor := newConstructor(ctx, &dep)
	constructor.verify()

	return newInstaller(constructor, &dep)
}

func verify(s interface{}) {
	cdValue := reflect.Indirect(reflect.ValueOf(s))
	cdType := cdValue.Type()

	for i := 0; i < cdValue.NumField(); i++ {
		fieldMeta := cdValue.Field(i)

		if (fieldMeta.Kind() == reflect.Interface || fieldMeta.Kind() == reflect.Ptr) && fieldMeta.IsNil() {
			panic(fmt.Sprintf("%s field %s is nil", cdType.Name(), cdType.Field(i).Name))
		}
	}
}

type Dep struct {
	KeyProcessor       insolar.KeyProcessor
	Scheme             insolar.PlatformCryptographyScheme
	CertificateManager insolar.CertificateManager
	KeyStore           insolar.KeyStore

	NodeKeeper        network.NodeKeeper
	DatagramTransport transport.DatagramTransport

	StateGetter         adapters.StateGetter
	PulseChanger        adapters.PulseChanger
	StateUpdater        adapters.StateUpdater
	EphemeralController adapters.EphemeralController
}

func (cd *Dep) verify() {
	verify(cd)
}

type constructor struct {
	consensusConfiguration       census.ConsensusConfiguration
	mandateRegistry              census.MandateRegistry
	misbehaviorRegistry          census.MisbehaviorRegistry
	offlinePopulation            census.OfflinePopulation
	versionedRegistries          census.VersionedRegistries
	nodeProfileFactory           profiles.Factory
	localNodeConfiguration       api.LocalNodeConfiguration
	roundStrategyFactory         core.RoundStrategyFactory
	transportCryptographyFactory transport2.CryptographyAssistant
	packetBuilder                transport2.PacketBuilder
	packetSender                 transport2.PacketSender
	transportFactory             transport2.Factory
}

func newConstructor(ctx context.Context, dep *Dep) *constructor {
	c := &constructor{}

	c.consensusConfiguration = adapters.NewConsensusConfiguration()
	c.mandateRegistry = adapters.NewMandateRegistry(
		cryptkit.NewDigest(
			longbits.NewBits512FromBytes(
				make([]byte, 64),
			),
			adapters.SHA3512Digest,
		).AsDigestHolder(),
		c.consensusConfiguration,
	)
	c.misbehaviorRegistry = adapters.NewMisbehaviorRegistry()
	c.offlinePopulation = adapters.NewOfflinePopulation(
		dep.NodeKeeper,
		dep.CertificateManager,
		dep.KeyProcessor,
	)
	c.versionedRegistries = adapters.NewVersionedRegistries(
		c.mandateRegistry,
		c.misbehaviorRegistry,
		c.offlinePopulation,
	)
	c.nodeProfileFactory = adapters.NewNodeProfileFactory(dep.KeyProcessor)
	c.localNodeConfiguration = adapters.NewLocalNodeConfiguration(
		ctx,
		dep.KeyStore,
	)
	c.roundStrategyFactory = adapters.NewRoundStrategyFactory()
	c.transportCryptographyFactory = adapters.NewTransportCryptographyFactory(dep.Scheme)
	c.packetBuilder = serialization.NewPacketBuilder(
		c.transportCryptographyFactory,
		c.localNodeConfiguration,
	)
	c.packetSender = adapters.NewPacketSender(dep.DatagramTransport)
	c.transportFactory = adapters.NewTransportFactory(
		c.transportCryptographyFactory,
		c.packetBuilder,
		c.packetSender,
	)

	return c
}

func (c *constructor) verify() {
	verify(c)
}

type Installer struct {
	dep       *Dep
	consensus *constructor
}

func newInstaller(constructor *constructor, dep *Dep) Installer {
	return Installer{
		dep:       dep,
		consensus: constructor,
	}
}

func (c Installer) ControllerFor(mode Mode, setters ...packetProcessorSetter) Controller {
	controlFeederInterceptor := adapters.InterceptConsensusControl(
		adapters.NewConsensusControlFeeder(),
	)
	candidateFeeder := &coreapi.SequentialCandidateFeeder{}

	var ephemeralFeeder api.EphemeralControlFeeder
	if c.dep.EphemeralController.EphemeralMode(c.dep.NodeKeeper.GetAccessor(insolar.GenesisPulse.PulseNumber).GetActiveNodes()) {
		ephemeralFeeder = adapters.NewEphemeralControlFeeder(c.dep.EphemeralController)
	}

	upstreamController := adapters.NewUpstreamPulseController(
		c.dep.StateGetter,
		c.dep.PulseChanger,
		c.dep.StateUpdater,
	)

	consensusChronicles := c.createConsensusChronicles(mode)
	consensusController := c.createConsensusController(
		consensusChronicles,
		controlFeederInterceptor.Feeder(),
		candidateFeeder,
		ephemeralFeeder,
		upstreamController,
	)
	packetParserFactory := c.createPacketParserFactory()

	c.bind(setters, consensusController, packetParserFactory)

	consensusController.Prepare()

	return newController(controlFeederInterceptor, candidateFeeder, consensusController, upstreamController)
}

func (c *Installer) createCensus(mode Mode) *censusimpl.PrimingCensusTemplate {
	certificate := c.dep.CertificateManager.GetCertificate()
	origin := c.dep.NodeKeeper.GetOrigin()
	knownNodes := c.dep.NodeKeeper.GetAccessor(insolar.GenesisPulse.PulseNumber).GetActiveNodes()

	node := adapters.NewStaticProfile(origin, certificate, c.dep.KeyProcessor)
	nodes := adapters.NewStaticProfileList(knownNodes, certificate, c.dep.KeyProcessor)

	if mode == Joiner {
		return adapters.NewCensusForJoiner(
			node,
			c.consensus.versionedRegistries,
			c.consensus.transportCryptographyFactory,
		)
	}

	return adapters.NewCensus(
		node,
		nodes,
		c.consensus.versionedRegistries,
		c.consensus.transportCryptographyFactory,
	)
}

func (c *Installer) createConsensusChronicles(mode Mode) censusimpl.LocalConsensusChronicles {
	consensusChronicles := adapters.NewChronicles(c.consensus.nodeProfileFactory)
	c.createCensus(mode).SetAsActiveTo(consensusChronicles)
	return consensusChronicles
}

func (c *Installer) createConsensusController(
	consensusChronicles censusimpl.LocalConsensusChronicles,
	controlFeeder api.ConsensusControlFeeder,
	candidateFeeder api.CandidateControlFeeder,
	ephemeralFeeder api.EphemeralControlFeeder,
	upstreamController api.UpstreamController,
) api.ConsensusController {
	return gcpv2.NewConsensusMemberController(
		consensusChronicles,
		upstreamController,
		core.NewPhasedRoundControllerFactory(
			c.consensus.localNodeConfiguration,
			c.consensus.transportFactory,
			c.consensus.roundStrategyFactory,
		),
		candidateFeeder,
		controlFeeder,
		ephemeralFeeder,
	)
}

func (c *Installer) createPacketParserFactory() adapters.PacketParserFactory {
	return serialization.NewPacketParserFactory(
		c.consensus.transportCryptographyFactory.GetDigestFactory().CreatePacketDigester(),
		c.consensus.transportCryptographyFactory.CreateNodeSigner(c.consensus.localNodeConfiguration.GetSecretKeyStore()).GetSignMethod(),
		c.dep.KeyProcessor,
	)
}

func (c *Installer) bind(
	setters []packetProcessorSetter,
	packetProcessor adapters.PacketProcessor,
	packetParserFactory adapters.PacketParserFactory,
) {
	for _, setter := range setters {
		setter.SetPacketProcessor(packetProcessor)
		setter.SetPacketParserFactory(packetParserFactory)
	}
}

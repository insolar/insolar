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

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	transport2 "github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/gcpv2"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/transport"
)

type Dep struct {
	PrimingCloudStateHash [64]byte

	KeyProcessor       insolar.KeyProcessor
	Scheme             insolar.PlatformCryptographyScheme
	CertificateManager insolar.CertificateManager
	KeyStore           insolar.KeyStore
	NodeKeeper         network.NodeKeeper
	DatagramTransport  transport.DatagramTransport

	StateGetter  adapters.StateGetter
	PulseChanger adapters.PulseChanger
	StateUpdater adapters.StateUpdater

	// TODO: remove it from here
	PacketBuilder func(transport2.CryptographyFactory, api.LocalNodeConfiguration) transport2.PacketBuilder
	PacketSender  transport2.PacketSender
}

func (cd *Dep) verify() {
	verify(cd)
}

type Consensus struct {
	population                   censusimpl.ManyNodePopulation // TODO: there should be interface
	consensusConfiguration       census.ConsensusConfiguration // TODO: there should be interface
	mandateRegistry              census.MandateRegistry
	misbehaviorRegistry          census.MisbehaviorRegistry
	offlinePopulation            census.OfflinePopulation
	versionedRegistries          census.VersionedRegistries
	nodeProfileFactory           profiles.Factory
	consensusChronicles          api.ConsensusChronicles
	localNodeConfiguration       api.LocalNodeConfiguration
	upstreamPulseController      api.UpstreamController
	roundStrategyFactory         core.RoundStrategyFactory
	transportCryptographyFactory transport2.CryptographyFactory
	packetBuilder                transport2.PacketBuilder
	packetSender                 transport2.PacketSender
	transportFactory             transport2.Factory
	consensusController          api.ConsensusController
}

func New(ctx context.Context, dep Dep) Consensus {
	dep.verify()

	consensus := Consensus{}

	certificate := dep.CertificateManager.GetCertificate()
	origin := dep.NodeKeeper.GetOrigin()
	knownNodes := dep.NodeKeeper.GetAccessor().GetActiveNodes()

	consensus.population = adapters.NewPopulation(
		adapters.NewNodeIntroProfile(origin, certificate, dep.KeyProcessor),
		adapters.NewNodeIntroProfileList(knownNodes, certificate, dep.KeyProcessor),
	)
	consensus.consensusConfiguration = adapters.NewConsensusConfiguration()
	consensus.mandateRegistry = adapters.NewMandateRegistry(
		cryptkit.NewDigest(
			longbits.NewBits512FromBytes(
				dep.PrimingCloudStateHash[:],
			),
			adapters.SHA3512Digest,
		).AsDigestHolder(),
		consensus.consensusConfiguration,
	)
	consensus.misbehaviorRegistry = adapters.NewMisbehaviorRegistry()
	consensus.offlinePopulation = adapters.NewOfflinePopulation(
		dep.NodeKeeper,
		dep.CertificateManager,
		dep.KeyProcessor,
	)
	consensus.versionedRegistries = adapters.NewVersionedRegistries(
		consensus.mandateRegistry,
		consensus.misbehaviorRegistry,
		consensus.offlinePopulation,
	)
	consensus.nodeProfileFactory = adapters.NewNodeProfileFactory(dep.KeyProcessor)
	consensus.consensusChronicles = adapters.NewChronicles(
		consensus.population,
		consensus.nodeProfileFactory,
		consensus.versionedRegistries,
	)
	consensus.localNodeConfiguration = adapters.NewLocalNodeConfiguration(
		ctx,
		dep.KeyStore,
	)
	consensus.upstreamPulseController = adapters.NewUpstreamPulseController(
		dep.StateGetter,
		dep.PulseChanger,
		dep.StateUpdater,
	)
	consensus.roundStrategyFactory = adapters.NewRoundStrategyFactory()
	consensus.transportCryptographyFactory = adapters.NewTransportCryptographyFactory(dep.Scheme)
	consensus.packetBuilder = dep.PacketBuilder(
		consensus.transportCryptographyFactory,
		consensus.localNodeConfiguration,
	)
	// TODO: comment until serialization ready
	// consensus.packetSender = NewPacketSender(dep.DatagramTransport)
	consensus.packetSender = dep.PacketSender
	consensus.transportFactory = adapters.NewTransportFactory(
		consensus.transportCryptographyFactory,
		consensus.packetBuilder,
		consensus.packetSender,
	)
	consensus.consensusController = gcpv2.NewConsensusMemberController(
		consensus.consensusChronicles,
		consensus.upstreamPulseController,
		core.NewPhasedRoundControllerFactory(
			consensus.localNodeConfiguration,
			consensus.transportFactory,
			consensus.roundStrategyFactory,
		),
		&core.SequentialCandidateFeeder{},
		adapters.NewConsensusControlFeeder(),
	)

	consensus.verify()
	return consensus
}

type packetProcessorSetter interface {
	SetPacketProcessor(adapters.PacketProcessor)
}

type Controller interface {
	// Leave
	// SetPower
}

func (c Consensus) Install(setters ...packetProcessorSetter) Controller {
	for _, setter := range setters {
		setter.SetPacketProcessor(c.consensusController)
	}
	return c.consensusController
}

func (c *Consensus) verify() {
	verify(c)
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

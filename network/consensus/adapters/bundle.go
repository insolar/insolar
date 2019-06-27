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
	"context"
	"fmt"
	"reflect"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	common2 "github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/transport"
)

type ConsensusDep struct {
	PrimingCloudStateHash [64]byte

	Scheme             insolar.PlatformCryptographyScheme
	CertificateManager insolar.CertificateManager
	KeyStore           insolar.KeyStore
	NodeKeeper         network.NodeKeeper
	DatagramTransport  transport.DatagramTransport

	Stater       stater
	PulseChanger pulseChanger

	// TODO: remove it from here
	PacketBuilder func(core.TransportCryptographyFactory, core.LocalNodeConfiguration) core.PacketBuilder
	PacketSender  core.PacketSender
}

func (cd *ConsensusDep) verify() {
	verify(cd)
}

type ConsensusBundle struct {
	population                   census.ManyNodePopulation // TODO: there should be interface
	mandateRegistry              census.MandateRegistry
	misbehaviorRegistry          census.MisbehaviorRegistry
	offlinePopulation            census.OfflinePopulation
	versionedRegistries          census.VersionedRegistries
	consensusChronicles          census.ConsensusChronicles
	localNodeConfiguration       core.LocalNodeConfiguration
	upstreamPulseController      core.UpstreamPulseController
	roundStrategyFactory         core.RoundStrategyFactory
	transportCryptographyFactory core.TransportCryptographyFactory
	packetBuilder                core.PacketBuilder
	packetSender                 core.PacketSender
	transportFactory             core.TransportFactory
	consensusController          core.ConsensusController
}

func NewConsensusBundle(ctx context.Context, dep ConsensusDep) ConsensusBundle {
	dep.verify()

	bundle := ConsensusBundle{}

	certificate := dep.CertificateManager.GetCertificate()
	origin := dep.NodeKeeper.GetOrigin()
	knownNodes := dep.NodeKeeper.GetAccessor().GetActiveNodes()

	bundle.population = NewPopulation(NewNodeIntroProfile(origin, certificate), NewNodeIntroProfileList(knownNodes, certificate))
	bundle.mandateRegistry = NewMandateRegistry(
		common2.NewDigest(
			common2.NewBits512FromBytes(dep.PrimingCloudStateHash[:]), SHA3512Digest,
		).AsDigestHolder(),
	)
	bundle.misbehaviorRegistry = NewMisbehaviorRegistry()
	bundle.offlinePopulation = NewOfflinePopulation(dep.NodeKeeper, dep.CertificateManager)
	bundle.versionedRegistries = NewVersionedRegistries(
		bundle.mandateRegistry,
		bundle.misbehaviorRegistry,
		bundle.offlinePopulation,
	)
	bundle.consensusChronicles = NewChronicles(bundle.population, bundle.versionedRegistries)
	bundle.localNodeConfiguration = NewLocalNodeConfiguration(ctx, dep.KeyStore)
	bundle.upstreamPulseController = NewUpstreamPulseController(dep.Stater, dep.PulseChanger)
	bundle.roundStrategyFactory = NewRoundStrategyFactory()
	bundle.transportCryptographyFactory = NewTransportCryptographyFactory(dep.Scheme)
	bundle.packetBuilder = dep.PacketBuilder(bundle.transportCryptographyFactory, bundle.localNodeConfiguration)
	// TODO: comment until serialization ready
	// bundle.packetSender = NewPacketSender(dep.DatagramTransport)
	bundle.packetSender = dep.PacketSender
	bundle.transportFactory = NewTransportFactory(
		bundle.transportCryptographyFactory,
		bundle.packetBuilder,
		bundle.packetSender,
	)
	bundle.consensusController = gcpv2.NewConsensusMemberController(
		bundle.consensusChronicles,
		bundle.upstreamPulseController,
		core.NewPhasedRoundControllerFactory(
			bundle.localNodeConfiguration,
			bundle.transportFactory,
			bundle.roundStrategyFactory,
		),
	)

	bundle.verify()
	return bundle
}

func (cb *ConsensusBundle) Controller() core.ConsensusController {
	return cb.consensusController
}

func (cb *ConsensusBundle) verify() {
	verify(cb)
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

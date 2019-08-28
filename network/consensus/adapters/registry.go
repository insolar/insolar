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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

type MisbehaviorRegistry struct{}

func NewMisbehaviorRegistry() *MisbehaviorRegistry {
	return &MisbehaviorRegistry{}
}

func (mr *MisbehaviorRegistry) AddReport(report misbehavior.Report) {
}

type MandateRegistry struct {
	cloudHash              proofs.CloudStateHash
	consensusConfiguration census.ConsensusConfiguration
}

func NewMandateRegistry(cloudHash proofs.CloudStateHash, consensusConfiguration census.ConsensusConfiguration) *MandateRegistry {
	return &MandateRegistry{
		cloudHash:              cloudHash,
		consensusConfiguration: consensusConfiguration,
	}
}

func (mr *MandateRegistry) FindRegisteredProfile(host endpoints.Inbound) profiles.Host {
	return nil
}

func (mr *MandateRegistry) GetCloudIdentity() cryptkit.DigestHolder {
	return mr.cloudHash
}

func (mr *MandateRegistry) GetConsensusConfiguration() census.ConsensusConfiguration {
	return mr.consensusConfiguration
}

func (mr *MandateRegistry) GetPrimingCloudHash() proofs.CloudStateHash {
	return mr.cloudHash
}

type OfflinePopulation struct {
	// TODO: should't use nodekeeper here.
	nodeKeeper   network.NodeKeeper
	manager      insolar.CertificateManager
	keyProcessor insolar.KeyProcessor
}

func NewOfflinePopulation(nodeKeeper network.NodeKeeper, manager insolar.CertificateManager, keyProcessor insolar.KeyProcessor) *OfflinePopulation {
	return &OfflinePopulation{
		nodeKeeper:   nodeKeeper,
		manager:      manager,
		keyProcessor: keyProcessor,
	}
}

func (op *OfflinePopulation) FindRegisteredProfile(identity endpoints.Inbound) profiles.Host {
	node := op.nodeKeeper.GetAccessor(insolar.GenesisPulse.PulseNumber).GetActiveNodeByAddr(identity.GetNameAddress().String())
	if node == nil {
		return nil
	}
	cert := op.manager.GetCertificate()
	return NewStaticProfile(node, cert, op.keyProcessor)
}

type VersionedRegistries struct {
	mandateRegistry     census.MandateRegistry
	misbehaviorRegistry census.MisbehaviorRegistry
	offlinePopulation   census.OfflinePopulation

	pulseData pulse.Data
}

func (c *VersionedRegistries) GetNearestValidPulseData() pulse.Data {
	return pulse.Data{}
}

func NewVersionedRegistries(
	mandateRegistry census.MandateRegistry,
	misbehaviorRegistry census.MisbehaviorRegistry,
	offlinePopulation census.OfflinePopulation,
) *VersionedRegistries {
	return &VersionedRegistries{
		mandateRegistry:     mandateRegistry,
		misbehaviorRegistry: misbehaviorRegistry,
		offlinePopulation:   offlinePopulation,
	}
}

func (c *VersionedRegistries) CommitNextPulse(pd pulse.Data, population census.OnlinePopulation) census.VersionedRegistries {
	pd.EnsurePulseData()
	cp := *c
	cp.pulseData = pd
	return &cp
}

func (c *VersionedRegistries) GetMisbehaviorRegistry() census.MisbehaviorRegistry {
	return c.misbehaviorRegistry
}

func (c *VersionedRegistries) GetMandateRegistry() census.MandateRegistry {
	return c.mandateRegistry
}

func (c *VersionedRegistries) GetOfflinePopulation() census.OfflinePopulation {
	return c.offlinePopulation
}

func (c *VersionedRegistries) GetVersionPulseData() pulse.Data {
	return c.pulseData
}

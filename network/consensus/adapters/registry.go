// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

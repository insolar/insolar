// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package core

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

var _ api.RoundControllerFactory = &PhasedRoundControllerFactory{}

func NewPhasedRoundControllerFactory(config api.LocalNodeConfiguration, t transport.Factory, strategyFactory RoundStrategyFactory) *PhasedRoundControllerFactory {
	return &PhasedRoundControllerFactory{strategyFactory: strategyFactory, transport: t, config: config}
}

type PhasedRoundControllerFactory struct {
	strategyFactory RoundStrategyFactory
	transport       transport.Factory
	config          api.LocalNodeConfiguration
}

func (c *PhasedRoundControllerFactory) GetLocalConfiguration() api.LocalNodeConfiguration {
	return c.config
}

func (c *PhasedRoundControllerFactory) CreateConsensusRound(chronicle api.ConsensusChronicles, controlFeeder api.ConsensusControlFeeder,
	candidateFeeder api.CandidateControlFeeder, ephemeralFeeder api.EphemeralControlFeeder) api.RoundController {

	latest, _ := chronicle.GetLatestCensus()
	strategy, bundle := c.strategyFactory.CreateRoundStrategy(latest.GetOnlinePopulation(), c.config)
	return NewPhasedRoundController(strategy, chronicle, bundle, c.transport, c.config, controlFeeder, candidateFeeder, ephemeralFeeder)
}

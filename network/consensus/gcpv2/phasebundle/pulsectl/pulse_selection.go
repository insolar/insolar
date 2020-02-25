// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsectl

import (
	"context"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type PulseSelectionStrategyFactory interface {
	CreatePulseSelectionStrategy(population census.OnlinePopulation, config api.LocalNodeConfiguration) PulseSelectionStrategy
}

type PulseSelectionStrategy interface {
	HandlePulsarPacket(ctx context.Context, p transport.PulsePacketReader,
		from endpoints.Inbound, fromPulsar bool) (bool, error)
}

var _ PulseSelectionStrategyFactory = &takeFirstStrategyFactory{}

func NewTakeFirstSelectionStrategyFactory() PulseSelectionStrategyFactory {
	return &takeFirstStrategyFactory{}
}

type takeFirstStrategyFactory struct {
}

func (p *takeFirstStrategyFactory) CreatePulseSelectionStrategy(population census.OnlinePopulation,
	config api.LocalNodeConfiguration) PulseSelectionStrategy {

	// if population.GetLocalProfile().IsJoiner() {
	return p
}

func (p *takeFirstStrategyFactory) HandlePulsarPacket(ctx context.Context, packet transport.PulsePacketReader,
	from endpoints.Inbound, fromPulsar bool) (bool, error) {

	return true, nil
}

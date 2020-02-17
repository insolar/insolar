// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package core

import (
	"context"

	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"
)

type PerNodePacketDispatcherFactory interface {
	CreatePerNodePacketHandler(perNodeContext context.Context, node *population.NodeAppearance) (context.Context, population.DispatchMemberPacketFunc)
}

type PrepPhaseController interface {
	GetPacketType() []phases.PacketType
	CreatePacketDispatcher(pt phases.PacketType, realm *PrepRealm) population.PacketDispatcher

	BeforeStart(ctx context.Context, realm *PrepRealm)
	StartWorker(ctx context.Context, realm *PrepRealm)
}

/* realm is provided for this handler to avoid being replicated in individual handlers */
type PhaseController interface {
	GetPacketType() []phases.PacketType
	CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *FullRealm) (population.PacketDispatcher, PerNodePacketDispatcherFactory)

	BeforeStart(ctx context.Context, realm *FullRealm)
	StartWorker(ctx context.Context, realm *FullRealm)
}

type PhaseControllersBundle interface {
	IsDynamicPopulationRequired() bool
	CreatePrepPhaseControllers() []PrepPhaseController
	CreateFullPhaseControllers(nodeCount int) ([]PhaseController, NodeUpdateCallback)
}

type PhaseControllersBundleFactory interface {
	CreateControllersBundle(population census.OnlinePopulation, config api.LocalNodeConfiguration /* strategy RoundStrategy */) PhaseControllersBundle
}

type NodeUpdateCallback interface {
	population.EventDispatcher
}

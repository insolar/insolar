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

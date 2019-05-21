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

package bootstrap

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

type NetworkBootstrapper interface {
	Bootstrap(ctx context.Context) (*network.BootstrapResult, error)
	SetLastPulse(number insolar.PulseNumber)
	GetLastPulse() insolar.PulseNumber
	AuthenticateToDiscoveryNode(ctx context.Context, discovery *DiscoveryNode) error
}

type networkBootstrapper struct {
	Certificate    insolar.Certificate     `inject:""`
	Bootstrapper   Bootstrapper            `inject:""`
	NodeKeeper     network.NodeKeeper      `inject:""`
	SessionManager SessionManager          `inject:""`
	AuthController AuthorizationController `inject:""`
	Gatewayer      network.Gatewayer       `inject:""`
}

func (nb *networkBootstrapper) Bootstrap(ctx context.Context) (*network.BootstrapResult, error) {
	ctx, span := instracer.StartSpan(ctx, "NetworkBootstrapper.Bootstrap")
	defer span.End()
	if len(nb.Certificate.GetDiscoveryNodes()) == 0 {
		return nb.Bootstrapper.ZeroBootstrap(ctx)
	}
	var err error
	var result *network.BootstrapResult
	if utils.OriginIsDiscovery(nb.Certificate) {
		result, err = nb.bootstrapDiscovery(ctx)
		// if the network is up and complete, we return discovery nodes via consensus
		if err == ErrReconnectRequired {
			log.Debugf("[ Bootstrap ] Connecting discovery node %s as joiner", nb.NodeKeeper.GetOrigin().ID())
			nb.NodeKeeper.GetOrigin().(node.MutableNode).SetState(insolar.NodePending)
			result, err = nb.bootstrapJoiner(ctx)
		}
	} else {
		result, err = nb.bootstrapJoiner(ctx)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap")
	}
	nb.NodeKeeper.SetIsBootstrapped(true)
	return result, nil
}

func (nb *networkBootstrapper) SetLastPulse(number insolar.PulseNumber) {
	nb.Bootstrapper.SetLastPulse(number)
}

func (nb *networkBootstrapper) GetLastPulse() insolar.PulseNumber {
	return nb.Bootstrapper.GetLastPulse()
}

func (nb *networkBootstrapper) bootstrapJoiner(ctx context.Context) (*network.BootstrapResult, error) {
	ctx, span := instracer.StartSpan(ctx, "NetworkBootstrapper.bootstrapJoiner")
	defer span.End()
	nb.NodeKeeper.GetConsensusInfo().SetIsJoiner(true)
	result, discoveryNode, err := nb.Bootstrapper.Bootstrap(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error bootstrapping to discovery node")
	}
	return result, nb.AuthenticateToDiscoveryNode(ctx, discoveryNode)
}

func (nb *networkBootstrapper) AuthenticateToDiscoveryNode(ctx context.Context, discovery *DiscoveryNode) error {
	data, err := nb.AuthController.Authorize(ctx, discovery, nb.Certificate)
	if err != nil {
		return errors.Wrap(err, "Error authorizing on discovery node")
	}
	// TODO: fix Short ID assignment logic
	// origin := nb.NodeKeeper.GetOrigin()
	// mutableOrigin := origin.(nodenetwork.MutableNode)
	// mutableOrigin.SetShortID(data.AssignShortID)
	return nb.AuthController.Register(ctx, discovery, SessionID(data.SessionID))
}

func (nb *networkBootstrapper) bootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
	return nb.Bootstrapper.BootstrapDiscovery(ctx)
}

func NewNetworkBootstrapper() NetworkBootstrapper {
	return &networkBootstrapper{}
}

/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package bootstrap

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

type NetworkBootstrapper interface {
	Bootstrap(ctx context.Context) (*network.BootstrapResult, error)
	SetLastPulse(number core.PulseNumber)
	GetLastPulse() core.PulseNumber
}

type networkBootstrapper struct {
	Certificate         core.Certificate            `inject:""`
	Bootstrapper        Bootstrapper                `inject:""`
	NodeKeeper          network.NodeKeeper          `inject:""`
	SessionManager      SessionManager              `inject:""`
	AuthController      AuthorizationController     `inject:""`
	ChallengeController ChallengeResponseController `inject:""`
}

func (nb *networkBootstrapper) Bootstrap(ctx context.Context) (*network.BootstrapResult, error) {
	ctx, span := instracer.StartSpan(ctx, "NetworkBootstrapper.Bootstrap")
	defer span.End()
	if len(nb.Certificate.GetDiscoveryNodes()) == 0 {
		host, err := host.NewHostN(nb.NodeKeeper.GetOrigin().Address(), nb.NodeKeeper.GetOrigin().ID())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create a host")
		}
		log.Info("[ Bootstrap ] Zero bootstrap")
		return &network.BootstrapResult{
			Host: host,
			// FirstPulseTime: nb.Bootstrapper.GetFirstFakePulseTime(),
		}, nil
	}
	var err error
	var result *network.BootstrapResult
	if utils.OriginIsDiscovery(nb.Certificate) {
		result, err = nb.bootstrapDiscovery(ctx)
		// if the network is up and complete, we return discovery nodes via consensus
		if err == ErrReconnectRequired {
			log.Debugf("[ Bootstrap ] Connecting discovery node %s as joiner", nb.NodeKeeper.GetOrigin().ID())
			nb.NodeKeeper.SetState(core.WaitingNodeNetworkState)
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

func (nb *networkBootstrapper) SetLastPulse(number core.PulseNumber) {
	nb.Bootstrapper.SetLastPulse(number)
}

func (nb *networkBootstrapper) GetLastPulse() core.PulseNumber {
	return nb.Bootstrapper.GetLastPulse()
}

func (nb *networkBootstrapper) bootstrapJoiner(ctx context.Context) (*network.BootstrapResult, error) {
	ctx, span := instracer.StartSpan(ctx, "NetworkBootstrapper.bootstrapJoiner")
	defer span.End()
	result, discoveryNode, err := nb.Bootstrapper.Bootstrap(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error bootstrapping to discovery node")
	}
	sessionID, err := nb.AuthController.Authorize(ctx, discoveryNode, nb.Certificate)
	if err != nil {
		return nil, errors.Wrap(err, "Error authorizing on discovery node")
	}

	_, err = nb.ChallengeController.Execute(ctx, discoveryNode, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "Error executing double challenge response")
	}
	// TODO: fix Short ID assignment logic
	// origin := nb.NodeKeeper.GetOrigin()
	// mutableOrigin := origin.(nodenetwork.MutableNode)
	// mutableOrigin.SetShortID(data.AssignShortID)
	return result, nb.AuthController.Register(ctx, discoveryNode, sessionID)
}

func (nb *networkBootstrapper) bootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
	return nb.Bootstrapper.BootstrapDiscovery(ctx)
}

func NewNetworkBootstrapper() NetworkBootstrapper {
	return &networkBootstrapper{}
}

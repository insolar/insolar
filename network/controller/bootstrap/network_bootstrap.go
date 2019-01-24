/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package bootstrap

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
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
	if len(nb.Certificate.GetDiscoveryNodes()) == 0 {
		host, err := host.NewHostN(nb.NodeKeeper.GetOrigin().Address(), nb.NodeKeeper.GetOrigin().ID())
		if err != nil {
			return nil, errors.Wrap(err, "[ Bootstrap ] failed to create a host")
		}
		log.Info("Zero bootstrap")
		return &network.BootstrapResult{
			Host:           host,
			FirstPulseTime: nb.Bootstrapper.GetFirstFakePulseTime(),
		}, nil
	}
	if utils.OriginIsDiscovery(nb.Certificate) {
		result, err := nb.bootstrapDiscovery(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "[ Bootstrap ] Couldn't OriginIsDiscovery")
		}
		nb.NodeKeeper.SetIsBootstrapped(true)
		return result, nil
	}
	return nb.bootstrapJoiner(ctx)
}

func (nb *networkBootstrapper) SetLastPulse(number core.PulseNumber) {
	nb.Bootstrapper.SetLastPulse(number)
}

func (nb *networkBootstrapper) GetLastPulse() core.PulseNumber {
	return nb.Bootstrapper.GetLastPulse()
}

func (nb *networkBootstrapper) bootstrapJoiner(ctx context.Context) (*network.BootstrapResult, error) {
	result, discoveryNode, err := nb.Bootstrapper.Bootstrap(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error bootstrapping to discovery node")
	}
	sessionID, err := nb.AuthController.Authorize(ctx, discoveryNode, nb.Certificate)
	if err != nil {
		return nil, errors.Wrap(err, "Error authorizing on discovery node")
	}

	data, err := nb.ChallengeController.Execute(ctx, discoveryNode, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "Error executing double challenge response")
	}
	origin := nb.NodeKeeper.GetOrigin()
	mutableOrigin := origin.(nodenetwork.MutableNode)
	mutableOrigin.SetShortID(data.AssignShortID)
	return result, nb.AuthController.Register(ctx, discoveryNode, sessionID)
}

func (nb *networkBootstrapper) bootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
	return nb.Bootstrapper.BootstrapDiscovery(ctx)
}

func NewNetworkBootstrapper() NetworkBootstrapper {
	return &networkBootstrapper{}
}

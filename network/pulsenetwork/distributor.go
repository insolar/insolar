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

package pulsenetwork

import (
	"context"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/pkg/errors"
)
type distributor struct {
	Transport transport.Transport `inject:""`

	pulsarHost     *host.Host
	bootstrapHosts []*host.Host
}

func NewDistributor(conf configuration.PulseDistributor) (core.PulseDistributor, error) {
	bootstrapHosts := make([]*host.Host, len(conf.BootstrapHosts))

	for _, node := range conf.BootstrapHosts {
		bootstrapHost, err := host.NewHost(node)
		if err != nil {
			return nil, errors.Wrap(err, "[ NewDistributor ] failed to create bootstrap node host")
		}
		bootstrapHosts = append(bootstrapHosts, bootstrapHost)
	}

	return &distributor{
		bootstrapHosts: bootstrapHosts,
	}, nil
}

func (d *distributor) Start(ctx context.Context) error {
	pulsarHost, err := host.NewHost(d.Transport.PublicAddress())
	if err != nil {
		return errors.Wrap(err, "[ NewDistributor ] failed to create pulsar host")
	}
	pulsarHost.NodeID = core.RecordRef{}

	return nil
}

func (d *distributor) Distribute(ctx context.Context, pulse *core.Pulse) {
	panic("not implemented")
func (d *distributor) pause(ctx context.Context) {
	inslogger.FromContext(ctx).Info("[ Pause ] Pause distribution, stopping transport")
	d.Transport.Stop()
}

func (d *distributor) resume(ctx context.Context) {
	inslogger.FromContext(ctx).Info("[ Resume ] Resume distribution, starting transport")

	go func(ctx context.Context, t transport.Transport) {
		err := t.Start(ctx)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
	}(ctx, d.Transport)
}

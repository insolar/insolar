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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type distributor struct {
	Transport transport.Transport `inject:""`

	pingTimeout        time.Duration
	randomNodesTimeout time.Duration
	pulseTimeout       time.Duration
	randomNodesCount   int

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
		pingTimeout:        time.Duration(conf.PingTimeout) * time.Millisecond,
		randomNodesTimeout: time.Duration(conf.RandomNodesTimeout) * time.Millisecond,
		pulseTimeout:       time.Duration(conf.PulseTimeout) * time.Millisecond,
		randomNodesCount:   conf.RandomNodesCount,

		bootstrapHosts: bootstrapHosts,
	}, nil
}

func (d *distributor) Start(ctx context.Context) error {
	pulsarHost, err := host.NewHost(d.Transport.PublicAddress())
	if err != nil {
		return errors.Wrap(err, "[ NewDistributor ] failed to create pulsar host")
	}
	pulsarHost.NodeID = core.RecordRef{}

	d.pause(ctx)
	return nil
}

func (d *distributor) Distribute(ctx context.Context, pulse *core.Pulse) {
	logger := inslogger.FromContext(ctx)

	d.resume(ctx)
	defer d.pause(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("sendPulseToNetwork failed with panic: %v", r)
		}
	}()

func (d *distributor) sendPulseToHost(ctx context.Context, pulse *core.Pulse, host *host.Host) error {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("sendPulseToHost failed with panic: %v", x)
		}
	}()

	pb := packet.NewBuilder(d.pulsarHost)
	pulseRequest := pb.Receiver(host).Request(&packet.RequestPulse{Pulse: *pulse}).Type(types.Pulse).Build()
	call, err := d.Transport.SendRequest(pulseRequest)
	if err != nil {
		return err
	}
	result, err := call.GetResult(d.pulseTimeout)
	if err != nil {
		return err
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

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

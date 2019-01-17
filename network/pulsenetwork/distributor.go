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
	"fmt"
	"sync"
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

	pingRequestTimeout        time.Duration
	randomHostsRequestTimeout time.Duration
	pulseRequestTimeout       time.Duration
	randomNodesCount          int

	pulsarHost     *host.Host
	bootstrapHosts []*host.Host
}

func NewDistributor(conf configuration.PulseDistributor) (core.PulseDistributor, error) {
	bootstrapHosts := make([]*host.Host, 0, len(conf.BootstrapHosts))

	for _, node := range conf.BootstrapHosts {
		bootstrapHost, err := host.NewHost(node)
		if err != nil {
			return nil, errors.Wrap(err, "[ NewDistributor ] failed to create bootstrap node host")
		}
		bootstrapHosts = append(bootstrapHosts, bootstrapHost)
	}

	return &distributor{
		pingRequestTimeout:        time.Duration(conf.PingRequestTimeout) * time.Millisecond,
		randomHostsRequestTimeout: time.Duration(conf.RandomHostsRequestTimeout) * time.Millisecond,
		pulseRequestTimeout:       time.Duration(conf.PulseRequestTimeout) * time.Millisecond,
		randomNodesCount:          conf.RandomNodesCount,

		bootstrapHosts: bootstrapHosts,
	}, nil
}

func (d *distributor) Start(ctx context.Context) error {
	pulsarHost, err := host.NewHost(d.Transport.PublicAddress())
	if err != nil {
		return errors.Wrap(err, "[ NewDistributor ] failed to create pulsar host")
	}
	pulsarHost.NodeID = core.RecordRef{}

	d.pulsarHost = pulsarHost
	return nil
}

func (d *distributor) Distribute(ctx context.Context, pulse *core.Pulse) {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("sendPulseToNetwork failed with panic: %v", r)
		}
	}()

	d.resume(ctx)
	defer d.pause(ctx)

	wg := sync.WaitGroup{}
	wg.Add(len(d.bootstrapHosts))

	for _, bootstrapHost := range d.bootstrapHosts {
		go func(bootstrapHost host.Host) {
			defer wg.Done()

			if bootstrapHost.NodeID.IsEmpty() {
				err := d.pingHost(ctx, &bootstrapHost)
				if err != nil {
					logger.Error("[ Distribute ] failed to ping and fill node id", err)
					return
				}
			}

			hosts, err := d.getRandomHosts(ctx, &bootstrapHost)
			if err != nil {
				logger.Errorf(
					"[ Distribute ] Failed to send pulse to host: %s, error: %s",
					bootstrapHost.String(),
					err.Error(),
				)
			}

			if len(hosts) == 0 {
				err := d.sendPulseToHost(ctx, pulse, &bootstrapHost)
				if err != nil {
					logger.Error(err)
				}
				return
			}

			d.sendPulseToHosts(ctx, pulse, hosts)
		}(*bootstrapHost)
	}

	wg.Wait()
}

func (d *distributor) pingHost(ctx context.Context, host *host.Host) error {
	logger := inslogger.FromContext(ctx)

	builder := packet.NewBuilder(d.pulsarHost)
	pingPacket := builder.Receiver(host).Type(types.Ping).Build()
	pingCall, err := d.Transport.SendRequest(pingPacket)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "[ pingHost ] failed to send ping request")
	}

	logger.Debugf("before ping request")
	result, err := pingCall.GetResult(d.pingRequestTimeout)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "[ pingHost ] failed to get ping result")
	}

	if result.Error != nil {
		logger.Error(result.Error)
		return errors.Wrap(err, "[ pingHost ] ping result returned error")
	}

	host.NodeID = result.Sender.NodeID
	logger.Debugf("ping request is done")

	return nil
}

func (d *distributor) getRandomHosts(ctx context.Context, host *host.Host) ([]host.Host, error) {
	logger := inslogger.FromContext(ctx)

	builder := packet.NewBuilder(d.pulsarHost)
	request := builder.
		Receiver(host).
		Request(&packet.RequestGetRandomHosts{HostsNumber: d.randomNodesCount}).
		Type(types.GetRandomHosts).
		Build()

	logger.Debugf("[ getRandomHosts ] before get random hosts request")
	call, err := d.Transport.SendRequest(request)
	if err != nil {
		logger.Errorf("[ getRandomHosts ] Failed to send request to host: %s, error: %s", host.String(), err)
		return nil, errors.Wrap(err, "[ getRandomHosts ] failed to send getRandomHosts request")
	}

	result, err := call.GetResult(d.randomHostsRequestTimeout)
	if err != nil {
		logger.Errorf("[ getRandomHosts ] Failed to get result from host: %s, error: %s", host.String(), err)
		return nil, errors.Wrap(err, "[ getRandomHosts ] failed to get getRandomHosts result")
	}

	if result.Error != nil {
		logger.Errorf("[ getRandomHosts ] Host %s returned error: %s", host.String(), result.Error.Error())
		return nil, errors.Wrap(result.Error, "[ getRandomHosts ] getRandomHosts result returned error")
	}

	logger.Debugf("[ getRandomHosts ] getRandomHosts request is done")

	body := result.Data.(*packet.ResponseGetRandomHosts)
	if len(body.Error) != 0 {
		logger.Errorf("[ getRandomHosts ] Body result from host %s is error %s", host.String(), body.Error)
		return nil, fmt.Errorf("[ getRandomHosts ] getRandomHosts data returned error: %s", body.Error)
	}

	return body.Hosts, nil
}

func (d *distributor) sendPulseToHosts(ctx context.Context, pulse *core.Pulse, hosts []host.Host) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("Before sending pulse to nodes - %v", hosts)

	wg := sync.WaitGroup{}
	wg.Add(len(hosts))

	for _, pulseReceiver := range hosts {
		go func(host host.Host) {
			defer wg.Done()
			err := d.sendPulseToHost(ctx, pulse, &host)
			if err != nil {
				logger.Errorf(
					"[ sendPulseToHosts ] Failed to send pulse to host: %s, error: %s",
					host.String(),
					err.Error(),
				)
			}
		}(pulseReceiver)
	}

	wg.Wait()
}

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
	result, err := call.GetResult(d.pulseRequestTimeout)
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
	go d.Transport.Stop()
	<-d.Transport.Stopped()
	d.Transport.Close()
}

func (d *distributor) resume(ctx context.Context) {
	inslogger.FromContext(ctx).Info("[ Resume ] Resume distribution, starting transport")
	transport.ListenAndWaitUntilReady(ctx, d.Transport)
}

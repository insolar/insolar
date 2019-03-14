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

package pulsenetwork

import (
	"context"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type distributor struct {
	Transport   transport.Transport `inject:""`
	idGenerator sequence.Generator

	pingRequestTimeout        time.Duration
	randomHostsRequestTimeout time.Duration
	pulseRequestTimeout       time.Duration
	randomNodesCount          int

	pulsarHost     *host.Host
	bootstrapHosts []string
}

// NewDistributor creates a new distributor object of pulses
func NewDistributor(conf configuration.PulseDistributor) (core.PulseDistributor, error) {
	return &distributor{
		idGenerator: sequence.NewGeneratorImpl(),

		pingRequestTimeout:        time.Duration(conf.PingRequestTimeout) * time.Millisecond,
		randomHostsRequestTimeout: time.Duration(conf.RandomHostsRequestTimeout) * time.Millisecond,
		pulseRequestTimeout:       time.Duration(conf.PulseRequestTimeout) * time.Millisecond,
		randomNodesCount:          conf.RandomNodesCount,

		bootstrapHosts: conf.BootstrapHosts,
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

// Distribute starts a fire-and-forget process of pulse distribution to bootstrap hosts
func (d *distributor) Distribute(ctx context.Context, pulse core.Pulse) {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("sendPulseToNetwork failed with panic: %v", r)
		}
	}()

	ctx, span := instracer.StartSpan(ctx, "distributor.Distribute")
	defer span.End()

	bootstrapHosts := make([]*host.Host, 0, len(d.bootstrapHosts))
	for _, node := range d.bootstrapHosts {
		bootstrapHost, err := host.NewHost(node)
		if err != nil {
			logger.Error(err, "[ Distribute ] failed to create bootstrap node host")
			continue
		}
		bootstrapHosts = append(bootstrapHosts, bootstrapHost)
	}

	if len(bootstrapHosts) == 0 {
		logger.Error("[ Distribute ] no bootstrap hosts to distribute")
		return
	}

	if err := d.resume(ctx); err != nil {
		logger.Error("[ Distribute ] resume distribution error: " + err.Error())
		return
	}
	defer d.pause(ctx)

	wg := sync.WaitGroup{}
	wg.Add(len(bootstrapHosts))

	for _, bootstrapHost := range bootstrapHosts {
		go func(ctx context.Context, pulse core.Pulse, bootstrapHost *host.Host) {
			defer wg.Done()

			if bootstrapHost.NodeID.IsEmpty() {
				err := d.pingHost(ctx, bootstrapHost)
				if err != nil {
					logger.Errorf("[ Distribute pulse %d ] Failed to ping and fill node id: %s", pulse.PulseNumber, err)
					return
				}
			}

			err := d.sendPulseToHost(ctx, &pulse, bootstrapHost)
			if err != nil {
				logger.Errorf("[ Distribute pulse %d ] Failed to send pulse: %s", pulse.PulseNumber, err)
				return
			}
			logger.Infof("[ Distribute pulse %d ] Successfully sent pulse to node %s", pulse.PulseNumber, bootstrapHost)
		}(ctx, pulse, bootstrapHost)
	}

	wg.Wait()
}

func (d *distributor) generateID() network.RequestID {
	return network.RequestID(d.idGenerator.Generate())
}

func (d *distributor) pingHost(ctx context.Context, host *host.Host) error {
	logger := inslogger.FromContext(ctx)

	ctx, span := instracer.StartSpan(ctx, "distributor.pingHost")
	defer span.End()
	builder := packet.NewBuilder(d.pulsarHost)
	pingPacket := builder.Receiver(host).Type(types.Ping).RequestID(d.generateID()).Build()
	pingCall, err := d.Transport.SendRequest(ctx, pingPacket)
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

func (d *distributor) sendPulseToHost(ctx context.Context, pulse *core.Pulse, host *host.Host) error {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("sendPulseToHost failed with panic: %v", x)
		}
	}()

	ctx, span := instracer.StartSpan(ctx, "distributor.sendPulseToHosts")
	defer span.End()
	pb := packet.NewBuilder(d.pulsarHost)
	pulseRequest := pb.Receiver(host).Request(&packet.RequestPulse{Pulse: *pulse}).RequestID(d.generateID()).Type(types.Pulse).Build()
	call, err := d.Transport.SendRequest(ctx, pulseRequest)
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
	_, span := instracer.StartSpan(ctx, "distributor.pause")
	defer span.End()
	go d.Transport.Stop()
	<-d.Transport.Stopped()
	d.Transport.Close()
}

func (d *distributor) resume(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("[ Resume ] Resume distribution, starting transport")
	ctx, span := instracer.StartSpan(ctx, "distributor.resume")
	defer span.End()
	return transport.ListenAndWaitUntilReady(ctx, d.Transport)
}

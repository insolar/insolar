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

package pulsenetwork

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/hostnetwork/pool"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
)

type distributor struct {
	Factory     transport.Factory `inject:""`
	transport   transport.StreamTransport
	idGenerator sequence.Generator

	pingRequestTimeout        time.Duration
	randomHostsRequestTimeout time.Duration
	pulseRequestTimeout       time.Duration
	randomNodesCount          int

	publicAddress   string
	pulsarHost      *host.Host
	bootstrapHosts  []string
	futureManager   future.Manager
	responseHandler future.PacketHandler
	pool            pool.ConnectionPool
}

// NewDistributor creates a new distributor object of pulses
func NewDistributor(conf configuration.PulseDistributor) (insolar.PulseDistributor, error) {

	futureManager := future.NewManager()

	result := &distributor{
		idGenerator: sequence.NewGenerator(),

		pingRequestTimeout:        time.Duration(conf.PingRequestTimeout) * time.Millisecond,
		randomHostsRequestTimeout: time.Duration(conf.RandomHostsRequestTimeout) * time.Millisecond,
		pulseRequestTimeout:       time.Duration(conf.PulseRequestTimeout) * time.Millisecond,
		randomNodesCount:          conf.RandomNodesCount,

		bootstrapHosts:  conf.BootstrapHosts,
		futureManager:   futureManager,
		responseHandler: future.NewPacketHandler(futureManager),
	}

	return result, nil
}

func (d *distributor) Init(ctx context.Context) error {
	handler := hostnetwork.NewStreamHandler(func(p *packet.Packet) {}, d.responseHandler)

	var err error
	d.transport, err = d.Factory.CreateStreamTransport(handler)
	if err != nil {
		return errors.Wrap(err, "Failed to create transport")
	}
	d.pool = pool.NewConnectionPool(d.transport)

	return nil
}

func (d *distributor) Start(ctx context.Context) error {

	err := d.transport.Start(ctx)
	if err != nil {
		return err
	}
	d.publicAddress = d.transport.Address()

	pulsarHost, err := host.NewHost(d.publicAddress)
	if err != nil {
		return errors.Wrap(err, "[ NewDistributor ] failed to create pulsar host")
	}
	pulsarHost.NodeID = insolar.Reference{}

	d.pulsarHost = pulsarHost
	return nil
}

// Distribute starts a fire-and-forget process of pulse distribution to bootstrap hosts
func (d *distributor) Distribute(ctx context.Context, pulse insolar.Pulse) {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("sendPulseToNetwork failed with panic: %v", r)
		}
	}()

	traceID := "pulse_" + strconv.FormatUint(uint64(pulse.PulseNumber), 10)
	ctx, logger = inslogger.WithTraceField(ctx, traceID)

	ctx, span := instracer.StartSpan(ctx, "Pulsar.Distribute")
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(pulse.PulseNumber)),
	)
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
		go func(ctx context.Context, pulse insolar.Pulse, bootstrapHost *host.Host) {
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

func (d *distributor) generateID() types.RequestID {
	return types.RequestID(d.idGenerator.Generate())
}

func (d *distributor) pingHost(ctx context.Context, host *host.Host) error {
	logger := inslogger.FromContext(ctx)

	ctx, span := instracer.StartSpan(ctx, "distributor.pingHost")
	defer span.End()

	pingPacket := packet.NewPacket(d.pulsarHost, host, types.Ping, uint64(d.generateID()))
	pingPacket.SetRequest(&packet.Ping{})
	pingCall, err := d.sendRequestToHost(ctx, pingPacket, host)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "[ pingHost ] failed to send ping request")
	}

	logger.Debugf("before ping request")
	result, err := pingCall.WaitResponse(d.pingRequestTimeout)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "[ pingHost ] failed to get ping result")
	}

	host.NodeID = result.GetSender()
	logger.Debugf("ping request is done")

	return nil
}

func (d *distributor) sendPulseToHost(ctx context.Context, p *insolar.Pulse, host *host.Host) error {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("sendPulseToHost failed with panic: %v", x)
		}
	}()

	ctx, span := instracer.StartSpan(ctx, "distributor.sendPulseToHosts")
	defer span.End()
	pulseRequest := packet.NewPacket(d.pulsarHost, host, types.Pulse, uint64(d.generateID()))
	request := &packet.PulseRequest{
		TraceSpanData: instracer.MustSerialize(ctx),
		Pulse:         pulse.ToProto(p),
	}
	pulseRequest.SetRequest(request)
	call, err := d.sendRequestToHost(ctx, pulseRequest, host)
	if err != nil {
		return err
	}
	_, err = call.WaitResponse(d.pulseRequestTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (d *distributor) pause(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Pause ] Pause distribution, stopping transport")
	d.pool.Reset()
	err := d.transport.Stop(ctx)
	if err != nil {
		logger.Errorf("Failed to stop network: %s", err.Error())
	}
}

func (d *distributor) resume(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("[ Resume ] Resume distribution, starting transport")
	return d.transport.Start(ctx)
}

func (d *distributor) sendRequestToHost(ctx context.Context, packet *packet.Packet, receiver *host.Host) (network.Future, error) {
	inslogger.FromContext(ctx).Debugf("Send %s request to %s with RequestID = %d",
		packet.GetType(), receiver.String(), packet.GetRequestID())

	f := d.futureManager.Create(packet)
	err := hostnetwork.SendPacket(ctx, d.pool, packet)
	if err != nil {
		f.Cancel()
		return nil, errors.Wrap(err, "Failed to send transport packet")
	}
	metrics.NetworkPacketSentTotal.WithLabelValues(packet.GetType().String()).Inc()
	return f, nil
}

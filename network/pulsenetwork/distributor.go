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
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar/pulse"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
)

type distributor struct {
	Factory     transport.Factory `inject:""`
	transport   transport.StreamTransport
	idGenerator sequence.Generator

	pulseRequestTimeout time.Duration

	publicAddress   string
	pulsarHost      *host.Host
	bootstrapHosts  []string
	futureManager   future.Manager
	responseHandler future.PacketHandler
}

// NewDistributor creates a new distributor object of pulses
func NewDistributor(conf configuration.PulseDistributor) (insolar.PulseDistributor, error) {

	futureManager := future.NewManager()

	result := &distributor{
		idGenerator: sequence.NewGenerator(),

		pulseRequestTimeout: time.Duration(conf.PulseRequestTimeout) * time.Millisecond,

		bootstrapHosts:  conf.BootstrapHosts,
		futureManager:   futureManager,
		responseHandler: future.NewPacketHandler(futureManager),
	}

	return result, nil
}

func (d *distributor) Init(ctx context.Context) error {
	handler := hostnetwork.NewStreamHandler(func(context.Context, *packet.ReceivedPacket) {}, d.responseHandler)

	var err error
	d.transport, err = d.Factory.CreateStreamTransport(handler)
	if err != nil {
		return errors.Wrap(err, "Failed to create transport")
	}

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
	pulsarHost.NodeID = *insolar.NewEmptyReference()

	d.pulsarHost = pulsarHost
	return nil
}

func (d *distributor) Stop(ctx context.Context) error {
	return d.transport.Stop(ctx)
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

	wg := sync.WaitGroup{}
	wg.Add(len(d.bootstrapHosts))

	distributed := 0
	for _, nodeAddr := range d.bootstrapHosts {
		go func(ctx context.Context, pulse insolar.Pulse, nodeAddr string) {
			defer wg.Done()

			addr, err := net.ResolveTCPAddr("tcp", nodeAddr)
			if err != nil {
				logger.Warnf("failed to resolve bootstrap node address %s, %s", nodeAddr, err.Error())
				return
			}

			err = d.sendPulseToHost(ctx, &pulse, addr)
			if err != nil {
				stats.Record(ctx, statSendPulseErrorsCount.M(1))
				logger.Warnf("Failed to send pulse %d to host: %s %s", pulse.PulseNumber, nodeAddr, err)
				return
			}

			distributed++
			logger.Infof("Successfully sent pulse %d to node %s", pulse.PulseNumber, nodeAddr)
		}(ctx, pulse, nodeAddr)
	}
	wg.Wait()

	if distributed == 0 {
		logger.Warn("No bootstrap hosts to distribute")
	} else {
		logger.Infof("Pulse distributed to %d hosts", distributed)
	}

}

func (d *distributor) generateID() types.RequestID {
	return types.RequestID(d.idGenerator.Generate())
}

func (d *distributor) sendPulseToHost(ctx context.Context, p *insolar.Pulse, host *net.TCPAddr) error {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("sendPulseToHost failed with panic: %v", x)
		}
	}()

	ctx, span := instracer.StartSpan(ctx, "distributor.sendPulseToHosts")
	defer span.End()

	pulseRequest := NewPulsePacketWithTrace(ctx, p, d.pulsarHost, host, uint64(d.generateID()))

	err := d.sendRequestToHost(ctx, pulseRequest, host)
	if err != nil {
		return err
	}
	return nil
}

func (d *distributor) sendRequestToHost(ctx context.Context, p *packet.Packet, receiver *net.TCPAddr) error {
	inslogger.FromContext(ctx).Debugf("Send %s request to %s with RequestID = %d",
		p.GetType(), receiver.String(), p.GetRequestID())

	data, err := packet.SerializePacket(p)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize packet")
	}

	conn, err := d.transport.Dial(ctx, receiver.String())
	if err != nil {
		return errors.Wrap(err, "Unable to connect")
	}
	defer conn.Close() //nolint

	n, err := conn.Write(data)
	if err != nil {
		return errors.Wrap(err, "[ SendPacket ] Failed to write data")
	}

	metrics.NetworkSentSize.Observe(float64(n))
	metrics.NetworkPacketSentTotal.WithLabelValues(p.GetType().String()).Inc()

	return nil
}

func NewPulsePacket(p *insolar.Pulse, pulsarHost *host.Host, to *net.TCPAddr, id uint64) *packet.Packet {
	rcv := host.Host{Address: &host.Address{net.UDPAddr{IP: to.IP, Port: to.Port, Zone: to.Zone}}} //nolint
	pulseRequest := packet.NewPacket(pulsarHost, &rcv, types.Pulse, id)
	request := &packet.PulseRequest{
		Pulse: pulse.ToProto(p),
	}
	pulseRequest.SetRequest(request)
	return pulseRequest
}

func NewPulsePacketWithTrace(ctx context.Context, p *insolar.Pulse, pulsarHost *host.Host, to *net.TCPAddr, id uint64) *packet.Packet {
	pulsePacket := NewPulsePacket(p, pulsarHost, to, id)
	pulsePacket.TraceSpanData = instracer.MustSerialize(ctx)
	return pulsePacket
}

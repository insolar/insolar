// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsenetwork

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

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

	pulseCtx := inslogger.SetLogger(context.Background(), logger)

	traceID := strconv.FormatUint(uint64(pulse.PulseNumber), 10) + "_pulse"
	pulseCtx, logger = inslogger.WithTraceField(pulseCtx, traceID)

	pulseCtx, span := instracer.StartSpan(pulseCtx, "Pulsar.Distribute")
	span.LogFields(
		log.Int64("pulse.PulseNumber", int64(pulse.PulseNumber)),
	)
	defer span.Finish()

	wg := sync.WaitGroup{}
	wg.Add(len(d.bootstrapHosts))

	distributed := int32(0)
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

			atomic.AddInt32(&distributed, 1)
			logger.Infof("Successfully sent pulse %d to node %s", pulse.PulseNumber, nodeAddr)
		}(pulseCtx, pulse, nodeAddr)
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
	defer span.Finish()

	pulseRequest := NewPulsePacketWithTrace(ctx, p, d.pulsarHost, host, uint64(d.generateID()))

	err := d.sendRequestToHost(ctx, pulseRequest, host)
	if err != nil {
		return err
	}
	return nil
}

func (d *distributor) sendRequestToHost(ctx context.Context, p *packet.Packet, receiver fmt.Stringer) error {
	rcv := receiver.String()
	inslogger.FromContext(ctx).Debugf("Send %s request to %s with RequestID = %d",
		p.GetType(), rcv, p.GetRequestID())

	data, err := packet.SerializePacket(p)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize packet")
	}

	conn, err := d.transport.Dial(ctx, rcv)
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

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

package hostnetwork

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/hostnetwork/pool"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
)

// NewHostNetwork constructor creates new NewHostNetwork component
func NewHostNetwork(nodeRef string) (network.HostNetwork, error) {

	id, err := insolar.NewReferenceFromBase58(nodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "invalid nodeRef")
	}

	futureManager := future.NewManager()

	result := &hostNetwork{
		handlers:          make(map[types.PacketType]network.RequestHandler),
		sequenceGenerator: sequence.NewGenerator(),
		nodeID:            *id,
		futureManager:     futureManager,
		responseHandler:   future.NewPacketHandler(futureManager),
	}

	return result, nil
}

type hostNetwork struct {
	Resolver network.RoutingTable `inject:""`
	Factory  transport.Factory    `inject:""`

	nodeID            insolar.Reference
	started           uint32
	transport         transport.StreamTransport
	sequenceGenerator sequence.Generator
	handlers          map[types.PacketType]network.RequestHandler
	futureManager     future.Manager
	responseHandler   future.PacketHandler
	pool              pool.ConnectionPool

	muOrigin sync.RWMutex
	origin   *host.Host
}

func (hn *hostNetwork) Init(ctx context.Context) error {

	handler := NewStreamHandler(hn.handleRequest, hn.responseHandler)

	var err error
	hn.transport, err = hn.Factory.CreateStreamTransport(handler)
	if err != nil {
		return errors.Wrap(err, "Failed to create stream transport")
	}

	hn.pool = pool.NewConnectionPool(hn.transport)
	return err
}

// Start listening to network requests, should be started in goroutine.
func (hn *hostNetwork) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&hn.started, 0, 1) {
		inslogger.FromContext(ctx).Warn("HostNetwork component already started")
		return nil
	}

	hn.muOrigin.Lock()
	defer hn.muOrigin.Unlock()

	if err := hn.transport.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start stream transport")
	}

	h, err := host.NewHostN(hn.transport.Address(), hn.nodeID)
	if err != nil {
		return errors.Wrap(err, "failed to create host")
	}

	hn.origin = h

	return nil
}

// Stop listening to network requests.
func (hn *hostNetwork) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&hn.started, 1, 0) {
		err := hn.transport.Stop(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to stop transport.")
		}
	}
	return nil
}

func (hn *hostNetwork) buildRequest(ctx context.Context, packetType types.PacketType,
	requestData interface{}, receiver *host.Host) *packet.PacketBackend {
	result := &packet.PacketBackend{
		Sender:   hn.getOrigin(),
		Receiver: receiver,
		// TODO: replace in protobuf with our type
		RequestID: uint64(hn.sequenceGenerator.Generate()),
		TraceID:   inslogger.TraceID(ctx),
		Type:      uint32(packetType),
	}
	result.SetRequest(requestData)
	return result
}

// PublicAddress returns public address that can be published for all nodes.
func (hn *hostNetwork) PublicAddress() string {
	return hn.getOrigin().Address.String()
}

func (hn *hostNetwork) handleRequest(p *packet.PacketBackend) {
	ctx, logger := inslogger.WithTraceField(context.Background(), p.TraceID)
	logger.Debugf("Got %s request from host %s; RequestID: %d", p.Type, p.Sender, p.RequestID)
	handler, exist := hn.handlers[p.GetType()]
	if !exist {
		logger.Errorf("No handler set for packet type %s from node %s", p.Type, p.Sender.NodeID)
		return
	}
	ctx, span := instracer.StartSpan(ctx, "hostTransport.processMessage")
	span.AddAttributes(
		trace.StringAttribute("msg receiver", p.Receiver.Address.String()),
		trace.StringAttribute("msg trace", p.TraceID),
		trace.StringAttribute("msg type", p.GetType().String()),
	)
	defer span.End()
	response, err := handler(ctx, p)
	if err != nil {
		logger.Errorf("Error handling request %s from node %s: %s", p.Type, p.Sender.NodeID, err)
		return
	}

	responsePacket := response.(*packet.PacketBackend)
	responsePacket.RequestID = p.RequestID
	err = SendPacket(ctx, hn.pool, responsePacket)
	if err != nil {
		logger.Errorf("Failed to send response: %s", err.Error())
	}
}

// SendRequestToHost send request packet to a remote node.
func (hn *hostNetwork) SendRequestToHost(ctx context.Context, packetType types.PacketType,
	requestData interface{}, receiver *host.Host) (network.Future, error) {

	if atomic.LoadUint32(&hn.started) == 0 {
		return nil, errors.New("host network is not started")
	}

	p := hn.buildRequest(ctx, packetType, requestData, receiver)

	inslogger.FromContext(ctx).Debugf("Send %s request to %s with RequestID = %d", p.Type, p.Receiver, p.RequestID)

	f := hn.futureManager.Create(p)
	err := SendPacket(ctx, hn.pool, p)
	if err != nil {
		f.Cancel()
		return nil, errors.Wrap(err, "Failed to send transport packet")
	}
	metrics.NetworkPacketSentTotal.WithLabelValues(p.GetType().String()).Inc()
	return f, nil
}

// RegisterPacketHandler register a handler function to process incoming request packets of a specific type.
func (hn *hostNetwork) RegisterPacketHandler(t types.PacketType, handler network.RequestHandler) {
	_, exists := hn.handlers[t]
	if exists {
		log.Warnf("Multiple handlers for packet type %s are not supported! New handler will replace the old one!", t)
	}
	hn.handlers[t] = handler
}

// BuildResponse create response to an incoming request with Data set to responseData.
func (hn *hostNetwork) BuildResponse(ctx context.Context, request network.Packet, responseData interface{}) network.Packet {
	result := &packet.PacketBackend{
		Sender:   hn.getOrigin(),
		Receiver: request.GetSenderHost(),
		// TODO: replace in protobuf with our type
		RequestID: uint64(request.GetRequestID()),
		TraceID:   inslogger.TraceID(ctx),
	}
	result.SetResponse(responseData)
	return result
}

// SendRequest send request to a remote node.
func (hn *hostNetwork) SendRequest(ctx context.Context, packetType types.PacketType,
	requestData interface{}, receiver insolar.Reference) (network.Future, error) {

	h, err := hn.Resolver.Resolve(receiver)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving NodeID -> Address")
	}
	return hn.SendRequestToHost(ctx, packetType, requestData, h)
}

// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
func (hn *hostNetwork) RegisterRequestHandler(t types.PacketType, handler network.RequestHandler) {
	f := func(ctx context.Context, request network.Packet) (network.Packet, error) {
		hn.Resolver.AddToKnownHosts(request.GetSenderHost())
		return handler(ctx, request)
	}
	hn.RegisterPacketHandler(t, f)
}

func (hn *hostNetwork) getOrigin() *host.Host {
	hn.muOrigin.RLock()
	defer hn.muOrigin.RUnlock()

	return hn.origin
}

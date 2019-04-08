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
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
)

func NewHostNetwork(conf configuration.Configuration, nodeRef string) (network.HostNetwork, error) {
	tp, publicAddress, err := transport.NewTransport(conf.Host.Transport)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transport")
	}
	id, err := insolar.NewReferenceFromBase58(nodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "invalid nodeRef")
	}

	origin, err := host.NewHostN(publicAddress, *id)
	if err != nil {
		return nil, errors.Wrap(err, "error getting origin")
	}
	result := &hostNetwork{handlers: make(map[types.PacketType]network.RequestHandler)}
	result.sequenceGenerator = sequence.NewGeneratorImpl()
	result.transport = tp
	result.origin = origin
	result.messageProcessor = result.processMessage
	return result, nil
}

type hostNetwork struct {
	Resolver network.RoutingTable `inject:""`

	started           uint32
	transport         transport.Transport
	origin            *host.Host
	messageProcessor  func(msg *packet.Packet)
	sequenceGenerator sequence.Generator
	handlers          map[types.PacketType]network.RequestHandler
}

// Start start listening to network requests, should be started in goroutine.
func (hn *hostNetwork) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&hn.started, 0, 1) {
		return errors.New("Failed to start transport: double listen initiated")
	}
	if err := hn.transport.Start(ctx); err != nil {
		return errors.Wrap(err, "Failed to start transport: listen syscall failed")
	}

	go hn.listen(ctx)
	return nil
}

func (hn *hostNetwork) listen(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	for {
		select {
		case msg := <-hn.transport.Packets():
			if msg == nil {
				logger.Error("HostNetwork receiving channel is closed")
				break
			}
			if msg.Error != nil {
				logger.Warnf("Received error response: %s", msg.Error.Error())
			}
			go hn.messageProcessor(msg)
		case <-hn.transport.Stopped():
			return
		}
	}
}

// Disconnect stop listening to network requests.
func (hn *hostNetwork) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&hn.started, 1, 0) {
		go hn.transport.Stop()
		<-hn.transport.Stopped()
		hn.transport.Close()
	}
	return nil
}

func (hn *hostNetwork) buildRequest(ctx context.Context, request network.Request, receiver *host.Host) *packet.Packet {
	return packet.NewBuilder(hn.origin).Receiver(receiver).Type(request.GetType()).RequestID(request.GetRequestID()).
		Request(request.GetData()).TraceID(inslogger.TraceID(ctx)).Build()
}

// PublicAddress returns public address that can be published for all nodes.
func (hn *hostNetwork) PublicAddress() string {
	return hn.origin.Address.String()
}

// GetNodeID get current node ID.
func (hn *hostNetwork) GetNodeID() insolar.Reference {
	return hn.origin.NodeID
}

// NewRequestBuilder create packet Builder for an outgoing request with sender set to current node.
func (hn *hostNetwork) NewRequestBuilder() network.RequestBuilder {
	return &Builder{sender: hn.origin, id: network.RequestID(hn.sequenceGenerator.Generate())}
}

func (hn *hostNetwork) processMessage(msg *packet.Packet) {
	ctx, logger := inslogger.WithTraceField(context.Background(), msg.TraceID)
	logger.Debugf("Got %s request from host %s; RequestID: %d", msg.Type.String(), msg.Sender.String(), msg.RequestID)
	handler, exist := hn.handlers[msg.Type]
	if !exist {
		logger.Errorf("No handler set for packet type %s from node %s",
			msg.Type.String(), msg.Sender.NodeID.String())
		return
	}
	ctx, span := instracer.StartSpan(ctx, "hostTransport.processMessage")
	span.AddAttributes(
		trace.StringAttribute("msg receiver", msg.Receiver.Address.String()),
		trace.StringAttribute("msg trace", msg.TraceID),
		trace.StringAttribute("msg type", msg.Type.String()),
	)
	defer span.End()
	response, err := handler(ctx, msg)
	if err != nil {
		logger.Errorf("Error handling request %s from node %s: %s",
			msg.Type.String(), msg.Sender.NodeID.String(), err)
		return
	}
	err = hn.transport.SendResponse(ctx, msg.RequestID, response.(*packet.Packet))
	if err != nil {
		logger.Error(err)
	}
}

// SendRequestPacket send request packet to a remote node.
func (hn *hostNetwork) SendRequestPacket(ctx context.Context, request network.Request, receiver *host.Host) (network.Future, error) {
	inslogger.FromContext(ctx).Debugf("Send %s request to host %s", request.GetType().String(), receiver.String())
	f, err := hn.transport.SendRequest(ctx, hn.buildRequest(ctx, request, receiver))
	if err != nil {
		return nil, err
	}
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
func (hn *hostNetwork) BuildResponse(ctx context.Context, request network.Request, responseData interface{}) network.Response {
	sender := request.GetSenderHost()
	p := packet.NewBuilder(hn.origin).Type(request.GetType()).Receiver(sender).RequestID(request.GetRequestID()).
		Response(responseData).TraceID(inslogger.TraceID(ctx)).Build()
	return p
}

// SendRequest send request to a remote node.
func (hn *hostNetwork) SendRequest(ctx context.Context, request network.Request, receiver insolar.Reference) (network.Future, error) {
	h, err := hn.Resolver.Resolve(receiver)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving NodeID -> Address")
	}
	return hn.SendRequestPacket(ctx, request, h)
}

// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
func (hn *hostNetwork) RegisterRequestHandler(t types.PacketType, handler network.RequestHandler) {
	f := func(ctx context.Context, request network.Request) (network.Response, error) {
		hn.Resolver.AddToKnownHosts(request.GetSenderHost())
		return handler(ctx, request)
	}
	hn.RegisterPacketHandler(t, f)
}

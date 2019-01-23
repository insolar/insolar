/*
 *    Copyright 2019 Insolar
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

package hostnetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"
)

type hostTransport struct {
	transportBase
	handlers map[types.PacketType]network.RequestHandler
}

type packetWrapper packet.Packet

func (p *packetWrapper) GetSender() core.RecordRef {
	return p.Sender.NodeID
}

func (p *packetWrapper) GetSenderHost() *host.Host {
	return p.Sender
}

func (p *packetWrapper) GetType() types.PacketType {
	return p.Type
}

func (p *packetWrapper) GetData() interface{} {
	return p.Data
}

func (p *packetWrapper) GetRequestID() network.RequestID {
	return p.RequestID
}

type future struct {
	transport.Future
}

// Response get channel that receives response to sent request
func (f future) Response() <-chan network.Response {
	in := transport.Future(f).Result()
	out := make(chan network.Response, cap(in))
	go func(in <-chan *packet.Packet, out chan<- network.Response) {
		for packet := range in {
			out <- (*packetWrapper)(packet)
		}
		close(out)
	}(in, out)
	return out
}

// GetResponse get response to sent request with `duration` timeout
func (f future) GetResponse(duration time.Duration) (network.Response, error) {
	result, err := f.GetResult(duration)
	if err != nil {
		return nil, err
	}
	return (*packetWrapper)(result), nil
}

// GetRequest get initiating request.
func (f future) GetRequest() network.Request {
	request := transport.Future(f).Request()
	return (*packetWrapper)(request)
}

func (h *hostTransport) processMessage(msg *packet.Packet) {
	ctx, logger := inslogger.WithTraceField(context.Background(), msg.TraceID)
	logger.Debugf("Got %s request from host %s; RequestID: %d", msg.Type.String(), msg.Sender.String(), msg.RequestID)
	handler, exist := h.handlers[msg.Type]
	if !exist {
		logger.Errorf("No handler set for packet type %s from node %s",
			msg.Type.String(), msg.Sender.NodeID.String())
		return
	}
	response, err := handler(ctx, (*packetWrapper)(msg))
	if err != nil {
		logger.Errorf("Error handling request %s from node %s: %s",
			msg.Type.String(), msg.Sender.NodeID.String(), err)
		return
	}
	r := response.(*packetWrapper)
	err = h.transport.SendResponse(ctx, msg.RequestID, (*packet.Packet)(r))
	if err != nil {
		logger.Error(err)
	}
}

// SendRequestPacket send request packet to a remote node.
func (h *hostTransport) SendRequestPacket(ctx context.Context, request network.Request, receiver *host.Host) (network.Future, error) {
	inslogger.FromContext(ctx).Debugf("Send %s request to host %s", request.GetType().String(), receiver.String())
	f, err := h.transport.SendRequest(ctx, h.buildRequest(ctx, request, receiver))
	if err != nil {
		return nil, err
	}
	return future{Future: f}, nil
}

// RegisterPacketHandler register a handler function to process incoming request packets of a specific type.
func (h *hostTransport) RegisterPacketHandler(t types.PacketType, handler network.RequestHandler) {
	_, exists := h.handlers[t]
	if exists {
		panic(fmt.Sprintf("multiple handlers for packet type %s are not supported!", t.String()))
	}
	h.handlers[t] = handler
}

// BuildResponse create response to an incoming request with Data set to responseData.
func (h *hostTransport) BuildResponse(ctx context.Context, request network.Request, responseData interface{}) network.Response {
	sender := request.(*packetWrapper).Sender
	p := packet.NewBuilder(h.origin).Type(request.GetType()).Receiver(sender).RequestID(request.GetRequestID()).
		Response(responseData).TraceID(inslogger.TraceID(ctx)).Build()
	return (*packetWrapper)(p)
}

func NewInternalTransport(conf configuration.Configuration, nodeRef string) (network.InternalTransport, error) {
	tp, err := transport.NewTransport(conf.Host.Transport, relay.NewProxy())
	if err != nil {
		return nil, errors.Wrap(err, "error creating transport")
	}
	origin, err := getOrigin(tp, nodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "error getting origin")
	}
	result := &hostTransport{handlers: make(map[types.PacketType]network.RequestHandler)}
	result.sequenceGenerator = sequence.NewGeneratorImpl()
	result.transport = tp
	result.origin = origin
	result.messageProcessor = result.processMessage
	return result, nil
}

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

package hostnetwork

import (
	"fmt"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type HostNetwork struct {
	transport transport.Transport
	origin    *host.Host
	handlers  map[types.PacketType]network.RequestHandler
}

type builder struct {
	p *packet.Packet
}

func newBuilder(origin *host.Host) *builder {
	data := &packet.Packet{Sender: origin}
	return &builder{p: data}
}

func (b *builder) Receiver(ref core.RecordRef) network.RequestBuilder {
	b.p.Receiver = &host.Host{NodeID: ref}
	return b
}

func (b *builder) Type(packetType types.PacketType) network.RequestBuilder {
	b.p.Type = packetType
	return b
}

func (b *builder) Data(data interface{}) network.RequestBuilder {
	b.p.Data = data
	b.p.IsResponse = false
	return b
}

func (b *builder) Build() network.Request {
	return (*racket)(b.p)
}

type racket packet.Packet

func (p *racket) GetSender() core.RecordRef {
	return p.Sender.NodeID
}

func (p *racket) GetReceiver() core.RecordRef {
	return p.Receiver.NodeID
}

func (p *racket) GetType() types.PacketType {
	return p.Type
}

func (p *racket) GetData() interface{} {
	return p.Data
}

type future struct {
	transport.Future
}

func (f future) Response() <-chan network.Response {
	in := transport.Future(f).Result()
	out := make(chan network.Response, cap(in))
	go func() {
		for packet := range in {
			out <- (*racket)(packet)
		}
	}()
	return out
}

func (f future) GetResponse(duration time.Duration) (network.Response, error) {
	select {
	case result, ok := <-f.Result():
		if !ok {
			return nil, transport.ErrChannelClosed
		}
		return (*racket)(result), nil
	case <-time.After(duration):
		f.Cancel()
		return nil, transport.ErrTimeout
	}
}

func (f future) GetRequest() network.Request {
	request := transport.Future(f).Request()
	return (*racket)(request)
}

// Listen start listening to network requests, should be started in goroutine.
func (h *HostNetwork) Listen() error {
	go h.listen()
	return h.transport.Start()
}

func (h *HostNetwork) listen() {
	for {
		select {
		case msg := <-h.transport.Packets():
			if msg == nil {
				log.Error("HostNetwork receiving channel is closed, disconnecting")
				h.Disconnect()
				break
			}
			if msg.Error != nil {
				log.Warn("Received error response")
			}
			h.processMessage(msg)
		case <-h.transport.Stopped():
			h.transport.Close()
			return
		}
	}
}

func (h *HostNetwork) processMessage(msg *packet.Packet) {
	handler, exist := h.handlers[msg.Type]
	if !exist {
		log.Errorf("No handler set for packet type %s from node %s",
			msg.Type.String(), msg.Sender.NodeID.String())
		return
	}
	response, err := handler((*racket)(msg))
	if err != nil {
		log.Errorf("Error handling request %s from node %s: %s",
			msg.Type.String(), msg.Sender.NodeID.String(), err)
		return
	}
	r := response.(*racket)
	h.transport.SendResponse(msg.RequestID, (*packet.Packet)(r))
}

// Disconnect stop listening to network requests.
func (h *HostNetwork) Disconnect() error {
	h.transport.Stop()
	return nil
}

// PublicAddress returns public address that can be published for all nodes.
func (h *HostNetwork) PublicAddress() string {
	return h.origin.Address.String()
}

// SendRequest send request to a remote node.
func (h *HostNetwork) SendRequest(request network.Request) (network.Future, error) {
	f, err := h.transport.SendRequest(requestToPacket(request))
	if err != nil {
		return nil, err
	}
	// TODO: resolve NodeID -> Address
	return future{Future: f}, nil
}

// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
func (h *HostNetwork) RegisterRequestHandler(t types.PacketType, handler network.RequestHandler) {
	_, exists := h.handlers[t]
	if exists {
		panic(fmt.Sprintf("multiple handlers for packet type %s are not supported!", t.String()))
	}
	h.handlers[t] = handler
}

// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
func (h *HostNetwork) NewRequestBuilder() network.RequestBuilder {
	return newBuilder(h.origin)
}

// BuildResponse create response to an incoming request with Data set to responseData
func (h *HostNetwork) BuildResponse(request network.Request, responseData interface{}) network.Response {
	sender := requestToPacket(request).Sender
	p := packet.NewBuilder(h.origin).Type(request.GetType()).
		Receiver(sender).Response(responseData).Build()
	return (*racket)(p)
}

func requestToPacket(request network.Request) *packet.Packet {
	return (*packet.Packet)(request.(*racket))
}

func NewHostNetwork() network.HostNetwork {
	return nil
}

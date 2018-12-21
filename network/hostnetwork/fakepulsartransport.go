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
	"context"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/network/transport/relay"
)

type msgHandler func(msg *packet.Packet)

type FakePulsarTransport interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
	SendRequest(request network.Request, receiver *host.Host) error
	NewBuilder() network.RequestBuilder
}

func NewFakePulsarTransport(handler msgHandler, cfg configuration.Configuration) FakePulsarTransport {
	if handler == nil {
		panic("[ NewFakePulsarTransport ] message handler is nil")
	}
	return &FPTransport{
		handler: handler,
		Cfg:     cfg,
	}
}

type FPTransport struct {
	transportBase
	Cfg configuration.Configuration

	tr      transport.Transport
	origin  *host.Host
	handler msgHandler
}

func (fpt *FPTransport) Start(ctx context.Context) {
	conf := configuration.Transport{}
	conf.Address = fpt.Cfg.Service.FakePulsarAddress
	conf.Protocol = "PURE_UDP"
	conf.BehindNAT = false

	var err error
	fpt.tr, err = transport.NewTransport(conf, relay.NewProxy())
	if err != nil {
		log.Error("[ FakePulsarTransport ] failed to create a transport")
	}
}

func (fpt *FPTransport) Stop(ctx context.Context) {
	fpt.tr.Stop()
	fpt.tr.Close()
}

func (fpt *FPTransport) SendRequest(request network.Request, receiver *host.Host) error {
	p := fpt.buildRequest(request, receiver)
	return fpt.tr.SendPacket(p)
}

func (ftp *FPTransport) NewBuilder() network.RequestBuilder {
	return ftp.NewRequestBuilder()
}

func (fpt *FPTransport) processMessage(ctx context.Context, msg *packet.Packet) {
	switch msg.Type {
	case types.FakePulsarResponse:
	case types.FakePulsarRequest:
		fpt.handler(msg)
	default:
		log.Debugf("[ FPTransport.processMessage ] received a wrong type packet")
	}
}

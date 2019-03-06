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

package transport

import (
	"context"
	"net"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/connection"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/network/transport/resolver"
	"github.com/insolar/insolar/network/utils"

	"github.com/pkg/errors"
)

// Transport is an interface for network transport.
type Transport interface {
	// SendRequest sends packet to destination. Sequence number is generated automatically.
	SendRequest(context.Context, *packet.Packet) (Future, error)

	// SendResponse sends response packet for request with passed request id.
	SendResponse(context.Context, network.RequestID, *packet.Packet) error

	// SendPacket low-level send packet without requestId and without spawning a waiting future
	SendPacket(ctx context.Context, p *packet.Packet) error

	// Listen starts thread to listen incoming packets.
	Listen(ctx context.Context, started chan struct{}) error

	// Stop gracefully stops listening.
	Stop()

	// Close disposing all transport underlying structures after stopped are called.
	Close()

	// Packets returns channel to listen incoming packets.
	Packets() <-chan *packet.Packet

	// Stopped returns signal channel to support graceful shutdown.
	Stopped() <-chan bool

	// PublicAddress returns PublicAddress
	PublicAddress() string
}

// NewTransport creates new Transport with particular configuration
func NewTransport(cfg configuration.Transport, proxy relay.Proxy) (Transport, error) {
	// TODO: let each transport creates connection in their constructor
	conn, publicAddress, err := NewConnection(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewTransport ] Failed to create connection.")
	}

	switch cfg.Protocol {
	case "TCP":
		// TODO: little hack: It's better to change interface for NewConnection
		utils.CloseVerbose(conn)

		return newTCPTransport(conn.LocalAddr().String(), proxy, publicAddress)
	case "PURE_UDP":
		// TODO: not little hack: @AndreyBronin rewrite all this mess, please!
		localAddress := conn.LocalAddr().String()
		utils.CloseVerbose(conn)

		return newUDPTransport(localAddress, proxy, publicAddress)
	case "QUIC":
		return newQuicTransport(conn, proxy, publicAddress)
	default:
		utils.CloseVerbose(conn)
		return nil, errors.New("invalid transport configuration")
	}
}

// NewConnection creates new Connection from configuration and returns connection and public address
func NewConnection(cfg configuration.Transport) (net.PacketConn, string, error) {
	conn, err := connection.NewConnectionFactory().Create(cfg.Address)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ NewConnection ] Failed to create connection")
	}
	resolver, err := createResolver(cfg)
	if err != nil {
		utils.CloseVerbose(conn)
		return nil, "", errors.Wrap(err, "[ NewConnection ] Failed to create resolver")
	}
	publicAddress, err := resolver.Resolve(conn)
	if err != nil {
		utils.CloseVerbose(conn)
		return nil, "", errors.Wrap(err, "[ NewConnection ] Failed to resolve public address")
	}
	return conn, publicAddress, nil
}

func createResolver(cfg configuration.Transport) (resolver.PublicAddressResolver, error) {
	if cfg.BehindNAT && cfg.FixedPublicAddress != "" {
		return nil, errors.New("BehindNAT and fixedPublicAddress cannot be set both")
	}

	if cfg.BehindNAT {
		return resolver.NewStunResolver(""), nil
	} else if cfg.FixedPublicAddress != "" {
		return resolver.NewFixedAddressResolver(cfg.FixedPublicAddress), nil
	}
	return resolver.NewExactResolver(), nil
}

func ListenAndWaitUntilReady(ctx context.Context, transport Transport) {
	started := make(chan struct{}, 1)
	go func(ctx context.Context, t Transport, started chan struct{}) {
		err := t.Listen(ctx, started)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
	}(ctx, transport, started)
	<-started
}

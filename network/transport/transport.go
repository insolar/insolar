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

package transport

import (
	"context"
	"net"

	"github.com/insolar/insolar/configuration"
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
	SendRequest(*packet.Packet) (Future, error)

	// SendResponse sends response packet for request with passed request id.
	SendResponse(packet.RequestID, *packet.Packet) error

	// SendPacket low-level send packet without requestId and without spawning a waiting future
	SendPacket(p *packet.Packet) error

	// Listen starts thread to listen incoming packets.
	Listen(ctx context.Context) error

	// Stop gracefully stops listening.
	Stop()

	// Close disposing all transport underlying structures after stop are called.
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
	case "UTP":
		return newUTPTransport(conn, proxy, publicAddress)
	case "PURE_UDP":
		return newUDPTransport(conn, proxy, publicAddress)
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
	publicAddress, err := createResolver(cfg.BehindNAT).Resolve(conn)
	if err != nil {
		utils.CloseVerbose(conn)
		return nil, "", errors.Wrap(err, "[ NewConnection ] Failed to create resolver")
	}
	return conn, publicAddress, nil
}

func createResolver(stun bool) resolver.PublicAddressResolver {
	if stun {
		return resolver.NewStunResolver("")
	}
	return resolver.NewExactResolver()
}

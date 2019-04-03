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

package transport

import (
	"context"
	"net"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/resolver"

	"github.com/pkg/errors"
)

// Transport is an interface for network transport.
type Transport interface {
	// SendRequest sends packet to destination. Sequence number is generated automatically.
	SendRequest(context.Context, *packet.Packet) (future.Future, error)

	// SendResponse sends response packet for request with passed request id.
	SendResponse(context.Context, network.RequestID, *packet.Packet) error

	// SendPacket low-level send packet without requestId and without spawning a waiting future
	SendPacket(ctx context.Context, p *packet.Packet) error

	// Listen starts thread to listen incoming packets.
	Listen(ctx context.Context) error

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
func NewTransport(cfg configuration.Transport) (Transport, error) {
	switch cfg.Protocol {
	case "TCP":
		listener, err := net.Listen("tcp", cfg.Address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to listen UDP")
		}
		publicAddress, err := Resolve(cfg, listener.Addr().String())
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve public address")
		}
		return newTCPTransport(listener, publicAddress)
	case "PURE_UDP":
		conn, err := net.ListenPacket("udp", cfg.Address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to listen TCP")
		}
		publicAddress, err := Resolve(cfg, conn.LocalAddr().String())
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve public address")
		}
		return newUDPTransport(conn, publicAddress)
	default:
		return nil, errors.New("invalid transport configuration")
	}
}

// Resolve resolves public address
func Resolve(cfg configuration.Transport, address string) (string, error) {
	resolver, err := createResolver(cfg)
	if err != nil {
		return "", errors.Wrap(err, "[ Resolve ] Failed to create resolver")
	}
	publicAddress, err := resolver.Resolve(address)
	if err != nil {
		return "", errors.Wrap(err, "[ Resolve ] Failed to resolve public address")
	}
	return publicAddress, nil
}

func createResolver(cfg configuration.Transport) (resolver.PublicAddressResolver, error) {
	if cfg.FixedPublicAddress != "" {
		return resolver.NewFixedAddressResolver(cfg.FixedPublicAddress), nil
	}
	return resolver.NewExactResolver(), nil
}

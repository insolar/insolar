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
	"fmt"
	"net"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
)

const udpMaxPacketSize = 1400

type udpTransport struct {
	conn      net.PacketConn
	address   string
	processor DatagramProcessor
}

func newUDPTransport(listenAddress, fixedPublicAddress string) (*udpTransport, string, error) {
	conn, err := net.ListenPacket("udp", listenAddress)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to listen UDP")
	}
	publicAddress, err := resolver.Resolve(fixedPublicAddress, conn.LocalAddr().String())
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to resolve public address")
	}

	transport := &udpTransport{conn: conn}
	return transport, publicAddress, nil
}

// SendDatagram sends datagram to remote host
func (t *udpTransport) SendDatagram(ctx context.Context, address string, buff []byte) error {
	return t.send(ctx, address, buff)
}

// SetDatagramProcessor registers callback to process received datagram
func (t *udpTransport) SetDatagramProcessor(processor DatagramProcessor) {
	t.processor = processor
}

func (t *udpTransport) send(ctx context.Context, recvAddress string, data []byte) error {
	log.Debug("Sending PURE_UDP request")
	if len(data) > udpMaxPacketSize {
		return errors.New(fmt.Sprintf("udpTransport.send: too big input data. Maximum: %d. Current: %d",
			udpMaxPacketSize, len(data)))
	}

	// TODO: may be try to send second time if error
	// TODO: skip resolving every time by caching result
	udpAddr, err := net.ResolveUDPAddr("udp", recvAddress)
	if err != nil {
		return errors.Wrap(err, "udpTransport.send")
	}

	log.Debug("udpTransport.send: len = ", len(data))
	n, err := t.conn.WriteTo(data, udpAddr)
	stats.Record(ctx, consensus.SentSize.M(int64(n)))
	return errors.Wrap(err, "Failed to write data")
}

func (t *udpTransport) prepareListen() (net.PacketConn, error) {
	if t.conn != nil {
		t.address = t.conn.LocalAddr().String()
	} else {
		var err error
		t.conn, err = net.ListenPacket("udp", t.address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to listen UDP")
		}
	}

	return t.conn, nil
}

// Start starts networking.
func (t *udpTransport) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Start ] Start UDP transport")

	conn, err := t.prepareListen()
	if err != nil {
		logger.Infof("[ Start ] Failed to prepare UDP transport: " + err.Error())
		return err
	}

	go func() {
		for {
			//todo handle stop
			buf := make([]byte, udpMaxPacketSize)
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				logger.Error("failed to read UDP: ", err.Error())
				return // TODO: we probably shouldn't return here
			}

			stats.Record(ctx, consensus.RecvSize.M(int64(n)))
			go func() {
				err := t.processor.ProcessDatagram(addr.String(), buf[:n])
				if err != nil {
					logger.Error("failed to process UDP packet: ", err.Error())
				}
			}()
		}
	}()

	return nil
}

// Stop stops networking.
func (t *udpTransport) Stop(ctx context.Context) error {
	log.Info("Stop UDP transport")
	return t.conn.Close()
}

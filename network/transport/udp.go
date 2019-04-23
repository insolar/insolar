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
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
	"github.com/insolar/insolar/network/utils"
)

const udpMaxPacketSize = 1400

type udpTransport struct {
	conn    net.PacketConn
	handler DatagramHandler
	started uint32
	address string
	cancel  context.CancelFunc
}

func newUDPTransport(listenAddress, fixedPublicAddress string, handler DatagramHandler) (*udpTransport, error) {
	conn, err := net.ListenPacket("udp", listenAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen UDP")
	}
	publicAddress, err := resolver.Resolve(fixedPublicAddress, conn.LocalAddr().String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve public address")
	}

	transport := &udpTransport{conn: conn, address: publicAddress, handler: handler}
	return transport, nil
}

// SendDatagram sends datagram to remote host
func (t *udpTransport) SendDatagram(ctx context.Context, address string, data []byte) error {
	logger := inslogger.FromContext(ctx)
	if len(data) > udpMaxPacketSize {
		return errors.New(fmt.Sprintf("udpTransport.send: too big input data. Maximum: %d. Current: %d",
			udpMaxPacketSize, len(data)))
	}

	// TODO: may be try to send second time if error
	// TODO: skip resolving every time by caching result
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return errors.Wrap(err, "Failed to resolve UDP address")
	}

	logger.Debug("udpTransport.send: len = ", len(data))
	n, err := t.conn.WriteTo(data, udpAddr) // Write(data)
	if err != nil {
		return errors.Wrap(err, "========================================== Failed to write data")
	}
	stats.Record(ctx, consensus.SentSize.M(int64(n)))
	return nil
}

func (t *udpTransport) Address() string {
	return t.address
}

// Start starts networking.
func (t *udpTransport) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	if !atomic.CompareAndSwapUint32(&t.started, 0, 1) {
		var err error
		t.conn, err = net.ListenPacket("udp", t.address)
		if err != nil {
			return errors.Wrap(err, "failed to listen UDP")
		}
	}

	logger.Info("[ Start ] Start UDP transport")
	ctx, t.cancel = context.WithCancel(ctx)
	go t.loop(ctx)

	return nil
}

func (t *udpTransport) loop(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		err := t.conn.SetDeadline(time.Now().Add(time.Second * 12))
		if err != nil {
			logger.Error(err.Error())
		}

		buf := make([]byte, udpMaxPacketSize)
		n, addr, err := t.conn.ReadFrom(buf)

		if err != nil {
			if utils.IsConnectionClosed(err) {
				logger.Info("Connection closed, quiting ReadFrom loop")
				return
			}

			logger.Error("failed to read UDP: ", err.Error())
			continue
		}

		stats.Record(ctx, consensus.RecvSize.M(int64(n)))
		go t.handler.HandleDatagram(addr.String(), buf[:n])
	}
}

// Stop stops networking.
func (t *udpTransport) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	//if atomic.CompareAndSwapUint32(&t.started, 1, 0) {
	logger.Warn("Stop UDP transport")
	t.cancel()
	err := t.conn.Close()

	if err != nil {
		if utils.IsConnectionClosed(err) {
			logger.Error("Connection already closed")
		} else {
			return err
		}
	}
	// } else {
	// 	logger.Warn("Failed to stop transport")
	// }
	return nil
}

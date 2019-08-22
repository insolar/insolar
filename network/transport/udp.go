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
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
)

const (
	udpMaxPacketSize = 1450
)

type udpTransport struct {
	mutex              sync.RWMutex
	conn               net.PacketConn
	handler            DatagramHandler
	started            uint32
	fixedPublicAddress string
	cancel             context.CancelFunc
	address            string
}

func newUDPTransport(listenAddress, fixedPublicAddress string, handler DatagramHandler) *udpTransport {
	return &udpTransport{address: listenAddress, fixedPublicAddress: fixedPublicAddress, handler: handler}
}

// SendDatagram sends datagram to remote host
func (t *udpTransport) SendDatagram(ctx context.Context, address string, data []byte) error {
	if atomic.LoadUint32(&t.started) != 1 {
		return errors.New("failed to send datagram: transport is not started")
	}

	if len(data) > udpMaxPacketSize {
		return fmt.Errorf(
			"failed to send datagram: too big input data. Maximum: %d. Current: %d",
			udpMaxPacketSize,
			len(data),
		)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return errors.Wrap(err, "failed to resolve UDP address")
	}

	n, err := t.conn.WriteTo(data, udpAddr)
	if err != nil {
		// TODO: may be try to send second time if error
		return errors.Wrap(err, "failed to write data")
	}
	stats.Record(ctx, network.SentSize.M(int64(n)))
	return nil
}

func (t *udpTransport) Address() string {
	return t.address
}

// Start starts networking.
func (t *udpTransport) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	if atomic.CompareAndSwapUint32(&t.started, 0, 1) {

		t.mutex.Lock()
		defer t.mutex.Unlock()

		var err error
		t.conn, err = net.ListenPacket("udp", t.address)
		if err != nil {
			return errors.Wrap(err, "failed to listen UDP")
		}

		t.address, err = resolver.Resolve(t.fixedPublicAddress, t.conn.LocalAddr().String())
		if err != nil {
			return errors.Wrap(err, "failed to resolve public address")
		}

		logger.Info("[ Start ] Start UDP transport")
		ctx, t.cancel = context.WithCancel(ctx)
		go t.loop(ctx)
	}

	return nil
}

func (t *udpTransport) loop(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		buf := make([]byte, udpMaxPacketSize)
		n, addr, err := t.conn.ReadFrom(buf)

		if err != nil {
			if network.IsConnectionClosed(err) {
				logger.Info("[ loop ] Connection closed, quiting ReadFrom loop")
				return
			}

			logger.Warn("[ loop ] failed to read UDP: ", err)
			continue
		}

		stats.Record(ctx, network.RecvSize.M(int64(n)))
		go t.handler.HandleDatagram(ctx, addr.String(), buf[:n])
	}
}

// Stop stops networking.
func (t *udpTransport) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	if atomic.CompareAndSwapUint32(&t.started, 1, 0) {
		logger.Info("[ Stop ] Stop UDP transport")

		t.cancel()
		err := t.conn.Close()
		if err != nil {
			if !network.IsConnectionClosed(err) {
				return err
			}
			logger.Warn("[ Stop ] Connection already closed")
		}
	}
	return nil
}

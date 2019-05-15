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
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
	"github.com/insolar/insolar/network/utils"
)

const (
	keepAlivePeriod = 10 * time.Second
)

type tcpTransport struct {
	listener           *net.TCPListener
	address            string
	started            uint32
	fixedPublicAddress string
	handler            StreamHandler
	cancel             context.CancelFunc
}

func newTCPTransport(listenAddress, fixedPublicAddress string, handler StreamHandler) *tcpTransport {
	return &tcpTransport{
		address:            listenAddress,
		fixedPublicAddress: fixedPublicAddress,
		handler:            handler,
	}
}

func (t *tcpTransport) Address() string {
	return t.address
}

func (t *tcpTransport) Dial(ctx context.Context, address string) (io.ReadWriteCloser, error) {
	logger := inslogger.FromContext(ctx).WithField("address", address)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, errors.New("[ Dial ] Failed to get tcp address")
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		logger.Error("[ Dial ] Failed to open connection: ", err)
		return nil, errors.Wrap(err, "[ Dial ] Failed to open connection")
	}

	setupConnection(ctx, conn)

	return conn, nil
}

// Start starts networking.
func (t *tcpTransport) Start(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&t.started, 0, 1) {

		logger := inslogger.FromContext(ctx)
		logger.Info("[ Start ] Start TCP transport")
		ctx, t.cancel = context.WithCancel(ctx)

		addr, err := net.ResolveTCPAddr("tcp", t.address)
		if err != nil {
			return errors.Wrap(err, "Failed to resolve TCP addr")
		}

		t.listener, err = net.ListenTCP("tcp", addr)
		if err != nil {
			return errors.Wrap(err, "Failed to Listen TCP ")
		}

		t.address, err = resolver.Resolve(t.fixedPublicAddress, t.listener.Addr().String())
		if err != nil {
			return errors.Wrap(err, "Failed to resolve public address")
		}

		go t.listen(ctx)
	}
	return nil
}

func (t *tcpTransport) listen(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, err := t.listener.AcceptTCP()
		if err != nil {
			if utils.IsConnectionClosed(err) {
				logger.Info("[ listen ] Connection closed, quiting accept loop")
				return
			}

			logger.Error("[ listen ] Failed to accept connection: ", err)
			return
		}
		logger = logger.WithField("address", conn.RemoteAddr())
		logger.Infof("[ listen ] Accepted new connection")
		setupConnection(ctx, conn)

		go t.handler.HandleStream(conn.RemoteAddr().String(), conn)
	}
}

// Stop stops networking.
func (t *tcpTransport) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	if atomic.CompareAndSwapUint32(&t.started, 1, 0) {
		logger.Info("[ Stop ] Stop TCP transport")

		t.cancel()
		err := t.listener.Close()
		if err != nil {
			if !utils.IsConnectionClosed(err) {
				return err
			}
			logger.Info("[ Stop ] Connection already closed")
		}
	}
	return nil
}

func setupConnection(ctx context.Context, conn *net.TCPConn) {
	logger := inslogger.FromContext(ctx).WithField("address", conn.RemoteAddr())

	if err := conn.SetNoDelay(true); err != nil {
		logger.Error("[ setupConnection ] Failed to set connection no delay: ", err)
	}

	if err := conn.SetKeepAlivePeriod(keepAlivePeriod); err != nil {
		logger.Error("[ setupConnection ] Failed to set keep alive period", err)
	}

	if err := conn.SetKeepAlive(true); err != nil {
		logger.Error("[ setupConnection ] Failed to set keep alive", err)
	}
}

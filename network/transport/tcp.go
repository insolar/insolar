// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package transport

import (
	"context"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
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
	dialer             net.Dialer
}

func newTCPTransport(listenAddress, fixedPublicAddress string, handler StreamHandler) *tcpTransport {
	return &tcpTransport{
		address:            listenAddress,
		fixedPublicAddress: fixedPublicAddress,
		handler:            handler,
		dialer:             net.Dialer{Timeout: 3 * time.Second},
	}
}

func (t *tcpTransport) Address() string {
	return t.address
}

func (t *tcpTransport) Dial(ctx context.Context, address string) (io.ReadWriteCloser, error) {
	logger := inslogger.FromContext(ctx).WithField("address", address)

	conn, err := t.dialer.Dial("tcp", address)
	if err != nil {
		logger.Warn("[ Dial ] Failed to open connection: ", err)
		return nil, errors.Wrap(err, "[ Dial ] Failed to open connection")
	}

	setupConnection(ctx, conn.(*net.TCPConn))

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
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			if network.IsConnectionClosed(err) {
				logger.Info("[ listen ] Connection closed, quiting accept loop")
				return
			}

			logger.Warn("[ listen ] Failed to accept connection: ", err)
			return
		}
		logger.Infof("[ listen ] Accepted new connection")
		setupConnection(ctx, conn)

		go t.handler.HandleStream(ctx, conn.RemoteAddr().String(), conn)
	}
}

// Stop stops networking.
func (t *tcpTransport) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	t.cancel()

	if atomic.CompareAndSwapUint32(&t.started, 1, 0) {
		logger.Info("[ Stop ] Stop TCP transport")

		err := t.listener.Close()
		if err != nil {
			if !network.IsConnectionClosed(err) {
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
		logger.Warn("[ setupConnection ] Failed to set connection no delay: ", err)
	}

	if err := conn.SetKeepAlivePeriod(keepAlivePeriod); err != nil {
		logger.Warn("[ setupConnection ] Failed to set keep alive period", err)
	}

	if err := conn.SetKeepAlive(true); err != nil {
		logger.Warn("[ setupConnection ] Failed to set keep alive", err)
	}
}

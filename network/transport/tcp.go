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
	"io"
	"net"

	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/transport/pool"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/network/utils"
)

const (
	keepAlivePeriod = 10 * time.Second
)

type tcpTransport struct {
	baseTransport

	pool            pool.ConnectionPool
	openConnections *sync.WaitGroup
	listenAddr      *net.TCPAddr
	listener        *net.TCPListener
}

func newTCPTransport(listenAddress string, proxy relay.Proxy, publicAddress string) (*tcpTransport, error) {
	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		return nil, errors.Wrap(err, "[ newTCPTransport ] Failed to resolve listenAddress")
	}

	transport := &tcpTransport{
		baseTransport:   newBaseTransport(proxy, publicAddress),
		pool:            pool.NewConnectionPool(&tcpConnectionFactory{}),
		listenAddr:      listenAddr,
		openConnections: &sync.WaitGroup{},
	}

	transport.sendFunc = transport.send

	return transport, nil
}

func (t *tcpTransport) send(address string, data []byte) error {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx)

	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return errors.Wrap(err, "[ send ] Failed to resolve net address")
	}

	conn, err := t.openConnection(ctx, addr)
	if err != nil {
		return err
	}

	logger.Debug("[ send ] len = ", len(data))

	_, err = conn.Write(data)

	if err != nil {
		// if netErr, ok := err.(*net.OpError); ok {
		// 	switch realNetErr := netErr.Err.(type) {
		// 	case *os.SyscallError:
		// 		if realNetErr.Err == syscall.EPIPE {
		t.pool.CloseConnection(ctx, addr)
		conn, err = t.openConnection(ctx, addr)
		if err != nil {
			return errors.Wrap(err, "[ send ] Failed to get connection")
		}
		_, err = conn.Write(data)
		// 		}
		// 	}
		// }
	}

	return errors.Wrap(err, "[ send ] Failed to write data")
}

func (t *tcpTransport) prepareListen() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.openConnections = &sync.WaitGroup{}
	t.disconnectStarted = make(chan bool, 1)

	listener, err := net.ListenTCP("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = listener

	return nil
}

// Start starts networking.
func (t *tcpTransport) Listen(ctx context.Context, started chan struct{}) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Listen ] Start TCP transport")

	if err := t.prepareListen(); err != nil {
		logger.Info("[ Listen ] Failed to prepare TCP transport")
		return err
	}

	started <- struct{}{}
	for {
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			logger.Debugf("[ Listen ] Failed to accept connection: %s", err.Error())
			return errors.Wrap(err, "[ Listen ] Failed to accept connection")
		}

		go func(conn *net.TCPConn) {
			logger.Debugf("[ Listen ] Accepted new connection from %s", conn.RemoteAddr())

			setupConnection(ctx, conn)
			t.handleAcceptedConnection(ctx, conn, false, nil)
		}(conn)
	}
}

// Stop stops networking.
func (t *tcpTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	ctx := context.TODO()
	inslogger.FromContext(ctx).Info("[ Stop ] Stop TCP transport")

	t.prepareDisconnect()

	utils.CloseVerbose(t.listener)
	t.pool.Reset(ctx)
	t.openConnections.Wait()

	<-t.disconnectFinished
	t.pool.Reset(ctx) // Second reset to ensure all connection is closed after transport clients are called close.

	inslogger.FromContext(ctx).Info("[ Stop ] TCP transport stopped")
}

func (t *tcpTransport) handleAcceptedConnection(
	ctx context.Context,
	conn *net.TCPConn,
	registered bool,
	remoteAddr *net.TCPAddr) {

	t.mutex.RLock()
	wg := t.openConnections
	t.mutex.RUnlock()

	logger := inslogger.FromContext(ctx)
	wg.Add(1)

	var alreadyClosed bool

	defer func() {
		logger.Debugf("[ handleAcceptedConnection ] Closing connection %p - alreadyClosed: %v, registered: %v", conn, alreadyClosed, registered)

		if !alreadyClosed {
			if registered {
				t.pool.CloseConnection(ctx, remoteAddr)
			} else {
				utils.CloseVerbose(conn)
			}
		}
		wg.Done()
	}()

	for {
		p, err := t.serializer.DeserializePacket(conn)

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				logger.Warn("[ handleAcceptedConnection ] Connection closed by peer")
				return
			}

			if strings.Contains(err.Error(), "use of closed network connection") {
				logger.Debug("[ handleAcceptedConnection ] Connection closed in another goroutine")
				alreadyClosed = true
				return
			}

			logger.Error("[ handleAcceptedConnection ] Failed to deserialize packet: ", err.Error())
			continue
		}

		ctx, logger := inslogger.WithTraceField(context.Background(), p.TraceID)
		logger.Debug("[ handleAcceptedConnection ] Handling packet: ", p.RequestID)
		go t.packetHandler.Handle(ctx, p)

		if registered {
			continue
		}

		remoteAddr, err = net.ResolveTCPAddr("tcp", p.RemoteAddress)
		if err != nil {
			logger.Errorf("[ handleAcceptedConnection ] Failed to register connection: %s", err.Error())
		}

		var closeConnection bool
		if closeConnection, registered = t.registerConnection(ctx, remoteAddr, conn); closeConnection {
			return
		}
	}
}

func (t *tcpTransport) registerConnection(
	ctx context.Context,
	remoteAddr *net.TCPAddr,
	conn *net.TCPConn,
) (closeConnection bool, registered bool) {
	logger := inslogger.FromContext(ctx)

	if !t.pool.RegisterConnection(context.Background(), remoteAddr, conn) && remoteAddr.String() < t.publicAddress {
		logger.Infof("[ registerConnection ] Connection %p to %s already registered", conn, remoteAddr)
		closeConnection = true
		return
	}

	registered = true
	return
}

func (t *tcpTransport) openConnection(ctx context.Context, addr *net.TCPAddr) (net.Conn, error) {
	created, conn, err := t.pool.GetConnection(ctx, addr)
	if err != nil {
		return nil, errors.Wrap(err, "[ openConnection ] Failed to get connection")
	}
	if created {
		go t.handleAcceptedConnection(ctx, conn.(*net.TCPConn), true, addr)
	}

	return conn, nil
}

func setupConnection(ctx context.Context, conn *net.TCPConn) {
	logger := inslogger.FromContext(ctx)

	err := conn.SetKeepAlivePeriod(keepAlivePeriod)
	if err != nil {
		logger.Error("[ setupConnection ] Failed to set keep alive")
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		logger.Error("[ setupConnection ] Failed to set keep alive")
	}

	err = conn.SetNoDelay(true)
	if err != nil {
		logger.Errorln("[ setupConnection ] Failed to set connection no delay: ", err.Error())
	}
}

type tcpConnectionFactory struct{}

func (*tcpConnectionFactory) CreateConnection(ctx context.Context, address net.Addr) (net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	tcpAddress, ok := address.(*net.TCPAddr)
	if !ok {
		return nil, errors.New("[ createConnection ] Failed to get tcp address")
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		logger.Errorf("[ createConnection ] Failed to open connection to %s: %s", address, err.Error())
		return nil, errors.Wrap(err, "[ createConnection ] Failed to open connection")
	}

	setupConnection(ctx, conn)

	return conn, nil
}

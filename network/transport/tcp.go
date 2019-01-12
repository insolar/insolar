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
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/transport/pool"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

const (
	keepAlivePeriod = 10 * time.Second
)

type tcpTransport struct {
	baseTransport

	openConnections sync.WaitGroup
	pool            pool.ConnectionPool
	listenAddr      *net.TCPAddr
	listener        *net.TCPListener
}

func newTCPTransport(listenAddress string, proxy relay.Proxy, publicAddress string) (*tcpTransport, error) {
	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		return nil, errors.Wrap(err, "[ newTCPTransport ] Failed to resolve listenAddress")
	}

	transport := &tcpTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		pool:          pool.NewConnectionPool(&tcpConnectionFactory{}),
		listenAddr:    listenAddr,
	}

	transport.sendFunc = transport.send

	return transport, nil
}

func (t *tcpTransport) send(address string, data []byte) error {
	ctx := context.TODO()
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
		// All this to check is error EPIPE
		if netErr, ok := err.(*net.OpError); ok {
			switch realNetErr := netErr.Err.(type) {
			case *os.SyscallError:
				if realNetErr.Err == syscall.EPIPE {
					t.pool.CloseConnection(ctx, addr)
					conn, err = t.openConnection(ctx, addr)
					if err != nil {
						return err
					}
					_, err = conn.Write(data)
				}
			}
		}
	}

	return errors.Wrap(err, "[ send ] Failed to write data")
}

// Start starts networking.
func (t *tcpTransport) Listen(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Listen ] Start TCP transport")

	t.mutex.Lock()

	t.openConnections = sync.WaitGroup{}
	t.disconnectStarted = make(chan bool, 1)

	listener, err := net.ListenTCP("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = listener

	t.mutex.Unlock()

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

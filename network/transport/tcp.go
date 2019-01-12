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
	"sync/atomic"
	"syscall"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
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

	stopped    uint32
	pool       pool.ConnectionPool
	listenAddr *net.TCPAddr
	listener   *net.TCPListener
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

	_, conn, err := t.pool.GetConnection(ctx, addr)
	if err != nil {
		return errors.Wrap(err, "[ send ] Failed to get connection")
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
					_, conn, err = t.pool.GetConnection(ctx, addr)
					if err != nil {
						return errors.Wrap(err, "[ send ] Failed to get connection")
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

	t.stopped = 0
	t.disconnectStarted = make(chan bool, 1)
	t.disconnectFinished = make(chan bool, 1)

	listener, err := net.ListenTCP("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = listener

	t.mutex.Unlock()

	for {
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			<-t.disconnectFinished
			logger.Error("[ Listen ] Failed to accept connection: ", err.Error())
			return errors.Wrap(err, "[ Listen ] Failed to accept connection")
		}

		go func(conn *net.TCPConn) {
			logger.Debugf("[ Listen ] Accepted new connection from %s", conn.RemoteAddr())

			setupConnection(ctx, conn)
			t.handleAcceptedConnection(conn)
		}(conn)
	}
}

// Stop stops networking.
func (t *tcpTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Info("[ Stop ] Stop TCP transport")
	t.prepareDisconnect()

	t.stopped = 1
	utils.CloseVerbose(t.listener)
	t.pool.Reset(ctx)
}

func (t *tcpTransport) handleAcceptedConnection(conn net.Conn) {
	defer utils.CloseVerbose(conn)
	closed := false

	for {
		if atomic.LoadUint32(&t.stopped) == 1 && closed {
			closed = true
			log.Debugf("[ handleAcceptedConnection ] Stop handling connection: %s", conn.RemoteAddr().String())
		}

		err := conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
		if err != nil {
			log.Errorf("[ handleAcceptedConnection ] Failed to set read deadline", err.Error())
		}

		msg, err := t.serializer.DeserializePacket(conn)

		if err != nil {
			if err == io.EOF {
				log.Warn("[ handleAcceptedConnection ] Connection closed by peer")
				return
			}

			if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
				if closed {
					return
				}
				continue
			}

			log.Error("[ handleAcceptedConnection ] Failed to deserialize packet: ", err.Error())
		} else {
			log.Debug("[ handleAcceptedConnection ] Handling packet: ", msg.RequestID)

			go t.packetHandler.Handle(context.TODO(), msg)
		}
	}
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

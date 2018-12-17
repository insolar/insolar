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
	"sync"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"
)

type tcpTransport struct {
	baseTransport
	l net.Listener

	conns     map[string]net.Conn
	connMutex sync.RWMutex
}

func newTCPTransport(addr string, proxy relay.Proxy, publicAddress string) (*tcpTransport, error) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	transport := &tcpTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		l:             listener,
		conns:         make(map[string]net.Conn),
	}

	transport.sendFunc = transport.send

	return transport, nil
}

func (tcp *tcpTransport) send(recvAddress string, data []byte) error {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx)

	tcpAddr, err := net.ResolveTCPAddr("tcp", recvAddress)
	if err != nil {
		return errors.Wrap(err, "[ send ] Failed to resolve tcp address")
	}

	tcp.connMutex.RLock()
	conn, ok := tcp.conns[tcpAddr.String()]
	tcp.connMutex.RUnlock()

	if !ok || tcp.connectionClosed(conn) {
		tcp.connMutex.Lock()

		conn, ok = tcp.conns[tcpAddr.String()]
		if !ok || tcp.connectionClosed(conn) {
			logger.Debugf("[ send ] Failed to retrieve connection to %s", tcpAddr)

			conn, err = tcp.openTCP(ctx, tcpAddr)
			if err != nil {
				tcp.connMutex.Unlock()
				return errors.Wrap(err, "[ send ] Failed to create TCP connection")
			}
			tcp.conns[conn.RemoteAddr().String()] = conn
			logger.Debugf("[ openTCP ] Added connection for %s. Current pool size: %d", conn.RemoteAddr(), len(tcp.conns))
		}

		tcp.connMutex.Unlock()
	}

	log.Debug("[ send ] len = ", len(data))
	_, err = conn.Write(data)
	return errors.Wrap(err, "[ send ] Failed to write data")
}

func (tcp *tcpTransport) openTCP(ctx context.Context, addr *net.TCPAddr) (net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Warnf("[ openTCP ] Failed to open connection to %s", addr)
		return nil, errors.Wrap(err, "[ openTCP ] Failed to open connection")
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		logger.Error("[ openTCP ] Failed to set keep alive")
	}

	return conn, nil
}

// Consuming 1 byte; only usable for outgoing connections.
func (tcp *tcpTransport) connectionClosed(conn net.Conn) bool {
	err := conn.SetReadDeadline(time.Now())
	if err != nil {
		log.Errorln("[ connectionClosed ] Failed to set connection deadline: ", err.Error())
	}

	n, err := conn.Read(make([]byte, 1))

	if err == io.EOF || n > 0 {
		err := conn.Close()
		if err != nil {
			log.Errorln("[ connectionClosed ] Failed to close connection: ", err.Error())
		} else {
			log.Debug("[ connectionClosed ] Close connection to %s", conn.RemoteAddr())
		}

		delete(tcp.conns, conn.RemoteAddr().String())
		return true
	}

	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		log.Errorln("[ connectionClosed ] Failed to set connection deadline: ", err.Error())
	}

	return false
}

// Start starts networking.
func (tcp *tcpTransport) Listen(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Listen ] Start TCP transport")
	for {
		conn, err := tcp.l.Accept()
		if err != nil {
			<-tcp.disconnectFinished
			logger.Error("[ Listen ] Failed to accept connection: ", err.Error())
			return errors.Wrap(err, "[ Listen ] Failed to accept connection")
		}

		logger.Debugf("[ Listen ] Accepted new connection from %s", conn.RemoteAddr())

		go tcp.handleAcceptedConnection(conn)
	}
}

// Stop stops networking.
func (tcp *tcpTransport) Stop() {
	tcp.mutex.Lock()
	defer tcp.mutex.Unlock()

	log.Info("[ Stop ] Stop TCP transport")
	tcp.prepareDisconnect()

	err := tcp.l.Close()
	if err != nil {
		log.Errorln("[ Stop ] Failed to close socket: ", err.Error())
	}

	for addr, conn := range tcp.conns {
		err := conn.Close()
		if err != nil {
			log.Errorln("[ Stop ] Failed to close outgoing connection: ", err.Error())
		}
		delete(tcp.conns, addr)
	}
}

func (tcp *tcpTransport) handleAcceptedConnection(conn net.Conn) {
	for {
		msg, err := tcp.serializer.DeserializePacket(conn)

		if err != nil {
			if err == io.EOF {
				log.Warn("[ handleAcceptedConnection ] Connection closed by sender")
				return
			}

			log.Error("[ handleAcceptedConnection ] Failed to deserialize packet: ", err.Error())
		} else {
			log.Info("[ handleAcceptedConnection ] Handling packet")
			tcp.handlePacket(msg)
		}
	}
}

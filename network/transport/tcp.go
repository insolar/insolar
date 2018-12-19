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
	l         net.Listener
	addr      string
	isStarted bool

	conns     map[string]net.Conn
	connMutex sync.RWMutex
}

func newTCPTransport(addr string, proxy relay.Proxy, publicAddress string) (*tcpTransport, error) {

	transport := &tcpTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		addr:          addr,
		conns:         make(map[string]net.Conn),
	}

	transport.sendFunc = transport.send

	return transport, nil
}

func (t *tcpTransport) send(recvAddress string, data []byte) error {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx)

	tcpAddr, err := net.ResolveTCPAddr("tcp", recvAddress)
	if err != nil {
		return errors.Wrap(err, "[ send ] Failed to resolve t address")
	}

	t.connMutex.RLock()
	conn, ok := t.conns[tcpAddr.String()]
	t.connMutex.RUnlock()

	if !ok || t.connectionClosed(conn) {
		t.connMutex.Lock()
		if ok && t.connectionClosed(conn) {
			logger.Debugf("[ send ] Delete addr %s from map of connections: %s", tcpAddr)
			delete(t.conns, conn.RemoteAddr().String())
		}

		conn, ok = t.conns[tcpAddr.String()]
		logger.Debugf("[ send ] Finding addr %s in map of connections: %s", tcpAddr, ok)
		if !ok || t.connectionClosed(conn) {
			if ok && t.connectionClosed(conn) {
				logger.Debugf("[ send ] Delete addr %s from map of connections: %s", tcpAddr)
				delete(t.conns, conn.RemoteAddr().String())
			}

			logger.Debugf("[ send ] Failed to retrieve connection to %s", tcpAddr)

			conn, err = t.openTCP(ctx, tcpAddr)
			if err != nil {
				t.connMutex.Unlock()
				return errors.Wrap(err, "[ send ] Failed to create TCP connection")
			}
			t.conns[conn.RemoteAddr().String()] = conn
			logger.Debugf("[ openTCP ] Added connection for %s. Current pool size: %d", conn.RemoteAddr(), len(t.conns))
		}

		t.connMutex.Unlock()
	}

	log.Debug("[ send ] len = ", len(data))
	_, err = conn.Write(data)
	return errors.Wrap(err, "[ send ] Failed to write data")
}

func (t *tcpTransport) openTCP(ctx context.Context, addr *net.TCPAddr) (net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Errorf("[ openTCP ] Failed to open connection to %s: %s", addr, err.Error())
		return nil, errors.Wrap(err, "[ openTCP ] Failed to open connection")
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		logger.Error("[ openTCP ] Failed to set keep alive")
	}

	return conn, nil
}

// Consuming 1 byte; only usable for outgoing connections.
func (t *tcpTransport) connectionClosed(conn net.Conn) bool {
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

		return true
	}

	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		log.Errorln("[ connectionClosed ] Failed to set connection deadline: ", err.Error())
	}

	return false
}

func (t *tcpTransport) startListen() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.isStarted {
		return errors.New("[ startListen ] Failed to listen, already listening")
	}

	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		return errors.Wrap(err, "[ startListen ] Failed open listen socket")
	}

	t.isStarted = true
	t.l = listener
	t.prepareListen()

	return nil
}

// Start starts networking.
func (t *tcpTransport) Listen(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Listen ] Start TCP transport")

	if err := t.startListen(); err != nil {
		return errors.Wrap(err, "[ Listen ] failed to startListen")
	}

	t.prepareListen()
	for {
		conn, err := t.l.Accept()
		if err != nil {
			<-t.disconnectFinished
			logger.Error("[ Listen ] Failed to accept connection: ", err.Error())
			return errors.Wrap(err, "[ Listen ] Failed to accept connection")
		}

		logger.Debugf("[ Listen ] Accepted new connection from %s", conn.RemoteAddr())

		go t.handleAcceptedConnection(conn)
	}
}

// Stop stops networking.
func (t *tcpTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Info("[ Stop ] Stop TCP transport")
	t.prepareDisconnect()

	if !t.isStarted {
		return
	}

	err := t.l.Close()
	if err != nil {
		log.Errorln("[ Stop ] Failed to close socket: ", err.Error())
	}
	t.isStarted = false

	t.connMutex.Lock()
	defer t.connMutex.Unlock()

	for addr, conn := range t.conns {
		err := conn.Close()
		if err != nil {
			log.Errorln("[ Stop ] Failed to close outgoing connection: ", err.Error())
		}
		delete(t.conns, addr)
	}
}

func (t *tcpTransport) handleAcceptedConnection(conn net.Conn) {
	for {
		msg, err := t.serializer.DeserializePacket(conn)

		if err != nil {
			if err == io.EOF {
				log.Warn("[ handleAcceptedConnection ] Connection closed by sender")
				return
			}

			log.Error("[ handleAcceptedConnection ] Failed to deserialize packet: ", err.Error())
		} else {
			log.Debug("[ handleAcceptedConnection ] Handling packet: ", msg.RequestID)

			go t.handlePacket(msg)
		}
	}
}

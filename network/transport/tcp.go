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

const defaultTimeout = 100 * time.Millisecond

type tcpTransport struct {
	baseTransport
	l       net.Listener
	maxChan chan bool

	conns     map[net.Addr]net.Conn
	connMutex *sync.Mutex
}

func newTCPTransport(addr string, proxy relay.Proxy, publicAddress string) (*tcpTransport, error) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	transport := &tcpTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		l:             listener,
		maxChan:       make(chan bool, 1000),
		conns:         make(map[net.Addr]net.Conn),
		connMutex:     &sync.Mutex{},
	}

	transport.sendFunc = transport.send

	return transport, nil
}

func (tcp *tcpTransport) send(recvAddress string, data []byte) error {
	logger := inslogger.FromContext(context.Background())

	tcpAddr, err := net.ResolveTCPAddr("tcp", recvAddress)
	if err != nil {
		return errors.Wrap(err, "tcpTransport.send")
	}

	tcp.connMutex.Lock()
	conn, ok := tcp.conns[tcpAddr]
	if !ok || tcp.connectionClosed(conn) {
		addr := tcpAddr.String()

		logger.Debugf("[ send ] Failed to retrieve connection to %s", addr)
		conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Warnf("[ send ] Failed to open connection to %s", addr)
			tcp.connMutex.Unlock()
			return errors.Wrap(err, "tcpTransport.send")
		}
		tcp.conns[conn.RemoteAddr()] = conn
	}
	tcp.connMutex.Unlock()

	tcp.maxChan <- true

	log.Debug("tcpTransport.send: len = ", len(data))
	_, err = conn.Write(data)
	return errors.Wrap(err, "Failed to write data")
}

func (tcp *tcpTransport) connectionClosed(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now())
	if _, err := conn.Read([]byte{}); err == io.EOF {
		conn.Close()
		<-tcp.maxChan
		return true
	}

	return false
}

// Start starts networking.
func (tcp *tcpTransport) Listen(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("Start TCP transport")
	for {

		tcp.maxChan <- true

		conn, err := tcp.l.Accept()
		if err != nil {
			<-tcp.maxChan
			<-tcp.disconnectFinished
			return errors.Wrap(err, "[ Listen ] Failed to accept connection")
		}

		tcp.connMutex.Lock()
		addr := conn.RemoteAddr()
		tcp.conns[addr] = conn
		logger.Debugf("[ Listen ] Accepted new connection from %s", addr.String())
		tcp.connMutex.Unlock()

		go tcp.handleAcceptedConnection(conn)
	}
}

// Stop stops networking.
func (tcp *tcpTransport) Stop() {
	tcp.mutex.Lock()
	defer tcp.mutex.Unlock()

	log.Info("Stop TCP transport")
	tcp.prepareDisconnect()

	err := tcp.l.Close()
	if err != nil {
		log.Errorln("Failed to close socket:", err.Error())
	}
}

func (tcp *tcpTransport) handleAcceptedConnection(conn net.Conn) {
	for {
		conn.SetReadDeadline(time.Now().Add(defaultTimeout))
		msg, err := tcp.serializer.DeserializePacket(conn)
		if err != nil {
			log.Error("[ handleAcceptedConnection ] ", err)
			return
		}

		tcp.handlePacket(msg)
	}
}

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
	"net"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"
)

type tcpTransport struct {
	baseTransport
	l       net.Listener
	maxChan chan bool
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
	}

	transport.sendFunc = transport.send

	return transport, nil
}

func (tcp *tcpTransport) send(recvAddress string, data []byte) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", recvAddress)
	if err != nil {
		return errors.Wrap(err, "tcpTransport.send")
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return errors.Wrap(err, "tcpTransport.send")
	}
	tcp.maxChan <- true
	defer func() { <-tcp.maxChan }()
	defer tcpConn.Close()

	log.Debug("tcpTransport.send: len = ", len(data))
	_, err = tcpConn.Write(data)
	return errors.Wrap(err, "Failed to write data")
}

// Start starts networking.
func (tcp *tcpTransport) Start(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("Start TCP transport")
	for {

		tcp.maxChan <- true

		conn, err := tcp.l.Accept()
		if err != nil {
			<-tcp.maxChan
			<-tcp.disconnectFinished
			return errors.Wrap(err, "[ Start ]")
		}

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
	defer conn.Close()
	msg, err := tcp.serializer.DeserializePacket(conn)
	if err != nil {
		log.Error("[ handleAcceptedConnection ] ", err)
		return
	}

	tcp.handlePacket(msg)

	<-tcp.maxChan
}

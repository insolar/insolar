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
	addr    string
	maxChan chan bool
}

func newTCPTransport(addr string, proxy relay.Proxy, publicAddress string) (*tcpTransport, error) {
	transport := &tcpTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		addr:          addr,
		maxChan:       make(chan bool, 1000),
	}

	transport.sendFunc = transport.send

	return transport, nil
}

func (t *tcpTransport) send(recvAddress string, data []byte) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", recvAddress)
	if err != nil {
		return errors.Wrap(err, "tcpTransport.send")
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return errors.Wrap(err, "tcpTransport.send")
	}
	t.maxChan <- true
	defer func() { <-t.maxChan }()
	defer tcpConn.Close()

	log.Debug("tcpTransport.send: len = ", len(data))
	_, err = tcpConn.Write(data)
	return errors.Wrap(err, "Failed to write data")
}

// Start starts networking.
func (t *tcpTransport) Listen(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("Start TCP transport")

	t.mutex.Lock()
	t.isStarted = true

	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}

	t.l = listener
	t.mutex.Unlock()

	for {

		t.maxChan <- true

		conn, err := t.l.Accept()
		if err != nil {
			<-t.maxChan
			<-t.disconnectFinished
			return errors.Wrap(err, "[ Start ]")
		}

		go t.handleAcceptedConnection(conn)
	}
}

// Stop stops networking.
func (t *tcpTransport) Stop() {
	log.Info("Stop TCP transport")

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.prepareDisconnect()

	if t.isStarted {
		err := t.l.Close()
		if err != nil {
			log.Errorln("Failed to close socket:", err.Error())
		}
	}
}

func (t *tcpTransport) handleAcceptedConnection(conn net.Conn) {
	defer conn.Close()
	msg, err := t.serializer.DeserializePacket(conn)
	if err != nil {
		log.Error("[ handleAcceptedConnection ] ", err)
		return
	}

	t.handlePacket(msg)

	<-t.maxChan
}

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
	"crypto/tls"
	"net"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/lucas-clemente/quic-go"
)

type quicTransport struct {
	baseTransport
	l quic.Listener
}

func newQuicTransport(conn net.PacketConn, proxy relay.Proxy, publicAddress string) (*quicTransport, error) {
	listener, err := quic.Listen(conn, generateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}

	transport := &quicTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		l:             listener,
	}

	transport.sendFunc = transport.send
	return transport, nil
}

func (tcp *quicTransport) send(recvAddress string, data []byte) error {
	ctx := context.Background()
	session, err := quic.DialAddrContext(ctx, recvAddress, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync()
	if err != nil {
		return err
	}

	_, err = stream.Write(data)
	if err != nil {
		return err
	}

	return err
}

// Start starts networking.
func (q *quicTransport) Start(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("Start TCP transport")
	for {

		//q.maxChan <- true

		session, err := q.l.Accept()
		if err != nil {
			//<-q.maxChan
			<-q.disconnectFinished
			return err
		}

		go q.handleAcceptedConnection(session)
	}
}

// Stop stops networking.
func (tcp *quicTransport) Stop() {
	tcp.mutex.Lock()
	defer tcp.mutex.Unlock()

	log.Info("Stop TCP transport")
	tcp.prepareDisconnect()

	err := tcp.l.Close()
	if err != nil {
		log.Errorln("Failed to close socket:", err.Error())
	}
}

func (tcp *quicTransport) handleAcceptedConnection(conn quic.Session) {
	/*
		defer conn.Close()
		msg, err := tcp.serializer.DeserializePacket(conn)
		if err != nil {
			log.Error("[ handleAcceptedConnection ] ", err)
			return
		}

		tcp.handlePacket(msg)
	*/
}

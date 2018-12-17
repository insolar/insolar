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
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"

	"github.com/anacrolix/utp"
)

type utpTransport struct {
	baseTransport

	socket *utp.Socket
}

func newUTPTransport(conn net.PacketConn, proxy relay.Proxy, publicAddress string) (*utpTransport, error) {
	socket, err := utp.NewSocketFromPacketConn(conn)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create socket")
	}

	transport := &utpTransport{
		socket:        socket,
		baseTransport: newBaseTransport(proxy, publicAddress),
	}
	transport.sendFunc = transport.send
	return transport, nil
}

// Start starts networking.
func (t *utpTransport) Listen(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("Start UTP transport")
	for {
		conn, err := t.socket.Accept()

		if err != nil {
			<-t.disconnectFinished
			return err
		}

		go t.handleAcceptedConnection(conn)
	}
}

// Stop stops networking.
func (t *utpTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Info("Stop UTP transport")
	t.prepareDisconnect()

	err := t.socket.CloseNow()
	if err != nil {
		log.Errorln("Failed to close socket:", err.Error())
	}
}

func (t *utpTransport) send(recvAddress string, data []byte) error {
	conn, err := t.socketDialTimeout(recvAddress, time.Second)
	if err != nil {
		return errors.Wrap(err, "Failed to socket dial")
	}
	defer conn.Close()

	_, err = conn.Write(data)
	return errors.Wrap(err, "Failed to write data")
}

func (t *utpTransport) socketDialTimeout(addr string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	return t.socket.DialContext(ctx, "", addr)
}

func (t *utpTransport) handleAcceptedConnection(conn net.Conn) {
	for {
		// Wait for Packets
		msg, err := t.serializer.DeserializePacket(conn)
		if err != nil {
			// TODO should we penalize this Host somehow ? Ban it ?
			// if err.Error() != "EOF" {
			// }
			return
		}
		msg.RemoteAddress = t.getRemoteAddress(conn)
		t.handlePacket(msg)
	}
}

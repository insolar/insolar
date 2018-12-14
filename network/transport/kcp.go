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
	"github.com/xtaci/kcp-go"
)

type kcpTransport struct {
	baseTransport
	listener   *kcp.Listener
	blockCrypt kcp.BlockCrypt
}

func newKCPTransport(conn net.PacketConn, proxy relay.Proxy, publicAddress string) (*kcpTransport, error) {
	crypt, err := kcp.NewNoneBlockCrypt([]byte{})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create KCP transport")
	}

	lis, err := kcp.ServeConn(crypt, 0, 0, conn)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serve connection")
	}

	transport := &kcpTransport{
		listener:      lis,
		baseTransport: newBaseTransport(proxy, publicAddress),
		blockCrypt:    crypt,
	}
	transport.sendFunc = transport.send
	return transport, nil
}

// Start starts networking.
func (t *kcpTransport) Listen(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("Start KCP transport")
	for {
		if session, err := t.listener.AcceptKCP(); err == nil {
			go t.handleAcceptedConnection(session)
		} else {
			<-t.disconnectFinished
			return err
		}
	}
}

// Stop stops networking.
func (t *kcpTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Info("Stop KCP transport")
	err := t.listener.Close()
	if err != nil {
		log.Errorln("Failed to close socket:", err.Error())
	}

	t.prepareDisconnect()
}

func (t *kcpTransport) socketDialTimeout(addr string, timeout time.Duration) (*kcp.UDPSession, error) {
	return kcp.DialWithOptions(addr, t.blockCrypt, 0, 0)
}

func (t *kcpTransport) send(recvAddress string, data []byte) error {
	session, err := t.socketDialTimeout(recvAddress, time.Second)
	if err != nil {
		return errors.Wrap(err, "Failed to socket dial")
	}
	// No need explicit close KCP session.
	// defer conn.Close()

	_, err = session.Write(data)
	return errors.Wrap(err, "Failed to session write data")
}

func (t *kcpTransport) handleAcceptedConnection(session *kcp.UDPSession) {
	for {
		err := session.SetDeadline(time.Now().Add(time.Millisecond * 50))
		if err != nil {
			log.Errorln(err.Error())
		}
		// Wait for Packets
		msg, err := t.serializer.DeserializePacket(session)
		if err != nil {
			// TODO should we penalize this Host somehow ? Ban it ?
			// if err.Error() != "EOF" {
			// }
			return
		}
		msg.RemoteAddress = t.getRemoteAddress(session)
		log.Debugln("Handle connection from ", msg.RemoteAddress)
		t.handlePacket(msg)
	}
}

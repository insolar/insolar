/*
 *    Copyright 2019 Insolar
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
	"strings"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"
)

type transportSerializer interface {
	SerializePacket(q *packet.Packet) ([]byte, error)
	DeserializePacket(conn io.Reader) (*packet.Packet, error)
}

type baseSerializer struct{}

func (b *baseSerializer) SerializePacket(q *packet.Packet) ([]byte, error) {
	return packet.SerializePacket(q)
}

func (b *baseSerializer) DeserializePacket(conn io.Reader) (*packet.Packet, error) {
	return packet.DeserializePacket(conn)
}

type baseTransport struct {
	futureManager futureManager
	serializer    transportSerializer
	proxy         relay.Proxy
	packetHandler packetHandler

	disconnectStarted  chan bool
	disconnectFinished chan bool

	mutex *sync.RWMutex

	publicAddress string
	sendFunc      func(recvAddress string, data []byte) error
}

func newBaseTransport(proxy relay.Proxy, publicAddress string) baseTransport {
	futureManager := newFutureManager()
	return baseTransport{
		futureManager: futureManager,
		packetHandler: newPacketHandler(futureManager),
		proxy:         proxy,
		serializer:    &baseSerializer{},

		mutex: &sync.RWMutex{},

		disconnectStarted:  make(chan bool, 1),
		disconnectFinished: make(chan bool, 1),

		publicAddress: publicAddress,
	}
}

// SendRequest sends request packet and returns future.
func (t *baseTransport) SendRequest(ctx context.Context, msg *packet.Packet) (Future, error) {
	future := t.futureManager.Create(msg)
	err := t.SendPacket(ctx, msg)
	if err != nil {
		future.Cancel()
		return nil, errors.Wrap(err, "Failed to send transport packet")
	}
	metrics.NetworkPacketSentTotal.WithLabelValues(msg.Type.String()).Inc()
	return future, nil
}

// SendResponse sends response packet.
func (t *baseTransport) SendResponse(ctx context.Context, requestID network.RequestID, msg *packet.Packet) error {
	msg.RequestID = requestID

	return t.SendPacket(ctx, msg)
}

// Close closes packet channels.
func (t *baseTransport) Close() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	close(t.disconnectFinished)
}

// Packets returns incoming packets channel.
func (t *baseTransport) Packets() <-chan *packet.Packet {
	return t.packetHandler.Received()
}

// Stopped checks if networking is stopped already.
func (t *baseTransport) Stopped() <-chan bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.disconnectStarted
}

func (t *baseTransport) prepareDisconnect() {
	t.disconnectStarted <- true
	close(t.disconnectStarted)
}

func (t *baseTransport) getRemoteAddress(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

// PublicAddress returns transport public ip address
func (t *baseTransport) PublicAddress() string {
	return t.publicAddress
}

func (t *baseTransport) SendPacket(ctx context.Context, p *packet.Packet) error {
	var recvAddress string
	if t.proxy.ProxyHostsCount() > 0 {
		recvAddress = t.proxy.GetNextProxyAddress()
	}
	if len(recvAddress) == 0 {
		recvAddress = p.Receiver.Address.String()
	}

	data, err := t.serializer.SerializePacket(p)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize packet")
	}

	inslogger.FromContext(ctx).Debugf("Send %s packet to %s with RequestID = %d", p.Type, recvAddress, p.RequestID)
	return t.sendFunc(recvAddress, data)
}

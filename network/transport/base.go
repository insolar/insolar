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
	"io"
	"net"
	"strings"
	"sync"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
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
	sequenceGenerator sequenceGenerator
	futureManager     futureManager
	serializer        transportSerializer
	proxy             relay.Proxy

	received chan *packet.Packet

	disconnectStarted  chan bool
	disconnectFinished chan bool

	mutex *sync.RWMutex

	publicAddress string
	sendFunc      func(recvAddress string, data []byte) error
}

func newBaseTransport(proxy relay.Proxy, publicAddress string) baseTransport {
	return baseTransport{
		sequenceGenerator: newSequenceGenerator(),
		futureManager:     newFutureManager(),

		received: make(chan *packet.Packet),

		disconnectStarted:  make(chan bool, 1),
		disconnectFinished: make(chan bool, 1),

		mutex: &sync.RWMutex{},

		proxy:         proxy,
		publicAddress: publicAddress,
		serializer:    &baseSerializer{},
	}
}

// SendRequest sends request packet and returns future.
func (t *baseTransport) SendRequest(msg *packet.Packet) (Future, error) {
	msg.RequestID = packet.RequestID(t.sequenceGenerator.Generate())

	future := t.futureManager.Create(msg)

	go func(msg *packet.Packet, f Future) {
		err := t.SendPacket(msg)
		if err != nil {
			f.Cancel()
			log.Error(err)
		}
	}(msg, future)

	return future, nil
}

// SendResponse sends response packet.
func (t *baseTransport) SendResponse(requestID packet.RequestID, msg *packet.Packet) error {
	msg.RequestID = requestID

	return t.SendPacket(msg)
}

// Close closes packet channels.
func (t *baseTransport) Close() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	close(t.received)
	close(t.disconnectFinished)
}

// Packets returns incoming packets channel.
func (t *baseTransport) Packets() <-chan *packet.Packet {
	return t.received
}

// Stopped checks if networking is stopped already.
func (t *baseTransport) Stopped() <-chan bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.disconnectStarted
}

func (t *baseTransport) prepareListen() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.disconnectStarted = make(chan bool, 1)
}

func (t *baseTransport) prepareDisconnect() {
	t.disconnectStarted <- true
	close(t.disconnectStarted)
}

func (t *baseTransport) getRemoteAddress(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

func (t *baseTransport) handlePacket(msg *packet.Packet) {
	if msg.IsResponse {
		t.processResponse(msg)
		return
	}

	t.processRequest(msg)
}

func (t *baseTransport) processResponse(msg *packet.Packet) {
	log.Debugf("[ processResponse ] Process response %s with RequestID = %d", msg.RemoteAddress, msg.RequestID)

	future := t.futureManager.Get(msg)
	if future != nil {
		if shouldProcessPacket(future, msg) {
			log.Debugf("[ processResponse ] Processing future with RequestID = %s", msg.RequestID)
			future.SetResult(msg)
		} else {
			log.Debugf("[ processResponse ] Canceling future with RequestID = %s", msg.RequestID)
		}
		future.Cancel()
	}
}

func (t *baseTransport) processRequest(msg *packet.Packet) {
	log.Debugf("[ processRequest ] Process request %s with RequestID = %d", msg.RemoteAddress, msg.RequestID)
	t.received <- msg
}

// PublicAddress returns transport public ip address
func (t *baseTransport) PublicAddress() string {
	return t.publicAddress
}

func (t *baseTransport) SendPacket(p *packet.Packet) error {
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

	log.Debugf("Send packet to %s with RequestID = %d", recvAddress, p.RequestID)
	return t.sendFunc(recvAddress, data)
}

func shouldProcessPacket(future Future, msg *packet.Packet) bool {
	typesShouldBeEqual := msg.Type == future.Request().Type
	responseIsForRightSender := future.Actor().Equal(*msg.Sender)

	return typesShouldBeEqual && (responseIsForRightSender || msg.Type == types.Ping)
}

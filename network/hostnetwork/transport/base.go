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
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/pkg/errors"
)

type baseTransport struct {
	received chan *packet.Packet
	sequence *uint64

	disconnectStarted  chan bool
	disconnectFinished chan bool

	mutex   *sync.RWMutex
	futures map[packet.RequestID]Future

	proxy         relay.Proxy
	publicAddress string
	sendFunc      func(recvAddress string, data []byte) error
}

func newBaseTransport(proxy relay.Proxy, publicAddress string) baseTransport {
	return baseTransport{
		received: make(chan *packet.Packet),
		sequence: new(uint64),

		disconnectStarted:  make(chan bool),
		disconnectFinished: make(chan bool),

		mutex:   &sync.RWMutex{},
		futures: make(map[packet.RequestID]Future),

		proxy:         proxy,
		publicAddress: publicAddress,
	}
}

// SendRequest sends request packet and returns future.
func (t *baseTransport) SendRequest(msg *packet.Packet) (Future, error) {
	if !msg.IsValid() {
		return nil, errors.New("invalid packet")
	}

	msg.RequestID = t.generateID()

	future := t.createFuture(msg)

	err := t.sendPacket(msg)
	if err != nil {
		future.Cancel()
		return nil, errors.Wrap(err, "Failed to send packet")
	}

	return future, nil
}

// SendResponse sends response packet.
func (t *baseTransport) SendResponse(requestID packet.RequestID, msg *packet.Packet) error {
	msg.RequestID = requestID

	return t.sendPacket(msg)
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
	return t.disconnectStarted
}

func (t *baseTransport) generateID() packet.RequestID {
	id := AtomicLoadAndIncrementUint64(t.sequence)
	return packet.RequestID(id)
}

func (t *baseTransport) getRemoteAddress(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

func (t *baseTransport) createFuture(msg *packet.Packet) Future {
	newFuture := NewFuture(msg.RequestID, msg.Receiver, msg, func(f Future) {
		t.mutex.Lock()
		defer t.mutex.Unlock()

		delete(t.futures, f.Request().RequestID)
	})

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.futures[msg.RequestID] = newFuture

	metrics.NetworkFutures.WithLabelValues(msg.Type.String()).Set(float64(len(t.futures)))
	return newFuture
}

func (t *baseTransport) getFuture(msg *packet.Packet) Future {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.futures[msg.RequestID]
}

func (t *baseTransport) handlePacket(msg *packet.Packet) {
	if msg.IsResponse {
		t.processResponse(msg)
	} else {
		t.processRequest(msg)
	}
}

func (t *baseTransport) processResponse(msg *packet.Packet) {
	log.Debugf("Process response %s with RequestID = %s", msg.RemoteAddress, msg.RequestID)

	future := t.getFuture(msg)
	if future != nil && !shouldProcessPacket(future, msg) {
		future.SetResult(msg)
	}
	future.Cancel()
}

func (t *baseTransport) processRequest(msg *packet.Packet) {
	if msg.IsValid() {
		log.Debugf("Process request %s with RequestID = %s", msg.RemoteAddress, msg.RequestID)
		t.received <- msg
	}
}

// PublicAddress returns transport public ip address
func (t *baseTransport) PublicAddress() string {
	return t.publicAddress
}

func (t *baseTransport) sendPacket(msg *packet.Packet) error {
	var recvAddress string
	if t.proxy.ProxyHostsCount() > 0 {
		recvAddress = t.proxy.GetNextProxyAddress()
	}
	if len(recvAddress) == 0 {
		recvAddress = msg.Receiver.Address.String()
	}

	data, err := packet.SerializePacket(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize packet")
	}

	log.Debugf("Send packet to %s with RequestID = %s", recvAddress, msg.RequestID)
	return t.sendFunc(recvAddress, data)
}

func shouldProcessPacket(future Future, msg *packet.Packet) bool {
	return !future.Actor().Equal(*msg.Sender) && msg.Type != packet.TypePing || msg.Type != future.Request().Type
}

// AtomicLoadAndIncrementUint64 performs CAS loop, increments counter and returns old value.
func AtomicLoadAndIncrementUint64(addr *uint64) uint64 {
	for {
		val := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, val, val+1) {
			return val
		}
	}
}

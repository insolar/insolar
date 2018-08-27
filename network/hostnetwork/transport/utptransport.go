/*
 *    Copyright 2018 INS Ecosystem
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
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"

	"github.com/anacrolix/utp"
)

type utpTransport struct {
	socket *utp.Socket

	received chan *packet.Packet
	sequence *uint64

	disconnectStarted  chan bool
	disconnectFinished chan bool

	mutex   *sync.RWMutex
	futures map[packet.RequestID]Future

	proxy relay.Proxy
	publicAddress string
}

// NewUTPTransport creates utpTransport.
func NewUTPTransport(conn net.PacketConn, proxy relay.Proxy, publicAddress string) (Transport, error) {
	return newUTPTransport(conn, proxy, publicAddress)
}

func newUTPTransport(conn net.PacketConn, proxy relay.Proxy, publicAddress string) (*utpTransport, error) {
	socket, err := utp.NewSocketFromPacketConn(conn)
	if err != nil {
		return nil, err
	}

	transport := &utpTransport{
		socket: socket,

		received: make(chan *packet.Packet),
		sequence: new(uint64),

		disconnectStarted:  make(chan bool),
		disconnectFinished: make(chan bool),

		mutex:   &sync.RWMutex{},
		futures: make(map[packet.RequestID]Future),

		proxy: proxy,
		publicAddress: publicAddress,
	}

	return transport, nil
}

// SendRequest sends request packet and returns future.
func (t *utpTransport) SendRequest(msg *packet.Packet) (Future, error) {
	if !msg.IsValid() {
		return nil, errors.New("invalid packet")
	}

	msg.RequestID = t.generateID()

	future := t.createFuture(msg)

	err := t.sendPacket(msg)
	if err != nil {
		future.Cancel()
		return nil, err
	}

	return future, nil
}

// SendResponse sends response packet.
func (t *utpTransport) SendResponse(requestID packet.RequestID, msg *packet.Packet) error {
	msg.RequestID = requestID

	return t.sendPacket(msg)
}

// Start starts networking.
func (t *utpTransport) Start() error {
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

	t.disconnectStarted <- true
	close(t.disconnectStarted)

	err := t.socket.CloseNow()
	if err != nil {
		log.Println("Failed to close socket:", err.Error())
	}
}

// Close closes packet channels.
func (t *utpTransport) Close() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	close(t.received)
	close(t.disconnectFinished)
}

// Packets returns incoming packets channel.
func (t *utpTransport) Packets() <-chan *packet.Packet {
	return t.received
}

// Stopped checks if networking is stopped already.
func (t *utpTransport) Stopped() <-chan bool {
	return t.disconnectStarted
}

func (t *utpTransport) socketDialTimeout(addr string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	return t.socket.DialContext(ctx, "", addr)
}

func (t *utpTransport) generateID() packet.RequestID {
	id := AtomicLoadAndIncrementUint64(t.sequence)
	return packet.RequestID(id)
}

func (t *utpTransport) createFuture(msg *packet.Packet) Future {
	newFuture := NewFuture(msg.RequestID, msg.Receiver, msg, func(f Future) {
		t.mutex.Lock()
		defer t.mutex.Unlock()

		delete(t.futures, f.Request().RequestID)
	})

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.futures[msg.RequestID] = newFuture

	return newFuture
}

func (t *utpTransport) getFuture(msg *packet.Packet) Future {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.futures[msg.RequestID]
}

func (t *utpTransport) sendPacket(msg *packet.Packet) error {
	var recvAddress string
	if t.proxy.ProxyHostsCount() > 0 {
		recvAddress = t.proxy.GetNextProxyAddress()
	}
	if len(recvAddress) == 0 {
		recvAddress = msg.Receiver.Address.String()
	}
	conn, err := t.socketDialTimeout(recvAddress, time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := packet.SerializePacket(msg)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func (t *utpTransport) getRemoteAddress(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

func (t *utpTransport) handleAcceptedConnection(conn net.Conn) {
	for {
		// Wait for Packets
		msg, err := packet.DeserializePacket(conn)
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

func (t *utpTransport) handlePacket(msg *packet.Packet) {
	if msg.IsResponse {
		t.processResponse(msg)
	} else {
		t.processRequest(msg)
	}
}

func (t *utpTransport) processResponse(msg *packet.Packet) {
	future := t.getFuture(msg)
	if future != nil && !shouldProcessPacket(future, msg) {
		future.SetResult(msg)
	}
	future.Cancel()
}

func (t *utpTransport) processRequest(msg *packet.Packet) {
	if msg.IsValid() {
		t.received <- msg
	}
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

func (t *utpTransport) PublicAddress() string {
	return t.publicAddress
}
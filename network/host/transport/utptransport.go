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

	"github.com/insolar/insolar/network/host/message"
	"github.com/insolar/insolar/network/host/relay"

	"github.com/anacrolix/utp"
)

type utpTransport struct {
	socket *utp.Socket

	received chan *message.Message
	sequence *uint64

	disconnectStarted  chan bool
	disconnectFinished chan bool

	mutex   *sync.RWMutex
	futures map[message.RequestID]Future

	proxy relay.Proxy
}

// NewUTPTransport creates utpTransport.
func NewUTPTransport(conn net.PacketConn, proxy relay.Proxy) (Transport, error) {
	return newUTPTransport(conn, proxy)
}

func newUTPTransport(conn net.PacketConn, proxy relay.Proxy) (*utpTransport, error) {
	socket, err := utp.NewSocketFromPacketConn(conn)
	if err != nil {
		return nil, err
	}

	transport := &utpTransport{
		socket: socket,

		received: make(chan *message.Message),
		sequence: new(uint64),

		disconnectStarted:  make(chan bool),
		disconnectFinished: make(chan bool),

		mutex:   &sync.RWMutex{},
		futures: make(map[message.RequestID]Future),

		proxy: proxy,
	}

	return transport, nil
}

// SendRequest sends request message and returns future.
func (t *utpTransport) SendRequest(msg *message.Message) (Future, error) {
	if !msg.IsValid() {
		return nil, errors.New("invalid message")
	}

	msg.RequestID = t.generateID()

	future := t.createFuture(msg)

	err := t.sendMessage(msg)
	if err != nil {
		future.Cancel()
		return nil, err
	}

	return future, nil
}

// SendResponse sends response message.
func (t *utpTransport) SendResponse(requestID message.RequestID, msg *message.Message) error {
	msg.RequestID = requestID

	return t.sendMessage(msg)
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

// Close closes message channels.
func (t *utpTransport) Close() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	close(t.received)
	close(t.disconnectFinished)
}

// Messages returns incoming messages channel.
func (t *utpTransport) Messages() <-chan *message.Message {
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

func (t *utpTransport) generateID() message.RequestID {
	id := AtomicLoadAndIncrementUint64(t.sequence)
	return message.RequestID(id)
}

func (t *utpTransport) createFuture(msg *message.Message) Future {
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

func (t *utpTransport) getFuture(msg *message.Message) Future {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.futures[msg.RequestID]
}

func (t *utpTransport) sendMessage(msg *message.Message) error {
	var recvAddress string
	if t.proxy.ProxyNodesCount() > 0 {
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

	data, err := message.SerializeMessage(msg)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func (t *utpTransport) handleAcceptedConnection(conn net.Conn) {
	for {
		// Wait for Messages
		msg, err := message.DeserializeMessage(conn)
		if err != nil {
			// TODO should we penalize this Node somehow ? Ban it ?
			// if err.Error() != "EOF" {
			// }
			return
		}
		msg.RemoteAddress = strings.Split(conn.RemoteAddr().String(), ":")[0]
		t.handleMessage(msg)
	}
}

func (t *utpTransport) handleMessage(msg *message.Message) {
	if msg.IsResponse {
		t.processResponse(msg)
	} else {
		t.processRequest(msg)
	}
}

func (t *utpTransport) processResponse(msg *message.Message) {
	future := t.getFuture(msg)
	if future != nil && !shouldProcessMessage(future, msg) {
		future.SetResult(msg)
	}
	future.Cancel()
}

func (t *utpTransport) processRequest(msg *message.Message) {
	if msg.IsValid() {
		t.received <- msg
	}
}

func shouldProcessMessage(future Future, msg *message.Message) bool {
	return !future.Actor().Equal(*msg.Sender) && msg.Type != message.TypePing || msg.Type != future.Request().Type
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

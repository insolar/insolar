/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package transport

import (
	"context"
	"io"
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

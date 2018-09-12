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

package hostnetwork

import (
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/pkg/errors"
)

type mockHostHandler struct {
	AuthenticatedHost string
	ReceivedKey       string
}

func newMockHostHandler() *mockHostHandler {
	return &mockHostHandler{}
}

func (hh *mockHostHandler) GetHighKnownHostID() string {
	return ""
}

func (hh *mockHostHandler) GetOuterHostsCount() int {
	return 0
}

func (hh *mockHostHandler) ConfirmNodeRole(role string) bool {
	return false
}

func (hh *mockHostHandler) StoreRetrieve(key store.Key) ([]byte, bool) {
	return nil, false
}

func (hh *mockHostHandler) HtFromCtx(ctx hosthandler.Context) *routing.HashTable {
	address, _ := host.NewAddress("0.0.0.0:0")
	id1, _ := id.NewID()
	ht, _ := routing.NewHashTable(id1, address)
	return ht
}

func (hh *mockHostHandler) EqualAuthSentKey(targetID string, key []byte) bool {
	return false
}

func (hh *mockHostHandler) SendRequest(packet1 *packet.Packet) (transport.Future, error) {
	t := newMockTransport()
	sequenceNumber := transport.AtomicLoadAndIncrementUint64(t.sequence)

	if t.failNext {
		t.failNext = false
		return nil, errors.New("MockNetworking Error")
	}
	t.recv <- packet1

	return &mockFuture{result: t.send, request: packet1, actor: packet1.Receiver, requestID: packet.RequestID(sequenceNumber)}, nil
}

func (hh *mockHostHandler) FindHost(ctx hosthandler.Context, targetID string) (*host.Host, bool, error) {
	return nil, false, nil
}

func (hh *mockHostHandler) InvokeRPC(sender *host.Host, method string, args [][]byte) ([]byte, error) {
	if strings.EqualFold(method, "error") {
		return nil, errors.New("invoke error")
	}
	return nil, nil
}

func (hh *mockHostHandler) Store(key store.Key, data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	return nil
}

func (hh *mockHostHandler) AddPossibleProxyID(id string) {
}

func (hh *mockHostHandler) AddPossibleRelayID(id string) {
}

func (hh *mockHostHandler) AddProxyHost(targetID string) {
}

func (hh *mockHostHandler) AddSubnetID(ip, targetID string) {
}

func (hh *mockHostHandler) AddAuthSentKey(id string, key []byte) {
}

func (hh *mockHostHandler) AddRelayClient(host *host.Host) error {
	return nil
}

func (hh *mockHostHandler) AddReceivedKey(target string, key []byte) {
}

func (hh *mockHostHandler) AddHost(ctx hosthandler.Context, host *routing.RouteHost) {
}

func (hh *mockHostHandler) RemoveAuthHost(key string) {
}

func (hh *mockHostHandler) RemoveProxyHost(targetID string) {
}

func (hh *mockHostHandler) RemovePossibleProxyID(id string) {
}

func (hh *mockHostHandler) RemoveAuthSentKeys(targetID string) {
}

func (hh *mockHostHandler) RemoveRelayClient(host *host.Host) error {
	return nil
}

func (hh *mockHostHandler) SetHighKnownHostID(id string) {
}

func (hh *mockHostHandler) SetOuterHostsCount(hosts int) {
}

func (hh *mockHostHandler) SetAuthStatus(targetID string, status bool) {
}

func (hh *mockHostHandler) GetProxyHostsCount() int {
	return 0
}

func (hh *mockHostHandler) GetSelfKnownOuterHosts() int {
	return 0
}

func (hh *mockHostHandler) GetPacketTimeout() time.Duration {
	return 0
}

func (hh *mockHostHandler) GetReplicationTime() time.Duration {
	return 0
}

func (hh *mockHostHandler) GetExpirationTime(ctx hosthandler.Context, key []byte) time.Time {
	return time.Now()
}

func (hh *mockHostHandler) KeyIsReceived(targetID string) ([]byte, bool) {
	if hh.ReceivedKey == targetID {
		return []byte(targetID), true
	}
	return nil, false
}

func (hh *mockHostHandler) HostIsAuthenticated(targetID string) bool {
	if targetID == hh.AuthenticatedHost {
		return true
	}
	return false
}

func (hh *mockHostHandler) GetOriginHost() *host.Origin {
	address, _ := host.NewAddress("0.0.0.0:0")
	var ids []id.ID
	id1, _ := id.NewID()
	ids = append(ids, id1)
	origin, _ := host.NewOrigin(ids, address)
	return origin
}

func TestDispatchPacketType(t *testing.T) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	hh := newMockHostHandler()
	builder := packet.NewBuilder()
	authenticatedSenderAddress, _ := host.NewAddress("0.0.0.0:0")
	authenticatedSender := host.NewHost(authenticatedSenderAddress)
	authenticatedSender.ID, _ = id.NewID()
	hh.AuthenticatedHost = authenticatedSender.ID.KeyString()
	hh.ReceivedKey = authenticatedSender.ID.KeyString()

	t.Run("ping", func(t *testing.T) {
		pckt := packet.NewPingPacket(sender, receiver)
		DispatchPacketType(hh, getDefaultCtx(nil), pckt, builder)
	})

	t.Run("check node priv", func(t *testing.T) {
		pckt := packet.NewCheckNodePrivPacket(sender, receiver, "role")
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("authentication", func(t *testing.T) {
		pckt := packet.NewAuthPacket(packet.Unknown, sender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewAuthPacket(packet.BeginAuth, sender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewAuthPacket(packet.RevokeAuth, sender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewAuthPacket(packet.BeginAuth, authenticatedSender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewAuthPacket(packet.RevokeAuth, authenticatedSender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("check origin", func(t *testing.T) {
		pckt := packet.NewCheckOriginPacket(sender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewCheckOriginPacket(authenticatedSender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("known outer hosts", func(t *testing.T) {
		pckt := packet.NewKnownOuterHostsPacket(sender, receiver, 1)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("obtain ip", func(t *testing.T) {
		pckt := packet.NewObtainIPPacket(sender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("relay ownership", func(t *testing.T) {
		pckt := packet.NewRelayOwnershipPacket(sender, receiver, true)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewRelayOwnershipPacket(sender, receiver, false)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("relay", func(t *testing.T) {
		pckt := packet.NewRelayPacket(packet.Unknown, sender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewRelayPacket(packet.StartRelay, authenticatedSender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewRelayPacket(packet.StopRelay, authenticatedSender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = packet.NewRelayPacket(packet.Unknown, authenticatedSender, receiver)
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("rpc", func(t *testing.T) {
		pckt := builder.Type(packet.TypeRPC).Request(&packet.RequestDataRPC{}).Build()
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
		pckt = builder.Type(packet.TypeRPC).Request(&packet.RequestDataRPC{Method: "error"}).Build()
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("store", func(t *testing.T) {
		pckt := builder.Type(packet.TypeStore).Request(&packet.RequestDataStore{}).Build()
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("find host", func(t *testing.T) {
		pckt := builder.Type(packet.TypeFindHost).Request(&packet.RequestDataFindHost{Target: receiver.ID.GetKey()}).Build()
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})

	t.Run("find value", func(t *testing.T) {
		pckt := builder.Type(packet.TypeFindValue).Request(&packet.RequestDataFindValue{Target: sender.ID.GetKey()}).Build()
		DispatchPacketType(hh, getDefaultCtx(hh), pckt, packet.NewBuilder())
	})
}

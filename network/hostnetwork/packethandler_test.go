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
	return nil
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
	return nil, false
}

func (hh *mockHostHandler) HostIsAuthenticated(targetID string) bool {
	return false
}

func TestDispatchPacketType(t *testing.T) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	hh := newMockHostHandler()

	pckt := &packet.Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     packet.TypePing,
	}
	DispatchPacketType(hh, getDefaultCtx(nil), pckt, packet.NewBuilder())
}

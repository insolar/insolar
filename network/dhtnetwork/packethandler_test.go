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

package dhtnetwork

import (
	"crypto/ecdsa"
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/dhtnetwork/hosthandler"
	"github.com/insolar/insolar/network/dhtnetwork/routing"
	"github.com/insolar/insolar/network/dhtnetwork/rpc"
	"github.com/insolar/insolar/network/dhtnetwork/signhandler"
	"github.com/insolar/insolar/network/dhtnetwork/store"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/id"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"

	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockNetworkCommonFacade struct {
	cascade *cascade.Cascade
}

func newMockNetworkCommonFacade() hosthandler.NetworkCommonFacade {
	var c cascade.Cascade
	c.SendMessage = func(data core.Cascade, method string, args [][]byte) error {
		return nil
	}
	return &mockNetworkCommonFacade{
		cascade: &c,
	}
}

func (fac *mockNetworkCommonFacade) GetRPC() rpc.RPC {
	return nil
}

func (fac *mockNetworkCommonFacade) GetCascade() *cascade.Cascade {
	return fac.cascade
}

func (fac *mockNetworkCommonFacade) GetConsensus() consensus.Processor {
	return nil
}

func (fac *mockNetworkCommonFacade) SetConsensus(consensus.Processor) {
}

func (fac *mockNetworkCommonFacade) GetNetworkCoordinator() core.NetworkCoordinator {
	return nil
}

func (fac *mockNetworkCommonFacade) SetNetworkCoordinator(core.NetworkCoordinator) {
}

func (fac *mockNetworkCommonFacade) GetSignHandler() signhandler.SignHandler {
	return nil
}

func (fac *mockNetworkCommonFacade) OnPulse(pulse core.Pulse) {
}

type mockHostHandler struct {
	AuthenticatedHost string
	ReceivedKey       string
	FoundHost         *host.Host
	ncf               hosthandler.NetworkCommonFacade
	keeper            network.NodeKeeper
}

func newMockHostHandler() *mockHostHandler {
	return &mockHostHandler{ncf: newMockNetworkCommonFacade(), keeper: nodekeeper.NewNodeKeeper(nodekeeper.NewNode(core.RecordRef{}, nil, nil, 0, 0, "", ""))}
}

func (hh *mockHostHandler) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
}

func (hh *mockHostHandler) GetNetworkCommonFacade() hosthandler.NetworkCommonFacade {
	return hh.ncf
}

func (hh *mockHostHandler) RemoteProcedureCall(ctx hosthandler.Context, targetID string, method string, args [][]byte) (result []byte, err error) {
	return nil, nil
}

func (hh *mockHostHandler) Disconnect() {

}

func (hh *mockHostHandler) GetPrivateKey() *ecdsa.PrivateKey {
	return nil
}

func (hh *mockHostHandler) Listen() error {
	return nil
}

func (hh *mockHostHandler) Bootstrap() error {
	return nil
}

func (hh *mockHostHandler) GetHostsFromBootstrap() {
}

func (hh *mockHostHandler) ObtainIP() error {
	return nil
}

func (hh *mockHostHandler) NumHosts(ctx hosthandler.Context) int {
	return 0
}

func (hh *mockHostHandler) AnalyzeNetwork(ctx hosthandler.Context) error {
	return nil
}

func (hh *mockHostHandler) GetHighKnownHostID() string {
	return ""
}

func (hh *mockHostHandler) GetOuterHostsCount() int {
	return 0
}

func (hh *mockHostHandler) GetNodeID() core.RecordRef {
	return testutils.RandomRef()
}

func (hh *mockHostHandler) ConfirmNodeRole(role string) bool {
	return false
}

func (hh *mockHostHandler) StoreRetrieve(key store.Key) ([]byte, bool) {
	return nil, false
}

func (hh *mockHostHandler) CascadeSendMessage(data core.Cascade, targetID string, method string, args [][]byte) error {
	return nil
}

func (hh *mockHostHandler) GetActiveNodesList() []core.Node {
	return nil
}

func (hh *mockHostHandler) AddActiveNodes(activeNodes []core.Node) error {
	return nil
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

func getOriginHost(hh *mockHostHandler) *host.Host {
	origin := host.Host{}
	origin.ID = hh.GetOriginHost().IDs[0]
	origin.Address = hh.GetOriginHost().Address
	return &origin
}

func (hh *mockHostHandler) SendRequest(request *packet.Packet) (transport.Future, error) {
	t := newMockTransport()
	sequenceNumber := transport.AtomicLoadAndIncrementUint64(t.sequence)

	future := &mockFuture{result: t.send, request: request, actor: request.Receiver, requestID: packet.RequestID(sequenceNumber)}
	var response *packet.Packet
	origin := getOriginHost(hh)
	builder := packet.NewBuilder(origin)

	switch request.Type {
	case types.TypeRelay:
		data := request.Data.(*packet.RequestRelay)
		switch data.Command {
		case packet.StartRelay:
			response = builder.Response(&packet.ResponseRelay{State: relay.Started}).Build()
		case packet.StopRelay:
			response = builder.Response(&packet.ResponseRelay{State: relay.Stopped}).Build()
		case packet.BeginAuthentication:
			response = builder.Response(&packet.ResponseRelay{State: relay.NoAuth}).Build()
		case packet.RevokeAuthentication:
			response = builder.Response(&packet.ResponseRelay{State: relay.NoAuth}).Build()
		case packet.Unknown:
			response = builder.Response(&packet.ResponseRelay{State: relay.Unknown}).Build()
		default:
			response = builder.Response(&packet.ResponseRelay{State: relay.Error}).Build()
		}
	case types.TypeObtainIP:
		response = builder.Response(&packet.ResponseObtainIP{IP: "0.0.0.0"}).Build()
	case types.TypeCheckOrigin:
		response = builder.Response(&packet.ResponseCheckOrigin{AuthUniqueKey: []byte("asd")}).Build()
	case types.TypeAuthentication:
		response = builder.Response(&packet.ResponseAuthentication{Success: true, AuthUniqueKey: []byte("asd")}).Build()
	case types.TypeRelayOwnership:
		response = builder.Response(&packet.ResponseRelayOwnership{Accepted: true}).Build()
	}

	go future.SetResult(response)
	return future, nil
}

func (hh *mockHostHandler) FindHost(ctx hosthandler.Context, targetID string) (*host.Host, bool, error) {
	if hh.FoundHost == nil {
		return nil, false, nil
	}
	if strings.EqualFold(targetID, hh.FoundHost.ID.String()) {
		return hh.FoundHost, true, nil
	}
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

func (hh *mockHostHandler) GetActiveNodes() error {
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
	return 10 * time.Microsecond
}

func (hh *mockHostHandler) GetReplicationTime() time.Duration {
	return 20 * time.Microsecond
}

func (hh *mockHostHandler) GetExpirationTime(ctx hosthandler.Context, key []byte) time.Time {
	return time.Now()
}

func (hh *mockHostHandler) AddUncheckedNode(hostID id.ID, nonce []byte, ref core.RecordRef) {

}

func (hh *mockHostHandler) StartAuthorize() error {
	return nil
}

func (hh *mockHostHandler) AddUnsync(nodeID core.RecordRef, roles []core.NodeRole, address string,
	version string /*, publicKey *ecdsa.PublicKey*/) (chan core.Node, error) {
	return nil, nil
}

func (hh *mockHostHandler) SetNodeKeeper(keeper network.NodeKeeper) {

}

func (hh *mockHostHandler) RemoveUncheckedNode(hostID id.ID) {

}

func (hh *mockHostHandler) UncheckedNodeExist(hostID id.ID, nonce []byte) bool {
	return true
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

func (hh *mockHostHandler) Sign(nonce []byte) []byte {
	return nil
}

func mockSenderReceiver() (sender, receiver *host.Host) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender = host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver = host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	return
}

func TestDispatchPacketType(t *testing.T) {
	sender, receiver := mockSenderReceiver()
	hh := newMockHostHandler()
	builder := packet.NewBuilder(sender)
	authenticatedSenderAddress, _ := host.NewAddress("0.0.0.0:0")
	authenticatedSender := host.NewHost(authenticatedSenderAddress)
	authenticatedSender.ID, _ = id.NewID()
	hh.AuthenticatedHost = authenticatedSender.ID.String()
	hh.ReceivedKey = authenticatedSender.ID.String()
	origin := getOriginHost(hh)

	t.Run("ping", func(t *testing.T) {
		pckt := packet.NewPingPacket(sender, receiver)
		DispatchPacketType(hh, GetDefaultCtx(nil), pckt, builder)
	})

	t.Run("check node priv", func(t *testing.T) {
		pckt := builder.Type(types.TypeCheckNodePriv).Receiver(receiver).Request(&packet.RequestCheckNodePriv{RoleKey: "test string"}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("authentication", func(t *testing.T) {
		pckt := builder.Type(types.TypeAuthentication).
			Receiver(receiver).
			Request(&packet.RequestAuthentication{Command: packet.Unknown}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = builder.Type(types.TypeAuthentication).
			Receiver(receiver).
			Request(&packet.RequestAuthentication{Command: packet.BeginAuthentication}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = builder.Type(types.TypeAuthentication).
			Receiver(receiver).
			Request(&packet.RequestAuthentication{Command: packet.RevokeAuthentication}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = packet.NewBuilder(authenticatedSender).Type(types.TypeAuthentication).
			Receiver(receiver).
			Request(&packet.RequestAuthentication{Command: packet.BeginAuthentication}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = packet.NewBuilder(authenticatedSender).Type(types.TypeAuthentication).
			Receiver(receiver).
			Request(&packet.RequestAuthentication{Command: packet.RevokeAuthentication}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("check origin", func(t *testing.T) {
		pckt := builder.Type(types.TypeCheckOrigin).
			Receiver(receiver).
			Request(&packet.RequestCheckOrigin{}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = packet.NewBuilder(sender).Type(types.TypeCheckOrigin).
			Receiver(receiver).
			Request(&packet.RequestCheckOrigin{}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("known outer hosts", func(t *testing.T) {
		pckt := builder.Type(types.TypeKnownOuterHosts).
			Receiver(receiver).
			Request(&packet.RequestKnownOuterHosts{
				ID:         sender.ID.String(),
				OuterHosts: 1},
			).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("obtain ip", func(t *testing.T) {
		pckt := builder.Type(types.TypeObtainIP).
			Receiver(receiver).
			Request(&packet.RequestObtainIP{}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("relay ownership", func(t *testing.T) {
		pckt := builder.Type(types.TypeRelayOwnership).
			Receiver(receiver).
			Request(&packet.RequestRelayOwnership{Ready: true}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = builder.Type(types.TypeRelayOwnership).
			Receiver(receiver).
			Request(&packet.RequestRelayOwnership{Ready: false}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("relay", func(t *testing.T) {
		pckt := builder.Type(types.TypeRelay).
			Receiver(receiver).
			Request(&packet.RequestRelay{Command: packet.Unknown}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = builder.Type(types.TypeRelay).
			Receiver(receiver).
			Request(&packet.RequestRelay{Command: packet.StartRelay}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = builder.Type(types.TypeRelay).
			Receiver(receiver).
			Request(&packet.RequestRelay{Command: packet.StopRelay}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = builder.Type(types.TypeRelay).
			Receiver(receiver).
			Request(&packet.RequestRelay{Command: packet.Unknown}).
			Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("rpc", func(t *testing.T) {
		pckt := builder.Type(types.TypeRPC).Request(&packet.RequestDataRPC{}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		pckt = builder.Type(types.TypeRPC).Request(&packet.RequestDataRPC{Method: "error"}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("store", func(t *testing.T) {
		pckt := builder.Type(types.TypeStore).Request(&packet.RequestDataStore{}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("find host", func(t *testing.T) {
		pckt := builder.Type(types.TypeFindHost).Request(&packet.RequestDataFindHost{Target: receiver.ID.Bytes()}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("find value", func(t *testing.T) {
		pckt := builder.Type(types.TypeFindValue).Request(&packet.RequestDataFindValue{Target: sender.ID.Bytes()}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("send cascade", func(t *testing.T) {
		pckt := builder.Type(types.TypeCascadeSend).Request(&packet.RequestCascadeSend{
			Data: core.Cascade{}, RPC: packet.RequestDataRPC{}}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("pulse", func(t *testing.T) {
		pckt := builder.Type(types.TypePulse).Request(&packet.RequestPulse{Pulse: core.Pulse{}}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("get random hosts", func(t *testing.T) {
		pckt := builder.Type(types.TypeGetRandomHosts).Request(&packet.RequestGetRandomHosts{HostsNumber: 2}).Build()
		DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
	})

	t.Run("broken packet", func(t *testing.T) {
		pckt := builder.Type(types.TypeGetRandomHosts * 1024).Request(&packet.RequestGetRandomHosts{}).Build()
		_, err := DispatchPacketType(hh, GetDefaultCtx(hh), pckt, packet.NewBuilder(origin))
		assert.NotNil(t, err)
	})
}

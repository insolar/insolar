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

package hosthandler

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/rpc"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
)

// Context type is localized for future purposes.
// Network Host can have multiple IDs, but each action must be executed with only one ID.
// Context is used in all actions to select specific ID to work with.
type Context context.Context

// NetworkCommonFacade is used to implement additional features in network to evade showing them in the main HostHandler interface
type NetworkCommonFacade interface {
	GetRPC() rpc.RPC
	GetCascade() *cascade.Cascade
	GetPulseManager() core.PulseManager
	SetPulseManager(manager core.PulseManager)
	SetConsensus(insolarConsensus consensus.InsolarConsensus)
	GetConsensus() consensus.InsolarConsensus
}

// commonFacade implements a NetworkCommonFacade.
type commonFacade struct {
	rpcPtr  rpc.RPC
	cascade *cascade.Cascade
	pm      core.PulseManager
	ic      consensus.InsolarConsensus
}

// NewNetworkCommonFacade creates a NetworkCommonFacade.
func NewNetworkCommonFacade(r rpc.RPC, casc *cascade.Cascade) NetworkCommonFacade {
	return &commonFacade{rpcPtr: r, cascade: casc, pm: nil}
}

// GetRPC return an RPC pointer.
func (fac *commonFacade) GetRPC() rpc.RPC {
	return fac.rpcPtr
}

// GetCascade returns a cascade pointer.
func (fac *commonFacade) GetCascade() *cascade.Cascade {
	return fac.cascade
}

// GetPulseManager returns a pulse manager pointer.
func (fac *commonFacade) GetPulseManager() core.PulseManager {
	return fac.pm
}

// SetPulseManager sets a pulse manager to common facade.
func (fac *commonFacade) SetPulseManager(manager core.PulseManager) {
	fac.pm = manager
}

func (fac *commonFacade) SetConsensus(insolarConsensus consensus.InsolarConsensus) {
	fac.ic = insolarConsensus
}

func (fac *commonFacade) GetConsensus() consensus.InsolarConsensus {
	return fac.ic
}

// HostHandler is an interface which uses for host network implementation.
type HostHandler interface {
	Disconnect()
	Listen() error
	ObtainIP() error
	Bootstrap() error
	GetActiveNodes() error
	GetHostsFromBootstrap()
	NumHosts(ctx Context) int
	AnalyzeNetwork(ctx Context) error
	ConfirmNodeRole(role string) bool
	StoreRetrieve(key store.Key) ([]byte, bool)
	HtFromCtx(ctx Context) *routing.HashTable
	EqualAuthSentKey(targetID string, key []byte) bool
	SendRequest(packet *packet.Packet) (transport.Future, error)
	FindHost(ctx Context, targetID string) (*host.Host, bool, error)
	RemoteProcedureRegister(name string, method core.RemoteProcedure)
	InvokeRPC(sender *host.Host, method string, args [][]byte) ([]byte, error)
	CascadeSendMessage(data core.Cascade, targetID string, method string, args [][]byte) error
	Store(key store.Key, data []byte, replication time.Time, expiration time.Time, publisher bool) error
	RemoteProcedureCall(ctx Context, targetID string, method string, args [][]byte) (result []byte, err error)

	AddPossibleProxyID(id string)
	AddPossibleRelayID(id string)
	AddProxyHost(targetID string)
	AddSubnetID(ip, targetID string)
	AddAuthSentKey(id string, key []byte)
	AddRelayClient(host *host.Host) error
	AddReceivedKey(target string, key []byte)
	AddHost(ctx Context, host *routing.RouteHost)
	AddActiveNodes(activeNode []*core.ActiveNode) error

	RemoveAuthHost(key string)
	RemoveProxyHost(targetID string)
	RemovePossibleProxyID(id string)
	RemoveAuthSentKeys(targetID string)
	RemoveRelayClient(host *host.Host) error

	SetHighKnownHostID(id string)
	SetOuterHostsCount(hosts int)
	SetSignChecker(func(msg core.Message) bool)
	SetAuthStatus(targetID string, status bool)

	GetProxyHostsCount() int
	GetOuterHostsCount() int
	GetNodeID() core.RecordRef
	GetHighKnownHostID() string
	GetSelfKnownOuterHosts() int
	GetOriginHost() *host.Origin
	GetPacketTimeout() time.Duration
	GetReplicationTime() time.Duration
	HostIsAuthenticated(targetID string) bool
	KeyIsReceived(targetID string) ([]byte, bool)
	GetNetworkCommonFacade() NetworkCommonFacade
	GetExpirationTime(ctx Context, key []byte) time.Time
	GetActiveNodesList() []*core.ActiveNode
}

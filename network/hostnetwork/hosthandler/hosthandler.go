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

package hosthandler

import (
	"context"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
)

// Context type is localized for future purposes.
// Network Host can have multiple IDs, but each action must be executed with only one ID.
// Context is used in all actions to select specific ID to work with.
type Context context.Context

// HostHandler is an interface which uses for host network implementation.
type HostHandler interface {
	ConfirmNodeRole(role string) bool
	StoreRetrieve(key store.Key) ([]byte, bool)
	HtFromCtx(ctx Context) *routing.HashTable
	EqualAuthSentKey(targetID string, key []byte) bool
	SendRequest(packet *packet.Packet) (transport.Future, error)
	FindHost(ctx Context, targetID string) (*host.Host, bool, error)
	InvokeRPC(sender *host.Host, method string, args [][]byte) ([]byte, error)
	Store(key store.Key, data []byte, replication time.Time, expiration time.Time, publisher bool) error

	AddPossibleProxyID(id string)
	AddPossibleRelayID(id string)
	AddProxyHost(targetID string)
	AddSubnetID(ip, targetID string)
	AddAuthSentKey(id string, key []byte)
	AddRelayClient(host *host.Host) error
	AddReceivedKey(target string, key []byte)
	AddHost(ctx Context, host *routing.RouteHost)

	RemoveAuthHost(key string)
	RemoveProxyHost(targetID string)
	RemovePossibleProxyID(id string)
	RemoveAuthSentKeys(targetID string)
	RemoveRelayClient(host *host.Host) error

	SetHighKnownHostID(id string)
	SetOuterHostsCount(hosts int)
	SetAuthStatus(targetID string, status bool)

	GetProxyHostsCount() int
	GetOuterHostsCount() int
	GetSelfKnownOuterHosts() int
	GetHighKnownHostID() string
	GetPacketTimeout() time.Duration
	GetReplicationTime() time.Duration
	GetExpirationTime(ctx Context, key []byte) time.Time
	KeyIsReceived(targetID string) ([]byte, bool)
	HostIsAuthenticated(targetID string) bool
}

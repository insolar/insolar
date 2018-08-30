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

package nodenetwork

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/rpc"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/stretchr/testify/assert"
)

func realDhtParams(ids []id.ID, address string) (store.Store, *host.Origin, transport.Transport, rpc.RPC, error) {
	st := store.NewMemoryStore()
	addr, _ := host.NewAddress(address)
	origin, _ := host.NewOrigin(ids, addr)
	cfg := configuration.NewConfiguration().Host.Transport
	cfg.Address = address
	cfg.BehindNAT = false
	tp, err := transport.NewTransport(cfg, relay.NewProxy())
	r := rpc.NewRPC()
	return st, origin, tp, r, err
}

func TestNewNode(t *testing.T) {
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st, s, tp, r, err := realDhtParams(ids1, "127.0.0.1:16001")
	dht1, _ := hostnetwork.NewDHT(st, s, tp, r, &hostnetwork.Options{}, relay.NewProxy())
	assert.NoError(t, err)
	node := NewNode("id", nil, dht1, "domainID")

	assert.NotNil(t, node)
}

func TestNode_GetDomainID(t *testing.T) {
	expectedDomain := "domain id"
	node := Node{
		id:       "id",
		dht:      nil,
		domainID: expectedDomain,
	}

	assert.Equal(t, expectedDomain, node.GetDomainID())
}

func TestNode_GetNodeID(t *testing.T) {
	expectedID := "id"
	node := Node{
		domainID: "domain id",
		id:       expectedID,
		dht:      nil,
		host:     nil,
	}
	assert.Equal(t, expectedID, node.GetNodeID())
}

func TestNode_GetNodeRole(t *testing.T) {
	expectedRole := "role"
	node := Node{
		id:       "id",
		dht:      nil,
		domainID: "domain id",
		host:     nil,
	}

	node.setRole(expectedRole)
	assert.Equal(t, expectedRole, node.GetNodeRole())
}

func TestNode_SendPacket(t *testing.T) {
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st, s, tp, r, err := realDhtParams(ids1, "127.0.0.1:16002")
	dht1, _ := hostnetwork.NewDHT(st, s, tp, r, &hostnetwork.Options{}, relay.NewProxy())
	assert.NoError(t, err)
	node := NewNode("id", nil, dht1, "domainID")

	args := make([][]byte, 2)

	err = node.SendPacket("target", "method", args)
	assert.Error(t, err)
}

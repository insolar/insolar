package nodenetwork

import (
	"testing"

	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/connection"
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
	conn, _ := connection.NewConnectionFactory().Create(address)
	tp, err := transport.NewUTPTransport(conn, relay.NewProxy())
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
	node, err := NewNode(id1, dht1, "domainID", "role")

	assert.Error(t, err, "bootstrap node not exist")
	node.SetDHT(dht1)
}

func TestNode_GetDomainID(t *testing.T) {
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st, s, tp, r, err := realDhtParams(ids1, "127.0.0.1:16002")
	dht1, _ := hostnetwork.NewDHT(st, s, tp, r, &hostnetwork.Options{}, relay.NewProxy())
	assert.NoError(t, err)
	node := Node{id: id1,
		role:     "role",
		dht:      dht1,
		domainID: "domain id",
	}

	args := []struct {
		name     string
		expected string
		actual   string
	}{
		{"equal", "domain ID", "domain ID"},
		{"not equal", "domain ID", "another domain ID"},
		{"equal", "another domain ID", "another domain ID"},
	}
	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			node.SetDomainID(arg.actual)
			if arg.expected == arg.actual {
				assert.Equal(t, arg.expected, node.GetDomainID())
			} else {
				assert.NotEqual(t, arg.expected, node.GetDomainID())
			}
		})
	}
}

func TestNode_GetNodeID(t *testing.T) {
	id2, _ := id.NewID(id.GetRandomKey())
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st, s, tp, r, err := realDhtParams(ids1, "127.0.0.1:16003")
	dht1, _ := hostnetwork.NewDHT(st, s, tp, r, &hostnetwork.Options{}, relay.NewProxy())
	assert.NoError(t, err)
	node := Node{id: id1,
		role:     "role",
		dht:      dht1,
		domainID: "domain id",
	}

	args := []struct {
		name     string
		expected id.ID
		actual   id.ID
	}{
		{"equal", id1, id1},
		{"not equal", id2, id1},
		{"equal", id2, id2},
	}
	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			node.SetNodeID(arg.actual)
			if arg.expected.HashEqual(arg.actual.GetHash()) {
				assert.Equal(t, arg.expected.GetHash(), node.GetNodeID().GetHash())
			} else {
				assert.NotEqual(t, arg.expected.GetHash(), node.GetNodeID().GetHash())
			}
		})
	}
}

func TestNode_GetNodeRole(t *testing.T) {
	ids1 := make([]id.ID, 0)
	id1, _ := id.NewID(id.GetRandomKey())
	ids1 = append(ids1, id1)
	st, s, tp, r, err := realDhtParams(ids1, "127.0.0.1:16004")
	dht1, _ := hostnetwork.NewDHT(st, s, tp, r, &hostnetwork.Options{}, relay.NewProxy())
	assert.NoError(t, err)
	node := Node{id: id1,
		role:     "role",
		dht:      dht1,
		domainID: "domain id",
	}

	args := []struct {
		name     string
		expected string
		actual   string
	}{
		{"equal", "node role", "node role"},
		{"not equal", "node role", "another node role"},
		{"equal", "another node role", "another node role"},
	}
	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			node.SetNodeRole(arg.actual)
			if arg.expected == arg.actual {
				assert.Equal(t, arg.expected, node.GetNodeRole())
			} else {
				assert.NotEqual(t, arg.expected, node.GetNodeRole())
			}
		})
	}
}

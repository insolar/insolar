package nodenetwork

import (
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	address, _ := host.NewAddress("127.0.0.1:50001")
	host1 := host.NewHost(address)
	node := NewNode(id1, host1, "domainID", "role")

	assert.NotNil(t, node)
}

func TestNode_GetDomainID(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	address, _ := host.NewAddress("127.0.0.1:50001")
	node := Node{id: id1,
		role:     "role",
		host:     host.NewHost(address),
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

func TestNode_GetNodeHost(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	address, _ := host.NewAddress("127.0.0.1:50001")
	host1 := host.NewHost(address)
	address, _ = host.NewAddress("127.0.0.1:50002")
	host2 := host.NewHost(address)
	node := Node{id: id1,
		role:     "role",
		host:     host1,
		domainID: "domain id",
	}

	args := []struct {
		name     string
		expected *host.Host
		actual   *host.Host
	}{
		{"equal", host1, host1},
		{"not equal", host2, host1},
		{"equal", host2, host2},
	}
	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			node.SetHost(arg.actual)
			if arg.expected.ID.HashEqual(arg.actual.ID.GetHash()) {
				assert.Equal(t, arg.expected.ID.GetHash(), node.GetNodeHost().ID.GetHash())
			} else {
				assert.NotEqual(t, arg.expected.ID.GetHash(), node.GetNodeHost().ID.GetHash())
			}
		})
	}
}

func TestNode_GetNodeID(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id2, _ := id.NewID(id.GetRandomKey())
	address, _ := host.NewAddress("127.0.0.1:50001")
	node := Node{id: id1,
		role:     "role",
		host:     host.NewHost(address),
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
	id1, _ := id.NewID(id.GetRandomKey())
	address, _ := host.NewAddress("127.0.0.1:50001")
	node := Node{id: id1,
		role:     "role",
		host:     host.NewHost(address),
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

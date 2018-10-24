package nodekeeper

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/testutils/network"
)

func GetTestNodekeeper() consensus.NodeKeeper {
	nw := network.GetTestNetwork()
	keeper := nodekeeper.NewNodeKeeper(core.NewRefFromBase58("test"))

	// dirty hack - we need 3 nodes as validators, pass one node 3 times
	nodes := []*core.ActiveNode{
		{NodeID: nw.GetNodeID(), State: core.NodeActive, Roles: []core.NodeRole{core.RoleVirtual, core.RoleLightMaterial}},
		{NodeID: nw.GetNodeID(), State: core.NodeActive, Roles: []core.NodeRole{core.RoleVirtual, core.RoleLightMaterial}},
		{NodeID: nw.GetNodeID(), State: core.NodeActive, Roles: []core.NodeRole{core.RoleVirtual, core.RoleLightMaterial}},
	}
	keeper.AddActiveNodes(nodes)

	return keeper
}

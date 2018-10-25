package nodekeeper

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodekeeper"
	testNetwork "github.com/insolar/insolar/testutils/network"
)

func GetTestNodekeeper() network.NodeKeeper {
	nw := testNetwork.GetTestNetwork()
	keeper := nodekeeper.NewNodeKeeper(
		&core.Node{
			NodeID:   nw.GetNodeID(),
			PulseNum: 0,
			State:    core.NodeJoined,
			Roles:    []core.NodeRole{core.RoleVirtual, core.RoleHeavyMaterial, core.RoleLightMaterial},

			// TODO implement later
			//Address:  publicAddress,
			//Version:  version.Version,
		})

	// dirty hack - we need 3 nodes as validators, pass one node 3 times
	nodes := []*core.Node{
		{NodeID: nw.GetNodeID(), State: core.NodeActive, Roles: []core.NodeRole{core.RoleVirtual, core.RoleLightMaterial}},
		{NodeID: nw.GetNodeID(), State: core.NodeActive, Roles: []core.NodeRole{core.RoleVirtual, core.RoleLightMaterial}},
		{NodeID: nw.GetNodeID(), State: core.NodeActive, Roles: []core.NodeRole{core.RoleVirtual, core.RoleLightMaterial}},
	}
	keeper.AddActiveNodes(nodes)

	return keeper
}

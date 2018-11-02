package nodekeeper

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
	testNetwork "github.com/insolar/insolar/testutils/network"
)

func GetTestNodekeeper() network.NodeKeeper {
	nw := testNetwork.GetTestNetwork()
	keeper := nodenetwork.NewNodeKeeper(
		nodenetwork.NewNode(
			nw.GetNodeID(),
			[]core.NodeRole{core.RoleVirtual, core.RoleHeavyMaterial, core.RoleLightMaterial},
			nil,
			core.PulseNumber(0),
			core.NodeJoined,
			// TODO implement later
			"",
			"",
		))

	// dirty hack - we need 3 nodes as validators, pass one node 3 times
	getValidator := func() core.Node {
		return nodenetwork.NewNode(
			nw.GetNodeID(),
			[]core.NodeRole{core.RoleVirtual, core.RoleLightMaterial},
			nil,
			core.PulseNumber(0),
			core.NodeActive,
			// TODO implement later
			"",
			"",
		)
	}
	nodes := []core.Node{getValidator(), getValidator(), getValidator()}
	keeper.AddActiveNodes(nodes)

	return keeper
}

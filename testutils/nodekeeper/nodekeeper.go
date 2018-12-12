package nodekeeper

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
)

func GetTestNodekeeper(cs core.CryptographyService) network.NodeKeeper {
	pk, err := cs.GetPublicKey()
	if err != nil {
		panic(err)
	}

	keeper := nodenetwork.NewNodeKeeper(
		nodenetwork.NewNode(
			core.NewRefFromBase58("v1"),
			core.StaticRoleVirtual,
			pk,
			// TODO implement later
			"",
			"",
		))

	// dirty hack - we need 3 nodes as validators, pass one node 3 times
	getValidator := func() core.Node {
		return nodenetwork.NewNode(
			core.NewRefFromBase58("v1"),
			core.StaticRoleVirtual,
			pk,
			// TODO implement later
			"",
			"",
		)
	}
	nodes := []core.Node{getValidator(), getValidator(), getValidator()}
	keeper.AddActiveNodes(nodes)

	return keeper
}

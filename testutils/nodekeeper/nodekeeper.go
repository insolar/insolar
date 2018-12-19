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

	ref, err := core.NewRefFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	if err != nil {
		panic(err)
	}

	keeper := nodenetwork.NewNodeKeeper(
		nodenetwork.NewNode(
			*ref,
			core.StaticRoleVirtual,
			pk,
			// TODO implement later
			"",
			"",
		))

	// dirty hack - we need 3 nodes as validators, pass one node 3 times
	getValidator := func() core.Node {
		return nodenetwork.NewNode(
			*ref,
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

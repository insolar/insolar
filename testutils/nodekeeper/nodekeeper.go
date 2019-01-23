/*
 *    Copyright 2019 Insolar
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

// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package rules

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/node"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
)

func TestRules_CheckMinRole(t *testing.T) {
	cert := testutils.NewCertificateMock(t)
	nodes := []insolar.NetworkNode{
		node.NewNode(gen.Reference(), insolar.StaticRoleHeavyMaterial, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleLightMaterial, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleLightMaterial, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "", ""),
	}
	cert.GetMinRolesMock.Set(func() (r uint, r1 uint, r2 uint) {
		return 1, 0, 0
	})
	err := CheckMinRole(cert, nodes)
	require.NoError(t, err)

	cert.GetMinRolesMock.Set(func() (r uint, r1 uint, r2 uint) {
		return 3, 2, 4
	})
	err = CheckMinRole(cert, nodes)
	require.Error(t, err)
}

func TestRules_CheckMajorityRule(t *testing.T) {
	discNodesCount := 5
	netNodes, discoveryNodes := getDiscoveryNodes(discNodesCount)
	cert := testutils.NewCertificateMock(t)
	cert.GetDiscoveryNodesMock.Set(func() (r []insolar.DiscoveryNode) {
		return discoveryNodes
	})
	cert.GetMajorityRuleMock.Set(func() (r int) {
		return discNodesCount
	})
	count, err := CheckMajorityRule(cert, netNodes)
	require.NoError(t, err)

	require.Equal(t, discNodesCount, count)

	netNodes = netNodes[:len(netNodes)-len(netNodes)/2]
	count, err = CheckMajorityRule(cert, netNodes)
	require.Error(t, err)
	require.Equal(t, len(netNodes), count)
}

func getDiscoveryNodes(count int) ([]insolar.NetworkNode, []insolar.DiscoveryNode) {
	netNodes := make([]insolar.NetworkNode, count)
	discoveryNodes := make([]insolar.DiscoveryNode, count)
	for i := 0; i < count; i++ {
		n := newNode(gen.Reference(), i)
		d := certificate.NewBootstrapNode(nil, "", n.Address(), n.ID().String(), n.Role().String())
		netNodes[i] = n
		discoveryNodes[i] = d
	}
	return netNodes, discoveryNodes
}

func newNode(ref insolar.Reference, i int) insolar.NetworkNode {
	return node.NewNode(ref, insolar.AllStaticRoles[i%len(insolar.AllStaticRoles)], nil,
		"127.0.0.1:"+strconv.Itoa(30000+i), "")
}
